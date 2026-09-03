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
	"regexp"
	"strings"
	"time"

	"scorp-agent/config"
	"scorp-agent/internal/helpers"
	"scorp-agent/registry"
)

// ──────────────────────────────────────────────
// Command Code API Provider (/alpha/generate)
// High-efficiency, cost-optimized gateway with prompt caching & streaming
// ──────────────────────────────────────────────

const (
	CommandCodeGenerateURL = "https://api.commandcode.ai/alpha/generate"
	CommandCodeVersion     = "1.39.2"
)

// commandCodePayload matches the Hono/Zod schema expected by Command Code backend
type commandCodePayload struct {
	Memory string                 `json:"memory"`
	Params commandCodeParams      `json:"params"`
	Config map[string]interface{} `json:"config"`
}

type commandCodeParams struct {
	Model           string             `json:"model"`
	Messages        []commandCodeMsg   `json:"messages"`
	Tools           []commandCodeTool  `json:"tools,omitempty"`
	ReasoningEffort string             `json:"reasoningEffort,omitempty"`
}

type commandCodeTool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"input_schema"`
}

type commandCodeMsg struct {
	Role    string        `json:"role"`
	Content []interface{} `json:"content"`
}

type textContentPart struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type toolCallContentPart struct {
	Type       string                 `json:"type"`
	ToolCallID string                 `json:"toolCallId"`
	ToolName   string                 `json:"toolName"`
	Input      map[string]interface{} `json:"input"`
}

// activeToolCall holds state for concurrent tool calls in SSE stream
type activeToolCall struct {
	id         string
	name       string
	jsonBuffer strings.Builder
	arguments  map[string]interface{}
	ended      bool
}

// buildCommandCodePayload transforms ChatMessages and optional tools into Command Code schema
func buildCommandCodePayload(model *ModelConfig, messages []ChatMessage, tools []registry.ToolSchema) (*commandCodePayload, error) {
	var systemParts []string
	var convertedMsgs []commandCodeMsg

	for _, msg := range messages {
		if msg.Role == "system" {
			if msg.Content != "" {
				systemParts = append(systemParts, msg.Content)
			}
			continue
		}

		// Assistant message with native tool calls
		if msg.Role == "assistant" && len(msg.ToolCalls) > 0 {
			var parts []interface{}
			if msg.Content != "" {
				parts = append(parts, textContentPart{Type: "text", Text: msg.Content})
			}
			for _, tc := range msg.ToolCalls {
				var args map[string]interface{}
				if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err != nil {
					args = map[string]interface{}{}
				}
				callID := tc.ID
				if callID == "" {
					callID = fmt.Sprintf("call_%d", time.Now().UnixNano())
				}
				parts = append(parts, toolCallContentPart{
					Type:       "tool-call",
					ToolCallID: callID,
					ToolName:   tc.Function.Name,
					Input:      args,
				})
			}
			convertedMsgs = append(convertedMsgs, commandCodeMsg{
				Role:    "assistant",
				Content: parts,
			})
			continue
		}

		// Standard user or assistant message
		role := msg.Role
		if role != "assistant" && role != "tool" {
			role = "user"
		}
		convertedMsgs = append(convertedMsgs, commandCodeMsg{
			Role: role,
			Content: []interface{}{
				textContentPart{Type: "text", Text: msg.Content},
			},
		})
	}

	// Build tools schema if requested
	var ccTools []commandCodeTool
	for _, t := range tools {
		ccTools = append(ccTools, commandCodeTool{
			Name:        t.Function.Name,
			Description: t.Function.Description,
			InputSchema: t.Function.Parameters,
		})
	}

	cwd, _ := os.Getwd()
	if cwd == "" {
		cwd = config.HomeDir()
	}

	payload := &commandCodePayload{
		Memory: strings.Join(systemParts, "\n\n"),
		Params: commandCodeParams{
			Model:    model.Model,
			Messages: convertedMsgs,
			Tools:    ccTools,
		},
		Config: map[string]interface{}{
			"workingDir":    cwd,
			"date":          time.Now().Format("2006-01-02"),
			"environment":   "linux",
			"structure":     []string{},
			"isGitRepo":     false,
			"currentBranch": "main",
			"mainBranch":    "main",
			"gitStatus":     "",
			"recentCommits": []string{},
		},
	}

	return payload, nil
}

// createCommandCodeRequest creates the HTTP request with full spoofing headers
func createCommandCodeRequest(ctx context.Context, apiKey string, payload *commandCodePayload) (*http.Request, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal command code payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", CommandCodeGenerateURL, bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "cli")
	req.Header.Set("x-command-code-version", CommandCodeVersion)
	req.Header.Set("x-cli-environment", "terminal")
	req.Header.Set("x-taste-learning", "false")
	req.Header.Set("x-session-id", "sess_scorp_cli")

	return req, nil
}

// CallCommandCode sends a non-tool completion request
func CallCommandCode(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, error) {
	content, _, err := CallCommandCodeWithTools(ctx, model, messages)
	return content, err
}

// CallCommandCodeWithTools sends a completion request with tools and parses SSE stream
func CallCommandCodeWithTools(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, []ToolCall, error) {
	apiKey := ResolveAPIKey(model)
	if apiKey == "" {
		return "", nil, fmt.Errorf("no API key for provider '%s' — set %s", model.Provider, KeySourceLabel(model))
	}

	payload, err := buildCommandCodePayload(model, messages, registry.GenerateNativeToolsSchema())
	if err != nil {
		return "", nil, err
	}

	req, err := createCommandCodeRequest(ctx, apiKey, payload)
	if err != nil {
		return "", nil, err
	}

	client := GetAIClient(CommandCodeGenerateURL)
	resp, err := client.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("API call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", nil, fmt.Errorf("API error (HTTP %d): %s", resp.StatusCode, helpers.TruncateStr(string(body), 300))
	}

	// Parse SSE Stream line-by-line
	scanner := bufio.NewScanner(resp.Body)
	// Allow 2MB buffer for large SSE chunks
	scanner.Buffer(make([]byte, 64*1024), 2*1024*1024)

	var textBuilder strings.Builder
	activeCalls := make(map[string]*activeToolCall)
	var toolCalls []ToolCall
	var promptTokens, completionTokens int

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "data:") {
			line = strings.TrimSpace(line[5:])
		}
		if line == "" || line == "[DONE]" {
			continue
		}

		var evt map[string]interface{}
		if err := json.Unmarshal([]byte(line), &evt); err != nil {
			continue
		}

		evtType, _ := evt["type"].(string)

		switch evtType {
		case "text-delta":
			if txt, ok := evt["text"].(string); ok {
				textBuilder.WriteString(txt)
			}
		case "tool-input-start":
			id, _ := evt["id"].(string)
			toolName, _ := evt["toolName"].(string)
			if id != "" {
				activeCalls[id] = &activeToolCall{
					id:        id,
					name:      toolName,
					arguments: make(map[string]interface{}),
				}
			}
		case "tool-input-delta":
			id, _ := evt["id"].(string)
			delta, _ := evt["delta"].(string)
			if ac, ok := activeCalls[id]; ok && delta != "" {
				ac.jsonBuffer.WriteString(delta)
			}
		case "tool-input-end":
			id, _ := evt["id"].(string)
			if ac, ok := activeCalls[id]; ok {
				if ac.jsonBuffer.Len() > 0 {
					var args map[string]interface{}
					if err := json.Unmarshal([]byte(ac.jsonBuffer.String()), &args); err == nil {
						ac.arguments = args
					}
				}
			}
		case "tool-call":
			id, _ := evt["toolCallId"].(string)
			if id == "" {
				id, _ = evt["id"].(string)
			}
			toolName, _ := evt["toolName"].(string)
			var args map[string]interface{}
			if inputMap, ok := evt["input"].(map[string]interface{}); ok {
				args = inputMap
			} else if ac, ok := activeCalls[id]; ok {
				args = ac.arguments
			} else {
				args = make(map[string]interface{})
			}

			if ac, ok := activeCalls[id]; ok {
				ac.ended = true
			}

			toolCalls = append(toolCalls, ToolCall{
				Name: toolName,
				Args: args,
			})
		case "finish-step", "finish":
			if usage, ok := evt["usage"].(map[string]interface{}); ok {
				if inTok, ok := usage["inputTokens"].(float64); ok {
					promptTokens = int(inTok)
				}
				if outTok, ok := usage["outputTokens"].(float64); ok {
					completionTokens = int(outTok)
				}
			}
		case "error":
			errMsg := "Command Code API stream error"
			if errObj, ok := evt["error"].(map[string]interface{}); ok {
				if m, ok := errObj["message"].(string); ok {
					errMsg = m
				}
			}
			return "", nil, fmt.Errorf("%s", errMsg)
		}
	}

	if err := scanner.Err(); err != nil {
		return "", nil, fmt.Errorf("read stream error: %w", err)
	}

	// Finalize any unfinalized tool calls in activeCalls that didn't emit a separate tool-call event
	for _, ac := range activeCalls {
		if !ac.ended {
			ac.ended = true
			if len(ac.arguments) == 0 && ac.jsonBuffer.Len() > 0 {
				var args map[string]interface{}
				if err := json.Unmarshal([]byte(ac.jsonBuffer.String()), &args); err == nil {
					ac.arguments = args
				}
			}
			toolCalls = append(toolCalls, ToolCall{
				Name: ac.name,
				Args: ac.arguments,
			})
		}
	}

	reply := textBuilder.String()

	// ── Auto-Healing Fallback Parser: DSML, XML, or pseudo bracket tool calls ──
	fallbackCalls, cleanReply := extractFallbackToolCalls(reply)
	if len(fallbackCalls) > 0 {
		log.Printf("[models/command-code] Detected %d fallback tool calls in text", len(fallbackCalls))
		toolCalls = append(toolCalls, fallbackCalls...)
		reply = cleanReply
	}

	// Track usage and cost
	if promptTokens > 0 || completionTokens > 0 {
		TrackModelUsage(model.Model, promptTokens, completionTokens)
		RecordCost(model.Model, promptTokens, completionTokens)
	}

	return reply, toolCalls, nil
}

// CallCommandCodeStream streams tokens as StreamChunk
func CallCommandCodeStream(ctx context.Context, model *ModelConfig, messages []ChatMessage) (<-chan StreamChunk, error) {
	apiKey := ResolveAPIKey(model)
	if apiKey == "" {
		return nil, fmt.Errorf("no API key for provider '%s' — set %s", model.Provider, KeySourceLabel(model))
	}

	payload, err := buildCommandCodePayload(model, messages, nil)
	if err != nil {
		return nil, err
	}

	req, err := createCommandCodeRequest(ctx, apiKey, payload)
	if err != nil {
		return nil, err
	}

	client := GetAIClient(CommandCodeGenerateURL)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API call failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("API error (HTTP %d): %s", resp.StatusCode, helpers.TruncateStr(string(body), 300))
	}

	ch := make(chan StreamChunk, 32)

	go func() {
		defer resp.Body.Close()
		defer close(ch)

		scanner := bufio.NewScanner(resp.Body)
		scanner.Buffer(make([]byte, 64*1024), 2*1024*1024)

		for scanner.Scan() {
			select {
			case <-ctx.Done():
				ch <- StreamChunk{Error: ctx.Err()}
				return
			default:
			}

			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}
			if strings.HasPrefix(line, "data:") {
				line = strings.TrimSpace(line[5:])
			}
			if line == "" || line == "[DONE]" {
				continue
			}

			var evt map[string]interface{}
			if err := json.Unmarshal([]byte(line), &evt); err != nil {
				continue
			}

			evtType, _ := evt["type"].(string)
			switch evtType {
			case "text-delta":
				if txt, ok := evt["text"].(string); ok && txt != "" {
					ch <- StreamChunk{Content: txt}
				}
			case "finish-step", "finish":
				ch <- StreamChunk{Finish: true}
			case "error":
				errMsg := "Command Code stream error"
				if errObj, ok := evt["error"].(map[string]interface{}); ok {
					if m, ok := errObj["message"].(string); ok {
						errMsg = m
					}
				}
				ch <- StreamChunk{Error: fmt.Errorf("%s", errMsg)}
				return
			}
		}

		if err := scanner.Err(); err != nil {
			ch <- StreamChunk{Error: err}
		}
	}()

	return ch, nil
}

// extractFallbackToolCalls parses DSML (<｜｜DSML｜｜invoke>), XML (<tool_call>), or bracket format
func extractFallbackToolCalls(text string) ([]ToolCall, string) {
	var toolCalls []ToolCall
	cleanText := text

	// 1. DSML Parser: <｜｜DSML｜｜invoke name="...">...</｜｜DSML｜｜invoke>
	dsmlRegex := regexp.MustCompile(`(?s)<[｜|]{2}DSML[｜|]{2}invoke name="([^"]+)">([\s\S]*?)<\/[｜|]{2}DSML[｜|]{2}invoke>`)
	matches := dsmlRegex.FindAllStringSubmatch(cleanText, -1)
	for _, m := range matches {
		toolName := m[1]
		body := m[2]
		params := make(map[string]interface{})

		paramRegex := regexp.MustCompile(`(?s)<[｜|]{2}DSML[｜|]{2}parameter name="([^"]+)"[^>]*>([\s\S]*?)<\/[｜|]{2}DSML[｜|]{2}parameter>`)
		pMatches := paramRegex.FindAllStringSubmatch(body, -1)
		for _, pm := range pMatches {
			pName := pm[1]
			pVal := strings.TrimSpace(pm[2])
			var parsedVal interface{}
			if err := json.Unmarshal([]byte(pVal), &parsedVal); err == nil {
				params[pName] = parsedVal
			} else {
				params[pName] = pVal
			}
		}
		toolCalls = append(toolCalls, ToolCall{Name: toolName, Args: params})
	}

	// 2. Generic XML: <tool_call>...</tool_call>
	xmlRegex := regexp.MustCompile(`(?s)<tool_call>([\s\S]*?)<\/tool_call>`)
	xmlMatches := xmlRegex.FindAllStringSubmatch(cleanText, -1)
	for _, m := range xmlMatches {
		body := strings.TrimSpace(m[1])
		var parsed struct {
			Name      string                 `json:"name"`
			Arguments map[string]interface{} `json:"arguments"`
		}
		if err := json.Unmarshal([]byte(body), &parsed); err == nil && parsed.Name != "" {
			toolCalls = append(toolCalls, ToolCall{Name: parsed.Name, Args: parsed.Arguments})
		}
	}

	// 3. Bracket: [Tool Call: name({...})]
	bracketRegex := regexp.MustCompile(`(?s)\[Tool Call:\s*([a-zA-Z0-9_-]+)\(([\s\S]*?)\)\]`)
	bMatches := bracketRegex.FindAllStringSubmatch(cleanText, -1)
	for _, m := range bMatches {
		toolName := strings.TrimSpace(m[1])
		rawArgs := strings.TrimSpace(m[2])
		args := make(map[string]interface{})
		if err := json.Unmarshal([]byte(rawArgs), &args); err == nil {
			toolCalls = append(toolCalls, ToolCall{Name: toolName, Args: args})
		}
	}

	// Clean out tags from user-facing text
	cleanText = regexp.MustCompile(`(?s)<[｜|]{2}DSML[｜|]{2}tool_calls>[\s\S]*?<\/[｜|]{2}DSML[｜|]{2}tool_calls>`).ReplaceAllString(cleanText, "")
	cleanText = dsmlRegex.ReplaceAllString(cleanText, "")
	cleanText = xmlRegex.ReplaceAllString(cleanText, "")
	cleanText = bracketRegex.ReplaceAllString(cleanText, "")
	cleanText = strings.TrimSpace(cleanText)

	return toolCalls, cleanText
}
