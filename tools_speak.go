package main

import (
	"fmt"
	"log"
	"os"
)

// ──────────────────────────────────────────────
// speak tool — TTS as agent tool
// ──────────────────────────────────────────────

func executeSpeak(args map[string]interface{}, chatID int64) (string, bool) {
	text, _ := args["text"].(string)
	if text == "" {
		return "Missing 'text' parameter", false
	}

	voice, _ := args["voice"].(string)

	// Synthesize
	oggPath, err := synthesizeSpeech(text, voice)
	if err != nil {
		log.Printf("[speak] TTS error: %v", err)
		return fmt.Sprintf("TTS error: %v", err), false
	}
	defer os.Remove(oggPath)

	// Send as voice message
	if err := sendVoiceMessage(chatID, oggPath); err != nil {
		log.Printf("[speak] send error: %v", err)
		return fmt.Sprintf("Failed to send voice: %v", err), false
	}

	truncMsg := text
	if len(truncMsg) > 80 {
		truncMsg = truncMsg[:77] + "..."
	}
	return fmt.Sprintf("✅ Voice sent: \"%s\" (voice=%s)", truncMsg, voice), true
}

func init() {
	registerTool(ToolDef{
		Name:        "speak",
		Description: "Convert text to speech and send as a voice message. Use for voice replies or audio output.",
		Category:    "other",
		Native:      true,
		Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
			return executeSpeak(args, chatID)
		},
		Arguments: map[string]ArgDef{
			"text":  {Type: "string", Description: "Text to convert to speech", Required: true},
			"voice": {Type: "string", Description: "Voice name (default: id-ID-ArdiNeural). Examples: id-ID-GadisNeural, en-US-JennyNeural"},
		},
	})
}
