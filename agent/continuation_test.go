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
		if !IsContinuationDirective(p) {
			t.Errorf("Expected IsContinuationDirective(%q) = true, got false", p)
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
		if IsContinuationDirective(n) {
			t.Errorf("Expected IsContinuationDirective(%q) = false, got true", n)
		}
	}
}

func TestPureInformationalQuery(t *testing.T) {
	info := []string{
		"Jelaskan apa itu Docker container",
		"Explain what is Kubernetes",
		"Apa perbedaan PostgreSQL dan MySQL?",
		"Kenapa Go sangat cepat?",
	}
	for _, q := range info {
		if !IsPureInformationalQuery(q) {
			t.Errorf("Expected IsPureInformationalQuery(%q) = true, got false", q)
		}
	}

	action := []string{
		"Buat file /tmp/test.py",
		"Jalankan command df -h",
		"Hapus direktori /tmp/test",
		"Cek status git dan edit code",
	}
	for _, a := range action {
		if IsPureInformationalQuery(a) {
			t.Errorf("Expected IsPureInformationalQuery(%q) = false, got true", a)
		}
	}
}
