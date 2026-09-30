---
title: "관리자 API"
description: "특정 기능에 집중하는 manager 인터페이스로 구성된 체계적인 API 구조"
slug: "concepts/manager-api"
sourcePath: "concepts/manager-api.md"
---

Wails v3 관리자 API는 `*application.App`의 공개 필드 아래에 그룹화된, 특정 기능에 집중하는 manager 구조체를 통해 애플리케이션 기능에 체계적이고 쉽게 탐색할 수 있는 방식으로 접근하도록 합니다. Wails 3은 v2와의 호환성을 완전히 단절했습니다. 기존 `app.NewWebviewWindow(...)` 스타일 API와의 호환성을 유지하기 위한 호출별 래퍼 계층이 없으므로, 아래 manager를 사용하여 앱을 제어해야 합니다.

## 개요

관리자 API는 애플리케이션 기능을 열두 가지 영역(로거 하나와 manager 열하나)으로 구성합니다.

- **`app.Window`** - 창 생성, 관리 및 콜백
- **`app.ContextMenu`** - 컨텍스트 메뉴 등록 및 관리\
- **`app.KeyBinding`** - 전역 키 바인딩 관리
- **`app.Browser`** - 브라우저 통합(URL 및 파일 열기)
- **`app.Env`** - 환경 정보 및 시스템 상태
- **`app.Dialog`** - 파일 및 메시지 대화 상자 작업
- **`app.Event`** - 사용자 지정 이벤트 처리 및 애플리케이션 이벤트
- **`app.Menu`** - 애플리케이션 메뉴 관리
- **`app.Screen`** - 화면 관리 및 좌표 변환
- **`app.Clipboard`** - 클립보드 텍스트 작업
- **`app.SystemTray`** - 시스템 트레이 아이콘 생성 및 관리
- **`app.Autostart`** - 사용자 로그인 시 애플리케이션이 실행되도록 등록

## 장점

- **향상된 탐색성** - IDE 자동 완성에서 체계적으로 구성된 API 표면을 표시합니다.
- **개선된 코드 구성** - 관련 메서드를 함께 그룹화합니다.
- **향상된 유지보수성** - manager 간에 관심사를 분리합니다.
- **향후 확장성** - 특정 영역에 새 기능을 더 쉽게 추가할 수 있습니다.

## 사용법

관리자 API를 사용하면 모든 애플리케이션 기능에 체계적으로 접근할 수 있습니다.

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

## Manager 참조

### 창 관리자

창 생성, 조회 및 수명 주기 콜백을 관리합니다.

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

### 이벤트 관리자

사용자 지정 이벤트와 애플리케이션 이벤트 수신을 처리합니다.

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

### 브라우저 관리자

URL과 파일을 열기 위한 브라우저 통합 기능을 제공합니다.

```go
// Open URLs and files in default browser
err := app.Browser.OpenURL("https://wails.io")
err := app.Browser.OpenFile("/path/to/document.pdf")
```

### 환경 관리자

시스템 환경 정보에 접근할 수 있습니다.

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

### 대화 상자 관리자

파일 및 메시지 대화 상자에 체계적으로 접근할 수 있습니다.

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

### 메뉴 관리자

애플리케이션 메뉴를 생성하고 관리합니다.

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

### 키 바인딩 관리자

전역 키 바인딩을 동적으로 관리합니다.

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

### 컨텍스트 메뉴 관리자

고급 컨텍스트 메뉴 관리 기능입니다(라이브러리 작성자용).

```go
// Create and register context menu
menu := app.ContextMenu.New()
app.ContextMenu.Add("myMenu", menu)

// Retrieve context menu
menu, exists := app.ContextMenu.Get("myMenu")

// Remove context menu
app.ContextMenu.Remove("myMenu")
```

### 화면 관리자

다중 모니터 구성에서 화면을 관리하고 좌표를 변환합니다.

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

### 클립보드 관리자

텍스트를 읽고 쓰는 클립보드 작업을 수행합니다.

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

### 시스템 트레이 관리자

시스템 트레이 아이콘을 생성하고 관리합니다.

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

### 자동 시작 관리자

사용자 로그인 시 애플리케이션이 실행되도록 등록합니다. 플랫폼별로 적합한 네이티브 메커니즘을 선택합니다. macOS에서는 SMAppService 또는 LaunchAgent plist, Windows에서는 `HKCU\…\Run` 레지스트리 키, Linux에서는 XDG `.desktop` 항목을 사용합니다.

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

플랫폼별 동작, 식별자 규칙 및 오래된 등록 감지 보장에 관한 자세한 내용은 [Autostart 기능 페이지](/features/autostart/basics/)를 참조하세요.
