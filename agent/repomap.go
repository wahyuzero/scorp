package agent

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// ──────────────────────────────────────────────
// Repository Map Prefix (Prompt Caching & Context)
// Generates a compact structural tree of the workspace
// with cached invalidation to maximize prompt cache hits.
// ──────────────────────────────────────────────

var (
	repoMapCache    string
	repoMapCacheTime time.Time
	repoMapMu       sync.Mutex
	repoMapTTL      = 45 * time.Second
)

var ignoredDirs = map[string]bool{
	".git":         true,
	"node_modules": true,
	"vendor":       true,
	".venv":        true,
	"venv":         true,
	"__pycache__":  true,
	"dist":         true,
	"build":        true,
	"graphify-out": true,
	".idea":        true,
	".vscode":      true,
	".cache":       true,
}

var ignoredExts = map[string]bool{
	".exe": true,
	".bin": true,
	".so":  true,
	".dylib": true,
	".png": true,
	".jpg": true,
	".jpeg": true,
	".gif": true,
	".ico": true,
	".pdf": true,
	".zip": true,
	".tar": true,
	".gz":  true,
}

// GetRepoMap returns a cached compact directory layout of the current working directory.
func GetRepoMap() string {
	repoMapMu.Lock()
	defer repoMapMu.Unlock()

	if repoMapCache != "" && time.Since(repoMapCacheTime) < repoMapTTL {
		return repoMapCache
	}

	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}

	var sb strings.Builder
	fileCount := 0
	const maxFiles = 60

	_ = filepath.Walk(cwd, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		rel, err := filepath.Rel(cwd, path)
		if err != nil || rel == "." {
			return nil
		}

		parts := strings.Split(rel, string(filepath.Separator))
		// Check ignored dirs
		for _, part := range parts {
			if ignoredDirs[part] || strings.HasPrefix(part, ".git") {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		// Depth limit: 3 levels
		if len(parts) > 3 {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if info.IsDir() {
			indent := strings.Repeat("  ", len(parts)-1)
			sb.WriteString(fmt.Sprintf("%s📁 %s/\n", indent, parts[len(parts)-1]))
			return nil
		}

		// Filter extensions
		ext := strings.ToLower(filepath.Ext(info.Name()))
		if ignoredExts[ext] {
			return nil
		}

		if fileCount < maxFiles {
			indent := strings.Repeat("  ", len(parts)-1)
			sb.WriteString(fmt.Sprintf("%s📄 %s\n", indent, parts[len(parts)-1]))
			fileCount++
		}

		return nil
	})

	if sb.Len() == 0 {
		return ""
	}

	repoMapCache = fmt.Sprintf("```\n%s```", sb.String())
	repoMapCacheTime = time.Now()
	return repoMapCache
}

// InvalidateRepoMap forces refresh on next request
func InvalidateRepoMap() {
	repoMapMu.Lock()
	repoMapCache = ""
	repoMapMu.Unlock()
}
