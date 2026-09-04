package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"

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
				fmt.Print("\033[1B\033[2K")
			}
			fmt.Printf("\033[%dA", popupRenderedLines)
			popupRenderedLines = 0
		}
	}

	render := func() {
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
		if strings.HasPrefix(currentStr, "/") && !strings.Contains(currentStr, " ") {
			matches := filterCommands(currentStr)
			if len(matches) > 0 {
				if selectedIndex >= len(matches) {
					selectedIndex = 0
				}
				if selectedIndex < 0 {
					selectedIndex = len(matches) - 1
				}
				lines = renderPopupBox(matches, selectedIndex)
			}
		}

		// Add persistent bottom statusline below prompt (or below popup if open)
		statusline := renderStatusFooter(currentSessionID)
		lines = append(lines, statusline)
		popupRenderedLines = len(lines)

		// Render below cursor, then restore cursor
		fmt.Print("\033[s") // Save cursor position
		for _, line := range lines {
			fmt.Print("\r\n\033[K" + line)
		}
		fmt.Print("\033[u") // Restore cursor position

		// Move cursor to proper position in buffer
		moveBack := len(buf) - cursorPos
		if moveBack > 0 {
			fmt.Printf("\033[%dD", moveBack)
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

// renderPopupBox generates styled ANSI lines for the slash commands popup box
func renderPopupBox(commands []SlashCommand, selected int) []string {
	var lines []string

	header := "\033[1;34m┌─ Slash Commands ────────────────────────────────────────────────────────┐\033[0m"
	footer := "\033[1;34m└─────────────────────────────────────────────────────────────────────────┘\033[0m"
	tip := "\033[2m  (Use ↑/↓ to navigate, Tab to complete, Enter to select, Esc to cancel)\033[0m"

	lines = append(lines, header)

	maxVisible := 6
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

		leftCol := fmt.Sprintf("%-30s", cmdStr)
		if len(leftCol) > 30 {
			leftCol = leftCol[:29] + " "
		}

		desc := cmd.Description
		if len(desc) > 38 {
			desc = desc[:35] + "..."
		}
		rightCol := fmt.Sprintf("%-38s", desc)

		if i == selected {
			line := fmt.Sprintf("\033[1;34m│\033[0m \033[1;37;44m ▶ %-28s — %-38s \033[0m \033[1;34m│\033[0m", leftCol, rightCol)
			lines = append(lines, line)
		} else {
			line := fmt.Sprintf("\033[1;34m│\033[0m   \033[1;36m%-28s\033[0m \033[2m—\033[0m %-38s \033[1;34m│\033[0m", leftCol, rightCol)
			lines = append(lines, line)
		}
	}

	lines = append(lines, footer)
	lines = append(lines, tip)
	return lines
}
