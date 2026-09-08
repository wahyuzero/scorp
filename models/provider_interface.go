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

// StreamingProvider is an optional interface for providers supporting SSE token streaming.
type StreamingProvider interface {
	CallStream(ctx context.Context, model *ModelConfig, messages []ChatMessage) (<-chan StreamChunk, error)
}

// KeyResolver is an optional interface for providers with custom credential resolution
// (e.g. reading credentials from disk, CLI sessions, SQLite databases, or OAuth stores).
type KeyResolver interface {
	ResolveKey(cfg *ModelConfig) string
}

// ProviderSpec holds metadata and preset configuration for an LLM provider.
type ProviderSpec struct {
	Name             string   // Canonical provider identifier (e.g. "mistral", "bedrock")
	Aliases          []string // Alternative identifiers or aliases
	DisplayName      string   // Human-readable title (e.g. "Mistral AI")
	DefaultBaseURL   string   // Default API endpoint URL
	DefaultAPIFormat string   // "openai" | "anthropic" | "gemini" | custom format
	KeyEnvs          []string // Environment variable names checked in order
	NoAuth           bool     // True for local models without auth (e.g. ollama)
	ExtraHeaders     bool     // True if vendor requires special headers (e.g. openrouter)
}

var (
	providersMu      sync.RWMutex
	providerAdapters = make(map[string]LLMProvider)
	defaultOpenAI    = &OpenAIProvider{}
)

// RegisterProvider registers a provider with its metadata and implementation adapter.
// This enables adding new providers simply by creating a new .go file with an init() hook,
// without modifying any core engine or routing files.
func RegisterProvider(spec ProviderSpec, adapter LLMProvider) {
	providersMu.Lock()
	defer providersMu.Unlock()

	canonicalName := strings.ToLower(strings.TrimSpace(spec.Name))
	if canonicalName != "" && adapter != nil {
		providerAdapters[canonicalName] = adapter
	}

	apiFormat := strings.ToLower(strings.TrimSpace(spec.DefaultAPIFormat))
	if apiFormat != "" && adapter != nil {
		if _, exists := providerAdapters[apiFormat]; !exists {
			providerAdapters[apiFormat] = adapter
		}
	}

	// Register aliases
	for _, alias := range spec.Aliases {
		a := strings.ToLower(strings.TrimSpace(alias))
		if a != "" && adapter != nil {
			providerAdapters[a] = adapter
		}
	}

	// Also register into ProviderRegistry presets
	preset := ProviderPreset{
		KeyEnvs:      spec.KeyEnvs,
		BaseURL:      spec.DefaultBaseURL,
		API:          spec.DefaultAPIFormat,
		NoAuth:       spec.NoAuth,
		ExtraHeaders: spec.ExtraHeaders,
		DisplayName:  spec.DisplayName,
	}

	if canonicalName != "" {
		ProviderRegistry[canonicalName] = preset
	}
	for _, alias := range spec.Aliases {
		a := strings.ToLower(strings.TrimSpace(alias))
		if a != "" {
			ProviderRegistry[a] = preset
		}
	}
}

// RegisterProviderAdapter registers an LLMProvider implementation for a format.
func RegisterProviderAdapter(format string, adapter LLMProvider) {
	providersMu.Lock()
	defer providersMu.Unlock()
	providerAdapters[strings.ToLower(strings.TrimSpace(format))] = adapter
}

// GetProvider returns the matching LLMProvider for a given API format or provider name.
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
