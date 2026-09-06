package agent

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"scorp-agent/internal/helpers"
)

// ──────────────────────────────────────────────
// Task Ledger — enforceable completion contract
//
// Product requirement: the agent must NEVER stop mid-task on its own. A
// free-text "done" from the model is not verifiable, so every multi-step
// request is pinned to an explicit plan. The runtime rejects complete_task
// while any plan item is unfinished and auto-resumes the loop instead of
// ending the turn — the user never has to type "continue".
// ──────────────────────────────────────────────

// PlanItem is one step of the Task Ledger.
type PlanItem struct {
	ID     string `json:"id"`
	Desc   string `json:"desc"`
	Status string `json:"status"` // pending | in_progress | done
}

// TaskPlan is the session-scoped ledger for the current outermost request.
type TaskPlan struct {
	Goal      string     `json:"goal"`
	Items     []PlanItem `json:"items"`
	UpdatedAt time.Time  `json:"updated_at"`
}

var (
	taskPlans   = make(map[string]*TaskPlan)
	taskPlansMu sync.Mutex
)

// SetTaskPlan stores (or replaces) the ledger for a session.
func SetTaskPlan(sessionID string, p *TaskPlan) {
	p.UpdatedAt = time.Now()
	taskPlansMu.Lock()
	defer taskPlansMu.Unlock()
	taskPlans[sessionID] = p
}

// GetTaskPlan returns the session ledger or nil.
func GetTaskPlan(sessionID string) *TaskPlan {
	taskPlansMu.Lock()
	defer taskPlansMu.Unlock()
	return taskPlans[sessionID]
}

// ClearTaskPlan drops the ledger (contract fulfilled or invalidated).
func ClearTaskPlan(sessionID string) {
	taskPlansMu.Lock()
	defer taskPlansMu.Unlock()
	delete(taskPlans, sessionID)
}

// Unfinished returns items not yet marked done, in plan order.
func (p *TaskPlan) Unfinished() []PlanItem {
	var out []PlanItem
	for _, it := range p.Items {
		if it.Status != "done" {
			out = append(out, it)
		}
	}
	return out
}

// Progress returns (done, total).
func (p *TaskPlan) Progress() (int, int) {
	done := 0
	for _, it := range p.Items {
		if it.Status == "done" {
			done++
		}
	}
	return done, len(p.Items)
}

// Render is the compact checklist injected into model context.
func (p *TaskPlan) Render() string {
	var sb strings.Builder
	sb.WriteString("Goal: ")
	sb.WriteString(p.Goal)
	for _, it := range p.Items {
		mark := "[ ]"
		switch it.Status {
		case "done":
			mark = "[x]"
		case "in_progress":
			mark = "[~]"
		}
		sb.WriteString(fmt.Sprintf("\n%s #%s %s", mark, it.ID, it.Desc))
	}
	return sb.String()
}

// maxAutoResumes bounds how many times the runtime may override a premature
// stop within one run; together with maxTurnTimeout it guarantees termination.
const maxAutoResumes = 6

// execTaskPlanTool handles the task_plan tool call inside the agent loop so
// the ledger is scoped to the SESSION (the tool layer only sees chatID).
func execTaskPlanTool(sessionID string, args map[string]interface{}) (string, bool) {
	action := strings.ToLower(strings.TrimSpace(helpers.GetStringArg(args, "action", "")))
	switch action {
	case "create":
		goal := strings.TrimSpace(helpers.GetStringArg(args, "goal", ""))
		if goal == "" {
			return "Error: 'goal' is required for action=create", false
		}
		rawItems, _ := args["items"].([]interface{})
		if len(rawItems) == 0 {
			return "Error: 'items' (array of step descriptions) is required for action=create", false
		}
		plan := &TaskPlan{Goal: goal}
		for i, ri := range rawItems {
			desc, _ := ri.(string)
			desc = strings.TrimSpace(desc)
			if desc == "" {
				continue
			}
			plan.Items = append(plan.Items, PlanItem{ID: strconv.Itoa(i + 1), Desc: desc, Status: "pending"})
		}
		if len(plan.Items) == 0 {
			return "Error: 'items' contained no usable steps", false
		}
		SetTaskPlan(sessionID, plan)
		return fmt.Sprintf("Task plan created (%d steps). The runtime BLOCKS complete_task until every step is done. Current plan:\n%s", len(plan.Items), plan.Render()), true

	case "update":
		plan := GetTaskPlan(sessionID)
		if plan == nil {
			return "Error: no task plan exists for this session. Call task_plan with action=create first.", false
		}
		id := strings.TrimSpace(helpers.GetStringArg(args, "item", ""))
		status := strings.ToLower(strings.TrimSpace(helpers.GetStringArg(args, "status", "")))
		if status != "pending" && status != "in_progress" && status != "done" {
			return "Error: 'status' must be one of: pending, in_progress, done", false
		}
		for i := range plan.Items {
			if plan.Items[i].ID == id {
				plan.Items[i].Status = status
				SetTaskPlan(sessionID, plan)
				done, total := plan.Progress()
				out := fmt.Sprintf("Item #%s → %s. Progress: %d/%d done.", id, status, done, total)
				if un := plan.Unfinished(); len(un) > 0 {
					out += "\nStill unfinished:\n" + renderPlanItems(un)
				} else {
					out += "\nALL STEPS DONE — you may now call complete_task with the final verified report."
				}
				return out, true
			}
		}
		return fmt.Sprintf("Error: item #%s not found. Current plan:\n%s", id, plan.Render()), false

	default:
		return "Error: 'action' must be 'create' or 'update'", false
	}
}

// renderPlanItems formats unfinished items for rejection nudges.
func renderPlanItems(items []PlanItem) string {
	var sb strings.Builder
	for _, it := range items {
		sb.WriteString(fmt.Sprintf("- #%s (%s): %s\n", it.ID, it.Status, it.Desc))
	}
	return strings.TrimRight(sb.String(), "\n")
}
