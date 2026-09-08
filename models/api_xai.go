package models

import (
	"context"
)

// ──────────────────────────────────────────────
// xAI (Grok) Provider
// Dedicated provider for xAI Grok inference
// ──────────────────────────────────────────────

type XAIProvider struct {
	OpenAIProvider
}

func init() {
	p := &XAIProvider{}
	RegisterProvider(ProviderSpec{
		Name:             "xai",
		Aliases:          []string{"grok"},
		DisplayName:      "xAI (Grok)",
		DefaultBaseURL:   "https://api.x.ai/v1",
		DefaultAPIFormat: "xai",
		KeyEnvs:          []string{"XAI_API_KEY", "GROK_API_KEY"},
	}, p)

	RegisterCatalog("xai", []CatalogEntry{
		{"grok-2-latest", 32768, true, "grok-2"},
		{"grok-2-vision-1212", 32768, true, "grok-2-vision"},
		{"grok-beta", 16384, false, "grok-beta"},
		{"grok-3-preview", 32768, true, "grok-3"},
	})
}

func (p *XAIProvider) Format() string {
	return "xai"
}

func (p *XAIProvider) Call(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, error) {
	return CallOpenAI(ctx, model, messages)
}

func (p *XAIProvider) CallWithTools(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, []ToolCall, error) {
	return CallOpenAIWithTools(ctx, model, messages)
}

func (p *XAIProvider) CallStream(ctx context.Context, model *ModelConfig, messages []ChatMessage) (<-chan StreamChunk, error) {
	return CallOpenAIStream(ctx, model, messages)
}
