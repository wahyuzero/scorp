package agent

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"scorp-agent/internal/helpers"
)

// ──────────────────────────────────────────────
// Thinking Stream & Status Formatter
// ──────────────────────────────────────────────

// buildThinkingMessage builds the thinking stream display
func buildThinkingMessage(lines []string, elapsed time.Duration, done bool) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("🧠 <b>Agent</b> [%s]\n\n", elapsed.Round(time.Second)))

	for _, line := range lines {
		sb.WriteString(line + "\n")
	}

	if !done {
		sb.WriteString("\n⏳ <i>working...</i>")
	}

	return sb.String()
}

// shouldUpdateThinking returns true if we should update the thinking message.
// Batches updates: every 2 tool calls OR every 2 seconds (whichever comes first).
func shouldUpdateThinking(toolCount int, lastUpdate time.Time) bool {
	if toolCount%2 == 0 {
		return true
	}
	if time.Since(lastUpdate) > 2*time.Second {
		return true
	}
	return false
}

// toolDescription returns a human-readable description of a tool call
func toolDescription(tc ToolCall) string {
	switch tc.Name {
	case "shell":
		cmd := helpers.GetStringArg(tc.Args, "command", "")
		if len(cmd) > 80 {
			cmd = cmd[:80] + "..."
		}
		return fmt.Sprintf("🖥 shell: %s", cmd)
	case "read_file":
		return fmt.Sprintf("📖 read: %s", helpers.GetStringArg(tc.Args, "path", ""))
	case "write_file":
		return fmt.Sprintf("✏️ write: %s", helpers.GetStringArg(tc.Args, "path", ""))
	case "replace_file_content":
		return fmt.Sprintf("✂️ replace: %s", helpers.GetStringArg(tc.Args, "target_file", helpers.GetStringArg(tc.Args, "path", "")))
	case "patch":
		return fmt.Sprintf("✂️ patch: %s", helpers.GetStringArg(tc.Args, "path", ""))
	case "list_dir":
		return fmt.Sprintf("📂 list: %s", helpers.GetStringArg(tc.Args, "path", "."))
	case "system_info":
		return fmt.Sprintf("ℹ️ sysinfo: %s", helpers.GetStringArg(tc.Args, "type", "full"))
	case "send_file":
		return fmt.Sprintf("📤 send: %s", helpers.GetStringArg(tc.Args, "path", ""))
	case "read_url":
		return fmt.Sprintf("🌐 read_url: %s", helpers.GetStringArg(tc.Args, "url", ""))
	case "web_fetch":
		return fmt.Sprintf("🌐 fetch: %s", helpers.GetStringArg(tc.Args, "url", ""))
	case "web_search":
		return fmt.Sprintf("🔍 search: %s", helpers.GetStringArg(tc.Args, "query", ""))
	case "memory":
		action := helpers.GetStringArg(tc.Args, "action", "")
		key := helpers.GetStringArg(tc.Args, "key", "")
		return fmt.Sprintf("🧠 memory.%s(%s)", action, key)
	case "search_code":
		return fmt.Sprintf("🔎 search_code: %s in %s", helpers.GetStringArg(tc.Args, "pattern", "?"), helpers.GetStringArg(tc.Args, "path", "."))
	case "git":
		return fmt.Sprintf("📦 git.%s (%s)", helpers.GetStringArg(tc.Args, "command", "?"), helpers.GetStringArg(tc.Args, "repo", "."))
	case "http":
		return fmt.Sprintf("📡 http.%s → %s", helpers.GetStringArg(tc.Args, "method", "GET"), helpers.GetStringArg(tc.Args, "url", "?"))
	case "log":
		return fmt.Sprintf("📋 log.%s(%s)", helpers.GetStringArg(tc.Args, "source", "?"), helpers.GetStringArg(tc.Args, "target", "?"))
	case "sql":
		query := helpers.GetStringArg(tc.Args, "query", "?")
		if len(query) > 50 {
			query = query[:47] + "..."
		}
		return fmt.Sprintf("🗄 sql: %s", query)
	case "process":
		return fmt.Sprintf("⚙️ process.%s", helpers.GetStringArg(tc.Args, "action", "?"))
	case "browser":
		action := helpers.GetStringArg(tc.Args, "action", "")
		if action == "goto" {
			return fmt.Sprintf("🌐 browser→%s", helpers.GetStringArg(tc.Args, "url", ""))
		}
		return fmt.Sprintf("🌐 browser.%s", action)
	case "analyze_image":
		return fmt.Sprintf("👁 analyze_image: %s", helpers.GetStringArg(tc.Args, "path", "?"))
	case "mcp_tool":
		server := helpers.GetStringArg(tc.Args, "server", "")
		tool := helpers.GetStringArg(tc.Args, "tool", "")
		return fmt.Sprintf("🔌 mcp: %s.%s", server, tool)
	case "delegate":
		task := helpers.GetStringArg(tc.Args, "task", "?")
		if len(task) > 60 {
			task = task[:60] + "..."
		}
		return fmt.Sprintf("🤖 delegate: %s", task)
	case "sop":
		return fmt.Sprintf("📋 sop.%s", helpers.GetStringArg(tc.Args, "action", "list"))
	default:
		return fmt.Sprintf("🔧 %s", tc.Name)
	}
}

// toolCallSignature returns a string key identifying a tool call by name + args.
func toolCallSignature(tc ToolCall) string {
	argsJSON, err := json.Marshal(tc.Args)
	if err != nil {
		return tc.Name
	}
	return fmt.Sprintf("%s:%s", tc.Name, string(argsJSON))
}
