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
	"time"

	"scorp-agent/agent"
	"scorp-agent/bootstrap"
	"scorp-agent/config"
	"scorp-agent/models"
	"scorp-agent/registry"
	"scorp-agent/skills"
	"scorp-agent/sop"
	"scorp-agent/tools"
)

const cliChatID int64 = 0

var (
	cliMu             sync.Mutex
	lastPrintedCount  int
	isTurnActive      bool
	currentThinkingID int64
	currentSessionID  = "default"
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

	// ── Autonomy & SOPs initialization ──
	if envMode := os.Getenv("SCORP_AUTONOMY"); envMode != "" {
		config.SetAutonomyLevel(envMode)
	}
	sop.InitDefaultSOPs()

	// ── Session selection (--session <name> or SCORP_SESSION env) ──
	for i := 1; i < len(os.Args); i++ {
		if (os.Args[i] == "--session" || os.Args[i] == "-s") && i+1 < len(os.Args) {
			currentSessionID = os.Args[i+1]
		} else if strings.HasPrefix(os.Args[i], "--session=") {
			currentSessionID = strings.TrimPrefix(os.Args[i], "--session=")
		}
	}
	if envSess := os.Getenv("SCORP_SESSION"); envSess != "" {
		currentSessionID = envSess
	}

	chatIDStr := currentSessionID

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
			if currentSessionID == "default" {
				fmt.Print("\n\033[1;36mscorp\033[0m \033[1;32m❯\033[0m ")
			} else {
				fmt.Printf("\n\033[1;36mscorp\033[0m \033[2m[%s]\033[0m \033[1;32m❯\033[0m ", currentSessionID)
			}
		}
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		// ── Multiline input mode (""" or ``` block, or trailing \) ──
		if strings.HasPrefix(input, `"""`) || strings.HasPrefix(input, "```") {
			delim := `"""`
			if strings.HasPrefix(input, "```") {
				delim = "```"
			}
			var lines []string
			firstLine := strings.TrimPrefix(input, delim)
			if firstLine != "" && !strings.HasSuffix(firstLine, delim) {
				lines = append(lines, firstLine)
			}
			if !strings.HasSuffix(input[len(delim):], delim) {
				for {
					if isTerminal(os.Stdin) {
						fmt.Print("\033[2m... ❯\033[0m ")
					}
					if !scanner.Scan() {
						break
					}
					l := scanner.Text()
					if strings.Contains(l, delim) {
						lines = append(lines, strings.TrimSuffix(l, delim))
						break
					}
					lines = append(lines, l)
				}
				input = strings.TrimSpace(strings.Join(lines, "\n"))
			} else {
				input = strings.TrimSuffix(firstLine, delim)
			}
		} else if strings.HasSuffix(input, "\\") {
			var lines []string
			lines = append(lines, strings.TrimSuffix(input, "\\"))
			for {
				if isTerminal(os.Stdin) {
					fmt.Print("\033[2m... ❯\033[0m ")
				}
				if !scanner.Scan() {
					break
				}
				l := scanner.Text()
				if strings.HasSuffix(l, "\\") {
					lines = append(lines, strings.TrimSuffix(l, "\\"))
				} else {
					lines = append(lines, l)
					break
				}
			}
			input = strings.TrimSpace(strings.Join(lines, "\n"))
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
			case "/mode":
				if len(parts) < 2 {
					fmt.Printf("🛡️ Current Autonomy Mode: %s\nUsage: /mode [readonly|supervised|yolo]\n", config.GetAutonomyLevel())
				} else {
					targetMode := parts[1]
					config.SetAutonomyLevel(targetMode)
					fmt.Printf("✓ Autonomy mode switched to: %s\n", config.GetAutonomyLevel())
				}
				continue
			case "/sop":
				handleCLISOP(parts[1:], chatIDStr)
				continue
			case "/receipts":
				printReceipts()
				continue
			case "/session", "/sessions":
				handleCLISession(parts[1:])
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
		case "/mode":
			if len(parts) < 2 {
				fmt.Printf("🛡️ Current Autonomy Mode: %s\nUsage: /mode [readonly|supervised|yolo]\n", config.GetAutonomyLevel())
			} else {
				targetMode := parts[1]
				config.SetAutonomyLevel(targetMode)
				fmt.Printf("✓ Autonomy mode switched to: %s\n", config.GetAutonomyLevel())
			}
			return
		case "/sop":
			handleCLISOP(parts[1:], chatIDStr)
			return
		case "/receipts":
			printReceipts()
			return
		case "/session", "/sessions":
			handleCLISession(parts[1:])
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

	agent.EnterAgentMode(currentSessionID)
	agent.RunAgentSessionLoop(currentSessionID, cliChatID, prompt, 0)

	cliMu.Lock()
	isTurnActive = false
	cliMu.Unlock()
}

// printBanner prints a clean, modern startup banner with ANSI colors
func printBanner() {
	primaryModel := "unknown"
	provider := "unknown"
	if m := models.RouteModel("agent"); m != nil {
		primaryModel = m.Model
		provider = m.Provider
	}

	mode := string(config.GetAutonomyLevel())
	modeColor := "\033[33m" // yellow
	if mode == "readonly" {
		modeColor = "\033[36m" // cyan
	} else if mode == "yolo" {
		modeColor = "\033[31m" // red
	}

	sessInfo := ""
	if currentSessionID != "default" {
		sessInfo = fmt.Sprintf(" | 📁 Session: \033[1m%s\033[0m", currentSessionID)
	}

	fmt.Println("\033[1;36m━━━ scorp ━━━\033[0m \033[2mAutonomous AI Coding Agent (v2.0)\033[0m")
	fmt.Printf("\033[32m●\033[0m Model: \033[1m%s\033[0m \033[2m(%s)\033[0m | %s🛡️ Mode: %s\033[0m%s\n", primaryModel, provider, modeColor, mode, sessInfo)
	fmt.Println("\033[2m💡 Type your task (paste with \"\"\" for multiline). Type /help for commands, /exit to quit.\033[0m")
}

// printCLIHelp displays commands and usage
func printCLIHelp() {
	fmt.Print(`
Commands:
  /help          — Show this help message
  /models        — List all available models and their status
  /model <name>  — Switch active agent model (e.g. /model deepseek/deepseek-v4-flash)
  /mode <level>  — View or set autonomy mode (readonly, supervised, yolo)
  /session       — Manage chat sessions (list, new, use, rename, delete)
  /sop [run]     — List or run Standard Operating Procedures
  /receipts      — View recent cryptographic tool execution receipts
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
  scorp> /sop run health_audit
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
	fmt.Printf("  • Active Session:    %s\n", currentSessionID)
	fmt.Printf("  • Autonomy Level:    %s\n", config.GetAutonomyLevel())
	fmt.Printf("  • Session Chat ID:   %s\n", chatIDStr)
	fmt.Printf("  • Agent Mode:        %v\n", agent.IsAgentMode(chatIDStr))
	if agent.HasPendingConfirmation(chatIDStr) {
		cmd, tool, _ := agent.GetPendingConfirmationDetails(chatIDStr)
		fmt.Printf("  • Pending Confirm:   tool=%s cmd=%s\n", tool, cmd)
	}
}

func handleCLISOP(args []string, chatIDStr string) {
	if len(args) == 0 || args[0] == "list" {
		sops := sop.ListSOPs()
		fmt.Println("\n📋 Standard Operating Procedures (SOPs):")
		for _, s := range sops {
			fmt.Printf("  • %-18s — %s\n", s.Name, s.Description)
		}
		fmt.Println("To run an SOP: /sop run <name>")
		return
	}

	if args[0] == "run" && len(args) > 1 {
		sopName := args[1]
		s, err := sop.GetSOP(sopName)
		if err != nil {
			fmt.Printf("❌ %v\n", err)
			return
		}
		fmt.Printf("🚀 Executing SOP '%s'...\n", s.Name)
		executeTurn(chatIDStr, s.Prompt)
		return
	}

	fmt.Println("Usage: /sop [list|run <name>]")
}

func handleCLISession(args []string) {
	if len(args) == 0 || args[0] == "list" {
		sessions := agent.ListSessions()
		fmt.Printf("\n📋 Saved Chat Sessions (%d total):\n", len(sessions))
		fmt.Printf("  %-18s %-10s %-16s %-30s\n", "NAME", "MESSAGES", "LAST MODIFIED", "PREVIEW")
		fmt.Println("  " + strings.Repeat("─", 78))

		foundCurrent := false
		for _, s := range sessions {
			marker := " "
			if s.ID == currentSessionID {
				marker = "\033[1;32m●\033[0m"
				foundCurrent = true
			}
			fmt.Printf("%s %-18s %-10s %-16s \033[2m%-30s\033[0m\n",
				marker, s.ID, fmt.Sprintf("%d msgs", s.MsgCount), s.LastModified.Format("02 Jan 15:04"), s.LastPreview)
		}
		if !foundCurrent {
			fmt.Printf("\033[1;32m●\033[0m %-18s %-10s %-16s \033[2m%-30s\033[0m\n",
				currentSessionID, "0 msgs", time.Now().Format("02 Jan 15:04"), "(current active session)")
		}

		fmt.Println("\nCommands:")
		fmt.Println("  /session list                — List all sessions")
		fmt.Println("  /session new <name>          — Create & switch to a new session")
		fmt.Println("  /session use <name>          — Switch to existing session")
		fmt.Println("  /session rename <old> <new>  — Rename a session")
		fmt.Println("  /session delete <name>       — Delete a session")
		fmt.Println()
		return
	}

	cmd := args[0]
	switch cmd {
	case "new":
		if len(args) < 2 {
			fmt.Println("Usage: /session new <name>")
			return
		}
		name := strings.TrimSpace(args[1])
		currentSessionID = name
		fmt.Printf("✓ Switched to new session: \033[1m%s\033[0m\n", currentSessionID)

	case "use", "switch":
		if len(args) < 2 {
			fmt.Println("Usage: /session use <name>")
			return
		}
		name := strings.TrimSpace(args[1])
		currentSessionID = name
		fmt.Printf("✓ Active session switched to: \033[1m%s\033[0m\n", currentSessionID)

	case "rename":
		if len(args) < 3 {
			fmt.Println("Usage: /session rename <old> <new>")
			return
		}
		oldName := strings.TrimSpace(args[1])
		newName := strings.TrimSpace(args[2])
		if err := agent.RenameSession(oldName, newName); err != nil {
			fmt.Printf("❌ %v\n", err)
			return
		}
		if currentSessionID == oldName {
			currentSessionID = newName
		}
		fmt.Printf("✓ Session renamed from '%s' to '%s'\n", oldName, newName)

	case "delete", "rm":
		if len(args) < 2 {
			fmt.Println("Usage: /session delete <name>")
			return
		}
		name := strings.TrimSpace(args[1])
		if err := agent.DeleteSession(name); err != nil {
			fmt.Printf("❌ %v\n", err)
			return
		}
		if currentSessionID == name {
			currentSessionID = "default"
		}
		fmt.Printf("✓ Session '%s' deleted. Active session is now 'default'.\n", name)

	default:
		fmt.Printf("Unknown session action '%s'. Use /session list, new, use, rename, or delete.\n", cmd)
	}
}

func printReceipts() {
	receipts := tools.GetRecentReceipts()
	if len(receipts) == 0 {
		fmt.Println("No tool execution receipts recorded yet.")
		return
	}
	fmt.Printf("\n🧾 Recent Tool Execution Receipts (%d total):\n", len(receipts))
	fmt.Printf("%-18s %-20s %-10s %-18s\n", "RECEIPT ID", "TOOL", "STATUS", "TIME")
	fmt.Println(strings.Repeat("─", 70))
	for _, r := range receipts {
		status := "✓ OK"
		if !r.Success {
			status = "❌ FAIL"
		}
		fmt.Printf("%-18s %-20s %-10s %-18s\n", r.ReceiptID, r.Tool, status, r.Timestamp.Format("15:04:05"))
	}
	fmt.Println()
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
