---
title: "アプリケーション API"
description: "アプリケーション API の完全なリファレンス"
slug: "reference/application"
sourcePath: "reference/application.md"
---

## 概要

`Application`は Wails アプリの中核です。ウィンドウ、サービス、イベントを管理し、すべてのプラットフォーム機能へのアクセスを提供します。

## アプリケーションの作成

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

## コアメソッド

### Run()

アプリケーションのイベントループを開始します。

```go
func (a *App) Run() error
```

**例：**

```go
err := app.Run()
if err != nil {
    log.Fatal(err)
}
```

<strong>戻り値：</strong>起動に失敗した場合はエラー

### Quit()

アプリケーションを正常に終了します。

```go
func (a *App) Quit()
```

**例：**

```go
// In a menu handler
menu.Add("Quit").OnClick(func(ctx *application.Context) {
    app.Quit()
})
```

### Config()

アプリケーションの設定を返します。

```go
func (a *App) Config() Options
```

**例：**

```go
config := app.Config()
fmt.Println("App name:", config.Name)
```

## ウィンドウ管理

### app.Window.New()

デフォルトのオプションで新しい WebView ウィンドウを作成します。

```go
func (wm *WindowManager) New() *WebviewWindow
```

**例：**

```go
window := app.Window.New()
window.Show()
```

### app.Window.NewWithOptions()

カスタムオプションで新しい WebView ウィンドウを作成します。

```go
func (wm *WindowManager) NewWithOptions(options WebviewWindowOptions) *WebviewWindow
```

**例：**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "My Window",
    Width:  800,
    Height: 600,
    BackgroundColour: application.NewRGB(255, 255, 255),
})
```

### app.Window.GetByName()

名前でウィンドウを取得します。ウィンドウと、ウィンドウが見つかったかどうかを返します。

```go
func (wm *WindowManager) GetByName(name string) (Window, bool)
```

**例：**

```go
if window, ok := app.Window.GetByName("main"); ok {
    window.Show()
}
```

### app.Window.GetAll()

アプリケーションのすべてのウィンドウを返します。

```go
func (wm *WindowManager) GetAll() []Window
```

**例：**

```go
windows := app.Window.GetAll()
for _, window := range windows {
    fmt.Println("Window:", window.Name())
}
```

## マネージャー

Application では、プロパティを通じて各種マネージャーにアクセスできます。

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

### 使用例

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

## サービス管理

### RegisterService()

アプリケーションにサービスを登録します。

```go
func (a *App) RegisterService(service Service)
```

`RegisterService`は何も返しません。サービスの初期化エラーは、`app.Run()`中に発生した`ServiceStartup`の失敗として表面化します。

**例：**

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

## イベント管理

### app.Event.Emit()

カスタムイベントを発行します。フックによってイベントの発行がキャンセルされた場合は、`true`を返します。

```go
func (em *EventManager) Emit(name string, data ...any) bool
```

**例：**

```go
// Emit event with data
app.Event.Emit("user-logged-in", map[string]interface{}{
    "username": "john",
    "timestamp": time.Now(),
})
```

### app.Event.On()

カスタムイベントをリッスンします。購読を解除するための`func()`を返します。

```go
func (em *EventManager) On(name string, callback func(*CustomEvent)) func()
```

**例：**

```go
app.Event.On("user-logged-in", func(e *application.CustomEvent) {
    data := e.Data.(map[string]interface{})
    username := data["username"].(string)
    fmt.Println("User logged in:", username)
})
```

### app.Event.OnApplicationEvent()

アプリケーションのライフサイクルイベントをリッスンします。`eventType`パラメーターは、`events`パッケージの`events.ApplicationEventType`です。

```go
func (em *EventManager) OnApplicationEvent(
    eventType events.ApplicationEventType,
    callback func(*ApplicationEvent),
) func()
```

**例：**

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

## ダイアログメソッド

ダイアログには`app.Dialog`マネージャーを通じてアクセスします。完全なリファレンスについては、[ダイアログ API](/reference/dialogs/)を参照してください。

### メッセージダイアログ

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

### 質問ダイアログ

質問ダイアログでは、ユーザーの応答を処理するためにボタンのコールバックを使用します。

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

### ファイルダイアログ

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

## ロガー

アプリケーションは構造化ロガーを提供します。

```go
app.Logger.Info("Message", "key", "value")
app.Logger.Error("Error occurred", "error", err)
app.Logger.Debug("Debug info")
app.Logger.Warn("Warning message")
```

**例：**

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

## 生メッセージの処理

フロントエンドからバックエンドへの通信を低レベルで直接制御する必要があるアプリケーション向けに、Wails は`RawMessageHandler`オプションを提供します。このオプションを使用すると、標準のバインディングシステムを迂回します。

@note{type="info"}
生メッセージは、最後の手段としてのみ使用してください。標準のバインディングシステムは高度に最適化されており、ほぼすべてのアプリケーションに十分です。アプリケーションをプロファイリングし、バインディングがボトルネックであることを確認した場合にのみ、生メッセージを使用してください。

@end

### RawMessageHandler

`RawMessageHandler`はメソッドではなく、`application.Options`のフィールドです。フロントエンドから`System.invoke()`を介して生メッセージが送信されるたびに、ランタイムがこのフィールドを呼び出します。

```go
type Options struct {
    // ... other fields ...
    RawMessageHandler func(window Window, message string, originInfo *OriginInfo)
}
```

`OriginInfo`には、`Origin`、`TopOrigin`、`IsMainFrame`が含まれます。プラットフォームごとに設定される項目の組み合わせが異なります。プラットフォーム別の対応表については、[生メッセージガイド](/guides/raw-messages/)を参照してください。

**例：**

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

詳細については、[生メッセージガイド](/guides/raw-messages/)を参照してください。

## プラットフォーム固有のオプション

### Windows オプション

アプリケーションレベルで Windows 固有の動作を設定します。

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

**ブラウザーフラグ：**

- `EnabledFeatures` - 有効にする WebView2 機能フラグ
- `DisabledFeatures` - 無効にする WebView2 機能フラグ
- `AdditionalBrowserArgs` - Chromium のコマンドライン引数

詳細については、[ウィンドウオプション - アプリケーションレベルの Windows オプション](/features/windows/options/#application-level-windows-options)を参照してください。

### Mac オプション

```go
app := application.New(application.Options{
    Name: "My App",
    Mac: application.MacOptions{
        ActivationPolicy: application.ActivationPolicyRegular,
        ApplicationShouldTerminateAfterLastWindowClosed: true,
    },
})
```

### Linux オプション

```go
app := application.New(application.Options{
    Name: "My App",
    Linux: application.LinuxOptions{
        ProgramName:                   "my-app",
        DisableQuitOnLastWindowClosed: false,
    },
})
```

## 完全なアプリケーションの例

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
