---
title: "Pemaketan MSIX"
description: "Memaketkan aplikasi Wails v3 Anda sebagai paket MSIX"
slug: "guides/build/msix"
sourcePath: "guides/build/msix.md"
---

MSIX adalah format modern untuk pemaketan aplikasi Windows. Wails dapat menghasilkan paket MSIX sebagai bagian dari proses build Windows Anda.

Petunjuk pemaketan MSIX didokumentasikan dalam panduan [Pemaketan Windows](/guides/build/windows/#msix-package).

## Pengemasan langsung melalui CLI

Jalankan alat MSIX di Windows, termasuk dalam CI. Backend bawaan menggunakan `MakeAppx.exe`; penandatanganan juga membutuhkan `signtool.exe`. Keduanya disertakan dalam Windows SDK. Pembantu instalasi membuka Microsoft Store dan, jika diperlukan, halaman unduhan SDK; selesaikan instalasi sebelum membuat paket.

Tetapkan identitas di `build/config.yml`, lalu bangun dan kemas berkas executable Anda:

```yaml
info:
  companyName: "Example Corp"
  productName: "MyApp"
  productIdentifier: "com.example.myapp"
  description: "MyApp"
  version: "1.0.0"
```

```powershell
wails3 tool msix-install-tools
wails3 build GOOS=windows
wails3 tool msix --executable bin/myapp.exe --name myapp.exe
```

Perintah langsung menulis `MyApp.msix` di direktori saat ini. Sebaliknya, [tugas pengemasan Windows](/guides/build/windows/#msix-package) menyediakan jalur keluarannya sendiri.

## Opsi CLI

| Opsi | Arti |
| --- | --- |
| `--config` | Berkas konfigurasi; nilai bawaan `build/config.yml`. |
| `--executable`, `--name` | Berkas executable yang sudah ada dan nama berkasnya di dalam paket; keduanya wajib diisi. |
| `--out` | Berkas keluaran; nilai bawaan `<ProductName>.msix`. |
| `--arch` | Arsitektur paket: `x64` (bawaan), `x86`, `arm`, `arm64`, `x86a64`, atau `neutral`. Alias Go `amd64` dan `386` diterima. Sesuaikan dengan arsitektur executable. |
| `--publisher` | Identitas penerbit; nilai bawaan `CN=<companyName>`. |
| `--cert`, `--cert-password` | Jalur sertifikat PFX dan kata sandi untuk penandatanganan. |
| `--use-makeappx` | Gunakan pengemas bawaan Windows SDK. |
| `--use-msix-tool` | Pilih secara eksplisit `MsixPackagingTool.exe`, yang harus tersedia di `PATH`. |

## Penandatanganan dan CI

Untuk distribusi di luar Store, tanda tangani dengan sertifikat yang dipercaya pada mesin tujuan. Subject sertifikat harus sama persis dengan `--publisher`. Backend MakeAppx memanggil SignTool dengan SHA256 ketika `--cert` diberikan. Lihat [panduan penandatanganan Microsoft](https://learn.microsoft.com/en-us/windows/msix/package/sign-msix-package-guide).

Langkah alur kerja Windows ini mengasumsikan Wails dan SDK sudah terpasang serta langkah sebelumnya telah menyediakan berkas PFX secara aman di `CERT_PATH`. Secret yang hanya berisi jalur tidak mengunggah sertifikat dengan sendirinya:

```yaml
- name: MSIX
  if: runner.os == 'Windows'
  shell: pwsh
  run: |
    wails3 build GOOS=windows
    wails3 tool msix --executable bin/myapp.exe --name myapp.exe --publisher "$env:MSIX_PUBLISHER" --cert "$env:CERT_PATH" --cert-password "$env:CERT_PASSWORD"
  env:
    MSIX_PUBLISHER: ${{ vars.MSIX_PUBLISHER }}
    CERT_PATH: ${{ secrets.WINDOWS_CERT_PATH }}
    CERT_PASSWORD: ${{ secrets.WINDOWS_CERT_PASSWORD }}
```

## Asosiasi berkas dan aset

Tambahkan ekstensi tanpa titik di awal ke `build/config.yml`; manifes yang dihasilkan akan menambahkan titik:

```yaml
fileAssociations:
  - ext: myext
    name: MyApp Document
    description: MyApp Document
    iconName: fileicon
```

Tangani pembukaan berkas saat aplikasi berjalan sebagaimana dijelaskan dalam [Asosiasi Berkas](/guides/file-associations/). Backend MakeAppx saat ini hanya menyalin executable dan menghasilkan gambar pengganti transparan. Backend ini tidak mengimpor berkas `Assets/` proyek atau mengonversi ikon `iconName`. Gunakan alur pengemasan khusus untuk aset bermerek atau DLL tambahan.

| Aset yang dihasilkan | Ukuran (piksel) |
| --- | --- |
| `Square150x150Logo.png` | 150×150 |
| `Square44x44Logo.png` | 44×44 |
| `Wide310x150Logo.png` | 310×150 |
| `StoreLogo.png` | 50×50 |
| `SplashScreen.png` | 620×300 |
| `FileIcon.png` | 44×44 |

`FileIcon.png` hanya dihasilkan jika asosiasi berkas dikonfigurasi. Berkas-berkas ini berada di direktori `Assets/` dalam paket.

## Pengiriman ke Store dan pemecahan masalah

Daftarkan reservasi aplikasi Anda di [portal Partner Center](https://partner.microsoft.com/dashboard) dan gunakan identitas paket serta penerbit dari sana saat menyiapkan pengiriman. Store menandatangani paket MSIX selama proses pengiriman; Anda tidak perlu membeli sertifikat penandatanganan untuk jalur ini. Lihat [persyaratan paket Microsoft](https://learn.microsoft.com/en-us/windows/apps/publish/publish-your-app/msix/app-package-requirements).

Jika `MakeAppx.exe` atau `signtool.exe` tidak ditemukan, pasang atau perbaiki Windows SDK. Wails mencari di `PATH` dan lokasi SDK standar. Untuk kegagalan penandatanganan, periksa Subject sertifikat, masa berlaku, dan kepercayaan pada mesin tujuan; lihat [pemecahan masalah MSIX](https://learn.microsoft.com/en-us/windows/msix/msix-troubleshooting-guide).
