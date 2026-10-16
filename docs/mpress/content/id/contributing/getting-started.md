---
title: "Memulai"
description: "Cara mulai berkontribusi pada Wails v3"
slug: "contributing/getting-started"
sourcePath: "contributing/getting-started.md"
---

## Selamat Datang, Kontributor!

Terima kasih atas minat Anda untuk berkontribusi pada Wails! Panduan ini akan membantu Anda memberikan kontribusi pertama.

## Prasyarat

Sebelum memulai, pastikan Anda memiliki:

- **Go 1.25+** yang telah terinstal ([unduh](https://go.dev/dl/))
- **Node.js 20+** dan **npm** ([unduh](https://nodejs.org/))
- **Git** yang telah dikonfigurasi dengan akun GitHub Anda
- Pemahaman dasar tentang Go dan JavaScript/TypeScript

### Persyaratan Khusus Platform

**macOS:**

- Xcode Command Line Tools: `xcode-select --install`

**Windows:**

- Disarankan menggunakan MSYS2 atau lingkungan serupa Unix lainnya
- Runtime WebView2 (biasanya sudah terinstal di Windows 11)

**Linux:**

- `gcc`, `pkg-config`, `libgtk-4-dev`, `libwebkitgtk-6.0-dev` (tumpukan GTK4 bawaan)
- Instal melalui: `sudo apt install build-essential pkg-config libgtk-4-dev libwebkitgtk-6.0-dev` (Debian/Ubuntu)
- Untuk jalur build `-tags gtk3` lama, instal juga `libgtk-3-dev` dan `libwebkit2gtk-4.1-dev`

## Ikhtisar Proses Kontribusi

Alur kerja kontribusi umumnya mengikuti langkah-langkah berikut:

1. **Fork & Klona** - Buat salinan repositori Wails milik Anda sendiri
2. **Siapkan** - Build CLI Wails dan verifikasi lingkungan Anda
3. **Buat Branch** - Buat branch fitur untuk perubahan Anda
4. **Kembangkan** - Buat perubahan sesuai standar pengodean kami
5. **Uji** - Jalankan pengujian untuk memastikan semuanya berfungsi
6. **Commit** - Lakukan commit dengan pesan commit konvensional yang jelas
7. **Kirim** - Buka pull request untuk ditinjau
8. **Iterasikan** - Tanggapi umpan balik dan lakukan penyesuaian
9. **Gabungkan** - Setelah disetujui, perubahan Anda menjadi bagian dari Wails!

## Panduan Langkah demi Langkah

Pilih jenis kontribusi Anda:

@tabs
[Perbaikan Bug]
@steps
### Temukan atau Laporkan Bug
- Periksa apakah bug tersebut sudah dilaporkan di [GitHub Issues](https://github.com/wailsapp/wails/issues)
- Jika belum, buat issue baru yang menyertakan langkah-langkah untuk mereproduksinya
- Tunggu konfirmasi sebelum mulai mengerjakannya

### Fork dan Klona
Fork repositori di [github.com/wailsapp/wails/fork](https://github.com/wailsapp/wails/fork)

Klona fork Anda:

```bash
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails
git remote add upstream https://github.com/wailsapp/wails.git
```

### Build dan Verifikasi
Build Wails dan verifikasi bahwa Anda dapat mereproduksi bug tersebut:

```bash
cd v3
go build -o ../wails3 ./cmd/wails3

# Reproduce the bug to understand it
```

### Buat Branch Perbaikan Bug
Buat branch untuk perbaikan Anda:

```bash
git checkout -b fix/issue-123-window-crash
```

### Perbaiki Bug
- Buat perubahan seminimal mungkin yang diperlukan untuk memperbaiki bug
- Jangan merefaktor kode yang tidak terkait
- Tambahkan atau perbarui pengujian untuk mencegah regresi

```bash
# Make your changes
# Add tests in *_test.go files
```

### Uji Perbaikan Anda
Jalankan pengujian untuk memastikan perbaikan berfungsi:

```bash
go test ./...

# Test the specific package
go test ./pkg/application -v

# Run with race detector
go test ./... -race
```

### Commit Perbaikan Anda
Lakukan commit dengan pesan yang jelas:

```bash
git commit -m "fix: prevent window crash when closing during initialization

Fixes #123"
```

### Kirim Pull Request
Lakukan push dan buat PR:

```bash
git push origin fix/issue-123-window-crash
```

Dalam deskripsi PR Anda:

- Jelaskan bug dan akar penyebabnya
- Jelaskan perbaikan Anda
- Cantumkan referensi issue: "Fixes #123"
- Sertakan perilaku sebelum dan sesudah perbaikan

### Tanggapi Umpan Balik
Tindak lanjuti komentar tinjauan dan perbarui PR Anda sesuai kebutuhan.

@end

[WEP (Peningkatan)]
@steps
### Tulis WEP
- Baca [proses WEP (Wails Enhancement Proposal)](https://github.com/wailsapp/wails/blob/master/v3/wep/README.md)
- Salin templat WEP ke `v3/wep/proposals/<name>/proposal.md`
- Buka draf PR berjudul `[WEP] <title>` yang hanya berisi WEP
- Tunggu keputusan pengelola sebelum mengimplementasikannya

### Fork dan Klona
Fork repositori di [github.com/wailsapp/wails/fork](https://github.com/wailsapp/wails/fork)

Klona fork Anda:

```bash
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails
git remote add upstream https://github.com/wailsapp/wails.git
```

### Siapkan Lingkungan Pengembangan
Build Wails dan verifikasi lingkungan Anda:

```bash
cd v3
go build -o ../wails3 ./cmd/wails3

# Run tests to ensure everything works
go test ./...
```

### Buat Branch Fitur
Buat branch dengan nama yang deskriptif:

```bash
git checkout -b feat/window-transparency-support
```

### Implementasikan Fitur
- Ikuti [Standar Penulisan Kode](/contributing/standards/) kami
- Pastikan perubahan tetap berfokus pada fitur tersebut
- Tulis kode yang rapi dan terdokumentasi
- Tambahkan pengujian yang menyeluruh

```bash
# Example: Adding a new window method
# 1. Add to window.go interface
# 2. Implement in platform files (darwin, windows, linux)
# 3. Add tests
# 4. Update documentation
```

### Uji Secara Menyeluruh
Uji fitur Anda:

```bash
# Unit tests
go test ./pkg/application -v

# Integration test - create a test app
cd ..
./wails3 init -n feature-test
cd feature-test
# Add code using your new feature
../wails3 dev
```

### Dokumentasikan Fitur Anda
- Tambahkan docstring ke semua API publik
- Perbarui dokumentasi yang relevan di `/docs/mpress/content/`
- Tambahkan contoh jika relevan

### Buat Commit Sesuai Konvensi
Gunakan conventional commit:

```bash
git commit -m "feat: add window transparency support

- Add SetTransparent() method to Window API
- Implement for macOS, Windows, and Linux
- Add tests and documentation

Closes #456"
```

### Kirim Pull Request
Lakukan push dan buat PR:

```bash
git push origin feat/window-transparency-support
```

Dalam PR Anda:

- Jelaskan fitur dan kasus penggunaannya
- Tampilkan contoh atau tangkapan layar
- Cantumkan setiap breaking change
- Referensikan PR WEP yang telah diterima

### Lakukan Iterasi Berdasarkan Hasil Tinjauan
Maintainer mungkin meminta perubahan. Bersabarlah dan bekerjalah secara kolaboratif.

@end

[Dokumentasi]
PR koreksi dapat langsung diajukan tanpa perlu membuka issue terlebih dahulu. Koreksi yang hanya menyangkut dokumentasi tidak memerlukan pengujian kode yang gagal. Ikuti [Perbaiki dokumentasi](/contributing/documentation/) untuk langkah-langkah instalasi M-Press, path sumber, pratinjau, validasi, dan PR.

@end

## Menemukan Issue untuk Dikerjakan

- Cari label [`good first issue`](https://github.com/wailsapp/wails/labels/good%20first%20issue)
- Periksa issue [`help wanted`](https://github.com/wailsapp/wails/labels/help%20wanted)
- Telusuri [issue yang terbuka](https://github.com/wailsapp/wails/issues) dan mintalah agar issue tersebut ditugaskan kepada Anda

## Mendapatkan Bantuan

- **Discord:** Bergabunglah dengan [Discord Wails](https://discord.gg/JDdSxwjhGf)
- **Diskusi:** Buat postingan di [GitHub Discussions](https://github.com/wailsapp/wails/discussions)
- **Issue:** Buka issue untuk bug yang dapat direproduksi; gunakan Discussions untuk pertanyaan dan PR WEP untuk peningkatan

## Kode Etik

Bersikaplah penuh hormat, konstruktif, dan ramah. Kita membangun komunitas yang bersahabat dengan fokus menciptakan perangkat lunak hebat bersama-sama.

## Langkah Berikutnya

- Siapkan [Lingkungan Pengembangan](/contributing/setup/) Anda
- Tinjau [Standar Penulisan Kode](/contributing/standards/) kami
- Jelajahi [Dokumentasi Teknis](/contributing/overview/)
