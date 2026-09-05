package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

// sessionLockFile holds the file descriptor for an acquired session lock
type sessionLockFile struct {
	file *os.File
	path string
}

// acquireSessionLock attempts to acquire an exclusive non-blocking advisory file lock
// for a session ID. If already locked by another running process, it returns an error
// with the other process's PID.
func acquireSessionLock(sessionID string) (*sessionLockFile, error) {
	safeID := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			return r
		}
		return '_'
	}, sessionID)

	lockDir := filepath.Join(os.TempDir(), "scorp_locks")
	_ = os.MkdirAll(lockDir, 0755)
	lockPath := filepath.Join(lockDir, safeID+".lock")

	f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open session lockfile: %w", err)
	}

	err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		// Locked by another active process. Read existing PID
		buf := make([]byte, 32)
		n, _ := f.Read(buf)
		otherPID := strings.TrimSpace(string(buf[:n]))
		_ = f.Close()
		return nil, fmt.Errorf("session '%s' is already running in another process (PID %s). To avoid runaway processes or corrupted state, wait or terminate that PID first", sessionID, otherPID)
	}

	// Write current PID into lockfile
	_ = f.Truncate(0)
	_, _ = f.Seek(0, 0)
	_, _ = f.WriteString(strconv.Itoa(os.Getpid()) + "\n")
	_ = f.Sync()

	return &sessionLockFile{file: f, path: lockPath}, nil
}

// Release releases the lock and cleans up the file
func (l *sessionLockFile) Release() {
	if l != nil && l.file != nil {
		_ = syscall.Flock(int(l.file.Fd()), syscall.LOCK_UN)
		_ = l.file.Close()
		_ = os.Remove(l.path)
		l.file = nil
	}
}
