# 🦂 Scorp

**Small agent, big tools.** An ultra-light, ultra-fast Go autonomous agent that connects modern LLMs (Command Code, DeepSeek, Gemini, Claude, OpenAI, Ollama, Qwen, and custom OpenAI-compatible endpoints) directly to your operating system — shell, files, web, database, code, and cron — accessible via terminal, Telegram, or an embedded Web Dashboard.

Runs natively with a memory footprint of **under 25MB RAM** on Linux VPS, edge nodes, and Android Termux as a single standalone static binary.

---

## ⚡ Core Architecture & Philosophy

Scorp is engineered around modern production agent principles:
- **State in Files, Not in Chat** — Active plans, shadow checkpoints, and long-term memory live on disk (`plans/<session>.plan.json`, `refs/scorp/ckpt`, `MEMORY.md`), completely surviving daemon restarts and session resets.
- **Evidence Over Model Claims (Merge-Rate Mindset)** — Verbal assertions like *"all tests pass"* or *"file deleted"* are rejected by the Test-Integrity and Claim Gates unless backed by cryptographic SHA-256 receipts.
- **Multi-Layered Trust & Safety** — Hard deny rules, Bubblewrap container sandboxing, and sensitive path protections stay enforced across *all* autonomy modes, including YOLO.
- **Radical Token Economy** — Context-isolated subagent delegations (`delegate`), deferred-by-default MCP tools with TTL caching, and preservation-aware token compaction slash token consumption by 40–90%.

---

## 🚀 Quick Start (Under 60 Seconds)

### 1. Installation

#### Option A: One-Line Universal Install (Recommended)
Automatically detects your OS and architecture, downloads the matching prebuilt static binary, installs to system PATH, and launches setup:
```bash
bash -c "$(curl -fsSL https://raw.githubusercontent.com/wahyuzero/scorp/master/install.sh)"
```
> Supported out-of-the-box: **Linux** (x86_64, aarch64/ARM64, ARMv7), **macOS** (Apple Silicon M-series & Intel x64), and **Android** (Termux).

#### Option B: Build from Source
```bash
git clone https://github.com/wahyuzero/scorp.git
cd scorp
make
./scorp setup
```

### 2. Interactive Setup Wizard (`scorp setup`)
Run `scorp setup` at any time to configure or reconfigure your agent:
1. **Operating Mode:** CLI Terminal Mode, 24/7 Telegram Bot Daemon, or Dual Mode.
2. **AI Provider:**
   - 🟢 **Verified & Tested:** Google Gemini (3.7 Flash, 3.8 Flash, 3.1 Pro), OpenCode Zen (`big-pickle`, `mimo-v2.5-free`), Command Code (DeepSeek v4 Flash, Laguna, GLM).
   - 🔵 **Local & Offline:** Ollama Local (no API key required), Claude CLI Bridge.
   - 🟡 **Cloud & Compatible:** OpenAI (GPT-4o, o3-mini), Anthropic Claude, DeepSeek Official, Mistral, xAI Grok, Perplexity Sonar, or Custom OpenAI-compatible endpoints.
3. **API Keys & Security:** Auto-detects existing `.env` credentials or securely stores new ones.
4. **Autonomy Profile:** `supervised` (default), `auto` (smart risk classifier), `readonly` (audit-only), or `yolo` (unattended).
5. **Systemd Service (Linux/VPS):** Optional one-click installation and enablement of `scorp.service`.

### 3. Choose Your Interface

* **Interactive Terminal REPL:**
  ```bash
  scorp --cli
  ```
* **One-Shot CLI Execution:**
  ```bash
  scorp "inspect disk usage and memory, then report"
  ```
* **24/7 Background Telegram Daemon:**
  ```bash
  scorp daemon
  # or manage via systemd if installed:
  sudo systemctl status scorp
  ```
* **Embedded Web Gateway & Dashboard (< 2MB RAM):**
  ```bash
  scorp gateway --port 8080
  ```

---

## 🛡️ Trust & Safety Stack (Defense in Depth)

Scorp implements a multi-layer gate stack that evaluates every tool call before execution:

```
[Tool Call]
    ↓
1. Deny-Rule Engine (SCORP_DENY_RULES — unbypassable in ANY mode)
    ↓
2. Sensitive Path Sandbox (/etc/shadow, ~/.ssh, credentials, .env)
    ↓
3. 4-Tier Autonomy Gate (readonly | supervised | auto | yolo)
    ↓
4. PreToolUse Hooks (SCORP_HOOKS_PRE — exit 2 hard block)
    ↓
5. Execution Sandbox (Bubblewrap bwrap: read-only root, isolated network)
    ↓
6. Evidence & Receipt Logging (SHA-256 hash + test-integrity validation)
    ↓
7. PostToolUse Hooks & Secret Redaction (masking API keys before LLM sees output)
```

### 4 Autonomy Tiers (`--mode=<level>`)
| Mode | Behavior | Use Case |
| :--- | :--- | :--- |
| `readonly` | Only non-mutating tools allowed (`read_file`, `web_search`, `log`, etc.). Write/exec tools blocked. | Safe code auditing, research, read-only analysis |
| `supervised` | **Default.** Non-destructive commands run automatically; dangerous operations prompt for confirmation. | Everyday interactive development & server admin |
| `auto` | Fast LLM classifier + deterministic heuristics grade actions: safe → run, risky → prompt, destructive → deny (unless allowlisted via `SCORP_AUTO_ALLOW`). | Semi-autonomous workflows with human fallback |
| `yolo` | Full unattended autonomy. Confirmation prompts bypassed, but deny rules & sandboxing still apply. | Automated CI/CD, batch processing, scripted pipelines |

---

## 🧩 Key Innovations

### 1. Plan Mode & Persistent Task Ledger
* `/plan <goal>` enters a sandboxed read-only drafting phase.
* Scorp constructs an explicit multi-step ledger and renders an interactive approval prompt (Terminal or Telegram inline keyboard).
* Once approved, execution resumes with the same ledger. The agent loop automatically resumes until all items are marked `done`.
* Ledgers are atomically persisted to disk (`~/.scorp/plans/<session>.plan.json`) and survive daemon reboots.

### 2. Checkpoint & Rewind (`/undo`)
* Before any modifying turn, Scorp captures an isolated Git shadow commit under `refs/scorp/ckpt/<session>/<timestamp>`.
* Never pollutes your working tree or branch history.
* Type `/undo` in the CLI or Telegram to roll back changes step-by-step (up to 20 checkpoints per session).

### 3. Subagent Delegation (`delegate` & `delegate_batch`)
* Spawn independent subagents for complex research, codebase exploration, or parallel checks.
* Runs in a fresh, isolated context window with its own turn cap and strict 6-minute wall-clock limit.
* Prevents token quadratic explosion, collapsing large multi-file investigations into a single concise summary.

### 4. Test-Integrity & Evidence Claim Gates
* If a session touches test files or CI configurations, `complete_task` is blocked until a successful test-suite run receipt exists *after* the latest file modification.
* Verbal model claims like *"all tests are green"* are cross-checked against actual execution receipts in `receipts.json`.

### 5. Durable Project Memory (`MEMORY.md`)
* Upon completing tasks, Scorp automatically distills architectural decisions, operational state, and preferences into `MEMORY.md`.
* Injected into subsequent sessions without manual bookkeeping, capped at ~200 lines to keep context lean.

### 6. Deterministic Hooks (`SCORP_HOOKS_PRE` / `SCORP_HOOKS_POST`)
* Deterministic shell hooks running at tool boundaries.
* PreToolUse hooks receiving JSON on stdin: exit `0` adds context to the prompt, exit `2` hard-blocks the tool call, timeouts (10s) fail-safe.

---

## 🔌 MCP Ecosystem & Tri-Option Marketplace

Scorp features a built-in Model Context Protocol (MCP) engine (Stdio and SSE HTTP JSON-RPC 2.0):

* **Deferred-by-Default (Zero Context Bloat):** MCP tool schemas are not injected into the system prompt. The model discovers them on-demand via `tool_search`, caching active tools with an automatic TTL eviction policy.
* **Contract Watch & Watchdog:** Tracks SHA-256 tool schema fingerprints in `mcp_contracts.json` to detect silent upstream API drift. Automatic crash recovery with exponential backoff.
* **Tri-Option Marketplace (`/mcp`):**
  1. **Pure Go Native:** Zero runtime dependencies, instant execution.
  2. **AI-Transpiler Local Rebuild:** Transpiles Node.js/Python MCP packages into compiled Go binaries or isolated sandboxes.
  3. **Upstream Runner:** Standard Stdio/SSE subprocesses (`npx`, `uvx`, `docker`).
* Commands: `/mcp search [query]`, `/mcp install <name> [1|2|3]`, `/mcp share <name>`, `/mcp restart <server>`.

---

## 🧰 45+ Built-in Tools Catalog

| Category | Tools | Description |
| :--- | :--- | :--- |
| **Core OS & Files** | `read_file`, `write_file`, `replace_file_content`, `patch`, `list_dir`, `search_code`, `system_info`, `send_file` | Surgical diff replacements, unified diffs, regex code search, and file transfers. |
| **Shell & Sandboxing** | `shell`, `exec_code`, `process`, `git` | Bash execution wrapped in Bubblewrap (`bwrap`), process-group kill on timeout, and Python orchestration bridge (`import scorp`). |
| **Web & Search** | `web_search`, `read_url`, `web_fetch`, `fastscraper` | Zero-API-key embedded metasearch (DuckDuckGo, SearXNG, Tavily, Brave) + Readability Markdown scraper (<5MB RAM). |
| **Database** | `sql` | Query SQLite, PostgreSQL, and MySQL with safety gates requiring confirmation on mutating queries. |
| **RAG & Memory** | `memory`, `index_*`, `ragvec_*`, `session_search` | `MEMORY.md` durable storage, Simhash fast indexing, local cosine vector search, and SQLite FTS5 session search. |
| **Agent Orchestration**| `delegate`, `delegate_batch`, `task_plan`, `complete_task`, `clarify`, `sop`, `skill_manage` | Subagent spawning, task ledger management, human clarification prompts, and SOP playbooks. |
| **Scheduler & Android**| `schedule_manage`, `termux_api`, `bg` | In-process cron scheduler with overlap protection, background jobs, and Android WakeLock/sensors/SMS. |

---

## 🖥️ CLI Reference & Slash Commands

### Subcommands & Binary Flags
```bash
./scorp                         # Launch interactive CLI or Telegram daemon
./scorp "prompt"                # One-shot CLI task
./scorp --mode=auto             # Set autonomy level (readonly | supervised | auto | yolo)
./scorp -s <session_name>       # Attach or create a named conversation session
./scorp setup                   # Run the interactive configuration wizard
./scorp gateway --port 8080     # Start micro Web Dashboard & REST API (<2MB RAM)
./scorp sop list / sop run <id> # List or execute automated SOP playbooks
./scorp eval [--live]           # Run the private evaluation arena & quality gates
./scorp --mcp-server            # Expose Scorp as a standard MCP server to other agents
./scorp update                  # Self-update binary from GitHub releases
./scorp version                 # Print version and target architecture
```

### Interactive Slash Commands
| Command | Description |
| :--- | :--- |
| `/help` | Show command list and usage guidelines |
| `/plan <goal>` | Enter Plan Mode: draft read-only plan, approve, and execute |
| `/undo` | Roll back repository changes to the previous shadow checkpoint |
| `/stop` | Cooperative kill-switch: abort current execution cycle immediately |
| `/mode [level]` | View or switch autonomy level (`readonly`, `supervised`, `auto`, `yolo`) |
| `/models` / `/model <name>`| List available LLM models or switch active model on the fly |
| `/mcp` | MCP dashboard: health, marketplace search, install, transpile, restart |
| `/cron` | List, pause, resume, or trigger scheduled tasks |
| `/compact` | Manually trigger token compaction while preserving task state |
| `/session` / `/sessions` | List, switch, rename, or delete conversation sessions |
| `/receipts` | Inspect recent cryptographic SHA-256 tool execution audit receipts |
| `/status` | Display system metrics, memory, active locks, and uptime |
| `/cost` / `/usage` | Show token consumption and estimated cost for the current session |
| `/tools` | List all registered native and dynamic tools |
| `/sop [run <name>]` | Execute or inspect Standard Operating Procedures |
| `/skills` | Inspect and toggle domain-specific instruction skills |
| `/clear` | Clear current session conversation history |
| `/exit`, `/quit`, `/q` | Exit the CLI session |

---

## ⚙️ Environment Variables

| Variable | Default | Description |
| :--- | :--- | :--- |
| `SCORP_AUTONOMY` | `supervised` | Default autonomy mode: `readonly`, `supervised`, `auto`, `yolo`. |
| `SCORP_AUTO_ALLOW`| - | Semicolon-separated regex allowlist for auto-mode shell commands (e.g. `^rm -rf /tmp/.*`). |
| `SCORP_DENY_RULES`| - | Hard deny rules enforced in all modes (e.g. `shell(command:mkfs);write_file(path:^/etc/)`). |
| `SCORP_HOOKS_PRE` | - | PreToolUse hook specs: `tool_pattern:command` (exit 2 to block). |
| `SCORP_HOOKS_POST`| - | PostToolUse hook specs: `tool_pattern:command` (informational context). |
| `SCORP_SESSION`   | `default` | Default session identifier. |
| `SCORP_DEBUG`     | `0` | Set to `1` to stream diagnostic logs to stderr instead of `~/.scorp/scorp.log`. |
| `TELEGRAM_BOT_TOKEN`| - | Telegram bot token for 24/7 background daemon mode. |
| `TELEGRAM_CHAT_ID`  | - | Authorized Telegram Chat ID (locks bot to owner). |
| `COMMANDCODE_API_KEY`| - | Command Code gateway API key. |
| `DEEPSEEK_API_KEY`  | - | DeepSeek API key (V3 / R1). |
| `GEMINI_API_KEY`    | - | Google Gemini API key. |
| `OPENAI_API_KEY`    | - | OpenAI API key. |
| `ANTHROPIC_API_KEY` | - | Anthropic Claude API key. |

---

## 🏗️ Build Targets & Deployment

### Build Targets (`Makefile`)
* **Standard Production Build:**
  ```bash
  make
  ```
* **Lightweight Headless (no Chromedp browser dependencies):**
  ```bash
  make minimal
  ```
* **Android Termux (Pure static ARM64 binary, zero libc dependencies):**
  ```bash
  make termux
  ```
  *(Produces `scorp-termux`, ready to run inside Termux with Termux:API).*

### Automated Eval-Gated Deployment
Scorp includes a production deploy script (`scripts/deploy.sh`) with an automated pre-flight quality gate:
```bash
./scripts/deploy.sh
```
Pipeline stages:
1. Local build, vet, and unit tests.
2. Source sync via rsync to target VPS.
3. Remote CGO build & remote test suite execution.
4. **Pre-Deploy Evaluation Gate:** Runs `scorp eval` on the candidate binary.
5. Atomic binary swap with MD5 verification and systemd restart.
*(If any test or evaluation fails, the active production binary remains untouched).*

---

## 📄 License

MIT License. Crafted with precision for high autonomy and minimal footprint.
