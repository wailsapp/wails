---
title: "API менеджеров"
description: "Организованная структура API со специализированными интерфейсами менеджеров"
slug: "concepts/manager-api"
sourcePath: "concepts/manager-api.md"
---

API менеджеров Wails v3 предоставляет организованный и удобный для изучения доступ к функциям приложения через специализированные структуры менеджеров, сгруппированные в открытых полях `*application.App`. В Wails 3 отказались от обратной совместимости с API v2: прослойки-обёртки для каждого вызова, обеспечивающей совместимость со старым API в стиле `app.NewWebviewWindow(...)`, нет, поэтому управлять приложением следует с помощью перечисленных ниже менеджеров.

## Обзор

API менеджеров разделяет функции приложения на двенадцать специализированных областей (один логгер и одиннадцать менеджеров):

- **`app.Window`** — создание окон, управление ими и обратные вызовы
- **`app.ContextMenu`** — регистрация контекстных меню и управление ими\
- **`app.KeyBinding`** — управление глобальными сочетаниями клавиш
- **`app.Browser`** — интеграция с браузером (открытие URL-адресов и файлов)
- **`app.Env`** — сведения о среде и состоянии системы
- **`app.Dialog`** — операции с диалоговыми окнами выбора файлов и диалоговыми окнами сообщений
- **`app.Event`** — обработка пользовательских событий и событий приложения
- **`app.Menu`** — управление меню приложения
- **`app.Screen`** — управление экранами и преобразование координат
- **`app.Clipboard`** — операции с текстом в буфере обмена
- **`app.SystemTray`** — создание значков в области уведомлений и управление ими
- **`app.Autostart`** — регистрация приложения для запуска при входе пользователя в систему

## Преимущества

- **Упрощённый поиск возможностей** — автодополнение в IDE отображает организованный интерфейс API
- **Улучшенная организация кода** — связанные методы сгруппированы вместе
- **Улучшенная сопровождаемость** — ответственность разделена между менеджерами
- **Возможность дальнейшего расширения** — новые функции проще добавлять в определённые области

## Использование

API менеджеров предоставляет организованный доступ ко всем функциям приложения:

```go
// Events and custom event handling
app.Event.Emit("custom", data)
app.Event.On("custom", func(e *CustomEvent) { ... })

// Window management
window, _ := app.Window.GetByName("main")
app.Window.OnCreate(func(window Window) { ... })

// Browser integration
app.Browser.OpenURL("https://wails.io")

// Menu management
menu := app.Menu.New()
app.Menu.Set(menu)

// System tray
systray := app.SystemTray.New()
```

## Справочник по менеджерам

### Менеджер окон

Управляет созданием и получением окон, а также обратными вызовами их жизненного цикла.

```go
// Create windows
window := app.Window.New()
window := app.Window.NewWithOptions(options)
current := app.Window.Current()

// Find windows
window, exists := app.Window.GetByName("main")
windows := app.Window.GetAll()

// Window callbacks
app.Window.OnCreate(func(window Window) {
    // Handle window creation
})
```

### Менеджер событий

Обрабатывает пользовательские события и отслеживает события приложения.

```go
// Custom events
app.Event.Emit("userAction", data)
cancelFunc := app.Event.On("userAction", func(e *CustomEvent) {
    // Handle event
})
app.Event.Off("userAction")
app.Event.Reset() // Remove all listeners

// Application events
app.Event.OnApplicationEvent(events.Common.ThemeChanged, func(e *ApplicationEvent) {
    // Handle system theme change
})
```

### Менеджер браузера

Обеспечивает интеграцию с браузером для открытия URL-адресов и файлов.

```go
// Open URLs and files in default browser
err := app.Browser.OpenURL("https://wails.io")
err := app.Browser.OpenFile("/path/to/document.pdf")
```

### Менеджер среды

Предоставляет доступ к сведениям о системной среде.

```go
// Get environment info
env := app.Env.Info()
fmt.Printf("OS: %s, Arch: %s\n", env.OS, env.Arch)

// Check system theme
if app.Env.IsDarkMode() {
    // Dark mode is active
}

// Open file manager
err := app.Env.OpenFileManager("/path/to/folder", false)
```

### Менеджер диалоговых окон

Предоставляет организованный доступ к диалоговым окнам выбора файлов и диалоговым окнам сообщений.

```go
// File dialogs
result, err := app.Dialog.OpenFile().
    AddFilter("Text Files", "*.txt").
    PromptForSingleSelection()

result, err = app.Dialog.SaveFile().
    SetFilename("document.txt").
    PromptForSingleSelection()

// Message dialogs
app.Dialog.Info().
    SetTitle("Information").
    SetMessage("Operation completed successfully").
    Show()

app.Dialog.Error().
    SetTitle("Error").
    SetMessage("An error occurred").
    Show()
```

### Менеджер меню

Создание меню приложения и управление ими.

```go
// Create and set application menu
menu := app.Menu.New()
fileMenu := menu.AddSubmenu("File")
fileMenu.Add("New").OnClick(func(ctx *Context) {
    // Handle menu click
})

app.Menu.Set(menu)

// Show about dialog
app.Menu.ShowAbout()
```

### Менеджер сочетаний клавиш

Динамическое управление глобальными сочетаниями клавиш.

```go
// Add key bindings
app.KeyBinding.Add("ctrl+n", func(window application.Window) {
    // Handle Ctrl+N
})

app.KeyBinding.Add("ctrl+q", func(window application.Window) {
    app.Quit()
})

// Remove key bindings
app.KeyBinding.Remove("ctrl+n")

// Get all bindings
bindings := app.KeyBinding.GetAll()
```

### Менеджер контекстных меню

Расширенное управление контекстными меню (для авторов библиотек).

```go
// Create and register context menu
menu := app.ContextMenu.New()
app.ContextMenu.Add("myMenu", menu)

// Retrieve context menu
menu, exists := app.ContextMenu.Get("myMenu")

// Remove context menu
app.ContextMenu.Remove("myMenu")
```

### Менеджер экранов

Управление экранами и преобразование координат в конфигурациях с несколькими мониторами.

```go
// Get screen information
screens := app.Screen.GetAll()
primary := app.Screen.GetPrimary()

// Coordinate transformations
physicalPoint := app.Screen.DipToPhysicalPoint(logicalPoint)
logicalPoint := app.Screen.PhysicalToDipPoint(physicalPoint)

// Screen detection
screen := app.Screen.ScreenNearestDipPoint(point)
screen = app.Screen.ScreenNearestDipRect(rect)
```

### Менеджер буфера обмена

Операции чтения и записи текста в буфере обмена.

```go
// Set text to clipboard
success := app.Clipboard.SetText("Hello World")
if !success {
    // Handle error
}

// Get text from clipboard
text, ok := app.Clipboard.Text()
if !ok {
    // Handle error
} else {
    // Use the text
}
```

### Менеджер SystemTray

Создание значков в области уведомлений и управление ими.

```go
// Create system tray
systray := app.SystemTray.New()
systray.SetLabel("My App")
systray.SetIcon(iconBytes)

// Add menu to system tray
menu := app.Menu.New()
menu.Add("Open").OnClick(func(ctx *Context) {
    // Handle click
})
systray.SetMenu(menu)

// Destroy system tray when done
systray.Destroy()
```

### Менеджер автозапуска

Регистрирует приложение для запуска при входе пользователя в систему. Для каждой платформы выбирается подходящий нативный механизм: SMAppService или plist-файл LaunchAgent в macOS, раздел реестра `HKCU\…\Run` в Windows либо запись XDG `.desktop` в Linux.

```go
// Register to launch at login
err := app.Autostart.Enable()

// With extra launch-time arguments and a custom identifier
err = app.Autostart.EnableWithOptions(application.AutostartOptions{
    Identifier: "com.example.myapp",
    Arguments:  []string{"--hidden"},
})

// Check / remove
enabled, err := app.Autostart.IsEnabled()
status, err := app.Autostart.Status()  // includes Path + Strategy
err = app.Autostart.Disable()
```

Сведения о поведении на разных платформах, правилах для идентификаторов и гарантии обнаружения устаревших записей см. на [странице функции автозапуска](/features/autostart/basics/).
