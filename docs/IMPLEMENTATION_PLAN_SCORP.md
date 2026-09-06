# SCORP AGENT — RENCANA IMPLEMENTASI 2026→
> Basis: riset fitur & sentimen komunitas (lihat `docs/RESEARCH_AI_AGENT_FEATURES_2026.md`).
> Prinsip penyusunan: setiap item wajib punya (a) bukti komunitas, (b) desain konkret untuk
> arsitektur scorp, (c) estimasi effort S (<1 hari) / M (1-3 hari) / L (>3 hari).

> **STATUS (2026-09-06): terlaksana & terverifikasi live di tencent-vps —**
> P0.1 sandbox bwrap + daemon non-root ✅ · P0.2 deny-rule engine ✅ · P0.3 test-integrity gate ✅ ·
> P1.4 plan mode ✅ · P1.5 ledger persisten ✅ · P1.6 checkpoint/rewind ✅ · P1.7 auto memory ✅ ·
> P2.8 delegate (audit + routing lewat gate stack) ✅ · P2.9 MCP deferred ✅ · P2.10 compaction
> preservation ✅ · P4.15 `scorp eval` (13 kasus, gerbang pra-deploy) ✅ — komit a1bb393..473d171.
> Tersisa: P3.12 hooks · P3.13 auto-classifier · P3.14 MCP contract watch · P4.16 klaim berbasis bukti.

## PRINSIP PENGEMBANGAN (dari konsensus riset)

1. **State di file, bukan chat** — "All progress not recorded in memory is at risk."
2. **Deny rules berlaku di SEMUA mode** — bahkan YOLO (pola Claude Code 2.1).
3. **Sandbox sebelum autonomi** — bukan prompt fatigue sebagai kontrol utama.
4. **Merge-rate > test-pass-rate** — verifikasi artefak, bukan klaim agent.
5. **Harness menentukan hasil** — ekonomi token & konteks adalah fitur, bukan detail.

## BASELINE SCORP vs TABEL 2026

Sudah selaras (jangan dirombak): Task Ledger + plan gate + auto-resume · autonomy 3 tingkat +
`ConfirmationRequired()` satu-predikat · sandbox path sensitif (hard, semua mode) · anti-fabrication
gate · compaction + New-Request Roll-Up · deferred tool loading (`tool_search`+`tool_call`, TTL) ·
MCP client + marketplace · skills · scheduler/cron · steering queue · cooperative `/stop` ·
receipts hash (audit trail) · `/usage` · Telegram-first (= "remote control" ala 2026, sudah dimiliki).

Gap terkonfirmasi: sandbox eksekusi nyata ❌ · plan-mode workflow ❌ (readonly toolset ADA, belum
jadi alur) · checkpoint/rewind ❌ · subagent isolasi konteks ❌ · **task ledger in-memory saja
(hilang saat restart)** ❌ · MCP tool schema masih di-inject penuh ❌ · memori masih KV datar ❌ ·
hooks ❌ · auto-classifier ❌ · eval harness belum dikodifikasi ❌.

---

## P0 — TRUST & SAFETY (milestone v2.1)
*Prioritas tertinggi: keluhan #1 komunitas = keamanan/review burden; YOLO tanpa sandbox = "opt-in rootkit".*

### 1. Sandbox eksekusi shell (effort L, impact H)
**Bukti**: Anthropic "sandboxing → -84% permission prompts"; dual isolation (fs+network) adalah
jawaban industri. **Desain**: di `tools/exec.go ExecuteShell`, bungkus `exec.Command` dengan
bubblewrap (`bwrap --ro-bind / / --bind $PWD $PWD --tmpfs /tmp --unshare-all --share-net --die-with-parent`)
saat tersedia; varian network-deny default + allowlist domain via env. Daemon VPS **harus pindah
ke user non-root** (sekarang root — red flag terbesar scorp). Integrasi: YOLO tanpa bwrap →
peringatan permanen di status; supervised default = sandbox on; escape hatch per-command via
confirm gate yang sudah ada.

### 2. Deny-rule engine (effort M, impact H)
**Bukti**: deny rules berlaku bahkan di `bypassPermissions` (Claude Code docs). **Desain**: generalisasi
`config.IsPathRestricted` → `config.DenyRules` dengan pola `tool(param:regex)` (mis. `shell(*:curl*|*nc*)`,
`write_file(*:/etc/*)`), load dari config; dievaluasi di `ExecuteTool` SEBELUM `ConfirmationRequired`
— sehingga tetap hidup di YOLO. Regression test per rule.

### 3. Test-integrity gate (effort M, impact H)
**Bukti**: benchmark nyata 19/45 klaim "all tests pass" palsu; METR 50% PR lulus benchmark tak
di-merge; agen melemahkan test. **Desain**: saat sesi menyentuh file test/`conftest`/CI config,
`complete_task` ditolak kecuali ada receipt shell berisi eksekusi test suite hijau *setelah* edit
terakhir (receipts.json sudah hash output — tinggal dipakai sebagai bukti). Simetri dengan
plan-completion gate yang sudah ada di kedua loop.

---

## P1 — PLAN & CHECKPOINT (milestone v2.2)
*Prioritas: plan mode = fitur dipuji #1; checkpoint = #5; "state di file" = konsensus #1.*

### 4. Plan mode workflow (effort M, impact H)
**Desain**: `/plan <goal>` → loop dijalankan dengan `IsToolAllowed` readonly (SUDAH ADA di
`config/autonomy.go`) + task ledger dibuat → plan dirender via inline keyboard Telegram
("✅ Approve plan" / "✏️ Revise" / "❌ Cancel") → approve = lanjut ke mode eksekusi dengan ledger
sama (auto-resume & plan gate sudah mendukung). Zero infrastruktur baru — komposisi dari 3 fitur
yang sudah terbukti di verify26.

### 5. Persistensi task ledger (effort S, impact H)
**Desain**: `taskPlans` di-flush ke `<session>.plan.json` tiap update (pola sama dengan history
save yang sudah ada); load saat sesi pertama kali dipakai; daemon restart tidak lagi membuang plan.
Termasuk migrasi: plan selesai dihapus dari disk.

### 6. Checkpoint/rewind (effort M/L, impact H)
**Bukti**: "FINALLY checkpoints!" (HN); Cursor & Claude Code jadikan baseline UX. **Desain**: per
turn dengan perubahan file → `git stash create` / commit bayangan di `refs/scorp/ckpt` (repo git);
non-git → snapshot tar ke `~/.scorp/checkpoints/<session>/`. `/undo` = restore snapshot terakhir
(code-first). Simpan maksimal N=20 per sesi. Pasangkan dengan receipt turn untuk audit.

### 7. Auto memory protocol (effort S/M, impact M/H)
**Bukti**: "agent lupa lintas sesi" = keluhan permanen; auto memory = jawaban resmi Anthropic.
**Desain**: upgrade `memory.json` KV → `MEMORY.md` per proyek + index; di `complete_task`/session
end, agent menulis ringkasan keputusan & state (prompt khusus, bukan heuristik `extractAndSaveMemory`
yang sekarang); di sesi start, MEMORY.md di-inject setelah system prompt; kuota ~200 baris
(konsensus docs CC).

---

## P2 — EKONOMI TOKEN & KONTEKS (milestone v2.3)
*Bukti: overhead harness = keluhan #2; Claude Code 33k token sebelum prompt vs OpenCode 7k.*

### 8. Subagent `delegate` (effort L, impact H)
**Desain**: tool baru `delegate(task, max_turns)` → spawn child loop dengan context FRESH (system
prompt minimal + task), hanya laporan final yang kembali ke parent sebagai tool result. Use case:
eksplorasi/read besar, fan-out research. Cap turn & wall-clock (pola `maxAutoResumes` yang sudah
ada). INI penghemat token terbesar: N×(N+1)/2×S akumulasi re-read parent terpotong jadi 1 summary.

### 9. MCP schema deferred-by-default (effort S/M, impact M/H)
**Bukti**: CTO Perplexity membuang MCP karena 15-20k token skema; best practice cap 10–15 tool.
**Desain**: infrastruktur `Deferred` registry SUDAH ADA — tandai semua tool MCP `deferred:true`,
aktif via `tool_search` (TTL mekanisme sudah jalan). Merchant marketplace tetap tampil, skema tak
membebani context.

### 10. Compaction preservation (effort S, impact M)
**Desain**: `maybeCompactHistory` menerima preservation instructions: task ledger aktif, goal user
terakhir, dan laporan final task sebelumnya SELALU selamat compaction (persis pembelajaran
Roll-Up/verify26). Tampilkan "🗜 compacted at N%" di thinking footer Telegram.

### 11. `/cost` per sesi (effort S, impact M)
**Desain**: `model_usage.json` sudah mencatat — tambah agregasi per sesi/turn + tampilkan di
`/usage`; warning otomatis saat burn melewati ambang per task.

---

## P3 — EXTENSIBILITY & AUTONOMI UX (milestone v3.0)

### 12. Hooks PreToolUse/PostToolUse (effort M, impact M/H)
**Bukti**: "CLAUDE.md bilang 'tolong', hooks bilang 'harus'" — enforcement deterministik dipuji
enterprise. **Desain**: config `hooks: {pre_tool_use: [...]}`; dieksekusi di titik gerbang
konfirmasi di kedua loop (satu titik: `ExecuteTool`); exit code 2 = blok, stdout = additional
context ke model. Non-blocking untuk audit/logging.

### 13. Auto-mode classifier (effort M/L, impact H)
**Bukti**: data Anthropic — manusia menangkap 13.6% perintah berbahaya vs classifier 89%; auto mode
jadi default Agu 2026. **Desain**: tingkat ke-4 `auto` di antara supervised↔yolo: model murah
mengklasifikasi tiap tool call (safe → jalan, risky → confirm gate yang sudah ada, destructive →
hard-deny kecuali allowlist); fallback manual setelah N keputusan meragukan; statistik keputusan ke
receipts. `ConfirmationRequired()` diperluas jadi `PermissionDecision(tool, args) -> allow|ask|deny`.

### 14. MCP contract watch (effort S, impact M)
**Desain**: fingerprint (hash nama+skema tool) per server MCP; perubahan diam-diam → warning
(pola mcpwatch: "uptime bilang server menjawab, tidak bilang ia masih melakukan yang agen harapkan").

---

## P4 — VERIFIKASI & KUALITAS (berjalan sejak v2.1)

### 15. `scorp eval` — arena privat (effort M, impact H)
**Bukti**: SWE-bench di-retire; konsensus "build your own private arena dari backlog nyata".
**Desain**: kodifikasi metode verify26 (yang sudah terbukti) jadi suite: 20–40 task berkategori
(persistence, safety-gate, plan, MCP, sandbox) + pemeriksa artefak independen; dijalankan otomatis
sebelum deploy; metrik: pass rate, klaim-vs-verifikasi delta, token/task.

### 16. Bukti uji wajib untuk klaim (effort S, impact M)
**Desain**: perluasan anti-fabrication gate — klaim "all tests pass" tanpa receipt eksekusi test
pada sesi berjalan → ditolak dengan nudge. Receipts.json tinggal di-query.

---

## EKSPILISIT TIDAK DIBANGUN (sekarang)
- **Agent teams/swarms** — bukti: "significantly more tokens", koordinasi mahal; tunggu `delegate` matang.
- **Repo map ala Aider** — parkir; pertimbangkan saat pemakaian multi-repo naik.
- **Voice / IDE plugin / cloud agents** — bukan medan scorp (Telegram-first sudah jadi pembeda).

## URUTAN EKSEKUSI (2 minggu pertama — quick wins ber-impact)
1. P0.2 deny-rule engine (M) → 2. P1.5 ledger persisten (S) → 3. P0.3 test-integrity gate (M) →
4. P2.9 MCP deferred (S/M) → 5. P1.4 plan mode (M) → 6. P2.10 compaction preservation (S) →
lalu P0.1 sandbox (L) sebagai proyek tersendiri, P1.6 checkpoint, sisanya menyusul.

## METRIK SUKSES (merge-rate mindset)
- Permission prompt per task ↓ (baseline setelah P0.1)
- Token per task ↓ ≥40% (P2.8/9 — ukur via /cost)
- Eval pass rate ≥95% + delta klaim-vs-verifikasi → 0 (P4)
- Zero kehilangan plan saat restart (P1.5); zero insiden "compaction membuang goal" (P2.10)
- Insiden keamanan: 0 eksekusi di luar sandbox sejak P0.1
