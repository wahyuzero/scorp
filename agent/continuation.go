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

// looksLikeContinuation detects if the model's response indicates intent to continue
// but didn't actually call tools (e.g., "Let me...", "I'll try...", "Mari coba...")
func looksLikeContinuation(text string) bool {
	lower := strings.ToLower(text)
	patterns := []string{
		"let me ", "i'll ", "i will ", "i'm going to ", "going to ",
		"mar i coba", "mari coba", "saya akan ", "akan coba ", "coba ",
		"saya perlu ", "perlu saya ", "perlu dicek", "perlu periksa", "harus ", "mari kita ",
		"akan saya ", "akan memeriksa", "akan mengecek ", "cek dulu ", "tunggu sebentar",
		"sekarang saya ", "sekarang kita ", "sekarang akan ", "saya buat ", "saya tulis ",
		"buat file ", "tulis file ", "berikutnya ", "selanjutnya ", "langkah berikutnya ",
		"akan saya buat ", "akan saya jalankan ", "now i will ", "next, ",
		"i'll try", "let me try", "try to ", "attempt to ",
		"next i", "then i", "now i", "continue to ",
		"proceed to ", "follow up", "followup",
		"i need to ", "still need to ", "first, let me",
		"after that", "once that", "now let's",
		"step 1", "step 2", "step 3",
		"first,", "second,", "third,",
		"i haven't", "not yet", "still working",
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
