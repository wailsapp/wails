---
title: "창 옵션"
description: "WebviewWindowOptions 전체 참조"
slug: "features/windows/options"
sourcePath: "features/windows/options.md"
---

## 창 구성 옵션

Wails는 크기, 위치, 모양, 동작을 위한 수십 가지 옵션으로 포괄적인 창 구성을 제공합니다. 이 참조에서는 `WebviewWindowOptions`에 대한 <strong>전체 참조</strong>로서 Windows, macOS, Linux에서 사용할 수 있는 모든 옵션을 다룹니다. 모든 옵션과 모든 플랫폼을 예제 및 제약 조건과 함께 설명합니다.

## WebviewWindowOptions 구조

```go
type WebviewWindowOptions struct {
    // Identity
    Name  string
    Title string

    // Size and Position
    Width           int
    Height          int
    X               int
    Y               int
    MinWidth        int
    MinHeight       int
    MaxWidth        int
    MaxHeight       int
    InitialPosition WindowStartPosition // WindowCentered (default) or WindowXY
    Screen          *Screen             // target screen for initial placement

    // Initial State
    Hidden        bool
    Frameless     bool
    DisableResize bool        // inverted vs v2's `Resizable`
    AlwaysOnTop   bool
    StartState    WindowState // WindowStateNormal | Minimised | Maximised | Fullscreen

    // Appearance
    BackgroundColour RGBA
    BackgroundType   BackgroundType
    Zoom             float64
    ZoomControlEnabled bool

    // Content
    URL  string
    HTML string
    JS   string
    CSS  string

    // Behaviour
    EnableFileDrop              bool
    IgnoreMouseEvents           bool
    HideOnFocusLost             bool
    HideOnEscape                bool
    DevToolsEnabled             bool
    DefaultContextMenuDisabled  bool
    ContentProtectionEnabled    bool
    KeyBindings                 map[string]func(window *WebviewWindow)

    // Permissions
    Permissions map[PermissionType]Permission

    // Window-control button states
    MinimiseButtonState ButtonState
    MaximiseButtonState ButtonState
    CloseButtonState    ButtonState

    // Menu
    UseApplicationMenu bool

    // Platform-specific (per-window)
    Mac     MacWindow
    Windows WindowsWindow
    Linux   LinuxWindow
}
```

`WebviewWindowOptions`에는 `Parent` 필드가 **없습니다**. 부모/모달 관계에는 `parentWindow.AttachModal(childWindow)`을 사용하세요. 또한 `Assets` 필드도 **없습니다**. 에셋 구성은 `application.Options`(`Assets AssetOptions`)에 있습니다.

전체 소스: [`v3/pkg/application/webview_window_options.go`](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window_options.go).

## 핵심 옵션

### Name

**유형:** `string` **기본값:** 자동 생성 UUID **플랫폼:** 모두

```go
Name: "main-window"
```

**용도:** 나중에 창을 찾을 때 사용하는 고유 식별자입니다.

**권장 사항:**

- 설명적인 이름을 사용하세요: `"main"`, `"settings"`, `"about"`
- kebab-case를 사용하세요: `"file-browser"`, `"color-picker"`
- 짧고 기억하기 쉽게 유지하세요

**예:**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name: "settings-window",
})

// Later...
if settings, ok := app.Window.GetByName("settings-window"); ok {
    settings.Focus()
}
```

### Title

**유형:** `string` **기본값:** 애플리케이션 이름 **플랫폼:** 모두

```go
Title: "My Application"
```

**용도:** 제목 표시줄과 작업 표시줄에 표시되는 텍스트입니다.

**동적 업데이트:**

```go
window.SetTitle("My Application - Document.txt")
```

### Width / Height

**유형:** `int`(픽셀) **기본값:** 800 x 600 **플랫폼:** 모두 **제약 조건:** 양수여야 함

```go
Width:  1200,
Height: 800,
```

**용도:** 논리 픽셀 단위의 초기 창 크기입니다.

**참고:**

- Wails가 DPI 배율 조정을 자동으로 처리합니다
- 물리 픽셀이 아닌 논리 픽셀을 사용하세요
- 최소 화면 해상도(1024x768)를 고려하세요

**크기 예시:**

| 사용 사례 | 너비 | 높이 |
| --- | --- | --- |
| 소형 유틸리티 | 400 | 300 |
| 표준 애플리케이션 | 1024 | 768 |
| 대형 애플리케이션 | 1440 | 900 |
| Full HD | 1920 | 1080 |

### X / Y

**유형:** `int`(픽셀) **기본값:** 화면 중앙 **플랫폼:** 모두

```go
X: 100,  // 100px from left edge
Y: 100,  // 100px from top edge
```

**용도:** 창의 초기 위치입니다.

**좌표계:**

- (0, 0)은 기본 화면의 왼쪽 위 모서리입니다
- X의 양의 방향은 오른쪽입니다
- Y의 양의 방향은 아래쪽입니다

**예:**

`X` 및 `Y`은 `InitialPosition: application.WindowXY`이 설정된 경우에만 적용됩니다. 설정하지 않으면 `InitialPosition`의 기본값은 `WindowCentered`이며 `X`/`Y`은 무시됩니다.

```go
settings := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:            "coordinate-window",
    InitialPosition: application.WindowXY, // opt into X/Y coordinates
    X:               100,
    Y:               100,
})
```

**권장 사항:** 특정 좌표가 중요하지 않다면 창을 만든 후 `Center()`을 사용하여 중앙에 배치하세요:

```go
window := app.Window.New()
window.Center()
```

### MinWidth / MinHeight

**유형:** `int`(픽셀) **기본값:** 0(최솟값 없음) **플랫폼:** 모두

```go
MinWidth:  400,
MinHeight: 300,
```

**용도:** 창이 너무 작아지지 않도록 합니다.

**사용 사례:**

- 레이아웃 깨짐 방지
- 사용 편의성 보장
- 가로세로 비율 유지

**예:**

```go
// Prevent window smaller than 400x300
MinWidth:  400,
MinHeight: 300,
```

### MaxWidth / MaxHeight

**유형:** `int`(픽셀) **기본값:** 0(최댓값 없음) **플랫폼:** 모두

```go
MaxWidth:  1920,
MaxHeight: 1080,
```

**용도:** 창이 너무 커지지 않도록 합니다.

**사용 사례:**

- 고정 크기 애플리케이션
- 과도한 리소스 사용 방지
- 디자인 제약 조건 유지

## 상태 옵션

### Hidden

**타입:** `bool` **기본값:** `false` **플랫폼:** 모두

```go
Hidden: true,
```

**용도:** 창을 표시하지 않고 생성합니다.

**사용 사례:**

- 백그라운드 창
- 필요할 때 표시하는 창
- 스플래시 화면(생성하고 로드한 후 표시)
- 콘텐츠를 로드하는 동안 흰색 화면이 깜박이는 현상 방지

**플랫폼별 개선 사항:**

- **Windows:** 흰색 창이 깜박이는 문제를 해결했습니다. `Show()`이 호출될 때까지 창이 보이지 않습니다.
- **macOS:** 완전히 지원
- **Linux:** 완전히 지원

**매끄러운 로딩을 위한 권장 패턴:**

```go
// Create hidden window
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:             "main-window",
    Hidden:           true,
    BackgroundColour: application.NewRGB(30, 30, 30), // Match your theme
})

// Load content while hidden
// ... content loads ...

// Show when ready (no flash!)
window.Show()
```

**예:**

```go
settings := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:   "settings",
    Hidden: true,
})

// Show when needed
settings.Show()
```

### Frameless

**타입:** `bool` **기본값:** `false` **플랫폼:** 모두

```go
Frameless: true,
```

**용도:** 제목 표시줄과 창 테두리를 제거합니다.

**사용 사례:**

- 사용자 정의 창 장식
- 스플래시 화면
- 키오스크 애플리케이션
- 사용자 정의 디자인 창

**중요:** 다음 기능을 직접 구현해야 합니다.

- 창 끌기
- 닫기/최소화/최대화 버튼
- 크기 조절 핸들(창 크기를 조절할 수 있는 경우)

**자세한 내용은 [프레임 없는 창](/features/windows/frameless/)을 참조하세요.**

### DisableResize

**타입:** `bool` **기본값:** `false`(기본적으로 창 크기를 조절할 수 있음) **플랫폼:** 모두

```go
DisableResize: true,
```

**용도:** 창 크기 조절을 방지합니다. 이 필드는 v2의 `Resizable`와 <strong>반대</strong>라는 점에 유의하세요. 창 크기를 조절할 수 없게 하려면 `DisableResize: true`으로 설정하세요.

**사용 사례:**

- 고정 크기 애플리케이션
- 스플래시 화면
- 대화 상자

**참고:** `MaximiseButtonState` 또는 컬렉션 동작을 통해 해당 기능도 비활성화하지 않으면 사용자가 창을 최대화하거나 전체 화면으로 전환할 수 있습니다.

### AlwaysOnTop

**타입:** `bool` **기본값:** `false` **플랫폼:** 모두

```go
AlwaysOnTop: true,
```

**용도:** 창을 다른 모든 창 위에 유지합니다.

**사용 사례:**

- 플로팅 도구 모음
- 알림
- 화면 속 화면
- 타이머

**플랫폼 참고 사항:**

- **macOS:** 완전히 지원
- **Windows:** 완전히 지원
- **Linux:** 창 관리자에 따라 다름

### StartState

**타입:** `WindowState` 열거형 **기본값:** `WindowStateNormal` **플랫폼:** 모두

```go
StartState: application.WindowStateMaximised,
```

**용도:** 창이 표시될 때의 초기 상태입니다.

**값:**

- `WindowStateNormal` - 일반 창
- `WindowStateMinimised` - 최소화
- `WindowStateMaximised` - 최대화
- `WindowStateFullscreen` - 전체 화면

`WindowStateHidden` 상수는 없습니다. 창을 보이지 않는 상태로 시작하려면 `Hidden` 불리언 필드를 사용하세요.

**런타임에 전체 화면 전환:**

```go
window.Fullscreen()
window.UnFullscreen()
window.ToggleFullscreen() // there is no SetFullscreen(bool)
```

## 모양 옵션

### BackgroundColour

**유형:** `RGBA` 구조체 **기본값:** 흰색 **플랫폼:** 모두

```go
BackgroundColour: application.RGBA{Red: 0, Green: 0, Blue: 0, Alpha: 255},
```

`RGBA` 필드는 `Red, Green, Blue, Alpha`(uint8)입니다. `application.NewRGB(r, g, b)`(알파 255) 또는 `application.NewRGBA(r, g, b, a)` 헬퍼를 사용하는 것이 좋습니다.

**용도:** 콘텐츠가 로드되기 전의 창 배경색입니다.

**사용 사례:**

- 앱 테마와 일치시키기
- 어두운 테마에서 흰색 화면이 깜박이는 현상 방지
- 매끄러운 로딩 환경 제공

**예:**

```go
// Dark theme
BackgroundColour: application.NewRGB(30, 30, 30),

// Light theme
BackgroundColour: application.NewRGB(255, 255, 255),
```

**헬퍼 메서드:**

```go
window.SetBackgroundColour(application.NewRGB(30, 30, 30))
```

### BackgroundType

**유형:** `BackgroundType` 열거형 **기본값:** `BackgroundTypeSolid` **플랫폼:** macOS, Windows(일부 지원)

```go
BackgroundType: application.BackgroundTypeTranslucent,
```

**값:**

- `BackgroundTypeSolid` - 단색
- `BackgroundTypeTransparent` - 완전 투명
- `BackgroundTypeTranslucent` - 반투명 블러

**플랫폼 지원:**

- **macOS:** `Mac.Backdrop`을 구성하세요. 웹뷰를 투명하게 하려면 [`private_mac_apis`](/guides/build/private-macos-apis/#webview-transparency-and-background)이 필요합니다. 이 설정이 없으면 웹뷰는 불투명하게 유지됩니다.
- **Windows:** 투명 및 반투명(Windows 11 이상)
- **Linux:** 단색만 지원

**예(macOS):**

```go
BackgroundType: application.BackgroundTypeTranslucent,
Mac: application.MacWindow{
    Backdrop: application.MacBackdropTranslucent,
},
```

### OpenInspectorOnStartup 및 OpenDevTools

**macOS의 비공개 API:** `OpenInspectorOnStartup: true`, Go `window.OpenDevTools()` 및 JavaScript `Window.OpenDevTools()`에서 프로그래밍 방식으로 인스펙터를 열려면 `private_mac_apis`이 필요합니다. 이 설정이 없으면 이러한 작업은 아무 동작도 하지 않습니다. 프로덕션 빌드에는 `devtools`도 필요합니다. macOS 13.3 이상에서 공개 Safari 검사 기능을 사용할 때는 비공개 API가 필요하지 않지만, 이전 버전의 macOS에서 인스펙터를 활성화할 때는 필요합니다. [Web Inspector 빌드 매트릭스](/guides/build/private-macos-apis/#web-inspector)를 참조하세요.

## 콘텐츠 옵션

### URL

**유형:** `string` **기본값:** 비어 있음(Assets에서 로드) **플랫폼:** 모두

```go
URL: "https://example.com",
```

**용도:** 임베드된 애셋 대신 외부 URL을 로드합니다.

**사용 사례:**

- 개발 환경(개발 서버에서 로드)
- 웹 기반 애플리케이션
- 하이브리드 애플리케이션

**예:**

```go
// Development — point the window at the Vite dev server
URL: "http://localhost:9245",

// Production — embedded assets are configured at the application level
// (Assets is application.Options.Assets, not a WebviewWindowOptions field).
```

프로덕션 사례의 애플리케이션 수준 코드 조각:

```go
app := application.New(application.Options{
    Name: "My App",
    Assets: application.AssetOptions{
        Handler: application.AssetFileServerFS(assets),
    },
})
```

### HTML

**유형:** `string` **기본값:** 비어 있음 **플랫폼:** 모두

```go
HTML: "<h1>Hello World</h1>",
```

**용도:** HTML 문자열을 직접 로드합니다.

**사용 사례:**

- 단순한 창
- 생성된 콘텐츠
- 테스트

**예:**

```go
HTML: `
<!DOCTYPE html>
<html>
<head><title>Simple Window</title></head>
<body><h1>Hello from Wails!</h1></body>
</html>
`,
```

### Assets(애플리케이션 수준에서만 사용)

애셋 구성은 **`WebviewWindowOptions` 필드가 아닙니다**. 프런트엔드 애셋은 애플리케이션 자체가 `application.Options.Assets`(`AssetOptions`)을 통해 제공하며, 모든 창은 해당 애셋 서버를 상속합니다.

```go
app := application.New(application.Options{
    Assets: application.AssetOptions{
        Handler: application.AssetFileServerFS(assets),
    },
})
```

**자세한 내용은 [빌드 시스템](/concepts/build-system/)을 참조하세요.**

### UseApplicationMenu

**유형:** `bool` **기본값:** `false` **플랫폼:** Windows, Linux(macOS에서는 효과 없음)

```go
UseApplicationMenu: true,
```

**용도:** 이 창에 `app.Menu.Set()`을 통해 설정한 애플리케이션 메뉴를 사용합니다.

<strong>macOS</strong>에서는 항상 화면 상단의 전역 애플리케이션 메뉴를 사용하므로 이 옵션은 아무런 효과가 없습니다.

**Windows** 및 <strong>Linux</strong>에서는 기본적으로 창에 메뉴가 표시되지 않습니다. `UseApplicationMenu: true`을 설정하면 창에서 애플리케이션 수준 메뉴를 사용하므로 간단한 크로스 플랫폼 솔루션을 구현할 수 있습니다.

**예:**

```go
// Set the application menu once
menu := app.NewMenu()
menu.AddRole(application.FileMenu)
menu.AddRole(application.EditMenu)
app.Menu.Set(menu)

// All windows with UseApplicationMenu will display this menu
app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:              "Main Window",
    UseApplicationMenu: true,
})

app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:              "Second Window",
    UseApplicationMenu: true,  // Also gets the app menu
})
```

**참고:**

- `UseApplicationMenu`과 창별 메뉴를 모두 설정하면 창별 메뉴가 우선 적용됩니다
- 런타임에 OS를 확인할 필요가 없어 크로스 플랫폼 코드가 간소화됩니다
- 메뉴에 관한 전체 문서는 [애플리케이션 메뉴](/features/menus/application/)를 참조하세요

## 입력 옵션

### EnableFileDrop

**유형:** `bool` **기본값:** `false` **플랫폼:** 모두

```go
EnableFileDrop: true,
```

**용도:** 운영 체제에서 창으로 파일을 끌어다 놓을 수 있게 합니다.

활성화하면 다음과 같이 동작합니다:

- 파일 관리자에서 끌어온 파일을 애플리케이션에 놓을 수 있습니다
- 파일을 놓으면 해당 파일 경로와 함께 `WindowFilesDropped` 이벤트가 발생합니다
- `data-file-drop-target` 특성이 있는 요소에서 자세한 파일 놓기 정보를 제공합니다

**사용 사례:**

- 파일 업로드 인터페이스
- 문서 편집기
- 미디어 가져오기 도구
- 파일을 받는 모든 앱

**예:**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:          "File Uploader",
    EnableFileDrop: true,
})

// Handle dropped files
window.OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {
    files := event.Context().DroppedFiles()
    details := event.Context().DropTargetDetails()
    
    for _, file := range files {
        fmt.Println("Dropped:", file)
    }
})
```

**HTML 파일 놓기 영역:**

```html
<!-- Mark elements as drop targets -->
<div id="upload" data-file-drop-target>
    Drop files here
</div>
```

**전체 문서는 [파일 놓기](/features/drag-and-drop/files/)를 참조하세요.**

## 보안 옵션

### ContentProtectionEnabled

**유형:** `bool` **기본값:** `false` **플랫폼:** Windows(10 이상), macOS

```go
ContentProtectionEnabled: true,
```

**용도:** 창 콘텐츠의 화면 캡처를 방지합니다.

**플랫폼 지원:**

- **Windows:** Windows 10 빌드 19041 이상(완전 지원), 이전 버전(부분 지원)
- **macOS:** 완전 지원
- **Linux:** 지원하지 않음

**사용 사례:**

- 뱅킹 애플리케이션
- 암호 관리자
- 의료 기록
- 기밀 문서

**중요 사항:**

1. 카메라로 직접 촬영하는 것은 방지하지 못합니다
2. 일부 도구는 보호 기능을 우회할 수 있습니다
3. 포괄적인 보안 체계의 일부일 뿐, 유일한 보호 수단은 아닙니다
4. DevTools 창은 자동으로 보호되지 않습니다

**예:**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Secure Window",
    ContentProtectionEnabled: true,
})

// Toggle at runtime
window.SetContentProtection(true)
```

### Permissions

**유형:** `map[PermissionType]Permission` **기본값:** `nil`(플랫폼 기본 처리) **플랫폼:** Linux, Windows(macOS에서는 TCC에 위임)

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionAllow,
    application.PermissionCamera:     application.PermissionAllow,
},
```

**용도:** 플랫폼별 코드를 사용하지 않고 창의 웹 콘텐츠에서 요청하는 기능(카메라, 마이크, 위치 정보, 알림, 클립보드 읽기)을 처리하는 방법을 선언적으로 제어합니다.

**PermissionType 값:** `PermissionMicrophone`, `PermissionCamera`, `PermissionGeolocation`, `PermissionNotifications`, `PermissionClipboardRead`

**권한 값:**

- `PermissionDefault`(0) — 플랫폼의 기본 처리 방식: macOS/Windows에서는 OS/WebView2 프롬프트를 표시하며, Linux에서는 카메라와 마이크를 허용하고 나머지는 모두 거부합니다
- `PermissionAllow`(1) — 프롬프트 없이 허용합니다(Linux에서는 카메라와 마이크만 구현되어 있으며 다른 유형은 계속 거부됩니다)
- `PermissionDeny`(2) — 프롬프트 없이 거부합니다

**중요 — Windows:** 이 옵션이 도입되기 전에는 Wails가 모든 WebView2 기능을 사용자에게 알리지 않고 허용했습니다. 이제 `Permissions`에 항목을 하나라도 설정하면 이러한 일괄 허용이 비활성화됩니다. 목록에 없는 기능은 자동 허용되지 않고 WebView2의 기본 프롬프트가 표시됩니다. 앱에 필요한 모든 기능을 명시적으로 나열하세요.

**전체 가이드, 플랫폼 지원표 및 예제는 [권한](/features/windows/permissions/)을 참조하세요.**

## 창 수명 주기 이벤트

창 수명 주기 이벤트는 `OnWindowEvent` 및 `RegisterHook`을 사용하여 처리합니다. 이 메서드를 사용하면 창을 닫고 제거하는 동작을 세밀하게 제어할 수 있습니다.

### 창 닫기 취소

창이 닫히지 않도록 하려면(예: 저장하지 않은 변경 사항이 있는 경우) `WindowClosing` 이벤트와 함께 `RegisterHook`을 사용하고 `event.Cancel()`을 호출하세요:

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:  "main-window",
    Title: "My Application",
})

// Register a hook to intercept the closing event
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

**핵심 사항:**

- `RegisterHook`은 이벤트가 발생하기 전에 이를 가로챕니다.
- 창이 닫히지 않도록 하려면 `event.Cancel()`을 호출하세요.
- 취소한 후에도 창은 열린 상태로 유지됩니다.

### 창 닫기 처리

창이 닫힐 때 정리 작업을 수행하려면 `WindowClosing` 이벤트와 함께 `OnWindowEvent`을 사용하세요.

```go
window.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
    // Cleanup code runs here
    fmt.Printf("Window %s is closing\n", window.Name())

    // Close database connection
    if db != nil {
        db.Close()
    }

    // Remove from window list
    removeWindow(window.ID())
})
```

**핵심 사항:**

- `OnWindowEvent`은 곧 발생할 이벤트를 처리합니다.
- 창이 제거되기 전에 정리 작업이 실행됩니다.
- 여기에서는 닫기를 취소할 수 없습니다(취소하려면 `RegisterHook`을 사용하세요).

### 싱글턴 창 정리 패턴

싱글턴 창(인스턴스가 하나만 존재하도록 보장하는 창)의 참조를 정리하려면 `WindowClosing`을 사용하세요.

```go
var settingsWindow *application.WebviewWindow

func ShowSettings(app *application.App) {
    // Create if doesn't exist
    if settingsWindow == nil {
        settingsWindow = app.Window.NewWithOptions(application.WebviewWindowOptions{
            Name:  "settings",
            Title: "Settings",
            Width: 600,
            Height: 400,
        })

        // Cleanup on close
        settingsWindow.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
            settingsWindow = nil
        })
    }

    // Show and focus
    settingsWindow.Show()
    settingsWindow.Focus()
}
```

## 플랫폼별 옵션

### Mac 옵션

```go
Mac: application.MacWindow{
    TitleBar: application.MacTitleBar{
        AppearsTransparent: true,
        Hide:               false,
        HideTitle:          true,
        FullSizeContent:    true,
    },
    Backdrop:                application.MacBackdropTranslucent,
    InvisibleTitleBarHeight: 50,
    WindowClass:             application.MacWindowClassPanel,
    PanelPreferences: application.MacPanelPreferences{
        NonActivating:          true,
        FloatingPanel:          true,
        BecomesKeyOnlyIfNeeded: false,
        UtilityWindow:          false,
    },
    WindowLevel:             application.MacWindowLevelFloating,
    CollectionBehavior:      application.MacWindowCollectionBehaviorDefault,
    TabbingMode:             application.MacWindowTabbingModeDisallowed,
},
```

**TitleBar**(`MacTitleBar`)

- `AppearsTransparent` - 제목 표시줄을 투명하게 만들고 콘텐츠를 제목 표시줄 영역까지 확장합니다.
- `Hide` - 제목 표시줄을 완전히 숨깁니다.
- `HideTitle` - 제목 텍스트만 숨깁니다.
- `FullSizeContent` - 콘텐츠를 창 전체 크기로 확장합니다.

macOS에서 webview 투명도와 프로그래밍 방식의 인스펙터 열기에는 `private_mac_apis` 빌드 태그가 필요합니다. 이 태그가 없어도 같은 옵션은 유효하지만 비공개 API만 사용하는 작업은 아무 동작도 하지 않습니다. Liquid Glass 그룹화도 무시되며 스타일에는 공개 API 기반 대안이 사용됩니다. 빌드 명령과 정확한 동작은 [비공개 macOS API](/guides/build/private-macos-apis/)를 참조하세요.

**Backdrop**(`MacBackdrop`)

- `MacBackdropNormal` - 표준 불투명 배경
- `MacBackdropTranslucent` - **webview 투명도에는 비공개 API가 필요합니다.** 태그가 없으면 네이티브 블러가 불투명한 webview 뒤에 그대로 유지됩니다.
- `MacBackdropTransparent` - **webview 투명도에는 비공개 API가 필요합니다.** 태그가 없으면 webview는 불투명한 상태로 유지됩니다.
- `MacBackdropLiquidGlass` - **webview 투명도에는 비공개 API가 필요합니다.** 태그가 없으면 유리 레이어가 불투명한 webview 뒤에 그대로 유지되며 스타일에는 공개 API 기반 대안이 사용됩니다.

**LiquidGlass**(`MacLiquidGlass`)

| 필드 또는 값 | macOS 비공개 API 의존성 |
| --- | --- |
| `Style: LiquidGlassStyleAutomatic` | 네이티브 일반 스타일은 공개 API이며, 배경 webview 투명도에는 `private_mac_apis`이 필요합니다. |
| `Style: LiquidGlassStyleLight` | 이 태그를 사용하면 기존 네이티브 투명 스타일 매핑이 유지됩니다. 태그가 없으면 Wails는 Aqua 모양의 일반 유리를 사용합니다. |
| `Style: LiquidGlassStyleDark` | **비공개 API:** 문서화되지 않은 네이티브 스타일 값 `2`. 태그가 없으면 Dark Aqua 모양의 일반 유리를 사용합니다. |
| `Style: LiquidGlassStyleVibrant` | 네이티브 투명 스타일 매핑은 공개 API이며, 배경 webview 투명도에는 태그가 필요합니다. |
| `GroupID` | **비공개 API:** 비어 있지 않은 값은 그룹화를 요청하며, 태그가 없으면 무시됩니다. |
| `GroupSpacing` | **비공개 API:** 양수 값은 그룹 간격을 요청하며, 태그가 없으면 무시됩니다. |
| `Material`, `CornerRadius`, `TintColor` | 이 필드들 자체에는 비공개 API 의존성이 없습니다. |

네이티브 스타일 값과 OS 가용성은 [Liquid Glass 값](/guides/build/private-macos-apis/#liquid-glass-values)을 참조하세요.

**InvisibleTitleBarHeight**(`int`)

- 보이지 않는 제목 표시줄 영역의 높이(드래그용)
- 네이티브 제목 표시줄의 드래그 영역이 숨겨진 경우, 즉 창이 프레임리스(`Frameless: true`)이거나 투명한 제목 표시줄(`AppearsTransparent: true`)을 사용하는 경우에만 적용됩니다.
- 제목 표시줄이 보이는 표준 창에는 영향을 주지 않습니다.

**WindowClass**(`MacWindowClass`)

- `MacWindowClassWindow` - 표준 `NSWindow` 동작(기본값)
- `MacWindowClassPanel` - 애플리케이션의 메인 창이 되지 않는 보조 `NSPanel`

`PanelPreferences`은 `MacWindowClassPanel`에만 적용됩니다.

- `NonActivating`은 `NSWindowStyleMaskNonactivatingPanel`을 추가합니다. 패널을 표시하거나 포커스해도 Wails 애플리케이션은 활성화되지 않지만, 컨트롤 및 텍스트 입력을 위해 패널이 키 윈도우가 될 수는 있습니다.
- `FloatingPanel`은 AppKit의 플로팅 패널 동작을 활성화합니다.
- `BecomesKeyOnlyIfNeeded`은 클릭한 뷰가 키보드 입력을 요청할 때만 키 상태를 얻습니다.
- `UtilityWindow`은 네이티브 유틸리티 윈도우 스타일을 적용합니다.

Wails 패널은 애플리케이션이 비활성화되어도 계속 표시되며, 닫히면 해제되어 `WebviewWindow`에서 예상하는 수명 주기와 일치합니다. 이는 반대 동작을 하는 `NSPanel` 기본값을 의도적으로 재정의한 것입니다.

윈도우 클래스, 레벨, 활성화 정책 및 컬렉션 동작은 각각 별개의 문제를 해결합니다:

- `WindowClass`은 `NSWindow` 또는 `NSPanel`을 선택하고 메인/키 윈도우의 의미 체계를 제어합니다.
- `WindowLevel`은 Z 순서를 제어합니다. 메뉴 막대 오버레이에는 `MacWindowLevelPopUpMenu`을 사용하세요.
- `MacOptions.ActivationPolicy`은 Dock 및 메뉴 막대 표시를 포함하여 애플리케이션 전체를 제어합니다. 애플리케이션을 활성화하지 않는 패널에는 액세서리 활성화 정책이 필요하지 않지만, 앱에서 Dock 아이콘을 숨기기 위해 이 정책을 사용할 수는 있습니다.
- `CollectionBehavior`은 Spaces 및 전체 화면 참여를 제어합니다.

```go
// Spotlight/menu-bar panel that leaves the current application active.
Mac: application.MacWindow{
    WindowClass: application.MacWindowClassPanel,
    PanelPreferences: application.MacPanelPreferences{
        NonActivating: true,
    },
    WindowLevel: application.MacWindowLevelPopUpMenu,
    CollectionBehavior: application.MacWindowCollectionBehaviorCanJoinAllSpaces |
        application.MacWindowCollectionBehaviorFullScreenAuxiliary |
        application.MacWindowCollectionBehaviorStationary,
},
```

**WindowLevel** (`MacWindowLevel`)

- `MacWindowLevelNormal` - 표준 윈도우 레벨(기본값)
- `MacWindowLevelFloating` - 일반 윈도우 위에 떠 있음
- `MacWindowLevelTornOffMenu` - 분리된 메뉴 레벨
- `MacWindowLevelModalPanel` - 모달 패널 레벨
- `MacWindowLevelMainMenu` - 메인 메뉴 레벨
- `MacWindowLevelStatus` - 상태 윈도우 레벨
- `MacWindowLevelPopUpMenu` - 팝업 메뉴 레벨
- `MacWindowLevelScreenSaver` - 화면 보호기 레벨

명시적인 `WindowLevel`은 `AlwaysOnTop` 및 `PanelPreferences.FloatingPanel`보다 우선합니다. 레벨을 명시하지 않으면 `AlwaysOnTop` 또는 플로팅 패널은 `MacWindowLevelFloating`로 결정되고, 그 외에는 레벨이 `MacWindowLevelNormal`로 결정됩니다. 이후에 `SetAlwaysOnTop`을 호출하는 것은 여전히 명시적인 런타임 변경으로 처리됩니다.

**CollectionBehavior** (`MacWindowCollectionBehavior`)

macOS Spaces와 전체 화면에서 윈도우가 동작하는 방식을 제어합니다. 이 값들은 비트 OR(`|`)을 사용해 조합할 수 있는 비트 마스크 값입니다.

**Space 동작:**

- `MacWindowCollectionBehaviorDefault` - FullScreenPrimary 사용(기본값, 이전 버전과 호환)
- `MacWindowCollectionBehaviorCanJoinAllSpaces` - 모든 Spaces에 윈도우 표시
- `MacWindowCollectionBehaviorMoveToActiveSpace` - 표시될 때 활성 Space로 이동
- `MacWindowCollectionBehaviorManaged` - 기본 관리형 윈도우 동작
- `MacWindowCollectionBehaviorTransient` - 임시/일시적 윈도우
- `MacWindowCollectionBehaviorStationary` - Space를 전환하는 동안 제자리에 유지

**윈도우 순환:**

- `MacWindowCollectionBehaviorParticipatesInCycle` - Cmd+` 순환에 포함
- `MacWindowCollectionBehaviorIgnoresCycle` - Cmd+` 순환에서 제외

**전체 화면 동작:**

- `MacWindowCollectionBehaviorFullScreenPrimary` - 전체 화면 모드로 전환 가능
- `MacWindowCollectionBehaviorFullScreenAuxiliary` - 전체 화면 앱 위에 오버레이 가능
- `MacWindowCollectionBehaviorFullScreenNone` - 전체 화면 기능 비활성화
- `MacWindowCollectionBehaviorFullScreenAllowsTiling` - 나란히 타일링 허용(macOS 10.11 이상)
- `MacWindowCollectionBehaviorFullScreenDisallowsTiling` - 타일링 방지(macOS 10.11 이상)

**예제 - Spotlight와 유사한 윈도우:**

```go
// Window that appears on all Spaces AND can overlay fullscreen apps
Mac: application.MacWindow{
	WindowClass: application.MacWindowClassPanel,
	PanelPreferences: application.MacPanelPreferences{
		NonActivating: true,
	},
    CollectionBehavior: application.MacWindowCollectionBehaviorCanJoinAllSpaces |
                        application.MacWindowCollectionBehaviorFullScreenAuxiliary,
    WindowLevel:        application.MacWindowLevelFloating,
},
```

**예제 - 단일 동작:**

```go
// Window that can appear over fullscreen applications
Mac: application.MacWindow{
    CollectionBehavior: application.MacWindowCollectionBehaviorFullScreenAuxiliary,
},
```

**TabbingMode** (`MacWindowTabbingMode`)

macOS 10.12 이상에서 윈도우 탭 동작을 제어합니다. 윈도우 탭 기능을 사용하면 여러 윈도우를 탭으로 그룹화할 수 있습니다.

**옵션:**

- `MacWindowTabbingModeDefault` - 영 값 센티널(명시적으로 설정되지 않음). 런타임에는 기본적으로 탭 기능을 허용하지 않음
- `MacWindowTabbingModeAutomatic` - 시스템에서 탭 동작 결정
- `MacWindowTabbingModePreferred` - 윈도우가 탭 모드를 우선 사용
- `MacWindowTabbingModeDisallowed` - 윈도우 탭 기능 비활성화

**예제 - 윈도우 탭 기능 비활성화:**

```go
Mac: application.MacWindow{
    TabbingMode: application.MacWindowTabbingModeDisallowed,
},
```

**예제 - 윈도우 탭 기능 우선 사용:**

```go
Mac: application.MacWindow{
    TabbingMode: application.MacWindowTabbingModePreferred,
},
```

**WebviewPreferences** (`MacWebviewPreferences`)

기반 `WKWebView` 구성을 세밀하게 제어합니다. 모든 필드는 선택 사항이며, 설정하지 않은 필드는 WebKit 기본값을 변경하지 않습니다.

```go
Mac: application.MacWindow{
    WebviewPreferences: application.MacWebviewPreferences{
        TabFocusesLinks:                       optional.True,
        TextInteractionEnabled:                optional.True,
        FullscreenEnabled:                     optional.False,
        AllowsBackForwardNavigationGestures:   optional.False,
        AllowsMagnification:                   optional.True,
        AllowsAirPlayForMediaPlayback:         optional.True,
        JavaScriptCanOpenWindowsAutomatically: optional.False,
        MinimumFontSize:                       optional.NewVar(12.0),
        ApplicationNameForUserAgent:           "MyApp",
        EnableAutoplayWithoutUserAction:       optional.True,
    },
},
```

- `TabFocusesLinks` — `true`이면 Tab 키를 누를 때 링크와 폼 컨트롤로 포커스가 이동합니다(기본값: `false`).
- `TextInteractionEnabled` — `true`이면 사용자가 웹뷰의 텍스트를 선택하고 상호 작용할 수 있습니다(기본값: `true`).
- `FullscreenEnabled` — `true`이면 웹 콘텐츠가 HTML Fullscreen API를 통해 전체 화면으로 전환될 수 있습니다(기본값: `false`). macOS 12.3 이상이 필요합니다.
- `AllowsBackForwardNavigationGestures` — `true`이면 가로 스와이프 제스처로 뒤로/앞으로 탐색할 수 있습니다(기본값: `false`).
- `AllowsMagnification` — `true`이면 웹뷰에서 핀치 투 줌을 사용할 수 있습니다(기본값: `false`).
- `AllowsAirPlayForMediaPlayback` — `true`이면 미디어를 AirPlay 기기로 스트리밍할 수 있습니다(기본값: `true`).
- `JavaScriptCanOpenWindowsAutomatically` — `true`이면 사용자 제스처 없이 JavaScript에서 새 창을 열 수 있습니다(기본값: `false`).
- `MinimumFontSize` — 포인트 단위의 최소 글꼴 크기입니다. 설정하려면 `optional.NewVar(12.0)`을 사용하세요. 설정하지 않으면 WebKit 기본값이 유지됩니다.
- `ApplicationNameForUserAgent` — WebKit 사용자 에이전트 문자열의 애플리케이션 이름 접미사를 재정의합니다. 사이트에서 기본 `"wails.io"` 식별자를 거부할 때(예: YouTube 임베드) 유용합니다. 기본값을 유지하려면 비워 두세요.
- `EnableAutoplayWithoutUserAction` — `true`이면 사용자 제스처 없이 오디오와 비디오를 자동 재생할 수 있습니다. `WKWebViewConfiguration.mediaTypesRequiringUserActionForPlayback = WKAudiovisualMediaTypeNone`에 매핑됩니다(기본값: `false`).

### Windows 옵션(창별)

창별 구조체는 `application.WindowsWindow`입니다. `WindowsOptions`이 **아닙니다**(이는 *애플리케이션* 수준 구조체입니다).

```go
Windows: application.WindowsWindow{
    DisableIcon:                       false,
    DisableMenu:                       false,
    BackdropType:                      application.Auto,
    CustomTheme:                       application.ThemeSettings{},
    DisableFramelessWindowDecorations: false,
    NonClientRegionSupport:            false,
    WebView2CompositionHosting:        false,
},
```

**DisableIcon**(`bool`)

- 제목 표시줄에서 아이콘을 제거합니다.

**DisableMenu**(`bool`)

- 창의 메뉴 표시줄을 비활성화합니다. `true`이면 메뉴 표시줄을 구성했더라도 창에 표시되지 않습니다.
- 기본값: `false`

**BackdropType**(`BackdropType`)

- `application.Auto` - 시스템 기본값
- `application.None` - 배경 효과 없음
- `application.Mica` - Mica 재질(Windows 11)
- `application.Acrylic` - Acrylic 재질(Windows 11)
- `application.Tabbed` - Tabbed 재질(Windows 11)

`WindowsBackdropTypeMica` 형식의 상수는 없습니다. `application.Mica` 등을 사용하세요.

**CustomTheme**(`ThemeSettings`)

- 포인터가 아닌 값입니다. 창 테두리, 제목 표시줄의 텍스트와 배경, 메뉴 표시줄에 사용할 사용자 지정 다크/라이트 모드 색상입니다.

**DisableFramelessWindowDecorations**(`bool`)

- 기본 프레임리스 장식(Aero 그림자, 둥근 모서리)을 비활성화합니다.

**NonClientRegionSupport**(`bool`)

- 프레임리스 사용자 지정 제목 표시줄에 WebView2 네이티브 `app-region: drag` / `app-region: no-drag` 지원을 활성화합니다.
- 이는 단순한 네이티브 앱 드래그만 지원합니다. 네이티브 사용자 지정 캡션 버튼 동작이나 사용자 지정 최대화 버튼의 Windows 11 Snap Assist / Snap Layouts는 제공하지 않습니다.

**WebView2CompositionHosting**(`bool`)

- 사용자 지정 최대화 버튼의 Windows 11 Snap Assist / Snap Layouts를 비롯하여, Windows 네이티브 동작을 제공하는 사용자 지정 캡션 버튼에 대해 Wails에서 관리하는 `--wails-non-client-region` 지원을 활성화합니다.
- 실험적 기능입니다. 기본 HWND 호스팅 컨트롤러 대신 `ICoreWebView2CompositionController` 및 DirectComposition을 통해 WebView2를 호스팅합니다.
- 창에 WebView2 네이티브 `app-region` 지원과 Wails에서 관리하는 사용자 지정 캡션 버튼 영역이 모두 필요한 경우 `NonClientRegionSupport`과 함께 사용할 수 있습니다.

**예:**

```go
Windows: application.WindowsWindow{
    BackdropType: application.Mica,
    DisableIcon:  true,
},
```

**예 - 사용자 지정 Windows 제목 표시줄 영역:**

```go
Windows: application.WindowsWindow{
    NonClientRegionSupport:    true,
    WebView2CompositionHosting: true,
},
```

자세한 동작, 장단점 및 이에 맞는 CSS는 [프레임리스 Windows 창](/features/windows/frameless/#native-non-client-regions-on-windows)을 참조하세요.

### Linux 옵션(창별)

창별 구조체는 `application.LinuxWindow`입니다. `LinuxOptions`이 **아닙니다**.

```go
Linux: application.LinuxWindow{
    Icon:                []byte{/* PNG data */},
    WindowIsTranslucent: false,
},
```

**Icon**(`[]byte`)

- 창 아이콘(PNG 형식)입니다.

**WindowIsTranslucent**(`bool`)

- 컴포지터 지원이 필요합니다.

**예:**

```go
//go:embed icon.png
var icon []byte

Linux: application.LinuxWindow{
    Icon: icon,
},
```

## 애플리케이션 수준 Windows 옵션

일부 Windows 전용 옵션은 창별로 설정하지 않고 애플리케이션 수준에서 설정해야 합니다. WebView2가 사용자 데이터 경로별로 하나의 브라우저 환경을 공유하기 때문입니다.

### 브라우저 플래그

WebView2 브라우저 플래그는 애플리케이션의 <strong>모든 창</strong>에 적용되는 실험적 기능과 동작을 제어합니다. 이 플래그는 `application.Options.Windows`에서 설정해야 합니다.

```go
app := application.New(application.Options{
    Name: "My App",
    Windows: application.WindowsOptions{
        // Enable experimental WebView2 features
        EnabledFeatures: []string{
            "msWebView2EnableDraggableRegions",
        },

        // Disable specific features
        DisabledFeatures: []string{
            "msSmartScreenProtection",  // Always disabled by Wails
        },

        // Additional Chromium command-line arguments
        AdditionalBrowserArgs: []string{
            "--disable-gpu",
            "--remote-debugging-port=9222",
        },
    },
})
```

**EnabledFeatures** (`[]string`)

- 활성화할 WebView2 기능 플래그 목록
- 사용 가능한 플래그는 [WebView2 브라우저 플래그](https://learn.microsoft.com/en-us/microsoft-edge/webview2/concepts/webview-features-flags)를 참조하세요.
- 예: `"msWebView2EnableDraggableRegions"`

**DisabledFeatures** (`[]string`)

- 비활성화할 WebView2 기능 플래그 목록
- Wails는 `msSmartScreenProtection`을 자동으로 비활성화합니다.
- 예: `"msExperimentalFeature"`

**AdditionalBrowserArgs** (`[]string`)

- 브라우저 프로세스에 전달할 Chromium 명령줄 인수
- `--` 접두사를 반드시 포함해야 합니다(예: `"--remote-debugging-port=9222"`).
- 사용 가능한 인수는 [Chromium 명령줄 스위치](https://peter.sh/experiments/chromium-command-line-switches/)를 참조하세요.

@note{type="caution" title="중요"}
WebView2는 사용자 데이터 경로별로 하나의 브라우저 환경을 공유하므로 이 플래그는 모든 창에 전역으로 적용됩니다. 창마다 서로 다른 브라우저 플래그를 사용할 수 없습니다.

@end

**전체 예제:**

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "My App",
        Windows: application.WindowsOptions{
            // Enable draggable regions feature
            EnabledFeatures: []string{
                "msWebView2EnableDraggableRegions",
            },
            // Enable remote debugging
            AdditionalBrowserArgs: []string{
                "--remote-debugging-port=9222",
            },
        },
    })

    // All windows will use the browser flags configured above
    window := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:  "Main Window",
        Width:  1024,
        Height: 768,
    })

    window.Show()
    app.Run()
}
```

## 전체 예제

다음은 프로덕션 환경에서 바로 사용할 수 있는 창 구성입니다.

```go
package main

import (
    _ "embed"
    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var assets embed.FS

//go:embed icon.png
var icon []byte

func main() {
    app := application.New(application.Options{
        Name: "My Application",
    })

    window := app.Window.NewWithOptions(application.WebviewWindowOptions{
        // Identity
        Name:  "main-window",
        Title: "My Application",

        // Size and Position
        Width:     1200,
        Height:    800,
        MinWidth:  800,
        MinHeight: 600,

        // Initial State
        StartState: application.WindowStateNormal,

        // Appearance
        BackgroundColour: application.NewRGB(255, 255, 255),

        // Platform-Specific (per-window structs)
        Mac: application.MacWindow{
            TitleBar: application.MacTitleBar{
                AppearsTransparent: true,
            },
            Backdrop: application.MacBackdropTranslucent,
        },

        Windows: application.WindowsWindow{
            BackdropType: application.Mica,
            DisableIcon:  false,
        },

        Linux: application.LinuxWindow{
            Icon: icon,
        },
    })

    window.Center()
    window.Show()

    app.Run()
}
```

프런트엔드 애셋은 창별로 제공되지 않고 애플리케이션 수준(`application.Options.Assets`)에서 제공됩니다.

## 다음 단계

- [창 기본 사항](/features/windows/basics/) - 창 생성 및 제어
- [여러 창](/features/windows/multiple/) - 다중 창 패턴
- [프레임 없는 창](/features/windows/frameless/) - 사용자 지정 창 장식
- [창 이벤트](/features/windows/events/) - 수명 주기 이벤트

---

**궁금한 점이 있으신가요?** [Discord](https://discord.gg/JDdSxwjhGf)에서 질문하거나 [예제](https://github.com/wailsapp/wails/tree/master/v3/examples)를 확인하세요.
