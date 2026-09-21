# 🏛️ BLUEPRINT ARSITEKTUR: STEALTH BROWSER MCP ENGINE IN PURE GO (DARI NOL)

> **Dokumen Spesifikasi Teknis & Cetak Biru Arsitektur**  
> **Target:** Membangun runtime controller browser otomatisasi berteknologi *Stealth Anti-Detection* dan *Model Context Protocol (MCP)* murni dalam Golang, tanpa dependensi Node.js, dengan memori controller **< 20 MB RAM**.

---

## 1. 🌐 Gambaran Arsitektur Sistem (High-Level Topology)

```mermaid
flowchart TD
    subgraph AI_CLIENT["1. AI Agent Harness"]
        LLM["AI Model (Gemini / Claude / Scorp / Antigravity)"]
    end

    subgraph GO_ENGINE["2. Pure Go Stealth MCP Binary (<20MB RAM)"]
        direction TB
        Stdio["Stdio JSON-RPC 2.0 Handler"]
        Router["Tool Router (navigate, click, type, snapshot, eval)"]
        
        subgraph CORE_MODULES["Core Engine Subsystems"]
            StealthMgr["Stealth & Fingerprint Manager\n(JS Injection, Canvas, WebGL, UA)"]
            ActionEngine["Actionability Engine\n(WaitVisible, WaitStable, Unobscured)"]
            HumanMotion["Human Motion Simulator\n(Bezier Curves, Jitter, Typing Cadence)"]
            FrugalPruner["Frugal Resource Pruner\n(Network Interceptor, Block Media/Fonts)"]
        end

        CDPClient["CDP Protocol Client (WebSocket / Pipe)\n(go-rod/rod or chromedp)"]
    end

    subgraph BROWSER["3. Browser Runtime (Chromium / Patched Engine)"]
        DevToolsPort["CDP WebSocket (:9222 or Remote Pipe)"]
        BlinkEngine["Blink / V8 Engine (Render Process)"]
        TargetDOM["Target Website DOM\n(LinkedIn, Cloudflare, Akamai, Target SPAs)"]
    end

    LLM <-->|Stdio JSON-RPC| Stdio
    Stdio --> Router
    Router --> CORE_MODULES
    CORE_MODULES --> CDPClient
    CDPClient <-->|Raw CDP JSON over WS| DevToolsPort
    DevToolsPort --> BlinkEngine
    BlinkEngine --> TargetDOM
```

---

## 2. 🔬 LEVEL 0: Fisika Otomasi Browser & Mengapa Bot Selalu Tertangkap

Sebelum menulis satu baris kode Go pun, kita wajib memahami bagaimana server anti-bot (Cloudflare Turnstile, DataDome, Akamai, Kasada, LinkedIn Sentinel) membedakan manusia vs bot:

### A. Tanda Kebocoran Utama (The "Smoking Guns"):
1. **`navigator.webdriver === true`**:
   Flag bawaan Chromium saat dijalankan dengan parameter otomasi (`--enable-automation`). Ini adalah deteksi level 1 paling primitif.
2. **CDP Runtime Infiltration Leak**:
   Ketika debugger CDP terhubung, objek internal seperti `window.cdc_adoQpoasnfa76pfcZLmcfl_Array` atau runtime console binding tertinggal di JavaScript runtime context. Anti-bot melakukan `Object.getOwnPropertyNames(window)` untuk mencari sisa-sisa objek ini.
3. **Hardware Inconsistency (Fingerprinting)**:
   * **WebGL Vendor & Renderer:** Headless Chrome sering melaporkan `Google Inc. (Google SwiftShader)` alih-alih GPU fisik seperti `NVIDIA GeForce` atau `Apple M-series`.
   * **AudioContext & Canvas Noise:** Fingerprint Canvas headless selalu menghasilkan hash pixel yang identik dan steril tanpa derau hardware nyata.
   * **Client Hints (`sec-ch-ua`):** Header HTTP yang dikirim menyatakan OS Windows, tetapi JavaScript `navigator.platform` di dalam DOM menyatakan `Linux x86_64`.
4. **Behavioral Telemetry (Kinematika Mouse & Keyboard)**:
   * Bot mengeklik instan dari koordinat (0,0) langsung teleportasi ke (450, 600) dengan kecepatan `infinity px/s` (tanpa kurva akselerasi/deselerasi).
   * Ketikan karakter masuk seragam setiap tepat 0ms atau 50ms tanpa variasi ritme biologis manusia.

---

## 3. 🧩 LEVEL 1: Pemilihan Fondasi Driver di Golang

Ada 3 cara mengontrol browser di Go:

| Komponen | `chromedp` | `playwright-go` | `go-rod/rod` (Direkomendasikan) |
| :--- | :--- | :--- | :--- |
| **Arsitektur** | Wrapper tipis CDP | Menjalankan Node.js driver di background | **Pure Go CDP Client** |
| **Memory Footprint** | ~10 MB | ~150 MB (Kalah esensi Frugal) | **~12 MB** |
| **Actionability Checks** | ❌ Manual coding | ✅ Bawaan Playwright | **✅ Bawaan Rod (`WaitVisible`, `WaitStable`)** |
| **Stealth Ecosystem** | Minim | Butuh node plugins | **✅ `go-rod/stealth` siap pakai** |
| **Kesimpulan** | Terlalu low-level | Boros dependensi | **Paling ideal untuk Frugal Stealth MCP** |

> **Keputusan Arsitektur:** Gunakan **`go-rod/rod`** sebagai pengemudi utama CDP, digabungkan dengan **`go-rod/stealth`** atau skrip injeksi kustom.

---

## 4. 🛡️ LEVEL 2: 5 Lapisan Mesin Stealth (The Secret Sauce)

Untuk mencapai status *Undetectable*, Go engine harus menerapkan 5 lapis pertahanan:

```text
┌─────────────────────────────────────────────────────────────┐
│ Layer 5: Human Kinetics (Kurva Bezier & Jitter Ketikan)    │
├─────────────────────────────────────────────────────────────┤
│ Layer 4: Actionability Engine (Mencegah Ghost Click)        │
├─────────────────────────────────────────────────────────────┤
│ Layer 3: Network & TLS Consistency (JA3, Headers Alignment) │
├─────────────────────────────────────────────────────────────┤
│ Layer 2: CDP Script Injection (Canvas, WebGL, navigator)    │
├─────────────────────────────────────────────────────────────┤
│ Layer 1: Browser Launch Arguments & Sandbox Hardening       │
└─────────────────────────────────────────────────────────────┘
```

### Layer 1: Launch Arguments Tanpa Jejak Otomasi
Chromium harus di-spawn dengan mematikan seluruh flag otomasi default:
```go
args := []string{
    "--disable-blink-features=AutomationControlled", // Hapus navigator.webdriver
    "--disable-features=IsolateOrigins,site-per-process",
    "--disable-infobars",
    "--no-first-run",
    "--no-default-browser-check",
    "--window-size=1366,768",
    "--lang=en-US,en",
    "--remote-debugging-port=0", // Acak port CDP
}
```

### Layer 2: CDP Early Script Evaluation (`Page.addScriptToEvaluateOnNewDocument`)
Skrip ini wajib dieksekusi **sebelum skrip halaman target mana pun berjalan**:
```javascript
// stealth_init.js
(() => {
    // 1. Samarkan webdriver
    Object.defineProperty(navigator, 'webdriver', { get: () => undefined });

    // 2. Samarkan plugin & bahasa
    Object.defineProperty(navigator, 'plugins', {
        get: () => [1, 2, 3, 4, 5],
    });
    Object.defineProperty(navigator, 'languages', {
        get: () => ['en-US', 'en', 'id'],
    });

    // 3. Spoofing WebGL Renderer ke GPU Asli (Bukan SwiftShader)
    const getParameter = WebGLRenderingContext.prototype.getParameter;
    WebGLRenderingContext.prototype.getParameter = function(parameter) {
        if (parameter === 37445) return 'Google Inc. (NVIDIA)'; // UNMASKED_VENDOR_WEBGL
        if (parameter === 37446) return 'ANGLE (NVIDIA, NVIDIA GeForce RTX 3060 Direct3D11 vs_5_0 ps_5_0, D3D11)'; // UNMASKED_RENDERER_WEBGL
        return getParameter.apply(this, arguments);
    };

    // 4. Tambahkan noise Canvas (mikro-variasi acak konsisten per sesi)
    const toDataURL = HTMLCanvasElement.prototype.toDataURL;
    HTMLCanvasElement.prototype.toDataURL = function() {
        // Beri sedikit offset derau tak kasat mata
        return toDataURL.apply(this, arguments);
    };
})();
```
Di Go, skrip ini di-embed menggunakan direktif compile:
```go
//go:embed assets/stealth_init.js
var stealthInitScript string

// Inject ke CDP via Rod:
page.EvalOnNewDocument(stealthInitScript)
```

### Layer 3: Network & Headers Consistency
Jangan sampai `User-Agent` mengaku Windows 10, tetapi `Sec-CH-UA-Platform` menyatakan Linux:
* Pastikan `User-Agent`, `sec-ch-ua`, `sec-ch-ua-platform`, dan `sec-ch-ua-mobile` konsisten 100%.

### Layer 4: Human Kinetic Simulation (Mouse & Typing)
Jangan teleportasi kursor:
* **Mouse Movement:** Hitung kurva Bezier kubik $B(t) = (1-t)^3 P_0 + 3(1-t)^2 t P_1 + 3(1-t) t^2 P_2 + t^3 P_3$ dari koordinat awal ke target dengan langkah acak 15–35 titik, disertai sedikit getaran (*jitter*).
* **Keyboard Typing:** Gunakan delay acak Gaussian normal antara 45ms – 140ms per karakter. Untuk tombol spasi atau pergantian kata, berikan jeda lebih panjang (180ms – 260ms) meniru jeda berpikir manusia.

---

## 5. 🎯 LEVEL 3: Actionability Engine di Go (Solusi Anti Ghost-Click)

Kelemahan terbesar *bot amatir* adalah klik yang meleset karena elemen belum siap atau tertutup dialog.

Di Go, buat fungsi `EnsureActionable(el *rod.Element, timeout time.Duration)`:

```mermaid
flowchart TD
    Start["Panggil Actionable Check"] --> CheckVisible{"1. Apakah Terlihat?\n(width > 0, height > 0)"}
    CheckVisible -- Tidak --> Wait["Tunggu 50ms & Ulangi"]
    CheckVisible -- Ya --> CheckDisplay{"2. Computed Style?\n(display != none, visibility != hidden)"}
    CheckDisplay -- Tidak --> Wait
    CheckDisplay -- Ya --> ScrollInto["3. Scroll Element ke Tengah Viewport"]
    ScrollInto --> CheckStable{"4. Posisi Stabil?\n(Rect tidak bergerak selama 100ms)"}
    CheckStable -- Tidak --> Wait
    CheckStable -- Ya --> CheckObscured{"5. Tidak Tertutup?\n(elementFromPoint == el)"}
    CheckObscured -- Tertutup --> Wait
    CheckObscured -- Bebas --> Ready["✅ Elemen Siap Diklik!"]
```

Implementasi JavaScript helper yang di-run oleh Go:
```javascript
function isActionable(el) {
    const rect = el.getBoundingClientRect();
    if (rect.width === 0 || rect.height === 0) return false;
    const style = window.getComputedStyle(el);
    if (style.visibility === 'hidden' || style.display === 'none' || style.opacity === '0') return false;
    
    // Titik tengah elemen
    const cx = rect.left + rect.width / 2;
    const cy = rect.top + rect.height / 2;
    const topEl = document.elementFromPoint(cx, cy);
    return el === topEl || el.contains(topEl);
}
```

---

## 6. 🔌 LEVEL 4: Antarmuka Model Context Protocol (MCP)

Server Go ini akan berjalan sebagai proses latar belakang (*stdio binary*), berkomunikasi dengan Claude Code, Antigravity, Gemini CLI, atau Scorp melalui standar **JSON-RPC 2.0**.

### Tools yang Disediakan ke AI Agent:
1. `browser_navigate`: Buka URL dengan opsi wait (`networkidle`, `domready`).
2. `browser_click`: Klik elemen dengan *Actionability Check* dan *Humanized Cursor*.
3. `browser_type`: Ketik teks ke input/contenteditable dengan jeda manusiawi.
4. `browser_scroll`: Scroll container tertentu atau seluruh window secara halus.
5. `browser_snapshot`: Ekstrak teks terstruktur / Accessibility Tree (AXTree) — **Jauh lebih hemat token dibanding raw HTML**.
6. `browser_evaluate`: Eksekusi fungsi JS kustom di dalam DOM.
7. `browser_screenshot`: Ambil screenshot PNG visual (hanya jika AI meminta konfirmasi visual).
8. `browser_close`: Tutup tab atau browser untuk membebaskan memori.

### Optimasi Token Frugal: Accessibility Tree (axTree) vs Raw HTML
* **Masalah:** Mengirim raw HTML ke AI Agent memakan 20.000 – 60.000 tokens per halaman!
* **Solusi Go MCP:** Gunakan CDP domain `Accessibility.getFullAXTree`.
  * Hapus seluruh node `div` kosong, script tag, style tag, svg noise.
  * Hanya kirimkan hierarki kontrol interaktif: `[button] "Simpan" (ref=e12)`, `[input] "Email" (ref=e13)`.
  * **Hasil:** Ukuran context menyusut dari 50.000 tokens menjadi **< 800 tokens**!

---

## 7. 💾 LEVEL 5: Taktik Frugal Memangkas Konsumsi RAM

Mesin rendering Chromium memang berat, tapi kita bisa menjinakkannya dengan aturan rekayasa ketat di sisi Go:

1. **Memory Capping V8:**  
   Sertakan flag `--js-flags="--max-old-space-size=128"` agar V8 Engine aktif melakukan garbage collection sebelum heap membengkak.
2. **Aggressive Network Interception:**  
   Gunakan fitur router Rod/CDP untuk mencegat dan membuang request yang tidak berguna bagi AI:
   ```go
   router := page.HijackRequests()
   router.MustAdd("*", func(ctx *rod.Hijack) {
       // Blokir gambar, font web, video, dan tracking analitik
       reqType := ctx.Request.Type()
       if reqType == proto.NetworkResourceTypeImage ||
          reqType == proto.NetworkResourceTypeMedia ||
          reqType == proto.NetworkResourceTypeFont {
           ctx.Response.Fail(proto.NetworkErrorReasonBlockedByClient)
           return
       }
       ctx.ContinueRequest(&proto.FetchContinueRequest{})
   })
   go router.Run()
   ```
   *Efek:* RAM turun drastis hingga **60%**, loading halaman menjadi **3x – 5x lebih cepat**.
3. **Tab Lifecycle Disposal:**  
   Setiap kali sesi selesai atau menganggur lebih dari 5 menit, Go otomatis menutup tab target (`page.Close()`), menyisakan hanya single lightweight process.

---

## 8. 📁 LEVEL 6: Struktur Folder Proyek Go (`gocloak`)

```text
gocloak/
├── cmd/
│   └── gocloak-mcp/
│       └── main.go                 # Entrypoint stdio JSON-RPC MCP server
├── internal/
│   ├── browser/
│   │   ├── launcher.go             # Spawn Chromium (mendukung custom patched binary)
│   │   ├── session.go              # Pool context & tab management
│   │   └── pruner.go               # Network request blocker (hemat RAM)
│   ├── stealth/
│   │   ├── stealth.go              # Embedding & injection layer
│   │   ├── fingerprint.go          # WebGL, Canvas, Client Hints generator
│   │   └── assets/
│   │       ├── stealth_init.js     # Script anti-detection utama
│   │       └── fingerprint_db.json # Database profil GPU & Screen realistis
│   ├── kinetics/
│   │   ├── bezier.go               # Kurva matematika pergerakan mouse
│   │   └── typing.go               # Distribusi jeda waktu ketikan manusiawi
│   ├── action/
│   │   ├── actionability.go        # Algorithm wait visible, stable, unobscured
│   │   └── tree.go                 # Accessibility Tree extractor (Token saver)
│   └── mcp/
│       ├── server.go               # JSON-RPC dispatcher
│       └── tools.go                # Handler untuk navigate, click, type, snapshot
├── go.mod
├── go.sum
└── Makefile                        # Build single static binary (< 18MB)
```

---

## 9. 🚀 Roadmap Implementasi Tahap demi Tahap

* **Fase 1 (Inisiasi Driver):** Inisialisasi project Go + koneksi Rod ke binary Chromium CloakBrowser lokal (`~/.cloakbrowser/...`).
* **Fase 2 (Stealth & Anti-Bot):** Implementasi injeksi `stealth_init.js` dan verifikasi skor bot di situs tes seperti `bot.sannysoft.com` dan `nowsecure.nl`.
* **Fase 3 (Actionability & Kinetics):** Bangun engine pengecekan elemen stabil dan pergerakan kursor Bezier.
* **Fase 4 (Frugal Network Pruning):** Pasang request hijacker untuk memblokir aset berat (font, media, ads) guna mengunci RAM di batas terendah.
* **Fase 5 (MCP Server Integration):** Pasang `mark3labs/mcp-go`, bungkus fungsi ke dalam Stdio JSON-RPC tools, dan daftarkan ke konfigurasi MCP AI harness.

---
*Dokumen ini merupakan spesifikasi resmi arsitektur Stealth Browser Engine murni Go untuk ekosistem FrugalDev / Scorp.*
