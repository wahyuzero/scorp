package tools

import (
	"scorp-agent/registry"
)

// init registers the complete_task termination tool.
// In Agent Mode, calling this tool is the explicit, deterministic contract
// to signal that all sub-tasks and verifications are 100% finished.
func init() {
	registry.RegisterTool(registry.ToolDef{
		Name:        "complete_task",
		Description: "Call this tool ONLY when you have fully completed and verified the user's task. Output the final conclusive answer, report, or verification result in the 'result' argument. This officially concludes the agent execution loop.",
		Category:    "core",
		Native:      true,
		Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
			result, _ := args["result"].(string)
			if result == "" {
				result, _ = args["summary"].(string)
			}
			if result == "" {
				return "Task marked as completed.", true
			}
			return result, true
		},
		Arguments: map[string]registry.ArgDef{
			"result": {Type: "string", Description: "The complete, polished final response, answer, or report to present to the user.", Required: true},
		},
	})
}
