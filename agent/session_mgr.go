package agent

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"scorp-agent/config"
	"scorp-agent/internal/helpers"
)

// ──────────────────────────────────────────────
// Multi-Session Management Engine
// Supports list, new, switch, rename, and delete
// ──────────────────────────────────────────────

// SessionMeta holds summary metadata for a saved chat session
type SessionMeta struct {
	ID           string    `json:"id"`
	MsgCount     int       `json:"msg_count"`
	LastModified time.Time `json:"last_modified"`
	Size         int64     `json:"size"`
	LastPreview  string    `json:"last_preview"`
}

// ListSessions returns all saved sessions sorted by last modified descending
func ListSessions() []SessionMeta {
	dir := config.HistoryDirPath()
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	var list []SessionMeta
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		name := strings.TrimSuffix(e.Name(), ".json")
		info, err := e.Info()
		if err != nil {
			continue
		}

		// Read file to inspect message count & preview
		filePath := filepath.Join(dir, e.Name())
		var msgs []AgentMessage
		lastPreview := "-"
		if data, err := os.ReadFile(filePath); err == nil {
			if err := json.Unmarshal(data, &msgs); err == nil && len(msgs) > 0 {
				// Find last user/assistant message for preview
				for i := len(msgs) - 1; i >= 0; i-- {
					if content, ok := msgs[i].Content.(string); ok && content != "" {
						lastPreview = helpers.TruncateStr(strings.ReplaceAll(content, "\n", " "), 50)
						break
					}
				}
			}
		}

		list = append(list, SessionMeta{
			ID:           name,
			MsgCount:     len(msgs),
			LastModified: info.ModTime(),
			Size:         info.Size(),
			LastPreview:  lastPreview,
		})
	}

	// Sort by last modified descending
	sort.Slice(list, func(i, j int) bool {
		return list[i].LastModified.After(list[j].LastModified)
	})

	return list
}

// RenameSession renames an existing session on disk and in memory
func RenameSession(oldID, newID string) error {
	oldID = sanitizeSessionID(oldID)
	newID = sanitizeSessionID(newID)

	if oldID == "" || newID == "" {
		return fmt.Errorf("session names cannot be empty")
	}
	if oldID == newID {
		return nil
	}

	oldPath := historyFilePath(oldID)
	newPath := historyFilePath(newID)

	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		return fmt.Errorf("session '%s' does not exist", oldID)
	}
	if _, err := os.Stat(newPath); err == nil {
		return fmt.Errorf("session '%s' already exists", newID)
	}

	// Rename on disk
	if err := os.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("failed to rename session file: %w", err)
	}

	// Move in memory
	oldSess := getSession(oldID)
	if oldSess != nil {
		setSession(newID, oldSess)
		m, mu := getSessionMap(oldID)
		mu.Lock()
		delete(m, oldID)
		mu.Unlock()
	}

	return nil
}

// DeleteSession deletes a session file and purges its memory
func DeleteSession(sessionID string) error {
	sessionID = sanitizeSessionID(sessionID)
	if sessionID == "" {
		return fmt.Errorf("session name cannot be empty")
	}

	ClearChatSession(sessionID)

	filePath := historyFilePath(sessionID)
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete session file: %w", err)
	}

	return nil
}

// SessionExists checks if a session file or memory exists
func SessionExists(sessionID string) bool {
	sessionID = sanitizeSessionID(sessionID)
	if sessionID == "" {
		return false
	}
	if sess := getSession(sessionID); sess != nil && len(sess.history) > 0 {
		return true
	}
	filePath := historyFilePath(sessionID)
	_, err := os.Stat(filePath)
	return err == nil
}

func sanitizeSessionID(id string) string {
	id = strings.TrimSpace(id)
	id = strings.ReplaceAll(id, "/", "_")
	id = strings.ReplaceAll(id, "\\", "_")
	id = strings.ReplaceAll(id, "..", "")
	return id
}
