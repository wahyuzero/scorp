package tools

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// ──────────────────────────────────────────────
// Background Process Management
// ──────────────────────────────────────────────

type BGProcess struct {
	ID        string
	SessionID string      // Session that owns this process
	Cmd       *exec.Cmd
	Stdin     io.WriteCloser
	Stdout    *bytes.Buffer
	Stderr    *bytes.Buffer
	StartTime time.Time
	Done      bool
	ExitCode  int
	mu        sync.Mutex
}

var (
	bgProcesses   = make(map[string]*BGProcess)
	bgProcessesMu sync.RWMutex
	bgProcessID   int64
)

func nextBGProcessID() string {
	bgProcessID++
	return fmt.Sprintf("bg_%d", bgProcessID)
}

// executeBgProcess handles the "bg" tool
func ExecuteBgProcess(args map[string]interface{}) (string, bool) {
	action, _ := args["action"].(string)
	sessionID, _ := args["session_id"].(string)
	command, _ := args["command"].(string)
	workdir, _ := args["workdir"].(string)

	switch action {
	case "spawn":
		return bgSpawn(command, workdir)
	case "list":
		return bgList()
	case "poll":
		return bgPoll(sessionID)
	case "wait":
		timeout, _ := args["timeout"].(float64)
		return bgWait(sessionID, int(timeout))
	case "kill":
		return bgKill(sessionID)
	case "write":
		data, _ := args["data"].(string)
		return bgWrite(sessionID, data, false)
	case "submit":
		data, _ := args["data"].(string)
		return bgWrite(sessionID, data, true)
	default:
		return fmt.Sprintf("Unknown action: %s. Use: spawn, list, poll, wait, kill, write, submit", action), false
	}
}

func bgSpawn(command, workdir string) (string, bool) {
	if command == "" {
		return "Missing 'command' parameter", false
	}

	// Create the command
	cmd := exec.Command("/bin/bash", "-c", command)
	if workdir != "" {
		cmd.Dir = workdir
	}

	// Capture stdout and stderr
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	// Stdin pipe for write/submit
	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Sprintf("Failed to create stdin pipe: %v", err), false
	}

	// Start the process
	if err := cmd.Start(); err != nil {
		return fmt.Sprintf("Failed to start process: %v", err), false
	}

	id := nextBGProcessID()
	proc := &BGProcess{
		ID:        id,
		Cmd:       cmd,
		Stdin:     stdinPipe,
		Stdout:    &stdoutBuf,
		Stderr:    &stderrBuf,
		StartTime: time.Now(),
	}

	bgProcessesMu.Lock()
	bgProcesses[id] = proc
	bgProcessesMu.Unlock()

	// Monitor process completion in background
	go func() {
		err := cmd.Wait()
		proc.mu.Lock()
		proc.Done = true
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				proc.ExitCode = exitErr.ExitCode()
			} else {
				proc.ExitCode = -1
			}
		} else {
			proc.ExitCode = 0
		}
		proc.mu.Unlock()

		// Auto-cleanup after 30 min idle
		go func() {
			time.Sleep(30 * time.Minute)
			bgProcessesMu.Lock()
			delete(bgProcesses, id)
			bgProcessesMu.Unlock()
		}()
	}()

	return fmt.Sprintf("✅ Spawned process <b>%s</b>\nCommand: <code>%s</code>\nSession: <code>%s</code>", id, command, id), true
}

func bgList() (string, bool) {
	bgProcessesMu.RLock()
	defer bgProcessesMu.RUnlock()

	if len(bgProcesses) == 0 {
		return "No background processes.", true
	}

	var sb strings.Builder
	sb.WriteString("📋 <b>Background Processes</b>\n\n")
	for id, proc := range bgProcesses {
		proc.mu.Lock()
		status := "🏃 Running"
		if proc.Done {
			if proc.ExitCode == 0 {
				status = "✅ Completed"
			} else {
				status = fmt.Sprintf("❌ Failed (exit %d)", proc.ExitCode)
			}
		}
		running := time.Since(proc.StartTime).Round(time.Second).String()
		proc.mu.Unlock()

		sb.WriteString(fmt.Sprintf("• <b>%s</b> — %s (%s)\n", id, status, running))
	}

	return sb.String(), true
}

func bgPoll(sessionID string) (string, bool) {
	proc, ok := getBGProcess(sessionID)
	if !ok {
		return fmt.Sprintf("Process not found: %s", sessionID), false
	}

	proc.mu.Lock()
	defer proc.mu.Unlock()

	status := "running"
	if proc.Done {
		if proc.ExitCode == 0 {
			status = "completed"
		} else {
			status = fmt.Sprintf("failed (exit %d)", proc.ExitCode)
		}
	}

	stdout := proc.Stdout.String()
	stderr := proc.Stderr.String()

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📋 <b>Process %s</b> — %s\n\n", sessionID, status))

	if stdout != "" {
		sb.WriteString("📤 <b>stdout:</b>\n")
		sb.WriteString(fmt.Sprintf("<pre>%s</pre>\n", truncateString(stdout, 1000)))
	}
	if stderr != "" {
		sb.WriteString("📥 <b>stderr:</b>\n")
		sb.WriteString(fmt.Sprintf("<pre>%s</pre>\n", truncateString(stderr, 1000)))
	}
	if stdout == "" && stderr == "" {
		sb.WriteString("(no output yet)")
	}

	return sb.String(), true
}

func bgWait(sessionID string, timeout int) (string, bool) {
	proc, ok := getBGProcess(sessionID)
	if !ok {
		return fmt.Sprintf("Process not found: %s", sessionID), false
	}

	if timeout <= 0 {
		timeout = 60 // default 60s
	}

	// Poll until done or timeout
	start := time.Now()
	for {
		proc.mu.Lock()
		done := proc.Done
		proc.mu.Unlock()

		if done {
			return bgPoll(sessionID)
		}

		if time.Since(start) > time.Duration(timeout)*time.Second {
			return bgPoll(sessionID)
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func bgKill(sessionID string) (string, bool) {
	proc, ok := getBGProcess(sessionID)
	if !ok {
		return fmt.Sprintf("Process not found: %s", sessionID), false
	}

	proc.mu.Lock()
	defer proc.mu.Unlock()

	if proc.Done {
		return fmt.Sprintf("Process %s already finished (exit code %d)", sessionID, proc.ExitCode), true
	}

	if err := proc.Cmd.Process.Kill(); err != nil {
		return fmt.Sprintf("Failed to kill process: %v", err), false
	}

	proc.Done = true
	proc.ExitCode = -1

	return fmt.Sprintf("✅ Process <b>%s</b> killed", sessionID), true
}

func bgWrite(sessionID, data string, newline bool) (string, bool) {
	proc, ok := getBGProcess(sessionID)
	if !ok {
		return fmt.Sprintf("Process not found: %s", sessionID), false
	}

	proc.mu.Lock()
	defer proc.mu.Unlock()

	if proc.Done {
		return fmt.Sprintf("Process %s already finished", sessionID), false
	}

	if newline {
		_, err := fmt.Fprintln(proc.Stdin, data)
		if err != nil {
			return fmt.Sprintf("Failed to write stdin: %v", err), false
		}
	} else {
		_, err := fmt.Fprint(proc.Stdin, data)
		if err != nil {
			return fmt.Sprintf("Failed to write stdin: %v", err), false
		}
	}

	return fmt.Sprintf("✅ Sent %d bytes to process %s", len(data), sessionID), true
}

// closeStdin closes stdin for a background process (sends EOF)
func closeStdin(sessionID string) (string, bool) {
	proc, ok := getBGProcess(sessionID)
	if !ok {
		return fmt.Sprintf("Process not found: %s", sessionID), false
	}

	proc.mu.Lock()
	defer proc.mu.Unlock()

	if proc.Done {
		return fmt.Sprintf("Process %s already finished", sessionID), false
	}

	if err := proc.Stdin.Close(); err != nil {
		return fmt.Sprintf("Failed to close stdin: %v", err), false
	}

	return fmt.Sprintf("✅ Stdin closed for process %s", sessionID), true
}

// ──────────────────────────────────────────────
// Helpers
// ──────────────────────────────────────────────

func getBGProcess(id string) (*BGProcess, bool) {
	bgProcessesMu.RLock()
	defer bgProcessesMu.RUnlock()
	proc, ok := bgProcesses[id]
	return proc, ok
}

func truncateString(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) > maxLen {
		return string(runes[:maxLen]) + "..."
	}
	return s
}

// ──────────────────────────────────────────────
// Helpers
// ──────────────────────────────────────────────