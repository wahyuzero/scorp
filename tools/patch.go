package tools

import (
	"scorp-agent/internal/helpers"
	"fmt"
	"os"
	"strings"
)

// ──────────────────────────────────────────────
// Surgical Diff / Chunk Replacement Tool
// ──────────────────────────────────────────────

// ExecuteReplaceFileContent implements modern chunk replacement (exact or fuzzy match)
// to perform surgical edits without rewriting entire files.
func ExecuteReplaceFileContent(args map[string]interface{}) (string, bool) {
	return patchReplace(args)
}

// ExecutePatch handles patch invocations (backward compatible)
func ExecutePatch(args map[string]interface{}) (string, bool) {
	mode := helpers.GetStringArg(args, "mode", "replace")

	switch mode {
	case "replace", "":
		return patchReplace(args)
	default:
		return "Error: mode must be 'replace'", false
	}
}

// patchReplace finds target_content/old_string in a file and replaces with replacement_content/new_string.
// Tries 3 matching strategies in order: exact, trim-trailing-WS, normalize-all-WS.
func patchReplace(args map[string]interface{}) (string, bool) {
	// Support both naming schemes: modern blueprint (target_file, target_content, replacement_content)
	// and classical patch (path, old_string, new_string)
	path := helpers.GetStringArg(args, "target_file", "")
	if path == "" {
		path = helpers.GetStringArg(args, "path", "")
	}
	if path == "" {
		path = helpers.GetStringArg(args, "file", "")
	}

	oldStr := helpers.GetStringArg(args, "target_content", "")
	if oldStr == "" {
		oldStr = helpers.GetStringArg(args, "old_string", "")
	}
	if oldStr == "" {
		oldStr = helpers.GetStringArg(args, "old_content", "")
	}

	newStr := helpers.GetStringArg(args, "replacement_content", "")
	if newStr == "" {
		newStr = helpers.GetStringArg(args, "new_string", "")
	}
	if newStr == "" {
		newStr = helpers.GetStringArg(args, "new_content", "")
	}

	startLine := helpers.GetIntArg(args, "start_line", 0)
	endLine := helpers.GetIntArg(args, "end_line", 0)
	replaceAll := helpers.GetBoolArg(args, "replace_all", false)

	if path == "" {
		return "Error: 'target_file' (or 'path') is required", false
	}
	if oldStr == "" {
		return "Error: 'target_content' (or 'old_string') is required", false
	}

	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Sprintf("Error reading file: %v", err), false
	}
	content := string(data)

	// If startLine/endLine specified, attempt scoped replacement first
	if startLine > 0 {
		if scopedResult, ok, msg := scopedLineReplace(content, oldStr, newStr, startLine, endLine); ok {
			if err := os.WriteFile(path, []byte(scopedResult), 0644); err != nil {
				return fmt.Sprintf("Error writing file: %v", err), false
			}
			diff := buildDiffPreview(oldStr, newStr)
			return fmt.Sprintf("✅ Surgically edited %s (lines %d-%d)\n\n%s", path, startLine, endLine, diff), true
		} else if msg != "" && endLine > 0 {
			// If line-scoped attempt explicitly failed, log but fall through to whole-file match
		}
	}

	// Try matching strategies in order of precision across whole file
	result, matchCount, strategy := tryMatchStrategies(content, oldStr, newStr, replaceAll)

	if matchCount == 0 {
		return fmt.Sprintf("Error: target_content not found in %s.\nTried: exact match, trimmed match, whitespace-normalized match. Make sure target_content matches the file content.", path), false
	}

	if !replaceAll && matchCount > 1 {
		return fmt.Sprintf("Error: target_content is not unique — found %d matches in %s.\nSpecify start_line/end_line, add more surrounding lines to make it unique, or set replace_all=true.", matchCount, path), false
	}

	// Write
	if err := os.WriteFile(path, []byte(result), 0644); err != nil {
		return fmt.Sprintf("Error writing file: %v", err), false
	}

	// Build diff preview
	diff := buildDiffPreview(oldStr, newStr)

	stratLabel := map[int]string{0: "exact", 1: "trim", 2: "normalize"}[strategy]

	return fmt.Sprintf("✅ Surgically edited %s (%d match, %s strategy)\n\n%s",
		path, matchCount, stratLabel, diff), true
}

// scopedLineReplace replaces oldStr with newStr within specific 1-based line bounds
func scopedLineReplace(content, oldStr, newStr string, startLine, endLine int) (string, bool, string) {
	lines := splitLines(content)
	totalLines := len(lines)
	if startLine < 1 || startLine > totalLines {
		return "", false, "start_line out of bounds"
	}
	if endLine < startLine || endLine > totalLines {
		endLine = totalLines
	}

	// Extract lines within range (0-indexed)
	prefix := lines[:startLine-1]
	window := lines[startLine-1 : endLine]
	suffix := lines[endLine:]

	windowContent := strings.Join(window, "\n")
	res, count, _ := tryMatchStrategies(windowContent, oldStr, newStr, false)
	if count == 1 {
		var finalLines []string
		if len(prefix) > 0 {
			finalLines = append(finalLines, prefix...)
		}
		finalLines = append(finalLines, res)
		if len(suffix) > 0 {
			finalLines = append(finalLines, suffix...)
		}
		return strings.Join(finalLines, "\n"), true, ""
	}
	return "", false, "no unique match within line scope"
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
