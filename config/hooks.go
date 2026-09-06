package config

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"sync"
)

// ──────────────────────────────────────────────
// Hooks Engine (P3.12)
//
// User-configurable PreToolUse/PostToolUse hook commands with the pattern
// `tool_pattern:shell_command` (split at the FIRST colon, so the command may
// contain further colons):
//
//	SCORP_HOOKS_PRE='shell:~/.scorp/hooks/audit.sh;write_file(*secret*):block-secret-writes.sh'
//	SCORP_HOOKS_POST='*:~/.scorp/hooks/log-result.sh'
//
// A `*` (or empty) tool pattern matches every tool. Multiple hooks for one
// event run in spec order; a PreToolUse block short-circuits the rest.
// NOTE: `;` separates entries, so hook commands cannot contain a raw
// semicolon — use `&&` / `||` or a wrapper script instead.
//
// Contract (mirrors Claude Code hooks):
//   - hook receives a JSON payload on stdin:
//     {hook_event, session_id, tool_name, tool_args[, tool_result, tool_ok]}
//   - PreToolUse: exit 0 = allow (stdout becomes additional context for the
//     model), exit 2 = BLOCK (stderr is the reason shown to the model),
//     anything else = non-blocking, logged only.
//   - PostToolUse: informational — exit codes never block; stdout is appended
//     to the tool result as additional context.
//   - a hook that hangs is killed after hookTimeoutSec and treated as
//     non-blocking — a broken hook must never brick the agent.
//
// "CLAUDE.md says 'please', hooks say 'must'": hooks run at the single
// execution choke point (agent.ExecuteTool) and on the user-confirmed resume
// path, in every autonomy mode, AFTER deny rules / autonomy / path checks —
// they add enforcement, they do not relax it.
// ──────────────────────────────────────────────

// HookEntry is one parsed hook spec.
type HookEntry struct {
	Raw     string // original spec text, e.g. "shell:audit.sh"
	Pattern string // tool name pattern, or "*" for any tool
	Command string // shell command to run
}

var (
	preHooks     []HookEntry
	postHooks    []HookEntry
	hooksOnce    sync.Once
	hookEntriesN int
)

// loadHooks parses SCORP_HOOKS_PRE / SCORP_HOOKS_POST once. Invalid specs are
// logged and skipped — a typo must never brick every tool call.
func loadHooks() {
	hooksOnce.Do(func() {
		preHooks = parseHookEnv("SCORP_HOOKS_PRE")
		postHooks = parseHookEnv("SCORP_HOOKS_POST")
	})
}

// parseHookEnv reads one env var into hook entries.
func parseHookEnv(envVar string) []HookEntry {
	raw := os.Getenv(envVar)
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	var entries []HookEntry
	for _, part := range strings.Split(raw, ";") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		entry, err := ParseHookSpec(part)
		if err != nil {
			log.Printf("[hooks] ignoring invalid %s entry %q: %v", envVar, part, err)
			continue
		}
		entries = append(entries, entry)
	}
	hookEntriesN += len(entries)
	if n := len(entries); n > 0 {
		log.Printf("[hooks] loaded %d %s hook(s) from %s", n, envVar, envVar)
	}
	return entries
}

// ReloadHooks re-reads SCORP_HOOKS_PRE / SCORP_HOOKS_POST. Used by tests and
// available for future config-reload paths; in normal operation hooks load
// once at startup.
func ReloadHooks() {
	hooksOnce = sync.Once{}
	preHooks = nil
	postHooks = nil
	hookEntriesN = 0
	loadHooks()
}

// ParseHookSpec parses one `tool_pattern:shell_command` spec.
func ParseHookSpec(spec string) (HookEntry, error) {
	spec = strings.TrimSpace(spec)
	colon := strings.Index(spec, ":")
	if colon < 0 {
		return HookEntry{}, fmt.Errorf("expected tool_pattern:shell_command, got %q", spec)
	}
	pattern := strings.TrimSpace(spec[:colon])
	command := strings.TrimSpace(spec[colon+1:])
	if command == "" {
		return HookEntry{}, fmt.Errorf("empty hook command in %q", spec)
	}
	if pattern == "" {
		pattern = "*"
	}
	return HookEntry{Raw: spec, Pattern: pattern, Command: command}, nil
}

// PreToolHooks returns the parsed PreToolUse hooks.
func PreToolHooks() []HookEntry {
	loadHooks()
	return preHooks
}

// PostToolHooks returns the parsed PostToolUse hooks.
func PostToolHooks() []HookEntry {
	loadHooks()
	return postHooks
}

// HookMatches reports whether a hook entry applies to the given tool name.
// "*" (or empty) matches everything; otherwise the pattern must equal the
// tool name or, if it contains a wildcard, glob-match it.
func HookMatches(entry HookEntry, toolName string) bool {
	if entry.Pattern == "*" || entry.Pattern == "" {
		return true
	}
	if !strings.Contains(entry.Pattern, "*") {
		return entry.Pattern == toolName
	}
	ok, err := path.Match(entry.Pattern, toolName)
	return err == nil && ok
}
