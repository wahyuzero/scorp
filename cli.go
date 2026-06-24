package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"scorp-agent/agent"
	"scorp-agent/bootstrap"
	"scorp-agent/config"
	"scorp-agent/models"
	"scorp-agent/skills"
	"scorp-agent/tools"
)

const cliChatID int64 = 0

// startCLI runs the interactive REPL mode (no Telegram needed).
func startCLI() {
	fmt.Println("━━━ scorp — CLI mode ━━━")
	fmt.Println("Type your message, or /help for commands. /exit to quit.")
	fmt.Println()

	// ── Config ──
	config.InitConfigManager()

	// ── Models ──
	models.LoadModelConfig()
	models.InitModelUsage()

	// ── Skills ──
	skills.Load()

	// ── Bootstrap tool registry ──
	bootstrap.RegisterAutonomous()

	// ── Wire CLI callbacks (replace Telegram with terminal) ──
	wireCLICallbacks()

	// ── Init agent state ──
	agent.LoadAutonomousConfig()
	agent.LoadAutoLog()

	// ── REPL loop ──
	chatIDStr := fmt.Sprintf("%d", cliChatID)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024) // 1MB buffer

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		// ── Built-in commands ──
		switch input {
		case "/exit", "/quit", "/q":
			fmt.Println("Bye!")
			return
		case "/help":
			printCLIHelp()
			continue
		case "/clear":
			agent.ClearChatSession(chatIDStr)
			fmt.Println("✓ History cleared.")
			continue
		case "/stop":
			agent.ExitAgentMode(chatIDStr)
			fmt.Println("✓ Agent mode stopped.")
			continue
		case "/models":
			fmt.Println("Use /models in Telegram, or edit ~/.scorp/models.json")
			continue
		}

		// ── Agent loop ──
		agent.EnterAgentMode(chatIDStr)
		agent.RunAgentLoop(cliChatID, input, 0)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Input error: %v\n", err)
	}
}

// wireCLICallbacks replaces all Telegram callbacks with terminal output.
func wireCLICallbacks() {
	var mu sync.Mutex

	// SendMessage → print to terminal
	tools.SendMessage = func(text string, keyboard map[string]interface{}) bool {
		mu.Lock()
		defer mu.Unlock()
		fmt.Println(stripHTML(text))
		fmt.Println()
		return true
	}

	// SendMessageGetID → print, return fake msgID
	tools.SendMessageGetID = func(text string, chatID int64) int64 {
		mu.Lock()
		defer mu.Unlock()
		fmt.Println(stripHTML(text))
		return 1
	}

	// EditMessageByID → overwrite (just print)
	tools.EditMessageByID = func(chatID int64, messageID int64, text string, keyboard map[string]interface{}) bool {
		mu.Lock()
		defer mu.Unlock()
		fmt.Printf("\r\033[2K%s\n", stripHTML(text))
		return true
	}

	// SendChatAction → noop
	tools.SendChatAction = func(chatID int64, action string) {}

	// TgPost → noop (return success)
	tools.TgPost = func(method string, payload map[string]interface{}) (tools.TgResponse, error) {
		return tools.TgResponse{OK: true}, nil
	}

	// Agent callbacks
	tools.StorePendingConfirmation = func(chatID, toolName, command string, _ []tools.AgentMessage) {
		agent.StorePendingConfirmation(chatID, toolName, command, nil)
	}
	tools.IsDangerousCommand = func(cmd string) bool {
		return agent.IsDangerousCommand(cmd)
	}

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

	// Model env callback
	models.UpdateEnvFile = func(key, value string) {
		// noop in CLI mode — user edits .env manually
	}
}

// stripHTML removes simple HTML tags for terminal display.
func stripHTML(s string) string {
	var sb strings.Builder
	inTag := false
	for _, r := range s {
		switch {
		case r == '<':
			inTag = true
		case r == '>':
			inTag = false
		case !inTag:
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

func printCLIHelp() {
	fmt.Print(`
Commands:
  /clear   — Clear conversation history
  /stop    — Exit agent mode (reset)
  /exit    — Quit scorp
  /help    — This message

Or just type your message and press Enter.
`)
}
