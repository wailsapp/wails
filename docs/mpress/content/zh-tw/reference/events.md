---
title: "事件 API"
description: "事件 API 完整參考資料"
slug: "reference/events"
sourcePath: "reference/events.md"
---

## 概觀

事件 API 提供發出及監聽事件的方法，讓應用程式的不同部分能夠互相通訊。

**事件類型：**

- **應用程式事件**－應用程式生命週期事件（啟動、關閉）
- **視窗事件**－視窗狀態變更（取得焦點、失去焦點、調整大小）
- **自訂事件**－使用者為應用程式特定通訊定義的事件

**通訊模式：**

- **Go 到前端**－從 Go 發出事件，並在 JavaScript 中監聽
- **前端到 Go**－不支援直接通訊（請改用服務繫結）
- **前端到前端**－透過 Go 或本機執行階段事件
- **視窗到視窗**－指定特定視窗，或廣播至所有視窗

## 事件方法（Go）

### app.Event.Emit()

向所有視窗發出自訂事件。如果掛鉤取消發出事件，則傳回`true`。

```go
func (em *EventManager) Emit(name string, data ...any) bool
```

**參數：**

- `name`－事件名稱
- `data`－要隨事件傳送的選用資料

**範例：**

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

在 Go 中監聽自訂事件。

```go
func (em *EventManager) On(name string, callback func(*CustomEvent)) func()
```

**參數：**

- `name`－要監聽的事件名稱
- `callback`－發出事件時呼叫的函式

<strong>傳回值：</strong>用於移除事件監聽器的清理函式

**範例：**

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

### 視窗特定事件

向特定視窗發出事件：

```go
// Emit to specific window
window.EmitEvent("notification", "Hello from Go!")

// Emit to all windows
app.Event.Emit("global-update", data)
```

## 事件方法（前端）

### On()

監聽來自 Go 的事件。

```javascript
import { Events } from '@wailsio/runtime'

Events.On(eventName, callback)
```

**參數：**

- `eventName`－要監聽的事件名稱
- `callback`－收到事件時呼叫的函式

<strong>傳回值：</strong>清理函式

**範例：**

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

監聽事件的一次發生。

```javascript
import { Events } from '@wailsio/runtime'

Events.Once(eventName, callback)
```

**範例：**

```javascript
import { Events } from '@wailsio/runtime'

// Listen for first occurrence only
Events.Once('initialization-complete', (data) => {
    console.log('App initialized!', data)
    // This will only fire once
})
```

### Off()

移除一或多個事件的所有監聽器。`Off`接受可變數量的事件名稱字串，<strong>不</strong>接受回呼。若要移除單一監聽器，請保留`Events.On(...)`傳回的取消訂閱函式，並呼叫該函式。

```typescript
import { Events } from '@wailsio/runtime'

Events.Off(...eventNames: string[]): void
```

**範例：**

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

移除<strong>所有</strong>事件監聽器。不接受任何引數。

```typescript
import { Events } from '@wailsio/runtime'

Events.OffAll(): void
```

**範例：**

```javascript
import { Events } from '@wailsio/runtime'

// Remove all listeners — typically used during teardown.
Events.OffAll()
```

### OnMultiple()

監聽事件最多`max`次，之後自動取消訂閱。

```typescript
Events.OnMultiple(eventName: string, callback, max: number): () => void
```

**範例：**

```javascript
Events.OnMultiple('progress', (data) => {
    console.log('progress', data)
}, 5)
```

## 應用程式事件

### app.Event.OnApplicationEvent()

監聽應用程式生命週期事件。

```go
func (em *EventManager) OnApplicationEvent(
    eventType events.ApplicationEventType,
    callback func(*ApplicationEvent),
) func()
```

事件常數位於`events`套件中：跨平台事件使用`events.Common.*`，平台特定事件則使用`events.Mac.*`／`events.Windows.*`／`events.Linux.*`。

**常見應用程式事件：**

- `events.Common.ApplicationStarted`－應用程式已完成啟動。
- `events.Common.ThemeChanged`－系統佈景主題已在淺色與深色之間切換。
- `events.Common.ApplicationOpenedWithFile`－透過檔案關聯啟動。
- `events.Common.ApplicationLaunchedWithUrl`－透過 URL 通訊協定啟動。

並<strong>沒有</strong>通用的「應用程式關閉」事件常數；請透過`application.Options.OnShutdown`或`app.OnShutdown(func())`註冊關閉時的清理作業。

**範例：**

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

## 視窗事件

### OnWindowEvent()

監聽視窗特定事件。

```go
func (w *WebviewWindow) OnWindowEvent(
    eventType events.WindowEventType,
    callback func(*WindowEvent),
) func()
```

視窗事件常數位於`events`套件中：跨平台事件使用`events.Common.*`（平台特定事件則使用`events.Mac.*`／`events.Windows.*`／`events.Linux.*`）。

**常見視窗事件：**

- `events.Common.WindowFocus`－視窗已取得焦點。
- `events.Common.WindowLostFocus`－視窗失去焦點。
- `events.Common.WindowClosing`－視窗即將關閉（可透過`RegisterHook`取消）。
- `events.Common.WindowDidResize`－視窗大小已變更。
- `events.Common.WindowDidMove`－視窗已移動。
- `events.Common.WindowMinimise`／`WindowUnMinimise`／`WindowMaximise`／`WindowUnMaximise`／`WindowFullscreen`／`WindowUnFullscreen`。
- `events.Common.WindowRuntimeReady`－此時可安全地向視窗內的執行階段發出事件。

**範例：**

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

## 常見模式

這些模式展示了在實際應用程式中使用事件的成熟方法。每種模式都能解決 Go 後端與前端之間特定的通訊難題，協助您建置反應靈敏且結構良好的應用程式。

### 要求／回應模式

當您想在後端作業完成時通知前端（例如擷取資料、處理檔案或執行背景工作之後），請使用此模式。服務繫結會直接傳回資料，而事件則提供額外通知，以便更新 UI，例如顯示快顯通知或重新整理清單。

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

### 進度更新

適合檔案上傳、批次處理、大量資料匯入或影片編碼等長時間執行的作業。在作業期間發出進度事件，以更新 UI 中的進度列、狀態文字或步驟指示器，為使用者提供即時回饋。

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

### 多視窗通訊

非常適合具有設定面板、儀表板或文件檢視器等多個視窗的應用程式。您可以廣播事件，在所有視窗間同步狀態（例如佈景主題變更和使用者偏好設定）；也可以將事件傳送至特定視窗，進行該視窗專屬的更新。

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

### 狀態同步

當您需要讓前端與後端狀態保持同步時使用此模式，例如同步使用者工作階段、應用程式設定或協作功能。後端狀態變更時，發出事件以更新所有已連線的前端，確保整個應用程式保持一致。

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

### 事件驅動通知

最適合顯示成功確認、錯誤警示或資訊訊息等使用者回饋。不要從服務直接呼叫 UI 程式碼，而應發出通知事件，由前端以一致的方式處理。這樣即可輕鬆變更通知樣式，或加入通知記錄等功能。

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

## 完整範例

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

## 內建事件

Wails 為應用程式和視窗生命週期提供內建系統事件。這些事件會由框架自動發出。

### 通用事件與平台原生事件

Wails 提供兩種類型的系統事件：

**通用事件**（`events.Common.*`）是跨平台抽象，在 macOS、Windows 和 Linux 上都能以一致的方式運作。為獲得最高的可攜性，您的應用程式應使用這些事件。

**平台原生事件**（`events.Mac.*`、`events.Windows.*`、`events.Linux.*`）是通用事件所對應的底層作業系統特定事件。這些事件可供您存取平台特定行為與邊界情況。

**運作方式：**

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

**事件對應：**

平台原生事件會自動對應至通用事件：

- macOS：`events.Mac.WindowShouldClose` → `events.Common.WindowClosing`
- Windows：`events.Windows.WindowClosing` → `events.Common.WindowClosing`
- Linux：`events.Linux.WindowDeleteEvent` → `events.Common.WindowClosing`

此對應會在背景自動進行，因此無論使用哪個平台，只要監聽`events.Common.WindowClosing`，就會收到該事件。

**各自的使用時機：**

- <strong>使用通用事件</strong>處理應用程式程式碼中99% 的情況－這些事件可在不同平台上提供一致的行為
- 只有在需要通用事件未提供的平台特定功能時，才<strong>使用平台原生事件</strong>（例如 macOS 特有的視窗生命週期事件或 Windows 電源管理事件）

### 應用程式事件

| 事件 | 說明 | 發出時機 | 可取消 |
| --- | --- | --- | --- |
| `ApplicationOpenedWithFile` | 透過檔案開啟應用程式 | 應用程式隨檔案啟動時（例如透過檔案關聯） | 否 |
| `ApplicationStarted` | 應用程式已完成啟動 | 應用程式初始化完成並就緒後 | 否 |
| `ApplicationLaunchedWithUrl` | 應用程式透過 URL 啟動 | 應用程式透過 URL 通訊協定啟動時 | 否 |
| `ThemeChanged` | 系統佈景主題已變更 | 作業系統的佈景主題在淺色與深色模式之間切換時 | 否 |
| `SystemWillSleep` | 系統即將進入暫停狀態 | 作業系統進入暫停狀態前一刻（macOS／Windows／使用 logind 的 Linux） | 否 |
| `SystemDidWake` | 系統已從暫停狀態恢復 | 從睡眠狀態恢復時立即觸發（macOS／Windows／使用 logind 的 Linux） | 否 |

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

### 視窗事件

| 事件 | 說明 | 觸發時機 | 可取消 |
| --- | --- | --- | --- |
| `WindowClosing` | 視窗即將關閉 | 視窗關閉前（使用者按一下 X 或呼叫 Close()） | 是 |
| `WindowDidMove` | 視窗已移至新位置 | 視窗位置變更後（已套用防彈跳） | 否 |
| `WindowDidResize` | 視窗大小已調整 | 視窗大小變更後 | 否 |
| `WindowDPIChanged` | 視窗的 DPI 縮放比例已變更 | 在具有不同 DPI 的顯示器之間移動時（Windows） | 否 |
| `WindowFilesDropped` | 透過作業系統原生拖放功能置入檔案 | 從作業系統將檔案拖放至視窗後 | 否 |
| `WindowFocus` | 視窗已取得焦點 | 視窗成為作用中視窗時 | 否 |
| `WindowFullscreen` | 視窗已進入全螢幕模式 | 呼叫 Fullscreen() 或使用者進入全螢幕模式後 | 否 |
| `WindowHide` | 視窗已隱藏 | 呼叫 Hide() 或視窗被遮蔽後 | 否 |
| `WindowLostFocus` | 視窗已失去焦點 | 視窗變為非作用中時 | 否 |
| `WindowMaximise` | 視窗已最大化 | 呼叫 Maximise() 或使用者將視窗最大化後 | 是（macOS） |
| `WindowMinimise` | 視窗已最小化 | 呼叫 Minimise() 或使用者將視窗最小化後 | 是（macOS） |
| `WindowRestore` | 視窗已從最小化／最大化狀態還原 | 呼叫 Restore() 後（主要適用於 Windows） | 否 |
| `WindowRuntimeReady` | Wails 執行階段已載入並就緒 | JavaScript 執行階段完成初始化時 | 否 |
| `WindowShow` | 視窗已變為可見 | 呼叫 Show() 後或視窗變為可見時 | 否 |
| `WindowUnFullscreen` | 視窗已退出全螢幕模式 | 呼叫 UnFullscreen() 後或使用者退出全螢幕模式時 | 否 |
| `WindowUnMaximise` | 視窗已退出最大化狀態 | 呼叫 UnMaximise() 後或使用者取消最大化時 | 是（macOS） |
| `WindowUnMinimise` | 視窗已退出最小化狀態 | 呼叫 UnMinimise()/Restore() 後或使用者還原視窗時 | 是（macOS） |
| `WindowZoomIn` | 視窗內容的縮放比例已提高 | 呼叫 ZoomIn() 後（主要適用於 macOS） | 是（macOS） |
| `WindowZoomOut` | 視窗內容的縮放比例已降低 | 呼叫 ZoomOut() 後（主要適用於 macOS） | 是（macOS） |
| `WindowZoomReset` | 視窗內容的縮放比例已重設為100% | 呼叫 ZoomReset() 後（主要適用於 macOS） | 是（macOS） |
| `WindowDropZoneFilesDropped` | 檔案已拖放至 JS 定義的拖放區域 | 將檔案拖放至設有拖放區域的元素時 | 否 |

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

**重要注意事項：**

- **WindowRuntimeReady** 至關重要——請先等待此事件，再向前端發出事件
- **WindowDidMove** 和 **WindowDidResize** 會進行防彈跳處理（預設為50毫秒），以避免事件大量湧入
- 在`RegisterHook()`處理常式中呼叫`event.Cancel()`，即可阻止<strong>可取消事件</strong>
- **WindowFilesDropped** 用於原生作業系統的檔案拖放；**WindowDropZoneFilesDropped** 用於網頁式拖放區域
- 某些事件僅適用於特定平台（例如，WindowDPIChanged 適用於 Windows，而縮放事件主要適用於 macOS）

## 事件命名慣例

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

## 效能考量

### 高頻事件的防彈跳處理

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

### 事件節流

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
