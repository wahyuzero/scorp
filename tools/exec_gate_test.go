package tools

import (
	"strings"
	"testing"

	"scorp-agent/config"
)

// TestShellDangerGateRespectsAutonomy pins the autonomy contract for the
// tool-layer shell gate: supervised stores a pending confirmation and refuses
// to run, YOLO executes unattended. Regression for YOLO mode being dead —
// the loop gate honored it but exec.go re-blocked the same command.
//
// The EXEC_MARKER_$((6*7)) construct only evaluates to EXEC_MARKER_42 when the
// command actually runs; the danger notice echoes the raw command text, so a
// literal substring of the command cannot distinguish echo from execution.
func TestShellDangerGateRespectsAutonomy(t *testing.T) {
	origDanger := IsDangerousCommand
	origStore := StorePendingConfirmation
	origAutonomy := config.GetAutonomyLevel()
	defer func() {
		IsDangerousCommand = origDanger
		StorePendingConfirmation = origStore
		config.SetAutonomyLevel(string(origAutonomy))
	}()

	IsDangerousCommand = func(cmd string) bool { return strings.Contains(cmd, "rm -rf") }
	stored := ""
	StorePendingConfirmation = func(chatID, toolName, command string, messages []AgentMessage, promptMsgID ...int64) {
		stored = toolName + ":" + command
	}

	config.SetAutonomyLevel("supervised")
	stored = ""
	out, ok := ExecuteShell(map[string]interface{}{
		"command": "rm -rf /tmp/scorp-gate-nonexistent && echo EXEC_MARKER_$((6*7))",
	}, 0)
	if ok {
		t.Fatalf("supervised mode must not execute a dangerous command, got ok=true output=%q", out)
	}
	if !strings.Contains(out, "DANGEROUS COMMAND DETECTED") {
		t.Fatalf("supervised mode must return the danger notice, got %q", out)
	}
	if !strings.HasPrefix(stored, "shell:") {
		t.Fatalf("supervised mode must store a pending confirmation, stored=%q", stored)
	}
	if strings.Contains(out, "EXEC_MARKER_42") {
		t.Fatalf("command must not run in supervised mode, output=%q", out)
	}

	config.SetAutonomyLevel("yolo")
	stored = ""
	out, ok = ExecuteShell(map[string]interface{}{
		"command": "rm -rf /tmp/scorp-gate-nonexistent && echo EXEC_MARKER_$((6*7))",
	}, 0)
	if !ok {
		t.Fatalf("yolo mode must execute a dangerous command without confirmation, output=%q", out)
	}
	if !strings.Contains(out, "EXEC_MARKER_42") {
		t.Fatalf("yolo mode command did not actually execute, output=%q", out)
	}
	if stored != "" {
		t.Fatalf("yolo mode must not store a pending confirmation, stored=%q", stored)
	}
}

// TestShellSensitivePathSandboxAllModes pins the sensitive-path sandbox as a
// hard layer above both the confirmation gate and the autonomy level: reading
// protected credentials through the shell must be denied in supervised AND
// yolo mode, even when a confirmation flag is attached. Regression for the
// sandbox escape where read_file was blocked but `cat` ran fine.
func TestShellSensitivePathSandboxAllModes(t *testing.T) {
	origDanger := IsDangerousCommand
	origAutonomy := config.GetAutonomyLevel()
	defer func() {
		IsDangerousCommand = origDanger
		config.SetAutonomyLevel(string(origAutonomy))
	}()

	IsDangerousCommand = func(cmd string) bool { return false }

	for _, mode := range []string{"supervised", "yolo"} {
		config.SetAutonomyLevel(mode)
		for _, cmd := range []string{
			"cat /etc/shadow",
			"cat ~/.ssh/id_rsa",
			"tar czf /tmp/x.tgz /etc/sudoers",
		} {
			out, ok := ExecuteShell(map[string]interface{}{
				"command":   cmd,
				"confirmed": true,
			}, 0)
			if ok {
				t.Fatalf("[%s] sandbox must block %q, got ok=true output=%q", mode, cmd, out)
			}
			if !strings.Contains(out, "Security Sandbox") {
				t.Fatalf("[%s] expected sandbox denial for %q, got %q", mode, cmd, out)
			}
		}
	}
}
