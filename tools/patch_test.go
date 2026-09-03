package tools

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExecuteReplaceFileContent(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")

	initialContent := `line 1: func hello() {
line 2:     fmt.Println("old")
line 3: }
line 4: func world() {
line 5:     fmt.Println("test")
line 6: }`

	if err := os.WriteFile(filePath, []byte(initialContent), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Test 1: Modern args (target_file, target_content, replacement_content)
	res, ok := ExecuteReplaceFileContent(map[string]interface{}{
		"target_file":         filePath,
		"target_content":      `fmt.Println("old")`,
		"replacement_content": `fmt.Println("new")`,
	})
	if !ok {
		t.Fatalf("Replace failed: %s", res)
	}
	if !strings.Contains(res, "✅ Surgically edited") {
		t.Errorf("Expected success message, got %s", res)
	}

	updated, _ := os.ReadFile(filePath)
	if !strings.Contains(string(updated), `fmt.Println("new")`) {
		t.Errorf("Content not updated properly: %s", string(updated))
	}

	// Test 2: Scoped line replace
	res2, ok2 := ExecuteReplaceFileContent(map[string]interface{}{
		"target_file":         filePath,
		"target_content":      `fmt.Println("test")`,
		"replacement_content": `fmt.Println("scoped")`,
		"start_line":          4,
		"end_line":            6,
	})
	if !ok2 {
		t.Fatalf("Scoped replace failed: %s", res2)
	}

	updated2, _ := os.ReadFile(filePath)
	if !strings.Contains(string(updated2), `fmt.Println("scoped")`) {
		t.Errorf("Scoped content not updated: %s", string(updated2))
	}

	// Test 3: Backward compatible ExecutePatch
	res3, ok3 := ExecutePatch(map[string]interface{}{
		"path":       filePath,
		"old_string": `fmt.Println("new")`,
		"new_string": `fmt.Println("patched")`,
	})
	if !ok3 {
		t.Fatalf("ExecutePatch failed: %s", res3)
	}
	updated3, _ := os.ReadFile(filePath)
	if !strings.Contains(string(updated3), `fmt.Println("patched")`) {
		t.Errorf("Patch not updated: %s", string(updated3))
	}
}
