package tools

import (
	"os"
	"regexp"
	"strings"
)

// ──────────────────────────────────────────────
// Outbound Secret Redactor (PicoClaw Parity)
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

	// 1. Redact known active environment secrets from process environment
	knownEnvKeys := []string{
		"COMMAND_CODE_API_KEY",
		"OPENAI_API_KEY",
		"ANTHROPIC_API_KEY",
		"GEMINI_API_KEY",
		"GOOGLE_API_KEY",
		"DEEPSEEK_API_KEY",
		"GROQ_API_KEY",
		"OPENROUTER_API_KEY",
		"TELEGRAM_BOT_TOKEN",
		"COOLIFY_API_TOKEN",
		"FIRECRAWL_API_KEY",
		"TAVILY_API_KEY",
		"GITHUB_TOKEN",
	}

	for _, key := range knownEnvKeys {
		if val := os.Getenv(key); len(val) >= 8 {
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
