---
title: "コンテキストメニュー"
description: "アプリケーション用の右クリックコンテキストメニューを作成する"
slug: "features/menus/context"
sourcePath: "features/menus/context.md"
---

## 課題

ユーザーは、状況に応じた操作を実行できる右クリックメニューを期待します。要素ごとに異なるメニューが必要です。

- **テキスト**：切り取り、コピー、貼り付け
- **画像**：保存、コピー、開く
- **カスタム要素**：アプリケーション固有の操作

コンテキストメニューを手動で構築する場合は、マウスイベント、表示位置、プラットフォーム間の違いを処理する必要があります。

## Wailsによる解決策

Wailsでは、CSSプロパティを使用する<strong>宣言的なコンテキストメニュー</strong>を利用できます。メニューをHTML要素に関連付け、データを渡し、クリックを処理できます。これらはすべて、各プラットフォームにネイティブな動作で実行されます。

![macOS上のWebViewに重ねて表示されたWailsのカスタムコンテキストメニュー](/assets/screenshots/context-menu-macos.png)

メニューはプラットフォームにネイティブですが、メニューを開いた要素はWebViewの一部のままです。このmacOSのキャプチャでは、以下の例にあるカスタムコンテキストメニューの登録を使用しています。

## クイックスタート

**Goコード：**

```go
// Create context menu
contextMenu := app.ContextMenu.New()
contextMenu.Add("Cut").OnClick(handleCut)
contextMenu.Add("Copy").OnClick(handleCopy)
contextMenu.Add("Paste").OnClick(handlePaste)

// Register with ID
app.ContextMenu.Add("editor-menu", contextMenu)
```

**HTML：**

```html
<textarea style="--custom-contextmenu: editor-menu">
    Right-click me!
</textarea>
```

<strong>これだけです。</strong>テキストエリアを右クリックすると、カスタムメニューが表示されます。

## コンテキストメニューの作成

### 基本的なコンテキストメニュー

```go
// Create menu
contextMenu := app.ContextMenu.New()

// Add items
contextMenu.Add("Cut").SetAccelerator("CmdOrCtrl+X").OnClick(func(ctx *application.Context) {
    // Handle cut
})

contextMenu.Add("Copy").SetAccelerator("CmdOrCtrl+C").OnClick(func(ctx *application.Context) {
    // Handle copy
})

contextMenu.Add("Paste").SetAccelerator("CmdOrCtrl+V").OnClick(func(ctx *application.Context) {
    // Handle paste
})

// Register with unique ID
app.ContextMenu.Add("text-menu", contextMenu)
```

<strong>メニューID：</strong>一意でなければなりません。メニューをHTML要素に関連付けるために使用します。

### サブメニューの使用

```go
contextMenu := app.ContextMenu.New()

// Add regular items
contextMenu.Add("Open").OnClick(handleOpen)
contextMenu.Add("Delete").OnClick(handleDelete)

contextMenu.AddSeparator()

// Add submenu
exportMenu := contextMenu.AddSubmenu("Export As")
exportMenu.Add("PNG").OnClick(exportPNG)
exportMenu.Add("JPEG").OnClick(exportJPEG)
exportMenu.Add("SVG").OnClick(exportSVG)

app.ContextMenu.Add("image-menu", contextMenu)
```

### チェックボックスとラジオグループの使用

```go
contextMenu := app.ContextMenu.New()

// Checkbox
contextMenu.AddCheckbox("Show Grid", true).OnClick(func(ctx *application.Context) {
    showGrid := ctx.ClickedMenuItem().Checked()
    // Toggle grid
})

contextMenu.AddSeparator()

// Radio group
contextMenu.AddRadio("Small", false).OnClick(handleSize)
contextMenu.AddRadio("Medium", true).OnClick(handleSize)
contextMenu.AddRadio("Large", false).OnClick(handleSize)

app.ContextMenu.Add("view-menu", contextMenu)
```

<strong>すべてのメニュー項目タイプ</strong>については、[メニューリファレンス](/features/menus/reference/)を参照してください。

## HTML要素との関連付け

CSSカスタムプロパティを使用して、コンテキストメニューを関連付けます。

### 基本的な関連付け

```html
<div style="--custom-contextmenu: menu-id">
    Right-click me!
</div>
```

**CSSプロパティ：** `--custom-contextmenu: <menu-id>`

### コンテキストデータの使用

HTMLからGoにデータを渡します。

```html
<div style="--custom-contextmenu: file-menu; --custom-contextmenu-data: file-123">
    Right-click this file
</div>
```

**Goハンドラー：**

```go
contextMenu := app.ContextMenu.New()
contextMenu.Add("Open").OnClick(func(ctx *application.Context) {
    fileID := ctx.ContextMenuData()  // "file-123"
    openFile(fileID)
})

app.ContextMenu.Add("file-menu", contextMenu)
```

**CSSプロパティ：**

- `--custom-contextmenu: <menu-id>` - 表示するメニュー
- `--custom-contextmenu-data: <data>` - ハンドラーに渡すデータ

### 動的データ

JavaScriptでデータを動的に生成します。

```html
<div id="file-item" style="--custom-contextmenu: file-menu">
    File.txt
</div>

<script>
// Set data dynamically
const fileItem = document.getElementById('file-item')
fileItem.style.setProperty('--custom-contextmenu-data', 'file-' + fileId)
</script>
```

### 複数の要素で同じメニューを使用する

```html
<div class="file-item" style="--custom-contextmenu: file-menu; --custom-contextmenu-data: file-1">
    Document.pdf
</div>

<div class="file-item" style="--custom-contextmenu: file-menu; --custom-contextmenu-data: file-2">
    Image.png
</div>

<div class="file-item" style="--custom-contextmenu: file-menu; --custom-contextmenu-data: file-3">
    Video.mp4
</div>
```

**1つのメニューを使用し、要素ごとに異なるデータを渡します。**

## コンテキストデータ

### コンテキストデータへのアクセス

```go
contextMenu.Add("Process").OnClick(func(ctx *application.Context) {
    data := ctx.ContextMenuData()  // Get data from HTML
    
    // Use the data
    processItem(data)
})
```

<strong>データ型：</strong>常に`string`です。必要に応じて解析してください。

### 複雑なデータの受け渡し

複雑なデータにはJSONを使用します。

```html
<div style="--custom-contextmenu: item-menu; --custom-contextmenu-data: {&quot;id&quot;:123,&quot;type&quot;:&quot;image&quot;}">
    Image.png
</div>
```

**Goハンドラー：**

```go
import "encoding/json"

type ItemData struct {
    ID   int    `json:"id"`
    Type string `json:"type"`
}

contextMenu.Add("Process").OnClick(func(ctx *application.Context) {
    dataStr := ctx.ContextMenuData()
    
    var data ItemData
    if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
        log.Printf("Invalid data: %v", err)
        return
    }
    
    processItem(data.ID, data.Type)
})
```

@note{type="caution" title="セキュリティ"}
フロントエンドから受け取ったコンテキストデータは、**必ず検証してください**。ユーザーはCSSプロパティを操作できるため、このデータを信頼できない入力として扱ってください。

@end

### 検証例

```go
contextMenu.Add("Delete").OnClick(func(ctx *application.Context) {
    fileID := ctx.ContextMenuData()
    
    // Validate
    if !isValidFileID(fileID) {
        log.Printf("Invalid file ID: %s", fileID)
        return
    }
    
    // Check permissions
    if !canDeleteFile(fileID) {
        showError("Permission denied")
        return
    }
    
    // Safe to proceed
    deleteFile(fileID)
})
```

## デフォルトのコンテキストメニュー

WebViewには、標準操作（コピー、貼り付け、検査）用のコンテキストメニューが組み込まれています。`--default-contextmenu`で制御します。

### デフォルトメニューを非表示にする

```html
<div style="--default-contextmenu: hide">
    No default menu here
</div>
```

<strong>使用例：</strong>デフォルトメニューが適さないカスタムUI要素。

### デフォルトメニューを表示する

```html
<div style="--default-contextmenu: show">
    Default menu always shown
</div>
```

<strong>使用例：</strong>テキストエリア、入力フィールド、編集可能なコンテンツ。

### 自動（スマート）モード

```html
<div style="--default-contextmenu: auto">
    Smart context menu
</div>
```

<strong>デフォルトの動作です。</strong>次の場合にデフォルトメニューを表示します。

- テキストが選択されている場合
- テキスト入力フィールド内の場合
- 編集可能なコンテンツ（`contenteditable`）内の場合

それ以外の場合は、デフォルトメニューを非表示にします。

### カスタムメニューとデフォルトメニューの併用

```html
<!-- Custom menu + default menu -->
<textarea style="--custom-contextmenu: editor-menu; --default-contextmenu: show">
    Both menus available
</textarea>
```

**動作：**

1. 最初にカスタムメニューを表示する
2. カスタムメニューが空か見つからない場合は、デフォルトメニューを表示する
3. 両方を併用できる（プラットフォームによって異なる）

## 動的コンテキストメニュー

アプリケーションの状態に応じてメニューを更新します。

### 項目の有効化／無効化

```go
var cutMenuItem *application.MenuItem
var copyMenuItem *application.MenuItem

func createContextMenu() {
    contextMenu := app.ContextMenu.New()
    
    cutMenuItem = contextMenu.Add("Cut")
    cutMenuItem.SetEnabled(false)  // Initially disabled
    cutMenuItem.OnClick(handleCut)
    
    copyMenuItem = contextMenu.Add("Copy")
    copyMenuItem.SetEnabled(false)
    copyMenuItem.OnClick(handleCopy)
    
    app.ContextMenu.Add("editor-menu", contextMenu)
}

func onSelectionChanged(hasSelection bool) {
    cutMenuItem.SetEnabled(hasSelection)
    copyMenuItem.SetEnabled(hasSelection)
    contextMenu.Update()  // Important!
}
```

@note{type="caution" title="必ずUpdate()を呼び出す"}
メニューの状態を変更した後は、<strong>`contextMenu.Update()`</strong>を呼び出してください。これはWindowsでは非常に重要です。

詳細については、[メニューリファレンス](/features/menus/reference/#enabled-state)を参照してください。

@end

### ラベルの変更

```go
playMenuItem := contextMenu.Add("Play")

playMenuItem.OnClick(func(ctx *application.Context) {
    if isPlaying {
        playMenuItem.SetLabel("Pause")
    } else {
        playMenuItem.SetLabel("Play")
    }
    contextMenu.Update()
})
```

### メニューの再構築

大幅な変更を行う場合は、メニュー全体を再構築します。

```go
func rebuildContextMenu(fileType string) {
    contextMenu := app.ContextMenu.New()
    
    // Common items
    contextMenu.Add("Open").OnClick(handleOpen)
    contextMenu.Add("Delete").OnClick(handleDelete)
    
    contextMenu.AddSeparator()
    
    // Type-specific items
    switch fileType {
    case "image":
        contextMenu.Add("Edit Image").OnClick(editImage)
        contextMenu.Add("Set as Wallpaper").OnClick(setWallpaper)
    case "video":
        contextMenu.Add("Play").OnClick(playVideo)
        contextMenu.Add("Extract Audio").OnClick(extractAudio)
    case "document":
        contextMenu.Add("Print").OnClick(printDocument)
        contextMenu.Add("Export PDF").OnClick(exportPDF)
    }
    
    app.ContextMenu.Add("file-menu", contextMenu)
}
```

## プラットフォームごとの動作

コンテキストメニューは<strong>各プラットフォームのネイティブ仕様</strong>に従います。

@tabs{sync-key="platform"}
[macOS]
**macOS ネイティブのコンテキストメニュー：**

- システム標準のアニメーションとトランジション
- 右クリック = Control+Click（自動対応）
- システムの外観（ライト／ダーク）に適応
- デフォルトメニューで標準的なテキスト操作を提供
- 長いメニューではネイティブスクロールを使用

**macOS の慣例：**

- メニュー項目には文頭のみ大文字の表記を使用する
- ダイアログを開く項目には省略記号（...）を使用する
- 一般的なショートカット：⌘C（コピー）、⌘V（貼り付け）

[Windows]
**Windows ネイティブのコンテキストメニュー：**

- Windows ネイティブのスタイル
- Windows のテーマに準拠
- デフォルトメニューで Windows 標準の操作を提供
- タッチ入力とペン入力をサポート

**Windows の慣例：**

- メニュー項目には各単語の先頭を大文字にする表記を使用する
- ダイアログを開く項目には省略記号（...）を使用する
- 一般的なショートカット：Ctrl+C（コピー）、Ctrl+V（貼り付け）

[Linux]
**デスクトップ環境との統合：**

- デスクトップテーマ（GTK、Qt など）に適応
- 右クリックの動作はシステム設定に準拠
- デフォルトメニューの内容は環境によって異なる
- 表示位置はデスクトップ環境の慣例に準拠

**Linux での考慮事項：**

- 対象のデスクトップ環境でテストする
- GTK と Qt では動作が異なる
- 一部のデスクトップ環境ではコンテキストメニューがカスタマイズされる

@end

## 完全な例

**Go コード：**

```go
package main

import (
    "encoding/json"
    "log"
    "github.com/wailsapp/wails/v3/pkg/application"
)

type FileData struct {
    ID   string `json:"id"`
    Type string `json:"type"`
    Name string `json:"name"`
}

func main() {
    app := application.New(application.Options{
        Name: "Context Menu Demo",
    })

    // Create file context menu
    fileMenu := createFileMenu(app)
    app.ContextMenu.Add("file-menu", fileMenu)

    // Create image context menu
    imageMenu := createImageMenu(app)
    app.ContextMenu.Add("image-menu", imageMenu)

    // Create text context menu
    textMenu := createTextMenu(app)
    app.ContextMenu.Add("text-menu", textMenu)

    app.Window.New()
    app.Run()
}

func createFileMenu(app *application.App) *application.ContextMenu {
    menu := app.ContextMenu.New()
    
    menu.Add("Open").OnClick(func(ctx *application.Context) {
        data := parseFileData(ctx.ContextMenuData())
        openFile(data.ID)
    })
    
    menu.Add("Rename").OnClick(func(ctx *application.Context) {
        data := parseFileData(ctx.ContextMenuData())
        renameFile(data.ID)
    })
    
    menu.AddSeparator()
    
    menu.Add("Delete").OnClick(func(ctx *application.Context) {
        data := parseFileData(ctx.ContextMenuData())
        deleteFile(data.ID)
    })
    
    return menu
}

func createImageMenu(app *application.App) *application.ContextMenu {
    menu := app.ContextMenu.New()
    
    menu.Add("View").OnClick(func(ctx *application.Context) {
        data := parseFileData(ctx.ContextMenuData())
        viewImage(data.ID)
    })
    
    menu.Add("Edit").OnClick(func(ctx *application.Context) {
        data := parseFileData(ctx.ContextMenuData())
        editImage(data.ID)
    })
    
    menu.AddSeparator()
    
    exportMenu := menu.AddSubmenu("Export As")
    exportMenu.Add("PNG").OnClick(exportPNG)
    exportMenu.Add("JPEG").OnClick(exportJPEG)
    exportMenu.Add("WebP").OnClick(exportWebP)
    
    return menu
}

func createTextMenu(app *application.App) *application.ContextMenu {
    menu := app.ContextMenu.New()
    
    menu.Add("Cut").SetAccelerator("CmdOrCtrl+X").OnClick(handleCut)
    menu.Add("Copy").SetAccelerator("CmdOrCtrl+C").OnClick(handleCopy)
    menu.Add("Paste").SetAccelerator("CmdOrCtrl+V").OnClick(handlePaste)
    
    return menu
}

func parseFileData(dataStr string) FileData {
    var data FileData
    if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
        log.Printf("Invalid file data: %v", err)
    }
    return data
}

// Handler implementations...
func openFile(id string) { /* ... */ }
func renameFile(id string) { /* ... */ }
func deleteFile(id string) { /* ... */ }
func viewImage(id string) { /* ... */ }
func editImage(id string) { /* ... */ }
func exportPNG(ctx *application.Context) { /* ... */ }
func exportJPEG(ctx *application.Context) { /* ... */ }
func exportWebP(ctx *application.Context) { /* ... */ }
func handleCut(ctx *application.Context) { /* ... */ }
func handleCopy(ctx *application.Context) { /* ... */ }
func handlePaste(ctx *application.Context) { /* ... */ }
```

**HTML：**

```html
<!DOCTYPE html>
<html>
<head>
    <style>
        .file-item {
            padding: 10px;
            margin: 5px;
            border: 1px solid #ccc;
            cursor: pointer;
        }
        
        .file-item:hover {
            background: #f0f0f0;
        }
        
        textarea {
            width: 100%;
            height: 200px;
        }
    </style>
</head>
<body>
    <h2>Files</h2>
    
    <!-- Regular file -->
    <div class="file-item" 
         style="--custom-contextmenu: file-menu; 
                --custom-contextmenu-data: {&quot;id&quot;:&quot;file-1&quot;,&quot;type&quot;:&quot;document&quot;,&quot;name&quot;:&quot;Report.pdf&quot;}">
        📄 Report.pdf
    </div>
    
    <!-- Image file -->
    <div class="file-item" 
         style="--custom-contextmenu: image-menu; 
                --custom-contextmenu-data: {&quot;id&quot;:&quot;file-2&quot;,&quot;type&quot;:&quot;image&quot;,&quot;name&quot;:&quot;Photo.jpg&quot;}">
        🖼️ Photo.jpg
    </div>
    
    <h2>Text Editor</h2>
    
    <!-- Text area with custom menu + default menu -->
    <textarea 
        style="--custom-contextmenu: text-menu; --default-contextmenu: show"
        placeholder="Type here, then right-click...">
    </textarea>
    
    <h2>No Context Menu</h2>
    
    <!-- Disable default menu -->
    <div style="--default-contextmenu: hide; padding: 20px; border: 1px solid #ccc;">
        Right-click here - no menu appears
    </div>
</body>
</html>
```

## ベストプラクティス

### ✅ 推奨事項

- **メニューの目的を絞る** — その要素に関連する操作だけを含める
- **コンテキストデータを検証する** — 信頼できない入力として扱う
- **明確なラベルを使用する** — 「削除」ではなく「ファイルを削除」とする
- **menu.Update() を呼び出す** — メニューの状態を変更した後に呼び出す
- **すべてのプラットフォームでテストする** — 動作はプラットフォームによって異なる
- **キーボードショートカットを用意する** — よく使う操作に設定する
- **関連する項目をグループ化する** — セパレーターを使用する

### ❌ 禁止事項

- **コンテキストデータを信頼しない** — 必ず検証する
- **メニューを長くしすぎない** — 最大 7-10 項目までにする
- **menu.Update() の呼び出しを忘れない** — 忘れるとメニューが正しく動作しない
- **階層を深くしすぎない** — 最大 2 階層までにする
- **専門用語を使用しない** — ユーザーに分かりやすいラベルにする
- **ハンドラーをブロックしない** — 処理を短時間で完了させる

## トラブルシューティング

### コンテキストメニューが表示されない

**考えられる原因：**

1. メニュー ID が一致していない
2. CSS プロパティに入力ミスがある
3. ランタイムが初期化されていない

**解決策：**

```go
// Check menu is registered
app.ContextMenu.Add("my-menu", contextMenu)
```

```html
<!-- Check ID matches -->
<div style="--custom-contextmenu: my-menu">
```

### コンテキストデータを受信できない

**考えられる原因：**

1. CSS プロパティが設定されていない
2. データに特殊文字が含まれている

**解決策：**

```html
<!-- Escape quotes in JSON -->
<div style="--custom-contextmenu-data: {&quot;id&quot;:123}">
```

または、JavaScript を使用します。

```javascript
element.style.setProperty('--custom-contextmenu-data', JSON.stringify(data))
```

### メニュー項目が反応しない

**原因：** 有効化した後に `menu.Update()` を呼び出していない

**解決策：**

```go
menuItem.SetEnabled(true)
contextMenu.Update()  // Add this!
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
★ システムトレイメニュー
システムトレイ／メニューバーとの統合を追加します。

[詳細を見る →](/features/menus/systray/)

---
📖 メニューパターン
一般的なメニューパターンとベストプラクティスを紹介します。

[詳細を見る →](/guides/menus/)

@end

---

**ご質問がありますか？** [Discord](https://discord.gg/JDdSxwjhGf) で質問するか、[コンテキストメニューの例](https://github.com/wailsapp/wails/tree/master/v3/examples/contextmenus)を確認してください。
