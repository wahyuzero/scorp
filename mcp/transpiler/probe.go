// Package transpiler implements the Scorp AI Transpiler: contract-driven
// transpilation of upstream Node.js/Python MCP servers into idiomatic Go
// binaries built on mark3labs/mcp-go.
//
// Pipeline (blueprint §5A):
//  1. Probe    — spawn the upstream server ephemerally, capture tools/list
//                schemas as the contract benchmark.
//  2. Generate — LLM codegen (routing: complex/premium) emitting a single
//                decoupled main.go against the mcp-go SDK.
//  3. Build    — sandboxed `go build` in an isolated temp module.
//  4. Verify   — boot the built binary and diff its contract against the
//                benchmark. Target: 100% schema match.
//
// Every failure degrades gracefully: the blocker is reported transparently
// and the upstream runtime install remains available as the fallback.
package transpiler

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"scorp-agent/mcp"
)

// ProbeSpec describes how to spawn the ephemeral upstream server.
type ProbeSpec struct {
	Name    string
	Command string   // npx | uvx | python3 | node | absolute binary path
	Args    []string
	Env     map[string]string
	Timeout time.Duration
}

// Benchmark is the Phase-1 contract capture: the authoritative tool schemas
// the transpiled binary must reproduce. ProbeArgs records how the upstream
// was booted so verification can mirror the same environment.
type Benchmark struct {
	ServerName   string            `json:"server_name"`
	UpstreamInfo map[string]string `json:"upstream_info,omitempty"`
	ProbedAt     string            `json:"probed_at"`
	Command      string            `json:"command,omitempty"`
	Args         []string          `json:"args,omitempty"`
	Tools        []mcp.MCPTool     `json:"tools"`
}

// Probe spawns the upstream server in an ephemeral session, performs the
// JSON-RPC handshake, and captures the tools/list benchmark.
func Probe(ctx context.Context, spec ProbeSpec) (*Benchmark, error) {
	if spec.Command == "" {
		return nil, fmt.Errorf("probe: no command given")
	}
	if err := detectRuntime(spec.Command); err != nil {
		return nil, err
	}

	timeout := spec.Timeout
	if timeout == 0 {
		timeout = 90 * time.Second
	}
	probeCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// The probe itself must honor the deadline: JSON-RPC calls inside
	// startMCPServer block on pipes, so spawn in a goroutine and abandon it
	// if the upstream runtime hangs during boot.
	type result struct {
		srv   *mcp.MCPServer
		tools []mcp.MCPTool
		err   error
	}
	done := make(chan result, 1)
	go func() {
		srv, tools, err := mcp.ProbeServer("probe-"+spec.Name, mcp.MCPServerConfig{
			Command: spec.Command,
			Args:    spec.Args,
			Env:     spec.Env,
		})
		done <- result{srv, tools, err}
	}()

	select {
	case <-probeCtx.Done():
		return nil, fmt.Errorf("probe timed out after %s — upstream server did not complete its handshake (is the package name correct?)", timeout)
	case r := <-done:
		if r.err != nil {
			return nil, fmt.Errorf("probe failed: %w", r.err)
		}
		if len(r.tools) == 0 {
			return nil, fmt.Errorf("upstream server exposed zero tools — nothing to transpile")
		}
		return &Benchmark{
			ServerName: spec.Name,
			ProbedAt:   time.Now().UTC().Format(time.RFC3339),
			Command:    spec.Command,
			Args:       spec.Args,
			Tools:      r.tools,
		}, nil
	}
}

// detectRuntime verifies the runtime binary exists with an actionable message.
func detectRuntime(command string) error {
	if command == "" {
		return nil
	}
	if strings.ContainsAny(command, "/.") && !isKnownRuntime(command) {
		return nil // explicit path or dotted name — LookPath will resolve it
	}
	if _, err := exec.LookPath(command); err != nil {
		hint := map[string]string{
			"npx":      "Install Node.js (https://nodejs.org) so `npx` is on PATH.",
			"bunx":     "Install Bun (https://bun.sh) so `bunx` is on PATH.",
			"uvx":      "Install uv (https://docs.astral.sh/uv/) so `uvx` is on PATH.",
			"uv":       "Install uv (https://docs.astral.sh/uv/) so `uv` is on PATH.",
			"python3":  "Install Python 3 so `python3` is on PATH.",
			"python":   "Install Python 3 so `python` is on PATH.",
			"node":     "Install Node.js (https://nodejs.org) so `node` is on PATH.",
		}[command]
		if hint == "" {
			hint = "Make sure the runtime is installed and on PATH."
		}
		return fmt.Errorf("runtime %q not found: %s", command, hint)
	}
	return nil
}

func isKnownRuntime(cmd string) bool {
	switch cmd {
	case "npx", "uvx", "bunx", "python3", "python", "node", "uv":
		return true
	}
	return false
}

// MarshalBenchmark serializes a benchmark for the codegen prompt / debug dump.
func MarshalBenchmark(b *Benchmark) (string, error) {
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// DumpBenchmark persists the benchmark next to the transpile workspace for
// auditability.
func DumpBenchmark(workDir string, b *Benchmark) error {
	data, err := MarshalBenchmark(b)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(workDir, "benchmark.json"), []byte(data), 0o644)
}
