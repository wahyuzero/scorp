package bootstrap


import (
	"scorp-agent/delegate"
	"scorp-agent/mcp"
	"scorp-agent/skills"
	"scorp-agent/tools"
	"scorp-agent/registry"
	"scorp-agent/config"
	"scorp-agent/session"
	"scorp-agent/rag"
)

// init() — register remaining tools (web, git, docker, mcp, vision, browser)
func init() {
	// ── Web ──
	registry.RegisterTool(registry.ToolDef{
		Name:        "web_fetch",
		Description: "Fetch a URL and return HTML/text content",
		Category:    "network",
		Native:      true,
		Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
			return tools.ExecuteWebFetch(args)
		},
		Arguments: map[string]registry.ArgDef{
			"url": {Type: "string", Description: "URL to fetch", Required: true},
		},
	})
	registry.RegisterTool(registry.ToolDef{
		Name:        "web_search",
		Description: "Search the web and return results",
		Category:    "network",
		Native:      true,
		Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
			return tools.ExecuteWebSearch(args)
		},
		Arguments: map[string]registry.ArgDef{
			"query": {Type: "string", Description: "Search query", Required: true},
		},
	})
	registry.RegisterTool(registry.ToolDef{
		Name:        "http",
		Description: "Make HTTP request with method, headers, body",
		Category:    "network",
		Native:      true,
		Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
			return tools.ExecuteHTTP(args)
		},
		Arguments: map[string]registry.ArgDef{
			"url":     {Type: "string", Description: "URL", Required: true},
			"method":  {Type: "string", Description: "HTTP method", Default: "GET", Enum: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD"}},
			"headers": {Type: "object", Description: "Headers JSON"},
			"body":    {Type: "string", Description: "Request body"},
		},
	})

	// ── Git ──
	registry.RegisterTool(registry.ToolDef{
		Name:        "git",
		Description: "Run git commands",
		Category:    "code",
		Native:      true,
		Execute: tools.ExecuteGit,
		Arguments: map[string]registry.ArgDef{
			"command": {Type: "string", Description: "Git command (e.g. 'status', 'log')", Required: true},
			"repo":    {Type: "string", Description: "Repo path", Default: "."},
		},
	})

	// ── Logs ──
	registry.RegisterTool(registry.ToolDef{
		Name:        "log",
		Description: "Fetch logs from docker, journal, or file",
		Category:    "system",
		Native:      true,
		Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
			return tools.ExecuteLog(args)
		},
		Arguments: map[string]registry.ArgDef{
			"source": {Type: "string", Description: "docker, journal, or file", Required: true, Enum: []string{"docker", "journal", "file"}},
			"target": {Type: "string", Description: "Container/unit/file path", Required: true},
			"lines":  {Type: "integer", Description: "Lines to fetch", Default: 50},
		},
	})

	// ── SQL ──
	registry.RegisterTool(registry.ToolDef{
		Name:        "sql",
		Description: "Execute SQL query (MySQL/PostgreSQL)",
		Category:    "database",
		Native:      true,
		Execute: tools.ExecuteSQL,
		Arguments: map[string]registry.ArgDef{
			"query":    {Type: "string", Description: "SQL query", Required: true},
			"database": {Type: "string", Description: "Database name"},
		},
	})

	// ── Process ──
	registry.RegisterTool(registry.ToolDef{
		Name:        "process",
		Description: "Manage processes: list, kill, info",
		Category:    "system",
		Native:      true,
		Execute: tools.ExecuteProcess,
		Arguments: map[string]registry.ArgDef{
			"action": {Type: "string", Description: "list, kill, or info", Required: true, Enum: []string{"list", "kill", "info"}},
			"pid":    {Type: "integer", Description: "Process ID (for kill/info)"},
		},
	})

	// ── Memory ──
	registry.RegisterTool(registry.ToolDef{
		Name:        "memory",
		Description: "Store/retrieve persistent key-value memory",
		Category:    "other",
		Native:      true,
		Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
			return tools.ExecuteMemory(args)
		},
		Arguments: map[string]registry.ArgDef{
			"action": {Type: "string", Description: "set, get, delete, list", Required: true, Enum: []string{"set", "get", "delete", "list"}},
			"key":    {Type: "string", Description: "Memory key"},
			"value":  {Type: "string", Description: "Value (for set)"},
		},
	})

	// ── Vision ──
	registry.RegisterTool(registry.ToolDef{
		Name:        "analyze_image",
		Description: "Analyze an image using a vision model. Use after browser screenshots.",
		Category:    "vision",
		Native:      true,
		Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
			return tools.ExecuteAnalyzeImage(args)
		},
		Arguments: map[string]registry.ArgDef{
			"path":    {Type: "string", Description: "Image file path", Required: true},
			"question": {Type: "string", Description: "What to analyze", Default: "Describe this image"},
		},
	})

	// ── Agent (delegate, delegate_batch) ──
	registry.RegisterTool(registry.ToolDef{
		Name:        "delegate",
		Description: "Spawn ONE subagent to run a focused task with restricted tools. Each subagent gets its own model (role-based routing: auto/coding/research/cheap) and up to 20 iterations. Subagents CANNOT re-delegate.",
		Category:    "agent",
		Native:      true,
		Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
			return delegate.ExecuteDelegate(args)
		},
		Arguments: map[string]registry.ArgDef{
			"task":       {Type: "string", Description: "What the subagent should accomplish", Required: true},
			"context":    {Type: "string", Description: "Background info, constraints, or specific requirements"},
			"tools":      {Type: "array", Description: "Allowed tools (default: read-only: read_file, search_code, system_info, log, web_fetch, web_search, list_dir, index_search)"},
			"max_iters":  {Type: "integer", Description: "Max iterations (1-20, default 10)"},
			"return_raw": {Type: "boolean", Description: "Return raw tool results instead of formatted summary"},
			"role":       {Type: "string", Description: "Model routing: auto (default), coding, research, cheap", Enum: []string{"auto", "coding", "research", "cheap"}},
			"model":      {Type: "string", Description: "Explicit model name override (e.g. '9router-glm5-agent'). Takes priority over role."},
		},
	})
	registry.RegisterTool(registry.ToolDef{
		Name:        "delegate_batch",
		Description: "Spawn MULTIPLE subagents in PARALLEL (max 5 concurrent). Each task can have its own role/model. Use for independent workstreams that benefit from concurrency.",
		Category:    "agent",
		Native:      true,
		Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
			return delegate.ExecuteDelegateBatch(args)
		},
		Arguments: map[string]registry.ArgDef{
			"tasks":     {Type: "array", Description: "Array of task objects, each with: task, context, tools, role, model, max_iters", Required: true},
			"max_batch": {Type: "integer", Description: "Max concurrent subagents (1-5, default 5)"},
		},
	})

	// ── RAG / Semantic Search ──
	registry.RegisterTool(registry.ToolDef{
		Name:        "index_add",
		Description: "Index a file or directory for semantic search. Supports text files, code, configs, logs (max 50 files, 1MB each).",
		Category:    "rag",
		Native:      false,
		Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
			return rag.RagIndexAdd(args)
		},
		Arguments: map[string]registry.ArgDef{
			"path": {Type: "string", Description: "File or directory path to index", Required: true},
		},
	})
	registry.RegisterTool(registry.ToolDef{
		Name:        "index_search",
		Description: "Search indexed content using TF-IDF semantic search. Returns most relevant chunks with scores.",
		Category:    "rag",
		Native:      false,
		Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
			return rag.RagIndexSearch(args)
		},
		Arguments: map[string]registry.ArgDef{
			"query":  {Type: "string", Description: "Natural language search query", Required: true},
			"top_k":  {Type: "integer", Description: "Max results (default 5)"},
		},
	})
	registry.RegisterTool(registry.ToolDef{
		Name:        "index_list",
		Description: "List all indexed sources and their chunk counts.",
		Category:    "rag",
		Native:      false,
		Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
			return rag.RagIndexList(args)
		},
		Arguments: map[string]registry.ArgDef{},
	})
	registry.RegisterTool(registry.ToolDef{
		Name:        "index_remove",
		Description: "Remove all indexed chunks from a specific source file.",
		Category:    "rag",
		Native:      false,
		Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
			return rag.RagIndexRemove(args)
		},
		Arguments: map[string]registry.ArgDef{
			"source": {Type: "string", Description: "Source file path to remove", Required: true},
		},
	})

	// ── Vector RAG / Semantic Search (SimHash) ──
	// ── Docker Compose ──
	registry.RegisterTool(registry.ToolDef{
		Name:        "compose",
		Description: "Manage Docker Compose projects: up, down, restart, logs, ps, pull, config, validate",
		Category:    "docker",
		Native:      true,
		Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
			return tools.ExecuteCompose(args)
		},
		Arguments: map[string]registry.ArgDef{
			"action":       {Type: "string", Description: "Action: up, down, restart, logs, ps, pull, config, validate", Required: true},
			"project":      {Type: "string", Description: "Project directory (default: current)"},
			"file":         {Type: "string", Description: "Compose file path (optional)"},
			"detach":       {Type: "boolean", Description: "Run in detached mode (default: true, for up)"},
			"rebuild":      {Type: "boolean", Description: "Rebuild images (for up)"},
			"services":     {Type: "string", Description: "Service names (space-separated, for restart/logs/pull)"},
			"tail":         {Type: "integer", Description: "Lines to show (for logs, default 50)"},
			"follow":       {Type: "boolean", Description: "Follow logs (for logs)"},
			"timeout":      {Type: "integer", Description: "Stop timeout in seconds (for restart, default 10)"},
			"volumes":      {Type: "boolean", Description: "Remove named volumes (for down, requires confirm=true)"},
			"remove_orphans": {Type: "boolean", Description: "Remove orphan containers (for down, default true)"},
			"all":          {Type: "boolean", Description: "Show all containers including stopped (for ps)"},
			"confirm":      {Type: "boolean", Description: "Must be true for destructive actions (down)"},
		},
	})
	// ── MCP Management ──
	registry.RegisterTool(registry.ToolDef{
		Name:        "mcp_manage",
		Description: "Manage MCP servers: list configured servers, add new server (auto-reloads), remove server, or hot-reload all servers without restarting. For Python MCP servers needing user packages, add env: {\"PYTHONPATH\": \"" + config.PythonSitePackages() + "\"}",
		Category:    "system",
		Native:      true,
		Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
			return mcp.ExecuteMCPManage(args)
		},
		Arguments: map[string]registry.ArgDef{
			"action":  {Type: "string", Description: "Action: list, add, remove, reload", Required: true, Enum: []string{"list", "add", "remove", "reload"}},
			"name":    {Type: "string", Description: "Server name — short identifier (e.g. 'drawio', 'mermaid', 'borrowip'). Used in native tool names as mcp_{name}_{tool}."},
			"command": {Type: "string", Description: "Executable to run (for add). e.g. 'npx', 'python3', 'node', 'uvx'"},
			"args":    {Type: "array", Description: "Command arguments array (for add). e.g. [\"-y\", \"@drawio/mcp-server\"]"},
			"env":     {Type: "object", Description: "Environment variables (for add). JSON object of key-value pairs"},
		},
	})

	// ── Skill Management ──
	registry.RegisterTool(registry.ToolDef{
		Name:        "skill_manage",
		Description: "Manage agent skills: list, view, create, update, delete. Skills are stored as JSON files in ~/.scorp/skills/. Built-in skills cannot be deleted.",
		Category:    "system",
		Native:      true,
		Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
			return skills.ExecuteSkillManage(args)
		},
		Arguments: map[string]registry.ArgDef{
			"action":  {Type: "string", Description: "Action: list, view, create, update, delete", Required: true, Enum: []string{"list", "view", "create", "update", "delete"}},
			"name":    {Type: "string", Description: "Skill name (required for view, create, update, delete)"},
			"content": {Type: "string", Description: "Skill JSON content (required for create, update)"},
		},
	})

	// ── Patch (file editing) ──
	registry.RegisterTool(registry.ToolDef{
		Name:        "patch",
		Description: "Surgical file edit: find old_string and replace with new_string. Supports fuzzy matching (exact, trim whitespace, normalize whitespace). Returns unified diff preview.",
		Category:    "code",
		Native:      true,
		Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
			return tools.ExecutePatch(args)
		},
		Arguments: map[string]registry.ArgDef{
			"mode":         {Type: "string", Description: "Mode: replace (default)", Default: "replace", Enum: []string{"replace"}},
			"path":         {Type: "string", Description: "File path to edit", Required: true},
			"old_string":   {Type: "string", Description: "Text to find and replace", Required: true},
			"new_string":   {Type: "string", Description: "Replacement text", Required: true},
			"replace_all":  {Type: "boolean", Description: "Replace all occurrences (default: false)"},
		},
	})

	// ── Todo (task tracking) ──
	registry.RegisterTool(registry.ToolDef{
		Name:        "todo",
		Description: "Manage a task list for multi-step work. No args: show list. With todos array: create/replace (merge=false) or update by id (merge=true). Enforces single in_progress item.",
		Category:    "other",
		Native:      true,
		Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
			return tools.ExecuteTodo(args)
		},
		Arguments: map[string]registry.ArgDef{
			"todos": {
				Type:        "array",
				Description: "Array of {id, content, status} objects. status: pending|in_progress|completed|cancelled",
			},
			"merge": {Type: "boolean", Description: "If true, update existing items by id instead of replacing list", Default: false},
		},
	})

	// ── Background Process Management ──
	registry.RegisterTool(registry.ToolDef{
		Name:        "bg",
		Description: "Manage background processes: spawn, list, poll, wait, kill, write, submit. Use for long-running commands, servers, or interactive programs.",
		Category:    "system",
		Native:      true,
		Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
			return tools.ExecuteBgProcess(args)
		},
		Arguments: map[string]registry.ArgDef{
			"action": {
				Type:        "string",
				Description: "Action: spawn, list, poll, wait, kill, write, submit",
				Required:    true,
				Enum:        []string{"spawn", "list", "poll", "wait", "kill", "write", "submit"},
			},
			"command":    {Type: "string", Description: "Shell command to run (for spawn)"},
			"session_id": {Type: "string", Description: "Process session ID (for poll, wait, kill, write, submit)"},
			"workdir":    {Type: "string", Description: "Working directory for the process (for spawn)"},
			"timeout":    {Type: "integer", Description: "Wait timeout in seconds (for wait; default 60)"},
			"data":       {Type: "string", Description: "Data to send to stdin (for write/submit)"},
		},
	})

	// ── Session Search ──
	registry.RegisterTool(registry.ToolDef{
		Name:        "session_search",
		Description: "Search past conversation history using full-text search (SQLite FTS5). Finds messages across all chats or within current chat. Use for recalling what was discussed, decided, or done previously.",
		Category:    "other",
		Native:      true,
		Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
			return session.ExecuteSessionSearch(args)
		},
		Arguments: map[string]registry.ArgDef{
			"query": {Type: "string", Description: "Search keywords or phrase", Required: true},
			"scope": {Type: "string", Description: "Search scope: 'all' (all chats) or 'current' (current chat only)", Default: "all", Enum: []string{"all", "current"}},
			"limit": {Type: "integer", Description: "Max results (default 10, max 50)"},
		},
	})

	// ── Uptime Monitor ──
	registry.RegisterTool(registry.ToolDef{
		Name:        "uptime",
		Description: "Manage uptime/health monitoring targets: list, add, remove, check",
		Category:    "monitoring",
		Native:      true,
		Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
			return tools.ExecuteUptime(args)
		},
		Arguments: map[string]registry.ArgDef{
			"action":          {Type: "string", Description: "Action: list, add, remove, check", Required: true},
			"name":            {Type: "string", Description: "Target name (for add/remove)"},
			"url":             {Type: "string", Description: "URL to monitor (for add)"},
			"expected_status": {Type: "integer", Description: "Expected HTTP status code (default: 200)"},
			"timeout":         {Type: "integer", Description: "Check timeout in seconds (default: 10)"},
		},
	})

	// ── Tool Discovery (for deferred/on-demand tools) ──
	registry.RegisterTool(registry.ToolDef{
		Name:        "tool_search",
		Description: "Search available tools by keyword. Returns tool name, description, and arguments. Use this to discover tools not shown in your initial tool list.",
		Category:    "meta",
		Native:      true,
		Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
			return tools.ExecuteToolSearch(args, chatID)
		},
		Arguments: map[string]registry.ArgDef{
			"query": {Type: "string", Description: "Search keywords (matches tool name, description, or category)", Required: true},
			"limit": {Type: "integer", Description: "Max results (default 10, max 50)"},
		},
	})
	registry.RegisterTool(registry.ToolDef{
		Name:        "tool_call",
		Description: "Invoke a tool by name with the given arguments. Use after finding a tool via tool_search.",
		Category:    "meta",
		Native:      true,
		Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
			return tools.ExecuteToolCall(args, chatID)
		},
		Arguments: map[string]registry.ArgDef{
			"name":      {Type: "string", Description: "The tool name to invoke", Required: true},
			"arguments": {Type: "object", Description: "Arguments object for the tool"},
		},
	})
	registry.RegisterTool(registry.ToolDef{
		Name:        "tool_list",
		Description: "List all available tools grouped by category, including deferred ones.",
		Category:    "meta",
		Native:      true,
		Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
			return tools.ExecuteToolList(args, chatID)
		},
		Arguments: map[string]registry.ArgDef{
			"category": {Type: "string", Description: "Filter by category (optional)"},
		},
	})

	// ── Vector RAG ──
	registry.RegisterTool(registry.ToolDef{
		Name:        "ragvec_ingest",
		Description: "Ingest file or directory into vector index (generates embeddings, chunks, stores).",
		Category:    "rag",
		Native:      true,
		Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
			return rag.RagVecIngest(args)
		},
		Arguments: map[string]registry.ArgDef{
			"path": {Type: "string", Description: "File or directory path to ingest", Required: true},
			"chunk_size": {Type: "integer", Description: "Max chars per chunk (default 1000)"},
			"overlap": {Type: "integer", Description: "Overlap chars between chunks (default 200)"},
		},
	})
	registry.RegisterTool(registry.ToolDef{
		Name:        "ragvec_search",
		Description: "Semantic search over vector index using embeddings. Returns top-K most similar chunks.",
		Category:    "rag",
		Native:      true,
		Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
			return rag.RagVecSearch(args)
		},
		Arguments: map[string]registry.ArgDef{
			"query": {Type: "string", Description: "Search query", Required: true},
			"top_k": {Type: "integer", Description: "Max results (default 5, max 50)"},
			"hybrid": {Type: "boolean", Description: "Combine with TF-IDF scores (default true)"},
			"vector_weight": {Type: "number", Description: "Weight for vector score in hybrid (0-1, default 0.7)"},
		},
	})
	registry.RegisterTool(registry.ToolDef{
		Name:        "ragvec_list",
		Description: "List all sources in vector index with chunk counts.",
		Category:    "rag",
		Native:      true,
		Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
			return rag.RagVecList(args)
		},
		Arguments: map[string]registry.ArgDef{},
	})
	registry.RegisterTool(registry.ToolDef{
		Name:        "ragvec_remove",
		Description: "Remove all chunks from a source in vector index.",
		Category:    "rag",
		Native:      true,
		Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
			return rag.RagVecRemove(args)
		},
		Arguments: map[string]registry.ArgDef{
			"source": {Type: "string", Description: "Source path/URL to remove", Required: true},
		},
	})
	registry.RegisterTool(registry.ToolDef{
		Name:        "ragvec_provider",
		Description: "Get or set the active embedding provider (tfidf, local, 9router, openai).",
		Category:    "rag",
		Native:      true,
		Execute: func(args map[string]interface{}, chatID int64) (string, bool) {
			return rag.RagVecProvider(args)
		},
		Arguments: map[string]registry.ArgDef{
			"provider": {Type: "string", Description: "Provider to set (tfidf, local, 9router, openai). Omit to get current."},
		},
	})
}