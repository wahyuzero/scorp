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
	"strings"

	"scorp-agent/internal/helpers"
	"scorp-agent/registry"
)

// ──────────────────────────────────────────────
// Gemini API (Google) — :generateContent format
// Complete Google AI Studio & Gemini 2.5 / 3.x Features:
// - Direct API key authentication
// - Native function calling with sanitized tool schemas
// - Multi-turn tool execution loop (functionResponse matching)
// - Thinking & reasoning separation (thought parts filtering)
// - ThinkingConfig support (thinkingLevel: low/medium/high, thinkingBudget)
// - ExtraBody injection (safetySettings, searchGrounding)
// - SSE token streaming via :streamGenerateContent?alt=sse
// ──────────────────────────────────────────────

type GeminiProvider struct{}

func init() {
	p := &GeminiProvider{}
	RegisterProvider(ProviderSpec{
		Name:             "gemini",
		Aliases:          []string{"google"},
		DisplayName:      "Google Gemini",
		DefaultBaseURL:   "https://generativelanguage.googleapis.com/v1beta",
		DefaultAPIFormat: "gemini",
		KeyEnvs:          []string{"GOOGLE_API_KEY", "GEMINI_API_KEY"},
	}, p)

	RegisterCatalog("gemini", []CatalogEntry{
		{"gemini-3.8-flash", 65536, false, "gemini-3.8-flash"},
		{"gemini-3.7-flash", 65536, false, "gemini-3.7-flash"},
		{"gemini-3.1-pro-preview", 65536, true, "gemini-3.1-pro"},
		{"gemini-3.1-flash-lite", 65536, false, "gemini-3.1-flash-lite"},
		{"gemini-flash-latest", 65536, false, "gemini-flash-latest"},
		{"gemini-pro-latest", 65536, true, "gemini-pro-latest"},
	})
}

func (p *GeminiProvider) Format() string {
	return "gemini"
}

func (p *GeminiProvider) Call(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, error) {
	return callGemini(ctx, model, messages)
}

func (p *GeminiProvider) CallWithTools(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, []ToolCall, error) {
	return CallGeminiWithTools(ctx, model, messages)
}

func (p *GeminiProvider) CallStream(ctx context.Context, model *ModelConfig, messages []ChatMessage) (<-chan StreamChunk, error) {
	return callGeminiStream(ctx, model, messages)
}

// geminiRequest is the request body for the Gemini generateContent API.
type geminiRequest struct {
	Contents          []geminiContent  `json:"contents"`
	SystemInstruction *geminiContent   `json:"systemInstruction,omitempty"`
	Tools             []geminiToolSet  `json:"tools,omitempty"`
	GenerationConfig  *geminiGenConfig `json:"generationConfig,omitempty"`
}

type geminiContent struct {
	Role  string       `json:"role"`
	Parts []geminiPart `json:"parts"`
}

type geminiInlineData struct {
	MimeType string `json:"mimeType"`
	Data     string `json:"data"`
}

type geminiPart struct {
	Text             string                  `json:"text,omitempty"`
	InlineData       *geminiInlineData       `json:"inlineData,omitempty"`
	Thought          bool                    `json:"thought,omitempty"`
	ThoughtSignature string                  `json:"thoughtSignature,omitempty"`
	FunctionCall     *geminiFunctionCall     `json:"functionCall,omitempty"`
	FunctionResponse *geminiFunctionResponse `json:"functionResponse,omitempty"`
}

func parseDataURL(dataURL string) *geminiInlineData {
	dataURL = strings.TrimSpace(dataURL)
	if strings.HasPrefix(dataURL, "data:") {
		idx := strings.Index(dataURL, ";base64,")
		if idx > 5 {
			mimeType := dataURL[5:idx]
			b64Data := dataURL[idx+8:]
			return &geminiInlineData{
				MimeType: mimeType,
				Data:     b64Data,
			}
		}
	}
	return nil
}

type geminiFunctionCall struct {
	Name string                 `json:"name"`
	Args map[string]interface{} `json:"args"`
}

type geminiFunctionResponse struct {
	Name     string                 `json:"name"`
	Response map[string]interface{} `json:"response"`
}

type geminiToolSet struct {
	FunctionDeclarations []geminiFuncDecl `json:"functionDeclarations"`
}

type geminiFuncDecl struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
}

type geminiGenConfig struct {
	MaxOutputTokens int                   `json:"maxOutputTokens,omitempty"`
	Temperature     float64               `json:"temperature,omitempty"`
	ThinkingConfig  *geminiThinkingConfig `json:"thinkingConfig,omitempty"`
}

type geminiThinkingConfig struct {
	ThinkingLevel  string `json:"thinkingLevel,omitempty"`  // "low" | "medium" | "high"
	ThinkingBudget int    `json:"thinkingBudget,omitempty"` // 0 disables thinking
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Role  string       `json:"role"`
			Parts []geminiPart `json:"parts"`
		} `json:"content"`
		FinishReason string `json:"finishReason"`
	} `json:"candidates"`
	UsageMetadata struct {
		PromptTokenCount     int `json:"promptTokenCount"`
		CandidatesTokenCount int `json:"candidatesTokenCount"`
		TotalTokenCount      int `json:"totalTokenCount"`
	} `json:"usageMetadata"`
	Error *struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
		Status  string `json:"status"`
	} `json:"error,omitempty"`
}

// geminiMessages converts chatMessage slice to Gemini's content format.
// Handles system instruction extraction, assistant functionCalls, and tool responses.
func geminiMessages(messages []ChatMessage) ([]geminiContent, *geminiContent) {
	var contents []geminiContent
	var systemParts []geminiPart
	lastToolName := "tool"

	for _, m := range messages {
		switch m.Role {
		case "system":
			systemParts = append(systemParts, geminiPart{Text: m.Content})

		case "assistant":
			var parts []geminiPart
			if m.Content != "" {
				parts = append(parts, geminiPart{Text: m.Content})
			}
			for _, tc := range m.ToolCalls {
				var parsedArgs map[string]interface{}
				if err := json.Unmarshal([]byte(tc.Function.Arguments), &parsedArgs); err != nil {
					parsedArgs = make(map[string]interface{})
				}
				parts = append(parts, geminiPart{
					FunctionCall: &geminiFunctionCall{
						Name: tc.Function.Name,
						Args: parsedArgs,
					},
					ThoughtSignature: tc.ThoughtSignature,
				})
				lastToolName = tc.Function.Name
			}
			if len(parts) > 0 {
				contents = append(contents, geminiContent{
					Role:  "model",
					Parts: parts,
				})
			}

		case "tool":
			// Tool output returned as functionResponse
			contents = append(contents, geminiContent{
				Role: "user",
				Parts: []geminiPart{
					{
						FunctionResponse: &geminiFunctionResponse{
							Name: lastToolName,
							Response: map[string]interface{}{
								"result": m.Content,
							},
						},
					},
				},
			})

		default: // "user" or anything else
			trimmed := strings.TrimSpace(m.Content)

			// Native Gemini functionResponse detection from tool result strings
			if strings.HasPrefix(trimmed, "[Tool Result: ") {
				toolName := "tool"
				if endIdx := strings.Index(trimmed, "]"); endIdx > 14 {
					toolName = strings.TrimSpace(trimmed[14:endIdx])
				}
				resBody := strings.TrimSpace(trimmed[strings.Index(trimmed, "\n")+1:])
				contents = append(contents, geminiContent{
					Role: "user",
					Parts: []geminiPart{
						{
							FunctionResponse: &geminiFunctionResponse{
								Name: toolName,
								Response: map[string]interface{}{
									"result": resBody,
								},
							},
						},
					},
				})
				continue
			}

			if strings.HasPrefix(trimmed, "[") && (strings.Contains(trimmed, "image_url") || strings.Contains(trimmed, "inlineData") || strings.Contains(trimmed, "image")) {
				var rawParts []map[string]interface{}
				if err := json.Unmarshal([]byte(trimmed), &rawParts); err == nil {
					var userParts []geminiPart
					for _, rp := range rawParts {
						pType, _ := rp["type"].(string)
						switch pType {
						case "text":
							if txt, ok := rp["text"].(string); ok && txt != "" {
								userParts = append(userParts, geminiPart{Text: txt})
							}
						case "image_url":
							var dataURL string
							if imgURLMap, ok := rp["image_url"].(map[string]interface{}); ok {
								dataURL, _ = imgURLMap["url"].(string)
							} else if strURL, ok := rp["image_url"].(string); ok {
								dataURL = strURL
							}
							if dataURL != "" {
								inline := parseDataURL(dataURL)
								if inline != nil {
									userParts = append(userParts, geminiPart{InlineData: inline})
								}
							}
						}
					}
					if len(userParts) > 0 {
						contents = append(contents, geminiContent{
							Role:  "user",
							Parts: userParts,
						})
						continue
					}
				}
			}

			// Plain text content
			contents = append(contents, geminiContent{
				Role:  "user",
				Parts: []geminiPart{{Text: m.Content}},
			})
		}
	}

	// Gemini API requires alternating turns and strictly forbids ending with a model turn
	if len(contents) > 0 && contents[len(contents)-1].Role == "model" {
		contents = append(contents, geminiContent{
			Role:  "user",
			Parts: []geminiPart{{Text: "Continue."}},
		})
	}

	var sys *geminiContent
	if len(systemParts) > 0 {
		sys = &geminiContent{Parts: systemParts}
	}

	return contents, sys
}

// geminiBuildRequest constructs the geminiRequest with defaults and thinking config.
func geminiBuildRequest(model *ModelConfig, messages []ChatMessage, withTools bool) geminiRequest {
	maxTokens := model.MaxTokens
	if maxTokens == 0 {
		maxTokens = 4096
	}

	contents, sys := geminiMessages(messages)

	genConfig := &geminiGenConfig{
		MaxOutputTokens: maxTokens,
		Temperature:     0.7,
	}

	// Thinking configuration (Gemini 2.5 / 3.x)
	if model.ExtraBody != nil {
		if tl, ok := model.ExtraBody["thinking_level"].(string); ok && tl != "" {
			genConfig.ThinkingConfig = &geminiThinkingConfig{ThinkingLevel: tl}
		}
	} else if strings.Contains(strings.ToLower(model.Model), "flash-lite") {
		// Only flash-lite models support "minimal" thinking level
		genConfig.ThinkingConfig = &geminiThinkingConfig{ThinkingLevel: "minimal"}
	}

	req := geminiRequest{
		Contents:         contents,
		GenerationConfig: genConfig,
	}

	if sys != nil {
		req.SystemInstruction = sys
	}

	if withTools {
		funcDecls := make([]geminiFuncDecl, 0, 32)
		for _, td := range registry.GenerateNativeToolsSchema() {
			funcDecls = append(funcDecls, geminiFuncDecl{
				Name:        td.Function.Name,
				Description: td.Function.Description,
				Parameters:  SanitizeToolSchema(td.Function.Parameters, "simple"),
			})
		}
		req.Tools = []geminiToolSet{{FunctionDeclarations: funcDecls}}
	}

	return req
}

func resolveGeminiBaseURL(model *ModelConfig) string {
	base := ResolveBaseURL(model)
	if base == "" {
		base = "https://generativelanguage.googleapis.com"
	}
	base = strings.TrimRight(base, "/")
	base = strings.TrimSuffix(base, "/v1beta")
	return base
}

// geminiDoRequest sends the request and returns the parsed response.
func geminiDoRequest(ctx context.Context, model *ModelConfig, reqBody geminiRequest) (*geminiResponse, error) {
	apiKey := ResolveAPIKey(model)
	if apiKey == "" {
		return nil, fmt.Errorf("no API key for provider '%s' — set %s",
			model.Provider, KeySourceLabel(model))
	}

	// Marshal and merge ExtraBody if configured
	reqMap := make(map[string]interface{})
	reqData, _ := json.Marshal(reqBody)
	_ = json.Unmarshal(reqData, &reqMap)

	for k, v := range model.ExtraBody {
		if k != "thinking_level" { // handled in GenerationConfig
			reqMap[k] = v
		}
	}

	jsonData, err := json.Marshal(reqMap)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	base := resolveGeminiBaseURL(model)
	endpoint := base + "/v1beta/models/" + model.Model + ":generateContent"

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-goog-api-key", apiKey)

	for k, v := range model.CustomHeaders {
		req.Header.Set(k, v)
	}

	client := GetAIClient(base)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API call failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read error: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API error (HTTP %d): %s", resp.StatusCode, helpers.TruncateStr(string(body), 300))
	}

	var apiResp geminiResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("parse error: %s", helpers.TruncateStr(string(body), 200))
	}

	if apiResp.Error != nil {
		return nil, fmt.Errorf("API error: %s", apiResp.Error.Message)
	}

	return &apiResp, nil
}

// callGemini sends a chat completion request to a Gemini-compatible API.
func callGemini(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, error) {
	reqBody := geminiBuildRequest(model, messages, false)
	apiResp, err := geminiDoRequest(ctx, model, reqBody)
	if err != nil {
		return "", err
	}

	if len(apiResp.Candidates) == 0 {
		return "", fmt.Errorf("no response candidates")
	}

	// Concatenate text parts (skipping internal thought tokens)
	var sb strings.Builder
	for _, part := range apiResp.Candidates[0].Content.Parts {
		if part.Thought {
			continue
		}
		if part.Text != "" {
			sb.WriteString(part.Text)
		}
	}

	TrackModelUsage(model.Model, apiResp.UsageMetadata.PromptTokenCount, apiResp.UsageMetadata.CandidatesTokenCount)
	RecordCost(model.Model, apiResp.UsageMetadata.PromptTokenCount, apiResp.UsageMetadata.CandidatesTokenCount)

	return sb.String(), nil
}

// CallGeminiWithTools sends a request with native tool definitions.
func CallGeminiWithTools(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, []ToolCall, error) {
	reqBody := geminiBuildRequest(model, messages, true)
	if os.Getenv("SCORP_DEBUG") != "" {
		dbgJSON, _ := json.Marshal(reqBody)
		log.Printf("[gemini/debug] Request payload: %s", string(dbgJSON))
	}
	apiResp, err := geminiDoRequest(ctx, model, reqBody)
	if err != nil {
		return "", nil, err
	}

	if len(apiResp.Candidates) == 0 {
		return "", nil, fmt.Errorf("no response candidates")
	}

	var sb strings.Builder
	var toolCalls []ToolCall

	for _, part := range apiResp.Candidates[0].Content.Parts {
		if part.Thought {
			continue
		}
		if part.Text != "" {
			sb.WriteString(part.Text)
		}
		if part.FunctionCall != nil {
			args := part.FunctionCall.Args
			if args == nil {
				args = make(map[string]interface{})
			}
			toolCalls = append(toolCalls, ToolCall{
				Name:             part.FunctionCall.Name,
				Args:             args,
				ThoughtSignature: part.ThoughtSignature,
			})
			log.Printf("[agent] Gemini functionCall: %s(%v)", part.FunctionCall.Name, args)
		}
	}

	TrackModelUsage(model.Model, apiResp.UsageMetadata.PromptTokenCount, apiResp.UsageMetadata.CandidatesTokenCount)
	RecordCost(model.Model, apiResp.UsageMetadata.PromptTokenCount, apiResp.UsageMetadata.CandidatesTokenCount)

	return sb.String(), toolCalls, nil
}

// callGeminiStream streams tokens from Gemini SSE endpoint
func callGeminiStream(ctx context.Context, model *ModelConfig, messages []ChatMessage) (<-chan StreamChunk, error) {
	apiKey := ResolveAPIKey(model)
	if apiKey == "" {
		return nil, fmt.Errorf("no API key for provider '%s' — set %s",
			model.Provider, KeySourceLabel(model))
	}

	reqBody := geminiBuildRequest(model, messages, false)

	reqMap := make(map[string]interface{})
	reqData, _ := json.Marshal(reqBody)
	_ = json.Unmarshal(reqData, &reqMap)

	for k, v := range model.ExtraBody {
		if k != "thinking_level" {
			reqMap[k] = v
		}
	}

	jsonData, err := json.Marshal(reqMap)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	base := resolveGeminiBaseURL(model)
	endpoint := base + "/v1beta/models/" + model.Model + ":streamGenerateContent?alt=sse"

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-goog-api-key", apiKey)
	req.Header.Set("Accept", "text/event-stream")

	for k, v := range model.CustomHeaders {
		req.Header.Set(k, v)
	}

	client := GetAIClient(base)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API call failed: %w", err)
	}

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("Gemini API error (HTTP %d): %s", resp.StatusCode, helpers.TruncateStr(string(body), 300))
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
			var streamResp geminiResponse
			if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
				continue
			}

			if streamResp.Error != nil {
				ch <- StreamChunk{Error: fmt.Errorf("stream error: %s", streamResp.Error.Message)}
				return
			}

			for _, candidate := range streamResp.Candidates {
				for _, part := range candidate.Content.Parts {
					if part.Thought {
						continue // skip internal thought tokens from user text stream
					}
					if part.Text != "" {
						ch <- StreamChunk{Content: part.Text}
					}
				}
				if candidate.FinishReason != "" {
					ch <- StreamChunk{Finish: true}
				}
			}
		}
	}()

	return ch, nil
}
