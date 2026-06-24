// +build !nobrowser

package bootstrap

import (
	"scorp-agent/browser"
	"scorp-agent/registry"
)


// init() — register browser tool (only when not nobrowser)
func init() {
	// ── Browser ──
	registry.RegisterTool(registry.ToolDef{
		Name:        "browser",
		Description: "Stateful headless Chrome with session persistence. Actions: goto, screenshot, click, type, evaluate, scroll, extract, snapshot (interactive elements), console (JS errors), reset (clear session), list. Browser context persists across actions within same chat.",
		Category:    "browser",
		Native:      true,
		Execute:     browser.ExecuteBrowser,
		Arguments: map[string]registry.ArgDef{
			"action":   {Type: "string", Description: "goto, screenshot, click, type, evaluate, scroll, extract, snapshot, console, reset, list", Required: true},
			"url":      {Type: "string", Description: "URL (required for goto, optional for others — uses current page if omitted)"},
			"selector": {Type: "string", Description: "CSS selector (for click/type/extract)"},
			"text":     {Type: "string", Description: "Text to type"},
			"code":     {Type: "string", Description: "JavaScript code (for evaluate)"},
			"direction": {Type: "string", Description: "Scroll direction: up, down (default: down)"},
		},
	})
}