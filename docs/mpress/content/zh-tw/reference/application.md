---
title: "應用程式 API"
description: "應用程式 API 完整參考文件"
slug: "reference/application"
sourcePath: "reference/application.md"
---

## 概觀

`Application` 是 Wails 應用程式的核心。它負責管理視窗、服務與事件，並提供所有平台功能的存取方式。

## 建立應用程式

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

啟動應用程式事件迴圈。

```go
func (a *App) Run() error
```

**範例：**

```go
err := app.Run()
if err != nil {
    log.Fatal(err)
}
```

<strong>傳回值：</strong>若啟動失敗，則傳回錯誤

### Quit()

正常關閉應用程式。

```go
func (a *App) Quit()
```

**範例：**

```go
// In a menu handler
menu.Add("Quit").OnClick(func(ctx *application.Context) {
    app.Quit()
})
```

### Config()

傳回應用程式組態。

```go
func (a *App) Config() Options
```

**範例：**

```go
config := app.Config()
fmt.Println("App name:", config.Name)
```

## 視窗管理

### app.Window.New()

使用預設選項建立新的 WebView 視窗。

```go
func (wm *WindowManager) New() *WebviewWindow
```

**範例：**

```go
window := app.Window.New()
window.Show()
```

### app.Window.NewWithOptions()

使用自訂選項建立新的 WebView 視窗。

```go
func (wm *WindowManager) NewWithOptions(options WebviewWindowOptions) *WebviewWindow
```

**範例：**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "My Window",
    Width:  800,
    Height: 600,
    BackgroundColour: application.NewRGB(255, 255, 255),
})
```

### app.Window.GetByName()

依名稱取得視窗。傳回該視窗以及是否找到該視窗。

```go
func (wm *WindowManager) GetByName(name string) (Window, bool)
```

**範例：**

```go
if window, ok := app.Window.GetByName("main"); ok {
    window.Show()
}
```

### app.Window.GetAll()

傳回應用程式的所有視窗。

```go
func (wm *WindowManager) GetAll() []Window
```

**範例：**

```go
windows := app.Window.GetAll()
for _, window := range windows {
    fmt.Println("Window:", window.Name())
}
```

## 管理器

Application 透過屬性提供各種管理器的存取方式：

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

### 使用範例

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

## 服務管理

### RegisterService()

向應用程式註冊服務。

```go
func (a *App) RegisterService(service Service)
```

`RegisterService`不傳回任何內容；服務初始化錯誤會在`app.Run()`期間透過`ServiceStartup`失敗呈現。

**範例：**

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

發出自訂事件。如果某個掛鉤取消了事件發出，則傳回`true`。

```go
func (em *EventManager) Emit(name string, data ...any) bool
```

**範例：**

```go
// Emit event with data
app.Event.Emit("user-logged-in", map[string]interface{}{
    "username": "john",
    "timestamp": time.Now(),
})
```

### app.Event.On()

監聽自訂事件。傳回用於取消訂閱的`func()`。

```go
func (em *EventManager) On(name string, callback func(*CustomEvent)) func()
```

**範例：**

```go
app.Event.On("user-logged-in", func(e *application.CustomEvent) {
    data := e.Data.(map[string]interface{})
    username := data["username"].(string)
    fmt.Println("User logged in:", username)
})
```

### app.Event.OnApplicationEvent()

監聽應用程式生命週期事件。`eventType`參數為`events.ApplicationEventType`（來自`events`套件）。

```go
func (em *EventManager) OnApplicationEvent(
    eventType events.ApplicationEventType,
    callback func(*ApplicationEvent),
) func()
```

**範例：**

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

## 對話方塊方法

請透過`app.Dialog`管理器存取對話方塊。完整參考資料請參閱[對話方塊 API](/reference/dialogs/)。

### 訊息對話方塊

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

### 詢問對話方塊

詢問對話方塊使用按鈕回呼處理使用者回應：

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

### 檔案對話方塊

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

## 記錄器

應用程式提供結構化記錄器：

```go
app.Logger.Info("Message", "key", "value")
app.Logger.Error("Error occurred", "error", err)
app.Logger.Debug("Debug info")
app.Logger.Warn("Warning message")
```

**範例：**

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

## 原始訊息處理

對於需要直接、低階控制前端到後端通訊的應用程式，Wails 提供`RawMessageHandler`選項。這會略過標準繫結系統。

@note{type="info"}
原始訊息應僅作為最後手段使用。標準繫結系統經過高度最佳化，足以滿足幾乎所有應用程式的需求。只有在分析應用程式效能並確認繫結是效能瓶頸後，才使用原始訊息。

@end

### RawMessageHandler

`RawMessageHandler`是`application.Options`上的欄位，而非方法。對於前端透過`System.invoke()`傳送的每則原始訊息，執行階段都會叫用此欄位。

```go
type Options struct {
    // ... other fields ...
    RawMessageHandler func(window Window, message string, originInfo *OriginInfo)
}
```

`OriginInfo`包含`Origin`、`TopOrigin`和`IsMainFrame`（各平台會填入不同的子集；各平台的對照矩陣請參閱[原始訊息指南](/guides/raw-messages/)）。

**範例：**

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

如需更多詳細資訊，請參閱[原始訊息指南](/guides/raw-messages/)。

## 平台特定選項

### Windows 選項

在應用程式層級設定 Windows 特定行為：

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

**瀏覽器旗標：**

- `EnabledFeatures` - 要啟用的 WebView2 功能旗標
- `DisabledFeatures` - 要停用的 WebView2 功能旗標
- `AdditionalBrowserArgs` - Chromium 命令列引數

如需詳細文件，請參閱[視窗選項 - 應用程式層級的 Windows 選項](/features/windows/options/#application-level-windows-options)。

### Mac 選項

```go
app := application.New(application.Options{
    Name: "My App",
    Mac: application.MacOptions{
        ActivationPolicy: application.ActivationPolicyRegular,
        ApplicationShouldTerminateAfterLastWindowClosed: true,
    },
})
```

### Linux 選項

```go
app := application.New(application.Options{
    Name: "My App",
    Linux: application.LinuxOptions{
        ProgramName:                   "my-app",
        DisableQuitOnLastWindowClosed: false,
    },
})
```

## 完整的應用程式範例

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
