---
title: "Menu Baki Sistem"
description: "Tambahkan integrasi baki sistem (area notifikasi) ke aplikasi Anda"
slug: "features/menus/systray"
sourcePath: "features/menus/systray.md"
---

## Menu Baki Sistem

Wails menyediakan **API baki sistem terpadu** yang berfungsi di semua platform. Buat ikon baki dengan menu, lampirkan jendela, dan tangani klik dengan perilaku native platform untuk aplikasi latar belakang, layanan, dan utilitas akses cepat.

![Menu baki sistem Wails yang dibuka dari bar menu macOS](/assets/screenshots/systray-menu-macos.png)

Di macOS, item baki sistem Wails muncul di bar menu dan membuka menu native. Contoh ini mencakup item yang dinonaktifkan, kotak centang, tombol radio, submenu, dan tindakan.

## Mulai Cepat

```go
package main

import (
    _ "embed"
    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets/icon.png
var icon []byte

func main() {
    app := application.New(application.Options{
        Name: "Tray App",
    })

    // Create system tray
    systray := app.SystemTray.New()
    systray.SetIcon(icon)
    systray.SetLabel("My App")

    // Add menu
    menu := app.NewMenu()
    menu.Add("Show").OnClick(func(ctx *application.Context) {
        // Show main window
    })
    menu.Add("Quit").OnClick(func(ctx *application.Context) {
        app.Quit()
    })
    systray.SetMenu(menu)

    // Create hidden window
    window := app.Window.New()
    window.Hide()

    app.Run()
}
```

**Hasil:** Ikon baki sistem dengan menu di semua platform.

## Membuat Baki Sistem

### Baki Sistem Dasar

```go
// Create system tray
systray := app.SystemTray.New()

// Set icon
systray.SetIcon(iconBytes)

// Set label (macOS) / tooltip (Windows)
systray.SetLabel("My Application")
```

### Dengan Ikon

Ikon sebaiknya disematkan:

```go
import _ "embed"

//go:embed assets/icon.png
var icon []byte

//go:embed assets/icon-dark.png
var iconDark []byte

func main() {
    app := application.New(application.Options{
        Name: "My App",
    })

    systray := app.SystemTray.New()
    systray.SetIcon(icon)
    systray.SetDarkModeIcon(iconDark)  // Windows and macOS dark mode
    
    app.Run()
}
```

**Persyaratan ikon:**

| Platform | Ukuran | Format | Catatan |
| --- | --- | --- | --- |
| **Windows** | 16x16 atau 32x32 | PNG, ICO | Area notifikasi |
| **macOS** | 18x18 hingga 22x22 | PNG | Bar menu, templat disarankan |
| **Linux** | 22x22 hingga 48x48 | PNG, SVG | Bervariasi menurut DE |

### Ikon Templat (macOS)

Ikon templat secara otomatis menyesuaikan dengan mode terang/gelap:

```go
systray.SetTemplateIcon(iconBytes)
```

**Panduan ikon templat:**

- Gunakan hanya warna hitam dan bening (transparan)
- Warna hitam berubah menjadi putih dalam mode gelap
- Beri nama file dengan akhiran `Template`: `iconTemplate.png`
- [Panduan desain](https://bjango.com/articles/designingmenubarextras/)

## Menambahkan Menu

Menu baki sistem berfungsi seperti menu aplikasi:

```go
menu := app.NewMenu()

// Add items
menu.Add("Open").OnClick(func(ctx *application.Context) {
    showMainWindow()
})

menu.AddSeparator()

menu.AddCheckbox("Start at Login", false).OnClick(func(ctx *application.Context) {
    enabled := ctx.ClickedMenuItem().Checked()
    setStartAtLogin(enabled)
})

menu.AddSeparator()

menu.Add("Quit").OnClick(func(ctx *application.Context) {
    app.Quit()
})

// Set menu
systray.SetMenu(menu)
```

Untuk **semua jenis item menu**, lihat [Referensi Menu](/features/menus/reference/).

## Melampirkan Jendela

Lampirkan jendela ke ikon baki agar ditampilkan/disembunyikan secara otomatis:

```go
// Create window
window := app.Window.New()

// Attach to tray
systray.AttachWindow(window)

// Configure behaviour — these are setters that return the receiver for chaining.
systray.WindowOffset(10)                          // Pixels from tray icon
systray.WindowDebounce(200 * time.Millisecond)    // Click debounce
```

**Perilaku:**

- Jendela dimulai dalam keadaan tersembunyi
- **Klik kiri ikon baki** → Alihkan visibilitas jendela
- **Klik kanan ikon baki** → Tampilkan menu (jika ditetapkan)
- Jendela ditempatkan di dekat ikon baki

**Contoh: Jendela popup**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:           "Quick Access",
    Width:           300,
    Height:          400,
    Frameless:       true, // No title bar
    AlwaysOnTop:     true, // Stay on top
    HideOnFocusLost: true, // Dismiss when another window receives focus
    HideOnEscape:    true, // Dismiss when the user presses Escape
})

systray.AttachWindow(window)
systray.WindowOffset(5)
```

`HideOnFocusLost` berguna untuk popup baki pada Windows, macOS, dan desktop Linux dengan fokus melalui klik. Wails menonaktifkan perilaku tersebut pada lingkungan Linux yang menerapkan fokus mengikuti tetikus (termasuk konfigurasi umum Hyprland, Sway, dan i3), karena meninggalkan popup dalam lingkungan tersebut dapat menyembunyikannya sebelum sempat digunakan. `HideOnEscape` tetap tersedia dalam lingkungan tersebut.

Perilaku klik kiri dan klik kanan di atas merupakan default cerdas. Handler `OnClick` atau `OnRightClick` yang ditetapkan secara eksplisit menggantikan default terkait. Untuk pemeriksaan platform dan kasus khusus, lihat [rangkaian pengujian manual systray](https://github.com/wailsapp/wails/tree/master/v3/test/manual/systray) dan [contoh pengujian beban systray](https://github.com/wailsapp/wails/tree/master/v3/examples/systray-stress).

## Handler Klik

Tangani klik pada ikon baki:

```go
systray := app.SystemTray.New()

// Left click
systray.OnClick(func() {
    fmt.Println("Tray icon clicked")
})

// Right click
systray.OnRightClick(func() {
    fmt.Println("Tray icon right-clicked")
})

// Double click
systray.OnDoubleClick(func() {
    fmt.Println("Tray icon double-clicked")
})

// Mouse enter/leave
systray.OnMouseEnter(func() {
    fmt.Println("Mouse entered tray icon")
})

systray.OnMouseLeave(func() {
    fmt.Println("Mouse left tray icon")
})
```

**Dukungan platform:**

| Peristiwa | Windows | macOS | Linux |
| --- | --- | --- | --- |
| OnClick | ✅ | ✅ | ✅ |
| OnRightClick | ✅ | ✅ | ✅ |
| OnDoubleClick | ✅ | ✅ | ⚠️ Bervariasi |
| OnMouseEnter | ✅ | ✅ | ⚠️ Bervariasi |
| OnMouseLeave | ✅ | ✅ | ⚠️ Bervariasi |

## Pembaruan Dinamis

Perbarui ikon dan menu baki secara dinamis:

### Mengubah Ikon

```go
var isActive bool

func updateTrayIcon() {
    if isActive {
        systray.SetIcon(activeIcon)
        systray.SetLabel("Active")
    } else {
        systray.SetIcon(inactiveIcon)
        systray.SetLabel("Inactive")
    }
}
```

### Memperbarui Menu

```go
var isPaused bool

pauseMenuItem := menu.Add("Pause")

pauseMenuItem.OnClick(func(ctx *application.Context) {
    isPaused = !isPaused
    
    if isPaused {
        pauseMenuItem.SetLabel("Resume")
    } else {
        pauseMenuItem.SetLabel("Pause")
    }
    
    menu.Update()  // Important!
})
```

@note{type="caution" title="Selalu Panggil Update()"}
Setelah mengubah status menu, **panggil `menu.Update()`**. Lihat [Referensi Menu](/features/menus/reference/#enabled-state).

@end

### Membangun Ulang Menu

Untuk perubahan besar, buat ulang seluruh menu:

```go
func rebuildTrayMenu(status string) {
    menu := app.NewMenu()
    
    // Status-specific items
    switch status {
    case "syncing":
        menu.Add("Syncing...").SetEnabled(false)
        menu.Add("Pause Sync").OnClick(pauseSync)
    case "synced":
        menu.Add("Up to date ✓").SetEnabled(false)
        menu.Add("Sync Now").OnClick(startSync)
    case "error":
        menu.Add("Sync Error").SetEnabled(false)
        menu.Add("Retry").OnClick(retrySync)
    }
    
    menu.AddSeparator()
    menu.Add("Quit").OnClick(func(ctx *application.Context) {
        app.Quit()
    })
    
    systray.SetMenu(menu)
}
```

## Fitur Khusus Platform

@tabs{sync-key="platform"}
[macOS]
**Integrasi bilah menu:**

```go
// Set label (appears next to icon)
systray.SetLabel("My App")

// Use template icon (adapts to dark mode)
systray.SetTemplateIcon(iconBytes)

// Set icon position — uses AppKit NSImage placement constants.
systray.SetIconPosition(application.NSImageRight)
```

**Posisi ikon** (mencerminkan `NSImagePosition`):

- `application.NSImageLeft` - Ikon di sebelah kiri label.
- `application.NSImageRight` - Ikon di sebelah kanan label.
- `application.NSImageOnly` - Hanya ikon, tanpa label.
- `application.NSImageNone` - Hanya label, tanpa ikon.

**Praktik terbaik:**

- Gunakan ikon templat (hitam + transparan)
- Gunakan label singkat (3-5 karakter)
- 18x18 hingga 22x22 piksel untuk layar Retina
- Uji dalam mode terang dan gelap

[Windows]
**Integrasi area notifikasi:**

```go
// Set tooltip (appears on hover)
systray.SetTooltip("My Application")

// Or use SetLabel (same as tooltip on Windows)
systray.SetLabel("My Application")

// Show/Hide functionality (fully functional)
systray.Show()  // Show tray icon
systray.Hide()  // Hide tray icon
```

**Persyaratan ikon:**

- 16x16 atau 32x32 piksel
- Format PNG atau ICO
- Latar belakang transparan

**Batas tooltip:**

- Maksimum 127 karakter UTF-16
- Tooltip yang lebih panjang akan dipotong
- Buat tetap ringkas untuk pengalaman terbaik

**Fitur platform:**

- Ikon tray tetap tersedia setelah Windows Explorer dimulai ulang
- Metode Show() dan Hide() berfungsi sepenuhnya
- Pengelolaan siklus hidup yang tepat

**Praktik terbaik:**

- Gunakan 32x32 untuk layar ber-DPI tinggi
- Batasi tooltip hingga kurang dari 127 karakter
- Uji pada berbagai versi Windows
- Pertimbangkan luapan area notifikasi
- Gunakan Show/Hide untuk mengatur visibilitas tray secara kondisional

[Linux]
**Integrasi tray sistem:**

Menggunakan spesifikasi StatusNotifierItem (sebagian besar DE modern).

```go
systray.SetIcon(iconBytes)
systray.SetLabel("My App")
```

**Dukungan lingkungan desktop:**

- **GNOME**: Bilah atas (dengan ekstensi)
- **KDE Plasma**: Tray sistem
- **XFCE**: Area notifikasi
- **Lainnya**: Bervariasi

**Praktik terbaik:**

- Gunakan ukuran 22x22 atau 24x24 piksel
- Ikon SVG dapat diskalakan dengan lebih baik
- Uji pada lingkungan desktop target
- Sediakan opsi fallback untuk DE yang tidak didukung

@end

## Contoh Lengkap

Berikut adalah aplikasi tray sistem yang siap digunakan dalam produksi:

```go
package main

import (
    _ "embed"
    "fmt"
    "time"
    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets/icon.png
var icon []byte

//go:embed assets/icon-active.png
var iconActive []byte

type TrayApp struct {
    app     *application.App
    systray *application.SystemTray
    window  *application.WebviewWindow
    menu    *application.Menu
    isActive bool
}

func main() {
    app := application.New(application.Options{
        Name: "Tray Application",
        Mac: application.MacOptions{
            ApplicationShouldTerminateAfterLastWindowClosed: false,
        },
    })

    trayApp := &TrayApp{app: app}
    trayApp.setup()

    app.Run()
}

func (t *TrayApp) setup() {
    // Create system tray
    t.systray = t.app.SystemTray.New()
    t.systray.SetIcon(icon)
    t.systray.SetLabel("Inactive")
    
    // Create menu
    t.createMenu()
    
    // Create window (hidden by default)
    t.window = t.app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:  "Tray Application",
        Width:  400,
        Height: 600,
        Hidden: true,
    })
    
    // Attach window to tray
    t.systray.AttachWindow(t.window)
    t.systray.WindowOffset(10)
    
    // Handle tray clicks
    t.systray.OnRightClick(func() {
        t.systray.OpenMenu()
    })
    
    // Start background task
    go t.backgroundTask()
}

func (t *TrayApp) createMenu() {
    t.menu = t.app.NewMenu()
    
    // Status item (disabled)
    statusItem := t.menu.Add("Status: Inactive")
    statusItem.SetEnabled(false)
    
    t.menu.AddSeparator()
    
    // Toggle active
    t.menu.Add("Start").OnClick(func(ctx *application.Context) {
        t.toggleActive()
    })
    
    // Show window
    t.menu.Add("Show Window").OnClick(func(ctx *application.Context) {
        t.window.Show()
        t.window.Focus()
    })
    
    t.menu.AddSeparator()
    
    // Settings
    t.menu.AddCheckbox("Start at Login", false).OnClick(func(ctx *application.Context) {
        enabled := ctx.ClickedMenuItem().Checked()
        t.setStartAtLogin(enabled)
    })
    
    t.menu.AddSeparator()
    
    // Quit
    t.menu.Add("Quit").OnClick(func(ctx *application.Context) {
        t.app.Quit()
    })
    
    t.systray.SetMenu(t.menu)
}

func (t *TrayApp) toggleActive() {
    t.isActive = !t.isActive
    t.updateTray()
}

func (t *TrayApp) updateTray() {
    if t.isActive {
        t.systray.SetIcon(iconActive)
        t.systray.SetLabel("Active")
    } else {
        t.systray.SetIcon(icon)
        t.systray.SetLabel("Inactive")
    }
    
    // Rebuild menu with new status
    t.createMenu()
}

func (t *TrayApp) backgroundTask() {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        if t.isActive {
            fmt.Println("Background task running...")
            // Do work
        }
    }
}

func (t *TrayApp) setStartAtLogin(enabled bool) {
    // Implementation varies by platform
    fmt.Printf("Start at login: %v\n", enabled)
}
```

## Kontrol Visibilitas

Tampilkan atau sembunyikan ikon tray secara dinamis:

```go
// Hide tray icon
systray.Hide()

// Show tray icon
systray.Show()
```

Tidak ada getter `IsVisible()` — jika diperlukan, lacak visibilitas dalam status aplikasi Anda sendiri.

**Dukungan Platform:**

| Platform | Hide() | Show() | Catatan |
| --- | --- | --- | --- |
| **Windows** | ✅ | ✅ | Berfungsi sepenuhnya - ikon muncul atau menghilang dari area notifikasi |
| **macOS** | ✅ | ✅ | Item bilah menu ditampilkan atau disembunyikan |
| **Linux** | ✅ | ✅ | Bervariasi menurut lingkungan desktop |

**Kasus penggunaan:**

- Sembunyikan ikon tray untuk sementara berdasarkan preferensi pengguna
- Mode headless dengan ikon tray yang hanya muncul saat diperlukan
- Alihkan visibilitas berdasarkan status aplikasi

**Contoh - Visibilitas Tray Kondisional:**

```go
func (t *TrayApp) setTrayVisibility(visible bool) {
    if visible {
        t.systray.Show()
    } else {
        t.systray.Hide()
    }
}

// Show tray only when updates are available
func (t *TrayApp) checkForUpdates() {
    if hasUpdates {
        t.systray.Show()
        t.systray.SetLabel("Update Available")
    } else {
        t.systray.Hide()
    }
}
```

## Pembersihan

Hancurkan ikon tray setelah selesai:

```go
// In OnShutdown
app := application.New(application.Options{
    OnShutdown: func() {
        if systray != nil {
            systray.Destroy()
        }
    },
})
```

**Penting:** Selalu hancurkan tray sistem saat aplikasi dimatikan untuk melepaskan sumber daya.

## Praktik Terbaik

### ✅ Lakukan

- **Gunakan ikon templat di macOS** - Menyesuaikan dengan mode gelap
- **Buat label tetap singkat** - Maksimum 3-5 karakter
- **Sediakan tooltip di Windows** - Membantu pengguna mengenali aplikasi Anda
- **Uji di semua platform** - Perilakunya berbeda-beda
- **Tangani klik dengan tepat** - Klik kiri untuk tindakan utama, klik kanan untuk menu
- **Perbarui ikon sesuai status** - Umpan balik visual itu penting
- **Hancurkan saat aplikasi ditutup** - Bebaskan sumber daya

### ❌ Jangan

- **Jangan gunakan ikon berukuran besar** - Ikuti pedoman platform
- **Jangan gunakan label panjang** - Label akan terpotong
- **Jangan lupakan mode gelap** - Uji dalam mode gelap Windows dan macOS
- **Jangan memblokir handler klik** - Pastikan handler tetap cepat
- **Jangan lupa memanggil menu.Update()** - Setelah mengubah status menu
- **Jangan berasumsi bahwa baki sistem didukung** - Beberapa lingkungan desktop Linux tidak mendukungnya

## Pemecahan Masalah

### Ikon Baki Tidak Muncul

**Kemungkinan penyebab:**

1. Format ikon tidak didukung
2. Ukuran ikon terlalu besar/kecil
3. Baki sistem tidak didukung (Linux)

**Solusi:**

Tidak ada helper `SystemTraySupported()`; sebagai gantinya, buat baki, periksa platform, dan turunkan fungsionalitas secara wajar:

```go
// Probe support: on Linux without a notification-area extension, the tray
// will simply not appear. Defensive code can fall back to window-only mode
// based on runtime.GOOS or after a short timeout if no tray events arrive.
systray := app.SystemTray.New()
systray.SetIcon(iconBytes)
```

### Ikon Terlihat Tidak Tepat di macOS

**Penyebab:** Tidak menggunakan ikon templat

**Solusi:**

```go
// Use template icon
systray.SetTemplateIcon(iconBytes)

// Or design icon as template (black + transparent)
```

### Menu Tidak Diperbarui

**Penyebab:** Lupa memanggil `menu.Update()`

**Solusi:**

```go
menuItem.SetLabel("New Label")
menu.Update()  // Add this!
```

## Langkah Berikutnya

@cards{cols="2"}
📖 Referensi Menu
Referensi lengkap untuk jenis dan properti item menu.

[Pelajari Lebih Lanjut →](/features/menus/reference/)

---
☰ Menu Aplikasi
Buat bilah menu aplikasi.

[Pelajari Lebih Lanjut →](/features/menus/application/)

---
◆ Menu Konteks
Buat menu konteks klik kanan.

[Pelajari Lebih Lanjut →](/features/menus/context/)

---
📖 Contoh Baki Sistem
Jelajahi aplikasi baki sistem yang lengkap.

[Pelajari Lebih Lanjut →](https://github.com/wailsapp/wails/tree/master/v3/examples/systray-basic)

@end

---

**Ada pertanyaan?** Tanyakan di [Discord](https://discord.gg/JDdSxwjhGf) atau lihat [contoh baki sistem](https://github.com/wailsapp/wails/tree/master/v3/examples/systray-basic).
