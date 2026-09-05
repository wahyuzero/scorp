package agent

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"scorp-agent/models"
)

// ──────────────────────────────────────────────
// Automatic Contextual Session Auto-Titling
// Generates clean, human-readable session titles with spaces
// on the first turn of a conversation (e.g. "Deploy Docker VPS", "Audit Security")
// ──────────────────────────────────────────────

var invalidTitleChars = regexp.MustCompile(`[<>:"/\\|?*\x00-\x1f]`)

// ShouldAutoTitleSession determines if a session is unnamed/default and needs a title
func ShouldAutoTitleSession(sessionID string) bool {
	s := strings.TrimSpace(strings.ToLower(sessionID))
	if s == "" || s == "default" || s == "main_chat" || s == "main" || strings.HasPrefix(s, "chat-") || strings.HasPrefix(s, "untitled") {
		return true
	}
	// Also if it's purely a numeric chatID
	isNumeric := true
	for _, r := range s {
		if r < '0' || r > '9' {
			isNumeric = false
			break
		}
	}
	return isNumeric
}

// GenerateContextualSessionTitle generates a 2-4 word natural title based on user message
func GenerateContextualSessionTitle(userMessage string) string {
	msg := strings.TrimSpace(userMessage)
	if msg == "" {
		return ""
	}

	// Truncate long prompts for the title generator
	if len(msg) > 300 {
		msg = msg[:300]
	}

	prompt := fmt.Sprintf("Generate a short 2-4 word title in Title Case with plain spaces (examples: 'Setup Nginx SSL', 'Calculate Sphere Volume', 'Debug Memory Leak') describing the topic of the following conversation. Reply with ONLY the title itself — no quotes, no punctuation, no preamble:\n\n%s", msg)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	chatMsgs := []models.ChatMessage{
		{Role: "system", Content: "You are a concise conversation titler. Output only a 2-4 word title in Title Case with spaces."},
		{Role: "user", Content: prompt},
	}

	reply, _, err := models.CallModelWithFallback(ctx, "fast", chatMsgs)
	if err != nil || strings.TrimSpace(reply) == "" {
		return fallbackTitleFromText(msg)
	}

	clean := sanitizeSessionTitle(reply)
	if clean == "" || len(clean) < 3 {
		return fallbackTitleFromText(msg)
	}
	return clean
}

func sanitizeSessionTitle(title string) string {
	title = strings.TrimSpace(title)
	title = strings.Trim(title, "\"`'*_#[]()")
	title = invalidTitleChars.ReplaceAllString(title, "")
	fields := strings.Fields(title)
	if len(fields) > 5 {
		fields = fields[:5]
	}
	res := strings.Join(fields, " ")
	if len(res) > 35 {
		res = res[:35]
	}
	return strings.TrimSpace(res)
}

func fallbackTitleFromText(msg string) string {
	fields := strings.Fields(msg)
	if len(fields) == 0 {
		return "New Conversation"
	}
	count := len(fields)
	if count > 4 {
		count = 4
	}
	var titleWords []string
	for _, w := range fields[:count] {
		clean := invalidTitleChars.ReplaceAllString(w, "")
		if clean != "" {
			titleWords = append(titleWords, strings.Title(strings.ToLower(clean)))
		}
	}
	if len(titleWords) == 0 {
		return "New Conversation"
	}
	return strings.Join(titleWords, " ")
}
