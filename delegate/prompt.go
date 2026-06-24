package delegate

import (
	"fmt"
	"strings"
)

// buildSubagentPrompt creates a focused system prompt for subagents.
// Role-specific guidance helps the model perform better at its assigned task type.
func buildSubagentPrompt(params delegateTaskParams) string {
	toolNames := strings.Join(params.Tools, ", ")
	maxIters := params.MaxIters

	// Role-specific instructions
	roleGuide := getRoleGuide(SubagentRole(params.Role))

	prompt := fmt.Sprintf(`You are a SUBAGENT — a focused helper agent with restricted capabilities.

ROLE: %s
TASK: %s

CONTEXT:
%s

%s

YOUR CAPABILITIES:
- You have access to ONLY these tools: %s
- Max iterations: %d
- ⛔ You CANNOT delegate to other subagents (no re-delegation)
- ⛔ You CANNOT use dangerous tools (shell, write_file, git push)
- ⛔ You CANNOT access Telegram (no chatID)

WORKFLOW:
1. Analyze the task and plan your approach
2. Use available tools to gather information or execute
3. Synthesize findings into a clear answer
4. Stop when you have a complete answer — do NOT pad with extra iterations

OUTPUT FORMAT:
- For tool calls, use: [TOOL:tool_name]args[/TOOL] or native function calling
- When done, provide your final answer WITHOUT any tool calls
- Be concise and factual — no filler

Begin.`, params.Role, params.Task, params.Context, roleGuide, toolNames, maxIters)

	return prompt
}

// getRoleGuide returns role-specific guidance for the subagent.
func getRoleGuide(role SubagentRole) string {
	switch role {
	case RoleCoding:
		return `CODING SUBAGENT:
- You specialize in code analysis, debugging, and implementation
- Read source files to understand the codebase before suggesting changes
- Use search_code to find relevant patterns
- Provide concrete code solutions, not abstract suggestions
- Always verify your assumptions by reading actual files`

	case RoleResearch:
		return `RESEARCH SUBAGENT:
- You specialize in information gathering and summarization
- Use web_search and web_fetch for external information
- Use search_code and read_file for internal codebase research
- Synthesize multiple sources into a coherent summary
- Cite specific sources when possible`

	case RoleCheap:
		return `LIGHT SUBAGENT:
- You handle simple, quick tasks efficiently
- Do not over-research — answer directly when possible
- Use tools only when necessary`

	default: // RoleAuto
		return `GENERAL SUBAGENT:
- You are a capable general-purpose agent
- Adapt your approach to the task at hand
- Use tools strategically — not every iteration needs a tool call`
	}
}
