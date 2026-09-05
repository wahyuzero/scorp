// Package marketplace implements the Scorp MCP Marketplace client:
// registry index synchronization, manifest parsing/validation, search with
// health/drift disclosure, and the Tri-Option installation flow
// (prebuilt binary / local rebuild / upstream runtime).
//
// Architectural note: every install path converges on the same terminal state
// as mcp_manage(action="add") — an entry in ~/.scorp/mcp.json. The existing
// infrastructure (watchdog, hot-reload, native tool registration) takes over
// from there automatically.
package marketplace

import "time"

// Registry location and cache policy.
const (
	// DefaultIndexURL is the canonical registry index (override with
	// SCORP_MCP_REGISTRY_URL, which may also be a file:// path for local
	// development against a checkout of scorp-mcp-registry).
	DefaultIndexURL = "https://raw.githubusercontent.com/wahyuzero/scorp-mcp-registry/main/registry.json"

	// CacheTTL bounds how long a cached index stays fresh.
	CacheTTL = 1 * time.Hour

	// HTTPTimeout bounds registry/manifest/artifact downloads.
	HTTPTimeout = 60 * time.Second
)

// InstallOption is one of the Tri-Option choices presented during install.
type InstallOption int

const (
	OptionUnknown InstallOption = iota
	OptionPrebuilt           // ⚡ download CI-built binary, SHA-256 pinned
	OptionRebuild            // 🛠 probe upstream → transpile → go build locally
	OptionUpstream           // 📦 register the original npx/uvx runtime
)

func (o InstallOption) String() string {
	switch o {
	case OptionPrebuilt:
		return "prebuilt"
	case OptionRebuild:
		return "rebuild"
	case OptionUpstream:
		return "upstream"
	default:
		return "unknown"
	}
}

// ParseInstallOption maps a user selection ("1"/"2"/"3") to an InstallOption.
func ParseInstallOption(s string) InstallOption {
	switch s {
	case "1":
		return OptionPrebuilt
	case "2":
		return OptionRebuild
	case "3":
		return OptionUpstream
	default:
		return OptionUnknown
	}
}

// artifactKey returns the manifest artifacts map key for the running platform.
func artifactKey(goos, goarch string) string { return goos + "_" + goarch }
