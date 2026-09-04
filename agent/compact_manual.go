package agent

import (
	"fmt"
	"log"
	"strings"
)

// ──────────────────────────────────────────────
// Manual Compaction Engine
// Triggers age-aware tool result compression and LLM-based summarization on-demand.
// ──────────────────────────────────────────────

// CompactStats holds before-and-after metrics of manual compaction
type CompactStats struct {
	InitialTokens   int
	FinalTokens     int
	TokensSaved     int
	PercentSaved    float64
	InitialMessages int
	FinalMessages   int
	ToolPrunedCount int
	Summarized      bool
	Error           error
}

// CompactSessionHistory performs immediate manual compaction for a given session.
// Returns detailed metrics on token reduction and message pruning.
func CompactSessionHistory(sessionID string) CompactStats {
	sess := getOrCreateSession(sessionID)
	history := sess.history
	if len(history) == 0 {
		history = loadHistoryFromDisk(sessionID)
		sess.history = history
	}

	initialMsgs := len(history)
	initialTokens := estimateHistoryTokens(history)

	if initialMsgs == 0 {
		return CompactStats{
			InitialTokens: 0,
			FinalTokens:   0,
		}
	}

	// 1. Force age-aware tool result pruning
	prunedHistory, prunedCount := truncateToolResultsInHistory(history)
	sess.history = prunedHistory
	setSession(sessionID, sess)

	summarized := false
	// 2. If there are enough messages to summarize (> 4 messages), perform LLM summarization
	if len(prunedHistory) > 4 {
		log.Printf("[compaction] Manual compact: summarizing %d messages for %s", len(prunedHistory), sessionID)
		summarizeHistory(sessionID)
		summarized = true
	}

	// Reload updated history from session
	sess = getSession(sessionID)
	finalHistory := sess.history
	finalMsgs := len(finalHistory)
	finalTokens := estimateHistoryTokens(finalHistory)

	tokensSaved := initialTokens - finalTokens
	if tokensSaved < 0 {
		tokensSaved = 0
	}
	pctSaved := 0.0
	if initialTokens > 0 {
		pctSaved = (float64(tokensSaved) / float64(initialTokens)) * 100.0
	}

	return CompactStats{
		InitialTokens:   initialTokens,
		FinalTokens:     finalTokens,
		TokensSaved:     tokensSaved,
		PercentSaved:    pctSaved,
		InitialMessages: initialMsgs,
		FinalMessages:   finalMsgs,
		ToolPrunedCount: prunedCount,
		Summarized:      summarized,
	}
}

// FormatCompactStats builds a clean terminal or Telegram friendly report
func FormatCompactStats(sessionID string, stats CompactStats) string {
	if stats.InitialMessages == 0 {
		return fmt.Sprintf("🗜️ Session <code>%s</code> is already empty. Nothing to compact.", sessionID)
	}

	var sb strings.Builder
	sb.WriteString("🗜️ <b>Session Compaction Complete!</b>\n\n")
	sb.WriteString(fmt.Sprintf("• Session: <b>%s</b>\n", sessionID))
	sb.WriteString(fmt.Sprintf("• Tool Results Pruned: <b>%d</b>\n", stats.ToolPrunedCount))
	sb.WriteString(fmt.Sprintf("• Messages: <b>%d ➔ %d</b>\n", stats.InitialMessages, stats.FinalMessages))
	sb.WriteString(fmt.Sprintf("• Token Footprint: <b>%d ➔ %d tokens</b>\n", stats.InitialTokens, stats.FinalTokens))
	sb.WriteString(fmt.Sprintf("• Reduction: <b>%.1f%% saved</b> (~%d tokens)\n\n", stats.PercentSaved, stats.TokensSaved))
	sb.WriteString("✓ <i>Context memory refreshed. Essential facts and summaries preserved!</i>")

	return sb.String()
}
