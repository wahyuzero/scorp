package config

import (
	"strings"
	"testing"
)

// TestIsToolAllowedPlanMode pins the P1.4 tool gate: while a plan draft is
// running, every non-read-only tool is blocked regardless of the autonomy
// level — YOLO included. task_plan/complete_task never reach this gate (both
// are intercepted inside the agent loop).
func TestIsToolAllowedPlanMode(t *testing.T) {
	orig := GetAutonomyLevel()
	defer func() {
		SetAutonomyLevel(string(orig))
		SetPlanningMode(false)
	}()

	SetPlanningMode(true)

	// read-only tools pass in plan mode
	if ok, _ := IsToolAllowed("read_file"); !ok {
		t.Fatal("read_file must be allowed in plan mode")
	}
	if ok, _ := IsToolAllowed("web_search"); !ok {
		t.Fatal("web_search must be allowed in plan mode")
	}

	// mutating tools blocked — even in YOLO
	SetAutonomyLevel("yolo")
	if ok, reason := IsToolAllowed("shell"); ok {
		t.Fatal("shell must be blocked in plan mode even in YOLO")
	} else if !strings.Contains(reason, "Plan Mode") {
		t.Fatalf("plan-mode denial must mention Plan Mode, got %q", reason)
	}
	if ok, _ := IsToolAllowed("write_file"); ok {
		t.Fatal("write_file must be blocked in plan mode")
	}

	// lifting the flag restores normal behavior
	SetPlanningMode(false)
	if ok, _ := IsToolAllowed("shell"); !ok {
		t.Fatal("shell must be allowed again after plan mode ends")
	}
}
