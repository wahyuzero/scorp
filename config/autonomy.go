package config

import (
	"strings"
	"sync"
)

// ──────────────────────────────────────────────
// 3-Tier Autonomy System (ZeroClaw Parity)
// - readonly   : Strict read-only audit & inspection mode
// - supervised : Default mode. Safe tools run auto, dangerous require confirm
// - yolo       : Full unattended autonomy (bypasses confirmations)
// ──────────────────────────────────────────────

type AutonomyLevel string

const (
	AutonomyReadOnly   AutonomyLevel = "readonly"
	AutonomySupervised AutonomyLevel = "supervised"
	AutonomyYOLO       AutonomyLevel = "yolo"
)

var (
	currentAutonomy   AutonomyLevel = AutonomySupervised
	currentAutonomyMu sync.RWMutex

	// Tools permitted in ReadOnly mode
	readOnlyTools = map[string]bool{
		"read_file":      true,
		"list_dir":       true,
		"search_code":    true,
		"system_info":    true,
		"read_url":       true,
		"web_fetch":      true,
		"web_search":     true,
		"log":            true,
		"tool_search":    true,
		"tool_list":      true,
		"index_search":   true,
		"ragvec_search":  true,
		"session_search": true,
		"analyze_image":  true,
		"uptime":         true,
	}

	// Paths strictly forbidden from access across all modes (Sandboxing)
	sensitivePathPatterns = []string{
		"/etc/shadow",
		"/etc/sudoers",
		"id_rsa",
		"id_ed25519",
		"id_ecdsa",
		".ssh/authorized_keys",
		"/proc/kcore",
		"vault.key",
	}
)

// SetAutonomyLevel sets the global autonomy level (readonly, supervised, yolo)
func SetAutonomyLevel(level string) {
	currentAutonomyMu.Lock()
	defer currentAutonomyMu.Unlock()

	normalized := strings.ToLower(strings.TrimSpace(level))
	switch normalized {
	case "readonly", "ro", "audit":
		currentAutonomy = AutonomyReadOnly
	case "yolo", "full", "auto":
		currentAutonomy = AutonomyYOLO
	default:
		currentAutonomy = AutonomySupervised
	}
}

// GetAutonomyLevel returns the active autonomy level
func GetAutonomyLevel() AutonomyLevel {
	currentAutonomyMu.RLock()
	defer currentAutonomyMu.RUnlock()
	return currentAutonomy
}

// IsToolAllowed checks if a tool is permitted under the active autonomy mode
func IsToolAllowed(toolName string) (bool, string) {
	level := GetAutonomyLevel()
	if level == AutonomyReadOnly {
		if !readOnlyTools[toolName] {
			return false, "Blocked by ReadOnly mode: tool '" + toolName + "' can modify system or files"
		}
	}
	return true, ""
}

// IsPathRestricted checks if a file path targets sensitive system credentials
func IsPathRestricted(targetPath string) (bool, string) {
	lower := strings.ToLower(targetPath)
	for _, pattern := range sensitivePathPatterns {
		if strings.Contains(lower, pattern) {
			return true, "Access denied by Security Sandbox: sensitive path '" + pattern + "'"
		}
	}
	return false, ""
}
