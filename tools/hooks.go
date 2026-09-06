package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"scorp-agent/config"
)

// ──────────────────────────────────────────────
// Hook Runner (P3.12)
//
// Executes user-configured PreToolUse/PostToolUse hooks (config/hooks.go) at
// the single execution choke point. Protocol per hook invocation:
//
//   stdin  : JSON payload {hook_event, session_id, tool_name, tool_args
//            [, tool_result, tool_ok]} — secrets redacted
//   env    : SCORP_HOOK_EVENT / SCORP_HOOK_TOOL / SCORP_HOOK_SESSION
//   exit 0 : allow; stdout is additional context for the model
//   exit 2 : PreToolUse BLOCKS the call (stderr = reason shown to the model);
//            PostToolUse is informational, logged + surfaced as context
//   other  : non-blocking, stderr logged only
//   hang   : killed after hookRunTimeout, treated as non-blocking — a broken
//            hook must never brick the agent
//
// Hooks run OUTSIDE the bwrap sandbox: they are admin-configured trusted
// commands (audit loggers, policy checks), not model-driven work.
// ──────────────────────────────────────────────

// hookRunTimeout bounds each hook invocation. Overridable in tests.
var hookRunTimeout = 10 * time.Second

// hookExitTimeout marks a killed hanging hook (impossible as a real exit code
// since bash caps at 255 but >2 keeps it out of the semantic 0/2 range).
const hookExitTimeout = 99

type hookPayload struct {
	HookEvent  string                 `json:"hook_event"`
	SessionID  string                 `json:"session_id"`
	ToolName   string                 `json:"tool_name"`
	ToolArgs   map[string]interface{} `json:"tool_args"`
	ToolResult string                 `json:"tool_result,omitempty"`
	ToolOK     bool                   `json:"tool_ok,omitempty"`
}

// runHookCommand invokes one hook command with the JSON payload on stdin.
// Returns trimmed stdout, trimmed stderr and the exit code (hookExitTimeout on
// hang/kill).
func runHookCommand(hook config.HookEntry, event, toolName string, args map[string]interface{}, chatID int64, result string, ok bool) (string, string, int) {
	payload := hookPayload{
		HookEvent:  event,
		SessionID:  fmt.Sprintf("%d", chatID),
		ToolName:   toolName,
		ToolArgs:   args,
		ToolResult: result,
		ToolOK:     ok,
	}
	pj, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Sprintf("hook payload marshal error: %v", err), -1
	}
	// Defense-in-depth: hooks are user-configured, but their stdin/logs must
	// never become a secret-exfiltration channel either.
	pj = []byte(RedactSecrets(string(pj)))

	ctx, cancel := context.WithTimeout(context.Background(), hookRunTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "bash", "-c", hook.Command)
	cmd.Stdin = bytes.NewReader(pj)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	cmd.Env = append(os.Environ(),
		"SCORP_HOOK_EVENT="+event,
		"SCORP_HOOK_TOOL="+toolName,
		"SCORP_HOOK_SESSION="+payload.SessionID,
	)
	// Own process group so a hanging hook's children die with it.
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if err := cmd.Start(); err != nil {
		return "", err.Error(), -1
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	select {
	case <-ctx.Done():
		syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		<-done
		return strings.TrimSpace(outBuf.String()), strings.TrimSpace(errBuf.String()), hookExitTimeout
	case err := <-done:
		code := 0
		if ee, ok := err.(*exec.ExitError); ok {
			code = ee.ExitCode()
		} else if err != nil {
			return strings.TrimSpace(outBuf.String()), err.Error(), -1
		}
		return strings.TrimSpace(outBuf.String()), strings.TrimSpace(errBuf.String()), code
	}
}

func appendHookContext(a, b string) string {
	if b == "" {
		return a
	}
	if a == "" {
		return b
	}
	return a + "\n" + b
}

func truncateHookLog(s string) string {
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) > 200 {
		return s[:200] + "…"
	}
	return s
}

// RunPreToolUseHooks runs every matching SCORP_HOOKS_PRE hook. Returns
// accumulated additional context, whether the call is blocked, and the block
// reason (stderr of the exit-2 hook).
func RunPreToolUseHooks(toolName string, args map[string]interface{}, chatID int64) (context string, blocked bool, reason string) {
	hooks := config.PreToolHooks()
	if len(hooks) == 0 {
		return "", false, ""
	}
	for _, hook := range hooks {
		if !config.HookMatches(hook, toolName) {
			continue
		}
		stdout, stderr, code := runHookCommand(hook, "pre_tool_use", toolName, args, chatID, "", false)
		switch {
		case code == 2:
			r := stderr
			if r == "" {
				r = fmt.Sprintf("blocked by PreToolUse hook '%s'", hook.Raw)
			}
			log.Printf("[hooks] pre_tool_use hook BLOCKED %s: %s", toolName, truncateHookLog(r))
			return context, true, r
		case code == 0:
			context = appendHookContext(context, stdout)
		case code == hookExitTimeout:
			log.Printf("[hooks] pre_tool_use hook timed out after %s: %s", hookRunTimeout, hook.Raw)
			context = appendHookContext(context,
				fmt.Sprintf("⚠️ hook '%s' timed out after %s (non-blocking)", hook.Raw, hookRunTimeout))
		default:
			log.Printf("[hooks] pre_tool_use hook %s exited %d (non-blocking): %s",
				hook.Raw, code, truncateHookLog(stderr))
		}
	}
	return context, false, ""
}

// RunPostToolUseHooks runs every matching SCORP_HOOKS_POST hook. Post hooks
// are informational: exit codes never block, stdout becomes additional context
// appended to the tool result.
func RunPostToolUseHooks(toolName string, args map[string]interface{}, result string, ok bool, chatID int64) (context string) {
	hooks := config.PostToolHooks()
	if len(hooks) == 0 {
		return ""
	}
	for _, hook := range hooks {
		if !config.HookMatches(hook, toolName) {
			continue
		}
		stdout, stderr, code := runHookCommand(hook, "post_tool_use", toolName, args, chatID, result, ok)
		switch {
		case code == 2:
			log.Printf("[hooks] post_tool_use hook %s rejected the result (tool already executed, non-blocking): %s",
				hook.Raw, truncateHookLog(stderr))
			context = appendHookContext(context,
				fmt.Sprintf("⚠️ post-tool-use hook '%s' flagged this result: %s", hook.Raw, stderr))
		case code == 0:
			context = appendHookContext(context, stdout)
		case code == hookExitTimeout:
			log.Printf("[hooks] post_tool_use hook timed out after %s: %s", hookRunTimeout, hook.Raw)
			context = appendHookContext(context,
				fmt.Sprintf("⚠️ hook '%s' timed out after %s (non-blocking)", hook.Raw, hookRunTimeout))
		default:
			log.Printf("[hooks] post_tool_use hook %s exited %d (non-blocking): %s",
				hook.Raw, code, truncateHookLog(stderr))
		}
	}
	return context
}
