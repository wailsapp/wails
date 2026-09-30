---
title: "Asosiasi File"
description: "Konfigurasikan asosiasi file untuk aplikasi Wails Anda"
slug: "guides/file-associations"
sourcePath: "guides/file-associations.md"
---

Platform yang Relevan: <span class="mpress-badge mpress-badge-note">Windows</span> <span class="mpress-badge mpress-badge-success">macOS</span>

<br/>

Asosiasi file memungkinkan aplikasi Anda menangani jenis file tertentu saat pengguna membukanya. Fitur ini sangat berguna untuk editor teks, penampil gambar, atau aplikasi apa pun yang menangani format file tertentu. Panduan ini menjelaskan cara menerapkan asosiasi file dalam aplikasi Wails v3 Anda.

## Ikhtisar

Dukungan asosiasi file di Wails v3 saat ini tersedia untuk:

- Windows (paket penginstal NSIS)
- macOS (bundel aplikasi)

## Konfigurasi

Asosiasi file dikonfigurasikan dalam file `config.yml` yang berada di direktori `build` proyek Anda.

### Konfigurasi Dasar

Untuk menyiapkan asosiasi file:

1. Buka `build/config.yml`
2. Tambahkan asosiasi file Anda di bagian `fileAssociations`
3. Jalankan `wails3 update build-assets` untuk memperbarui aset build
4. Atur bidang `FileAssociations` dalam opsi aplikasi
5. Kemas aplikasi Anda menggunakan `wails3 package`

Berikut contoh konfigurasinya:

```yaml
fileAssociations:
  - ext: myapp
    name: MyApp Document
    description: MyApp Document File
    iconName: myappFileIcon
    role: Editor
  - ext: custom
    name: Custom Format
    description: Custom File Format
    iconName: customFileIcon
    role: Editor
```

### Properti Konfigurasi

| Properti | Deskripsi | Platform |
| --- | --- | --- |
| ext | Ekstensi file tanpa tanda titik di awal (misalnya, `txt`) | Semua |
| name | Nama tampilan untuk jenis file | Semua |
| description | Deskripsi yang ditampilkan dalam properti file | Windows |
| iconName | Nama file ikon (tanpa ekstensi) dalam folder build | Semua |
| role | Peran aplikasi untuk jenis file ini (misalnya, `Editor`, `Viewer`) | macOS |
| mimeType | Jenis MIME untuk file (misalnya, `image/jpeg`) | macOS |

## Mendengarkan Peristiwa Pembukaan File

Untuk menangani peristiwa pembukaan file dalam aplikasi, Anda dapat mendengarkan peristiwa `events.Common.ApplicationOpenedWithFile`:

```go {title="main.go"}
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)

func main() {
    app := application.New(application.Options{
        Name: "MyApp",
        FileAssociations: []string{".txt", ".md"}, // Specify supported extensions
    })

    // Listen for files being used to open the application
    app.Event.OnApplicationEvent(events.Common.ApplicationOpenedWithFile, func(event *application.ApplicationEvent) {
        associatedFile := event.Context().Filename()
        app.Dialog.Info().SetMessage("Application opened with file: " + associatedFile).Show()
    })

    // Create your window and run the app...
}

```

## Tutorial Langkah demi Langkah

Mari siapkan asosiasi file untuk editor teks sederhana langkah demi langkah:

@steps
### Buat Ikon
- Buat ikon untuk jenis file Anda (ukuran yang disarankan: 16x16, 32x32, 48x48, 256x256)
- Simpan ikon dalam folder `build` proyek Anda
- Beri nama ikon sesuai dengan konfigurasi `iconName` Anda (misalnya, `textFileIcon.png`)

@note{type="tip"}
Anda dapat menggunakan `wails3 generate icons` untuk menghasilkan ikon yang diperlukan. Jalankan `wails3 generate icons --help` untuk informasi selengkapnya.

@end

- Untuk macOS, tambahkan pernyataan penyalinan seperti `cp build/darwin/documenticon.icns {{.BIN_DIR}}/{{.APP_NAME}}.app/Contents/Resources` dalam tugas `create:app:bundle:`.

### Konfigurasikan Asosiasi File
Edit file `build/config.yml` untuk menambahkan asosiasi file Anda:

```yaml
# build/config.yml
fileAssociations:
  - ext: txt
    name: Text Document
    description: Plain Text Document
    iconName: textFileIcon
    role: Editor
```

### Perbarui Aset Build
Jalankan perintah berikut untuk memperbarui aset build:

```bash
wails3 update build-assets
```

### Atur Asosiasi File dalam Opsi Aplikasi
Dalam file `main.go`, atur bidang `FileAssociations` dalam opsi aplikasi:

```go
app := application.New(application.Options{
  Name: "MyApp",
  FileAssociations: []string{".txt", ".md"}, // Specify supported extensions
})
```

@note{type="tip" title="Mengapa ekstensi file diperlukan dalam konfigurasi aplikasi dan config.yml?"}
Di Windows, ketika file dibuka melalui asosiasi file, aplikasi diluncurkan dengan nama file sebagai argumen pertama aplikasi. Aplikasi tidak dapat mengetahui apakah argumen pertama merupakan file atau argumen baris perintah. Oleh karena itu, aplikasi menggunakan bidang `FileAssociations` dalam opsi aplikasi untuk menentukan apakah argumen pertama merupakan file yang diasosiasikan.

@end

### Kemas Aplikasi Anda
Kemas aplikasi Anda menggunakan perintah berikut:

```bash
wails3 package
```

Aplikasi yang telah dikemas akan dibuat dalam direktori `bin`. Setelah itu, Anda dapat menginstal dan menguji aplikasi tersebut.

## Catatan Tambahan

- Ikon sebaiknya disediakan dalam format PNG di folder build
- Pengujian asosiasi file mengharuskan aplikasi yang telah dikemas untuk diinstal

@end
