---
title: "螢幕資訊"
description: "取得顯示器與監視器的相關資訊"
slug: "features/screens/info"
sourcePath: "features/screens/info.md"
---

## 螢幕資訊

Wails 提供可跨所有平台運作的<strong>統一螢幕 API</strong>。您可以使用一致的程式碼取得螢幕資訊、偵測多台監視器、查詢螢幕屬性（大小、位置、DPI）、識別主要顯示器，以及處理 DPI 縮放。

## 快速開始

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

<strong>就是這麼簡單！</strong>跨平台螢幕資訊。

## 取得螢幕資訊

### 所有螢幕

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

### 主要螢幕

```go
primary := app.Screen.GetPrimary()

fmt.Printf("Primary screen: %s\n", primary.Name)
fmt.Printf("Resolution: %dx%d\n", primary.Size.Width, primary.Size.Height)
fmt.Printf("Scale factor: %.2f\n", primary.ScaleFactor)
```

### 目前螢幕

取得視窗所在的螢幕：

```go
screen, err := window.GetScreen()
if err != nil {
    // handle error
}
fmt.Printf("Window is on: %s\n", screen.Name)
```

### 依 ID 取得螢幕

```go
screen := app.Screen.GetByID("screen-id")
if screen != nil {
    fmt.Printf("Found screen: %s\n", screen.Name)
}
```

## 螢幕屬性

### 螢幕結構

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

透過`screen.Size.Width` / `screen.Size.Height`讀取像素尺寸；頂層沒有`Width`/`Height`欄位。

### 實體像素與邏輯像素

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

**常見縮放比例：**

- `1.0`－標準 DPI（96 DPI）
- `1.25`－125% 縮放（120 DPI）
- `1.5`－150% 縮放（144 DPI）
- `2.0`－200% 縮放（192 DPI）－Retina
- `3.0`－300% 縮放（288 DPI）－4K/5K

## 視窗定位

### 在螢幕上置中

```go
func centreOnScreen(window *application.WebviewWindow, screen *Screen) {
    windowWidth, windowHeight := window.Size()
    
    x := screen.X + (screen.Size.Width-windowWidth)/2
    y := screen.Y + (screen.Size.Height-windowHeight)/2
    
    window.SetPosition(x, y)
}
```

### 在指定螢幕上定位

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

### 相對於螢幕定位

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

## 多螢幕支援

### 偵測多台監視器

```go
func hasMultipleMonitors() bool {
    return len(app.Screen.GetAll()) > 1
}

func getMonitorCount() int {
    return len(app.Screen.GetAll())
}
```

### 列出所有監視器

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

### 選擇監視器

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

## 完整範例

### 多螢幕視窗管理器

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

### 偵測螢幕變更

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

### 可感知 DPI 的視窗大小調整

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

### 螢幕配置視覺化工具

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

## 最佳實務

### ✅ 建議做法

- **檢查螢幕數量**－處理單一與多台監視器的情況
- **使用邏輯像素**－Wails 會自動處理 DPI
- **將視窗置中**－比固定位置提供更好的使用者體驗
- **驗證位置**－確保視窗可見
- **處理螢幕變更**－監視器可能會新增或移除
- **在不同 DPI 下測試**－100%、125%、150%、200%

### ❌ 避免做法

- **不要將位置寫死**－請使用螢幕尺寸
- **不要假設使用主要螢幕**－使用者可能有多台監視器
- **不要忽略縮放比例**－這對 DPI 感知能力很重要
- **不要將視窗定位在螢幕範圍外**－請驗證座標
- **不要忘記處理螢幕變更**－筆記型電腦可能會接上或拔除擴充基座
- **不要使用實體像素**－請使用邏輯像素

## 平台差異

### macOS

- Retina 顯示器（2 倍縮放比例）
- 常見多顯示器配置
- 座標系統：原點（0,0）位於左下角
- Spaces（虛擬桌面）會影響定位

### Windows

- 多種 DPI 縮放比例（100%、125%、150%、200%）
- 常見多顯示器配置
- 座標系統：原點（0,0）位於左上角
- 個別監視器 DPI 感知能力

### Linux

- 依桌面環境而異
- X11 與 Wayland 之間有所差異
- DPI 縮放支援程度不一
- 支援多台顯示器

## 後續步驟

@cards{cols="2"}
▣ 視窗
瞭解視窗管理。

[深入瞭解 →](/features/windows/basics/)

---
⚙ 視窗選項
設定視窗外觀。

[深入瞭解 →](/features/windows/options/)

---
◆ 多個視窗
多視窗模式。

[深入瞭解 →](/features/windows/multiple/)

---
🚀 繫結
從 JavaScript 呼叫 Go 函式。

[深入瞭解 →](/features/bindings/methods/)

@end

---

<strong>有問題嗎？</strong>請在[Discord](https://discord.gg/JDdSxwjhGf)中提問，或查看[螢幕範例](https://github.com/wailsapp/wails/tree/master/v3/examples)。
