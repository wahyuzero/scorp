package wizard

// ──────────────────────────────────────────────
// Model Manager v2 — Provider-first interactive management
// Features:
//   - /model menu with role badges
//   - Add Provider → auto-populate catalog models from API key
//   - Add Custom Model (manual)
//   - Role assignment: 💬 primary, 🤖 agent, 🎯 delegation
//   - Fallback chain editor (add/remove/reorder)
//   - API key management via .env
// ──────────────────────────────────────────────

import (
	"scorp-agent/config"
	"scorp-agent/models"
	"scorp-agent/tools"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

// CustomProvider type is now in models package
type CustomProvider = models.CustomProvider

// ──────────────────────────────────────────────
// Wizard State Machine
// ──────────────────────────────────────────────

type modelWizard struct {
	Step      string    // current step name
	Provider  string    // selected provider key
	ModelName string    // friendly name (map key)
	ModelID   string    // API model identifier
	APIKey    string    // raw API key (transient)
	BaseURL   string    // custom provider base URL
	APIFormat string    // openai | anthropic | gemini
	IsCustom  bool      // adding custom provider
	Mode      string    // "addprov" | "addmodel" | "setkey"
	Started   time.Time // for timeout
}

var (
	wizards   = make(map[int64]*modelWizard) // keyed by chatID
	wizardsMu sync.Mutex
)

const wizardTimeout = 3 * time.Minute

func GetModelWizard(chatID int64) *modelWizard {
	wizardsMu.Lock()
	defer wizardsMu.Unlock()
	w, ok := wizards[chatID]
	if !ok {
		return nil
	}
	if time.Since(w.Started) > wizardTimeout {
		delete(wizards, chatID)
		return nil
	}
	return w
}

func SetModelWizard(chatID int64, w *modelWizard) {
	wizardsMu.Lock()
	defer wizardsMu.Unlock()
	w.Started = time.Now()
	wizards[chatID] = w
}

func ClearModelWizard(chatID int64) {
	wizardsMu.Lock()
	defer wizardsMu.Unlock()
	delete(wizards, chatID)
}

// ──────────────────────────────────────────────
// .env file management
// ──────────────────────────────────────────────

func EnvFilePath() string {
	return config.ProjectDir() + "/.env"
}

func UpdateEnvFile(key, value string) error {
	os.Setenv(key, value)
	path := EnvFilePath()
	data, _ := os.ReadFile(path)
	lines := strings.Split(string(data), "\n")
	found := false
	keyPrefix := key + "="
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), keyPrefix) {
			lines[i] = keyPrefix + value
			found = true
			break
		}
	}
	var output string
	if found {
		output = strings.Join(lines, "\n")
	} else {
		if len(data) > 0 && !strings.HasSuffix(string(data), "\n") {
			output = string(data) + "\n" + keyPrefix + value + "\n"
		} else {
			output = string(data) + keyPrefix + value + "\n"
		}
	}
	return os.WriteFile(path, []byte(output), 0600)
}

// ──────────────────────────────────────────────
// Provider helpers
// ──────────────────────────────────────────────

func AllProviderNames() []string {
	seen := make(map[string]bool)
	for name := range models.ProviderRegistry {
		seen[name] = true
	}
	models.ModelCfgMu.RLock()
	if models.ModelCfg != nil {
		for name := range models.ModelCfg.CustomProviders {
			seen[name] = true
		}
	}
	models.ModelCfgMu.RUnlock()
	var names []string
	for n := range seen {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

func ProviderDisplayName(name string) string {
	if preset, ok := models.ProviderRegistry[name]; ok && preset.DisplayName != "" {
		return preset.DisplayName
	}
	models.ModelCfgMu.RLock()
	defer models.ModelCfgMu.RUnlock()
	if models.ModelCfg != nil {
		if cp, ok := models.ModelCfg.CustomProviders[name]; ok && cp.DisplayName != "" {
			return cp.DisplayName
		}
	}
	return name
}

// ──────────────────────────────────────────────
// Keyboard Generators
// ──────────────────────────────────────────────

func ModelMenuKeyboard() map[string]interface{} {
	models.ModelCfgMu.RLock()
	defer models.ModelCfgMu.RUnlock()

	var rows [][]map[string]string

	// ── Role cards at top — 1 click → model picker ──
	primaryLabel := "💬 Primary: " + OrNA(models.ModelCfg.DefaultModel)
	agentLabel := "🤖 Agent: " + OrNA(models.ModelCfg.AgentModel)
	// Truncate long labels for button display
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
	if len(delegLabel) > 35 {
		delegLabel = delegLabel[:35] + "…"
	}
	rows = append(rows, []map[string]string{
		{"text": delegLabel, "callback_data": "mdl:pick:role:delegation:_"},
	})

	// ── Model list — tap for details ──
	var modelNames []string
	for name := range models.ModelCfg.Models {
		modelNames = append(modelNames, name)
	}
	sort.Strings(modelNames)

	for _, name := range modelNames {
		m := models.ModelCfg.Models[name]
		// Clean: just name + provider, no emoji stacking
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

// providerListKeyboard — for "Add Provider" flow
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

// modelPickerKeyboard — pick a model for role assignment or fallback
// purpose: "role:chat", "role:agent", "role:delegation", "fallback"
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
			}
		}

		// Clean label: just model name + provider
		label := name + "  (" + m.Provider + ")"

		// Only one indicator: ● = currently selected, nothing = available
		prefix := "  "
		if isCurrent {
			prefix = "● "
		}

		// Determine callback based on purpose
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

	// Back button
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

// modelInfoKeyboard shows actions for a specific model
func ModelInfoKeyboard(name string) map[string]interface{} {
	var rows [][]map[string]string

	// Role assignment row
	rows = append(rows, []map[string]string{
		{"text": "💬 Primary", "callback_data": "mdl:use:" + name},
		{"text": "🤖 Agent", "callback_data": "mdl:ag:" + name},
		{"text": "🎯 Delegation", "callback_data": "mdl:dlg:" + name},
	})

	// Fallback + key
	rows = append(rows, []map[string]string{
		{"text": "🔄 Add to Fallback", "callback_data": "mdl:fba:" + name},
		{"text": "🔑 Set API Key", "callback_data": "mdl:key:" + name},
	})

	// Delete + back
	rows = append(rows, []map[string]string{
		{"text": "🗑️ Delete", "callback_data": "mdl:del:" + name},
		{"text": "⬅️ Back", "callback_data": "mdl:menu"},
	})

	return map[string]interface{}{"inline_keyboard": rows}
}

// fallbackEditorKeyboard — shows fallback chain with reorder controls
func FallbackEditorKeyboard() map[string]interface{} {
	models.ModelCfgMu.RLock()
	defer models.ModelCfgMu.RUnlock()

	var rows [][]map[string]string

	// Show current fallback chain
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
			// Up/down/remove
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

// apiFormatPickerKeyboard for custom provider
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

// ──────────────────────────────────────────────
// Text formatters
// ──────────────────────────────────────────────

func ModelMenuText() string {
	models.ModelCfgMu.RLock()
	defer models.ModelCfgMu.RUnlock()

	var sb strings.Builder
	sb.WriteString("🤖 <b>Model Manager</b>\n\n")

	if len(models.ModelCfg.Models) == 0 {
		sb.WriteString("<i>No models yet. Tap \"Add Provider\" to auto-populate from catalog.</i>")
		return sb.String()
	}

	sb.WriteString("<i>Tap a role to change model, or tap a model for details.</i>\n")
	sb.WriteString(fmt.Sprintf("\n<b>%d Models</b>\n", len(models.ModelCfg.Models)))

	names := make([]string, 0)
	for n := range models.ModelCfg.Models {
		names = append(names, n)
	}
	sort.Strings(names)

	for _, name := range names {
		m := models.ModelCfg.Models[name]
		// Show which roles use this model — clean text, no emoji
		roles := []string{}
		if name == models.ModelCfg.DefaultModel {
			roles = append(roles, "Primary")
		}
		if name == models.ModelCfg.AgentModel && name != models.ModelCfg.DefaultModel {
			roles = append(roles, "Agent")
		}
		if name == models.ModelCfg.DelegationModel && name != models.ModelCfg.DefaultModel && name != models.ModelCfg.AgentModel {
			roles = append(roles, "Delegation")
		}
		roleStr := ""
		if len(roles) > 0 {
			roleStr = " [" + strings.Join(roles, ", ") + "]"
		}
		sb.WriteString(fmt.Sprintf("• <code>%s%s</code> → %s (%s)\n", name, roleStr, m.Model, m.Provider))
	}

	return sb.String()
}

func OrNA(s string) string {
	if s == "" {
		return "(not set)"
	}
	return s
}

func ModelInfoText(name string) string {
	models.ModelCfgMu.RLock()
	defer models.ModelCfgMu.RUnlock()

	m, ok := models.ModelCfg.Models[name]
	if !ok {
		return fmt.Sprintf("❌ Model '%s' not found.", name)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("🔧 <b>%s</b>\n\n", name))
	sb.WriteString(fmt.Sprintf("Model ID: <code>%s</code>\n", m.Model))
	sb.WriteString(fmt.Sprintf("Provider: <code>%s</code>\n", m.Provider))
	sb.WriteString(fmt.Sprintf("Base URL: <code>%s</code>\n", m.BaseURL))
	sb.WriteString(fmt.Sprintf("API: <code>%s</code>\n", m.API))
	sb.WriteString(fmt.Sprintf("Max Tokens: <code>%d</code>\n", m.MaxTokens))

	if m.KeyEnv != "" {
		keyStatus := "❌ not set"
		if os.Getenv(m.KeyEnv) != "" {
			keyStatus = "✅ set"
		}
		sb.WriteString(fmt.Sprintf("Key: <code>%s</code> (%s)\n", m.KeyEnv, keyStatus))
	}

	sb.WriteString("\n")
	if name == models.ModelCfg.DefaultModel {
		sb.WriteString("💬 <b>Primary model</b>\n")
	}
	if name == models.ModelCfg.AgentModel {
		sb.WriteString("🤖 <b>Agent model</b>\n")
	}
	if name == models.ModelCfg.DelegationModel {
		sb.WriteString("🎯 <b>Delegation model</b>\n")
	}

	// Check if in fallback chain
	for i, f := range models.ModelCfg.FallbackModels {
		if f == name {
			sb.WriteString(fmt.Sprintf("🔄 <b>Fallback #%d</b>\n", i+1))
			break
		}
	}

	return sb.String()
}

func FallbackText() string {
	models.ModelCfgMu.RLock()
	defer models.ModelCfgMu.RUnlock()

	var sb strings.Builder
	sb.WriteString("🔄 <b>Fallback Chain</b>\n\n")
	if len(models.ModelCfg.FallbackModels) == 0 {
		sb.WriteString("<i>No fallback models configured.</i>\n\n")
		sb.WriteString("Fallback models are tried automatically when the primary model fails.\n")
		sb.WriteString("Triggers: <code>" + strings.Join(models.ModelCfg.FallbackOnError, ", ") + "</code>")
	} else {
		for i, name := range models.ModelCfg.FallbackModels {
			sb.WriteString(fmt.Sprintf("%d️⃣ <code>%s</code>\n", i+1, name))
		}
	}
	return sb.String()
}

func FormatAPIKeysList() string {
	models.ModelCfgMu.RLock()
	defer models.ModelCfgMu.RUnlock()

	var sb strings.Builder
	sb.WriteString("🔑 <b>API Keys</b>\n\n")

	keyEnvs := make(map[string]string)
	for _, m := range models.ModelCfg.Models {
		if m.KeyEnv != "" {
			if os.Getenv(m.KeyEnv) != "" {
				keyEnvs[m.KeyEnv] = "✅ set"
			} else {
				keyEnvs[m.KeyEnv] = "❌ not set"
			}
		}
	}

	if len(keyEnvs) == 0 {
		sb.WriteString("<i>No API keys configured.</i>\n")
	} else {
		names := make([]string, 0)
		for k := range keyEnvs {
			names = append(names, k)
		}
		sort.Strings(names)
		for _, k := range names {
			sb.WriteString(fmt.Sprintf("<code>%s</code> — %s\n", k, keyEnvs[k]))
		}
	}

	sb.WriteString("\n<i>Tap a model → \"Set API Key\" to update.</i>")
	return sb.String()
}

func formatProvidersList() string {
	var sb strings.Builder
	sb.WriteString("🔌 <b>Providers</b>\n\n")
	sb.WriteString("<b>Built-in:</b>\n")
	names := AllProviderNames()
	builtInCount := 0
	for _, name := range names {
		if preset, ok := models.ProviderRegistry[name]; ok {
			keyStatus := "❌"
			if models.ProviderHasAPIKey(name) {
				keyStatus = "✅"
			}
			catStr := ""
			if models.HasCatalog(name) {
				catStr = fmt.Sprintf(" [%d models]", len(models.CatalogModels(name)))
			}
			sb.WriteString(fmt.Sprintf("%s <code>%s</code> — %s (%s)%s\n", keyStatus, name, preset.DisplayName, preset.API, catStr))
			builtInCount++
		}
	}

	models.ModelCfgMu.RLock()
	customCount := 0
	if models.ModelCfg != nil && len(models.ModelCfg.CustomProviders) > 0 {
		sb.WriteString("\n<b>Custom:</b>\n")
		for name, cp := range models.ModelCfg.CustomProviders {
			disp := cp.DisplayName
			if disp == "" {
				disp = name
			}
			keyStatus := "❌"
			if models.ProviderHasAPIKey(name) {
				keyStatus = "✅"
			}
			sb.WriteString(fmt.Sprintf("%s <code>%s</code> — %s (%s)\n", keyStatus, name, disp, cp.API))
			customCount++
		}
	}
	models.ModelCfgMu.RUnlock()

	sb.WriteString(fmt.Sprintf("\n📊 %d built-in, %d custom", builtInCount, customCount))
	return sb.String()
}

// ──────────────────────────────────────────────
// Fallback chain operations
// ──────────────────────────────────────────────

func addToFallback(name string) string {
	models.ModelCfgMu.Lock()
	defer models.ModelCfgMu.Unlock()
	if _, ok := models.ModelCfg.Models[name]; !ok {
		return fmt.Sprintf("❌ Model '%s' not found.", name)
	}
	for _, f := range models.ModelCfg.FallbackModels {
		if f == name {
			return fmt.Sprintf("ℹ️ <code>%s</code> already in fallback chain.", name)
		}
	}
	models.ModelCfg.FallbackModels = append(models.ModelCfg.FallbackModels, name)
	models.SaveModelConfig()
	return fmt.Sprintf("✅ Added <code>%s</code> to fallback chain.", name)
}

func removeFromFallback(name string) string {
	models.ModelCfgMu.Lock()
	defer models.ModelCfgMu.Unlock()
	var newList []string
	for _, f := range models.ModelCfg.FallbackModels {
		if f != name {
			newList = append(newList, f)
		}
	}
	models.ModelCfg.FallbackModels = newList
	models.SaveModelConfig()
	return fmt.Sprintf("✅ Removed <code>%s</code> from fallback.", name)
}

func moveFallback(name string, direction string) string {
	models.ModelCfgMu.Lock()
	defer models.ModelCfgMu.Unlock()
	idx := -1
	for i, f := range models.ModelCfg.FallbackModels {
		if f == name {
			idx = i
			break
		}
	}
	if idx < 0 {
		return "❌ Not in fallback chain."
	}
	if direction == "up" && idx > 0 {
		models.ModelCfg.FallbackModels[idx], models.ModelCfg.FallbackModels[idx-1] = models.ModelCfg.FallbackModels[idx-1], models.ModelCfg.FallbackModels[idx]
	} else if direction == "down" && idx < len(models.ModelCfg.FallbackModels)-1 {
		models.ModelCfg.FallbackModels[idx], models.ModelCfg.FallbackModels[idx+1] = models.ModelCfg.FallbackModels[idx+1], models.ModelCfg.FallbackModels[idx]
	}
	models.SaveModelConfig()
	return ""
}

// ──────────────────────────────────────────────
// Custom Model Wizard
// ──────────────────────────────────────────────

func startCustomModelWizard(chatID int64) {
	ClearModelWizard(chatID)
	text := "🔧 <b>Add Custom Model</b>\n\nChoose a provider:"
	tools.SendMessage(text, ProviderListKeyboard())
}

// ──────────────────────────────────────────────
// Wizard Text Handler
// ──────────────────────────────────────────────

func handleModelWizardText(text string, chatID int64) bool {
	w := GetModelWizard(chatID)
	if w == nil {
		return false
	}

	if text == "/cancel" || text == "/model cancel" {
		ClearModelWizard(chatID)
		tools.SendMessage("❌ <b>Wizard dibatalkan.</b>", ModelMenuKeyboard())
		return true
	}

	switch w.Step {
	case "prov_name":
		name := strings.TrimSpace(strings.ToLower(text))
		name = strings.ReplaceAll(name, " ", "-")
		if name == "" || name == "custom" {
			tools.SendMessage("❌ Invalid name. Try again:", nil)
			return true
		}
		w.Provider = name
		w.Step = "prov_url"
		tools.SendMessage(fmt.Sprintf("✅ Provider: <code>%s</code>\n\nEnter base URL (e.g., <code>https://api.example.com/v1</code>):", name), nil)
		return true

	case "prov_url":
		url := strings.TrimSpace(text)
		if !strings.HasPrefix(url, "http") {
			tools.SendMessage("❌ URL must start with http:// or https://. Try again:", nil)
			return true
		}
		w.BaseURL = url
		w.Step = "prov_api"
		tools.SendMessage(fmt.Sprintf("✅ Base URL: <code>%s</code>\n\nChoose API format:", url), ApiFormatPickerKeyboard())
		return true

	case "model_name":
		parts := strings.Fields(text)
		if len(parts) == 0 {
			tools.SendMessage("❌ Empty input. Enter name (and optionally model ID):", nil)
			return true
		}
		w.ModelName = strings.ToLower(parts[0])
		if len(parts) > 1 {
			w.ModelID = strings.Join(parts[1:], " ")
		}
		w.Step = "model_id"
		if w.ModelID != "" {
			return askAPIKey(w, chatID)
		}
		tools.SendMessage(fmt.Sprintf("✅ Name: <code>%s</code>\n\nEnter model ID (e.g., <code>gpt-4o</code>):", w.ModelName), nil)
		return true

	case "model_id":
		id := strings.TrimSpace(text)
		if id == "" {
			tools.SendMessage("❌ Empty model ID. Try again:", nil)
			return true
		}
		w.ModelID = id
		return askAPIKey(w, chatID)

	case "api_key":
		key := strings.TrimSpace(text)
		if key == "/skip" || key == "skip" {
			if w.Mode == "addprov" {
				return finalizeProviderKeySave(w, chatID)
			}
			if w.ModelName != "" && w.Mode == "setkey" {
				return finalizeModelKeySave(w, chatID)
			}
			return showWizardSummary(w, chatID)
		}
		if len(key) < 10 {
			tools.SendMessage("❌ API key looks too short. Try again or type /skip:", nil)
			return true
		}
		w.APIKey = key

		if w.Mode == "addprov" {
			return finalizeProviderKeySave(w, chatID)
		}
		if w.Mode == "setkey" {
			return finalizeModelKeySave(w, chatID)
		}
		return showWizardSummary(w, chatID)
	}

	return false
}

// askAPIKey transitions to API key input step
func askAPIKey(w *modelWizard, chatID int64) bool {
	preset, hasPreset := models.ProviderRegistry[w.Provider]
	noAuth := false
	if hasPreset {
		noAuth = preset.NoAuth
	}

	if noAuth {
		w.APIKey = ""
		return showWizardSummary(w, chatID)
	}

	keyHint := models.ProviderKeyEnv(w.Provider)

	if keyHint != "" && os.Getenv(keyHint) != "" {
		w.Step = "api_key"
		tools.SendMessage(fmt.Sprintf("✅ Model ID: <code>%s</code>\n\n🔑 Key <code>%s</code> already set.\nSend new key or type /skip:", w.ModelID, keyHint), nil)
		return true
	}

	w.Step = "api_key"
	prompt := fmt.Sprintf("✅ Model ID: <code>%s</code>\n\n", w.ModelID)
	if keyHint != "" {
		prompt += fmt.Sprintf("🔑 Enter API key (saved as <code>%s</code>):\n", keyHint)
	} else {
		prompt += "🔑 Enter API key:\n"
	}
	prompt += "\n<i>⚠️ Your message will NOT be auto-deleted. Delete it manually after.</i>"
	tools.SendMessage(prompt, nil)
	return true
}

// showWizardSummary for custom model add
func showWizardSummary(w *modelWizard, chatID int64) bool {
	w.Step = "confirm"
	var sb strings.Builder
	sb.WriteString("📋 <b>Review Model</b>\n\n")
	sb.WriteString(fmt.Sprintf("🔤 Name: <code>%s</code>\n", w.ModelName))
	sb.WriteString(fmt.Sprintf("🤖 Model ID: <code>%s</code>\n", w.ModelID))
	sb.WriteString(fmt.Sprintf("🔌 Provider: <code>%s</code>\n", w.Provider))
	if w.IsCustom {
		sb.WriteString(fmt.Sprintf("🌐 Base URL: <code>%s</code>\n", w.BaseURL))
		sb.WriteString(fmt.Sprintf("📡 API Format: <code>%s</code>\n", w.APIFormat))
	}
	keyEnv := models.ProviderKeyEnv(w.Provider)
	if w.APIKey != "" {
		sb.WriteString(fmt.Sprintf("🔑 API Key: <code>%s</code> = %s...%s\n", keyEnv, w.APIKey[:4], w.APIKey[len(w.APIKey)-4:]))
	} else if keyEnv != "" && os.Getenv(keyEnv) != "" {
		sb.WriteString(fmt.Sprintf("🔑 API Key: <code>%s</code> (existing)\n", keyEnv))
	} else {
		sb.WriteString("🔑 API Key: ⚠️ <i>not set</i>\n")
	}
	sb.WriteString("\nSave this model?")
	tools.SendMessage(sb.String(), map[string]interface{}{
		"inline_keyboard": [][]map[string]string{
			{
				{"text": "✅ Save", "callback_data": "mdl:save:" + w.ModelName},
				{"text": "❌ Cancel", "callback_data": "mdl:cancel"},
			},
		},
	})
	return true
}

// finalizeProviderKeySave — after API key entered in "addprov" mode
func finalizeProviderKeySave(w *modelWizard, chatID int64) bool {
	ClearModelWizard(chatID)

	if w.APIKey != "" {
		keyEnv := models.ProviderKeyEnv(w.Provider)
		UpdateEnvFile(keyEnv, w.APIKey)
	}

	// Auto-populate from catalog
	if models.HasCatalog(w.Provider) {
		added := models.AutoPopulateFromCatalog(w.Provider, "")
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("✅ <b>%s connected!</b>\n\n", ProviderDisplayName(w.Provider)))
		sb.WriteString(fmt.Sprintf("📦 %d models auto-added:\n", len(added)))
		for _, name := range added {
			sb.WriteString(fmt.Sprintf("  • <code>%s</code>\n", name))
		}
		sb.WriteString("\nTap a model to set it as Primary/Agent/Delegation.")
		tools.SendMessage(sb.String(), ModelMenuKeyboard())
	} else {
		tools.SendMessage(fmt.Sprintf("✅ <b>%s</b> API key disimpan!\n\nProvider ini tidak punya model catalog. Gunakan \"Tambah Model\" untuk menambah", ProviderDisplayName(w.Provider)), ModelMenuKeyboard())
	}
	return true
}

// finalizeModelKeySave — after API key entered in "setkey" mode
func finalizeModelKeySave(w *modelWizard, chatID int64) bool {
	ClearModelWizard(chatID)
	keyEnv := models.ProviderKeyEnv(w.Provider)
	if w.APIKey != "" {
		UpdateEnvFile(keyEnv, w.APIKey)
		msg := fmt.Sprintf("✅ API key updated for <code>%s</code>\nSaved as: <code>%s</code>", w.ModelName, keyEnv)
		tools.SendMessage(msg, ModelMenuKeyboard())
	}
	return true
}

// finalizeWizardSave — save custom model from wizard confirmation
func finalizeWizardSave(w *modelWizard, chatID int64) string {
	models.ModelCfgMu.Lock()
	defer models.ModelCfgMu.Unlock()

	// Save custom provider if needed
	if w.IsCustom {
		if models.ModelCfg.CustomProviders == nil {
			models.ModelCfg.CustomProviders = make(map[string]CustomProvider)
		}
		keyEnv := models.ProviderKeyEnv(w.Provider)
		cp := CustomProvider{
			BaseURL:     w.BaseURL,
			API:         w.APIFormat,
			DisplayName: w.Provider,
		}
		if keyEnv != "" {
			cp.KeyEnvs = []string{keyEnv}
		}
		models.ModelCfg.CustomProviders[w.Provider] = cp
		models.ProviderRegistry[w.Provider] = models.ProviderPreset{
			BaseURL:     w.BaseURL,
			API:         w.APIFormat,
			KeyEnvs:     cp.KeyEnvs,
			DisplayName: w.Provider,
		}
	}

	apiFormat := "openai"
	if w.IsCustom && w.APIFormat != "" {
		apiFormat = w.APIFormat
	} else if preset, ok := models.ProviderRegistry[w.Provider]; ok && preset.API != "" {
		apiFormat = preset.API
	}

	baseURL := w.BaseURL
	if baseURL == "" {
		if preset, ok := models.ProviderRegistry[w.Provider]; ok {
			baseURL = preset.BaseURL
		}
	}

	keyEnv := models.ProviderKeyEnv(w.Provider)
	if w.APIKey != "" && keyEnv != "" {
		go UpdateEnvFile(keyEnv, w.APIKey)
	}

	mc := models.ModelConfig{
		Provider:  w.Provider,
		Model:     w.ModelID,
		BaseURL:   baseURL,
		MaxTokens: 4096,
		API:       apiFormat,
	}
	if keyEnv != "" {
		mc.KeyEnv = keyEnv
	}
	models.ModelCfg.Models[w.ModelName] = mc

	if models.ModelCfg.DefaultModel == "" {
		models.ModelCfg.DefaultModel = w.ModelName
		if models.ModelCfg.RoutingRules == nil {
			models.ModelCfg.RoutingRules = make(map[string]string)
		}
		models.ModelCfg.RoutingRules["chat"] = w.ModelName
	}
	if models.ModelCfg.AgentModel == "" {
		models.ModelCfg.AgentModel = w.ModelName
	}

	models.SaveModelConfig()
	return fmt.Sprintf("✅ <b>Model saved!</b>\n\n<code>%s</code> → %s (%s)\nKey: <code>%s</code>",
		w.ModelName, w.ModelID, w.Provider, keyEnv)
}

// deleteModel removes a model from config
func deleteModel(name string) string {
	models.ModelCfgMu.Lock()
	defer models.ModelCfgMu.Unlock()
	if _, ok := models.ModelCfg.Models[name]; !ok {
		return fmt.Sprintf("❌ Model '%s' not found.", name)
	}
	delete(models.ModelCfg.Models, name)
	if models.ModelCfg.DefaultModel == name {
		models.ModelCfg.DefaultModel = ""
		delete(models.ModelCfg.RoutingRules, "chat")
	}
	if models.ModelCfg.AgentModel == name {
		models.ModelCfg.AgentModel = ""
		delete(models.ModelCfg.RoutingRules, "agent")
	}
	if models.ModelCfg.DelegationModel == name {
		models.ModelCfg.DelegationModel = ""
	}
	if models.ModelCfg.PremiumModel == name {
		models.ModelCfg.PremiumModel = ""
	}
	var newFallback []string
	for _, f := range models.ModelCfg.FallbackModels {
		if f != name {
			newFallback = append(newFallback, f)
		}
	}
	models.ModelCfg.FallbackModels = newFallback
	models.SaveModelConfig()
	return fmt.Sprintf("✅ Model <code>%s</code> deleted.", name)
}

// ──────────────────────────────────────────────
// Callback Handler
// ──────────────────────────────────────────────

func HandleModelCallback(data string, chatID int64, msgID int64) (string, map[string]interface{}, bool) {
	parts := strings.SplitN(data, ":", 4)
	if len(parts) < 2 || parts[0] != "mdl" {
		return "", nil, false
	}

	action := parts[1]
	arg1 := ""
	arg2 := ""
	if len(parts) > 2 {
		arg1 = parts[2]
	}
	if len(parts) > 3 {
		arg2 = parts[3]
	}
	// Strip trailing placeholder from callbacks like mdl:pick:role:agent:_
	arg1 = strings.TrimSuffix(arg1, ":_")
	arg2 = strings.TrimSuffix(arg2, ":_")

	switch action {

	// ── Main menu ──
	case "menu":
		return ModelMenuText(), ModelMenuKeyboard(), true

	case "close":
		return "✅ <b>Model menu closed.</b>", nil, true

	// ── Model info ──
	case "info":
		return ModelInfoText(arg1), ModelInfoKeyboard(arg1), true

	// ── Set roles ──
	case "use":
		msg := models.SwitchModel("chat", arg1)
		return msg, ModelMenuKeyboard(), true

	case "ag":
		msg := models.SwitchModel("agent", arg1)
		return msg, ModelMenuKeyboard(), true

	case "dlg":
		msg := models.SwitchModel("delegation", arg1)
		return msg, ModelMenuKeyboard(), true

	// ── Delete model ──
	case "del":
		msg := deleteModel(arg1)
		return msg, ModelMenuKeyboard(), true

	// ── Add Provider (show provider list) ──
	case "aprov":
		text := "🔌 <b>Add Provider</b>\n\n"
		text += "✅ = key set | ❌ = no key | 📦 = has model catalog\n\n"
		text += "Pick a provider to add/replace its API key:"
		return text, ProviderListKeyboard(), true

	// ── Provider picked → ask for API key ──
	case "pk":
		w := &modelWizard{
			Step:     "api_key",
			Provider: arg1,
			Mode:     "addprov",
		}
		SetModelWizard(chatID, w)

		keyEnv := models.ProviderKeyEnv(arg1)
		keyStatus := "❌ not set"
		if os.Getenv(keyEnv) != "" {
			keyStatus = "✅ already set"
		}

		text := fmt.Sprintf("🔌 <b>%s</b>\n\n", ProviderDisplayName(arg1))
		text += fmt.Sprintf("Key: <code>%s</code> (%s)\n", keyEnv, keyStatus)
		if models.HasCatalog(arg1) {
			text += fmt.Sprintf("📦 %d models will be auto-added after key is set.\n\n", len(models.CatalogModels(arg1)))
		}
		text += "Enter API key (or /skip to keep existing):"
		return text, nil, true

	// ── New custom provider ──
	case "newprov":
		w := &modelWizard{
			Step:     "prov_name",
			IsCustom: true,
			Mode:     "addprov",
		}
		SetModelWizard(chatID, w)
		return "🔧 <b>New Custom Provider</b>\n\nEnter provider name (lowercase, no spaces):", nil, true

	// ── Add Custom Model ──
	case "acustom":
		startCustomModelWizard(chatID)
		return "", nil, true

	// ── API format picked (custom provider wizard) ──
	case "af":
		w := GetModelWizard(chatID)
		if w == nil {
			return "❌ Wizard expired. Start again.", ModelMenuKeyboard(), true
		}
		w.APIFormat = arg1
		w.Step = "model_name"
		return fmt.Sprintf("✅ API format: <code>%s</code>\n\nEnter name and model ID.\nExample: <code>myclaude claude-3-5-sonnet</code>", arg1), nil, true

	// ── Save wizard ──
	case "save":
		w := GetModelWizard(chatID)
		if w == nil {
			return "❌ Wizard expired.", ModelMenuKeyboard(), true
		}
		msg := finalizeWizardSave(w, chatID)
		ClearModelWizard(chatID)
		return msg, ModelMenuKeyboard(), true

	// ── Cancel ──
	case "cancel":
		ClearModelWizard(chatID)
		return "❌ <b>Cancelled.</b>", ModelMenuKeyboard(), true

	// ── Roles menu (redirect to main — role cards are on main menu now) ──
	case "roles":
		return ModelMenuText(), ModelMenuKeyboard(), true

	// ── Show model picker (role assignment) ──
	// mdl:pick:role:chat:_ → picker for chat role
	case "pick":
		if arg1 == "role" {
			roleType := arg2
			current := "—"
			models.ModelCfgMu.RLock()
			switch roleType {
			case "chat":
				current = OrNA(models.ModelCfg.DefaultModel)
			case "agent":
				current = OrNA(models.ModelCfg.AgentModel)
			case "delegation":
				current = OrNA(models.ModelCfg.DelegationModel)
			}
			models.ModelCfgMu.RUnlock()
			title := fmt.Sprintf("%s — pick a model\nCurrent: <code>%s</code>", roleLabel(roleType), current)
			return title, ModelPickerKeyboard("role:" + roleType), true
		}
		if arg1 == "fallback" {
			return "Pick model to add to fallback chain:", ModelPickerKeyboard("fallback"), true
		}
		return "", nil, false

	// ── Role assignment: model picked from picker ──
	// mdl:rdo:chat:modelName → models.SwitchModel("chat", modelName)
	case "rdo":
		roleType := arg1
		modelName := arg2
		_ = models.SwitchModel(roleType, modelName)
		// Return to main model menu with updated role cards
		return ModelMenuText(), ModelMenuKeyboard(), true

	// ── Fallback editor ──
	case "fb":
		return FallbackText(), FallbackEditorKeyboard(), true

	// ── Fallback add (from model info) ──
	case "fba":
		msg := addToFallback(arg1)
		return msg, FallbackEditorKeyboard(), true

	// ── Fallback remove ──
	case "fbrm":
		// arg1 = name:idx
		fbParts := strings.SplitN(arg1, ":", 2)
		if len(fbParts) == 2 {
			msg := removeFromFallback(fbParts[0])
			return msg, FallbackEditorKeyboard(), true
		}
		return "", nil, false

	// ── Fallback move up ──
	case "fbup":
		fbParts := strings.SplitN(arg1, ":", 2)
		if len(fbParts) == 2 {
			moveFallback(fbParts[0], "up")
			return FallbackText(), FallbackEditorKeyboard(), true
		}
		return "", nil, false

	// ── Fallback move down ──
	case "fbdn":
		fbParts := strings.SplitN(arg1, ":", 2)
		if len(fbParts) == 2 {
			moveFallback(fbParts[0], "down")
			return FallbackText(), FallbackEditorKeyboard(), true
		}
		return "", nil, false

	// ── Health check ──
	case "check":
		tools.SendMessage("⏳ <i>Checking all models...</i>", nil)
		go func() {
			result := models.FormatModelListWithHealth()
			tools.SendMessage(result, ModelMenuKeyboard())
		}()
		return "", nil, true

	// ── API key management ──
	case "key":
		models.ModelCfgMu.RLock()
		mc, exists := models.ModelCfg.Models[arg1]
		models.ModelCfgMu.RUnlock()
		if !exists {
			return fmt.Sprintf("❌ Model '%s' not found.", arg1), ModelMenuKeyboard(), true
		}
		w := &modelWizard{
			Step:      "api_key",
			Provider:  mc.Provider,
			ModelName: arg1,
			ModelID:   mc.Model,
			Mode:      "setkey",
		}
		SetModelWizard(chatID, w)
		keyEnv := mc.KeyEnv
		if keyEnv == "" {
			keyEnv = models.ProviderKeyEnv(mc.Provider)
		}
		prompt := fmt.Sprintf("🔑 <b>Set API Key for %s</b>\n\n", arg1)
		if keyEnv != "" {
			if os.Getenv(keyEnv) != "" {
				prompt += fmt.Sprintf("Current: <code>%s</code> = ✅ set\n", keyEnv)
			} else {
				prompt += fmt.Sprintf("Will save as: <code>%s</code>\n", keyEnv)
			}
		}
		prompt += "\nEnter new API key:"
		return prompt, nil, true

	// ── API keys list ──
	case "keys":
		return FormatAPIKeysList(), ModelMenuKeyboard(), true

	// ── Providers list ──
	case "prov":
		return formatProvidersList(), ModelMenuKeyboard(), true
	}

	return "", nil, false
}

func roleLabel(role string) string {
	switch role {
	case "chat":
		return "💬 Primary"
	case "agent":
		return "🤖 Agent"
	case "delegation":
		return "🎯 Delegation"
	}
	return role
}

// ──────────────────────────────────────────────
// Wizard text router — called from handleAction default case
// ──────────────────────────────────────────────

func HandleModelWizardTextRouter(text string, chatID int64) bool {
	w := GetModelWizard(chatID)
	if w == nil {
		return false
	}
	return handleModelWizardText(text, chatID)
}
