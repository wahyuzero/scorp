package scheduler

import (
	"scorp-agent/internal/helpers"
	"context"
	"fmt"
	"os/exec"
	"time"
)

// ──────────────────────────────────────────────
// Scheduler Extensions
// ──────────────────────────────────────────────

// runShellTaskConfig runs a shell task with per-job timeout config.
func runShellTaskConfig(task ScheduledTask) (string, string) {
	timeoutSec := task.Timeout
	if timeoutSec <= 0 {
		timeoutSec = 30
	}
	timeoutSec = min(timeoutSec, 600)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSec)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "bash", "-c", task.Prompt)
	out, err := cmd.CombinedOutput()
	output := string(out)

	if err != nil {
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
