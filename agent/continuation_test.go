package agent

import (
	"testing"
)

func TestContinuationDirectives(t *testing.T) {
	positives := []string{
		"Lanjutkan",
		"lanjutkan",
		"Lanjut",
		"lanjut",
		"continue",
		"Continue",
		"next",
		"teruskan",
		"proceed",
		"keep going",
		"gas",
		"lanjut dong",
		"coba lanjutkan",
		"oke lanjutkan",
		"ok lanjutkan",
		"lanjutkan investigasi",
		"continue search",
	}

	for _, p := range positives {
		if !isContinuationDirective(p) {
			t.Errorf("Expected isContinuationDirective(%q) = true, got false", p)
		}
	}

	negatives := []string{
		"Halo apa kabar?",
		"Jelaskan cara kerja Docker",
		"Buat file app.py",
		"Siapa penemu bahasa Go?",
		"lanjutkan membuatkan sistem microservices lengkap dengan database postgres dan redis di folder /var/www/myproject",
	}

	for _, n := range negatives {
		if isContinuationDirective(n) {
			t.Errorf("Expected isContinuationDirective(%q) = false, got true", n)
		}
	}
}

func TestSemanticForwardIntent(t *testing.T) {
	// The exact phrases from the real-world Telegram incident
	telegramIncidentPhrases := []string{
		"Aku mulai melihat gambaran. Proses utama server (node) memakai 32% memori — itu cukup tinggi. Sekarang aku lihat log server untuk lihat detil entri dekat dengan waktu pesan 'lanjutkan' (sekitar 14:20–14:32) dan ada nggak error/timeout di situ.",
		"Aku lanjutkan investigasi. Lihat log sistem/server untuk cari jejak yang terkait dengan penghentian alur saat pesan 'lanjutkan' diproses.",
		"Aku sudah menemukan akar masalahnya... Aku coba satu cara langsung: fetch URL Twitter-nya lewat read_url biar dapat output nyata dan bukan dead-end seperti sebelumnya.",
		"Ada temuan penting: URL yang kupakai (status ID palsu 1850000000000000000) ternyata 404 — itu bukan URL asli dari postingan yang kita baca tadi. Aku perlu ambil URL yang benar dari riwayat sesi biar tidak dead-end terus. Aku cari URL postingan asli tadi di riwayat sesi.",
		"Mari kita periksa file config.json untuk melihat variabel environment.",
		"Next I will inspect the docker container logs to find the root cause.",
		"Sekarang saya akan buat file hello.txt",
		"Kucari URL di riwayat sesi.",
	}

	for _, phrase := range telegramIncidentPhrases {
		if !looksLikeContinuation(phrase) {
			t.Errorf("Expected looksLikeContinuation = true for phrase:\n%s\nGot false", phrase)
		}
	}

	completedPhrases := []string{
		"Berikut hasil analisanya: Server berjalan normal, RAM terpakai 40%, load average 0.2. Semua proses sudah selesai diverifikasi.",
		"File /tmp/data.txt sudah berhasil dibuat dan terverifikasi.",
		"Semua tugas selesai. Tidak ada tindakan lebih lanjut yang diperlukan.",
		"All done. The container is running and healthy.",
	}

	for _, phrase := range completedPhrases {
		if looksLikeContinuation(phrase) {
			t.Errorf("Expected looksLikeContinuation = false for completed phrase:\n%s\nGot true", phrase)
		}
	}
}
