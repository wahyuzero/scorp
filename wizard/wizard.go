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
	"fmt"
	"os"
	"scorp-agent/config"
	"scorp-agent/models"
	"scorp-agent/tools"
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
		tools.SendMessage("❌ <b>Wizard cancelled.</b>", ModelMenuKeyboard())
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
		tools.SendMessage(fmt.Sprintf("✅ <b>%s</b> API key saved!\n\nThis provider has no model catalog. Use \"Add Model\" to add one", ProviderDisplayName(w.Provider)), ModelMenuKeyboard())
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
			case "vision":
				current = OrNA(models.ModelCfg.VisionModel)
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
	case "vision":
		return "👁 Vision"
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
