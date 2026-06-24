package agent

import (
	"scorp-agent/tools"
	"scorp-agent/models"
	"scorp-agent/internal/helpers"
	"scorp-agent/registry"
	"scorp-agent/collectors"
	"scorp-agent/config"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// ──────────────────────────────────────────────
// Phase 7 — Autonomous Agent
// ──────────────────────────────────────────────

// AutonomousConfig holds runtime config for the autonomous agent.
type AutonomousConfig = tools.AutonomousConfig

// AutonomousAction is a single action in the LLM's action plan.
type AutonomousAction struct {
	Tool   string                 `json:"tool"`
	Args   map[string]interface{} `json:"args"`
	Reason string                 `json:"reason"`
	Risk   string                 `json:"risk"` // low, medium, high
}

// AutonomousDecision is the parsed LLM response.
type AutonomousDecision struct {
	Analysis string              `json:"analysis"`
	Actions  []AutonomousAction  `json:"actions"`
	Notify   bool                `json:"notify"`
}

// AutonomousLogEntry records a single executed action for audit.
type AutonomousLogEntry = tools.AutonomousLogEntry

var (
	AutoConfig   AutonomousConfig
	AutoMu       sync.Mutex
	AutoLog      []AutonomousLogEntry
	AutoLogFile  = config.ScorpPath("autonomous_log.json")
	AutoKillFile = config.ScorpPath("autonomous_kill")
	AutoCycleNum int
)

// ──────────────────────────────────────────────
// Config Persistence (using ConfigManager)
// ──────────────────────────────────────────────

func LoadAutonomousConfig() {
	AutoMu.Lock()
	defer AutoMu.Unlock()

	if err := config.ConfigMgr().Load("autonomous_config.json", &AutoConfig); err != nil {
		log.Printf("[autonomous] Load error: %v", err)
	}

	// Check kill switch file
	if _, err := os.Stat(AutoKillFile); err == nil {
		AutoConfig.KillSwitch = true
	}

	// Defaults
	if AutoConfig.Interval == 0 {
		AutoConfig.Interval = 10 * time.Minute
	}
	if AutoConfig.MaxActions == 0 {
		AutoConfig.MaxActions = 5
	}
	if AutoConfig.ApprovalLevel == "" {
		AutoConfig.ApprovalLevel = "medium"
	}
}

func SaveAutonomousConfig() {
	AutoMu.Lock()
	defer AutoMu.Unlock()
	saveAutonomousConfigLocked()
}

// saveAutonomousConfigLocked saves config WITHOUT locking — caller must hold AutoMu
func saveAutonomousConfigLocked() {
	if err := config.ConfigMgr().SavePerm("autonomous_config.json", AutoConfig, 0600); err != nil {
		log.Printf("[autonomous] Save error: %v", err)
	}
}

// ──────────────────────────────────────────────
// Kill Switch — file-based, works even if agent crashes
// ──────────────────────────────────────────────

func CheckKillSwitch() bool {
	// Re-check file every time (allows external toggle)
	_, err := os.Stat(AutoKillFile)
	return err == nil
}

func SetKillSwitch(active bool) {
	AutoMu.Lock()
	defer AutoMu.Unlock()

	if active {
		os.MkdirAll(filepath.Dir(AutoKillFile), 0700)
		os.WriteFile(AutoKillFile, []byte("KILL\n"), 0600)
		AutoConfig.KillSwitch = true
		AutoConfig.Enabled = false
	} else {
		os.Remove(AutoKillFile)
		AutoConfig.KillSwitch = false
	}
	saveAutonomousConfigLocked()
}

// ──────────────────────────────────────────────
// Audit Log
// ──────────────────────────────────────────────

func LoadAutoLog() {
	data, err := os.ReadFile(AutoLogFile)
	if err == nil {
		json.Unmarshal(data, &AutoLog)
	}
	// Keep max 500 entries
	if len(AutoLog) > 500 {
		AutoLog = AutoLog[len(AutoLog)-500:]
	}
}

func AppendAutoLog(entry AutonomousLogEntry) {
	AutoLog = append(AutoLog, entry)
	if len(AutoLog) > 500 {
		AutoLog = AutoLog[len(AutoLog)-500:]
	}
	// Ensure directory exists
	os.MkdirAll(filepath.Dir(AutoLogFile), 0700)
	data, _ := json.MarshalIndent(AutoLog, "", "  ")
	os.WriteFile(AutoLogFile, data, 0600)
}

// ──────────────────────────────────────────────
// Context Builder — gather system state snapshot
// ──────────────────────────────────────────────

type AutonomousContext struct {
	Timestamp   string `json:"timestamp"`
	CPU         float64 `json:"cpu"`
	RAM         float64 `json:"ram"`
	Disk        float64 `json:"disk"`
	Load1       float64 `json:"load1"`
	SwapUsedMB  float64 `json:"swap_mb"`
	Uptime      string `json:"uptime"`

	ContainerCount   int      `json:"container_count"`
	UnhealthyCount   int      `json:"unhealthy_count"`
	UnhealthyNames   []string `json:"unhealthy_names"`
	HighCPUContainers []string `json:"high_cpu_containers"`

	Fail2banBanned  int      `json:"fail2ban_banned"`
	RecentSSHFailed int      `json:"ssh_failed"`

	RecentAlerts    []string `json:"recent_alerts"`
	LastActions     []string `json:"last_actions"`
}

func gatherContext() AutonomousContext {
	ctx := AutonomousContext{
		Timestamp: time.Now().Format(time.RFC3339),
	}

	// System metrics
	sys := collectors.CollectSystem()
	ctx.CPU = sys.CPUPercent
	ctx.RAM = sys.RAMPercent
	ctx.Disk = sys.DiskPercent
	ctx.Load1 = sys.LoadAvg[0]
	ctx.SwapUsedMB = sys.SwapUsedGB * 1024
	ctx.Uptime = sys.Uptime

	// Docker
	docker := collectors.CollectDocker()
	ctx.ContainerCount = len(docker.Containers)
	for _, c := range docker.Containers {
		if strings.Contains(strings.ToLower(c.Health), "unhealthy") {
			ctx.UnhealthyCount++
			ctx.UnhealthyNames = append(ctx.UnhealthyNames, c.Name)
		}
		if c.CPUPercent > 80 {
			ctx.HighCPUContainers = append(ctx.HighCPUContainers,
				fmt.Sprintf("%s(%.0f%%)", c.Name, c.CPUPercent))
		}
	}

	// Security
	sec := collectors.CollectSecurity()
	totalBanned := 0
	for _, jail := range sec.Fail2ban.Jails {
		totalBanned += jail.Banned
	}
	ctx.Fail2banBanned = totalBanned
	ctx.RecentSSHFailed = sec.SSHFailedCount

	// Recent alerts (from last 5 log entries that contain "alert")
	for i := len(AutoLog) - 1; i >= 0 && len(ctx.RecentAlerts) < 5; i-- {
		e := AutoLog[i]
		ctx.LastActions = append(ctx.LastActions,
			fmt.Sprintf("[%s] %s: %s", e.Timestamp.Format("15:04"), e.Tool, e.Reason))
	}

	return ctx
}

func (c AutonomousContext) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "Time: %s\n", c.Timestamp)
	fmt.Fprintf(&b, "CPU: %.1f%%, RAM: %.1f%%, Disk: %.1f%%, Load: %.1f\n", c.CPU, c.RAM, c.Disk, c.Load1)
	if c.SwapUsedMB > 0 {
		fmt.Fprintf(&b, "Swap: %.0fMB, ", c.SwapUsedMB)
	}
	fmt.Fprintf(&b, "Uptime: %s\n", c.Uptime)
	fmt.Fprintf(&b, "Containers: %d", c.ContainerCount)
	if c.UnhealthyCount > 0 {
		fmt.Fprintf(&b, ", ⚠️ Unhealthy: %d (%s)", c.UnhealthyCount, strings.Join(c.UnhealthyNames, ", "))
	}
	if len(c.HighCPUContainers) > 0 {
		fmt.Fprintf(&b, ", 🔥 High CPU: %s", strings.Join(c.HighCPUContainers, ", "))
	}
	b.WriteString("\n")
	if c.Fail2banBanned > 0 {
		fmt.Fprintf(&b, "Fail2ban banned: %d", c.Fail2banBanned)
	}
	if c.RecentSSHFailed > 0 {
		fmt.Fprintf(&b, ", SSH failed (24h): %d", c.RecentSSHFailed)
	}
	if b.Len() > 0 && b.String()[b.Len()-1] != '\n' {
		b.WriteString("\n")
	}
	if len(c.LastActions) > 0 {
		fmt.Fprintf(&b, "Recent autonomous actions:\n%s\n", strings.Join(c.LastActions, "\n"))
	}
	return b.String()
}

// ──────────────────────────────────────────────
// Decision Engine — LLM call → structured action plan
// ──────────────────────────────────────────────

const autonomousSystemPrompt = `You are an autonomous monitoring agent. You receive the current system state and must decide what actions to take.

RULES:
1. Only take action when something needs attention (anomalies, thresholds exceeded, unhealthy containers).
2. If everything is normal, return empty actions with a brief analysis.
3. Risk levels: "low" = read-only/monitoring, "medium" = service restart/cleanup, "high" = destructive/security changes.
4. Available tools: exec (shell command), restart (docker container), ban (fail2ban IP), ragvec_add (knowledge).
5. Max %d actions per cycle.
6. Set "notify": true only for important events the user should see.

Respond in JSON only:
{
  "analysis": "brief assessment of current state",
  "actions": [
    {"tool": "exec", "args": {"command": "..."}, "reason": "...", "risk": "low"}
  ],
  "notify": false
}`

func makeDecision(ctxData AutonomousContext) (*AutonomousDecision, error) {
	// Build the LLM prompt
	userPrompt := fmt.Sprintf("Current system state:\n\n%s\n\nWhat actions should be taken?", ctxData.String())

	maxActions := 5
	AutoMu.Lock()
	maxActions = AutoConfig.MaxActions
	AutoMu.Unlock()

	sysPrompt := fmt.Sprintf(autonomousSystemPrompt, maxActions)

	messages := []models.ChatMessage{
		{Role: "system", Content: sysPrompt},
		{Role: "user", Content: userPrompt},
	}

	// Call LLM with 45s timeout
	llmCtx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	resp, _, err := models.CallModelWithFallback(llmCtx, "autonomous", messages)
	if err != nil {
		return nil, fmt.Errorf("LLM call failed: %w", err)
	}

	// Extract JSON from response
	jsonStr := extractJSON(resp)
	var decision AutonomousDecision
	if err := json.Unmarshal([]byte(jsonStr), &decision); err != nil {
		// Fallback: no actions, just use raw response as analysis
		return &AutonomousDecision{
			Analysis: helpers.TruncateStr(resp, 300),
			Actions:  nil,
		}, nil
	}

	// Enforce max actions
	if len(decision.Actions) > maxActions {
		decision.Actions = decision.Actions[:maxActions]
	}

	return &decision, nil
}

// ──────────────────────────────────────────────
// Action Executor — validate → execute → audit
// ──────────────────────────────────────────────

// Tools available to the autonomous agent
func executeAutonomousAction(action AutonomousAction, cycle int) (string, bool) {
	// Validate risk level
	risk := strings.ToLower(action.Risk)
	if risk == "" {
		risk = "medium"
	}

	AutoMu.Lock()
	approvalLevel := AutoConfig.ApprovalLevel
	AutoMu.Unlock()

	// Check if action needs approval
	needsApproval := false
	switch approvalLevel {
	case "high":
		needsApproval = true // everything needs approval
	case "medium":
		needsApproval = risk == "high"
	case "low":
		needsApproval = risk == "high" || risk == "medium"
	}

	if needsApproval {
		// Log as blocked
		entry := AutonomousLogEntry{
			Timestamp: time.Now(),
			Cycle:     cycle,
			Analysis:  "",
			Tool:      action.Tool,
			Args:      action.Args,
			Reason:    action.Reason,
			Risk:      risk,
			Result:    "BLOCKED: requires approval (risk=" + risk + ")",
			Success:   false,
			Approved:  false,
		}
		AppendAutoLog(entry)
		return "BLOCKED: requires user approval", false
	}

	// Execute based on tool type
	var result string
	var ok bool

	switch strings.ToLower(action.Tool) {
	case "exec":
		cmd, _ := action.Args["command"].(string)
		if cmd == "" {
			return "missing command", false
		}
		// Block dangerous commands
		if IsDangerousCommand(cmd) {
			result = "BLOCKED: dangerous command"
			ok = false
		} else {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			out, err := exec.CommandContext(ctx, "bash", "-c", cmd).CombinedOutput()
			result = helpers.TruncateStr(string(out), 500)
			if err != nil {
				result += "\nError: " + err.Error()
				ok = false
			} else {
				ok = true
			}
		}

	case "restart":
		container, _ := action.Args["container"].(string)
		if container == "" {
			return "missing container name", false
		}
		out, err := exec.Command("docker", "restart", container).CombinedOutput()
		result = helpers.TruncateStr(string(out), 300)
		ok = (err == nil)
		if err != nil {
			result += "\nError: " + err.Error()
		}

	case "ban":
		ip, _ := action.Args["ip"].(string)
		jail, _ := action.Args["jail"].(string)
		if jail == "" {
			jail = "sshd"
		}
		if ip == "" {
			return "missing ip", false
		}
		out, err := exec.Command("sudo", "fail2ban-client", "set", jail, "banip", ip).CombinedOutput()
		result = helpers.TruncateStr(string(out), 300)
		ok = (err == nil)

	case "ragvec_add":
		// Delegate to the existing tool
		result, ok = registry.ExecuteToolByName("ragvec_add", action.Args, 0)

	default:
		result = fmt.Sprintf("unknown tool: %s", action.Tool)
		ok = false
	}

	// Audit log
	entry := AutonomousLogEntry{
		Timestamp: time.Now(),
		Cycle:     cycle,
		Tool:      action.Tool,
		Args:      action.Args,
		Reason:    action.Reason,
		Risk:      risk,
		Result:    result,
		Success:   ok,
		Approved:  true,
	}
	AppendAutoLog(entry)

	AutoMu.Lock()
	AutoConfig.TotalActions++
	AutoMu.Unlock()

	return result, ok
}

// ──────────────────────────────────────────────
// Autonomous Loop
// ──────────────────────────────────────────────

func AutonomousLoop(done <-chan struct{}) {
	// Wait for initial setup
	time.Sleep(30 * time.Second)

	log.Println("[autonomous] Loop started")

	for {
		AutoMu.Lock()
		enabled := AutoConfig.Enabled
		interval := AutoConfig.Interval
		AutoMu.Unlock()

		if !enabled || CheckKillSwitch() {
			// Sleep longer when disabled, check periodically
			select {
			case <-done:
				return
			case <-time.After(60 * time.Second):
				continue
			}
		}

		// Run one autonomous cycle
		RunAutonomousCycle()

		// Wait for next cycle
		select {
		case <-done:
			return
		case <-time.After(interval):
		}
	}
}

func RunAutonomousCycle() {
	AutoCycleNum++

	log.Printf("[autonomous] Cycle #%d starting", AutoCycleNum)

	// 1. Gather context
	ctxData := gatherContext()

	// 2. Make decision
	decision, err := makeDecision(ctxData)
	if err != nil {
		log.Printf("[autonomous] Decision error: %v", err)
		// Log the error
		AppendAutoLog(AutonomousLogEntry{
			Timestamp: time.Now(),
			Cycle:     AutoCycleNum,
			Tool:      "decision_engine",
			Result:    "ERROR: " + err.Error(),
			Success:   false,
			Approved:  true,
		})
		return
	}

	log.Printf("[autonomous] Cycle #%d: %s (actions: %d)", AutoCycleNum,
		helpers.TruncateStr(decision.Analysis, 100), len(decision.Actions))

	// 3. Execute actions
	actionResults := []string{}
	for _, action := range decision.Actions {
		result, ok := executeAutonomousAction(action, AutoCycleNum)
		status := "✅"
		if !ok {
			status = "❌"
		}
		actionResults = append(actionResults,
			fmt.Sprintf("%s %s: %s", status, action.Tool, helpers.TruncateStr(result, 100)))
	}

	// 4. Update stats
	AutoMu.Lock()
	AutoConfig.LastCycle = time.Now()
	AutoConfig.TotalCycles++
	AutoMu.Unlock()
	SaveAutonomousConfig()

	// 5. Notify user if needed
	if decision.Notify {
		var msg strings.Builder
		msg.WriteString("🤖 <b>Autonomous Agent</b>\n\n")
		msg.WriteString(fmt.Sprintf("<i>%s</i>\n\n", helpers.EscapeHTML(decision.Analysis)))
		if len(actionResults) > 0 {
			msg.WriteString("<b>Actions:</b>\n")
			for _, r := range actionResults {
				msg.WriteString(helpers.EscapeHTML(r) + "\n")
			}
		} else {
			msg.WriteString("✅ All systems normal. No action needed.\n")
		}
		SendMessageSmart(msg.String(), nil)
	}
}

// ──────────────────────────────────────────────
// Helpers — truncate, escapeHTML, isDangerousCommand
// are defined in scheduler.go / agent_prompt.go
// ──────────────────────────────────────────────

func extractJSON(s string) string {
	// Find first { and last }
	start := strings.Index(s, "{")
	if start == -1 {
		return s
	}
	end := strings.LastIndex(s, "}")
	if end == -1 || end < start {
		return s
	}
	return s[start : end+1]
}

// truncate defined in phase6_test.go
// escapeHTML defined in scheduler.go
