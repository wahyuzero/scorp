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


// SplitMessage splits text into chunks of maxLen runes, trying to break on newlines.
func SplitMessage(text string, maxLen int) []string {
	if len([]rune(text)) <= maxLen {
		return []string{text}
	}
	var chunks []string
	var current strings.Builder
	for _, line := range strings.Split(text, "\n") {
		lineLen := len([]rune(line))
		if len([]rune(current.String()))+lineLen+1 > maxLen {
			if current.Len() > 0 {
				chunks = append(chunks, current.String())
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
		chunks = append(chunks, current.String())
	}
	return chunks
}
