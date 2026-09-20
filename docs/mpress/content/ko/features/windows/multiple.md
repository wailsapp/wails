---
title: "다중 창"
description: "다중 창 애플리케이션을 위한 패턴과 모범 사례"
slug: "features/windows/multiple"
sourcePath: "features/windows/multiple.md"
---

## 다중 창 애플리케이션

Wails v3는 설정 창, 문서 창, 도구 팔레트 및 검사기 창을 만들 수 있도록 <strong>네이티브 다중 창 지원</strong>을 제공합니다. 간단하고 일관된 API로 창을 추적하고, 창 간 통신을 지원하며, 창의 수명 주기를 관리할 수 있습니다.

### 메인 창 + 설정 창

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

**핵심 사항:**

- 메인 창은 항상 표시
- 설정 창은 생성하되 숨김
- 필요할 때 설정 창 표시
- 같은 창을 재사용(여러 개를 만들지 않음)

## 창 추적

### 모든 창 가져오기

```go
windows := app.Window.GetAll()
fmt.Printf("Total windows: %d\n", len(windows))

for _, window := range windows {
    fmt.Printf("- %s (ID: %d)\n", window.Name(), window.ID())
}
```

### 특정 창 찾기

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

### 창 레지스트리 패턴

애플리케이션의 창을 추적합니다:

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

## 창 간 통신

### 이벤트 사용

창은 이벤트 시스템을 통해 통신합니다:

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

### 공유 상태 패턴

공유 상태 관리자를 사용합니다:

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

### 창 간 메시지

특정 창 간에 메시지를 전송합니다:

```go
// Get target window
if targetWindow, ok := app.Window.GetByName("preview"); ok {
    // Emit event to specific window
    targetWindow.EmitEvent("update-preview", previewData)
}
```

## 일반적인 패턴

### 패턴 1: 싱글턴 창

창 인스턴스가 하나만 존재하도록 합니다:

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

### 패턴 2: 문서 창

동일한 창 유형의 인스턴스를 여러 개 사용합니다:

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

### 패턴 3: 도구 팔레트

항상 위에 표시되는 부동 창입니다:

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

### 패턴 4: 모달 대화 상자(macOS 전용)

부모 창을 차단하는 자식 창입니다:

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

### 패턴 5: 검사기/미리 보기 창

함께 업데이트되는 연결된 창입니다:

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

## 부모-자식 관계

### 자식 창 만들기

```go
childWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "Child Window",
})

parentWindow.AttachModal(childWindow)
```

**동작:**

- 자식 창이 부모 창 위에 유지됨
- 자식 창이 부모 창과 함께 이동함
- 자식 창이 부모 창과의 상호 작용을 차단함

**플랫폼 지원:**

| macOS | Windows | Linux |
| --- | --- | --- |
| ✅ | ❌ | ❌ |

### 모달 동작

모달과 유사한 동작을 구현합니다:

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

## 창 수명 주기 관리

### 생성 콜백

창이 생성될 때 알림을 받습니다:

```go
app.Window.OnCreate(func(window application.Window) {
    fmt.Printf("Window created: %s\n", window.Name())

    // Configure all new windows
    window.SetMinSize(400, 300)
})
```

### 소멸 콜백

창이 소멸될 때 정리합니다:

```go
window.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
    fmt.Printf("Window %s is closing\n", window.Name())
    
    // Cleanup resources
    cleanup(window.ID())
    
    // Remove from tracking
    removeFromRegistry(window.Name())
})
```

### 애플리케이션 종료 동작

애플리케이션이 종료되는 시점을 제어합니다:

```go
app := application.New(application.Options{
    Mac: application.MacOptions{
        // Don't quit when last window closes
        ApplicationShouldTerminateAfterLastWindowClosed: false,
    },
})
```

**사용 사례:**

- 시스템 트레이 애플리케이션
- 백그라운드 서비스
- 메뉴 막대 애플리케이션(macOS)

## 메모리 관리

### 누수 방지

창 참조를 항상 정리합니다:

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

### 창 닫기

```go
// Close — dispatches WindowClosing; a RegisterHook can call e.Cancel().
window.Close()
```

v3에는 `window.Destroy()`이 없습니다. 창이 닫히지 않도록 하려면 `RegisterHook(events.Common.WindowClosing, ...)`을 사용하고 `event.Cancel()`을 호출합니다.

### 리소스 정리

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

## 고급 패턴

### 창 풀

새 창을 만드는 대신 기존 창을 재사용합니다:

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

### 창 그룹

관련 창을 함께 관리합니다:

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

### 워크스페이스 관리

창 레이아웃을 저장하고 복원합니다:

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

## 전체 예제

다음은 프로덕션 환경에서 바로 사용할 수 있는 다중 창 애플리케이션입니다:

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

## 모범 사례

### ✅ 권장 사항

- **창 추적하기** - 쉽게 접근할 수 있도록 참조를 유지하세요
- **소멸 시 정리하기** - 메모리 누수를 방지하세요
- **통신에 이벤트 사용하기** - 결합도가 낮은 아키텍처를 구성하세요
- **창 재사용하기** - 중복 창을 만들지 마세요
- **레이아웃 저장 및 복원하기** - 더 나은 사용자 경험을 제공하세요
- **창 닫기 처리하기** - 저장하지 않은 데이터가 있으면 닫기 전에 확인하세요

### ❌ 금지 사항

- **창을 무제한으로 만들지 않기** - 메모리 및 성능 문제가 발생합니다
- **정리를 잊지 않기** - 메모리 누수가 발생합니다
- **전역 변수를 부주의하게 사용하지 않기** - 스레드 안전성 문제가 발생합니다
- **창 생성 시 실행을 블로킹하지 마세요** - 필요한 경우 비동기 방식으로 생성하세요
- **플랫폼 차이를 무시하지 않기** - 모든 플랫폼에서 테스트하세요

## 다음 단계

- [창 기초](/features/windows/basics/) - 창 관리의 기본 사항을 알아보세요
- [창 이벤트](/features/windows/events/) - 창 수명 주기 이벤트를 처리하세요
- [이벤트 시스템](/features/events/system/) - 이벤트 시스템을 자세히 알아보세요
- [프레임리스 창](/features/windows/frameless/) - 사용자 지정 창 프레임을 만드세요

---

**궁금한 점이 있으신가요?** [Discord](https://discord.gg/JDdSxwjhGf)에서 질문하거나 [다중 창 예제](https://github.com/wailsapp/wails/tree/master/v3/examples/multi-window)를 확인하세요.
