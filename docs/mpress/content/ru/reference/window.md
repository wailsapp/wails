---
title: "API окон"
description: "Полный справочник по API окон"
slug: "reference/window"
sourcePath: "reference/window.md"
---

## Обзор

API окон предоставляет методы для управления внешним видом, поведением и жизненным циклом окон. Доступ к ним осуществляется через экземпляры окон или диспетчер `app.Window`.

`Window` — интерфейс, которому соответствует `*application.WebviewWindow`; приведённые ниже сигнатуры методов относятся к `*WebviewWindow`. Многие изменяющие методы возвращают `Window`, что позволяет объединять вызовы в цепочки. Возвращаемые значения указаны в описании каждого метода.

**Основные операции:**

- Создание и отображение окон
- Управление размером, положением и состоянием
- Обработка событий окна
- Управление содержимым окна
- Настройка внешнего вида и поведения

## Видимость

### Show()

Отображает окно. Если окно было скрыто, оно становится видимым. Возвращает получатель, позволяя объединять вызовы в цепочку.

```go
func (w *WebviewWindow) Show() Window
```

**Пример:**

```go
window := app.Window.New()
window.Show()
```

### Hide()

Скрывает окно, не закрывая его. Окно остаётся в памяти, и его можно снова отобразить. Возвращает получатель, позволяя объединять вызовы в цепочку.

```go
func (w *WebviewWindow) Hide() Window
```

**Пример:**

```go
// Hide window temporarily
window.Hide()

// Show it again later
window.Show()
```

**Сценарии использования:**

- Приложения в области уведомлений, которые сворачиваются туда при скрытии
- Пошаговые мастера, в которых окна используются повторно
- Временное скрытие во время выполнения операций

### Close()

Закрывает окно. При этом возникает событие `WindowClosing`.

```go
func (w *WebviewWindow) Close()
```

**Пример:**

```go
window.Close()
```

**Примечание:** если зарегистрированный перехватчик вызовет `event.Cancel()`, закрытие будет предотвращено.

## Свойства окна

### SetTitle()

Задаёт текст в строке заголовка окна. Возвращает получатель, позволяя объединять вызовы в цепочку.

```go
func (w *WebviewWindow) SetTitle(title string) Window
```

**Параметры:**

- `title` — новый заголовок окна

**Пример:**

```go
window.SetTitle("My Application - Document.txt")
```

### Name()

Возвращает уникальный идентификатор имени окна.

```go
func (w *WebviewWindow) Name() string
```

**Пример:**

```go
name := window.Name()
fmt.Println("Window name:", name)

// Retrieve window by name later
if w, ok := app.Window.GetByName(name); ok {
    w.Focus()
}
```

## Размер и положение

### SetSize()

Задаёт размеры окна в пикселях. Возвращает получатель, позволяя объединять вызовы в цепочку.

```go
func (w *WebviewWindow) SetSize(width, height int) Window
```

**Параметры:**

- `width` — ширина окна в пикселях
- `height` — высота окна в пикселях

**Пример:**

```go
window.SetSize(1024, 768)
```

### Size()

Возвращает текущие размеры окна.

```go
func (w *WebviewWindow) Size() (width, height int)
```

**Пример:**

```go
width, height := window.Size()
fmt.Printf("Window is %dx%d\n", width, height)
```

### SetMinSize() / SetMaxSize()

Задаёт минимальные и максимальные размеры окна. Оба метода возвращают получатель, позволяя объединять вызовы в цепочку.

```go
func (w *WebviewWindow) SetMinSize(width, height int) Window
func (w *WebviewWindow) SetMaxSize(width, height int) Window
```

**Пример:**

```go
// Prevent window from being too small
window.SetMinSize(800, 600)

// Prevent window from being too large
window.SetMaxSize(1920, 1080)
```

### SetPosition()

Задаёт положение окна относительно левого верхнего угла экрана.

```go
func (w *WebviewWindow) SetPosition(x, y int)
```

**Параметры:**

- `x` — положение по горизонтали в пикселях
- `y` — положение по вертикали в пикселях

**Пример:**

```go
// Position window at top-left
window.SetPosition(0, 0)

// Position window 100px from top-left
window.SetPosition(100, 100)
```

### Position()

Возвращает текущее положение окна.

```go
func (w *WebviewWindow) Position() (x, y int)
```

**Пример:**

```go
x, y := window.Position()
fmt.Printf("Window is at (%d, %d)\n", x, y)
```

### Center()

Размещает окно по центру экрана.

```go
func (w *WebviewWindow) Center()
```

**Пример:**

```go
window := app.Window.New()
window.Center()
window.Show()
```

**Примечание:** окно центрируется на основном мониторе. Для конфигураций с несколькими мониторами используйте API экранов.

### Focus()

Выводит окно на передний план и передаёт ему фокус клавиатуры.

```go
func (w *WebviewWindow) Focus()
```

**Пример:**

```go
// Bring window to front
window.Focus()
```

## Состояние окна

### Minimise() / UnMinimise()

Сворачивает окно на панель задач или в Dock либо восстанавливает его. `Minimise()` возвращает получатель, позволяя объединять вызовы в цепочку; `UnMinimise()` ничего не возвращает.

```go
func (w *WebviewWindow) Minimise() Window
func (w *WebviewWindow) UnMinimise()
```

**Пример:**

```go
// Minimise window
window.Minimise()

// Restore from minimised state
window.UnMinimise()
```

### Maximise() / UnMaximise()

Разворачивает окно на весь экран либо восстанавливает его прежний размер. `Maximise()` возвращает получатель, позволяя объединять вызовы в цепочку; `UnMaximise()` ничего не возвращает.

```go
func (w *WebviewWindow) Maximise() Window
func (w *WebviewWindow) UnMaximise()
```

**Пример:**

```go
// Maximise window
window.Maximise()

// Restore to previous size
window.UnMaximise()
```

### Fullscreen() / UnFullscreen() / ToggleFullscreen()

Включает или выключает полноэкранный режим. `Fullscreen()` возвращает получатель, позволяя объединять вызовы в цепочку.

```go
func (w *WebviewWindow) Fullscreen() Window
func (w *WebviewWindow) UnFullscreen()
func (w *WebviewWindow) ToggleFullscreen()
```

**Пример:**

```go
// Enter fullscreen
window.Fullscreen()

// Exit fullscreen
window.UnFullscreen()

// Or toggle
window.ToggleFullscreen()
```

Метода `SetFullscreen(bool)` не существует.

### IsMinimised() / IsMaximised() / IsFullscreen()

Проверяет текущее состояние окна.

```go
func (w *WebviewWindow) IsMinimised() bool
func (w *WebviewWindow) IsMaximised() bool
func (w *WebviewWindow) IsFullscreen() bool
```

**Пример:**

```go
if window.IsMinimised() {
    window.UnMinimise()
}

if window.IsMaximised() {
    fmt.Println("Window is maximised")
}

if window.IsFullscreen() {
    window.UnFullscreen()
}
```

## Содержимое окна

### SetURL()

Переходит в окне по указанному URL. Возвращает получатель, позволяя объединять вызовы в цепочку.

```go
func (w *WebviewWindow) SetURL(url string) Window
```

**Параметры:**

- `url` — URL, по которому нужно перейти (для встроенных ресурсов можно использовать `http://wails.localhost/`)

**Пример:**

```go
// Navigate to embedded page
window.SetURL("http://wails.localhost/settings.html")

// Navigate to external URL (if allowed)
window.SetURL("https://wails.io")
```

### SetHTML()

Задаёт содержимое окна непосредственно из строки HTML. Возвращает получатель, позволяя объединять вызовы в цепочку.

```go
func (w *WebviewWindow) SetHTML(html string) Window
```

**Параметры:**

- `html` — отображаемое HTML-содержимое

**Пример:**

```go
html := `
<!DOCTYPE html>
<html>
<head><title>Dynamic Content</title></head>
<body>
    <h1>Hello from Go!</h1>
    <p>This content was generated dynamically.</p>
</body>
</html>
`
window.SetHTML(html)
```

**Сценарии использования:**

- Динамическое создание содержимого
- Простые окна без процесса сборки клиентской части
- Страницы ошибок или заставки

### Reload()

Перезагружает текущее содержимое окна.

```go
func (w *WebviewWindow) Reload()
```

**Пример:**

```go
// Reload current page
window.Reload()
```

**Примечание:** Полезно во время разработки или когда требуется обновить содержимое.

## События окна

Wails предоставляет два метода обработки событий окна:

- **OnWindowEvent()** — отслеживает события окна (не может их предотвратить).
- **RegisterHook()** — перехватывает события окна (их можно предотвратить, вызвав `event.Cancel()`).

### OnWindowEvent()

Регистрирует функцию обратного вызова для событий окна. Возвращает функцию отмены подписки.

```go
func (w *WebviewWindow) OnWindowEvent(
    eventType events.WindowEventType,
    callback func(event *WindowEvent),
) func()
```

**Пример:**

```go
import "github.com/wailsapp/wails/v3/pkg/events"

// Listen for window focus
window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
    app.Logger.Info("Window gained focus")
})

// Listen for window lost focus
window.OnWindowEvent(events.Common.WindowLostFocus, func(e *application.WindowEvent) {
    app.Logger.Info("Window lost focus")
})

// Listen for window resize
window.OnWindowEvent(events.Common.WindowDidResize, func(e *application.WindowEvent) {
    app.Logger.Info("Window resized")
})
```

**Распространённые события окна:**

- `events.Common.WindowClosing` — окно готовится к закрытию
- `events.Common.WindowFocus` — окно получило фокус
- `events.Common.WindowLostFocus` — окно потеряло фокус
- `events.Common.WindowDidMove` — окно перемещено
- `events.Common.WindowDidResize` — размер окна изменён
- `events.Common.WindowMinimise` — окно свёрнуто
- `events.Common.WindowMaximise` — окно развёрнуто
- `events.Common.WindowFullscreen` — окно перешло в полноэкранный режим
- `events.Common.WindowRuntimeReady` — среда выполнения внутри окна инициализирована

### RegisterHook()

Регистрирует перехватчик событий окна. Перехватчики выполняются раньше обработчиков и могут предотвратить событие, вызвав `event.Cancel()`. Возвращает функцию отмены подписки.

```go
func (w *WebviewWindow) RegisterHook(
    eventType events.WindowEventType,
    callback func(event *WindowEvent),
) func()
```

**Пример — предотвращение закрытия окна:**

```go
import "github.com/wailsapp/wails/v3/pkg/events"

window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    confirm := app.Dialog.Question().
        SetTitle("Confirm Close").
        SetMessage("Are you sure you want to close this window?")

    yes := confirm.AddButton("Yes")
    no := confirm.AddButton("No")
    confirm.SetDefaultButton(yes)
    confirm.SetCancelButton(no)

    no.OnClick(func() {
        e.Cancel() // Prevent window from closing
    })

    confirm.Show()
})
```

**Пример — сохранение перед закрытием:**

```go
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    if !hasUnsavedChanges {
        return
    }

    dlg := app.Dialog.Question().
        SetTitle("Unsaved Changes").
        SetMessage("Save changes before closing?")

    save := dlg.AddButton("Save")
    discard := dlg.AddButton("Don't Save")
    cancel := dlg.AddButton("Cancel")
    dlg.SetDefaultButton(save)
    dlg.SetCancelButton(cancel)

    save.OnClick(func() { saveData() })
    cancel.OnClick(func() { e.Cancel() })
    _ = discard // allow close

    dlg.Show()
})
```

### EmitEvent()

Отправляет пользовательское событие во фронтенд окна. Возвращает `true`, если отправка была отменена перехватчиком.

```go
func (w *WebviewWindow) EmitEvent(name string, data ...any) bool
```

**Параметры:**

- `name` — имя события
- `data` — необязательные данные, отправляемые вместе с событием

**Пример:**

```go
// Send data to specific window
window.EmitEvent("data-updated", map[string]any{
    "count":  42,
    "status": "success",
})
```

**Фронтенд (JavaScript):**

```javascript
import { Events } from '@wailsio/runtime'

Events.On('data-updated', (data) => {
    console.log('Count:', data.count)
    console.log('Status:', data.status)
})
```

## Другие методы

### SetEnabled()

Включает или отключает взаимодействие пользователя с окном.

```go
func (w *WebviewWindow) SetEnabled(enabled bool)
```

**Пример:**

```go
// Disable window during long operation
window.SetEnabled(false)

// Perform operation
performLongOperation()

// Re-enable window
window.SetEnabled(true)
```

### SetBackgroundColour()

Задаёт цвет фона окна, отображаемый до загрузки содержимого. Возвращает получатель для объединения вызовов в цепочку.

```go
func (w *WebviewWindow) SetBackgroundColour(colour RGBA) Window
```

`RGBA` имеет тип `application.RGBA{Red, Green, Blue, Alpha uint8}`. Используйте вспомогательные функции `application.NewRGB(r, g, b)` (альфа-канал 255) или `application.NewRGBA(r, g, b, a)`.

**Пример:**

```go
// White background
window.SetBackgroundColour(application.NewRGB(255, 255, 255))

// Dark background with full alpha
window.SetBackgroundColour(application.NewRGBA(30, 30, 30, 255))
```

### SetResizable()

Определяет, может ли пользователь изменять размер окна. Возвращает получатель для объединения вызовов в цепочку.

```go
func (w *WebviewWindow) SetResizable(resizable bool) Window
```

**Пример:**

```go
// Make window fixed size
window.SetResizable(false)
```

### SetAlwaysOnTop()

Определяет, должно ли окно оставаться поверх других окон. Возвращает получатель для объединения вызовов в цепочку.

```go
func (w *WebviewWindow) SetAlwaysOnTop(alwaysOnTop bool) Window
```

**Пример:**

```go
// Keep window on top
window.SetAlwaysOnTop(true)
```

### SetMenu()

Задаёт строку меню окна в Windows и Linux, заменяя отображаемую строку меню.

```go
func (w *WebviewWindow) SetMenu(menu *Menu)
```

**Параметры:**
- `menu` — отображаемое меню или `nil` для удаления строки меню

При передаче `nil` строка меню удаляется: в окне не отображается никакое меню, включая меню приложения, а сочетания клавиш удалённого меню перестают работать в этом окне. Если вызвать метод до запуска окна, окно будет создано без строки меню.

**Поддержка платформ:**
- **Windows**: полная поддержка
- **Linux**: полная поддержка
- **macOS**: не действует, поскольку строка меню принадлежит приложению. Вместо этого используйте `app.Menu.Set`.

**Пример:**

```go
window.SetMenu(menu) // show a menu bar
window.SetMenu(nil)  // remove it
```

### Print()

Открывает системное диалоговое окно печати для содержимого окна.

```go
func (w *WebviewWindow) Print() error
```

**Возвращаемое значение:** ошибка при сбое печати.

**Пример:**

```go
if err := window.Print(); err != nil {
    log.Println("Print failed:", err)
}
```

### AttachModal()

Присоединяет второе окно как модальное окно-лист.

```go
func (w *WebviewWindow) AttachModal(modalWindow Window)
```

**Параметры:**

- `modalWindow` — окно, присоединяемое как модальное

**Поддержка платформ:**

- **macOS**: полная поддержка (отображается как лист)
- **Windows**: не поддерживается
- **Linux**: не поддерживается

**Пример:**

```go
modalWindow := app.Window.New()
window.AttachModal(modalWindow)
```

## Платформозависимые параметры

### Linux

Окна в Linux поддерживают следующие платформозависимые параметры через `LinuxWindow`:

#### MenuStyle

Определяет способ отображения меню приложения. Этот параметр доступен в стандартной сборке GTK4 и игнорируется в устаревших сборках `-tags gtk3`.

| Значение | Описание |
| --- | --- |
| `LinuxMenuStyleMenuBar` | Традиционная строка меню под строкой заголовка (по умолчанию) |
| `LinuxMenuStylePrimaryMenu` | Кнопка главного меню в строке заголовка (в стиле GNOME) |

**Пример:**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "My Application",
    Linux: application.LinuxWindow{
        MenuStyle: application.LinuxMenuStylePrimaryMenu,
    },
})
window.SetMenu(menu)
```

**Примечание:** При использовании стиля главного меню в строке заголовка отображается кнопка-гамбургер (☰) в соответствии с рекомендациями GNOME по проектированию пользовательских интерфейсов. Этот стиль рекомендуется для современных приложений GNOME.

## Полный пример

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)

func main() {
    app := application.New(application.Options{
        Name: "Window API Demo",
    })

    // Create window with options
    window := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:            "My Application",
        Width:            1024,
        Height:           768,
        MinWidth:         800,
        MinHeight:        600,
        BackgroundColour: application.NewRGB(255, 255, 255),
        URL:              "http://wails.localhost/",
    })

    // Configure window behaviour
    window.SetResizable(true)
    window.SetMinSize(800, 600)
    window.SetMaxSize(1920, 1080)

    // Confirm-before-close hook
    window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
        dlg := app.Dialog.Question().
            SetTitle("Confirm Close").
            SetMessage("Are you sure you want to close this window?")

        yes := dlg.AddButton("Yes")
        no := dlg.AddButton("No")
        dlg.SetDefaultButton(yes)
        dlg.SetCancelButton(no)
        no.OnClick(func() { e.Cancel() })

        dlg.Show()
    })

    // Listen for window events
    window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
        window.SetTitle("My Application (Active)")
        app.Logger.Info("Window gained focus")
    })

    window.OnWindowEvent(events.Common.WindowLostFocus, func(e *application.WindowEvent) {
        window.SetTitle("My Application")
        app.Logger.Info("Window lost focus")
    })

    // Position and show window
    window.Center()
    window.Show()

    app.Run()
}
```
