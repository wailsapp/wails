---
title: "Dock \u0026 Bilah Tugas"
description: "Kelola visibilitas ikon dock dan tampilkan lencana di macOS dan Windows"
slug: "features/platform/dock"
sourcePath: "features/platform/dock.md"
---

## Pendahuluan

Wails menyediakan layanan Dock lintas platform untuk aplikasi desktop. Layanan ini memungkinkan Anda untuk:

- Menyembunyikan dan menampilkan ikon aplikasi di Dock macOS
- Menampilkan lencana pada ubin aplikasi atau ikon dock/bilah tugas (macOS dan Windows)

## Penggunaan Dasar

### Membuat Layanan

Pertama, inisialisasikan layanan dock:

```go
import "github.com/wailsapp/wails/v3/pkg/application"
import "github.com/wailsapp/wails/v3/pkg/services/dock"

// Create a new Dock service
dockService := dock.New()

// Register the service with the application
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(dockService),
    },
})
```

### Membuat Layanan dengan Opsi Lencana Khusus (Khusus Windows)

Di Windows, Anda dapat menyesuaikan tampilan lencana dengan berbagai opsi:

```go
import "github.com/wailsapp/wails/v3/pkg/application"
import "github.com/wailsapp/wails/v3/pkg/services/dock"
import "image/color"

// Create a dock service with custom badge options
options := dock.BadgeOptions{
    TextColour:       color.RGBA{255, 255, 255, 255}, // White text
    BackgroundColour: color.RGBA{0, 0, 255, 255},     // Blue background
    FontName:         "consolab.ttf",                 // Bold Consolas font
    FontSize:         20,                             // Font size for single character
    SmallFontSize:    14,                             // Font size for multiple characters
}

dockService := dock.NewWithOptions(options)

// Register the service with the application
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(dockService),
    },
})
```

## Operasi Dock

### Menyembunyikan ikon aplikasi di dock

Sembunyikan ikon aplikasi dari Dock macOS:

```go
// Hide the app icon
dockService.HideAppIcon()
```

### Menampilkan ikon aplikasi di dock

Tampilkan ikon aplikasi di Dock macOS:

```go
// Show the app icon
dockService.ShowAppIcon()
```

## Operasi Lencana

### Menetapkan Lencana

Tetapkan lencana pada ubin aplikasi/ikon dock:

```go
// Set a default badge
dockService.SetBadge("")

// Set a numeric badge
dockService.SetBadge("3")

// Set a text badge
dockService.SetBadge("New")
```

### Menetapkan Lencana Khusus (Khusus Windows)

Tetapkan lencana dengan opsi sekali pakai yang diterapkan:

```go
options := dock.BadgeOptions{
    BackgroundColour: color.RGBA{0, 255, 255, 255},
    FontName:         "arialb.ttf", // System font
    FontSize:         16,
    SmallFontSize:    10,
    TextColour:       color.RGBA{0, 0, 0, 255},
}

// Set a default badge
dockService.SetCustomBadge("", options)

// Set a numeric badge
dockService.SetCustomBadge("3", options)

// Set a text badge
dockService.SetCustomBadge("New", options)
```

### Menghapus Lencana

Hapus lencana dari ikon aplikasi:

```go
dockService.RemoveBadge()
```

### Mendapatkan lencana yang ditetapkan

```go
dockService.GetBadge()
```

## Pertimbangan Platform

@tabs
[macOS]
Di macOS:

- Ikon dock dapat **disembunyikan** dan **ditampilkan**
- Lencana ditampilkan langsung pada ikon dock
- Opsi lencana **tidak dapat disesuaikan** (semua opsi yang diteruskan ke `NewWithOptions`/`SetCustomBadge` akan diabaikan)
- Gaya lencana dock standar macOS digunakan dan secara otomatis menyesuaikan dengan tampilan
- Teks label yang meluap ditangani oleh sistem
- Memberikan label kosong akan menampilkan lencana default "●"

[Windows]
Di Windows:

- Layanan ini saat ini tidak mendukung penyembunyian/penampilan ikon bilah tugas
- Lencana ditampilkan sebagai ikon overlay di bilah tugas
- Lencana mendukung nilai teks
- Tampilan lencana dapat disesuaikan melalui `BadgeOptions`
- Aplikasi harus memiliki jendela agar lencana dapat ditampilkan
- Ukuran font yang lebih kecil digunakan secara otomatis untuk label dengan beberapa karakter
- Teks label yang meluap tidak ditangani
- Opsi penyesuaian:
  - **TextColour**: Warna teks (default: putih)
  - **BackgroundColour**: Warna latar belakang lencana (default: merah)
  - **FontName**: Nama file font (default: "segoeuib.ttf")
  - **FontSize**: Ukuran font untuk satu karakter (default: 18)
  - **SmallFontSize**: Ukuran font untuk beberapa karakter (default: 14)


[Linux]
Di Linux:

- Fungsi visibilitas ikon dock dan lencana tidak tersedia

@end

## Praktik Terbaik

1. **Saat menyembunyikan ikon dock (macOS):**
  - Pastikan pengguna tetap dapat mengakses aplikasi Anda (misalnya melalui [baki sistem](/features/menus/systray/))
  - Sertakan opsi "Keluar" dalam antarmuka alternatif Anda
  - Aplikasi tidak akan muncul di pengalih Command+Tab
  - Jendela yang terbuka tetap terlihat dan berfungsi
  - Menutup semua jendela mungkin tidak menghentikan aplikasi (perilaku macOS bervariasi)
  - Pengguna kehilangan cara standar untuk keluar melalui klik kanan pada Dock


2. **Gunakan lencana secukupnya:**
  - Terlalu banyak pembaruan lencana dapat mengganggu pengguna
  - Gunakan lencana hanya untuk notifikasi penting


3. **Gunakan teks lencana yang singkat:**
  - Lencana numerik paling efektif
  - Di macOS, teks lencana sebaiknya singkat


4. **Untuk penyesuaian lencana Windows:**
  - Pastikan kontras yang tinggi antara warna teks dan latar belakang
  - Uji dengan panjang teks yang berbeda karena ukuran font mengecil seiring bertambahnya panjang teks
  - Gunakan font sistem yang umum untuk memastikan ketersediaannya


## Referensi API

### Pengelolaan Layanan

| Metode | Deskripsi |
| --- | --- |
| `New()` | Membuat layanan Dock baru |
| `NewWithOptions(options BadgeOptions)` | Membuat layanan Dock baru dengan opsi badge khusus (khusus Windows; opsi diabaikan di macOS dan Linux) |

### Operasi Dock

| Metode | Deskripsi |
| --- | --- |
| `HideAppIcon()` | Menyembunyikan ikon aplikasi dari Dock macOS (khusus macOS) |
| `ShowAppIcon()` | Menampilkan ikon aplikasi di Dock macOS (khusus macOS) |

### Operasi Badge

| Metode | Deskripsi |
| --- | --- |
| `SetBadge(label string) error` | Menetapkan badge dengan label yang ditentukan |
| `SetCustomBadge(label string, options BadgeOptions) error` | Menetapkan badge dengan label dan opsi gaya khusus yang ditentukan (khusus Windows) |
| `RemoveBadge() error` | Menghapus badge dari ikon aplikasi |
| `GetBadge() *string` | Mendapatkan badge saat ini |

### Struct dan Tipe

```go
// Options for customizing badge appearance (Windows only)
type BadgeOptions struct {
    TextColour       color.RGBA  // Color of the badge text
    BackgroundColour color.RGBA  // Color of the badge background
    FontName         string      // Font file name (e.g., "segoeuib.ttf")
    FontSize         int         // Font size for single character
    SmallFontSize    int         // Font size for multiple characters
}
```
