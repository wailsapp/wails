---
title: "系统托盘菜单"
description: "为应用程序添加系统托盘（通知区域）集成"
slug: "features/menus/systray"
sourcePath: "features/menus/systray.md"
---

## 系统托盘菜单

Wails 提供可在所有平台上使用的<strong>统一系统托盘 API</strong>。您可以创建带菜单的托盘图标、附加窗口并处理点击操作，从而让后台应用程序、服务和快速访问实用工具具备平台原生行为。

![从 macOS 菜单栏打开的 Wails 系统托盘菜单](/assets/screenshots/systray-menu-macos.png)

在 macOS 上，Wails 系统托盘项会显示在菜单栏中，并打开原生菜单。该示例包含禁用项、复选项、单选项、子菜单和操作项。

## 快速入门

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

<strong>结果：</strong>在所有平台上显示带菜单的系统托盘图标。

## 创建系统托盘

### 基本系统托盘

```go
// Create system tray
systray := app.SystemTray.New()

// Set icon
systray.SetIcon(iconBytes)

// Set label (macOS) / tooltip (Windows)
systray.SetLabel("My Application")
```

### 添加图标

图标应嵌入应用程序：

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

**图标要求：**

| 平台 | 尺寸 | 格式 | 备注 |
| --- | --- | --- | --- |
| **Windows** | 16x16 或 32x32 | PNG、ICO | 通知区域 |
| **macOS** | 18x18 至 22x22 | PNG | 菜单栏，建议使用模板图标 |
| **Linux** | 22x22 至 48x48 | PNG、SVG | 因桌面环境而异 |

### 模板图标（macOS）

模板图标会自动适应浅色或深色模式：

```go
systray.SetTemplateIcon(iconBytes)
```

**模板图标指南：**

- 仅使用黑色和无色（透明）
- 在深色模式下，黑色会变为白色
- 为文件名添加`Template`后缀：`iconTemplate.png`
- [设计指南](https://bjango.com/articles/designingmenubarextras/)

## 添加菜单

系统托盘菜单的用法与应用程序菜单相同：

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

有关<strong>所有菜单项类型</strong>，请参阅[菜单参考](/features/menus/reference/)。

## 附加窗口

将窗口附加到托盘图标，以便自动显示或隐藏窗口：

```go
// Create window
window := app.Window.New()

// Attach to tray
systray.AttachWindow(window)

// Configure behaviour — these are setters that return the receiver for chaining.
systray.WindowOffset(10)                          // Pixels from tray icon
systray.WindowDebounce(200 * time.Millisecond)    // Click debounce
```

**行为：**

- 窗口启动时处于隐藏状态
- **左键单击托盘图标** → 切换窗口可见性
- **右键单击托盘图标** → 显示菜单（如果已设置）
- 窗口位于托盘图标附近

**示例：弹出窗口**

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

`HideOnFocusLost`适用于 Windows、macOS 和采用点击聚焦方式的 Linux 桌面上的托盘弹出窗口。Wails 会在采用焦点跟随鼠标方式的 Linux 环境（包括常见的 Hyprland、Sway 和 i3 配置）中禁用该行为，否则鼠标离开弹出窗口时，窗口可能会在用户能够操作之前隐藏。`HideOnEscape`在这些环境中仍然可用。

上述左键和右键单击行为是智能默认设置。显式设置`OnClick`或`OnRightClick`处理程序会替代相应的默认行为。有关平台检查和边界情况，请参阅[手动系统托盘测试套件](https://github.com/wailsapp/wails/tree/master/v3/test/manual/systray)和[系统托盘压力测试示例](https://github.com/wailsapp/wails/tree/master/v3/examples/systray-stress)。

## 点击处理程序

处理托盘图标点击操作：

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

**平台支持：**

| 事件 | Windows | macOS | Linux |
| --- | --- | --- | --- |
| OnClick | ✅ | ✅ | ✅ |
| OnRightClick | ✅ | ✅ | ✅ |
| OnDoubleClick | ✅ | ✅ | ⚠️ 因环境而异 |
| OnMouseEnter | ✅ | ✅ | ⚠️ 因环境而异 |
| OnMouseLeave | ✅ | ✅ | ⚠️ 因环境而异 |

## 动态更新

动态更新托盘图标和菜单：

### 更改图标

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

### 更新菜单

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

@note{type="caution" title="始终调用 Update()"}
更改菜单状态后，**调用`menu.Update()`**。请参阅[菜单参考](/features/menus/reference/#enabled-state)。

@end

### 重建菜单

对于重大更改，请重新构建整个菜单：

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

## 平台特定功能

@tabs{sync-key="platform"}
[macOS]
**菜单栏集成：**

```go
// Set label (appears next to icon)
systray.SetLabel("My App")

// Use template icon (adapts to dark mode)
systray.SetTemplateIcon(iconBytes)

// Set icon position — uses AppKit NSImage placement constants.
systray.SetIconPosition(application.NSImageRight)
```

**图标位置**（对应`NSImagePosition`）：

- `application.NSImageLeft` - 图标位于标签左侧。
- `application.NSImageRight` - 图标位于标签右侧。
- `application.NSImageOnly` - 仅显示图标，不显示标签。
- `application.NSImageNone` - 仅显示标签，不显示图标。

**最佳实践：**

- 使用模板图标（黑色 + 透明）
- 保持标签简短（3-5个字符）
- Retina 显示屏使用 18x18 至 22x22 像素
- 同时在浅色和深色模式下测试

[Windows]
**通知区域集成：**

```go
// Set tooltip (appears on hover)
systray.SetTooltip("My Application")

// Or use SetLabel (same as tooltip on Windows)
systray.SetLabel("My Application")

// Show/Hide functionality (fully functional)
systray.Show()  // Show tray icon
systray.Hide()  // Hide tray icon
```

**图标要求：**

- 16x16 或 32x32 像素
- PNG 或 ICO 格式
- 透明背景

**工具提示限制：**

- 最多 127 个 UTF-16 字符
- 较长的工具提示将被截断
- 保持简洁以获得最佳体验

**平台功能：**

- Windows 资源管理器重启后，托盘图标仍会保留
- Show() 和 Hide() 方法功能完整
- 妥善的生命周期管理

**最佳实践：**

- 高 DPI 显示屏使用 32x32 像素
- 工具提示应少于 127 个字符
- 在不同 Windows 版本上测试
- 考虑通知区域溢出区
- 使用 Show/Hide 按条件控制托盘图标的可见性

[Linux]
**系统托盘集成：**

使用 StatusNotifierItem 规范（大多数现代桌面环境）。

```go
systray.SetIcon(iconBytes)
systray.SetLabel("My App")
```

**桌面环境支持：**

- **GNOME**：顶部栏（需要扩展）
- **KDE Plasma**：系统托盘
- **XFCE**：通知区域
- **其他**：视具体情况而定

**最佳实践：**

- 使用 22x22 或 24x24 像素
- SVG 图标的缩放效果更好
- 在目标桌面环境上测试
- 为不受支持的桌面环境提供回退方案

@end

## 完整示例

下面是一个可用于生产环境的系统托盘应用程序：

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

## 可见性控制

动态显示或隐藏托盘图标：

```go
// Hide tray icon
systray.Hide()

// Show tray icon
systray.Show()
```

没有`IsVisible()`获取方法；如有需要，请在应用程序自身的状态中跟踪可见性。

**平台支持：**

| 平台 | Hide() | Show() | 说明 |
| --- | --- | --- | --- |
| **Windows** | ✅ | ✅ | 功能完整——图标会在通知区域中显示或消失 |
| **macOS** | ✅ | ✅ | 显示或隐藏菜单栏项目 |
| **Linux** | ✅ | ✅ | 因桌面环境而异 |

**使用场景：**

- 根据用户偏好暂时隐藏托盘图标
- 无界面模式，仅在需要时显示托盘图标
- 根据应用程序状态切换可见性

**示例——按条件控制托盘图标的可见性：**

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

使用完毕后销毁托盘图标：

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

<strong>重要：</strong>关闭应用程序时，务必销毁系统托盘以释放资源。

## 最佳实践

### ✅ 应该做

- **在 macOS 上使用模板图标**——可适配深色模式
- **保持标签简短**——最多 3-5 个字符
- **在 Windows 上提供工具提示**——帮助用户识别你的应用
- **在所有平台上测试**——行为因平台而异
- **正确处理点击事件**——左键单击执行主要操作，右键单击打开菜单
- **根据状态更新图标**——视觉反馈很重要
- **关闭时销毁系统托盘**——释放资源

### ❌ 不要这样做

- **不要使用过大的图标**——遵循平台规范
- **不要使用过长的标签**——标签会被截断
- **不要忽略深色模式**——在 Windows 和 macOS 的深色模式下进行测试
- **不要阻塞点击处理程序**——确保处理程序快速执行
- **不要忘记调用 menu.Update()**——应在更改菜单状态后调用
- **不要假定系统支持托盘**——某些 Linux 桌面环境不支持系统托盘

## 故障排除

### 托盘图标未显示

**可能的原因：**

1. 不支持该图标格式
2. 图标尺寸过大或过小
3. 系统不支持托盘（Linux）

**解决方案：**

没有 `SystemTraySupported()` 辅助函数；请改为创建托盘、检查平台，并在不受支持时进行优雅降级：

```go
// Probe support: on Linux without a notification-area extension, the tray
// will simply not appear. Defensive code can fall back to window-only mode
// based on runtime.GOOS or after a short timeout if no tray events arrive.
systray := app.SystemTray.New()
systray.SetIcon(iconBytes)
```

### 图标在 macOS 上显示异常

<strong>原因：</strong>未使用模板图标

**解决方案：**

```go
// Use template icon
systray.SetTemplateIcon(iconBytes)

// Or design icon as template (black + transparent)
```

### 菜单未更新

<strong>原因：</strong>忘记调用 `menu.Update()`

**解决方案：**

```go
menuItem.SetLabel("New Label")
menu.Update()  // Add this!
```

## 后续步骤

@cards{cols="2"}
📖 菜单参考
菜单项类型和属性的完整参考。

[了解更多 →](/features/menus/reference/)

---
☰ 应用菜单
创建应用菜单栏。

[了解更多 →](/features/menus/application/)

---
◆ 上下文菜单
创建右键上下文菜单。

[了解更多 →](/features/menus/context/)

---
📖 系统托盘示例
探索一个完整的系统托盘应用。

[了解更多 →](https://github.com/wailsapp/wails/tree/master/v3/examples/systray-basic)

@end

---

<strong>有疑问？</strong>请在 [Discord](https://discord.gg/JDdSxwjhGf) 中提问，或查看[系统托盘示例](https://github.com/wailsapp/wails/tree/master/v3/examples/systray-basic)。
