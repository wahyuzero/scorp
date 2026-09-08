package tools

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExecuteSendFile_SafetyAndDispatch(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "scorp_test_send_")
	if err != nil {
		t.Fatalf("MkdirTemp failed: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testFilePath := filepath.Join(tmpDir, "test_doc.pdf")
	if err := os.WriteFile(testFilePath, []byte("%PDF-1.4 test document content"), 0644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	// 1. Missing path arg
	out, ok := ExecuteSendFile(map[string]interface{}{}, 12345)
	if ok || out != "Error: 'path' argument is required" {
		t.Errorf("expected missing path error, got: %s (ok=%v)", out, ok)
	}

	// 2. Non-existent file
	out, ok = ExecuteSendFile(map[string]interface{}{"path": filepath.Join(tmpDir, "ghost.png")}, 12345)
	if ok {
		t.Errorf("expected error for non-existent file, got: %s", out)
	}

	// 3. Directory path instead of file
	out, ok = ExecuteSendFile(map[string]interface{}{"path": tmpDir}, 12345)
	if ok {
		t.Errorf("expected error when sending a directory, got: %s", out)
	}

	// 4. Standalone (nil SendMedia & nil SendDocumentBytes) must NOT panic
	SendMedia = nil
	SendDocumentBytes = nil
	out, ok = ExecuteSendFile(map[string]interface{}{"path": testFilePath}, 12345)
	if !ok {
		t.Errorf("expected standalone verification success, got: %s", out)
	}

	// 5. Mock SendMedia callback dispatch
	var capturedChatID, capturedPath, capturedCaption string
	var capturedAsDoc bool
	SendMedia = func(chatID string, filePath string, caption string, asDocument bool) (bool, string) {
		capturedChatID = chatID
		capturedPath = filePath
		capturedCaption = caption
		capturedAsDoc = asDocument
		return true, "document"
	}

	out, ok = ExecuteSendFile(map[string]interface{}{
		"path":        testFilePath,
		"caption":     "My Monthly Report",
		"as_document": true,
	}, 98765)

	if !ok {
		t.Errorf("expected send success, got: %s", out)
	}
	if capturedChatID != "98765" {
		t.Errorf("expected chatID 98765, got: %s", capturedChatID)
	}
	if capturedPath != testFilePath {
		t.Errorf("expected path %s, got: %s", testFilePath, capturedPath)
	}
	if capturedCaption != "My Monthly Report" {
		t.Errorf("expected caption 'My Monthly Report', got: %s", capturedCaption)
	}
	if !capturedAsDoc {
		t.Errorf("expected as_document to be true")
	}

	// Reset callbacks
	SendMedia = nil
	SendDocumentBytes = nil
}
