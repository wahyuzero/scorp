package marketplace

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"scorp-agent/mcp"
)

// RebuildHook is wired by the transpiler package (M2). When nil, Option 2
// reports transparently that local rebuild is not yet available.
var RebuildHook func(ctx context.Context, entry *IndexEntry, m *Manifest) (string, bool)

// binariesDir stores prebuilt marketplace binaries (Layer 4 target path).
func binariesDir() string { return filepath.Join(os.Getenv("HOME"), ".scorp", "mcp-binaries") }

// Install executes one Tri-Option installation path for a registry entry and
// returns the human-readable dialog output plus success.
func Install(ctx context.Context, entry *IndexEntry, opt InstallOption) (string, bool) {
	m, err := GetManifest(ctx, *entry)
	if err != nil {
		return fmt.Sprintf("❌ %v", err), false
	}

	switch opt {
	case OptionPrebuilt:
		return installPrebuilt(ctx, entry, m)
	case OptionRebuild:
		if RebuildHook != nil {
			return RebuildHook(ctx, entry, m)
		}
		return "🛠️ Local rebuild (AI Transpiler) is coming in v2.5.\n" +
			"For now, choose option 1 (prebuilt binary) or 3 (upstream runtime).", false
	case OptionUpstream:
		return installUpstream(entry, m)
	default:
		return "❌ Unknown installation option", false
	}
}

// installPrebuilt downloads the CI-built artifact, verifies its SHA-256 pin
// BEFORE execution (Layer 4), installs it into ~/.scorp/mcp-binaries/, and
// registers it via the shared mcp.json entry path.
func installPrebuilt(ctx context.Context, entry *IndexEntry, m *Manifest) (string, bool) {
	key := artifactKey(runtime.GOOS, runtime.GOARCH)
	art, ok := m.Artifacts[key]
	if !ok {
		available := make([]string, 0, len(m.Artifacts))
		for k := range m.Artifacts {
			available = append(available, k)
		}
		return fmt.Sprintf("❌ No prebuilt artifact for %s/%s (available: %s).\n"+
			"Try option 2 (local rebuild, needs Go toolchain) or option 3 (upstream runtime).",
			runtime.GOOS, runtime.GOARCH, strings.Join(available, ", ")), false
	}

	binPath := filepath.Join(binariesDir(), entry.Name)
	fmt.Fprintf(os.Stderr, "⬇️  Downloading %s (%s/%s)...\n", entry.Name, runtime.GOOS, runtime.GOARCH)
	if err := downloadVerified(ctx, art.URL, art.SHA256, binPath); err != nil {
		_ = os.Remove(binPath) // never keep an unverified artifact
		return fmt.Sprintf("❌ Prebuilt install failed: %v", err), false
	}

	out, ok := mcp.AddServerEntry(entry.Name, mcp.MCPServerConfig{Command: binPath})
	if !ok {
		return out, false
	}
	SaveInstalledManifest(m)
	return out, true
}

// installUpstream registers the original Node/Python runtime from the
// manifest's upstream.install block.
func installUpstream(entry *IndexEntry, m *Manifest) (string, bool) {
	if m.Upstream == nil || m.Upstream.Install == nil || m.Upstream.Install.Command == "" {
		return fmt.Sprintf("❌ No upstream runtime registered for '%s'.\n"+
			"Register one manually: mcp_manage(action=\"add\", ...) or /mcp install upstream npx:<package>", entry.Name), false
	}
	out, ok := mcp.AddServerEntry(entry.Name, mcp.MCPServerConfig{
		Command: m.Upstream.Install.Command,
		Args:    m.Upstream.Install.Args,
	})
	if ok {
		SaveInstalledManifest(m)
	}
	return out, ok
}

// specPattern matches npx:<pkg> / uvx:<pkg> / bare executable paths.
var specPattern = regexp.MustCompile(`^(npx|uvx|bunx|python3?|node):(.+)$`)

// InstallUpstreamSpec registers an unlisted upstream server directly from a
// spec string — Case B of the blueprint decision flow:
//
//	npx:@modelcontextprotocol/server-filesystem /tmp
//	uvx:mcp-server-fetch
func InstallUpstreamSpec(spec string) (string, bool) {
	spec = strings.TrimSpace(spec)
	m := specPattern.FindStringSubmatch(spec)
	if m == nil {
		return "❌ Unrecognized upstream spec.\n" +
			"Use: <runtime>:<package> — e.g. npx:@modelcontextprotocol/server-filesystem or uvx:mcp-server-fetch", false
	}
	command, pkg := m[1], m[2]

	// Everything after the package name becomes runtime arguments
	// (npx:@org/server /tmp → npx -y @org/server /tmp).
	parts := strings.Fields(pkg)
	if len(parts) == 0 {
		return "❌ Empty package spec.", false
	}

	var args []string
	switch command {
	case "npx", "bunx":
		args = append([]string{"-y"}, parts...)
	default:
		args = parts
	}

	name := deriveSpecName(command, parts[0])
	return mcp.AddServerEntry(name, mcp.MCPServerConfig{Command: command, Args: args})
}

// deriveSpecName produces a stable config key from a spec
// (npx:@org/server-x -> org-server-x).
func deriveSpecName(command, pkg string) string {
	clean := strings.NewReplacer("@", "", "/", "-", ":", "-").Replace(pkg)
	if len(clean) > 40 {
		clean = clean[:40]
	}
	return command + "-" + strings.Trim(clean, "-")
}

// downloadVerified streams url to dest and enforces the SHA-256 pin before the
// artifact becomes executable (Layer 4 — abort on a single divergent bit).
func downloadVerified(ctx context.Context, url, wantSHA, dest string) error {
	if !sha256Pattern.MatchString(wantSHA) {
		return fmt.Errorf("registry artifact has malformed sha256 pin")
	}

	var reader io.Reader
	if strings.HasPrefix(url, "file://") {
		f, err := os.Open(strings.TrimPrefix(url, "file://"))
		if err != nil {
			return err
		}
		defer f.Close()
		reader = f
	} else {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return err
		}
		req.Header.Set("User-Agent", "scorp-agent-marketplace")

		resp, err := httpClient.Do(req)
		if err != nil {
			return fmt.Errorf("download: %w", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("download: HTTP %s", resp.Status)
		}
		reader = resp.Body
	}

	_ = os.MkdirAll(binariesDir(), 0o755)
	tmp, err := os.CreateTemp(binariesDir(), ".download-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName) // no-op after successful rename

	hasher := sha256.New()
	size, err := io.Copy(io.MultiWriter(tmp, hasher), io.LimitReader(reader, 128<<20))
	if err := tmp.Close(); err != nil {
		return err
	}
	if err != nil {
		return fmt.Errorf("download: %w", err)
	}

	got := hex.EncodeToString(hasher.Sum(nil))
	if got != wantSHA {
		return fmt.Errorf("SHA-256 MISMATCH — expected %s, got %s (artifact %d bytes). Install aborted.",
			wantSHA[:16]+"…", got[:16]+"…", size)
	}

	if err := os.Chmod(tmpName, 0o755); err != nil {
		return err
	}
	if err := os.Rename(tmpName, dest); err != nil {
		return err
	}
	log.Printf("[marketplace] installed %s (%d bytes, sha256 verified)", dest, size)
	return nil
}
