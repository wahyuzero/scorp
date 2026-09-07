//go:build windows

package tools

import (
	"os/exec"
)

func SetProcessGroup(cmd *exec.Cmd) {
	// Process groups are Unix-specific; standard process spawning on Windows
}

func KillProcessGroup(cmd *exec.Cmd) {
	if cmd != nil && cmd.Process != nil {
		_ = cmd.Process.Kill()
	}
}
