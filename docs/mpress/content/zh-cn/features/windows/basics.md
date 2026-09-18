---
title: "窗口基础"
description: "在 Wails 中创建和管理应用程序窗口"
slug: "features/windows/basics"
sourcePath: "features/windows/basics.md"
---

## 窗口管理

Wails 提供了一套可跨所有平台使用的<strong>统一窗口管理 API</strong>。您可以创建窗口、控制窗口行为并管理多个窗口，全面掌控窗口的创建、外观、行为和生命周期。

## 快速入门

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

<strong>就这么简单！</strong>您现在已经拥有了一个跨平台窗口。

## 创建窗口

### 基本窗口

创建窗口最简单的方法：

```go
window := app.Window.New()
```

**您将获得：**

- 默认尺寸（800x600）
- 默认标题（应用程序名称）
- 可供前端使用的 WebView
- 平台原生外观

### 带选项的窗口

使用自定义配置创建窗口：

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

**常用选项：**

| 选项 | 类型 | 说明 |
| --- | --- | --- |
| `Title` | `string` | 窗口标题 |
| `Width` | `int` | 窗口宽度（像素） |
| `Height` | `int` | 窗口高度（像素） |
| `X` | `int` | X 位置（从左侧起） |
| `Y` | `int` | Y 位置（从顶部起） |
| `AlwaysOnTop` | `bool` | 使窗口保持在其他窗口上方 |
| `Frameless` | `bool` | 移除标题栏和边框 |
| `Hidden` | `bool` | 启动时隐藏 |
| `MinWidth` | `int` | 最小宽度 |
| `MinHeight` | `int` | 最小高度 |
| `MaxWidth` | `int` | 最大宽度 |
| `MaxHeight` | `int` | 最大高度 |

**完整列表请参阅[窗口选项](/features/windows/options/)。**

### 命名窗口

为窗口命名，以便轻松检索：

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

**使用场景：**

- 多个窗口（主窗口、设置窗口、关于窗口）
- 从代码的不同位置查找窗口
- 窗口通信

## 控制窗口

### 显示和隐藏

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

**使用场景：**

- 启动画面（先显示，然后隐藏）
- 设置窗口（不需要时隐藏）
- 弹出窗口（按需显示）

### 位置和尺寸

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

**坐标系：**

- (0, 0) 位于主屏幕左上角
- X 轴正方向向右
- Y 轴正方向向下

### 窗口状态

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

**状态转换：**

```
Normal ←→ Minimised
Normal ←→ Maximised
Normal ←→ Fullscreen
```

### 标题和外观

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

### 关闭窗口

```go
// Close window — dispatches WindowClosing; a RegisterHook can call e.Cancel().
window.Close()
```

v3 中没有`window.Destroy()`方法——请使用`Close()`，并通过`OnWindowEvent`监听（无法取消），或通过`RegisterHook`挂接处理程序（可调用`e.Cancel()`让窗口保持打开）。

## 查找窗口

### 按名称

```go
if window, ok := app.Window.GetByName("settings"); ok {
    window.Show()
}
```

### 按 ID

每个窗口都有唯一 ID：

```go
id := window.ID()
fmt.Printf("Window ID: %d\n", id)

// Find by ID
if found, ok := app.Window.GetByID(id); ok {
    found.Focus()
}
```

### 当前窗口

获取当前获得焦点的窗口：

```go
current := app.Window.Current()
if current != nil {
    current.SetTitle("Active Window")
}
```

### 所有窗口

获取所有窗口：

```go
windows := app.Window.GetAll()
fmt.Printf("Total windows: %d\n", len(windows))

for _, w := range windows {
    fmt.Printf("Window: %s (ID: %d)\n", w.Name(), w.ID())
}
```

## 窗口生命周期

### 创建

```go
app.Window.OnCreate(func(window application.Window) {
    fmt.Printf("Window created: %s\n", window.Name())

    // Configure new windows
    window.SetMinSize(400, 300)
})
```

### 关闭

要阻止窗口关闭，请对`WindowClosing`事件使用`RegisterHook`：

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

**重要：**`RegisterHook`会在关闭事件发生前将其拦截。调用`event.Cancel()`可阻止窗口关闭。此方法适用于用户发起的关闭操作（单击 X 按钮）。

### 销毁

要在窗口关闭时执行清理，请对`WindowClosing`事件使用`OnWindowEvent`：

```go
window.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
    fmt.Println("Window is closing")
    // Cleanup resources
})
```

## 多个窗口

### 创建多个窗口

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

### 窗口通信

窗口可通过事件进行通信：

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

**有关更多信息，请参阅[事件](/features/events/system/)。**

### 父子窗口

`WebviewWindowOptions`没有`Parent`字段。请将子窗口创建为普通窗口，然后以表单式模态窗口的形式将其附加到父窗口：

```go
// Create child window
childWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Child Window",
})

// Attach to the parent — presents as a sheet on macOS.
mainWindow.AttachModal(childWindow)
```

**行为：**

- 子窗口保持在父窗口上方。
- 子窗口是模态窗口——会阻止与父窗口交互。

**平台支持：**

- <strong>macOS：</strong>完全支持（显示为附着于父窗口的对话框）。
- <strong>Windows：</strong>不支持。
- <strong>Linux：</strong>不支持。

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

没有针对单个窗口的`SetIcon`——应用程序图标通过`app.SetIcon([]byte)`在应用级别设置（对于 Linux 特定的窗口图标，也可在创建窗口时使用`application.LinuxWindow.Icon`字段）。

**贴靠助手：** 通过系统快捷方式路径显示 Windows 11贴靠布局选项。对于需要在原生悬停时显示贴靠布局的自定义 HTML 最大化按钮，请改用[Windows 上的原生非客户区](/features/windows/frameless/#native-non-client-regions-on-windows)。

**任务栏闪烁：** 适合在窗口最小化时发出通知。

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

**背景类型：**

- `MacBackdropNormal` - 标准窗口
- `MacBackdropTranslucent` - 半透明背景；要使 WebView 透明，**需要使用私有 API**。
- `MacBackdropTransparent` - 完全透明；要使 WebView 透明，**需要使用私有 API**。
- `MacBackdropLiquidGlass` - 玻璃背景；要使 WebView 透明，**需要使用私有 API**。

请使用`-tags private_mac_apis`进行构建，以便这些效果能透过 WebView 显示。否则，WebView 将保持不透明。`TitleBar.AppearsTransparent`本身使用公共 API。请参阅[私有 macOS API](/guides/build/private-macos-apis/)。

**集合行为：** 控制窗口在多个 Space 间的行为：

- `MacWindowCollectionBehaviorCanJoinAllSpaces` - 在所有 Space 中可见
- `MacWindowCollectionBehaviorFullScreenAuxiliary` - 可覆盖全屏应用

**原生全屏：** macOS 全屏模式会创建一个新的 Space（虚拟桌面）。

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

**桌面环境说明：**

- GNOME：完全支持
- KDE Plasma：完全支持
- XFCE：部分支持
- 其他：支持情况各异

**平铺式窗口管理器（Hyprland、Sway、i3 等）：**

- `Minimise()`和`Maximise()`可能无法按预期工作——窗口几何属性由窗口管理器控制
- `SetSize()`和`SetPosition()`请求仅供参考，可能会被忽略
- `Fullscreen()`通常会按预期工作
- 某些窗口管理器不支持窗口置顶

@end

## 常见模式

### 启动画面

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

### 设置窗口

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

### 关闭前确认

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

## 最佳实践

### ✅ 应该做

- **为重要窗口命名** — 便于日后查找
- **设置最小尺寸** — 防止布局变得无法使用
- **将窗口居中** — 用户体验优于随机定位
- **处理关闭事件** — 防止数据丢失
- **在所有平台上测试** — 行为因平台而异
- **使用合适的尺寸** — 考虑不同的屏幕尺寸

### ❌ 不应该做

- **不要创建过多窗口** — 会让用户感到困惑
- **不要忘记关闭窗口** — 否则会导致内存泄漏
- **不要硬编码窗口位置** — 屏幕尺寸各不相同
- **不要忽略平台差异** — 应进行充分测试
- **不要阻塞 UI 线程** — 对耗时操作使用 goroutine

## 故障排除

### 窗口未显示

**可能的原因：**

1. 创建窗口时将其设为了隐藏状态
2. 窗口位于屏幕之外
3. 窗口被其他窗口遮挡

**解决方案：**

```go
window.Show()
window.Center()
window.Focus()
```

### 窗口尺寸不正确

**原因：** Windows/Linux 上的 DPI 缩放

**解决方案：**

```go
// Wails handles DPI automatically
// Just use logical pixels
window.SetSize(800, 600)
```

### 窗口立即关闭

**原因：** 最后一个窗口关闭时应用程序退出

**解决方案：**

```go
app := application.New(application.Options{
    Mac: application.MacOptions{
        ApplicationShouldTerminateAfterLastWindowClosed: false,
    },
})
```

## 后续步骤

@cards{cols="2"}
⚙ 窗口选项
所有窗口选项的完整参考。

[了解更多 →](/features/windows/options/)

---
▣ 多窗口
多窗口应用程序的设计模式。

[了解更多 →](/features/windows/multiple/)

---
★ 无边框窗口
创建自定义窗口装饰。

[了解更多 →](/features/windows/frameless/)

---
🚀 窗口事件
处理窗口生命周期事件。

[了解更多 →](/features/windows/events/)

@end

---

<strong>有疑问？</strong>请在[Discord](https://discord.gg/JDdSxwjhGf)中提问，或查看[窗口示例](https://github.com/wailsapp/wails/tree/master/v3/examples)。
