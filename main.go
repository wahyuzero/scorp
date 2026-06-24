package main

import (
	"scorp-agent/agent"
	"scorp-agent/bootstrap"
	"scorp-agent/browser"
	"scorp-agent/mcp"
	"scorp-agent/models"
	"scorp-agent/skills"
	"scorp-agent/testutil"
	"scorp-agent/telegram"
	"scorp-agent/tools"
	"scorp-agent/wizard"
	"scorp-agent/loop"
	"scorp-agent/metrics"
	"scorp-agent/collectors"
	"scorp-agent/config"
	"scorp-agent/session"
	"scorp-agent/rag"
	"scorp-agent/scheduler"
	"scorp-agent/updater"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"
)
func main() {
	log.SetFlags(log.Ldate | log.Ltime)

	// Update subcommand: scorp update
	if len(os.Args) > 1 && os.Args[1] == "update" {
		msg, err := updater.SelfUpdate()
		if err != nil {
			fmt.Printf("❌ %v\n", err)
			os.Exit(1)
		}
		fmt.Println(msg)
		return
	}

	// Version subcommand: scorp version
	if len(os.Args) > 1 && (os.Args[1] == "version" || os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("scorp %s (%s/%s)\n", updater.Version, runtime.GOOS, runtime.GOARCH)
		repo := os.Getenv("GITHUB_REPO")
		if repo != "" {
			fmt.Printf("repo: %s\n", repo)
		}
		return
	}

	// MCP-only mode: skip everything, just serve MCP over stdio
	if len(os.Args) > 1 && os.Args[1] == "--mcp-server" {
		log.Println("[mcp] MCP-only mode starting...")
		mcp.StartMCPServerMode()
		return
	}

	// CLI mode: --cli flag OR no Telegram token configured
	if isCLIMode() {
		if err := config.LoadConfig(); err != nil {
			// CLI mode can work without TELEGRAM_* vars
			config.Cfg.TelegramBotToken = ""
			config.Cfg.TelegramChatID = ""
		}
		startCLI()
		return
	}

	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Config error: %v", err)
	}

	// Wire callback for models package
	models.UpdateEnvFile = func(key, value string) {
		wizard.UpdateEnvFile(key, value)
	}

	// Wire callback for browser package
	browser.SendFile = func(chatID, filePath string) bool {
		return telegram.SendFile(chatID, filePath)
	}

	// Initialize unified config manager
	config.InitConfigManager()

	telegram.InitTelegram()

	// Wire callback for telegram → main.go
	telegram.HandleAction = handleAction

	// Wire tools package callbacks (avoid import cycles)
	tools.TgPost = func(method string, payload map[string]interface{}) (tools.TgResponse, error) {
		resp, err := telegram.TgPost(method, payload)
		if err != nil || resp == nil {
			return tools.TgResponse{}, err
		}
		// telegram.go tgResponse only has OK, Description — Result is always nil
		return tools.TgResponse{OK: resp.OK, Description: resp.Description, Result: nil}, nil
	}

	// Telegram callbacks → telegram.go functions
	tools.SendMessage = func(text string, keyboard map[string]interface{}) bool {
		return telegram.SendMessage(text, keyboard)
	}
	tools.SendMessageGetID = func(text string, chatID int64) int64 {
		return telegram.SendMessageGetID(text, chatID)
	}
	tools.EditMessageByID = func(chatID int64, messageID int64, text string, keyboard map[string]interface{}) bool {
		return telegram.EditMessageByID(chatID, messageID, text, keyboard)
	}
	tools.SendChatAction = func(chatID int64, action string) {
		telegram.SendChatAction(chatID, action)
	}

	// Agent callbacks → agent package
	tools.StorePendingConfirmation = func(chatID, toolName, command string, _ []tools.AgentMessage) {
		agent.StorePendingConfirmation(chatID, toolName, command, nil)
	}
	tools.IsDangerousCommand = func(cmd string) bool {
		return agent.IsDangerousCommand(cmd)
	}

	// Autonomous callbacks → agent package variables (type alias makes this work)
	tools.AutoConfig = &agent.AutoConfig
	tools.AutoMu = &agent.AutoMu
	tools.AutoLog = &agent.AutoLog
	tools.AutoKillFile = agent.AutoKillFile
	tools.AutoCycleNum = &agent.AutoCycleNum
	tools.SaveAutonomousConfig = agent.SaveAutonomousConfig
	tools.SetKillSwitch = agent.SetKillSwitch
	tools.RunAutonomousCycle = agent.RunAutonomousCycle

	// Load dynamic skills from JSON files
	skills.Load()

	// Init credential vault
	// Handled by bootstrap/vault.go init() (runs before main)
	
	// Init browser monitor (scheduled scraping + change detection)
	tools.InitMonitor()

	// Init autonomous agent (Phase 7)
	agent.LoadAutonomousConfig()
	agent.LoadAutoLog()
	log.Printf("[autonomous] Config loaded: enabled=%v interval=%v approval=%s (cycles=%d actions=%d)",
		agent.AutoConfig.Enabled, agent.AutoConfig.Interval,
		agent.AutoConfig.ApprovalLevel, agent.AutoConfig.TotalCycles, agent.AutoConfig.TotalActions)
	bootstrap.RegisterAutonomous()

	// Periodic browser session cleanup (every 5 min)
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			browser.CleanupStaleBrowserSessions()
		}
	}()

	log.Println("==================================================")
	log.Println("Scorp Agent starting...")
	log.Printf("Telegram Chat ID: %s", config.Cfg.TelegramChatID)
	log.Printf("Report interval: %ds", config.Cfg.ReportInterval)
	log.Println("==================================================")

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	var wg sync.WaitGroup
	done := make(chan struct{})

	// Start background samplers
	collectors.StartCPUSampler(done)
	collectors.StartDockerStatsSampler(done)

	// Start webhook server if configured, otherwise use long polling
	if config.Cfg.TelegramWebhookURL != "" {
		if err := telegram.StartWebhookServer(config.Cfg.TelegramWebhookURL); err != nil {
			log.Fatalf("Failed to start webhook: %v", err)
		}
		log.Println("Webhook mode enabled")
	} else {
		// 5. Command/callback poller (long polling mode)
		wg.Add(1)
		go func() {
			defer wg.Done()
			commandLoop(done)
		}()
		log.Println("Long polling mode enabled")
	}

	// 1. Hourly report loop (only if scheduled reports enabled)
	if config.Cfg.ScheduledReportsEnabled {
		loop.RepLoopCtl.Start(func(d chan struct{}) { hourlyLoop(d) })
		log.Println("[main] Scheduled reports: ON")
	} else {
		log.Println("[main] Scheduled reports: OFF (toggle via SCHEDULED_REPORTS_ENABLED or Settings)")
	}

	// 2. Resource alert loop (only if monitoring enabled)
	if config.Cfg.MonitoringEnabled {
		loop.MonLoopCtl.Start(func(d chan struct{}) { resourceAlertLoop(d) })
		log.Println("[main] Monitoring (resource alerts): ON")
	} else {
		log.Println("[main] Monitoring (resource alerts): OFF (toggle via MONITORING_ENABLED or Settings)")
	}

	// 3. Journal watcher (always on — feeds security event processor)
	wg.Add(1)
	go func() {
		defer wg.Done()
		collectors.WatchJournal(done)
	}()

	// 4. Security event processor (only if security alerts enabled)
	if config.Cfg.SecurityAlertsEnabled {
		loop.SecLoopCtl.Start(func(d chan struct{}) { securityEventLoop(d) })
		log.Println("[main] Security alerts: ON")
	} else {
		log.Println("[main] Security alerts: OFF (toggle via SECURITY_ALERTS_ENABLED or Settings)")
	}

	// 5. Command/callback poller (handled above in webhook/long-polling branch)

	// 6. Session cleanup loop
	wg.Add(1)
	go func() {
		defer wg.Done()
		agent.CleanupSessionsLoop(done)
	}()

	// 7. Start MCP servers (blocking - tools must be registered before LLM calls)
	log.Println("[main] Starting MCP servers...")
	mcp.StartMCPServers()
	log.Println("[main] MCP servers started")

	// 7b. Start test endpoint (localhost only)
	testutil.StartTestEndpoint()

	// 8. Start MCP server mode (scorp-agent as MCP server)
	wg.Add(1)
	go func() {
		defer wg.Done()
		mcp.StartMCPServerMode()
	}()
	// 8. Start scheduler
	// Wire callbacks (scheduler package can't import main)
	scheduler.SendMessage = agent.SendMessageSmart
	scheduler.SendMessageGetID = telegram.SendMessageGetID
	scheduler.TgPost = func(method string, payload map[string]interface{}) {
		telegram.TgPost(method, payload) // ignore return
	}
	scheduler.RunAgentLoop = agent.RunAgentLoop
	scheduler.LoadTasks()
	wg.Add(1)
	go func() {
		defer wg.Done()
		scheduler.Loop(done)
	}()

	// 9. Load model router
	models.LoadModelConfig()
	models.InitModelUsage()
	models.LoadCostConfig()
	models.LoadCostTracker()

	// 12. Initialize memory cache
	tools.InitMemoryCache()

	// 12b. Initialize session search DB (SQLite + FTS5)
	session.InitSessionDB()

	// 13. Initialize RAG index
	rag.InitRAG()
	rag.InitVectorRAG()

	// 13b. Start autonomous agent loop (Phase 7)
	wg.Add(1)
	go func() {
		defer wg.Done()
		agent.AutonomousLoop(done)
	}()

	// 14. Start pprof server (localhost only, for debugging)
	go func() {
		log.Println("[pprof] Listening on localhost:6060")
		if err := http.ListenAndServe("127.0.0.1:6060", nil); err != nil {
			log.Printf("[pprof] Failed to start: %v", err)
		}
	}()

	// 15. Start Prometheus metrics server
	metrics.StartServer()

	// Setup inline query mode
	tools.SetupInlineMode()

	// 16. Start uptime monitor
	wg.Add(1)
	go func() {
		defer wg.Done()
		tools.UptimeLoop(done)
	}()

	// Wait for shutdown signal
	<-stop
	log.Println("Shutting down...")
	close(done)

	// Stop metrics server
	metrics.StopServer()

	// Stop webhook server if running
	telegram.StopWebhookServer()

	// Stop MCP server mode
	mcp.StopMCPServerMode()

	// Stop MCP servers
	mcp.StopMCPServers()
	telegram.SendMessage("🔴 <b>Scorp Agent Stopped</b>", nil)
	wg.Wait()
	log.Println("Goodbye.")
}

// isCLIMode returns true if running in CLI mode (--cli flag or no Telegram token)
func isCLIMode() bool {
	if len(os.Args) > 1 && os.Args[1] == "--cli" {
		return true
	}
	// If no token in config or env, default to CLI
	if config.Cfg.TelegramBotToken == "" && os.Getenv("TELEGRAM_BOT_TOKEN") == "" {
		return true
	}
	return false
}

// ──────────────────────────────────────────────
// Hourly Loop
// ──────────────────────────────────────────────

func hourlyLoop(done chan struct{}) {
	time.Sleep(5 * time.Second)

	// Startup message (only if user is not actively using agent)
	if !agent.IsUserActive() {
		welcome := "🟢 <b>Scorp Agent Active</b>\n\n" +
			"I'm ready to help. Send any message to get started.\n" +
			"Use the menu below for quick access:"
		telegram.SendMessage(welcome, telegram.MainMenuKeyboard())
		sendHourlyReport()
	} else {
		log.Println("Startup report deferred: user is active")
	}

	// Async update check (non-blocking)
	go func() {
		msg, err := updater.CheckAndNotify()
		if err != nil {
			log.Printf("[updater] check failed: %v", err)
			return
		}
		if msg != "" {
			telegram.SendMessage(msg, nil)
		}
	}()

	// Clock-aligned: wait until the next full hour
	for {
		now := time.Now()
		nextHour := now.Truncate(time.Hour).Add(time.Hour)
		wait := time.Until(nextHour)
		log.Printf("Next hourly report at %s (in %s)", nextHour.Format("15:04"), wait.Round(time.Second))

		select {
		case <-done:
			return
		case <-time.After(wait):
			// Defer report if user is actively chatting/using agent
			if agent.IsUserActive() {
				log.Println("Hourly report deferred: user is active")
				continue
			}
			sendHourlyReport()
		}
	}
}

func sendHourlyReport() {
	log.Println("Generating hourly report...")

	// Collect all data concurrently
	var sys collectors.SystemData
	var docker collectors.DockerData
	var coolify collectors.CoolifyData
	var sec collectors.SecuritySummary
	var stor collectors.StorageData
	var net collectors.NetworkData

	var wg sync.WaitGroup

	wg.Add(1)
	go func() { defer wg.Done(); sys = collectors.CollectSystem() }()

	wg.Add(1)
	go func() { defer wg.Done(); docker = collectors.CollectDocker() }()

	wg.Add(1)
	go func() { defer wg.Done(); coolify = collectors.CollectCoolify() }()

	wg.Add(1)
	go func() { defer wg.Done(); sec = collectors.CollectSecurity() }()

	wg.Add(1)
	go func() { defer wg.Done(); stor = collectors.CollectStorage() }()

	wg.Add(1)
	go func() { defer wg.Done(); net = collectors.CollectNetwork() }()

	wg.Wait()

	msg := scheduler.FormatHourlyReport(sys, docker, coolify, sec, stor, net)
	ok := telegram.SendMessage(msg, telegram.BackButtonKeyboard())
	if ok {
		log.Println("Hourly report sent successfully")
	} else {
		log.Println("Failed to send hourly report")
	}
}

// ──────────────────────────────────────────────
// Resource Alert Loop
// ──────────────────────────────────────────────

func resourceAlertLoop(done chan struct{}) {
	time.Sleep(15 * time.Second)

	for {
		select {
		case <-done:
			return
		default:
		}

		// Collect concurrently
		var sys collectors.SystemData
		var docker collectors.DockerData
		var stor collectors.StorageData
		var net collectors.NetworkData

		var wg sync.WaitGroup
		wg.Add(4)
		go func() { defer wg.Done(); sys = collectors.CollectSystem() }()
		go func() { defer wg.Done(); docker = collectors.CollectDocker() }()
		go func() { defer wg.Done(); stor = collectors.CollectStorage() }()
		go func() { defer wg.Done(); net = collectors.CollectNetwork() }()
		wg.Wait()

		var allAlerts []string
		allAlerts = append(allAlerts, scheduler.CheckSystemAlerts(sys)...)
		allAlerts = append(allAlerts, scheduler.CheckDockerAlerts(docker)...)
		allAlerts = append(allAlerts, scheduler.CheckStorageAlerts(stor)...)
		allAlerts = append(allAlerts, scheduler.CheckNetworkAlerts(net)...)

		for _, vnc := range collectors.CheckVNCConnections() {
			allAlerts = append(allAlerts, fmt.Sprintf("🖥️ <b>VNC CONNECTION</b>\n🕐 Time: %s", vnc.Time))
		}

		for _, alert := range allAlerts {
			log.Println("Resource alert fired")
			agent.SendMessageSmart(alert, nil)
			time.Sleep(500 * time.Millisecond)
		}

		timer := time.NewTimer(30 * time.Second)
		select {
		case <-done:
			timer.Stop()
			return
		case <-timer.C:
		}
	}
}

// ──────────────────────────────────────────────
// Security Event Loop
// ──────────────────────────────────────────────

func securityEventLoop(done chan struct{}) {
	time.Sleep(3 * time.Second)

	for {
		select {
		case <-done:
			return
		case event := <-collectors.EventChan:
			// Only ssh_login events come through now (real-time)
			log.Printf("Security alert: %s by %s from %s", event.Type, event.User, event.IP)

			basicMsg := formatBasicAlert(event)
			if basicMsg == "" {
				continue
			}

			agent.SendMessageSmart(basicMsg, nil)

			if event.IP != "" {
				go collectors.EnrichWithGeo(event)
			}
		}
	}
}

func formatBasicAlert(event collectors.SecurityEvent) string {
	switch event.Type {
	case "ssh_login":
		return fmt.Sprintf("🔐 <b>SSH LOGIN</b>\n"+
			"👤 User: <code>%s</code>\n"+
			"🌐 IP: <code>%s</code>\n"+
			"🔑 Method: %s\n"+
			"🕐 %s\n"+
			"📍 Looking up location...",
			event.User, event.IP, event.Method, event.Time)
	case "ssh_failed":
		return fmt.Sprintf("❌ <b>SSH FAILED</b>\n"+
			"👤 Tried: <code>%s</code>\n"+
			"🌐 IP: <code>%s</code>\n"+
			"🕐 %s\n"+
			"📍 Looking up location...",
			event.User, event.IP, event.Time)
	}
	return ""
}

// ──────────────────────────────────────────────
// Command Loop
// ──────────────────────────────────────────────

func commandLoop(done chan struct{}) {
	time.Sleep(3 * time.Second)
	log.Println("[telegram] Command loop started")
	telegram.SetupBotCommands()
	os.MkdirAll(config.UploadsDir(), 0755)

	for {
		select {
		case <-done:
			return
		default:
		}

		commands, callbacks, documents, inlineQueries := telegram.PollUpdates()

		for _, cmd := range commands {
			log.Printf("Command: %s from %d", cmd.Text, cmd.ChatID)
			handleAction(cmd.Text, cmd.ChatID, cmd.MsgID, "")
		}

		for _, cb := range callbacks {
			log.Printf("Callback: %s from %d", cb.Data, cb.ChatID)
			telegram.AnswerCallback(cb.CBID, "")
			handleAction(cb.Data, cb.ChatID, cb.MsgID, cb.CBID)
		}

		for _, doc := range documents {
			log.Printf("Upload received: %s (%s, photo=%v)", doc.FileName, telegram.HumanSize(doc.FileSize), doc.IsPhoto)

			// All uploads go through agent (vision for photos, analysis for files)
			go agent.HandleUploadInAgentMode(doc)
		}

		for _, iq := range inlineQueries {
			log.Printf("Inline query: %s from %d", iq.Query, iq.UserID)
			tools.HandleInlineQuery(iq)
		}
	}
}

func handleAction(action string, chatID int64, messageID int64, callbackID string) {
	isCallback := callbackID != ""
	edit := isCallback && messageID != 0

	switch {
	case action == "/start" || action == "mn:main":
		text := "🤖 <b>Scorp Agent</b>\n\n" +
			"💬 Type anything to chat with AI (agent-first)\n" +
			"🛠 <code>/agent</code> for shell, file, web access\n" +
			"━━━━━━━━━━━━━━━━━━━"
		kb := telegram.MainMenuKeyboard()
		if edit {
			telegram.EditMessage(chatID, messageID, text, kb)
		} else {
			telegram.SendMessage(text, kb)
		}

	case action == "mn:mon":
		text := "📊 <b>Monitor Server</b>"
		kb := telegram.MonitorMenuKeyboard()
		if edit {
			telegram.EditMessage(chatID, messageID, text, kb)
		} else {
			telegram.SendMessage(text, kb)
		}

	case action == "mn:sys":
		text := "🔧 <b>System &amp; Tools</b>"
		kb := telegram.SystemMenuKeyboard()
		if edit {
			telegram.EditMessage(chatID, messageID, text, kb)
		} else {
			telegram.SendMessage(text, kb)
		}

	case action == "mn:set":
		text := telegram.SettingsMenuText()
		kb := telegram.SettingsMenuKeyboard()
		if edit {
			telegram.EditMessage(chatID, messageID, text, kb)
		} else {
			telegram.SendMessage(text, kb)
		}

	case strings.HasPrefix(action, "set:"):
		sub := action[4:]
		switch sub {
		case "mon":
			config.Cfg.MonitoringEnabled = !config.Cfg.MonitoringEnabled
			if config.Cfg.MonitoringEnabled {
				loop.MonLoopCtl.Start(func(d chan struct{}) { resourceAlertLoop(d) })
				telegram.AnswerCallback(callbackID, "✅ Monitoring ON")
			} else {
				loop.MonLoopCtl.Stop()
				telegram.AnswerCallback(callbackID, "⏹️ Monitoring OFF")
			}
		case "sec":
			config.Cfg.SecurityAlertsEnabled = !config.Cfg.SecurityAlertsEnabled
			if config.Cfg.SecurityAlertsEnabled {
				loop.SecLoopCtl.Start(func(d chan struct{}) { securityEventLoop(d) })
				telegram.AnswerCallback(callbackID, "✅ Security ON")
			} else {
				loop.SecLoopCtl.Stop()
				telegram.AnswerCallback(callbackID, "⏹️ Security OFF")
			}
		case "rep":
			config.Cfg.ScheduledReportsEnabled = !config.Cfg.ScheduledReportsEnabled
			if config.Cfg.ScheduledReportsEnabled {
				loop.RepLoopCtl.Start(func(d chan struct{}) { hourlyLoop(d) })
				telegram.AnswerCallback(callbackID, "✅ Reports ON")
			} else {
				loop.RepLoopCtl.Stop()
				telegram.AnswerCallback(callbackID, "⏹️ Reports OFF")
			}
		default:
			telegram.AnswerCallback(callbackID, "")
		}
		log.Printf("[settings] Toggle %s → mon=%v sec=%v rep=%v",
			sub, config.Cfg.MonitoringEnabled, config.Cfg.SecurityAlertsEnabled, config.Cfg.ScheduledReportsEnabled)
		// Refresh settings page
		telegram.EditMessage(chatID, messageID, telegram.SettingsMenuText(), telegram.SettingsMenuKeyboard())

	case action == "/status" || action == "status":
		sys := collectors.CollectSystem()
		docker := collectors.CollectDocker()
		stor := collectors.CollectStorage()
		text := scheduler.FormatStatusResponse(sys, docker, stor)
		kb := telegram.BackAndRefreshKeyboard("status")
		if edit {
			telegram.EditMessage(chatID, messageID, text, kb)
		} else {
			telegram.SendMessage(text, kb)
		}

	case action == "report":
		if isCallback {
			telegram.AnswerCallback(callbackID, "⏳ Generating full report...")
		}
		sendHourlyReport()

	case action == "/containers" || action == "containers":
		docker := collectors.CollectDocker()
		text := scheduler.SectionDocker(docker)
		kb := telegram.BackAndRefreshKeyboard("containers")
		if edit {
			telegram.EditMessage(chatID, messageID, text, kb)
		} else {
			telegram.SendMessage(text, kb)
		}

	case action == "/coolify" || action == "coolify":
		coolify := collectors.CollectCoolify()
		text := scheduler.SectionCoolify(coolify)
		kb := telegram.BackAndRefreshKeyboard("coolify")
		if edit {
			telegram.EditMessage(chatID, messageID, text, kb)
		} else {
			telegram.SendMessage(text, kb)
		}

	case action == "/security" || action == "security":
		sec := collectors.CollectSecurityWithPeek()
		text := scheduler.SectionSecurity(sec)
		kb := telegram.BackAndRefreshKeyboard("security")
		if edit {
			telegram.EditMessage(chatID, messageID, text, kb)
		} else {
			telegram.SendMessage(text, kb)
		}

	case action == "/storage" || action == "storage":
		stor := collectors.CollectStorage()
		text := scheduler.SectionStorage(stor)
		kb := telegram.BackAndRefreshKeyboard("storage")
		if edit {
			telegram.EditMessage(chatID, messageID, text, kb)
		} else {
			telegram.SendMessage(text, kb)
		}

	case action == "/network" || action == "network":
		net := collectors.CollectNetwork()
		text := scheduler.SectionNetwork(net)
		kb := telegram.BackAndRefreshKeyboard("network")
		if edit {
			telegram.EditMessage(chatID, messageID, text, kb)
		} else {
			telegram.SendMessage(text, kb)
		}

	case action == "/help" || action == "help":
		text := "🤖 <b>Scorp Agent</b>\n\n" +
			"<b>Agent-First Mode:</b>\n" +
			"💬 Type anything → processed by AI with tools\n" +
			"🛠 <code>/agent</code> — enable agent mode explicitly\n" +
			"🛑 <code>/stop</code> — disable agent mode\n" +
			"🧹 <code>/clear</code> — clear chat history\n\n" +
			"<b>Quick Commands:</b>\n"+
			"🏠 <code>/start</code> — interactive menu\n"+
			"🤖 <code>/model</code> — change AI model\n"+
			"📦 <code>/update</code> — check & install updates\n"+
			"📋 <code>/version</code> — show current version\n\n"+
			"<b>Menu:</b> Monitor, System, Models — all via /start"
		kb := telegram.BackButtonKeyboard()
		if edit {
			telegram.EditMessage(chatID, messageID, text, kb)
		} else {
			telegram.SendMessage(text, kb)
		}

	case action == "/agent" || strings.HasPrefix(action, "/agent "):
			cid := fmt.Sprintf("%d", chatID)
			task := strings.TrimPrefix(action, "/agent")
			task = strings.TrimSpace(task)

			agent.EnterAgentMode(cid)

					if task == "" {
						telegram.SendMessage("🤖 <b>Agent Mode Active!</b>\n\n" +
							"Agent can now run commands, read/write files, and access system.\n" +
							"Send any message to start.\n" +
							"Type /stop to exit.", nil)
			} else {
				go func() {
					msgID := telegram.SendMessageGetID("⏳ <i>Thinking...</i>", chatID)
					if msgID == 0 {
						return
					}
					agent.RunAgentLoop(chatID, task, msgID)
				}()
			}

	case action == "/stop":
		cid := fmt.Sprintf("%d", chatID)
		stopped := false
		if agent.ExitAgentMode(cid) {
			telegram.SendMessage("🛑 Agent mode disabled.", nil)
			stopped = true
		}
		if !stopped {
			telegram.SendMessage("ℹ️ No active mode.", nil)
		}

	case action == "/clear":
		cid := fmt.Sprintf("%d", chatID)
		agent.ClearChatSession(cid)
		agent.ExitAgentMode(cid)
		telegram.SendMessage("🧹 Chat history and agent mode reset.", nil)

	case action == "/sessions":
		// List all saved sessions from disk
		entries, err := os.ReadDir(config.HistoryDirPath())
		if err != nil || len(entries) == 0 {
			telegram.SendMessage("📂 <b>No saved sessions.</b>", nil)
			break
		}
		var sb strings.Builder
		sb.WriteString("📂 <b>Saved Sessions:</b>\n\n")
		for _, e := range entries {
			info, err := e.Info()
			if err != nil {
				continue
			}
			name := strings.TrimSuffix(e.Name(), ".json")
			size := info.Size()
			modTime := info.ModTime().Format("2006-01-02 15:04")
			sb.WriteString(fmt.Sprintf("• <code>%s</code> — %s (%s)\n", name, modTime, telegram.HumanSize(size)))
		}
		sb.WriteString("\n💡 <i>Use /forget to delete current session.</i>")
		telegram.SendMessage(sb.String(), nil)

	case action == "/mcp":
		summary := mcp.MCPToolsSummary()
		telegram.SendMessage("🔌 <b>MCP Servers</b>\n\n"+summary, nil)

	case action == "/cron":
		telegram.SendMessage(scheduler.FormatTasksList(), nil)

	case strings.HasPrefix(action, "/cron "):
		subCmd := strings.TrimPrefix(action, "/cron ")
		parts := strings.SplitN(subCmd, " ", 2)
		switch parts[0] {
		case "run":
			if len(parts) < 2 {
				telegram.SendMessage("❓ Usage: <code>/cron run t1</code>", nil)
				break
			}
			task := scheduler.GetTask(strings.TrimSpace(parts[1]))
			if task == nil {
				telegram.SendMessage(fmt.Sprintf("❌ Task %s not found.", parts[1]), nil)
			} else {
				telegram.SendMessage(fmt.Sprintf("🚀 Running task %s...", task.ID), nil)
				go scheduler.RunTask(*task)
			}
		case "del", "delete", "rm":
			if len(parts) < 2 {
				telegram.SendMessage("❓ Usage: <code>/cron del t1</code>", nil)
				break
			}
			if scheduler.RemoveTask(strings.TrimSpace(parts[1])) {
				telegram.SendMessage(fmt.Sprintf("✅ Task %s deleted.", parts[1]), nil)
			} else {
				telegram.SendMessage(fmt.Sprintf("❌ Task %s not found.", parts[1]), nil)
			}
		case "pause":
			if len(parts) < 2 {
				telegram.SendMessage("❓ Usage: <code>/cron pause t1</code>", nil)
				break
			}
			if scheduler.ToggleTask(strings.TrimSpace(parts[1]), false) {
				telegram.SendMessage(fmt.Sprintf("⏸ Task %s paused.", parts[1]), nil)
			} else {
				telegram.SendMessage(fmt.Sprintf("❌ Task %s not found.", parts[1]), nil)
			}
		case "resume":
			if len(parts) < 2 {
				telegram.SendMessage("❓ Usage: <code>/cron resume t1</code>", nil)
				break
			}
			if scheduler.ToggleTask(strings.TrimSpace(parts[1]), true) {
				telegram.SendMessage(fmt.Sprintf("▶️ Task %s resumed.", parts[1]), nil)
			} else {
				telegram.SendMessage(fmt.Sprintf("❌ Task %s not found.", parts[1]), nil)
			}
		default:
			telegram.SendMessage("❓ Usage:\n"+
				"<code>/cron</code> — List all tasks\n"+
				"<code>/cron run t1</code> — Run now\n"+
				"<code>/cron del t1</code> — Delete\n"+
				"<code>/cron pause t1</code> — Pause\n"+
				"<code>/cron resume t1</code> — Resume", nil)
		}

	case action == "/model" || action == "/model list":
		if isCallback && edit {
			telegram.EditMessage(chatID, messageID, wizard.ModelMenuText(), wizard.ModelMenuKeyboard())
		} else {
			telegram.SendMessage(wizard.ModelMenuText(), wizard.ModelMenuKeyboard())
		}

	case action == "/usage":
		telegram.SendMessage(models.FormatUsageStats(), nil)

	case action == "/version" || action == "/update":
		go func() {
			if action == "/version" {
				telegram.SendMessage(fmt.Sprintf("📦 <b>Version</b>\n\n"+
					"scorp %s\nPlatform: %s/%s", updater.Version, runtime.GOOS, runtime.GOARCH), nil)
				return
			}
			// /update
			telegram.SendMessage("⏳ <i>Checking for updates...</i>", nil)
			msg, err := updater.SelfUpdate()
			if err != nil {
				telegram.SendMessage(fmt.Sprintf("❌ Update failed: %v", err), nil)
			} else {
				telegram.SendMessage(msg, nil)
			}
		}()

	case action == "/files" || action == "files":
		text := "📂 <b>File Manager</b>\n\nChoose a directory to browse:"
		kb := telegram.RootsKeyboard()
		if edit {
			telegram.EditMessage(chatID, messageID, text, kb)
		} else {
			telegram.SendMessage(text, kb)
		}

	case strings.HasPrefix(action, "fb:"):
		pid := action[3:]
		path := telegram.GetPath(pid)
		if path != "" {
			text, kb := telegram.DirKeyboard(path)
			if edit {
				telegram.EditMessage(chatID, messageID, text, kb)
			} else {
				telegram.SendMessage(text, kb)
			}
		}

	case strings.HasPrefix(action, "fd:"):
		pid := action[3:]
		path := telegram.GetPath(pid)
		if path != "" {
			text, kb := telegram.FileDetailKeyboard(path)
			if edit {
				telegram.EditMessage(chatID, messageID, text, kb)
			} else {
				telegram.SendMessage(text, kb)
			}
		}

	case strings.HasPrefix(action, "zp:"):
		pid := action[3:]
		path := telegram.GetPath(pid)
		if path != "" {
			text, kb := telegram.FolderZipInfo(path)
			if edit {
				telegram.EditMessage(chatID, messageID, text, kb)
			} else {
				telegram.SendMessage(text, kb)
			}
		}

	case strings.HasPrefix(action, "zc:"):
		pid := action[3:]
		path := telegram.GetPath(pid)
		if path != "" {
			if isCallback {
				telegram.AnswerCallback(callbackID, "📦 Creating ZIP...")
			}
			telegram.SendFolderAsZip(config.Cfg.TelegramChatID, path)
		}

	case strings.HasPrefix(action, "dl:"):
		pid := action[3:]
		path := telegram.GetPath(pid)
		if path != "" {
			if isCallback {
				telegram.AnswerCallback(callbackID, "⬇️ Sending file...")
			}
			if !telegram.SendFile(config.Cfg.TelegramChatID, path) {
				telegram.SendMessage("❌ Failed to send file (may be too large)", nil)
			}
		}

	case action == "upload":
		text := "📤 <b>Upload File</b>\n\nKirim file apa saja dan saya akan menganalisisnya.\n\n" +
			"Foto → analisis gambar\nFile → analisis konten\n(maks 20MB via Telegram Bot API)"
		kb := telegram.BackButtonKeyboard()
		if edit {
			telegram.EditMessage(chatID, messageID, text, kb)
		} else {
			telegram.SendMessage(text, kb)
		}

	default:
		if !isCallback {
			cid := fmt.Sprintf("%d", chatID)

			// Check if there's a pending clarify question
			if tools.HasPendingClarify(cid) {
				tools.ResolveClarify(action, cid, "")
				break
			}

			// Check if user is in model wizard
			if wizard.GetModelWizard(chatID) != nil && !strings.HasPrefix(action, "/") {
				if wizard.HandleModelWizardTextRouter(action, chatID) {
					break
				}
			}
			// Allow /cancel from wizard even if starts with /
			if action == "/cancel" && wizard.GetModelWizard(chatID) != nil {
				wizard.HandleModelWizardTextRouter(action, chatID)
				break
			}

			// All plain text goes through agent loop (agent-first architecture)
			if !strings.HasPrefix(action, "/") {
				go func() {
					msgID := telegram.SendMessageGetID("⏳ <i>Thinking...</i>", chatID)
					if msgID == 0 {
						return
					}
					// Auto-detect skill mentions in message
					msg := action
					if skillCtx := skills.GetPromptForMessage(action); skillCtx != "" {
						msg = action + skillCtx
					}
					agent.RunAgentLoop(chatID, msg, msgID)
				}()
			} else {
				telegram.SendMessage(fmt.Sprintf("❓ Unknown command: %s\nUse /start to open the menu.", action),
					telegram.MainMenuKeyboard())
			}
		} else {
			// Handle confirmation callbacks
			if action == "confirm_yes" {
				if isCallback {
					telegram.AnswerCallback(callbackID, "✅ Confirmed")
				}
				go agent.HandleConfirmation(chatID, true)
			} else if action == "confirm_no" {
				if isCallback {
					telegram.AnswerCallback(callbackID, "❌ Cancelled")
				}
				go agent.HandleConfirmation(chatID, false)
			} else if strings.HasPrefix(action, "clarify:") {
						// Handle clarify callback responses
						if isCallback {
							telegram.AnswerCallback(callbackID, "")
						}
						tools.ResolveClarify(action, fmt.Sprintf("%d", chatID), callbackID)
			} else if strings.HasPrefix(action, "mdl:") {
				// Model manager callbacks
				if isCallback {
					telegram.AnswerCallback(callbackID, "")
				}
				text, kb, handled := wizard.HandleModelCallback(action, chatID, messageID)
				if handled && text != "" {
					if edit {
						telegram.EditMessage(chatID, messageID, text, kb)
					} else {
						telegram.SendMessage(text, kb)
					}
				}
			}
		}
	}
}
