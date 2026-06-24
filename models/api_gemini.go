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
// Gemini API (Google) — :generateContent format
// ──────────────────────────────────────────────

// geminiRequest is the request body for the Gemini generateContent API.
type geminiRequest struct {
	Contents          []geminiContent      `json:"contents"`
	SystemInstruction *geminiContent       `json:"systemInstruction,omitempty"`
	Tools             []geminiToolSet      `json:"tools,omitempty"`
	GenerationConfig  *geminiGenConfig     `json:"generationConfig,omitempty"`
}

type geminiContent struct {
	Role  string       `json:"role"`
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text         string                 `json:"text,omitempty"`
	FunctionCall *geminiFunctionCall    `json:"functionCall,omitempty"`
}

type geminiFunctionCall struct {
	Name string                 `json:"name"`
	Args map[string]interface{} `json:"args"`
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
	MaxOutputTokens int     `json:"maxOutputTokens,omitempty"`
	Temperature     float64 `json:"temperature,omitempty"`
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
// Gemini uses "model" as the assistant role name.
func geminiMessages(messages []ChatMessage) ([]geminiContent, *geminiContent) {
	var contents []geminiContent
	var systemParts []geminiPart

	for _, m := range messages {
		switch m.Role {
		case "system":
			systemParts = append(systemParts, geminiPart{Text: m.Content})
		case "assistant":
			contents = append(contents, geminiContent{
				Role:  "model",
				Parts: []geminiPart{{Text: m.Content}},
			})
		default: // "user" or anything else
			contents = append(contents, geminiContent{
				Role:  "user",
				Parts: []geminiPart{{Text: m.Content}},
			})
		}
	}

	var sys *geminiContent
	if len(systemParts) > 0 {
		sys = &geminiContent{Parts: systemParts}
	}

	return contents, sys
}

// geminiBuildRequest constructs the geminiRequest with defaults.
func geminiBuildRequest(model *ModelConfig, messages []ChatMessage, withTools bool) geminiRequest {
	maxTokens := model.MaxTokens
	if maxTokens == 0 {
		maxTokens = 4096
	}

	contents, sys := geminiMessages(messages)

	req := geminiRequest{
		Contents: contents,
		GenerationConfig: &geminiGenConfig{
			MaxOutputTokens: maxTokens,
			Temperature:     0.7,
		},
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
				Parameters:  td.Function.Parameters,
			})
		}
		req.Tools = []geminiToolSet{{FunctionDeclarations: funcDecls}}
	}

	return req
}

// geminiDoRequest sends the request and returns the parsed response.
func geminiDoRequest(ctx context.Context, model *ModelConfig, reqBody geminiRequest) (*geminiResponse, error) {
	apiKey := ResolveAPIKey(model)
	if apiKey == "" {
		return nil, fmt.Errorf("no API key for provider '%s' — set %s",
			model.Provider, KeySourceLabel(model))
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	// Gemini endpoint: {base}/v1beta/models/{model}:generateContent
	endpoint := strings.TrimRight(model.BaseURL, "/") + "/v1beta/models/" + model.Model + ":generateContent"

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-goog-api-key", apiKey)

	resp, err := GetAIClient(model.BaseURL).Do(req)
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

	// Concatenate text parts
	var sb strings.Builder
	for _, part := range apiResp.Candidates[0].Content.Parts {
		if part.Text != "" {
			sb.WriteString(part.Text)
		}
	}

	TrackModelUsage(model.Model, apiResp.UsageMetadata.PromptTokenCount, apiResp.UsageMetadata.CandidatesTokenCount)
	RecordCost(model.Model, apiResp.UsageMetadata.PromptTokenCount, apiResp.UsageMetadata.CandidatesTokenCount)

	return sb.String(), nil
}

// callGeminiWithTools sends a request with native tool definitions.
func CallGeminiWithTools(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, []ToolCall, error) {
	reqBody := geminiBuildRequest(model, messages, true)
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
		if part.Text != "" {
			sb.WriteString(part.Text)
		}
		if part.FunctionCall != nil {
			args := part.FunctionCall.Args
			if args == nil {
				args = make(map[string]interface{})
			}
			toolCalls = append(toolCalls, ToolCall{
				Name: part.FunctionCall.Name,
				Args: args,
			})
			log.Printf("[agent] Gemini functionCall: %s(%v)", part.FunctionCall.Name, args)
		}
	}

	TrackModelUsage(model.Model, apiResp.UsageMetadata.PromptTokenCount, apiResp.UsageMetadata.CandidatesTokenCount)
	RecordCost(model.Model, apiResp.UsageMetadata.PromptTokenCount, apiResp.UsageMetadata.CandidatesTokenCount)

	return sb.String(), toolCalls, nil
}
