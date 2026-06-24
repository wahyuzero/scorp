package session

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

func scorpDir() string {
	return filepath.Join(homeDir(), ".scorp")
}

func scorpPath(parts ...string) string {
	all := append([]string{scorpDir()}, parts...)
	return filepath.Join(all...)
}
