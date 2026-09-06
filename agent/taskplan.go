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
// The mutex guards the struct contents (not just the taskPlans map): the
// agent loop renders the plan on the model thread while task_plan tool calls
// mutate it concurrently — a data race here can corrupt completion-gate
// decisions, so every read/write of Goal/Items/UpdatedAt goes through the
// locked accessors below.
type TaskPlan struct {
	mu        sync.RWMutex `json:"-"`
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
	p.mu.Lock()
	p.UpdatedAt = time.Now()
	p.mu.Unlock()
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

// snapshot returns a copy of the items under the read lock — callers may keep
// or pass around the result without holding any lock.
func (p *TaskPlan) snapshot() []PlanItem {
	p.mu.RLock()
	defer p.mu.RUnlock()
	out := make([]PlanItem, len(p.Items))
	copy(out, p.Items)
	return out
}

// Total returns the number of ledger items.
func (p *TaskPlan) Total() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.Items)
}

// Unfinished returns items not yet marked done, in plan order.
func (p *TaskPlan) Unfinished() []PlanItem {
	var out []PlanItem
	for _, it := range p.snapshot() {
		if it.Status != "done" {
			out = append(out, it)
		}
	}
	return out
}

// Progress returns (done, total).
func (p *TaskPlan) Progress() (int, int) {
	items := p.snapshot()
	done := 0
	for _, it := range items {
		if it.Status == "done" {
			done++
		}
	}
	return done, len(items)
}

// Render is the compact checklist injected into model context.
func (p *TaskPlan) Render() string {
	p.mu.RLock()
	goal := p.Goal
	items := make([]PlanItem, len(p.Items))
	copy(items, p.Items)
	p.mu.RUnlock()

	var sb strings.Builder
	sb.WriteString("Goal: ")
	sb.WriteString(goal)
	for _, it := range items {
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

// UpdateItemStatus applies a status change atomically and returns whether the
// item exists, the resulting progress, and the unfinished items after the
// change — all captured under the write lock.
func (p *TaskPlan) UpdateItemStatus(id, status string) (found bool, done, total int, unfinished []PlanItem) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for i := range p.Items {
		if p.Items[i].ID == id {
			p.Items[i].Status = status
			p.UpdatedAt = time.Now()
			for _, it := range p.Items {
				if it.Status == "done" {
					done++
				}
			}
			for _, it := range p.Items {
				if it.Status != "done" {
					unfinished = append(unfinished, it)
				}
			}
			return true, done, len(p.Items), unfinished
		}
	}
	return false, 0, 0, nil
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
		found, done, total, un := plan.UpdateItemStatus(id, status)
		if !found {
			return fmt.Sprintf("Error: item #%s not found. Current plan:\n%s", id, plan.Render()), false
		}
		SetTaskPlan(sessionID, plan)
		out := fmt.Sprintf("Item #%s → %s. Progress: %d/%d done.", id, status, done, total)
		if len(un) > 0 {
			out += "\nStill unfinished:\n" + renderPlanItems(un)
		} else {
			out += "\nALL STEPS DONE — you may now call complete_task with the final verified report."
		}
		return out, true

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
