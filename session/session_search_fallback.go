// Session Search Fallback — SQLite without FTS5
// Provides basic session storage and simple LIKE-based search
// when FTS5 is not available in the SQLite build.
// Build with: CGO_ENABLED=1 go build ... (no tags)
//
//go:build !fts5
// +build !fts5

package session

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// ──────────────────────────────────────────────
// Session Search Fallback
// ──────────────────────────────────────────────

var (
	sessionDB   *sql.DB
	sessionDBMu sync.Mutex
)

// initSessionDB opens the SQLite database and creates tables (no FTS5).
func InitSessionDB() {
	dbPath := scorpPath("sessions.db")

	var err error
	sessionDB, err = sql.Open("sqlite3", dbPath+"?_journal_mode=WAL&_busy_timeout=5000&_foreign_keys=on")
	if err != nil {
		log.Printf("[session-search] Failed to open DB: %v", err)
		return
	}

	// Optimize for low-contention
	sessionDB.SetMaxOpenConns(5)

	// Create tables (no FTS5 virtual table)
	schema := `
	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		chat_id TEXT NOT NULL,
		role TEXT NOT NULL,
		content TEXT NOT NULL,
		timestamp INTEGER NOT NULL,
		msg_index INTEGER NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_messages_chat ON messages(chat_id);
	CREATE INDEX IF NOT EXISTS idx_messages_time ON messages(timestamp DESC);
	`

	if _, err := sessionDB.Exec(schema); err != nil {
		log.Printf("[session-search] Failed to create schema: %v", err)
		return
	}

	log.Printf("[session-search] SQLite DB initialized at %s (no FTS5)", dbPath)
}

// indexMessage stores a message in the session DB for later search.
// Called with 3 args from chat.go: IndexMessage(chatID, role, content)
func IndexMessage(chatID, role, content string, msgIndex ...int) {
	if sessionDB == nil {
		return
	}
	idx := 0
	if len(msgIndex) > 0 {
		idx = msgIndex[0]
	}
	sessionDBMu.Lock()
	defer sessionDBMu.Unlock()

	_, err := sessionDB.Exec(
		`INSERT INTO messages (chat_id, role, content, timestamp, msg_index) VALUES (?, ?, ?, ?, ?)`,
		chatID, role, content, time.Now().Unix(), idx,
	)
	if err != nil {
		log.Printf("[session-search] Failed to index message: %v", err)
	}
}

// ──────────────────────────────────────────────
// Helper functions (used by session_search tool)
// ──────────────────────────────────────────────

// searchSessions performs simple LIKE-based search across messages.
// Returns matching messages with context.
func SearchSessions(query string, limit int) []SessionResult {
	if sessionDB == nil || query == "" {
		return nil
	}

	sessionDBMu.Lock()
	defer sessionDBMu.Unlock()

	// Simple LIKE search
	likeQuery := "%" + strings.ReplaceAll(query, " ", "%") + "%"
	rows, err := sessionDB.Query(`
		SELECT id, chat_id, role, content, timestamp, msg_index
		FROM messages
		WHERE content LIKE ?
		ORDER BY timestamp DESC
		LIMIT ?
	`, likeQuery, limit)
	if err != nil {
		log.Printf("[session-search] Search error: %v", err)
		return nil
	}
	defer rows.Close()

	var results []SessionResult
	for rows.Next() {
		var r SessionResult
		if err := rows.Scan(&r.ID, &r.ChatID, &r.Role, &r.Content, &r.Timestamp, &r.MsgIndex); err != nil {
			continue
		}
		results = append(results, r)
	}
	return results
}

// SessionResult represents a search result.
type SessionResult struct {
	ID        int64
	ChatID    string
	Role      string
	Content   string
	Timestamp int64
	MsgIndex  int
}

// executeSessionSearch is the tool handler for session_search.
func ExecuteSessionSearch(args map[string]interface{}) (string, bool) {
	query := getStringArg(args, "query", "")
	if query == "" {
		return "Error: 'query' is required", false
	}
	scope := getStringArg(args, "scope", "all")
	limit := getIntArg(args, "limit", 10)
	if limit > 50 {
		limit = 50
	}

	results := SearchSessions(query, limit)
	if len(results) == 0 {
		return "No matching sessions found.", true
	}

	var out strings.Builder
	out.WriteString(fmt.Sprintf("Found %d result(s) for \"%s\":\n\n", len(results), query))
	for i, r := range results {
		out.WriteString(fmt.Sprintf("%d. [%s] %s: %s\n", i+1, r.Role, time.Unix(r.Timestamp, 0).Format("Jan 2 15:04"), truncateString(r.Content, 200)))
		if scope == "all" {
			out.WriteString(fmt.Sprintf("   (chat: %s)\n", r.ChatID))
		}
		out.WriteString("\n")
	}
	return out.String(), true
}

// SessionStats returns database statistics.
func SessionStats() string {
	if sessionDB == nil {
		return "Session DB not initialized"
	}
	sessionDBMu.Lock()
	defer sessionDBMu.Unlock()

	var count int
	err := sessionDB.QueryRow(`SELECT COUNT(*) FROM messages`).Scan(&count)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	return fmt.Sprintf("Messages in session DB: %d", count)
}

// cleanupOldSessions removes messages older than retention period.
func CleanupOldSessions(retentionDays int) {
	if sessionDB == nil {
		return
	}
	sessionDBMu.Lock()
	defer sessionDBMu.Unlock()

	cutoff := time.Now().AddDate(0, 0, -retentionDays).Unix()
	_, err := sessionDB.Exec(`DELETE FROM messages WHERE timestamp < ?`, cutoff)
	if err != nil {
		log.Printf("[session-search] Cleanup error: %v", err)
	}
}