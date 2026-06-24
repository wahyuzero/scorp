package bootstrap

import (
	"scorp-agent/config"
	"scorp-agent/tools"
)

import "log"

func init() {
	// Initialize vault from tools package
	tools.Vault = &tools.CredentialVault{
		Path: config.ScorpPath("vault.json"),
	}
	tools.Vault.LoadMasterKey()
	tools.Vault.Load()
	log.Printf("[vault] Loaded %d credential entries", len(tools.Vault.Entries))
}