---
title: "Berkontribusi"
description: "Berkontribusi pada Wails"
slug: "contributing"
sourcePath: "contributing/index.md"
---

## Selamat Datang, Kontributor!

Kami menyambut kontribusi untuk Wails! Baik Anda memperbaiki bug, menambahkan fitur, maupun menyempurnakan dokumentasi, bantuan Anda sangat kami hargai.

## Cara Berkontribusi

### 1. Laporkan Masalah

Menemukan bug? [Buka issue](https://github.com/wailsapp/wails/issues/new) dengan menyertakan:

- Deskripsi yang jelas
- Langkah-langkah untuk mereproduksi masalah
- Perilaku yang diharapkan dibandingkan dengan perilaku aktual
- Informasi sistem
- Contoh kode

### 2. Sempurnakan Dokumentasi

PR koreksi diterima tanpa perlu issue sebelumnya atau pengujian kode yang gagal.  
Ikuti [Menulis dokumentasi dengan M-Press](/contributing/documentation/) untuk mengunduh, menginstal di macOS/Linux/Windows, melihat pratinjau langsung, dan mengirimkan perubahan Anda.

Penyempurnaan dokumentasi selalu diterima:

- Perbaiki kesalahan ketik dan kekeliruan
- Tambahkan contoh
- Perjelas penjelasan
- Terjemahkan konten

### 3. Kirimkan Kode

Kontribusikan kode melalui pull request:

- Perbaikan bug
- Fitur baru
- Peningkatan performa
- Pengujian

### 4. Ajukan Penyempurnaan (WEP)

Fungsionalitas baru dan perubahan pada perilaku publik menggunakan proses Wails Enhancement Proposal (WEP). Proses ini menjaga pengembangan fitur tetap transparan dan memastikan setiap proposal yang diterima memiliki pelaksana. Jangan membuka issue permintaan fitur.

1. Jika ingin, sampaikan terlebih dahulu ide Anda di kategori [Ideas](https://github.com/wailsapp/wails/discussions/categories/ideas) pada GitHub Discussions atau di [Discord](https://discord.gg/JDdSxwjhGf) untuk mengukur minat.
2. Salin [`v3/wep/WEP_TEMPLATE.md`](https://github.com/wailsapp/wails/blob/releases/v3.0.0/v3/wep/WEP_TEMPLATE.md) ke `v3/wep/proposals/<proposal name>/proposal.md` dan isi setiap bagiannya.
3. Buka draf pull request berjudul `[WEP] <title>` yang hanya berisi proposal tersebut. PR itu merupakan tempat resmi untuk membahasnya.
4. Kumpulkan masukan dan dukungan (komentar dan reaksi jempol pada PR). Berikan waktu setidaknya dua minggu untuk berdiskusi, lalu sepakati siapa yang akan mengimplementasikan proposal tersebut.
5. Tandai PR sebagai siap ditinjau. Pengelola membuat keputusan akhir: proposal yang diterima akan diberi nomor WEP dan digabungkan.

Proses lengkap didokumentasikan dalam [`v3/wep/README.md`](https://github.com/wailsapp/wails/blob/releases/v3.0.0/v3/wep/README.md).

## Memulai

### Fork dan Kloning

```bash
# Fork the repository on GitHub
# Then clone your fork
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails

# Add upstream remote
git remote add upstream https://github.com/wailsapp/wails.git
```

### Build dari Kode Sumber

```bash
# Build the v3 CLI. Go downloads any required modules automatically.
cd v3
go build -o ../wails3 ./cmd/wails3

# Confirm the built CLI runs and reports its version.
../wails3 version
```

### Jalankan Pengujian

Pengujian merupakan bagian dari perubahan, bukan sekadar pemeriksaan awal terakhir. Letakkan pengujian unit di dekat kode yang diuji dan utamakan pengujian berbasis tabel ketika satu perilaku diperiksa dengan beberapa input atau kasus ekstrem. Beri nama pada setiap kasus agar kegagalan dapat menjelaskan skenarionya.

Logika baru dan yang diubah harus tercakup sepenuhnya. Targetkan cakupan pernyataan Go sebesar 100% untuk kode yang ditambahkan atau diubah oleh PR Anda; jangan menganggap persentase untuk seluruh repositori sebagai pengganti pengujian perubahan tersebut. Kekurangan cakupan bisa saja valid—misalnya, jalur kesalahan khusus OS atau kondisi yang tidak praktis untuk direproduksi tanpa perangkat keras nyata—tetapi jelaskan kekurangan tersebut dan alasan mengapa hal itu tidak dapat diuji secara wajar dalam deskripsi PR.

```bash
# Run all v3 tests
cd v3
go test ./...

# Run specific package tests
go test ./pkg/application

# Inspect coverage for the packages you changed
go test ./pkg/application -coverprofile=coverage.out
go tool cover -func=coverage.out
cd ..
```

Untuk rangkaian pengujian integrasi, deteksi race condition, dan perintah lengkap yang setara dengan CI, lihat [Pengujian & Integrasi Berkelanjutan](/contributing/testing-ci/).

## Membuat Perubahan

### Buat Branch

```bash
# Update master
git checkout master
git pull upstream master

# Create feature branch
git checkout -b feature/my-feature
```

### Buat Perubahan Anda

1. **Tulis kode** sesuai konvensi Go
2. **Tambahkan pengujian** untuk fungsionalitas baru
3. **Perbarui dokumentasi** jika diperlukan
4. **Jalankan pengujian** untuk memastikan tidak ada yang rusak
5. **Commit perubahan** dengan pesan yang jelas

### Panduan Commit

```bash
# Good commit messages
git commit -m "fix: resolve window focus issue on macOS"
git commit -m "feat: add support for custom window chrome"
git commit -m "docs: improve bindings documentation"

# Use conventional commits:
# - feat: New feature
# - fix: Bug fix
# - docs: Documentation
# - test: Tests
# - refactor: Code refactoring
# - chore: Maintenance
```

### Kirimkan Pull Request

```bash
# Push to your fork
git push origin feature/my-feature

# Open pull request on GitHub
# Provide clear description
# Reference related issues
```

## Panduan Pull Request

### Deskripsi PR yang Baik

```markdown
## Description
Brief description of changes

## Changes
- Added feature X
- Fixed bug Y
- Updated documentation

## Testing
- Tested on macOS 14
- Tested on Windows 11
- All tests passing

## Related Issues
Fixes #123
```

### Daftar Periksa PR

- [ ] Kode mengikuti konvensi Go
- [ ] Pengujian ditambahkan/diperbarui
- [ ] Dokumentasi diperbarui
- [ ] Semua pengujian berhasil
- [ ] Tidak ada perubahan yang merusak kompatibilitas (atau sudah didokumentasikan)
- [ ] Pesan commit jelas

## Panduan Kode

### Gaya Kode Go

```go
// ✅ Good: Clear, documented, tested
// ProcessData processes the input data and returns the result.
// It returns an error if the data is invalid.
func ProcessData(data string) (string, error) {
    if data == "" {
        return "", errors.New("data cannot be empty")
    }
    
    result := process(data)
    return result, nil
}

// ❌ Bad: No docs, no error handling
func ProcessData(data string) string {
    return process(data)
}
```

### Pengujian

```go
func TestProcessData(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"valid input", "test", "processed", false},
        {"empty input", "", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ProcessData(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("ProcessData() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("ProcessData() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

## Dokumentasi

### Menulis Dokumentasi

Ikuti [Menulis dokumentasi dengan M-Press](/contributing/documentation/) untuk alur kerja dokumentasi lengkap. Panduan tersebut mencakup versi yang disematkan, opsi penginstalan berkas yang dapat dieksekusi dan Go, berkas sumber yang dapat diedit, navigasi, pemeriksaan terjemahan, serta pengiriman PR terhadap `master`.

Pekerjaan yang hanya berkaitan dengan dokumentasi tidak memerlukan build aplikasi Wails atau dependensi platform. M-Press menampilkan pratinjau situs yang ada menggunakan `mpress dev` dari root repositori.

### Gaya Dokumentasi

- Gunakan ejaan bahasa Inggris Internasional
- Mulailah dengan masalahnya
- Berikan contoh yang berfungsi
- Sertakan pemecahan masalah
- Berikan referensi silang ke konten terkait

## Komunitas

### Dapatkan Bantuan

- **Discord:** [Bergabunglah dengan komunitas kami](https://discord.gg/JDdSxwjhGf)
- **GitHub Discussions:** Ajukan pertanyaan
- **GitHub Issues:** Laporkan bug

### Kode Etik

Bersikaplah penuh hormat, inklusif, dan profesional. Kita semua berada di sini untuk bersama-sama membangun perangkat lunak yang hebat. Lihat [Kode Etik](https://github.com/wailsapp/wails/blob/master/CODE_OF_CONDUCT.md) untuk detailnya.

## Pengakuan

Kontributor dicantumkan dalam:

- Catatan rilis
- Daftar kontributor
- Insight GitHub

Terima kasih telah berkontribusi pada Wails! 🎉

## Langkah Berikutnya

@cards{cols="2"}
◆ Repositori GitHub
Kunjungi repositori Wails.

[Lihat di GitHub →](https://github.com/wailsapp/wails)

---
◆ Komunitas Discord
Bergabunglah dengan komunitas.

[Bergabung dengan Discord →](https://discord.gg/JDdSxwjhGf)

---
📖 Dokumentasi
Baca dokumentasi.

[Jelajahi Dokumentasi →](/quick-start/why-wails/)

@end
