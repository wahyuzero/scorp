package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// ──────────────────────────────────────────────
// Phase 7 — Autonomous Agent Tests
// ──────────────────────────────────────────────

func setupTestPaths(t *testing.T) {
	t.Helper()
	tmpDir := t.TempDir()
	autoCfgFile = filepath.Join(tmpDir, "autonomous_config.json")
	autoLogFile = filepath.Join(tmpDir, "autonomous_log.json")
	autoKillFile = filepath.Join(tmpDir, "autonomous_kill")
}

func TestPhase7_KillSwitch(t *testing.T) {
	setupTestPaths(t)
	// Save state
	origEnabled := autoConfig.Enabled

	// Enable
	autoConfig.Enabled = true
	autoConfig.KillSwitch = false
	os.Remove(autoKillFile)

	// Activate kill switch
	setKillSwitch(true)

	if !autoConfig.KillSwitch {
		t.Error("Kill switch should be active")
	}
	if autoConfig.Enabled {
		t.Error("Kill switch should disable autonomous mode")
	}
	if !checkKillSwitch() {
		t.Error("checkKillSwitch should return true when kill file exists")
	}

	// Deactivate
	setKillSwitch(false)
	if autoConfig.KillSwitch {
		t.Error("Kill switch should be deactivated")
	}
	if checkKillSwitch() {
		t.Error("checkKillSwitch should return false when kill file removed")
	}

	// Restore
	autoConfig.Enabled = origEnabled
	os.Remove(autoKillFile)
}

func TestPhase7_ConfigPersistence(t *testing.T) {
	setupTestPaths(t)
	// Set config
	autoMu.Lock()
	autoConfig.Enabled = true
	autoConfig.Interval = 5 * time.Minute
	autoConfig.VoiceMode = "important"
	autoConfig.ApprovalLevel = "low"
	autoConfig.MaxActions = 3
	autoMu.Unlock()

	saveAutonomousConfig()

	// Verify file exists
	if _, err := os.Stat(autoCfgFile); err != nil {
		t.Fatal("Config file not created:", err)
	}

	// Modify and reload
	autoMu.Lock()
	autoConfig.Enabled = false
	autoConfig.Interval = 10 * time.Minute
	autoMu.Unlock()

	loadAutonomousConfig()

	autoMu.Lock()
	defer autoMu.Unlock()
	if !autoConfig.Enabled {
		t.Error("Enabled should be true from saved config")
	}
	if autoConfig.Interval != 5*time.Minute {
		t.Errorf("Interval = %v, want 5m", autoConfig.Interval)
	}
	if autoConfig.VoiceMode != "important" {
		t.Errorf("VoiceMode = %s, want important", autoConfig.VoiceMode)
	}
	if autoConfig.ApprovalLevel != "low" {
		t.Errorf("ApprovalLevel = %s, want low", autoConfig.ApprovalLevel)
	}
	if autoConfig.MaxActions != 3 {
		t.Errorf("MaxActions = %d, want 3", autoConfig.MaxActions)
	}
}

func TestPhase7_AuditLog(t *testing.T) {
	setupTestPaths(t)
	// Clear log
	autoLog = nil

	// Add entries
	for i := 0; i < 3; i++ {
		appendAutoLog(AutonomousLogEntry{
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

	if len(autoLog) != 3 {
		t.Fatalf("Log should have 3 entries, got %d", len(autoLog))
	}

	// Verify file
	data, err := os.ReadFile(autoLogFile)
	if err != nil {
		t.Fatal("Log file not created:", err)
	}
	if len(data) < 10 {
		t.Error("Log file too small")
	}

	// Reload
	autoLog = nil
	loadAutoLog()
	if len(autoLog) != 3 {
		t.Errorf("After reload, log should have 3 entries, got %d", len(autoLog))
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

	// String representation should not be empty
	s := ctx.String()
	if s == "" {
		t.Error("Context String() should not be empty")
	}
}

func TestPhase7_DecisionParsing(t *testing.T) {
	// Test extractJSON
	tests := []struct {
		input string
		want  string
	}{
		{`{"a":1}`, `{"a":1}`},
		{`some text {"b":2} more text`, `{"b":2}`},
		{`{"c":3`, `{"c":3`},  // no closing brace — returns as-is
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
	// Test low-risk exec action
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
	if result != "hello_phase7\n" && result != "hello_phase7" {
		// truncation might trim, but for short output it should be exact
		t.Logf("Exec result: %q (ok=%v)", result, ok)
	}
}

func TestPhase7_AutonomousAction_BlockedHighRisk(t *testing.T) {
	// Set approval to medium — high risk should be blocked
	autoMu.Lock()
	origApproval := autoConfig.ApprovalLevel
	autoConfig.ApprovalLevel = "medium"
	autoMu.Unlock()
	defer func() {
		autoMu.Lock()
		autoConfig.ApprovalLevel = origApproval
		autoMu.Unlock()
	}()

	action := AutonomousAction{
		Tool:   "exec",
		Args:   map[string]interface{}{"command": "echo dangerous"},
		Reason: "high risk test",
		Risk:   "high",
	}

	result, ok := executeAutonomousAction(action, 1)

	// Should be blocked
	if ok {
		t.Error("High risk action should be blocked with medium approval")
	}
	if result == "" {
		t.Error("Result should explain why blocked")
	}
}

func TestPhase7_AutonomousAction_DangerousCommand(t *testing.T) {
	// Test that dangerous commands are always blocked
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

func TestPhase7_Tool_Status(t *testing.T) {
	result, ok := executeAutonomous(map[string]interface{}{"action": "status"}, 12345)
	if !ok {
		t.Error("Status action should succeed")
	}
	if result == "" {
		t.Error("Status should return non-empty result")
	}
}

func TestPhase7_Tool_EnableDisable(t *testing.T) {
	// Test enable
	_, ok := executeAutonomous(map[string]interface{}{"action": "enable"}, 12345)
	if !ok {
		t.Error("Enable should succeed")
	}
	autoMu.Lock()
	enabledAfterEnable := autoConfig.Enabled
	autoMu.Unlock()
	if !enabledAfterEnable {
		t.Error("Should be enabled after enable")
	}

	// Test disable
	_, ok = executeAutonomous(map[string]interface{}{"action": "disable"}, 12345)
	if !ok {
		t.Error("Disable should succeed")
	}
	autoMu.Lock()
	enabledAfterDisable := autoConfig.Enabled
	autoMu.Unlock()
	if enabledAfterDisable {
		t.Error("Should be disabled after disable")
	}
}

func TestPhase7_Tool_KillRevive(t *testing.T) {
	setupTestPaths(t)
	// Kill
	_, ok := executeAutonomous(map[string]interface{}{"action": "kill"}, 12345)
	if !ok {
		t.Error("Kill should succeed")
	}
	if !checkKillSwitch() {
		t.Error("Kill switch should be active after kill")
	}

	// Revive
	_, ok = executeAutonomous(map[string]interface{}{"action": "revive"}, 12345)
	if !ok {
		t.Error("Revive should succeed")
	}
	if checkKillSwitch() {
		t.Error("Kill switch should be inactive after revive")
	}
	autoMu.Lock()
	enabled := autoConfig.Enabled
	autoMu.Unlock()
	if !enabled {
		t.Error("Should be enabled after revive")
	}
}

func TestPhase7_Tool_Config(t *testing.T) {
	result, ok := executeAutonomous(map[string]interface{}{
		"action":  "config",
		"voice":   "always",
		"approval": "high",
	}, 12345)
	if !ok {
		t.Error("Config action should succeed")
	}
	if result == "" {
		t.Error("Config should return non-empty result")
	}

	autoMu.Lock()
	defer autoMu.Unlock()
	if autoConfig.VoiceMode != "always" {
		t.Errorf("VoiceMode = %s, want always", autoConfig.VoiceMode)
	}
	if autoConfig.ApprovalLevel != "high" {
		t.Errorf("ApprovalLevel = %s, want high", autoConfig.ApprovalLevel)
	}

	// Reset for other tests
	autoConfig.VoiceMode = "off"
	autoConfig.ApprovalLevel = "medium"
}

func TestPhase7_Tool_Log(t *testing.T) {
	result, ok := executeAutonomous(map[string]interface{}{
		"action": "log",
		"count":  5,
	}, 12345)
	if !ok {
		t.Error("Log action should succeed")
	}
	if result == "" {
		t.Error("Log should return non-empty result")
	}
}

func TestPhase7_Tool_Actions(t *testing.T) {
	result, ok := executeAutonomous(map[string]interface{}{
		"action": "actions",
	}, 12345)
	if !ok {
		t.Error("Actions action should succeed")
	}
	if result == "" {
		t.Error("Actions should return non-empty result")
	}
}
