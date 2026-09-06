package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"scorp-agent/config"
)

// setHookEnv installs hook env vars and restores a clean hook state afterwards.
func setHookEnv(t *testing.T, pre, post string) {
	t.Helper()
	if pre != "" {
		os.Setenv("SCORP_HOOKS_PRE", pre)
	} else {
		os.Unsetenv("SCORP_HOOKS_PRE")
	}
	if post != "" {
		os.Setenv("SCORP_HOOKS_POST", post)
	} else {
		os.Unsetenv("SCORP_HOOKS_POST")
	}
	config.ReloadHooks()
	t.Cleanup(func() {
		os.Unsetenv("SCORP_HOOKS_PRE")
		os.Unsetenv("SCORP_HOOKS_POST")
		config.ReloadHooks()
	})
}

func TestPreToolUseHookExit2Blocks(t *testing.T) {
	setHookEnv(t, "shell:echo 'policy violation: no rm -rf' >&2 && exit 2", "")
	ctx, blocked, reason := RunPreToolUseHooks("shell", map[string]interface{}{"command": "rm -rf /tmp/x"}, 42)
	if !blocked {
		t.Fatalf("expected block, got context=%q", ctx)
	}
	if !strings.Contains(reason, "policy violation") {
		t.Errorf("expected stderr as reason, got %q", reason)
	}
}

func TestPreToolUseHookExit2DefaultReason(t *testing.T) {
	setHookEnv(t, "shell:exit 2", "")
	_, blocked, reason := RunPreToolUseHooks("shell", nil, 42)
	if !blocked {
		t.Fatal("expected block")
	}
	if !strings.Contains(reason, "blocked by PreToolUse hook") {
		t.Errorf("expected default reason, got %q", reason)
	}
}

func TestPreToolUseHookExit0AddsContext(t *testing.T) {
	setHookEnv(t, "*:echo advisory-check-ok", "")
	ctx, blocked, _ := RunPreToolUseHooks("read_file", nil, 42)
	if blocked {
		t.Fatal("exit 0 must not block")
	}
	if !strings.Contains(ctx, "advisory-check-ok") {
		t.Errorf("expected stdout as context, got %q", ctx)
	}
}

func TestPreToolUseHookNonZeroNonTwoIsNonBlocking(t *testing.T) {
	setHookEnv(t, "*:echo oops >&2 && exit 1", "")
	ctx, blocked, _ := RunPreToolUseHooks("read_file", nil, 42)
	if blocked {
		t.Fatal("exit 1 must be non-blocking")
	}
	if strings.Contains(ctx, "oops") {
		t.Errorf("stderr of non-blocking failure must not leak into context, got %q", ctx)
	}
}

func TestPreToolUseHookTimeoutIsNonBlocking(t *testing.T) {
	old := hookRunTimeout
	hookRunTimeout = 300 * time.Millisecond
	t.Cleanup(func() { hookRunTimeout = old })

	setHookEnv(t, "shell:sleep 30", "")
	ctx, blocked, _ := RunPreToolUseHooks("shell", nil, 42)
	if blocked {
		t.Fatal("timed-out hook must be non-blocking (a broken hook must not brick the agent)")
	}
	if !strings.Contains(ctx, "timed out") {
		t.Errorf("expected timeout warning in context, got %q", ctx)
	}
}

func TestHookMatcherFiltersTools(t *testing.T) {
	marker := filepath.Join(t.TempDir(), "fired.txt")
	setHookEnv(t, fmt.Sprintf("shell:touch %s", marker), "")
	RunPreToolUseHooks("read_file", nil, 42)
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatal("hook scoped to 'shell' must not fire for read_file")
	}
	RunPreToolUseHooks("shell", nil, 42)
	if _, err := os.Stat(marker); err != nil {
		t.Fatalf("hook scoped to 'shell' did not fire: %v", err)
	}
}

func TestHookPayloadFieldsAndRedaction(t *testing.T) {
	dir := t.TempDir()
	payloadFile := filepath.Join(dir, "payload.json")
	cmd := fmt.Sprintf(`cat > %s`, payloadFile)
	setHookEnv(t, "shell:"+cmd, "")

	t.Setenv("SCORP_TEST_FAKE_KEY", "sk-fake-secret-value-1234567890")
	secret := os.Getenv("SCORP_TEST_FAKE_KEY")

	args := map[string]interface{}{"command": "echo hello", "extra": secret}
	RunPreToolUseHooks("shell", args, 777)

	raw, err := os.ReadFile(payloadFile)
	if err != nil {
		t.Fatalf("hook did not receive stdin payload: %v", err)
	}
	if strings.Contains(string(raw), secret) {
		t.Errorf("payload leaked the secret value: %s", raw)
	}
	var p map[string]interface{}
	if err := json.Unmarshal(raw, &p); err != nil {
		t.Fatalf("payload is not valid JSON: %v\n%s", err, raw)
	}
	if p["hook_event"] != "pre_tool_use" || p["tool_name"] != "shell" || p["session_id"] != "777" {
		t.Errorf("payload fields wrong: %v", p)
	}
	ta, ok := p["tool_args"].(map[string]interface{})
	if !ok || ta["command"] != "echo hello" {
		t.Errorf("tool_args wrong: %v", p["tool_args"])
	}
}

func TestPostToolUseHookAddsContextAndNeverBlocks(t *testing.T) {
	dir := t.TempDir()
	payloadFile := filepath.Join(dir, "post.json")
	// tee echoes the payload to stdout (becomes hook context) and stores it.
	setHookEnv(t, "", "*:tee "+payloadFile)

	ctx := RunPostToolUseHooks("write_file", map[string]interface{}{"path": "/tmp/x"}, "file written", true, 42)
	if !strings.Contains(ctx, "post_tool_use") {
		t.Errorf("expected post hook stdout as context, got %q", ctx)
	}
	raw, _ := os.ReadFile(payloadFile)
	var p map[string]interface{}
	if err := json.Unmarshal(raw, &p); err != nil {
		t.Fatalf("post payload invalid: %v", err)
	}
	if p["hook_event"] != "post_tool_use" || p["tool_result"] != "file written" || p["tool_ok"] != true {
		t.Errorf("post payload fields wrong: %v", p)
	}
}

func TestPostToolUseHookExit2IsAdvisoryOnly(t *testing.T) {
	setHookEnv(t, "", "*:echo policy-flag >&2 && exit 2")
	ctx := RunPostToolUseHooks("shell", nil, "done", true, 42)
	if !strings.Contains(ctx, "flagged this result") {
		t.Errorf("expected advisory flag context, got %q", ctx)
	}
}

func TestNoHooksConfiguredIsNoop(t *testing.T) {
	setHookEnv(t, "", "")
	if ctx, blocked, _ := RunPreToolUseHooks("shell", nil, 42); blocked || ctx != "" {
		t.Errorf("no hooks configured: got ctx=%q blocked=%v", ctx, blocked)
	}
	if ctx := RunPostToolUseHooks("shell", nil, "r", true, 42); ctx != "" {
		t.Errorf("no hooks configured: got ctx=%q", ctx)
	}
}
