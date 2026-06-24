package telegram

import (
	"scorp-agent/agent"
	"scorp-agent/config"
	"scorp-agent/internal/helpers"
	"scorp-agent/tools"
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


// HandleAction callback — set from main.go to avoid import cycle
var HandleAction func(action string, chatID int64, messageID int64, callbackID string)

var (
	TgBase         string
	TgFileBase     string
	HttpClient     *http.Client // existing — 15s, for Telegram API calls
	HttpShort      *http.Client // 15s timeout — general API calls, web tools
	HttpLong       *http.Client // 5min timeout — AI model calls, large file ops
	HttpPoll       *http.Client // 35s timeout — Telegram long polling
	HttpLongAI     *http.Client // 5min timeout — dedicated for AI model calls
	LastUpdateID   int64
)

// Separate transports to avoid head-of-line blocking
var (
	transportAI  = &http.Transport{MaxIdleConns: 10, MaxIdleConnsPerHost: 5, IdleConnTimeout: 90 * time.Second}
	transportWeb = &http.Transport{MaxIdleConns: 20, MaxIdleConnsPerHost: 10, IdleConnTimeout: 90 * time.Second}
	transportTG  = &http.Transport{MaxIdleConns: 10, MaxIdleConnsPerHost: 5, IdleConnTimeout: 90 * time.Second}
)

func InitTelegram() {
	TgBase = "https://api.telegram.org/bot" + config.Cfg.TelegramBotToken
	TgFileBase = "https://api.telegram.org/file/bot" + config.Cfg.TelegramBotToken
	HttpClient = &http.Client{
		Timeout:   15 * time.Second,
		Transport: transportTG,
	}
	HttpShort = &http.Client{
		Timeout:   15 * time.Second,
		Transport: transportWeb,
	}
	HttpLong = &http.Client{
		Timeout:   5 * time.Minute,
		Transport: transportWeb,
	}
	HttpLongAI = &http.Client{
		Timeout:   5 * time.Minute,
		Transport: transportAI,
	}
	HttpPoll = &http.Client{
		Timeout:   35 * time.Second,
		Transport: transportTG,
	}
}

// ──────────────────────────────────────────────
// Keyboard Layouts
// ──────────────────────────────────────────────

func MainMenuKeyboard() map[string]interface{} {
	return map[string]interface{}{
		"inline_keyboard": []interface{}{
			[]interface{}{
				map[string]string{"text": "📊 Monitor", "callback_data": "mn:mon"},
				map[string]string{"text": "🤖 Models", "callback_data": "/model"},
			},
			[]interface{}{
				map[string]string{"text": "🔧 System", "callback_data": "mn:sys"},
				map[string]string{"text": "⚙️ Settings", "callback_data": "mn:set"},
			},
			[]interface{}{
				map[string]string{"text": "❓ Help", "callback_data": "help"},
			},
		},
	}
}

func MonitorMenuKeyboard() map[string]interface{} {
	return map[string]interface{}{
		"inline_keyboard": []interface{}{
			[]interface{}{
				map[string]string{"text": "⚡ Status", "callback_data": "status"},
				map[string]string{"text": "📊 Report", "callback_data": "report"},
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
				map[string]string{"text": "◀️ Menu", "callback_data": "mn:main"},
			},
		},
	}
}

func SystemMenuKeyboard() map[string]interface{} {
	return map[string]interface{}{
		"inline_keyboard": []interface{}{
			[]interface{}{
				map[string]string{"text": "📂 Files", "callback_data": "files"},
				map[string]string{"text": "📊 Usage", "callback_data": "/usage"},
			},
			[]interface{}{
				map[string]string{"text": "🔌 MCP", "callback_data": "/mcp"},
				map[string]string{"text": "⏰ Cron", "callback_data": "/cron"},
			},
			[]interface{}{
				map[string]string{"text": "📋 Sessions", "callback_data": "/sessions"},
				map[string]string{"text": "◀️ Menu", "callback_data": "mn:main"},
			},
		},
	}
}

func BackButtonKeyboard() map[string]interface{} {
	return map[string]interface{}{
		"inline_keyboard": []interface{}{
			[]interface{}{
				map[string]string{"text": "◀️ Back to Menu", "callback_data": "mn:main"},
			},
		},
	}
}

// ─── Settings Menu ───

func SettingsMenuText() string {
	monStatus := "OFF"
	if config.Cfg.MonitoringEnabled {
		monStatus = "ON"
	}
	secStatus := "OFF"
	if config.Cfg.SecurityAlertsEnabled {
		secStatus = "ON"
	}
	repStatus := "OFF"
	if config.Cfg.ScheduledReportsEnabled {
		repStatus = "ON"
	}
	return "⚙️ <b>Settings</b>\n\n" +
		"<i>Toggle monitoring modules on/off:</i>\n" +
		"━━━━━━━━━━━━━━━━━━━\n" +
		fmt.Sprintf("📊 <b>Resource Monitoring</b>: %s\n", monStatus) +
		fmt.Sprintf("🔐 <b>Security Alerts</b>: %s\n", secStatus) +
		fmt.Sprintf("📅 <b>Scheduled Reports</b>: %s", repStatus)
}

func SettingsMenuKeyboard() map[string]interface{} {
	monLabel := "📊 Monitoring: OFF"
	if config.Cfg.MonitoringEnabled {
		monLabel = "📊 Monitoring: ON"
	}
	secLabel := "🔐 Security: OFF"
	if config.Cfg.SecurityAlertsEnabled {
		secLabel = "🔐 Security: ON"
	}
	repLabel := "📅 Reports: OFF"
	if config.Cfg.ScheduledReportsEnabled {
		repLabel = "📅 Reports: ON"
	}
	return map[string]interface{}{
		"inline_keyboard": []interface{}{
			[]interface{}{
				map[string]string{"text": monLabel, "callback_data": "set:mon"},
			},
			[]interface{}{
				map[string]string{"text": secLabel, "callback_data": "set:sec"},
			},
			[]interface{}{
				map[string]string{"text": repLabel, "callback_data": "set:rep"},
			},
			[]interface{}{
				map[string]string{"text": "◀️ Menu", "callback_data": "mn:main"},
			},
		},
	}
}

func BackAndRefreshKeyboard(section string) map[string]interface{} {
	return map[string]interface{}{
		"inline_keyboard": []interface{}{
			[]interface{}{
				map[string]string{"text": "🔄 Refresh", "callback_data": section},
				map[string]string{"text": "◀️ Back to Menu", "callback_data": "mn:main"},
			},
		},
	}
}

// ──────────────────────────────────────────────
// Setup bot commands
// ──────────────────────────────────────────────

func SetupBotCommands() {
	commands := []map[string]string{
		{"command": "start", "description": "🏠 Main menu"},
		{"command": "agent", "description": "🛠 Agent mode (tools)"},
		{"command": "stop", "description": "🛑 Stop chat/agent"},
		{"command": "clear", "description": "🧹 Clear history"},
		{"command": "model", "description": "🤖 Model manager"},
		{"command": "help", "description": "❓ Help"},
	}
	payload := map[string]interface{}{"commands": commands}
	_, err := TgPost("/setMyCommands", payload)
	if err != nil {
		log.Printf("[telegram] set commands error: %v", err)
	} else {
		log.Println("[telegram] Bot commands registered")
	}
}

// ──────────────────────────────────────────────
// Send / Edit Messages
// ──────────────────────────────────────────────

func SendMessage(text string, keyboard map[string]interface{}) bool {
	chunks := helpers.SplitMessage(text, 4000)
	ok := true
	for i, chunk := range chunks {
		payload := map[string]interface{}{
			"chat_id":                  config.Cfg.TelegramChatID,
			"text":                     chunk,
			"parse_mode":               "HTML",
			"disable_web_page_preview": true,
		}
		if keyboard != nil && i == len(chunks)-1 {
			payload["reply_markup"] = keyboard
		}
		resp, err := TgPost("/sendMessage", payload)
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

func EditMessage(chatID int64, messageID int64, text string, keyboard map[string]interface{}) bool {
	return EditMessageByID(chatID, messageID, text, keyboard)
}

func EditMessageByID(chatID int64, messageID int64, text string, keyboard map[string]interface{}) bool {
	if len([]rune(text)) > 4096 {
		text = string([]rune(text)[:4096])
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
	resp, err := TgPost("/editMessageText", payload)
	if err != nil || !resp.OK {
		if resp != nil && strings.Contains(resp.Description, "message is not modified") {
			return true
		}
		return false
	}
	return true
}

// sendMessageGetID sends a message and returns the message_id (0 on failure).
func SendMessageGetID(text string, chatID int64) int64 {
	payload := map[string]interface{}{
		"chat_id":                  chatID,
		"text":                     text,
		"parse_mode":               "HTML",
		"disable_web_page_preview": true,
	}
	data, _ := json.Marshal(payload)
	resp, err := HttpClient.Post(TgBase+"/sendMessage", "application/json", bytes.NewReader(data))
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

// sendChatAction sends a chat action indicator (e.g. "typing") to Telegram.
// The indicator expires after ~5 seconds, so it must be sent periodically.
func SendChatAction(chatID int64, action string) {
	payload := map[string]interface{}{
		"chat_id": chatID,
		"action":  action,
	}
	TgPost("/sendChatAction", payload)
}

func AnswerCallback(callbackID string, text string) {
	payload := map[string]interface{}{
		"callback_query_id": callbackID,
	}
	if text != "" {
		payload["text"] = text
	}
	TgPost("/answerCallbackQuery", payload)
}

// ──────────────────────────────────────────────
// Send Document
// ──────────────────────────────────────────────

func SendDocument(chatID string, filePath string, caption string) bool {
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

	client := HttpLong
	resp, err := client.Post(TgBase+"/sendDocument", writer.FormDataContentType(), body)
	if err != nil {
		log.Printf("[telegram] sendDocument error: %v", err)
		return false
	}
	defer resp.Body.Close()

	var tgResp TgResponse
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

// TGDocument moved to agent package

func PollUpdates() ([]TGCommand, []TGCallback, []agent.TGDocument, []tools.TGInlineQuery) {
	var commands []TGCommand
	var callbacks []TGCallback
	var documents []agent.TGDocument
	var inlineQueries []tools.TGInlineQuery

	url := fmt.Sprintf("%s/getUpdates?offset=%d&timeout=30&allowed_updates=[\"message\",\"callback_query\",\"inline_query\"]",
		TgBase, LastUpdateID+1)

	client := HttpPoll
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
		LastUpdateID = u.UpdateID

		if u.Message != nil {
			chatID := u.Message.Chat.ID

			// Handle photo uploads (pick largest resolution = last element)
			if len(u.Message.Photo) > 0 {
				photo := u.Message.Photo[len(u.Message.Photo)-1]
				documents = append(documents, agent.TGDocument{
					FileID:   photo.FileID,
					FileName: fmt.Sprintf("photo_%d.jpg", u.Message.MessageID),
					FileSize: photo.FileSize,
					ChatID:   chatID,
					MsgID:    u.Message.MessageID,
					Caption:  u.Message.Caption,
					IsPhoto:  true,
				})
			} else if u.Message.Voice != nil {
				voice := u.Message.Voice
				documents = append(documents, agent.TGDocument{
					FileID:   voice.FileID,
					FileName: fmt.Sprintf("voice_%d.ogg", u.Message.MessageID),
					FileSize: voice.FileSize,
					ChatID:   chatID,
					MsgID:    u.Message.MessageID,
					Caption:  u.Message.Caption,
				})
			} else if u.Message.Document != nil {
				doc := u.Message.Document
				documents = append(documents, agent.TGDocument{
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
				inlineQueries = append(inlineQueries, tools.TGInlineQuery{
					ID:     u.InlineQuery.ID,
					Query:  u.InlineQuery.Query,
					UserID: u.InlineQuery.From.ID,
					Offset: u.InlineQuery.Offset,
				})
			}
	}

	return commands, callbacks, documents, inlineQueries
	}

	// Helpers
	// ──────────────────────────────────────────────

	type TgResponse struct {
		OK          bool   `json:"ok"`
		Description string `json:"description"`
	}

	func TgPost(method string, payload map[string]interface{}) (*TgResponse, error) {
	data, _ := json.Marshal(payload)
	resp, err := HttpClient.Post(TgBase+method, "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tgResp TgResponse
	json.NewDecoder(resp.Body).Decode(&tgResp)
	return &tgResp, nil
}

// splitMessage moved to helpers.SplitMessage

// ──────────────────────────────────────────────
// Webhook Mode (optional, replaces long polling)
// ──────────────────────────────────────────────

var webhookServer *http.Server

// setWebhook registers the webhook URL with Telegram
func SetWebhook(url string) error {
	payload := map[string]interface{}{
		"url":             url,
		"allowed_updates": []string{"message", "callback_query"},
		"drop_pending_updates": true,
	}
	resp, err := TgPost("setWebhook", payload)
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
func DeleteWebhook() error {
	resp, err := TgPost("deleteWebhook", map[string]interface{}{"drop_pending_updates": true})
	if err != nil {
		return err
	}
	if !resp.OK {
		return fmt.Errorf("deleteWebhook failed: %s", resp.Description)
	}
	log.Println("[telegram] Webhook deleted")
	return nil
}

// WebhookHandler handles incoming webhook updates from Telegram
func WebhookHandler(w http.ResponseWriter, r *http.Request) {
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
			HandleAction(update.Message.Text, chatID, update.Message.MessageID, "")
		} else if len(update.Message.Photo) > 0 || update.Message.Document != nil {
			// Handle file/photo upload
			doc := agent.TGDocument{
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
			// All uploads go through agent
			go agent.HandleUploadInAgentMode(doc)
		}
	}

	// Process callback query
	if update.CallbackQuery != nil {
		AnswerCallback(update.CallbackQuery.ID, "")
		chatID := update.CallbackQuery.Message.Chat.ID
		HandleAction(update.CallbackQuery.Data, chatID, update.CallbackQuery.Message.MessageID, update.CallbackQuery.ID)
	}

	// Process inline query
	if update.InlineQuery != nil {
		tools.HandleInlineQuery(tools.TGInlineQuery{
			ID:     update.InlineQuery.ID,
			Query:  update.InlineQuery.Query,
			UserID: update.InlineQuery.From.ID,
			Offset: update.InlineQuery.Offset,
		})
	}

	w.WriteHeader(http.StatusOK)
}

// startWebhookServer starts the HTTP server for webhook mode
func StartWebhookServer(webhookURL string) error {
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
	mux.HandleFunc(path, WebhookHandler)

	webhookServer = &http.Server{
		Addr:    ":8080", // Listen on port 8080
		Handler: mux,
	}

	log.Printf("[telegram] Starting webhook server on :8080%s", path)

	// Set webhook with Telegram
	if err := SetWebhook(webhookURL); err != nil {
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
func StopWebhookServer() {
	if webhookServer != nil {
		webhookServer.Close()
		DeleteWebhook()
		log.Println("[telegram] Webhook server stopped")
	}
}
