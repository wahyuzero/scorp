package agent

import (
	"strings"
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

	// Conversational greetings or casual remarks that require no actions
	greetings := []string{
		"hai", "halo", "hello", "hi", "hey", "hei", "pagi", "siang", "sore", "malam",
		"selamat pagi", "selamat siang", "selamat sore", "selamat malam", "assalamualaikum",
		"tes", "test", "ping", "hola", "sup", "yo",
	}
	for _, g := range greetings {
		if lower == g {
			return true
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
