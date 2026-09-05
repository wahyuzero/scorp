package transpiler

import (
	"context"
	"strings"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// goldenFilesystemBinary builds (or reuses) the mcp-filesystem reference port
// from a local scorp-mcp-registry checkout.
func goldenFilesystemBinary(t *testing.T) string {
	t.Helper()
	checkout := os.Getenv("SCORP_MCP_REGISTRY_DIR")
	if checkout == "" {
		defaultCheckout := filepath.Join("..", "..", "..", "scorp-mcp-registry")
		if abs, err := filepath.Abs(defaultCheckout); err == nil {
			if _, statErr := os.Stat(filepath.Join(abs, "servers", "filesystem")); statErr == nil {
				checkout = abs
			}
		}
	}
	if checkout == "" {
		t.Skip("no local scorp-mcp-registry checkout (set SCORP_MCP_REGISTRY_DIR)")
	}

	src := filepath.Join(checkout, "servers", "filesystem")
	bin := filepath.Join(t.TempDir(), "mcp-filesystem")
	buildCtx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(buildCtx, "go", "build", "-o", bin, ".")
	cmd.Dir = src
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build golden port: %v\n%s", err, out)
	}
	return bin
}

// TestProbeAndVerifyGoldenRoundTrip validates the non-LLM pipeline end to end:
// probe the golden port binary → sandbox-build the golden source → verify the
// contract diff is 100%. This is the deterministic stand-in for codegen.
func TestProbeAndVerifyGoldenRoundTrip(t *testing.T) {
	checkout := os.Getenv("SCORP_MCP_REGISTRY_DIR")
	if checkout == "" {
		defaultCheckout := filepath.Join("..", "..", "..", "scorp-mcp-registry")
		if abs, err := filepath.Abs(defaultCheckout); err == nil {
			if _, statErr := os.Stat(filepath.Join(abs, "servers", "filesystem")); statErr == nil {
				checkout = abs
			}
		}
	}
	if checkout == "" {
		t.Skip("no local scorp-mcp-registry checkout (set SCORP_MCP_REGISTRY_DIR)")
	}

	bin := goldenFilesystemBinary(t)

	// Phase 1: capture the contract benchmark from the golden binary.
	bench, err := Probe(context.Background(), ProbeSpec{Name: "filesystem", Command: bin, Args: []string{t.TempDir()}})
	if err != nil {
		t.Fatalf("probe: %v", err)
	}
	if len(bench.Tools) != 11 {
		t.Fatalf("expected 11 tools in benchmark, got %d", len(bench.Tools))
	}

	// Phases 2-3: sandbox-build the golden source (the deterministic
	 // stand-in for LLM codegen output).
	goldenSource, err := os.ReadFile(filepath.Join(checkout, "servers", "filesystem", "main.go"))
	if err != nil {
		t.Fatalf("read golden source: %v", err)
	}
	sandbox, err := NewSandbox("roundtrip", string(goldenSource))
	if err != nil {
		t.Fatal(err)
	}
	defer sandbox.Close()
	if err := sandbox.Build(); err != nil {
		t.Fatalf("sandbox build: %v", err)
	}

	// Phase 4: verify the contract matches 100%.
	report, err := VerifyContract(context.Background(), bench, sandbox.BinPath)
	if err != nil {
		t.Fatalf("verify boot: %v", err)
	}
	if !report.Pass {
		t.Fatalf("contract mismatch: coverage %.2f, missing=%v mismatches=%v",
			report.Coverage, report.Missing, report.ParamMismatches)
	}
}

func TestDetectRuntimeMessages(t *testing.T) {
	err := detectRuntime("definitely-not-a-real-runtime-xyz")
	if err == nil {
		t.Fatal("expected error for missing runtime")
	}
	for _, want := range []string{"not found"} {
		if !contains(err.Error(), want) {
			t.Errorf("error %q missing %q", err.Error(), want)
		}
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (func() bool {
		for i := 0; i+len(sub) <= len(s); i++ {
			if s[i:i+len(sub)] == sub {
				return true
			}
		}
		return false
	})()
}

// TestProbeUpstreamNpx is the live upstream probe (network + npx required).
// Run manually: SCORP_TRANSPILER_E2E=1 go test ./mcp/transpiler/ -run TestProbeUpstreamNpx
func TestProbeUpstreamNpx(t *testing.T) {
	if os.Getenv("SCORP_TRANSPILER_E2E") == "" {
		t.Skip("set SCORP_TRANSPILER_E2E=1 for the live upstream probe")
	}
	bench, err := Probe(context.Background(), ProbeSpec{
		Name:    "filesystem",
		Command: "npx",
		Args:    []string{"-y", "@modelcontextprotocol/server-filesystem", "/tmp"},
		Timeout: 3 * time.Minute,
	})
	if err != nil {
		t.Fatalf("live probe: %v", err)
	}
	if len(bench.Tools) == 0 {
		t.Fatal("no tools captured from live upstream")
	}
	t.Logf("captured %d tools from live npx upstream", len(bench.Tools))
}

func TestStripCodeFences(t *testing.T) {
	simple := "prose\n```go\npackage main\n\nfunc main() {}\n```\n"
	got := stripCodeFences(simple)
	if !strings.Contains(got, "func main()") || strings.Contains(got, "```") {
		t.Fatalf("simple fence not stripped: %q", got)
	}

	// Generated markdown-emitting servers legitimately contain inline
	// triple backticks inside string literals — they must survive.
	withInline := "```go\npackage main\n\nvar s = \"x\\n```\"\n\nfunc main() {}\n```\n"
	got = stripCodeFences(withInline)
	if !strings.Contains(got, "func main()") {
		t.Fatalf("inline backticks broke the stripper: %q", got)
	}

	// Bare Go (no fences) must pass through untouched.
	bare := "package main\n\nfunc main() {}\n"
	if got := stripCodeFences(bare); got != strings.TrimSpace(bare) {
		t.Fatalf("bare source modified: %q", got)
	}
}
