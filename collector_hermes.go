package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"time"
)

// HermesData holds complete Hermes Agent state.
type HermesData struct {
	Model         string
	Provider      string
	BaseURL       string
	GatewayStatus string
	IsPeak        bool

	// Agent settings
	MaxTurns       int
	ReasoningEffort string
	Plugins        []string

	// Auxiliary models
	AuxVision      string
	AuxCompress    string
	AuxWebExtract  string
	AuxTitle       string

	// OMH routing (role → model)
	OMHRoles       map[string]string

	// MCP servers
	MCPServers     []string

	// Model aliases
	Aliases        map[string]string

	// Fallback providers
	Fallbacks      []string
}

// Snapshot string for change detection
func (d HermesData) modelSnapshot() string {
	return d.Model + "|" + d.Provider
}

var (
	lastHermesSnapshot string
	hermesConfigPath   string
	omhConfigPath      string

	// Regexes for main config
	rxModelDefault   = regexp.MustCompile(`(?m)^\s{2}default:\s*(.+)`)
	rxModelProvider  = regexp.MustCompile(`(?m)^\s{2}provider:\s*(.+)`)
	rxModelBaseURL   = regexp.MustCompile(`(?m)^\s{2}base_url:\s*(.+)`)
	rxMaxTurns       = regexp.MustCompile(`(?m)max_turns:\s*(\d+)`)
	rxReasoning      = regexp.MustCompile(`(?m)reasoning_effort:\s*(\S+)`)
)

func init() {
	hermesConfigPath = os.Getenv("HERMES_CONFIG_PATH")
	if hermesConfigPath == "" {
		hermesConfigPath = hermesConfigPath_()
	}
	omhConfigPath = os.Getenv("OMH_CONFIG_PATH")
	if omhConfigPath == "" {
		omhConfigPath = omhConfigPath_()
	}
}

// ──────────────────────────────────────────────
// Collectors
// ──────────────────────────────────────────────

func collectHermes() HermesData {
	var d HermesData

	// Read main config
	raw, err := os.ReadFile(hermesConfigPath)
	if err != nil {
		log.Printf("[hermes] config read error: %v", err)
		d.Model = "unknown"
		d.GatewayStatus = getGatewayStatus()
		d.IsPeak = isPeakHour()
		return d
	}
	cfg := string(raw)

	if m := rxModelDefault.FindStringSubmatch(cfg); len(m) > 1 {
		d.Model = strings.TrimSpace(m[1])
	}
	if m := rxModelProvider.FindStringSubmatch(cfg); len(m) > 1 {
		d.Provider = strings.TrimSpace(m[1])
	}
	if m := rxModelBaseURL.FindStringSubmatch(cfg); len(m) > 1 {
		d.BaseURL = strings.TrimSpace(m[1])
	}
	if m := rxMaxTurns.FindStringSubmatch(cfg); len(m) > 1 {
		fmt.Sscanf(m[1], "%d", &d.MaxTurns)
	}
	if m := rxReasoning.FindStringSubmatch(cfg); len(m) > 1 {
		d.ReasoningEffort = strings.TrimSpace(m[1])
	}

	// Auxiliary models
	d.AuxVision = extractAuxModel(cfg, "vision")
	d.AuxCompress = extractAuxModel(cfg, "compression")
	d.AuxWebExtract = extractAuxModel(cfg, "web_extract")
	d.AuxTitle = extractAuxModel(cfg, "title_generation")

	// Plugins
	d.Plugins = extractYAMLList(cfg, "plugins", "enabled")

	// MCP servers
	d.MCPServers = extractMCPServers(cfg)

	// Model aliases
	d.Aliases = extractAliases(cfg)

	// Fallback providers
	d.Fallbacks = extractFallbacks(cfg)

	// OMH routing
	d.OMHRoles = collectOMHRoles()

	// Runtime
	d.GatewayStatus = getGatewayStatus()
	d.IsPeak = isPeakHour()

	return d
}

// extractAuxModel finds auxiliary.MODEL.model: value
func extractAuxModel(cfg, section string) string {
	// Pattern: section:\n    ...\n    model: VALUE
	re := regexp.MustCompile(`(?m)^\s{2}` + section + `:\s*\n(?:\s{4}\w.*\n)*?\s{4}model:\s*(.+)`)
	if m := re.FindStringSubmatch(cfg); len(m) > 1 {
		v := strings.TrimSpace(m[1])
		if v == "''" || v == "\"\"" {
			return "auto"
		}
		return v
	}
	return "auto"
}

// extractYAMLList finds a list under parent/key
func extractYAMLList(cfg, parent, key string) []string {
	// Find "parent:" then "  key:" then list items until un-indent
	re := regexp.MustCompile(`(?m)^\s{2}` + parent + `:\s*\n(?:^\s{4}.*\n)*?^\s{4}` + key + `:\s*\n((?:^\s{6}-\s.*\n)*)`)
	m := re.FindStringSubmatch(cfg)
	if len(m) < 2 {
		return nil
	}
	itemRe := regexp.MustCompile(`(?m)^\s*-\s*(.+)`)
	var items []string
	for _, match := range itemRe.FindAllStringSubmatch(m[1], -1) {
		items = append(items, strings.TrimSpace(match[1]))
	}
	return items
}

// extractMCPServers finds mcp_servers keys — each top-level key at exactly 2-space indent
func extractMCPServers(cfg string) []string {
	// Find the mcp_servers: section and extract only top-level keys (2-space indent, word chars + colon at EOL)
	re := regexp.MustCompile(`(?ms)^mcp_servers:\n(.*?)^[a-z_]`)
	m := re.FindStringSubmatch(cfg)
	if len(m) < 2 {
		return nil
	}
	// Only match keys at exactly 2-space indent that end with ':' and nothing else on the line
	keyRe := regexp.MustCompile(`(?m)^  ([a-zA-Z][\w-]*):$`)
	var servers []string
	for _, match := range keyRe.FindAllStringSubmatch(m[1], -1) {
		servers = append(servers, match[1])
	}
	return servers
}

// extractAliases finds model_aliases section
func extractAliases(cfg string) map[string]string {
	re := regexp.MustCompile(`(?ms)^model_aliases:\s*\n((?:^\s{2}\w.*\n)*)`)
	m := re.FindStringSubmatch(cfg)
	if len(m) < 2 {
		return nil
	}
	entryRe := regexp.MustCompile(`(?m)^\s{2}(\w[\w-]*):\s*(.+)`)
	aliases := make(map[string]string)
	for _, match := range entryRe.FindAllStringSubmatch(m[1], -1) {
		aliases[match[1]] = strings.TrimSpace(match[2])
	}
	return aliases
}

// extractFallbacks finds fallback_providers list
func extractFallbacks(cfg string) []string {
	re := regexp.MustCompile(`(?ms)^fallback_providers:\s*\n((?:^\s{2}-.*\n(?:^\s{4}.*\n)*)*)`)
	m := re.FindStringSubmatch(cfg)
	if len(m) < 2 {
		return nil
	}
	modelRe := regexp.MustCompile(`(?m)^\s{4}-\s*model:\s*(.+)`)
	provRe := regexp.MustCompile(`(?m)^\s{4}provider:\s*(.+)`)
	models := modelRe.FindAllStringSubmatch(m[1], -1)
	provs := provRe.FindAllStringSubmatch(m[1], -1)
	var fallbacks []string
	for i, mm := range models {
		fb := strings.TrimSpace(mm[1])
		if i < len(provs) {
			fb += " (" + strings.TrimSpace(provs[i][1]) + ")"
		}
		fallbacks = append(fallbacks, fb)
	}
	return fallbacks
}

// collectOMHRoles reads OMH plugin config for role→model mapping
func collectOMHRoles() map[string]string {
	raw, err := os.ReadFile(omhConfigPath)
	if err != nil {
		return nil
	}
	cfg := string(raw)
	roles := make(map[string]string)

	// Pattern: roleName: {provider: "...", model: "..."}
	re := regexp.MustCompile(`(?m)^\s{2}([\w-]+):\s*\{[^}]*model:\s*"([^"]+)"`)
	for _, m := range re.FindAllStringSubmatch(cfg, -1) {
		roles[m[1]] = m[2]
	}
	return roles
}

// ──────────────────────────────────────────────
// Helpers
// ──────────────────────────────────────────────

func getGatewayStatus() string {
	out, err := exec.Command("systemctl", "-M", "ubuntu@", "is-active", "hermes-gateway").Output()
	if err == nil {
		s := strings.TrimSpace(string(out))
		if s == "active" || s == "activating" || s == "deactivating" {
			return s
		}
	}
	out, err = exec.Command("pgrep", "-f", "hermes.*gateway").Output()
	if err == nil && strings.TrimSpace(string(out)) != "" {
		return "active"
	}
	return "inactive"
}

func isPeakHour() bool {
	utcHour := time.Now().UTC().Hour()
	return utcHour >= 6 && utcHour < 10
}

func quotaMultiplier(model string, isPeak bool) string {
	switch model {
	case "glm-5.2", "glm-5-turbo":
		if isPeak {
			return "3×"
		}
		return "1× (promo)"
	case "glm-5.1", "glm-4.7", "glm-4.5-air", "glm-4.6", "glm-4.5":
		return "1×"
	default:
		return "—"
	}
}

func modelInfo(model string) string {
	switch model {
	case "glm-5.2":
		return "Flagship • 1M ctx • 131K output • 296 tok/s"
	case "glm-5.1":
		return "Solid all-rounder • 1× always"
	case "glm-4.7":
		return "Fast routine tasks • 1× always"
	case "glm-4.5-air":
		return "Lightweight • 1× always"
	default:
		return ""
	}
}

// ──────────────────────────────────────────────
// Formatters
// ──────────────────────────────────────────────

// formatModelChangeAlert — compact notification for model switch
func formatModelChangeAlert(d HermesData, oldModel string) string {
	period := "Off-peak"
	periodIcon := "🌙"
	if d.IsPeak {
		period = "PEAK (13-17 WIB)"
		periodIcon = "🔥"
	}

	mult := quotaMultiplier(d.Model, d.IsPeak)
	info := modelInfo(d.Model)
	if info != "" {
		info = "\n💡 " + info
	}

	gwIcon := "✅"
	if d.GatewayStatus != "active" {
		gwIcon = "⚠️"
	}

	now := time.Now().UTC().Format("15:04:05 UTC")

	// Count OMH roles and aux models
	omhCount := len(d.OMHRoles)

	return fmt.Sprintf("🔄 <b>HERMES MODEL SWITCHED</b>\n"+
		"🕐 %s • %s %s\n"+
		"━━━━━━━━━━━━━━━━━\n"+
		"Previous: <code>%s</code>\n"+
		"Current:  <code>%s</code> ⚡ <b>%s</b>\n"+
		"━━━━━━━━━━━━━━━━━\n"+
		"🔌 Gateway: %s <code>%s</code>\n"+
		"🏷️ Provider: <code>%s</code>\n"+
		"🧩 Plugins: <code>%s</code>\n"+
		"🤖 OMH roles: <b>%d</b> | MCP: <b>%d</b> | Aliases: <b>%d</b>%s",
		now, periodIcon, period,
		oldModel,
		d.Model, mult,
		gwIcon, d.GatewayStatus,
		d.Provider,
		strings.Join(d.Plugins, ", "),
		omhCount, len(d.MCPServers), len(d.Aliases),
		info,
	)
}

// formatHermesStatus — full status report (for /hermes command)
func formatHermesStatus(d HermesData) string {
	period := "🌙 Off-peak"
	if d.IsPeak {
		period = "🔥 PEAK (13-17 WIB)"
	}

	mult := quotaMultiplier(d.Model, d.IsPeak)
	gwIcon := "✅"
	if d.GatewayStatus != "active" {
		gwIcon = "⚠️"
	}

	now := time.Now().UTC().Format("02 Jan 2006, 15:04 UTC")

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("🤖 <b>Hermes Agent Status</b>\n🕐 %s • %s\n%s\n\n",
		now, period, strings.Repeat("─", 28)))

	// Main model
	sb.WriteString(fmt.Sprintf("<b>📌 Main Model</b>\n"))
	sb.WriteString(fmt.Sprintf("  Model: <code>%s</code> ⚡ <b>%s</b>\n", d.Model, mult))
	sb.WriteString(fmt.Sprintf("  Provider: <code>%s</code>\n", d.Provider))
	if info := modelInfo(d.Model); info != "" {
		sb.WriteString(fmt.Sprintf("  💡 %s\n", info))
	}

	// Agent settings
	sb.WriteString(fmt.Sprintf("\n<b>⚙️ Agent</b>\n"))
	sb.WriteString(fmt.Sprintf("  Max turns: <b>%d</b> | Reasoning: <b>%s</b>\n", d.MaxTurns, d.ReasoningEffort))
	sb.WriteString(fmt.Sprintf("  Gateway: %s <code>%s</code>\n", gwIcon, d.GatewayStatus))
	sb.WriteString(fmt.Sprintf("  Plugins: <code>%s</code>\n", strings.Join(d.Plugins, ", ")))

	// OMH routing
	if len(d.OMHRoles) > 0 {
		sb.WriteString(fmt.Sprintf("\n<b>🧩 OMH Role Routing (%d roles)</b>\n", len(d.OMHRoles)))
		// Sort roles for consistent display
		roles := make([]string, 0, len(d.OMHRoles))
		for r := range d.OMHRoles {
			roles = append(roles, r)
		}
		sort.Strings(roles)
		for _, r := range roles {
			sb.WriteString(fmt.Sprintf("  <code>%-18s</code> → %s\n", r, d.OMHRoles[r]))
		}
	}

	// Auxiliary models
	sb.WriteString(fmt.Sprintf("\n<b>🔧 Auxiliary Models</b>\n"))
	sb.WriteString(fmt.Sprintf("  👁️ Vision:     <code>%s</code>\n", d.AuxVision))
	sb.WriteString(fmt.Sprintf("  🗜️ Compress:   <code>%s</code>\n", d.AuxCompress))
	sb.WriteString(fmt.Sprintf("  🌐 Web Extract: <code>%s</code>\n", d.AuxWebExtract))
	sb.WriteString(fmt.Sprintf("  📝 Title Gen:  <code>%s</code>\n", d.AuxTitle))

	// MCP servers
	if len(d.MCPServers) > 0 {
		sb.WriteString(fmt.Sprintf("\n<b>🔌 MCP Servers (%d)</b>\n", len(d.MCPServers)))
		sb.WriteString(fmt.Sprintf("  <code>%s</code>\n", strings.Join(d.MCPServers, ", ")))
	}

	// Model aliases
	if len(d.Aliases) > 0 {
		sb.WriteString(fmt.Sprintf("\n<b>🔖 Aliases (%d)</b>\n", len(d.Aliases)))
		// Show as compact list
		keys := make([]string, 0, len(d.Aliases))
		for k := range d.Aliases {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			sb.WriteString(fmt.Sprintf("  <code>/%-10s</code> → %s\n", k, d.Aliases[k]))
		}
	}

	// Fallbacks
	if len(d.Fallbacks) > 0 {
		sb.WriteString(fmt.Sprintf("\n<b>🔄 Fallback Chain</b>\n"))
		for i, f := range d.Fallbacks {
			sb.WriteString(fmt.Sprintf("  %d. <code>%s</code>\n", i+1, f))
		}
	}

	sb.WriteString("\n" + strings.Repeat("─", 28) + "\n")
	sb.WriteString("🔥 13:00 WIB → glm-5.1 (peak)\n")
	sb.WriteString("🌙 17:00 WIB → glm-5.2 (off-peak promo)")

	return sb.String()
}
