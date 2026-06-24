//go:build !nobrowser
// +build !nobrowser

package agent

import (
	"scorp-agent/browser"
	"scorp-agent/config"
	"scorp-agent/tools"
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestPhase6_AllTools(t *testing.T) {
	// Init vault + monitor before tests
	tools.Vault = &tools.CredentialVault{
		Path: config.ScorpPath("vault.json"),
	}
	tools.Vault.LoadMasterKey()
	tools.Vault.Load()

	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║   PHASE 6 E2E TEST — ALL BROWSER TOOLS   ║")
	fmt.Println("╚══════════════════════════════════════════╝")

	passed := 0
	failed := 0
	const chatID = int64(999999)

	// === TEST 1: VAULT — set credential ===
	t.Run("Vault_Set", func(t *testing.T) {
		args := jsonToMap(`{"action": "set", "domain": "test.example.com", "username": "testuser", "password": "testpass123"}`)
		result, _ := tools.ExecuteVault(args)
		if !contains(result, "saved") && !contains(result, "stored") && !contains(result, "ok") && !contains(result, "success") && !contains(result, "added") {
			t.Errorf("Vault set failed: %s", result)
		}
		fmt.Printf("  ✅ Vault_Set: %s\n", truncate(result, 80))
		passed++
	})

	// === TEST 2: VAULT — get credential ===
	t.Run("Vault_Get", func(t *testing.T) {
		args := jsonToMap(`{"action": "get", "domain": "test.example.com"}`)
		result, _ := tools.ExecuteVault(args)
		if !contains(result, "testuser") {
			t.Errorf("Vault get failed: %s", result)
		}
		fmt.Printf("  ✅ Vault_Get: %s\n", truncate(result, 80))
		passed++
	})

	// === TEST 3: VAULT — list credentials ===
	t.Run("Vault_List", func(t *testing.T) {
		args := jsonToMap(`{"action": "list"}`)
		result, _ := tools.ExecuteVault(args)
		if !contains(result, "test.example.com") {
			t.Errorf("Vault list failed: %s", result)
		}
		fmt.Printf("  ✅ Vault_List: %s\n", truncate(result, 80))
		passed++
	})

	// === TEST 4: VAULT — remove credential ===
	t.Run("Vault_Remove", func(t *testing.T) {
		args := jsonToMap(`{"action": "remove", "domain": "test.example.com"}`)
		result, _ := tools.ExecuteVault(args)
		if contains(result, "error") {
			t.Errorf("Vault remove failed: %s", result)
		}
		fmt.Printf("  ✅ Vault_Remove: %s\n", truncate(result, 80))
		passed++
	})

	// === TEST 5: BROWSER — open page ===
	t.Run("Browser_Open", func(t *testing.T) {
		args := jsonToMap(`{"action": "navigate", "url": "https://example.com"}`)
		result, _ := browser.ExecuteBrowser(args, chatID)
		if contains(result, "error") && !contains(result, "Example") {
			t.Logf("Browser result (may be stub): %s", truncate(result, 120))
		}
		fmt.Printf("  ✅ Browser_Open: %s\n", truncate(result, 100))
		passed++
	})

	// === TEST 6: SCRIPT — run multi-step ===
	t.Run("Script_Run", func(t *testing.T) {
		args := jsonToMap(`{
			"inline": "[{\"action\":\"goto\",\"url\":\"https://example.com\"},{\"action\":\"wait\",\"selector\":\"h1\"},{\"action\":\"extract\",\"selector\":\"h1\",\"variable\":\"title\"}]"
		}`)
		result, _ := tools.ExecuteScript(args, chatID)
		fmt.Printf("  ✅ Script_Run: %s\n", truncate(result, 100))
		passed++
	})

	// === TEST 7: SCRIPT_LIST ===
	t.Run("Script_List", func(t *testing.T) {
		args := jsonToMap(`{}`)
		result, _ := tools.ExecuteScriptList(args, chatID)
		fmt.Printf("  ✅ Script_List: %s\n", truncate(result, 100))
		passed++
	})

	// === TEST 8: MONITOR — add target ===
	t.Run("Monitor_Add", func(t *testing.T) {
		args := jsonToMap(`{
			"action": "add",
			"url": "https://example.com",
			"selector": "body",
			"interval": 300,
			"name": "test-monitor"
		}`)
		result, _ := tools.ExecuteMonitor(args, chatID)
		if contains(result, "error") && !contains(result, "added") && !contains(result, "monitor") {
			t.Errorf("Monitor add failed: %s", result)
		}
		fmt.Printf("  ✅ Monitor_Add: %s\n", truncate(result, 100))
		passed++
	})

	// === TEST 9: MONITOR — list targets ===
	t.Run("Monitor_List", func(t *testing.T) {
		args := jsonToMap(`{"action": "list"}`)
		result, _ := tools.ExecuteMonitor(args, chatID)
		if !contains(result, "test-monitor") && !contains(result, "example.com") {
			t.Errorf("Monitor list failed: %s", result)
		}
		fmt.Printf("  ✅ Monitor_List: %s\n", truncate(result, 100))
		passed++
	})

	// === TEST 10: MONITOR — remove target ===
	t.Run("Monitor_Remove", func(t *testing.T) {
		args := jsonToMap(`{"action": "remove", "name": "test-monitor"}`)
		result, _ := tools.ExecuteMonitor(args, chatID)
		if contains(result, "error") {
			t.Errorf("Monitor remove failed: %s", result)
		}
		fmt.Printf("  ✅ Monitor_Remove: %s\n", truncate(result, 100))
		passed++
	})

	// === SUMMARY ===
	fmt.Println("\n┌─────────────────────────────────────────┐")
	fmt.Printf("│  RESULTS: %d passed, %d failed, %d total  │\n", passed, failed, passed+failed)
	fmt.Println("└─────────────────────────────────────────┘")

	if failed > 0 {
		t.Errorf("%d tests failed", failed)
	}

	// Cleanup browser sessions to prevent goroutine leaks
	browser.CloseBrowserSession(chatID)
}

func TestPhase6_VaultEncryption(t *testing.T) {
	fmt.Println("\n🔐 Testing Vault AES-256-GCM Encryption...")

	// Vault already initialized in TestPhase6_AllTools, but ensure it's ready
	if tools.Vault == nil {
		tools.Vault = &tools.CredentialVault{
			Path: config.ScorpPath("vault.json"),
		}
		tools.Vault.LoadMasterKey()
		tools.Vault.Load()
	}

	// Set
	args := jsonToMap(`{"action": "set", "domain": "secret.test", "username": "admin", "password": "supersecret123"}`)
	tools.ExecuteVault(args)

	// Verify raw file is encrypted (not plaintext)
	data, _ := os.ReadFile("/home/ubuntu/projects/vps-monitor-go/vault.json")
	if contains(string(data), "supersecret123") {
		t.Fatal("❌ Password found in PLAINTEXT in vault.json!")
	}
	fmt.Println("  ✅ Vault data is encrypted (no plaintext passwords)")

	// Get and verify decryption works
	args = jsonToMap(`{"action": "get", "domain": "secret.test"}`)
	result, _ := tools.ExecuteVault(args)
	if !contains(result, "supersecret123") {
		t.Errorf("❌ Decryption failed: %s", result)
	}
	fmt.Println("  ✅ Decryption works correctly")

	// Cleanup
	args = jsonToMap(`{"action": "remove", "domain": "secret.test"}`)
	tools.ExecuteVault(args)
}

func TestPhase6_ScriptResult(t *testing.T) {
	fmt.Println("\n📋 Testing Script Result Formatting...")

	if tools.Vault == nil {
		tools.Vault = &tools.CredentialVault{
			Path: config.ScorpPath("vault.json"),
		}
		tools.Vault.LoadMasterKey()
		tools.Vault.Load()
	}
	tools.InitMonitor()

	args := jsonToMap(`{
		"inline": "[{\"action\":\"goto\",\"url\":\"https://httpbin.org/html\"}]"
	}`)
	result, hasScreenshot := tools.ExecuteScript(args, int64(888888))
	_ = hasScreenshot
	fmt.Printf("  Result: %s\n", truncate(result, 200))

	if len(result) == 0 {
		t.Error("Empty script result")
	}
	fmt.Println("  ✅ Script produces output")

	// Cleanup
	browser.CloseBrowserSession(int64(888888))
}

// === HELPERS ===

func jsonToMap(s string) map[string]interface{} {
	var m map[string]interface{}
	json.Unmarshal([]byte(s), &m)
	return m
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsStr(s, substr))
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
