---
title: "系統匣選單"
description: "為應用程式新增系統匣（通知區域）整合功能"
slug: "features/menus/systray"
sourcePath: "features/menus/systray.md"
---

## 系統匣選單

Wails 提供可在所有平台運作的<strong>統一系統匣 API</strong>。您可以建立含選單的系統匣圖示、連結視窗及處理點擊，並以平台原生行為支援背景應用程式、服務和快速存取工具。

![從 macOS 選單列開啟的 Wails 系統匣選單](/assets/screenshots/systray-menu-macos.png)

在 macOS 上，Wails 系統匣項目會顯示於選單列中，並開啟原生選單。此範例包含停用項目、核取方塊項目、單選項目、子選單及動作項目。

## 快速入門

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

<strong>結果：</strong>在所有平台上顯示含選單的系統匣圖示。

## 建立系統匣

### 基本系統匣

```go
// Create system tray
systray := app.SystemTray.New()

// Set icon
systray.SetIcon(iconBytes)

// Set label (macOS) / tooltip (Windows)
systray.SetLabel("My Application")
```

### 加入圖示

圖示應嵌入應用程式中：

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

**圖示需求：**

| 平台 | 尺寸 | 格式 | 備註 |
| --- | --- | --- | --- |
| **Windows** | 16x16 或 32x32 | PNG、ICO | 通知區域 |
| **macOS** | 18x18 至 22x22 | PNG | 選單列，建議使用範本圖示 |
| **Linux** | 22x22 至 48x48 | PNG、SVG | 依桌面環境而異 |

### 範本圖示（macOS）

範本圖示會自動配合淺色／深色模式：

```go
systray.SetTemplateIcon(iconBytes)
```

**範本圖示準則：**

- 僅使用黑色及無色（透明）
- 在深色模式下，黑色會變成白色
- 檔名應使用`Template`後綴：`iconTemplate.png`
- [設計指南](https://bjango.com/articles/designingmenubarextras/)

## 新增選單

系統匣選單的運作方式與應用程式選單相同：

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

如需瞭解<strong>所有選單項目類型</strong>，請參閱[選單參考](/features/menus/reference/)。

## 連結視窗

將視窗連結至系統匣圖示，以便自動顯示／隱藏：

```go
// Create window
window := app.Window.New()

// Attach to tray
systray.AttachWindow(window)

// Configure behaviour — these are setters that return the receiver for chaining.
systray.WindowOffset(10)                          // Pixels from tray icon
systray.WindowDebounce(200 * time.Millisecond)    // Click debounce
```

**行為：**

- 視窗啟動時為隱藏狀態
- **以滑鼠左鍵點擊系統匣圖示** → 切換視窗顯示狀態
- **以滑鼠右鍵點擊系統匣圖示** → 顯示選單（若已設定）
- 視窗會置於系統匣圖示附近

**範例：快顯視窗**

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

`HideOnFocusLost`適合用於 Windows、macOS 及採用點擊聚焦模式的 Linux 桌面上的系統匣快顯視窗。Wails 會在採用焦點跟隨滑鼠模式的 Linux 環境（包括常見的 Hyprland、Sway 和 i3 設定）中停用該行為，否則游標離開快顯視窗時，視窗可能會在使用者能夠操作前便隱藏。在這些環境中，`HideOnEscape`仍可使用。

上述滑鼠左鍵和右鍵點擊行為是智慧型預設值。明確指定的`OnClick`或`OnRightClick`處理常式會取代對應的預設行為。如需瞭解平台檢查及邊界案例，請參閱[手動系統匣測試套件](https://github.com/wailsapp/wails/tree/master/v3/test/manual/systray)和[系統匣壓力測試範例](https://github.com/wailsapp/wails/tree/master/v3/examples/systray-stress)。

## 點擊處理常式

處理系統匣圖示的點擊事件：

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

**平台支援：**

| 事件 | Windows | macOS | Linux |
| --- | --- | --- | --- |
| OnClick | ✅ | ✅ | ✅ |
| OnRightClick | ✅ | ✅ | ✅ |
| OnDoubleClick | ✅ | ✅ | ⚠️ 視情況而異 |
| OnMouseEnter | ✅ | ✅ | ⚠️ 視情況而異 |
| OnMouseLeave | ✅ | ✅ | ⚠️ 視情況而異 |

## 動態更新

動態更新系統匣圖示和選單：

### 變更圖示

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

### 更新選單

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

@note{type="caution" title="一律呼叫 Update()"}
變更選單狀態後，**請呼叫`menu.Update()`**。請參閱[選單參考](/features/menus/reference/#enabled-state)。

@end

### 重建選單

若要進行重大變更，請重建整個選單：

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

## 平台特有功能

@tabs{sync-key="platform"}
[macOS]
**選單列整合：**

```go
// Set label (appears next to icon)
systray.SetLabel("My App")

// Use template icon (adapts to dark mode)
systray.SetTemplateIcon(iconBytes)

// Set icon position — uses AppKit NSImage placement constants.
systray.SetIconPosition(application.NSImageRight)
```

**圖示位置**（對應`NSImagePosition`）：

- `application.NSImageLeft`－圖示位於標籤左側。
- `application.NSImageRight`－圖示位於標籤右側。
- `application.NSImageOnly`－僅顯示圖示，不顯示標籤。
- `application.NSImageNone`－僅顯示標籤，不顯示圖示。

**最佳做法：**

- 使用樣板圖示（黑色＋透明）
- 標籤應保持簡短（3-5個字元）
- Retina 顯示器使用 18x18 至 22x22 像素
- 在淺色和深色模式下都進行測試

[Windows]
**通知區域整合：**

```go
// Set tooltip (appears on hover)
systray.SetTooltip("My Application")

// Or use SetLabel (same as tooltip on Windows)
systray.SetLabel("My Application")

// Show/Hide functionality (fully functional)
systray.Show()  // Show tray icon
systray.Hide()  // Hide tray icon
```

**圖示需求：**

- 16x16 或 32x32 像素
- PNG 或 ICO 格式
- 透明背景

**工具提示限制：**

- 最多 127 個 UTF-16 字元
- 較長的工具提示會遭到截斷
- 保持簡潔，以提供最佳使用體驗

**平台功能：**

- Windows 檔案總管重新啟動後，系統匣圖示仍會保留
- Show() 和 Hide() 方法功能完整
- 完善的生命週期管理

**最佳做法：**

- 高 DPI 顯示器使用 32x32 像素
- 工具提示應少於 127 個字元
- 在不同 Windows 版本上進行測試
- 將通知區域的溢位情況納入考量
- 使用 Show/Hide 依條件控制系統匣的可見性

[Linux]
**系統匣整合：**

使用 StatusNotifierItem 規範（大多數現代桌面環境）。

```go
systray.SetIcon(iconBytes)
systray.SetLabel("My App")
```

**桌面環境支援：**

- **GNOME**：頂端列（需使用擴充套件）
- **KDE Plasma**：系統匣
- **XFCE**：通知區域
- **其他**：視情況而異

**最佳做法：**

- 使用 22x22 或 24x24 像素
- SVG 圖示的縮放效果較佳
- 在目標桌面環境上進行測試
- 為不受支援的桌面環境提供備用方案

@end

## 完整範例

以下是可用於正式環境的系統匣應用程式：

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

## 可見性控制

動態顯示或隱藏系統匣圖示：

```go
// Hide tray icon
systray.Hide()

// Show tray icon
systray.Show()
```

沒有 `IsVisible()` getter；如有需要，請在應用程式自己的狀態中追蹤可見性。

**平台支援：**

| 平台 | Hide() | Show() | 備註 |
| --- | --- | --- | --- |
| **Windows** | ✅ | ✅ | 功能完整－圖示會在通知區域中出現或消失 |
| **macOS** | ✅ | ✅ | 顯示或隱藏選單列項目 |
| **Linux** | ✅ | ✅ | 視桌面環境而異 |

**使用情境：**

- 根據使用者偏好暫時隱藏系統匣圖示
- 在無介面模式下，僅於需要時顯示系統匣圖示
- 根據應用程式狀態切換可見性

**範例－依條件控制系統匣可見性：**

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

## 清理

完成後銷毀系統匣圖示：

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

<strong>重要：</strong>關閉應用程式時，務必銷毀系統匣以釋放資源。

## 最佳做法

### ✅ 建議做法

- **在 macOS 上使用範本圖示** — 可配合深色模式調整
- **保持標籤簡短** — 最多 3-5 個字元
- **在 Windows 上提供工具提示** — 協助使用者識別您的應用程式
- **在所有平台上測試** — 行為因平台而異
- **適當處理點擊事件** — 按一下滑鼠左鍵執行主要動作，按一下滑鼠右鍵開啟選單
- **更新圖示以反映狀態** — 視覺回饋很重要
- **關閉時銷毀系統匣** — 釋放資源

### ❌ 不要這樣做

- **不要使用大型圖示** — 請遵循平台準則
- **不要使用過長的標籤** — 內容會遭到截斷
- **不要忽略深色模式** — 請在 Windows 和 macOS 的深色模式下測試
- **不要阻塞點擊事件處理函式** — 應快速完成處理
- **不要忘記呼叫 menu.Update()** — 請在變更選單狀態後呼叫
- **不要假設系統支援系統匣** — 部分 Linux 桌面環境不支援

## 疑難排解

### 系統匣圖示未顯示

**可能的原因：**

1. 不支援該圖示格式
2. 圖示尺寸過大或過小
3. 不支援系統匣（Linux）

**解決方法：**

沒有 `SystemTraySupported()` 輔助函式；請改為建立系統匣、檢查平台，並妥善採用降級處理：

```go
// Probe support: on Linux without a notification-area extension, the tray
// will simply not appear. Defensive code can fall back to window-only mode
// based on runtime.GOOS or after a short timeout if no tray events arrive.
systray := app.SystemTray.New()
systray.SetIcon(iconBytes)
```

### 圖示在 macOS 上顯示異常

<strong>原因：</strong>未使用範本圖示

**解決方法：**

```go
// Use template icon
systray.SetTemplateIcon(iconBytes)

// Or design icon as template (black + transparent)
```

### 選單未更新

<strong>原因：</strong>忘記呼叫 `menu.Update()`

**解決方法：**

```go
menuItem.SetLabel("New Label")
menu.Update()  // Add this!
```

## 後續步驟

@cards{cols="2"}
📖 選單參考資料
選單項目類型與屬性的完整參考資料。

[深入瞭解 →](/features/menus/reference/)

---
☰ 應用程式選單
建立應用程式選單列。

[深入瞭解 →](/features/menus/application/)

---
◆ 快顯功能表
建立按一下滑鼠右鍵時顯示的快顯功能表。

[深入瞭解 →](/features/menus/context/)

---
📖 系統匣範例
探索完整的系統匣應用程式。

[深入瞭解 →](https://github.com/wailsapp/wails/tree/master/v3/examples/systray-basic)

@end

---

<strong>有問題嗎？</strong>請在 [Discord](https://discord.gg/JDdSxwjhGf) 中提問，或查看[系統匣範例](https://github.com/wailsapp/wails/tree/master/v3/examples/systray-basic)。
