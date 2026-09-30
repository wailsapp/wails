---
title: "Параметры окна"
description: "Полное справочное руководство по WebviewWindowOptions"
slug: "features/windows/options"
sourcePath: "features/windows/options.md"
---

## Параметры конфигурации окна

Wails предоставляет широкие возможности настройки окна с десятками параметров размера, положения, внешнего вида и поведения. Это **полное справочное руководство** по `WebviewWindowOptions` охватывает все доступные параметры для Windows, macOS и Linux. Все параметры и платформы с примерами и ограничениями.

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
    KeyBindings                 map[string]func(window *WebviewWindow)

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

В `WebviewWindowOptions` **нет** поля `Parent` — для связей между родительскими и модальными окнами используйте `parentWindow.AttachModal(childWindow)`. В этой структуре также **нет** поля `Assets` — ресурсы настраиваются в `application.Options` (`Assets AssetOptions`).

Полный исходный код: [`v3/pkg/application/webview_window_options.go`](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window_options.go).

## Основные параметры

### Name

**Тип:** `string` **Значение по умолчанию:** автоматически сгенерированный UUID **Платформы:** все

```go
Name: "main-window"
```

**Назначение:** уникальный идентификатор, позволяющий впоследствии найти окно.

**Рекомендации:**

- Используйте информативные имена: `"main"`, `"settings"`, `"about"`
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

**Тип:** `string` **Значение по умолчанию:** имя приложения **Платформы:** все

```go
Title: "My Application"
```

**Назначение:** текст, отображаемый в строке заголовка и на панели задач.

**Динамическое обновление:**

```go
window.SetTitle("My Application - Document.txt")
```

### Width / Height

**Тип:** `int` (пиксели) **Значение по умолчанию:** 800 × 600 **Платформы:** все **Ограничения:** значение должно быть положительным

```go
Width:  1200,
Height: 800,
```

**Назначение:** начальный размер окна в логических пикселях.

**Примечания:**

- Wails автоматически учитывает масштабирование DPI
- Используйте логические, а не физические пиксели
- Учитывайте минимальное разрешение экрана (1024x768)

**Примеры размеров:**

| Сценарий использования | Ширина | Высота |
| --- | --- | --- |
| Небольшая утилита | 400 | 300 |
| Стандартное приложение | 1024 | 768 |
| Большое приложение | 1440 | 900 |
| Full HD | 1920 | 1080 |

### X / Y

**Тип:** `int` (пиксели) **Значение по умолчанию:** по центру экрана **Платформы:** все

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

`X` и `Y` действуют только в том случае, если задано `InitialPosition: application.WindowXY`. В противном случае для `InitialPosition` по умолчанию используется `WindowCentered`, а `X`/`Y` игнорируются.

```go
settings := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:            "coordinate-window",
    InitialPosition: application.WindowXY, // opt into X/Y coordinates
    X:               100,
    Y:               100,
})
```

**Рекомендация:** если конкретные координаты не важны, после создания окна центрируйте его с помощью `Center()`:

```go
window := app.Window.New()
window.Center()
```

### MinWidth / MinHeight

**Тип:** `int` (пиксели) **Значение по умолчанию:** 0 (минимум не задан) **Платформы:** все

```go
MinWidth:  400,
MinHeight: 300,
```

**Назначение:** не допускать чрезмерного уменьшения окна.

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

**Тип:** `int` (пиксели) **Значение по умолчанию:** 0 (максимум не задан) **Платформы:** все

```go
MaxWidth:  1920,
MaxHeight: 1080,
```

**Назначение:** не допускать чрезмерного увеличения окна.

**Сценарии использования:**

- Приложения с фиксированным размером
- Предотвращение чрезмерного потребления ресурсов
- Соблюдение ограничений дизайна

## Параметры состояния

### Hidden

**Тип:** `bool` **Значение по умолчанию:** `false` **Платформы:** все

```go
Hidden: true,
```

**Назначение:** создать окно, не отображая его.

**Варианты использования:**

- Фоновые окна
- Окна, отображаемые по запросу
- Заставки (создать, загрузить, затем показать)
- Предотвращение белой вспышки при загрузке содержимого

**Улучшения для отдельных платформ:**

- **Windows:** устранена белая вспышка окна — окно остаётся невидимым до вызова `Show()`
- **macOS:** полная поддержка
- **Linux:** полная поддержка

**Рекомендуемый способ плавной загрузки:**

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

**Тип:** `bool` **Значение по умолчанию:** `false` **Платформы:** все

```go
Frameless: true,
```

**Назначение:** убрать строку заголовка и рамки окна.

**Варианты использования:**

- Нестандартное оформление окна
- Заставки
- Киоск-приложения
- Окна с индивидуальным дизайном

**Важно:** вам потребуется реализовать:

- Перетаскивание окна
- Кнопки закрытия, сворачивания и разворачивания
- Маркеры изменения размера (если размер окна можно изменять)

**Подробности см. в разделе [«Окна без рамок»](/features/windows/frameless/).**

### DisableResize

**Тип:** `bool` **Значение по умолчанию:** `false` (по умолчанию размер окна можно изменять) **Платформы:** все

```go
DisableResize: true,
```

**Назначение:** запретить изменение размера окна. Обратите внимание: это поле имеет значение, **обратное** значению поля `Resizable` в v2. Чтобы запретить изменение размера окна, задайте `DisableResize: true`.

**Варианты использования:**

- Приложения с фиксированным размером окна
- Заставки
- Диалоговые окна

**Примечание:** пользователи по-прежнему смогут разворачивать окно или переводить его в полноэкранный режим, если эти возможности также не отключены с помощью `MaximiseButtonState` или настройки поведения коллекции.

### AlwaysOnTop

**Тип:** `bool` **Значение по умолчанию:** `false` **Платформы:** все

```go
AlwaysOnTop: true,
```

**Назначение:** отображать окно поверх всех остальных окон.

**Варианты использования:**

- Плавающие панели инструментов
- Уведомления
- Режим «картинка в картинке»
- Таймеры

**Примечания для отдельных платформ:**

- **macOS:** полная поддержка
- **Windows:** полная поддержка
- **Linux:** зависит от оконного менеджера

### StartState

**Тип:** перечисление `WindowState` **Значение по умолчанию:** `WindowStateNormal` **Платформы:** все

```go
StartState: application.WindowStateMaximised,
```

**Назначение:** начальное состояние окна при его отображении.

**Значения:**

- `WindowStateNormal` — обычное окно
- `WindowStateMinimised` — свёрнутое окно
- `WindowStateMaximised` — развёрнутое окно
- `WindowStateFullscreen` — полноэкранный режим

Константы `WindowStateHidden` нет — чтобы при запуске окно было невидимым, используйте логическое поле `Hidden`.

**Переключение полноэкранного режима во время выполнения:**

```go
window.Fullscreen()
window.UnFullscreen()
window.ToggleFullscreen() // there is no SetFullscreen(bool)
```

## Параметры внешнего вида

### BackgroundColour

**Тип:** структура `RGBA` **Значение по умолчанию:** белый цвет **Платформы:** все

```go
BackgroundColour: application.RGBA{Red: 0, Green: 0, Blue: 0, Alpha: 255},
```

Поля `RGBA` — `Red, Green, Blue, Alpha` (тип uint8). Рекомендуется использовать вспомогательные функции `application.NewRGB(r, g, b)` (альфа-канал 255) или `application.NewRGBA(r, g, b, a)`.

**Назначение:** цвет фона окна до загрузки содержимого.

**Варианты использования:**

- Соответствие теме приложения
- Предотвращение белой вспышки при использовании тёмных тем
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

**Тип:** перечисление `BackgroundType` **Значение по умолчанию:** `BackgroundTypeSolid` **Платформы:** macOS, Windows (частичная поддержка)

```go
BackgroundType: application.BackgroundTypeTranslucent,
```

**Значения:**

- `BackgroundTypeSolid` — сплошной цвет
- `BackgroundTypeTransparent` — полная прозрачность
- `BackgroundTypeTranslucent` — полупрозрачное размытие

**Поддержка платформ:**

- **macOS:** настройте `Mac.Backdrop`; для прозрачности веб-представления требуется [`private_mac_apis`](/guides/build/private-macos-apis/#webview-transparency-and-background). Без этого веб-представление остаётся непрозрачным.
- **Windows:** прозрачный и полупрозрачный режимы (Windows 11 и более поздние версии)
- **Linux:** только сплошной цвет

**Пример (macOS):**

```go
BackgroundType: application.BackgroundTypeTranslucent,
Mac: application.MacWindow{
    Backdrop: application.MacBackdropTranslucent,
},
```

### OpenInspectorOnStartup и OpenDevTools

**Закрытый API в macOS:** для программного открытия инспектора с помощью `OpenInspectorOnStartup: true`, Go `window.OpenDevTools()` и JavaScript `Window.OpenDevTools()` требуется `private_mac_apis`. Без него эти операции ничего не делают. Для производственных сборок также требуется `devtools`. Общедоступная инспекция Safari в macOS 13.3 и более поздних версиях не требует закрытых API; для включения инспектора в более ранних версиях macOS они необходимы. См. [матрицу сборок Web Inspector](/guides/build/private-macos-apis/#web-inspector).

## Параметры содержимого

### URL

**Тип:** `string` **Значение по умолчанию:** пустое (загрузка из Assets) **Платформы:** все

```go
URL: "https://example.com",
```

**Назначение:** загрузка внешнего URL вместо встроенных ресурсов.

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

Фрагмент конфигурации уровня приложения для производственной среды:

```go
app := application.New(application.Options{
    Name: "My App",
    Assets: application.AssetOptions{
        Handler: application.AssetFileServerFS(assets),
    },
})
```

### HTML

**Тип:** `string` **Значение по умолчанию:** пустое **Платформы:** все

```go
HTML: "<h1>Hello World</h1>",
```

**Назначение:** непосредственная загрузка строки HTML.

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

### Assets (только на уровне приложения)

Конфигурация ресурсов **не является** полем `WebviewWindowOptions`. Ресурсы фронтенда обслуживает само приложение через `application.Options.Assets` (`AssetOptions`); каждое окно использует этот сервер ресурсов.

```go
app := application.New(application.Options{
    Assets: application.AssetOptions{
        Handler: application.AssetFileServerFS(assets),
    },
})
```

**Подробные сведения см. в разделе [«Система сборки»](/concepts/build-system/).**

### UseApplicationMenu

**Тип:** `bool` **Значение по умолчанию:** `false` **Платформы:** Windows, Linux (не действует в macOS)

```go
UseApplicationMenu: true,
```

**Назначение:** использовать для этого окна меню приложения, заданное через `app.Menu.Set()`.

В **macOS** этот параметр не действует, поскольку macOS всегда использует глобальное меню приложения в верхней части экрана.

В **Windows** и **Linux** окна по умолчанию не отображают меню. Значение `UseApplicationMenu: true` предписывает окну использовать меню уровня приложения, обеспечивая простое кроссплатформенное решение.

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
- Это упрощает кроссплатформенный код, устраняя необходимость проверять операционную систему во время выполнения
- Полную документацию по меню см. в разделе [«Меню приложения»](/features/menus/application/)

## Параметры ввода

### EnableFileDrop

**Тип:** `bool` **По умолчанию:** `false` **Платформа:** все

```go
EnableFileDrop: true,
```

**Назначение:** разрешает перетаскивание файлов из операционной системы в окно.

Если этот параметр включён:

- Файлы из файловых менеджеров можно перетаскивать в приложение
- Событие `WindowFilesDropped` срабатывает и передаёт пути к перетащенным файлам
- Элементы с атрибутом `data-file-drop-target` предоставляют подробные сведения о перетаскивании

**Варианты использования:**

- Интерфейсы загрузки файлов
- Редакторы документов
- Средства импорта мультимедиа
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

**Тип:** `bool` **По умолчанию:** `false` **Платформа:** Windows (10+), macOS

```go
ContentProtectionEnabled: true,
```

**Назначение:** предотвращает захват содержимого окна с экрана.

**Поддержка платформ:**

- **Windows:** Windows 10, сборка 19041+ (полная поддержка), более ранние версии (частичная поддержка)
- **macOS:** полная поддержка
- **Linux:** не поддерживается

**Варианты использования:**

- Банковские приложения
- Менеджеры паролей
- Медицинские карты
- Конфиденциальные документы

**Важные замечания:**

1. Не предотвращает фотографирование экрана физической камерой
2. Некоторые инструменты могут обходить защиту
3. Является частью комплексной системы безопасности, а не единственной мерой защиты
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

**Тип:** `map[PermissionType]Permission` **По умолчанию:** `nil` (обработка по умолчанию для платформы) **Платформа:** Linux, Windows (в macOS обработка передаётся TCC)

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionAllow,
    application.PermissionCamera:     application.PermissionAllow,
},
```

**Назначение:** позволяет декларативно, без платформозависимого кода, управлять обработкой запросов веб-содержимого окна на доступ к возможностям устройства: камере, микрофону, геолокации, уведомлениям и чтению буфера обмена.

**Значения PermissionType:** `PermissionMicrophone`, `PermissionCamera`, `PermissionGeolocation`, `PermissionNotifications`, `PermissionClipboardRead`

**Значения разрешений:**

- `PermissionDefault` (0) — встроенная обработка платформы: системный запрос ОС/WebView2 в macOS/Windows; в Linux доступ к камере и микрофону разрешён, а ко всему остальному — запрещён
- `PermissionAllow` (1) — предоставить доступ без запроса (в Linux реализован только доступ к камере и микрофону; доступ к остальным типам по-прежнему запрещён)
- `PermissionDeny` (2) — запретить доступ без запроса

**Важно — Windows:** до появления этого параметра Wails без уведомления предоставлял доступ ко всем возможностям WebView2. Теперь добавление любой записи в `Permissions` отключает такое безусловное предоставление доступа. Для возможностей, не указанных в списке, WebView2 будет показывать встроенный запрос вместо автоматического предоставления доступа. Явно перечислите все возможности, необходимые приложению.

**Полное руководство, таблицу поддержки платформ и примеры см. в разделе [«Разрешения»](/features/windows/permissions/).**

## События жизненного цикла окна

События жизненного цикла окна обрабатываются с помощью `OnWindowEvent` и `RegisterHook`. Эти методы позволяют точно управлять поведением при закрытии и уничтожении окна.

### Отмена закрытия окна

Чтобы предотвратить закрытие окна (например, при наличии несохранённых изменений), используйте `RegisterHook` с событием `WindowClosing` и вызовите `event.Cancel()`:

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
- После отмены закрытия окно останется открытым

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

- `OnWindowEvent` обрабатывает события, которые вот-вот произойдут
- Очистка выполняется до уничтожения окна
- Отменить закрытие здесь нельзя (для этого используйте `RegisterHook`)

### Шаблон очистки для окна в единственном экземпляре

Для окон в единственном экземпляре (чтобы существовал только один экземпляр) используйте `WindowClosing` для очистки ссылки:

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

## Параметры для отдельных платформ

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
- `FullSizeContent` — расширяет содержимое до полного размера окна

В macOS для прозрачности webview и программного открытия инспектора требуется тег сборки `private_mac_apis`. Без него те же параметры остаются допустимыми, но операции, использующие закрытые API, ничего не делают. Группировка Liquid Glass также игнорируется, а для стилей используются общедоступные альтернативы. Команды сборки и точное описание поведения приведены в разделе [Закрытые API macOS](/guides/build/private-macos-apis/).

**Backdrop** (`MacBackdrop`)

- `MacBackdropNormal` — стандартный непрозрачный фон
- `MacBackdropTranslucent` — **Для прозрачности webview требуется закрытый API.** Без тега нативное размытие остаётся позади непрозрачного webview.
- `MacBackdropTransparent` — **Для прозрачности webview требуется закрытый API.** Без тега webview остаётся непрозрачным.
- `MacBackdropLiquidGlass` — **Для прозрачности webview требуется закрытый API.** Без тега стеклянный слой остаётся позади непрозрачного webview, а для стилей используются общедоступные альтернативы.

**LiquidGlass** (`MacLiquidGlass`)

| Поле или значение | Зависимость от закрытого API в macOS |
| --- | --- |
| `Style: LiquidGlassStyleAutomatic` | Нативный обычный стиль является общедоступным; для прозрачности webview с фоном требуется `private_mac_apis`. |
| `Style: LiquidGlassStyleLight` | Тег сохраняет существующее сопоставление с нативным прозрачным стилем; без него Wails использует обычное стекло с оформлением Aqua. |
| `Style: LiquidGlassStyleDark` | **Закрытый API:** недокументированное значение нативного стиля `2`; без тега используется обычное стекло с оформлением Dark Aqua. |
| `Style: LiquidGlassStyleVibrant` | Сопоставление с нативным прозрачным стилем является общедоступным; для прозрачности webview с фоном требуется тег. |
| `GroupID` | **Закрытый API:** непустые значения запрашивают группировку; без тега они игнорируются. |
| `GroupSpacing` | **Закрытый API:** положительные значения запрашивают установку интервала для группы; без тега они игнорируются. |
| `Material`, `CornerRadius`, `TintColor` | Сами по себе не зависят от закрытых API. |

Нативные значения стилей и сведения о доступности в версиях ОС приведены в разделе [Значения Liquid Glass](/guides/build/private-macos-apis/#liquid-glass-values).

**InvisibleTitleBarHeight** (`int`)

- Высота невидимой области строки заголовка (для перетаскивания)
- Действует только тогда, когда нативная область строки заголовка для перетаскивания скрыта, то есть когда окно не имеет рамки (`Frameless: true`) или использует прозрачную строку заголовка (`AppearsTransparent: true`)
- Не влияет на стандартные окна с видимой строкой заголовка

**WindowClass** (`MacWindowClass`)

- `MacWindowClassWindow` — стандартное поведение `NSWindow` (по умолчанию)
- `MacWindowClassPanel` — вспомогательное `NSPanel`, которое никогда не становится главным окном приложения

`PanelPreferences` применяется только к `MacWindowClassPanel`:

- `NonActivating` добавляет `NSWindowStyleMaskNonactivatingPanel`. Отображение панели или перевод на неё фокуса не активирует приложение Wails, однако панель по-прежнему может становиться ключевой для работы с элементами управления и ввода текста.
- `FloatingPanel` включает поведение плавающей панели AppKit.
- `BecomesKeyOnlyIfNeeded` получает статус ключевого окна, только если представление, по которому щёлкнули, запрашивает ввод с клавиатуры.
- `UtilityWindow` применяет нативный стиль служебного окна.

Панели Wails остаются видимыми после деактивации приложения и освобождаются при закрытии, что соответствует жизненному циклу, ожидаемому для `WebviewWindow`. Эти настройки намеренно переопределяют противоположные значения `NSPanel` по умолчанию.

Класс окна, уровень, политика активации и поведение коллекции решают разные задачи:

- `WindowClass` выбирает `NSWindow` или `NSPanel` и управляет семантикой главного и ключевого окна.
- `WindowLevel` управляет порядком по оси Z. Для наложений над строкой меню используйте `MacWindowLevelPopUpMenu`.
- `MacOptions.ActivationPolicy` управляет приложением в целом, включая отображение в Dock и строке меню. Для неактивирующей панели не требуется политика активации вспомогательного приложения, хотя приложение всё же может использовать её, чтобы скрыть свой значок в Dock.
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

- `MacWindowLevelNormal` — Стандартный уровень окна (по умолчанию)
- `MacWindowLevelFloating` — Располагается поверх обычных окон
- `MacWindowLevelTornOffMenu` — Уровень отделённого меню
- `MacWindowLevelModalPanel` — Уровень модальной панели
- `MacWindowLevelMainMenu` — Уровень главного меню
- `MacWindowLevelStatus` — Уровень окна состояния
- `MacWindowLevelPopUpMenu` — Уровень всплывающего меню
- `MacWindowLevelScreenSaver` — Уровень экранной заставки

Явно заданное значение `WindowLevel` имеет приоритет над `AlwaysOnTop` и `PanelPreferences.FloatingPanel`. Если уровень не задан явно, при включённом `AlwaysOnTop` или использовании плавающей панели уровень окна становится `MacWindowLevelFloating`; в противном случае уровень — `MacWindowLevelNormal`. Последующий вызов `SetAlwaysOnTop` по-прежнему считается явным изменением во время выполнения.

**CollectionBehavior** (`MacWindowCollectionBehavior`)

Управляет поведением окна в пространствах macOS Spaces и полноэкранном режиме. Это значения битовой маски, которые можно объединять с помощью побитового ИЛИ (`|`).

**Поведение в Spaces:**

- `MacWindowCollectionBehaviorDefault` — Использует FullScreenPrimary (по умолчанию, с обратной совместимостью)
- `MacWindowCollectionBehaviorCanJoinAllSpaces` — Окно отображается во всех пространствах Spaces
- `MacWindowCollectionBehaviorMoveToActiveSpace` — При отображении перемещается в активное пространство Space
- `MacWindowCollectionBehaviorManaged` — Стандартное поведение управляемого окна
- `MacWindowCollectionBehaviorTransient` — Временное/переходное окно
- `MacWindowCollectionBehaviorStationary` — Остаётся неподвижным при переключении между пространствами Spaces

**Циклическое переключение окон:**

- `MacWindowCollectionBehaviorParticipatesInCycle` — Включено в цикл переключения по Cmd+`
- `MacWindowCollectionBehaviorIgnoresCycle` — Исключено из цикла переключения по Cmd+`

**Поведение в полноэкранном режиме:**

- `MacWindowCollectionBehaviorFullScreenPrimary` — Может переходить в полноэкранный режим
- `MacWindowCollectionBehaviorFullScreenAuxiliary` — Может отображаться поверх полноэкранных приложений
- `MacWindowCollectionBehaviorFullScreenNone` — Отключает возможность полноэкранного режима
- `MacWindowCollectionBehaviorFullScreenAllowsTiling` — Разрешает размещение рядом (macOS 10.11+)
- `MacWindowCollectionBehaviorFullScreenDisallowsTiling` — Запрещает размещение рядом (macOS 10.11+)

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

Управляет поведением вкладок окон в macOS 10.12 и более поздних версиях. Вкладки окон позволяют группировать несколько окон в виде вкладок.

**Варианты:**

- `MacWindowTabbingModeDefault` — Сигнальное нулевое значение (явно не задано). Во время выполнения по умолчанию запрещает вкладки окон
- `MacWindowTabbingModeAutomatic` — Система определяет поведение вкладок окон
- `MacWindowTabbingModePreferred` — Для окна предпочтителен режим вкладок
- `MacWindowTabbingModeDisallowed` — Отключает вкладки окон

**Пример — отключение вкладок окон:**

```go
Mac: application.MacWindow{
    TabbingMode: application.MacWindowTabbingModeDisallowed,
},
```

**Пример — предпочтительное использование вкладок окон:**

```go
Mac: application.MacWindow{
    TabbingMode: application.MacWindowTabbingModePreferred,
},
```

**WebviewPreferences** (`MacWebviewPreferences`)

Точное управление базовой конфигурацией `WKWebView`. Все поля необязательны: незаданные поля не изменяют значения WebKit по умолчанию.

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
- `AllowsAirPlayForMediaPlayback` — если задано значение `true`, мультимедиа можно передавать на устройства AirPlay (по умолчанию: `true`)
- `JavaScriptCanOpenWindowsAutomatically` — если задано значение `true`, JavaScript может открывать новые окна без действия пользователя (по умолчанию: `false`)
- `MinimumFontSize` — минимальный размер шрифта в пунктах. Для задания значения используйте `optional.NewVar(12.0)`. Если значение не задано, сохраняется значение WebKit по умолчанию.
- `ApplicationNameForUserAgent` — переопределяет суффикс имени приложения в строке user-agent WebKit. Полезно, когда сайты отклоняют стандартный идентификатор `"wails.io"` (например, во встроенных проигрывателях YouTube). Оставьте пустым, чтобы сохранить значение по умолчанию.
- `EnableAutoplayWithoutUserAction` — если задано значение `true`, аудио и видео могут воспроизводиться автоматически без действия пользователя. Соответствует `WKWebViewConfiguration.mediaTypesRequiringUserActionForPlayback = WKAudiovisualMediaTypeNone` (по умолчанию: `false`)

### Параметры Windows (для отдельного окна)

Структура для отдельного окна — `application.WindowsWindow`, а **не** `WindowsOptions` (это структура уровня *приложения*).

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

- Отключает строку меню окна. При значении `true` окно не отображает строку меню, даже если она настроена.
- По умолчанию: `false`

**BackdropType** (`BackdropType`)

- `application.Auto` — системное значение по умолчанию
- `application.None` — без фонового материала
- `application.Mica` — материал Mica (Windows 11)
- `application.Acrylic` — материал Acrylic (Windows 11)
- `application.Tabbed` — материал с вкладками (Windows 11)

Констант в стиле `WindowsBackdropTypeMica` нет — используйте `application.Mica` и т. д.

**CustomTheme** (`ThemeSettings`)

- Значение (не указатель). Пользовательские цвета для тёмного и светлого режимов: цвета рамки окна, текста и фона строки заголовка, а также строки меню.

**DisableFramelessWindowDecorations** (`bool`)

- Отключает стандартное оформление безрамочного окна (тень Aero, скруглённые углы).

**NonClientRegionSupport** (`bool`)

- Включает встроенную в WebView2 поддержку `app-region: drag` / `app-region: no-drag` для пользовательских строк заголовка безрамочных окон.
- Этот параметр предназначен только для простого перетаскивания приложения средствами ОС. Он не обеспечивает нативное поведение пользовательских кнопок заголовка окна или поддержку функций Windows 11 Snap Assist / Snap Layouts для пользовательских кнопок развёртывания.

**WebView2CompositionHosting** (`bool`)

- Включает управляемую Wails поддержку `--wails-non-client-region` для пользовательских кнопок заголовка окна с нативным поведением Windows, включая функции Windows 11 Snap Assist / Snap Layouts для пользовательских кнопок развёртывания.
- Экспериментальная возможность. При её использовании WebView2 размещается через `ICoreWebView2CompositionController` и DirectComposition вместо стандартного контроллера, размещённого в HWND.
- Можно использовать совместно с `NonClientRegionSupport`, если окну одновременно требуются встроенная в WebView2 поддержка `app-region` и управляемые Wails пользовательские области кнопок заголовка окна.

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

Подробное описание поведения, компромиссов и соответствующих стилей CSS см. в разделе [«Безрамочные окна»](/features/windows/frameless/#native-non-client-regions-on-windows).

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

## Параметры Windows на уровне приложения

Некоторые параметры, относящиеся к Windows, необходимо настраивать на уровне приложения, а не для каждого окна. Это связано с тем, что WebView2 использует одну общую среду браузера для каждого пути к пользовательским данным.

### Флаги браузера

Флаги браузера WebView2 управляют экспериментальными функциями и поведением **всех окон** приложения. Их необходимо задать в `application.Options.Windows`:

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

- Список включаемых флагов функций WebView2
- Доступные флаги перечислены в разделе [«Флаги браузера WebView2»](https://learn.microsoft.com/en-us/microsoft-edge/webview2/concepts/webview-features-flags)
- Пример: `"msWebView2EnableDraggableRegions"`

**DisabledFeatures** (`[]string`)

- Список отключаемых флагов функций WebView2
- Wails автоматически отключает `msSmartScreenProtection`
- Пример: `"msExperimentalFeature"`

**AdditionalBrowserArgs** (`[]string`)

- Аргументы командной строки Chromium, передаваемые процессу браузера
- Необходимо указывать префикс `--` (например, `"--remote-debugging-port=9222"`)
- Доступные аргументы перечислены в разделе [«Ключи командной строки Chromium»](https://peter.sh/experiments/chromium-command-line-switches/)

@note{type="caution" title="Важно"}
Эти флаги применяются глобально ко ВСЕМ окнам, поскольку WebView2 использует одну общую среду браузера для каждого пути к пользовательским данным. Для разных окон нельзя задать разные флаги браузера.

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

Ниже приведена готовая к использованию в рабочей среде конфигурация окна:

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

Ресурсы фронтенда обслуживаются на уровне приложения (`application.Options.Assets`), а не для каждого окна.

## Дальнейшие шаги

- [Основы работы с окнами](/features/windows/basics/) — создание окон и управление ими
- [Несколько окон](/features/windows/multiple/) — шаблоны работы с несколькими окнами
- [Окна без рамки](/features/windows/frameless/) — пользовательское оформление окна
- [События окна](/features/windows/events/) — события жизненного цикла

---

**Есть вопросы?** Задайте их в [Discord](https://discord.gg/JDdSxwjhGf) или ознакомьтесь с [примерами](https://github.com/wailsapp/wails/tree/master/v3/examples).
