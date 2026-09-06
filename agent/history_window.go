package agent

// ──────────────────────────────────────────────
// New-Request Roll-Up
//
// Failure mode (verify26 turn 19): a brand-new user request landed on a
// session whose recent history was saturated by the PREVIOUS (closed) task.
// The tail — dozens of tool results from the finished task — outweighed the
// single new message, and the model re-executed the old task instead of
// reading the new request.
//
// Fix: when a NEW request arrives (not a continuation, no active Task Ledger),
// the previous closed task is represented compactly (system prompt + its final
// report + a bounded recent window) and a hard turn boundary is injected so
// the new instruction is the dominant, unambiguous signal.
// ──────────────────────────────────────────────

import (
	"fmt"
	"strings"
)

// newRequestTokenBudget: above this estimate the old-task history is rolled up
// before a new request is processed. Kept well below compactionThreshold so
// the roll-up happens before generic compaction would mutilate context.
const newRequestTokenBudget = 6000

// keepRecentWindow: how many recent messages survive the roll-up. Enough for
// conversational follow-ups ("what did you just do?"), small enough that the
// new request dominates the tail.
const keepRecentWindow = 8

// turnBoundaryMarker is injected as a user-role system note (mid-conversation
// role conventions are already used by gate nudges) directly before the new
// user message.
const turnBoundaryMarker = "[📋 TURN BOUNDARY] The task above is COMPLETED and CLOSED. The next user message is a NEW, SEPARATE request. Do NOT repeat, redo, or continue any previous work. Read the new request below and respond to IT only."

// prepareNewTurnHistory returns the history slice to send to the model for a
// brand-new user request. Continuations and sessions with an active Task
// Ledger keep their full context (the contract depends on it).
func prepareNewTurnHistory(sessionID string, history []AgentMessage, isContinuation bool) []AgentMessage {
	if isContinuation {
		return history
	}
	if plan := GetTaskPlan(sessionID); plan != nil && len(plan.Unfinished()) > 0 {
		return history // active work: full context IS the contract
	}
	if estimateHistoryTokens(history) <= newRequestTokenBudget {
		return history
	}

	// Locate the last substantial assistant message (the previous task's final
	// report) so the roll-up preserves a summary of what was done.
	lastReportIdx := -1
	for i := len(history) - 1; i >= 0; i-- {
		if history[i].Role != "assistant" {
			continue
		}
		if c, ok := history[i].Content.(string); ok && len(strings.TrimSpace(c)) > 80 {
			lastReportIdx = i
			break
		}
	}

	var out []AgentMessage
	if len(history) > 0 && history[0].Role == "system" {
		out = append(out, history[0])
	}
	if lastReportIdx > 0 {
		out = append(out, AgentMessage{Role: "user", Content: "[context note] For reference, the final report of the previous (completed) task:"})
		out = append(out, history[lastReportIdx])
	}

	// Bounded recent window (excluding the report already copied).
	start := len(history) - keepRecentWindow
	if start < lastReportIdx+1 {
		start = lastReportIdx + 1
	}
	if start < 0 {
		start = 0
	}
	if start > len(history) {
		start = len(history)
	}
	dropped := start - len(out) + 1
	if dropped > 0 {
		out = append(out, AgentMessage{Role: "user", Content: fmt.Sprintf("[context note] %d older execution messages were rolled up to keep the context window focused.", dropped)})
	}
	out = append(out, history[start:]...)
	out = append(out, AgentMessage{Role: "user", Content: turnBoundaryMarker})
	return out
}
