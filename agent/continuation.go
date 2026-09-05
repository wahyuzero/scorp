package agent

import (
	"fmt"
	"strings"
	"unicode"
)

// ──────────────────────────────────────────────
// Agent Task Completion & Directive Protocol
// Clean, non-heuristic helpers for user directives.
// All continuation decisions are governed by the Task State Machine and complete_task tool.
// ──────────────────────────────────────────────

// IsContinuationDirective checks if the user's message is an explicit command
// telling the agent to resume/continue (e.g. "lanjutkan", "continue", "next", "teruskan").
func IsContinuationDirective(msg string) bool {
	lower := strings.ToLower(strings.TrimSpace(msg))
	lower = strings.TrimRight(lower, ".!?,; \t\n")

	exactCommands := []string{
		"lanjutkan", "lanjut", "continue", "next", "teruskan", "proceed",
		"keep going", "go on", "gas", "lanjut dong", "coba lanjutkan",
		"oke lanjutkan", "ok lanjutkan", "ayoo lanjutkan", "lanjutkan lagi",
		"teruskan lagi", "lanjutkan ya", "lanjut ya", "continue task",
	}
	for _, cmd := range exactCommands {
		if lower == cmd {
			return true
		}
	}

	if strings.HasPrefix(lower, "lanjutkan ") || strings.HasPrefix(lower, "continue ") ||
		strings.HasPrefix(lower, "teruskan ") || strings.HasPrefix(lower, "lanjut ") {
		if len(lower) < 35 {
			return true
		}
	}

	return false
}

// IsPureInformationalQuery checks if the query is a simple conversational/conceptual question
// that does not involve system actions, filesystem, or deep multi-step tools.
func IsPureInformationalQuery(msg string) bool {
	lower := strings.TrimSpace(strings.ToLower(msg))

	actionKeywords := []string{
		"buat", "hapus", "jalankan", "tulis", "ubah", "edit", "ganti", "bersihkan",
		"create", "delete", "write", "run", "execute", "modify", "remove", "clean",
		"cek status", "cek sistem", "check system", "test", "uji", "deploy", "install",
		"investigasi", "investigate", "periksa", "inspect", "telusuri", "gali lebih lanjut",
		"/tmp", ".txt", ".go", ".py", ".sh", ".json", ".md",
	}
	for _, kw := range actionKeywords {
		if strings.Contains(lower, kw) {
			return false
		}
	}

	infoPrefixes := []string{
		"jelaskan ", "apa itu ", "kenapa ", "mengapa ", "bagaimana cara ", "apa saja ",
		"explain ", "what is ", "why ", "how does ", "how to ", "what are ",
		"ceritakan ", "pendapatmu ", "menurutmu ", "bedanya ", "perbedaan ", "apa perbedaan ",
		"bandingkan ", "compare ", "mana yang lebih ", "which is ",
	}
	for _, p := range infoPrefixes {
		if strings.HasPrefix(lower, p) {
			return true
		}
	}
	return false
}

// CountStepsInMessage counts requested steps if user formatted a numbered list or sequential connectors.
func CountStepsInMessage(msg string) int {
	lower := strings.ToLower(msg)

	numberedCount := 0
	for i := 1; i <= 20; i++ {
		if strings.Contains(lower, fmt.Sprintf("%d.", i)) || strings.Contains(lower, fmt.Sprintf("%d)", i)) {
			numberedCount++
		}
	}
	if numberedCount >= 2 {
		return numberedCount
	}

	seqCount := 0
	seqConnectors := []string{
		" then ", " after that ", " next,", "\nthen ", "→",
		" lalu ", " kemudian ", " setelah ", " setelah itu ",
	}
	for _, sc := range seqConnectors {
		seqCount += strings.Count(lower, sc)
	}

	if seqCount >= 1 {
		return seqCount + 1
	}

	return 0
}

// ContainsAnyWord checks if a string contains any target token bounded by non-alphanumeric chars
func ContainsAnyWord(s string, words []string) bool {
	fields := strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' && r != '\''
	})

	fieldSet := make(map[string]bool, len(fields))
	for _, f := range fields {
		fieldSet[f] = true
	}

	for _, w := range words {
		if fieldSet[w] {
			return true
		}
	}
	return false
}
