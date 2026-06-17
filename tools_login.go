// +build !nobrowser

package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
)

// ──────────────────────────────────────────────
// Credential Vault — Encrypted credential storage
// ──────────────────────────────────────────────

// CredentialVault stores login credentials with AES-GCM encryption.
type CredentialVault struct {
	Path    string
	master  []byte // AES-256 key (32 bytes)
	entries []CredentialEntry
	mu      sync.Mutex
}

// CredentialEntry holds credentials for one service.
type CredentialEntry struct {
	ID         string    `json:"id"`
	Domain     string    `json:"domain"`     // e.g. "mail.frugaldev.biz.id"
	Username   string    `json:"username"`   // encrypted base64
	Password   string    `json:"password"`   // encrypted base64
	TOTPSecret string    `json:"totp_secret"` // encrypted base64
	Created    time.Time `json:"created"`
	Updated    time.Time `json:"updated"`
}

var vault *CredentialVault
func (v *CredentialVault) loadMasterKey() {
	keyPath := scorpPath("vault.key")
	if data, err := os.ReadFile(keyPath); err == nil {
		v.master = data
		return
	}
	// Generate new key
	key := make([]byte, 32)
	rand.Read(key)
	v.master = key
	os.WriteFile(keyPath, key, 0600)
	log.Printf("[vault] Generated new master key")
}

func (v *CredentialVault) encrypt(plaintext string) string {
	if v.master == nil {
		return plaintext
	}
	block, _ := aes.NewCipher(v.master)
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	rand.Read(nonce)
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext)
}

func (v *CredentialVault) decrypt(cipherB64 string) string {
	if v.master == nil {
		return cipherB64
	}
	block, _ := aes.NewCipher(v.master)
	gcm, _ := cipher.NewGCM(block)
	data, _ := base64.StdEncoding.DecodeString(cipherB64)
	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	plain, _ := gcm.Open(nil, nonce, ciphertext, nil)
	return string(plain)
}

func (v *CredentialVault) load() {
	v.mu.Lock()
	defer v.mu.Unlock()
	if data, err := os.ReadFile(v.Path); err == nil {
		json.Unmarshal(data, &v.entries)
	}
	if v.entries == nil {
		v.entries = []CredentialEntry{}
	}
}

func (v *CredentialVault) persist() {
	v.mu.Lock()
	defer v.mu.Unlock()
	data, _ := json.MarshalIndent(v.entries, "", "  ")
	os.MkdirAll(filepath.Dir(v.Path), 0700)
	os.WriteFile(v.Path, data, 0600)
}

func (v *CredentialVault) add(domain, username, password string) {
	v.mu.Lock()
	defer v.mu.Unlock()

	// Check existing
	for i, e := range v.entries {
		if e.Domain == domain && e.Username == username {
			v.entries[i].Password = v.encrypt(password)
			v.entries[i].Updated = time.Now()
			return
		}
	}

	v.entries = append(v.entries, CredentialEntry{
		ID:       fmt.Sprintf("cred_%d", time.Now().UnixNano()),
		Domain:   domain,
		Username: v.encrypt(username),
		Password: v.encrypt(password),
		Created:  time.Now(),
		Updated:  time.Now(),
	})
}

func (v *CredentialVault) get(domain string) (username, password string, ok bool) {
	v.mu.Lock()
	defer v.mu.Unlock()
	for _, e := range v.entries {
		if e.Domain == domain {
			return v.decrypt(e.Username), v.decrypt(e.Password), true
		}
	}
	return "", "", false
}

// ──────────────────────────────────────────────
// Login Tools (exposed to agent)
// ──────────────────────────────────────────────

// executeCredential handles vault_get and vault_set operations.
func executeVault(args map[string]interface{}) (string, bool) {
	action := getStringArg(args, "action", "")
	domain := getStringArg(args, "domain", "")
	username := getStringArg(args, "username", "")
	password := getStringArg(args, "password", "")

	switch action {
	case "get":
		if domain == "" {
			return "Error: domain required for vault get", false
		}
		u, p, ok := vault.get(domain)
		if !ok {
			return fmt.Sprintf("No credentials found for domain: %s", domain), false
		}
		return fmt.Sprintf("🔐 Domain: %s\nUsername: %s\nPassword: %s\n", domain, u, p), true
	case "set":
		if domain == "" || username == "" || password == "" {
			return "Error: domain, username, and password required for vault set", false
		}
		vault.add(domain, username, password)
		vault.persist()
		return fmt.Sprintf("✅ Credentials saved for %s [encrypted]", domain), true
	case "list":
		vault.mu.Lock()
		defer vault.mu.Unlock()
		if len(vault.entries) == 0 {
			return "No credentials stored.", true
		}
		var result string
		for _, e := range vault.entries {
			result += fmt.Sprintf("- %s (user: %s, updated: %s)\n", e.Domain, vault.decrypt(e.Username), e.Updated.Format("02 Jan 15:04"))
		}
		return result, true
	case "remove":
		if domain == "" {
			return "Error: domain required", false
		}
		vault.mu.Lock()
		for i, e := range vault.entries {
			if e.Domain == domain {
				vault.entries = append(vault.entries[:i], vault.entries[i+1:]...)
				break
			}
		}
		vault.mu.Unlock()
		vault.persist()
		return fmt.Sprintf("✅ Removed credentials for %s", domain), true
	default:
		return "Unknown vault action (use get, set, list, remove)", false
	}
}

// executeAutoLogin detects login forms and auto-fills credentials.
func executeAutoLogin(args map[string]interface{}, chatID int64) (string, bool) {
	domain := getStringArg(args, "domain", "")
	if domain == "" && len(args) > 0 {
		// Try to extract domain from URL arg
		url := getStringArg(args, "url", "")
		if url == "" {
			return "Error: domain or url required", false
		}
		domain = extractDomain(url)
	}

	username, password, ok := vault.get(domain)
	if !ok {
		return fmt.Sprintf("No stored credentials for %s. Use vault_set to add them.", domain), false
	}

	sess := getOrCreateBrowserSession(chatID)

	log.Printf("[autologin] Attempting login for %s with user %s", domain, username)

	// Auto-fill common login form selectors
	selectors := []string{
		"input[name='username']",
		"input[name='email']",
		"input[name='login']",
		"input[name='user']",
		"input[type='email']",
		"#username",
		"#email",
		"#user",
		"#login",
	}

	// Check if we're on a login page (look for password field first)
	formFields := make(map[string]string)
	var detectScript = `
		(function() {
			var fields = {};
			var pwdInput = document.querySelector('input[type="password"]');
			if (pwdInput) fields.password = pwdInput;
			var inputs = document.querySelectorAll('input');
			for (var i = 0; i < inputs.length; i++) {
				if (inputs[i].type === 'password') continue;
				var name = inputs[i].name || inputs[i].id || 'field' + i;
				if (inputs[i].type === 'text' || inputs[i].type === 'email' || !inputs[i].type) {
					if (!fields.username) fields.username = inputs[i];
					break;
				}
			}
			return JSON.stringify({hasPassword: !!pwdInput, fieldCount: Object.keys(fields).length});
		})();
	`

	var detectResult string
	err := chromedp.Run(sess.Ctx, chromedp.Evaluate(detectScript, &detectResult))
	if err != nil {
		return fmt.Sprintf("Error detecting login form: %v", err), false
	}

	if detectResult == "" || detectResult == `{"hasPassword":false,"fieldCount":0}` {
		return fmt.Sprintf("No login form detected on current page for %s.", domain), false
	}

	// Try typing username into each potential field
	for _, sel := range selectors {
		var found bool
		checkScript := fmt.Sprintf(`document.querySelector('%s') !== null`, sel)
		chromedp.Run(sess.Ctx, chromedp.Evaluate(checkScript, &found))
		if found {
			chromedp.Run(sess.Ctx,
				chromedp.Clear(sel, chromedp.ByQuery),
				chromedp.SendKeys(sel, username, chromedp.ByQuery),
			)
			log.Printf("[autologin] Filled username using selector: %s", sel)
			formFields["username"] = sel
			break
		}
	}

	// Fill password
	chromedp.Run(sess.Ctx,
		chromedp.Clear("input[type='password']", chromedp.ByQuery),
		chromedp.SendKeys("input[type='password']", password, chromedp.ByQuery),
	)

	return fmt.Sprintf("🔐 Auto-filled login form for %s (user: %s). Submit manually via browser action=click or check + submit.", domain, username), true
}

// extractDomain extracts the domain from a URL string.
func extractDomain(url string) string {
	// Simple domain extraction
	domain := url
	if len(domain) > 7 && domain[:7] == "http://" {
		domain = domain[7:]
	}
	if len(domain) > 8 && domain[:8] == "https://" {
		domain = domain[8:]
	}
	for i := 0; i < len(domain); i++ {
		if domain[i] == '/' || domain[i] == '?' || domain[i] == ':' {
			domain = domain[:i]
			break
		}
	}
	return domain
}
