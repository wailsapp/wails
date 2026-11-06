---
title: "API Aplikasi"
description: "Referensi lengkap untuk API Aplikasi"
slug: "reference/application"
sourcePath: "reference/application.md"
---

## Ikhtisar

`Application` merupakan inti aplikasi Wails Anda. Komponen ini mengelola jendela, layanan, dan peristiwa, serta menyediakan akses ke semua fitur platform.

## Membuat Aplikasi

```go
import "github.com/wailsapp/wails/v3/pkg/application"

app := application.New(application.Options{
    Name:        "My App",
    Description: "My awesome application",
    Services: []application.Service{
        application.NewService(&MyService{}),
    },
})
```

## Metode Inti

### Run()

Memulai loop peristiwa aplikasi.

```go
func (a *App) Run() error
```

**Contoh:**

```go
err := app.Run()
if err != nil {
    log.Fatal(err)
}
```

**Nilai kembalian:** Error jika proses awal gagal

### Quit()

Menghentikan aplikasi dengan aman.

```go
func (a *App) Quit()
```

**Contoh:**

```go
// In a menu handler
menu.Add("Quit").OnClick(func(ctx *application.Context) {
    app.Quit()
})
```

### Config()

Mengembalikan konfigurasi aplikasi.

```go
func (a *App) Config() Options
```

**Contoh:**

```go
config := app.Config()
fmt.Println("App name:", config.Name)
```

## Pengelolaan Jendela

### app.Window.New()

Membuat jendela webview baru dengan opsi default.

```go
func (wm *WindowManager) New() *WebviewWindow
```

**Contoh:**

```go
window := app.Window.New()
window.Show()
```

### app.Window.NewWithOptions()

Membuat jendela webview baru dengan opsi khusus.

```go
func (wm *WindowManager) NewWithOptions(options WebviewWindowOptions) *WebviewWindow
```

**Contoh:**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "My Window",
    Width:  800,
    Height: 600,
    BackgroundColour: application.NewRGB(255, 255, 255),
})
```

### app.Window.GetByName()

Mengambil jendela berdasarkan namanya. Mengembalikan jendela tersebut dan status apakah jendela ditemukan.

```go
func (wm *WindowManager) GetByName(name string) (Window, bool)
```

**Contoh:**

```go
if window, ok := app.Window.GetByName("main"); ok {
    window.Show()
}
```

### app.Window.GetAll()

Mengembalikan semua jendela aplikasi.

```go
func (wm *WindowManager) GetAll() []Window
```

**Contoh:**

```go
windows := app.Window.GetAll()
for _, window := range windows {
    fmt.Println("Window:", window.Name())
}
```

## Pengelola

Application menyediakan akses ke berbagai pengelola melalui properti:

```go
app.Window       // Window management
app.Menu         // Menu management
app.Dialog       // Dialog management
app.Event        // Event management
app.Clipboard    // Clipboard operations
app.Screen       // Screen information
app.SystemTray   // System tray
app.Browser      // Browser operations
app.Env          // Environment variables
app.ContextMenu  // Context-menu management
app.KeyBinding   // Global keyboard shortcuts
app.Logger       // *slog.Logger
```

### Contoh Penggunaan

```go
// Create window
window := app.Window.New()

// Show dialog
app.Dialog.Info().SetMessage("Hello!").Show()

// Copy to clipboard
app.Clipboard.SetText("Copied text")

// Get screens
screens := app.Screen.GetAll()
```

## Pengelolaan Layanan

### RegisterService()

Mendaftarkan layanan pada aplikasi.

```go
func (a *App) RegisterService(service Service)
```

`RegisterService` tidak mengembalikan nilai; kesalahan inisialisasi layanan muncul melalui kegagalan `ServiceStartup` selama `app.Run()`.

**Contoh:**

```go
type MyService struct {
    app *application.App
}

func NewMyService(app *application.App) *MyService {
    return &MyService{app: app}
}

// Register after app creation
app.RegisterService(application.NewService(NewMyService(app)))
```

## Pengelolaan Peristiwa

### app.Event.Emit()

Memancarkan peristiwa khusus. Mengembalikan `true` jika sebuah hook membatalkan pemancaran tersebut.

```go
func (em *EventManager) Emit(name string, data ...any) bool
```

**Contoh:**

```go
// Emit event with data
app.Event.Emit("user-logged-in", map[string]interface{}{
    "username": "john",
    "timestamp": time.Now(),
})
```

### app.Event.On()

Memantau peristiwa khusus. Mengembalikan `func()` untuk berhenti berlangganan.

```go
func (em *EventManager) On(name string, callback func(*CustomEvent)) func()
```

**Contoh:**

```go
app.Event.On("user-logged-in", func(e *application.CustomEvent) {
    data := e.Data.(map[string]interface{})
    username := data["username"].(string)
    fmt.Println("User logged in:", username)
})
```

### app.Event.OnApplicationEvent()

Memantau peristiwa siklus hidup aplikasi. Parameter `eventType` bertipe `events.ApplicationEventType` (dari paket `events`).

```go
func (em *EventManager) OnApplicationEvent(
    eventType events.ApplicationEventType,
    callback func(*ApplicationEvent),
) func()
```

**Contoh:**

```go
import "github.com/wailsapp/wails/v3/pkg/events"

// Listen for app-started
app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(e *application.ApplicationEvent) {
    fmt.Println("Application started")
})

// Application shutdown is NOT an event constant; register cleanup via:
app.OnShutdown(func() {
    fmt.Println("Application shutting down")
})
```

## Metode Dialog

Dialog diakses melalui pengelola `app.Dialog`. Lihat [API Dialog](/reference/dialogs/) untuk referensi lengkap.

### Dialog Pesan

```go
// Information dialog
app.Dialog.Info().
    SetTitle("Success").
    SetMessage("Operation completed!").
    Show()

// Error dialog
app.Dialog.Error().
    SetTitle("Error").
    SetMessage("Something went wrong.").
    Show()

// Warning dialog
app.Dialog.Warning().
    SetTitle("Warning").
    SetMessage("This action cannot be undone.").
    Show()
```

### Dialog Pertanyaan

Dialog pertanyaan menggunakan callback tombol untuk menangani respons pengguna:

```go
dialog := app.Dialog.Question().
    SetTitle("Confirm").
    SetMessage("Continue?")

yes := dialog.AddButton("Yes")
yes.OnClick(func() {
    // Handle yes
})

no := dialog.AddButton("No")
no.OnClick(func() {
    // Handle no
})

dialog.SetDefaultButton(yes)
dialog.SetCancelButton(no)
dialog.Show()
```

### Dialog File

```go
// Open file dialog
path, err := app.Dialog.OpenFile().
    SetTitle("Select File").
    AddFilter("Images", "*.png;*.jpg").
    PromptForSingleSelection()

// Save file dialog
path, err := app.Dialog.SaveFile().
    SetTitle("Save File").
    SetFilename("document.pdf").
    AddFilter("PDF", "*.pdf").
    PromptForSingleSelection()

// Folder selection (use OpenFile with directory options)
path, err := app.Dialog.OpenFile().
    SetTitle("Select Folder").
    CanChooseDirectories(true).
    CanChooseFiles(false).
    PromptForSingleSelection()
```

## Pencatat Log

Aplikasi menyediakan pencatat log terstruktur:

```go
app.Logger.Info("Message", "key", "value")
app.Logger.Error("Error occurred", "error", err)
app.Logger.Debug("Debug info")
app.Logger.Warn("Warning message")
```

**Contoh:**

```go
func (s *MyService) ProcessData(data string) error {
    s.app.Logger.Info("Processing data", "length", len(data))
    
    if err := process(data); err != nil {
        s.app.Logger.Error("Processing failed", "error", err)
        return err
    }
    
    s.app.Logger.Info("Processing complete")
    return nil
}
```

## Penanganan Pesan Mentah

Untuk aplikasi yang memerlukan kontrol langsung tingkat rendah atas komunikasi dari frontend ke backend, Wails menyediakan opsi `RawMessageHandler`. Opsi ini melewati sistem binding standar.

@note{type="info"}
Pesan mentah sebaiknya hanya digunakan sebagai pilihan terakhir. Sistem binding standar sangat dioptimalkan dan memadai untuk hampir semua aplikasi. Gunakan pesan mentah hanya jika Anda telah melakukan profiling terhadap aplikasi dan memastikan bahwa binding menjadi hambatan performa.

@end

### RawMessageHandler

`RawMessageHandler` adalah field pada `application.Options`, bukan metode. Runtime memanggilnya untuk setiap pesan mentah yang dikirim dari frontend melalui `System.invoke()`.

```go
type Options struct {
    // ... other fields ...
    RawMessageHandler func(window Window, message string, originInfo *OriginInfo)
}
```

`OriginInfo` memuat `Origin`, `TopOrigin`, dan `IsMainFrame` (setiap platform mengisi subset yang berbeda—lihat [Panduan Pesan Mentah](/guides/raw-messages/) untuk matriks per platform).

**Contoh:**

```go
app := application.New(application.Options{
    Name: "My App",
    RawMessageHandler: func(window application.Window, message string, originInfo *application.OriginInfo) {
        // Handle the raw message
        fmt.Printf("Received from %s (%s): %s\n", window.Name(), originInfo.Origin, message)

        // You can respond using events
        window.EmitEvent("response", processMessage(message))
    },
})
```

Untuk detail selengkapnya, lihat [Panduan Pesan Mentah](/guides/raw-messages/).

## Opsi Khusus Platform

### Opsi Windows

Konfigurasikan perilaku khusus Windows pada tingkat aplikasi:

```go
app := application.New(application.Options{
    Name: "My App",
    Windows: application.WindowsOptions{
        // WebView2 browser flags (apply to ALL windows)
        EnabledFeatures:       []string{"msWebView2EnableDraggableRegions"},
        DisabledFeatures:      []string{"msExperimentalFeature"},
        AdditionalBrowserArgs: []string{"--remote-debugging-port=9222"},

        // Other Windows options
        WndClass:                      "MyAppClass",
        WebviewUserDataPath:           "",  // Default: %APPDATA%\[BinaryName.exe]
        WebviewBrowserPath:            "",  // Default: system WebView2
        DisableQuitOnLastWindowClosed: false,
    },
})
```

**Flag Browser:**

- `EnabledFeatures` - Flag fitur WebView2 yang akan diaktifkan
- `DisabledFeatures` - Flag fitur WebView2 yang akan dinonaktifkan
- `AdditionalBrowserArgs` - Argumen baris perintah Chromium

Lihat [Opsi Jendela - Opsi Windows Tingkat Aplikasi](/features/windows/options/#application-level-windows-options) untuk dokumentasi terperinci.

### Opsi Mac

```go
app := application.New(application.Options{
    Name: "My App",
    Mac: application.MacOptions{
        ActivationPolicy: application.ActivationPolicyRegular,
        ApplicationShouldTerminateAfterLastWindowClosed: true,
    },
})
```

### Opsi Linux

```go
app := application.New(application.Options{
    Name: "My App",
    Linux: application.LinuxOptions{
        ProgramName:                   "my-app",
        DisableQuitOnLastWindowClosed: false,
    },
})
```

## Contoh Aplikasi Lengkap

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name:        "My Application",
        Description: "A demo application",
        Mac: application.MacOptions{
            ApplicationShouldTerminateAfterLastWindowClosed: true,
        },
    })

    // Create main window
    window := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:            "My App",
        Width:            1024,
        Height:           768,
        MinWidth:         800,
        MinHeight:        600,
        BackgroundColour: application.NewRGB(255, 255, 255),
        URL:              "http://wails.localhost/",
    })

    window.Center()
    window.Show()

    app.Run()
}
```
