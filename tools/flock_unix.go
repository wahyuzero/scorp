//go:build !windows

package tools

import (
	"os"
	"syscall"
)

// FileLockExclusive acquires an exclusive advisory lock on the file
func FileLockExclusive(f *os.File) error {
	return syscall.Flock(int(f.Fd()), syscall.LOCK_EX)
}

// FileUnlock releases an advisory lock on the file
func FileUnlock(f *os.File) error {
	return syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
}

func fileLockExclusive(f *os.File) error {
	return FileLockExclusive(f)
}

func fileUnlock(f *os.File) error {
	return FileUnlock(f)
}
