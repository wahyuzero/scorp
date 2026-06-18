package main

import (
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
// Native Function Calling (OpenAI-compatible)
// ──────────────────────────────────────────────

// getNativeToolDefs returns tool definitions in OpenAI function calling format
func getNativeToolDefs() []toolDef {
	return []toolDef{
		{
			Type: "function",
			Function: toolFunction{
				Name:        "shell",
				Description: "Execute a shell command on the VPS. Use for system tasks, package management, service control, disk/memory checks, docker, etc.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"command": map[string]interface{}{
							"type":        "string",
							"description": "The shell command to execute",
						},
						"timeout": map[string]interface{}{
							"type":        "integer",
							"description": "Timeout in seconds (default 30)",
							"default":     30,
						},
					},
					"required": []string{"command"},
				},
			},
		},
		{
			Type: "function",
			Function: toolFunction{
				Name:        "read_file",
				Description: "Read the contents of a file.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"path": map[string]interface{}{
							"type":        "string",
							"description": "Path to the file",
						},
						"lines": map[string]interface{}{
							"type":        "integer",
							"description": "Max lines to read (optional)",
						},
					},
					"required": []string{"path"},
				},
			},
		},
		{
			Type: "function",
			Function: toolFunction{
				Name:        "write_file",
				Description: "Write content to a file. Creates parent directories if needed.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"path":    map[string]interface{}{"type": "string", "description": "File path"},
						"content": map[string]interface{}{"type": "string", "description": "File content"},
					},
					"required": []string{"path", "content"},
				},
			},
		},
		{
			Type: "function",
			Function: toolFunction{
				Name:        "list_dir",
				Description: "List directory contents with details.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"path":      map[string]interface{}{"type": "string", "description": "Directory path"},
						"recursive": map[string]interface{}{"type": "boolean", "description": "List recursively", "default": false},
					},
					"required": []string{"path"},
				},
			},
		},
		{
			Type: "function",
			Function: toolFunction{
				Name:        "system_info",
				Description: "Get system information: CPU, memory, disk, network, docker containers, services.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"type": map[string]interface{}{
							"type":        "string",
							"description": "Type of info: full, cpu, memory, disk, network, docker, services",
							"default":     "full",
						},
					},
				},
			},
		},
		{
			Type: "function",
			Function: toolFunction{
				Name:        "process",
				Description: "Inspect and manage processes and services: list, top, kill, service status/restart, ports.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"action": map[string]interface{}{
							"type":        "string",
							"description": "Action: list, top, kill, service_status, service_restart, service_list, ports",
						},
						"filter":  map[string]interface{}{"type": "string", "description": "Filter for list"},
						"service": map[string]interface{}{"type": "string", "description": "Service name for service_*"},
						"pid":     map[string]interface{}{"type": "string", "description": "PID for kill"},
						"sort_by": map[string]interface{}{"type": "string", "description": "Sort: mem or cpu"},
					},
					"required": []string{"action"},
				},
			},
		},
		{
			Type: "function",
			Function: toolFunction{
				Name:        "log",
				Description: "Fetch logs from docker containers, systemd journal, or files.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"source": map[string]interface{}{"type": "string", "description": "docker, journal, or file"},
						"target": map[string]interface{}{"type": "string", "description": "Container name, unit name, or file path"},
						"lines":  map[string]interface{}{"type": "integer", "description": "Number of lines (default 50)"},
					},
					"required": []string{"source", "target"},
				},
			},
		},
		{
			Type: "function",
			Function: toolFunction{
				Name:        "search_code",
				Description: "Search codebase using ripgrep (fast regex search).",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"pattern":     map[string]interface{}{"type": "string", "description": "Regex pattern"},
						"path":        map[string]interface{}{"type": "string", "description": "Search path", "default": "."},
						"glob":        map[string]interface{}{"type": "string", "description": "File glob filter"},
						"max_results": map[string]interface{}{"type": "integer", "description": "Max results", "default": 20},
					},
					"required": []string{"pattern"},
				},
			},
		},
		{
			Type: "function",
			Function: toolFunction{
				Name:        "http",
				Description: "Make HTTP requests to APIs or websites.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"method": map[string]interface{}{"type": "string", "description": "HTTP method", "default": "GET"},
						"url":    map[string]interface{}{"type": "string", "description": "URL to request"},
						"body":   map[string]interface{}{"type": "string", "description": "Request body (JSON)"},
						"headers": map[string]interface{}{"type": "string", "description": "JSON object of headers"},
						"timeout": map[string]interface{}{"type": "integer", "description": "Timeout seconds", "default": 15},
					},
					"required": []string{"url"},
				},
			},
		},
		{
			Type: "function",
			Function: toolFunction{
				Name:        "send_file",
				Description: "Send a file to the user via Telegram.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"path":    map[string]interface{}{"type": "string", "description": "File path"},
						"caption": map[string]interface{}{"type": "string", "description": "Optional caption"},
					},
					"required": []string{"path"},
				},
				},
				},
				{
				Type: "function",
				Function: toolFunction{
					Name:        "analyze_image",
					Description: "Analyze an image file using a vision model. Use after browser screenshots or for any image file analysis.",
					Parameters: map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"path": map[string]interface{}{"type": "string", "description": "Path to the image file (from browser screenshot or upload)"},
							"question": map[string]interface{}{"type": "string", "description": "What to look for in the image (default: describe in detail)"},
						},
						"required": []string{"path"},
					},
				},
				},
				}
}

// callModelWithTools calls the API with native function calling enabled.
// Returns: (textContent, nativeToolCalls, error)
func callModelWithTools(ctx context.Context, model *ModelConfig, messages []chatMessage) (string, []ToolCall, error) {
	if model == nil {
		return "", nil, fmt.Errorf("no model configured")
	}

	if model.APIKey == "" {
		return "", nil, fmt.Errorf("no API key configured for %s", model.Provider)
	}

	maxTokens := model.MaxTokens
	if maxTokens == 0 {
		maxTokens = 4096
	}

	reqBody := chatRequest{
		Model:       model.Model,
		Messages:    messages,
		MaxTokens:   maxTokens,
		Temperature: 0.7,
		Tools:       generateNativeToolsSchema(),
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
	httpReq.Header.Set("Authorization", "Bearer "+model.APIKey)
	if model.Provider == "openrouter" {
		httpReq.Header.Set("HTTP-Referer", "https://scorp-agent.local")
		httpReq.Header.Set("X-Title", "ScorpAgent")
	}

	resp, err := getAIClient(model.BaseURL).Do(httpReq)
	if err != nil {
		return "", nil, fmt.Errorf("API call failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, fmt.Errorf("read error: %w", err)
	}

	if resp.StatusCode != 200 {
		return "", nil, fmt.Errorf("API error (HTTP %d): %s", resp.StatusCode, truncateStr(string(body), 300))
	}

	var chatResp chatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", nil, fmt.Errorf("parse error: %s", truncateStr(string(body), 200))
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
		log.Printf("[agent] Native tool call: %s(%v)", tc.Function.Name, args)
	}

	trackModelUsage(model.Model, chatResp.Usage.PromptTokens, chatResp.Usage.CompletionTokens)
	recordCost(model.Model, chatResp.Usage.PromptTokens, chatResp.Usage.CompletionTokens)

	return content, toolCalls, nil
}

// isRateLimitError checks if an error is caused by HTTP 429 (rate limit) or 402 (quota exhausted).
func isRateLimitError(err error) bool {
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

// callModelWithToolsAndFallback tries with tools first, falls back to plain callModel.
// On rate-limit errors (429/402), retries with exponential backoff before giving up.
func callModelWithToolsAndFallback(ctx context.Context, taskType string, messages []chatMessage) (string, []ToolCall, string, error) {
	primary := routeModel(taskType)

	// Rate-limit retry loop: try primary model up to 3 times with backoff
	maxRateLimitRetries := 3
	backoffSteps := []time.Duration{15 * time.Second, 30 * time.Second, 60 * time.Second}

	if primary != nil {
		var lastErr error
		for attempt := 0; attempt <= maxRateLimitRetries; attempt++ {
			// Try with tools first
			text, toolCalls, err := callModelWithTools(ctx, primary, messages)
			if err == nil {
				return text, toolCalls, primary.Model, nil
			}
			lastErr = err

			if !isRateLimitError(err) {
				// Non-rate-limit error — try plain call, then fallback
				log.Printf("[models] Primary model %s with tools failed: %v, trying plain", primary.Model, err)
				text2, err2 := callModel(ctx, primary, messages)
				if err2 == nil {
					return text2, nil, primary.Model, nil
				}
				log.Printf("[models] Primary plain also failed: %v", err2)
				break
			}

			// Rate-limit error — backoff and retry
			if attempt < maxRateLimitRetries {
				wait := backoffSteps[attempt]
				log.Printf("[models] ⏳ Rate-limit detected (HTTP 429/402), backing off for %v before retry %d/%d",
					wait, attempt+1, maxRateLimitRetries)
				select {
				case <-time.After(wait):
				case <-ctx.Done():
					return "", nil, "", ctx.Err()
				}
			}
		}

		// If we exhausted rate-limit retries, try plain call one more time
		if isRateLimitError(lastErr) {
			log.Printf("[models] Rate-limit retries exhausted, trying plain call once more")
			text2, err2 := callModel(ctx, primary, messages)
			if err2 == nil {
				return text2, nil, primary.Model, nil
			}
			if !isRateLimitError(err2) {
				lastErr = err2
			}
		}
		_ = lastErr
	}

	// Fallback to default (chat) model
	fallback := routeModel("chat")
	if fallback != nil && (primary == nil || fallback.Model != primary.Model) {
		text, err := callModel(ctx, fallback, messages)
		if err == nil {
			return text, nil, fallback.Model, nil
		}
	}

	// All models failed — return descriptive error
	if primary != nil {
		return "", nil, "", fmt.Errorf("all models failed (primary: %s) — if rate-limited, wait and retry later", primary.Model)
	}
	return "", nil, "", fmt.Errorf("all models failed — no models configured")
}

// ──────────────────────────────────────────────
// Code-Block Fallback Parser
// ──────────────────────────────────────────────

// parseCodeBlockFallback extracts tool calls from markdown code blocks
// when the model doesn't use <tool_call> tags or native function calling.
//
// Pattern: ```shell\ncommand\n``` or ```\ncommand\n```
var codeBlockRe = regexp.MustCompile("(?s)```(?:shell|bash|sh)?\\s*\n(.*?)```")

func parseCodeBlockFallback(text string) ([]ToolCall, string) {
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

		// Take only the first line if multi-line (commands can be chained with &&)
		// Actually keep multi-line for complex commands
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

// parseAllToolCalls tries native tool calls, then XML tags, then code blocks
func parseAllToolCalls(text string, nativeCalls []ToolCall) ([]ToolCall, string) {
	// 1. Native function calling (already parsed by callModelWithTools)
	if len(nativeCalls) > 0 {
		return nativeCalls, text
	}

	// 2. XML tag format: <tool_call>{"name": "...", "args": {...}}</tool_call>
	xmlCalls, cleanText := parseToolCalls(text)
	if len(xmlCalls) > 0 {
		return xmlCalls, cleanText
	}

	// 3. Code-block fallback: ```shell\ncommand\n```
	codeCalls, cleanText := parseCodeBlockFallback(text)
	if len(codeCalls) > 0 {
		return codeCalls, cleanText
	}

	return nil, text
}
