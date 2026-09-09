# 🥊 Architecture Comparison: Scorp vs PicoClaw vs ZeroClaw

This document provides an in-depth head-to-head architectural analysis between three ultra-lightweight autonomous AI agents based on native compiled binaries: **Scorp** (Go), **PicoClaw** (Sipeed / Go), and **ZeroClaw** (ZeroClaw Labs / Rust).

---

## 📊 1. Comparison Matrix Summary

| Dimension / Feature | 🦂 **Scorp** | 🦞 **PicoClaw** (Sipeed) | 🦀 **ZeroClaw** (ZeroClaw Labs) |
| :--- | :--- | :--- | :--- |
| **Language** | **Go 1.25+** (Static binary) | **Go** (Static binary, ~95% AI-bootstrap) | **Rust** (Memory safety, no GC) |
| **Average RAM Footprint** | **~10 – 25 MB** | **< 10 MB** (Optimized for $10 SBCs) | **< 5 MB** (Lowest thanks to native Rust) |
| **Primary Focus & DNA** | **DevOps, Sysadmin, Linux Server Ops & AI Coding** | **Embedded Linux, IoT, SBC & Home Assistant** | **Local-First Agent, Sandboxed Runtime & Security** |
| **Target Hardware** | VPS (x86/ARM64), Cloud VMs, Linux PCs, Android Termux | $10 SBCs (LicheeRV Nano, NanoKVM, RPi Zero, RISC-V) | Linux Servers, Desktops, Hardened Containers |
| **CPU Architectures** | x86_64, ARM64 | **x86_64, ARM64, ARMv6 (32-bit), RISC-V, LoongArch** | x86_64, ARM64 |
| **Primary Interfaces** | **Interactive Telegram Daemon** (Inline Keyboards, Health Checks, Session UI) + Interactive CLI + Web Gateway | **Multi-Channel Gateway** (Telegram, Discord, DingTalk, Feishu, WeCom, QQ) + Flutter FUI | **TUI (`zerocode`)**, CLI + **HTTP Gateway Daemon** (`port 42617`) |
| **LLM Providers** | **Hybrid:** Command Code, DeepSeek, Gemini, OpenAI, Claude, Ollama, OpenRouter, Groq, custom OpenAI APIs | OpenAI, Anthropic, DeepSeek, OpenRouter, Groq, Zhipu | 20+ Providers (OpenAI, Anthropic, Ollama, Groq, DeepSeek, etc.) |
| **Prompt Caching** | ✅ **Native Telemetry & 90% Discount Tracking** (Cloudflare / DeepSeek / Vercel AI SDK) | ⚠️ Dependent on upstream API standard | ⚠️ Dependent on crate provider |
| **System Collectors** | **Native Go Collectors** (CPU, RAM, Disk, Systemd, Process, Thermal, Network) | Standard Linux CLI execution | Standard Linux CLI execution |
| **Hardware Tools** | Remote VPS & CLI | **Physical Sensors (I2C/SPI) & GPIO Pins** | Standard local tools |
| **Security & Guardrails**| 4 Autonomy Modes (`readonly`, `supervised`, `auto`, `yolo`), Bubblewrap Sandbox, Deny Rules, Cryptographic Receipts | Standard permissions, embedded sandbox | **Defense-in-depth**, Workspace Jail, Path Traversal Blocker, Private IP Blocker, Distroless Docker |
| **RAG & Knowledge** | **Local Simhash & Vector RAG** (0 external DB) | Basic file context / search | Pluggable memory engines |

---

## 🌐 2. Web Search Capabilities & Third-Party Strategy Analysis

In the realm of web search, the three projects take noticeably different strategic approaches:

### 🦀 ZeroClaw (Third-Party & Commercial Heavy)
ZeroClaw provides flexible search engine integrations, but relies heavily on paid commercial APIs:
* **Tavily Multi-Key Round-Robin:** Uses paid Tavily API keys with key rotation and load-balancing to mitigate rate limits.
* **Brave Search API:** Relies on commercial Brave API keys.
* **Jina Reader & Exa:** Leverages external reader services for LLM context transformation.
* **SearXNG (Self-Hosted):** A strong non-third-party option connecting to a personal metasearch server to avoid scraper bans.

### 🦞 PicoClaw (Multi-Region & Auto-Fallback Focus)
* Integrates commercial **Brave Search** API as its primary backend.
* Routes Asian/Chinese queries to **Baidu & Sogou APIs**.
* Features `provider: "auto"`, falling back to DuckDuckGo scraping when third-party quotas are exhausted.

### 🦂 Scorp (Native Embedded Metasearch Champion)
Scorp focuses on clean, self-contained local engineering:
* **Embedded Native MetaSearch Cluster (Zero-RAM, Zero-Docker):**
  Scorp includes a built-in concurrent metasearch engine querying **Bing**, **DuckDuckGo**, **Wikipedia OpenSearch**, and **GitHub Repositories** in parallel via Go goroutines.
  - *Consensus Ranking:* URLs discovered across multiple engines receive higher relevance weight.
  - *Smart Deduplication:* Groups identical URLs and merges the most comprehensive snippet.
  - *Anti-Single-Point-of-Failure:* If one engine rate-limits (e.g. DDG 403), the other engines (Bing/Wiki/GitHub) continue returning results seamlessly.
  - *Pluggable Hooks:* If the operator configures a **personal SearXNG instance** (`SEARXNG_URL`), **Brave API** (`BRAVE_API_KEY`), or **Tavily API** (`TAVILY_API_KEY`), Scorp auto-connects them immediately.
* **Clean Page Extraction Engine (`read_url`):**
  Rather than raw page dumps, Scorp employs a **Tiered Web Engine (< 5MB RAM)**:
  1. *Local Zero-RAM:* Streams raw HTML via native HTTP streams.
  2. *Mozilla Readability Parser:* Strips boilerplate, ads, and navigation using `go-shiori/go-readability`.
  3. *AST Markdown Converter:* Converts clean DOM into Markdown using `JohannesKaufmann/html-to-markdown`.
  4. *Cloud Scraper Fallback:* Only if a page is Cloudflare-protected or JavaScript-heavy does Scorp offload to remote headless scrapers (Firecrawl / Tavily API).

---

## 🎯 3. Selection Recommendations

1. **Choose SCORP if:**
   * Your primary workflows are **Linux server operations, DevOps, VPS automation**, and autonomous remote coding from mobile via **Telegram** or terminal CLI.
   * You demand genuine cost efficiency via prompt caching telemetry, native Go single-binary execution, and low memory usage (<25MB RAM).

2. **Choose PICOCLAW if:**
   * You need to deploy an AI agent onto **$10 single-board computers** (LicheeRV Nano, NanoKVM, Raspberry Pi Zero) or **RISC-V** hardware for IoT and GPIO sensor control.
   * You require integration with Asian messaging platforms (DingTalk, Feishu, WeCom, LINE).

3. **Choose ZEROCLAW if:**
   * You need **heavy container-grade sandboxing and process isolation** for corporate or multi-tenant hosting.
   * You prefer a pure **Rust** stack with sub-5MB RAM targets and paid API integrations (Tavily/Brave).

---

## 🚀 4. Completed Scorp Feature Enhancements

Scorp's web search capabilities have been fully modernized without adding memory bloat:
- [x] **Embedded Native Multi-Engine Metasearch** (`tools/metasearch.go`) running in parallel (Bing + DuckDuckGo + Wikipedia + GitHub).
- [x] **Consensus Scoring & Deduplication** (URLs found across multiple sources rank first).
- [x] **Regional & Anti-Pollution Filters** (Enforcing clean locale headers on search queries).
- [x] **Integrated Pluggable Hooks:**
  - `SEARXNG_URL` for local/remote SearXNG instances.
  - `BRAVE_API_KEY` for commercial Brave Search integration.
  - `TAVILY_API_KEY` for commercial Tavily integration.

---

## 📈 5. Empirical Benchmarks

For empirical head-to-head testing on live VPS infrastructure, see:
* [🥊 Benchmark Komparasi 20 Tugas Ringan: Scorp vs PicoClaw vs ZeroClaw](./BENCHMARK_LIGHT_20_SCORP_PICOCLAW_ZEROCLAW.md) — Live benchmark covering factual QA, arithmetic, logic riddles, code generation, formatting, and DevOps tasks on Google Gemini 3.5 Flash-Lite.
* [🏗️ Benchmark Komparasi 15 Tugas Agentic Berat: Scorp vs PicoClaw vs ZeroClaw](./BENCHMARK_HEAVY_AGENTIC_15_SCORP_PICOCLAW_ZEROCLAW.md) — Live autonomous benchmark testing tool calling, multi-step directory structures, code execution, bug diagnosis & self-repair, unit testing (TDD), CSV-to-JSON transformations, cryptographic hashing, and system telemetry on real VPS.
* [🏛️ Benchmark Komparasi 15 Tugas Agentic Brutal & Kompleks: Scorp vs PicoClaw vs ZeroClaw](./BENCHMARK_BRUTAL_15_SCORP_PICOCLAW_ZEROCLAW.md) — Deep engineering benchmark testing multi-file Go compilation, SQLite transactions, background REST API processes, signal traps, HMAC authentication, Git multi-branch trees, and system cgroup audits.
* [🌌 Benchmark Komparasi 15 Tugas Rekayasa Ekstrem & Arsitektur Sistem: Scorp vs PicoClaw vs ZeroClaw](./BENCHMARK_EXTREME_15_SCORP_PICOCLAW_ZEROCLAW.md) — Extreme engineering benchmark testing native C shared libraries compilation via gcc + ctypes, async SQLite multi-worker task queues, self-signed TLS certificates & HTTPS server daemon, recursive JSON schema diff engines, Linux memory limits & OOM handling, non-blocking TCP socket suites, and live Prometheus metric exporters.




