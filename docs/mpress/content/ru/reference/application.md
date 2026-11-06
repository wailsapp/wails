---
title: "API приложения"
description: "Полное справочное руководство по API приложения"
slug: "reference/application"
sourcePath: "reference/application.md"
---

## Обзор

`Application` — это ядро вашего приложения Wails. Оно управляет окнами, службами и событиями, а также предоставляет доступ ко всем возможностям платформы.

## Создание приложения

```go
import "github.com/wailsapp/wails/v3/pkg/application"

app := application.New(application.Options{
    Name:        "My App",
    Description: "My awesome application",
    Services: []application.Service{
        application.NewService(&MyService{}),
    },
})
```

## Основные методы

### Run()

Запускает цикл обработки событий приложения.

```go
func (a *App) Run() error
```

**Пример:**

```go
err := app.Run()
if err != nil {
    log.Fatal(err)
}
```

**Возвращаемое значение:** ошибка, если запуск завершился неудачей

### Quit()

Корректно завершает работу приложения.

```go
func (a *App) Quit()
```

**Пример:**

```go
// In a menu handler
menu.Add("Quit").OnClick(func(ctx *application.Context) {
    app.Quit()
})
```

### Config()

Возвращает конфигурацию приложения.

```go
func (a *App) Config() Options
```

**Пример:**

```go
config := app.Config()
fmt.Println("App name:", config.Name)
```

## Управление окнами

### app.Window.New()

Создаёт новое окно webview с параметрами по умолчанию.

```go
func (wm *WindowManager) New() *WebviewWindow
```

**Пример:**

```go
window := app.Window.New()
window.Show()
```

### app.Window.NewWithOptions()

Создаёт новое окно webview с пользовательскими параметрами.

```go
func (wm *WindowManager) NewWithOptions(options WebviewWindowOptions) *WebviewWindow
```

**Пример:**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "My Window",
    Width:  800,
    Height: 600,
    BackgroundColour: application.NewRGB(255, 255, 255),
})
```

### app.Window.GetByName()

Получает окно по его имени. Возвращает окно и признак того, было ли оно найдено.

```go
func (wm *WindowManager) GetByName(name string) (Window, bool)
```

**Пример:**

```go
if window, ok := app.Window.GetByName("main"); ok {
    window.Show()
}
```

### app.Window.GetAll()

Возвращает все окна приложения.

```go
func (wm *WindowManager) GetAll() []Window
```

**Пример:**

```go
windows := app.Window.GetAll()
for _, window := range windows {
    fmt.Println("Window:", window.Name())
}
```

## Менеджеры

Application предоставляет через свойства доступ к различным менеджерам:

```go
app.Window       // Window management
app.Menu         // Menu management
app.Dialog       // Dialog management
app.Event        // Event management
app.Clipboard    // Clipboard operations
app.Screen       // Screen information
app.SystemTray   // System tray
app.Browser      // Browser operations
app.Env          // Environment variables
app.ContextMenu  // Context-menu management
app.KeyBinding   // Global keyboard shortcuts
app.Logger       // *slog.Logger
```

### Пример использования

```go
// Create window
window := app.Window.New()

// Show dialog
app.Dialog.Info().SetMessage("Hello!").Show()

// Copy to clipboard
app.Clipboard.SetText("Copied text")

// Get screens
screens := app.Screen.GetAll()
```

## Управление службами

### RegisterService()

Регистрирует службу в приложении.

```go
func (a *App) RegisterService(service Service)
```

`RegisterService` ничего не возвращает; ошибки инициализации службы проявляются как сбои `ServiceStartup` во время `app.Run()`.

**Пример:**

```go
type MyService struct {
    app *application.App
}

func NewMyService(app *application.App) *MyService {
    return &MyService{app: app}
}

// Register after app creation
app.RegisterService(application.NewService(NewMyService(app)))
```

## Управление событиями

### app.Event.Emit()

Генерирует пользовательское событие. Возвращает `true`, если перехватчик отменил генерацию события.

```go
func (em *EventManager) Emit(name string, data ...any) bool
```

**Пример:**

```go
// Emit event with data
app.Event.Emit("user-logged-in", map[string]interface{}{
    "username": "john",
    "timestamp": time.Now(),
})
```

### app.Event.On()

Отслеживает пользовательские события. Возвращает `func()` для отмены подписки.

```go
func (em *EventManager) On(name string, callback func(*CustomEvent)) func()
```

**Пример:**

```go
app.Event.On("user-logged-in", func(e *application.CustomEvent) {
    data := e.Data.(map[string]interface{})
    username := data["username"].(string)
    fmt.Println("User logged in:", username)
})
```

### app.Event.OnApplicationEvent()

Отслеживает события жизненного цикла приложения. Параметр `eventType` имеет тип `events.ApplicationEventType` (из пакета `events`).

```go
func (em *EventManager) OnApplicationEvent(
    eventType events.ApplicationEventType,
    callback func(*ApplicationEvent),
) func()
```

**Пример:**

```go
import "github.com/wailsapp/wails/v3/pkg/events"

// Listen for app-started
app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(e *application.ApplicationEvent) {
    fmt.Println("Application started")
})

// Application shutdown is NOT an event constant; register cleanup via:
app.OnShutdown(func() {
    fmt.Println("Application shutting down")
})
```

## Методы диалоговых окон

Доступ к диалоговым окнам осуществляется через менеджер `app.Dialog`. Полное справочное руководство см. в разделе [API диалоговых окон](/reference/dialogs/).

### Диалоговые окна сообщений

```go
// Information dialog
app.Dialog.Info().
    SetTitle("Success").
    SetMessage("Operation completed!").
    Show()

// Error dialog
app.Dialog.Error().
    SetTitle("Error").
    SetMessage("Something went wrong.").
    Show()

// Warning dialog
app.Dialog.Warning().
    SetTitle("Warning").
    SetMessage("This action cannot be undone.").
    Show()
```

### Диалоговые окна вопросов

Для обработки ответов пользователя диалоговые окна вопросов используют функции обратного вызова кнопок:

```go
dialog := app.Dialog.Question().
    SetTitle("Confirm").
    SetMessage("Continue?")

yes := dialog.AddButton("Yes")
yes.OnClick(func() {
    // Handle yes
})

no := dialog.AddButton("No")
no.OnClick(func() {
    // Handle no
})

dialog.SetDefaultButton(yes)
dialog.SetCancelButton(no)
dialog.Show()
```

### Диалоговые окна выбора файлов

```go
// Open file dialog
path, err := app.Dialog.OpenFile().
    SetTitle("Select File").
    AddFilter("Images", "*.png;*.jpg").
    PromptForSingleSelection()

// Save file dialog
path, err := app.Dialog.SaveFile().
    SetTitle("Save File").
    SetFilename("document.pdf").
    AddFilter("PDF", "*.pdf").
    PromptForSingleSelection()

// Folder selection (use OpenFile with directory options)
path, err := app.Dialog.OpenFile().
    SetTitle("Select Folder").
    CanChooseDirectories(true).
    CanChooseFiles(false).
    PromptForSingleSelection()
```

## Журналирование

Приложение предоставляет логгер для структурированного журналирования:

```go
app.Logger.Info("Message", "key", "value")
app.Logger.Error("Error occurred", "error", err)
app.Logger.Debug("Debug info")
app.Logger.Warn("Warning message")
```

**Пример:**

```go
func (s *MyService) ProcessData(data string) error {
    s.app.Logger.Info("Processing data", "length", len(data))
    
    if err := process(data); err != nil {
        s.app.Logger.Error("Processing failed", "error", err)
        return err
    }
    
    s.app.Logger.Info("Processing complete")
    return nil
}
```

## Обработка необработанных сообщений

Для приложений, которым требуется прямое низкоуровневое управление обменом данными между фронтендом и бэкендом, Wails предоставляет параметр `RawMessageHandler`. Он позволяет обойти стандартную систему привязок.

@note{type="info"}
Необработанные сообщения следует использовать только в крайнем случае. Стандартная система привязок хорошо оптимизирована и подходит почти для всех приложений. Используйте необработанные сообщения, только если вы выполнили профилирование приложения и подтвердили, что привязки являются узким местом.

@end

### RawMessageHandler

`RawMessageHandler` — это поле структуры `application.Options`, а не метод. Среда выполнения вызывает его для каждого необработанного сообщения, отправленного из фронтенда через `System.invoke()`.

```go
type Options struct {
    // ... other fields ...
    RawMessageHandler func(window Window, message string, originInfo *OriginInfo)
}
```

`OriginInfo` содержит `Origin`, `TopOrigin` и `IsMainFrame` (на разных платформах заполняются разные подмножества полей — матрицу по платформам см. в [руководстве по необработанным сообщениям](/guides/raw-messages/)).

**Пример:**

```go
app := application.New(application.Options{
    Name: "My App",
    RawMessageHandler: func(window application.Window, message string, originInfo *application.OriginInfo) {
        // Handle the raw message
        fmt.Printf("Received from %s (%s): %s\n", window.Name(), originInfo.Origin, message)

        // You can respond using events
        window.EmitEvent("response", processMessage(message))
    },
})
```

Дополнительные сведения см. в [руководстве по необработанным сообщениям](/guides/raw-messages/).

## Параметры для отдельных платформ

### Параметры Windows

Настройте поведение для Windows на уровне приложения:

```go
app := application.New(application.Options{
    Name: "My App",
    Windows: application.WindowsOptions{
        // WebView2 browser flags (apply to ALL windows)
        EnabledFeatures:       []string{"msWebView2EnableDraggableRegions"},
        DisabledFeatures:      []string{"msExperimentalFeature"},
        AdditionalBrowserArgs: []string{"--remote-debugging-port=9222"},

        // Other Windows options
        WndClass:                      "MyAppClass",
        WebviewUserDataPath:           "",  // Default: %APPDATA%\[BinaryName.exe]
        WebviewBrowserPath:            "",  // Default: system WebView2
        DisableQuitOnLastWindowClosed: false,
    },
})
```

**Флаги браузера:**

- `EnabledFeatures` — включаемые флаги функций WebView2
- `DisabledFeatures` — отключаемые флаги функций WebView2
- `AdditionalBrowserArgs` — аргументы командной строки Chromium

Подробную документацию см. в разделе [«Параметры окна — параметры Windows уровня приложения»](/features/windows/options/#application-level-windows-options).

### Параметры macOS

```go
app := application.New(application.Options{
    Name: "My App",
    Mac: application.MacOptions{
        ActivationPolicy: application.ActivationPolicyRegular,
        ApplicationShouldTerminateAfterLastWindowClosed: true,
    },
})
```

### Параметры Linux

```go
app := application.New(application.Options{
    Name: "My App",
    Linux: application.LinuxOptions{
        ProgramName:                   "my-app",
        DisableQuitOnLastWindowClosed: false,
    },
})
```

## Полный пример приложения

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name:        "My Application",
        Description: "A demo application",
        Mac: application.MacOptions{
            ApplicationShouldTerminateAfterLastWindowClosed: true,
        },
    })

    // Create main window
    window := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:            "My App",
        Width:            1024,
        Height:           768,
        MinWidth:         800,
        MinHeight:        600,
        BackgroundColour: application.NewRGB(255, 255, 255),
        URL:              "http://wails.localhost/",
    })

    window.Center()
    window.Show()

    app.Run()
}
```
