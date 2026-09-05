package telegram

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"scorp-agent/agent"
	"scorp-agent/bootstrap"
	"scorp-agent/browser"
	"scorp-agent/collectors"
	"scorp-agent/config"
	"scorp-agent/mcp"
	"scorp-agent/mcp/marketplace"
	"scorp-agent/metrics"
	"scorp-agent/models"
	"scorp-agent/rag"
	"scorp-agent/scheduler"
	"scorp-agent/session"
	"scorp-agent/skills"
	"scorp-agent/sop"
	"scorp-agent/testutil"
	"scorp-agent/tools"
	"scorp-agent/wizard"
)

// ──────────────────────────────────────────────
// Telegram 24/7 Daemon & Action Handler
// Modularized from main.go
// ──────────────────────────────────────────────

// StartDaemon initializes all background subsystems and runs the Telegram bot.
func StartDaemon() {
	// Wire callback for models package
	models.UpdateEnvFile = func(key, value string) {
		wizard.UpdateEnvFile(key, value)
	}

	// Wire callback for browser package
	browser.SendFile = func(chatID, filePath string) bool {
		return SendFile(chatID, filePath)
	}

	// Initialize unified config manager
	config.InitConfigManager()

	InitTelegram()

	// Wire callback for telegram actions
	HandleAction = HandleTelegramAction

	// Wire tools package callbacks (avoid import cycles)
	tools.TgPost = func(method string, payload map[string]interface{}) (tools.TgResponse, error) {
		resp, err := TgPost(method, payload)
		if err != nil || resp == nil {
			return tools.TgResponse{}, err
		}
		return tools.TgResponse{OK: resp.OK, Description: resp.Description, Result: nil}, nil
	}

	tools.SendMessage = func(text string, keyboard map[string]interface{}) bool {
		return SendMessage(text, keyboard)
	}
	tools.SendMessageGetID = func(text string, chatID int64) int64 {
		return SendMessageGetID(text, chatID)
	}
	tools.SendMessageGetIDWithKeyboard = func(text string, chatID int64, keyboard map[string]interface{}) int64 {
		return SendMessageGetIDWithKeyboard(text, chatID, keyboard)
	}
	tools.EditMessageByID = func(chatID int64, messageID int64, text string, keyboard map[string]interface{}) bool {
		return EditMessageByID(chatID, messageID, text, keyboard)
	}
	tools.SendChatAction = func(chatID int64, action string) {
		SendChatAction(chatID, action)
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

	// Load dynamic skills
	skills.EnsureBuiltinSkills()
	skills.LoadAllSkills()
	skills.Load()

	// Bootstrap tool registry
	bootstrap.RegisterAutonomous()

	// Start CPU/Docker stats samplers
	var wg sync.WaitGroup
	done := make(chan struct{})
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	collectors.StartCPUSampler(done)
	collectors.StartDockerStatsSampler(done)

	// Webhook or long-polling mode
	if config.Cfg.TelegramWebhookURL != "" {
		if err := StartWebhookServer(config.Cfg.TelegramWebhookURL); err != nil {
			log.Fatalf("Failed to start webhook: %v", err)
		}
		log.Println("Webhook mode enabled")
	} else {
		wg.Add(1)
		go func() {
			defer wg.Done()
			runCommandLoop(done)
		}()
		log.Println("Long polling mode enabled")
	}

	// Session cleanup loop
	wg.Add(1)
	go func() {
		defer wg.Done()
		agent.CleanupSessionsLoop(done)
	}()

	// Start MCP servers
	log.Println("[main] Starting MCP servers...")
	mcp.StartMCPServers()
	log.Println("[main] MCP servers started")

	// Start test endpoint
	testutil.StartTestEndpoint()

	// Start MCP server mode
	wg.Add(1)
	go func() {
		defer wg.Done()
		mcp.StartMCPServerMode()
	}()

	// Start Cron scheduler
	scheduler.SendMessage = agent.SendMessageSmart
	scheduler.SendMessageGetID = SendMessageGetID
	scheduler.TgPost = func(method string, payload map[string]interface{}) {
		TgPost(method, payload)
	}
	scheduler.RunAgentLoop = agent.RunAgentLoop
	scheduler.LoadTasks()
	wg.Add(1)
	go func() {
		defer wg.Done()
		scheduler.Loop(done)
	}()

	// Models & RAG & Session DB
	models.LoadModelConfig()
	models.InitModelUsage()
	models.LoadCostConfig()
	models.LoadCostTracker()
	tools.InitMemoryCache()
	session.InitSessionDB()
	rag.InitRAG()
	rag.InitVectorRAG()
	sop.InitDefaultSOPs()

	// Wire ExecuteTool callback for tools bridging
	tools.ExecuteTool = func(tc models.ToolCall, chatID int64) (string, bool) {
		return agent.ExecuteTool(tc, chatID)
	}

	// Start autonomous agent loop
	wg.Add(1)
	go func() {
		defer wg.Done()
		agent.AutonomousLoop(done)
	}()

	// Start metrics Prometheus server
	wg.Add(1)
	go func() {
		defer wg.Done()
		metrics.StartServer()
	}()

	// Wait for shutdown signal
	<-stop
	log.Println("Shutting down Scorp Telegram Daemon...")
	close(done)

	metrics.StopServer()
	StopWebhookServer()
	mcp.StopMCPServerMode()
	mcp.StopMCPServers()
	SendMessage("🔴 <b>Scorp Agent Stopped</b>", nil)
	wg.Wait()
	log.Println("Goodbye.")
}

func runCommandLoop(done chan struct{}) {
	time.Sleep(3 * time.Second)
	log.Println("[telegram] Command loop started")
	SetupBotCommands()
	os.MkdirAll(config.UploadsDir(), 0755)

	for {
		select {
		case <-done:
			return
		default:
		}

		commands, callbacks, documents, inlineQueries := PollUpdates()

		for _, cmd := range commands {
			log.Printf("Command: %s from %d", cmd.Text, cmd.ChatID)
			HandleTelegramAction(cmd.Text, cmd.ChatID, cmd.MsgID, "")
		}

		for _, cb := range callbacks {
			log.Printf("Callback: %s from %d", cb.Data, cb.ChatID)
			AnswerCallback(cb.CBID, "")
			HandleTelegramAction(cb.Data, cb.ChatID, cb.MsgID, cb.CBID)
		}

		for _, doc := range documents {
			log.Printf("Upload received: %s (%s, photo=%v)", doc.FileName, HumanSize(doc.FileSize), doc.IsPhoto)
			go agent.HandleUploadInAgentMode(doc)
		}

		for _, iq := range inlineQueries {
			log.Printf("Inline query: %s from %d", iq.Query, iq.UserID)
			tools.HandleInlineQuery(iq)
		}
	}
}

// HandleTelegramAction dispatches all incoming Telegram button clicks, slash commands, and user messages.
func HandleTelegramAction(action string, chatID int64, messageID int64, callbackID string) {
	isCallback := callbackID != ""
	edit := isCallback && messageID != 0
	chatIDStr := fmt.Sprintf("%d", chatID)

	switch {
	case action == "/start" || action == "mn:main":
		text := "🤖 <b>Scorp Agent (v2.0)</b>\n\n" +
			"💬 Type anything to chat with AI (agent-first)\n" +
			"🛠 <code>/agent</code> for autonomous shell, file, and web access\n" +
			"⏰ <code>/cron</code> to view and manage background tasks\n" +
			"━━━━━━━━━━━━━━━━━━━"
		kb := MainMenuKeyboard()
		if edit {
			EditMessage(chatID, messageID, text, kb)
		} else {
			SendMessage(text, kb)
		}

	case action == "mn:sys":
		text := "🔧 <b>System &amp; Tools</b>"
		kb := SystemMenuKeyboard()
		if edit {
			EditMessage(chatID, messageID, text, kb)
		} else {
			SendMessage(text, kb)
		}

	case action == "mn:set":
		text := SettingsMenuText()
		kb := SettingsMenuKeyboard()
		if edit {
			EditMessage(chatID, messageID, text, kb)
		} else {
			SendMessage(text, kb)
		}

	case strings.HasPrefix(action, "set:"):
		AnswerCallback(callbackID, "⚙️ Settings updated")

	case action == "/status" || action == "status":
		sys := collectors.CollectSystem()
		docker := collectors.CollectDocker()
		stor := collectors.CollectStorage()
		text := scheduler.FormatStatusResponse(sys, docker, stor)
		kb := BackAndRefreshKeyboard("status")
		if edit {
			EditMessage(chatID, messageID, text, kb)
		} else {
			SendMessage(text, kb)
		}

	case action == "/containers" || action == "containers":
		docker := collectors.CollectDocker()
		text := scheduler.SectionDocker(docker)
		kb := BackAndRefreshKeyboard("containers")
		if edit {
			EditMessage(chatID, messageID, text, kb)
		} else {
			SendMessage(text, kb)
		}

	case action == "/coolify" || action == "coolify":
		coolify := collectors.CollectCoolify()
		text := scheduler.SectionCoolify(coolify)
		kb := BackAndRefreshKeyboard("coolify")
		if edit {
			EditMessage(chatID, messageID, text, kb)
		} else {
			SendMessage(text, kb)
		}

	case action == "/security" || action == "security":
		sec := collectors.CollectSecurityWithPeek()
		text := scheduler.SectionSecurity(sec)
		kb := BackAndRefreshKeyboard("security")
		if edit {
			EditMessage(chatID, messageID, text, kb)
		} else {
			SendMessage(text, kb)
		}

	case action == "/storage" || action == "storage":
		stor := collectors.CollectStorage()
		text := scheduler.SectionStorage(stor)
		kb := BackAndRefreshKeyboard("storage")
		if edit {
			EditMessage(chatID, messageID, text, kb)
		} else {
			SendMessage(text, kb)
		}

	case action == "/network" || action == "network":
		net := collectors.CollectNetwork()
		text := scheduler.SectionNetwork(net)
		kb := BackAndRefreshKeyboard("network")
		if edit {
			EditMessage(chatID, messageID, text, kb)
		} else {
			SendMessage(text, kb)
		}

	case action == "/files" || action == "files":
		text := "📂 <b>File Manager</b>\n\nChoose a directory to browse:"
		kb := RootsKeyboard()
		if edit {
			EditMessage(chatID, messageID, text, kb)
		} else {
			SendMessage(text, kb)
		}

	case strings.HasPrefix(action, "fb:"):
		pid := action[3:]
		path := GetPath(pid)
		if path != "" {
			text, kb := DirKeyboard(path)
			if edit {
				EditMessage(chatID, messageID, text, kb)
			} else {
				SendMessage(text, kb)
			}
		}

	case strings.HasPrefix(action, "fd:"):
		pid := action[3:]
		path := GetPath(pid)
		if path != "" {
			text, kb := FileDetailKeyboard(path)
			if edit {
				EditMessage(chatID, messageID, text, kb)
			} else {
				SendMessage(text, kb)
			}
		}

	case strings.HasPrefix(action, "zp:"):
		pid := action[3:]
		path := GetPath(pid)
		if path != "" {
			text, kb := FolderZipInfo(path)
			if edit {
				EditMessage(chatID, messageID, text, kb)
			} else {
				SendMessage(text, kb)
			}
		}

	case strings.HasPrefix(action, "zc:"):
		pid := action[3:]
		path := GetPath(pid)
		if path != "" {
			if isCallback {
				AnswerCallback(callbackID, "📦 Creating ZIP...")
			}
			SendFolderAsZip(config.Cfg.TelegramChatID, path)
		}

	case strings.HasPrefix(action, "dl:"):
		pid := action[3:]
		path := GetPath(pid)
		if path != "" {
			if isCallback {
				AnswerCallback(callbackID, "⬇️ Sending file...")
			}
			if !SendFile(config.Cfg.TelegramChatID, path) {
				SendMessage("❌ Failed to send file (may be too large)", nil)
			}
		}

	case action == "/mcp" || strings.HasPrefix(action, "/mcp "):
		sub := strings.Fields(action)
		if len(sub) > 2 && sub[1] == "restart" {
			target := sub[2]
			if err := mcp.RestartServer(target); err != nil {
				SendMessage(fmt.Sprintf("❌ Failed to restart MCP server '%s': %v", target, err), nil)
			} else {
				SendMessage(fmt.Sprintf("✓ MCP server <code>%s</code> successfully restarted!", target), nil)
			}
		} else if len(sub) >= 2 && sub[1] == "search" {
			term := strings.Join(sub[2:], " ")
			out, err := marketplace.CLISearch(term)
			if err != nil {
				SendMessage(fmt.Sprintf("❌ Marketplace search failed: %v", err), BackButtonKeyboard())
			} else {
				SendMessage(out, BackButtonKeyboard())
			}
		} else if len(sub) >= 3 && sub[1] == "install" {
			name := sub[2]
			opt := strings.Join(sub[3:], " ") // keep multi-token specs intact (e.g. "upstream npx:pkg /tmp")
			handleTelegramMarketplaceInstall(chatID, messageID, name, opt, isCallback)
		} else {
			SendMessage(mcp.GetServerHealthStatus(), BackButtonKeyboard())
		}

	case strings.HasPrefix(action, "mcpm:"):
		// Marketplace install callbacks: mcpm:<name>:<1|2|3>
		payload := strings.TrimPrefix(action, "mcpm:")
		seg := strings.SplitN(payload, ":", 2)
		name := seg[0]
		opt := ""
		if len(seg) > 1 {
			opt = seg[1]
		}
		AnswerCallback(callbackID, "🛒 Processing…")
		handleTelegramMarketplaceInstall(chatID, messageID, name, opt, edit)

	case action == "/skills" || action == "/skill" || strings.HasPrefix(action, "/skill "):
		sub := strings.Fields(action)
		if len(sub) > 1 && sub[1] != "list" {
			target := sub[1]
			if _, err := skills.ActivateSkill(target, 5); err != nil {
				SendMessage(fmt.Sprintf("❌ Failed to activate skill: %v", err), nil)
			} else {
				SendMessage(fmt.Sprintf("✅ Skill <code>%s</code> activated for 5 turns!", target), nil)
			}
		} else {
			SendMessage(skills.ListSkillsOverview(), BackButtonKeyboard())
		}

	case action == "/sops":
		text, _ := tools.ExecuteSOP(map[string]interface{}{"action": "list"})
		SendMessage(text, BackButtonKeyboard())

	case action == "/cron":
		SendMessage(scheduler.FormatTasksList(), nil)

	case action == "/usage":
		SendMessage(models.FormatUsageStats(), BackButtonKeyboard())

	case action == "/session" || action == "/sessions" || action == "sessions":
		text := FormatSessionMenuText(chatIDStr)
		kb := BuildSessionMenuKeyboard(chatIDStr)
		if edit {
			EditMessage(chatID, messageID, text, kb)
		} else {
			SendMessage(text, kb)
		}

	case strings.HasPrefix(action, "/session "):
		subArgs := strings.Fields(strings.TrimPrefix(action, "/session "))
		if len(subArgs) == 0 {
			SendMessage(FormatSessionMenuText(chatIDStr), BuildSessionMenuKeyboard(chatIDStr))
			break
		}
		switch subArgs[0] {
		case "list":
			SendMessage(FormatSessionMenuText(chatIDStr), BuildSessionMenuKeyboard(chatIDStr))
		case "new":
			newSess := ""
			if len(subArgs) >= 2 {
				newSess = strings.TrimSpace(subArgs[1])
			} else {
				// Auto-generate temporary name that will be auto-titled on turn 1
				newSess = fmt.Sprintf("chat-%s", time.Now().Format("0102-150405"))
			}
			SetActiveSessionID(chatIDStr, newSess)
			SendMessage(fmt.Sprintf("✓ Started new conversation: 🟢 <b>%s</b>\n<i>This session will be auto-named after your first message.</i>", newSess), BuildSessionMenuKeyboard(chatIDStr))
		case "switch", "use":
			if len(subArgs) < 2 {
				SendMessage("⚠️ Usage: <code>/session use &lt;name&gt;</code>", nil)
				break
			}
			newSess := strings.TrimSpace(subArgs[1])
			SetActiveSessionID(chatIDStr, newSess)
			SendMessage(fmt.Sprintf("✓ Active session switched to: 🟢 <b>%s</b>", newSess), BuildSessionMenuKeyboard(chatIDStr))
		case "rename":
			if len(subArgs) < 3 {
				SendMessage("⚠️ Usage: <code>/session rename &lt;old&gt; &lt;new&gt;</code>", nil)
				break
			}
			oldName := strings.TrimSpace(subArgs[1])
			newName := strings.TrimSpace(subArgs[2])
			if err := agent.RenameSession(oldName, newName); err != nil {
				SendMessage(fmt.Sprintf("❌ %v", err), nil)
				break
			}
			if GetActiveSessionID(chatIDStr) == oldName {
				SetActiveSessionID(chatIDStr, newName)
			}
			SendMessage(fmt.Sprintf("✓ Session renamed from <code>%s</code> to <code>%s</code>", oldName, newName), BuildSessionMenuKeyboard(chatIDStr))
		case "delete", "rm":
			if len(subArgs) < 2 {
				SendMessage("⚠️ Usage: <code>/session delete &lt;name&gt;</code>", nil)
				break
			}
			targetName := strings.TrimSpace(subArgs[1])
			if err := agent.DeleteSession(targetName); err != nil {
				SendMessage(fmt.Sprintf("❌ %v", err), nil)
				break
			}
			if GetActiveSessionID(chatIDStr) == targetName {
				SetActiveSessionID(chatIDStr, chatIDStr)
			}
			SendMessage(fmt.Sprintf("✓ Session <code>%s</code> deleted.", targetName), BuildSessionMenuKeyboard(chatIDStr))
		default:
			SendMessage(FormatSessionMenuText(chatIDStr), BuildSessionMenuKeyboard(chatIDStr))
		}

	case strings.HasPrefix(action, "sess:"):
		if isCallback {
			AnswerCallback(callbackID, "")
		}
		text, kb, handled := HandleSessionCallback(action, chatIDStr)
		if handled && text != "" {
			if edit {
				EditMessage(chatID, messageID, text, kb)
			} else {
				SendMessage(text, kb)
			}
		}

	case action == "/compact":
		activeSess := GetActiveSessionID(chatIDStr)
		SendMessage("🗜️ <i>Compacting conversation history in background...</i>", nil)
		go func(sessName string) {
			stats := agent.CompactSessionHistory(sessName)
			SendMessage(agent.FormatCompactStats(sessName, stats), BuildSessionMenuKeyboard(chatIDStr))
		}(activeSess)

	case action == "/clear":
		agent.ClearChatSession(chatIDStr)
		SendMessage("🧹 Chat history cleared.", nil)

	case action == "/stop":
		agent.ExitAgentMode(chatIDStr)
		SendMessage("🛑 Agent mode stopped. Back to normal chat.", nil)

	case action == "/agent":
		agent.EnterAgentMode(chatIDStr)
		SendMessage("🛠 <b>Agent Mode Activated</b>\n\nI can execute shell commands, inspect files, and automate tasks.\nWhat should I do?", nil)

	case action == "/confirm_yes" || action == "confirm_yes":
		agent.HandleConfirmation(chatID, true, messageID)

	case action == "/confirm_no" || action == "confirm_no":
		agent.HandleConfirmation(chatID, false, messageID)

	case action == "/model" || action == "/model list":
		if isCallback && edit {
			EditMessage(chatID, messageID, wizard.ModelMenuText(), wizard.ModelMenuKeyboard())
		} else {
			SendMessage(wizard.ModelMenuText(), wizard.ModelMenuKeyboard())
		}

	case strings.HasPrefix(action, "mdl:"):
		if isCallback {
			AnswerCallback(callbackID, "")
		}
		text, kb, handled := wizard.HandleModelCallback(action, chatID, messageID)
		if handled && text != "" {
			if edit {
				EditMessage(chatID, messageID, text, kb)
			} else {
				SendMessage(text, kb)
			}
		}

	case strings.HasPrefix(action, "cron:"):
		parts := strings.Split(action, ":")
		if len(parts) >= 3 {
			cmd := parts[1]
			taskID := parts[2]
			switch cmd {
			case "run":
				task := scheduler.GetTask(taskID)
				if task != nil {
					go scheduler.RunTask(*task)
					AnswerCallback(callbackID, "Executing task...")
				}
			case "del":
				if scheduler.RemoveTask(taskID) {
					AnswerCallback(callbackID, "Task deleted")
					SendMessage(scheduler.FormatTasksList(), nil)
				}
			case "pause":
				if scheduler.ToggleTask(taskID, false) {
					AnswerCallback(callbackID, "Task paused")
					SendMessage(scheduler.FormatTasksList(), nil)
				}
			case "resume":
				if scheduler.ToggleTask(taskID, true) {
					AnswerCallback(callbackID, "Task resumed")
					SendMessage(scheduler.FormatTasksList(), nil)
				}
			}
		}

	case action == "help" || action == "/help":
		helpText := "📖 <b>Scorp Agent Commands</b>\n\n" +
			"/agent — Enter autonomous agent mode\n" +
			"/model — Manage AI models & providers\n" +
			"/cron — View & manage scheduled cron tasks\n" +
			"/sops — View standard operating procedures\n" +
			"/mcp — Model Context Protocol servers\n" +
			"/status — System & container status\n" +
			"/usage — Token usage & cost tracking\n" +
			"/clear — Clear conversation history\n" +
			"/stop — Stop running task"
		SendMessage(helpText, BackButtonKeyboard())

	default:
		// Regular user text message
		// Check if user is in model wizard
		if wizard.GetModelWizard(chatID) != nil && !strings.HasPrefix(action, "/") {
			if wizard.HandleModelWizardTextRouter(action, chatID) {
				return
			}
		}
		if action == "/cancel" && wizard.GetModelWizard(chatID) != nil {
			wizard.HandleModelWizardTextRouter(action, chatID)
			return
		}

		// Check if user is clarifying a question
		activeSess := GetActiveSessionID(chatIDStr)
		if tools.HasPendingClarify(chatIDStr) {
			tools.ResolveClarify(action, chatIDStr, "")
			return
		}
		if tools.HasPendingClarify(activeSess) {
			tools.ResolveClarify(action, activeSess, "")
			return
		}

		// Real-time Steering Queue check: if agent loop is active for this chat, queue as steering message!
		if agent.IsLoopActive(chatIDStr) || agent.IsLoopActive(activeSess) {
			agent.QueueSteeringMessage(chatIDStr, action)
			agent.QueueSteeringMessage(activeSess, action)
			SendMessage("⚡ <i>Instruction received — steering agent mid-run...</i>", nil)
			return
		}

		// Check if session is in agent mode, default to agent loop using current active session!
		agent.EnterAgentMode(activeSess)
		go agent.RunAgentSessionLoop(activeSess, chatID, action, 0)
	}
}
