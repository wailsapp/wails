---
title: "應用程式生命週期"
description: "瞭解 Wails 應用程式從啟動到關閉的生命週期"
slug: "concepts/lifecycle"
sourcePath: "concepts/lifecycle.md"
---

## 瞭解應用程式生命週期

桌面應用程式具有從啟動到關閉的生命週期。Wails v3 提供<strong>服務</strong>、<strong>事件</strong>和<strong>鉤子</strong>，以有效管理此生命週期。

## 生命週期階段

```d2
direction: down

Start: 應用程式啟動 {
  shape: oval
  style.fill: "#10B981"
}

Init: 初始化 {
  Parse: 解析選項 {
    shape: rectangle
  }
  Register: 註冊服務 {
    shape: rectangle
  }
  Setup: 設定執行階段 {
    shape: rectangle
  }
}

AppRun: app.Run() {
  shape: rectangle
  style.fill: "#3B82F6"
}

ServiceStartup: 啟動服務 {
  shape: rectangle
  style.fill: "#8B5CF6"
}

EventLoop: 事件迴圈 {
  Process: 處理事件 {
    shape: rectangle
  }
  Handle: 處理訊息 {
    shape: rectangle
  }
  Update: 更新使用者介面 {
    shape: rectangle
  }
}

QuitSignal: 結束訊號 {
  shape: diamond
  style.fill: "#F59E0B"
}

ShouldQuit: ShouldQuit 檢查 {
  shape: rectangle
  style.fill: "#3B82F6"
}

OnShutdown: OnShutdown 回呼 {
  shape: rectangle
  style.fill: "#3B82F6"
}

ServiceShutdown: 關閉服務 {
  shape: rectangle
  style.fill: "#8B5CF6"
}

Cleanup: 清理 {
  Close: 關閉視窗 {
    shape: rectangle
  }
  Release: 釋放資源 {
    shape: rectangle
  }
}

End: 應用程式結束 {
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
EventLoop.Update -> EventLoop.Process: 迴圈
EventLoop.Process -> QuitSignal: 使用者結束應用程式
QuitSignal -> ShouldQuit: 檢查是否允許？
ShouldQuit -> EventLoop.Process: 已拒絕
ShouldQuit -> OnShutdown: 已允許
OnShutdown -> ServiceShutdown
ServiceShutdown -> Cleanup.Close
Cleanup.Close -> Cleanup.Release
Cleanup.Release -> End
```

### 1. 建立應用程式

使用`application.New()`建立應用程式：

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

**執行的作業：**

1. 剖析並驗證選項
2. 註冊服務（但尚未啟動）
3. 設定資產伺服器
4. 設定執行階段

### 2. 執行應用程式

呼叫`app.Run()`以啟動應用程式：

```go
err := app.Run()  // Blocks until quit
if err != nil {
    log.Fatal(err)
}
```

**執行的作業：**

1. 依註冊順序啟動服務
2. 啟用事件監聽器
3. 此時可以建立視窗
4. 事件迴圈開始執行

### 3. 事件迴圈

應用程式會進入事件迴圈，並在其中度過大部分執行時間：

- 處理作業系統事件（滑鼠、鍵盤和視窗事件）
- 處理 Go 傳送至 JS 的訊息
- 執行 JS 對 Go 的呼叫
- 轉譯 UI 更新

### 4. 關閉

應用程式結束時：

1. 檢查`ShouldQuit`回呼（若已設定）
2. 執行`OnShutdown`回呼
3. 依相反順序關閉服務
4. 關閉視窗
5. 釋放資源

## 服務生命週期

服務是 Wails v3 中管理生命週期的主要方式。服務透過介面提供啟動和關閉鉤子。如需完整的服務文件，請參閱[服務指南](/features/bindings/services/)。

### 建立服務

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

### 註冊服務

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&MyService{}),
        application.NewService(&AnotherService{}),
    },
})
```

**重點：**

- 服務依註冊順序啟動
- 服務依註冊的<strong>相反</strong>順序關閉
- 如果服務的`ServiceStartup`傳回錯誤，應用程式便會中止
- 關閉開始時，傳遞給`ServiceStartup`的`ctx`會被取消

### 使用應用程式上下文

傳遞給`ServiceStartup`的 context 在應用程式的整個生命週期內都有效：

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

您也可以從應用程式實例存取 context：

```go
app := application.Get()
ctx := app.Context()
```

## 應用程式層級的 Hook

這些是`application.Options`中的便利回呼，可讓您掛接應用程式生命週期，而不必建立完整的服務。它們適用於簡單的清理工作、結束確認，或需要在關閉程序中的特定時點執行程式碼時。

若需要包含啟動邏輯、相依性注入或具狀態資源的較複雜生命週期管理，請改用[服務](#heading-2)。

### ShouldQuit

每當收到結束要求時，都會呼叫`ShouldQuit`回呼，無論是使用者關閉最後一個視窗、按下 Cmd+Q（macOS）／Alt+F4（Windows），還是以程式方式呼叫`app.Quit()`。

**傳回值：**

- 傳回`true`以允許繼續結束（應用程式將會關閉）
- 傳回`false`以取消結束（應用程式會繼續執行）

您可以在此攔截結束要求，並選擇是否阻止結束，例如提示使用者尚有未儲存的變更：

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

如果未設定`ShouldQuit`，應用程式會在收到要求時立即結束。

**呼叫 ShouldQuit 的時機：**

- 使用者關閉最後一個視窗（除非已設定`DisableQuitOnLastWindowClosed`）
- 使用者在 macOS 上按下 Cmd+Q
- 使用者在 Windows 上按下 Alt+F4（焦點位於最後一個視窗時）
- 程式碼呼叫`app.Quit()`

**不會呼叫 ShouldQuit 的情況：**

- 處理程序遭到終止（SIGKILL、透過工作管理員強制結束）
- 直接呼叫`os.Exit()`

### OnShutdown

確認應用程式即將結束時，會呼叫`OnShutdown`回呼（若已設定`ShouldQuit`，則會在其傳回`true`後呼叫）。請用它執行清理工作，例如儲存狀態、關閉資料庫連線或釋放資源。

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

您也可以在應用程式生命週期中的任何時點，以程式方式註冊其他關閉回呼：

```go
app.OnShutdown(func() {
    log.Println("Application shutting down...")
})
```

多個回呼會依註冊順序執行。關閉程序會封鎖，直到所有回呼執行完畢。

<strong>重要：</strong>關閉回呼應快速完成（少於1秒）。作業系統可能會強制終止耗時過久仍未結束的應用程式，因而中斷清理工作並造成資料遺失。

### PostShutdown

所有關閉工作完成後、處理程序終止前，會呼叫`PostShutdown`回呼。此時已無法再使用應用程式實例——視窗已關閉、服務已停止，資源也已釋放。

這主要適用於：

- 必須在所有其他清理工作之後進行的最終記錄
- 測試及偵錯關閉行為
- `app.Run()`不會傳回的平台（此回呼可確保您的程式碼執行）

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

<strong>注意：</strong>請勿嘗試在`PostShutdown`中使用應用程式功能（視窗、對話方塊等），因為這些功能已無法使用。

## 事件式生命週期

Wails 提供事件系統，在應用程式中發生各種情況時通知您，例如視窗開啟、應用程式啟動、佈景主題變更等。您可以監聽這些事件以回應生命週期變更，而不會封鎖或攔截它們。

對於視窗事件，需要攔截並取消動作（例如防止視窗關閉）時，也可以使用`RegisterHook`取代`OnWindowEvent`。請參閱下方的[視窗 Hook](#-hook-1)。

如需事件系統的完整文件，請參閱[事件指南](/features/events/system/)。

### 應用程式事件

監聽應用程式生命週期事件：

```go
app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(event *application.ApplicationEvent) {
    app.Logger.Info("Application has started!")
})
```

此外也有平台特定事件可供使用：

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

### 視窗事件

監聽視窗生命週期事件：

```go
window := app.Window.New()

window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
    app.Logger.Info("Window gained focus")
})

window.OnWindowEvent(events.Common.WindowClosing, func(e *application.WindowEvent) {
    app.Logger.Info("Window is closing")
})
```

### 視窗 Hook（可取消的事件）

需要<strong>取消</strong>事件時，請使用`RegisterHook`，而非`OnWindowEvent`：

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

**OnWindowEvent 與 RegisterHook 的差異：**

- `OnWindowEvent`：在事件發生時通知您（無法取消）
- `RegisterHook`：讓您攔截事件，並視情況取消事件

## 視窗生命週期

視窗有自己的生命週期，從建立到銷毀。每個視窗都會獨立載入其前端內容，並可隨時顯示、隱藏或關閉。當使用者嘗試關閉視窗時，您可以使用`RegisterHook`攔截此操作，以提示使用者確認，或隱藏視窗而不將其銷毀。

如需完整的視窗說明文件，請參閱[視窗指南](/features/windows/basics/)。

```d2
direction: down

Create: 建立視窗 {
  shape: oval
  style.fill: "#10B981"
}

Load: 載入前端 {
  shape: rectangle
}

Show: 顯示視窗 {
  shape: rectangle
}

Active: 視窗作用中 {
  Events: 處理事件 {
    shape: rectangle
  }
}

CloseRequest: 關閉要求 {
  shape: diamond
  style.fill: "#F59E0B"
}

Hook: WindowClosing 掛鉤 {
  shape: rectangle
  style.fill: "#3B82F6"
}

Destroy: 銷毀視窗 {
  shape: rectangle
}

End: 視窗已關閉 {
  shape: oval
  style.fill: "#EF4444"
}

Create -> Load
Load -> Show
Show -> Active.Events
Active.Events -> Active.Events: 迴圈
Active.Events -> CloseRequest: 使用者關閉視窗
CloseRequest -> Hook
Hook -> Active.Events: 已取消
Hook -> Destroy: 已允許
Destroy -> End
```

### 建立視窗

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "My Window",
    Width:  800,
    Height: 600,
})
```

### 防止視窗關閉

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

### 隱藏而非關閉

系統匣應用程式的常見模式：

```go
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    window.Hide()  // Hide instead of destroy
    e.Cancel()     // Prevent actual close
})
```

## 多視窗生命週期

使用多個視窗時：

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

**預設行為因平台而異：**

| 平台 | 最後一個視窗關閉時的預設行為 |
| --- | --- |
| macOS | 應用程式繼續執行（選單列仍保留） |
| Windows | 應用程式結束 |
| Linux | 應用程式結束 |

macOS 遵循原生平台慣例，即使沒有任何視窗，應用程式通常仍會在選單列中保持作用中。Windows 和 Linux 則預設會結束應用程式。

**讓所有平台都在最後一個視窗關閉時結束應用程式：**

```go
app := application.New(application.Options{
    Mac: application.MacOptions{
        ApplicationShouldTerminateAfterLastWindowClosed: true,
    },
})
```

**讓所有平台都在最後一個視窗關閉時繼續執行應用程式：**

這適用於系統匣應用程式，或應在背景持續執行的應用程式。

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

## 常見模式

### 模式1：資料庫服務

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

### 模式2：組態服務

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

### 模式3：背景工作程式

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

## 生命週期參考

| 掛鉤／介面 | 呼叫時機 | 可否取消？ | 用途 |
| --- | --- | --- | --- |
| `ServiceStartup` | 在`app.Run()`期間、事件迴圈開始前 | 否（傳回錯誤以中止） | 初始化 |
| `ServiceShutdown` | 關閉期間，在`OnShutdown`之後 | 否 | 清理 |
| `OnShutdown` | 確認結束時 | 否 | 應用程式清理 |
| `ShouldQuit` | 要求結束時 | 是（傳回 false） | 確認結束 |
| `RegisterHook(WindowClosing)` | 要求關閉視窗時 | 是（`e.Cancel()`） | 防止視窗關閉 |
| `OnWindowEvent` | 事件發生時 | 否 | 回應事件 |
| `OnApplicationEvent` | 事件發生時 | 否 | 回應事件 |

## 平台差異

### macOS

- 即使沒有視窗，<strong>應用程式選單</strong>仍會保留
- <strong>Cmd+Q</strong>會觸發結束（會經過`ShouldQuit`）
- 除非隱藏，否則<strong>Dock 圖示</strong>會保留
- 使用`ApplicationShouldTerminateAfterLastWindowClosed`控制結束行為

### Windows

- 沒有視窗時，**不會有應用程式選單**
- <strong>Alt+F4</strong>會關閉視窗（可使用`RegisterHook`阻止）
- <strong>系統匣</strong>可讓應用程式繼續執行

### Linux

- **行為會因桌面環境而異**
- **通常與 Windows 類似**

## 偵錯生命週期問題

### 問題：應用程式無法結束

**原因：**

1. `ShouldQuit`傳回`false`
2. `OnShutdown`耗時過久
3. 背景 goroutine 未停止

**解決方案：**

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

### 問題：服務啟動失敗

<strong>解決方案：</strong>傳回描述清楚的錯誤：

```go
func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    if err := s.init(); err != nil {
        return fmt.Errorf("failed to initialise: %w", err)
    }
    return nil
}
```

系統會記錄錯誤，且應用程式不會啟動。

## 最佳實務

### 建議做法

- **使用服務管理生命週期** — 服務提供適當的啟動與關閉掛鉤
- **快速完成關閉** — 所有清理工作的目標時間應少於1秒
- **使用 context 進行取消** — 正確停止背景工作
- **處理啟動期間的錯誤** — 傳回錯誤以乾淨地中止啟動
- **記錄生命週期事件** — 有助於偵錯

### 避免事項

- **不要在服務啟動期間阻塞** — 快速完成初始化（少於2秒）
- **不要在關閉期間顯示對話方塊** — 應用程式正在結束，使用者介面可能無法運作
- **不要忽略 context** — 一律在 goroutine 中檢查`ctx.Done()`
- **不要洩漏資源** — 一律實作`ServiceShutdown`

## 後續步驟

**服務** — 進一步瞭解服務系統 [深入瞭解 →](/features/bindings/services/)

**事件系統** — 使用事件進行通訊 [深入瞭解 →](/features/events/system/)

**視窗管理** — 建立及管理視窗 [深入瞭解 →](/features/windows/basics/)

---

<strong>對生命週期有疑問嗎？</strong>請在[Discord](https://discord.gg/JDdSxwjhGf)中提問，或查看[範例](https://github.com/wailsapp/wails/tree/master/v3/examples)。
