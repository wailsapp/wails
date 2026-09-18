---
title: "Server Aset"
description: "Cara Wails v3 menyajikan dan menyematkan aset web Anda dalam pengembangan dan produksi"
slug: "contributing/asset-server"
sourcePath: "contributing/asset-server.md"
---

## Ikhtisar

Setiap aplikasi Wails didistribusikan sebagai **satu executable native** yang menggabungkan:

1. Backend *Go* Anda
2. Frontend *Web* (HTML + JS + CSS)

**Server Aset** adalah penghubung yang memungkinkan hal ini. Server ini memiliki **dua mode operasi** yang dipilih pada waktu kompilasi melalui tag build Go:

| Mode | Tag | Tujuan |
| --- | --- | --- |
| **Pengembangan** | `//go:build !production` | Iterasi cepat dengan hot reload |
| **Produksi** | `//go:build production` | Aset tersemat tanpa dependensi |

Implementasinya berada di `v3/internal/assetserver/` dengan pemisahan file yang jelas:

```
build_dev.go              # ⬅️ dev-only entrypoint (!production build tag)
build_production.go       # ⬅️ production-only entrypoint (production build tag)
assetserver.go            # Shared core
assetserver_dev.go        # Dev proxy/disk handler
assetserver_webview.go    # WebView-side adapter
assetserver_darwin.go     # OS-specific helpers (also linux/windows variants)
asset_fileserver.go       # Shared static file logic
content_type_sniffer.go   # MIME type detection
mimecache.go              # Cached MIME lookups
ringqueue.go              # Tiny in-memory LRU
options.go                # Configuration struct
middleware.go             # http.Handler middleware type
bundled_assetserver.go    # Hand-written wrapper around embedded bundles
bundledassets/            # Embedded runtime JS assets
```

---

## Mode Pengembangan

### Siklus Hidup

1. `wails3 dev` melakukan boot dan **menjalankan server pengembangan frontend Anda** (Vite, SvelteKit, React-SWC …) dengan menjalankan task yang ditentukan di `build/Taskfile.yml` (biasanya `npm run dev`).
2. CLI menetapkan `WAILS_VITE_PORT` ke port pengembangan Wails dan `FRONTEND_DEVSERVER_URL` ke URL **lengkap** (`http://host:port` / `https://host:port`) yang mengarah ke server pengembangan framework yang sedang berjalan. Lihat `internal/commands/dev.go`.
3. Server aset pengembangan (dikompilasi melalui `//go:build !production` di `build_dev.go`) membaca `FRONTEND_DEVSERVER_URL` melalui `GetDevServerURL()` dan meneruskan traffic non-runtime kepadanya melalui reverse proxy.
4. File statis (`/assets/logo.svg`) dapat **disajikan langsung dari disk** melalui `asset_fileserver.go` (agar cepat), sedangkan semua yang tidak dikenal akan **diteruskan melalui proxy** ke server pengembangan framework sehingga Anda memperoleh hot module replacement secara *instan*.

```
┌─────────┐  /wails/runtime.js     ┌─────────────┐
│ Browser │ ── embedded runtime ──▶│   Runtime   │
├─────────┤                        └─────────────┘
│   JS    │  / (index.html)        proxy / -> Vite via FRONTEND_DEVSERVER_URL
└─────────┘ ◀─────────────┐
              AssetServer │
                          ▼
                   ┌────────────┐
                   │  Vite Dev  │
                   │   Server   │
                   └────────────┘
```

### Fitur

- **Pemuatan Ulang Langsung** — Vite/SvelteKit/dll. menyuntikkan HMR melalui WebSocket; server aset pengembangan meneruskannya secara transparan melalui proxy.
- **Dukungan Source Map** — karena aset tidak dibundel, devtools browser Anda memetakan kembali error ke sumber asli.
- **Tanpa Kompilasi Ulang Go** — Hanya frontend yang dibangun ulang; kode Go tetap berjalan hingga Anda mengubah file `.go`.

### Beralih Framework

Proxy pengembangan **tidak bergantung pada framework**. CLI Wails menyediakan dua variabel lingkungan saat menjalankan task pengembangan Anda:

| Variabel lingkungan | Sumber | Arti |
| --- | --- | --- |
| `WAILS_VITE_PORT` | `internal/commands/dev.go` (konstanta `wailsVitePort`) | Port pengembangan default (9245 kecuali `--port` diteruskan) — konfigurasi Vite Anda sebaiknya mematuhinya |
| `FRONTEND_DEVSERVER_URL` | `internal/commands/dev.go` | URL lengkap yang akan menjadi tujuan proxy Wails; dibaca di Go melalui `assetserver.GetDevServerURL()` (`build_dev.go`) |

Tidak ada variabel lingkungan `VITE_PORT`, `FRONTEND_DEV_PORT`, atau `WAILSDEV_VERBOSE` di tree v3.

Tambahkan template baru → tentukan task pengembangannya → server aset langsung berfungsi.

---

## Mode Produksi

Saat Anda menjalankan `wails3 build`, pipeline akan:

1. Menjalankan **build produksi** frontend (`npm run build`) yang menghasilkan `frontend/dist/**`.
2. **Menyematkan** direktori tersebut ke dalam aplikasi melalui `go:embed` di package aplikasi Anda sendiri (biasanya `//go:embed all:frontend/dist` di samping `main.go`).
3. Mengompilasi binary Go dengan `-tags production` (diteruskan melalui `EXTRA_TAGS` oleh wrapper Taskfile).

`internal/assetserver/build_production.go` adalah stub tag build yang menggantikan jalur kode dengan jalur produksi. `internal/assetserver/bundled_assetserver.go` **ditulis secara manual** — file ini membungkus JS runtime dalam `bundledassets/` dan bukan file yang dihasilkan secara otomatis.

### Penanganan Permintaan

Handler yang sebenarnya adalah `internal/assetserver/assetserver.go` / `asset_fileserver.go`. Secara konseptual:

1. Coba akses aset statis tersemat di path yang diminta.
2. Jika gagal, gunakan `index.html` untuk routing SPA.
3. Deteksi jenis konten jika ekstensinya tidak diketahui (`content_type_sniffer.go`).
4. Tetapkan header cache yang wajar.

- **Deteksi MIME** — jenis konten file tanpa ekstensi dideteksi dari sekitar 512 byte pertama (`content_type_sniffer.go`), lalu hasilnya disimpan dalam cache di `mimecache.go` / `ringqueue.go`.
- **Header Keamanan** — melarang navigasi `file://` dan menetapkan `nosniff`.

Karena semuanya disematkan, binary yang didistribusikan **tidak memiliki dependensi eksternal** (bahkan di Windows).

---

## Menjembatani Pengembangan ↔ Produksi

Dari sudut pandang `pkg/application`, kedua mode mengekspos **antarmuka publik yang sama**: struct `AssetOptions` dengan `Handler http.Handler`, ditambah middleware dan pengawatan siklus hidup di dalam `internal/assetserver/`. Peralihan antara pengembangan dan produksi sepenuhnya dilakukan melalui tag build Go, sehingga kode aplikasi identik pada kedua mode.

```go
app := application.New(application.Options{
    Assets: application.AssetOptions{
        Handler: application.AssetFileServerFS(assetsFS),
    },
})
```

---

## Cara Framework Frontend Berintegrasi

### Template

Setiap template yang disertakan (React, Vue, Svelte, Solid, Vanilla, …) berisi:

- `build/Taskfile.yml`
- `frontend/vite.config.ts` (atau yang setara)

Konfigurasi Vite atau yang setara membaca `WAILS_VITE_PORT` dan mengikat server pengembangan ke port tersebut. CLI kemudian memublikasikan `FRONTEND_DEVSERVER_URL` aktual untuk digunakan oleh proksi dalam aplikasi.

Framework tetap sepenuhnya terpisah dari Go:

- Tidak perlu mengimpor Wails JS SDK apa pun pada waktu build — `/wails/runtime.js` disajikan oleh server aset saat runtime.
- Framework apa pun yang memiliki server pengembangan HTTP dapat diintegrasikan.

---

## Memperluas / Menyesuaikan

Memerlukan header khusus, autentikasi, atau gzip?

1. Definisikan `middleware.Middleware` (alias untuk `func(http.Handler) http.Handler`, yang dideklarasikan dalam `internal/assetserver/middleware.go`).
2. Hubungkan ke `application.AssetOptions` Anda melalui konfigurasi yang diekspos oleh `internal/assetserver/options.go`.
3. Perilakunya identik dalam mode pengembangan dan produksi — tidak ada daftar middleware terpisah untuk setiap mode.

---

## File Sumber Utama

| File | Peran |
| --- | --- |
| `build_dev.go` / `build_production.go` | Pembungkus tag build yang memilih mode pengembangan atau produksi |
| `assetserver.go` / `asset_fileserver.go` | Handler HTTP inti |
| `assetserver_dev.go` | Proksi balik ke `FRONTEND_DEVSERVER_URL` |
| `bundled_assetserver.go` | Pembungkus yang ditulis secara manual di sekitar `bundledassets/` |
| `options.go` | Konfigurasi yang ditujukan untuk `application.AssetOptions` |
| `mimecache.go` / `ringqueue.go` | Cache MIME + LRU kecil |

---

## Hal yang Perlu Diwaspadai & Debugging

- **Layar Putih dalam Produksi** — biasanya disebabkan oleh perutean SPA: pastikan server pengembangan Anda menyajikan `index.html` untuk jalur yang tidak dikenal dan fallback handler produksi tersemat tercapai.
- **404 dalam Pengembangan** — konfigurasi Vite Anda tidak mengikat ke `WAILS_VITE_PORT`, atau CLI tidak dapat menjangkau server pengembangan untuk mengisi `FRONTEND_DEVSERVER_URL`.
- **Aset Besar** — penyematan memperbesar ukuran biner. Sajikan media berukuran besar dari origin terpisah atau alirkan melalui `http.Handler` khusus.

---

Sekarang Anda mengetahui cara **Server Aset** Wails menyalurkan kode web Anda ke jendela native, baik dalam mode **pengembangan** maupun **produksi**. Kuasai lapisan ini agar Anda dapat men-debug masalah pemuatan, menambahkan middleware, atau bahkan mengganti seluruh toolchain frontend dengan percaya diri.
