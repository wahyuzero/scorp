package tools

import (
	"scorp-agent/internal/helpers"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// ──────────────────────────────────────────────
// Code Search Tool — ripgrep wrapper
// ──────────────────────────────────────────────

func ExecuteSearchCode(args map[string]interface{}) (string, bool) {
	pattern := helpers.GetStringArg(args, "pattern", "")
	if pattern == "" {
		return "Error: 'pattern' argument is required", false
	}

	path := helpers.GetStringArg(args, "path", ".")
	fileGlob := helpers.GetStringArg(args, "glob", "")
	maxResults := helpers.GetIntArg(args, "max_results", 20)
	contextLines := helpers.GetIntArg(args, "context", 0)

	if maxResults > 100 {
		maxResults = 100
	}

	// Check if rg is available; fallback to standard grep if absent
	hasRg := false
	if _, err := exec.LookPath("rg"); err == nil {
		hasRg = true
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var cmd *exec.Cmd

	if hasRg {
		// Build rg command
		rgArgs := []string{
			"--no-heading",
			"--line-number",
			"--color", "never",
			fmt.Sprintf("--max-count=%d", maxResults),
		}

		if contextLines > 0 {
			rgArgs = append(rgArgs, fmt.Sprintf("--context=%d", contextLines))
		}

		if fileGlob != "" {
			rgArgs = append(rgArgs, "--glob", fileGlob)
		}

		rgArgs = append(rgArgs, pattern, path)
		cmd = exec.CommandContext(ctx, "rg", rgArgs...)
	} else {
		// Fallback to standard grep
		grepArgs := []string{"-rn"}
		if contextLines > 0 {
			grepArgs = append(grepArgs, fmt.Sprintf("-C%d", contextLines))
		}
		if fileGlob != "" {
			grepArgs = append(grepArgs, fmt.Sprintf("--include=%s", fileGlob))
		}
		grepArgs = append(grepArgs, "-E", pattern, path)
		cmd = exec.CommandContext(ctx, "grep", grepArgs...)
	}

	output, err := cmd.CombinedOutput()
	result := string(output)

	if ctx.Err() == context.DeadlineExceeded {
		return "Search timed out after 15s", false
	}

	// rg/grep returns exit code 1 for no matches — not an error
	if err != nil && !strings.Contains(err.Error(), "exit status 1") {
		return fmt.Sprintf("Search error: %v\n%s", err, helpers.TruncOutput(result, helpers.MaxToolOutput)), false
	}

	if result == "" {
		return fmt.Sprintf("No matches found for '%s' in %s", pattern, path), true
	}

	// Count matches
	matchCount := strings.Count(result, "\n") + 1
	header := fmt.Sprintf("Found %d matches for '%s' in %s:\n\n", matchCount, pattern, path)

	return helpers.TruncOutput(header+result, helpers.MaxToolOutput), true
}
