package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ──────────────────────────────────────────────
// Clarify Tool — Ask user questions
// ──────────────────────────────────────────────

type PendingClarify struct {
	ChatID   string
	Question string
	Choices  []string
	Response chan string // Agent blocks until response arrives
	Created  time.Time
	MsgID    int64 // Message ID of the question sent
}

var (
	pendingClarifies   = make(map[string]*PendingClarify)
	pendingClarifiesMu sync.Mutex
)

func executeClarify(args map[string]interface{}) (string, bool) {
	// This function is called by the agent. It should:
	// 1. Send the question to the user
	// 2. Wait for response via channel
	// 3. Return the response

	question, _ := args["question"].(string)
	if question == "" {
		return "Missing 'question' parameter", false
	}

	// choices is optional
	var choices []string
	if choicesRaw, ok := args["choices"]; ok {
		if choicesArr, ok := choicesRaw.([]interface{}); ok {
			for _, c := range choicesArr {
				if s, ok := c.(string); ok {
					choices = append(choices, s)
				}
			}
		}
	}

	// Create response channel
	respChan := make(chan string, 1)

	// We need the chat ID — it's passed to executeTool via the tool execution context
	// But executeClarify doesn't receive chatID directly.
	// We use a global mapping from the goroutine context.
	// The chat ID is set via setClarifyChatID when the tool is dispatched.

	chatID := getClarifyChatID()
	if chatID == "" {
		return "Error: no chat context", false
	}

	// Send question to user
	chatIDInt, _ := strconv.ParseInt(chatID, 10, 64)
	msgID := sendClarifyMessage(chatIDInt, question, choices)
	if msgID == 0 {
		return "Failed to send question", false
	}

	// Store pending clarify
	pendingClarifiesMu.Lock()
	pendingClarifies[chatID] = &PendingClarify{
		ChatID:   chatID,
		Question: question,
		Choices:  choices,
		Response: respChan,
		Created:  time.Now(),
		MsgID:    msgID,
	}
	pendingClarifiesMu.Unlock()

	// Wait for response with timeout
	select {
	case resp := <-respChan:
		// Clean up
		pendingClarifiesMu.Lock()
		delete(pendingClarifies, chatID)
		pendingClarifiesMu.Unlock()
		return fmt.Sprintf("[User Response] %s", resp), true

	case <-time.After(120 * time.Second):
		// Timeout
		pendingClarifiesMu.Lock()
		delete(pendingClarifies, chatID)
		pendingClarifiesMu.Unlock()
		return "[User Response] No response (timeout)", true
	}
}

// sendClarifyMessage sends a question to the user, optionally with inline keyboard
func sendClarifyMessage(chatID int64, question string, choices []string) int64 {
	if len(choices) > 0 {
		// Build inline keyboard
		var rows []interface{}
		for _, choice := range choices {
			rows = append(rows, []interface{}{
				map[string]string{"text": choice, "callback_data": "clarify:" + choice},
			})
		}
		// Add cancel button
		rows = append(rows, []interface{}{
			map[string]string{"text": "❌ Cancel", "callback_data": "clarify:cancel"},
		})

		keyboard := map[string]interface{}{
			"inline_keyboard": rows,
		}

		// Send as inline keyboard
		msgID := sendMessageGetID("🤔 <b>Clarify</b>\n\n"+question, chatID)
		if msgID != 0 {
			editMessageByID(chatID, msgID, "🤔 <b>Clarify</b>\n\n"+question, keyboard)
		}
		return msgID
	}

	// Open-ended: just send message, wait for text reply
	return sendMessageGetID("🤔 <b>Clarify</b>\n\n"+question+"\n\n<i>Reply with your answer.</i>", chatID)
}

// ──────────────────────────────────────────────
// Clarify Chat ID Context
// ──────────────────────────────────────────────

var (
	clarifyChatID   string
	clarifyChatIDMu sync.Mutex
)

func setClarifyChatID(chatID string) {
	clarifyChatIDMu.Lock()
	clarifyChatID = chatID
	clarifyChatIDMu.Unlock()
}

func getClarifyChatID() string {
	clarifyChatIDMu.Lock()
	defer clarifyChatIDMu.Unlock()
	return clarifyChatID
}

// ──────────────────────────────────────────────
// Clarify Response Handlers
// ──────────────────────────────────────────────

// handleClarifyResponse processes a user response to a clarify question
// and resumes the agent loop with the response.
func handleClarifyResponse(chatID string, response string) {
	pendingClarifiesMu.Lock()
	pending, ok := pendingClarifies[chatID]
	if !ok {
		pendingClarifiesMu.Unlock()
		return
	}
	pendingClarifiesMu.Unlock()

	// Send response through channel (non-blocking)
	select {
	case pending.Response <- response:
	default:
	}
}

// resolveClarify checks for pending clarify and returns true if it was resolved
func resolveClarify(action string, chatID string, callbackID string) bool {
	pendingClarifiesMu.Lock()
	_, ok := pendingClarifies[chatID]
	pendingClarifiesMu.Unlock()

	if !ok {
		return false
	}

	// If it's a callback with "clarify:" prefix
	if strings.HasPrefix(action, "clarify:") {
		response := strings.TrimPrefix(action, "clarify:")
		if response == "cancel" {
			response = "cancelled"
		}
		if callbackID != "" {
			answerCallback(callbackID, "")
		}
		handleClarifyResponse(chatID, response)
		return true
	}

	// Open-ended — any text message from the user while a clarify is pending
	if callbackID == "" {
		handleClarifyResponse(chatID, action)
		return true
	}

	return false
}

// hasPendingClarify checks if a chat has a pending clarify question
func hasPendingClarify(chatID string) bool {
	pendingClarifiesMu.Lock()
	defer pendingClarifiesMu.Unlock()
	_, ok := pendingClarifies[chatID]
	return ok
}

// init — register the clarify tool
func init() {
	registerTool(ToolDef{
		Name:        "clarify",
		Description: "Ask the user a question when you need clarification, feedback, or a decision. Supports multiple choice (up to 4) or open-ended text. Blocks until user responds or 120s timeout.",
		Category:    "other",
		Native:      true,
		Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
			setClarifyChatID(fmt.Sprintf("%d", chatID))
			return executeClarify(args)
		},
		Arguments: map[string]ArgDef{
			"question": {Type: "string", Description: "The question to ask the user", Required: true},
			"choices":  {Type: "array", Description: "Up to 4 answer choices for inline keyboard (optional). If omitted, user types free-form response."},
		},
	})
}

// clean up expired clarifies periodically
func cleanupStaleClarifies() {
	pendingClarifiesMu.Lock()
	defer pendingClarifiesMu.Unlock()

	now := time.Now()
	for chatID, pc := range pendingClarifies {
		if now.Sub(pc.Created) > 130*time.Second {
			select {
			case pc.Response <- "timeout":
			default:
			}
			delete(pendingClarifies, chatID)
			log.Printf("[clarify] Cleaned up stale clarify for %s", chatID)
		}
	}
}
