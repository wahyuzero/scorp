package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"scorp-agent/config"
	"scorp-agent/internal/helpers"
	"strings"
	"sync"
	"time"
	"unicode/utf16"
)

// shellQuote safely quotes a string for use in shell commands
func shellQuote(s string) string {
	// If string is empty or contains only safe chars, return as-is or with single quotes
	if s == "" {
		return "''"
	}
	// Check if string contains only safe characters (alphanumeric, underscore, dash, dot, slash, colon, @)
	safe := true
	for _, c := range s {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') ||
			c == '_' || c == '-' || c == '.' || c == '/' || c == ':' || c == '@' || c == '=') {
			safe = false
			break
		}
	}
	if safe {
		return s
	}
	// Escape single quotes and wrap in single quotes
	return "'" + strings.ReplaceAll(s, "'", "'\"'\"'") + "'"
}

// needsShellExecution detects if a command requires shell features (pipes, redirects, variables, etc.)
// Returns false for simple command + args even if they contain parentheses, brackets, etc.
func needsShellExecution(command string) bool {
	// Trim leading/trailing whitespace
	command = strings.TrimSpace(command)
	if command == "" {
		return false
	}

	// Check for shell constructs that require bash -c
	// These must appear OUTSIDE of quotes
	inSingleQuote := false
	inDoubleQuote := false

	for i := 0; i < len(command); i++ {
		c := command[i]

		// Handle quotes
		if c == '\'' && !inDoubleQuote {
			inSingleQuote = !inSingleQuote
			continue
		}
		if c == '"' && !inSingleQuote {
			inDoubleQuote = !inDoubleQuote
			continue
		}

		// Skip if inside quotes
		if inSingleQuote || inDoubleQuote {
			continue
		}

		// Shell metacharacters that require shell execution
		switch c {
		case '|', '&', ';', '<', '>', '$', '`', '*', '?', '[', ']', '{', '}', '\\':
			return true
		}
	}

	return false
}

// ── Shell Executor ──

// ExecuteShell runs a shell command
func ExecuteShell(args map[string]interface{}, chatID int64) (string, bool) {
	command := helpers.GetStringArg(args, "command", "")
	if command == "" {
		return "Error: 'command' argument is required", false
	}

	timeout := helpers.GetIntArg(args, "timeout", 30)
	if timeout > 300 {
		timeout = 300
	}

	// Deny rules are the absolute outer layer (P0.2): they precede the
	// sensitive-path sandbox, the confirmation gate and the autonomy level, so
	// they hold in YOLO and on user-confirmed resumes too.
	if blocked, reason := config.CheckDenyRules("shell", args); blocked {
		return "🚫 " + reason, false
	}

	// Sensitive-path sandbox applies in EVERY autonomy mode (including YOLO):
	// structured tools are already blocked from touching protected credentials,
	// so the shell must not become the bypass around them. Runs before the
	// confirmation gate — even a human-confirmed command may not read these.
	if restricted, reason := config.IsPathRestricted(command); restricted {
		return "🛡️ " + reason, false
	}

	// Check for dangerous commands (skip if already explicitly confirmed by user,
	// or when YOLO autonomy runs unattended — see config.ConfirmationRequired)
	confirmed := helpers.GetBoolArg(args, "confirmed", false)
	if !confirmed && config.ConfirmationRequired() && IsDangerousCommand(command) {
		// Store for confirmation
		chatIDStr := fmt.Sprintf("%d", chatID)
		StorePendingConfirmation(chatIDStr, "shell", command, nil)
		return fmt.Sprintf("⚠️ DANGEROUS COMMAND DETECTED:\n%s\n\nPlease confirm execution.", command), false
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	var cmd *exec.Cmd
	if argv, wrapped := SandboxWrap(command); wrapped {
		// P0.1: run inside the bubblewrap sandbox (system trees read-only).
		cmd = exec.Command(argv[0], argv[1:]...)
	} else {
		cmd = exec.Command("bash", "-c", command)
	}
	// Own process group so the timeout can kill backgrounded grandchildren —
	// otherwise `server &` inherits the stdout pipe and CombinedOutput blocks
	// forever even after the direct child exits.
	SetProcessGroup(cmd)

	// Sanitize LD_PRELOAD for PRoot Linux environments inside Termux:
	// When running inside a PRoot Linux distro (Debian, Ubuntu, Arch, etc.) on Android,
	// Termux's Android Bionic libtermux-exec-ld-preload.so crashes glibc child processes with static TLS errors.
	// Only strip when running inside a guest Linux rootfs (/etc/os-release exists) with Termux LD_PRELOAD present.
	if lp := os.Getenv("LD_PRELOAD"); strings.Contains(lp, "libtermux-exec") {
		if isGuestLinuxRootfs() {
			var cleanEnv []string
			for _, e := range os.Environ() {
				if !strings.HasPrefix(e, "LD_PRELOAD=") {
					cleanEnv = append(cleanEnv, e)
				}
			}
			cmd.Env = cleanEnv
		}
	}

	// Set parent tracking env vars so self-spawned Scorp sub-processes know they are children
	cleanEnv := cmd.Env
	if len(cleanEnv) == 0 {
		cleanEnv = os.Environ()
	}
	cleanEnv = append(cleanEnv, fmt.Sprintf("SCORP_PARENT_PID=%d", os.Getpid()))
	cleanEnv = append(cleanEnv, fmt.Sprintf("SCORP_PARENT_SESSION=%d", chatID))

	// Inject persistent session environment variables
	if sessEnvs := GetSessionEnv(chatID); len(sessEnvs) > 0 {
		for k, v := range sessEnvs {
			cleanEnv = append(cleanEnv, fmt.Sprintf("%s=%s", k, v))
		}
	}
	cmd.Env = cleanEnv

	// Auto-capture any explicit static export KEY=VAL from command to persist for next turns
	captureSessionExports(chatID, command)

	// Non-interactive stdin guard: prevent child processes from blocking indefinitely
	// waiting for interactive keyboard input (e.g. read -p, sudo password, python input()).
	cmd.Stdin = strings.NewReader("")

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Sprintf("Failed to create stdout pipe: %v", err), false
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Sprintf("Failed to create stderr pipe: %v", err), false
	}

	if err := cmd.Start(); err != nil {
		return fmt.Sprintf("Command failed to start: %v", err), false
	}

	var buf bytes.Buffer
	var bufMu sync.Mutex

	// Read stdout and stderr concurrently with a buffer limit
	readDone := make(chan struct{}, 2)
	streamReader := func(r io.Reader) {
		defer func() { readDone <- struct{}{} }()
		temp := make([]byte, 4096)
		for {
			n, rErr := r.Read(temp)
			if n > 0 {
				bufMu.Lock()
				if buf.Len() < helpers.MaxToolOutput*3 {
					buf.Write(temp[:n])
				}
				bufMu.Unlock()
			}
			if rErr != nil {
				break
			}
		}
	}

	go streamReader(stdoutPipe)
	go streamReader(stderrPipe)

	waitDone := make(chan error, 1)
	go func() {
		waitDone <- cmd.Wait()
	}()

	var waitErr error
	timedOut := false

	select {
	case <-ctx.Done():
		timedOut = true
		if cmd.Process != nil {
			KillProcessGroup(cmd)
		}
		waitErr = fmt.Errorf("timeout")
	case waitErr = <-waitDone:
		// Process exited. If background children still hold the pipe, give them a short grace period
		// then close the pipes so cmd doesn't hang the loop forever.
		select {
		case <-readDone:
		case <-time.After(500 * time.Millisecond):
			stdoutPipe.Close()
			stderrPipe.Close()
		}
	}

	bufMu.Lock()
	result := buf.String()
	bufMu.Unlock()

	if timedOut {
		// Clean up any lingering POSIX shared memory files created in /dev/shm during aborted runs
		cleanAbortedShm()
		return fmt.Sprintf("Command timed out after %ds (background processes holding the pipe are killed — for daemons use: nohup CMD >/tmp/log 2>&1 & disown):\n%s",
			timeout, helpers.TruncOutput(result, helpers.MaxToolOutput)), false
	}

	if waitErr != nil {
		return fmt.Sprintf("Command failed: %v\nOutput:\n%s", waitErr, helpers.TruncOutput(result, helpers.MaxToolOutput)), false
	}

	trimmedOut := strings.TrimSpace(result)
	if trimmedOut == "" {
		return "(Command executed successfully with exit code 0 and empty output)", true
	}

	return helpers.TruncOutput(result, helpers.MaxToolOutput), true
}

// ── Session Scoped Environment Storage ──

var (
	sessionEnvMap = make(map[int64]map[string]string)
	sessionEnvMu  sync.RWMutex
)

func sessionEnvFilePath(chatID int64) string {
	dir := config.ScorpPath("session_envs")
	_ = os.MkdirAll(dir, 0755)
	return filepath.Join(dir, fmt.Sprintf("env_%d.json", chatID))
}

func loadSessionEnvFromDisk(chatID int64) map[string]string {
	p := sessionEnvFilePath(chatID)
	data, err := os.ReadFile(p)
	if err != nil {
		return nil
	}
	var res map[string]string
	if err := json.Unmarshal(data, &res); err == nil {
		return res
	}
	return nil
}

func saveSessionEnvToDisk(chatID int64, envs map[string]string) {
	p := sessionEnvFilePath(chatID)
	data, err := json.Marshal(envs)
	if err == nil {
		_ = WriteFileAtomic(p, data, 0644)
	}
}

// SetSessionEnv sets a persistent environment variable scoped to a session
func SetSessionEnv(chatID int64, key, val string) {
	sessionEnvMu.Lock()
	defer sessionEnvMu.Unlock()
	if sessionEnvMap[chatID] == nil {
		sessionEnvMap[chatID] = loadSessionEnvFromDisk(chatID)
		if sessionEnvMap[chatID] == nil {
			sessionEnvMap[chatID] = make(map[string]string)
		}
	}
	sessionEnvMap[chatID][key] = val
	saveSessionEnvToDisk(chatID, sessionEnvMap[chatID])
}

// GetSessionEnv returns the map of persistent environment variables for a session
func GetSessionEnv(chatID int64) map[string]string {
	sessionEnvMu.RLock()
	defer sessionEnvMu.RUnlock()
	m := sessionEnvMap[chatID]
	if m == nil {
		m = loadSessionEnvFromDisk(chatID)
		if m != nil {
			sessionEnvMap[chatID] = m
		}
	}
	if m == nil {
		return nil
	}
	cp := make(map[string]string, len(m))
	for k, v := range m {
		cp[k] = v
	}
	return cp
}

// captureSessionExports extracts top-level 'export KEY=VALUE' assignments from commands
func captureSessionExports(chatID int64, command string) {
	lines := strings.Split(command, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		parts := strings.Split(trimmed, ";")
		for _, part := range parts {
			subParts := strings.Split(part, "&&")
			for _, sp := range subParts {
				spTrim := strings.TrimSpace(sp)
				if strings.HasPrefix(spTrim, "export ") {
					assign := strings.TrimPrefix(spTrim, "export ")
					eqIdx := strings.Index(assign, "=")
					if eqIdx > 0 {
						k := strings.TrimSpace(assign[:eqIdx])
						v := strings.TrimSpace(assign[eqIdx+1:])
						if (strings.HasPrefix(v, "\"") && strings.HasSuffix(v, "\"")) ||
							(strings.HasPrefix(v, "'") && strings.HasSuffix(v, "'")) {
							v = v[1 : len(v)-1]
						}
						if !strings.Contains(v, "$(") && !strings.Contains(v, "`") {
							SetSessionEnv(chatID, k, v)
						}
					}
				}
			}
		}
	}
}

// ── File Reader ──

var allowedReadPaths = []string{
	config.HomeDir() + "/",
	"/root/",
	"/home/",
	"/tmp/",
	"/etc/",
	"/var/log/",
	"/data/coolify/",
	"/opt/",
}

func isPathAllowed(path string, allowedPrefixes []string) bool {
	// If Sandbox is disabled, host operations have unrestricted filesystem access
	if !SandboxModeEnabled() {
		return true
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	cwd, err := os.Getwd()
	if err == nil && cwd != "" {
		cleanCwd := filepath.Clean(cwd)
		if strings.HasPrefix(absPath, cleanCwd) {
			return true
		}
	}
	for _, prefix := range allowedPrefixes {
		cleanPrefix := filepath.Clean(prefix)
		if strings.HasPrefix(absPath, cleanPrefix) {
			return true
		}
	}
	return false
}

// ExecuteReadFile reads a file with support for offset, lines/limit, and line numbering
func ExecuteReadFile(args map[string]interface{}) (string, bool) {
	path := helpers.GetStringArg(args, "path", "")
	if path == "" {
		return "Error: 'path' argument is required", false
	}

	if !isPathAllowed(path, allowedReadPaths) {
		return fmt.Sprintf("Error: path '%s' is not in allowed directories", path), false
	}

	// Guard against unbounded memory allocation on large or binary files (e.g. 50MB+ dumps).
	// Read up to 2MB max directly through a bounded limit reader.
	f, err := os.Open(path)
	if err != nil {
		return fmt.Sprintf("Error opening file: %v", err), false
	}
	defer f.Close()

	const maxReadSize = 2 * 1024 * 1024 // 2MB
	limitedReader := io.LimitReader(f, maxReadSize+1)
	data, err := io.ReadAll(limitedReader)
	if err != nil {
		return fmt.Sprintf("Error reading file: %v", err), false
	}
	truncatedNotice := ""
	if len(data) > maxReadSize {
		data = data[:maxReadSize]
		truncatedNotice = "\n\n... [File content truncated: exceeded 2MB max read window] ..."
	}

	// Auto-detect and decode UTF-16LE / UTF-16BE with Byte Order Mark (BOM)
	content := decodeUTFString(data)
	lines := strings.Split(content, "\n")
	totalLines := len(lines)

	offset := helpers.GetIntArg(args, "offset", 1)
	if offset < 1 {
		offset = 1
	}

	maxLines := helpers.GetIntArg(args, "lines", 0)
	if maxLines == 0 {
		maxLines = helpers.GetIntArg(args, "limit", 0)
	}

	if offset > 1 || maxLines > 0 {
		startIdx := offset - 1
		if startIdx >= totalLines {
			return fmt.Sprintf("(File has %d lines, offset %d is past end of file)", totalLines, offset), true
		}
		endIdx := totalLines
		if maxLines > 0 && startIdx+maxLines < endIdx {
			endIdx = startIdx + maxLines
		}
		var sb strings.Builder
		for i := startIdx; i < endIdx; i++ {
			sb.WriteString(fmt.Sprintf("%4d | %s\n", i+1, lines[i]))
		}
		if endIdx < totalLines {
			sb.WriteString(fmt.Sprintf("... (%d more lines, total %d)\n", totalLines-endIdx, totalLines))
		}
		return sb.String(), true
	}

	if len(content) > 10000 {
		var sb strings.Builder
		limit := 200
		if totalLines < limit {
			limit = totalLines
		}
		for i := 0; i < limit; i++ {
			sb.WriteString(fmt.Sprintf("%4d | %s\n", i+1, lines[i]))
		}
		if totalLines > limit {
			sb.WriteString(fmt.Sprintf("... (%d more lines, total %d. Use offset and limit to read more)\n", totalLines-limit, totalLines))
		}
		return sb.String() + truncatedNotice, true
	}

	return content + truncatedNotice, true
}

// ── File Writer ──

var allowedWritePaths = []string{
	config.HomeDir() + "/",
	"/tmp/",
}

// ExecuteWriteFile writes a file
func ExecuteWriteFile(args map[string]interface{}) (string, bool) {
	path := helpers.GetStringArg(args, "path", "")
	content := helpers.GetStringArg(args, "content", "")
	if path == "" || content == "" {
		return "Error: 'path' and 'content' are required", false
	}

	if !isPathAllowed(path, allowedWritePaths) {
		return fmt.Sprintf("Error: path '%s' is not in allowed write directories", path), false
	}

	// Create parent directory
	dir := filepath.Dir(path)
	os.MkdirAll(dir, 0755)

	if err := WriteFileAtomic(path, []byte(content), 0644); err != nil {
		return fmt.Sprintf("Error writing file: %v", err), false
	}

	return fmt.Sprintf("File written: %s (%d bytes)", path, len(content)), true
}

// ── Directory Lister ──

// ExecuteListDir lists a directory
func ExecuteListDir(args map[string]interface{}) (string, bool) {
	path := helpers.GetStringArg(args, "path", ".")
	recursive := false
	if v, ok := args["recursive"]; ok {
		if b, ok := v.(bool); ok {
			recursive = b
		}
	}

	if !isPathAllowed(path, allowedReadPaths) {
		return fmt.Sprintf("Error: path '%s' is not in allowed directories", path), false
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Directory: %s\n\n", path))

	if recursive {
		visitedInodes := make(map[string]bool)
		filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			// Inode cycle detection to prevent ELOOP & infinite recursion on cross-symlinks
			realPath, rErr := filepath.EvalSymlinks(p)
			if rErr == nil {
				if visitedInodes[realPath] && info.IsDir() {
					return filepath.SkipDir
				}
				visitedInodes[realPath] = true
			}
			rel, _ := filepath.Rel(path, p)
			if info.IsDir() {
				sb.WriteString(fmt.Sprintf("📁 %s/\n", rel))
			} else {
				sb.WriteString(fmt.Sprintf("📄 %s (%d bytes)\n", rel, info.Size()))
			}
			return nil
		})
	} else {
		entries, err := os.ReadDir(path)
		if err != nil {
			return fmt.Sprintf("Error: %v", err), false
		}
		for _, entry := range entries {
			info, _ := entry.Info()
			if entry.IsDir() {
				sb.WriteString(fmt.Sprintf("📁 %s/\n", entry.Name()))
			} else if info != nil {
				sb.WriteString(fmt.Sprintf("📄 %s (%d bytes)\n", entry.Name(), info.Size()))
			}
		}
	}

	return helpers.TruncOutput(sb.String(), helpers.MaxToolOutput), true
}

// ── System Info ──

// ExecuteSystemInfo returns system information
func ExecuteSystemInfo(args map[string]interface{}) (string, bool) {
	infoType := helpers.GetStringArg(args, "type", "full")

	var result string
	switch infoType {
	case "cpu":
		out, _ := exec.Command("bash", "-c", "top -bn1 | head -5; echo '---'; lscpu | head -10").CombinedOutput()
		result = string(out)
	case "memory":
		out, _ := exec.Command("bash", "-c", "free -h; echo '---'; cat /proc/meminfo | head -5").CombinedOutput()
		result = string(out)
	case "disk":
		out, _ := exec.Command("bash", "-c", "df -h; echo '---'; lsblk").CombinedOutput()
		result = string(out)
	case "network":
		out, _ := exec.Command("bash", "-c", "ip addr show | head -30; echo '---'; ss -tuln | head -20").CombinedOutput()
		result = string(out)
	case "docker":
		out, _ := exec.Command("bash", "-c", "docker ps --format 'table {{.Names}}\t{{.Status}}\t{{.Ports}}' 2>/dev/null || echo 'Docker not available'").CombinedOutput()
		result = string(out)
	case "services":
		out, _ := exec.Command("bash", "-c", "systemctl list-units --type=service --state=running --no-pager | head -30").CombinedOutput()
		result = string(out)
	default:
		out, _ := exec.Command("bash", "-c", "echo '=== CPU ==='; top -bn1 | head -3; echo '=== Memory ==='; free -h; echo '=== Disk ==='; df -h | head -5; echo '=== Docker ==='; docker ps --format 'table {{.Names}}\t{{.Status}}' 2>/dev/null | head -10").CombinedOutput()
		result = string(out)
	}

	return helpers.TruncOutput(result, helpers.MaxToolOutput), true
}

// ── Send File ──

func formatSendFileSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%dB", size)
	} else if size < 1024*1024 {
		return fmt.Sprintf("%.1fKB", float64(size)/1024)
	} else if size < 1024*1024*1024 {
		return fmt.Sprintf("%.1fMB", float64(size)/(1024*1024))
	}
	return fmt.Sprintf("%.2fGB", float64(size)/(1024*1024*1024))
}

func ExecuteSendFile(args map[string]interface{}, chatID int64) (string, bool) {
	path := helpers.GetStringArg(args, "path", "")
	caption := helpers.GetStringArg(args, "caption", "")
	asDocument := helpers.GetBoolArg(args, "as_document", false)

	if path == "" {
		return "Error: 'path' argument is required", false
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Sprintf("Error resolving path: %v", err), false
	}

	if !isPathAllowed(absPath, allowedReadPaths) {
		return fmt.Sprintf("Error: path '%s' is not in allowed directories", absPath), false
	}

	info, err := os.Stat(absPath)
	if err != nil {
		return fmt.Sprintf("Error: %v", err), false
	}
	if info.IsDir() {
		return fmt.Sprintf("Error: '%s' is a directory. Use a file path or zip the directory first.", absPath), false
	}

	if info.Size() > 50*1024*1024 {
		return fmt.Sprintf("Error: file too large (%s) — maximum size for Telegram bot upload is 50MB", formatSendFileSize(info.Size())), false
	}

	if caption == "" {
		caption = fmt.Sprintf("📄 %s (%s)", filepath.Base(absPath), formatSendFileSize(info.Size()))
	}

	chatIDStr := fmt.Sprintf("%d", chatID)

	// Primary dispatcher: SendMedia
	if SendMedia != nil {
		ok, mediaType := SendMedia(chatIDStr, absPath, caption, asDocument)
		if !ok {
			return fmt.Sprintf("Error sending file via Telegram: %s", mediaType), false
		}
		return fmt.Sprintf("✅ File sent successfully to Telegram as %s: %s (%s)", mediaType, filepath.Base(absPath), formatSendFileSize(info.Size())), true
	}

	// Secondary fallback: SendDocumentBytes
	if SendDocumentBytes != nil {
		data, err := os.ReadFile(absPath)
		if err != nil {
			return fmt.Sprintf("Error reading file: %v", err), false
		}
		ok := SendDocumentBytes(chatIDStr, data, filepath.Base(absPath), caption)
		if !ok {
			return "Error sending file via Telegram", false
		}
		return fmt.Sprintf("✅ File sent successfully to Telegram: %s (%s)", filepath.Base(absPath), formatSendFileSize(info.Size())), true
	}

	// No sender wired (e.g. standalone test or non-Telegram mode)
	return fmt.Sprintf("✅ File verified for sending: %s (%s)", filepath.Base(absPath), formatSendFileSize(info.Size())), true
}

// isGuestLinuxRootfs detects if execution is occurring inside a guest Linux container (e.g. PRoot/Chroot)
// on Android. Native Termux has no /etc/os-release or /etc/debian_version (only $PREFIX/etc).
func isGuestLinuxRootfs() bool {
	if _, err := os.Stat("/etc/os-release"); err == nil {
		return true
	}
	if _, err := os.Stat("/etc/debian_version"); err == nil {
		return true
	}
	return false
}

// decodeUTFString decodes raw file bytes into a clean UTF-8 string,
// automatically handling UTF-16LE / UTF-16BE Byte Order Marks (BOM).
func decodeUTFString(b []byte) string {
	if len(b) >= 2 {
		// UTF-16LE BOM: FF FE
		if b[0] == 0xFF && b[1] == 0xFE {
			return decodeUTF16(b[2:], false)
		}
		// UTF-16BE BOM: FE FF
		if b[0] == 0xFE && b[1] == 0xFF {
			return decodeUTF16(b[2:], true)
		}
	}
	return string(b)
}

func decodeUTF16(b []byte, bigEndian bool) string {
	if len(b)%2 != 0 {
		b = b[:len(b)-1]
	}
	u16s := make([]uint16, len(b)/2)
	for i := 0; i < len(u16s); i++ {
		if bigEndian {
			u16s[i] = uint16(b[2*i])<<8 | uint16(b[2*i+1])
		} else {
			u16s[i] = uint16(b[2*i+1])<<8 | uint16(b[2*i])
		}
	}
	return string(utf16.Decode(u16s))
}

// cleanAbortedShm scans /dev/shm for abandoned scratch artifacts created during tests
func cleanAbortedShm() {
	entries, err := os.ReadDir("/dev/shm")
	if err != nil {
		return
	}
	for _, e := range entries {
		name := e.Name()
		if strings.HasPrefix(name, "scorp_") || strings.HasPrefix(name, "test_") {
			_ = os.Remove(filepath.Join("/dev/shm", name))
		}
	}
}
