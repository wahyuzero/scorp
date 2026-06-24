// +build nobrowser

package browser

// executeBrowser is a stub when browser tool is disabled via build tag
func ExecuteBrowser(args map[string]interface{}, chatID int64) (string, bool) {
	return "❌ Browser tool disabled (build without -tags nobrowser to enable)", false
}

// createBrowserContext is a stub
func createBrowserContext() (interface{}, func()) {
	return nil, func() {}
}

func initVault() {
	// No-op when browser tools are disabled
}

func CleanupStaleBrowserSessions() {
	// No-op when browser tools are disabled
}

func initMonitor() {
	// No-op when browser tools are disabled
}