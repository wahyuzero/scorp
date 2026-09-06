package tools

import (
	"strings"
	"testing"
)

// stubDanger installs an inert dangerous-command classifier for the test —
// production wires tools.IsDangerousCommand at startup (daemon/CLI), tests
// must provide it themselves (same convention as exec_gate_test.go).
func stubDanger(t *testing.T) {
	t.Helper()
	orig := IsDangerousCommand
	IsDangerousCommand = func(cmd string) bool { return false }
	t.Cleanup(func() { IsDangerousCommand = orig })
}

// TestSandboxWrapArgvShape validates the bwrap argv contract when the sandbox
// is active (bubblewrap present + smoke test passed). When it is not, the test
// pins the fallback: SandboxWrap must report "not wrapped".
func TestSandboxWrapArgvShape(t *testing.T) {
	argv, wrapped := SandboxWrap("echo hi")
	if !SandboxActive() {
		if wrapped {
			t.Fatalf("sandbox inactive but SandboxWrap wrapped the command: %v", argv)
		}
		t.Skip("sandbox inactive on this host (no bwrap or smoke test failed) — argv shape test skipped")
	}
	if !wrapped {
		t.Fatal("sandbox active but SandboxWrap did not wrap")
	}
	joined := strings.Join(argv, " ")
	for _, want := range []string{
		"bwrap",
		"--ro-bind / /",
		"--dev-bind /dev /dev",
		"--proc /proc",
		"--unshare-all",
		"--share-net",
		"--die-with-parent",
		"-- bash -c echo hi",
	} {
		if !strings.Contains(joined, want) {
			t.Errorf("sandbox argv missing %q, got: %s", want, joined)
		}
	}
	// home and cwd must be bound read-write (later --bind overrides --ro-bind)
	if !strings.Contains(joined, "--bind /tmp /tmp") {
		t.Errorf("sandbox argv must rw-bind /tmp, got: %s", joined)
	}
}

// TestSandboxedShellExecutes proves real execution still works inside the
// sandbox: arithmetic evaluation only happens if the command actually ran.
func TestSandboxedShellExecutes(t *testing.T) {
	stubDanger(t)
	if !SandboxActive() {
		t.Skip("sandbox inactive on this host")
	}
	out, ok := ExecuteShell(map[string]interface{}{"command": "echo SBX_$((4*9))"}, 0)
	if !ok {
		t.Fatalf("sandboxed command failed: %q", out)
	}
	if !strings.Contains(out, "SBX_36") {
		t.Fatalf("sandboxed command did not actually execute: %q", out)
	}
}

// TestSandboxBlocksSystemWrite is the core isolation assertion: with the
// sandbox active, writing under /usr must fail even for root (the daemon's
// case on the VPS). Skipped when the sandbox is inactive.
func TestSandboxBlocksSystemWrite(t *testing.T) {
	stubDanger(t)
	if !SandboxActive() {
		t.Skip("sandbox inactive on this host")
	}
	out, ok := ExecuteShell(map[string]interface{}{
		"command": "touch /usr/local/bin/scorp-sbx-probe && echo TOUCH_OK || echo TOUCH_DENIED:$?",
	}, 0)
	if ok && strings.Contains(out, "TOUCH_OK") {
		t.Fatalf("sandbox allowed a write to /usr — isolation broken: %q", out)
	}
	if !strings.Contains(out, "TOUCH_DENIED") {
		t.Fatalf("expected in-shell write failure under /usr, got ok=%v out=%q", ok, out)
	}
}

// TestSandboxCrossCommandTmpfsFile guards the workflow guarantee: files
// written to /tmp in one sandboxed command must still exist in the next one
// (/tmp is rw-BIND, not a per-command tmpfs).
func TestSandboxCrossCommandTmpfsFile(t *testing.T) {
	stubDanger(t)
	if !SandboxActive() {
		t.Skip("sandbox inactive on this host")
	}
	if _, ok := ExecuteShell(map[string]interface{}{"command": "echo cross-marker > /tmp/scorp-sbx-cross.txt"}, 0); !ok {
		t.Fatal("failed to write the cross-command marker")
	}
	out, ok := ExecuteShell(map[string]interface{}{"command": "cat /tmp/scorp-sbx-cross.txt"}, 0)
	if !ok || !strings.Contains(out, "cross-marker") {
		t.Fatalf("cross-command /tmp persistence broken inside sandbox: ok=%v out=%q", ok, out)
	}
}

func TestSandboxStatusNotice(t *testing.T) {
	notice := SandboxStatusNotice()
	if !strings.Contains(notice, "Shell sandbox") {
		t.Fatalf("status notice must mention the sandbox, got %q", notice)
	}
}
