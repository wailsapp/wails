---
title: "事件 API"
description: "事件 API 完整参考"
slug: "reference/events"
sourcePath: "reference/events.md"
---

## 概述

事件 API 提供触发和监听事件的方法，让应用程序的不同部分能够相互通信。

**事件类型：**

- **应用程序事件** — 应用程序生命周期事件（启动、关闭）
- **窗口事件** — 窗口状态变化（获得焦点、失去焦点、调整大小）
- **自定义事件** — 用户定义的事件，用于应用程序特定的通信

**通信模式：**

- **Go 到前端** — 从 Go 触发事件，在 JavaScript 中监听
- **前端到 Go** — 不支持直接通信（请改用服务绑定）
- **前端到前端** — 通过 Go 或本地运行时事件通信
- **窗口到窗口** — 定向发送到特定窗口，或广播到所有窗口

## 事件方法（Go）

### app.Event.Emit()

向所有窗口触发一个自定义事件。如果钩子取消了事件触发，则返回`true`。

```go
func (em *EventManager) Emit(name string, data ...any) bool
```

**参数：**

- `name` — 事件名称
- `data` — 随事件发送的可选数据

**示例：**

```go
// Emit simple event
app.Event.Emit("user-logged-in")

// Emit with data
app.Event.Emit("data-updated", map[string]interface{}{
    "count": 42,
    "status": "success",
})

// Emit multiple values
app.Event.Emit("progress", 75, "Processing files...")
```

### app.Event.On()

在 Go 中监听自定义事件。

```go
func (em *EventManager) On(name string, callback func(*CustomEvent)) func()
```

**参数：**

- `name` — 要监听的事件名称
- `callback` — 触发事件时调用的函数

<strong>返回值：</strong>用于移除事件监听器的清理函数

**示例：**

```go
// Listen for events
cleanup := app.Event.On("user-action", func(e *application.CustomEvent) {
    data := e.Data.(map[string]interface{})
    action := data["action"].(string)
    app.Logger.Info("User action", "action", action)
})

// Later, remove listener
cleanup()
```

### 特定窗口的事件

向特定窗口触发事件：

```go
// Emit to specific window
window.EmitEvent("notification", "Hello from Go!")

// Emit to all windows
app.Event.Emit("global-update", data)
```

## 事件方法（前端）

### On()

监听来自 Go 的事件。

```javascript
import { Events } from '@wailsio/runtime'

Events.On(eventName, callback)
```

**参数：**

- `eventName` — 要监听的事件名称
- `callback` — 收到事件时调用的函数

<strong>返回值：</strong>清理函数

**示例：**

```javascript
import { Events } from '@wailsio/runtime'

// Listen for events
const cleanup = Events.On('data-updated', (data) => {
    console.log('Count:', data.count)
    console.log('Status:', data.status)
    updateUI(data)
})

// Later, remove listener
cleanup()
```

### Once()

监听一次事件。

```javascript
import { Events } from '@wailsio/runtime'

Events.Once(eventName, callback)
```

**示例：**

```javascript
import { Events } from '@wailsio/runtime'

// Listen for first occurrence only
Events.Once('initialization-complete', (data) => {
    console.log('App initialized!', data)
    // This will only fire once
})
```

### Off()

移除一个或多个事件的所有监听器。`Off`接受数量可变的事件名称字符串，<strong>不</strong>接受回调。要移除单个监听器，请保留`Events.On(...)`返回的取消订阅函数并调用它。

```typescript
import { Events } from '@wailsio/runtime'

Events.Off(...eventNames: string[]): void
```

**示例：**

```javascript
import { Events } from '@wailsio/runtime'

// Preferred: keep the unsubscribe fn from On()
const unsubscribe = Events.On('my-event', (data) => {
    console.log('Event received:', data)
})

// Later — remove just this listener
unsubscribe()

// Or: remove every listener for one or more events
Events.Off('my-event', 'another-event')
```

### OffAll()

移除<strong>所有</strong>事件监听器。不接受任何参数。

```typescript
import { Events } from '@wailsio/runtime'

Events.OffAll(): void
```

**示例：**

```javascript
import { Events } from '@wailsio/runtime'

// Remove all listeners — typically used during teardown.
Events.OffAll()
```

### OnMultiple()

最多监听`max`次事件，随后自动取消订阅。

```typescript
Events.OnMultiple(eventName: string, callback, max: number): () => void
```

**示例：**

```javascript
Events.OnMultiple('progress', (data) => {
    console.log('progress', data)
}, 5)
```

## 应用程序事件

### app.Event.OnApplicationEvent()

监听应用程序生命周期事件。

```go
func (em *EventManager) OnApplicationEvent(
    eventType events.ApplicationEventType,
    callback func(*ApplicationEvent),
) func()
```

事件常量位于`events`包中：跨平台事件使用`events.Common.*`，平台特定事件使用`events.Mac.*`、`events.Windows.*`或`events.Linux.*`。

**常见应用程序事件：**

- `events.Common.ApplicationStarted` — 应用程序已完成启动。
- `events.Common.ThemeChanged` — 系统主题在浅色与深色之间切换。
- `events.Common.ApplicationOpenedWithFile` — 通过文件关联启动。
- `events.Common.ApplicationLaunchedWithUrl` — 通过 URL 方案启动。

并<strong>没有</strong>通用的“应用程序关闭”事件常量；请通过`application.Options.OnShutdown`或`app.OnShutdown(func())`注册关闭清理操作。

**示例：**

```go
import "github.com/wailsapp/wails/v3/pkg/events"

// Handle application startup
app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(e *application.ApplicationEvent) {
    app.Logger.Info("Application started")
})

// React to theme changes
app.Event.OnApplicationEvent(events.Common.ThemeChanged, func(e *application.ApplicationEvent) {
    if e.Context().IsDarkMode() {
        app.Logger.Info("Dark mode enabled")
    }
})

// Shutdown cleanup is configured on the application, not as an event:
app.OnShutdown(func() {
    database.Close()
    saveSettings()
})
```

## 窗口事件

### OnWindowEvent()

监听特定窗口的事件。

```go
func (w *WebviewWindow) OnWindowEvent(
    eventType events.WindowEventType,
    callback func(*WindowEvent),
) func()
```

窗口事件常量位于`events`包中：跨平台事件使用`events.Common.*`（平台特定事件则使用`events.Mac.*`、`events.Windows.*`或`events.Linux.*`）。

**常见窗口事件：**

- `events.Common.WindowFocus` — 窗口获得焦点。
- `events.Common.WindowLostFocus` - 窗口失去焦点。
- `events.Common.WindowClosing` - 窗口即将关闭（可通过`RegisterHook`取消）。
- `events.Common.WindowDidResize` - 窗口大小已调整。
- `events.Common.WindowDidMove` - 窗口已移动。
- `events.Common.WindowMinimise` / `WindowUnMinimise` / `WindowMaximise` / `WindowUnMaximise` / `WindowFullscreen` / `WindowUnFullscreen`。
- `events.Common.WindowRuntimeReady` - 此时可以安全地向窗口内运行时发出事件。

**示例：**

```go
import "github.com/wailsapp/wails/v3/pkg/events"

// Handle window focus
window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
    app.Logger.Info("Window focused")
})

// Handle window resize
window.OnWindowEvent(events.Common.WindowDidResize, func(e *application.WindowEvent) {
    width, height := window.Size()
    app.Logger.Info("Window resized", "width", width, "height", height)
})
```

## 常用模式

以下模式展示了在实际应用中使用事件的成熟方法。每种模式都解决 Go 后端与前端之间特定的通信难题，帮助你构建响应迅速、结构清晰的应用。

### 请求/响应模式

如果要在后端操作完成时通知前端，例如数据获取、文件处理或后台任务完成后，请使用此模式。服务绑定会直接返回数据，而事件可为 UI 更新提供额外通知，例如显示提示消息或刷新列表。

**Go：**

```go
// Service method
type DataService struct {
    app *application.App
}

func (s *DataService) FetchData(query string) ([]Item, error) {
    items := fetchFromDatabase(query)

    // Emit event when done
    s.app.Event.Emit("data-fetched", map[string]interface{}{
        "query": query,
        "count": len(items),
    })

    return items, nil
}
```

**JavaScript：**

```javascript
import { FetchData } from './bindings/DataService'
import { Events } from '@wailsio/runtime'

// Listen for completion event
Events.On('data-fetched', (data) => {
    console.log(`Fetched ${data.count} items for query: ${data.query}`)
    showNotification(`Found ${data.count} results`)
})

// Call service method
const items = await FetchData("search term")
displayItems(items)
```

### 进度更新

此模式非常适合文件上传、批处理、大型数据导入或视频编码等耗时操作。在操作期间发出进度事件，以更新 UI 中的进度条、状态文本或步骤指示器，从而为用户提供实时反馈。

**Go：**

```go
func (s *Service) ProcessFiles(files []string) error {
    total := len(files)

    for i, file := range files {
        // Process file
        processFile(file)

        // Emit progress event
        s.app.Event.Emit("progress", map[string]interface{}{
            "current": i + 1,
            "total":   total,
            "percent": float64(i+1) / float64(total) * 100,
            "file":    file,
        })
    }

    s.app.Event.Emit("processing-complete")
    return nil
}
```

**JavaScript：**

```javascript
import { Events } from '@wailsio/runtime'

// Update progress bar
Events.On('progress', (data) => {
    progressBar.style.width = `${data.percent}%`
    statusText.textContent = `Processing ${data.file}... (${data.current}/${data.total})`
})

// Handle completion
Events.Once('processing-complete', () => {
    progressBar.style.width = '100%'
    statusText.textContent = 'Complete!'
    setTimeout(() => hideProgressBar(), 2000)
})
```

### 多窗口通信

此模式非常适合具有设置面板、仪表板或文档查看器等多个窗口的应用。你可以广播事件，使所有窗口之间的状态（如主题更改和用户偏好设置）保持同步；也可以向特定窗口发送定向事件，以进行窗口专属更新。

**Go：**

```go
// Broadcast to all windows
app.Event.Emit("theme-changed", "dark")

// Send to specific window
preferencesWindow.EmitEvent("settings-updated", settings)

// Per-window listener — receive on the global event bus, but
// gate by the event's source-window name (set automatically when
// a window emits via window.EmitEvent).
app.Event.On("request-data", func(e *application.CustomEvent) {
    if e.Sender != window1.Name() {
        return
    }
    window1.EmitEvent("data-response", data)
})
```

**JavaScript：**

```javascript
import { Events } from '@wailsio/runtime'

// Listen in any window
Events.On('theme-changed', (theme) => {
    document.body.className = theme
})
```

### 状态同步

需要保持前端与后端状态同步时，例如处理用户会话、应用配置或协作功能时，请使用此模式。当后端状态发生变化时，发出事件以更新所有已连接的前端，确保整个应用中的状态一致。

**Go：**

```go
type StateService struct {
    app   *application.App
    state map[string]interface{}
    mu    sync.RWMutex
}

func (s *StateService) UpdateState(key string, value interface{}) {
    s.mu.Lock()
    s.state[key] = value
    s.mu.Unlock()

    // Notify all windows
    s.app.Event.Emit("state-updated", map[string]interface{}{
        "key":   key,
        "value": value,
    })
}

func (s *StateService) GetState(key string) interface{} {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.state[key]
}
```

**JavaScript：**

```javascript
import { Events } from '@wailsio/runtime'
import { GetState } from './bindings/StateService'

// Keep local state in sync
let localState = {}

Events.On('state-updated', async (data) => {
    localState[data.key] = data.value
    updateUI(data.key, data.value)
})

// Initialize state
const initialState = await GetState("all")
localState = initialState
```

### 事件驱动通知

此模式最适合显示成功确认、错误警报或信息消息等用户反馈。不要从服务直接调用 UI 代码，而应发出由前端统一处理的通知事件，以便轻松更改通知样式或添加通知历史记录等功能。

**Go：**

```go
type NotificationService struct {
    app *application.App
}

func (s *NotificationService) Success(message string) {
    s.app.Event.Emit("notification", map[string]interface{}{
        "type":    "success",
        "message": message,
    })
}

func (s *NotificationService) Error(message string) {
    s.app.Event.Emit("notification", map[string]interface{}{
        "type":    "error",
        "message": message,
    })
}

func (s *NotificationService) Info(message string) {
    s.app.Event.Emit("notification", map[string]interface{}{
        "type":    "info",
        "message": message,
    })
}
```

**JavaScript：**

```javascript
import { Events } from '@wailsio/runtime'

// Unified notification handler
Events.On('notification', (data) => {
    const toast = document.createElement('div')
    toast.className = `toast toast-${data.type}`
    toast.textContent = data.message

    document.body.appendChild(toast)

    setTimeout(() => {
        toast.classList.add('fade-out')
        setTimeout(() => toast.remove(), 300)
    }, 3000)
})
```

## 完整示例

**Go：**

```go
package main

import (
    "sync"
    "time"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)

type EventDemoService struct {
    app *application.App
    mu  sync.Mutex
}

func NewEventDemoService(app *application.App) *EventDemoService {
    service := &EventDemoService{app: app}

    // Listen for custom events
    app.Event.On("user-action", func(e *application.CustomEvent) {
        data := e.Data.(map[string]interface{})
        app.Logger.Info("User action received", "data", data)
    })

    return service
}

func (s *EventDemoService) StartLongTask() {
    go func() {
        s.app.Event.Emit("task-started")

        for i := 1; i <= 10; i++ {
            time.Sleep(500 * time.Millisecond)

            s.app.Event.Emit("task-progress", map[string]interface{}{
                "step":    i,
                "total":   10,
                "percent": i * 10,
            })
        }

        s.app.Event.Emit("task-completed", map[string]interface{}{
            "message": "Task finished successfully!",
        })
    }()
}

func (s *EventDemoService) BroadcastMessage(message string) {
    s.app.Event.Emit("broadcast", message)
}

func main() {
    app := application.New(application.Options{
        Name: "Event Demo",
    })

    // Handle application lifecycle
    app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(e *application.ApplicationEvent) {
        app.Logger.Info("Application started!")
    })

    app.OnShutdown(func() {
        app.Logger.Info("Application shutting down...")
    })

    // Register service (RegisterService returns nothing).
    service := NewEventDemoService(app)
    app.RegisterService(application.NewService(service))

    // Create window
    window := app.Window.New()

    // Handle window events
    window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
        window.EmitEvent("window-state", "focused")
    })

    window.OnWindowEvent(events.Common.WindowLostFocus, func(e *application.WindowEvent) {
        window.EmitEvent("window-state", "blurred")
    })

    window.Show()
    app.Run()
}
```

**JavaScript：**

```javascript
import { Events } from '@wailsio/runtime'
import { StartLongTask, BroadcastMessage } from './bindings/EventDemoService'

// Task events
Events.On('task-started', () => {
    console.log('Task started...')
    document.getElementById('status').textContent = 'Running...'
})

Events.On('task-progress', (data) => {
    const progressBar = document.getElementById('progress')
    progressBar.style.width = `${data.percent}%`
    console.log(`Step ${data.step} of ${data.total}`)
})

Events.Once('task-completed', (data) => {
    console.log('Task completed!', data.message)
    document.getElementById('status').textContent = data.message
})

// Broadcast events
Events.On('broadcast', (message) => {
    console.log('Broadcast:', message)
    alert(message)
})

// Window state events
Events.On('window-state', (state) => {
    console.log('Window is now:', state)
    document.body.dataset.windowState = state
})

// Trigger long task
document.getElementById('startTask').addEventListener('click', async () => {
    await StartLongTask()
})

// Send broadcast
document.getElementById('broadcast').addEventListener('click', async () => {
    const message = document.getElementById('message').value
    await BroadcastMessage(message)
})
```

## 内置事件

Wails 为应用和窗口生命周期提供内置系统事件。这些事件由框架自动发出。

### 通用事件与平台原生事件

Wails 提供两类系统事件：

**通用事件**（`events.Common.*`）是跨平台抽象，在 macOS、Windows 和 Linux 上具有一致的行为。为了最大限度地提高可移植性，你应在应用中使用这些事件。

**平台原生事件**（`events.Mac.*`、`events.Windows.*`、`events.Linux.*`）是通用事件所映射自的底层操作系统特有事件。通过这些事件可以处理平台特有的行为和边界情况。

**工作原理：**

```go
import "github.com/wailsapp/wails/v3/pkg/events"

// ✅ RECOMMENDED: Use Common Events for cross-platform code
window.OnWindowEvent(events.Common.WindowClosing, func(e *application.WindowEvent) {
    // This works on all platforms
})

// Platform-specific events for advanced use cases
window.OnWindowEvent(events.Mac.WindowWillClose, func(e *application.WindowEvent) {
    // macOS-specific "will close" event (before WindowClosing)
})

window.OnWindowEvent(events.Windows.WindowClosing, func(e *application.WindowEvent) {
    // Windows-specific close event
})
```

**事件映射：**

平台原生事件会自动映射到通用事件：

- macOS：`events.Mac.WindowShouldClose` → `events.Common.WindowClosing`
- Windows：`events.Windows.WindowClosing` → `events.Common.WindowClosing`
- Linux：`events.Linux.WindowDeleteEvent` → `events.Common.WindowClosing`

此映射会在后台自动完成，因此无论在哪个平台上，监听`events.Common.WindowClosing`时都能收到该事件。

**各自的适用场景：**

- <strong>使用通用事件</strong>处理应用代码中99% 的场景——它们可在不同平台上提供一致的行为
- 仅当需要通用事件未提供的平台特有功能时，才<strong>使用平台原生事件</strong>（例如 macOS 特有的窗口生命周期事件或 Windows 电源管理事件）

### 应用事件

| 事件 | 说明 | 发出时机 | 可取消 |
| --- | --- | --- | --- |
| `ApplicationOpenedWithFile` | 应用通过文件打开 | 应用通过文件启动时（例如通过文件关联） | 否 |
| `ApplicationStarted` | 应用已完成启动 | 应用初始化完成并准备就绪后 | 否 |
| `ApplicationLaunchedWithUrl` | 应用通过 URL 启动 | 应用通过 URL 方案启动时 | 否 |
| `ThemeChanged` | 系统主题已更改 | 操作系统主题在浅色与深色模式之间切换时 | 否 |
| `SystemWillSleep` | 系统即将挂起 | 操作系统挂起前一刻（macOS / Windows / 使用 logind 的 Linux） | 否 |
| `SystemDidWake` | 系统已从挂起状态恢复 | 从睡眠状态恢复后立即触发（macOS / Windows / 使用 logind 的 Linux） | 否 |

**用法：**

```go
import "github.com/wailsapp/wails/v3/pkg/events"

app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(e *application.ApplicationEvent) {
    app.Logger.Info("Application ready!")
})

app.Event.OnApplicationEvent(events.Common.ThemeChanged, func(e *application.ApplicationEvent) {
    // Update app theme
})
```

### 窗口事件

| 事件 | 说明 | 触发时机 | 可取消 |
| --- | --- | --- | --- |
| `WindowClosing` | 窗口即将关闭 | 窗口关闭前（用户单击 X 或调用 Close()） | 是 |
| `WindowDidMove` | 窗口已移动到新位置 | 窗口位置更改后（已防抖） | 否 |
| `WindowDidResize` | 窗口大小已调整 | 窗口大小更改后 | 否 |
| `WindowDPIChanged` | 窗口 DPI 缩放比例已更改 | 在 DPI 不同的显示器之间移动窗口时（Windows） | 否 |
| `WindowFilesDropped` | 通过操作系统原生拖放功能拖入文件 | 从操作系统将文件拖放到窗口后 | 否 |
| `WindowFocus` | 窗口已获得焦点 | 窗口变为活动状态时 | 否 |
| `WindowFullscreen` | 窗口已进入全屏模式 | 调用 Fullscreen() 或用户进入全屏模式后 | 否 |
| `WindowHide` | 窗口已隐藏 | 调用 Hide() 或窗口被遮挡后 | 否 |
| `WindowLostFocus` | 窗口已失去焦点 | 窗口变为非活动状态时 | 否 |
| `WindowMaximise` | 窗口已最大化 | 调用 Maximise() 或用户最大化窗口后 | 是（macOS） |
| `WindowMinimise` | 窗口已最小化 | 调用 Minimise() 或用户最小化窗口后 | 是（macOS） |
| `WindowRestore` | 窗口已从最小化/最大化状态还原 | 调用 Restore() 后（主要适用于 Windows） | 否 |
| `WindowRuntimeReady` | Wails 运行时已加载并准备就绪 | JavaScript 运行时初始化完成时 | 否 |
| `WindowShow` | 窗口已变为可见 | 调用 Show() 后或窗口变为可见时 | 否 |
| `WindowUnFullscreen` | 窗口退出全屏模式 | 调用 UnFullscreen() 后或用户退出全屏模式时 | 否 |
| `WindowUnMaximise` | 窗口退出最大化状态 | 调用 UnMaximise() 后或用户取消最大化时 | 是（macOS） |
| `WindowUnMinimise` | 窗口退出最小化状态 | 调用 UnMinimise()/Restore() 后或用户还原窗口时 | 是（macOS） |
| `WindowZoomIn` | 窗口内容缩放比例增大 | 调用 ZoomIn() 后（主要适用于 macOS） | 是（macOS） |
| `WindowZoomOut` | 窗口内容缩放比例减小 | 调用 ZoomOut() 后（主要适用于 macOS） | 是（macOS） |
| `WindowZoomReset` | 窗口内容缩放比例重置为 100% | 调用 ZoomReset() 后（主要适用于 macOS） | 是（macOS） |
| `WindowDropZoneFilesDropped` | 文件被拖放到由 JS 定义的放置区 | 文件被拖放到设有放置区的元素上时 | 否 |

**用法：**

```go
import "github.com/wailsapp/wails/v3/pkg/events"

// Listen for window events
window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
    app.Logger.Info("Window focused")
})

// Cancel window close
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    dlg := app.Dialog.Question().SetMessage("Close window?")
    yes := dlg.AddButton("Yes")
    no := dlg.AddButton("No")
    dlg.SetDefaultButton(yes)
    dlg.SetCancelButton(no)
    no.OnClick(func() { e.Cancel() })

    dlg.Show()
})

// Wait for runtime ready
window.OnWindowEvent(events.Common.WindowRuntimeReady, func(e *application.WindowEvent) {
    app.Logger.Info("Runtime ready, safe to emit events to frontend")
    window.EmitEvent("app-initialized", data)
})
```

**重要说明：**

- **WindowRuntimeReady** 至关重要——向前端发出事件之前，请等待此事件
- **WindowDidMove** 和 **WindowDidResize** 经过防抖处理（默认为 50ms），以防止事件泛滥
- 可以在 `RegisterHook()` 处理程序中调用 `event.Cancel()`，阻止<strong>可取消事件</strong>
- **WindowFilesDropped** 用于操作系统原生文件拖放；**WindowDropZoneFilesDropped** 用于基于 Web 的放置区
- 部分事件仅适用于特定平台（例如，WindowDPIChanged 适用于 Windows，缩放事件主要适用于 macOS）

## 事件命名约定

```go
// Good - descriptive and specific
app.Event.Emit("user:logged-in", user)
app.Event.Emit("data:fetch:complete", results)
app.Event.Emit("ui:theme:changed", theme)

// Bad - vague and unclear
app.Event.Emit("event1", data)
app.Event.Emit("update", stuff)
app.Event.Emit("e", value)
```

## 性能注意事项

### 高频事件防抖

```go
type Service struct {
    app            *application.App
    lastEmit       time.Time
    debounceWindow time.Duration
}

func (s *Service) EmitWithDebounce(event string, data interface{}) {
    now := time.Now()
    if now.Sub(s.lastEmit) < s.debounceWindow {
        return // Skip this emission
    }

    s.app.Event.Emit(event, data)
    s.lastEmit = now
}
```

### 事件节流

```javascript
import { Events } from '@wailsio/runtime'

let lastUpdate = 0
const throttleMs = 100

Events.On('high-frequency-event', (data) => {
    const now = Date.now()
    if (now - lastUpdate < throttleMs) {
        return // Skip this update
    }

    processUpdate(data)
    lastUpdate = now
})
```
