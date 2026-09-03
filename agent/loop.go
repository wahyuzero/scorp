package agent

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"scorp-agent/config"
	"scorp-agent/internal/helpers"
	"scorp-agent/models"
	"scorp-agent/rag"
	"scorp-agent/registry"
	"scorp-agent/tools"
)

// maxIterations returns the max agent iterations, configurable via SCORP_MAX_ITERATIONS env (default 25)
func maxIterations() int {
	const defaultMax = 25
	if v := os.Getenv("SCORP_MAX_ITERATIONS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return defaultMax
}

type AgentMessage struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
}

// ──────────────────────────────────────────────
// ReAct Agent Execution Loop
// ──────────────────────────────────────────────

// RunAgentLoop executes the multi-turn agent loop for a chat ID.
func RunAgentLoop(chatID int64, userMessage string, msgID int64) {
	RunAgentSessionLoop(fmt.Sprintf("%d", chatID), chatID, userMessage, msgID)
}

// RunAgentSessionLoop executes the multi-turn agent loop for a named or custom session.
func RunAgentSessionLoop(sessionID string, chatID int64, userMessage string, msgID int64) {
	chatIDStr := sessionID

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load history
	history := getSessionHistory(chatIDStr)

	// Add system prompt if first message
	if len(history) == 0 {
		history = append(history, AgentMessage{
			Role:    "system",
			Content: getAgentSystemPrompt(),
		})
	}

	// Add user message
	history = append(history, AgentMessage{
		Role:    "user",
		Content: userMessage,
	})
	appendSessionHistory(chatIDStr, AgentMessage{Role: "user", Content: userMessage})

	// Mark session as active
	setLoopActive(chatIDStr, true)
	defer setLoopActive(chatIDStr, false)

	// In Termux environments: hold wake lock to prevent CPU sleep during multi-step runs
	if tools.IsTermux() {
		tools.AcquireTermuxWakeLock()
		defer func() {
			tools.ReleaseTermuxWakeLock()
			tools.SendTermuxNotification("Scorp Agent", "Task execution finished")
		}()
	}

	// ── Auto-RAG: Search indexed context for user query ──
	if rag.VecIndex != nil && len(rag.VecIndex.Chunks) > 0 {
		queryFP := rag.ComputeSimhash(userMessage)
		ragResults := rag.VecIndex.HybridSearch(queryFP, userMessage, 3, 0.7)
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
			if len(history) > 0 && history[0].Role == "system" {
				if sysStr, ok := history[0].Content.(string); ok {
					history[0].Content = sysStr + ragContext.String()
				}
			}
		}
	}

	// ── Compact history if getting long ──
	history = maybeCompactHistory(chatIDStr, history)

	// Send initial thinking message
	start := time.Now()
	if msgID == 0 {
		msgID = tools.SendMessageGetID("🧠 <b>Agent</b>\n\n⏳ <i>berpikir...</i>", chatID)
	}

	// ── Real-time typing indicator ──
	tools.SendChatAction(chatID, "typing")
	typingTicker := time.NewTicker(4 * time.Second)
	go func() {
		for {
			select {
			case <-ctx.Done():
				typingTicker.Stop()
				return
			case <-typingTicker.C:
				tools.SendChatAction(chatID, "typing")
			}
		}
	}()

	var thinkingLines []string
	toolCount := 0
	lastThinkingUpdate := time.Now()
	noToolRetries := 0
	recentToolSignatures := make(map[string]int)

	expectedSteps := countStepsInMessage(userMessage)

	for iter := 0; iter < maxIterations(); iter++ {
		// Convert history to ChatMessage format
		chatMsgs := make([]models.ChatMessage, len(history))
		for i, m := range history {
			switch c := m.Content.(type) {
			case string:
				chatMsgs[i] = models.ChatMessage{Role: m.Role, Content: c}
			default:
				chatMsgs[i] = models.ChatMessage{Role: m.Role, Content: fmt.Sprintf("%v", c)}
			}
		}

		// Call model with tools and automatic fallback
		reply, toolCalls, modelUsed, err := models.CallModelWithToolsAndFallback(ctx, "agent", chatMsgs)
		if err != nil {
			errMsg := fmt.Sprintf("❌ Error calling model: %v", err)
			tools.EditMessageByID(chatID, msgID, errMsg, nil)
			return
		}

		if iter == 0 && modelUsed != "" {
			log.Printf("[agent] Using model: %s", modelUsed)
		}

		// Check if done (no tool calls)
		if len(toolCalls) == 0 {
			cleanReply := cleanToolCallTags(reply)

			shouldRetry := false
			retryReason := ""

			if noToolRetries < 5 {
				if mentionsBrowserTask(userMessage) && !screenshotWasTaken(history) {
					shouldRetry = true
					retryReason = "⚠️ INCOMPLETE TASK: The user asked for a browser task. You MUST take a screenshot (browser action=screenshot) before completing."
				} else if looksLikeContinuation(cleanReply) {
					shouldRetry = true
					retryReason = "⚠️ CONTINUATION DETECTED: You stated what you intend to do, but did NOT call any tools. You MUST execute the actions by calling tools NOW."
				} else if toolCount > 0 && expectedSteps >= 2 && toolCount < expectedSteps && !looksLikeContinuation(cleanReply) && !hasCompletionIndicators(cleanReply) {
					shouldRetry = true
					retryReason = fmt.Sprintf("⚠️ PARTIAL COMPLETION: The user asked for %d distinct steps, but you only executed %d tool(s).", expectedSteps, toolCount)
				}
			}

			if shouldRetry {
				noToolRetries++
				log.Printf("[agent] Model responded without tool calls (retry %d/5): %s", noToolRetries, retryReason)
				thinkingLines = append(thinkingLines, fmt.Sprintf("⚠️ Retry %d/5: %s", noToolRetries, strings.Split(retryReason, ":")[0]))
				tools.EditMessageByID(chatID, msgID, buildThinkingMessage(thinkingLines, time.Since(start), false), nil)
				history = append(history, AgentMessage{Role: "assistant", Content: reply})
				history = append(history, AgentMessage{Role: "user", Content: retryReason})
				appendSessionHistory(chatIDStr, AgentMessage{Role: "assistant", Content: reply})
				appendSessionHistory(chatIDStr, AgentMessage{Role: "user", Content: retryReason})
				continue
			}

			// Final answer
			history = append(history, AgentMessage{Role: "assistant", Content: reply})
			appendSessionHistory(chatIDStr, AgentMessage{Role: "assistant", Content: reply})

			// Send final answer
			sendScorpReply(chatID, msgID, cleanReply)

			// Trigger self-review in background
			maybeRunSelfReview(chatID, chatIDStr)
			return
		}

		// Reset retry counter on successful tool call
		noToolRetries = 0

		// Execute tools
		history = append(history, AgentMessage{Role: "assistant", Content: reply})
		appendSessionHistory(chatIDStr, AgentMessage{Role: "assistant", Content: reply})

		for _, tc := range toolCalls {
			// Real-time Steering Queue check (PicoClaw Parity)
			if steerMsg, hasSteer := PopSteeringMessage(chatIDStr); hasSteer {
				log.Printf("[steering] User redirected execution mid-run: %s", steerMsg)
				thinkingLines = append(thinkingLines, fmt.Sprintf("⚡ [INTERRUPT] Steered: %s", helpers.TruncateStr(steerMsg, 50)))
				steerTurnMsg := fmt.Sprintf("⚡ [USER INTERRUPT]: %s", steerMsg)
				history = append(history, AgentMessage{Role: "user", Content: steerTurnMsg})
				appendSessionHistory(chatIDStr, AgentMessage{Role: "user", Content: steerTurnMsg})
				break
			}

			toolCount++
			desc := toolDescription(tc)
			log.Printf("[agent] Executing tool: %s", desc)
			thinkingLines = append(thinkingLines, desc)

			// Check for dangerous commands needing confirmation (bypassed in YOLO mode)
			if tc.Name == "shell" && config.GetAutonomyLevel() != config.AutonomyYOLO && IsDangerousCommand(helpers.GetStringArg(tc.Args, "command", "")) {
				cmd := helpers.GetStringArg(tc.Args, "command", "")
				if shouldUpdateThinking(toolCount, lastThinkingUpdate) {
					tools.EditMessageByID(chatID, msgID, buildThinkingMessage(thinkingLines, time.Since(start), false), nil)
					lastThinkingUpdate = time.Now()
				}

				thinkingLines = append(thinkingLines, "⚠️ Awaiting confirmation...")
				tools.EditMessageByID(chatID, msgID, buildThinkingMessage(thinkingLines, time.Since(start), false), nil)
				lastThinkingUpdate = time.Now()

				promptMsgID := tools.SendMessageGetIDWithKeyboard(
					fmt.Sprintf("⚠️ <b>Dangerous Command</b>\n\n<pre>%s</pre>\n\nAllow execution?", helpers.EscapeHTML(cmd)),
					chatID, confirmKeyboard())

				StorePendingConfirmation(chatIDStr, "shell", cmd, history, promptMsgID)
				return
			}

			// Check for repeated identical actions (loop prevention)
			tcSig := toolCallSignature(tc)
			if tcSig != "" && recentToolSignatures[tcSig] >= 2 {
				warnMsg := fmt.Sprintf("⚠️ STOP: You already executed '%s' %d times with the same arguments. Try a DIFFERENT approach.", desc, recentToolSignatures[tcSig])
				history = append(history, AgentMessage{Role: "user", Content: warnMsg})
				thinkingLines = append(thinkingLines, fmt.Sprintf("  ⚠️ Repeat #%d", recentToolSignatures[tcSig]))
			}
			recentToolSignatures[tcSig]++

			// Update thinking message
			if shouldUpdateThinking(toolCount, lastThinkingUpdate) {
				tools.EditMessageByID(chatID, msgID, buildThinkingMessage(thinkingLines, time.Since(start), false), nil)
				lastThinkingUpdate = time.Now()
			}

			// Execute tool via registry
			result, ok := ExecuteTool(tc, chatID)
			if !ok {
				log.Printf("[agent] Tool %s returned error: %s", tc.Name, helpers.TruncateStr(result, 200))
			}

			// Show preview in thinking stream
			preview := result
			if len(preview) > 60 {
				preview = preview[:57] + "..."
			}
			preview = strings.ReplaceAll(preview, "\n", "\n  • ")
			thinkingLines = append(thinkingLines, fmt.Sprintf("     ↳ %s", preview))

			// Add tool result to history
			toolResult := fmt.Sprintf("[Tool Result: %s]\n%s", tc.Name, result)
			history = append(history, AgentMessage{Role: "user", Content: toolResult})
			appendSessionHistory(chatIDStr, AgentMessage{Role: "user", Content: toolResult})
		}

		if toolCount > 0 {
			tools.EditMessageByID(chatID, msgID, buildThinkingMessage(thinkingLines, time.Since(start), false), nil)
			lastThinkingUpdate = time.Now()
		}

		// Tick dynamic tool discovery TTLs (PicoClaw Parity)
		registry.TickToolTTL()
	}

	// Max iterations reached
	tools.EditMessageByID(chatID, msgID, fmt.Sprintf("⚠️ Agent reached maximum iterations (%d). Last results have been saved to history.", maxIterations()), nil)
}

func cleanToolCallTags(reply string) string {
	_, clean := models.ParseAllToolCalls(reply, nil)
	return clean
}

// resumeAgentLoop continues the agent loop after user confirms a dangerous command
func resumeAgentLoop(chatID int64, messages []AgentMessage, msgID int64) {
	chatIDStr := fmt.Sprintf("%d", chatID)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if msgID == 0 {
		msgID = tools.SendMessageGetID("🧠 <b>Agent</b>\n\n⏳ <i>melanjutkan...</i>", chatID)
	}

	tools.SendChatAction(chatID, "typing")
	typingTicker := time.NewTicker(4 * time.Second)
	go func() {
		for {
			select {
			case <-ctx.Done():
				typingTicker.Stop()
				return
			case <-typingTicker.C:
				tools.SendChatAction(chatID, "typing")
			}
		}
	}()

	start := time.Now()
	var thinkingLines []string
	toolCount := 0
	lastThinkingUpdate := time.Now()
	noToolRetries := 0
	recentToolSignatures := make(map[string]int)

	setLoopActive(chatIDStr, true)
	defer setLoopActive(chatIDStr, false)

	for iter := 0; iter < maxIterations(); iter++ {
		chatMsgs := make([]models.ChatMessage, len(messages))
		for i, m := range messages {
			switch c := m.Content.(type) {
			case string:
				chatMsgs[i] = models.ChatMessage{Role: m.Role, Content: c}
			default:
				chatMsgs[i] = models.ChatMessage{Role: m.Role, Content: fmt.Sprintf("%v", c)}
			}
		}

		reply, toolCalls, _, err := models.CallModelWithToolsAndFallback(ctx, "agent", chatMsgs)
		if err != nil {
			tools.EditMessageByID(chatID, msgID, fmt.Sprintf("❌ Error calling model: %v", err), nil)
			return
		}

		if len(toolCalls) == 0 {
			cleanReply := cleanToolCallTags(reply)

			shouldRetry := false
			retryReason := ""

			if noToolRetries < 5 {
				if looksLikeContinuation(cleanReply) {
					shouldRetry = true
					retryReason = "⚠️ CONTINUATION DETECTED: You stated what you intend to do, but did NOT call any tools. You MUST execute the actions by calling tools NOW."
				}
			}

			if shouldRetry {
				noToolRetries++
				thinkingLines = append(thinkingLines, fmt.Sprintf("⚠️ Retry %d/5: %s", noToolRetries, strings.Split(retryReason, ":")[0]))
				tools.EditMessageByID(chatID, msgID, buildThinkingMessage(thinkingLines, time.Since(start), false), nil)
				messages = append(messages, AgentMessage{Role: "assistant", Content: reply})
				messages = append(messages, AgentMessage{Role: "user", Content: retryReason})
				appendSessionHistory(chatIDStr, AgentMessage{Role: "assistant", Content: reply})
				appendSessionHistory(chatIDStr, AgentMessage{Role: "user", Content: retryReason})
				continue
			}

			messages = append(messages, AgentMessage{Role: "assistant", Content: reply})
			appendSessionHistory(chatIDStr, AgentMessage{Role: "assistant", Content: reply})
			sendScorpReply(chatID, msgID, cleanReply)
			maybeRunSelfReview(chatID, chatIDStr)
			return
		}

		messages = append(messages, AgentMessage{Role: "assistant", Content: reply})
		appendSessionHistory(chatIDStr, AgentMessage{Role: "assistant", Content: reply})

		for _, tc := range toolCalls {
			// Real-time Steering Queue check (PicoClaw Parity)
			if steerMsg, hasSteer := PopSteeringMessage(chatIDStr); hasSteer {
				log.Printf("[steering] User redirected execution mid-run: %s", steerMsg)
				thinkingLines = append(thinkingLines, fmt.Sprintf("⚡ [INTERRUPT] Steered: %s", helpers.TruncateStr(steerMsg, 50)))
				steerTurnMsg := fmt.Sprintf("⚡ [USER INTERRUPT]: %s", steerMsg)
				messages = append(messages, AgentMessage{Role: "user", Content: steerTurnMsg})
				appendSessionHistory(chatIDStr, AgentMessage{Role: "user", Content: steerTurnMsg})
				break
			}

			toolCount++
			desc := toolDescription(tc)
			thinkingLines = append(thinkingLines, desc)

			if tc.Name == "shell" && config.GetAutonomyLevel() != config.AutonomyYOLO && IsDangerousCommand(helpers.GetStringArg(tc.Args, "command", "")) {
				cmd := helpers.GetStringArg(tc.Args, "command", "")
				thinkingLines = append(thinkingLines, "⚠️ Awaiting confirmation...")
				tools.EditMessageByID(chatID, msgID, buildThinkingMessage(thinkingLines, time.Since(start), false), nil)
				lastThinkingUpdate = time.Now()

				promptMsgID := tools.SendMessageGetIDWithKeyboard(
					fmt.Sprintf("⚠️ <b>Dangerous Command</b>\n\n<pre>%s</pre>\n\nAllow execution?", helpers.EscapeHTML(cmd)),
					chatID, confirmKeyboard())

				StorePendingConfirmation(chatIDStr, "shell", cmd, messages, promptMsgID)
				return
			}

			tcSig := toolCallSignature(tc)
			if tcSig != "" && recentToolSignatures[tcSig] >= 2 {
				warnMsg := fmt.Sprintf("⚠️ STOP: You already executed '%s' %d times with the same arguments. Try a DIFFERENT approach.", desc, recentToolSignatures[tcSig])
				messages = append(messages, AgentMessage{Role: "user", Content: warnMsg})
				thinkingLines = append(thinkingLines, fmt.Sprintf("  ⚠️ Repeat #%d", recentToolSignatures[tcSig]))
			}
			recentToolSignatures[tcSig]++

			if shouldUpdateThinking(toolCount, lastThinkingUpdate) {
				tools.EditMessageByID(chatID, msgID, buildThinkingMessage(thinkingLines, time.Since(start), false), nil)
				lastThinkingUpdate = time.Now()
			}

			result, _ := ExecuteTool(tc, chatID)
			preview := result
			if len(preview) > 60 {
				preview = preview[:57] + "..."
			}
			preview = strings.ReplaceAll(preview, "\n", "\n  • ")
			thinkingLines = append(thinkingLines, fmt.Sprintf("     ↳ %s", preview))

			toolResult := fmt.Sprintf("[Tool Result: %s]\n%s", tc.Name, result)
			messages = append(messages, AgentMessage{Role: "user", Content: toolResult})
			appendSessionHistory(chatIDStr, AgentMessage{Role: "user", Content: toolResult})
		}

		if toolCount > 0 {
			tools.EditMessageByID(chatID, msgID, buildThinkingMessage(thinkingLines, time.Since(start), false), nil)
			lastThinkingUpdate = time.Now()
		}

		// Tick dynamic tool discovery TTLs (PicoClaw Parity)
		registry.TickToolTTL()
	}

	tools.EditMessageByID(chatID, msgID, "⚠️ Agent reached maximum iterations.", nil)
}
