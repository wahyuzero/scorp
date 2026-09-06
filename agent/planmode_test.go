package agent

import (
	"testing"

	"scorp-agent/config"
)

func TestPlanningStateTransitions(t *testing.T) {
	t.Cleanup(EndPlanning)

	if _, _, active := PlanningState(); active {
		t.Fatal("plan mode must start inactive")
	}

	BeginPlanning("sess-plan-1", 42)
	sess, chatID, active := PlanningState()
	if !active || sess != "sess-plan-1" || chatID != 42 {
		t.Fatalf("BeginPlanning state wrong: sess=%q chatID=%d active=%v", sess, chatID, active)
	}
	if !config.PlanningModeActive() {
		t.Fatal("config flag must be set while planning")
	}

	EndPlanning()
	if _, _, active := PlanningState(); active {
		t.Fatal("EndPlanning must deactivate plan mode")
	}
	if config.PlanningModeActive() {
		t.Fatal("config flag must be cleared after EndPlanning")
	}
}

func TestApprovePlanWithoutLedger(t *testing.T) {
	sess := "zz-approve-no-plan"
	t.Cleanup(func() { ClearTaskPlan(sess) })
	ClearTaskPlan(sess)

	if ApprovePlan(sess, 0) {
		t.Fatal("ApprovePlan must report false when no plan is pending")
	}
}

func TestCancelPlanDropsLedger(t *testing.T) {
	sess := "zz-cancel-plan"
	t.Cleanup(func() { ClearTaskPlan(sess) })

	SetTaskPlan(sess, &TaskPlan{
		Goal:  "cancelled goal",
		Items: []PlanItem{{ID: "1", Desc: "step", Status: "pending"}},
	})

	CancelPlan(sess)
	if p := GetTaskPlan(sess); p != nil {
		t.Fatal("CancelPlan must drop the ledger")
	}
}
