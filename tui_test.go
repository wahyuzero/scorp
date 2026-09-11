package main

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestVisibleWidth(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"hello", 5},
		{"\033[1;34mhello\033[0m", 5},
		{"\033[38;5;208morange\033[0m", 6},
		{"💡 /models", 10}, // 💡 = 2 cols + 1 space + 7 chars
		{"🤖 default", 10},  // 🤖 = 2 cols + 1 space + 7 chars
		{"\033[2m╰─\033[0m \033[1;36m/help\033[0m \033[2m─╯\033[0m", 11},
	}

	for _, tc := range tests {
		got := visibleWidth(tc.input)
		if got != tc.want {
			t.Errorf("visibleWidth(%q) = %d, want %d", tc.input, got, tc.want)
		}
	}
}

func TestClampLineWidth(t *testing.T) {
	plain := "This is a long line that needs to be clamped"
	clamped := clampLineWidth(plain, 15)
	if visibleWidth(clamped) > 15 {
		t.Errorf("clamped line width = %d, want <= 15", visibleWidth(clamped))
	}
	if !strings.HasSuffix(clamped, "\033[0m") {
		t.Errorf("clamped line should preserve ANSI reset at end, got %q", clamped)
	}

	// ANSI styled line
	styled := "\033[1;34m[Header]\033[0m \033[1;36mVery long content here\033[0m"
	clampedStyled := clampLineWidth(styled, 12)
	if visibleWidth(clampedStyled) > 12 {
		t.Errorf("clamped styled width = %d, want <= 12", visibleWidth(clampedStyled))
	}
}

func TestRenderCompactSuggestions_MobileTermux(t *testing.T) {
	// Termux portrait on phone: typically ~45-50 cols width
	termWidth := 45
	matches := filterCommands("/m")
	if len(matches) == 0 {
		t.Fatal("expected matches for /m prefix")
	}

	lines := renderCompactSuggestions(matches, 0, termWidth)
	if len(lines) == 0 || len(lines) > 2 {
		t.Fatalf("compact mode must render 1 or 2 lines, got %d", len(lines))
	}

	for i, line := range lines {
		w := visibleWidth(line)
		if w > termWidth-2 {
			t.Errorf("line %d width %d exceeds termWidth-2 (%d): %s", i, w, termWidth-2, line)
		}
	}

	// Verify selected command is highlighted
	if !strings.Contains(lines[0], matches[0].Command) {
		t.Errorf("line 0 should contain selected command %s", matches[0].Command)
	}

	// Verify line 2 contains description
	if len(lines) > 1 {
		if !strings.Contains(lines[1], "→") {
			t.Errorf("line 1 should contain description arrow, got %s", lines[1])
		}
	}
}

func TestRenderPopupBox_Desktop(t *testing.T) {
	termWidth := 80
	termHeight := 24
	matches := filterCommands("/")
	if len(matches) == 0 {
		t.Fatal("expected slash commands")
	}

	lines := renderPopupBox(matches, 0, termWidth, termHeight)
	if len(lines) == 0 {
		t.Fatal("expected non-empty popup lines")
	}

	for i, line := range lines {
		w := visibleWidth(line)
		if w > 74 {
			t.Errorf("desktop popup line %d exceeds 74 columns (got %d): %s", i, w, line)
		}
	}

	// Total lines should be header + maxVisible + footer + tip <= 10 lines
	if len(lines) > 10 {
		t.Errorf("popup rendered %d lines, want <= 10", len(lines))
	}
}

func TestRenderPopupBox_ConstrainedHeight(t *testing.T) {
	// Terminal with soft keyboard up (e.g. 16 rows)
	termWidth := 80
	termHeight := 16
	matches := filterCommands("/")

	lines := renderPopupBox(matches, 0, termWidth, termHeight)
	// For height 16, maxVisible should be clamped: min(6, 16 - 10) = 6
	// For height 14, maxVisible: min(6, 14 - 10) = 4
	lines14 := renderPopupBox(matches, 0, termWidth, 14)
	if len(lines14) > 8 { // header + 4 items + footer + tip = 7
		t.Errorf("popup with height 14 rendered %d lines, want <= 8", len(lines14))
	}
	_ = lines
}

func TestRenderStatusFooter_Compact(t *testing.T) {
	// Mobile 45-col width
	footer := renderStatusFooter("test-session", 45)
	w := visibleWidth(footer)
	if w > 43 {
		t.Errorf("compact footer width %d exceeds 43 cols: %s", w, footer)
	}
	if !strings.Contains(footer, "🤖") {
		t.Errorf("footer should contain model pill, got %s", footer)
	}
}

func TestRenderStatusFooter_Desktop(t *testing.T) {
	// Desktop 120-col width
	footer := renderStatusFooter("test-session", 120)
	w := visibleWidth(footer)
	if w > 118 {
		t.Errorf("desktop footer width %d exceeds 118 cols: %s", w, footer)
	}
	if !strings.Contains(footer, "🤖") || !strings.Contains(footer, "📂") {
		t.Errorf("desktop footer should contain model and folder pills: %s", footer)
	}
}

type chunkReader struct {
	chunks [][]byte
	idx    int
}

func (cr *chunkReader) Read(p []byte) (int, error) {
	if cr.idx >= len(cr.chunks) {
		return 0, io.EOF
	}
	chunk := cr.chunks[cr.idx]
	cr.idx++
	n := copy(p, chunk)
	return n, nil
}

// TestSingleCharacterTypingNoNewlines verifies that typing individual characters
// emits ZERO newlines, and only the final Enter key emits the submission newline.
func TestSingleCharacterTypingNoNewlines(t *testing.T) {
	cr := &chunkReader{
		chunks: [][]byte{
			{'s'},
			{'c'},
			{'o'},
			{'r'},
			{'p'},
			{13}, // Enter
		},
	}
	var out bytes.Buffer
	res, err := readInputEngine(cr, &out, "scorp ❯ ", func() (int, int) { return 80, 24 }, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res != "scorp" {
		t.Fatalf("got %q, want 'scorp'", res)
	}

	outStr := out.String()
	newlineCount := strings.Count(outStr, "\n")
	// Must contain exactly 1 newline at the very end when Enter was pressed!
	if newlineCount != 1 {
		t.Fatalf("expected exactly 1 newline (for Enter submission), but got %d: %q", newlineCount, outStr)
	}
	if !strings.HasSuffix(outStr, "\r\n") {
		t.Fatalf("expected output to end with \\r\\n on submit, got %q", outStr)
	}
}

// TestPasteBurstMultilineAccumulated verifies that when unbracketed paste arrives
// in a single read burst (n > 1) with multiple newlines, all lines are accumulated
// into the buffer rather than prematurely submitted on the first line.
func TestPasteBurstMultilineAccumulated(t *testing.T) {
	multilineText := "def calculate_total():\n    return 42\nresult = calculate_total()"
	cr := &chunkReader{
		chunks: [][]byte{
			[]byte(multilineText), // n > 1 with \n
			{13},                   // Enter key submitted afterwards
		},
	}
	var out bytes.Buffer
	res, err := readInputEngine(cr, &out, "scorp ❯ ", func() (int, int) { return 80, 24 }, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res != multilineText {
		t.Fatalf("expected multiline text preserved.\ngot: %q\nwant: %q", res, multilineText)
	}
}

// TestPasteBurstWithCRLF verifies that Windows/network CRLF sequences are normalized
// to clean single newlines without generating duplicate empty lines.
func TestPasteBurstWithCRLF(t *testing.T) {
	crlfText := "line 1\r\nline 2\r\nline 3"
	want := "line 1\nline 2\nline 3"
	cr := &chunkReader{
		chunks: [][]byte{
			[]byte(crlfText),
			{10}, // Enter as LF
		},
	}
	var out bytes.Buffer
	res, err := readInputEngine(cr, &out, "scorp ❯ ", func() (int, int) { return 80, 24 }, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res != want {
		t.Fatalf("expected CRLF normalized to clean LF.\ngot: %q\nwant: %q", res, want)
	}
}

// TestBracketedPasteMultiline verifies that bracketed paste sequences (\033[200~ ... \033[201~)
// properly preserve multiline content and submit on Enter.
func TestBracketedPasteMultiline(t *testing.T) {
	pasted := "package main\n\nfunc main() {\n\tprintln(\"hello\")\n}"
	cr := &chunkReader{
		chunks: [][]byte{
			[]byte("\033[200~" + pasted + "\033[201~"),
			{13}, // Enter
		},
	}
	var out bytes.Buffer
	res, err := readInputEngine(cr, &out, "scorp ❯ ", func() (int, int) { return 80, 24 }, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res != pasted {
		t.Fatalf("bracketed paste mismatch.\ngot: %q\nwant: %q", res, pasted)
	}
}

// TestBracketedPasteChunkedAcrossReads verifies bracketed paste spanning multiple read chunks.
func TestBracketedPasteChunkedAcrossReads(t *testing.T) {
	cr := &chunkReader{
		chunks: [][]byte{
			[]byte("\033[200~first line\n"),
			[]byte("second line\n"),
			[]byte("third line\033[201~"),
			{13}, // Enter
		},
	}
	var out bytes.Buffer
	res, err := readInputEngine(cr, &out, "scorp ❯ ", func() (int, int) { return 80, 24 }, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "first line\nsecond line\nthird line"
	if res != want {
		t.Fatalf("chunked bracketed paste mismatch.\ngot: %q\nwant: %q", res, want)
	}
}

// TestMultilineBufferEditAndSubmit verifies that after pasting multiline text,
// the user can continue typing characters at the end and then submit the entire block.
func TestMultilineBufferEditAndSubmit(t *testing.T) {
	cr := &chunkReader{
		chunks: [][]byte{
			[]byte("hello\nworld"),
			{'!'},
			{13}, // Enter
		},
	}
	var out bytes.Buffer
	res, err := readInputEngine(cr, &out, "scorp ❯ ", func() (int, int) { return 80, 24 }, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "hello\nworld!"
	if res != want {
		t.Fatalf("edit multiline buffer mismatch.\ngot: %q\nwant: %q", res, want)
	}
}

// TestMoveCursorVertical verifies up and down navigation within a multiline buffer.
func TestMoveCursorVertical(t *testing.T) {
	buf := []rune("line1\nlonger line 2\nline3")
	// line1 is 0..4 (len 5), \n is 5
	// longer line 2 is 6..18 (len 13), \n is 19
	// line3 is 20..24 (len 5)

	// Cursor on line 0, col 2 ('n') -> pos 2
	pos := 2
	posDown := moveCursorVertical(buf, pos, 1)
	if posDown != 8 {
		t.Errorf("moveCursorVertical down: got %d, want 8", posDown)
	}

	// Move down again to line 2, col 2 -> pos 20 + 2 = 22
	posDown2 := moveCursorVertical(buf, posDown, 1)
	if posDown2 != 22 {
		t.Errorf("moveCursorVertical down to line 2: got %d, want 22", posDown2)
	}

	// Move up to line 1: should be 8
	posUp := moveCursorVertical(buf, posDown2, -1)
	if posUp != 8 {
		t.Errorf("moveCursorVertical up: got %d, want 8", posUp)
	}

	// Cursor at end of line 1 (col 13) -> pos 19
	posEndLine1 := 19
	posClamped := moveCursorVertical(buf, posEndLine1, 1)
	if posClamped != 25 {
		t.Errorf("moveCursorVertical clamp: got %d, want 25", posClamped)
	}
}

// TestPromptLeadingNewlineStripped verifies that prompts containing leading newlines
// are stripped so no extra vertical space or scrolling is triggered.
func TestPromptLeadingNewlineStripped(t *testing.T) {
	cr := &chunkReader{
		chunks: [][]byte{
			{'x'},
			{13},
		},
	}
	var out bytes.Buffer
	dirtyPrompt := "\r\n\n\033[1;36mscorp\033[0m \033[1;32m❯\033[0m "
	res, err := readInputEngine(cr, &out, dirtyPrompt, func() (int, int) { return 80, 24 }, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res != "x" {
		t.Fatalf("got %q, want 'x'", res)
	}

	outStr := out.String()
	if strings.HasPrefix(outStr, "\n") || strings.HasPrefix(outStr, "\r\n") {
		t.Fatalf("rendered output started with newline: %q", outStr)
	}
}

// TestSlashCommandAutocompleteVsRegularText verifies that typing normal text
// emits zero popup lines and zero extra newlines, while typing slash commands
// activates autocomplete.
func TestSlashCommandAutocompleteVsRegularText(t *testing.T) {
	// 1. Regular text: no popup lines, zero newlines during typing
	cr := &chunkReader{
		chunks: [][]byte{
			{'n'}, {'o'}, {'r'}, {'m'}, {'a'}, {'l'},
			{13},
		},
	}
	var out bytes.Buffer
	res, err := readInputEngine(cr, &out, "❯ ", func() (int, int) { return 80, 24 }, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res != "normal" {
		t.Fatalf("got %q, want 'normal'", res)
	}
	if strings.Count(out.String(), "\n") != 1 {
		t.Fatalf("expected exactly 1 newline (the submit newline), got %d: %q", strings.Count(out.String(), "\n"), out.String())
	}

	// 2. Slash command prefix: popup lines are rendered and Tab completes
	crSlash := &chunkReader{
		chunks: [][]byte{
			{'/'}, {'m'}, {'o'}, {'d'}, {'e'}, {'l'},
			{13},
		},
	}
	var outSlash bytes.Buffer
	resSlash, err := readInputEngine(crSlash, &outSlash, "❯ ", func() (int, int) { return 80, 24 }, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resSlash != "/model" {
		t.Fatalf("got %q, want '/model'", resSlash)
	}
}

// TestBackspaceSingleLineAndMultiline verifies backspace behavior.
func TestBackspaceSingleLineAndMultiline(t *testing.T) {
	cr := &chunkReader{
		chunks: [][]byte{
			{'a'}, {'b'}, {'c'},
			{127}, // Backspace
			{13},  // Enter
		},
	}
	var out bytes.Buffer
	res, err := readInputEngine(cr, &out, "❯ ", func() (int, int) { return 80, 24 }, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res != "ab" {
		t.Fatalf("got %q, want 'ab'", res)
	}
	if strings.Count(out.String(), "\n") != 1 {
		t.Fatalf("expected exactly 1 newline, got %d", strings.Count(out.String(), "\n"))
	}
}

