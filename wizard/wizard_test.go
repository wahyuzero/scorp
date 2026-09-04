package wizard

import (
	"testing"

	"scorp-agent/models"
)

func TestWizardViewsAndKeyboards(t *testing.T) {
	// Initialize sample model config
	models.ModelCfg = &models.ModelRouterConfig{
		DefaultModel:    "deepseek/deepseek-v4-flash",
		AgentModel:      "deepseek/deepseek-v4-flash",
		DelegationModel: "deepseek/deepseek-v4-flash",
		VisionModel:     "opencode/mimo-v2.5-free",
		Models: map[string]models.ModelConfig{
			"deepseek/deepseek-v4-flash": {
				Provider: "command-code",
				Model:    "deepseek/deepseek-v4-flash",
			},
			"opencode/mimo-v2.5-free": {
				Provider: "opencode",
				Model:    "mimo-v2.5-free",
			},
		},
		FallbackModels: []string{"opencode/mimo-v2.5-free"},
	}

	// 1. Test ModelMenuText
	txt := ModelMenuText()
	if txt == "" {
		t.Errorf("ModelMenuText returned empty string")
	}

	// 2. Test ModelMenuKeyboard
	kb := ModelMenuKeyboard()
	if kb == nil || kb["inline_keyboard"] == nil {
		t.Errorf("ModelMenuKeyboard returned nil keyboard")
	}

	// 3. Test OrNA helper
	if OrNA("") != "(not set)" {
		t.Errorf("expected '(not set)', got '%s'", OrNA(""))
	}
	if OrNA("custom-model") != "custom-model" {
		t.Errorf("expected 'custom-model', got '%s'", OrNA("custom-model"))
	}

	// 4. Test ModelInfoText
	info := ModelInfoText("opencode/mimo-v2.5-free")
	if info == "" {
		t.Errorf("ModelInfoText returned empty")
	}

	// 5. Test Fallback operations
	res := addToFallback("deepseek/deepseek-v4-flash")
	if res == "" {
		t.Errorf("addToFallback failed")
	}
	res2 := addToFallback("deepseek/deepseek-v4-flash")
	if res2 == "" {
		t.Errorf("expected already in fallback response")
	}

	resRem := removeFromFallback("deepseek/deepseek-v4-flash")
	if resRem == "" {
		t.Errorf("removeFromFallback failed")
	}
}
