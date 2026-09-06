package agent

import (
	"fmt"
	"sync"
	"time"

	"scorp-agent/internal/helpers"
	"scorp-agent/tools"
)

// ──────────────────────────────────────────────
// Dangerous Command Confirmation System
// Pauses loop execution until user confirms y/n
// ──────────────────────────────────────────────

type pendingConfirmation struct {
	toolName    string
	command     string
	args        map[string]interface{} // full tool args (P3.13 auto-mode confirmations); nil = legacy shell path
	messages    []AgentMessage
	created     time.Time
	promptMsgID int64
}

var (
	pendingConfirms   = make(map[string]*pendingConfirmation)
	pendingConfirmsMu sync.Mutex
)

// StorePendingConfirmation records a pending command confirmation
func StorePendingConfirmation(chatID, toolName, command string, messages []AgentMessage, promptMsgID ...int64) {
	StorePendingConfirmationArgs(chatID, toolName, command, nil, messages, promptMsgID...)
}

// StorePendingConfirmationArgs records a pending confirmation carrying full
// tool args (P3.13): on approval, any tool can be executed — not just shell.
func StorePendingConfirmationArgs(chatID, toolName, command string, args map[string]interface{}, messages []AgentMessage, promptMsgID ...int64) {
	pendingConfirmsMu.Lock()
	defer pendingConfirmsMu.Unlock()
	var pMsgID int64
	if len(promptMsgID) > 0 {
		pMsgID = promptMsgID[0]
	}
	pendingConfirms[chatID] = &pendingConfirmation{
		toolName:    toolName,
		command:     command,
		args:        args,
		messages:    messages,
		created:     time.Now(),
		promptMsgID: pMsgID,
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

// HasPendingConfirmation checks if there is an active pending confirmation for a chat
func HasPendingConfirmation(chatID string) bool {
	return getPendingConfirmation(chatID) != nil
}

// GetPendingConfirmationDetails returns details of pending confirmation
func GetPendingConfirmationDetails(chatID string) (command string, toolName string, exists bool) {
	pc := getPendingConfirmation(chatID)
	if pc == nil {
		return "", "", false
	}
	return pc.command, pc.toolName, true
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

// HandleConfirmation processes a user's yes/no response to a pending confirmation
func HandleConfirmation(chatID int64, confirmed bool, callbackMsgID ...int64) {
	chatIDStr := fmt.Sprintf("%d", chatID)
	pc := getPendingConfirmation(chatIDStr)

	if pc == nil {
		tools.SendMessage("❌ No pending confirmation found. It may have expired.", nil)
		return
	}

	clearPendingConfirmation(chatIDStr)

	// Clean up inline keyboard buttons if callback message ID provided
	if len(callbackMsgID) > 0 && callbackMsgID[0] != 0 {
		statusText := "❌ Command Denied"
		if confirmed {
			statusText = "✅ Command Approved"
		}
		tools.EditMessageByID(chatID, callbackMsgID[0], fmt.Sprintf("%s\n\n<pre>%s</pre>", statusText, helpers.EscapeHTML(pc.command)), nil)
	}

	if !confirmed {
		tools.SendMessage("❌ Command rejected.", nil)
		toolResult := fmt.Sprintf("[Tool Result: %s]\nUser REJECTED the command: %s\nPlease suggest an alternative approach.", pc.toolName, pc.command)
		pc.messages = append(pc.messages, AgentMessage{Role: "user", Content: toolResult})
		appendSessionHistory(chatIDStr, AgentMessage{Role: "user", Content: toolResult})
		go resumeAgentLoop(chatID, pc.messages, 0)
		return
	}

	// P3.13: auto-mode confirmations carry full args and may target any tool.
	// Execute via ExecuteTool so deny rules, hooks, receipts all apply. The
	// user just approved THIS call — preset AutoDecision so ExecuteTool's auto
	// gate trusts it instead of re-classifying (a re-grade could contradict
	// the human's approval); deny rules and hooks still bind.
	if pc.args != nil {
		tc := ToolCall{Name: pc.toolName, Args: pc.args, AutoDecision: "auto:user-approved"}
		result, ok := ExecuteTool(tc, chatID)
		status := "✅"
		if !ok {
			status = "❌"
		}
		toolResult := fmt.Sprintf("[Tool Result: %s]\n%s%s", pc.toolName, status, result)
		pc.messages = append(pc.messages, AgentMessage{Role: "user", Content: toolResult})
		appendSessionHistory(chatIDStr, AgentMessage{Role: "user", Content: toolResult})
		tools.SendMessage(fmt.Sprintf("%s Executed after approval. Continuing...", status), nil)
		go resumeAgentLoop(chatID, pc.messages, 0)
		return
	}

	// User confirmed — execute command
	shellArgs := map[string]interface{}{"command": pc.command, "timeout": 60, "confirmed": true}
	// P3.12: PreToolUse hooks also bind user-confirmed resumes (same contract
	// as deny rules) — user confirmation relaxes the danger gate, not hook
	// policy.
	if _, hookBlocked, hookReason := tools.RunPreToolUseHooks("shell", shellArgs, chatID); hookBlocked {
		toolResult := fmt.Sprintf("[Tool Result: %s]\n🪝 BLOCKED by PreToolUse hook: %s\nThe user approved the command, but a hook policy blocked it. Please suggest an alternative approach.", pc.toolName, hookReason)
		pc.messages = append(pc.messages, AgentMessage{Role: "user", Content: toolResult})
		appendSessionHistory(chatIDStr, AgentMessage{Role: "user", Content: toolResult})
		tools.SendMessage(fmt.Sprintf("🪝 Blocked by hook policy: %s", hookReason), nil)
		go resumeAgentLoop(chatID, pc.messages, 0)
		return
	}
	result, ok := tools.ExecuteShell(shellArgs, chatID)
	// Record the receipt here too: this path bypasses agent.ExecuteTool, and
	// the test-integrity gate (P0.3) needs confirmed test runs as evidence.
	tools.RecordToolReceipt("shell", shellArgs, result, ok)
	// P3.12: post-tool-use hook context on the confirmed path as well.
	if postCtx := tools.RunPostToolUseHooks("shell", shellArgs, result, ok, chatID); postCtx != "" {
		result += "\n\n🪝 " + postCtx
	}
	status := "✅"
	if !ok {
		status = "❌"
	}

	toolResult := fmt.Sprintf("[Tool Result: %s]\n%s%s", pc.toolName, status, result)
	pc.messages = append(pc.messages, AgentMessage{Role: "user", Content: toolResult})
	appendSessionHistory(chatIDStr, AgentMessage{Role: "user", Content: toolResult})

	tools.SendMessage(fmt.Sprintf("%s Command executed. Continuing...", status), nil)
	go resumeAgentLoop(chatID, pc.messages, 0)
}
