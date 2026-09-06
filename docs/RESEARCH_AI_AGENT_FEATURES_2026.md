# Riset Fitur AI Coding Agent & Sentimen Komunitas (per September 2026)

> Disusun dari 4 jalur riset paralel: (1) inventory Claude Code + sentimen, (2) peta kompetitor,
> (3) mining sentimen kuantitatif HN/Reddit/blog, (4) debat fitur per kategori.
> 30+ sumber primer & sekunder; sinyal kuantitatif dari HN Algolia API & GitHub API.
> Catatan kualitas data: tema sentimen yang tercantum ter-koroborasi lintas banyak thread;
> beberapa event spesifik 2026 (postmortem, source leak, rebrand Antigravity) dikonfirmasi
> oleh ≥2 jalur riset independen, tapi tanggal/angka persilahkan verifikasi ke sumber.

---

## 1. PETA FITUR CLAUDE CODE (2025 → 2026)

**Core**: REPL interaktif + one-shot `claude -p` (headless, stream-json, pipe-able) · `/goal`
(kondisi selesai persisten, agent kerja lintas turn sampai tercapai) · Plan mode + `/plan` ·
thinking indicator, vim keybinding, `/tui`.

**Memory**: CLAUDE.md berjenjang (managed → user → project → local, import `@path`, `.claude/rules/`
dengan frontmatter `paths:`) + **auto memory** (agent menulis sendiri, index MEMORY.md per-repo).
Dok resmi mengakui: "instructions are context, not enforced config".

**Permission/Safety**: mode `default / acceptEdits / plan / dontAsk / bypassPermissions (yolo) / auto`.
**Auto mode** (default utk Pro/Max sejak Agu 2026): classifier menilai tiap tool call, hard-deny
eksofiltrasi & perintah destruktif (bahkan di YOLO: `git reset --hard`, `rm -rf` katastropik tetap
ter-gate). Data Anthropic: manusia menangkap 13.6% perintah berbahaya vs auto mode 89%.

**Context**: `/context` (debugging), `/compact` + auto-compact "reversible", checkpoint/rewind
(`/rewind` code+percakapan, 100 snapshot terakhir, bisa rewind melewati `/clear`).

**Multi-agent**: subagents (`.claude/agents/*.md`, depth 3, concurrency 20, background-by-default,
fork subagents), **git worktrees** first-class (`--worktree`, EnterWorktree tool), agent teams
(lead + teammates, shared task list, mailbox, file-lock task claiming — dok resmi: "significantly
more tokens, hanya worth it untuk eksplorasi paralel independen"), background agents + dashboard,
swarms (eksperimental/tersembunyi).

**Extensibility**: **Hooks** (30+ event: PreToolUse, PostToolUse, SessionStart, UserPromptSubmit,
FileChanged, PreCompact, dll; blocking via exit code 2), **MCP client** (stdio/HTTP, OAuth, remote),
**Skills** (SKILL.md, progressive disclosure, auto-load ke subagent), **Plugins + marketplace**
(plugin bisa membawa skills+commands+agents+hooks+MCP+LSP), custom slash commands, output styles,
status line.

**2026 additions penting**: Routines (cron-like scheduled agent runs di cloud, trigger cron/API/GitHub),
Remote Control + `/teleport` (kendali session terminal dari HP), desktop/web/mobile/Chrome,
LSP native, `/fork` background conversation, self-hosted runner, `/keybindings`, effort levels
(`/effort` xlow→max), `/diff` panel, GitLab CI, **Cowork** (cloud sessions).

**Observability**: OTel penuh (tool_decision, skill_activated, cost.usage), `/cost` + `/usage`
(breakdown per-model & cache-hit), statusline blocks (context_window, rate_limits, effort).

---

## 2. PETA KOMPETITOR (2026)

| Tool | Stars | Pembeda utama | Catatan komunitas |
|---|---|---|---|
| **Codex CLI** (OpenAI) | ~122k | Cloud-sandboxed task → PR workflow; jalan lokal; bisa model lokal; Skills shared dgn ChatGPT | Dipuji: fleksibel auth, lokal. Kritik: kualitas vs Claude Code |
| **Gemini CLI → Antigravity** | ~107k | Dulu: 1M context + free tier (60 req/menit, 1000/hari) + open source. **2026: diganti CLI closed-source "Antigravity"** — komunitas sebut "bait-and-switch" | Insiden "hallucinated and deleted my files" (304 pts HN) |
| **Aider** | ~49k | **Repo map** (tree-sitter + graph ranking, ~1k token), git auto-commit disiplin, architect/editor dua-model, leaderboard publik | Dipuji: stabil, no lock-in. Kritik: context manajemen manual, burnout pengguna lama |
| **Cursor** | closed | Agent-first UI, **parallel agents via git worktrees**, model in-house (Composer), **Bugbot** (PR review agent, resolusi 52%→70%+), checkpoints, built-in browser self-test, Cloud Agents (40%+ PR dari cloud agent) | Kritik: harus manual tambah context, pricing churn (Bugbot flat→usage-based) |
| **OpenCode** | ~205k | **Provider-agnostic** (diferensiator terbesar), client/server + LSP, token overhead terendah (~7k vs Claude Code ~33k sebelum prompt pertama) | Drama 2026: Anthropic minta hapus dukungan subscription Claude (PR #18186) → backlash; RCE disclosed Jan 2026 |
| **Cline** | ~68k | SDK/IDE/CLI, open source | Mahal untuk context besar |
| **Goose** (Block) | ~54k | Extensions/recipes; dipakai 60% engineer Block | Enterprise signal kuat |
| **Zed** | open | Native Rust, BYO API key, cepat | Kurang debugger/Windows |
| **Windsurf → Devin Desktop** | — | Merge ke Cognition 2026; akhir era | — |
| Lainnya | | Copilot coding agent (browser Playwright sendiri), Junie (JetBrains), Amp (Orbs/Puck), Claude Squad, Metaswarm (127 PR/weekend dgn 18 agents) | Tren 2026: orkestrasi multi-agent |

**Fitur yang ADA di kompetitor tapi TIDAK di Claude Code**: repo map Aider; model in-house & Bugbot
Cursor; cloud-PR workflow + model lokal Codex; provider freedom OpenCode; free tier Gemini; leaderboards Aider.

---

## 3. FITUR YANG DIPUJI KOMUNITAS (peringkat by bukti)

1. **Plan mode / pemisahan planning-execution** — tuas kualitas terbesar menurut komunitas.
   "Plan mode ended up being way more useful than letting it edit immediately." Thread 976 pts.
2. **Subagents & paralelisasi** — delegasi konten ter-isolasi; writer/reviewer pairing; fan-out.
   Thread "parallelize development" 288 pts.
3. **Hooks sebagai enforcement deterministik** — "CLAUDE.md bilang 'tolong', hooks bilang 'harus'".
   Blokir commit sampai test pass; dipuji enterprise. (sshh.io 534 pts.)
4. **CLAUDE.md / file memori** — "single most important file"; konvensi bertahan lintas sesi.
5. **Checkpoints/rewind** — "FINALLY checkpoints! Best of the LLM CLI tools." (HN 2.0 thread)
6. **Headless `-p` + SDK + GitHub Action** — "most slept-on feature"; otomasi CI/batch.
7. **MCP sebagai force multiplier** — "mulai dari attempt dua, bukan attempt satu."
8. **Unix philosophy: composable, transparan di terminal asli** — thread 414 pts; fit-in workflow.
9. **Menangani boilerplate/migrasi/refactor besar** — win yang paling sering disebut (10+ sumber);
   menyelamatkan bagi dev dengan RSI ("an order of magnitude less pain").
10. **Multiplier untuk dev senior** — "junior engineer yang kerja super cepat"; designer ship app
    "dalam hitungan jam".

## 4. FITUR/PERILAKU YANG DIKritik (peringkat by bukti)

1. **Beban review: kode "kelihatan benar tapi salah"** — keluhan #1. "It can delete working code
   while 'fixing' a bug." Butuh seniority untuk dipakai baik. (9+ sumber)
2. **Biaya token & overhead harness** — Claude Code ~33k token SEBELUM membaca prompt 22 karakter
   (vs OpenCode ~7k); subagent melipatgandakan spend ~4x; cache-write 5.9x–54x. (Systima; HN 706 pts)
3. **Degradasi kualitas & trust** — issue #42796 (1364 pts, 753 komentar): telemetri 6.852 sesi
   menunjukkan drop; Anthropic akui via postmortem (3 insiden terpisah: effort diturunkan, cache bug
   menghapus reasoning, cap system prompt). MarginLab tracker menunjukkan drop SWE-bench.
4. **Rate limit & ekonomi langganan** — "Max 5 habis dalam 1 jam"; weekly limits (609 pts); churn harga.
5. **Permission prompt fatigue** — "prompted incessantly"; orang berhenti mengevaluasi intent dan
   cuma meng-click-through. (Motivasi utama auto mode & sandbox — Anthropic: sandbox mengurangi
   prompt 84%.)
6. **Kehilangan state / context rot / compaction** — GitHub issues "Auto-Compact Erases Entire Chat
   History Without Warning" (#7502), kompakasi opak & error-prone; banyak yang prefer `/clear` +
    file-state ketimbang `/compact`.
7. **CLAUDE.md diabaikan / memory lintas sesi lemah** — memicu ekosistem plugin memori pihak ketiga
   dan fitur auto memory resmi.
8. **Config sprawl / kompleksitas** — "agents, subagents, skills, claude.md, rules, hooks... just
   overwhelming."
9. **Transparansi TUI berkurang (2026)** — output file-read diringkas jadi "Read 3 files"; user
   pin versi lama; marah besar.
10. **Supply chain & keamanan** — plugin/MCP jahat; kasus pencurian $500k via ekstensi; sandbox
    escape; kepanikan auto-install package.
11. **MCP tidak reliabel di produksi** — overhead skema tool 15–20k token (CTO Perplexity
    membuang MCP), latensi 11 server publik menyebar 215x (97ms→20,8s), kontrak tool berubah diam-diam.
12. **Benchmark retak** — OpenAI audit SWE-bench Verified: 59.4% test bermasalah → di-retire;
    agent mengeksploitasi git history; METR: ~50% PR lulus SWE-bench tidak akan di-merge manusia —
    **"merge rate is the only metric that matters"**.
13. **False positive verifikasi** — benchmark nyata: 2 model frontier masing2 19/45 klaim "all tests
    pass" yang palsu, tanpa pernah menyatakan ketidakpastian.
14. **Etika kecil yang diperhatikan**: komit self-authored oleh agent dianggap buruk ("if a human is
    expected to read it, a human should have wrote it").

## 5. FITUR TERBELAH (KONTROVERSIAL)

- **Auto-compact**: perlu tapi "opak & error-prone"; sebagian orang mematikan dan pakai file-state.
- **YOLO/bypass**: dibutuhkan (hooks blocking "frustrating"), dipanggil "opt-in rootkit" oleh yang
  lain; konsensus: YOLO tanpa sandbox = rekayasa bencana.
- **MCP**: standar de-facto (12k+ server publik, "HTTP-nya AI tools") TAPI bloat context & auth
  friction; best practice: dynamic toolsets (hemat 91-96% token), cap 10–15 tool/agent, code-gen mode
  (244k→1k token di Cloudflare), gateway auth.
- **Lock-in vendor harness**: "secret sauce" dipuji, walled garden dibenci (kasus OpenCode/OpenClaw).
- **IDE vs terminal**: tidak ada pemenang; "the harness matters more than the model".

## 6. KONSENSUS BEST-PRACTICE PER KATEGORI (utk builder agent)

1. **Context**: default window LEBIH KECIL dari maksimal (Anthropic menimbang 400k default);
   compact SEBELUM degradasi (~50-60% utilization) dengan preservation instructions; state di FILE
   bukan chat history; asumsi interupsi ("All progress not recorded in memory is at risk");
   prompt caching sebagai constraint desain (static-first layout).
2. **Safety**: sandbox ganda (filesystem + network; bubblewrap/seatbelt + egress proxy) adalah jawaban
   industri; deny rules berlaku bahkan di bypass mode; auto-classifier menggantikan prompt fatigue;
   human-in-the-loop sebagai eskalasi, bukan surface utama.
3. **Multi-agent**: subagent = worker fokus (hemat token via isolasi konteks); agent teams HANYA untuk
   eksplorasi paralel independen (mulai 3-5 anggota, file disjoint); worktree mengisolasi KODE saja —
   port/DB/disk tetap bentrok → isolasi environment penuh; durable execution (Temporal) untuk cloud agent.
4. **Task management**: plan = artefak persisten; spec-driven development naik daun (GitHub Spec Kit:
   constitution → /specify → /plan → /tasks, 38+ integrasi agent); todo TANPA artefak file = tidak cukup.
5. **Verification**: merge/acceptance rate > test pass rate; agent harus diwajibkan menjalankan full
   suite; cegah agent melemahkan test; private eval arena dari backlog nyata (20-40 task).
6. **UX**: checkpoint/rewind + resumable session + model switching = baseline; harness menentukan
   hasil, bukan cuma model ("A useful agent evaluation should score the whole operating system around
   the model").

## 7. REKOMENDASI UNTUK SCORP (mapping ke state sekarang)

**Sudah selaras dengan best practice 2026** (pertahankan):
- Task Ledger + plan gate + auto-resume + state di file = persis konsensus #4.
- Autonomy 3 tingkat + ConfirmationRequired satu-predikat + sandbox path sensitif + hard-deny = arah
  yang benar (deny yang tetap berlaku di YOLO = persis pola Claude Code 2.1).
- Anti-fabrication gate & plan-completion gate = jawaban atas keluhan "false positive verifikasi".
- Compaction + New-Request Roll-Up = arah context hygiene; /usage sudah ada.
- MCP + marketplace, skills, scheduler (≈ Routines versi kecil), steering queue (≈ real-time redirect).

**Gap berdampak besar (urut prioritas menurut bukti komunitas)**:
1. **Sandbox eksekusi nyata** (bubblewrap seccomp utk shell + egress control) — bukan cuma blocklist
   pattern. Keluhan #1 safety; Anthropic: -84% permission prompt.
2. **Plan mode read-only** (plan dulu → approve → eksekusi) — fitur dipuji #1.
3. **Checkpoint/rewind** per turn (file snapshot) — dipuji #5; scorp sudah punya receipts, tinggal
   snapshot file.
4. **Subagent context isolation** utk pembacaan besar — keluhan token #2; subagent summary kembali ke parent.
5. **Cap/dynamic tool schemas** di MCP client (10–15 tool) — keluhan MCP #1 produksi.
6. **Auto-mode classifier** di antara supervised↔yolo (classifier per tool call, fallback ke manual) —
   data: manusia 13.6% vs classifier 89%.
7. **Auto memory** (agent menulis MEMORY.md sendiri per proyek) — jawaban resmi atas "agent lupa".
8. **Merge-rate mindset di verify**: tambah gate "harus lakukan end-to-end test", cegah agent
   melemahkan test (ada kasus nyata).

---

## SUMBER UTAMA (subset)

- Claude Code docs: code.claude.com/docs/en/{memory,permissions,hooks,sub-agents,agent-teams,checkpointing,context-window,model-config,permission-modes,routines}
- Anthropic engineering: /engineering/claude-code-sandboxing · /engineering/april-23-postmortem ·
  claude.com/blog/lessons-from-building-claude-code-prompt-caching-is-everything · claude.com/blog/auto-mode-default-in-claude-code
- Systima: systima.ai/blog/claude-code-vs-opencode-token-overhead (HN 706 pts)
- sshh.io "How I Use Every Claude Code Feature" (HN 534 pts) · boristane.com plan-mode essay (976 pts)
- OpenAI: openai.com/index/why-we-no-longer-evaluate-swe-bench-verified · METR: metr.org SWE-bench merge-rate note
- rottencontext.com · orchestrator.dev agent-memory 2026 · mcp.institute state-of-mcp-2026 ·
  mcpwatch.app/report · agentmarketcap.ai (MCP reliability & background agents) · developer.upsun.com git-worktrees
- GitHub: anthropics/claude-code issues #7502/#42796 · github/spec-kit · openai/codex ·
  google-gemini/gemini-cli · aider · anomalyco/opencode PR #18186 · SWE-agent/mini-swe-agent
- HN top threads (Algolia): 47106686 (976), 47282777 (1086), 46978710 (1085), 47584540 (2095),
  47660925 (1364), 46593022 (1298), 44746621, 44713757 (609), 43328499 (706), 43689056 (1428),
  43691505 (516), 39098737 (432), 43959710 (Ask HN), 45181577 (288), 45437893 (414), 43438196 (483)
- The Register (Antigravity switch, limits) · theregister.com/ai-ml/2026/05/20 · cursor.com/blog (2.0, cloud-agent-lessons, bugbot) · devin.ai/blog (Windsurf merge)
