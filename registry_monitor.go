// +build !nobrowser

package main

func init() {
	// ── Browser Monitor ──
	registerTool(ToolDef{
		Name:        "monitor",
		Description: "Scheduled browser scraping with change detection. Actions: add (create target), remove, list, check (force check). Targets can auto-capture screenshots and ingest changes to RAG.",
		Category:    "browser",
		Native:      true,
		Execute:     executeMonitor,
		Arguments: map[string]ArgDef{
			"action":       {Type: "string", Description: "add, remove, list, check", Required: true},
			"id":           {Type: "string", Description: "Target ID (for remove/check)"},
			"name":         {Type: "string", Description: "Target name (for add/remove/check)"},
			"url":          {Type: "string", Description: "URL to monitor (for add)"},
			"selector":     {Type: "string", Description: "CSS selector to extract specific content (optional)"},
			"interval_min": {Type: "integer", Description: "Check interval in minutes (default 30)"},
			"screenshot":   {Type: "boolean", Description: "Capture screenshot on each check (default false)"},
			"ingest_rag":   {Type: "boolean", Description: "Auto-ingest changes to RAG (default true)"},
		},
	})
}
