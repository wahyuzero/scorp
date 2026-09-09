//go:build windows

package tools

import (
	"os"
)

func fileLockExclusive(f *os.File) error {
	return nil
}

func fileUnlock(f *os.File) error {
	return nil
}
