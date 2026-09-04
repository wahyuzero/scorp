//go:build nobrowser
// +build nobrowser

package tools

// InitMonitor stub for nobrowser builds
func InitMonitor() {}

// ExecuteMonitor stub for nobrowser builds
func ExecuteMonitor(args map[string]interface{}, chatID int64) (string, bool) {
	return "Error: monitor feature disabled in this build (compiled with nobrowser tag)", false
}
