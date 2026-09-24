---
title: "API событий"
description: "Полный справочник по API событий"
slug: "reference/events"
sourcePath: "reference/events.md"
---

## Обзор

API событий предоставляет методы для отправки и прослушивания событий, обеспечивая обмен данными между различными частями приложения.

**Типы событий:**

- **События приложения** — события жизненного цикла приложения (запуск, завершение работы)
- **События окна** — изменения состояния окна (получение и потеря фокуса, изменение размера)
- **Пользовательские события** — определяемые пользователем события для обмена данными с учётом особенностей приложения

**Схемы обмена данными:**

- **Из Go во фронтенд** — отправка событий из Go и их прослушивание в JavaScript
- **Из фронтенда в Go** — напрямую не поддерживается (вместо этого используйте привязки сервисов)
- **Между частями фронтенда** — через Go или локальные события среды выполнения
- **Между окнами** — адресная отправка в определённые окна или широковещательная отправка во все окна

## Методы событий (Go)

### app.Event.Emit()

Отправляет пользовательское событие всем окнам. Возвращает `true`, если обработчик-перехватчик отменил отправку.

```go
func (em *EventManager) Emit(name string, data ...any) bool
```

**Параметры:**

- `name` — имя события
- `data` — необязательные данные, отправляемые вместе с событием

**Пример:**

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

Прослушивает пользовательские события в Go.

```go
func (em *EventManager) On(name string, callback func(*CustomEvent)) func()
```

**Параметры:**

- `name` — имя прослушиваемого события
- `callback` — функция, вызываемая при отправке события

**Возвращает:** функцию очистки для удаления слушателя события

**Пример:**

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

### События определённого окна

Отправка событий определённому окну:

```go
// Emit to specific window
window.EmitEvent("notification", "Hello from Go!")

// Emit to all windows
app.Event.Emit("global-update", data)
```

## Методы событий (фронтенд)

### On()

Прослушивает события из Go.

```javascript
import { Events } from '@wailsio/runtime'

Events.On(eventName, callback)
```

**Параметры:**

- `eventName` — имя прослушиваемого события
- `callback` — функция, вызываемая при получении события

**Возвращает:** функцию очистки

**Пример:**

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

Прослушивает одно срабатывание события.

```javascript
import { Events } from '@wailsio/runtime'

Events.Once(eventName, callback)
```

**Пример:**

```javascript
import { Events } from '@wailsio/runtime'

// Listen for first occurrence only
Events.Once('initialization-complete', (data) => {
    console.log('App initialized!', data)
    // This will only fire once
})
```

### Off()

Удаляет всех слушателей одного или нескольких событий. `Off` принимает переменное число строк с именами событий и **не** принимает функцию обратного вызова. Чтобы удалить отдельного слушателя, сохраните функцию отмены подписки, возвращённую `Events.On(...)`, и вызовите её.

```typescript
import { Events } from '@wailsio/runtime'

Events.Off(...eventNames: string[]): void
```

**Пример:**

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

Удаляет **всех** слушателей событий. Не принимает аргументов.

```typescript
import { Events } from '@wailsio/runtime'

Events.OffAll(): void
```

**Пример:**

```javascript
import { Events } from '@wailsio/runtime'

// Remove all listeners — typically used during teardown.
Events.OffAll()
```

### OnMultiple()

Прослушивает событие не более `max` раз, а затем автоматически отменяет подписку.

```typescript
Events.OnMultiple(eventName: string, callback, max: number): () => void
```

**Пример:**

```javascript
Events.OnMultiple('progress', (data) => {
    console.log('progress', data)
}, 5)
```

## События приложения

### app.Event.OnApplicationEvent()

Прослушивает события жизненного цикла приложения.

```go
func (em *EventManager) OnApplicationEvent(
    eventType events.ApplicationEventType,
    callback func(*ApplicationEvent),
) func()
```

Константы событий находятся в пакете `events`: `events.Common.*` содержит кроссплатформенные события, а `events.Mac.*`, `events.Windows.*` и `events.Linux.*` — события для отдельных платформ.

**Распространённые события приложения:**

- `events.Common.ApplicationStarted` — запуск приложения завершён.
- `events.Common.ThemeChanged` — системная тема переключилась со светлой на тёмную или наоборот.
- `events.Common.ApplicationOpenedWithFile` — приложение запущено через ассоциацию с файлом.
- `events.Common.ApplicationLaunchedWithUrl` — приложение запущено через схему URL.

Общей константы события «завершение работы приложения» **не** существует: зарегистрируйте очистку при завершении работы через `application.Options.OnShutdown` или `app.OnShutdown(func())`.

**Пример:**

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

## События окна

### OnWindowEvent()

Прослушивает события определённого окна.

```go
func (w *WebviewWindow) OnWindowEvent(
    eventType events.WindowEventType,
    callback func(*WindowEvent),
) func()
```

Константы событий окна находятся в пакете `events`: `events.Common.*` содержит кроссплатформенные события, а `events.Mac.*`, `events.Windows.*` и `events.Linux.*` — события для отдельных платформ.

**Распространённые события окна:**

- `events.Common.WindowFocus` — окно получило фокус.
- `events.Common.WindowLostFocus` — Окно потеряло фокус.
- `events.Common.WindowClosing` — Окно готовится к закрытию (закрытие можно отменить с помощью `RegisterHook`).
- `events.Common.WindowDidResize` — Размер окна был изменён.
- `events.Common.WindowDidMove` — Окно было перемещено.
- `events.Common.WindowMinimise` / `WindowUnMinimise` / `WindowMaximise` / `WindowUnMaximise` / `WindowFullscreen` / `WindowUnFullscreen`.
- `events.Common.WindowRuntimeReady` — Можно безопасно отправлять события среде выполнения внутри окна.

**Пример:**

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

## Распространённые шаблоны

Эти шаблоны демонстрируют проверенные подходы к использованию событий в реальных приложениях. Каждый шаблон решает определённую задачу взаимодействия между бэкендом на Go и фронтендом, помогая создавать отзывчивые и хорошо структурированные приложения.

### Шаблон «Запрос — ответ»

Используйте этот шаблон, когда нужно уведомить фронтенд о завершении операций бэкенда, например получения данных, обработки файлов или фоновых задач. Привязка сервиса возвращает данные напрямую, а события обеспечивают дополнительные уведомления для обновления пользовательского интерфейса, например показа всплывающих сообщений или обновления списков.

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

### Обновления хода выполнения

Идеально подходит для длительных операций, таких как загрузка файлов, пакетная обработка, импорт больших объёмов данных или кодирование видео. Во время операции отправляйте события о ходе выполнения, чтобы обновлять индикаторы выполнения, текст состояния или индикаторы этапов в пользовательском интерфейсе и предоставлять пользователям обратную связь в реальном времени.

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

### Взаимодействие между несколькими окнами

Отлично подходит для приложений с несколькими окнами, например панелями настроек, информационными панелями или средствами просмотра документов. Рассылайте события, чтобы синхронизировать состояние во всех окнах (изменения темы, пользовательские настройки), или отправляйте адресные события в определённые окна для их обновления.

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

### Синхронизация состояния

Используйте этот шаблон, когда необходимо синхронизировать состояние фронтенда и бэкенда, например пользовательские сеансы, конфигурацию приложения или функции совместной работы. При изменении состояния в бэкенде отправляйте события для обновления всех подключённых фронтендов, обеспечивая согласованность во всём приложении.

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

### Уведомления на основе событий

Лучше всего подходит для отображения обратной связи пользователю, например подтверждений успешного выполнения, предупреждений об ошибках или информационных сообщений. Вместо прямого вызова кода пользовательского интерфейса из сервисов отправляйте события уведомлений, которые фронтенд обрабатывает единообразно. Это упрощает изменение оформления уведомлений и добавление таких функций, как история уведомлений.

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

## Полный пример

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

## Встроенные события

Wails предоставляет встроенные системные события жизненного цикла приложения и окон. Эти события автоматически отправляются фреймворком.

### Общие и платформенные события

Wails предоставляет два типа системных событий:

**Общие события** (`events.Common.*`) — это кроссплатформенные абстракции, которые одинаково работают в macOS, Windows и Linux. Для максимальной переносимости приложения следует использовать именно эти события.

**Платформенные события** (`events.Mac.*`, `events.Windows.*`, `events.Linux.*`) — это базовые события, специфичные для операционной системы, на основе которых формируются общие события. Они предоставляют доступ к платформенно-зависимому поведению и особым случаям.

**Принцип работы:**

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

**Сопоставление событий:**

Платформенные события автоматически сопоставляются с общими событиями:

- macOS: `events.Mac.WindowShouldClose` → `events.Common.WindowClosing`
- Windows: `events.Windows.WindowClosing` → `events.Common.WindowClosing`
- Linux: `events.Linux.WindowDeleteEvent` → `events.Common.WindowClosing`

Это сопоставление выполняется автоматически в фоновом режиме, поэтому при подписке на `events.Common.WindowClosing` вы получите это событие независимо от платформы.

**Когда использовать каждый тип:**

- **Используйте общие события** в 99% кода приложения — они обеспечивают единообразное поведение на разных платформах
- **Используйте платформенные события** только тогда, когда необходимы платформенно-зависимые функции, недоступные через общие события (например, специфичные для macOS события жизненного цикла окна или события управления питанием Windows)

### События приложения

| Событие | Описание | Когда отправляется | Можно отменить |
| --- | --- | --- | --- |
| `ApplicationOpenedWithFile` | Приложение открыто с файлом | При запуске приложения с файлом (например, через ассоциацию файлов) | Нет |
| `ApplicationStarted` | Запуск приложения завершён | После завершения инициализации, когда приложение готово к работе | Нет |
| `ApplicationLaunchedWithUrl` | Приложение запущено с URL-адресом | При запуске приложения через схему URL | Нет |
| `ThemeChanged` | Системная тема изменена | При переключении темы ОС между светлым и тёмным режимами | Нет |
| `SystemWillSleep` | Система готовится перейти в спящий режим | Непосредственно перед переходом ОС в спящий режим (macOS / Windows / Linux с logind) | Нет |
| `SystemDidWake` | Система возобновила работу после сна | Сразу после выхода из спящего режима (macOS / Windows / Linux с logind) | Нет |

**Использование:**

```go
import "github.com/wailsapp/wails/v3/pkg/events"

app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(e *application.ApplicationEvent) {
    app.Logger.Info("Application ready!")
})

app.Event.OnApplicationEvent(events.Common.ThemeChanged, func(e *application.ApplicationEvent) {
    // Update app theme
})
```

### События окна

| Событие | Описание | Когда генерируется | Можно отменить |
| --- | --- | --- | --- |
| `WindowClosing` | Окно готовится к закрытию | Перед закрытием окна (пользователь нажал X или был вызван Close()) | Да |
| `WindowDidMove` | Окно перемещено в новое положение | После изменения положения окна (с устранением дребезга) | Нет |
| `WindowDidResize` | Размер окна изменён | После изменения размера окна | Нет |
| `WindowDPIChanged` | Масштабирование окна по DPI изменено | При перемещении между мониторами с разными значениями DPI (Windows) | Нет |
| `WindowFilesDropped` | Файлы перетащены средствами ОС | После перетаскивания файлов из ОС в окно | Нет |
| `WindowFocus` | Окно получило фокус | Когда окно становится активным | Нет |
| `WindowFullscreen` | Окно перешло в полноэкранный режим | После вызова Fullscreen() или перехода пользователя в полноэкранный режим | Нет |
| `WindowHide` | Окно скрыто | После вызова Hide() или перекрытия окна | Нет |
| `WindowLostFocus` | Окно потеряло фокус | Когда окно становится неактивным | Нет |
| `WindowMaximise` | Окно развёрнуто | После вызова Maximise() или разворачивания окна пользователем | Да (macOS) |
| `WindowMinimise` | Окно свёрнуто | После вызова Minimise() или сворачивания окна пользователем | Да (macOS) |
| `WindowRestore` | Окно восстановлено из свёрнутого или развёрнутого состояния | После вызова Restore() (преимущественно в Windows) | Нет |
| `WindowRuntimeReady` | Среда выполнения Wails загружена и готова к работе | После завершения инициализации среды выполнения JavaScript | Нет |
| `WindowShow` | Окно стало видимым | После вызова Show() или когда окно становится видимым | Нет |
| `WindowUnFullscreen` | Окно вышло из полноэкранного режима | После вызова UnFullscreen() или выхода пользователя из полноэкранного режима | Нет |
| `WindowUnMaximise` | Окно вышло из развёрнутого состояния | После вызова UnMaximise() или отмены развёртывания пользователем | Да (macOS) |
| `WindowUnMinimise` | Окно вышло из свёрнутого состояния | После вызова UnMinimise()/Restore() или восстановления окна пользователем | Да (macOS) |
| `WindowZoomIn` | Масштаб содержимого окна увеличен | После вызова ZoomIn() (преимущественно в macOS) | Да (macOS) |
| `WindowZoomOut` | Масштаб содержимого окна уменьшен | После вызова ZoomOut() (преимущественно в macOS) | Да (macOS) |
| `WindowZoomReset` | Масштаб содержимого окна сброшен до 100% | После вызова ZoomReset() (преимущественно в macOS) | Да (macOS) |
| `WindowDropZoneFilesDropped` | Файлы перетащены в зону сброса, заданную в JS | Когда файлы перетаскиваются на элемент с зоной сброса | Нет |

**Использование:**

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

**Важные замечания:**

- **WindowRuntimeReady** имеет критически важное значение: дождитесь этого события, прежде чем отправлять события во фронтенд
- Для событий **WindowDidMove** и **WindowDidResize** применяется устранение дребезга (по умолчанию 50 мс), чтобы предотвратить лавину событий
- **Отменяемые события** можно предотвратить, вызвав `event.Cancel()` в обработчике `RegisterHook()`
- **WindowFilesDropped** предназначено для нативного перетаскивания файлов средствами ОС, а **WindowDropZoneFilesDropped** — для веб-зон сброса
- Некоторые события зависят от платформы (например, WindowDPIChanged в Windows, а события масштабирования — преимущественно в macOS)

## Соглашения об именовании событий

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

## Рекомендации по производительности

### Устранение дребезга для высокочастотных событий

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

### Ограничение частоты событий

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
