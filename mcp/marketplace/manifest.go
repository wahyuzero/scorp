package marketplace

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"time"

	"scorp-agent/config"
)

// UpstreamInstall describes how to run the original (Node/Python) server.
type UpstreamInstall struct {
	Command string   `json:"command,omitempty"`
	Args    []string `json:"args,omitempty"`
}

// Upstream captures attribution and runtime info for ported servers.
type Upstream struct {
	Repository       string          `json:"repository,omitempty"`
	Author           string          `json:"author,omitempty"`
	OriginalLanguage string          `json:"original_language,omitempty"`
	License          string          `json:"license,omitempty"`
	PinnedCommit     string          `json:"pinned_commit,omitempty"`
	PinnedVersion    string          `json:"pinned_version,omitempty"`
	Install          *UpstreamInstall `json:"install,omitempty"`
}

// PortInfo attributes the Go port author.
type PortInfo struct {
	Author            string `json:"author,omitempty"`
	GitHubHandle      string `json:"github_handle,omitempty"`
	TranspilerVersion string `json:"transpiler_version,omitempty"`
	PortRepository    string `json:"port_repository,omitempty"`
	License           string `json:"license,omitempty"`
}

// Health is the capability disclosure block.
type Health struct {
	Status            string   `json:"status"` // full | partial
	CoverageScore     float64  `json:"coverage_score"`
	ActiveTools       []string `json:"active_tools"`
	UnsupportedTools  []string `json:"unsupported_tools,omitempty"`
	UnsupportedReason string   `json:"unsupported_reason,omitempty"`
}

// Security carries the CI audit attestations.
type Security struct {
	AstAuditPassed        bool   `json:"ast_audit_passed,omitempty"`
	PromptInjectionScan   bool   `json:"prompt_injection_scanned,omitempty"`
	BuildProvenance       string `json:"build_provenance,omitempty"`
}

// Build describes the toolchain contract.
type Build struct {
	GoVersion string `json:"go_version,omitempty"`
	SDK       string `json:"sdk,omitempty"`
	CGOMade   bool   `json:"cgo,omitempty"`
}

// Variant marks flavors/variants of a canonical port.
type Variant struct {
	IsFlavor       bool     `json:"is_flavor,omitempty"`
	CanonicalBase  string   `json:"canonical_base,omitempty"`
	FlavorName     string   `json:"flavor_name,omitempty"`
	CustomFeatures []string `json:"custom_features,omitempty"`
}

// ToolRef mirrors the tools[] entries of a manifest.
type ToolRef struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Contributor is one attribution changelog entry.
type Contributor struct {
	Handle  string `json:"handle"`
	Role    string `json:"role"`
	Version string `json:"version"`
	Change  string `json:"change,omitempty"`
}

// Drift is the upstream synchronization state maintained by livecheck.
type Drift struct {
	Status               string `json:"status,omitempty"` // synchronized | outdated | unknown
	LastChecked          string `json:"last_checked,omitempty"`
	UpstreamLatestCommit string `json:"upstream_latest_commit,omitempty"`
	UpstreamLatestVer    string `json:"upstream_latest_version,omitempty"`
}

// Manifest is the Go representation of scorp-mcp.json (schema v1).
type Manifest struct {
	Schema       string              `json:"$schema"`
	ID           string              `json:"id"`
	Name         string              `json:"name"`
	Version      string              `json:"version"`
	Description  string              `json:"description"`
	Origin       string              `json:"origin"`
	Upstream     *Upstream           `json:"upstream,omitempty"`
	Port         *PortInfo           `json:"port,omitempty"`
	Health       Health              `json:"health"`
	Security     *Security           `json:"security,omitempty"`
	Build        Build               `json:"build"`
	Variant      *Variant            `json:"variant,omitempty"`
	Artifacts    map[string]Artifact `json:"artifacts,omitempty"`
	Tools        []ToolRef           `json:"tools"`
	Contributors []Contributor       `json:"contributors,omitempty"`
	Drift        *Drift              `json:"drift,omitempty"`
}

var (
	sha256Pattern = regexp.MustCompile(`^[0-9a-f]{64}$`)
	idPattern     = regexp.MustCompile(`^scorp\.mcp\.[a-z0-9][a-z0-9_-]*$`)
)

// Validate enforces the structural invariants of schema v1 that matter at
// install time (the full JSON Schema contract is checked registry-side in CI).
func (m *Manifest) Validate() error {
	if !idPattern.MatchString(m.ID) {
		return fmt.Errorf("invalid manifest id %q", m.ID)
	}
	if m.Name == "" || m.Version == "" {
		return fmt.Errorf("manifest name and version are required")
	}
	switch m.Origin {
	case "port", "native":
	default:
		return fmt.Errorf("origin must be 'port' or 'native', got %q", m.Origin)
	}
	if m.Origin == "port" && m.Upstream == nil {
		return fmt.Errorf("origin=port requires the upstream block")
	}
	if m.Health.Status != "full" && m.Health.Status != "partial" {
		return fmt.Errorf("health.status must be 'full' or 'partial', got %q", m.Health.Status)
	}
	if m.Health.Status == "partial" && len(m.Health.UnsupportedTools) == 0 {
		return fmt.Errorf("health.status=partial requires unsupported_tools")
	}
	if len(m.Tools) == 0 {
		return fmt.Errorf("manifest must declare at least one tool")
	}
	for key, art := range m.Artifacts {
		if !sha256Pattern.MatchString(art.SHA256) {
			return fmt.Errorf("artifact %s: invalid sha256", key)
		}
	}
	return nil
}

// GetManifest resolves and validates the full manifest for an index entry.
// Local file:// manifest URLs and a checkout-local fallback (servers/<name>/)
// are honored so development against a local registry repo works offline.
func GetManifest(ctx context.Context, entry IndexEntry) (*Manifest, error) {
	var data []byte
	var err error

	if entry.ManifestURL != "" {
		data, err = fetchURL(ctx, entry.ManifestURL, "")
	} else {
		err = fmt.Errorf("no manifest url")
	}

	// Local registry checkout fallback: index URLs are the published raw
	// GitHub files; when installing from a local checkout the sibling
	// servers/<name>/manifest.json is authoritative.
	if err != nil && isLocalCheckout() {
		local := localCheckoutDir() + "/" + entry.SourceDir + "/manifest.json"
		if d, lerr := os.ReadFile(local); lerr == nil {
			data, err = d, nil
		}
	}
	if err != nil {
		return nil, fmt.Errorf("fetch manifest for %s: %w", entry.Name, err)
	}

	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parse manifest for %s: %w", entry.Name, err)
	}
	if err := m.Validate(); err != nil {
		return nil, fmt.Errorf("manifest %s failed validation: %w", entry.Name, err)
	}
	return &m, nil
}

// isLocalCheckout reports whether SCORP_MCP_REGISTRY_URL points at a file://
// registry checkout.
func isLocalCheckout() bool {
	return len(IndexURL()) > 7 && IndexURL()[:7] == "file://"
}

func localCheckoutDir() string {
	// SCORP_MCP_REGISTRY_URL=file:///home/user/scorp-mcp-registry/registry.json
	u := IndexURL()
	path := u[len("file://"):]
	if base := tailTrim(path, "registry.json"); base != "" {
		return base
	}
	if base := tailTrim(path, "/"); base != "" {
		return base
	}
	return path
}

func tailTrim(s, suffix string) string {
	if len(s) > len(suffix) && s[len(s)-len(suffix):] == suffix {
		return s[:len(s)-len(suffix)]
	}
	return ""
}

// manifestCachePath returns where validated manifests are persisted after a
// successful install (provenance record under ~/.scorp/marketplace/).
func manifestCachePath(name string) string {
	return config.ScorpPath("marketplace", "servers", name+".manifest.json")
}

// SaveInstalledManifest records the manifest of a successfully installed
// server for later provenance/audit lookups.
func SaveInstalledManifest(m *Manifest) {
	path := manifestCachePath(m.Name)
	_ = os.MkdirAll(config.ScorpPath("marketplace", "servers"), 0o755)
	data, err := json.MarshalIndent(m, "", "  ")
	if err == nil {
		_ = os.WriteFile(path, data, 0o644)
	}
}

// LoadInstalledManifest returns the provenance manifest of an installed server.
func LoadInstalledManifest(name string) (*Manifest, error) {
	data, err := os.ReadFile(manifestCachePath(name))
	if err != nil {
		return nil, err
	}
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

// driftBadge renders the upstream synchronization indicator.
func driftBadge(d *Drift) string {
	if d == nil {
		return "⚪ drift unknown"
	}
	switch d.Status {
	case "synchronized":
		return "🟢 synchronized with upstream"
	case "outdated":
		return fmt.Sprintf("🟡 outdated — upstream advanced (checked %s)", shortDate(d.LastChecked))
	default:
		return "⚪ drift unknown"
	}
}

func shortDate(iso string) string {
	t, err := time.Parse("2006-01-02T15:04:05Z", iso)
	if err != nil {
		return iso
	}
	return t.Format("2006-01-02")
}
