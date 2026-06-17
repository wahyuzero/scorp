package main

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
)

// ──────────────────────────────────────────────
// Tool Definitions & System Prompt
// ──────────────────────────────────────────────

// getAgentSystemPrompt returns the system prompt, dynamically including MCP tools
func getAgentSystemPrompt() string {
	prompt := `You are Scorp, a powerful AI agent running on a VPS (Ubuntu 24.04, 8 vCPU, 24GB RAM).
You can execute tools to help the user. Always respond in the same language the user uses.

## CRITICAL: Always USE tools for multi-step tasks!
When a user asks you to check, run, fix, search, or monitor something — you MUST call the appropriate tool.
NEVER just describe what you would do. Actually DO it by calling tools.
For example, if asked to "check disk space", call the shell tool with "df -h" — don't just say you would run it.

## MULTI-STEP REASONING (MANDATORY)
Most tasks require MULTIPLE tool calls in sequence. Do NOT stop after one tool call unless the task is truly complete.
- If you need to search → call search tool → analyze results → call another tool if needed → repeat
- If you fetch a page → read it → if it references other files/links → fetch those too
- If a command output suggests next steps → execute them
- Continue calling tools until you have a COMPLETE answer for the user.

## Tool calling
You have tools available via function calling. Call them directly.
You can call multiple tools in one response.
After receiving tool results, analyze them and decide: continue with more tools OR give final answer.
`

	// Add MCP tool descriptions if available
	mcpToolList := GetMCPTools()
	if len(mcpToolList) > 0 {
		var mcpDesc strings.Builder
		for _, t := range mcpToolList {
			mcpDesc.WriteString(fmt.Sprintf("- %s.%s: %s\n", t.ServerName, t.Name, t.Description))
		}
		prompt += "\n### Available MCP Tools\n" + mcpDesc.String() + "\n"
	}

	// Add skill context
	skillDesc := getSkillPromptForMessage("")
	if skillDesc != "" {
		prompt += "\n### Skills\n" + skillDesc + "\n"
	}

	// Add memory summary (auto-injected)
	memSummary := getMemorySummary()
	if memSummary != "" {
		prompt += "\n### Memory (Auto-injected)\n" + memSummary + "\n"
	}

	prompt += `
## Important Rules
1. Always use tools when asked to perform system tasks - don't guess or make up information
2. For dangerous commands (rm -rf, format, DROP TABLE, etc.), ask for confirmation first
3. Keep responses concise and in the user's language
4. When a task is complete, provide a clear summary
5. If a command fails, explain why and suggest alternatives
6. Use [REMEMBER:key:value] to save important information to memory
7. You can use multiple tools in a single response

## BROWSER WORKFLOW (CRITICAL — read this carefully)
When using the browser tool, follow these rules STRICTLY:

### Navigation
1. browser goto → browser snapshot (to see interactive elements)
2. Always snapshot after goto — NEVER type or click blind

### Form Filling & Login
1. browser type fills ONE field and reports available submit buttons
2. After filling, READ the submit button info from the result
3. browser click the submit button that was reported

### After EVERY click — analyze the result!
The click result tells you:
- "URL: X → Y" — did the URL change? If yes, you navigated somewhere new
- "⚠️ URL unchanged" — login may have FAILED, check the page text for errors
- "📋 Elements:" — what interactive elements are now available
- "📄 Page text:" — what the page says after clicking

### Screenshot
- ALWAYS take a screenshot as the LAST step of any browser task
- browser action=screenshot
- If you were asked to screenshot something, the task is NOT done until you call screenshot

### Loop Prevention
- If you type the same thing and click the same button TWICE with no change, STOP
- The URL did NOT change after your click → login FAILED or page didn't respond
- Do NOT repeat the same action — try a different approach or report the failure
- Check the page text for error messages after each click
`

	return prompt
}

// ──────────────────────────────────────────────
// Tool Call Parser
// ──────────────────────────────────────────────

type ToolCall struct {
	Name string                 `json:"name"`
	Args map[string]interface{} `json:"args"`
}

var toolCallRe = regexp.MustCompile(`<tool_call>(.*?)</tool_call>`)

// parseToolCalls extracts tool calls from LLM response
func parseToolCalls(text string) ([]ToolCall, string) {
	matches := toolCallRe.FindAllStringSubmatch(text, -1)
	if len(matches) == 0 {
		return nil, text
	}

	var calls []ToolCall
	cleanText := text

	for _, m := range matches {
		cleanText = strings.Replace(cleanText, m[0], "", 1)
		var tc ToolCall
		if err := json.Unmarshal([]byte(m[1]), &tc); err != nil {
			log.Printf("[scorp] Failed to parse tool call: %v", err)
			continue
		}
		calls = append(calls, tc)
	}

	return calls, strings.TrimSpace(cleanText)
}

// ──────────────────────────────────────────────
// Dangerous Command Detection
// ──────────────────────────────────────────────

var dangerousPatterns []string

func init() {
	raw := []string{
		"rm -rf /", "rm -rf /*", "mkfs", "dd if=", ":(){ :|:& };:",
		"drop table", "drop database", "delete from",
		"kill -9", "killall", "pkill",
		"systemctl stop", "systemctl disable",
		"apt remove", "apt purge", "pip uninstall",
		"docker rm", "docker rmi", "docker prune",
		"docker-compose down", "docker compose down",
		"> /dev/", "chmod 777",
	}
	dangerousPatterns = make([]string, len(raw))
	for i, p := range raw {
		dangerousPatterns[i] = strings.ToLower(p)
	}
}

func isDangerousCommand(cmd string) bool {
	lower := strings.ToLower(cmd)
	for _, p := range dangerousPatterns {
		if strings.Contains(lower, p) {
			return true
		}
	}
	return false
}

// ──────────────────────────────────────────────
// Tool Executors
// ──────────────────────────────────────────────

const maxToolOutput = 3000

func truncOutput(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "\n... (truncated)"
}

// executeTool runs a tool via the registry
func executeTool(tc ToolCall, chatID int64) (string, bool) {
	return executeToolByName(tc.Name, tc.Args, chatID)
}

func getStringArg(args map[string]interface{}, key, defaultVal string) string {
	if v, ok := args[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return defaultVal
}

func getIntArg(args map[string]interface{}, key string, defaultVal int) int {
	if v, ok := args[key]; ok {
		switch n := v.(type) {
		case float64:
			return int(n)
		case int:
			return n
		}
	}
	return defaultVal
}

func getBoolArg(args map[string]interface{}, key string, defaultVal bool) bool {
	if v, ok := args[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return defaultVal
}

func getFloatArg(args map[string]interface{}, key string, defaultVal float64) float64 {
	if v, ok := args[key]; ok {
		switch n := v.(type) {
		case float64:
			return n
		case float32:
			return float64(n)
		case int:
			return float64(n)
		case int64:
			return float64(n)
		}
	}
	return defaultVal
}

func getStringSliceArg(args map[string]interface{}, key string) []string {
	if v, ok := args[key]; ok {
		if arr, ok := v.([]interface{}); ok {
			result := make([]string, 0, len(arr))
			for _, item := range arr {
				if s, ok := item.(string); ok {
					result = append(result, s)
				}
			}
			return result
		}
	}
	return nil
}

func getInt64Arg(args map[string]interface{}, key string, defaultVal int64) int64 {
	if v, ok := args[key]; ok {
		switch n := v.(type) {
		case float64:
			return int64(n)
		case int64:
			return n
		case int:
			return int64(n)
		case json.Number:
			if i, err := n.Int64(); err == nil {
				return i
			}
		}
	}
	return defaultVal
}
