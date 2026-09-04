//go:build !nobrowser
// +build !nobrowser

package tools

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"scorp-agent/browser"
	"scorp-agent/internal/helpers"

	"github.com/chromedp/chromedp"
)

// ──────────────────────────────────────────────
// Browser AutoLogin Tool (requires chromedp)
// ──────────────────────────────────────────────

// ExecuteAutoLogin detects login forms and auto-fills credentials from Vault.
func ExecuteAutoLogin(args map[string]interface{}, chatID int64) (string, bool) {
	domain := helpers.GetStringArg(args, "domain", "")
	if domain == "" && len(args) > 0 {
		rawURL := helpers.GetStringArg(args, "url", "")
		if rawURL == "" {
			return "Error: domain or url required", false
		}
		domain = extractDomain(rawURL)
	}

	username, password, ok := Vault.Get(domain)
	if !ok {
		return fmt.Sprintf("No stored credentials for %s. Use vault_set to add them.", domain), false
	}

	sess := browser.GetOrCreateBrowserSession(chatID)

	log.Printf("[autologin] Attempting login for %s with user %s", domain, username)

	// Check if we're on a login page (look for password field first)
	var detectScript = `
		(function() {
			var pwdInput = document.querySelector('input[type="password"]');
			return JSON.stringify({hasPassword: !!pwdInput});
		})();
	`

	var detectJSON string
	if err := chromedp.Run(sess.Ctx, chromedp.Evaluate(detectScript, &detectJSON)); err != nil {
		return fmt.Sprintf("Error detecting login form: %v", err), false
	}

	var detectResult struct {
		HasPassword bool `json:"hasPassword"`
	}
	json.Unmarshal([]byte(detectJSON), &detectResult)

	if !detectResult.HasPassword {
		return fmt.Sprintf("No login form detected on current page. Navigate to login page first."), false
	}

	// Try filling username
	selectors := []string{
		"input[name='username']",
		"input[name='email']",
		"input[name='login']",
		"input[name='user']",
		"input[type='email']",
		"#username",
		"#email",
		"#user",
		"#login",
	}

	filledUser := false
	for _, sel := range selectors {
		err := chromedp.Run(sess.Ctx,
			chromedp.WaitVisible(sel, chromedp.ByQuery),
			chromedp.SendKeys(sel, username, chromedp.ByQuery),
		)
		if err == nil {
			filledUser = true
			log.Printf("[autologin] Filled username using %s", sel)
			break
		}
	}

	// Fill password
	pwdSelectors := []string{
		"input[type='password']",
		"input[name='password']",
		"#password",
	}

	filledPwd := false
	for _, sel := range pwdSelectors {
		err := chromedp.Run(sess.Ctx,
			chromedp.WaitVisible(sel, chromedp.ByQuery),
			chromedp.SendKeys(sel, password, chromedp.ByQuery),
		)
		if err == nil {
			filledPwd = true
			log.Printf("[autologin] Filled password using %s", sel)
			break
		}
	}

	if !filledUser && !filledPwd {
		return fmt.Sprintf("Could not find username or password fields for %s", domain), false
	}

	// Try clicking submit button
	submitSelectors := []string{
		"button[type='submit']",
		"input[type='submit']",
		"#submit",
		"#login-button",
		"#signin",
	}

	for _, sel := range submitSelectors {
		err := chromedp.Run(sess.Ctx, chromedp.Click(sel, chromedp.ByQuery))
		if err == nil {
			log.Printf("[autologin] Clicked submit button: %s", sel)
			time.Sleep(2 * time.Second) // wait for navigation
			break
		}
	}

	return fmt.Sprintf("✅ Auto-filled credentials for %s (user: %s). Submitted login form.", domain, username), true
}

func extractDomain(rawURL string) string {
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	host := u.Hostname()
	return host
}
