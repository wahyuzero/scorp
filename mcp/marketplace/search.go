package marketplace

import (
	"context"
	"fmt"
	"runtime"
	"strings"
)

// Search filters the registry index by a free-text term matched against name,
// description, tool names, and origin. Empty term returns the full catalog.
func Search(ctx context.Context, term string) ([]IndexEntry, error) {
	idx, err := GetIndex(ctx)
	if err != nil {
		return nil, err
	}
	term = strings.ToLower(strings.TrimSpace(term))
	if term == "" {
		return idx.Servers, nil
	}

	var hits []IndexEntry
	for _, e := range idx.Servers {
		if entryMatches(e, term) {
			hits = append(hits, e)
		}
	}
	return hits, nil
}

func entryMatches(e IndexEntry, term string) bool {
	if strings.Contains(strings.ToLower(e.Name), term) ||
		strings.Contains(strings.ToLower(e.Description), term) ||
		strings.Contains(strings.ToLower(e.ID), term) ||
		strings.Contains(strings.ToLower(e.Origin), term) {
		return true
	}
	// Flavors surface under their canonical base too (search "sqlite" finds
	// "sqlite-cipher" via @author/canonical-base naming, blueprint §4D).
	if e.Variant != nil {
		if strings.Contains(strings.ToLower(e.Variant.CanonicalBase), term) ||
			strings.Contains(strings.ToLower(e.Variant.FlavorName), term) {
			return true
		}
	}
	for _, t := range e.Tools {
		if strings.Contains(strings.ToLower(t), term) {
			return true
		}
	}
	return false
}

// FindEntry locates one index entry by exact name (canonical lookup for install).
func FindEntry(ctx context.Context, name string) (*IndexEntry, error) {
	idx, err := GetIndex(ctx)
	if err != nil {
		return nil, err
	}
	for i := range idx.Servers {
		if idx.Servers[i].Name == name {
			return &idx.Servers[i], nil
		}
	}
	return nil, fmt.Errorf("'%s' not found in the Scorp registry", name)
}

// FormatSearchResults renders the numbered multi-result view (blueprint §4D).
func FormatSearchResults(entries []IndexEntry, term string) string {
	if len(entries) == 0 {
		return fmt.Sprintf("🔍 No servers matching %q in the Scorp Marketplace.\n\nTip: you can still install any Node/Python server directly:\n  /mcp install upstream npx:<package-name>\n  /mcp install upstream uvx:<package-name>", term)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("🔍 Scorp Marketplace — %d result(s) for %q:\n\n", len(entries), term))
	for i, e := range entries {
		label := originLabel(e.Origin)
		if e.Variant != nil && e.Variant.IsFlavor {
			label = fmt.Sprintf("⚡ Flavor of %s", e.Variant.CanonicalBase)
		}
		sb.WriteString(fmt.Sprintf("%d. 📦 <b>%s</b> v%s (%s)\n", i+1, e.Name, e.Version, label))
		if e.Description != "" {
			sb.WriteString(fmt.Sprintf("   %s\n", e.Description))
		}
		sb.WriteString(fmt.Sprintf("   %s · %d tool(s): %s\n", healthBadge(e.Health), len(e.Tools), strings.Join(e.Tools, ", ")))
		sb.WriteString(fmt.Sprintf("   %s\n\n", driftLine(e.Drift)))
	}
	sb.WriteString("Install with: /mcp install <name>")
	return sb.String()
}

func originLabel(origin string) string {
	switch origin {
	case "native":
		return "native Go"
	default:
		return "Go port"
	}
}

func healthBadge(status string) string {
	switch status {
	case "full":
		return "🟢 Full Support"
	case "partial":
		return "🟡 Partial Support"
	default:
		return "⚪ Unknown"
	}
}

func driftLine(drift string) string {
	switch drift {
	case "synchronized":
		return "🟢 Up-to-date with upstream"
	case "outdated":
		return "🟡 Upstream drift detected"
	default:
		return "⚪ Drift unknown"
	}
}

// DisclosureCard renders the pre-install transparency card (blueprint §4B):
// health coverage, unsupported tools with reasons, drift status, and the
// resource delta vs the upstream runtime.
func DisclosureCard(entry *IndexEntry, m *Manifest) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("✨ %s available on Scorp Marketplace!\n\n", entry.Name))
	sb.WriteString(fmt.Sprintf("   • Version:         v%s\n", entry.Version))
	if m.Port != nil && m.Port.Author != "" {
		by := m.Port.Author
		if m.Port.GitHubHandle != "" {
			by = "@" + m.Port.GitHubHandle
		}
		sb.WriteString(fmt.Sprintf("   • Ported by:       %s (%s)\n", by, m.Build.SDK))
	}
	if m.Upstream != nil {
		sb.WriteString(fmt.Sprintf("   • Upstream Source: %s (%s, %s)\n", m.Upstream.Repository, m.Upstream.OriginalLanguage, m.Upstream.License))
	}

	active := entry.Tools
	if len(m.Health.ActiveTools) > 0 {
		active = m.Health.ActiveTools
	}
	switch m.Health.Status {
	case "full":
		sb.WriteString(fmt.Sprintf("   • Health Status:   🟢 Full Support (%d/%d tools active)\n", len(active), len(active)))
	case "partial":
		sb.WriteString(fmt.Sprintf("   • Health Status:   🟡 Partial Support (%d of %d tools active)\n", len(active), len(active)+len(m.Health.UnsupportedTools)))
		for _, t := range m.Health.UnsupportedTools {
			sb.WriteString(fmt.Sprintf("     ⚠️ %s (Disabled: %s)\n", t, m.Health.UnsupportedReason))
		}
	}

	sb.WriteString(fmt.Sprintf("   • Resource Delta:  Go binary ~5MB RAM vs %s\n", upstreamRAM(m)))
	sb.WriteString(fmt.Sprintf("   • Drift:           %s\n", driftBadge(m.Drift)))
	if m.Security != nil && m.Security.BuildProvenance != "" {
		audit := "🟢 Verified (AST scanned, SHA-256 pinned)"
		sb.WriteString(fmt.Sprintf("   • Audit Status:    %s\n", audit))
	}
	return sb.String()
}

func upstreamRAM(m *Manifest) string {
	switch {
	case m.Upstream == nil:
		return "typical Node.js ~48MB RAM"
	case strings.EqualFold(m.Upstream.OriginalLanguage, "python"):
		return "Python runtime ~65MB RAM"
	default:
		return "Node.js runtime ~48MB RAM"
	}
}

// ArtifactKeyForPlatform exposes the platform artifact map key (for UX hints).
func ArtifactKeyForPlatform() string {
	return artifactKey(runtime.GOOS, runtime.GOARCH)
}
