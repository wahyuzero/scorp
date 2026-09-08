package models

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"scorp-agent/internal/helpers"
	"scorp-agent/registry"
)

// ──────────────────────────────────────────────
// Azure OpenAI Provider
// Supports deployment-based endpoints and api-key headers
// ──────────────────────────────────────────────

type AzureProvider struct{}

func init() {
	p := &AzureProvider{}
	RegisterProvider(ProviderSpec{
		Name:             "azure",
		Aliases:          []string{"azure-openai"},
		DisplayName:      "Azure OpenAI",
		DefaultBaseURL:   "",
		DefaultAPIFormat: "azure",
		KeyEnvs:          []string{"AZURE_OPENAI_API_KEY", "AZURE_API_KEY"},
	}, p)

	RegisterCatalog("azure", []CatalogEntry{
		{"gpt-4o", 16384, true, "azure-gpt4o"},
		{"gpt-4o-mini", 16384, false, "azure-gpt4o-mini"},
		{"gpt-4-turbo", 4096, true, "azure-gpt4-turbo"},
		{"o1", 32768, true, "azure-o1"},
		{"o3-mini", 16384, true, "azure-o3-mini"},
	})
}

func (p *AzureProvider) Format() string {
	return "azure"
}

func (p *AzureProvider) Call(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, error) {
	return callAzure(ctx, model, messages)
}

func (p *AzureProvider) CallWithTools(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, []ToolCall, error) {
	return callAzureWithTools(ctx, model, messages)
}

func (p *AzureProvider) CallStream(ctx context.Context, model *ModelConfig, messages []ChatMessage) (<-chan StreamChunk, error) {
	return callAzureStream(ctx, model, messages)
}

func resolveAzureEndpoint(model *ModelConfig) string {
	base := strings.TrimRight(model.BaseURL, "/")
	if base == "" {
		return ""
	}

	// If already a full chat/completions endpoint
	if strings.Contains(base, "/chat/completions") {
		return base
	}

	// If deployment URL: https://{resource}.openai.azure.com/openai/deployments/{deployment}
	if strings.Contains(base, "openai.azure.com") {
		if !strings.Contains(base, "/openai/deployments/") {
			base = fmt.Sprintf("%s/openai/deployments/%s", base, model.Model)
		}
		return base + "/chat/completions?api-version=2024-06-01"
	}

	return base + "/chat/completions"
}

func applyAzureHeaders(req *http.Request, model *ModelConfig, apiKey string) {
	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("api-key", apiKey)
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}
	for k, v := range model.CustomHeaders {
		req.Header.Set(k, v)
	}
}

func callAzure(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, error) {
	apiKey := ResolveAPIKey(model)
	if apiKey == "" {
		return "", fmt.Errorf("no API key for Azure OpenAI — set %s", KeySourceLabel(model))
	}

	endpoint := resolveAzureEndpoint(model)
	if endpoint == "" {
		return "", fmt.Errorf("base_url required for Azure OpenAI (e.g. https://<resource>.openai.azure.com)")
	}

	reqBody := buildOpenAIRequestBody(model, messages, nil, false)
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal error: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(jsonData))
	if err != nil {
		return "", fmt.Errorf("request error: %w", err)
	}
	applyAzureHeaders(req, model, apiKey)

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
		return "", fmt.Errorf("Azure API error (HTTP %d): %s", resp.StatusCode, helpers.TruncateStr(string(body), 300))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("parse error: %s", helpers.TruncateStr(string(body), 200))
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no response choices from Azure")
	}

	reply := chatResp.Choices[0].Message.Content
	TrackModelUsage(model.Model, chatResp.Usage.PromptTokens, chatResp.Usage.CompletionTokens)
	RecordCost(model.Model, chatResp.Usage.PromptTokens, chatResp.Usage.CompletionTokens)

	return reply, nil
}

func callAzureWithTools(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, []ToolCall, error) {
	apiKey := ResolveAPIKey(model)
	if apiKey == "" {
		return "", nil, fmt.Errorf("no API key for Azure OpenAI — set %s", KeySourceLabel(model))
	}

	endpoint := resolveAzureEndpoint(model)
	if endpoint == "" {
		return "", nil, fmt.Errorf("base_url required for Azure OpenAI")
	}

	reqBody := buildOpenAIRequestBody(model, messages, registry.GenerateNativeToolsSchema(), false)
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", nil, fmt.Errorf("marshal error: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(jsonData))
	if err != nil {
		return "", nil, fmt.Errorf("request error: %w", err)
	}
	applyAzureHeaders(req, model, apiKey)

	client := GetAIClient(model.BaseURL)
	resp, err := client.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("API call failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, fmt.Errorf("read error: %w", err)
	}

	if resp.StatusCode != 200 {
		return "", nil, fmt.Errorf("Azure API error (HTTP %d): %s", resp.StatusCode, helpers.TruncateStr(string(body), 300))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", nil, fmt.Errorf("parse error: %s", helpers.TruncateStr(string(body), 200))
	}

	if len(chatResp.Choices) == 0 {
		return "", nil, fmt.Errorf("no response choices from Azure")
	}

	choice := chatResp.Choices[0]
	var toolCalls []ToolCall
	for _, tc := range choice.Message.ToolCalls {
		var args map[string]interface{}
		if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err != nil {
			args = make(map[string]interface{})
		}
		toolCalls = append(toolCalls, ToolCall{
			Name: tc.Function.Name,
			Args: args,
		})
	}

	TrackModelUsage(model.Model, chatResp.Usage.PromptTokens, chatResp.Usage.CompletionTokens)
	RecordCost(model.Model, chatResp.Usage.PromptTokens, chatResp.Usage.CompletionTokens)

	return choice.Message.Content, toolCalls, nil
}

func callAzureStream(ctx context.Context, model *ModelConfig, messages []ChatMessage) (<-chan StreamChunk, error) {
	apiKey := ResolveAPIKey(model)
	if apiKey == "" {
		return nil, fmt.Errorf("no API key for Azure OpenAI — set %s", KeySourceLabel(model))
	}

	endpoint := resolveAzureEndpoint(model)
	if endpoint == "" {
		return nil, fmt.Errorf("base_url required for Azure OpenAI")
	}

	reqBody := buildOpenAIRequestBody(model, messages, nil, true)
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	applyAzureHeaders(req, model, apiKey)
	req.Header.Set("Accept", "text/event-stream")

	client := GetAIClient(model.BaseURL)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API call failed: %w", err)
	}

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("Azure API error (HTTP %d): %s", resp.StatusCode, helpers.TruncateStr(string(body), 300))
	}

	ch := make(chan StreamChunk, 16)
	go func() {
		defer resp.Body.Close()
		defer close(ch)

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" || !strings.HasPrefix(line, "data: ") {
				continue
			}
			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				ch <- StreamChunk{Finish: true}
				return
			}
			var streamResp struct {
				Choices []struct {
					Delta struct {
						Content   string         `json:"content"`
						ToolCalls []ToolCallResp `json:"tool_calls"`
					} `json:"delta"`
					FinishReason string `json:"finish_reason"`
				} `json:"choices"`
			}
			if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
				continue
			}
			if len(streamResp.Choices) > 0 {
				c := streamResp.Choices[0]
				ch <- StreamChunk{
					Content:   c.Delta.Content,
					ToolCalls: c.Delta.ToolCalls,
					Finish:    c.FinishReason != "",
				}
			}
		}
	}()

	return ch, nil
}
