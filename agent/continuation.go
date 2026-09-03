package agent

import (
	"fmt"
	"strings"
)

// ──────────────────────────────────────────────
// Agent Continuation & Intent Heuristics
// Detects if the model's response indicates intent to continue
// but didn't actually call tools (preventing premature stopping).
// ──────────────────────────────────────────────

// hasCompletionIndicators detects if the model has explicitly finished the task.
func hasCompletionIndicators(text string) bool {
	lower := strings.ToLower(text)
	indicators := []string{
		"selesai", "berhasil", "terverifikasi", "sudah dihapus", "sudah dibuat",
		"sudah selesai", "done", "completed", "verified", "tidak ada tindakan lebih lanjut",
		"telah berhasil", "semua selesai", "hasil verifikasi", "ringkasnya:", "kesimpulannya",
	}
	for _, ind := range indicators {
		if strings.Contains(lower, ind) {
			return true
		}
	}
	return false
}

// looksLikeContinuation detects if the model's response indicates intent to continue
// but didn't actually call tools (e.g., "Let me...", "I'll try...", "Mari coba...")
func looksLikeContinuation(text string) bool {
	if hasCompletionIndicators(text) {
		return false
	}

	lower := strings.ToLower(text)
	patterns := []string{
		"let me ", "i'll ", "i will ", "i'm going to ", "going to ",
		"mari coba", "saya akan ", "akan coba ", "coba jalankan ", "coba periksa ",
		"saya perlu ", "perlu saya ", "perlu dicek", "perlu periksa", "mari kita ",
		"akan saya ", "akan memeriksa", "akan mengecek ", "cek dulu ", "tunggu sebentar",
		"sekarang saya ", "sekarang kita ", "sekarang akan ", "sekarang compile", "sekarang jalankan",
		"saya susun ulang", "saya perbaiki", "saya buat file", "saya tulis file",
		"langkah berikutnya:", "tahap selanjutnya:", "akan saya buat ", "akan saya jalankan ",
		"now i will ", "next, let me", "i'll try", "let me try", "try to execute",
		"next i will", "then i will", "continue to ", "proceed to ",
		"i need to ", "still need to ", "first, let me", "still working on",
	}
	for _, p := range patterns {
		if strings.Contains(lower, p) {
			return true
		}
	}
	return false
}

// countStepsInMessage estimates how many distinct steps/actions the user's message asks for.
func countStepsInMessage(msg string) int {
	lower := strings.ToLower(msg)
	count := 0

	// Count numbered list items: "1)", "1.", "2)", "2.", etc.
	for i := 1; i <= 20; i++ {
		if strings.Contains(lower, fmt.Sprintf("%d.", i)) || strings.Contains(lower, fmt.Sprintf("%d)", i)) {
			count++
		}
	}
	if count >= 2 {
		return count
	}

	// Count "then" / "after that" / "next" chains (EN + ID)
	count += strings.Count(lower, " then ")
	count += strings.Count(lower, " after that ")
	count += strings.Count(lower, " next,")
	count += strings.Count(lower, "\nthen ")
	count += strings.Count(lower, "→")
	count += strings.Count(lower, " lalu ")
	count += strings.Count(lower, " kemudian ")
	count += strings.Count(lower, " setelah ")
	count += strings.Count(lower, " lalu\n")
	count += strings.Count(lower, " kemudian\n")
	count += strings.Count(lower, " setelah\n")

	if count >= 2 {
		return count + 1
	}

	return 0
}

// mentionsBrowserTask checks if user message references browser actions
func mentionsBrowserTask(msg string) bool {
	lower := strings.ToLower(msg)
	keywords := []string{
		"screenshot", "login", "log in", "sign in", "scrape",
		"ambil gambar", "tangkap layar", "masuk", "buka web",
		"buka halaman", "cek halaman", "cek web", "navigate",
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
