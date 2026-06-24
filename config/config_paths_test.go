package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigPaths_ScorpDir(t *testing.T) {
	// Test that ScorpDir() returns expected path based on HOME
	os.Setenv("HOME", "/tmp/testhome")
	defer os.Unsetenv("HOME")

	dir := ScorpDir()
	expected := filepath.Join("/tmp/testhome", ".scorp")
	if dir != expected {
		t.Errorf("ScorpDir() = %q, want %q", dir, expected)
	}
}

func TestConfigPaths_ScreenshotsDir(t *testing.T) {
	os.Setenv("HOME", "/tmp/testhome")
	defer os.Unsetenv("HOME")

	dir := ScreenshotsDir()
	expected := filepath.Join("/tmp/testhome", ".scorp", "screenshots")
	if dir != expected {
		t.Errorf("ScreenshotsDir() = %q, want %q", dir, expected)
	}
}

func TestConfigPaths_RagVectorDBPath(t *testing.T) {
	os.Setenv("HOME", "/tmp/testhome")
	defer os.Unsetenv("HOME")

	dir := RagVectorDBPath()
	expected := filepath.Join("/tmp/testhome", ".scorp", "rag", "vector_index.json")
	if dir != expected {
		t.Errorf("RagVectorDBPath() = %q, want %q", dir, expected)
	}
}