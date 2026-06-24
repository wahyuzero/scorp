package tools

import (
	"scorp-agent/internal/helpers"
	"fmt"
	"os"
	"strings"
)

// ──────────────────────────────────────────────
// Patch Tool — Fuzzy find-and-replace file editing
// ──────────────────────────────────────────────

func ExecutePatch(args map[string]interface{}) (string, bool) {
	mode := helpers.GetStringArg(args, "mode", "replace")

	switch mode {
	case "replace":
		return patchReplace(args)
	default:
		return "Error: mode must be 'replace'", false
	}
}

// patchReplace finds old_string in a file and replaces with new_string.
// Tries 3 matching strategies in order: exact, trim-trailing-WS, normalize-all-WS.
func patchReplace(args map[string]interface{}) (string, bool) {
	path := helpers.GetStringArg(args, "path", "")
	oldStr := helpers.GetStringArg(args, "old_string", "")
	newStr := helpers.GetStringArg(args, "new_string", "")
	replaceAll := helpers.GetBoolArg(args, "replace_all", false)

	if path == "" {
		return "Error: 'path' is required", false
	}
	if oldStr == "" {
		return "Error: 'old_string' is required", false
	}

	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Sprintf("Error reading file: %v", err), false
	}
	content := string(data)

	// Try matching strategies in order of precision
	result, matchCount, strategy := tryMatchStrategies(content, oldStr, newStr, replaceAll)

	if matchCount == 0 {
		return fmt.Sprintf("Error: old_string not found in %s.\nTried: exact match, trimmed match, whitespace-normalized match.", path), false
	}

	if !replaceAll && matchCount > 1 {
		// Show first 3 occurrences for context
		return fmt.Sprintf("Error: old_string is not unique — found %d matches in %s.\nAdd more surrounding lines to make it unique, or set replace_all=true.", matchCount, path), false
	}

	// Write
	if err := os.WriteFile(path, []byte(result), 0644); err != nil {
		return fmt.Sprintf("Error writing file: %v", err), false
	}

	// Build diff preview
	diff := buildDiffPreview(oldStr, newStr)

	stratLabel := map[int]string{0: "exact", 1: "trim", 2: "normalize"}[strategy]

	return fmt.Sprintf("✅ Patched %s (%d match, %s strategy)\n\n%s",
		path, matchCount, stratLabel, diff), true
}

// tryMatchStrategies tries 3 strategies, returns first that matches.
// Returns: (modifiedContent, matchCount, strategyUsed)
func tryMatchStrategies(content, oldStr, newStr string, replaceAll bool) (string, int, int) {
	// Strategy 0: Exact match
	if strings.Contains(content, oldStr) {
		count := strings.Count(content, oldStr)
		if replaceAll {
			return strings.ReplaceAll(content, oldStr, newStr), count, 0
		}
		return strings.Replace(content, oldStr, newStr, 1), count, 0
	}

	// Strategy 1: Trim trailing whitespace per line
	if result, count, ok := lineWindowMatch(content, oldStr, newStr, replaceAll, false); ok {
		return result, count, 1
	}

	// Strategy 2: Normalize all internal whitespace
	if result, count, ok := lineWindowMatch(content, oldStr, newStr, replaceAll, true); ok {
		return result, count, 2
	}

	return content, 0, 0
}

// lineWindowMatch does a sliding window over content lines.
// If collapseWS=true, normalizes whitespace within each line before comparing.
// Returns (modifiedContent, matchCount, found).
func lineWindowMatch(content, oldStr, newStr string, replaceAll bool, collapseWS bool) (string, int, bool) {
	contentLines := strings.Split(content, "\n")
	oldLines := splitLines(oldStr)
	newLines := splitLines(newStr)

	if len(oldLines) == 0 || len(oldLines) > len(contentLines) {
		return content, 0, false
	}

	// Pre-normalize old lines for comparison
	normOldLines := make([]string, len(oldLines))
	for i, l := range oldLines {
		normOldLines[i] = normalizeForCompare(l, collapseWS)
	}

	count := 0
	result := make([]string, 0, len(contentLines))
	i := 0
	for i < len(contentLines) {
		// Check if window matches
		if i+len(oldLines) <= len(contentLines) {
			matched := true
			for j := 0; j < len(oldLines); j++ {
				cl := normalizeForCompare(contentLines[i+j], collapseWS)
				if cl != normOldLines[j] {
					matched = false
					break
				}
			}
			if matched {
				count++
				result = append(result, newLines...)
				i += len(oldLines)
				if !replaceAll {
					result = append(result, contentLines[i:]...)
					return strings.Join(result, "\n"), count, true
				}
				continue
			}
		}
		result = append(result, contentLines[i])
		i++
	}

	if count == 0 {
		return content, 0, false
	}
	return strings.Join(result, "\n"), count, true
}

// normalizeForCompare normalizes a line for fuzzy matching.
// If collapseWS=true, collapses all whitespace to single spaces.
// Always trims leading/trailing whitespace.
func normalizeForCompare(line string, collapseWS bool) string {
	line = strings.TrimSpace(line)
	if collapseWS {
		fields := strings.Fields(line)
		return strings.Join(fields, " ")
	}
	return line
}

// splitLines splits on \n, preserving the behavior of handling \r\n
func splitLines(s string) []string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	return strings.Split(s, "\n")
}

// buildDiffPreview generates a simple unified diff preview
func buildDiffPreview(oldStr, newStr string) string {
	var sb strings.Builder
	sb.WriteString("```diff\n")
	for _, line := range splitLines(oldStr) {
		sb.WriteString("- " + line + "\n")
	}
	for _, line := range splitLines(newStr) {
		sb.WriteString("+ " + line + "\n")
	}
	sb.WriteString("```")
	return sb.String()
}
