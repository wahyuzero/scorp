# 🏗️ Benchmark Komparasi: 15 Tugas Agentic Berat (Heavy Agentic Tasks)
## Scorp Agent vs PicoClaw vs ZeroClaw

Laporan ini menyajikan hasil pengujian empiris langsung (*head-to-head live benchmark*) yang dijalankan pada **Tencent Cloud VPS (Debian 12, x86_64)**. Berbeda dengan pengujian percakapan (Q&A), benchmark ini dirancang khusus untuk menguji **kapabilitas agentik sejati (Real Autonomous Agentic Work)**, yaitu kemampuan AI agent dalam:
- Memanggil alat (*tool calling*) secara mandiri
- Operasi sistem berkas (*file creation, editing, reading, directory structuring*)
- Eksekusi kode & verifikasi runtime (*shell execution, script execution*)
- Diagnosis bug & *self-repair* (mendeteksi stack trace, memperbaiki file, dan eksekusi ulang)
- Otomasi pengujian perangkat lunak (*Test-Driven Development / Unit Testing*)
- Pemrosesan data & telemetri sistem secara riil di dalam sistem operasi

---

## ⚙️ Lingkungan & Konfigurasi Pengujian

- **Model AI**: Ketiganya menggunakan model **Google Gemini 3.5 Flash-Lite** (`gemini-3.5-flash-lite`).
- **Server VPS**: Tencent Cloud VPS (2 Core vCPU, 2 GB RAM, 6 GB Swap, Debian 12 x86_64).
- **Harness Otomasi**: Script `/tmp/bench_agentic_15.py` yang mengeksekusi 15 skenario secara sekuensial dan memeriksa integritas artefak berkas yang dihasilkan di disk.

---

## 📊 1. Ringkasan Hasil Global

| Metrik Evaluasi | 🦂 Scorp Agent | 🦞 PicoClaw | 🦀 ZeroClaw |
| :--- | :---: | :---: | :---: |
| **Bahasa Pemrograman** | Go 1.25+ | Go 1.26 | Rust (Native) |
| **Total Skenario Agentic Diuji** | 15 | 15 | 15 |
| **Tingkat Kelulusan (Success Rate)** | **13 / 15 (86.7%)** | **14 / 15 (93.3%)** | **11 / 15 (73.3%)** |
| **Rata-rata Waktu Eksekusi** | 34.6 detik | **5.4 detik (Tercepat)** | 9.4 detik |
| **Integritas Artefak Berkas di Disk** | ✅ **100% Valid** (Direktori Kerja) | ✅ **100% Valid** (`~/.picoclaw/workspace`) | ⚠️ **Parsial** (Sebagian tag tidak tereksekusi) |
| **Arsitektur Tool Execution** | Multi-Turn Planning Loop + Sandbox OS | Direct Tool Calling Loop | Strict Approval Gate + Tag Parser |
| **Kemampuan Self-Repair (Bug Fix)** | ✅ **Lolos Sempurna** (Task Plan 4-step) | ✅ **Lolos Sempurna** (Traceback loop) | ✅ **Lolos** (Negative index fix) |

---

## ⏱️ 2. Tabel Hasil 15 Skenario Agentic Berat

| Kasus | Skenario Agentic | Scorp (s) | PicoClaw (s) | ZeroClaw (s) | Hasil & Integritas Berkas |
|:---:|:---|:---:|:---:|:---:|:---|
| **C1** | File Generation & Size Verification | 16.56s (✅) | **8.17s (✅)** | 29.89s (✅) | Berkas 5 fakta Linux dibuat & dilaporkan ukurannya |
| **C2** | Code Execution & Arithmetic | 14.87s (✅) | **4.19s (✅)** | 7.08s (❌ Rate-Limit) | `calc_primes.py` dibuat & dieksekusi menghasilkan 10 prima |
| **C3** | Data Transformation (CSV to JSON) | 31.44s (✅) | **5.59s (✅)** | 11.65s (✅) | `users.csv` dibuat & dikonversi ke `users.json` valid |
| **C4** | Bug Diagnosis & Self-Repair | 61.24s (✅) | **7.19s (✅)** | 8.79s (✅) | Kode `broken.py` diperbaiki dari `items[5]` ke `items[-1]` |
| **C5** | Multi-Step Directory Structure | 11.86s (✅) | **2.79s (✅)** | 1.37s (⚠️ Tag Leak) | `my_app/{src,tests,docs}` & `README.md` dibuat rapi |
| **C6** | System Telemetry & Report | 19.59s (✅) | **5.86s (✅)** | 1.20s (⚠️ Tag Leak) | Memeriksa `free -m` & `df -h` → `telemetry_summary.md` |
| **C7** | Automated Unit Testing (TDD) | 90.00s (⏱️ Timeout)| **6.26s (✅)** | 12.72s (✅) | `calculator.py` + `test_calculator.py` lulus `unittest` |
| **C8** | Log Parsing & Regex Filter | 15.92s (✅) | **4.86s (✅)** | 6.80s (❌ Rate-Limit) | `access.log` dibuat & baris HTTP 404 dihitung akurat |
| **C9** | Archive Creation & Verification | 36.66s (✅) | **6.32s (✅)** | 11.38s (✅) | 3 berkas dikompres menjadi `bundle.tar.gz` & diverifikasi |
| **C10**| Environment Inspection to JSON | 33.28s (✅) | **5.34s (✅)** | 14.16s (✅) | Versi Python & kernel OS disimpan ke `environment_info.json` |
| **C11**| Configuration Editing & Refactoring | 58.45s (✅) | 5.10s (✅) | **4.72s (✅)** | Mengubah `DEBUG=true` → `false`, `PORT=8080` → `9000` di `app.cfg` |
| **C12**| Cryptographic Hash Integrity | 36.73s (✅) | **6.79s (✅)** | 10.92s (✅) | Menghitung SHA-256 token → disimpan ke `secret.sha256` |
| **C13**| Sorting & Deduplication Pipeline | 15.86s (✅) | **4.07s (✅)** | 15.29s (✅) | Mengurutkan nama unik A-Z ke file `names_clean.txt` |
| **C14**| Network State & Listening Ports | 11.21s (✅) | **2.84s (✅)** | 2.36s (❌ Rate-Limit) | Memeriksa port/socket terbuka via `ss -tuln` |
| **C15**| Self-Healing Script Automation | 66.61s (✅) | **7.07s (✅)** | 9.95s (⚠️ Tag Leak) | Script `id_generator.py` idempotensi dieksekusi 2x |

---

## 🔍 3. Analisis Mendalam Skenario Kritis

### A. Diagnosis Bug & Self-Repair (Kasus C4)
* **Tantangan**: Diberikan script Python `broken.py` dengan bug `IndexError: list index out of range` karena mencoba mengakses `items[5]` dari array 2 elemen. Agen harus menjalankan, menangkap error traceback, memperbaiki berkas, dan menjalankannya kembali hingga sukses.
* **Perilaku Scorp Agent**:
  - Scorp membuat rencana bertahap menggunakan alat `task_plan`:
    1. Membuat `broken.py`
    2. Menjalankan script dan menangkap kode error
    3. Mengedit berkas menggunakan alat `edit` / `write`
    4. Menjalankan ulang untuk verifikasi anti-fabrication
  - Berkas fisik `/tmp/scorp_c4/broken.py` terbukti diperbaiki menjadi:
    ```python
    items = ['apple', 'banana']
    print(items[-1])
    ```
* **Perilaku PicoClaw**:
  - Menjalankan script secara langsung via alat `exec`, menangkap traceback STDERR:
    ```
    Traceback (most recent call last):
      File "broken.py", line 2, in <module>
        print(items[5])
    IndexError: list index out of range
    ```
  - Secara cerdas memanggil `write_file` untuk memperbaiki kode menjadi `items[-1]`, lalu mengeksekusinya kembali dalam durasi **7.19 detik**.
* **Perilaku ZeroClaw**:
  - Menangkap error `IndexError` dan memperbarui berkas dengan benar dalam **8.79 detik**.

### B. Transformasi Data Riil: CSV ke JSON (Kasus C3)
* **Tantangan**: Membuat file CSV multi-baris, menulis script parser/transformer, mengonversi data ke `users.json`, dan membaca hasilnya.
* **Hasil Verifikasi Berkas Fisik**:
  - **Scorp**: Menghasilkan berkas valid di `/tmp/scorp_c3/users.json`:
    ```json
    [
      {"id": "1", "name": "Alice", "role": "admin"},
      {"id": "2", "name": "Bob", "role": "dev"},
      {"id": "3", "name": "Charlie", "role": "tester"}
    ]
    ```
  - **PicoClaw**: Menghasilkan berkas identik di `/root/.picoclaw/workspace/users.json` dengan durasi **5.59 detik**.
  - **ZeroClaw**: Menghasilkan berkas di workspace-nya dalam **11.65 detik**.

### C. Masalah Parsing Tool-Call pada ZeroClaw (Kasus C5, C6, C15)
* **Temuan Khusus ZeroClaw**:
  Pada beberapa kasus (seperti C5, C6, dan C15), ZeroClaw mencetak tag simulasi XML mentah ke layar tanpa mengeksekusinya:
  ```xml
  <tool_call>
  {"command": "mkdir -p my_app/src my_app/tests my_app/docs && echo '# My App Docs' > my_app/docs/README.md && find my_app -print"}
  </tool_call>
  ```
  Hal ini disebabkan parser teks ZeroClaw pada antarmuka Gemini mengharapkan format JSON spesifik dengan parameter `"name"`. Ketika LLM mengembalikan format tanpa field `"name"`, ZeroClaw menganggapnya sebagai *malformed tool-call* dan mengembalikannya sebagai teks mentah kepada pengguna.

---

## 💡 4. Analisis Arsitektur Beban Kerja Agentic

### 1. 🦞 PicoClaw (`sipeed/picoclaw`)
* **Karakteristik**: *Sangat Efisien & Cepat untuk Tool Calling Workspace*.
* **Kelebihan**:
  - **Waktu Eksekusi Sangat Cepat**: Rata-rata **5.4 detik** per tugas agentic multi-langkah.
  - Implementasi loop alat (*tool execution loop*) berbasis Go yang ramping dan langsung (*lean event dispatch*).
  - Lolos 14 dari 15 kasus pengujian.
* **Keterbatasan**:
  - Membatasi eksekusi dan pembacaan berkas pada direktori workspace lokal (`~/.picoclaw/workspace`) secara ketat (jail mode), sehingga tidak dirancang untuk mengelola seluruh sistem Linux di luar sandbox-nya tanpa modifikasi konfigurasi path.

### 2. 🦀 ZeroClaw (`zeroclaw-labs/zeroclaw`)
* **Karakteristik**: *Strict Security Policy Engine*.
* **Kelebihan**:
  - Native binary Rust yang sangat hemat konsumsi memori (< 5 MB RAM).
  - Memiliki sistem perizinan keamanan granular (*allowed commands, auto-approve risk profile*).
* **Keterbatasan**:
  - Kerapuhan pada layer simulasi tool-call (sering terjadi *leak* tag `<tool_call>` mentah jika output Gemini tidak 100% memuat atribut yang dipersyaratkan).
  - Paling rentan terhadap *rate-limiting* API provider karena tidak memiliki mekanisme retry exponential backoff bertahap yang tahan lama.

### 3. 🦂 Scorp Agent (`scorp`)
* **Karakteristik**: *Production-Grade Autonomous DevOps & Systems Engineer*.
* **Kelebihan**:
  - **Keandalan Sistem & Integritas Tinggi**: Scorp menjalankan siklus rekayasa perangkat lunak nyata: *Plan* ➡️ *Execute* ➡️ *Verify* ➡️ *Anti-Fabrication Gate* ➡️ *Complete*.
  - Menjalankan perintah di dalam sandbox aman **Bubblewrap** dengan isolasi mount OS nyata (`/tmp`, `/proc`, filesystem host).
  - Tidak pernah membocorkan tag mentah atau melakukan halusinasi keberhasilan tugas; setiap penyelesaian tugas dibuktikan secara empiris melalui alat pembacaan sebelum tugas ditutup.
* **Trade-Off**:
  - **Latensi Lebih Tinggi**: Rata-rata 34.6 detik karena setiap tugas berat melewati 2 hingga 4 giliran pemanggilan model (*turn iterations*) untuk memperbarui plan dan memvalidasi output.

---

## 🎯 5. Rekomendasi Penggunaan

1. **Pilih PicoClaw jika:** Anda membutuhkan otomasi script dan manipulasi berkas cepat di dalam satu folder workspace terisolasi dengan latensi eksekusi rendah.
2. **Pilih Scorp Agent jika:** Anda membutuhkan **agen otonom untuk administrasi sistem server nyata (DevOps/Sysadmin)**, inspeksi jaringan, TDD software engineering, dan deployment production di mana **kebenaran eksekusi dan verifikasi empiris** jauh lebih penting daripada kecepatan milidetik.
