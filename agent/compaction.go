package agent

import (
	"fmt"
	"log"
	"strings"
)

// ──────────────────────────────────────────────
// Context Compaction — token-aware conversation trimming
// Enhances the existing message-count-based summarizeHistory()
// with token estimation and smart tool-result truncation.
// ──────────────────────────────────────────────

const (
	// Approximate tokens-per-char ratio (rough estimate for mixed content)
	tokensPerChar = 0.25
	// Target token budget for conversation history (keeps requests light & fast)
	maxHistoryTokens = 16000
	// When token count exceeds this, trigger compaction
	compactionThreshold = 12000
	// Tool results older than this many messages get progressively trimmed
	oldToolResultTrimAge = 3
	// Tool results older than this get reduced to a one-line stub
	veryOldToolResultAge = 6
	// Max chars for recent tool results (within trim age)
	recentToolResultMax = 2000
	// Max chars for old tool results (past trim age)
	oldToolResultMax = 400
)

// estimateTokens calculates an accurate token count estimate for a string.
// Code, JSON, and non-ASCII characters have higher token density than plain English.
func estimateTokens(s string) int {
	if len(s) == 0 {
		return 0
	}
	base := float64(len(s)) * tokensPerChar

	// Check if content has high code/JSON density (braces, brackets, punctuation)
	codeDensityPunct := 0
	for _, r := range s {
		if r == '{' || r == '}' || r == '[' || r == ']' || r == ';' || r == ':' || r == '=' || r == '"' || r == '\\' {
			codeDensityPunct++
		}
	}
	// If punctuation accounts for >8% of text, bump token estimate by 25%
	if float64(codeDensityPunct)/float64(len(s)) > 0.08 {
		base *= 1.25
	}

	return int(base)
}

// estimateHistoryTokens calculates total estimated tokens for a history slice.
func estimateHistoryTokens(history []AgentMessage) int {
	total := 0
	for _, msg := range history {
		switch c := msg.Content.(type) {
		case string:
			total += estimateTokens(c)
		default:
			total += 200 // estimate for non-text content
		}
		// Overhead per message (role, formatting)
		total += 4
	}
	return total
}

// truncateToolResultsInHistory applies age-aware pruning to tool results.
// Recent tool results (within oldToolResultTrimAge messages from end) are kept
// at recentToolResultMax. Older ones are trimmed to oldToolResultMax. Very old
// ones are reduced to a one-line stub. This keeps conversation flow intact
// (model can still see what happened) while aggressively cutting token bloat
// from old tool outputs.
func truncateToolResultsInHistory(history []AgentMessage) ([]AgentMessage, int) {
	truncations := 0
	n := len(history)

	for i, msg := range history {
		if msg.Role != "user" {
			continue
		}
		content, ok := msg.Content.(string)
		if !ok {
			continue
		}
		if !strings.HasPrefix(content, "[Tool Result:") {
			continue
		}

		// Age = how many messages after this one (0 = most recent)
		age := n - 1 - i

		var maxLen int
		switch {
		case age < oldToolResultTrimAge:
			// Recent — keep generous
			maxLen = recentToolResultMax
		case age < veryOldToolResultAge:
			// Old — trim aggressively but keep some context
			maxLen = oldToolResultMax
		default:
			// Very old — reduce to stub only
			maxLen = 120
		}

		if len(content) <= maxLen {
			continue
		}

		var truncated string
		if maxLen <= 120 {
			// Extract tool name from header [Tool Result: <name>]
			toolName := "tool"
			if idx := strings.Index(content, "\n"); idx > 0 {
				header := content[:idx]
				header = strings.TrimPrefix(header, "[Tool Result:")
				header = strings.TrimSuffix(header, "]")
				toolName = strings.TrimSpace(header)
			}
			truncated = SummarizeOldToolResult(toolName, content)
		} else {
			// Keep head + tail
			head := maxLen * 3 / 4
			tail := maxLen / 4
			truncated = content[:head] + "\n...[trimmed]...\n" + content[len(content)-tail:]
		}
		history[i] = AgentMessage{Role: msg.Role, Content: truncated}
		truncations++
	}
	return history, truncations
}

// preservationMarker marks the always-survive context block injected after
// every compaction; re-compactions replace the old block instead of stacking.
const preservationMarker = "[🔒 PRESERVED CONTEXT]"

// preservedUserPrefixes are message shapes that are runtime machinery, not
// genuine user requests — skipped when hunting for the latest goal.
var preservedUserPrefixes = []string{
	"[Tool Result:", "⚠️", "🚨", "[⚡", "[📋", "[✅", "ℹ️ [Note:", preservationMarker,
}

// latestUserGoal returns the most recent genuine user request in history
// (runtime nudges and tool results excluded), truncated for context safety.
func latestUserGoal(history []AgentMessage) (string, int) {
	for i := len(history) - 1; i >= 0; i-- {
		if history[i].Role != "user" {
			continue
		}
		content, ok := history[i].Content.(string)
		if !ok {
			continue
		}
		skip := false
		for _, p := range preservedUserPrefixes {
			if strings.HasPrefix(content, p) {
				skip = true
				break
			}
		}
		if skip {
			continue
		}
		if len(content) > 600 {
			content = content[:600] + "..."
		}
		return content, i
	}
	return "", -1
}

// previousFinalReport returns the last assistant message that came before the
// current goal — the final report of the previous task (verify26: losing it
// made the model re-execute or contradict earlier results).
func previousFinalReport(history []AgentMessage) string {
	_, goalIdx := latestUserGoal(history)
	for i := goalIdx - 1; i >= 0; i-- {
		if history[i].Role != "assistant" {
			continue
		}
		content, ok := history[i].Content.(string)
		if !ok || strings.TrimSpace(content) == "" {
			continue
		}
		if len(content) > 600 {
			content = content[:600] + "..."
		}
		return content
	}
	return ""
}

// preservationNote assembles the block that MUST survive compaction: the
// active task ledger, the user's latest goal, and the previous task's final
// report. Returns "" when there is nothing worth preserving.
func preservationNote(chatID string, history []AgentMessage) string {
	var sb strings.Builder

	if plan := GetTaskPlan(chatID); plan != nil && plan.Total() > 0 {
		done, total := plan.Progress()
		sb.WriteString(fmt.Sprintf("📋 Active Task Ledger (%d/%d done) — keep working against it:\n%s\n", done, total, plan.Render()))
	}
	if goal, _ := latestUserGoal(history); goal != "" {
		sb.WriteString("🎯 Current user goal:\n" + goal + "\n")
	}
	if report := previousFinalReport(history); report != "" {
		sb.WriteString("📝 Previous task final report:\n" + report + "\n")
	}

	if sb.Len() == 0 {
		return ""
	}
	return preservationMarker + "\n" + sb.String()
}

// stripPreservationNotes drops old preserved blocks so a fresh one replaces
// them instead of accumulating.
func stripPreservationNotes(history []AgentMessage) []AgentMessage {
	out := history[:0]
	for _, msg := range history {
		if content, ok := msg.Content.(string); ok && strings.HasPrefix(content, preservationMarker) {
			continue
		}
		out = append(out, msg)
	}
	return out
}

// injectPreservationNote strips stale blocks and appends the given one (as a
// user-role message — system-role messages mid-history break some adapters).
// The note must be computed from the FULL pre-compaction history — after
// truncation the previous report is already gone.
func injectPreservationNote(history []AgentMessage, note string) []AgentMessage {
	if note == "" {
		return stripPreservationNotes(history)
	}
	history = stripPreservationNotes(history)
	return append(history, AgentMessage{Role: "user", Content: note})
}

// maybeCompactHistory checks if the conversation needs compaction based on
// token count. If so, it truncates tool results first, then applies
// structural truncation during active loops (keeping system + latest goal +
// recent window + the preservation block). Returns the (possibly new)
// history and a one-line notice for the Telegram thinking footer.
// LLM-based summarization only runs OUTSIDE active loops (race safety).
func maybeCompactHistory(chatID string, history []AgentMessage) ([]AgentMessage, string) {
	// Always apply age-aware tool result pruning (cheap, safe, prevents token bloat)
	history, n := truncateToolResultsInHistory(history)
	if n > 0 {
		log.Printf("[compaction] Age-pruned %d old tool results", n)
		sess := getOrCreateSession(chatID)
		sess.history = history
		setSession(chatID, sess)
	}

	tokens := estimateHistoryTokens(history)

	if tokens <= compactionThreshold {
		return history, ""
	}

	log.Printf("[compaction] History at ~%d tokens (threshold %d), compacting for %s",
		tokens, compactionThreshold, chatID)

	// Re-check token count
	tokens = estimateHistoryTokens(history)
	if tokens <= compactionThreshold {
		log.Printf("[compaction] After pruning: ~%d tokens (under threshold)", tokens)
		return history, ""
	}
	beforeTokens := tokens

	// Step 2: Check if we're in an active agent loop
	sess := getSession(chatID)
	if sess != nil && sess.loopActive {
		// During active loop: truncate to keep recent messages + system prompt
		// + preservation block (ledger, latest goal, previous final report).
		// This prevents context loss no matter how many iterations run.
		const keepDuringActiveLoop = 20 // keep last 20 messages
		if len(history) > keepDuringActiveLoop+2 {
			systemIdx := -1
			userGoalIdx := -1
			for i, msg := range history {
				if systemIdx < 0 && msg.Role == "system" {
					systemIdx = i
				}
				if userGoalIdx < 0 && msg.Role == "user" {
					if content, ok := msg.Content.(string); ok {
						if !strings.HasPrefix(content, "[Tool Result:") && !strings.HasPrefix(content, "⚠️") && !strings.HasPrefix(content, "🚨") {
							userGoalIdx = i
						}
					}
				}
			}

		// Compute the preservation note from the FULL history — after
		// truncation the previous report would already be gone.
		preserve := preservationNote(chatID, history)

		history = stripPreservationNotes(history)
		// indices shifted only if a note existed before userGoalIdx —
		// recompute the goal index against the stripped slice
		_, userGoalIdx = latestUserGoal(history)

			var newHistory []AgentMessage
			if systemIdx >= 0 {
				newHistory = append(newHistory, history[systemIdx])
			}
			if userGoalIdx >= 0 && userGoalIdx != systemIdx {
				newHistory = append(newHistory, history[userGoalIdx])
			}

			start := len(history) - keepDuringActiveLoop
			if start > userGoalIdx+1 {
				note := fmt.Sprintf("ℹ️ [Note: %d intermediate execution messages truncated to preserve context window. The agent loop is active and progressing.]", start-(userGoalIdx+1))
				newHistory = append(newHistory, AgentMessage{Role: "user", Content: note})
			} else {
				start = userGoalIdx + 1
			}

			if start < len(history) {
				newHistory = append(newHistory, history[start:]...)
			}
			newHistory = injectPreservationNote(newHistory, preserve)

			after := estimateHistoryTokens(newHistory)
			log.Printf("[compaction] Active loop: truncated from %d to %d messages (~%d → ~%d tokens; kept system + user goal + last %d + preserved context)",
				len(history), len(newHistory), beforeTokens, after, len(history)-start)

			// Update session
			sess.history = newHistory
			setSession(chatID, sess)
			return newHistory, fmt.Sprintf("🗜 compacted: %d → %d msgs (~%dk → ~%dk tokens)", len(history), len(newHistory), beforeTokens/1000, after/1000)
		}

		log.Printf("[compaction] Still at ~%d tokens but agent loop active and history small — deferring LLM summarization", tokens)
		return history, ""
	}

	// Step 3: Full summarization (only when NOT in agent loop) —
	// summarizeHistory computes the preservation note from the full history
	// itself and injects it into the summarized result.
	log.Printf("[compaction] Still at ~%d tokens after truncation, summarizing old messages", tokens)
	summarizeHistory(chatID)

	sess = getSession(chatID)
	if sess != nil {
		after := estimateHistoryTokens(sess.history)
		return sess.history, fmt.Sprintf("🗜 compacted: summarized (~%dk → ~%dk tokens)", beforeTokens/1000, after/1000)
	}
	return history, "🗜 compacted: summarized old messages"
}

// formatTokenEstimate returns a human-readable token estimate for debugging
func formatTokenEstimate(history []AgentMessage) string {
	tokens := estimateHistoryTokens(history)
	chars := 0
	for _, msg := range history {
		if c, ok := msg.Content.(string); ok {
			chars += len(c)
		}
	}
	return fmt.Sprintf("%d messages, ~%d chars, ~%d tokens", len(history), chars, tokens)
}
