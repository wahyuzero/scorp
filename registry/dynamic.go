package registry

import (
	"log"
	"os"
	"sync"
)

// ──────────────────────────────────────────────
// Dynamic Tool Discovery with TTL Injection (PicoClaw Parity)
// Keeps core tools active while dynamically injecting specialized tools
// with a Turns-To-Live (TTL) counter when discovered or invoked.
// ──────────────────────────────────────────────

var (
	dynamicTTL   = make(map[string]int)
	dynamicTTLMu sync.Mutex

	// Core tools always included in native schema
	coreTools = map[string]bool{
		"read_file":            true,
		"replace_file_content": true,
		"write_file":           true,
		"list_dir":             true,
		"search_code":          true,
		"shell":                true,
		"read_url":             true,
		"sop":                  true,
		"tool_search":          true,
		"tool_call":            true,
	}
)

// IsDynamicModeEnabled checks if dynamic tool discovery is enabled
func IsDynamicModeEnabled() bool {
	return os.Getenv("SCORP_DYNAMIC_TOOLS") == "true" || os.Getenv("SCORP_DYNAMIC_TOOLS") == "1"
}

// IsCoreTool checks if a tool is a primary core tool
func IsCoreTool(name string) bool {
	return coreTools[name]
}

// ActivateToolWithTTL temporarily injects a tool into active schema for N turns
func ActivateToolWithTTL(name string, ttl int) {
	if ttl <= 0 {
		ttl = 3
	}
	dynamicTTLMu.Lock()
	defer dynamicTTLMu.Unlock()

	dynamicTTL[name] = ttl
	ResetNativeToolCache()
	if os.Getenv("SCORP_DEBUG") != "" {
		log.Printf("[dynamic-tools] Activated tool '%s' with TTL=%d turns", name, ttl)
	}
}

// TickToolTTL decrements the TTL of dynamically activated tools at turn boundaries
func TickToolTTL() {
	dynamicTTLMu.Lock()
	defer dynamicTTLMu.Unlock()

	changed := false
	for name, ttl := range dynamicTTL {
		ttl--
		if ttl <= 0 {
			delete(dynamicTTL, name)
			changed = true
			if os.Getenv("SCORP_DEBUG") != "" {
				log.Printf("[dynamic-tools] Tool '%s' TTL expired, returning to deferred", name)
			}
		} else {
			dynamicTTL[name] = ttl
		}
	}

	if changed {
		ResetNativeToolCache()
	}
}

// IsToolActive determines if a tool definition should be included in native LLM schema
func IsToolActive(def ToolDef) bool {
	// Static mode (default): all non-deferred native tools are active, and
	// deferred tools (e.g. MCP, P2.9) pop in for their TTL window once
	// discovered via tool_search or invoked via tool_call.
	if !IsDynamicModeEnabled() {
		if def.Native && !def.Deferred {
			return true
		}
		dynamicTTLMu.Lock()
		defer dynamicTTLMu.Unlock()
		return dynamicTTL[def.Name] > 0
	}

	// Dynamic mode: core tools are always active
	if coreTools[def.Name] {
		return true
	}

	// Check if temporarily activated with TTL
	dynamicTTLMu.Lock()
	defer dynamicTTLMu.Unlock()
	return dynamicTTL[def.Name] > 0
}

// ResetDynamicTools clears all TTL injections
func ResetDynamicTools() {
	dynamicTTLMu.Lock()
	defer dynamicTTLMu.Unlock()
	dynamicTTL = make(map[string]int)
	ResetNativeToolCache()
}
