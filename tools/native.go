package tools

import (
	"scorp-agent/registry"
)

// ──────────────────────────────────────────────
// Native Function Calling (OpenAI-compatible)
// ──────────────────────────────────────────────

// getNativeToolDefs returns tool definitions in OpenAI function calling format
func getNativeToolDefs() []registry.ToolSchema {
	return []registry.ToolSchema{
		{
			Type: "function",
			Function: registry.ToolSchemaFunc{
				Name:        "shell",
				Description: "Execute a shell command on the VPS. Use for system tasks, package management, service control, disk/memory checks, docker, etc.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"command": map[string]interface{}{
							"type":        "string",
							"description": "The shell command to execute",
						},
						"timeout": map[string]interface{}{
							"type":        "integer",
							"description": "Timeout in seconds (default 30)",
							"default":     30,
						},
					},
					"required": []string{"command"},
				},
			},
		},
		{
			Type: "function",
			Function: registry.ToolSchemaFunc{
				Name:        "read_file",
				Description: "Read the contents of a file.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"path": map[string]interface{}{
							"type":        "string",
							"description": "Path to the file",
						},
						"lines": map[string]interface{}{
							"type":        "integer",
							"description": "Max lines to read (optional)",
						},
					},
					"required": []string{"path"},
				},
			},
		},
		{
			Type: "function",
			Function: registry.ToolSchemaFunc{
				Name:        "write_file",
				Description: "Write content to a file. Creates parent directories if needed.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"path":    map[string]interface{}{"type": "string", "description": "File path"},
						"content": map[string]interface{}{"type": "string", "description": "File content"},
					},
					"required": []string{"path", "content"},
				},
			},
		},
		{
			Type: "function",
			Function: registry.ToolSchemaFunc{
				Name:        "list_dir",
				Description: "List directory contents with details.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"path":      map[string]interface{}{"type": "string", "description": "Directory path"},
						"recursive": map[string]interface{}{"type": "boolean", "description": "List recursively", "default": false},
					},
					"required": []string{"path"},
				},
			},
		},
		{
			Type: "function",
			Function: registry.ToolSchemaFunc{
				Name:        "system_info",
				Description: "Get system information: CPU, memory, disk, network, docker containers, services.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"type": map[string]interface{}{
							"type":        "string",
							"description": "Type of info: full, cpu, memory, disk, network, docker, services",
							"default":     "full",
						},
					},
				},
			},
		},
		{
			Type: "function",
			Function: registry.ToolSchemaFunc{
				Name:        "process",
				Description: "Inspect and manage processes and services: list, top, kill, service status/restart, ports.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"action": map[string]interface{}{
							"type":        "string",
							"description": "Action: list, top, kill, service_status, service_restart, service_list, ports",
						},
						"filter":  map[string]interface{}{"type": "string", "description": "Filter for list"},
						"service": map[string]interface{}{"type": "string", "description": "Service name for service_*"},
						"pid":     map[string]interface{}{"type": "string", "description": "PID for kill"},
						"sort_by": map[string]interface{}{"type": "string", "description": "Sort: mem or cpu"},
					},
					"required": []string{"action"},
				},
			},
		},
		{
			Type: "function",
			Function: registry.ToolSchemaFunc{
				Name:        "log",
				Description: "Fetch logs from docker containers, systemd journal, or files.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"source": map[string]interface{}{"type": "string", "description": "docker, journal, or file"},
						"target": map[string]interface{}{"type": "string", "description": "Container name, unit name, or file path"},
						"lines":  map[string]interface{}{"type": "integer", "description": "Number of lines (default 50)"},
					},
					"required": []string{"source", "target"},
				},
			},
		},
		{
			Type: "function",
			Function: registry.ToolSchemaFunc{
				Name:        "search_code",
				Description: "Search codebase using ripgrep (fast regex search).",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"pattern":     map[string]interface{}{"type": "string", "description": "Regex pattern"},
						"path":        map[string]interface{}{"type": "string", "description": "Search path", "default": "."},
						"glob":        map[string]interface{}{"type": "string", "description": "File glob filter"},
						"max_results": map[string]interface{}{"type": "integer", "description": "Max results", "default": 20},
					},
					"required": []string{"pattern"},
				},
			},
		},
		{
			Type: "function",
			Function: registry.ToolSchemaFunc{
				Name:        "http",
				Description: "Make HTTP requests to APIs or websites.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"method": map[string]interface{}{"type": "string", "description": "HTTP method", "default": "GET"},
						"url":    map[string]interface{}{"type": "string", "description": "URL to request"},
						"body":   map[string]interface{}{"type": "string", "description": "Request body (JSON)"},
						"headers": map[string]interface{}{"type": "string", "description": "JSON object of headers"},
						"timeout": map[string]interface{}{"type": "integer", "description": "Timeout seconds", "default": 15},
					},
					"required": []string{"url"},
				},
			},
		},
		{
			Type: "function",
			Function: registry.ToolSchemaFunc{
				Name:        "send_file",
				Description: "Send a file to the user via Telegram.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"path":    map[string]interface{}{"type": "string", "description": "File path"},
						"caption": map[string]interface{}{"type": "string", "description": "Optional caption"},
					},
					"required": []string{"path"},
				},
				},
				},
				{
				Type: "function",
				Function: registry.ToolSchemaFunc{
					Name:        "analyze_image",
					Description: "Analyze an image file using a vision model. Use after browser screenshots or for any image file analysis.",
					Parameters: map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"path": map[string]interface{}{"type": "string", "description": "Path to the image file (from browser screenshot or upload)"},
							"question": map[string]interface{}{"type": "string", "description": "What to look for in the image (default: describe in detail)"},
						},
						"required": []string{"path"},
					},
				},
				},
				}
}
