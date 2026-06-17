package main

import (
	"sync"
	"time"
)

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
