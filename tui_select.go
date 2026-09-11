package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"

	"golang.org/x/term"
	"scorp-agent/agent"
	"scorp-agent/models"
)

// SelectChoice represents an option in the interactive selection menu.
type SelectChoice struct {
	ID     string
	Render func(isSelected bool) string
}

// runInteractiveSelectEngine is the pure I/O engine for the interactive selector, testable with mock readers/writers.
func runInteractiveSelectEngine(
	in io.Reader,
	out io.Writer,
	title string,
	items []SelectChoice,
	initialIndex int,
	allowNew bool,
	getTermSize func() (int, int),
) (action string, selectedIndex int, err error) {
	if len(items) == 0 {
		return "cancel", 0, nil
	}
	if initialIndex < 0 {
		initialIndex = 0
	}
	if initialIndex >= len(items) {
		initialIndex = len(items) - 1
	}

	selectedIndex = initialIndex

	termHeight := 24
	if getTermSize != nil {
		_, h := getTermSize()
		if h > 0 {
			termHeight = h
		}
	}
	maxVisible := termHeight - 4
	if maxVisible > 15 {
		maxVisible = 15
	}
	if maxVisible < 3 {
		maxVisible = 3
	}

	visibleCount := len(items)
	if visibleCount > maxVisible {
		visibleCount = maxVisible
	}

	windowStart := 0
	adjustWindow := func() {
		if selectedIndex < windowStart {
			windowStart = selectedIndex
		}
		if selectedIndex >= windowStart+maxVisible {
			windowStart = selectedIndex - maxVisible + 1
		}
		if windowStart < 0 {
			windowStart = 0
		}
		if windowStart > len(items)-maxVisible && len(items) > maxVisible {
			windowStart = len(items) - maxVisible
		}
	}
	adjustWindow()

	// Initial render
	fmt.Fprintf(out, "\r\n\033[1;33m%s\033[0m\r\n", title)
	for i := 0; i < visibleCount; i++ {
		actualIdx := windowStart + i
		fmt.Fprint(out, "\033[2K"+items[actualIdx].Render(actualIdx == selectedIndex)+"\r\n")
	}

	clearMenu := func() {
		totalLines := visibleCount + 2
		fmt.Fprintf(out, "\033[%dA\r", totalLines)
		for i := 0; i < totalLines; i++ {
			fmt.Fprint(out, "\033[2K\r\n")
		}
		fmt.Fprintf(out, "\033[%dA\r", totalLines)
	}

	redraw := func() {
		adjustWindow()
		fmt.Fprintf(out, "\033[%dA\r", visibleCount)
		for i := 0; i < visibleCount; i++ {
			actualIdx := windowStart + i
			fmt.Fprint(out, "\033[2K"+items[actualIdx].Render(actualIdx == selectedIndex)+"\r\n")
		}
	}

	readBuf := make([]byte, 128)
	for {
		n, err := in.Read(readBuf)
		if err != nil {
			clearMenu()
			return "cancel", selectedIndex, err
		}
		if n == 0 {
			continue
		}

		slice := readBuf[:n]
		for len(slice) > 0 {
			// Ctrl+C (3)
			if slice[0] == 3 {
				clearMenu()
				return "cancel", selectedIndex, nil
			}
			// Enter (13 or 10)
			if slice[0] == 13 || slice[0] == 10 {
				clearMenu()
				return "select", selectedIndex, nil
			}
			// Esc / Escape sequences (27)
			if slice[0] == 27 {
				if len(slice) >= 3 && (slice[1] == '[' || slice[1] == 'O') {
					code := slice[2]
					slice = slice[3:]
					if code == 'A' { // UP
						if selectedIndex > 0 {
							selectedIndex--
							redraw()
						}
						continue
					} else if code == 'B' { // DOWN
						if selectedIndex < len(items)-1 {
							selectedIndex++
							redraw()
						}
						continue
					}
					continue
				}
				// Plain Esc
				clearMenu()
				return "cancel", selectedIndex, nil
			}
			// 'k' or 'K' -> UP
			if slice[0] == 'k' || slice[0] == 'K' {
				if selectedIndex > 0 {
					selectedIndex--
					redraw()
				}
				slice = slice[1:]
				continue
			}
			// 'j' or 'J' -> DOWN
			if slice[0] == 'j' || slice[0] == 'J' {
				if selectedIndex < len(items)-1 {
					selectedIndex++
					redraw()
				}
				slice = slice[1:]
				continue
			}
			// 'n' or 'N' (new session)
			if allowNew && (slice[0] == 'n' || slice[0] == 'N') {
				clearMenu()
				return "new", selectedIndex, nil
			}
			// 'q' or 'Q' (quit / cancel)
			if slice[0] == 'q' || slice[0] == 'Q' {
				clearMenu()
				return "cancel", selectedIndex, nil
			}

			slice = slice[1:]
		}
	}
}

// runInteractiveSelect wraps runInteractiveSelectEngine in terminal raw mode.
func runInteractiveSelect(
	title string,
	items []SelectChoice,
	initialIndex int,
	allowNew bool,
) (string, int, error) {
	fd := int(os.Stdin.Fd())
	if !term.IsTerminal(fd) {
		return "cancel", initialIndex, fmt.Errorf("stdin is not a terminal")
	}

	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return "cancel", initialIndex, err
	}
	defer func() {
		_ = term.Restore(fd, oldState)
	}()

	getTermSize := func() (int, int) {
		w, h, err := term.GetSize(fd)
		if err != nil || w <= 0 || h <= 0 {
			return 80, 24
		}
		return w, h
	}

	return runInteractiveSelectEngine(os.Stdin, os.Stdout, title, items, initialIndex, allowNew, getTermSize)
}

// selectInteractiveModel presents an interactive selector for choosing the active AI model.
func selectInteractiveModel() {
	models.ModelCfgMu.RLock()
	cfg := models.ModelCfg
	models.ModelCfgMu.RUnlock()

	if cfg == nil || len(cfg.Models) == 0 {
		fmt.Println("No models configured.")
		return
	}

	active := cfg.AgentModel
	var names []string
	for name := range cfg.Models {
		names = append(names, name)
	}
	sort.Strings(names)

	activeIdx := 0
	for i, name := range names {
		if name == active {
			activeIdx = i
			break
		}
	}

	var choices []SelectChoice
	for _, name := range names {
		m := cfg.Models[name]
		isActive := name == active
		keyStatus := models.KeySourceLabel(&m)

		activeLabel := ""
		if isActive {
			activeLabel = "★ ACTIVE"
		}

		choices = append(choices, SelectChoice{
			ID: name,
			Render: func(selected bool) string {
				if selected {
					return fmt.Sprintf("▶ \033[1;36m%-33s\033[0m  %-15s  %-10s  %s",
						name, m.Provider, activeLabel, keyStatus)
				}
				return fmt.Sprintf("  %-33s  %-15s  %-10s  %s",
					name, m.Provider, activeLabel, keyStatus)
			},
		})
	}

	action, idx, err := runInteractiveSelect(
		"Select Active AI Model (↑/↓ to navigate, Enter to select, Esc to cancel):",
		choices,
		activeIdx,
		false,
	)
	if err != nil || action != "select" {
		return
	}

	chosen := choices[idx].ID
	if err := models.SwitchActiveModel(chosen); err != nil {
		fmt.Printf("❌ Failed to switch model: %v\n", err)
	} else {
		fmt.Printf("✓ Active model switched to: %s\n", chosen)
	}
}

// selectInteractiveSession presents an interactive selector for choosing or creating a chat session.
func selectInteractiveSession() {
	sessions := agent.ListSessions()

	// Ensure currentSessionID is in the list
	foundCurrent := false
	for _, s := range sessions {
		if s.ID == currentSessionID {
			foundCurrent = true
			break
		}
	}
	if !foundCurrent {
		sessions = append([]agent.SessionMeta{
			{
				ID:           currentSessionID,
				MsgCount:     0,
				LastModified: time.Now(),
				LastPreview:  "(current active session)",
			},
		}, sessions...)
	}

	activeIdx := 0
	for i, s := range sessions {
		if s.ID == currentSessionID {
			activeIdx = i
			break
		}
	}

	var choices []SelectChoice
	for _, s := range sessions {
		sess := s
		isActive := sess.ID == currentSessionID

		marker := " "
		if isActive {
			marker = "●"
		}

		info := fmt.Sprintf("(%d msgs, %s, %s)",
			sess.MsgCount, sess.LastModified.Format("02 Jan 15:04"), sess.LastPreview)

		choices = append(choices, SelectChoice{
			ID: sess.ID,
			Render: func(selected bool) string {
				if selected {
					return fmt.Sprintf("▶ \033[1;36m%-24s\033[0m %s %s", sess.ID, marker, info)
				}
				return fmt.Sprintf("  %-24s %s %s", sess.ID, marker, info)
			},
		})
	}

	action, idx, err := runInteractiveSelect(
		"Select Active Session (↑/↓ to navigate, Enter to switch, 'n' for new, Esc to cancel):",
		choices,
		activeIdx,
		true,
	)
	if err != nil || action == "cancel" {
		return
	}

	if action == "new" {
		fmt.Print("Enter new session name (leave blank for auto): ")
		scanner := bufio.NewScanner(os.Stdin)
		name := ""
		if scanner.Scan() {
			name = strings.TrimSpace(scanner.Text())
		}
		if name == "" {
			name = fmt.Sprintf("chat-%s", time.Now().Format("0102-150405"))
		} else {
			name = strings.Trim(name, "\"'")
		}
		newLock, err := acquireSessionLock(name)
		if err != nil {
			fmt.Printf("❌ %v\n", err)
			return
		}
		if currentSessionLock != nil {
			currentSessionLock.Release()
		}
		currentSessionLock = newLock
		currentSessionID = name
		fmt.Printf("✓ Switched to new session: \033[1;32m%s\033[0m\n", currentSessionID)
		if strings.HasPrefix(name, "chat-") {
			fmt.Println("\033[2mℹ️ This session will be auto-named after your first message.\033[0m")
		}
		return
	}

	if action == "select" {
		chosenID := choices[idx].ID
		if chosenID == currentSessionID {
			fmt.Printf("ℹ️ Already on session '%s'\n", chosenID)
			return
		}
		newLock, err := acquireSessionLock(chosenID)
		if err != nil {
			fmt.Printf("❌ %v\n", err)
			return
		}
		if currentSessionLock != nil {
			currentSessionLock.Release()
		}
		currentSessionLock = newLock
		currentSessionID = chosenID
		fmt.Printf("✓ Active session switched to: \033[1;32m%s\033[0m\n", currentSessionID)
	}
}
