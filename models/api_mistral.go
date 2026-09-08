package models

import (
	"context"
)

// ──────────────────────────────────────────────
// Mistral AI Provider
// Dedicated provider for Mistral & Codestral models
// ──────────────────────────────────────────────

type MistralProvider struct {
	OpenAIProvider
}

func init() {
	p := &MistralProvider{}
	RegisterProvider(ProviderSpec{
		Name:             "mistral",
		Aliases:          []string{"codestral", "mistralai"},
		DisplayName:      "Mistral AI",
		DefaultBaseURL:   "https://api.mistral.ai/v1",
		DefaultAPIFormat: "mistral",
		KeyEnvs:          []string{"MISTRAL_API_KEY", "CODESTRAL_API_KEY"},
	}, p)

	RegisterCatalog("mistral", []CatalogEntry{
		{"codestral-latest", 32768, false, "codestral"},
		{"mistral-large-latest", 32768, true, "mistral-large"},
		{"mistral-small-latest", 16384, false, "mistral-small"},
		{"pixtral-large-latest", 32768, true, "pixtral-large"},
		{"pixtral-12b-2409", 16384, false, "pixtral-12b"},
		{"ministral-8b-latest", 16384, false, "ministral-8b"},
		{"ministral-3b-latest", 8192, false, "ministral-3b"},
	})
}

func (p *MistralProvider) Format() string {
	return "mistral"
}

func (p *MistralProvider) Call(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, error) {
	return CallOpenAI(ctx, model, messages)
}

func (p *MistralProvider) CallWithTools(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, []ToolCall, error) {
	return CallOpenAIWithTools(ctx, model, messages)
}

func (p *MistralProvider) CallStream(ctx context.Context, model *ModelConfig, messages []ChatMessage) (<-chan StreamChunk, error) {
	return CallOpenAIStream(ctx, model, messages)
}
