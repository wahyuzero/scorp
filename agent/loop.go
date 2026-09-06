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
	"scorp-agent/session"
	"scorp-agent/skills"
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

// maxTurnTimeout returns the maximum wall-clock duration for a single agent run (default 30m)
func maxTurnTimeout() time.Duration {
	const defaultTimeout = 30 * time.Minute
	if v := os.Getenv("SCORP_MAX_TURN_TIMEOUT"); v != "" {
		if d, err := time.ParseDuration(v); err == nil && d > 0 {
			return d
		}
	}
	return defaultTimeout
}

type AgentMessage struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
}

// ──────────────────────────────────────────────
// ReAct Agent Execution Loop with Task State Machine
// ──────────────────────────────────────────────

// RunAgentLoop executes the multi-turn agent loop for a chat ID.
func RunAgentLoop(chatID int64, userMessage string, msgID int64) {
	RunAgentSessionLoop(fmt.Sprintf("%d", chatID), chatID, userMessage, msgID)
}

// RunAgentSessionLoop executes the multi-turn agent loop for a named or custom session.
func RunAgentSessionLoop(sessionID string, chatID int64, userMessage string, msgID int64) {
	chatIDStr := sessionID

	// Wrap execution with safety wall-clock timeout to prevent runaway deadloops (Incident 2026-09-05)
	ctx, cancel := context.WithTimeout(context.Background(), maxTurnTimeout())
	defer cancel()

	// Load history
	history := getSessionHistory(chatIDStr)

	// Add system prompt if first message
	if len(history) == 0 {
		history = append(history, AgentMessage{
			Role:    "system",
			Content: getAgentSystemPrompt(chatID),
		})
	}

	// Check if user issued an explicit continuation directive (e.g. "continue", "lanjutkan")
	isContinuation := IsContinuationDirective(userMessage)
	activeUserPrompt := userMessage
	if isContinuation {
		activeUserPrompt = fmt.Sprintf("%s\n\n[⚡ SYSTEM DIRECTIVE: CONTINUE] The user wants previous work continued. Scan the conversation history: identify the OUTERMOST original request and finish ALL of its remaining parts NOW by invoking real action tool(s). A completed sub-step is NOT completion — only call 'complete_task' when every part of the original request is done and verified. If the history shows no original task at all (or nothing remains), briefly state that there is nothing to continue and call 'complete_task' with that status. NEVER reference tools that do not exist in your tool list.]", userMessage)
	}

	// Cross-turn Task Ledger: a previous run that hit its budget hands back an
	// unfinished plan — put the contract back in front of the model so ANY new
	// message (not just "continue") resumes precisely where the work stopped.
	if plan := GetTaskPlan(chatIDStr); plan != nil {
		if len(plan.Unfinished()) == 0 {
			ClearTaskPlan(chatIDStr)
		} else {
			activeUserPrompt = fmt.Sprintf("%s\n\n[📋 ACTIVE PLAN — %d item(s) unfinished. Finish ALL of them now via real tools, keep task_plan statuses updated, then complete_task.]\n%s", activeUserPrompt, len(plan.Unfinished()), plan.Render())
		}
	}

	// New-Request Roll-Up: a new request on a session saturated by a previous
	// CLOSED task gets a compact history + hard turn boundary (verify26 turn-19
	// lesson: without it the model re-executed the previous task).
	history = prepareNewTurnHistory(chatIDStr, history, isContinuation)

	// Add user message
	history = append(history, AgentMessage{
		Role:    "user",
		Content: activeUserPrompt,
	})
	appendSessionHistory(chatIDStr, AgentMessage{Role: "user", Content: userMessage})

	// Mark session and chat as active
	setLoopActive(chatIDStr, true)
	defer setLoopActive(chatIDStr, false)
	rawChatIDStr := fmt.Sprintf("%d", chatID)
	if rawChatIDStr != chatIDStr && chatID != 0 {
		setLoopActive(rawChatIDStr, true)
		defer setLoopActive(rawChatIDStr, false)
	}

	// Auto-title session in background if unnamed and on turn 1
	if ShouldAutoTitleSession(sessionID) && len(history) <= 3 {
		go func(oldID string, prompt string, targetChatID int64) {
			time.Sleep(3500 * time.Millisecond)
			newTitle := GenerateContextualSessionTitle(prompt)
			if newTitle != "" && newTitle != oldID {
				if err := RenameSession(oldID, newTitle); err == nil {
					log.Printf("[session] Auto-titled session '%s' -> '%s'", oldID, newTitle)
					if tools.OnSessionAutoTitled != nil {
						tools.OnSessionAutoTitled(oldID, newTitle, fmt.Sprintf("%d", targetChatID))
					}
				} else {
					log.Printf("[session] Failed to rename auto-titled session '%s' -> '%s': %v", oldID, newTitle, err)
				}
			}
		}(sessionID, userMessage, chatID)
	}

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

	// ── Session Search context (SQLite FTS5) ──
	if sessionCtx := getSessionSearchContext(chatIDStr, userMessage); sessionCtx != "" {
		if len(history) > 0 && history[0].Role == "system" {
			if sysStr, ok := history[0].Content.(string); ok {
				history[0].Content = sysStr + sessionCtx
			}
		}
	}

	// Send initial thinking indicator
	if msgID == 0 {
		msgID = tools.SendMessageGetID("🧠 <b>Agent</b>\n\n⏳ <i>memproses...</i>", chatID)
	}

	start := time.Now()
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
	completeTaskGateNudged := false
	autoResumes := 0      // Task Ledger auto-resume budget (see maxAutoResumes)
	lastFullThought := "" // full text of the last non-empty thought-only reply
	recentToolSignatures := make(map[string]int)

	isPureInfo := IsPureInformationalQuery(userMessage)

	for iter := 0; iter < maxIterations(); iter++ {
		// Real-time Steering Queue check at the start of each iteration
		steerMsg, hasSteer := PopSteeringMessage(chatIDStr)
		if !hasSteer && rawChatIDStr != chatIDStr && chatID != 0 {
			steerMsg, hasSteer = PopSteeringMessage(rawChatIDStr)
		}
		if hasSteer {
			log.Printf("[steering] User redirected execution before iteration %d: %s", iter, steerMsg)
			thinkingLines = append(thinkingLines, fmt.Sprintf("⚡ [INTERRUPT] Steered: %s", helpers.TruncateStr(steerMsg, 50)))
			steerTurnMsg := fmt.Sprintf("⚡ [USER INTERRUPT]: %s", steerMsg)
			history = append(history, AgentMessage{Role: "user", Content: steerTurnMsg})
			appendSessionHistory(chatIDStr, AgentMessage{Role: "user", Content: steerTurnMsg})
			tools.EditMessageByID(chatID, msgID, buildThinkingMessage(thinkingLines, time.Since(start), false), nil)
			userMessage = steerMsg
			isPureInfo = IsPureInformationalQuery(userMessage)
			noToolRetries = 0
		}

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

		// ── Check if Model Invoked 'complete_task' ──
		var explicitFinalResult string
		isTaskExplicitlyCompleted := false

		var actionToolCalls []ToolCall
		for _, tc := range toolCalls {
			if tc.Name == "complete_task" {
				isTaskExplicitlyCompleted = true
				if res, ok := tc.Args["result"].(string); ok && res != "" {
					explicitFinalResult = res
				} else if sum, ok := tc.Args["summary"].(string); ok && sum != "" {
					explicitFinalResult = sum
				}
			} else {
				actionToolCalls = append(actionToolCalls, tc)
			}
		}

		// If complete_task was explicitly called, conclude task immediately!
		if isTaskExplicitlyCompleted {
			// Anti-fabrication gate: on an action task, reject complete_task
			// when ZERO tools were executed this turn — the most common shape
			// of a hallucinated "Done" report is concluding from stale history
			// without doing the work. Allow one immediate re-conclusion so
			// genuinely tool-free action tasks are not bricked.
			if toolCount == 0 && !isPureInfo && !completeTaskGateNudged {
				completeTaskGateNudged = true
				noToolRetries = 0
				log.Printf("[agent] complete_task rejected — no tools executed this turn (anti-fabrication gate)")
				gateNudge := "⚠️ COMPLETE_TASK REJECTED: You concluded the task WITHOUT executing any tools this turn, and the user's request is an action task. Fabricating or reusing previous results is FORBIDDEN. Execute the required tools now and verify real outputs; only then call complete_task with evidence-backed results. (If — and only if — the task genuinely requires no system action, call complete_task again with the final answer.)"
				history = append(history, AgentMessage{Role: "assistant", Content: reply})
				history = append(history, AgentMessage{Role: "user", Content: gateNudge})
				continue
			}

			// Plan-completion gate (Task Ledger): an unfinished plan is
			// deterministic proof the request is not fulfilled — reject the
			// premature completion and keep executing, bounded by autoResumes.
			if plan := GetTaskPlan(chatIDStr); plan != nil {
				if un := plan.Unfinished(); len(un) > 0 {
					if autoResumes >= maxAutoResumes {
						planHandback := fmt.Sprintf("\n\n⏸ Auto-resume budget exhausted (%d attempts). Task is NOT fully finished — remaining items:\n%s\nSend \"continue\" to resume.", autoResumes, renderPlanItems(un))
						finalOutput := explicitFinalResult + planHandback
						history = append(history, AgentMessage{Role: "assistant", Content: reply})
						appendSessionHistory(chatIDStr, AgentMessage{Role: "assistant", Content: finalOutput})
						sendScorpReply(chatID, msgID, finalOutput)
						maybeRunSelfReview(chatID, chatIDStr)
						return
					}
					autoResumes++
					noToolRetries = 0
					log.Printf("[agent] complete_task rejected — %d unfinished plan items (Task Ledger gate, resume %d/%d)", len(un), autoResumes, maxAutoResumes)
					gateNudge := fmt.Sprintf("⚠️ COMPLETE_TASK REJECTED — the Task Ledger still has %d unfinished item(s):\n%s\nContinue executing them NOW with real tools (mark each 'in_progress' → 'done' via task_plan as you verify). complete_task is accepted ONLY when every item is done.", len(un), renderPlanItems(un))
					history = append(history, AgentMessage{Role: "assistant", Content: reply})
					history = append(history, AgentMessage{Role: "user", Content: gateNudge})
					thinkingLines = append(thinkingLines, fmt.Sprintf("🛑 complete_task rejected: %d/%d plan items left", len(un), len(plan.Items)))
					continue
				}
				ClearTaskPlan(chatIDStr) // every item verified — contract fulfilled
			}

			finalOutput := explicitFinalResult
			if finalOutput == "" {
				finalOutput = cleanToolCallTags(reply)
			}
			if finalOutput == "" {
				finalOutput = "Task completed successfully."
			}

			history = append(history, AgentMessage{Role: "assistant", Content: reply})
			appendSessionHistory(chatIDStr, AgentMessage{Role: "assistant", Content: finalOutput})

			sendScorpReply(chatID, msgID, finalOutput)
			maybeRunSelfReview(chatID, chatIDStr)
			return
		}

		// ── State Machine Check: Model emitted NO tool calls ──
		if len(actionToolCalls) == 0 {
			cleanReply := cleanToolCallTags(reply)

			shouldRetry := false
			retryReason := ""

			// Detect text-form tool syntax (DSML/XML) — a training prior of
			// some gateway models. Nudge specifically so the retry uses
			// native function calling instead of repeating the mistake.
			dsmlHint := ""
			if strings.Contains(cleanReply, "DSML") || strings.Contains(cleanReply, "<tool_call>") || strings.Contains(cleanReply, "tool_call>") {
				dsmlHint = " Your last reply contained raw tool-call TAGS (DSML/XML). That text is NOT executed. Emit the tool call through native function calling instead."
			}

			// If user asked a simple informational/conceptual question (e.g. "What is Docker?"), allow direct text answer.
			if isPureInfo && toolCount == 0 && iter == 0 {
				shouldRetry = false
			} else if noToolRetries < 4 {
				// Otherwise, this is an action task where model emitted thought text without calling an action tool or complete_task!
				shouldRetry = true
				if isContinuation {
					retryReason = "⚠️ CONTINUE PROTOCOL: Finish the remaining parts of the OUTERMOST original request found in the history — invoke a REAL tool from your tool list NOW, or call 'complete_task' only if truly nothing remains. Do NOT invent tools and do NOT emit plain text."
				} else {
					retryReason = "⚠️ TASK IN-PROGRESS: You outputted intermediate thought without calling any tools. In Agent Mode, you must either call an action tool to continue execution, or call 'complete_task' with your final report if all requested work is finished and verified."
				}
				retryReason += dsmlHint
			}

			if shouldRetry {
				noToolRetries++
				if strings.TrimSpace(cleanReply) != "" {
					lastFullThought = cleanReply
				}
				log.Printf("[agent] Model emitted intermediate thought without tool calls (retry %d/4): %s", noToolRetries, helpers.TruncateStr(cleanReply, 80))

				// Display intermediate thought in thinking stream instead of terminating!
				thoughtPreview := cleanReply
				if len(thoughtPreview) > 80 {
					thoughtPreview = thoughtPreview[:77] + "..."
				}
				thoughtPreview = strings.ReplaceAll(thoughtPreview, "\n", " ")
				thinkingLines = append(thinkingLines, fmt.Sprintf("💭 %s", thoughtPreview))
				tools.EditMessageByID(chatID, msgID, buildThinkingMessage(thinkingLines, time.Since(start), false), nil)

				history = append(history, AgentMessage{Role: "assistant", Content: reply})
				history = append(history, AgentMessage{Role: "user", Content: retryReason})
				continue
			}

			// Autonomous persistence (Task Ledger): never end the turn while
			// the ledger shows unfinished work — auto-resume execution instead
			// so the user never has to type "continue".
			planHandback := ""
			if plan := GetTaskPlan(chatIDStr); plan != nil && len(plan.Unfinished()) > 0 {
				if autoResumes < maxAutoResumes {
					autoResumes++
					noToolRetries = 0
					log.Printf("[agent] auto-resume %d/%d — %d unfinished plan items", autoResumes, maxAutoResumes, len(plan.Unfinished()))
					thinkingLines = append(thinkingLines, fmt.Sprintf("🔄 Auto-resume: %d/%d plan items left", len(plan.Unfinished()), len(plan.Items)))
					resumeMsg := fmt.Sprintf("⚠️ AUTONOMOUS PERSISTENCE: The task is NOT finished — %d plan item(s) remain:\n%s\nDo NOT stop for conversation. Execute the next item now with real tools; keep task_plan statuses updated.", len(plan.Unfinished()), renderPlanItems(plan.Unfinished()))
					history = append(history, AgentMessage{Role: "assistant", Content: reply})
					history = append(history, AgentMessage{Role: "user", Content: resumeMsg})
					continue
				}
				// Budget exhausted: transparent handback with exact remaining work.
				planHandback = fmt.Sprintf("\n\n⏸ Auto-resume budget exhausted (%d attempts). Task is NOT fully finished — remaining items:\n%s\nSend \"continue\" to resume.", autoResumes, renderPlanItems(plan.Unfinished()))
			}

			// Max retries reached or pure informational query: output clean response as final
			history = append(history, AgentMessage{Role: "assistant", Content: reply})
			appendSessionHistory(chatIDStr, AgentMessage{Role: "assistant", Content: cleanReply})

			// Never deliver an empty bubble: if the model exhausted its
			// no-tool retries without producing a report, first try one
			// tool-free synthesis call, then fall back to its last full
			// thought, then to a transparent execution summary.
			if strings.TrimSpace(cleanReply) == "" {
				if toolCount > 0 {
					// The model did real work but stopped narrating. Ask it to
					// summarize WITHOUT the tools schema so it must answer in text.
					log.Printf("[agent] Running tool-free synthesis call for final report (%d tools executed)", toolCount)
					synthMsgs := make([]models.ChatMessage, 0, len(history)+1)
					for _, m := range history {
						if c, ok := m.Content.(string); ok {
							synthMsgs = append(synthMsgs, models.ChatMessage{Role: m.Role, Content: c})
						}
					}
					synthMsgs = append(synthMsgs, models.ChatMessage{Role: "user", Content: "All tools above have ALREADY been executed. Write the final summary (3-8 sentences) describing the REAL results achieved, based on the tool results. Do NOT call any tools — reply with the summary text only."})
					synthCtx, synthCancel := context.WithTimeout(context.Background(), 90*time.Second)
					if reply2, _, err2 := models.CallModelWithFallback(synthCtx, "chat", synthMsgs); err2 == nil && strings.TrimSpace(reply2) != "" {
						cleanReply = strings.TrimSpace(reply2)
					}
					synthCancel()
				}
			}
			if strings.TrimSpace(cleanReply) == "" {
				if strings.TrimSpace(lastFullThought) != "" {
					// The model already WROTE its final report as thought text but
					// never called complete_task — deliver that text instead of a
					// truncated summary or an empty bubble.
					cleanReply = lastFullThought + "\n\n_(auto-completed: model did not call complete_task)_"
				} else {
					summary := "[NO-FINAL-REPORT] The model did not produce a final report. Summary of what actually executed:"
					tail := thinkingLines
					if len(tail) > 12 {
						tail = tail[len(tail)-12:]
					}
					for _, l := range tail {
						summary += "\n" + l
					}
					summary += "\n\n💡 Retry the command or continue with \"continue\"."
					cleanReply = summary
				}
			}

			if planHandback != "" {
				cleanReply = strings.TrimSpace(cleanReply) + planHandback
			}

			sendScorpReply(chatID, msgID, cleanReply)
			maybeRunSelfReview(chatID, chatIDStr)
			return
		}

		// Reset retry counter on successful tool call
		noToolRetries = 0

		// Record assistant message with tool calls
		history = append(history, AgentMessage{Role: "assistant", Content: reply})
		appendSessionHistory(chatIDStr, AgentMessage{Role: "assistant", Content: reply})

		// ── Execute Action Tools ──
		for _, tc := range actionToolCalls {
			// Real-time Steering Queue check (PicoClaw Parity)
			steerMsg, hasSteer := PopSteeringMessage(chatIDStr)
			if !hasSteer && rawChatIDStr != chatIDStr && chatID != 0 {
				steerMsg, hasSteer = PopSteeringMessage(rawChatIDStr)
			}
			if hasSteer {
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
			if tc.Name == "shell" && config.ConfirmationRequired() && IsDangerousCommand(helpers.GetStringArg(tc.Args, "command", "")) {
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

				StorePendingConfirmation(fmt.Sprintf("%d", chatID), "shell", cmd, history, promptMsgID)
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

			// Execute tool via registry; task_plan is intercepted so the
			// ledger is scoped to the session name, not the numeric chatID.
			var result string
			var ok bool
			if tc.Name == "task_plan" {
				result, ok = execTaskPlanTool(chatIDStr, tc.Args)
				if p := GetTaskPlan(chatIDStr); p != nil {
					done, total := p.Progress()
					thinkingLines = append(thinkingLines, fmt.Sprintf("📋 Plan: %d/%d done", done, total))
				}
			} else {
				result, ok = ExecuteTool(tc, chatID)
			}
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
		skills.TickActiveSkills()
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

	// Wrap resumed execution with safety wall-clock timeout
	ctx, cancel := context.WithTimeout(context.Background(), maxTurnTimeout())
	defer cancel()

	if msgID == 0 {
		msgID = tools.SendMessageGetID("🧠 <b>Agent</b>\n\n⏳ <i>continuing...</i>", chatID)
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

		var actionToolCalls []ToolCall
		isTaskCompleted := false
		var explicitFinalResult string

		for _, tc := range toolCalls {
			if tc.Name == "complete_task" {
				isTaskCompleted = true
				if res, ok := tc.Args["result"].(string); ok && res != "" {
					explicitFinalResult = res
				}
			} else {
				actionToolCalls = append(actionToolCalls, tc)
			}
		}

		if isTaskCompleted {
			// Task Ledger gate (resume path): same contract as the main loop —
			// an unfinished plan means the request is not fulfilled.
			if plan := GetTaskPlan(chatIDStr); plan != nil {
				if un := plan.Unfinished(); len(un) > 0 {
					log.Printf("[agent] resume complete_task rejected — %d unfinished plan items (Task Ledger gate)", len(un))
					gateNudge := fmt.Sprintf("⚠️ COMPLETE_TASK REJECTED — the Task Ledger still has %d unfinished item(s):\n%s\nContinue executing them NOW with real tools; complete_task is accepted ONLY when every item is done.", len(un), renderPlanItems(un))
					messages = append(messages, AgentMessage{Role: "assistant", Content: reply})
					messages = append(messages, AgentMessage{Role: "user", Content: gateNudge})
					continue
				}
				ClearTaskPlan(chatIDStr)
			}

			finalOutput := explicitFinalResult
			if finalOutput == "" {
				finalOutput = cleanToolCallTags(reply)
			}
			messages = append(messages, AgentMessage{Role: "assistant", Content: reply})
			appendSessionHistory(chatIDStr, AgentMessage{Role: "assistant", Content: finalOutput})
			sendScorpReply(chatID, msgID, finalOutput)
			maybeRunSelfReview(chatID, chatIDStr)
			return
		}

		if len(actionToolCalls) == 0 {
			cleanReply := cleanToolCallTags(reply)

			if noToolRetries < 4 {
				noToolRetries++
				retryReason := "⚠️ TASK IN-PROGRESS: In Agent Mode, you must call an action tool to continue, or call 'complete_task' with your final answer if finished."
				thinkingLines = append(thinkingLines, fmt.Sprintf("💭 %s", helpers.TruncateStr(cleanReply, 60)))
				tools.EditMessageByID(chatID, msgID, buildThinkingMessage(thinkingLines, time.Since(start), false), nil)
				messages = append(messages, AgentMessage{Role: "assistant", Content: reply})
				messages = append(messages, AgentMessage{Role: "user", Content: retryReason})
				continue
			}

			messages = append(messages, AgentMessage{Role: "assistant", Content: reply})
			appendSessionHistory(chatIDStr, AgentMessage{Role: "assistant", Content: cleanReply})
			sendScorpReply(chatID, msgID, cleanReply)
			maybeRunSelfReview(chatID, chatIDStr)
			return
		}

		messages = append(messages, AgentMessage{Role: "assistant", Content: reply})
		appendSessionHistory(chatIDStr, AgentMessage{Role: "assistant", Content: reply})

		for _, tc := range actionToolCalls {
			toolCount++
			desc := toolDescription(tc)
			thinkingLines = append(thinkingLines, desc)

			if tc.Name == "shell" && config.ConfirmationRequired() && IsDangerousCommand(helpers.GetStringArg(tc.Args, "command", "")) {
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

		registry.TickToolTTL()
		skills.TickActiveSkills()
	}

	tools.EditMessageByID(chatID, msgID, "⚠️ Agent reached maximum iterations.", nil)
}

func getSessionSearchContext(sessionID string, userMessage string) string {
	queryTokens := strings.Fields(userMessage)
	var meaningful []string
	for _, t := range queryTokens {
		t = strings.ToLower(strings.Trim(t, `.,!?:;"'()[]{}/*`))
		if len(t) > 2 {
			meaningful = append(meaningful, t)
		}
	}
	if len(meaningful) == 0 {
		return ""
	}
	query := strings.Join(meaningful, " OR ")
	results := session.SearchSessions(query, 3)
	if len(results) == 0 {
		return ""
	}
	var sb strings.Builder
	for _, r := range results {
		content := r.Content
		if len(content) > 200 {
			content = content[:200] + "..."
		}
		sb.WriteString(fmt.Sprintf("[%s]: %s\n", r.Role, content))
	}
	formatted := sb.String()
	if len(formatted) > 800 {
		formatted = formatted[:800] + "..."
	}
	return fmt.Sprintf("\n\n### Relevant Past Session Messages\n%s\n", formatted)
}
