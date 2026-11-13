---
title: "複数ウィンドウ"
description: "マルチウィンドウアプリケーションのパターンとベストプラクティス"
slug: "features/windows/multiple"
sourcePath: "features/windows/multiple.md"
---

## マルチウィンドウアプリケーション

Wails v3は、設定ウィンドウ、ドキュメントウィンドウ、ツールパレット、インスペクターウィンドウを作成するための**ネイティブなマルチウィンドウサポート**を提供します。シンプルで一貫性のある API を使用して、ウィンドウの追跡、ウィンドウ間の通信、ライフサイクルの管理を行えます。

### メインウィンドウと設定ウィンドウ

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

**要点：**

- メインウィンドウは常に表示する
- 設定ウィンドウは作成しておくが非表示にする
- 必要に応じて設定ウィンドウを表示する
- 同じウィンドウを再利用する（複数作成しない）

## ウィンドウの追跡

### すべてのウィンドウの取得

```go
windows := app.Window.GetAll()
fmt.Printf("Total windows: %d\n", len(windows))

for _, window := range windows {
    fmt.Printf("- %s (ID: %d)\n", window.Name(), window.ID())
}
```

### 特定のウィンドウの検索

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

### ウィンドウレジストリパターン

アプリケーション内のウィンドウを追跡します：

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

## ウィンドウ間の通信

### イベントの使用

ウィンドウはイベントシステムを介して通信します：

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

### 共有状態パターン

共有状態マネージャーを使用します：

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

### ウィンドウ間メッセージ

特定のウィンドウ間でメッセージを送信します：

```go
// Get target window
if targetWindow, ok := app.Window.GetByName("preview"); ok {
    // Emit event to specific window
    targetWindow.EmitEvent("update-preview", previewData)
}
```

## 一般的なパターン

### パターン 1：シングルトンウィンドウ

ウィンドウのインスタンスが 1 つだけ存在するようにします：

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

### パターン 2：ドキュメントウィンドウ

同じ種類のウィンドウのインスタンスを複数作成します：

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

### パターン 3：ツールパレット

常に手前に表示されるフローティングウィンドウです：

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

### パターン 4：モーダルダイアログ（macOS のみ）

親ウィンドウの操作をブロックする子ウィンドウです：

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

### パターン 5：インスペクター／プレビューウィンドウ

連動して更新されるウィンドウです：

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

## 親子関係

### 子ウィンドウの作成

```go
childWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "Child Window",
})

parentWindow.AttachModal(childWindow)
```

**動作：**

- 子ウィンドウは親ウィンドウより手前に表示される
- 子ウィンドウは親ウィンドウとともに移動する
- 子ウィンドウは親ウィンドウに対する操作をブロックする

**対応プラットフォーム：**

| macOS | Windows | Linux |
| --- | --- | --- |
| ✅ | ❌ | ❌ |

### モーダル動作

モーダルに似た動作を実現します：

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

## ウィンドウのライフサイクル管理

### 作成時のコールバック

ウィンドウが作成されたときに通知を受け取ります：

```go
app.Window.OnCreate(func(window application.Window) {
    fmt.Printf("Window created: %s\n", window.Name())

    // Configure all new windows
    window.SetMinSize(400, 300)
})
```

### 破棄時のコールバック

ウィンドウが破棄されたときにクリーンアップします：

```go
window.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
    fmt.Printf("Window %s is closing\n", window.Name())
    
    // Cleanup resources
    cleanup(window.ID())
    
    // Remove from tracking
    removeFromRegistry(window.Name())
})
```

### アプリケーション終了時の動作

アプリケーションを終了するタイミングを制御します：

```go
app := application.New(application.Options{
    Mac: application.MacOptions{
        // Don't quit when last window closes
        ApplicationShouldTerminateAfterLastWindowClosed: false,
    },
})
```

**ユースケース：**

- システムトレイアプリケーション
- バックグラウンドサービス
- メニューバーアプリケーション（macOS）

## メモリ管理

### リークの防止

ウィンドウへの参照は必ずクリーンアップします：

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

### ウィンドウを閉じる

```go
// Close — dispatches WindowClosing; a RegisterHook can call e.Cancel().
window.Close()
```

v3には `window.Destroy()` がありません。ウィンドウが閉じるのを防ぐには、`RegisterHook(events.Common.WindowClosing, ...)`を使用して `event.Cancel()` を呼び出します。

### リソースのクリーンアップ

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

## 高度なパターン

### ウィンドウプール

新しいウィンドウを作成せずに既存のウィンドウを再利用します：

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

### ウィンドウグループ

関連するウィンドウをまとめて管理します：

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

### ワークスペース管理

ウィンドウのレイアウトを保存し、復元します：

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

## 完全な例

以下は、本番環境で使用できるマルチウィンドウアプリケーションです：

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

## ベストプラクティス

### ✅ 推奨事項

- **ウィンドウを追跡する** - 簡単にアクセスできるように参照を保持します
- **破棄時にクリーンアップする** - メモリリークを防止します
- **通信にはイベントを使用する** - 疎結合なアーキテクチャを実現します
- **ウィンドウを再利用する** - 重複して作成しないでください
- **レイアウトを保存／復元する** - UX が向上します
- **ウィンドウを閉じる処理を実装する** - 未保存のデータがある場合は、閉じる前に確認します

### ❌ 非推奨事項

- **ウィンドウを無制限に作成しない** - メモリとパフォーマンスの問題が発生します
- **クリーンアップを忘れない** - メモリリークが発生します
- **グローバル変数を不用意に使用しない** - スレッドセーフティの問題が発生します
- **ウィンドウの作成をブロックしない** - 必要に応じて非同期で作成します
- **プラットフォーム間の違いを無視しない** - すべてのプラットフォームでテストします

## 次のステップ

- [ウィンドウの基本](/features/windows/basics/) - ウィンドウ管理の基礎を学びます
- [ウィンドウイベント](/features/windows/events/) - ウィンドウのライフサイクルイベントを処理します
- [イベントシステム](/features/events/system/) - イベントシステムを詳しく学びます
- [フレームレスウィンドウ](/features/windows/frameless/) - 独自のウィンドウ装飾を作成します

---

**ご質問がありますか？** [Discord](https://discord.gg/JDdSxwjhGf) で質問するか、[マルチウィンドウの例](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/window)を確認してください。
