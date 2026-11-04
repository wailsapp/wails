---
title: "Sistem Binding"
description: "Cara sistem binding mengumpulkan, memproses, dan menghasilkan kode JavaScript/TypeScript"
slug: "contributing/architecture/bindings"
sourcePath: "contributing/architecture/bindings.md"
---

Panduan ini menjelaskan cara kerja internal sistem binding Wails dan memberikan wawasan bagi pengembang yang ingin memahami mekanisme di balik pembuatan kode otomatis.

## Gambaran Umum Arsitektur

Sistem binding Wails terdiri atas tiga komponen utama:

1. **Pengumpulan**: Menganalisis kode Go untuk mengekstrak informasi tentang layanan, model, dan deklarasi lainnya
2. **Konfigurasi**: Mengelola pengaturan dan opsi untuk proses pembuatan binding
3. **Rendering**: Menghasilkan kode JavaScript/TypeScript berdasarkan informasi yang dikumpulkan

@filetree
- internal/generator/
  - collect/     # Analisis paket dan ekstraksi informasi
  - config/      # Struktur dan antarmuka konfigurasi
  - render/      # Pembuatan kode untuk JS/TS
@end

## Proses Pengumpulan

Proses pengumpulan bertugas menganalisis paket Go dan mengekstrak informasi tentang layanan, model, dan deklarasi lainnya. Proses ini ditangani oleh paket `collect`.

### Komponen Utama

- **Collector**: Mengelola informasi paket dan menyimpan data yang dikumpulkan dalam cache
- **Package**: Merepresentasikan paket Go yang sedang dianalisis serta menyimpan layanan, model, dan direktif yang dikumpulkan
- **Service**: Mengumpulkan informasi tentang tipe layanan dan metodenya
- **Model**: Mengumpulkan informasi terperinci tentang tipe model, termasuk field, nilai, dan parameter tipe
- **Directive**: Mengurai dan menafsirkan direktif `//wails:` dalam kode sumber Go

### Alur Pengumpulan

1. Collector memindai paket Go yang ditentukan dalam proyek
2. Collector mengidentifikasi tipe layanan (struct dengan metode yang akan diekspos ke frontend)
3. Untuk setiap layanan, collector mengumpulkan informasi tentang metodenya
4. Collector mengidentifikasi tipe model (struct yang digunakan sebagai parameter atau nilai kembalian dalam metode layanan)
5. Untuk setiap model, collector mengumpulkan informasi tentang field dan parameter tipenya
6. Collector memproses setiap direktif `//wails:` yang ditemukan dalam kode

## Proses Rendering

Proses rendering bertugas menghasilkan kode JavaScript/TypeScript berdasarkan informasi yang dikumpulkan. Proses ini ditangani oleh paket `render`.

### Komponen Utama

- **Renderer**: Mengoordinasikan rendering file layanan, model, dan indeks
- **Module**: Merepresentasikan satu modul JavaScript/TypeScript yang dihasilkan
- **Templates**: Templat teks yang digunakan untuk pembuatan kode

### Alur Rendering

1. Untuk setiap layanan, renderer menghasilkan file JavaScript/TypeScript berisi fungsi yang mencerminkan metode layanan tersebut
2. Untuk setiap model, renderer menghasilkan class JavaScript/TypeScript yang mencerminkan struct model tersebut
3. Renderer menghasilkan file indeks yang mengekspor ulang semua layanan dan model
4. Renderer menerapkan setiap injeksi kode khusus yang ditentukan oleh direktif `//wails:inject`

## Pemetaan Tipe

Salah satu aspek terpenting dari sistem binding adalah cara tipe Go dipetakan ke tipe JavaScript/TypeScript. Berikut ringkasan pemetaannya:

| Tipe Go | Tipe JavaScript | Tipe TypeScript |
| --- | --- | --- |
| `bool` | `boolean` | `boolean` |
| `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `float32`, `float64` | `number` | `number` |
| `string` | `string` | `string` |
| `[]byte` | `Uint8Array` | `Uint8Array` |
| `[]T` | `Array<T>` | `T[]` |
| `map[string]V` | `Object` | `{ [_: string]: V }` |
| `map[K]V` (`K` non-string) | `Object` | `{ [_ in K]?: V }` |
| `struct` | `Object` | Class khusus |
| `interface{}` | `any` | `any` |
| `*T` | `T \| null` | `T \| null` |
| `func` | Tidak didukung | Tidak didukung |
| `chan` | Tidak didukung | Tidak didukung |

## Sistem Direktif

Sistem binding mendukung beberapa direktif yang dapat digunakan untuk menyesuaikan kode yang dihasilkan. Direktif ini ditambahkan sebagai komentar dalam kode Go Anda.

### Direktif yang Tersedia

- `//wails:inject`: Menginjeksikan kode JavaScript/TypeScript khusus ke dalam binding yang dihasilkan
- `//wails:include`: Menyertakan file tambahan bersama binding yang dihasilkan
- `//wails:internal`: Menyembunyikan ekspor dalam SDK yang dihasilkan; direktif ini tidak mencegah pemanggilan metode Go yang diekspor saat runtime
- `//wails:ignore`: Sepenuhnya mengabaikan metode selama pembuatan binding
- `//wails:id`: Menentukan ID khusus untuk metode dan menggantikan ID default berbasis hash

Direktif ini mengontrol kode yang dihasilkan. Dalam build normal, registrasi runtime tetap merefleksikan metode layanan yang diekspor; gunakan metode yang tidak diekspor atau tipe yang tidak diikat untuk operasi privat, dan terapkan otorisasi dalam metode yang diekspos.

### Pemrosesan Direktif

1. Selama tahap pengumpulan, pengumpul mengidentifikasi dan mengurai direktif dalam kode Go
2. Direktif disimpan bersama deklarasi terkait (layanan, metode, model, dan sebagainya)
3. Selama tahap rendering, perender menerapkan direktif untuk menyesuaikan kode yang dihasilkan

## Fitur Lanjutan

### Pembuatan Kode Bersyarat

Sistem binding mendukung pembuatan kode bersyarat menggunakan prefiks kondisi dua karakter untuk direktif `include` dan `inject`:

```
<language><style>:<content>
```

Dengan ketentuan:

- `<language>` dapat berupa:
  - `*` - JavaScript dan TypeScript
  - `j` - Hanya JavaScript
  - `t` - Hanya TypeScript


- `<style>` dapat berupa:
  - `*` - Class dan interface
  - `c` - Hanya class
  - `i` - Hanya interface


Contoh:

```go
//wails:inject j*:console.log("JavaScript only");
//wails:inject t*:console.log("TypeScript only");
```

### ID Metode Kustom

Secara default, metode diidentifikasi berdasarkan ID berbasis hash. Namun, Anda dapat menentukan ID kustom menggunakan direktif `//wails:id`:

```go
//wails:id 42
func (s *Service) CustomIDMethod() {}
```

Ini dapat berguna untuk mempertahankan kompatibilitas saat melakukan refaktorisasi kode.

## Pertimbangan Performa

Generator binding dirancang agar efisien, tetapi ada beberapa hal yang perlu diperhatikan:

1. Proses pertama akan lebih lambat karena generator membuat cache paket yang akan dipindai
2. Proses berikutnya akan lebih cepat karena menggunakan informasi yang telah disimpan dalam cache
3. Generator memproses semua paket dalam proyek, yang dapat memakan waktu untuk proyek besar
4. Anda dapat menggunakan flag `-clean` untuk membersihkan direktori output sebelum pembuatan kode

## Penelusuran Kesalahan

Jika mengalami masalah saat membuat binding, Anda dapat menggunakan flag `-v` untuk mengaktifkan output debug:

```bash
wails3 generate bindings -v
```

Output ini akan memberikan informasi terperinci tentang proses pengumpulan dan rendering, yang dapat membantu mengidentifikasi sumber masalah.
