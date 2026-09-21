---
title: "시스템 트레이 메뉴"
description: "애플리케이션에 시스템 트레이(알림 영역) 통합 기능 추가"
slug: "features/menus/systray"
sourcePath: "features/menus/systray.md"
---

## 시스템 트레이 메뉴

Wails는 모든 플랫폼에서 작동하는 <strong>통합 시스템 트레이 API</strong>를 제공합니다. 메뉴가 있는 트레이 아이콘을 만들고, 창을 연결하고, 클릭을 처리할 수 있으며 백그라운드 애플리케이션, 서비스 및 빠른 실행 유틸리티에서 각 플랫폼의 네이티브 동작을 지원합니다.

![macOS 메뉴 막대에서 연 Wails 시스템 트레이 메뉴](/assets/screenshots/systray-menu-macos.png)

macOS에서 Wails 시스템 트레이 항목은 메뉴 막대에 표시되며 네이티브 메뉴를 엽니다. 이 예제에는 비활성화 항목, 체크박스 항목, 라디오 항목, 하위 메뉴 및 작업 항목이 포함되어 있습니다.

## 빠른 시작

```go
package main

import (
    _ "embed"
    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets/icon.png
var icon []byte

func main() {
    app := application.New(application.Options{
        Name: "Tray App",
    })

    // Create system tray
    systray := app.SystemTray.New()
    systray.SetIcon(icon)
    systray.SetLabel("My App")

    // Add menu
    menu := app.NewMenu()
    menu.Add("Show").OnClick(func(ctx *application.Context) {
        // Show main window
    })
    menu.Add("Quit").OnClick(func(ctx *application.Context) {
        app.Quit()
    })
    systray.SetMenu(menu)

    // Create hidden window
    window := app.Window.New()
    window.Hide()

    app.Run()
}
```

**결과:** 모든 플랫폼에서 메뉴가 있는 시스템 트레이 아이콘이 표시됩니다.

## 시스템 트레이 만들기

### 기본 시스템 트레이

```go
// Create system tray
systray := app.SystemTray.New()

// Set icon
systray.SetIcon(iconBytes)

// Set label (macOS) / tooltip (Windows)
systray.SetLabel("My Application")
```

### 아이콘 추가

아이콘은 임베드하는 것이 좋습니다:

```go
import _ "embed"

//go:embed assets/icon.png
var icon []byte

//go:embed assets/icon-dark.png
var iconDark []byte

func main() {
    app := application.New(application.Options{
        Name: "My App",
    })

    systray := app.SystemTray.New()
    systray.SetIcon(icon)
    systray.SetDarkModeIcon(iconDark)  // Windows and macOS dark mode
    
    app.Run()
}
```

**아이콘 요구 사항:**

| 플랫폼 | 크기 | 형식 | 참고 |
| --- | --- | --- | --- |
| **Windows** | 16x16 또는 32x32 | PNG, ICO | 알림 영역 |
| **macOS** | 18x18~22x22 | PNG | 메뉴 막대, 템플릿 권장 |
| **Linux** | 22x22~48x48 | PNG, SVG | 데스크톱 환경에 따라 다름 |

### 템플릿 아이콘(macOS)

템플릿 아이콘은 라이트/다크 모드에 자동으로 맞춰집니다.

```go
systray.SetTemplateIcon(iconBytes)
```

**템플릿 아이콘 지침:**

- 검은색과 투명색만 사용하세요.
- 다크 모드에서는 검은색이 흰색으로 바뀝니다.
- 파일 이름에 `Template` 접미사를 붙이세요: `iconTemplate.png`
- [디자인 가이드](https://bjango.com/articles/designingmenubarextras/)

## 메뉴 추가하기

시스템 트레이 메뉴는 애플리케이션 메뉴와 동일하게 작동합니다.

```go
menu := app.NewMenu()

// Add items
menu.Add("Open").OnClick(func(ctx *application.Context) {
    showMainWindow()
})

menu.AddSeparator()

menu.AddCheckbox("Start at Login", false).OnClick(func(ctx *application.Context) {
    enabled := ctx.ClickedMenuItem().Checked()
    setStartAtLogin(enabled)
})

menu.AddSeparator()

menu.Add("Quit").OnClick(func(ctx *application.Context) {
    app.Quit()
})

// Set menu
systray.SetMenu(menu)
```

<strong>모든 메뉴 항목 유형</strong>은 [메뉴 레퍼런스](/features/menus/reference/)를 참조하세요.

## 창 연결하기

자동으로 표시하거나 숨기려면 창을 트레이 아이콘에 연결하세요.

```go
// Create window
window := app.Window.New()

// Attach to tray
systray.AttachWindow(window)

// Configure behaviour — these are setters that return the receiver for chaining.
systray.WindowOffset(10)                          // Pixels from tray icon
systray.WindowDebounce(200 * time.Millisecond)    // Click debounce
```

**동작:**

- 창은 숨겨진 상태로 시작됩니다.
- **트레이 아이콘 왼쪽 클릭** → 창 표시 여부 전환
- **트레이 아이콘 오른쪽 클릭** → 메뉴 표시(설정된 경우)
- 창은 트레이 아이콘 근처에 배치됩니다.

**예제: 팝업 창**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:           "Quick Access",
    Width:           300,
    Height:          400,
    Frameless:       true, // No title bar
    AlwaysOnTop:     true, // Stay on top
    HideOnFocusLost: true, // Dismiss when another window receives focus
    HideOnEscape:    true, // Dismiss when the user presses Escape
})

systray.AttachWindow(window)
systray.WindowOffset(5)
```

`HideOnFocusLost`은 Windows, macOS 및 클릭하여 포커스를 설정하는 Linux 데스크톱의 트레이 팝업에 유용합니다. Wails는 포커스가 마우스를 따라가는 Linux 환경(일반적인 Hyprland, Sway 및 i3 설정 포함)에서는 이 동작을 비활성화합니다. 그렇지 않으면 팝업에서 마우스 포인터가 벗어날 때 팝업을 사용하기도 전에 숨겨질 수 있기 때문입니다. 이러한 환경에서도 `HideOnEscape`은 계속 사용할 수 있습니다.

위의 왼쪽 및 오른쪽 클릭 동작은 지능형 기본값입니다. 명시적인 `OnClick` 또는 `OnRightClick` 핸들러를 지정하면 해당 기본 동작을 대체합니다. 플랫폼 확인 및 예외 사례는 [수동 systray 테스트 모음](https://github.com/wailsapp/wails/tree/master/v3/test/manual/systray)과 [systray 스트레스 테스트 예제](https://github.com/wailsapp/wails/tree/master/v3/examples/systray-stress)를 참조하세요.

## 클릭 핸들러

트레이 아이콘 클릭을 처리하세요.

```go
systray := app.SystemTray.New()

// Left click
systray.OnClick(func() {
    fmt.Println("Tray icon clicked")
})

// Right click
systray.OnRightClick(func() {
    fmt.Println("Tray icon right-clicked")
})

// Double click
systray.OnDoubleClick(func() {
    fmt.Println("Tray icon double-clicked")
})

// Mouse enter/leave
systray.OnMouseEnter(func() {
    fmt.Println("Mouse entered tray icon")
})

systray.OnMouseLeave(func() {
    fmt.Println("Mouse left tray icon")
})
```

**플랫폼 지원:**

| 이벤트 | Windows | macOS | Linux |
| --- | --- | --- | --- |
| OnClick | ✅ | ✅ | ✅ |
| OnRightClick | ✅ | ✅ | ✅ |
| OnDoubleClick | ✅ | ✅ | ⚠️ 환경에 따라 다름 |
| OnMouseEnter | ✅ | ✅ | ⚠️ 환경에 따라 다름 |
| OnMouseLeave | ✅ | ✅ | ⚠️ 환경에 따라 다름 |

## 동적 업데이트

트레이 아이콘과 메뉴를 동적으로 업데이트하세요.

### 아이콘 변경하기

```go
var isActive bool

func updateTrayIcon() {
    if isActive {
        systray.SetIcon(activeIcon)
        systray.SetLabel("Active")
    } else {
        systray.SetIcon(inactiveIcon)
        systray.SetLabel("Inactive")
    }
}
```

### 메뉴 업데이트하기

```go
var isPaused bool

pauseMenuItem := menu.Add("Pause")

pauseMenuItem.OnClick(func(ctx *application.Context) {
    isPaused = !isPaused
    
    if isPaused {
        pauseMenuItem.SetLabel("Resume")
    } else {
        pauseMenuItem.SetLabel("Pause")
    }
    
    menu.Update()  // Important!
})
```

@note{type="caution" title="항상 Update() 호출"}
메뉴 상태를 변경한 후 <strong>`menu.Update()`</strong>을 호출하세요. [메뉴 레퍼런스](/features/menus/reference/#enabled-state)를 참조하세요.

@end

### 메뉴 다시 빌드하기

대규모로 변경하려면 메뉴 전체를 다시 빌드하세요:

```go
func rebuildTrayMenu(status string) {
    menu := app.NewMenu()
    
    // Status-specific items
    switch status {
    case "syncing":
        menu.Add("Syncing...").SetEnabled(false)
        menu.Add("Pause Sync").OnClick(pauseSync)
    case "synced":
        menu.Add("Up to date ✓").SetEnabled(false)
        menu.Add("Sync Now").OnClick(startSync)
    case "error":
        menu.Add("Sync Error").SetEnabled(false)
        menu.Add("Retry").OnClick(retrySync)
    }
    
    menu.AddSeparator()
    menu.Add("Quit").OnClick(func(ctx *application.Context) {
        app.Quit()
    })
    
    systray.SetMenu(menu)
}
```

## 플랫폼별 기능

@tabs{sync-key="platform"}
[macOS]
**메뉴 막대 통합:**

```go
// Set label (appears next to icon)
systray.SetLabel("My App")

// Use template icon (adapts to dark mode)
systray.SetTemplateIcon(iconBytes)

// Set icon position — uses AppKit NSImage placement constants.
systray.SetIconPosition(application.NSImageRight)
```

**아이콘 위치**(`NSImagePosition` 미러링):

- `application.NSImageLeft` - 레이블 왼쪽에 아이콘을 표시합니다.
- `application.NSImageRight` - 레이블 오른쪽에 아이콘을 표시합니다.
- `application.NSImageOnly` - 레이블 없이 아이콘만 표시합니다.
- `application.NSImageNone` - 아이콘 없이 레이블만 표시합니다.

**권장 사례:**

- 템플릿 아이콘(검은색 + 투명 배경)을 사용하세요
- 레이블을 짧게 유지하세요(3-5자)
- Retina 디스플레이에는 18x18~22x22픽셀을 사용하세요
- 라이트 모드와 다크 모드에서 모두 테스트하세요

[Windows]
**알림 영역 통합:**

```go
// Set tooltip (appears on hover)
systray.SetTooltip("My Application")

// Or use SetLabel (same as tooltip on Windows)
systray.SetLabel("My Application")

// Show/Hide functionality (fully functional)
systray.Show()  // Show tray icon
systray.Hide()  // Hide tray icon
```

**아이콘 요구 사항:**

- 16x16 또는 32x32픽셀
- PNG 또는 ICO 형식
- 투명 배경

**도구 설명 제한:**

- 최대 127개의 UTF-16 문자
- 이보다 긴 도구 설명은 잘립니다
- 최상의 사용자 경험을 위해 간결하게 작성하세요

**플랫폼 기능:**

- Windows Explorer가 다시 시작되어도 트레이 아이콘이 유지됩니다
- Show() 및 Hide() 메서드가 완전히 작동합니다
- 수명 주기를 올바르게 관리합니다

**권장 사례:**

- 고DPI 디스플레이에는 32x32를 사용하세요
- 도구 설명을 127자 미만으로 유지하세요
- 여러 Windows 버전에서 테스트하세요
- 알림 영역 오버플로를 고려하세요
- 조건에 따라 트레이 표시 여부를 제어하려면 Show/Hide를 사용하세요

[Linux]
**시스템 트레이 통합:**

StatusNotifierItem 사양을 사용합니다(대부분의 최신 데스크톱 환경).

```go
systray.SetIcon(iconBytes)
systray.SetLabel("My App")
```

**데스크톱 환경 지원:**

- **GNOME**: 상단 표시줄(확장 기능 필요)
- **KDE Plasma**: 시스템 트레이
- **XFCE**: 알림 영역
- **기타**: 환경에 따라 다름

**권장 사례:**

- 22x22 또는 24x24픽셀을 사용하세요
- SVG 아이콘은 크기 조정 시 품질이 더 잘 유지됩니다
- 대상 데스크톱 환경에서 테스트하세요
- 지원되지 않는 데스크톱 환경을 위한 대체 수단을 제공하세요

@end

## 전체 예제

다음은 프로덕션 환경에서 사용할 수 있는 시스템 트레이 애플리케이션입니다:

```go
package main

import (
    _ "embed"
    "fmt"
    "time"
    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets/icon.png
var icon []byte

//go:embed assets/icon-active.png
var iconActive []byte

type TrayApp struct {
    app     *application.App
    systray *application.SystemTray
    window  *application.WebviewWindow
    menu    *application.Menu
    isActive bool
}

func main() {
    app := application.New(application.Options{
        Name: "Tray Application",
        Mac: application.MacOptions{
            ApplicationShouldTerminateAfterLastWindowClosed: false,
        },
    })

    trayApp := &TrayApp{app: app}
    trayApp.setup()

    app.Run()
}

func (t *TrayApp) setup() {
    // Create system tray
    t.systray = t.app.SystemTray.New()
    t.systray.SetIcon(icon)
    t.systray.SetLabel("Inactive")
    
    // Create menu
    t.createMenu()
    
    // Create window (hidden by default)
    t.window = t.app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:  "Tray Application",
        Width:  400,
        Height: 600,
        Hidden: true,
    })
    
    // Attach window to tray
    t.systray.AttachWindow(t.window)
    t.systray.WindowOffset(10)
    
    // Handle tray clicks
    t.systray.OnRightClick(func() {
        t.systray.OpenMenu()
    })
    
    // Start background task
    go t.backgroundTask()
}

func (t *TrayApp) createMenu() {
    t.menu = t.app.NewMenu()
    
    // Status item (disabled)
    statusItem := t.menu.Add("Status: Inactive")
    statusItem.SetEnabled(false)
    
    t.menu.AddSeparator()
    
    // Toggle active
    t.menu.Add("Start").OnClick(func(ctx *application.Context) {
        t.toggleActive()
    })
    
    // Show window
    t.menu.Add("Show Window").OnClick(func(ctx *application.Context) {
        t.window.Show()
        t.window.Focus()
    })
    
    t.menu.AddSeparator()
    
    // Settings
    t.menu.AddCheckbox("Start at Login", false).OnClick(func(ctx *application.Context) {
        enabled := ctx.ClickedMenuItem().Checked()
        t.setStartAtLogin(enabled)
    })
    
    t.menu.AddSeparator()
    
    // Quit
    t.menu.Add("Quit").OnClick(func(ctx *application.Context) {
        t.app.Quit()
    })
    
    t.systray.SetMenu(t.menu)
}

func (t *TrayApp) toggleActive() {
    t.isActive = !t.isActive
    t.updateTray()
}

func (t *TrayApp) updateTray() {
    if t.isActive {
        t.systray.SetIcon(iconActive)
        t.systray.SetLabel("Active")
    } else {
        t.systray.SetIcon(icon)
        t.systray.SetLabel("Inactive")
    }
    
    // Rebuild menu with new status
    t.createMenu()
}

func (t *TrayApp) backgroundTask() {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        if t.isActive {
            fmt.Println("Background task running...")
            // Do work
        }
    }
}

func (t *TrayApp) setStartAtLogin(enabled bool) {
    // Implementation varies by platform
    fmt.Printf("Start at login: %v\n", enabled)
}
```

## 표시 여부 제어

트레이 아이콘을 동적으로 표시하거나 숨기세요:

```go
// Hide tray icon
systray.Hide()

// Show tray icon
systray.Show()
```

`IsVisible()` getter는 없습니다. 표시 여부를 확인해야 한다면 애플리케이션 상태에서 직접 추적하세요.

**플랫폼 지원:**

| 플랫폼 | Hide() | Show() | 참고 |
| --- | --- | --- | --- |
| **Windows** | ✅ | ✅ | 완전히 작동함 - 알림 영역에서 아이콘이 나타나거나 사라집니다 |
| **macOS** | ✅ | ✅ | 메뉴 막대 항목이 표시되거나 숨겨집니다 |
| **Linux** | ✅ | ✅ | 데스크톱 환경에 따라 다릅니다 |

**사용 사례:**

- 사용자 환경 설정에 따라 트레이 아이콘을 일시적으로 숨기기
- 필요할 때만 트레이 아이콘이 나타나는 헤드리스 모드
- 애플리케이션 상태에 따라 표시 여부 전환

**예제 - 조건부 트레이 표시:**

```go
func (t *TrayApp) setTrayVisibility(visible bool) {
    if visible {
        t.systray.Show()
    } else {
        t.systray.Hide()
    }
}

// Show tray only when updates are available
func (t *TrayApp) checkForUpdates() {
    if hasUpdates {
        t.systray.Show()
        t.systray.SetLabel("Update Available")
    } else {
        t.systray.Hide()
    }
}
```

## 정리

작업을 마치면 트레이 아이콘을 제거하세요:

```go
// In OnShutdown
app := application.New(application.Options{
    OnShutdown: func() {
        if systray != nil {
            systray.Destroy()
        }
    },
})
```

**중요:** 리소스를 해제하려면 종료 시 항상 시스템 트레이를 제거하세요.

## 권장 사례

### ✅ 권장 사항

- **macOS에서는 템플릿 아이콘을 사용하세요** - 다크 모드에 맞게 조정됩니다
- **레이블을 짧게 유지하세요** - 최대 3-5자
- **Windows에서는 도구 설명을 제공하세요** - 사용자가 앱을 식별하는 데 도움이 됩니다
- **모든 플랫폼에서 테스트하세요** - 동작이 플랫폼마다 다릅니다
- **클릭을 적절히 처리하세요** - 기본 동작에는 왼쪽 클릭을, 메뉴에는 오른쪽 클릭을 사용하세요
- **상태에 따라 아이콘을 업데이트하세요** - 시각적 피드백은 중요합니다
- **종료할 때 시스템 트레이를 제거하세요** - 리소스를 해제할 수 있습니다

### ❌ 하지 말아야 할 사항

- **큰 아이콘을 사용하지 마세요** - 플랫폼 지침을 따르세요
- **긴 레이블을 사용하지 마세요** - 일부가 잘립니다
- **다크 모드를 잊지 마세요** - Windows와 macOS의 다크 모드에서 테스트하세요
- **클릭 핸들러를 블로킹하지 마세요** - 빠르게 실행되도록 유지하세요
- **menu.Update()를 잊지 마세요** - 메뉴 상태를 변경한 후 호출하세요
- **시스템 트레이가 지원된다고 단정하지 마세요** - 일부 Linux 데스크톱 환경에서는 지원하지 않습니다

## 문제 해결

### 트레이 아이콘이 나타나지 않음

**가능한 원인:**

1. 지원되지 않는 아이콘 형식
2. 아이콘 크기가 너무 크거나 작음
3. 시스템 트레이가 지원되지 않음(Linux)

**해결 방법:**

`SystemTraySupported()` 헬퍼는 없습니다. 대신 트레이를 생성하고 플랫폼을 확인한 후, 지원되지 않는 경우에도 정상적으로 동작하도록 처리하세요:

```go
// Probe support: on Linux without a notification-area extension, the tray
// will simply not appear. Defensive code can fall back to window-only mode
// based on runtime.GOOS or after a short timeout if no tray events arrive.
systray := app.SystemTray.New()
systray.SetIcon(iconBytes)
```

### macOS에서 아이콘이 올바르게 표시되지 않음

**원인:** 템플릿 아이콘을 사용하지 않음

**해결 방법:**

```go
// Use template icon
systray.SetTemplateIcon(iconBytes)

// Or design icon as template (black + transparent)
```

### 메뉴가 업데이트되지 않음

**원인:** `menu.Update()` 호출을 잊음

**해결 방법:**

```go
menuItem.SetLabel("New Label")
menu.Update()  // Add this!
```

## 다음 단계

@cards{cols="2"}
📖 메뉴 참조
메뉴 항목 유형과 속성에 대한 전체 참조입니다.

[자세히 알아보기 →](/features/menus/reference/)

---
☰ 애플리케이션 메뉴
애플리케이션 메뉴 모음을 만드세요.

[자세히 알아보기 →](/features/menus/application/)

---
◆ 컨텍스트 메뉴
오른쪽 클릭 컨텍스트 메뉴를 만드세요.

[자세히 알아보기 →](/features/menus/context/)

---
📖 시스템 트레이 예제
완전한 시스템 트레이 애플리케이션을 살펴보세요.

[자세히 알아보기 →](https://github.com/wailsapp/wails/tree/master/v3/examples/systray-basic)

@end

---

**질문이 있으신가요?** [Discord](https://discord.gg/JDdSxwjhGf)에서 질문하거나 [시스템 트레이 예제](https://github.com/wailsapp/wails/tree/master/v3/examples/systray-basic)를 확인하세요.
