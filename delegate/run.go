package delegate

import (
	"scorp-agent/models"
	"scorp-agent/registry"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// ──────────────────────────────────────────────
// Subagent Execution Engine
// ──────────────────────────────────────────────

// runSubagent executes a single subagent with its own history, tool restrictions,
// and model routing based on role/model params.
func runSubagent(params delegateTaskParams) delegateResult {
	start := time.Now()
	subagentID := fmt.Sprintf("sub_%d", start.UnixNano())

	// S08: Create isolation context
	iso := defaultIsolation()
	iso.WorkDir = createSubagentSandbox(subagentID)
	registerSubagentIsolation(subagentID, iso)
	defer func() {
		cleanupSubagentSandbox(subagentID)
		unregisterSubagentIsolation(subagentID)
	}()

	// S08: Filter out state-modifying tools even if requested
	filteredTools := make([]string, 0, len(params.Tools))
	for _, t := range params.Tools {
		if isSubagentToolBlocked(t) {
			log.Printf("[isolation] Subagent %s: blocked tool '%s'", subagentID, t)
			continue
		}
		filteredTools = append(filteredTools, t)
	}
	params.Tools = filteredTools

	log.Printf("[isolation] Subagent %s sandbox=%s tools=%d", subagentID, iso.WorkDir, len(params.Tools))

	// Determine which model to use
	taskType := roleToTaskType(SubagentRole(params.Role))
	modelNameUsed := ""

	// Build system prompt for subagent
	subPrompt := buildSubagentPrompt(params)

	// Subagent history (independent from main agent)
	history := []AgentMessage{
		{Role: "system", Content: subPrompt},
		{Role: "user", Content: params.Task},
	}

	// Track tools used
	toolsUsed := make(map[string]bool)
	var toolNames []string

	for iteration := 0; iteration < params.MaxIters; iteration++ {
		// Build message history
		chatMsgs := make([]models.ChatMessage, len(history))
		for i, m := range history {
			switch c := m.Content.(type) {
			case string:
				chatMsgs[i] = models.ChatMessage{Role: m.Role, Content: c}
			default:
				jsonBytes, _ := json.Marshal(c)
				chatMsgs[i] = models.ChatMessage{Role: m.Role, Content: string(jsonBytes)}
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)

		var text string
		var _, _, _ = text, chatMsgs, ctx
		var toolCalls []models.ToolCall
		var modelLabel string
		var err error

		// If explicit model override, use it; otherwise route by role
		if params.Model != "" {
			m := models.GetModelByName(params.Model)
			if m != nil {
				text, err = models.CallModel(ctx, m, chatMsgs)
				if err != nil {
					// Fallback to routed model
					text, toolCalls, modelLabel, err = models.CallModelWithToolsAndFallback(ctx, taskType, chatMsgs)
				} else {
					modelLabel = m.Model
					toolCalls, _ = models.ParseAllToolCalls(text, nil)
				}
			} else {
				text, toolCalls, modelLabel, err = models.CallModelWithToolsAndFallback(ctx, taskType, chatMsgs)
			}
		} else {
			text, toolCalls, modelLabel, err = models.CallModelWithToolsAndFallback(ctx, taskType, chatMsgs)
		}

		modelNameUsed = modelLabel
		cancel()

		if err != nil {
			return delegateResult{
				SubagentID: subagentID,
				Task:       params.Task,
				Role:       params.Role,
				ModelUsed:  modelNameUsed,
				Status:     "failed",
				Error:      fmt.Sprintf("model error (iter %d): %v", iteration+1, err),
				Duration:   time.Since(start),
			}
		}

		// Parse tool calls
		if len(toolCalls) == 0 {
			// Check if model returned tool calls in text format
			toolCalls, text = models.ParseAllToolCalls(text, nil)
		}

		if len(toolCalls) == 0 {
			// Final answer — subagent is done
			return delegateResult{
				SubagentID: subagentID,
				Task:       params.Task,
				Role:       params.Role,
				ModelUsed:  modelNameUsed,
				Status:     "completed",
				Result:     text,
				ToolsUsed:  toolNames,
				Iterations: iteration + 1,
				Duration:   time.Since(start),
			}
		}

		// Execute tools (only allowed ones)
		history = append(history, AgentMessage{Role: "assistant", Content: text})

		for _, tc := range toolCalls {
			// Check if tool is allowed
			allowed := false
			for _, t := range params.Tools {
				if t == tc.Name {
					allowed = true
					break
				}
			}

			// CRITICAL: Block re-delegation except for orchestrator role
			if tc.Name == "delegate" || tc.Name == "delegate_batch" {
				if SubagentRole(params.Role) == RoleOrchestrator {
					// Orchestrator can spawn workers — allow it
					if !toolsUsed[tc.Name] {
						toolsUsed[tc.Name] = true
						toolNames = append(toolNames, tc.Name)
					}
					result, _ := registry.ExecuteToolByName(tc.Name, tc.Args, 0)
					if iso != nil {
						result = iso.truncateOutput(result)
					}
					history = append(history, AgentMessage{Role: "user", Content: fmt.Sprintf("[Tool Result: %s]\n%s", tc.Name, result)})
					continue
				}
				toolResult := fmt.Sprintf("[Tool Result: %s]\nError: Re-delegation is BLOCKED. You cannot spawn subagents. Complete the task yourself using your available tools.", tc.Name)
				history = append(history, AgentMessage{Role: "user", Content: toolResult})
				continue
			}

			if !allowed {
				toolResult := fmt.Sprintf("[Tool Result: %s]\nError: Tool '%s' not allowed for this subagent", tc.Name, tc.Name)
				history = append(history, AgentMessage{Role: "user", Content: toolResult})
				continue
			}

			if !toolsUsed[tc.Name] {
				toolsUsed[tc.Name] = true
				toolNames = append(toolNames, tc.Name)
			}

			result, _ := registry.ExecuteToolByName(tc.Name, tc.Args, 0) // chatID=0 for subagent (no Telegram)

			// S08: Enforce output size limit
			if iso != nil {
				result = iso.truncateOutput(result)
			}

			toolResult := fmt.Sprintf("[Tool Result: %s]\n%s", tc.Name, result)
			history = append(history, AgentMessage{Role: "user", Content: toolResult})
		}
	}

	// Max iterations reached
	log.Printf("[delegate] Subagent %s hit max iterations (%d)", subagentID, params.MaxIters)
	return delegateResult{
		SubagentID: subagentID,
		Task:       params.Task,
		Role:       params.Role,
		ModelUsed:  modelNameUsed,
		Status:     "timeout",
		Result:     "Max iterations reached. Partial result from last response.",
		ToolsUsed:  toolNames,
		Iterations: params.MaxIters,
		Duration:   time.Since(start),
	}
}
