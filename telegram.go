package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	tgBase         string
	tgFileBase     string
	httpClient     *http.Client // existing — 15s, for Telegram API calls
	httpShort      *http.Client // 15s timeout — general API calls, web tools
	httpLong       *http.Client // 5min timeout — AI model calls, large file ops
	httpPoll       *http.Client // 35s timeout — Telegram long polling
	httpLongAI     *http.Client // 5min timeout — dedicated for AI model calls
	lastUpdateID   int64
)

// Separate transports to avoid head-of-line blocking
var (
	transportAI  = &http.Transport{MaxIdleConns: 10, MaxIdleConnsPerHost: 5, IdleConnTimeout: 90 * time.Second}
	transportWeb = &http.Transport{MaxIdleConns: 20, MaxIdleConnsPerHost: 10, IdleConnTimeout: 90 * time.Second}
	transportTG  = &http.Transport{MaxIdleConns: 10, MaxIdleConnsPerHost: 5, IdleConnTimeout: 90 * time.Second}
)

func initTelegram() {
	tgBase = "https://api.telegram.org/bot" + cfg.TelegramBotToken
	tgFileBase = "https://api.telegram.org/file/bot" + cfg.TelegramBotToken
	httpClient = &http.Client{
		Timeout:   15 * time.Second,
		Transport: transportTG,
	}
	httpShort = &http.Client{
		Timeout:   15 * time.Second,
		Transport: transportWeb,
	}
	httpLong = &http.Client{
		Timeout:   5 * time.Minute,
		Transport: transportWeb,
	}
	httpLongAI = &http.Client{
		Timeout:   5 * time.Minute,
		Transport: transportAI,
	}
	httpPoll = &http.Client{
		Timeout:   35 * time.Second,
		Transport: transportTG,
	}
}

// ──────────────────────────────────────────────
// Keyboard Layouts
// ──────────────────────────────────────────────

func mainMenuKeyboard() map[string]interface{} {
	return map[string]interface{}{
		"inline_keyboard": []interface{}{
			[]interface{}{
				map[string]string{"text": "⚡ Quick Status", "callback_data": "status"},
				map[string]string{"text": "📊 Full Report", "callback_data": "report"},
			},
			[]interface{}{
				map[string]string{"text": "🐳 Containers", "callback_data": "containers"},
				map[string]string{"text": "☁️ Coolify", "callback_data": "coolify"},
			},
			[]interface{}{
				map[string]string{"text": "🔐 Security", "callback_data": "security"},
				map[string]string{"text": "📁 Storage", "callback_data": "storage"},
			},
			[]interface{}{
				map[string]string{"text": "🌐 Network", "callback_data": "network"},
				map[string]string{"text": "📂 Files", "callback_data": "files"},
			},
			[]interface{}{
				map[string]string{"text": "🤖 Hermes", "callback_data": "hermes"},
				map[string]string{"text": "❓ Help", "callback_data": "help"},
			},
		},
	}
}

func backButtonKeyboard() map[string]interface{} {
	return map[string]interface{}{
		"inline_keyboard": []interface{}{
			[]interface{}{
				map[string]string{"text": "◀️ Back to Menu", "callback_data": "menu"},
			},
		},
	}
}

func backAndRefreshKeyboard(section string) map[string]interface{} {
	return map[string]interface{}{
		"inline_keyboard": []interface{}{
			[]interface{}{
				map[string]string{"text": "🔄 Refresh", "callback_data": section},
				map[string]string{"text": "◀️ Back to Menu", "callback_data": "menu"},
			},
		},
	}
}

// ──────────────────────────────────────────────
// Setup bot commands
// ──────────────────────────────────────────────

func setupBotCommands() {
	commands := []map[string]string{
		{"command": "start", "description": "🏠 Main menu"},
		{"command": "status", "description": "⚡ Quick status"},
		{"command": "report", "description": "📊 Full hourly report"},
		{"command": "containers", "description": "🐳 Docker containers"},
		{"command": "security", "description": "🔐 Security summary"},
		{"command": "storage", "description": "📁 Storage health"},
		{"command": "hermes", "description": "🤖 Hermes Agent status"},
		{"command": "agent", "description": "🛠 Run agent (tools)"},
		{"command": "stop", "description": "🛑 Stop chat/agent mode"},
		{"command": "cron", "description": "⏰ Scheduled tasks"},
		{"command": "model", "description": "🤖 AI model settings"},
		{"command": "usage", "description": "📊 Token usage stats"},
		{"command": "skill", "description": "🎯 Run a skill"},
		{"command": "skills", "description": "📋 List skills"},
		{"command": "clear", "description": "🧹 Clear history"},
		{"command": "help", "description": "❓ Help"},
	}
	payload := map[string]interface{}{"commands": commands}
	_, err := tgPost("/setMyCommands", payload)
	if err != nil {
		log.Printf("[telegram] set commands error: %v", err)
	} else {
		log.Println("[telegram] Bot commands registered")
	}
}

// ──────────────────────────────────────────────
// Send / Edit Messages
// ──────────────────────────────────────────────

func sendMessage(text string, keyboard map[string]interface{}) bool {
	chunks := splitMessage(text, 4000)
	ok := true
	for i, chunk := range chunks {
		payload := map[string]interface{}{
			"chat_id":                  cfg.TelegramChatID,
			"text":                     chunk,
			"parse_mode":               "HTML",
			"disable_web_page_preview": true,
		}
		if keyboard != nil && i == len(chunks)-1 {
			payload["reply_markup"] = keyboard
		}
		resp, err := tgPost("/sendMessage", payload)
		if err != nil {
			log.Printf("[telegram] send error: %v", err)
			ok = false
		} else if !resp.OK {
			log.Printf("[telegram] send failed: %s", resp.Description)
			ok = false
		}
		if len(chunks) > 1 {
			time.Sleep(300 * time.Millisecond)
		}
	}
	return ok
}

func editMessage(chatID int64, messageID int64, text string, keyboard map[string]interface{}) bool {
	return editMessageByID(chatID, messageID, text, keyboard)
}

func editMessageByID(chatID int64, messageID int64, text string, keyboard map[string]interface{}) bool {
	if len(text) > 4096 {
		text = text[:4096]
	}
	payload := map[string]interface{}{
		"chat_id":                  chatID,
		"message_id":               messageID,
		"text":                     text,
		"parse_mode":               "HTML",
		"disable_web_page_preview": true,
	}
	if keyboard != nil {
		payload["reply_markup"] = keyboard
	}
	resp, err := tgPost("/editMessageText", payload)
	if err != nil || !resp.OK {
		if resp != nil && strings.Contains(resp.Description, "message is not modified") {
			return true
		}
		return false
	}
	return true
}

// sendMessageGetID sends a message and returns the message_id (0 on failure).
func sendMessageGetID(text string, chatID int64) int64 {
	payload := map[string]interface{}{
		"chat_id":                  chatID,
		"text":                     text,
		"parse_mode":               "HTML",
		"disable_web_page_preview": true,
	}
	data, _ := json.Marshal(payload)
	resp, err := httpClient.Post(tgBase+"/sendMessage", "application/json", bytes.NewReader(data))
	if err != nil {
		log.Printf("[telegram] sendMessageGetID error: %v", err)
		return 0
	}
	defer resp.Body.Close()

	var result struct {
		OK     bool `json:"ok"`
		Result struct {
			MessageID int64 `json:"message_id"`
		} `json:"result"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	if !result.OK {
		return 0
	}
	return result.Result.MessageID
}

func answerCallback(callbackID string, text string) {
	payload := map[string]interface{}{
		"callback_query_id": callbackID,
	}
	if text != "" {
		payload["text"] = text
	}
	tgPost("/answerCallbackQuery", payload)
}

// ──────────────────────────────────────────────
// Send Document
// ──────────────────────────────────────────────

func sendDocument(chatID string, filePath string, caption string) bool {
	f, err := os.Open(filePath)
	if err != nil {
		log.Printf("[telegram] open file error: %v", err)
		return false
	}
	defer f.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("chat_id", chatID)
	if caption != "" {
		writer.WriteField("caption", caption)
	}

	part, err := writer.CreateFormFile("document", baseName(filePath))
	if err != nil {
		return false
	}
	io.Copy(part, f)
	writer.Close()

	client := httpLong
	resp, err := client.Post(tgBase+"/sendDocument", writer.FormDataContentType(), body)
	if err != nil {
		log.Printf("[telegram] sendDocument error: %v", err)
		return false
	}
	defer resp.Body.Close()

	var tgResp tgResponse
	json.NewDecoder(resp.Body).Decode(&tgResp)
	return tgResp.OK
}

func baseName(path string) string {
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}

// ──────────────────────────────────────────────
// Poll Updates
// ──────────────────────────────────────────────

type TGCommand struct {
	Text   string
	ChatID int64
	MsgID  int64
}

type TGCallback struct {
	Data   string
	ChatID int64
	MsgID  int64
	CBID   string
}

type TGDocument struct {
	FileID   string
	FileName string
	FileSize int64
	ChatID   int64
	MsgID    int64
	Caption  string
	IsPhoto  bool // True if this is a photo (not a document)
	IsVoice  bool // True if this is a voice message (needs STT)
}

func pollUpdates() ([]TGCommand, []TGCallback, []TGDocument, []TGInlineQuery) {
	var commands []TGCommand
	var callbacks []TGCallback
	var documents []TGDocument
	var inlineQueries []TGInlineQuery

	url := fmt.Sprintf("%s/getUpdates?offset=%d&timeout=30&allowed_updates=[\"message\",\"callback_query\",\"inline_query\"]",
		tgBase, lastUpdateID+1)

	client := httpPoll
	resp, err := client.Get(url)
	if err != nil {
		if !strings.Contains(err.Error(), "context deadline") {
			log.Printf("[telegram] poll error: %v", err)
			time.Sleep(1 * time.Second)
		}
		return commands, callbacks, documents, nil
	}
	defer resp.Body.Close()

	log.Printf("[telegram] poll: status=%d", resp.StatusCode)

	var result struct {
		OK     bool `json:"ok"`
		Result []struct {
			UpdateID int64 `json:"update_id"`
			Message  *struct {
				MessageID int64 `json:"message_id"`
				Chat      struct {
					ID int64 `json:"id"`
				} `json:"chat"`
				Text     string `json:"text"`
				Caption  string `json:"caption"`
				Document *struct {
					FileID   string `json:"file_id"`
					FileName string `json:"file_name"`
					FileSize int64  `json:"file_size"`
				} `json:"document"`
				Voice *struct {
					FileID   string `json:"file_id"`
					Duration int    `json:"duration"`
					FileSize int64  `json:"file_size"`
				} `json:"voice"`
				Photo []struct {
					FileID   string `json:"file_id"`
					FileSize int64  `json:"file_size"`
					Width    int    `json:"width"`
					Height   int    `json:"height"`
				} `json:"photo"`
			} `json:"message"`
			CallbackQuery *struct {
				ID      string `json:"id"`
				Data    string `json:"data"`
				Message *struct {
					MessageID int64 `json:"message_id"`
					Chat      struct {
						ID int64 `json:"id"`
					} `json:"chat"`
				} `json:"message"`
			} `json:"callback_query"`
			InlineQuery *struct {
				ID     string `json:"id"`
				From   struct {
					ID int64 `json:"id"`
				} `json:"from"`
				Query  string `json:"query"`
				Offset string `json:"offset"`
			} `json:"inline_query"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return commands, callbacks, documents, nil
	}

	log.Printf("[telegram] poll: got %d updates", len(result.Result))

	for _, u := range result.Result {
		lastUpdateID = u.UpdateID

		if u.Message != nil {
			chatID := u.Message.Chat.ID

			// Handle photo uploads (pick largest resolution = last element)
			if len(u.Message.Photo) > 0 {
				photo := u.Message.Photo[len(u.Message.Photo)-1]
				documents = append(documents, TGDocument{
					FileID:   photo.FileID,
					FileName: fmt.Sprintf("photo_%d.jpg", u.Message.MessageID),
					FileSize: photo.FileSize,
					ChatID:   chatID,
					MsgID:    u.Message.MessageID,
					Caption:  u.Message.Caption,
					IsPhoto:  true,
				})
			} else if u.Message.Voice != nil {
				// Voice message → transcribe
				voice := u.Message.Voice
				documents = append(documents, TGDocument{
					FileID:   voice.FileID,
					FileName: fmt.Sprintf("voice_%d.ogg", u.Message.MessageID),
					FileSize: voice.FileSize,
					ChatID:   chatID,
					MsgID:    u.Message.MessageID,
					Caption:  u.Message.Caption,
					IsVoice:  true,
				})
			} else if u.Message.Document != nil {
				doc := u.Message.Document
				documents = append(documents, TGDocument{
					FileID:   doc.FileID,
					FileName: doc.FileName,
					FileSize: doc.FileSize,
					ChatID:   chatID,
					MsgID:    u.Message.MessageID,
					Caption:  u.Message.Caption,
				})
			} else if u.Message.Text != "" {
				text := u.Message.Text
				// For / commands, lowercase the command part but keep args original case
				if strings.HasPrefix(text, "/") {
					// Strip @botname suffix
					text = strings.Split(text, "@")[0] + text[len(strings.Split(text, "@")[0]):]
					parts := strings.SplitN(text, " ", 2)
					cmdPart := strings.ToLower(strings.Split(parts[0], "@")[0])
					if len(parts) > 1 {
						// Keep original case for arguments (important for /ask and /agent)
						text = cmdPart + " " + parts[1]
					} else {
						text = cmdPart
					}
				}
				commands = append(commands, TGCommand{
					Text:   strings.TrimSpace(text),
					ChatID: chatID,
					MsgID:  u.Message.MessageID,
				})
			}
		}

		if u.CallbackQuery != nil {
			cb := u.CallbackQuery
			var chatID, msgID int64
			if cb.Message != nil {
				chatID = cb.Message.Chat.ID
				msgID = cb.Message.MessageID
			}
			callbacks = append(callbacks, TGCallback{
				Data:   cb.Data,
				ChatID: chatID,
				MsgID:  msgID,
				CBID:   cb.ID,
			})
		}

		if u.InlineQuery != nil {
			inlineQueries = append(inlineQueries, TGInlineQuery{
				ID:     u.InlineQuery.ID,
				Query:  u.InlineQuery.Query,
				UserID: u.InlineQuery.From.ID,
				Offset: u.InlineQuery.Offset,
			})
		}
	}

	return commands, callbacks, documents, inlineQueries
}

// ──────────────────────────────────────────────
// Download Telegram file (for uploads)
// ──────────────────────────────────────────────

func downloadTelegramFile(fileID string, destPath string) bool {
	// Step 1: getFile
	url := fmt.Sprintf("%s/getFile?file_id=%s", tgBase, fileID)
	resp, err := httpClient.Get(url)
	if err != nil {
		log.Printf("[telegram] getFile error: %v", err)
		return false
	}
	defer resp.Body.Close()

	var result struct {
		OK     bool `json:"ok"`
		Result struct {
			FilePath string `json:"file_path"`
		} `json:"result"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	if !result.OK {
		return false
	}

	// Step 2: Download
	dlURL := fmt.Sprintf("%s/%s", tgFileBase, result.Result.FilePath)
	dlResp, err := httpClient.Get(dlURL)
	if err != nil {
		return false
	}
	defer dlResp.Body.Close()

	f, err := os.Create(destPath)
	if err != nil {
		return false
	}
	defer f.Close()

	_, err = io.Copy(f, dlResp.Body)
	return err == nil
}

// ──────────────────────────────────────────────
// Helpers
// ──────────────────────────────────────────────

type tgResponse struct {
	OK          bool   `json:"ok"`
	Description string `json:"description"`
}

func tgPost(method string, payload map[string]interface{}) (*tgResponse, error) {
	data, _ := json.Marshal(payload)
	resp, err := httpClient.Post(tgBase+method, "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tgResp tgResponse
	json.NewDecoder(resp.Body).Decode(&tgResp)
	return &tgResp, nil
}

func splitMessage(text string, maxLen int) []string {
	if len(text) <= maxLen {
		return []string{text}
	}
	var chunks []string
	var current strings.Builder
	for _, line := range strings.Split(text, "\n") {
		if current.Len()+len(line)+1 > maxLen {
			if current.Len() > 0 {
				chunks = append(chunks, current.String())
				current.Reset()
			}
			current.WriteString(line)
		} else {
			if current.Len() > 0 {
				current.WriteByte('\n')
			}
			current.WriteString(line)
		}
	}
	if current.Len() > 0 {
		chunks = append(chunks, current.String())
	}
	return chunks
}

// ──────────────────────────────────────────────
// Webhook Mode (optional, replaces long polling)
// ──────────────────────────────────────────────

var webhookServer *http.Server

// setWebhook registers the webhook URL with Telegram
func setWebhook(url string) error {
	payload := map[string]interface{}{
		"url":             url,
		"allowed_updates": []string{"message", "callback_query"},
		"drop_pending_updates": true,
	}
	resp, err := tgPost("setWebhook", payload)
	if err != nil {
		return err
	}
	if !resp.OK {
		return fmt.Errorf("setWebhook failed: %s", resp.Description)
	}
	log.Printf("[telegram] Webhook set: %s", url)
	return nil
}

// deleteWebhook removes the webhook (fallback to polling)
func deleteWebhook() error {
	resp, err := tgPost("deleteWebhook", map[string]interface{}{"drop_pending_updates": true})
	if err != nil {
		return err
	}
	if !resp.OK {
		return fmt.Errorf("deleteWebhook failed: %s", resp.Description)
	}
	log.Println("[telegram] Webhook deleted")
	return nil
}

// webhookHandler handles incoming webhook updates from Telegram
func webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var update struct {
		UpdateID int64 `json:"update_id"`
		Message  *struct {
			MessageID int64 `json:"message_id"`
			Chat      struct {
				ID int64 `json:"id"`
			} `json:"chat"`
			Text     string `json:"text"`
			Caption  string `json:"caption"`
			Document *struct {
				FileID   string `json:"file_id"`
				FileName string `json:"file_name"`
				FileSize int64  `json:"file_size"`
			} `json:"document"`
			Voice *struct {
				FileID   string `json:"file_id"`
				Duration int    `json:"duration"`
				FileSize int64  `json:"file_size"`
			} `json:"voice"`
			Photo []struct {
				FileID   string `json:"file_id"`
				FileSize int64  `json:"file_size"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
			} `json:"photo"`
		} `json:"message"`
		CallbackQuery *struct {
			ID      string `json:"id"`
			Data    string `json:"data"`
			Message *struct {
				MessageID int64 `json:"message_id"`
				Chat      struct {
					ID int64 `json:"id"`
				} `json:"chat"`
			} `json:"message"`
		} `json:"callback_query"`
		InlineQuery *struct {
			ID     string `json:"id"`
			Query  string `json:"query"`
			From   struct {
				ID int64 `json:"id"`
			} `json:"from"`
			Offset string `json:"offset"`
		} `json:"inline_query"`
	}

	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		log.Printf("[telegram] Webhook decode error: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Process message
	if update.Message != nil {
		chatID := update.Message.Chat.ID
		if update.Message.Text != "" {
			handleAction(update.Message.Text, chatID, update.Message.MessageID, "")
		} else if len(update.Message.Photo) > 0 || update.Message.Document != nil {
			// Handle file/photo upload
			doc := TGDocument{
				ChatID:   chatID,
				MsgID:    update.Message.MessageID,
				Caption:  update.Message.Caption,
				IsPhoto:  len(update.Message.Photo) > 0,
			}
			if update.Message.Document != nil {
				doc.FileID = update.Message.Document.FileID
				doc.FileName = update.Message.Document.FileName
				doc.FileSize = update.Message.Document.FileSize
			} else {
				// Photo - use largest resolution
				doc.FileID = update.Message.Photo[len(update.Message.Photo)-1].FileID
				doc.FileName = fmt.Sprintf("photo_%d.jpg", update.Message.MessageID)
				doc.FileSize = update.Message.Photo[len(update.Message.Photo)-1].FileSize
				doc.IsPhoto = true
			}
			if isAgentMode(fmt.Sprintf("%d", chatID)) {
				go handleUploadInAgentMode(doc)
			} else {
				dest := fmt.Sprintf("%s/%s", uploadsDir(), doc.FileName)
				os.MkdirAll(uploadsDir(), 0755)
				ok := downloadTelegramFile(doc.FileID, dest)
				fileType := "📄"
				if doc.IsPhoto {
					fileType = "🖼"
				}
				if ok {
					sendMessage(fmt.Sprintf("✅ <b>File Saved</b>\n%s %s\n📏 %s\n📂 <code>~/uploads/%s</code>\n\n💡 <i>Gunakan /agent untuk agent mode — AI bisa analisis gambar!</i>",
						fileType, doc.FileName, humanSize(doc.FileSize), doc.FileName), mainMenuKeyboard())
				} else {
					sendMessage(fmt.Sprintf("❌ Failed to save <b>%s</b>", doc.FileName), mainMenuKeyboard())
				}
			}
		}
	}

	// Process callback query
	if update.CallbackQuery != nil {
		answerCallback(update.CallbackQuery.ID, "")
		chatID := update.CallbackQuery.Message.Chat.ID
		handleAction(update.CallbackQuery.Data, chatID, update.CallbackQuery.Message.MessageID, update.CallbackQuery.ID)
	}

	// Process inline query
	if update.InlineQuery != nil {
		handleInlineQuery(TGInlineQuery{
			ID:     update.InlineQuery.ID,
			Query:  update.InlineQuery.Query,
			UserID: update.InlineQuery.From.ID,
			Offset: update.InlineQuery.Offset,
		})
	}

	w.WriteHeader(http.StatusOK)
}

// startWebhookServer starts the HTTP server for webhook mode
func startWebhookServer(webhookURL string) error {
	// Extract path from webhook URL for HTTP handler
	// e.g., https://example.com/webhook -> /webhook
	path := "/webhook"
	if idx := strings.Index(webhookURL, "/webhook"); idx >= 0 {
		path = webhookURL[idx:]
		// Remove query params if any
		if q := strings.Index(path, "?"); q >= 0 {
			path = path[:q]
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc(path, webhookHandler)

	webhookServer = &http.Server{
		Addr:    ":8080", // Listen on port 8080
		Handler: mux,
	}

	log.Printf("[telegram] Starting webhook server on :8080%s", path)

	// Set webhook with Telegram
	if err := setWebhook(webhookURL); err != nil {
		return fmt.Errorf("set webhook: %w", err)
	}

	go func() {
		if err := webhookServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("[telegram] Webhook server error: %v", err)
		}
	}()

	return nil
}

// stopWebhookServer stops the webhook server and deletes webhook
func stopWebhookServer() {
	if webhookServer != nil {
		webhookServer.Close()
		deleteWebhook()
		log.Println("[telegram] Webhook server stopped")
	}
}
