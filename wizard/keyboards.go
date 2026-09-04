package wizard

import (
	"fmt"
	"sort"
	"strings"

	"scorp-agent/models"
)

// ──────────────────────────────────────────────
// Keyboard Generators for Telegram Model Manager
// ──────────────────────────────────────────────

func ModelMenuKeyboard() map[string]interface{} {
	models.ModelCfgMu.RLock()
	defer models.ModelCfgMu.RUnlock()

	var rows [][]map[string]string

	// ── Role cards at top — 1 click → model picker ──
	primaryLabel := "💬 Primary: " + OrNA(models.ModelCfg.DefaultModel)
	agentLabel := "🤖 Agent: " + OrNA(models.ModelCfg.AgentModel)
	if len(primaryLabel) > 30 {
		primaryLabel = primaryLabel[:30] + "…"
	}
	if len(agentLabel) > 30 {
		agentLabel = agentLabel[:30] + "…"
	}
	rows = append(rows, []map[string]string{
		{"text": primaryLabel, "callback_data": "mdl:pick:role:chat:_"},
		{"text": agentLabel, "callback_data": "mdl:pick:role:agent:_"},
	})
	delegLabel := "🎯 Delegation: " + OrNA(models.ModelCfg.DelegationModel)
	if len(delegLabel) > 30 {
		delegLabel = delegLabel[:30] + "…"
	}
	visionLabel := "👁 Vision: " + OrNA(models.ModelCfg.VisionModel)
	if len(visionLabel) > 30 {
		visionLabel = visionLabel[:30] + "…"
	}
	rows = append(rows, []map[string]string{
		{"text": delegLabel, "callback_data": "mdl:pick:role:delegation:_"},
		{"text": visionLabel, "callback_data": "mdl:pick:role:vision:_"},
	})

	// ── Model list — tap for details ──
	var modelNames []string
	for name := range models.ModelCfg.Models {
		modelNames = append(modelNames, name)
	}
	sort.Strings(modelNames)

	for _, name := range modelNames {
		m := models.ModelCfg.Models[name]
		btn := map[string]string{
			"text":          name + "  (" + m.Provider + ")",
			"callback_data": "mdl:info:" + name,
		}
		rows = append(rows, []map[string]string{btn})
	}

	// ── Actions ──
	rows = append(rows, []map[string]string{
		{"text": "🔌 Add Provider", "callback_data": "mdl:aprov"},
		{"text": "🔧 Add Model", "callback_data": "mdl:acustom"},
	})
	rows = append(rows, []map[string]string{
		{"text": "🔄 Fallback", "callback_data": "mdl:fb"},
		{"text": "🔍 Health Check", "callback_data": "mdl:check"},
	})
	rows = append(rows, []map[string]string{
		{"text": "✏️ API Keys", "callback_data": "mdl:keys"},
		{"text": "⬅️ Menu", "callback_data": "mn:main"},
	})

	return map[string]interface{}{"inline_keyboard": rows}
}

// ProviderListKeyboard — for "Add Provider" flow
func ProviderListKeyboard() map[string]interface{} {
	var rows [][]map[string]string

	builtIn := make([]string, 0)
	custom := make([]string, 0)
	for _, n := range AllProviderNames() {
		if _, ok := models.ProviderRegistry[n]; ok {
			builtIn = append(builtIn, n)
		} else {
			custom = append(custom, n)
		}
	}

	for i := 0; i < len(builtIn); i++ {
		name := builtIn[i]
		disp := ProviderDisplayName(name)
		icon := "❌"
		if models.ProviderHasAPIKey(name) {
			icon = "✅"
		}
		catIcon := ""
		if models.HasCatalog(name) {
			catIcon = "📦"
		}
		btn := map[string]string{
			"text":          icon + " " + disp + " " + catIcon,
			"callback_data": "mdl:pk:" + name,
		}
		if i+1 < len(builtIn) {
			name2 := builtIn[i+1]
			disp2 := ProviderDisplayName(name2)
			icon2 := "❌"
			if models.ProviderHasAPIKey(name2) {
				icon2 = "✅"
			}
			catIcon2 := ""
			if models.HasCatalog(name2) {
				catIcon2 = "📦"
			}
			btn2 := map[string]string{
				"text":          icon2 + " " + disp2 + " " + catIcon2,
				"callback_data": "mdl:pk:" + name2,
			}
			rows = append(rows, []map[string]string{btn, btn2})
			i++
		} else {
			rows = append(rows, []map[string]string{btn})
		}
	}

	if len(custom) > 0 {
		for _, name := range custom {
			icon := "❌"
			if models.ProviderHasAPIKey(name) {
				icon = "✅"
			}
			rows = append(rows, []map[string]string{
				{"text": icon + " 🔧 " + ProviderDisplayName(name), "callback_data": "mdl:pk:" + name},
			})
		}
	}

	rows = append(rows, []map[string]string{
		{"text": "➕ New Custom Provider", "callback_data": "mdl:newprov"},
	})
	rows = append(rows, []map[string]string{
		{"text": "⬅️ Back", "callback_data": "mdl:menu"},
	})

	return map[string]interface{}{"inline_keyboard": rows}
}

// ModelPickerKeyboard — pick a model for role assignment or fallback
func ModelPickerKeyboard(purpose string) map[string]interface{} {
	models.ModelCfgMu.RLock()
	defer models.ModelCfgMu.RUnlock()

	var modelNames []string
	for name := range models.ModelCfg.Models {
		modelNames = append(modelNames, name)
	}
	sort.Strings(modelNames)

	var rows [][]map[string]string

	for _, name := range modelNames {
		m := models.ModelCfg.Models[name]

		// Mark currently active model for this role
		isCurrent := false
		role := strings.TrimPrefix(purpose, "role:")
		if strings.HasPrefix(purpose, "role:") {
			switch role {
			case "chat":
				isCurrent = (name == models.ModelCfg.DefaultModel)
			case "agent":
				isCurrent = (name == models.ModelCfg.AgentModel)
			case "delegation":
				isCurrent = (name == models.ModelCfg.DelegationModel)
			case "vision":
				isCurrent = (name == models.ModelCfg.VisionModel)
			}
		}

		label := name + "  (" + m.Provider + ")"
		prefix := "  "
		if isCurrent {
			prefix = "● "
		}

		var cb string
		if strings.HasPrefix(purpose, "role:") {
			cb = "mdl:rdo:" + role + ":" + name
		} else if purpose == "fallback" {
			cb = "mdl:fba:" + name
		} else {
			cb = "mdl:info:" + name
		}

		btn := map[string]string{
			"text":          prefix + label,
			"callback_data": cb,
		}
		rows = append(rows, []map[string]string{btn})
	}

	backTarget := "mdl:menu"
	if purpose == "fallback" {
		backTarget = "mdl:fb"
	} else if strings.HasPrefix(purpose, "role:") {
		backTarget = "mdl:menu"
	}
	rows = append(rows, []map[string]string{
		{"text": "⬅️ Back", "callback_data": backTarget},
	})

	return map[string]interface{}{"inline_keyboard": rows}
}

// ModelInfoKeyboard shows actions for a specific model
func ModelInfoKeyboard(name string) map[string]interface{} {
	var rows [][]map[string]string

	rows = append(rows, []map[string]string{
		{"text": "💬 Primary", "callback_data": "mdl:use:" + name},
		{"text": "🤖 Agent", "callback_data": "mdl:ag:" + name},
		{"text": "🎯 Delegation", "callback_data": "mdl:dlg:" + name},
		{"text": "👁 Vision", "callback_data": "mdl:vis:" + name},
	})

	rows = append(rows, []map[string]string{
		{"text": "🔄 Add to Fallback", "callback_data": "mdl:fba:" + name},
		{"text": "🔑 Set API Key", "callback_data": "mdl:key:" + name},
	})

	rows = append(rows, []map[string]string{
		{"text": "🗑️ Delete", "callback_data": "mdl:del:" + name},
		{"text": "⬅️ Back", "callback_data": "mdl:menu"},
	})

	return map[string]interface{}{"inline_keyboard": rows}
}

// FallbackEditorKeyboard — shows fallback chain with reorder controls
func FallbackEditorKeyboard() map[string]interface{} {
	models.ModelCfgMu.RLock()
	defer models.ModelCfgMu.RUnlock()

	var rows [][]map[string]string

	if len(models.ModelCfg.FallbackModels) == 0 {
		rows = append(rows, []map[string]string{
			{"text": "📭 No fallback models", "callback_data": "mdl:fb"},
		})
	} else {
		for i, name := range models.ModelCfg.FallbackModels {
			num := fmt.Sprintf("%d️⃣", i+1)
			btns := []map[string]string{
				{"text": num + " " + name, "callback_data": "mdl:info:" + name},
			}
			ctrlBtns := make([]map[string]string, 0)
			if i > 0 {
				ctrlBtns = append(ctrlBtns, map[string]string{"text": "⬆️", "callback_data": fmt.Sprintf("mdl:fbup:%s:%d", name, i)})
			}
			if i < len(models.ModelCfg.FallbackModels)-1 {
				ctrlBtns = append(ctrlBtns, map[string]string{"text": "⬇️", "callback_data": fmt.Sprintf("mdl:fbdn:%s:%d", name, i)})
			}
			ctrlBtns = append(ctrlBtns, map[string]string{"text": "❌", "callback_data": fmt.Sprintf("mdl:fbrm:%s:%d", name, i)})
			rows = append(rows, btns)
			rows = append(rows, ctrlBtns)
		}
	}

	rows = append(rows, []map[string]string{
		{"text": "➕ Add Model to Fallback", "callback_data": "mdl:pick:fallback:_"},
	})
	rows = append(rows, []map[string]string{
		{"text": "⬅️ Back", "callback_data": "mdl:menu"},
	})

	return map[string]interface{}{"inline_keyboard": rows}
}

// ApiFormatPickerKeyboard for custom provider
func ApiFormatPickerKeyboard() map[string]interface{} {
	return map[string]interface{}{
		"inline_keyboard": [][]map[string]string{
			{
				{"text": "OpenAI", "callback_data": "mdl:af:openai"},
				{"text": "Anthropic", "callback_data": "mdl:af:anthropic"},
			},
			{
				{"text": "Gemini", "callback_data": "mdl:af:gemini"},
				{"text": "❌ Cancel", "callback_data": "mdl:cancel"},
			},
		},
	}
}
