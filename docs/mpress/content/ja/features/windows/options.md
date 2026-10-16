---
title: "ウィンドウオプション"
description: "WebviewWindowOptions の完全なリファレンス"
slug: "features/windows/options"
sourcePath: "features/windows/options.md"
---

## ウィンドウ設定オプション

Wails には、サイズ、位置、外観、動作に関する数十のオプションが用意されており、ウィンドウを包括的に設定できます。このリファレンスでは、`WebviewWindowOptions`の<strong>完全なリファレンス</strong>として、Windows、macOS、Linux で利用可能なすべてのオプションを説明します。すべてのオプションと対応プラットフォームを、例および制約とともに示します。

## WebviewWindowOptions の構造

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

`WebviewWindowOptions`には<strong>`Parent`フィールドがありません</strong>。親子関係やモーダル関係には`parentWindow.AttachModal(childWindow)`を使用します。また、**`Assets`フィールドもありません**。アセットの設定は`application.Options`（`Assets AssetOptions`）にあります。

完全なソース：[`v3/pkg/application/webview_window_options.go`](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window_options.go)。

## 基本オプション

### 名前

**型：** `string` **デフォルト：** 自動生成された UUID **プラットフォーム：** すべて

```go
Name: "main-window"
```

**用途：** 後でウィンドウを検索するための一意の識別子。

**ベストプラクティス：**

- `"main"`、`"settings"`、`"about"`のように、用途が分かる名前を使用する
- `"file-browser"`、`"color-picker"`のように、ケバブケースを使用する
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

### タイトル

**型：** `string` **デフォルト：** アプリケーション名 **プラットフォーム：** すべて

```go
Title: "My Application"
```

**用途：** タイトルバーとタスクバーに表示されるテキスト。

**動的な更新：**

```go
window.SetTitle("My Application - Document.txt")
```

### 幅／高さ

**型：** `int`（ピクセル） **デフォルト：** 800 x 600 **プラットフォーム：** すべて **制約：** 正の値であること

```go
Width:  1200,
Height: 800,
```

**用途：** 論理ピクセル単位で指定するウィンドウの初期サイズ。

**注意事項：**

- Wails が DPI スケーリングを自動的に処理する
- 物理ピクセルではなく論理ピクセルを使用する
- 最小画面解像度（1024x768）を考慮する

**サイズの例：**

| 用途 | 幅 | 高さ |
| --- | --- | --- |
| 小規模なユーティリティ | 400 | 300 |
| 標準的なアプリ | 1024 | 768 |
| 大規模なアプリ | 1440 | 900 |
| フル HD | 1920 | 1080 |

### X／Y

**型：** `int`（ピクセル） **デフォルト：** 画面中央 **プラットフォーム：** すべて

```go
X: 100,  // 100px from left edge
Y: 100,  // 100px from top edge
```

**用途：** ウィンドウの初期位置。

**座標系：**

- （0, 0）はプライマリ画面の左上
- X の正方向は右
- Y の正方向は下

**例：**

`X`と`Y`は、`InitialPosition: application.WindowXY`が設定されている場合にのみ有効です。設定されていない場合、`InitialPosition`のデフォルトは`WindowCentered`となり、`X`／`Y`は無視されます。

```go
settings := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:            "coordinate-window",
    InitialPosition: application.WindowXY, // opt into X/Y coordinates
    X:               100,
    Y:               100,
})
```

**ベストプラクティス：** 特定の座標を指定する必要がない場合は、作成後に`Center()`を使用してウィンドウを中央に配置します：

```go
window := app.Window.New()
window.Center()
```

### MinWidth / MinHeight

**型：** `int`（ピクセル） **デフォルト：** 0（最小値なし） **プラットフォーム：** すべて

```go
MinWidth:  400,
MinHeight: 300,
```

**用途：** ウィンドウが小さくなりすぎるのを防ぐ。

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

**型：** `int`（ピクセル） **デフォルト：** 0（最大値なし） **プラットフォーム：** すべて

```go
MaxWidth:  1920,
MaxHeight: 1080,
```

**用途：** ウィンドウが大きくなりすぎるのを防ぐ。

**ユースケース：**

- 固定サイズのアプリケーション
- リソースの過剰な使用を防止
- 設計上の制約を維持

## 状態オプション

### Hidden

**型：** `bool` **デフォルト：** `false` **プラットフォーム：** すべて

```go
Hidden: true,
```

**目的：** ウィンドウを表示せずに作成します。

**用途：**

- バックグラウンドウィンドウ
- 必要に応じて表示するウィンドウ
- スプラッシュ画面（作成、読み込み後に表示）
- コンテンツの読み込み中に白い画面が一瞬表示されるのを防止

**プラットフォーム別の改善点：**

- **Windows：** ウィンドウが白く一瞬表示される問題を修正しました。`Show()`が呼び出されるまでウィンドウは非表示のままです
- **macOS：** 完全対応
- **Linux：** 完全対応

**滑らかに読み込むための推奨パターン：**

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

**型：** `bool` **デフォルト：** `false` **プラットフォーム：** すべて

```go
Frameless: true,
```

**目的：** タイトルバーとウィンドウの境界線を取り除きます。

**用途：**

- カスタムウィンドウ装飾
- スプラッシュ画面
- キオスクアプリケーション
- 独自デザインのウィンドウ

**重要：** 次の機能を実装する必要があります：

- ウィンドウのドラッグ移動
- 閉じる／最小化／最大化ボタン
- サイズ変更ハンドル（サイズ変更を許可する場合）

**詳細については、[フレームレスウィンドウ](/features/windows/frameless/)を参照してください。**

### DisableResize

**型：** `bool` **デフォルト：** `false`（デフォルトではウィンドウのサイズを変更可能） **プラットフォーム：** すべて

```go
DisableResize: true,
```

**目的：** ウィンドウのサイズ変更を禁止します。このフィールドは、v2の`Resizable`とは<strong>逆</strong>であることに注意してください。ウィンドウのサイズを変更できないようにするには、`DisableResize: true`を設定します。

**用途：**

- 固定サイズのアプリケーション
- スプラッシュ画面
- ダイアログ

**注意：** `MaximiseButtonState`またはコレクション動作でも無効にしない限り、ユーザーは引き続き最大化や全画面表示を行えます。

### AlwaysOnTop

**型：** `bool` **デフォルト：** `false` **プラットフォーム：** すべて

```go
AlwaysOnTop: true,
```

**目的：** ウィンドウをほかのすべてのウィンドウより前面に保ちます。

**用途：**

- フローティングツールバー
- 通知
- ピクチャーインピクチャー
- タイマー

**プラットフォーム別の注意事項：**

- **macOS：** 完全対応
- **Windows：** 完全対応
- **Linux：** ウィンドウマネージャーに依存

### StartState

**型：** `WindowState`列挙型 **デフォルト：** `WindowStateNormal` **プラットフォーム：** すべて

```go
StartState: application.WindowStateMaximised,
```

**目的：** ウィンドウを表示するときの初期状態を指定します。

**値：**

- `WindowStateNormal` - 通常のウィンドウ
- `WindowStateMinimised` - 最小化
- `WindowStateMaximised` - 最大化
- `WindowStateFullscreen` - 全画面表示

`WindowStateHidden`定数はありません。ウィンドウを非表示の状態で起動するには、`Hidden`ブールフィールドを使用します。

**実行時に全画面表示を切り替える：**

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

`RGBA` の各フィールドは `Red, Green, Blue, Alpha`（uint8）です。ヘルパー `application.NewRGB(r, g, b)`（アルファ値 255）または `application.NewRGBA(r, g, b, a)` の使用を推奨します。

**用途：** コンテンツが読み込まれるまでのウィンドウの背景色。

**ユースケース：**

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
- `BackgroundTypeTransparent` - 完全な透明
- `BackgroundTypeTranslucent` - 半透明のぼかし

**プラットフォームの対応状況：**

- **macOS：** `Mac.Backdrop` を設定します。WebViewを透明にするには [`private_mac_apis`](/guides/build/private-macos-apis/#webview-transparency-and-background) が必要です。これがない場合、WebViewは不透明なままです。
- **Windows：** 透明および半透明（Windows 11 以降）
- **Linux：** 単色のみ

**例（macOS）：**

```go
BackgroundType: application.BackgroundTypeTranslucent,
Mac: application.MacWindow{
    Backdrop: application.MacBackdropTranslucent,
},
```

### OpenInspectorOnStartupとOpenDevTools

**macOSでのプライベートAPI：** `OpenInspectorOnStartup: true`、Goの `window.OpenDevTools()`、JavaScriptの `Window.OpenDevTools()` からプログラムでインスペクターを開くには、`private_mac_apis` が必要です。これがない場合、これらの操作は何も行いません。プロダクションビルドでは `devtools` も必要です。macOS 13.3 以降のSafariによる公開インスペクションにはプライベートAPIは不要ですが、それより古いmacOSでインスペクターを有効にするには必要です。[Web Inspectorのビルドマトリックス](/guides/build/private-macos-apis/#web-inspector)を参照してください。

## コンテンツオプション

### URL

**型：** `string` **デフォルト：** 空（Assetsから読み込み） **プラットフォーム：** すべて

```go
URL: "https://example.com",
```

**用途：** 埋め込みアセットの代わりに外部URLを読み込みます。

**ユースケース：**

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

プロダクション環境向けのアプリケーションレベルのスニペット：

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

**用途：** HTML文字列を直接読み込みます。

**ユースケース：**

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

### Assets（アプリケーションレベルのみ）

アセット設定は **`WebviewWindowOptions` のフィールドではありません**。フロントエンドアセットは、アプリケーション自体が `application.Options.Assets`（`AssetOptions`）を介して配信します。すべてのウィンドウがそのアセットサーバーを継承します。

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

**用途：** このウィンドウで、`app.Menu.Set()` を介して設定したアプリケーションメニューを使用します。

**macOS** では、画面上部のグローバルアプリケーションメニューが常に使用されるため、このオプションは効果がありません。

**Windows** および **Linux** では、デフォルトでウィンドウにメニューは表示されません。`UseApplicationMenu: true` を設定すると、ウィンドウでアプリケーションレベルのメニューが使用され、シンプルなクロスプラットフォームソリューションを実現できます。

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

**注記：**

- `UseApplicationMenu`とウィンドウ固有のメニューの両方を設定した場合、ウィンドウ固有のメニューが優先されます
- 実行時にOSを判定する必要がなくなるため、クロスプラットフォーム対応のコードを簡潔にできます
- メニューの完全なドキュメントについては、[アプリケーションメニュー](/features/menus/application/)を参照してください

## 入力オプション

### EnableFileDrop

**型：** `bool` **デフォルト：** `false` **プラットフォーム：** すべて

```go
EnableFileDrop: true,
```

**目的：** オペレーティングシステムからウィンドウへのファイルのドラッグ＆ドロップを有効にします。

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

**HTMLのドロップゾーン：**

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

**目的：** ウィンドウ内容の画面キャプチャを防止します。

**プラットフォームの対応状況：**

- **Windows：** Windows 10 ビルド19041以降（完全対応）、それ以前のバージョン（一部対応）
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

**目的：** ウィンドウのWebコンテンツからの機能要求（カメラ、マイク、位置情報、通知、クリップボードの読み取り）の処理方法を、プラットフォーム固有のコードを使わずに宣言的に制御します。

**PermissionTypeの値：** `PermissionMicrophone`、`PermissionCamera`、`PermissionGeolocation`、`PermissionNotifications`、`PermissionClipboardRead`

**権限の値：**

- `PermissionDefault`（0）— プラットフォームのネイティブ処理。macOS／WindowsではOS／WebView2のプロンプトを表示します。Linuxではカメラとマイクを許可し、それ以外はすべて拒否します
- `PermissionAllow`（1）— プロンプトを表示せずに許可します（Linuxではカメラとマイクのみ実装されており、その他の種類は引き続き拒否されます）
- `PermissionDeny`（2）— プロンプトを表示せずに拒否します

<strong>重要 — Windows：</strong>このオプションが導入される前は、WailsはすべてのWebView2機能を通知せずに許可していました。現在は、`Permissions`にいずれかの項目を設定すると、その一括許可が無効になります。リストに含まれていない機能は自動的に許可されず、WebView2のネイティブプロンプトが表示されます。アプリに必要なすべての機能を明示的に列挙してください。

**完全なガイド、プラットフォーム対応表、例については、[権限](/features/windows/permissions/)を参照してください。**

## ウィンドウのライフサイクルイベント

ウィンドウのライフサイクルイベントは、`OnWindowEvent`と`RegisterHook`を使用して処理します。これらのメソッドを使用すると、ウィンドウを閉じる動作と破棄する動作を細かく制御できます。

### ウィンドウを閉じる処理のキャンセル

ウィンドウが閉じるのを防ぐには（未保存の変更がある場合など）、`WindowClosing`イベントに対して`RegisterHook`を使用し、`event.Cancel()`を呼び出します：

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

- `RegisterHook` は、イベントが発生する前にイベントをインターセプトします
- ウィンドウが閉じないようにするには、`event.Cancel()` を呼び出します
- キャンセル後もウィンドウは開いたままになります

### ウィンドウを閉じる処理

ウィンドウが閉じるときにクリーンアップを実行するには、`WindowClosing` イベントで `OnWindowEvent` を使用します：

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

- `OnWindowEvent` は、まもなく発生するイベントを処理します
- クリーンアップはウィンドウが破棄される前に実行されます
- ここから閉じる処理をキャンセルすることはできません（キャンセルするには `RegisterHook` を使用します）

### シングルトンウィンドウのクリーンアップパターン

シングルトンウィンドウ（インスタンスを1つだけに制限）では、参照のクリーンアップに `WindowClosing` を使用します：

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

### Macオプション

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

**TitleBar**（`MacTitleBar`）

- `AppearsTransparent` - タイトルバーを透明にし、コンテンツをタイトルバー領域まで拡張します
- `Hide` - タイトルバーを完全に非表示にします
- `HideTitle` - タイトルテキストのみを非表示にします
- `FullSizeContent` - コンテンツをウィンドウ全体のサイズまで拡張します

macOSでは、WebViewの透明化とプログラムによるインスペクターの起動に `private_mac_apis` ビルドタグが必要です。このタグがなくても同じオプションは有効ですが、プライベートAPIのみで可能な操作は何も実行しません。Liquid Glassのグループ化も無視され、スタイルには公開APIによる代替手段が使用されます。ビルドコマンドと正確な動作については、[macOSのプライベートAPI](/guides/build/private-macos-apis/)を参照してください。

**Backdrop**（`MacBackdrop`）

- `MacBackdropNormal` - 標準の不透明な背景
- `MacBackdropTranslucent` - <strong>WebViewの透明化にはプライベートAPIが必要です。</strong>このタグがない場合、ネイティブのぼかし効果は不透明なWebViewの背後に残ります。
- `MacBackdropTransparent` - <strong>WebViewの透明化にはプライベートAPIが必要です。</strong>このタグがない場合、WebViewは不透明なままです。
- `MacBackdropLiquidGlass` - <strong>WebViewの透明化にはプライベートAPIが必要です。</strong>このタグがない場合、ガラスレイヤーは不透明なWebViewの背後に残り、スタイルには公開APIによる代替手段が使用されます。

**LiquidGlass**（`MacLiquidGlass`）

| フィールドまたは値 | macOSでのプライベートAPIへの依存 |
| --- | --- |
| `Style: LiquidGlassStyleAutomatic` | ネイティブのregularスタイルは公開APIです。背景でWebViewを透明化するには `private_mac_apis` が必要です。 |
| `Style: LiquidGlassStyleLight` | このタグを指定すると、既存のネイティブclearスタイルへのマッピングが維持されます。指定しない場合、WailsはAquaアピアランスのregular glassを使用します。 |
| `Style: LiquidGlassStyleDark` | <strong>プライベートAPI：</strong>文書化されていないネイティブスタイル値 `2`。このタグがない場合は、Dark Aquaアピアランスのregular glassを使用します。 |
| `Style: LiquidGlassStyleVibrant` | ネイティブのclearスタイルへのマッピングは公開APIです。背景でWebViewを透明化するには、このタグが必要です。 |
| `GroupID` | <strong>プライベートAPI：</strong>空でない値を指定するとグループ化が要求されます。このタグがない場合は無視されます。 |
| `GroupSpacing` | <strong>プライベートAPI：</strong>正の値を指定するとグループ間隔が要求されます。このタグがない場合は無視されます。 |
| `Material`、`CornerRadius`、`TintColor` | これら自体はプライベートAPIに依存しません。 |

ネイティブスタイルの値と対応OSについては、[Liquid Glassの値](/guides/build/private-macos-apis/#liquid-glass-values)を参照してください。

**InvisibleTitleBarHeight**（`int`）

- 非表示のタイトルバー領域の高さ（ドラッグ用）
- ネイティブのタイトルバーのドラッグ領域が非表示の場合、つまりウィンドウがフレームレス（`Frameless: true`）であるか、透明なタイトルバー（`AppearsTransparent: true`）を使用している場合にのみ有効です
- タイトルバーが表示される標準ウィンドウには影響しません

**WindowClass**（`MacWindowClass`）

- `MacWindowClassWindow` - 標準の `NSWindow` の動作（デフォルト）
- `MacWindowClassPanel` - アプリケーションのメインウィンドウには決してならない補助的な `NSPanel`

`PanelPreferences` は `MacWindowClassPanel` にのみ適用されます：

- `NonActivating`は`NSWindowStyleMaskNonactivatingPanel`を追加します。パネルを表示またはフォーカスしても Wails アプリケーションはアクティブになりませんが、コントロールやテキスト入力のためにパネルをキーウインドウにすることはできます。
- `FloatingPanel`は、AppKit のフローティングパネル動作を有効にします。
- `BecomesKeyOnlyIfNeeded`は、クリックされたビューがキーボード入力を要求した場合にのみキーウインドウになります。
- `UtilityWindow`は、ネイティブのユーティリティウインドウスタイルを適用します。

Wails のパネルは、アプリケーションが非アクティブになっても表示されたままで、閉じると解放されます。これは`WebviewWindow`で想定されるライフサイクルに一致します。これらは、逆の動作をする`NSPanel`のデフォルト値を意図的に上書きするものです。

ウインドウクラス、ウインドウレベル、アクティベーションポリシー、コレクション動作は、それぞれ異なる問題を解決します：

- `WindowClass`は`NSWindow`または`NSPanel`を選択し、メインウインドウおよびキーウインドウのセマンティクスを制御します。
- `WindowLevel`は Z オーダーを制御します。メニューバーのオーバーレイには`MacWindowLevelPopUpMenu`を使用します。
- `MacOptions.ActivationPolicy`は、Dock とメニューバーの表示を含むアプリケーション全体を制御します。アプリケーションをアクティブにしないパネルにアクセサリアクティベーションポリシーは不要ですが、アプリケーションが Dock アイコンを非表示にするためにこのポリシーを使用することはできます。
- `CollectionBehavior`は、Spaces とフルスクリーンへの参加を制御します。

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

**WindowLevel**（`MacWindowLevel`）

- `MacWindowLevelNormal` - 標準のウインドウレベル（デフォルト）
- `MacWindowLevelFloating` - 通常のウインドウより前面に表示
- `MacWindowLevelTornOffMenu` - 切り離されたメニューのレベル
- `MacWindowLevelModalPanel` - モーダルパネルのレベル
- `MacWindowLevelMainMenu` - メインメニューのレベル
- `MacWindowLevelStatus` - ステータスウインドウのレベル
- `MacWindowLevelPopUpMenu` - ポップアップメニューのレベル
- `MacWindowLevelScreenSaver` - スクリーンセーバーのレベル

明示的な`WindowLevel`は、`AlwaysOnTop`および`PanelPreferences.FloatingPanel`より優先されます。レベルが明示されていない場合、`AlwaysOnTop`またはフローティングパネルは`MacWindowLevelFloating`に解決され、それ以外ではレベルは`MacWindowLevelNormal`になります。後から`SetAlwaysOnTop`を呼び出した場合も、明示的な実行時変更として扱われます。

**CollectionBehavior**（`MacWindowCollectionBehavior`）

macOS の Spaces およびフルスクリーンでのウインドウの動作を制御します。これらは、ビット単位 OR（`|`）で組み合わせられるビットマスク値です。

**Space の動作：**

- `MacWindowCollectionBehaviorDefault` - FullScreenPrimary を使用（デフォルト、後方互換）
- `MacWindowCollectionBehaviorCanJoinAllSpaces` - すべての Space にウインドウを表示
- `MacWindowCollectionBehaviorMoveToActiveSpace` - 表示時にアクティブな Space へ移動
- `MacWindowCollectionBehaviorManaged` - デフォルトの管理対象ウインドウ動作
- `MacWindowCollectionBehaviorTransient` - 一時的／過渡的なウインドウ
- `MacWindowCollectionBehaviorStationary` - Space の切り替え中も同じ位置に留まる

**ウインドウの巡回：**

- `MacWindowCollectionBehaviorParticipatesInCycle` - Cmd+`による巡回の対象に含める
- `MacWindowCollectionBehaviorIgnoresCycle` - Cmd+`による巡回の対象から除外

**フルスクリーンの動作：**

- `MacWindowCollectionBehaviorFullScreenPrimary` - フルスクリーンモードへの移行が可能
- `MacWindowCollectionBehaviorFullScreenAuxiliary` - フルスクリーンアプリケーションへのオーバーレイが可能
- `MacWindowCollectionBehaviorFullScreenNone` - フルスクリーン機能を無効化
- `MacWindowCollectionBehaviorFullScreenAllowsTiling` - 左右に並べるタイル表示を許可（macOS 10.11以降）
- `MacWindowCollectionBehaviorFullScreenDisallowsTiling` - タイル表示を禁止（macOS 10.11以降）

**例 - Spotlight のようなウインドウ：**

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

**例 - 単一の動作：**

```go
// Window that can appear over fullscreen applications
Mac: application.MacWindow{
    CollectionBehavior: application.MacWindowCollectionBehaviorFullScreenAuxiliary,
},
```

**TabbingMode**（`MacWindowTabbingMode`）

macOS 10.12以降でのウインドウのタブ化動作を制御します。ウインドウのタブ化により、複数のウインドウをタブとしてグループ化できます。

**オプション：**

- `MacWindowTabbingModeDefault` - ゼロ値の番兵値（明示的に未設定）。実行時のデフォルトではタブ化を許可しない
- `MacWindowTabbingModeAutomatic` - システムがタブ化の動作を決定
- `MacWindowTabbingModePreferred` - ウインドウでタブモードを優先
- `MacWindowTabbingModeDisallowed` - ウインドウのタブ化を無効化

**例 - ウインドウのタブ化を無効化：**

```go
Mac: application.MacWindow{
    TabbingMode: application.MacWindowTabbingModeDisallowed,
},
```

**例 - ウインドウのタブ化を優先：**

```go
Mac: application.MacWindow{
    TabbingMode: application.MacWindowTabbingModePreferred,
},
```

**WebviewPreferences**（`MacWebviewPreferences`）

基盤となる`WKWebView`の構成をきめ細かく制御します。すべてのフィールドは任意です。未設定のフィールドについては、WebKit のデフォルトから変更されません。

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

- `TabFocusesLinks` — `true`の場合、Tab キーを押すとリンクおよびフォームコントロールにフォーカスが移動します（デフォルト：`false`）
- `TextInteractionEnabled` — `true`の場合、ユーザーは WebView 内のテキストを選択して操作できます（デフォルト：`true`）
- `FullscreenEnabled` — `true`の場合、Web コンテンツは HTML Fullscreen API を介してフルスクリーン表示に移行できます（デフォルト：`false`）。macOS 12.3 以降が必要です。
- `AllowsBackForwardNavigationGestures` — `true`の場合、水平方向のスワイプジェスチャーで「戻る」または「進む」ナビゲーションが実行されます（デフォルト：`false`）
- `AllowsMagnification` — `true`の場合、WebView でピンチズームが有効になります（デフォルト：`false`）
- `AllowsAirPlayForMediaPlayback` — `true`の場合、メディアを AirPlay デバイスへストリーミングできます（デフォルト：`true`）
- `JavaScriptCanOpenWindowsAutomatically` — `true`の場合、JavaScript はユーザージェスチャーなしで新しいウィンドウを開けます（デフォルト：`false`）
- `MinimumFontSize` — ポイント単位の最小フォントサイズです。設定するには `optional.NewVar(12.0)` を使用します。未設定の場合、WebKit のデフォルト値が維持されます。
- `ApplicationNameForUserAgent` — WebKit のユーザーエージェント文字列に含まれるアプリケーション名のサフィックスを上書きします。サイトがデフォルトの `"wails.io"` 識別子を拒否する場合（YouTube の埋め込みなど）に便利です。デフォルトを維持するには空のままにします。
- `EnableAutoplayWithoutUserAction` — `true`の場合、音声と動画をユーザージェスチャーなしで自動再生できます。`WKWebViewConfiguration.mediaTypesRequiringUserActionForPlayback = WKAudiovisualMediaTypeNone` に対応します（デフォルト：`false`）

### Windows オプション（ウィンドウ単位）

ウィンドウ単位の構造体は `application.WindowsWindow` です。`WindowsOptions` では<strong>ありません</strong>（こちらは<em>アプリケーション</em>レベルの構造体です）。

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

- ウィンドウのメニューバーを無効にします。`true`の場合、メニューバーが構成されていてもウィンドウには表示されません。
- デフォルト：`false`

**BackdropType**（`BackdropType`）

- `application.Auto` - システムのデフォルト
- `application.None` - 背景効果なし
- `application.Mica` - Mica マテリアル（Windows 11）
- `application.Acrylic` - Acrylic マテリアル（Windows 11）
- `application.Tabbed` - Tabbed マテリアル（Windows 11）

`WindowsBackdropTypeMica` 形式の定数はありません。`application.Mica` などを使用してください。

**CustomTheme**（`ThemeSettings`）

- 値型（ポインターではありません）。ウィンドウの境界線、タイトルバーのテキストと背景、およびメニューバーに使用する、ダークモード／ライトモードのカスタムカラーです。

**DisableFramelessWindowDecorations**（`bool`）

- デフォルトのフレームレス装飾（Aero の影、角の丸み）を無効にします。

**NonClientRegionSupport**（`bool`）

- フレームレスのカスタムタイトルバー向けに、WebView2 ネイティブの `app-region: drag`／`app-region: no-drag` サポートを有効にします。
- これは、単純なネイティブアプリのドラッグ専用です。ネイティブのカスタムキャプションボタン動作や、カスタム最大化ボタンに対する Windows 11 Snap Assist／Snap Layouts は提供しません。

**WebView2CompositionHosting**（`bool`）

- ネイティブの Windows 動作を備えたカスタムキャプションボタン向けに、Wails が管理する `--wails-non-client-region` サポートを有効にします。これには、カスタム最大化ボタンに対する Windows 11 Snap Assist／Snap Layouts も含まれます。
- 試験的な機能です。デフォルトの HWND ホスト型コントローラーの代わりに、`ICoreWebView2CompositionController` と DirectComposition を介して WebView2 をホストします。
- ウィンドウで WebView2 ネイティブの `app-region` サポートと、Wails が管理するカスタムキャプションボタン領域の両方が必要な場合は、`NonClientRegionSupport` と組み合わせて使用できます。

**例：**

```go
Windows: application.WindowsWindow{
    BackdropType: application.Mica,
    DisableIcon:  true,
},
```

**例 - Windows のカスタムタイトルバー領域：**

```go
Windows: application.WindowsWindow{
    NonClientRegionSupport:    true,
    WebView2CompositionHosting: true,
},
```

動作の詳細、トレードオフ、および対応する CSS については、[フレームレスウィンドウ](/features/windows/frameless/#native-non-client-regions-on-windows)を参照してください。

### Linux オプション（ウィンドウ単位）

ウィンドウ単位の構造体は `application.LinuxWindow` です。`LinuxOptions` では<strong>ありません</strong>。

```go
Linux: application.LinuxWindow{
    Icon:                []byte{/* PNG data */},
    WindowIsTranslucent: false,
},
```

**Icon**（`[]byte`）

- ウィンドウアイコン（PNG 形式）。

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

## アプリケーションレベルの Windows オプション

Windows 固有の一部のオプションは、ウィンドウごとではなくアプリケーションレベルで設定する必要があります。これは、WebView2 がユーザーデータパスごとに単一のブラウザー環境を共有するためです。

### ブラウザーフラグ

WebView2 のブラウザーフラグは、アプリケーション内の<strong>すべてのウィンドウ</strong>にわたって、試験的な機能や動作を制御します。これらは `application.Options.Windows` で設定する必要があります：

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
- 利用可能なフラグについては、[WebView2 のブラウザーフラグ](https://learn.microsoft.com/en-us/microsoft-edge/webview2/concepts/webview-features-flags)を参照してください
- 例：`"msWebView2EnableDraggableRegions"`

**DisabledFeatures**（`[]string`）

- 無効にする WebView2 機能フラグの一覧
- Wails は `msSmartScreenProtection` を自動的に無効にします
- 例：`"msExperimentalFeature"`

**AdditionalBrowserArgs**（`[]string`）

- ブラウザープロセスに渡す Chromium のコマンドライン引数
- `--` プレフィックスを含める必要があります（例：`"--remote-debugging-port=9222"`）
- 利用可能な引数については、[Chromium のコマンドラインスイッチ](https://peter.sh/experiments/chromium-command-line-switches/)を参照してください

@note{type="caution" title="重要"}
WebView2 はユーザーデータパスごとに単一のブラウザー環境を共有するため、これらのフラグはすべてのウィンドウにグローバルに適用されます。ウィンドウごとに異なるブラウザーフラグを設定することはできません。

@end

**完全な例：**

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

以下は、本番環境で使用できるウィンドウ設定です：

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

フロントエンドアセットは、ウィンドウごとではなくアプリケーションレベル（`application.Options.Assets`）で配信されます。

## 次のステップ

- [ウィンドウの基本](/features/windows/basics/) - ウィンドウの作成と制御
- [複数ウィンドウ](/features/windows/multiple/) - マルチウィンドウのパターン
- [フレームレスウィンドウ](/features/windows/frameless/) - カスタムウィンドウ装飾
- [ウィンドウイベント](/features/windows/events/) - ライフサイクルイベント

---

**ご質問がありますか？** [Discord](https://discord.gg/JDdSxwjhGf) で質問するか、[サンプル](https://github.com/wailsapp/wails/tree/master/v3/examples)を確認してください。
