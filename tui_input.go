package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"golang.org/x/term"
)

// SlashCommand represents a registered CLI command with description and autocomplete support
type SlashCommand struct {
	Command     string
	Args        string
	Description string
}

var availableSlashCommands = []SlashCommand{
	{Command: "/help", Args: "", Description: "Show commands and interactive usage help"},
	{Command: "/models", Args: "", Description: "List all configured models and key status"},
	{Command: "/model", Args: "<name>", Description: "Switch or inspect active AI model"},
	{Command: "/mode", Args: "<level>", Description: "Set autonomy level: readonly, supervised, yolo"},
	{Command: "/session", Args: "[list|new|use|rename|delete]", Description: "Manage conversational chat sessions"},
	{Command: "/status", Args: "", Description: "Display agent status, workspace, and model"},
	{Command: "/cost", Args: "", Description: "Check token usage and daily cost dashboard"},
	{Command: "/tools", Args: "", Description: "List all registered agent tools and schemas"},
	{Command: "/sop", Args: "[list|run <name>]", Description: "Run Standard Operating Procedures"},
	{Command: "/cron", Args: "[list|run|del]", Description: "View and manage scheduled background cron tasks"},
	{Command: "/queue", Args: "", Description: "View real-time agent steering & message queue"},
	{Command: "/skills", Args: "[list|<name>]", Description: "Inspect and activate specialized agent skills"},
	{Command: "/mcp", Args: "[list|restart <server>]", Description: "Monitor MCP server health and crash recovery"},
	{Command: "/receipts", Args: "", Description: "View cryptographic tool execution receipts"},
	{Command: "/compact", Args: "", Description: "Manually compact and summarize current session history"},
	{Command: "/clear", Args: "", Description: "Clear current conversation session history"},
	{Command: "/stop", Args: "", Description: "Interrupt or reset running agent mode"},
	{Command: "/exit", Args: "", Description: "Quit interactive Scorp session"},
}

// readInteractiveInput reads a line or multiline block from terminal with live autocomplete popup.
// Supports bracketed paste mode (\033[200~ ... \033[201~) so pasting paragraphs with newlines
// is preserved as a single user turn instead of prematurely submitting each line.
func readInteractiveInput(prompt string) (string, error) {
	fd := int(os.Stdin.Fd())
	if !term.IsTerminal(fd) {
		fmt.Print(prompt)
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			return scanner.Text(), nil
		}
		return "", scanner.Err()
	}

	oldState, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Print(prompt)
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			return scanner.Text(), nil
		}
		return "", scanner.Err()
	}
	defer func() {
		disableBracketedPaste()
		_ = term.Restore(fd, oldState)
	}()

	enableBracketedPaste()

	var buf []rune
	cursorPos := 0
	selectedIndex := 0
	popupRenderedLines := 0
	isPasting := false

	clearPopup := func() {
		if popupRenderedLines > 0 {
			for i := 0; i < popupRenderedLines; i++ {
				fmt.Print("\r\n\033[2K")
			}
			fmt.Printf("\033[%dA\r", popupRenderedLines)
			popupRenderedLines = 0
		}
	}

	render := func() {
		w, h, err := term.GetSize(fd)
		if err != nil || w <= 0 || h <= 0 {
			w = 80
			h = 24
		}
		isCompact := w < 72 || h < 22

		clearPopup()
		fmt.Print("\r\033[K")
		fmt.Print(prompt)

		// If buffer contains newlines (e.g. from paste), print with visual continuation
		bufStr := string(buf)
		if strings.Contains(bufStr, "\n") {
			parts := strings.Split(bufStr, "\n")
			for idx, p := range parts {
				if idx > 0 {
					fmt.Print("\r\n\033[K\033[2m... ❯\033[0m ")
				}
				fmt.Print(p)
			}
		} else {
			fmt.Print(bufStr)
		}

		currentStr := string(buf)
		var lines []string
		isSlash := strings.HasPrefix(currentStr, "/") && !strings.Contains(currentStr, " ")
		if isSlash {
			matches := filterCommands(currentStr)
			if len(matches) > 0 {
				if selectedIndex >= len(matches) {
					selectedIndex = 0
				}
				if selectedIndex < 0 {
					selectedIndex = len(matches) - 1
				}

				if isCompact {
					// Compact 1-2 line inline suggestion (mobile Termux / small screen)
					lines = renderCompactSuggestions(matches, selectedIndex, w)
				} else {
					// Desktop dynamic popup box + status footer
					lines = renderPopupBox(matches, selectedIndex, w, h)
					statusline := renderStatusFooter(currentSessionID, w)
					if statusline != "" {
						lines = append(lines, statusline)
					}
				}
			} else {
				statusline := renderStatusFooter(currentSessionID, w)
				if statusline != "" {
					lines = append(lines, statusline)
				}
			}
		} else {
			statusline := renderStatusFooter(currentSessionID, w)
			if statusline != "" {
				lines = append(lines, statusline)
			}
		}

		// Clamp each line's visual width to w-2 to strictly avoid terminal line-wrapping
		for i := range lines {
			lines[i] = clampLineWidth(lines[i], w-2)
		}

		// Cap maximum lines rendered below prompt to prevent vertical scroll
		maxLinesBelow := h - 4
		if maxLinesBelow < 1 {
			maxLinesBelow = 1
		}
		if len(lines) > maxLinesBelow {
			lines = lines[:maxLinesBelow]
		}
		popupRenderedLines = len(lines)

		// Render below prompt, then restore cursor relatively (no ANSI \033[s / \033[u drift)
		for _, line := range lines {
			fmt.Print("\r\n\033[2K" + line)
		}
		if len(lines) > 0 {
			fmt.Printf("\033[%dA", len(lines))
		}

		// Return to column 0 on the prompt line, then advance cursor to cursorPos
		fmt.Print("\r")
		if strings.Contains(bufStr, "\n") {
			parts := strings.Split(bufStr, "\n")
			fmt.Print(prompt + parts[len(parts)-1])
		} else {
			fmt.Print(prompt + string(buf[:cursorPos]))
		}
	}

	render()

	readBuf := make([]byte, 512)
	for {
		n, err := os.Stdin.Read(readBuf)
		if err != nil {
			clearPopup()
			return "", err
		}
		if n == 0 {
			continue
		}

		// Check for Bracketed Paste Mode sequences: \033[200~ (start) and \033[201~ (end)
		slice := readBuf[:n]
		for len(slice) > 0 {
			if bytes.HasPrefix(slice, []byte("\033[200~")) {
				isPasting = true
				slice = slice[6:]
				continue
			}
			if bytes.HasPrefix(slice, []byte("\033[201~")) {
				isPasting = false
				slice = slice[6:]
				render()
				continue
			}

			b := slice[0]
			slice = slice[1:]

			// While inside bracketed paste block, treat newlines as literal \n without submitting
			if isPasting {
				if b == 13 || b == 10 {
					buf = append(buf[:cursorPos], append([]rune{'\n'}, buf[cursorPos:]...)...)
					cursorPos++
					continue
				}
				if b >= 32 || b == '\t' {
					r := rune(b)
					buf = append(buf[:cursorPos], append([]rune{r}, buf[cursorPos:]...)...)
					cursorPos++
					continue
				}
			}

			// Enter key outside paste mode: submit
			if b == 13 || b == 10 {
				currentStr := string(buf)
				clearPopup()
				disableBracketedPaste()
				_ = term.Restore(fd, oldState)
				fmt.Print("\r\n")

				// If popup is active and user typed a prefix that is not an exact match
				if strings.HasPrefix(currentStr, "/") && !strings.Contains(currentStr, " ") {
					matches := filterCommands(currentStr)
					if len(matches) > 0 {
						isExact := false
						for _, m := range matches {
							if m.Command == currentStr {
								isExact = true
								break
							}
						}
						if !isExact {
							matched := matches[selectedIndex]
							return matched.Command, nil
						}
					}
				}
				return string(buf), nil
			}

			// Ctrl+C
			if b == 3 {
				clearPopup()
				disableBracketedPaste()
				_ = term.Restore(fd, oldState)
				fmt.Print("\r\n")
				return "", fmt.Errorf("interrupted")
			}

			// Ctrl+D
			if b == 4 {
				clearPopup()
				disableBracketedPaste()
				_ = term.Restore(fd, oldState)
				fmt.Print("\r\n")
				return "/exit", nil
			}

			// Backspace: 127 or 8
			if b == 127 || b == 8 {
				if cursorPos > 0 {
					buf = append(buf[:cursorPos-1], buf[cursorPos:]...)
					cursorPos--
					selectedIndex = 0
					render()
				}
				continue
			}

			// Tab: 9
			if b == 9 {
				currentStr := string(buf)
				if strings.HasPrefix(currentStr, "/") && !strings.Contains(currentStr, " ") {
					matches := filterCommands(currentStr)
					if len(matches) > 0 {
						cmd := matches[selectedIndex].Command
						if matches[selectedIndex].Args != "" {
							cmd += " "
						}
						buf = []rune(cmd)
						cursorPos = len(buf)
						selectedIndex = 0
						render()
					}
				}
				continue
			}

			// ANSI escape sequences: 27, 91 (or 79)
			if b == 27 {
				if len(slice) >= 2 && (slice[0] == '[' || slice[0] == 'O') {
					code := slice[1]
					slice = slice[2:]
					switch code {
					case 'A': // UP
						currentStr := string(buf)
						if strings.HasPrefix(currentStr, "/") && !strings.Contains(currentStr, " ") {
							selectedIndex--
							render()
						}
					case 'B': // DOWN
						currentStr := string(buf)
						if strings.HasPrefix(currentStr, "/") && !strings.Contains(currentStr, " ") {
							selectedIndex++
							render()
						}
					case 'C': // RIGHT
						if cursorPos < len(buf) {
							cursorPos++
							render()
						}
					case 'D': // LEFT
						if cursorPos > 0 {
							cursorPos--
							render()
						}
					}
				} else {
					clearPopup()
					render()
				}
				continue
			}

			// Printable ASCII characters
			if b >= 32 && b <= 126 {
				r := rune(b)
				buf = append(buf[:cursorPos], append([]rune{r}, buf[cursorPos:]...)...)
				cursorPos++
				selectedIndex = 0
				render()
			}
		}
	}
}

// filterCommands filters available slash commands based on current prefix
func filterCommands(prefix string) []SlashCommand {
	lowerPrefix := strings.ToLower(prefix)
	var matches []SlashCommand
	for _, sc := range availableSlashCommands {
		if strings.HasPrefix(strings.ToLower(sc.Command), lowerPrefix) {
			matches = append(matches, sc)
		}
	}
	return matches
}

// visibleWidth calculates string visual display width in terminal columns (ignoring ANSI escape sequences)
func visibleWidth(s string) int {
	w := 0
	inEscape := false
	for i := 0; i < len(s); {
		if s[i] == '\033' {
			inEscape = true
			i++
			continue
		}
		if inEscape {
			if (s[i] >= 'a' && s[i] <= 'z') || (s[i] >= 'A' && s[i] <= 'Z') || s[i] == '~' {
				inEscape = false
			}
			i++
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		i += size
		if r == 0xfe0f || r == 0xfe0e {
			continue
		}
		if isWideRune(r) {
			w += 2
		} else {
			w += 1
		}
	}
	return w
}

func isWideRune(r rune) bool {
	return (r >= 0x1100 && r <= 0x115f) ||
		(r >= 0x2e80 && r <= 0xa4cf) ||
		(r >= 0xac00 && r <= 0xd7a3) ||
		(r >= 0xf900 && r <= 0xfaff) ||
		(r >= 0xfe10 && r <= 0xfe19) ||
		(r >= 0xfe30 && r <= 0xfe6f) ||
		(r >= 0xff00 && r <= 0xff60) ||
		(r >= 0xffe0 && r <= 0xffe6) ||
		(r >= 0x1f300 && r <= 0x1faff) ||
		r == 0x26a1 || r == 0x2705 || r == 0x274c
}

// clampLineWidth truncates string s so its visible width does not exceed maxWidth, preserving ANSI reset
func clampLineWidth(s string, maxWidth int) string {
	if maxWidth <= 0 {
		return s
	}
	if visibleWidth(s) <= maxWidth {
		return s
	}

	var sb strings.Builder
	curW := 0
	inEscape := false
	var escBuf strings.Builder

	for i := 0; i < len(s); {
		if s[i] == '\033' {
			inEscape = true
			escBuf.Reset()
			escBuf.WriteByte(s[i])
			i++
			continue
		}
		if inEscape {
			escBuf.WriteByte(s[i])
			if (s[i] >= 'a' && s[i] <= 'z') || (s[i] >= 'A' && s[i] <= 'Z') || s[i] == '~' {
				inEscape = false
				sb.WriteString(escBuf.String())
			}
			i++
			continue
		}

		r, size := utf8.DecodeRuneInString(s[i:])
		i += size
		if r == 0xfe0f || r == 0xfe0e {
			sb.WriteRune(r)
			continue
		}

		rw := 1
		if isWideRune(r) {
			rw = 2
		}

		if curW+rw > maxWidth {
			break
		}
		sb.WriteRune(r)
		curW += rw
	}

	sb.WriteString("\033[0m")
	return sb.String()
}

// renderCompactSuggestions renders a clean, compact 1-2 line inline suggestion for mobile/narrow terminals
func renderCompactSuggestions(commands []SlashCommand, selected int, termWidth int) []string {
	if len(commands) == 0 {
		return nil
	}
	if termWidth <= 0 {
		termWidth = 80
	}
	maxW := termWidth - 2

	// Line 1: 💡 [/models]  /model  /mode  (Tab)
	prefix := "💡 "
	var pillStrs []string
	for i, cmd := range commands {
		var p string
		if i == selected {
			p = fmt.Sprintf("\033[1;37;44m %s \033[0m", cmd.Command)
		} else {
			p = fmt.Sprintf("\033[1;36m%s\033[0m", cmd.Command)
		}
		pillStrs = append(pillStrs, p)
	}

	start := 0
	if selected > 2 {
		start = selected - 1
	}

	line1 := prefix
	tip := " \033[2m(Tab)\033[0m"
	tipW := visibleWidth(tip)

	added := 0
	for i := start; i < len(pillStrs); i++ {
		sep := " "
		if added == 0 {
			sep = ""
		}
		candidate := line1 + sep + pillStrs[i]
		if visibleWidth(candidate)+tipW > maxW {
			if added == 0 {
				line1 = candidate
			} else {
				line1 += " \033[2m…\033[0m"
			}
			break
		}
		line1 = candidate
		added++
	}
	if visibleWidth(line1)+tipW <= maxW {
		line1 += tip
	}

	var lines []string
	lines = append(lines, line1)

	// Line 2: Selected command description clamped to maxW
	if selected >= 0 && selected < len(commands) {
		cmd := commands[selected]
		args := ""
		if cmd.Args != "" {
			args = " " + cmd.Args
		}
		line2 := fmt.Sprintf("   \033[2m→ %s%s: %s\033[0m", cmd.Command, args, cmd.Description)
		lines = append(lines, clampLineWidth(line2, maxW))
	}

	return lines
}

// renderPopupBox generates styled ANSI lines for the slash commands popup box with dynamic width and height clamping
func renderPopupBox(commands []SlashCommand, selected int, termWidth, termHeight int) []string {
	var lines []string

	if termWidth <= 0 {
		termWidth = 80
	}
	if termHeight <= 0 {
		termHeight = 24
	}

	boxWidth := termWidth - 4
	if boxWidth > 74 {
		boxWidth = 74
	}
	if boxWidth < 40 {
		boxWidth = 40
	}

	maxVisible := 6
	if maxVisible > termHeight-10 {
		maxVisible = termHeight - 10
	}
	if maxVisible < 2 {
		maxVisible = 2
	}

	leftColWidth := 25
	if boxWidth < 74 {
		leftColWidth = boxWidth/2 - 5
		if leftColWidth < 12 {
			leftColWidth = 12
		}
	}
	rightColWidth := boxWidth - leftColWidth - 9
	if rightColWidth < 10 {
		rightColWidth = 10
	}

	dashCount := boxWidth - 19
	if dashCount < 1 {
		dashCount = 1
	}
	header := fmt.Sprintf("\033[1;34m┌─ Slash Commands %s┐\033[0m", strings.Repeat("─", dashCount))
	footer := fmt.Sprintf("\033[1;34m└%s┘\033[0m", strings.Repeat("─", boxWidth-2))
	tip := "\033[2m  (Use ↑/↓ to navigate, Tab to complete, Enter to select, Esc to cancel)\033[0m"
	if visibleWidth(tip) > boxWidth {
		tip = "\033[2m  (Tab: complete, ↑/↓: navigate, Enter: select)\033[0m"
	}
	tip = clampLineWidth(tip, boxWidth)

	lines = append(lines, header)

	start := 0
	if selected >= maxVisible {
		start = selected - maxVisible + 1
	}
	end := start + maxVisible
	if end > len(commands) {
		end = len(commands)
	}

	for i := start; i < end; i++ {
		cmd := commands[i]
		cmdStr := cmd.Command
		if cmd.Args != "" {
			cmdStr += " " + cmd.Args
		}

		if len(cmdStr) > leftColWidth {
			cmdStr = cmdStr[:leftColWidth-1] + " "
		}
		leftCol := fmt.Sprintf("%-*s", leftColWidth, cmdStr)

		desc := cmd.Description
		if len(desc) > rightColWidth {
			desc = desc[:rightColWidth-3] + "..."
		}
		rightCol := fmt.Sprintf("%-*s", rightColWidth, desc)

		var line string
		if i == selected {
			line = fmt.Sprintf("\033[1;34m│\033[0m \033[1;37;44m▶ %s — %s\033[0m \033[1;34m│\033[0m", leftCol, rightCol)
		} else {
			line = fmt.Sprintf("\033[1;34m│\033[0m   \033[1;36m%s\033[0m \033[2m—\033[0m %s \033[1;34m│\033[0m", leftCol, rightCol)
		}
		lines = append(lines, line)
	}

	lines = append(lines, footer)
	lines = append(lines, tip)
	return lines
}
