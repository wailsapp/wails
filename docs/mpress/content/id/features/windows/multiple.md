---
title: "Beberapa Jendela"
description: "Pola dan praktik terbaik untuk aplikasi dengan beberapa jendela"
slug: "features/windows/multiple"
sourcePath: "features/windows/multiple.md"
---

## Aplikasi dengan Beberapa Jendela

Wails v3 menyediakan **dukungan bawaan untuk beberapa jendela** guna membuat jendela pengaturan, jendela dokumen, palet alat, dan jendela pemeriksa. Lacak jendela, aktifkan komunikasi antarjendela, dan kelola siklus hidupnya dengan API yang sederhana dan konsisten.

### Jendela Utama + Pengaturan

```go
package main

import "github.com/wailsapp/wails/v3/pkg/application"

type App struct {
    app            *application.App
    mainWindow     *application.WebviewWindow
    settingsWindow *application.WebviewWindow
}

func main() {
    app := &App{}
    
    app.app = application.New(application.Options{
        Name: "Multi-Window App",
    })
    
    // Create main window
    app.mainWindow = app.app.Window.NewWithOptions(application.WebviewWindowOptions{
        Name:   "main",
        Title:  "Main Application",
        Width:  1200,
        Height: 800,
    })
    
    // Create settings window (hidden initially)
    app.settingsWindow = app.app.Window.NewWithOptions(application.WebviewWindowOptions{
        Name:   "settings",
        Title:  "Settings",
        Width:  600,
        Height: 400,
        Hidden: true,
    })
    
    app.app.Run()
}

// Show settings from main window
func (a *App) ShowSettings() {
    if a.settingsWindow != nil {
        a.settingsWindow.Show()
        a.settingsWindow.Focus()
    }
}
```

**Poin utama:**

- Jendela utama selalu terlihat
- Jendela pengaturan dibuat tetapi disembunyikan
- Tampilkan pengaturan saat diperlukan
- Gunakan kembali jendela yang sama (jangan buat beberapa jendela)

## Pelacakan Jendela

### Mendapatkan Semua Jendela

```go
windows := app.Window.GetAll()
fmt.Printf("Total windows: %d\n", len(windows))

for _, window := range windows {
    fmt.Printf("- %s (ID: %d)\n", window.Name(), window.ID())
}
```

### Menemukan Jendela Tertentu

```go
// By name
if settings, ok := app.Window.GetByName("settings"); ok {
    settings.Show()
}

// By ID
if window, ok := app.Window.GetByID(123); ok {
    window.Focus()
}

// Current (focused) window
current := app.Window.Current()
```

### Pola Registri Jendela

Lacak jendela dalam aplikasi Anda:

```go
type WindowManager struct {
    windows map[string]*application.WebviewWindow
    mu      sync.RWMutex
}

func (wm *WindowManager) Register(name string, window *application.WebviewWindow) {
    wm.mu.Lock()
    defer wm.mu.Unlock()
    wm.windows[name] = window
}

func (wm *WindowManager) Get(name string) *application.WebviewWindow {
    wm.mu.RLock()
    defer wm.mu.RUnlock()
    return wm.windows[name]
}

func (wm *WindowManager) Remove(name string) {
    wm.mu.Lock()
    defer wm.mu.Unlock()
    delete(wm.windows, name)
}
```

## Komunikasi Antarjendela

### Menggunakan Peristiwa

Jendela berkomunikasi melalui sistem peristiwa:

```go
// In main window - emit event
app.Event.Emit("settings-changed", map[string]interface{}{
    "theme": "dark",
    "fontSize": 14,
})

// In settings window - listen for event
app.Event.On("settings-changed", func(event *application.CustomEvent) {
    data := event.Data.(map[string]interface{})
    theme := data["theme"].(string)
    fontSize := data["fontSize"].(int)
    
    // Update UI
    updateSettings(theme, fontSize)
})
```

### Pola Status Bersama

Gunakan pengelola status bersama:

```go
type AppState struct {
    theme    string
    fontSize int
    mu       sync.RWMutex
}

var state = &AppState{
    theme:    "light",
    fontSize: 12,
}

func (s *AppState) SetTheme(theme string) {
    s.mu.Lock()
    s.theme = theme
    s.mu.Unlock()
    
    // Notify all windows
    app.Event.Emit("theme-changed", theme)
}

func (s *AppState) GetTheme() string {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.theme
}
```

### Pesan Antarjendela

Kirim pesan di antara jendela tertentu:

```go
// Get target window
if targetWindow, ok := app.Window.GetByName("preview"); ok {
    // Emit event to specific window
    targetWindow.EmitEvent("update-preview", previewData)
}
```

## Pola Umum

### Pola 1: Jendela Singleton

Pastikan hanya ada satu instans jendela:

```go
var settingsWindow *application.WebviewWindow

func ShowSettings(app *application.App) {
    // Create if doesn't exist
    if settingsWindow == nil {
        settingsWindow = app.Window.NewWithOptions(application.WebviewWindowOptions{
            Name:   "settings",
            Title:  "Settings",
            Width:  600,
            Height: 400,
        })
        
        // Cleanup on close
        settingsWindow.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
            settingsWindow = nil
        })
    }
    
    // Show and focus
    settingsWindow.Show()
    settingsWindow.Focus()
}
```

### Pola 2: Jendela Dokumen

Beberapa instans dari jenis jendela yang sama:

```go
type DocumentWindow struct {
    window   *application.WebviewWindow
    filePath string
    modified bool
}

var documents = make(map[string]*DocumentWindow)

func OpenDocument(app *application.App, filePath string) {
    // Check if already open
    if doc, exists := documents[filePath]; exists {
        doc.window.Show()
        doc.window.Focus()
        return
    }
    
    // Create new document window
    window := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:  filepath.Base(filePath),
        Width:  800,
        Height: 600,
    })
    
    doc := &DocumentWindow{
        window:   window,
        filePath: filePath,
        modified: false,
    }
    
    documents[filePath] = doc
    
    // Cleanup on close
    window.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
        delete(documents, filePath)
    })
    
    // Load document
    loadDocument(window, filePath)
}
```

### Pola 3: Palet Alat

Jendela mengambang yang tetap berada di atas:

```go
func CreateToolPalette(app *application.App) *application.WebviewWindow {
    palette := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Name:          "tools",
        Title:         "Tools",
        Width:         200,
        Height:        400,
        AlwaysOnTop:   true,
        DisableResize: true,
    })

    return palette
}
```

### Pola 4: Dialog modal (khusus macOS)

Jendela anak yang memblokir jendela induk:

```go
func ShowModaldialog(parent *application.WebviewWindow, title string) {
    dialog := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:         title,
        Width:         400,
        Height:        200,
        AlwaysOnTop:   true,
        DisableResize: true,
    })

    parent.AttachModal(dialog)
}
```

### Pola 5: Jendela Pemeriksa/Pratinjau

Jendela tertaut yang diperbarui secara bersamaan:

```go
type EditorApp struct {
    editor  *application.WebviewWindow
    preview *application.WebviewWindow
}

func (e *EditorApp) UpdatePreview(content string) {
    if e.preview != nil && e.preview.IsVisible() {
        e.preview.EmitEvent("content-changed", content)
    }
}

func (e *EditorApp) TogglePreview() {
    if e.preview == nil {
        e.preview = app.Window.NewWithOptions(application.WebviewWindowOptions{
            Name:   "preview",
            Title:  "Preview",
            Width:  600,
            Height: 800,
        })
        
        e.preview.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
            e.preview = nil
        })
    }
    
    if e.preview.IsVisible() {
        e.preview.Hide()
    } else {
        e.preview.Show()
    }
}
```

## Hubungan Induk-Anak

### Membuat Jendela Anak

```go
childWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "Child Window",
})

parentWindow.AttachModal(childWindow)
```

**Perilaku:**

- Jendela anak tetap berada di atas jendela induk
- Jendela anak bergerak bersama jendela induk
- Jendela anak memblokir interaksi dengan jendela induk

**Dukungan platform:**

| macOS | Windows | Linux |
| --- | --- | --- |
| ✅ | ❌ | ❌ |

### Perilaku Modal

Buat perilaku menyerupai modal:

```go
func ShowModaldialog(parent *application.WebviewWindow, title string) {
    dialog := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:       title,
        Width:       400,
        Height:      200,
    })
    
    parent.AttachModal(dialog)
}
```

## Pengelolaan Siklus Hidup Jendela

### Callback Pembuatan

Dapatkan notifikasi saat jendela dibuat:

```go
app.Window.OnCreate(func(window application.Window) {
    fmt.Printf("Window created: %s\n", window.Name())

    // Configure all new windows
    window.SetMinSize(400, 300)
})
```

### Callback Penghancuran

Lakukan pembersihan saat jendela dihancurkan:

```go
window.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
    fmt.Printf("Window %s is closing\n", window.Name())
    
    // Cleanup resources
    cleanup(window.ID())
    
    // Remove from tracking
    removeFromRegistry(window.Name())
})
```

### Perilaku Keluar Aplikasi

Kendalikan kapan aplikasi ditutup:

```go
app := application.New(application.Options{
    Mac: application.MacOptions{
        // Don't quit when last window closes
        ApplicationShouldTerminateAfterLastWindowClosed: false,
    },
})
```

**Kasus penggunaan:**

- Aplikasi baki sistem
- Layanan latar belakang
- Aplikasi bilah menu (macOS)

## Pengelolaan Memori

### Mencegah Kebocoran

Selalu bersihkan referensi jendela:

```go
var windows = make(map[string]*application.WebviewWindow)

func CreateWindow(name string) {
    window := app.Window.New()
    windows[name] = window
    
    // IMPORTANT: Clean up on close
    window.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
        delete(windows, name)
    })
}
```

### Menutup Jendela

```go
// Close — dispatches WindowClosing; a RegisterHook can call e.Cancel().
window.Close()
```

Tidak ada `window.Destroy()` di v3. Untuk mencegah penutupan, gunakan `RegisterHook(events.Common.WindowClosing, ...)` dan panggil `event.Cancel()`.

### Pembersihan Sumber Daya

```go
type ManagedWindow struct {
    window     *application.WebviewWindow
    resources  []io.Closer
}

func (mw *ManagedWindow) Destroy() {
    // Close all resources
    for _, resource := range mw.resources {
        resource.Close()
    }

    // Close the window — WindowClosing listeners run before the window goes away.
    mw.window.Close()
}
```

## Pola Lanjutan

### Pool Jendela

Gunakan kembali jendela alih-alih membuat yang baru:

```go
type WindowPool struct {
    available []*application.WebviewWindow
    inUse     map[uint]*application.WebviewWindow
    mu        sync.Mutex
}

func (wp *WindowPool) Acquire() *application.WebviewWindow {
    wp.mu.Lock()
    defer wp.mu.Unlock()
    
    // Reuse available window
    if len(wp.available) > 0 {
        window := wp.available[0]
        wp.available = wp.available[1:]
        wp.inUse[window.ID()] = window
        return window
    }
    
    // Create new window
    window := app.Window.New()
    wp.inUse[window.ID()] = window
    return window
}

func (wp *WindowPool) Release(window *application.WebviewWindow) {
    wp.mu.Lock()
    defer wp.mu.Unlock()
    
    delete(wp.inUse, window.ID())
    window.Hide()
    wp.available = append(wp.available, window)
}
```

### Grup Jendela

Kelola jendela yang saling terkait secara bersamaan:

```go
type WindowGroup struct {
    name    string
    windows []*application.WebviewWindow
}

func (wg *WindowGroup) Add(window *application.WebviewWindow) {
    wg.windows = append(wg.windows, window)
}

func (wg *WindowGroup) ShowAll() {
    for _, window := range wg.windows {
        window.Show()
    }
}

func (wg *WindowGroup) HideAll() {
    for _, window := range wg.windows {
        window.Hide()
    }
}

func (wg *WindowGroup) CloseAll() {
    for _, window := range wg.windows {
        window.Close()
    }
}
```

### Pengelolaan Ruang Kerja

Simpan dan pulihkan tata letak jendela:

```go
type WindowLayout struct {
    Windows []WindowState `json:"windows"`
}

type WindowState struct {
    Name   string `json:"name"`
    X      int    `json:"x"`
    Y      int    `json:"y"`
    Width  int    `json:"width"`
    Height int    `json:"height"`
}

func SaveLayout() *WindowLayout {
    layout := &WindowLayout{}
    
    for _, window := range app.Window.GetAll() {
        x, y := window.Position()
        width, height := window.Size()
        
        layout.Windows = append(layout.Windows, WindowState{
            Name:   window.Name(),
            X:      x,
            Y:      y,
            Width:  width,
            Height: height,
        })
    }
    
    return layout
}

func RestoreLayout(layout *WindowLayout) {
    for _, state := range layout.Windows {
        if window, ok := app.Window.GetByName(state.Name); ok {
            window.SetPosition(state.X, state.Y)
            window.SetSize(state.Width, state.Height)
        }
    }
}
```

## Contoh Lengkap

Berikut adalah aplikasi multijendela yang siap digunakan dalam produksi:

```go
package main

import (
    "encoding/json"
    "os"
    "sync"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)

type MultiWindowApp struct {
    app     *application.App
    windows map[string]*application.WebviewWindow
    mu      sync.RWMutex
}

func main() {
    mwa := &MultiWindowApp{
        windows: make(map[string]*application.WebviewWindow),
    }
    
    mwa.app = application.New(application.Options{
        Name: "Multi-Window Application",
        Mac: application.MacOptions{
            ApplicationShouldTerminateAfterLastWindowClosed: false,
        },
    })
    
    // Create main window
    mwa.CreateMainWindow()
    
    // Load saved layout
    mwa.LoadLayout()
    
    mwa.app.Run()
}

func (mwa *MultiWindowApp) CreateMainWindow() {
    window := mwa.app.Window.NewWithOptions(application.WebviewWindowOptions{
        Name:   "main",
        Title:  "Main Application",
        Width:  1200,
        Height: 800,
    })
    
    mwa.RegisterWindow("main", window)
}

func (mwa *MultiWindowApp) ShowSettings() {
    if window := mwa.GetWindow("settings"); window != nil {
        window.Show()
        window.Focus()
        return
    }
    
    window := mwa.app.Window.NewWithOptions(application.WebviewWindowOptions{
        Name:   "settings",
        Title:  "Settings",
        Width:  600,
        Height: 400,
    })
    
    mwa.RegisterWindow("settings", window)
}

func (mwa *MultiWindowApp) OpenDocument(path string) {
    name := "doc-" + path
    
    if window := mwa.GetWindow(name); window != nil {
        window.Show()
        window.Focus()
        return
    }
    
    window := mwa.app.Window.NewWithOptions(application.WebviewWindowOptions{
        Name:  name,
        Title: path,
        Width: 800,
        Height: 600,
    })
    
    mwa.RegisterWindow(name, window)
}

func (mwa *MultiWindowApp) RegisterWindow(name string, window *application.WebviewWindow) {
    mwa.mu.Lock()
    mwa.windows[name] = window
    mwa.mu.Unlock()
    
    window.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
        mwa.UnregisterWindow(name)
    })
}

func (mwa *MultiWindowApp) UnregisterWindow(name string) {
    mwa.mu.Lock()
    delete(mwa.windows, name)
    mwa.mu.Unlock()
}

func (mwa *MultiWindowApp) GetWindow(name string) *application.WebviewWindow {
    mwa.mu.RLock()
    defer mwa.mu.RUnlock()
    return mwa.windows[name]
}

func (mwa *MultiWindowApp) SaveLayout() {
    layout := make(map[string]WindowState)
    
    mwa.mu.RLock()
    for name, window := range mwa.windows {
        x, y := window.Position()
        width, height := window.Size()
        
        layout[name] = WindowState{
            X:      x,
            Y:      y,
            Width:  width,
            Height: height,
        }
    }
    mwa.mu.RUnlock()
    
    data, _ := json.Marshal(layout)
    os.WriteFile("layout.json", data, 0644)
}

func (mwa *MultiWindowApp) LoadLayout() {
    data, err := os.ReadFile("layout.json")
    if err != nil {
        return
    }
    
    var layout map[string]WindowState
    if err := json.Unmarshal(data, &layout); err != nil {
        return
    }
    
    for name, state := range layout {
        if window := mwa.GetWindow(name); window != nil {
            window.SetPosition(state.X, state.Y)
            window.SetSize(state.Width, state.Height)
        }
    }
}

type WindowState struct {
    X      int `json:"x"`
    Y      int `json:"y"`
    Width  int `json:"width"`
    Height int `json:"height"`
}
```

## Praktik Terbaik

### ✅ Lakukan

- **Lacak jendela** - Simpan referensi agar mudah diakses
- **Bersihkan saat dihancurkan** - Cegah kebocoran memori
- **Gunakan event untuk komunikasi** - Arsitektur yang tidak saling bergantung
- **Gunakan kembali jendela** - Jangan buat duplikat
- **Simpan/pulihkan tata letak** - UX yang lebih baik
- **Tangani penutupan jendela** - Minta konfirmasi sebelum menutup jika ada data yang belum disimpan

### ❌ Jangan Lakukan

- **Jangan buat jendela tanpa batas** - Menimbulkan masalah memori dan performa
- **Jangan lupa melakukan pembersihan** - Menyebabkan kebocoran memori
- **Jangan gunakan variabel global secara sembarangan** - Menimbulkan masalah keamanan thread
- **Jangan blokir pembuatan jendela** - Buat secara asinkron jika diperlukan
- **Jangan abaikan perbedaan antarplatform** - Uji di semua platform

## Langkah Berikutnya

- [Dasar-Dasar Jendela](/features/windows/basics/) - Pelajari dasar-dasar pengelolaan jendela
- [Event Jendela](/features/windows/events/) - Tangani event siklus hidup jendela
- [Sistem Event](/features/events/system/) - Pelajari sistem event secara mendalam
- [Jendela Tanpa Bingkai](/features/windows/frameless/) - Buat dekorasi jendela khusus

---

**Ada pertanyaan?** Tanyakan di [Discord](https://discord.gg/JDdSxwjhGf) atau lihat [contoh multijendela](https://github.com/wailsapp/wails/tree/master/v3/examples/multi-window).
