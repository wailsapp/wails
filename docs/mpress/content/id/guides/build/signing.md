---
title: "Penandatanganan Kode"
description: "Panduan untuk menandatangani aplikasi Wails Anda di semua platform"
slug: "guides/build/signing"
sourcePath: "guides/build/signing.md"
---

## Menandatangani Kode Aplikasi Anda

Panduan ini menjelaskan cara menandatangani aplikasi Wails Anda untuk macOS, Windows, dan Linux. Wails v3 menyediakan alat CLI bawaan untuk penandatanganan kode, notarisasi, dan pengelolaan kunci PGP.

- **macOS** - Tandatangani dan notarisasikan aplikasi macOS Anda
- **Windows** - Tandatangani berkas yang dapat dieksekusi dan paket Windows Anda
- **Linux** - Tandatangani paket DEB dan RPM dengan kunci PGP

## Matriks Penandatanganan Lintas Platform

Matriks ini menunjukkan format yang dapat Anda tandatangani dari setiap platform sumber:

| Format Target | Dari Windows | Dari macOS | Dari Linux |
| --- | :---: | :---: | :---: |
| EXE/MSI Windows | ✅ | ✅ | ✅ |
| Bundel .app macOS | ❌ | ✅ | ❌ |
| Notarisasi macOS | ❌ | ✅ | ❌ |
| DEB Linux | ✅ | ✅ | ✅ |
| RPM Linux | ✅ | ✅ | ✅ |

@note{type="tip"}
Paket Windows dan Linux dapat ditandatangani dari **platform apa pun**. Penandatanganan macOS memerlukan Mac karena persyaratan alat dari Apple.

@end

### Backend Penandatanganan

Wails secara otomatis memilih backend penandatanganan terbaik yang tersedia:

| Platform | Backend Native | Backend Lintas Platform |
| --- | --- | --- |
| Windows | `signtool.exe` (Windows SDK) | Bawaan |
| macOS | `codesign` (Xcode) | Tidak tersedia |
| Linux | Tidak berlaku | Bawaan |

Saat dijalankan di platform native, Wails menggunakan alat native untuk kompatibilitas maksimum. Saat melakukan kompilasi silang, Wails menggunakan dukungan penandatanganan bawaan.

## Mulai Cepat

Cara tercepat untuk mengonfigurasi penandatanganan adalah menggunakan wizard penyiapan:

```bash
wails3 setup
```

Tindakan ini menulis konfigurasi penandatanganan bersama ke `~/.config/wails/defaults.yaml` yang akan **digunakan pada saat penandatanganan di setiap platform** (lihat [Prioritas konfigurasi](#prioritas-konfigurasi)). Langkah penandatanganannya:

- Mendeteksi alat penandatanganan yang terinstal untuk platform target Anda — **dari host apa pun** — dan menampilkan perintah instalasi untuk OS Anda (misalnya `brew install gnupg`, `sudo apt install osslsigncode`, `winget install GnuPG.Gpg4win`).
- Di macOS, menampilkan daftar sertifikat Developer ID dari rantai kunci Anda.
- Untuk Linux, menampilkan daftar kunci GPG Anda dan dapat **membuat serta mengekspor** kunci baru.
- Untuk Windows, dapat **membuat sertifikat yang ditandatangani sendiri** (untuk pengujian) melalui OpenSSL.
- **Menyimpan kata sandi dengan aman di penyimpanan kredensial sistem Anda** (bukan di Taskfile).

@note{type="tip"}
Kata sandi disimpan di penyimpanan kredensial native sistem Anda (macOS Keychain, Windows Credential Manager, atau Linux Secret Service). Dengan demikian, konfigurasi penandatanganan Anda aman dan berfungsi di semua proyek Wails Anda.

@end

## Prioritas konfigurasi

Saat Anda menjalankan tugas penandatanganan, setiap opsi penandatanganan ditentukan dengan urutan berikut (kecocokan pertama yang digunakan):

1. Flag eksplisit yang diteruskan ke `wails3 tool sign` (misalnya `--pgp-key`, `--certificate`, `--identity`).
2. **Variabel Taskfile proyek** yang sesuai (`PGP_KEY`, `SIGN_CERTIFICATE`/`SIGN_THUMBPRINT`, `SIGN_IDENTITY`, …).
3. **Konfigurasi global** di `~/.config/wails/defaults.yaml` (ditulis oleh `wails3 setup`).

Artinya, variabel Taskfile merupakan **penggantian opsional**: jika suatu variabel tidak ditetapkan, kunci/sertifikat/identitas yang dikonfigurasi secara global akan digunakan. Jika tidak satu pun dari ketiganya memberikan nilai, perintah penandatanganan akan melaporkan kesalahan yang jelas dan memberi tahu Anda cara mengonfigurasinya.

## Konfigurasi per proyek

Untuk mengonfigurasi penandatanganan bagi satu proyek alih-alih secara global, jalankan wizard per proyek dari dalam proyek tersebut—wizard ini menulis `vars` ke berkas `build/<platform>/Taskfile.yml` proyek tersebut:

```bash
wails3 setup signing                                   # all detected platforms
wails3 setup signing --platform windows --platform linux
```

### Konfigurasi Manual

Sebagai alternatif, Anda dapat mengedit Taskfile khusus platform secara manual. Edit bagian `vars` di bagian atas setiap berkas:

@tabs
[macOS]
Edit `build/darwin/Taskfile.yml`:

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  KEYCHAIN_PROFILE: "my-notarize-profile"
  # ENTITLEMENTS: "build/darwin/entitlements.plist"
```

Kemudian jalankan:

```bash
wails3 task darwin:sign           # Sign only
wails3 task darwin:sign:notarize  # Sign and notarize
```

[Windows]
Edit `build/windows/Taskfile.yml`:

```yaml
vars:
  SIGN_CERTIFICATE: "path/to/certificate.pfx"
  # Or use thumbprint instead:
  # SIGN_THUMBPRINT: "certificate-thumbprint"
  # TIMESTAMP_SERVER: "http://timestamp.digicert.com"
```

Kata sandi diambil dari rantai kunci sistem (jalankan `wails3 setup signing` untuk mengonfigurasinya).

Kemudian jalankan:

```bash
wails3 task windows:sign           # Sign executable
wails3 task windows:sign:installer # Sign NSIS installer
```

[Linux]
Edit `build/linux/Taskfile.yml`:

```yaml
vars:
  PGP_KEY: "path/to/signing-key.asc"
  # SIGN_ROLE: "builder"  # Options: origin, maint, archive, builder
```

Kata sandi diambil dari rantai kunci sistem (jalankan `wails3 setup signing` untuk mengonfigurasinya).

Kemudian jalankan:

```bash
wails3 task linux:sign:deb       # Sign DEB package
wails3 task linux:sign:rpm       # Sign RPM package
wails3 task linux:sign:packages  # Sign all packages
```

@end

Anda juga dapat memeriksa status penandatanganan secara langsung dari sistem:

```bash
# List available macOS code-signing identities
security find-identity -v -p codesigning

# List PGP keys (Linux package signing)
gpg --list-keys
```

Untuk menyiapkan konfigurasi penandatanganan setiap platform secara interaktif, gunakan wizard:

```bash
wails3 setup signing
```

## Penandatanganan Kode macOS

### Prasyarat

- Akun Apple Developer ($99/tahun)
- Sertifikat Developer ID Application
- Xcode Command Line Tools telah terinstal

### Identitas Penandatanganan

Periksa identitas penandatanganan yang tersedia:

```bash
security find-identity -v -p codesigning
```

Keluaran:

```
Found 2 signing identities:

  Developer ID Application: Your Company (ABCD1234) [valid]
    Hash: ABC123DEF456...

  Apple Development: your@email.com (XYZ789) [valid]
    Hash: DEF789ABC123...
```

@note{type="tip"}
Untuk distribusi di luar App Store, Anda memerlukan sertifikat **Developer ID Application**.

@end

### Konfigurasi

Edit `build/darwin/Taskfile.yml` dan tetapkan variabel penandatanganan:

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  KEYCHAIN_PROFILE: "my-notarize-profile"
  ENTITLEMENTS: "build/darwin/entitlements.plist"
```

| Variabel | Wajib | Deskripsi |
| --- | --- | --- |
| `SIGN_IDENTITY` | Ya | Developer ID Anda (misalnya, "Developer ID Application: Perusahaan Anda (TEAMID)") |
| `KEYCHAIN_PROFILE` | Untuk notarisasi | Nama profil rantai kunci yang menyimpan kredensial |
| `ENTITLEMENTS` | Tidak | Path ke file entitlement |

Kemudian jalankan:

```bash
wails3 task darwin:sign           # Build, package, and sign
wails3 task darwin:sign:notarize  # Build, package, sign, and notarize
```

### Entitlement

Entitlement mengontrol kapabilitas yang dapat diakses aplikasi Anda. Aplikasi Wails biasanya memerlukan entitlement yang berbeda untuk pengembangan dan produksi:

- **Pengembangan**: Memerlukan entitlement JIT, memori tanpa tanda tangan, dan debugging
- **Produksi**: Entitlement minimal (hanya akses jaringan)

Gunakan wizard penyiapan interaktif untuk membuat kedua file:

```bash
wails3 setup entitlements
```

Tindakan ini membuat:

- `build/darwin/entitlements.dev.plist` - Untuk build pengembangan
- `build/darwin/entitlements.plist` - Untuk build produksi/bertanda tangan

**Preset yang tersedia:**\

| Preset | Deskripsi |
| --- | --- |
| Pengembangan | JIT, memori tanpa tanda tangan, debugging, jaringan |
| Produksi | Hanya jaringan (minimal, paling aman) |
| Keduanya | Membuat file pengembangan dan produksi (direkomendasikan) |
| App Store | Sandbox diaktifkan dengan akses jaringan dan file |
| Kustom | Pilih entitlement satu per satu |

@note{type="note"}
Tugas `run` dalam Taskfile darwin menggunakan `entitlements.dev.plist` secara otomatis. Tugas `sign` menggunakan `entitlements.plist` untuk build produksi.

@end

Kemudian tetapkan `ENTITLEMENTS` dalam variabel Taskfile Anda agar mengarah ke file yang sesuai.

### Notarisasi

Apple mewajibkan semua aplikasi yang didistribusikan untuk dinotarisasi.

@steps
### **Simpan kredensial Anda di rantai kunci** (penyiapan satu kali). Jalankan `wails3 setup signing` (perintah ini meminta nilai-nilai tersebut dan memanggil `notarytool` di balik layar), atau panggil `notarytool` secara langsung:
```bash
xcrun notarytool store-credentials "my-notarize-profile" \
  --apple-id "your@email.com" \
  --team-id "ABCD1234" \
  --password "app-specific-password"
```

### **Tetapkan KEYCHAIN_PROFILE dalam Taskfile Anda** agar sesuai dengan nama profil di atas.
### **Tandatangani dan notarisasikan aplikasi Anda**:
```bash
wails3 task darwin:sign:notarize
```

### **Verifikasi notarisasi**:
```bash
spctl --assess --verbose=2 bin/MyApp.app
```

@end

@note{type="note"}
Notarisasi biasanya memerlukan waktu 1-2 menit. Tiket secara otomatis dilampirkan ke aplikasi Anda.

@end

## Penandatanganan Kode Windows

### Prasyarat

- Sertifikat penandatanganan kode (dari DigiCert, Sectigo, dan sebagainya)
- Untuk penandatanganan native di Windows: Windows SDK telah terinstal (untuk `signtool.exe`)
- Untuk penandatanganan lintas platform dari macOS/Linux: [`osslsigncode`](https://github.com/mtrojnar/osslsigncode) (langkah penandatanganan `wails3 setup` menampilkan perintah instalasi khusus untuk host tersebut)

### Membuat sertifikat yang ditandatangani sendiri (pengujian)

Jika Anda hanya perlu menguji alur penandatanganan, jalankan `wails3 setup`, buka tab **Windows** pada langkah penandatanganan, lalu pilih **Buat sertifikat yang ditandatangani sendiri**. Tindakan ini menggunakan OpenSSL untuk membuat `.pfx` penandatanganan kode dan mencatat lokasinya dalam konfigurasi global.

@note{type="caution"}
Sertifikat yang ditandatangani sendiri hanya untuk **pengujian dan distribusi internal**—sertifikat tersebut memicu peringatan SmartScreen bagi pengguna akhir. Rilis publik memerlukan sertifikat dari CA tepercaya.

@end

### Konfigurasi

Edit `build/windows/Taskfile.yml` dan tetapkan variabel penandatanganan:

```yaml
vars:
  SIGN_CERTIFICATE: "path/to/certificate.pfx"
  # Or use thumbprint instead:
  # SIGN_THUMBPRINT: "certificate-thumbprint"
  # TIMESTAMP_SERVER: "http://timestamp.digicert.com"
```

| Variabel | Wajib | Deskripsi |
| --- | --- | --- |
| `SIGN_CERTIFICATE` | Penggantian | Lokasi file sertifikat .pfx/.p12 (menggunakan konfigurasi global jika tidak ditetapkan) |
| `SIGN_THUMBPRINT` | Penggantian | Sidik jari sertifikat dalam penyimpanan sertifikat Windows (alternatif untuk `SIGN_CERTIFICATE`) |
| `TIMESTAMP_SERVER` | Tidak | URL server stempel waktu (bawaan: http://timestamp.digicert.com) |

@note{type="note"}
Variabel-variabel ini adalah **penggantian** opsional. Jika tidak ditetapkan, sertifikat yang dikonfigurasi secara global melalui `wails3 setup` akan digunakan. Lihat [Prioritas konfigurasi](#prioritas-konfigurasi).

@end

@note{type="note"}
Kata sandi sertifikat disimpan dalam gantungan kunci sistem Anda, bukan dalam Taskfile. Jalankan `wails3 setup signing` untuk mengonfigurasinya, atau tetapkan variabel lingkungan `WAILS_WINDOWS_CERT_PASSWORD` di CI.

@end

Kemudian jalankan:

```bash
wails3 task windows:sign           # Build and sign executable
wails3 task windows:sign:installer # Build and sign NSIS installer
```

### Penandatanganan Lintas Platform

File yang dapat dieksekusi Windows dapat ditandatangani dari platform apa pun. Konfigurasi Taskfile dan perintah yang sama dapat digunakan di macOS dan Linux.

### Format Windows yang Didukung

| Format | Ekstensi | Catatan |
| --- | --- | --- |
| File yang dapat dieksekusi | .exe | Penandatanganan PE standar |
| Penginstal | .msi | Paket Windows Installer |
| Paket Aplikasi | .msix, .appx | Aplikasi Windows modern |

## Penandatanganan Paket Linux

Paket Linux (DEB dan RPM) ditandatangani menggunakan kunci PGP/GPG. Tidak seperti penandatanganan kode Windows dan macOS, penandatanganan paket Linux membuktikan bahwa paket berasal dari sumber tepercaya, bukan bahwa kode tersebut dipercaya oleh OS.

### Prasyarat

- Pasangan kunci PGP (dapat dibuat dengan Wails)

### Membuat Kunci PGP

Cara termudah adalah menggunakan wisaya penyiapan—jalankan `wails3 setup`, buka tab **Linux** pada langkah penandatanganan, lalu pilih **Buat kunci GPG baru**. Wisaya tersebut:

- membuat kunci RSA 4096 dalam keyring GPG Anda (biarkan frasa sandi kosong untuk kunci yang dapat digunakan tanpa pengawasan dan cocok untuk CI),
- **mengekspornya ke `~/.wails/signing/<keyid>.asc`** (file yang digunakan build untuk menandatangani), dan
- mencatat ID kunci dan lokasi hasil ekspor dalam `~/.config/wails/defaults.yaml` agar kunci tersebut digunakan secara otomatis saat penandatanganan.

@note{type="note"}
Jika sebuah kunci dikonfigurasi **hanya berdasarkan ID** (misalnya, dibuat sebelum wisaya mengekspor kunci secara otomatis), tab Linux akan mengekspornya ke sebuah file saat Anda membukanya lagi dan mengisi lokasinya—tidak diperlukan langkah manual. Build melakukan penandatanganan dengan *file* kunci, sehingga lokasinyalah yang penting.

@end

Anda juga dapat melakukannya secara manual dengan `gpg`:

```bash
# Interactive — the wizard will prompt for name, email, key size and expiry.
gpg --full-generate-key

# Export the key pair to ASCII-armoured files for the Taskfile to consume.
gpg --armor --export-secret-keys "your@email.com" > signing-key.asc
gpg --armor --export "your@email.com" > signing-key.pub.asc
```

Pengaturan yang disarankan: RSA 4096-bit, masa berlaku 1 tahun, dan dilindungi dengan kata sandi yang kuat.

@note{type="caution"}
Jaga keamanan kunci privat Anda! Simpan dalam bentuk terenkripsi dan buat cadangannya dengan aman.

@end

### Konfigurasi

Edit `build/linux/Taskfile.yml` dan tetapkan variabel penandatanganan:

```yaml
vars:
  PGP_KEY: "path/to/signing-key.asc"
  # SIGN_ROLE: "builder"  # Options: origin, maint, archive, builder
```

| Variabel | Wajib | Deskripsi |
| --- | --- | --- |
| `PGP_KEY` | Penimpaan | Jalur ke file kunci privat PGP yang diekspor (menggunakan kunci yang dikonfigurasi secara global melalui `wails3 setup` jika tidak ditetapkan) |
| `SIGN_ROLE` | Tidak | Peran penandatanganan DEB (default: builder) |

@note{type="note"}
Kata sandi kunci PGP disimpan di keychain sistem Anda, bukan di Taskfile. Jalankan `wails3 setup signing` untuk mengonfigurasinya, atau tetapkan variabel lingkungan `WAILS_PGP_PASSWORD` di CI.

@end

Kemudian jalankan:

```bash
wails3 task linux:sign:deb       # Build and sign DEB package
wails3 task linux:sign:rpm       # Build and sign RPM package
wails3 task linux:sign:packages  # Build and sign all packages
```

### Peran Penandatanganan DEB

Untuk paket DEB, Anda dapat menentukan peran penandatanganan melalui `SIGN_ROLE`:

- `origin`: Tanda tangan dari sumber asal paket
- `maint`: Tanda tangan dari pemelihara paket
- `archive`: Tanda tangan dari pemelihara arsip
- `builder`: Tanda tangan dari pembuat paket (default)

### Penandatanganan Lintas Platform

Paket Linux dapat ditandatangani dari platform apa pun. Konfigurasi dan perintah Taskfile yang sama dapat digunakan di Windows dan macOS.

### Melihat Informasi Kunci

```bash
gpg --show-keys signing-key.asc
```

Keluaran:

```
pub   rsa4096 2024-01-15 [SC] [expires: 2025-01-15]
      1234 5678 90AB CDEF 1234 5678 90AB CDEF 1234 5678
uid                      Your Name <your@email.com>
```

### Memverifikasi Paket Linux

```bash
# Verify DEB signature
dpkg-sig --verify myapp_1.0.0_amd64.deb

# Verify RPM signature
rpm --checksig myapp-1.0.0.x86_64.rpm
```

### Mendistribusikan Kunci Publik Anda

Pengguna memerlukan kunci publik Anda untuk memverifikasi paket:

```bash
# Export public key for distribution
gpg --armor --export "your@email.com" > myapp-signing.pub.asc

# Users can import it:
# For DEB (apt):
sudo apt-key add myapp-signing.pub.asc
# Or for modern apt:
sudo cp myapp-signing.pub.asc /etc/apt/trusted.gpg.d/

# For RPM:
sudo rpm --import myapp-signing.pub.asc
```

## Integrasi GitHub Actions

Di lingkungan CI, kata sandi diberikan melalui variabel lingkungan, bukan keychain sistem:

| Variabel Lingkungan | Deskripsi |
| --- | --- |
| `WAILS_WINDOWS_CERT_PASSWORD` | Kata sandi sertifikat Windows |
| `WAILS_PGP_PASSWORD` | Kata sandi kunci PGP untuk paket Linux |

Anda juga dapat meneruskan variabel Taskfile secara langsung:

```bash
wails3 task darwin:sign SIGN_IDENTITY="$SIGN_IDENTITY" KEYCHAIN_PROFILE="$KEYCHAIN_PROFILE"
```

### Alur Kerja macOS

```yaml
name: Build and Sign macOS

on:
  push:
    tags: ['v*']

jobs:
  build:
    runs-on: macos-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Install Wails
        run: go install github.com/wailsapp/wails/v3/cmd/wails3@latest

      - name: Import Certificate
        env:
          CERTIFICATE_BASE64: ${{ secrets.MACOS_CERTIFICATE }}
          CERTIFICATE_PASSWORD: ${{ secrets.MACOS_CERTIFICATE_PASSWORD }}
        run: |
          echo $CERTIFICATE_BASE64 | base64 --decode > certificate.p12
          security create-keychain -p "" build.keychain
          security default-keychain -s build.keychain
          security unlock-keychain -p "" build.keychain
          security import certificate.p12 -k build.keychain -P "$CERTIFICATE_PASSWORD" -T /usr/bin/codesign
          security set-key-partition-list -S apple-tool:,apple:,codesign: -s -k "" build.keychain

      - name: Store Notarization Credentials
        env:
          APPLE_ID: ${{ secrets.APPLE_ID }}
          APPLE_TEAM_ID: ${{ secrets.APPLE_TEAM_ID }}
          APPLE_APP_PASSWORD: ${{ secrets.APPLE_APP_PASSWORD }}
        run: |
          xcrun notarytool store-credentials "notarize-profile" \
            --apple-id "$APPLE_ID" \
            --team-id "$APPLE_TEAM_ID" \
            --password "$APPLE_APP_PASSWORD"

      - name: Build, Sign, and Notarize
        env:
          SIGN_IDENTITY: ${{ secrets.MACOS_SIGN_IDENTITY }}
        run: |
          wails3 task darwin:sign:notarize \
            SIGN_IDENTITY="$SIGN_IDENTITY" \
            KEYCHAIN_PROFILE="notarize-profile"

      - name: Upload Artifact
        uses: actions/upload-artifact@v4
        with:
          name: MyApp-macOS
          path: bin/*.app
```

### Alur Kerja Windows

```yaml
name: Build and Sign Windows

on:
  push:
    tags: ['v*']

jobs:
  build:
    runs-on: windows-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Install Wails
        run: go install github.com/wailsapp/wails/v3/cmd/wails3@latest

      - name: Import Certificate
        env:
          CERTIFICATE_BASE64: ${{ secrets.WINDOWS_CERTIFICATE }}
        run: |
          $certBytes = [Convert]::FromBase64String($env:CERTIFICATE_BASE64)
          [IO.File]::WriteAllBytes("certificate.pfx", $certBytes)

      - name: Build and Sign
        env:
          WAILS_WINDOWS_CERT_PASSWORD: ${{ secrets.WINDOWS_CERTIFICATE_PASSWORD }}
        run: |
          wails3 task windows:sign SIGN_CERTIFICATE=certificate.pfx

      - name: Upload Artifact
        uses: actions/upload-artifact@v4
        with:
          name: MyApp-Windows
          path: bin/*.exe
```

### Alur Kerja Lintas Platform (Runner Linux)

Tandatangani paket Windows dan Linux dari satu runner Linux:

```yaml
name: Build and Sign (Cross-Platform)

on:
  push:
    tags: ['v*']

jobs:
  build-and-sign:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Install Wails
        run: go install github.com/wailsapp/wails/v3/cmd/wails3@latest

      - name: Install Build Dependencies
        run: |
          sudo apt-get update
          sudo apt-get install -y nsis rpm

      # Import certificates
      - name: Import Certificates
        env:
          WINDOWS_CERT_BASE64: ${{ secrets.WINDOWS_CERTIFICATE }}
          PGP_KEY_BASE64: ${{ secrets.PGP_PRIVATE_KEY }}
        run: |
          echo "$WINDOWS_CERT_BASE64" | base64 -d > certificate.pfx
          echo "$PGP_KEY_BASE64" | base64 -d > signing-key.asc

      # Build and sign Windows
      - name: Build and Sign Windows
        env:
          WAILS_WINDOWS_CERT_PASSWORD: ${{ secrets.WINDOWS_CERTIFICATE_PASSWORD }}
        run: |
          wails3 task windows:sign SIGN_CERTIFICATE=certificate.pfx

      # Build and sign Linux packages
      - name: Build and Sign Linux Packages
        env:
          WAILS_PGP_PASSWORD: ${{ secrets.PGP_PASSWORD }}
        run: |
          wails3 task linux:sign:packages PGP_KEY=signing-key.asc

      # Cleanup secrets
      - name: Cleanup
        if: always()
        run: rm -f certificate.pfx signing-key.asc

      - name: Upload Artifacts
        uses: actions/upload-artifact@v4
        with:
          name: signed-binaries
          path: |
            bin/*.exe
            bin/*.deb
            bin/*.rpm
```

@note{type="note"}
Menggunakan runner Linux untuk penandatanganan lintas platform menyederhanakan CI/CD karena tidak memerlukan runner Windows terpisah. Penandatanganan macOS tetap memerlukan runner macOS karena persyaratan alat native Apple.

@end

## Referensi CLI

### wails3 setup signing

Wizard interaktif untuk mengonfigurasi penandatanganan proyek Anda.

```bash
wails3 setup signing [flags]

Flags:
  --platform    Platform to configure (darwin, windows, linux). Repeatable.
                If omitted, auto-detects which platforms to configure from the build directory.
```

Wizard akan memandu Anda untuk:

- **macOS**: Memilih sertifikat Developer ID dan mengonfigurasi kredensial notarisasi (memanggil `xcrun notarytool store-credentials`).
- **Windows**: Memilih antara file sertifikat atau thumbprint, serta menetapkan kata sandi dan server stempel waktu.
- **Linux**: Menggunakan kunci PGP yang ada atau membuat kunci baru (memanggil `gpg`), serta mengonfigurasi peran penandatanganan.

### wails3 setup entitlements

Wizard interaktif untuk mengonfigurasi entitlement macOS.

```bash
wails3 setup entitlements [flags]

Flags:
  --output    Output path for entitlements.plist (default: build/darwin/entitlements.plist)
```

**Preset:**

- **Pengembangan**: Membuat `entitlements.dev.plist` dengan JIT, debugging, dan jaringan
- **Produksi**: Membuat `entitlements.plist` dengan entitlement minimal
- **Keduanya**: Membuat kedua file (disarankan)
- **App Store**: Membuat entitlement dengan sandbox untuk Mac App Store
- **Kustom**: Pilih setiap entitlement dan file target

### wails3 sign

Tandatangani biner dan paket untuk platform saat ini atau platform yang ditentukan. Perintah ini merupakan wrapper yang memanggil tugas penandatanganan khusus platform yang sesuai.

```bash
wails3 sign
wails3 sign GOOS=darwin
wails3 sign GOOS=windows
wails3 sign GOOS=linux
```

Perintah ini menjalankan tugas `<platform>:sign` yang sesuai, yang menggunakan konfigurasi penandatanganan dari Taskfile Anda.

### wails3 tool sign

Perintah tingkat rendah untuk menandatangani file tertentu secara langsung. Digunakan secara internal oleh Taskfile.

```bash
wails3 tool sign [flags]
```

**Flag Umum:**\

| Flag | Deskripsi |
| --- | --- |
| `--input` | Jalur ke file yang akan ditandatangani |
| `--output` | Jalur keluaran (opsional, secara default menimpa di lokasi yang sama) |
| `--verbose` | Aktifkan keluaran mendetail |

**Flag Windows/macOS:**\

| Flag | Deskripsi |
| --- | --- |
| `--certificate` | Jalur ke sertifikat PKCS#12 (.pfx/.p12) |
| `--password` | Kata sandi sertifikat |
| `--timestamp` | URL server stempel waktu |

**Flag Khusus macOS:**\

| Flag | Deskripsi |
| --- | --- |
| `--identity` | Identitas penandatanganan (gunakan '-' untuk ad hoc) |
| `--entitlements` | Jalur ke plist hak akses |
| `--hardened-runtime` | Aktifkan runtime yang diperkuat (default: true) |
| `--notarize` | Kirim untuk notarisasi |
| `--keychain-profile` | Profil Keychain untuk notarisasi |

**Flag Khusus Windows:**\

| Flag | Deskripsi |
| --- | --- |
| `--thumbprint` | Thumbprint sertifikat di penyimpanan Windows |

**Flag Khusus Linux:**\

| Flag | Deskripsi |
| --- | --- |
| `--pgp-key` | Jalur ke kunci privat PGP |
| `--pgp-password` | Kata sandi kunci PGP |
| `--role` | Peran penandatanganan DEB (origin/maint/archive/builder) |

### Memeriksa status penandatanganan (alat native)

Di v3, **tidak ada** perintah `wails3 signing`. Untuk memeriksa status penandatanganan, gunakan alat native secara langsung:

| Tugas | Perintah |
| --- | --- |
| Cantumkan identitas penandatanganan kode macOS | `security find-identity -v -p codesigning` |
| Simpan kredensial notarisasi | `xcrun notarytool store-credentials "<profile>" --apple-id … --team-id … --password …` |
| Periksa file kunci PGP | `gpg --show-keys <key.asc>` |
| Buat pasangan kunci PGP | `gpg --full-generate-key` |
| Ekspor kunci publik | `gpg --armor --export <email>` |

## Pemecahan Masalah

### Masalah macOS

**"Sertifikat Developer ID tidak ditemukan"**

- Pastikan sertifikat Anda terpasang di Keychain
- Pastikan sertifikat belum kedaluwarsa dengan `security find-identity -v -p codesigning`
- Pastikan Anda memiliki sertifikat "Developer ID Application" (bukan hanya "Apple Development")

**"Notarisasi gagal"**

- Periksa log notarisasi: `xcrun notarytool log <submission-id> --keychain-profile <profile>`
- Pastikan runtime yang diperkuat telah diaktifkan
- Pastikan aplikasi Anda tidak menyertakan biner yang belum ditandatangani

**"Codesign gagal"**

- Pastikan keychain tidak terkunci: `security unlock-keychain`
- Periksa izin file pada bundel aplikasi

### Masalah Windows

**"Sertifikat tidak ditemukan"**

- Pastikan jalur sertifikat sudah benar
- Periksa kata sandi sertifikat
- Pastikan sertifikat masih valid (belum kedaluwarsa atau dicabut)

**"Kesalahan server stempel waktu"**

- Coba server stempel waktu lain:
  - `http://timestamp.digicert.com`
  - `http://timestamp.sectigo.com`
  - `http://timestamp.comodoca.com`


### Masalah Linux

**"Kunci PGP tidak valid"**

- Pastikan file kunci menggunakan format ASCII-armored
- Periksa dengan `gpg --show-keys <key.asc>` bahwa kunci belum kedaluwarsa
- Pastikan kata sandi sudah benar

**"Verifikasi tanda tangan gagal"**

- Pastikan kunci publik telah diimpor dengan benar
- Pastikan paket tidak diubah setelah ditandatangani

## Sumber Daya Tambahan

### Dokumentasi Resmi

- [Panduan Penandatanganan Kode Apple](https://developer.apple.com/support/code-signing/)
- [Dokumentasi Notarisasi Apple](https://developer.apple.com/documentation/security/notarizing_macos_software_before_distribution)
- [Penandatanganan Kode Microsoft](https://docs.microsoft.com/en-us/windows-hardware/drivers/dashboard/get-a-code-signing-certificate)
- [Penandatanganan Paket Debian](https://wiki.debian.org/SecureApt)
- [Penandatanganan Paket RPM](https://rpm-software-management.github.io/rpm/manual/signatures.html)
