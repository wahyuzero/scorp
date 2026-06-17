package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ──────────────────────────────────────────────
// Scheduled Tasks Engine
// ──────────────────────────────────────────────

type ScheduledTask struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`      // "agent", "shell", "script"
	Prompt     string    `json:"prompt"`    // agent prompt or shell command
	Schedule   string    `json:"schedule"`  // "every 5m", "0 9 * * *"
	NextRun    time.Time `json:"next_run"`
	LastRun    time.Time `json:"last_run"`
	LastResult string    `json:"last_result"`
	LastStatus string    `json:"last_status"` // "ok", "error", ""
	Enabled    bool      `json:"enabled"`
	RunCount   int       `json:"run_count"`
	CreatedAt  time.Time `json:"created_at"`

	// S06: Per-job config
	MaxRetries   int    `json:"max_retries,omitempty"`    // retry on failure (default 0)
	Timeout      int    `json:"timeout,omitempty"`        // timeout in seconds (default 120)
	ChatTarget   int64  `json:"chat_target,omitempty"`    // override chat ID for results
	NotifyOnError bool  `json:"notify_on_error,omitempty"` // always notify on error
	NotifyOnSuccess bool `json:"notify_on_success,omitempty"` // also notify on success
	// S06: Context chaining
	ChainContext bool   `json:"chain_context,omitempty"` // feed previous result as context
	PrevResult   string `json:"prev_result,omitempty"`   // last execution result for chaining
}

var (
	scheduledTasks   []ScheduledTask
	scheduledTasksMu sync.Mutex
	schedulerFile    = os.ExpandEnv("$HOME") + "/.scorp-agent/scheduler.json"
	taskIDCounter    int
)

// ──────────────────────────────────────────────
// Persistence
// ──────────────────────────────────────────────

func loadScheduledTasks() {
	scheduledTasksMu.Lock()
	defer scheduledTasksMu.Unlock()

	data, err := os.ReadFile(schedulerFile)
	if err != nil {
		log.Printf("[scheduler] No saved tasks: %v", err)
		return
	}

	if err := json.Unmarshal(data, &scheduledTasks); err != nil {
		log.Printf("[scheduler] Failed to parse tasks: %v", err)
		return
	}

	// Recalculate next run times and find max ID
	maxID := 0
	for i := range scheduledTasks {
		if scheduledTasks[i].Enabled {
			next, err := nextRunTime(scheduledTasks[i].Schedule, time.Now())
			if err == nil {
				scheduledTasks[i].NextRun = next
			}
		}
		// Extract numeric ID for counter
		idStr := strings.TrimPrefix(scheduledTasks[i].ID, "t")
		if n, err := strconv.Atoi(idStr); err == nil && n > maxID {
			maxID = n
		}
	}
	taskIDCounter = maxID

	log.Printf("[scheduler] Loaded %d tasks", len(scheduledTasks))
}

func saveScheduledTasks() {
	os.MkdirAll(os.ExpandEnv("$HOME")+"/.scorp-agent", 0755)
	data, err := json.MarshalIndent(scheduledTasks, "", "  ")
	if err != nil {
		log.Printf("[scheduler] Failed to marshal tasks: %v", err)
		return
	}
	if err := os.WriteFile(schedulerFile, data, 0644); err != nil {
		log.Printf("[scheduler] Failed to save tasks: %v", err)
	}
}

// ──────────────────────────────────────────────
// CRUD
// ──────────────────────────────────────────────

func addScheduledTask(name, taskType, schedule, prompt string) (*ScheduledTask, error) {
	return addScheduledTaskEx(name, taskType, schedule, prompt, nil)
}

// addScheduledTaskEx creates a scheduled task with optional per-job config from args.
func addScheduledTaskEx(name, taskType, schedule, prompt string, args map[string]interface{}) (*ScheduledTask, error) {
	// Validate schedule
	next, err := nextRunTime(schedule, time.Now())
	if err != nil {
		return nil, fmt.Errorf("invalid schedule '%s': %w", schedule, err)
	}

	scheduledTasksMu.Lock()
	defer scheduledTasksMu.Unlock()

	taskIDCounter++
	task := ScheduledTask{
		ID:        fmt.Sprintf("t%d", taskIDCounter),
		Name:      name,
		Type:      taskType,
		Prompt:    prompt,
		Schedule:  schedule,
		NextRun:   next,
		Enabled:   true,
		CreatedAt: time.Now(),
	}

	// S06: Parse per-job config from args
	if args != nil {
		task.MaxRetries = getIntArg(args, "max_retries", 0)
		task.Timeout = getIntArg(args, "timeout", 0)
		task.ChainContext = getBoolArg(args, "chain_context", false)
		task.NotifyOnError = getBoolArg(args, "notify_on_error", false)
		task.NotifyOnSuccess = getBoolArg(args, "notify_on_success", false)
	}

	scheduledTasks = append(scheduledTasks, task)
	saveScheduledTasks()

	log.Printf("[scheduler] Added task %s: %s (%s) type=%s retries=%d timeout=%d chain=%v",
		task.ID, task.Name, task.Schedule, task.Type, task.MaxRetries, task.Timeout, task.ChainContext)
	return &task, nil
}

func removeScheduledTask(id string) bool {
	scheduledTasksMu.Lock()
	defer scheduledTasksMu.Unlock()

	for i, t := range scheduledTasks {
		if t.ID == id {
			scheduledTasks = append(scheduledTasks[:i], scheduledTasks[i+1:]...)
			saveScheduledTasks()
			log.Printf("[scheduler] Removed task %s", id)
			return true
		}
	}
	return false
}

func toggleScheduledTask(id string, enabled bool) bool {
	scheduledTasksMu.Lock()
	defer scheduledTasksMu.Unlock()

	for i, t := range scheduledTasks {
		if t.ID == id {
			scheduledTasks[i].Enabled = enabled
			if enabled {
				next, err := nextRunTime(t.Schedule, time.Now())
				if err == nil {
					scheduledTasks[i].NextRun = next
				}
			}
			saveScheduledTasks()
			return true
		}
	}
	return false
}

func getScheduledTask(id string) *ScheduledTask {
	scheduledTasksMu.Lock()
	defer scheduledTasksMu.Unlock()
	for i, t := range scheduledTasks {
		if t.ID == id {
			return &scheduledTasks[i]
		}
	}
	return nil
}

// ──────────────────────────────────────────────
// Scheduler Loop
// ──────────────────────────────────────────────

func schedulerLoop(done <-chan struct{}) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			log.Println("[scheduler] Stopping")
			return
		case now := <-ticker.C:
			scheduledTasksMu.Lock()
			var dueTasks []ScheduledTask
			for i := range scheduledTasks {
				if scheduledTasks[i].Enabled && !scheduledTasks[i].NextRun.IsZero() && now.After(scheduledTasks[i].NextRun) {
					dueTasks = append(dueTasks, scheduledTasks[i])
				}
			}
			scheduledTasksMu.Unlock()

			for _, task := range dueTasks {
				go runScheduledTask(task)
			}
		}
	}
}

func runScheduledTask(task ScheduledTask) {
	log.Printf("[scheduler] Running task %s: %s", task.ID, task.Name)

	var result string
	var status string

	// S06: Retry logic
	maxRetries := task.MaxRetries
	if maxRetries < 0 {
		maxRetries = 0
	}

	for attempt := 0; attempt <= maxRetries; attempt++ {
		switch task.Type {
		case "shell":
			result, status = runShellTaskConfig(task)
		case "script":
			result, status = runScriptTask(task)
		case "agent":
			// S06: Context chaining — prepend previous result
			prompt := task.Prompt
			if task.ChainContext && task.PrevResult != "" {
				prompt = fmt.Sprintf("Previous run result:\n%s\n\nTask:\n%s",
					truncateStr(task.PrevResult, 2000), task.Prompt)
			}
			result, status = ruScorpAgentTask(prompt)
		default:
			result = fmt.Sprintf("Unknown task type: %s", task.Type)
			status = "error"
		}

		if status == "ok" || attempt == maxRetries {
			break
		}
		log.Printf("[scheduler] Task %s attempt %d/%d failed, retrying...",
			task.ID, attempt+1, maxRetries+1)
		time.Sleep(time.Duration(attempt+1) * 2 * time.Second) // exponential backoff
	}

	// Update task state
	scheduledTasksMu.Lock()
	for i := range scheduledTasks {
		if scheduledTasks[i].ID == task.ID {
			scheduledTasks[i].LastRun = time.Now()
			scheduledTasks[i].LastStatus = status
			scheduledTasks[i].RunCount++

			// Truncate result for storage
			storedResult := result
			if len(storedResult) > 500 {
				storedResult = storedResult[:500] + "..."
			}
			scheduledTasks[i].LastResult = storedResult

			// S06: Context chaining — store previous result
			if scheduledTasks[i].ChainContext {
				chainResult := result
				if len(chainResult) > 2000 {
					chainResult = chainResult[:2000]
				}
				scheduledTasks[i].PrevResult = chainResult
			}

			// Calculate next run
			next, err := nextRunTime(task.Schedule, time.Now())
			if err == nil {
				scheduledTasks[i].NextRun = next
			}
			break
		}
	}
	saveScheduledTasks()
	scheduledTasksMu.Unlock()

	log.Printf("[scheduler] Task %s completed: %s", task.ID, status)
}

func runShellTask(command string) (string, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "bash", "-c", command)
	out, err := cmd.CombinedOutput()
	output := string(out)

	if err != nil {
		// Send error to Telegram
		msg := fmt.Sprintf("⏰ <b>Scheduled Task (Error)</b>\n\n"+
			"<b>Command:</b> <code>%s</code>\n"+
			"<b>Error:</b> %s\n\n<pre>%s</pre>",
			escapeHTML(truncateStr(command, 100)),
			escapeHTML(err.Error()),
			escapeHTML(truncateStr(output, 2000)))
		sendMessageSmart(msg, nil)
		return output, "error"
	}

	// Send result to Telegram
	if len(output) > 0 {
		msg := fmt.Sprintf("⏰ <b>Scheduled Task</b>\n\n"+
			"<b>Command:</b> <code>%s</code>\n\n<pre>%s</pre>",
			escapeHTML(truncateStr(command, 100)),
			escapeHTML(truncateStr(output, 3000)))
		sendMessageSmart(msg, nil)
	}

	return output, "ok"
}

func ruScorpAgentTask(prompt string) (string, string) {
	chatID, err := strconv.ParseInt(cfg.TelegramChatID, 10, 64)
	if err != nil {
		return "Invalid chat ID config", "error"
	}

	// Send a "running" indicator
	msgID := sendMessageGetID(fmt.Sprintf("⏰ <i>Scheduled: %s</i>", escapeHTML(truncateStr(prompt, 80))), chatID)
	if msgID == 0 {
		return "Failed to send message", "error"
	}

	// Run the full agent loop with Thinking Stream
	runAgentLoop(chatID, prompt, msgID)

	return "completed", "ok"
}

// ──────────────────────────────────────────────
// Schedule Parsing: "every Xm/h/d" + cron
// ──────────────────────────────────────────────

var intervalRe = regexp.MustCompile(`^every\s+(\d+)\s*(m|min|h|hour|d|day|s|sec)$`)

func nextRunTime(schedule string, from time.Time) (time.Time, error) {
	schedule = strings.TrimSpace(strings.ToLower(schedule))

	// Try interval format: "every 5m", "every 1h", "every 24h"
	if matches := intervalRe.FindStringSubmatch(schedule); len(matches) == 3 {
		n, _ := strconv.Atoi(matches[1])
		var dur time.Duration
		switch matches[2] {
		case "s", "sec":
			dur = time.Duration(n) * time.Second
		case "m", "min":
			dur = time.Duration(n) * time.Minute
		case "h", "hour":
			dur = time.Duration(n) * time.Hour
		case "d", "day":
			dur = time.Duration(n) * 24 * time.Hour
		}
		if dur < 30*time.Second {
			return time.Time{}, fmt.Errorf("interval too short (min 30s)")
		}
		return from.Add(dur), nil
	}

	// Try cron format: "M H D MON DOW" (5 fields)
	return nextCronTime(schedule, from)
}

// nextCronTime parses standard 5-field cron expression
func nextCronTime(expr string, from time.Time) (time.Time, error) {
	fields := strings.Fields(expr)
	if len(fields) != 5 {
		return time.Time{}, fmt.Errorf("expected 5 cron fields, got %d", len(fields))
	}

	// Parse each field
	minute, err := parseCronField(fields[0], 0, 59)
	if err != nil {
		return time.Time{}, fmt.Errorf("minute: %w", err)
	}
	hour, err := parseCronField(fields[1], 0, 23)
	if err != nil {
		return time.Time{}, fmt.Errorf("hour: %w", err)
	}
	dom, err := parseCronField(fields[2], 1, 31)
	if err != nil {
		return time.Time{}, fmt.Errorf("day of month: %w", err)
	}
	month, err := parseCronField(fields[3], 1, 12)
	if err != nil {
		return time.Time{}, fmt.Errorf("month: %w", err)
	}
	dow, err := parseCronField(fields[4], 0, 6)
	if err != nil {
		return time.Time{}, fmt.Errorf("day of week: %w", err)
	}

	// Brute-force search for next matching time (up to 1 year ahead)
	t := from.Add(1 * time.Minute)
	t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())

	limit := from.Add(366 * 24 * time.Hour)
	for t.Before(limit) {
		if month[int(t.Month())] &&
			dom[t.Day()] &&
			hour[t.Hour()] &&
			minute[t.Minute()] &&
			dow[int(t.Weekday())] {
			return t, nil
		}
		t = t.Add(1 * time.Minute)
	}

	return time.Time{}, fmt.Errorf("no matching time found within 1 year")
}

// parseCronField parses a single cron field into a boolean array of valid values
func parseCronField(field string, min, max int) (map[int]bool, error) {
	result := make(map[int]bool)

	for _, part := range strings.Split(field, ",") {
		part = strings.TrimSpace(part)

		// Handle */N (step)
		if strings.HasPrefix(part, "*/") {
			step, err := strconv.Atoi(strings.TrimPrefix(part, "*/"))
			if err != nil || step <= 0 {
				return nil, fmt.Errorf("invalid step: %s", part)
			}
			for i := min; i <= max; i += step {
				result[i] = true
			}
			continue
		}

		// Handle * (all)
		if part == "*" {
			for i := min; i <= max; i++ {
				result[i] = true
			}
			continue
		}

		// Handle N-M (range)
		if strings.Contains(part, "-") {
			parts := strings.SplitN(part, "-", 2)
			lo, err1 := strconv.Atoi(parts[0])
			hi, err2 := strconv.Atoi(parts[1])
			if err1 != nil || err2 != nil || lo < min || hi > max || lo > hi {
				return nil, fmt.Errorf("invalid range: %s", part)
			}
			for i := lo; i <= hi; i++ {
				result[i] = true
			}
			continue
		}

		// Handle single value
		val, err := strconv.Atoi(part)
		if err != nil || val < min || val > max {
			return nil, fmt.Errorf("invalid value: %s (must be %d-%d)", part, min, max)
		}
		result[val] = true
	}

	return result, nil
}

// ──────────────────────────────────────────────
// Schedule Tool (for agent)
// ──────────────────────────────────────────────

func executeSchedule(args map[string]interface{}) (string, bool) {
	action := getStringArg(args, "action", "")

	switch action {
	case "add":
		name := getStringArg(args, "name", "unnamed")
		taskType := getStringArg(args, "type", "agent")
		schedule := getStringArg(args, "schedule", "")
		prompt := getStringArg(args, "task", "")

		if schedule == "" {
			return "Error: 'schedule' is required (e.g. 'every 1h', '0 9 * * *')", false
		}
		if prompt == "" {
			return "Error: 'task' is required (the prompt or command to execute)", false
		}
		if taskType != "agent" && taskType != "shell" && taskType != "script" {
			taskType = "agent"
		}

		task, err := addScheduledTaskEx(name, taskType, schedule, prompt, args)
		if err != nil {
			return fmt.Sprintf("Error creating schedule: %v", err), false
		}

		return fmt.Sprintf("✅ Scheduled task created:\n"+
			"ID: %s\nName: %s\nType: %s\nSchedule: %s\nNext run: %s\nTask: %s",
			task.ID, task.Name, task.Type, task.Schedule,
			task.NextRun.Format("2006-01-02 15:04:05"),
			truncateStr(task.Prompt, 100)), true

	case "list":
		scheduledTasksMu.Lock()
		defer scheduledTasksMu.Unlock()

		if len(scheduledTasks) == 0 {
			return "No scheduled tasks.", true
		}

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("📋 %d scheduled tasks:\n\n", len(scheduledTasks)))
		for _, t := range scheduledTasks {
			status := "⏸"
			if t.Enabled {
				status = "▶️"
			}
			lastInfo := "never"
			if !t.LastRun.IsZero() {
				lastInfo = fmt.Sprintf("%s (%s)", t.LastRun.Format("01/02 15:04"), t.LastStatus)
			}
			sb.WriteString(fmt.Sprintf("%s [%s] %s\n  Schedule: %s | Next: %s\n  Last: %s | Runs: %d\n  Task: %s\n\n",
				status, t.ID, t.Name,
				t.Schedule, t.NextRun.Format("01/02 15:04"),
				lastInfo, t.RunCount,
				truncateStr(t.Prompt, 60)))
		}
		return sb.String(), true

	case "delete":
		id := getStringArg(args, "id", "")
		if id == "" {
			return "Error: 'id' is required (e.g. 't1')", false
		}
		if removeScheduledTask(id) {
			return fmt.Sprintf("✅ Task %s deleted.", id), true
		}
		return fmt.Sprintf("Error: task %s not found.", id), false

	case "pause":
		id := getStringArg(args, "id", "")
		if id == "" {
			return "Error: 'id' is required", false
		}
		if toggleScheduledTask(id, false) {
			return fmt.Sprintf("⏸ Task %s paused.", id), true
		}
		return fmt.Sprintf("Error: task %s not found.", id), false

	case "resume":
		id := getStringArg(args, "id", "")
		if id == "" {
			return "Error: 'id' is required", false
		}
		if toggleScheduledTask(id, true) {
			return fmt.Sprintf("▶️ Task %s resumed.", id), true
		}
		return fmt.Sprintf("Error: task %s not found.", id), false

	case "run":
		id := getStringArg(args, "id", "")
		if id == "" {
			return "Error: 'id' is required", false
		}
		task := getScheduledTask(id)
		if task == nil {
			return fmt.Sprintf("Error: task %s not found.", id), false
		}
		go runScheduledTask(*task)
		return fmt.Sprintf("🚀 Task %s triggered manually.", id), true

	default:
		return "Error: action must be one of: add, list, delete, pause, resume, run", false
	}
}

// ──────────────────────────────────────────────
// Display helpers
// ──────────────────────────────────────────────

func formatScheduledTasksList() string {
	scheduledTasksMu.Lock()
	defer scheduledTasksMu.Unlock()

	if len(scheduledTasks) == 0 {
		return "📋 <b>Scheduled Tasks</b>\n\nBelum ada jadwal. Buat via agent:\n<i>\"setiap jam cek status docker\"</i>"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📋 <b>Scheduled Tasks</b> (%d)\n\n", len(scheduledTasks)))

	for _, t := range scheduledTasks {
		icon := "⏸"
		if t.Enabled {
			icon = "▶️"
		}
		statusIcon := ""
		if t.LastStatus == "ok" {
			statusIcon = "✅"
		} else if t.LastStatus == "error" {
			statusIcon = "❌"
		}

		sb.WriteString(fmt.Sprintf("%s <code>%s</code> <b>%s</b> %s\n", icon, t.ID, escapeHTML(t.Name), statusIcon))
		sb.WriteString(fmt.Sprintf("   ⏰ %s → next: %s\n", t.Schedule, t.NextRun.Format("01/02 15:04")))
		if t.RunCount > 0 {
			sb.WriteString(fmt.Sprintf("   📊 %d runs, last: %s\n", t.RunCount, t.LastRun.Format("01/02 15:04")))
		}
		sb.WriteString(fmt.Sprintf("   📝 %s: <i>%s</i>\n\n", t.Type, escapeHTML(truncateStr(t.Prompt, 50))))
	}

	sb.WriteString("Commands: <code>/cron run t1</code> | <code>/cron del t1</code>")
	return sb.String()
}

func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}
