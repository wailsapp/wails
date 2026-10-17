---
title: "Система событий"
description: "Обмен данными между компонентами с помощью системы событий"
slug: "features/events/system"
sourcePath: "features/events/system.md"
---

## Система событий

Wails предоставляет **единую систему событий** для обмена данными по модели публикации и подписки. Генерируйте события где угодно и обрабатывайте их где угодно — из Go в JavaScript, из JavaScript в Go, из одного окна в другое. Это позволяет создавать слабосвязанную архитектуру с типизированными событиями и хуками жизненного цикла.

## Быстрый старт

**Go (генерация события):**

```go
app.Event.Emit("user-logged-in", map[string]interface{}{
    "userId": 123,
    "name": "Alice",
})
```

**JavaScript (обработка события):**

```javascript
import { Events } from '@wailsio/runtime'

Events.On("user-logged-in", (event) => {
    console.log(`User ${event.data.name} logged in`)
})
```

**Вот и всё!** Публикация и подписка между языками.

## Типы событий

### Пользовательские события

События, относящиеся к вашему приложению:

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

### Системные события

Встроенные события операционной системы и приложения:

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

### События окна

События, относящиеся к окну:

```go
window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
    app.Logger.Info("Window focused")
})

window.OnWindowEvent(events.Common.WindowClosing, func(e *application.WindowEvent) {
    app.Logger.Info("Window closing")
})
```

## Генерация событий

### Из Go

**Базовая генерация события:**

```go
app.Event.Emit("event-name", data)
```

**С данными разных типов:**

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

**Для конкретного окна:**

```go
window.EmitEvent("window-specific-event", data)
```

### Из JavaScript

```javascript
import { Events } from '@wailsio/runtime'

// Emit to Go
Events.Emit("button-clicked", { buttonId: "submit" })

// Emit to all windows
Events.Emit("broadcast-message", "Hello everyone")
```

## Обработка событий

### В Go

**События приложения:**

```go
app.Event.On("custom-event", func(e *application.CustomEvent) {
    data := e.Data
    // Handle event
})
```

**С утверждением типа:**

```go
app.Event.On("user-updated", func(e *application.CustomEvent) {
    user := e.Data.(User)
    app.Logger.Info("User updated", "name", user.Name)
})
```

**Несколько обработчиков:**

```go
// All handlers will be called
app.Event.On("order-created", logOrder)
app.Event.On("order-created", sendEmail)
app.Event.On("order-created", updateInventory)
```

### В JavaScript

**Базовый обработчик:**

```javascript
import { Events } from '@wailsio/runtime'

Events.On("event-name", (event) => {
    console.log("Event received:", event.data)
})
```

**С очисткой:**

```javascript
const unsubscribe = Events.On("event-name", handleEvent)

// Later, stop listening
unsubscribe()
```

**Несколько обработчиков:**

```javascript
Events.On("data-updated", updateUI)
Events.On("data-updated", saveToCache)
Events.On("data-updated", logChange)
```

**Одноразовые обработчики:**

```javascript
Events.Once("data-updated", updateVariable)
```

## Системные события

### События приложения

**Общие события (кроссплатформенные):**

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

**События для отдельных платформ:**

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

### События окна

**Общие события окна:**

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

## Хуки событий

Хуки выполняются **до** стандартных обработчиков и могут **отменять** события:

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

**Основные различия:**

| Характеристика | Хуки | Стандартные обработчики |
| --- | --- | --- |
| Порядок выполнения | Сначала, в порядке регистрации | После хуков, порядок не гарантирован |
| Блокировка | Синхронные, блокируют следующий хук | Асинхронные, неблокирующие |
| Возможность отмены | Да | Нет (событие уже распространено) |
| Сценарий использования | Управление потоком, проверка данных | Ведение журнала, побочные эффекты |

## Шаблоны использования событий

### Шаблон публикации и подписки

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

### Шаблон «запрос — ответ»

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

**Примечание:** для запросов и ответов **лучше подходят привязки**. Используйте события для уведомлений.

### Шаблон широковещательной рассылки

```go
// Broadcast to all windows
app.Event.Emit("global-notification", "System update available")

// Each window handles it
Events.On("global-notification", (event) => {
    const message = event.data
    showNotification(message)
})
```

### Агрегация событий

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

## Полный пример

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

## Рекомендации

### ✅ Следует

- **Используйте события для уведомлений** — односторонний обмен данными
- **Используйте привязки для запросов** — двусторонний обмен данными
- **Соблюдайте единообразие в именах событий** — используйте kebab-case
- **Документируйте данные событий** — какие поля в них входят?
- **Отписывайтесь, когда подписка больше не нужна** — предотвращайте утечки памяти
- **Используйте хуки для проверки данных** — управляйте потоком событий

### ❌ Не следует

- **Не используйте события для RPC** — вместо них используйте привязки
- **Не отправляйте события слишком часто** — при необходимости объединяйте их в пакеты
- **Не блокируйте выполнение в обработчиках** — они должны работать быстро
- **Не забывайте отписываться** — иначе возможны утечки памяти
- **Не используйте события для больших объёмов данных** — используйте привязки
- **Не создавайте циклы событий** — A отправляет B, а B отправляет A

## Дальнейшие шаги

@cards{cols="2"}
★ Руководство по событиям
Узнайте о шаблонах работы с событиями и генерации типобезопасных событий.

[Подробнее →](/guides/events-reference/)

---
◆ API событий
Изучите полный API событий приложения и окон.

[Подробнее →](/reference/events/)

---
🚀 Привязки
Используйте привязки для взаимодействия по схеме «запрос — ответ».

[Подробнее →](/features/bindings/methods/)

---
▣ События окна
Обрабатывайте события жизненного цикла окна.

[Подробнее →](/features/windows/events/)

@end

---

**Есть вопросы?** Задайте их в [Discord](https://discord.gg/JDdSxwjhGf) или ознакомьтесь с [примерами событий](https://github.com/wailsapp/wails/tree/master/v3/examples/events).
