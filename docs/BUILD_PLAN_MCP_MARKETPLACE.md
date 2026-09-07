# 🦂 Scorp MCP Marketplace — Build Plan & Execution Roadmap

> **Status:** Executed — M0/M1/M2 completed and validated; M3 partial (container hardening deferred)  
> **Blueprint Reference:** [`docs/MCP_MARKETPLACE_BLUEPRINT.md`](MCP_MARKETPLACE_BLUEPRINT.md)  
> **Target Version:** Scorp Agent v2.5 / v3.0  
> **Owner:** Wahyu  
> **Created:** September 2026  

---

## 🎯 1. Executive Summary

This document translates the MCP Marketplace blueprint into an executable build plan, grounded in an actual audit of the scorp-agent codebase.

The central finding of the audit: **this blueprint leverages existing infrastructure rather than building from scratch.** All three installation options (prebuilt / rebuild / upstream) converge at the exact same point — an entry in `~/.scorp/mcp.json`. Once that entry exists, the entire existing infrastructure activates automatically: watchdog, hot-reload, health monitoring, and native tool registration in `registry/`. Therefore, Tri-Option is a **new sourcing layer in front of the existing `mcp_manage`**, not an entirely new subsystem.

**Agreed Execution Strategy:**

1. **Vertical Slice First** — Implement `/mcp search` + `/mcp install` with Tri-Option for seed ports (prebuilt + upstream) prior to the transpiler. The marketplace feels functional from M1 with minimal operational risk.
2. **Handwritten Reference Ports** — `mcp-fetch`, `mcp-sqlite`, and `mcp-filesystem` written idiomatically as the golden standard, serving as the first marketplace catalog entries and transpiler validation test cases.
3. **Lightweight Sandbox Build** — `go build` in an isolated temp directory with checksum verification; container hardening follows later.
4. **Self-Hosted Codegen** — The transpiler utilizes Scorp's internal LLM gateway (routing through *complex/premium* models in `models.json`) without external CLI dependencies.

---

## 🔍 2. Codebase Audit: Existing vs. Gaps

### 2.1 Existing Components (No Need to Rebuild)

| Blueprint Component | Existing Implementation | Location |
| :--- | :--- | :--- |
| MCP client (JSON-RPC 2.0, stdio + SSE) | Custom client, server spawning, hot-reload | `mcp/client.go`, `mcp/sse.go` |
| Server lifecycle & crash recovery | Watchdog + health status | `mcp/watchdog.go`, `mcp/RestartServer` |
| Tool registration to LLM | Native function calling + deferred tools + raw JSON Schema | `registry/registry.go` |
| Server management (add/remove/list/reload) | Agent tool `mcp_manage` with hot-reload | `mcp/manage.go` |
| **Layer 5: Secret Redactor** ✅ | Regex tool output sanitizer | `tools/redact.go` (+ `redact_test.go`) |
| Transpiler codegen engine | Multi-provider gateway, routing rules, fallback on error | `gateway/gateway.go`, `models.json` |
| UX Interfaces | `/mcp` CLI command + Telegram MCP inline buttons | `cli.go`, `telegram/telegram.go` |
| Centralized config paths | `mcpConfigFilePath()` → `~/.scorp/mcp.json` | `config/config_paths.go` |

### 2.2 Gaps to Build

- ❌ Registry client: fetch + cache `registry.json`, local search, manifest parser
- ❌ JSON Schema `scorp-mcp.json` v1 (health matrix, flavors, contributors, origin)
- ❌ Tri-Option install flow (interactive CLI + Telegram inline keyboard)
- ❌ SHA-256 pinning & artifact verification (Cosign deferred)
- ❌ Transpiler pipeline: probe → codegen → sandbox build → contract verify
- ❌ Drift livecheck & upstream sync indicator
- ❌ PR generator (port contributions / maintenance updates)
- ❌ Public repo `scorp-mcp-registry` + 3 CI workflows
- ❌ 3 Go reference ports (fetch, sqlite, filesystem)

---

## 🏛️ 3. Key Architectural Insights

### 3.1 Tri-Option Convergence Point

```
Option 1 Prebuilt  ──┐  download binary → SHA-256 verify → ~/.scorp/mcp-binaries/
                     ├──>  Entry in ~/.scorp/mcp.json  ──>  Existing infra takes over
Option 2 Rebuild   ──┤  (transpile → go build → local bin)       (watchdog, hot-reload,
                     │                                            native tool registry)
Option 3 Upstream  ──┘  register direct npx/uvx (already supported)
```

Implication: Existing `mcp_manage(action="add")` is the ultimate endpoint for ALL install paths. New components merely decide *how to generate* that configuration entry.

### 3.2 Transpiler = Self-Hosted Codegen

- **Phase 1 (Probe):** Reuse JSON-RPC client `mcp/client.go` to spawn an ephemeral server (detecting `npx`/`uvx`/`python3` runtime) and capture `initialize` + `tools/list` as benchmark contracts.
- **Phase 2 (Generate):** Structured codegen prompt for `mark3labs/mcp-go`, executed via `gateway/` with *complex* (premium) model routing.
- **Phase 3 (Verify):** `go build` in a temp directory, boot the compiled binary, diff tool schema against Phase 1 benchmark. On failure → graceful degradation to upstream runtime.

### 3.3 Security Layers: Initial Status

| Layer | Status | Notes |
| :--- | :--- | :--- |
| 1. Source-Only Registry | ⬜ New (M0) | Registry repo policy + CI |
| 2. AST & Prompt Scan | ⬜ New (M0) | `gosec` + `govulncheck` + prompt linter in CI |
| 3. Hermetic CI/CD | ⬜ New (M0) | GitHub Actions, multi-arch |
| 4. SHA-256 Pinning | ⬜ New (M1) | Client-side pre-execution verification |
| 5. Secret Redactor | ✅ **Existing** | Audit coverage on MCP client output paths |

---

## ⚠️ 4. Decisions & Blueprint Alignment

1. **Port Source Location:** Section 5C placed `main.go` directly in the registry repo, whereas section 5B pointed `port_repository` to a separate `scorp-mcp-ports` repo.
   **Decision:** Port sources live directly in `scorp-mcp-registry/servers/<name>/main.go` for the MVP (single source of truth, single CI audit point). The `port_repository` field remains optional.
2. **Cosign:** Deferred to M3. MVP relies on pure SHA-256; Cosign key management is a dedicated sub-project.
3. **Configuration Cleanup:** Migrate hardcoded `$HOME/.scorp/mcp.json` references to `config.ConfigFilePathMCP()`.

---

## 🚀 5. Milestones & Deliverables

### M0 — Foundation: Schema, Registry Repo & Reference Ports

*Scope: Outside the scorp-agent core repository.*

- [x] Finalize JSON Schema `scorp-mcp.json` v1: fields `origin` (port/native), `upstream.pinned_commit`, `health` (status, coverage_score, active/unsupported tools), `variant` (flavor), `contributors`, `artifacts` (multi-arch + sha256), `build` (go_version, sdk).
- [x] Publish schema to `scorp-mcp-registry/schema/v1.json`.
- [x] Scaffold public repo `wahyuzero/scorp-mcp-registry` following layout 5C (`registry.json`, `servers/`, `.github/workflows/`).
- [x] CI Workflows: `security-audit.yml` (gosec, govulncheck, prompt-injection linter, JSON-RPC contract validation), `release-binaries.yml` (cross-compile linux amd64/arm64 + checksums), `drift-livecheck.yml` (stub, enabled in M3).
- [x] Implement 3 idiomatic Go reference ports (`mark3labs/mcp-go`, Go 1.25+):
  - [x] `mcp-fetch` — web scraping & markdown extraction
  - [x] `mcp-sqlite` — local database inspection & querying
  - [x] `mcp-filesystem` — scoped local filesystem access
- [x] `scorp-mcp.json` manifests for all three ports (artifacts populated by release CI on first tag).

**Acceptance:** `registry.json` indexed; all three binaries pass security-audit CI and run manual JSON-RPC handshakes successfully.

### M1 — Vertical Slice: Marketplace Client & Tri-Option Install (MVP)

*Scope: New package `mcp/marketplace/` + CLI & Telegram integration.*

- [x] `mcp/marketplace/registry.go` — fetch `registry.json` from GitHub raw, local cache with ETag/TTL, fallback to stale cache when offline.
- [x] `mcp/marketplace/manifest.go` — parser + manifest validation against schema v1.
- [x] `mcp/marketplace/search.go` — local search (name, description, tool, flavor).
- [x] `mcp/marketplace/install.go` — Tri-Option orchestration:
  - [x] Option 1 (Prebuilt): download per-architecture artifact, **verify SHA-256 before execution**, store in `~/.scorp/mcp-binaries/`, register via `mcp_manage add`.
  - [x] Option 2 (Rebuild): displayed with status "coming in v2.5 — run transpiler manually" (fully enabled in M2).
  - [x] Option 3 (Upstream): register `npx`/`uvx` directly — reuse existing pathway.
- [x] Pre-install disclosure: health badge (🟢 full / 🟡 partial), disabled tools list + rationale, resource delta vs upstream.
- [x] CLI: subcommands `/mcp search <term>` and `/mcp install <target>` in `cli_mcp.go` (interactive `[1/2/3]` prompts).
- [x] Telegram: inline Tri-Option keyboard (using existing `callback_data` pattern) and narrative dialogs.
- [x] Layer 5 coverage audit: MCP tool outputs route through `tools/redact.go` in `registerMCPToolsAsNative`.

**Acceptance:** End-to-end `/mcp install fetch` → select option 1 → verified binary registered and tools exposed to the agent; option 3 works for arbitrary Node/Python servers; Telegram dialog matches CLI.

### M2 — AI Transpiler

- [x] `mcp/transpiler/probe.go` — spawn ephemeral upstream server in sandbox dir, detect runtime (`npx`/`uvx`/`python3` with actionable failure messages), capture benchmark: tool schemas, parameters, required fields.
- [x] `mcp/transpiler/generate.go` — structured codegen prompt for `mark3labs/mcp-go`, executed via `gateway/` (*complex* routing), outputting a decoupled standalone `main.go`.
- [x] `mcp/transpiler/build.go` — `go build` in an isolated temp directory; pinning `go.sum`; controlled module cache.
- [x] `mcp/transpiler/verify.go` — boot compiled binary, send mock JSON-RPC requests, diff output schema against Phase 1 benchmark (target: 100% match).
- [x] Graceful degradation: compile/verify failures → transparent blocker report → prompt fallback to upstream runtime.
- [x] Enable Option 2 in Tri-Option flow (CLI + Telegram) via `marketplace.RebuildHook`.
- [x] "Share to Marketplace" workflow: manifest generator + draft PR (with original author attribution) via `gh` (`/mcp share <name>`).
- [x] **Golden standard validation:** transpile upstream fetch/sqlite/filesystem, comparing output against M0 handwritten ports.

**Acceptance:** At least 1 of 3 upstream ports successfully transpiled end-to-end with 100% contract match; failures always provide clear fallback options.

### M3 — Governance & Hardening

- [x] Full `drift-livecheck.yml`: daily checks against upstream tags/commits, update drift fields in `registry.json`, display 🟢/🟡 indicators in search/install + resync notifications.
- [x] Cosign verification step prepared (disabled pending key management) — SHA-256 remains primary Layer 4 defense.
- [ ] Build sandbox hardening: evaluate optional container isolation (Docker if available) for network/filesystem isolation.
- [x] Multi-maintainer collaboration workflow: every update PR undergoes full Layer 2 checks (`security-audit.yml` on pull_request); attribution changelog in manifest.
- [x] Flavor/variant support in search (namespace `@author/name-flavor`) — `variant` indexed in `registry.json`.
- [x] Classification `origin: native` for original Go submissions (bypassing drift check).

---

## 📊 6. Risks & Mitigations

| Risk | Impact | Mitigation |
| :--- | :--- | :--- |
| Variable LLM codegen quality | Transpiler outputs code failing compilation or contract diff | M0 manual ports serve as golden test cases; `verify.go` requires 100% schema match; graceful degradation to upstream |
| Upstream runtime missing on host (`npx`/`uvx`) | Transpile & upstream options fail | Runtime detection + actionable error messages before pipeline execution |
| Uncontained sandbox accessing network during `go build` | Broader supply-chain exposure | Pinned `go.sum` + controlled module cache in M2; optional container isolation in M3 |
| Cosign key management overhead | Delays signed releases | Pure SHA-256 sufficient for M1/M2; Cosign isolated to M3 |
| Registry repo unpopulated at M1 completion | Client lacks data source | M0 precedes M1; seed `registry.json` with 3 ports before client merge |

---

## 🗺️ 7. Execution Order

```
M0 (Foundation)  ──>  M1 (MVP Install)  ──>  M2 (Transpiler)  ──>  M3 (Governance)
   schema, repo,        search + install,      probe, generate,        livecheck, cosign,
   3 manual ports       SHA-256, UX            verify, share-PR        flavors, hardening
```

Cross-milestone prerequisite: public `scorp-mcp-registry` repo before M1 merge; `mark3labs/mcp-go` added as a dependency **only for reference ports & transpiled outputs** (Scorp core client remains custom hand-rolled).
