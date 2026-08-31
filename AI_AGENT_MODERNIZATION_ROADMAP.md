# AI Coding Agent Modernization Roadmap & Blueprint (v2.0)

> **State-of-the-Art (SOTA) Research & Comprehensive Engineering Upgrade Plan**  
> **Target Repositories:** [`wahyuzero/scorp`](https://github.com/wahyuzero/scorp) (Go) & [`wahyuzero/termagent`](https://github.com/wahyuzero/termagent) (TypeScript/Termux)  
> **Focus:** Memutakhirkan engine agent dari era awal (GPT-4/GLM) ke era modern (Multi-Model, Prompt Caching, MCP Standard, Surgical Diffs, Native Termux).

---

## 📌 Daftar Isi
1. [Lanskap Model AI Modern & Verifikasi Data (2026)](#1-lanskap-model-ai-modern--verifikasi-data-2026)
2. [Evaluasi Arsitektur Scorp & Termagent Saat Ini](#2-evaluasi-arsitektur-scorp--termagent-saat-ini)
3. [Pilar Upgrade 1: Multi-Model & Fallback Engine](#3-pilar-upgrade-1-multi-model--fallback-engine)
4. [Pilar Upgrade 2: Prompt Caching & Token Economics](#4-pilar-upgrade-2-prompt-caching--token-economics)
5. [Pilar Upgrade 3: Surgical Diff Editing (Bukan Full Overwrite)](#5-pilar-upgrade-3-surgical-diff-editing-bukan-full-overwrite)
6. [Pilar Upgrade 4: Standardisasi MCP (Model Context Protocol)](#6-pilar-upgrade-4-standardisasi-mcp-model-context-protocol)
7. [Pilar Upgrade 5: Subagent & Task Delegation System](#7-pilar-upgrade-5-subagent--task-delegation-system)
8. [Pilar Upgrade 6: Deep Termux & Mobile Performance Optimization](#8-pilar-upgrade-6-deep-termux--mobile-performance-optimization)
9. [Roadmap Eksekusi Teknis (Fase 1 s.d. 4)](#9-roadmap-eksekusi-teknis-fase-1-sd-4)

---

## 1. Lanskap Model AI Modern & Verifikasi Data (2026)

### 📊 Tabel Perbandingan Model AI Coding Saat Ini:

| Model | Provider | Jendela Konteks | Harga Input / Output (per 1M token) | Cache Hit Discount | Karakteristik Utama |
|---|---|---|---|---|---|
| **Gemini 1.5 / 2.0 Flash** | Google | **1,000,000+ token** | **~$0.075 / $0.30** | Diskon hingga **~90%** ($0.0075) | Respon sub-detik (< 1s), telan seluruh repo, gratis kuota harian. |
| **DeepSeek V3** | DeepSeek / OpenRouter | **128,000 token** | **~$0.14 / $0.28** | Diskon **~50%-80%** ($0.014) | SOTA reasoning, sangat murah, akurat untuk refactoring besar. |
| **DeepSeek R1** | DeepSeek / OpenRouter | **128,000 token** | **~$0.55 / $2.19** | - | Model penalaran CoT (*Chain of Thought*) untuk debugging rumit. |
| **Qwen 2.5 Coder (32B/72B)** | Alibaba / OpenRouter | **128,000 token** | **~$0.20 / $0.60** | - | Khusus sintaks kode, bisa di-host lokal via Ollama. |
| **Claude 3.5 Sonnet** | Anthropic | **200,000 token** | **$3.00 / $15.00** | Diskon **90%** ($0.30) | Standar industri tertinggi untuk coding complex, tapi lebih mahal. |

---

## 2. Evaluasi Arsitektur Scorp & Termagent Saat Ini

### 🦂 Status `scorp` (Golang Engine):
* **Kelebihan:** Ditulis dengan Go (performa native, konsumsi memori rendah, single binary, startup kilat, mudah di-compile untuk ARM64 Android).
* **Area Upgrade:** Loop model perlu dukungan dynamic streaming SSE untuk Gemini & DeepSeek; parser tool calls perlu standarisasi JSON Schema; peremajaan MCP client.

### 📱 Status `termagent` (TypeScript Engine):
* **Kelebihan:** Integrasi brilian dengan `Termux:API` (notifikasi Android, dialog konfirmasi modal, clipboard sync, haptic vibration).
* **Area Upgrade:** Migrasi runtime ke **Bun / Node 22+ Native ARM64** untuk memangkas *memory footprint* `node_modules`; implementasi chunk replacement editing.

---

## 3. Pilar Upgrade 1: Multi-Model & Fallback Engine

Menerapkan sistem routing model cerdas bertingkat (*Tiered Model Routing*):

```
                       [ User Input / Task ]
                                 ⬇️
            [ Complexity Router / Task Classifier ]
             /                   |                   \
   [ Simple Query / Read ]   [ Code Edit / Gen ]   [ Deep Bug / Architecture ]
             ⬇️                  ⬇️                     ⬇️
      Gemini Flash          DeepSeek V3 /          DeepSeek R1 /
      ($0.075/1M token)     Qwen 2.5 Coder         Claude Sonnet
```

### 🛡️ Automatic Fallback Chain:
Jika API utama mengalami rate limit (HTTP 429) atau downtime (HTTP 503), agent otomatis beralih ke provider cadangan tanpa memutus sesi:
`Gemini Flash ➔ DeepSeek V3 ➔ Qwen Coder ➔ OpenRouter / 9Router`.

---

## 4. Pilar Upgrade 2: Prompt Caching & Token Economics

Pada agent coding biasa, setiap langkah (*step*) mengirim ulang seluruh riwayat chat, instruksi sistem, dan file-file yang dibaca. Ini membuat biaya token membengkak secara eksponensial.

### ⚡ Strategi Caching:
1. **Static System Prefix:** Instruksi sistem, format tools, dan aturan coding diletakkan di bagian paling awal prompt (posisi tetap).
2. **Repository Map Prefix:** Ringkasan pohon file proyek diletakkan setelah system prompt.
3. **Hasil:** Setiap turn percakapan hanya membayar token baru (*cache read hit*), menghemat **75%–90% total biaya API**.

---

## 5. Pilar Upgrade 3: Surgical Diff Editing (Bukan Full Overwrite)

### 🛑 Masalah Cara Lama (Full File Rewrite):
* File `index.html` (500 baris) diedit 2 baris $\rightarrow$ AI menulis ulang 500 baris.
* Memakan 2.000+ output token, lambat (15-30 detik), dan rawan memotong kode yang tidak sengaja terhapus.

### ✂️ Solusi Modern (Chunk Replacement / Exact Match):
Tool `replace_file_content` hanya menerima 4 parameter:
```typescript
interface ReplaceFileContentArgs {
  target_file: string;
  start_line: number;
  end_line: number;
  target_content: string;      // Teks eksak yang ingin diganti (5-10 baris)
  replacement_content: string; // Teks pengganti baru
}
```
* **Hasil:** AI hanya mengeluarkan 50 token. Waktu eksekusi turun dari 20 detik menjadi **0.8 detik** di HP!

---

## 6. Pilar Upgrade 4: Standardisasi MCP (Model Context Protocol)

Mengadopsi protokol terbuka **Anthropic MCP (2025/2026 Standard)**:

```
[ Scorp / TermAgent Core ]
        │
        ├── Stdio / SSE Transport (JSON-RPC 2.0)
        │
        ├── 📂 MCP Filesystem Server
        ├── 🌐 MCP CloakBrowser / Headless Browser
        ├── 🗄️ MCP SQLite / Database
        ├── 🐙 MCP GitHub Server
        └── 📱 MCP Termux:API Bridge
```

Dengan MCP, `scorp` dan `termagent` tidak perlu menulis ulang integrasi tool baru. Cukup pasang konfigurasi server di `mcp.json`.

---

## 7. Pilar Upgrade 5: Subagent & Task Delegation System

Untuk proyek besar, chat history utama cepat penuh (*context bloat*). Solusinya adalah delegasi subagent:

1. **Main Planner Agent:** Mengatur strategi, membagi task, dan berinteraksi dengan user.
2. **Research Subagent (Read-Only):** Membaca puluhan file dan mencari grep di background, lalu hanya mengembalikan ringkasan 3 baris ke agent utama.
3. **Execution Subagent (Worker):** Menjalankan build/test di workspace terisolasi.

---

## 8. Pilar Upgrade 6: Deep Termux & Mobile Performance Optimization

1. **Kompilasi Binary Native Android (Go):**
   ```bash
   CGO_ENABLED=0 GOOS=android GOARCH=arm64 go build -ldflags="-s -w" -o scorp-termux main.go
   ```
   * Menghasilkan binary tunggal ~12MB tanpa dependensi libc eksternal, jalan langsung di `/data/data/com.termux/files/usr/bin/scorp`.
2. **Optimasi I/O untuk Penyimpanan eMMC / UFS HP:**
   * Menggunakan `fd` dan `ripgrep` native binary Termux yang diakses lewat child process streams.
3. **Background Notification & WakeLock:**
   * Memanfaatkan `termux-wake-lock` saat agent mengerjakan task multi-step agar Android tidak mematikan proses saat layar HP mati.
   * `termux-notification` untuk memberitahu user saat task selesai.

---

## 9. Roadmap Eksekusi Teknis

```mermaid
gantt
    title Roadmap Modernisasi Scorp & Termagent v2.0
    dateFormat  YYYY-MM-DD
    section Fase 1 (Core)
    Multi-Model Client (Gemini 2.0 / DeepSeek) :a1, 2026-09-01, 7d
    Surgical Diff Tool (Replace Chunk)          :a2, after a1, 5d
    section Fase 2 (Protocol)
    MCP Client Stdio/SSE Implementation         :b1, after a2, 7d
    Prompt Caching Optimization Engine          :b2, after b1, 4d
    section Fase 3 (Mobile)
    Native Termux:API Bridge Tools              :c1, after b2, 5d
    WakeLock & Background Task Notifier         :c2, after c1, 3d
    section Fase 4 (Advanced)
    Subagent Delegation Protocol                :d1, after c2, 7d
    Release Scorp v2.0 & Termagent v2.0         :d2, after d1, 3d
```

---

*Cetak biru ini disusun berdasarkan riset model AI dan arsitektur agent termutakhir. Siap diimplementasikan bertahap.*
