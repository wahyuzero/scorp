package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime)

	// MCP-only mode: skip Telegram, just serve MCP over stdio
	if len(os.Args) > 1 && os.Args[1] == "--mcp-server" {
		log.Println("[mcp] MCP-only mode starting...")
		StartMCPServerMode()
		return
	}

	if err := loadConfig(); err != nil {
		log.Fatalf("Config error: %v", err)
	}

	initTelegram()

	// Load dynamic skills from JSON files
	loadSkills()

	// Init credential vault
	initVault()

	// Init browser monitor (scheduled scraping + change detection)
	initMonitor()

	// Init autonomous agent (Phase 7)
	initAutonomous()

	// Periodic browser session cleanup (every 5 min)
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			cleanupStaleBrowserSessions()
		}
	}()

	log.Println("==================================================")
	log.Println("VPS Monitor (Go) starting...")
	log.Printf("Telegram Chat ID: %s", cfg.TelegramChatID)
	log.Printf("Report interval: %ds", cfg.ReportInterval)
	log.Println("==================================================")

	// Start background samplers
	startCPUSampler()
	startDockerStatsSampler()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	var wg sync.WaitGroup
	done := make(chan struct{})

	// Start webhook server if configured, otherwise use long polling
	if cfg.TelegramWebhookURL != "" {
		if err := startWebhookServer(cfg.TelegramWebhookURL); err != nil {
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

	// 1. Hourly report loop
	wg.Add(1)
	go func() {
		defer wg.Done()
		hourlyLoop(done)
	}()

	// 2. Resource alert loop
	wg.Add(1)
	go func() {
		defer wg.Done()
		resourceAlertLoop(done)
	}()

	// 3. Journal watcher
	wg.Add(1)
	go func() {
		defer wg.Done()
		watchJournal(done)
	}()

	// 4. Security event processor
	wg.Add(1)
	go func() {
		defer wg.Done()
		securityEventLoop(done)
	}()

	// 5. Command/callback poller (handled above in webhook/long-polling branch)

	// 6. Session cleanup loop
	wg.Add(1)
	go func() {
		defer wg.Done()
		cleanupSessionsLoop(done)
	}()

	// 7. Start MCP servers (blocking - tools must be registered before LLM calls)
	log.Println("[main] Starting MCP servers...")
	StartMCPServers()
	log.Println("[main] MCP servers started")

	// 7b. Start test endpoint (localhost only)
	startTestEndpoint()

	// 8. Start MCP server mode (scorp-agent as MCP server)
	wg.Add(1)
	go func() {
		defer wg.Done()
		StartMCPServerMode()
	}()

	// 8. Start scheduler
	loadScheduledTasks()
	wg.Add(1)
	go func() {
		defer wg.Done()
		schedulerLoop(done)
	}()

	// 9. Load model router
	loadModelConfig()
	initModelUsage()
	loadCostConfig()
	loadCostTracker()

	// 11. Hermes model watcher
	wg.Add(1)
	go func() {
		defer wg.Done()
		hermesMonitorLoop(done)
	}()

	// 12. Initialize memory cache
	initMemoryCache()

	// 12b. Initialize session search DB (SQLite + FTS5)
	initSessionDB()

	// 13. Initialize RAG index
	initRAG()
	initVectorRAG()

	// 13b. Start autonomous agent loop (Phase 7)
	wg.Add(1)
	go func() {
		defer wg.Done()
		autonomousLoop(done)
	}()

	// 14. Start pprof server (localhost only, for debugging)
	go func() {
		log.Println("[pprof] Listening on localhost:6060")
		if err := http.ListenAndServe("127.0.0.1:6060", nil); err != nil {
			log.Printf("[pprof] Failed to start: %v", err)
		}
	}()

	// 15. Start Prometheus metrics server
	startMetricsServer()

	// Setup inline query mode
	setupInlineMode()

	// 16. Start uptime monitor
	wg.Add(1)
	go func() {
		defer wg.Done()
		uptimeLoop(done)
	}()

	// Wait for shutdown signal
	<-stop
	log.Println("Shutting down...")
	close(done)

	// Stop metrics server
	stopMetricsServer()

	// Stop webhook server if running
	stopWebhookServer()

	// Stop MCP server mode
	StopMCPServerMode()

	// Stop MCP servers
	StopMCPServers()
	sendMessage("🔴 <b>VPS Monitor Stopped</b>", nil)
	wg.Wait()
	log.Println("Goodbye.")
}

// ──────────────────────────────────────────────
// Hourly Loop
// ──────────────────────────────────────────────

func hourlyLoop(done chan struct{}) {
	time.Sleep(5 * time.Second)

	// Startup message (only if user is not actively using agent)
	if !isUserActive() {
		welcome := "🟢 <b>VPS Monitor Started</b>\n\n" +
			"I'll send hourly status reports and real-time alerts.\n" +
			"Use the menu below to check status anytime:"
		sendMessage(welcome, mainMenuKeyboard())
		sendHourlyReport()
	} else {
		log.Println("Startup report deferred: user is active")
	}

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
			if isUserActive() {
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
	var sys SystemData
	var docker DockerData
	var coolify CoolifyData
	var sec SecuritySummary
	var stor StorageData
	var net NetworkData

	var wg sync.WaitGroup

	wg.Add(1)
	go func() { defer wg.Done(); sys = collectSystem() }()

	wg.Add(1)
	go func() { defer wg.Done(); docker = collectDocker() }()

	wg.Add(1)
	go func() { defer wg.Done(); coolify = collectCoolify() }()

	wg.Add(1)
	go func() { defer wg.Done(); sec = collectSecurity() }()

	wg.Add(1)
	go func() { defer wg.Done(); stor = collectStorage() }()

	wg.Add(1)
	go func() { defer wg.Done(); net = collectNetwork() }()

	wg.Wait()

	msg := formatHourlyReport(sys, docker, coolify, sec, stor, net)
	ok := sendMessage(msg, backButtonKeyboard())
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
		var sys SystemData
		var docker DockerData
		var stor StorageData
		var net NetworkData

		var wg sync.WaitGroup
		wg.Add(4)
		go func() { defer wg.Done(); sys = collectSystem() }()
		go func() { defer wg.Done(); docker = collectDocker() }()
		go func() { defer wg.Done(); stor = collectStorage() }()
		go func() { defer wg.Done(); net = collectNetwork() }()
		wg.Wait()

		var allAlerts []string
		allAlerts = append(allAlerts, checkSystemAlerts(sys)...)
		allAlerts = append(allAlerts, checkDockerAlerts(docker)...)
		allAlerts = append(allAlerts, checkStorageAlerts(stor)...)
		allAlerts = append(allAlerts, checkNetworkAlerts(net)...)

		for _, vnc := range checkVNCConnections() {
			allAlerts = append(allAlerts, fmt.Sprintf("🖥️ <b>VNC CONNECTION</b>\n🕐 Time: %s", vnc.Time))
		}

		for _, alert := range allAlerts {
			log.Println("Resource alert fired")
			sendMessageSmart(alert, nil)
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
		case event := <-eventChan:
			// Only ssh_login events come through now (real-time)
			log.Printf("Security alert: %s by %s from %s", event.Type, event.User, event.IP)

			basicMsg := formatBasicAlert(event)
			if basicMsg == "" {
				continue
			}

			sendMessageSmart(basicMsg, nil)

			if event.IP != "" {
				go enrichWithGeo(event)
			}
		}
	}
}

func formatBasicAlert(event SecurityEvent) string {
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
	setupBotCommands()
	os.MkdirAll(uploadsDir(), 0755)

	for {
		select {
		case <-done:
			return
		default:
		}

		commands, callbacks, documents, inlineQueries := pollUpdates()

		for _, cmd := range commands {
			log.Printf("Command: %s from %d", cmd.Text, cmd.ChatID)
			handleAction(cmd.Text, cmd.ChatID, cmd.MsgID, "")
		}

		for _, cb := range callbacks {
			log.Printf("Callback: %s from %d", cb.Data, cb.ChatID)
			answerCallback(cb.CBID, "")
			handleAction(cb.Data, cb.ChatID, cb.MsgID, cb.CBID)
		}

		for _, doc := range documents {
			log.Printf("Upload received: %s (%s, photo=%v, voice=%v)", doc.FileName, humanSize(doc.FileSize), doc.IsPhoto, doc.IsVoice)
			cid := fmt.Sprintf("%d", doc.ChatID)

			// Voice message → STT → process as text
			if doc.IsVoice {
				go handleVoiceMessage(doc)
				continue
			}

			// If in agent mode, route through agent (with vision for photos)
			if isAgentMode(cid) {
				go handleUploadInAgentMode(doc)
				continue
			}

			// Default behavior: just save the file
			dest := fmt.Sprintf("%s/%s", uploadsDir(), doc.FileName)
			os.MkdirAll(uploadsDir(), 0755)
			ok := downloadTelegramFile(doc.FileID, dest)
			if ok {
				fileType := "📄"
				if doc.IsPhoto {
					fileType = "🖼"
				}
				sendMessage(fmt.Sprintf("✅ <b>File Saved</b>\n%s %s\n📏 %s\n📂 <code>~/%s</code>\n\n💡 <i>Gunakan /agent untuk agent mode — AI bisa analisis gambar!</i>",
					fileType, doc.FileName, humanSize(doc.FileSize), doc.FileName), mainMenuKeyboard())
			} else {
				sendMessage(fmt.Sprintf("❌ Failed to save <b>%s</b>", doc.FileName), mainMenuKeyboard())
			}
		}

		for _, iq := range inlineQueries {
			log.Printf("Inline query: %s from %d", iq.Query, iq.UserID)
			handleInlineQuery(iq)
		}
	}
}

func handleAction(action string, chatID int64, messageID int64, callbackID string) {
	isCallback := callbackID != ""
	edit := isCallback && messageID != 0

	switch {
	case action == "/start" || action == "menu":
		text := "🤖 <b>VPS Monitor</b>\n\nChoose an option from the menu below:"
		kb := mainMenuKeyboard()
		if edit {
			editMessage(chatID, messageID, text, kb)
		} else {
			sendMessage(text, kb)
		}

	case action == "/status" || action == "status":
		sys := collectSystem()
		docker := collectDocker()
		stor := collectStorage()
		text := formatStatusResponse(sys, docker, stor)
		kb := backAndRefreshKeyboard("status")
		if edit {
			editMessage(chatID, messageID, text, kb)
		} else {
			sendMessage(text, kb)
		}

	case action == "/report" || action == "report":
		if isCallback {
			answerCallback(callbackID, "⏳ Generating full report...")
		}
		sendHourlyReport()

	case action == "/containers" || action == "containers":
		docker := collectDocker()
		text := sectionDocker(docker)
		kb := backAndRefreshKeyboard("containers")
		if edit {
			editMessage(chatID, messageID, text, kb)
		} else {
			sendMessage(text, kb)
		}

	case action == "/coolify" || action == "coolify":
		coolify := collectCoolify()
		text := sectionCoolify(coolify)
		kb := backAndRefreshKeyboard("coolify")
		if edit {
			editMessage(chatID, messageID, text, kb)
		} else {
			sendMessage(text, kb)
		}

	case action == "/security" || action == "security":
		sec := collectSecurityWithPeek()
		text := sectionSecurity(sec)
		kb := backAndRefreshKeyboard("security")
		if edit {
			editMessage(chatID, messageID, text, kb)
		} else {
			sendMessage(text, kb)
		}

	case action == "/storage" || action == "storage":
		stor := collectStorage()
		text := sectionStorage(stor)
		kb := backAndRefreshKeyboard("storage")
		if edit {
			editMessage(chatID, messageID, text, kb)
		} else {
			sendMessage(text, kb)
		}

	case action == "/network" || action == "network":
		net := collectNetwork()
		text := sectionNetwork(net)
		kb := backAndRefreshKeyboard("network")
		if edit {
			editMessage(chatID, messageID, text, kb)
		} else {
			sendMessage(text, kb)
		}

	case action == "/hermes" || action == "hermes":
		hd := collectHermes()
		text := formatHermesStatus(hd)
		kb := backAndRefreshKeyboard("hermes")
		if edit {
			editMessage(chatID, messageID, text, kb)
		} else {
			sendMessage(text, kb)
		}

	case action == "/help" || action == "help":
		text := fmt.Sprintf("❓ <b>VPS Monitor Help</b>\n\n"+
			"<b>Interactive Menu:</b>\nTap buttons to navigate.\n\n"+
			"<b>Monitor Commands:</b>\n"+
			"/status — Quick status\n/report — Full report\n"+
			"/containers — Docker\n/coolify — Coolify apps\n"+
			"/security — Security\n/storage — Storage\n/network — Network\n\n"+
			"<b>🤖 AI:</b>\n"+
			"Ketik pesan apapun → chat langsung dengan AI\n"+
			"/agent &lt;task&gt; — Agent mode (file/shell/web)\n"+
			"/stop — Stop chat/agent mode\n"+
			"/clear — Clear chat history\n"+
			"/forget — Hapus sesi sepenuhnya\n"+
			"/sessions — Lihat semua sesi tersimpan\n\n"+
			"<b>🎯 Skills:</b>\n"+
			"/skills — List skills\n"+
			"/skill &lt;name&gt; &lt;task&gt; — Run a skill\n\n"+
			"<b>Auto:</b>\n"+
			"📊 Report: setiap jam\n🔐 SSH login: instant alert\n"+
			"⚡ Resource check: 30s\n🔕 Cooldown: %d min",
			cfg.AlertCooldown/60)
		kb := backButtonKeyboard()
		if edit {
			editMessage(chatID, messageID, text, kb)
		} else {
			sendMessage(text, kb)
		}

	case action == "/agent" || strings.HasPrefix(action, "/agent "):
		cid := fmt.Sprintf("%d", chatID)
		// Block if in chat mode
		if isChatMode(cid) {
			sendMessage("⚠️ Kamu sedang dalam chat mode. Ketik /stop dulu untuk beralih ke agent.", nil)
			break
		}
		task := strings.TrimPrefix(action, "/agent")
		task = strings.TrimSpace(task)

		enterAgentMode(cid)

		if task == "" {
			sendMessage("🤖 <b>Agent Mode aktif!</b>\n\n"+
				"Agent sekarang bisa menjalankan perintah, baca/tulis file, dan akses sistem.\n"+
				"Kirim pesan apapun untuk mulai.\n"+
				"Ketik /stop untuk keluar.", nil)
		} else {
			go func() {
				msgID := sendMessageGetID("⏳ <i>Thinking...</i>", chatID)
				if msgID == 0 {
					return
				}
				runAgentLoop(chatID, task, msgID)
			}()
		}

	case action == "/stop":
		cid := fmt.Sprintf("%d", chatID)
		stopped := false
		if exitAgentMode(cid) {
			sendMessage("🛑 Agent mode dimatikan.", nil)
			stopped = true
		}
		if exitChatMode(cid) {
			sendMessage("🛑 Chat mode dimatikan.", nil)
			stopped = true
		}
		if !stopped {
			sendMessage("ℹ️ Tidak ada mode aktif.", nil)
		}

	case action == "/clear":
		cid := fmt.Sprintf("%d", chatID)
		clearChatSession(cid)
		exitAgentMode(cid)
		exitChatMode(cid)
		sendMessage("🧹 Chat history dan semua mode di-reset.", nil)

	case action == "/forget":
		cid := fmt.Sprintf("%d", chatID)
		clearChatSession(cid)
		exitAgentMode(cid)
		exitChatMode(cid)
		sendMessage("🧹 <b>Sesi dihapus!</b>\n\nHistory agent & chat telah di-reset sepenuhnya.\nGunakan /agent untuk memulai sesi baru.", nil)

	case action == "/voice":
		enabled := toggleVoiceReply(chatID)
		if enabled {
			sendMessage("🎤 <b>Voice replies ON</b>\n\nBot akan membalas dengan voice message. Kirim /voice lagi untuk matikan.", nil)
		} else {
			sendMessage("🔇 <b>Voice replies OFF</b>\n\nBot akan membalas dengan text.", nil)
		}

	case action == "/sessions":
		// List all saved sessions from disk
		entries, err := os.ReadDir(historyDirPath())
		if err != nil || len(entries) == 0 {
			sendMessage("📂 <b>Tidak ada sesi tersimpan.</b>", nil)
			break
		}
		var sb strings.Builder
		sb.WriteString("📂 <b>Sesi Tersimpan:</b>\n\n")
		for _, e := range entries {
			info, err := e.Info()
			if err != nil {
				continue
			}
			name := strings.TrimSuffix(e.Name(), ".json")
			size := info.Size()
			modTime := info.ModTime().Format("2006-01-02 15:04")
			sb.WriteString(fmt.Sprintf("• <code>%s</code> — %s (%s)\n", name, modTime, humanSize(size)))
		}
		sb.WriteString("\n💡 <i>Gunakan /forget untuk hapus sesi saat ini.</i>")
		sendMessage(sb.String(), nil)

	case action == "/mcp":
		summary := MCPToolsSummary()
		sendMessage("🔌 <b>MCP Servers</b>\n\n"+summary+"\nUse <code>/mcp reload</code> to restart servers.", nil)

	case action == "/mcp reload":
		sendMessage("🔄 Reloading MCP servers...", nil)
		go func() {
			ReloadMCPServers()
			sendMessage("✅ MCP servers reloaded.\n\n"+MCPToolsSummary(), nil)
		}()

	case action == "/cron":
		sendMessage(formatScheduledTasksList(), nil)

	case strings.HasPrefix(action, "/cron "):
		subCmd := strings.TrimPrefix(action, "/cron ")
		parts := strings.SplitN(subCmd, " ", 2)
		switch parts[0] {
		case "run":
			if len(parts) < 2 {
				sendMessage("❓ Usage: <code>/cron run t1</code>", nil)
				break
			}
			task := getScheduledTask(strings.TrimSpace(parts[1]))
			if task == nil {
				sendMessage(fmt.Sprintf("❌ Task %s not found.", parts[1]), nil)
			} else {
				sendMessage(fmt.Sprintf("🚀 Running task %s...", task.ID), nil)
				go runScheduledTask(*task)
			}
		case "del", "delete", "rm":
			if len(parts) < 2 {
				sendMessage("❓ Usage: <code>/cron del t1</code>", nil)
				break
			}
			if removeScheduledTask(strings.TrimSpace(parts[1])) {
				sendMessage(fmt.Sprintf("✅ Task %s deleted.", parts[1]), nil)
			} else {
				sendMessage(fmt.Sprintf("❌ Task %s not found.", parts[1]), nil)
			}
		case "pause":
			if len(parts) < 2 {
				sendMessage("❓ Usage: <code>/cron pause t1</code>", nil)
				break
			}
			if toggleScheduledTask(strings.TrimSpace(parts[1]), false) {
				sendMessage(fmt.Sprintf("⏸ Task %s paused.", parts[1]), nil)
			} else {
				sendMessage(fmt.Sprintf("❌ Task %s not found.", parts[1]), nil)
			}
		case "resume":
			if len(parts) < 2 {
				sendMessage("❓ Usage: <code>/cron resume t1</code>", nil)
				break
			}
			if toggleScheduledTask(strings.TrimSpace(parts[1]), true) {
				sendMessage(fmt.Sprintf("▶️ Task %s resumed.", parts[1]), nil)
			} else {
				sendMessage(fmt.Sprintf("❌ Task %s not found.", parts[1]), nil)
			}
		default:
			sendMessage("❓ Usage:\n"+
				"<code>/cron</code> — List all tasks\n"+
				"<code>/cron run t1</code> — Run now\n"+
				"<code>/cron del t1</code> — Delete\n"+
				"<code>/cron pause t1</code> — Pause\n"+
				"<code>/cron resume t1</code> — Resume", nil)
		}

	case action == "/skills":
		sendMessage(listSkills(), nil)

	case action == "/model" || action == "/model list":
		sendMessage(formatModelList(), nil)

	case action == "/model check":
		sendMessage("⏳ <i>Checking all API keys...</i>", nil)
		go func() {
			result := formatModelListWithHealth()
			sendMessage(result, nil)
		}()

	case strings.HasPrefix(action, "/model use "):
		name := strings.TrimPrefix(action, "/model use ")
		sendMessage(switchModel("chat", strings.TrimSpace(name)), nil)

	case strings.HasPrefix(action, "/model agent "):
		name := strings.TrimPrefix(action, "/model agent ")
		sendMessage(switchModel("agent", strings.TrimSpace(name)), nil)

	case action == "/usage":
		sendMessage(formatUsageStats(), nil)

	case strings.HasPrefix(action, "/skill ") || action == "/skill":
		if action == "/skill" {
			sendMessage(listSkills(), nil)
			return
		}
		skillPrompt, task, found := handleSkillCommand(action)
		if !found {
			sendMessage("❌ Skill tidak ditemukan. Ketik /skills untuk melihat daftar.", nil)
			return
		}

		cid := fmt.Sprintf("%d", chatID)
		enterAgentMode(cid)

		if task == "" {
			sendMessage("🤖 <b>Skill mode aktif!</b>\nKirim pesan untuk mulai.", nil)
		} else {
			// Run agent with skill prompt injected
			msgID := sendMessageGetID("⏳ <i>Processing with skill...</i>", chatID)
			if msgID == 0 {
				return
			}
			go func() {
				// Inject skill context into user message
				enhancedTask := fmt.Sprintf("[Skill Context]\n%s\n\n[User Task]\n%s", skillPrompt, task)
				runAgentLoop(chatID, enhancedTask, msgID)
			}()
		}

	case action == "/files" || action == "files":
		text := "📂 <b>File Manager</b>\n\nChoose a directory to browse:"
		kb := rootsKeyboard()
		if edit {
			editMessage(chatID, messageID, text, kb)
		} else {
			sendMessage(text, kb)
		}

	case strings.HasPrefix(action, "fb:"):
		pid := action[3:]
		path := getPath(pid)
		if path != "" {
			text, kb := dirKeyboard(path)
			if edit {
				editMessage(chatID, messageID, text, kb)
			} else {
				sendMessage(text, kb)
			}
		}

	case strings.HasPrefix(action, "fd:"):
		pid := action[3:]
		path := getPath(pid)
		if path != "" {
			text, kb := fileDetailKeyboard(path)
			if edit {
				editMessage(chatID, messageID, text, kb)
			} else {
				sendMessage(text, kb)
			}
		}

	case strings.HasPrefix(action, "zp:"):
		pid := action[3:]
		path := getPath(pid)
		if path != "" {
			text, kb := folderZipInfo(path)
			if edit {
				editMessage(chatID, messageID, text, kb)
			} else {
				sendMessage(text, kb)
			}
		}

	case strings.HasPrefix(action, "zc:"):
		pid := action[3:]
		path := getPath(pid)
		if path != "" {
			if isCallback {
				answerCallback(callbackID, "📦 Creating ZIP...")
			}
			sendFolderAsZip(cfg.TelegramChatID, path)
		}

	case strings.HasPrefix(action, "dl:"):
		pid := action[3:]
		path := getPath(pid)
		if path != "" {
			if isCallback {
				answerCallback(callbackID, "⬇️ Sending file...")
			}
			if !sendFile(cfg.TelegramChatID, path) {
				sendMessage("❌ Failed to send file (may be too large)", nil)
			}
		}

	case action == "upload":
		text := "📤 <b>Upload File</b>\n\nSend me any file and I'll save it to the server.\n\n" +
			"It will be saved to <code>~/uploads/</code>\n(max 20MB via Telegram Bot API)"
		kb := backButtonKeyboard()
		if edit {
			editMessage(chatID, messageID, text, kb)
		} else {
			sendMessage(text, kb)
		}

	default:
		if !isCallback {
			cid := fmt.Sprintf("%d", chatID)

			// Check if there's a pending clarify question
			if hasPendingClarify(cid) {
				resolveClarify(action, cid, "")
				break
			}

			// If in agent mode, forward plain text through agentic loop
			if isAgentMode(cid) && !strings.HasPrefix(action, "/") {
				go func() {
					msgID := sendMessageGetID("⏳ <i>Thinking...</i>", chatID)
					if msgID == 0 {
						return
					}
					// Auto-detect skill mentions in message
					msg := action
					if skillCtx := getSkillPromptForMessage(action); skillCtx != "" {
						msg = action + skillCtx
					}
					runAgentLoop(chatID, msg, msgID)
				}()
			} else if !strings.HasPrefix(action, "/") {
				// Natural text input → auto-enter chat mode and send to AI
				if !isChatMode(cid) {
					enterChatMode(cid)
				}
				go scorpStreamChat(chatID, action)
			} else {
				sendMessage(fmt.Sprintf("❓ Unknown command: %s\nUse /start to open the menu.", action),
					mainMenuKeyboard())
			}
		} else {
			// Handle confirmation callbacks
			if action == "confirm_yes" {
				if isCallback {
					answerCallback(callbackID, "✅ Confirmed")
				}
				go handleConfirmation(chatID, true)
			} else if action == "confirm_no" {
				if isCallback {
					answerCallback(callbackID, "❌ Cancelled")
				}
				go handleConfirmation(chatID, false)
			} else if strings.HasPrefix(action, "clarify:") {
				// Handle clarify callback responses
				if isCallback {
					answerCallback(callbackID, "")
				}
				resolveClarify(action, fmt.Sprintf("%d", chatID), callbackID)
			}
		}
	}
}

// ──────────────────────────────────────────────
// Hermes Model Monitor Loop
// ──────────────────────────────────────────────

func hermesMonitorLoop(done chan struct{}) {
	time.Sleep(10 * time.Second)

	// Initialize baseline snapshot
	initial := collectHermes()
	lastHermesSnapshot = initial.modelSnapshot()
	log.Printf("[hermes] Monitor started. Model: %s | Provider: %s | OMH roles: %d | MCP: %d | Gateway: %s",
		initial.Model, initial.Provider, len(initial.OMHRoles), len(initial.MCPServers), initial.GatewayStatus)

	// Send startup notification
	startupMsg := fmt.Sprintf("🤖 <b>Hermes Monitor Active</b>\n"+
		"Model: <code>%s</code> (%s)\n"+
		"Gateway: %s\n"+
		"OMH roles: <b>%d</b> | MCP: <b>%d</b> | Aliases: <b>%d</b>\n"+
		"Plugins: <code>%s</code>\n"+
		"━━━━━━━━━━━━━━━━━\n"+
		"Auto-switch:\n"+
		"🔥 13:00 WIB → glm-5.1 (peak)\n"+
		"🌙 17:00 WIB → glm-5.2 (off-peak)",
		initial.Model, quotaMultiplier(initial.Model, initial.IsPeak),
		initial.GatewayStatus,
		len(initial.OMHRoles), len(initial.MCPServers), len(initial.Aliases),
		strings.Join(initial.Plugins, ", "))
	sendMessageSmart(startupMsg, nil)

	for {
		select {
		case <-done:
			return
		default:
		}

		data := collectHermes()

		snap := data.modelSnapshot()
		if snap != lastHermesSnapshot {
			oldModel := strings.SplitN(lastHermesSnapshot, "|", 2)[0]
			log.Printf("[hermes] Change detected: %s → %s", lastHermesSnapshot, snap)
			alert := formatModelChangeAlert(data, oldModel)
			sendMessageSmart(alert, nil)
			lastHermesSnapshot = snap
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
