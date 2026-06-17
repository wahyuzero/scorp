package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// ──────────────────────────────────────────────
// Docker Compose Tool
// ──────────────────────────────────────────────

func executeCompose(args map[string]interface{}) (string, bool) {
	action := getStringArg(args, "action", "up")
	project := getStringArg(args, "project", ".")
	composeFile := getStringArg(args, "file", "")

	// Safety: never delete containers permanently
	if action == "down" {
		if !getBoolArg(args, "confirm", false) {
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
		detach := getBoolArg(args, "detach", true)
		rebuild := getBoolArg(args, "rebuild", false)
		cmdArgs = append(cmdArgs, "up")
		if detach {
			cmdArgs = append(cmdArgs, "-d")
		}
		if rebuild {
			cmdArgs = append(cmdArgs, "--build")
		}
	case "down":
		volumes := getBoolArg(args, "volumes", false)
		removeOrphans := getBoolArg(args, "remove_orphans", true)
		cmdArgs = append(cmdArgs, "down")
		if volumes {
			cmdArgs = append(cmdArgs, "-v")
		}
		if removeOrphans {
			cmdArgs = append(cmdArgs, "--remove-orphans")
		}
	case "restart":
		services := getStringArg(args, "services", "")
		timeout := getIntArg(args, "timeout", 10)
		cmdArgs = append(cmdArgs, "restart", "-t", fmt.Sprintf("%d", timeout))
		if services != "" {
			cmdArgs = append(cmdArgs, strings.Fields(services)...)
		}
	case "logs":
		services := getStringArg(args, "services", "")
		tail := getIntArg(args, "tail", 50)
		follow := getBoolArg(args, "follow", false)
		cmdArgs = append(cmdArgs, "logs", "--tail", fmt.Sprintf("%d", tail))
		if follow {
			cmdArgs = append(cmdArgs, "-f")
		}
		if services != "" {
			cmdArgs = append(cmdArgs, strings.Fields(services)...)
		}
	case "ps":
		all := getBoolArg(args, "all", false)
		cmdArgs = append(cmdArgs, "ps")
		if all {
			cmdArgs = append(cmdArgs, "-a")
		}
	case "pull":
		services := getStringArg(args, "services", "")
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
		return fmt.Sprintf("❌ Compose error: %v\n%s", err, truncateStr(outStr, 2000)), false
	}
	if outStr == "" {
		outStr = "✅ Done (no output)"
	}
	return fmt.Sprintf("✅ `compose %s` completed\n<code>%s</code>", action, escapeHTML(outStr)), true
}

func truncateOutput(s string, max int) string {
	if len(s) > max {
		return s[:max] + fmt.Sprintf("\n... [truncated %d bytes]", len(s)-max)
	}
	return s
}
