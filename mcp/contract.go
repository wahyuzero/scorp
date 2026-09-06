package mcp

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"scorp-agent/config"
)

// ──────────────────────────────────────────────
// MCP Contract Watch (P3.14)
//
// Fingerprint (SHA-256 of name+description+inputSchema per tool, sorted) of
// every registered MCP server, persisted in ~/.scorp/mcp_contracts.json.
// On each startup the fingerprints are re-computed and compared: a server
// that still answers but silently ships a different toolset triggers a
// warning (logged + surfaced in /status). "Uptime tells you the server
// responds — not that it still does what the agent expects."
//
// The stored contract is updated after the warning, so each change warns
// exactly once. A server that disappears from the registry while still
// configured keeps warning until it returns or is removed from mcp.json.
// ──────────────────────────────────────────────

// ServerContract is the persisted fingerprint record for one MCP server.
type ServerContract struct {
	Server      string    `json:"server"`
	Fingerprint string    `json:"fingerprint"`
	ToolCount   int       `json:"tool_count"`
	Tools       []string  `json:"tools"`
	CheckedAt   time.Time `json:"checked_at"`
}

var contractWarnings = struct {
	sync.Mutex
	list []string
}{}

// contractPathOverride is set by tests/eval to isolate the persisted contract
// file; empty means the default ~/.scorp/mcp_contracts.json.
var contractPathOverride string

// SetContractPath overrides the persisted contract file location (tests/eval);
// empty restores the default.
func SetContractPath(p string) {
	contractPathOverride = p
}

func contractFile() string {
	if contractPathOverride != "" {
		return contractPathOverride
	}
	return config.ScorpPath("mcp_contracts.json")
}

// serverFingerprint hashes the toolset of one server deterministically.
func serverFingerprint(tools []MCPTool) (string, []string) {
	lines := make([]string, 0, len(tools))
	names := make([]string, 0, len(tools))
	for _, t := range tools {
		schema, _ := json.Marshal(t.InputSchema)
		lines = append(lines, t.Name+"|"+t.Description+"|"+string(schema))
		names = append(names, t.Name)
	}
	sort.Strings(lines)
	sort.Strings(names)
	h := sha256.Sum256([]byte(strings.Join(lines, "\n")))
	return hex.EncodeToString(h[:16]), names
}

// CheckServerContracts compares the live toolset of every registered MCP
// server against the persisted contract and returns the warnings.
func CheckServerContracts() []string {
	mcpServersMu.RLock()
	snap := make(map[string][]MCPTool, len(mcpServers))
	for name, srv := range mcpServers {
		snap[name] = srv.tools
	}
	mcpServersMu.RUnlock()

	path := contractFile()
	stored := map[string]ServerContract{}
	if data, err := os.ReadFile(path); err == nil {
		_ = json.Unmarshal(data, &stored)
	}

	names := make([]string, 0, len(snap))
	for name := range snap {
		names = append(names, name)
	}
	sort.Strings(names)

	var warnings []string
	dirty := false
	now := time.Now()
	for _, name := range names {
		tools := snap[name]
		fp, toolNames := serverFingerprint(tools)
		old, had := stored[name]
		if had && old.Fingerprint != fp {
			w := fmt.Sprintf("MCP contract changed for %q: %d → %d tools (fp %s… → %s…)",
				name, old.ToolCount, len(tools), old.Fingerprint[:8], fp[:8])
			warnings = append(warnings, w)
			log.Printf("[mcp-watch] ⚠️ %s", w)
		}
		if !had || old.Fingerprint != fp {
			stored[name] = ServerContract{Server: name, Fingerprint: fp, ToolCount: len(tools), Tools: toolNames, CheckedAt: now}
			dirty = true
		}
	}
	// Registered servers that vanished from the live registry keep warning.
	for name, old := range stored {
		if _, ok := snap[name]; !ok {
			w := fmt.Sprintf("MCP server %q no longer registered (contract still on file, %d tools)", name, old.ToolCount)
			warnings = append(warnings, w)
		}
	}
	if dirty {
		if data, err := json.MarshalIndent(stored, "", "  "); err == nil {
			_ = os.WriteFile(path, data, 0644)
		}
	}

	contractWarnings.Lock()
	contractWarnings.list = warnings
	contractWarnings.Unlock()
	return warnings
}

// ContractWarnings returns the warnings of the most recent contract check.
func ContractWarnings() []string {
	contractWarnings.Lock()
	defer contractWarnings.Unlock()
	out := make([]string, len(contractWarnings.list))
	copy(out, contractWarnings.list)
	return out
}

// ContractStatusNotice formats the warnings for chat status surfaces (empty
// string when the contract is unchanged).
func ContractStatusNotice() string {
	ws := ContractWarnings()
	if len(ws) == 0 {
		return ""
	}
	out := "\n\n⚠️ <b>MCP contract watch</b>\n"
	for _, w := range ws {
		out += "  • " + w + "\n"
	}
	return strings.TrimRight(out, "\n")
}
