---
title: "창 API"
description: "창 API 전체 레퍼런스"
slug: "reference/window"
sourcePath: "reference/window.md"
---

## 개요

창 API는 창의 모양, 동작 및 수명 주기를 제어하는 메서드를 제공합니다. 창 인스턴스 또는 `app.Window` 관리자를 통해 접근합니다.

`Window`은 `*application.WebviewWindow`이 구현하는 인터페이스이며, 아래의 메서드 시그니처는 `*WebviewWindow`에 정의되어 있습니다. 많은 변경 메서드는 연속 호출이 가능하도록 `Window`을 반환합니다. 반환 값은 각 메서드에서 설명합니다.

**일반적인 작업:**

- 창 생성 및 표시
- 크기, 위치 및 상태 제어
- 창 이벤트 처리
- 창 콘텐츠 관리
- 모양 및 동작 구성

## 표시 여부

### Show()

창을 표시합니다. 창이 숨겨져 있었다면 다시 표시됩니다. 연속 호출이 가능하도록 리시버를 반환합니다.

```go
func (w *WebviewWindow) Show() Window
```

**예:**

```go
window := app.Window.New()
window.Show()
```

### Hide()

창을 닫지 않고 숨깁니다. 창은 메모리에 유지되며 다시 표시할 수 있습니다. 연속 호출이 가능하도록 리시버를 반환합니다.

```go
func (w *WebviewWindow) Hide() Window
```

**예:**

```go
// Hide window temporarily
window.Hide()

// Show it again later
window.Show()
```

**사용 사례:**

- 트레이로 숨기는 시스템 트레이 애플리케이션
- 창을 재사용하는 마법사 흐름
- 작업 중 일시적으로 숨기기

### Close()

창을 닫습니다. 이때 `WindowClosing` 이벤트가 트리거됩니다.

```go
func (w *WebviewWindow) Close()
```

**예:**

```go
window.Close()
```

**참고:** 등록된 훅에서 `event.Cancel()`을 호출하면 창이 닫히지 않습니다.

## 창 속성

### SetTitle()

창 제목 표시줄의 텍스트를 설정합니다. 연속 호출이 가능하도록 리시버를 반환합니다.

```go
func (w *WebviewWindow) SetTitle(title string) Window
```

**매개변수:**

- `title` - 새 창 제목

**예:**

```go
window.SetTitle("My Application - Document.txt")
```

### Name()

창의 고유 이름 식별자를 반환합니다.

```go
func (w *WebviewWindow) Name() string
```

**예:**

```go
name := window.Name()
fmt.Println("Window name:", name)

// Retrieve window by name later
if w, ok := app.Window.GetByName(name); ok {
    w.Focus()
}
```

## 크기 및 위치

### SetSize()

창 크기를 픽셀 단위로 설정합니다. 연속 호출이 가능하도록 리시버를 반환합니다.

```go
func (w *WebviewWindow) SetSize(width, height int) Window
```

**매개변수:**

- `width` - 창 너비(픽셀)
- `height` - 창 높이(픽셀)

**예:**

```go
window.SetSize(1024, 768)
```

### Size()

현재 창 크기를 반환합니다.

```go
func (w *WebviewWindow) Size() (width, height int)
```

**예:**

```go
width, height := window.Size()
fmt.Printf("Window is %dx%d\n", width, height)
```

### SetMinSize() / SetMaxSize()

창의 최소 및 최대 크기를 설정합니다. 두 메서드 모두 연속 호출이 가능하도록 리시버를 반환합니다.

```go
func (w *WebviewWindow) SetMinSize(width, height int) Window
func (w *WebviewWindow) SetMaxSize(width, height int) Window
```

**예:**

```go
// Prevent window from being too small
window.SetMinSize(800, 600)

// Prevent window from being too large
window.SetMaxSize(1920, 1080)
```

### SetPosition()

화면의 왼쪽 위 모서리를 기준으로 창 위치를 설정합니다.

```go
func (w *WebviewWindow) SetPosition(x, y int)
```

**매개변수:**

- `x` - 가로 위치(픽셀)
- `y` - 세로 위치(픽셀)

**예:**

```go
// Position window at top-left
window.SetPosition(0, 0)

// Position window 100px from top-left
window.SetPosition(100, 100)
```

### Position()

현재 창 위치를 반환합니다.

```go
func (w *WebviewWindow) Position() (x, y int)
```

**예:**

```go
x, y := window.Position()
fmt.Printf("Window is at (%d, %d)\n", x, y)
```

### Center()

창을 화면 중앙에 배치합니다.

```go
func (w *WebviewWindow) Center()
```

**예:**

```go
window := app.Window.New()
window.Center()
window.Show()
```

**참고:** 기본 모니터의 중앙에 배치됩니다. 다중 모니터 구성에서는 화면 API를 참조하세요.

### Focus()

창을 맨 앞으로 가져오고 키보드 포커스를 부여합니다.

```go
func (w *WebviewWindow) Focus()
```

**예:**

```go
// Bring window to front
window.Focus()
```

## 창 상태

### Minimise() / UnMinimise()

창을 작업 표시줄/독으로 최소화하거나 복원합니다. `Minimise()`은 연속 호출이 가능하도록 리시버를 반환하고, `UnMinimise()`은 아무것도 반환하지 않습니다.

```go
func (w *WebviewWindow) Minimise() Window
func (w *WebviewWindow) UnMinimise()
```

**예:**

```go
// Minimise window
window.Minimise()

// Restore from minimised state
window.UnMinimise()
```

### Maximise() / UnMaximise()

창을 화면에 꽉 차도록 최대화하거나 이전 크기로 복원합니다. `Maximise()`은 연속 호출이 가능하도록 리시버를 반환하고, `UnMaximise()`은 아무것도 반환하지 않습니다.

```go
func (w *WebviewWindow) Maximise() Window
func (w *WebviewWindow) UnMaximise()
```

**예:**

```go
// Maximise window
window.Maximise()

// Restore to previous size
window.UnMaximise()
```

### Fullscreen() / UnFullscreen() / ToggleFullscreen()

전체 화면 모드로 전환하거나 전체 화면 모드를 종료합니다. `Fullscreen()`은 연속 호출이 가능하도록 리시버를 반환합니다.

```go
func (w *WebviewWindow) Fullscreen() Window
func (w *WebviewWindow) UnFullscreen()
func (w *WebviewWindow) ToggleFullscreen()
```

**예:**

```go
// Enter fullscreen
window.Fullscreen()

// Exit fullscreen
window.UnFullscreen()

// Or toggle
window.ToggleFullscreen()
```

`SetFullscreen(bool)` 메서드는 없습니다.

### IsMinimised() / IsMaximised() / IsFullscreen()

현재 창 상태를 확인합니다.

```go
func (w *WebviewWindow) IsMinimised() bool
func (w *WebviewWindow) IsMaximised() bool
func (w *WebviewWindow) IsFullscreen() bool
```

**예:**

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

## 창 콘텐츠

### SetURL()

창에서 지정된 URL로 이동합니다. 메서드 체이닝을 위해 수신 객체를 반환합니다.

```go
func (w *WebviewWindow) SetURL(url string) Window
```

**매개변수:**

- `url` - 이동할 URL(임베디드 에셋에는 `http://wails.localhost/` 사용 가능)

**예:**

```go
// Navigate to embedded page
window.SetURL("http://wails.localhost/settings.html")

// Navigate to external URL (if allowed)
window.SetURL("https://wails.io")
```

### SetHTML()

HTML 문자열로 창 콘텐츠를 직접 설정합니다. 메서드 체이닝을 위해 수신 객체를 반환합니다.

```go
func (w *WebviewWindow) SetHTML(html string) Window
```

**매개변수:**

- `html` - 표시할 HTML 콘텐츠

**예:**

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

**사용 사례:**

- 동적 콘텐츠 생성
- 프런트엔드 빌드 과정이 없는 간단한 창
- 오류 페이지 또는 시작 화면

### Reload()

현재 창 콘텐츠를 다시 로드합니다.

```go
func (w *WebviewWindow) Reload()
```

**예:**

```go
// Reload current page
window.Reload()
```

**참고:** 개발 중이거나 콘텐츠를 새로 고쳐야 할 때 유용합니다.

## 창 이벤트

Wails는 창 이벤트를 처리하는 두 가지 메서드를 제공합니다.

- **OnWindowEvent()** - 창 이벤트를 수신합니다(이벤트를 막을 수 없음).
- **RegisterHook()** - 창 이벤트에 훅을 연결합니다(`event.Cancel()`을 호출하여 이벤트를 막을 수 있음).

### OnWindowEvent()

창 이벤트에 대한 콜백을 등록합니다. 구독 해제 함수를 반환합니다.

```go
func (w *WebviewWindow) OnWindowEvent(
    eventType events.WindowEventType,
    callback func(event *WindowEvent),
) func()
```

**예:**

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

**일반적인 창 이벤트:**

- `events.Common.WindowClosing` - 창이 곧 닫힘
- `events.Common.WindowFocus` - 창이 포커스를 얻음
- `events.Common.WindowLostFocus` - 창이 포커스를 잃음
- `events.Common.WindowDidMove` - 창이 이동됨
- `events.Common.WindowDidResize` - 창 크기가 변경됨
- `events.Common.WindowMinimise` - 창이 최소화됨
- `events.Common.WindowMaximise` - 창이 최대화됨
- `events.Common.WindowFullscreen` - 창이 전체 화면 모드로 전환됨
- `events.Common.WindowRuntimeReady` - 창 내부 런타임이 초기화됨

### RegisterHook()

창 이벤트에 대한 훅을 등록합니다. 훅은 리스너보다 먼저 실행되며 `event.Cancel()`을 호출하여 이벤트를 막을 수 있습니다. 구독 해제 함수를 반환합니다.

```go
func (w *WebviewWindow) RegisterHook(
    eventType events.WindowEventType,
    callback func(event *WindowEvent),
) func()
```

**예 - 창 닫기 방지:**

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

**예 - 닫기 전에 저장:**

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

창의 프런트엔드로 사용자 지정 이벤트를 발생시킵니다. 훅에 의해 이벤트 발생이 취소되면 `true`을 반환합니다.

```go
func (w *WebviewWindow) EmitEvent(name string, data ...any) bool
```

**매개변수:**

- `name` - 이벤트 이름
- `data` - 이벤트와 함께 전송할 선택적 데이터

**예:**

```go
// Send data to specific window
window.EmitEvent("data-updated", map[string]any{
    "count":  42,
    "status": "success",
})
```

**프런트엔드(JavaScript):**

```javascript
import { Events } from '@wailsio/runtime'

Events.On('data-updated', (data) => {
    console.log('Count:', data.count)
    console.log('Status:', data.status)
})
```

## 기타 메서드

### SetEnabled()

사용자가 창과 상호 작용할 수 있도록 설정하거나 상호 작용을 비활성화합니다.

```go
func (w *WebviewWindow) SetEnabled(enabled bool)
```

**예:**

```go
// Disable window during long operation
window.SetEnabled(false)

// Perform operation
performLongOperation()

// Re-enable window
window.SetEnabled(true)
```

### SetBackgroundColour()

창 배경색을 설정합니다. 이 색상은 콘텐츠가 로드되기 전에 표시됩니다. 메서드 체이닝을 위해 수신 객체를 반환합니다.

```go
func (w *WebviewWindow) SetBackgroundColour(colour RGBA) Window
```

`RGBA`은 `application.RGBA{Red, Green, Blue, Alpha uint8}`입니다. `application.NewRGB(r, g, b)`(알파 255) 또는 `application.NewRGBA(r, g, b, a)` 헬퍼를 사용하세요.

**예:**

```go
// White background
window.SetBackgroundColour(application.NewRGB(255, 255, 255))

// Dark background with full alpha
window.SetBackgroundColour(application.NewRGBA(30, 30, 30, 255))
```

### SetResizable()

사용자가 창 크기를 조절할 수 있는지를 제어합니다. 메서드 체이닝을 위해 수신 객체를 반환합니다.

```go
func (w *WebviewWindow) SetResizable(resizable bool) Window
```

**예:**

```go
// Make window fixed size
window.SetResizable(false)
```

### SetAlwaysOnTop()

창을 다른 창보다 항상 위에 표시할지를 설정합니다. 메서드 체이닝을 위해 수신 객체를 반환합니다.

```go
func (w *WebviewWindow) SetAlwaysOnTop(alwaysOnTop bool) Window
```

**예:**

```go
// Keep window on top
window.SetAlwaysOnTop(true)
```

### Print()

창 콘텐츠에 대한 네이티브 인쇄 대화 상자를 엽니다.

```go
func (w *WebviewWindow) Print() error
```

**반환값:** 인쇄에 실패하면 오류를 반환합니다.

**예:**

```go
if err := window.Print(); err != nil {
    log.Println("Print failed:", err)
}
```

### AttachModal()

두 번째 창을 시트 모달로 연결합니다.

```go
func (w *WebviewWindow) AttachModal(modalWindow Window)
```

**매개변수:**

- `modalWindow` - 모달로 연결할 창

**플랫폼 지원:**

- **macOS**: 완전 지원(시트로 표시)
- **Windows**: 지원하지 않음
- **Linux**: 지원하지 않음

**예:**

```go
modalWindow := app.Window.New()
window.AttachModal(modalWindow)
```

## 플랫폼별 옵션

### Linux

Linux 창에서는 `LinuxWindow`을(를) 통해 다음과 같은 플랫폼별 옵션을 사용할 수 있습니다.

#### MenuStyle

애플리케이션 메뉴의 표시 방식을 제어합니다. 이 옵션은 기본 GTK4 빌드에서 사용할 수 있으며, 레거시 `-tags gtk3` 빌드에서는 무시됩니다.

| 값 | 설명 |
| --- | --- |
| `LinuxMenuStyleMenuBar` | 제목 표시줄 아래의 기존 메뉴 모음(기본값) |
| `LinuxMenuStylePrimaryMenu` | 헤더 표시줄의 기본 메뉴 버튼(GNOME 스타일) |

**예:**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "My Application",
    Linux: application.LinuxWindow{
        MenuStyle: application.LinuxMenuStylePrimaryMenu,
    },
})
window.SetMenu(menu)
```

**참고:** 기본 메뉴 스타일은 GNOME Human Interface Guidelines에 따라 헤더 표시줄에 햄버거 버튼(☰)을 표시합니다. 최신 GNOME 애플리케이션에는 이 스타일을 사용하는 것이 좋습니다.

## 전체 예제

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
