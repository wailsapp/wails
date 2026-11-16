---
title: "Pemaketan Windows"
description: "Kemas aplikasi Wails Anda untuk didistribusikan di Windows"
slug: "guides/build/windows"
sourcePath: "guides/build/windows.md"
---

## Penginstal NSIS

Format pemaketan default membuat penginstal NSIS:

```bash
wails3 package GOOS=windows
```

Perintah ini menjalankan `wails3 task windows:package`, yang:

1. Mem-build aplikasi
2. Membuat bootstrapper WebView2
3. Membuat penginstal NSIS

Keluaran: `bin/<AppName>-<arch>-installer.exe`. Skrip sumber NSIS tetap berada di `build/windows/nsis/`.

### Paket MSIX

Untuk distribusi melalui Microsoft Store atau penerapan Windows modern:

```bash
wails3 package GOOS=windows FORMAT=msix
```

Keluaran: `bin/<AppName>-<arch>.msix`

@note{type="note"}
MSIX memerlukan `makeappx.exe` (Windows SDK) atau alat MSIX mandiri. Taskfile Windows menyediakan tugas instalasi sebagai `wails3 task install:msix:tools`.

@end

## Menyesuaikan Penginstal

Konfigurasi NSIS berada di `build/windows/nsis/project.nsi`. Edit file ini untuk menyesuaikan:

- Antarmuka pengguna dan pencitraan merek penginstal
- Direktori instalasi
- Pintasan menu Start dan desktop
- Asosiasi file
- Perjanjian lisensi

Metadata aplikasi berasal dari `build/windows/info.json`:

```json
{
  "fixed": {
    "file_version": "1.0.0"
  },
  "info": {
    "0409": {
      "ProductVersion": "1.0.0",
      "CompanyName": "My Company",
      "FileDescription": "My Application",
      "ProductName": "MyApp"
    }
  }
}
```

## Penandatanganan Kode

Tandatangani file executable dan penginstal Anda untuk menghindari peringatan SmartScreen:

```bash
# Using the wrapper (auto-detects platform)
wails3 sign GOOS=windows

# Or using tasks directly
wails3 task windows:sign
wails3 task windows:sign:installer
```

Konfigurasikan penandatanganan di `build/windows/Taskfile.yml`:

```yaml
vars:
  SIGN_CERTIFICATE: "path/to/certificate.pfx"
  # Or use thumbprint for certificates in Windows store
  SIGN_THUMBPRINT: "certificate-thumbprint"
  TIMESTAMP_SERVER: "http://timestamp.digicert.com"
```

Simpan kata sandi sertifikat Anda dengan aman:

```bash
wails3 signing
```

Lihat [Menandatangani Aplikasi](/guides/build/signing/) untuk detail selengkapnya.

## Mem-build untuk ARM

```bash
wails3 build GOOS=windows GOARCH=arm64
wails3 package GOOS=windows GOARCH=arm64
```

## Pemecahan Masalah

### makensis tidak ditemukan

Instal NSIS:

```bash
# Windows
winget install NSIS.NSIS

# Or download from https://nsis.sourceforge.io/
```

### Peringatan SmartScreen

File executable Anda belum ditandatangani. Lihat [Penandatanganan Kode](#penandatanganan-kode) di atas.

### WebView2 tidak ditemukan

Penginstal menyertakan bootstrapper WebView2 yang mengunduh runtime jika diperlukan. Jika Anda memerlukan instalasi offline, unduh Evergreen Standalone Installer dari Microsoft.
