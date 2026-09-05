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

## Rekomendasi Perbaikan

1. **Timeout pada sesi verify CLI** — bungkus operasi verify dengan
   `context.WithTimeout` (misal 30–60 menit), sehingga proses bunuh
   diri bila tidak selesai.
2. **Deteksi heartbeat/stall** — log progres per langkah verify; jika
   tidak ada progres dalam N menit, anggap stall dan keluar dengan
   error.
3. **Guard proses lama** — saat `--cli` start, cek PID file / lock
   untuk session name yang sama; tolak atau gantikan proses lama
   alih-alih menumpuk.
4. **Cleanup saat re-build** — ingat bahwa `go build` mengganti inode;
   proses lama yang masih jalan memakai binary basi. Pertimbangkan
   `pkill -f 'scorp --cli'` sebelum rebuild, atau dokumentasikan
   risikonya.
5. **Opsional: system-level guard** — pasang batas CPU quota via
   systemd slice / `ulimit` untuk CLI berumur panjang.

## Pelajaran Umum

- `(deleted)` di `/proc/<pid>/exe` adalah penanda cepat bahwa proses
  menjalankan binary versi lama.
- Proses dengan parent `systemd --user` + CPU tinggi + tanpa network
  = kandidat kuat orphaned runaway.
- Load average tinggi persisten tanpa aktivitas user yang jelas layak
  diselidiki dengan `ps aux --sort=-%cpu` sebelum laptop dianggap
  "normal panas".
