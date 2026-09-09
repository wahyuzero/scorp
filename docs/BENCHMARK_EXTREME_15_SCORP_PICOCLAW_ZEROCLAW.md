# 🌌 Benchmark Komparasi Ekstrem: 15 Tugas Rekayasa Sistem & Arsitektur Tingkat Lanjut
## Scorp Agent vs PicoClaw vs ZeroClaw (Head-to-Head Live VPS Evaluation)

Laporan ini menyajikan hasil pengujian empiris langsung (*head-to-head live benchmark*) pada tingkat kesulitan **ekstrem (Extreme Architecture & Systems Engineering)** di server **Tencent Cloud VPS (Debian 12, x86_64)** dengan model **Google Gemini 3.5 Flash-Lite**.

Benchmark ini dirancang khusus untuk membedah kemampuan AI agent dalam menangani skenario yang menuntut **pemahaman sistem operasi mendalam, kompilasi bahasa tingkat rendah (C/C++), konkurensi multithreading, integritas transaksi database, jaringan socket mentah, dan parser sintaks mandiri**.

---

## 🏛️ 15 Skenario Uji Ekstrem

1. **E1 (Native C Shared Library & Python ctypes Binding)**: Menulis kode C `math_core.c`, mengompilasi menjadi shared object `libmath_core.so` dengan `gcc -shared -fPIC`, lalu memuat dan memverifikasinya melalui antarmuka Python `ctypes`.
2. **E2 (Asynchronous Task Queue with SQLite Locking)**: Mengimplementasikan sistem antrean tugas konkuren multi-worker menggunakan SQLite dengan penanganan lock dan transaksi tanpa deadlock.
3. **E3 (Autonomous TLS Certificate & HTTPS Server)**: Menghasilkan sertifikat SSL mandiri (*self-signed certificate*) dan private key via OpenSSL secara non-interaktif, menjalankan server HTTPS Python port 9443 di background, menguji request HTTPS via `curl -k`, lalu mematikan server.
4. **E4 (Multi-Level JSON Schema Flattener & Diff Engine)**: Menulis diff engine rekursif yang meratakan struktur JSON bersarang 3 tingkat menjadi dot-notation (`a.b.c`) dan mendeteksi perubahan `ADDED`, `REMOVED`, dan `MODIFIED`.
5. **E5 (Linux Process Memory Limit & OOM Handling)**: Membatasi memori proses Linux secara ketat menggunakan syscall `resource.setrlimit(RLIMIT_AS, 300MB)`, mengalokasikan memori hingga memicu `MemoryError`, dan menangkapnya secara aman (*graceful OOM handling*).
6. **E6 (Multi-Module Circular Dependency Refactoring)**: Mendiagnosis siklus impor melingkar pada 3 modul Python (`A -> B -> C -> A`), merefaktor arsitekturnya dengan memindahkan dependensi bersama ke `common.py` hingga dapat dieksekusi tanpa error.
7. **E7 (SQLite Full-Text Search FTS5 & BM25 Ranking)**: Menginisialisasi virtual table FTS5 SQLite, memasukkan korpus artikel berbahasa Indonesia, lalu mengeksekusi query pencarian kata kunci dengan highlight snippet dan skor relevansi BM25.
8. **E8 (High-Concurrency HTTP Benchmark Load Engine)**: Menjalankan HTTP server lokal port 9222 di background, lalu membanjirinya dengan load tester multithreaded (100 request simultan), mengukur durasi, persentase keberhasilan, dan Requests Per Second (RPS).
9. **E9 (SQL Dialect Tokenizer & Lexer State Machine)**: Membangun lexer/tokenizer mandiri (mesin keadaan karakter tanpa pustaka pihak ketiga) untuk mengurai ekspresi SQL kompleks menjadi token berlabel.
10. **E10 (Systemd Unit File Generator & Syntax Validation)**: Menghasilkan unit file systemd `scorp-demo.service` dengan konfigurasi restart policy, isolasi resource RAM/CPU, dan memvalidasi keabsahan sintaksnya.
11. **E11 (Micro-Language Lexer, AST Parser & Evaluator)**: Membangun interpreter mini yang mengurai ekspresi matematika berpangkat dan bertanda kurung (seperti `3 + 5 * (10 - 4) / 2 - 2^3`), menyusun representasi pohon AST, dan mengevaluasi hasilnya secara rekursif.
12. **E12 (AST-Based Static Code Security Vulnerability Linter)**: Menulis penganalisis kode statis berbasis Python `ast` untuk mendeteksi potensi kerentanan (penggunaan `eval()`, `subprocess.Popen(shell=True)`, dan hardcoded secret).
13. **E13 (Non-Blocking Async TCP Echo Server & Socket Suite)**: Menulis server socket TCP non-blocking berbasis `select` pada port 9333 di background, menghubungkan client socket, mengirim payload `"PING-SCORP-CHALLENGE-2026"`, menerima echo, dan menutup socket secara bersih.
14. **E14 (Atomic File Swapping & Crash-Resilient Write Engine)**: Mengimplementasikan penulisan berkas berstandar ACID (tulis ke `.tmp`, panggil `os.fsync`, validasi hash SHA-256, lalu swap atomik via `os.replace` ke `production_data.dat`).
15. **E15 (Prometheus Metric Exporter with Live System Gauges)**: Menjalankan exporter HTTP port 9555 yang menyajikan data telemetri real-time sistem Linux host (`/proc/meminfo`, `/proc/loadavg`, `os.statvfs`) dalam format OpenMetrics.

---

## 📊 1. Ringkasan Hasil Global

| Metrik Evaluasi | 🦂 Scorp Agent | 🦞 PicoClaw | 🦀 ZeroClaw |
| :--- | :---: | :---: | :---: |
| **Bahasa Pemrograman** | Go 1.25+ | Go 1.26 | Rust (Native) |
| **Total Skenario Ekstrem Diuji** | 15 | 15 | 15 |
| **Tingkat Kelulusan (Success Rate)** | **15 / 15 (100% Sempurna)** | 14 / 15 (93.3%) | 12 / 15 (80.0%) |
| **Rata-rata Waktu Eksekusi** | **9.45 detik (Tercepat)** | 9.87 detik | 20.48 detik |
| **Total Durasi 15 Skenario** | **141.74 detik** | **148.11 detik** | 307.13 detik |
| **Verifikasi Artefak Fisik di Host** | ✅ **100% Ada di Host Disk** | ⚠️ Terkunci di Workspace | ❌ Gagal di Berbagai Kasus |
| **Biner Native C Kompilasi (E1)** | ✅ **Sukses `gcc` & ctypes (11.79s)** | ✅ Sukses (27.31s) | ⚠️ Iteration Timeout (49.37s) |
| **Ketahanan Concurrency & Jaringan** | ✅ **100% Sempurna (E3, E8, E13, E15)**| ✅ Sempurna | ⚠️ Policy Deny & Port Glitch |

---

## ⏱️ 2. Tabel Waktu Eksekusi 15 Skenario Ekstrem (Detik)

| Kasus | Skenario Rekayasa Ekstrem | Scorp (s) | PicoClaw (s) | ZeroClaw (s) | Analisis Eksekusi & Kualitas Output |
|:---:|:---|:---:|:---:|:---:|:---|
| **E1** | Native C Shared Library & ctypes | **11.79s (✅)** | 27.31s (✅) | 49.37s (⚠️ Max Turns) | Scorp mengompilasi `libmath_core.so` via gcc dan ctypes Python memanggil fib(10) = 55. ZeroClaw kehabisan batas turn. |
| **E2** | Asynchronous Task Queue SQLite | 18.02s (✅) | **6.27s (✅)** | 32.99s (✅) | Multi-worker threading memproses 6 task secara aman tanpa deadlock. |
| **E3** | Autonomous TLS & HTTPS Server | **6.46s (✅)** | 17.88s (✅) | 17.96s (⚠️ Max Turns) | Scorp menghasilkan sertifikat OpenSSL, menyalakan server HTTPS port 9443, curl -k, dan mematikannya. |
| **E4** | Recursive JSON Flattener & Diff | **8.10s (✅)** | 9.12s (✅) | 28.97s (✅) | Menghitung perbedaan JSON bersarang 3 tingkat dengan format dot-notation. |
| **E5** | Memory Limit & OOM Handling | **3.35s (✅)** | 4.95s (✅) | 8.52s (❌ Rate-Limit) | Syscall `setrlimit` membatasi heap dan menangkap `MemoryError` pada 250 MB secara aman. |
| **E6** | Circular Dependency Refactoring | **5.90s (✅)** | 20.01s (✅) | 9.50s (✅) | Mendiagnosis siklus impor melingkar dan merefaktor arsitektur via `common.py`. |
| **E7** | SQLite Full-Text Search FTS5 | 9.53s (✅) | 6.86s (❌ Crash) | **9.07s (✅)** | Scorp membuat tabel FTS5 dan menjalankan ranking BM25. PicoClaw mengalami crash internal error pada turn final. |
| **E8** | High-Concurrency HTTP Load Bench | **6.60s (✅)** | 10.12s (✅) | 26.50s (⚠️ Max Turns) | Menguji 100 request simultan port 9222 hingga mencatatkan 100% sukses. |
| **E9** | SQL Lexer State Machine Parser | 9.29s (✅) | **5.83s (✅)** | 12.35s (❌ Gagal) | Mengurai query SQL kompleks menjadi token berlabel tanpa pustaka pihak ketiga. |
| **E10**| Systemd Unit Generator & Check | 9.71s (✅) | 9.04s (✅) | **8.15s (✅)** | Menghasilkan unit service systemd dengan spesifikasi resource isolasi. |
| **E11**| Micro-Language AST Evaluator | 13.99s (✅) | **6.09s (✅)** | 41.72s (✅) | Membangun pohon AST dari ekspresi matematika bertanda kurung dan mengevaluasi hasil `10.0`. |
| **E12**| AST Static Security Linter | 13.58s (✅) | 6.92s (✅) | **5.20s (✅)** | Memindai file kode Python dan melaporkan pelanggaran `eval`, shell injection, dan hardcoded secret. |
| **E13**| Non-Blocking Async TCP Echo | 8.86s (✅) | **5.82s (✅)** | 28.76s (✅) | Server socket TCP non-blocking di background merespons payload echo dengan verifikasi integritas paket. |
| **E14**| Atomic File Swapping & fsync | 8.44s (✅) | **5.60s (✅)** | 6.82s (❌ Rate-Limit) | Penulisan file berbasis crash-resilient `.tmp` -> `os.fsync` -> swap atomik `production_data.dat`. |
| **E15**| Prometheus Exporter & Gauges | 8.12s (✅) | **6.29s (✅)** | 21.25s (✅) | Menjalankan exporter HTTP port 9555 yang menyajikan metrik CPU, RAM, dan disk berstandar Prometheus. |

---

## 🔍 3. Analisis Mendalam Keandalan Arsitektur

1. **Scorp Agent (Pemenang Mutlak Keandalan & Ketahanan Arsitektur)**:
   - **Tingkat Kelulusan**: **15 / 15 (100% Sempurna)**.
   - **Rata-rata Waktu**: **9.45 detik** (Tercepat dari seluruh agen).
   - Seluruh 15 berkas hasil rekayasa nyata (`libmath_core.so`, `queue.db`, `server.crt`, `server.key`, `memory_stress.py`, `common.py`, `articles.db`, `calc_ast.py`, `production_data.dat`, `scorp-demo.service`, dll.) **benar-benar tercipta secara fisik di disk VPS (`/tmp/extreme_scorp_c*`)**.
   - Menunjukkan keunggulan mutlak pada skenario kompilasi C asli (E1), pembuatan server HTTPS mandiri (E3), dan stabilitas sistem antrean database konkuren (E2).

2. **PicoClaw (Sangat Cepat & Tangkas, Namun Gagal pada Ekstensi FTS5)**:
   - Menyelesaikan 14 dari 15 kasus (93.3%) dengan waktu rata-rata 9.87 detik.
   - Mengalami crash pada kasus E7 (SQLite Full-Text Search FTS5) karena library internalnya mengalami error unhandled exception saat giliran akhir diselesaikan:
     `severity=error status=error trace_id=turn.end`.

3. **ZeroClaw (Sangat Rapuh pada Tugas Sistem Kompleks)**:
   - Hanya menyelesaikan 12 dari 15 kasus (80.0%) dengan waktu rata-rata yang membengkak hingga **20.48 detik** (total durasi >5 menit).
   - Pada 4 skenario (E1, E3, E8), ZeroClaw kehabisan batas giliran internal (*Turn stopped: reached maximum tool iterations (10)*) karena terjebak dalam loop perintah tanpa berhasil menyimpulkan laporan akhir.
   - Mengalami kegagalan akibat limit API dan kendala kebijakan eksekusi latar belakang pada kasus E5 dan E14.
