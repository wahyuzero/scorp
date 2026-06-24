package tools

import (
	"scorp-agent/internal/helpers"
	"scorp-agent/config"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
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

	// Check for dangerous commands
	if IsDangerousCommand(command) {
		// Store for confirmation
		chatIDStr := fmt.Sprintf("%d", chatID)
		StorePendingConfirmation(chatIDStr, "shell", command, nil)
		return fmt.Sprintf("⚠️ DANGEROUS COMMAND DETECTED:\n%s\n\nPlease confirm execution.", command), false
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	// For simple commands without shell metacharacters, execute directly (avoids quoting issues)
	// Detect if command needs shell features (pipes, redirects, variables, etc.)
	// But allow parentheses in arguments (common in output/summaries)
	needsShell := needsShellExecution(command)
	
	var cmd *exec.Cmd
	if needsShell {
		cmd = exec.CommandContext(ctx, "bash", "-c", command)
	} else {
		// Split into command + args for direct execution (safer, no shell parsing)
		parts := strings.Fields(command)
		if len(parts) > 0 {
			cmd = exec.CommandContext(ctx, parts[0], parts[1:]...)
		} else {
			cmd = exec.CommandContext(ctx, "bash", "-c", command)
		}
	}
	
	output, err := cmd.CombinedOutput()
	result := string(output)

	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Sprintf("Command timed out after %ds:\n%s", timeout, helpers.TruncOutput(result, helpers.MaxToolOutput)), false
	}

	if err != nil {
		return fmt.Sprintf("Command failed: %v\nOutput:\n%s", err, helpers.TruncOutput(result, helpers.MaxToolOutput)), false
	}

	return helpers.TruncOutput(result, helpers.MaxToolOutput), true
}

// ── File Reader ──

var allowedReadPaths = []string{
	config.HomeDir() + "/",
	"/tmp/",
	"/etc/",
	"/var/log/",
	"/data/coolify/",
	"/opt/",
}

func isPathAllowed(path string, allowedPrefixes []string) bool {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	for _, prefix := range allowedPrefixes {
		if strings.HasPrefix(absPath, prefix) {
			return true
		}
	}
	return false
}

// ExecuteReadFile reads a file
func ExecuteReadFile(args map[string]interface{}) (string, bool) {
	path := helpers.GetStringArg(args, "path", "")
	if path == "" {
		return "Error: 'path' argument is required", false
	}

	if !isPathAllowed(path, allowedReadPaths) {
		return fmt.Sprintf("Error: path '%s' is not in allowed directories", path), false
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Sprintf("Error reading file: %v", err), false
	}

	content := string(data)
	maxLines := helpers.GetIntArg(args, "lines", 0)
	if maxLines > 0 {
		lines := strings.Split(content, "\n")
		if len(lines) > maxLines {
			content = strings.Join(lines[:maxLines], "\n") + fmt.Sprintf("\n... (%d more lines)", len(lines)-maxLines)
		}
	}

	return helpers.TruncOutput(content, helpers.MaxToolOutput), true
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

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
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
		filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
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

func ExecuteSendFile(args map[string]interface{}, chatID int64) (string, bool) {
	path := helpers.GetStringArg(args, "path", "")
	caption := helpers.GetStringArg(args, "caption", "")
	if path == "" {
		return "Error: 'path' argument is required", false
	}

	if !isPathAllowed(path, allowedReadPaths) {
		return fmt.Sprintf("Error: path '%s' is not in allowed directories", path), false
	}

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Sprintf("Error: %v", err), false
	}

	if info.Size() > 50*1024*1024 {
		return "Error: file too large (>50MB)", false
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Sprintf("Error reading file: %v", err), false
	}

	if caption == "" {
		caption = filepath.Base(path)
	}

	chatIDStr := fmt.Sprintf("%d", chatID)
	ok := SendDocumentBytes(chatIDStr, data, filepath.Base(path), caption)
	if !ok {
		return "Error sending file via Telegram", false
	}

	return fmt.Sprintf("File sent: %s (%d bytes)", filepath.Base(path), info.Size()), true
}
