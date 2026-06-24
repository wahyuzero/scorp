package agent

import (
	"scorp-agent/config"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// ──────────────────────────────────────────────
// Phase 7 — Autonomous Agent Tests (main package internals)
// ──────────────────────────────────────────────

func setupTestPaths(t *testing.T) {
	t.Helper()
	tmpDir := t.TempDir()
	config.ConfigMgr().OverrideBaseDir(tmpDir)
	AutoKillFile = filepath.Join(tmpDir, "autonomous_kill")
	AutoLogFile = filepath.Join(tmpDir, "autonomous_log.json")
}

func TestPhase7_KillSwitch(t *testing.T) {
	setupTestPaths(t)

	origEnabled := AutoConfig.Enabled
	AutoConfig.Enabled = true
	AutoConfig.KillSwitch = false
	os.Remove(AutoKillFile)

	SetKillSwitch(true)

	if !AutoConfig.KillSwitch {
		t.Error("Kill switch should be active")
	}
	if AutoConfig.Enabled {
		t.Error("Kill switch should disable autonomous mode")
	}
	if !CheckKillSwitch() {
		t.Error("CheckKillSwitch should return true when kill file exists")
	}

	SetKillSwitch(false)
	if AutoConfig.KillSwitch {
		t.Error("Kill switch should be deactivated")
	}
	if CheckKillSwitch() {
		t.Error("CheckKillSwitch should return false when kill file removed")
	}

	AutoConfig.Enabled = origEnabled
	os.Remove(AutoKillFile)
}

func TestPhase7_ConfigPersistence(t *testing.T) {
	setupTestPaths(t)

	AutoMu.Lock()
	AutoConfig.Enabled = true
	AutoConfig.Interval = 5 * time.Minute
	AutoConfig.ApprovalLevel = "low"
	AutoConfig.MaxActions = 3
	AutoMu.Unlock()

	SaveAutonomousConfig()

	if _, err := os.Stat(config.ConfigMgr().Path("autonomous_config.json")); err != nil {
		t.Fatal("Config file not created:", err)
	}

	AutoMu.Lock()
	AutoConfig.Enabled = false
	AutoConfig.Interval = 10 * time.Minute
	AutoMu.Unlock()

	LoadAutonomousConfig()

	AutoMu.Lock()
	defer AutoMu.Unlock()
	if !AutoConfig.Enabled {
		t.Error("Enabled should be true from saved config")
	}
	if AutoConfig.Interval != 5*time.Minute {
		t.Errorf("Interval = %v, want 5m", AutoConfig.Interval)
	}
	if AutoConfig.ApprovalLevel != "low" {
		t.Errorf("ApprovalLevel = %s, want low", AutoConfig.ApprovalLevel)
	}
	if AutoConfig.MaxActions != 3 {
		t.Errorf("MaxActions = %d, want 3", AutoConfig.MaxActions)
	}
}

func TestPhase7_AuditLog(t *testing.T) {
	setupTestPaths(t)
	AutoLog = nil

	for i := 0; i < 3; i++ {
		AppendAutoLog(AutonomousLogEntry{
			Timestamp: time.Now(),
			Cycle:     i + 1,
			Tool:      "exec",
			Reason:    "test action",
			Risk:      "low",
			Result:    "success",
			Success:   true,
			Approved:  true,
		})
	}

	if len(AutoLog) != 3 {
		t.Fatalf("Log should have 3 entries, got %d", len(AutoLog))
	}

	data, err := os.ReadFile(AutoLogFile)
	if err != nil {
		t.Fatal("Log file not created:", err)
	}
	if len(data) < 10 {
		t.Error("Log file too small")
	}

	AutoLog = nil
	LoadAutoLog()
	if len(AutoLog) != 3 {
		t.Errorf("After reload, log should have 3 entries, got %d", len(AutoLog))
	}
}

func TestPhase7_GatherContext(t *testing.T) {
	ctx := gatherContext()

	if ctx.Timestamp == "" {
		t.Error("Timestamp should not be empty")
	}
	if ctx.CPU < 0 || ctx.CPU > 100 {
		t.Errorf("CPU = %.1f, should be 0-100", ctx.CPU)
	}
	if ctx.RAM < 0 || ctx.RAM > 100 {
		t.Errorf("RAM = %.1f, should be 0-100", ctx.RAM)
	}
	if ctx.ContainerCount < 0 {
		t.Errorf("ContainerCount = %d, should be >= 0", ctx.ContainerCount)
	}

	s := ctx.String()
	if s == "" {
		t.Error("Context String() should not be empty")
	}
}

func TestPhase7_DecisionParsing(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{`{"a":1}`, `{"a":1}`},
		{`some text {"b":2} more text`, `{"b":2}`},
		{`{"c":3`, `{"c":3`},
		{`no json here`, `no json here`},
		{`{"nested": {"d":4}}`, `{"nested": {"d":4}}`},
	}

	for _, tc := range tests {
		got := extractJSON(tc.input)
		if got != tc.want {
			t.Errorf("extractJSON(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestPhase7_AutonomousAction_Exec(t *testing.T) {
	action := AutonomousAction{
		Tool:   "exec",
		Args:   map[string]interface{}{"command": "echo hello_phase7"},
		Reason: "test echo",
		Risk:   "low",
	}

	result, ok := executeAutonomousAction(action, 1)
	if !ok {
		t.Errorf("Exec should succeed, result: %s", result)
	}
}

func TestPhase7_AutonomousAction_BlockedHighRisk(t *testing.T) {
	AutoMu.Lock()
	origApproval := AutoConfig.ApprovalLevel
	AutoConfig.ApprovalLevel = "medium"
	AutoMu.Unlock()
	defer func() {
		AutoMu.Lock()
		AutoConfig.ApprovalLevel = origApproval
		AutoMu.Unlock()
	}()

	action := AutonomousAction{
		Tool:   "exec",
		Args:   map[string]interface{}{"command": "echo dangerous"},
		Reason: "high risk test",
		Risk:   "high",
	}

	result, ok := executeAutonomousAction(action, 1)

	if ok {
		t.Error("High risk action should be blocked with medium approval")
	}
	if result == "" {
		t.Error("Result should explain why blocked")
	}
}

func TestPhase7_AutonomousAction_DangerousCommand(t *testing.T) {
	action := AutonomousAction{
		Tool:   "exec",
		Args:   map[string]interface{}{"command": "rm -rf /"},
		Reason: "dangerous test",
		Risk:   "low",
	}

	result, ok := executeAutonomousAction(action, 1)

	if ok {
		t.Error("Dangerous command should be blocked")
	}
	if result == "" {
		t.Error("Result should explain blocking")
	}
}

func TestPhase7_AutonomousAction_UnknownTool(t *testing.T) {
	action := AutonomousAction{
		Tool:   "nonexistent_tool",
		Args:   map[string]interface{}{},
		Reason: "test unknown",
		Risk:   "low",
	}

	result, ok := executeAutonomousAction(action, 1)

	if ok {
		t.Error("Unknown tool should fail")
	}
	if result == "" {
		t.Error("Result should explain failure")
	}
}
