---
title: "应用程序 API"
description: "应用程序 API 完整参考"
slug: "reference/application"
sourcePath: "reference/application.md"
---

## 概述

`Application` 是 Wails 应用的核心。它管理窗口、服务和事件，并提供对所有平台功能的访问。

## 创建应用程序

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

## 核心方法

### Run()

启动应用程序事件循环。

```go
func (a *App) Run() error
```

**示例：**

```go
err := app.Run()
if err != nil {
    log.Fatal(err)
}
```

<strong>返回值：</strong>启动失败时返回错误

### Quit()

正常关闭应用程序。

```go
func (a *App) Quit()
```

**示例：**

```go
// In a menu handler
menu.Add("Quit").OnClick(func(ctx *application.Context) {
    app.Quit()
})
```

### Config()

返回应用程序配置。

```go
func (a *App) Config() Options
```

**示例：**

```go
config := app.Config()
fmt.Println("App name:", config.Name)
```

## 窗口管理

### app.Window.New()

使用默认选项创建新的 WebView 窗口。

```go
func (wm *WindowManager) New() *WebviewWindow
```

**示例：**

```go
window := app.Window.New()
window.Show()
```

### app.Window.NewWithOptions()

使用自定义选项创建新的 WebView 窗口。

```go
func (wm *WindowManager) NewWithOptions(options WebviewWindowOptions) *WebviewWindow
```

**示例：**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "My Window",
    Width:  800,
    Height: 600,
    BackgroundColour: application.NewRGB(255, 255, 255),
})
```

### app.Window.GetByName()

按名称获取窗口。返回该窗口以及是否找到该窗口。

```go
func (wm *WindowManager) GetByName(name string) (Window, bool)
```

**示例：**

```go
if window, ok := app.Window.GetByName("main"); ok {
    window.Show()
}
```

### app.Window.GetAll()

返回应用程序的所有窗口。

```go
func (wm *WindowManager) GetAll() []Window
```

**示例：**

```go
windows := app.Window.GetAll()
for _, window := range windows {
    fmt.Println("Window:", window.Name())
}
```

## 管理器

Application 通过属性提供对各种管理器的访问：

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

### 用法示例

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

## 服务管理

### RegisterService()

向应用程序注册服务。

```go
func (a *App) RegisterService(service Service)
```

`RegisterService` 不返回任何内容；服务初始化错误会在 `app.Run()` 期间通过 `ServiceStartup` 失败显现。

**示例：**

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

## 事件管理

### app.Event.Emit()

发出自定义事件。如果钩子取消了此次发出，则返回 `true`。

```go
func (em *EventManager) Emit(name string, data ...any) bool
```

**示例：**

```go
// Emit event with data
app.Event.Emit("user-logged-in", map[string]interface{}{
    "username": "john",
    "timestamp": time.Now(),
})
```

### app.Event.On()

监听自定义事件。返回用于取消订阅的 `func()`。

```go
func (em *EventManager) On(name string, callback func(*CustomEvent)) func()
```

**示例：**

```go
app.Event.On("user-logged-in", func(e *application.CustomEvent) {
    data := e.Data.(map[string]interface{})
    username := data["username"].(string)
    fmt.Println("User logged in:", username)
})
```

### app.Event.OnApplicationEvent()

监听应用程序生命周期事件。`eventType` 参数的类型为 `events.ApplicationEventType`（来自 `events` 包）。

```go
func (em *EventManager) OnApplicationEvent(
    eventType events.ApplicationEventType,
    callback func(*ApplicationEvent),
) func()
```

**示例：**

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

## 对话框方法

通过 `app.Dialog` 管理器访问对话框。完整参考请参阅[对话框 API](/reference/dialogs/)。

### 消息对话框

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

### 询问对话框

询问对话框使用按钮回调处理用户响应：

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

### 文件对话框

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

## 日志记录器

应用程序提供结构化日志记录器：

```go
app.Logger.Info("Message", "key", "value")
app.Logger.Error("Error occurred", "error", err)
app.Logger.Debug("Debug info")
app.Logger.Warn("Warning message")
```

**示例：**

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

## 原始消息处理

对于需要直接从底层控制前端到后端通信的应用程序，Wails 提供了`RawMessageHandler`选项。此选项会绕过标准绑定系统。

@note{type="info"}
仅应将原始消息作为最后手段使用。标准绑定系统经过高度优化，足以满足几乎所有应用程序的需求。只有在对应用程序进行性能分析并确认绑定是瓶颈后，才应使用原始消息。

@end

### RawMessageHandler

`RawMessageHandler` 是 `application.Options` 的字段，而不是方法。对于前端通过 `System.invoke()` 发送的每条原始消息，运行时都会调用该字段。

```go
type Options struct {
    // ... other fields ...
    RawMessageHandler func(window Window, message string, originInfo *OriginInfo)
}
```

`OriginInfo` 包含 `Origin`、`TopOrigin` 和 `IsMainFrame`（不同平台会填充不同的子集——有关各平台的对应矩阵，请参阅[原始消息指南](/guides/raw-messages/)）。

**示例：**

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

有关更多详细信息，请参阅[原始消息指南](/guides/raw-messages/)。

## 平台特定选项

### Windows 选项

在应用程序级别配置 Windows 特定行为：

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

**浏览器标志：**

- `EnabledFeatures` - 要启用的 WebView2 功能标志
- `DisabledFeatures` - 要禁用的 WebView2 功能标志
- `AdditionalBrowserArgs` - Chromium 命令行参数

有关详细文档，请参阅[窗口选项 - 应用程序级 Windows 选项](/features/windows/options/#application-level-windows-options)。

### Mac 选项

```go
app := application.New(application.Options{
    Name: "My App",
    Mac: application.MacOptions{
        ActivationPolicy: application.ActivationPolicyRegular,
        ApplicationShouldTerminateAfterLastWindowClosed: true,
    },
})
```

### Linux 选项

```go
app := application.New(application.Options{
    Name: "My App",
    Linux: application.LinuxOptions{
        ProgramName:                   "my-app",
        DisableQuitOnLastWindowClosed: false,
    },
})
```

## 完整应用程序示例

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
