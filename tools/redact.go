package tools

import (
	"os"
	"regexp"
	"strings"
)

// ──────────────────────────────────────────────
// Outbound Secret Redactor
// Sanitizes tool outputs before feeding them to the LLM or chat history
// to prevent accidental credential leakage.
// ──────────────────────────────────────────────

var secretRegexes = []*regexp.Regexp{
	// Private Keys (RSA, EC, OpenSSH, PGP)
	regexp.MustCompile(`(?s)-----BEGIN[ A-Z0-9_-]*PRIVATE KEY-----.*?-----END[ A-Z0-9_-]*PRIVATE KEY-----`),

	// Provider Specific Tokens
	regexp.MustCompile(`\b(sk-[a-zA-Z0-9_-]{20,})\b`),         // OpenAI, Groq, Anthropic
	regexp.MustCompile(`\b(user_[a-zA-Z0-9_-]{40,})\b`),       // Command Code tokens
	regexp.MustCompile(`\b(gh[pousr]_[A-Za-z0-9_]{36,})\b`),   // GitHub tokens
	regexp.MustCompile(`\b(github_pat_[A-Za-z0-9_]{60,})\b`),  // GitHub fine-grained PAT
	regexp.MustCompile(`\b(AIza[0-9A-Za-z-_]{35})\b`),         // Google / Gemini API Keys
	regexp.MustCompile(`\b(glpat-[0-9a-zA-Z_-]{20,})\b`),      // GitLab PAT
	regexp.MustCompile(`\b(xox[baprs]-[0-9a-zA-Z]{10,48})\b`), // Slack tokens

	// HTTP Authorization Headers
	regexp.MustCompile(`(?i)(authorization:\s*bearer\s+)([a-zA-Z0-9_.-]{20,})`),

	// Key-value pairs in config/env formats
	regexp.MustCompile(`(?i)\b(API[_-]?KEY|SECRET|PASSWORD|TOKEN|ACCESS[_-]?KEY|AUTH[_-]?TOKEN)\s*([:=])\s*(["']?)([^\s"']{8,})`),
}

// RedactSecrets scans a text string and replaces detected credentials with [REDACTED_SECRET].
func RedactSecrets(input string) string {
	if len(input) == 0 {
		return input
	}

	result := input

	// 1. Redact ANY environment variable whose name implies a secret, key, token, or password
	for _, envEntry := range os.Environ() {
		parts := strings.SplitN(envEntry, "=", 2)
		if len(parts) != 2 {
			continue
		}
		k := strings.ToUpper(parts[0])
		val := strings.TrimSpace(parts[1])
		if len(val) < 6 {
			continue
		}
		// Match generic secret suffixes/prefixes (e.g. AGENTON_API_KEY, DB_PASSWORD, MY_AUTH_TOKEN)
		if strings.Contains(k, "_KEY") || strings.Contains(k, "_TOKEN") ||
			strings.Contains(k, "_SECRET") || strings.Contains(k, "_PASSWORD") ||
			strings.Contains(k, "APIKEY") || strings.Contains(k, "TOKEN") ||
			strings.Contains(k, "SECRET") || strings.Contains(k, "PASSWORD") {
			result = strings.ReplaceAll(result, val, "[REDACTED_SECRET]")
		}
	}

	// 2. Pattern-based redactions
	// Private keys
	result = secretRegexes[0].ReplaceAllString(result, "[REDACTED_PRIVATE_KEY]")

	// Provider tokens
	for i := 1; i <= 7; i++ {
		result = secretRegexes[i].ReplaceAllString(result, "[REDACTED_SECRET]")
	}

	// Bearer tokens
	result = secretRegexes[8].ReplaceAllString(result, "${1}[REDACTED_SECRET]")

	// Key-value password/secret pairs
	result = secretRegexes[9].ReplaceAllString(result, "${1}${2}${3}[REDACTED_SECRET]")

	return result
}
