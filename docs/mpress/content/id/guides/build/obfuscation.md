---
title: "Build Terobfuscasi"
description: "Build aplikasi Wails Anda dengan Garble untuk melindungi kode sumber dari rekayasa balik"
slug: "guides/build/obfuscation"
sourcePath: "guides/build/obfuscation.md"
---

[Garble](https://github.com/burrowers/garble) adalah alat build Go yang menggantikan `go build` untuk mengganti nama simbol, mengobfuscasi konstanta, dan menghapus informasi debug dari berkas biner yang dihasilkan. Wails v3 menyediakan dukungan kelas satu untuk Garble melalui dua perintah baru.

## Prasyarat

- **Go 1.26.2 atau yang lebih baru** — diperlukan oleh Garble v0.16.0
- **Garble v0.16.0**

```bash
go install mvdan.cc/garble@v0.16.0
```

@note{type="tip"}
Versi minimum Go yang diperlukan Garble berubah di antara rilis. Jika Anda menggunakan toolchain Go yang lebih lama, sebelum menginstal, periksa [halaman rilis Garble](https://github.com/burrowers/garble/releases) untuk menemukan versi yang sesuai dengan toolchain Anda.

@end

## Wajib: tambahkan tag JSON ke tipe layanan Anda

Setiap struct yang dikembalikan atau diterima oleh metode layanan yang di-binding harus memiliki tag JSON eksplisit pada setiap field yang diekspor:

```go
// Without tags — breaks under Garble
type OrderSummary struct {
    ID        int
    Total     float64
    LineItems []LineItem
}

// With tags — safe under Garble
type OrderSummary struct {
    ID        int       `json:"id"`
    Total     float64   `json:"total"`
    LineItems []LineItem `json:"lineItems"`
}
```

@note{type="caution"}
Garble mengganti nama field struct yang diekspor, sedangkan Wails meneruskan struct tersebut ke `json.Marshal` melalui parameter `interface{}` yang tidak dapat dilacak secara statis oleh Garble. Build terobfuscasi tanpa tag JSON akan berhasil dikompilasi, tetapi frontend Anda akan menerima nama field yang kacau atau kosong saat runtime. Tambahkan tag sebelum menjalankan build terobfuscasi.

@end

Tipe bawaan Wails — `Screen`, `Rect`, `Point`, `Size`, `EnvironmentInfo`, `OSInfo`, `Capabilities` — sudah memiliki tag. Anda hanya perlu menambahkan tag ke tipe Anda sendiri.

## Melakukan build dengan obfuscation

@steps
### Buat berkas ID stabil
Jalankan perintah ini setiap kali Anda menambahkan, mengganti nama, atau menghapus metode layanan yang di-binding:

```bash
wails3 generate bindings -obfuscated
```

Perintah ini membuat `wails_obfuscated.gen.go` di direktori package utama Anda — commit berkas ini.

### Build dengan Garble
```bash
wails3 build --obfuscated
```

Mem-build aplikasi menggunakan binding yang terobfuscasi.

@end

## Meneruskan flag tambahan ke Garble

Gunakan `--garbleargs` untuk meneruskan opsi secara langsung ke `garble`:

```bash
# Obfuscate string literals and reduce binary size
wails3 build --obfuscated --garbleargs "-literals -tiny"

# Reproducible output — same seed produces the same binary
wails3 build --obfuscated --garbleargs "-seed=deadbeef"
```

Lihat [dokumentasi Garble](https://github.com/burrowers/garble#flags) untuk mengetahui daftar lengkap flag yang didukung.

## Lanjutan: menulis berkas ID ke package lain

Secara default, `wails_obfuscated.gen.go` ditulis di samping package `main` Anda. Jika proyek Anda menyimpan layanan dalam subpackage yang diimpor oleh `main`, Anda dapat menulis berkas tersebut di sana menggunakan `-obfuscated-output`:

```bash
wails3 generate bindings -obfuscated -obfuscated-output ./internal/services
```

@note{type="caution"}
Package tujuan harus diimpor — secara langsung atau transitif — oleh package `main` Anda agar `init()` miliknya dijalankan saat startup. Jika package tersebut tidak dapat dijangkau, ID stabil tidak akan pernah didaftarkan dan pemanggilan binding akan gagal (misalnya, dengan error `binding not found` saat runtime).

@end

## Pemecahan masalah

### `garble: command not found`

Garble belum diinstal atau `$(go env GOPATH)/bin` tidak ada dalam `PATH` Anda.

```bash
go install mvdan.cc/garble@v0.16.0
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Frontend menerima nilai field yang salah atau kosong

Tipe nilai kembalian layanan Anda tidak memiliki tag `json:"..."`. Periksa setiap struct yang dikembalikan oleh metode yang di-binding dan tambahkan tag eksplisit ke setiap field yang diekspor.

### Error `binding not found` di konsol browser

Berkas ID stabil tidak ada atau tidak disertakan saat kompilasi. Periksa hal berikut:

- `wails_obfuscated.gen.go` tersedia di direktori package utama Anda (atau direktori yang Anda teruskan ke `-obfuscated-output`)
- Anda menjalankan `wails3 build --obfuscated`, yang menambahkan tag build `wails_obfuscated`
- Jika Anda menggunakan `-obfuscated-output`, package tujuan diimpor oleh `main`

### Windows Defender menandai hasil build sebagai virus

Berkas biner Go yang diobfuscasi dengan Garble ditandai secara heuristik oleh Windows Defender selama proses build karena tidak memiliki simbol debug dan menyerupai executable yang dikemas. Build gagal dengan pesan berikut:

```
open C:\Users\...\AppData\Local\Temp\go-build...\a.out.exe: The file contains a virus or potentially unwanted software.
```

Tambahkan direktori sementara Anda (tempat Go menulis artefak build perantara) dan direktori proyek Anda ke daftar pengecualian Defender:

```powershell
Add-MpPreference -ExclusionPath "$env:TEMP"
Add-MpPreference -ExclusionPath "C:\path\to\your\project"
```

Pengecualian ini hanya berlaku untuk path yang ditentukan dan tidak menonaktifkan Defender secara global.

### Build gagal dengan `unsupported Go version`

Garble v0.16.0 memerlukan Go 1.26.2 atau yang lebih baru. Upgrade Go atau lihat [halaman rilis Garble](https://github.com/burrowers/garble/releases) untuk menemukan versi yang kompatibel dengan toolchain Anda.
