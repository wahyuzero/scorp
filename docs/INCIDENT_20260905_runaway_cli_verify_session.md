# Incident Report: Runaway CLI Verify Session

- **Tanggal insiden**: 4–5 September 2026
- **Dokumentasi**: 5 September 2026, 14:05 WIB
- **Severity**: Medium (tidak ada data loss, murni pemborosan resource)
- **Status**: Resolved (proses di-kill manual)

## Ringkasan

Dua proses `scorp --cli` ditemukan berjalan runaway selama ±15 jam
sejak malam 4 September, mengonsumsi ~170% CPU masing-masing (total
~340% dari kapasitas core) tanpa melakukan pekerjaan yang berguna.

| PID | Command | CPU | CPU-time total | Umur proses |
|---|---|---|---|---|
| 87914 | `./scorp --cli --session=verify_all_ok` | ~170% | 25j 44m | 15j 07m |
| 88650 | `./scorp --cli --session=verify_all_final` | ~170% | 25j 49m | 15j 07m |

Total pemborosan: **±51 core-hours** untuk busy-loop kosong.
Load average sistem sempat menyentuh 4.7–6.1.

## Ciri-ciri Runaway yang Terdeteksi

1. **Binary `(deleted)`** — `readlink /proc/<pid>/exe` menunjuk ke
   `/home/wxsys/Project/scorp/scorp (deleted)`. Binary sudah di-rebuild
   (07:46) tapi proses lama masih memegang inode versi malam sebelumnya.
   Build baru tidak menyentuh proses yang sudah jalan.
2. **Ter-adopsi systemd-user** — parent proses hanya
   `systemd --user`; tab konsole asal proses sudah ditutup, tidak ada
   yang memantau hasilnya.
3. **Tanpa network & tanpa child process** — tidak ada koneksi TCP
   aktif dan tidak ada child, mengindikasikan busy-loop murni.
4. **Nama session menyesatkan** — `verify_all_ok` dan
   `verify_all_final` tidak pernah menyelesaikan verifikasinya.

## Akar Masalah (dugaan)

Session CLI verify dari malam 4 September mengalami *deadloop* /
kondisi tanpa timeout. Tidak ada mekanisme yang menghentikan proses
ketika pekerjaan tidak kunjung selesai, sehingga proses bertahan
hingga di-kill manual keesokan harinya.

Log aplikasi (`~/.scorp/scorp.log`) menunjukkan aktivitas normal hingga
07:48 (model loading, cost tracking, tool execution), namun proses
verify lama tidak berhenti setelah build berikutnya.

## Penanganan

```bash
kill 87914 88650
```

Verifikasi pasca-kill:
- CPU usage turun dari 39.7% → 6.0% (idle 91.7%)
- Load average menurun bertahap dari 4.13
- Tidak ada data loss; binary hasil build 07:46 tetap utuh di disk

## Rekomendasi Perbaikan & Status Implementasi

1. **Timeout pada sesi verify CLI** — [RESOLVED] Loop eksekusi agent (`agent/loop.go`) kini dibungkus dengan `context.WithTimeout` default 30 menit (bisa dikonfigurasi lewat env `SCORP_MAX_TURN_TIMEOUT`), sehingga proses tidak bisa deadloop selamanya.
2. **Deteksi heartbeat/stall** — [RESOLVED] Diatur lewat context deadline pada HTTP client, retry retries max 4, dan tick active skills.
3. **Guard proses lama (Session Lock)** — [RESOLVED] Mekanisme exclusive process lock berbasis kernel advisory file lock (`cli_lock.go`) kini aktif di CLI. Jika ada proses CLI lain yang mencoba mengakses session yang sama (seperti insiden `verify_all_ok` & `verify_all_final`), proses kedua langsung ditolak dengan pesan jelas dan PID proses aktif.
4. **Cleanup saat re-build** — [RESOLVED] `pkill` dan advisory lock mencegah collision file history dan orphan process.
5. **State Machine Non-Heuristic** — [RESOLVED] Menghapus regex/grammar loops dan menggantikannya dengan explicit tool contract `complete_task`.


## Pelajaran Umum

- `(deleted)` di `/proc/<pid>/exe` adalah penanda cepat bahwa proses
  menjalankan binary versi lama.
- Proses dengan parent `systemd --user` + CPU tinggi + tanpa network
  = kandidat kuat orphaned runaway.
- Load average tinggi persisten tanpa aktivitas user yang jelas layak
  diselidiki dengan `ps aux --sort=-%cpu` sebelum laptop dianggap
  "normal panas".
