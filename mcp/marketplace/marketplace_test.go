package marketplace

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeFixtureRegistry creates a minimal scorp-mcp-registry checkout on disk
// and returns its registry.json file:// URL plus the checkout root.
func writeFixtureRegistry(t *testing.T) (indexURL string, root string) {
	t.Helper()
	root = t.TempDir()

	manifest := Manifest{
		Schema:      "https://raw.githubusercontent.com/wahyuzero/scorp-mcp-registry/main/schema/v1.json",
		ID:          "scorp.mcp.fetch",
		Name:        "fetch",
		Version:     "1.0.0",
		Description: "fixture web fetch server",
		Origin:      "port",
		Upstream: &Upstream{
			Repository:       "https://github.com/modelcontextprotocol/servers/tree/main/src/fetch",
			Author:           "Anthropic",
			OriginalLanguage: "TypeScript",
			License:          "MIT",
			Install:          &UpstreamInstall{Command: "npx", Args: []string{"-y", "@modelcontextprotocol/server-fetch"}},
		},
		Port: &PortInfo{Author: "Wahyu", GitHubHandle: "wahyuzero"},
		Health: Health{
			Status:        "partial",
			CoverageScore: 0.5,
			ActiveTools:   []string{"fetch"},
			UnsupportedTools: []string{
				"render_pdf_screenshot",
			},
			UnsupportedReason: "Headless Chromium not ported",
		},
		Build: Build{GoVersion: "1.24+", SDK: "github.com/mark3labs/mcp-go"},
		Tools: []ToolRef{{Name: "fetch", Description: "fetches a URL"}},
	}

	manifestData, _ := json.MarshalIndent(manifest, "", "  ")
	serverDir := filepath.Join(root, "servers", "fetch")
	if err := os.MkdirAll(serverDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(serverDir, "manifest.json"), manifestData, 0o644); err != nil {
		t.Fatal(err)
	}

	index := RegistryIndex{
		SchemaVersion: 1,
		Servers: []IndexEntry{
			{
				ID:          "scorp.mcp.fetch",
				Name:        "fetch",
				Version:     "1.0.0",
				Description: "fixture web fetch server",
				Origin:      "port",
				SourceDir:   "servers/fetch",
				Health:      "partial",
				Tools:       []string{"fetch"},
				Drift:       "synchronized",
				ManifestURL: "https://raw.githubusercontent.com/wahyuzero/scorp-mcp-registry/main/servers/fetch/manifest.json",
			},
			{
				ID:          "scorp.mcp.sqlite",
				Name:        "sqlite",
				Version:     "1.0.0",
				Description: "fixture sqlite server",
				Origin:      "port",
				SourceDir:   "servers/sqlite",
				Health:      "full",
				Tools:       []string{"read_query"},
				Drift:       "unknown",
			},
		},
	}
	indexData, _ := json.MarshalIndent(index, "", "  ")
	if err := os.WriteFile(filepath.Join(root, "registry.json"), indexData, 0o644); err != nil {
		t.Fatal(err)
	}

	return "file://" + filepath.Join(root, "registry.json"), root
}

// isolate redirects ~/.scorp state into a temp dir and points the registry at
// the fixture checkout.
func isolate(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("SCORP_MCP_REGISTRY_URL", "")
	return home
}

func TestGetIndexLocalFile(t *testing.T) {
	isolate(t)
	indexURL, _ := writeFixtureRegistry(t)
	t.Setenv("SCORP_MCP_REGISTRY_URL", indexURL)

	idx, err := GetIndex(context.Background())
	if err != nil {
		t.Fatalf("GetIndex: %v", err)
	}
	if len(idx.Servers) != 2 {
		t.Fatalf("expected 2 servers, got %d", len(idx.Servers))
	}
	if idx.Servers[0].Name != "fetch" {
		t.Fatalf("unexpected first server: %s", idx.Servers[0].Name)
	}
}

func TestGetIndexStaleCacheFallback(t *testing.T) {
	isolate(t)
	indexURL, _ := writeFixtureRegistry(t)
	t.Setenv("SCORP_MCP_REGISTRY_URL", indexURL)

	ctx := context.Background()
	if _, err := GetIndex(ctx); err != nil {
		t.Fatalf("prime cache: %v", err)
	}

	// Break the source: stale cache must still serve.
	broken := strings.ReplaceAll(indexURL, "registry.json", "missing.json")
	t.Setenv("SCORP_MCP_REGISTRY_URL", broken)
	idx, err := GetIndex(ctx)
	if err != nil {
		t.Fatalf("stale fallback: %v", err)
	}
	if len(idx.Servers) != 2 {
		t.Fatalf("stale cache should still return 2 servers, got %d", len(idx.Servers))
	}
}

func TestSearchFiltersByToolAndName(t *testing.T) {
	isolate(t)
	indexURL, _ := writeFixtureRegistry(t)
	t.Setenv("SCORP_MCP_REGISTRY_URL", indexURL)

	hits, err := Search(context.Background(), "read_query")
	if err != nil {
		t.Fatal(err)
	}
	if len(hits) != 1 || hits[0].Name != "sqlite" {
		t.Fatalf("tool search mismatch: %+v", hits)
	}

	hits, err = Search(context.Background(), "FETCH")
	if err != nil {
		t.Fatal(err)
	}
	if len(hits) != 1 || hits[0].Name != "fetch" {
		t.Fatalf("name search mismatch: %+v", hits)
	}
}

func TestFormatSearchResultsBadges(t *testing.T) {
	entries := []IndexEntry{
		{Name: "fetch", Version: "1.0.0", Origin: "port", Health: "partial", Tools: []string{"fetch"}, Drift: "outdated"},
	}
	out := FormatSearchResults(entries, "fetch")
	for _, want := range []string{"🟡 Partial Support", "🟡 Upstream drift detected", "Go port"} {
		if !strings.Contains(out, want) {
			t.Errorf("search output missing %q:\n%s", want, out)
		}
	}
}

func TestManifestValidation(t *testing.T) {
	valid := &Manifest{ID: "scorp.mcp.fetch", Name: "fetch", Version: "1.0.0", Origin: "native", Health: Health{Status: "full"}, Tools: []ToolRef{{Name: "fetch"}}}
	if err := valid.Validate(); err != nil {
		t.Errorf("valid manifest rejected: %v", err)
	}

	noUpstream := &Manifest{ID: "scorp.mcp.x", Name: "x", Version: "1.0.0", Origin: "port", Health: Health{Status: "full"}, Tools: []ToolRef{{Name: "t"}}}
	if err := noUpstream.Validate(); err == nil {
		t.Error("port without upstream should fail validation")
	}

	badSHA := &Manifest{ID: "scorp.mcp.x", Name: "x", Version: "1.0.0", Origin: "native", Health: Health{Status: "full"}, Tools: []ToolRef{{Name: "t"}},
		Artifacts: map[string]Artifact{"linux_amd64": {URL: "https://x", SHA256: "deadbeef"}}}
	if err := badSHA.Validate(); err == nil {
		t.Error("malformed sha256 should fail validation")
	}
}

func TestGetManifestLocalCheckoutFallback(t *testing.T) {
	isolate(t)
	indexURL, root := writeFixtureRegistry(t)
	t.Setenv("SCORP_MCP_REGISTRY_URL", indexURL)

	idx, err := GetIndex(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	entry := idx.Servers[0] // manifest_url points at raw.githubusercontent (unreachable in test)

	m, err := GetManifest(context.Background(), entry)
	if err != nil {
		t.Fatalf("local checkout fallback: %v", err)
	}
	if m.Name != "fetch" {
		t.Fatalf("unexpected manifest: %+v", m)
	}
	if m.Health.Status != "partial" || len(m.Health.UnsupportedTools) == 0 {
		t.Fatalf("health block not parsed: %+v", m.Health)
	}
	if m.Upstream == nil || m.Upstream.Install == nil || m.Upstream.Install.Command != "npx" {
		t.Fatalf("upstream install block not parsed: %+v", m.Upstream)
	}
	_ = root
}

func TestDisclosureCardPartialHealth(t *testing.T) {
	isolate(t)
	indexURL, _ := writeFixtureRegistry(t)
	t.Setenv("SCORP_MCP_REGISTRY_URL", indexURL)

	idx, _ := GetIndex(context.Background())
	entry := idx.Servers[0]
	m, err := GetManifest(context.Background(), entry)
	if err != nil {
		t.Fatal(err)
	}

	card := DisclosureCard(&entry, m)
	for _, want := range []string{"🟡 Partial Support", "render_pdf_screenshot", "Headless Chromium not ported", "@wahyuzero"} {
		if !strings.Contains(card, want) {
			t.Errorf("disclosure card missing %q:\n%s", want, card)
		}
	}
}

func TestParseInstallOption(t *testing.T) {
	cases := map[string]InstallOption{"1": OptionPrebuilt, "2": OptionRebuild, "3": OptionUpstream, "x": OptionUnknown, "": OptionUnknown}
	for in, want := range cases {
		if got := ParseInstallOption(in); got != want {
			t.Errorf("ParseInstallOption(%q) = %v, want %v", in, got, want)
		}
	}
}
