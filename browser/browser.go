// +build !nobrowser

package browser

import (
	"scorp-agent/internal/helpers"
	"fmt"
)

// ──────────────────────────────────────────────
// Browser Control Tool — Headless Chrome via chromedp
// ──────────────────────────────────────────────

// executeBrowser handles the "browser" tool — navigate, screenshot, click, type, extract text
func ExecuteBrowser(args map[string]interface{}, chatID int64) (string, bool) {
	action := helpers.GetStringArg(args, "action", "")
	url := helpers.GetStringArg(args, "url", "")
	selector := helpers.GetStringArg(args, "selector", "")
	text := helpers.GetStringArg(args, "text", "")
	code := helpers.GetStringArg(args, "code", "")
	direction := helpers.GetStringArg(args, "direction", "")

	switch action {
	case "goto", "navigate":
		if url == "" {
			return "Error: url required for goto", false
		}
		return browserSessionNavigate(url, chatID)
	case "screenshot":
		return browserSessionScreenshot(url, chatID)
	case "click":
		if selector == "" {
			return "Error: selector required for click", false
		}
		return browserSessionClick(selector, url, chatID)
	case "type":
		if selector == "" || text == "" {
			return "Error: selector and text required for type", false
		}
		return browserSessionFill(selector, text, url, chatID)
	case "evaluate":
		if code == "" {
			return "Error: code required for evaluate", false
		}
		return browserSessionEvaluate(code, url, chatID)
	case "scroll":
		if direction == "" {
			direction = "down"
		}
		return browserSessionScroll(direction, url, chatID)
	case "extract":
		if selector == "" {
			return "Error: selector required for extract", false
		}
		return browserSessionExtract(selector, url, chatID)
	case "snapshot":
		return browserSnapshot(chatID)
	case "console":
		return browserConsole(chatID)
	case "reset":
		if CloseBrowserSession(chatID) {
			return "🔄 Browser session reset (cookies + state cleared).", true
		}
		return "No active session to reset.", true
	case "list":
		return listBrowserSessions(), true
	default:
		return fmt.Sprintf("Unknown browser action: %s (use goto, screenshot, click, type, evaluate, scroll, extract, snapshot, console, reset, list)", action), false
	}
}

