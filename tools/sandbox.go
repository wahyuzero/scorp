package tools

import (
	"context"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"scorp-agent/config"
)

// ──────────────────────────────────────────────
// Shell Sandbox (P0.1) — bubblewrap isolation for model-initiated shell calls
//
// Design (docs/IMPLEMENTATION_PLAN_SCORP.md P0.1, evidence: Anthropic reports
// sandboxing cuts permission prompts by 84%):
//   --ro-bind / /     → system trees (/etc, /usr, /var, /opt, /boot) read-only
//   --dev-bind /dev   → device access preserved (fuse, tty, docker via /run)
//   --proc /proc      → real proc (ps, top, systemctl status read sysfs)
//   rw binds          → $HOME, cwd, /tmp, plus SCORP_SANDBOX_RW extras
//                       (default /run,/var/run,/data) so systemd, docker
//                       sockets, coolify data and cross-command /tmp flows
//                       keep working — this is an ops agent, not a CI runner.
//   --unshare-all --share-net → fresh namespaces, network intact (v1)
//   --die-with-parent --new-session → no orphaned sandbox processes
//
// Active by default ("on"); SCORP_SANDBOX=off is the escape hatch. A one-time
// smoke test gates activation: if bubblewrap is missing or the kernel refuses
// user namespaces, the sandbox silently disables and SandboxStatusNotice
// surfaces the fact in /status — including the permanent YOLO warning.
// ──────────────────────────────────────────────

const (
	sandboxBin      = "bwrap"
	sandboxRWDefault = "/run,/var/run,/data"
)

var (
	sandboxSmokeOnce sync.Once
	sandboxSmokeOK   bool
)

// SandboxModeEnabled reports whether SCORP_SANDBOX requests the sandbox.
// Unset or any value except off/false/0/no means enabled (default on).
func SandboxModeEnabled() bool {
	switch strings.ToLower(os.Getenv("SCORP_SANDBOX")) {
	case "off", "false", "0", "no":
		return false
	}
	return true
}

// sandboxRWPaths returns resolved, existing writable bind targets: home, cwd,
// /tmp and the configured extras, symlink-resolved and de-duplicated.
func sandboxRWPaths() []string {
	requested := []string{
		os.TempDir(),
		config.HomeDir(),
	}
	if wd, err := os.Getwd(); err == nil {
		requested = append(requested, wd)
	}
	extras := os.Getenv("SCORP_SANDBOX_RW")
	if strings.TrimSpace(extras) == "" {
		extras = sandboxRWDefault
	}
	requested = append(requested, strings.Split(extras, ",")...)

	var paths []string
	seen := map[string]bool{}
	for _, p := range requested {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if resolved, err := filepath.EvalSymlinks(p); err == nil {
			p = resolved
		}
		if fi, err := os.Stat(p); err != nil || !fi.IsDir() {
			continue
		}
		if seen[p] {
			continue
		}
		seen[p] = true
		paths = append(paths, p)
	}
	return paths
}

// sandboxArgv builds the full bubblewrap argv for a shell command.
func sandboxArgv(command string) []string {
	argv := []string{
		sandboxBin,
		"--unshare-all", "--share-net",
		"--die-with-parent", "--new-session",
		"--ro-bind", "/", "/",
		"--dev-bind", "/dev", "/dev",
		"--proc", "/proc",
	}
	for _, p := range sandboxRWPaths() {
		argv = append(argv, "--bind", p, p)
	}
	return append(argv, "--", "bash", "-c", command)
}

// sandboxSmokeTest runs the actual sandbox once with a trivial command; its
// result is cached for the process lifetime. A sandbox that cannot start must
// disable itself rather than break every shell call.
func sandboxSmokeTest() bool {
	sandboxSmokeOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		argv := sandboxArgv("echo scorp-sandbox-ok")
		out, err := exec.CommandContext(ctx, argv[0], argv[1:]...).CombinedOutput()
		sandboxSmokeOK = err == nil && strings.Contains(string(out), "scorp-sandbox-ok")
		if !sandboxSmokeOK {
			log.Printf("[sandbox] smoke test FAILED — shell sandbox disabled. err=%v out=%s",
				err, strings.TrimSpace(string(out)))
		} else {
			log.Printf("[sandbox] smoke test passed — shell sandbox ACTIVE (bubblewrap)")
		}
	})
	return sandboxSmokeOK
}

// SandboxActive reports whether model shell commands are currently wrapped in
// the bubblewrap sandbox.
func SandboxActive() bool {
	return SandboxModeEnabled() && sandboxSmokeTest()
}

// SandboxWrap returns the argv to execute a command inside the sandbox, or
// (nil, false) when the sandbox is disabled/unavailable and the caller should
// run its plain exec path.
func SandboxWrap(command string) ([]string, bool) {
	if !SandboxActive() {
		return nil, false
	}
	return sandboxArgv(command), true
}

// SandboxVersion returns the bubblewrap version string, or "" if unavailable.
func SandboxVersion() string {
	out, err := exec.Command(sandboxBin, "--version").CombinedOutput()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// SandboxStatusNotice renders the sandbox line for /status and /agent info,
// including the permanent warning when YOLO runs without a sandbox (plan P0.1:
// "YOLO tanpa bwrap → peringatan permanen di status").
func SandboxStatusNotice() string {
	var sb strings.Builder
	switch {
	case SandboxActive():
		v := SandboxVersion()
		if v == "" {
			v = "bubblewrap"
		}
		sb.WriteString("\n🧱 Shell sandbox: <b>ON</b> (" + v + ") — /etc, /usr, /var, /opt read-only; writable: home, cwd, /tmp" +
			" (extra rw: SCORP_SANDBOX_RW)")
	case !SandboxModeEnabled():
		sb.WriteString("\n🧱 Shell sandbox: <b>OFF</b> (SCORP_SANDBOX=off)")
	default:
		sb.WriteString("\n🧱 Shell sandbox: <b>OFF</b> (bubblewrap unavailable or smoke test failed — see daemon log)")
	}
	if config.GetAutonomyLevel() == config.AutonomyYOLO && !SandboxActive() {
		sb.WriteString("\n⚠️ <b>YOLO is running WITHOUT a sandbox</b> — every shell command has full root-level system access. Install bubblewrap or unset SCORP_SANDBOX=off.")
	}
	return sb.String()
}
