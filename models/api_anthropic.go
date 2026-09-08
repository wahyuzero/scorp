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
// Anthropic API (Claude) — /v1/messages format
// Modern Messages API architecture:
// - Exact tool_result / tool_use block conversion
// - Message alternation normalization
// - StreamingProvider (SSE content_block_delta)
// - CustomHeaders & ExtraBody injection
// ──────────────────────────────────────────────

type AnthropicProvider struct{}

func init() {
	p := &AnthropicProvider{}
	RegisterProvider(ProviderSpec{
		Name:             "anthropic",
		Aliases:          []string{"claude"},
		DisplayName:      "Anthropic",
		DefaultBaseURL:   "https://api.anthropic.com",
		DefaultAPIFormat: "anthropic",
		KeyEnvs:          []string{"ANTHROPIC_API_KEY"},
	}, p)
}

func (p *AnthropicProvider) Format() string {
	return "anthropic"
}

func (p *AnthropicProvider) Call(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, error) {
	return callAnthropic(ctx, model, messages)
}

func (p *AnthropicProvider) CallWithTools(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, []ToolCall, error) {
	return CallAnthropicWithTools(ctx, model, messages)
}

func (p *AnthropicProvider) CallStream(ctx context.Context, model *ModelConfig, messages []ChatMessage) (<-chan StreamChunk, error) {
	return callAnthropicStream(ctx, model, messages)
}

type anthropicTool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"input_schema"`
}

type anthropicResponse struct {
	Content []struct {
		Type  string                 `json:"type"`
		Text  string                 `json:"text,omitempty"`
		ID    string                 `json:"id,omitempty"`
		Name  string                 `json:"name,omitempty"`
		Input map[string]interface{} `json:"input,omitempty"`
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

// buildAnthropicMessages converts internal ChatMessage array to Anthropic's specification
// extracting the system prompt and converting tool results/calls to Anthropic blocks.
func buildAnthropicMessages(messages []ChatMessage) (string, []map[string]interface{}) {
	var systemPrompt string
	var apiMessages []map[string]interface{}

	for _, msg := range messages {
		switch msg.Role {
		case "system":
			if systemPrompt != "" {
				systemPrompt += "\n\n" + msg.Content
			} else {
				systemPrompt = msg.Content
			}

		case "tool":
			// Anthropic represents tool outputs as a user message with tool_result block
			toolResultBlock := map[string]interface{}{
				"type":    "tool_result",
				"content": msg.Content,
			}

			// If previous message was user with content slice, merge into it to maintain alternation
			if len(apiMessages) > 0 && apiMessages[len(apiMessages)-1]["role"] == "user" {
				prevContent, ok := apiMessages[len(apiMessages)-1]["content"].([]map[string]interface{})
				if ok {
					apiMessages[len(apiMessages)-1]["content"] = append(prevContent, toolResultBlock)
					continue
				}
			}

			apiMessages = append(apiMessages, map[string]interface{}{
				"role":    "user",
				"content": []map[string]interface{}{toolResultBlock},
			})

		case "assistant":
			var contentBlocks []map[string]interface{}
			if msg.Content != "" {
				contentBlocks = append(contentBlocks, map[string]interface{}{
					"type": "text",
					"text": msg.Content,
				})
			}
			for i, tc := range msg.ToolCalls {
				args := tc.Function.Arguments
				var parsedArgs map[string]interface{}
				if err := json.Unmarshal([]byte(args), &parsedArgs); err != nil {
					parsedArgs = make(map[string]interface{})
				}
				callID := tc.ID
				if callID == "" {
					callID = fmt.Sprintf("call_%d", i)
				}
				contentBlocks = append(contentBlocks, map[string]interface{}{
					"type":  "tool_use",
					"id":    callID,
					"name":  tc.Function.Name,
					"input": parsedArgs,
				})
			}
			if len(contentBlocks) > 0 {
				apiMessages = append(apiMessages, map[string]interface{}{
					"role":    "assistant",
					"content": contentBlocks,
				})
			}

		case "user":
			apiMessages = append(apiMessages, map[string]interface{}{
				"role":    "user",
				"content": msg.Content,
			})
		}
	}

	return systemPrompt, apiMessages
}

func buildAnthropicRequestBody(model *ModelConfig, messages []ChatMessage, tools []anthropicTool, stream bool) map[string]interface{} {
	maxTokens := model.MaxTokens
	if maxTokens == 0 {
		maxTokens = 4096
	}

	systemPrompt, apiMessages := buildAnthropicMessages(messages)

	reqBody := map[string]interface{}{
		"model":      model.Model,
		"max_tokens": maxTokens,
		"messages":   apiMessages,
	}

	if systemPrompt != "" {
		reqBody["system"] = systemPrompt
	}

	if len(tools) > 0 {
		reqBody["tools"] = tools
	}

	if stream {
		reqBody["stream"] = true
	}

	// Inject ExtraBody if configured
	for k, v := range model.ExtraBody {
		reqBody[k] = v
	}

	return reqBody
}

func applyAnthropicHeaders(req *http.Request, model *ModelConfig, apiKey string) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	for k, v := range model.CustomHeaders {
		req.Header.Set(k, v)
	}
}

// callAnthropic sends a chat completion request to an Anthropic-compatible API.
func callAnthropic(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, error) {
	apiKey := ResolveAPIKey(model)
	if apiKey == "" {
		return "", fmt.Errorf("no API key for provider '%s' — set %s",
			model.Provider, KeySourceLabel(model))
	}

	reqBody := buildAnthropicRequestBody(model, messages, nil, false)

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal error: %w", err)
	}

	endpoint := strings.TrimRight(model.BaseURL, "/") + "/v1/messages"

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(jsonData))
	if err != nil {
		return "", fmt.Errorf("request error: %w", err)
	}
	applyAnthropicHeaders(req, model, apiKey)

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

// CallAnthropicWithTools sends a request with native tool definitions to Anthropic Messages API.
func CallAnthropicWithTools(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, []ToolCall, error) {
	apiKey := ResolveAPIKey(model)
	if apiKey == "" {
		return "", nil, fmt.Errorf("no API key for provider '%s' — set %s", model.Provider, KeySourceLabel(model))
	}

	// Generate and sanitize native tools schema
	nativeTools := registry.GenerateNativeToolsSchema()
	if model.ToolSchemaTransform != "" {
		nativeTools = TransformToolDefinitions(nativeTools, model.ToolSchemaTransform)
	}

	var anthropicTools []anthropicTool
	for _, td := range nativeTools {
		anthropicTools = append(anthropicTools, anthropicTool{
			Name:        td.Function.Name,
			Description: td.Function.Description,
			InputSchema: td.Function.Parameters,
		})
	}

	reqBody := buildAnthropicRequestBody(model, messages, anthropicTools, false)

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", nil, fmt.Errorf("marshal error: %w", err)
	}

	endpoint := strings.TrimRight(model.BaseURL, "/") + "/v1/messages"

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(jsonData))
	if err != nil {
		return "", nil, fmt.Errorf("request error: %w", err)
	}
	applyAnthropicHeaders(req, model, apiKey)

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

	var textParts []string
	var toolCalls []ToolCall

	for _, block := range apiResp.Content {
		switch block.Type {
		case "text":
			textParts = append(textParts, block.Text)
		case "tool_use":
			toolCalls = append(toolCalls, ToolCall{
				Name: block.Name,
				Args: block.Input,
			})
		}
	}

	TrackModelUsage(model.Model, apiResp.Usage.InputTokens, apiResp.Usage.OutputTokens)
	RecordCost(model.Model, apiResp.Usage.InputTokens, apiResp.Usage.OutputTokens)

	return strings.Join(textParts, "\n"), toolCalls, nil
}

// callAnthropicStream streams chat completion tokens from Anthropic SSE endpoint.
func callAnthropicStream(ctx context.Context, model *ModelConfig, messages []ChatMessage) (<-chan StreamChunk, error) {
	apiKey := ResolveAPIKey(model)
	if apiKey == "" {
		return nil, fmt.Errorf("no API key for provider '%s' — set %s",
			model.Provider, KeySourceLabel(model))
	}

	reqBody := buildAnthropicRequestBody(model, messages, nil, true)

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	endpoint := strings.TrimRight(model.BaseURL, "/") + "/v1/messages"

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	applyAnthropicHeaders(req, model, apiKey)
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

		var inputTokens, outputTokens int

		for scanner.Scan() {
			line := scanner.Text()
			if line == "" || !strings.HasPrefix(line, "data: ") {
				continue
			}

			data := strings.TrimPrefix(line, "data: ")

			var event struct {
				Type  string `json:"type"`
				Delta *struct {
					Type string `json:"type"`
					Text string `json:"text"`
				} `json:"delta,omitempty"`
				Usage *struct {
					OutputTokens int `json:"output_tokens"`
				} `json:"usage,omitempty"`
				Message *struct {
					Usage struct {
						InputTokens int `json:"input_tokens"`
					} `json:"usage"`
				} `json:"message,omitempty"`
			}

			if err := json.Unmarshal([]byte(data), &event); err != nil {
				continue
			}

			if event.Message != nil && event.Message.Usage.InputTokens > 0 {
				inputTokens = event.Message.Usage.InputTokens
			}
			if event.Usage != nil && event.Usage.OutputTokens > 0 {
				outputTokens = event.Usage.OutputTokens
			}

			switch event.Type {
			case "content_block_delta":
				if event.Delta != nil && event.Delta.Text != "" {
					ch <- StreamChunk{Content: event.Delta.Text}
				}
			case "message_stop":
				ch <- StreamChunk{Finish: true}
				if inputTokens > 0 || outputTokens > 0 {
					TrackModelUsage(model.Model, inputTokens, outputTokens)
					RecordCost(model.Model, inputTokens, outputTokens)
				}
				return
			}
		}

		if err := scanner.Err(); err != nil && err != io.EOF {
			log.Printf("[anthropic/stream] scan error: %v", err)
			ch <- StreamChunk{Error: err}
		}
	}()

	return ch, nil
}
