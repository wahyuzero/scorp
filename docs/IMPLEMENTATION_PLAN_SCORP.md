# SCORP AGENT — IMPLEMENTATION PLAN 2026→
> Basis: Feature research & community sentiment (see `docs/RESEARCH_AI_AGENT_FEATURES_2026.md`).  
> Structuring principles: Every item must have (a) community evidence, (b) a concrete design for Scorp's architecture, (c) effort estimate S (<1 day) / M (1–3 days) / L (>3 days).  

> **STATUS (2026-09-06): ALL ITEMS IMPLEMENTED & live-verified on tencent-vps —**  
> P0.1 bwrap sandbox + non-root daemon ✅ · P0.2 deny-rule engine ✅ · P0.3 test-integrity gate ✅ ·  
> P1.4 plan mode ✅ · P1.5 persistent ledger ✅ · P1.6 checkpoint/rewind ✅ · P1.7 auto memory ✅ ·  
> P2.8 delegate (audit + routing through gate stack) ✅ · P2.9 MCP deferred ✅ · P2.10 compaction preservation ✅ ·  
> P4.15 `scorp eval` (14-case arena + automated pre-deploy gate via `scripts/deploy.sh`) ✅ · P3.12 hooks ✅ ·  
> P3.13 auto-classifier (`auto` mode, PermissionDecision) ✅ · P3.14 MCP contract watch (fingerprint + warning) ✅ ·  
> P4.16 evidence-based claims (claim gate) ✅ — commits a1bb393..5310066.  

## CORE PRINCIPLES (From Research Consensus)

1. **State in Files, Not in Chat** — "All progress not recorded in memory is at risk."
2. **Deny Rules Apply in ALL Modes** — Even in YOLO (Claude Code 2.1 pattern).
3. **Sandboxing Before Autonomy** — Not prompt fatigue as the primary control.
4. **Merge-Rate > Test-Pass-Rate** — Verify artifacts, not agent claims.
5. **Harness Determines Results** — Token & context economics are core features, not implementation details.

## SCORP BASELINE vs 2026 ARCHITECTURE

Already aligned (do not touch): Task Ledger + plan gate + auto-resume · 3-tier autonomy + single-predicate `ConfirmationRequired()` · Sensitive path sandbox (hard, all modes) · Anti-fabrication gate · Compaction + New-Request Roll-Up · Deferred tool loading (`tool_search` + `tool_call`, TTL) · MCP client + marketplace · Skills · Scheduler/cron · Steering queue · Cooperative `/stop` · Receipts hash (audit trail) · `/usage` · Telegram-first (remote control model).

Resolved gaps: Real execution sandbox ✅ · Plan-mode workflow ✅ · Checkpoint/rewind ✅ · Subagent context isolation ✅ · Persistent task ledger ✅ · MCP tool schema deferred-by-default ✅ · Structured markdown memory (`MEMORY.md`) ✅ · Hooks ✅ · Auto-classifier ✅ · Encoded eval harness ✅.

---

## P0 — TRUST & SAFETY
*Highest priority: #1 community complaint is security/review burden; YOLO without sandbox is an "opt-in rootkit".*

### 1. Shell Execution Sandbox (Effort L, Impact H)
**Evidence**: Anthropic reports "sandboxing → -84% permission prompts"; dual isolation (fs + network) is the industry standard.  
**Design**: In `tools/exec.go ExecuteShell`, wrap `exec.Command` with bubblewrap (`bwrap --ro-bind / / --bind $PWD $PWD --tmpfs /tmp --unshare-all --share-net --die-with-parent`) when available; default network-deny + domain allowlists via environment. VPS daemon **migrated to a non-root user**.  
**Integration**: YOLO without bwrap triggers a permanent status warning; supervised default = sandbox ON; per-command escape hatch via existing confirmation gates.

### 2. Deny-Rule Engine (Effort M, Impact H)
**Evidence**: Deny rules remain enforced even under `bypassPermissions` (Claude Code docs).  
**Design**: Generalized `config.IsPathRestricted` → `config.DenyRules` with `tool(param:regex)` syntax (e.g. `shell(*:curl*|*nc*)`, `write_file(*:/etc/*)`), loaded from config; evaluated in `ExecuteTool` BEFORE `ConfirmationRequired` — active even in YOLO. Regression test per rule.

### 3. Test-Integrity Gate (Effort M, Impact H)
**Evidence**: Empirical audits found 19/45 "all tests pass" claims were fabricated; METR found 50% of PRs passing benchmarks were unmergeable; agents weaken test suites to turn CI green.  
**Design**: When a session touches test or CI files, `complete_task` is rejected unless a shell receipt with a passing test-suite run exists *after* the latest file edit (receipts.json hashes output as cryptographic proof). Symmetrical to the plan-completion gate.

---

## P1 — PLAN & CHECKPOINT
*Priority: Plan mode is the #1 praised feature; checkpointing is #5; "state in files" is #1 consensus.*

### 4. Plan Mode Workflow (Effort M, Impact H)
**Design**: `/plan <goal>` runs the loop with `IsToolAllowed` restricted to read-only tools + task ledger created → plan rendered via Telegram inline keyboard ("✅ Approve plan" / "✏️ Revise" / "❌ Cancel") or CLI prompt → approval transitions directly to execution mode with the same ledger.

### 5. Task Ledger Persistence (Effort S, Impact H)
**Design**: `taskPlans` flushed to `<session>.plan.json` on each update; reloaded when the session is touched; daemon restarts preserve active plans. Completed plans are cleanly purged from disk.

### 6. Checkpoint / Rewind (Effort M/L, Impact H)
**Evidence**: "FINALLY checkpoints!" (HN); Cursor & Claude Code made it baseline UX.  
**Design**: Shadow commits created under `refs/scorp/ckpt` before modifying turns (in Git repos); `/undo` restores the latest snapshot (code-first). Up to N=20 checkpoints retained per session.

### 7. Auto Memory Protocol (Effort S/M, Impact M/H)
**Evidence**: "Agents forgetting across sessions" is a permanent complaint; auto-memory files are Anthropic's official solution.  
**Design**: Upgraded flat KV store to project-scoped `MEMORY.md`; at `complete_task`, the agent distills decisions and state via a dedicated extraction prompt; on session start, `MEMORY.md` is injected after the system prompt (~200-line quota).

---

## P2 — TOKEN & CONTEXT ECONOMY
*Evidence: Harness overhead is the #2 complaint; Claude Code consumes 33k tokens before prompts vs. OpenCode's 7k.*

### 8. Subagent `delegate` (Effort L, Impact H)
**Design**: Tool `delegate(task, max_turns)` spawns a child loop with a FRESH context (minimal system prompt + task), returning only the final summary to the parent as the tool result. Bounded by turn caps and strict 6-minute wall-clock limits. Cuts quadratic context accumulation ($N \times (N+1)/2$).

### 9. Deferred-by-Default MCP Tool Schemas (Effort S/M, Impact M/H)
**Evidence**: CTO of Perplexity discarded MCP due to 15–20k token schema overhead; best practice recommends 10–15 active tools maximum.  
**Design**: Mark MCP tools `deferred:true`, loaded dynamically via `tool_search` with TTL auto-eviction.

### 10. Compaction Preservation (Effort S, Impact M)
**Design**: `maybeCompactHistory` preserves key state: active task ledger, user goals, and previous final summaries ALWAYS survive compaction passes.

### 11. Per-Session Cost Tracking (Effort S, Impact M)
**Design**: Session and turn aggregation in `model_usage.json`, exposed via `/usage` with threshold warnings when spending exceeds limits.

---

## P3 — EXTENSIBILITY & AUTONOMY UX

### 12. PreToolUse/PostToolUse Hooks (Effort M, Impact M/H) — ✅ Implemented
**Evidence**: "CLAUDE.md says 'please', hooks say 'must'" — deterministic enforcement praised by enterprise teams.  
**Design**: Configuration via `SCORP_HOOKS_PRE` and `SCORP_HOOKS_POST` (`tool_pattern:shell_command`). Evaluated at the single execution bottleneck (`ExecuteTool`); exit code 2 blocks execution, stdout injects additional context.

### 13. Auto-Mode Classifier (Effort M/L, Impact H) — ✅ Implemented
**Evidence**: Anthropic data indicates humans catch 13.6% of dangerous commands vs. classifier catching 89%; auto-mode became default in late 2026.  
**Design**: 4th tier `auto` between supervised and yolo: cheap model + deterministic heuristics grade every tool call (safe → execute, risky → confirm, destructive → hard-deny unless allowlisted). Falls back to supervised behavior after 3 consecutive uncertain decisions.

### 14. MCP Contract Watch (Effort S, Impact M) — ✅ Implemented
**Design**: SHA-256 fingerprinting (name, description, inputSchema) per MCP server in `~/.scorp/mcp_contracts.json`. Schema drift triggers immediate warnings.

---

## P4 — VERIFICATION & QUALITY GATES

### 15. `scorp eval` — Private Evaluation Arena (Effort M, Impact H) — ✅ Implemented
**Evidence**: SWE-bench was retired; modern consensus is "build private arenas from real-world backlog issues".  
**Design**: Encoded suite: deterministic safety/persistence checks + live agent tasks verified by independent artifact checkers; integrated into `scripts/deploy.sh` as an automated deployment gate.

### 16. Evidence-Based Claim Gate (Effort S, Impact M) — ✅ Implemented
**Design**: Claims like "all tests pass" without corresponding execution receipts in the current session window are rejected with a verification nudge.

---

## EXPLICITLY NOT BUILT
- **Agent teams / Swarms** — Community evidence shows high token expenditure and coordination overhead; deferred until delegation matures.
- **Full-repo graphs (Aider style)** — Deferred until multi-repo demands rise.
- **Voice / IDE plugins / Cloud agents** — Out of scope (Telegram-first remains Scorp's core differentiator).

## SUCCESS METRICS (Merge-Rate Mindset)
- Permission prompts per task decreased significantly.
- Token consumption per task reduced by ≥40% (via `delegate` and deferred MCP).
- Eval pass rate at 100% across all regression cases.
- Zero loss of task plans across daemon restarts.
- Zero uncontained executions outside the sandbox.
