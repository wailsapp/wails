---
title: "メニュー"
description: "Wails v3 でメニューを作成およびカスタマイズするためのガイド"
slug: "guides/menus"
sourcePath: "guides/menus.md"
---

Wails v3 は、アプリケーションメニューとコンテキストメニューの両方を作成できる強力なメニューシステムを提供します。このガイドでは、メニューシステムのさまざまな機能を説明します。

## メニューの作成

新しいメニューを作成するには、Menus マネージャーの `New()` メソッドを使用します。

```go
menu := app.Menu.New()
```

### メニュー項目の追加

Wails は、特定の用途に応じた複数の種類のメニュー項目をサポートしています。

#### 通常のメニュー項目

通常のメニュー項目は、メニューの基本的な構成要素です。テキストを表示し、クリック時にアクションを実行できます。

```go
menuItem := menu.Add("Click Me")
```

#### チェックボックス

チェックボックスのメニュー項目には切り替え可能な状態があり、機能や設定の有効化と無効化に便利です。

```go
checkbox := menu.AddCheckbox("My checkbox", true)  // true = initially checked
```

#### ラジオグループ

ラジオグループでは、相互に排他的な選択肢の中から1つを選択できます。ラジオ項目を隣り合わせに配置すると、自動的に作成されます。

```go
menu.AddRadio("Option 1", true)   // true = initially selected
menu.AddRadio("Option 2", false)
menu.AddRadio("Option 3", false)
```

#### セパレーター

セパレーターは、メニュー項目を論理的なグループに整理するための水平線です。

```go
menu.AddSeparator()
```

#### サブメニュー

サブメニューは、メニュー項目にカーソルを合わせるかクリックすると表示される、入れ子になったメニューです。複雑なメニュー構造の整理に便利です。

```go
submenu := menu.AddSubmenu("File")
submenu.Add("Open")
submenu.Add("Save")
```

#### メニューの結合

メニューを末尾または先頭に追加することで、別のメニューに組み込めます。

```go
menu := app.Menu.New()
menu.Add("First Menu")

secondaryMenu := app.Menu.New()
secondaryMenu.Add("Second Menu")

// insert 'secondaryMenu' after 'menu'
menu.Append(secondaryMenu)

// insert 'secondaryMenu' before 'menu'
menu.Prepend(secondaryMenu)

// update the menu
menu.Update()
```

@note{type="info"}
デフォルトでは、`prepend` と `append` は元のメニューと状態を共有します。独自の状態を持つ新しいメニューを作成するには、メニューで `.Clone()` を呼び出します。

例：`menu.Append(secondaryMenu.Clone())`

@end

#### メニューのクリア

メニュー項目の数が変動する場合は、メニュー全体を新しく構築した方がよいことがあります。

これにより、既存のメニューにあるすべての項目がクリアされ、項目を再度追加できるようになります。

```go
menu := app.Menu.New()
menu.Add("Waiting for update...")

// after certain logic, the menu has to be updated
menu.Clear()
menu.Add("Update complete!")
menu.Update()
```

@note{type="info"}
メニューをクリアしても、クリアされるのは最上位のメニュー項目だけです。サブメニューは表示されなくなりますが、引き続きメモリを占有するため、メニューは慎重に管理してください。

@end

#### メニューの破棄

メニューをクリアして解放するには、`Destroy()` メソッドを使用します。

```go
menu := app.Menu.New()
menu.Add("Waiting for update...")

// after certain logic, the menu has to be destroyed
menu.Destroy()
```

### メニュー項目のプロパティ

メニュー項目では、複数のプロパティを設定できます。

| プロパティ | メソッド | 説明 |
| --- | --- | --- |
| ラベル | `SetLabel(string)` | 表示テキストを設定します |
| 有効 | `SetEnabled(bool)` | 項目を有効化または無効化します |
| チェック状態 | `SetChecked(bool)` | チェック状態を設定します（チェックボックス／ラジオ項目用） |
| ツールチップ | `SetTooltip(string)` | ツールチップのテキストを設定します |
| 非表示 | `SetHidden(bool)` | 項目を表示または非表示にします |
| アクセラレーター | `SetAccelerator(string)` | キーボードショートカットを設定します |

### メニュー項目の状態

メニュー項目には、表示と操作の可否を制御するさまざまな状態があります。

#### 表示状態

`SetHidden()` メソッドを使用して、メニュー項目の表示と非表示を動的に切り替えられます。

```go
menuItem := menu.Add("Dynamic Item")

// Hide the menu item
menuItem.SetHidden(true)

// Show the menu item
menuItem.SetHidden(false)

// Check current visibility
isHidden := menuItem.Hidden()
```

非表示にしたメニュー項目は、再表示されるまでメニューから完全に取り除かれます。これは、アプリケーションが特定の状態にある場合にのみ表示すべき、状況に応じたメニュー項目に便利です。

#### 有効状態

`SetEnabled()` メソッドを使用して、メニュー項目を有効または無効にできます。

```go
menuItem := menu.Add("Save")

// Disable the menu item
menuItem.SetEnabled(false)  // Item appears grayed out and cannot be clicked

// Enable the menu item
menuItem.SetEnabled(true)   // Item becomes clickable again

// Check current enabled state
isEnabled := menuItem.Enabled()
```

無効にしたメニュー項目は表示されたままですが、グレー表示になり、クリックできません。これは、次のようにアクションが現在利用できないことを示すためによく使用されます。

- 保存する変更がない場合に「保存」を無効にする
- 何も選択されていない場合に「コピー」を無効にする
- 元に戻す操作がない場合に「元に戻す」を無効にする

#### 動的な状態管理

これらの状態をイベントハンドラーと組み合わせて、動的なメニューを作成できます。

```go
saveMenuItem := menu.Add("Save")

// Initially disable the Save menu item
saveMenuItem.SetEnabled(false)

// Enable Save only when there are unsaved changes
documentChanged := func() {
    saveMenuItem.SetEnabled(true)
    menu.Update()  // Remember to update the menu after changing states
}

// Disable Save after saving
documentSaved := func() {
    saveMenuItem.SetEnabled(false)
    menu.Update()
}
```

### イベント処理

メニュー項目では、`OnClick` メソッドを使用してクリックイベントを処理できます。

```go
menuItem.OnClick(func(ctx *application.Context) {
    // Handle the click event
    println("Menu item clicked!")
})
```

コンテキストには、クリックされたメニュー項目に関する情報が含まれます。

```go
menuItem.OnClick(func(ctx *application.Context) {
    // Get the clicked menu item
    clickedItem := ctx.ClickedMenuItem()
    // Get its current state
    isChecked := clickedItem.Checked()
})
```

### ロールベースのメニュー項目

Wails には、標準機能を備えたメニュー項目を自動的に作成する、定義済みのメニューロールが用意されています。サポートされているメニューロールは次のとおりです。

#### 完全なメニュー構造

これらのロールは、一般的な機能を備えたメニュー構造全体を作成します。

| ロール | 説明 | プラットフォームに関する注記 |
| --- | --- | --- |
| `AppMenu` | 「このアプリケーションについて」、サービス、非表示／表示、終了を含むアプリケーションメニュー | macOS のみ |
| `EditMenu` | 元に戻す、やり直し、切り取り、コピー、貼り付けなどを含む標準の編集メニュー | すべてのプラットフォーム |
| `ViewMenu` | 再読み込み、ズーム、フルスクリーンの各コントロールを含む表示メニュー | すべてのプラットフォーム |
| `WindowMenu` | ウィンドウのコントロール（最小化、ズームなど） | すべてのプラットフォーム |
| `HelpMenu` | Wails Web サイトへの「詳細」リンクを含むヘルプメニュー | すべてのプラットフォーム |

#### 個別のメニュー項目

これらのロールを使用して、個別のメニュー項目を追加できます。

| ロール | 説明 | プラットフォームに関する注記 |
| --- | --- | --- |
| `About` | アプリケーションの「このアプリケーションについて」ダイアログを表示 | すべてのプラットフォーム |
| `Hide` | アプリケーションを非表示 | macOS のみ |
| `HideOthers` | ほかのアプリケーションを非表示 | macOS のみ |
| `UnHide` | 非表示のアプリケーションを表示 | macOS のみ |
| `CloseWindow` | 現在のウィンドウを閉じる | すべてのプラットフォーム |
| `Minimise` | ウィンドウを最小化 | すべてのプラットフォーム |
| `Zoom` | ウィンドウをズーム | macOS のみ |
| `Front` | ウィンドウを最前面に移動 | macOS のみ |
| `Quit` | アプリケーションを終了 | すべてのプラットフォーム |
| `Undo` | 直前の操作を元に戻す | すべてのプラットフォーム |
| `Redo` | 直前の操作をやり直す | すべてのプラットフォーム |
| `Cut` | 選択範囲を切り取り | すべてのプラットフォーム |
| `Copy` | 選択範囲をコピー | すべてのプラットフォーム |
| `Paste` | クリップボードから貼り付け | すべてのプラットフォーム |
| `PasteAndMatchStyle` | スタイルを合わせて貼り付け | macOS のみ |
| `SelectAll` | すべて選択 | すべてのプラットフォーム |
| `Delete` | 選択範囲を削除 | すべてのプラットフォーム |
| `Reload` | 現在のページを再読み込み | すべてのプラットフォーム |
| `ForceReload` | 現在のページを強制再読み込み | すべてのプラットフォーム |
| `ToggleFullscreen` | 全画面表示を切り替え | すべてのプラットフォーム |
| `ResetZoom` | ズームレベルをリセット | すべてのプラットフォーム |
| `ZoomIn` | 拡大 | すべてのプラットフォーム |
| `ZoomOut` | 縮小 | すべてのプラットフォーム |

完全なメニューと個別のロールの両方を使用する例を次に示します。

```go
menu := app.Menu.New()

// Add complete menu structures
menu.AddRole(application.AppMenu)    // macOS only
menu.AddRole(application.EditMenu)   // Common edit operations
menu.AddRole(application.ViewMenu)   // View controls
menu.AddRole(application.WindowMenu) // Window controls

// Add individual role-based items to a custom menu
fileMenu := menu.AddSubmenu("File")
fileMenu.AddRole(application.CloseWindow)
fileMenu.AddSeparator()
fileMenu.AddRole(application.Quit)
```

## アプリケーションメニュー

アプリケーションメニューは、アプリケーションウィンドウの上部（Windows/Linux）または画面上部（macOS）に表示されるメニューです。

### アプリケーションメニューの動作

`app.Menu.Set()` を使用してアプリケーションメニューを設定すると、macOS ではメインメニューになります。 Windows/Linux では、メニューはウィンドウごとに設定されます。

```go
app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Custom Menu Window",
    Windows: application.WindowsWindow{
        Menu: customMenu,  // Override application menu for this window
    },
})
```

これらの異なるメニュー動作を示す完全な例を次に示します。

```go
func main() {
    app := application.New(application.Options{})

    // Create application menu
    appMenu := app.Menu.New()
    fileMenu := appMenu.AddSubmenu("File")
    fileMenu.Add("New").OnClick(func(ctx *application.Context) {
        // This will be available in all windows unless overridden
        window := app.Window.Current()
        window.SetTitle("New Window")
    })
    
    // Set as application menu - this is for macOS
    app.Menu.Set(appMenu)

    // Window with custom menu on Windows
    customMenu := app.Menu.New()
    customMenu.Add("Custom Action")
    app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title: "Custom Menu",
        Windows: application.WindowsWindow{
            Menu: customMenu,
        },
    })

    app.Run()
}
```

## コンテキストメニュー

コンテキストメニューは、アプリケーション内の要素を右クリックしたときに表示されるポップアップメニューです。クリックした要素に関連する操作へすばやくアクセスできます。

### デフォルトのコンテキストメニュー

デフォルトのコンテキストメニューは WebView に組み込まれたコンテキストメニューで、次のようなシステムレベルの操作を提供します。

- テキスト操作のためのコピー、切り取り、貼り付け
- テキスト選択コントロール
- スペルチェックのオプション

#### デフォルトのコンテキストメニューの制御

`--default-contextmenu` CSS プロパティを使用して、デフォルトのコンテキストメニューを表示するタイミングを制御できます。

```html
<!-- Always show default context menu -->
<div style="--default-contextmenu: show">
    <input type="text" placeholder="Right-click for text operations"/>
    <textarea>Standard text operations available here</textarea>
</div>

<!-- Hide default context menu -->
<div style="--default-contextmenu: hide">
    <div class="custom-component">Custom context menu only</div>
</div>

<!-- Smart context menu behaviour (default) -->
<div style="--default-contextmenu: auto">
    <!-- Shows default menu when text is selected or in input fields -->
    <p>Select this text to see the default menu</p>
    <input type="text" placeholder="Default menu for input operations"/>
</div>
```

@note{type="info"}
この機能が想定どおりに動作するのは、[フロントエンドランタイムの準備が完了](/reference/frontend-runtime/)した後だけです。

@end

#### ネストされたコンテキストメニューの動作

ネストされた要素で `--default-contextmenu` プロパティを使用する場合は、次のルールが適用されます。

1. 明示的に上書きしない限り、子要素は親要素のコンテキストメニュー設定を継承します
2. 最も具体的な（最も近い）設定が優先されます
3. `auto` 値を使用すると、デフォルトの動作にリセットできます

ネストされたコンテキストメニューの動作例：

```html
<!-- Parent sets hide -->
<div style="--default-contextmenu: hide">
    <!-- This inherits hide -->
    <p>No context menu here</p>
    
    <!-- This overrides to show -->
    <div style="--default-contextmenu: show">
        <p>Context menu shown here</p>
        
        <!-- This inherits show -->
        <span>Also has context menu</span>
        
        <!-- This resets to automatic behaviour -->
        <div style="--default-contextmenu: auto">
            <p>Shows menu only when text is selected</p>
        </div>
    </div>
</div>
```

### カスタムコンテキストメニュー

カスタムコンテキストメニューを使用すると、クリックされた要素に関連するアプリケーション固有の操作を提供できます。特に、次の用途に役立ちます。

- ドキュメントマネージャーでのファイル操作
- 画像編集ツール
- データグリッドでのカスタム操作
- コンポーネント固有の操作

#### カスタムコンテキストメニューの作成

カスタムコンテキストメニューを作成するときは、メニューを HTML 要素に関連付ける一意の識別子（名前）を指定します。

```go
// Create a context menu with identifier "imageMenu"
contextMenu := app.ContextMenu.New()
app.ContextMenu.Add("imageMenu", contextMenu)
```

name パラメーター（この例では「imageMenu」）は、次の目的で使用する一意の識別子です。

1. HTML 要素をこの特定のコンテキストメニューに関連付ける
2. 右クリック時に表示するメニューを特定する
3. メニューの更新とクリーンアップを可能にする

#### コンテキストデータ

コンテキストメニューイベントを処理するときは、クリックされたメニュー項目と、それに関連付けられたコンテキストデータの両方にアクセスできます。

```go
contextMenu.Add("Process").OnClick(func(ctx *application.Context) {
    // Get the clicked menu item
    menuItem := ctx.ClickedMenuItem()
    
    // Get the context data as a string
    contextData := ctx.ContextMenuData()
    
    // Check if the menu item is checked (for checkbox/radio items)
    isChecked := ctx.IsChecked()
    
    // Use the data
    if contextData != "" {
        processItem(contextData)
    }
})
```

コンテキストデータは HTML 要素の `--custom-contextmenu-data` プロパティから渡され、クリックハンドラーでは `ctx.ContextMenuData()` を介して利用できます。これは特に、次の場合に役立ちます。

- 各項目を一意に識別する必要があるリストやグリッドを扱う場合
- 特定のコンポーネントや要素に対する操作を処理する場合
- 状態またはメタデータをフロントエンドからバックエンドに渡す場合

#### コンテキストメニューの管理

コンテキストメニューを変更した後は、`Update()` メソッドを呼び出して変更を適用します。

```go
contextMenu.Update()
```

コンテキストメニューが不要になったら、破棄できます。

```go
contextMenu.Destroy()
```

@note{type="danger" title="警告"}
`Destroy()` を呼び出した後にコンテキストメニューの参照を再度使用すると、パニックが発生します。

@end

### 実践例：画像ギャラリー

画像ギャラリー用のカスタムコンテキストメニューを実装する完全な例を次に示します。

```go
// Backend: Create the context menu
imageMenu := app.ContextMenu.New()
app.ContextMenu.Add("imageMenu", imageMenu)

// Add relevant operations
imageMenu.Add("View Full Size").OnClick(func(ctx *application.Context) {
    // Get the image ID from context data
    if imageID := ctx.ContextMenuData(); imageID != "" {
        openFullSizeImage(imageID)
    }
})

imageMenu.Add("Download").OnClick(func(ctx *application.Context) {
    if imageID := ctx.ContextMenuData(); imageID != "" {
        downloadImage(imageID)
    }
})

imageMenu.Add("Share").OnClick(func(ctx *application.Context) {
    if imageID := ctx.ContextMenuData(); imageID != "" {
        showShareDialog(imageID)
    }
})
```

```html
<!-- Frontend: Image gallery implementation -->
<div class="gallery">
    <!-- Each image container with context menu -->
    <div class="image-container" 
         style="--custom-contextmenu: imageMenu; --custom-contextmenu-data: img_123">
        <img src="/images/img_123.jpg" alt="Gallery Image"/>
        <span class="caption">Nature Photo</span>
    </div>
    
    <div class="image-container"
         style="--custom-contextmenu: imageMenu; --custom-contextmenu-data: img_124">
        <img src="/images/img_124.jpg" alt="Gallery Image"/>
        <span class="caption">City Photo</span>
    </div>
</div>
```

この例では、次のように動作します。

1. コンテキストメニューは識別子「imageMenu」で作成されます
2. 各画像コンテナは `--custom-contextmenu: imageMenu` を使用してメニューに関連付けられます
3. 各コンテナは `--custom-contextmenu-data` を使用して、その画像 ID をコンテキストデータとして提供します
4. バックエンドはクリックハンドラーで画像 ID を受け取り、その画像に固有の操作を実行できます
5. すべての画像で同じメニューを再利用しますが、操作対象の画像はコンテキストデータによって判別できます

このパターンは、特に次の用途で有効です。

- 行ごとに固有の操作が必要なデータグリッド
- ファイルごとにコンテキストに応じたアクションが必要なファイルマネージャー
- 要素ごとに異なる操作が必要なデザインツール
- 同じ操作を複数のインスタンスに適用するあらゆるコンポーネント
