---
title: "应用生命周期"
description: "了解 Wails 应用从启动到关闭的生命周期"
slug: "concepts/lifecycle"
sourcePath: "concepts/lifecycle.md"
---

## 了解应用生命周期

桌面应用具有从启动到关闭的生命周期。Wails v3 提供<strong>服务</strong>、<strong>事件</strong>和<strong>钩子</strong>，用于有效管理此生命周期。

## 生命周期阶段

```d2
direction: down

Start: 应用启动 {
  shape: oval
  style.fill: "#10B981"
}

Init: 初始化 {
  Parse: 解析选项 {
    shape: rectangle
  }
  Register: 注册服务 {
    shape: rectangle
  }
  Setup: 设置运行时 {
    shape: rectangle
  }
}

AppRun: app.Run() {
  shape: rectangle
  style.fill: "#3B82F6"
}

ServiceStartup: 启动服务 {
  shape: rectangle
  style.fill: "#8B5CF6"
}

EventLoop: 事件循环 {
  Process: 处理事件 {
    shape: rectangle
  }
  Handle: 处理消息 {
    shape: rectangle
  }
  Update: 更新界面 {
    shape: rectangle
  }
}

QuitSignal: 退出信号 {
  shape: diamond
  style.fill: "#F59E0B"
}

ShouldQuit: ShouldQuit 检查 {
  shape: rectangle
  style.fill: "#3B82F6"
}

OnShutdown: OnShutdown 回调 {
  shape: rectangle
  style.fill: "#3B82F6"
}

ServiceShutdown: 关闭服务 {
  shape: rectangle
  style.fill: "#8B5CF6"
}

Cleanup: 清理 {
  Close: 关闭窗口 {
    shape: rectangle
  }
  Release: 释放资源 {
    shape: rectangle
  }
}

End: 应用结束 {
  shape: oval
  style.fill: "#EF4444"
}

Start -> Init.Parse
Init.Parse -> Init.Register
Init.Register -> Init.Setup
Init.Setup -> AppRun
AppRun -> ServiceStartup
ServiceStartup -> EventLoop.Process
EventLoop.Process -> EventLoop.Handle
EventLoop.Handle -> EventLoop.Update
EventLoop.Update -> EventLoop.Process: 循环
EventLoop.Process -> QuitSignal: 用户退出
QuitSignal -> ShouldQuit: 检查是否允许？
ShouldQuit -> EventLoop.Process: 已拒绝
ShouldQuit -> OnShutdown: 已允许
OnShutdown -> ServiceShutdown
ServiceShutdown -> Cleanup.Close
Cleanup.Close -> Cleanup.Release
Cleanup.Release -> End
```

### 1. 创建应用

使用`application.New()`创建应用：

```go
app := application.New(application.Options{
    Name:        "My App",
    Description: "An application built with Wails",
    Services: []application.Service{
        application.NewService(&MyService{}),
    },
    Assets: application.AssetOptions{
        Handler: application.BundledAssetFileServer(assets),
    },
})
```

**执行的操作：**

1. 解析并验证选项
2. 注册服务（但尚未启动）
3. 配置资源服务器
4. 设置运行时

### 2. 运行应用

调用`app.Run()`启动应用：

```go
err := app.Run()  // Blocks until quit
if err != nil {
    log.Fatal(err)
}
```

**执行的操作：**

1. 按注册顺序启动服务
2. 激活事件监听器
3. 可以创建窗口
4. 开始事件循环

### 3. 事件循环

应用进入事件循环，并在其中运行大部分时间：

- 处理操作系统事件（鼠标、键盘和窗口事件）
- 处理 Go 到 JS 的消息
- 执行 JS 到 Go 的调用
- 渲染 UI 更新

### 4. 关闭

应用退出时：

1. 检查`ShouldQuit`回调（如果已设置）
2. 执行`OnShutdown`回调
3. 按相反顺序关闭服务
4. 关闭窗口
5. 释放资源

## 服务生命周期

服务是 Wails v3 中管理生命周期的主要方式。服务通过接口提供启动和关闭钩子。有关服务的完整文档，请参阅[服务指南](/features/bindings/services/)。

### 创建服务

```go
type MyService struct {
    db *sql.DB
}

// ServiceStartup is called when the application starts
func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    var err error
    s.db, err = sql.Open("sqlite3", "app.db")
    if err != nil {
        return err  // Startup aborts if error returned
    }

    // Run migrations
    if err := s.runMigrations(); err != nil {
        return err
    }

    return nil
}

// ServiceShutdown is called when the application shuts down
func (s *MyService) ServiceShutdown() error {
    if s.db != nil {
        return s.db.Close()
    }
    return nil
}
```

### 注册服务

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&MyService{}),
        application.NewService(&AnotherService{}),
    },
})
```

**要点：**

- 服务按注册顺序启动
- 服务按注册顺序的<strong>相反顺序</strong>关闭
- 如果服务的`ServiceStartup`返回错误，应用将中止
- 传递给`ServiceStartup`的`ctx`会在关闭开始时被取消

### 使用应用上下文

传递给`ServiceStartup`的上下文在应用程序的整个生命周期内均有效：

```go
func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    // Start a background task that respects shutdown
    go func() {
        ticker := time.NewTicker(5 * time.Minute)
        defer ticker.Stop()

        for {
            select {
            case <-ticker.C:
                s.performBackgroundSync()
            case <-ctx.Done():
                // Application is shutting down
                return
            }
        }
    }()

    return nil
}
```

也可以从应用程序实例访问该上下文：

```go
app := application.Get()
ctx := app.Context()
```

## 应用程序级钩子

这些是`application.Options`中的便捷回调，无需创建完整的服务即可挂接到应用程序生命周期。它们适用于简单的清理任务、退出确认，或需要在关闭流程的特定阶段运行代码的情况。

如果需要通过启动逻辑、依赖注入或有状态资源进行更复杂的生命周期管理，请改用[服务](#heading-2)。

### ShouldQuit

每当收到退出请求时，都会调用`ShouldQuit`回调——无论是用户关闭最后一个窗口、按下 Cmd+Q (macOS) / Alt+F4 (Windows)，还是以编程方式调用`app.Quit()`。

**返回值：**

- 返回`true`以允许继续退出（应用程序将关闭）
- 返回`false`以取消退出（应用程序将继续运行）

你可以借此拦截退出请求，并选择阻止退出，例如提示用户存在未保存的更改：

```go
app := application.New(application.Options{
    ShouldQuit: func() bool {
        if !hasUnsavedChanges() {
            return true // No unsaved changes, allow quit
        }

        // Prompt the user — MessageDialog.Show() blocks and returns nothing.
        // The button's OnClick callback fires for whichever button the user picks.
        shouldQuit := false

        dlg := application.Get().Dialog.Question().
            SetTitle("Unsaved Changes").
            SetMessage("You have unsaved changes. Quit anyway?")

        quit := dlg.AddButton("Quit")
        cancel := dlg.AddButton("Cancel")
        dlg.SetDefaultButton(cancel)
        dlg.SetCancelButton(cancel)

        quit.OnClick(func() { shouldQuit = true })

        dlg.Show()
        return shouldQuit
    },
})
```

如果未设置`ShouldQuit`，应用程序将在收到请求时立即退出。

**调用 ShouldQuit 的情况：**

- 用户关闭最后一个窗口（除非设置了`DisableQuitOnLastWindowClosed`）
- 用户在 macOS 上按下 Cmd+Q
- 用户在 Windows 上按下 Alt+F4（焦点位于最后一个窗口时）
- 代码调用`app.Quit()`

**不调用 ShouldQuit 的情况：**

- 进程被终止（SIGKILL、通过任务管理器强制退出）
- 直接调用`os.Exit()`

### OnShutdown

确认应用程序即将退出时，将调用`OnShutdown`回调（如果设置了`ShouldQuit`，则在其返回`true`后调用）。可使用此回调执行保存状态、关闭数据库连接或释放资源等清理任务。

```go
app := application.New(application.Options{
    OnShutdown: func() {
        // Save application state
        saveState()

        // Close connections
        cleanup()
    },
})
```

在应用程序生命周期内，也可以随时以编程方式注册其他关闭回调：

```go
app.OnShutdown(func() {
    log.Println("Application shutting down...")
})
```

多个回调按注册顺序执行。关闭流程会阻塞，直到所有回调执行完毕。

<strong>重要提示：</strong>请确保关闭回调能快速完成（少于1秒）。操作系统可能会强制终止退出耗时过长的应用程序，这可能会中断清理工作并导致数据丢失。

### PostShutdown

所有关闭任务完成后、进程终止前，将调用`PostShutdown`回调。此时应用程序实例已无法使用——窗口已关闭、服务已停止、资源已释放。

此回调主要用于：

- 在完成所有其他清理工作后必须执行的最终日志记录
- 测试和调试关闭行为
- 在`app.Run()`不会返回的平台上确保代码得以运行

```go
app := application.New(application.Options{
    PostShutdown: func() {
        // Final logging
        log.Println("Application terminated cleanly")

        // Flush any buffered logs
        logger.Sync()
    },
})
```

<strong>注意：</strong>请勿尝试在`PostShutdown`中使用应用程序功能（窗口、对话框等），因为这些功能已不可用。

## 基于事件的生命周期

Wails 提供了一个事件系统，可在应用程序中发生各种情况时通知你，例如窗口打开、应用程序启动、主题更改等。你可以监听这些事件，以便在不阻塞或拦截事件的情况下响应生命周期变化。

对于窗口事件，如果需要拦截并取消操作（例如阻止窗口关闭），还可以使用`RegisterHook`而不是`OnWindowEvent`。请参阅下文的[窗口钩子](#heading-10)。

有关事件系统的完整文档，请参阅[事件指南](/features/events/system/)。

### 应用程序事件

监听应用程序生命周期事件：

```go
app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(event *application.ApplicationEvent) {
    app.Logger.Info("Application has started!")
})
```

还可以使用平台特定的事件：

```go
// macOS
app.Event.OnApplicationEvent(events.Mac.ApplicationDidFinishLaunching, func(event *application.ApplicationEvent) {
    // Handle macOS launch
})

app.Event.OnApplicationEvent(events.Mac.ApplicationWillTerminate, func(event *application.ApplicationEvent) {
    // Handle macOS termination
})

// Windows
app.Event.OnApplicationEvent(events.Windows.ApplicationStarted, func(event *application.ApplicationEvent) {
    // Handle Windows start
})
```

### 窗口事件

监听窗口生命周期事件：

```go
window := app.Window.New()

window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
    app.Logger.Info("Window gained focus")
})

window.OnWindowEvent(events.Common.WindowClosing, func(e *application.WindowEvent) {
    app.Logger.Info("Window is closing")
})
```

### 窗口钩子（可取消事件）

需要<strong>取消</strong>事件时，请使用`RegisterHook`而不是`OnWindowEvent`：

```go
window := app.Window.New()

var countdown = 3

window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    countdown--
    if countdown > 0 {
        app.Logger.Info("Not closing yet!", "remaining", countdown)
        e.Cancel()  // Prevent the window from closing
        return
    }
    app.Logger.Info("Window closing now")
})
```

**OnWindowEvent 与 RegisterHook 的区别：**

- `OnWindowEvent`：在事件发生时通知你（无法取消）
- `RegisterHook`：允许你拦截并视情况取消事件

## 窗口生命周期

窗口有自己的生命周期，从创建到销毁。每个窗口独立加载其前端内容，并且可以随时显示、隐藏或关闭。当用户尝试关闭窗口时，你可以使用`RegisterHook`拦截此操作，以提示用户确认，或隐藏窗口而不将其销毁。

有关窗口的完整文档，请参阅[窗口指南](/features/windows/basics/)。

```d2
direction: down

Create: 创建窗口 {
  shape: oval
  style.fill: "#10B981"
}

Load: 加载前端 {
  shape: rectangle
}

Show: 显示窗口 {
  shape: rectangle
}

Active: 窗口处于活动状态 {
  Events: 处理事件 {
    shape: rectangle
  }
}

CloseRequest: 关闭请求 {
  shape: diamond
  style.fill: "#F59E0B"
}

Hook: WindowClosing 钩子 {
  shape: rectangle
  style.fill: "#3B82F6"
}

Destroy: 销毁窗口 {
  shape: rectangle
}

End: 窗口已关闭 {
  shape: oval
  style.fill: "#EF4444"
}

Create -> Load
Load -> Show
Show -> Active.Events
Active.Events -> Active.Events: 循环
Active.Events -> CloseRequest: 用户关闭窗口
CloseRequest -> Hook
Hook -> Active.Events: 已取消
Hook -> Destroy: 已允许
Destroy -> End
```

### 创建窗口

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "My Window",
    Width:  800,
    Height: 600,
})
```

### 阻止窗口关闭

```go
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    if !hasUnsavedChanges() {
        return
    }

    // MessageDialog.Show() returns nothing; per-button OnClick handlers fire.
    dlg := application.Get().Dialog.Question().
        SetTitle("Unsaved Changes").
        SetMessage("Save before closing?")

    save := dlg.AddButton("Save")
    discard := dlg.AddButton("Discard")
    cancel := dlg.AddButton("Cancel")
    dlg.SetDefaultButton(save)
    dlg.SetCancelButton(cancel)

    save.OnClick(func() { saveChanges() })
    cancel.OnClick(func() { e.Cancel() }) // Prevent close
    _ = discard                            // "Discard" falls through and allows close

    dlg.Show()
})
```

### 隐藏而非关闭

系统托盘应用常用以下模式：

```go
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    window.Hide()  // Hide instead of destroy
    e.Cancel()     // Prevent actual close
})
```

## 多窗口生命周期

使用多个窗口时：

```go
mainWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Main Window",
})

settingsWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "Settings",
    Width:  400,
    Height: 600,
    Hidden: true,  // Start hidden
})
```

**默认行为因平台而异：**

| 平台 | 最后一个窗口关闭时的默认行为 |
| --- | --- |
| macOS | 应用继续运行（菜单栏保留） |
| Windows | 应用退出 |
| Linux | 应用退出 |

macOS 遵循平台原生惯例，即使没有窗口，应用通常也会在菜单栏中保持活动状态。Windows 和 Linux 默认会退出应用。

**让所有平台在最后一个窗口关闭时退出应用：**

```go
app := application.New(application.Options{
    Mac: application.MacOptions{
        ApplicationShouldTerminateAfterLastWindowClosed: true,
    },
})
```

**让所有平台在最后一个窗口关闭时继续运行应用：**

这适用于系统托盘应用或需要在后台继续运行的应用。

```go
app := application.New(application.Options{
    Windows: application.WindowsOptions{
        DisableQuitOnLastWindowClosed: true,
    },
    Linux: application.LinuxOptions{
        DisableQuitOnLastWindowClosed: true,
    },
})
```

## 常见模式

### 模式 1：数据库服务

```go
type DatabaseService struct {
    db *sql.DB
}

func (s *DatabaseService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    var err error
    s.db, err = sql.Open("sqlite3", "app.db")
    if err != nil {
        return fmt.Errorf("failed to open database: %w", err)
    }

    if err := s.db.PingContext(ctx); err != nil {
        return fmt.Errorf("failed to connect to database: %w", err)
    }

    return nil
}

func (s *DatabaseService) ServiceShutdown() error {
    if s.db != nil {
        return s.db.Close()
    }
    return nil
}

// Exported methods are available to the frontend
func (s *DatabaseService) GetUsers() ([]User, error) {
    // Query implementation
}
```

### 模式 2：配置服务

```go
type ConfigService struct {
    config *Config
    path   string
}

func (s *ConfigService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    s.path = "config.json"

    data, err := os.ReadFile(s.path)
    if err != nil {
        if os.IsNotExist(err) {
            s.config = &Config{} // Default config
            return nil
        }
        return err
    }

    return json.Unmarshal(data, &s.config)
}

func (s *ConfigService) ServiceShutdown() error {
    data, err := json.MarshalIndent(s.config, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(s.path, data, 0644)
}
```

### 模式 3：后台工作程序

```go
type WorkerService struct {
    cancel context.CancelFunc
}

func (s *WorkerService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    workerCtx, cancel := context.WithCancel(ctx)
    s.cancel = cancel

    go s.runWorker(workerCtx)

    return nil
}

func (s *WorkerService) ServiceShutdown() error {
    if s.cancel != nil {
        s.cancel()
    }
    return nil
}

func (s *WorkerService) runWorker(ctx context.Context) {
    ticker := time.NewTicker(time.Minute)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            s.doWork()
        case <-ctx.Done():
            return
        }
    }
}
```

## 生命周期参考

| 钩子/接口 | 调用时机 | 可取消？ | 用途 |
| --- | --- | --- | --- |
| `ServiceStartup` | 在`app.Run()`期间、事件循环开始前 | 否（返回错误可中止） | 初始化 |
| `ServiceShutdown` | 关闭期间，在`OnShutdown`之后 | 否 | 清理 |
| `OnShutdown` | 确认退出时 | 否 | 应用清理 |
| `ShouldQuit` | 请求退出时 | 是（返回 false） | 确认退出 |
| `RegisterHook(WindowClosing)` | 请求关闭窗口时 | 是（`e.Cancel()`） | 阻止窗口关闭 |
| `OnWindowEvent` | 事件发生时 | 否 | 响应事件 |
| `OnApplicationEvent` | 事件发生时 | 否 | 响应事件 |

## 平台差异

### macOS

- 即使没有窗口，<strong>应用程序菜单</strong>也会保留
- <strong>Cmd+Q</strong>会触发退出（经过`ShouldQuit`）
- 除非隐藏，否则<strong>程序坞图标</strong>会一直保留
- 使用`ApplicationShouldTerminateAfterLastWindowClosed`控制退出行为

### Windows

- 没有窗口时<strong>不会显示应用程序菜单</strong>
- <strong>Alt+F4</strong>会关闭窗口（可使用`RegisterHook`阻止）
- <strong>系统托盘</strong>可使应用保持运行

### Linux

- **行为会因桌面环境而异**
- **通常与 Windows 类似**

## 调试生命周期问题

### 问题：应用无法退出

**原因：**

1. `ShouldQuit`返回`false`
2. `OnShutdown`耗时过长
3. 后台 goroutine 未停止

**解决方案：**

```go
// 1. Check ShouldQuit logic
ShouldQuit: func() bool {
    log.Println("ShouldQuit called")
    return true
}

// 2. Keep OnShutdown fast
OnShutdown: func() {
    log.Println("OnShutdown started")
    // Fast cleanup only
    log.Println("OnShutdown finished")
}

// 3. Use context for background tasks
func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    go func() {
        <-ctx.Done()
        log.Println("Context cancelled, stopping background work")
    }()
    return nil
}
```

### 问题：服务启动失败

<strong>解决方案：</strong>返回描述清晰的错误：

```go
func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    if err := s.init(); err != nil {
        return fmt.Errorf("failed to initialise: %w", err)
    }
    return nil
}
```

系统将记录该错误，并且应用不会启动。

## 最佳实践

### 应该做

- **使用服务管理生命周期**——服务提供适当的启动和关闭钩子
- **快速完成关闭**——争取在1秒内完成所有清理工作
- **使用 context 取消操作**——正确停止后台任务
- **处理启动过程中的错误**——返回错误以彻底中止启动
- **记录生命周期事件**——有助于调试

### 不应该做

- **不要在服务启动期间阻塞**——快速完成初始化（控制在2秒内）
- **不要在关闭期间显示对话框**——应用正在退出，UI 可能无法正常工作
- **不要忽略 context**——始终在 goroutine 中检查`ctx.Done()`
- **不要泄漏资源**——始终实现`ServiceShutdown`

## 后续步骤

**服务**——进一步了解服务系统 [了解更多 →](/features/bindings/services/)

**事件系统**——使用事件进行通信 [了解更多 →](/features/events/system/)

**窗口管理**——创建和管理窗口 [了解更多 →](/features/windows/basics/)

---

<strong>有生命周期方面的问题？</strong>请在[Discord](https://discord.gg/JDdSxwjhGf)中提问或查看[示例](https://github.com/wailsapp/wails/tree/master/v3/examples)。
