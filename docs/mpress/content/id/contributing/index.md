---
title: "Berkontribusi"
description: "Berkontribusi pada Wails"
slug: "contributing"
sourcePath: "contributing/index.md"
---

## Selamat Datang, Para Kontributor!

Kami menyambut kontribusi untuk Wails! Baik Anda memperbaiki bug, menambahkan fitur, maupun menyempurnakan dokumentasi, bantuan Anda sangat kami hargai.

## Cara Berkontribusi

### 1. Laporkan Masalah

Menemukan bug? [Buka issue](https://github.com/wailsapp/wails/issues/new) yang menyertakan:

- Deskripsi yang jelas
- Langkah-langkah untuk mereproduksi masalah
- Perilaku yang diharapkan dan perilaku aktual
- Informasi sistem
- Contoh kode

### 2. Sempurnakan Dokumentasi

PR perbaikan dokumentasi dapat diajukan tanpa issue sebelumnya atau pengujian kode yang gagal.  
Ikuti [Perbaiki dokumentasi](/contributing/documentation/) untuk melihat pratinjau dan memvalidasi perubahan dengan M-Press.

Penyempurnaan dokumentasi selalu kami sambut:

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

Fungsionalitas baru dan perubahan pada perilaku publik menggunakan proses Wails Enhancement Proposal (WEP). Proses ini menjaga transparansi pengembangan fitur dan memastikan setiap proposal yang diterima memiliki pelaksana. Jangan buka issue permintaan fitur.

1. Jika ingin, sampaikan terlebih dahulu ide Anda dalam kategori [Ideas](https://github.com/wailsapp/wails/discussions/categories/ideas) di GitHub Discussions atau di [Discord](https://discord.gg/JDdSxwjhGf) untuk mengukur minat.
2. Salin [`v3/wep/WEP_TEMPLATE.md`](https://github.com/wailsapp/wails/blob/master/v3/wep/WEP_TEMPLATE.md) ke `v3/wep/proposals/<proposal name>/proposal.md`, lalu isi setiap bagiannya.
3. Buka draft pull request berjudul `[WEP] <title>` yang hanya berisi proposal tersebut. PR itu merupakan tempat resmi untuk membahasnya.
4. Kumpulkan masukan dan dukungan (komentar dan reaksi jempol pada PR). Berikan waktu setidaknya dua minggu untuk berdiskusi dan sepakati siapa yang akan mengimplementasikan proposal tersebut.
5. Tandai PR sebagai siap ditinjau. Maintainer mengambil keputusan akhir: proposal yang diterima akan diberi nomor WEP dan digabungkan.

Proses lengkapnya didokumentasikan dalam [`v3/wep/README.md`](https://github.com/wailsapp/wails/blob/master/v3/wep/README.md).

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

Pengujian merupakan bagian dari perubahan, bukan sekadar pemeriksaan dasar di tahap akhir untuk memastikan perangkat lunak berfungsi. Tempatkan unit test di dekat kode yang diuji dan utamakan pengujian berbasis tabel setiap kali satu perilaku diuji dengan beberapa input atau kasus ekstrem. Beri nama setiap kasus agar kegagalan menjelaskan skenarionya.

Logika baru dan yang diubah sebaiknya tercakup sepenuhnya. Targetkan cakupan statement Go sebesar 100% untuk kode yang ditambahkan atau diubah oleh PR Anda; jangan gunakan persentase cakupan seluruh repositori sebagai pengganti pengujian perubahan tersebut. Kekurangan cakupan dapat dibenarkan—misalnya, jalur kesalahan yang hanya berlaku pada OS tertentu atau kondisi yang tidak praktis untuk direproduksi tanpa perangkat keras nyata—tetapi jelaskan kekurangan itu dan alasan hal tersebut tidak dapat diuji secara wajar dalam deskripsi PR.

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

### Pedoman Commit

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

## Pedoman Pull Request

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
- [ ] Pengujian ditambahkan atau diperbarui
- [ ] Dokumentasi diperbarui
- [ ] Semua pengujian lulus
- [ ] Tidak ada breaking change (atau sudah didokumentasikan)
- [ ] Pesan commit jelas

## Pedoman Kode

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

Dokumentasi menggunakan M-Press. Edit file `.md` di bawah  
`docs/mpress/content/`, lalu pratinjau dan validasi dari direktori akar repositori:

```bash
go install github.com/leaanthony/mpress/cmd/mpress@v1.0.17
mpress dev
mpress build --strict --no-purge-css
mpress check
# Use `python` on Windows if `python3` is unavailable.
python3 docs/mpress/scripts/check_site.py docs/mpress/site
```

### Gaya Dokumentasi

- Gunakan ejaan bahasa Inggris internasional
- Mulailah dengan masalah yang akan diselesaikan
- Berikan contoh yang dapat dijalankan
- Sertakan panduan pemecahan masalah
- Cantumkan referensi silang ke konten terkait

## Komunitas

### Dapatkan Bantuan

- **Discord:** [Bergabunglah dengan komunitas kami](https://discord.gg/JDdSxwjhGf)
- **GitHub Discussions:** Ajukan pertanyaan
- **GitHub Issues:** Laporkan bug

### Kode Etik

Bersikaplah saling menghormati, inklusif, dan profesional. Kita semua berada di sini untuk bersama-sama membangun perangkat lunak yang hebat. Lihat [Kode Etik](https://github.com/wailsapp/wails/blob/master/CODE_OF_CONDUCT.md) untuk detail selengkapnya.

## Pengakuan

Kontributor mendapat pengakuan dalam:

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

[Bergabung ke Discord →](https://discord.gg/JDdSxwjhGf)

---
📖 Dokumentasi
Baca dokumentasi.

[Jelajahi Dokumentasi →](/quick-start/why-wails/)

@end
