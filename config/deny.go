package config

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
)

// ──────────────────────────────────────────────
// Deny-Rule Engine (P0.2)
//
// User-configurable hard denials with the pattern `tool(param:regex)`:
//   shell(command:mkfs|dd if=/dev/zero)   — regex on the shell command arg
//   write_file(path:^/etc/)               — regex on a named argument
//   *(*:AKIA[0-9A-Z]{16})                 — any tool, any string argument
//
// Contract: deny rules are evaluated BEFORE the confirmation gates and do NOT
// consult the autonomy level, so they hold in readonly, supervised AND yolo.
// A user-confirmed resume cannot bypass them either — tools.ExecuteShell
// re-checks them on the confirmation path. This mirrors Claude Code's deny
// rules, which stay in force under bypassPermissions.
//
// Configuration: SCORP_DENY_RULES env var, semicolon-separated specs.
// The value part is a Go RE2 regex, unanchored — anchor with ^ / $ explicitly.
// ──────────────────────────────────────────────

// DenyRule is one parsed and compiled deny rule.
type DenyRule struct {
	Raw    string // original spec text, e.g. "shell(command:^rm\s)"
	Tool   string // tool name, or "*" for any tool
	Param  string // argument name, or "*" for any string argument
	Regexp *regexp.Regexp
}

var (
	denyRules     []DenyRule
	denyRulesOnce sync.Once
)

// loadDenyRules parses SCORP_DENY_RULES once. Invalid specs are logged and
// skipped — a typo must never brick every tool call.
func loadDenyRules() {
	denyRulesOnce.Do(func() {
		spec := os.Getenv("SCORP_DENY_RULES")
		if strings.TrimSpace(spec) == "" {
			return
		}
		for _, part := range strings.Split(spec, ";") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			rule, err := ParseDenyRule(part)
			if err != nil {
				log.Printf("[deny] ignoring invalid deny rule %q: %v", part, err)
				continue
			}
			denyRules = append(denyRules, rule)
		}
		if n := len(denyRules); n > 0 {
			log.Printf("[deny] loaded %d deny rule(s) from SCORP_DENY_RULES", n)
		}
	})
}

// ReloadDenyRules re-reads SCORP_DENY_RULES. Used by tests and available for
// future config-reload paths; in normal operation rules load once at startup.
func ReloadDenyRules() {
	denyRulesOnce = sync.Once{}
	denyRules = nil
	loadDenyRules()
}

// ParseDenyRule parses one `tool(param:regex)` spec.
func ParseDenyRule(spec string) (DenyRule, error) {
	spec = strings.TrimSpace(spec)
	open := strings.Index(spec, "(")
	if open <= 0 || !strings.HasSuffix(spec, ")") {
		return DenyRule{}, fmt.Errorf("expected tool(param:regex), got %q", spec)
	}
	tool := strings.TrimSpace(spec[:open])
	inner := spec[open+1 : len(spec)-1]
	colon := strings.Index(inner, ":")
	if colon <= 0 {
		return DenyRule{}, fmt.Errorf("expected param:regex inside parentheses, got %q", inner)
	}
	param := strings.TrimSpace(inner[:colon])
	pattern := inner[colon+1:]
	if tool == "" || param == "" || pattern == "" {
		return DenyRule{}, fmt.Errorf("empty tool/param/regex in %q", spec)
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return DenyRule{}, fmt.Errorf("bad regex %q: %w", pattern, err)
	}
	return DenyRule{Raw: spec, Tool: tool, Param: param, Regexp: re}, nil
}

// CheckDenyRules reports whether a tool call matches any deny rule. Args values
// are inspected as strings only (command, path, url, …). Returns the matching
// rule spec for the denial message.
func CheckDenyRules(toolName string, args map[string]interface{}) (bool, string) {
	loadDenyRules()
	for _, rule := range denyRules {
		if rule.Tool != "*" && rule.Tool != toolName {
			continue
		}
		for key, val := range args {
			if rule.Param != "*" && rule.Param != key {
				continue
			}
			s, ok := val.(string)
			if !ok || s == "" {
				continue
			}
			if rule.Regexp.MatchString(s) {
				return true, fmt.Sprintf(
					"Denied by deny-rule '%s': %s(%s) matched /%s/. Deny rules apply in ALL autonomy modes (including YOLO) and cannot be bypassed by user confirmation.",
					rule.Raw, toolName, key, rule.Regexp.String())
			}
		}
	}
	return false, ""
}
