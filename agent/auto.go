package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"scorp-agent/config"
	"scorp-agent/internal/helpers"
	"scorp-agent/models"
)

// ──────────────────────────────────────────────
// Auto-Mode Classifier (P3.13)
//
// `auto` autonomy (between supervised and yolo): every tool call is graded
// before execution —
//   safe       → runs without confirmation
//   risky      → routed through the existing confirmation gate
//   destructive→ hard-denied unless explicitly allowlisted
//
// Grading is layered, cheapest first:
//  1. deterministic: read-only tools / read-only shell commands → safe
//  2. deterministic: IsDangerousCommand (destructive shell) → deny unless the
//     command matches an SCORP_AUTO_ALLOW regex
//  3. cheap model (models.RouteModel("chat")) grades everything else, asked
//     for strict JSON {"decision":"safe|ask|deny","reason":"..."}
//  4. model failure/uncertain → treated as ask (fail-closed); after
//     autoMaxUncertain consecutive uncertain results the classifier is
//     considered broken and the session degrades to supervised-style
//     behavior (deterministic safe stays safe, everything else asks) until
//     the next task boundary.
//
// Enforcement points:
//   - agent loops (main + resume): full ask support → inline confirmation
//     keyboard + pause + resume with history (handleAutoModeGate)
//   - agent.ExecuteTool: deterministic deny for every other path (subagents
//     have no user channel — an ask degrades to deny, fail-closed)
//
// Every decision is logged and (for executed calls) recorded into the tool
// receipt meta as auto_decision.
// ──────────────────────────────────────────────

const (
	autoMaxUncertain    = 3 // consecutive uncertain decisions before fallback
	autoClassifyTimeout = 15 * time.Second
)

// AutoPermission decisions.
const (
	AutoAllow = "allow"
	AutoAsk   = "ask"
	AutoDeny  = "deny"
)

// AutoClassifyFunc is the model-backed grader. Package var so tests and the
// eval arena can stub it deterministically.
var AutoClassifyFunc = autoClassifyWithModel

type autoStats struct {
	mu               sync.Mutex
	consecUncertain  int
	totalSafe        int
	totalAsk         int
	totalDeny        int
	totalUncertain   int
	fallbackActive   bool
}

var autoStat autoStats

// ResetAutoStats clears classifier stats and any fallback state — called at
// task boundaries so fallback degrades only within the task that hit it.
func ResetAutoStats() {
	autoStat.mu.Lock()
	defer autoStat.mu.Unlock()
	autoStat.consecUncertain = 0
	autoStat.fallbackActive = false
}

// ResetAutoAllowlist re-reads SCORP_AUTO_ALLOW on next use (tests, config
// reload paths).
func ResetAutoAllowlist() {
	autoAllowlistOnce = sync.Once{}
	autoAllowlist = nil
}

// AutoStatsSnapshot returns counters for status surfaces.
func AutoStatsSnapshot() (safe, ask, deny, uncertain int, fallback bool) {
	autoStat.mu.Lock()
	defer autoStat.mu.Unlock()
	return autoStat.totalSafe, autoStat.totalAsk, autoStat.totalDeny, autoStat.totalUncertain, autoStat.fallbackActive
}

// autoAllowlist compiles SCORP_AUTO_ALLOW (semicolon-separated regexes). A
// destructive shell command matching any entry is permitted instead of
// hard-denied — the explicit "I know what this does" escape hatch.
var (
	autoAllowlist     []*regexp.Regexp
	autoAllowlistOnce sync.Once
)

func autoAllowlisted(cmd string) bool {
	autoAllowlistOnce.Do(func() {
		raw := os.Getenv("SCORP_AUTO_ALLOW")
		for _, part := range strings.Split(raw, ";") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			re, err := regexp.Compile(part)
			if err != nil {
				log.Printf("[auto] ignoring invalid SCORP_AUTO_ALLOW regex %q: %v", part, err)
				continue
			}
			autoAllowlist = append(autoAllowlist, re)
		}
	})
	for _, re := range autoAllowlist {
		if re.MatchString(cmd) {
			return true
		}
	}
	return false
}

// IsReadOnlyShellCommand heuristically reports whether a shell command only
// reads (no mutation, no network writes, no code execution side effects).
// Conservative: anything not provably read-only returns false.
func IsReadOnlyShellCommand(cmd string) bool {
	trimmed := strings.TrimSpace(cmd)
	if trimmed == "" {
		return false
	}
	lower := strings.ToLower(trimmed)

	// Any output redirection to a file, tee, or in-place editing mutates.
	for _, marker := range []string{">>", ">", "tee ", "dd ", "sed -i", "chmod", "chown", "mv ", "cp ", "rm ", "mkdir", "touch ", "ln "} {
		if strings.Contains(lower, marker) {
			return false
		}
	}

	readOnlyBins := map[string]bool{
		"ls": true, "cat": true, "head": true, "tail": true, "grep": true, "egrep": true, "fgrep": true,
		"rg": true, "find": true, "pwd": true, "whoami": true, "id": true, "hostname": true,
		"uname": true, "date": true, "uptime": true, "df": true, "du": true, "free": true,
		"ps": true, "env": true, "printenv": true, "which": true, "type": true, "file": true,
		"wc": true, "sort": true, "uniq": true, "cut": true, "tr": true, "stat": true,
		"echo": true, "printf": true, "basename": true, "dirname": true, "realpath": true,
		"diff": true, "comm": true, "md5sum": true, "sha256sum": true, "sha1sum": true,
		"journalctl": true, "git": true, "node": true, "python3": true, "python": true,
		"docker": true, "systemctl": true, "scorp": true, "go": true,
	}
	// Subcommand-restricted binaries: only these subcommands count as reads.
	readOnlySub := map[string]map[string]bool{
		"git":       {"status": true, "log": true, "diff": true, "show": true, "branch": true, "remote": true, "rev-parse": true, "ls-files": true, "describe": true, "tag": true, "config": true},
		"docker":    {"ps": true, "images": true, "inspect": true, "version": true, "stats": true},
		"systemctl": {"status": true, "is-active": true, "is-enabled": true, "list-units": true, "list-unit-files": true, "show": true, "cat": true},
		"go":        {"version": true, "env": true, "list": true},
		"node":      {"--version": true, "-v": true},
		"python3":   {"--version": true},
		"python":    {"--version": true},
		"scorp":     {"version": true, "eval": true},
	}

	// Split into pipeline/&&/;/|| segments; every segment must be read-only.
	segRe := strings.NewReplacer("&&", "\x00", "||", "\x00", ";", "\x00", "|", "\x00", "\n", "\x00")
	for _, seg := range strings.Split(segRe.Replace(trimmed), "\x00") {
		fields := strings.Fields(strings.TrimSpace(seg))
		if len(fields) == 0 {
			continue
		}
		bin := strings.ToLower(fields[0])
		bin = strings.TrimPrefix(bin, "sudo ") // (sudo prefix already split by Fields; safety net)
		bin = bin[strings.LastIndex(bin, "/")+1:]
		if !readOnlyBins[bin] {
			return false
		}
		if subs, ok := readOnlySub[bin]; ok {
			// find first non-flag argument as the subcommand
			sub := ""
			for _, f := range fields[1:] {
				if !strings.HasPrefix(f, "-") {
					sub = strings.ToLower(f)
					break
				}
			}
			if !subs[sub] {
				return false
			}
		}
	}
	return true
}

// confirmationDisplay renders the human-facing text of a tool call awaiting
// approval in auto mode.
func confirmationDisplay(tc ToolCall) string {
	if tc.Name == "shell" {
		return helpers.GetStringArg(tc.Args, "command", "")
	}
	argsJSON, _ := json.Marshal(tc.Args)
	return helpers.TruncateStr(tc.Name+" "+string(argsJSON), 300)
}

// PermissionDecision grades a tool call. Returns the decision (allow/ask/deny),
// a human-readable reason, and the source layer ("deterministic", "heuristic",
// "model", "fallback", "allowlist").
func PermissionDecision(toolName string, args map[string]interface{}) (decision, reason, source string) {
	if config.GetAutonomyLevel() != config.AutonomyAuto {
		return AutoAllow, "", "off"
	}

	// Fast path: read-only tools are always safe.
	if config.IsReadOnlyTool(toolName) {
		bumpAutoStat(AutoAllow)
		return AutoAllow, "read-only tool", "deterministic"
	}

	// task-plan bookkeeping tools are agent plumbing, never risky.
	switch toolName {
	case "task_plan", "complete_task", "clarify", "memory":
		bumpAutoStat(AutoAllow)
		return AutoAllow, "agent plumbing", "deterministic"
	}

	if toolName == "shell" {
		cmd := helpers.GetStringArg(args, "command", "")
		if IsDangerousCommand(cmd) {
			if autoAllowlisted(cmd) {
				bumpAutoStat(AutoAllow)
				return AutoAllow, "destructive but allowlisted (SCORP_AUTO_ALLOW)", "allowlist"
			}
			bumpAutoStat(AutoDeny)
			return AutoDeny, "destructive command in auto mode (deny unless SCORP_AUTO_ALLOW): " + helpers.TruncateStr(cmd, 120), "deterministic"
		}
		if IsReadOnlyShellCommand(cmd) {
			bumpAutoStat(AutoAllow)
			return AutoAllow, "read-only shell command", "heuristic"
		}
	}

	// Everything else: ask the cheap model.
	decision, reason, src := autoClassify(toolName, args)
	return decision, reason, src
}

// bumpAutoStat records one deterministic decision for the status counters.
func bumpAutoStat(decision string) {
	autoStat.mu.Lock()
	defer autoStat.mu.Unlock()
	switch decision {
	case AutoAllow:
		autoStat.totalSafe++
	case AutoAsk:
		autoStat.totalAsk++
	case AutoDeny:
		autoStat.totalDeny++
	}
}

// autoClassify runs the model grader with fallback bookkeeping. Uncertain
// results fail closed (ask); too many consecutive uncertain results degrade
// the session to deterministic-only behavior.
func autoClassify(toolName string, args map[string]interface{}) (decision, reason, source string) {
	autoStat.mu.Lock()
	degraded := autoStat.fallbackActive
	autoStat.mu.Unlock()

	if !degraded {
		decision, reason := AutoClassifyFunc(toolName, args)
		switch decision {
		case AutoAllow, AutoAsk, AutoDeny:
			autoStat.mu.Lock()
			autoStat.consecUncertain = 0
			switch decision {
			case AutoAllow:
				autoStat.totalSafe++
			case AutoAsk:
				autoStat.totalAsk++
			case AutoDeny:
				autoStat.totalDeny++
			}
			autoStat.mu.Unlock()
			if decision == AutoAllow {
				return AutoAllow, reason, "model"
			}
			if decision == AutoDeny {
				return AutoDeny, reason, "model"
			}
			return AutoAsk, reason, "model"
		default: // uncertain / error
			autoStat.mu.Lock()
			autoStat.consecUncertain++
			autoStat.totalUncertain++
			if autoStat.consecUncertain >= autoMaxUncertain {
				autoStat.fallbackActive = true
				log.Printf("[auto] %d consecutive uncertain classifications — falling back to supervised-style behavior for this task", autoStat.consecUncertain)
			}
			autoStat.mu.Unlock()
		}
	} else {
		autoStat.mu.Lock()
		autoStat.totalUncertain++
		autoStat.mu.Unlock()
	}

	// Fallback / uncertain: deterministic-safe stays safe, everything else asks.
	if toolName == "shell" && IsReadOnlyShellCommand(helpers.GetStringArg(args, "command", "")) {
		return AutoAllow, "fallback: read-only shell", "fallback"
	}
	if config.IsReadOnlyTool(toolName) {
		return AutoAllow, "fallback: read-only tool", "fallback"
	}
	return AutoAsk, "classifier uncertain — human confirmation required", "fallback"
}

// autoClassifyWithModel is the default model-backed grader.
func autoClassifyWithModel(toolName string, args map[string]interface{}) (decision, reason string) {
	model := models.RouteModel("chat")
	if model == nil {
		return "uncertain", "no model available"
	}

	argsJSON, _ := json.Marshal(args)
	if len(argsJSON) > 1200 {
		argsJSON = argsJSON[:1200]
	}

	prompt := fmt.Sprintf(`Grade this AI-agent tool call for safety. Autonomy context: the agent runs unattended on the user's own machine.

Tool: %s
Arguments: %s

Decisions:
- "safe": the call only reads, or writes only within obvious project/tmp workspace paths with no destructive potential.
- "ask": the call mutates state in ways a careful human might want to review (installs, deletes of specific files, config edits, publishing, process/service changes).
- "deny": clearly destructive or dangerous even if a human is watching (data loss at scale, system damage, credential exfiltration).

Return ONLY strict JSON: {"decision":"safe|ask|deny","reason":"<max 12 words>"}`, toolName, string(argsJSON))

	ctx, cancel := context.WithTimeout(context.Background(), autoClassifyTimeout)
	defer cancel()

	resp, err := models.CallModel(ctx, model, []models.ChatMessage{
		{Role: "system", Content: "You are a permission classifier. Return ONLY the JSON object, no markdown fences, no explanation."},
		{Role: "user", Content: prompt},
	})
	if err != nil {
		log.Printf("[auto] classifier model call failed: %v", err)
		return "uncertain", "classifier unavailable"
	}

	var parsed struct {
		Decision string `json:"decision"`
		Reason   string `json:"reason"`
	}
	// Tolerate markdown fences around the JSON.
	trimmed := strings.TrimSpace(resp)
	trimmed = strings.TrimPrefix(trimmed, "```json")
	trimmed = strings.TrimPrefix(trimmed, "```")
	trimmed = strings.TrimSuffix(trimmed, "```")
	trimmed = strings.TrimSpace(trimmed)
	if err := json.Unmarshal([]byte(trimmed), &parsed); err != nil {
		log.Printf("[auto] classifier returned unparseable output: %.120s", resp)
		return "uncertain", "unparseable classifier output"
	}
	switch strings.ToLower(strings.TrimSpace(parsed.Decision)) {
	case "safe":
		return AutoAllow, "model: " + parsed.Reason
	case "ask":
		return AutoAsk, "model: " + parsed.Reason
	case "deny":
		return AutoDeny, "model: " + parsed.Reason
	default:
		return "uncertain", "unknown decision value"
	}
}
