---
title: "ウィンドウの基本"
description: "Wailsでアプリケーションウィンドウを作成、管理する方法"
slug: "features/windows/basics"
sourcePath: "features/windows/basics.md"
---

## ウィンドウ管理

Wailsは、すべてのプラットフォームで動作する**統一ウィンドウ管理API**を提供します。ウィンドウを作成して動作を制御できるほか、作成、外観、動作、ライフサイクルを完全に制御しながら複数のウィンドウを管理できます。

## クイックスタート

```go
package main

import "github.com/wailsapp/wails/v3/pkg/application"

func main() {
    app := application.New(application.Options{
        Name: "My App",
    })

    // Create a window
    window := app.Window.New()
    
    // Configure it
    window.SetTitle("Hello Wails")
    window.SetSize(800, 600)
    window.Center()
    
    // Show it
    window.Show()

    app.Run()
}
```

<strong>これだけです！</strong>クロスプラットフォームのウィンドウが作成できました。

## ウィンドウの作成

### 基本的なウィンドウ

ウィンドウを作成する最も簡単な方法は次のとおりです。

```go
window := app.Window.New()
```

**作成されるもの：**

- デフォルトサイズ（800x600）
- デフォルトのタイトル（アプリケーション名）
- フロントエンドで使用できるWebView
- プラットフォームネイティブの外観

### オプションを指定したウィンドウ

カスタム設定を指定してウィンドウを作成します。

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "My Application",
    Width:  1200,
    Height: 800,
    X:      100,   // Position from left
    Y:      100,   // Position from top
    AlwaysOnTop: false,
    Frameless: false,
    Hidden: false,
    MinWidth: 400,
    MinHeight: 300,
    MaxWidth: 1920,
    MaxHeight: 1080,
})
```

**よく使うオプション：**

| オプション | 型 | 説明 |
| --- | --- | --- |
| `Title` | `string` | ウィンドウのタイトル |
| `Width` | `int` | ウィンドウの幅（ピクセル） |
| `Height` | `int` | ウィンドウの高さ（ピクセル） |
| `X` | `int` | X座標（左端から） |
| `Y` | `int` | Y座標（上端から） |
| `AlwaysOnTop` | `bool` | ウィンドウをほかのウィンドウより前面に維持 |
| `Frameless` | `bool` | タイトルバーと境界線を削除 |
| `Hidden` | `bool` | 非表示で起動 |
| `MinWidth` | `int` | 最小幅 |
| `MinHeight` | `int` | 最小高さ |
| `MaxWidth` | `int` | 最大幅 |
| `MaxHeight` | `int` | 最大高さ |

**完全な一覧については、[ウィンドウオプション](/features/windows/options/)を参照してください。**

### 名前付きウィンドウ

簡単に取得できるよう、ウィンドウに名前を付けます。

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:  "main-window",
    Title: "Main Application",
})

// Later, find it by name
if mainWindow, ok := app.Window.GetByName("main-window"); ok {
    mainWindow.Show()
}
```

**ユースケース：**

- 複数のウィンドウ（メイン、設定、バージョン情報）
- コード内の異なる場所からウィンドウを検索
- ウィンドウ間の通信

## ウィンドウの制御

### 表示と非表示

```go
// Show window
window.Show()

// Hide window
window.Hide()

// Check if visible
if window.IsVisible() {
    fmt.Println("Window is visible")
}
```

**ユースケース：**

- スプラッシュ画面（表示後に非表示）
- 設定ウィンドウ（不要なときは非表示）
- ポップアップウィンドウ（必要に応じて表示）

### 位置とサイズ

```go
// Set size
window.SetSize(1024, 768)

// Set position
window.SetPosition(100, 100)

// Centre on screen
window.Center()

// Get current size
width, height := window.Size()

// Get current position
x, y := window.Position()
```

**座標系：**

- （0, 0）はプライマリ画面の左上
- Xの正方向は右
- Yの正方向は下

### ウィンドウの状態

```go
// Minimise
window.Minimise()

// Maximise
window.Maximise()

// Fullscreen
window.Fullscreen()

// Restore to normal
window.Restore()

// Check state
if window.IsMinimised() {
    fmt.Println("Window is minimised")
}

if window.IsMaximised() {
    fmt.Println("Window is maximised")
}

if window.IsFullscreen() {
    fmt.Println("Window is fullscreen")
}
```

**状態遷移：**

```
Normal ←→ Minimised
Normal ←→ Maximised
Normal ←→ Fullscreen
```

### タイトルと外観

```go
// Set title
window.SetTitle("My Application - Document.txt")

// Set background colour — RGBA value (helper for RGB)
window.SetBackgroundColour(application.NewRGBA(0, 0, 0, 255))

// Set always on top
window.SetAlwaysOnTop(true)

// Set resizable
window.SetResizable(false)
```

### ウィンドウを閉じる

```go
// Close window — dispatches WindowClosing; a RegisterHook can call e.Cancel().
window.Close()
```

v3には`window.Destroy()`メソッドがありません。`Close()`を使用し、`OnWindowEvent`でリッスンするか（キャンセル不可）、`RegisterHook`でフックします（`e.Cancel()`を呼び出してウィンドウを開いたままにできます）。

## ウィンドウの検索

### 名前で検索

```go
if window, ok := app.Window.GetByName("settings"); ok {
    window.Show()
}
```

### IDで検索

各ウィンドウには一意のIDがあります。

```go
id := window.ID()
fmt.Printf("Window ID: %d\n", id)

// Find by ID
if found, ok := app.Window.GetByID(id); ok {
    found.Focus()
}
```

### 現在のウィンドウ

現在フォーカスされているウィンドウを取得します。

```go
current := app.Window.Current()
if current != nil {
    current.SetTitle("Active Window")
}
```

### すべてのウィンドウ

すべてのウィンドウを取得します。

```go
windows := app.Window.GetAll()
fmt.Printf("Total windows: %d\n", len(windows))

for _, w := range windows {
    fmt.Printf("Window: %s (ID: %d)\n", w.Name(), w.ID())
}
```

## ウィンドウのライフサイクル

### 作成

```go
app.Window.OnCreate(func(window application.Window) {
    fmt.Printf("Window created: %s\n", window.Name())

    // Configure new windows
    window.SetMinSize(400, 300)
})
```

### 終了

ウィンドウが閉じるのを防ぐには、`WindowClosing`イベントで`RegisterHook`を使用します。

```go
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

**重要：**`RegisterHook`は、クローズイベントが発生する前にインターセプトします。ウィンドウが閉じるのを防ぐには、`event.Cancel()`を呼び出します。これは、ユーザーが開始した終了操作（Xボタンのクリック）で機能します。

### 破棄

ウィンドウが閉じるときにクリーンアップを実行するには、`WindowClosing`イベントで`OnWindowEvent`を使用します。

```go
window.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
    fmt.Println("Window is closing")
    // Cleanup resources
})
```

## 複数のウィンドウ

### 複数のウィンドウの作成

```go
// Main window
mainWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:   "main",
    Title:  "Main Application",
    Width:  1200,
    Height: 800,
})

// Settings window
settingsWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:   "settings",
    Title:  "Settings",
    Width:  600,
    Height: 400,
    Hidden: true,  // Start hidden
})

// Show settings when needed
settingsWindow.Show()
```

### ウィンドウ間の通信

ウィンドウはイベントを介して通信できます。

```go
// In main window
app.Event.Emit("data-updated", map[string]interface{}{
    "value": 42,
})

// In settings window
app.Event.On("data-updated", func(event *application.CustomEvent) {
    data := event.Data.(map[string]interface{})
    value := data["value"].(int)
    fmt.Printf("Received: %d\n", value)
})
```

**詳しくは、[イベント](/features/events/system/)を参照してください。**

### 親子ウィンドウ

`WebviewWindowOptions`には`Parent`フィールドがありません。子ウィンドウを通常のウィンドウとして作成し、シートモーダルとして親ウィンドウに関連付けます。

```go
// Create child window
childWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Child Window",
})

// Attach to the parent — presents as a sheet on macOS.
mainWindow.AttachModal(childWindow)
```

**動作：**

- 子ウィンドウは親ウィンドウより前面に維持されます。
- 子ウィンドウはモーダルであり、親ウィンドウの操作をブロックします。

**プラットフォームのサポート状況：**

- <strong>macOS：</strong>完全にサポートされています（シートとして表示されます）。
- <strong>Windows：</strong>サポートされていません。
- <strong>Linux：</strong>サポートされていません。

## プラットフォーム固有の機能

@tabs{sync-key="platform"}
[Windows]
**Windows 固有の機能：**

```go
// Flash taskbar button
window.Flash(true)  // Start flashing
window.Flash(false) // Stop flashing

// Trigger Windows 11 Snap Assist (Win+Z)
window.SnapAssist()
```

ウィンドウごとの `SetIcon` はありません。アプリケーションアイコンは `app.SetIcon([]byte)` を使用してアプリに設定します（Linux 固有のウィンドウアイコンについては、ウィンドウ作成時に `application.LinuxWindow.Icon` フィールドを使用します）。

<strong>Snap Assist：</strong>システムのショートカット経由で Windows の 11 スナップレイアウトオプションを表示します。ネイティブのホバー時 Snap Layouts に対応したカスタム HTML 最大化ボタンを使用する場合は、代わりに [Windows のネイティブ非クライアント領域](/features/windows/frameless/#native-non-client-regions-on-windows)を使用してください。

<strong>タスクバーの点滅：</strong>ウィンドウが最小化されているときの通知に便利です。

[macOS]
**macOS 固有の機能：**

```go
// Transparent title bar
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Mac: application.MacWindow{
        TitleBar: application.MacTitleBar{
            AppearsTransparent: true,
        },
        Backdrop: application.MacBackdropTranslucent,
    },
})
```

**背景効果の種類：**

- `MacBackdropNormal` - 標準ウィンドウ
- `MacBackdropTranslucent` - 半透明の背景。WebView を透明にするには**プライベート API が必要です**。
- `MacBackdropTransparent` - 完全に透明。WebView を透明にするには**プライベート API が必要です**。
- `MacBackdropLiquidGlass` - ガラス風の背景。WebView を透明にするには**プライベート API が必要です**。

これらの効果を WebView 越しに表示するには、`-tags private_mac_apis` を指定してビルドします。指定しない場合、WebView は不透明なままです。`TitleBar.AppearsTransparent` 自体はパブリック API を使用します。[macOS のプライベート API](/guides/build/private-macos-apis/)を参照してください。

<strong>コレクション動作：</strong>複数の操作スペースにまたがるウィンドウの動作を制御します：

- `MacWindowCollectionBehaviorCanJoinAllSpaces` - すべての操作スペースに表示
- `MacWindowCollectionBehaviorFullScreenAuxiliary` - フルスクリーンアプリの上に重ねて表示可能

<strong>ネイティブフルスクリーン：</strong>macOS のフルスクリーンでは、新しい操作スペース（仮想デスクトップ）が作成されます。

[Linux]
**Linux 固有の機能：**

```go
// Set window icon (per-window struct is LinuxWindow, not the app-level LinuxOptions)
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Linux: application.LinuxWindow{
        Icon: iconBytes,
    },
})
```

**デスクトップ環境に関する注意事項：**

- GNOME：完全にサポート
- KDE Plasma：完全にサポート
- XFCE：部分的にサポート
- その他：環境によって異なります

**タイル型ウィンドウマネージャー（Hyprland、Sway、i3 など）：**

- `Minimise()` と `Maximise()` は期待どおりに動作しない場合があります。ウィンドウの位置とサイズはウィンドウマネージャーが制御します
- `SetSize()` と `SetPosition()` の要求は参考情報として扱われ、無視される場合があります
- `Fullscreen()` は通常、期待どおりに動作します
- 一部のウィンドウマネージャーは常に最前面に表示する機能をサポートしていません

@end

## 一般的なパターン

### スプラッシュ画面

```go
// Create splash screen
splash := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:     "Loading...",
    Width:     400,
    Height:    300,
    Frameless: true,
    AlwaysOnTop: true,
})

// Show splash
splash.Show()

// Initialise application
time.Sleep(2 * time.Second)

// Hide splash, show main window
splash.Close()
mainWindow.Show()
```

### 設定ウィンドウ

```go
var settingsWindow *application.WebviewWindow

func showSettings() {
    if settingsWindow == nil {
        settingsWindow = app.Window.NewWithOptions(application.WebviewWindowOptions{
            Name:   "settings",
            Title:  "Settings",
            Width:  600,
            Height: 400,
        })
    }
    
    settingsWindow.Show()
    settingsWindow.Focus()
}
```

### 閉じる前の確認

```go
window.RegisterHook(events.Common.WindowClosing, func(event *application.WindowEvent) {
    if hasUnsavedChanges() {
        // Show dialog
        result := showConfirmDialog("Unsaved changes. Close anyway?")
        if result != "yes" {
            // Cancel the close event
            event.Cancel()
        }
    }
})
```

## ベストプラクティス

### ✅ 推奨事項

- **重要なウィンドウに名前を付ける** - 後で簡単に検索できます
- **最小サイズを設定する** - 使用不能なレイアウトを防ぎます
- **ウィンドウを中央に配置する** - ランダムな位置に配置するよりもユーザー体験が向上します
- **クローズイベントを処理する** - データの損失を防ぎます
- **すべてのプラットフォームでテストする** - 動作はプラットフォームによって異なります
- **適切なサイズを使用する** - さまざまな画面サイズを考慮します

### ❌ 禁止事項

- **ウィンドウを作成しすぎない** - ユーザーの混乱を招きます
- **ウィンドウを閉じ忘れない** - メモリリークの原因になります
- **位置をハードコードしない** - 画面サイズは環境によって異なります
- **プラットフォーム間の違いを無視しない** - 十分にテストしてください
- **UI スレッドをブロックしない** - 長時間かかる処理には goroutine を使用してください

## トラブルシューティング

### ウィンドウが表示されない

**考えられる原因：**

1. ウィンドウが非表示の状態で作成されている
2. ウィンドウが画面外にある
3. ウィンドウがほかのウィンドウの背後にある

**解決策：**

```go
window.Show()
window.Center()
window.Focus()
```

### ウィンドウのサイズが正しくない

<strong>原因：</strong>Windows/Linux の DPI スケーリング

**解決策：**

```go
// Wails handles DPI automatically
// Just use logical pixels
window.SetSize(800, 600)
```

### ウィンドウがすぐに閉じる

<strong>原因：</strong>最後のウィンドウを閉じるとアプリケーションが終了する

**解決策：**

```go
app := application.New(application.Options{
    Mac: application.MacOptions{
        ApplicationShouldTerminateAfterLastWindowClosed: false,
    },
})
```

## 次のステップ

@cards{cols="2"}
⚙ ウィンドウオプション
すべてのウィンドウオプションの完全なリファレンスです。

[詳細を見る →](/features/windows/options/)

---
▣ 複数ウィンドウ
マルチウィンドウアプリケーションのパターンです。

[詳細を見る →](/features/windows/multiple/)

---
★ フレームレスウィンドウ
カスタムのウィンドウ装飾を作成します。

[詳細を見る →](/features/windows/frameless/)

---
🚀 ウィンドウイベント
ウィンドウのライフサイクルイベントを処理します。

[詳細を見る →](/features/windows/events/)

@end

---

**ご質問がありますか？** [Discord](https://discord.gg/JDdSxwjhGf) で質問するか、[ウィンドウのサンプル](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples)を確認してください。
