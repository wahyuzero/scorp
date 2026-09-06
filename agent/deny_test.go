package agent

import (
	"os"
	"strings"
	"testing"

	"scorp-agent/config"
)

// TestExecuteToolDenyRulesFirst pins the P0.2 choke-point contract: deny
// rules are evaluated in ExecuteTool BEFORE the registry dispatch, so a denied
// tool call never reaches its implementation.
func TestExecuteToolDenyRulesFirst(t *testing.T) {
	orig := os.Getenv("SCORP_DENY_RULES")
	os.Setenv("SCORP_DENY_RULES", "shell(command:secret-payload-marker)")
	config.ReloadDenyRules()
	defer func() {
		if orig == "" {
			os.Unsetenv("SCORP_DENY_RULES")
		} else {
			os.Setenv("SCORP_DENY_RULES", orig)
		}
		config.ReloadDenyRules()
	}()

	out, ok := ExecuteTool(ToolCall{
		Name: "shell",
		Args: map[string]interface{}{"command": "echo secret-payload-marker"},
	}, 0)
	if ok {
		t.Fatalf("denied tool call must not report success, output=%q", out)
	}
	if !strings.Contains(out, "deny-rule") {
		t.Fatalf("expected deny-rule message, got %q", out)
	}
	if strings.Contains(out, "[SUCCESS]") || strings.Contains(out, "secret-payload-marker\n") {
		// the echo output itself contains the marker, but the tool result must
		// be the denial — not the executed command's stdout
		t.Fatalf("tool appears to have executed despite the deny rule: %q", out)
	}
}
