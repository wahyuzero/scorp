package scheduler

import (
	"os"
	"strings"
	"testing"
	"time"

	"scorp-agent/config"
)

func TestCronJobScheduler(t *testing.T) {
	config.InitConfigManager()

	// 1. Test NextRunTime with intervals
	now := time.Date(2026, 9, 4, 12, 0, 0, 0, time.UTC)

	next5m, err := NextRunTime("every 5m", now)
	if err != nil {
		t.Fatalf("NextRunTime 'every 5m' failed: %v", err)
	}
	if next5m.Sub(now) != 5*time.Minute {
		t.Errorf("Expected 5m delta, got: %v", next5m.Sub(now))
	}

	next15m, err := NextRunTime("every 15m", now)
	if err != nil {
		t.Fatalf("NextRunTime 'every 15m' failed: %v", err)
	}
	if next15m.Sub(now) != 15*time.Minute {
		t.Errorf("Expected 15m delta, got: %v", next15m.Sub(now))
	}

	// 2. Test NextRunTime with 5-field cron (at minute 30)
	nextCron, err := NextRunTime("30 * * * *", now)
	if err != nil {
		t.Fatalf("NextRunTime cron expression failed: %v", err)
	}
	if nextCron.Minute() != 30 {
		t.Errorf("Expected minute 30, got %d", nextCron.Minute())
	}

	// 3. Test Task CRUD
	task, err := AddTaskEx("daily_checkin", "shell", "every 15m", "echo checkin", nil)
	if err != nil {
		t.Fatalf("AddTaskEx failed: %v", err)
	}
	if task.Name != "daily_checkin" {
		t.Errorf("Expected task name 'daily_checkin', got '%s'", task.Name)
	}

	gotTask := GetTask(task.ID)
	if gotTask == nil || gotTask.Name != "daily_checkin" {
		t.Errorf("Expected GetTask to find task")
	}

	// 4. Test format task list
	desc := FormatTasksList()
	if !strings.Contains(desc, "daily_checkin") {
		t.Errorf("Expected FormatTasksList to contain task name")
	}

	// 5. Test RemoveTask
	ok := RemoveTask(task.ID)
	if !ok {
		t.Errorf("Expected RemoveTask to succeed")
	}
}

// TestDispatchDueTasksOverlapGuard pins the C14 fix: a due task is claimed
// exactly once (in-flight guard + NextRun advanced at dispatch). Before the
// fix, a 150s task on "every 2m" re-fired on every 30s tick — six concurrent
// runs in three minutes.
func TestDispatchDueTasksOverlapGuard(t *testing.T) {
	config.InitConfigManager()

	scheduledTasksMu.Lock()
	savedTasks := scheduledTasks
	savedRunning := make(map[string]bool)
	for k, v := range runningTasks {
		savedRunning[k] = v
	}
	taskIDCounterSaved := taskIDCounter
	scheduledTasks = nil
	taskIDCounter = 0
	runningTasks = map[string]bool{}
	scheduledTasksMu.Unlock()
	defer func() {
		scheduledTasksMu.Lock()
		scheduledTasks = savedTasks
		runningTasks = savedRunning
		taskIDCounter = taskIDCounterSaved
		scheduledTasksMu.Unlock()
	}()

	task, err := AddTaskEx("overlap_guard", "shell", "every 2m", "echo overlap", nil)
	if err != nil {
		t.Fatalf("AddTaskEx failed: %v", err)
	}

	now := time.Now()
	// Force the task due regardless of when AddTaskEx computed NextRun.
	scheduledTasksMu.Lock()
	for i := range scheduledTasks {
		scheduledTasks[i].NextRun = now.Add(-time.Minute)
	}
	scheduledTasksMu.Unlock()

	due := dispatchDueTasks(now)
	if len(due) != 1 {
		t.Fatalf("first dispatch must return exactly the due task, got %d", len(due))
	}

	// Same tick again: the in-flight guard + advanced NextRun must suppress
	// re-dispatch (this is the pile-up that produced 6 concurrent runs).
	if again := dispatchDueTasks(now); len(again) != 0 {
		t.Fatalf("second dispatch on the same tick must be empty, got %d", len(again))
	}
	if again := dispatchDueTasks(now.Add(30 * time.Second)); len(again) != 0 {
		t.Fatalf("dispatch while run is in flight must be empty, got %d", len(again))
	}

	// Run completes → guard released → next dispatch waits for the cadence
	// (NextRun was advanced a full interval at dispatch), not an immediate
	// re-fire.
	scheduledTasksMu.Lock()
	delete(runningTasks, task.ID)
	var nextRun time.Time
	for i := range scheduledTasks {
		if scheduledTasks[i].ID == task.ID {
			nextRun = scheduledTasks[i].NextRun
		}
	}
	scheduledTasksMu.Unlock()
	if !nextRun.After(now.Add(30 * time.Second)) {
		t.Fatalf("NextRun must be advanced ~one full interval, got %v (now=%v)", nextRun, now)
	}
	if again := dispatchDueTasks(nextRun.Add(-time.Second)); len(again) != 0 {
		t.Fatalf("dispatch just before NextRun must be empty, got %d", len(again))
	}
	if due := dispatchDueTasks(nextRun.Add(time.Minute)); len(due) != 1 {
		t.Fatalf("dispatch after NextRun must fire exactly once, got %d", len(due))
	}
}

// TestScheduledShellGateAndGroupKill pins the C14/F-20+F-19 fixes: scheduled
// shell tasks go through deny rules (no silent gate bypass) and a timed-out
// task gets its whole process group killed instead of leaving orphans that
// hold the output pipe.
func TestScheduledShellGateAndGroupKill(t *testing.T) {
	config.InitConfigManager()

	t.Run("deny rule blocks scheduled shell", func(t *testing.T) {
		t.Setenv("SCORP_DENY_RULES", "shell(command:f20deny-marker)")
		config.ReloadDenyRules()
		defer config.ReloadDenyRules()

		task := ScheduledTask{ID: "t_f20", Name: "f20", Type: "shell", Prompt: "echo f20deny-marker"}
		out, status := runShellTaskConfig(task)
		if status != "error" {
			t.Fatalf("deny-ruled command must not run, status=%q out=%q", status, out)
		}
		if !strings.Contains(out, "Denied by deny-rule") {
			t.Fatalf("expected deny-rule denial, got %q", out)
		}
	})

	t.Run("timeout kills whole process group", func(t *testing.T) {
		if os.Getenv("CI_SANITY") == "" && testing.Short() {
			t.Skip("short mode")
		}
		task := ScheduledTask{ID: "t_f19", Name: "f19", Type: "shell",
			Prompt: "sleep 30 && echo never-should-print", Timeout: 1}
		start := time.Now()
		out, status := runShellTaskConfig(task)
		elapsed := time.Since(start)
		if status != "error" {
			t.Fatalf("timed-out task must be an error, status=%q", status)
		}
		if elapsed > 10*time.Second {
			t.Fatalf("timeout must fire at ~1s, took %v (orphan sleep held the pipe)", elapsed)
		}
		if strings.Contains(out, "never-should-print") {
			t.Fatalf("grandchild output must never surface")
		}
	})
}
