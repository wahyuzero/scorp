package models

import (
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
// OpenAI-Compatible Provider
// Handles: OpenAI, DeepSeek, Groq, OpenRouter,
// Zhipu (Z.ai), Ollama, vLLM, LM Studio, etc.
// ──────────────────────────────────────────────

type OpenAIProvider struct{}

func (p *OpenAIProvider) Format() string {
	return "openai"
}

func (p *OpenAIProvider) Call(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, error) {
	return CallOpenAI(ctx, model, messages)
}

func (p *OpenAIProvider) CallWithTools(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, []ToolCall, error) {
	return CallOpenAIWithTools(ctx, model, messages)
}

// formatOpenAIMessages converts ChatMessage list to OpenAI-compatible messages,
// unpacking multimodal array content (e.g. vision image_url) if present.
func formatOpenAIMessages(messages []ChatMessage) []map[string]interface{} {
	var formatted []map[string]interface{}
	for _, m := range messages {
		msgMap := map[string]interface{}{
			"role": m.Role,
		}
		// Check if content is a JSON array of multimodal content parts
		trimmed := strings.TrimSpace(m.Content)
		if strings.HasPrefix(trimmed, "[") && strings.Contains(trimmed, "image") {
			var parts []interface{}
			if err := json.Unmarshal([]byte(trimmed), &parts); err == nil {
				msgMap["content"] = parts
			} else {
				msgMap["content"] = m.Content
			}
		} else {
			msgMap["content"] = m.Content
		}

		if len(m.ToolCalls) > 0 {
			msgMap["tool_calls"] = m.ToolCalls
		}
		formatted = append(formatted, msgMap)
	}
	return formatted
}

// CallOpenAI sends a plain chat completion request to an OpenAI-compatible API.
func CallOpenAI(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, error) {
	apiKey := ResolveAPIKey(model)
	if apiKey == "" {
		return "", fmt.Errorf("no API key for provider '%s' — set %s",
			model.Provider, KeySourceLabel(model))
	}

	maxTokens := model.MaxTokens
	if maxTokens == 0 {
		maxTokens = 4096
	}

	reqBody := map[string]interface{}{
		"model":       model.Model,
		"messages":    formatOpenAIMessages(messages),
		"max_tokens":  maxTokens,
		"temperature": 0.7,
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
	} else if model.Provider == "opencode" || model.Provider == "opencode-zen" {
		req.Header.Set("User-Agent", "opencode/1.0.0")
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

	reply := chatResp.Choices[0].Message.Content

	// Track usage + cost
	cachedTokens := 0
	if chatResp.Usage.PromptTokensDetails != nil {
		cachedTokens = chatResp.Usage.PromptTokensDetails.CachedTokens
	}
	TrackModelUsageWithCache(model.Model, chatResp.Usage.PromptTokens, chatResp.Usage.CompletionTokens, cachedTokens)
	RecordCostWithCache(model.Model, chatResp.Usage.PromptTokens, chatResp.Usage.CompletionTokens, cachedTokens)

	return reply, nil
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

	reqBody := map[string]interface{}{
		"model":       model.Model,
		"messages":    formatOpenAIMessages(messages),
		"max_tokens":  maxTokens,
		"temperature": 0.7,
		"tools":       registry.GenerateNativeToolsSchema(),
		"tool_choice": "auto",
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
	} else if model.Provider == "opencode" || model.Provider == "opencode-zen" {
		httpReq.Header.Set("User-Agent", "opencode/1.0.0")
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

	// Track usage + cost
	cachedTokensWithTools := 0
	if chatResp.Usage.PromptTokensDetails != nil {
		cachedTokensWithTools = chatResp.Usage.PromptTokensDetails.CachedTokens
	}
	TrackModelUsageWithCache(model.Model, chatResp.Usage.PromptTokens, chatResp.Usage.CompletionTokens, cachedTokensWithTools)
	RecordCostWithCache(model.Model, chatResp.Usage.PromptTokens, chatResp.Usage.CompletionTokens, cachedTokensWithTools)

	return content, toolCalls, nil
}
