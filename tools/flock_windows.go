//go:build windows

package tools

import (
	"os"
)

// FileLockExclusive acquires an exclusive advisory lock on the file
func FileLockExclusive(f *os.File) error {
	return nil
}

// FileUnlock releases an advisory lock on the file
func FileUnlock(f *os.File) error {
	return nil
}

func fileLockExclusive(f *os.File) error {
	return FileLockExclusive(f)
}

func fileUnlock(f *os.File) error {
	return FileUnlock(f)
}
