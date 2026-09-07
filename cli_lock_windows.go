//go:build windows

package main

import (
	"os"
)

func lockFileExclusive(f *os.File) error {
	// Standard non-blocking file access
	return nil
}

func unlockFile(f *os.File) error {
	return nil
}
