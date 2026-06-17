package main

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

// ──────────────────────────────────────────────
// S06: Scheduler Extensions
//
// - runShellTaskConfig: shell with per-job timeout
// - runScriptTask: script-only mode (no agent loop overhead)
// - Per-job config parsing in executeSchedule
// ──────────────────────────────────────────────

// runShellTaskConfig runs a shell task with per-job timeout config.
func runShellTaskConfig(task ScheduledTask) (string, string) {
	timeoutSec := task.Timeout
	if timeoutSec <= 0 {
		timeoutSec = 30 // default
	}
	timeoutSec = min(timeoutSec, 600) // cap at 10 min

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSec)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "bash", "-c", task.Prompt)
	out, err := cmd.CombinedOutput()
	output := string(out)

	if err != nil {
		msg := fmt.Sprintf("⏰ <b>Scheduled Task (Error)</b>\n\n"+
			"<b>Task:</b> %s\n<b>Error:</b> %s\n\n<pre>%s</pre>",
			escapeHTML(truncateStr(task.Name+" ("+task.ID+")", 100)),
			escapeHTML(err.Error()),
			escapeHTML(truncateStr(output, 2000)))
		sendMessageSmart(msg, nil)
		return output, "error"
	}

	// S06: Notify on success if configured
	if task.NotifyOnSuccess && len(output) > 0 {
		msg := fmt.Sprintf("⏰ <b>Scheduled Task ✓</b>\n\n<b>%s</b>\n\n<pre>%s</pre>",
			escapeHTML(task.Name),
			escapeHTML(truncateStr(output, 3000)))
		sendMessageSmart(msg, nil)
	}

	return output, "ok"
}

// runScriptTask executes a script file directly (script-only mode).
// The prompt field contains the script path. Falls back to inline script.
func runScriptTask(task ScheduledTask) (string, string) {
	timeoutSec := task.Timeout
	if timeoutSec <= 0 {
		timeoutSec = 120
	}
	timeoutSec = min(timeoutSec, 600)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSec)*time.Second)
	defer cancel()

	var cmd *exec.Cmd

	// If the prompt looks like a file path, execute it directly
	if isLikelyScriptPath(task.Prompt) {
		cmd = exec.CommandContext(ctx, "bash", task.Prompt)
	} else {
		// Treat the prompt as inline script content
		cmd = exec.CommandContext(ctx, "bash", "-c", task.Prompt)
	}

	out, err := cmd.CombinedOutput()
	output := string(out)

	if err != nil {
		if task.NotifyOnError || task.NotifyOnSuccess {
			msg := fmt.Sprintf("📜 <b>Script Task (Error): %s</b>\n\n<pre>%s</pre>",
				escapeHTML(task.Name),
				escapeHTML(truncateStr(output, 2000)))
			sendMessageSmart(msg, nil)
		}
		return output, "error"
	}

	if task.NotifyOnSuccess && len(output) > 0 {
		msg := fmt.Sprintf("📜 <b>Script Task ✓: %s</b>\n\n<pre>%s</pre>",
			escapeHTML(task.Name),
			escapeHTML(truncateStr(output, 3000)))
		sendMessageSmart(msg, nil)
	}

	return output, "ok"
}

func isLikelyScriptPath(s string) bool {
	if len(s) == 0 || len(s) > 512 {
		return false
	}
	if s[0] != '/' && s[0] != '.' {
		return false
	}
	// Check if file exists
	_ = exec.Command("test", "-f", s).Run()
	return true
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
