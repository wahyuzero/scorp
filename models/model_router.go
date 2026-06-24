package models

import (
	"scorp-agent/internal/helpers"
	"scorp-agent/registry"
	"scorp-agent/config"
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// ──────────────────────────────────────────────
// Multi-Model Router
// ──────────────────────────────────────────────

type ModelConfig struct {
	Provider  string `json:"provider"`             // provider name (openai, deepseek, 9router, etc)
	Model     string `json:"model"`                // model ID
	APIKey    string `json:"api_key,omitempty"`    // DEPRECATED — use key_env instead
	KeyEnv    string `json:"key_env,omitempty"`    // env var name for API key
	BaseURL   string `json:"base_url"`             // OpenAI-compatible endpoint
	MaxTokens int    `json:"max_tokens"`           // max output tokens
	API       string `json:"api,omitempty"`        // "openai" | "anthropic" | "gemini"
}

type ModelRouterConfig struct {
	DefaultModel    string `json:"default_model"`    // 💬 primary chat model
	AgentModel      string `json:"agent_model"`      // 🤖 agent mode model
	DelegationModel string `json:"delegation_model"` // 🎯 delegation/subagent model
	PremiumModel    string `json:"premium_model"`    // 💎 complex tasks (optional)
	Models       map[string]ModelConfig `json:"models"`
	RoutingRules map[string]string      `json:"routing_rules"` // taskType → modelName

	// Phase 2: Fallback chain configuration
	FallbackModels []string `json:"fallback_models,omitempty"`      // ordered list of model names to try
	FallbackOnError []string `json:"fallback_on_error,omitempty"`   // error types that trigger fallback: "rate_limit", "timeout", "server_error", "auth_error"

	// Custom providers (user-defined via /model add)
	CustomProviders map[string]CustomProvider `json:"custom_providers,omitempty"`
}

// Usage tracking
type ModelUsage struct {
	Model        string    `json:"model"`
	InputTokens  int       `json:"input_tokens"`
	OutputTokens int       `json:"output_tokens"`
	Calls        int       `json:"calls"`
	LastUsed     time.Time `json:"last_used"`
}

var (
	ModelCfg       *ModelRouterConfig
	ModelCfgMu     sync.RWMutex
	ModelUsageMap  map[string]*ModelUsage
	ModelUsageMu   sync.Mutex
)

// ──────────────────────────────────────────────
// Config loading (using ConfigManager)
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
	if ModelCfg == nil {
		ModelCfg = defaultModelConfig()
		return
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
		DefaultModel: "",
		AgentModel:   "",
		PremiumModel: "",
		Models:       map[string]ModelConfig{},
		RoutingRules: map[string]string{},
	}
}

// ──────────────────────────────────────────────
// Model routing
// ──────────────────────────────────────────────

func RouteModel(taskType string) *ModelConfig {
	ModelCfgMu.RLock()
	defer ModelCfgMu.RUnlock()

	if ModelCfg == nil {
		return nil
	}

	// Check routing rules first
	var modelName string
	if name, ok := ModelCfg.RoutingRules[taskType]; ok {
		modelName = name
	} else {
		// Fallback by task type
		switch taskType {
		case "agent":
			modelName = ModelCfg.AgentModel
		case "complex":
			modelName = ModelCfg.PremiumModel
		default:
			modelName = ModelCfg.DefaultModel
		}
	}

	if m, ok := ModelCfg.Models[modelName]; ok {
		return &m
	}

	// Final fallback: return first available
	for _, m := range ModelCfg.Models {
		mc := m
		return &mc
	}
	return nil
}

func GetModelByName(name string) *ModelConfig {
	ModelCfgMu.RLock()
	defer ModelCfgMu.RUnlock()
	if ModelCfg == nil {
		return nil
	}
	if m, ok := ModelCfg.Models[name]; ok {
		return &m
	}
	return nil
}

// ──────────────────────────────────────────────
// Unified OpenAI-compatible API caller
// ──────────────────────────────────────────────

type ChatRequest struct {
	Model       string          `json:"model"`
	Messages    []ChatMessage   `json:"messages"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	Temperature float64         `json:"temperature,omitempty"`
	Stream      bool            `json:"stream"`
	Tools       []registry.ToolSchema       `json:"tools,omitempty"`
	ToolChoice  string          `json:"tool_choice,omitempty"`
}


// toolCallResp represents a native tool call from the API response
type ToolCallResp struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"` // JSON string
	} `json:"function"`
}

type ChatMessage struct {
	Role      string         `json:"role"`
	Content   string         `json:"content"`
	ToolCalls []ToolCallResp `json:"tool_calls,omitempty"`
}

type ChatResponse struct {
	ID      string `json:"id"`
	Choices []struct {
		Message struct {
			Role      string         `json:"role"`
			Content   string         `json:"content"`
			ToolCalls []ToolCallResp `json:"tool_calls"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}

// callModel sends a chat completion request, dispatching to the correct API format.
func CallModel(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, error) {
	if model == nil {
		return "", fmt.Errorf("no model configured")
	}

	switch ResolveAPIFormat(model) {
	case "anthropic":
		return callAnthropic(ctx, model, messages)
	case "gemini":
		return callGemini(ctx, model, messages)
	default:
		return CallOpenAI(ctx, model, messages)
	}
}

// callOpenAI sends a chat completion request to an OpenAI-compatible API.
func CallOpenAI(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, error) {
	// Resolve API key (4-tier: key_env → preset → generic → deprecated inline)
	apiKey := ResolveAPIKey(model)
	if apiKey == "" {
		return "", fmt.Errorf("no API key for provider '%s' — set %s",
			model.Provider, KeySourceLabel(model))
	}

	maxTokens := model.MaxTokens
	if maxTokens == 0 {
		maxTokens = 4096
	}

	reqBody := ChatRequest{
		Model:       model.Model,
		Messages:    messages,
		MaxTokens:   maxTokens,
		Temperature: 0.7,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal error: %w", err)
	}

	endpoint := strings.TrimRight(model.BaseURL, "/") + "/chat/completions"

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(jsonData))
	if err != nil {
		return "", fmt.Errorf("request error: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Provider-specific headers
	if model.Provider == "openrouter" {
		req.Header.Set("HTTP-Referer", "https://scorp-agent.local")
		req.Header.Set("X-Title", "ScorpAgent")
	}

	// Use per-provider transport pool
	client := GetAIClient(model.BaseURL)
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("API call failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read error: %w", err)
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("API error (HTTP %d): %s", resp.StatusCode, helpers.TruncateStr(string(body), 300))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("parse error: %s", helpers.TruncateStr(string(body), 200))
	}

	if chatResp.Error != nil {
		return "", fmt.Errorf("API error: %s", chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no response choices")
	}

	content := chatResp.Choices[0].Message.Content

	// Debug: log response details if content is empty but no error
	if content == "" {
		log.Printf("[models] ⚠️ Empty content from %s — finish=%s, choices=%d, usage=%+v, body_preview=%s",
			model.Model, chatResp.Choices[0].FinishReason, len(chatResp.Choices), chatResp.Usage,
			helpers.TruncateStr(string(body), 500))
	}

	// Track usage + cost
	TrackModelUsage(model.Model, chatResp.Usage.PromptTokens, chatResp.Usage.CompletionTokens)
	RecordCost(model.Model, chatResp.Usage.PromptTokens, chatResp.Usage.CompletionTokens)

	return content, nil
}

// ──────────────────────────────────────────────
// Streaming support (SSE)
// ──────────────────────────────────────────────

// StreamChunk represents a single token chunk from streaming response
type StreamChunk struct {
	Content   string
	Finish    bool
	ToolCalls []ToolCallResp
	Error     error
}

// callModelStream streams chat completion from an OpenAI-compatible API.
// Returns a channel that yields StreamChunk for each token.
// Caller must drain the channel completely.
// Note: Does not support native tool calls in streaming mode yet.
func CallModelStream(ctx context.Context, model *ModelConfig, messages []ChatMessage) (<-chan StreamChunk, error) {
	if model == nil {
		return nil, fmt.Errorf("no model configured")
	}

	// For non-OpenAI formats, fall back to non-streaming and emit as single chunk.
	apiFormat := ResolveAPIFormat(model)
	if apiFormat != "openai" {
		ch := make(chan StreamChunk, 2)
		go func() {
			defer close(ch)
			content, err := CallModel(ctx, model, messages)
			if err != nil {
				ch <- StreamChunk{Error: err}
				return
			}
			ch <- StreamChunk{Content: content}
			ch <- StreamChunk{Finish: true}
		}()
		return ch, nil
	}

	// Resolve API key
	apiKey := ResolveAPIKey(model)
	if apiKey == "" {
		return nil, fmt.Errorf("no API key for provider '%s' — set %s",
			model.Provider, KeySourceLabel(model))
	}

	maxTokens := model.MaxTokens
	if maxTokens == 0 {
		maxTokens = 4096
	}

	reqBody := ChatRequest{
		Model:       model.Model,
		Messages:    messages,
		MaxTokens:   maxTokens,
		Temperature: 0.7,
		Stream:      true, // Enable streaming
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	endpoint := strings.TrimRight(model.BaseURL, "/") + "/chat/completions"

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "text/event-stream")

	// Provider-specific headers
	if model.Provider == "openrouter" {
		req.Header.Set("HTTP-Referer", "https://scorp-agent.local")
		req.Header.Set("X-Title", "ScorpAgent")
	}

	client := GetAIClient(model.BaseURL)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API call failed: %w", err)
	}

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("API error (HTTP %d): %s", resp.StatusCode, helpers.TruncateStr(string(body), 300))
	}

	// Channel for streaming chunks
	ch := make(chan StreamChunk, 16)

	go func() {
		defer resp.Body.Close()
		defer close(ch)

		scanner := bufio.NewScanner(resp.Body)
		buf := make([]byte, 0, 64*1024)
		scanner.Buffer(buf, 1024*1024)

		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				continue
			}

			// SSE format: "data: {...}"
			if !strings.HasPrefix(line, "data: ") {
				continue
			}

			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				ch <- StreamChunk{Finish: true}
				return
			}

			var streamResp struct {
				ID      string `json:"id"`
				Choices []struct {
					Delta struct {
						Content   string         `json:"content"`
						ToolCalls []ToolCallResp `json:"tool_calls"`
					} `json:"delta"`
					FinishReason string `json:"finish_reason"`
				} `json:"choices"`
				Usage struct {
					PromptTokens     int `json:"prompt_tokens"`
					CompletionTokens int `json:"completion_tokens"`
					TotalTokens      int `json:"total_tokens"`
				} `json:"usage"`
				Error *struct {
					Message string `json:"message"`
					Type    string `json:"type"`
				} `json:"error,omitempty"`
			}

			if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
				continue // Skip malformed chunks
			}

			if streamResp.Error != nil {
				ch <- StreamChunk{Error: fmt.Errorf("API error: %s", streamResp.Error.Message)}
				return
			}

			if len(streamResp.Choices) > 0 {
				delta := streamResp.Choices[0].Delta
				if delta.Content != "" {
					ch <- StreamChunk{Content: delta.Content}
				}
				if len(delta.ToolCalls) > 0 {
					ch <- StreamChunk{ToolCalls: delta.ToolCalls}
				}
				if streamResp.Choices[0].FinishReason != "" {
					// Track usage from final chunk
					if streamResp.Usage.TotalTokens > 0 {
						TrackModelUsage(model.Model, streamResp.Usage.PromptTokens, streamResp.Usage.CompletionTokens)
						RecordCost(model.Model, streamResp.Usage.PromptTokens, streamResp.Usage.CompletionTokens)
					}
					ch <- StreamChunk{Finish: true}
					return
				}
			}
		}

		if err := scanner.Err(); err != nil {
			ch <- StreamChunk{Error: fmt.Errorf("stream read error: %w", err)}
		}
	}()

	return ch, nil
}

// callModelWithFallback tries models in order until one succeeds.
// Phase 2: Supports N-tier fallback chain from config (fallback_models).
// Phase 3 fix: Uses *ModelConfig pointers directly to avoid model-key mismatch.
func CallModelWithFallback(ctx context.Context, taskType string, messages []ChatMessage) (string, string, error) {
	ModelCfgMu.RLock()
	cfg := ModelCfg
	ModelCfgMu.RUnlock()

	if cfg == nil {
		return "", "", fmt.Errorf("no model config loaded")
	}

	// Build ordered list of models to try (as pointers, not name strings)
	var models []*ModelConfig
	var primaryLabel string

	// 1. Primary model (cost-aware)
	primary := RouteModelCostAware(taskType)
	if primary != nil {
		primaryLabel = primary.Model
		models = append(models, primary)
	}

	// 2. Fallback models from config (looked up by map key)
	for _, name := range cfg.FallbackModels {
		m := GetModelByName(name)
		if m == nil {
			log.Printf("[models] Fallback model '%s' not found in config, skipping", name)
			continue
		}
		// Skip duplicates
		dup := false
		for _, existing := range models {
			if existing.Model == m.Model {
				dup = true
				break
			}
		}
		if !dup {
			models = append(models, m)
		}
	}

	// 3. Default chat model as last resort
	if dm := RouteModel("chat"); dm != nil {
		dup := false
		for _, existing := range models {
			if existing.Model == dm.Model {
				dup = true
				break
			}
		}
		if !dup {
			models = append(models, dm)
		}
	}

	if len(models) == 0 {
		return "", "", fmt.Errorf("no models available")
	}

	// Try each model in order
	var lastErr error
	for i, model := range models {
		if i > 0 {
			log.Printf("[models] Trying fallback model %s (%d/%d)", model.Model, i+1, len(models))
		}

		result, err := CallModel(ctx, model, messages)
		if err == nil {
			if i == 0 {
				return result, model.Model, nil
			}
			return result, model.Model + " (fallback)", nil
		}

		lastErr = err
		log.Printf("[models] Model %s failed: %v", model.Model, err)

		// Check if error type should trigger fallback (configurable)
		if !ShouldFallbackOnError(err, cfg.FallbackOnError) {
			log.Printf("[models] Error type does not match fallback triggers, stopping chain")
			break
		}
	}

	// All models failed
	return "", "", fmt.Errorf("all models failed (primary: %s): %w", primaryLabel, lastErr)
}

// shouldFallbackOnError checks if an error matches the configured fallback triggers.
func ShouldFallbackOnError(err error, triggers []string) bool {
	if err == nil || len(triggers) == 0 {
		return true // fallback on any error if no specific triggers configured
	}
	msg := strings.ToLower(err.Error())
	for _, trigger := range triggers {
		trigger = strings.ToLower(trigger)
		switch trigger {
		case "rate_limit", "rate-limit":
			if strings.Contains(msg, "rate limit") || strings.Contains(msg, "rate_limit") ||
				strings.Contains(msg, "http 429") || strings.Contains(msg, "too many requests") ||
				strings.Contains(msg, "quota") {
				return true
			}
		case "timeout":
			if strings.Contains(msg, "timeout") || strings.Contains(msg, "deadline exceeded") ||
				strings.Contains(msg, "context deadline") || strings.Contains(msg, "i/o timeout") {
				return true
			}
		case "server_error", "server-error":
			if strings.Contains(msg, "http 5") || strings.Contains(msg, "server error") ||
				strings.Contains(msg, "internal error") || strings.Contains(msg, "service unavailable") {
				return true
			}
		case "auth_error", "auth-error":
			if strings.Contains(msg, "unauthorized") || strings.Contains(msg, "http 401") ||
				strings.Contains(msg, "http 403") || strings.Contains(msg, "invalid api key") ||
				strings.Contains(msg, "authentication") {
				return true
			}
		case "network_error", "network-error":
			if strings.Contains(msg, "connection refused") || strings.Contains(msg, "no such host") ||
				strings.Contains(msg, "network") || strings.Contains(msg, "dial tcp") {
				return true
			}
		}
	}
	return false
}

// ──────────────────────────────────────────────
// Usage tracking
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
	u.Calls++
	u.LastUsed = time.Now()

	// Save
	path := config.ModelUsageFilePath()
	data, _ := json.MarshalIndent(ModelUsageMap, "", "  ")
	os.WriteFile(path, data, 0644)
}

// ──────────────────────────────────────────────
// Display helpers
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

// checkModelHealth sends a tiny test request to verify the API key works
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

// formatModelListWithHealth checks all models and returns with health status
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
		// model: {input_per_1M, output_per_1M}
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
			cost = (float64(u.InputTokens)/1e6)*p[0] + (float64(u.OutputTokens)/1e6)*p[1]
		}
		totalCost += cost

		costStr := "FREE"
		if cost > 0 {
			costStr = fmt.Sprintf("$%.4f", cost)
		}

		sb.WriteString(fmt.Sprintf("<b>%s</b>\n", u.Model))
		sb.WriteString(fmt.Sprintf("  Calls: %d | In: %d | Out: %d\n", u.Calls, u.InputTokens, u.OutputTokens))
		sb.WriteString(fmt.Sprintf("  Cost: %s | Last: %s\n\n", costStr, u.LastUsed.Format("01/02 15:04")))
	}

	if totalCost > 0 {
		sb.WriteString(fmt.Sprintf("💰 <b>Total estimated cost: $%.4f</b>", totalCost))
	} else {
		sb.WriteString("💰 <b>Total cost: $0 (all free models!)</b>")
	}

	return sb.String()
}

// switchModel changes the default or agent model
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
	case "premium":
		ModelCfg.PremiumModel = modelName
	default:
		ModelCfg.DefaultModel = modelName
	}

	SaveModelConfig()
	return fmt.Sprintf("✅ %s model set to <code>%s</code>", role, modelName)
}
