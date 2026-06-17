package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// ──────────────────────────────────────────────
// Phase 7 — Autonomous Agent Tools (user-facing)
// ──────────────────────────────────────────────

func executeAutonomous(args map[string]interface{}, chatID int64) (string, bool) {
	action, _ := args["action"].(string)
	if action == "" {
		action = "status"
	}

	switch strings.ToLower(action) {
	case "status":
		return autoStatus(), true

	case "enable":
		autoMu.Lock()
		autoConfig.Enabled = true
		autoConfig.KillSwitch = false
		autoMu.Unlock()
		os.Remove(autoKillFile)
		saveAutonomousConfig()
		return "✅ Autonomous agent **ENABLED**.\n\n" + autoStatus(), true

	case "disable":
		autoMu.Lock()
		autoConfig.Enabled = false
		autoMu.Unlock()
		saveAutonomousConfig()
		return "⏸️ Autonomous agent **DISABLED**.\n\n" + autoStatus(), true

	case "kill":
		setKillSwitch(true)
		return "🛑 **KILL SWITCH ACTIVATED**.\n\nAutonomous mode is now OFF and locked. Use `autonomous revive` to restart.", true

	case "revive":
		setKillSwitch(false)
		autoMu.Lock()
		autoConfig.Enabled = true
		autoMu.Unlock()
		saveAutonomousConfig()
		return "🔄 **REVIVED** — Kill switch removed, autonomous agent enabled.\n\n" + autoStatus(), true

	case "run":
		// Manual cycle trigger
		go func() {
			runAutonomousCycle()
			// Send result to user
			autoMu.Lock()
			lastEntry := AutonomousLogEntry{}
			if len(autoLog) > 0 {
				lastEntry = autoLog[len(autoLog)-1]
			}
			autoMu.Unlock()

			msg := fmt.Sprintf("🤖 <b>Manual autonomous cycle #%d complete</b>\n%s",
				autoCycleNum, lastEntry.Result)
			sendMessage(msg, nil)
		}()
		return "▶️ Triggering autonomous cycle now...", true

	case "config":
		return autoShowConfig(args), true

	case "log":
		return autoShowLog(args), true

	case "actions":
		return autoShowActions(), true

	default:
		return "Unknown action: " + action + "\n\nAvailable: status, enable, disable, kill, revive, run, config, log, actions", false
	}
}

func autoStatus() string {
	autoMu.Lock()
	defer autoMu.Unlock()

	enabled := "❌ Disabled"
	if autoConfig.Enabled && !autoConfig.KillSwitch {
		enabled = "✅ Enabled"
	}
	if autoConfig.KillSwitch {
		enabled = "🛑 KILL SWITCH ACTIVE"
	}

	lastCycle := "never"
	if !autoConfig.LastCycle.IsZero() {
		ago := time.Since(autoConfig.LastCycle).Round(time.Second)
		lastCycle = ago.String() + " ago"
	}

	return fmt.Sprintf(`🤖 **Autonomous Agent Status**

State: %s
Interval: %v
Approval: %s
Voice: %s
Max actions/cycle: %d

Cycles: %d
Actions taken: %d
Last cycle: %s

Commands:
• **autonomous enable** — Start autonomous mode
• **autonomous disable** — Pause
• **autonomous kill** — Emergency stop (locked)
• **autonomous revive** — Remove kill switch + enable
• **autonomous run** — Trigger manual cycle
• **autonomous config** interval=5m voice=important
• **autonomous log 10** — Show recent actions
• **autonomous actions** — Show action history`,
		enabled, autoConfig.Interval, autoConfig.ApprovalLevel,
		autoConfig.VoiceMode, autoConfig.MaxActions,
		autoConfig.TotalCycles, autoConfig.TotalActions, lastCycle)
}

func autoShowConfig(args map[string]interface{}) string {
	// Parse config changes
	if interval, ok := args["interval"].(string); ok && interval != "" {
		d, err := time.ParseDuration(interval)
		if err == nil && d >= 1*time.Minute {
			autoMu.Lock()
			autoConfig.Interval = d
			autoMu.Unlock()
			saveAutonomousConfig()
		}
	}
	if voice, ok := args["voice"].(string); ok && voice != "" {
		if voice == "off" || voice == "important" || voice == "always" {
			autoMu.Lock()
			autoConfig.VoiceMode = voice
			autoMu.Unlock()
			saveAutonomousConfig()
		}
	}
	if approval, ok := args["approval"].(string); ok && approval != "" {
		if approval == "low" || approval == "medium" || approval == "high" {
			autoMu.Lock()
			autoConfig.ApprovalLevel = approval
			autoMu.Unlock()
			saveAutonomousConfig()
		}
	}
	if maxStr, ok := args["max_actions"].(float64); ok {
		autoMu.Lock()
		autoConfig.MaxActions = int(maxStr)
		autoMu.Unlock()
		saveAutonomousConfig()
	}

	autoMu.Lock()
	defer autoMu.Unlock()
	return fmt.Sprintf(`⚙️ **Autonomous Config**

Enabled: %v
Interval: %v
Approval level: %s (low=auto-execute all, medium=approve high risk, high=approve all)
Voice mode: %s
Max actions/cycle: %d

To change: **autonomous config** interval=5m voice=important approval=medium`,
		autoConfig.Enabled, autoConfig.Interval, autoConfig.ApprovalLevel,
		autoConfig.VoiceMode, autoConfig.MaxActions)
}

func autoShowLog(args map[string]interface{}) string {
	count := 10
	if n, ok := args["count"].(float64); ok && int(n) > 0 {
		count = int(n)
	}
	if count > len(autoLog) {
		count = len(autoLog)
	}

	if count == 0 {
		return "📝 No autonomous actions logged yet."
	}

	start := len(autoLog) - count
	var b strings.Builder
	b.WriteString(fmt.Sprintf("📝 **Last %d autonomous actions:**\n\n", count))

	for _, e := range autoLog[start:] {
		status := "✅"
		if !e.Success {
			if !e.Approved {
				status = "🚫"
			} else {
				status = "❌"
			}
		}
		b.WriteString(fmt.Sprintf("%s [%s] **%s** (%s)\n", status,
			e.Timestamp.Format("01-02 15:04"), e.Tool, e.Risk))
		if e.Reason != "" {
			b.WriteString("   " + e.Reason + "\n")
		}
		b.WriteString("   → " + truncateStr(e.Result, 80) + "\n\n")
	}

	return b.String()
}

func autoShowActions() string {
	if len(autoLog) == 0 {
		return "📝 No autonomous actions yet."
	}

	// Group by tool
	type statEntry struct {
		count   int
		success int
		blocked int
	}
	stats := make(map[string]*statEntry)
	for _, e := range autoLog {
		s, ok := stats[e.Tool]
		if !ok {
			s = &statEntry{}
			stats[e.Tool] = s
		}
		s.count++
		if e.Success {
			s.success++
		}
		if !e.Approved {
			s.blocked++
		}
	}

	var b strings.Builder
	b.WriteString("📊 **Autonomous Action Stats:**\n\n")
	b.WriteString(fmt.Sprintf("Total cycles: %d, Total actions: %d\n\n", autoConfig.TotalCycles, autoConfig.TotalActions))

	for tool, s := range stats {
		b.WriteString(fmt.Sprintf("• **%s**: %d total, %d success", tool, s.count, s.success))
		if s.blocked > 0 {
			b.WriteString(fmt.Sprintf(", %d blocked", s.blocked))
		}
		b.WriteString("\n")
	}

	return b.String()
}
