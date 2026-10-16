---
title: "창 기본 사항"
description: "Wails에서 애플리케이션 창 만들기 및 관리"
slug: "features/windows/basics"
sourcePath: "features/windows/basics.md"
---

## 창 관리

Wails는 모든 플랫폼에서 작동하는 <strong>통합 창 관리 API</strong>를 제공합니다. 창을 만들고 동작을 제어하며, 생성, 모양, 동작, 수명 주기를 완전히 제어하면서 여러 창을 관리할 수 있습니다.

## 빠른 시작

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

**이것으로 끝입니다!** 이제 크로스 플랫폼 창이 만들어졌습니다.

## 창 만들기

### 기본 창

창을 만드는 가장 간단한 방법은 다음과 같습니다.

```go
window := app.Window.New()
```

**생성되는 항목:**

- 기본 크기(800x600)
- 기본 제목(애플리케이션 이름)
- 프런트엔드에 사용할 준비가 된 WebView
- 플랫폼 고유의 모양

### 옵션을 지정한 창

사용자 지정 구성으로 창을 만드세요.

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

**일반 옵션:**

| 옵션 | 유형 | 설명 |
| --- | --- | --- |
| `Title` | `string` | 창 제목 |
| `Width` | `int` | 창 너비(픽셀) |
| `Height` | `int` | 창 높이(픽셀) |
| `X` | `int` | X 위치(왼쪽 기준) |
| `Y` | `int` | Y 위치(위쪽 기준) |
| `AlwaysOnTop` | `bool` | 창을 다른 창보다 위에 유지 |
| `Frameless` | `bool` | 제목 표시줄과 테두리 제거 |
| `Hidden` | `bool` | 숨긴 상태로 시작 |
| `MinWidth` | `int` | 최소 너비 |
| `MinHeight` | `int` | 최소 높이 |
| `MaxWidth` | `int` | 최대 너비 |
| `MaxHeight` | `int` | 최대 높이 |

**전체 목록은 [창 옵션](/features/windows/options/)을 참조하세요.**

### 이름이 지정된 창

창을 쉽게 찾을 수 있도록 이름을 지정하세요.

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

**사용 사례:**

- 여러 창(메인, 설정, 정보)
- 코드의 여러 부분에서 창 찾기
- 창 간 통신

## 창 제어

### 표시 및 숨기기

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

**사용 사례:**

- 스플래시 화면(표시한 다음 숨기기)
- 설정 창(필요하지 않을 때 숨기기)
- 팝업 창(필요할 때 표시)

### 위치 및 크기

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

**좌표계:**

- (0, 0)은 기본 화면의 왼쪽 위 모서리입니다.
- X의 양의 방향은 오른쪽입니다.
- 양의 Y 방향은 아래쪽입니다.

### 창 상태

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

**상태 전환:**

```
Normal ←→ Minimised
Normal ←→ Maximised
Normal ←→ Fullscreen
```

### 제목 및 모양

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

### 창 닫기

```go
// Close window — dispatches WindowClosing; a RegisterHook can call e.Cancel().
window.Close()
```

v3에는 `window.Destroy()` 메서드가 없습니다. `Close()`을 사용하고, `OnWindowEvent`로 수신하거나(취소할 수 없음) `RegisterHook`로 후킹하세요(`e.Cancel()`을 호출하여 창을 열린 상태로 유지할 수 있음).

## 창 찾기

### 이름으로 찾기

```go
if window, ok := app.Window.GetByName("settings"); ok {
    window.Show()
}
```

### ID로 찾기

모든 창에는 고유한 ID가 있습니다:

```go
id := window.ID()
fmt.Printf("Window ID: %d\n", id)

// Find by ID
if found, ok := app.Window.GetByID(id); ok {
    found.Focus()
}
```

### 현재 창

현재 포커스된 창을 가져옵니다:

```go
current := app.Window.Current()
if current != nil {
    current.SetTitle("Active Window")
}
```

### 모든 창

모든 창을 가져옵니다:

```go
windows := app.Window.GetAll()
fmt.Printf("Total windows: %d\n", len(windows))

for _, w := range windows {
    fmt.Printf("Window: %s (ID: %d)\n", w.Name(), w.ID())
}
```

## 창 수명 주기

### 생성

```go
app.Window.OnCreate(func(window application.Window) {
    fmt.Printf("Window created: %s\n", window.Name())

    // Configure new windows
    window.SetMinSize(400, 300)
})
```

### 닫기

창이 닫히지 않도록 하려면 `WindowClosing` 이벤트와 함께 `RegisterHook`을 사용하세요:

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

**중요:** `RegisterHook`은 닫기 이벤트가 발생하기 전에 이를 가로챕니다. 창이 닫히지 않도록 하려면 `event.Cancel()`을 호출하세요. 이는 사용자가 시작한 닫기 동작(X 버튼 클릭)에 적용됩니다.

### 소멸

창이 닫힐 때 정리 작업을 수행하려면 `WindowClosing` 이벤트와 함께 `OnWindowEvent`을 사용하세요:

```go
window.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
    fmt.Println("Window is closing")
    // Cleanup resources
})
```

## 여러 창

### 여러 창 만들기

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

### 창 간 통신

창은 이벤트를 통해 통신할 수 있습니다:

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

**자세한 내용은 [이벤트](/features/events/system/)를 참조하세요.**

### 부모-자식 창

`WebviewWindowOptions`에는 `Parent` 필드가 없습니다. 자식 창을 일반 창으로 만든 다음 시트 모달로 부모 창에 연결하세요:

```go
// Create child window
childWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Child Window",
})

// Attach to the parent — presents as a sheet on macOS.
mainWindow.AttachModal(childWindow)
```

**동작:**

- 자식 창은 부모 창 위에 유지됩니다.
- 자식 창은 모달이므로 부모 창과의 상호 작용을 차단합니다.

**플랫폼 지원:**

- **macOS:** 완전히 지원됩니다(시트로 표시됨).
- **Windows:** 지원되지 않습니다.
- **Linux:** 지원되지 않습니다.

## 플랫폼별 기능

@tabs{sync-key="platform"}
[Windows]
**Windows 전용 기능:**

```go
// Flash taskbar button
window.Flash(true)  // Start flashing
window.Flash(false) // Stop flashing

// Trigger Windows 11 Snap Assist (Win+Z)
window.SnapAssist()
```

창별 `SetIcon`은 없습니다. 애플리케이션 아이콘은 `app.SetIcon([]byte)`을 통해 앱에 설정합니다. Linux 전용 창 아이콘의 경우에는 창을 생성할 때 `application.LinuxWindow.Icon` 필드를 사용합니다.

**Snap Assist:** 시스템 바로 가기 경로를 통해 Windows 11 스냅 레이아웃 옵션을 표시합니다. 마우스를 올렸을 때 네이티브 Snap Layouts를 표시하는 사용자 지정 HTML 최대화 버튼을 사용하려면 대신 [Windows의 네이티브 비클라이언트 영역](/features/windows/frameless/#native-non-client-regions-on-windows)을 사용하세요.

**작업 표시줄 깜박임:** 창이 최소화되어 있을 때 알림을 표시하는 데 유용합니다.

[macOS]
**macOS 전용 기능:**

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

**배경막 유형:**

- `MacBackdropNormal` - 표준 창
- `MacBackdropTranslucent` - 반투명 배경. webview를 투명하게 하려면 **비공개 API가 필요합니다**.
- `MacBackdropTransparent` - 완전히 투명한 배경. webview를 투명하게 하려면 **비공개 API가 필요합니다**.
- `MacBackdropLiquidGlass` - 유리 효과 배경막. webview를 투명하게 하려면 **비공개 API가 필요합니다**.

이러한 효과가 webview를 통해 보이도록 하려면 `-tags private_mac_apis`을 사용하여 빌드하세요. 이를 사용하지 않으면 webview가 불투명하게 유지됩니다. `TitleBar.AppearsTransparent` 자체는 공개 API를 사용합니다. [비공개 macOS API](/guides/build/private-macos-apis/)를 참조하세요.

**컬렉션 동작:** 여러 Space에서 창이 동작하는 방식을 제어합니다:

- `MacWindowCollectionBehaviorCanJoinAllSpaces` - 모든 Space에 표시
- `MacWindowCollectionBehaviorFullScreenAuxiliary` - 전체 화면 앱 위에 표시 가능

**네이티브 전체 화면:** macOS의 전체 화면 모드는 새 Space(가상 데스크톱)를 만듭니다.

[Linux]
**Linux 전용 기능:**

```go
// Set window icon (per-window struct is LinuxWindow, not the app-level LinuxOptions)
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Linux: application.LinuxWindow{
        Icon: iconBytes,
    },
})
```

**데스크톱 환경 참고 사항:**

- GNOME: 완전히 지원됨
- KDE Plasma: 완전히 지원됨
- XFCE: 부분적으로 지원됨
- 기타: 환경에 따라 다름

**타일링 창 관리자(Hyprland, Sway, i3 등):**

- `Minimise()` 및 `Maximise()`은 예상대로 작동하지 않을 수 있습니다. 창 관리자가 창의 위치와 크기를 제어하기 때문입니다.
- `SetSize()` 및 `SetPosition()` 요청은 권고 사항이므로 무시될 수 있습니다
- `Fullscreen()`은 일반적으로 예상대로 작동합니다
- 일부 WM은 항상 위에 표시 기능을 지원하지 않습니다

@end

## 일반적인 패턴

### 시작 화면

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

### 설정 창

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

### 닫기 전 확인

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

## 권장 사항

### ✅ 권장

- **중요한 창에 이름을 지정하세요** - 나중에 더 쉽게 찾을 수 있습니다
- **최소 크기를 설정하세요** - 사용할 수 없는 레이아웃을 방지합니다
- **창을 가운데에 배치하세요** - 임의의 위치보다 사용자 경험이 좋습니다
- **닫기 이벤트를 처리하세요** - 데이터 손실을 방지합니다
- **모든 플랫폼에서 테스트하세요** - 동작이 플랫폼마다 다릅니다
- **적절한 크기를 사용하세요** - 다양한 화면 크기를 고려하세요

### ❌ 비권장

- **창을 너무 많이 만들지 마세요** - 사용자가 혼란스러울 수 있습니다
- **창을 닫는 것을 잊지 마세요** - 메모리 누수가 발생합니다
- **위치를 하드코딩하지 마세요** - 화면 크기가 서로 다릅니다
- **플랫폼 간 차이를 무시하지 마세요** - 철저히 테스트하세요
- **UI 스레드를 차단하지 마세요** - 오래 걸리는 작업에는 goroutine을 사용하세요

## 문제 해결

### 창이 표시되지 않음

**가능한 원인:**

1. 창이 숨김 상태로 생성됨
2. 창이 화면 밖에 있음
3. 창이 다른 창 뒤에 있음

**해결 방법:**

```go
window.Show()
window.Center()
window.Focus()
```

### 창 크기가 잘못됨

**원인:** Windows/Linux의 DPI 배율 조정

**해결 방법:**

```go
// Wails handles DPI automatically
// Just use logical pixels
window.SetSize(800, 600)
```

### 창이 즉시 닫힘

**원인:** 마지막 창이 닫히면 애플리케이션이 종료됨

**해결 방법:**

```go
app := application.New(application.Options{
    Mac: application.MacOptions{
        ApplicationShouldTerminateAfterLastWindowClosed: false,
    },
})
```

## 다음 단계

@cards{cols="2"}
⚙ 창 옵션
모든 창 옵션에 대한 전체 참조 문서입니다.

[자세히 알아보기 →](/features/windows/options/)

---
▣ 다중 창
다중 창 애플리케이션을 위한 패턴입니다.

[자세히 알아보기 →](/features/windows/multiple/)

---
★ 프레임 없는 창
사용자 정의 창 장식을 만듭니다.

[자세히 알아보기 →](/features/windows/frameless/)

---
🚀 창 이벤트
창 수명 주기 이벤트를 처리합니다.

[자세히 알아보기 →](/features/windows/events/)

@end

---

**궁금한 점이 있나요?** [Discord](https://discord.gg/JDdSxwjhGf)에서 질문하거나 [창 예제](https://github.com/wailsapp/wails/tree/master/v3/examples)를 확인하세요.
