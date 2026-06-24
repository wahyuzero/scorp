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
	"strings"
)

// ──────────────────────────────────────────────
// Anthropic API (Claude) — /v1/messages format
// ──────────────────────────────────────────────

// anthropicRequest is the request body for the Anthropic Messages API.
type anthropicRequest struct {
	Model     string          `json:"model"`
	MaxTokens int             `json:"max_tokens"`
	Messages  []ChatMessage   `json:"messages"`
	System    string          `json:"system,omitempty"`
	Tools     []anthropicTool `json:"tools,omitempty"`
}

type anthropicTool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"input_schema"`
}

type anthropicResponse struct {
	Content []struct {
		Type  string `json:"type"`
		Text  string `json:"text,omitempty"`
		ID    string `json:"id,omitempty"`
		Name  string `json:"name,omitempty"`
		Input json.RawMessage `json:"input,omitempty"`
	} `json:"content"`
	StopReason string `json:"stop_reason"`
	Usage      struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}

// callAnthropic sends a chat completion request to an Anthropic-compatible API.
func callAnthropic(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, error) {
	apiKey := ResolveAPIKey(model)
	if apiKey == "" {
		return "", fmt.Errorf("no API key for provider '%s' — set %s",
			model.Provider, KeySourceLabel(model))
	}

	maxTokens := model.MaxTokens
	if maxTokens == 0 {
		maxTokens = 4096
	}

	// Extract system message (Anthropic puts it at top-level, not in messages)
	systemMsg := ""
	var filtered []ChatMessage
	for _, m := range messages {
		if m.Role == "system" {
			if systemMsg != "" {
				systemMsg += "\n\n"
			}
			systemMsg += m.Content
		} else {
			filtered = append(filtered, m)
		}
	}

	reqBody := anthropicRequest{
		Model:     model.Model,
		MaxTokens: maxTokens,
		Messages:  filtered,
		System:    systemMsg,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal error: %w", err)
	}

	endpoint := strings.TrimRight(model.BaseURL, "/") + "/v1/messages"

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(jsonData))
	if err != nil {
		return "", fmt.Errorf("request error: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

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

	var apiResp anthropicResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return "", fmt.Errorf("parse error: %s", helpers.TruncateStr(string(body), 200))
	}

	if apiResp.Error != nil {
		return "", fmt.Errorf("API error: %s", apiResp.Error.Message)
	}

	// Concatenate all text blocks (skip tool_use blocks)
	var sb strings.Builder
	for _, block := range apiResp.Content {
		if block.Type == "text" {
			sb.WriteString(block.Text)
		}
	}

	TrackModelUsage(model.Model, apiResp.Usage.InputTokens, apiResp.Usage.OutputTokens)
	RecordCost(model.Model, apiResp.Usage.InputTokens, apiResp.Usage.OutputTokens)

	return sb.String(), nil
}

// callAnthropicWithTools sends a request with native tool definitions.
func CallAnthropicWithTools(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, []ToolCall, error) {
	apiKey := ResolveAPIKey(model)
	if apiKey == "" {
		return "", nil, fmt.Errorf("no API key for provider '%s' — set %s", model.Provider, KeySourceLabel(model))
	}

	maxTokens := model.MaxTokens
	if maxTokens == 0 {
		maxTokens = 4096
	}

	// Extract system message
	systemMsg := ""
	var filtered []ChatMessage
	for _, m := range messages {
		if m.Role == "system" {
			if systemMsg != "" {
				systemMsg += "\n\n"
			}
			systemMsg += m.Content
		} else {
			filtered = append(filtered, m)
		}
	}

	// Convert OpenAI tool defs → Anthropic format
	var anthropicTools []anthropicTool
	for _, td := range registry.GenerateNativeToolsSchema() {
		anthropicTools = append(anthropicTools, anthropicTool{
			Name:        td.Function.Name,
			Description: td.Function.Description,
			InputSchema: td.Function.Parameters,
		})
	}

	reqBody := anthropicRequest{
		Model:     model.Model,
		MaxTokens: maxTokens,
		Messages:  filtered,
		System:    systemMsg,
		Tools:     anthropicTools,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", nil, fmt.Errorf("marshal error: %w", err)
	}

	endpoint := strings.TrimRight(model.BaseURL, "/") + "/v1/messages"

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(jsonData))
	if err != nil {
		return "", nil, fmt.Errorf("request error: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := GetAIClient(model.BaseURL).Do(req)
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

	var apiResp anthropicResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return "", nil, fmt.Errorf("parse error: %s", helpers.TruncateStr(string(body), 200))
	}

	if apiResp.Error != nil {
		return "", nil, fmt.Errorf("API error: %s", apiResp.Error.Message)
	}

	// Extract text content and tool_use blocks
	var sb strings.Builder
	var toolCalls []ToolCall
	for _, block := range apiResp.Content {
		switch block.Type {
		case "text":
			sb.WriteString(block.Text)
		case "tool_use":
			var args map[string]interface{}
			if len(block.Input) > 0 {
				if err := json.Unmarshal(block.Input, &args); err != nil {
					log.Printf("[agent] Failed to parse Anthropic tool args: %v", err)
					args = make(map[string]interface{})
				}
			} else {
				args = make(map[string]interface{})
			}
			toolCalls = append(toolCalls, ToolCall{
				Name: block.Name,
				Args: args,
			})
			log.Printf("[agent] Anthropic tool_use: %s(%v)", block.Name, args)
		}
	}

	TrackModelUsage(model.Model, apiResp.Usage.InputTokens, apiResp.Usage.OutputTokens)
	RecordCost(model.Model, apiResp.Usage.InputTokens, apiResp.Usage.OutputTokens)

	return sb.String(), toolCalls, nil
}
