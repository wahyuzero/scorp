package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"scorp-agent/agent"
	"scorp-agent/models"
	"scorp-agent/tools"
)

// wireCLICallbacks replaces Telegram callbacks with clean terminal output
func wireCLICallbacks() {
	// Plan Mode approval — stdin prompt replaces the Telegram inline keyboard.
	// Runs inside the planning loop's presentation call, so EndPlanning must
	// be called before executing (the loop's own defer is idempotent).
	agent.PresentPlanHook = func(sessionID string, chatID int64, goal string, plan *agent.TaskPlan) {
		if goal != "" {
			fmt.Printf("\n🧭 PLAN READY — goal: %s\n", goal)
		} else {
			fmt.Println("\n🧭 PLAN READY (revised)")
		}
		fmt.Println(plan.Render())
		fmt.Print("\nExecute this plan with full tools? (y = execute / n = cancel): ")
		reader := bufio.NewReader(os.Stdin)
		ans, _ := reader.ReadString('\n')
		agent.EndPlanning()
		switch strings.ToLower(strings.TrimSpace(ans)) {
		case "y", "yes":
			fmt.Println("✅ Plan approved — executing…")
			agent.RunAgentSessionLoop(sessionID, chatID, "[✅ PLAN APPROVED] Execute every step of the task ledger now — mark each step done via task_plan as you verify it, then complete_task with the final verified report.", 0)
		default:
			agent.CancelPlan(sessionID)
			fmt.Println("❌ Plan cancelled — nothing will execute.")
		}
	}

	// SendMessage → print cleanly to terminal
	tools.SendMessage = func(text string, keyboard map[string]interface{}) bool {
		cliMu.Lock()
		defer cliMu.Unlock()

		clean := formatTerminalText(text)
		if clean != "" {
			fmt.Println(clean)
		}
		return true
	}

	// SendMessageGetID → initial thinking notice
	tools.SendMessageGetID = func(text string, chatID int64) int64 {
		cliMu.Lock()
		defer cliMu.Unlock()

		currentThinkingID = 1
		lastPrintedCount = 0

		// Show initial working indicator if it's the start message
		if strings.Contains(text, "processing") || strings.Contains(text, "🧠") {
			fmt.Println("⏳ Thinking...")
		} else {
			clean := formatTerminalText(text)
			if clean != "" {
				fmt.Println(clean)
			}
		}
		return currentThinkingID
	}

	tools.SendMessageGetIDWithKeyboard = func(text string, chatID int64, keyboard map[string]interface{}) int64 {
		cliMu.Lock()
		defer cliMu.Unlock()

		clean := formatTerminalText(text)
		if clean != "" {
			fmt.Println("\n" + clean)
		}
		return 1
	}

	// EditMessageByID → smart update handler
	tools.EditMessageByID = func(chatID int64, messageID int64, text string, keyboard map[string]interface{}) bool {
		cliMu.Lock()
		defer cliMu.Unlock()

		// Case 1: Final Scorp reply
		if strings.Contains(text, "Scorp:") || strings.Contains(text, "🤖") {
			clean := formatFinalResponse(text)
			if clean != "" {
				fmt.Printf("\n%s\n", clean)
			}
			return true
		}

		// Case 2: Confirmation required
		if strings.Contains(text, "Allow execution") || (strings.Contains(text, "Dangerous Command") && !strings.Contains(text, "APPROVED") && !strings.Contains(text, "REJECTED")) {
			clean := formatTerminalText(text)
			fmt.Printf("\n%s\nType 'y'/'yes' to allow or 'n'/'no' to reject.\n", clean)
			return true
		}

		// Case 3: Incremental thinking stream updates
		if strings.Contains(text, "🧠") {
			lines := strings.Split(text, "\n")
			// Filter out header and footer
			var contentLines []string
			for _, line := range lines {
				trimmed := strings.TrimSpace(line)
				if trimmed == "" || strings.HasPrefix(trimmed, "🧠") || strings.Contains(trimmed, "working...") || strings.Contains(trimmed, "processing...") {
					continue
				}
				// Filter out internal engine retry/continuation noise from user view
				if strings.Contains(trimmed, "Retry") || strings.Contains(trimmed, "continuation detected") || strings.Contains(trimmed, "CONTINUATION") {
					continue
				}
				contentLines = append(contentLines, formatTerminalText(trimmed))
			}

			// Print only newly added lines
			if len(contentLines) > lastPrintedCount {
				useColor := isTerminal(os.Stdout)
				for i := lastPrintedCount; i < len(contentLines); i++ {
					line := contentLines[i]
					if strings.HasPrefix(line, "🔧") || strings.HasPrefix(line, "🤖") || strings.HasPrefix(line, "⚡") {
						toolName := strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(line, "🔧 "), "🤖 "), "⚡ ")
						if useColor {
							fmt.Printf("  \033[1;33m⚡ %s\033[0m\n", toolName)
						} else {
							fmt.Printf("  ⚡ %s\n", toolName)
						}
					} else if strings.HasPrefix(line, "→") || strings.HasPrefix(line, "↳") {
						preview := strings.TrimSpace(strings.TrimPrefix(line, "→"))
						if useColor {
							fmt.Printf("     \033[2m↳ %s\033[0m\n", preview)
						} else {
							fmt.Printf("     ↳ %s\n", preview)
						}
					} else {
						if useColor {
							fmt.Printf("  \033[36m•\033[0m %s\n", line)
						} else {
							fmt.Printf("  • %s\n", line)
						}
					}
				}
				lastPrintedCount = len(contentLines)
			}
			return true
		}

		// Default fallback
		clean := formatTerminalText(text)
		if clean != "" {
			fmt.Println(clean)
		}
		return true
	}

	// SendChatAction → noop in CLI
	tools.SendChatAction = func(chatID int64, action string) {}

	// TgPost → noop
	tools.TgPost = func(method string, payload map[string]interface{}) (tools.TgResponse, error) {
		return tools.TgResponse{OK: true}, nil
	}

	// Agent callbacks
	tools.StorePendingConfirmation = func(chatID, toolName, command string, _ []tools.AgentMessage, promptMsgID ...int64) {
		agent.StorePendingConfirmation(chatID, toolName, command, nil, promptMsgID...)
	}
	tools.IsDangerousCommand = func(cmd string) bool {
		return agent.IsDangerousCommand(cmd)
	}
	tools.AppendDurableMemory = agent.AppendMemoryMD

	// Autonomous callbacks
	tools.AutoConfig = &agent.AutoConfig
	tools.AutoMu = &agent.AutoMu
	tools.AutoLog = &agent.AutoLog
	tools.AutoKillFile = agent.AutoKillFile
	tools.AutoCycleNum = &agent.AutoCycleNum
	tools.SaveAutonomousConfig = agent.SaveAutonomousConfig
	tools.SetKillSwitch = agent.SetKillSwitch
	tools.RunAutonomousCycle = agent.RunAutonomousCycle

	// ExecuteTool callback
	tools.ExecuteTool = func(tc models.ToolCall, chatID int64) (string, bool) {
		return agent.ExecuteTool(tc, chatID)
	}

	models.UpdateEnvFile = func(key, value string) {
		// CLI users manage .env manually
	}

	// Auto-titling callback for live CLI prompt update
	tools.OnSessionAutoTitled = func(oldID, newID, chatIDStr string) {
		cliMu.Lock()
		defer cliMu.Unlock()
		if currentSessionID == oldID || currentSessionID == "default" {
			currentSessionID = newID
			fmt.Printf("\n\033[1;34mℹ️ Session auto-titled:\033[0m \033[1;32m%s\033[0m\n", newID)
		}
	}
}
