---
title: "API Jendela"
description: "Referensi lengkap untuk API Jendela"
slug: "reference/window"
sourcePath: "reference/window.md"
---

## Ikhtisar

API Jendela menyediakan berbagai metode untuk mengontrol tampilan, perilaku, dan siklus hidup jendela. Akses API ini melalui instans jendela atau pengelola `app.Window`.

`Window` adalah antarmuka yang diimplementasikan oleh `*application.WebviewWindow`; tanda tangan metode di bawah ini terdapat pada `*WebviewWindow`. Banyak metode pengubah mengembalikan `Window` agar dapat dirangkai—nilai kembali didokumentasikan untuk setiap metode.

**Operasi umum:**

- Membuat dan menampilkan jendela
- Mengontrol ukuran, posisi, dan status
- Menangani peristiwa jendela
- Mengelola konten jendela
- Mengonfigurasi tampilan dan perilaku

## Visibilitas

### Show()

Menampilkan jendela. Jika sebelumnya disembunyikan, jendela akan terlihat. Mengembalikan penerima agar dapat dirangkai.

```go
func (w *WebviewWindow) Show() Window
```

**Contoh:**

```go
window := app.Window.New()
window.Show()
```

### Hide()

Menyembunyikan jendela tanpa menutupnya. Jendela tetap berada dalam memori dan dapat ditampilkan kembali. Mengembalikan penerima agar dapat dirangkai.

```go
func (w *WebviewWindow) Hide() Window
```

**Contoh:**

```go
// Hide window temporarily
window.Hide()

// Show it again later
window.Show()
```

**Kasus penggunaan:**

- Aplikasi baki sistem yang disembunyikan ke baki
- Alur wizard yang menggunakan kembali jendela
- Menyembunyikan jendela sementara selama operasi

### Close()

Menutup jendela. Tindakan ini memicu peristiwa `WindowClosing`.

```go
func (w *WebviewWindow) Close()
```

**Contoh:**

```go
window.Close()
```

**Catatan:** Jika hook yang terdaftar memanggil `event.Cancel()`, penutupan akan dicegah.

## Properti Jendela

### SetTitle()

Mengatur teks bilah judul jendela. Mengembalikan penerima agar dapat dirangkai.

```go
func (w *WebviewWindow) SetTitle(title string) Window
```

**Parameter:**

- `title` - Judul jendela baru

**Contoh:**

```go
window.SetTitle("My Application - Document.txt")
```

### Name()

Mengembalikan pengenal nama unik jendela.

```go
func (w *WebviewWindow) Name() string
```

**Contoh:**

```go
name := window.Name()
fmt.Println("Window name:", name)

// Retrieve window by name later
if w, ok := app.Window.GetByName(name); ok {
    w.Focus()
}
```

## Ukuran dan Posisi

### SetSize()

Mengatur dimensi jendela dalam piksel. Mengembalikan penerima agar dapat dirangkai.

```go
func (w *WebviewWindow) SetSize(width, height int) Window
```

**Parameter:**

- `width` - Lebar jendela dalam piksel
- `height` - Tinggi jendela dalam piksel

**Contoh:**

```go
window.SetSize(1024, 768)
```

### Size()

Mengembalikan dimensi jendela saat ini.

```go
func (w *WebviewWindow) Size() (width, height int)
```

**Contoh:**

```go
width, height := window.Size()
fmt.Printf("Window is %dx%d\n", width, height)
```

### SetMinSize() / SetMaxSize()

Mengatur dimensi minimum dan maksimum jendela. Keduanya mengembalikan penerima agar dapat dirangkai.

```go
func (w *WebviewWindow) SetMinSize(width, height int) Window
func (w *WebviewWindow) SetMaxSize(width, height int) Window
```

**Contoh:**

```go
// Prevent window from being too small
window.SetMinSize(800, 600)

// Prevent window from being too large
window.SetMaxSize(1920, 1080)
```

### SetPosition()

Mengatur posisi jendela relatif terhadap sudut kiri atas layar.

```go
func (w *WebviewWindow) SetPosition(x, y int)
```

**Parameter:**

- `x` - Posisi horizontal dalam piksel
- `y` - Posisi vertikal dalam piksel

**Contoh:**

```go
// Position window at top-left
window.SetPosition(0, 0)

// Position window 100px from top-left
window.SetPosition(100, 100)
```

### Position()

Mengembalikan posisi jendela saat ini.

```go
func (w *WebviewWindow) Position() (x, y int)
```

**Contoh:**

```go
x, y := window.Position()
fmt.Printf("Window is at (%d, %d)\n", x, y)
```

### Center()

Menempatkan jendela di tengah layar.

```go
func (w *WebviewWindow) Center()
```

**Contoh:**

```go
window := app.Window.New()
window.Center()
window.Show()
```

**Catatan:** Jendela ditempatkan di tengah monitor utama. Untuk konfigurasi multi-monitor, lihat API layar.

### Focus()

Membawa jendela ke depan dan memberinya fokus papan ketik.

```go
func (w *WebviewWindow) Focus()
```

**Contoh:**

```go
// Bring window to front
window.Focus()
```

## Status Jendela

### Minimise() / UnMinimise()

Meminimalkan jendela ke bilah tugas/dock atau memulihkannya. `Minimise()` mengembalikan penerima agar dapat dirangkai; `UnMinimise()` tidak mengembalikan apa pun.

```go
func (w *WebviewWindow) Minimise() Window
func (w *WebviewWindow) UnMinimise()
```

**Contoh:**

```go
// Minimise window
window.Minimise()

// Restore from minimised state
window.UnMinimise()
```

### Maximise() / UnMaximise()

Memaksimalkan jendela hingga memenuhi layar atau memulihkannya ke ukuran sebelumnya. `Maximise()` mengembalikan penerima agar dapat dirangkai; `UnMaximise()` tidak mengembalikan apa pun.

```go
func (w *WebviewWindow) Maximise() Window
func (w *WebviewWindow) UnMaximise()
```

**Contoh:**

```go
// Maximise window
window.Maximise()

// Restore to previous size
window.UnMaximise()
```

### Fullscreen() / UnFullscreen() / ToggleFullscreen()

Memasuki atau keluar dari mode layar penuh. `Fullscreen()` mengembalikan penerima agar dapat dirangkai.

```go
func (w *WebviewWindow) Fullscreen() Window
func (w *WebviewWindow) UnFullscreen()
func (w *WebviewWindow) ToggleFullscreen()
```

**Contoh:**

```go
// Enter fullscreen
window.Fullscreen()

// Exit fullscreen
window.UnFullscreen()

// Or toggle
window.ToggleFullscreen()
```

Tidak ada metode `SetFullscreen(bool)`.

### IsMinimised() / IsMaximised() / IsFullscreen()

Memeriksa status jendela saat ini.

```go
func (w *WebviewWindow) IsMinimised() bool
func (w *WebviewWindow) IsMaximised() bool
func (w *WebviewWindow) IsFullscreen() bool
```

**Contoh:**

```go
if window.IsMinimised() {
    window.UnMinimise()
}

if window.IsMaximised() {
    fmt.Println("Window is maximised")
}

if window.IsFullscreen() {
    window.UnFullscreen()
}
```

## Konten Jendela

### SetURL()

Membuka URL tertentu di dalam jendela. Mengembalikan receiver agar pemanggilan metode dapat dirangkai.

```go
func (w *WebviewWindow) SetURL(url string) Window
```

**Parameter:**

- `url` - URL yang akan dibuka (dapat berupa `http://wails.localhost/` untuk aset tersemat)

**Contoh:**

```go
// Navigate to embedded page
window.SetURL("http://wails.localhost/settings.html")

// Navigate to external URL (if allowed)
window.SetURL("https://wails.io")
```

### SetHTML()

Mengatur konten jendela secara langsung dari string HTML. Mengembalikan receiver agar pemanggilan metode dapat dirangkai.

```go
func (w *WebviewWindow) SetHTML(html string) Window
```

**Parameter:**

- `html` - Konten HTML yang akan ditampilkan

**Contoh:**

```go
html := `
<!DOCTYPE html>
<html>
<head><title>Dynamic Content</title></head>
<body>
    <h1>Hello from Go!</h1>
    <p>This content was generated dynamically.</p>
</body>
</html>
`
window.SetHTML(html)
```

**Kasus penggunaan:**

- Pembuatan konten dinamis
- Jendela sederhana tanpa proses build frontend
- Halaman kesalahan atau layar pembuka

### Reload()

Memuat ulang konten jendela saat ini.

```go
func (w *WebviewWindow) Reload()
```

**Contoh:**

```go
// Reload current page
window.Reload()
```

**Catatan:** Berguna selama pengembangan atau ketika konten perlu dimuat ulang.

## Peristiwa Jendela

Wails menyediakan dua metode untuk menangani peristiwa jendela:

- **OnWindowEvent()** - Mendengarkan peristiwa jendela (tidak dapat mencegahnya).
- **RegisterHook()** - Memasang hook pada peristiwa jendela (dapat mencegahnya dengan memanggil `event.Cancel()`).

### OnWindowEvent()

Mendaftarkan callback untuk peristiwa jendela. Mengembalikan fungsi untuk berhenti berlangganan.

```go
func (w *WebviewWindow) OnWindowEvent(
    eventType events.WindowEventType,
    callback func(event *WindowEvent),
) func()
```

**Contoh:**

```go
import "github.com/wailsapp/wails/v3/pkg/events"

// Listen for window focus
window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
    app.Logger.Info("Window gained focus")
})

// Listen for window lost focus
window.OnWindowEvent(events.Common.WindowLostFocus, func(e *application.WindowEvent) {
    app.Logger.Info("Window lost focus")
})

// Listen for window resize
window.OnWindowEvent(events.Common.WindowDidResize, func(e *application.WindowEvent) {
    app.Logger.Info("Window resized")
})
```

**Peristiwa Jendela yang Umum:**

- `events.Common.WindowClosing` - Jendela akan ditutup
- `events.Common.WindowFocus` - Jendela mendapatkan fokus
- `events.Common.WindowLostFocus` - Jendela kehilangan fokus
- `events.Common.WindowDidMove` - Jendela dipindahkan
- `events.Common.WindowDidResize` - Ukuran jendela diubah
- `events.Common.WindowMinimise` - Jendela diminimalkan
- `events.Common.WindowMaximise` - Jendela dimaksimalkan
- `events.Common.WindowFullscreen` - Jendela memasuki mode layar penuh
- `events.Common.WindowRuntimeReady` - Runtime di dalam jendela diinisialisasi

### RegisterHook()

Mendaftarkan hook untuk peristiwa jendela. Hook dijalankan sebelum listener dan dapat mencegah peristiwa dengan memanggil `event.Cancel()`. Mengembalikan fungsi untuk berhenti berlangganan.

```go
func (w *WebviewWindow) RegisterHook(
    eventType events.WindowEventType,
    callback func(event *WindowEvent),
) func()
```

**Contoh - Mencegah jendela ditutup:**

```go
import "github.com/wailsapp/wails/v3/pkg/events"

window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    confirm := app.Dialog.Question().
        SetTitle("Confirm Close").
        SetMessage("Are you sure you want to close this window?")

    yes := confirm.AddButton("Yes")
    no := confirm.AddButton("No")
    confirm.SetDefaultButton(yes)
    confirm.SetCancelButton(no)

    no.OnClick(func() {
        e.Cancel() // Prevent window from closing
    })

    confirm.Show()
})
```

**Contoh - Menyimpan sebelum menutup:**

```go
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    if !hasUnsavedChanges {
        return
    }

    dlg := app.Dialog.Question().
        SetTitle("Unsaved Changes").
        SetMessage("Save changes before closing?")

    save := dlg.AddButton("Save")
    discard := dlg.AddButton("Don't Save")
    cancel := dlg.AddButton("Cancel")
    dlg.SetDefaultButton(save)
    dlg.SetCancelButton(cancel)

    save.OnClick(func() { saveData() })
    cancel.OnClick(func() { e.Cancel() })
    _ = discard // allow close

    dlg.Show()
})
```

### EmitEvent()

Memancarkan peristiwa khusus ke frontend jendela. Mengembalikan `true` jika pemancaran dibatalkan oleh hook.

```go
func (w *WebviewWindow) EmitEvent(name string, data ...any) bool
```

**Parameter:**

- `name` - Nama peristiwa
- `data` - Data opsional yang akan dikirim bersama peristiwa

**Contoh:**

```go
// Send data to specific window
window.EmitEvent("data-updated", map[string]any{
    "count":  42,
    "status": "success",
})
```

**Frontend (JavaScript):**

```javascript
import { Events } from '@wailsio/runtime'

Events.On('data-updated', (data) => {
    console.log('Count:', data.count)
    console.log('Status:', data.status)
})
```

## Metode Lainnya

### SetEnabled()

Mengaktifkan atau menonaktifkan interaksi pengguna dengan jendela.

```go
func (w *WebviewWindow) SetEnabled(enabled bool)
```

**Contoh:**

```go
// Disable window during long operation
window.SetEnabled(false)

// Perform operation
performLongOperation()

// Re-enable window
window.SetEnabled(true)
```

### SetBackgroundColour()

Mengatur warna latar belakang jendela (ditampilkan sebelum konten dimuat). Mengembalikan receiver agar pemanggilan metode dapat dirangkai.

```go
func (w *WebviewWindow) SetBackgroundColour(colour RGBA) Window
```

`RGBA` adalah `application.RGBA{Red, Green, Blue, Alpha uint8}`. Gunakan fungsi pembantu `application.NewRGB(r, g, b)` (alfa 255) atau `application.NewRGBA(r, g, b, a)`.

**Contoh:**

```go
// White background
window.SetBackgroundColour(application.NewRGB(255, 255, 255))

// Dark background with full alpha
window.SetBackgroundColour(application.NewRGBA(30, 30, 30, 255))
```

### SetResizable()

Mengatur apakah ukuran jendela dapat diubah oleh pengguna. Mengembalikan receiver agar pemanggilan metode dapat dirangkai.

```go
func (w *WebviewWindow) SetResizable(resizable bool) Window
```

**Contoh:**

```go
// Make window fixed size
window.SetResizable(false)
```

### SetAlwaysOnTop()

Mengatur apakah jendela tetap berada di atas jendela lain. Mengembalikan receiver agar pemanggilan metode dapat dirangkai.

```go
func (w *WebviewWindow) SetAlwaysOnTop(alwaysOnTop bool) Window
```

**Contoh:**

```go
// Keep window on top
window.SetAlwaysOnTop(true)
```

### Print()

Membuka dialog cetak native untuk konten jendela.

```go
func (w *WebviewWindow) Print() error
```

**Mengembalikan:** Error jika pencetakan gagal.

**Contoh:**

```go
if err := window.Print(); err != nil {
    log.Println("Print failed:", err)
}
```

### AttachModal()

Melampirkan jendela kedua sebagai modal lembar.

```go
func (w *WebviewWindow) AttachModal(modalWindow Window)
```

**Parameter:**

- `modalWindow` - Jendela yang akan dilampirkan sebagai modal

**Dukungan platform:**

- **macOS**: Dukungan penuh (ditampilkan sebagai lembar)
- **Windows**: Tidak didukung
- **Linux**: Tidak didukung

**Contoh:**

```go
modalWindow := app.Window.New()
window.AttachModal(modalWindow)
```

## Opsi Khusus Platform

### Linux

Jendela Linux mendukung opsi khusus platform berikut melalui `LinuxWindow`:

#### MenuStyle

Mengontrol cara menu aplikasi ditampilkan. Opsi ini tersedia pada build GTK4 default dan diabaikan pada build `-tags gtk3` lama.

| Nilai | Deskripsi |
| --- | --- |
| `LinuxMenuStyleMenuBar` | Bilah menu tradisional di bawah bilah judul (default) |
| `LinuxMenuStylePrimaryMenu` | Tombol menu utama pada bilah header (gaya GNOME) |

**Contoh:**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "My Application",
    Linux: application.LinuxWindow{
        MenuStyle: application.LinuxMenuStylePrimaryMenu,
    },
})
window.SetMenu(menu)
```

**Catatan:** Gaya menu utama menampilkan tombol hamburger (☰) pada bilah header, sesuai dengan Panduan Antarmuka Manusia GNOME. Ini adalah gaya yang direkomendasikan untuk aplikasi GNOME modern.

## Contoh Lengkap

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)

func main() {
    app := application.New(application.Options{
        Name: "Window API Demo",
    })

    // Create window with options
    window := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:            "My Application",
        Width:            1024,
        Height:           768,
        MinWidth:         800,
        MinHeight:        600,
        BackgroundColour: application.NewRGB(255, 255, 255),
        URL:              "http://wails.localhost/",
    })

    // Configure window behaviour
    window.SetResizable(true)
    window.SetMinSize(800, 600)
    window.SetMaxSize(1920, 1080)

    // Confirm-before-close hook
    window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
        dlg := app.Dialog.Question().
            SetTitle("Confirm Close").
            SetMessage("Are you sure you want to close this window?")

        yes := dlg.AddButton("Yes")
        no := dlg.AddButton("No")
        dlg.SetDefaultButton(yes)
        dlg.SetCancelButton(no)
        no.OnClick(func() { e.Cancel() })

        dlg.Show()
    })

    // Listen for window events
    window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
        window.SetTitle("My Application (Active)")
        app.Logger.Info("Window gained focus")
    })

    window.OnWindowEvent(events.Common.WindowLostFocus, func(e *application.WindowEvent) {
        window.SetTitle("My Application")
        app.Logger.Info("Window lost focus")
    })

    // Position and show window
    window.Center()
    window.Show()

    app.Run()
}
```
