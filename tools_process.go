package main

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// ──────────────────────────────────────────────
// Process Manager Tool — ps, top, kill, systemctl
// ──────────────────────────────────────────────

func executeProcess(args map[string]interface{}, chatID int64) (string, bool) {
	action := getStringArg(args, "action", "")
	if action == "" {
		return "Error: 'action' is required (list, top, kill, service_status, service_restart, service_start, service_stop)", false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var cmd *exec.Cmd

	switch action {
	case "list":
		filter := getStringArg(args, "filter", "")
		if filter != "" {
			cmd = exec.CommandContext(ctx, "bash", "-c",
				fmt.Sprintf("ps aux | grep -i '%s' | grep -v grep | head -30", filter))
		} else {
			cmd = exec.CommandContext(ctx, "bash", "-c",
				"ps aux --sort=-%%mem | head -20")
		}

	case "top":
		sortBy := getStringArg(args, "sort_by", "mem")
		switch sortBy {
		case "cpu":
			cmd = exec.CommandContext(ctx, "bash", "-c",
				"ps aux --sort=-%%cpu | head -15")
		default:
			cmd = exec.CommandContext(ctx, "bash", "-c",
				"ps aux --sort=-%%mem | head -15")
		}

	case "kill":
		pid := getStringArg(args, "pid", "")
		if pid == "" {
			return "Error: 'pid' argument required for kill", false
		}
		// Kill requires confirmation
		chatIDStr := fmt.Sprintf("%d", chatID)
		signal := getStringArg(args, "signal", "TERM")
		storePendingConfirmation(chatIDStr, "process", fmt.Sprintf("kill -%s %s", signal, pid), nil)
		return fmt.Sprintf("⚠️ KILL DETECTED:\nkill -%s %s\n\nThis will send signal to the process. Please confirm.", signal, pid), false

	case "service_status":
		service := getStringArg(args, "service", "")
		if service == "" {
			return "Error: 'service' argument required", false
		}
		cmd = exec.CommandContext(ctx, "systemctl", "status", service, "--no-pager", "-l")

	case "service_restart":
		service := getStringArg(args, "service", "")
		if service == "" {
			return "Error: 'service' argument required", false
		}
		// Restart requires confirmation
		chatIDStr := fmt.Sprintf("%d", chatID)
		storePendingConfirmation(chatIDStr, "process", fmt.Sprintf("systemctl restart %s", service), nil)
		return fmt.Sprintf("⚠️ SERVICE RESTART:\nsystemctl restart %s\n\nThis will restart the service. Please confirm.", service), false

	case "service_start":
		service := getStringArg(args, "service", "")
		if service == "" {
			return "Error: 'service' argument required", false
		}
		cmd = exec.CommandContext(ctx, "systemctl", "start", service)

	case "service_stop":
		service := getStringArg(args, "service", "")
		if service == "" {
			return "Error: 'service' argument required", false
		}
		// Stop requires confirmation
		chatIDStr := fmt.Sprintf("%d", chatID)
		storePendingConfirmation(chatIDStr, "process", fmt.Sprintf("systemctl stop %s", service), nil)
		return fmt.Sprintf("⚠️ SERVICE STOP:\nsystemctl stop %s\n\nThis will stop the service. Please confirm.", service), false

	case "service_list":
		cmd = exec.CommandContext(ctx, "bash", "-c",
			"systemctl list-units --type=service --state=running --no-pager | head -30")

	case "ports":
		cmd = exec.CommandContext(ctx, "bash", "-c",
			"ss -tlnp | head -30")

	default:
		return fmt.Sprintf("Error: unknown action '%s'. Available: list, top, kill, service_status, service_restart, service_start, service_stop, service_list, ports", action), false
	}

	output, err := cmd.CombinedOutput()
	result := string(output)

	if ctx.Err() == context.DeadlineExceeded {
		return "Process command timed out", false
	}

	if err != nil {
		// systemctl status returns non-zero when service is inactive — still useful
		if strings.Contains(action, "status") && result != "" {
			return truncOutput(result, maxToolOutput), true
		}
		return fmt.Sprintf("Command failed: %v\n%s", err, truncOutput(result, maxToolOutput)), false
	}

	return truncOutput(result, maxToolOutput), true
}
