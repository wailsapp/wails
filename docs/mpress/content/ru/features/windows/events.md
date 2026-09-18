---
title: "События окна"
description: "Обработка событий жизненного цикла окна и изменения его состояния"
slug: "features/windows/events"
sourcePath: "features/windows/events.md"
---

## События окна

Wails отправляет события жизненного цикла окна и изменения его состояния через единый API: `OnWindowEvent` — для пассивных слушателей, а `RegisterHook` — для отменяемых перехватчиков, которые могут предотвратить действие по умолчанию.

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

Оба вызова возвращают `unsubscribe func()`, который можно вызвать, чтобы удалить обработчик.

Кроссплатформенные события находятся в `events.Common.*`. События для отдельных платформ находятся в `events.Mac.*`, `events.Windows.*` и `events.Linux.*`. Полный список генерируется в `v3/pkg/events/events.go`.

## События жизненного цикла

### Создание окна

Чтобы выполнять функцию обратного вызова при каждом создании окна, используйте `app.Window.OnCreate`:

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

Функция обратного вызова получает `application.Window` (интерфейс). Функция обратного вызова `OnCreate` вызывается один раз для каждого окна после инициализации среды выполнения этого окна.

### WindowClosing

Срабатывает, когда пользователь пытается закрыть окно (нажимает X, ⌘W, Alt+F4 и т. д.).

Используйте **перехватчик** (`RegisterHook`) — только перехватчики могут отменить закрытие. Слушатели отслеживают событие, но не могут предотвратить закрытие.

```go
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    if hasUnsavedChanges() {
        // Cancel the close — the window stays open.
        e.Cancel()
    }
})
```

**Важно:**

- `WindowClosing` отправляется при попытках закрытия, инициированных пользователем.
- Функция обратного вызова `RegisterHook` может вызвать `e.Cancel()`, чтобы оставить окно открытым.
- Функции обратного вызова `OnWindowEvent` для `WindowClosing` являются наблюдателями: они срабатывают, но не могут отменить закрытие.

**Программное закрытие:** вызовите `window.Close()`. `window.Destroy()` не существует.

### WindowRuntimeReady

Срабатывает после завершения инициализации среды выполнения внутри окна — с этого момента можно безопасно обращаться к JS-контексту окна:

```go
window.OnWindowEvent(events.Common.WindowRuntimeReady, func(e *application.WindowEvent) {
    window.EmitEvent("app-ready", nil)
})
```

## События фокуса

### WindowFocus

Вызывается, когда окно получает фокус:

```go
window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
    fmt.Println("Window gained focus")
    updateTitleBar(true)
    app.Event.Emit("window-focused", window.ID())
})
```

### WindowLostFocus

Вызывается, когда окно теряет фокус:

```go
window.OnWindowEvent(events.Common.WindowLostFocus, func(e *application.WindowEvent) {
    fmt.Println("Window lost focus")
    updateTitleBar(false)
    saveCurrentState()
})
```

**Пример: интерфейс, учитывающий фокус:**

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

## События изменения состояния

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

Чтобы программно перейти в полноэкранный режим или выйти из него, используйте `window.Fullscreen()` / `window.UnFullscreen()` / `window.ToggleFullscreen()` — `SetFullscreen(bool)` не существует. Проверяйте состояние с помощью `window.IsFullscreen() bool`.

## События изменения положения и размера

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

Сами события `WindowDidResize` и `WindowDidMove` не содержат координат — получайте их из окна с помощью `window.Size()` / `window.Position()` внутри функции обратного вызова.

**Пример: адаптивная компоновка:**

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

## Полный пример

Готовое к эксплуатации окно с полной обработкой событий:

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

## Координация событий

### События между окнами

Координируйте работу нескольких окон с помощью шины событий приложения:

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

### Цепочки событий

```go
window.OnWindowEvent(events.Common.WindowMaximise, func(e *application.WindowEvent) {
    saveWindowState()
    window.EmitEvent("layout-changed", "maximised")
    app.Event.Emit("window-maximised", window.ID())
})
```

### События с устранением дребезга

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

## Рекомендации

### ✅ Делайте так

- **Используйте перехватчики для отмены** — только функции обратного вызова `RegisterHook` могут вызывать `e.Cancel()`.
- **Сохраняйте состояние при закрытии** — восстанавливайте положение и размер окна при следующем запуске.
- **Устраняйте дребезг частых событий** — `WindowDidResize` и `WindowDidMove` срабатывают очень часто.
- **Обрабатывайте изменения фокуса** — соответствующим образом обновляйте интерфейс.
- **Координируйте работу через события** — используйте `app.Event.Emit` для обмена сообщениями между окнами.
- **Отписывайтесь**, когда обработчики больше не нужны: и `OnWindowEvent`, и `RegisterHook` возвращают функцию отписки `func()`.

### ❌ Не делайте так

- **Не блокируйте обработчики событий** — они должны выполняться быстро.
- **Не пытайтесь отменить событие из `OnWindowEvent`** — используйте `RegisterHook`.
- **Не используйте `window.Destroy()`** — такого API не существует; используйте `window.Close()`.
- **Не сохраняйте данные при каждом событии** — сначала устраните дребезг.
- **Не предполагайте, что событие содержит координаты** — вызывайте `window.Position()` / `window.Size()`.

## Устранение неполадок

### Перехватчик WindowClosing не предотвращает закрытие

**Причина:** вы использовали `OnWindowEvent` вместо `RegisterHook`.

**Решение:** только функции обратного вызова `RegisterHook` могут вызывать `e.Cancel()`. Функции обратного вызова `OnWindowEvent` только наблюдают за событием и не могут его отменить.

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

### События не срабатывают

**Причина:** обработчик зарегистрирован после возникновения события.

**Решение:** регистрируйте обработчики сразу после создания окна (или в обратном вызове `app.Window.OnCreate`).

```go
window := app.Window.New()
window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) { /* ... */ })
```

### Утечки памяти

**Причина:** обработчики с длительным временем жизни, хранящиеся в состоянии с коротким временем жизни.

**Решение:** сохраните функцию отписки `func()`, возвращаемую `OnWindowEvent`/`RegisterHook`, и вызовите её при очистке.

```go
unsub := window.OnWindowEvent(events.Common.WindowDidResize, handler)
// ...later:
unsub()
```

## Дальнейшие шаги

**Основы работы с окнами** — изучите основы управления окнами [Подробнее →](/features/windows/basics/)

**Несколько окон** — шаблоны для многооконных приложений [Подробнее →](/features/windows/multiple/)

**Система событий** — углублённое изучение системы событий [Подробнее →](/features/events/system/)

**Жизненный цикл приложения** — разберитесь в жизненном цикле приложения [Подробнее →](/concepts/lifecycle/)

---

**Есть вопросы?** Задайте их в [Discord](https://discord.gg/JDdSxwjhGf) или ознакомьтесь с [примерами](https://github.com/wailsapp/wails/tree/master/v3/examples).
