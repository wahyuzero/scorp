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
	"os"
	"os/exec"
	"strings"
	"time"

	"scorp-agent/config"
	"scorp-agent/internal/helpers"
	"scorp-agent/registry"
)

// ──────────────────────────────────────────────
// OpenCode Zen & Free Models Dedicated Provider
//
// Endpoint: https://opencode.ai/zen/v1 (or alternate https://opencode.ai/zen/go/v1)
// Free Models Catalog:
//   - big-pickle               (Fastest, 200K ctx, 32K out, best for chat & fallback)
//   - mimo-v2.5-free           (Xiaomi MiMo 2.5, 200K ctx, 32K out, best for code & tools)
//   - ling-3.0-flash-fin-free  (256K ctx, 32K out, fast structured analysis)
//   - laguna-s-2.1-free        (Poolside Laguna 2.1, 256K ctx, 32K out, deep reasoning)
//
// Required Gateway Headers:
//   - Authorization: Bearer <key>
//   - User-Agent: opencode/1.0.0
//   - x-opencode-session: <sessionId> (Strictly required by OpenCode Zen routing)
// ──────────────────────────────────────────────

const (
	OpenCodeZenBaseURL      = "https://opencode.ai/zen/v1"
	OpenCodeZenAltURL       = "https://opencode.ai/zen/go/v1"
	OpenCodeDefaultUserAgent = "opencode/1.0.0"
)

// OpenCodeProvider implements LLMProvider for OpenCode Zen
type OpenCodeProvider struct{}

func init() {
	p := &OpenCodeProvider{}
	RegisterProvider(ProviderSpec{
		Name:             "opencode",
		Aliases:          []string{"opencode-zen", "opencode-free"},
		DisplayName:      "OpenCode Zen (Free AI Gateway)",
		DefaultBaseURL:   OpenCodeZenBaseURL,
		DefaultAPIFormat: "opencode",
		KeyEnvs:          []string{"OPENCODE_API_KEY", "OPENCODE_ZEN_API_KEY"},
	}, p)
}

func (p *OpenCodeProvider) Format() string {
	return "opencode"
}

func (p *OpenCodeProvider) Call(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, error) {
	return CallOpenCode(ctx, model, messages)
}

func (p *OpenCodeProvider) CallWithTools(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, []ToolCall, error) {
	return CallOpenCodeWithTools(ctx, model, messages)
}

func (p *OpenCodeProvider) CallStream(ctx context.Context, model *ModelConfig, messages []ChatMessage) (<-chan StreamChunk, error) {
	return CallOpenCodeStream(ctx, model, messages)
}

func (p *OpenCodeProvider) ResolveKey(cfg *ModelConfig) string {
	return resolveOpenCodeKeyFromDisk()
}

// resolveOpenCodeKeyFromDisk checks local opencode SQLite database for OpenCode Zen key
func resolveOpenCodeKeyFromDisk() string {
	home := config.HomeDir()
	dbPath := home + "/.local/share/opencode/opencode.db"
	if _, err := os.Stat(dbPath); err != nil {
		return ""
	}
	out, err := exec.Command("sqlite3", dbPath, "SELECT value FROM credential WHERE integration_id = 'opencode' LIMIT 1;").Output()
	if err == nil && len(out) > 0 {
		var cred struct {
			Key string `json:"key"`
		}
		if err := json.Unmarshal(out, &cred); err == nil && cred.Key != "" {
			return cred.Key
		}
	}
	return ""
}

// resolveOpenCodeSessionID generates a unique session trace for OpenCode Zen gateway routing
func resolveOpenCodeSessionID() string {
	return fmt.Sprintf("sess_scorp_%d", time.Now().UnixNano())
}

// resolveOpenCodeBaseURL ensures valid base URL defaulting to official Zen gateway
func resolveOpenCodeBaseURL(model *ModelConfig) string {
	if model != nil && model.BaseURL != "" {
		return strings.TrimRight(model.BaseURL, "/")
	}
	return OpenCodeZenBaseURL
}

// CallOpenCode sends a standard text chat completion request to OpenCode Zen
func CallOpenCode(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, error) {
	apiKey := ResolveAPIKey(model)
	if apiKey == "" {
		return "", fmt.Errorf("no API key for OpenCode Zen — set OPENCODE_API_KEY or login with opencode CLI")
	}

	maxTokens := model.MaxTokens
	if maxTokens <= 0 {
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

	endpoint := resolveOpenCodeBaseURL(model) + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(jsonData))
	if err != nil {
		return "", fmt.Errorf("request error: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("User-Agent", OpenCodeDefaultUserAgent)
	req.Header.Set("x-opencode-session", resolveOpenCodeSessionID())

	client := GetAIClient(resolveOpenCodeBaseURL(model))
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("API call failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read error: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
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

	// Track usage and Cloudflare Workers AI prompt caching
	cachedTokens := 0
	if chatResp.Usage.PromptTokensDetails != nil {
		cachedTokens = chatResp.Usage.PromptTokensDetails.CachedTokens
	}
	TrackModelUsageWithCache(model.Model, chatResp.Usage.PromptTokens, chatResp.Usage.CompletionTokens, cachedTokens)
	RecordCostWithCache(model.Model, chatResp.Usage.PromptTokens, chatResp.Usage.CompletionTokens, cachedTokens)

	return reply, nil
}

// CallOpenCodeWithTools sends a completion request with native tool schemas to OpenCode Zen
func CallOpenCodeWithTools(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, []ToolCall, error) {
	apiKey := ResolveAPIKey(model)
	if apiKey == "" {
		return "", nil, fmt.Errorf("no API key for OpenCode Zen — set OPENCODE_API_KEY or login with opencode CLI")
	}

	maxTokens := model.MaxTokens
	if maxTokens <= 0 {
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

	endpoint := resolveOpenCodeBaseURL(model) + "/chat/completions"
	httpReq, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(jsonData))
	if err != nil {
		return "", nil, fmt.Errorf("request error: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("User-Agent", OpenCodeDefaultUserAgent)
	httpReq.Header.Set("x-opencode-session", resolveOpenCodeSessionID())

	client := GetAIClient(resolveOpenCodeBaseURL(model))
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", nil, fmt.Errorf("API call failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, fmt.Errorf("read error: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
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

	// Parse native tool calls returned by OpenCode Zen
	var toolCalls []ToolCall
	for _, tc := range choice.Message.ToolCalls {
		var args map[string]interface{}
		if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err != nil {
			log.Printf("[opencode] Failed to parse tool args '%s': %v", tc.Function.Arguments, err)
			args = make(map[string]interface{})
		}
		toolCalls = append(toolCalls, ToolCall{
			Name: tc.Function.Name,
			Args: args,
		})
	}

	// Track usage and Cloudflare Workers AI prompt caching
	cachedTokens := 0
	if chatResp.Usage.PromptTokensDetails != nil {
		cachedTokens = chatResp.Usage.PromptTokensDetails.CachedTokens
	}
	TrackModelUsageWithCache(model.Model, chatResp.Usage.PromptTokens, chatResp.Usage.CompletionTokens, cachedTokens)
	RecordCostWithCache(model.Model, chatResp.Usage.PromptTokens, chatResp.Usage.CompletionTokens, cachedTokens)

	return content, toolCalls, nil
}

// CallOpenCodeStream streams chat completion chunks from OpenCode Zen via SSE
func CallOpenCodeStream(ctx context.Context, model *ModelConfig, messages []ChatMessage) (<-chan StreamChunk, error) {
	apiKey := ResolveAPIKey(model)
	if apiKey == "" {
		return nil, fmt.Errorf("no API key for OpenCode Zen — set OPENCODE_API_KEY")
	}

	maxTokens := model.MaxTokens
	if maxTokens <= 0 {
		maxTokens = 4096
	}

	reqBody := ChatRequest{
		Model:       model.Model,
		Messages:    messages,
		MaxTokens:   maxTokens,
		Temperature: 0.7,
		Stream:      true,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	endpoint := resolveOpenCodeBaseURL(model) + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("User-Agent", OpenCodeDefaultUserAgent)
	req.Header.Set("x-opencode-session", resolveOpenCodeSessionID())

	client := GetAIClient(resolveOpenCodeBaseURL(model))
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API call failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("API error (HTTP %d): %s", resp.StatusCode, helpers.TruncateStr(string(body), 300))
	}

	ch := make(chan StreamChunk, 16)
	go func() {
		defer resp.Body.Close()
		defer close(ch)

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if !strings.HasPrefix(line, "data: ") {
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
						Content string `json:"content"`
					} `json:"delta"`
					FinishReason *string `json:"finish_reason"`
				} `json:"choices"`
			}

			if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
				continue
			}

			if len(streamResp.Choices) > 0 {
				choice := streamResp.Choices[0]
				if choice.Delta.Content != "" {
					ch <- StreamChunk{Content: choice.Delta.Content}
				}
				if choice.FinishReason != nil && *choice.FinishReason != "" {
					ch <- StreamChunk{Finish: true}
					return
				}
			}
		}

		if err := scanner.Err(); err != nil {
			ch <- StreamChunk{Error: err}
		}
	}()

	return ch, nil
}
