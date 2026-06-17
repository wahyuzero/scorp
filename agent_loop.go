package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sort"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// maxIterations returns the max agent iterations, configurable via SCORP_MAX_ITERATIONS env (default 20)
func maxIterations() int {
	const defaultMax = 20
	if v := os.Getenv("SCORP_MAX_ITERATIONS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return defaultMax
}

type agentMessage struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
}

// ──────────────────────────────────────────────
// MCP Tool execution
// ──────────────────────────────────────────────

// contentPart for OpenAI vision format
type contentPart struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *imageURL `json:"image_url,omitempty"`
}

type imageURL struct {
	URL string `json:"url"`
}

// handleUploadInAgentMode handles file/photo uploads when agent mode is active
func handleUploadInAgentMode(doc TGDocument) {
	chatIDStr := fmt.Sprintf("%d", doc.ChatID)
	touchSession(chatIDStr)

	// Download the file
	fileURL := fmt.Sprintf("https://api.telegram.org/bot%s/getFile?file_id=%s", cfg.TelegramBotToken, doc.FileID)
	resp, err := httpShort.Get(fileURL)
	if err != nil {
		sendMessage(fmt.Sprintf("❌ Error getting file: %v", err), nil)
		return
	}
	defer resp.Body.Close()

	var fileResp struct {
		OK     bool `json:"ok"`
		Result struct {
			FilePath string `json:"file_path"`
		} `json:"result"`
	}
	json.NewDecoder(resp.Body).Decode(&fileResp)

	if !fileResp.OK || fileResp.Result.FilePath == "" {
		sendMessage("❌ Could not get file path", nil)
		return
	}

	downloadURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", cfg.TelegramBotToken, fileResp.Result.FilePath)
	fileResp2, err := httpLong.Get(downloadURL)
	if err != nil {
		sendMessage(fmt.Sprintf("❌ Error downloading: %v", err), nil)
		return
	}
	defer fileResp2.Body.Close()

	fileData, _ := io.ReadAll(fileResp2.Body)

	// Determine if it's an image
	ext := strings.ToLower(filepath.Ext(fileResp.Result.FilePath))
	isImage := ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".webp"

	if isImage {
		// Send as vision message
		b64 := base64Encode(fileData)
		mimeType := "image/jpeg"
		if ext == ".png" {
			mimeType = "image/png"
		} else if ext == ".gif" {
			mimeType = "image/gif"
		} else if ext == ".webp" {
			mimeType = "image/webp"
		}

		parts := []contentPart{
			{Type: "image_url", ImageURL: &imageURL{URL: fmt.Sprintf("data:%s;base64,%s", mimeType, b64)}},
			{Type: "text", Text: "Analyze this image. Describe what you see."},
		}

		if doc.Caption != "" {
			parts[1].Text = doc.Caption
		}

		// Build message with vision content
		msgs := getSessionHistory(chatIDStr)
		msgs = append(msgs, agentMessage{Role: "user", Content: parts})
		appendSessionHistory(chatIDStr, agentMessage{Role: "user", Content: parts})

		msgID := sendMessageGetID("🔍 <i>Analyzing image...</i>", doc.ChatID)

		// Call Scorp API directly (no tool loop needed for vision)
		chatMsgs := make([]chatMessage, len(msgs))
		for i, m := range msgs {
			switch c := m.Content.(type) {
			case string:
				chatMsgs[i] = chatMessage{Role: m.Role, Content: c}
			default:
				jsonBytes, _ := json.Marshal(c)
				chatMsgs[i] = chatMessage{Role: m.Role, Content: string(jsonBytes)}
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		reply, _, err := callModelWithFallback(ctx, "agent", chatMsgs)
		if err != nil {
			editMessageByID(doc.ChatID, msgID, fmt.Sprintf("❌ Error: %v", err), nil)
			return
		}

		appendSessionHistory(chatIDStr, agentMessage{Role: "assistant", Content: reply})
		sendScorpReply(doc.ChatID, msgID, reply)
	} else {
		// Non-image file: save and inform agent
		savePath := fmt.Sprintf("/tmp/scorp_upload_%d_%s", time.Now().Unix(), filepath.Base(fileResp.Result.FilePath))
		os.WriteFile(savePath, fileData, 0644)

		userMsg := fmt.Sprintf("User uploaded a file: %s (%d bytes, saved to %s)", filepath.Base(fileResp.Result.FilePath), len(fileData), savePath)
		if doc.Caption != "" {
			userMsg += "\nCaption: " + doc.Caption
		}

		runAgentLoop(doc.ChatID, userMsg, 0)
	}
}

// base64Encode encodes bytes to base64 string
func base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// toolDescription returns a human-readable description of a tool call
func toolDescription(tc ToolCall) string {
	switch tc.Name {
	case "shell":
		cmd := getStringArg(tc.Args, "command", "")
		if len(cmd) > 80 {
			cmd = cmd[:80] + "..."
		}
		return fmt.Sprintf("🖥 shell: %s", cmd)
	case "read_file":
		return fmt.Sprintf("📖 read: %s", getStringArg(tc.Args, "path", ""))
	case "write_file":
		return fmt.Sprintf("✏️ write: %s", getStringArg(tc.Args, "path", ""))
	case "list_dir":
		return fmt.Sprintf("📂 list: %s", getStringArg(tc.Args, "path", "."))
	case "system_info":
		return fmt.Sprintf("ℹ️ sysinfo: %s", getStringArg(tc.Args, "type", "full"))
	case "send_file":
		return fmt.Sprintf("📤 send: %s", getStringArg(tc.Args, "path", ""))
	case "web_fetch":
		return fmt.Sprintf("🌐 fetch: %s", getStringArg(tc.Args, "url", ""))
	case "web_search":
		return fmt.Sprintf("🔍 search: %s", getStringArg(tc.Args, "query", ""))
	case "memory":
		action := getStringArg(tc.Args, "action", "")
		key := getStringArg(tc.Args, "key", "")
		return fmt.Sprintf("🧠 memory.%s(%s)", action, key)
	case "search_code":
		return fmt.Sprintf("🔎 search_code: %s in %s", getStringArg(tc.Args, "pattern", "?"), getStringArg(tc.Args, "path", "."))
	case "git":
		return fmt.Sprintf("📦 git.%s (%s)", getStringArg(tc.Args, "action", "?"), getStringArg(tc.Args, "repo", "."))
	case "http":
		return fmt.Sprintf("📡 http.%s → %s", getStringArg(tc.Args, "method", "GET"), getStringArg(tc.Args, "url", "?"))
	case "log":
		return fmt.Sprintf("📋 log.%s(%s)", getStringArg(tc.Args, "source", "?"), getStringArg(tc.Args, "target", "?"))
	case "sql":
		query := getStringArg(tc.Args, "query", "?")
		if len(query) > 50 {
			query = query[:47] + "..."
		}
		return fmt.Sprintf("🗄 sql: %s", query)
	case "process":
		return fmt.Sprintf("⚙️ process.%s", getStringArg(tc.Args, "action", "?"))
	case "browser":
		action := getStringArg(tc.Args, "action", "")
		if action == "goto" {
			return fmt.Sprintf("🌐 browser→%s", getStringArg(tc.Args, "url", ""))
		}
		return fmt.Sprintf("🌐 browser.%s", action)
	case "analyze_image":
		return fmt.Sprintf("👁 analyze_image: %s", getStringArg(tc.Args, "path", "?"))
	case "mcp_tool":
		server := getStringArg(tc.Args, "server", "")
		tool := getStringArg(tc.Args, "tool", "")
		return fmt.Sprintf("🔌 mcp: %s.%s", server, tool)
	case "delegate":
		task := getStringArg(tc.Args, "task", "?")
		if len(task) > 60 {
			task = task[:60] + "..."
		}
		return fmt.Sprintf("🤖 delegate: %s", task)
	default:
		return fmt.Sprintf("🔧 %s", tc.Name)
	}
}

// buildThinkingMessage builds the thinking stream display
func buildThinkingMessage(lines []string, elapsed time.Duration, done bool) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("🧠 <b>Agent</b> [%s]\n\n", elapsed.Round(time.Second)))

	for _, line := range lines {
		sb.WriteString(line + "\n")
	}

	if !done {
		sb.WriteString("\n⏳ <i>working...</i>")
	}

	return sb.String()
}

// shouldUpdateThinking returns true if we should update the thinking message.
// Batches updates: every 2 tool calls OR every 2 seconds (whichever comes first).
func shouldUpdateThinking(toolCount int, lastUpdate time.Time) bool {
	if toolCount%2 == 0 {
		return true
	}
	if time.Since(lastUpdate) > 2*time.Second {
		return true
	}
	return false
}

// runAgentLoop executes the multi-turn agent loop.
// It loads conversation history from session, sends to Scorp, parses tool calls,
// executes tools, and loops until final answer. History is persisted across turns.
func runAgentLoop(chatID int64, userMessage string, msgID int64) {
	chatIDStr := fmt.Sprintf("%d", chatID)

	// Create context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load history
	history := getSessionHistory(chatIDStr)

	// Add system prompt if first message
	if len(history) == 0 {
		history = append(history, agentMessage{
			Role:    "system",
			Content: getAgentSystemPrompt(),
		})
	}

	// Add user message
	history = append(history, agentMessage{
		Role:    "user",
		Content: userMessage,
	})
	appendSessionHistory(chatIDStr, agentMessage{Role: "user", Content: userMessage})

	// Mark session as "agent loop active" to prevent async summarization (race condition fix)
	setLoopActive(chatIDStr, true)
	defer setLoopActive(chatIDStr, false)

	// ── Auto-RAG: Search indexed context for user query ──
	if vecIndex != nil && len(vecIndex.Chunks) > 0 {
		queryFP := computeSimhash(userMessage)
		ragResults := vecIndex.hybridSearch(queryFP, userMessage, 3, 0.7)
		if len(ragResults) > 0 {
			var ragContext strings.Builder
			ragContext.WriteString("\n\n### Relevant Context (auto-RAG)\n")
			ragContext.WriteString("The following information was retrieved from indexed knowledge that may be relevant:\n\n")
			for i, r := range ragResults {
				ragContext.WriteString(fmt.Sprintf("**%d.** (score: %.2f)\n```\n", i+1, r.Final))
				preview := r.Chunk.Content
				if len(preview) > 600 {
					preview = preview[:600] + "..."
				}
				ragContext.WriteString(preview)
				ragContext.WriteString("\n```\n")
			}
			// Inject context into system prompt (first message)
			if len(history) > 0 && history[0].Role == "system" {
				if sysStr, ok := history[0].Content.(string); ok {
					history[0].Content = sysStr + ragContext.String()
				}
			}
		}
	}

	// Create thinking display
	if msgID == 0 {
		msgID = sendMessageGetID("🧠 <b>Agent</b>\n\n⏳ <i>thinking...</i>", chatID)
	}

	start := time.Now()
	var thinkingLines []string
	var toolCount int
	var lastThinkingUpdate = time.Now()
	noToolRetries := 0 // consecutive iterations where model returned 0 tool calls
	recentToolSignatures := map[string]int{}

	for iteration := 0; iteration < maxIterations(); iteration++ {
		// Compact history if token budget exceeded
		history = maybeCompactHistory(chatIDStr, history)

		// Convert to chat messages
		chatMsgs := make([]chatMessage, len(history))
		for i, m := range history {
			switch c := m.Content.(type) {
			case string:
				chatMsgs[i] = chatMessage{Role: m.Role, Content: c}
			default:
				jsonBytes, _ := json.Marshal(c)
				chatMsgs[i] = chatMessage{Role: m.Role, Content: string(jsonBytes)}
			}
		}
		// Call Scorp API with full conversation + native function calling
		log.Printf("[agent] Iteration %d: calling model (%d messages)...", iteration, len(chatMsgs))
		reply, nativeToolCalls, modelUsed, err := callModelWithToolsAndFallback(ctx, "agent", chatMsgs)
		if err != nil {
			log.Printf("[agent] Model error: %v", err)
			if strings.Contains(err.Error(), "token") || strings.Contains(err.Error(), "401") {
				sendMessage(
					"🔑 <b>Token expired!</b>\n\nJalankan di server:\n<pre>scorp login</pre>", nil)
				return
			}
			editMessageByID(chatID, msgID, fmt.Sprintf("❌ Error: %v", err), nil)
			return
		}

		_ = modelUsed

		// Parse tool calls: native → XML tags → code-block fallback
		toolCalls, cleanResponse := parseAllToolCalls(reply, nativeToolCalls)
		log.Printf("[agent] Model replied: %d tool calls (native=%d), response len=%d", len(toolCalls), len(nativeToolCalls), len(reply))

		// ── Handle "0 tool calls" responses with retry logic ──
		if len(toolCalls) == 0 && iteration < maxIterations()-1 {
			noToolRetries++
			shouldRetry := false
			reason := ""

			// Case 1: Model said it would continue but didn't call tools
			if looksLikeContinuation(reply) {
				shouldRetry = true
				reason = "continuation detected"
			}

			// Case 2: Tools were called before but task seems incomplete
			// If the agent has used tools and then suddenly stops, it likely hasn't finished.
			// This catches cases where the model says "done" but the task isn't actually complete.
			if !shouldRetry && toolCount > 0 && noToolRetries <= 2 {
				// Check if user explicitly asked for multiple steps
				if userSteps := countStepsInMessage(userMessage); userSteps >= 2 && toolCount < userSteps {
					shouldRetry = true
					reason = fmt.Sprintf("incomplete multi-step task (%d/%d steps)", toolCount, userSteps)
				}
				// Check if user mentioned browser/screenshot keywords but no screenshot was taken
				if !shouldRetry && mentionsBrowserTask(userMessage) && !screenshotWasTaken(history) {
					shouldRetry = true
					reason = "browser task not complete (no screenshot taken)"
				}
			}

			// Cap retries at 3 to prevent infinite loops
			if shouldRetry && noToolRetries > 3 {
				log.Printf("[agent] Max no-tool retries (%d) reached, accepting as final answer", noToolRetries)
				shouldRetry = false
			}

			if shouldRetry {
				log.Printf("[agent] 0 tool calls at iteration %d (%s), retry %d/3 — injecting nudge", iteration, reason, noToolRetries)
				history = append(history, agentMessage{Role: "assistant", Content: reply})
				appendSessionHistory(chatIDStr, agentMessage{Role: "assistant", Content: reply})

				// Escalating urgency
				var reminder string
				if noToolRetries <= 2 {
					reminder = fmt.Sprintf("⚠️ You haven't completed the task yet (%s). You must CALL TOOLS to finish. Do NOT describe what you would do — EXECUTE the tools NOW.", reason)
				} else {
					reminder = fmt.Sprintf("🚨 CRITICAL: The task is NOT complete (%s). This is your LAST CHANCE. Call the necessary tools immediately or explain clearly why you cannot proceed.", reason)
				}
				history = append(history, agentMessage{Role: "user", Content: reminder})
				appendSessionHistory(chatIDStr, agentMessage{Role: "user", Content: reminder})

				thinkingLines = append(thinkingLines, fmt.Sprintf("⚠️ Retry %d/3: %s", noToolRetries, reason))
				editMessageByID(chatID, msgID, buildThinkingMessage(thinkingLines, time.Since(start), false), nil)
				continue
			}
		} else if len(toolCalls) > 0 {
			noToolRetries = 0 // Reset on successful tool call
		}

		if len(toolCalls) == 0 {
			// Final answer — no more tool calls
			history = append(history, agentMessage{Role: "assistant", Content: reply})
			appendSessionHistory(chatIDStr, agentMessage{Role: "assistant", Content: reply})

			cleanReply := extractAndSaveMemory(cleanResponse)
			sendScorpReply(chatID, msgID, cleanReply)
			return
		}

		// Execute tools
		history = append(history, agentMessage{Role: "assistant", Content: reply})
		appendSessionHistory(chatIDStr, agentMessage{Role: "assistant", Content: reply})

		for _, tc := range toolCalls {
			toolCount++
			desc := toolDescription(tc)
			log.Printf("[agent] Executing tool: %s", desc)
			thinkingLines = append(thinkingLines, desc)

			// Check for dangerous commands needing confirmation
			if tc.Name == "shell" && isDangerousCommand(getStringArg(tc.Args, "command", "")) {
				cmd := getStringArg(tc.Args, "command", "")
				if shouldUpdateThinking(toolCount, lastThinkingUpdate) {
					editMessageByID(chatID, msgID, buildThinkingMessage(thinkingLines, time.Since(start), false), nil)
					lastThinkingUpdate = time.Now()
				}

				// Store pending confirmation with full context
				storePendingConfirmation(chatIDStr, "shell", cmd, history)

				thinkingLines = append(thinkingLines, "⚠️ Awaiting confirmation...")
				editMessageByID(chatID, msgID, buildThinkingMessage(thinkingLines, time.Since(start), false), nil)
				lastThinkingUpdate = time.Now()

				sendMessage(
					fmt.Sprintf("⚠️ <b>Dangerous Command</b>\n\n<pre>%s</pre>\n\nAllow execution?", escapeHTML(cmd)),
					confirmKeyboard())
				return
			}

			// Check for repeated identical actions (browser loop detection)
			tcSig := toolCallSignature(tc)
			if tcSig != "" && recentToolSignatures[tcSig] >= 2 {
				warnMsg := fmt.Sprintf("⚠️ STOP: You already executed '%s' %d times with the same arguments. The action is NOT working — trying again will not help. Try a DIFFERENT approach: check the page state, try a different selector, or report what's happening.", desc, recentToolSignatures[tcSig])
				log.Printf("[agent] Repeated action detected: %s (count=%d), injecting warning", tcSig, recentToolSignatures[tcSig])
				history = append(history, agentMessage{Role: "user", Content: warnMsg})
				appendSessionHistory(chatIDStr, agentMessage{Role: "user", Content: warnMsg})
				thinkingLines = append(thinkingLines, fmt.Sprintf("  ⚠️ Repeat #%d blocked", recentToolSignatures[tcSig]))
				// Still execute but with warning so model sees fresh feedback
			}
			recentToolSignatures[tcSig]++

			// Execute tool
			result, _ := executeTool(tc, chatID)
			thinkingLines = append(thinkingLines, fmt.Sprintf("  → %s", truncateStr(result, 60)))

			// Update thinking display (batched)
			if shouldUpdateThinking(toolCount, lastThinkingUpdate) {
				editMessageByID(chatID, msgID, buildThinkingMessage(thinkingLines, time.Since(start), false), nil)
				lastThinkingUpdate = time.Now()
			}

			// Add tool result to history
			toolResult := fmt.Sprintf("[Tool Result: %s]\n%s", tc.Name, result)
			history = append(history, agentMessage{Role: "user", Content: toolResult})
			appendSessionHistory(chatIDStr, agentMessage{Role: "user", Content: toolResult})
		}
		// Final update after all tools in this iteration
		if toolCount > 0 {
			editMessageByID(chatID, msgID, buildThinkingMessage(thinkingLines, time.Since(start), false), nil)
			lastThinkingUpdate = time.Now()
		}
	}

	// Max iterations reached
	editMessageByID(chatID, msgID, fmt.Sprintf("⚠️ Agent reached maximum iterations (%d). Last results have been saved to history.", maxIterations()), nil)
}

// ──────────────────────────────────────────────
// Confirmation System
// ──────────────────────────────────────────────

type pendingConfirmation struct {
	toolName string
	command  string
	messages []agentMessage
	created  time.Time
}

var (
	pendingConfirms = make(map[string]*pendingConfirmation)
	pendingConfirmsMu sync.Mutex
)

func storePendingConfirmation(chatID, toolName, command string, messages []agentMessage) {
	pendingConfirmsMu.Lock()
	defer pendingConfirmsMu.Unlock()
	pendingConfirms[chatID] = &pendingConfirmation{
		toolName: toolName,
		command:  command,
		messages: messages,
		created:  time.Now(),
	}
}

func getPendingConfirmation(chatID string) *pendingConfirmation {
	pendingConfirmsMu.Lock()
	defer pendingConfirmsMu.Unlock()
	pc, ok := pendingConfirms[chatID]
	if !ok {
		return nil
	}
	// Expire after 5 minutes
	if time.Since(pc.created) > 5*time.Minute {
		delete(pendingConfirms, chatID)
		return nil
	}
	return pc
}

func clearPendingConfirmation(chatID string) {
	pendingConfirmsMu.Lock()
	defer pendingConfirmsMu.Unlock()
	delete(pendingConfirms, chatID)
}

func confirmKeyboard() map[string]interface{} {
	return map[string]interface{}{
		"inline_keyboard": [][]map[string]string{
			{
				{"text": "✅ Yes", "callback_data": "/confirm_yes"},
				{"text": "❌ No", "callback_data": "/confirm_no"},
			},
		},
	}
}

// handleConfirmation processes a user's yes/no response to a pending confirmation
func handleConfirmation(chatID int64, confirmed bool) {
	chatIDStr := fmt.Sprintf("%d", chatID)
	pc := getPendingConfirmation(chatIDStr)

	if pc == nil {
		sendMessage("❌ No pending confirmation found. It may have expired.", nil)
		return
	}

	if !confirmed {
		// User rejected
		clearPendingConfirmation(chatIDStr)

		if pc.messages != nil {
			// Resume agent with rejection
			toolResult := fmt.Sprintf("[Tool Result: %s]\nUser REJECTED the command: %s\nPlease suggest an alternative approach.", pc.toolName, pc.command)
			pc.messages = append(pc.messages, agentMessage{Role: "user", Content: toolResult})
			appendSessionHistory(chatIDStr, agentMessage{Role: "user", Content: toolResult})

			sendMessage("❌ Command rejected. Agent will suggest alternatives...", nil)
			resumeAgentLoop(chatID, pc.messages, 0)
		} else {
			sendMessage("❌ Command rejected.", nil)
		}
		return
	}

	// User confirmed — execute the command
	clearPendingConfirmation(chatIDStr)

	result, ok := executeShell(map[string]interface{}{"command": pc.command, "timeout": 60}, chatID)
	status := "✅"
	if !ok {
		status = "❌"
	}

	if pc.messages != nil {
		// Resume agent with result
		toolResult := fmt.Sprintf("[Tool Result: %s]\n%s%s", pc.toolName, status, result)
		pc.messages = append(pc.messages, agentMessage{Role: "user", Content: toolResult})
		appendSessionHistory(chatIDStr, agentMessage{Role: "user", Content: toolResult})

		sendMessage(fmt.Sprintf("%s Command executed. Continuing...", status), nil)
		resumeAgentLoop(chatID, pc.messages, 0)
	} else {
		sendMessage(fmt.Sprintf("%s Result:\n<pre>%s</pre>", status, escapeHTML(truncateStr(result, 2000))), nil)
	}
}

// resumeAgentLoop continues the agent loop after confirmation
func resumeAgentLoop(chatID int64, messages []agentMessage, msgID int64) {
	chatIDStr := fmt.Sprintf("%d", chatID)

	// Create context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if msgID == 0 {
		msgID = sendMessageGetID("🧠 <b>Agent</b>\n\n⏳ <i>continuing...</i>", chatID)
	}

	start := time.Now()
	var thinkingLines []string
	var toolCount int
	var lastThinkingUpdate = time.Now()
	noToolRetries := 0
	recentToolSignatures := map[string]int{}

	// Mark session as "agent loop active" to prevent async summarization
	setLoopActive(chatIDStr, true)
	defer setLoopActive(chatIDStr, false)

	for iteration := 0; iteration < maxIterations(); iteration++ {
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

		reply, nativeToolCalls, _, err := callModelWithToolsAndFallback(ctx, "agent", chatMsgs)
		if err != nil {
			editMessageByID(chatID, msgID, fmt.Sprintf("❌ Error: %v", err), nil)
			return
		}

		toolCalls, cleanResponse := parseAllToolCalls(reply, nativeToolCalls)
		log.Printf("[agent-continue] Model replied: %d tool calls (native=%d)", len(toolCalls), len(nativeToolCalls))

		// ── Handle "0 tool calls" with retry logic (same as runAgentLoop) ──
		if len(toolCalls) == 0 && iteration < maxIterations()-1 {
			noToolRetries++
			shouldRetry := false
			reason := ""

			if looksLikeContinuation(reply) {
				shouldRetry = true
				reason = "continuation detected"
			}
			if shouldRetry && noToolRetries > 3 {
				log.Printf("[agent-continue] Max no-tool retries (%d) reached, accepting as final answer", noToolRetries)
				shouldRetry = false
			}

			if shouldRetry {
				log.Printf("[agent-continue] 0 tool calls at iteration %d (%s), retry %d/3", iteration, reason, noToolRetries)
				messages = append(messages, agentMessage{Role: "assistant", Content: reply})
				appendSessionHistory(chatIDStr, agentMessage{Role: "assistant", Content: reply})

				reminder := fmt.Sprintf("⚠️ You haven't completed the task yet (%s). CALL THE APPROPRIATE TOOL(S) NOW.", reason)
				messages = append(messages, agentMessage{Role: "user", Content: reminder})
				appendSessionHistory(chatIDStr, agentMessage{Role: "user", Content: reminder})

				thinkingLines = append(thinkingLines, fmt.Sprintf("⚠️ Retry %d/3: %s", noToolRetries, reason))
				editMessageByID(chatID, msgID, buildThinkingMessage(thinkingLines, time.Since(start), false), nil)
				continue
			}
		} else if len(toolCalls) > 0 {
			noToolRetries = 0
		}

		if len(toolCalls) == 0 {
			messages = append(messages, agentMessage{Role: "assistant", Content: reply})
			appendSessionHistory(chatIDStr, agentMessage{Role: "assistant", Content: reply})
			cleanReply := extractAndSaveMemory(cleanResponse)
			sendScorpReply(chatID, msgID, cleanReply)
			return
		}

		messages = append(messages, agentMessage{Role: "assistant", Content: reply})
		appendSessionHistory(chatIDStr, agentMessage{Role: "assistant", Content: reply})

		for _, tc := range toolCalls {
			toolCount++
			desc := toolDescription(tc)
			thinkingLines = append(thinkingLines, desc)

			if tc.Name == "shell" && isDangerousCommand(getStringArg(tc.Args, "command", "")) {
				cmd := getStringArg(tc.Args, "command", "")
				storePendingConfirmation(chatIDStr, "shell", cmd, messages)
				thinkingLines = append(thinkingLines, "⚠️ Awaiting confirmation...")
				editMessageByID(chatID, msgID, buildThinkingMessage(thinkingLines, time.Since(start), false), nil)
				lastThinkingUpdate = time.Now()
				sendMessage(
					fmt.Sprintf("⚠️ <b>Dangerous Command</b>\n\n<pre>%s</pre>\n\nAllow execution?", escapeHTML(cmd)),
					confirmKeyboard())
				return
			}

			// Check for repeated identical actions (browser loop detection)
			tcSig := toolCallSignature(tc)
			if tcSig != "" && recentToolSignatures[tcSig] >= 2 {
				warnMsg := fmt.Sprintf("⚠️ STOP: You already executed '%s' %d times with the same arguments. The action is NOT working — trying again will not help. Try a DIFFERENT approach.", desc, recentToolSignatures[tcSig])
				log.Printf("[agent-continue] Repeated action detected: %s (count=%d)", tcSig, recentToolSignatures[tcSig])
				messages = append(messages, agentMessage{Role: "user", Content: warnMsg})
				thinkingLines = append(thinkingLines, fmt.Sprintf("  ⚠️ Repeat #%d", recentToolSignatures[tcSig]))
			}
			recentToolSignatures[tcSig]++

			result, _ := executeTool(tc, chatID)
			thinkingLines = append(thinkingLines, fmt.Sprintf("  → %s", truncateStr(result, 60)))

			// Update thinking display (batched)
			if shouldUpdateThinking(toolCount, lastThinkingUpdate) {
				editMessageByID(chatID, msgID, buildThinkingMessage(thinkingLines, time.Since(start), false), nil)
				lastThinkingUpdate = time.Now()
			}

			toolResult := fmt.Sprintf("[Tool Result: %s]\n%s", tc.Name, result)
			messages = append(messages, agentMessage{Role: "user", Content: toolResult})
			appendSessionHistory(chatIDStr, agentMessage{Role: "user", Content: toolResult})
		}
		// Final update after all tools in this iteration
		if toolCount > 0 {
			editMessageByID(chatID, msgID, buildThinkingMessage(thinkingLines, time.Since(start), false), nil)
			lastThinkingUpdate = time.Now()
		}
	}

	editMessageByID(chatID, msgID, "⚠️ Agent reached maximum iterations.", nil)
}

// looksLikeContinuation detects if the model's response indicates intent to continue
// but didn't actually call tools (e.g., "Let me...", "I'll try...", "Mari coba...")
func looksLikeContinuation(text string) bool {
	lower := strings.ToLower(text)
	patterns := []string{
		"let me ", "i'll ", "i will ", "i'm going to ", "going to ",
		"mar i coba", "mari coba", "saya akan ", "akan coba ", "coba ",
		"i'll try", "let me try", "try to ", "attempt to ",
		"next i", "then i", "now i", "continue to ",
		"proceed to ", "follow up", "followup",
		// Additional patterns — model says it will do something but hasn't yet
		"i need to ", "still need to ", "first, let me",
		"after that", "once that", "now let's",
		"step 1", "step 2", "step 3",
		"first,", "second,", "third,",
		"i haven't", "not yet", "still working",
	}
	for _, p := range patterns {
		if strings.Contains(lower, p) {
			return true
		}
	}
	return false
}

// countStepsInMessage estimates how many distinct steps/actions the user's message asks for.
// Detects numbered lists (1. 2. 3.), "then" chains, and explicit multi-action words.
func countStepsInMessage(msg string) int {
	lower := strings.ToLower(msg)
	count := 0

	// Count numbered list items: "1)", "1.", "2)", "2.", etc.
	for i := 1; i <= 20; i++ {
		if strings.Contains(lower, fmt.Sprintf("%d.", i)) || strings.Contains(lower, fmt.Sprintf("%d)", i)) {
			count++
		}
	}
	if count >= 2 {
		return count
	}

	// Count "then" / "after that" / "next" chains
	count += strings.Count(lower, " then ")
	count += strings.Count(lower, " after that ")
	count += strings.Count(lower, " next,")
	count += strings.Count(lower, "\nthen ")
	count += strings.Count(lower, "→")

	if count >= 2 {
		return count + 1 // "do X then Y" = 2 steps minimum
	}

	return 0
}

// mentionsBrowserTask checks if the user's message references browser actions
// like login, screenshot, scrape, or navigate — used to detect incomplete browser tasks.
func mentionsBrowserTask(msg string) bool {
	lower := strings.ToLower(msg)
	keywords := []string{
		"screenshot", "login", "log in", "sign in", "scrape",
		"ambil gambar", "tangkap layar", "masuk", "buka web",
		"buka halaman", "cek halaman", "cek web", "navigate",
	}
	for _, kw := range keywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}

// screenshotWasTaken checks if a browser.screenshot tool call exists in the history.
func screenshotWasTaken(history []agentMessage) bool {
	for _, msg := range history {
		content, ok := msg.Content.(string)
		if !ok {
			continue
		}
		// Tool results are stored as "[Tool Result: browser.screenshot]"
		if strings.Contains(content, "[Tool Result: browser.screenshot]") ||
			strings.Contains(content, "browser.screenshot") {
			return true
		}
		// Also check assistant messages that contain tool call markup
		if msg.Role == "assistant" && strings.Contains(strings.ToLower(content), "screenshot") {
			// Only count if it's a tool call, not just text mentioning screenshot
			if strings.Contains(content, "<tool") || strings.Contains(content, "```tool") {
				return true
			}
		}
	}
	return false
}

// toolCallSignature returns a string key identifying a tool call by name + args.
// Used for detecting repeated identical actions. Returns "" for non-trackable calls.
func toolCallSignature(tc ToolCall) string {
	if tc.Name == "" {
		return ""
	}
	// Build signature from name + sorted key args
	sig := tc.Name
	keys := make([]string, 0, len(tc.Args))
	for k := range tc.Args {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		sig += "|" + k + "=" + fmt.Sprintf("%v", tc.Args[k])
	}
	return sig
}
