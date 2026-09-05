package transpiler

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// goModTemplate pins the dependency graph for generated ports. Only the SDK
// plus optional audited helpers live here; the sandbox has no other sources.
const goModTemplate = `module scorp.transpiled/%s

go 1.24.0

require github.com/mark3labs/mcp-go v1.0.0
`

// buildTimeouts bound the sandbox toolchain (misbehaving codegen must not
// spin forever).
const (
	tidyTimeout  = 5 * time.Minute
	buildTimeout = 5 * time.Minute
)

// Sandbox is the isolated build workspace for one transpilation.
type Sandbox struct {
	Dir     string
	BinPath string
	MainGo  string
}

// NewSandbox creates an isolated temp module for building generated code.
func NewSandbox(name, mainGo string) (*Sandbox, error) {
	dir, err := os.MkdirTemp("", "scorp-transpile-"+name+"-*")
	if err != nil {
		return nil, fmt.Errorf("create sandbox: %w", err)
	}
	return &Sandbox{
		Dir:     dir,
		MainGo:  mainGo,
		BinPath: filepath.Join(dir, "server_bin"),
	}, nil
}

// Build writes the module files and compiles. go.sum is produced by
// `go mod tidy` against the pinned go.mod; the module cache stays the user's
// default GOPATH cache (container hardening is an M3 concern).
func (s *Sandbox) Build() error {
	if err := os.WriteFile(filepath.Join(s.Dir, "main.go"), []byte(s.MainGo), 0o644); err != nil {
		return fmt.Errorf("write main.go: %w", err)
	}
	goMod := fmt.Sprintf(goModTemplate, sanitizeModuleName(s.Dir))
	if err := os.WriteFile(filepath.Join(s.Dir, "go.mod"), []byte(goMod), 0o644); err != nil {
		return fmt.Errorf("write go.mod: %w", err)
	}

	if out, err := runIn(s.Dir, tidyTimeout, "go", "mod", "tidy"); err != nil {
		return fmt.Errorf("go mod tidy failed (dependency graph could not be resolved):\n%s", tail(out, 40))
	}

	if out, err := runIn(s.Dir, buildTimeout, "go", "build", "-trimpath", "-o", s.BinPath, "."); err != nil {
		return fmt.Errorf("go build failed:\n%s", tail(out, 60))
	}
	return nil
}

// Close removes the sandbox workspace (after install/verify completes).
func (s *Sandbox) Close() {
	_ = os.RemoveAll(s.Dir)
}

func runIn(dir string, timeout time.Duration, name string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func sanitizeModuleName(path string) string {
	base := filepath.Base(path)
	clean := make([]rune, 0, len(base))
	for _, r := range base {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			clean = append(clean, r)
		}
	}
	if len(clean) == 0 {
		return "server"
	}
	return string(clean)
}

func tail(s string, lines int) string {
	all := strings.Split(strings.TrimRight(s, "\n"), "\n")
	if len(all) > lines {
		all = all[len(all)-lines:]
	}
	return strings.Join(all, "\n")
}
