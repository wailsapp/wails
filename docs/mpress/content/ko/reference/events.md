---
title: "이벤트 API"
description: "이벤트 API 전체 레퍼런스"
slug: "reference/events"
sourcePath: "reference/events.md"
---

## 개요

이벤트 API는 이벤트를 발생시키고 수신하는 메서드를 제공하여 애플리케이션의 여러 부분이 서로 통신할 수 있게 합니다.

**이벤트 유형:**

- **애플리케이션 이벤트** - 앱 수명 주기 이벤트(시작, 종료)
- **창 이벤트** - 창 상태 변경(포커스 획득, 포커스 상실, 크기 조정)
- **사용자 정의 이벤트** - 앱별 통신을 위해 사용자가 정의하는 이벤트

**통신 패턴:**

- **Go에서 프런트엔드로** - Go에서 이벤트를 발생시키고 JavaScript에서 수신
- **프런트엔드에서 Go로** - 직접 통신할 수 없음(대신 서비스 바인딩 사용)
- **프런트엔드에서 프런트엔드로** - Go 또는 로컬 런타임 이벤트를 통해 통신
- **창에서 창으로** - 특정 창을 대상으로 지정하거나 모든 창에 브로드캐스트

## 이벤트 메서드(Go)

### app.Event.Emit()

모든 창에 사용자 정의 이벤트를 발생시킵니다. 훅이 이벤트 발생을 취소하면 `true`을 반환합니다.

```go
func (em *EventManager) Emit(name string, data ...any) bool
```

**매개변수:**

- `name` - 이벤트 이름
- `data` - 이벤트와 함께 보낼 선택적 데이터

**예제:**

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

Go에서 사용자 정의 이벤트를 수신합니다.

```go
func (em *EventManager) On(name string, callback func(*CustomEvent)) func()
```

**매개변수:**

- `name` - 수신할 이벤트의 이름
- `callback` - 이벤트가 발생할 때 호출되는 함수

**반환값:** 이벤트 리스너를 제거하는 정리 함수

**예제:**

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

### 창별 이벤트

특정 창에 이벤트를 발생시킵니다:

```go
// Emit to specific window
window.EmitEvent("notification", "Hello from Go!")

// Emit to all windows
app.Event.Emit("global-update", data)
```

## 이벤트 메서드(프런트엔드)

### On()

Go에서 보낸 이벤트를 수신합니다.

```javascript
import { Events } from '@wailsio/runtime'

Events.On(eventName, callback)
```

**매개변수:**

- `eventName` - 수신할 이벤트의 이름
- `callback` - 이벤트를 수신할 때 호출되는 함수

**반환값:** 정리 함수

**예제:**

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

이벤트를 한 번만 수신합니다.

```javascript
import { Events } from '@wailsio/runtime'

Events.Once(eventName, callback)
```

**예제:**

```javascript
import { Events } from '@wailsio/runtime'

// Listen for first occurrence only
Events.Once('initialization-complete', (data) => {
    console.log('App initialized!', data)
    // This will only fire once
})
```

### Off()

하나 이상의 이벤트에 등록된 모든 리스너를 제거합니다. `Off`은 가변 개수의 이벤트 이름 문자열을 받으며, 콜백은 받지 **않습니다**. 리스너 하나만 제거하려면 `Events.On(...)`이 반환한 구독 해제 함수를 보관해 두었다가 호출합니다.

```typescript
import { Events } from '@wailsio/runtime'

Events.Off(...eventNames: string[]): void
```

**예제:**

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

이벤트 리스너를 **모두** 제거합니다. 인수를 받지 않습니다.

```typescript
import { Events } from '@wailsio/runtime'

Events.OffAll(): void
```

**예제:**

```javascript
import { Events } from '@wailsio/runtime'

// Remove all listeners — typically used during teardown.
Events.OffAll()
```

### OnMultiple()

이벤트를 최대 `max`회 수신한 후 자동으로 구독을 해제합니다.

```typescript
Events.OnMultiple(eventName: string, callback, max: number): () => void
```

**예제:**

```javascript
Events.OnMultiple('progress', (data) => {
    console.log('progress', data)
}, 5)
```

## 애플리케이션 이벤트

### app.Event.OnApplicationEvent()

애플리케이션 수명 주기 이벤트를 수신합니다.

```go
func (em *EventManager) OnApplicationEvent(
    eventType events.ApplicationEventType,
    callback func(*ApplicationEvent),
) func()
```

이벤트 상수는 `events` 패키지에 있습니다. 크로스 플랫폼 이벤트에는 `events.Common.*`을 사용하고, 플랫폼별 이벤트에는 `events.Mac.*` / `events.Windows.*` / `events.Linux.*`을 사용합니다.

**일반적인 애플리케이션 이벤트:**

- `events.Common.ApplicationStarted` - 애플리케이션 시작이 완료되었습니다.
- `events.Common.ThemeChanged` - 시스템 테마가 라이트 모드와 다크 모드 사이에서 전환되었습니다.
- `events.Common.ApplicationOpenedWithFile` - 파일 연결을 통해 실행되었습니다.
- `events.Common.ApplicationLaunchedWithUrl` - URL 스킴을 통해 실행되었습니다.

일반적인 "애플리케이션 종료" 이벤트 상수는 **없습니다**. 종료 시 정리 작업은 `application.Options.OnShutdown` 또는 `app.OnShutdown(func())`을 통해 등록합니다.

**예제:**

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

## 창 이벤트

### OnWindowEvent()

창별 이벤트를 수신합니다.

```go
func (w *WebviewWindow) OnWindowEvent(
    eventType events.WindowEventType,
    callback func(*WindowEvent),
) func()
```

창 이벤트 상수는 `events` 패키지에 있습니다. 크로스 플랫폼 이벤트에는 `events.Common.*`을 사용하고, 플랫폼별 이벤트에는 `events.Mac.*` / `events.Windows.*` / `events.Linux.*`을 사용합니다.

**일반적인 창 이벤트:**

- `events.Common.WindowFocus` - 창이 포커스를 얻었습니다.
- `events.Common.WindowLostFocus` - 창이 포커스를 잃었습니다.
- `events.Common.WindowClosing` - 창이 닫히기 직전입니다(`RegisterHook`을 통해 취소할 수 있음).
- `events.Common.WindowDidResize` - 창 크기가 조정되었습니다.
- `events.Common.WindowDidMove` - 창이 이동되었습니다.
- `events.Common.WindowMinimise` / `WindowUnMinimise` / `WindowMaximise` / `WindowUnMaximise` / `WindowFullscreen` / `WindowUnFullscreen`.
- `events.Common.WindowRuntimeReady` - 창 내 런타임으로 이벤트를 안전하게 내보낼 수 있습니다.

**예:**

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

## 일반적인 패턴

다음 패턴은 실제 애플리케이션에서 이벤트를 사용하는 검증된 방식을 보여 줍니다. 각 패턴은 Go 백엔드와 프런트엔드 간의 특정 통신 문제를 해결하여 응답성이 뛰어나고 구조가 잘 갖춰진 애플리케이션을 구축하는 데 도움이 됩니다.

### 요청/응답 패턴

데이터 가져오기, 파일 처리 또는 백그라운드 작업 등이 끝난 후 백엔드 작업의 완료를 프런트엔드에 알릴 때 사용합니다. 서비스 바인딩은 데이터를 직접 반환하며, 이벤트는 토스트 메시지 표시나 목록 새로 고침과 같은 UI 업데이트를 위한 추가 알림을 제공합니다.

**Go:**

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

**JavaScript:**

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

### 진행률 업데이트

파일 업로드, 일괄 처리, 대용량 데이터 가져오기 또는 동영상 인코딩처럼 오래 실행되는 작업에 적합합니다. 작업 중에 진행률 이벤트를 내보내 UI의 진행률 표시줄, 상태 텍스트 또는 단계 표시기를 업데이트하고 사용자에게 실시간 피드백을 제공합니다.

**Go:**

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

**JavaScript:**

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

### 다중 창 통신

설정 패널, 대시보드 또는 문서 뷰어처럼 여러 창이 있는 애플리케이션에 적합합니다. 이벤트를 브로드캐스트하여 모든 창에서 상태(테마 변경, 사용자 환경 설정)를 동기화하거나, 특정 창에 대상 지정 이벤트를 보내 해당 창만 업데이트합니다.

**Go:**

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

**JavaScript:**

```javascript
import { Events } from '@wailsio/runtime'

// Listen in any window
Events.On('theme-changed', (theme) => {
    document.body.className = theme
})
```

### 상태 동기화

사용자 세션, 애플리케이션 구성 또는 협업 기능처럼 프런트엔드와 백엔드의 상태를 동기화해야 할 때 사용합니다. 백엔드의 상태가 변경되면 이벤트를 내보내 연결된 모든 프런트엔드를 업데이트하여 애플리케이션 전체의 일관성을 보장합니다.

**Go:**

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

**JavaScript:**

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

### 이벤트 기반 알림

성공 확인, 오류 경고 또는 정보 메시지와 같은 사용자 피드백을 표시하는 데 가장 적합합니다. 서비스에서 UI 코드를 직접 호출하는 대신 프런트엔드가 일관되게 처리하는 알림 이벤트를 내보내면, 알림 스타일을 변경하거나 알림 기록과 같은 기능을 쉽게 추가할 수 있습니다.

**Go:**

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

**JavaScript:**

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

## 전체 예제

**Go:**

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

**JavaScript:**

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

## 기본 제공 이벤트

Wails는 애플리케이션 및 창 수명 주기를 위한 기본 제공 시스템 이벤트를 제공합니다. 이러한 이벤트는 프레임워크에서 자동으로 내보냅니다.

### 공통 이벤트와 플랫폼 네이티브 이벤트

Wails는 두 가지 유형의 시스템 이벤트를 제공합니다.

**공통 이벤트**(`events.Common.*`)는 macOS, Windows 및 Linux에서 일관되게 작동하는 크로스 플랫폼 추상화입니다. 이식성을 극대화하려면 애플리케이션에서 이러한 이벤트를 사용하는 것이 좋습니다.

**플랫폼 네이티브 이벤트**(`events.Mac.*`, `events.Windows.*`, `events.Linux.*`)는 공통 이벤트로 매핑되는 기반 OS별 이벤트입니다. 이를 통해 플랫폼별 동작과 예외적인 상황에 접근할 수 있습니다.

**작동 방식:**

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

**이벤트 매핑:**

플랫폼 네이티브 이벤트는 공통 이벤트에 자동으로 매핑됩니다.

- macOS: `events.Mac.WindowShouldClose` → `events.Common.WindowClosing`
- Windows: `events.Windows.WindowClosing` → `events.Common.WindowClosing`
- Linux: `events.Linux.WindowDeleteEvent` → `events.Common.WindowClosing`

이 매핑은 백그라운드에서 자동으로 이루어지므로 `events.Common.WindowClosing`을 수신하도록 등록하면 플랫폼과 관계없이 해당 이벤트를 받습니다.

**각 유형을 사용해야 하는 경우:**

- 애플리케이션 코드의 99%에는 **공통 이벤트를 사용하세요**. 공통 이벤트는 플랫폼 전반에서 일관된 동작을 제공합니다.
- 공통 이벤트에서 제공하지 않는 플랫폼별 기능이 필요한 경우에만 **플랫폼 네이티브 이벤트를 사용하세요**(예: macOS 전용 창 수명 주기 이벤트, Windows 전원 관리 이벤트).

### 애플리케이션 이벤트

| 이벤트 | 설명 | 발생 시점 | 취소 가능 여부 |
| --- | --- | --- | --- |
| `ApplicationOpenedWithFile` | 파일과 함께 애플리케이션을 열었음 | 파일 연결 등을 통해 파일과 함께 앱이 실행될 때 | 아니요 |
| `ApplicationStarted` | 애플리케이션 시작이 완료됨 | 앱 초기화가 완료되어 사용할 준비가 된 후 | 아니요 |
| `ApplicationLaunchedWithUrl` | URL로 애플리케이션이 실행됨 | URL 스킴을 통해 앱이 실행될 때 | 아니요 |
| `ThemeChanged` | 시스템 테마가 변경됨 | OS 테마가 라이트 모드와 다크 모드 사이에서 전환될 때 | 아니요 |
| `SystemWillSleep` | 시스템이 곧 절전 모드로 전환됨 | OS가 절전 모드로 전환되기 직전(macOS / Windows / logind를 사용하는 Linux) | 아니요 |
| `SystemDidWake` | 시스템이 절전 모드에서 복귀함 | 절전 모드에서 복귀하는 즉시(macOS / Windows / logind를 사용하는 Linux) | 아니요 |

**사용법:**

```go
import "github.com/wailsapp/wails/v3/pkg/events"

app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(e *application.ApplicationEvent) {
    app.Logger.Info("Application ready!")
})

app.Event.OnApplicationEvent(events.Common.ThemeChanged, func(e *application.ApplicationEvent) {
    // Update app theme
})
```

### 창 이벤트

| 이벤트 | 설명 | 발생 시점 | 취소 가능 여부 |
| --- | --- | --- | --- |
| `WindowClosing` | 창이 곧 닫힘 | 창이 닫히기 전(사용자가 X 버튼을 클릭하거나 Close()가 호출됨) | 예 |
| `WindowDidMove` | 창이 새 위치로 이동함 | 창 위치가 변경된 후(디바운싱됨) | 아니요 |
| `WindowDidResize` | 창 크기가 조정됨 | 창 크기가 변경된 후 | 아니요 |
| `WindowDPIChanged` | 창의 DPI 배율이 변경됨 | DPI가 서로 다른 모니터 사이로 창을 이동할 때(Windows) | 아니요 |
| `WindowFilesDropped` | 네이티브 OS 끌어서 놓기로 파일이 놓임 | OS에서 창으로 파일을 끌어다 놓은 후 | 아니요 |
| `WindowFocus` | 창이 포커스를 얻음 | 창이 활성화될 때 | 아니요 |
| `WindowFullscreen` | 창이 전체 화면 모드로 전환됨 | Fullscreen()이 호출되거나 사용자가 전체 화면 모드로 전환한 후 | 아니요 |
| `WindowHide` | 창이 숨겨짐 | Hide()가 호출되거나 창이 다른 창에 가려진 후 | 아니요 |
| `WindowLostFocus` | 창이 포커스를 잃음 | 창이 비활성화될 때 | 아니요 |
| `WindowMaximise` | 창이 최대화됨 | Maximise()가 호출되거나 사용자가 창을 최대화한 후 | 예(macOS) |
| `WindowMinimise` | 창이 최소화됨 | Minimise()가 호출되거나 사용자가 창을 최소화한 후 | 예(macOS) |
| `WindowRestore` | 최소화 또는 최대화 상태에서 창이 복원됨 | Restore()가 호출된 후(주로 Windows) | 아니요 |
| `WindowRuntimeReady` | Wails 런타임이 로드되어 사용할 준비가 됨 | JavaScript 런타임 초기화가 완료될 때 | 아니요 |
| `WindowShow` | 창이 표시됨 | Show() 호출 후 또는 창이 표시될 때 | 아니요 |
| `WindowUnFullscreen` | 창이 전체 화면 모드에서 해제됨 | UnFullscreen() 호출 후 또는 사용자가 전체 화면 모드를 종료할 때 | 아니요 |
| `WindowUnMaximise` | 창의 최대화 상태가 해제됨 | UnMaximise() 호출 후 또는 사용자가 최대화를 해제할 때 | 예(macOS) |
| `WindowUnMinimise` | 창의 최소화 상태가 해제됨 | UnMinimise()/Restore() 호출 후 또는 사용자가 창을 복원할 때 | 예(macOS) |
| `WindowZoomIn` | 창 콘텐츠 확대 비율이 증가함 | ZoomIn() 호출 후(주로 macOS) | 예(macOS) |
| `WindowZoomOut` | 창 콘텐츠 확대 비율이 감소함 | ZoomOut() 호출 후(주로 macOS) | 예(macOS) |
| `WindowZoomReset` | 창 콘텐츠 확대 비율이 100%로 재설정됨 | ZoomReset() 호출 후(주로 macOS) | 예(macOS) |
| `WindowDropZoneFilesDropped` | JS에서 정의한 드롭 영역에 파일이 드롭됨 | 드롭 영역이 있는 요소에 파일을 드롭할 때 | 아니요 |

**사용법:**

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

**중요 참고 사항:**

- <strong>WindowRuntimeReady</strong>는 매우 중요합니다. 프런트엔드로 이벤트를 내보내기 전에 이 이벤트를 기다리세요.
- **WindowDidMove** 및 **WindowDidResize** 이벤트에는 이벤트 폭주를 방지하기 위해 디바운싱이 적용됩니다(기본값: 50ms).
- <strong>취소 가능한 이벤트</strong>는 `RegisterHook()` 핸들러에서 `event.Cancel()`을 호출하여 차단할 수 있습니다.
- <strong>WindowFilesDropped</strong>는 네이티브 OS 파일 드롭용이며, <strong>WindowDropZoneFilesDropped</strong>는 웹 기반 드롭 영역용입니다.
- 일부 이벤트는 플랫폼별로 제공됩니다(예: Windows의 WindowDPIChanged, 주로 macOS에서 제공되는 확대/축소 이벤트).

## 이벤트 명명 규칙

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

## 성능 고려 사항

### 빈도가 높은 이벤트 디바운싱

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

### 이벤트 스로틀링

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
