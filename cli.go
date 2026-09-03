package main

import (
	"bufio"
	"fmt"
	"html"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"

	"scorp-agent/agent"
	"scorp-agent/bootstrap"
	"scorp-agent/config"
	"scorp-agent/models"
	"scorp-agent/registry"
	"scorp-agent/skills"
	"scorp-agent/tools"
)

const cliChatID int64 = 0

var (
	cliMu             sync.Mutex
	lastPrintedCount  int
	isTurnActive      bool
	currentThinkingID int64
)

// startCLI runs the CLI mode (one-shot or interactive REPL).
func startCLI(initialPrompts ...string) {
	// ── Setup clean CLI logging (write to ~/.scorp/scorp.log unless SCORP_DEBUG=1 or --debug) ──
	setupCLILogging()

	// ── Config ──
	config.InitConfigManager()

	// ── Models ──
	models.LoadModelConfig()
	models.InitModelUsage()
	models.LoadCostConfig()
	models.LoadCostTracker()

	// ── Skills ──
	skills.Load()

	// ── Bootstrap tool registry ──
	bootstrap.RegisterAutonomous()

	// ── Wire CLI callbacks ──
	wireCLICallbacks()

	// ── Init agent state ──
	agent.LoadAutonomousConfig()
	agent.LoadAutoLog()

	chatIDStr := fmt.Sprintf("%d", cliChatID)

	// ── Check for one-shot argument prompt ──
	initialPrompt := strings.TrimSpace(strings.Join(initialPrompts, " "))
	if initialPrompt != "" {
		executeOneShot(chatIDStr, initialPrompt)
		return
	}

	// ── Interactive REPL / Piped Line-by-Line ──
	if isTerminal(os.Stdin) {
		printBanner()
	}

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024) // 1MB buffer

	for {
		if isTerminal(os.Stdin) {
			fmt.Print("\nscorp> ")
		}
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		// Check if awaiting confirmation for dangerous command
		if agent.HasPendingConfirmation(chatIDStr) {
			lower := strings.ToLower(input)
			if lower == "y" || lower == "yes" || lower == "/confirm_yes" || lower == "allow" {
				fmt.Println("✓ Command approved. Resuming execution...")
				agent.HandleConfirmation(cliChatID, true)
				continue
			} else if lower == "n" || lower == "no" || lower == "/confirm_no" || lower == "deny" {
				fmt.Println("❌ Command rejected.")
				agent.HandleConfirmation(cliChatID, false)
				continue
			} else {
				fmt.Println("⚠️ A dangerous command is pending confirmation. Type 'y'/'yes' to allow, or 'n'/'no' to reject.")
				continue
			}
		}

		// Built-in slash commands
		if strings.HasPrefix(input, "/") {
			parts := strings.Fields(input)
			cmd := parts[0]

			switch cmd {
			case "/exit", "/quit", "/q":
				fmt.Println("Goodbye! 👋")
				return
			case "/help":
				printCLIHelp()
				continue
			case "/clear":
				agent.ClearChatSession(chatIDStr)
				fmt.Println("✓ Conversation history cleared.")
				continue
			case "/stop":
				agent.ExitAgentMode(chatIDStr)
				fmt.Println("✓ Agent mode stopped.")
				continue
			case "/models":
				printModelList()
				continue
			case "/model":
				if len(parts) < 2 {
					printCurrentModel()
				} else {
					targetModel := parts[1]
					if err := models.SwitchActiveModel(targetModel); err != nil {
						fmt.Printf("❌ Failed to switch model: %v\n", err)
					} else {
						fmt.Printf("✓ Active model switched to: %s\n", targetModel)
					}
				}
				continue
			case "/tools":
				printToolList()
				continue
			case "/status":
				printStatus(chatIDStr)
				continue
			case "/cost", "/usage":
				printCostUsage()
				continue
			default:
				fmt.Printf("Unknown command '%s'. Type /help for available commands.\n", cmd)
				continue
			}
		}

		// Execute prompt in agent loop
		executeTurn(chatIDStr, input)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Input error: %v\n", err)
	}
}

// executeOneShot executes a single prompt or slash command and exits cleanly
func executeOneShot(chatIDStr, prompt string) {
	if strings.HasPrefix(prompt, "/") {
		parts := strings.Fields(prompt)
		cmd := parts[0]
		switch cmd {
		case "/exit", "/quit", "/q":
			return
		case "/help":
			printCLIHelp()
			return
		case "/clear":
			agent.ClearChatSession(chatIDStr)
			fmt.Println("✓ Conversation history cleared.")
			return
		case "/models":
			printModelList()
			return
		case "/model":
			if len(parts) < 2 {
				printCurrentModel()
			} else {
				targetModel := parts[1]
				if err := models.SwitchActiveModel(targetModel); err != nil {
					fmt.Printf("❌ Failed to switch model: %v\n", err)
				} else {
					fmt.Printf("✓ Active model switched to: %s\n", targetModel)
				}
			}
			return
		case "/tools":
			printToolList()
			return
		case "/status":
			printStatus(chatIDStr)
			return
		case "/cost", "/usage":
			printCostUsage()
			return
		}
	}

	executeTurn(chatIDStr, prompt)
}

// setupCLILogging directs debug logs to ~/.scorp/scorp.log unless debug is enabled
func setupCLILogging() {
	if os.Getenv("SCORP_DEBUG") != "" || hasDebugFlag() {
		return
	}
	logDir := config.ScorpDir()
	_ = os.MkdirAll(logDir, 0755)
	logFile, err := os.OpenFile(config.ScorpPath("scorp.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err == nil {
		log.SetOutput(logFile)
	} else {
		log.SetOutput(io.Discard)
	}
}

func hasDebugFlag() bool {
	for _, arg := range os.Args {
		if arg == "--debug" || arg == "-d" {
			return true
		}
	}
	return false
}

// executeTurn sets up state and triggers the multi-turn agent loop
func executeTurn(chatIDStr, prompt string) {
	cliMu.Lock()
	lastPrintedCount = 0
	isTurnActive = true
	currentThinkingID = 0
	cliMu.Unlock()

	agent.EnterAgentMode(chatIDStr)
	agent.RunAgentLoop(cliChatID, prompt, 0)

	cliMu.Lock()
	isTurnActive = false
	cliMu.Unlock()
}

// printBanner prints a clean startup banner
func printBanner() {
	primaryModel := "unknown"
	provider := "unknown"
	if m := models.RouteModel("agent"); m != nil {
		primaryModel = m.Model
		provider = m.Provider
	}

	fmt.Println("━━━ scorp — Autonomous AI Agent CLI ━━━")
	fmt.Printf("🤖 Active Model: %s (%s)\n", primaryModel, provider)
	fmt.Println("💡 Type your message and press Enter. Type /help for commands, /exit to quit.")
}

// printCLIHelp displays commands and usage
func printCLIHelp() {
	fmt.Print(`
Commands:
  /help          — Show this help message
  /models        — List all available models and their status
  /model <name>  — Switch active agent model (e.g. /model deepseek/deepseek-v4-flash)
  /tools         — List all registered agent tools
  /status        — Show agent session status and working directory
  /cost          — Show today's token usage and estimated cost
  /clear         — Clear conversation history
  /stop          — Reset/stop agent mode
  /exit          — Quit scorp

Usage Examples:
  scorp> Jelaskan isi file main.go dan cari fungsi yang menangani CLI
  scorp> Buat file hello.py yang mencetak angka genap 1-20 lalu jalankan dengan python
  scorp> Cek penggunaan disk dan memori sistem saat ini
`)
}

// printModelList displays configured models
func printModelList() {
	models.ModelCfgMu.RLock()
	cfg := models.ModelCfg
	models.ModelCfgMu.RUnlock()

	if cfg == nil || len(cfg.Models) == 0 {
		fmt.Println("No models configured.")
		return
	}

	active := cfg.AgentModel
	fmt.Println("\nAvailable Models:")
	fmt.Printf("%-35s %-15s %-10s %-15s\n", "MODEL", "PROVIDER", "ACTIVE", "KEY STATUS")
	fmt.Println(strings.Repeat("─", 80))

	var names []string
	for name := range cfg.Models {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		m := cfg.Models[name]
		isActive := ""
		if name == active {
			isActive = "★ ACTIVE"
		}
		keyStatus := models.KeySourceLabel(&m)
		fmt.Printf("%-35s %-15s %-10s %-15s\n", name, m.Provider, isActive, keyStatus)
	}
	fmt.Println("\nSwitch model using: /model <model-name>")
}

// printCurrentModel shows active model details
func printCurrentModel() {
	if m := models.RouteModel("agent"); m != nil {
		fmt.Printf("Active Model: %s\nProvider: %s | API: %s | Key: %s\n",
			m.Model, m.Provider, models.ResolveAPIFormat(m), models.KeySourceLabel(m))
	} else {
		fmt.Println("No active model resolved.")
	}
}

// printToolList lists all registered tools
func printToolList() {
	allTools := registry.GetAllTools()
	if len(allTools) == 0 {
		fmt.Println("No tools registered.")
		return
	}

	// Group by category
	categories := make(map[string][]registry.ToolDef)
	for _, t := range allTools {
		categories[t.Category] = append(categories[t.Category], t)
	}

	var catNames []string
	for c := range categories {
		catNames = append(catNames, c)
	}
	sort.Strings(catNames)

	fmt.Printf("\nRegistered Tools (%d total):\n", len(allTools))
	for _, c := range catNames {
		fmt.Printf("\n[%s]\n", strings.ToUpper(c))
		toolsList := categories[c]
		for _, t := range toolsList {
			fmt.Printf("  • %-18s — %s\n", t.Name, t.Description)
		}
	}
}

// printStatus shows session status
func printStatus(chatIDStr string) {
	cwd, _ := os.Getwd()
	activeModel := "unknown"
	if m := models.RouteModel("agent"); m != nil {
		activeModel = fmt.Sprintf("%s (%s)", m.Model, m.Provider)
	}

	fmt.Printf("\nAgent Status:\n")
	fmt.Printf("  • Working Directory: %s\n", cwd)
	fmt.Printf("  • Active Model:      %s\n", activeModel)
	fmt.Printf("  • Session Chat ID:   %s\n", chatIDStr)
	fmt.Printf("  • Agent Mode:        %v\n", agent.IsAgentMode(chatIDStr))
	if agent.HasPendingConfirmation(chatIDStr) {
		cmd, tool, _ := agent.GetPendingConfirmationDetails(chatIDStr)
		fmt.Printf("  • Pending Confirm:   tool=%s cmd=%s\n", tool, cmd)
	}
}

// printCostUsage displays token usage and cost stats
func printCostUsage() {
	usage := models.FormatUsageStats()
	cost := models.FormatDailyCostSummary()
	fmt.Printf("\n%s\n\n%s\n", formatTerminalText(usage), formatTerminalText(cost))
}

// wireCLICallbacks replaces Telegram callbacks with clean terminal output
func wireCLICallbacks() {
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
				contentLines = append(contentLines, formatTerminalText(trimmed))
			}

			// Print only newly added lines
			if len(contentLines) > lastPrintedCount {
				for i := lastPrintedCount; i < len(contentLines); i++ {
					line := contentLines[i]
					if strings.HasPrefix(line, "🔧") || strings.HasPrefix(line, "🤖") || strings.HasPrefix(line, "⚡") {
						fmt.Printf("  ⚡ %s\n", strings.TrimPrefix(strings.TrimPrefix(line, "🔧 "), "🤖 "))
					} else if strings.HasPrefix(line, "→") || strings.HasPrefix(line, "↳") {
						fmt.Printf("     ↳ %s\n", strings.TrimSpace(strings.TrimPrefix(line, "→")))
					} else {
						fmt.Printf("  • %s\n", line)
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
}

// formatFinalResponse cleans up the markdown/HTML header and formats response
func formatFinalResponse(text string) string {
	clean := stripHTML(text)
	// Remove Scorp header if present to avoid redundant prefixes
	reHeader := regexp.MustCompile(`(?i)^(?:🤖\s*)?(?:Scorp:\s*)?`)
	clean = reHeader.ReplaceAllString(strings.TrimSpace(clean), "")
	return strings.TrimSpace(clean)
}

// formatTerminalText converts HTML formatting to clean readable terminal text
func formatTerminalText(s string) string {
	return strings.TrimSpace(stripHTML(s))
}

// stripHTML removes HTML tags, decodes entities, and cleans whitespace
func stripHTML(s string) string {
	// Replace <pre><code> and </code></pre> with code block backticks
	s = regexp.MustCompile(`(?i)<pre><code>([\s\S]*?)</code></pre>`).ReplaceAllString(s, "\n```\n$1\n```\n")
	s = regexp.MustCompile(`(?i)<pre>([\s\S]*?)</pre>`).ReplaceAllString(s, "\n```\n$1\n```\n")
	s = regexp.MustCompile(`(?i)<code>([\s\S]*?)</code>`).ReplaceAllString(s, "`$1`")

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
	return html.UnescapeString(sb.String())
}

// isTerminal checks if a file descriptor is an interactive terminal
func isTerminal(f *os.File) bool {
	stat, err := f.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) != 0
}
