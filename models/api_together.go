package models

import (
	"context"
)

// ──────────────────────────────────────────────
// Fast Open-Source Model Inference Providers
// Covers: Together AI, Fireworks AI, SiliconFlow,
// Cerebras, SambaNova, Novita AI, DeepInfra
// ──────────────────────────────────────────────

type FastInferenceProvider struct {
	OpenAIProvider
}

func init() {
	p := &FastInferenceProvider{}

	// Together AI
	RegisterProvider(ProviderSpec{
		Name:             "together",
		Aliases:          []string{"togetherai"},
		DisplayName:      "Together AI",
		DefaultBaseURL:   "https://api.together.xyz/v1",
		DefaultAPIFormat: "openai",
		KeyEnvs:          []string{"TOGETHER_API_KEY", "TOGETHERAI_API_KEY"},
	}, p)
	RegisterCatalog("together", []CatalogEntry{
		{"meta-llama/Llama-3.3-70B-Instruct-Turbo", 16384, false, "llama-3.3-70b"},
		{"deepseek-ai/DeepSeek-V3", 8192, false, "deepseek-v3"},
		{"deepseek-ai/DeepSeek-R1", 16384, true, "deepseek-r1"},
		{"Qwen/Qwen2.5-Coder-32B-Instruct", 16384, false, "qwen-coder-32b"},
	})

	// Fireworks AI
	RegisterProvider(ProviderSpec{
		Name:             "fireworks",
		Aliases:          []string{"fireworks-ai"},
		DisplayName:      "Fireworks AI",
		DefaultBaseURL:   "https://api.fireworks.ai/inference/v1",
		DefaultAPIFormat: "openai",
		KeyEnvs:          []string{"FIREWORKS_API_KEY"},
	}, p)
	RegisterCatalog("fireworks", []CatalogEntry{
		{"accounts/fireworks/models/llama-v3p3-70b-instruct", 16384, false, "llama-3.3-70b"},
		{"accounts/fireworks/models/deepseek-v3", 8192, false, "deepseek-v3"},
		{"accounts/fireworks/models/deepseek-r1", 16384, true, "deepseek-r1"},
		{"accounts/fireworks/models/qwen2p5-coder-32b-instruct", 16384, false, "qwen-coder-32b"},
	})

	// SiliconFlow
	RegisterProvider(ProviderSpec{
		Name:             "siliconflow",
		Aliases:          []string{"silicon-cloud"},
		DisplayName:      "SiliconFlow",
		DefaultBaseURL:   "https://api.siliconflow.cn/v1",
		DefaultAPIFormat: "openai",
		KeyEnvs:          []string{"SILICONFLOW_API_KEY", "SILICON_API_KEY"},
	}, p)
	RegisterCatalog("siliconflow", []CatalogEntry{
		{"deepseek-ai/DeepSeek-V3", 8192, false, "deepseek-v3"},
		{"deepseek-ai/DeepSeek-R1", 16384, true, "deepseek-r1"},
		{"Qwen/Qwen2.5-Coder-32B-Instruct", 8192, false, "qwen-coder-32b"},
		{"Pro/deepseek-ai/DeepSeek-V3", 8192, true, "deepseek-v3-pro"},
	})

	// Cerebras
	RegisterProvider(ProviderSpec{
		Name:             "cerebras",
		DisplayName:      "Cerebras",
		DefaultBaseURL:   "https://api.cerebras.ai/v1",
		DefaultAPIFormat: "openai",
		KeyEnvs:          []string{"CEREBRAS_API_KEY"},
	}, p)
	RegisterCatalog("cerebras", []CatalogEntry{
		{"llama3.3-70b", 8192, false, "llama3.3-70b"},
		{"llama3.1-8b", 8192, false, "llama3.1-8b"},
		{"deepseek-r1-distill-llama-70b", 8192, true, "r1-distill-70b"},
	})

	// SambaNova
	RegisterProvider(ProviderSpec{
		Name:             "sambanova",
		DisplayName:      "SambaNova Systems",
		DefaultBaseURL:   "https://api.sambanova.ai/v1",
		DefaultAPIFormat: "openai",
		KeyEnvs:          []string{"SAMBANOVA_API_KEY"},
	}, p)

	// Novita AI
	RegisterProvider(ProviderSpec{
		Name:             "novita",
		Aliases:          []string{"novita-ai"},
		DisplayName:      "Novita AI",
		DefaultBaseURL:   "https://api.novita.ai/v3/openai",
		DefaultAPIFormat: "openai",
		KeyEnvs:          []string{"NOVITA_API_KEY"},
	}, p)

	// DeepInfra
	RegisterProvider(ProviderSpec{
		Name:             "deepinfra",
		DisplayName:      "DeepInfra",
		DefaultBaseURL:   "https://api.deepinfra.com/v1/openai",
		DefaultAPIFormat: "openai",
		KeyEnvs:          []string{"DEEPINFRA_API_KEY"},
	}, p)
}

func (p *FastInferenceProvider) Format() string {
	return "openai"
}

func (p *FastInferenceProvider) Call(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, error) {
	return CallOpenAI(ctx, model, messages)
}

func (p *FastInferenceProvider) CallWithTools(ctx context.Context, model *ModelConfig, messages []ChatMessage) (string, []ToolCall, error) {
	return CallOpenAIWithTools(ctx, model, messages)
}

func (p *FastInferenceProvider) CallStream(ctx context.Context, model *ModelConfig, messages []ChatMessage) (<-chan StreamChunk, error) {
	return CallOpenAIStream(ctx, model, messages)
}
