---
title: "이벤트 시스템"
description: "이벤트 시스템을 사용하여 컴포넌트 간에 통신합니다"
slug: "features/events/system"
sourcePath: "features/events/system.md"
---

## 이벤트 시스템

Wails는 게시/구독 통신을 위한 <strong>통합 이벤트 시스템</strong>을 제공합니다. Go에서 JavaScript로, JavaScript에서 Go로, 창에서 다른 창으로 어디서든 이벤트를 발생시키고 수신할 수 있으므로, 형식이 지정된 이벤트와 수명 주기 훅을 활용하는 결합도 낮은 아키텍처를 구현할 수 있습니다.

## 빠른 시작

**Go(발생):**

```go
app.Event.Emit("user-logged-in", map[string]interface{}{
    "userId": 123,
    "name": "Alice",
})
```

**JavaScript(수신):**

```javascript
import { Events } from '@wailsio/runtime'

Events.On("user-logged-in", (event) => {
    console.log(`User ${event.data.name} logged in`)
})
```

**이게 전부입니다!** 언어 간 게시/구독이 가능합니다.

## 이벤트 유형

### 사용자 정의 이벤트

애플리케이션별 이벤트:

```go
// Emit from Go
app.Event.Emit("order-created", order)
app.Event.Emit("payment-processed", payment)
app.Event.Emit("notification", message)
```

```javascript
// Listen in JavaScript
Events.On("order-created", handleOrder)
Events.On("payment-processed", handlePayment)
Events.On("notification", showNotification)
```

### 시스템 이벤트

기본 제공되는 OS 및 애플리케이션 이벤트:

```go
import "github.com/wailsapp/wails/v3/pkg/events"

// Theme changes
app.Event.OnApplicationEvent(events.Common.ThemeChanged, func(e *application.ApplicationEvent) {
    if e.Context().IsDarkMode() {
        app.Logger.Info("Dark mode enabled")
    }
})

// Application lifecycle
app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(e *application.ApplicationEvent) {
    app.Logger.Info("Application started")
})
```

### 창 이벤트

창별 이벤트:

```go
window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
    app.Logger.Info("Window focused")
})

window.OnWindowEvent(events.Common.WindowClosing, func(e *application.WindowEvent) {
    app.Logger.Info("Window closing")
})
```

## 이벤트 발생시키기

### Go에서

**기본 발생:**

```go
app.Event.Emit("event-name", data)
```

**다양한 데이터 유형 사용:**

```go
// String
app.Event.Emit("message", "Hello")

// Number
app.Event.Emit("count", 42)

// Struct
app.Event.Emit("user", User{ID: 1, Name: "Alice"})

// Map
app.Event.Emit("config", map[string]interface{}{
    "theme": "dark",
    "fontSize": 14,
})

// Array
app.Event.Emit("items", []string{"a", "b", "c"})
```

**특정 창으로 전송:**

```go
window.EmitEvent("window-specific-event", data)
```

### JavaScript에서

```javascript
import { Events } from '@wailsio/runtime'

// Emit to Go
Events.Emit("button-clicked", { buttonId: "submit" })

// Emit to all windows
Events.Emit("broadcast-message", "Hello everyone")
```

## 이벤트 수신하기

### Go에서

**애플리케이션 이벤트:**

```go
app.Event.On("custom-event", func(e *application.CustomEvent) {
    data := e.Data
    // Handle event
})
```

**형식 단언 사용:**

```go
app.Event.On("user-updated", func(e *application.CustomEvent) {
    user := e.Data.(User)
    app.Logger.Info("User updated", "name", user.Name)
})
```

**여러 핸들러:**

```go
// All handlers will be called
app.Event.On("order-created", logOrder)
app.Event.On("order-created", sendEmail)
app.Event.On("order-created", updateInventory)
```

### JavaScript에서

**기본 리스너:**

```javascript
import { Events } from '@wailsio/runtime'

Events.On("event-name", (event) => {
    console.log("Event received:", event.data)
})
```

**정리 작업 포함:**

```javascript
const unsubscribe = Events.On("event-name", handleEvent)

// Later, stop listening
unsubscribe()
```

**여러 핸들러:**

```javascript
Events.On("data-updated", updateUI)
Events.On("data-updated", saveToCache)
Events.On("data-updated", logChange)
```

**일회성 핸들러:**

```javascript
Events.Once("data-updated", updateVariable)
```

## 시스템 이벤트

### 애플리케이션 이벤트

**공통 이벤트(크로스 플랫폼):**

```go
import "github.com/wailsapp/wails/v3/pkg/events"

// Application started
app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(e *application.ApplicationEvent) {
    app.Logger.Info("App started")
})

// Theme changed
app.Event.OnApplicationEvent(events.Common.ThemeChanged, func(e *application.ApplicationEvent) {
    isDark := e.Context().IsDarkMode()
    app.Event.Emit("theme-changed", isDark)
})

// System about to suspend
app.Event.OnApplicationEvent(events.Common.SystemWillSleep, func(e *application.ApplicationEvent) {
    flushPendingWrites()
})

// System resumed from suspend
app.Event.OnApplicationEvent(events.Common.SystemDidWake, func(e *application.ApplicationEvent) {
    reconnectSockets()
    refreshState()
})

// File opened
app.Event.OnApplicationEvent(events.Common.ApplicationOpenedWithFile, func(e *application.ApplicationEvent) {
    // Single-file: e.Context().Filename() returns ""; for multi-file launches
    // (Finder "Open With..." → multiple selection) use OpenedFiles().
    filePath := e.Context().Filename()
    openFile(filePath)
})
```

**플랫폼별 이벤트:**

@tabs{sync-key="platform"}
[macOS]
```go
// Application became active
app.Event.OnApplicationEvent(events.Mac.ApplicationDidBecomeActive, func(e *application.ApplicationEvent) {
    app.Logger.Info("App became active")
})

// Application will terminate
app.Event.OnApplicationEvent(events.Mac.ApplicationWillTerminate, func(e *application.ApplicationEvent) {
    cleanup()
})

// System about to sleep / resumed
app.Event.OnApplicationEvent(events.Mac.ApplicationWillSleep, func(e *application.ApplicationEvent) {
    flushPendingWrites()
})
app.Event.OnApplicationEvent(events.Mac.ApplicationDidWake, func(e *application.ApplicationEvent) {
    refreshState()
})

// Displays sleeping/waking — distinct from system sleep (e.g. lid lowered)
app.Event.OnApplicationEvent(events.Mac.ApplicationScreensDidSleep, func(e *application.ApplicationEvent) {
    pauseRendering()
})
app.Event.OnApplicationEvent(events.Mac.ApplicationScreensDidWake, func(e *application.ApplicationEvent) {
    resumeRendering()
})
```

[Windows]
```go
// Power status changed
app.Event.OnApplicationEvent(events.Windows.APMPowerStatusChange, func(e *application.ApplicationEvent) {
    app.Logger.Info("Power status changed")
})

// System suspending
app.Event.OnApplicationEvent(events.Windows.APMSuspend, func(e *application.ApplicationEvent) {
    saveState()
})

// System resumed. APMResumeAutomatic always fires on resume;
// APMResumeSuspend fires after it when the wake was triggered by
// user input. Prefer Common.SystemDidWake if you don't need to
// distinguish.
app.Event.OnApplicationEvent(events.Windows.APMResumeAutomatic, func(e *application.ApplicationEvent) {
    refreshState()
})
```

[Linux]
```go
// Application startup
app.Event.OnApplicationEvent(events.Linux.ApplicationStartup, func(e *application.ApplicationEvent) {
    app.Logger.Info("App starting")
})

// Theme changed
app.Event.OnApplicationEvent(events.Linux.SystemThemeChanged, func(e *application.ApplicationEvent) {
    updateTheme()
})

// System sleep/wake — sourced from systemd-logind's PrepareForSleep
// signal on the system bus. Requires logind/elogind exposing
// org.freedesktop.login1, so it may be unavailable on distros/setups
// without it (Alpine, Void, some Devuan setups).
app.Event.OnApplicationEvent(events.Linux.SystemWillSleep, func(e *application.ApplicationEvent) {
    flushPendingWrites()
})
app.Event.OnApplicationEvent(events.Linux.SystemDidWake, func(e *application.ApplicationEvent) {
    refreshState()
})
```

@end

### 창 이벤트

**공통 창 이벤트:**

```go
// Window focus
window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
    app.Logger.Info("Window focused")
})

// Window blur
window.OnWindowEvent(events.Common.WindowLostFocus, func(e *application.WindowEvent) {
    app.Logger.Info("Window blurred")
})

// Window closing
window.OnWindowEvent(events.Common.WindowClosing, func(e *application.WindowEvent) {
    if hasUnsavedChanges() {
        e.Cancel()  // Prevent close
    }
})

// Window closed
window.OnWindowEvent(events.Common.WindowClosing, func(e *application.WindowEvent) {
    cleanup()
})
```

## 이벤트 훅

훅은 표준 리스너보다 **먼저** 실행되며 이벤트를 <strong>취소</strong>할 수 있습니다.

```go
// Hook - runs first, can cancel
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    if hasUnsavedChanges() {
        result := showConfirmdialog("Unsaved changes. Close anyway?")
        if result != "yes" {
            e.Cancel()  // Prevent window close
        }
    }
})

// Standard listener - runs after hooks
window.OnWindowEvent(events.Common.WindowClosing, func(e *application.WindowEvent) {
    app.Logger.Info("Window closing")
})
```

**주요 차이점:**

| 기능 | 훅 | 표준 리스너 |
| --- | --- | --- |
| 실행 순서 | 먼저 실행되며, 등록 순서를 따름 | 훅 이후 실행되며 순서는 보장되지 않음 |
| 블로킹 | 동기식이며 다음 훅을 차단 | 비동기식이며 차단하지 않음 |
| 취소 가능 여부 | 가능 | 불가능(이미 전파됨) |
| 사용 사례 | 제어 흐름, 유효성 검사 | 로깅, 부수 효과 |

## 이벤트 패턴

### 게시/구독 패턴

```go
// Publisher (service)
type OrderService struct {
    app *application.App
}

func (o *OrderService) CreateOrder(items []Item) (*Order, error) {
    order := &Order{Items: items}
    
    if err := o.saveOrder(order); err != nil {
        return nil, err
    }
    
    // Publish event
    o.app.Event.Emit("order-created", order)
    
    return order, nil
}

// Subscribers
app.Event.On("order-created", func(e *application.CustomEvent) {
    order := e.Data.(*Order)
    sendConfirmationEmail(order)
})

app.Event.On("order-created", func(e *application.CustomEvent) {
    order := e.Data.(*Order)
    updateInventory(order)
})

app.Event.On("order-created", func(e *application.CustomEvent) {
    order := e.Data.(*Order)
    logOrder(order)
})
```

### 요청/응답 패턴

```go
// Frontend requests data
Emit("get-user-data", { userId: 123 })

// Backend responds
app.Event.On("get-user-data", func(e *application.CustomEvent) {
    data := e.Data.(map[string]interface{})
    userId := int(data["userId"].(float64))
    
    user := getUserFromDB(userId)
    
    // Send response
    app.Event.Emit("user-data-response", user)
})
```

```javascript
// Frontend receives response
Events.On("user-data-response", (event) => {
    const user = event.data
    displayUser(user)
})
```

**참고:** 요청/응답에는 **바인딩이 더 적합합니다**. 알림에는 이벤트를 사용하세요.

### 브로드캐스트 패턴

```go
// Broadcast to all windows
app.Event.Emit("global-notification", "System update available")

// Each window handles it
Events.On("global-notification", (event) => {
    const message = event.data
    showNotification(message)
})
```

### 이벤트 집계

```go
type EventAggregator struct {
    events []Event
    mu     sync.Mutex
}

func (ea *EventAggregator) Add(event Event) {
    ea.mu.Lock()
    defer ea.mu.Unlock()
    
    ea.events = append(ea.events, event)
    
    // Emit batch every 100 events
    if len(ea.events) >= 100 {
        app.Event.Emit("event-batch", ea.events)
        ea.events = nil
    }
}
```

## 전체 예제

**Go:**

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)

type NotificationService struct {
    app *application.App
}

func (n *NotificationService) Notify(message string) {
    // Emit to all windows
    n.app.Event.Emit("notification", map[string]interface{}{
        "message":   message,
        "timestamp": time.Now(),
    })
}

func main() {
    app := application.New(application.Options{
        Name: "Event Demo",
    })
    
    notifService := &NotificationService{app: app}
    
    // System events
    app.Event.OnApplicationEvent(events.Common.ThemeChanged, func(e *application.ApplicationEvent) {
        isDark := e.Context().IsDarkMode()
        app.Event.Emit("theme-changed", isDark)
    })
    
    // Custom events from frontend
    app.Event.On("user-action", func(e *application.CustomEvent) {
        data := e.Data.(map[string]interface{})
        action := data["action"].(string)
        
        app.Logger.Info("User action", "action", action)
        
        // Respond
        notifService.Notify("Action completed: " + action)
    })
    
    // Window events
    window := app.Window.New()
    
    window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
        app.Event.Emit("window-focused", window.Name())
    })
    
    window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
        // Confirm before close
        app.Event.Emit("confirm-close", nil)
        e.Cancel()  // Wait for confirmation
    })
    
    app.Run()
}
```

**JavaScript:**

```javascript
import { Events } from '@wailsio/runtime'

// Listen for notifications
Events.On("notification", (event) => {
    showNotification(event.data.message)
})

// Listen for theme changes
Events.On("theme-changed", (event) => {
    const isDark = event.data
    document.body.classList.toggle('dark', isDark)
})

// Listen for window focus
Events.On("window-focused", (event) => {
    const windowName = event.data
    console.log(`Window ${windowName} focused`)
})

// Handle close confirmation
Events.On("confirm-close", (event) => {
    if (confirm("Close window?")) {
        Events.Emit("close-confirmed", true)
    }
})

// Emit user actions
document.getElementById('button').addEventListener('click', () => {
    Events.Emit("user-action", { action: "button-clicked" })
})
```

## 모범 사례

### ✅ 권장 사항

- **알림에는 이벤트를 사용하세요** - 단방향 통신
- **요청에는 바인딩을 사용하세요** - 양방향 통신
- **이벤트 이름을 일관되게 유지하세요** - kebab-case를 사용하세요
- **이벤트 데이터를 문서화하세요** - 어떤 필드가 포함되나요?
- **사용이 끝나면 구독을 해제하세요** - 메모리 누수를 방지합니다
- **유효성 검사에는 훅을 사용하세요** - 이벤트 흐름을 제어합니다

### ❌ 금지 사항

- **RPC에 이벤트를 사용하지 마세요** - 대신 바인딩을 사용하세요
- **너무 자주 내보내지 마세요** - 필요한 경우 일괄 처리하세요
- **핸들러에서 실행을 차단하지 마세요** - 빠르게 처리되도록 유지하세요
- **구독 해제를 잊지 마세요** - 메모리 누수가 발생합니다
- **대용량 데이터에 이벤트를 사용하지 마세요** - 바인딩을 사용하세요
- **이벤트 루프를 만들지 마세요** - A가 B를 내보내고 B가 A를 내보내는 방식입니다

## 다음 단계

@cards{cols="2"}
★ 이벤트 가이드
이벤트 패턴과 형식이 안전한 이벤트 생성에 대해 알아보세요.

[자세히 알아보기 →](/guides/events-reference/)

---
◆ 이벤트 API
애플리케이션 및 창 이벤트 API 전체를 살펴보세요.

[자세히 알아보기 →](/reference/events/)

---
🚀 바인딩
요청/응답에는 바인딩을 사용하세요.

[자세히 알아보기 →](/features/bindings/methods/)

---
▣ 창 이벤트
창 수명 주기 이벤트를 처리하세요.

[자세히 알아보기 →](/features/windows/events/)

@end

---

**궁금한 점이 있으신가요?** [Discord](https://discord.gg/JDdSxwjhGf)에서 질문하거나 [이벤트 예제](https://github.com/wailsapp/wails/tree/master/v3/examples/events)를 확인하세요.
