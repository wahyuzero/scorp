// +build nobrowser

package main

// executeBrowser is a stub when browser tool is disabled via build tag
func executeBrowser(args map[string]interface{}, chatID int64) (string, bool) {
	return "❌ Browser tool disabled (build without -tags nobrowser to enable)", false
}

// createBrowserContext is a stub
func createBrowserContext() (interface{}, func()) {
	return nil, func() {}
}