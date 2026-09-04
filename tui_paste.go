package main

import "fmt"

// Enable bracketed paste mode on raw terminal
func enableBracketedPaste() {
	fmt.Print("\033[?2004h")
}

// Disable bracketed paste mode
func disableBracketedPaste() {
	fmt.Print("\033[?2004l")
}
