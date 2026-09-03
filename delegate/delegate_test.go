package delegate

import (
	"scorp-agent/registry"
	"testing"
)

func TestValidateSubagentTools(t *testing.T) {
	// Register dummy test tools
	registry.RegisterTool(registry.ToolDef{Name: "delegate"})
	registry.RegisterTool(registry.ToolDef{Name: "read_file"})
	registry.RegisterTool(registry.ToolDef{Name: "search_code"})

	// Standard subagent should not have access to delegate
	tools := []string{"read_file", "search_code", "delegate", "non_existent_tool"}
	validated := ValidateSubagentTools(tools, false)

	for _, tool := range validated {
		if tool == "delegate" {
			t.Errorf("Standard subagent should not have access to delegate")
		}
		if tool == "non_existent_tool" {
			t.Errorf("Non-existent tool should be filtered out")
		}
	}

	// Orchestrator subagent can re-delegate
	orchestratorTools := ValidateSubagentTools(tools, true)
	hasDelegate := false
	for _, tool := range orchestratorTools {
		if tool == "delegate" {
			hasDelegate = true
			break
		}
	}
	if !hasDelegate {
		t.Errorf("Orchestrator subagent should have access to delegate")
	}

	// Cleanup
	registry.UnregisterTool("delegate")
	registry.UnregisterTool("read_file")
	registry.UnregisterTool("search_code")
}

func TestParseDelegateParams(t *testing.T) {
	args := map[string]interface{}{
		"task":      "Analyze logs",
		"role":      "research",
		"max_iters": 25, // Should clamp to MaxSubagentIters (15)
	}

	params := ParseDelegateParams(args)
	if params.Task != "Analyze logs" {
		t.Errorf("Expected task 'Analyze logs', got '%s'", params.Task)
	}
	if params.Role != "research" {
		t.Errorf("Expected role 'research', got '%s'", params.Role)
	}
	if params.MaxIters != MaxSubagentIters {
		t.Errorf("Expected max_iters clamped to %d, got %d", MaxSubagentIters, params.MaxIters)
	}
}
