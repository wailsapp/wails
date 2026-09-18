---
title: "API Menu"
description: "Referensi lengkap untuk API Menu"
slug: "reference/menu"
sourcePath: "reference/menu.md"
---

## Ikhtisar

API Menu menyediakan metode untuk membuat dan mengelola menu aplikasi, menu konteks, dan menu baki sistem.

**Jenis Menu:**

- **Menu Aplikasi** - Bilah menu atas (File, Edit, dan sebagainya)
- **Menu Konteks** - Menu klik kanan
- **Menu Baki Sistem** - Menu di baki sistem/area notifikasi

## Membuat Menu

### NewMenu()

Membuat menu baru.

```go
func (a *App) NewMenu() *Menu
```

**Contoh:**

```go
menu := app.NewMenu()
```

## Metode Menu

### Add()

Menambahkan item menu ke menu.

```go
func (m *Menu) Add(label string) *MenuItem
```

**Parameter:**

- `label` - Teks yang ditampilkan untuk item menu

**Mengembalikan:** Item menu yang dibuat

**Contoh:**

```go
item := menu.Add("Open File")
item.OnClick(func(ctx *application.Context) {
    // Handle click
})
```

### AddSubmenu()

Menambahkan submenu ke menu.

```go
func (m *Menu) AddSubmenu(label string) *Menu
```

**Parameter:**

- `label` - Label submenu

**Mengembalikan:** Submenu yang dibuat

**Contoh:**

```go
fileMenu := menu.AddSubmenu("File")
fileMenu.Add("New")
fileMenu.Add("Open")
fileMenu.Add("Save")
```

### AddSeparator()

Menambahkan garis pemisah visual di antara item menu.

```go
func (m *Menu) AddSeparator()
```

**Contoh:**

```go
menu.Add("Copy")
menu.Add("Paste")
menu.AddSeparator()
menu.Add("Select All")
```

**Praktik terbaik:** Gunakan pemisah untuk mengelompokkan item menu yang saling berkaitan.

### AddCheckbox()

Menambahkan item menu yang dapat dicentang.

```go
func (m *Menu) AddCheckbox(label string, checked bool) *MenuItem
```

**Parameter:**

- `label` - Label kotak centang
- `checked` - Status centang awal

**Contoh:**

```go
darkMode := menu.AddCheckbox("Dark Mode", false)
darkMode.OnClick(func(ctx *application.Context) {
    isChecked := darkMode.Checked()
    // Toggle dark mode
})
```

### AddRadio()

Menambahkan item menu tombol radio (grup yang saling eksklusif).

```go
func (m *Menu) AddRadio(label string, checked bool) *MenuItem
```

**Parameter:**

- `label` - Label tombol radio
- `checked` - Status centang awal

**Contoh:**

```go
// Create radio group for view modes
viewMenu := menu.AddSubmenu("View")
listView := viewMenu.AddRadio("List View", true)
gridView := viewMenu.AddRadio("Grid View", false)
treeView := viewMenu.AddRadio("Tree View", false)

listView.OnClick(func(ctx *application.Context) {
    setViewMode("list")
})
gridView.OnClick(func(ctx *application.Context) {
    setViewMode("grid")
})
```

### Update()

Memperbarui menu agar mencerminkan setiap perubahan yang dibuat pada item menu.

```go
func (m *Menu) Update()
```

**Contoh:**

```go
item.SetEnabled(false)
menu.Update()  // Must call to apply changes
```

**Penting:** Selalu panggil `Update()` setelah mengubah properti item menu.

## Metode Item Menu

### OnClick()

Mendaftarkan pengendali klik untuk item menu.

```go
func (mi *MenuItem) OnClick(callback func(ctx *application.Context)) *MenuItem
```

**Parameter:**

- `callback` - Fungsi yang dipanggil saat item diklik

**Mengembalikan:** Item menu (untuk perantaian)

**Contoh:**

```go
item.OnClick(func(ctx *application.Context) {
    fmt.Println("Menu item clicked")
    app.Logger.Info("User clicked menu item")
})
```

### SetLabel()

Mengubah label item menu.

```go
func (mi *MenuItem) SetLabel(label string) *MenuItem
```

**Contoh:**

```go
item.SetLabel("Save As...")
menu.Update()
```

### SetEnabled()

Mengaktifkan atau menonaktifkan item menu.

```go
func (mi *MenuItem) SetEnabled(enabled bool) *MenuItem
```

**Contoh:**

```go
// Disable save when no document is open
saveItem.SetEnabled(hasOpenDocument)
menu.Update()
```

**Pola umum:**

```go
// Update menu state based on application state
func updateMenuState() {
    saveItem.SetEnabled(hasUnsavedChanges)
    undoItem.SetEnabled(canUndo)
    redoItem.SetEnabled(canRedo)
    menu.Update()
}
```

### SetChecked()

Mengatur status centang untuk item menu kotak centang/tombol radio.

```go
func (mi *MenuItem) SetChecked(checked bool) *MenuItem
```

**Contoh:**

```go
darkModeItem.SetChecked(isDarkModeEnabled)
menu.Update()
```

### Checked()

Mengembalikan status centang saat ini.

```go
func (mi *MenuItem) Checked() bool
```

**Contoh:**

```go
if darkModeItem.Checked() {
    // Dark mode is enabled
}
```

### SetAccelerator()

Mengatur pintasan papan ketik untuk item menu.

```go
func (mi *MenuItem) SetAccelerator(accelerator string) *MenuItem
```

**Parameter:**

- `accelerator` - Pintasan papan ketik (misalnya, "Ctrl+S", "Cmd+Q")

**Format akselerator:**

- **Tombol pengubah:** `Ctrl`, `Cmd`, `Alt`, `Shift`
- **Tombol:** `A-Z`, `0-9`, `F1-F12`, `Enter`, `Backspace`, dan sebagainya.
- **Platform:** Gunakan `Cmd` di macOS, `Ctrl` di Windows/Linux

**Contoh:**

```go
saveItem.SetAccelerator("Ctrl+S")
quitItem.SetAccelerator("Ctrl+Q")
newItem.SetAccelerator("Ctrl+N")
```

**Contoh yang menyesuaikan platform:**

```go
import "runtime"

var quitShortcut string
if runtime.GOOS == "darwin" {
    quitShortcut = "Cmd+Q"
} else {
    quitShortcut = "Ctrl+Q"
}
quitItem.SetAccelerator(quitShortcut)
```

### SetTooltip()

Mengatur tooltip yang muncul saat penunjuk diarahkan ke item menu.

```go
func (mi *MenuItem) SetTooltip(tooltip string) *MenuItem
```

**Contoh:**

```go
item.SetTooltip("Opens a file from disk")
```

### SetHidden()

Menampilkan atau menyembunyikan item menu.

```go
func (mi *MenuItem) SetHidden(hidden bool) *MenuItem
```

**Contoh:**

```go
// Hide debug menu in production
debugItem.SetHidden(!isDevelopment)
menu.Update()
```

## Menu Aplikasi

### app.Menu.Set()

Mengatur bilah menu utama aplikasi.

```go
func (mm *MenuManager) Set(menu *Menu)
```

**Contoh:**

```go
menu := app.NewMenu()

// File menu
fileMenu := menu.AddSubmenu("File")
fileMenu.Add("New").SetAccelerator("Ctrl+N").OnClick(newFile)
fileMenu.Add("Open").SetAccelerator("Ctrl+O").OnClick(openFile)
fileMenu.Add("Save").SetAccelerator("Ctrl+S").OnClick(saveFile)
fileMenu.AddSeparator()
fileMenu.Add("Exit").SetAccelerator("Ctrl+Q").OnClick(func(ctx *application.Context) {
    app.Quit()
})

// Edit menu
editMenu := menu.AddSubmenu("Edit")
editMenu.Add("Undo").SetAccelerator("Ctrl+Z").OnClick(undo)
editMenu.Add("Redo").SetAccelerator("Ctrl+Y").OnClick(redo)
editMenu.AddSeparator()
editMenu.Add("Cut").SetAccelerator("Ctrl+X").OnClick(cut)
editMenu.Add("Copy").SetAccelerator("Ctrl+C").OnClick(copy)
editMenu.Add("Paste").SetAccelerator("Ctrl+V").OnClick(paste)

app.Menu.Set(menu)
```

**Catatan platform:**

- **macOS:** Menu muncul di bilah menu bagian atas
- **Windows/Linux:** Menu muncul di bilah judul jendela
- **macOS:** Secara otomatis menambahkan menu aplikasi dengan nama aplikasi

## Menu Konteks

### app.ContextMenu.New() / app.ContextMenu.Add()

Buat `*ContextMenu` melalui pengelola, lalu daftarkan dengan suatu nama. `ContextMenuManager.Add` menerima `*ContextMenu` — **bukan** `*Menu`, dan **tidak ada** metode `app.RegisterContextMenu`.

```go
func (cm *ContextMenuManager) New() *ContextMenu
func (cm *ContextMenuManager) Add(name string, menu *ContextMenu)
func (cm *ContextMenuManager) Get(name string) (*ContextMenu, bool)
func (cm *ContextMenuManager) Remove(name string)
```

Sebagai alternatif, `application.NewContextMenu(name string) *ContextMenu` pada tingkat paket membuat dan mendaftarkan menu konteks dalam satu langkah.

**Go:**

```go
// Build the context menu via the manager
contextMenu := app.ContextMenu.New()
contextMenu.Add("Cut").OnClick(cut)
contextMenu.Add("Copy").OnClick(copy)
contextMenu.Add("Paste").OnClick(paste)
contextMenu.AddSeparator()
contextMenu.Add("Select All").OnClick(selectAll)

// Register under a name; HTML opts into it via the CSS custom property below.
app.ContextMenu.Add("editor", contextMenu)
```

**HTML/CSS:**

Runtime memicu menu konteks terdaftar ketika target klik kanan (atau salah satu leluhurnya) memiliki properti kustom CSS `--custom-contextmenu` yang diatur ke nama tersebut. Properti opsional `--custom-contextmenu-data` diteruskan ke callback Go melalui `ctx.ContextMenuData()`. Untuk menonaktifkan menu konteks bawaan browser, atur `--default-contextmenu: hide` (atau `auto`/`show`).

```html
<!-- Trigger context menu on right-click -->
<div style="--custom-contextmenu: editor; --default-contextmenu: hide">
    Right-click here for context menu
</div>
```

Tidak ada atribut `data-wails-context-menu="..."` — atribut tersebut tidak pernah dihubungkan dalam runtime.

**Menu konteks dinamis:**

```go
// Update context menu based on selection
func updateContextMenu() {
    contextMenu := app.ContextMenu.New()

    if hasSelection {
        contextMenu.Add("Cut").OnClick(cut)
        contextMenu.Add("Copy").OnClick(copy)
    }

    contextMenu.Add("Paste").SetEnabled(hasClipboardContent).OnClick(paste)

    app.ContextMenu.Add("editor", contextMenu)
}
```

## Menu Baki Sistem

### app.SystemTray.New()

Membuat ikon baki sistem baru.

```go
func (sm *SystemTrayManager) New() *SystemTray
```

**Contoh:**

```go
tray := app.SystemTray.New()
```

### SetIcon()

Mengatur ikon baki sistem.

```go
func (st *SystemTray) SetIcon(icon []byte) *SystemTray
```

**Contoh:**

```go
iconData, _ := os.ReadFile("icon.png")
tray.SetIcon(iconData)
```

### SetMenu()

Mengatur menu untuk baki sistem.

```go
func (st *SystemTray) SetMenu(menu *Menu) *SystemTray
```

**Contoh:**

```go
trayMenu := app.NewMenu()
trayMenu.Add("Show Window").OnClick(func(ctx *application.Context) {
    window.Show()
    window.Focus()
})
trayMenu.Add("Settings").OnClick(openSettings)
trayMenu.AddSeparator()
trayMenu.Add("Quit").OnClick(func(ctx *application.Context) {
    app.Quit()
})

tray.SetMenu(trayMenu)
```

### SetTooltip()

Mengatur tooltip yang ditampilkan saat penunjuk diarahkan ke ikon baki sistem. Tidak mengembalikan nilai.

```go
func (st *SystemTray) SetTooltip(tooltip string)
```

**Contoh:**

```go
tray.SetTooltip("My Application - Running")
```

### OnClick()

Menangani klik kiri pada ikon baki sistem.

```go
func (st *SystemTray) OnClick(callback func()) *SystemTray
```

**Contoh:**

```go
tray.OnClick(func() {
    if window.IsVisible() {
        window.Hide()
    } else {
        window.Show()
        window.Focus()
    }
})
```

## Contoh Lengkap

### Menu Aplikasi Standar

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
)

func createMenu(app *application.App) *application.Menu {
    menu := app.NewMenu()

    // File menu
    fileMenu := menu.AddSubmenu("File")
    fileMenu.Add("New").
        SetAccelerator("Ctrl+N").
        OnClick(func(ctx *application.Context) {
            // Create new document
        })
    fileMenu.Add("Open").
        SetAccelerator("Ctrl+O").
        OnClick(func(ctx *application.Context) {
            // Open file dialog
        })
    fileMenu.Add("Save").
        SetAccelerator("Ctrl+S").
        OnClick(func(ctx *application.Context) {
            // Save document
        })
    fileMenu.AddSeparator()
    fileMenu.Add("Exit").OnClick(func(ctx *application.Context) {
        app.Quit()
    })

    // Edit menu
    editMenu := menu.AddSubmenu("Edit")
    editMenu.Add("Undo").SetAccelerator("Ctrl+Z")
    editMenu.Add("Redo").SetAccelerator("Ctrl+Y")
    editMenu.AddSeparator()
    editMenu.Add("Cut").SetAccelerator("Ctrl+X")
    editMenu.Add("Copy").SetAccelerator("Ctrl+C")
    editMenu.Add("Paste").SetAccelerator("Ctrl+V")

    // View menu
    viewMenu := menu.AddSubmenu("View")
    darkMode := viewMenu.AddCheckbox("Dark Mode", false)
    darkMode.OnClick(func(ctx *application.Context) {
        // Toggle dark mode
        isChecked := darkMode.Checked()
        app.Logger.Info("Dark mode", "enabled", isChecked)
    })
    viewMenu.AddSeparator()
    viewMenu.AddRadio("List View", true)
    viewMenu.AddRadio("Grid View", false)
    viewMenu.AddRadio("Detail View", false)

    // Help menu
    helpMenu := menu.AddSubmenu("Help")
    helpMenu.Add("Documentation").OnClick(func(ctx *application.Context) {
        // Open docs
    })
    helpMenu.Add("About").OnClick(func(ctx *application.Context) {
        // Show about dialog
    })

    return menu
}

func main() {
    app := application.New(application.Options{
        Name: "Menu Demo",
    })

    menu := createMenu(app)
    app.Menu.Set(menu)

    window := app.Window.New()
    window.Show()

    app.Run()
}
```

### Aplikasi Baki Sistem

```go
func setupSystemTray(app *application.App, window application.Window) {
    // Create system tray
    tray := app.SystemTray.New()

    // Set icon
    iconData, _ := os.ReadFile("icon.png")
    tray.SetIcon(iconData)
    tray.SetTooltip("My App - Running")

    // Handle left-click on tray icon
    tray.OnClick(func() {
        if window.IsVisible() {
            window.Hide()
        } else {
            window.Show()
            window.Focus()
        }
    })

    // Create tray menu
    trayMenu := app.NewMenu()

    showItem := trayMenu.Add("Show Window")
    showItem.OnClick(func(ctx *application.Context) {
        window.Show()
        window.Focus()
    })

    trayMenu.AddSeparator()

    trayMenu.Add("Settings").OnClick(func(ctx *application.Context) {
        // Open settings window
    })

    trayMenu.Add("About").OnClick(func(ctx *application.Context) {
        // Show about dialog
    })

    trayMenu.AddSeparator()

    trayMenu.Add("Quit").OnClick(func(ctx *application.Context) {
        app.Quit()
    })

    tray.SetMenu(trayMenu)
}
```

### Pembaruan Menu Dinamis

```go
type Editor struct {
    app         *application.App
    menu        *application.Menu
    undoItem    *application.MenuItem
    redoItem    *application.MenuItem
    saveItem    *application.MenuItem
    undoStack   []string
    redoStack   []string
    hasChanges  bool
}

func (e *Editor) createMenu() {
    e.menu = e.app.NewMenu()

    fileMenu := e.menu.AddSubmenu("File")
    e.saveItem = fileMenu.Add("Save").SetAccelerator("Ctrl+S")
    e.saveItem.OnClick(func(ctx *application.Context) {
        e.save()
    })

    editMenu := e.menu.AddSubmenu("Edit")
    e.undoItem = editMenu.Add("Undo").SetAccelerator("Ctrl+Z")
    e.undoItem.OnClick(func(ctx *application.Context) {
        e.undo()
    })

    e.redoItem = editMenu.Add("Redo").SetAccelerator("Ctrl+Y")
    e.redoItem.OnClick(func(ctx *application.Context) {
        e.redo()
    })

    e.updateMenuState()
    e.app.Menu.Set(e.menu)
}

func (e *Editor) updateMenuState() {
    // Update menu items based on current state
    e.saveItem.SetEnabled(e.hasChanges)
    e.undoItem.SetEnabled(len(e.undoStack) > 0)
    e.redoItem.SetEnabled(len(e.redoStack) > 0)
    e.menu.Update()
}

func (e *Editor) onChange() {
    e.hasChanges = true
    e.updateMenuState()
}

func (e *Editor) save() {
    // Save logic
    e.hasChanges = false
    e.updateMenuState()
}
```

## Praktik Terbaik

### ✅ Lakukan

- **Gunakan akselerator standar** — Ikuti konvensi platform (Ctrl+C untuk menyalin, dan sebagainya)
- **Panggil Update() setelah melakukan perubahan** — Jika tidak, menu tidak akan mencerminkan perubahan
- **Kelompokkan item terkait** — Gunakan pemisah untuk menata item menu
- **Nonaktifkan tindakan yang tidak tersedia** — Jangan sembunyikan; nonaktifkan dengan SetEnabled(false)
- **Gunakan label yang jelas** — Buat label tetap ringkas dan deskriptif
- **Ikuti konvensi platform** — Pola menu macOS berbeda dengan Windows/Linux

### ❌ Jangan Lakukan

- **Jangan lupa memanggil Update()** — Ini kesalahan yang paling umum
- **Jangan membuat menu terlalu bertingkat** — Batasi menu hingga maksimum 2-3 tingkat
- **Jangan gunakan label yang ambigu** — "Proses" dibandingkan dengan "Proses Dokumen"
- **Jangan membuatnya terlalu rumit** — Buat menu tetap sederhana dan terfokus
- **Jangan mencampur metafora** — Gunakan penamaan dan penataan yang konsisten

## Catatan Khusus Platform

### macOS

- Menu aplikasi ditambahkan secara otomatis dengan nama aplikasi
- Gunakan `Cmd`, bukan `Ctrl`, untuk akselerator
- Secara bawaan, "Tentang", "Preferensi", dan "Keluar" berada di menu aplikasi

### Windows/Linux

- Tidak ada menu aplikasi otomatis
- Gunakan `Ctrl` untuk tombol pintasan
- "Keluar" biasanya berada di menu File
