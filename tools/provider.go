package tools

import (
	"scorp-agent/models"
	"scorp-agent/internal/helpers"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// ──────────────────────────────────────────────
// Provider Tools — Phase 2: Auto-Provisioning
// ──────────────────────────────────────────────

// models.ProviderRegistry maps provider name -> preset info (from providers.go)
// Re-export for tool use

// getProviderInfo returns display info for a provider
func getProviderInfo(name string) (models.ProviderPreset, bool) {
	preset, ok := models.ProviderRegistry[strings.ToLower(name)]
	return preset, ok
}

// listProvidersWithKeyStatus returns a formatted list of all providers with key detection
func listProvidersWithKeyStatus() string {
	var sb strings.Builder
	sb.WriteString("📋 **Available Providers**\n\n")

	// Sort keys for consistent output
	names := make([]string, 0, len(models.ProviderRegistry))
	for name := range models.ProviderRegistry {
		names = append(names, name)
	}
	// Simple sort
	for i := 0; i < len(names); i++ {
		for j := i + 1; j < len(names); j++ {
			if names[i] > names[j] {
				names[i], names[j] = names[j], names[i]
			}
		}
	}

	categories := map[string][]string{
		"☁️  Cloud APIs":       {"openai", "anthropic", "gemini", "groq", "deepseek", "mistral", "cohere", "together", "fireworks", "nvidia"},
		"🔀 Router/Aggregators": {"openrouter", "9router", "huggingface"},
		"🤖 Chinese Models":     {"zai", "kimi", "minimax"},
		"💻 Local/Offline":      {"ollama", "lmstudio"},
		"🛠  Dev Tools":         {"copilot"},
	}

	for cat, providers := range categories {
		sb.WriteString(fmt.Sprintf("**%s**\n", cat))
		for _, name := range providers {
			preset, ok := models.ProviderRegistry[name]
			if !ok {
				continue
			}

			// Check key status
			keyStatus := "❌ No key"
			if preset.NoAuth {
				keyStatus = "🟢 No auth needed (local)"
			} else {
				for _, env := range preset.KeyEnvs {
					if os.Getenv(env) != "" {
						keyStatus = fmt.Sprintf("🟢 %s set", env)
						break
					}
				}
				if keyStatus == "❌ No key" {
					generic := "SCORP_" + strings.ToUpper(name) + "_API_KEY"
					if os.Getenv(generic) != "" {
						keyStatus = fmt.Sprintf("🟢 %s set", generic)
					}
				}
			}

			sb.WriteString(fmt.Sprintf("  • **%s** (%s) — %s\n", preset.DisplayName, name, keyStatus))
			sb.WriteString(fmt.Sprintf("    Base: %s | API: %s\n", preset.BaseURL, preset.API))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("💡 **Usage:**\n")
	sb.WriteString("  `provider_add <name>` — Add model from provider (prompts for env var if needed)\n")
	sb.WriteString("  `provider_test <name>` — Test API connectivity\n")
	sb.WriteString("  `provider_remove <model_name>` — Remove a model from config\n")

	return sb.String()
}

// HandleProviderCommand is the main entry point for provider subcommands.
// Actions: list, add, test, remove, models
func HandleProviderCommand(args map[string]interface{}, chatID int64) (string, bool) {
	action := helpers.GetStringArg(args, "provider_action", "list")

	switch action {
	case "list":
		return listProvidersWithKeyStatus(), true

	case "add":
		providerName := strings.ToLower(helpers.GetStringArg(args, "provider_name", ""))
		if providerName == "" {
			return "Usage: `provider_add <provider_name>`\nExample: `provider_add deepseek`", false
		}
		return providerAddInteractive(providerName), true

	case "test":
		providerName := strings.ToLower(helpers.GetStringArg(args, "provider_name", ""))
		modelName := helpers.GetStringArg(args, "model_name", "")
		if providerName == "" {
			return "Usage: `provider_test <provider_name> [model_name]`", false
		}
		return providerTest(providerName, modelName), true

	case "remove":
		modelName := helpers.GetStringArg(args, "model_name", "")
		if modelName == "" {
			return "Usage: `provider_remove <model_name>`", false
		}
		return providerRemove(modelName), true

	case "models":
		return listConfiguredModels(), true

	default:
		return "Unknown action: " + action + "\nAvailable: list, add, test, remove, models", false
	}
}

// providerAddInteractive adds a model from a provider interactively.
// It prompts for model ID and detects/guides API key setup.
func providerAddInteractive(providerName string) string {
	preset, ok := getProviderInfo(providerName)
	if !ok {
		return fmt.Sprintf("❌ Unknown provider: %s\nUse `provider_list` to see available providers.", providerName)
	}

	// Determine model ID from preset defaults or ask user
	// For now, we'll use common model IDs per provider
	commonModels := map[string][]string{
		"openai":     {"gpt-4o", "gpt-4o-mini", "gpt-4-turbo", "gpt-3.5-turbo"},
		"deepseek":   {"deepseek-chat", "deepseek-coder"},
		"groq":       {"llama-3.1-70b-versatile", "llama-3.1-8b-instant", "mixtral-8x7b-32768"},
		"gemini":     {"gemini-1.5-pro", "gemini-1.5-flash", "gemini-2.0-flash-exp"},
		"anthropic":  {"claude-3-5-sonnet-20241022", "claude-3-5-haiku-20241022", "claude-3-opus-20240229"},
		"openrouter": {"openrouter/auto", "anthropic/claude-3.5-sonnet", "google/gemini-pro-1.5"},
		"zai":        {"glm-4.5", "glm-4-flash", "glm-4-air"},
		"kimi":       {"moonshot-v1-8k", "moonshot-v1-32k", "moonshot-v1-128k"},
		"minimax":    {"abab6.5s-chat-chat", "abab6.5-chat-32k"},
		"ollama":     {"llama3.1", "llama3.2", "qwen2.5", "mistral"},
		"lmstudio":   {"(any model loaded in LM Studio)"},
		"copilot":    {"gpt-4o", "gpt-4o-mini", "o1-preview", "o1-mini"},
		"nvidia":     {"nemotron-3-ultra", "nemotron-3-ultra-128k"},
		"huggingface": {"meta-llama/llama-3.1-70b-instruct", "microsoft/phi-3.5-mini"},
		"mistral":    {"mistral-large-latest", "mistral-small-latest", "pixtral-large-latest"},
		"cohere":     {"command-r-plus", "command-r", "command-r7b-128k"},
		"together":   {"meta-llama/Llama-3.1-70B-Instruct", "mistralai/Mixtral-8x7B-Instruct-v0.1"},
		"fireworks":  {"accounts/fireworks/models/llama-v3p1-70b-instruct"},
	}

	suggested := commonModels[providerName]
	modelID := ""
	if len(suggested) > 0 {
		modelID = suggested[0]
	}

	// Build the model config entry
	models.ModelCfgMu.Lock()
	defer models.ModelCfgMu.Unlock()

	if models.ModelCfg == nil {
		models.LoadModelConfig()
	}
	if models.ModelCfg.Models == nil {
		models.ModelCfg.Models = make(map[string]models.ModelConfig)
	}

	// Generate unique model name
	baseName := providerName
	if modelID != "" {
		// Sanitize model ID for config key
		cleanID := strings.ReplaceAll(modelID, "/", "-")
		cleanID = strings.ReplaceAll(cleanID, ":", "-")
		baseName = providerName + "-" + cleanID
	}

	// Ensure unique name
	name := baseName
	counter := 1
	for {
		if _, exists := models.ModelCfg.Models[name]; !exists {
			break
		}
		counter++
		name = fmt.Sprintf("%s-%d", baseName, counter)
	}

	// Determine key_env
	keyEnv := ""
	if !preset.NoAuth && len(preset.KeyEnvs) > 0 {
		keyEnv = preset.KeyEnvs[0]
		// Check if already set
		if os.Getenv(keyEnv) == "" {
			genericKey := "SCORP_" + strings.ToUpper(providerName) + "_API_KEY"
			if os.Getenv(genericKey) != "" {
				keyEnv = genericKey
			}
		}
	}

	models.ModelCfg.Models[name] = models.ModelConfig{
		Provider:  providerName,
		Model:     modelID,
		KeyEnv:    keyEnv,
		BaseURL:   preset.BaseURL,
		MaxTokens: 8192,
		API:       preset.API,
	}

	models.SaveModelConfig()

	// Build response
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("✅ **Model added: %s**\n\n", name))
	sb.WriteString(fmt.Sprintf("Provider: %s (%s)\n", preset.DisplayName, providerName))
	sb.WriteString(fmt.Sprintf("Model ID: %s\n", modelID))
	sb.WriteString(fmt.Sprintf("Base URL: %s\n", preset.BaseURL))
	sb.WriteString(fmt.Sprintf("API Format: %s\n", preset.API))
	sb.WriteString(fmt.Sprintf("Max Tokens: 8192\n"))

	if preset.NoAuth {
		sb.WriteString("\n🟢 **No API key required** (local provider)\n")
		sb.WriteString("Ready to use!\n")
	} else if keyEnv != "" {
		if os.Getenv(keyEnv) != "" {
			sb.WriteString(fmt.Sprintf("\n🟢 **API key detected**: %s is set\n", keyEnv))
			sb.WriteString("Ready to use!\n")
		} else {
			sb.WriteString(fmt.Sprintf("\n⚠️ **API key needed**: Set `%s` environment variable\n", keyEnv))
			sb.WriteString("Example: `export " + keyEnv + "=\"sk-xxxx\"`\n")
			sb.WriteString("Then restart scorp-agent or run `models reload`\n")
		}
	} else {
		genericKey := "SCORP_" + strings.ToUpper(providerName) + "_API_KEY"
		sb.WriteString(fmt.Sprintf("\n⚠️ **API key needed**: Set `%s` or provider-specific env var\n", genericKey))
		sb.WriteString(fmt.Sprintf("Available env vars for %s: %s\n", providerName, strings.Join(preset.KeyEnvs, ", ")))
	}

	sb.WriteString(fmt.Sprintf("\n📝 **Config entry** (in models.json):\n```json\n"))
	sb.WriteString(fmt.Sprintf("  \"%s\": {\n", name))
	sb.WriteString(fmt.Sprintf("    \"provider\": \"%s\",\n", providerName))
	sb.WriteString(fmt.Sprintf("    \"model\": \"%s\",\n", modelID))
	if keyEnv != "" {
		sb.WriteString(fmt.Sprintf("    \"key_env\": \"%s\",\n", keyEnv))
	}
	sb.WriteString(fmt.Sprintf("    \"base_url\": \"%s\",\n", preset.BaseURL))
	sb.WriteString(fmt.Sprintf("    \"max_tokens\": 8192,\n"))
	sb.WriteString(fmt.Sprintf("    \"api\": \"%s\"\n", preset.API))
	sb.WriteString(fmt.Sprintf("  }\n```\n"))

	if len(suggested) > 1 {
		sb.WriteString(fmt.Sprintf("\n💡 **Other common models**: %s\n", strings.Join(suggested[1:], ", ")))
		sb.WriteString("To use a different model, edit models.json and change the `model` field.\n")
	}

	sb.WriteString(fmt.Sprintf("\n🧪 Test it: `provider_test %s %s`", providerName, name))

	return sb.String()
}

// providerTest tests API connectivity for a provider/model.
func providerTest(providerName, modelName string) string {
	preset, ok := getProviderInfo(providerName)
	if !ok {
		return fmt.Sprintf("❌ Unknown provider: %s", providerName)
	}

	models.ModelCfgMu.RLock()
	var model *models.ModelConfig
	if modelName != "" {
		model = models.GetModelByName(modelName)
	} else {
		// Find first model for this provider
		for _, m := range models.ModelCfg.Models {
			if strings.EqualFold(m.Provider, providerName) {
				model = &m
				break
			}
		}
	}
	models.ModelCfgMu.RUnlock()

	if model == nil {
		return fmt.Sprintf("❌ No model configured for provider %s. Run `provider_add %s` first.", providerName, providerName)
	}

	// Check API key
	apiKey := models.ResolveAPIKey(model)
	if apiKey == "" && !preset.NoAuth {
		return fmt.Sprintf("❌ No API key for %s (%s). Set %s or run `provider_add` to configure.", modelName, providerName, models.KeySourceLabel(model))
	}

	// Make a simple test request
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	testMsg := []models.ChatMessage{
		{Role: "user", Content: "Say 'OK' if you receive this."},
	}

	log.Printf("[provider_test] Testing %s (%s)...", modelName, model.Model)

	result, err := models.CallModel(ctx, model, testMsg)
	if err != nil {
		return fmt.Sprintf("❌ **Test failed for %s**\n\nError: %v\n\nCheck:\n  • API key is valid and has credits\n  • Base URL is correct: %s\n  • Network connectivity", modelName, err, model.BaseURL)
	}

	return fmt.Sprintf("✅ **Test passed for %s**\n\nModel: %s\nProvider: %s\nResponse: %s\nLatency: <15s", modelName, model.Model, preset.DisplayName, helpers.TruncateStr(result, 200))
}

// providerRemove removes a model from config.
func providerRemove(modelName string) string {
	models.ModelCfgMu.Lock()
	defer models.ModelCfgMu.Unlock()

	if models.ModelCfg == nil {
		models.LoadModelConfig()
	}
	if models.ModelCfg.Models == nil {
		return "❌ No models configured"
	}

	if _, exists := models.ModelCfg.Models[modelName]; !exists {
		return fmt.Sprintf("❌ Model '%s' not found in config", modelName)
	}

	// Don't allow removing if it's the default/agent/premium model
	protected := map[string]bool{
		models.ModelCfg.DefaultModel: true,
		models.ModelCfg.AgentModel:   true,
		models.ModelCfg.PremiumModel: true,
	}
	if protected[modelName] {
		return fmt.Sprintf("❌ Cannot remove '%s' — it's set as default/agent/premium model.\nChange it first with `models set <type> <name>`", modelName)
	}

	delete(models.ModelCfg.Models, modelName)
	models.SaveModelConfig()

	return fmt.Sprintf("✅ Removed model '%s' from config", modelName)
}

// listConfiguredModels lists all models currently in config.
func listConfiguredModels() string {
	models.ModelCfgMu.RLock()
	defer models.ModelCfgMu.RUnlock()

	if models.ModelCfg == nil || len(models.ModelCfg.Models) == 0 {
		return "📭 No models configured.\nRun `provider_add <provider>` to add one."
	}

	var sb strings.Builder
	sb.WriteString("📋 **Configured Models**\n\n")

	for name, m := range models.ModelCfg.Models {
		status := "🟢"
		keyStatus := "OK"
		if !strings.Contains(m.API, "local") {
			apiKey := models.ResolveAPIKey(&m)
			if apiKey == "" {
				status = "🔴"
				keyStatus = "NO KEY"
			}
		}

		tags := []string{}
		if name == models.ModelCfg.DefaultModel {
			tags = append(tags, "🏠 default")
		}
		if name == models.ModelCfg.AgentModel {
			tags = append(tags, "🤖 agent")
		}
		if name == models.ModelCfg.PremiumModel {
			tags = append(tags, "⭐ premium")
		}

		sb.WriteString(fmt.Sprintf("%s **%s** %s\n", status, name, strings.Join(tags, " ")))
		sb.WriteString(fmt.Sprintf("  Provider: %s | Model: %s | Key: %s\n", m.Provider, m.Model, keyStatus))
		sb.WriteString(fmt.Sprintf("  Base: %s | API: %s | MaxTokens: %d\n\n", m.BaseURL, m.API, m.MaxTokens))
	}

	sb.WriteString(fmt.Sprintf("**Routing**: default=%s | agent=%s | premium=%s\n",
		models.ModelCfg.DefaultModel, models.ModelCfg.AgentModel, models.ModelCfg.PremiumModel))

	if len(models.ModelCfg.FallbackModels) > 0 {
		sb.WriteString(fmt.Sprintf("**Fallback Chain**: %s\n", strings.Join(models.ModelCfg.FallbackModels, " → ")))
	}
	if len(models.ModelCfg.FallbackOnError) > 0 {
		sb.WriteString(fmt.Sprintf("**Fallback Triggers**: %s\n", strings.Join(models.ModelCfg.FallbackOnError, ", ")))
	}

	return sb.String()
}