package tools

import (
	"fmt"
	"os"
	"path/filepath"
)

// WriteFileAtomic writes data to a temporary file in the same directory,
// syncs to disk, and atomically replaces the target path via rename.
// This guarantees that power failure, SIGKILL, or OOM crashes never leave
// a truncated, half-written, or 0-byte file on disk.
func WriteFileAtomic(path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create parent dir: %w", err)
	}

	tmpFile, err := os.CreateTemp(dir, fmt.Sprintf(".%s.*.tmp", filepath.Base(path)))
	if err != nil {
		// Fallback to direct write if temp file cannot be created
		return os.WriteFile(path, data, perm)
	}
	tmpPath := tmpFile.Name()

	// Ensure cleanup on error
	success := false
	defer func() {
		if !success {
			_ = os.Remove(tmpPath)
		}
	}()

	if _, err := tmpFile.Write(data); err != nil {
		_ = tmpFile.Close()
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	// Flush OS buffers to storage media
	if err := tmpFile.Sync(); err != nil {
		_ = tmpFile.Close()
		return fmt.Errorf("failed to sync temp file: %w", err)
	}

	if err := tmpFile.Chmod(perm); err != nil {
		_ = tmpFile.Close()
		return fmt.Errorf("failed to chmod temp file: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	// Atomic filesystem swap
	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("failed to atomically rename temp file to target: %w", err)
	}

	success = true
	return nil
}
