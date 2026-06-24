package tools

import (
	"scorp-agent/internal/helpers"
	"fmt"
	"os/exec"
	"strings"
)

// ──────────────────────────────────────────────
// Docker Compose Tool
// ──────────────────────────────────────────────

func ExecuteCompose(args map[string]interface{}) (string, bool) {
	action := helpers.GetStringArg(args, "action", "up")
	project := helpers.GetStringArg(args, "project", ".")
	composeFile := helpers.GetStringArg(args, "file", "")

	// Safety: never delete containers permanently
	if action == "down" {
		if !helpers.GetBoolArg(args, "confirm", false) {
			return "⚠️ Use `confirm=true` to run `compose down`. This stops all services in the project.", true
		}
	}

	// Build docker compose command
	cmdArgs := []string{"compose"}

	if composeFile != "" {
		cmdArgs = append(cmdArgs, "-f", composeFile)
	}

	switch action {
	case "up":
		detach := helpers.GetBoolArg(args, "detach", true)
		rebuild := helpers.GetBoolArg(args, "rebuild", false)
		cmdArgs = append(cmdArgs, "up")
		if detach {
			cmdArgs = append(cmdArgs, "-d")
		}
		if rebuild {
			cmdArgs = append(cmdArgs, "--build")
		}
	case "down":
		volumes := helpers.GetBoolArg(args, "volumes", false)
		removeOrphans := helpers.GetBoolArg(args, "remove_orphans", true)
		cmdArgs = append(cmdArgs, "down")
		if volumes {
			cmdArgs = append(cmdArgs, "-v")
		}
		if removeOrphans {
			cmdArgs = append(cmdArgs, "--remove-orphans")
		}
	case "restart":
		services := helpers.GetStringArg(args, "services", "")
		timeout := helpers.GetIntArg(args, "timeout", 10)
		cmdArgs = append(cmdArgs, "restart", "-t", fmt.Sprintf("%d", timeout))
		if services != "" {
			cmdArgs = append(cmdArgs, strings.Fields(services)...)
		}
	case "logs":
		services := helpers.GetStringArg(args, "services", "")
		tail := helpers.GetIntArg(args, "tail", 50)
		follow := helpers.GetBoolArg(args, "follow", false)
		cmdArgs = append(cmdArgs, "logs", "--tail", fmt.Sprintf("%d", tail))
		if follow {
			cmdArgs = append(cmdArgs, "-f")
		}
		if services != "" {
			cmdArgs = append(cmdArgs, strings.Fields(services)...)
		}
	case "ps":
		all := helpers.GetBoolArg(args, "all", false)
		cmdArgs = append(cmdArgs, "ps")
		if all {
			cmdArgs = append(cmdArgs, "-a")
		}
	case "pull":
		services := helpers.GetStringArg(args, "services", "")
		cmdArgs = append(cmdArgs, "pull")
		if services != "" {
			cmdArgs = append(cmdArgs, strings.Fields(services)...)
		}
	case "config":
		cmdArgs = append(cmdArgs, "config")
	case "validate":
		cmdArgs = append(cmdArgs, "config", "--quiet")
	default:
		return fmt.Sprintf("Unknown compose action: %s (use up/down/restart/logs/ps/pull/config/validate)", action), false
	}

	// Execute
	cmd := exec.Command("docker", cmdArgs...)
	cmd.Dir = project
	output, err := cmd.CombinedOutput()
	outStr := strings.TrimSpace(string(output))
	if err != nil {
		return fmt.Sprintf("❌ Compose error: %v\n%s", err, helpers.TruncateStr(outStr, 2000)), false
	}
	if outStr == "" {
		outStr = "✅ Done (no output)"
	}
	return fmt.Sprintf("✅ `compose %s` completed\n<code>%s</code>", action, helpers.EscapeHTML(outStr)), true
}

func truncateOutput(s string, max int) string {
	if len(s) > max {
		return s[:max] + fmt.Sprintf("\n... [truncated %d bytes]", len(s)-max)
	}
	return s
}
