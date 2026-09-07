package agent

import (
	"os"
	"strings"
	"testing"

	"scorp-agent/config"
	"scorp-agent/registry"
	"scorp-agent/tools"
)

// TestAutoAllowlistPrefixAnchoring pins the A3/F-5 fix: an allowlist entry
// matches the allowlisted path and its SUBPATHS, but not prefix-adjacent
// paths that merely start with the same text.
func TestAutoAllowlistPrefixAnchoring(t *testing.T) {
	t.Setenv("SCORP_AUTO_ALLOW", "rm -rf /tmp/allowed-only")
	ResetAutoAllowlist()
	defer ResetAutoAllowlist()

	cases := []struct {
		cmd  string
		want bool
	}{
		{"rm -rf /tmp/allowed-only", true},
		{"rm -rf /tmp/allowed-only/a", true},
		{"rm -rf /tmp/allowed-only/sub/dir", true},
		{"rm -rf /tmp/allowed-onlyx/b", false},     // prefix-adjacent leak
		{"rm -rf /tmp/allowed-only-backup", false}, // dash-suffixed sibling
		{"echo rm -rf /tmp/allowed-only", false},   // not a command prefix
		{"rm -rf /somewhere/else", false},
	}
	for _, tc := range cases {
		if got := autoAllowlisted(tc.cmd); got != tc.want {
			t.Errorf("autoAllowlisted(%q) = %v, want %v", tc.cmd, got, tc.want)
		}
	}
}

// TestExecuteToolAutoAllowlistReachesExec pins the A3/F-4 fix: in auto mode a
// classifier-ALLOW decision (destructive but allowlisted) must reach real
// execution. Before the fix the exec-layer dangerous gate re-blocked what
// auto mode deliberately allowed, making SCORP_AUTO_ALLOW dead code.
func TestExecuteToolAutoAllowlistReachesExec(t *testing.T) {
	defer config.SetAutonomyLevel(string(config.GetAutonomyLevel()))
	config.SetAutonomyLevel("auto")

	t.Setenv("SCORP_AUTO_ALLOW", "rm -rf /tmp/agent-allowlist-test")
	ResetAutoAllowlist()
	defer ResetAutoAllowlist()
	defer ResetAutoStats()

	target := "/tmp/agent-allowlist-test/x"
	os.RemoveAll(target)
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(target+"/f", []byte("1"), 0o644); err != nil {
		t.Fatal(err)
	}

	registry.RegisterTool(registry.ToolDef{
		Name:     "shell",
		Category: "test",
		Execute:  tools.ExecuteShell,
	})
	defer registry.UnregisterTool("shell")

	tc := ToolCall{Name: "shell", Args: map[string]interface{}{"command": "rm -rf " + target}}
	out, ok := ExecuteTool(tc, 0)
	if !ok {
		t.Fatalf("allowlisted destructive must execute in auto mode, out=%.300q", out)
	}
	if _, err := os.Stat(target); !os.IsNotExist(err) {
		t.Fatal("allowlisted target directory should be deleted")
	}
}

// TestExecuteToolAutoDenyNotBypassedByConfirmedArgs keeps the A2 guarantee
// intact after the allowlist fix: a model-supplied confirmed:true arg must
// NOT bypass the deterministic auto-mode deny.
func TestExecuteToolAutoDenyNotBypassedByConfirmedArgs(t *testing.T) {
	defer config.SetAutonomyLevel(string(config.GetAutonomyLevel()))
	config.SetAutonomyLevel("auto")
	defer ResetAutoStats()

	registry.RegisterTool(registry.ToolDef{
		Name:     "shell",
		Category: "test",
		Execute:  tools.ExecuteShell,
	})
	defer registry.UnregisterTool("shell")

	tc := ToolCall{Name: "shell", Args: map[string]interface{}{
		"command":   "rm -rf /tmp/agent-deny-no-bypass",
		"confirmed": true, // forged by the model — must be ignored
	}}
	out, ok := ExecuteTool(tc, 0)
	if ok {
		t.Fatalf("model-supplied confirmed:true must not bypass auto deny, out=%.200q", out)
	}
	if !strings.Contains(out, "Denied by auto-mode classifier") {
		t.Fatalf("expected classifier deny, got %.200q", out)
	}
}
