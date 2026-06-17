// +build !nobrowser

package main

import "log"

func init() {
	// ── Vault (credential store) ──
	registerTool(ToolDef{
		Name:        "vault",
		Description: "Encrypted credential store. Actions: get (retrieve credentials by domain), set (save credentials), list (show saved domains), remove (delete). Credentials are encrypted at rest with AES-256-GCM.",
		Category:    "browser",
		Native:      true,
		Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
			return executeVault(args)
		},
		Arguments: map[string]ArgDef{
			"action":   {Type: "string", Description: "get, set, list, remove", Required: true},
			"domain":   {Type: "string", Description: "Domain name (e.g. mail.frugaldev.biz.id)"},
			"username": {Type: "string", Description: "Username/email (for set)"},
			"password": {Type: "string", Description: "Password (for set)"},
		},
	})

	// ── Autologin ──
	registerTool(ToolDef{
		Name:        "autologin",
		Description: "Auto-detect login forms on the current browser page and fill in stored credentials. Detects username/email/password fields and auto-fills from vault. Requires browser session to be on the login page first.",
		Category:    "browser",
		Native:      true,
		Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
			return executeAutoLogin(args, chatID)
		},
		Arguments: map[string]ArgDef{
			"domain": {Type: "string", Description: "Domain to look up credentials (auto-extracted from url if omitted)"},
			"url":    {Type: "string", Description: "URL of login page (optional)"},
		},
	})
}

func initVault() {
	vault = &CredentialVault{
		Path: scorpPath("vault.json"),
	}
	vault.loadMasterKey()
	vault.load()
	log.Printf("[vault] Loaded %d credential entries", len(vault.entries))
}
