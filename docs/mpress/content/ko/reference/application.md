---
title: "애플리케이션 API"
description: "Application API 전체 참조"
slug: "reference/application"
sourcePath: "reference/application.md"
---

## 개요

`Application`은 Wails 앱의 핵심입니다. 창, 서비스, 이벤트를 관리하며 모든 플랫폼 기능에 대한 접근을 제공합니다.

## 애플리케이션 만들기

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

## 핵심 메서드

### Run()

애플리케이션 이벤트 루프를 시작합니다.

```go
func (a *App) Run() error
```

**예:**

```go
err := app.Run()
if err != nil {
    log.Fatal(err)
}
```

**반환값:** 시작에 실패한 경우 오류

### Quit()

애플리케이션을 정상적으로 종료합니다.

```go
func (a *App) Quit()
```

**예:**

```go
// In a menu handler
menu.Add("Quit").OnClick(func(ctx *application.Context) {
    app.Quit()
})
```

### Config()

애플리케이션 구성을 반환합니다.

```go
func (a *App) Config() Options
```

**예:**

```go
config := app.Config()
fmt.Println("App name:", config.Name)
```

## 창 관리

### app.Window.New()

기본 옵션으로 새 웹뷰 창을 만듭니다.

```go
func (wm *WindowManager) New() *WebviewWindow
```

**예:**

```go
window := app.Window.New()
window.Show()
```

### app.Window.NewWithOptions()

사용자 지정 옵션으로 새 웹뷰 창을 만듭니다.

```go
func (wm *WindowManager) NewWithOptions(options WebviewWindowOptions) *WebviewWindow
```

**예:**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "My Window",
    Width:  800,
    Height: 600,
    BackgroundColour: application.NewRGB(255, 255, 255),
})
```

### app.Window.GetByName()

이름으로 창을 가져옵니다. 창과 함께 창을 찾았는지 여부를 반환합니다.

```go
func (wm *WindowManager) GetByName(name string) (Window, bool)
```

**예:**

```go
if window, ok := app.Window.GetByName("main"); ok {
    window.Show()
}
```

### app.Window.GetAll()

애플리케이션의 모든 창을 반환합니다.

```go
func (wm *WindowManager) GetAll() []Window
```

**예:**

```go
windows := app.Window.GetAll()
for _, window := range windows {
    fmt.Println("Window:", window.Name())
}
```

## 관리자

Application은 속성을 통해 다양한 관리자에 대한 접근을 제공합니다.

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

### 사용 예

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

## 서비스 관리

### RegisterService()

애플리케이션에 서비스를 등록합니다.

```go
func (a *App) RegisterService(service Service)
```

`RegisterService`은 아무것도 반환하지 않습니다. 서비스 초기화 오류는 `app.Run()` 중 발생하는 `ServiceStartup` 실패를 통해 드러납니다.

**예:**

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

## 이벤트 관리

### app.Event.Emit()

사용자 지정 이벤트를 발생시킵니다. 훅이 이벤트 발생을 취소하면 `true`을 반환합니다.

```go
func (em *EventManager) Emit(name string, data ...any) bool
```

**예:**

```go
// Emit event with data
app.Event.Emit("user-logged-in", map[string]interface{}{
    "username": "john",
    "timestamp": time.Now(),
})
```

### app.Event.On()

사용자 지정 이벤트를 수신합니다. 구독을 해제하는 `func()`을 반환합니다.

```go
func (em *EventManager) On(name string, callback func(*CustomEvent)) func()
```

**예:**

```go
app.Event.On("user-logged-in", func(e *application.CustomEvent) {
    data := e.Data.(map[string]interface{})
    username := data["username"].(string)
    fmt.Println("User logged in:", username)
})
```

### app.Event.OnApplicationEvent()

애플리케이션 수명 주기 이벤트를 수신합니다. `eventType` 매개변수는 `events` 패키지의 `events.ApplicationEventType`입니다.

```go
func (em *EventManager) OnApplicationEvent(
    eventType events.ApplicationEventType,
    callback func(*ApplicationEvent),
) func()
```

**예:**

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

## 대화 상자 메서드

대화 상자는 `app.Dialog` 관리자를 통해 접근합니다. 전체 참조는 [대화 상자 API](/reference/dialogs/)를 참조하세요.

### 메시지 대화 상자

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

### 질문 대화 상자

질문 대화 상자는 버튼 콜백을 사용하여 사용자 응답을 처리합니다.

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

### 파일 대화 상자

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

## 로거

애플리케이션은 구조화된 로거를 제공합니다.

```go
app.Logger.Info("Message", "key", "value")
app.Logger.Error("Error occurred", "error", err)
app.Logger.Debug("Debug info")
app.Logger.Warn("Warning message")
```

**예:**

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

## 원시 메시지 처리

프런트엔드에서 백엔드로의 통신을 저수준에서 직접 제어해야 하는 애플리케이션을 위해 Wails는 `RawMessageHandler` 옵션을 제공합니다. 이 옵션은 표준 바인딩 시스템을 우회합니다.

@note{type="info"}
원시 메시지는 최후의 수단으로만 사용해야 합니다. 표준 바인딩 시스템은 고도로 최적화되어 있으며 거의 모든 애플리케이션에 충분합니다. 애플리케이션을 프로파일링하여 바인딩이 병목 지점임을 확인한 경우에만 원시 메시지를 사용하세요.

@end

### RawMessageHandler

`RawMessageHandler`은 메서드가 아니라 `application.Options`의 필드입니다. 런타임은 프런트엔드에서 `System.invoke()`를 통해 전송된 모든 원시 메시지에 대해 이 필드를 호출합니다.

```go
type Options struct {
    // ... other fields ...
    RawMessageHandler func(window Window, message string, originInfo *OriginInfo)
}
```

`OriginInfo`에는 `Origin`, `TopOrigin` 및 `IsMainFrame`이 포함됩니다. 플랫폼마다 채우는 하위 집합이 다릅니다. 플랫폼별 매트릭스는 [원시 메시지 가이드](/guides/raw-messages/)를 참조하세요.

**예:**

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

자세한 내용은 [원시 메시지 가이드](/guides/raw-messages/)를 참조하세요.

## 플랫폼별 옵션

### Windows 옵션

애플리케이션 수준에서 Windows별 동작을 구성합니다.

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

**브라우저 플래그:**

- `EnabledFeatures` - 활성화할 WebView2 기능 플래그
- `DisabledFeatures` - 비활성화할 WebView2 기능 플래그
- `AdditionalBrowserArgs` - Chromium 명령줄 인수

자세한 문서는 [창 옵션 - 애플리케이션 수준 Windows 옵션](/features/windows/options/#application-level-windows-options)을 참조하세요.

### Mac 옵션

```go
app := application.New(application.Options{
    Name: "My App",
    Mac: application.MacOptions{
        ActivationPolicy: application.ActivationPolicyRegular,
        ApplicationShouldTerminateAfterLastWindowClosed: true,
    },
})
```

### Linux 옵션

```go
app := application.New(application.Options{
    Name: "My App",
    Linux: application.LinuxOptions{
        ProgramName:                   "my-app",
        DisableQuitOnLastWindowClosed: false,
    },
})
```

## 전체 애플리케이션 예제

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
