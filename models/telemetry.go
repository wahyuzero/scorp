package models

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"scorp-agent/config"
	"scorp-agent/internal/helpers"
)

// ──────────────────────────────────────────────
// Usage tracking & persistence
// ──────────────────────────────────────────────

func InitModelUsage() {
	ModelUsageMu.Lock()
	defer ModelUsageMu.Unlock()

	ModelUsageMap = make(map[string]*ModelUsage)

	// Load from file
	path := config.ModelUsageFilePath()
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	json.Unmarshal(data, &ModelUsageMap)
}

func TrackModelUsage(model string, inputTokens, outputTokens int) {
	TrackModelUsageWithCache(model, inputTokens, outputTokens, 0)
}

func TrackModelUsageWithCache(model string, inputTokens, outputTokens, cachedTokens int) {
	ModelUsageMu.Lock()
	defer ModelUsageMu.Unlock()

	if ModelUsageMap == nil {
		ModelUsageMap = make(map[string]*ModelUsage)
	}

	u, ok := ModelUsageMap[model]
	if !ok {
		u = &ModelUsage{Model: model}
		ModelUsageMap[model] = u
	}

	u.InputTokens += inputTokens
	u.OutputTokens += outputTokens
	u.CachedTokens += cachedTokens
	u.Calls++
	u.LastUsed = time.Now()

	// Save
	path := config.ModelUsageFilePath()
	data, _ := json.MarshalIndent(ModelUsageMap, "", "  ")
	os.WriteFile(path, data, 0644)
}

// ──────────────────────────────────────────────
// Display & Health Check Helpers
// ──────────────────────────────────────────────

func FormatModelList() string {
	ModelCfgMu.RLock()
	defer ModelCfgMu.RUnlock()

	if ModelCfg == nil {
		return "❌ No model config loaded."
	}

	var sb strings.Builder
	sb.WriteString("🤖 <b>AI Models</b>\n\n")

	for name, m := range ModelCfg.Models {
		isDefault := ""
		if name == ModelCfg.DefaultModel {
			isDefault += " 💬chat"
		}
		if name == ModelCfg.AgentModel {
			isDefault += " 🤖agent"
		}
		if name == ModelCfg.PremiumModel {
			isDefault += " 🧠premium"
		}
		if name == ModelCfg.VisionModel {
			isDefault += " 👁vision"
		}

		apiStatus := "⏳"
		if !hasAPIKey(&m) {
			apiStatus = "⚠️ " + KeySourceLabel(&m)
		}

		sb.WriteString(fmt.Sprintf("%s <code>%s</code> — %s (%s)%s\n",
			apiStatus, name, m.Model, m.Provider, isDefault))
	}

	sb.WriteString("\n<b>Routing:</b>\n")
	for task, model := range ModelCfg.RoutingRules {
		sb.WriteString(fmt.Sprintf("  %s → <code>%s</code>\n", task, model))
	}

	sb.WriteString("\nCommands:\n")
	sb.WriteString("<code>/model check</code> — test all API keys\n")
	sb.WriteString("<code>/model use [name]</code> — switch default\n")
	sb.WriteString("<code>/model agent [name]</code> — switch agent\n")
	sb.WriteString("<code>/usage</code> — token usage stats")

	return sb.String()
}

// CheckModelHealth sends a tiny test request to verify the API key works
func CheckModelHealth(name string, m *ModelConfig) (bool, string) {
	apiKey := ResolveAPIKey(m)
	if apiKey == "" && !ProviderRegistry[m.Provider].NoAuth {
		return false, "no API key (" + KeySourceLabel(m) + ")"
	}

	// Use callModel for format-agnostic health check.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := CallModel(ctx, m, []ChatMessage{{Role: "user", Content: "hi"}})
	if err != nil {
		return false, helpers.TruncateStr(err.Error(), 60)
	}
	return true, "ok"
}

// FormatModelListWithHealth checks all models and returns with health status
func FormatModelListWithHealth() string {
	ModelCfgMu.RLock()
	cfgCopy := *ModelCfg
	ModelCfgMu.RUnlock()

	if len(cfgCopy.Models) == 0 {
		return "❌ No models configured."
	}

	// Check all models in parallel
	type healthResult struct {
		Name    string
		Ok      bool
		Detail  string
		Latency time.Duration
	}

	results := make(chan healthResult, len(cfgCopy.Models))

	for name, m := range cfgCopy.Models {
		go func(n string, mc ModelConfig) {
			start := time.Now()
			ok, detail := CheckModelHealth(n, &mc)
			results <- healthResult{Name: n, Ok: ok, Detail: detail, Latency: time.Since(start)}
		}(name, m)
	}

	// Collect results
	healthMap := make(map[string]healthResult)
	for i := 0; i < len(cfgCopy.Models); i++ {
		r := <-results
		healthMap[r.Name] = r
	}

	var sb strings.Builder
	sb.WriteString("🤖 <b>AI Models — Health Check</b>\n\n")

	for name, m := range cfgCopy.Models {
		h := healthMap[name]

		isDefault := ""
		if name == cfgCopy.DefaultModel {
			isDefault += " 💬"
		}
		if name == cfgCopy.AgentModel {
			isDefault += " 🤖"
		}
		if name == cfgCopy.PremiumModel {
			isDefault += " 🧠"
		}
		if name == cfgCopy.VisionModel {
			isDefault += " 👁"
		}

		status := "✅"
		latency := fmt.Sprintf("%.1fs", h.Latency.Seconds())
		if !h.Ok {
			status = "❌"
			latency = h.Detail
		}

		sb.WriteString(fmt.Sprintf("%s <code>%s</code> %s\n   %s (%s) [%s]\n",
			status, name, isDefault, m.Model, m.Provider, latency))
	}

	sb.WriteString("\n<b>Routing:</b>\n")
	for task, model := range cfgCopy.RoutingRules {
		sb.WriteString(fmt.Sprintf("  %s → <code>%s</code>\n", task, model))
	}

	return sb.String()
}

func FormatUsageStats() string {
	ModelUsageMu.Lock()
	defer ModelUsageMu.Unlock()

	if len(ModelUsageMap) == 0 {
		return "📊 No usage data yet."
	}

	// Pricing per 1M tokens (approximate)
	pricing := map[string][]float64{
		"deepseek/deepseek-v4-flash":                {0.22, 0.66},
		"deepseek/deepseek-v4-pro":                  {0.66, 1.98},
		"poolside/laguna-s-2.1-free":                {0, 0},
		"gpt-5.6-luna":                              {0.20, 1.20},
		"meta/muse-spark-1.2-contributor":           {0.10, 0.20},
		"z-ai/glm-5.3-flash":                        {0.15, 0.50},
		"xiaomi/mimo-v2.5":                          {0.14, 0.28},
		"Qwen/Qwen3.8-Flash":                        {0.16, 0.47},
		"opencode/big-pickle":                       {0, 0},
		"opencode/mimo-v2.5-free":                   {0, 0},
		"opencode/ling-3.0-flash-fin-free":          {0, 0},
		"llama-3.3-70b-versatile":                   {0, 0},
		"qwen/qwen3-32b":                            {0, 0},
		"meta-llama/llama-4-scout-17b-16e-instruct": {0, 0},
		"gemini-2.5-flash":                          {0, 0},
		"glm-4.7":                                   {0, 0},
		"MiniMax-M2":                                {0.25, 1.00},
		"deepseek-chat":                             {0.05, 0.50},
	}

	var sb strings.Builder
	sb.WriteString("📊 <b>Token Usage</b>\n\n")

	totalCost := 0.0
	for _, u := range ModelUsageMap {
		cost := 0.0
		if p, ok := pricing[u.Model]; ok && len(p) == 2 {
			noCache := u.InputTokens - u.CachedTokens
			if noCache < 0 {
				noCache = 0
			}
			cacheRate := p[0] * 0.10 // 90% discount on cached tokens
			cost = (float64(noCache)/1e6)*p[0] + (float64(u.CachedTokens)/1e6)*cacheRate + (float64(u.OutputTokens)/1e6)*p[1]
		}
		totalCost += cost

		costStr := "FREE"
		if cost > 0 {
			costStr = fmt.Sprintf("$%.4f", cost)
		}

		cachedInfo := ""
		if u.CachedTokens > 0 && u.InputTokens > 0 {
			cachedInfo = fmt.Sprintf(" | Cached: %d (%.0f%%)", u.CachedTokens, float64(u.CachedTokens)*100.0/float64(u.InputTokens))
		}

		sb.WriteString(fmt.Sprintf("<b>%s</b>\n", u.Model))
		sb.WriteString(fmt.Sprintf("  Calls: %d | In: %d%s | Out: %d\n", u.Calls, u.InputTokens, cachedInfo, u.OutputTokens))
		sb.WriteString(fmt.Sprintf("  Cost: %s | Last: %s\n\n", costStr, u.LastUsed.Format("01/02 15:04")))
	}

	if totalCost > 0 {
		sb.WriteString(fmt.Sprintf("💰 <b>Total estimated cost: $%.4f</b>", totalCost))
	} else {
		sb.WriteString("💰 <b>Total cost: $0 (all free models!)</b>")
	}

	return sb.String()
}

// SwitchModel changes the active model for a given role
func SwitchModel(role, modelName string) string {
	ModelCfgMu.Lock()
	defer ModelCfgMu.Unlock()

	if ModelCfg == nil {
		return "❌ No model config loaded."
	}

	if _, ok := ModelCfg.Models[modelName]; !ok {
		names := make([]string, 0)
		for n := range ModelCfg.Models {
			names = append(names, n)
		}
		return fmt.Sprintf("❌ Model '%s' not found. Available: %s", modelName, strings.Join(names, ", "))
	}

	switch role {
	case "default", "chat", "use":
		ModelCfg.DefaultModel = modelName
		ModelCfg.RoutingRules["chat"] = modelName
	case "agent":
		ModelCfg.AgentModel = modelName
		ModelCfg.RoutingRules["agent"] = modelName
	case "delegation", "delegate":
		ModelCfg.DelegationModel = modelName
	case "vision":
		ModelCfg.VisionModel = modelName
		if ModelCfg.RoutingRules == nil {
			ModelCfg.RoutingRules = make(map[string]string)
		}
		ModelCfg.RoutingRules["vision"] = modelName
	case "premium":
		ModelCfg.PremiumModel = modelName
	default:
		ModelCfg.DefaultModel = modelName
	}

	SaveModelConfig()
	return fmt.Sprintf("✅ %s model set to <code>%s</code>", role, modelName)
}
