---
title: "マネージャー API"
description: "目的別のマネージャーインターフェースで整理された API 構造"
slug: "concepts/manager-api"
sourcePath: "concepts/manager-api.md"
---

Wails v3 のマネージャー API では、`*application.App` の公開フィールド下に目的別にまとめられたマネージャー構造体を通じて、アプリケーションの機能に整理された分かりやすい方法でアクセスできます。Wails 3 は v2 との互換性を完全に断っており、従来の `app.NewWebviewWindow(...)` 形式の API との互換性を維持するための呼び出し単位のラッパーレイヤーはありません。そのため、アプリの操作には以下のマネージャーを使用します。

## 概要

マネージャー API は、アプリケーションの機能を目的別の 12 の領域（1 つのロガーと 11 のマネージャー）に整理します。

- **`app.Window`** - ウィンドウの作成、管理、コールバック
- **`app.ContextMenu`** - コンテキストメニューの登録と管理\
- **`app.KeyBinding`** - グローバルキーバインドの管理
- **`app.Browser`** - ブラウザー連携（URL とファイルを開く）
- **`app.Env`** - 環境情報とシステム状態
- **`app.Dialog`** - ファイルダイアログとメッセージダイアログの操作
- **`app.Event`** - カスタムイベントの処理とアプリケーションイベント
- **`app.Menu`** - アプリケーションメニューの管理
- **`app.Screen`** - 画面の管理と座標変換
- **`app.Clipboard`** - クリップボードのテキスト操作
- **`app.SystemTray`** - システムトレイアイコンの作成と管理
- **`app.Autostart`** - ユーザーログイン時にアプリケーションを起動するよう登録

## 利点

- **見つけやすさの向上** - IDE の自動補完に整理された API が表示されます
- **コード構成の改善** - 関連するメソッドがまとめられています
- **保守性の向上** - マネージャー間で関心事が分離されています
- **将来の拡張性** - 特定の領域に新機能を追加しやすくなります

## 使用方法

マネージャー API を使用すると、アプリケーションのすべての機能に整理された形でアクセスできます。

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

## マネージャーリファレンス

### ウィンドウマネージャー

ウィンドウの作成、取得、ライフサイクルコールバックを管理します。

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

### イベントマネージャー

カスタムイベントとアプリケーションイベントのリスニングを処理します。

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

### ブラウザーマネージャー

URL やファイルを開くためのブラウザー連携機能を提供します。

```go
// Open URLs and files in default browser
err := app.Browser.OpenURL("https://wails.io")
err := app.Browser.OpenFile("/path/to/document.pdf")
```

### 環境マネージャー

システム環境情報にアクセスできます。

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

### ダイアログマネージャー

ファイルダイアログとメッセージダイアログに整理された形でアクセスできます。

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

### メニューマネージャー

アプリケーションメニューを作成および管理します。

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

### キーバインドマネージャー

グローバルキーバインドを動的に管理します。

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

### コンテキストメニューマネージャー

高度なコンテキストメニュー管理機能です（ライブラリ作者向け）。

```go
// Create and register context menu
menu := app.ContextMenu.New()
app.ContextMenu.Add("myMenu", menu)

// Retrieve context menu
menu, exists := app.ContextMenu.Get("myMenu")

// Remove context menu
app.ContextMenu.Remove("myMenu")
```

### 画面マネージャー

マルチモニター構成における画面の管理と座標変換を行います。

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

### クリップボードマネージャー

クリップボードでテキストを読み書きするための操作を提供します。

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

### SystemTray マネージャー

システムトレイアイコンを作成および管理します。

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

### 自動起動マネージャー

ユーザーログイン時にアプリケーションを起動するよう登録します。プラットフォームごとに適切なネイティブ機構を選択します。macOS では SMAppService または LaunchAgent plist、Windows では `HKCU\…\Run` レジストリキー、Linux では XDG `.desktop` エントリを使用します。

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

プラットフォームごとの動作、識別子の規則、古い登録を検出する保証については、[自動起動機能のページ](/features/autostart/basics/)を参照してください。
