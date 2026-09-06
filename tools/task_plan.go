package tools

import (
	"scorp-agent/registry"
)

// init registers the task_plan schema. Execution is intercepted by the agent
// loop so the ledger is scoped to the SESSION (the tool layer only sees the
// numeric chatID, which is not unique across CLI sessions).
func init() {
	registry.RegisterTool(registry.ToolDef{
		Name:        "task_plan",
		Description: "Create or update the structured plan for the current request. MANDATORY FIRST STEP for every multi-step task: decompose the request into concrete verifiable steps. Mark an item 'in_progress' when you start it and 'done' ONLY after verifying its result with real tool output. The runtime REJECTS complete_task while any item is unfinished.",
		Category:    "core",
		Native:      true,
		Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
			return "Handled by the agent runtime loop.", true
		},
		Arguments: map[string]registry.ArgDef{
			"action": {Type: "string", Description: "'create' a new plan or 'update' one item's status", Required: true},
			"goal":   {Type: "string", Description: "action=create: the user's original request in one sentence"},
			"items":  {Type: "array", Description: "action=create: ordered step descriptions, e.g. [\"write script\", \"run and verify\", \"save report\"]"},
			"item":   {Type: "string", Description: "action=update: item id to update, e.g. '2'"},
			"status": {Type: "string", Description: "action=update: pending | in_progress | done"},
		},
	})
}
