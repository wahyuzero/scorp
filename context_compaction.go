package main

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
	maxHistoryTokens = 24000
	// When token count exceeds this, trigger compaction
	compactionThreshold = 18000
	// Tool results longer than this get truncated in the history
	maxToolResultChars = 2000
)

// estimateTokens gives a rough token count for a string.
func estimateTokens(s string) int {
	return int(float64(len(s)) * tokensPerChar)
}

// estimateHistoryTokens calculates total estimated tokens for a history slice.
func estimateHistoryTokens(history []agentMessage) int {
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

// truncateToolResultsInHistory shortens overly long tool results to save context.
// Returns the modified history and number of truncations made.
func truncateToolResultsInHistory(history []agentMessage) ([]agentMessage, int) {
	truncations := 0
	for i, msg := range history {
		if msg.Role != "user" {
			continue
		}
		content, ok := msg.Content.(string)
		if !ok {
			continue
		}
		// Check if it's a tool result (they start with "[Tool Result:")
		if strings.HasPrefix(content, "[Tool Result:") && len(content) > maxToolResultChars {
			// Keep first and last portions, truncate middle
			head := maxToolResultChars * 3 / 4
			tail := maxToolResultChars / 4
			truncated := content[:head] + "\n...[truncated for context efficiency]...\n" + content[len(content)-tail:]
			history[i] = agentMessage{Role: msg.Role, Content: truncated}
			truncations++
		}
	}
	return history, truncations
}

// maybeCompactHistory checks if the conversation needs compaction based on token count.
// If so, it triggers truncation of tool results first, then summarization if still over budget.
// Called from the agent loop before each model call.
func maybeCompactHistory(chatID string, history []agentMessage) []agentMessage {
	tokens := estimateHistoryTokens(history)

	if tokens <= compactionThreshold {
		return history
	}

	log.Printf("[compaction] History at ~%d tokens (threshold %d), compacting for %s",
		tokens, compactionThreshold, chatID)

	// Step 1: Truncate long tool results
	history, n := truncateToolResultsInHistory(history)
	if n > 0 {
		log.Printf("[compaction] Truncated %d long tool results", n)
		// Update session with truncated history
		sess := getOrCreateSession(chatID)
		sess.history = history
		setSession(chatID, sess)
	}

	// Re-check token count
	tokens = estimateHistoryTokens(history)
	if tokens <= compactionThreshold {
		log.Printf("[compaction] After truncation: ~%d tokens (under threshold)", tokens)
		return history
	}

	// Step 2: Force summarization if still over budget
	log.Printf("[compaction] Still at ~%d tokens after truncation, summarizing old messages", tokens)
	summarizeHistory(chatID)

	// Return updated history from session
	sess := getSession(chatID)
	if sess != nil {
		return sess.history
	}
	return history
}

// formatTokenEstimate returns a human-readable token estimate for debugging
func formatTokenEstimate(history []agentMessage) string {
	tokens := estimateHistoryTokens(history)
	chars := 0
	for _, msg := range history {
		if c, ok := msg.Content.(string); ok {
			chars += len(c)
		}
	}
	return fmt.Sprintf("%d messages, ~%d chars, ~%d tokens", len(history), chars, tokens)
}
