---
title: "視窗事件"
description: "處理視窗生命週期與狀態變更事件"
slug: "features/windows/events"
sourcePath: "features/windows/events.md"
---

## 視窗事件

Wails 透過單一 API 分派視窗生命週期與狀態變更事件：使用`OnWindowEvent`註冊被動監聽器，使用`RegisterHook`註冊可取消的勾點，以阻止預設動作。

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

這兩種呼叫都會傳回一個`unsubscribe func()`，您可以呼叫它來移除處理常式。

跨平台事件位於`events.Common.*`。平台特定事件位於`events.Mac.*`、`events.Windows.*`和`events.Linux.*`。完整清單會產生於`v3/pkg/events/events.go`。

## 生命週期事件

### 建立視窗

使用`app.Window.OnCreate`，在每次建立視窗時執行回呼：

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

回呼會接收`application.Window`（介面）。每個視窗的`OnCreate`回呼都會在其執行階段初始化後呼叫一次。

### WindowClosing

當使用者嘗試關閉視窗時觸發（例如按一下 X、⌘W 或 Alt+F4）。

使用<strong>勾點</strong>（`RegisterHook`）——只有勾點能取消關閉動作。監聽器可以觀察事件，但無法阻止關閉。

```go
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    if hasUnsavedChanges() {
        // Cancel the close — the window stays open.
        e.Cancel()
    }
})
```

**重要：**

- `WindowClosing`會在使用者發起關閉嘗試時分派。
- `RegisterHook`回呼可以呼叫`e.Cancel()`，讓視窗保持開啟。
- 針對`WindowClosing`的`OnWindowEvent`回呼是觀察者——它們會觸發，但無法取消。

<strong>以程式方式關閉：</strong>呼叫`window.Close()`。不存在`window.Destroy()`。

### WindowRuntimeReady

視窗內執行階段完成初始化後觸發——此時可安全地呼叫視窗 JS 執行環境中的程式碼：

```go
window.OnWindowEvent(events.Common.WindowRuntimeReady, func(e *application.WindowEvent) {
    window.EmitEvent("app-ready", nil)
})
```

## 焦點事件

### WindowFocus

當視窗取得焦點時呼叫：

```go
window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
    fmt.Println("Window gained focus")
    updateTitleBar(true)
    app.Event.Emit("window-focused", window.ID())
})
```

### WindowLostFocus

當視窗失去焦點時呼叫：

```go
window.OnWindowEvent(events.Common.WindowLostFocus, func(e *application.WindowEvent) {
    fmt.Println("Window lost focus")
    updateTitleBar(false)
    saveCurrentState()
})
```

**範例：感知焦點的 UI：**

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

## 狀態變更事件

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

使用`window.Fullscreen()` / `window.UnFullscreen()` / `window.ToggleFullscreen()`以程式方式進入或退出全螢幕——不存在`SetFullscreen(bool)`。使用`window.IsFullscreen() bool`查詢狀態。

## 位置與大小事件

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

`WindowDidResize`和`WindowDidMove`的事件本身不包含座標——請在回呼內透過`window.Size()` / `window.Position()`查詢視窗。

**範例：響應式版面配置：**

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

## 完整範例

具備完整事件處理、可用於生產環境的視窗：

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

## 事件協調

### 跨視窗事件

使用應用程式事件匯流排協調多個視窗：

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

### 事件鏈

```go
window.OnWindowEvent(events.Common.WindowMaximise, func(e *application.WindowEvent) {
    saveWindowState()
    window.EmitEvent("layout-changed", "maximised")
    app.Event.Emit("window-maximised", window.ID())
})
```

### 防抖事件

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

## 最佳做法

### ✅ 建議做法

- **使用勾點取消動作**——只有`RegisterHook`回呼能呼叫`e.Cancel()`。
- **關閉時儲存狀態**——下次啟動時還原視窗的位置和大小。
- **對頻繁事件進行防抖處理**——`WindowDidResize`和`WindowDidMove`會快速連續觸發。
- **處理焦點變更**——適當更新 UI。
- **透過事件進行協調**——使用`app.Event.Emit`進行跨視窗通訊。
- <strong>取消訂閱</strong>不再需要的處理常式——`OnWindowEvent`和`RegisterHook`都會傳回用於取消訂閱的`func()`。

### ❌ 請勿這樣做

- **不要阻塞事件處理常式**——應讓它們快速完成。
- <strong>不要嘗試從`OnWindowEvent`</strong>取消動作——請使用`RegisterHook`。
- **不要使用`window.Destroy()`**——它並不存在；請使用`window.Close()`。
- **不要在每次事件觸發時都儲存**——請先進行防抖處理。
- **不要假設事件包含座標**——請呼叫`window.Position()` / `window.Size()`。

## 疑難排解

### WindowClosing 勾點無法阻止關閉

<strong>原因：</strong>您使用了`OnWindowEvent`，而非`RegisterHook`。

<strong>解決方案：</strong>只有`RegisterHook`回呼能呼叫`e.Cancel()`。`OnWindowEvent`回呼只能觀察，無法取消。

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

### 事件未觸發

<strong>原因：</strong>事件發生後才註冊處理常式。

<strong>解決方法：</strong>請在建立視窗後立即註冊處理常式（或在`app.Window.OnCreate`回呼中註冊）。

```go
window := app.Window.New()
window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) { /* ... */ })
```

### 記憶體洩漏

<strong>原因：</strong>短生命週期的狀態持有長生命週期的處理常式。

<strong>解決方法：</strong>保留`OnWindowEvent`/`RegisterHook`傳回的取消訂閱`func()`，並在清理時呼叫它。

```go
unsub := window.OnWindowEvent(events.Common.WindowDidResize, handler)
// ...later:
unsub()
```

## 後續步驟

**視窗基礎** - 瞭解視窗管理的基礎知識 [深入瞭解 →](/features/windows/basics/)

**多視窗** - 多視窗應用程式的設計模式 [深入瞭解 →](/features/windows/multiple/)

**事件系統** - 深入探索事件系統 [深入瞭解 →](/features/events/system/)

**應用程式生命週期** - 瞭解應用程式生命週期 [深入瞭解 →](/concepts/lifecycle/)

---

<strong>有問題嗎？</strong>請在[Discord](https://discord.gg/JDdSxwjhGf)中提問，或查看[範例](https://github.com/wailsapp/wails/tree/master/v3/examples)。
