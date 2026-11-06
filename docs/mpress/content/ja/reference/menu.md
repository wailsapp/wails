---
title: "メニュー API"
description: "メニュー API の完全なリファレンス"
slug: "reference/menu"
sourcePath: "reference/menu.md"
---

## 概要

メニュー API には、アプリケーションメニュー、コンテキストメニュー、システムトレイメニューを作成および管理するためのメソッドが用意されています。

**メニューの種類：**

- **アプリケーションメニュー** - 上部のメニューバー（ファイル、編集など）
- **コンテキストメニュー** - 右クリックメニュー
- **システムトレイメニュー** - システムトレイ／通知領域のメニュー

## メニューの作成

### NewMenu()

新しいメニューを作成します。

```go
func (a *App) NewMenu() *Menu
```

**例：**

```go
menu := app.NewMenu()
```

## メニューのメソッド

### Add()

メニューにメニュー項目を追加します。

```go
func (m *Menu) Add(label string) *MenuItem
```

**パラメーター：**

- `label` - メニュー項目に表示するテキスト

**戻り値：** 作成されたメニュー項目

**例：**

```go
item := menu.Add("Open File")
item.OnClick(func(ctx *application.Context) {
    // Handle click
})
```

### AddSubmenu()

メニューにサブメニューを追加します。

```go
func (m *Menu) AddSubmenu(label string) *Menu
```

**パラメーター：**

- `label` - サブメニューのラベル

**戻り値：** 作成されたサブメニュー

**例：**

```go
fileMenu := menu.AddSubmenu("File")
fileMenu.Add("New")
fileMenu.Add("Open")
fileMenu.Add("Save")
```

### AddSeparator()

メニュー項目の間に視覚的な区切り線を追加します。

```go
func (m *Menu) AddSeparator()
```

**例：**

```go
menu.Add("Copy")
menu.Add("Paste")
menu.AddSeparator()
menu.Add("Select All")
```

**ベストプラクティス：** 関連するメニュー項目をまとめるには、区切り線を使用します。

### AddSectionHeader()

**macOS 14 以降では、ネイティブの操作不可なセクション見出しを追加します**。
それより古いバージョンの macOS およびその他すべてのプラットフォーム（サーバービルドを含む）では、
何も行われず、プレースホルダー項目も挿入されません。

```go
func (m *Menu) AddSectionHeader(label string)
```

トレイメニュー、コンテキストメニュー、またはアプリケーションメニュー内のグループにラベルを付けるには、見出しを使用します。

```go
menu.AddSectionHeader("mac-node1")
menu.Add("Windows VM").OnClick(openWindowsVM)
menu.Add("Linux VM").OnClick(openLinuxVM)
menu.AddSectionHeader("ai-master")
menu.Add("Build VM").OnClick(openBuildVM)
```

見出しには、AppKit 標準のセクションヘッダー表示が使用されます。すでに表示されているメニューを再構築した後は、ほかの構造変更の場合と同様に
`menu.Update()`を呼び出してください。`os.Hostname()`から取得した見出しについては、
`systray-menu`の例を参照してください。

省略された見出しは、ラジオグループを分割しません。すべてのプラットフォームでラジオ
グループを分離しておく必要がある場合は、`AddSeparator()`を使用してください。

### AddCheckbox()

チェック可能なメニュー項目を追加します。

```go
func (m *Menu) AddCheckbox(label string, checked bool) *MenuItem
```

**パラメーター：**

- `label` - チェックボックスのラベル
- `checked` - チェック状態の初期値

**例：**

```go
darkMode := menu.AddCheckbox("Dark Mode", false)
darkMode.OnClick(func(ctx *application.Context) {
    isChecked := darkMode.Checked()
    // Toggle dark mode
})
```

### AddRadio()

ラジオメニュー項目（相互排他的なグループ）を追加します。

```go
func (m *Menu) AddRadio(label string, checked bool) *MenuItem
```

**パラメーター：**

- `label` - ラジオボタンのラベル
- `checked` - チェック状態の初期値

**例：**

```go
// Create radio group for view modes
viewMenu := menu.AddSubmenu("View")
listView := viewMenu.AddRadio("List View", true)
gridView := viewMenu.AddRadio("Grid View", false)
treeView := viewMenu.AddRadio("Tree View", false)

listView.OnClick(func(ctx *application.Context) {
    setViewMode("list")
})
gridView.OnClick(func(ctx *application.Context) {
    setViewMode("grid")
})
```

### Update()

メニュー項目に加えた変更が反映されるように、メニューを更新します。

```go
func (m *Menu) Update()
```

**例：**

```go
item.SetEnabled(false)
menu.Update()  // Must call to apply changes
```

**重要：** メニュー項目のプロパティを変更した後は、必ず`Update()`を呼び出してください。

## メニュー項目のメソッド

### OnClick()

メニュー項目のクリックハンドラーを登録します。

```go
func (mi *MenuItem) OnClick(callback func(ctx *application.Context)) *MenuItem
```

**パラメーター：**

- `callback` - 項目がクリックされたときに呼び出される関数

**戻り値：** メニュー項目（メソッドチェーン用）

**例：**

```go
item.OnClick(func(ctx *application.Context) {
    fmt.Println("Menu item clicked")
    app.Logger.Info("User clicked menu item")
})
```

### SetLabel()

メニュー項目のラベルを変更します。

```go
func (mi *MenuItem) SetLabel(label string) *MenuItem
```

**例：**

```go
item.SetLabel("Save As...")
menu.Update()
```

### SetEnabled()

メニュー項目を有効または無効にします。

```go
func (mi *MenuItem) SetEnabled(enabled bool) *MenuItem
```

**例：**

```go
// Disable save when no document is open
saveItem.SetEnabled(hasOpenDocument)
menu.Update()
```

**一般的なパターン：**

```go
// Update menu state based on application state
func updateMenuState() {
    saveItem.SetEnabled(hasUnsavedChanges)
    undoItem.SetEnabled(canUndo)
    redoItem.SetEnabled(canRedo)
    menu.Update()
}
```

### SetChecked()

チェックボックス／ラジオメニュー項目のチェック状態を設定します。

```go
func (mi *MenuItem) SetChecked(checked bool) *MenuItem
```

**例：**

```go
darkModeItem.SetChecked(isDarkModeEnabled)
menu.Update()
```

### Checked()

現在のチェック状態を返します。

```go
func (mi *MenuItem) Checked() bool
```

**例：**

```go
if darkModeItem.Checked() {
    // Dark mode is enabled
}
```

### SetAccelerator()

メニュー項目にキーボードショートカットを設定します。

```go
func (mi *MenuItem) SetAccelerator(accelerator string) *MenuItem
```

**パラメーター：**

- `accelerator` - キーボードショートカット（例："Ctrl+S"、"Cmd+Q"）

**アクセラレーターの形式：**

- **修飾キー：** `Ctrl`、`Cmd`、`Alt`、`Shift`
- **キー：** `A-Z`、`0-9`、`F1-F12`、`Enter`、`Backspace`など
- **プラットフォーム：** macOS では`Cmd`、Windows/Linux では`Ctrl`を使用します

**例：**

```go
saveItem.SetAccelerator("Ctrl+S")
quitItem.SetAccelerator("Ctrl+Q")
newItem.SetAccelerator("Ctrl+N")
```

**プラットフォームを考慮した例：**

```go
import "runtime"

var quitShortcut string
if runtime.GOOS == "darwin" {
    quitShortcut = "Cmd+Q"
} else {
    quitShortcut = "Ctrl+Q"
}
quitItem.SetAccelerator(quitShortcut)
```

### SetTooltip()

メニュー項目にカーソルを合わせたときに表示されるツールチップを設定します。

```go
func (mi *MenuItem) SetTooltip(tooltip string) *MenuItem
```

**例：**

```go
item.SetTooltip("Opens a file from disk")
```

### SetHidden()

メニュー項目を表示または非表示にします。

```go
func (mi *MenuItem) SetHidden(hidden bool) *MenuItem
```

**例：**

```go
// Hide debug menu in production
debugItem.SetHidden(!isDevelopment)
menu.Update()
```

## アプリケーションメニュー

### app.Menu.Set()

アプリケーションのメインメニューバーを設定します。

```go
func (mm *MenuManager) Set(menu *Menu)
```

**例：**

```go
menu := app.NewMenu()

// File menu
fileMenu := menu.AddSubmenu("File")
fileMenu.Add("New").SetAccelerator("Ctrl+N").OnClick(newFile)
fileMenu.Add("Open").SetAccelerator("Ctrl+O").OnClick(openFile)
fileMenu.Add("Save").SetAccelerator("Ctrl+S").OnClick(saveFile)
fileMenu.AddSeparator()
fileMenu.Add("Exit").SetAccelerator("Ctrl+Q").OnClick(func(ctx *application.Context) {
    app.Quit()
})

// Edit menu
editMenu := menu.AddSubmenu("Edit")
editMenu.Add("Undo").SetAccelerator("Ctrl+Z").OnClick(undo)
editMenu.Add("Redo").SetAccelerator("Ctrl+Y").OnClick(redo)
editMenu.AddSeparator()
editMenu.Add("Cut").SetAccelerator("Ctrl+X").OnClick(cut)
editMenu.Add("Copy").SetAccelerator("Ctrl+C").OnClick(copy)
editMenu.Add("Paste").SetAccelerator("Ctrl+V").OnClick(paste)

app.Menu.Set(menu)
```

**プラットフォームに関する注記：**

- <strong>macOS：</strong>メニューは画面上部のメニューバーに表示されます
- <strong>Windows/Linux：</strong>メニューはウィンドウのタイトルバーに表示されます
- <strong>macOS：</strong>アプリ名を使用したアプリケーションメニューが自動的に追加されます

## コンテキストメニュー

### app.ContextMenu.New() / app.ContextMenu.Add()

マネージャーを介して`*ContextMenu`を作成し、名前を付けて登録します。`ContextMenuManager.Add`が受け取るのは`*ContextMenu`であり、****`*Menu`ではありません。また、****`app.RegisterContextMenu`メソッドはありません。

```go
func (cm *ContextMenuManager) New() *ContextMenu
func (cm *ContextMenuManager) Add(name string, menu *ContextMenu)
func (cm *ContextMenuManager) Get(name string) (*ContextMenu, bool)
func (cm *ContextMenuManager) Remove(name string)
```

代わりに、パッケージレベルの`application.NewContextMenu(name string) *ContextMenu`を使用すると、コンテキストメニューの作成と登録を一度に行えます。

**Go：**

```go
// Build the context menu via the manager
contextMenu := app.ContextMenu.New()
contextMenu.Add("Cut").OnClick(cut)
contextMenu.Add("Copy").OnClick(copy)
contextMenu.Add("Paste").OnClick(paste)
contextMenu.AddSeparator()
contextMenu.Add("Select All").OnClick(selectAll)

// Register under a name; HTML opts into it via the CSS custom property below.
app.ContextMenu.Add("editor", contextMenu)
```

**HTML / CSS：**

右クリックの対象要素（またはその祖先要素）の`--custom-contextmenu` CSSカスタムプロパティに登録名が設定されていると、ランタイムは登録済みのコンテキストメニューを表示します。省略可能な`--custom-contextmenu-data`プロパティは、`ctx.ContextMenuData()`を介してGoコールバックに渡されます。ブラウザーのデフォルトのコンテキストメニューを抑止するには、`--default-contextmenu: hide`（または`auto`/`show`）を設定します。

```html
<!-- Trigger context menu on right-click -->
<div style="--custom-contextmenu: editor; --default-contextmenu: hide">
    Right-click here for context menu
</div>
```

`data-wails-context-menu="..."`属性はありません。これはランタイムに一度も接続されていません。

**動的コンテキストメニュー：**

```go
// Update context menu based on selection
func updateContextMenu() {
    contextMenu := app.ContextMenu.New()

    if hasSelection {
        contextMenu.Add("Cut").OnClick(cut)
        contextMenu.Add("Copy").OnClick(copy)
    }

    contextMenu.Add("Paste").SetEnabled(hasClipboardContent).OnClick(paste)

    app.ContextMenu.Add("editor", contextMenu)
}
```

## システムトレイメニュー

### app.SystemTray.New()

新しいシステムトレイアイコンを作成します。

```go
func (sm *SystemTrayManager) New() *SystemTray
```

**例：**

```go
tray := app.SystemTray.New()
```

### SetIcon()

システムトレイアイコンを設定します。

```go
func (st *SystemTray) SetIcon(icon []byte) *SystemTray
```

**例：**

```go
iconData, _ := os.ReadFile("icon.png")
tray.SetIcon(iconData)
```

### SetMenu()

システムトレイのメニューを設定します。

```go
func (st *SystemTray) SetMenu(menu *Menu) *SystemTray
```

**例：**

```go
trayMenu := app.NewMenu()
trayMenu.Add("Show Window").OnClick(func(ctx *application.Context) {
    window.Show()
    window.Focus()
})
trayMenu.Add("Settings").OnClick(openSettings)
trayMenu.AddSeparator()
trayMenu.Add("Quit").OnClick(func(ctx *application.Context) {
    app.Quit()
})

tray.SetMenu(trayMenu)
```

### SetTooltip()

トレイアイコンにカーソルを合わせたときに表示されるツールチップを設定します。戻り値はありません。

```go
func (st *SystemTray) SetTooltip(tooltip string)
```

**例：**

```go
tray.SetTooltip("My Application - Running")
```

### OnClick()

トレイアイコンの左クリックを処理します。

```go
func (st *SystemTray) OnClick(callback func()) *SystemTray
```

**例：**

```go
tray.OnClick(func() {
    if window.IsVisible() {
        window.Hide()
    } else {
        window.Show()
        window.Focus()
    }
})
```

## 完全な例

### 標準的なアプリケーションメニュー

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
)

func createMenu(app *application.App) *application.Menu {
    menu := app.NewMenu()

    // File menu
    fileMenu := menu.AddSubmenu("File")
    fileMenu.Add("New").
        SetAccelerator("Ctrl+N").
        OnClick(func(ctx *application.Context) {
            // Create new document
        })
    fileMenu.Add("Open").
        SetAccelerator("Ctrl+O").
        OnClick(func(ctx *application.Context) {
            // Open file dialog
        })
    fileMenu.Add("Save").
        SetAccelerator("Ctrl+S").
        OnClick(func(ctx *application.Context) {
            // Save document
        })
    fileMenu.AddSeparator()
    fileMenu.Add("Exit").OnClick(func(ctx *application.Context) {
        app.Quit()
    })

    // Edit menu
    editMenu := menu.AddSubmenu("Edit")
    editMenu.Add("Undo").SetAccelerator("Ctrl+Z")
    editMenu.Add("Redo").SetAccelerator("Ctrl+Y")
    editMenu.AddSeparator()
    editMenu.Add("Cut").SetAccelerator("Ctrl+X")
    editMenu.Add("Copy").SetAccelerator("Ctrl+C")
    editMenu.Add("Paste").SetAccelerator("Ctrl+V")

    // View menu
    viewMenu := menu.AddSubmenu("View")
    darkMode := viewMenu.AddCheckbox("Dark Mode", false)
    darkMode.OnClick(func(ctx *application.Context) {
        // Toggle dark mode
        isChecked := darkMode.Checked()
        app.Logger.Info("Dark mode", "enabled", isChecked)
    })
    viewMenu.AddSeparator()
    viewMenu.AddRadio("List View", true)
    viewMenu.AddRadio("Grid View", false)
    viewMenu.AddRadio("Detail View", false)

    // Help menu
    helpMenu := menu.AddSubmenu("Help")
    helpMenu.Add("Documentation").OnClick(func(ctx *application.Context) {
        // Open docs
    })
    helpMenu.Add("About").OnClick(func(ctx *application.Context) {
        // Show about dialog
    })

    return menu
}

func main() {
    app := application.New(application.Options{
        Name: "Menu Demo",
    })

    menu := createMenu(app)
    app.Menu.Set(menu)

    window := app.Window.New()
    window.Show()

    app.Run()
}
```

### システムトレイアプリケーション

```go
func setupSystemTray(app *application.App, window application.Window) {
    // Create system tray
    tray := app.SystemTray.New()

    // Set icon
    iconData, _ := os.ReadFile("icon.png")
    tray.SetIcon(iconData)
    tray.SetTooltip("My App - Running")

    // Handle left-click on tray icon
    tray.OnClick(func() {
        if window.IsVisible() {
            window.Hide()
        } else {
            window.Show()
            window.Focus()
        }
    })

    // Create tray menu
    trayMenu := app.NewMenu()

    showItem := trayMenu.Add("Show Window")
    showItem.OnClick(func(ctx *application.Context) {
        window.Show()
        window.Focus()
    })

    trayMenu.AddSeparator()

    trayMenu.Add("Settings").OnClick(func(ctx *application.Context) {
        // Open settings window
    })

    trayMenu.Add("About").OnClick(func(ctx *application.Context) {
        // Show about dialog
    })

    trayMenu.AddSeparator()

    trayMenu.Add("Quit").OnClick(func(ctx *application.Context) {
        app.Quit()
    })

    tray.SetMenu(trayMenu)
}
```

### 動的なメニュー更新

```go
type Editor struct {
    app         *application.App
    menu        *application.Menu
    undoItem    *application.MenuItem
    redoItem    *application.MenuItem
    saveItem    *application.MenuItem
    undoStack   []string
    redoStack   []string
    hasChanges  bool
}

func (e *Editor) createMenu() {
    e.menu = e.app.NewMenu()

    fileMenu := e.menu.AddSubmenu("File")
    e.saveItem = fileMenu.Add("Save").SetAccelerator("Ctrl+S")
    e.saveItem.OnClick(func(ctx *application.Context) {
        e.save()
    })

    editMenu := e.menu.AddSubmenu("Edit")
    e.undoItem = editMenu.Add("Undo").SetAccelerator("Ctrl+Z")
    e.undoItem.OnClick(func(ctx *application.Context) {
        e.undo()
    })

    e.redoItem = editMenu.Add("Redo").SetAccelerator("Ctrl+Y")
    e.redoItem.OnClick(func(ctx *application.Context) {
        e.redo()
    })

    e.updateMenuState()
    e.app.Menu.Set(e.menu)
}

func (e *Editor) updateMenuState() {
    // Update menu items based on current state
    e.saveItem.SetEnabled(e.hasChanges)
    e.undoItem.SetEnabled(len(e.undoStack) > 0)
    e.redoItem.SetEnabled(len(e.redoStack) > 0)
    e.menu.Update()
}

func (e *Editor) onChange() {
    e.hasChanges = true
    e.updateMenuState()
}

func (e *Editor) save() {
    // Save logic
    e.hasChanges = false
    e.updateMenuState()
}
```

## ベストプラクティス

### ✅ 推奨事項

- **標準的なアクセラレーターを使用する** - プラットフォームの慣例に従ってください（コピーにはCtrl+Cなど）
- **変更後にUpdate()を呼び出す** - 呼び出さないと、変更がメニューに反映されません
- **関連する項目をグループ化する** - セパレーターを使用してメニュー項目を整理してください
- **利用できないアクションを無効化する** - 非表示にせず、SetEnabled(false)で無効化してください
- **明確なラベルを使用する** - 簡潔で内容が伝わる表現にしてください
- **プラットフォームの慣例に従う** - macOSとWindows/Linuxではメニューのパターンが異なります

### ❌ 禁止事項

- **Update()の呼び出しを忘れない** - 最もよくある間違いです
- **深くネストしすぎない** - メニューは最大2-3階層にしてください
- **曖昧なラベルを使用しない** - 「処理」ではなく「ドキュメントを処理」のように明確にしてください
- **複雑にしすぎない** - メニューはシンプルで目的を絞った構成にしてください
- **異なる概念を混在させない** - 命名と構成に一貫性を持たせてください

## プラットフォーム固有の注意事項

### macOS

- アプリ名を使用したアプリケーションメニューが自動的に追加されます
- アクセラレーターには`Ctrl`ではなく`Cmd`を使用します
- デフォルトでは、アプリケーションメニューに「このアプリケーションについて」、「環境設定」、「終了」が含まれます

### Windows/Linux

- アプリケーションメニューは自動的に追加されません
- アクセラレーターには`Ctrl`を使用します
- 通常、「終了」は「ファイル」メニューに配置します
