---
title: "アプリケーションメニュー"
description: "デスクトップアプリケーション用のネイティブメニューバーを作成する"
slug: "features/menus/application"
sourcePath: "features/menus/application.md"
---

## 課題

本格的なデスクトップアプリケーションには、ファイル、編集、表示、ヘルプなどのメニューバーが必要です。しかし、メニューの動作はプラットフォームごとに異なります。

- **macOS**：画面上部のグローバルメニューバー
- **Windows**：ウィンドウのタイトルバーにあるメニューバー
- **Linux**：デスクトップ環境によって異なる

各プラットフォームに適したメニューを手作業で構築するのは煩雑で、エラーも起きやすくなります。

## Wailsによる解決策

Wailsは、各プラットフォームのネイティブメニューを自動的に作成する<strong>統一API</strong>を提供します。一度記述すれば、すべてのプラットフォームでネイティブな動作が得られます。

![標準項目、チェックボックス項目、ラジオ項目、サブメニュー項目を含むmacOS上のWailsアプリケーションメニュー](/assets/screenshots/application-menu-macos.png)

macOSでは、アプリケーションメニューはグローバルメニューバーに配置されます。このキャプチャは、無効な項目、チェックボックス項目、ラジオ項目、サブメニュー項目を含め、WailsのメニューAPIによってレンダリングされたネイティブメニューを示しています。

## クイックスタート

```go
package main

import (
    "runtime"
    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "My App",
    })

    // Create menu
    menu := app.NewMenu()

    // Add standard menus (platform-appropriate)
    if runtime.GOOS == "darwin" {
        menu.AddRole(application.AppMenu)  // macOS only
    }
    menu.AddRole(application.FileMenu)
    menu.AddRole(application.EditMenu)
    menu.AddRole(application.WindowMenu)
    menu.AddRole(application.HelpMenu)

    // Set the application menu
    app.Menu.Set(menu)

    // Create window with UseApplicationMenu to inherit the menu on Windows/Linux
    app.Window.NewWithOptions(application.WebviewWindowOptions{
        UseApplicationMenu: true,
    })

    app.Run()
}
```

**これだけです！** 標準項目を備えた、各プラットフォームのネイティブメニューが利用できるようになりました。`UseApplicationMenu`オプションにより、追加のコードなしでWindowsとLinuxのウィンドウにメニューが表示されます。

## メニューの作成

### 基本的なメニューの作成

```go
// Create a new menu
menu := app.NewMenu()

// Add a top-level menu
fileMenu := menu.AddSubmenu("File")

// Add menu items
fileMenu.Add("New").OnClick(func(ctx *application.Context) {
    // Handle New
})

fileMenu.Add("Open").OnClick(func(ctx *application.Context) {
    // Handle Open
})

fileMenu.AddSeparator()

fileMenu.Add("Quit").OnClick(func(ctx *application.Context) {
    app.Quit()
})
```

### メニューの設定

**推奨される方法** — プラットフォーム間で一貫した動作を実現するには、`UseApplicationMenu`を使用します。

```go
// Set the application menu once
app.Menu.Set(menu)

// Create windows that inherit the menu on Windows/Linux
app.Window.NewWithOptions(application.WebviewWindowOptions{
    UseApplicationMenu: true,  // Window uses the app menu
})
```

この方法では、次のように動作します。

- **macOS**：メニューは画面上部に表示されます（標準動作）。
- **Windows/Linux**：`UseApplicationMenu: true`が設定された各ウィンドウにアプリケーションメニューが表示されます。

**プラットフォーム固有の詳細：**

@tabs{sync-key="platform"}
[macOS]
**グローバルメニューバー**（アプリケーションごとに1つ）：

```go
app.Menu.Set(menu)
```

メニューは画面上部に表示され、すべてのウィンドウを閉じた後も表示され続けます。macOSではすべてのアプリがグローバルメニューを使用するため、`UseApplicationMenu`オプションは効果がありません。

[Windows]
**ウィンドウごとのメニューバー**：

```go
// Option 1: Use application menu (recommended)
app.Menu.Set(menu)
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    UseApplicationMenu: true,
})

// Option 2: Set menu directly on window
window.SetMenu(menu)
```

各ウィンドウには独自のメニューを設定することも、アプリケーションメニューを継承させることもできます。メニューはウィンドウのタイトルバーに表示されます。

[Linux]
**ウィンドウごとのメニューバー**（通常）：

```go
// Option 1: Use application menu (recommended)
app.Menu.Set(menu)
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    UseApplicationMenu: true,
})

// Option 2: Set menu directly on window
window.SetMenu(menu)
```

動作はデスクトップ環境によって異なります。Unityなど、一部の環境はグローバルメニューに対応しています。

@end

@note{type="tip" title="クロスプラットフォームメニューの簡素化"}
`UseApplicationMenu: true`を使用すると、次のようなプラットフォーム固有のコードは不要になります。

```go
// Old approach - no longer needed
if runtime.GOOS == "darwin" {
    app.Menu.Set(menu)
} else {
    window.SetMenu(menu)
}
```

@end

**ウィンドウごとのカスタムメニュー：**

ウィンドウにアプリケーションメニューとは異なるメニューが必要な場合は、そのウィンドウに直接設定します。

```go
window.SetMenu(customMenu)  // Overrides UseApplicationMenu
```

## メニューのロール

Wailsは、各プラットフォームに適したメニュー構造を自動的に作成する<strong>定義済みのメニューロール</strong>を提供します。

### 利用可能なロール

| ロール | 説明 | プラットフォームに関する注記 |
| --- | --- | --- |
| `AppMenu` | 「このアプリケーションについて」、「環境設定」、「終了」を含むアプリケーションメニュー | **macOSのみ** |
| `FileMenu` | ファイル操作（新規、開く、保存など） | すべてのプラットフォーム |
| `EditMenu` | テキスト編集（取り消す、やり直す、カット、コピー、ペースト） | すべてのプラットフォーム |
| `WindowMenu` | ウィンドウ管理（しまう、拡大／縮小など） | すべてのプラットフォーム |
| `HelpMenu` | ヘルプと情報 | すべてのプラットフォーム |

### ロールの使用

```go
menu := app.NewMenu()

// macOS: Add application menu
if runtime.GOOS == "darwin" {
    menu.AddRole(application.AppMenu)
}

// All platforms: Add standard menus
menu.AddRole(application.FileMenu)
menu.AddRole(application.EditMenu)
menu.AddRole(application.WindowMenu)
menu.AddRole(application.HelpMenu)
```

**作成される内容：**

@tabs{sync-key="platform"}
[macOS]
**AppMenu**（アプリ名を含む）：

- [App Name]について
- 環境設定... (⌘,)
- ---
- サービス
- ---
- [App Name]を隠す (⌘H)
- ほかを隠す (⌥⌘H)
- すべてを表示
- ---
- [App Name]を終了 (⌘Q)

**FileMenu**：

- 新規 (⌘N)
- 開く... (⌘O)
- ---
- ウィンドウを閉じる (⌘W)

**EditMenu**：

- 取り消す (⌘Z)
- やり直す (⇧⌘Z)
- ---
- カット（⌘X）
- コピー（⌘C）
- ペースト（⌘V）
- すべてを選択（⌘A）

**WindowMenu**：

- 最小化（⌘M）
- 拡大／縮小
- ---
- すべてを手前に移動

**HelpMenu**：

- [App Name] ヘルプ

[Windows]
**FileMenu**：

- 新規（Ctrl+N）
- 開く...（Ctrl+O）
- ---
- 終了（Alt+F4）

**EditMenu**：

- 元に戻す（Ctrl+Z）
- やり直す（Ctrl+Y）
- ---
- 切り取り（Ctrl+X）
- コピー（Ctrl+C）
- 貼り付け（Ctrl+V）
- すべて選択（Ctrl+A）

**WindowMenu**：

- 最小化
- 最大化

**HelpMenu**：

- [App Name] について

[Linux]
Windows と似ていますが、キーボードショートカットはデスクトップ環境によって異なる場合があります。

@end

### ロールメニューのカスタマイズ

`Menu.AddRole(role)` が返すのは、ロールのサブメニューでは<strong>なく</strong>、<strong>レシーバー</strong>のメニュー（トップレベルメニュー）です。ロールのサブメニューに項目を追加するには、挿入されたロール項目を `FindByRole` で検索し、その項目に対して `GetSubmenu()` を呼び出します：

```go
menu.AddRole(application.FileMenu)

fileMenu := menu.FindByRole(application.FileMenu).GetSubmenu()
fileMenu.Add("Import...").OnClick(handleImport)
fileMenu.Add("Export...").OnClick(handleExport)
```

## カスタムメニュー

アプリケーション固有の機能用に独自のメニューを作成します：

```go
// Add a custom top-level menu
toolsMenu := menu.AddSubmenu("Tools")

// Add items
toolsMenu.Add("Settings").OnClick(func(ctx *application.Context) {
    showSettingsWindow()
})

toolsMenu.AddSeparator()

// Add checkbox
toolsMenu.AddCheckbox("Dark Mode", false).OnClick(func(ctx *application.Context) {
    isDark := ctx.ClickedMenuItem().Checked()
    setTheme(isDark)
})

// Add radio group
toolsMenu.AddRadio("Small", true).OnClick(handleFontSize)
toolsMenu.AddRadio("Medium", false).OnClick(handleFontSize)
toolsMenu.AddRadio("Large", false).OnClick(handleFontSize)

// Add submenu
advancedMenu := toolsMenu.AddSubmenu("Advanced")
advancedMenu.Add("Configure...").OnClick(showAdvancedSettings)
```

<strong>その他のメニュー項目タイプ</strong>については、[メニューリファレンス](/features/menus/reference/)を参照してください。

## 動的メニュー

アプリケーションの状態に応じてメニューを更新します：

### 項目の有効化／無効化

```go
var saveMenuItem *application.MenuItem

func createMenu() {
    menu := app.NewMenu()
    fileMenu := menu.AddSubmenu("File")
    
    saveMenuItem = fileMenu.Add("Save")
    saveMenuItem.SetEnabled(false)  // Initially disabled
    saveMenuItem.OnClick(handleSave)
    
    app.Menu.Set(menu)
}

func onDocumentChanged() {
    saveMenuItem.SetEnabled(hasUnsavedChanges())
    menu.Update()  // Important!
}
```

@note{type="caution" title="必ず menu.Update() を呼び出す"}
メニューの状態（有効化／無効化、ラベル、チェック状態）を変更した後は、**必ず `menu.Update()`** を呼び出してください。メニューが再構築される Windows では特に重要です。

詳細については、[メニューリファレンス](/features/menus/reference/#enabled-state)を参照してください。

@end

### ラベルの変更

```go
updateMenuItem := menu.Add("Check for Updates")

updateMenuItem.OnClick(func(ctx *application.Context) {
    updateMenuItem.SetLabel("Checking...")
    menu.Update()
    
    checkForUpdates()
    
    updateMenuItem.SetLabel("Check for Updates")
    menu.Update()
})
```

### メニューの再構築

大幅な変更を行う場合は、メニュー全体を再構築します：

```go
func rebuildFileMenu() {
    menu := app.NewMenu()
    fileMenu := menu.AddSubmenu("File")
    
    fileMenu.Add("New").OnClick(handleNew)
    fileMenu.Add("Open").OnClick(handleOpen)
    
    // Add recent files dynamically
    if hasRecentFiles() {
        recentMenu := fileMenu.AddSubmenu("Open Recent")
        for _, file := range getRecentFiles() {
            filePath := file  // Capture for closure
            recentMenu.Add(filepath.Base(file)).OnClick(func(ctx *application.Context) {
                openFile(filePath)
            })
        }
        recentMenu.AddSeparator()
        recentMenu.Add("Clear Recent").OnClick(clearRecentFiles)
    }
    
    fileMenu.AddSeparator()
    fileMenu.Add("Quit").OnClick(func(ctx *application.Context) {
        app.Quit()
    })
    
    app.Menu.Set(menu)
}
```

## メニューからのウィンドウ制御

メニュー項目からウィンドウを制御できます：

```go
viewMenu := menu.AddSubmenu("View")

// Toggle fullscreen
viewMenu.Add("Toggle Fullscreen").OnClick(func(ctx *application.Context) {
    if window, ok := app.Window.GetByName("main"); ok {
        window.ToggleFullscreen()
    }
})

// Zoom controls
viewMenu.Add("Zoom In").SetAccelerator("CmdOrCtrl++").OnClick(func(ctx *application.Context) {
    // Increase zoom
})

viewMenu.Add("Zoom Out").SetAccelerator("CmdOrCtrl+-").OnClick(func(ctx *application.Context) {
    // Decrease zoom
})

viewMenu.Add("Reset Zoom").SetAccelerator("CmdOrCtrl+0").OnClick(func(ctx *application.Context) {
    // Reset zoom
})
```

**アクティブウィンドウを取得：**

```go
menuItem.OnClick(func(ctx *application.Context) {
    window := application.Get().Window.Current() // the window the menu was invoked from
    // Use window
})
```

## プラットフォーム固有の考慮事項

### macOS

**メニューバーの動作：**

- <strong>画面上部</strong>に表示される（グローバル）
- すべてのウィンドウを閉じても表示され続ける
- 最初のメニューは<strong>常にアプリケーションメニュー</strong>になる
- 標準項目には `menu.AddRole(application.AppMenu)` を使用する

**標準的な配置場所：**

- **このアプリケーションについて**：アプリケーションメニュー
- **環境設定**：アプリケーションメニュー（⌘,）
- **終了**：アプリケーションメニュー（⌘Q）
- **ヘルプ**：ヘルプメニュー

**例：**

```go
if runtime.GOOS == "darwin" {
    menu.AddRole(application.AppMenu)  // Adds About, Preferences, Quit
    
    // Don't add Quit to File menu on macOS
    // Don't add About to Help menu on macOS
}
```

### Windows

**メニューバーの動作：**

- <strong>ウィンドウのタイトルバー</strong>に表示される
- 各ウィンドウに固有のメニューがある
- アプリケーションメニューはない

**標準的な配置場所：**

- **終了**：ファイルメニュー（Alt+F4）
- **設定**：ツールメニューまたは編集メニュー
- **バージョン情報**：ヘルプメニュー

**例：**

```go
if runtime.GOOS == "windows" {
    menu.AddRole(application.FileMenu) // Exit is added automatically
    menu.AddRole(application.HelpMenu) // About is added automatically
}
```

### Linux

**メニューバーの動作：**

- 通常はウィンドウごとに表示される（Windows と同様）
- 一部のデスクトップ環境では、グローバルメニューがサポートされています（Unity、拡張機能を導入した GNOME）
- 外観はデスクトップ環境によって異なります

<strong>ベストプラクティス：</strong>Windows の慣例に従い、対象のデスクトップ環境でテストしてください。

## 完全な例

以下は、本番環境で使用できるメニュー構成です。

```go
package main

import (
    "runtime"
    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "My Application",
    })

    // Create and set menu
    createMenu(app)

    // Create main window with UseApplicationMenu for cross-platform menu support
    app.Window.NewWithOptions(application.WebviewWindowOptions{
        UseApplicationMenu: true,
    })

    app.Run()
}

func createMenu(app *application.App) {
    menu := app.NewMenu()

    // Platform-specific application menu (macOS only)
    if runtime.GOOS == "darwin" {
        menu.AddRole(application.AppMenu)
    }

    // File menu — AddRole returns the receiver menu, not the role submenu.
    // To add items into the File submenu, look it up via FindByRole + GetSubmenu.
    menu.AddRole(application.FileMenu)
    fileMenu := menu.FindByRole(application.FileMenu).GetSubmenu()
    fileMenu.Add("Import...").SetAccelerator("CmdOrCtrl+I").OnClick(handleImport)
    fileMenu.Add("Export...").SetAccelerator("CmdOrCtrl+E").OnClick(handleExport)

    // Edit menu
    menu.AddRole(application.EditMenu)

    // View menu
    viewMenu := menu.AddSubmenu("View")
    viewMenu.Add("Toggle Fullscreen").SetAccelerator("F11").OnClick(toggleFullscreen)
    viewMenu.AddSeparator()
    viewMenu.AddCheckbox("Show Sidebar", true).OnClick(toggleSidebar)
    viewMenu.AddCheckbox("Show Toolbar", true).OnClick(toggleToolbar)

    // Tools menu
    toolsMenu := menu.AddSubmenu("Tools")
    
    // Settings location varies by platform
    if runtime.GOOS == "darwin" {
        // On macOS, Preferences is in Application menu (added by AppMenu role)
    } else {
        toolsMenu.Add("Settings").SetAccelerator("CmdOrCtrl+,").OnClick(showSettings)
    }
    
    toolsMenu.AddSeparator()
    toolsMenu.AddCheckbox("Dark Mode", false).OnClick(toggleDarkMode)

    // Window menu
    menu.AddRole(application.WindowMenu)

    // Help menu
    helpMenu := menu.AddRole(application.HelpMenu)
    helpMenu.Add("Documentation").OnClick(openDocumentation)
    
    // About location varies by platform
    if runtime.GOOS == "darwin" {
        // On macOS, About is in Application menu (added by AppMenu role)
    } else {
        helpMenu.AddSeparator()
        helpMenu.Add("About").OnClick(showAbout)
    }

    // Set the application menu
    app.Menu.Set(menu)
}

func handleImport(ctx *application.Context) {
    // Implementation
}

func handleExport(ctx *application.Context) {
    // Implementation
}

func toggleFullscreen(ctx *application.Context) {
    window := application.Get().Window.Current()
    window.ToggleFullscreen()
}

func toggleSidebar(ctx *application.Context) {
    // Implementation
}

func toggleToolbar(ctx *application.Context) {
    // Implementation
}

func showSettings(ctx *application.Context) {
    // Implementation
}

func toggleDarkMode(ctx *application.Context) {
    isDark := ctx.ClickedMenuItem().Checked()
    // Apply theme
}

func openDocumentation(ctx *application.Context) {
    // Open browser
}

func showAbout(ctx *application.Context) {
    // Show about dialog
}
```

## ベストプラクティス

### ✅ 推奨事項

- 標準メニュー（ファイル、編集など）には<strong>メニューロールを使用する</strong>
- メニュー構成は<strong>各プラットフォームの慣例に従う</strong>
- よく使う操作には<strong>キーボードショートカットを追加する</strong>
- メニューの状態を変更した後は<strong>menu.Update() を呼び出す</strong>
- 動作が異なるため、**すべてのプラットフォームでテストする**
- **メニュー階層を浅く保つ** — 最大2-3階層
- **明確なラベルを使用する** — 「保存」ではなく「プロジェクトを保存」

### ❌ 非推奨事項

- **プラットフォーム固有のショートカットをハードコードしない** — `CmdOrCtrl`を使用してください
- **macOS では「終了」を「ファイル」メニューに配置しない** — 「アプリケーション」メニューに配置されます
- **macOS では「このアプリケーションについて」を「ヘルプ」メニューに配置しない** — 「アプリケーション」メニューに配置されます
- **menu.Update() の呼び出しを忘れない** — メニューが正しく動作しなくなります
- **階層を深くしすぎない** — ユーザーが迷ってしまいます
- **専門用語を使用しない** — ユーザーに分かりやすいラベルにしてください

## 次のステップ

@cards{cols="2"}
📖 メニューのリファレンス
メニュー項目の種類とプロパティに関する完全なリファレンスです。

[詳細を見る →](/features/menus/reference/)

---
◆ コンテキストメニュー
右クリックで表示するコンテキストメニューを作成します。

[詳細を見る →](/features/menus/context/)

---
★ システムトレイメニュー
システムトレイ／メニューバーとの連携を追加します。

[詳細を見る →](/features/menus/systray/)

---
📖 メニューパターン
一般的なメニューパターンとベストプラクティスです。

[詳細を見る →](/guides/menus/)

@end

---

**ご質問がありますか？**[Discord](https://discord.gg/JDdSxwjhGf) で質問するか、[メニューのサンプル](https://github.com/wailsapp/wails/tree/master/v3/examples/menu)を確認してください。
