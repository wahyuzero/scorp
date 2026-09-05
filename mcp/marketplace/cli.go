package marketplace

import (
	"context"
	"fmt"
	"strings"
)

// CLISearch renders the marketplace search view for the terminal.
func CLISearch(term string) (string, error) {
	entries, err := Search(context.Background(), term)
	if err != nil {
		return "", err
	}
	return FormatSearchResults(entries, term), nil
}

// CLIInstall runs the install flow. When opt is empty it renders the
// disclosure card + Tri-Option menu and returns needsSelection=true so the
// caller can prompt for a choice interactively; otherwise it executes
// immediately. The reserved target "upstream" installs an unlisted server
// straight from a spec string (blueprint Case B).
func CLIInstall(name, opt string) (out string, ok bool, needsSelection bool) {
	ctx := context.Background()

	// Case B: unlisted upstream runtime from a raw spec.
	if name == "upstream" {
		out, ok = InstallUpstreamSpec(opt)
		return out, ok, false
	}

	entry, err := FindEntry(ctx, name)
	if err != nil {
		// Helpful Case B guidance when the target is unknown to the registry.
		return fmt.Sprintf("💡 No existing Go port found for %q in the Scorp Marketplace.\n\n"+
			"Install the upstream runtime directly:\n"+
			"  /mcp install upstream npx:@scope/some-mcp-server\n"+
			"  /mcp install upstream uvx:some-mcp-server\n\n"+
			"Or search the marketplace: /mcp search <term>", name), false, false
	}

	m, err := GetManifest(ctx, *entry)
	if err != nil {
		return fmt.Sprintf("❌ %v", err), false, false
	}

	if opt == "" {
		return RenderTriOptionMenu(entry, m), true, true
	}

	choice := ParseInstallOption(opt)
	if choice == OptionUnknown {
		return fmt.Sprintf("❌ Invalid selection %q. Choose 1, 2, or 3.", opt), false, false
	}

	out, ok = Install(ctx, entry, choice)
	return out, ok, false
}

// RenderTriOptionMenu builds the blueprint §2 Scenario A dialog.
func RenderTriOptionMenu(entry *IndexEntry, m *Manifest) string {
	var sb strings.Builder
	sb.WriteString(DisclosureCard(entry, m))
	sb.WriteString("\nSelect Installation Method:\n")

	prebuilt := "1. ⚡ [Prebuilt Marketplace Binary]\n" +
		"      Download official precompiled Go binary (Instant startup, verified SHA-256)."
	if _, hasArtifact := m.Artifacts[ArtifactKeyForPlatform()]; !hasArtifact {
		prebuilt += " [unavailable for this platform]"
	}
	sb.WriteString(prebuilt + "\n")

	sb.WriteString("2. 🛠️ [Local Machine Rebuild from Source]\n" +
		"      Probe upstream, transpile to Go, inspect, and compile locally (Zero-Trust / Audit-Ready).\n")
	sb.WriteString("3. 📦 [Upstream Original Runtime]\n" +
		"      Execute original implementation via " + upstreamRuntime(m) + " (Full original parity, higher RAM).\n")
	sb.WriteString("\nSelection [1/2/3]: ")
	return sb.String()
}

func upstreamRuntime(m *Manifest) string {
	if m.Upstream != nil && strings.EqualFold(m.Upstream.OriginalLanguage, "python") {
		return "uvx"
	}
	return "npx"
}
