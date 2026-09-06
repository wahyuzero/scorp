package eval

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"scorp-agent/agent"
	"scorp-agent/config"
	"scorp-agent/registry"
	"scorp-agent/tools"
)

// CoreCases returns the deterministic, model-free checks that guard every
// deploy. Each case re-verifies a production safety/persistence layer the
// same way its regression tests do — from OUTSIDE the packages.
func CoreCases() []Case {
	return []Case{
		{Name: "deny_rule_blocks_shell_even_yolo_confirmed", Category: "safety", Run: caseDenyShellYOLO},
		{Name: "deny_rule_invalid_specs_skipped", Category: "safety", Run: caseDenyInvalidSkipped},
		{Name: "sensitive_path_sandbox_all_modes", Category: "safety", Run: caseSensitivePath},
		{Name: "plan_mode_blocks_writes_even_yolo", Category: "safety", Run: casePlanModeGate},
		{Name: "danger_gate_supervised_needs_confirm", Category: "safety", Run: caseDangerGate},
		{Name: "hooks_block_and_context", Category: "safety", Run: caseHooksBlockAndContext},
		{Name: "ledger_persisted_to_disk", Category: "persistence", Run: caseLedgerPersisted},
		{Name: "ledger_clear_removes_file", Category: "persistence", Run: caseLedgerClear},
		{Name: "memory_md_dedup_and_quota", Category: "memory", Run: caseMemoryMD},
		{Name: "checkpoint_lifecycle_and_undo", Category: "checkpoint", Run: caseCheckpoint},
		{Name: "sandbox_contract_consistent", Category: "sandbox", Run: caseSandbox},
	}
}

// withEnv sets an env var for the duration of a case.
func withEnv(key, val string) (restore func()) {
	orig, had := os.LookupEnv(key)
	if val == "" {
		os.Unsetenv(key)
	} else if err := os.Setenv(key, val); err != nil {
		return func() {}
	}
	return func() {
		if had {
			os.Setenv(key, orig)
		} else {
			os.Unsetenv(key)
		}
	}
}

func caseDenyShellYOLO() error {
	restore := withEnv("SCORP_DENY_RULES", "shell(command:eval-forbidden-marker)")
	defer restore()
	config.ReloadDenyRules()
	defer config.ReloadDenyRules()

	orig := config.GetAutonomyLevel()
	defer config.SetAutonomyLevel(string(orig))
	config.SetAutonomyLevel("yolo")

	out, ok := tools.ExecuteShell(map[string]interface{}{
		"command":   "echo eval-forbidden-marker",
		"confirmed": true,
	}, 0)
	if ok {
		return fmt.Errorf("deny rule did not block a YOLO confirmed shell call: %q", out)
	}
	if !strings.Contains(out, "deny-rule") {
		return fmt.Errorf("expected deny-rule message, got %q", out)
	}
	return nil
}

func caseDenyInvalidSkipped() error {
	restore := withEnv("SCORP_DENY_RULES", "broken;also-broken(param:x);shell(command:eval-good-marker)")
	defer restore()
	config.ReloadDenyRules()
	defer config.ReloadDenyRules()

	blocked, _ := config.CheckDenyRules("shell", map[string]interface{}{"command": "run eval-good-marker now"})
	if !blocked {
		return fmt.Errorf("valid rule after invalid specs did not match")
	}
	blocked, _ = config.CheckDenyRules("shell", map[string]interface{}{"command": "harmless"})
	if blocked {
		return fmt.Errorf("invalid specs must not match everything")
	}
	return nil
}

func caseSensitivePath() error {
	orig := config.GetAutonomyLevel()
	defer config.SetAutonomyLevel(string(orig))
	for _, mode := range []string{"supervised", "yolo"} {
		config.SetAutonomyLevel(mode)
		out, ok := tools.ExecuteShell(map[string]interface{}{
			"command":   "cat /etc/shadow",
			"confirmed": true,
		}, 0)
		if ok || !strings.Contains(out, "Security Sandbox") {
			return fmt.Errorf("[%s] /etc/shadow read not blocked: ok=%v out=%q", mode, ok, out)
		}
	}
	return nil
}

func casePlanModeGate() error {
	orig := config.GetAutonomyLevel()
	defer config.SetAutonomyLevel(string(orig))
	config.SetPlanningMode(true)
	defer config.SetPlanningMode(false)
	config.SetAutonomyLevel("yolo")

	if ok, reason := config.IsToolAllowed("write_file"); ok {
		return fmt.Errorf("write_file allowed in plan mode (yolo)")
	} else if !strings.Contains(reason, "Plan Mode") {
		return fmt.Errorf("unexpected denial reason: %q", reason)
	}
	if ok, _ := config.IsToolAllowed("read_file"); !ok {
		return fmt.Errorf("read_file must stay allowed in plan mode")
	}
	return nil
}

func caseDangerGate() error {
	orig := config.GetAutonomyLevel()
	defer config.SetAutonomyLevel(string(orig))
	config.SetAutonomyLevel("supervised")

	origDanger := tools.IsDangerousCommand
	origStore := tools.StorePendingConfirmation
	defer func() {
		tools.IsDangerousCommand = origDanger
		tools.StorePendingConfirmation = origStore
	}()
	tools.IsDangerousCommand = func(cmd string) bool { return strings.Contains(cmd, "eval-danger") }
	stored := ""
	tools.StorePendingConfirmation = func(chatID, toolName, command string, _ []tools.AgentMessage, _ ...int64) {
		stored = toolName + ":" + command
	}

	out, ok := tools.ExecuteShell(map[string]interface{}{"command": "eval-danger-op"}, 0)
	if ok || !strings.Contains(out, "DANGEROUS COMMAND DETECTED") || stored == "" {
		return fmt.Errorf("supervised danger gate broken: ok=%v stored=%q out=%.80q", ok, stored, out)
	}
	return nil
}

func caseHooksBlockAndContext() error {
	// Deterministic probe tool so the case never depends on real tools.
	registry.RegisterTool(registry.ToolDef{
		Name:     "eval_hook_probe",
		Category: "test",
		Execute: func(map[string]interface{}, int64) (string, bool) {
			return "eval-probe-ran", true
		},
	})
	defer registry.UnregisterTool("eval_hook_probe")

	// Phase 1: PreToolUse exit-2 must block the call, stderr = reason.
	restore := withEnv("SCORP_HOOKS_PRE", "eval_hook_probe:echo 'eval hook policy' >&2 && exit 2")
	config.ReloadHooks()
	out, ok := agent.ExecuteTool(agent.ToolCall{Name: "eval_hook_probe"}, 0)
	restore()
	config.ReloadHooks()
	if ok || !strings.Contains(out, "eval hook policy") {
		return fmt.Errorf("pre-hook did not block: ok=%v out=%.120q", ok, out)
	}

	// Phase 2: exit-0 stdout (pre) and post-hook stdout surface as context.
	r1 := withEnv("SCORP_HOOKS_PRE", "eval_hook_probe:echo pre-eval-ctx")
	defer r1()
	r2 := withEnv("SCORP_HOOKS_POST", "eval_hook_probe:echo post-eval-ctx")
	defer r2()
	config.ReloadHooks()
	defer config.ReloadHooks()
	out2, ok2 := agent.ExecuteTool(agent.ToolCall{Name: "eval_hook_probe"}, 0)
	if !ok2 {
		return fmt.Errorf("exit-0 hooks must not block: %q", out2)
	}
	for _, want := range []string{"eval-probe-ran", "pre-eval-ctx", "post-eval-ctx"} {
		if !strings.Contains(out2, want) {
			return fmt.Errorf("result missing %q: %.160q", want, out2)
		}
	}
	return nil
}

func caseLedgerPersisted() error {
	sess := fmt.Sprintf("eval-ledger-%d", time.Now().UnixNano())
	defer agent.ClearTaskPlan(sess)

	plan := &agent.TaskPlan{
		Goal:  "eval goal",
		Items: []agent.PlanItem{{ID: "1", Desc: "step", Status: "pending"}},
	}
	agent.SetTaskPlan(sess, plan)

	path := filepath.Join(config.ScorpPath("plans"), sess+".plan.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("plan file not persisted: %v", err)
	}
	var onDisk agent.TaskPlan
	if err := json.Unmarshal(data, &onDisk); err != nil {
		return fmt.Errorf("persisted plan unreadable: %v", err)
	}
	if onDisk.Goal != "eval goal" || len(onDisk.Items) != 1 {
		return fmt.Errorf("persisted plan content mismatch: goal=%q items=%d", onDisk.Goal, len(onDisk.Items))
	}
	return nil
}

func caseLedgerClear() error {
	sess := fmt.Sprintf("eval-ledger-clear-%d", time.Now().UnixNano())
	agent.SetTaskPlan(sess, &agent.TaskPlan{
		Goal:  "x",
		Items: []agent.PlanItem{{ID: "1", Desc: "s", Status: "pending"}},
	})
	agent.ClearTaskPlan(sess)
	path := filepath.Join(config.ScorpPath("plans"), sess+".plan.json")
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return fmt.Errorf("cleared plan file still exists")
	}
	return nil
}

func caseMemoryMD() error {
	dir, err := os.MkdirTemp("", "scorp-eval-mem-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	agent.SetMemoryMDFile(filepath.Join(dir, "MEMORY.md"))

	if n := agent.AppendMemoryMD([]string{"eval fact A", "eval fact A", ""}); n != 1 {
		return fmt.Errorf("dedup broken: added %d (want 1)", n)
	}
	if n := agent.AppendMemoryMD([]string{"EVAL FACT A"}); n != 0 {
		return fmt.Errorf("case-insensitive dedup broken: added %d", n)
	}
	content := agent.ReadMemoryMD()
	if !strings.Contains(content, "eval fact A") {
		return fmt.Errorf("entry missing from MEMORY.md")
	}
	return nil
}

func caseCheckpoint() error {
	dir, err := os.MkdirTemp("", "scorp-eval-ckpt-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	if err := os.Chdir(dir); err != nil {
		return err
	}

	run := func(args ...string) error {
		cmd := exec.Command("git", args...)
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("git %s: %v (%s)", strings.Join(args, " "), err, out)
		}
		return nil
	}
	for _, a := range [][]string{{"init", "-q"}, {"config", "user.email", "e@e"}, {"config", "user.name", "e"}} {
		if err := run(a...); err != nil {
			return err
		}
	}
	if err := os.WriteFile(filepath.Join(dir, "f.txt"), []byte("v1\n"), 0644); err != nil {
		return err
	}

	_, created, err := tools.CreateCheckpoint("eval-ckpt")
	if err != nil || !created {
		return fmt.Errorf("checkpoint 1: created=%v err=%v", created, err)
	}
	if err := os.WriteFile(filepath.Join(dir, "f.txt"), []byte("v2\n"), 0644); err != nil {
		return err
	}
	cks, err := tools.ListCheckpoints("eval-ckpt")
	if err != nil || len(cks) == 0 {
		return fmt.Errorf("list failed: %v", err)
	}
	n, err := tools.RestoreCheckpoint("eval-ckpt", cks[0].Ref)
	if err != nil {
		return fmt.Errorf("restore failed: %v", err)
	}
	got, _ := os.ReadFile(filepath.Join(dir, "f.txt"))
	if string(got) != "v1\n" {
		return fmt.Errorf("undo restored wrong content: %q", got)
	}
	if n < 1 {
		return fmt.Errorf("restore count %d", n)
	}
	return nil
}

func caseSandbox() error {
	active := tools.SandboxActive()
	argv, wrapped := tools.SandboxWrap("echo x")
	if active && !wrapped {
		return fmt.Errorf("sandbox active but Wrap declined")
	}
	if !active && wrapped {
		return fmt.Errorf("sandbox inactive but Wrap wrapped: %v", argv)
	}
	if active {
		joined := strings.Join(argv, " ")
		for _, want := range []string{"--ro-bind / /", "--unshare-all", "-- bash -c echo x"} {
			if !strings.Contains(joined, want) {
				return fmt.Errorf("sandbox argv missing %q", want)
			}
		}
	}
	return nil
}
