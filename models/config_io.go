package models

import (
	"log"
	"sync"
	"time"

	"scorp-agent/config"
)

// ──────────────────────────────────────────────
// Model Configuration Types
// ──────────────────────────────────────────────

type ModelConfig struct {
	Provider  string `json:"provider"`          // provider name (openai, deepseek, 9router, etc)
	Model     string `json:"model"`             // model ID
	APIKey    string `json:"api_key,omitempty"` // DEPRECATED — use key_env instead
	KeyEnv    string `json:"key_env,omitempty"` // env var name for API key
	BaseURL   string `json:"base_url"`          // OpenAI-compatible endpoint
	MaxTokens int    `json:"max_tokens"`        // max output tokens
	API       string `json:"api,omitempty"`     // "openai" | "anthropic" | "gemini"
}

type ModelRouterConfig struct {
	DefaultModel    string                 `json:"default_model"`          // 💬 primary chat model
	AgentModel      string                 `json:"agent_model"`            // 🤖 agent mode model
	DelegationModel string                 `json:"delegation_model"`       // 🎯 delegation/subagent model
	VisionModel     string                 `json:"vision_model,omitempty"` // 👁 vision/image analysis model
	PremiumModel    string                 `json:"premium_model"`          // 💎 complex tasks (optional)
	Models          map[string]ModelConfig `json:"models"`
	RoutingRules    map[string]string      `json:"routing_rules"` // taskType → modelName

	// Fallback chain configuration
	FallbackModels  []string `json:"fallback_models,omitempty"`   // ordered list of model names to try
	FallbackOnError []string `json:"fallback_on_error,omitempty"` // error types that trigger fallback: "rate_limit", "timeout", "server_error", "auth_error"

	// Custom providers (user-defined via /model add)
	CustomProviders map[string]CustomProvider `json:"custom_providers,omitempty"`
}

// ModelUsage tracks token consumption and cost per model
type ModelUsage struct {
	Model        string    `json:"model"`
	InputTokens  int       `json:"input_tokens"`
	OutputTokens int       `json:"output_tokens"`
	CachedTokens int       `json:"cached_tokens"`
	Calls        int       `json:"calls"`
	LastUsed     time.Time `json:"last_used"`
}

var (
	ModelCfg      *ModelRouterConfig
	ModelCfgMu    sync.RWMutex
	ModelUsageMap map[string]*ModelUsage
	ModelUsageMu  sync.Mutex
)

// ──────────────────────────────────────────────
// Config loading & persistence (using ConfigManager)
// ──────────────────────────────────────────────

func LoadModelConfig() {
	ModelCfgMu.Lock()
	defer ModelCfgMu.Unlock()

	if err := config.ConfigMgr().Load("models.json", &ModelCfg); err != nil {
		log.Printf("[models] Load error: %v, using defaults", err)
		ModelCfg = defaultModelConfig()
		SaveModelConfig()
		return
	}
	if ModelCfg == nil || len(ModelCfg.Models) == 0 || ModelCfg.AgentModel == "" {
		ModelCfg = defaultModelConfig()
		SaveModelConfig()
	}

	// Auto-migrate plaintext api_key → key_env, fill provider defaults
	migrateModelConfigs(ModelCfg)

	// Merge custom providers into runtime registry
	for name, cp := range ModelCfg.CustomProviders {
		ProviderRegistry[name] = ProviderPreset{
			BaseURL:     cp.BaseURL,
			API:         cp.API,
			KeyEnvs:     cp.KeyEnvs,
			DisplayName: cp.DisplayName,
			NoAuth:      cp.NoAuth,
		}
		log.Printf("[models] Custom provider: %s (%s)", name, cp.API)
	}

	log.Printf("[models] Loaded %d models (default=%s, agent=%s)",
		len(ModelCfg.Models), ModelCfg.DefaultModel, ModelCfg.AgentModel)
}

func SaveModelConfig() {
	if ModelCfg == nil {
		return
	}
	if err := config.ConfigMgr().Save("models.json", ModelCfg); err != nil {
		log.Printf("[models] Save error: %v", err)
	}
}

func defaultModelConfig() *ModelRouterConfig {
	return &ModelRouterConfig{
		DefaultModel:    "deepseek/deepseek-v4-flash",
		AgentModel:      "deepseek/deepseek-v4-flash",
		DelegationModel: "deepseek/deepseek-v4-flash",
		PremiumModel:    "deepseek/deepseek-v4-pro",
		RoutingRules: map[string]string{
			"agent":   "deepseek/deepseek-v4-flash",
			"chat":    "deepseek/deepseek-v4-flash",
			"complex": "deepseek/deepseek-v4-pro",
		},
		FallbackModels: []string{
			"z-ai/glm-5.3-flash",
			"meta/muse-spark-1.2-contributor",
			"poolside/laguna-s-2.1-free",
		},
		FallbackOnError: []string{"rate_limit", "timeout", "server_error", "auth_error", "network_error"},
		Models: map[string]ModelConfig{
			"deepseek/deepseek-v4-flash": {
				Provider:  "command-code",
				Model:     "deepseek/deepseek-v4-flash",
				KeyEnv:    "COMMAND_CODE_API_KEY",
				BaseURL:   "https://api.commandcode.ai",
				MaxTokens: 16384,
				API:       "command-code",
			},
			"deepseek/deepseek-v4-pro": {
				Provider:  "command-code",
				Model:     "deepseek/deepseek-v4-pro",
				KeyEnv:    "COMMAND_CODE_API_KEY",
				BaseURL:   "https://api.commandcode.ai",
				MaxTokens: 32768,
				API:       "command-code",
			},
			"poolside/laguna-s-2.1-free": {
				Provider:  "command-code",
				Model:     "poolside/laguna-s-2.1-free",
				KeyEnv:    "COMMAND_CODE_API_KEY",
				BaseURL:   "https://api.commandcode.ai",
				MaxTokens: 16384,
				API:       "command-code",
			},
			"gpt-5.6-luna": {
				Provider:  "command-code",
				Model:     "gpt-5.6-luna",
				KeyEnv:    "COMMAND_CODE_API_KEY",
				BaseURL:   "https://api.commandcode.ai",
				MaxTokens: 32768,
				API:       "command-code",
			},
			"meta/muse-spark-1.2-contributor": {
				Provider:  "command-code",
				Model:     "meta/muse-spark-1.2-contributor",
				KeyEnv:    "COMMAND_CODE_API_KEY",
				BaseURL:   "https://api.commandcode.ai",
				MaxTokens: 16384,
				API:       "command-code",
			},
			"z-ai/glm-5.3-flash": {
				Provider:  "command-code",
				Model:     "z-ai/glm-5.3-flash",
				KeyEnv:    "COMMAND_CODE_API_KEY",
				BaseURL:   "https://api.commandcode.ai",
				MaxTokens: 16384,
				API:       "command-code",
			},
			"xiaomi/mimo-v2.5": {
				Provider:  "command-code",
				Model:     "xiaomi/mimo-v2.5",
				KeyEnv:    "COMMAND_CODE_API_KEY",
				BaseURL:   "https://api.commandcode.ai",
				MaxTokens: 16384,
				API:       "command-code",
			},
			"Qwen/Qwen3.8-Flash": {
				Provider:  "command-code",
				Model:     "Qwen/Qwen3.8-Flash",
				KeyEnv:    "COMMAND_CODE_API_KEY",
				BaseURL:   "https://api.commandcode.ai",
				MaxTokens: 16384,
				API:       "command-code",
			},
		},
	}
}
