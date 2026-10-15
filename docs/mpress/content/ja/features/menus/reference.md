---
title: "メニューリファレンス"
description: "メニュー項目の種類、プロパティ、メソッドに関する完全なリファレンス"
slug: "features/menus/reference"
sourcePath: "features/menus/reference.md"
---

## メニューリファレンス

メニュー項目の種類、プロパティ、動的な動作に関する完全なリファレンスです。チェックボックス、ラジオグループ、セパレーター、動的更新を使用して、プロフェッショナルで応答性の高いメニューを構築します。

## メニュー項目の種類

### 通常のメニュー項目

最も一般的な種類です。テキストを表示し、アクションを実行します：

```go
menuItem := menu.Add("Click Me")
menuItem.OnClick(func(ctx *application.Context) {
    fmt.Println("Menu item clicked!")
})
```

**用途：** コマンド、アクション、ウィンドウを開く操作

### チェックボックス

オンとオフを切り替えられる、チェック済み／未チェック状態を持つメニュー項目です：

```go
checkbox := menu.AddCheckbox("Enable Feature", true)  // true = initially checked
checkbox.OnClick(func(ctx *application.Context) {
    isChecked := ctx.ClickedMenuItem().Checked()
    fmt.Printf("Feature is now: %v\n", isChecked)
})
```

**用途：** ブール値の設定、機能の切り替え、表示オプション

**重要：** クリックすると、チェック状態が自動的に切り替わります。

### ラジオグループ

相互排他的なオプションで、一度に1つだけ選択できます：

```go
menu.AddRadio("Small", true)   // true = initially selected
menu.AddRadio("Medium", false)
menu.AddRadio("Large", false)
```

**用途：** 相互排他的な選択肢（サイズ、テーマ、モード）

**グループ化の仕組み：**

- 隣接するラジオ項目は自動的にグループを形成します
- 1つを選択すると、グループ内のほかの項目は選択解除されます
- セパレーターまたは通常の項目でグループを分けます

**複数のグループを使用する例：**

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

項目を整理するための、入れ子になったメニュー構造です：

```go
submenu := menu.AddSubmenu("More Options")
submenu.Add("Submenu Item 1").OnClick(func(ctx *application.Context) {
    // Handle click
})
submenu.Add("Submenu Item 2")
```

**用途：** 関連項目のグループ化、煩雑さの軽減

**入れ子の上限：** ほとんどのプラットフォームは2-3階層までサポートします。それより深い入れ子は避けてください。

### セパレーター

メニュー項目間に表示する視覚的な区切りです：

```go
menu.Add("Item 1")
menu.AddSeparator()
menu.Add("Item 2")
```

**用途：** 関連項目を視覚的にグループ化する

**ベストプラクティス：** メニューの先頭や末尾にセパレーターを配置しないでください。

## メニュー項目のプロパティ

### ラベル

メニュー項目に表示されるテキストです：

```go
menuItem := menu.Add("Initial Label")
menuItem.SetLabel("New Label")

// Get current label
label := menuItem.Label()
```

**動的ラベル：**

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

メニュー項目を操作できるかどうかを制御します：

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
Windowsでは、メニューの状態が変わったときにメニューを再構築する必要があります。メニュー項目を有効化または無効化した後は<strong>必ず`menu.Update()`を呼び出してください</strong>。特に、その項目が無効な状態で作成された場合は必須です。

**理由：** Windowsのメニューは、更新時に一から再構築されます。`Update()`を呼び出さないと、クリックハンドラーが正しく実行されません。

@end

**例：動的な有効化／無効化**

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

**一般的なパターン：条件に応じて有効化**

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

チェックボックス項目とラジオ項目では、チェック状態を制御または照会できます：

```go
checkbox := menu.AddCheckbox("Feature", false)
checkbox.SetChecked(true)
menu.Update()

// Query state
isChecked := checkbox.Checked()
```

**自動切り替え：** チェックボックスはクリックすると自動的に切り替わります。クリックハンドラー内で`SetChecked()`を呼び出す必要はありません。

**手動制御：**

```go
checkbox := menu.AddCheckbox("Auto-save", false)

// Sync with external state
func syncAutoSave(enabled bool) {
    checkbox.SetChecked(enabled)
    menu.Update()
}
```

### アクセラレーター（キーボードショートカット）

メニュー項目にキーボードショートカットを追加します：

```go
saveMenuItem := menu.Add("Save")
saveMenuItem.SetAccelerator("CmdOrCtrl+S")

quitMenuItem := menu.Add("Quit")
quitMenuItem.SetAccelerator("CmdOrCtrl+Q")
```

**アクセラレーターの形式：**

- `CmdOrCtrl` - macOSではCmd、Windows/LinuxではCtrl
- `Shift`、`Alt`、`Option` - 修飾キー
- `A-Z`、`0-9` - 英字／数字キー
- `F1-F12` - ファンクションキー
- `Enter`、`Space`、`Backspace`など - 特殊キー

**例：**

```go
"CmdOrCtrl+S"           // Save
"CmdOrCtrl+Shift+S"     // Save As
"CmdOrCtrl+W"           // Close Window
"CmdOrCtrl+Q"           // Quit
"F5"                    // Refresh
"CmdOrCtrl+,"           // Preferences (macOS convention)
"Alt+F4"                // Close (Windows convention)
```

**プラットフォーム固有のアクセラレーター：**

```go
if runtime.GOOS == "darwin" {
    prefsMenuItem.SetAccelerator("Cmd+,")
} else {
    prefsMenuItem.SetAccelerator("Ctrl+P")
}
```

### ツールチップ

メニュー項目に、ポインターを合わせたときに表示されるテキストを追加します（対応状況はプラットフォームによって異なります）：

```go
menuItem := menu.Add("Advanced Options")
menuItem.SetTooltip("Configure advanced settings")
```

**プラットフォームの対応状況：**

- **Windows：** ✅ 対応
- **macOS：** ❌ 非対応（メニューではツールチップが標準ではないため）
- **Linux：** ⚠️ デスクトップ環境によって異なります

### 非表示状態

メニュー項目を削除せずに非表示にします：

```go
debugMenuItem := menu.Add("Debug Mode")
debugMenuItem.SetHidden(true)  // Hidden

// Show in debug builds
if isDebugBuild {
    debugMenuItem.SetHidden(false)
    menu.Update()
}
```

**用途：** デバッグオプション、機能フラグ、条件付き機能

## イベント処理

### OnClickハンドラー

メニュー項目がクリックされたときにコードを実行します：

```go
menuItem := menu.Add("Click Me")
menuItem.OnClick(func(ctx *application.Context) {
    // Handle click
    fmt.Println("Clicked!")
})
```

**コンテキストから取得できる情報：**

- `ctx.ClickedMenuItem()` - クリックされたメニュー項目
- ウィンドウのコンテキスト（ウィンドウメニューからの場合）
- アプリケーションのコンテキスト

**例：ハンドラー内でメニュー項目にアクセスする**

```go
checkbox := menu.AddCheckbox("Feature", false)
checkbox.OnClick(func(ctx *application.Context) {
    item := ctx.ClickedMenuItem()
    isChecked := item.Checked()
    fmt.Printf("Feature is now: %v\n", isChecked)
})
```

### 複数のハンドラー

複数のハンドラーを設定できます（最後に設定したものが優先されます）：

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

**ベストプラクティス：** ハンドラーは一度だけ設定し、必要に応じてその内部で条件分岐を使用します。

## 動的メニュー

### メニュー項目の更新

**鉄則：** メニューの状態を変更した後は、必ず `menu.Update()` を呼び出します。

```go
// ✅ Correct
menuItem.SetEnabled(true)
menu.Update()

// ❌ Wrong (especially on Windows)
menuItem.SetEnabled(true)
// Forgot to call Update() - click handlers may not work!
```

**これが重要な理由：**

- **Windows：** 更新時にメニューが再構築されます
- **macOS/Linux：** 重要度は低いものの、引き続き推奨されます
- **クリックハンドラー：** Update() を呼び出さないと正しく実行されません

### メニューの再構築

大幅な変更を行う場合は、メニュー全体を再構築します：

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

- 最近使用したファイルの一覧が変わる場合
- プラグインのメニューが変わる場合
- 大幅な状態遷移が発生する場合

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
| **Linux** | ウィンドウ上部 | ウィンドウごと（通常） |

### 標準メニュー

**macOS：**

- アプリ名を表示した「アプリケーション」メニューがあります
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

- ほとんどのショートカットでは `Cmd+`
- 「環境設定」では `Cmd+,`
- 「終了」では `Cmd+Q`

**Windows：**

- ほとんどのショートカットでは `Ctrl+`
- 「環境設定」では `Ctrl+P` または `Ctrl+,`
- 「終了」では `Alt+F4`（または `Ctrl+Q`）

**Linux：**

- 通常は Windows の慣例に従います
- デスクトップ環境によって上書きされる場合があります

## ベストプラクティス

### ✅ 推奨事項

- メニューの状態を変更した後は、**menu.Update() を呼び出します**（特に Windows）
- 相互排他的なオプションには<strong>ラジオグループを使用します</strong>
- オンとオフを切り替えられる機能には<strong>チェックボックスを使用します</strong>
- 一般的な操作には<strong>アクセラレータを追加する</strong>
- <strong>関連する項目を</strong>セパレーターでグループ化する
- 動作が異なるため、**すべてのプラットフォームでテストする**

### ❌ 避けるべきこと

- **menu.Update() の呼び出しを忘れない** — クリックハンドラーが正しく動作しなくなります
- **深くネストしすぎない** — 最大2-3階層までにします
- **先頭や末尾にセパレーターを配置しない** — 見栄えが悪くなります
- **macOS ではツールチップを使用しない** — サポートされていません
- **プラットフォーム固有のショートカットをハードコードしない** — `CmdOrCtrl`を使用します

## トラブルシューティング

### メニュー項目が反応しない

<strong>症状：</strong>クリックハンドラーが呼び出されない

<strong>原因：</strong>項目を有効にした後、`menu.Update()`の呼び出しを忘れている

**解決方法：**

```go
menuItem.SetEnabled(true)
menu.Update()  // Add this!
```

### メニュー項目がグレー表示される

<strong>症状：</strong>メニュー項目をクリックできない

<strong>原因：</strong>項目が無効になっている

**解決方法：**

```go
menuItem.SetEnabled(true)
menu.Update()
```

### アクセラレータが機能しない

<strong>症状：</strong>キーボードショートカットでメニュー項目が実行されない

**原因：**

1. アクセラレータの形式が正しくない
2. システムショートカットと競合している
3. ウィンドウにフォーカスがない

**解決方法：**

```go
// Check format
menuItem.SetAccelerator("CmdOrCtrl+S")  // ✅ Correct
menuItem.SetAccelerator("Ctrl+S")       // ❌ Wrong (macOS uses Cmd)

// Avoid conflicts
// ❌ Cmd+H (Hide Window on macOS - system shortcut)
// ✅ Cmd+Shift+H (Custom shortcut)
```

## 次のステップ

- [アプリケーションメニュー](/features/menus/application/) — アプリケーションのメニューバーを作成する
- [コンテキストメニュー](/features/menus/context/) — 右クリックで表示するコンテキストメニューを作成する
- [システムトレイメニュー](/features/menus/systray/) — システムトレイ／メニューバーのメニューを作成する
- [メニューパターン](/guides/menus/) — 一般的なメニューパターンとベストプラクティス

---

**質問がありますか？** [Discord](https://discord.gg/JDdSxwjhGf) で質問するか、[メニューのサンプル](https://github.com/wailsapp/wails/tree/master/v3/examples/menu)を確認してください。
