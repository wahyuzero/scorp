// +build !nobrowser

package bootstrap

import (
	"scorp-agent/registry"
	"scorp-agent/tools"
)

func init() {
	// ── Script Runner ──
	registry.RegisterTool(registry.ToolDef{
		Name:        "script",
		Description: "Run a browser automation script. Scripts are JSON arrays of steps. Each step has: action (goto, click, type, submit, wait, extract, screenshot, evaluate, scroll), url, selector, value, wait_ms, timeout, name, retry. Provide either 'source' (file path) or 'inline' (JSON string).",
		Category:    "browser",
		Native:      true,
		Execute:     tools.ExecuteScript,
		Arguments: map[string]registry.ArgDef{
			"source": {Type: "string", Description: "Path to a .json script file"},
			"inline": {Type: "string", Description: "JSON string defining the script inline"},
			},
		})
	// ── Script List ──
	registry.RegisterTool(registry.ToolDef{
		Name:        "script_list",
		Description: "List saved browser automation scripts in ~/.scorp/scripts/",
		Category:    "browser",
		Native:      true,
		Execute:     tools.ExecuteScriptList,
	})
}