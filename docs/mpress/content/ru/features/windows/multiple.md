---
title: "Несколько окон"
description: "Шаблоны и рекомендации по разработке многооконных приложений"
slug: "features/windows/multiple"
sourcePath: "features/windows/multiple.md"
---

## Многооконные приложения

Wails v3 предоставляет **встроенную поддержку нескольких окон** для создания окон настроек, окон документов, палитр инструментов и окон инспектора. Отслеживайте окна, обеспечивайте обмен данными между ними и управляйте их жизненным циклом с помощью простых и единообразных API.

### Главное окно и окно настроек

```go
package main

import "github.com/wailsapp/wails/v3/pkg/application"

type App struct {
    app            *application.App
    mainWindow     *application.WebviewWindow
    settingsWindow *application.WebviewWindow
}

func main() {
    app := &App{}
    
    app.app = application.New(application.Options{
        Name: "Multi-Window App",
    })
    
    // Create main window
    app.mainWindow = app.app.Window.NewWithOptions(application.WebviewWindowOptions{
        Name:   "main",
        Title:  "Main Application",
        Width:  1200,
        Height: 800,
    })
    
    // Create settings window (hidden initially)
    app.settingsWindow = app.app.Window.NewWithOptions(application.WebviewWindowOptions{
        Name:   "settings",
        Title:  "Settings",
        Width:  600,
        Height: 400,
        Hidden: true,
    })
    
    app.app.Run()
}

// Show settings from main window
func (a *App) ShowSettings() {
    if a.settingsWindow != nil {
        a.settingsWindow.Show()
        a.settingsWindow.Focus()
    }
}
```

**Ключевые моменты:**

- Главное окно всегда отображается
- Окно настроек создано, но скрыто
- Окно настроек отображается по запросу
- Повторно используйте одно и то же окно (не создавайте несколько окон)

## Отслеживание окон

### Получение всех окон

```go
windows := app.Window.GetAll()
fmt.Printf("Total windows: %d\n", len(windows))

for _, window := range windows {
    fmt.Printf("- %s (ID: %d)\n", window.Name(), window.ID())
}
```

### Поиск определённого окна

```go
// By name
if settings, ok := app.Window.GetByName("settings"); ok {
    settings.Show()
}

// By ID
if window, ok := app.Window.GetByID(123); ok {
    window.Focus()
}

// Current (focused) window
current := app.Window.Current()
```

### Шаблон реестра окон

Отслеживайте окна приложения:

```go
type WindowManager struct {
    windows map[string]*application.WebviewWindow
    mu      sync.RWMutex
}

func (wm *WindowManager) Register(name string, window *application.WebviewWindow) {
    wm.mu.Lock()
    defer wm.mu.Unlock()
    wm.windows[name] = window
}

func (wm *WindowManager) Get(name string) *application.WebviewWindow {
    wm.mu.RLock()
    defer wm.mu.RUnlock()
    return wm.windows[name]
}

func (wm *WindowManager) Remove(name string) {
    wm.mu.Lock()
    defer wm.mu.Unlock()
    delete(wm.windows, name)
}
```

## Обмен данными между окнами

### Использование событий

Окна обмениваются данными через систему событий:

```go
// In main window - emit event
app.Event.Emit("settings-changed", map[string]interface{}{
    "theme": "dark",
    "fontSize": 14,
})

// In settings window - listen for event
app.Event.On("settings-changed", func(event *application.CustomEvent) {
    data := event.Data.(map[string]interface{})
    theme := data["theme"].(string)
    fontSize := data["fontSize"].(int)
    
    // Update UI
    updateSettings(theme, fontSize)
})
```

### Шаблон общего состояния

Используйте общий диспетчер состояния:

```go
type AppState struct {
    theme    string
    fontSize int
    mu       sync.RWMutex
}

var state = &AppState{
    theme:    "light",
    fontSize: 12,
}

func (s *AppState) SetTheme(theme string) {
    s.mu.Lock()
    s.theme = theme
    s.mu.Unlock()
    
    // Notify all windows
    app.Event.Emit("theme-changed", theme)
}

func (s *AppState) GetTheme() string {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.theme
}
```

### Сообщения между окнами

Отправляйте сообщения между определёнными окнами:

```go
// Get target window
if targetWindow, ok := app.Window.GetByName("preview"); ok {
    // Emit event to specific window
    targetWindow.EmitEvent("update-preview", previewData)
}
```

## Распространённые шаблоны

### Шаблон 1: Окна в единственном экземпляре

Обеспечьте существование только одного экземпляра окна:

```go
var settingsWindow *application.WebviewWindow

func ShowSettings(app *application.App) {
    // Create if doesn't exist
    if settingsWindow == nil {
        settingsWindow = app.Window.NewWithOptions(application.WebviewWindowOptions{
            Name:   "settings",
            Title:  "Settings",
            Width:  600,
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

### Шаблон 2: Окна документов

Несколько экземпляров окна одного типа:

```go
type DocumentWindow struct {
    window   *application.WebviewWindow
    filePath string
    modified bool
}

var documents = make(map[string]*DocumentWindow)

func OpenDocument(app *application.App, filePath string) {
    // Check if already open
    if doc, exists := documents[filePath]; exists {
        doc.window.Show()
        doc.window.Focus()
        return
    }
    
    // Create new document window
    window := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:  filepath.Base(filePath),
        Width:  800,
        Height: 600,
    })
    
    doc := &DocumentWindow{
        window:   window,
        filePath: filePath,
        modified: false,
    }
    
    documents[filePath] = doc
    
    // Cleanup on close
    window.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
        delete(documents, filePath)
    })
    
    // Load document
    loadDocument(window, filePath)
}
```

### Шаблон 3: Палитры инструментов

Плавающие окна, которые остаются поверх остальных:

```go
func CreateToolPalette(app *application.App) *application.WebviewWindow {
    palette := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Name:          "tools",
        Title:         "Tools",
        Width:         200,
        Height:        400,
        AlwaysOnTop:   true,
        DisableResize: true,
    })

    return palette
}
```

### Шаблон 4: Модальные диалоговые окна (только macOS)

Дочерние окна, блокирующие родительское окно:

```go
func ShowModaldialog(parent *application.WebviewWindow, title string) {
    dialog := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:         title,
        Width:         400,
        Height:        200,
        AlwaysOnTop:   true,
        DisableResize: true,
    })

    parent.AttachModal(dialog)
}
```

### Шаблон 5: Окна инспектора и предварительного просмотра

Связанные окна, которые обновляются вместе:

```go
type EditorApp struct {
    editor  *application.WebviewWindow
    preview *application.WebviewWindow
}

func (e *EditorApp) UpdatePreview(content string) {
    if e.preview != nil && e.preview.IsVisible() {
        e.preview.EmitEvent("content-changed", content)
    }
}

func (e *EditorApp) TogglePreview() {
    if e.preview == nil {
        e.preview = app.Window.NewWithOptions(application.WebviewWindowOptions{
            Name:   "preview",
            Title:  "Preview",
            Width:  600,
            Height: 800,
        })
        
        e.preview.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
            e.preview = nil
        })
    }
    
    if e.preview.IsVisible() {
        e.preview.Hide()
    } else {
        e.preview.Show()
    }
}
```

## Родительско-дочерние отношения

### Создание дочерних окон

```go
childWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "Child Window",
})

parentWindow.AttachModal(childWindow)
```

**Поведение:**

- Дочернее окно остаётся поверх родительского
- Дочернее окно перемещается вместе с родительским
- Дочернее окно блокирует взаимодействие с родительским

**Поддержка платформ:**

| macOS | Windows | Linux |
| --- | --- | --- |
| ✅ | ❌ | ❌ |

### Модальное поведение

Реализуйте поведение, подобное модальному:

```go
func ShowModaldialog(parent *application.WebviewWindow, title string) {
    dialog := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:       title,
        Width:       400,
        Height:      200,
    })
    
    parent.AttachModal(dialog)
}
```

## Управление жизненным циклом окон

### Обратные вызовы при создании

Получайте уведомления о создании окон:

```go
app.Window.OnCreate(func(window application.Window) {
    fmt.Printf("Window created: %s\n", window.Name())

    // Configure all new windows
    window.SetMinSize(400, 300)
})
```

### Обратные вызовы при уничтожении

Выполняйте очистку при уничтожении окон:

```go
window.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
    fmt.Printf("Window %s is closing\n", window.Name())
    
    // Cleanup resources
    cleanup(window.ID())
    
    // Remove from tracking
    removeFromRegistry(window.Name())
})
```

### Поведение при завершении работы приложения

Управляйте моментом завершения работы приложения:

```go
app := application.New(application.Options{
    Mac: application.MacOptions{
        // Don't quit when last window closes
        ApplicationShouldTerminateAfterLastWindowClosed: false,
    },
})
```

**Варианты использования:**

- Приложения в области уведомлений
- Фоновые службы
- Приложения в строке меню (macOS)

## Управление памятью

### Предотвращение утечек

Всегда удаляйте ссылки на окна:

```go
var windows = make(map[string]*application.WebviewWindow)

func CreateWindow(name string) {
    window := app.Window.New()
    windows[name] = window
    
    // IMPORTANT: Clean up on close
    window.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
        delete(windows, name)
    })
}
```

### Закрытие окон

```go
// Close — dispatches WindowClosing; a RegisterHook can call e.Cancel().
window.Close()
```

В v3 нет `window.Destroy()`. Чтобы предотвратить закрытие, используйте `RegisterHook(events.Common.WindowClosing, ...)` и вызовите `event.Cancel()`.

### Освобождение ресурсов

```go
type ManagedWindow struct {
    window     *application.WebviewWindow
    resources  []io.Closer
}

func (mw *ManagedWindow) Destroy() {
    // Close all resources
    for _, resource := range mw.resources {
        resource.Close()
    }

    // Close the window — WindowClosing listeners run before the window goes away.
    mw.window.Close()
}
```

## Расширенные шаблоны

### Пул окон

Повторно используйте окна вместо создания новых:

```go
type WindowPool struct {
    available []*application.WebviewWindow
    inUse     map[uint]*application.WebviewWindow
    mu        sync.Mutex
}

func (wp *WindowPool) Acquire() *application.WebviewWindow {
    wp.mu.Lock()
    defer wp.mu.Unlock()
    
    // Reuse available window
    if len(wp.available) > 0 {
        window := wp.available[0]
        wp.available = wp.available[1:]
        wp.inUse[window.ID()] = window
        return window
    }
    
    // Create new window
    window := app.Window.New()
    wp.inUse[window.ID()] = window
    return window
}

func (wp *WindowPool) Release(window *application.WebviewWindow) {
    wp.mu.Lock()
    defer wp.mu.Unlock()
    
    delete(wp.inUse, window.ID())
    window.Hide()
    wp.available = append(wp.available, window)
}
```

### Группы окон

Управляйте связанными окнами как единым целым:

```go
type WindowGroup struct {
    name    string
    windows []*application.WebviewWindow
}

func (wg *WindowGroup) Add(window *application.WebviewWindow) {
    wg.windows = append(wg.windows, window)
}

func (wg *WindowGroup) ShowAll() {
    for _, window := range wg.windows {
        window.Show()
    }
}

func (wg *WindowGroup) HideAll() {
    for _, window := range wg.windows {
        window.Hide()
    }
}

func (wg *WindowGroup) CloseAll() {
    for _, window := range wg.windows {
        window.Close()
    }
}
```

### Управление рабочими пространствами

Сохраняйте и восстанавливайте расположение окон:

```go
type WindowLayout struct {
    Windows []WindowState `json:"windows"`
}

type WindowState struct {
    Name   string `json:"name"`
    X      int    `json:"x"`
    Y      int    `json:"y"`
    Width  int    `json:"width"`
    Height int    `json:"height"`
}

func SaveLayout() *WindowLayout {
    layout := &WindowLayout{}
    
    for _, window := range app.Window.GetAll() {
        x, y := window.Position()
        width, height := window.Size()
        
        layout.Windows = append(layout.Windows, WindowState{
            Name:   window.Name(),
            X:      x,
            Y:      y,
            Width:  width,
            Height: height,
        })
    }
    
    return layout
}

func RestoreLayout(layout *WindowLayout) {
    for _, state := range layout.Windows {
        if window, ok := app.Window.GetByName(state.Name); ok {
            window.SetPosition(state.X, state.Y)
            window.SetSize(state.Width, state.Height)
        }
    }
}
```

## Полный пример

Ниже приведён готовый к промышленному использованию пример многооконного приложения:

```go
package main

import (
    "encoding/json"
    "os"
    "sync"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)

type MultiWindowApp struct {
    app     *application.App
    windows map[string]*application.WebviewWindow
    mu      sync.RWMutex
}

func main() {
    mwa := &MultiWindowApp{
        windows: make(map[string]*application.WebviewWindow),
    }
    
    mwa.app = application.New(application.Options{
        Name: "Multi-Window Application",
        Mac: application.MacOptions{
            ApplicationShouldTerminateAfterLastWindowClosed: false,
        },
    })
    
    // Create main window
    mwa.CreateMainWindow()
    
    // Load saved layout
    mwa.LoadLayout()
    
    mwa.app.Run()
}

func (mwa *MultiWindowApp) CreateMainWindow() {
    window := mwa.app.Window.NewWithOptions(application.WebviewWindowOptions{
        Name:   "main",
        Title:  "Main Application",
        Width:  1200,
        Height: 800,
    })
    
    mwa.RegisterWindow("main", window)
}

func (mwa *MultiWindowApp) ShowSettings() {
    if window := mwa.GetWindow("settings"); window != nil {
        window.Show()
        window.Focus()
        return
    }
    
    window := mwa.app.Window.NewWithOptions(application.WebviewWindowOptions{
        Name:   "settings",
        Title:  "Settings",
        Width:  600,
        Height: 400,
    })
    
    mwa.RegisterWindow("settings", window)
}

func (mwa *MultiWindowApp) OpenDocument(path string) {
    name := "doc-" + path
    
    if window := mwa.GetWindow(name); window != nil {
        window.Show()
        window.Focus()
        return
    }
    
    window := mwa.app.Window.NewWithOptions(application.WebviewWindowOptions{
        Name:  name,
        Title: path,
        Width: 800,
        Height: 600,
    })
    
    mwa.RegisterWindow(name, window)
}

func (mwa *MultiWindowApp) RegisterWindow(name string, window *application.WebviewWindow) {
    mwa.mu.Lock()
    mwa.windows[name] = window
    mwa.mu.Unlock()
    
    window.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
        mwa.UnregisterWindow(name)
    })
}

func (mwa *MultiWindowApp) UnregisterWindow(name string) {
    mwa.mu.Lock()
    delete(mwa.windows, name)
    mwa.mu.Unlock()
}

func (mwa *MultiWindowApp) GetWindow(name string) *application.WebviewWindow {
    mwa.mu.RLock()
    defer mwa.mu.RUnlock()
    return mwa.windows[name]
}

func (mwa *MultiWindowApp) SaveLayout() {
    layout := make(map[string]WindowState)
    
    mwa.mu.RLock()
    for name, window := range mwa.windows {
        x, y := window.Position()
        width, height := window.Size()
        
        layout[name] = WindowState{
            X:      x,
            Y:      y,
            Width:  width,
            Height: height,
        }
    }
    mwa.mu.RUnlock()
    
    data, _ := json.Marshal(layout)
    os.WriteFile("layout.json", data, 0644)
}

func (mwa *MultiWindowApp) LoadLayout() {
    data, err := os.ReadFile("layout.json")
    if err != nil {
        return
    }
    
    var layout map[string]WindowState
    if err := json.Unmarshal(data, &layout); err != nil {
        return
    }
    
    for name, state := range layout {
        if window := mwa.GetWindow(name); window != nil {
            window.SetPosition(state.X, state.Y)
            window.SetSize(state.Width, state.Height)
        }
    }
}

type WindowState struct {
    X      int `json:"x"`
    Y      int `json:"y"`
    Width  int `json:"width"`
    Height int `json:"height"`
}
```

## Рекомендации

### ✅ Рекомендуется

- **Отслеживайте окна** — сохраняйте ссылки на них для быстрого доступа
- **Освобождайте ресурсы при уничтожении** — предотвращайте утечки памяти
- **Используйте события для обмена данными** — это обеспечивает слабую связанность архитектуры
- **Повторно используйте окна** — не создавайте дубликаты
- **Сохраняйте и восстанавливайте расположение окон** — это улучшает удобство работы
- **Обрабатывайте закрытие окна** — при наличии несохранённых данных запрашивайте подтверждение перед закрытием

### ❌ Не рекомендуется

- **Не создавайте неограниченное количество окон** — это приводит к проблемам с памятью и производительностью
- **Не забывайте освобождать ресурсы** — иначе возникают утечки памяти
- **Не используйте глобальные переменные без должной осторожности** — это может привести к проблемам с потокобезопасностью
- **Не блокируйте создание окон** — при необходимости создавайте их асинхронно
- **Не игнорируйте различия между платформами** — тестируйте приложение на всех платформах

## Дальнейшие шаги

- [Основы работы с окнами](/features/windows/basics/) — изучите основы управления окнами
- [События окон](/features/windows/events/) — обрабатывайте события жизненного цикла окон
- [Система событий](/features/events/system/) — подробно изучите систему событий
- [Окна без рамки](/features/windows/frameless/) — создавайте собственное оформление окон

---

**Есть вопросы?** Задайте их в [Discord](https://discord.gg/JDdSxwjhGf) или ознакомьтесь с [примером многооконного приложения](https://github.com/wailsapp/wails/tree/master/v3/examples/multi-window).
