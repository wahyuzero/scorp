# 🥊 Komparasi Arsitektur: Scorp vs PicoClaw vs ZeroClaw

Dokumen ini memuat analisis perbandingan mendalam (*head-to-head*) antara tiga agen AI otonom ultra-ringan berbasis native compiled binary: **Scorp** (Go), **PicoClaw** (Sipeed / Go), dan **ZeroClaw** (ZeroClaw Labs / Rust).

---

## 📊 1. Ringkasan Matriks Perbandingan

| Dimensi / Fitur | 🦂 **Scorp** | 🦞 **PicoClaw** (Sipeed) | 🦀 **ZeroClaw** (ZeroClaw Labs) |
| :--- | :--- | :--- | :--- |
| **Bahasa Pemrograman** | **Go 1.26+** (Static binary) | **Go** (Static binary, ~95% AI-bootstrap) | **Rust** (Memory safety, no GC) |
| **Alokasi RAM Rata-rata** | **~10 – 15 MB** | **< 10 MB** (Dioptimalkan untuk SBC $10) | **< 5 MB** (Terendah berkat Rust native) |
| **Fokus & DNA Utama** | **DevOps, Sysadmin, Linux Server Ops & AI Coding** | **Embedded Linux, IoT, SBC & Home Assistant** | **Local-First Agent, Sandboxed Runtime & Security** |
| **Target Hardware** | VPS (x86/ARM64), Cloud VM, Linux PC, Termux Android | SBC $10 (LicheeRV Nano, NanoKVM, RPi Zero, RISC-V) | Linux Server, Desktop, Hardened Containers |
| **Arsitektur CPU** | x86_64, ARM64 | **x86_64, ARM64, ARMv6 (32-bit), RISC-V, LoongArch** | x86_64, ARM64 |
| **Antarmuka Utama** | **Telegram Daemon Interaktif** (Inline Role Picker, Health Check, Session UI) + Interactive CLI | **Multi-Channel Gateway** (Telegram, Discord, DingTalk, Feishu, WeCom, QQ) + Flutter FUI | **TUI (`zerocode`)**, CLI + **HTTP Gateway Daemon** (`port 42617`) |
| **Provider LLM** | **Dual Hybrid:** Command Code (Vercel AI Gateway) + OpenCode Zen (Free) + Direct OpenAI/Anthropic/Gemini | OpenAI, Anthropic, DeepSeek, OpenRouter, Groq, Zhipu | 20+ Provider (OpenAI, Anthropic, Ollama, Groq, DeepSeek, dll.) |
| **Prompt Caching** | ✅ **Native Telemetry & 90% Discount Tracking** (Cloudflare / DeepSeek / Vercel AI SDK) | ⚠️ Tergantung upstream API standard | ⚠️ Tergantung crate provider |
| **System Collectors** | **Native Go Collectors** (CPU, RAM, Disk, Systemd, Process, Thermal, Network) | Standard Linux CLI execution | Standard Linux CLI execution |
| **Hardware Tools** | Remote VPS & CLI | **Sensor Fisik (I2C/SPI) & Pin GPIO** | Standard local tools |
| **Keamanan & Guardrails** | Autonomy Modes (`readonly`, `supervised`, `yolo`), Cryptographic Receipts, Command Blacklist | Standard permissions, embedded sandbox | **Defense-in-depth**, Workspace Jail, Path Traversal Blocker, Private IP Blocker, Distroless Docker |
| **RAG & Knowledge** | **Local Simhash & Vector RAG** (0 external DB) | Basic file context / search | Pluggable memory engines |

---

## 🌐 2. Analisis Kemampuan Web Search & "Trik Pihak Ketiga"

Dalam ranah pencarian web (*web search*), terdapat perbedaan strategi yang sangat menarik antara ketiga proyek:

### 🦀 ZeroClaw (Raja Pihak Ketiga & Metasearch)
ZeroClaw menduduki peringkat teratas dalam hal fleksibilitas search engine, namun kekuatannya didorong oleh ketergantungan pada layanan API komersial:
* **Tavily Multi-Key Round-Robin:** Menggunakan API berbayar Tavily dengan mekanisme pemutaran kunci (*load-balancing*) agar tidak terkena limit.
* **Brave Search API:** Menggunakan API komersial Brave.
* **Jina Reader & Exa:** Memanfaatkan web reader eksternal untuk konversi LLM context.
* **SearXNG (Self-Hosted):** Satu-satunya keunggulan non-pihak-ketiga di mana ZeroClaw dapat dihubungkan ke server metasearch pribadi untuk menghindari blokir scraper.

### 🦞 PicoClaw (Fokus Multi-Region & Auto-Fallback)
* Mengintegrasikan API komersial **Brave Search** sebagai backend utama.
* Memiliki rute pencarian wilayah Asia/China ke **Baidu & Sogou API**.
* Menyediakan fitur `provider: "auto"` yang otomatis turun (*fallback*) ke scraping DuckDuckGo jika kuota API pihak ketiga habis.

### 🦂 Scorp (The Native Embedded Metasearch Winner)
Scorp menerapkan arsitektur *clean engineering* yang sangat mandiri di mesin lokal:
* **Embedded Native MetaSearch Cluster (Zero-RAM, Zero-Docker):**
  Scorp kini memiliki metasearch engine bawaan yang berjalan paralel via Go goroutines menembak **Bing**, **DuckDuckGo**, **Wikipedia OpenSearch**, dan **GitHub Repositories**.
  - *Consensus Ranking:* URL yang muncul di beberapa engine sekaligus otomatis mendapatkan skor relevansi lebih tinggi.
  - *Smart Deduplication:* Mengelompokkan URL identik dan menggabungkan snippet terpanjang.
  - *Anti-Single-Point-of-Failure:* Jika satu engine terkena limit/blokir (misal DDG 403), engine lain (Bing/Wiki/GitHub) tetap mengembalikan hasil tanpa membuat agen gagal.
  - *Pluggable Hooks:* Jika user ingin memakai server **SearXNG pribadi** (`SEARXNG_URL`), **Brave API** (`BRAVE_API_KEY`), atau **Tavily API** (`TAVILY_API_KEY`), Scorp langsung menghubungkannya secara otomatis!
* **Ekstraksi Konten Halaman Terbersih (`read_url`):**
  Alih-alih langsung menyedot data mentah, Scorp membangun arsitektur **Tiered Web Engine (< 5MB RAM)**:
  1. *Local Zero-RAM:* Mengambil HTML mentah via HTTP stream native.
  2. *Mozilla Readability Parser:* Membersihkan boilerplate, iklan, dan navigasi menggunakan engine `go-shiori/go-readability`.
  3. *AST Markdown Converter:* Mengubah konten bersih menjadi Markdown menggunakan `JohannesKaufmann/html-to-markdown`.
  4. *Cloud Scraper Fallback:* Hanya jika halaman memblokir bot (Cloudflare) atau mewajibkan eksekusi JavaScript berat, barulah Scorp mengalihkan ke remote headless scraper (Firecrawl / Tavily API).

---

## 🎯 3. Posisi & Rekomendasi Pilihan

1. **Gunakan SCORP jika:**
   * Kebutuhan utama adalah **mengelola server Linux, DevOps, automation VPS**, dan coding mandiri jarak jauh langsung dari HP via **Telegram Bot interaktif** atau terminal CLI.
   * Menginginkan efisiensi biaya nyata lewat integrasi **Command Code + OpenCode Zen** dengan prompt caching telemetri riil.

2. **Gunakan PICOCLAW jika:**
   * Anda ingin memasang AI agent di **hardware mini $10** (LicheeRV Nano, NanoKVM, Raspberry Pi Zero) atau arsitektur **RISC-V** untuk membaca sensor IoT.
   * Memerlukan integrasi ke chat messenger Asia (DingTalk, Feishu, WeCom, LINE).

3. **Gunakan ZEROCLAW jika:**
   * Memerlukan **isolasi keamanan tingkat tinggi (sandboxing/jailing)** untuk lingkungan korporat/multi-tenant.
   * Menginginkan ekosistem **Rust** murni dengan konsumsi RAM sub-5MB dan integrasi API search berbayar (Tavily/Brave).

---

## 🚀 4. Status Implementasi Fitur Baru Scorp (Completed)

Sektor pencarian web Scorp telah resmi ditingkatkan tanpa menambah footprint memori:
- [x] **Embedded Native Multi-Engine Metasearch** (`tools/metasearch.go`) berjalan paralel (Bing + DuckDuckGo + Wikipedia + GitHub).
- [x] **Consensus Scoring & Deduplication** (URL yang muncul di banyak engine mendapat ranking teratas).
- [x] **Anti-Pollution Regional Filters** (Enforce English locale headers & parameters pada search engine).
- [x] **Pluggable Hooks Terintegrasi:**
  - `SEARXNG_URL` untuk menghubungkan instance SearXNG lokal/eksternal.
  - `BRAVE_API_KEY` untuk opsi komersial Brave Search.
  - `TAVILY_API_KEY` untuk opsi komersial Tavily.
