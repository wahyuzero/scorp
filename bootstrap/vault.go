package bootstrap

import (
	"log"
	"os"

	"scorp-agent/config"
	"scorp-agent/tools"
)

func init() {
	// Initialize vault from tools package
	tools.Vault = &tools.CredentialVault{
		Path: config.ScorpPath("vault.json"),
	}
	tools.Vault.LoadMasterKey()
	tools.Vault.Load()
	if os.Getenv("SCORP_DEBUG") != "" {
		log.Printf("[vault] Loaded %d credential entries", len(tools.Vault.Entries))
	}
}