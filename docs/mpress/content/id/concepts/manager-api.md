---
title: "API Pengelola"
description: "Struktur API yang tertata dengan antarmuka pengelola yang berfokus pada fungsi tertentu"
slug: "concepts/manager-api"
sourcePath: "concepts/manager-api.md"
---

API Pengelola Wails v3 menyediakan cara yang tertata dan mudah ditemukan untuk mengakses fungsionalitas aplikasi melalui struct pengelola yang berfokus pada fungsi tertentu dan dikelompokkan dalam field publik pada `*application.App`. Wails 3 sepenuhnya berbeda dari v2—tidak ada lapisan pembungkus per panggilan untuk mempertahankan kompatibilitas dengan API lama bergaya `app.NewWebviewWindow(...)`, sehingga pengelola di bawah ini merupakan cara untuk mengendalikan aplikasi.

## Ikhtisar

API Pengelola menata fungsionalitas aplikasi ke dalam dua belas area yang berfokus pada fungsi tertentu (satu pencatat log dan sebelas pengelola):

- **`app.Window`** - Pembuatan dan pengelolaan jendela serta callback
- **`app.ContextMenu`** - Pendaftaran dan pengelolaan menu konteks\
- **`app.KeyBinding`** - Pengelolaan pengikatan tombol global
- **`app.Browser`** - Integrasi browser (membuka URL dan file)
- **`app.Env`** - Informasi lingkungan dan status sistem
- **`app.Dialog`** - Operasi dialog file dan pesan
- **`app.Event`** - Penanganan peristiwa khusus dan peristiwa aplikasi
- **`app.Menu`** - Pengelolaan menu aplikasi
- **`app.Screen`** - Pengelolaan layar dan transformasi koordinat
- **`app.Clipboard`** - Operasi teks papan klip
- **`app.SystemTray`** - Pembuatan dan pengelolaan ikon baki sistem
- **`app.Autostart`** - Mendaftarkan aplikasi agar dijalankan saat pengguna masuk

## Manfaat

- **Lebih mudah ditemukan** - Pelengkapan otomatis IDE menampilkan cakupan API yang tertata
- **Penataan kode yang lebih baik** - Metode terkait dikelompokkan bersama
- **Kemudahan pemeliharaan yang lebih baik** - Pemisahan tanggung jawab antar-pengelola
- **Kemudahan perluasan di masa mendatang** - Fitur baru lebih mudah ditambahkan ke area tertentu

## Penggunaan

API Pengelola menyediakan akses yang tertata ke seluruh fungsionalitas aplikasi:

```go
// Events and custom event handling
app.Event.Emit("custom", data)
app.Event.On("custom", func(e *CustomEvent) { ... })

// Window management
window, _ := app.Window.GetByName("main")
app.Window.OnCreate(func(window Window) { ... })

// Browser integration
app.Browser.OpenURL("https://wails.io")

// Menu management
menu := app.Menu.New()
app.Menu.Set(menu)

// System tray
systray := app.SystemTray.New()
```

## Referensi Pengelola

### Pengelola Jendela

Mengelola pembuatan dan pengambilan jendela serta callback siklus hidupnya.

```go
// Create windows
window := app.Window.New()
window := app.Window.NewWithOptions(options)
current := app.Window.Current()

// Find windows
window, exists := app.Window.GetByName("main")
windows := app.Window.GetAll()

// Window callbacks
app.Window.OnCreate(func(window Window) {
    // Handle window creation
})
```

### Pengelola Peristiwa

Menangani peristiwa khusus dan pemantauan peristiwa aplikasi.

```go
// Custom events
app.Event.Emit("userAction", data)
cancelFunc := app.Event.On("userAction", func(e *CustomEvent) {
    // Handle event
})
app.Event.Off("userAction")
app.Event.Reset() // Remove all listeners

// Application events
app.Event.OnApplicationEvent(events.Common.ThemeChanged, func(e *ApplicationEvent) {
    // Handle system theme change
})
```

### Pengelola Browser

Menyediakan integrasi browser untuk membuka URL dan file.

```go
// Open URLs and files in default browser
err := app.Browser.OpenURL("https://wails.io")
err := app.Browser.OpenFile("/path/to/document.pdf")
```

### Pengelola Lingkungan

Menyediakan akses ke informasi lingkungan sistem.

```go
// Get environment info
env := app.Env.Info()
fmt.Printf("OS: %s, Arch: %s\n", env.OS, env.Arch)

// Check system theme
if app.Env.IsDarkMode() {
    // Dark mode is active
}

// Open file manager
err := app.Env.OpenFileManager("/path/to/folder", false)
```

### Pengelola Dialog

Menyediakan akses yang tertata ke dialog file dan pesan.

```go
// File dialogs
result, err := app.Dialog.OpenFile().
    AddFilter("Text Files", "*.txt").
    PromptForSingleSelection()

result, err = app.Dialog.SaveFile().
    SetFilename("document.txt").
    PromptForSingleSelection()

// Message dialogs
app.Dialog.Info().
    SetTitle("Information").
    SetMessage("Operation completed successfully").
    Show()

app.Dialog.Error().
    SetTitle("Error").
    SetMessage("An error occurred").
    Show()
```

### Pengelola Menu

Pembuatan dan pengelolaan menu aplikasi.

```go
// Create and set application menu
menu := app.Menu.New()
fileMenu := menu.AddSubmenu("File")
fileMenu.Add("New").OnClick(func(ctx *Context) {
    // Handle menu click
})

app.Menu.Set(menu)

// Show about dialog
app.Menu.ShowAbout()
```

### Pengelola Pengikatan Tombol

Pengelolaan dinamis untuk pengikatan tombol global.

```go
// Add key bindings
app.KeyBinding.Add("ctrl+n", func(window application.Window) {
    // Handle Ctrl+N
})

app.KeyBinding.Add("ctrl+q", func(window application.Window) {
    app.Quit()
})

// Remove key bindings
app.KeyBinding.Remove("ctrl+n")

// Get all bindings
bindings := app.KeyBinding.GetAll()
```

### Pengelola Menu Konteks

Pengelolaan menu konteks tingkat lanjut (untuk pembuat pustaka).

```go
// Create and register context menu
menu := app.ContextMenu.New()
app.ContextMenu.Add("myMenu", menu)

// Retrieve context menu
menu, exists := app.ContextMenu.Get("myMenu")

// Remove context menu
app.ContextMenu.Remove("myMenu")
```

### Pengelola Layar

Pengelolaan layar dan transformasi koordinat untuk konfigurasi multi-monitor.

```go
// Get screen information
screens := app.Screen.GetAll()
primary := app.Screen.GetPrimary()

// Coordinate transformations
physicalPoint := app.Screen.DipToPhysicalPoint(logicalPoint)
logicalPoint := app.Screen.PhysicalToDipPoint(physicalPoint)

// Screen detection
screen := app.Screen.ScreenNearestDipPoint(point)
screen = app.Screen.ScreenNearestDipRect(rect)
```

### Pengelola Papan Klip

Operasi papan klip untuk membaca dan menulis teks.

```go
// Set text to clipboard
success := app.Clipboard.SetText("Hello World")
if !success {
    // Handle error
}

// Get text from clipboard
text, ok := app.Clipboard.Text()
if !ok {
    // Handle error
} else {
    // Use the text
}
```

### Pengelola SystemTray

Pembuatan dan pengelolaan ikon baki sistem.

```go
// Create system tray
systray := app.SystemTray.New()
systray.SetLabel("My App")
systray.SetIcon(iconBytes)

// Add menu to system tray
menu := app.Menu.New()
menu.Add("Open").OnClick(func(ctx *Context) {
    // Handle click
})
systray.SetMenu(menu)

// Destroy system tray when done
systray.Destroy()
```

### Pengelola Autostart

Mendaftarkan aplikasi agar dijalankan saat pengguna masuk. Memilih mekanisme native yang sesuai untuk setiap platform: SMAppService atau plist LaunchAgent di macOS, kunci registry `HKCU\…\Run` di Windows, dan entri XDG `.desktop` di Linux.

```go
// Register to launch at login
err := app.Autostart.Enable()

// With extra launch-time arguments and a custom identifier
err = app.Autostart.EnableWithOptions(application.AutostartOptions{
    Identifier: "com.example.myapp",
    Arguments:  []string{"--hidden"},
})

// Check / remove
enabled, err := app.Autostart.IsEnabled()
status, err := app.Autostart.Status()  // includes Path + Strategy
err = app.Autostart.Disable()
```

Lihat [halaman fitur Autostart](/features/autostart/basics/) untuk mengetahui perilaku pada setiap platform, aturan pengidentifikasi, dan jaminan pendeteksian data usang.
