package tools

import (
	"fmt"
	"scorp-agent/registry"
	"scorp-agent/skills"
)

// init registers the activate_skill tool for on-demand Level-2 progressive disclosure
func init() {
	registry.RegisterTool(registry.ToolDef{
		Name:        "activate_skill",
		Description: "Load and activate specialized skill instructions/playbook into context for complex tasks (e.g. docker-ops, vps-devops, git-workflow, golang-pro). Call this when a task matches an available skill.",
		Category:    "skills",
		Native:      true,
		Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
			name, _ := args["name"].(string)
			if name == "" {
				return "Error: 'name' argument is required", false
			}

			body, err := skills.ActivateSkill(name, 5) // Active for 5 turns
			if err != nil {
				return fmt.Sprintf("Error activating skill '%s': %v", name, err), false
			}

			preview := body
			if len(preview) > 300 {
				preview = preview[:300] + "\n...[full instructions loaded into context]..."
			}

			return fmt.Sprintf("✅ Skill '%s' successfully activated for 5 turns!\n\n%s", name, preview), true
		},
		Arguments: map[string]registry.ArgDef{
			"name": {Type: "string", Description: "The exact name of the skill to activate (e.g. docker-ops, vps-devops, git-workflow, golang-pro)", Required: true},
		},
	})
}
