package scheduler

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
	"scorp-agent/config"
	"scorp-agent/internal/helpers"
	"scorp-agent/tools"
	"time"
)

// ──────────────────────────────────────────────
// Scheduler Extensions
// ──────────────────────────────────────────────

// gateScheduledShell applies the outer layers of the interactive gate stack
// (deny rules + sensitive-path sandbox) to scheduled shell tasks. A cron task
// used to be a silent bypass for the entire gate stack — the same commands
// that require confirmation (or are outright denied) interactively would run
// unattended on a timer. Scheduled write-ops are fail-closed: denied tasks
// are logged and surfaced via the error-notification path.
func gateScheduledShell(task ScheduledTask) (string, bool) {
	args := map[string]interface{}{"command": task.Prompt}
	if blocked, reason := config.CheckDenyRules("shell", args); blocked {
		return "🚫 " + reason, false
	}
	if restricted, reason := config.IsPathRestricted(task.Prompt); restricted {
		return "🛡️ " + reason, false
	}
	return "", true
}

// runShellTaskConfig runs a shell task with per-job timeout config.
func runShellTaskConfig(task ScheduledTask) (string, string) {
	if denyReason, ok := gateScheduledShell(task); !ok {
		log.Printf("[scheduler] Task %s blocked by gate: %s", task.ID, denyReason)
		return denyReason, "error"
	}

	timeoutSec := task.Timeout
	if timeoutSec <= 0 {
		timeoutSec = 30
	}
	timeoutSec = min(timeoutSec, 600)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSec)*time.Second)
	defer cancel()

	// Same sandbox wrap as interactive shell execution (no-op when bwrap is
	// unavailable).
	argv, wrapped := tools.SandboxWrap(task.Prompt)
	var cmd *exec.Cmd
	if wrapped {
		cmd = exec.Command(argv[0], argv[1:]...)
	} else {
		cmd = exec.Command("bash", "-c", task.Prompt)
	}

	// Own process group so the timeout can kill backgrounded grandchildren:
	// CommandContext only signals the direct child, and a `sleep` grandchild
	// holding the stdout pipe kept CombinedOutput blocked long past the
	// deadline (mirror of tools/exec.go).
	tools.SetProcessGroup(cmd)

	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	if err := cmd.Start(); err != nil {
		return fmt.Sprintf("failed to start: %v", err), "error"
	}

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	var err error
	timedOut := false
	select {
	case <-ctx.Done():
		timedOut = true
		if cmd.Process != nil {
			// Negative PID kills the whole group (bash + backgrounded children).
			tools.KillProcessGroup(cmd)
		}
		<-done
		err = fmt.Errorf("timeout after %ds", timeoutSec)
	case err = <-done:
	}
	output := buf.String()

	// Receipts: scheduled execution must leave the same audit trail as
	// interactive runs.
	tools.RecordToolReceipt("shell", map[string]interface{}{
		"command":        task.Prompt,
		"scheduled_task": task.ID,
	}, output, err == nil)

	if err != nil {
		if timedOut {
			err = fmt.Errorf("timeout after %ds (process group killed)", timeoutSec)
		}
		if task.NotifyOnError || task.NotifyOnSuccess {
			msg := fmt.Sprintf("⏰ <b>Scheduled Task (Error)</b>\n\n"+
				"<b>Task:</b> %s\n<b>Error:</b> %s\n\n<pre>%s</pre>",
				helpers.EscapeHTML(helpers.TruncateStr(task.Name+" ("+task.ID+")", 100)),
				helpers.EscapeHTML(err.Error()),
				helpers.EscapeHTML(helpers.TruncateStr(output, 2000)))
			notifyTaskResult(task, msg)
		}
		return output, "error"
	}

	if task.NotifyOnSuccess && len(output) > 0 {
		msg := fmt.Sprintf("⏰ <b>Scheduled Task ✓</b>\n\n<b>%s</b>\n\n<pre>%s</pre>",
			helpers.EscapeHTML(task.Name),
			helpers.EscapeHTML(helpers.TruncateStr(output, 3000)))
		notifyTaskResult(task, msg)
	}

	return output, "ok"
}

// runScriptTask executes a script file directly (script-only mode).
func runScriptTask(task ScheduledTask) (string, string) {
	timeoutSec := task.Timeout
	if timeoutSec <= 0 {
		timeoutSec = 120
	}
	timeoutSec = min(timeoutSec, 600)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSec)*time.Second)
	defer cancel()

	var cmd *exec.Cmd

	if isLikelyScriptPath(task.Prompt) {
		cmd = exec.CommandContext(ctx, "bash", task.Prompt)
	} else {
		cmd = exec.CommandContext(ctx, "bash", "-c", task.Prompt)
	}

	out, err := cmd.CombinedOutput()
	output := string(out)

	if err != nil {
		if task.NotifyOnError || task.NotifyOnSuccess {
			msg := fmt.Sprintf("📜 <b>Script Task (Error): %s</b>\n\n<pre>%s</pre>",
				helpers.EscapeHTML(task.Name),
				helpers.EscapeHTML(helpers.TruncateStr(output, 2000)))
			notifyTaskResult(task, msg)
		}
		return output, "error"
	}

	if task.NotifyOnSuccess && len(output) > 0 {
		msg := fmt.Sprintf("📜 <b>Script Task ✓: %s</b>\n\n<pre>%s</pre>",
			helpers.EscapeHTML(task.Name),
			helpers.EscapeHTML(helpers.TruncateStr(output, 3000)))
		notifyTaskResult(task, msg)
	}

	return output, "ok"
}

// notifyTaskResult sends a notification respecting per-job ChatTarget override.
func notifyTaskResult(task ScheduledTask, msg string) {
	if task.ChatTarget != 0 {
		chunks := splitMessage(msg, 4000)
		for _, chunk := range chunks {
			payload := map[string]interface{}{
				"chat_id":                  task.ChatTarget,
				"text":                     chunk,
				"parse_mode":               "HTML",
				"disable_web_page_preview": true,
			}
			if TgPost != nil {
				TgPost("/sendMessage", payload)
			}
		}
	} else {
		if SendMessage != nil {
			SendMessage(msg, nil)
		}
	}
}

func splitMessage(text string, maxLen int) []string {
	if len([]rune(text)) <= maxLen {
		return []string{text}
	}

	var chunks []string
	runes := []rune(text)
	for i := 0; i < len(runes); i += maxLen {
		end := i + maxLen
		if end > len(runes) {
			end = len(runes)
		}
		chunks = append(chunks, string(runes[i:end]))
	}
	return chunks
}

func isLikelyScriptPath(s string) bool {
	if len(s) == 0 || len(s) > 512 {
		return false
	}
	if s[0] != '/' && s[0] != '.' {
		return false
	}
	_ = exec.Command("test", "-f", s).Run()
	return true
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
