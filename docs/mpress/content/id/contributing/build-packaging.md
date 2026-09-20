---
title: "Pipeline Build \u0026 Pemaketan"
description: "Proses internal yang terjadi saat Anda menjalankan `wails3 build`, cara menghasilkan biner lintas platform, dan cara membuat penginstal untuk setiap OS."
slug: "contributing/build-packaging"
sourcePath: "contributing/build-packaging.md"
---

`wails3 build` sengaja dibuat **ringkas**: ini adalah pembungkus Taskfile yang meneruskan tag build tambahan ke tugas `build` milik proyek host. Pekerjaan utamanya berada di `build/Taskfile.yml` milik proyek itu sendiri (yang dihasilkan oleh `wails3 init`), di `internal/commands/build-assets.go` (yang mengelola aset pada waktu bake), serta di `internal/packager` (pemaketan nfpm Linux) dan `internal/commands/appimage.go`, `internal/commands/msix.go`, `internal/commands/dmg/dmg.go`, dan `internal/commands/dot_desktop.go` (penginstal per platform).

Halaman ini membahas:

1. Titik masuk CLI yang sebenarnya
2. Alur build berbasis Taskfile
3. Bake aset & injeksi informasi build
4. Back-end pemaketan per platform
5. Penyesuaian pipeline
6. Pemecahan masalah

---

## 1. Titik Masuk CLI yang Sebenarnya

```
wails3 build       → internal/commands.Build       (in task_wrapper.go)
wails3 package     → internal/commands.Package     (in task_wrapper.go)
wails3 generate build-assets → GenerateBuildAssets (in build-assets.go)
wails3 update build-assets   → UpdateBuildAssets   (in build-assets.go)
wails3 tool buildinfo        → BuildInfoOptions    (in tool_buildinfo.go)
wails3 tool package          → internal/packager   (nfpm wrapper)
wails3 generate .desktop     → in dot_desktop.go
```

`internal/commands/task_wrapper.go`:

```go
func Build(buildFlags *flags.Build, otherArgs []string) error {
    // forwards --tags / EXTRA_TAGS, then defers to a Taskfile target
    return wrapTask("build", otherArgs)
}
```

`flags.Build` hanya menyediakan **satu** flag — `--tags` (diteruskan sebagai `EXTRA_TAGS=`). `wails3 build` **tidak** memiliki flag `-platform`, `-o`, `-skipbindings`, `-skip-package`, `-package`, `-ldflags`, `-verbose`, `-debug`, `-devbuild`, `-icon`, atau `-clean`. Kompilasi silang, jalur keluaran, ikon, dan sebagainya dikonfigurasi di **`Taskfile.yml`**, **`build/config.yml`**, serta perintah pendukung `wails3 generate icons` / `wails3 generate build-assets`.

`build/build.json` **bukan** bagian dari v3 — konfigurasinya menggunakan `Taskfile.yml` serta `build/config.yml`.

---

## 2. Alur Build Berbasis Taskfile

Proyek yang baru diinisialisasi menyertakan `build/Taskfile.yml` dengan namespace yang kurang lebih sebagai berikut:

| Namespace | Tugas (pilihan) |
| --- | --- |
| `darwin:` | `build`, `build:universal`, `package`, `run`, `dev` |
| `windows:` | `build`, `package`, `run`, `dev` |
| `linux:` | `build`, `package`, `run`, `dev` |
| `common:` | `update:build-assets`, `generate:icons`, `generate:syso` |

Secara default, `wails3 build` menjalankan namespace `build` milik OS host; kemudian Taskfile proyek menjalankan `go build` melalui shell dengan flag khusus host. Untuk membangun bagi OS lain, jalankan tugas OS tersebut secara langsung (misalnya `wails3 task darwin:build:universal`), bukan dengan meneruskan flag ke `wails3 build`.

Direktori keluaran default adalah **`bin/<APP_NAME>`** (tanpa prefiks `build/bin/`).

---

## 3. Aset pada Waktu Bake dan Informasi Build

| Aspek | File |
| --- | --- |
| Pembuatan/pembaruan aset build | `internal/commands/build-assets.go` |
| Pencetak informasi build (CLI: `wails3 tool buildinfo`) | `internal/commands/tool_buildinfo.go` — mencetak informasi; ini **bukan** injektor `ldflags` |
| Stub produksi | `internal/assetserver/build_production.go` — `//go:build production` |
| Bundel frontend | disematkan melalui `//go:embed` dalam package milik aplikasi itu sendiri (misalnya di sebelah `main.go`) |
| Sumber daya Windows (`.syso`) | `internal/commands/syso.go` — menghasilkan `rsrc_windows_<arch>.syso` |
| MSIX Windows | `internal/commands/msix.go` + `internal/commands/webview2/` |
| Input DMG macOS | `internal/commands/dmg/` |
| `.desktop` Linux | `internal/commands/dot_desktop.go` |

CLI tidak melakukan bake otomatis terhadap `bundled_assetserver.go` untuk aplikasi Anda — `internal/assetserver/bundled_assetserver.go` **ditulis secara manual** dan membungkus runtime JS yang disematkan di bawah `bundledassets/`.

---

## 4. Back-End Pemaketan

### Linux

Pemaketan Linux ditangani oleh **nfpm** (bukan `fpm`):

- `internal/packager/packager.go` membungkus `github.com/goreleaser/nfpm/v2` dan menyediakan `CreatePackageFromConfig(pkgType, configPath, output)` / `CreatePackageFromConfigWriter(...)`.
- Proyek yang dihasilkan menyertakan konfigurasi bergaya nfpm `myapp.DEB`, `myapp.RPM`, dan `myapp.ARCHLINUX` di bawah `internal/commands/` (digunakan oleh `wails3 tool package`).
- Pembuatan AppImage berada di `internal/commands/appimage.go`, yang memanggil `linuxdeploy` + `linuxdeploy-plugin-gtk` (plugin disertakan di `internal/commands/linuxdeploy-plugin-gtk.sh`).

`wails3 build` **tidak** memiliki flag `-package deb`/`rpm`. Gunakan `wails3 tool package` atau target Taskfile khusus platform.

### macOS

- `darwin:package` dari Taskfile proyek menghasilkan bundel `.app`.
- Aset DMG berada di bawah `internal/commands/dmg/`; proyek dapat membungkus bundel menjadi DMG dengan `hdiutil` setelah `darwin:package` selesai (Taskfile dalam templat yang lebih baru menyertakan pembantu `dmg`).
- Pengidentifikasi CFBundle, versi, dan hak cipta berasal dari flag `-product*` selama `wails3 init` serta dari `build/config.yml`.

### Windows

- Paket Windows menargetkan **MSIX** (bukan WiX/MSI). Lihat `internal/commands/msix.go` dan `internal/commands/webview2/` untuk alur kerja lengkapnya.
- **Tidak ada** `internal/commands/packager.go` dan **tidak ada** direktori `internal/commands/windows_resources/`.
- Penandatanganan kode opsional untuk berkas yang dapat dieksekusi dijalankan melalui `wails3 tool sign` (Authenticode) — lihat `internal/commands/sign.go`.

---

## 5. Menyesuaikan Pipeline

| Kebutuhan | Pendekatan |
| --- | --- |
| Tag build tambahan | `wails3 build --tags myFeature,otherTag` |
| Linter/langkah prabuild | Tambahkan tugas ke `build/Taskfile.yml` dan jadikan tugas `build` khusus OS bergantung padanya |
| Kompilasi silang | Jalankan tugas OS yang relevan (misalnya `wails3 task linux:build`) — tidak ada flag `-platform` |
| Lewati pemaketan | Cukup jalankan tugas `build`; `package` merupakan tugas terpisah |
| Pembuat paket khusus | Letakkan konfigurasi di bawah `internal/commands/myapp.*` dan panggil `wails3 tool package` dengan `-config <file>` |
| Hapus simbol | Edit tugas `darwin:/windows:/linux:` `build` agar meneruskan `-ldflags "-s -w"` secara langsung ke `go build` — `wails3 build` sendiri tidak memiliki flag `-ldflags` |

Semua target Taskfile mematuhi variabel lingkungan yang dipublikasikan Wails (`APP_NAME`, `WAILS_VITE_PORT`, `FRONTEND_DEVSERVER_URL`, …), sehingga tugas khusus dapat mengandalkannya.

---

## 6. Pemecahan Masalah

| Gejala | Kemungkinan Penyebab | Perbaikan |
| --- | --- | --- |
| **`ld: framework not found WebKit` (mac)** | Alat CLI Xcode tidak tersedia | `xcode-select --install` |
| **Jendela kosong dalam build produksi** | Build frontend gagal atau masalah perutean SPA | Pastikan `frontend/dist/index.html` tersedia dan penangan aset Anda beralih menggunakannya sebagai fallback |
| **Alat pemaketan MSIX tidak tersedia** | SDK `WebView2`/alat MSIX belum diinstal | Jalankan `wails3 task install:msix:tools` |
| **`linuxdeploy` tidak ditemukan** | Plugin tidak tersedia di PATH | Instal `linuxdeploy` + jalankan `internal/commands/linuxdeploy-plugin-gtk.sh` melalui langkah instalasi otomatis CLI |

`wails3 build` tidak memiliki flag `-verbose`. Atur `TASK_X_VERBOSE=1` (Taskfile) atau periksa target tugas secara langsung untuk melihat perintah yang dijalankan.

---

## 7. Peta Sumber Utama

| Aspek | Berkas |
| --- | --- |
| Pembungkus build | `internal/commands/task_wrapper.go` (`Build`, `Package`, `SignWrapper`, `wrapTask`) |
| Pembuatan aset build | `internal/commands/build-assets.go` (`GenerateBuildAssets`, `UpdateBuildAssets`) |
| Pencetak informasi build | `internal/commands/tool_buildinfo.go` |
| Pembuat AppImage | `internal/commands/appimage.go` |
| Pemaketan Linux (nfpm) | `internal/packager/packager.go`, `internal/commands/myapp.{DEB,RPM,ARCHLINUX}` |
| MSIX Windows | `internal/commands/msix.go`, `internal/commands/webview2/` |
| Generator sumber daya Windows | `internal/commands/syso.go` |
| Aset DMG macOS | `internal/commands/dmg/` |
| Generator `.desktop` | `internal/commands/dot_desktop.go` |
| Konstanta versi | `internal/version/version.go` |

Simpan tabel ini agar mudah dirujuk saat Anda melacak kegagalan build.

---

Sekarang Anda memiliki gambaran lengkap, mulai dari **kode sumber** hingga **penginstal**. Singkatnya, `wails3 build` sendiri hanyalah pembungkus tipis — hampir semua penyesuaian dilakukan dalam `Taskfile.yml`/`build/config.yml` proyek, atau melalui subperintah `wails3 generate …`/`wails3 tool …` yang eksplisit. Selamat merilis!
