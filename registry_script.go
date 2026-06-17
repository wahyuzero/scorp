// +build !nobrowser

package main

func init() {
	// ── Script Runner ──
	registerTool(ToolDef{
		Name:        "script",
		Description: "Run a browser automation script. Scripts are JSON arrays of steps. Each step has: action (goto, click, type, submit, wait, extract, screenshot, evaluate, scroll), url, selector, value, wait_ms, timeout, name, retry. Provide either 'source' (file path) or 'inline' (JSON string).",
		Category:    "browser",
		Native:      true,
		Execute:     executeScript,
		Arguments: map[string]ArgDef{
			"source": {Type: "string", Description: "Path to a .json script file"},
			"inline": {Type: "string", Description: "JSON string defining the script inline"},
		},
	})

	// ── Script List ──
	registerTool(ToolDef{
		Name:        "script_list",
		Description: "List saved browser automation scripts in ~/.scorp-agent/scripts/",
		Category:    "browser",
		Native:      true,
		Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
			return executeScriptList(args)
		},
	})
}
