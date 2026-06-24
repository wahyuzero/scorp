package bootstrap

import (
	"scorp-agent/tools"
	"scorp-agent/registry"
)


// init() runs before main() — register all tools
func init() {
	// ── Shell ──
	registry.RegisterTool(registry.ToolDef{
		Name:        "shell",
		Description: "Jalankan shell command. Mengembalikan stdout+stderr. Gunakan timeout untuk command lama.",
		Category:    "shell",
		Native:      true,
		Execute: tools.ExecuteShell,
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
		Execute: tools.ExecuteSendFile,
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
			"path":    {Type: "string", Description: "File path", Required: true},
			"max_len": {Type: "integer", Description: "Max chars to return", Default: 50000},
		},
	})
	registry.RegisterTool(registry.ToolDef{
		Name:        "write_file",
		Description: "Write content to a file (overwrites)",
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