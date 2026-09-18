---
title: "Informasi Layar"
description: "Dapatkan informasi tentang layar dan monitor"
slug: "features/screens/info"
sourcePath: "features/screens/info.md"
---

## Informasi Layar

Wails menyediakan **API layar terpadu** yang berfungsi di semua platform. Dapatkan informasi layar, deteksi beberapa monitor, kueri properti layar (ukuran, posisi, DPI), identifikasi layar utama, dan tangani penskalaan DPI dengan kode yang konsisten.

## Mulai Cepat

```go
// Get all screens
screens := app.Screen.GetAll()

for _, screen := range screens {
    fmt.Printf("Screen: %s (%dx%d)\n", 
        screen.Name, screen.Size.Width, screen.Size.Height)
}

// Get primary screen
primary := app.Screen.GetPrimary()
fmt.Printf("Primary: %s\n", primary.Name)
```

**Selesai!** Informasi layar lintas platform.

## Mendapatkan Informasi Layar

### Semua Layar

```go
screens := app.Screen.GetAll()

for _, screen := range screens {
    fmt.Printf("ID: %s\n", screen.ID)
    fmt.Printf("Name: %s\n", screen.Name)
    fmt.Printf("Size: %dx%d\n", screen.Size.Width, screen.Size.Height)
    fmt.Printf("Position: %d,%d\n", screen.X, screen.Y)
    fmt.Printf("Scale: %.2f\n", screen.ScaleFactor)
    fmt.Printf("Primary: %v\n", screen.IsPrimary)
    fmt.Println("---")
}
```

### Layar Utama

```go
primary := app.Screen.GetPrimary()

fmt.Printf("Primary screen: %s\n", primary.Name)
fmt.Printf("Resolution: %dx%d\n", primary.Size.Width, primary.Size.Height)
fmt.Printf("Scale factor: %.2f\n", primary.ScaleFactor)
```

### Layar Saat Ini

Dapatkan layar yang memuat sebuah jendela:

```go
screen, err := window.GetScreen()
if err != nil {
    // handle error
}
fmt.Printf("Window is on: %s\n", screen.Name)
```

### Layar berdasarkan ID

```go
screen := app.Screen.GetByID("screen-id")
if screen != nil {
    fmt.Printf("Found screen: %s\n", screen.Name)
}
```

## Properti Layar

### Struktur Layar

```go
type Screen struct {
    ID               string  // Unique identifier
    Name             string  // Display name
    ScaleFactor      float32 // DPI scale (1.0, 1.5, 2.0, etc.)
    X, Y             int     // Position (top-left, logical pixels)
    Size             Size    // Logical size (Width, Height)
    Bounds           Rect    // Logical bounds (X, Y, Width, Height)
    PhysicalBounds   Rect    // Physical bounds (scaled)
    WorkArea         Rect    // Logical work area excluding taskbars/menubars
    PhysicalWorkArea Rect    // Physical work area
    IsPrimary        bool    // Is this the primary screen?
    Rotation         float32 // Screen rotation in degrees (e.g. 0, 90, 180, 270)
}

// Size and Rect are simple value types:
type Size struct{ Width, Height int }
type Rect struct{ X, Y, Width, Height int }
```

Baca ukuran piksel melalui `screen.Size.Width` / `screen.Size.Height` — tidak ada field `Width`/`Height` di tingkat teratas.

### Piksel Fisik vs. Logis

```go
screen := app.Screen.GetPrimary()

// Logical pixels (what you use)
logicalWidth := screen.Size.Width
logicalHeight := screen.Size.Height

// Physical pixels (actual display)
physicalWidth := int(float32(screen.Size.Width) * screen.ScaleFactor)
physicalHeight := int(float32(screen.Size.Height) * screen.ScaleFactor)

fmt.Printf("Logical: %dx%d\n", logicalWidth, logicalHeight)
fmt.Printf("Physical: %dx%d\n", physicalWidth, physicalHeight)
fmt.Printf("Scale: %.2f\n", screen.ScaleFactor)
```

**Faktor skala yang umum:**

- `1.0` - DPI standar (96 DPI)
- `1.25` - Penskalaan 125% (120 DPI)
- `1.5` - Penskalaan 150% (144 DPI)
- `2.0` - Penskalaan 200% (192 DPI) - Retina
- `3.0` - Penskalaan 300% (288 DPI) - 4K/5K

## Penempatan Jendela

### Tengahkan di Layar

```go
func centreOnScreen(window *application.WebviewWindow, screen *Screen) {
    windowWidth, windowHeight := window.Size()
    
    x := screen.X + (screen.Size.Width-windowWidth)/2
    y := screen.Y + (screen.Size.Height-windowHeight)/2
    
    window.SetPosition(x, y)
}
```

### Tempatkan di Layar Tertentu

```go
func moveToScreen(window *application.WebviewWindow, screenIndex int) {
    screens := app.Screen.GetAll()
    
    if screenIndex < 0 || screenIndex >= len(screens) {
        return
    }
    
    screen := screens[screenIndex]
    
    // Centre on target screen
    centreOnScreen(window, screen)
}
```

### Tempatkan secara Relatif terhadap Layar

```go
// Top-left corner
func positionTopLeft(window *application.WebviewWindow, screen *Screen) {
    window.SetPosition(screen.X+10, screen.Y+10)
}

// Top-right corner
func positionTopRight(window *application.WebviewWindow, screen *Screen) {
    windowWidth, _ := window.Size()
    window.SetPosition(screen.X+screen.Size.Width-windowWidth-10, screen.Y+10)
}

// Bottom-right corner
func positionBottomRight(window *application.WebviewWindow, screen *Screen) {
    windowWidth, windowHeight := window.Size()
    window.SetPosition(
        screen.X+screen.Size.Width-windowWidth-10,
        screen.Y+screen.Size.Height-windowHeight-10,
    )
}
```

## Dukungan Multi-Monitor

### Deteksi Beberapa Monitor

```go
func hasMultipleMonitors() bool {
    return len(app.Screen.GetAll()) > 1
}

func getMonitorCount() int {
    return len(app.Screen.GetAll())
}
```

### Cantumkan Semua Monitor

```go
func listMonitors() {
    screens := app.Screen.GetAll()
    
    fmt.Printf("Found %d monitor(s):\n", len(screens))
    
    for i, screen := range screens {
        primary := ""
        if screen.IsPrimary {
            primary = " (Primary)"
        }
        
        fmt.Printf("%d. %s%s\n", i+1, screen.Name, primary)
        fmt.Printf("   Resolution: %dx%d\n", screen.Size.Width, screen.Size.Height)
        fmt.Printf("   Position: %d,%d\n", screen.X, screen.Y)
        fmt.Printf("   Scale: %.2fx\n", screen.ScaleFactor)
    }
}
```

### Pilih Monitor

```go
func chooseMonitor() (*Screen, error) {
    screens := app.Screen.GetAll()
    
    if len(screens) == 1 {
        return screens[0], nil
    }
    
    // Show dialog to choose
    var options []string
    for i, screen := range screens {
        primary := ""
        if screen.IsPrimary {
            primary = " (Primary)"
        }
        options = append(options, 
            fmt.Sprintf("%d. %s%s - %dx%d", 
                i+1, screen.Name, primary, screen.Size.Width, screen.Size.Height))
    }
    
    // Use dialog to select
    // (Implementation depends on your dialog system)
    
    return screens[0], nil
}
```

## Contoh Lengkap

### Pengelola Jendela Multi-Monitor

```go
type MultiMonitorManager struct {
    app     *application.App
    windows map[int]*application.WebviewWindow
}

func NewMultiMonitorManager(app *application.App) *MultiMonitorManager {
    return &MultiMonitorManager{
        app:     app,
        windows: make(map[int]*application.WebviewWindow),
    }
}

func (m *MultiMonitorManager) CreateWindowOnScreen(screenIndex int) error {
    screens := m.app.Screen.GetAll()
    
    if screenIndex < 0 || screenIndex >= len(screens) {
        return errors.New("invalid screen index")
    }
    
    screen := screens[screenIndex]
    
    // Create window
    window := m.app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:  fmt.Sprintf("Window on %s", screen.Name),
        Width:  800,
        Height: 600,
    })
    
    // Centre on screen
    x := screen.X + (screen.Size.Width-800)/2
    y := screen.Y + (screen.Size.Height-600)/2
    window.SetPosition(x, y)
    
    window.Show()
    
    m.windows[screenIndex] = window
    return nil
}

func (m *MultiMonitorManager) CreateWindowOnEachScreen() {
    screens := m.app.Screen.GetAll()
    
    for i := range screens {
        m.CreateWindowOnScreen(i)
    }
}
```

### Deteksi Perubahan Layar

```go
type ScreenMonitor struct {
    app           *application.App
    lastScreens   []*Screen
    changeHandler func([]*Screen)
}

func NewScreenMonitor(app *application.App) *ScreenMonitor {
    return &ScreenMonitor{
        app:         app,
        lastScreens: app.Screen.GetAll(),
    }
}

func (sm *ScreenMonitor) OnScreenChange(handler func([]*Screen)) {
    sm.changeHandler = handler
}

func (sm *ScreenMonitor) Start() {
    ticker := time.NewTicker(2 * time.Second)
    
    go func() {
        for range ticker.C {
            sm.checkScreens()
        }
    }()
}

func (sm *ScreenMonitor) checkScreens() {
    current := sm.app.Screen.GetAll()
    
    if len(current) != len(sm.lastScreens) {
        sm.lastScreens = current
        if sm.changeHandler != nil {
            sm.changeHandler(current)
        }
    }
}
```

### Ukuran Jendela yang Memperhitungkan DPI

```go
func createDPIAwareWindow(screen *Screen) *application.WebviewWindow {
    // Base size at 1.0 scale
    baseWidth := 800
    baseHeight := 600
    
    // Adjust for DPI
    width := int(float32(baseWidth) * screen.ScaleFactor)
    height := int(float32(baseHeight) * screen.ScaleFactor)
    
    window := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:  "DPI-Aware Window",
        Width:  width,
        Height: height,
    })
    
    // Centre on screen
    x := screen.X + (screen.Size.Width-width)/2
    y := screen.Y + (screen.Size.Height-height)/2
    window.SetPosition(x, y)
    
    return window
}
```

### Visualisator Tata Letak Layar

```go
func visualiseScreenLayout() string {
    screens := app.Screen.GetAll()
    
    var layout strings.Builder
    layout.WriteString("Screen Layout:\n\n")
    
    for i, screen := range screens {
        primary := ""
        if screen.IsPrimary {
            primary = " [PRIMARY]"
        }
        
        layout.WriteString(fmt.Sprintf("Screen %d: %s%s\n", i+1, screen.Name, primary))
        layout.WriteString(fmt.Sprintf("  Position: (%d, %d)\n", screen.X, screen.Y))
        layout.WriteString(fmt.Sprintf("  Size: %dx%d\n", screen.Size.Width, screen.Size.Height))
        layout.WriteString(fmt.Sprintf("  Scale: %.2fx\n", screen.ScaleFactor))
        layout.WriteString(fmt.Sprintf("  Physical: %dx%d\n", 
            int(float32(screen.Size.Width)*screen.ScaleFactor),
            int(float32(screen.Size.Height)*screen.ScaleFactor)))
        layout.WriteString("\n")
    }
    
    return layout.String()
}
```

## Praktik Terbaik

### ✅ Lakukan

- **Periksa jumlah layar** - Tangani satu maupun beberapa monitor
- **Gunakan piksel logis** - Wails menangani DPI secara otomatis
- **Tengahkan jendela** - Memberikan pengalaman pengguna yang lebih baik daripada posisi tetap
- **Validasi posisi** - Pastikan jendela terlihat
- **Tangani perubahan layar** - Monitor dapat ditambahkan atau dilepas
- **Uji pada DPI yang berbeda** - 100%, 125%, 150%, 200%

### ❌ Jangan Lakukan

- **Jangan melakukan hardcode pada posisi** - Gunakan dimensi layar
- **Jangan berasumsi layar utama selalu digunakan** - Pengguna mungkin memiliki beberapa layar
- **Jangan mengabaikan faktor skala** - Penting untuk kesadaran DPI
- **Jangan menempatkan jendela di luar layar** - Validasi koordinat
- **Jangan melupakan perubahan layar** - Laptop dapat dipasang ke atau dilepas dari dok
- **Jangan gunakan piksel fisik** - Gunakan piksel logis

## Perbedaan Antarplatform

### macOS

- Layar Retina (faktor skala 2x)
- Penggunaan beberapa layar lazim
- Sistem koordinat: (0,0) di kiri bawah
- Spaces (desktop virtual) memengaruhi penempatan

### Windows

- Beragam penskalaan DPI (100%, 125%, 150%, 200%)
- Penggunaan beberapa layar lazim
- Sistem koordinat: (0,0) di kiri atas
- Kesadaran DPI per monitor

### Linux

- Bervariasi menurut lingkungan desktop
- Perbedaan antara X11 dan Wayland
- Dukungan penskalaan DPI bervariasi
- Mendukung beberapa layar

## Langkah Berikutnya

@cards{cols="2"}
▣ Jendela
Pelajari pengelolaan jendela.

[Pelajari Selengkapnya →](/features/windows/basics/)

---
⚙ Opsi Jendela
Konfigurasikan tampilan jendela.

[Pelajari Selengkapnya →](/features/windows/options/)

---
◆ Beberapa Jendela
Pola untuk beberapa jendela.

[Pelajari Selengkapnya →](/features/windows/multiple/)

---
🚀 Binding
Panggil fungsi Go dari JavaScript.

[Pelajari Selengkapnya →](/features/bindings/methods/)

@end

---

**Ada pertanyaan?** Tanyakan di [Discord](https://discord.gg/JDdSxwjhGf) atau lihat [contoh layar](https://github.com/wailsapp/wails/tree/master/v3/examples).
