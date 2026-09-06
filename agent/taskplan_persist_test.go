package agent

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// uniqueSess gives each test its own ledger file so parallel test runs and
// real sessions on this machine never collide; ClearTaskPlan removes the
// persisted file during cleanup.
func uniqueSess(t *testing.T) string {
	t.Helper()
	sess := fmt.Sprintf("zz-p15-%d", time.Now().UnixNano())
	t.Cleanup(func() { ClearTaskPlan(sess) })
	return sess
}

// simulateRestart drops every in-memory ledger — persistence must carry state
// across this boundary exactly as it does across a daemon restart.
func simulateRestart() {
	taskPlansMu.Lock()
	taskPlans = make(map[string]*TaskPlan)
	taskPlansMu.Unlock()
}

func mkTestPlan(goal string) *TaskPlan {
	return &TaskPlan{
		Goal: goal,
		Items: []PlanItem{
			{ID: "1", Desc: "step one", Status: "pending"},
			{ID: "2", Desc: "step two", Status: "pending"},
		},
	}
}

func TestPlanPersistenceRoundtripAcrossRestart(t *testing.T) {
	sess := uniqueSess(t)
	SetTaskPlan(sess, mkTestPlan("deploy the thing"))

	if _, err := os.Stat(planFilePath(sess)); err != nil {
		t.Fatalf("SetTaskPlan must persist the ledger to disk: %v", err)
	}

	simulateRestart()
	p := GetTaskPlan(sess)
	if p == nil {
		t.Fatal("ledger must be reloaded from disk after a simulated restart")
	}
	if p.Goal != "deploy the thing" || p.Total() != 2 {
		t.Fatalf("reloaded ledger corrupted: goal=%q items=%d", p.Goal, p.Total())
	}
	if len(p.Unfinished()) != 2 {
		t.Fatalf("reloaded ledger must keep unfinished items, got %d", len(p.Unfinished()))
	}
}

func TestClearTaskPlanRemovesPersistedFile(t *testing.T) {
	sess := uniqueSess(t)
	SetTaskPlan(sess, mkTestPlan("done contract"))
	ClearTaskPlan(sess)

	if _, err := os.Stat(planFilePath(sess)); !os.IsNotExist(err) {
		t.Fatal("ClearTaskPlan must delete the persisted plan file")
	}
	simulateRestart()
	if GetTaskPlan(sess) != nil {
		t.Fatal("cleared plan must not resurrect after restart")
	}
}

func TestStalePlanExpiredOnLoad(t *testing.T) {
	sess := uniqueSess(t)
	SetTaskPlan(sess, mkTestPlan("ancient task"))

	// backdate the on-disk UpdatedAt beyond the expiry window
	path := planFilePath(sess)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read plan file: %v", err)
	}
	var p TaskPlan
	if err := json.Unmarshal(data, &p); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	p.UpdatedAt = time.Now().Add(-8 * 24 * time.Hour)
	aged, _ := json.Marshal(&p)
	if err := os.WriteFile(path, aged, 0644); err != nil {
		t.Fatalf("backdate write: %v", err)
	}

	simulateRestart()
	if GetTaskPlan(sess) != nil {
		t.Fatal("stale plan beyond expiry must be discarded on load")
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("expired plan file must be removed on load")
	}
}

func TestPlanFilePathSanitizesSessionID(t *testing.T) {
	plansDir := filepath.Dir(planFilePath("anything"))
	path := planFilePath("../..//evil-session")
	if filepath.Dir(path) != plansDir {
		t.Fatalf("session id must be sanitized: path escaped the plans dir: %q", path)
	}
	if strings.Contains(path, "..") {
		t.Fatalf("sanitized path must not contain traversal segments: %q", path)
	}
}

func TestUpdateItemStatusPersistsThroughRestart(t *testing.T) {
	sess := uniqueSess(t)
	p := mkTestPlan("persist statuses")
	SetTaskPlan(sess, p)

	if found, done, total, un := p.UpdateItemStatus("1", "done"); !found || done != 1 || total != 2 || len(un) != 1 {
		t.Fatalf("UpdateItemStatus unexpected: found=%v done=%d total=%d un=%d", found, done, total, len(un))
	}
	SetTaskPlan(sess, p) // mirrors execTaskPlanTool's post-update persist

	simulateRestart()
	reloaded := GetTaskPlan(sess)
	if reloaded == nil {
		t.Fatal("ledger must survive restart after an item update")
	}
	done, total := reloaded.Progress()
	if done != 1 || total != 2 {
		t.Fatalf("status change lost across restart: done=%d total=%d", done, total)
	}
}
