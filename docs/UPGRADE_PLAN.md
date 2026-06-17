# Scorp-Agent → Hermes-Level Upgrade Plan

> **Tanggal:** 16 Juni 2026
> **Author:** Wahyu × Hermes Agent (GLM-5.2)
> **Goal:** Menutup gap fitur scorp-agent vs hermes-agent, prioritaskan impact tinggi & effort rendah

---

## 1. Current State — Scorp-Agent v4

### 1.1 Arsitektur
- **Bahasa:** Go (single binary ~12MB, build tag `nobrowser` → ~9.7MB)
- **Build:** `make` / `make minimal` / `go build -ldflags="-s -w -buildid=" -trimpath`
- **Runtime:** systemd service (`scorp-agent`), runs as root
- **Platform:** Telegram bot (single chat ID)
- **Model:** Multi-provider router (kimi, groq, deepseek, gemini, openrouter, 9router)
- **MCP:** Config-driven via `~/.scorp-agent/mcp.json`, native tool registration

### 1.2 Inventory — 28 Native Tools (saat ini)

| # | Tool | File | Deskripsi |
|---|------|------|-----------|
| 1 | `shell` | tools_native.go | Execute shell command + timeout |
| 2 | `read_file` | tools_native.go | Read file contents (opt. line limit) |
| 3 | `write_file` | tools_native.go | Write/overwrite file |
| 4 | `list_dir` | tools_native.go | List directory (opt. recursive) |
| 5 | `system_info` | tools_native.go | CPU/mem/disk/net/docker/services |
| 6 | `process` | tools_native.go | Process list/kill, service status/restart |
| 7 | `log` | tools_native.go | Docker/journal/file logs |
| 8 | `search_code` | tools_native.go | Ripgrep regex search |
| 9 | `http` | tools_native.go | HTTP requests (any method) |
| 10 | `send_file` | tools_native.go | Send file via Telegram |
| 11 | `analyze_image` | tools_vision.go | Vision model (claude-sonnet-4 via 9router) |
| 12 | `browser` | browser.go | chromedp: goto/click/type/screenshot/scroll/extract/evaluate |
| 13 | `web_fetch` | tools_web.go | Fetch URL → text |
| 14 | `web_search` | tools_web.go | Search engine query |
| 15 | `git` | tools_git.go | Git operations |
| 16 | `sql` | tools_db.go | MySQL/PostgreSQL queries |
| 17 | `compose` | tools_docker.go | Docker compose management |
| 18 | `schedule` | scheduler.go | Cron task add/list/delete |
| 19 | `memory` | memory.go | KV store (set/get/delete/list) |
| 20 | `index_add` | rag.go | TF-IDF index from files/dirs |
| 21 | `index_search` | rag.go | TF-IDF cosine similarity search |
| 22 | `index_list` | rag.go | List indexed sources |
| 23 | `index_remove` | rag.go | Remove indexed source |
| 24 | `delegate` | delegate.go | Single subagent (role-based routing) |
| 25 | `delegate_batch` | delegate.go | Parallel subagents (max 5) |
| 26 | `subagent_run` | delegate_impl.go | Subagent execution engine |
| 27 | `mcp_manage` | mcp_manage.go | MCP server add/remove/reload/list |
| 28 | `mcp_borrowip_*` | (MCP) | 6 borrowip tools (auto-registered) |

### 1.3 Telegram Commands
`/start` `/status` `/report` `/containers` `/security` `/storage` `/hermes`
`/agent` `/stop` `/cron` `/model` `/usage` `/skill` `/skills` `/clear` `/forget` `/sessions` `/help`

### 1.4 Built-in Skills (8, hardcoded)
`docker` `git` `network` `coolify` `backup` `security` `disk` `performance`

### 1.5 Config Files (`~/.scorp-agent/`)
- `mcp.json` — MCP server definitions
- `models.json` — Model router config
- `scheduler.json` — Scheduled tasks
- `rag_index.json` — TF-IDF index
- `memory.json` — KV memory store

---

## 2. Gap Analysis — Scorp vs Hermes

### 2.1 GAP BESAR (Hermes punya, Scorp tidak)

| ID | Feature | Hermes Equivalent | Dampak | Effort |
|----|---------|-------------------|--------|--------|
| G01 | **`patch` tool** — fuzzy find-replace, targeted edit tanpa rewrite file | `patch` (9 strategies, V4A multi-file) | 🔥🔥🔥 | Medium |
| G02 | **`todo` tool** — task list untuk multi-step work | `todo` (merge/replace, priority order) | 🔥🔥 | Low |
| G03 | **Dynamic skills** — create/edit/delete tanpa recompile, JSON file-based | `skill_manage` (CRUD, YAML frontmatter, linked files) | 🔥🔥🔥 | Medium |
| G04 | **Memory auto-inject** — agent selalu aware tanpa explicit call | Auto-injected MEMORY.md + USER.md per turn | 🔥🔥 | Low |
| G05 | **Background process management** — spawn, poll, kill, PTY, stdin write | `process` (bg=true, notify, watch_patterns) | 🔥🔥 | Medium |
| G06 | **`execute_code`** — Python script dengan tool imports | `execute_code` (hermes_tools, batch ops, loops) | 🔥🔥🔥 | High |
| G07 | **Session search** — search isi percakapan masa lalu | `session_search` (FTS5, scroll, discovery) | 🔥 | High |
| G08 | **Context compaction** — handle long conversations | Auto-compaction + summary injection | 🔥🔥 | High |
| G09 | **Deferred tool loading** — load tools on-demand | `tool_search` + `tool_call` (149+ tools) | 🔥 | Medium |
| G10 | **`clarify` tool** — ask user when ambiguous | `clarify` (multiple choice / open-ended) | 🔥 | Low |
| G11 | **Multi-platform messaging** — Discord, Slack, Matrix, dll | `send_message` (multi-platform routing) | 🔥 | High |
| G12 | **ACP delegation** — delegate ke Claude Code, Codex, OpenCode | `delegate_task` (ACP subprocess transport) | 🔥 | High |

### 2.2 GAP SEDANG (ada di Scorp tapi jauh lebih primitif)

| ID | Feature | Scorp Sekarang | Hermes Equivalent | Gap |
|----|---------|----------------|-------------------|-----|
| S01 | **Skills** | 8 hardcoded prompt strings di Go | 50+ skills, YAML frontmatter, categories, linked files, auto-load | Tidak bisa create/edit runtime, tidak ada metadata |
| S02 | **RAG** | TF-IDF, 2K fixed chunks, no embeddings | (via skills/plugins) Embedding-based, dynamic chunking | No semantic search, no re-ranking |
| S03 | **Browser** | Stateless — setiap action = new chromedp context | Session persistence, cookies, snapshot dengan ref IDs, console capture | No login flows, no multi-step, fragile selectors |
| S04 | **Voice/TTS** | STT + TTS work tapi tidak exposed sebagai tool | `text_to_speech` tool (multi-provider) | Agent tidak bisa proactively speak |
| S05 | **Model Router** | Manual routing rules, usage tracking | Cost optimization, time-based auto-switching, model pinning per job | No cost-aware routing |
| S06 | **Scheduler** | Basic add/list/delete, results → Telegram | LLM-driven or script-only, context chaining, per-job model/profile/workdir, delivery targeting | No script-only mode, no chaining |
| S07 | **Memory** | KV store, must explicit call | Auto-injected per turn, user vs agent split | Not automatic |
| S08 | **Delegate** | In-process subagents, shared history | Fully isolated context per subagent, toolset restriction, orchestrator role | Shared context, no isolation |
| S09 | **Telegram** | Single chat ID | Multi-user, group chat, RBAC, threaded topics | Single user only |

### 2.3 FITUR NANGGUNG / Technical Debt

| ID | Issue | Detail | Fix |
|----|-------|--------|-----|
| T01 | Dead code `executeMCPTool()` | `agent_loop.go:38-75`, generic mcp_tool sudah dihapus | Hapus fungsi |
| T02 | Vision hardcoded model | `tools_vision.go:81` → `kr/claude-sonnet-4` hardcoded | Baca dari models.json atau arg |
| T03 | Skills tidak auto-load | Agent harus `/skill docker` atau ketik "pakai skill docker" | Auto-detect by context |
| T04 | Browser save path hardcoded | `browser.go:114` → `/home/ubuntu/.scorp-agent/screenshots` | Use config or $HOME |
| T05 | RAG path hardcoded root | `rag.go:46` → `/root/.scorp-agent/rag_index.json` | Use $HOME |
| T06 | Voice paths hardcoded | `voice.go:23-26` hardcoded paths | Use config |
| T07 | No graceful shutdown | MCP servers tidak di-clean-up saat service stop | Context cancel + wait |

---

## 3. Implementation Plan

### Phase 1 — Quick Wins (P0, 1-2 hari)
> Impact tinggi, effort rendah. Fundament yang dibutuhkan sebelum fitur kompleks.

#### P1.1: `patch` Tool (G01)
**Goal:** Edit file tanpa rewrite penuh, fuzzy matching untuk toleransi whitespace.

```go
// tools_patch.go
func executePatch(args map[string]interface{}) (string, bool) {
    // action: "replace" | "patch"
    // replace mode: path, old_string, new_string, replace_all
    // patch mode: V4A format (*** Begin Patch / *** End Patch)
    //
    // Fuzzy matching strategies (minimal 3):
    //   1. Exact match
    //   2. Whitespace-normalized match
    //   3. Line-by-line fuzzy (ignore leading/trailing whitespace)
    //
    // Return: unified diff of changes
}
```

**Schema:**
```json
{
  "name": "patch",
  "parameters": {
    "action": "replace|patch",
    "path": "string (required for replace)",
    "old_string": "string (required for replace)",
    "new_string": "string (required for replace)",
    "replace_all": "boolean (default false)",
    "patch": "string (required for patch mode, V4A format)"
  }
}
```

**Acceptance:**
- [x] `patch` action=replace dengan exact match works
- [x] Fuzzy match (whitespace beda) tetap match
- [x] `patch` action=patch dengan V4A format works
- [x] Return unified diff
- [x] Error jika old_string tidak unik (tanpa replace_all)

---

#### P1.2: `todo` Tool (G02)
**Goal:** Task list untuk multi-step agent work.

```
// tools_todo.go
func executeTodo(args map[string]interface{}) (string, bool) {
    // No args → return current list
    // todos: [{id, content, status}] → replace or merge
    // merge: false (default) → replace entire list
    // merge: true → update by id, add new
    //
    // Display: ordered list with status icons
    // 🔲 pending | 🔄 in_progress | ✅ completed | ❌ cancelled
}
```

**Schema:**
```json
{
  "name": "todo",
  "parameters": {
    "todos": [{
      "id": "string (required)",
      "content": "string (required)",
      "status": "pending|in_progress|completed|cancelled"
    }],
    "merge": "boolean (default false)"
  }
}
```

**State:** In-memory `[]TodoItem` per chat session, no persistence needed.

**Acceptance:**
- [x] Create list with `merge=false`
- [x] Update item status with `merge=true`
- [x] Only ONE `in_progress` at a time (enforce in code)
- [x] No args → return formatted list

---

### P1.3: Memory Auto-Inject (G04, S07)
**Goal:** Memory dump otomatis di system prompt setiap agent turn.

**Changes:**
1. `agent_prompt.go` — tambahkan section `[MEMORY]` di system prompt
2. Baca `memory.json` saat build system prompt
3. Format: `§ {key}: {value}` (compact, section-delimited)

```go
// agent_prompt.go — di buildSystemPrompt()
func buildMemorySection() string {
    data, _ := os.ReadFile(memoryFilePath)
    var mem map[string]string
    json.Unmarshal(data, &mem)
    // Format compact, inject ke system prompt
    var sb strings.Builder
    for k, v := range mem {
        sb.WriteString(fmt.Sprintf("%s: %s\n", k, v))
    }
    return sb.String()
}
```

**Acceptance:**
- [x] Memory entries muncul di system prompt agent
- [x] Agent tidak perlu explicit `memory_get` untuk facts dasar
- [x] Update memory → reflected di turn berikutnya

---

### Phase 2 — Core Improvements (P1, 3-5 hari)
> Membuat agent lebih powerful & self-managing.

#### P2.1: Dynamic Skills System (G03, S01)
**Goal:** Skills tersimpan sebagai file JSON/MD, bisa CRUD via agent.

**Design:**
```
~/.scorp-agent/skills/
├── docker.json
├── git.json
├── network.json
├── coolify.json
├── backup.json
├── security.json
├── disk.json
├── performance.json
└── (user-created).json
```

**Skill JSON format:**
```json
{
  "name": "docker",
  "emoji": "🐳",
  "description": "Manage Docker containers",
  "category": "devops",
  "prompt": "You are a Docker expert...",
  "examples": ["cek status container", "restart X"],
  "auto_load_keywords": ["docker", "container", "compose"]
}
```

**New tool: `skill_manage`**
```json
{
  "name": "skill_manage",
  "parameters": {
    "action": "create|list|view|update|delete",
    "name": "string",
    "content": "string (for create/update)"
  }
}
```

**Migration:** Import 8 hardcoded skills ke JSON files saat first run.

**Acceptance:**
- [x] Skills loaded from JSON files at startup
- [x] `skill_manage action=create` creates new skill file
- [x] `skill_manage action=update` edits existing skill
- [x] `skill_manage action=delete` removes skill file
- [x] `/skill` command still works with file-based skills
- [x] Auto-load keywords trigger skill prompt injection

---

#### P2.2: Background Process Management (G05)
**Goal:** Spawn long-running processes, poll output, send stdin.

```go
// tools_bg.go
type BGProcess struct {
    ID        string
    Cmd       *exec.Cmd
    StartTime time.Time
    Output    *bytes.Buffer
    Done      bool
    ExitCode  int
}

var bgProcesses = make(map[string]*BGProcess)
var bgProcessMu sync.RWMutex

func executeBgProcess(args map[string]interface{}) (string, bool) {
    // action: spawn|list|poll|wait|kill|write|submit
    //
    // spawn: command, workdir → returns session_id
    // list: → all bg processes with status
    // poll: session_id → new output since last poll
    // wait: session_id, timeout → block until done or timeout
    // kill: session_id → terminate
    // write: session_id, data → send raw stdin
    // submit: session_id, data → send data + Enter
}
```

**Acceptance:**
- [x] Spawn bg process returns session_id
- [x] Poll returns incremental output
- [x] Wait blocks until completion or timeout
- [x] Kill terminates process
- [x] List shows all active bg processes

---

#### P2.3: `clarify` Tool (G10)
**Goal:** Agent bisa tanya user saat stuck/ambiguous.

```go
// tools_clarify.go
func executeClarify(args map[string]interface{}) (string, bool) {
    // question: string (required)
    // choices: []string (optional, max 4)
    //
    // Behavior:
    //   - Jika ada choices: kirim inline keyboard buttons
    //   - Jika tidak: kirim pesan, tunggu reply
    //   - Block agent loop sampai user jawab
    //   - Return jawaban user ke agent
}
```

**Challenge:** Agent loop perlu support blocking wait untuk user input. Perlu channel-based mechanism.

**Acceptance:**
- [x] Multiple choice (inline keyboard) works
- [x] Open-ended (wait for text reply) works
- [x] Timeout (configurable, default 60s)
- [x] Agent receives user's answer as tool result

---

### Phase 3 — Advanced (P2, 1-2 minggu)
> Membuat agent setara Hermes untuk use cases utama.

#### P3.1: `execute_code` — Python Orchestration (G06)
**Goal:** Jalankan Python script yang bisa call agent tools.

**Design:**
- Spawn Python subprocess dengan injected `hermes_tools`-like module
- Module wraps tool functions: `shell()`, `read_file()`, `write_file()`, `search_files()`, `patch()`
- Output stdout → tool result
- 5-minute timeout, 50KB stdout cap, max 50 tool calls

**Acceptance:**
- [x] Python script bisa call `read_file()` and get content
- [x] Python script bisa call `terminal()` and get output
- [x] Batch processing (loop over files) works
- [x] Timeout enforced
- [x] Tool call count enforced

---

#### P3.2: Browser Session Persistence (S03)
**Goal:** Stateful browser — login flows, multi-step, cookies.

**Changes to `browser.go`:**
- Pool of named browser contexts (keyed by chat session)
- Context persists across actions within same session
- `/browser_reset` to clear session
- Cookie/storage state save/load

**Also add:**
- `browser_snapshot` — accessibility tree with ref IDs
- `browser_console` — capture JS errors + console.log

**Acceptance:**
- [x] Login flow: goto → fill → click → screenshot (same session)
- [x] Cookies persist across actions
- [x] Snapshot returns ref IDs
- [x] Console captures errors

---

#### P3.3: TTS as Agent Tool (S04)
**Goal:** Agent bisa proactively generate voice output.

```json
{
  "name": "speak",
  "parameters": {
    "text": "string (required)",
    "voice": "string (optional, default id-ID-ArdiNeural)"
  }
}
```

**Implementation:** Reuse existing `synthesizeSpeech()` from `voice.go`, just expose as tool.

**Acceptance:**
- [x] Agent calls `speak` → voice message sent to Telegram
- [x] Voice configurable via arg
- [x] Reuses existing edge-tts pipeline
---

### Phase 4 — Long-term (P3-P4, 2-4 minggu)
> Advanced features untuk parity penuh.

#### P4.1: Session Search (G07)
- SQLite + FTS5 untuk conversation history
- Migrasi dari in-memory session ke SQLite
- Tool: `session_search(query)` → matching messages

#### P4.2: Context Compaction (G08)
- Detect context window approaching limit
- Summarize older messages
- Inject summary as system message

#### P4.3: Deferred Tool Loading (G09)
- Split tools into "always loaded" (core ~10) vs "on-demand"
- Tool registry with metadata-only schemas (no full description)
- `tool_search(query)` → list matching tools
- `tool_call(name, args)` → execute

#### P4.4: Multi-Platform (G11)
- Abstract messaging layer
- Discord adapter, Slack adapter
- Per-platform routing

#### P4.5: Technical Debt Cleanup (T01-T07)
- [x] T01: Remove `executeMCPTool()` dead code
- [x] T02: Vision model from config
- [x] T03: Skills auto-load by keyword
- [x] T04: Browser path from config
- [x] T05: RAG path from $HOME
- [x] T06: Voice paths from config
- [x] T07: Graceful MCP shutdown

---

### Phase 5 — Semantic RAG ✅ (DONE)
> SimHash-based vector search + TF-IDF hybrid fusion. Pure Go, zero dependencies.

#### P5.1: Vector Store Setup
- [x] Evaluasi: SimHash (64-bit, Go native, no deps) vs sentence-transformers (PyTorch, heavy compile on ARM64 — **skipped**)
- [x] Schema: `VecChunk` struct (id, source, content, simhash uint64, created) → JSON persistence
- [x] Embedding: SimHash (FNV-1a, term-frequency weighted) + TF-IDF hybrid fusion
- [x] Config: `~/.scorp-agent/rag_vector.json`, auto-persist on write

#### P5.2: Ingestion Pipeline
- [x] File/dir ingestion: `SmartChunk` (paragraph → sentence → hard-split, configurable max-chunk)
- [x] Incremental: auto-remove old chunks from same source before re-index
- [x] Tools: `ragvec_ingest(path, chunk_size)`, `ragvec_list()`, `ragvec_remove(source)`

#### P5.3: Semantic Search + Re-ranking
- [x] Vector search: SimHash → Hamming distance → similarity (0.0-1.0), threshold filter (0.75)
- [x] Re-ranking: weighted fusion (simhash weight + TF-IDF weight)
- [x] Hybrid search: configurable `vector_weight` (default 0.7 simhash + 0.3 TF-IDF)
- [x] Tools: `ragvec_search(query, top_k, hybrid, vector_weight)` → ranked chunks with source + score

#### P5.4: RAG Integration
- [x] Agent auto-RAG: auto-search user query → inject top-3 chunks ke system prompt setiap turn
- [x] Citation: per-result source label + preview (300 chars) di result format
- [x] Tool: `ragvec_provider()` untuk cek status embedding provider

---

### Phase 6 — Browser Sessions (DONE ✅)
> Persistent browser contexts: cookies, login flows, multi-step automation.

#### P6.1: Session Persistence (DONE ✅)
- [x] User data dir: persistent chromedp context per chat ID (`~/.scorp-agent/browser_data/<chatID>`)
- [x] Cookie jar: automatic via Chrome UserDataDir (survives across tool calls)
- [x] Session reuse: `getOrCreateBrowserSession()` caches contexts
- [x] Auto-cleanup: stale sessions pruned every 5 min (10 min idle timeout)

#### P6.2: Login Flow Automation (DONE ✅)
- [x] Credential store: AES-256-GCM encrypted vault (`~/.scorp-agent/vault.json`)
- [x] Master key: auto-generated, stored at `~/.scorp-agent/vault.key`
- [x] Form detection: JS-based username/email/password field auto-detect
- [x] Tool: `vault` (get/set/list/remove) + `autologin` (auto-fill from vault)

#### P6.3: Multi-step Scripting (DONE ✅)
- [x] JSON script format: array of steps with action/url/selector/value/wait_ms/timeout/name/retry
- [x] Actions: goto, click, type, submit, wait, extract, screenshot, evaluate, scroll
- [x] Error handling: per-step retry count, screenshot on fail

#### P6.4: Monitoring & Scheduled Scraping (DONE ✅)
- [x] Tool: `monitor` (add/list/remove targets with URL + selector + interval)
- [x] Background loop: periodic scrape + change detection (SHA-256 diff)
- [x] Feed to RAG: ragvec_add on change detected
- [x] Goroutine-safe: loadMonitorTargets() separated from initMonitor() to prevent leaks

---

### Phase 7 — Autonomous Agent (DONE ✅)
> Self-governing VPS management: LLM-driven decisions, proactive actions, safety gates.

#### P7.1: Autonomous Loop Framework (DONE ✅)
- [x] Background goroutine: `autonomousLoop()` with configurable interval (default 10m)
- [x] Context builder: `gatherContext()` — system metrics, Docker containers, security, alerts
- [x] Decision engine: LLM prompt → JSON action plan (analysis + actions[] + notify + speak)
- [x] Action executor: `executeAutonomousAction()` — validates risk, blocks danger, runs tool
- [x] Memory: `autonomous_log.json` — persistent audit trail (max 500 entries)

#### P7.2: Proactive Capabilities (DONE ✅)
- [x] LLM can issue `exec` commands: docker restart, fail2ban, systemctl, etc.
- [x] LLM can call `ragvec_add` to ingest findings
- [x] LLM can call any registered tool via autonomous action
- [x] Self-healing: detect issues from metrics → propose fix → execute if low/medium risk

#### P7.3: Proactive Voice (DONE ✅)
- [x] Config: `voice: off|important|always`
- [x] After cycle: if speak=true → `synthesizeSpeech()` → send voice to Telegram
- [x] After cycle: if notify=true → send text summary to Telegram

#### P7.4: Safety & Observability (DONE ✅)
- [x] Approval levels: `low` (auto all), `medium` (approve high), `high` (approve all)
- [x] Kill switch: file-based `autonomous_kill` — works even if agent crashes
- [x] Dangerous command blocklist: rm -rf /, dd, mkfs, shutdown, halt, reboot, etc.
- [x] Max actions/cycle: configurable (default 5, max 10)
- [x] Tool: `autonomous` — status/enable/disable/kill/revive/run/config/log/actions
- [x] Test suite: 15 tests (all PASS)

---

### Gap Closure — Advanced Features

#### GC1: ACP Delegation (G12)
- [ ] Subprocess transport: delegate ke external CLI (Claude Code, Codex, OpenCode)
- [ ] Per-task ACP command override
- [ ] Toolset restriction per delegated task

#### GC2: Model Router Cost Optimization (S05)
- [ ] Cost-aware routing: track token cost per model, auto-select cheapest adequate
- [ ] Time-based auto-switching: peak/off-peak model rotation
- [ ] Per-job model pinning via config

#### GC3: Scheduler Improvements (S06)
- [ ] Script-only mode (no LLM, script stdout = output)
- [ ] Context chaining: inject output from previous job as context
- [ ] Per-job config: model, profile, workdir, delivery target

#### GC4: Delegate Isolation (S08)
- [ ] Fully isolated context per subagent (separate session, no shared history)
- [ ] Toolset restriction per subagent
- [ ] Orchestrator role: can spawn own workers

#### GC5: Post-Implementation Review
- [ ] Binary size check
- [ ] Startup time check
- [ ] Regression test suite
- [ ] Git commit (feat: gap-closure)
- [ ] Update MEMORY entry

---

## 4. Dependency Graph

```
P1.1 patch ──────────────────────────┐
P1.2 todo ───────────────────────────┤
P1.3 memory auto-inject ─────────────┤── Phase 1 (no deps)
                                     │
P2.1 dynamic skills ─────────────────┤
P2.2 bg process ─────────────────────┤── Phase 2 (needs Phase 1)
P2.3 clarify ────────────────────────┤
                                     │
P3.1 execute_code ───────────────────┤
P3.2 browser session ────────────────┤── Phase 3 (needs Phase 1+2)
P3.3 TTS tool ───────────────────────┤
                                     │
P4.1 session search ─────────────────┤
P4.2 context compaction ─────────────┤── Phase 4 (needs Phase 2+3)
P4.3 deferred tools ─────────────────┤
P4.4 multi-platform (skipped) ────────┤
P4.5 tech debt ──────────────────────┘
                                      │
P5.1 vector store ────────────────────┤
P5.2 ingestion ───────────────────────┤── Phase 5 (DONE ✅)
P5.3 search + rerank ─────────────────┤
P5.4 RAG integration ─────────────────┘
                                      │
P6.1 session persistence ─────────────┤
P6.2 login automation ────────────────┤
P6.3 multi-step scripting ────────────┤── Phase 6 (DONE ✅)
P6.4 monitoring & scraping ───────────┘
                                      │
P7.1 autonomous loop ─────────────────┤
P7.2 proactive capabilities ──────────┤── Phase 7 (needs Phase 5+6)
P7.3 proactive voice ─────────────────┤
P7.4 safety & observability ──────────┘
```

---

## 5. Risk Assessment

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| Binary size bloat (fitur baru) | Medium | Low | Build tags untuk optional features |
| Tool count overload (LLM confused) | Medium | Medium | Phase 4.3 (deferred loading) menyelesaikan |
| Breaking existing tools saat refactor | Low | High | Test sebelum deploy, rollback binary |
| Memory inject bloat system prompt | Low | Medium | Compact format, max N entries |
| Browser session leak (zombie Chrome) | Medium | Medium | Timeout + cleanup per session |

---

## 6. Naming & Convention

- **File naming:** `tools_{domain}.go` (contoh: `tools_patch.go`, `tools_todo.go`, `tools_bg.go`)
- **Function naming:** `execute{ToolName}(args map[string]interface{}) (string, bool)`
- **Registration:** di `registry_init2.go` dengan `registerTool()`
- **Schema:** di `registry_init2.go` atau inline di file tool
- **Build tags:** fitur berat (browser, python) → optional via `+build` tag
- **Config:** `~/.scorp-agent/` directory, JSON format
- **Error return:** `(message, false)` untuk error, `(message, true)` untuk success

---

## 7. Testing Strategy

Setiap phase:
1. **Build:** `make` → pastikan 0 errors
2. **Deploy:** `sudo cp scorp-agent /usr/local/bin/ && sudo systemctl restart scorp-agent`
3. **Smoke test via Telegram:** kirim test command ke agent mode
4. **Verify logs:** `journalctl -u scorp-agent --since "1 min ago"`
5. **Regression:** test 2-3 tool lama yang ada di file yang sama

---

*End of Planning Document*
