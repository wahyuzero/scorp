package main

import (
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
