package models

import (
	"log"
	"os"
	"strings"
)

// ──────────────────────────────────────────────
// Provider Registry — Phase 1 Multi-Provider System
// ──────────────────────────────────────────────

// ProviderPreset defines a built-in provider with known endpoint + env vars.
type ProviderPreset struct {
	KeyEnvs       []string // env var names to try in order (first non-empty wins)
	BaseURL       string   // default API endpoint
	API           string   // "openai" | "anthropic" | "gemini"
	NoAuth        bool     // true for local providers (ollama)
	ExtraHeaders  bool     // true for openrouter (HTTP-Referer, X-Title)
	DisplayName   string   // human-readable name
}

// providerRegistry is the built-in provider registry.
// Users can reference these by name in models.json without specifying base_url.
var ProviderRegistry = map[string]ProviderPreset{
	"openai": {
		KeyEnvs:     []string{"OPENAI_API_KEY"},
		BaseURL:     "https://api.openai.com/v1",
		API:         "openai",
		DisplayName: "OpenAI",
	},
	"zai": {
		KeyEnvs:     []string{"GLM_API_KEY", "ZAI_API_KEY", "Z_AI_API_KEY"},
		BaseURL:     "https://api.z.ai/api/paas/v4",
		API:         "openai",
		DisplayName: "Z.AI (GLM) — Pay-per-Token",
	},
	"zai-coding": {
		KeyEnvs:     []string{"ZAI_CODING_API_KEY", "GLM_CODING_API_KEY", "GLM_API_KEY"},
		BaseURL:     "https://api.z.ai/api/coding/paas/v4",
		API:         "openai",
		DisplayName: "Z.AI Coding Plan ($18/mo)",
	},
	"deepseek": {
		KeyEnvs:     []string{"DEEPSEEK_API_KEY"},
		BaseURL:     "https://api.deepseek.com/v1",
		API:         "openai",
		DisplayName: "DeepSeek",
	},
	"groq": {
		KeyEnvs:     []string{"GROQ_API_KEY"},
		BaseURL:     "https://api.groq.com/openai/v1",
		API:         "openai",
		DisplayName: "Groq",
	},
	"openrouter": {
		KeyEnvs:      []string{"OPENROUTER_API_KEY"},
		BaseURL:      "https://openrouter.ai/api/v1",
		API:          "openai",
		ExtraHeaders: true,
		DisplayName:  "OpenRouter",
	},
	"gemini": {
		KeyEnvs:     []string{"GOOGLE_API_KEY", "GEMINI_API_KEY"},
		BaseURL:     "https://generativelanguage.googleapis.com/v1beta",
		API:         "gemini",
		DisplayName: "Google Gemini",
	},
	"anthropic": {
		KeyEnvs:     []string{"ANTHROPIC_API_KEY"},
		BaseURL:     "https://api.anthropic.com",
		API:         "anthropic",
		DisplayName: "Anthropic",
	},
	"copilot": {
		KeyEnvs:     []string{"COPILOT_GITHUB_TOKEN", "GH_TOKEN", "GITHUB_TOKEN"},
		BaseURL:     "https://models.inference.ai.azure.com",
		API:         "openai",
		DisplayName: "GitHub Copilot",
	},
	"kimi": {
		KeyEnvs:     []string{"KIMI_API_KEY", "MOONSHOT_API_KEY"},
		BaseURL:     "https://api.moonshot.ai/v1",
		API:         "openai",
		DisplayName: "Moonshot (Kimi)",
	},
	"minimax": {
		KeyEnvs:     []string{"MINIMAX_API_KEY"},
		BaseURL:     "https://api.minimax.io/anthropic",
		API:         "anthropic",
		DisplayName: "MiniMax",
	},
	"nvidia": {
		KeyEnvs:     []string{"NVIDIA_API_KEY"},
		BaseURL:     "https://integrate.api.nvidia.com/v1",
		API:         "openai",
		DisplayName: "NVIDIA NIM",
	},
	"huggingface": {
		KeyEnvs:     []string{"HF_TOKEN"},
		BaseURL:     "https://router.huggingface.co/v1",
		API:         "openai",
		DisplayName: "Hugging Face",
	},
	"ollama": {
		KeyEnvs:     []string{},
		BaseURL:     "http://127.0.0.1:11434/v1",
		API:         "openai",
		NoAuth:      true,
		DisplayName: "Ollama (local)",
	},
	"lmstudio": {
		KeyEnvs:     []string{"LM_API_KEY"},
		BaseURL:     "http://127.0.0.1:1234/v1",
		API:         "openai",
		NoAuth:      true,
		DisplayName: "LM Studio (local)",
	},
}

// resolveAPIKey resolves the API key using a 4-tier fallback:
// 1. key_env field (explicit env var name in config)
// 2. provider preset KeyEnvs (registry lookup)
// 3. SCORP_{PROVIDER}_API_KEY (generic pattern)
// 4. inline api_key (deprecated — logs warning)
func ResolveAPIKey(cfg *ModelConfig) string {
	if cfg == nil {
		return ""
	}

	// Tier 1: explicit key_env in config
	if cfg.KeyEnv != "" {
		if v := os.Getenv(cfg.KeyEnv); v != "" {
			return v
		}
		log.Printf("[models] WARNING: key_env '%s' set but env var is empty for provider %s", cfg.KeyEnv, cfg.Provider)
	}

	// Tier 2: provider preset registry lookup
	if preset, ok := ProviderRegistry[cfg.Provider]; ok {
		for _, envName := range preset.KeyEnvs {
			if v := os.Getenv(envName); v != "" {
				return v
			}
		}
	}

	// Tier 3: generic SCORP_{PROVIDER}_API_KEY pattern
	genericKey := "SCORP_" + strings.ToUpper(cfg.Provider) + "_API_KEY"
	if v := os.Getenv(genericKey); v != "" {
		return v
	}

	// Tier 4: inline api_key (deprecated)
	if cfg.APIKey != "" {
		log.Printf("[models] WARNING: using plaintext api_key for %s — migrate to key_env", cfg.Provider)
		return cfg.APIKey
	}

	return ""
}

// resolveBaseURL fills in base_url from the provider registry if not set in config.
func ResolveBaseURL(cfg *ModelConfig) string {
	if cfg == nil {
		return ""
	}
	if cfg.BaseURL != "" {
		return cfg.BaseURL
	}
	if preset, ok := ProviderRegistry[cfg.Provider]; ok {
		return preset.BaseURL
	}
	return ""
}

// resolveAPIFormat fills in the API format ("openai", "anthropic", "gemini").
func ResolveAPIFormat(cfg *ModelConfig) string {
	if cfg == nil {
		return "openai"
	}
	if cfg.API != "" {
		return cfg.API
	}
	if preset, ok := ProviderRegistry[cfg.Provider]; ok && preset.API != "" {
		return preset.API
	}
	return "openai"
}

// hasAPIKey checks whether a key can be resolved (for health display without leaking the key).
func hasAPIKey(cfg *ModelConfig) bool {
	return ResolveAPIKey(cfg) != ""
}

// keySourceLabel returns a human-readable description of where the key comes from.
// Used in /model list to show key status without revealing the key itself.
func KeySourceLabel(cfg *ModelConfig) string {
	if cfg == nil {
		return "no config"
	}

	// Check explicit key_env
	if cfg.KeyEnv != "" {
		if os.Getenv(cfg.KeyEnv) != "" {
			return "env:" + cfg.KeyEnv
		}
		return "⚠️ env:" + cfg.KeyEnv + " (empty)"
	}

	// Check provider preset
	if preset, ok := ProviderRegistry[cfg.Provider]; ok {
		for _, envName := range preset.KeyEnvs {
			if os.Getenv(envName) != "" {
				return "env:" + envName
			}
		}
		if preset.NoAuth {
			return "no auth"
		}
	}

	// Check generic pattern
	genericKey := "SCORP_" + strings.ToUpper(cfg.Provider) + "_API_KEY"
	if os.Getenv(genericKey) != "" {
		return "env:" + genericKey
	}

	// Inline key
	if cfg.APIKey != "" {
		return "⚠️ plaintext (deprecated)"
	}

	return "❌ no key"
}

// applyProviderDefaults fills in missing base_url and api fields from the registry.
// Called during config load to normalize configs that only specify provider name.
func applyProviderDefaults(cfg *ModelConfig) {
	if cfg == nil || cfg.Provider == "" {
		return
	}

	preset, ok := ProviderRegistry[cfg.Provider]
	if !ok {
		return // custom provider, nothing to fill
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = preset.BaseURL
	}
	if cfg.API == "" {
		cfg.API = preset.API
	}
}

// migrateModelConfigs auto-migrates plaintext api_key → key_env where possible.
// Logs warnings for each migration. Does NOT clear api_key if no preset match
// (keeps working, just logs warning at call time).
func migrateModelConfigs(cfg *ModelRouterConfig) {
	if cfg == nil {
		return
	}

	migrated := 0
	for name, m := range cfg.Models {
		// Fill defaults from registry
		applyProviderDefaults(&m)

		// If has plaintext key but no key_env, try to set key_env from preset
		if m.APIKey != "" && m.KeyEnv == "" {
			if preset, ok := ProviderRegistry[m.Provider]; ok && len(preset.KeyEnvs) > 0 {
				m.KeyEnv = preset.KeyEnvs[0]
				m.APIKey = "" // Clear plaintext
				log.Printf("[models] Migrated '%s': api_key → key_env=%s", name, m.KeyEnv)
				migrated++
			}
		}

		// Update the map entry (Go maps: need to write back)
		cfg.Models[name] = m
	}

	if migrated > 0 {
		log.Printf("[models] Auto-migrated %d model(s) from plaintext to key_env", migrated)
	}
}
