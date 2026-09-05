package agent

import (
	"fmt"
	"strings"
	"unicode"
)

// ──────────────────────────────────────────────
// Agent Task Completion & Continuation Protocol
// Robust grammatical intent classifier and continuation directive handler
// replacing fragile string-matching heuristics.
// ──────────────────────────────────────────────

// isContinuationDirective checks if user message is an explicit command to continue
// an ongoing or prematurely halted task (e.g. "lanjutkan", "continue", "next").
func isContinuationDirective(msg string) bool {
	lower := strings.ToLower(strings.TrimSpace(msg))
	// Strip trailing punctuation
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

	// Prefixes like "lanjutkan investigasi", "continue with the search"
	if strings.HasPrefix(lower, "lanjutkan ") || strings.HasPrefix(lower, "continue ") ||
		strings.HasPrefix(lower, "teruskan ") || strings.HasPrefix(lower, "lanjut ") {
		// If short (< 35 chars), it's a continuation directive rather than a fresh prompt
		if len(lower) < 35 {
			return true
		}
	}

	return false
}

// isPureInformationalQuery checks if the user message is an informational,
// conceptual, or educational question that does NOT require tool execution.
func isPureInformationalQuery(msg string) bool {
	lower := strings.TrimSpace(strings.ToLower(msg))

	// If it asks to execute, create, modify, delete, run, or search files/tools, it is NOT purely informational
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

	// Explanatory question prefixes
	infoPrefixes := []string{
		"jelaskan ", "apa itu ", "kenapa ", "mengapa ", "bagaimana cara ", "apa saja ",
		"explain ", "what is ", "why ", "how does ", "how to ", "what are ",
		"ceritakan ", "pendapatmu ", "menurutmu ", "bedanya ", "perbedaan ",
	}
	for _, p := range infoPrefixes {
		if strings.HasPrefix(lower, p) {
			return true
		}
	}
	return false
}

// hasCompletionIndicators detects if the model explicitly declares the requested task finished.
func hasCompletionIndicators(text string) bool {
	lower := strings.ToLower(text)
	indicators := []string{
		"selesai", "berhasil", "terverifikasi", "sudah dihapus", "sudah dibuat",
		"sudah selesai", "done", "completed", "verified", "tidak ada tindakan lebih lanjut",
		"telah berhasil", "semua selesai", "hasil verifikasi", "ringkasnya:", "kesimpulannya",
		"all done", "all tasks completed", "successfully finished", "berikut kesimpulan",
		"berikut rangkuman", "hasil investigasi",
	}
	for _, ind := range indicators {
		if strings.Contains(lower, ind) {
			return true
		}
	}
	return false
}

// hasForwardIntent detects if the text expresses an intention to perform an unexecuted action
// using robust semantic and grammatical clause matching.
func hasForwardIntent(text string) bool {
	lower := strings.ToLower(text)

	// 1. Direct explicit future action phrases
	explicitPhrases := []string{
		"sekarang saya ", "sekarang aku ", "sekarang kita ", "sekarang akan ",
		"sekarang jalankan", "sekarang hapus", "sekarang buat", "sekarang bersihkan", "sekarang compile",
		"sekarang lihat", "sekarang cek", "sekarang periksa", "sekarang cari", "sekarang baca",
		"langkah berikutnya:", "tahap selanjutnya:", "tahap berikutnya", "langkah selanjutnya",
		"akan saya ", "akan aku ", "akan kita ", "akan coba ", "akan jalankan", "akan periksa",
		"akan membuat", "akan menjalankan", "akan membaca", "akan menulis", "akan menguji", "akan mencari",
		"saya akan ", "aku akan ", "kita akan ", "kami akan ",
		"saya coba ", "aku coba ", "kita coba ",
		"saya lanjutkan ", "aku lanjutkan ", "mari kita lanjutkan ",
		"saya cari ", "aku cari ", "saya lihat ", "aku lihat ", "saya periksa ", "aku periksa ",
		"saya baca ", "aku baca ", "saya tulis ", "aku tulis ", "saya buat ", "aku buat ",
		"saya jalankan ", "aku jalankan ", "saya hapus ", "aku hapus ",
		"perlu saya ", "perlu aku ", "perlu dicek", "perlu periksa", "perlu diambil",
		"kucoba ", "kucari ", "kulihat ", "kuperiksa ", "kubaca ", "kuhapus ",
		"now deleting", "now removing", "now running", "now executing", "now cleaning", "now creating",
		"now checking", "now inspecting", "now reading", "now searching", "now looking",
		"deleting the file", "removing the file", "creating the file", "running the script",
		"i will now ", "now i will ", "now i'll ", "next, i will ", "let me execute",
		"let me check", "let me inspect", "let me read", "let me search", "let me run",
		"let's check", "let's inspect", "i'm going to ", "proceed to ", "still need to ",
	}

	for _, p := range explicitPhrases {
		if strings.Contains(lower, p) {
			return true
		}
	}

	// 2. Clause-level grammatical classifier
	// Break text into sentences/clauses
	clauses := strings.FieldsFunc(lower, func(r rune) bool {
		return r == '.' || r == '\n' || r == '!' || r == '?' || r == ';' || r == '—' || r == '-'
	})

	for _, clause := range clauses {
		c := strings.TrimSpace(clause)
		if len(c) < 10 {
			continue
		}

		// Check if clause has Actor/Temporal + Action Verb + Target Entity
		hasActorOrTemporal := containsAnyWord(c, []string{
			"aku", "saya", "kita", "kami", "gue", "mari",
			"sekarang", "nanti", "berikutnya", "selanjutnya", "next", "now", "then",
			"i", "we", "let's",
		})

		hasActionVerb := containsAnyWord(c, []string{
			"akan", "mau", "coba", "perlu", "harus", "bisa", "lanjut", "lanjutkan",
			"lihat", "cek", "periksa", "baca", "tulis", "buat", "hapus", "jalankan",
			"eksekusi", "cari", "ambil", "fetch", "telusuri", "investigasi", "buka",
			"will", "shall", "try", "check", "inspect", "read", "write", "delete",
			"run", "execute", "search", "find", "open", "investigate",
		})

		hasTargetEntity := containsAnyWord(c, []string{
			"log", "logs", "server", "file", "berkas", "url", "link", "web", "website",
			"data", "sesi", "session", "riwayat", "history", "code", "kode", "script",
			"output", "hasil", "error", "status", "postingan", "tweet", "source", "repo",
			"proses", "process", "jejak", "docker", "service", "system", "sistem",
		})

		if hasActorOrTemporal && hasActionVerb && hasTargetEntity {
			return true
		}
	}

	return false
}

// containsAnyWord checks if a string contains any of the target words bounded by word boundaries
func containsAnyWord(s string, words []string) bool {
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

// looksLikeContinuation detects if the model output indicates unfinished work.
func looksLikeContinuation(text string) bool {
	if hasForwardIntent(text) {
		return true
	}
	if hasCompletionIndicators(text) {
		return false
	}

	lower := strings.ToLower(text)
	patterns := []string{
		"let me ", "i'll ", "i will ", "i'm going to ",
		"saya akan ", "akan saya ", "akan coba ", "coba jalankan ", "coba periksa ",
		"aku akan ", "akan aku ", "aku coba ", "coba aku ",
		"saya perlu ", "aku perlu ", "perlu saya ", "perlu aku ", "perlu dicek", "perlu periksa",
		"now i will ", "next, let me", "i'll try", "let me try",
		"proceed to ", "still need to ", "sekarang aku ", "sekarang saya ",
	}
	for _, p := range patterns {
		if strings.Contains(lower, p) {
			return true
		}
	}
	return false
}

// countStepsInMessage calculates the number of distinct operations requested in the prompt.
func countStepsInMessage(msg string) int {
	lower := strings.ToLower(msg)

	// 1. Numbered lists: "1.", "1)", "2.", etc.
	numberedCount := 0
	for i := 1; i <= 20; i++ {
		if strings.Contains(lower, fmt.Sprintf("%d.", i)) || strings.Contains(lower, fmt.Sprintf("%d)", i)) {
			numberedCount++
		}
	}
	if numberedCount >= 2 {
		return numberedCount
	}

	// 2. Sequential connectors
	seqCount := 0
	seqConnectors := []string{
		" then ", " after that ", " next,", "\nthen ", "→",
		" lalu ", " kemudian ", " setelah ", " setelah itu ",
		", lalu ", ", kemudian ",
	}
	for _, sc := range seqConnectors {
		seqCount += strings.Count(lower, sc)
	}

	if seqCount >= 1 {
		return seqCount + 1
	}

	return 0
}

// mentionsBrowserTask checks if user message references browser actions
func mentionsBrowserTask(msg string) bool {
	lower := strings.ToLower(msg)
	keywords := []string{
		"screenshot", "login", "log in", "sign in", "scrape",
		"ambil gambar", "tangkap layar", "buka web",
		"buka halaman", "cek halaman", "cek web",
	}
	for _, kw := range keywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}

// screenshotWasTaken checks if a browser.screenshot tool call exists in history
func screenshotWasTaken(history []AgentMessage) bool {
	for _, msg := range history {
		content, ok := msg.Content.(string)
		if !ok {
			continue
		}
		if strings.Contains(content, "screenshot") || strings.Contains(content, "browser(screenshot)") {
			return true
		}
	}
	return false
}
