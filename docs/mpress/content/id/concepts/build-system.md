---
title: "Sistem Build"
description: "Memahami cara Wails membangun dan mengemas aplikasi Anda"
slug: "concepts/build-system"
sourcePath: "concepts/build-system.md"
---

## Sistem Build Terpadu

Wails menyediakan **sistem build terpadu** yang mengompilasi kode Go, membundel aset frontend, menyematkan semuanya ke dalam satu berkas executable, dan menangani build khusus platform—semuanya dengan satu perintah.

```bash
wails3 build
```

**Keluaran:** Berkas executable native dengan semua komponen yang telah disematkan.

## Ringkasan Proses Build

**[Placeholder Diagram Proses Build]**

## Tahapan Build

### 1. Tahap Analisis

Wails memindai kode Go Anda untuk memahami layanan Anda:

```go
type GreetService struct {
    prefix string
}

func (g *GreetService) Greet(name string) string {
    return g.prefix + name + "!"
}
```

**Yang diekstrak Wails:**

- Nama layanan: `GreetService`
- Nama metode: `Greet`
- Tipe parameter: `string`
- Tipe nilai kembalian: `string`

**Digunakan untuk:** Menghasilkan binding TypeScript

### 2. Tahap Pembuatan

#### Binding TypeScript

Wails menghasilkan binding dengan keamanan tipe:

```javascript
// Auto-generated: frontend/bindings/<full-go-import-path>/greetservice.js
// (TypeScript is also generated when you pass `-ts`. The shape below is the real
// runtime call format — numeric IDs via $Call.ByID, imported from /wails/runtime.js.)
import { Call as $Call } from "/wails/runtime.js";

export function Greet($0) {
    return $Call.ByID(1234567890, $0);
}
```

**Manfaat:**

- Keamanan tipe penuh
- Pelengkapan otomatis di IDE
- Kesalahan pada waktu kompilasi
- Komentar JSDoc

#### Build Frontend

Bundler frontend Anda dijalankan (Vite, webpack, dan sebagainya):

```bash
# Vite example
vite build --outDir dist
```

**Yang terjadi:**

- JavaScript/TypeScript dikompilasi
- CSS diproses dan diminifikasi
- Aset dioptimalkan
- Source map dihasilkan (khusus pengembangan)
- Keluaran disimpan ke `frontend/dist/`

### 3. Tahap Kompilasi

#### Kompilasi Go

Kode Go dikompilasi dengan pengoptimalan:

```bash
go build -ldflags="-s -w" -o myapp.exe
```

**Flag:**

- `-s`: Menghapus tabel simbol
- `-w`: Menghapus informasi debug DWARF
- Hasil: Berkas biner lebih kecil (berkurang ~30%)

**Khusus platform:**

- Windows: `.exe` dengan ikon yang disematkan
- macOS: Struktur bundle `.app`
- Linux: Berkas biner ELF

#### Penyematan Aset

Aset frontend disematkan ke dalam berkas biner Go:

```go
//go:embed frontend/dist
var assets embed.FS
```

**Hasil:** Satu berkas executable yang memuat semuanya.

### 4. Keluaran

**Satu berkas biner native:**

- Windows: `myapp.exe` (~15MB)
- macOS: `myapp.app` (~15MB)
- Linux: `myapp` (~15MB)

**Tanpa dependensi** (kecuali WebView sistem).

## Pengembangan vs Produksi

@tabs{sync-key="mode"}
[Pengembangan (wails3 dev)]
**Dioptimalkan untuk kecepatan:**

```bash
wails3 dev
```

**Yang terjadi:**

1. Memulai server pengembangan frontend (Vite menggunakan porta 9245 secara default)
2. Mengompilasi Go tanpa pengoptimalan
3. Meluncurkan aplikasi yang mengarah ke server pengembangan
4. Mengaktifkan hot reload
5. Menyertakan source map

**Karakteristik:**

- **Build ulang cepat** (&lt;1 dtk untuk perubahan frontend)
- **Tanpa penyematan aset** (disajikan dari server pengembangan)
- **Simbol debug** disertakan
- **Source map** diaktifkan
- **Pencatatan log mendetail**

**Ukuran file:** Lebih besar (~50 MB dengan simbol debug)

[Produksi (wails3 build)]
**Dioptimalkan untuk ukuran dan performa:**

```bash
wails3 build
```

**Yang terjadi:**

1. Membangun frontend untuk produksi (diminifikasi)
2. Mengompilasi Go dengan optimisasi
3. Menghapus simbol debug
4. Menyematkan aset
5. Membuat satu berkas biner

**Karakteristik:**

- **Kode yang dioptimalkan** (diminifikasi, menerapkan tree-shaking)
- **Aset disematkan** (tanpa file eksternal)
- **Simbol debug dihapus**
- **Tanpa peta sumber**
- **Pencatatan log minimal**

**Ukuran file:** Lebih kecil (~15 MB)

@end

## Perintah Build

### Build Dasar

```bash
wails3 build
```

**Keluaran:** `bin/<APP_NAME>` (atau `bin/<APP_NAME>.exe` di Windows). Direktori `bin/` berada di root proyek.

`wails3 build` adalah pembungkus tipis untuk `wails3 task build`. Satu-satunya flag waktu build yang diteruskannya adalah `--tags`, yang menjadi variabel Taskfile `EXTRA_TAGS`:

```bash
# Build with extra Go build tags
wails3 build --tags "myfeature,gtk4"
```

`wails3 build` tidak memiliki flag `-platform`, `-o`, `-skipbindings`, `-clean`, `-debug`, `-devbuild`, `-icon`, `-ldflags`, atau `-package`. Kompilasi silang, jalur keluaran, ikon, dan pengemasan dikendalikan melalui Taskfile proyek (`Taskfile.yml` + `build/config.yml`).

### Build lintas platform dan khusus platform

Build platform tersedia sebagai tugas Taskfile dalam namespace `darwin:` / `windows:` / `linux:` (ditentukan dalam `build/Taskfile.<platform>.yml`). Contoh:

```bash
# macOS — universal binary
wails3 task darwin:build:universal

# macOS — current arch
wails3 task darwin:build

# Windows
wails3 task windows:build

# Linux
wails3 task linux:build
```

Untuk melihat semua tugas yang tersedia dalam proyek saat ini:

```bash
wails3 task --list
```

### Ikon dan pengemasan

Buat ikon platform (`build/icons.icns`, `build/icon.ico`, dan sebagainya) dari PNG sumber:

```bash
wails3 generate icons -input appicon.png
```

Buat penginstal/paket khusus platform:

```bash
wails3 package           # uses the current Go build env
wails3 task linux:create:deb
wails3 task windows:package
wails3 task darwin:package:universal
```

## Konfigurasi Build

### Taskfile.yml

Proyek Wails 3 menggunakan [Taskfile](https://taskfile.dev/) sebagai orkestrator build. `Taskfile.yml` di root menyertakan file tugas per platform dari `build/`:

```yaml
# Taskfile.yml (excerpt — the real templates are richer)
version: '3'

includes:
  common: ./build/Taskfile.yml
  darwin: ./build/Taskfile.darwin.yml
  windows: ./build/Taskfile.windows.yml
  linux: ./build/Taskfile.linux.yml

tasks:
  build:
    desc: Build the application
    cmds:
      - task: "{{OS}}:build"
```

Jalankan tugas dengan `wails3 task <name>` atau `task <name>`:

```bash
wails3 task windows:build
wails3 task darwin:package:universal
wails3 task linux:create:appimage
```

### Konfigurasi proyek: `build/config.yml`

Metadata proyek (nama, pengidentifikasi, versi, nilai info-plist, pengaturan NSIS, kolom `.desktop`, protokol kustom, dan sebagainya) berada dalam `build/config.yml`. Taskfile membaca file ini saat membuat ikon, manifes, penginstal, dan sejenisnya. Wails 3 **tidak** memiliki file `build/build.json`.

```yaml
# build/config.yml (illustrative)
info:
  productName: "My App"
  productIdentifier: "com.example.myapp"
  productVersion: "1.0.0"
  companyName: "Example Ltd."
  productDescription: "An application built with Wails"
```

Jalankan `wails3 generate build-assets` (atau `wails3 update build-assets`) untuk memperbarui aset build khusus platform dari konfigurasi ini.

## Penyematan Aset

### Cara Kerjanya

Wails menggunakan paket `embed` milik Go:

```go
package main

import (
    "embed"
    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
    app := application.New(application.Options{
        Name:   "My App",
        Assets: application.AssetOptions{
            Handler: application.AssetFileServerFS(assets),
        },
    })
    
    app.Window.New()
    app.Run()
}
```

**Pada waktu build:**

1. Frontend dibangun ke `frontend/dist/`
2. Direktif `//go:embed` menyertakan file
3. File dikompilasi ke dalam berkas biner
4. Berkas biner memuat semuanya

**Saat runtime:**

1. Aplikasi dimulai
2. Aset disajikan dari memori
3. Tidak ada I/O disk untuk aset
4. Pemuatan cepat

### Aset Kustom

Sematkan file tambahan:

```go
//go:embed frontend/dist
var frontendAssets embed.FS

//go:embed data/*.json
var dataAssets embed.FS

//go:embed templates/*.html
var templateAssets embed.FS
```

## Optimisasi Build

### Optimisasi Frontend

**Vite (bawaan):**

```javascript
// vite.config.js
export default {
  build: {
    minify: 'terser',
    terserOptions: {
      compress: {
        drop_console: true,  // Remove console.log
        drop_debugger: true,
      },
    },
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ['react', 'react-dom'],  // Separate vendor bundle
        },
      },
    },
  },
}
```

**Hasil:**

- JavaScript diminifikasi (pengurangan ~70%)
- CSS diminifikasi (pengurangan ~60%)
- Gambar dioptimalkan
- Tree-shaking diterapkan

### Optimisasi Go

**Flag compiler:**

```bash
-ldflags="-s -w"
```

- `-s`: Hapus tabel simbol (pengurangan ~10%)
- `-w`: Hapus informasi debug DWARF (pengurangan ~20%)

**Pengoptimalan tambahan:**

```bash
-ldflags="-s -w -X main.version=1.0.0"
```

- `-X`: Tetapkan nilai variabel pada waktu build
- Berguna untuk nomor versi dan tanggal build

### Kompresi Biner

**UPX (opsional):**

```bash
# After building
upx --best bin/myapp.exe
```

**Hasil:**

- Pengurangan ukuran ~50%
- Waktu mulai sedikit lebih lambat (~100 ms)
- Tidak disarankan untuk macOS (masalah penandatanganan kode)

## Build Khusus Platform

### Windows

**Keluaran:** `myapp.exe`

**Mencakup:**

- Ikon aplikasi
- Informasi versi
- Manifes (pengaturan UAC)

**Ikon:**

```bash
# Generate platform icons from a source PNG
wails3 generate icons -input appicon.png -windowsfilename build/icon.ico
```

Langkah `tool package` Windows kemudian menyematkan `.ico` yang dihasilkan ke dalam berkas yang dapat dieksekusi.

**Manifes:**

```xml
<!-- build/windows/manifest.xml -->
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly xmlns="urn:schemas-microsoft-com:asm.v1" manifestVersion="1.0">
  <assemblyIdentity version="1.0.0.0" name="MyApp"/>
  <trustInfo xmlns="urn:schemas-microsoft-com:asm.v3">
    <security>
      <requestedPrivileges>
        <requestedExecutionLevel level="asInvoker" uiAccess="false"/>
      </requestedPrivileges>
    </security>
  </trustInfo>
</assembly>
```

### macOS

**Keluaran:** `myapp.app` (bundel aplikasi)

**Struktur:**

```
myapp.app/
├── Contents/
│   ├── Info.plist          # App metadata
│   ├── MacOS/
│   │   └── myapp           # Binary
│   ├── Resources/
│   │   └── icon.icns       # Icon
│   └── _CodeSignature/     # Code signature (if signed)
```

**Info.plist:**

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleName</key>
    <string>My App</string>
    <key>CFBundleIdentifier</key>
    <string>com.example.myapp</string>
    <key>CFBundleVersion</key>
    <string>1.0.0</string>
</dict>
</plist>
```

**Biner Universal:**

Taskfile macOS menyediakan tugas `darwin:build:universal` (dan `darwin:package:universal`) yang membangun kedua arsitektur dan menggabungkannya melalui `wails3 tool lipo`:

```bash
wails3 task darwin:build:universal
```

### Linux

**Keluaran:** `myapp` (biner ELF)

**Dependensi:**

- GTK3
- WebKitGTK

**Berkas desktop:**

```ini
# myapp.desktop
[Desktop Entry]
Name=My App
Exec=/usr/bin/myapp
Icon=myapp
Type=Application
Categories=Utility;
```

**Instalasi:**

```bash
# Copy binary
sudo cp myapp /usr/bin/

# Copy desktop file
sudo cp myapp.desktop /usr/share/applications/

# Copy icon
sudo cp icon.png /usr/share/icons/hicolor/256x256/apps/myapp.png
```

## Performa Build

### Waktu Build Umum

| Fase | Waktu | Catatan |
| --- | --- | --- |
| Analisis | &lt;1 dtk | Pemindaian kode Go |
| Pembuatan Binding | &lt;1 dtk | Pembuatan TypeScript |
| Build Frontend | 5-30 dtk | Bergantung pada ukuran proyek |
| Kompilasi Go | 2-10 dtk | Bergantung pada ukuran kode |
| Penyematan Aset | &lt;1 dtk | Penyematan frontend |
| **Total** | **10-45 dtk** | Build pertama |
| **Inkremental** | **5-15 dtk** | Build berikutnya |

### Mempercepat Build

**1. Gunakan cache build:**

```bash
# Go build cache is automatic
# Frontend cache (Vite)
npm run build  # Uses cache by default
```

**2. Jalankan hanya yang Anda perlukan:**

```bash
# Pick the specific Taskfile target you actually need
wails3 task common:build:frontend   # rebuild only the frontend
wails3 task windows:build           # rebuild only the Windows binary
```

**3. Build paralel (beberapa mesin/CI):**

Kompilasi silang antara Linux/Windows/macOS di v3 umumnya dilakukan dalam kontainer Docker `wails-cross` atau pada runner khusus untuk setiap platform — `wails3 build` sendiri menargetkan OS host. Lihat [Build Lintas Platform](/guides/build/cross-platform/) untuk alur kerja yang didukung.

**4. Gunakan alat yang lebih cepat:**

```bash
# Use esbuild instead of webpack
# (Vite uses esbuild by default)
```

## Pemecahan Masalah

### Build Gagal

**Gejala:** `wails3 build` berhenti dengan pesan kesalahan

**Penyebab umum:**

1. **Kesalahan kompilasi Go**
  ```bash
  # Check Go code compiles
  go build
  ```


2. **Kesalahan build frontend**
  ```bash
  # Check frontend builds
  cd frontend
  npm run build
  ```


3. **Dependensi tidak tersedia**
  ```bash
  # Install dependencies
  npm install
  go mod download
  ```


### Biner Terlalu Besar

**Gejala:** Ukuran biner >50 MB

**Solusi:**

1. **Hapus simbol debug** (Taskfile yang disertakan sudah meneruskan `-ldflags="-s -w"` ke `go build`).

2. **Periksa aset yang disematkan**
  ```bash
  # Remove unnecessary files from frontend/dist/
  # Check for large images, videos, etc.
  ```


3. **Gunakan kompresi UPX**
  ```bash
  upx --best bin/myapp.exe
  ```


### Build Lambat

**Gejala:** Build memerlukan waktu >1 menit

**Solusi:**

1. **Gunakan cache build**
  - Cache Go berjalan secara otomatis
  - Cache frontend (Vite) berjalan secara otomatis


2. **Jalankan hanya tugas yang Anda perlukan**
  ```bash
  wails3 task common:build:frontend
  wails3 task windows:build
  ```


3. **Optimalkan build frontend**
  ```javascript
  // vite.config.js
  export default {
    build: {
      minify: 'esbuild',  // Faster than terser
    },
  }
  ```


## Praktik Terbaik

### ✅ Lakukan

- **Gunakan `wails3 dev` selama pengembangan** — Iterasi cepat
- **Gunakan `wails3 build` untuk rilis** — Keluaran yang dioptimalkan
- **Beri versi pada build Anda** — Gunakan `-ldflags` untuk menyematkan versi
- **Uji build pada platform target** — Kompilasi silang tidak selalu sempurna
- **Pastikan build frontend tetap cepat** — Optimalkan konfigurasi bundler
- **Gunakan cache build** — Mempercepat build berikutnya

### ❌ Jangan Lakukan

- **Jangan commit direktori `build/`** — Tambahkan ke `.gitignore`
- **Jangan melewatkan pengujian build** — Selalu uji sebelum rilis
- **Jangan sematkan aset yang tidak diperlukan** — Jaga agar ukuran biner tetap kecil
- **Jangan gunakan build debug untuk produksi** — Gunakan build yang dioptimalkan
- **Jangan lupakan penandatanganan kode** — Diperlukan untuk distribusi

## Langkah Berikutnya

**Membangun Aplikasi** — Panduan terperinci untuk membangun dan mengemas [Pelajari Selengkapnya →](/guides/build/building/)

**Build Lintas Platform** — Lakukan build untuk semua platform dari satu mesin [Pelajari Selengkapnya →](/guides/build/cross-platform/)

**Membuat Penginstal** — Buat penginstal untuk pengguna akhir [Pelajari Selengkapnya →](/guides/installers/)

---

**Ada pertanyaan tentang proses build?** Tanyakan di [Discord](https://discord.gg/JDdSxwjhGf) atau lihat [contoh build](https://github.com/wailsapp/wails/tree/master/v3/examples/build).
