package main

import (
	"html"
	"os"
	"regexp"
	"strings"
)

// formatFinalResponse cleans up the markdown/HTML header and formats response
func formatFinalResponse(text string) string {
	clean := stripHTML(text)
	// Remove Scorp header if present to avoid redundant prefixes
	reHeader := regexp.MustCompile(`(?i)^(?:🤖\s*)?(?:Scorp:\s*)?`)
	clean = reHeader.ReplaceAllString(strings.TrimSpace(clean), "")
	return strings.TrimSpace(clean)
}

// formatTerminalText converts HTML formatting to clean readable terminal text
func formatTerminalText(s string) string {
	return strings.TrimSpace(stripHTML(s))
}

// stripHTML removes HTML tags, decodes entities, and cleans whitespace
func stripHTML(s string) string {
	// Replace <pre><code> and </code></pre> with code block backticks
	s = regexp.MustCompile(`(?i)<pre><code>([\s\S]*?)</code></pre>`).ReplaceAllString(s, "\n```\n$1\n```\n")
	s = regexp.MustCompile(`(?i)<pre>([\s\S]*?)</pre>`).ReplaceAllString(s, "\n```\n$1\n```\n")
	s = regexp.MustCompile(`(?i)<code>([\s\S]*?)</code>`).ReplaceAllString(s, "`$1`")

	var sb strings.Builder
	inTag := false
	for _, r := range s {
		switch {
		case r == '<':
			inTag = true
		case r == '>':
			inTag = false
		case !inTag:
			sb.WriteRune(r)
		}
	}
	return html.UnescapeString(sb.String())
}

// isTerminal checks if a file descriptor is an interactive terminal
func isTerminal(f *os.File) bool {
	stat, err := f.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) != 0
}
