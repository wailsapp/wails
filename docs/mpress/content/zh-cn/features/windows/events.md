---
title: "窗口事件"
description: "处理窗口生命周期和状态变更事件"
slug: "features/windows/events"
sourcePath: "features/windows/events.md"
---

## 窗口事件

Wails 通过统一的 API 分发窗口生命周期和状态变更事件：`OnWindowEvent`用于被动监听器，`RegisterHook`用于可取消默认操作的钩子。

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

这两种调用都会返回一个`unsubscribe func()`，你可以调用它来移除处理程序。

跨平台事件位于`events.Common.*`中。平台特定事件位于`events.Mac.*`、`events.Windows.*`和`events.Linux.*`中。完整列表生成在`v3/pkg/events/events.go`中。

## 生命周期事件

### 窗口创建

使用`app.Window.OnCreate`在每次创建窗口时运行回调：

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

该回调接收`application.Window`（一个接口）。每个窗口的`OnCreate`回调都会在窗口运行时初始化完成后调用一次。

### WindowClosing

当用户尝试关闭窗口时触发（单击 X、按 ⌘W、Alt+F4 等）。

使用<strong>钩子</strong>（`RegisterHook`）——只有钩子才能取消关闭操作。监听器只能观察事件，无法阻止关闭。

```go
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    if hasUnsavedChanges() {
        // Cancel the close — the window stays open.
        e.Cancel()
    }
})
```

**重要：**

- `WindowClosing`在用户发起关闭尝试时分发。
- `RegisterHook`回调可以调用`e.Cancel()`以保持窗口打开。
- `WindowClosing`的`OnWindowEvent`回调是观察者——它们会触发，但无法取消操作。

<strong>以编程方式关闭：</strong>调用`window.Close()`。不存在`window.Destroy()`。

### WindowRuntimeReady

窗口内运行时完成初始化后触发——此时可以安全地调用窗口的 JS 上下文：

```go
window.OnWindowEvent(events.Common.WindowRuntimeReady, func(e *application.WindowEvent) {
    window.EmitEvent("app-ready", nil)
})
```

## 焦点事件

### WindowFocus

窗口获得焦点时调用：

```go
window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
    fmt.Println("Window gained focus")
    updateTitleBar(true)
    app.Event.Emit("window-focused", window.ID())
})
```

### WindowLostFocus

窗口失去焦点时调用：

```go
window.OnWindowEvent(events.Common.WindowLostFocus, func(e *application.WindowEvent) {
    fmt.Println("Window lost focus")
    updateTitleBar(false)
    saveCurrentState()
})
```

**示例：感知焦点的 UI：**

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

## 状态变更事件

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

使用`window.Fullscreen()` / `window.UnFullscreen()` / `window.ToggleFullscreen()`以编程方式进入或退出全屏——不存在`SetFullscreen(bool)`。使用`window.IsFullscreen() bool`查询状态。

## 位置和大小事件

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

`WindowDidResize`和`WindowDidMove`事件本身不携带坐标——请在回调中通过`window.Size()` / `window.Position()`查询窗口。

**示例：响应式布局：**

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

## 完整示例

一个具有完整事件处理功能、可用于生产环境的窗口：

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

## 事件协调

### 跨窗口事件

使用应用程序事件总线在多个窗口之间进行协调：

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

### 事件链

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

## 最佳实践

### ✅ 应该做

- **使用钩子取消操作**——只有`RegisterHook`回调可以调用`e.Cancel()`。
- **关闭时保存状态**——下次启动时恢复窗口的位置和大小。
- **对高频事件进行防抖**——`WindowDidResize`和`WindowDidMove`会快速连续触发。
- **处理焦点变化**——相应地更新 UI。
- **通过事件进行协调**——使用`app.Event.Emit`进行跨窗口消息传递。
- <strong>取消订阅</strong>不再需要的处理程序——`OnWindowEvent`和`RegisterHook`都会返回一个用于取消订阅的`func()`。

### ❌ 不应该做

- **不要阻塞事件处理程序**——应让它们快速完成。
- <strong>不要尝试从`OnWindowEvent`</strong>中取消操作——请使用`RegisterHook`。
- **不要使用`window.Destroy()`**——它不存在；请使用`window.Close()`。
- **不要在每次事件触发时都保存**——应先进行防抖。
- **不要假定事件中包含坐标**——请调用`window.Position()` / `window.Size()`。

## 故障排除

### WindowClosing 钩子无法阻止关闭

<strong>原因：</strong>你使用了`OnWindowEvent`，而不是`RegisterHook`。

<strong>解决方案：</strong>只有`RegisterHook`回调可以调用`e.Cancel()`。`OnWindowEvent`回调只能观察，无法取消操作。

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

### 事件未触发

<strong>原因：</strong>事件发生后才注册处理程序。

<strong>解决方案：</strong>创建窗口后立即注册处理程序（或在`app.Window.OnCreate`回调中注册）。

```go
window := app.Window.New()
window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) { /* ... */ })
```

### 内存泄漏

<strong>原因：</strong>短生命周期状态持有了长生命周期的处理程序。

<strong>解决方案：</strong>保存`OnWindowEvent`/`RegisterHook`返回的取消订阅`func()`，并在清理时调用它。

```go
unsub := window.OnWindowEvent(events.Common.WindowDidResize, handler)
// ...later:
unsub()
```

## 后续步骤

**窗口基础** - 了解窗口管理的基础知识 [了解更多 →](/features/windows/basics/)

**多窗口** - 了解多窗口应用的设计模式 [了解更多 →](/features/windows/multiple/)

**事件系统** - 深入了解事件系统 [了解更多 →](/features/events/system/)

**应用生命周期** - 了解应用生命周期 [了解更多 →](/concepts/lifecycle/)

---

<strong>有疑问？</strong>请在[Discord](https://discord.gg/JDdSxwjhGf)中提问，或查看[示例](https://github.com/wailsapp/wails/tree/master/v3/examples)。
