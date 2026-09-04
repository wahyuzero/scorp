package telegram

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTelegramHelpers(t *testing.T) {
	// 1. Test HumanSize
	tests := []struct {
		size     int64
		expected string
	}{
		{500, "500B"},
		{1024, "1.0KB"},
		{1024 * 1024 * 5, "5.0MB"},
		{1024 * 1024 * 1024 * 2, "2.00GB"},
	}

	for _, tt := range tests {
		got := HumanSize(tt.size)
		if got != tt.expected {
			t.Errorf("HumanSize(%d) = %s, expected %s", tt.size, got, tt.expected)
		}
	}

	// 2. Test PathID and GetPath roundtrip
	testPath := "/tmp/scorp/test.txt"
	pid := PathID(testPath)
	if pid == "" {
		t.Errorf("PathID returned empty string")
	}
	recovered := GetPath(pid)
	if recovered != testPath {
		t.Errorf("GetPath(%s) = %s, expected %s", pid, recovered, testPath)
	}

	// 3. Test baseName
	if baseName("/etc/systemd/system/scorp.service") != "scorp.service" {
		t.Errorf("baseName failed on full path")
	}
	if baseName("file.go") != "file.go" {
		t.Errorf("baseName failed on filename")
	}

	// 4. Test Keyboards
	if MainMenuKeyboard() == nil {
		t.Errorf("MainMenuKeyboard returned nil")
	}
	if MonitorMenuKeyboard() == nil {
		t.Errorf("MonitorMenuKeyboard returned nil")
	}
	if SystemMenuKeyboard() == nil {
		t.Errorf("SystemMenuKeyboard returned nil")
	}
	if SettingsMenuKeyboard() == nil {
		t.Errorf("SettingsMenuKeyboard returned nil")
	}
}

func TestWebhookHandlerMalformedJSON(t *testing.T) {
	req := httptest.NewRequest("POST", "/webhook", bytes.NewBufferString("{malformed json"))
	w := httptest.NewRecorder()

	WebhookHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request on malformed JSON, got %d", resp.StatusCode)
	}
}
