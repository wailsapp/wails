---
title: "Menu Aplikasi"
description: "Buat bilah menu native untuk aplikasi desktop Anda"
slug: "features/menus/application"
sourcePath: "features/menus/application.md"
---

## Masalah

Aplikasi desktop profesional memerlukan bilah menu—File, Edit, Tampilan, Bantuan. Namun, menu bekerja secara berbeda di setiap platform:

- **macOS**: Bilah menu global di bagian atas layar
- **Windows**: Bilah menu pada bilah judul jendela
- **Linux**: Berbeda-beda menurut lingkungan desktop

Membuat menu yang sesuai untuk setiap platform secara manual itu merepotkan dan rawan kesalahan.

## Solusi Wails

Wails menyediakan **API terpadu** yang secara otomatis membuat menu native untuk setiap platform. Tulis sekali dan dapatkan perilaku native di semua platform.

![Menu aplikasi Wails di macOS dengan item standar, kotak centang, tombol radio, dan submenu](/assets/screenshots/application-menu-macos.png)

Di macOS, menu aplikasi ditempatkan pada bilah menu global. Tangkapan layar ini menampilkan menu native yang dirender oleh API menu Wails, termasuk item yang dinonaktifkan, kotak centang, tombol radio, dan submenu.

## Mulai Cepat

```go
package main

import (
    "runtime"
    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "My App",
    })

    // Create menu
    menu := app.NewMenu()

    // Add standard menus (platform-appropriate)
    if runtime.GOOS == "darwin" {
        menu.AddRole(application.AppMenu)  // macOS only
    }
    menu.AddRole(application.FileMenu)
    menu.AddRole(application.EditMenu)
    menu.AddRole(application.WindowMenu)
    menu.AddRole(application.HelpMenu)

    // Set the application menu
    app.Menu.Set(menu)

    // Create window with UseApplicationMenu to inherit the menu on Windows/Linux
    app.Window.NewWithOptions(application.WebviewWindowOptions{
        UseApplicationMenu: true,
    })

    app.Run()
}
```

**Selesai!** Kini Anda memiliki menu native untuk setiap platform dengan item standar. Opsi `UseApplicationMenu` memastikan jendela Windows dan Linux menampilkan menu tanpa kode tambahan.

## Membuat Menu

### Pembuatan Menu Dasar

```go
// Create a new menu
menu := app.NewMenu()

// Add a top-level menu
fileMenu := menu.AddSubmenu("File")

// Add menu items
fileMenu.Add("New").OnClick(func(ctx *application.Context) {
    // Handle New
})

fileMenu.Add("Open").OnClick(func(ctx *application.Context) {
    // Handle Open
})

fileMenu.AddSeparator()

fileMenu.Add("Quit").OnClick(func(ctx *application.Context) {
    app.Quit()
})
```

### Mengatur Menu

**Pendekatan yang disarankan** — Gunakan `UseApplicationMenu` untuk menjaga konsistensi lintas platform:

```go
// Set the application menu once
app.Menu.Set(menu)

// Create windows that inherit the menu on Windows/Linux
app.Window.NewWithOptions(application.WebviewWindowOptions{
    UseApplicationMenu: true,  // Window uses the app menu
})
```

Pendekatan ini menghasilkan perilaku berikut:

- Di **macOS**: Menu muncul di bagian atas layar (perilaku standar)
- Di **Windows/Linux**: Setiap jendela dengan `UseApplicationMenu: true` menampilkan menu aplikasi

**Detail khusus platform:**

@tabs{sync-key="platform"}
[macOS]
**Bilah menu global** (satu per aplikasi):

```go
app.Menu.Set(menu)
```

Menu muncul di bagian atas layar dan tetap ada meskipun semua jendela ditutup. Opsi `UseApplicationMenu` tidak berpengaruh di macOS karena semua aplikasi menggunakan menu global.

[Windows]
**Bilah menu per jendela**:

```go
// Option 1: Use application menu (recommended)
app.Menu.Set(menu)
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    UseApplicationMenu: true,
})

// Option 2: Set menu directly on window
window.SetMenu(menu)
```

Setiap jendela dapat memiliki menu sendiri atau mewarisi menu aplikasi. Menu muncul pada bilah judul jendela.

[Linux]
**Bilah menu per jendela** (biasanya):

```go
// Option 1: Use application menu (recommended)
app.Menu.Set(menu)
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    UseApplicationMenu: true,
})

// Option 2: Set menu directly on window
window.SetMenu(menu)
```

Perilakunya berbeda-beda menurut lingkungan desktop. Beberapa lingkungan (seperti Unity) mendukung menu global.

@end

@note{type="tip" title="Sederhanakan Menu Lintas Platform"}
Penggunaan `UseApplicationMenu: true` meniadakan kebutuhan akan kode khusus platform seperti:

```go
// Old approach - no longer needed
if runtime.GOOS == "darwin" {
    app.Menu.Set(menu)
} else {
    window.SetMenu(menu)
}
```

@end

**Menu khusus per jendela:**

Jika suatu jendela memerlukan menu yang berbeda dari menu aplikasi, atur menu tersebut secara langsung:

```go
window.SetMenu(customMenu)  // Overrides UseApplicationMenu
```

## Peran Menu

Wails menyediakan **peran menu bawaan** yang secara otomatis membuat struktur menu yang sesuai untuk setiap platform.

### Peran yang Tersedia

| Peran | Deskripsi | Catatan Platform |
| --- | --- | --- |
| `AppMenu` | Menu aplikasi dengan Tentang, Preferensi, Keluar | **Khusus macOS** |
| `FileMenu` | Operasi file (Baru, Buka, Simpan, dan sebagainya) | Semua platform |
| `EditMenu` | Penyuntingan teks (Urungkan, Ulangi, Potong, Salin, Tempel) | Semua platform |
| `WindowMenu` | Pengelolaan jendela (Minimalkan, Perbesar, dan sebagainya) | Semua platform |
| `HelpMenu` | Bantuan dan informasi | Semua platform |

### Menggunakan Peran

```go
menu := app.NewMenu()

// macOS: Add application menu
if runtime.GOOS == "darwin" {
    menu.AddRole(application.AppMenu)
}

// All platforms: Add standard menus
menu.AddRole(application.FileMenu)
menu.AddRole(application.EditMenu)
menu.AddRole(application.WindowMenu)
menu.AddRole(application.HelpMenu)
```

**Yang Anda dapatkan:**

@tabs{sync-key="platform"}
[macOS]
**AppMenu** (dengan nama aplikasi):

- Tentang [Nama Aplikasi]
- Preferensi... (⌘,)
- ---
- Layanan
- ---
- Sembunyikan [Nama Aplikasi] (⌘H)
- Sembunyikan Lainnya (⌥⌘H)
- Tampilkan Semua
- ---
- Keluar dari [Nama Aplikasi] (⌘Q)

**FileMenu**:

- Baru (⌘N)
- Buka... (⌘O)
- ---
- Tutup Jendela (⌘W)

**EditMenu**:

- Urungkan (⌘Z)
- Ulangi (⇧⌘Z)
- ---
- Potong (⌘X)
- Salin (⌘C)
- Tempel (⌘V)
- Pilih Semua (⌘A)

**WindowMenu**:

- Minimalkan (⌘M)
- Perbesar
- ---
- Bawa Semua ke Depan

**HelpMenu**:

- Bantuan [Nama Aplikasi]

[Windows]
**FileMenu**:

- Baru (Ctrl+N)
- Buka... (Ctrl+O)
- ---
- Keluar (Alt+F4)

**EditMenu**:

- Urungkan (Ctrl+Z)
- Ulangi (Ctrl+Y)
- ---
- Potong (Ctrl+X)
- Salin (Ctrl+C)
- Tempel (Ctrl+V)
- Pilih Semua (Ctrl+A)

**WindowMenu**:

- Minimalkan
- Maksimalkan

**HelpMenu**:

- Tentang [Nama Aplikasi]

[Linux]
Serupa dengan Windows, tetapi pintasan keyboard dapat berbeda menurut lingkungan desktop.

@end

### Menyesuaikan Menu Peran

`Menu.AddRole(role)` mengembalikan menu **penerima** (menu tingkat teratas), **bukan** submenu peran. Untuk menambahkan item ke submenu peran, cari item peran yang disisipkan dengan `FindByRole`, lalu panggil `GetSubmenu()` pada item tersebut:

```go
menu.AddRole(application.FileMenu)

fileMenu := menu.FindByRole(application.FileMenu).GetSubmenu()
fileMenu.Add("Import...").OnClick(handleImport)
fileMenu.Add("Export...").OnClick(handleExport)
```

## Menu Khusus

Buat menu sendiri untuk fitur khusus aplikasi:

```go
// Add a custom top-level menu
toolsMenu := menu.AddSubmenu("Tools")

// Add items
toolsMenu.Add("Settings").OnClick(func(ctx *application.Context) {
    showSettingsWindow()
})

toolsMenu.AddSeparator()

// Add checkbox
toolsMenu.AddCheckbox("Dark Mode", false).OnClick(func(ctx *application.Context) {
    isDark := ctx.ClickedMenuItem().Checked()
    setTheme(isDark)
})

// Add radio group
toolsMenu.AddRadio("Small", true).OnClick(handleFontSize)
toolsMenu.AddRadio("Medium", false).OnClick(handleFontSize)
toolsMenu.AddRadio("Large", false).OnClick(handleFontSize)

// Add submenu
advancedMenu := toolsMenu.AddSubmenu("Advanced")
advancedMenu.Add("Configure...").OnClick(showAdvancedSettings)
```

**Untuk jenis item menu lainnya**, lihat [Referensi Menu](/features/menus/reference/).

## Menu Dinamis

Perbarui menu berdasarkan status aplikasi:

### Mengaktifkan/Menonaktifkan Item

```go
var saveMenuItem *application.MenuItem

func createMenu() {
    menu := app.NewMenu()
    fileMenu := menu.AddSubmenu("File")
    
    saveMenuItem = fileMenu.Add("Save")
    saveMenuItem.SetEnabled(false)  // Initially disabled
    saveMenuItem.OnClick(handleSave)
    
    app.Menu.Set(menu)
}

func onDocumentChanged() {
    saveMenuItem.SetEnabled(hasUnsavedChanges())
    menu.Update()  // Important!
}
```

@note{type="caution" title="Selalu Panggil menu.Update()"}
Setelah mengubah status menu (aktif/nonaktif, label, status centang), **selalu panggil `menu.Update()`**. Hal ini sangat penting terutama di Windows, tempat menu direkonstruksi.

Lihat [Referensi Menu](/features/menus/reference/#enabled-state) untuk detailnya.

@end

### Mengubah Label

```go
updateMenuItem := menu.Add("Check for Updates")

updateMenuItem.OnClick(func(ctx *application.Context) {
    updateMenuItem.SetLabel("Checking...")
    menu.Update()
    
    checkForUpdates()
    
    updateMenuItem.SetLabel("Check for Updates")
    menu.Update()
})
```

### Membangun Ulang Menu

Untuk perubahan besar, bangun ulang seluruh menu:

```go
func rebuildFileMenu() {
    menu := app.NewMenu()
    fileMenu := menu.AddSubmenu("File")
    
    fileMenu.Add("New").OnClick(handleNew)
    fileMenu.Add("Open").OnClick(handleOpen)
    
    // Add recent files dynamically
    if hasRecentFiles() {
        recentMenu := fileMenu.AddSubmenu("Open Recent")
        for _, file := range getRecentFiles() {
            filePath := file  // Capture for closure
            recentMenu.Add(filepath.Base(file)).OnClick(func(ctx *application.Context) {
                openFile(filePath)
            })
        }
        recentMenu.AddSeparator()
        recentMenu.Add("Clear Recent").OnClick(clearRecentFiles)
    }
    
    fileMenu.AddSeparator()
    fileMenu.Add("Quit").OnClick(func(ctx *application.Context) {
        app.Quit()
    })
    
    app.Menu.Set(menu)
}
```

## Mengontrol Jendela dari Menu

Item menu dapat mengontrol jendela:

```go
viewMenu := menu.AddSubmenu("View")

// Toggle fullscreen
viewMenu.Add("Toggle Fullscreen").OnClick(func(ctx *application.Context) {
    if window, ok := app.Window.GetByName("main"); ok {
        window.ToggleFullscreen()
    }
})

// Zoom controls
viewMenu.Add("Zoom In").SetAccelerator("CmdOrCtrl++").OnClick(func(ctx *application.Context) {
    // Increase zoom
})

viewMenu.Add("Zoom Out").SetAccelerator("CmdOrCtrl+-").OnClick(func(ctx *application.Context) {
    // Decrease zoom
})

viewMenu.Add("Reset Zoom").SetAccelerator("CmdOrCtrl+0").OnClick(func(ctx *application.Context) {
    // Reset zoom
})
```

**Dapatkan jendela aktif:**

```go
menuItem.OnClick(func(ctx *application.Context) {
    window := application.Get().Window.Current() // the window the menu was invoked from
    // Use window
})
```

## Pertimbangan Khusus Platform

### macOS

**Perilaku bilah menu:**

- Muncul di **bagian atas layar** (global)
- Tetap ada ketika semua jendela ditutup
- Menu pertama **selalu merupakan menu aplikasi**
- Gunakan `menu.AddRole(application.AppMenu)` untuk item standar

**Lokasi standar:**

- **Tentang**: Menu aplikasi
- **Preferensi**: Menu aplikasi (⌘,)
- **Keluar**: Menu aplikasi (⌘Q)
- **Bantuan**: Menu Bantuan

**Contoh:**

```go
if runtime.GOOS == "darwin" {
    menu.AddRole(application.AppMenu)  // Adds About, Preferences, Quit
    
    // Don't add Quit to File menu on macOS
    // Don't add About to Help menu on macOS
}
```

### Windows

**Perilaku bilah menu:**

- Muncul di **bilah judul jendela**
- Setiap jendela memiliki menu sendiri
- Tidak ada menu aplikasi

**Lokasi standar:**

- **Keluar**: Menu File (Alt+F4)
- **Pengaturan**: Menu Alat atau Edit
- **Tentang**: Menu Bantuan

**Contoh:**

```go
if runtime.GOOS == "windows" {
    menu.AddRole(application.FileMenu) // Exit is added automatically
    menu.AddRole(application.HelpMenu) // About is added automatically
}
```

### Linux

**Perilaku bilah menu:**

- Biasanya tersedia per jendela (seperti Windows)
- Beberapa DE mendukung menu global (Unity, GNOME dengan ekstensi)
- Tampilan bervariasi menurut lingkungan desktop

**Praktik terbaik:** Ikuti konvensi Windows dan uji pada DE target.

## Contoh Lengkap

Berikut struktur menu yang siap digunakan dalam produksi:

```go
package main

import (
    "runtime"
    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "My Application",
    })

    // Create and set menu
    createMenu(app)

    // Create main window with UseApplicationMenu for cross-platform menu support
    app.Window.NewWithOptions(application.WebviewWindowOptions{
        UseApplicationMenu: true,
    })

    app.Run()
}

func createMenu(app *application.App) {
    menu := app.NewMenu()

    // Platform-specific application menu (macOS only)
    if runtime.GOOS == "darwin" {
        menu.AddRole(application.AppMenu)
    }

    // File menu — AddRole returns the receiver menu, not the role submenu.
    // To add items into the File submenu, look it up via FindByRole + GetSubmenu.
    menu.AddRole(application.FileMenu)
    fileMenu := menu.FindByRole(application.FileMenu).GetSubmenu()
    fileMenu.Add("Import...").SetAccelerator("CmdOrCtrl+I").OnClick(handleImport)
    fileMenu.Add("Export...").SetAccelerator("CmdOrCtrl+E").OnClick(handleExport)

    // Edit menu
    menu.AddRole(application.EditMenu)

    // View menu
    viewMenu := menu.AddSubmenu("View")
    viewMenu.Add("Toggle Fullscreen").SetAccelerator("F11").OnClick(toggleFullscreen)
    viewMenu.AddSeparator()
    viewMenu.AddCheckbox("Show Sidebar", true).OnClick(toggleSidebar)
    viewMenu.AddCheckbox("Show Toolbar", true).OnClick(toggleToolbar)

    // Tools menu
    toolsMenu := menu.AddSubmenu("Tools")
    
    // Settings location varies by platform
    if runtime.GOOS == "darwin" {
        // On macOS, Preferences is in Application menu (added by AppMenu role)
    } else {
        toolsMenu.Add("Settings").SetAccelerator("CmdOrCtrl+,").OnClick(showSettings)
    }
    
    toolsMenu.AddSeparator()
    toolsMenu.AddCheckbox("Dark Mode", false).OnClick(toggleDarkMode)

    // Window menu
    menu.AddRole(application.WindowMenu)

    // Help menu
    helpMenu := menu.AddRole(application.HelpMenu)
    helpMenu.Add("Documentation").OnClick(openDocumentation)
    
    // About location varies by platform
    if runtime.GOOS == "darwin" {
        // On macOS, About is in Application menu (added by AppMenu role)
    } else {
        helpMenu.AddSeparator()
        helpMenu.Add("About").OnClick(showAbout)
    }

    // Set the application menu
    app.Menu.Set(menu)
}

func handleImport(ctx *application.Context) {
    // Implementation
}

func handleExport(ctx *application.Context) {
    // Implementation
}

func toggleFullscreen(ctx *application.Context) {
    window := application.Get().Window.Current()
    window.ToggleFullscreen()
}

func toggleSidebar(ctx *application.Context) {
    // Implementation
}

func toggleToolbar(ctx *application.Context) {
    // Implementation
}

func showSettings(ctx *application.Context) {
    // Implementation
}

func toggleDarkMode(ctx *application.Context) {
    isDark := ctx.ClickedMenuItem().Checked()
    // Apply theme
}

func openDocumentation(ctx *application.Context) {
    // Open browser
}

func showAbout(ctx *application.Context) {
    // Show about dialog
}
```

## Praktik Terbaik

### ✅ Lakukan

- **Gunakan peran menu** untuk menu standar (File, Edit, dan sebagainya)
- **Ikuti konvensi platform** untuk struktur menu
- **Tambahkan pintasan papan ketik** untuk tindakan yang umum digunakan
- **Panggil menu.Update()** setelah mengubah status menu
- **Uji pada semua platform** — perilakunya bervariasi
- **Jaga agar hierarki menu tetap dangkal** — maksimum 2-3 tingkat
- **Gunakan label yang jelas** — "Simpan Proyek", bukan "Simpan"

### ❌ Jangan Lakukan

- **Jangan melakukan hardcode pada pintasan khusus platform** — Gunakan `CmdOrCtrl`
- **Jangan menempatkan Keluar di menu File pada macOS** — Opsi tersebut berada di menu aplikasi
- **Jangan menempatkan Tentang di menu Bantuan pada macOS** — Opsi tersebut berada di menu aplikasi
- **Jangan lupa memanggil menu.Update()** — Menu tidak akan berfungsi dengan semestinya
- **Jangan membuat menu bertingkat terlalu dalam** — Pengguna akan kebingungan
- **Jangan gunakan jargon** — Gunakan label yang mudah dipahami pengguna

## Langkah Berikutnya

@cards{cols="2"}
📖 Referensi Menu
Referensi lengkap untuk jenis dan properti item menu.

[Pelajari Selengkapnya →](/features/menus/reference/)

---
◆ Menu Konteks
Buat menu konteks yang dibuka dengan klik kanan.

[Pelajari Selengkapnya →](/features/menus/context/)

---
★ Menu Baki Sistem
Tambahkan integrasi dengan baki sistem/bilah menu.

[Pelajari Selengkapnya →](/features/menus/systray/)

---
📖 Pola Menu
Pola menu umum dan praktik terbaik.

[Pelajari Selengkapnya →](/guides/menus/)

@end

---

**Ada pertanyaan?** Tanyakan di [Discord](https://discord.gg/JDdSxwjhGf) atau lihat [contoh menu](https://github.com/wailsapp/wails/tree/master/v3/examples/menu).
