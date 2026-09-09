package tools

import (
	"os"
	"scorp-agent/internal/helpers"
)

// WriteFileAtomic delegates to internal/helpers.WriteFileAtomic
func WriteFileAtomic(path string, data []byte, perm os.FileMode) error {
	return helpers.WriteFileAtomic(path, data, perm)
}
