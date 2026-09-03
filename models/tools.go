package models

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"
)

// ──────────────────────────────────────────────
// Shared Model Calling Functions
// ──────────────────────────────────────────────

// CallModelWithTools calls a model with native tool definitions via the modular LLMProvider
func CallModelWithTools(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, []ToolCall, error) {
	if model == nil {
		return "", nil, fmt.Errorf("no model configured")
	}

	provider := GetProvider(ResolveAPIFormat(model))
	return provider.CallWithTools(ctx, model, messages)
}



// IsRateLimitError checks if an error indicates rate limiting (429/402)
func IsRateLimitError(err error) bool {
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



// CallModelWithToolsAndFallback tries with tools first, falls back to plain models.CallModel.
// On rate-limit errors (429/402), retries with exponential backoff before giving up.
// Supports N-tier fallback chain from config (fallback_models).
func CallModelWithToolsAndFallback(ctx context.Context, taskType string, messages []ChatMessage) (string, []ToolCall, string, error) {
	ModelCfgMu.RLock()
	cfg := ModelCfg
	ModelCfgMu.RUnlock()

	if cfg == nil {
		return "", nil, "", fmt.Errorf("no model config loaded")
	}

	// Build ordered list of models to try (as pointers, not name strings)
	var modelList []*ModelConfig
	var primaryLabel string

	// 1. Primary model (cost-aware)
	primary := RouteModelCostAware(taskType)
	if primary != nil {
		primaryLabel = primary.Model
		modelList = append(modelList, primary)
	}

	// 2. Fallback models from config (looked up by map key)
	for _, name := range cfg.FallbackModels {
		m := GetModelByName(name)
		if m == nil {
			log.Printf("[models] Fallback model '%s' not found in config, skipping", name)
			continue
		}
		dup := false
		for _, existing := range modelList {
			if existing.Model == m.Model {
				dup = true
				break
			}
		}
		if !dup {
			modelList = append(modelList, m)
		}
	}

	// 3. Default chat model as last resort
	if dm := RouteModel("chat"); dm != nil {
		dup := false
		for _, existing := range modelList {
			if existing.Model == dm.Model {
				dup = true
				break
			}
		}
		if !dup {
			modelList = append(modelList, dm)
		}
	}

	if len(modelList) == 0 {
		return "", nil, "", fmt.Errorf("no models available")
	}

	// Try each model in order
	var lastErr error
	for i, model := range modelList {
		label := model.Model

		if i > 0 {
			log.Printf("[models] Trying fallback model %s (%d/%d) with tools", label, i+1, len(modelList))
		}

		// Rate-limit retry loop for this model
		maxRateLimitRetries := 3
		backoffSteps := []time.Duration{15 * time.Second, 30 * time.Second, 60 * time.Second}
		var modelLastErr error

		for attempt := 0; attempt <= maxRateLimitRetries; attempt++ {
			// Try with tools first
			text, toolCalls, err := CallModelWithTools(ctx, model, messages)
			if err == nil {
				if i == 0 {
					return text, toolCalls, label, nil
				}
				return text, toolCalls, label + " (fallback)", nil
			}
			modelLastErr = err



			if !IsRateLimitError(err) {
				// Non-rate-limit error — try plain call, then continue to next fallback
				log.Printf("[models] Model %s with tools failed: %v, trying plain", label, err)
				text2, err2 := CallModel(ctx, model, messages)
				if err2 == nil {
					if i == 0 {
						return text2, nil, label, nil
					}
					return text2, nil, label + " (fallback)", nil
				}
				log.Printf("[models] Model %s plain also failed: %v", label, err2)
				modelLastErr = err2
				break
			}

			// Rate-limit error — backoff and retry
			if attempt < maxRateLimitRetries {
				wait := backoffSteps[attempt]
				log.Printf("[models] ⏳ Rate-limit detected on %s (HTTP 429/402), backing off for %v before retry %d/%d",
					label, wait, attempt+1, maxRateLimitRetries)
				select {
				case <-time.After(wait):
				case <-ctx.Done():
					return "", nil, "", ctx.Err()
				}
			}
		}

		// If we exhausted rate-limit retries, try plain call one more time
		if IsRateLimitError(modelLastErr) {
			log.Printf("[models] Rate-limit retries exhausted for %s, trying plain call once more", label)
			text2, err2 := CallModel(ctx, model, messages)
			if err2 == nil {
				if i == 0 {
					return text2, nil, label, nil
				}
				return text2, nil, label + " (fallback)", nil
			}
			if !IsRateLimitError(err2) {
				modelLastErr = err2
			}
		}

		lastErr = modelLastErr
		log.Printf("[models] Model %s failed after all retries: %v", label, modelLastErr)

		// Check if error type should trigger fallback
		if !ShouldFallbackOnError(lastErr, cfg.FallbackOnError) {
			log.Printf("[models] Error type does not match fallback triggers, stopping chain")
			break
		}
	}

	// All models failed
	return "", nil, "", fmt.Errorf("all models failed (primary: %s): %w", primaryLabel, lastErr)
}



// ──────────────────────────────────────────────
// Code-Block Fallback Parser
// ──────────────────────────────────────────────

var codeBlockRe = regexp.MustCompile("(?s)`{3}(?:shell|bash|sh)?\\s*\\n(.*?)`{3}")
func ParseCodeBlockFallback(text string) ([]ToolCall, string) {
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

// ParseAllToolCalls tries native tool calls, then XML tags, then code blocks
func ParseAllToolCalls(text string, nativeCalls []ToolCall) ([]ToolCall, string) {
	// 1. Native function calling (already parsed by CallModelWithTools)
	if len(nativeCalls) > 0 {
		return nativeCalls, text
	}

	// 2. XML tag format: <tool_call>{"name": "...", "args": {...}}
	xmlCalls, cleanText := ParseToolCalls(text)
	if len(xmlCalls) > 0 {
		return xmlCalls, cleanText
	}

	// 3. Code-block fallback: parse shell commands from code blocks
	codeCalls, cleanText := ParseCodeBlockFallback(text)
	if len(codeCalls) > 0 {
		return codeCalls, cleanText
	}

	return nil, text
}



// ──────────────────────────────────────────────
// Tool Call Parser (XML format)
// ──────────────────────────────────────────────

var ToolCallRe = regexp.MustCompile(`<tool_call>(.*?)</tool_call>`)

// ParseToolCalls extracts tool calls from LLM response in XML format
func ParseToolCalls(text string) ([]ToolCall, string) {
	matches := ToolCallRe.FindAllStringSubmatch(text, -1)
	if len(matches) == 0 {
		return nil, text
	}

	var calls []ToolCall
	cleanText := text

	for _, m := range matches {
		cleanText = strings.Replace(cleanText, m[0], "", 1)
		var tc ToolCall
		if err := json.Unmarshal([]byte(m[1]), &tc); err != nil {
			log.Printf("[scorp] Failed to parse tool call: %v", err)
			continue
		}
		calls = append(calls, tc)
	}

	return calls, strings.TrimSpace(cleanText)
}

