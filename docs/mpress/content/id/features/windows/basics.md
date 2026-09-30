---
title: "Dasar-Dasar Jendela"
description: "Membuat dan mengelola jendela aplikasi di Wails"
slug: "features/windows/basics"
sourcePath: "features/windows/basics.md"
---

## Pengelolaan Jendela

Wails menyediakan **API pengelolaan jendela terpadu** yang berfungsi di semua platform. Buat jendela, kendalikan perilakunya, dan kelola beberapa jendela dengan kendali penuh atas pembuatan, tampilan, perilaku, dan siklus hidupnya.

## Mulai Cepat

```go
package main

import "github.com/wailsapp/wails/v3/pkg/application"

func main() {
    app := application.New(application.Options{
        Name: "My App",
    })

    // Create a window
    window := app.Window.New()
    
    // Configure it
    window.SetTitle("Hello Wails")
    window.SetSize(800, 600)
    window.Center()
    
    // Show it
    window.Show()

    app.Run()
}
```

**Selesai!** Anda kini memiliki jendela lintas platform.

## Membuat Jendela

### Jendela Dasar

Cara paling sederhana untuk membuat jendela:

```go
window := app.Window.New()
```

**Yang Anda dapatkan:**

- Ukuran default (800x600)
- Judul default (nama aplikasi)
- WebView yang siap digunakan oleh frontend Anda
- Tampilan native platform

### Jendela dengan Opsi

Buat jendela dengan konfigurasi khusus:

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "My Application",
    Width:  1200,
    Height: 800,
    X:      100,   // Position from left
    Y:      100,   // Position from top
    AlwaysOnTop: false,
    Frameless: false,
    Hidden: false,
    MinWidth: 400,
    MinHeight: 300,
    MaxWidth: 1920,
    MaxHeight: 1080,
})
```

**Opsi umum:**

| Opsi | Tipe | Deskripsi |
| --- | --- | --- |
| `Title` | `string` | Judul jendela |
| `Width` | `int` | Lebar jendela dalam piksel |
| `Height` | `int` | Tinggi jendela dalam piksel |
| `X` | `int` | Posisi X (dari kiri) |
| `Y` | `int` | Posisi Y (dari atas) |
| `AlwaysOnTop` | `bool` | Pertahankan jendela di atas jendela lain |
| `Frameless` | `bool` | Hapus bilah judul dan bingkai |
| `Hidden` | `bool` | Mulai dalam keadaan tersembunyi |
| `MinWidth` | `int` | Lebar minimum |
| `MinHeight` | `int` | Tinggi minimum |
| `MaxWidth` | `int` | Lebar maksimum |
| `MaxHeight` | `int` | Tinggi maksimum |

**Lihat [Opsi Jendela](/features/windows/options/) untuk daftar lengkap.**

### Jendela Bernama

Beri nama pada jendela agar mudah ditemukan kembali:

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:  "main-window",
    Title: "Main Application",
})

// Later, find it by name
if mainWindow, ok := app.Window.GetByName("main-window"); ok {
    mainWindow.Show()
}
```

**Kasus penggunaan:**

- Beberapa jendela (utama, pengaturan, tentang)
- Menemukan jendela dari berbagai bagian kode Anda
- Komunikasi antarjendela

## Mengendalikan Jendela

### Menampilkan dan Menyembunyikan

```go
// Show window
window.Show()

// Hide window
window.Hide()

// Check if visible
if window.IsVisible() {
    fmt.Println("Window is visible")
}
```

**Kasus penggunaan:**

- Layar pembuka (tampilkan, lalu sembunyikan)
- Jendela pengaturan (sembunyikan saat tidak diperlukan)
- Jendela pop-up (tampilkan sesuai kebutuhan)

### Posisi dan Ukuran

```go
// Set size
window.SetSize(1024, 768)

// Set position
window.SetPosition(100, 100)

// Centre on screen
window.Center()

// Get current size
width, height := window.Size()

// Get current position
x, y := window.Position()
```

**Sistem koordinat:**

- (0, 0) adalah sudut kiri atas layar utama
- Nilai X positif mengarah ke kanan
- Nilai Y positif mengarah ke bawah

### Status Jendela

```go
// Minimise
window.Minimise()

// Maximise
window.Maximise()

// Fullscreen
window.Fullscreen()

// Restore to normal
window.Restore()

// Check state
if window.IsMinimised() {
    fmt.Println("Window is minimised")
}

if window.IsMaximised() {
    fmt.Println("Window is maximised")
}

if window.IsFullscreen() {
    fmt.Println("Window is fullscreen")
}
```

**Transisi status:**

```
Normal ←→ Minimised
Normal ←→ Maximised
Normal ←→ Fullscreen
```

### Judul dan Tampilan

```go
// Set title
window.SetTitle("My Application - Document.txt")

// Set background colour — RGBA value (helper for RGB)
window.SetBackgroundColour(application.NewRGBA(0, 0, 0, 255))

// Set always on top
window.SetAlwaysOnTop(true)

// Set resizable
window.SetResizable(false)
```

### Menutup Jendela

```go
// Close window — dispatches WindowClosing; a RegisterHook can call e.Cancel().
window.Close()
```

Tidak ada metode `window.Destroy()` di v3 — gunakan `Close()` dan dengarkan dengan `OnWindowEvent` (tidak dapat membatalkan penutupan) atau pasang hook dengan `RegisterHook` (dapat memanggil `e.Cancel()` agar jendela tetap terbuka).

## Menemukan Jendela

### Berdasarkan Nama

```go
if window, ok := app.Window.GetByName("settings"); ok {
    window.Show()
}
```

### Berdasarkan ID

Setiap jendela memiliki ID unik:

```go
id := window.ID()
fmt.Printf("Window ID: %d\n", id)

// Find by ID
if found, ok := app.Window.GetByID(id); ok {
    found.Focus()
}
```

### Jendela Saat Ini

Dapatkan jendela yang sedang memiliki fokus:

```go
current := app.Window.Current()
if current != nil {
    current.SetTitle("Active Window")
}
```

### Semua Jendela

Dapatkan semua jendela:

```go
windows := app.Window.GetAll()
fmt.Printf("Total windows: %d\n", len(windows))

for _, w := range windows {
    fmt.Printf("Window: %s (ID: %d)\n", w.Name(), w.ID())
}
```

## Siklus Hidup Jendela

### Pembuatan

```go
app.Window.OnCreate(func(window application.Window) {
    fmt.Printf("Window created: %s\n", window.Name())

    // Configure new windows
    window.SetMinSize(400, 300)
})
```

### Penutupan

Untuk mencegah jendela ditutup, gunakan `RegisterHook` dengan peristiwa `WindowClosing`:

```go
window.RegisterHook(events.Common.WindowClosing, func(event *application.WindowEvent) {
    if hasUnsavedChanges() {
        // Ask user for confirmation
        result := showConfirmDialog("Unsaved changes. Close anyway?")
        if result != "yes" {
            // Cancel the close event
            event.Cancel()
        }
    }
})
```

**Penting:** `RegisterHook` mencegat peristiwa penutupan sebelum terjadi. Panggil `event.Cancel()` untuk mencegah jendela ditutup. Ini berlaku untuk penutupan yang dimulai oleh pengguna (mengeklik tombol X).

### Pemusnahan

Untuk melakukan pembersihan saat jendela ditutup, gunakan `OnWindowEvent` dengan peristiwa `WindowClosing`:

```go
window.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
    fmt.Println("Window is closing")
    // Cleanup resources
})
```

## Beberapa Jendela

### Membuat Beberapa Jendela

```go
// Main window
mainWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:   "main",
    Title:  "Main Application",
    Width:  1200,
    Height: 800,
})

// Settings window
settingsWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:   "settings",
    Title:  "Settings",
    Width:  600,
    Height: 400,
    Hidden: true,  // Start hidden
})

// Show settings when needed
settingsWindow.Show()
```

### Komunikasi Antarjendela

Jendela dapat berkomunikasi melalui peristiwa:

```go
// In main window
app.Event.Emit("data-updated", map[string]interface{}{
    "value": 42,
})

// In settings window
app.Event.On("data-updated", func(event *application.CustomEvent) {
    data := event.Data.(map[string]interface{})
    value := data["value"].(int)
    fmt.Printf("Received: %d\n", value)
})
```

**Lihat [Peristiwa](/features/events/system/) untuk informasi selengkapnya.**

### Jendela Induk-Anak

`WebviewWindowOptions` tidak memiliki bidang `Parent`. Buat jendela anak sebagai jendela biasa, lalu kaitkan ke induk sebagai lembar modal:

```go
// Create child window
childWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Child Window",
})

// Attach to the parent — presents as a sheet on macOS.
mainWindow.AttachModal(childWindow)
```

**Perilaku:**

- Jendela anak tetap berada di atas jendela induk.
- Jendela anak bersifat modal — memblokir interaksi dengan jendela induk.

**Dukungan platform:**

- **macOS:** Dukungan penuh (ditampilkan sebagai lembar).
- **Windows:** Tidak didukung.
- **Linux:** Tidak didukung.

## Fitur Khusus Platform

@tabs{sync-key="platform"}
[Windows]
**Fitur khusus Windows:**

```go
// Flash taskbar button
window.Flash(true)  // Start flashing
window.Flash(false) // Stop flashing

// Trigger Windows 11 Snap Assist (Win+Z)
window.SnapAssist()
```

Tidak ada `SetIcon` per jendela — ikon aplikasi ditetapkan pada aplikasi melalui `app.SetIcon([]byte)` (atau, untuk ikon jendela khusus Linux, melalui bidang `application.LinuxWindow.Icon` saat jendela dibuat).

**Snap Assist:** Menampilkan opsi tata letak snap Windows 11 melalui jalur pintasan sistem. Untuk tombol maksimalkan HTML khusus dengan Snap Layouts native saat penunjuk diarahkan ke atasnya, gunakan [Wilayah Nonklien Native di Windows](/features/windows/frameless/#native-non-client-regions-on-windows) sebagai gantinya.

**Kedipan bilah tugas:** Berguna untuk notifikasi saat jendela diminimalkan.

[macOS]
**Fitur khusus macOS:**

```go
// Transparent title bar
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Mac: application.MacWindow{
        TitleBar: application.MacTitleBar{
            AppearsTransparent: true,
        },
        Backdrop: application.MacBackdropTranslucent,
    },
})
```

**Jenis latar belakang:**

- `MacBackdropNormal` - Jendela standar
- `MacBackdropTranslucent` - Latar belakang tembus cahaya; **memerlukan API privat** agar webview transparan.
- `MacBackdropTransparent` - Sepenuhnya transparan; **memerlukan API privat** agar webview transparan.
- `MacBackdropLiquidGlass` - Latar belakang kaca; **memerlukan API privat** agar webview transparan.

Lakukan build dengan `-tags private_mac_apis` agar efek ini terlihat melalui webview. Tanpanya, webview tetap tidak transparan. `TitleBar.AppearsTransparent` sendiri menggunakan API publik. Lihat [API macOS Privat](/guides/build/private-macos-apis/).

**Perilaku koleksi:** Atur perilaku jendela di berbagai Space:

- `MacWindowCollectionBehaviorCanJoinAllSpaces` - Terlihat di semua Space
- `MacWindowCollectionBehaviorFullScreenAuxiliary` - Dapat ditampilkan di atas aplikasi layar penuh

**Layar penuh native:** Mode layar penuh macOS membuat Space (desktop virtual) baru.

[Linux]
**Fitur khusus Linux:**

```go
// Set window icon (per-window struct is LinuxWindow, not the app-level LinuxOptions)
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Linux: application.LinuxWindow{
        Icon: iconBytes,
    },
})
```

**Catatan lingkungan desktop:**

- GNOME: Dukungan penuh
- KDE Plasma: Dukungan penuh
- XFCE: Dukungan sebagian
- Lainnya: Bervariasi

**Pengelola jendela tiling (Hyprland, Sway, i3, dan sebagainya):**

- `Minimise()` dan `Maximise()` mungkin tidak berfungsi sebagaimana mestinya — WM mengendalikan geometri jendela
- Permintaan `SetSize()` dan `SetPosition()` hanya bersifat saran dan dapat diabaikan
- `Fullscreen()` biasanya berfungsi sesuai harapan
- Beberapa WM tidak mendukung mode selalu di atas

@end

## Pola Umum

### Layar Pembuka

```go
// Create splash screen
splash := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:     "Loading...",
    Width:     400,
    Height:    300,
    Frameless: true,
    AlwaysOnTop: true,
})

// Show splash
splash.Show()

// Initialise application
time.Sleep(2 * time.Second)

// Hide splash, show main window
splash.Close()
mainWindow.Show()
```

### Jendela Pengaturan

```go
var settingsWindow *application.WebviewWindow

func showSettings() {
    if settingsWindow == nil {
        settingsWindow = app.Window.NewWithOptions(application.WebviewWindowOptions{
            Name:   "settings",
            Title:  "Settings",
            Width:  600,
            Height: 400,
        })
    }
    
    settingsWindow.Show()
    settingsWindow.Focus()
}
```

### Konfirmasi Sebelum Menutup

```go
window.RegisterHook(events.Common.WindowClosing, func(event *application.WindowEvent) {
    if hasUnsavedChanges() {
        // Show dialog
        result := showConfirmDialog("Unsaved changes. Close anyway?")
        if result != "yes" {
            // Cancel the close event
            event.Cancel()
        }
    }
})
```

## Praktik Terbaik

### ✅ Lakukan

- **Beri nama pada jendela penting** — Agar lebih mudah ditemukan nanti
- **Tetapkan ukuran minimum** — Mencegah tata letak yang tidak dapat digunakan
- **Posisikan jendela di tengah** — Memberikan pengalaman pengguna yang lebih baik daripada posisi acak
- **Tangani peristiwa penutupan** — Mencegah kehilangan data
- **Uji di semua platform** — Perilakunya berbeda-beda
- **Gunakan ukuran yang sesuai** — Pertimbangkan berbagai ukuran layar

### ❌ Jangan Lakukan

- **Jangan membuat terlalu banyak jendela** — Membingungkan pengguna
- **Jangan lupa menutup jendela** — Menyebabkan kebocoran memori
- **Jangan menetapkan posisi secara hardcode** — Ukuran layar berbeda-beda
- **Jangan abaikan perbedaan antarplatform** — Uji secara menyeluruh
- **Jangan memblokir thread UI** — Gunakan goroutine untuk operasi yang berlangsung lama

## Pemecahan Masalah

### Jendela Tidak Muncul

**Kemungkinan penyebab:**

1. Jendela dibuat dalam keadaan tersembunyi
2. Jendela berada di luar layar
3. Jendela berada di belakang jendela lain

**Solusi:**

```go
window.Show()
window.Center()
window.Focus()
```

### Ukuran Jendela Salah

**Penyebab:** penskalaan DPI di Windows/Linux

**Solusi:**

```go
// Wails handles DPI automatically
// Just use logical pixels
window.SetSize(800, 600)
```

### Jendela Langsung Tertutup

**Penyebab:** aplikasi berhenti ketika jendela terakhir ditutup

**Solusi:**

```go
app := application.New(application.Options{
    Mac: application.MacOptions{
        ApplicationShouldTerminateAfterLastWindowClosed: false,
    },
})
```

## Langkah Berikutnya

@cards{cols="2"}
⚙ Opsi Jendela
Referensi lengkap untuk semua opsi jendela.

[Pelajari Selengkapnya →](/features/windows/options/)

---
▣ Beberapa Jendela
Pola untuk aplikasi multijendela.

[Pelajari Selengkapnya →](/features/windows/multiple/)

---
★ Jendela Tanpa Bingkai
Buat bingkai jendela khusus.

[Pelajari Selengkapnya →](/features/windows/frameless/)

---
🚀 Peristiwa Jendela
Tangani peristiwa siklus hidup jendela.

[Pelajari Selengkapnya →](/features/windows/events/)

@end

---

**Ada pertanyaan?** Tanyakan di [Discord](https://discord.gg/JDdSxwjhGf) atau lihat [contoh jendela](https://github.com/wailsapp/wails/tree/master/v3/examples).
