---
title: "システムトレイメニュー"
description: "アプリケーションにシステムトレイ（通知領域）との統合を追加する"
slug: "features/menus/systray"
sourcePath: "features/menus/systray.md"
---

## システムトレイメニュー

Wails は、すべてのプラットフォームで動作する<strong>統一されたシステムトレイ API</strong> を提供します。メニュー付きのトレイアイコンの作成、ウィンドウの関連付け、クリックの処理が可能で、バックグラウンドアプリケーション、サービス、クイックアクセスユーティリティで各プラットフォームのネイティブな動作を利用できます。

![macOS のメニューバーから開いた Wails のシステムトレイメニュー](/assets/screenshots/systray-menu-macos.png)

macOS では、Wails のシステムトレイ項目がメニューバーに表示され、ネイティブメニューが開きます。この例には、無効化された項目、チェックボックス項目、ラジオ項目、サブメニュー項目、アクション項目が含まれています。

## クイックスタート

```go
package main

import (
    _ "embed"
    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets/icon.png
var icon []byte

func main() {
    app := application.New(application.Options{
        Name: "Tray App",
    })

    // Create system tray
    systray := app.SystemTray.New()
    systray.SetIcon(icon)
    systray.SetLabel("My App")

    // Add menu
    menu := app.NewMenu()
    menu.Add("Show").OnClick(func(ctx *application.Context) {
        // Show main window
    })
    menu.Add("Quit").OnClick(func(ctx *application.Context) {
        app.Quit()
    })
    systray.SetMenu(menu)

    // Create hidden window
    window := app.Window.New()
    window.Hide()

    app.Run()
}
```

<strong>結果：</strong>すべてのプラットフォームでメニュー付きのシステムトレイアイコンが表示されます。

## システムトレイの作成

### 基本的なシステムトレイ

```go
// Create system tray
systray := app.SystemTray.New()

// Set icon
systray.SetIcon(iconBytes)

// Set label (macOS) / tooltip (Windows)
systray.SetLabel("My Application")
```

### アイコンの設定

アイコンは埋め込んでください：

```go
import _ "embed"

//go:embed assets/icon.png
var icon []byte

//go:embed assets/icon-dark.png
var iconDark []byte

func main() {
    app := application.New(application.Options{
        Name: "My App",
    })

    systray := app.SystemTray.New()
    systray.SetIcon(icon)
    systray.SetDarkModeIcon(iconDark)  // Windows and macOS dark mode
    
    app.Run()
}
```

**アイコンの要件：**

| プラットフォーム | サイズ | 形式 | 備考 |
| --- | --- | --- | --- |
| **Windows** | 16x16 または 32x32 | PNG、ICO | 通知領域 |
| **macOS** | 18x18～22x22 | PNG | メニューバー、テンプレートを推奨 |
| **Linux** | 22x22～48x48 | PNG、SVG | デスクトップ環境によって異なる |

### テンプレートアイコン（macOS）

テンプレートアイコンはライトモードとダークモードに自動的に適応します：

```go
systray.SetTemplateIcon(iconBytes)
```

**テンプレートアイコンのガイドライン：**

- 黒と透明色のみを使用する
- ダークモードでは黒が白になる
- ファイル名に`Template`サフィックスを付ける：`iconTemplate.png`
- [デザインガイド](https://bjango.com/articles/designingmenubarextras/)

## メニューの追加

システムトレイメニューはアプリケーションメニューと同様に動作します：

```go
menu := app.NewMenu()

// Add items
menu.Add("Open").OnClick(func(ctx *application.Context) {
    showMainWindow()
})

menu.AddSeparator()

menu.AddCheckbox("Start at Login", false).OnClick(func(ctx *application.Context) {
    enabled := ctx.ClickedMenuItem().Checked()
    setStartAtLogin(enabled)
})

menu.AddSeparator()

menu.Add("Quit").OnClick(func(ctx *application.Context) {
    app.Quit()
})

// Set menu
systray.SetMenu(menu)
```

<strong>すべてのメニュー項目タイプ</strong>については、[メニューリファレンス](/features/menus/reference/)を参照してください。

## ウィンドウの関連付け

ウィンドウをトレイアイコンに関連付けると、表示と非表示を自動的に切り替えられます：

```go
// Create window
window := app.Window.New()

// Attach to tray
systray.AttachWindow(window)

// Configure behaviour — these are setters that return the receiver for chaining.
systray.WindowOffset(10)                          // Pixels from tray icon
systray.WindowDebounce(200 * time.Millisecond)    // Click debounce
```

**動作：**

- ウィンドウは非表示の状態で起動する
- **トレイアイコンを左クリック** → ウィンドウの表示／非表示を切り替える
- **トレイアイコンを右クリック** → 設定されている場合はメニューを表示する
- ウィンドウはトレイアイコンの近くに配置される

**例：ポップアップウィンドウ**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:           "Quick Access",
    Width:           300,
    Height:          400,
    Frameless:       true, // No title bar
    AlwaysOnTop:     true, // Stay on top
    HideOnFocusLost: true, // Dismiss when another window receives focus
    HideOnEscape:    true, // Dismiss when the user presses Escape
})

systray.AttachWindow(window)
systray.WindowOffset(5)
```

`HideOnFocusLost` は、Windows、macOS、およびクリックフォーカス方式の Linux デスクトップでトレイポップアップを表示する場合に便利です。マウス追従フォーカス方式の Linux 環境（一般的な Hyprland、Sway、i3 の構成を含む）では、ポップアップからポインターが外れたときに、使用する前にポップアップが非表示になる可能性があるため、Wails はこの動作を無効にします。`HideOnEscape` は、これらの環境でも引き続き使用できます。

上記の左クリックと右クリックの動作は、状況に応じたデフォルトです。明示的な`OnClick`ハンドラーまたは`OnRightClick`ハンドラーを指定すると、対応するデフォルト動作が置き換えられます。プラットフォームの確認方法とエッジケースについては、[手動 systray テストスイート](https://github.com/wailsapp/wails/tree/master/v3/test/manual/systray)および[systray ストレステストの例](https://github.com/wailsapp/wails/tree/master/v3/examples/systray-stress)を参照してください。

## クリックハンドラー

トレイアイコンのクリックを処理します：

```go
systray := app.SystemTray.New()

// Left click
systray.OnClick(func() {
    fmt.Println("Tray icon clicked")
})

// Right click
systray.OnRightClick(func() {
    fmt.Println("Tray icon right-clicked")
})

// Double click
systray.OnDoubleClick(func() {
    fmt.Println("Tray icon double-clicked")
})

// Mouse enter/leave
systray.OnMouseEnter(func() {
    fmt.Println("Mouse entered tray icon")
})

systray.OnMouseLeave(func() {
    fmt.Println("Mouse left tray icon")
})
```

**プラットフォームの対応状況：**

| イベント | Windows | macOS | Linux |
| --- | --- | --- | --- |
| OnClick | ✅ | ✅ | ✅ |
| OnRightClick | ✅ | ✅ | ✅ |
| OnDoubleClick | ✅ | ✅ | ⚠️ 環境によって異なる |
| OnMouseEnter | ✅ | ✅ | ⚠️ 環境によって異なる |
| OnMouseLeave | ✅ | ✅ | ⚠️ 環境によって異なる |

## 動的な更新

トレイアイコンとメニューを動的に更新します：

### アイコンの変更

```go
var isActive bool

func updateTrayIcon() {
    if isActive {
        systray.SetIcon(activeIcon)
        systray.SetLabel("Active")
    } else {
        systray.SetIcon(inactiveIcon)
        systray.SetLabel("Inactive")
    }
}
```

### メニューの更新

```go
var isPaused bool

pauseMenuItem := menu.Add("Pause")

pauseMenuItem.OnClick(func(ctx *application.Context) {
    isPaused = !isPaused
    
    if isPaused {
        pauseMenuItem.SetLabel("Resume")
    } else {
        pauseMenuItem.SetLabel("Pause")
    }
    
    menu.Update()  // Important!
})
```

@note{type="caution" title="必ず Update() を呼び出す"}
メニューの状態を変更した後は、<strong>`menu.Update()`</strong>を呼び出してください。[メニューリファレンス](/features/menus/reference/#enabled-state)を参照してください。

@end

### メニューの再構築

大規模な変更を行う場合は、メニュー全体を再構築します。

```go
func rebuildTrayMenu(status string) {
    menu := app.NewMenu()
    
    // Status-specific items
    switch status {
    case "syncing":
        menu.Add("Syncing...").SetEnabled(false)
        menu.Add("Pause Sync").OnClick(pauseSync)
    case "synced":
        menu.Add("Up to date ✓").SetEnabled(false)
        menu.Add("Sync Now").OnClick(startSync)
    case "error":
        menu.Add("Sync Error").SetEnabled(false)
        menu.Add("Retry").OnClick(retrySync)
    }
    
    menu.AddSeparator()
    menu.Add("Quit").OnClick(func(ctx *application.Context) {
        app.Quit()
    })
    
    systray.SetMenu(menu)
}
```

## プラットフォーム固有の機能

@tabs{sync-key="platform"}
[macOS]
**メニューバーとの統合：**

```go
// Set label (appears next to icon)
systray.SetLabel("My App")

// Use template icon (adapts to dark mode)
systray.SetTemplateIcon(iconBytes)

// Set icon position — uses AppKit NSImage placement constants.
systray.SetIconPosition(application.NSImageRight)
```

**アイコンの位置**（`NSImagePosition`に対応）：

- `application.NSImageLeft` - ラベルの左側にアイコンを表示。
- `application.NSImageRight` - ラベルの右側にアイコンを表示。
- `application.NSImageOnly` - アイコンのみを表示し、ラベルは表示しない。
- `application.NSImageNone` - ラベルのみを表示し、アイコンは表示しない。

**ベストプラクティス：**

- テンプレートアイコン（黒色＋透明）を使用する
- ラベルは短くする（3-5文字）
- Retinaディスプレイでは18x18～22x22ピクセルにする
- ライトモードとダークモードの両方でテストする

[Windows]
**通知領域との統合：**

```go
// Set tooltip (appears on hover)
systray.SetTooltip("My Application")

// Or use SetLabel (same as tooltip on Windows)
systray.SetLabel("My Application")

// Show/Hide functionality (fully functional)
systray.Show()  // Show tray icon
systray.Hide()  // Hide tray icon
```

**アイコンの要件：**

- 16x16または32x32ピクセル
- PNGまたはICO形式
- 透明な背景

**ツールチップの制限：**

- 最大127文字（UTF-16）
- これより長いツールチップは切り詰められる
- 最適な操作性を得るため、簡潔にする

**プラットフォームの機能：**

- Windows Explorerを再起動してもトレイアイコンが維持される
- Show()メソッドとHide()メソッドが完全に機能する
- 適切なライフサイクル管理

**ベストプラクティス：**

- 高DPIディスプレイでは32x32を使用する
- ツールチップは127文字未満にする
- 異なるバージョンのWindowsでテストする
- 通知領域のオーバーフローを考慮する
- トレイの表示を条件に応じて切り替えるにはShow/Hideを使用する

[Linux]
**システムトレイとの統合：**

StatusNotifierItem仕様を使用します（最新のデスクトップ環境の大半）。

```go
systray.SetIcon(iconBytes)
systray.SetLabel("My App")
```

**デスクトップ環境のサポート：**

- **GNOME**：トップバー（拡張機能が必要）
- **KDE Plasma**：システムトレイ
- **XFCE**：通知領域
- **その他**：環境によって異なる

**ベストプラクティス：**

- 22x22または24x24ピクセルを使用する
- SVGアイコンのほうが適切に拡大縮小される
- 対象のデスクトップ環境でテストする
- サポートされていないデスクトップ環境向けのフォールバックを用意する

@end

## 完全な例

以下は、本番環境に対応したシステムトレイアプリケーションです。

```go
package main

import (
    _ "embed"
    "fmt"
    "time"
    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets/icon.png
var icon []byte

//go:embed assets/icon-active.png
var iconActive []byte

type TrayApp struct {
    app     *application.App
    systray *application.SystemTray
    window  *application.WebviewWindow
    menu    *application.Menu
    isActive bool
}

func main() {
    app := application.New(application.Options{
        Name: "Tray Application",
        Mac: application.MacOptions{
            ApplicationShouldTerminateAfterLastWindowClosed: false,
        },
    })

    trayApp := &TrayApp{app: app}
    trayApp.setup()

    app.Run()
}

func (t *TrayApp) setup() {
    // Create system tray
    t.systray = t.app.SystemTray.New()
    t.systray.SetIcon(icon)
    t.systray.SetLabel("Inactive")
    
    // Create menu
    t.createMenu()
    
    // Create window (hidden by default)
    t.window = t.app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:  "Tray Application",
        Width:  400,
        Height: 600,
        Hidden: true,
    })
    
    // Attach window to tray
    t.systray.AttachWindow(t.window)
    t.systray.WindowOffset(10)
    
    // Handle tray clicks
    t.systray.OnRightClick(func() {
        t.systray.OpenMenu()
    })
    
    // Start background task
    go t.backgroundTask()
}

func (t *TrayApp) createMenu() {
    t.menu = t.app.NewMenu()
    
    // Status item (disabled)
    statusItem := t.menu.Add("Status: Inactive")
    statusItem.SetEnabled(false)
    
    t.menu.AddSeparator()
    
    // Toggle active
    t.menu.Add("Start").OnClick(func(ctx *application.Context) {
        t.toggleActive()
    })
    
    // Show window
    t.menu.Add("Show Window").OnClick(func(ctx *application.Context) {
        t.window.Show()
        t.window.Focus()
    })
    
    t.menu.AddSeparator()
    
    // Settings
    t.menu.AddCheckbox("Start at Login", false).OnClick(func(ctx *application.Context) {
        enabled := ctx.ClickedMenuItem().Checked()
        t.setStartAtLogin(enabled)
    })
    
    t.menu.AddSeparator()
    
    // Quit
    t.menu.Add("Quit").OnClick(func(ctx *application.Context) {
        t.app.Quit()
    })
    
    t.systray.SetMenu(t.menu)
}

func (t *TrayApp) toggleActive() {
    t.isActive = !t.isActive
    t.updateTray()
}

func (t *TrayApp) updateTray() {
    if t.isActive {
        t.systray.SetIcon(iconActive)
        t.systray.SetLabel("Active")
    } else {
        t.systray.SetIcon(icon)
        t.systray.SetLabel("Inactive")
    }
    
    // Rebuild menu with new status
    t.createMenu()
}

func (t *TrayApp) backgroundTask() {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        if t.isActive {
            fmt.Println("Background task running...")
            // Do work
        }
    }
}

func (t *TrayApp) setStartAtLogin(enabled bool) {
    // Implementation varies by platform
    fmt.Printf("Start at login: %v\n", enabled)
}
```

## 表示制御

トレイアイコンの表示と非表示を動的に切り替えます。

```go
// Hide tray icon
systray.Hide()

// Show tray icon
systray.Show()
```

`IsVisible()`ゲッターはありません。必要な場合は、アプリケーション独自の状態で表示状態を管理してください。

**プラットフォームのサポート：**

| プラットフォーム | Hide() | Show() | 備考 |
| --- | --- | --- | --- |
| **Windows** | ✅ | ✅ | 完全に機能する - 通知領域でアイコンが表示／非表示になる |
| **macOS** | ✅ | ✅ | メニューバー項目が表示／非表示になる |
| **Linux** | ✅ | ✅ | デスクトップ環境によって異なる |

**ユースケース：**

- ユーザー設定に基づいてトレイアイコンを一時的に非表示にする
- 必要な場合にのみトレイアイコンを表示するヘッドレスモード
- アプリケーションの状態に基づいて表示／非表示を切り替える

**例 - 条件に応じたトレイの表示制御：**

```go
func (t *TrayApp) setTrayVisibility(visible bool) {
    if visible {
        t.systray.Show()
    } else {
        t.systray.Hide()
    }
}

// Show tray only when updates are available
func (t *TrayApp) checkForUpdates() {
    if hasUpdates {
        t.systray.Show()
        t.systray.SetLabel("Update Available")
    } else {
        t.systray.Hide()
    }
}
```

## クリーンアップ

使用後はトレイアイコンを破棄します。

```go
// In OnShutdown
app := application.New(application.Options{
    OnShutdown: func() {
        if systray != nil {
            systray.Destroy()
        }
    },
})
```

**重要：** リソースを解放するため、終了時には必ずシステムトレイを破棄してください。

## ベストプラクティス

### ✅ 推奨事項

- **macOSではテンプレートアイコンを使用する** - ダークモードに適応します
- **ラベルを短くする** - 最大3-5文字
- **Windowsではツールチップを表示する** - ユーザーがアプリを識別しやすくなります
- **すべてのプラットフォームでテストする** - 動作はプラットフォームによって異なります
- **クリックを適切に処理する** - 左クリックにはメイン操作、右クリックにはメニューを割り当てます
- **状態に応じてアイコンを更新する** - 視覚的なフィードバックは重要です
- **終了時に破棄する** - リソースを解放します

### ❌ 避けるべきこと

- **大きなアイコンを使用しない** - プラットフォームのガイドラインに従ってください
- **長いラベルを使用しない** - 途中で切り詰められます
- **ダークモードへの対応を忘れない** - WindowsとmacOSのダークモードでテストしてください
- **クリックハンドラーをブロックしない** - 処理を高速に保ってください
- **menu.Update()の呼び出しを忘れない** - メニューの状態を変更した後に呼び出してください
- **システムトレイがサポートされていると決めつけない** - 一部のLinuxデスクトップ環境ではサポートされていません

## トラブルシューティング

### トレイアイコンが表示されない

**考えられる原因：**

1. アイコン形式がサポートされていない
2. アイコンのサイズが大きすぎる、または小さすぎる
3. システムトレイがサポートされていない（Linux）

**解決策：**

`SystemTraySupported()`ヘルパーはありません。代わりにトレイを作成し、プラットフォームを確認して、適切に機能を縮退させてください：

```go
// Probe support: on Linux without a notification-area extension, the tray
// will simply not appear. Defensive code can fall back to window-only mode
// based on runtime.GOOS or after a short timeout if no tray events arrive.
systray := app.SystemTray.New()
systray.SetIcon(iconBytes)
```

### macOSでアイコンが正しく表示されない

**原因：** テンプレートアイコンを使用していない

**解決策：**

```go
// Use template icon
systray.SetTemplateIcon(iconBytes)

// Or design icon as template (black + transparent)
```

### メニューが更新されない

**原因：** `menu.Update()`の呼び出しを忘れている

**解決策：**

```go
menuItem.SetLabel("New Label")
menu.Update()  // Add this!
```

## 次のステップ

@cards{cols="2"}
📖 メニューリファレンス
メニュー項目の種類とプロパティに関する完全なリファレンスです。

[詳細を見る →](/features/menus/reference/)

---
☰ アプリケーションメニュー
アプリケーションのメニューバーを作成します。

[詳細を見る →](/features/menus/application/)

---
◆ コンテキストメニュー
右クリックで表示するコンテキストメニューを作成します。

[詳細を見る →](/features/menus/context/)

---
📖 システムトレイのサンプル
完全なシステムトレイアプリケーションを確認します。

[詳細を見る →](https://github.com/wailsapp/wails/tree/master/v3/examples/systray-basic)

@end

---

**ご質問がありますか？** [Discord](https://discord.gg/JDdSxwjhGf)で質問するか、[システムトレイのサンプル](https://github.com/wailsapp/wails/tree/master/v3/examples/systray-basic)を確認してください。
