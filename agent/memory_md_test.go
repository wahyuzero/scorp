package agent

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func memoryMDTestFile(t *testing.T) string {
	t.Helper()
	orig := memoryMDFile
	tmp := filepath.Join(t.TempDir(), "MEMORY.md")
	memoryMDFile = tmp
	t.Cleanup(func() { memoryMDFile = orig })
	return tmp
}

func TestAppendMemoryMD_DedupQuotaAndHeader(t *testing.T) {
	memoryMDTestFile(t)

	added := AppendMemoryMD([]string{"Deploy target is tencent-vps", "", "Deploy target is TENCENT-VPS", "Uses supervised autonomy"})
	if added != 2 {
		t.Fatalf("expected 2 unique entries added, got %d", added)
	}

	// exact-duplicate (case-insensitive) across calls adds nothing
	if n := AppendMemoryMD([]string{"deploy target is Tencent-VPS"}); n != 0 {
		t.Fatalf("duplicate entry must not be added, got %d", n)
	}

	// header lines always survive
	AppendMemoryMD([]string{"# custom header"})
	content := ReadMemoryMD()
	if !strings.Contains(content, "# custom header") && !strings.Contains(content, "# Durable Memory") {
		t.Fatalf("expected a header in MEMORY.md:\n%s", content)
	}
	if !strings.Contains(content, "Uses supervised autonomy") {
		t.Fatalf("entry missing from file:\n%s", content)
	}
}

func TestAppendMemoryMD_QuotaDropsOldest(t *testing.T) {
	memoryMDTestFile(t)

	entries := make([]string, 0, memoryMDMaxLines+10)
	for i := 0; i < memoryMDMaxLines+10; i++ {
		entries = append(entries, fmt.Sprintf("e entry number %d", i))
	}
	AppendMemoryMD(entries)

	content := ReadMemoryMD()
	lines := nonEmptyLines(content)
	// header + at most 200 bullets
	if len(lines) > memoryMDMaxLines+1 {
		t.Fatalf("quota exceeded: %d lines", len(lines))
	}
	// oldest entries must be gone — the first generated entry can't survive
	if strings.Contains(content, "entry number 0") && !strings.Contains(content, fmt.Sprintf("entry number %d", memoryMDMaxLines+9)) {
		t.Fatal("oldest entries should have been dropped in favor of newest")
	}
}

func TestReadMemoryMD_Bounded(t *testing.T) {
	path := memoryMDTestFile(t)
	big := strings.Repeat("x", memoryMDMaxInject+5000)
	if err := os.WriteFile(path, []byte(big), 0644); err != nil {
		t.Fatalf("setup write: %v", err)
	}
	out := ReadMemoryMD()
	if len(out) > memoryMDMaxInject+50 {
		t.Fatalf("ReadMemoryMD must bound output, got %d chars", len(out))
	}
}

func TestParseMemoryEntries(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want int
	}{
		{"plain array", `["a","b"]`, 2},
		{"fenced array", "```json\n[\"fact one\"]\n```", 1},
		{"objects with content", `[{"content":"decided X"},{"memory":"uses Y"}]`, 2},
		{"empty", `[]`, 0},
		{"garbage", `not json at all`, 0},
		{"over limit capped", `["1","2","3","4","5","6","7"]`, memoryExtractLimit},
	}
	for _, tc := range cases {
		got := parseMemoryEntries(tc.in)
		if len(got) != tc.want {
			t.Errorf("%s: parseMemoryEntries got %d entries (%v), want %d", tc.name, len(got), got, tc.want)
		}
	}
}

func nonEmptyLines(s string) []string {
	var out []string
	for _, l := range strings.Split(s, "\n") {
		if strings.TrimSpace(l) != "" {
			out = append(out, l)
		}
	}
	return out
}
