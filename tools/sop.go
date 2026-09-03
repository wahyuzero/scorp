package tools

import (
	"fmt"
	"strings"

	"scorp-agent/internal/helpers"
	"scorp-agent/sop"
)

// ──────────────────────────────────────────────
// SOP Tool — Run and inspect Standard Operating Procedures
// ──────────────────────────────────────────────

// ExecuteSOP handles listing, reading, and running SOP playbooks
func ExecuteSOP(args map[string]interface{}) (string, bool) {
	action := helpers.GetStringArg(args, "action", "list")
	name := helpers.GetStringArg(args, "name", "")

	switch action {
	case "list":
		sops := sop.ListSOPs()
		if len(sops) == 0 {
			return "No SOPs found. Built-in SOPs can be created with action='init'.", true
		}
		var sb strings.Builder
		sb.WriteString("📋 Standard Operating Procedures (SOPs):\n\n")
		for _, s := range sops {
			sb.WriteString(fmt.Sprintf("• **%s** — %s\n", s.Name, s.Description))
			for _, step := range s.Steps {
				sb.WriteString(fmt.Sprintf("    %s\n", step))
			}
			sb.WriteString("\n")
		}
		return sb.String(), true

	case "get", "view":
		if name == "" {
			return "Error: 'name' is required for action='get'", false
		}
		s, err := sop.GetSOP(name)
		if err != nil {
			return fmt.Sprintf("Error: %v", err), false
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("📋 SOP: %s\n", s.Name))
		sb.WriteString(fmt.Sprintf("Description: %s\n\n", s.Description))
		sb.WriteString("Steps:\n")
		for _, st := range s.Steps {
			sb.WriteString(fmt.Sprintf("  %s\n", st))
		}
		sb.WriteString(fmt.Sprintf("\nPrompt template:\n%s\n", s.Prompt))
		return sb.String(), true

	case "init":
		sop.InitDefaultSOPs()
		return "✅ Default SOP templates initialized in ~/.scorp/sops/", true

	default:
		return fmt.Sprintf("Unknown SOP action '%s'. Use 'list', 'get', or 'init'.", action), false
	}
}
