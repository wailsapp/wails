---
title: "Параметры окна"
description: "Параметры окна, значения по умолчанию и ограничения платформ"
slug: "features/windows/options"
sourcePath: "features/windows/options.md"
---

## Параметры конфигурации окна

Wails предоставляет широкие возможности настройки размера, положения, внешнего вида и поведения окна с помощью десятков параметров. В этом справочном разделе описаны часто используемые поля `WebviewWindowOptions` в Windows, macOS и Linux. Приведённая ниже структура содержит только некоторые поля; полное определение типа см. в исходном коде Go по ссылке.

## Структура WebviewWindowOptions

```go
type WebviewWindowOptions struct {
    // Identity
    Name  string
    Title string

    // Size and Position
    Width           int
    Height          int
    X               int
    Y               int
    MinWidth        int
    MinHeight       int
    MaxWidth        int
    MaxHeight       int
    InitialPosition WindowStartPosition // WindowCentered (default) or WindowXY
    Screen          *Screen             // target screen for initial placement

    // Initial State
    Hidden        bool
    Frameless     bool
    DisableResize bool        // inverted vs v2's `Resizable`
    AlwaysOnTop   bool
    StartState    WindowState // WindowStateNormal | Minimised | Maximised | Fullscreen

    // Appearance
    BackgroundColour RGBA
    BackgroundType   BackgroundType
    Zoom             float64
    ZoomControlEnabled bool

    // Content
    URL  string
    HTML string
    JS   string
    CSS  string

    // Behaviour
    EnableFileDrop              bool
    IgnoreMouseEvents           bool
    HideOnFocusLost             bool
    HideOnEscape                bool
    DevToolsEnabled             bool
    DefaultContextMenuDisabled  bool
    ContentProtectionEnabled    bool
    KeyBindings                 map[string]func(window Window)

    // Permissions
    Permissions map[PermissionType]Permission

    // Window-control button states
    MinimiseButtonState ButtonState
    MaximiseButtonState ButtonState
    CloseButtonState    ButtonState

    // Menu
    UseApplicationMenu bool

    // Platform-specific (per-window)
    Mac     MacWindow
    Windows WindowsWindow
    Linux   LinuxWindow
}
```

В `WebviewWindowOptions` нет поля `Parent` — для родительских и модальных отношений используйте `parentWindow.AttachModal(childWindow)`. В этой структуре также нет поля `Assets` — ресурсы настраиваются в `application.Options` (`Assets AssetOptions`).

Полный исходный код: [`v3/pkg/application/webview_window_options.go`](https://github.com/wailsapp/wails/blob/releases/v3.0.0/v3/pkg/application/webview_window_options.go).

## Основные параметры

### Name

**Тип:** `string` **По умолчанию:** `window-<numeric ID>` **Платформы:** все

```go
Name: "main-window"
```

**Назначение:** уникальный идентификатор, позволяющий впоследствии найти окно.

**Рекомендации:**

- Используйте понятные имена: `"main"`, `"settings"`, `"about"`
- Используйте kebab-case: `"file-browser"`, `"color-picker"`
- Выбирайте короткое и запоминающееся имя

**Пример:**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name: "settings-window",
})

// Later...
if settings, ok := app.Window.GetByName("settings-window"); ok {
    settings.Focus()
}
```

### Title

**Тип:** `string` **По умолчанию:** имя приложения **Платформы:** все

```go
Title: "My Application"
```

**Назначение:** текст, отображаемый в строке заголовка и на панели задач.

**Динамическое обновление:**

```go
window.SetTitle("My Application - Document.txt")
```

### Width / Height

**Тип:** `int` (пиксели) **По умолчанию:** 800 × 600 **Платформы:** все **Ограничения:** значения должны быть положительными

```go
Width:  1200,
Height: 800,
```

**Назначение:** начальный размер окна в логических пикселях.

**Примечания:**

- Wails автоматически учитывает масштабирование DPI
- Используйте логические, а не физические пиксели
- Учитывайте минимальное разрешение экрана (1024 × 768)

**Примеры размеров:**

| Сценарий использования | Ширина | Высота |
| --- | --- | --- |
| Небольшая утилита | 400 | 300 |
| Стандартное приложение | 1024 | 768 |
| Крупное приложение | 1440 | 900 |
| Full HD | 1920 | 1080 |

### X / Y

**Тип:** `int` (пиксели) **По умолчанию:** по центру экрана **Платформы:** все

```go
X: 100,  // 100px from left edge
Y: 100,  // 100px from top edge
```

**Назначение:** начальное положение окна.

**Система координат:**

- Точка (0, 0) находится в левом верхнем углу основного экрана
- Положительное направление оси X — вправо
- Положительное направление оси Y — вниз

**Пример:**

`X` и `Y` действуют только при заданном `InitialPosition: application.WindowXY`. В противном случае `InitialPosition` по умолчанию получает значение `WindowCentered`, а `X`/`Y` игнорируются.

```go
settings := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:            "coordinate-window",
    InitialPosition: application.WindowXY, // opt into X/Y coordinates
    X:               100,
    Y:               100,
})
```

**Рекомендация:** если точные координаты не важны, после создания окна центрируйте его с помощью `Center()`:

```go
window := app.Window.New()
window.Center()
```

### MinWidth / MinHeight

**Тип:** `int` (пиксели) **По умолчанию:** 0 (минимум не задан) **Платформы:** все

```go
MinWidth:  400,
MinHeight: 300,
```

**Назначение:** не позволяет сделать окно слишком маленьким.

**Сценарии использования:**

- Предотвращение нарушения макета
- Обеспечение удобства использования
- Сохранение соотношения сторон

**Пример:**

```go
// Prevent window smaller than 400x300
MinWidth:  400,
MinHeight: 300,
```

### MaxWidth / MaxHeight

**Тип:** `int` (пиксели) **По умолчанию:** 0 (максимум не задан) **Платформы:** все

```go
MaxWidth:  1920,
MaxHeight: 1080,
```

**Назначение:** не позволяет сделать окно слишком большим.

**Сценарии использования:**

- Приложения с фиксированным размером окна
- Предотвращение чрезмерного потребления ресурсов
- Соблюдение ограничений дизайна

## Параметры состояния

### Hidden

**Тип:** `bool` **По умолчанию:** `false` **Платформы:** все

```go
Hidden: true,
```

**Назначение:** создать окно, не отображая его.

**Сценарии использования:**

- Фоновые окна
- Окна, отображаемые по требованию
- Заставки (создать, загрузить, затем показать)
- Предотвращение белой вспышки при загрузке содержимого

**Улучшения для платформ:**

- **Windows:** устранена белая вспышка окна — окно остаётся невидимым до вызова `Show()`
- **macOS:** полная поддержка
- **Linux:** полная поддержка

**Рекомендуемый шаблон для плавной загрузки:**

```go
// Create hidden window
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:             "main-window",
    Hidden:           true,
    BackgroundColour: application.NewRGB(30, 30, 30), // Match your theme
})

// Load content while hidden
// ... content loads ...

// Show when ready (no flash!)
window.Show()
```

**Пример:**

```go
settings := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:   "settings",
    Hidden: true,
})

// Show when needed
settings.Show()
```

### Frameless

**Тип:** `bool` **По умолчанию:** `false` **Платформы:** все

```go
Frameless: true,
```

**Назначение:** убрать строку заголовка и границы окна.

**Сценарии использования:**

- Пользовательское оформление окна
- Экраны-заставки
- Киоск-приложения
- Окна с индивидуальным дизайном

**Важно:** вам потребуется реализовать:

- Перетаскивание окна
- Кнопки закрытия, сворачивания и разворачивания
- Области изменения размера (если размер окна можно изменять)

**Подробности см. в разделе [«Окна без рамки»](/features/windows/frameless/).**

### DisableResize

**Тип:** `bool` **По умолчанию:** `false` (по умолчанию размер окна можно изменять) **Платформы:** все

```go
DisableResize: true,
```

**Назначение:** запретить изменение размера окна. Обратите внимание, что это поле имеет **обратный** смысл по отношению к `Resizable` в v2: чтобы запретить изменение размера окна, задайте `DisableResize: true`.

**Варианты использования:**

- Приложения с фиксированным размером окна
- Экраны-заставки
- Диалоговые окна

**Примечание:** пользователи всё равно смогут развернуть окно или перейти в полноэкранный режим, если эти возможности также не отключены через `MaximiseButtonState` или поведение коллекции.

### AlwaysOnTop

**Тип:** `bool` **По умолчанию:** `false` **Платформы:** все

```go
AlwaysOnTop: true,
```

**Назначение:** удерживать окно поверх всех остальных окон.

**Варианты использования:**

- Плавающие панели инструментов
- Уведомления
- Режим «картинка в картинке»
- Таймеры

**Примечания по платформам:**

- **macOS:** полная поддержка
- **Windows:** полная поддержка
- **Linux:** зависит от оконного менеджера

### StartState

**Тип:** перечисление `WindowState` **По умолчанию:** `WindowStateNormal` **Платформы:** все

```go
StartState: application.WindowStateMaximised,
```

**Назначение:** начальное состояние окна при его отображении.

**Значения:**

- `WindowStateNormal` — обычное окно
- `WindowStateMinimised` — свёрнутое окно
- `WindowStateMaximised` — развёрнутое окно
- `WindowStateFullscreen` — полноэкранный режим

Константы `WindowStateHidden` не существует. Чтобы окно при запуске было невидимым, используйте логическое поле `Hidden`.

**Переключение полноэкранного режима во время выполнения:**

```go
window.Fullscreen()
window.UnFullscreen()
window.ToggleFullscreen() // there is no SetFullscreen(bool)
```

## Параметры внешнего вида

### BackgroundColour

**Тип:** структура `RGBA` **По умолчанию:** белый **Платформы:** все

```go
BackgroundColour: application.RGBA{Red: 0, Green: 0, Blue: 0, Alpha: 255},
```

Поля `RGBA` имеют тип `Red, Green, Blue, Alpha` (uint8). Рекомендуется использовать вспомогательные функции `application.NewRGB(r, g, b)` (альфа-канал 255) или `application.NewRGBA(r, g, b, a)`.

**Назначение:** цвет фона окна до загрузки содержимого.

**Варианты использования:**

- Соответствие теме приложения
- Предотвращение белой вспышки при использовании тёмной темы
- Плавное отображение при загрузке

**Пример:**

```go
// Dark theme
BackgroundColour: application.NewRGB(30, 30, 30),

// Light theme
BackgroundColour: application.NewRGB(255, 255, 255),
```

**Вспомогательный метод:**

```go
window.SetBackgroundColour(application.NewRGB(30, 30, 30))
```

### BackgroundType

**Тип:** перечисление `BackgroundType` **По умолчанию:** `BackgroundTypeSolid` **Платформы:** macOS, Windows (частично)

```go
BackgroundType: application.BackgroundTypeTranslucent,
```

**Значения:**

- `BackgroundTypeSolid` — сплошной цвет
- `BackgroundTypeTransparent` — полная прозрачность
- `BackgroundTypeTranslucent` — полупрозрачное размытие

**Поддержка платформ:**

- **macOS:** настройте `Mac.Backdrop`; для прозрачности webview требуется [`private_mac_apis`](/guides/build/private-macos-apis/#webview-transparency-and-background). Без этого webview останется непрозрачным.
- **Windows:** прозрачный и полупрозрачный фон (Windows 11+)
- **Linux:** только сплошной фон

**Пример (macOS):**

```go
BackgroundType: application.BackgroundTypeTranslucent,
Mac: application.MacWindow{
    Backdrop: application.MacBackdropTranslucent,
},
```

### OpenInspectorOnStartup и OpenDevTools

**Закрытый API в macOS:** для программного открытия инспектора через `OpenInspectorOnStartup: true`, Go `window.OpenDevTools()` и JavaScript `Window.OpenDevTools()` требуется `private_mac_apis`. Без него эти операции ничего не делают. Для производственных сборок также требуется `devtools`. Для общедоступной инспекции через Safari в macOS 13.3+ закрытые API не требуются; в более ранних версиях macOS они необходимы для включения инспектора. См. [матрицу сборок Web Inspector](/guides/build/private-macos-apis/#web-inspector).

## Параметры содержимого

### URL

**Тип:** `string` **По умолчанию:** пустая строка (загрузка из Assets) **Платформы:** все

```go
URL: "https://example.com",
```

**Назначение:** загружать внешний URL вместо встроенных ресурсов.

**Варианты использования:**

- Разработка (загрузка с сервера разработки)
- Веб-приложения
- Гибридные приложения

**Пример:**

```go
// Development — point the window at the Vite dev server
URL: "http://localhost:9245",

// Production — embedded assets are configured at the application level
// (Assets is application.Options.Assets, not a WebviewWindowOptions field).
```

Фрагмент кода уровня приложения для производственной среды:

```go
app := application.New(application.Options{
    Name: "My App",
    Assets: application.AssetOptions{
        Handler: application.AssetFileServerFS(assets),
    },
})
```

### HTML

**Тип:** `string` **По умолчанию:** пустая строка **Платформы:** все

```go
HTML: "<h1>Hello World</h1>",
```

**Назначение:** напрямую загружать строку HTML.

**Варианты использования:**

- Простые окна
- Сгенерированное содержимое
- Тестирование

**Пример:**

```go
HTML: `
<!DOCTYPE html>
<html>
<head><title>Simple Window</title></head>
<body><h1>Hello from Wails!</h1></body>
</html>
`,
```

### Ресурсы (только на уровне приложения)

Конфигурация ресурсов **не** является полем `WebviewWindowOptions`. Ресурсы фронтенда предоставляет само приложение через `application.Options.Assets` (`AssetOptions`); каждое окно использует этот сервер ресурсов.

```go
app := application.New(application.Options{
    Assets: application.AssetOptions{
        Handler: application.AssetFileServerFS(assets),
    },
})
```

**Подробнее см. в разделе [«Система сборки»](/concepts/build-system/).**

### UseApplicationMenu

**Тип:** `bool` **Значение по умолчанию:** `false` **Платформы:** Windows, Linux (не действует в macOS)

```go
UseApplicationMenu: true,
```

**Назначение:** использовать для этого окна меню приложения, заданное через `app.Menu.Set()`.

В **macOS** этот параметр не действует, поскольку macOS всегда использует глобальное меню приложения в верхней части экрана.

В **Windows** и **Linux** окна по умолчанию не отображают меню. Параметр `UseApplicationMenu: true` указывает окну использовать меню уровня приложения, обеспечивая простое кроссплатформенное решение.

**Пример:**

```go
// Set the application menu once
menu := app.NewMenu()
menu.AddRole(application.FileMenu)
menu.AddRole(application.EditMenu)
app.Menu.Set(menu)

// All windows with UseApplicationMenu will display this menu
app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:              "Main Window",
    UseApplicationMenu: true,
})

app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:              "Second Window",
    UseApplicationMenu: true,  // Also gets the app menu
})
```

**Примечания:**

- Если заданы и `UseApplicationMenu`, и меню конкретного окна, приоритет имеет меню окна
- Это упрощает кроссплатформенный код, устраняя необходимость проверять ОС во время выполнения
- Полную документацию по меню см. в разделе [«Меню приложения»](/features/menus/application/)

## Параметры ввода

### EnableFileDrop

**Тип:** `bool` **Значение по умолчанию:** `false` **Платформы:** все

```go
EnableFileDrop: true,
```

**Назначение:** разрешить перетаскивание файлов из операционной системы в окно.

Если параметр включён:

- Файлы можно перетаскивать из файловых менеджеров в приложение
- Событие `WindowFilesDropped` срабатывает и передаёт пути к перемещённым файлам
- Элементы с атрибутом `data-file-drop-target` предоставляют подробные сведения о перетаскивании

**Варианты использования:**

- Интерфейсы загрузки файлов
- Редакторы документов
- Средства импорта медиафайлов
- Любые приложения, принимающие файлы

**Пример:**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:          "File Uploader",
    EnableFileDrop: true,
})

// Handle dropped files
window.OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {
    files := event.Context().DroppedFiles()
    details := event.Context().DropTargetDetails()
    
    for _, file := range files {
        fmt.Println("Dropped:", file)
    }
})
```

**HTML-зоны перетаскивания:**

```html
<!-- Mark elements as drop targets -->
<div id="upload" data-file-drop-target>
    Drop files here
</div>
```

**Полную документацию см. в разделе [«Перетаскивание файлов»](/features/drag-and-drop/files/).**

## Параметры безопасности

### ContentProtectionEnabled

**Тип:** `bool` **Значение по умолчанию:** `false` **Платформы:** Windows (10+), macOS

```go
ContentProtectionEnabled: true,
```

**Назначение:** запретить захват содержимого окна с экрана.

**Поддержка платформ:**

- **Windows:** Windows 10, сборка 19041 и новее (полная поддержка); более ранние версии (частичная поддержка)
- **macOS:** полная поддержка
- **Linux:** не поддерживается

**Варианты использования:**

- Банковские приложения
- Менеджеры паролей
- Медицинские карты
- Конфиденциальные документы

**Важные примечания:**

1. Не предотвращает фотографирование экрана физической камерой
2. Некоторые инструменты могут обходить защиту
3. Это часть комплексной системы безопасности, а не единственная мера защиты
4. Окна DevTools не защищаются автоматически

**Пример:**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Secure Window",
    ContentProtectionEnabled: true,
})

// Toggle at runtime
window.SetContentProtection(true)
```

### Разрешения

**Тип:** `map[PermissionType]Permission` **Значение по умолчанию:** `nil` (обработка по умолчанию для платформы) **Платформы:** Linux, Windows (в macOS обработка передаётся TCC)

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionAllow,
    application.PermissionCamera:     application.PermissionAllow,
},
```

**Назначение:** декларативно управлять обработкой запросов возможностей (камеры, микрофона, геолокации, уведомлений и чтения буфера обмена) из веб-содержимого окна без платформозависимого кода.

**Значения PermissionType:** `PermissionMicrophone`, `PermissionCamera`, `PermissionGeolocation`, `PermissionNotifications`, `PermissionClipboardRead`

**Значения разрешений:**

- `PermissionDefault` (0) — встроенная обработка платформы: запрос ОС/WebView2 в macOS/Windows; в Linux доступ к камере и микрофону разрешается, а всё остальное запрещается
- `PermissionAllow` (1) — предоставить доступ без запроса (в Linux реализовано только для камеры и микрофона; доступ остальных типов остаётся запрещённым)
- `PermissionDeny` (2) — запретить доступ без запроса

**Важно для Windows:** до появления этого параметра Wails без уведомления предоставлял доступ ко всем возможностям WebView2. Теперь добавление любой записи в `Permissions` отключает это общее разрешение. Для возможностей, которых нет в списке, WebView2 будет показывать встроенный запрос вместо автоматического предоставления доступа. Явно перечислите все возможности, необходимые приложению.

**Полное руководство, таблицу поддержки платформ и примеры см. в разделе [«Разрешения»](/features/windows/permissions/).**

## События жизненного цикла окна

События жизненного цикла окна обрабатываются с помощью `OnWindowEvent` и `RegisterHook`. Эти методы позволяют точно управлять поведением при закрытии и уничтожении окна.

### Отмена закрытия окна

Чтобы предотвратить закрытие окна, например при наличии несохранённых изменений, используйте `RegisterHook` с событием `WindowClosing` и вызовите `event.Cancel()`:

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:  "main-window",
    Title: "My Application",
})

// Register a hook to intercept the closing event
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

**Ключевые моменты:**

- `RegisterHook` перехватывает события до того, как они произойдут
- Чтобы предотвратить закрытие окна, вызовите `event.Cancel()`
- После отмены окно останется открытым

### Обработка закрытия окна

Чтобы выполнить очистку при закрытии окна, используйте `OnWindowEvent` с событием `WindowClosing`:

```go
window.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
    // Cleanup code runs here
    fmt.Printf("Window %s is closing\n", window.Name())

    // Close database connection
    if db != nil {
        db.Close()
    }

    // Remove from window list
    removeWindow(window.ID())
})
```

**Ключевые моменты:**

- `OnWindowEvent` обрабатывает события, которые должны произойти в ближайший момент
- Очистка выполняется до уничтожения окна
- Отменить закрытие отсюда нельзя (для этого используйте `RegisterHook`)

### Шаблон очистки для окна-одиночки

Для окон-одиночек (существующих только в одном экземпляре) используйте `WindowClosing`, чтобы очистить ссылку:

```go
var settingsWindow *application.WebviewWindow

func ShowSettings(app *application.App) {
    // Create if doesn't exist
    if settingsWindow == nil {
        settingsWindow = app.Window.NewWithOptions(application.WebviewWindowOptions{
            Name:  "settings",
            Title: "Settings",
            Width: 600,
            Height: 400,
        })

        // Cleanup on close
        settingsWindow.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
            settingsWindow = nil
        })
    }

    // Show and focus
    settingsWindow.Show()
    settingsWindow.Focus()
}
```

## Платформозависимые параметры

### Параметры macOS

```go
Mac: application.MacWindow{
    TitleBar: application.MacTitleBar{
        AppearsTransparent: true,
        Hide:               false,
        HideTitle:          true,
        FullSizeContent:    true,
    },
    Backdrop:                application.MacBackdropTranslucent,
    InvisibleTitleBarHeight: 50,
    WindowClass:             application.MacWindowClassPanel,
    PanelPreferences: application.MacPanelPreferences{
        NonActivating:          true,
        FloatingPanel:          true,
        BecomesKeyOnlyIfNeeded: false,
        UtilityWindow:          false,
    },
    WindowLevel:             application.MacWindowLevelFloating,
    CollectionBehavior:      application.MacWindowCollectionBehaviorDefault,
    TabbingMode:             application.MacWindowTabbingModeDisallowed,
},
```

**TitleBar** (`MacTitleBar`)

- `AppearsTransparent` — делает строку заголовка прозрачной, а содержимое распространяется на её область
- `Hide` — полностью скрывает строку заголовка
- `HideTitle` — скрывает только текст заголовка
- `FullSizeContent` — распространяет содержимое на всю область окна

В macOS для прозрачности webview и программного открытия инспектора требуется тег сборки `private_mac_apis`. Без него те же параметры остаются допустимыми, но операции, доступные только через закрытые API, ничего не делают. Группировка Liquid Glass также игнорируется, а для стилей используются общедоступные альтернативы. Команды сборки и точное описание поведения см. в разделе [Закрытые API macOS](/guides/build/private-macos-apis/).

**Backdrop** (`MacBackdrop`)

- `MacBackdropNormal` — стандартный непрозрачный фон
- `MacBackdropTranslucent` — **Для прозрачности webview требуется закрытый API.** Без этого тега нативное размытие остаётся позади непрозрачного webview.
- `MacBackdropTransparent` — **Для прозрачности webview требуется закрытый API.** Без этого тега webview остаётся непрозрачным.
- `MacBackdropLiquidGlass` — **Для прозрачности webview требуется закрытый API.** Без этого тега стеклянный слой остаётся позади непрозрачного webview, а для стилей используются общедоступные альтернативы.

**LiquidGlass** (`MacLiquidGlass`)

| Поле или значение | Зависимость от закрытого API в macOS |
| --- | --- |
| `Style: LiquidGlassStyleAutomatic` | Нативный обычный стиль является общедоступным; для прозрачности фонового webview требуется `private_mac_apis`. |
| `Style: LiquidGlassStyleLight` | Тег сохраняет существующее нативное сопоставление с прозрачным стилем; без него Wails использует обычное стекло с оформлением Aqua. |
| `Style: LiquidGlassStyleDark` | **Закрытый API:** недокументированное значение нативного стиля `2`; без тега используется обычное стекло с оформлением Dark Aqua. |
| `Style: LiquidGlassStyleVibrant` | Нативное сопоставление с прозрачным стилем является общедоступным; для прозрачности фонового webview требуется тег. |
| `GroupID` | **Закрытый API:** непустые значения запрашивают группировку; без тега они игнорируются. |
| `GroupSpacing` | **Закрытый API:** положительные значения задают интервал между группами; без тега они игнорируются. |
| `Material`, `CornerRadius`, `TintColor` | Сами по себе от закрытых API не зависят. |

Нативные значения стилей и их доступность в разных версиях ОС см. в разделе [Значения Liquid Glass](/guides/build/private-macos-apis/#liquid-glass-values).

**InvisibleTitleBarHeight** (`int`)

- Высота невидимой области строки заголовка (для перетаскивания)
- Действует только тогда, когда нативная область перетаскивания строки заголовка скрыта, то есть когда окно не имеет рамки (`Frameless: true`) или использует прозрачную строку заголовка (`AppearsTransparent: true`)
- Не влияет на стандартные окна с видимой строкой заголовка

**WindowClass** (`MacWindowClass`)

- `MacWindowClassWindow` — стандартное поведение `NSWindow` (по умолчанию)
- `MacWindowClassPanel` — вспомогательное окно `NSPanel`, которое никогда не становится главным окном приложения

`PanelPreferences` применяется только к `MacWindowClassPanel`:

- `NonActivating` добавляет `NSWindowStyleMaskNonactivatingPanel`. Отображение панели или передача ей фокуса не активирует приложение Wails, но панель всё равно может стать ключевой для элементов управления и ввода текста.
- `FloatingPanel` включает поведение плавающей панели AppKit.
- `BecomesKeyOnlyIfNeeded` получает статус ключевого окна, только если представление, по которому щёлкнули, запрашивает ввод с клавиатуры.
- `UtilityWindow` применяет нативный стиль служебного окна.

Панели Wails остаются видимыми после деактивации приложения и освобождаются при закрытии в соответствии с жизненным циклом, ожидаемым для `WebviewWindow`. Это намеренное переопределение противоположных значений `NSPanel` по умолчанию.

Класс и уровень окна, политика активации и поведение коллекции решают разные задачи:

- `WindowClass` выбирает `NSWindow` или `NSPanel` и управляет семантикой главного и ключевого окна.
- `WindowLevel` управляет порядком по оси Z. Для наложений строки меню используйте `MacWindowLevelPopUpMenu`.
- `MacOptions.ActivationPolicy` управляет приложением в целом, включая его представление в Dock и строке меню. Для неактивирующей панели не требуется политика активации вспомогательного приложения, хотя приложение всё же может использовать её, чтобы скрыть свой значок в Dock.
- `CollectionBehavior` управляет участием в Spaces и полноэкранном режиме.

```go
// Spotlight/menu-bar panel that leaves the current application active.
Mac: application.MacWindow{
    WindowClass: application.MacWindowClassPanel,
    PanelPreferences: application.MacPanelPreferences{
        NonActivating: true,
    },
    WindowLevel: application.MacWindowLevelPopUpMenu,
    CollectionBehavior: application.MacWindowCollectionBehaviorCanJoinAllSpaces |
        application.MacWindowCollectionBehaviorFullScreenAuxiliary |
        application.MacWindowCollectionBehaviorStationary,
},
```

**WindowLevel** (`MacWindowLevel`)

- `MacWindowLevelNormal` — стандартный уровень окна (по умолчанию)
- `MacWindowLevelFloating` — располагается поверх обычных окон
- `MacWindowLevelTornOffMenu` — уровень отделённого меню
- `MacWindowLevelModalPanel` — уровень модальной панели
- `MacWindowLevelMainMenu` — уровень главного меню
- `MacWindowLevelStatus` — уровень окна состояния
- `MacWindowLevelPopUpMenu` — уровень всплывающего меню
- `MacWindowLevelScreenSaver` — уровень экранной заставки

Явно заданное значение `WindowLevel` имеет приоритет над `AlwaysOnTop` и `PanelPreferences.FloatingPanel`. Если уровень не задан явно, `AlwaysOnTop` или плавающая панель разрешается в `MacWindowLevelFloating`; в противном случае используется уровень `MacWindowLevelNormal`. Последующий вызов `SetAlwaysOnTop` по-прежнему считается явным изменением во время выполнения.

**CollectionBehavior** (`MacWindowCollectionBehavior`)

Управляет поведением окна в пространствах macOS Spaces и полноэкранном режиме. Это значения битовой маски, которые можно объединять побитовым ИЛИ (`|`).

**Поведение в Spaces:**

- `MacWindowCollectionBehaviorDefault` — использует FullScreenPrimary (по умолчанию, с обратной совместимостью)
- `MacWindowCollectionBehaviorCanJoinAllSpaces` — окно отображается во всех пространствах Spaces
- `MacWindowCollectionBehaviorMoveToActiveSpace` — при отображении перемещается в активное пространство Space
- `MacWindowCollectionBehaviorManaged` — стандартное поведение управляемого окна
- `MacWindowCollectionBehaviorTransient` — временное окно
- `MacWindowCollectionBehaviorStationary` — остаётся неподвижным при переключении между пространствами Spaces

**Переключение между окнами:**

- `MacWindowCollectionBehaviorParticipatesInCycle` — включается в перебор окон по Cmd+`
- `MacWindowCollectionBehaviorIgnoresCycle` — исключается из перебора окон по Cmd+`

**Поведение в полноэкранном режиме:**

- `MacWindowCollectionBehaviorFullScreenPrimary` — позволяет переходить в полноэкранный режим
- `MacWindowCollectionBehaviorFullScreenAuxiliary` — позволяет отображать окно поверх полноэкранных приложений
- `MacWindowCollectionBehaviorFullScreenNone` — отключает возможность перехода в полноэкранный режим
- `MacWindowCollectionBehaviorFullScreenAllowsTiling` — разрешает размещение окон рядом друг с другом (macOS 10.11+)
- `MacWindowCollectionBehaviorFullScreenDisallowsTiling` — запрещает размещение окон рядом друг с другом (macOS 10.11+)

**Пример — окно в стиле Spotlight:**

```go
// Window that appears on all Spaces AND can overlay fullscreen apps
Mac: application.MacWindow{
	WindowClass: application.MacWindowClassPanel,
	PanelPreferences: application.MacPanelPreferences{
		NonActivating: true,
	},
    CollectionBehavior: application.MacWindowCollectionBehaviorCanJoinAllSpaces |
                        application.MacWindowCollectionBehaviorFullScreenAuxiliary,
    WindowLevel:        application.MacWindowLevelFloating,
},
```

**Пример — один вариант поведения:**

```go
// Window that can appear over fullscreen applications
Mac: application.MacWindow{
    CollectionBehavior: application.MacWindowCollectionBehaviorFullScreenAuxiliary,
},
```

**TabbingMode** (`MacWindowTabbingMode`)

Управляет группировкой окон во вкладки в macOS 10.12 и более поздних версиях. Эта функция позволяет объединять несколько окон в виде вкладок.

**Варианты:**

- `MacWindowTabbingModeDefault` — служебное нулевое значение (явно не задано). Во время выполнения по умолчанию группировка во вкладки запрещена
- `MacWindowTabbingModeAutomatic` — поведение при группировке во вкладки определяет система
- `MacWindowTabbingModePreferred` — для окна предпочтителен режим вкладок
- `MacWindowTabbingModeDisallowed` — отключает группировку окон во вкладки

**Пример — отключение группировки окон во вкладки:**

```go
Mac: application.MacWindow{
    TabbingMode: application.MacWindowTabbingModeDisallowed,
},
```

**Пример — предпочтительный режим вкладок для окна:**

```go
Mac: application.MacWindow{
    TabbingMode: application.MacWindowTabbingModePreferred,
},
```

**WebviewPreferences** (`MacWebviewPreferences`)

Обеспечивает детальное управление базовой конфигурацией `WKWebView`. Все поля необязательны: если поле не задано, соответствующее значение WebKit по умолчанию не изменяется.

```go
Mac: application.MacWindow{
    WebviewPreferences: application.MacWebviewPreferences{
        TabFocusesLinks:                       optional.True,
        TextInteractionEnabled:                optional.True,
        FullscreenEnabled:                     optional.False,
        AllowsBackForwardNavigationGestures:   optional.False,
        AllowsMagnification:                   optional.True,
        AllowsAirPlayForMediaPlayback:         optional.True,
        JavaScriptCanOpenWindowsAutomatically: optional.False,
        MinimumFontSize:                       optional.NewVar(12.0),
        ApplicationNameForUserAgent:           "MyApp",
        EnableAutoplayWithoutUserAction:       optional.True,
    },
},
```

- `TabFocusesLinks` — если задано значение `true`, нажатие Tab перемещает фокус на ссылки и элементы управления форм (по умолчанию: `false`)
- `TextInteractionEnabled` — если задано значение `true`, пользователи могут выделять текст в веб-представлении и взаимодействовать с ним (по умолчанию: `true`)
- `FullscreenEnabled` — если задано значение `true`, веб-содержимое может переходить в полноэкранный режим через HTML Fullscreen API (по умолчанию: `false`). Требуется macOS 12.3 или более поздней версии.
- `AllowsBackForwardNavigationGestures` — если задано значение `true`, горизонтальные жесты смахивания запускают переход назад или вперёд (по умолчанию: `false`)
- `AllowsMagnification` — если задано значение `true`, в веб-представлении включено масштабирование жестом щипка (по умолчанию: `false`)
- `AllowsAirPlayForMediaPlayback` — если задано значение `true`, мультимедийные данные можно передавать потоково на устройства AirPlay (по умолчанию: `true`)
- `JavaScriptCanOpenWindowsAutomatically` — если задано значение `true`, JavaScript может открывать новые окна без действия пользователя (по умолчанию: `false`)
- `MinimumFontSize` — минимальный размер шрифта в пунктах. Для задания используйте `optional.NewVar(12.0)`. Если значение не задано, сохраняется значение WebKit по умолчанию.
- `ApplicationNameForUserAgent` — переопределяет суффикс имени приложения в строке пользовательского агента WebKit. Это полезно, если сайты отклоняют стандартный идентификатор `"wails.io"` (например, при встраивании YouTube). Оставьте поле пустым, чтобы сохранить значение по умолчанию.
- `EnableAutoplayWithoutUserAction` — если задано значение `true`, аудио и видео могут воспроизводиться автоматически без действия пользователя. Соответствует `WKWebViewConfiguration.mediaTypesRequiringUserActionForPlayback = WKAudiovisualMediaTypeNone` (по умолчанию: `false`)

### Параметры Windows (для отдельного окна)

Структура для отдельного окна — `application.WindowsWindow`, а **не** `WindowsOptions` (последняя является структурой уровня *приложения*).

```go
Windows: application.WindowsWindow{
    DisableIcon:                       false,
    DisableMenu:                       false,
    BackdropType:                      application.Auto,
    CustomTheme:                       application.ThemeSettings{},
    DisableFramelessWindowDecorations: false,
    NonClientRegionSupport:            false,
    WebView2CompositionHosting:        false,
},
```

**DisableIcon** (`bool`)

- Удаляет значок из строки заголовка.

**DisableMenu** (`bool`)

- Отключает строку меню для окна. При значении `true` строка меню не отображается в окне, даже если она настроена.
- По умолчанию: `false`

**BackdropType** (`BackdropType`)

- `application.Auto` — системный фон по умолчанию
- `application.None` — без фонового материала
- `application.Mica` — материал Mica (Windows 11)
- `application.Acrylic` — материал Acrylic (Windows 11)
- `application.Tabbed` — материал с вкладками (Windows 11)

Констант в стиле `WindowsBackdropTypeMica` нет — используйте `application.Mica` и т. п.

**CustomTheme** (`ThemeSettings`)

- Значение (не указатель). Пользовательские цвета для тёмного и светлого режимов: цвет рамки окна, текста и фона строки заголовка, а также строки меню.

**DisableFramelessWindowDecorations** (`bool`)

- Отключает стандартное оформление безрамочного окна (тень Aero и скруглённые углы).

**NonClientRegionSupport** (`bool`)

- Включает встроенную в WebView2 поддержку `app-region: drag` / `app-region: no-drag` для пользовательских строк заголовка безрамочных окон.
- Этот параметр предназначен только для простого перетаскивания приложения нативными средствами. Он не обеспечивает нативное поведение пользовательских кнопок заголовка окна или поддержку Windows 11 Snap Assist / Snap Layouts для пользовательских кнопок разворачивания.

**WebView2CompositionHosting** (`bool`)

- Включает управляемую Wails поддержку `--wails-non-client-region` для пользовательских кнопок заголовка с нативным поведением Windows, включая Windows 11 Snap Assist / Snap Layouts для пользовательских кнопок разворачивания.
- Экспериментальная функция. WebView2 размещается через `ICoreWebView2CompositionController` и DirectComposition вместо стандартного контроллера, размещённого в HWND.
- Можно использовать совместно с `NonClientRegionSupport`, если окну одновременно требуются встроенная в WebView2 поддержка `app-region` и управляемые Wails пользовательские области кнопок заголовка.

**Пример:**

```go
Windows: application.WindowsWindow{
    BackdropType: application.Mica,
    DisableIcon:  true,
},
```

**Пример — пользовательские области строки заголовка Windows:**

```go
Windows: application.WindowsWindow{
    NonClientRegionSupport:    true,
    WebView2CompositionHosting: true,
},
```

Подробное описание поведения, компромиссов и соответствующего CSS см. в разделе [«Безрамочные окна Windows»](/features/windows/frameless/#native-non-client-regions-on-windows).

### Параметры Linux (для отдельного окна)

Структура для отдельного окна — `application.LinuxWindow`, а **не** `LinuxOptions`.

```go
Linux: application.LinuxWindow{
    Icon:                []byte{/* PNG data */},
    WindowIsTranslucent: false,
},
```

**Icon** (`[]byte`)

- Значок окна (в формате PNG).

**WindowIsTranslucent** (`bool`)

- Требуется поддержка со стороны композитора.

**Пример:**

```go
//go:embed icon.png
var icon []byte

Linux: application.LinuxWindow{
    Icon: icon,
},
```

## Параметры Windows уровня приложения

Некоторые параметры, относящиеся к Windows, необходимо настраивать на уровне приложения, а не для отдельных окон. Причина в том, что WebView2 использует одну общую среду браузера для каждого пути к пользовательским данным.

### Флаги браузера

Флаги браузера WebView2 управляют экспериментальными функциями и поведением во **всех окнах** приложения. Их необходимо задавать в `application.Options.Windows`:

```go
app := application.New(application.Options{
    Name: "My App",
    Windows: application.WindowsOptions{
        // Enable experimental WebView2 features
        EnabledFeatures: []string{
            "msWebView2EnableDraggableRegions",
        },

        // Disable specific features
        DisabledFeatures: []string{
            "msSmartScreenProtection",  // Always disabled by Wails
        },

        // Additional Chromium command-line arguments
        AdditionalBrowserArgs: []string{
            "--disable-gpu",
            "--remote-debugging-port=9222",
        },
    },
})
```

**EnabledFeatures** (`[]string`)

- Список флагов функций WebView2, которые необходимо включить
- Доступные флаги перечислены в разделе [«Флаги браузера WebView2»](https://learn.microsoft.com/en-us/microsoft-edge/webview2/concepts/webview-features-flags)
- Пример: `"msWebView2EnableDraggableRegions"`

**DisabledFeatures** (`[]string`)

- Список флагов функций WebView2, которые необходимо отключить
- Wails автоматически отключает `msSmartScreenProtection`
- Пример: `"msExperimentalFeature"`

**AdditionalBrowserArgs** (`[]string`)

- Аргументы командной строки Chromium, передаваемые процессу браузера
- Необходимо указывать префикс `--` (например, `"--remote-debugging-port=9222"`)
- Доступные аргументы перечислены в разделе [«Параметры командной строки Chromium»](https://peter.sh/experiments/chromium-command-line-switches/)

@note{type="caution" title="Важно"}
Эти флаги применяются глобально ко ВСЕМ окнам, поскольку WebView2 использует одну среду браузера для каждого пути к пользовательским данным. Задать разные флаги браузера для разных окон невозможно.

@end

**Полный пример:**

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "My App",
        Windows: application.WindowsOptions{
            // Enable draggable regions feature
            EnabledFeatures: []string{
                "msWebView2EnableDraggableRegions",
            },
            // Enable remote debugging
            AdditionalBrowserArgs: []string{
                "--remote-debugging-port=9222",
            },
        },
    })

    // All windows will use the browser flags configured above
    window := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:  "Main Window",
        Width:  1024,
        Height: 768,
    })

    window.Show()
    app.Run()
}
```

## Полный пример

Ниже приведена конфигурация окна, готовая к использованию в рабочей среде:

```go
package main

import (
    _ "embed"
    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var assets embed.FS

//go:embed icon.png
var icon []byte

func main() {
    app := application.New(application.Options{
        Name: "My Application",
    })

    window := app.Window.NewWithOptions(application.WebviewWindowOptions{
        // Identity
        Name:  "main-window",
        Title: "My Application",

        // Size and Position
        Width:     1200,
        Height:    800,
        MinWidth:  800,
        MinHeight: 600,

        // Initial State
        StartState: application.WindowStateNormal,

        // Appearance
        BackgroundColour: application.NewRGB(255, 255, 255),

        // Platform-Specific (per-window structs)
        Mac: application.MacWindow{
            TitleBar: application.MacTitleBar{
                AppearsTransparent: true,
            },
            Backdrop: application.MacBackdropTranslucent,
        },

        Windows: application.WindowsWindow{
            BackdropType: application.Mica,
            DisableIcon:  false,
        },

        Linux: application.LinuxWindow{
            Icon: icon,
        },
    })

    window.Center()
    window.Show()

    app.Run()
}
```

Ресурсы фронтенда обслуживаются на уровне приложения (`application.Options.Assets`), а не отдельно для каждого окна.

## Дальнейшие шаги

- [Основы работы с окнами](/features/windows/basics/) — создание окон и управление ими
- [Несколько окон](/features/windows/multiple/) — шаблоны работы с несколькими окнами
- [Окна без рамки](/features/windows/frameless/) — пользовательское оформление окна
- [События окна](/features/windows/events/) — события жизненного цикла

---

**Остались вопросы?** Задайте их в [Discord](https://discord.gg/JDdSxwjhGf) или ознакомьтесь с [примерами](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples).
