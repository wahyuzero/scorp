package agent

import (
	"encoding/json"
	"fmt"
	"log"
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

## LANGUAGE (CRITICAL)
- ALWAYS respond in the SAME LANGUAGE the user's message is written in. Indonesian prompt → Indonesian reply, English prompt → English reply. NEVER switch languages mid-conversation unless the user does.

## HIGH-THROUGHPUT EXECUTION & PARALLEL BATCHING (CRITICAL)
- TURBO SPEED (1-2 TURNS ONLY): For tasks that involve creating and running scripts/programs, prefer writing the script and running it via 'shell' in a single command pipeline (e.g. 'mkdir -p /path && cat << "EOF" > /path/script.py ... EOF && python3 /path/script.py') or emit both tool calls together. Never create empty folders or perform multiple separate inspection turns when a direct execution can accomplish the task.
- DIRECT ACTIONS (Fast-Path): For focused, direct, or atomic tasks (creating/editing files, running commands/scripts, system diagnostics, QA, single bugs): DO NOT call 'task_plan'. Directly emit the action tools immediately.
- MULTI-PHASE PROJECTS ONLY: Call 'task_plan' ONLY for massive, multi-phase projects (e.g. migrating 10+ modules, multi-service architecture overhauls).
- AUTO-COMPLETION: When all necessary action tools have executed successfully and you have the final results, conclude immediately by outputting your final comprehensive summary directly in text, or call 'complete_task' with the result. Do not perform redundant re-reading or directory listing once the command output already proves success. Never attempt to rewrite or repeat work that was already successfully edited.
- REPORTING ERRORS DIRECTLY: If a command was specifically intended to verify an error, trigger an expected non-zero exit code, or demonstrate a failure (e.g., checking invalid database files, nonexistent paths, or broken syntax), report the captured error output immediately as the proof. Do NOT treat expected test errors as an unexpected failure that requires retrying or debugging.

## MULTI-STEP TASKS & EXECUTION RULES (CRITICAL)
- ACTION-FIRST: If an action is required, EMIT THE TOOL CALL IMMEDIATELY.
- NATIVE TOOL CALLING ONLY: invoke tools exclusively through native function calling (tool_calls). NEVER write tool invocations as text (no DSML, no <tool_call>, no XML tags).
- NEVER narrate future actions in text (e.g. do NOT say "I will now create...", "Now running..."). Output the tool call directly.
- SILENT INTERMEDIATE STEPS: Any text output between tool calls is treated as internal thought.
- Do NOT stop midway. If more verification or actions are needed, call the next tools. If finished, deliver the final verified report.
- When asked to check, run, fix, compile, or test — call the appropriate tool.
- Sequence for engineering tasks:
  1. Inspect/Prepare if needed (read_file, list_dir).
  2. Modify/Create code (write_file / replace_file_content) and Execute (shell) — batch together if possible!
  3. Verify outputs and deliver the verified final report.

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
	skillsIndex := skills.FormatSkillsIndexForSystemPrompt()
	if skillsIndex != "" {
		sb.WriteString(skillsIndex + "\n")
	}

	activeSkillsCtx := skills.GetActiveSkillsContext()
	if activeSkillsCtx != "" {
		sb.WriteString(activeSkillsCtx + "\n")
	}

	memSummary := tools.GetMemorySummary()
	if memSummary != "" {
		sb.WriteString("\n## PERSISTENT MEMORY\n" + memSummary + "\n")
	}

	// Durable memory file (P1.7): decisions & state extracted at task
	// completions — the agent's long-term notes across sessions.
	if durable := ReadMemoryMD(); durable != "" {
		sb.WriteString("\n## DURABLE MEMORY (MEMORY.md)\n" + durable + "\n")
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

// ExecuteTool runs a tool via the registry with deny-rule, autonomy & sandbox enforcement
func ExecuteTool(tc ToolCall, chatID int64) (string, bool) {
	// 0. Deny rules — the absolute outer layer (P0.2). Evaluated before
	// autonomy/confirmation logic so they hold in readonly, supervised, YOLO,
	// and on user-confirmed resumes alike.
	if blocked, reason := config.CheckDenyRules(tc.Name, tc.Args); blocked {
		return "🚫 " + reason, false
	}

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

	// 3. PreToolUse hooks (P3.12): user-configured enforcement at the single
	// choke point. Fires only when the call is actually about to execute —
	// deny rules, autonomy and path checks remain the outer layers.

	// 2.5 Auto-mode enforcement (P3.13) for paths that bypass the agent loops
	// (subagents). The loops classify before calling ExecuteTool and preset
	// tc.AutoDecision; a preset decision is trusted. Un-preset calls are
	// graded here with fail-closed semantics — this path has no confirmation
	// channel, so an "ask" degrades to a denial.
	if config.GetAutonomyLevel() == config.AutonomyAuto {
		if tc.AutoDecision == "" {
			decision, reason, source := PermissionDecision(tc.Name, tc.Args)
			switch decision {
			case AutoAllow:
				tc.AutoDecision = "auto:" + source
			case AutoDeny:
				log.Printf("[auto] DENY %s (subagent/direct path): %s", tc.Name, reason)
				return "⛔ Denied by auto-mode classifier: " + reason, false
			default:
				log.Printf("[auto] ASK→DENY %s (no confirmation channel on this path): %s", tc.Name, reason)
				return "⛔ Auto-mode requires human confirmation for this action, and this execution path has no confirmation channel. Escalate to the parent agent/user instead.", false
			}
		}
	}

	// 3.9 Loop/human-approved auto decisions authorize the exec-layer danger
	// gate for shell: the classifier (or the human on the auto-ask path) has
	// already approved THIS exact call, and re-prompting inside ExecuteShell
	// would make SCORP_AUTO_ALLOW unreachable end-to-end. AutoDecision is
	// preset by the loop/resume code only (json:"-" — a model cannot forge
	// it), and deterministic denies never reach this point.
	if config.GetAutonomyLevel() == config.AutonomyAuto && tc.AutoDecision != "" && tc.Name == "shell" {
		if tc.Args == nil {
			tc.Args = make(map[string]interface{})
		}
		tc.Args["confirmed"] = true
	}

	hookCtx, hookBlocked, hookReason := tools.RunPreToolUseHooks(tc.Name, tc.Args, chatID)
	if hookBlocked {
		return "🪝 " + hookReason, false
	}

	// 4. Execute Tool
	out, ok := registry.ExecuteToolByName(tc.Name, tc.Args, chatID)

	// 5. Outbound Secret Redaction (prevent API key leakage)
	out = tools.RedactSecrets(out)

	// 6. Record Cryptographic Receipt (ZeroClaw Parity) — before hook context
	// is appended, so receipts capture the tool's own output. The auto-mode
	// decision (P3.13) travels with the receipt when present.
	var extraMeta []map[string]string
	if tc.AutoDecision != "" {
		extraMeta = append(extraMeta, map[string]string{"auto_decision": tc.AutoDecision})
	}
	tools.RecordToolReceipt(tc.Name, tc.Args, out, ok, extraMeta...)

	// 7. PostToolUse hooks + pre-hook additional context (P3.12): never block,
	// stdout is surfaced to the model.
	if postCtx := tools.RunPostToolUseHooks(tc.Name, tc.Args, out, ok, chatID); postCtx != "" {
		out += "\n\n🪝 " + postCtx
	}
	if hookCtx != "" {
		out += "\n\n🪝 " + hookCtx
	}

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
