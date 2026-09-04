package agent

import (
	"encoding/json"
	"fmt"
	"scorp-agent/config"
	"scorp-agent/mcp"
	"scorp-agent/models"
	"scorp-agent/registry"
	"scorp-agent/skills"
	"scorp-agent/tools"
	"strings"
)

// ──────────────────────────────────────────────
// Tool Definitions & System Prompt (with Prompt Caching)
// ──────────────────────────────────────────────

// staticSystemPrefix provides a constant, byte-stable system prompt prefix.
// Modern LLMs (Gemini, Claude 3.5, DeepSeek V3) cache this prefix, saving 75-90% of token costs.
const staticSystemPrefix = `You are Scorp Agent (v2.0) — an intelligent, ultra-fast, and lightweight AI coding & automation agent.

## IDENTITY
You are a versatile AI agent built for programming, deep research, system administration, automation, and DevOps.
You run natively with ultra-low memory footprint on Linux VPS, Termux Android, and edge devices.
Communication style: direct, efficient, technically precise, no fluff. Respond in the same language as the user.
You are interacting directly with the user in this active session. All your textual responses are displayed directly to the user in their interface. When asked to report, summarize, or present findings, output them directly here — never claim you lack messaging credentials or cannot deliver results to this session.

## FILE EDITING & CODING (CRITICAL)
- When modifying existing files, ALWAYS use 'replace_file_content' (or 'patch') to perform surgical diffs.
- DO NOT rewrite whole files with 'write_file' if only editing parts of the file.
- Use 'write_file' ONLY when creating brand new files or complete rewrites.
- Chunk replacement saves massive tokens, avoids token limits, and executes in sub-seconds.

## WEB BROWSING & RESEARCH (ULTRA-LOW RAM)
- For reading documentation, articles, GitHub files, API specs, and websites, ALWAYS prefer 'read_url'.
- 'read_url' uses the Zero-RAM reader engine (<5MB RAM) and outputs clean Markdown.
- Use the heavy 'browser' tool ONLY when you genuinely need to click UI buttons, fill interactive forms, or take graphical screenshots.

## MULTI-STEP TASKS & ACTION-FIRST PROTOCOL (CRITICAL)
- ACTION-FIRST: If an action is required, EMIT THE TOOL CALL IMMEDIATELY.
- NEVER narrate future actions in text (e.g. do NOT say "I will now delete...", "Now running the script...", "Sekarang saya akan buat...").
- SILENT INTERMEDIATE STEPS: Do not output chit-chat or explanatory text between sequential tool calls.
- ONLY output final conversational text when ALL requested steps from the user are 100% finished and verified.
- When a user asks you to check, run, fix, search, or monitor something — you MUST call the appropriate tool.
- Most tasks REQUIRE multiple tool calls in sequence. Do NOT stop after one tool call unless the task is truly complete:
  1. Search / inspect first (search_code, list_dir, read_file).
  2. Make surgical modifications (replace_file_content).
  3. Verify changes with tests or build commands (shell).
  4. Continue calling tools until you have verified results and a COMPLETE answer.

## FORBIDDEN
- NEVER substitute plausible-looking fabricated output for results you couldn't actually produce.
- Reporting a blocker honestly is always better than inventing a result.
- For dangerous commands (rm -rf /, mkfs, dd, DROP TABLE, systemctl stop scorp-agent, etc.), ask for confirmation first.

## BROWSER WORKFLOW (When using headless browser)
### Navigation
1. browser goto → browser snapshot (to see interactive elements)
2. Always snapshot after goto — NEVER type or click blind

### Form Filling & Login
1. browser type fills ONE field and reports available submit buttons
2. After filling, READ the submit button info from the result
3. browser click the submit button that was reported

### Loop Prevention
- If you type the same thing and click the same button TWICE with no change, STOP.
- The URL did NOT change after your click → login FAILED or page didn't respond.
- Do NOT repeat the same action — try a different approach or report the failure.
- Take EXACTLY ONE screenshot at the very end of browser tasks if requested.
`

// getAgentSystemPrompt returns the complete system prompt.
// The layout is ordered: Static Prefix -> Repository Map -> Tools -> Dynamic Context
// to maximize prefix prompt caching hits.
func getAgentSystemPrompt(chatID ...int64) string {
	var sb strings.Builder

	// 1. Static Prefix (Fixed & constant)
	sb.WriteString(staticSystemPrefix)

	// 2. Repository Map Prefix (Cached workspace overview)
	repoMap := GetRepoMap()
	if repoMap != "" {
		sb.WriteString("\n## WORKSPACE STRUCTURE\n" + repoMap + "\n")
	}

	// 3. MCP Tools List (if any connected)
	mcpToolList := mcp.GetMCPTools()
	if len(mcpToolList) > 0 {
		sb.WriteString("\n## AVAILABLE MCP TOOLS\n")
		for _, t := range mcpToolList {
			sb.WriteString(fmt.Sprintf("- %s.%s: %s\n", t.ServerName, t.Name, t.Description))
		}
	}

	// 4. Dynamic Context (Skills & Memory at tail to keep prefix stable)
	skillDesc := skills.GetPromptForMessage("")
	if skillDesc != "" {
		sb.WriteString("\n## ACTIVE SKILLS\n" + skillDesc + "\n")
	}

	memSummary := tools.GetMemorySummary()
	if memSummary != "" {
		sb.WriteString("\n## PERSISTENT MEMORY\n" + memSummary + "\n")
	}

	// 5. Interface Environment Context (Tail only)
	if len(chatID) > 0 && chatID[0] != 0 {
		sb.WriteString("\n## SESSION CONTEXT\n- Active interface: Telegram Messenger\n")
	} else {
		sb.WriteString("\n## SESSION CONTEXT\n- Active interface: Terminal / CLI Console\n")
	}

	return sb.String()
}

// ToolCall type alias
type ToolCall = models.ToolCall

// Dangerous Command Filtering moved to agent/safety.go (argument tokenizer)

// ──────────────────────────────────────────────
// Tool Executors
// ──────────────────────────────────────────────

const maxToolOutput = 3000

// ExecuteTool runs a tool via the registry with autonomy & sandbox enforcement
func ExecuteTool(tc ToolCall, chatID int64) (string, bool) {
	// 1. Check Autonomy Level permission (ReadOnly mode restrictions)
	if allowed, reason := config.IsToolAllowed(tc.Name); !allowed {
		return "⚠️ " + reason, false
	}

	// 2. Check Sandbox Path Restrictions
	for _, key := range []string{"path", "target_file", "file"} {
		if pathVal, ok := tc.Args[key].(string); ok && pathVal != "" {
			if restricted, reason := config.IsPathRestricted(pathVal); restricted {
				return "🛡️ " + reason, false
			}
		}
	}

	// 3. Execute Tool
	out, ok := registry.ExecuteToolByName(tc.Name, tc.Args, chatID)

	// 4. Outbound Secret Redaction (PicoClaw Parity: prevent API key leakage)
	out = tools.RedactSecrets(out)

	// 5. Record Cryptographic Receipt (ZeroClaw Parity)
	tools.RecordToolReceipt(tc.Name, tc.Args, out, ok)

	return out, ok
}

// FormatToolResult formats a tool result for the conversation
func FormatToolResult(tc ToolCall, result string, ok bool) string {
	status := "SUCCESS"
	if !ok {
		status = "FAILED"
	}

	argsJSON, _ := json.Marshal(tc.Args)
	return fmt.Sprintf("[%s] %s(%s)\n%s", status, tc.Name, string(argsJSON), result)
}
