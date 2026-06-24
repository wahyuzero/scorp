package models

import (
	"scorp-agent/internal/helpers"
	"scorp-agent/registry"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// ──────────────────────────────────────────────
// Shared Model Calling Functions
// ──────────────────────────────────────────────

// CallModelWithTools calls a model with native tool definitions (OpenAI/Anthropic/Gemini)
func CallModelWithTools(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, []ToolCall, error) {
	if model == nil {
		return "", nil, fmt.Errorf("no model configured")
	}

	switch ResolveAPIFormat(model) {
	case "anthropic":
		return CallAnthropicWithTools(ctx, model, messages)
	case "gemini":
		return CallGeminiWithTools(ctx, model, messages)
	default:
		return CallOpenAIWithTools(ctx, model, messages)
	}
}

// CallOpenAIWithTools sends a chat completion with native tool definitions to an OpenAI-compatible API.
func CallOpenAIWithTools(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, []ToolCall, error) {
	apiKey := ResolveAPIKey(model)
	if apiKey == "" {
		return "", nil, fmt.Errorf("no API key for provider '%s' — %s", model.Provider, KeySourceLabel(model))
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
		Tools:       registry.GenerateNativeToolsSchema(),
		ToolChoice:  "auto",
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", nil, fmt.Errorf("marshal error: %w", err)
	}

	endpoint := strings.TrimRight(model.BaseURL, "/") + "/chat/completions"

	httpReq, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(jsonData))
	if err != nil {
		return "", nil, fmt.Errorf("request error: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	if model.Provider == "openrouter" {
		httpReq.Header.Set("HTTP-Referer", "https://scorp-agent.local")
		httpReq.Header.Set("X-Title", "ScorpAgent")
	}

	resp, err := GetAIClient(model.BaseURL).Do(httpReq)
	if err != nil {
		return "", nil, fmt.Errorf("API call failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, fmt.Errorf("read error: %w", err)
	}

	if resp.StatusCode != 200 {
		return "", nil, fmt.Errorf("API error (HTTP %d): %s", resp.StatusCode, helpers.TruncateStr(string(body), 300))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", nil, fmt.Errorf("parse error: %s", helpers.TruncateStr(string(body), 200))
	}

	if chatResp.Error != nil {
		return "", nil, fmt.Errorf("API error: %s", chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 {
		return "", nil, fmt.Errorf("no response choices")
	}

	choice := chatResp.Choices[0]
	content := choice.Message.Content

	// Parse native tool calls
	var toolCalls []ToolCall
	for _, tc := range choice.Message.ToolCalls {
		var args map[string]interface{}
		if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err != nil {
			log.Printf("[agent] Failed to parse tool args '%s': %v", tc.Function.Arguments, err)
			args = make(map[string]interface{})
		}
		toolCalls = append(toolCalls, ToolCall{
			Name: tc.Function.Name,
			Args: args,
		})
	}

	return content, toolCalls, nil
}



// IsRateLimitError checks if an error indicates rate limiting (429/402)
func IsRateLimitError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "HTTP 429") ||
		strings.Contains(msg, "HTTP 402") ||
		strings.Contains(msg, "rate_limit") ||
		strings.Contains(msg, "rate limit") ||
		strings.Contains(msg, "quota") ||
		strings.Contains(msg, "suspicious activity")
}



// CallModelWithToolsAndFallback tries with tools first, falls back to plain models.CallModel.
// On rate-limit errors (429/402), retries with exponential backoff before giving up.
// Phase 2: Supports N-tier fallback chain from config (fallback_models).
// Phase 3 fix: Uses *models.ModelConfig pointers directly to avoid model-key mismatch.
func CallModelWithToolsAndFallback(ctx context.Context, taskType string, messages []ChatMessage) (string, []ToolCall, string, error) {
	ModelCfgMu.RLock()
	cfg := ModelCfg
	ModelCfgMu.RUnlock()

	if cfg == nil {
		return "", nil, "", fmt.Errorf("no model config loaded")
	}

	// Build ordered list of models to try (as pointers, not name strings)
	var modelList []*ModelConfig
	var primaryLabel string

	// 1. Primary model (cost-aware)
	primary := RouteModelCostAware(taskType)
	if primary != nil {
		primaryLabel = primary.Model
		modelList = append(modelList, primary)
	}

	// 2. Fallback models from config (looked up by map key)
	for _, name := range cfg.FallbackModels {
		m := GetModelByName(name)
		if m == nil {
			log.Printf("[models] Fallback model '%s' not found in config, skipping", name)
			continue
		}
		dup := false
		for _, existing := range modelList {
			if existing.Model == m.Model {
				dup = true
				break
			}
		}
		if !dup {
			modelList = append(modelList, m)
		}
	}

	// 3. Default chat model as last resort
	if dm := RouteModel("chat"); dm != nil {
		dup := false
		for _, existing := range modelList {
			if existing.Model == dm.Model {
				dup = true
				break
			}
		}
		if !dup {
			modelList = append(modelList, dm)
		}
	}

	if len(modelList) == 0 {
		return "", nil, "", fmt.Errorf("no models available")
	}

	// Try each model in order
	var lastErr error
	for i, model := range modelList {
		label := model.Model

		if i > 0 {
			log.Printf("[models] Trying fallback model %s (%d/%d) with tools", label, i+1, len(modelList))
		}

		// Rate-limit retry loop for this model
		maxRateLimitRetries := 3
		backoffSteps := []time.Duration{15 * time.Second, 30 * time.Second, 60 * time.Second}
		var modelLastErr error

		for attempt := 0; attempt <= maxRateLimitRetries; attempt++ {
			// Try with tools first
			text, toolCalls, err := CallModelWithTools(ctx, model, messages)
			if err == nil {
				if i == 0 {
					return text, toolCalls, label, nil
				}
				return text, toolCalls, label + " (fallback)", nil
			}
			modelLastErr = err



			if !IsRateLimitError(err) {
				// Non-rate-limit error — try plain call, then continue to next fallback
				log.Printf("[models] Model %s with tools failed: %v, trying plain", label, err)
				text2, err2 := CallModel(ctx, model, messages)
				if err2 == nil {
					if i == 0 {
						return text2, nil, label, nil
					}
					return text2, nil, label + " (fallback)", nil
				}
				log.Printf("[models] Model %s plain also failed: %v", label, err2)
				modelLastErr = err2
				break
			}

			// Rate-limit error — backoff and retry
			if attempt < maxRateLimitRetries {
				wait := backoffSteps[attempt]
				log.Printf("[models] ⏳ Rate-limit detected on %s (HTTP 429/402), backing off for %v before retry %d/%d",
					label, wait, attempt+1, maxRateLimitRetries)
				select {
				case <-time.After(wait):
				case <-ctx.Done():
					return "", nil, "", ctx.Err()
				}
			}
		}

		// If we exhausted rate-limit retries, try plain call one more time
		if IsRateLimitError(modelLastErr) {
			log.Printf("[models] Rate-limit retries exhausted for %s, trying plain call once more", label)
			text2, err2 := CallModel(ctx, model, messages)
			if err2 == nil {
				if i == 0 {
					return text2, nil, label, nil
				}
				return text2, nil, label + " (fallback)", nil
			}
			if !IsRateLimitError(err2) {
				modelLastErr = err2
			}
		}

		lastErr = modelLastErr
		log.Printf("[models] Model %s failed after all retries: %v", label, modelLastErr)

		// Check if error type should trigger fallback
		if !ShouldFallbackOnError(lastErr, cfg.FallbackOnError) {
			log.Printf("[models] Error type does not match fallback triggers, stopping chain")
			break
		}
	}

	// All models failed
	return "", nil, "", fmt.Errorf("all models failed (primary: %s): %w", primaryLabel, lastErr)
}



// ──────────────────────────────────────────────
// Code-Block Fallback Parser
// ──────────────────────────────────────────────

var codeBlockRe = regexp.MustCompile("(?s)`{3}(?:shell|bash|sh)?\\s*\\n(.*?)`{3}")
func ParseCodeBlockFallback(text string) ([]ToolCall, string) {
	matches := codeBlockRe.FindAllStringSubmatch(text, -1)
	if len(matches) == 0 {
		return nil, text
	}

	var calls []ToolCall
	cleanText := text

	for _, m := range matches {
		cmd := strings.TrimSpace(m[1])
		if cmd == "" {
			continue
		}

		// Skip if it looks like output (has shell prompt or is too long for a command)
		if strings.Contains(cmd, "\n\n") || len(cmd) > 500 {
			continue
		}

		calls = append(calls, ToolCall{
			Name: "shell",
			Args: map[string]interface{}{
				"command": cmd,
			},
		})
		cleanText = strings.Replace(cleanText, m[0], "", 1)
	}

	return calls, strings.TrimSpace(cleanText)
}

// ParseAllToolCalls tries native tool calls, then XML tags, then code blocks
func ParseAllToolCalls(text string, nativeCalls []ToolCall) ([]ToolCall, string) {
	// 1. Native function calling (already parsed by CallModelWithTools)
	if len(nativeCalls) > 0 {
		return nativeCalls, text
	}

	// 2. XML tag format: <tool_call>{"name": "...", "args": {...}}
	xmlCalls, cleanText := ParseToolCalls(text)
	if len(xmlCalls) > 0 {
		return xmlCalls, cleanText
	}

	// 3. Code-block fallback: parse shell commands from code blocks
	codeCalls, cleanText := ParseCodeBlockFallback(text)
	if len(codeCalls) > 0 {
		return codeCalls, cleanText
	}

	return nil, text
}



// ──────────────────────────────────────────────
// Tool Call Parser (XML format)
// ──────────────────────────────────────────────

var ToolCallRe = regexp.MustCompile(`<tool_call>(.*?)</tool_call>`)

// ParseToolCalls extracts tool calls from LLM response in XML format
func ParseToolCalls(text string) ([]ToolCall, string) {
	matches := ToolCallRe.FindAllStringSubmatch(text, -1)
	if len(matches) == 0 {
		return nil, text
	}

	var calls []ToolCall
	cleanText := text

	for _, m := range matches {
		cleanText = strings.Replace(cleanText, m[0], "", 1)
		var tc ToolCall
		if err := json.Unmarshal([]byte(m[1]), &tc); err != nil {
			log.Printf("[scorp] Failed to parse tool call: %v", err)
			continue
		}
		calls = append(calls, tc)
	}

	return calls, strings.TrimSpace(cleanText)
}

