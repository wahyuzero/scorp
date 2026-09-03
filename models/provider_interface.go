package models

import (
	"context"
	"strings"
	"sync"
)

// ──────────────────────────────────────────────
// LLM Provider Architecture (Clean Strategy Pattern)
// Isolates each provider into its own file so changes
// or bug fixes to Vendor A never affect Vendor B.
// ──────────────────────────────────────────────

// LLMProvider defines the uniform interface all AI inference backends must satisfy.
type LLMProvider interface {
	Format() string
	Call(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, error)
	CallWithTools(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, []ToolCall, error)
}

var (
	providersMu       sync.RWMutex
	providerAdapters  = make(map[string]LLMProvider)
	defaultOpenAI     = &OpenAIProvider{}
)

func init() {
	RegisterProviderAdapter("openai", &OpenAIProvider{})
	RegisterProviderAdapter("anthropic", &AnthropicProvider{})
	RegisterProviderAdapter("gemini", &GeminiProvider{})
	RegisterProviderAdapter("command-code", &CommandCodeProvider{})
	RegisterProviderAdapter("commandcode", &CommandCodeProvider{})
}

// RegisterProviderAdapter registers an LLMProvider implementation for a format.
func RegisterProviderAdapter(format string, adapter LLMProvider) {
	providersMu.Lock()
	defer providersMu.Unlock()
	providerAdapters[strings.ToLower(strings.TrimSpace(format))] = adapter
}

// GetProvider returns the matching LLMProvider for a given API format.
// Defaults to OpenAI-compatible provider if not explicitly matched.
func GetProvider(format string) LLMProvider {
	providersMu.RLock()
	defer providersMu.RUnlock()

	key := strings.ToLower(strings.TrimSpace(format))
	if p, ok := providerAdapters[key]; ok {
		return p
	}
	return defaultOpenAI
}
