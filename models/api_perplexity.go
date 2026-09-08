package models

import (
	"context"
)

// ──────────────────────────────────────────────
// Perplexity Provider (Sonar Search Models)
// ──────────────────────────────────────────────

type PerplexityProvider struct {
	OpenAIProvider
}

func init() {
	p := &PerplexityProvider{}
	RegisterProvider(ProviderSpec{
		Name:             "perplexity",
		Aliases:          []string{"pplx"},
		DisplayName:      "Perplexity AI",
		DefaultBaseURL:   "https://api.perplexity.ai",
		DefaultAPIFormat: "perplexity",
		KeyEnvs:          []string{"PERPLEXITY_API_KEY", "PPLX_API_KEY"},
	}, p)

	RegisterCatalog("perplexity", []CatalogEntry{
		{"sonar", 4096, false, "sonar"},
		{"sonar-pro", 8192, true, "sonar-pro"},
		{"sonar-reasoning", 8192, true, "sonar-reasoning"},
		{"sonar-reasoning-pro", 16384, true, "sonar-reasoning-pro"},
	})
}

func (p *PerplexityProvider) Format() string {
	return "perplexity"
}

func (p *PerplexityProvider) Call(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, error) {
	return CallOpenAI(ctx, model, messages)
}

func (p *PerplexityProvider) CallWithTools(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, []ToolCall, error) {
	return CallOpenAIWithTools(ctx, model, messages)
}

func (p *PerplexityProvider) CallStream(ctx context.Context, model *ModelConfig, messages []ChatMessage) (<-chan StreamChunk, error) {
	return CallOpenAIStream(ctx, model, messages)
}
