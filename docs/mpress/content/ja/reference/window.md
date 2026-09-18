---
title: "ウィンドウ API"
description: "ウィンドウ API の完全なリファレンス"
slug: "reference/window"
sourcePath: "reference/window.md"
---

## 概要

ウィンドウ API は、ウィンドウの外観、動作、ライフサイクルを制御するためのメソッドを提供します。ウィンドウインスタンスまたは `app.Window` マネージャーを介してアクセスします。

`Window` は `*application.WebviewWindow` が実装するインターフェースです。以下のメソッドシグネチャは `*WebviewWindow` に定義されています。多くの変更メソッドは、メソッドチェーンを可能にするために `Window` を返します。戻り値については、各メソッドの説明を参照してください。

**一般的な操作：**

- ウィンドウの作成と表示
- サイズ、位置、状態の制御
- ウィンドウイベントの処理
- ウィンドウコンテンツの管理
- 外観と動作の設定

## 表示と非表示

### Show()

ウィンドウを表示します。ウィンドウが非表示の場合は、表示状態になります。メソッドチェーン用にレシーバーを返します。

```go
func (w *WebviewWindow) Show() Window
```

**例：**

```go
window := app.Window.New()
window.Show()
```

### Hide()

ウィンドウを閉じずに非表示にします。ウィンドウはメモリ上に残り、再度表示できます。メソッドチェーン用にレシーバーを返します。

```go
func (w *WebviewWindow) Hide() Window
```

**例：**

```go
// Hide window temporarily
window.Hide()

// Show it again later
window.Show()
```

**ユースケース：**

- システムトレイに格納するアプリケーション
- ウィンドウを再利用するウィザードフロー
- 処理中の一時的な非表示

### Close()

ウィンドウを閉じます。これにより、`WindowClosing` イベントが発生します。

```go
func (w *WebviewWindow) Close()
```

**例：**

```go
window.Close()
```

**注：** 登録済みのフックが `event.Cancel()` を呼び出すと、ウィンドウは閉じられません。

## ウィンドウのプロパティ

### SetTitle()

ウィンドウのタイトルバーに表示するテキストを設定します。メソッドチェーン用にレシーバーを返します。

```go
func (w *WebviewWindow) SetTitle(title string) Window
```

**パラメーター：**

- `title` - 新しいウィンドウタイトル

**例：**

```go
window.SetTitle("My Application - Document.txt")
```

### Name()

ウィンドウを一意に識別する名前を返します。

```go
func (w *WebviewWindow) Name() string
```

**例：**

```go
name := window.Name()
fmt.Println("Window name:", name)

// Retrieve window by name later
if w, ok := app.Window.GetByName(name); ok {
    w.Focus()
}
```

## サイズと位置

### SetSize()

ウィンドウの寸法をピクセル単位で設定します。メソッドチェーン用にレシーバーを返します。

```go
func (w *WebviewWindow) SetSize(width, height int) Window
```

**パラメーター：**

- `width` - ウィンドウの幅（ピクセル単位）
- `height` - ウィンドウの高さ（ピクセル単位）

**例：**

```go
window.SetSize(1024, 768)
```

### Size()

現在のウィンドウの寸法を返します。

```go
func (w *WebviewWindow) Size() (width, height int)
```

**例：**

```go
width, height := window.Size()
fmt.Printf("Window is %dx%d\n", width, height)
```

### SetMinSize() / SetMaxSize()

ウィンドウの最小寸法と最大寸法を設定します。どちらもメソッドチェーン用にレシーバーを返します。

```go
func (w *WebviewWindow) SetMinSize(width, height int) Window
func (w *WebviewWindow) SetMaxSize(width, height int) Window
```

**例：**

```go
// Prevent window from being too small
window.SetMinSize(800, 600)

// Prevent window from being too large
window.SetMaxSize(1920, 1080)
```

### SetPosition()

画面の左上隅を基準としてウィンドウの位置を設定します。

```go
func (w *WebviewWindow) SetPosition(x, y int)
```

**パラメーター：**

- `x` - 水平方向の位置（ピクセル単位）
- `y` - 垂直方向の位置（ピクセル単位）

**例：**

```go
// Position window at top-left
window.SetPosition(0, 0)

// Position window 100px from top-left
window.SetPosition(100, 100)
```

### Position()

現在のウィンドウ位置を返します。

```go
func (w *WebviewWindow) Position() (x, y int)
```

**例：**

```go
x, y := window.Position()
fmt.Printf("Window is at (%d, %d)\n", x, y)
```

### Center()

ウィンドウを画面の中央に配置します。

```go
func (w *WebviewWindow) Center()
```

**例：**

```go
window := app.Window.New()
window.Center()
window.Show()
```

**注：** プライマリモニターの中央に配置されます。マルチモニター環境については、画面 API を参照してください。

### Focus()

ウィンドウを最前面に移動し、キーボードフォーカスを与えます。

```go
func (w *WebviewWindow) Focus()
```

**例：**

```go
// Bring window to front
window.Focus()
```

## ウィンドウの状態

### Minimise() / UnMinimise()

ウィンドウをタスクバー／Dock に最小化するか、元の状態に戻します。`Minimise()` はメソッドチェーン用にレシーバーを返しますが、`UnMinimise()` は何も返しません。

```go
func (w *WebviewWindow) Minimise() Window
func (w *WebviewWindow) UnMinimise()
```

**例：**

```go
// Minimise window
window.Minimise()

// Restore from minimised state
window.UnMinimise()
```

### Maximise() / UnMaximise()

ウィンドウを画面全体に最大化するか、以前のサイズに戻します。`Maximise()` はメソッドチェーン用にレシーバーを返しますが、`UnMaximise()` は何も返しません。

```go
func (w *WebviewWindow) Maximise() Window
func (w *WebviewWindow) UnMaximise()
```

**例：**

```go
// Maximise window
window.Maximise()

// Restore to previous size
window.UnMaximise()
```

### Fullscreen() / UnFullscreen() / ToggleFullscreen()

フルスクリーンモードに切り替えるか、フルスクリーンモードを終了します。`Fullscreen()` はメソッドチェーン用にレシーバーを返します。

```go
func (w *WebviewWindow) Fullscreen() Window
func (w *WebviewWindow) UnFullscreen()
func (w *WebviewWindow) ToggleFullscreen()
```

**例：**

```go
// Enter fullscreen
window.Fullscreen()

// Exit fullscreen
window.UnFullscreen()

// Or toggle
window.ToggleFullscreen()
```

`SetFullscreen(bool)` メソッドはありません。

### IsMinimised() / IsMaximised() / IsFullscreen()

現在のウィンドウの状態を確認します。

```go
func (w *WebviewWindow) IsMinimised() bool
func (w *WebviewWindow) IsMaximised() bool
func (w *WebviewWindow) IsFullscreen() bool
```

**例：**

```go
if window.IsMinimised() {
    window.UnMinimise()
}

if window.IsMaximised() {
    fmt.Println("Window is maximised")
}

if window.IsFullscreen() {
    window.UnFullscreen()
}
```

## ウィンドウのコンテンツ

### SetURL()

ウィンドウ内で指定した URL に移動します。メソッドチェーン用にレシーバーを返します。

```go
func (w *WebviewWindow) SetURL(url string) Window
```

**パラメーター：**

- `url` - 移動先の URL（埋め込みアセットには `http://wails.localhost/` を使用できます）

**例：**

```go
// Navigate to embedded page
window.SetURL("http://wails.localhost/settings.html")

// Navigate to external URL (if allowed)
window.SetURL("https://wails.io")
```

### SetHTML()

HTML 文字列からウィンドウのコンテンツを直接設定します。メソッドチェーン用にレシーバーを返します。

```go
func (w *WebviewWindow) SetHTML(html string) Window
```

**パラメーター：**

- `html` - 表示する HTML コンテンツ

**例：**

```go
html := `
<!DOCTYPE html>
<html>
<head><title>Dynamic Content</title></head>
<body>
    <h1>Hello from Go!</h1>
    <p>This content was generated dynamically.</p>
</body>
</html>
`
window.SetHTML(html)
```

**使用例：**

- 動的なコンテンツの生成
- フロントエンドのビルドプロセスを使用しないシンプルなウィンドウ
- エラーページまたはスプラッシュ画面

### Reload()

現在のウィンドウのコンテンツを再読み込みします。

```go
func (w *WebviewWindow) Reload()
```

**例：**

```go
// Reload current page
window.Reload()
```

**注：** 開発中やコンテンツを更新する必要がある場合に便利です。

## ウィンドウイベント

Wails には、ウィンドウイベントを処理するためのメソッドが 2 つあります。

- **OnWindowEvent()** - ウィンドウイベントをリッスンします（イベントを阻止することはできません）。
- **RegisterHook()** - ウィンドウイベントにフックします（`event.Cancel()` を呼び出すことでイベントを阻止できます）。

### OnWindowEvent()

ウィンドウイベントのコールバックを登録します。購読解除関数を返します。

```go
func (w *WebviewWindow) OnWindowEvent(
    eventType events.WindowEventType,
    callback func(event *WindowEvent),
) func()
```

**例：**

```go
import "github.com/wailsapp/wails/v3/pkg/events"

// Listen for window focus
window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
    app.Logger.Info("Window gained focus")
})

// Listen for window lost focus
window.OnWindowEvent(events.Common.WindowLostFocus, func(e *application.WindowEvent) {
    app.Logger.Info("Window lost focus")
})

// Listen for window resize
window.OnWindowEvent(events.Common.WindowDidResize, func(e *application.WindowEvent) {
    app.Logger.Info("Window resized")
})
```

**一般的なウィンドウイベント：**

- `events.Common.WindowClosing` - ウィンドウが閉じようとしている
- `events.Common.WindowFocus` - ウィンドウがフォーカスを取得した
- `events.Common.WindowLostFocus` - ウィンドウがフォーカスを失った
- `events.Common.WindowDidMove` - ウィンドウが移動した
- `events.Common.WindowDidResize` - ウィンドウのサイズが変更された
- `events.Common.WindowMinimise` - ウィンドウが最小化された
- `events.Common.WindowMaximise` - ウィンドウが最大化された
- `events.Common.WindowFullscreen` - ウィンドウがフルスクリーンになった
- `events.Common.WindowRuntimeReady` - ウィンドウ内のランタイムが初期化された

### RegisterHook()

ウィンドウイベントのフックを登録します。フックはリスナーより先に実行され、`event.Cancel()` を呼び出すことでイベントを阻止できます。購読解除関数を返します。

```go
func (w *WebviewWindow) RegisterHook(
    eventType events.WindowEventType,
    callback func(event *WindowEvent),
) func()
```

**例 — ウィンドウが閉じるのを阻止する：**

```go
import "github.com/wailsapp/wails/v3/pkg/events"

window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    confirm := app.Dialog.Question().
        SetTitle("Confirm Close").
        SetMessage("Are you sure you want to close this window?")

    yes := confirm.AddButton("Yes")
    no := confirm.AddButton("No")
    confirm.SetDefaultButton(yes)
    confirm.SetCancelButton(no)

    no.OnClick(func() {
        e.Cancel() // Prevent window from closing
    })

    confirm.Show()
})
```

**例 — 閉じる前に保存する：**

```go
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    if !hasUnsavedChanges {
        return
    }

    dlg := app.Dialog.Question().
        SetTitle("Unsaved Changes").
        SetMessage("Save changes before closing?")

    save := dlg.AddButton("Save")
    discard := dlg.AddButton("Don't Save")
    cancel := dlg.AddButton("Cancel")
    dlg.SetDefaultButton(save)
    dlg.SetCancelButton(cancel)

    save.OnClick(func() { saveData() })
    cancel.OnClick(func() { e.Cancel() })
    _ = discard // allow close

    dlg.Show()
})
```

### EmitEvent()

ウィンドウのフロントエンドにカスタムイベントを送出します。送出がフックによってキャンセルされた場合は `true` を返します。

```go
func (w *WebviewWindow) EmitEvent(name string, data ...any) bool
```

**パラメーター：**

- `name` - イベント名
- `data` - イベントとともに送信する任意のデータ

**例：**

```go
// Send data to specific window
window.EmitEvent("data-updated", map[string]any{
    "count":  42,
    "status": "success",
})
```

**フロントエンド（JavaScript）：**

```javascript
import { Events } from '@wailsio/runtime'

Events.On('data-updated', (data) => {
    console.log('Count:', data.count)
    console.log('Status:', data.status)
})
```

## その他のメソッド

### SetEnabled()

ウィンドウに対するユーザー操作を有効または無効にします。

```go
func (w *WebviewWindow) SetEnabled(enabled bool)
```

**例：**

```go
// Disable window during long operation
window.SetEnabled(false)

// Perform operation
performLongOperation()

// Re-enable window
window.SetEnabled(true)
```

### SetBackgroundColour()

ウィンドウの背景色（コンテンツが読み込まれる前に表示される色）を設定します。メソッドチェーン用にレシーバーを返します。

```go
func (w *WebviewWindow) SetBackgroundColour(colour RGBA) Window
```

`RGBA` は `application.RGBA{Red, Green, Blue, Alpha uint8}` です。ヘルパー `application.NewRGB(r, g, b)`（アルファ値 255）または `application.NewRGBA(r, g, b, a)` を使用してください。

**例：**

```go
// White background
window.SetBackgroundColour(application.NewRGB(255, 255, 255))

// Dark background with full alpha
window.SetBackgroundColour(application.NewRGBA(30, 30, 30, 255))
```

### SetResizable()

ユーザーがウィンドウのサイズを変更できるかどうかを制御します。メソッドチェーン用にレシーバーを返します。

```go
func (w *WebviewWindow) SetResizable(resizable bool) Window
```

**例：**

```go
// Make window fixed size
window.SetResizable(false)
```

### SetAlwaysOnTop()

ウィンドウをほかのウィンドウより常に手前に表示するかどうかを設定します。メソッドチェーン用にレシーバーを返します。

```go
func (w *WebviewWindow) SetAlwaysOnTop(alwaysOnTop bool) Window
```

**例：**

```go
// Keep window on top
window.SetAlwaysOnTop(true)
```

### Print()

ウィンドウのコンテンツに対するネイティブの印刷ダイアログを開きます。

```go
func (w *WebviewWindow) Print() error
```

**戻り値：** 印刷に失敗した場合はエラー。

**例：**

```go
if err := window.Print(); err != nil {
    log.Println("Print failed:", err)
}
```

### AttachModal()

2つ目のウィンドウをシートモーダルとしてアタッチします。

```go
func (w *WebviewWindow) AttachModal(modalWindow Window)
```

**パラメーター：**

- `modalWindow` - モーダルとしてアタッチするウィンドウ

**プラットフォーム対応：**

- **macOS**：完全対応（シートとして表示）
- **Windows**：未対応
- **Linux**：未対応

**例：**

```go
modalWindow := app.Window.New()
window.AttachModal(modalWindow)
```

## プラットフォーム固有のオプション

### Linux

Linuxのウィンドウでは、`LinuxWindow`を介して次のプラットフォーム固有のオプションを使用できます。

#### MenuStyle

アプリケーションメニューの表示方法を制御します。このオプションはデフォルトのGTK4ビルドで使用でき、従来の`-tags gtk3`ビルドでは無視されます。

| 値 | 説明 |
| --- | --- |
| `LinuxMenuStyleMenuBar` | タイトルバーの下に表示される従来型のメニューバー（デフォルト） |
| `LinuxMenuStylePrimaryMenu` | ヘッダーバー内のプライマリメニューボタン（GNOMEスタイル） |

**例：**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "My Application",
    Linux: application.LinuxWindow{
        MenuStyle: application.LinuxMenuStylePrimaryMenu,
    },
})
window.SetMenu(menu)
```

<strong>注：</strong>プライマリメニュースタイルでは、GNOME Human Interface Guidelinesに従って、ヘッダーバーにハンバーガーボタン（☰）が表示されます。これは最新のGNOMEアプリケーションに推奨されるスタイルです。

## 完全な例

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)

func main() {
    app := application.New(application.Options{
        Name: "Window API Demo",
    })

    // Create window with options
    window := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:            "My Application",
        Width:            1024,
        Height:           768,
        MinWidth:         800,
        MinHeight:        600,
        BackgroundColour: application.NewRGB(255, 255, 255),
        URL:              "http://wails.localhost/",
    })

    // Configure window behaviour
    window.SetResizable(true)
    window.SetMinSize(800, 600)
    window.SetMaxSize(1920, 1080)

    // Confirm-before-close hook
    window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
        dlg := app.Dialog.Question().
            SetTitle("Confirm Close").
            SetMessage("Are you sure you want to close this window?")

        yes := dlg.AddButton("Yes")
        no := dlg.AddButton("No")
        dlg.SetDefaultButton(yes)
        dlg.SetCancelButton(no)
        no.OnClick(func() { e.Cancel() })

        dlg.Show()
    })

    // Listen for window events
    window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
        window.SetTitle("My Application (Active)")
        app.Logger.Info("Window gained focus")
    })

    window.OnWindowEvent(events.Common.WindowLostFocus, func(e *application.WindowEvent) {
        window.SetTitle("My Application")
        app.Logger.Info("Window lost focus")
    })

    // Position and show window
    window.Center()
    window.Show()

    app.Run()
}
```
