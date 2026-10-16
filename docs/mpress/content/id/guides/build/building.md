---
title: "Membangun Aplikasi"
description: "Bangun dan kemas aplikasi Wails Anda"
slug: "guides/build/building"
sourcePath: "guides/build/building.md"
---

Wails v3 menggunakan [Task](https://taskfile.dev) sebagai sistem build-nya. Perintah `wails3 build` dan `wails3 package` merupakan pembungkus praktis untuk Task.

## Membangun

Bangun untuk platform saat ini:

```bash
wails3 build
```

Bangun untuk platform tertentu:

```bash
wails3 build GOOS=windows
wails3 build GOOS=darwin
wails3 build GOOS=linux

# With architecture
wails3 build GOOS=darwin GOARCH=arm64

# Environment variable style works too
GOOS=windows wails3 build
```

Hasil build disimpan di direktori `bin/`.

@note{type="tip"}
Kompilasi silang ke macOS atau Linux dari platform lain memerlukan Docker. Lihat [Build Lintas Platform](/guides/build/cross-platform/) untuk petunjuk penyiapan.

@end

## Pengembangan

Jalankan aplikasi Anda dengan pemuatan ulang otomatis:

```bash
wails3 dev
```

Tindakan ini memulai pemantau berkas yang membangun ulang dan memulai ulang aplikasi Anda saat ada perubahan. Secara default, server pengembangan frontend berjalan pada port 9245.

```bash
# Custom port
wails3 dev -port 3000

# Enable HTTPS
wails3 dev -s
```

## Pengemasan

Kemas aplikasi Anda untuk didistribusikan:

```bash
wails3 package
wails3 package GOOS=windows
wails3 package GOOS=darwin
wails3 package GOOS=linux
```

Tindakan ini membuat paket khusus platform:

- **Windows**: Penginstal NSIS — lihat [Pengemasan Windows](/guides/build/windows/)
- **macOS**: Bundel aplikasi (`.app`) — lihat [Pengemasan macOS](/guides/build/macos/)
- **Linux**: AppImage, deb, dan rpm — lihat [Pengemasan Linux](/guides/build/linux/)

## Tag Build Kustom

Teruskan tag build Go kustom dengan flag `-tags`:

```bash
# Build with legacy GTK3 + WebKit2GTK 4.1 on Linux (default is GTK4 + WebKitGTK 6.0)
wails3 build -tags gtk3

# Build in server mode (no GUI, CGO-free)
wails3 build -tags server

# Combine multiple tags
wails3 build -tags gtk3,customtag
```

Tag diteruskan sebagai `EXTRA_TAGS` ke Taskfile yang mendasarinya. Lihat [Build Server](/guides/server-build/) dan [Pengemasan Linux - Dukungan GTK3 Lama](/guides/build/linux/#legacy-gtk3-support) untuk detailnya.

## Menggunakan Task Secara Langsung

Untuk kontrol lebih lanjut, gunakan Task secara langsung:

```bash
# List available tasks
wails3 task --list

# Verbose output
wails3 task build -v

# Dry run
wails3 task --dry

# Force rebuild
wails3 task build -f

# Pass variables
wails3 task darwin:build ARCH=amd64
```

Task khusus platform seperti `linux:create:deb` atau `darwin:build:universal` hanya tersedia melalui Task.

## Menghasilkan Aset

Buat ulang ikon atau perbarui konfigurasi build:

```bash
wails3 generate icons -input build/appicon.png
wails3 update build-assets -name "MyApp" -config build/config.yml -dir build
```
