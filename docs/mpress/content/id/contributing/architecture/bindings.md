---
title: "Sistem Binding"
description: "Cara sistem binding mengumpulkan, memproses, dan menghasilkan kode JavaScript/TypeScript"
slug: "contributing/architecture/bindings"
sourcePath: "contributing/architecture/bindings.md"
---

Panduan ini menjelaskan cara kerja internal sistem binding Wails dan memberikan wawasan bagi pengembang yang ingin memahami mekanisme di balik pembuatan kode otomatis.

## Ikhtisar Arsitektur

Sistem binding Wails terdiri dari tiga komponen utama:

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
- **Package**: Merepresentasikan paket Go yang sedang dianalisis dan menyimpan layanan, model, serta direktif yang dikumpulkan
- **Service**: Mengumpulkan informasi tentang tipe layanan dan metode-metodenya
- **Model**: Mengumpulkan informasi terperinci tentang tipe model, termasuk field, nilai, dan parameter tipe
- **Directive**: Mengurai dan menafsirkan direktif `//wails:` dalam kode sumber Go

### Alur Pengumpulan

1. Collector memindai paket Go yang ditentukan dalam proyek
2. Collector mengidentifikasi tipe layanan (struct dengan metode yang akan diekspos ke frontend)
3. Untuk setiap layanan, collector mengumpulkan informasi tentang metode-metodenya
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

1. Untuk setiap layanan, renderer menghasilkan file JavaScript/TypeScript berisi fungsi-fungsi yang mencerminkan metode layanan tersebut
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

- `//wails:inject`: Menyisipkan kode JavaScript/TypeScript khusus ke dalam binding yang dihasilkan
- `//wails:include`: Menyertakan file tambahan bersama binding yang dihasilkan
- `//wails:internal`: Menandai tipe atau metode sebagai internal sehingga tidak diekspor ke frontend
- `//wails:ignore`: Mengabaikan metode sepenuhnya selama pembuatan binding
- `//wails:id`: Menentukan ID khusus untuk metode, menggantikan ID default berbasis hash

### Pemrosesan Direktif

1. Selama tahap pengumpulan, kolektor mengidentifikasi dan mengurai direktif dalam kode Go
2. Direktif disimpan bersama deklarasi yang terkait (layanan, metode, model, dan sebagainya)
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
  - `*` - Kelas dan antarmuka
  - `c` - Hanya kelas
  - `i` - Hanya antarmuka


Contoh:

```go
//wails:inject j*:console.log("JavaScript only");
//wails:inject t*:console.log("TypeScript only");
```

### ID Metode Khusus

Secara default, metode diidentifikasi dengan ID berbasis hash. Namun, Anda dapat menentukan ID khusus menggunakan direktif `//wails:id`:

```go
//wails:id 42
func (s *Service) CustomIDMethod() {}
```

Ini dapat berguna untuk mempertahankan kompatibilitas saat memfaktorkan ulang kode.

## Pertimbangan Performa

Generator binding dirancang agar efisien, tetapi ada beberapa hal yang perlu diperhatikan:

1. Eksekusi pertama akan lebih lambat karena membangun cache paket yang akan dipindai
2. Eksekusi berikutnya akan lebih cepat karena menggunakan informasi dalam cache
3. Generator memproses semua paket dalam proyek, yang dapat memakan waktu untuk proyek besar
4. Anda dapat menggunakan flag `-clean` untuk membersihkan direktori output sebelum pembuatan kode

## Debugging

Jika mengalami masalah dalam pembuatan binding, Anda dapat menggunakan flag `-v` untuk mengaktifkan output debug:

```bash
wails3 generate bindings -v
```

Ini akan memberikan informasi terperinci tentang proses pengumpulan dan rendering, yang dapat membantu mengidentifikasi sumber masalah.
