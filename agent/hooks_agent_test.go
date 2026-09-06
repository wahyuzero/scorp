package agent

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"scorp-agent/config"
	"scorp-agent/registry"
)

// registerHookProbeTool registers a deterministic probe tool and cleans up.
func registerHookProbeTool(t *testing.T) {
	t.Helper()
	registry.RegisterTool(registry.ToolDef{
		Name:     "hook_probe",
		Category: "test",
		Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
			return "probe-ran", true
		},
	})
	t.Cleanup(func() { registry.UnregisterTool("hook_probe") })
}

func setHookEnvAgent(t *testing.T, pre, post string) {
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

func TestExecuteToolPreHookBlocks(t *testing.T) {
	registerHookProbeTool(t)
	setHookEnvAgent(t, "hook_probe:echo 'company policy: probe forbidden' >&2 && exit 2", "")

	out, ok := ExecuteTool(ToolCall{Name: "hook_probe", Args: map[string]interface{}{}}, 1)
	if ok {
		t.Fatalf("hook must block the tool, got ok=true out=%q", out)
	}
	if !strings.Contains(out, "🪝") || !strings.Contains(out, "company policy") {
		t.Errorf("expected 🪝 block with hook stderr as reason, got %q", out)
	}
	if strings.Contains(out, "probe-ran") {
		t.Error("tool must not have executed when a PreToolUse hook blocks")
	}
}

func TestExecuteToolHookContextAppended(t *testing.T) {
	registerHookProbeTool(t)
	// Pre-hook stdout = additional context; post-hook stdout likewise.
	setHookEnvAgent(t, "hook_probe:echo pre-context-here", "hook_probe:echo post-context-here")

	out, ok := ExecuteTool(ToolCall{Name: "hook_probe", Args: map[string]interface{}{}}, 1)
	if !ok {
		t.Fatalf("expected success, got %q", out)
	}
	if !strings.Contains(out, "probe-ran") {
		t.Errorf("missing tool output: %q", out)
	}
	if !strings.Contains(out, "pre-context-here") || !strings.Contains(out, "post-context-here") {
		t.Errorf("expected pre+post hook context appended, got %q", out)
	}
}

func TestExecuteToolHookScopedToOtherToolDoesNotFire(t *testing.T) {
	registerHookProbeTool(t)
	marker := filepath.Join(t.TempDir(), "fired.txt")
	setHookEnvAgent(t, "shell:echo fired >> "+marker, "")

	out, ok := ExecuteTool(ToolCall{Name: "hook_probe", Args: map[string]interface{}{}}, 1)
	if !ok || !strings.Contains(out, "probe-ran") {
		t.Fatalf("probe should run untouched, got ok=%v out=%q", ok, out)
	}
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Error("hook scoped to shell must not fire for hook_probe")
	}
}
