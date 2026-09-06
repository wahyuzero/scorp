package agent

import (
	"strings"
	"testing"

	"scorp-agent/config"
	"scorp-agent/registry"
	"scorp-agent/tools"
)

func setAutoMode(t *testing.T, level string) {
	t.Helper()
	orig := string(config.GetAutonomyLevel())
	config.SetAutonomyLevel(level)
	t.Cleanup(func() { config.SetAutonomyLevel(orig) })
	ResetAutoStats()
	ResetAutoAllowlist()
	t.Cleanup(ResetAutoStats)
	t.Cleanup(ResetAutoAllowlist)
}

func TestIsReadOnlyShellCommand(t *testing.T) {
	readOnly := []string{
		"ls -la",
		"cat /etc/hostname",
		"cat app.log | grep error | wc -l",
		"git status",
		"git log --oneline -5",
		"docker ps",
		"systemctl is-active scorp",
		"go version",
		"echo hello",
		"find . -name '*.go' -type f",
		"df -h && free -m",
	}
	for _, cmd := range readOnly {
		if !IsReadOnlyShellCommand(cmd) {
			t.Errorf("expected read-only: %q", cmd)
		}
	}
	mutating := []string{
		"echo hi > /tmp/out.txt",
		"cat x >> y",
		"rm -rf /tmp/x",
		"touch /tmp/newfile",
		"mkdir /tmp/d",
		"git push origin master",
		"git checkout main",
		"curl https://example.com",
		"apt install htop",
		"go test ./...",
		"sed -i 's/a/b/' f.txt",
		"docker rm x",
	}
	for _, cmd := range mutating {
		if IsReadOnlyShellCommand(cmd) {
			t.Errorf("expected mutating: %q", cmd)
		}
	}
}

func TestPermissionDecisionReadOnlyFastPath(t *testing.T) {
	setAutoMode(t, "auto")
	decision, _, source := PermissionDecision("read_file", map[string]interface{}{"path": "/tmp/x"})
	if decision != AutoAllow || source != "deterministic" {
		t.Errorf("read_file in auto: got %s/%s", decision, source)
	}
}

func TestPermissionDecisionDestructiveDenyAndAllowlist(t *testing.T) {
	setAutoMode(t, "auto")

	decision, reason, _ := PermissionDecision("shell", map[string]interface{}{"command": "rm -rf /tmp/eval-auto-deny"})
	if decision != AutoDeny || !strings.Contains(reason, "destructive") {
		t.Errorf("destructive must deny, got %s: %s", decision, reason)
	}

	t.Setenv("SCORP_AUTO_ALLOW", "rm -rf /tmp/eval-allowed")
	ResetAutoAllowlist()
	decision, _, source := PermissionDecision("shell", map[string]interface{}{"command": "rm -rf /tmp/eval-allowed/cache"})
	if decision != AutoAllow || source != "allowlist" {
		t.Errorf("allowlisted destructive must allow, got %s/%s", decision, source)
	}
	decision, _, _ = PermissionDecision("shell", map[string]interface{}{"command": "rm -rf /tmp/other"})
	if decision != AutoDeny {
		t.Errorf("non-allowlisted must still deny, got %s", decision)
	}
}

func TestPermissionDecisionModelPaths(t *testing.T) {
	setAutoMode(t, "auto")
	orig := AutoClassifyFunc
	t.Cleanup(func() { AutoClassifyFunc = orig })

	cases := []struct {
		stub string
		want string
	}{
		{AutoAllow, AutoAllow},
		{AutoAsk, AutoAsk},
		{AutoDeny, AutoDeny},
	}
	for _, c := range cases {
		AutoClassifyFunc = func(string, map[string]interface{}) (string, string) {
			return c.stub, "eval stub"
		}
		decision, _, source := PermissionDecision("shell", map[string]interface{}{"command": "apt install htop"})
		if decision != c.want || source != "model" {
			t.Errorf("stub %s: got %s/%s want %s/model", c.stub, decision, source, c.want)
		}
	}
}

func TestAutoFallbackAfterRepeatedUncertain(t *testing.T) {
	setAutoMode(t, "auto")
	orig := AutoClassifyFunc
	t.Cleanup(func() { AutoClassifyFunc = orig })
	AutoClassifyFunc = func(string, map[string]interface{}) (string, string) {
		return "uncertain", "classifier broken"
	}

	for i := 0; i < autoMaxUncertain; i++ {
		decision, _, source := PermissionDecision("shell", map[string]interface{}{"command": "apt install htop"})
		if decision != AutoAsk || source != "fallback" {
			t.Fatalf("uncertain #%d must fail closed to ask/fallback, got %s/%s", i+1, decision, source)
		}
	}
	if _, _, _, _, fallback := AutoStatsSnapshot(); !fallback {
		t.Fatal("fallback must be active after 3 consecutive uncertain results")
	}
	// A mutating call while degraded still asks, without calling the model.
	calls := 0
	AutoClassifyFunc = func(string, map[string]interface{}) (string, string) {
		calls++
		return AutoAllow, "should not be consulted while degraded"
	}
	if decision, _, _ := PermissionDecision("shell", map[string]interface{}{"command": "apt install htop"}); decision != AutoAsk {
		t.Errorf("degraded mode must ask, got %s", decision)
	}
	if calls != 0 {
		t.Errorf("degraded mode must not call the model, got %d calls", calls)
	}
}

func TestExecuteToolAutoDeniesOnNoChannelPath(t *testing.T) {
	setAutoMode(t, "auto")

	// Deterministic destructive deny via the subagent/direct path.
	out, ok := ExecuteTool(ToolCall{Name: "shell", Args: map[string]interface{}{"command": "rm -rf /tmp/eval-auto-deny"}}, 0)
	if ok || !strings.Contains(out, "Denied by auto-mode classifier") {
		t.Errorf("expected auto-deny, got ok=%v out=%q", ok, out)
	}

	// Model "ask" degrades to deny on the no-confirmation-channel path.
	orig := AutoClassifyFunc
	t.Cleanup(func() { AutoClassifyFunc = orig })
	AutoClassifyFunc = func(string, map[string]interface{}) (string, string) {
		return AutoAsk, "eval stub risky"
	}
	out, ok = ExecuteTool(ToolCall{Name: "shell", Args: map[string]interface{}{"command": "apt install htop"}}, 0)
	if ok || !strings.Contains(out, "no confirmation channel") {
		t.Errorf("expected fail-closed deny, got ok=%v out=%q", ok, out)
	}
}

func TestExecuteToolAutoPresetTrustedAndRecorded(t *testing.T) {
	setAutoMode(t, "auto")

	// Tests don't run bootstrap: register a deterministic shell stub.
	registry.RegisterTool(registry.ToolDef{
		Name:     "shell",
		Category: "test",
		Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
			return "auto-preset-ok", true
		},
	})
	t.Cleanup(func() { registry.UnregisterTool("shell") })

	tc := ToolCall{Name: "shell", Args: map[string]interface{}{"command": "echo auto-preset-ok"}, AutoDecision: "auto:heuristic"}
	out, ok := ExecuteTool(tc, 0)
	if !ok || !strings.Contains(out, "auto-preset-ok") {
		t.Fatalf("preset allow must execute, got ok=%v out=%q", ok, out)
	}
	receipts := tools.GetRecentReceipts()
	last := receipts[len(receipts)-1]
	if last.Meta["auto_decision"] != "auto:heuristic" {
		t.Errorf("receipt must carry auto_decision, got meta=%v", last.Meta)
	}
}

func TestStorePendingConfirmationArgs(t *testing.T) {
	chatID := "auto-confirm-test"
	StorePendingConfirmationArgs(chatID, "write_file", "write_file {path:...}", map[string]interface{}{"path": "/tmp/x"}, nil)
	defer clearPendingConfirmation(chatID)

	pc := getPendingConfirmation(chatID)
	if pc == nil || pc.args == nil || pc.args["path"] != "/tmp/x" || pc.toolName != "write_file" {
		t.Fatalf("args confirmation not stored: %+v", pc)
	}
}
