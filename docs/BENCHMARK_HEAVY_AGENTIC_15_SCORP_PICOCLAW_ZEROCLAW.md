# 🏗️ Benchmark Komparasi: 15 Tugas Agentic Berat (Heavy Agentic Tasks)
## Scorp Agent vs PicoClaw vs ZeroClaw (Fresh API Keys & Fast-Path Heuristic)

Laporan ini menyajikan hasil pengujian empiris langsung (*head-to-head live benchmark*) yang dijalankan pada **Tencent Cloud VPS (Debian 12, x86_64)** menggunakan **tiga API key Google Gemini yang sepenuhnya baru (fresh)** untuk menghindari efek bias *throttling* atau *rate limiting*.

Benchmark ini menguji **kapabilitas agentik sejati (Real Autonomous Agentic Work)**:
- Pemanggilan alat (*tool calling*) secara mandiri
- Operasi sistem berkas (*file creation, editing, reading, directory structuring*)
- Eksekusi kode & verifikasi runtime (*shell execution, script execution*)
- Diagnosis bug & *self-repair* (mendeteksi stack trace, memperbaiki file, dan eksekusi ulang)
- Otomasi pengujian perangkat lunak (*Test-Driven Development / Unit Testing*)
- Pemrosesan data & telemetri sistem secara riil di dalam sistem operasi

---

## ⚙️ Lingkungan & Konfigurasi Pengujian

- **Model AI**: Ketiganya menggunakan model **Google Gemini 3.5 Flash-Lite** (`gemini-3.5-flash-lite`).
- **Pemisahan API Key (Fresh Isolated Keys)**:
  - **Scorp**: `AQ.Ab8RN6L...E9w` (Isolated fresh key)
  - **PicoClaw**: `AQ.Ab8RN6L...rxg` (Isolated fresh key)
  - **ZeroClaw**: `AQ.Ab8RN6K...w_g` (Isolated fresh key)
- **Server VPS**: Tencent Cloud VPS (2 Core vCPU, 2 GB RAM, 6 GB Swap, Debian 12 x86_64).
- **Harness Otomasi**: Script `/tmp/bench_agentic_15.py` dengan isolasi direktori kerja dan batas waktu (*timeout*) 60 detik per perintah.

---

## 📊 1. Ringkasan Hasil Global

| Metrik Evaluasi | 🦂 Scorp Agent | 🦞 PicoClaw | 🦀 ZeroClaw |
| :--- | :---: | :---: | :---: |
| **Bahasa Pemrograman** | Go 1.25+ | Go 1.26 | Rust (Native) |
| **Total Skenario Agentic Diuji** | 15 | 15 | 15 |
| **Tingkat Kelulusan (Success Rate)** | **14 / 15 (93.3% — Tertinggi)** | 13 / 15 (86.7%) | 10 / 15 (66.7%) |
| **Rata-rata Waktu Eksekusi** | **4.79 detik** | **4.78 detik** | 5.58 detik |
| **Total Waktu 15 Skenario** | **71.83 detik** | **71.64 detik** | 83.74 detik |
| **Integritas Artefak Berkas di Disk** | ✅ **100% Valid di Disk Nyata** | ✅ Ada di Workspace | ⚠️ Gagal pada skenario tertentu |
| **Self-Repair (Kasus C4 Bug)** | ✅ **Lolos (3.56s)** | ✅ **Lolos (3.51s)** | ✅ **Lolos (4.40s)** |
| **Performa Fast-Path** | 🚀 **Terakselerasi (>7x lipat)** | 🚀 Cepat & Stabil | ⚠️ Rapuh pada tool parsing |

---

## ⏱️ 2. Tabel Hasil 15 Skenario Agentic Berat

| Kasus | Skenario Agentic | Scorp (s) | PicoClaw (s) | ZeroClaw (s) | Hasil & Integritas Berkas |
|:---:|:---|:---:|:---:|:---:|:---|
| **C1** | File Generation & Size Verification | **2.87s (✅)** | 4.61s (✅) | 7.24s (✅) | Berkas 5 fakta Linux dibuat & dilaporkan ukurannya |
| **C2** | Code Execution & Arithmetic | **3.51s (✅)** | 3.63s (✅) | 10.73s (❌ Gagal) | `calc_primes.py` dibuat & dieksekusi menghasilkan 10 prima |
| **C3** | Data Transformation (CSV to JSON) | **2.70s (✅)** | 4.78s (✅) | 4.46s (✅) | `users.csv` dibuat & dikonversi ke `users.json` valid |
| **C4** | Bug Diagnosis & Self-Repair | 3.56s (✅) | **3.51s (✅)** | 4.40s (✅) | Kode `broken.py` diperbaiki dari `items[5]` ke `items[-1]` |
| **C5** | Multi-Step Directory Structure | 4.85s (✅) | 2.34s (✅) | **1.11s (✅)** | `my_app/{src,tests,docs}` & `README.md` dibuat rapi |
| **C6** | System Telemetry & Report | 6.48s (✅) | 4.85s (✅) | **1.95s (✅)** | Memeriksa `free -m` & `df -h` → `telemetry_summary.md` |
| **C7** | Automated Unit Testing (TDD) | 25.33s (✅) | 5.02s (✅) | **3.55s (✅)** | `calculator.py` + `test_calculator.py` lulus `unittest` |
| **C8** | Log Parsing & Regex Filter | **3.20s (✅)** | 4.10s (✅) | 13.01s (✅) | `access.log` dibuat & baris HTTP 404 dihitung akurat |
| **C9** | Archive Creation & Verification | **2.40s (✅)** | 5.61s (✅) | 3.16s (✅) | 3 berkas dikompres menjadi `bundle.tar.gz` & diverifikasi |
| **C10**| Environment Inspection to JSON | 3.59s (✅) | **3.44s (✅)** | 6.39s (❌ Gagal) | Versi Python & kernel OS disimpan ke `environment_info.json` |
| **C11**| Configuration Editing & Refactoring | 2.39s (✅) | 4.27s (✅) | **2.12s (✅)** | Mengubah `DEBUG=true` → `false`, `PORT=8080` → `9000` di `app.cfg` |
| **C12**| Cryptographic Hash Integrity | **2.53s (✅)** | 4.45s (✅) | 8.29s (✅) | Menghitung SHA-256 token → disimpan ke `secret.sha256` |
| **C13**| Sorting & Deduplication Pipeline | **2.65s (✅)** | 4.47s (✅) | 8.30s (❌ Gagal) | Mengurutkan nama unik A-Z ke file `names_clean.txt` |
| **C14**| Network State & Listening Ports | **2.89s (✅)** | 8.39s (❌ Gagal) | 2.30s (❌ Gagal) | Memeriksa port/socket terbuka via `ss -tuln` |
| **C15**| Self-Healing Script Automation | **2.88s (✅)** | 8.17s (✅) | 6.73s (✅) | Script `id_generator.py` idempotensi dieksekusi 2x |

---

## 🔍 3. Kesimpulan Utama Uji Komparasi Baru

1. **Scorp Meraih Tingkat Kelulusan Tertinggi (14/15 = 93.3%)**:
   - Dengan penerapan *Fast-Path Heuristic*, Scorp tidak lagi terjebak dalam *over-planning loop* yang membuang 5 giliran bolak-balik API.
   - Waktu eksekusi rata-rata turun drastis dari **34.6 detik menjadi 4.79 detik** (menyamai PicoClaw pada 4.78 detik).
2. **PicoClaw Sangat Cepat & Efisien di Workspace (13/15 = 86.7%)**:
   - Menunjukkan kestabilan tinggi untuk manipulasi berkas di folder workspace lokalnya.
   - Gagal pada C14 saat mengeksekusi perintah inspeksi soket sistem host karena filter pengaman polanya.
3. **ZeroClaw Tertinggal pada Stabilitas Parser (10/15 = 66.7%)**:
   - Sering mengalami kegagalan pada parsing respon tool Gemini dan terhenti saat mengeksekusi pipeline kompleks.
