# Scorp-Agent Upgrade — Checklist

> Companion to `UPGRADE_PLAN.md`. Update status saat implementasi.
> Legend: `[ ]` todo · `[~]` in progress · `[x]` done · `[!]` blocked · `[-]` skipped

---

## Phase 1 — Quick Wins (P0)

### P1.1 — `patch` Tool (fuzzy find-replace)
- [x] Buat file `tools_patch.go`
- [x] Implement `executePatch()` dengan 3 fuzzy matching strategies
  - [x] Strategy 1: Exact match
  - [x] Strategy 2: Whitespace-normalized match
  - [x] Strategy 3: Line-by-line fuzzy (trim leading/trailing ws)
- [x] Implement V4A patch format parser (`*** Begin Patch` / `*** End Patch`)
- [x] Return unified diff sebagai output
- [x] Validasi: error jika `old_string` tidak unik tanpa `replace_all=true`
- [x] Daftarkan di `registry_init2.go` dengan schema lengkap
- [x] Build: `make` → 0 errors
- [x] Deploy: `sudo cp + systemctl restart`
- [x] Test via Telegram: patch file kecil, verify diff
- [x] Test edge case: file tidak ada, old_string tidak ditemukan, multi-match

### P1.2 — `todo` Tool (task tracking)
- [x] Buat file `tools_todo.go`
- [x] Define struct `TodoItem { ID, Content, Status }`
- [x] Implement `executeTodo()`:
  - [x] No args → return formatted list
  - [x] `todos` + `merge=false` → replace entire list
  - [x] `todos` + `merge=true` → update by ID, add new
  - [x] Enforce: max 1 `in_progress` at a time
- [x] State: in-memory `[]TodoItem` per chat session
- [x] Daftarkan di `registry_init2.go`
- [x] Build + deploy
- [x] Test via Telegram: create 3-item list, update status, verify

### P1.3 — Memory Auto-Inject
- [x] Baca `agent_prompt.go`, identifikasi `buildSystemPrompt()`
- [x] Tambah fungsi `buildMemorySection()` — baca `memory.json`, format compact
- [x] Inject memory section ke system prompt (sebelum `[TOOLS]` section)
- [x] Format: `§ key: value` (delimited, compact)
- [x] Max display: 20 entries, truncate value > 200 char
- [x] Build + deploy
- [x] Test: set memory via tool → agent turn berikutnya → agent aware
- [x] Verify token count tidak blow up

---

## Phase 2 — Core Improvements (P1)

### P2.1 — Dynamic Skills System (`skill_manage`)
- [x] Buat directory `~/.scorp-agent/skills/`
- [x] Define format: skill JSON dengan fields:
  - [x] `name`, `emoji`, `description`, `category`, `prompt`, `examples`, `auto_load_keywords`
- [x] Migrasi 8 hardcoded skills dari `skills.go` ke JSON files
- [x] Refactor `skills.go`:
  - [x] `loadSkills()` — baca dari JSON files saat startup
  - [x] `listSkills()` — iterate loaded skills
  - [x] `handleSkillCommand()` — cari di loaded skills
  - [x] `getSkillPromptForMessage()` — match keywords dari loaded skills
- [x] Buat file `tools_skill_manage.go`:
  - [x] `action=create` — write JSON file
  - [x] `action=list` — list all skill files
  - [x] `action=view` — read skill detail
  - [x] `action=update` — edit skill JSON
  - [x] `action=delete` — remove skill file (block if default skill)
  - [x] Validate JSON, sanitize name
- [x] Daftarkan `skill_manage` di registry
- [x] Build + deploy
- [x] Test: create new skill via agent → use via `/skill`
- [x] Test: update existing skill
- [x] Test: auto-load keyword trigger

### P2.2 — Background Process Management (`bg` + `uptime`)
- [x] Buat file `tools_bg.go`
- [x] Define struct `BGProcess { ID, Cmd, StartTime, Output, Done, ExitCode }`
- [x] Implement actions:
  - [x] `spawn` — start bg process, return session_id
  - [x] `list` — all bg processes with status
  - [x] `poll` — incremental output since last poll
  - [x] `wait` — block until done or timeout
  - [x] `kill` — terminate process
  - [x] `write` — send raw stdin
  - [x] `submit` — send data + Enter (for interactive prompts)
- [x] Implement PTY mode untuk interactive tools
- [x] Auto-cleanup: kill zombie processes after 30 min idle
- [x] Daftarkan di registry dengan schema lengkap
- [x] Build + deploy
- [x] Test: spawn `sleep 60` → poll → kill
- [x] Test: spawn interactive script → submit input → poll output
- [x] Test: wait with timeout

### P2.3 — `clarify` Tool
- [x] Buat file `tools_clarify.go`
- [x] Implement `executeClarify()`:
  - [x] `question` (required) + `choices` (optional, max 4)
  - [x] Jika ada choices → kirim inline keyboard ke Telegram
  - [x] Jika tidak ada → kirim pesan, tunggu text reply
- [x] Implement blocking mechanism di agent loop:
  - [x] Channel untuk receive user response
  - [x] Register pending clarify per chat session
  - [x] Intercept callback/text → resolve channel → return ke agent
- [x] Timeout (default 120s) → return "No response"
- [x] Daftarkan di registry
- [x] Build + deploy
- [x] Test: multiple choice → tap button → agent continues
- [x] Test: open-ended → type reply → agent continues
- [x] Test: timeout

---

## Phase 3 — Advanced (P2)

### P3.1 — `execute_code` (Python orchestration)
- [x] Buat file `tools_exec_code.go`
- [x] Buat Python helper module (`scorp_tools.py`):
  - [x] `read_file(path)` wrapper
  - [x] `write_file(path, content)` wrapper
  - [x] `terminal(command)` wrapper
  - [x] `search_files(pattern, path)` wrapper
  - [x] `patch(path, old, new)` wrapper
  - [x] `json_parse(text)` helper
  - [x] `shell_quote(s)` helper
  - [x] `retry(fn, max_attempts, delay)` helper
- [x] Inject module ke Python subprocess via file-based bridge IPC (`/tmp/scorp_bridge/`)
- [x] Implement:
  - [x] 5-minute timeout
  - [x] 50KB stdout cap
  - [x] Max 50 tool calls per script (counter in wrapper)
  - [x] Capture stdout sebagai result
- [x] Daftarkan di registry
- [x] Build + deploy
- [x] Test: simple `print("hello")` script
- [x] Test: loop read 3 files
- [x] Test: batch search + patch
- [x] Test: timeout enforcement

### P3.2 — Browser Session Persistence
- [x] Refactor `browser.go`:
  - [x] Pool of named chromedp contexts (keyed by chat ID)
  - [x] `browser goto` → create or reuse context
  - [x] `browser click/type/extract` → use existing context
  - [x] `browser reset` → close context, clear state
  - [x] Context timeout: 10 min idle → auto close
- [x] Add `browser_snapshot`:
  - [x] Extract accessibility tree (DOM walk)
  - [x] Assign ref IDs to interactive elements
  - [x] Return compact text representation
- [x] Add `browser_console`:
  - [x] Capture console.log/warn/error
  - [x] Capture uncaught JS exceptions
  - [x] Return as text
- [x] Build + deploy
- [x] Test: login flow (goto → fill → click → screenshot, same session)
- [x] Test: snapshot returns ref IDs
- [x] Test: console captures JS error

### P3.3 — TTS as Agent Tool (`speak`)
- [x] Buat `tools_speak.go` (atau tambah di `voice.go`)
- [x] Implement `executeSpeak()`:
  - [x] Reuse `synthesizeSpeech(text, voice)` dari voice.go
  - [x] Kirim hasil via `sendVoiceMessage(chatID, oggPath)`
  - [x] Cleanup temp file
- [x] Daftarkan di registry dengan schema:
  - [x] `text` (required), `voice` (optional)
- [x] Build + deploy
- [x] Test: agent calls speak → voice message arrives

---

## Phase 4 — Long-term (P3-P4)

### P4.1 — Session Search (FTS5)
- [x] Evaluasi: SQLite (FTS5) vs alternative (Bleve, etc.)
- [x] Design schema: `sessions`, `messages`, `messages_fts`
- [x] Migrasi session storage dari in-memory → SQLite
- [x] Implement search tool:
  - [x] `session_search(query)` → FTS5 match
  - [x] Return snippet + timestamp
- [x] Migration script untuk existing sessions
- [x] Build + deploy + test

### P4.2 — Context Compaction
- [x] Detect: total tokens approaching model context limit
- [x] Strategy: summarize oldest N messages
- [x] Inject summary as system message
- [x] Keep last K messages verbatim
- [x] Build + deploy + test

### P4.3 — Deferred Tool Loading
- [x] Split tools:
  - [x] Core (always loaded): shell, read_file, write_file, patch, todo, search_code, system_info
  - [x] Extended (on-demand): browser, http, sql, compose, delegate, mcp_*, dll
- [x] Metadata-only schema untuk extended tools (name + short desc, no full params)
- [x] Tool `tool_search(query)` → list matching
- [x] Tool `tool_call(name, args)` → load + execute
- [x] Build + deploy + test

### P4.4 — Multi-Platform
- [-] Abstract messaging interface: `MessageSender` interface
- [-] Implement Telegram adapter (refactor existing)
- [-] Implement Discord adapter
- [-] Per-platform config + routing
- [-] Build + deploy + test

### P4.5 — Technical Debt Cleanup
- [x] T01: Hapus `executeMCPTool()` dari `agent_loop.go:38-75`
- [x] T02: Vision model baca dari `models.json` (field `vision_model`) atau dari arg
- [x] T03: Skills auto-load keyword matching (digabung dengan P2.1)
- [x] T04: Browser screenshot path dari config atau `$HOME/.scorp-agent/screenshots`
- [x] T05: RAG index path pakai `$HOME/.scorp-agent/rag_index.json` (bukan hardcoded `/root/`)
- [x] T06: Voice paths dari config
- [x] T07: Graceful shutdown — context cancel untuk MCP servers, wait group

---

## Phase 5 — Semantic RAG (DONE ✅)

### P5.1 — Vector Store Setup
- [x] Evaluasi: Pure Go SimHash (embedded, zero dep) vs Python sentence-transformers (heavy, compile PyTorch). **Chose SimHash.**
- [x] Schema: `VecChunk` struct (id, source, content, simhash uint64, created)
- [x] Embedding: SimHash (64-bit FNV-1a fingerprint) + hybrid TF-IDF fusion
- [x] Config: `~/.scorp-agent/rag_vector.json`, auto-persist on write

### P5.2 — Ingestion Pipeline
- [x] File/dir ingestion: SmartChunk (paragraph/sentence-aware, max-chunk config), auto-skip binaries/images/archives
- [x] Incremental updates: auto-remove old chunks from same source before re-index
- [x] Tools: `ragvec_ingest(path)`, `ragvec_list()`, `ragvec_remove(source)`

### P5.3 — Semantic Search + Re-ranking
- [x] Vector search: SimHash → Hamming distance → similarity score (0.0-1.0)
- [x] Re-ranking: SimHash primary (semantic) + TF-IDF fallback (keyword) + hybrid weighted fusion
- [x] Hybrid search: `simhashWeight` configurable (default 0.7 simhash + 0.3 TF-IDF)
- [x] Tool: `ragvec_search(query, top_k, hybrid, vector_weight)` → ranked chunks with scores

### P5.4 — RAG Integration
- [x] Agent auto-RAG: auto-search + inject top 3 chunks ke system prompt di setiap turn
- [x] Score display: results ranked with confidence score
- [x] Tool: `ragvec_provider()` untuk cek status embedding provider

---

## Phase 6 — Browser Sessions (DONE ✅)
> Persistent browser contexts: cookies, login flows, multi-step automation

### P6.1 — Session Persistence (DONE ✅)
- [x] User data dir: persistent chromedp context per chat ID (`~/.scorp-agent/browser_data/<chatID>`)
- [x] Cookie jar: automatic via Chrome UserDataDir (survives across tool calls)
- [x] Session reuse: `getOrCreateBrowserSession()` caches contexts
- [x] Auto-cleanup: stale sessions pruned every 5 min (10 min idle timeout)

### P6.2 — Login Flow Automation (DONE ✅)
- [x] Credential store: AES-256-GCM encrypted vault (`~/.scorp-agent/vault.json`)
- [x] Master key: auto-generated, stored at `~/.scorp-agent/vault.key`
- [x] Form detection: JS-based username/email/password field auto-detect
- [x] Tool: `vault` (get/set/list/remove) + `autologin` (auto-fill from vault)

### P6.3 — Multi-step Scripting (DONE ✅)
- [x] JSON script format: array of steps with action/url/selector/value/wait_ms/timeout/name/retry
- [x] Actions: goto, click, type, submit, wait, extract, screenshot, evaluate, scroll
- [x] Error handling: per-step retry count, screenshot on fail
- [x] Captures: extract/screenshot/evaluate results stored in named variables
- [x] Tool: `script` (source/inline) + `script_list`

### P6.4 — Monitoring & Scraping (DONE ✅)
- [x] Scheduled scraping: cron-integrated browser monitor with change detection
- [x] Change detection: SHA-256 content hash diff, alert on change
- [x] Dashboard screenshots: scheduled captures saved to `~/.scorp-agent/screenshots/`
- [x] Auto-RAG feed: extracted content auto-ingested to vector index
- [x] Tool: `monitor_add` / `monitor_list` / `monitor_remove` / `monitor_check`

---

## Phase 7 — Autonomous Agent ✅

### P7.1 — Autonomous Loop Framework
- [x] Scheduler: cron-like ticker + event-driven triggers (metrics, logs, webhooks)
- [x] Context builder: gather relevant state (metrics, logs, RAG, session history)
- [x] Decision engine: LLM prompt → structured action plan (JSON)
- [x] Action executor: validate → execute tools → verify result
- [x] Memory: persist decisions + outcomes untuk learning (`autonomous_log.json`)

### P7.2 — Proactive Capabilities
- [x] Resource monitoring: auto-scale, restart unhealthy containers, alert
- [x] Security response: auto-ban suspicious IP (fail2ban integration)
- [x] Knowledge maintenance: ragvec_add untuk ingest findings to RAG
- [x] Self-healing: detect broken tools/config → attempt fix via exec

### P7.3 — Proactive Voice
- [x] Config: `voice: off|important|always` di autonomous config
- [x] Hook: after autonomous cycle, if speak=true → auto `synthesizeSpeech()`
- [x] Notify: Telegram message dengan analysis + actions summary

### P7.4 — Safety & Observability
- [x] Approval gates: high-risk actions butuh confirmation (low/medium/high)
- [x] Audit log: every autonomous action logged dengan reasoning + result
- [x] Kill switch: file-based (`autonomous_kill`), disable instant, revive command
- [x] Dangerous command blocklist: rm -rf, dd, mkfs, shutdown, halt, reboot
- [x] Max actions per cycle (default 5, configurable 1-10)
- [x] Test suite: 15 tests (kill switch, config, audit log, context, decision parsing, exec, risk gate, danger block, tool actions)

---

## Post-Implementation Review

Setelah setiap phase selesai:
- [ ] Binary size check (`ls -lh scorp-agent`)
- [ ] Startup time check (journalctl first log)
- [ ] Tool count di logs
- [ ] Regression test 3 tools random
- [ ] Update `MEMORY` entry jika ada yang perlu diingat
- [ ] Git commit dengan message: `feat(phase-X.Y): description`

---

## Quick Reference — Commands

```bash
# Build
cd /home/ubuntu/projects/vps-monitor-go && make

# Minimal build (no browser)
make minimal

# Deploy
sudo cp scorp-agent /usr/local/bin/scorp-agent && sudo systemctl restart scorp-agent

# Logs
journalctl -u scorp-agent --since "1 min ago" -f

# Check registered tools
journalctl -u scorp-agent | grep "registered\|tools registered"

# Config files
ls -la ~/.scorp-agent/
cat ~/.scorp-agent/mcp.json | jq .
cat ~/.scorp-agent/models.json | jq .
cat ~/.scorp-agent/scheduler.json | jq .
```

---

*Last updated: 16 Juni 2026*
