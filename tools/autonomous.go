package tools

import (
	"scorp-agent/internal/helpers"
	"fmt"
	"os"
	"strings"
	"time"
)

// ──────────────────────────────────────────────
// Phase 7 — Autonomous Agent Tools (user-facing)
// ──────────────────────────────────────────────

// ExecuteAutonomous handles autonomous agent control commands
func ExecuteAutonomous(args map[string]interface{}, chatID int64) (string, bool) {
	action, _ := args["action"].(string)
	if action == "" {
		action = "status"
	}

	switch strings.ToLower(action) {
	case "status":
		return autoStatus(), true

	case "enable":
		AutoMu.Lock()
		AutoConfig.Enabled = true
		AutoConfig.KillSwitch = false
		AutoMu.Unlock()
		os.Remove(AutoKillFile)
		SaveAutonomousConfig()
		return "✅ Autonomous agent **ENABLED**.\n\n" + autoStatus(), true

	case "disable":
		AutoMu.Lock()
		AutoConfig.Enabled = false
		AutoMu.Unlock()
		SaveAutonomousConfig()
		return "⏸️ Autonomous agent **DISABLED**.\n\n" + autoStatus(), true

	case "kill":
		SetKillSwitch(true)
		return "🛑 **KILL SWITCH ACTIVATED**.\n\nAutonomous mode is now OFF and locked. Use `autonomous revive` to restart.", true

	case "revive":
		SetKillSwitch(false)
		AutoMu.Lock()
		AutoConfig.Enabled = true
		AutoMu.Unlock()
		SaveAutonomousConfig()
		return "🔄 **REVIVED** — Kill switch removed, autonomous agent enabled.\n\n" + autoStatus(), true

	case "run":
		// Manual cycle trigger
		go func() {
			RunAutonomousCycle()
			// Send result to user
			AutoMu.Lock()
			lastEntry := AutonomousLogEntry{}
			if len(*AutoLog) > 0 {
				lastEntry = (*AutoLog)[len(*AutoLog)-1]
			}
			AutoMu.Unlock()

			msg := fmt.Sprintf("🤖 <b>Manual autonomous cycle #%d complete</b>\n%s",
				AutoCycleNum, lastEntry.Result)
			SendMessage(msg, nil)
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
	AutoMu.Lock()
	defer AutoMu.Unlock()

	enabled := "❌ Disabled"
	if AutoConfig.Enabled && !AutoConfig.KillSwitch {
		enabled = "✅ Enabled"
	}
	if AutoConfig.KillSwitch {
		enabled = "🛑 KILL SWITCH ACTIVE"
	}

	lastCycle := "never"
	if !AutoConfig.LastCycle.IsZero() {
		ago := time.Since(AutoConfig.LastCycle).Round(time.Second)
		lastCycle = ago.String() + " ago"
	}

	return fmt.Sprintf(`🤖 **Autonomous Agent Status**

State: %s
Interval: %v
Approval: %s
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
• **autonomous config** interval=5m approval=medium
• **autonomous log 10** — Show recent actions
• **autonomous actions** — Show action history`,
		enabled, AutoConfig.Interval, AutoConfig.ApprovalLevel,
		AutoConfig.MaxActions,
		AutoConfig.TotalCycles, AutoConfig.TotalActions, lastCycle)
}

func autoShowConfig(args map[string]interface{}) string {
	// Parse config changes
	if interval, ok := args["interval"].(string); ok && interval != "" {
		d, err := time.ParseDuration(interval)
		if err == nil && d >= 1*time.Minute {
			AutoMu.Lock()
			AutoConfig.Interval = d
			AutoMu.Unlock()
			SaveAutonomousConfig()
		}
	}
	if approval, ok := args["approval"].(string); ok && approval != "" {
		if approval == "low" || approval == "medium" || approval == "high" {
			AutoMu.Lock()
			AutoConfig.ApprovalLevel = approval
			AutoMu.Unlock()
			SaveAutonomousConfig()
		}
	}
	if maxStr, ok := args["max_actions"].(float64); ok {
		AutoMu.Lock()
		AutoConfig.MaxActions = int(maxStr)
		AutoMu.Unlock()
		SaveAutonomousConfig()
	}

	AutoMu.Lock()
	defer AutoMu.Unlock()
	return fmt.Sprintf(`⚙️ **Autonomous Config**

Enabled: %v
Interval: %v
Approval level: %s (low=auto-execute all, medium=approve high risk, high=approve all)
Max actions/cycle: %d

To change: **autonomous config** interval=5m approval=medium`,
		AutoConfig.Enabled, AutoConfig.Interval, AutoConfig.ApprovalLevel,
		AutoConfig.MaxActions)
}

func autoShowLog(args map[string]interface{}) string {
	count := 10
	if n, ok := args["count"].(float64); ok && int(n) > 0 {
		count = int(n)
	}
	if count > len(*AutoLog) {
		count = len(*AutoLog)
	}

	if count == 0 {
		return "📝 No autonomous actions logged yet."
	}

	start := len(*AutoLog) - count
	var b strings.Builder
	b.WriteString(fmt.Sprintf("📝 **Last %d autonomous actions:**\n\n", count))

	for _, e := range (*AutoLog)[start:] {
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
		b.WriteString("   → " + helpers.TruncateStr(e.Result, 80) + "\n\n")
	}

	return b.String()
}

func autoShowActions() string {
	if len(*AutoLog) == 0 {
		return "📝 No autonomous actions yet."
	}

	// Group by tool
	type statEntry struct {
		count   int
		success int
		blocked int
	}
	stats := make(map[string]*statEntry)
	for _, e := range *AutoLog {
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
	b.WriteString(fmt.Sprintf("Total cycles: %d, Total actions: %d\n\n", AutoConfig.TotalCycles, AutoConfig.TotalActions))

	for tool, s := range stats {
		b.WriteString(fmt.Sprintf("• **%s**: %d total, %d success", tool, s.count, s.success))
		if s.blocked > 0 {
			b.WriteString(fmt.Sprintf(", %d blocked", s.blocked))
		}
		b.WriteString("\n")
	}

	return b.String()
}
