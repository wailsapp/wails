---
title: "ウィンドウオプション"
description: "ウィンドウオプション、既定値、プラットフォームの制約"
slug: "features/windows/options"
sourcePath: "features/windows/options.md"
---

## ウィンドウ構成オプション

Wails では、サイズ、位置、外観、動作に関する多数のオプションを使用して、ウィンドウを包括的に構成できます。このリファレンスでは、Windows、macOS、Linux に共通する `WebviewWindowOptions` のフィールドについて説明します。以下の構造体は一部のフィールドを抜粋した概要です。完全な型については、リンク先の Go ソースを参照してください。

## WebviewWindowOptions 構造体

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

`WebviewWindowOptions` には `Parent` フィールドがありません。親子関係やモーダル関係には `parentWindow.AttachModal(childWindow)` を使用してください。また、`Assets` フィールドもありません。アセットの構成は `application.Options`（`Assets AssetOptions`）にあります。

完全なソース：[`v3/pkg/application/webview_window_options.go`](https://github.com/wailsapp/wails/blob/releases/v3.0.0/v3/pkg/application/webview_window_options.go)。

## 基本オプション

### Name

**型：** `string` **既定値：** `window-<numeric ID>` **プラットフォーム：** すべて

```go
Name: "main-window"
```

<strong>用途：</strong>後でウィンドウを検索するための一意の識別子。

**ベストプラクティス：**

- `"main"`、`"settings"`、`"about"` のような、内容が分かる名前を使用する
- `"file-browser"`、`"color-picker"` のようなケバブケースを使用する
- 短く覚えやすい名前にする

**例：**

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

**型：** `string` **既定値：** アプリケーション名 **プラットフォーム：** すべて

```go
Title: "My Application"
```

<strong>用途：</strong>タイトルバーとタスクバーに表示するテキスト。

**動的な更新：**

```go
window.SetTitle("My Application - Document.txt")
```

### Width / Height

**型：** `int`（ピクセル）**既定値：** 800 × 600 **プラットフォーム：** すべて **制約：** 正の値であること

```go
Width:  1200,
Height: 800,
```

<strong>用途：</strong>論理ピクセル単位の初期ウィンドウサイズ。

**注意事項：**

- Wails が DPI スケーリングを自動的に処理する
- 物理ピクセルではなく論理ピクセルを使用する
- 最低画面解像度（1024×768）を考慮する

**サイズの例：**

| 用途 | 幅 | 高さ |
| --- | --- | --- |
| 小型ユーティリティ | 400 | 300 |
| 標準アプリケーション | 1024 | 768 |
| 大型アプリケーション | 1440 | 900 |
| フルHD | 1920 | 1080 |

### X / Y

**型：** `int`（ピクセル）**既定値：** 画面中央 **プラットフォーム：** すべて

```go
X: 100,  // 100px from left edge
Y: 100,  // 100px from top edge
```

<strong>用途：</strong>ウィンドウの初期位置。

**座標系：**

- （0, 0）はプライマリ画面の左上
- X の正方向は右
- Y の正方向は下

**例：**

`X` と `Y` は、`InitialPosition: application.WindowXY` が設定されている場合にのみ有効です。設定されていない場合、`InitialPosition` の既定値は `WindowCentered` となり、`X`/`Y` は無視されます。

```go
settings := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:            "coordinate-window",
    InitialPosition: application.WindowXY, // opt into X/Y coordinates
    X:               100,
    Y:               100,
})
```

<strong>ベストプラクティス：</strong>特定の座標を指定する必要がない場合は、ウィンドウの作成後に `Center()` を使用して中央に配置します。

```go
window := app.Window.New()
window.Center()
```

### MinWidth / MinHeight

**型：** `int`（ピクセル）**既定値：** 0（最小値なし）**プラットフォーム：** すべて

```go
MinWidth:  400,
MinHeight: 300,
```

<strong>用途：</strong>ウィンドウが小さくなりすぎるのを防ぐ。

**ユースケース：**

- レイアウトの崩れを防ぐ
- 使いやすさを確保する
- アスペクト比を維持する

**例：**

```go
// Prevent window smaller than 400x300
MinWidth:  400,
MinHeight: 300,
```

### MaxWidth / MaxHeight

**型：** `int`（ピクセル）**既定値：** 0（最大値なし）**プラットフォーム：** すべて

```go
MaxWidth:  1920,
MaxHeight: 1080,
```

<strong>用途：</strong>ウィンドウが大きくなりすぎるのを防ぐ。

**ユースケース：**

- 固定サイズのアプリケーション
- 過剰なリソース使用を防ぐ
- 設計上の制約を維持する

## 状態オプション

### Hidden

**型：** `bool` **既定値：** `false` **プラットフォーム：** すべて

```go
Hidden: true,
```

<strong>用途：</strong>表示せずにウィンドウを作成する。

**ユースケース：**

- バックグラウンドウィンドウ
- 必要に応じて表示するウィンドウ
- スプラッシュ画面（作成、読み込み、表示の順に実行）
- コンテンツの読み込み中に白い画面が一瞬表示されるのを防ぐ

**プラットフォーム別の改善点：**

- <strong>Windows：</strong>白いウィンドウが一瞬表示される問題を修正。`Show()` が呼び出されるまでウィンドウは非表示のままになります
- <strong>macOS：</strong>完全対応
- <strong>Linux：</strong>完全対応

**スムーズに読み込むための推奨パターン：**

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

**例：**

```go
settings := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:   "settings",
    Hidden: true,
})

// Show when needed
settings.Show()
```

### Frameless

**型：** `bool` **既定値：** `false` **プラットフォーム：** すべて

```go
Frameless: true,
```

<strong>用途：</strong>タイトルバーとウィンドウの境界線を削除する。

**ユースケース：**

- カスタムウィンドウ装飾
- スプラッシュ画面
- キオスクアプリケーション
- 独自デザインのウィンドウ

**重要：** 次の機能を実装する必要があります。

- ウィンドウのドラッグ
- 閉じる／最小化／最大化ボタン
- サイズ変更用ハンドル（サイズ変更可能な場合）

**詳しくは、[フレームレスウィンドウ](/features/windows/frameless/)を参照してください。**

### DisableResize

**型：** `bool` **デフォルト：** `false`（デフォルトではウィンドウのサイズを変更可能） **プラットフォーム：** すべて

```go
DisableResize: true,
```

**目的：** ウィンドウのサイズ変更を禁止します。このフィールドは、v2の`Resizable`とは**逆**であることに注意してください。ウィンドウのサイズを変更不可にするには、`DisableResize: true`を設定します。

**使用例：**

- 固定サイズのアプリケーション
- スプラッシュ画面
- ダイアログ

**注：** `MaximiseButtonState`またはコレクションの動作で無効にしない限り、ユーザーは引き続き最大化やフルスクリーン表示を行えます。

### AlwaysOnTop

**型：** `bool` **デフォルト：** `false` **プラットフォーム：** すべて

```go
AlwaysOnTop: true,
```

**目的：** ウィンドウをほかのすべてのウィンドウより前面に維持します。

**使用例：**

- フローティングツールバー
- 通知
- ピクチャーインピクチャー
- タイマー

**プラットフォームに関する注意事項：**

- **macOS：** 完全対応
- **Windows：** 完全対応
- **Linux：** ウィンドウマネージャーによって異なる

### StartState

**型：** `WindowState` 列挙型 **デフォルト：** `WindowStateNormal` **プラットフォーム：** すべて

```go
StartState: application.WindowStateMaximised,
```

**目的：** 表示時のウィンドウの初期状態を指定します。

**値：**

- `WindowStateNormal` - 通常のウィンドウ
- `WindowStateMinimised` - 最小化
- `WindowStateMaximised` - 最大化
- `WindowStateFullscreen` - フルスクリーン

`WindowStateHidden`定数はありません。ウィンドウを非表示の状態で起動するには、真偽値フィールド`Hidden`を使用してください。

**実行時にフルスクリーンを切り替える：**

```go
window.Fullscreen()
window.UnFullscreen()
window.ToggleFullscreen() // there is no SetFullscreen(bool)
```

## 外観オプション

### BackgroundColour

**型：** `RGBA` 構造体 **デフォルト：** 白 **プラットフォーム：** すべて

```go
BackgroundColour: application.RGBA{Red: 0, Green: 0, Blue: 0, Alpha: 255},
```

`RGBA`の各フィールドは`Red, Green, Blue, Alpha`（uint8）です。ヘルパー`application.NewRGB(r, g, b)`（アルファ値255）または`application.NewRGBA(r, g, b, a)`の使用を推奨します。

**目的：** コンテンツが読み込まれる前のウィンドウの背景色を指定します。

**使用例：**

- アプリのテーマに合わせる
- ダークテーマで白く点滅するのを防ぐ
- 滑らかな読み込み体験を実現する

**例：**

```go
// Dark theme
BackgroundColour: application.NewRGB(30, 30, 30),

// Light theme
BackgroundColour: application.NewRGB(255, 255, 255),
```

**ヘルパーメソッド：**

```go
window.SetBackgroundColour(application.NewRGB(30, 30, 30))
```

### BackgroundType

**型：** `BackgroundType` 列挙型 **デフォルト：** `BackgroundTypeSolid` **プラットフォーム：** macOS、Windows（一部対応）

```go
BackgroundType: application.BackgroundTypeTranslucent,
```

**値：**

- `BackgroundTypeSolid` - 単色
- `BackgroundTypeTransparent` - 完全に透明
- `BackgroundTypeTranslucent` - 半透明のぼかし

**プラットフォーム対応：**

- **macOS：** `Mac.Backdrop`を設定します。WebViewを透明にするには[`private_mac_apis`](/guides/build/private-macos-apis/#webview-transparency-and-background)が必要です。これがない場合、WebViewは不透明のままです。
- **Windows：** 透明および半透明（Windows 11以降）
- **Linux：** 単色のみ

**例（macOS）：**

```go
BackgroundType: application.BackgroundTypeTranslucent,
Mac: application.MacWindow{
    Backdrop: application.MacBackdropTranslucent,
},
```

### OpenInspectorOnStartupとOpenDevTools

**macOSの非公開API：** `OpenInspectorOnStartup: true`、Goの`window.OpenDevTools()`、JavaScriptの`Window.OpenDevTools()`からプログラムでインスペクターを開くには、`private_mac_apis`が必要です。これがない場合、これらの操作では何も実行されません。プロダクションビルドでは`devtools`も必要です。macOS 13.3以降のSafariによる公開インスペクションには非公開APIは不要ですが、それより古いmacOSでインスペクターを有効にするには必要です。[Webインスペクターのビルドマトリックス](/guides/build/private-macos-apis/#web-inspector)を参照してください。

## コンテンツオプション

### URL

**型：** `string` **デフォルト：** 空（Assetsから読み込む） **プラットフォーム：** すべて

```go
URL: "https://example.com",
```

**目的：** 埋め込みアセットの代わりに外部URLを読み込みます。

**使用例：**

- 開発時（開発サーバーから読み込む）
- Webベースのアプリケーション
- ハイブリッドアプリケーション

**例：**

```go
// Development — point the window at the Vite dev server
URL: "http://localhost:9245",

// Production — embedded assets are configured at the application level
// (Assets is application.Options.Assets, not a WebviewWindowOptions field).
```

プロダクション環境向けのアプリケーションレベルのコード例：

```go
app := application.New(application.Options{
    Name: "My App",
    Assets: application.AssetOptions{
        Handler: application.AssetFileServerFS(assets),
    },
})
```

### HTML

**型：** `string` **デフォルト：** 空 **プラットフォーム：** すべて

```go
HTML: "<h1>Hello World</h1>",
```

**目的：** HTML文字列を直接読み込みます。

**使用例：**

- シンプルなウィンドウ
- 生成されたコンテンツ
- テスト

**例：**

```go
HTML: `
<!DOCTYPE html>
<html>
<head><title>Simple Window</title></head>
<body><h1>Hello from Wails!</h1></body>
</html>
`,
```

### アセット（アプリケーションレベルのみ）

アセット設定は<strong>`WebviewWindowOptions`のフィールドではありません</strong>。フロントエンドアセットは、アプリケーション自体が`application.Options.Assets`（`AssetOptions`）を介して配信します。すべてのウィンドウはそのアセットサーバーを継承します。

```go
app := application.New(application.Options{
    Assets: application.AssetOptions{
        Handler: application.AssetFileServerFS(assets),
    },
})
```

**詳細については、[ビルドシステム](/concepts/build-system/)を参照してください。**

### UseApplicationMenu

**型：** `bool` **デフォルト：** `false` **プラットフォーム：** Windows、Linux（macOSでは効果なし）

```go
UseApplicationMenu: true,
```

<strong>目的：</strong>このウィンドウで、`app.Menu.Set()`を介して設定したアプリケーションメニューを使用します。

**macOS**では常に画面上部のグローバルアプリケーションメニューが使用されるため、このオプションは効果がありません。

**Windows**および**Linux**では、デフォルトでウィンドウにメニューが表示されません。`UseApplicationMenu: true`を設定すると、ウィンドウでアプリケーションレベルのメニューが使用され、シンプルなクロスプラットフォーム対応が可能になります。

**例：**

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

**注意事項：**

- `UseApplicationMenu`とウィンドウ固有のメニューの両方が設定されている場合は、ウィンドウ固有のメニューが優先されます
- 実行時にOSを確認する必要がなくなるため、クロスプラットフォームコードを簡素化できます
- メニューに関する完全なドキュメントについては、[アプリケーションメニュー](/features/menus/application/)を参照してください

## 入力オプション

### EnableFileDrop

<strong>型：</strong> `bool` <strong>デフォルト：</strong> `false` <strong>プラットフォーム：</strong>すべて

```go
EnableFileDrop: true,
```

<strong>目的：</strong>オペレーティングシステムからウィンドウへのファイルのドラッグ＆ドロップを有効にします。

有効にすると、次のように動作します。

- ファイルマネージャーからドラッグしたファイルをアプリケーションにドロップできます
- ドロップされたファイルのパスとともに`WindowFilesDropped`イベントが発生します
- `data-file-drop-target`属性を持つ要素から、ドロップに関する詳細情報を取得できます

**ユースケース：**

- ファイルアップロードインターフェース
- ドキュメントエディター
- メディアインポーター
- ファイルを受け付けるあらゆるアプリ

**例：**

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

**HTMLドロップゾーン：**

```html
<!-- Mark elements as drop targets -->
<div id="upload" data-file-drop-target>
    Drop files here
</div>
```

**完全なドキュメントについては、[ファイルのドロップ](/features/drag-and-drop/files/)を参照してください。**

## セキュリティオプション

### ContentProtectionEnabled

**型：** `bool` **デフォルト：** `false` **プラットフォーム：** Windows（10以降）、macOS

```go
ContentProtectionEnabled: true,
```

<strong>目的：</strong>ウィンドウ内容の画面キャプチャを防止します。

**プラットフォーム対応：**

- **Windows：** Windows 10ビルド19041以降（完全対応）、それより前のバージョン（一部対応）
- <strong>macOS：</strong>完全対応
- <strong>Linux：</strong>未対応

**ユースケース：**

- 銀行アプリケーション
- パスワードマネージャー
- 医療記録
- 機密文書

**重要な注意事項：**

1. 物理的なカメラによる撮影は防止できません
2. 一部のツールでは保護を回避できる場合があります
3. 包括的なセキュリティ対策の一部であり、これだけで保護できるものではありません
4. DevToolsウィンドウは自動的には保護されません

**例：**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Secure Window",
    ContentProtectionEnabled: true,
})

// Toggle at runtime
window.SetContentProtection(true)
```

### 権限

**型：** `map[PermissionType]Permission` **デフォルト：** `nil`（プラットフォームのデフォルト処理） **プラットフォーム：** Linux、Windows（macOSではTCCに委ねられます）

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionAllow,
    application.PermissionCamera:     application.PermissionAllow,
},
```

<strong>目的：</strong>ウィンドウ内のWebコンテンツからの機能要求（カメラ、マイク、位置情報、通知、クリップボードの読み取り）を処理する方法を、プラットフォーム固有のコードを使用せずに宣言的に制御します。

**PermissionTypeの値：** `PermissionMicrophone`、`PermissionCamera`、`PermissionGeolocation`、`PermissionNotifications`、`PermissionClipboardRead`

**Permissionの値：**

- `PermissionDefault`（0）— プラットフォームのネイティブ処理：macOS／WindowsではOS／WebView2のプロンプトを表示します。Linuxではカメラとマイクを許可し、それ以外はすべて拒否します
- `PermissionAllow`（1）— プロンプトを表示せずに許可します（Linuxではカメラとマイクのみ実装されており、その他の種類は引き続き拒否されます）
- `PermissionDeny`（2）— プロンプトを表示せずに拒否します

<strong>重要 — Windows：</strong>このオプションが導入される前は、WailsはすべてのWebView2機能を暗黙的に許可していました。現在は、`Permissions`にいずれかの項目を設定すると、その一括許可が無効になります。リストに含まれていない機能は自動的に許可されず、WebView2のネイティブプロンプトが表示されます。アプリに必要な機能をすべて明示的に列挙してください。

**完全なガイド、プラットフォーム対応表、例については、[権限](/features/windows/permissions/)を参照してください。**

## ウィンドウのライフサイクルイベント

ウィンドウのライフサイクルイベントは、`OnWindowEvent`と`RegisterHook`を使用して処理します。これらのメソッドにより、ウィンドウを閉じる処理と破棄処理を細かく制御できます。

### ウィンドウを閉じる操作のキャンセル

未保存の変更がある場合など、ウィンドウが閉じないようにするには、`WindowClosing`イベントで`RegisterHook`を使用し、`event.Cancel()`を呼び出します。

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

**要点：**

- `RegisterHook`はイベントが発生する前にインターセプトします
- ウィンドウが閉じないようにするには、`event.Cancel()`を呼び出します
- キャンセル後もウィンドウは開いたままになります

### ウィンドウを閉じる処理

ウィンドウを閉じる際にクリーンアップを実行するには、`WindowClosing`イベントで`OnWindowEvent`を使用します。

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

**要点：**

- `OnWindowEvent`は、まさに発生しようとしているイベントを処理します
- ウィンドウが破棄される前にクリーンアップが実行されます
- ここから閉じる処理をキャンセルすることはできません（キャンセルするには`RegisterHook`を使用してください）

### シングルトンウィンドウのクリーンアップパターン

シングルトンウィンドウ（インスタンスを1つだけに制限するウィンドウ）では、参照のクリーンアップに `WindowClosing` を使用します。

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

## プラットフォーム固有のオプション

### Mac のオプション

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

**タイトルバー**（`MacTitleBar`）

- `AppearsTransparent` - タイトルバーを透明にし、コンテンツをタイトルバー領域まで拡張します
- `Hide` - タイトルバーを完全に非表示にします
- `HideTitle` - タイトルテキストのみを非表示にします
- `FullSizeContent` - コンテンツをウィンドウ全体のサイズまで拡張します

macOS で WebView の透過とプログラムによるインスペクターの起動を使用するには、`private_mac_apis` ビルドタグが必要です。このタグがなくても同じオプションは有効ですが、プライベート API 専用の操作は何も行いません。Liquid Glass のグループ化も無視され、スタイルには公開 API による代替手段が使用されます。ビルドコマンドと正確な動作については、[macOS のプライベート API](/guides/build/private-macos-apis/)を参照してください。

**背景効果**（`MacBackdrop`）

- `MacBackdropNormal` - 標準の不透明な背景
- `MacBackdropTranslucent` - <strong>WebView の透過にはプライベート API が必要です。</strong>このタグがない場合、ネイティブのぼかし効果は不透明な WebView の背後に残ります。
- `MacBackdropTransparent` - <strong>WebView の透過にはプライベート API が必要です。</strong>このタグがない場合、WebView は不透明なままです。
- `MacBackdropLiquidGlass` - <strong>WebView の透過にはプライベート API が必要です。</strong>このタグがない場合、ガラスレイヤーは不透明な WebView の背後に残り、スタイルには公開 API による代替手段が使用されます。

**Liquid Glass**（`MacLiquidGlass`）

| フィールドまたは値 | macOS でのプライベート API への依存 |
| --- | --- |
| `Style: LiquidGlassStyleAutomatic` | ネイティブの通常スタイルは公開 API です。背景効果を適用した WebView の透過には `private_mac_apis` が必要です。 |
| `Style: LiquidGlassStyleLight` | このタグを指定すると、既存のネイティブなクリアスタイルへのマッピングが維持されます。指定しない場合、Wails は Aqua アピアランスの通常のガラス効果を使用します。 |
| `Style: LiquidGlassStyleDark` | <strong>プライベート API：</strong>文書化されていないネイティブスタイル値 `2`。このタグがない場合は、Dark Aqua アピアランスの通常のガラス効果になります。 |
| `Style: LiquidGlassStyleVibrant` | ネイティブのクリアスタイルへのマッピングは公開 API です。背景効果を適用した WebView の透過には、このタグが必要です。 |
| `GroupID` | <strong>プライベート API：</strong>空でない値はグループ化を要求します。このタグがない場合は無視されます。 |
| `GroupSpacing` | <strong>プライベート API：</strong>正の値はグループ間隔を要求します。このタグがない場合は無視されます。 |
| `Material`, `CornerRadius`, `TintColor` | これら自体はプライベート API に依存しません。 |

ネイティブスタイルの値と対応 OS については、[Liquid Glass の値](/guides/build/private-macos-apis/#liquid-glass-values)を参照してください。

**非表示タイトルバーの高さ**（`int`）

- 非表示のタイトルバー領域の高さ（ドラッグ用）
- ネイティブのタイトルバードラッグ領域が非表示の場合、つまりウィンドウがフレームレス（`Frameless: true`）であるか、透明なタイトルバー（`AppearsTransparent: true`）を使用している場合にのみ有効です
- タイトルバーが表示される標準ウィンドウには効果がありません

**ウィンドウクラス**（`MacWindowClass`）

- `MacWindowClassWindow` - 標準の `NSWindow` 動作（デフォルト）
- `MacWindowClassPanel` - アプリケーションのメインウィンドウにはならない補助的な `NSPanel`

`PanelPreferences` は `MacWindowClassPanel` にのみ適用されます。

- `NonActivating` は `NSWindowStyleMaskNonactivatingPanel` を追加します。パネルを表示またはフォーカスしても Wails アプリケーションはアクティブになりませんが、コントロールやテキスト入力のためにパネルがキーウィンドウになることはできます。
- `FloatingPanel` は AppKit のフローティングパネル動作を有効にします。
- `BecomesKeyOnlyIfNeeded` は、クリックされたビューがキーボード入力を要求した場合にのみキー状態になります。
- `UtilityWindow` はネイティブのユーティリティウィンドウスタイルを適用します。

Wails のパネルは、アプリケーションが非アクティブになっても表示されたままで、閉じると解放されます。これは `WebviewWindow` で想定されるライフサイクルと一致します。これらは、`NSPanel` の反対のデフォルト動作を意図的に上書きするものです。

ウィンドウクラス、ウィンドウレベル、アクティベーションポリシー、コレクション動作は、それぞれ異なる問題を解決します。

- `WindowClass` は `NSWindow` または `NSPanel` を選択し、メインウィンドウとキーウィンドウのセマンティクスを制御します。
- `WindowLevel` は Z オーダーを制御します。メニューバーのオーバーレイには `MacWindowLevelPopUpMenu` を使用します。
- `MacOptions.ActivationPolicy` は、Dock やメニューバーの表示を含め、アプリケーション全体を制御します。非アクティブ化パネルにアクセサリアクティベーションポリシーは必要ありませんが、アプリで Dock アイコンを非表示にするために使用することはできます。
- `CollectionBehavior` は Spaces とフルスクリーンへの参加を制御します。

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

**ウィンドウレベル**（`MacWindowLevel`）

- `MacWindowLevelNormal` - 標準のウィンドウレベル（デフォルト）
- `MacWindowLevelFloating` - 通常のウィンドウより前面に表示されます
- `MacWindowLevelTornOffMenu` - 切り離されたメニューのレベル
- `MacWindowLevelModalPanel` - モーダルパネルのレベル
- `MacWindowLevelMainMenu` - メインメニューのレベル
- `MacWindowLevelStatus` - ステータスウィンドウのレベル
- `MacWindowLevelPopUpMenu` - ポップアップメニューのレベル
- `MacWindowLevelScreenSaver` - スクリーンセーバーのレベル

明示的な `WindowLevel` は、`AlwaysOnTop` と `PanelPreferences.FloatingPanel` より優先されます。レベルが明示されていない場合、`AlwaysOnTop` またはフローティングパネルは `MacWindowLevelFloating` に解決され、それ以外の場合のレベルは `MacWindowLevelNormal` になります。後から `SetAlwaysOnTop` を呼び出した場合も、明示的なランタイム変更として扱われます。

**コレクション動作**（`MacWindowCollectionBehavior`）

macOS の Spaces とフルスクリーンでウィンドウがどのように動作するかを制御します。これらはビットマスク値であり、ビット単位 OR（`|`）で組み合わせることができます。

**Space での動作：**

- `MacWindowCollectionBehaviorDefault` - FullScreenPrimary を使用します（デフォルト、後方互換）
- `MacWindowCollectionBehaviorCanJoinAllSpaces` - すべての Space にウィンドウを表示します
- `MacWindowCollectionBehaviorMoveToActiveSpace` - 表示時にアクティブな Space へ移動します
- `MacWindowCollectionBehaviorManaged` - デフォルトの管理対象ウィンドウ動作
- `MacWindowCollectionBehaviorTransient` - 一時的なウィンドウ
- `MacWindowCollectionBehaviorStationary` - Space の切り替え中も同じ位置に留まります

**ウィンドウの循環切り替え：**

- `MacWindowCollectionBehaviorParticipatesInCycle` - Cmd+` による循環切り替えの対象に含めます
- `MacWindowCollectionBehaviorIgnoresCycle` - Cmd+` による循環切り替えの対象から除外します

**フルスクリーン動作：**

- `MacWindowCollectionBehaviorFullScreenPrimary` - フルスクリーンモードに移行可能
- `MacWindowCollectionBehaviorFullScreenAuxiliary` - フルスクリーンアプリの上に重ねて表示可能
- `MacWindowCollectionBehaviorFullScreenNone` - フルスクリーン機能を無効化
- `MacWindowCollectionBehaviorFullScreenAllowsTiling` - 左右に並べるタイル表示を許可（macOS 10.11以降）
- `MacWindowCollectionBehaviorFullScreenDisallowsTiling` - タイル表示を禁止（macOS 10.11以降）

**例 — Spotlight風のウィンドウ：**

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

**例 — 単一の動作：**

```go
// Window that can appear over fullscreen applications
Mac: application.MacWindow{
    CollectionBehavior: application.MacWindowCollectionBehaviorFullScreenAuxiliary,
},
```

**TabbingMode**（`MacWindowTabbingMode`）

macOS 10.12以降で、ウィンドウのタブ化動作を制御します。ウィンドウのタブ化により、複数のウィンドウをタブとしてグループ化できます。

**オプション：**

- `MacWindowTabbingModeDefault` - ゼロ値を示す番兵値（明示的に未設定）。実行時には、デフォルトでタブ化が禁止されます
- `MacWindowTabbingModeAutomatic` - システムがタブ化動作を決定
- `MacWindowTabbingModePreferred` - ウィンドウでタブ化モードを優先
- `MacWindowTabbingModeDisallowed` - ウィンドウのタブ化を無効化

**例 — ウィンドウのタブ化を無効にする：**

```go
Mac: application.MacWindow{
    TabbingMode: application.MacWindowTabbingModeDisallowed,
},
```

**例 — ウィンドウのタブ化を優先する：**

```go
Mac: application.MacWindow{
    TabbingMode: application.MacWindowTabbingModePreferred,
},
```

**WebviewPreferences**（`MacWebviewPreferences`）

基盤となる`WKWebView`設定をきめ細かく制御します。すべてのフィールドは任意です。未設定のフィールドについては、WebKitのデフォルトがそのまま使用されます。

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

- `TabFocusesLinks` — `true`の場合、Tabキーを押すとリンクおよびフォームコントロールへフォーカスが移動します（デフォルト：`false`）
- `TextInteractionEnabled` — `true`の場合、ユーザーはWebView内のテキストを選択して操作できます（デフォルト：`true`）
- `FullscreenEnabled` — `true`の場合、WebコンテンツはHTML Fullscreen APIを介してフルスクリーン表示に移行できます（デフォルト：`false`）。macOS 12.3以降が必要です。
- `AllowsBackForwardNavigationGestures` — `true`の場合、水平方向のスワイプジェスチャーで前後のページへ移動します（デフォルト：`false`）
- `AllowsMagnification` — `true`の場合、WebViewでピンチズームが有効になります（デフォルト：`false`）
- `AllowsAirPlayForMediaPlayback` — `true`の場合、メディアをAirPlayデバイスへストリーミングできます（デフォルト：`true`）
- `JavaScriptCanOpenWindowsAutomatically` — `true`の場合、JavaScriptはユーザージェスチャーなしで新しいウィンドウを開けます（デフォルト：`false`）
- `MinimumFontSize` — ポイント単位の最小フォントサイズ。設定するには`optional.NewVar(12.0)`を使用します。未設定の場合はWebKitのデフォルトが使用されます。
- `ApplicationNameForUserAgent` — WebKitのユーザーエージェント文字列に含まれるアプリケーション名のサフィックスを上書きします。サイトがデフォルトの`"wails.io"`識別子を拒否する場合（YouTubeの埋め込みなど）に便利です。デフォルトを維持するには空のままにします。
- `EnableAutoplayWithoutUserAction` — `true`の場合、音声と動画をユーザージェスチャーなしで自動再生できます。`WKWebViewConfiguration.mediaTypesRequiringUserActionForPlayback = WKAudiovisualMediaTypeNone`に対応します（デフォルト：`false`）

### Windowsオプション（ウィンドウ単位）

ウィンドウ単位の構造体は`application.WindowsWindow`です。`WindowsOptions`では**ありません**（こちらは*アプリケーション*レベルの構造体です）。

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

**DisableIcon**（`bool`）

- タイトルバーからアイコンを削除します。

**DisableMenu**（`bool`）

- ウィンドウのメニューバーを無効にします。`true`の場合、メニューバーが設定されていてもウィンドウには表示されません。
- デフォルト：`false`

**BackdropType**（`BackdropType`）

- `application.Auto` - システムのデフォルト
- `application.None` - 背景効果なし
- `application.Mica` - Micaマテリアル（Windows 11）
- `application.Acrylic` - Acrylicマテリアル（Windows 11）
- `application.Tabbed` - Tabbedマテリアル（Windows 11）

`WindowsBackdropTypeMica`形式の定数はありません。`application.Mica`などを使用してください。

**CustomTheme**（`ThemeSettings`）

- 値型（ポインターではありません）。ウィンドウの境界線、タイトルバーのテキストと背景、およびメニューバーに、ダーク／ライトモード用のカスタム色を設定します。

**DisableFramelessWindowDecorations**（`bool`）

- フレームレスウィンドウのデフォルト装飾（Aeroシャドウ、角丸）を無効にします。

**NonClientRegionSupport**（`bool`）

- フレームレスのカスタムタイトルバーで、WebView2ネイティブの`app-region: drag`／`app-region: no-drag`サポートを有効にします。
- これは、単純なネイティブアプリのドラッグ専用です。ネイティブのカスタムキャプションボタンの動作や、カスタム最大化ボタンでのWindows 11 Snap Assist／Snap Layoutsは提供しません。

**WebView2CompositionHosting**（`bool`）

- カスタムキャプションボタンに対するWails管理の`--wails-non-client-region`サポートを有効にし、カスタム最大化ボタンでのWindows 11 Snap Assist／Snap Layoutsなど、Windowsネイティブの動作を提供します。
- 試験的な機能です。デフォルトのHWNDホスト型コントローラーではなく、`ICoreWebView2CompositionController`とDirectCompositionを介してWebView2をホストします。
- WebView2ネイティブの`app-region`サポートと、Wails管理のカスタムキャプションボタン領域の両方がウィンドウに必要な場合は、`NonClientRegionSupport`と組み合わせて使用できます。

**例：**

```go
Windows: application.WindowsWindow{
    BackdropType: application.Mica,
    DisableIcon:  true,
},
```

**例 — Windowsのカスタムタイトルバー領域：**

```go
Windows: application.WindowsWindow{
    NonClientRegionSupport:    true,
    WebView2CompositionHosting: true,
},
```

動作の詳細、トレードオフ、および対応するCSSについては、[フレームレスWindows](/features/windows/frameless/#native-non-client-regions-on-windows)を参照してください。

### Linuxオプション（ウィンドウ単位）

ウィンドウ単位の構造体は`application.LinuxWindow`です。`LinuxOptions`では**ありません**。

```go
Linux: application.LinuxWindow{
    Icon:                []byte{/* PNG data */},
    WindowIsTranslucent: false,
},
```

**Icon**（`[]byte`）

- ウィンドウアイコン（PNG形式）。

**WindowIsTranslucent**（`bool`）

- コンポジターのサポートが必要です。

**例：**

```go
//go:embed icon.png
var icon []byte

Linux: application.LinuxWindow{
    Icon: icon,
},
```

## アプリケーションレベルのWindowsオプション

Windows固有の一部のオプションは、ウィンドウ単位ではなくアプリケーションレベルで設定する必要があります。これは、WebView2がユーザーデータパスごとに単一のブラウザー環境を共有するためです。

### ブラウザーフラグ

WebView2のブラウザーフラグは、アプリケーション内の**すべてのウィンドウ**にわたって試験的機能と動作を制御します。これらは`application.Options.Windows`で設定する必要があります：

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

**EnabledFeatures**（`[]string`）

- 有効にする WebView2 機能フラグの一覧
- 利用可能なフラグについては、[WebView2 のブラウザフラグ](https://learn.microsoft.com/en-us/microsoft-edge/webview2/concepts/webview-features-flags)を参照してください
- 例: `"msWebView2EnableDraggableRegions"`

**DisabledFeatures**（`[]string`）

- 無効にする WebView2 機能フラグの一覧
- Wails は `msSmartScreenProtection` を自動的に無効にします
- 例: `"msExperimentalFeature"`

**AdditionalBrowserArgs**（`[]string`）

- ブラウザプロセスに渡す Chromium コマンドライン引数
- `--` プレフィックスを含める必要があります（例: `"--remote-debugging-port=9222"`）
- 利用可能な引数については、[Chromium のコマンドラインスイッチ](https://peter.sh/experiments/chromium-command-line-switches/)を参照してください

@note{type="caution" title="重要"}
WebView2 はユーザーデータパスごとに単一のブラウザ環境を共有するため、これらのフラグはすべてのウィンドウにグローバルに適用されます。ウィンドウごとに異なるブラウザフラグを指定することはできません。

@end

**完全な例:**

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

## 完全な例

以下は、本番環境で使用できるウィンドウ設定です:

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

フロントエンドアセットは、ウィンドウ単位ではなく、アプリケーションレベル（`application.Options.Assets`）で配信されます。

## 次のステップ

- [ウィンドウの基本](/features/windows/basics/) - ウィンドウの作成と制御
- [複数のウィンドウ](/features/windows/multiple/) - マルチウィンドウのパターン
- [フレームレスウィンドウ](/features/windows/frameless/) - カスタムウィンドウフレーム
- [ウィンドウイベント](/features/windows/events/) - ライフサイクルイベント

---

**質問がありますか？** [Discord](https://discord.gg/JDdSxwjhGf) で質問するか、[サンプル](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples)を確認してください。
