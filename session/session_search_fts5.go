// Session Search FTS5 — SQLite with FTS5 full-text search
// Provides ranked full-text search with snippets when FTS5 is available.
// Build with: CGO_ENABLED=1 go build -tags fts5 ...
//
//go:build fts5
// +build fts5

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
// Session Search FTS5
// ──────────────────────────────────────────────

var (
	sessionDB   *sql.DB
	sessionDBMu sync.Mutex
)

// initSessionDB opens the SQLite database and creates FTS5 virtual table.
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

	// Create FTS5 virtual table with content table for storage
	schema := `
CREATE TABLE IF NOT EXISTS messages_content (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	chat_id TEXT NOT NULL,
	role TEXT NOT NULL,
	content TEXT NOT NULL,
	timestamp INTEGER NOT NULL,
	msg_index INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_messages_content_chat ON messages_content(chat_id);
CREATE INDEX IF NOT EXISTS idx_messages_content_time ON messages_content(timestamp DESC);

CREATE VIRTUAL TABLE IF NOT EXISTS messages_fts USING fts5(
	chat_id, role, content,
	content='messages_content',
	content_rowid='id'
);

-- Triggers to keep FTS in sync with content table
CREATE TRIGGER IF NOT EXISTS messages_ai AFTER INSERT ON messages_content BEGIN
	INSERT INTO messages_fts(rowid, chat_id, role, content) VALUES (new.id, new.chat_id, new.role, new.content);
END;

CREATE TRIGGER IF NOT EXISTS messages_ad AFTER DELETE ON messages_content BEGIN
	INSERT INTO messages_fts(messages_fts, rowid, chat_id, role, content) VALUES ('delete', old.id, old.chat_id, old.role, old.content);
END;

CREATE TRIGGER IF NOT EXISTS messages_au AFTER UPDATE ON messages_content BEGIN
	INSERT INTO messages_fts(messages_fts, rowid, chat_id, role, content) VALUES ('delete', old.id, old.chat_id, old.role, old.content);
	INSERT INTO messages_fts(rowid, chat_id, role, content) VALUES (new.id, new.chat_id, new.role, new.content);
END;
`

	if _, err := sessionDB.Exec(schema); err != nil {
		log.Printf("[session-search] Failed to create schema: %v", err)
		return
	}

	log.Printf("[session-search] SQLite DB initialized at %s (FTS5 ENABLED)", dbPath)
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
		`INSERT INTO messages_content (chat_id, role, content, timestamp, msg_index) VALUES (?, ?, ?, ?, ?)`,
		chatID, role, content, time.Now().Unix(), idx,
	)
	if err != nil {
		log.Printf("[session-search] Failed to index message: %v", err)
	}
}

// ──────────────────────────────────────────────
// Helper functions (used by session_search tool)
// ──────────────────────────────────────────────

// searchSessions performs FTS5 full-text search across messages.
// Returns ranked results with snippets.
func SearchSessions(query string, limit int) []SessionResult {
	if sessionDB == nil || query == "" {
		return nil
	}

	sessionDBMu.Lock()
	defer sessionDBMu.Unlock()

	// FTS5 query - supports phrase search, prefix, NEAR, etc.
	// Escape special FTS5 chars
	ftsQuery := EscapeFTS5Query(query)

	rows, err := sessionDB.Query(`
		SELECT m.id, m.chat_id, m.role, m.content, m.timestamp, m.msg_index,
		       snippet(messages_fts, 2, '<b>', '</b>', '...', 64) as snippet
		FROM messages_fts
		JOIN messages_content m ON messages_fts.rowid = m.id
		WHERE messages_fts MATCH ?
		ORDER BY rank
		LIMIT ?
	`, ftsQuery, limit)
	if err != nil {
		log.Printf("[session-search] FTS5 search error: %v", err)
		// Fallback to LIKE if FTS5 fails
		return searchSessionsFallback(query, limit)
	}
	defer rows.Close()

	var results []SessionResult
	for rows.Next() {
		var r SessionResult
		var snippet sql.NullString
		if err := rows.Scan(&r.ID, &r.ChatID, &r.Role, &r.Content, &r.Timestamp, &r.MsgIndex, &snippet); err != nil {
			continue
		}
		if snippet.Valid {
			r.Snippet = snippet.String
		}
		results = append(results, r)
	}
	return results
}

// searchSessionsFallback is used when FTS5 query fails
func searchSessionsFallback(query string, limit int) []SessionResult {
	if sessionDB == nil {
		return nil
	}

	likeQuery := "%" + strings.ReplaceAll(query, " ", "%") + "%"
	rows, err := sessionDB.Query(`
		SELECT id, chat_id, role, content, timestamp, msg_index
		FROM messages_content
		WHERE content LIKE ?
		ORDER BY timestamp DESC
		LIMIT ?
	`, likeQuery, limit)
	if err != nil {
		log.Printf("[session-search] Fallback search error: %v", err)
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

// escapeFTS5Query escapes special FTS5 characters
func EscapeFTS5Query(query string) string {
	// FTS5 special chars: " * + - ( ) { } ^ ~ : < >
	replacer := strings.NewReplacer(
		`"`, `""`,
		`*`, ``,
		`+`, ` `,
		`-`, ` `,
		`(`, ` `,
		`)`, ` `,
		`{`, ` `,
		`}`, ` `,
		`^`, ` `,
		`~`, ` `,
		`:`, ` `,
		`<`, ` `,
		`>`, ` `,
	)
	return replacer.Replace(query)
}

// SessionResult represents a search result.
type SessionResult struct {
	ID        int64
	ChatID    string
	Role      string
	Content   string
	Timestamp int64
	MsgIndex  int
	Snippet   string
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
		displayContent := r.Snippet
		if displayContent == "" {
			displayContent = truncateString(r.Content, 200)
		}
		out.WriteString(fmt.Sprintf("%d. [%s] %s: %s\n", i+1, r.Role, time.Unix(r.Timestamp, 0).Format("Jan 2 15:04"), displayContent))
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
	err := sessionDB.QueryRow(`SELECT COUNT(*) FROM messages_content`).Scan(&count)
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
	_, err := sessionDB.Exec(`DELETE FROM messages_content WHERE timestamp < ?`, cutoff)
	if err != nil {
		log.Printf("[session-search] Cleanup error: %v", err)
	}
}

// migrateFromFallback migrates data from old non-FTS5 schema to FTS5 schema
func MigrateFromFallback() {
	if sessionDB == nil {
		return
	}
	sessionDBMu.Lock()
	defer sessionDBMu.Unlock()

	// Check if old 'messages' table exists (from fallback)
	var tableExists int
	err := sessionDB.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='messages'`).Scan(&tableExists)
	if err != nil || tableExists == 0 {
		return // No old table to migrate
	}

	log.Printf("[session-search] Migrating from fallback schema to FTS5...")

	// Copy data from old messages table to new messages_content
	_, err = sessionDB.Exec(`
		INSERT OR IGNORE INTO messages_content (id, chat_id, role, content, timestamp, msg_index)
		SELECT id, chat_id, role, content, timestamp, msg_index FROM messages
	`)
	if err != nil {
		log.Printf("[session-search] Migration error: %v", err)
		return
	}

	// Rebuild FTS5 index
	_, err = sessionDB.Exec(`INSERT INTO messages_fts(messages_fts) VALUES('rebuild')`)
	if err != nil {
		log.Printf("[session-search] FTS5 rebuild error: %v", err)
	}

	log.Printf("[session-search] Migration to FTS5 complete")
}