package bootstrap

import (
	"scorp-agent/registry"
	"scorp-agent/tools"
)

// RegisterAutonomous registers the autonomous agent tool.
// Config loading must be done by the caller (main package) before calling this.
func RegisterAutonomous() {
	registry.RegisterTool(registry.ToolDef{
		Name:        "autonomous",
		Description: "Autonomous agent control: enable/disable, manual trigger, configure interval, view audit log and action stats.",
		Category:    "autonomous",
		Native:      true,
		Arguments: map[string]registry.ArgDef{
			"action": {
				Type:        "string",
				Description: "Action to perform: status, enable, disable, kill, revive, run, config, log, actions",
				Required:    false,
				Default:     "status",
				Enum:        []string{"status", "enable", "disable", "kill", "revive", "run", "config", "log", "actions"},
			},
			"interval": {
				Type:        "string",
				Description: "Interval between cycles (e.g. 5m, 10m, 1h). Used with action=config.",
				Required:    false,
			},
			"approval": {
				Type:        "string",
				Description: "Approval level: low (auto-execute all), medium (approve high risk), high (approve all). Used with action=config.",
				Required:    false,
				Enum:        []string{"low", "medium", "high"},
			},
			"max_actions": {
				Type:        "integer",
				Description: "Max actions per cycle (1-10). Used with action=config.",
				Required:    false,
			},
			"count": {
				Type:        "integer",
				Description: "Number of log entries to show. Used with action=log.",
				Required:    false,
				Default:     10,
			},
		},
		Execute: tools.ExecuteAutonomous,
	})
}
