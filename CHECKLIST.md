# scorp-agent Development Checklist

> Progress tracker for DEVELOPMENT_PLAN.md  
| ☐ = Not started | ☐→ = In progress | ☑ = Done

---

## Phase 0: Finance Removal

### File Deletion
- ☑ Delete `finance.go` (605 lines)
- ☑ Delete `finance.json`

### Reference Cleanup
- ☑ `main.go:99-100` — Remove `loadFinanceConfig()` call + comment
- ☑ `main.go:617-621` — Remove `/market` command handler
- ☑ `agent_prompt.go:61-64` — Remove `### finance` from system prompt
- ☑ `agent_prompt.go:208-209` — Remove `case "finance"` from `executeTool()`
- ☑ `agent_loop.go:210-212` — Remove `case "finance"` from `toolDescription()`
- ☑ `telegram.go:121` — Remove `/market` from bot command list
- ☑ `model_router.go:107` — Remove `"finance": "kimi"` routing rule

### Verification
- ☑ `grep -rn "finance\|Finance\|finCfg\|PriceAlert\|CoinGecko\|AlphaVantage" *.go` → zero results
- ☑ `go build` → compiles clean (14 MB binary)
- ☑ Restart service → no crash
- ☑ `/help` → no `/market` command
- ☑ Agent mode → system prompt has no finance tool

---

## Phase 1: New Tools

### 1.1 Code Search (`search_code`)
- ☑ Install ripgrep (`sudo apt install ripgrep`)
- ☑ Create `tools_search.go`
- ☑ Add `case "search_code"` to `executeTool()`
- ☑ Add tool description to system prompt
- ☑ Add `case "search_code"` to `toolDescription()`
- ☑ Build, deploy, test

### 1.2 Git Tool (`git`)
- ☑ Create `tools_git.go`
- ☑ Implement: status, log, diff, commit, branch, stash, pull, push
- ☑ Safety guard: push requires confirmation
- ☑ Add `case "git"` to `executeTool()`
- ☑ Add tool description to system prompt
- ☑ Add `case "git"` to `toolDescription()`
- ☑ Build, deploy, test

### 1.3 HTTP/API Testing (`http`)
- ☑ Create `tools_http.go`
- ☑ Implement: all methods (GET/POST/PUT/PATCH/DELETE)
- ☑ Implement: bearer/basic/api_key auth
- ☑ Implement: JSON body + auto-pretty-print
- ☑ Add `case "http"` to `executeTool()`
- ☑ Add tool description to system prompt
- ☑ Add `case "http"` to `toolDescription()`
- ☑ Build, deploy, test

### 1.4 Log Tail (`log`)
- ☑ Create `tools_log.go`
- ☑ Implement sources: docker, journal, file
- ☑ Implement: follow mode with auto-stop timeout
- ☑ Add `case "log"` to `executeTool()`
- ☑ Add tool description to system prompt
- ☑ Add `case "log"` to `toolDescription()`
- ☑ Build, deploy, test

### 1.5 Database Query (`sql`)
- ☑ Add Go module deps: `go-sqlite3`, `lib/pq`, `go-sql-driver/mysql`
- ☑ Create `tools_db.go`
- ☑ Implement: SELECT queries with row limit
- ☑ Implement: write protection (INSERT/UPDATE/DELETE/DDL → confirm)
- ☑ Create `db_connections.json` config support
- ☑ Add `case "sql"` to `executeTool()`
- ☑ Add tool description to system prompt
- ☑ Add `case "sql"` to `toolDescription()`
- ☑ Build, deploy, test

### 1.6 Process Manager (`process`)
- ☑ Create `tools_process.go`
- ☑ Implement: list, top, kill, service_status, service_restart
- ☑ Safety: kill/restart requires confirmation
- ☑ Add `case "process"` to `executeTool()`
- ☑ Add tool description to system prompt
- ☑ Add `case "process"` to `toolDescription()`
- ☑ Build, deploy, test

### 1.7 Native Function Calling (BONUS)
- ☑ Add `tools` parameter to API request (OpenAI-compatible)
- ☑ Add `tool_calls` parsing in API response
- ☑ Create `tools_native.go` with tool definitions
- ☑ Add 3-layer fallback parser: native → XML tags → code-block
- ☑ Update `agent_loop.go` to use `callModelWithToolsAndFallback()`
- ☑ Switch model to `kr/glm-5-agentic`
- ☑ Build, deploy, test — CONFIRMED WORKING (2 tools called natively)

---

## Phase 2: Fix "Nanggung" Features

### 2.1 Browser Screenshot → Vision
- ☑ Create `analyze_image` tool in new `tools_vision.go`
- ☑ Implement: read image file → base64 → call vision model → return text
- ☑ Update `browserScreenshot()` to save + return file path for analysis
- ☑ Add `case "analyze_image"` to `executeTool()`
- ☑ Add tool description to system prompt
- ☑ Add `case "analyze_image"` to `toolDescription()`
- ☑ Add `analyze_image` to native tool definitions in tools_native.go
- ☑ Build, deploy, test — CONFIRMED WORKING

### 2.2 MCP Server Mode
- ☑ Create `mcp_server.go` → moved to `mcp_client.go`
- ☑ Implement: JSON-RPC 2.0 server over stdio
- ☑ Implement: `initialize` handler
- ☑ Implement: `tools/list` handler
- ☑ Implement: `tools/call` handler
- ☑ Create `mcp.json` config with `mcpServerMode` section
- ☑ Expose safe subset: shell, system_info, search_code, log
- ☑ Add startup hook in `main.go` (StartMCPServerMode, StopMCPServerMode)
- ☑ Build, deploy, test — CONFIRMED RUNNING

### 2.3 Plugin / Tool Registry
- ☑ Create `registry.go` with `ToolDef`, `ArgDef`, `registerTool()`, `getTool()`
- ☑ Register all 17 tools via `registry_init.go` + `registry_init2.go`
- ☑ Auto-generate native function calling schema via `generateNativeToolsSchema()`
- ☑ Replace hardcoded `getNativeToolDefs()` with `generateNativeToolsSchema()`
- ☑ Replace `executeTool()` switch with `executeToolByName()` dispatch
- ☑ Build, deploy — 17 tools registered, service active

### 2.4 Session Management
- ☑ Add `/forget` command — full session wipe (history + modes)
- ☑ Add `/sessions` command — list all saved sessions from disk
- ☑ Update `/help` text with new commands
- ☑ Build, deploy — service active

---

## Phase 3: Agent Capability Upgrades

### 3.1 Subagent / Delegation (v2 — Upgraded)
- ☑ Create `delegate.go` with `executeDelegate()`
- ☑ Implement: parallel goroutine spawning (max 5 concurrent)
- ☑ Implement: per-subagent max 15 iterations (down from 20)
- ☑ Implement: role-based model routing (cheap/coding/research/auto)
- ☑ Implement: read-only tool restriction for subagents
- ☑ Implement: result aggregation
- ☑ Add tool to registry + system prompt
- ☑ Build, deploy — delegate registered
- ✅ **VERIFIED WORKING** — tested 2026-06-15: subagent spawned, executed shell(df -h), returned result; 2nd subagent scanned /home/ubuntu with 3 parallel shell calls
- ✅ **v2 UPGRADE** — OhMyOpenAgent-inspired patterns:
  - Category-based model routing: role=auto/coding/research/cheap
  - Per-agent model override: explicit `model` param
  - `delegate_batch` tool for parallel execution (max 5 concurrent)
  - No-re-delegation enforced at BOTH prompt + execution level
  - Role-specific system prompts (coding/research/cheap guidance)
  - Routing rules: coding→glm-5-agentic, research→glm-5, cheap→glm-5
  - 22 tools registered total

### 3.2 RAG / Semantic Search
- ✅ **Approach**: Pure Go TF-IDF (zero deps, no API needed)
- ☑ Create `rag.go` — TF-IDF in-memory index with disk persistence
- ☑ Implement: `index_add` (add file/dir to index, auto-chunk 2000 chars)
- ☑ Implement: `index_search` (TF-IDF cosine similarity, top-K results)
- ☑ Implement: `index_list` (show indexed sources)
- ☑ Implement: `index_remove` (remove source from index)
- ☑ Add 4 tools to registry (category: rag)
- ☑ Initialize RAG in main.go (`initRAG()`)
- ☑ Build, deploy — 4 RAG tools registered, index loads from disk

### 3.3 Voice / Audio
- ✅ **STT**: faster-whisper (base model, CPU int8) via subprocess
- ✅ **TTS**: edge-tts via subprocess, default voice id-ID-ArdiNeural
- ☑ Create `voice.go` — STT + TTS + voice message handling
- ☑ Create `stt.py` — faster-whisper wrapper script
- ☑ Implement: inbound voice message → ffmpeg → STT → agent
- ☑ Implement: outbound TTS → edge-tts → ffmpeg opus → sendVoice
- ☑ Add `/voice` toggle command for voice replies
- ☑ Add Voice parsing in pollUpdates + webhook handler
- ☑ Add `handleVoiceMessage()` in main.go
- ☑ Pre-download whisper base model (verified)
- ☑ Build, deploy — STT tested ✅, TTS tested ✅

### 3.4 Webhook Mode
- ☑ Add `TELEGRAM_WEBHOOK_URL` env support in `config.go`
- ☑ Implement: webhook HTTP server (if URL set)
- ☑ Implement: fallback to long polling (if URL not set)
- ☑ Build, deploy, test

---

## Phase 4: Advanced / Nice to Have

### 4.1 Prometheus Metrics
- ☑ Create `metrics.go`
- ☑ Implement counters: agent_iterations, tool_calls, messages
- ☑ Implement gauges: sessions, memory_items, scheduler_tasks
- ☑ Implement histograms: response_time, tool_execution_time
- ☑ Add `/metrics` HTTP endpoint on `127.0.0.1:9091`
- ☑ Build, deploy, test with `curl localhost:9091/metrics` ✅ VERIFIED

### 4.2 Docker Compose Tool
- ☑ Create `tools_compose.go`
- ☑ Implement: up, down, restart, logs, ps, pull, config, validate
- ☑ Add to registry + native tool schema
- ☑ Build, deploy, test ✅ VERIFIED

### 4.3 Backup Tool
- ☐ Create `tools_backup.go`
- ☐ Implement: file backup (tar+gzip)
- ☐ Implement: database dump (mysqldump/pg_dump)
- ☐ Implement: docker volume backup
- ☐ Implement: rclone upload to GDrive/R3
- ☐ Add to registry + system prompt
- ☐ Build, deploy, test

### 4.4 Uptime / Health Check Monitor
- ☑ Create `tools_uptime.go`
- ☑ Implement: HTTP health check (GET request + status validation)
- ☑ Implement: concurrent checks with configurable interval (5 min default)
- ☑ Implement: auto-alert via Telegram on DOWN
- ☑ Implement tool: add/list/remove/check via agent command
- ☑ Add to registry + native tool schema
- ☑ Build, deploy, test ✅ VERIFIED

### 4.5 Telegram Inline Query
- ☑ Create `tools_inline.go`
- ☑ Implement: `InlineQuery` handler + results builder
- ☑ Implement: `answerInlineQuery` API
- ☑ Implement: `buildInlineResults` for status/docker/storage/network + safe commands
- ☑ Implement: safe read-only command execution (df, free, uptime, date)
- ☑ Add inline query to polling loop
- ☑ Add inline query to webhook handler
- ☑ Setup inline mode via `setupInlineMode()` (bot description, name, short description)
- ☑ Build, deploy, test ✅ VERIFIED

| Phase | Items | Est. Effort |
|---|---|---|
| Phase 0 | 15 | 30 min |
| Phase 1 | 42 | 6-9 hours |
| Phase 2 | 25 | 8-12 hours |
| Phase 3 | 28 | 15-22 hours | ✅ **DONE** |
| Phase 4 | 23 | 8-10 hours | ✅ **21/23 DONE** (Backup Tool pending) |
| **Total** | **133** | **37-53 hours** |
