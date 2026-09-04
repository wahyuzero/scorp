package agent

import (
	"fmt"
	"strings"
)

// ──────────────────────────────────────────────
// Agent Task Completion & Continuation Protocol
// Replaces fragile string-matching heuristics with structural validation
// and autonomous system self-correction prompts.
// ──────────────────────────────────────────────

// isPureInformationalQuery checks if the user message is an informational,
// conceptual, or educational question that does NOT require tool execution.
// This prevents false positives where explanatory answers get mistaken for pending tasks.
func isPureInformationalQuery(msg string) bool {
	lower := strings.TrimSpace(strings.ToLower(msg))

	// If it asks to execute, create, modify, delete, run, or search files/tools, it is NOT purely informational
	actionKeywords := []string{
		"buat", "hapus", "jalankan", "tulis", "ubah", "edit", "ganti", "bersihkan",
		"create", "delete", "write", "run", "execute", "modify", "remove", "clean",
		"cek status", "cek sistem", "check system", "test", "uji", "deploy", "install",
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

// hasForwardIntent detects if the model states an intention to perform an action
// right now instead of having called the tool (violating Action-First protocol).
func hasForwardIntent(text string) bool {
	lower := strings.ToLower(text)

	// Explicit forward action verbs without tool execution
	forwardPhrases := []string{
		"sekarang saya ", "sekarang kita ", "sekarang akan ", "sekarang jalankan",
		"sekarang hapus", "sekarang buat", "sekarang bersihkan", "sekarang compile",
		"langkah berikutnya:", "tahap selanjutnya:", "tahap berikutnya",
		"akan saya buat", "akan saya jalankan", "akan saya hapus", "akan saya bersihkan",
		"akan menghapus", "akan membersihkan", "akan membuat", "akan menjalankan",
		"akan membaca", "akan menulis", "akan menguji",
		"saya tulis ", "saya buatkan ", "saya jalankan ", "saya hapus ",
		"now deleting", "now removing", "now running", "now executing", "now cleaning", "now creating",
		"deleting the file", "removing the file", "creating the file", "running the script",
		"i will now ", "now i will ", "now i'll ", "next, i will ", "let me execute",
	}

	for _, p := range forwardPhrases {
		if strings.Contains(lower, p) {
			return true
		}
	}
	return false
}

// hasCompletionIndicators detects if the model explicitly declares the requested task finished.
func hasCompletionIndicators(text string) bool {
	if hasForwardIntent(text) {
		return false
	}

	lower := strings.ToLower(text)
	indicators := []string{
		"selesai", "berhasil", "terverifikasi", "sudah dihapus", "sudah dibuat",
		"sudah selesai", "done", "completed", "verified", "tidak ada tindakan lebih lanjut",
		"telah berhasil", "semua selesai", "hasil verifikasi", "ringkasnya:", "kesimpulannya",
		"all done", "all tasks completed", "successfully finished",
	}
	for _, ind := range indicators {
		if strings.Contains(lower, ind) {
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

	// Short unfinished intent markers
	patterns := []string{
		"let me ", "i'll ", "i will ", "i'm going to ",
		"saya akan ", "akan saya ", "akan coba ", "coba jalankan ", "coba periksa ",
		"saya perlu ", "perlu saya ", "perlu dicek", "perlu periksa",
		"now i will ", "next, let me", "i'll try", "let me try",
		"proceed to ", "still need to ",
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
