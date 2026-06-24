// +build !nobrowser

package browser

import (
	"scorp-agent/internal/helpers"
	"scorp-agent/config"
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
)

var whitespaceRe = regexp.MustCompile(`\s+`)

// ──────────────────────────────────────────────
// Browser Session Management — persistent chromedp contexts per chat
// ──────────────────────────────────────────────

type BrowserSession struct {
	Ctx       context.Context
	Cancel    context.CancelFunc
	ChatID    int64
	Created   time.Time
	LastUsed  time.Time
	CurrentURL string
}

var (
	browserSessions   = make(map[int64]*BrowserSession)
	browserSessionsMu sync.Mutex
)

// getOrCreateBrowserSession returns the browser session for the given chat ID,
// creating a new one if needed.
// GetOrCreateBrowserSession returns existing session or creates new one
func GetOrCreateBrowserSession(chatID int64) *BrowserSession {
	browserSessionsMu.Lock()
	defer browserSessionsMu.Unlock()

	if sess, ok := browserSessions[chatID]; ok {
		sess.LastUsed = time.Now()
		return sess
	}

	// Persistent user data dir for cookies, localStorage, etc.
	userDataDir := config.BrowserSessionPath(fmt.Sprintf("chat_%d", chatID))
	os.MkdirAll(userDataDir, 0755)

	// Create new browser context with persistent user data
	allocOpts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", "true"),
		chromedp.Flag("disable-gpu", "true"),
		chromedp.Flag("no-sandbox", "true"),
		chromedp.Flag("disable-dev-shm-usage", "true"),
		chromedp.Flag("no-first-run", "true"),
		chromedp.Flag("disable-extensions", "true"),
		chromedp.UserDataDir(userDataDir),
	)
	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), allocOpts...)

	ctx, cancel := chromedp.NewContext(allocCtx)

	sess := &BrowserSession{
		Ctx:      ctx,
		Cancel:   func() { cancel(); allocCancel() },
		ChatID:   chatID,
		Created:  time.Now(),
		LastUsed: time.Now(),
	}
	browserSessions[chatID] = sess

	// Initialize by running no-op
	chromedp.Run(ctx)

	log.Printf("[browser] Created persistent session for chat %d (data dir: %s)", chatID, userDataDir)
	return sess
}

// closeBrowserSession closes the browser session for the given chat ID.
// CloseBrowserSession closes the browser session for the given chat ID
func CloseBrowserSession(chatID int64) bool {
	browserSessionsMu.Lock()
	defer browserSessionsMu.Unlock()

	sess, ok := browserSessions[chatID]
	if !ok {
		return false
	}

	sess.Cancel()
	delete(browserSessions, chatID)
	log.Printf("[browser] Closed session for chat %d", chatID)
	return true
}

// listBrowserSessions returns info about all active sessions.
func listBrowserSessions() string {
	browserSessionsMu.Lock()
	defer browserSessionsMu.Unlock()

	if len(browserSessions) == 0 {
		return "No active browser sessions."
	}

	var sb strings.Builder
	sb.WriteString("🌐 <b>Browser Sessions</b>\n\n")
	for _, sess := range browserSessions {
		age := time.Since(sess.Created).Round(time.Second)
		sb.WriteString(fmt.Sprintf("• Chat <code>%d</code>: %s (age %s)\n",
			sess.ChatID,
			sess.CurrentURL,
			age,
		))
	}
	return sb.String()
}

// cleanupStaleBrowserSessions removes sessions idle for more than 10 minutes.
// CleanupStaleBrowserSessions closes idle browser sessions
func CleanupStaleBrowserSessions() {
	browserSessionsMu.Lock()
	defer browserSessionsMu.Unlock()

	for id, sess := range browserSessions {
		if time.Since(sess.LastUsed) > 10*time.Minute {
			sess.Cancel()
			delete(browserSessions, id)
			log.Printf("[browser] Cleaned up stale session for chat %d", id)
		}
	}
}

// ──────────────────────────────────────────────
// Session-based browser actions
// ──────────────────────────────────────────────

// browserSessionNavigate navigates using an existing or new session.
func browserSessionNavigate(url string, chatID int64) (string, bool) {
	sess := GetOrCreateBrowserSession(chatID)

	var title string
	err := chromedp.Run(sess.Ctx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
		chromedp.Title(&title),
		chromedp.Sleep(1*time.Second),
	)
	if err != nil {
		return fmt.Sprintf("Navigate error: %v", err), false
	}

	sess.CurrentURL = url

	var bodyText string
	chromedp.Run(sess.Ctx, chromedp.Text("body", &bodyText, chromedp.ByQuery))
	bodyText = whitespaceRe.ReplaceAllString(bodyText, " ")
	bodyText = strings.TrimSpace(bodyText)

	return helpers.TruncOutput(fmt.Sprintf("🌐 Navigated to: %s\nTitle: %s\n\n%s", url, title, bodyText), helpers.MaxToolOutput), true
}

// browserSessionScreenshot takes a screenshot using existing session (no re-navigate).
func browserSessionScreenshot(url string, chatID int64) (string, bool) {
	sess := GetOrCreateBrowserSession(chatID)

	// Only navigate if URL is provided and different from current
	if url != "" && url != sess.CurrentURL {
		if _, ok := browserSessionNavigate(url, chatID); !ok {
			return "Navigate failed", false
		}
	}

	var buf []byte
	err := chromedp.Run(sess.Ctx,
		chromedp.Sleep(1*time.Second),
		chromedp.FullScreenshot(&buf, 90),
	)
	if err != nil {
		return fmt.Sprintf("Screenshot error: %v", err), false
	}

	// Save screenshot
	ssDir := config.ScreenshotsDir()
	os.MkdirAll(ssDir, 0755)
	filename := fmt.Sprintf("%s/shot_%d.png", ssDir, time.Now().Unix())
	if err := os.WriteFile(filename, buf, 0644); err != nil {
		return fmt.Sprintf("Error saving screenshot: %v", err), false
	}

	SendFile(fmt.Sprintf("%d", chatID), filename)
	return fmt.Sprintf("📸 Screenshot saved & sent (%d KB). Path: %s", len(buf)/1024, filename), true
}

// browserSessionClick clicks a selector in the existing session.
// After clicking, returns rich feedback: URL change detection, title,
// interactive elements snapshot, and body text excerpt.
func browserSessionClick(selector, url string, chatID int64) (string, bool) {
	sess := GetOrCreateBrowserSession(chatID)

	if url != "" && url != sess.CurrentURL {
		browserSessionNavigate(url, chatID)
	}

	// Capture URL before click
	beforeURL := sess.CurrentURL

	var afterText, afterTitle, afterURL string

	// Use a timeout context — click may trigger navigation that can hang
	// Form submits can take longer, use 30s timeout
	clickCtx, cancel := context.WithTimeout(sess.Ctx, 30*time.Second)
	defer cancel()

	err := chromedp.Run(clickCtx,
		chromedp.Click(selector, chromedp.ByQuery),
		chromedp.Sleep(2*time.Second),
		chromedp.Title(&afterTitle),
		chromedp.Location(&afterURL),
		chromedp.Text("body", &afterText, chromedp.ByQuery),
	)
	if err != nil {
		// Click triggered navigation and page might still be loading — try to grab what we can
		log.Printf("[browser] click error (likely page navigation): %v", err)
		// Still return useful info
		chromedp.Run(clickCtx, chromedp.Location(&afterURL))
		sess.CurrentURL = afterURL
		urlChange := "⚠️ Page may still be loading (navigation triggered)"
		if afterURL != "" && beforeURL != afterURL {
			urlChange = "✅ Page changed! (URL redirect detected)"
		}
		return fmt.Sprintf("👆 Clicked '%s'\nURL: %s → %s\n%s\nTitle: %s\n\n📄 Note: Page was navigating. Use browser.snapshot or browser.extract to inspect current state.", selector, beforeURL, afterURL, urlChange, afterTitle), true
	}

	// Update session URL
	sess.CurrentURL = afterURL

	// Build URL change indicator
	urlChange := "✅ Page changed!"
	if beforeURL == afterURL {
		urlChange = "⚠️ URL unchanged (still on same page)"
	}

	afterText = whitespaceRe.ReplaceAllString(afterText, " ")
	afterText = strings.TrimSpace(afterText)
	if len(afterText) > 600 {
		afterText = afterText[:600] + "..."
	}

	// Auto-capture interactive elements snapshot (with timeout)
	var interactiveSnapshot string
	snapScript := `
		(function() {
			var interactive = document.querySelectorAll('a, button, input, select, textarea, [role="button"], [onclick]');
			var lines = [];
			interactive.forEach(function(el, i) {
				if (i >= 15) return; // limit
				var tag = el.tagName.toLowerCase();
				var type = el.type || '';
				var text = (el.innerText || el.value || el.placeholder || el.getAttribute('aria-label') || '').trim();
				if (text.length > 50) text = text.substring(0, 47) + '...';
				lines.push('  [' + i + '] <' + tag + (type ? ' type=' + type : '') + '>: ' + text);
			});
			return lines.join('\n') || '  (no interactive elements)';
		})()
	`
	snapCtx, snapCancel := context.WithTimeout(sess.Ctx, 5*time.Second)
	defer snapCancel()
	chromedp.Run(snapCtx, chromedp.Evaluate(snapScript, &interactiveSnapshot))

	return helpers.TruncOutput(fmt.Sprintf("👆 Clicked '%s'\nURL: %s → %s\n%s\nTitle: %s\n\n📋 Elements:\n%s\n\n📄 Page text:\n%s",
		selector, beforeURL, afterURL, urlChange, afterTitle, interactiveSnapshot, afterText), helpers.MaxToolOutput), true
}

// browserSessionFill fills a form field, returns context about what's next on the page.
func browserSessionFill(selector, text, url string, chatID int64) (string, bool) {
	sess := GetOrCreateBrowserSession(chatID)

	if url != "" && url != sess.CurrentURL {
		browserSessionNavigate(url, chatID)
	}

	err := chromedp.Run(sess.Ctx,
		chromedp.WaitReady(selector),
		chromedp.Clear(selector, chromedp.ByQuery),
		chromedp.SendKeys(selector, text, chromedp.ByQuery),
	)
	if err != nil {
		return fmt.Sprintf("Fill error: %v", err), false
	}

	// Check what submit buttons are available after filling
	var submitInfo string
	checkScript := `
		(function() {
			var btns = document.querySelectorAll('button[type="submit"], input[type="submit"], button:not([type]), [role="button"]');
			var lines = [];
			btns.forEach(function(b, i) {
				if (i >= 5) return;
				var t = (b.innerText || b.value || '').trim();
				var sel = b.id ? '#' + b.id : (b.className ? 'button.' + (typeof b.className === 'string' ? b.className.split(' ')[0] : '') : 'button');
				lines.push(sel + ': "' + t + '"');
			});
			return lines.join(', ') || 'no submit found';
		})()
	`
	chromedp.Run(sess.Ctx, chromedp.Evaluate(checkScript, &submitInfo))

	return fmt.Sprintf("✏️ Filled '%s' with '%s'\n💡 Ready to submit: %s", selector, helpers.TruncateStr(text, 50), submitInfo), true
}

// browserSessionEvaluate runs JS in the existing session.
func browserSessionEvaluate(code, url string, chatID int64) (string, bool) {
	sess := GetOrCreateBrowserSession(chatID)

	if url != "" && url != sess.CurrentURL {
		browserSessionNavigate(url, chatID)
	}

	var result interface{}
	err := chromedp.Run(sess.Ctx,
		chromedp.Evaluate(code, &result),
	)
	if err != nil {
		return fmt.Sprintf("Evaluate error: %v", err), false
	}

	return helpers.TruncOutput(fmt.Sprintf("JS result: %v", result), helpers.MaxToolOutput), true
}

// browserSessionScroll scrolls in the existing session.
func browserSessionScroll(direction, url string, chatID int64) (string, bool) {
	sess := GetOrCreateBrowserSession(chatID)

	if url != "" && url != sess.CurrentURL {
		browserSessionNavigate(url, chatID)
	}

	scrollScript := "window.scrollBy(0, 500)"
	if direction == "up" {
		scrollScript = "window.scrollBy(0, -500)"
	}

	err := chromedp.Run(sess.Ctx,
		chromedp.Evaluate(scrollScript, nil),
		chromedp.Sleep(1*time.Second),
	)
	if err != nil {
		return fmt.Sprintf("Scroll error: %v", err), false
	}

	return fmt.Sprintf("📜 Scrolled %s", direction), true
}

// browserSessionExtract extracts text from a selector in the existing session.
func browserSessionExtract(selector, url string, chatID int64) (string, bool) {
	sess := GetOrCreateBrowserSession(chatID)

	if url != "" && url != sess.CurrentURL {
		browserSessionNavigate(url, chatID)
	}

	var text string
	err := chromedp.Run(sess.Ctx,
		chromedp.Text(selector, &text, chromedp.ByQuery),
	)
	if err != nil {
		return fmt.Sprintf("Extract error: %v", err), false
	}

	text = whitespaceRe.ReplaceAllString(text, " ")
	text = strings.TrimSpace(text)
	return helpers.TruncOutput(text, helpers.MaxToolOutput), true
}

// browserSnapshot extracts a compact accessibility tree of interactive elements.
func browserSnapshot(chatID int64) (string, bool) {
	sess := GetOrCreateBrowserSession(chatID)

	// Extract interactive elements with their attributes
	var elements string
	snapScript := `
		(function() {
			var interactive = document.querySelectorAll('a, button, input, select, textarea, [role="button"], [onclick]');
			var lines = [];
			interactive.forEach(function(el, i) {
				var tag = el.tagName.toLowerCase();
				var type = el.type || '';
				var text = (el.innerText || el.value || el.placeholder || el.getAttribute('aria-label') || '').trim();
				if (text.length > 60) text = text.substring(0, 57) + '...';
				var id = el.id ? '#' + el.id : '';
				var cls = el.className ? '.' + (typeof el.className === 'string' ? el.className.split(' ')[0] : '') : '';
				var name = el.name ? '[name=' + el.name + ']' : '';
				var href = el.href ? ' → ' + el.href.substring(0, 60) : '';
				lines.push('[' + i + '] <' + tag + (type ? ' type=' + type : '') + id + cls + name + '>: ' + text + href);
			});
			return lines.join('\n') || 'No interactive elements found';
		})()
	`
	err := chromedp.Run(sess.Ctx,
		chromedp.Evaluate(snapScript, &elements),
	)
	if err != nil {
		return fmt.Sprintf("Snapshot error: %v", err), false
	}

	var title, url string
	chromedp.Run(sess.Ctx,
		chromedp.Title(&title),
		chromedp.Location(&url),
	)

	return helpers.TruncOutput(fmt.Sprintf("📸 <b>Snapshot</b>\nURL: %s\nTitle: %s\n\n%s", url, title, elements), helpers.MaxToolOutput), true
}

// browserConsole captures console messages and JS errors.
func browserConsole(chatID int64) (string, bool) {
	sess := GetOrCreateBrowserSession(chatID)

	// Inject console capture and collect errors
	var result string
	consoleScript := `
		(function() {
			var msgs = window.__scorpConsole || [];
			var errors = window.__scorpErrors || [];
			var lines = [];
			if (msgs.length > 0) {
				lines.push('Console:');
				msgs.slice(-20).forEach(function(m) { lines.push('  ' + m); });
			}
			if (errors.length > 0) {
				lines.push('\nErrors:');
				errors.slice(-10).forEach(function(e) { lines.push('  ❌ ' + e); });
			}
			if (lines.length === 0) lines.push('(no console output captured)');
			return lines.join('\n');
		})()
	`

	// First inject the console/error capture
	injectScript := `
		(function() {
			if (window.__scorpConsoleSet) return;
			window.__scorpConsoleSet = true;
			window.__scorpConsole = [];
			window.__scorpErrors = [];
			var origLog = console.log;
			var origWarn = console.warn;
			var origError = console.error;
			console.log = function() {
				var args = Array.from(arguments).map(function(a) { return typeof a === 'object' ? JSON.stringify(a) : String(a); }).join(' ');
				window.__scorpConsole.push('[LOG] ' + args);
				origLog.apply(console, arguments);
			};
			console.warn = function() {
				var args = Array.from(arguments).map(function(a) { return typeof a === 'object' ? JSON.stringify(a) : String(a); }).join(' ');
				window.__scorpConsole.push('[WARN] ' + args);
				origWarn.apply(console, arguments);
			};
			console.error = function() {
				var args = Array.from(arguments).map(function(a) { return typeof a === 'object' ? JSON.stringify(a) : String(a); }).join(' ');
				window.__scorpConsole.push('[ERROR] ' + args);
				origError.apply(console, arguments);
			};
			window.addEventListener('error', function(e) {
				window.__scorpErrors.push(e.message + ' at ' + e.filename + ':' + e.lineno);
			});
		})()
	`

	chromedp.Run(sess.Ctx, chromedp.Evaluate(injectScript, nil))
	chromedp.Run(sess.Ctx, chromedp.Evaluate(consoleScript, &result))

	return helpers.TruncOutput("🖥 <b>Console Output</b>\n\n"+result, helpers.MaxToolOutput), true
}
