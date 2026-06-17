package main

import "log"

// ──────────────────────────────────────────────
// Phase 7 — Register Autonomous Tools
// ──────────────────────────────────────────────

func initAutonomous() {
	loadAutonomousConfig()
	loadAutoLog()
	log.Printf("[autonomous] Config loaded: enabled=%v interval=%v voice=%s approval=%s (cycles=%d actions=%d)",
		autoConfig.Enabled, autoConfig.Interval, autoConfig.VoiceMode,
		autoConfig.ApprovalLevel, autoConfig.TotalCycles, autoConfig.TotalActions)

	registerTool(ToolDef{
		Name:        "autonomous",
		Description: "Control the autonomous agent: enable/disable autonomous VPS management, trigger manual cycles, configure intervals, view audit log and action stats.",
		Category:    "autonomous",
		Native:      true,
		Arguments: map[string]ArgDef{
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
			"voice": {
				Type:        "string",
				Description: "Voice alert mode: off, important, always. Used with action=config.",
				Required:    false,
				Enum:        []string{"off", "important", "always"},
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
		Execute: executeAutonomous,
	})
}
