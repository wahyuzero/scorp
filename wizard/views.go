package wizard

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"scorp-agent/models"
)

// ──────────────────────────────────────────────
// Text Formatters for Telegram Model Manager
// ──────────────────────────────────────────────

func ModelMenuText() string {
	models.ModelCfgMu.RLock()
	defer models.ModelCfgMu.RUnlock()

	var sb strings.Builder
	sb.WriteString("🤖 <b>Model Manager</b>\n\n")

	if len(models.ModelCfg.Models) == 0 {
		sb.WriteString("<i>No models yet. Tap \"Add Provider\" to auto-populate from catalog.</i>")
		return sb.String()
	}

	sb.WriteString("<i>Tap a role to change model, or tap a model for details.</i>\n")
	sb.WriteString(fmt.Sprintf("\n<b>%d Models</b>\n", len(models.ModelCfg.Models)))

	names := make([]string, 0)
	for n := range models.ModelCfg.Models {
		names = append(names, n)
	}
	sort.Strings(names)

	for _, name := range names {
		m := models.ModelCfg.Models[name]
		roles := []string{}
		if name == models.ModelCfg.DefaultModel {
			roles = append(roles, "Primary")
		}
		if name == models.ModelCfg.AgentModel && name != models.ModelCfg.DefaultModel {
			roles = append(roles, "Agent")
		}
		if name == models.ModelCfg.DelegationModel && name != models.ModelCfg.DefaultModel && name != models.ModelCfg.AgentModel {
			roles = append(roles, "Delegation")
		}
		if name == models.ModelCfg.VisionModel && name != models.ModelCfg.DefaultModel && name != models.ModelCfg.AgentModel && name != models.ModelCfg.DelegationModel {
			roles = append(roles, "Vision")
		}
		roleStr := ""
		if len(roles) > 0 {
			roleStr = " [" + strings.Join(roles, ", ") + "]"
		}
		sb.WriteString(fmt.Sprintf("• <code>%s%s</code> → %s (%s)\n", name, roleStr, m.Model, m.Provider))
	}

	return sb.String()
}

func OrNA(s string) string {
	if s == "" {
		return "(not set)"
	}
	return s
}

func ModelInfoText(name string) string {
	models.ModelCfgMu.RLock()
	defer models.ModelCfgMu.RUnlock()

	m, ok := models.ModelCfg.Models[name]
	if !ok {
		return fmt.Sprintf("❌ Model '%s' not found.", name)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("🔧 <b>%s</b>\n\n", name))
	sb.WriteString(fmt.Sprintf("Model ID: <code>%s</code>\n", m.Model))
	sb.WriteString(fmt.Sprintf("Provider: <code>%s</code>\n", m.Provider))
	sb.WriteString(fmt.Sprintf("Base URL: <code>%s</code>\n", m.BaseURL))
	sb.WriteString(fmt.Sprintf("API: <code>%s</code>\n", m.API))
	sb.WriteString(fmt.Sprintf("Max Tokens: <code>%d</code>\n", m.MaxTokens))

	if m.KeyEnv != "" {
		keyStatus := "❌ not set"
		if os.Getenv(m.KeyEnv) != "" {
			keyStatus = "✅ set"
		}
		sb.WriteString(fmt.Sprintf("Key: <code>%s</code> (%s)\n", m.KeyEnv, keyStatus))
	}

	sb.WriteString("\n")
	if name == models.ModelCfg.DefaultModel {
		sb.WriteString("💬 <b>Primary model</b>\n")
	}
	if name == models.ModelCfg.AgentModel {
		sb.WriteString("🤖 <b>Agent model</b>\n")
	}
	if name == models.ModelCfg.DelegationModel {
		sb.WriteString("🎯 <b>Delegation model</b>\n")
	}
	if name == models.ModelCfg.VisionModel {
		sb.WriteString("👁 <b>Vision model</b>\n")
	}

	for i, f := range models.ModelCfg.FallbackModels {
		if f == name {
			sb.WriteString(fmt.Sprintf("🔄 <b>Fallback #%d</b>\n", i+1))
			break
		}
	}

	return sb.String()
}

func FallbackText() string {
	models.ModelCfgMu.RLock()
	defer models.ModelCfgMu.RUnlock()

	var sb strings.Builder
	sb.WriteString("🔄 <b>Fallback Chain</b>\n\n")
	if len(models.ModelCfg.FallbackModels) == 0 {
		sb.WriteString("<i>No fallback models configured.</i>\n\n")
		sb.WriteString("Fallback models are tried automatically when the primary model fails.\n")
		sb.WriteString("Triggers: <code>" + strings.Join(models.ModelCfg.FallbackOnError, ", ") + "</code>")
	} else {
		for i, name := range models.ModelCfg.FallbackModels {
			sb.WriteString(fmt.Sprintf("%d️⃣ <code>%s</code>\n", i+1, name))
		}
	}
	return sb.String()
}

func FormatAPIKeysList() string {
	models.ModelCfgMu.RLock()
	defer models.ModelCfgMu.RUnlock()

	var sb strings.Builder
	sb.WriteString("🔑 <b>API Keys</b>\n\n")

	keyEnvs := make(map[string]string)
	for _, m := range models.ModelCfg.Models {
		if m.KeyEnv != "" {
			if os.Getenv(m.KeyEnv) != "" {
				keyEnvs[m.KeyEnv] = "✅ set"
			} else {
				keyEnvs[m.KeyEnv] = "❌ not set"
			}
		}
	}

	if len(keyEnvs) == 0 {
		sb.WriteString("<i>No API keys configured.</i>\n")
	} else {
		names := make([]string, 0)
		for k := range keyEnvs {
			names = append(names, k)
		}
		sort.Strings(names)
		for _, k := range names {
			sb.WriteString(fmt.Sprintf("<code>%s</code> — %s\n", k, keyEnvs[k]))
		}
	}

	sb.WriteString("\n<i>Tap a model → \"Set API Key\" to update.</i>")
	return sb.String()
}

func formatProvidersList() string {
	var sb strings.Builder
	sb.WriteString("🔌 <b>Providers</b>\n\n")
	sb.WriteString("<b>Built-in:</b>\n")
	names := AllProviderNames()
	builtInCount := 0
	for _, name := range names {
		if preset, ok := models.ProviderRegistry[name]; ok {
			keyStatus := "❌"
			if models.ProviderHasAPIKey(name) {
				keyStatus = "✅"
			}
			catStr := ""
			if models.HasCatalog(name) {
				catStr = fmt.Sprintf(" [%d models]", len(models.CatalogModels(name)))
			}
			sb.WriteString(fmt.Sprintf("%s <code>%s</code> — %s (%s)%s\n", keyStatus, name, preset.DisplayName, preset.API, catStr))
			builtInCount++
		}
	}

	models.ModelCfgMu.RLock()
	customCount := 0
	if models.ModelCfg != nil && len(models.ModelCfg.CustomProviders) > 0 {
		sb.WriteString("\n<b>Custom:</b>\n")
		for name, cp := range models.ModelCfg.CustomProviders {
			disp := cp.DisplayName
			if disp == "" {
				disp = name
			}
			keyStatus := "❌"
			if models.ProviderHasAPIKey(name) {
				keyStatus = "✅"
			}
			sb.WriteString(fmt.Sprintf("%s <code>%s</code> — %s (%s)\n", keyStatus, name, disp, cp.API))
			customCount++
		}
	}
	models.ModelCfgMu.RUnlock()

	sb.WriteString(fmt.Sprintf("\n📊 %d built-in, %d custom", builtInCount, customCount))
	return sb.String()
}

// ──────────────────────────────────────────────
// Fallback Chain Operations
// ──────────────────────────────────────────────

func addToFallback(name string) string {
	models.ModelCfgMu.Lock()
	defer models.ModelCfgMu.Unlock()
	if _, ok := models.ModelCfg.Models[name]; !ok {
		return fmt.Sprintf("❌ Model '%s' not found.", name)
	}
	for _, f := range models.ModelCfg.FallbackModels {
		if f == name {
			return fmt.Sprintf("ℹ️ <code>%s</code> already in fallback chain.", name)
		}
	}
	models.ModelCfg.FallbackModels = append(models.ModelCfg.FallbackModels, name)
	models.SaveModelConfig()
	return fmt.Sprintf("✅ Added <code>%s</code> to fallback chain.", name)
}

func removeFromFallback(name string) string {
	models.ModelCfgMu.Lock()
	defer models.ModelCfgMu.Unlock()
	var newList []string
	for _, f := range models.ModelCfg.FallbackModels {
		if f != name {
			newList = append(newList, f)
		}
	}
	models.ModelCfg.FallbackModels = newList
	models.SaveModelConfig()
	return fmt.Sprintf("✅ Removed <code>%s</code> from fallback.", name)
}

func moveFallback(name string, direction string) string {
	models.ModelCfgMu.Lock()
	defer models.ModelCfgMu.Unlock()
	idx := -1
	for i, f := range models.ModelCfg.FallbackModels {
		if f == name {
			idx = i
			break
		}
	}
	if idx < 0 {
		return "❌ Not in fallback chain."
	}
	if direction == "up" && idx > 0 {
		models.ModelCfg.FallbackModels[idx], models.ModelCfg.FallbackModels[idx-1] = models.ModelCfg.FallbackModels[idx-1], models.ModelCfg.FallbackModels[idx]
	} else if direction == "down" && idx < len(models.ModelCfg.FallbackModels)-1 {
		models.ModelCfg.FallbackModels[idx], models.ModelCfg.FallbackModels[idx+1] = models.ModelCfg.FallbackModels[idx+1], models.ModelCfg.FallbackModels[idx]
	}
	models.SaveModelConfig()
	return ""
}
