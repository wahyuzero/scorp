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
		"clarify":        true,
	}

	// Plan Mode (P1.4): while a /plan draft runs, IsToolAllowed restricts
	// every tool dispatch to the read-only set regardless of autonomy level.
	planningMode   bool
	planningModeMu sync.RWMutex

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

// ConfirmationRequired reports whether the active autonomy mode still routes
// risky actions through human confirmation. Every confirmation gate (agent
// loop, shell, sql, process, git) must consult this single predicate so the
// autonomy contract stays uniform: only YOLO runs unattended; readonly and
// supervised always require confirmation.
func ConfirmationRequired() bool {
	return GetAutonomyLevel() != AutonomyYOLO
}

// SetPlanningMode toggles Plan Mode (P1.4): while a /plan draft is running,
// IsToolAllowed restricts every tool dispatch to the read-only set regardless
// of the autonomy level. task_plan and complete_task never reach this gate —
// both are intercepted inside the agent loop.
func SetPlanningMode(on bool) {
	planningModeMu.Lock()
	defer planningModeMu.Unlock()
	planningMode = on
}

// PlanningModeActive reports whether a plan draft is restricting tools.
func PlanningModeActive() bool {
	planningModeMu.RLock()
	defer planningModeMu.RUnlock()
	return planningMode
}

// IsToolAllowed checks if a tool is permitted under the active autonomy mode
func IsToolAllowed(toolName string) (bool, string) {
	level := GetAutonomyLevel()
	if level == AutonomyReadOnly {
		if !readOnlyTools[toolName] {
			return false, "Blocked by ReadOnly mode: tool '" + toolName + "' can modify system or files"
		}
	}
	if PlanningModeActive() && !readOnlyTools[toolName] {
		return false, "Blocked by Plan Mode: only read-only tools while drafting a plan — execution starts after the plan is approved"
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
