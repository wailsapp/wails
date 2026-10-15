---
title: "Penyiapan Pengembangan"
description: "Siapkan lingkungan pengembangan Anda untuk pengembangan Wails v3"
slug: "contributing/setup"
sourcePath: "contributing/setup.md"
---

## Penyiapan Lingkungan Pengembangan

Panduan ini memandu Anda menyiapkan lingkungan pengembangan lengkap untuk mengerjakan Wails v3.

## Alat yang Diperlukan

### Pengembangan Go

1. **Instal Go 1.25 atau versi yang lebih baru:**
  ```bash
  # Download from https://go.dev/dl/
  go version  # Verify installation
  ```


2. **Konfigurasikan lingkungan Go:**
  ```bash
  # Add to your shell profile (.bashrc, .zshrc, etc.)
  export GOPATH=$HOME/go
  export PATH=$PATH:$GOPATH/bin
  ```


3. **Instal alat Go yang berguna:**
  ```bash
  go install golang.org/x/tools/cmd/goimports@latest
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
  ```


### Node.js dan npm

Hanya diperlukan untuk contoh integrasi frontend.

```bash
# Install Node.js 20+ and npm
node --version  # Should be 20+
npm --version
```

### Dependensi Khusus Platform

**macOS:**

```bash
# Install Xcode Command Line Tools
xcode-select --install

# Verify installation
xcode-select -p  # Should output a path
```

**Windows:**

1. Instal [MSYS2](https://www.msys2.org/) untuk mendapatkan lingkungan menyerupai Unix
2. WebView2 Runtime (sudah terinstal di Windows 11, [unduh](https://developer.microsoft.com/en-us/microsoft-edge/webview2/) untuk Windows 10)
3. Opsional: Instal [Git for Windows](https://git-scm.com/download/win)

**Linux (Debian/Ubuntu):**

```bash
sudo apt update
# Default GTK4 + WebKitGTK 6.0 stack (Ubuntu 24.04+ / Debian 13+)
sudo apt install build-essential pkg-config libgtk-4-dev libwebkitgtk-6.0-dev
# For the legacy -tags gtk3 path:
sudo apt install libgtk-3-dev libwebkit2gtk-4.1-dev
```

**Linux (Fedora/RHEL):**

```bash
# Default GTK4 stack
sudo dnf install gcc pkg-config gtk4-devel webkitgtk6.0-devel
# Legacy GTK3 path:
sudo dnf install gtk3-devel webkit2gtk4.1-devel
```

**Linux (Arch):**

```bash
# Default GTK4 stack
sudo pacman -S base-devel gtk4 webkitgtk-6.0
# Legacy GTK3 path:
sudo pacman -S gtk3 webkit2gtk-4.1
```

## Penyiapan Repositori

### Kloning dan Konfigurasi

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails

# Add upstream remote
git remote add upstream https://github.com/wailsapp/wails.git

# Verify remotes
git remote -v
```

### Bangun CLI Wails

```bash
# Navigate to v3 directory
cd v3

# Build the CLI
go build -o ../wails3 ./cmd/wails3

# Test the build
cd ..
./wails3 version
```

### Tambahkan ke PATH (Opsional)

**Linux/macOS:**

```bash
# Add to ~/.bashrc or ~/.zshrc
export PATH=$PATH:/path/to/wails
```

**Windows:**

Tambahkan direktori Wails ke variabel lingkungan PATH melalui System Properties.

## Penyiapan IDE

### VS Code (Disarankan)

1. **Instal VS Code:** [Unduh](https://code.visualstudio.com/)

2. **Instal ekstensi:**
  - Go (oleh Go Team di Google)
  - ESLint
  - Prettier
  - MDX (untuk dokumentasi)


3. **Konfigurasikan pengaturan ruang kerja** (`.vscode/settings.json`):
  ```json
  {
    "go.useLanguageServer": true,
    "go.lintTool": "golangci-lint",
    "go.lintOnSave": "workspace",
    "editor.formatOnSave": true,
    "go.formatTool": "goimports"
  }
  ```


### GoLand

1. **Instal GoLand:** [Unduh](https://www.jetbrains.com/go/)

2. **Konfigurasikan:**
  - Aktifkan dukungan modul Go
  - Siapkan pemantau file untuk `goimports`
  - Konfigurasikan gaya kode agar sesuai dengan konvensi proyek


## Verifikasi Penyiapan Anda

Jalankan perintah berikut untuk memastikan semuanya berfungsi:

```bash
# Go version check
go version

# Build Wails
cd v3
go build ./cmd/wails3

# Run tests
go test ./pkg/...

# Create a test app
cd ..
./wails3 init -n mytest -t vanilla
cd mytest
../wails3 dev
```

Jika aplikasi pengujian berhasil dibangun dan dijalankan, lingkungan Anda sudah siap!

## Menjalankan Pengujian

### Pengujian Unit

```bash
cd v3
go test ./...
cd ..
```

### Pengujian Paket Tertentu

```bash
cd v3
go test ./pkg/application
go test ./pkg/events -v  # Verbose output
```

### Jalankan dengan Cakupan Kode

```bash
cd v3
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Jalankan dengan Pendeteksi Race Condition

```bash
cd v3
go test ./... -race
```

## Mengelola Dokumentasi

Dokumentasi Wails v3 ditulis dalam M-Press. File sumber bahasa Inggris berada di  
`docs/mpress/content/`; terjemahan berada di direktori bahasa seperti  
`fr/` dan `id/`.

Pratinjau dan validasi perubahan dokumentasi dari direktori utama repositori:

```bash
go install github.com/leaanthony/mpress/cmd/mpress@v1.0.17
mpress dev
mpress build --strict --no-purge-css
mpress check
# Use `python` on Windows if `python3` is unavailable.
python3 docs/mpress/scripts/check_site.py docs/mpress/site
```

Situs produksi bersifat statis. Node.js, penyedia layanan terjemahan, dan kredensial Cloudflare  
tidak diperlukan untuk mengerjakan dokumentasi secara lokal.

## Debugging

### Men-debug Kode Go

**VS Code:**

Buat `.vscode/launch.json`:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Debug Wails CLI",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/v3/cmd/wails3",
      "args": ["dev"]
    }
  ]
}
```

**Baris Perintah:**

```bash
# Use Delve debugger
go install github.com/go-delve/delve/cmd/dlv@latest
dlv debug ./cmd/wails3 -- dev
```

### Men-debug Kode Platform

Debugging khusus platform memerlukan alat platform:

- **macOS:** Xcode Instruments
- **Windows:** Visual Studio Debugger
- **Linux:** GDB

## Masalah Umum

### "command not found: wails3"

Tambahkan direktori Wails ke PATH Anda atau gunakan `./wails3` dari direktori utama proyek.

### "webkitgtk-6.0 tidak ditemukan" atau "webkit2gtk tidak ditemukan" (Linux)

Instal paket pengembangan untuk stack yang Anda gunakan dalam proses build:

```bash
# Default GTK4 stack (Debian/Ubuntu):
sudo apt install libwebkitgtk-6.0-dev

# Legacy GTK3 path:
sudo apt install libwebkit2gtk-4.1-dev
```

### Build gagal karena kesalahan modul Go

```bash
cd v3
go mod tidy
go mod download
```

### Kesalahan "CGO_ENABLED" di Windows

Pastikan compiler C (MinGW-w64 melalui MSYS2) tersedia di PATH Anda.

## Langkah Berikutnya

- Tinjau [Standar Penulisan Kode](/contributing/standards/)
- Jelajahi [Dokumentasi Teknis](/contributing/)
- Temukan masalah untuk dikerjakan: [Masalah yang Cocok untuk Kontributor Pemula](https://github.com/wailsapp/wails/labels/good%20first%20issue)
