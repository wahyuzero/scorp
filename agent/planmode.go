package agent

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"scorp-agent/config"
	"scorp-agent/internal/helpers"
	"scorp-agent/models"
	"scorp-agent/tools"
)

// ──────────────────────────────────────────────
// Plan Mode (P1.4) — /plan <goal>
//
// The #1 community-praised agent feature of 2026, composed entirely from
// pieces scorp already proved in verify26:
//   - read-only tool gating: config.IsToolAllowed while a draft runs
//   - the Task Ledger: the draft IS the ledger (persistence + gates apply)
//   - Telegram inline approval: ✅ Approve / ✏️ Revise / ❌ Cancel
//
// Approve resumes execution on the SAME ledger through RunAgentSessionLoop —
// the ACTIVE PLAN injection and completion gates take over from there, so
// there is zero new execution machinery. Revise re-runs the planning loop
// against the current draft. Cancel drops the ledger.
// ──────────────────────────────────────────────

const maxPlanningIterations = 20

const planningDirective = `

## PLANNING MODE (ACTIVE)
You are drafting an execution plan for the user's goal. Rules:
- ONLY read-only tools are allowed (read_file, list_dir, search_code, web_search, read_url, system_info, tool_search). Writes and shell execution are BLOCKED until the user approves the plan.
- Explore enough context to make the plan concrete (files, paths, current state), then record it with task_plan (action=create) as 3-8 concrete, independently verifiable steps. Put specific paths/commands in step descriptions.
- Do NOT mark steps done during planning — the ledger is a draft awaiting approval.
- When the ledger matches the goal, call complete_task with a one-paragraph summary of your approach and any risks.
`

var (
	planningMu      sync.Mutex
	planningActive  bool
	planningSession string
	planningChatID  int64
)

// BeginPlanning marks plan mode active for a session and flips the config
// flag that restricts tool dispatch to the read-only set.
func BeginPlanning(sessionID string, chatID int64) {
	planningMu.Lock()
	planningActive = true
	planningSession = sessionID
	planningChatID = chatID
	planningMu.Unlock()
	config.SetPlanningMode(true)
}

// EndPlanning lifts the read-only restriction (loop exit, approval, cancel).
func EndPlanning() {
	planningMu.Lock()
	planningActive = false
	planningSession = ""
	planningChatID = 0
	planningMu.Unlock()
	config.SetPlanningMode(false)
}

// PlanningState reports the current plan-mode session, if any.
func PlanningState() (sessionID string, chatID int64, active bool) {
	planningMu.Lock()
	defer planningMu.Unlock()
	return planningSession, planningChatID, planningActive
}

// PresentPlanHook renders the drafted plan plus its approval UI. The default
// implementation sends the Telegram inline keyboard; the CLI overrides it
// with a stdin prompt (cli_callbacks.go).
var PresentPlanHook = presentPlanTelegram

// RunPlanningLoop drafts a plan for goal: explores read-only, records the
// ledger via task_plan, and presents it for approval when the model calls
// complete_task (or the iteration budget runs out).
func RunPlanningLoop(sessionID string, chatID int64, goal string, msgID int64) {
	BeginPlanning(sessionID, chatID)
	defer EndPlanning()
	setLoopActive(sessionID, true)
	defer setLoopActive(sessionID, false)

	// A fresh /plan replaces any prior contract for this session.
	ClearTaskPlan(sessionID)
	ClearTaskPlan(fmt.Sprintf("%d", chatID))

	history := []AgentMessage{
		{Role: "system", Content: getAgentSystemPrompt(chatID) + planningDirective},
		{Role: "user", Content: goal},
	}

	runPlanningTurns(sessionID, chatID, msgID, history, goal)
}

// RevisePlan re-enters planning against the existing draft. Returns false if
// there is no pending plan.
func RevisePlan(sessionID string, chatID int64, msgID int64) bool {
	plan := GetTaskPlan(sessionID)
	if plan == nil || plan.Total() == 0 {
		return false
	}
	BeginPlanning(sessionID, chatID)
	go func() {
		defer EndPlanning()
		setLoopActive(sessionID, true)
		defer setLoopActive(sessionID, false)

		history := []AgentMessage{
			{Role: "system", Content: getAgentSystemPrompt(chatID) + planningDirective},
			{Role: "user", Content: "[✏️ REVISION REQUESTED] The user reviewed this plan and wants it improved — more concrete, verifiable steps; fix ordering or gaps. If you need specifics, ask via clarify. Current draft:\n" + plan.Render() + "\n\nUpdate the ledger via task_plan (action=create replaces it), then call complete_task."},
		}
		runPlanningTurns(sessionID, chatID, msgID, history, "")
	}()
	return true
}

// runPlanningTurns is the shared planning iteration loop (fresh draft and
// revision rounds differ only in their starting history).
func runPlanningTurns(sessionID string, chatID int64, msgID int64, history []AgentMessage, goal string) {
	ctx, cancel := context.WithTimeout(context.Background(), maxTurnTimeout())
	defer cancel()

	start := time.Now()
	var thinkingLines []string

	for iter := 0; iter < maxPlanningIterations; iter++ {
		// Cooperative stop
		if ConsumeStopRequest(sessionID) || ConsumeStopRequest(fmt.Sprintf("%d", chatID)) {
			tools.EditMessageByID(chatID, msgID, "⏹ Planning stopped — no plan was kept.", nil)
			return
		}

		// Real-time steering: user redirections land in the draft context
		for {
			steerMsg, ok := PopSteeringMessage(sessionID)
			if !ok {
				break
			}
			log.Printf("[plan] steering while drafting: %s", helpers.TruncateStr(steerMsg, 120))
			history = append(history, AgentMessage{Role: "user", Content: "[⚡ USER STEERING] " + steerMsg})
		}

		chatMsgs := make([]models.ChatMessage, len(history))
		for i, m := range history {
			chatMsgs[i] = models.ChatMessage{Role: m.Role, Content: fmt.Sprintf("%v", m.Content)}
		}

		reply, toolCalls, _, err := models.CallModelWithToolsAndFallback(ctx, "agent", chatMsgs)
		if err != nil {
			tools.EditMessageByID(chatID, msgID, fmt.Sprintf("❌ Planning failed: %v", err), nil)
			return
		}

		// complete_task = draft ready → present for approval
		completed := false
		for _, tc := range toolCalls {
			if tc.Name == "complete_task" {
				completed = true
			}
		}
		if completed {
			plan := GetTaskPlan(sessionID)
			if plan == nil || plan.Total() == 0 {
				history = append(history, AgentMessage{Role: "assistant", Content: reply})
				history = append(history, AgentMessage{Role: "user", Content: "⚠️ No ledger yet — record the plan via task_plan (action=create) with concrete steps, then call complete_task again."})
				continue
			}
			log.Printf("[plan] Plan ready for session %s (%d steps) — awaiting approval", sessionID, plan.Total())
			PresentPlanHook(sessionID, chatID, goal, plan)
			return
		}

		history = append(history, AgentMessage{Role: "assistant", Content: reply})

		if len(toolCalls) == 0 {
			history = append(history, AgentMessage{Role: "user", Content: "⚠️ PLANNING MODE: explore read-only if needed, record the plan via task_plan (action=create), then call complete_task."})
			thinkingLines = append(thinkingLines, "💭 (no tool call)")
		}

		for _, tc := range toolCalls {
			var result string
			var ok bool
			if tc.Name == "task_plan" {
				result, ok = execTaskPlanTool(sessionID, tc.Args)
			} else {
				result, ok = ExecuteTool(tc, chatID)
			}
			if !ok {
				log.Printf("[plan] Tool %s: %s", tc.Name, helpers.TruncateStr(result, 160))
			}
			thinkingLines = append(thinkingLines, "  • "+strings.ReplaceAll(helpers.TruncateStr(result, 60), "\n", " "))

			toolResult := fmt.Sprintf("[Tool Result: %s]\n%s", tc.Name, result)
			history = append(history, AgentMessage{Role: "user", Content: toolResult})
		}

		tools.EditMessageByID(chatID, msgID, "🧭 <b>Drafting plan (read-only)</b> — step "+fmt.Sprintf("%d/%d", iter+1, maxPlanningIterations)+"\n"+buildThinkingMessage(thinkingLines, time.Since(start), false), nil)
	}

	// Budget exhausted — present whatever ledger exists
	plan := GetTaskPlan(sessionID)
	if plan != nil && plan.Total() > 0 {
		log.Printf("[plan] Iteration budget reached — presenting partial draft")
		PresentPlanHook(sessionID, chatID, goal, plan)
		return
	}
	tools.EditMessageByID(chatID, msgID, "⚠️ Planning did not converge — try a more specific <code>/plan</code> goal.", nil)
}

// ApprovePlan executes the pending plan ledger through the main agent loop
// (ACTIVE PLAN injection + completion gates drive it). Returns false when no
// plan is pending.
func ApprovePlan(sessionID string, chatID int64) bool {
	plan := GetTaskPlan(sessionID)
	if plan == nil || plan.Total() == 0 {
		return false
	}
	go RunAgentSessionLoop(sessionID, chatID, "[✅ PLAN APPROVED] Execute every step of the task ledger now — mark each step done via task_plan as you verify it, then complete_task with the final verified report.", 0)
	return true
}

// CancelPlan drops the pending plan draft (both key shapes, since the ledger
// may have been recorded under the session name or the raw chat ID).
func CancelPlan(sessionID string) {
	ClearTaskPlan(sessionID)
}

func presentPlanTelegram(sessionID string, chatID int64, goal string, plan *TaskPlan) {
	// Fresh drafts show the goal in the header — strip Render's duplicate
	// "Goal:" line. Revisions (goal == "") keep it since their header is generic.
	rendered := plan.Render()
	if goal != "" {
		if lines := strings.SplitN(rendered, "\n", 2); len(lines) == 2 && strings.HasPrefix(lines[0], "Goal:") {
			rendered = strings.TrimSpace(lines[1])
		}
	}
	text := "🧭 <b>Plan ready</b>\n" +
		"Goal: " + helpers.EscapeHTML(goal) + "\n\n" +
		"<pre>" + helpers.EscapeHTML(rendered) + "</pre>\n\n" +
		"✅ <b>Approve</b> → execute with full tools · ✏️ <b>Revise</b> → re-draft · ❌ <b>Cancel</b> → discard"
	kb := map[string]interface{}{
		"inline_keyboard": [][]map[string]string{
			{
				{"text": "✅ Approve & Execute", "callback_data": "plan:approve"},
				{"text": "✏️ Revise", "callback_data": "plan:revise"},
			},
			{
				{"text": "❌ Cancel", "callback_data": "plan:cancel"},
			},
		},
	}
	tools.SendMessageGetIDWithKeyboard(text, chatID, kb)
}
