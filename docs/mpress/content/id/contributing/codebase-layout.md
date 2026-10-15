---
title: "Tata Letak Basis Kode"
description: "Cara repositori Wails v3 diatur dan bagaimana setiap bagiannya saling terhubung"
slug: "contributing/codebase-layout"
sourcePath: "contributing/codebase-layout.md"
---

Wails v3 berada dalam **monorepo** yang memuat runtime framework, CLI, contoh, dokumentasi, dan rangkaian alat build.  
Halaman ini menguraikan *struktur direktori* yang penting bagi siapa pun yang ingin mendalami bagian internalnya.

## Ikhtisar Tingkat Atas

```
wails/
├── v3/               # ⬅️ Everything specific to Wails v3 lives here
├── v2/               # Legacy v2 implementation (can be ignored for v3 work)
├── docs/             # M-Press-powered v3 docs site (this page!)
├── website/          # Docusaurus v2 site and marketing pages (main site)
├── scripts/          # Misc helper scripts (e.g. sponsor image generator)
└── *.md              # Project-wide meta files (CHANGELOG, LICENSE, …)
```

Mulai dari sini, kita akan menelusuri lebih dalam struktur **`v3/`**.

## Root `v3/`

```
v3/
├── cmd/          # Compilable commands (currently only the wails3 CLI)
├── internal/     # Framework implementation (not public API)
├── pkg/          # Public Go packages — the API surface
├── tasks/        # Taskfile-based release / generation utilities
├── wep/          # RFC-style proposals (Wails Enhancement Proposals)
├── tests/        # Integration test harness
├── go.mod
└── go.sum
```

> Templat proyek disertakan di bawah `internal/templates/` (satu folder untuk setiap stack framework
>
> serta `base/`, `_common/`, dan `ios/`). Tidak ada direktori `v3/templates/`
>
> di tingkat teratas.

### Model Mental

1. **`pkg/`** mengekspos *apa yang diimpor oleh pengembang aplikasi*\
2. **`internal/`** memuat *cara mekanisme internal diimplementasikan*\
3. **`cmd/wails3`** menggerakkan *siklus hidup proyek & build*\

Semua bagian lainnya mendukung ketiga pilar tersebut.

---

## `cmd/` – Perintah

| Path | Catatan |
| --- | --- |
| `v3/cmd/wails3` | **Titik masuk CLI**. Sebuah `main.go` berukuran kecil mendelegasikan seluruh logika ke paket-paket dalam `internal/commands`. |
| `internal/commands/*` | Subperintah (init, dev, build, doctor, …). Masing-masing berada dalam file tersendiri agar mudah ditemukan. |
| `internal/commands/task_wrapper.go` | Menjembatani flag CLI dengan pipeline build Taskfile. |

CLI menangani:

- **Pembuatan kerangka proyek** (`init`, pembuatan templat)\
- **Orkestrasi server pengembangan** (`dev`, pemuatan ulang langsung)\
- **Build produksi & pemaketan** (`build`, `package`, pembungkus platform)\
- **Diagnostik** (`doctor`)\

---

## `internal/` – Ruang Mesin

```
internal/
├── assetserver/  # Serving & embedding web assets
├── buildinfo/    # Reproducible build metadata
├── commands/     # CLI mechanics (see above)
├── runtime/      # Build-tag glue + embedded JS runtime sources
├── generator/    # Static analysis & binding generator
├── templates/    # Project templates (frontend stacks)
├── packager/     # nfpm wrapper used by `wails3 tool package`
├── capabilities/ # Host OS capability probing
├── dbus/         # Generic D-Bus helper
├── service/      # Service-template scaffolding (`wails3 generate service`)
└── ...           # [other helper sub-packages: flags, hash, term, …]
```

### Subpaket Utama

| Paket | Tanggung Jawab | Tempat Terhubungnya |
| --- | --- | --- |
| `runtime` | Menampung kode penghubung kecil dengan tag build `runtime{,_darwin,_linux,_windows,_android,_dev,_prod}.go` serta runtime JS tertanam di bawah `runtime/desktop/`. Kode aktual per OS untuk jendela, papan klip, dialog, dan baki sistem berada di `pkg/application/*_{darwin,linux,windows}.go`. | Diimpor secara tidak langsung melalui `pkg/application`. |
| `assetserver` | Server file bermode ganda:<br />• Pengembangan: menyajikan dari disk & mem-proxy Vite (`build_dev.go`)<br />• Produksi: menyematkan aset melalui `go:embed` (`build_production.go`) | Diinisialisasi oleh `pkg/application` saat proses awal. |
| `generator` | Mengurai kode sumber Go untuk membuat **metadata binding** yang kemudian menghasilkan file stub TypeScript/JS dan konstanta peristiwa. Titik masuk: `generator.Generate` / `generator.Generator` di atas `collect/` + `render/`. | Dipicu oleh `wails3 generate bindings`. |
| `packager` | Pembungkus `nfpm` yang digunakan untuk menghasilkan artefak Linux `deb`/`rpm`/`archlinux` (digerakkan oleh konfigurasi nfpm `myapp.DEB`/`.RPM`/`.ARCHLINUX` di bawah `internal/commands/`). | Dipanggil oleh `wails3 tool package`. DMG macOS / MSIX Windows berada di bawah `internal/commands/{dmg,msix.go,webview2/}`. |

Utilitas pendukung (misalnya `s/`, `hash/`, `flags/`) menjaga agar kepentingan internal tetap terpisah.

---

## `pkg/` – API Publik

```
pkg/
├── application/  # Core API: App, windows, menus, dialogs, events, managers
├── events/       # Event constants (Common/Mac/Windows/Linux) + generator
├── services/     # Optional built-in services (notifications, kvstore, …)
├── doctor-ng/    # New-style `wails3 doctor-ng` checks
├── errs/         # Shared error types
├── icons/        # Default platform icons
├── mac/          # macOS-only helpers
└── w32/          # Windows Win32 helpers
```

> Tidak ada paket `pkg/runtime/`, `pkg/options/`, atau `pkg/menu/`. Opsi jendela/menu
>
> berada bersama `pkg/application` (misalnya `WebviewWindowOptions`, `Menu`,
>
> `MenuItem`), dan `assetserver/` berada di bawah `internal/`.

`pkg/application` melakukan bootstrap pada program Wails:

```go
func main() {
    app := application.New(application.Options{
        Name: "MyApp",
        Assets: application.AssetOptions{
            Handler: application.AssetFileServerFS(assetsFS),
        },
    })
    window := app.Window.New()
    window.SetTitle("Hello").SetSize(1024, 768)
    _ = app.Run()
}
```

Di balik layar, komponen ini:

1. Menghubungkan kode penghubung tag build di `internal/runtime` dan kode khusus setiap OS di `pkg/application/`
2. Menyiapkan instans `internal/assetserver`
3. Mendaftarkan semua pemroses pesan berbasis binding
4. Memasuki thread utama OS

---

## `internal/templates/` – Cetak Biru Kerangka Proyek

`internal/templates/` menyertakan **templat dasar** (tata letak Go di bawah `base/`,  
`_common/`, `ios/`) dan **tema frontend** (`vanilla[-ts]`, `react[-ts]`,  
`react-swc[-ts]`, `lit[-ts]`, `preact[-ts]`, `qwik[-ts]`, `solid[-ts]`,  
`svelte[-ts]`, `sveltekit[-ts]`, `vue[-ts]`).

Pada `wails3 init -t react`, CLI:

1. Menyalin file Go `_common`
2. Menggabungkan paket frontend yang diinginkan
3. Menjalankan `go mod tidy` (dapat dilewati dengan `--skipgomodtidy`)

Mengedit templat **tidak** memengaruhi aplikasi yang sudah ada, hanya `init` berikutnya. Contoh publik berada di bawah `v3/examples/`; contoh tersebut bukan pengganti rangkaian pengujian otomatis yang dijelaskan dalam dokumentasi kontributor.

---

## `tasks/` – Otomatisasi Rilis

Taskfile membungkus kompilasi silang yang kompleks, pembaruan versi, dan pembuatan changelog. Taskfile digunakan secara terprogram oleh `internal/commands/task.go` sehingga logika yang sama mendukung **CLI** dan **CI**.

---

## Cara Komponen Saling Berinteraksi

```d2
direction: down
CLI: CLI wails3
Generator: internal/generator
AssetDev: assetserver (pengembangan)
Packager: internal/packager
AppRuntime: {
  label: Runtime aplikasi
  ApplicationPkg: pkg.application
  InternalRuntime: internal.runtime
  OSAPIs: API OS
}
CLI -> Generator: build / buat
CLI -> AssetDev: pengembangan
CLI -> Packager: kemas
Generator -> ApplicationPkg: binding
ApplicationPkg -> InternalRuntime
InternalRuntime -> OSAPIs
ApplicationPkg -> AssetDev
ApplicationPkg.label: ApplicationPkg
InternalRuntime.label: InternalRuntime
OSAPIs.label: OSAPIs
```

*CLI → generator → runtime* membentuk jalur inti dari **kode sumber** hingga **aplikasi desktop yang berjalan**.

---

## Kiat Orientasi

| Perlu memahami… | Lihat… |
| --- | --- |
| Shim platform | `pkg/application/*_darwin.go`, `*_linux.go`, `*_windows.go` (window, clipboard, dialogs, systray, mainthread, events_common). cgo Linux: `pkg/application/linux_cgo*.go`. |
| Protokol bridge | `pkg/application/messageprocessor*.go` |
| Alur kerja aset | `internal/assetserver/` (`build_dev.go` dibandingkan dengan `build_production.go`) |
| Alur pemaketan | `internal/commands/{appimage,msix,dot_desktop,dmg/}.go`, `internal/packager/` |
| Mesin templat | `internal/templates/` (`templates.Install`, `templates.GetDefaultTemplates`) |
| Analisis statis | `internal/generator/{generate.go,collect/,render/}` |

---

Kini Anda memiliki **peta mental** repositori ini. Gunakan bersama `ripgrep`, fitur “Go to File/Symbol” di IDE Anda, dan aplikasi contoh untuk menelusuri fitur apa pun lebih dalam. Selamat mengutak-atik!
