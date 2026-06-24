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
	// Target token budget for conversation history (leaves room for system prompt + tools)
	maxHistoryTokens = 32000
	// When token count exceeds this, trigger compaction
	compactionThreshold = 28000
	// Tool results older than this many messages get progressively trimmed
	oldToolResultTrimAge = 6
	// Tool results older than this get reduced to a one-line stub
	veryOldToolResultAge = 12
	// Max chars for recent tool results (within trim age)
	recentToolResultMax = 3000
	// Max chars for old tool results (past trim age)
	oldToolResultMax = 500
)

// estimateTokens gives a rough token count for a string.
func estimateTokens(s string) int {
	return int(float64(len(s)) * tokensPerChar)
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
			// Stub: first line + char count
			firstLine := content
			if idx := strings.Index(content, "\n"); idx > 0 {
				firstLine = content[:idx]
			}
			if len(firstLine) > 100 {
				firstLine = firstLine[:100]
			}
			truncated = firstLine + fmt.Sprintf(" ...[%d chars trimmed]", len(content))
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

// maybeCompactHistory checks if the conversation needs compaction based on token count.
// If so, it triggers truncation of tool results first, then summarization if still over budget.
// Called from the agent loop before each model call.
// During an active agent loop, we allow aggressive message truncation (keep recent N + system)
// but defer LLM-based summarization to avoid race conditions and context loss mid-task.
func maybeCompactHistory(chatID string, history []AgentMessage) []AgentMessage {
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
		return history
	}

	log.Printf("[compaction] History at ~%d tokens (threshold %d), compacting for %s",
		tokens, compactionThreshold, chatID)

	// Re-check token count
	tokens = estimateHistoryTokens(history)
	if tokens <= compactionThreshold {
		log.Printf("[compaction] After pruning: ~%d tokens (under threshold)", tokens)
		return history
	}

	// Step 2: Check if we're in an active agent loop
	sess := getSession(chatID)
	if sess != nil && sess.loopActive {
		// During active loop: aggressively truncate to keep only recent messages + system prompt
		// This prevents session bloat without LLM summarization race conditions
		const keepDuringActiveLoop = 30 // keep last 30 messages + system prompt
		if len(history) > keepDuringActiveLoop+1 { // +1 for system prompt
			// Find system prompt (first message with role="system")
			systemIdx := -1
			for i, msg := range history {
				if msg.Role == "system" {
					systemIdx = i
					break
				}
			}

			var newHistory []AgentMessage
			if systemIdx >= 0 {
				// Keep system prompt + last keepDuringActiveLoop messages
				newHistory = append(newHistory, history[systemIdx])
				start := len(history) - keepDuringActiveLoop
				if start <= systemIdx {
					start = systemIdx + 1
				}
				newHistory = append(newHistory, history[start:]...)
			} else {
				// No system prompt, just keep last keepDuringActiveLoop messages
				newHistory = history[len(history)-keepDuringActiveLoop:]
			}

			log.Printf("[compaction] Active loop: truncated from %d to %d messages (kept system + last %d)",
				len(history), len(newHistory), keepDuringActiveLoop)

			// Update session
			sess.history = newHistory
			setSession(chatID, sess)
			return newHistory
		}

		log.Printf("[compaction] Still at ~%d tokens but agent loop active and history small — deferring LLM summarization", tokens)
		return history
	}

	// Step 3: Full summarization (only when NOT in agent loop)
	log.Printf("[compaction] Still at ~%d tokens after truncation, summarizing old messages", tokens)
	summarizeHistory(chatID)

	// Return updated history from session
	sess = getSession(chatID)
	if sess != nil {
		return sess.history
	}
	return history
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
