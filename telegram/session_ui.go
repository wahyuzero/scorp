package telegram

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"scorp-agent/agent"
	"scorp-agent/config"
	"scorp-agent/internal/helpers"
	"scorp-agent/tools"
)

// ──────────────────────────────────────────────
// Telegram Session Switcher & State Manager
// Maps each Telegram user/chat to an active conversation session key.
// ──────────────────────────────────────────────

var (
	tgSessionMu      sync.RWMutex
	tgActiveSessions = make(map[string]string) // chatIDStr -> sessionID
)

func init() {
	loadTgSessionMapping()
	tools.OnSessionAutoTitled = func(oldID, newID, chatIDStr string) {
		tgSessionMu.Lock()
		defer tgSessionMu.Unlock()
		for cid, s := range tgActiveSessions {
			if s == oldID {
				tgActiveSessions[cid] = newID
			}
		}
		if tgActiveSessions[chatIDStr] == "" || tgActiveSessions[chatIDStr] == oldID {
			tgActiveSessions[chatIDStr] = newID
		}
		saveTgSessionMapping()
	}
}

func tgSessionMappingPath() string {
	return filepath.Join(config.ScorpDir(), "telegram_sessions.json")
}

func loadTgSessionMapping() {
	tgSessionMu.Lock()
	defer tgSessionMu.Unlock()

	data, err := os.ReadFile(tgSessionMappingPath())
	if err != nil {
		return
	}
	_ = json.Unmarshal(data, &tgActiveSessions)
}

func saveTgSessionMapping() {
	data, err := json.MarshalIndent(tgActiveSessions, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(tgSessionMappingPath(), data, 0644)
}

// GetActiveSessionID returns current active session for a telegram chat, defaulting to chatIDStr
func GetActiveSessionID(chatIDStr string) string {
	tgSessionMu.RLock()
	sess, ok := tgActiveSessions[chatIDStr]
	tgSessionMu.RUnlock()

	if ok && sess != "" {
		return sess
	}
	return chatIDStr
}

// SetActiveSessionID changes current active session for a telegram chat
func SetActiveSessionID(chatIDStr, sessionID string) {
	tgSessionMu.Lock()
	tgActiveSessions[chatIDStr] = sessionID
	saveTgSessionMapping()
	tgSessionMu.Unlock()
}

// FormatSessionMenuText builds the HTML status text for the session menu
func FormatSessionMenuText(chatIDStr string) string {
	activeSess := GetActiveSessionID(chatIDStr)
	sessions := agent.ListSessions()

	var sb strings.Builder
	sb.WriteString("📂 <b>Chat Sessions Manager</b>\n\n")
	sb.WriteString(fmt.Sprintf("Active Session: 🟢 <b>%s</b>\n", activeSess))
	sb.WriteString("━━━━━━━━━━━━━━━━━━━━\n")

	if len(sessions) == 0 {
		sb.WriteString("<i>No saved sessions yet. Type <code>/session new &lt;name&gt;</code> to create one.</i>\n")
		return sb.String()
	}

	foundActive := false
	for _, s := range sessions {
		marker := "•"
		if s.ID == activeSess {
			marker = "👉 <b>[ACTIVE]</b>"
			foundActive = true
		}
		preview := helpers.EscapeHTML(s.LastPreview)
		if len(preview) > 30 {
			preview = preview[:27] + "..."
		}
		if preview == "" {
			preview = "<i>(empty)</i>"
		}
		sb.WriteString(fmt.Sprintf("%s <code>%s</code> — %d msgs (%s)\n  ↳ %s\n",
			marker, s.ID, s.MsgCount, s.LastModified.Format("02 Jan 15:04"), preview))
	}

	if !foundActive {
		sb.WriteString(fmt.Sprintf("👉 <b>[ACTIVE]</b> <code>%s</code> — <i>(new active session)</i>\n", activeSess))
	}

	sb.WriteString("\n<b>Usage:</b>\n")
	sb.WriteString("• Click buttons below to switch sessions\n")
	sb.WriteString("• <code>/session new &lt;name&gt;</code> — create &amp; switch\n")
	sb.WriteString("• <code>/session use &lt;name&gt;</code> — switch session\n")
	sb.WriteString("• <code>/session delete &lt;name&gt;</code> — delete session\n")
	sb.WriteString("• <code>/session rename &lt;old&gt; &lt;new&gt;</code>\n")

	return sb.String()
}

// BuildSessionMenuKeyboard generates inline buttons for switching & managing sessions
func BuildSessionMenuKeyboard(chatIDStr string) map[string]interface{} {
	activeSess := GetActiveSessionID(chatIDStr)
	sessions := agent.ListSessions()

	var rows [][]map[string]string

	// Rows of switch buttons (max 6 sessions shown for clean keyboard)
	count := 0
	for _, s := range sessions {
		if count >= 6 {
			break
		}
		label := s.ID
		if s.ID == activeSess {
			label = "🟢 " + s.ID
		} else {
			label = "📁 " + s.ID
		}
		rows = append(rows, []map[string]string{
			{
				"text":          label,
				"callback_data": fmt.Sprintf("sess:use:%s", s.ID),
			},
		})
		count++
	}

	// Action row
	rows = append(rows, []map[string]string{
		{"text": "➕ Sesi Baru", "callback_data": "sess:new"},
		{"text": "🗜️ Compact", "callback_data": "sess:compact"},
		{"text": "🔄 Refresh", "callback_data": "sess:refresh"},
		{"text": "🧹 Reset", "callback_data": "sess:clear"},
	})

	// Navigation row
	rows = append(rows, []map[string]string{
		{"text": "◀️ Back to Menu", "callback_data": "mn:main"},
	})

	return map[string]interface{}{
		"inline_keyboard": rows,
	}
}

// HandleSessionCallback handles inline button callbacks with sess:* prefix
func HandleSessionCallback(action, chatIDStr string) (string, map[string]interface{}, bool) {
	if !strings.HasPrefix(action, "sess:") {
		return "", nil, false
	}

	parts := strings.Split(action, ":")
	if len(parts) < 2 {
		return "", nil, false
	}

	cmd := parts[1]
	switch cmd {
	case "refresh":
		return FormatSessionMenuText(chatIDStr), BuildSessionMenuKeyboard(chatIDStr), true

	case "new":
		newSess := fmt.Sprintf("chat-%s", time.Now().Format("0102-150405"))
		SetActiveSessionID(chatIDStr, newSess)
		return fmt.Sprintf("✓ Started new conversation: 🟢 <b>%s</b>\n<i>Topik sesi ini akan otomatis dinamai sesuai pesan pertamamu.</i>\n\n%s", newSess, FormatSessionMenuText(chatIDStr)),
			BuildSessionMenuKeyboard(chatIDStr), true

	case "use":
		if len(parts) >= 3 {
			targetSess := parts[2]
			SetActiveSessionID(chatIDStr, targetSess)
			return FormatSessionMenuText(chatIDStr), BuildSessionMenuKeyboard(chatIDStr), true
		}

	case "compact":
		activeSess := GetActiveSessionID(chatIDStr)
		stats := agent.CompactSessionHistory(activeSess)
		return agent.FormatCompactStats(activeSess, stats), BuildSessionMenuKeyboard(chatIDStr), true

	case "clear":
		activeSess := GetActiveSessionID(chatIDStr)
		agent.ClearChatSession(activeSess)
		return fmt.Sprintf("🧹 History cleared for session <code>%s</code>.\n\n%s", activeSess, FormatSessionMenuText(chatIDStr)),
			BuildSessionMenuKeyboard(chatIDStr), true
	}

	return "", nil, false
}
