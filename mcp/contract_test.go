package mcp

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func contractTestServer(t *testing.T) (restore func()) {
	t.Helper()
	dir := t.TempDir()
	oldOverride := contractPathOverride
	SetContractPath(filepath.Join(dir, "contracts.json"))
	return func() { SetContractPath(oldOverride) }
}

func TestServerFingerprintDeterministicAndSensitive(t *testing.T) {
	toolsA := []MCPTool{
		{Name: "read", Description: "read a file", InputSchema: map[string]interface{}{"type": "object"}},
		{Name: "write", Description: "write a file", InputSchema: map[string]interface{}{"type": "object", "required": []interface{}{"path"}}},
	}
	fp1, _ := serverFingerprint(toolsA)
	fp2, _ := serverFingerprint([]MCPTool{toolsA[1], toolsA[0]})
	if fp1 != fp2 {
		t.Fatal("fingerprint must be order-independent")
	}

	changed := []MCPTool{toolsA[0], toolsA[1]}
	changed[1] = MCPTool{Name: "write", Description: "write a file", InputSchema: map[string]interface{}{"type": "object", "required": []interface{}{"path", "data"}}}
	fp3, _ := serverFingerprint(changed)
	if fp1 == fp3 {
		t.Fatal("schema change must change the fingerprint")
	}
}

func TestCheckContractsDetectsSilentChange(t *testing.T) {
	restore := contractTestServer(t)
	defer restore()

	toolsA := []MCPTool{{Name: "read", Description: "d", InputSchema: map[string]interface{}{}}}
	setTestServerTools("eval-srv", toolsA)
	defer deleteTestServerTools("eval-srv")

	// First check: baseline persisted, no warnings.
	if ws := CheckServerContracts(); len(ws) != 0 {
		t.Fatalf("baseline check must not warn: %v", ws)
	}
	data, err := os.ReadFile(contractFile())
	if err != nil {
		t.Fatalf("contract file not persisted: %v", err)
	}
	var stored map[string]ServerContract
	if err := json.Unmarshal(data, &stored); err != nil || stored["eval-srv"].Fingerprint == "" {
		t.Fatalf("contract not stored: %v %v", err, stored)
	}
	origFP := stored["eval-srv"].Fingerprint

	// Silent schema change → exactly one warning; new fingerprint persisted.
	toolsB := []MCPTool{{Name: "read", Description: "d", InputSchema: map[string]interface{}{"type": "object"}}}
	setTestServerTools("eval-srv", toolsB)
	ws := CheckServerContracts()
	if len(ws) != 1 || !strings.Contains(ws[0], "contract changed") {
		t.Fatalf("silent change must warn once, got %v", ws)
	}
	data, _ = os.ReadFile(contractFile())
	_ = json.Unmarshal(data, &stored)
	if stored["eval-srv"].Fingerprint == origFP {
		t.Fatal("new fingerprint must be persisted after the warning")
	}

	// Unchanged state → no new warnings.
	if ws := CheckServerContracts(); len(ws) != 0 {
		t.Fatalf("unchanged contract must not warn: %v", ws)
	}
}

func TestCheckContractsWarnsForVanishedServer(t *testing.T) {
	restore := contractTestServer(t)
	defer restore()

	toolsA := []MCPTool{{Name: "read", Description: "d", InputSchema: map[string]interface{}{}}}
	setTestServerTools("eval-gone", toolsA)
	CheckServerContracts() // baseline

	deleteTestServerTools("eval-gone")
	ws := CheckServerContracts()
	if len(ws) != 1 || !strings.Contains(ws[0], "no longer registered") {
		t.Fatalf("vanished server must warn, got %v", ws)
	}
}

// Direct registry manipulation helpers (same package).
func setTestServerTools(name string, tools []MCPTool) {
	mcpServersMu.Lock()
	defer mcpServersMu.Unlock()
	mcpServers[name] = &MCPServer{tools: tools}
}

func deleteTestServerTools(name string) {
	mcpServersMu.Lock()
	defer mcpServersMu.Unlock()
	delete(mcpServers, name)
}
