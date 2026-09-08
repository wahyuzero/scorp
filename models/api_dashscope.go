package models

import (
	"context"
)

// ──────────────────────────────────────────────
// Alibaba DashScope & Qwen Provider
// Supports DashScope portal, international endpoints, and Alibaba Coding Plan
// ──────────────────────────────────────────────

type DashScopeProvider struct {
	OpenAIProvider
}

func init() {
	p := &DashScopeProvider{}
	RegisterProvider(ProviderSpec{
		Name:             "qwen",
		Aliases:          []string{"dashscope", "qwen-portal", "qwen-intl", "qwen-us", "alibaba-coding"},
		DisplayName:      "Alibaba Qwen (DashScope)",
		DefaultBaseURL:   "https://dashscope.aliyuncs.com/compatible-mode/v1",
		DefaultAPIFormat: "qwen",
		KeyEnvs:          []string{"DASHSCOPE_API_KEY", "ALIBABA_API_KEY", "QWEN_API_KEY"},
	}, p)

	RegisterCatalog("qwen", []CatalogEntry{
		{"qwen-max-latest", 8192, true, "qwen-max"},
		{"qwen-plus-latest", 8192, false, "qwen-plus"},
		{"qwen-turbo-latest", 8192, false, "qwen-turbo"},
		{"qwen2.5-coder-32b-instruct", 8192, false, "qwen-coder-32b"},
		{"qwen2.5-72b-instruct", 8192, true, "qwen-72b"},
		{"qvq-72b-preview", 8192, true, "qvq-72b"},
	})
}

func (p *DashScopeProvider) Format() string {
	return "qwen"
}

func (p *DashScopeProvider) Call(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, error) {
	return CallOpenAI(ctx, model, messages)
}

func (p *DashScopeProvider) CallWithTools(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, []ToolCall, error) {
	return CallOpenAIWithTools(ctx, model, messages)
}

func (p *DashScopeProvider) CallStream(ctx context.Context, model *ModelConfig, messages []ChatMessage) (<-chan StreamChunk, error) {
	return CallOpenAIStream(ctx, model, messages)
}
