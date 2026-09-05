# 🦂 Scorp MCP Marketplace — Build Plan & Execution Roadmap

> **Status:** Executed — M0/M1/M2 selesai & tervalidasi; M3 sebagian (container hardening ditunda)
> **Blueprint Reference:** [`docs/MCP_MARKETPLACE_BLUEPRINT.md`](MCP_MARKETPLACE_BLUEPRINT.md)
> **Target Version:** Scorp Agent v2.5 / v3.0
> **Owner:** Wahyu
> **Created:** September 2026

---

## 🎯 1. Ringkasan Eksekutif

Dokumen ini menerjemahkan blueprint MCP Marketplace menjadi rencana build yang dapat dieksekusi, berdasarkan audit aktual codebase scorp-agent.

Temuan sentral dari audit: **blueprint ini menumpang pada infrastruktur yang sudah ada, bukan membangun dari nol.** Ketiga opsi instalasi (prebuilt / rebuild / upstream) berakhir di titik yang sama — sebuah entry di `~/.scorp/mcp.json`. Setelah entry itu ada, seluruh infrastruktur existing otomatis bekerja: watchdog, hot-reload, health status, dan registrasi tool sebagai native function di `registry/`. Dengan demikian, Tri-Option adalah **lapisan sourcing baru di depan `mcp_manage` yang sudah ada**, bukan sistem baru.

**Strategi eksekusi yang disepakati:**

1. **Vertical slice dulu** — `/mcp search` + `/mcp install` dengan Tri-Option untuk port seed (prebuilt + upstream) sebelum transpiler. Marketplace terasa hidup sejak M1 dengan risiko rendah.
2. **Port referensi ditulis manual** — `mcp-fetch`, `mcp-sqlite`, `mcp-filesystem` dibuat idiomatic sebagai golden standard, sekaligus isi marketplace pertama dan test case validasi transpiler.
3. **Sandbox build ringan** — `go build` di temp directory terisolasi dengan verifikasi checksum; hardening container menyusul.
4. **Codegen self-hosted** — transpiler memakai LLM gateway internal scorp (routing model *complex/premium* dari `models.json`), tanpa dependency CLI eksternal.

---

## 🔍 2. Audit Codebase: Existing vs Gap

### 2.1 Yang Sudah Ada (Tidak Perlu Dibangun)

| Komponen Blueprint | Implementasi Existing | Lokasi |
| :--- | :--- | :--- |
| MCP client (JSON-RPC 2.0, stdio + SSE) | Hand-rolled client, spawn server, hot-reload | `mcp/client.go`, `mcp/sse.go` |
| Server lifecycle & crash recovery | Watchdog + health status | `mcp/watchdog.go`, `mcp/RestartServer` |
| Tool registration ke LLM | Native function calling + deferred tools + raw JSON Schema | `registry/registry.go` |
| Manajemen server (add/remove/list/reload) | Agent tool `mcp_manage` dengan hot-reload | `mcp/manage.go` |
| **Layer 5: Secret Redactor** ✅ | Regex sanitizer output tool | `tools/redact.go` (+ `redact_test.go`) |
| Engine codegen transpiler | Multi-provider gateway, routing rules, fallback on error | `gateway/gateway.go`, `models.json` |
| Interface UX | `/mcp` command CLI + tombol MCP Telegram (`callback_data`) | `cli.go:278`, `telegram/telegram.go` |
| Path konfigurasi terpusat | `mcpConfigFilePath()` → `~/.scorp/mcp.json` | `config/config_paths.go:93` |

### 2.2 Gap yang Harus Dibangun

- ❌ Registry client: fetch + cache `registry.json`, search lokal, parser manifest
- ❌ JSON Schema `scorp-mcp.json` v1 (health matrix, flavor, contributors, origin)
- ❌ Tri-Option install flow (CLI interaktif + Telegram inline keyboard)
- ❌ SHA-256 pinning & verifikasi artifact (Cosign menyusul)
- ❌ Pipeline transpiler: probe → codegen → sandbox build → contract verify
- ❌ Drift livecheck & indikator sinkronisasi upstream
- ❌ PR generator (kontribusi port / maintenance update)
- ❌ Repo publik `scorp-mcp-registry` + 3 workflow CI
- ❌ 3 port referensi Go (fetch, sqlite, filesystem)

---

## 🏛️ 3. Insight Arsitektur Kunci

### 3.1 Titik Konvergensi Tri-Option

```
Opsi 1 Prebuilt  ──┐  download binary → SHA-256 verify → ~/.scorp/mcp-binaries/
                   ├──>  Entry di ~/.scorp/mcp.json  ──>  infra existing jalan sendiri
Opsi 2 Rebuild   ──┤  (transpile → go build → binary lokal)     (watchdog, hot-reload,
                   │                                             native tool registry)
Opsi 3 Upstream  ──┘  register npx/uvx langsung (sudah didukung client)
```

Implikasi: `mcp_manage(action="add")` yang sudah ada adalah mekanisme akhir dari SEMUA jalur instalasi. Komponen baru hanya memutuskan *cara menghasilkan* entry tersebut.

### 3.2 Transpiler = Self-Hosted Codegen

- **Phase 1 (Probe):** reuse JSON-RPC client `mcp/client.go` untuk spawn server ephemeral (deteksi runtime `npx`/`uvx`/`python3`) dan capture `initialize` + `tools/list` sebagai benchmark kontrak.
- **Phase 2 (Generate):** prompt codegen khusus `mark3labs/mcp-go`, dieksekusi via `gateway/` dengan routing model *complex* (premium).
- **Phase 3 (Verify):** `go build` di temp dir, boot binary hasil build, diff skema tool terhadap benchmark Phase 1. Gagal → graceful degradation ke upstream runtime.

### 3.3 Security Layers: Status Awal

| Layer | Status | Catatan |
| :--- | :--- | :--- |
| 1. Source-Only Registry | ⬜ Baru (M0) | Kebijakan repo registry + CI |
| 2. AST & Prompt Scan | ⬜ Baru (M0) | `gosec` + `govulncheck` + prompt linter di CI |
| 3. Hermetic CI/CD | ⬜ Baru (M0) | GitHub Actions, multi-arch |
| 4. SHA-256 Pinning | ⬜ Baru (M1) | Verifikasi di sisi client sebelum eksekusi |
| 5. Secret Redactor | ✅ **Sudah Ada** | Tinggal audit coverage path output MCP client |

---

## ⚠️ 4. Keputusan Terbuka & Inkonsistensi Blueprint

1. **Lokasi source port (inkonsistensi blueprint):** Bagian 5C blueprint menaruh `main.go` langsung di repo registry, tetapi bagian 5B manifest menunjuk `port_repository` ke repo terpisah `scorp-mcp-ports`.
   **Keputusan:** Source port hidup langsung di `scorp-mcp-registry/servers/<nama>/main.go` untuk MVP (satu sumber kebenaran, satu tempat CI audit). Field `port_repository` menjadi opsional.
2. **Cosign:** ditunda ke M3. MVP mengandalkan SHA-256 murni; key management Cosign adalah proyek kecil tersendiri.
3. **Cleanup kecil:** `mcp/manage.go` masih hardcode `$HOME/.scorp/mcp.json` — migrasikan ke `config.ConfigFilePathMCP()` saat menyentuh file tersebut.

---

## 🚀 5. Milestone & Deliverables

### M0 — Foundation: Schema, Registry Repo & Port Referensi

*Scope: tanpa menyentuh codebase scorp-agent.*

- [x] Finalisasi JSON Schema `scorp-mcp.json` v1: field `origin` (port/native), `upstream.pinned_commit`, `health` (status, coverage_score, active/unsupported tools), `variant` (flavor), `contributors`, `artifacts` (multi-arch + sha256), `build` (go_version, sdk)
- [x] Publish schema ke `scorp-mcp-registry/schema/v1.json`
- [x] Scaffold repo publik `wahyuzero/scorp-mcp-registry` dengan layout blueprint 5C (`registry.json`, `servers/`, `.github/workflows/`)
- [x] Workflow CI: `security-audit.yml` (gosec, govulncheck, prompt-injection linter, JSON-RPC contract validation), `release-binaries.yml` (cross-compile linux amd64/arm64 + checksum), `drift-livecheck.yml` (stub, diaktifkan penuh di M3)
- [x] Tulis 3 port referensi manual (idiomatic, `mark3labs/mcp-go`, Go 1.24+):
  - [x] `mcp-fetch` — web scraping & markdown extraction
  - [x] `mcp-sqlite` — inspeksi & query database lokal
  - [x] `mcp-filesystem` — akses file lokal ter-scope
- [x] Manifest `scorp-mcp.json` untuk ketiganya (artifacts diisi oleh release CI saat pertama tag)

**Acceptance:** `registry.json` ter-index; ketiga binary lolos security-audit CI dan bisa dijalankan manual via JSON-RPC handshake.

### M1 — Vertical Slice: Marketplace Client & Tri-Option Install (MVP)

*Scope: package baru `mcp/marketplace/` + integrasi CLI & Telegram.*

- [x] `mcp/marketplace/registry.go` — fetch `registry.json` dari GitHub raw, cache lokal dengan ETag/TTL, fallback ke cache stale saat offline
- [x] `mcp/marketplace/manifest.go` — parser + validasi manifest terhadap schema v1
- [x] `mcp/marketplace/search.go` — pencarian lokal (nama, deskripsi, tool, flavor)
- [x] `mcp/marketplace/install.go` — orkestrasi Tri-Option:
  - [x] Opsi 1 (Prebuilt): unduh artifact per-arsitektur, **verifikasi SHA-256 sebelum eksekusi**, simpan ke `~/.scorp/mcp-binaries/`, register via `mcp_manage add`
  - [x] Opsi 2 (Rebuild): tampil dengan status "coming in v2.5 — jalankan transpiler manual" (diaktifkan penuh di M2)
  - [x] Opsi 3 (Upstream): register `npx`/`uvx` langsung — reuse jalur existing
- [x] Pre-install disclosure: health badge (🟢 full / 🟡 partial), daftar tool nonaktif + alasan, delta resource vs upstream
- [x] CLI: subcommand `/mcp search <term>` dan `/mcp install <target>` di `cli_mcp.go` (prompt interaktif `[1/2/3]`)
- [x] Telegram: inline keyboard Tri-Option (pola `callback_data` existing), dialog naratif sesuai blueprint bagian 2
- [x] Audit coverage Layer 5: output tool MCP kini melewati `tools/redact.go` di `registerMCPToolsAsNative`

**Acceptance:** End-to-end `/mcp install fetch` → pilih opsi 1 → binary terverifikasi terdaftar dan tool-nya muncul di agent; opsi 3 berfungsi untuk server Node/Python apa pun; dialog Telegram setara CLI.

### M2 — AI Transpiler

- [x] `mcp/transpiler/probe.go` — spawn ephemeral upstream server di sandbox dir, deteksi runtime (`npx`/`uvx`/`python3` dengan graceful fail + pesan jelas bila tidak ada), capture benchmark: skema tool, parameter, required fields
- [x] `mcp/transpiler/generate.go` — prompt codegen terstruktur untuk `mark3labs/mcp-go`, dieksekusi via `gateway/` (routing *complex*), output `main.go` tunggal yang decoupled
- [x] `mcp/transpiler/build.go` — `go build` di temp dir terisolasi; pinning `go.sum`; module cache terkontrol
- [x] `mcp/transpiler/verify.go` — boot binary hasil build, kirim mock JSON-RPC, diff skema output vs benchmark Phase 1 (target: 100% match)
- [x] Graceful degradation: kegagalan compile/verify → laporan blocker transparan → tawaran fallback ke upstream runtime
- [x] Aktifkan Opsi 2 di Tri-Option (CLI + Telegram) via `marketplace.RebuildHook`
- [x] Flow "Share to Marketplace": generator manifest + draft PR (dengan atribusi penulis asli) via `gh` (`/mcp share <name>`)
- [x] **Validasi golden standard:** jalankan transpiler terhadap upstream fetch/sqlite/filesystem, bandingkan hasil dengan port manual M0

**Acceptance:** Minimal 1 dari 3 port upstream berhasil ditranspilasi end-to-end dengan contract match 100%; kegagalan selalu menghasilkan jalur fallback yang jelas.

### M3 — Governance & Hardening

- [x] `drift-livecheck.yml` penuh: cek harian upstream tags/commits, update field drift di `registry.json`, indikator 🟢/🟡 di hasil search/install + dialog notifikasi resync
- [x] Cosign step disiapkan (disabled, butuh key management) — SHA-256 tetap Layer 4 utama untuk sekarang
- [ ] Hardening sandbox build: evaluasi container opsional (Docker bila tersedia) untuk isolasi network/FS
- [x] Workflow kolaborasi multi-maintainer: setiap PR update melewati Layer 2 penuh (`security-audit.yml` on pull_request); attribution changelog di manifest
- [x] Dukungan flavor/variant di search (namespace `@author/nama-flavor`) — `variant` ikut terindeks di registry.json
- [x] Klasifikasi `origin: native` untuk submission Go orisinal (tanpa drift check)

---

## 📊 6. Risiko & Mitigasi

| Risiko | Dampak | Mitigasi |
| :--- | :--- | :--- |
| Kualitas codegen LLM variatif | Transpiler menghasilkan kode yang gagal compile / kontrak meleset | Port manual M0 sebagai golden test case; verify.go mewajibkan contract diff 100%; graceful degradation ke upstream |
| Runtime upstream tidak terpasang di mesin user (`npx`/`uvx`) | Opsi transpile & upstream gagal | Deteksi runtime + pesan error actionable sebelum pipeline berjalan |
| Sandbox tanpa container bisa akses network saat `go build` | Permukaan supply-chain lebih lebar | Pinning `go.sum` + module cache terkontrol di M2; container opsional di M3 |
| Key management Cosign | Menunda release bertanda tangan | SHA-256 murni cukup untuk M1/M2; Cosign terisolasi di M3 |
| Registry repo belum ada saat M1 selesai | Client tidak punya sumber data | M0 mendahului M1; `registry.json` seed dengan 3 port sebelum client di-merge |

---

## 🗺️ 7. Urutan Eksekusi

```
M0 (Foundation)  ──>  M1 (MVP Install)  ──>  M2 (Transpiler)  ──>  M3 (Governance)
   schema, repo,        search + install,      probe, generate,        livecheck, cosign,
   3 port manual        SHA-256, UX            verify, share-PR        flavors, hardening
```

Prasyarat lintas milestone: repo `scorp-mcp-registry` publik sebelum M1 merge; `mark3labs/mcp-go` ditambahkan sebagai dependency **hanya untuk port referensi & hasil transpiler** (client scorp tetap hand-rolled).
