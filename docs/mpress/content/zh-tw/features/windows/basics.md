---
title: "視窗基礎"
description: "在 Wails 中建立及管理應用程式視窗"
slug: "features/windows/basics"
sourcePath: "features/windows/basics.md"
---

## 視窗管理

Wails 提供可跨所有平台運作的<strong>統一視窗管理 API</strong>。您可以建立視窗、控制其行為，並完整掌控多個視窗的建立、外觀、行為及生命週期。

## 快速開始

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

<strong>就這麼簡單！</strong>您現在已有一個跨平台視窗。

## 建立視窗

### 基本視窗

建立視窗最簡單的方式：

```go
window := app.Window.New()
```

**您會得到：**

- 預設大小（800x600）
- 預設標題（應用程式名稱）
- 可供前端使用的 WebView
- 平台原生外觀

### 含選項的視窗

使用自訂設定建立視窗：

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

**常用選項：**

| 選項 | 型別 | 說明 |
| --- | --- | --- |
| `Title` | `string` | 視窗標題 |
| `Width` | `int` | 視窗寬度（像素） |
| `Height` | `int` | 視窗高度（像素） |
| `X` | `int` | X 位置（從左側起算） |
| `Y` | `int` | Y 位置（從頂端起算） |
| `AlwaysOnTop` | `bool` | 讓視窗保持在其他視窗上方 |
| `Frameless` | `bool` | 移除標題列和邊框 |
| `Hidden` | `bool` | 啟動時隱藏 |
| `MinWidth` | `int` | 最小寬度 |
| `MinHeight` | `int` | 最小高度 |
| `MaxWidth` | `int` | 最大寬度 |
| `MaxHeight` | `int` | 最大高度 |

**如需完整清單，請參閱[視窗選項](/features/windows/options/)。**

### 具名視窗

為視窗命名，以便輕鬆擷取：

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

**使用情境：**

- 多個視窗（主視窗、設定、關於）
- 從程式碼的不同部分尋找視窗
- 視窗間通訊

## 控制視窗

### 顯示與隱藏

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

**使用情境：**

- 啟動畫面（先顯示，再隱藏）
- 設定視窗（不需要時隱藏）
- 快顯視窗（依需求顯示）

### 位置與大小

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

**座標系統：**

- （0, 0）是主要螢幕的左上角
- X 軸正向朝右
- Y 正值向下

### 視窗狀態

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

**狀態轉換：**

```
Normal ←→ Minimised
Normal ←→ Maximised
Normal ←→ Fullscreen
```

### 標題與外觀

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

### 關閉視窗

```go
// Close window — dispatches WindowClosing; a RegisterHook can call e.Cancel().
window.Close()
```

v3 中沒有 `window.Destroy()` 方法——請使用 `Close()`，並以 `OnWindowEvent` 監聽（無法取消），或以 `RegisterHook` 掛接處理常式（可呼叫 `e.Cancel()` 讓視窗保持開啟）。

## 尋找視窗

### 依名稱

```go
if window, ok := app.Window.GetByName("settings"); ok {
    window.Show()
}
```

### 依 ID

每個視窗都有唯一的 ID：

```go
id := window.ID()
fmt.Printf("Window ID: %d\n", id)

// Find by ID
if found, ok := app.Window.GetByID(id); ok {
    found.Focus()
}
```

### 目前視窗

取得目前具有焦點的視窗：

```go
current := app.Window.Current()
if current != nil {
    current.SetTitle("Active Window")
}
```

### 所有視窗

取得所有視窗：

```go
windows := app.Window.GetAll()
fmt.Printf("Total windows: %d\n", len(windows))

for _, w := range windows {
    fmt.Printf("Window: %s (ID: %d)\n", w.Name(), w.ID())
}
```

## 視窗生命週期

### 建立

```go
app.Window.OnCreate(func(window application.Window) {
    fmt.Printf("Window created: %s\n", window.Name())

    // Configure new windows
    window.SetMinSize(400, 300)
})
```

### 關閉

若要防止視窗關閉，請搭配 `WindowClosing` 事件使用 `RegisterHook`：

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

**重要：**`RegisterHook` 會在關閉事件發生前攔截該事件。呼叫 `event.Cancel()` 可防止視窗關閉。這適用於使用者主動發起的關閉操作（按一下 X 按鈕）。

### 銷毀

若要在視窗關閉時執行清理，請搭配 `WindowClosing` 事件使用 `OnWindowEvent`：

```go
window.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
    fmt.Println("Window is closing")
    // Cleanup resources
})
```

## 多個視窗

### 建立多個視窗

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

### 視窗間通訊

視窗可透過事件互相通訊：

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

**如需更多資訊，請參閱[事件](/features/events/system/)。**

### 父子視窗

`WebviewWindowOptions`沒有`Parent`欄位。請將子視窗建立為一般視窗，並以表單式模態視窗的形式附加至父視窗：

```go
// Create child window
childWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Child Window",
})

// Attach to the parent — presents as a sheet on macOS.
mainWindow.AttachModal(childWindow)
```

**行為：**

- 子視窗會保持在父視窗上方。
- 子視窗是模態視窗——會阻止與父視窗互動。

**平台支援：**

- <strong>macOS：</strong>完整支援（以表單式對話框顯示）。
- <strong>Windows：</strong>不支援。
- <strong>Linux：</strong>不支援。

## 平台特定功能

@tabs{sync-key="platform"}
[Windows]
**Windows 特定功能：**

```go
// Flash taskbar button
window.Flash(true)  // Start flashing
window.Flash(false) // Stop flashing

// Trigger Windows 11 Snap Assist (Win+Z)
window.SnapAssist()
```

沒有個別視窗專用的 `SetIcon`——應透過 `app.SetIcon([]byte)` 設定應用程式圖示（若要設定 Linux 特定的視窗圖示，則在建立視窗時使用 `application.LinuxWindow.Icon` 欄位）。

**貼齊小幫手：** 透過系統捷徑路徑顯示 Windows 11 貼齊配置選項。若自訂 HTML 最大化按鈕需要原生的懸停貼齊配置，請改用 [Windows 上的原生非用戶端區域](/features/windows/frameless/#native-non-client-regions-on-windows)。

**工作列閃爍：** 適合在視窗最小化時發出通知。

[macOS]
**macOS 特定功能：**

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

**背景類型：**

- `MacBackdropNormal` - 標準視窗
- `MacBackdropTranslucent` - 半透明背景；若要讓 WebView 透明，**需要使用私有 API**。
- `MacBackdropTransparent` - 完全透明；若要讓 WebView 透明，**需要使用私有 API**。
- `MacBackdropLiquidGlass` - 玻璃背景；若要讓 WebView 透明，**需要使用私有 API**。

建置時使用 `-tags private_mac_apis`，才能讓這些效果透過 WebView 顯示。若未使用，WebView 會保持不透明。`TitleBar.AppearsTransparent` 本身使用公開 API。請參閱[私有 macOS API](/guides/build/private-macos-apis/)。

**集合行為：** 控制視窗在各個 Space 間的行為：

- `MacWindowCollectionBehaviorCanJoinAllSpaces` - 在所有 Space 上皆可見
- `MacWindowCollectionBehaviorFullScreenAuxiliary` - 可覆蓋全螢幕應用程式

**原生全螢幕：** macOS 全螢幕模式會建立新的 Space（虛擬桌面）。

[Linux]
**Linux 特定功能：**

```go
// Set window icon (per-window struct is LinuxWindow, not the app-level LinuxOptions)
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Linux: application.LinuxWindow{
        Icon: iconBytes,
    },
})
```

**桌面環境注意事項：**

- GNOME：完整支援
- KDE Plasma：完整支援
- XFCE：部分支援
- 其他：視情況而定

**平鋪式視窗管理員（Hyprland、Sway、i3 等）：**

- `Minimise()` 和 `Maximise()` 可能不會如預期運作——視窗幾何配置由視窗管理員控制
- `SetSize()`和`SetPosition()`請求僅供參考，可能會被忽略
- `Fullscreen()`通常會如預期運作
- 部分視窗管理員不支援視窗永遠置頂

@end

## 常見模式

### 啟動畫面

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

### 設定視窗

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

### 關閉前確認

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

## 最佳實務

### ✅ 建議做法

- **為重要視窗命名**——日後更容易尋找
- **設定最小尺寸**——避免版面無法使用
- **將視窗置中**——使用者體驗比隨機位置更好
- **處理關閉事件**——避免資料遺失
- **在所有平台上測試**——行為因平台而異
- **使用適當的尺寸**——考量不同的螢幕尺寸

### ❌ 請勿這樣做

- **不要建立過多視窗**——會讓使用者感到困惑
- **不要忘記關閉視窗**——否則會造成記憶體洩漏
- **不要將位置寫死**——螢幕尺寸各不相同
- **不要忽略平台差異**——請徹底測試
- **不要阻塞 UI 執行緒**——長時間執行的作業請使用 goroutine

## 疑難排解

### 視窗未顯示

**可能原因：**

1. 建立視窗時將其設為隱藏
2. 視窗位於螢幕範圍之外
3. 視窗位於其他視窗後方

**解決方法：**

```go
window.Show()
window.Center()
window.Focus()
```

### 視窗尺寸不正確

<strong>原因：</strong>Windows/Linux 上的 DPI 縮放

**解決方法：**

```go
// Wails handles DPI automatically
// Just use logical pixels
window.SetSize(800, 600)
```

### 視窗立即關閉

<strong>原因：</strong>最後一個視窗關閉時，應用程式會結束

**解決方法：**

```go
app := application.New(application.Options{
    Mac: application.MacOptions{
        ApplicationShouldTerminateAfterLastWindowClosed: false,
    },
})
```

## 後續步驟

@cards{cols="2"}
⚙ 視窗選項
所有視窗選項的完整參考資料。

[深入瞭解 →](/features/windows/options/)

---
▣ 多視窗
多視窗應用程式的常見模式。

[深入瞭解 →](/features/windows/multiple/)

---
★ 無邊框視窗
建立自訂視窗框架。

[深入瞭解 →](/features/windows/frameless/)

---
🚀 視窗事件
處理視窗生命週期事件。

[深入瞭解 →](/features/windows/events/)

@end

---

<strong>有問題嗎？</strong>請到[Discord](https://discord.gg/JDdSxwjhGf)提問，或查看[視窗範例](https://github.com/wailsapp/wails/tree/master/v3/examples)。
