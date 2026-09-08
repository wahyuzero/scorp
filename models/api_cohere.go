package models

import (
	"context"
)

// ──────────────────────────────────────────────
// Cohere Provider (Command-R Series)
// ──────────────────────────────────────────────

type CohereProvider struct {
	OpenAIProvider
}

func init() {
	p := &CohereProvider{}
	RegisterProvider(ProviderSpec{
		Name:             "cohere",
		Aliases:          []string{"cohere-v2"},
		DisplayName:      "Cohere",
		DefaultBaseURL:   "https://api.cohere.com/v2",
		DefaultAPIFormat: "cohere",
		KeyEnvs:          []string{"COHERE_API_KEY", "CO_API_KEY"},
	}, p)

	RegisterCatalog("cohere", []CatalogEntry{
		{"command-r-plus-08-2024", 4096, true, "command-r-plus"},
		{"command-r-08-2024", 4096, false, "command-r"},
		{"command-r7b-12-2024", 4096, false, "command-r7b"},
	})
}

func (p *CohereProvider) Format() string {
	return "cohere"
}

func (p *CohereProvider) Call(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, error) {
	return CallOpenAI(ctx, model, messages)
}

func (p *CohereProvider) CallWithTools(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, []ToolCall, error) {
	return CallOpenAIWithTools(ctx, model, messages)
}

func (p *CohereProvider) CallStream(ctx context.Context, model *ModelConfig, messages []ChatMessage) (<-chan StreamChunk, error) {
	return CallOpenAIStream(ctx, model, messages)
}
