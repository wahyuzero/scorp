package delegate

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// ──────────────────────────────────────────────
// S08: Delegate Isolation
//
// Provides fully isolated context per subagent:
// - Working directory sandbox
// - Read-only global state enforcement
// - Output capture and size limits
// - Per-subagent environment variables
// ──────────────────────────────────────────────

// SubagentIsolation controls what a subagent can access.
type SubagentIsolation struct {
	WorkDir      string            // sandboxed working directory
	ReadOnlyMode bool              // if true, all tools execute read-only
	EnvVars      map[string]string // per-subagent env overrides
	MaxOutputKB  int               // max output size per tool call (default 50KB)
	AllowedPaths []string          // whitelisted read paths (empty = all readable)
	BlockedPaths []string          // blacklisted paths (always blocked)
}

var (
	subagentIsolations   = make(map[string]*SubagentIsolation)
	subagentIsolationMu  sync.RWMutex
)

// defaultIsolation returns a safe default isolation context.
func defaultIsolation() *SubagentIsolation {
	return &SubagentIsolation{
		ReadOnlyMode: false,
		MaxOutputKB:  50,
		EnvVars:      make(map[string]string),
	}
}

// createSubagentSandbox creates a temp working directory for a subagent.
func createSubagentSandbox(subagentID string) string {
	dir := filepath.Join(os.TempDir(), "scorp-sandbox", subagentID)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Printf("[isolation] Failed to create sandbox %s: %v", dir, err)
		return "/tmp"
	}
	return dir
}

// cleanupSubagentSandbox removes the temp directory after subagent completes.
func cleanupSubagentSandbox(subagentID string) {
	dir := filepath.Join(os.TempDir(), "scorp-sandbox", subagentID)
	if err := os.RemoveAll(dir); err != nil {
		log.Printf("[isolation] Failed to cleanup sandbox %s: %v", dir, err)
	}
}

// registerSubagentIsolation stores isolation config for a subagent.
func registerSubagentIsolation(subagentID string, iso *SubagentIsolation) {
	subagentIsolationMu.Lock()
	defer subagentIsolationMu.Unlock()
	subagentIsolations[subagentID] = iso
}

// getSubagentIsolation retrieves isolation config. Returns nil if not found.
func getSubagentIsolation(subagentID string) *SubagentIsolation {
	subagentIsolationMu.RLock()
	defer subagentIsolationMu.RUnlock()
	return subagentIsolations[subagentID]
}

// unregisterSubagentIsolation removes isolation config.
func unregisterSubagentIsolation(subagentID string) {
	subagentIsolationMu.Lock()
	defer subagentIsolationMu.Unlock()
	delete(subagentIsolations, subagentID)
}

// isPathBlocked checks if a path is in the blocked list or outside allowed paths.
func (iso *SubagentIsolation) isPathBlocked(path string) bool {
	abs, err := filepath.Abs(path)
	if err != nil {
		return true
	}

	// Check blocked paths
	for _, bp := range iso.BlockedPaths {
		bpAbs, _ := filepath.Abs(bp)
		if strings.HasPrefix(abs, bpAbs) {
			return true
		}
	}

	// Check allowed paths whitelist
	if len(iso.AllowedPaths) > 0 {
		allowed := false
		for _, ap := range iso.AllowedPaths {
			apAbs, _ := filepath.Abs(ap)
			if strings.HasPrefix(abs, apAbs) {
				allowed = true
				break
			}
		}
		if !allowed {
			return true
		}
	}

	return false
}

// truncateOutput enforces the max output size.
func (iso *SubagentIsolation) truncateOutput(output string) string {
	maxBytes := iso.MaxOutputKB * 1024
	if maxBytes <= 0 {
		maxBytes = 50 * 1024
	}
	if len(output) <= maxBytes {
		return output
	}
	return output[:maxBytes] + "\n\n... [output truncated by isolation limit]"
}

// Tools that subagents should NEVER modify (global state mutations)
var subagentBlockedTools = map[string]bool{
	"schedule":    true, // don't let subagents create/modify cron jobs
	"vault":       true, // don't let subagents modify credential vault
	"monitor":     true, // don't let subagents modify monitoring targets
	"autonomous":  true, // don't let subagents toggle autonomous mode
	"models":      true, // don't let subagents modify model config
	"ragvec_add":  true, // don't let subagents pollute RAG index
	"ragvec_remove": true,
}

// isSubagentToolBlocked checks if a tool should be blocked for subagents.
// Even if the tool is in params.Tools, these state-modifying tools are always blocked.
func isSubagentToolBlocked(toolName string) bool {
	return subagentBlockedTools[toolName]
}

// formatIsolationInfo returns a human-readable summary of isolation config.
func formatIsolationInfo(iso *SubagentIsolation) string {
	if iso == nil {
		return "No isolation configured"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("  WorkDir: %s\n", iso.WorkDir))
	sb.WriteString(fmt.Sprintf("  ReadOnly: %v\n", iso.ReadOnlyMode))
	sb.WriteString(fmt.Sprintf("  MaxOutput: %dKB\n", iso.MaxOutputKB))
	if len(iso.AllowedPaths) > 0 {
		sb.WriteString(fmt.Sprintf("  AllowedPaths: %v\n", iso.AllowedPaths))
	}
	if len(iso.BlockedPaths) > 0 {
		sb.WriteString(fmt.Sprintf("  BlockedPaths: %v\n", iso.BlockedPaths))
	}
	return sb.String()
}
