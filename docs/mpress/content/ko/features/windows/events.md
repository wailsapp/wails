---
title: "창 이벤트"
description: "창 수명 주기 및 상태 변경 이벤트 처리"
slug: "features/windows/events"
sourcePath: "features/windows/events.md"
---

## 창 이벤트

Wails는 단일 API를 통해 창 수명 주기 및 상태 변경 이벤트를 전달합니다. 기본 동작을 막지 않고 이벤트를 관찰하는 리스너에는 `OnWindowEvent`을 사용하고, 기본 동작을 막을 수 있는 취소 가능 훅에는 `RegisterHook`을 사용합니다.

```go
import "github.com/wailsapp/wails/v3/pkg/events"

// Listener — observes the event, cannot cancel it.
window.OnWindowEvent(events.Common.WindowDidResize, func(e *application.WindowEvent) { /* ... */ })

// Hook — runs before listeners; can call e.Cancel() to suppress the default action.
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    if hasUnsavedChanges() {
        e.Cancel() // prevent the window from closing
    }
})
```

두 호출 모두 핸들러를 제거하기 위해 호출할 수 있는 `unsubscribe func()`을 반환합니다.

크로스 플랫폼 이벤트는 `events.Common.*`에 있습니다. 플랫폼별 이벤트는 `events.Mac.*`, `events.Windows.*` 및 `events.Linux.*`에 있습니다. 전체 목록은 `v3/pkg/events/events.go`에 생성됩니다.

## 수명 주기 이벤트

### 창 생성

`app.Window.OnCreate`을 사용하여 창이 생성될 때마다 콜백을 실행합니다.

```go
app.Window.OnCreate(func(window application.Window) {
    fmt.Printf("Window created: %s (ID: %d)\n", window.Name(), window.ID())

    // Configure all new windows
    window.SetMinSize(400, 300)

    // Register a hook for this window
    window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
        if !confirmClose() {
            e.Cancel()
        }
    })
})
```

콜백은 인터페이스인 `application.Window`을 받습니다. `OnCreate` 콜백은 창의 런타임이 초기화된 후 각 창마다 한 번씩 호출됩니다.

### WindowClosing

사용자가 창을 닫으려고 할 때(X 클릭, ⌘W, Alt+F4 등) 발생합니다.

**훅**(`RegisterHook`)을 사용하세요. 훅만 닫기를 취소할 수 있습니다. 리스너는 이벤트를 감지하지만 닫기를 막을 수 없습니다.

```go
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    if hasUnsavedChanges() {
        // Cancel the close — the window stays open.
        e.Cancel()
    }
})
```

**중요:**

- `WindowClosing`은 사용자가 창 닫기를 시도할 때 전달됩니다.
- `RegisterHook` 콜백은 `e.Cancel()`을 호출하여 창을 열린 상태로 유지할 수 있습니다.
- `WindowClosing`에 대한 `OnWindowEvent` 콜백은 옵서버입니다. 콜백이 실행되기는 하지만 취소할 수는 없습니다.

**프로그래밍 방식으로 닫기:** `window.Close()`을 호출하세요. `window.Destroy()`은 없습니다.

### WindowRuntimeReady

창 내부 런타임의 초기화가 완료되면 발생합니다. 이 시점부터 창의 JS 컨텍스트를 안전하게 호출할 수 있습니다:

```go
window.OnWindowEvent(events.Common.WindowRuntimeReady, func(e *application.WindowEvent) {
    window.EmitEvent("app-ready", nil)
})
```

## 포커스 이벤트

### WindowFocus

창이 포커스를 얻으면 호출됩니다.

```go
window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
    fmt.Println("Window gained focus")
    updateTitleBar(true)
    app.Event.Emit("window-focused", window.ID())
})
```

### WindowLostFocus

창이 포커스를 잃으면 호출됩니다.

```go
window.OnWindowEvent(events.Common.WindowLostFocus, func(e *application.WindowEvent) {
    fmt.Println("Window lost focus")
    updateTitleBar(false)
    saveCurrentState()
})
```

**예: 포커스를 인식하는 UI:**

```go
type FocusAwareWindow struct {
    window  *application.WebviewWindow
    focused bool
}

func (fw *FocusAwareWindow) Setup() {
    fw.window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
        fw.focused = true
        fw.updateAppearance()
    })

    fw.window.OnWindowEvent(events.Common.WindowLostFocus, func(e *application.WindowEvent) {
        fw.focused = false
        fw.updateAppearance()
    })
}

func (fw *FocusAwareWindow) updateAppearance() {
    if fw.focused {
        fw.window.EmitEvent("update-theme", "active")
    } else {
        fw.window.EmitEvent("update-theme", "inactive")
    }
}
```

## 상태 변경 이벤트

### WindowMinimise / WindowUnMinimise

```go
window.OnWindowEvent(events.Common.WindowMinimise, func(e *application.WindowEvent) {
    pauseRendering()
    saveWindowState()
})

window.OnWindowEvent(events.Common.WindowUnMinimise, func(e *application.WindowEvent) {
    resumeRendering()
    refreshContent()
})
```

### WindowMaximise / WindowUnMaximise

```go
window.OnWindowEvent(events.Common.WindowMaximise, func(e *application.WindowEvent) {
    window.EmitEvent("layout-mode", "maximised")
})

window.OnWindowEvent(events.Common.WindowUnMaximise, func(e *application.WindowEvent) {
    window.EmitEvent("layout-mode", "normal")
})
```

### WindowFullscreen / WindowUnFullscreen

```go
window.OnWindowEvent(events.Common.WindowFullscreen, func(e *application.WindowEvent) {
    window.EmitEvent("chrome-visibility", false)
    window.EmitEvent("layout-mode", "fullscreen")
})

window.OnWindowEvent(events.Common.WindowUnFullscreen, func(e *application.WindowEvent) {
    window.EmitEvent("chrome-visibility", true)
    window.EmitEvent("layout-mode", "normal")
})
```

`window.Fullscreen()` / `window.UnFullscreen()` / `window.ToggleFullscreen()`을 사용하여 프로그래밍 방식으로 전체 화면 모드에 진입하거나 종료합니다. `SetFullscreen(bool)`은 없습니다. `window.IsFullscreen() bool`을 사용하여 상태를 조회합니다.

## 위치 및 크기 이벤트

### WindowDidMove

```go
window.OnWindowEvent(events.Common.WindowDidMove, func(e *application.WindowEvent) {
    x, y := window.Position()
    fmt.Printf("Window moved to: %d, %d\n", x, y)
    saveWindowPosition(x, y)
})
```

### WindowDidResize

```go
window.OnWindowEvent(events.Common.WindowDidResize, func(e *application.WindowEvent) {
    width, height := window.Size()
    fmt.Printf("Window resized to: %dx%d\n", width, height)
    saveWindowSize(width, height)
    window.EmitEvent("window-size", map[string]int{
        "width":  width,
        "height": height,
    })
})
```

`WindowDidResize` 및 `WindowDidMove` 이벤트 자체에는 좌표가 포함되지 않습니다. 콜백 내에서 `window.Size()` / `window.Position()`을 통해 창의 정보를 조회하세요.

**예: 반응형 레이아웃:**

```go
window.OnWindowEvent(events.Common.WindowDidResize, func(e *application.WindowEvent) {
    width, _ := window.Size()
    var layout string
    switch {
    case width < 600:
        layout = "compact"
    case width < 1200:
        layout = "normal"
    default:
        layout = "wide"
    }
    window.EmitEvent("layout-changed", layout)
})
```

## 전체 예제

모든 이벤트 처리를 갖춘 프로덕션용 창의 예제입니다.

```go
package main

import (
    "encoding/json"
    "fmt"
    "os"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)

type WindowState struct {
    X          int  `json:"x"`
    Y          int  `json:"y"`
    Width      int  `json:"width"`
    Height     int  `json:"height"`
    Maximised  bool `json:"maximised"`
    Fullscreen bool `json:"fullscreen"`
}

type ManagedWindow struct {
    app    *application.App
    window *application.WebviewWindow
    state  WindowState
    dirty  bool
}

func main() {
    app := application.New(application.Options{
        Name: "Event Demo",
    })

    mw := &ManagedWindow{app: app}
    mw.CreateWindow()
    mw.LoadState()
    mw.SetupEventHandlers()

    app.Run()
}

func (mw *ManagedWindow) CreateWindow() {
    mw.window = mw.app.Window.NewWithOptions(application.WebviewWindowOptions{
        Name:   "main",
        Title:  "Event Demo",
        Width:  800,
        Height: 600,
    })
}

func (mw *ManagedWindow) SetupEventHandlers() {
    // Focus events
    mw.window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
        mw.window.EmitEvent("focus-state", true)
    })

    mw.window.OnWindowEvent(events.Common.WindowLostFocus, func(e *application.WindowEvent) {
        mw.window.EmitEvent("focus-state", false)
    })

    // State change events
    mw.window.OnWindowEvent(events.Common.WindowMinimise, func(e *application.WindowEvent) {
        mw.SaveState()
    })

    mw.window.OnWindowEvent(events.Common.WindowMaximise, func(e *application.WindowEvent) {
        mw.state.Maximised = true
        mw.dirty = true
    })

    mw.window.OnWindowEvent(events.Common.WindowUnMaximise, func(e *application.WindowEvent) {
        mw.state.Maximised = false
        mw.dirty = true
    })

    mw.window.OnWindowEvent(events.Common.WindowFullscreen, func(e *application.WindowEvent) {
        mw.state.Fullscreen = true
        mw.dirty = true
    })

    mw.window.OnWindowEvent(events.Common.WindowUnFullscreen, func(e *application.WindowEvent) {
        mw.state.Fullscreen = false
        mw.dirty = true
    })

    // Position and size events
    mw.window.OnWindowEvent(events.Common.WindowDidMove, func(e *application.WindowEvent) {
        mw.state.X, mw.state.Y = mw.window.Position()
        mw.dirty = true
    })

    mw.window.OnWindowEvent(events.Common.WindowDidResize, func(e *application.WindowEvent) {
        mw.state.Width, mw.state.Height = mw.window.Size()
        mw.dirty = true
    })

    // Cancellable close — use a hook, not a listener.
    mw.window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
        if mw.dirty {
            mw.SaveState()
        }
    })
}

func (mw *ManagedWindow) LoadState() {
    data, err := os.ReadFile("window-state.json")
    if err != nil {
        return
    }

    if err := json.Unmarshal(data, &mw.state); err != nil {
        return
    }

    // Restore window state
    mw.window.SetPosition(mw.state.X, mw.state.Y)
    mw.window.SetSize(mw.state.Width, mw.state.Height)

    if mw.state.Maximised {
        mw.window.Maximise()
    }

    if mw.state.Fullscreen {
        mw.window.Fullscreen()
    }
}

func (mw *ManagedWindow) SaveState() {
    data, err := json.Marshal(mw.state)
    if err != nil {
        return
    }

    os.WriteFile("window-state.json", data, 0644)
    mw.dirty = false

    fmt.Println("Window state saved")
}
```

## 이벤트 조정

### 창 간 이벤트

애플리케이션 이벤트 버스를 사용하여 여러 창을 조정합니다.

```go
// In main window
mainWindow.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
    app.Event.Emit("main-window-focused", nil)
})

// In other windows
app.Event.On("main-window-focused", func(event *application.CustomEvent) {
    updateRelativeToMain()
})
```

### 이벤트 체인

```go
window.OnWindowEvent(events.Common.WindowMaximise, func(e *application.WindowEvent) {
    saveWindowState()
    window.EmitEvent("layout-changed", "maximised")
    app.Event.Emit("window-maximised", window.ID())
})
```

### 디바운스된 이벤트

```go
var resizeTimer *time.Timer

window.OnWindowEvent(events.Common.WindowDidResize, func(e *application.WindowEvent) {
    if resizeTimer != nil {
        resizeTimer.Stop()
    }

    resizeTimer = time.AfterFunc(500*time.Millisecond, func() {
        w, h := window.Size()
        saveWindowSize(w, h)
    })
})
```

## 모범 사례

### ✅ 권장 사항

- **취소에는 훅을 사용하세요**. `RegisterHook` 콜백만 `e.Cancel()`을 호출할 수 있습니다.
- **닫을 때 상태를 저장하세요**. 다음 실행 시 창의 위치와 크기를 복원하세요.
- **자주 발생하는 이벤트를 디바운스하세요**. `WindowDidResize` 및 `WindowDidMove`은 빠르게 연속해서 발생합니다.
- **포커스 변경을 처리하세요**. UI를 적절히 업데이트하세요.
- **이벤트를 통해 조정하세요**. 창 간 메시징에는 `app.Event.Emit`을 사용하세요.
- **핸들러가 더 이상 필요하지 않으면 구독을 해제하세요**. `OnWindowEvent`와 `RegisterHook` 모두 구독 해제용 `func()`을 반환합니다.

### ❌ 금지 사항

- **이벤트 핸들러를 블로킹하지 마세요**. 빠르게 처리되도록 유지하세요.
- <strong>`OnWindowEvent`</strong>에서 취소하려고 하지 마세요. `RegisterHook`을 사용하세요.
- <strong>`window.Destroy()`</strong>은 사용하지 마세요. 존재하지 않으므로 `window.Close()`을 사용하세요.
- **모든 이벤트마다 저장하지 마세요**. 먼저 디바운스하세요.
- **이벤트에 좌표가 있다고 가정하지 마세요**. `window.Position()` / `window.Size()`을 호출하세요.

## 문제 해결

### WindowClosing 훅으로 닫기를 막을 수 없음

**원인:** `RegisterHook` 대신 `OnWindowEvent`을 사용했습니다.

**해결 방법:** `RegisterHook` 콜백만 `e.Cancel()`을 호출할 수 있습니다. `OnWindowEvent` 콜백은 이벤트를 감지할 뿐이며 취소할 수 없습니다.

```go
// ❌ Cannot cancel — this is a listener.
window.OnWindowEvent(events.Common.WindowClosing, func(e *application.WindowEvent) {
    e.Cancel() // no effect
})

// ✅ Can cancel — this is a hook.
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    e.Cancel()
})
```

### 이벤트가 발생하지 않음

**원인:** 이벤트가 발생한 후에 핸들러를 등록했습니다.

**해결 방법:** 창을 생성한 직후(또는 `app.Window.OnCreate` 콜백에서) 핸들러를 등록하세요.

```go
window := app.Window.New()
window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) { /* ... */ })
```

### 메모리 누수

**원인:** 수명이 짧은 상태가 수명이 긴 핸들러를 참조하고 있습니다.

**해결 방법:** `OnWindowEvent`/`RegisterHook`에서 반환된 구독 해제 `func()`을 보관하고 정리할 때 호출하세요.

```go
unsub := window.OnWindowEvent(events.Common.WindowDidResize, handler)
// ...later:
unsub()
```

## 다음 단계

**창 기본 사항** - 창 관리의 기초를 알아보세요 [자세히 알아보기 →](/features/windows/basics/)

**여러 창** - 다중 창 애플리케이션을 위한 패턴을 알아보세요 [자세히 알아보기 →](/features/windows/multiple/)

**이벤트 시스템** - 이벤트 시스템을 자세히 살펴보세요 [자세히 알아보기 →](/features/events/system/)

**애플리케이션 수명 주기** - 애플리케이션 수명 주기를 이해하세요 [자세히 알아보기 →](/concepts/lifecycle/)

---

**궁금한 점이 있나요?** [Discord](https://discord.gg/JDdSxwjhGf)에서 질문하거나 [예제](https://github.com/wailsapp/wails/tree/master/v3/examples)를 확인하세요.
