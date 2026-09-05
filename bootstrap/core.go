package bootstrap

import (
	"scorp-agent/registry"
	"scorp-agent/tools"
)

// init() runs before main() — register all tools
func init() {
	// ── Shell ──
	registry.RegisterTool(registry.ToolDef{
		Name:        "shell",
		Description: "Run a shell command. Returns stdout+stderr. Use timeout for long-running commands.",
		Category:    "shell",
		Native:      true,
		Execute:     tools.ExecuteShell,
		Arguments: map[string]registry.ArgDef{
			"command": {Type: "string", Description: "The shell command to execute", Required: true},
			"timeout": {Type: "integer", Description: "Timeout in seconds", Default: 30},
		},
	})
	registry.RegisterTool(registry.ToolDef{
		Name:        "send_file",
		Description: "Send a file to the user via Telegram",
		Category:    "other",
		Native:      false,
		Execute:     tools.ExecuteSendFile,
		Arguments: map[string]registry.ArgDef{
			"path": {Type: "string", Description: "File path", Required: true},
		},
	})

	// ── System ──
	registry.RegisterTool(registry.ToolDef{
		Name:        "system_info",
		Description: "Get system information: cpu, memory, disk, network, docker, services, or full",
		Category:    "system",
		Native:      true,
		Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
			return tools.ExecuteSystemInfo(args)
		},
		Arguments: map[string]registry.ArgDef{
			"type": {Type: "string", Description: "Type: full, cpu, memory, disk, network, docker, services", Default: "full"},
		},
	})

	// ── File ops ──
	registry.RegisterTool(registry.ToolDef{
		Name:        "read_file",
		Description: "Read a file's content",
		Category:    "shell",
		Native:      true,
		Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
			return tools.ExecuteReadFile(args)
		},
		Arguments: map[string]registry.ArgDef{
			"path":   {Type: "string", Description: "File path", Required: true},
			"offset": {Type: "integer", Description: "1-based line number to start reading from (optional)", Default: 1},
			"limit":  {Type: "integer", Description: "Maximum number of lines to read (optional)", Default: 100},
			"lines":  {Type: "integer", Description: "Alias for limit (optional)", Default: 100},
		},
	})
	registry.RegisterTool(registry.ToolDef{
		Name:        "write_file",
		Description: "Write content to a file (overwrites). Only use for brand new files. For editing existing files, use replace_file_content.",
		Category:    "shell",
		Native:      true,
		Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
			return tools.ExecuteWriteFile(args)
		},
		Arguments: map[string]registry.ArgDef{
			"path":    {Type: "string", Description: "File path", Required: true},
			"content": {Type: "string", Description: "Content to write", Required: true},
		},
	})
	// ── Surgical Diff / Chunk Replacement ──
	registry.RegisterTool(registry.ToolDef{
		Name:        "replace_file_content",
		Description: "Surgically replace exact or fuzzy content chunk in a file. PREFERRED over write_file when editing existing files to avoid rewriting the entire file and saving tokens.",
		Category:    "shell",
		Native:      true,
		Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
			return tools.ExecuteReplaceFileContent(args)
		},
		Arguments: map[string]registry.ArgDef{
			"target_file":         {Type: "string", Description: "File path to modify", Required: true},
			"target_content":      {Type: "string", Description: "Exact existing content chunk to replace", Required: true},
			"replacement_content": {Type: "string", Description: "New replacement content chunk", Required: true},
			"start_line":          {Type: "integer", Description: "Optional starting line number constraint"},
			"end_line":            {Type: "integer", Description: "Optional ending line number constraint"},
			"replace_all":         {Type: "boolean", Description: "Replace all occurrences instead of first/unique", Default: false},
		},
	})
	registry.RegisterTool(registry.ToolDef{
		Name:        "patch",
		Description: "Alias for surgical replacement (replace_file_content). Use path, old_string, new_string.",
		Category:    "shell",
		Native:      true,
		Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
			return tools.ExecutePatch(args)
		},
		Arguments: map[string]registry.ArgDef{
			"path":        {Type: "string", Description: "File path to patch", Required: true},
			"old_string":  {Type: "string", Description: "Exact text or chunk to find", Required: true},
			"new_string":  {Type: "string", Description: "New replacement text", Required: true},
			"replace_all": {Type: "boolean", Description: "Replace all occurrences", Default: false},
		},
	})
	registry.RegisterTool(registry.ToolDef{
		Name:        "list_dir",
		Description: "List directory contents",
		Category:    "shell",
		Native:      true,
		Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
			return tools.ExecuteListDir(args)
		},
		Arguments: map[string]registry.ArgDef{
			"path": {Type: "string", Description: "Directory path", Default: "."},
		},
	})

	// ── Code search ──
	registry.RegisterTool(registry.ToolDef{
		Name:        "search_code",
		Description: "Search codebase using ripgrep (regex)",
		Category:    "code",
		Native:      true,
		Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
			return tools.ExecuteSearchCode(args)
		},
		Arguments: map[string]registry.ArgDef{
			"pattern": {Type: "string", Description: "Regex pattern", Required: true},
			"path":    {Type: "string", Description: "Search path", Default: "."},
		},
	})
}
