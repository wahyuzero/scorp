// +build nobrowser

package tools

// ExecuteAutoLogin stub for nobrowser builds
func ExecuteAutoLogin(args map[string]interface{}, chatID int64) (string, bool) {
	return "Error: browser feature disabled in this build (compiled with nobrowser tag)", false
}
