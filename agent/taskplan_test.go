package agent

import (
	"strings"
	"sync"
	"testing"
)

func TestTaskPlanCreateUpdateLifecycle(t *testing.T) {
	sid := "test-plan-lifecycle"
	defer ClearTaskPlan(sid)

	if GetTaskPlan(sid) != nil {
		t.Fatal("expected no plan before create")
	}

	out, ok := execTaskPlanTool(sid, map[string]interface{}{
		"action": "create",
		"goal":   "Build 3 reports",
		"items":  []interface{}{"cpu.md", "mem.md", "disk.md"},
	})
	if !ok || !strings.Contains(out, "Task plan created (3 steps)") {
		t.Fatalf("create failed: ok=%v out=%q", ok, out)
	}

	plan := GetTaskPlan(sid)
	if plan == nil || plan.Goal != "Build 3 reports" || len(plan.Items) != 3 {
		t.Fatalf("plan state wrong: %+v", plan)
	}
	if len(plan.Unfinished()) != 3 {
		t.Fatalf("expected 3 unfinished, got %d", len(plan.Unfinished()))
	}

	if _, ok := execTaskPlanTool(sid, map[string]interface{}{"action": "update", "item": "1", "status": "done"}); !ok {
		t.Fatal("update item 1 failed")
	}
	if _, ok := execTaskPlanTool(sid, map[string]interface{}{"action": "update", "item": "2", "status": "in_progress"}); !ok {
		t.Fatal("update item 2 failed")
	}

	plan = GetTaskPlan(sid)
	done, total := plan.Progress()
	if done != 1 || total != 3 {
		t.Fatalf("progress wrong: %d/%d", done, total)
	}
	un := plan.Unfinished()
	if len(un) != 2 || un[0].ID != "2" || un[1].ID != "3" {
		t.Fatalf("unfinished wrong: %+v", un)
	}

	// Complete the rest — the plan must then be fully done.
	execTaskPlanTool(sid, map[string]interface{}{"action": "update", "item": "2", "status": "done"})
	out, ok = execTaskPlanTool(sid, map[string]interface{}{"action": "update", "item": "3", "status": "done"})
	if !ok || !strings.Contains(out, "ALL STEPS DONE") {
		t.Fatalf("final update should report all done: %q", out)
	}
	if len(GetTaskPlan(sid).Unfinished()) != 0 {
		t.Fatal("expected zero unfinished items")
	}
}

func TestTaskPlanValidation(t *testing.T) {
	sid := "test-plan-validation"
	defer ClearTaskPlan(sid)

	if _, ok := execTaskPlanTool(sid, map[string]interface{}{"action": "bogus"}); ok {
		t.Fatal("bogus action must fail")
	}
	if _, ok := execTaskPlanTool(sid, map[string]interface{}{"action": "create"}); ok {
		t.Fatal("create without goal must fail")
	}
	if _, ok := execTaskPlanTool(sid, map[string]interface{}{"action": "create", "goal": "g"}); ok {
		t.Fatal("create without items must fail")
	}
	if _, ok := execTaskPlanTool(sid, map[string]interface{}{"action": "update", "item": "1", "status": "done"}); ok {
		t.Fatal("update without plan must fail")
	}

	if _, ok := execTaskPlanTool(sid, map[string]interface{}{"action": "create", "goal": "g", "items": []interface{}{"a"}}); !ok {
		t.Fatal("valid create failed")
	}
	if _, ok := execTaskPlanTool(sid, map[string]interface{}{"action": "update", "item": "9", "status": "done"}); ok {
		t.Fatal("update of unknown item must fail")
	}
	if _, ok := execTaskPlanTool(sid, map[string]interface{}{"action": "update", "item": "1", "status": "finished"}); ok {
		t.Fatal("invalid status must fail")
	}
}

func TestTaskPlanRender(t *testing.T) {
	sid := "test-plan-render"
	defer ClearTaskPlan(sid)

	execTaskPlanTool(sid, map[string]interface{}{
		"action": "create",
		"goal":   "Ship it",
		"items":  []interface{}{"step one", "step two"},
	})
	execTaskPlanTool(sid, map[string]interface{}{"action": "update", "item": "1", "status": "done"})
	execTaskPlanTool(sid, map[string]interface{}{"action": "update", "item": "2", "status": "in_progress"})

	r := GetTaskPlan(sid).Render()
	for _, want := range []string{"Goal: Ship it", "[x] #1 step one", "[~] #2 step two"} {
		if !strings.Contains(r, want) {
			t.Fatalf("render missing %q in:\n%s", want, r)
		}
	}
}

func TestTaskPlanSessionsIsolated(t *testing.T) {
	defer ClearTaskPlan("sess-a")
	defer ClearTaskPlan("sess-b")

	execTaskPlanTool("sess-a", map[string]interface{}{"action": "create", "goal": "A", "items": []interface{}{"a1"}})
	execTaskPlanTool("sess-b", map[string]interface{}{"action": "create", "goal": "B", "items": []interface{}{"b1", "b2"}})

	if len(GetTaskPlan("sess-a").Items) != 1 || len(GetTaskPlan("sess-b").Items) != 2 {
		t.Fatal("sessions must have isolated ledgers")
	}
}

func TestTaskPlanConcurrentAccess(t *testing.T) {
	sid := "test-plan-concurrent"
	defer ClearTaskPlan(sid)

	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			execTaskPlanTool(sid, map[string]interface{}{"action": "create", "goal": "g", "items": []interface{}{"x", "y"}})
			GetTaskPlan(sid)
			execTaskPlanTool(sid, map[string]interface{}{"action": "update", "item": "1", "status": "done"})
		}()
	}
	wg.Wait()

	if p := GetTaskPlan(sid); p == nil || len(p.Items) != 2 {
		t.Fatal("concurrent access corrupted the ledger")
	}
}
