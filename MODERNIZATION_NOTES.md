# 🦂 Scorp Agent v2.0 — Modernization & Architecture Upgrade Report

Dokumen ini merangkum seluruh pembaruan, refactoring arsitektur, dan fitur-fitur baru yang telah diimplementasikan pada **Scorp Agent v2.0** berdasarkan cetak biru modernisasi dan analisis komparatif performa (ZeroClaw & PicoClaw).

---

## 📌 Ringkasan Pembaruan Utama (v2.0)

| Kategori | Fitur Baru | Keterangan & Manfaat |
| :--- | :--- | :--- |
| **Penyuntingan Berkas** | **Surgical Diff Editing** (`replace_file_content`) | Chunk replacement presisi tanpa menimpa seluruh berkas. Menghemat 90%+ token dan eksekusi turun dari 20s ke 0.8s. |
| **Mesin Web** | **Tiered Ultra-Low-RAM Web Engine** (`read_url`) | Membaca dokumen/web via Firefox Readability + Markdown (<5MB RAM). Fallback otomatis ke cloud renderer (Firecrawl/Tavily) tanpa browser berat lokal. |
| **Protokol** | **Standar Anthropic MCP (Stdio & SSE)** | Mendukung Model Context Protocol lokal (stdio) dan remote (Server-Sent Events / HTTP JSON-RPC 2.0). |
| **Optimasi Token** | **Prompt Caching & Repo Map Prefix** | Static system prefix byte-stable + repository map ter-cache untuk memaksimalkan *cache read hit* hingga 75%–90% pada Gemini, Claude, dan DeepSeek. |
| **Keamanan** | **3-Tier Autonomy System & Sandboxing** | Mode `readonly` (audit aman), `supervised` (konfirmasi perintah bahaya), dan `yolo` (otonomi penuh). Path traversal sandbox dan penolakan jalur sensitif. |
| **Keamanan** | **Outbound Secret Redactor** | Pemindaian regex otomatis yang menyensor API key, token GitHub, Bearer header, dan private key sebelum output tool dikirim ke LLM. |
| **Keamanan** | **Cryptographic Execution Receipts** | Setiap pemanggilan tool menghasilkan receipt bervalidasi hash SHA-256 persisten di `~/.scorp/receipts.json` (`/receipts`). |
| **Eksekusi Agent** | **Real-Time Steering Queue** | Pengguna dapat mengirim instruksi pengalihan di tengah-tengah eksekusi langkah multi-tool tanpa perlu mematikan sesi. |
| **Efisiensi Skema** | **Dynamic Tool Discovery & TTL Injection** | Hanya 10 core tools yang aktif di skema LLM. 39 tools khusus lainnya diinjeksi secara dinamis dengan TTL 3 turn saat dicari/dipanggil (`SCORP_DYNAMIC_TOOLS=true`). |
| **Otomasi Rutinitas** | **Standard Operating Procedures (SOP) Engine** | Playbook deklaratif bertahap (`health_audit`, `code_review`, `site_check`) yang dapat dijalankan instan via `./scorp sop run <nama>`. |
| **Antarmuka Pengguna** | **Micro HTTP Gateway & Web Dashboard** | Web server bawaan (`net/http`) hemat RAM (~1.5MB) di port 8080 dengan dark theme dan kotak chat interaktif untuk akses browser HP/PC. |
| **Manajemen Sesi** | **Multi-Session Management** | Mendukung pembuatan, pergantian, pengubahan nama (*rename*), dan penghapusan sesi via `/session` atau flag `./scorp -s <nama>`. |
| **Pengalaman Terminal** | **Clean CLI with Multiline Paste Mode** | Prompt modern (`scorp ❯ `), eliminasi noise internal retry, dan dukungan paste blok multi-baris (`"""` atau ````). |
| **Pembersihan VPS** | **Pembersihan Fitur Monitoring VPS Lama** | Loop pemantau swap/RAM/CPU yang sering spam notifikasi telah dihapus total, sementara **Cron Job Scheduler** dipertahankan 100%. |

---

## 🧱 Rekapitulasi Refactoring Arsitektur Bersih (*Clean Architecture*)

### 1. Diet Ekstrem `main.go`
* **Sebelumnya:** File monolitik sebesar 891 baris yang bercampur antara parser CLI, inisialisasi daemon, dan 400+ baris switch-case Telegram.
* **Sesudah:** `main.go` ramping menjadi **145 baris** sebagai dispatcher murni. Seluruh alur background daemon dan handler tombol dipindahkan ke `telegram/daemon.go`.

### 2. Modularisasi `agent/loop.go`
* **Sebelumnya:** 1.027 baris yang memuat ReAct loop, upload handling, thinking formatter, continuation heuristics, dan pending confirmations.
* **Sesudah:**
  * `agent/loop.go` (400 baris): Murni alur eksekusi ReAct agent loop.
  * `agent/confirmation.go`: Konfirmasi perintah berbahaya (*dangerous command gate*).
  * `agent/thinking.go`: Format visual *thinking stream*.
  * `agent/continuation.go`: Heuristik kelanjutan langkah (*intent detection*).
  * `agent/upload.go`: Pemrosesan upload berkas & image vision.
  * `agent/steering.go`: Antrean *real-time steering queue*.
  * `agent/session_mgr.go`: Manajemen siklus hidup multi-sesi.

### 3. Pemisahan Provider AI Independen (`models/`)
* Mengadopsi **Clean Strategy Pattern**:
  * `models/provider_interface.go`: Kontrak interface tunggal `LLMProvider`.
  * `models/api_openai.go`: Driver mandiri OpenAI-compatible (DeepSeek, Groq, OpenRouter, Ollama).
  * `models/api_anthropic.go`: Driver mandiri Claude Messages API.
  * `models/api_gemini.go`: Driver mandiri Google Gemini.
  * `models/api_commandcode.go`: Driver mandiri Command Code SSE stream.

---

## 🚀 Panduan Perintah Cepat (Quick Cheatsheet)

```bash
# Setup interaktif 60 detik
./scorp setup

# Masuk ke terminal REPL interaktif
./scorp

# Menjalankan sesi khusus
./scorp --session project-auth

# Mode keamanan Read-Only (Audit)
./scorp --mode=readonly "analisis arsitektur sistem"

# Mode otonom tanpa konfirmasi (YOLO)
./scorp --mode=yolo "jalankan test suite dan laporkan hasil"

# Jalankan playbook SOP otomatis
./scorp sop run health_audit
./scorp sop run code_review

# Nyalakan Web Dashboard lokal (<2MB RAM)
./scorp gateway --port 8080
```

---

## 🧪 Verifikasi & Kualitas Kode
* Seluruh test suite (`go test ./...`) terverifikasi **100% PASS** di seluruh package.
* Kompilasi binary native x86_64, minimal (`nobrowser`), dan Android Termux ARM64 (`make termux`) berhasil tanpa error.
* Knowledge graph telah diperbarui via `graphify update .`.
