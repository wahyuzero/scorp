package agent

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"scorp-agent/config"
	"scorp-agent/models"
	"scorp-agent/tools"
)

// ──────────────────────────────────────────────
// File & Vision Upload Handler in Agent Mode
// ──────────────────────────────────────────────

type contentPart struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *imageURL `json:"image_url,omitempty"`
}

type imageURL struct {
	URL string `json:"url"`
}

// HandleUploadInAgentMode handles file/photo uploads when agent mode is active
func HandleUploadInAgentMode(doc TGDocument) {
	chatIDStr := fmt.Sprintf("%d", doc.ChatID)
	touchSession(chatIDStr)

	// Download the file
	fileURL := fmt.Sprintf("https://api.telegram.org/bot%s/getFile?file_id=%s", config.Cfg.TelegramBotToken, doc.FileID)
	resp, err := tools.HttpShort.Get(fileURL)
	if err != nil {
		tools.SendMessage(fmt.Sprintf("❌ Error getting file: %v", err), nil)
		return
	}
	defer resp.Body.Close()

	var fileResp struct {
		OK     bool `json:"ok"`
		Result struct {
			FilePath string `json:"file_path"`
		} `json:"result"`
	}
	json.NewDecoder(resp.Body).Decode(&fileResp)

	if !fileResp.OK || fileResp.Result.FilePath == "" {
		tools.SendMessage("❌ Could not get file path", nil)
		return
	}

	downloadURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", config.Cfg.TelegramBotToken, fileResp.Result.FilePath)
	fileResp2, err := tools.HttpLong.Get(downloadURL)
	if err != nil {
		tools.SendMessage(fmt.Sprintf("❌ Error downloading: %v", err), nil)
		return
	}
	defer fileResp2.Body.Close()

	fileData, _ := io.ReadAll(fileResp2.Body)

	// Determine if it's an image
	ext := strings.ToLower(filepath.Ext(fileResp.Result.FilePath))
	isImage := ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".webp"

	if isImage {
		// Also cache image to /tmp for subsequent tool operations (send_file, inspection, etc.)
		savePath := fmt.Sprintf("/tmp/scorp_upload_%d_%s", time.Now().Unix(), filepath.Base(fileResp.Result.FilePath))
		_ = os.WriteFile(savePath, fileData, 0644)

		// Send as vision message
		b64 := base64Encode(fileData)
		mimeType := "image/jpeg"
		if ext == ".png" {
			mimeType = "image/png"
		} else if ext == ".gif" {
			mimeType = "image/gif"
		} else if ext == ".webp" {
			mimeType = "image/webp"
		}

		questionText := "Analyze this image. Describe what you see in detail."
		if doc.Caption != "" {
			questionText = doc.Caption
		}

		parts := []contentPart{
			{Type: "text", Text: questionText},
			{Type: "image_url", ImageURL: &imageURL{URL: fmt.Sprintf("data:%s;base64,%s", mimeType, b64)}},
		}

		// Build message with vision content for the immediate LLM call
		msgs := getSessionHistory(chatIDStr)
		msgs = append(msgs, AgentMessage{Role: "user", Content: parts})

		// For persistent session history, store lightweight summary without megabytes of base64
		historyNote := fmt.Sprintf("[Image uploaded: %s (%d bytes, saved to %s)]\n%s",
			filepath.Base(fileResp.Result.FilePath), len(fileData), savePath, questionText)
		appendSessionHistory(chatIDStr, AgentMessage{Role: "user", Content: historyNote})

		msgID := tools.SendMessageGetID("👁 Analyzing image with vision model...", doc.ChatID)

		// Keep only recent messages for vision context to prevent history confusion
		var recentMsgs []AgentMessage
		if len(msgs) > 6 {
			recentMsgs = msgs[len(msgs)-6:]
		} else {
			recentMsgs = msgs
		}

		// Convert to models.ChatMessage format
		chatMsgs := make([]models.ChatMessage, len(recentMsgs))
		for i, m := range recentMsgs {
			switch c := m.Content.(type) {
			case string:
				chatMsgs[i] = models.ChatMessage{Role: m.Role, Content: c}
			default:
				jsonBytes, _ := json.Marshal(c)
				chatMsgs[i] = models.ChatMessage{Role: m.Role, Content: string(jsonBytes)}
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		reply, _, err := models.CallModelWithFallback(ctx, "vision", chatMsgs)
		if err != nil {
			tools.EditMessageByID(doc.ChatID, msgID, fmt.Sprintf("❌ Vision Error: %v", err), nil)
			return
		}

		appendSessionHistory(chatIDStr, AgentMessage{Role: "assistant", Content: reply})
		sendScorpReply(doc.ChatID, msgID, reply)
	} else {
		// Non-image file: save and inform agent
		savePath := fmt.Sprintf("/tmp/scorp_upload_%d_%s", time.Now().Unix(), filepath.Base(fileResp.Result.FilePath))
		os.WriteFile(savePath, fileData, 0644)

		userMsg := fmt.Sprintf("User uploaded a file: %s (%d bytes, saved to %s)", filepath.Base(fileResp.Result.FilePath), len(fileData), savePath)
		if doc.Caption != "" {
			userMsg += "\nCaption: " + doc.Caption
		}

		RunAgentLoop(doc.ChatID, userMsg, 0)
	}
}

// base64Encode encodes bytes to base64 string
func base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}
