package main

import (
"context"
"fmt"
"os"
"os/exec"
"path/filepath"
"strings"
"time"
)


// ── Shell Executor ──

func executeShell(args map[string]interface{}, chatID int64) (string, bool) {
	command := getStringArg(args, "command", "")
	if command == "" {
		return "Error: 'command' argument is required", false
	}

	timeout := getIntArg(args, "timeout", 30)
	if timeout > 300 {
		timeout = 300
	}

	// Check for dangerous commands
	if isDangerousCommand(command) {
		// Store for confirmation
		chatIDStr := fmt.Sprintf("%d", chatID)
		storePendingConfirmation(chatIDStr, "shell", command, nil)
		return fmt.Sprintf("⚠️ DANGEROUS COMMAND DETECTED:\n%s\n\nPlease confirm execution.", command), false
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "bash", "-c", command)
	output, err := cmd.CombinedOutput()
	result := string(output)

	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Sprintf("Command timed out after %ds:\n%s", timeout, truncOutput(result, maxToolOutput)), false
	}

	if err != nil {
		return fmt.Sprintf("Command failed: %v\nOutput:\n%s", err, truncOutput(result, maxToolOutput)), false
	}

	return truncOutput(result, maxToolOutput), true
}

// ── File Reader ──

var allowedReadPaths = []string{
	homeDir() + "/",
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

func executeReadFile(args map[string]interface{}) (string, bool) {
	path := getStringArg(args, "path", "")
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
	maxLines := getIntArg(args, "lines", 0)
	if maxLines > 0 {
		lines := strings.Split(content, "\n")
		if len(lines) > maxLines {
			content = strings.Join(lines[:maxLines], "\n") + fmt.Sprintf("\n... (%d more lines)", len(lines)-maxLines)
		}
	}

	return truncOutput(content, maxToolOutput), true
}

// ── File Writer ──

var allowedWritePaths = []string{
	homeDir() + "/",
	"/tmp/",
}

func executeWriteFile(args map[string]interface{}) (string, bool) {
	path := getStringArg(args, "path", "")
	content := getStringArg(args, "content", "")
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

func executeListDir(args map[string]interface{}) (string, bool) {
	path := getStringArg(args, "path", ".")
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

	return truncOutput(sb.String(), maxToolOutput), true
}

// ── System Info ──

func executeSystemInfo(args map[string]interface{}) (string, bool) {
	infoType := getStringArg(args, "type", "full")

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

	return truncOutput(result, maxToolOutput), true
}

// ── Send File ──

func executeSendFile(args map[string]interface{}, chatID int64) (string, bool) {
	path := getStringArg(args, "path", "")
	caption := getStringArg(args, "caption", "")
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
	ok := sendDocumentBytes(chatIDStr, data, filepath.Base(path), caption)
	if !ok {
		return "Error sending file via Telegram", false
	}

	return fmt.Sprintf("File sent: %s (%d bytes)", filepath.Base(path), info.Size()), true
}
