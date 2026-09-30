---
title: "窗口 API"
description: "窗口 API 完整参考"
slug: "reference/window"
sourcePath: "reference/window.md"
---

## 概述

窗口 API 提供用于控制窗口外观、行为和生命周期的方法。可通过窗口实例或`app.Window`管理器访问。

`Window`是由`*application.WebviewWindow`实现的接口；下列方法签名属于`*WebviewWindow`。许多修改方法会返回`Window`以支持链式调用——各方法的返回值均在相应说明中注明。

**常用操作：**

- 创建并显示窗口
- 控制大小、位置和状态
- 处理窗口事件
- 管理窗口内容
- 配置外观和行为

## 可见性

### Show()

显示窗口。如果窗口处于隐藏状态，则使其可见。返回接收者以支持链式调用。

```go
func (w *WebviewWindow) Show() Window
```

**示例：**

```go
window := app.Window.New()
window.Show()
```

### Hide()

隐藏窗口，但不将其关闭。窗口仍保留在内存中，可以再次显示。返回接收者以支持链式调用。

```go
func (w *WebviewWindow) Hide() Window
```

**示例：**

```go
// Hide window temporarily
window.Hide()

// Show it again later
window.Show()
```

**使用场景：**

- 隐藏到系统托盘的应用程序
- 重复使用窗口的向导流程
- 在操作期间暂时隐藏窗口

### Close()

关闭窗口。此操作会触发`WindowClosing`事件。

```go
func (w *WebviewWindow) Close()
```

**示例：**

```go
window.Close()
```

<strong>注意：</strong>如果已注册的钩子调用`event.Cancel()`，则会阻止窗口关闭。

## 窗口属性

### SetTitle()

设置窗口标题栏文本。返回接收者以支持链式调用。

```go
func (w *WebviewWindow) SetTitle(title string) Window
```

**参数：**

- `title` - 新窗口标题

**示例：**

```go
window.SetTitle("My Application - Document.txt")
```

### Name()

返回窗口的唯一名称标识符。

```go
func (w *WebviewWindow) Name() string
```

**示例：**

```go
name := window.Name()
fmt.Println("Window name:", name)

// Retrieve window by name later
if w, ok := app.Window.GetByName(name); ok {
    w.Focus()
}
```

## 大小和位置

### SetSize()

设置窗口尺寸，单位为像素。返回接收者以支持链式调用。

```go
func (w *WebviewWindow) SetSize(width, height int) Window
```

**参数：**

- `width` - 窗口宽度，单位为像素
- `height` - 窗口高度，单位为像素

**示例：**

```go
window.SetSize(1024, 768)
```

### Size()

返回当前窗口尺寸。

```go
func (w *WebviewWindow) Size() (width, height int)
```

**示例：**

```go
width, height := window.Size()
fmt.Printf("Window is %dx%d\n", width, height)
```

### SetMinSize() / SetMaxSize()

设置窗口的最小和最大尺寸。这两个方法均返回接收者以支持链式调用。

```go
func (w *WebviewWindow) SetMinSize(width, height int) Window
func (w *WebviewWindow) SetMaxSize(width, height int) Window
```

**示例：**

```go
// Prevent window from being too small
window.SetMinSize(800, 600)

// Prevent window from being too large
window.SetMaxSize(1920, 1080)
```

### SetPosition()

设置窗口相对于屏幕左上角的位置。

```go
func (w *WebviewWindow) SetPosition(x, y int)
```

**参数：**

- `x` - 水平位置，单位为像素
- `y` - 垂直位置，单位为像素

**示例：**

```go
// Position window at top-left
window.SetPosition(0, 0)

// Position window 100px from top-left
window.SetPosition(100, 100)
```

### Position()

返回窗口的当前位置。

```go
func (w *WebviewWindow) Position() (x, y int)
```

**示例：**

```go
x, y := window.Position()
fmt.Printf("Window is at (%d, %d)\n", x, y)
```

### Center()

将窗口置于屏幕中央。

```go
func (w *WebviewWindow) Center()
```

**示例：**

```go
window := app.Window.New()
window.Center()
window.Show()
```

<strong>注意：</strong>窗口会在主显示器上居中。对于多显示器配置，请参阅屏幕 API。

### Focus()

将窗口置于最前端并使其获得键盘焦点。

```go
func (w *WebviewWindow) Focus()
```

**示例：**

```go
// Bring window to front
window.Focus()
```

## 窗口状态

### Minimise() / UnMinimise()

将窗口最小化到任务栏或程序坞，或将其还原。`Minimise()`返回接收者以支持链式调用；`UnMinimise()`不返回任何值。

```go
func (w *WebviewWindow) Minimise() Window
func (w *WebviewWindow) UnMinimise()
```

**示例：**

```go
// Minimise window
window.Minimise()

// Restore from minimised state
window.UnMinimise()
```

### Maximise() / UnMaximise()

将窗口最大化以填满屏幕，或将其还原到之前的大小。`Maximise()`返回接收者以支持链式调用；`UnMaximise()`不返回任何值。

```go
func (w *WebviewWindow) Maximise() Window
func (w *WebviewWindow) UnMaximise()
```

**示例：**

```go
// Maximise window
window.Maximise()

// Restore to previous size
window.UnMaximise()
```

### Fullscreen() / UnFullscreen() / ToggleFullscreen()

进入或退出全屏模式。`Fullscreen()`返回接收者以支持链式调用。

```go
func (w *WebviewWindow) Fullscreen() Window
func (w *WebviewWindow) UnFullscreen()
func (w *WebviewWindow) ToggleFullscreen()
```

**示例：**

```go
// Enter fullscreen
window.Fullscreen()

// Exit fullscreen
window.UnFullscreen()

// Or toggle
window.ToggleFullscreen()
```

没有`SetFullscreen(bool)`方法。

### IsMinimised() / IsMaximised() / IsFullscreen()

检查窗口的当前状态。

```go
func (w *WebviewWindow) IsMinimised() bool
func (w *WebviewWindow) IsMaximised() bool
func (w *WebviewWindow) IsFullscreen() bool
```

**示例：**

```go
if window.IsMinimised() {
    window.UnMinimise()
}

if window.IsMaximised() {
    fmt.Println("Window is maximised")
}

if window.IsFullscreen() {
    window.UnFullscreen()
}
```

## 窗口内容

### SetURL()

在窗口中导航到指定的 URL。返回接收者以支持链式调用。

```go
func (w *WebviewWindow) SetURL(url string) Window
```

**参数：**

- `url` - 要导航到的 URL（对于嵌入式资源，可以使用`http://wails.localhost/`）

**示例：**

```go
// Navigate to embedded page
window.SetURL("http://wails.localhost/settings.html")

// Navigate to external URL (if allowed)
window.SetURL("https://wails.io")
```

### SetHTML()

直接使用 HTML 字符串设置窗口内容。返回接收者以支持链式调用。

```go
func (w *WebviewWindow) SetHTML(html string) Window
```

**参数：**

- `html` - 要显示的 HTML 内容

**示例：**

```go
html := `
<!DOCTYPE html>
<html>
<head><title>Dynamic Content</title></head>
<body>
    <h1>Hello from Go!</h1>
    <p>This content was generated dynamically.</p>
</body>
</html>
`
window.SetHTML(html)
```

**使用场景：**

- 动态生成内容
- 无需前端构建流程的简单窗口
- 错误页面或启动画面

### Reload()

重新加载当前窗口内容。

```go
func (w *WebviewWindow) Reload()
```

**示例：**

```go
// Reload current page
window.Reload()
```

<strong>注意：</strong>此方法适用于开发期间或需要刷新内容时。

## 窗口事件

Wails 提供两种处理窗口事件的方法：

- **OnWindowEvent()** - 监听窗口事件（无法阻止事件）。
- **RegisterHook()** - 挂接窗口事件（可通过调用`event.Cancel()`阻止事件）。

### OnWindowEvent()

为窗口事件注册回调。返回取消订阅函数。

```go
func (w *WebviewWindow) OnWindowEvent(
    eventType events.WindowEventType,
    callback func(event *WindowEvent),
) func()
```

**示例：**

```go
import "github.com/wailsapp/wails/v3/pkg/events"

// Listen for window focus
window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
    app.Logger.Info("Window gained focus")
})

// Listen for window lost focus
window.OnWindowEvent(events.Common.WindowLostFocus, func(e *application.WindowEvent) {
    app.Logger.Info("Window lost focus")
})

// Listen for window resize
window.OnWindowEvent(events.Common.WindowDidResize, func(e *application.WindowEvent) {
    app.Logger.Info("Window resized")
})
```

**常见窗口事件：**

- `events.Common.WindowClosing` - 窗口即将关闭
- `events.Common.WindowFocus` - 窗口获得焦点
- `events.Common.WindowLostFocus` - 窗口失去焦点
- `events.Common.WindowDidMove` - 窗口已移动
- `events.Common.WindowDidResize` - 窗口大小已调整
- `events.Common.WindowMinimise` - 窗口已最小化
- `events.Common.WindowMaximise` - 窗口已最大化
- `events.Common.WindowFullscreen` - 窗口已进入全屏模式
- `events.Common.WindowRuntimeReady` - 窗口内运行时已初始化

### RegisterHook()

为窗口事件注册钩子。钩子先于监听器运行，并可通过调用`event.Cancel()`阻止事件。返回取消订阅函数。

```go
func (w *WebviewWindow) RegisterHook(
    eventType events.WindowEventType,
    callback func(event *WindowEvent),
) func()
```

**示例 - 阻止窗口关闭：**

```go
import "github.com/wailsapp/wails/v3/pkg/events"

window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    confirm := app.Dialog.Question().
        SetTitle("Confirm Close").
        SetMessage("Are you sure you want to close this window?")

    yes := confirm.AddButton("Yes")
    no := confirm.AddButton("No")
    confirm.SetDefaultButton(yes)
    confirm.SetCancelButton(no)

    no.OnClick(func() {
        e.Cancel() // Prevent window from closing
    })

    confirm.Show()
})
```

**示例 - 关闭前保存：**

```go
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    if !hasUnsavedChanges {
        return
    }

    dlg := app.Dialog.Question().
        SetTitle("Unsaved Changes").
        SetMessage("Save changes before closing?")

    save := dlg.AddButton("Save")
    discard := dlg.AddButton("Don't Save")
    cancel := dlg.AddButton("Cancel")
    dlg.SetDefaultButton(save)
    dlg.SetCancelButton(cancel)

    save.OnClick(func() { saveData() })
    cancel.OnClick(func() { e.Cancel() })
    _ = discard // allow close

    dlg.Show()
})
```

### EmitEvent()

向窗口的前端发出自定义事件。如果发出操作被钩子取消，则返回`true`。

```go
func (w *WebviewWindow) EmitEvent(name string, data ...any) bool
```

**参数：**

- `name` - 事件名称
- `data` - 随事件发送的可选数据

**示例：**

```go
// Send data to specific window
window.EmitEvent("data-updated", map[string]any{
    "count":  42,
    "status": "success",
})
```

**前端（JavaScript）：**

```javascript
import { Events } from '@wailsio/runtime'

Events.On('data-updated', (data) => {
    console.log('Count:', data.count)
    console.log('Status:', data.status)
})
```

## 其他方法

### SetEnabled()

启用或禁用用户与窗口交互。

```go
func (w *WebviewWindow) SetEnabled(enabled bool)
```

**示例：**

```go
// Disable window during long operation
window.SetEnabled(false)

// Perform operation
performLongOperation()

// Re-enable window
window.SetEnabled(true)
```

### SetBackgroundColour()

设置窗口的背景颜色（在内容加载前显示）。返回接收者以支持链式调用。

```go
func (w *WebviewWindow) SetBackgroundColour(colour RGBA) Window
```

`RGBA`为`application.RGBA{Red, Green, Blue, Alpha uint8}`。请使用辅助函数`application.NewRGB(r, g, b)`（alpha 值为255）或`application.NewRGBA(r, g, b, a)`。

**示例：**

```go
// White background
window.SetBackgroundColour(application.NewRGB(255, 255, 255))

// Dark background with full alpha
window.SetBackgroundColour(application.NewRGBA(30, 30, 30, 255))
```

### SetResizable()

控制用户能否调整窗口大小。返回接收者以支持链式调用。

```go
func (w *WebviewWindow) SetResizable(resizable bool) Window
```

**示例：**

```go
// Make window fixed size
window.SetResizable(false)
```

### SetAlwaysOnTop()

设置窗口是否始终位于其他窗口之上。返回接收者以支持链式调用。

```go
func (w *WebviewWindow) SetAlwaysOnTop(alwaysOnTop bool) Window
```

**示例：**

```go
// Keep window on top
window.SetAlwaysOnTop(true)
```

### Print()

为窗口内容打开原生打印对话框。

```go
func (w *WebviewWindow) Print() error
```

<strong>返回值：</strong>如果打印失败，则返回错误。

**示例：**

```go
if err := window.Print(); err != nil {
    log.Println("Print failed:", err)
}
```

### AttachModal()

将第二个窗口附加为表单式模态窗口。

```go
func (w *WebviewWindow) AttachModal(modalWindow Window)
```

**参数：**

- `modalWindow` - 要附加为模态窗口的窗口

**平台支持：**

- **macOS**：完全支持（以表单形式显示）
- **Windows**：不支持
- **Linux**：不支持

**示例：**

```go
modalWindow := app.Window.New()
window.AttachModal(modalWindow)
```

## 平台特定选项

### Linux

Linux 窗口通过`LinuxWindow`支持以下平台特定选项：

#### MenuStyle

控制应用程序菜单的显示方式。此选项适用于默认的 GTK4 构建，在旧版`-tags gtk3`构建中会被忽略。

| 值 | 说明 |
| --- | --- |
| `LinuxMenuStyleMenuBar` | 标题栏下方的传统菜单栏（默认） |
| `LinuxMenuStylePrimaryMenu` | 标题栏中的主菜单按钮（GNOME 风格） |

**示例：**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "My Application",
    Linux: application.LinuxWindow{
        MenuStyle: application.LinuxMenuStylePrimaryMenu,
    },
})
window.SetMenu(menu)
```

<strong>注意：</strong>主菜单样式遵循 GNOME 人机界面指南，在标题栏中显示汉堡菜单按钮（☰）。这是现代 GNOME 应用程序的推荐样式。

## 完整示例

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)

func main() {
    app := application.New(application.Options{
        Name: "Window API Demo",
    })

    // Create window with options
    window := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:            "My Application",
        Width:            1024,
        Height:           768,
        MinWidth:         800,
        MinHeight:        600,
        BackgroundColour: application.NewRGB(255, 255, 255),
        URL:              "http://wails.localhost/",
    })

    // Configure window behaviour
    window.SetResizable(true)
    window.SetMinSize(800, 600)
    window.SetMaxSize(1920, 1080)

    // Confirm-before-close hook
    window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
        dlg := app.Dialog.Question().
            SetTitle("Confirm Close").
            SetMessage("Are you sure you want to close this window?")

        yes := dlg.AddButton("Yes")
        no := dlg.AddButton("No")
        dlg.SetDefaultButton(yes)
        dlg.SetCancelButton(no)
        no.OnClick(func() { e.Cancel() })

        dlg.Show()
    })

    // Listen for window events
    window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
        window.SetTitle("My Application (Active)")
        app.Logger.Info("Window gained focus")
    })

    window.OnWindowEvent(events.Common.WindowLostFocus, func(e *application.WindowEvent) {
        window.SetTitle("My Application")
        app.Logger.Info("Window lost focus")
    })

    // Position and show window
    window.Center()
    window.Show()

    app.Run()
}
```
