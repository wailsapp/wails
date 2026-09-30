---
title: "Build Lintas Platform"
description: "Build untuk beberapa platform dari satu mesin"
slug: "guides/build/cross-platform"
sourcePath: "guides/build/cross-platform.md"
---

## Mulai Cepat

Wails v3 mendukung build untuk Windows, macOS, dan Linux dari sistem operasi host apa pun. Sistem build secara otomatis mendeteksi lingkungan Anda dan memilih metode kompilasi yang tepat.

**Ingin melakukan kompilasi silang ke macOS dan Linux?** Jalankan perintah ini sekali untuk menyiapkan image Docker (unduhan sekitar 800 MB):

```bash
wails3 task setup:docker
```

Kemudian, lakukan build untuk platform apa pun:

```bash
# Build for current platform (production by default)
wails3 build

# Build for specific platforms
wails3 build GOOS=windows
wails3 build GOOS=darwin
wails3 build GOOS=linux

# Build for ARM64 architecture
wails3 build GOOS=windows GOARCH=arm64
wails3 build GOOS=darwin GOARCH=arm64
wails3 build GOOS=linux GOARCH=arm64

# Environment variable style also works
GOOS=darwin GOARCH=arm64 wails3 build
```

### Windows

Windows adalah target kompilasi silang yang paling sederhana karena secara default tidak memerlukan CGO.

```bash
wails3 build GOOS=windows
```

Ini dapat dilakukan dari OS host apa pun tanpa penyiapan tambahan. Fitur kompilasi silang bawaan Go menangani semuanya.

**Jika aplikasi Anda memerlukan CGO** (misalnya, Anda menggunakan pustaka C atau paket yang bergantung pada CGO), Anda memerlukan Docker saat melakukan build dari macOS atau Linux:

```bash
# One-time setup
wails3 task setup:docker

# Build with CGO enabled
wails3 task windows:build CGO_ENABLED=1
```

Taskfile mendeteksi `CGO_ENABLED=1` pada host selain Windows dan secara otomatis menggunakan image Docker.

### macOS

Build macOS memerlukan CGO untuk integrasi WebView, yang berarti kompilasi silang memerlukan alat khusus.

```bash
# Build for Apple Silicon (arm64) - default
wails3 build GOOS=darwin

# Build for Intel (amd64)
wails3 build GOOS=darwin GOARCH=amd64

# Build universal binary (both architectures)
wails3 task darwin:build:universal
```

**Dari Linux atau Windows**, Anda harus menyiapkan Docker terlebih dahulu:

```bash
wails3 task setup:docker
```

Setelah image selesai dibuat, sistem build mendeteksi bahwa Anda tidak menggunakan macOS dan secara otomatis menggunakan Docker. Anda tidak perlu mengubah perintah build.

Perhatikan bahwa biner macOS hasil kompilasi silang belum memiliki tanda tangan digital kode. Anda harus menandatanganinya di macOS atau dalam CI sebelum didistribusikan.

### Linux

Build Linux memerlukan CGO untuk integrasi WebView.

```bash
wails3 build GOOS=linux

# Build for specific architecture
wails3 build GOOS=linux GOARCH=amd64
wails3 build GOOS=linux GOARCH=arm64
```

**Dari macOS atau Windows**, Anda harus menyiapkan Docker terlebih dahulu:

```bash
wails3 task setup:docker
```

Sistem build mendeteksi bahwa Anda tidak menggunakan Linux dan secara otomatis menggunakan Docker.

**Di Linux tanpa kompiler C**, sistem build memeriksa keberadaan `gcc` atau `clang`. Jika keduanya tidak ditemukan, sistem beralih menggunakan Docker. Ini berguna untuk kontainer minimal atau sistem yang tidak memiliki alat build terinstal. Anda dapat:

1. Menginstal kompiler C: `sudo apt install build-essential` (Debian/Ubuntu) atau `sudo pacman -S base-devel` (Arch)
2. Membuat image Docker dan membiarkannya digunakan secara otomatis

### Arsitektur ARM

Semua platform mendukung kompilasi silang ARM64 menggunakan `GOARCH`:

```bash
# Windows ARM64 (Surface Pro X, Windows on ARM)
wails3 build GOOS=windows GOARCH=arm64

# Linux ARM64 (Raspberry Pi 4/5, AWS Graviton)
wails3 build GOOS=linux GOARCH=arm64

# macOS ARM64 (Apple Silicon - this is the default on macOS)
wails3 build GOOS=darwin GOARCH=arm64

# macOS Intel (amd64)
wails3 build GOOS=darwin GOARCH=amd64
```

Image Docker menyertakan target kompiler silang Zig untuk amd64 dan arm64 di semua platform, sehingga build ARM dapat dilakukan dari host apa pun:

| Build ARM64 untuk | Dari Windows | Dari macOS | Dari Linux |
| --- | --- | --- | --- |
| **Windows ARM64** | Go Native | Go Native | Go Native |
| **macOS ARM64** | Docker | Native | Docker |
| **Linux ARM64** | Docker | Docker | Docker_ |

<em>Linux ARM64 dari Linux x86</em>64 menggunakan Docker karena kompilasi silang CGO memerlukan toolchain yang berbeda.

## Cara Kerjanya

### Matriks Kompilasi Silang

| Host → Sasaran | Windows | macOS | Linux |
| --- | --- | --- | --- |
| **Windows** | Native | Docker | Docker |
| **macOS** | Go Native | Native | Docker |
| **Linux** | Go Native | Docker | Native |

- **Native** = Toolchain native platform, tanpa penyiapan tambahan
- **Go Native** = Kompilasi silang bawaan Go (`CGO_ENABLED=0`)
- **Docker** = Image Docker dengan kompiler silang Zig

### Persyaratan CGO

| Target | Memerlukan CGO | Metode Kompilasi Silang |
| --- | --- | --- |
| Windows | Tidak (secara default) | Go Native. Docker hanya jika `CGO_ENABLED=1` |
| macOS | Ya | Docker dengan SDK macOS |
| Linux | Ya | Docker, atau secara native jika kompiler C tersedia |

### Deteksi Otomatis

Taskfile secara otomatis memilih metode build yang tepat berdasarkan lingkungan Anda:

- **Target Windows:** Secara default menggunakan kompilasi silang Go native. Jika Anda secara eksplisit menetapkan `CGO_ENABLED=1` pada host non-Windows, metode akan beralih ke Docker.
- **Target macOS:** Secara otomatis menggunakan Docker jika tidak dijalankan di macOS. Tidak diperlukan intervensi manual.
- **Target Linux:** Memeriksa keberadaan `gcc` atau `clang`. Menggunakan kompilasi native jika ditemukan; jika tidak, beralih ke Docker.

### Image Docker

Wails menggunakan satu image Docker (`wails-cross`) yang dapat melakukan build untuk semua platform. Image ini menggunakan [Zig](https://ziglang.org/) sebagai kompiler silang, sehingga dapat menargetkan platform apa pun dari host apa pun. SDK macOS disertakan untuk target darwin.

```bash
wails3 task setup:docker
```

Anda dapat memeriksa apakah image sudah siap dengan menjalankan `wails3 doctor`.

### SDK macOS

Image Docker mengunduh SDK macOS dari [wailsapp/macosx-sdks](https://github.com/wailsapp/macosx-sdks) selama proses build image. SDK ini diperlukan karena header macOS dibutuhkan untuk kompilasi CGO.

**Penting:** Wails tidak mendistribusikan SDK macOS. Pengguna bertanggung jawab untuk meninjau ketentuan lisensi SDK Apple sebelum menggunakan fitur ini.

## Build Image Anda Sendiri

Jika perlu menyesuaikan image Docker (misalnya, menggunakan versi SDK macOS yang berbeda, menambahkan alat lain, atau menggunakan SDK Anda sendiri), Anda dapat melakukan build image tersebut sendiri.

### Dockerfile

Buat `Dockerfile` dengan konten berikut:

```dockerfile
# syntax=docker/dockerfile:1
FROM golang:1.24-alpine

ARG ZIG_VERSION=0.14.0
ARG MACOS_SDK_VERSION=14.5
ARG IMAGE_VERSION=1.0.0

LABEL org.opencontainers.image.title="Wails Cross-Compiler"
LABEL org.opencontainers.image.description="Cross-compile Wails v3 apps to macOS, Linux, and Windows"
LABEL org.opencontainers.image.source="https://github.com/wailsapp/wails"
LABEL org.opencontainers.image.vendor="Wails"
LABEL org.opencontainers.image.version="${IMAGE_VERSION}"
LABEL io.wails.sdk.version="${MACOS_SDK_VERSION}"
LABEL io.wails.zig.version="${ZIG_VERSION}"

RUN apk add --no-cache curl xz nodejs npm gcompat

RUN curl -L "https://ziglang.org/download/${ZIG_VERSION}/zig-linux-x86_64-${ZIG_VERSION}.tar.xz" \
    | tar -xJ -C /opt \
    && ln -s /opt/zig-linux-x86_64-${ZIG_VERSION}/zig /usr/local/bin/zig

RUN curl -fL --retry 3 --retry-delay 5 -o /tmp/sdk.tar.xz \
    "https://github.com/wailsapp/macosx-sdks/releases/download/${MACOS_SDK_VERSION}/MacOSX${MACOS_SDK_VERSION}.sdk.tar.xz" \
    && tar -xJf /tmp/sdk.tar.xz -C /opt \
    && mv /opt/MacOSX${MACOS_SDK_VERSION}.sdk /opt/macos-sdk \
    && rm /tmp/sdk.tar.xz

ENV MACOS_SDK_PATH=/opt/macos-sdk

# Create zig cc wrappers for each target
# Darwin arm64
COPY <<'ZIGWRAP' /usr/local/bin/zcc-darwin-arm64
#!/bin/sh
ARGS=""
SKIP_NEXT=0
for arg in "$@"; do
    if [ $SKIP_NEXT -eq 1 ]; then
        SKIP_NEXT=0
        continue
    fi
    case "$arg" in
        -target) SKIP_NEXT=1 ;;
        -mmacosx-version-min=*) ;;
        *) ARGS="$ARGS $arg" ;;
    esac
done
exec zig cc -target aarch64-macos-none -isysroot /opt/macos-sdk -I/opt/macos-sdk/usr/include -L/opt/macos-sdk/usr/lib -F/opt/macos-sdk/System/Library/Frameworks -w $ARGS
ZIGWRAP
RUN chmod +x /usr/local/bin/zcc-darwin-arm64

# Darwin amd64
COPY <<'ZIGWRAP' /usr/local/bin/zcc-darwin-amd64
#!/bin/sh
ARGS=""
SKIP_NEXT=0
for arg in "$@"; do
    if [ $SKIP_NEXT -eq 1 ]; then
        SKIP_NEXT=0
        continue
    fi
    case "$arg" in
        -target) SKIP_NEXT=1 ;;
        -mmacosx-version-min=*) ;;
        *) ARGS="$ARGS $arg" ;;
    esac
done
exec zig cc -target x86_64-macos-none -isysroot /opt/macos-sdk -I/opt/macos-sdk/usr/include -L/opt/macos-sdk/usr/lib -F/opt/macos-sdk/System/Library/Frameworks -w $ARGS
ZIGWRAP
RUN chmod +x /usr/local/bin/zcc-darwin-amd64

# Linux amd64
COPY <<'ZIGWRAP' /usr/local/bin/zcc-linux-amd64
#!/bin/sh
ARGS=""
SKIP_NEXT=0
for arg in "$@"; do
    if [ $SKIP_NEXT -eq 1 ]; then
        SKIP_NEXT=0
        continue
    fi
    case "$arg" in
        -target) SKIP_NEXT=1 ;;
        *) ARGS="$ARGS $arg" ;;
    esac
done
exec zig cc -target x86_64-linux-musl $ARGS
ZIGWRAP
RUN chmod +x /usr/local/bin/zcc-linux-amd64

# Linux arm64
COPY <<'ZIGWRAP' /usr/local/bin/zcc-linux-arm64
#!/bin/sh
ARGS=""
SKIP_NEXT=0
for arg in "$@"; do
    if [ $SKIP_NEXT -eq 1 ]; then
        SKIP_NEXT=0
        continue
    fi
    case "$arg" in
        -target) SKIP_NEXT=1 ;;
        *) ARGS="$ARGS $arg" ;;
    esac
done
exec zig cc -target aarch64-linux-musl $ARGS
ZIGWRAP
RUN chmod +x /usr/local/bin/zcc-linux-arm64

# Windows amd64
COPY <<'ZIGWRAP' /usr/local/bin/zcc-windows-amd64
#!/bin/sh
ARGS=""
SKIP_NEXT=0
for arg in "$@"; do
    if [ $SKIP_NEXT -eq 1 ]; then
        SKIP_NEXT=0
        continue
    fi
    case "$arg" in
        -target) SKIP_NEXT=1 ;;
        -Wl,*) ;;
        *) ARGS="$ARGS $arg" ;;
    esac
done
exec zig cc -target x86_64-windows-gnu $ARGS
ZIGWRAP
RUN chmod +x /usr/local/bin/zcc-windows-amd64

# Windows arm64
COPY <<'ZIGWRAP' /usr/local/bin/zcc-windows-arm64
#!/bin/sh
ARGS=""
SKIP_NEXT=0
for arg in "$@"; do
    if [ $SKIP_NEXT -eq 1 ]; then
        SKIP_NEXT=0
        continue
    fi
    case "$arg" in
        -target) SKIP_NEXT=1 ;;
        -Wl,*) ;;
        *) ARGS="$ARGS $arg" ;;
    esac
done
exec zig cc -target aarch64-windows-gnu $ARGS
ZIGWRAP
RUN chmod +x /usr/local/bin/zcc-windows-arm64

# Build script
COPY <<'SCRIPT' /usr/local/bin/build.sh
#!/bin/sh
set -e

OS=${1:-darwin}
ARCH=${2:-arm64}

case "${OS}-${ARCH}" in
    darwin-arm64|darwin-aarch64) export CC=zcc-darwin-arm64; export GOARCH=arm64; export GOOS=darwin ;;
    darwin-amd64|darwin-x86_64)  export CC=zcc-darwin-amd64; export GOARCH=amd64; export GOOS=darwin ;;
    linux-arm64|linux-aarch64)   export CC=zcc-linux-arm64;  export GOARCH=arm64; export GOOS=linux ;;
    linux-amd64|linux-x86_64)    export CC=zcc-linux-amd64;  export GOARCH=amd64; export GOOS=linux ;;
    windows-arm64|windows-aarch64) export CC=zcc-windows-arm64; export GOARCH=arm64; export GOOS=windows ;;
    windows-amd64|windows-x86_64)  export CC=zcc-windows-amd64; export GOARCH=amd64; export GOOS=windows ;;
    *) echo "Usage: <os> <arch>"; echo "  os: darwin, linux, windows"; echo "  arch: amd64, arm64"; exit 1 ;;
esac

export CGO_ENABLED=1
export CGO_CFLAGS="-w"

# Build frontend if exists and not already built (host may have built it)
if [ -d "frontend" ] && [ -f "frontend/package.json" ] && [ ! -d "frontend/dist" ]; then
    (cd frontend && npm install --silent && npm run build --silent)
fi

# Build
APP=${APP_NAME:-$(basename $(pwd))}
mkdir -p bin

EXT=""
LDFLAGS="-s -w"
if [ "$GOOS" = "windows" ]; then
    EXT=".exe"
    LDFLAGS="-s -w -H windowsgui"
fi

go build -ldflags="$LDFLAGS" -o bin/${APP}-${GOOS}-${GOARCH}${EXT} .
echo "Built: bin/${APP}-${GOOS}-${GOARCH}${EXT}"
SCRIPT
RUN chmod +x /usr/local/bin/build.sh

WORKDIR /app
ENTRYPOINT ["/usr/local/bin/build.sh"]
CMD ["darwin", "arm64"]
```

### Melakukan Build Image

Simpan Dockerfile, lalu lakukan build image:

```bash
docker build -t wails-cross .
```

### Menggunakan Versi SDK yang Berbeda

Ubah argumen build `MACOS_SDK_VERSION`:

```bash
docker build -t wails-cross --build-arg MACOS_SDK_VERSION=15.0 .
```

Lihat [versi SDK yang tersedia](https://github.com/wailsapp/macosx-sdks/releases) untuk mengetahui pilihannya.

### Menggunakan SDK Anda Sendiri

Jika memiliki SDK macOS sendiri (misalnya, diekstrak dari Xcode), Anda dapat mengubah Dockerfile agar menggunakan file lokal alih-alih mengunduhnya:

```dockerfile
# Replace the curl/tar SDK download section with:
# (Replace 14.5 with your SDK version in both the COPY and mv commands)
COPY MacOSX14.5.sdk.tar.xz /tmp/sdk.tar.xz
RUN tar -xJf /tmp/sdk.tar.xz -C /opt \
    && mv /opt/MacOSX14.5.sdk /opt/macos-sdk \
    && rm /tmp/sdk.tar.xz
```

Tempatkan tarball SDK Anda di direktori yang sama dengan Dockerfile, lalu lakukan build.

## Integrasi CI/CD

Untuk rilis produksi, sebaiknya gunakan CI/CD dengan runner native untuk setiap platform. Cara ini sepenuhnya menghindari kompilasi silang dan memastikan Anda memperoleh berkas biner yang ditandatangani dengan benar.

```yaml
name: Build

on:
  push:
    branches: [main]

jobs:
  build:
    strategy:
      matrix:
        include:
          - os: ubuntu-latest
            goos: linux
          - os: macos-latest
            goos: darwin
          - os: windows-latest
            goos: windows

    runs-on: ${{ matrix.os }}

    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'

      - uses: actions/setup-node@v4
        with:
          node-version: '20'

      - name: Install Wails CLI
        run: go install github.com/wailsapp/wails/v3/cmd/wails3@latest

      - name: Install Task
        uses: arduino/setup-task@v2

      - name: Build
        run: wails3 build

      - uses: actions/upload-artifact@v4
        with:
          name: app-${{ matrix.goos }}
          path: bin/
```

## Pemecahan Masalah

### Image Docker tidak ditemukan

```
Docker image 'wails-cross' not found.
```

Jalankan `wails3 task setup:docker` untuk melakukan build image Docker. Anda hanya perlu melakukannya sekali.

### Daemon Docker tidak berjalan

```
Docker is required for cross-compilation. Please install Docker.
```

Jalankan Docker Desktop atau daemon Docker. Di Linux, Anda mungkin perlu menjalankan `sudo systemctl start docker`.

### Tidak ada kompiler C di Linux

Jika Anda melihat kesalahan terkait CGO saat melakukan build di Linux, tersedia dua opsi:

1. **Instal kompiler C:**
  - Debian/Ubuntu: `sudo apt install build-essential`
  - Arch Linux: `sudo pacman -S base-devel`
  - Fedora: `sudo dnf install gcc`


2. **Gunakan Docker sebagai gantinya:** Jalankan `wails3 task setup:docker`; Taskfile akan menggunakannya secara otomatis jika tidak ada kompiler yang terdeteksi.

### Berkas biner macOS tidak ditandatangani

Berkas biner macOS hasil kompilasi silang tidak ditandatangani dengan tanda tangan kode. Apple mewajibkan penandatanganan kode untuk distribusi, jadi Anda perlu:

1. Menandatangani berkas biner di mesin macOS, atau
2. Menandatanganinya di CI menggunakan runner macOS

Lihat [Menandatangani Aplikasi](/guides/build/signing/) untuk detailnya.

### Pembuatan berkas biner universal

Berkas biner universal (gabungan arm64 + amd64) dapat dibuat di platform apa pun:

```bash
wails3 task darwin:build:universal
```

Di Linux dan Windows, Wails menggunakan perintah bawaan `wails3 tool lipo` (yang didukung oleh [konoui/lipo](https://github.com/konoui/lipo)) untuk menggabungkan berkas biner. Proses ini menghasilkan satu berkas biner yang berjalan secara native di Mac Apple Silicon maupun Intel.

## Langkah Berikutnya

- [Melakukan Build Aplikasi](/guides/build/building/) - Perintah dan opsi build dasar
- [Menandatangani Aplikasi](/guides/build/signing/) - Penandatanganan kode untuk distribusi
