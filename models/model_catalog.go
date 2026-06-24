package models

import (
	"log"
	"os"
	"sort"
	"strings"
)

// ──────────────────────────────────────────────
// Model Catalog — predefined models per provider
// When user adds an API key, models auto-populate from here
// ──────────────────────────────────────────────

// CatalogEntry defines a single model in a provider's catalog
type CatalogEntry struct {
	ModelID    string // API model identifier
	MaxTokens  int    // max output tokens
	Premium    bool   // true = high-capability, expensive
	Alias      string // suggested friendly name (optional)
}

// providerCatalog holds the model catalog for each provider.
// Providers not listed here (ollama, lmstudio, 9router, custom) have no catalog
// — users must add models manually.
var providerCatalog = map[string][]CatalogEntry{
	"openai": {
		{"gpt-4o", 16384, true, "gpt4o"},
		{"gpt-4o-mini", 16384, false, "gpt4o-mini"},
		{"o3-mini", 16384, true, "o3-mini"},
		{"gpt-4-turbo", 4096, true, ""},
		{"gpt-3.5-turbo", 4096, false, ""},
	},
	"zai": {
		{"glm-4.6", 4096, true, ""},
		{"glm-4.5", 4096, true, ""},
		{"glm-4-plus", 4096, true, ""},
		{"glm-4-air", 4096, false, ""},
		{"glm-4-airx", 4096, false, ""},
		{"glm-4-flash", 4096, false, ""},
		{"glm-4-flashx", 4096, false, ""},
		{"glm-4-long", 4096, false, ""},
	},
	"zai-coding": {
		{"glm-5.2", 4096, true, ""},
		{"glm-5-turbo", 4096, true, ""},
		{"glm-4.7", 4096, false, ""},
	},
	"deepseek": {
		{"deepseek-chat", 8192, false, ""},
		{"deepseek-reasoner", 8192, true, ""},
		{"deepseek-coder", 8192, false, ""},
	},
	"anthropic": {
		{"claude-sonnet-4-20250514", 8192, true, "claude-sonnet4"},
		{"claude-3-5-sonnet-20241022", 8192, true, ""},
		{"claude-3-5-haiku-20241022", 8192, false, ""},
		{"claude-3-opus-20240229", 4096, true, ""},
	},
	"gemini": {
		{"gemini-2.5-flash", 8192, false, ""},
		{"gemini-2.0-flash", 8192, false, ""},
		{"gemini-2.0-flash-thinking-exp", 8192, true, ""},
		{"gemini-1.5-pro", 8192, true, ""},
		{"gemini-1.5-flash", 8192, false, ""},
		{"gemini-1.5-flash-8b", 8192, false, ""},
	},
	"groq": {
		{"llama-3.3-70b-versatile", 32768, false, ""},
		{"llama-3.1-8b-instant", 32768, false, ""},
		{"mixtral-8x7b-32768", 32768, false, ""},
		{"gemma2-9b-it", 8192, false, ""},
	},
	"openrouter": {
		{"anthropic/claude-3.5-sonnet", 8192, true, ""},
		{"openai/gpt-4o", 16384, true, ""},
		{"google/gemini-2.0-flash-exp:free", 8192, false, ""},
		{"meta-llama/llama-3.3-70b-instruct", 32768, false, ""},
		{"deepseek/deepseek-r1", 8192, true, ""},
	},
	"kimi": {
		{"moonshot-v1-8k", 8192, false, ""},
		{"moonshot-v1-32k", 32768, false, ""},
		{"moonshot-v1-128k", 131072, false, ""},
	},
	"mistral": {
		{"mistral-large-latest", 8192, true, ""},
		{"mistral-small-latest", 8192, false, ""},
		{"codestral-latest", 8192, false, ""},
		{"open-mistral-7b", 8192, false, ""},
	},
	"cohere": {
		{"command-r-plus", 4096, true, ""},
		{"command-r", 4096, false, ""},
		{"command-r7b-12-2024", 4096, false, ""},
	},
	"together": {
		{"meta-llama/Llama-3.3-70B-Instruct-Turbo", 4096, false, ""},
		{"meta-llama/Meta-Llama-3.1-8B-Instruct-Turbo", 4096, false, ""},
		{"Qwen/Qwen2.5-72B-Instruct-Turbo", 4096, true, ""},
	},
	"fireworks": {
		{"accounts/fireworks/models/llama-v3p1-70b-instruct", 4096, false, ""},
		{"accounts/fireworks/models/llama-v3p1-8b-instruct", 4096, false, ""},
		{"accounts/fireworks/models/qwen2p5-72b-instruct", 4096, true, ""},
	},
	"nvidia": {
		{"meta/llama-3.1-70b-instruct", 4096, false, ""},
		{"meta/llama-3.1-8b-instruct", 4096, false, ""},
		{"nvidia/llama-3.1-nemotron-70b-instruct", 4096, true, ""},
	},
	"minimax": {
		{"MiniMax-Text-01", 8192, false, ""},
	},
	"huggingface": {
		{"meta-llama/Llama-3.3-70B-Instruct", 4096, true, ""},
		{"Qwen/Qwen2.5-72B-Instruct", 4096, true, ""},
	},
}

// hasCatalog returns true if the provider has a predefined model catalog
func HasCatalog(provider string) bool {
	_, ok := providerCatalog[provider]
	return ok
}

// catalogModels returns the catalog entries for a provider
func CatalogModels(provider string) []CatalogEntry {
	return providerCatalog[provider]
}

// providerHasAPIKey checks if a provider has any API key configured
func ProviderHasAPIKey(provider string) bool {
	preset, ok := ProviderRegistry[provider]
	if ok {
		if preset.NoAuth {
			return true
		}
		for _, envName := range preset.KeyEnvs {
			if os.Getenv(envName) != "" {
				return true
			}
		}
	}
	// Check custom
	ModelCfgMu.RLock()
	if cp, ok := ModelCfg.CustomProviders[provider]; ok {
		for _, envName := range cp.KeyEnvs {
			if os.Getenv(envName) != "" {
				return true
			}
		}
	}
	ModelCfgMu.RUnlock()
	// Check generic
	genericKey := "SCORP_" + strings.ToUpper(provider) + "_API_KEY"
	return os.Getenv(genericKey) != ""
}

// providerKeyEnv returns the primary key env var for a provider
func ProviderKeyEnv(provider string) string {
	preset, ok := ProviderRegistry[provider]
	if ok && len(preset.KeyEnvs) > 0 {
		return preset.KeyEnvs[0]
	}
	ModelCfgMu.RLock()
	defer ModelCfgMu.RUnlock()
	if cp, ok := ModelCfg.CustomProviders[provider]; ok && len(cp.KeyEnvs) > 0 {
		return cp.KeyEnvs[0]
	}
	return "SCORP_" + strings.ToUpper(provider) + "_API_KEY"
}

// modelsForProvider returns all model names in config that belong to a provider
func ModelsForProvider(provider string) []string {
	ModelCfgMu.RLock()
	defer ModelCfgMu.RUnlock()
	var names []string
	for name, m := range ModelCfg.Models {
		if m.Provider == provider {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	return names
}

// autoPopulateFromCatalog adds all catalog models for a provider to the config.
// Returns the list of model names that were added (or already existed).
// If apiKey is non-empty, saves it to env.
func AutoPopulateFromCatalog(provider, apiKey string) []string {
	entries, ok := providerCatalog[provider]
	if !ok {
		return nil
	}

	ModelCfgMu.Lock()
	defer ModelCfgMu.Unlock()

	// Save API key if provided
	keyEnv := ProviderKeyEnv(provider)
	if apiKey != "" {
		// Defer to updateEnvFile outside the lock
		go func() {
			UpdateEnvFile(keyEnv, apiKey)
			log.Printf("[models] API key saved for provider %s as %s", provider, keyEnv)
		}()
	}

	preset, hasPreset := ProviderRegistry[provider]
	baseURL := ""
	apiFormat := "openai"
	if hasPreset {
		baseURL = preset.BaseURL
		apiFormat = preset.API
	}

	// Check custom provider for base URL/API
	if cp, ok := ModelCfg.CustomProviders[provider]; ok {
		baseURL = cp.BaseURL
		apiFormat = cp.API
	}

	var added []string
	for _, entry := range entries {
		// Generate friendly name
		friendlyName := entry.Alias
		if friendlyName == "" {
			friendlyName = entry.ModelID
			// Clean up: remove slashes, dots → dashes
			friendlyName = strings.ReplaceAll(friendlyName, "/", "-")
			friendlyName = strings.ReplaceAll(friendlyName, ":", "-")
		}
		// Avoid collision: if name taken by different provider's model, append provider prefix
		if existing, ok := ModelCfg.Models[friendlyName]; ok && existing.Provider != provider {
			friendlyName = provider + "-" + friendlyName
		}

		mc := ModelConfig{
			Provider:  provider,
			Model:     entry.ModelID,
			BaseURL:   baseURL,
			MaxTokens: entry.MaxTokens,
			KeyEnv:    keyEnv,
			API:       apiFormat,
		}

		if _, exists := ModelCfg.Models[friendlyName]; !exists {
			ModelCfg.Models[friendlyName] = mc
			added = append(added, friendlyName)
			log.Printf("[models] Auto-populated: %s → %s (%s)", friendlyName, entry.ModelID, provider)
		}
	}

	// If no default model set, pick the first non-premium one
	if ModelCfg.DefaultModel == "" && len(added) > 0 {
		ModelCfg.DefaultModel = added[0]
		if ModelCfg.RoutingRules == nil {
			ModelCfg.RoutingRules = make(map[string]string)
		}
		ModelCfg.RoutingRules["chat"] = added[0]
	}
	if ModelCfg.AgentModel == "" && len(added) > 0 {
		ModelCfg.AgentModel = added[0]
	}

	SaveModelConfig()
	return added
}

// removeProviderModels removes all models for a provider from config
func RemoveProviderModels(provider string) int {
	ModelCfgMu.Lock()
	defer ModelCfgMu.Unlock()

	count := 0
	for name, m := range ModelCfg.Models {
		if m.Provider == provider {
			delete(ModelCfg.Models, name)
			count++
			// Clear roles
			if ModelCfg.DefaultModel == name {
				ModelCfg.DefaultModel = ""
				delete(ModelCfg.RoutingRules, "chat")
			}
			if ModelCfg.AgentModel == name {
				ModelCfg.AgentModel = ""
				delete(ModelCfg.RoutingRules, "agent")
			}
			if ModelCfg.PremiumModel == name {
				ModelCfg.PremiumModel = ""
			}
		}
	}

	// Clean fallback list
	var newFallback []string
	for _, f := range ModelCfg.FallbackModels {
		if _, ok := ModelCfg.Models[f]; ok {
			newFallback = append(newFallback, f)
		}
	}
	ModelCfg.FallbackModels = newFallback

	SaveModelConfig()
	return count
}
