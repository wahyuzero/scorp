package bootstrap

import (
	"scorp-agent/registry"
	"scorp-agent/tools"
)


// Registry entry for provider tools
func init() {
	// ── Provider Management (Phase 2) ──
	registry.RegisterTool(registry.ToolDef{
		Name:        "provider",
		Description: "Manage LLM providers and models: list available providers, add/remove models, test connectivity",
		Category:    "system",
		Native:      true,
		Execute:     tools.HandleProviderCommand,
		Arguments: map[string]registry.ArgDef{
			"provider_action": {Type: "string", Description: "Action: list, add, test, remove, models", Required: true, Enum: []string{"list", "add", "test", "remove", "models"}},
			"provider_name":   {Type: "string", Description: "Provider name (e.g., deepseek, groq, openrouter, ollama, zai)"},
			"model_name":      {Type: "string", Description: "Model config name for test/remove actions"},
		},
	})
}