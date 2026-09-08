package models

import (
	"log"
	"os"
	"strings"
)

// ──────────────────────────────────────────────
// Provider Registry — Multi-Provider System
// Rich pluggable provider ecosystem
// ──────────────────────────────────────────────

// ProviderPreset defines a built-in provider with known endpoint + env vars.
type ProviderPreset struct {
	KeyEnvs      []string // env var names to try in order (first non-empty wins)
	BaseURL      string   // default API endpoint
	API          string   // "openai" | "anthropic" | "gemini" | custom format
	NoAuth       bool     // true for local providers (ollama)
	ExtraHeaders bool     // true for openrouter (HTTP-Referer, X-Title)
	DisplayName  string   // human-readable name
}

// ProviderRegistry holds provider presets.
// Can be extended dynamically via RegisterProvider or RegisterProviderPreset.
var ProviderRegistry = map[string]ProviderPreset{
	"command-code": {
		KeyEnvs:     []string{"COMMAND_CODE_API_KEY", "COMMANDCODE_API_KEY"},
		BaseURL:     "https://api.commandcode.ai",
		API:         "command-code",
		DisplayName: "Command Code (CLI Gateway)",
	},
	"commandcode": {
		KeyEnvs:     []string{"COMMAND_CODE_API_KEY", "COMMANDCODE_API_KEY"},
		BaseURL:     "https://api.commandcode.ai",
		API:         "command-code",
		DisplayName: "Command Code (CLI Gateway)",
	},
	"opencode": {
		KeyEnvs:     []string{"OPENCODE_API_KEY", "OPENCODE_ZEN_API_KEY"},
		BaseURL:     "https://opencode.ai/zen/v1",
		API:         "opencode",
		DisplayName: "OpenCode Zen (Free AI Gateway)",
	},
	"opencode-zen": {
		KeyEnvs:     []string{"OPENCODE_API_KEY", "OPENCODE_ZEN_API_KEY"},
		BaseURL:     "https://opencode.ai/zen/v1",
		API:         "opencode",
		DisplayName: "OpenCode Zen (Free AI Gateway)",
	},
	"opencode-free": {
		KeyEnvs:     []string{"OPENCODE_API_KEY", "OPENCODE_ZEN_API_KEY"},
		BaseURL:     "https://opencode.ai/zen/v1",
		API:         "opencode",
		DisplayName: "OpenCode Zen Free",
	},
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
	"mistral": {
		KeyEnvs:     []string{"MISTRAL_API_KEY"},
		BaseURL:     "https://api.mistral.ai/v1",
		API:         "openai",
		DisplayName: "Mistral AI",
	},
	"siliconflow": {
		KeyEnvs:     []string{"SILICONFLOW_API_KEY"},
		BaseURL:     "https://api.siliconflow.cn/v1",
		API:         "openai",
		DisplayName: "SiliconFlow",
	},
	"cerebras": {
		KeyEnvs:     []string{"CEREBRAS_API_KEY"},
		BaseURL:     "https://api.cerebras.ai/v1",
		API:         "openai",
		DisplayName: "Cerebras",
	},
	"novita": {
		KeyEnvs:     []string{"NOVITA_API_KEY"},
		BaseURL:     "https://api.novita.ai/v3/openai",
		API:         "openai",
		DisplayName: "Novita AI",
	},
	"together": {
		KeyEnvs:     []string{"TOGETHER_API_KEY"},
		BaseURL:     "https://api.together.xyz/v1",
		API:         "openai",
		DisplayName: "Together AI",
	},
	"fireworks": {
		KeyEnvs:     []string{"FIREWORKS_API_KEY"},
		BaseURL:     "https://api.fireworks.ai/inference/v1",
		API:         "openai",
		DisplayName: "Fireworks AI",
	},
	"volcengine": {
		KeyEnvs:     []string{"VOLCENGINE_API_KEY", "ARK_API_KEY"},
		BaseURL:     "https://ark.cn-beijing.volces.com/api/v3",
		API:         "openai",
		DisplayName: "Volcengine Ark",
	},
	"modelscope": {
		KeyEnvs:     []string{"MODELSCOPE_API_KEY"},
		BaseURL:     "https://api-inference.modelscope.cn/v1",
		API:         "openai",
		DisplayName: "ModelScope",
	},
	"qwen": {
		KeyEnvs:     []string{"DASHSCOPE_API_KEY"},
		BaseURL:     "https://dashscope.aliyuncs.com/compatible-mode/v1",
		API:         "openai",
		DisplayName: "Alibaba Qwen",
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

// RegisterProviderPreset adds or updates a provider preset in the registry.
func RegisterProviderPreset(name string, preset ProviderPreset) {
	providersMu.Lock()
	defer providersMu.Unlock()
	ProviderRegistry[strings.ToLower(strings.TrimSpace(name))] = preset
}

// ResolveAPIKey resolves the API key using a multi-tier fallback:
// 1. key_env field (explicit env var name in config)
// 2. KeyResolver interface on registered provider adapter (disk/cli/oauth)
// 3. provider preset KeyEnvs (registry lookup)
// 4. generic SCORP_{PROVIDER}_API_KEY pattern
// 5. inline api_key (deprecated — logs warning)
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

	// Tier 2: Dynamic KeyResolver on provider adapter
	format := ResolveAPIFormat(cfg)
	adapter := GetProvider(format)
	if resolver, ok := adapter.(KeyResolver); ok {
		if key := resolver.ResolveKey(cfg); key != "" {
			return key
		}
	}
	if cfg.Provider != "" && cfg.Provider != format {
		if directAdapter := GetProvider(cfg.Provider); directAdapter != nil && directAdapter != adapter {
			if resolver, ok := directAdapter.(KeyResolver); ok {
				if key := resolver.ResolveKey(cfg); key != "" {
					return key
				}
			}
		}
	}

	// Tier 3: provider preset registry lookup
	providersMu.RLock()
	preset, hasPreset := ProviderRegistry[cfg.Provider]
	providersMu.RUnlock()
	if hasPreset {
		for _, envName := range preset.KeyEnvs {
			if v := os.Getenv(envName); v != "" {
				return v
			}
		}
	}

	// Tier 4: generic SCORP_{PROVIDER}_API_KEY pattern
	cleanProvider := strings.ReplaceAll(cfg.Provider, "-", "_")
	genericKey := "SCORP_" + strings.ToUpper(cleanProvider) + "_API_KEY"
	if v := os.Getenv(genericKey); v != "" {
		return v
	}
	if legacyKey := "SCORP_" + strings.ToUpper(cfg.Provider) + "_API_KEY"; legacyKey != genericKey {
		if v := os.Getenv(legacyKey); v != "" {
			return v
		}
	}

	// Tier 5: inline api_key (deprecated)
	if cfg.APIKey != "" {
		log.Printf("[models] WARNING: using plaintext api_key for %s — migrate to key_env", cfg.Provider)
		return cfg.APIKey
	}

	return ""
}

// ResolveBaseURL fills in base_url from the provider registry if not set in config.
func ResolveBaseURL(cfg *ModelConfig) string {
	if cfg == nil {
		return ""
	}
	if cfg.BaseURL != "" {
		return cfg.BaseURL
	}
	providersMu.RLock()
	defer providersMu.RUnlock()
	if preset, ok := ProviderRegistry[cfg.Provider]; ok {
		return preset.BaseURL
	}
	return ""
}

// ResolveAPIFormat fills in the API format ("openai", "anthropic", "gemini", or custom).
// Automatically checks dynamically registered provider adapters first.
func ResolveAPIFormat(cfg *ModelConfig) string {
	if cfg == nil {
		return "openai"
	}
	// Auto-route OpenCode provider / endpoint to dedicated opencode adapter
	if cfg.Provider == "opencode" || cfg.Provider == "opencode-zen" || cfg.Provider == "opencode-free" || strings.Contains(cfg.BaseURL, "opencode.ai") {
		if cfg.API == "" || cfg.API == "openai" || cfg.API == "opencode" {
			return "opencode"
		}
	}
	if cfg.API != "" {
		return cfg.API
	}

	providersMu.RLock()
	defer providersMu.RUnlock()

	// Check if provider name directly matches a registered adapter
	providerKey := strings.ToLower(strings.TrimSpace(cfg.Provider))
	if _, ok := providerAdapters[providerKey]; ok {
		return providerKey
	}

	// Check provider preset registry
	if preset, ok := ProviderRegistry[cfg.Provider]; ok && preset.API != "" {
		return preset.API
	}

	return "openai"
}

// hasAPIKey checks whether a key can be resolved (for health display without leaking the key).
func hasAPIKey(cfg *ModelConfig) bool {
	return ResolveAPIKey(cfg) != ""
}

// KeySourceLabel returns a human-readable description of where the key comes from.
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

	// Check dynamic KeyResolver
	format := ResolveAPIFormat(cfg)
	adapter := GetProvider(format)
	if resolver, ok := adapter.(KeyResolver); ok {
		if key := resolver.ResolveKey(cfg); key != "" {
			return "disk/auth-store"
		}
	}
	if cfg.Provider != "" && cfg.Provider != format {
		if directAdapter := GetProvider(cfg.Provider); directAdapter != nil && directAdapter != adapter {
			if resolver, ok := directAdapter.(KeyResolver); ok {
				if key := resolver.ResolveKey(cfg); key != "" {
					return "disk/auth-store"
				}
			}
		}
	}

	// Check provider preset
	providersMu.RLock()
	preset, ok := ProviderRegistry[cfg.Provider]
	providersMu.RUnlock()
	if ok {
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
	cleanProvider := strings.ReplaceAll(cfg.Provider, "-", "_")
	genericKey := "SCORP_" + strings.ToUpper(cleanProvider) + "_API_KEY"
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

	providersMu.RLock()
	preset, ok := ProviderRegistry[cfg.Provider]
	providersMu.RUnlock()
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
func migrateModelConfigs(cfg *ModelRouterConfig) {
	if cfg == nil {
		return
	}

	migrated := 0
	for name, m := range cfg.Models {
		applyProviderDefaults(&m)

		if m.APIKey != "" && m.KeyEnv == "" {
			providersMu.RLock()
			preset, ok := ProviderRegistry[m.Provider]
			providersMu.RUnlock()
			if ok && len(preset.KeyEnvs) > 0 {
				m.KeyEnv = preset.KeyEnvs[0]
				m.APIKey = ""
				log.Printf("[models] Migrated '%s': api_key → key_env=%s", name, m.KeyEnv)
				migrated++
			}
		}

		cfg.Models[name] = m
	}

	if migrated > 0 {
		log.Printf("[models] Auto-migrated %d plaintext API key(s) to key_env", migrated)
	}
}
