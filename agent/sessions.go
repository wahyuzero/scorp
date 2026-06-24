package agent

import (
	"sync"
	"time"
)

// ──────────────────────────────────────────────
// Shared types
// ──────────────────────────────────────────────

// TGDocument represents an incoming file/photo from Telegram
type TGDocument struct {
	FileID   string
	FileName string
	FileSize int64
	ChatID   int64
	MsgID    int64
	Caption  string
	IsPhoto  bool
}

// ──────────────────────────────────────────────
// Agent session tracking
// ──────────────────────────────────────────────

type agentSession struct {
	active   bool
	lastUsed time.Time
}

var (
	agentSessions   = make(map[string]*agentSession)
	agentSessionsMu sync.Mutex
)

// cleanupAgentSessions removes old agent sessions.
func cleanupAgentSessions() {
	agentSessionsMu.Lock()
	defer agentSessionsMu.Unlock()
	cutoff := time.Now().Add(-30 * time.Minute)
	for id, sess := range agentSessions {
		if !sess.active && sess.lastUsed.Before(cutoff) {
			delete(agentSessions, id)
		}
	}
}
