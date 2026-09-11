package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
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
// Supports bracketed paste mode (\033[200~ ... \033[201~) and unbracketed paste bursts
// so pasting paragraphs with newlines is preserved as a single user turn instead of prematurely submitting each line.
func readInteractiveInput(prompt string) (string, error) {
	prompt = strings.TrimLeft(prompt, "\r\n")
	fd := int(os.Stdin.Fd())
	if !term.IsTerminal(fd) {
		if prompt != "" {
			fmt.Print(prompt)
		}
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			return scanner.Text(), nil
		}
		return "", scanner.Err()
	}

	oldState, err := term.MakeRaw(fd)
	if err != nil {
		if prompt != "" {
			fmt.Print(prompt)
		}
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			return scanner.Text(), nil
		}
		return "", scanner.Err()
	}

	var restoreOnce sync.Once
	cleanup := func() {
		restoreOnce.Do(func() {
			disableBracketedPaste()
			_ = term.Restore(fd, oldState)
		})
	}
	defer cleanup()

	enableBracketedPaste()

	getTermSize := func() (int, int) {
		w, h, err := term.GetSize(fd)
		if err != nil || w <= 0 || h <= 0 {
			return 80, 24
		}
		return w, h
	}

	return readInputEngine(os.Stdin, os.Stdout, prompt, getTermSize, cleanup)
}

// readInputEngine is the core input loop separated for direct testability with mocked readers/writers.
func readInputEngine(in io.Reader, out io.Writer, prompt string, getTermSize func() (int, int), onExit func()) (string, error) {
	prompt = strings.TrimLeft(prompt, "\r\n")

	var buf []rune
	cursorPos := 0
	selectedIndex := 0
	popupRenderedLines := 0
	renderedBufferLines := 1
	renderedCursorLine := 0
	isPasting := false

	clearPopup := func() {
		if popupRenderedLines > 0 {
			for i := 0; i < popupRenderedLines; i++ {
				fmt.Fprint(out, "\033[1B\033[2K")
			}
			fmt.Fprintf(out, "\033[%dA\r", popupRenderedLines)
			popupRenderedLines = 0
		}
	}

	render := func() {
		w, h := 80, 24
		if getTermSize != nil {
			w, h = getTermSize()
		}
		if w <= 0 {
			w = 80
		}
		if h <= 0 {
			h = 24
		}
		isCompact := w < 72 || h < 22

		bufStr := string(buf)
		isSlash := strings.HasPrefix(bufStr, "/") && !strings.Contains(bufStr, " ") && !strings.Contains(bufStr, "\n")

		// If popup was active and we're no longer in slash command mode, clear it
		if !isSlash && popupRenderedLines > 0 {
			clearPopup()
		}

		bufLines := strings.Split(bufStr, "\n")
		numBufferLines := len(bufLines)

		// Calculate cursor line and column
		cursorLine := 0
		lineStart := 0
		for i := 0; i < cursorPos && i < len(buf); i++ {
			if buf[i] == '\n' {
				cursorLine++
				lineStart = i + 1
			}
		}
		cursorCol := cursorPos - lineStart

		// FAST PATH: Single line, not slash command, previous was single line with no popup
		if numBufferLines == 1 && renderedBufferLines == 1 && !isSlash && popupRenderedLines == 0 {
			fmt.Fprint(out, "\r\033[K"+prompt+bufStr)
			if cursorPos < len(buf) {
				colsBack := visibleWidth(string(buf[cursorPos:]))
				if colsBack > 0 {
					fmt.Fprintf(out, "\033[%dD", colsBack)
				}
			}
			renderedBufferLines = 1
			renderedCursorLine = 0
			return
		}

		// If popup was active, clear it before redrawing buffer
		if popupRenderedLines > 0 {
			clearPopup()
		}

		// Move cursor back to line 0 if it was on a lower line
		if renderedCursorLine > 0 {
			fmt.Fprintf(out, "\033[%dA", renderedCursorLine)
		}
		fmt.Fprint(out, "\r")

		// Redraw buffer lines
		// Line 0:
		fmt.Fprint(out, "\033[K"+prompt+bufLines[0])
		// Lines 1..N-1:
		for i := 1; i < numBufferLines; i++ {
			fmt.Fprint(out, "\r\n\033[K\033[2m... ❯\033[0m "+bufLines[i])
		}
		// If buffer shrank, clear leftover lines below
		if renderedBufferLines > numBufferLines {
			for j := numBufferLines; j < renderedBufferLines; j++ {
				fmt.Fprint(out, "\r\n\033[K")
			}
			fmt.Fprintf(out, "\033[%dA", renderedBufferLines-numBufferLines)
		}

		// Now cursor is at end of line numBufferLines - 1
		// Move to cursorLine
		if cursorLine < numBufferLines-1 {
			fmt.Fprintf(out, "\033[%dA", (numBufferLines-1)-cursorLine)
		}
		fmt.Fprint(out, "\r")
		linePrompt := prompt
		if cursorLine > 0 {
			linePrompt = "\033[2m... ❯\033[0m "
		}
		curLineRunes := []rune(bufLines[cursorLine])
		if cursorCol > len(curLineRunes) {
			cursorCol = len(curLineRunes)
		}
		fmt.Fprint(out, linePrompt+string(curLineRunes[:cursorCol]))

		renderedBufferLines = numBufferLines
		renderedCursorLine = cursorLine

		// Autocomplete popup: ONLY when isSlash
		if isSlash {
			matches := filterCommands(bufStr)
			var popupLines []string
			if len(matches) > 0 {
				if selectedIndex >= len(matches) {
					selectedIndex = 0
				}
				if selectedIndex < 0 {
					selectedIndex = len(matches) - 1
				}
				if isCompact {
					popupLines = renderCompactSuggestions(matches, selectedIndex, w)
				} else {
					popupLines = renderPopupBox(matches, selectedIndex, w, h)
				}
			}

			// Clamp lines
			for i := range popupLines {
				popupLines[i] = clampLineWidth(popupLines[i], w-2)
			}
			maxLinesBelow := h - 4
			if maxLinesBelow < 1 {
				maxLinesBelow = 1
			}
			if len(popupLines) > maxLinesBelow {
				popupLines = popupLines[:maxLinesBelow]
			}
			popupRenderedLines = len(popupLines)

			// Render below prompt
			for _, line := range popupLines {
				fmt.Fprint(out, "\r\n\033[2K"+line)
			}
			if len(popupLines) > 0 {
				fmt.Fprintf(out, "\033[%dA\r", len(popupLines))
				fmt.Fprint(out, prompt+string(buf[:cursorPos]))
			}
		} else {
			popupRenderedLines = 0
		}
	}

	render()

	readBuf := make([]byte, 64*1024)
	for {
		n, err := in.Read(readBuf)
		if err != nil {
			if popupRenderedLines > 0 {
				clearPopup()
			}
			return "", err
		}
		if n == 0 {
			continue
		}

		rawChunk := readBuf[:n]
		isBurst := (n > 1 && bytes.ContainsAny(rawChunk, "\r\n") && !bytes.HasPrefix(rawChunk, []byte("\033")))

		slice := rawChunk
		for len(slice) > 0 {
			// Check for Bracketed Paste Mode sequences: \033[200~ (start) and \033[201~ (end)
			if bytes.HasPrefix(slice, []byte("\033[200~")) {
				isPasting = true
				slice = slice[6:]
				continue
			}
			if bytes.HasPrefix(slice, []byte("\033[201~")) {
				isPasting = false
				slice = slice[6:]
				continue
			}

			// While inside bracketed paste block or clipboard paste burst:
			if isPasting || isBurst {
				// Handle newlines: CRLF, LF, or CR
				if slice[0] == '\r' {
					if len(slice) > 1 && slice[1] == '\n' {
						slice = slice[2:]
					} else {
						slice = slice[1:]
					}
					buf = append(buf[:cursorPos], append([]rune{'\n'}, buf[cursorPos:]...)...)
					cursorPos++
					continue
				}
				if slice[0] == '\n' {
					slice = slice[1:]
					buf = append(buf[:cursorPos], append([]rune{'\n'}, buf[cursorPos:]...)...)
					cursorPos++
					continue
				}
				if slice[0] == '\t' {
					slice = slice[1:]
					buf = append(buf[:cursorPos], append([]rune{'\t'}, buf[cursorPos:]...)...)
					cursorPos++
					continue
				}
				if slice[0] >= 32 {
					r, sz := utf8.DecodeRune(slice)
					if r != utf8.RuneError && sz > 0 {
						buf = append(buf[:cursorPos], append([]rune{r}, buf[cursorPos:]...)...)
						cursorPos++
						slice = slice[sz:]
					} else {
						slice = slice[1:]
					}
					continue
				}
				// Ignore other control characters in paste
				slice = slice[1:]
				continue
			}

			// Shift+Enter & Alt+Enter escape sequences:
			// \033[13;2u (CSI u / Kitty keyboard protocol for Shift+Enter)
			// \033[13;5u (CSI u Ctrl+Enter)
			// \033[27;2;13~ (xterm format for Shift+Enter)
			// \033[27;5;13~ (xterm Ctrl+Enter)
			// \033\r (ESC then CR: Alt+Enter)
			// \033\n (ESC then LF: Alt+Enter)
			// \033OM (SS3 Enter)
			if bytes.HasPrefix(slice, []byte("\033[13;2u")) {
				buf = append(buf[:cursorPos], append([]rune{'\n'}, buf[cursorPos:]...)...)
				cursorPos++
				selectedIndex = 0
				slice = slice[len("\033[13;2u"):]
				render()
				continue
			}
			if bytes.HasPrefix(slice, []byte("\033[13;5u")) {
				buf = append(buf[:cursorPos], append([]rune{'\n'}, buf[cursorPos:]...)...)
				cursorPos++
				selectedIndex = 0
				slice = slice[len("\033[13;5u"):]
				render()
				continue
			}
			if bytes.HasPrefix(slice, []byte("\033[27;2;13~")) {
				buf = append(buf[:cursorPos], append([]rune{'\n'}, buf[cursorPos:]...)...)
				cursorPos++
				selectedIndex = 0
				slice = slice[len("\033[27;2;13~"):]
				render()
				continue
			}
			if bytes.HasPrefix(slice, []byte("\033[27;5;13~")) {
				buf = append(buf[:cursorPos], append([]rune{'\n'}, buf[cursorPos:]...)...)
				cursorPos++
				selectedIndex = 0
				slice = slice[len("\033[27;5;13~"):]
				render()
				continue
			}
			if bytes.HasPrefix(slice, []byte("\033\r")) {
				buf = append(buf[:cursorPos], append([]rune{'\n'}, buf[cursorPos:]...)...)
				cursorPos++
				selectedIndex = 0
				slice = slice[2:]
				render()
				continue
			}
			if bytes.HasPrefix(slice, []byte("\033\n")) {
				buf = append(buf[:cursorPos], append([]rune{'\n'}, buf[cursorPos:]...)...)
				cursorPos++
				selectedIndex = 0
				slice = slice[2:]
				render()
				continue
			}
			if bytes.HasPrefix(slice, []byte("\033OM")) {
				buf = append(buf[:cursorPos], append([]rune{'\n'}, buf[cursorPos:]...)...)
				cursorPos++
				selectedIndex = 0
				slice = slice[3:]
				render()
				continue
			}

			b := slice[0]

			// Enter key (CR / 13) outside paste mode: submit complete multiline buffer
			if b == 13 {
				currentStr := string(buf)
				if popupRenderedLines > 0 {
					clearPopup()
				}
				if onExit != nil {
					onExit()
				}
				if renderedCursorLine < renderedBufferLines-1 {
					fmt.Fprintf(out, "\033[%dB", (renderedBufferLines-1)-renderedCursorLine)
				}
				fmt.Fprint(out, "\r\n")

				// If popup is active and user typed a slash prefix that is not an exact match
				if strings.HasPrefix(currentStr, "/") && !strings.Contains(currentStr, " ") && !strings.Contains(currentStr, "\n") {
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

			// Ctrl+J (LF / 10): insert newline at cursorPos and re-render multiline buffer without submitting
			if b == 10 {
				buf = append(buf[:cursorPos], append([]rune{'\n'}, buf[cursorPos:]...)...)
				cursorPos++
				selectedIndex = 0
				slice = slice[1:]
				render()
				continue
			}

			// Ctrl+C
			if b == 3 {
				if popupRenderedLines > 0 {
					clearPopup()
				}
				if onExit != nil {
					onExit()
				}
				fmt.Fprint(out, "\r\n")
				return "", fmt.Errorf("interrupted")
			}

			// Ctrl+D
			if b == 4 {
				if len(buf) == 0 {
					if popupRenderedLines > 0 {
						clearPopup()
					}
					if onExit != nil {
						onExit()
					}
					fmt.Fprint(out, "\r\n")
					return "/exit", nil
				}
				slice = slice[1:]
				continue
			}

			// Ctrl+A (Home)
			if b == 1 {
				for cursorPos > 0 && buf[cursorPos-1] != '\n' {
					cursorPos--
				}
				slice = slice[1:]
				render()
				continue
			}

			// Ctrl+E (End)
			if b == 5 {
				for cursorPos < len(buf) && buf[cursorPos] != '\n' {
					cursorPos++
				}
				slice = slice[1:]
				render()
				continue
			}

			// Ctrl+U (kill line to start)
			if b == 21 {
				start := cursorPos
				for start > 0 && buf[start-1] != '\n' {
					start--
				}
				buf = append(buf[:start], buf[cursorPos:]...)
				cursorPos = start
				selectedIndex = 0
				slice = slice[1:]
				render()
				continue
			}

			// Ctrl+K (kill line to end)
			if b == 11 {
				end := cursorPos
				for end < len(buf) && buf[end] != '\n' {
					end++
				}
				buf = append(buf[:cursorPos], buf[end:]...)
				selectedIndex = 0
				slice = slice[1:]
				render()
				continue
			}

			// Backspace: 127 or 8
			if b == 127 || b == 8 {
				if cursorPos > 0 {
					buf = append(buf[:cursorPos-1], buf[cursorPos:]...)
					cursorPos--
					selectedIndex = 0
					render()
				}
				slice = slice[1:]
				continue
			}

			// Tab: 9
			if b == 9 {
				currentStr := string(buf)
				if strings.HasPrefix(currentStr, "/") && !strings.Contains(currentStr, " ") && !strings.Contains(currentStr, "\n") {
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
				slice = slice[1:]
				continue
			}

			// ANSI escape sequences: 27
			if b == 27 {
				if len(slice) >= 3 && (slice[1] == '[' || slice[1] == 'O') {
					code := slice[2]
					slice = slice[3:]
					switch code {
					case 'A': // UP
						currentStr := string(buf)
						if strings.HasPrefix(currentStr, "/") && !strings.Contains(currentStr, " ") && !strings.Contains(currentStr, "\n") {
							selectedIndex--
							render()
						} else if strings.Contains(currentStr, "\n") {
							cursorPos = moveCursorVertical(buf, cursorPos, -1)
							render()
						}
					case 'B': // DOWN
						currentStr := string(buf)
						if strings.HasPrefix(currentStr, "/") && !strings.Contains(currentStr, " ") && !strings.Contains(currentStr, "\n") {
							selectedIndex++
							render()
						} else if strings.Contains(currentStr, "\n") {
							cursorPos = moveCursorVertical(buf, cursorPos, 1)
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
					case 'H': // Home
						for cursorPos > 0 && buf[cursorPos-1] != '\n' {
							cursorPos--
						}
						render()
					case 'F': // End
						for cursorPos < len(buf) && buf[cursorPos] != '\n' {
							cursorPos++
						}
						render()
					}
				} else {
					clearPopup()
					render()
					slice = slice[1:]
				}
				continue
			}

			// Printable UTF-8 runes
			r, sz := utf8.DecodeRune(slice)
			if sz > 0 && r != utf8.RuneError && (r >= 32 || r == '\t') {
				buf = append(buf[:cursorPos], append([]rune{r}, buf[cursorPos:]...)...)
				cursorPos++
				selectedIndex = 0
				slice = slice[sz:]
				render()
			} else {
				slice = slice[1:]
			}
		}

		if isPasting || isBurst {
			render()
		}
	}
}

// moveCursorVertical moves the cursor up (dir=-1) or down (dir=1) within a multiline buffer
func moveCursorVertical(buf []rune, cursorPos int, dir int) int {
	if len(buf) == 0 {
		return 0
	}
	if cursorPos < 0 {
		cursorPos = 0
	}
	if cursorPos > len(buf) {
		cursorPos = len(buf)
	}

	lines := strings.Split(string(buf), "\n")
	if len(lines) <= 1 {
		return cursorPos
	}

	curLine := 0
	lineStart := 0
	for i := 0; i < cursorPos && i < len(buf); i++ {
		if buf[i] == '\n' {
			curLine++
			lineStart = i + 1
		}
	}
	curCol := cursorPos - lineStart

	targetLine := curLine + dir
	if targetLine < 0 || targetLine >= len(lines) {
		return cursorPos
	}

	targetStart := 0
	for l := 0; l < targetLine; l++ {
		targetStart += len([]rune(lines[l])) + 1
	}

	targetLineLen := len([]rune(lines[targetLine]))
	targetCol := curCol
	if targetCol > targetLineLen {
		targetCol = targetLineLen
	}

	return targetStart + targetCol
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
