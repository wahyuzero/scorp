package main

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
// Session Search — SQLite + FTS5 full-text search
// over conversation history. Allows the agent to
// recall past conversations by keyword/phrase.
// ──────────────────────────────────────────────

var (
	sessionDB   *sql.DB
	sessionDBMu sync.Mutex
)

// initSessionDB opens the SQLite database and creates tables.
func initSessionDB() {
	dbPath := scorpPath("sessions.db")

	var err error
	sessionDB, err = sql.Open("sqlite3", dbPath+"?_journal_mode=WAL&_busy_timeout=5000&_foreign_keys=on")
	if err != nil {
		log.Printf("[session-search] Failed to open DB: %v", err)
		return
	}

	// Optimize for low-contention
	sessionDB.SetMaxOpenConns(5)

	// Create tables
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

	CREATE VIRTUAL TABLE IF NOT EXISTS messages_fts USING fts5(
		content,
		content='messages',
		content_rowid='id',
		tokenize='porter unicode61'
	);

	-- Triggers to keep FTS in sync
	CREATE TRIGGER IF NOT EXISTS messages_ai AFTER INSERT ON messages BEGIN
		INSERT INTO messages_fts(rowid, content) VALUES (new.id, new.content);
	END;
	CREATE TRIGGER IF NOT EXISTS messages_ad AFTER DELETE ON messages BEGIN
		INSERT INTO messages_fts(messages_fts, rowid, content) VALUES ('delete', old.id, old.content);
	END;
	CREATE TRIGGER IF NOT EXISTS messages_au AFTER UPDATE ON messages BEGIN
		INSERT INTO messages_fts(messages_fts, rowid, content) VALUES ('delete', old.id, old.content);
		INSERT INTO messages_fts(rowid, content) VALUES (new.id, new.content);
	END;
	`

	if _, err := sessionDB.Exec(schema); err != nil {
		log.Printf("[session-search] Failed to create schema: %v", err)
		return
	}

	log.Printf("[session-search] SQLite DB initialized at %s", dbPath)
}

// indexMessage stores a message in the session DB for later FTS search.
// Called asynchronously to avoid blocking the agent loop.
func indexMessage(chatID, role, content string) {
	if sessionDB == nil {
		return
	}

	// Truncate very long messages
	if len(content) > 10000 {
		content = content[:10000]
	}

	msgIndex := time.Now().UnixNano()

	sessionDBMu.Lock()
	defer sessionDBMu.Unlock()

	_, err := sessionDB.Exec(
		"INSERT INTO messages (chat_id, role, content, timestamp, msg_index) VALUES (?, ?, ?, ?, ?)",
		chatID, role, content, time.Now().Unix(), msgIndex,
	)
	if err != nil {
		log.Printf("[session-search] Failed to index message: %v", err)
	}
}

// searchSessions performs FTS5 full-text search over conversation history.
// Returns matching messages with surrounding context.
func searchSessions(query string, chatID string, limit int) ([]SessionSearchResult, error) {
	if sessionDB == nil {
		return nil, fmt.Errorf("session DB not initialized")
	}

	if limit <= 0 || limit > 50 {
		limit = 10
	}

	// Build FTS5 query — support simple keyword search
	// Escape special FTS5 characters and build a prefix-match query
	ftsQuery := buildFTSQuery(query)
	if ftsQuery == "" {
		return nil, fmt.Errorf("invalid search query")
	}

	var rows *sql.Rows
	var err error

	if chatID != "" {
		rows, err = sessionDB.Query(`
			SELECT m.id, m.chat_id, m.role, m.content, m.timestamp,
			       bm25(messages_fts) as rank
			FROM messages_fts
			JOIN messages m ON m.id = messages_fts.rowid
			WHERE messages_fts MATCH ?
			  AND m.chat_id = ?
			ORDER BY rank
			LIMIT ?`,
			ftsQuery, chatID, limit)
	} else {
		rows, err = sessionDB.Query(`
			SELECT m.id, m.chat_id, m.role, m.content, m.timestamp,
			       bm25(messages_fts) as rank
			FROM messages_fts
			JOIN messages m ON m.id = messages_fts.rowid
			WHERE messages_fts MATCH ?
			ORDER BY rank
			LIMIT ?`,
			ftsQuery, limit)
	}

	if err != nil {
		return nil, fmt.Errorf("search query failed: %w", err)
	}
	defer rows.Close()

	var results []SessionSearchResult
	for rows.Next() {
		var r SessionSearchResult
		var ts int64
		var rank float64
		if err := rows.Scan(&r.ID, &r.ChatID, &r.Role, &r.Content, &ts, &rank); err != nil {
			continue
		}
		r.Timestamp = time.Unix(ts, 0)
		r.Rank = rank
		// Truncate for display
		if len(r.Content) > 500 {
			r.Content = r.Content[:500] + "..."
		}
		results = append(results, r)
	}

	return results, nil
}

// buildFTSQuery converts a plain-text query into an FTS5 MATCH expression.
// Supports multi-word queries (implicit AND) and prefix matching.
func buildFTSQuery(query string) string {
	query = strings.TrimSpace(query)
	if query == "" {
		return ""
	}

	// Split into words, build prefix-match query
	words := strings.Fields(query)
	var parts []string
	for _, w := range words {
		// Escape double quotes
		w = strings.ReplaceAll(w, "\"", "")
		if w == "" {
			continue
		}
		// Prefix match: word*  (matches "deploy" when searching "dep")
		parts = append(parts, w+"*")
	}

	if len(parts) == 0 {
		return ""
	}

	// Join with space (implicit AND in FTS5)
	return strings.Join(parts, " ")
}

// SessionSearchResult represents a single search hit
type SessionSearchResult struct {
	ID        int64
	ChatID    string
	Role      string
	Content   string
	Timestamp time.Time
	Rank      float64
}

// executeSessionSearch is the agent tool handler
func executeSessionSearch(args map[string]interface{}) (string, bool) {
	query := getStringArg(args, "query", "")
	if query == "" {
		return "Error: 'query' is required", false
	}

	limit := getIntArg(args, "limit", 10)
	scope := getStringArg(args, "scope", "all") // "all" or "current"

	var chatID string
	if scope == "current" {
		// Use the calling chat's ID — passed via args or empty for "all"
		chatID = getStringArg(args, "chat_id", "")
	}

	results, err := searchSessions(query, chatID, limit)
	if err != nil {
		return fmt.Sprintf("Search error: %v", err), false
	}

	if len(results) == 0 {
		return fmt.Sprintf("No matches found for \"%s\"", query), true
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("🔍 Found %d result(s) for \"%s\":\n\n", len(results), query))

	for i, r := range results {
		role := r.Role
		if role == "user" {
			role = "👤"
		} else if role == "assistant" {
			role = "🤖"
		}

		timeStr := r.Timestamp.Format("2006-01-02 15:04")
		sb.WriteString(fmt.Sprintf("**%d.** %s [%s, chat:%s]\n", i+1, role, timeStr, r.ChatID))
		sb.WriteString(fmt.Sprintf("   %s\n\n", r.Content))
	}

	return sb.String(), true
}

// browseSessions returns recent sessions chronologically
func executeSessionBrowse(args map[string]interface{}) (string, bool) {
	if sessionDB == nil {
		return "Session DB not initialized", false
	}

	limit := getIntArg(args, "limit", 20)
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	rows, err := sessionDB.Query(`
		SELECT chat_id, role, content, timestamp
		FROM messages
		ORDER BY timestamp DESC
		LIMIT ?`, limit)
	if err != nil {
		return fmt.Sprintf("Browse error: %v", err), false
	}
	defer rows.Close()

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📋 Recent conversations (last %d messages):\n\n", limit))

	for rows.Next() {
		var chatID, role, content string
		var ts int64
		if err := rows.Scan(&chatID, &role, &content, &ts); err != nil {
			continue
		}

		timeStr := time.Unix(ts, 0).Format("01-02 15:04")
		emoji := "👤"
		if role == "assistant" {
			emoji = "🤖"
		}
		sb.WriteString(fmt.Sprintf("%s [%s] %s\n", emoji, timeStr, truncateStr(content, 80)))
	}

	return sb.String(), true
}
