package tools

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

	"scorp-agent/config"
	"scorp-agent/internal/helpers"
)

// ──────────────────────────────────────────────
// Credential Vault — Encrypted credential storage (AES-GCM)
// ──────────────────────────────────────────────

// CredentialVault stores login credentials with AES-GCM encryption.
type CredentialVault struct {
	Path    string
	master  []byte // AES-256 key (32 bytes)
	Entries []CredentialEntry
	mu      sync.Mutex
}

// CredentialEntry holds credentials for one service.
type CredentialEntry struct {
	ID       string    `json:"id"`
	Domain   string    `json:"domain"`   // e.g. "mail.example.com"
	Username string    `json:"username"` // encrypted base64
	Password string    `json:"password"` // encrypted base64
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
}

var Vault *CredentialVault

func (v *CredentialVault) LoadMasterKey() {
	keyPath := config.ScorpPath("vault.key")
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

func (v *CredentialVault) Encrypt(plaintext string) string {
	if v.master == nil {
		return plaintext
	}
	block, err := aes.NewCipher(v.master)
	if err != nil {
		return plaintext
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return plaintext
	}
	nonce := make([]byte, gcm.NonceSize())
	rand.Read(nonce)
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext)
}

func (v *CredentialVault) Decrypt(cipherB64 string) string {
	if v.master == nil {
		return cipherB64
	}
	block, err := aes.NewCipher(v.master)
	if err != nil {
		return cipherB64
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return cipherB64
	}
	data, err := base64.StdEncoding.DecodeString(cipherB64)
	if err != nil || len(data) < gcm.NonceSize() {
		return cipherB64
	}
	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return cipherB64
	}
	return string(plain)
}

func (v *CredentialVault) Load() {
	v.mu.Lock()
	defer v.mu.Unlock()
	if data, err := os.ReadFile(v.Path); err == nil {
		json.Unmarshal(data, &v.Entries)
	}
	if v.Entries == nil {
		v.Entries = []CredentialEntry{}
	}
}

func (v *CredentialVault) Persist() {
	v.mu.Lock()
	defer v.mu.Unlock()
	data, _ := json.MarshalIndent(v.Entries, "", "  ")
	os.MkdirAll(filepath.Dir(v.Path), 0700)
	os.WriteFile(v.Path, data, 0600)
}

func (v *CredentialVault) Add(domain, username, password string) {
	v.mu.Lock()
	defer v.mu.Unlock()

	// Check existing
	for i, e := range v.Entries {
		if e.Domain == domain && e.Username == username {
			v.Entries[i].Password = v.Encrypt(password)
			v.Entries[i].Updated = time.Now()
			return
		}
	}

	v.Entries = append(v.Entries, CredentialEntry{
		ID:       fmt.Sprintf("cred_%d", time.Now().UnixNano()),
		Domain:   domain,
		Username: v.Encrypt(username),
		Password: v.Encrypt(password),
		Created:  time.Now(),
		Updated:  time.Now(),
	})
}

func (v *CredentialVault) Get(domain string) (username, password string, ok bool) {
	v.mu.Lock()
	defer v.mu.Unlock()
	for _, e := range v.Entries {
		if e.Domain == domain {
			return v.Decrypt(e.Username), v.Decrypt(e.Password), true
		}
	}
	return "", "", false
}

// ExecuteVault handles vault_get and vault_set operations.
func ExecuteVault(args map[string]interface{}) (string, bool) {
	action := helpers.GetStringArg(args, "action", "")
	domain := helpers.GetStringArg(args, "domain", "")
	username := helpers.GetStringArg(args, "username", "")
	password := helpers.GetStringArg(args, "password", "")

	switch action {
	case "get":
		if domain == "" {
			return "Error: domain required for vault get", false
		}
		u, p, ok := Vault.Get(domain)
		if !ok {
			return fmt.Sprintf("No credentials found for domain: %s", domain), false
		}
		return fmt.Sprintf("🔐 Domain: %s\nUsername: %s\nPassword: %s\n", domain, u, p), true
	case "set":
		if domain == "" || username == "" || password == "" {
			return "Error: domain, username, and password required for vault set", false
		}
		Vault.Add(domain, username, password)
		Vault.Persist()
		return fmt.Sprintf("✅ Credentials saved for %s [encrypted]", domain), true
	case "list":
		Vault.mu.Lock()
		defer Vault.mu.Unlock()
		if len(Vault.Entries) == 0 {
			return "No credentials stored.", true
		}
		var result string
		for _, e := range Vault.Entries {
			result += fmt.Sprintf("- %s (user: %s, updated: %s)\n", e.Domain, Vault.Decrypt(e.Username), e.Updated.Format("02 Jan 15:04"))
		}
		return result, true
	case "remove":
		if domain == "" {
			return "Error: domain required", false
		}
		Vault.mu.Lock()
		for i, e := range Vault.Entries {
			if e.Domain == domain {
				Vault.Entries = append(Vault.Entries[:i], Vault.Entries[i+1:]...)
				break
			}
		}
		Vault.mu.Unlock()
		Vault.Persist()
		return fmt.Sprintf("✅ Removed credentials for %s", domain), true
	default:
		return "Unknown vault action (use get, set, list, remove)", false
	}
}
