package helpers

import (
	"encoding/json"
	"strings"
)

// TruncateStr truncates a string to n characters
func TruncateStr(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

// GetStringArg extracts a string argument from a map with a default value
func GetStringArg(args map[string]interface{}, key, defaultVal string) string {
	if v, ok := args[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return defaultVal
}

// GetIntArg extracts an int argument from a map with a default value
func GetIntArg(args map[string]interface{}, key string, defaultVal int) int {
	if v, ok := args[key]; ok {
		switch n := v.(type) {
		case float64:
			return int(n)
		case int:
			return n
		}
	}
	return defaultVal
}

// GetFloatArg extracts a float argument from a map with a default value
func GetFloatArg(args map[string]interface{}, key string, def float64) float64 {
	if v, ok := args[key]; ok {
		switch n := v.(type) {
		case float64:
			return n
		case float32:
			return float64(n)
		case int:
			return float64(n)
		case int64:
			return float64(n)
		}
	}
	return def
}

// GetStringSliceArg extracts a string slice argument from a map
func GetStringSliceArg(args map[string]interface{}, key string) []string {
	if v, ok := args[key]; ok {
		if arr, ok := v.([]interface{}); ok {
			result := make([]string, 0, len(arr))
			for _, item := range arr {
				if s, ok := item.(string); ok {
					result = append(result, s)
				}
			}
			return result
		}
	}
	return nil
}

// GetBoolArg extracts a bool argument from a map with a default value
func GetBoolArg(args map[string]interface{}, key string, def bool) bool {
	if v, ok := args[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return def
}

// GetInt64Arg extracts an int64 argument from a map with a default value
func GetInt64Arg(args map[string]interface{}, key string, def int64) int64 {
	if v, ok := args[key]; ok {
		switch n := v.(type) {
		case float64:
			return int64(n)
		case int64:
			return n
		case int:
			return int64(n)
		case json.Number:
			if i, err := n.Int64(); err == nil {
				return i
			}
		}
	}
	return def
}

// TruncOutput truncates a string with "(truncated)" suffix
func TruncOutput(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "\n... (truncated)"
}

// EscapeHTML escapes HTML special characters for Telegram
func EscapeHTML(s string) string {
	r := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;")
	return r.Replace(s)
}


// MaxToolOutput is the maximum length for tool output before truncation
const MaxToolOutput = 3000

// TruncOutputTool truncates tool output to MaxToolOutput characters
func TruncOutputTool(output string) string {
	return TruncateStr(output, MaxToolOutput)
}


// SplitMessage splits text into chunks of maxLen runes, ensuring HTML tags are properly balanced across chunks.
func SplitMessage(text string, maxLen int) []string {
	if len([]rune(text)) <= maxLen {
		return []string{text}
	}
	var rawChunks []string
	var current strings.Builder
	for _, line := range strings.Split(text, "\n") {
		lineLen := len([]rune(line))
		if len([]rune(current.String()))+lineLen+1 > maxLen {
			if current.Len() > 0 {
				rawChunks = append(rawChunks, current.String())
				current.Reset()
			}
			current.WriteString(line)
		} else {
			if current.Len() > 0 {
				current.WriteByte('\n')
			}
			current.WriteString(line)
		}
	}
	if current.Len() > 0 {
		rawChunks = append(rawChunks, current.String())
	}

	// Balance HTML tags across chunks to prevent Telegram entity parse errors
	var balanced []string
	var openTags []string
	for i, chunk := range rawChunks {
		var prefix strings.Builder
		for _, tag := range openTags {
			prefix.WriteString("<" + tag + ">")
		}
		chunkText := prefix.String() + chunk

		openTags = getUnclosedTags(chunkText)
		var suffix strings.Builder
		for j := len(openTags) - 1; j >= 0; j-- {
			suffix.WriteString("</" + openTags[j] + ">")
		}
		if i < len(rawChunks)-1 && len(openTags) > 0 {
			chunkText += suffix.String()
		}
		balanced = append(balanced, chunkText)
	}

	return balanced
}

func getUnclosedTags(htmlText string) []string {
	supported := []string{"pre", "code", "b", "i", "u", "s"}
	var stack []string
	for _, tag := range supported {
		openCount := strings.Count(strings.ToLower(htmlText), "<"+tag+">")
		closeCount := strings.Count(strings.ToLower(htmlText), "</"+tag+">")
		if openCount > closeCount {
			for k := 0; k < openCount-closeCount; k++ {
				stack = append(stack, tag)
			}
		}
	}
	return stack
}
