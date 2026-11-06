---
title: "Penandatanganan Kode"
description: "Panduan untuk menandatangani aplikasi Wails Anda di semua platform"
slug: "guides/build/signing"
sourcePath: "guides/build/signing.md"
---

## Menandatangani Kode Aplikasi Anda

Panduan ini menjelaskan cara menandatangani aplikasi Wails Anda untuk macOS, Windows, dan Linux. Wails v3 menyediakan alat CLI bawaan untuk penandatanganan kode, notarisasi, dan pengelolaan kunci PGP.

- **macOS** - Tandatangani dan notarisasikan aplikasi macOS Anda
- **Windows** - Tandatangani berkas executable dan paket Windows Anda
- **Linux** - Tandatangani paket DEB dan RPM dengan kunci PGP

## Matriks Penandatanganan Lintas Platform

Matriks ini menunjukkan format yang dapat Anda tandatangani dari setiap platform sumber:

| Format Target | Dari Windows | Dari macOS | Dari Linux |
| --- | :---: | :---: | :---: |
| Windows EXE/DLL/MSIX/APPX, PFX yang diekspor | ✅ | ✅ | ✅ |
| Windows MSI atau identitas penyimpanan sertifikat/penyedia | ✅ | ❌ | ❌ |
| macOS .app / Mach-O / DMG, identitas yang diekspor | ✅ | ✅ | ✅ |
| Notarisasi macOS, kredensial API tim | ✅ | ✅ | ✅ |
| Identitas keychain atau profil notaris macOS | ❌ | ✅ | ❌ |
| Linux DEB | ✅ | ✅ | ✅ |
| Linux RPM | ✅ | ✅ | ✅ |

@note{type="tip"}
Identitas Apple yang diekspor dan kredensial API tim memungkinkan penandatanganan serta notarisasi portabel. Identitas keychain dan profil notaris macOS memerlukan Mac. Pembatasan khusus setiap format dijelaskan di bawah ini.

@end

### Backend Penandatanganan

Wails secara otomatis memilih backend penandatanganan terbaik yang tersedia:

| Platform | Backend Native | Backend Lintas Platform |
| --- | --- | --- |
| Windows | `signtool.exe` untuk identitas penyimpanan/penyedia, MSI, atau PFX RC2 lama | Bawaan untuk PFX RSA yang diekspor dan PE/MSIX |
| macOS | `codesign` untuk identitas keychain | Bawaan untuk identitas PKCS#12 yang diekspor atau penandatanganan ad hoc |
| Linux | Tidak tersedia | Bawaan |

Jenis kredensial menentukan backend penandatanganan Apple pada setiap host. Identitas yang diekspor menggunakan Go; identitas keychain macOS menggunakan `codesign`. Pengambilan, verifikasi, dan penyematan tiket menggunakan Go.

## Mulai Cepat

Cara tercepat untuk mengonfigurasi penandatanganan adalah menggunakan wizard penyiapan:

```bash
wails3 setup
```

Tindakan ini menulis konfigurasi penandatanganan bersama ke `~/.config/wails/defaults.yaml` yang **digunakan saat penandatanganan di setiap platform** (lihat [Prioritas konfigurasi](#prioritas-konfigurasi)). Langkah penandatanganannya:

- Menampilkan alat penyiapan untuk membuat atau mengekspor sertifikat dan kunci, seperti OpenSSL dan GnuPG. Identitas Apple dan Windows yang diekspor menggunakan penanda tangan bawaan; identitas penyimpanan sertifikat/penyedia Windows dan paket MSI memerlukan Windows SignTool.
- Di macOS, menampilkan daftar sertifikat Developer ID dari keychain Anda.
- Untuk Linux, menampilkan daftar kunci GPG Anda dan dapat **membuat serta mengekspor** kunci baru.
- Untuk Windows, dapat **membuat sertifikat yang ditandatangani sendiri** (untuk pengujian) melalui OpenSSL.
- **Menyimpan kata sandi dengan aman di keychain sistem Anda** (bukan di Taskfile).

@note{type="tip"}
Kata sandi disimpan di penyimpanan kredensial native sistem Anda (macOS Keychain, Windows Credential Manager, atau Linux Secret Service). Dengan demikian, konfigurasi penandatanganan Anda tetap aman dan dapat digunakan di semua proyek Wails Anda.

@end

## Prioritas konfigurasi

Saat Anda menjalankan tugas penandatanganan, setiap opsi penandatanganan ditentukan menurut urutan berikut (kecocokan pertama yang digunakan):

1. Flag eksplisit yang diteruskan ke `wails3 tool sign` (misalnya `--pgp-key`, `--certificate`, `--identity`).
2. **Variabel Taskfile proyek** yang sesuai (`PGP_KEY`, `SIGN_CERTIFICATE`/`SIGN_THUMBPRINT`, `SIGN_IDENTITY`, …).
3. **Konfigurasi global** di `~/.config/wails/defaults.yaml` (ditulis oleh `wails3 setup`).

Artinya, variabel Taskfile merupakan **penggantian opsional**: jika suatu variabel tidak ditetapkan, kunci/sertifikat/identitas yang dikonfigurasi secara global akan digunakan. Jika tidak satu pun dari ketiganya memberikan nilai, perintah penandatanganan akan melaporkan galat yang jelas dan memberi tahu Anda cara mengonfigurasinya.

## Konfigurasi per proyek

Untuk mengonfigurasi penandatanganan bagi satu proyek alih-alih secara global, jalankan wizard per proyek dari dalam proyek tersebut—wizard ini menulis `vars` ke berkas `build/<platform>/Taskfile.yml` proyek itu:

```bash
wails3 signing                                   # all detected platforms
wails3 signing --platform windows --platform linux
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

Kata sandi diambil dari keychain sistem (jalankan `wails3 signing` untuk mengonfigurasinya).

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

Kata sandi diambil dari keychain sistem (jalankan `wails3 signing` untuk mengonfigurasinya).

Kemudian jalankan:

```bash
wails3 task linux:sign:deb       # Sign DEB package
wails3 task linux:sign:rpm       # Sign RPM package
wails3 task linux:sign:packages  # Sign all packages
```

@end

Anda juga dapat memeriksa status penandatanganan langsung dari sistem:

```bash
# List available macOS code-signing identities
security find-identity -v -p codesigning

# List PGP keys (Linux package signing)
gpg --list-keys
```

Untuk menyiapkan konfigurasi penandatanganan setiap platform secara interaktif, gunakan wizard:

```bash
wails3 signing
```

## Penandatanganan Kode macOS

### Prasyarat

- Akun Apple Developer ($99/tahun)
- Sertifikat Developer ID Application
- Untuk identitas keychain macOS: Xcode Command Line Tools telah terinstal
- Untuk penandatanganan portabel: identitas RSA PKCS#12 yang diekspor; tidak memerlukan alat penandatanganan Apple

### Penandatanganan portabel dengan Wake

Wake membaca `wails.yaml` secara langsung. Wake tidak menjalankan Taskfile. Konfigurasikan
identitas Developer ID yang diekspor dan kunci API tim App Store Connect:

```yaml
darwin:
  signing:
    certificate: signing/developer-id.p12
    credential: APPLE_CERT_PASSWORD
    timestamp_server: https://timestamp.digicert.com
  notarization:
    api_key: signing/AuthKey.p8
    key_id: ABCDEFGHIJ
    issuer_id: 00000000-0000-0000-0000-000000000000
```

`credential` menamai variabel lingkungan yang menyimpan kata sandi sertifikat.
Simpan kedua berkas kunci privat di luar kontrol versi. Aktifkan penandatanganan dan
notarisasi dalam profil rilis yang dijelaskan di
[panduan build YAML](/experimental/yaml-builds/).

Penanda tangan portabel mendukung berkas Mach-O Intel 64-bit dan Apple Silicon, biner universal, bundel aplikasi kanonis dengan kode bertingkat, framework, serta DMG UDIF yang berukuran kurang dari 4 GiB. Tata letak yang tidak didukung akan ditolak sebelum output dipublikasikan. Penandatanganan ad-hoc menggunakan `identity: "-"` dan tidak dapat dinotarisasi. Identitas PKCS#12 dapat menggunakan enkripsi AES/PBKDF2 atau 3DES lama; ekspor RC2 harus diekspor ulang dengan cipher yang didukung.

Tanda tangan portabel yang dinotarisasi memerlukan layanan stempel waktu SHA-2 RFC 3161 yang ditentukan secara eksplisit. Tanda tangan stempel waktu SHA-1 lama ditolak. Pengiriman melalui API tim mendukung satu unggahan hingga 5 GiB; kunci API individual dan unggahan multipart tidak didukung. Verifikasi dan penyematan tiket menggunakan protokol layanan Apple yang teramati dengan Apple Root G3 yang disematkan, tanpa `xcrun stapler`. Gatekeeper tetap menjadi otoritas untuk kepercayaan distribusi dan pencabutan.

Bagian berikut menjelaskan alur kerja keychain macOS yang ada dan digunakan oleh proyek lama serta konfigurasi Wake yang memilih identitas keychain.

### Identitas Penandatanganan

Periksa identitas penandatanganan yang tersedia:

```bash
security find-identity -v -p codesigning
```

Output:

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

Edit `build/darwin/Taskfile.yml` dan atur variabel penandatanganan:

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  KEYCHAIN_PROFILE: "my-notarize-profile"
  ENTITLEMENTS: "build/darwin/entitlements.plist"
```

| Variabel | Wajib | Deskripsi |
| --- | --- | --- |
| `SIGN_IDENTITY` | Ya | Developer ID Anda (misalnya, "Developer ID Application: Perusahaan Anda (TEAMID)") |
| `KEYCHAIN_PROFILE` | Untuk notarisasi | Nama profil keychain yang berisi kredensial tersimpan |
| `ENTITLEMENTS` | Tidak | Jalur ke berkas entitlement |

Kemudian jalankan:

```bash
wails3 task darwin:sign           # Build, package, and sign
wails3 task darwin:sign:notarize  # Build, package, sign, and notarize
```

### Entitlement

Entitlement mengatur kapabilitas yang dapat diakses aplikasi Anda. Aplikasi Wails biasanya memerlukan entitlement yang berbeda untuk pengembangan dan produksi:

- **Pengembangan**: Memerlukan entitlement untuk JIT, memori tanpa tanda tangan, dan debugging
- **Produksi**: Entitlement minimal (hanya akses jaringan)

Gunakan wizard penyiapan interaktif untuk menghasilkan kedua berkas:

```bash
wails3 entitlements
```

Perintah ini membuat:

- `build/darwin/entitlements.dev.plist` - Untuk build pengembangan
- `build/darwin/entitlements.plist` - Untuk build produksi/bertanda tangan

**Preset yang tersedia:**\

| Preset | Deskripsi |
| --- | --- |
| Pengembangan | JIT, memori tanpa tanda tangan, debugging, jaringan |
| Produksi | Hanya jaringan (minimal, paling aman) |
| Keduanya | Membuat berkas pengembangan dan produksi (disarankan) |
| App Store | Sandbox diaktifkan dengan akses jaringan dan berkas |
| Kustom | Pilih setiap entitlement secara individual |

@note{type="note"}
Task `run` dalam Taskfile darwin menggunakan `entitlements.dev.plist` secara otomatis. Task `sign` menggunakan `entitlements.plist` untuk build produksi.

@end

Kemudian atur `ENTITLEMENTS` dalam variabel Taskfile Anda agar mengarah ke berkas yang sesuai.

### Notarisasi

Apple mewajibkan semua aplikasi yang didistribusikan untuk dinotarisasi.

@steps
### **Simpan kredensial Anda di keychain** (penyiapan satu kali). Jalankan `wails3 signing` (perintah ini meminta nilai tersebut dan memanggil `notarytool` di balik layar), atau panggil `notarytool` secara langsung:
```bash
xcrun notarytool store-credentials "my-notarize-profile" \
  --apple-id "your@email.com" \
  --team-id "ABCD1234" \
  --password "app-specific-password"
```

### **Atur KEYCHAIN_PROFILE dalam Taskfile Anda** agar sesuai dengan nama profil di atas.
### **Tandatangani dan notariskan aplikasi Anda**:
```bash
wails3 task darwin:sign:notarize
```

### **Verifikasi notarisasi**:
```bash
spctl --assess --verbose=2 bin/MyApp.app
```

@end

@note{type="note"}
Notarisasi biasanya memerlukan waktu 1-2 menit. Tiket secara otomatis disematkan ke aplikasi Anda.

@end

## Penandatanganan Kode Windows

### Prasyarat

- Sertifikat penandatanganan kode (dari DigiCert, Sectigo, dan sebagainya)
- Untuk penandatanganan portabel: identitas RSA PFX yang diekspor dan dienkripsi dengan AES/PBKDF2 atau 3DES lama; tidak memerlukan alat penandatanganan eksternal
- Untuk identitas penyimpanan sertifikat/penyedia, paket MSI, atau PFX RC2 lama: Windows dengan alat penandatanganan Windows SDK di `PATH`

### Membuat sertifikat yang ditandatangani sendiri (pengujian)

Jika Anda hanya perlu menguji pipeline penandatanganan, jalankan `wails3 setup`, buka tab **Windows** pada langkah penandatanganan, lalu pilih **Buat sertifikat yang ditandatangani sendiri**. Tindakan ini menggunakan OpenSSL untuk membuat `.pfx` penandatanganan kode dan mencatat jalurnya dalam konfigurasi global.

@note{type="caution"}
Sertifikat yang ditandatangani sendiri hanya ditujukan untuk **pengujian dan distribusi internal** — sertifikat tersebut memicu peringatan SmartScreen bagi pengguna akhir. Rilis publik memerlukan sertifikat dari CA tepercaya.

@end

### Konfigurasi

Edit `build/windows/Taskfile.yml` dan atur variabel penandatanganan:

```yaml
vars:
  SIGN_CERTIFICATE: "path/to/certificate.pfx"
  # Or use thumbprint instead:
  # SIGN_THUMBPRINT: "certificate-thumbprint"
  # TIMESTAMP_SERVER: "http://timestamp.digicert.com"
```

| Variabel | Wajib | Deskripsi |
| --- | --- | --- |
| `SIGN_CERTIFICATE` | Penggantian | Jalur ke berkas sertifikat .pfx/.p12 (menggunakan konfigurasi global sebagai fallback jika tidak diatur) |
| `SIGN_THUMBPRINT` | Penggantian | Thumbprint sertifikat dalam penyimpanan sertifikat Windows (alternatif untuk `SIGN_CERTIFICATE`) |
| `TIMESTAMP_SERVER` | Tidak | URL server stempel waktu (default: http://timestamp.digicert.com) |

@note{type="note"}
Variabel ini merupakan **penggantian** opsional. Jika tidak diatur, sertifikat yang dikonfigurasi secara global melalui `wails3 setup` akan digunakan. Lihat [Prioritas konfigurasi](#prioritas-konfigurasi).

@end

@note{type="note"}
Kata sandi sertifikat disimpan dalam keychain sistem Anda, bukan dalam Taskfile. Jalankan `wails3 signing` untuk mengonfigurasinya, atau atur variabel lingkungan `WAILS_WINDOWS_CERT_PASSWORD` di CI.

@end

Kemudian jalankan:

```bash
wails3 task windows:sign           # Build and sign executable
wails3 task windows:sign:installer # Build and sign NSIS installer
```

### Penandatanganan Lintas Platform

Identitas RSA PFX yang diekspor menandatangani file EXE/DLL, berkas executable penginstal NSIS, serta paket MSIX/APPX individual menggunakan Go di Windows, macOS, dan Linux. Tanda tangan yang sudah ada dapat diganti. Proses penandatanganan menempatkan output di area staging dan memublikasikannya hanya setelah verifikasi tanda tangan dan stempel waktu berhasil, sehingga file asli tetap dipertahankan jika terjadi kegagalan. Gunakan `--output` untuk menulis file bertanda tangan yang terpisah.

Wake membaca pengaturan ini langsung dari `wails.yaml`:

```yaml
windows:
  signing:
    certificate: signing/windows.pfx
    credential: WINDOWS_SIGNING_PASSWORD
    timestamp_server: http://timestamp.digicert.com
```

`credential` menentukan nama variabel lingkungan untuk kata sandi, bukan kata sandinya sendiri. Aktifkan `sign: true` untuk target rilis; lihat [panduan build YAML](/experimental/yaml-builds/). Wake tidak menjalankan Taskfile. Proyek lama dapat terus menggunakan konfigurasi Taskfile di atas.

Profil portabel menggunakan RSA/SHA-256 dan memvalidasi stempel waktu RFC 3161, termasuk rantai sertifikat tepercaya milik penanda tangan stempel waktu. Layanan default-nya adalah DigiCert. MSIX Publisher harus cocok dengan subjek sertifikat penandatanganan; bundel paket dan arsip berukuran 4 GiB atau lebih tidak didukung. Ukuran file PFX dibatasi hingga 16 MiB. Algoritma sertifikat lain dan enkripsi RC2 ditolak. Ekspor ulang PFX RC2 menggunakan AES atau 3DES untuk penandatanganan portabel; di Windows, PFX RC2 yang berhasil diautentikasi menggunakan SignTool sebagai gantinya. Kata sandi yang salah atau PFX yang rusak tidak pernah memicu mekanisme fallback tersebut.

Kunci dari penyimpanan sertifikat, token perangkat keras, dan penyedia tetap menggunakan alur kerja SignTool native Windows melalui `thumbprint`. Penandatanganan MSI juga memerlukan SignTool di Windows. Mengonfigurasi sertifikat yang diekspor tidak menjadikan alur kerja ini portabel.

### Format Windows yang Didukung

| Format | Ekstensi | Catatan |
| --- | --- | --- |
| File executable dan pustaka | .exe, .dll | Penandatanganan PE portabel, termasuk penginstal NSIS |
| Penginstal | .msi | Memerlukan Windows SignTool |
| Paket aplikasi | .msix, .appx | Penandatanganan paket individual secara portabel; bundel tidak didukung |

## Penandatanganan Paket Linux

Paket Linux (DEB dan RPM) ditandatangani menggunakan kunci PGP/GPG. Berbeda dari penandatanganan kode Windows dan macOS, penandatanganan paket Linux membuktikan bahwa paket berasal dari sumber tepercaya, bukan bahwa kode tersebut dipercaya oleh sistem operasi.

File kunci privat armored atau biner yang diekspor menggunakan penanda tangan native Go di Windows, macOS, dan Linux. Penandatanganan DEB atau RPM yang sudah ada mempertahankan isi paket tanpa membongkar atau mengompresinya ulang dan tidak memerlukan instalasi GnuPG, `dpkg-sig`, ataupun `rpmsign`. Paket Arch menggunakan pemuat kunci yang sama untuk menulis tanda tangan terpisah. GnuPG tetap berguna untuk membuat atau mengekspor kunci serta memverifikasi tanda tangan secara independen. Konfigurasi DEB/RPM lama yang menentukan identitas keyring, bukan file kunci, tetap menggunakan alat penandatanganan eksternalnya.

### Prasyarat

- Pasangan kunci PGP (dapat dibuat dengan Wails)

### Membuat Kunci PGP

Cara termudah adalah menggunakan wizard penyiapan—jalankan `wails3 setup`, buka tab **Linux** pada langkah penandatanganan, lalu pilih **Buat kunci GPG baru**. Wizard tersebut:

- membuat kunci RSA 4096 di keyring GPG Anda (biarkan frasa sandi kosong untuk kunci yang dapat digunakan tanpa interaksi dan cocok untuk CI),
- **mengekspornya ke `~/.wails/signing/<keyid>.asc`** (file yang digunakan build untuk menandatangani), dan
- mencatat ID kunci dan jalur hasil ekspor di `~/.config/wails/defaults.yaml` agar digunakan secara otomatis saat penandatanganan.

@note{type="note"}
Jika kunci dikonfigurasi **hanya berdasarkan ID** (misalnya dibuat sebelum wizard mengekspor kunci secara otomatis), tab Linux akan mengekspornya ke sebuah file saat Anda membukanya lagi dan mengisi jalurnya—tanpa memerlukan langkah manual. Build melakukan penandatanganan dengan *file* kunci, sehingga jalur file-lah yang penting.

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
Jaga keamanan kunci privat Anda! Simpan dalam keadaan terenkripsi dan buat cadangan dengan aman.

@end

### Konfigurasi

Edit `build/linux/Taskfile.yml` dan atur variabel penandatanganan:

```yaml
vars:
  PGP_KEY: "path/to/signing-key.asc"
  # SIGN_ROLE: "builder"  # Options: origin, maint, archive, builder
```

| Variabel | Wajib | Deskripsi |
| --- | --- | --- |
| `PGP_KEY` | Timpa | Jalur ke file kunci privat PGP yang diekspor (jika tidak diatur, menggunakan kunci yang dikonfigurasi secara global melalui `wails3 setup`) |
| `SIGN_ROLE` | Tidak | Peran penandatanganan DEB (default: builder) |

@note{type="note"}
Kata sandi kunci PGP disimpan di keychain sistem Anda, bukan di Taskfile. Jalankan `wails3 signing` untuk mengonfigurasinya, atau atur variabel lingkungan `WAILS_PGP_PASSWORD` di CI.

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
- `maint`: Tanda tangan dari pengelola paket
- `archive`: Tanda tangan dari pengelola arsip
- `builder`: Tanda tangan dari pembuat paket (default)

### Penandatanganan Lintas Platform

Paket Linux dapat ditandatangani dari platform apa pun. Konfigurasi dan perintah Taskfile yang sama dapat digunakan di Windows dan macOS.

### Melihat Informasi Kunci

```bash
gpg --show-keys signing-key.asc
```

Output:

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

Di lingkungan CI, kata sandi diberikan melalui variabel lingkungan, bukan melalui keychain sistem:

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
        run: go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0

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
        run: go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0

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
        run: go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0

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
Menggunakan runner Linux untuk penandatanganan lintas platform menyederhanakan CI/CD karena tidak memerlukan runner Windows terpisah. Identitas keychain Apple dan profil notaris memerlukan runner macOS. Identitas yang diekspor dan kredensial API tim mendukung alur kerja Wake portabel di atas.

@end

## Referensi CLI

### wails3 signing

Wizard interaktif untuk mengonfigurasi penandatanganan bagi proyek Anda.

```bash
wails3 signing [flags]

Flags:
  --platform    Platform to configure (darwin, windows, linux). Repeatable.
                If omitted, auto-detects which platforms to configure from the build directory.
```

Wizard akan memandu Anda melalui:

- **macOS**: Memilih sertifikat Developer ID dan mengonfigurasi kredensial notarisasi (memanggil `xcrun notarytool store-credentials`).
- **Windows**: Memilih antara berkas sertifikat atau thumbprint, serta menetapkan kata sandi dan server stempel waktu.
- **Linux**: Menggunakan kunci PGP yang sudah ada atau membuat kunci baru (memanggil `gpg`), serta mengonfigurasi peran penandatanganan.

### wails3 entitlements

Wizard interaktif untuk mengonfigurasi entitlement macOS.

```bash
wails3 entitlements [flags]

Flags:
  --output    Output path for entitlements.plist (default: build/darwin/entitlements.plist)
```

**Preset:**

- **Pengembangan**: Membuat `entitlements.dev.plist` dengan JIT, debugging, dan jaringan
- **Produksi**: Membuat `entitlements.plist` dengan entitlement minimal
- **Keduanya**: Membuat kedua berkas (disarankan)
- **App Store**: Membuat entitlement dengan sandbox untuk Mac App Store
- **Kustom**: Pilih setiap entitlement dan berkas target

### wails3 sign

Tandatangani biner dan paket untuk platform saat ini atau platform yang ditentukan. Perintah ini adalah pembungkus yang memanggil tugas penandatanganan khusus platform yang sesuai.

```bash
wails3 sign
wails3 sign GOOS=darwin
wails3 sign GOOS=windows
wails3 sign GOOS=linux
```

Perintah ini menjalankan tugas `<platform>:sign` yang sesuai, yang menggunakan konfigurasi penandatanganan dari Taskfile Anda.

### wails3 tool sign

Perintah tingkat rendah untuk menandatangani berkas tertentu secara langsung. Digunakan secara internal oleh Taskfile.

```bash
wails3 tool sign [flags]
```

**Flag Umum:**\

| Flag | Deskripsi |
| --- | --- |
| `--input` | Path ke berkas yang akan ditandatangani |
| `--output` | Path keluaran (opsional, secara default menimpa berkas di tempat) |
| `--verbose` | Aktifkan keluaran mendetail |

**Flag Windows/macOS:**\

| Flag | Deskripsi |
| --- | --- |
| `--certificate` | Path ke sertifikat PKCS#12 (.pfx/.p12) |
| `--password` | Kata sandi sertifikat |
| `--timestamp` | URL server stempel waktu |

**Flag Khusus macOS:**\

| Flag | Deskripsi |
| --- | --- |
| `--identity` | Identitas penandatanganan (gunakan '-' untuk ad hoc) |
| `--entitlements` | Path ke plist entitlement |
| `--hardened-runtime` | Aktifkan hardened runtime (default: true) |
| `--notarize` | Kirim untuk notarisasi |
| `--keychain-profile` | Profil Keychain untuk notarisasi |
| `--notary-key` | Berkas kunci privat `.p8` API tim App Store Connect; menggunakan klien Go Notary API |
| `--notary-key-id` | ID kunci API tim; wajib digunakan bersama `--notary-key` |
| `--notary-issuer-id` | ID penerbit tim; wajib digunakan bersama `--notary-key` |

**Flag Khusus Windows:**\

| Flag | Deskripsi |
| --- | --- |
| `--thumbprint` | Thumbprint sertifikat di penyimpanan Windows |

**Flag Khusus Linux:**\

| Flag | Deskripsi |
| --- | --- |
| `--pgp-key` | Path ke kunci privat PGP |
| `--pgp-password` | Kata sandi kunci PGP |
| `--role` | Peran penandatanganan DEB (origin/maint/archive/builder) |

### Memeriksa status penandatanganan (alat native)

`wails3 signing` mengonfigurasi kredensial penandatanganan. Untuk memeriksa status penandatanganan, gunakan alat native secara langsung:

| Tugas | Perintah |
| --- | --- |
| Cantumkan identitas penandatanganan kode macOS | `security find-identity -v -p codesigning` |
| Simpan kredensial notarisasi | `xcrun notarytool store-credentials "<profile>" --apple-id … --team-id … --password …` |
| Periksa berkas kunci PGP | `gpg --show-keys <key.asc>` |
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
- Pastikan hardened runtime diaktifkan
- Pastikan aplikasi Anda tidak menyertakan biner yang belum ditandatangani

**"Codesign gagal"**

- Pastikan Keychain tidak terkunci: `security unlock-keychain`
- Periksa izin berkas pada bundel aplikasi

### Masalah Windows

**"Sertifikat tidak ditemukan"**

- Pastikan path sertifikat sudah benar
- Periksa kata sandi sertifikat
- Pastikan sertifikat valid (belum kedaluwarsa atau dicabut)

**"Kesalahan server stempel waktu"**

- Coba server stempel waktu lain:
  - `http://timestamp.digicert.com`
  - `http://timestamp.sectigo.com`
  - `http://timestamp.comodoca.com`


### Masalah Linux

**"Kunci PGP tidak valid"**

- Pastikan berkas kunci menggunakan format ASCII-armored
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
