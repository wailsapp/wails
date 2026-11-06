---
title: "Меню в системном трее"
description: "Добавьте в приложение интеграцию с системным треем (областью уведомлений)"
slug: "features/menus/systray"
sourcePath: "features/menus/systray.md"
---

## Меню в системном трее

Wails предоставляет **унифицированные API системного трея**, которые работают на всех платформах. Создавайте значки в трее с меню, привязывайте окна и обрабатывайте щелчки с учётом нативного поведения платформы в фоновых приложениях, службах и утилитах быстрого доступа.

![Меню Wails в системном трее, открытое из строки меню macOS](/assets/screenshots/systray-menu-macos.png)

В macOS элемент Wails в системном трее отображается в строке меню и открывает нативное меню. В примере представлены отключённые элементы, флажки, переключатели, подменю и элементы действий.

## Быстрый старт

```go
package main

import (
    _ "embed"
    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets/icon.png
var icon []byte

func main() {
    app := application.New(application.Options{
        Name: "Tray App",
    })

    // Create system tray
    systray := app.SystemTray.New()
    systray.SetIcon(icon)
    systray.SetLabel("My App")

    // Add menu
    menu := app.NewMenu()
    menu.Add("Show").OnClick(func(ctx *application.Context) {
        // Show main window
    })
    menu.Add("Quit").OnClick(func(ctx *application.Context) {
        app.Quit()
    })
    systray.SetMenu(menu)

    // Create hidden window
    window := app.Window.New()
    window.Hide()

    app.Run()
}
```

**Результат:** значок в системном трее с меню на всех платформах.

## Создание системного трея

### Базовый системный трей

```go
// Create system tray
systray := app.SystemTray.New()

// Set icon
systray.SetIcon(iconBytes)

// Set label (macOS) / tooltip (Windows)
systray.SetLabel("My Application")
```

### Со значком

Значки следует встраивать:

```go
import _ "embed"

//go:embed assets/icon.png
var icon []byte

//go:embed assets/icon-dark.png
var iconDark []byte

func main() {
    app := application.New(application.Options{
        Name: "My App",
    })

    systray := app.SystemTray.New()
    systray.SetIcon(icon)
    systray.SetDarkModeIcon(iconDark)  // Windows and macOS dark mode
    
    app.Run()
}
```

**Требования к значкам:**

| Платформа | Размер | Формат | Примечания |
| --- | --- | --- | --- |
| **Windows** | 16x16 или 32x32 | PNG, ICO | Область уведомлений |
| **macOS** | от 18x18 до 22x22 | PNG | Строка меню; рекомендуется шаблонный значок |
| **Linux** | от 22x22 до 48x48 | PNG, SVG | Зависит от окружения рабочего стола |

### Шаблонные значки (macOS)

Шаблонные значки автоматически адаптируются к светлой и тёмной теме:

```go
systray.SetTemplateIcon(iconBytes)
```

**Рекомендации по шаблонным значкам:**

- Используйте только чёрный и прозрачный цвета
- В тёмной теме чёрный цвет становится белым
- Добавьте к имени файла суффикс `Template`: `iconTemplate.png`
- [Руководство по дизайну](https://bjango.com/articles/designingmenubarextras/)

## Добавление меню

Меню системного трея работают так же, как меню приложения:

```go
menu := app.NewMenu()

// Add items
menu.Add("Open").OnClick(func(ctx *application.Context) {
    showMainWindow()
})

menu.AddSeparator()

menu.AddCheckbox("Start at Login", false).OnClick(func(ctx *application.Context) {
    enabled := ctx.ClickedMenuItem().Checked()
    setStartAtLogin(enabled)
})

menu.AddSeparator()

menu.Add("Quit").OnClick(func(ctx *application.Context) {
    app.Quit()
})

// Set menu
systray.SetMenu(menu)
```

Описание **всех типов элементов меню** см. в [справочнике по меню](/features/menus/reference/).

## Привязка окон

Привяжите окно к значку в трее, чтобы автоматически показывать и скрывать его:

```go
// Create window
window := app.Window.New()

// Attach to tray
systray.AttachWindow(window)

// Configure behaviour — these are setters that return the receiver for chaining.
systray.WindowOffset(10)                          // Pixels from tray icon
systray.WindowDebounce(200 * time.Millisecond)    // Click debounce
```

**Поведение:**

- При запуске окно скрыто
- **Щелчок левой кнопкой по значку в трее** → Переключение видимости окна
- **Щелчок правой кнопкой по значку в трее** → Отображение меню (если оно задано)
- Окно располагается рядом со значком в трее

**Пример: всплывающее окно**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:           "Quick Access",
    Width:           300,
    Height:          400,
    Frameless:       true, // No title bar
    AlwaysOnTop:     true, // Stay on top
    HideOnFocusLost: true, // Dismiss when another window receives focus
    HideOnEscape:    true, // Dismiss when the user presses Escape
})

systray.AttachWindow(window)
systray.WindowOffset(5)
```

`HideOnFocusLost` удобно использовать для всплывающих окон трея в Windows, macOS и средах рабочего стола Linux, где фокус устанавливается щелчком. Wails отключает это поведение в средах Linux, где фокус следует за указателем мыши (включая распространённые конфигурации Hyprland, Sway и i3): иначе всплывающее окно могло бы скрыться при уходе указателя ещё до того, как им удастся воспользоваться. В таких средах `HideOnEscape` остаётся доступным.

Описанное выше поведение при щелчках левой и правой кнопками — это разумные настройки по умолчанию. Явно заданный обработчик `OnClick` или `OnRightClick` заменяет соответствующее поведение по умолчанию. Проверки для отдельных платформ и пограничные случаи см. в [наборе ручных тестов системного трея](https://github.com/wailsapp/wails/tree/v3.0.0/v3/test/manual/systray) и [примере стресс-тестирования системного трея](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/systray-stress).

## Обработчики щелчков

Обрабатывайте щелчки по значку в трее:

```go
systray := app.SystemTray.New()

// Left click
systray.OnClick(func() {
    fmt.Println("Tray icon clicked")
})

// Right click
systray.OnRightClick(func() {
    fmt.Println("Tray icon right-clicked")
})

// Double click
systray.OnDoubleClick(func() {
    fmt.Println("Tray icon double-clicked")
})

// Mouse enter/leave
systray.OnMouseEnter(func() {
    fmt.Println("Mouse entered tray icon")
})

systray.OnMouseLeave(func() {
    fmt.Println("Mouse left tray icon")
})
```

**Поддержка платформ:**

| Событие | Windows | macOS | Linux |
| --- | --- | --- | --- |
| OnClick | ✅ | ✅ | ✅ |
| OnRightClick | ✅ | ✅ | ✅ |
| OnDoubleClick | ✅ | ✅ | ⚠️ Зависит от среды |
| OnMouseEnter | ✅ | ✅ | ⚠️ Зависит от среды |
| OnMouseLeave | ✅ | ✅ | ⚠️ Зависит от среды |

## Динамическое обновление

Динамически обновляйте значок и меню в трее:

### Изменение значка

```go
var isActive bool

func updateTrayIcon() {
    if isActive {
        systray.SetIcon(activeIcon)
        systray.SetLabel("Active")
    } else {
        systray.SetIcon(inactiveIcon)
        systray.SetLabel("Inactive")
    }
}
```

### Обновление меню

```go
var isPaused bool

pauseMenuItem := menu.Add("Pause")

pauseMenuItem.OnClick(func(ctx *application.Context) {
    isPaused = !isPaused
    
    if isPaused {
        pauseMenuItem.SetLabel("Resume")
    } else {
        pauseMenuItem.SetLabel("Pause")
    }
    
    menu.Update()  // Important!
})
```

@note{type="caution" title="Всегда вызывайте Update()"}
После изменения состояния меню **вызовите `menu.Update()`**. См. [справочник по меню](/features/menus/reference/#enabled-state).

@end

### Пересоздание меню

При существенных изменениях пересоздайте всё меню:

```go
func rebuildTrayMenu(status string) {
    menu := app.NewMenu()
    
    // Status-specific items
    switch status {
    case "syncing":
        menu.Add("Syncing...").SetEnabled(false)
        menu.Add("Pause Sync").OnClick(pauseSync)
    case "synced":
        menu.Add("Up to date ✓").SetEnabled(false)
        menu.Add("Sync Now").OnClick(startSync)
    case "error":
        menu.Add("Sync Error").SetEnabled(false)
        menu.Add("Retry").OnClick(retrySync)
    }
    
    menu.AddSeparator()
    menu.Add("Quit").OnClick(func(ctx *application.Context) {
        app.Quit()
    })
    
    systray.SetMenu(menu)
}
```

## Возможности отдельных платформ

@tabs{sync-key="platform"}
[macOS]
**Интеграция со строкой меню:**

```go
// Set label (appears next to icon)
systray.SetLabel("My App")

// Use template icon (adapts to dark mode)
systray.SetTemplateIcon(iconBytes)

// Set icon position — uses AppKit NSImage placement constants.
systray.SetIconPosition(application.NSImageRight)
```

**Расположение значков** (аналогично `NSImagePosition`):

- `application.NSImageLeft` — значок слева от подписи.
- `application.NSImageRight` — значок справа от подписи.
- `application.NSImageOnly` — только значок, без подписи.
- `application.NSImageNone` — только подпись, без значка.

**Рекомендации:**

- Используйте шаблонные значки (чёрный + прозрачный)
- Используйте короткие подписи (3-5 символов)
- От 18x18 до 22x22 пикселей для дисплеев Retina
- Проверяйте как в светлой, так и в тёмной теме

[Windows]
**Интеграция с областью уведомлений:**

```go
// Set tooltip (appears on hover)
systray.SetTooltip("My Application")

// Or use SetLabel (same as tooltip on Windows)
systray.SetLabel("My Application")

// Show/Hide functionality (fully functional)
systray.Show()  // Show tray icon
systray.Hide()  // Hide tray icon
```

**Требования к значкам:**

- 16x16 или 32x32 пикселя
- Формат PNG или ICO
- Прозрачный фон

**Ограничения всплывающей подсказки:**

- Не более 127 символов UTF-16
- Более длинные всплывающие подсказки будут обрезаны
- Для удобства формулируйте кратко

**Возможности платформы:**

- Значок в области уведомлений сохраняется после перезапуска Проводника Windows
- Методы Show() и Hide() полностью функциональны
- Корректное управление жизненным циклом

**Рекомендации:**

- Используйте 32x32 для дисплеев с высокой плотностью пикселей
- Ограничивайте длину всплывающих подсказок значением 127 символов
- Тестируйте в разных версиях Windows
- Учитывайте переполнение области уведомлений
- Используйте Show/Hide для отображения значка в области уведомлений по условию

[Linux]
**Интеграция с областью уведомлений:**

Использует спецификацию StatusNotifierItem (в большинстве современных сред рабочего стола).

```go
systray.SetIcon(iconBytes)
systray.SetLabel("My App")
```

**Поддержка сред рабочего стола:**

- **GNOME**: верхняя панель (с расширением)
- **KDE Plasma**: область уведомлений
- **XFCE**: область уведомлений
- **Другие**: зависит от среды

**Рекомендации:**

- Используйте значки размером 22x22 или 24x24 пикселя
- Значки SVG лучше масштабируются
- Тестируйте в целевых средах рабочего стола
- Предусмотрите резервный вариант для неподдерживаемых сред рабочего стола

@end

## Полный пример

Ниже приведено готовое к эксплуатации приложение со значком в области уведомлений:

```go
package main

import (
    _ "embed"
    "fmt"
    "time"
    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets/icon.png
var icon []byte

//go:embed assets/icon-active.png
var iconActive []byte

type TrayApp struct {
    app     *application.App
    systray *application.SystemTray
    window  *application.WebviewWindow
    menu    *application.Menu
    isActive bool
}

func main() {
    app := application.New(application.Options{
        Name: "Tray Application",
        Mac: application.MacOptions{
            ApplicationShouldTerminateAfterLastWindowClosed: false,
        },
    })

    trayApp := &TrayApp{app: app}
    trayApp.setup()

    app.Run()
}

func (t *TrayApp) setup() {
    // Create system tray
    t.systray = t.app.SystemTray.New()
    t.systray.SetIcon(icon)
    t.systray.SetLabel("Inactive")
    
    // Create menu
    t.createMenu()
    
    // Create window (hidden by default)
    t.window = t.app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:  "Tray Application",
        Width:  400,
        Height: 600,
        Hidden: true,
    })
    
    // Attach window to tray
    t.systray.AttachWindow(t.window)
    t.systray.WindowOffset(10)
    
    // Handle tray clicks
    t.systray.OnRightClick(func() {
        t.systray.OpenMenu()
    })
    
    // Start background task
    go t.backgroundTask()
}

func (t *TrayApp) createMenu() {
    t.menu = t.app.NewMenu()
    
    // Status item (disabled)
    statusItem := t.menu.Add("Status: Inactive")
    statusItem.SetEnabled(false)
    
    t.menu.AddSeparator()
    
    // Toggle active
    t.menu.Add("Start").OnClick(func(ctx *application.Context) {
        t.toggleActive()
    })
    
    // Show window
    t.menu.Add("Show Window").OnClick(func(ctx *application.Context) {
        t.window.Show()
        t.window.Focus()
    })
    
    t.menu.AddSeparator()
    
    // Settings
    t.menu.AddCheckbox("Start at Login", false).OnClick(func(ctx *application.Context) {
        enabled := ctx.ClickedMenuItem().Checked()
        t.setStartAtLogin(enabled)
    })
    
    t.menu.AddSeparator()
    
    // Quit
    t.menu.Add("Quit").OnClick(func(ctx *application.Context) {
        t.app.Quit()
    })
    
    t.systray.SetMenu(t.menu)
}

func (t *TrayApp) toggleActive() {
    t.isActive = !t.isActive
    t.updateTray()
}

func (t *TrayApp) updateTray() {
    if t.isActive {
        t.systray.SetIcon(iconActive)
        t.systray.SetLabel("Active")
    } else {
        t.systray.SetIcon(icon)
        t.systray.SetLabel("Inactive")
    }
    
    // Rebuild menu with new status
    t.createMenu()
}

func (t *TrayApp) backgroundTask() {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        if t.isActive {
            fmt.Println("Background task running...")
            // Do work
        }
    }
}

func (t *TrayApp) setStartAtLogin(enabled bool) {
    // Implementation varies by platform
    fmt.Printf("Start at login: %v\n", enabled)
}
```

## Управление видимостью

Динамическое отображение и скрытие значка в области уведомлений:

```go
// Hide tray icon
systray.Hide()

// Show tray icon
systray.Show()
```

Метода чтения `IsVisible()` нет — при необходимости отслеживайте видимость в состоянии своего приложения.

**Поддержка платформ:**

| Платформа | Hide() | Show() | Примечания |
| --- | --- | --- | --- |
| **Windows** | ✅ | ✅ | Полностью функциональны — значок появляется в области уведомлений и исчезает из неё |
| **macOS** | ✅ | ✅ | Элемент строки меню отображается и скрывается |
| **Linux** | ✅ | ✅ | Зависит от среды рабочего стола |

**Сценарии использования:**

- Временно скрывать значок в области уведомлений в соответствии с настройками пользователя
- Безоконный режим, в котором значок в области уведомлений появляется только при необходимости
- Переключать видимость в зависимости от состояния приложения

**Пример — условная видимость значка в области уведомлений:**

```go
func (t *TrayApp) setTrayVisibility(visible bool) {
    if visible {
        t.systray.Show()
    } else {
        t.systray.Hide()
    }
}

// Show tray only when updates are available
func (t *TrayApp) checkForUpdates() {
    if hasUpdates {
        t.systray.Show()
        t.systray.SetLabel("Update Available")
    } else {
        t.systray.Hide()
    }
}
```

## Очистка ресурсов

Когда значок в области уведомлений больше не нужен, уничтожьте его:

```go
// In OnShutdown
app := application.New(application.Options{
    OnShutdown: func() {
        if systray != nil {
            systray.Destroy()
        }
    },
})
```

**Важно:** при завершении работы всегда уничтожайте значок в области уведомлений, чтобы освободить ресурсы.

## Рекомендации

### ✅ Рекомендуется

- **Используйте шаблонные значки в macOS** — они адаптируются к тёмной теме
- **Используйте короткие подписи** — не более 3-5 символов
- **Добавляйте всплывающие подсказки в Windows** — они помогают пользователям распознать ваше приложение
- **Тестируйте на всех платформах** — поведение различается
- **Обрабатывайте щелчки надлежащим образом** — щелчок левой кнопкой для основного действия, правой — для открытия меню
- **Обновляйте значок в соответствии с состоянием** — визуальная обратная связь важна
- **Уничтожайте значок при завершении работы** — это освобождает ресурсы

### ❌ Не рекомендуется

- **Не используйте крупные значки** — следуйте рекомендациям платформы
- **Не используйте длинные подписи** — они обрезаются
- **Не забывайте о тёмной теме** — тестируйте в тёмной теме Windows и macOS
- **Не блокируйте обработчики щелчков** — они должны выполняться быстро
- **Не забывайте вызывать menu.Update()** — после изменения состояния меню
- **Не рассчитывайте на поддержку области уведомлений** — некоторые среды рабочего стола Linux её не поддерживают

## Устранение неполадок

### Значок не появляется в области уведомлений

**Возможные причины:**

1. Формат значка не поддерживается
2. Размер значка слишком большой или слишком маленький
3. Область уведомлений не поддерживается (Linux)

**Решение:**

Вспомогательной функции `SystemTraySupported()` нет; вместо этого создайте значок в области уведомлений, проверьте платформу и предусмотрите корректную деградацию функциональности:

```go
// Probe support: on Linux without a notification-area extension, the tray
// will simply not appear. Defensive code can fall back to window-only mode
// based on runtime.GOOS or after a short timeout if no tray events arrive.
systray := app.SystemTray.New()
systray.SetIcon(iconBytes)
```

### Значок некорректно выглядит в macOS

**Причина:** не используется шаблонный значок

**Решение:**

```go
// Use template icon
systray.SetTemplateIcon(iconBytes)

// Or design icon as template (black + transparent)
```

### Меню не обновляется

**Причина:** не был вызван `menu.Update()`

**Решение:**

```go
menuItem.SetLabel("New Label")
menu.Update()  // Add this!
```

## Дальнейшие шаги

@cards{cols="2"}
📖 Справочник по меню
Полный справочник по типам и свойствам элементов меню.

[Подробнее →](/features/menus/reference/)

---
☰ Меню приложения
Создавайте строки меню приложения.

[Подробнее →](/features/menus/application/)

---
◆ Контекстные меню
Создавайте контекстные меню, открываемые щелчком правой кнопки.

[Подробнее →](/features/menus/context/)

---
📖 Пример области уведомлений
Изучите полное приложение со значком в области уведомлений.

[Подробнее →](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/systray-basic)

@end

---

**Есть вопросы?** Задайте их в [Discord](https://discord.gg/JDdSxwjhGf) или ознакомьтесь с [примерами работы с областью уведомлений](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/systray-basic).
