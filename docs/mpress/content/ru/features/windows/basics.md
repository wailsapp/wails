---
title: "Основы работы с окнами"
description: "Создание окон приложения и управление ими в Wails"
slug: "features/windows/basics"
sourcePath: "features/windows/basics.md"
---

## Управление окнами

Wails предоставляет **единый API управления окнами**, который работает на всех платформах. Создавайте окна, управляйте их поведением и работайте с несколькими окнами, полностью контролируя их создание, внешний вид, поведение и жизненный цикл.

## Быстрый старт

```go
package main

import "github.com/wailsapp/wails/v3/pkg/application"

func main() {
    app := application.New(application.Options{
        Name: "My App",
    })

    // Create a window
    window := app.Window.New()
    
    // Configure it
    window.SetTitle("Hello Wails")
    window.SetSize(800, 600)
    window.Center()
    
    // Show it
    window.Show()

    app.Run()
}
```

**Вот и всё!** Теперь у вас есть кроссплатформенное окно.

## Создание окон

### Простое окно

Самый простой способ создать окно:

```go
window := app.Window.New()
```

**Что вы получите:**

- Размер по умолчанию (800x600)
- Заголовок по умолчанию (имя приложения)
- WebView, готовый к работе с вашим фронтендом
- Внешний вид, соответствующий платформе

### Окно с параметрами

Создайте окно с пользовательской конфигурацией:

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "My Application",
    Width:  1200,
    Height: 800,
    X:      100,   // Position from left
    Y:      100,   // Position from top
    AlwaysOnTop: false,
    Frameless: false,
    Hidden: false,
    MinWidth: 400,
    MinHeight: 300,
    MaxWidth: 1920,
    MaxHeight: 1080,
})
```

**Распространённые параметры:**

| Параметр | Тип | Описание |
| --- | --- | --- |
| `Title` | `string` | Заголовок окна |
| `Width` | `int` | Ширина окна в пикселях |
| `Height` | `int` | Высота окна в пикселях |
| `X` | `int` | Координата X (от левого края) |
| `Y` | `int` | Координата Y (от верхнего края) |
| `AlwaysOnTop` | `bool` | Располагать окно поверх остальных |
| `Frameless` | `bool` | Убрать строку заголовка и рамки |
| `Hidden` | `bool` | Изначально скрыть окно |
| `MinWidth` | `int` | Минимальная ширина |
| `MinHeight` | `int` | Минимальная высота |
| `MaxWidth` | `int` | Максимальная ширина |
| `MaxHeight` | `int` | Максимальная высота |

**Полный список см. в разделе «[Параметры окна](/features/windows/options/)».**

### Именованные окна

Присваивайте окнам имена, чтобы их было легко находить:

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:  "main-window",
    Title: "Main Application",
})

// Later, find it by name
if mainWindow, ok := app.Window.GetByName("main-window"); ok {
    mainWindow.Show()
}
```

**Варианты использования:**

- Несколько окон (главное, настройки, сведения о приложении)
- Поиск окон из разных частей кода
- Обмен данными между окнами

## Управление окнами

### Отображение и скрытие

```go
// Show window
window.Show()

// Hide window
window.Hide()

// Check if visible
if window.IsVisible() {
    fmt.Println("Window is visible")
}
```

**Варианты использования:**

- Заставки (показать, затем скрыть)
- Окна настроек (скрывать, когда они не нужны)
- Всплывающие окна (показывать по запросу)

### Положение и размер

```go
// Set size
window.SetSize(1024, 768)

// Set position
window.SetPosition(100, 100)

// Centre on screen
window.Center()

// Get current size
width, height := window.Size()

// Get current position
x, y := window.Position()
```

**Система координат:**

- Точка (0, 0) находится в левом верхнем углу основного экрана
- Положительное направление оси X — вправо
- Положительное направление оси Y — вниз

### Состояние окна

```go
// Minimise
window.Minimise()

// Maximise
window.Maximise()

// Fullscreen
window.Fullscreen()

// Restore to normal
window.Restore()

// Check state
if window.IsMinimised() {
    fmt.Println("Window is minimised")
}

if window.IsMaximised() {
    fmt.Println("Window is maximised")
}

if window.IsFullscreen() {
    fmt.Println("Window is fullscreen")
}
```

**Переходы между состояниями:**

```
Normal ←→ Minimised
Normal ←→ Maximised
Normal ←→ Fullscreen
```

### Заголовок и внешний вид

```go
// Set title
window.SetTitle("My Application - Document.txt")

// Set background colour — RGBA value (helper for RGB)
window.SetBackgroundColour(application.NewRGBA(0, 0, 0, 255))

// Set always on top
window.SetAlwaysOnTop(true)

// Set resizable
window.SetResizable(false)
```

### Закрытие окон

```go
// Close window — dispatches WindowClosing; a RegisterHook can call e.Cancel().
window.Close()
```

В v3 нет метода `window.Destroy()` — используйте `Close()` и либо отслеживайте событие с помощью `OnWindowEvent` (закрытие нельзя отменить), либо перехватывайте его с помощью `RegisterHook` (чтобы оставить окно открытым, можно вызвать `e.Cancel()`).

## Поиск окон

### По имени

```go
if window, ok := app.Window.GetByName("settings"); ok {
    window.Show()
}
```

### По идентификатору

У каждого окна есть уникальный идентификатор:

```go
id := window.ID()
fmt.Printf("Window ID: %d\n", id)

// Find by ID
if found, ok := app.Window.GetByID(id); ok {
    found.Focus()
}
```

### Текущее окно

Получите окно, которое сейчас находится в фокусе:

```go
current := app.Window.Current()
if current != nil {
    current.SetTitle("Active Window")
}
```

### Все окна

Получите все окна:

```go
windows := app.Window.GetAll()
fmt.Printf("Total windows: %d\n", len(windows))

for _, w := range windows {
    fmt.Printf("Window: %s (ID: %d)\n", w.Name(), w.ID())
}
```

## Жизненный цикл окна

### Создание

```go
app.Window.OnCreate(func(window application.Window) {
    fmt.Printf("Window created: %s\n", window.Name())

    // Configure new windows
    window.SetMinSize(400, 300)
})
```

### Закрытие

Чтобы предотвратить закрытие окна, используйте `RegisterHook` с событием `WindowClosing`:

```go
window.RegisterHook(events.Common.WindowClosing, func(event *application.WindowEvent) {
    if hasUnsavedChanges() {
        // Ask user for confirmation
        result := showConfirmDialog("Unsaved changes. Close anyway?")
        if result != "yes" {
            // Cancel the close event
            event.Cancel()
        }
    }
})
```

**Важно:** `RegisterHook` перехватывает событие закрытия до того, как оно произойдёт. Чтобы предотвратить закрытие окна, вызовите `event.Cancel()`. Это работает при закрытии по инициативе пользователя (нажатии кнопки X).

### Уничтожение

Чтобы выполнить очистку при закрытии окна, используйте `OnWindowEvent` с событием `WindowClosing`:

```go
window.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
    fmt.Println("Window is closing")
    // Cleanup resources
})
```

## Несколько окон

### Создание нескольких окон

```go
// Main window
mainWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:   "main",
    Title:  "Main Application",
    Width:  1200,
    Height: 800,
})

// Settings window
settingsWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:   "settings",
    Title:  "Settings",
    Width:  600,
    Height: 400,
    Hidden: true,  // Start hidden
})

// Show settings when needed
settingsWindow.Show()
```

### Взаимодействие между окнами

Окна могут взаимодействовать посредством событий:

```go
// In main window
app.Event.Emit("data-updated", map[string]interface{}{
    "value": 42,
})

// In settings window
app.Event.On("data-updated", func(event *application.CustomEvent) {
    data := event.Data.(map[string]interface{})
    value := data["value"].(int)
    fmt.Printf("Received: %d\n", value)
})
```

**Подробнее см. в разделе [«События»](/features/events/system/).**

### Родительские и дочерние окна

У `WebviewWindowOptions` нет поля `Parent`. Создайте дочернее окно как обычное, а затем прикрепите его к родительскому в качестве модального окна-листа:

```go
// Create child window
childWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Child Window",
})

// Attach to the parent — presents as a sheet on macOS.
mainWindow.AttachModal(childWindow)
```

**Поведение:**

- Дочернее окно остаётся поверх родительского.
- Дочернее окно является модальным — оно блокирует взаимодействие с родительским окном.

**Поддержка платформ:**

- **macOS:** Полная поддержка (отображается как окно-лист).
- **Windows:** Не поддерживается.
- **Linux:** Не поддерживается.

## Возможности отдельных платформ

@tabs{sync-key="platform"}
[Windows]
**Возможности Windows:**

```go
// Flash taskbar button
window.Flash(true)  // Start flashing
window.Flash(false) // Stop flashing

// Trigger Windows 11 Snap Assist (Win+Z)
window.SnapAssist()
```

Отдельного `SetIcon` для каждого окна нет — значок приложения задаётся для всего приложения с помощью `app.SetIcon([]byte)` (а значок окна только для Linux — с помощью поля `application.LinuxWindow.Icon` при создании окна).

**Snap Assist:** Отображает варианты макетов привязки Windows 11 через системное сочетание клавиш. Если для собственной HTML-кнопки развёртывания нужны нативные макеты привязки, появляющиеся при наведении, используйте вместо этого [нативные неклиентские области в Windows](/features/windows/frameless/#native-non-client-regions-on-windows).

**Мигание кнопки на панели задач:** Полезно для уведомлений, когда окно свёрнуто.

[macOS]
**Возможности macOS:**

```go
// Transparent title bar
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Mac: application.MacWindow{
        TitleBar: application.MacTitleBar{
            AppearsTransparent: true,
        },
        Backdrop: application.MacBackdropTranslucent,
    },
})
```

**Типы фона:**

- `MacBackdropNormal` — стандартное окно
- `MacBackdropTranslucent` — полупрозрачный фон; для прозрачности веб-представления **требуется закрытый API**.
- `MacBackdropTransparent` — полностью прозрачный фон; для прозрачности веб-представления **требуется закрытый API**.
- `MacBackdropLiquidGlass` — стеклянный фон; для прозрачности веб-представления **требуется закрытый API**.

Чтобы эти эффекты были видны сквозь веб-представление, выполняйте сборку с `-tags private_mac_apis`. Без этого параметра веб-представление остаётся непрозрачным. Сам `TitleBar.AppearsTransparent` использует общедоступные API. См. раздел [«Закрытые API macOS»](/guides/build/private-macos-apis/).

**Поведение коллекции:** Управляйте поведением окон в разных пространствах Spaces:

- `MacWindowCollectionBehaviorCanJoinAllSpaces` — отображается во всех пространствах Spaces
- `MacWindowCollectionBehaviorFullScreenAuxiliary` — может отображаться поверх полноэкранных приложений

**Нативный полноэкранный режим:** В полноэкранном режиме macOS создаёт новое пространство Space (виртуальный рабочий стол).

[Linux]
**Возможности Linux:**

```go
// Set window icon (per-window struct is LinuxWindow, not the app-level LinuxOptions)
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Linux: application.LinuxWindow{
        Icon: iconBytes,
    },
})
```

**Примечания о средах рабочего стола:**

- GNOME: полная поддержка
- KDE Plasma: полная поддержка
- XFCE: частичная поддержка
- Другие: уровень поддержки различается

**Тайловые оконные менеджеры (Hyprland, Sway, i3 и другие):**

- `Minimise()` и `Maximise()` могут работать не так, как ожидается: геометрией окна управляет оконный менеджер
- Запросы `SetSize()` и `SetPosition()` носят рекомендательный характер и могут быть проигнорированы
- `Fullscreen()` обычно работает как ожидается
- Некоторые оконные менеджеры не поддерживают режим «поверх всех окон»

@end

## Распространённые шаблоны

### Заставка

```go
// Create splash screen
splash := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:     "Loading...",
    Width:     400,
    Height:    300,
    Frameless: true,
    AlwaysOnTop: true,
})

// Show splash
splash.Show()

// Initialise application
time.Sleep(2 * time.Second)

// Hide splash, show main window
splash.Close()
mainWindow.Show()
```

### Окно настроек

```go
var settingsWindow *application.WebviewWindow

func showSettings() {
    if settingsWindow == nil {
        settingsWindow = app.Window.NewWithOptions(application.WebviewWindowOptions{
            Name:   "settings",
            Title:  "Settings",
            Width:  600,
            Height: 400,
        })
    }
    
    settingsWindow.Show()
    settingsWindow.Focus()
}
```

### Подтверждение перед закрытием

```go
window.RegisterHook(events.Common.WindowClosing, func(event *application.WindowEvent) {
    if hasUnsavedChanges() {
        // Show dialog
        result := showConfirmDialog("Unsaved changes. Close anyway?")
        if result != "yes" {
            // Cancel the close event
            event.Cancel()
        }
    }
})
```

## Рекомендации

### ✅ Рекомендуется

- **Присваивайте важным окнам имена** — так их будет проще найти впоследствии
- **Задавайте минимальный размер** — это предотвращает появление непригодных для использования макетов
- **Размещайте окна по центру** — это удобнее для пользователей, чем случайное расположение
- **Обрабатывайте события закрытия** — это предотвращает потерю данных
- **Тестируйте на всех платформах** — поведение различается
- **Используйте подходящие размеры** — учитывайте разные размеры экранов

### ❌ Не рекомендуется

- **Не создавайте слишком много окон** — это запутывает пользователей
- **Не забывайте закрывать окна** — иначе возможны утечки памяти
- **Не задавайте расположение жёстко** — размеры экранов различаются
- **Не игнорируйте различия между платформами** — тщательно тестируйте приложение
- **Не блокируйте поток пользовательского интерфейса** — используйте горутины для длительных операций

## Устранение неполадок

### Окно не отображается

**Возможные причины:**

1. Окно создано скрытым
2. Окно находится за пределами экрана
3. Окно находится позади других окон

**Решение:**

```go
window.Show()
window.Center()
window.Focus()
```

### Неверный размер окна

**Причина:** масштабирование с учётом DPI в Windows/Linux

**Решение:**

```go
// Wails handles DPI automatically
// Just use logical pixels
window.SetSize(800, 600)
```

### Окно сразу закрывается

**Причина:** приложение завершает работу при закрытии последнего окна

**Решение:**

```go
app := application.New(application.Options{
    Mac: application.MacOptions{
        ApplicationShouldTerminateAfterLastWindowClosed: false,
    },
})
```

## Дальнейшие действия

@cards{cols="2"}
⚙ Параметры окна
Полный справочник по всем параметрам окна.

[Подробнее →](/features/windows/options/)

---
▣ Несколько окон
Шаблоны для многооконных приложений.

[Подробнее →](/features/windows/multiple/)

---
★ Окна без рамки
Создавайте собственное оформление окна.

[Подробнее →](/features/windows/frameless/)

---
🚀 События окна
Обрабатывайте события жизненного цикла окна.

[Подробнее →](/features/windows/events/)

@end

---

**Есть вопросы?** Задайте их в [Discord](https://discord.gg/JDdSxwjhGf) или ознакомьтесь с [примерами работы с окнами](https://github.com/wailsapp/wails/tree/master/v3/examples).
