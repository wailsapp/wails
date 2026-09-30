---
title: "管理器 API"
description: "通过职责明确的管理器接口提供结构清晰的 API"
slug: "concepts/manager-api"
sourcePath: "concepts/manager-api.md"
---

Wails v3 管理器 API 将各项应用功能划分到职责明确的管理器结构体中，并通过 `*application.App` 的公开字段提供访问方式，使 API 结构清晰且易于发现。Wails 3 与 v2 完全不同——它没有用于兼容旧式 `app.NewWebviewWindow(...)` API 的逐调用封装层，因此应使用以下管理器来操控应用。

## 概述

管理器 API 将应用功能划分为十二个职责明确的领域（一个日志记录器和十一个管理器）：

- **`app.Window`** - 窗口创建、管理和回调
- **`app.ContextMenu`** - 上下文菜单注册和管理\
- **`app.KeyBinding`** - 全局按键绑定管理
- **`app.Browser`** - 浏览器集成（打开 URL 和文件）
- **`app.Env`** - 环境信息和系统状态
- **`app.Dialog`** - 文件和消息对话框操作
- **`app.Event`** - 自定义事件处理和应用事件
- **`app.Menu`** - 应用菜单管理
- **`app.Screen`** - 屏幕管理和坐标转换
- **`app.Clipboard`** - 剪贴板文本操作
- **`app.SystemTray`** - 系统托盘图标的创建和管理
- **`app.Autostart`** - 注册应用，使其在用户登录时启动

## 优势

- **更易发现** - IDE 自动补全会显示结构清晰的 API 接口
- **改进代码组织** - 将相关方法归为一组
- **增强可维护性** - 在各管理器之间分离关注点
- **便于未来扩展** - 更容易向特定领域添加新功能

## 用法

管理器 API 以结构清晰的方式提供对所有应用功能的访问：

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

## 管理器参考

### 窗口管理器

管理窗口的创建、获取和生命周期回调。

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

### 事件管理器

处理自定义事件并监听应用事件。

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

### 浏览器管理器

提供用于打开 URL 和文件的浏览器集成功能。

```go
// Open URLs and files in default browser
err := app.Browser.OpenURL("https://wails.io")
err := app.Browser.OpenFile("/path/to/document.pdf")
```

### 环境管理器

提供对系统环境信息的访问。

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

### 对话框管理器

以结构清晰的方式提供对文件和消息对话框的访问。

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

### 菜单管理器

创建和管理应用菜单。

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

### 按键绑定管理器

动态管理全局按键绑定。

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

### 上下文菜单管理器

高级上下文菜单管理（面向库作者）。

```go
// Create and register context menu
menu := app.ContextMenu.New()
app.ContextMenu.Add("myMenu", menu)

// Retrieve context menu
menu, exists := app.ContextMenu.Get("myMenu")

// Remove context menu
app.ContextMenu.Remove("myMenu")
```

### 屏幕管理器

面向多显示器配置的屏幕管理和坐标转换。

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

### 剪贴板管理器

用于读取和写入文本的剪贴板操作。

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

### SystemTray 管理器

创建和管理系统托盘图标。

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

### 自动启动管理器

注册应用，使其在用户登录时启动。它会根据各平台选择合适的原生机制：macOS 上使用 SMAppService 或 LaunchAgent plist，Windows 上使用 `HKCU\…\Run` 注册表项，Linux 上使用 XDG `.desktop` 条目。

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

有关各平台的行为、标识符规则以及过期项检测保证，请参阅[自动启动功能页面](/features/autostart/basics/)。
