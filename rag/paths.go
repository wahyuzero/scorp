package rag

import (
	"os"
	"path/filepath"
)

// ──────────────────────────────────────────────
// Path helpers — inlined from config_paths.go
// Will be consolidated into config/ package in Phase 2
// ──────────────────────────────────────────────

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	if h, err := os.UserHomeDir(); err == nil {
		return h
	}
	return "/tmp"
}

func ragDirPath() string {
	return filepath.Join(homeDir(), ".scorp", "rag")
}

func ragIndexPath() string {
	return filepath.Join(ragDirPath(), "index.json")
}

func ragVectorDBPath() string {
	return filepath.Join(ragDirPath(), "vector_index.json")
}
