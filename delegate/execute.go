package delegate

import (
	"scorp-agent/internal/helpers"
	"scorp-agent/registry"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

// ──────────────────────────────────────────────
// Subagent / Delegation System v2
//
// Inspired by OhMyOpenAgent orchestration patterns:
// - Category-based model routing (role → model)
// - Parallel execution (max 5 concurrent)
// - No-re-delegation rule (subagents can't spawn subagents)
// - Per-agent model override
// - Up to 20 iterations per subagent
// ──────────────────────────────────────────────

const (
	MaxParallelSubagents = 5
	MaxSubagentIters     = 15
	DefaultSubagentIters = 10
)

// SubagentRole determines which model a subagent uses.
// Inspired by OmO's category-based routing.
type SubagentRole string

const (
	RoleAuto         SubagentRole = "auto"         // Use agent model (default)
	RoleCoding       SubagentRole = "coding"       // Powerful coding model
	RoleResearch     SubagentRole = "research"     // Fast model for search/summarize
	RoleCheap        SubagentRole = "cheap"        // Cheapest model
	RoleOrchestrator SubagentRole = "orchestrator" // Can spawn own workers (re-delegation allowed)
)

// roleToTaskType maps a subagent role to a model routing task type.
func roleToTaskType(role SubagentRole) string {
	switch role {
	case RoleCoding:
		return "coding"
	case RoleResearch:
		return "research"
	case RoleCheap:
		return "chat"
	case RoleOrchestrator:
		return "agent" // orchestrator uses full agent model
	default:
		return "agent"
	}
}

// delegateTaskParams defines the input for a single subagent task.
type delegateTaskParams struct {
	Task      string   `json:"task"`       // What the subagent should do
	Context   string   `json:"context"`    // Background info / constraints
	Tools     []string `json:"tools"`      // Allowed tools (empty = read-only default)
	MaxIters  int      `json:"max_iters"`  // Max iterations (default 10, max 20)
	ReturnRaw bool     `json:"return_raw"` // Return raw results instead of formatted
	Role       string   `json:"role"`        // auto/coding/research/cheap (model routing)
	Model      string   `json:"model"`       // Explicit model name override
	ACPCommand string   `json:"acp_command"` // External CLI binary (e.g., "claude", "codex")
	ACPArgs    []string `json:"acp_args"`    // CLI args (default: ["--acp", "--stdio"])
}

// delegateResult is the result of a subagent execution
type delegateResult struct {
	SubagentID  string
	Task        string
	Role        string
	ModelUsed   string
	Status      string // "completed", "failed", "timeout"
	Result      string
	ToolsUsed   []string
	Iterations  int
	Duration    time.Duration
	Error       string
}

// ──────────────────────────────────────────────
// Single-task delegate (backward compatible)
// ──────────────────────────────────────────────

func ExecuteDelegate(args map[string]interface{}) (string, bool) {
	params := ParseDelegateParams(args)

	if params.Task == "" {
		return "Error: 'task' is required", false
	}

	// Default tools if none specified
	if len(params.Tools) == 0 {
		params.Tools = DefaultSubagentTools()
	}

	// Validate tools
	isOrch := SubagentRole(params.Role) == RoleOrchestrator
	params.Tools = ValidateSubagentTools(params.Tools, isOrch)

	log.Printf("[delegate] Starting subagent: role=%s model=%s iters=%d task=%s",
		params.Role, params.Model, params.MaxIters, helpers.TruncateStr(params.Task, 60))

	// Route to ACP subprocess if acp_command is set
	var result delegateResult
	if params.ACPCommand != "" {
		log.Printf("[delegate] ACP mode: command=%s args=%v", params.ACPCommand, params.ACPArgs)
		result = runSubagentACP(params)
	} else {
		result = runSubagent(params)
	}

	return FormatDelegateResult(params, result)
}

// ParseDelegateParams extracts and validates parameters from args map.
func ParseDelegateParams(args map[string]interface{}) delegateTaskParams {
	p := delegateTaskParams{
		Task:       helpers.GetStringArg(args, "task", ""),
		Context:    helpers.GetStringArg(args, "context", ""),
		Tools:      helpers.GetStringSliceArg(args, "tools"),
		MaxIters:   helpers.GetIntArg(args, "max_iters", DefaultSubagentIters),
		ReturnRaw:  helpers.GetBoolArg(args, "return_raw", false),
		Role:       helpers.GetStringArg(args, "role", "auto"),
		Model:      helpers.GetStringArg(args, "model", ""),
		ACPCommand: helpers.GetStringArg(args, "acp_command", ""),
		ACPArgs:    helpers.GetStringSliceArg(args, "acp_args"),
	}

	// Clamp iterations
	if p.MaxIters < 1 {
		p.MaxIters = DefaultSubagentIters
	}
	if p.MaxIters > MaxSubagentIters {
		p.MaxIters = MaxSubagentIters
	}

	// Validate role
	switch SubagentRole(p.Role) {
	case RoleAuto, RoleCoding, RoleResearch, RoleCheap, RoleOrchestrator:
		// valid
	default:
		p.Role = string(RoleAuto)
	}

	return p
}

// DefaultSubagentTools returns the standard read-only toolset.
func DefaultSubagentTools() []string {
	return []string{"read_file", "search_code", "system_info", "log", "web_fetch", "web_search", "list_dir", "index_search"}
}

// ValidateSubagentTools filters out non-existent tools and blocks dangerous tools.
// S08/GC4: Orchestrator role can use delegate/delegate_batch for re-delegation.
func ValidateSubagentTools(tools []string, isOrchestrator bool) []string {
	// Tools that subagents can NEVER have access to (unless orchestrator)
	blocked := map[string]bool{
		"delegate":       true, // Blocked unless orchestrator
		"delegate_batch": true, // Blocked unless orchestrator
	}

	valid := []string{}
	for _, t := range tools {
		if blocked[t] && !isOrchestrator {
			continue
		}
		if _, ok := registry.GetTool(t); ok {
			valid = append(valid, t)
		}
	}

	if len(valid) == 0 {
		valid = []string{"read_file", "search_code", "system_info"}
	}

	return valid
}

// FormatDelegateResult formats the output for the caller.
func FormatDelegateResult(params delegateTaskParams, result delegateResult) (string, bool) {
	if result.Error != "" {
		return fmt.Sprintf("❌ Subagent failed: %s\n\n%s", result.Error, result.Result), false
	}

	if params.ReturnRaw {
		return result.Result, true
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("✅ <b>Subagent Completed</b>\n\n"))
	sb.WriteString(fmt.Sprintf("📋 <b>Task:</b> %s\n", helpers.TruncateStr(params.Task, 80)))
	sb.WriteString(fmt.Sprintf("🎭 <b>Role:</b> %s", result.Role))
	if result.ModelUsed != "" {
		sb.WriteString(fmt.Sprintf(" | 🤖 <b>Model:</b> %s", result.ModelUsed))
	}
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("⏱ <b>Duration:</b> %s | <b>Iters:</b> %d | <b>Tools:</b> %v\n\n",
		result.Duration.Round(time.Second), result.Iterations, result.ToolsUsed))
	sb.WriteString(result.Result)

	return sb.String(), true
}

// ──────────────────────────────────────────────
// Batch delegate (parallel execution)
// ──────────────────────────────────────────────

// delegateBatchParams defines input for parallel batch delegation.
type delegateBatchParams struct {
	Tasks    []delegateTaskParams `json:"tasks"`     // Multiple tasks to run in parallel
	MaxBatch int                  `json:"max_batch"` // Override max concurrent (default 5)
}

// ExecuteDelegateBatch spawns multiple subagents in parallel.
// Max concurrent: 5 (configurable). Each subagent runs independently.
func ExecuteDelegateBatch(args map[string]interface{}) (string, bool) {
	var params delegateBatchParams
	params.MaxBatch = helpers.GetIntArg(args, "max_batch", MaxParallelSubagents)

	if params.MaxBatch < 1 {
		params.MaxBatch = 1
	}
	if params.MaxBatch > MaxParallelSubagents {
		params.MaxBatch = MaxParallelSubagents
	}

	// Parse tasks array
	tasksRaw, ok := args["tasks"].([]interface{})
	if !ok || len(tasksRaw) == 0 {
		return "Error: 'tasks' array is required (min 1 task)", false
	}

	for _, raw := range tasksRaw {
		taskMap, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		p := ParseDelegateParams(taskMap)
		if p.Task == "" {
			continue
		}
		if len(p.Tools) == 0 {
			p.Tools = DefaultSubagentTools()
		}
		isOrch := SubagentRole(p.Role) == RoleOrchestrator
		p.Tools = ValidateSubagentTools(p.Tools, isOrch)
		params.Tasks = append(params.Tasks, p)
	}

	if len(params.Tasks) == 0 {
		return "Error: No valid tasks provided", false
	}

	log.Printf("[delegate_batch] Starting %d subagents (max %d parallel)",
		len(params.Tasks), params.MaxBatch)

	// Run in parallel with semaphore
	results := make([]delegateResult, len(params.Tasks))
	sem := make(chan struct{}, params.MaxBatch)
	var wg sync.WaitGroup

	for i, task := range params.Tasks {
		wg.Add(1)
		go func(idx int, t delegateTaskParams) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			log.Printf("[delegate_batch] Subagent %d/%d started: %s",
				idx+1, len(params.Tasks), helpers.TruncateStr(t.Task, 60))
			if t.ACPCommand != "" {
				results[idx] = runSubagentACP(t)
			} else {
				results[idx] = runSubagent(t)
			}
			log.Printf("[delegate_batch] Subagent %d/%d done: status=%s dur=%s",
				idx+1, len(params.Tasks), results[idx].Status,
				results[idx].Duration.Round(time.Second))
		}(i, task)
	}

	wg.Wait()

	// Aggregate results
	var sb strings.Builder
	completed := 0
	failed := 0
	totalDuration := time.Duration(0)

	sb.WriteString(fmt.Sprintf("🔀 <b>Batch Complete: %d subagents</b>\n\n", len(results)))

	for i, r := range results {
		totalDuration += r.Duration
		if r.Status == "completed" {
			completed++
		} else {
			failed++
		}

		statusIcon := "✅"
		if r.Status != "completed" {
			statusIcon = "❌"
		}

		sb.WriteString(fmt.Sprintf("%s <b>[%d] %s</b>", statusIcon, i+1, helpers.TruncateStr(r.Task, 60)))
		if r.ModelUsed != "" {
			sb.WriteString(fmt.Sprintf(" <i>(%s)</i>", r.ModelUsed))
		}
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf("   ⏱ %s | 🔧 %v\n", r.Duration.Round(time.Second), r.ToolsUsed))
		sb.WriteString(fmt.Sprintf("   %s\n\n", helpers.TruncateStr(r.Result, 500)))
	}

	sb.WriteString(fmt.Sprintf("📊 <b>Summary:</b> %d ✅ / %d ❌ | Total: %s",
		completed, failed, totalDuration.Round(time.Second)))

	return sb.String(), failed == 0
}

// AgentMessage represents a single message in the agent conversation history
type AgentMessage struct {
	Role    string
	Content interface{}
}
