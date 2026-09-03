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
	pendingConfirmsMu.Lock()
	defer pendingConfirmsMu.Unlock()
	var pMsgID int64
	if len(promptMsgID) > 0 {
		pMsgID = promptMsgID[0]
	}
	pendingConfirms[chatID] = &pendingConfirmation{
		toolName:    toolName,
		command:     command,
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

	// User confirmed — execute command
	result, ok := tools.ExecuteShell(map[string]interface{}{"command": pc.command, "timeout": 60, "confirmed": true}, chatID)
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
