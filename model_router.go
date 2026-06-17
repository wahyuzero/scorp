package main

import (
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
	Provider  string `json:"provider"` // groq, deepseek, gemini, openrouter, 9router
	Model     string `json:"model"`    // model ID
	APIKey    string `json:"api_key"`
	BaseURL   string `json:"base_url"`   // OpenAI-compatible endpoint
	MaxTokens int    `json:"max_tokens"` // max output tokens
}

type ModelRouterConfig struct {
	DefaultModel string                 `json:"default_model"` // for chat
	AgentModel   string                 `json:"agent_model"`   // for agent mode
	PremiumModel string                 `json:"premium_model"` // for complex tasks
	Models       map[string]ModelConfig `json:"models"`
	RoutingRules map[string]string      `json:"routing_rules"` // taskType → modelName
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
	modelCfg     *ModelRouterConfig
	modelCfgMu   sync.RWMutex
	modelUsage   map[string]*ModelUsage
	modelUsageMu sync.Mutex
	modelCfgPath = os.ExpandEnv("$HOME") + "/.scorp-agent/models.json"
)

// ──────────────────────────────────────────────
// Config loading
// ──────────────────────────────────────────────

func loadModelConfig() {
	modelCfgMu.Lock()
	defer modelCfgMu.Unlock()

	data, err := os.ReadFile(modelCfgPath)
	if err != nil {
		log.Printf("[models] No config found, using defaults: %v", err)
		modelCfg = defaultModelConfig()
		saveModelConfig()
		return
	}

	var cfg ModelRouterConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		log.Printf("[models] Parse error: %v, using defaults", err)
		modelCfg = defaultModelConfig()
		return
	}

	modelCfg = &cfg
	log.Printf("[models] Loaded %d models (default=%s, agent=%s)",
		len(cfg.Models), cfg.DefaultModel, cfg.AgentModel)
}

func saveModelConfig() {
	if modelCfg == nil {
		return
	}
	os.MkdirAll(os.ExpandEnv("$HOME")+"/.scorp-agent", 0755)
	data, _ := json.MarshalIndent(modelCfg, "", "  ")
	os.WriteFile(modelCfgPath, data, 0644)
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

func routeModel(taskType string) *ModelConfig {
	modelCfgMu.RLock()
	defer modelCfgMu.RUnlock()

	if modelCfg == nil {
		return nil
	}

	// Check routing rules first
	var modelName string
	if name, ok := modelCfg.RoutingRules[taskType]; ok {
		modelName = name
	} else {
		// Fallback by task type
		switch taskType {
		case "agent":
			modelName = modelCfg.AgentModel
		case "complex":
			modelName = modelCfg.PremiumModel
		default:
			modelName = modelCfg.DefaultModel
		}
	}

	if m, ok := modelCfg.Models[modelName]; ok {
		return &m
	}

	// Final fallback: return first available
	for _, m := range modelCfg.Models {
		mc := m
		return &mc
	}
	return nil
}

func getModelByName(name string) *ModelConfig {
	modelCfgMu.RLock()
	defer modelCfgMu.RUnlock()
	if modelCfg == nil {
		return nil
	}
	if m, ok := modelCfg.Models[name]; ok {
		return &m
	}
	return nil
}

// ──────────────────────────────────────────────
// Unified OpenAI-compatible API caller
// ──────────────────────────────────────────────

type chatRequest struct {
	Model       string          `json:"model"`
	Messages    []chatMessage   `json:"messages"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	Temperature float64         `json:"temperature,omitempty"`
	Stream      bool            `json:"stream"`
	Tools       []toolDef       `json:"tools,omitempty"`
	ToolChoice  string          `json:"tool_choice,omitempty"`
}

type toolDef struct {
	Type     string       `json:"type"` // always "function"
	Function toolFunction `json:"function"`
}

type toolFunction struct {
	Name        string                    `json:"name"`
	Description string                    `json:"description"`
	Parameters  map[string]interface{}    `json:"parameters"`
}

// toolCallResp represents a native tool call from the API response
type toolCallResp struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"` // JSON string
	} `json:"function"`
}

type chatMessage struct {
	Role      string         `json:"role"`
	Content   string         `json:"content"`
	ToolCalls []toolCallResp `json:"tool_calls,omitempty"`
}

type chatResponse struct {
	ID      string `json:"id"`
	Choices []struct {
		Message struct {
			Role      string         `json:"role"`
			Content   string         `json:"content"`
			ToolCalls []toolCallResp `json:"tool_calls"`
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

// callModel sends a chat completion request to any OpenAI-compatible API
func callModel(ctx context.Context, model *ModelConfig, messages []chatMessage) (string, error) {
	if model == nil {
		return "", fmt.Errorf("no model configured")
	}

	// Skip API call for providers with no API key
	if model.APIKey == "" {
		return "", fmt.Errorf("no API key configured for %s", model.Provider)
	}

	maxTokens := model.MaxTokens
	if maxTokens == 0 {
		maxTokens = 4096
	}

	reqBody := chatRequest{
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
	req.Header.Set("Authorization", "Bearer "+model.APIKey)

	// Provider-specific headers
	if model.Provider == "openrouter" {
		req.Header.Set("HTTP-Referer", "https://scorp-agent.local")
		req.Header.Set("X-Title", "ScorpAgent")
	}

	// Use per-provider transport pool
	client := getAIClient(model.BaseURL)
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
		return "", fmt.Errorf("API error (HTTP %d): %s", resp.StatusCode, truncateStr(string(body), 300))
	}

	var chatResp chatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("parse error: %s", truncateStr(string(body), 200))
	}

	if chatResp.Error != nil {
		return "", fmt.Errorf("API error: %s", chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no response choices")
	}

	content := chatResp.Choices[0].Message.Content

	// Track usage + cost
	trackModelUsage(model.Model, chatResp.Usage.PromptTokens, chatResp.Usage.CompletionTokens)
	recordCost(model.Model, chatResp.Usage.PromptTokens, chatResp.Usage.CompletionTokens)

	return content, nil
}

// ──────────────────────────────────────────────
// Streaming support (SSE)
// ──────────────────────────────────────────────

// StreamChunk represents a single token chunk from streaming response
type StreamChunk struct {
	Content   string
	Finish    bool
	ToolCalls []toolCallResp
	Error     error
}

// callModelStream streams chat completion from an OpenAI-compatible API.
// Returns a channel that yields StreamChunk for each token.
// Caller must drain the channel completely.
// Note: Does not support native tool calls in streaming mode yet.
func callModelStream(ctx context.Context, model *ModelConfig, messages []chatMessage) (<-chan StreamChunk, error) {
	if model == nil {
		return nil, fmt.Errorf("no model configured")
	}

	// Non-streaming fallback for models without streaming support
	if model.APIKey == "" {
		return nil, fmt.Errorf("no API key configured for %s", model.Provider)
	}

	maxTokens := model.MaxTokens
	if maxTokens == 0 {
		maxTokens = 4096
	}

	reqBody := chatRequest{
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
	req.Header.Set("Authorization", "Bearer "+model.APIKey)
	req.Header.Set("Accept", "text/event-stream")

	// Provider-specific headers
	if model.Provider == "openrouter" {
		req.Header.Set("HTTP-Referer", "https://scorp-agent.local")
		req.Header.Set("X-Title", "ScorpAgent")
	}

	client := getAIClient(model.BaseURL)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API call failed: %w", err)
	}

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("API error (HTTP %d): %s", resp.StatusCode, truncateStr(string(body), 300))
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
						ToolCalls []toolCallResp `json:"tool_calls"`
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
						trackModelUsage(model.Model, streamResp.Usage.PromptTokens, streamResp.Usage.CompletionTokens)
						recordCost(model.Model, streamResp.Usage.PromptTokens, streamResp.Usage.CompletionTokens)
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

// callModelWithFallback tries models in order until one succeeds
func callModelWithFallback(ctx context.Context, taskType string, messages []chatMessage) (string, string, error) {
	// Primary model (cost-aware)
	primary := routeModelCostAware(taskType)
	if primary != nil {
		result, err := callModel(ctx, primary, messages)
		if err == nil {
			return result, primary.Model, nil
		}
		log.Printf("[models] Primary model %s failed: %v, trying fallback", primary.Model, err)
	}

	// Fallback to default
	fallback := routeModel("chat")
	if fallback != nil && (primary == nil || fallback.Model != primary.Model) {
		result, err := callModel(ctx, fallback, messages)
		if err == nil {
			return result, fallback.Model + " (fallback)", nil
		}
		log.Printf("[models] Fallback model %s failed: %v", fallback.Model, err)
	}

	// All models failed
	return "", "", fmt.Errorf("all models failed (primary: %s)", routeModel(taskType).Model)
}

// ──────────────────────────────────────────────
// Usage tracking
// ──────────────────────────────────────────────

func initModelUsage() {
	modelUsageMu.Lock()
	defer modelUsageMu.Unlock()

	modelUsage = make(map[string]*ModelUsage)

	// Load from file
	path := os.ExpandEnv("$HOME") + "/.scorp-agent/model_usage.json"
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	json.Unmarshal(data, &modelUsage)
}

func trackModelUsage(model string, inputTokens, outputTokens int) {
	modelUsageMu.Lock()
	defer modelUsageMu.Unlock()

	if modelUsage == nil {
		modelUsage = make(map[string]*ModelUsage)
	}

	u, ok := modelUsage[model]
	if !ok {
		u = &ModelUsage{Model: model}
		modelUsage[model] = u
	}

	u.InputTokens += inputTokens
	u.OutputTokens += outputTokens
	u.Calls++
	u.LastUsed = time.Now()

	// Save
	path := os.ExpandEnv("$HOME") + "/.scorp-agent/model_usage.json"
	data, _ := json.MarshalIndent(modelUsage, "", "  ")
	os.WriteFile(path, data, 0644)
}

// ──────────────────────────────────────────────
// Display helpers
// ──────────────────────────────────────────────

func formatModelList() string {
	modelCfgMu.RLock()
	defer modelCfgMu.RUnlock()

	if modelCfg == nil {
		return "❌ No model config loaded."
	}

	var sb strings.Builder
	sb.WriteString("🤖 <b>AI Models</b>\n\n")

	for name, m := range modelCfg.Models {
		isDefault := ""
		if name == modelCfg.DefaultModel {
			isDefault += " 💬chat"
		}
		if name == modelCfg.AgentModel {
			isDefault += " 🤖agent"
		}
		if name == modelCfg.PremiumModel {
			isDefault += " 🧠premium"
		}

		apiStatus := "⏳"
		if m.APIKey == "" {
			apiStatus = "⚠️ no key"
		}

		sb.WriteString(fmt.Sprintf("%s <code>%s</code> — %s (%s)%s\n",
			apiStatus, name, m.Model, m.Provider, isDefault))
	}

	sb.WriteString("\n<b>Routing:</b>\n")
	for task, model := range modelCfg.RoutingRules {
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
func checkModelHealth(name string, m *ModelConfig) (bool, string) {
	if m.APIKey == "" {
		return false, "no API key"
	}

	reqBody := chatRequest{
		Model:     m.Model,
		Messages:  []chatMessage{{Role: "user", Content: "hi"}},
		MaxTokens: 5,
	}

	jsonData, _ := json.Marshal(reqBody)
	endpoint := strings.TrimRight(m.BaseURL, "/") + "/chat/completions"

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(jsonData))
	if err != nil {
		return false, err.Error()
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+m.APIKey)

	client := httpShort
	resp, err := client.Do(req)
	if err != nil {
		return false, truncateStr(err.Error(), 60)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return true, "ok"
	}

	body, _ := io.ReadAll(resp.Body)
	return false, fmt.Sprintf("HTTP %d: %s", resp.StatusCode, truncateStr(string(body), 80))
}

// formatModelListWithHealth checks all models and returns with health status
func formatModelListWithHealth() string {
	modelCfgMu.RLock()
	cfgCopy := *modelCfg
	modelCfgMu.RUnlock()

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
			ok, detail := checkModelHealth(n, &mc)
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

func formatUsageStats() string {
	modelUsageMu.Lock()
	defer modelUsageMu.Unlock()

	if len(modelUsage) == 0 {
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
	for _, u := range modelUsage {
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
func switchModel(role, modelName string) string {
	modelCfgMu.Lock()
	defer modelCfgMu.Unlock()

	if modelCfg == nil {
		return "❌ No model config loaded."
	}

	if _, ok := modelCfg.Models[modelName]; !ok {
		names := make([]string, 0)
		for n := range modelCfg.Models {
			names = append(names, n)
		}
		return fmt.Sprintf("❌ Model '%s' not found. Available: %s", modelName, strings.Join(names, ", "))
	}

	switch role {
	case "default", "chat", "use":
		modelCfg.DefaultModel = modelName
		modelCfg.RoutingRules["chat"] = modelName
	case "agent":
		modelCfg.AgentModel = modelName
		modelCfg.RoutingRules["agent"] = modelName
	case "premium":
		modelCfg.PremiumModel = modelName
	default:
		modelCfg.DefaultModel = modelName
	}

	saveModelConfig()
	return fmt.Sprintf("✅ %s model set to <code>%s</code>", role, modelName)
}
