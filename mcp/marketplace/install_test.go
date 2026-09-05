package marketplace

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// tempArtifact writes a payload file and returns its file:// URL + sha256.
func tempArtifact(t *testing.T, payload []byte) (url, sha string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "artifact")
	if err := os.WriteFile(path, payload, 0o644); err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(payload)
	return "file://" + path, hex.EncodeToString(sum[:])
}

func TestDownloadVerifiedSuccess(t *testing.T) {
	isolate(t)
	payload := []byte("#!/bin/sh\necho fake-mcp-server\n")
	url, sha := tempArtifact(t, payload)

	dest := filepath.Join(t.TempDir(), "server")
	if err := downloadVerified(context.Background(), url, sha, dest); err != nil {
		t.Fatalf("downloadVerified: %v", err)
	}

	data, err := os.ReadFile(dest)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != string(payload) {
		t.Fatal("artifact content drifted")
	}
	info, _ := os.Stat(dest)
	if info.Mode()&0o100 == 0 {
		t.Fatal("installed artifact is not executable")
	}
}

func TestDownloadVerifiedRejectsTamperedArtifact(t *testing.T) {
	isolate(t)
	url, realSHA := tempArtifact(t, []byte("legit payload"))

	// Simulate tampering-in-transit: the pin describes different bytes than
	// what the URL now serves.
	wrongSum := sha256.Sum256([]byte("malicious payload"))
	pinnedWrong := hex.EncodeToString(wrongSum[:])
	if pinnedWrong == realSHA {
		t.Fatal("test setup produced identical hashes")
	}

	dest := filepath.Join(t.TempDir(), "server")
	err := downloadVerified(context.Background(), url, pinnedWrong, dest)
	if err == nil || !strings.Contains(err.Error(), "MISMATCH") {
		t.Fatalf("expected SHA-256 mismatch error, got: %v", err)
	}
	if _, statErr := os.Stat(dest); !os.IsNotExist(statErr) {
		t.Fatal("unverified artifact must never land at the destination")
	}
}

func TestInstallUpstreamSpec(t *testing.T) {
	isolate(t)

	out, ok := InstallUpstreamSpec("npx:@modelcontextprotocol/server-filesystem /tmp")
	if !ok {
		t.Fatalf("npx spec install failed: %s", out)
	}

	// Verify the entry landed in the (isolated) mcp.json.
	data, err := os.ReadFile(filepath.Join(os.Getenv("HOME"), ".scorp", "mcp.json"))
	if err != nil {
		t.Fatalf("mcp.json not written: %v", err)
	}
	if !strings.Contains(string(data), "npx-") || !strings.Contains(string(data), "server-filesystem") || !strings.Contains(string(data), `"/tmp"`) {
		t.Fatalf("mcp.json missing upstream entry:\n%s", data)
	}

	if out, ok := InstallUpstreamSpec("garbage"); ok {
		t.Fatalf("garbage spec should fail, got: %s", out)
	}
}

func TestInstallOption2ReportsComingSoon(t *testing.T) {
	isolate(t)
	indexURL, _ := writeFixtureRegistry(t)
	t.Setenv("SCORP_MCP_REGISTRY_URL", indexURL)

	idx, err := GetIndex(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	entry := idx.Servers[0]

	RebuildHook = nil
	out, ok := Install(context.Background(), &entry, OptionRebuild)
	if ok {
		t.Fatal("rebuild should not succeed without a hook")
	}
	if !strings.Contains(out, "coming in v2.5") {
		t.Fatalf("unexpected rebuild message: %s", out)
	}

	// Wire the hook (as the M2 transpiler will) and confirm delegation.
	RebuildHook = func(ctx context.Context, entry *IndexEntry, m *Manifest) (string, bool) {
		return "HOOK-OK " + entry.Name, true
	}
	defer func() { RebuildHook = nil }()
	out, ok = Install(context.Background(), &entry, OptionRebuild)
	if !ok || !strings.Contains(out, "HOOK-OK fetch") {
		t.Fatalf("rebuild hook not delegated: %s (ok=%v)", out, ok)
	}
}
