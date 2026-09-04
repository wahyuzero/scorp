package models

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"scorp-agent/internal/helpers"
	"scorp-agent/registry"
)

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
		case "vision":
			if ModelCfg.VisionModel != "" {
				modelName = ModelCfg.VisionModel
			} else {
				modelName = findFirstVisionModel()
			}
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

// findFirstVisionModel returns the best available vision-capable model in config
func findFirstVisionModel() string {
	if ModelCfg == nil {
		return ""
	}
	// 1. High-priority known vision models
	candidates := []string{
		"opencode/mimo-v2.5-free",
		"gpt-5.6-luna",
		"z-ai/glm-5.3-flash",
		"gpt-4o",
		"gpt-4o-mini",
		"claude-3-5-sonnet",
	}
	for _, c := range candidates {
		if _, ok := ModelCfg.Models[c]; ok {
			return c
		}
	}
	// 2. Keyword heuristic
	for name := range ModelCfg.Models {
		lower := strings.ToLower(name)
		if strings.Contains(lower, "vision") || strings.Contains(lower, "luna") ||
			strings.Contains(lower, "omni") || strings.Contains(lower, "4o") ||
			strings.Contains(lower, "mimo") || strings.Contains(lower, "gemini") {
			return name
		}
	}
	return ModelCfg.DefaultModel
}

// isVisionModelName checks if a model name/ID supports vision/multimodal input
func isVisionModelName(name string) bool {
	lower := strings.ToLower(name)
	return strings.Contains(lower, "vision") || strings.Contains(lower, "luna") ||
		strings.Contains(lower, "omni") || strings.Contains(lower, "4o") ||
		strings.Contains(lower, "glm-5.3") || strings.Contains(lower, "mimo") ||
		strings.Contains(lower, "sonnet") || strings.Contains(lower, "gemini")
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

// SwitchActiveModel sets the active default and agent model, updating routing rules and saving config
func SwitchActiveModel(name string) error {
	ModelCfgMu.Lock()
	defer ModelCfgMu.Unlock()
	if ModelCfg == nil {
		return fmt.Errorf("model config not loaded")
	}
	if _, ok := ModelCfg.Models[name]; !ok {
		return fmt.Errorf("model '%s' not found in config", name)
	}
	ModelCfg.DefaultModel = name
	ModelCfg.AgentModel = name
	if ModelCfg.RoutingRules == nil {
		ModelCfg.RoutingRules = make(map[string]string)
	}
	ModelCfg.RoutingRules["agent"] = name
	ModelCfg.RoutingRules["chat"] = name
	SaveModelConfig()
	log.Printf("[models] Switched active model to: %s", name)
	return nil
}

// ──────────────────────────────────────────────
// Unified OpenAI-compatible API caller
// ──────────────────────────────────────────────

type ChatRequest struct {
	Model       string                `json:"model"`
	Messages    []ChatMessage         `json:"messages"`
	MaxTokens   int                   `json:"max_tokens,omitempty"`
	Temperature float64               `json:"temperature,omitempty"`
	Stream      bool                  `json:"stream"`
	Tools       []registry.ToolSchema `json:"tools,omitempty"`
	ToolChoice  string                `json:"tool_choice,omitempty"`
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
		PromptTokens        int `json:"prompt_tokens"`
		CompletionTokens    int `json:"completion_tokens"`
		TotalTokens         int `json:"total_tokens"`
		PromptTokensDetails *struct {
			CachedTokens int `json:"cached_tokens"`
		} `json:"prompt_tokens_details,omitempty"`
	} `json:"usage"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}

// CallModel sends a chat completion request, dispatching to the modular LLMProvider.
func CallModel(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, error) {
	if model == nil {
		return "", fmt.Errorf("no model configured")
	}

	provider := GetProvider(ResolveAPIFormat(model))
	return provider.Call(ctx, model, messages)
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

	// For non-OpenAI formats, dispatch or fall back
	apiFormat := ResolveAPIFormat(model)
	if apiFormat == "command-code" || apiFormat == "commandcode" {
		return CallCommandCodeStream(ctx, model, messages)
	}
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
	} else if model.Provider == "opencode" || model.Provider == "opencode-zen" {
		req.Header.Set("User-Agent", "opencode/1.0.0")
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
// Supports N-tier fallback chain from config (fallback_models).
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

	// For vision tasks, ensure only vision-capable models are in the chain
	if taskType == "vision" {
		var visionModels []*ModelConfig
		if primary != nil && isVisionModelName(primary.Model) {
			visionModels = append(visionModels, primary)
		}
		// Add known vision candidates from config
		visionCandidates := []string{"opencode/mimo-v2.5-free", "gpt-5.6-luna", "z-ai/glm-5.3-flash", "gpt-4o", "gpt-4o-mini"}
		for _, name := range visionCandidates {
			if m := GetModelByName(name); m != nil {
				dup := false
				for _, existing := range visionModels {
					if existing.Model == m.Model {
						dup = true
						break
					}
				}
				if !dup {
					visionModels = append(visionModels, m)
				}
			}
		}
		if len(visionModels) > 0 {
			models = visionModels
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
	// Always allow fallback if upstream rejected parameters or model failed
	if strings.Contains(msg, "invalid request") || strings.Contains(msg, "upstream request failed") || strings.Contains(msg, "http 400") {
		return true
	}
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
		case "insufficient_balance", "insufficient-quota":
			if strings.Contains(msg, "402") || strings.Contains(msg, "insufficient") ||
				strings.Contains(msg, "quota") || strings.Contains(msg, "balance") {
				return true
			}
		default:
			// Custom trigger from config — match as substring in error message
			if strings.Contains(msg, trigger) {
				return true
			}
		}
	}
	return false
}
