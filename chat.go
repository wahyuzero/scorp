package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

// ──────────────────────────────────────────────
// ScorpAgent Chat + Agent Mode
// ──────────────────────────────────────────────

// Session tracking
type chatSession struct {
	agentMode      bool // Agent mode active (tools enabled)
	chatMode       bool // Chat mode active (natural conversation)
	lastUsed       time.Time
	history        []agentMessage     // Conversation history for context
	autoStopCancel context.CancelFunc // Cancel previous auto-stop goroutine
	loopActive     bool               // True while runAgentLoop is executing — prevents async summarization
}

// Sharded session map to reduce lock contention
const numSessionShards = 16

var (
	chatSessionsShards [numSessionShards]map[string]*chatSession
	chatSessionsMu     [numSessionShards]sync.Mutex
)

func init() {
	for i := 0; i < numSessionShards; i++ {
		chatSessionsShards[i] = make(map[string]*chatSession)
	}
}

func sessionShard(chatID string) int {
	// FNV-1a hash for good distribution
	h := uint32(2166136261)
	for i := 0; i < len(chatID); i++ {
		h ^= uint32(chatID[i])
		h *= 16777619
	}
	return int(h % numSessionShards)
}

func getSessionMap(chatID string) (map[string]*chatSession, *sync.Mutex) {
	shard := sessionShard(chatID)
	return chatSessionsShards[shard], &chatSessionsMu[shard]
}

func getSession(chatID string) *chatSession {
	m, mu := getSessionMap(chatID)
	mu.Lock()
	defer mu.Unlock()
	return m[chatID]
}

func getOrCreateSession(chatID string) *chatSession {
	m, mu := getSessionMap(chatID)
	mu.Lock()
	defer mu.Unlock()
	if sess, ok := m[chatID]; ok {
		return sess
	}
	sess := &chatSession{}
	m[chatID] = sess
	return sess
}

func setSession(chatID string, sess *chatSession) {
	m, mu := getSessionMap(chatID)
	mu.Lock()
	defer mu.Unlock()
	m[chatID] = sess
}

func deleteSession(chatID string) {
	m, mu := getSessionMap(chatID)
	mu.Lock()
	defer mu.Unlock()
	delete(m, chatID)
}

// Pre-compiled regex for extracting [REMEMBER:key:value] tags
var rememberRe = regexp.MustCompile(`\[REMEMBER:([^:]+):([^\]]+)\]`)

// Pre-compiled regex for inline markdown conversion
var (
	boldRe   = regexp.MustCompile(`\*\*(.+?)\*\*`)
	italicRe = regexp.MustCompile(`(?:^|[^*])\*([^*]+?)\*(?:[^*]|$)`)
	codeRe   = regexp.MustCompile("`([^`]+)`")
)

// ──────────────────────────────────────────────
// Chat & Agent Mode management
// ──────────────────────────────────────────────

func isAgentMode(chatID string) bool {
	sess := getSession(chatID)
	if sess != nil {
		return sess.agentMode
	}
	return false
}

func isChatMode(chatID string) bool {
	sess := getSession(chatID)
	if sess != nil {
		return sess.chatMode
	}
	return false
}

func enterChatMode(chatID string) {
	sess := getOrCreateSession(chatID)
	// Cancel previous auto-stop goroutine to prevent leaks
	if sess.autoStopCancel != nil {
		sess.autoStopCancel()
	}
	sess.chatMode = true
	sess.lastUsed = time.Now()
	sess.history = nil
	ctx, cancel := context.WithCancel(context.Background())
	sess.autoStopCancel = cancel
	setSession(chatID, sess)

	go chatAutoStop(ctx, chatID)
}

func exitChatMode(chatID string) bool {
	sess := getSession(chatID)
	if sess != nil && sess.chatMode {
		sess.chatMode = false
		setSession(chatID, sess)
		go flushPendingMessages()
		return true
	}
	return false
}

// chatAutoStop auto-stops chat mode after 30 min idle (cancellable via context)
func chatAutoStop(ctx context.Context, chatID string) {
	const idleTimeout = 30 * time.Minute
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return // Cancelled — another enter*Mode was called
		case <-ticker.C:
		}
		sess := getSession(chatID)
		if sess == nil || !sess.chatMode {
			return
		}
		idle := time.Since(sess.lastUsed)

		if idle >= idleTimeout {
			exitChatMode(chatID)
			log.Printf("[chat] Auto-stopped for %s after %s idle", chatID, idle.Round(time.Second))
			sendMessage("⏰ Chat mode otomatis dimatikan karena 30 menit tidak aktif.", nil)
			return
		}
	}
}

func enterAgentMode(chatID string) {
	sess := getOrCreateSession(chatID)
	// Cancel previous auto-stop goroutine to prevent leaks
	if sess.autoStopCancel != nil {
		sess.autoStopCancel()
	}
	sess.agentMode = true
	sess.lastUsed = time.Now()
	ctx, cancel := context.WithCancel(context.Background())
	sess.autoStopCancel = cancel
	setSession(chatID, sess)

	go agentAutoStop(ctx, chatID)
}

// agentAutoStop monitors agent idle time and auto-stops after 30 min (cancellable)
func agentAutoStop(ctx context.Context, chatID string) {
	const idleTimeout = 30 * time.Minute
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return // Cancelled — another enter*Mode was called
		case <-ticker.C:
		}
		sess := getSession(chatID)
		if sess == nil || !sess.agentMode {
			return
		}
		idle := time.Since(sess.lastUsed)

		if idle >= idleTimeout {
			exitAgentMode(chatID)
			log.Printf("[agent] Auto-stopped for %s after %s idle", chatID, idle.Round(time.Second))
			sendMessage("⏰ Agent mode otomatis dimatikan karena 30 menit tidak aktif.", nil)
			return
		}
	}
}

func exitAgentMode(chatID string) bool {
	sess := getSession(chatID)
	if sess != nil && sess.agentMode {
		sess.agentMode = false
		setSession(chatID, sess)
		// Flush pending messages
		go flushPendingMessages()
		return true
	}
	return false
}

// ── Pending Message Queue (suppressed during agent mode) ──

var (
	pendingMessages   []string
	pendingMessagesMu sync.Mutex
)

// sendMessageSmart sends immediately if no active mode, otherwise queues
func sendMessageSmart(text string, keyboard map[string]interface{}) {
	if isAnyModeActive() {
		pendingMessagesMu.Lock()
		pendingMessages = append(pendingMessages, text)
		pendingMessagesMu.Unlock()
		log.Printf("[telegram] Message queued (mode active): %s", truncateStr(text, 60))
		return
	}
	sendMessage(text, keyboard)
}

// isAnyModeActive checks if any chat session has agent or chat mode active
func isAnyModeActive() bool {
	for i := 0; i < numSessionShards; i++ {
		chatSessionsMu[i].Lock()
		for _, sess := range chatSessionsShards[i] {
			if sess.agentMode || sess.chatMode {
				chatSessionsMu[i].Unlock()
				return true
			}
		}
		chatSessionsMu[i].Unlock()
	}
	return false
}

// flushPendingMessages sends all queued messages
func flushPendingMessages() {
	time.Sleep(1 * time.Second) // Small delay after agent stops
	pendingMessagesMu.Lock()
	msgs := make([]string, len(pendingMessages))
	copy(msgs, pendingMessages)
	pendingMessages = nil
	pendingMessagesMu.Unlock()

	if len(msgs) == 0 {
		return
	}
	log.Printf("[telegram] Flushing %d pending messages", len(msgs))
	for _, msg := range msgs {
		sendMessage(msg, nil)
		time.Sleep(300 * time.Millisecond)
	}
}

func touchSession(chatID string) {
	sess := getOrCreateSession(chatID)
	sess.lastUsed = time.Now()
	setSession(chatID, sess)
}

// ── Conversation History (with disk persistence) ──

const maxHistoryMessages = 50  // Trim when exceeds this
const keepHistoryMessages = 30 // Keep this many after trim
// historyDir is resolved dynamically via historyDirPath()

func historyFilePath(chatID string) string {
	return fmt.Sprintf("%s/%s.json", historyDirPath(), chatID)
}

func loadHistoryFromDisk(chatID string) []agentMessage {
	data, err := os.ReadFile(historyFilePath(chatID))
	if err != nil {
		return nil
	}
	var msgs []agentMessage
	if err := json.Unmarshal(data, &msgs); err != nil {
		log.Printf("[memory] Failed to parse history for %s: %v", chatID, err)
		return nil
	}
	log.Printf("[memory] Loaded %d messages from disk for %s", len(msgs), chatID)
	return msgs
}

func saveHistoryToDisk(chatID string, msgs []agentMessage) {
	os.MkdirAll(historyDirPath(), 0755)
	data, err := json.Marshal(msgs)
	if err != nil {
		log.Printf("[memory] Failed to save history for %s: %v", chatID, err)
		return
	}
	os.WriteFile(historyFilePath(chatID), data, 0644)
}

func getSessionHistory(chatID string) []agentMessage {
	sess := getSession(chatID)
	if sess != nil && len(sess.history) > 0 {
		result := make([]agentMessage, len(sess.history))
		copy(result, sess.history)
		return result
	}
	// Try loading from disk if not in memory
	msgs := loadHistoryFromDisk(chatID)
	if len(msgs) > 0 {
		sess := getOrCreateSession(chatID)
		sess.history = msgs
		setSession(chatID, sess)
		result := make([]agentMessage, len(msgs))
		copy(result, msgs)
		return result
	}
	return nil
}

func appendSessionHistory(chatID string, msgs ...agentMessage) {
	sess := getOrCreateSession(chatID)
	if sess.history == nil {
		sess.history = loadHistoryFromDisk(chatID)
	}
	sess.history = append(sess.history, msgs...)

	// Index messages for session search (async, non-blocking)
	for _, msg := range msgs {
		if content, ok := msg.Content.(string); ok && content != "" {
			go indexMessage(chatID, msg.Role, content)
		}
	}

	// Auto-summarize if too long — but NOT during an active agent loop (race condition)
	if len(sess.history) > maxHistoryMessages && !sess.loopActive {
		go summarizeHistory(chatID)
	}
	// Debounced persist — schedule save after 3 seconds
	scheduleHistorySave(chatID)
	setSession(chatID, sess)
}

// historyDirty tracks sessions needing a disk save
var (
	historyDirty   = make(map[string]bool)
	historyDirtyMu sync.Mutex
	historySaveChan = make(chan string, 64) // channel for async disk persistence
)

// init starts the async history writer
func init() {
	go historyWriterLoop()
}

// historyWriterLoop runs as a single goroutine to serialize disk writes
func historyWriterLoop() {
	for chatID := range historySaveChan {
		sess := getSession(chatID)
		if sess == nil || len(sess.history) == 0 {
			continue
		}
		snapshot := make([]agentMessage, len(sess.history))
		copy(snapshot, sess.history)
		saveHistoryToDisk(chatID, snapshot)
	}
}

// scheduleHistorySave marks a session as dirty and triggers a delayed save via channel
func scheduleHistorySave(chatID string) {
	historyDirtyMu.Lock()
	alreadyScheduled := historyDirty[chatID]
	historyDirty[chatID] = true
	historyDirtyMu.Unlock()

	if !alreadyScheduled {
		go func() {
			time.Sleep(3 * time.Second)
			historyDirtyMu.Lock()
			delete(historyDirty, chatID)
			historyDirtyMu.Unlock()

			// Send to async writer channel (non-blocking with buffer)
			select {
			case historySaveChan <- chatID:
			default:
				// Channel full, log and drop (should rarely happen with 64 buffer)
				log.Printf("[memory] Warning: history save channel full, dropping save for %s", chatID)
			}
		}()
	}
}

// setLoopActive marks the session's agent loop as active/inactive.
// When active, async summarization is suppressed to prevent race conditions.
func setLoopActive(chatID string, active bool) {
	sess := getOrCreateSession(chatID)
	sess.loopActive = active
	setSession(chatID, sess)
}

// summarizeHistory summarizes old conversation messages using the configured model
func summarizeHistory(chatID string) {
	sess := getSession(chatID)
	if sess == nil || len(sess.history) <= keepHistoryMessages {
		return
	}

	// Split: old messages to summarize, recent to keep
	cutPoint := len(sess.history) - keepHistoryMessages
	oldMessages := make([]agentMessage, cutPoint)
	copy(oldMessages, sess.history[:cutPoint])

	log.Printf("[memory] Summarizing %d old messages for %s", len(oldMessages), chatID)

	// Build summary request
	var sb strings.Builder
	sb.WriteString("Summarize this conversation concisely. Capture key facts, decisions, results, and context. ")
	sb.WriteString("Write in the same language as the conversation. Keep under 500 words.\n\n")
	for _, msg := range oldMessages {
		content := ""
		switch v := msg.Content.(type) {
		case string:
			content = v
		default:
			content = "[non-text content]"
		}
		if len(content) > 300 {
			content = content[:300] + "..."
		}
		sb.WriteString(fmt.Sprintf("[%s]: %s\n", msg.Role, content))
	}

	summary, err := scorpChatMultiTurn([]agentMessage{
		{Role: "system", Content: "You are a conversation summarizer. Output ONLY the summary, no preamble."},
		{Role: "user", Content: sb.String()},
	})
	if err != nil {
		log.Printf("[memory] Summarize failed: %v, falling back to simple trim", err)
		// Fallback: simple trim
		sess := getSession(chatID)
		if sess != nil && len(sess.history) > maxHistoryMessages {
			sess.history = sess.history[len(sess.history)-keepHistoryMessages:]
			// Queue async save
			select {
			case historySaveChan <- chatID:
			default:
				log.Printf("[memory] Warning: history save channel full, dropping save for %s", chatID)
			}
			setSession(chatID, sess)
		}
		return
	}

	// Replace old messages with summary
	sess = getSession(chatID)
	if sess != nil {
		summaryMsg := agentMessage{
			Role:    "system",
			Content: fmt.Sprintf("[Previous conversation summary]\n%s", summary),
		}
		// Keep summary + recent messages
		recentStart := len(sess.history) - keepHistoryMessages
		if recentStart < 0 {
			recentStart = 0
		}
		newHistory := make([]agentMessage, 0, keepHistoryMessages+1)
		newHistory = append(newHistory, summaryMsg)
		newHistory = append(newHistory, sess.history[recentStart:]...)
		sess.history = newHistory
		// Queue async save
		select {
		case historySaveChan <- chatID:
		default:
			log.Printf("[memory] Warning: history save channel full, dropping save for %s", chatID)
		}
		log.Printf("[memory] Summarized %d→%d messages for %s", cutPoint, len(sess.history), chatID)
		setSession(chatID, sess)
	}
}

func clearChatSession(chatID string) {
	sess := getSession(chatID)
	if sess != nil {
		sess.history = nil
		sess.agentMode = false
		setSession(chatID, sess)
	}
	// Remove from disk too
	os.Remove(historyFilePath(chatID))
	log.Printf("[memory] Cleared history for %s", chatID)
}

// ──────────────────────────────────────────────
// Chat API (non-streaming, for /ask)
// ──────────────────────────────────────────────

func scorpChat(chatID, userMessage string) (string, error) {
	touchSession(chatID)
	log.Printf("[scorp] Chat request from %s: %s", chatID, truncateStr(userMessage, 50))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	messages := []chatMessage{
		{Role: "system", Content: "You are Scorp, a helpful AI assistant. Respond in the same language as the user."},
		{Role: "user", Content: userMessage},
	}

	reply, _, err := callModelWithFallback(ctx, "chat", messages)
	if err != nil {
		return "", err
	}

	log.Printf("[scorp] Response received (%d chars)", len(reply))
	return reply, nil
}

// scorpChatMultiTurn sends a full message history using the model router
func scorpChatMultiTurn(messages []agentMessage) (string, error) {
	log.Printf("[scorp] Multi-turn request (%d messages)", len(messages))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Convert to chatMessage format
	chatMsgs := make([]chatMessage, len(messages))
	for i, m := range messages {
		switch c := m.Content.(type) {
		case string:
			chatMsgs[i] = chatMessage{Role: m.Role, Content: c}
		default:
			jsonBytes, _ := json.Marshal(c)
			chatMsgs[i] = chatMessage{Role: m.Role, Content: string(jsonBytes)}
		}
	}

	reply, _, err := callModelWithFallback(ctx, "chat", chatMsgs)
	if err != nil {
		return "", err
	}

	log.Printf("[scorp] Multi-turn response (%d chars)", len(reply))
	return reply, nil
}

// ──────────────────────────────────────────────
// Streaming chat (for agent mode)
// Sends initial message, then edits it as chunks arrive
// ──────────────────────────────────────────────

func scorpStreamChat(chatID int64, userMessage string) {
	cid := fmt.Sprintf("%d", chatID)
	touchSession(cid)
	log.Printf("[scorp] Chat from %s: %s", cid, truncateStr(userMessage, 50))

	// Send "Thinking..." placeholder
	msgID := sendMessageGetID("⏳ <i>Thinking...</i>", chatID)
	if msgID == 0 {
		sendMessage("❌ Failed to send message", nil)
		return
	}

	// Save user message to session history
	appendSessionHistory(cid, agentMessage{Role: "user", Content: userMessage})

	// Build messages with shared memory context
	history := getSessionHistory(cid)
	messages := make([]agentMessage, 0, len(history)+2)

	// System prompt with memory instructions
	sysPrompt := `Kamu adalah Scorp, AI assistant yang cerdas dan ramah.
Jawab dalam bahasa yang sama dengan pengguna.
Gunakan formatting markdown: **bold**, *italic*, ` + "`code`" + `, tabel, dll.

PENTING - Memori jangka panjang:
- Jika pengguna meminta kamu mengingat sesuatu (nama, preferensi, fakta penting), SELALU simpan dengan menambahkan tag di AKHIR jawabanmu:
  [REMEMBER:key:value]
  Contoh: [REMEMBER:nama_user:Wahyu]
  Contoh: [REMEMBER:server_ip:10.0.0.1]
- Tag ini akan otomatis disimpan dan tidak ditampilkan ke user.
- Gunakan fakta dari memori di bawah untuk menjawab pertanyaan.`

	// Inject shared memory if available
	memSummary := getSharedMemorySummary()
	if memSummary != "" {
		sysPrompt += "\n\n[Fakta yang kamu ingat]\n" + memSummary
	}

	messages = append(messages, agentMessage{Role: "system", Content: sysPrompt})
	messages = append(messages, history...)

	// Convert to chatMessage format for model router
	chatMsgs := make([]chatMessage, len(messages))
	for i, m := range messages {
		content, _ := m.Content.(string)
		chatMsgs[i] = chatMessage{Role: m.Role, Content: content}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Use model router with fallback
	reply, usedModel, err := callModelWithFallback(ctx, "chat", chatMsgs)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "invalid_grant") || strings.Contains(errMsg, "unauthorized") || strings.Contains(errMsg, "refresh") {
			editMessageByID(chatID, msgID,
				"🔑 <b>Token expired!</b>\n\nCek API key di <code>/model</code>", nil)
		} else {
			editMessageByID(chatID, msgID, fmt.Sprintf("❌ %s", errMsg), nil)
		}
		return
	}
	_ = usedModel // could log this if needed

	// Extract and save any [REMEMBER:key:value] tags from reply
	cleanReply := extractAndSaveMemory(reply)

	// Save assistant response (clean version)
	appendSessionHistory(cid, agentMessage{Role: "assistant", Content: cleanReply})

	sendScorpReply(chatID, msgID, cleanReply)
	log.Printf("[scorp] Chat response sent (%d chars)", len(cleanReply))
}

// extractAndSaveMemory finds [REMEMBER:key:value] tags, saves to memory, and returns clean text
func extractAndSaveMemory(reply string) string {
	matches := rememberRe.FindAllStringSubmatch(reply, -1)

	for _, match := range matches {
		if len(match) == 3 {
			key := strings.TrimSpace(match[1])
			value := strings.TrimSpace(match[2])
			saveToMemory(key, value)
			log.Printf("[memory] Auto-saved from chat: %s = %s", key, value)
		}
	}

	// Remove tags from display
	clean := rememberRe.ReplaceAllString(reply, "")
	clean = strings.TrimRight(clean, "\n ")
	return clean
}

// saveToMemory delegates to unified in-process cache (tools_memory.go)
func saveToMemory(key, value string) {
	setMemory(key, value)
}

// getSharedMemorySummary delegates to unified in-process cache (tools_memory.go)
func getSharedMemorySummary() string {
	return getMemorySummary()
}

// sendScorpReply converts markdown to Telegram HTML, splits long messages, and sends
// First chunk edits msgID, remaining chunks are sent as new messages
func sendScorpReply(chatID int64, msgID int64, reply string) {
	display := markdownToTelegramHTML(reply)
	header := "🤖 <b>Scorp:</b>\n\n"

	if len(display) <= 3800 {
		editMessageByID(chatID, msgID, header+display, nil)
	} else {
		chunks := splitMessage(display, 3800)
		editMessageByID(chatID, msgID, header+chunks[0], nil)
		for k := 1; k < len(chunks); k++ {
			time.Sleep(300 * time.Millisecond)
			sendMessage(chunks[k], nil)
		}
	}
}

// ──────────────────────────────────────────────
// Session management
// ──────────────────────────────────────────────

func cleanupChatSessions() {
	cutoff := time.Now().Add(-30 * time.Minute)
	for i := 0; i < numSessionShards; i++ {
		chatSessionsMu[i].Lock()
		for id, sess := range chatSessionsShards[i] {
			if sess.lastUsed.Before(cutoff) {
				delete(chatSessionsShards[i], id)
			}
		}
		chatSessionsMu[i].Unlock()
	}
}

func truncateStr(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

// markdownToTelegramHTML converts common Markdown to Telegram HTML
func markdownToTelegramHTML(md string) string {
	lines := strings.Split(md, "\n")
	var result []string
	inCodeBlock := false
	i := 0

	for i < len(lines) {
		line := lines[i]

		// Handle fenced code blocks
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			if inCodeBlock {
				result = append(result, "</pre>")
				inCodeBlock = false
			} else {
				result = append(result, "<pre>")
				inCodeBlock = true
			}
			i++
			continue
		}

		if inCodeBlock {
			result = append(result, line)
			i++
			continue
		}

		// Detect markdown tables (lines starting and ending with |)
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "|") && strings.HasSuffix(trimmed, "|") {
			tableLines := collectTableLines(lines, i)
			if len(tableLines) >= 2 {
				result = append(result, convertTableToList(tableLines)...)
				i += len(tableLines)
				continue
			}
		}

		// Headers: ### Title → bold
		if strings.HasPrefix(line, "### ") {
			title := convertInlineMarkdown(strings.TrimPrefix(line, "### "))
			result = append(result, fmt.Sprintf("\n<b>%s</b>", title))
			i++
			continue
		}
		if strings.HasPrefix(line, "## ") {
			title := convertInlineMarkdown(strings.TrimPrefix(line, "## "))
			result = append(result, fmt.Sprintf("\n<b>%s</b>", title))
			i++
			continue
		}
		if strings.HasPrefix(line, "# ") {
			title := convertInlineMarkdown(strings.TrimPrefix(line, "# "))
			result = append(result, fmt.Sprintf("\n<b>%s</b>", title))
			i++
			continue
		}

		// Horizontal rule
		if trimmed == "---" || trimmed == "***" {
			result = append(result, "━━━━━━━━━━━━━━━")
			i++
			continue
		}

		// Convert inline markdown
		result = append(result, convertInlineMarkdown(line))
		i++
	}

	if inCodeBlock {
		result = append(result, "</pre>")
	}

	return strings.Join(result, "\n")
}

// collectTableLines collects consecutive markdown table lines starting at index
func collectTableLines(lines []string, start int) []string {
	var table []string
	for j := start; j < len(lines); j++ {
		trimmed := strings.TrimSpace(lines[j])
		if strings.HasPrefix(trimmed, "|") && strings.HasSuffix(trimmed, "|") {
			table = append(table, trimmed)
		} else {
			break
		}
	}
	return table
}

// isSeparatorRow checks if a table row is a separator like |---|---|
func isSeparatorRow(row string) bool {
	cells := parseTableRow(row)
	for _, cell := range cells {
		cleaned := strings.ReplaceAll(strings.TrimSpace(cell), "-", "")
		cleaned = strings.ReplaceAll(cleaned, ":", "")
		if cleaned != "" {
			return false
		}
	}
	return true
}

// parseTableRow splits a markdown table row into cells
func parseTableRow(row string) []string {
	// Remove leading/trailing |
	row = strings.TrimSpace(row)
	if strings.HasPrefix(row, "|") {
		row = row[1:]
	}
	if strings.HasSuffix(row, "|") {
		row = row[:len(row)-1]
	}
	parts := strings.Split(row, "|")
	var cells []string
	for _, p := range parts {
		cells = append(cells, strings.TrimSpace(p))
	}
	return cells
}

// convertTableToList converts markdown table lines to a Telegram-friendly list
func convertTableToList(tableLines []string) []string {
	if len(tableLines) < 2 {
		return tableLines
	}

	// Parse header
	headers := parseTableRow(tableLines[0])

	// Find data rows (skip separator)
	var dataRows [][]string
	for k := 1; k < len(tableLines); k++ {
		if isSeparatorRow(tableLines[k]) {
			continue
		}
		dataRows = append(dataRows, parseTableRow(tableLines[k]))
	}

	if len(dataRows) == 0 {
		return nil
	}

	var result []string

	// For 2-column tables: compact "Header: Value" format
	if len(headers) == 2 {
		for _, row := range dataRows {
			val0 := convertInlineMarkdown(safeIndex(row, 0))
			val1 := convertInlineMarkdown(safeIndex(row, 1))
			result = append(result, fmt.Sprintf("▸ <b>%s</b> → %s", val0, val1))
		}
	} else {
		// Multi-column: each row as a card
		for _, row := range dataRows {
			var parts []string
			for j, header := range headers {
				val := convertInlineMarkdown(safeIndex(row, j))
				if val != "" {
					parts = append(parts, fmt.Sprintf("<b>%s:</b> %s", header, val))
				}
			}
			result = append(result, "▫️ "+strings.Join(parts, " · "))
		}
	}

	return result
}

func safeIndex(slice []string, i int) string {
	if i < len(slice) {
		return slice[i]
	}
	return ""
}

// convertInlineMarkdown converts inline markdown: **bold**, *italic*, `code`
func convertInlineMarkdown(line string) string {
	// Bold: **text** → <b>text</b>
	line = boldRe.ReplaceAllString(line, "<b>$1</b>")

	// Italic: *text* → <i>text</i> (but not ** which was already handled)
	for {
		loc := italicRe.FindStringIndex(line)
		if loc == nil {
			break
		}
		match := italicRe.FindStringSubmatch(line)
		if len(match) < 2 {
			break
		}
		old := "*" + match[1] + "*"
		line = strings.Replace(line, old, "<i>"+match[1]+"</i>", 1)
	}

	// Inline code: `text` → <code>text</code>
	line = codeRe.ReplaceAllString(line, "<code>$1</code>")

	return line
}

func cleanupSessionsLoop(done chan struct{}) {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			cleanupChatSessions()
			cleanupAgentSessions()
		}
	}
}

func isUserActive() bool {
	for i := 0; i < numSessionShards; i++ {
		chatSessionsMu[i].Lock()
		for _, sess := range chatSessionsShards[i] {
			// Check if user is in agent mode (always defer)
			if sess.agentMode {
				chatSessionsMu[i].Unlock()
				return true
			}
			// Check if user has been chatting recently
			if time.Since(sess.lastUsed) < 10*time.Minute {
				chatSessionsMu[i].Unlock()
				return true
			}
		}
		chatSessionsMu[i].Unlock()
	}
	return false
}
