---
title: "メニューリファレンス"
description: "メニュー項目の種類、プロパティ、メソッドの完全なリファレンス"
slug: "features/menus/reference"
sourcePath: "features/menus/reference.md"
---

## メニューリファレンス

メニュー項目の種類、プロパティ、動的な動作についての完全なリファレンスです。チェックボックス、ラジオグループ、セパレーター、動的更新を使用して、プロフェッショナルで応答性の高いメニューを構築します。

## メニュー項目の種類

### 通常のメニュー項目

最も一般的な種類です。テキストを表示し、アクションを実行します。

```go
menuItem := menu.Add("Click Me")
menuItem.OnClick(func(ctx *application.Context) {
    fmt.Println("Menu item clicked!")
})
```

**用途:** コマンド、アクション、ウィンドウを開く操作

### チェックボックス

オンとオフを切り替えられる、チェック済み／未チェック状態を持つメニュー項目です。

```go
checkbox := menu.AddCheckbox("Enable Feature", true)  // true = initially checked
checkbox.OnClick(func(ctx *application.Context) {
    isChecked := ctx.ClickedMenuItem().Checked()
    fmt.Printf("Feature is now: %v\n", isChecked)
})
```

**用途:** ブール値の設定、機能の切り替え、表示オプション

**重要:** クリックすると、チェック状態が自動的に切り替わります。

### ラジオグループ

相互排他的な選択肢です。選択できるのは1つだけです。

```go
menu.AddRadio("Small", true)   // true = initially selected
menu.AddRadio("Medium", false)
menu.AddRadio("Large", false)
```

**用途:** 相互排他的な選択肢（サイズ、テーマ、モード）

**グループ化の仕組み:**

- 隣接するラジオ項目は自動的にグループを形成します
- 1つを選択すると、グループ内のほかの項目の選択が解除されます
- セパレーターまたは通常の項目でグループを分けます

**複数のグループを使用する例:**

```go
// Group 1: Size
menu.AddRadio("Small", true)
menu.AddRadio("Medium", false)
menu.AddRadio("Large", false)

menu.AddSeparator()

// Group 2: Theme
menu.AddRadio("Light", true)
menu.AddRadio("Dark", false)
```

### サブメニュー

整理するための入れ子構造のメニューです。

```go
submenu := menu.AddSubmenu("More Options")
submenu.Add("Submenu Item 1").OnClick(func(ctx *application.Context) {
    // Handle click
})
submenu.Add("Submenu Item 2")
```

**用途:** 関連項目のグループ化、煩雑さの軽減

**入れ子の上限:** ほとんどのプラットフォームは2-3階層をサポートします。それより深い入れ子は避けてください。

### セパレーター

メニュー項目間の視覚的な区切りです。

```go
menu.Add("Item 1")
menu.AddSeparator()
menu.Add("Item 2")
```

**用途:** 関連項目を視覚的にグループ化する

**ベストプラクティス:** メニューの先頭または末尾にセパレーターを配置しないでください。

### セクション見出し

**macOS 14以降**で、メニュー項目のグループに使用できるネイティブの非対話型ラベルです。

```go
menu.AddSectionHeader("Build server")
menu.Add("Start build").OnClick(startBuild)
menu.Add("View logs").OnClick(viewLogs)
```

`AddSectionHeader`は、AppKit標準のセクション見出しの外観を使用します。メニュー項目は返さず、アクションも関連付けません。古いバージョンのmacOS、Windows、Linux、モバイルおよびサーバービルドでは、この呼び出しは何も行わず、行も挿入しません。

トレイメニュー、コンテキストメニュー、またはアプリケーションメニューで、関連するコマンドに見出しを付けるために使用します。省略された見出しではラジオグループが分断されないため、すべてのプラットフォームでラジオグループを分けるには`AddSeparator()`を使用してください。メソッドシグネチャについては[メニューAPI](/reference/menu#addsectionheader)を、ホスト名を見出しにする例については`systray-menu`の例を参照してください。

## メニュー項目のプロパティ

### ラベル

メニュー項目に表示されるテキストです。

```go
menuItem := menu.Add("Initial Label")
menuItem.SetLabel("New Label")

// Get current label
label := menuItem.Label()
```

**動的ラベル:**

```go
updateMenuItem := menu.Add("Check for Updates")
updateMenuItem.OnClick(func(ctx *application.Context) {
    updateMenuItem.SetLabel("Checking...")
    menu.Update()  // Important on Windows!
    
    // Perform update check
    checkForUpdates()
    
    updateMenuItem.SetLabel("Check for Updates")
    menu.Update()
})
```

### 有効状態

メニュー項目を操作できるかどうかを制御します。

```go
menuItem := menu.Add("Save")
menuItem.SetEnabled(false)  // Greyed out, can't click

// Enable it later
menuItem.SetEnabled(true)
menu.Update()  // Important: Call this after changing enabled state!

// Check current state
isEnabled := menuItem.Enabled()
```

@note{type="caution" title="Windowsでのメニューの動作"}
Windowsでは、メニューの状態が変化したときにメニューを再構築する必要があります。メニュー項目を有効化または無効化した後は**必ず`menu.Update()`を呼び出してください**。特に、その項目を無効な状態で作成した場合は必須です。

**理由:** Windowsのメニューは、更新時に一から再構築されます。`Update()`を呼び出さないと、クリックハンドラーが正しく実行されません。

@end

**例: 動的な有効化／無効化**

```go
var hasSelection bool

cutMenuItem := menu.Add("Cut")
cutMenuItem.SetEnabled(false)  // Initially disabled

copyMenuItem := menu.Add("Copy")
copyMenuItem.SetEnabled(false)

// When selection changes
func onSelectionChanged(selected bool) {
    hasSelection = selected
    cutMenuItem.SetEnabled(hasSelection)
    copyMenuItem.SetEnabled(hasSelection)
    menu.Update()  // Critical on Windows!
}
```

**一般的なパターン: 条件に応じて有効化**

```go
saveMenuItem := menu.Add("Save")

func updateSaveMenuItem() {
    canSave := hasUnsavedChanges() && !isSaving()
    saveMenuItem.SetEnabled(canSave)
    menu.Update()
}

// Call whenever state changes
onDocumentChanged(func() {
    updateSaveMenuItem()
})
```

### チェック状態

チェックボックス項目とラジオ項目では、チェック状態を制御または取得できます。

```go
checkbox := menu.AddCheckbox("Feature", false)
checkbox.SetChecked(true)
menu.Update()

// Query state
isChecked := checkbox.Checked()
```

**自動切り替え:** チェックボックスは、クリックすると自動的に切り替わります。クリックハンドラー内で`SetChecked()`を呼び出す必要はありません。

**手動制御:**

```go
checkbox := menu.AddCheckbox("Auto-save", false)

// Sync with external state
func syncAutoSave(enabled bool) {
    checkbox.SetChecked(enabled)
    menu.Update()
}
```

### アクセラレーター（キーボードショートカット）

メニュー項目にキーボードショートカットを追加します。

```go
saveMenuItem := menu.Add("Save")
saveMenuItem.SetAccelerator("CmdOrCtrl+S")

quitMenuItem := menu.Add("Quit")
quitMenuItem.SetAccelerator("CmdOrCtrl+Q")
```

**アクセラレーターの形式:**

- `CmdOrCtrl` - macOSではCmd、Windows/LinuxではCtrl
- `Shift`、`Alt`、`Option` - 修飾キー
- `A-Z`、`0-9` - 英字／数字キー
- `F1-F12` - ファンクションキー
- `Enter`、`Space`、`Backspace`など - 特殊キー

**例:**

```go
"CmdOrCtrl+S"           // Save
"CmdOrCtrl+Shift+S"     // Save As
"CmdOrCtrl+W"           // Close Window
"CmdOrCtrl+Q"           // Quit
"F5"                    // Refresh
"CmdOrCtrl+,"           // Preferences (macOS convention)
"Alt+F4"                // Close (Windows convention)
```

**プラットフォーム固有のアクセラレーター:**

```go
if runtime.GOOS == "darwin" {
    prefsMenuItem.SetAccelerator("Cmd+,")
} else {
    prefsMenuItem.SetAccelerator("Ctrl+P")
}
```

### ツールチップ

メニュー項目に、マウスポインターを重ねたときに表示するテキストを追加します（対応状況はプラットフォームによって異なります）。

```go
menuItem := menu.Add("Advanced Options")
menuItem.SetTooltip("Configure advanced settings")
```

**プラットフォームの対応状況:**

- **Windows:** ✅ 対応
- **macOS:** ❌ 非対応（メニューでのツールチップは標準ではありません）
- **Linux:** ⚠️ デスクトップ環境によって異なります

### 非表示状態

メニュー項目を削除せずに非表示にします。

```go
debugMenuItem := menu.Add("Debug Mode")
debugMenuItem.SetHidden(true)  // Hidden

// Show in debug builds
if isDebugBuild {
    debugMenuItem.SetHidden(false)
    menu.Update()
}
```

**用途:** デバッグオプション、機能フラグ、条件付き機能

## イベント処理

### OnClickハンドラー

メニュー項目がクリックされたときにコードを実行します。

```go
menuItem := menu.Add("Click Me")
menuItem.OnClick(func(ctx *application.Context) {
    // Handle click
    fmt.Println("Clicked!")
})
```

**コンテキストから取得できるもの:**

- `ctx.ClickedMenuItem()` - クリックされたメニュー項目
- ウィンドウのコンテキスト（ウィンドウメニューからの場合）
- アプリケーションのコンテキスト

**例: ハンドラー内でメニュー項目にアクセスする**

```go
checkbox := menu.AddCheckbox("Feature", false)
checkbox.OnClick(func(ctx *application.Context) {
    item := ctx.ClickedMenuItem()
    isChecked := item.Checked()
    fmt.Printf("Feature is now: %v\n", isChecked)
})
```

### 複数のハンドラー

複数のハンドラーを設定できます（最後に設定したものが優先されます）。

```go
menuItem := menu.Add("Action")
menuItem.OnClick(func(ctx *application.Context) {
    fmt.Println("First handler")
})

// This replaces the first handler
menuItem.OnClick(func(ctx *application.Context) {
    fmt.Println("Second handler - this one runs")
})
```

**ベストプラクティス:** ハンドラーは一度だけ設定し、必要に応じてその内部で条件分岐を使用してください。

## 動的メニュー

### メニュー項目の更新

**鉄則:** メニューの状態を変更した後は、必ず`menu.Update()`を呼び出してください。

```go
// ✅ Correct
menuItem.SetEnabled(true)
menu.Update()

// ❌ Wrong (especially on Windows)
menuItem.SetEnabled(true)
// Forgot to call Update() - click handlers may not work!
```

**これが重要な理由:**

- **Windows:** メニューは更新時に再構築されます
- <strong>macOS/Linux：</strong>重要度は低いものの、引き続き推奨されます
- <strong>クリックハンドラー：</strong>Update() を呼び出さないと正常に発火しません

### メニューの再構築

大幅に変更する場合は、メニュー全体を再構築します：

```go
func rebuildFileMenu() {
    menu := app.Menu.New()
    
    menu.Add("New").OnClick(handleNew)
    menu.Add("Open").OnClick(handleOpen)
    
    if hasRecentFiles() {
        recentMenu := menu.AddSubmenu("Open Recent")
        for _, file := range getRecentFiles() {
            recentMenu.Add(file).OnClick(func(ctx *application.Context) {
                openFile(file)
            })
        }
    }
    
    menu.AddSeparator()
    menu.Add("Quit").OnClick(handleQuit)
    
    // Set the new menu
    window.SetMenu(menu)
}
```

**再構築する場合：**

- 最近使ったファイルの一覧を変更する場合
- プラグインメニューを変更する場合
- 大きな状態遷移がある場合

**更新する場合：**

- 項目を有効化または無効化する場合
- ラベルを変更する場合
- チェックボックスの状態を切り替える場合

### コンテキスト依存メニュー

アプリケーションの状態に応じてメニューを調整します：

```go
func updateEditMenu() {
    cutMenuItem.SetEnabled(hasSelection())
    copyMenuItem.SetEnabled(hasSelection())
    pasteMenuItem.SetEnabled(hasClipboardContent())
    undoMenuItem.SetEnabled(canUndo())
    redoMenuItem.SetEnabled(canRedo())
    menu.Update()
}

// Call whenever state changes
onSelectionChanged(updateEditMenu)
onClipboardChanged(updateEditMenu)
onUndoStackChanged(updateEditMenu)
```

## プラットフォームによる違い

### メニューバーの位置

| プラットフォーム | 位置 | 備考 |
| --- | --- | --- |
| **macOS** | 画面上部 | グローバルメニューバー |
| **Windows** | ウィンドウ上部 | ウィンドウごとのメニュー |
| **Linux** | ウィンドウ上部 | 通常はウィンドウごと |

### 標準メニュー

**macOS：**

- アプリ名を表示する「アプリケーション」メニューがあります
- 「環境設定」は「アプリケーション」メニューにあります
- 「終了」は「アプリケーション」メニューにあります

**Windows/Linux：**

- 「アプリケーション」メニューはありません
- 「環境設定」は「編集」または「ツール」メニューにあります
- 「終了」は「ファイル」メニューにあります

**例：プラットフォームに適した構成**

```go
menu := app.Menu.New()

// macOS gets Application menu
if runtime.GOOS == "darwin" {
    menu.AddRole(application.AppMenu)
}

// File menu
fileMenu := menu.AddSubmenu("File")
fileMenu.Add("New")
fileMenu.Add("Open")

// Preferences location varies
if runtime.GOOS == "darwin" {
    // On macOS, preferences are in Application menu (added by AppMenu role)
} else {
    // On Windows/Linux, add to Edit or Tools menu
    editMenu := menu.AddSubmenu("Edit")
    editMenu.Add("Preferences")
}
```

### アクセラレーターの慣例

**macOS：**

- ほとんどのショートカットには `Cmd+` を使用します
- 「環境設定」には `Cmd+,` を使用します
- 「終了」には `Cmd+Q` を使用します

**Windows：**

- ほとんどのショートカットには `Ctrl+` を使用します
- 「環境設定」には `Ctrl+P` または `Ctrl+,` を使用します
- 「終了」には `Alt+F4`（または `Ctrl+Q`）を使用します

**Linux：**

- 通常は Windows の慣例に従います
- デスクトップ環境によって上書きされる場合があります

## ベストプラクティス

### ✅ 推奨事項

- メニューの状態を変更した後は、特に Windows では**menu.Update() を呼び出してください**
- 排他的な選択肢には**ラジオグループを使用してください**
- オンとオフを切り替えられる機能には**チェックボックスを使用してください**
- よく使う操作には**アクセラレーターを追加してください**
- 関連する項目は区切り線で**グループ化してください**
- 動作が異なるため、**すべてのプラットフォームでテストしてください**

### ❌ 禁止事項

- **menu.Update() の呼び出しを忘れないでください**。クリックハンドラーが正常に動作しなくなります
- **階層を深くしすぎないでください**。最大 2-3 階層にします
- **先頭または末尾に区切り線を置かないでください**。洗練されていない印象になります
- **macOS ではツールチップを使用しないでください**。サポートされていません
- **プラットフォーム固有のショートカットをハードコードしないでください**。`CmdOrCtrl` を使用します

## トラブルシューティング

### メニュー項目が反応しない

<strong>症状：</strong>クリックハンドラーが発火しない

<strong>原因：</strong>項目を有効化した後に `menu.Update()` を呼び出していない

**解決策：**

```go
menuItem.SetEnabled(true)
menu.Update()  // Add this!
```

### メニュー項目がグレーアウトしている

<strong>症状：</strong>メニュー項目をクリックできない

<strong>原因：</strong>項目が無効になっている

**解決策：**

```go
menuItem.SetEnabled(true)
menu.Update()
```

### アクセラレーターが動作しない

<strong>症状：</strong>キーボードショートカットを押してもメニュー項目が実行されない

**原因：**

1. アクセラレーターの形式が正しくない
2. システムのショートカットと競合している
3. ウィンドウにフォーカスがない

**解決策：**

```go
// Check format
menuItem.SetAccelerator("CmdOrCtrl+S")  // ✅ Correct
menuItem.SetAccelerator("Ctrl+S")       // ❌ Wrong (macOS uses Cmd)

// Avoid conflicts
// ❌ Cmd+H (Hide Window on macOS - system shortcut)
// ✅ Cmd+Shift+H (Custom shortcut)
```

## 次のステップ

- [アプリケーションメニュー](/features/menus/application/)：アプリケーションのメニューバーを作成する
- [コンテキストメニュー](/features/menus/context/)：右クリックのコンテキストメニューを作成する
- [システムトレイメニュー](/features/menus/systray/)：システムトレイ／メニューバーのメニューを作成する
- [メニューパターン](/guides/menus/)：一般的なメニューパターンとベストプラクティス

---

**ご質問がありますか？** [Discord](https://discord.gg/JDdSxwjhGf) で質問するか、[メニューの例](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/menu)を確認してください。
