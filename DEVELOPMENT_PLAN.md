# scorp-agent Development Plan

> **Project**: vps-monitor-go (scorp-agent) — Go-based lightweight VPS monitoring + AI agent  
> **Current state**: 9,662 lines across 27 `.go` files, ~20 MB RAM  
> **Last updated**: 2026-06-15

---

## Table of Contents

1. [Phase 0: Finance Removal](#phase-0-finance-removal)
2. [Phase 1: Quick Wins — New Tools](#phase-1-quick-wins--new-tools)
3. [Phase 2: Fix "Nanggung" Features](#phase-2-fix-nanggung-features)
4. [Phase 3: Agent Capability Upgrades](#phase-3-agent-capability-upgrades)
5. [Phase 4: Advanced / Nice to Have](#phase-4-advanced--nice-to-have)
6. [Architecture Notes](#architecture-notes)

---

## Phase 0: Finance Removal

### Rationale
Finance tool (crypto/forex/TA via CoinGecko + AlphaVantage) doesn't fit scorp-agent's core purpose as a **VPS monitoring & operations agent**. Removing it simplifies the codebase, eliminates external API dependencies, and removes hardcoded API keys.

### Scope of Removal

#### Files to Delete
| File | Lines | Description |
|---|---|---|
| `finance.go` | 605 | All finance logic: CoinGecko, AlphaVantage, price alerts, formatting helpers |
| `finance.json` | 12 | Config with API keys (CoinGecko + AlphaVantage) |

#### References to Clean (6 files)

| File | Line(s) | What to Remove | Action |
|---|---|---|---|
| `main.go` | 99-100 | `loadFinanceConfig()` call + comment | Delete block |
| `main.go` | 617-621 | `/market` command handler (`quickMarketOverview()`) | Delete case block |
| `agent_prompt.go` | 61-64 | `### finance` tool description in system prompt | Delete section |
| `agent_prompt.go` | 208-209 | `case "finance": return executeFinance(tc.Args)` | Delete case |
| `agent_loop.go` | 210-212 | `case "finance":` in `toolDescription()` | Delete case |
| `telegram.go` | 121 | `{"command": "market", ...}` in bot command list | Delete entry |
| `model_router.go` | 107 | `"finance": "kimi"` routing rule | Delete map entry |

#### Helper Functions Staying vs Going
These are defined in `finance.go` but used elsewhere:

| Function | Used Outside finance.go? | Action |
|---|---|---|
| `formatNum` | ❌ No | Delete |
| `formatNumIDR` | ❌ No | Delete |
| `formatLargeNum` | ❌ No | Delete |
| `getFloat` | ❌ No | Delete |
| `getFloatArg` | ❌ No | Delete |
| `getStringSliceArg` | ❌ No | Delete |
| `truncateStr` | ✅ Yes (defined in `chat.go:627`, not finance.go) | Keep — not affected |
| `round1` / `round2` | Defined in `utils.go` | Keep — not affected |

#### Verification Steps
1. `grep -rn "finance\|Finance\|finCfg\|PriceAlert\|CoinGecko\|AlphaVantage\|crypto_prices\|forex\|market_report\|quickMarketOverview\|loadFinanceConfig\|saveFinanceConfig\|executeFinance" *.go` → should return **zero results**
2. `go build` → compiles clean
3. Restart service → no crash, no missing function errors in logs
4. Test `/help` → no `/market` command shown
5. Test agent mode → system prompt has no finance tool

#### Expected Impact
- **Code reduction**: ~605 lines removed (~6.3% of codebase)
- **Dependencies**: No Go module changes needed (CoinGecko/AlphaVantage accessed via `net/http`, no SDK)
- **Config**: Remove `finance.json` from disk
- **RAM**: Negligible change (config struct was small)

---

## Phase 1: Quick Wins — New Tools

### 1.1 `grep` / Code Search Tool
**Goal**: Fast codebase search without raw shell commands.  
**Implementation**: Shell wrapper around `ripgrep` binary.

```go
// New file: tools_search.go
// Tool name: "search_code"
// Args: pattern (string), path (string, default "."), glob (string, optional), max_results (int, default 20)
// Uses: rg (ripgrep) — must be installed or bundled
```

**Steps**:
1. `sudo apt install ripgrep` (or compile from source)
2. Create `tools_search.go` with `executeSearchCode()`
3. Add `case "search_code"` to `executeTool()` in `agent_prompt.go`
4. Add tool description to system prompt in `agent_prompt.go`
5. Add `case "search_code"` to `toolDescription()` in `agent_loop.go`
6. Add to allowed skills if needed

**Estimated effort**: 30 minutes  
**RAM impact**: +0 MB (uses external binary)

---

### 1.2 Git Tool
**Goal**: Structured git operations with safety guards (no force push without confirmation).  
**Implementation**: Shell wrapper with structured output parsing.

```go
// New file: tools_git.go
// Tool name: "git"
// Args: action (string: "status", "log", "diff", "commit", "branch", "stash", "pull", "push"), 
//       repo (string, path), message (string, for commit), count (int, for log)
```

**Actions**:
- `status` → `git status --porcelain` → parse to human-readable
- `log` → `git log --oneline -N` → return formatted
- `diff` → `git diff` or `git diff --staged`
- `commit` → `git add -A && git commit -m "message"` (with message validation)
- `branch` → list/create/switch branches
- `stash` → stash/pop
- `pull` / `push` → with safety check (push requires confirmation)

**Estimated effort**: 1-2 hours  
**RAM impact**: +0 MB

---

### 1.3 HTTP/API Testing Tool
**Goal**: Full HTTP client for API testing (method, headers, body, auth, response parsing).  
**Implementation**: Native Go `net/http` with structured request builder.

```go
// New file: tools_http.go
// Tool name: "http"
// Args: method (string), url (string), headers (map), body (string), 
//       auth_type (string: "bearer", "basic", "api_key"), auth_value (string),
//       timeout (int, seconds), follow_redirects (bool)
```

**Features**:
- All HTTP methods (GET, POST, PUT, PATCH, DELETE)
- JSON body auto-formatting
- Bearer/Basic/API key auth templates
- Response: status, headers, body (truncated), timing
- Auto-pretty-print JSON responses

**Estimated effort**: 1-2 hours  
**RAM impact**: +0 MB

---

### 1.4 Log Tail / Follow Tool
**Goal**: Realtime log streaming from docker containers, journalctl, or files.  
**Implementation**: Goroutine-based tail with timeout.

```go
// New file: tools_log.go
// Tool name: "log"
// Args: source (string: "docker", "journal", "file"), 
//       target (string: container name / unit name / file path),
//       lines (int, default 50), follow (bool, default false), duration (int, seconds, default 10)
```

**Sources**:
- `docker` → `docker logs --tail N [-f] CONTAINER`
- `journal` → `journalctl -u UNIT --since "X min ago"`
- `file` → `tail -n L [-f] /path/to/file`

**Safety**: Follow mode auto-stops after `duration` seconds (default 10s).  
**Estimated effort**: 1 hour  
**RAM impact**: +0 MB

---

### 1.5 Database Query Tool
**Goal**: Direct SQL queries to SQLite/PostgreSQL/MySQL without raw shell.  
**Implementation**: Go `database/sql` with driver imports.

```go
// New file: tools_db.go
// Tool name: "sql"
// Args: db_type (string: "sqlite", "postgres", "mysql"), 
//       connection (string: file path or DSN), query (string), 
//       limit (int, default 100)
```

**Safety**:
- SELECT only by default; INSERT/UPDATE/DELETE/DDL requires confirmation
- Row limit enforced (default 100, max 1000)
- Timeout: 10 seconds
- Connection string from `.scorp-agent/db_connections.json` (avoid DSN in prompt)

**Drivers**: `github.com/mattn/go-sqlite3`, `github.com/lib/pq`, `github.com/go-sql-driver/mysql`  
**Estimated effort**: 2-3 hours  
**RAM impact**: +10 MB (driver overhead)

---

### 1.6 Process Manager Tool
**Goal**: Structured process inspection and management.  
**Implementation**: Shell wrapper around `ps`, `top`, `systemctl`.

```go
// New file: tools_process.go
// Tool name: "process"
// Args: action (string: "list", "top", "kill", "service_status", "service_restart"),
//       filter (string: process name filter), pid (int), service (string)
```

**Actions**:
- `list` → `ps aux` with filtering + formatted output
- `top` → top 10 by CPU or memory
- `kill` → kill by PID (confirmation required)
- `service_status` → `systemctl status SERVICE`
- `service_restart` → `systemctl restart SERVICE` (confirmation required)

**Estimated effort**: 1 hour  
**RAM impact**: +0 MB

---

## Phase 2: Fix "Nanggung" Features

### 2.1 Browser Screenshot Vision Integration
**Current state**: `browserScreenshot()` takes screenshot and sends to Telegram, but does **NOT** feed it to a vision model for analysis. The agent can't "see" what the browser captured.

**Fix**: After taking screenshot, encode as base64 and include in tool result as a vision message. The agent loop already supports vision (see `handleUploadInAgentMode`).

```go
// In browserScreenshot():
// After capturing screenshot bytes:
// 1. Save to disk (existing behavior)
// 2. Send to Telegram (existing behavior)  
// 3. NEW: Return base64-encoded image in tool result for agent vision analysis
```

**Challenge**: Current tool results are strings. Need to support image content parts in tool results. Options:
- **Option A**: Return a file path reference, agent can use a new `analyze_image` tool
- **Option B**: Modify agent loop to handle image+text tool results (more complex)
- **Recommendation**: Option A — simpler, less invasive

**New tool**: `analyze_image`
```go
// Tool name: "analyze_image"
// Args: path (string, image file path), question (string, what to look for)
// Reads file → base64 → call vision-capable model → return text
```

**Estimated effort**: 2 hours  
**RAM impact**: +0 MB (model call is external)

---

### 2.2 MCP Server Mode
**Current state**: scorp-agent is MCP **client** only (can call external MCP server tools). Cannot **expose** its own tools to other MCP clients (Claude Desktop, Cursor, VSCode).

**Fix**: Add `mcp_server.go` implementing JSON-RPC 2.0 server over stdio, exposing scorp tools.

```go
// New file: mcp_server.go
// Implements: initialize, tools/list, tools/call
// Exposes: shell (safe subset), system_info, docker_status, log_tail, search_code
// Config: ~/.scorp-agent/mcp_server.json (enable, exposed tools, auth)
```

**Protocol**:
- stdio JSON-RPC 2.0 (standard MCP transport)
- `initialize` → return server info + capabilities
- `tools/list` → return tool schemas
- `tools/call` → execute and return result

**Safety**: Only expose read-only/safe tools by default. Shell access requires explicit config.  
**Estimated effort**: 3-4 hours  
**RAM impact**: +5 MB

---

### 2.3 Plugin / Dynamic Tool System
**Current state**: Tools are hardcoded in `executeTool()` switch statement and system prompt. Adding a new tool requires editing 3 files (executeTool, system prompt, toolDescription) + recompiling.

**Fix**: Tool registry pattern with auto-discovery.

```go
// New file: registry.go

type ToolDef struct {
    Name        string
    Description string
    Args        []ArgDef
    Executor    func(args map[string]interface{}, chatID int64) (string, bool)
    Icon        string // for toolDescription display
}

type ArgDef struct {
    Name        string
    Type        string // "string", "int", "bool", "map"
    Required    bool
    Default     interface{}
    Description string
}

var toolRegistry = []ToolDef{}

func registerTool(def ToolDef) {
    toolRegistry = append(toolRegistry, def)
}

func init() {
    registerTool(ToolDef{Name: "shell", ...})
    registerTool(ToolDef{Name: "search_code", ...})
    // etc.
}
```

**System prompt auto-generation**: Loop over registry to build tool descriptions.  
**executeTool dispatch**: Map lookup instead of switch.  
**Plugin loading**: Scan `~/.scorp-agent/plugins/*.go` or support compiled `.so` plugins (Go `plugin` package).

**Estimated effort**: 3-4 hours (registry + refactor existing tools)  
**RAM impact**: +2 MB

---

### 2.4 Session History Auto-Summary Improvements
**Current state**: History is persisted to disk (`~/.scorp-agent/sessions/`) and old messages are summarized via LLM. This already works well.

**Remaining gaps**:
- No way to list/resume old sessions
- No `/history` command to see past conversations
- Session files may accumulate without cleanup

**Fix**:
```go
// New commands:
// /sessions — list saved sessions with timestamps
// /session <id> — switch to a specific session
// /forget — clear current session history

// Auto-cleanup: sessions older than 30 days → archive or delete
```

**Estimated effort**: 1-2 hours  
**RAM impact**: +0 MB

---

### 2.5 Multi-User / Per-Chat Config
**Current state**: All chats share the same model config, API keys, and allowed paths. Bot only serves one Telegram chat ID (from `.env`).

**Remaining**: For now this is fine — single-user design. If multi-user needed later:
- Per-chat model preference overrides
- Per-chat allowed paths
- Per-chat memory namespace

**Status**: Not needed now. Document as future consideration.

---

## Phase 3: Agent Capability Upgrades

### 3.1 Subagent / Task Delegation
**Goal**: Agent can spawn parallel child tasks (e.g., "check all container logs" → spawn 5 parallel log checks).

**Implementation**:
```go
// New tool: "delegate"
// Args: tasks ([]struct{ goal, tool_subset })
// Each task runs in a separate goroutine with its own mini agent loop
// Results aggregated and returned

func executeDelegate(args map[string]interface{}, chatID int64) (string, bool) {
    // Spawn goroutines for each task
    // Each runs a simplified agent loop (max 5 iterations)
    // Aggregate results
    // Return combined output
}
```

**Safety**: Max 3 concurrent subagents. Each has 5 iteration limit. Read-only tools only.  
**Estimated effort**: 4-5 hours  
**RAM impact**: +10-20 MB (goroutine overhead)

---

### 3.2 RAG / Semantic Search
**Goal**: Index docs, logs, and code → semantic search ("find that error from last week about nginx").

**Implementation**:
```go
// New file: rag.go
// Embedding model: local (ONNX/llama.cpp) or API (OpenAI embeddings)
// Vector store: SQLite with sqlite-vss extension, or in-memory HNSW
// Index sources: ~/.scorp-agent/docs/, /var/log/, project READMEs

// Tools:
// "index_add" — add file/directory to index
// "index_search" — semantic search across indexed content  
// "index_list" — show indexed sources
```

**Options**:
- **Lightweight**: Use API-based embeddings (OpenAI/Gemini) + in-memory cosine similarity
- **Self-contained**: Local ONNX model (all-MiniLM-L6-v2, ~22 MB)

**Trade-off**: Local model adds ~50-100 MB RAM. API model adds latency + cost.  
**Estimated effort**: 5-8 hours  
**RAM impact**: +50-200 MB (depends on approach)

---

### 3.3 Voice / Audio Support
**Goal**: Accept voice messages (STT) and respond with voice (TTS).

**Implementation**:
```go
// Inbound: Telegram voice message → download .ogg → whisper.cpp STT → text → agent
// Outbound: Agent response → edge-tts or Piper TTS → .ogg → send as voice message

// Requires: whisper.cpp (or API), Piper/edge-tts
// Trigger: /voice command to toggle voice replies
```

**Estimated effort**: 4-6 hours  
**RAM impact**: +100-300 MB (if local Whisper), +0 MB (if API-based)

---

### 3.4 Webhook Mode (Replace Long Polling)
**Goal**: Better scalability and lower latency.

**Current**: Long polling with 35s timeout. Works fine for single-user but wastes resources.

**Fix**: Option to use Telegram webhook mode. Requires HTTPS endpoint (already have via Cloudflare).

```go
// Config: TELEGRAM_WEBHOOK_URL env var
// If set: register webhook, start HTTP server on configurable port
// If not set: fallback to long polling (current behavior)
```

**Estimated effort**: 2-3 hours  
**RAM impact**: +0 MB

---

## Phase 4: Advanced / Nice to Have

### 4.1 Metrics Exporter (Prometheus)
**Goal**: Expose `/metrics` endpoint for Grafana dashboards.

```go
// New file: metrics.go
// GET /metrics → Prometheus format
// Counters: agent_iterations_total, tool_calls_total, messages_processed_total
// Gauges: session_count, memory_items, scheduler_tasks
// Histograms: agent_response_seconds, tool_execution_seconds
```

**Estimated effort**: 2 hours  
**RAM impact**: +5 MB

---

### 4.2 Docker Compose Management
**Goal**: Structured compose operations (up, down, restart, logs, config validate).

```go
// Tool name: "compose"
// Args: action, project_dir, service, options
```

**Estimated effort**: 1-2 hours  
**RAM impact**: +0 MB

---

### 4.3 Backup Tool Integration
**Goal**: Structured backup operations (rclone, tar, database dumps).

```go
// Tool name: "backup"
// Args: source, destination, type ("file", "db", "docker_volume"), compress (bool)
// Uses existing rclone config + GDrive mount
```

**Estimated effort**: 2 hours  
**RAM impact**: +0 MB

---

### 4.4 Uptime/Health Check Monitor
**Goal**: Periodic HTTP/TCP health checks for external services.

```go
// Config: ~/.scorp-agent/uptime.json
// Monitors: [{ name, url, method, expected_status, interval }]
// Alert via Telegram on status change
```

**Estimated effort**: 2-3 hours  
**RAM impact**: +5 MB

---

### 4.5 Telegram Inline Query Support
**Goal**: Type `@bot query` in any chat to get quick answers.

**Estimated effort**: 1-2 hours  
**RAM impact**: +0 MB

---

## Architecture Notes

### Build & Deploy
```bash
# Build (ARM64, Oracle VPS)
cd /home/ubuntu/projects/vps-monitor-go
PATH=$PATH:/usr/local/go/bin go build -o scorp-agent .

# Deploy
sudo systemctl stop scorp-agent.service
sudo cp scorp-agent /usr/local/bin/scorp-agent
sudo systemctl start scorp-agent.service

# Verify
sudo systemctl status scorp-agent.service
sudo journalctl -u scorp-agent.service -n 20
```

### Adding a New Tool (Current Process)
1. Create `tools_<name>.go` with `execute<Name>(args map[string]interface{}, chatID int64) (string, bool)`
2. Add `case "<name>": return execute<Name>(tc.Args, chatID)` to `executeTool()` in `agent_prompt.go`
3. Add tool description to `getAgentSystemPrompt()` in `agent_prompt.go`
4. Add `case "<name>":` to `toolDescription()` in `agent_loop.go`
5. Build, deploy, test

> **Note**: Phase 2.3 (Plugin Registry) will consolidate steps 2-4 into auto-generation.

### File Organization
```
collector_*.go  → Monitoring data collectors (system, docker, coolify, hermes, security)
tools_*.go      → Agent tools (exec, web, memory, search, git, http, log, db, process)
agent_*.go      → Agent loop, prompt, session management
config.go       → Env-based configuration
telegram.go     → Telegram bot (polling, keyboards, file ops)
scheduler.go    → Cron/interval task engine
skills.go       → Predefined expertise prompts
mcp_*.go        → MCP client + server (future)
```

### Config Files (Runtime State)
```
~/.scorp-agent/
├── memory.json           # Persistent KV memory
├── scheduler.json        # Scheduled tasks
├── model_usage.json      # Token usage tracking
├── finance.json          # ← TO BE DELETED (Phase 0)
├── mcp.json              # MCP server configs
├── sessions/             # Chat history per chat ID
│   └── <chatID>.json
├── screenshots/          # Browser screenshots
├── plugins/              # Dynamic plugins (future)
├── db_connections.json   # SQL tool connection profiles (future)
└── uptime.json           # Health check monitors (future)
```

### RAM Budget
| Component | Current | After All Phases |
|---|---|---|
| Base (Go runtime + monitoring) | 20 MB | 20 MB |
| SQL drivers | — | +10 MB |
| MCP server | — | +5 MB |
| Subagent goroutines | — | +15 MB |
| RAG (local embeddings) | — | +100 MB |
| Metrics exporter | — | +5 MB |
| **Total** | **20 MB** | **~155 MB** |

Still **~12x lighter** than Hermes (~1.8 GB).

---

## Priority Summary

| Phase | Items | Effort | Impact |
|---|---|---|---|
| **Phase 0** | Finance removal | 30 min | Code cleanup |
| **Phase 1** | 6 new tools (grep, git, http, log, db, process) | 6-9 hours | Core VPS agent |
| **Phase 2** | Fix nanggung (vision, MCP server, registry, sessions) | 8-12 hours | Agent maturity |
| **Phase 3** | Advanced (subagent, RAG, voice, webhook) | 15-22 hours | Power agent |
| **Phase 4** | Nice to have (Prometheus, uptime, compose, backup) | 8-10 hours | Ops completeness |

**Recommended order**: Phase 0 → Phase 1 → Phase 2.1 (vision) → Phase 2.3 (registry) → Phase 2.2 (MCP server) → Phase 3+ as needed.
