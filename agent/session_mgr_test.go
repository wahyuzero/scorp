package agent

import (
	"os"
	"path/filepath"
	"testing"

	"scorp-agent/config"
)

func TestSessionManager(t *testing.T) {
	config.InitConfigManager()
	dir := config.HistoryDirPath()
	os.MkdirAll(dir, 0755)

	test1 := "test_session_mgr_1"
	test2 := "test_session_mgr_2"

	// Cleanup
	_ = DeleteSession(test1)
	_ = DeleteSession(test2)

	// Save dummy session
	msgs := []AgentMessage{
		{Role: "user", Content: "Hello world"},
		{Role: "assistant", Content: "Hi there!"},
	}
	saveHistoryToDisk(test1, msgs)

	if !SessionExists(test1) {
		t.Fatalf("Expected session %s to exist", test1)
	}

	// List sessions
	list := ListSessions()
	found := false
	for _, s := range list {
		if s.ID == test1 {
			found = true
			if s.MsgCount != 2 {
				t.Errorf("Expected 2 messages, got %d", s.MsgCount)
			}
			break
		}
	}
	if !found {
		t.Errorf("Expected %s in ListSessions()", test1)
	}

	// Rename session
	if err := RenameSession(test1, test2); err != nil {
		t.Fatalf("RenameSession failed: %v", err)
	}
	if SessionExists(test1) {
		t.Errorf("Expected old session %s to no longer exist", test1)
	}
	if !SessionExists(test2) {
		t.Errorf("Expected renamed session %s to exist", test2)
	}

	// Delete session
	if err := DeleteSession(test2); err != nil {
		t.Fatalf("DeleteSession failed: %v", err)
	}
	if SessionExists(test2) {
		t.Errorf("Expected session %s to be deleted", test2)
	}

	// Verify file is gone from disk
	p := filepath.Join(dir, test2+".json")
	if _, err := os.Stat(p); !os.IsNotExist(err) {
		t.Errorf("Expected file %s to be deleted", p)
	}
}
