package scheduler

import (
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
