---
title: "システムトレイメニュー"
description: "アプリケーションにシステムトレイ（通知領域）との統合を追加します"
slug: "features/menus/systray"
sourcePath: "features/menus/systray.md"
---

## システムトレイメニュー

Wails は、すべてのプラットフォームで動作する**統一されたシステムトレイ API** を提供します。メニュー付きのトレイアイコンを作成し、ウィンドウを関連付け、クリックを処理できます。バックグラウンドアプリケーション、サービス、クイックアクセスユーティリティで、各プラットフォームのネイティブな動作を利用できます。

![macOS のメニューバーから開いた Wails のシステムトレイメニュー](/assets/screenshots/systray-menu-macos.png)

macOS では、Wails のシステムトレイ項目がメニューバーに表示され、ネイティブメニューが開きます。この例には、無効な項目、チェックボックス項目、ラジオ項目、サブメニュー、アクション項目が含まれています。

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

- 黒と無色（透明）のみを使用する
- ダークモードでは黒が白になる
- ファイル名に `Template` サフィックスを付ける：`iconTemplate.png`
- [デザインガイド](https://bjango.com/articles/designingmenubarextras/)

## メニューの追加

システムトレイメニューはアプリケーションメニューと同じように動作します：

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

**すべてのメニュー項目の種類**については、[メニューリファレンス](/features/menus/reference/)を参照してください。

## ウィンドウの関連付け

ウィンドウをトレイアイコンに関連付けると、表示と非表示が自動的に切り替わります：

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
- **トレイアイコンを右クリック** → メニューを設定している場合は表示する
- ウィンドウをトレイアイコンの近くに配置する

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

`HideOnFocusLost` は、Windows、macOS、およびクリックフォーカス方式の Linux デスクトップでトレイポップアップを表示する際に便利です。Wails は、フォーカスがマウスに追従する Linux 環境（一般的な Hyprland、Sway、i3 の構成を含む）では、この動作を無効にします。このような環境では、無効にしないと、ポップアップからポインターを外しただけで、操作する前にポップアップが非表示になる可能性があるためです。これらの環境でも `HideOnEscape` は引き続き利用できます。

上記の左クリックと右クリックの動作は、適切なデフォルト動作です。`OnClick` または `OnRightClick` のハンドラーを明示的に指定すると、対応するデフォルト動作が置き換えられます。プラットフォームの確認方法とエッジケースについては、[手動 systray テストスイート](https://github.com/wailsapp/wails/tree/v3.0.0/v3/test/manual/systray)および[systray ストレステストの例](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/systray-stress)を参照してください。

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

**プラットフォームのサポート状況：**

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
メニューの状態を変更した後は、**必ず `menu.Update()`** を呼び出してください。[メニューリファレンス](/features/menus/reference/#enabled-state)を参照してください。

@end

### メニューの再構築

大幅に変更する場合は、メニュー全体を再構築します：

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

`NSImagePosition` に対応する**アイコンの位置**：

- `application.NSImageLeft` - ラベルの左側にアイコンを表示。
- `application.NSImageRight` - ラベルの右側にアイコンを表示。
- `application.NSImageOnly` - アイコンのみを表示し、ラベルは表示しない。
- `application.NSImageNone` - ラベルのみを表示し、アイコンは表示しない。

**ベストプラクティス：**

- テンプレートアイコン（黒＋透明）を使用する
- ラベルは短くする（3-5 文字）
- Retina ディスプレイでは 18x18～22x22 ピクセルにする
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

- 16x16 または 32x32 ピクセル
- PNG または ICO 形式
- 背景を透明にする

**ツールチップの制限：**

- 最大 127 UTF-16 文字
- これより長いツールチップは切り詰められます
- 最適な使い勝手を得るため、簡潔にしてください

**プラットフォーム機能：**

- Windows Explorer が再起動してもトレイアイコンは維持されます
- Show() メソッドと Hide() メソッドが完全に機能します
- 適切なライフサイクル管理

**ベストプラクティス：**

- 高 DPI ディスプレイには 32x32 を使用してください
- ツールチップは 127 文字未満にしてください
- 複数の Windows バージョンでテストしてください
- 通知領域からあふれる場合も考慮してください
- 条件に応じてトレイの表示／非表示を切り替えるには、Show/Hide を使用してください

[Linux]
**システムトレイとの統合：**

StatusNotifierItem 仕様を使用します（最新の DE の大半）。

```go
systray.SetIcon(iconBytes)
systray.SetLabel("My App")
```

**デスクトップ環境のサポート：**

- **GNOME**：トップバー（拡張機能が必要）
- **KDE Plasma**：システムトレイ
- **XFCE**：通知領域
- **その他**：環境によって異なります

**ベストプラクティス：**

- 22x22 または 24x24 ピクセルを使用してください
- SVG アイコンの方が適切に拡大縮小されます
- 対象のデスクトップ環境でテストしてください
- サポートされていない DE 向けのフォールバックを用意してください

@end

## 完全な例

以下は、本番環境で使用できるシステムトレイアプリケーションです：

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

## 表示の制御

トレイアイコンの表示／非表示を動的に切り替えます：

```go
// Hide tray icon
systray.Hide()

// Show tray icon
systray.Show()
```

`IsVisible()` ゲッターはありません。必要な場合は、アプリケーション独自の状態で表示状態を管理してください。

**プラットフォームのサポート：**

| プラットフォーム | Hide() | Show() | 備考 |
| --- | --- | --- | --- |
| **Windows** | ✅ | ✅ | 完全に機能します。アイコンが通知領域に表示されるか、通知領域から消えます |
| **macOS** | ✅ | ✅ | メニューバー項目の表示／非表示が切り替わります |
| **Linux** | ✅ | ✅ | デスクトップ環境によって異なります |

**使用例：**

- ユーザー設定に基づいてトレイアイコンを一時的に非表示にする
- 必要な場合にのみトレイアイコンが表示されるヘッドレスモード
- アプリケーションの状態に応じて表示／非表示を切り替える

**例 — 条件に応じたトレイの表示：**

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

使用後はトレイアイコンを破棄します：

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

<strong>重要：</strong>リソースを解放するため、シャットダウン時には必ずシステムトレイを破棄してください。

## ベストプラクティス

### ✅ 推奨事項

- **macOS ではテンプレートアイコンを使用する** — ダークモードに適応します
- **ラベルを短くする** — 最大 3-5 文字です
- **Windows ではツールチップを設定する** — ユーザーがアプリを識別しやすくなります
- **すべてのプラットフォームでテストする** — 動作はプラットフォームによって異なります
- **クリックを適切に処理する** — メイン操作には左クリック、メニューには右クリックを使用します
- **状態に応じてアイコンを更新する** — 視覚的なフィードバックは重要です
- **シャットダウン時に破棄する** — リソースを解放します

### ❌ 禁止事項

- **大きなアイコンを使用しない** — プラットフォームのガイドラインに従ってください
- **長いラベルを使用しない** — 途中で切り詰められます
- **ダークモードを考慮する** — Windows と macOS のダークモードでテストしてください
- **クリックハンドラーをブロックしない** — 処理を高速に保ってください
- **menu.Update() の呼び出しを忘れない** — メニューの状態を変更した後に呼び出します
- **トレイがサポートされていると決めつけない** — 一部の Linux DE ではサポートされていません

## トラブルシューティング

### トレイアイコンが表示されない

**考えられる原因：**

1. アイコン形式がサポートされていない
2. アイコンのサイズが大きすぎる、または小さすぎる
3. システムトレイがサポートされていない（Linux）

**解決策：**

`SystemTraySupported()` ヘルパーはありません。代わりにトレイを作成し、プラットフォームを確認したうえで、適切に機能を縮退させてください：

```go
// Probe support: on Linux without a notification-area extension, the tray
// will simply not appear. Defensive code can fall back to window-only mode
// based on runtime.GOOS or after a short timeout if no tray events arrive.
systray := app.SystemTray.New()
systray.SetIcon(iconBytes)
```

### macOS でアイコンが正しく表示されない

<strong>原因：</strong>テンプレートアイコンを使用していません

**解決策：**

```go
// Use template icon
systray.SetTemplateIcon(iconBytes)

// Or design icon as template (black + transparent)
```

### メニューが更新されない

**原因：**`menu.Update()` の呼び出しを忘れています

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
📖 システムトレイの例
完全なシステムトレイアプリケーションを確認します。

[詳細を見る →](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/systray-basic)

@end

---

**ご質問がありますか？**[Discord](https://discord.gg/JDdSxwjhGf) で質問するか、[システムトレイの例](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/systray-basic)を確認してください。
