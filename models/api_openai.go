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
// OpenAI-Compatible Provider
// Handles: OpenAI, DeepSeek, Groq, OpenRouter,
// Zhipu (Z.ai), Ollama, vLLM, LM Studio, etc.
// Enhanced multi-vendor capabilities:
// - MaxTokensField (e.g. max_completion_tokens)
// - CustomHeaders injection
// - ExtraBody injection (e.g. reasoning_split)
// - Tool schema transformation
// - SSE streaming support
// ──────────────────────────────────────────────

type OpenAIProvider struct{}

func init() {
	p := &OpenAIProvider{}
	RegisterProvider(ProviderSpec{
		Name:             "openai",
		DisplayName:      "OpenAI",
		DefaultBaseURL:   "https://api.openai.com/v1",
		DefaultAPIFormat: "openai",
		KeyEnvs:          []string{"OPENAI_API_KEY"},
	}, p)
}

func (p *OpenAIProvider) Format() string {
	return "openai"
}

func (p *OpenAIProvider) Call(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, error) {
	return CallOpenAI(ctx, model, messages)
}

func (p *OpenAIProvider) CallWithTools(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, []ToolCall, error) {
	return CallOpenAIWithTools(ctx, model, messages)
}

func (p *OpenAIProvider) CallStream(ctx context.Context, model *ModelConfig, messages []ChatMessage) (<-chan StreamChunk, error) {
	return CallOpenAIStream(ctx, model, messages)
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

func buildOpenAIRequestBody(model *ModelConfig, messages []ChatMessage, tools []registry.ToolSchema, stream bool) map[string]interface{} {
	maxTokens := model.MaxTokens
	if maxTokens == 0 {
		maxTokens = 4096
	}

	maxTokensField := model.MaxTokensField
	if maxTokensField == "" {
		lowerModel := strings.ToLower(model.Model)
		if strings.Contains(lowerModel, "glm") || strings.Contains(lowerModel, "o1") ||
			strings.Contains(lowerModel, "o3") || strings.Contains(lowerModel, "gpt-5") {
			maxTokensField = "max_completion_tokens"
		} else {
			maxTokensField = "max_tokens"
		}
	}

	reqBody := map[string]interface{}{
		"model":         model.Model,
		"messages":      formatOpenAIMessages(messages),
		maxTokensField: maxTokens,
		"temperature":   0.7,
	}

	if stream {
		reqBody["stream"] = true
	}

	if len(tools) > 0 {
		transformedTools := TransformToolDefinitions(tools, model.ToolSchemaTransform)
		reqBody["tools"] = transformedTools
		reqBody["tool_choice"] = "auto"
	}

	// Merge extra body fields configured per-model/provider (e.g. reasoning_split, thinking, etc.)
	for k, v := range model.ExtraBody {
		reqBody[k] = v
	}

	return reqBody
}

func applyOpenAIHeaders(req *http.Request, model *ModelConfig, apiKey string) {
	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	// Default openrouter headers
	if model.Provider == "openrouter" {
		req.Header.Set("HTTP-Referer", "https://scorp-agent.local")
		req.Header.Set("X-Title", "ScorpAgent")
	}

	// Custom headers from ModelConfig
	for k, v := range model.CustomHeaders {
		req.Header.Set(k, v)
	}
}

// CallOpenAI sends a plain chat completion request to an OpenAI-compatible API.
func CallOpenAI(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, error) {
	apiKey := ResolveAPIKey(model)
	if apiKey == "" {
		return "", fmt.Errorf("no API key for provider '%s' — set %s",
			model.Provider, KeySourceLabel(model))
	}

	reqBody := buildOpenAIRequestBody(model, messages, nil, false)

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal error: %w", err)
	}

	endpoint := strings.TrimRight(model.BaseURL, "/") + "/chat/completions"

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(jsonData))
	if err != nil {
		return "", fmt.Errorf("request error: %w", err)
	}

	applyOpenAIHeaders(req, model, apiKey)

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

	reqBody := buildOpenAIRequestBody(model, messages, registry.GenerateNativeToolsSchema(), false)

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", nil, fmt.Errorf("marshal error: %w", err)
	}

	endpoint := strings.TrimRight(model.BaseURL, "/") + "/chat/completions"

	httpReq, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(jsonData))
	if err != nil {
		return "", nil, fmt.Errorf("request error: %w", err)
	}
	applyOpenAIHeaders(httpReq, model, apiKey)

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

// CallOpenAIStream streams chat completion tokens from an OpenAI-compatible API.
func CallOpenAIStream(ctx context.Context, model *ModelConfig, messages []ChatMessage) (<-chan StreamChunk, error) {
	apiKey := ResolveAPIKey(model)
	if apiKey == "" {
		return nil, fmt.Errorf("no API key for provider '%s' — set %s",
			model.Provider, KeySourceLabel(model))
	}

	reqBody := buildOpenAIRequestBody(model, messages, nil, true)

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	endpoint := strings.TrimRight(model.BaseURL, "/") + "/chat/completions"

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	applyOpenAIHeaders(req, model, apiKey)
	req.Header.Set("Accept", "text/event-stream")

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

	ch := make(chan StreamChunk, 16)

	go func() {
		defer resp.Body.Close()
		defer close(ch)

		scanner := bufio.NewScanner(resp.Body)
		buf := make([]byte, 0, 64*1024)
		scanner.Buffer(buf, 1024*1024)

		var totalPromptTokens, totalCompletionTokens int

		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				continue
			}

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
				continue
			}

			if streamResp.Error != nil {
				ch <- StreamChunk{Error: fmt.Errorf("stream error: %s", streamResp.Error.Message)}
				return
			}

			if streamResp.Usage.TotalTokens > 0 {
				totalPromptTokens = streamResp.Usage.PromptTokens
				totalCompletionTokens = streamResp.Usage.CompletionTokens
			}

			if len(streamResp.Choices) > 0 {
				choice := streamResp.Choices[0]
				chunk := StreamChunk{
					Content:   choice.Delta.Content,
					ToolCalls: choice.Delta.ToolCalls,
					Finish:    choice.FinishReason != "",
				}
				ch <- chunk
			}
		}

		if err := scanner.Err(); err != nil && err != io.EOF {
			ch <- StreamChunk{Error: fmt.Errorf("stream scan error: %w", err)}
			return
		}

		if totalPromptTokens > 0 || totalCompletionTokens > 0 {
			TrackModelUsage(model.Model, totalPromptTokens, totalCompletionTokens)
			RecordCost(model.Model, totalPromptTokens, totalCompletionTokens)
		}
	}()

	return ch, nil
}
