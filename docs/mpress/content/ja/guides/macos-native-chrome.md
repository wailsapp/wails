---
title: "macOS のネイティブウィンドウ装飾"
description: "Wails ウィンドウにネイティブの AppKit ツールバー、サイドバー、コンテンツリスト、インスペクタ、アクセサリ、ウィンドウタブを構築する"
slug: "guides/macos-native-chrome"
sourcePath: "guides/macos-native-chrome.md"
---

対象プラットフォーム: macOS

Wails v3 では、WebView を本物の AppKit ウィンドウ装飾で囲めます。ツールバー、サイドバー、コンテンツリスト、インスペクタ、タイトルバーの各領域は、Go から作成するネイティブコントロールです。HTML は使われないため、AppKit のマテリアル、キーボード操作、アニメーション、状態の永続化をそのまま利用でき、フロントエンドはコンテンツに集中できます。

このページの API はすべて `application` パッケージにあり、名前に `Mac` が付いています。同じコードは Windows と Linux でもコンパイルできます。コンストラクタとセッターはどの環境でも動作しますが、アタッチ操作は何もしないかエラーを返し、ウィンドウは通常どおり単一の WebView を保持します。

## ウィンドウの構成

すべての要素を備えたウィンドウは、先頭側から末尾側に向かって次の部分で構成されます。

| 部分 | 型 | AppKit クラス |
|------|------|--------------|
| ツールバー | `MacToolbar` | `NSToolbar` |
| サイドバー | `MacSidebar` | サイドバーの分割項目内にあるソースリストの `NSOutlineView` |
| コンテンツリスト | `MacContentList` | コンテンツリストの分割項目内にある `NSTableView` |
| メインコンテンツ | WebView、または `MacTextEditor` | `WKWebView` または `NSTextView` |
| インスペクタ | `MacInspector` | インスペクタの分割項目内にあるネイティブのプロパティコントロール |
| アクセサリ | `MacAccessory` | `NSTitlebarAccessoryViewController` または `NSSplitViewItemAccessoryViewController` |

ペインは `NSSplitViewController` である `MacSplitView` によって配置されます。まず各部を作成し、先頭側から末尾側の順に分割ビューへ追加します。次に分割ビューをウィンドウへ、最後にツールバーをアタッチします。`SetSplitView` と `SetToolbar` を使うか、`Mac.SplitView` と `Mac.Toolbar` のウィンドウオプションを使って一度に設定できます。

一般的なウィンドウ設定では、統合されたツールバーの下までコンテンツをスクロールできるよう、装飾と次のウィンドウオプションを組み合わせます。

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Mac: application.MacWindow{
        ContentLayout: application.MacContentLayoutEdgeToEdge,
        TitleBar: application.MacTitleBar{
            FullSizeContent:      true,
            HideToolbarSeparator: true,
            ToolbarStyle:         application.MacToolbarStyleUnified,
        },
    },
})
```

## ツールバー

`NewMacToolbar` は `NSToolbar` を作成します。`Add` メソッドで項目を追加し、返されたハンドルに対してセッターとコールバックを連続して設定してから、`SetToolbar` でツールバーをアタッチします。識別子は自動生成されます。

```go
toolbar := application.NewMacToolbar().
    SetDisplayMode(application.MacToolbarDisplayModeIconOnly)

// Standard AppKit items. They have no handle because AppKit owns them.
toolbar.AddSidebarToggle()
toolbar.AddSidebarTrackingSeparator()

toolbar.AddButton("New").
    SetSymbol("square.and.pencil").
    SetTooltip("Create a note").
    SetBordered(true).
    OnClick(func(*application.Context) {
        // create a note
    })

toolbar.AddSearch("Search").
    SetSearchPlaceholder("Search notes").
    SetSearchIncremental(true).
    OnSearch(func(_ *application.Context, query string) {
        filterNotes(query)
    })

toolbar.AddFlexibleSpace()

mode := toolbar.AddGroup("Mode", application.ToolbarGroupSelectOne)
mode.AddButton("Write").SetSymbol("pencil").OnClick(func(*application.Context) {
    mode.SetSelectedIndex(0)
})
mode.AddButton("Preview").SetSymbol("doc.richtext").OnClick(func(*application.Context) {
    mode.SetSelectedIndex(1)
})

actions := application.NewMenu()
actions.Add("Export...").OnClick(func(*application.Context) {})
toolbar.AddMenu("Actions", actions).SetSymbol("ellipsis.circle")

window.SetToolbar(toolbar)
```

項目の種類は次のとおりです。

- `AddButton` はプッシュボタンを追加します。ツールバーをアタッチする前に、すべてのボタンに `OnClick` が必要です。設定されていない場合、`SetToolbar` はエラーを報告し、以前のツールバーをそのまま残します。
- `AddSearch` は `NSSearchToolbarItem` を追加します。`OnSearch` が必須です。`SetSearchPlaceholder`、`SetSearchIncremental`、最近の検索を永続化するメニュー用の `SetSearchRecentsKey`、虫眼鏡の背後にあるカスタムメニュー用の `SetSearchMenu` を使用します。
- `AddShare` はシステムの共有項目を追加します（後述）。
- `AddGroup` はセグメント化された `NSToolbarItemGroup` を追加します。グループの `AddButton` でメンバーを追加し、`ToolbarGroupSelectOne`、`ToolbarGroupMomentary`、`ToolbarGroupSelectAny` のいずれかを選びます。
- `AddMenu` は通常の `Menu` で構成するドロップダウンの `NSMenuToolbarItem` を追加します。`SetShowsIndicator(false)` で山形アイコンを非表示にします。
- `AddSpace` と `AddFlexibleSpace` は標準のスペーサーを追加します。
- `AddSidebarToggle` と `AddSidebarTrackingSeparator` は AppKit のサイドバー項目を追加します。セパレーターは、それより前の項目をサイドバーの境界線の上に揃えるため、ウィンドウにはサイドバーペインを含む分割ビューが必要です。
- `AddInspectorToggle` と `AddInspectorTrackingSeparator` はインスペクタペインに対して同じ機能を提供します。

`SetDisplayMode` では、`MacToolbarDisplayModeIconAndLabel`（デフォルト）、`MacToolbarDisplayModeIconOnly`、`MacToolbarDisplayModeLabelOnly`、`MacToolbarDisplayModeDefault` から選択します。

### 実行中の更新

各ハンドルは実行中の更新にも使えます。アタッチ後に呼び出したセッターは、アプリケーションスレッド上でネイティブ項目を更新します。

```go
save := toolbar.AddButton("Save").SetSymbol("checkmark.circle")
save.OnClick(func(*application.Context) {
    save.SetBadgeCount(0).SetProminent(false)
})

// Later, when the document changes:
save.SetBadgeCount(1).SetProminent(true)
save.SetVisibilityPriority(application.MacToolbarVisibilityPriorityHigh)

toolbar.Move(save, 0)
toolbar.Remove(save)
```

ほかに便利なセッターとして、`SetLabel`、`SetTooltip`、`SetEnabled`、`SetHidden`、`SetTintColor`、`SetNavigational` があります。`SetNavigational` は、Safari の戻るボタンや進むボタンと同様に、項目を先頭側に配置し続けます。`SetVisibilityPriority` は、ウィンドウが狭くなったときに、どの項目を先にオーバーフローメニューへ移すかを決めます。

### ユーザーによるカスタマイズ

`SetCustomizable` は標準の「ツールバーをカスタマイズ...」シートを有効にし、指定したキーでユーザーの配置を保存します。各項目に安定した `SetPersistenceKey` を設定すると、保存した配置がアプリの再起動後も維持されます。ユーザーが追加したときだけ表示する項目には `SetInDefaultSet(false)` を使用します。ツールバーをアタッチする前に `SetCustomizable` を呼び出してください。

```go
toolbar := application.NewMacToolbar().SetCustomizable("myapp.main-toolbar")

newNote := toolbar.AddButton("New").
    SetPersistenceKey("new").
    OnClick(func(*application.Context) {})

// Offered in the palette but hidden until the user adds it.
toolbar.AddButton("Archive").
    SetPersistenceKey("archive").
    SetInDefaultSet(false).
    OnClick(func(*application.Context) {})

toolbar.SetCenteredItems(newNote)

// From a menu item, for example:
toolbar.RunCustomizationPalette()
```

### 共有

`AddShare` は `MacToolbarShareItem` を返します。`MacShareProvider` が少なくとも 1 つの表現形式を提供するまで、この項目は無効のままです。Wails は共有サービスが要求したときにだけプロバイダーからバイト列を取得するため、大きなエクスポートは遅延して生成されます。`MacShareProviderFunc` は 2 つの関数をプロバイダーに変換します。状態を持つアプリケーションでは、インターフェースを直接実装できます。

```go
share := toolbar.AddShare("Share")
share.SetProvider(application.MacShareProviderFunc{
    Available: []application.MacShareRepresentation{
        {ContentType: application.MacShareTypePDF},
        {ContentType: application.MacShareTypePlainText},
    },
    Load: func(request application.MacShareRequest) ([]byte, error) {
        if request.ContentType == application.MacShareTypePDF {
            return renderPDF()
        }
        return []byte(currentText()), nil
    },
}).SetSubject("Saturday, slowly").SetSuggestedName("Note")

share.OnShared(func(_ *application.Context, service string) {
    // service is the localised name of the sharing service
})
share.OnShareError(func(_ *application.Context, service string, err error) {
    window.Error("share via %s failed: %s", service, err)
})
```

一般的なコンテンツ型は `MacShareTypePlainText`、`MacShareTypeHTML`、`MacShareTypePDF`、`MacShareTypePNG`、`MacShareTypeJPEG` です。そのほかの UTI 文字列も使用できます。

### アタッチとデタッチ

ツールバーは一度に 1 つのウィンドウに属します。`WebviewWindow.SetToolbar` と `NativeWindow.SetToolbar` は、ウィンドウの作成前でも作成後でもツールバーを受け取ります。`nil` を渡すとツールバーが取り外され、別の場所で使えるようになります。`WebviewWindow` の `SetToolbar` は検証エラーを `Window.Error` を通じて報告し、`NativeWindow` 版はエラーを返します。`Mac.Toolbar` ウィンドウオプションは作成時にツールバーをアタッチします。トラッキングセパレーターが揃える対象のサイドバーを見つけられるよう、`Mac.SplitView` の後に適用されます。

## 分割ビュー

`MacSplitView` はペインを配置します。先頭側から末尾側の順に追加してください。`AddSidebar`、`AddContentList`、`AddInspector` は、収容するネイティブモデルを受け取り、サイズと折りたたみを制御する `MacSplitPane` を返します。`AddPrimaryContent` はウィンドウの既存の WebView を配置し、`SetContentLayout` を備えた `MacSplitWebviewPane` を返します。

```go
split := application.NewMacSplitView().SetAutosaveName("myapp.main-window")

sidebarPane := split.AddSidebar(sidebar).
    SetMinimumThickness(200).
    SetMaximumThickness(320).
    SetCollapsible(true)

split.AddContentList(list).
    SetMinimumThickness(240).
    SetCollapsible(true)

split.AddPrimaryContent().
    SetContentLayout(application.MacContentLayoutEdgeToEdge)

split.AddInspector(inspector).
    SetPreferredThicknessFraction(0.25).
    SetHoldingPriority(300).
    SetCollapsible(true).
    SetCanCollapseFromWindowResize(false).
    SetCollapsed(true)

sidebarPane.OnCollapsedChange(func(_ *application.Context, collapsed bool) {
    // fired for toggles, gestures, menu items and SetCollapsed alike
})

window.SetSplitView(split)

// At runtime:
sidebarPane.Toggle()
```

`SetSplitView` はネイティブウィンドウの作成前でも作成後でも使えます。作成前に呼び出すと、レイアウトは保留され、ウィンドウの作成時にインストールされます。作成後、たとえば実行中のアプリケーションのメニューやトレイのコールバックから呼び出すと、レイアウトは直ちにインストールされます。ウィンドウの既存の WebView がメインペインになり、トラッキングセパレーターが揃うよう現在のツールバーが再アタッチされ、保留中のアクセサリもアタッチされます。

`app.Run` の後に作成するウィンドウは、`Mac.SplitView` と `Mac.Toolbar` オプションを使って 1 回の呼び出しで設定できます。

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Notes",
    URL:   "/",
    Mac: application.MacWindow{
        SplitView: split,
        Toolbar:   toolbar,
    },
})
```

レイアウトの規則は次のとおりです。

- レイアウトには少なくとも 2 つのペインと、ちょうど 1 つのメインペインが必要です（`WebviewWindow` では `AddPrimaryContent`、`NativeWindow` では `AddTextEditor`）。
- コンテンツリストは最大 1 つで、サイドバーの後、メインペインの前に配置します。
- 分割ビューをアタッチすると、ペイン構造は固定されます。ペインの設定、折りたたみ状態、サイドバー、リスト、インスペクタの内容は、その後もいつでも変更できます。
- インストール済みのレイアウトは置き換えられません。同じウィンドウで 2 回目の `SetSplitView` を呼び出すと `ErrMacSplitViewAlreadyInstalled` が報告され（`WebviewWindow` では `Window.Error` を通じて、`NativeWindow` では戻り値として）、ウィンドウは変更されません。インストール前に `nil` を渡すと、保留中のレイアウトが消去されます。
- サイドバー、リスト、インスペクタ、分割ビューは、それぞれ一度に 1 つのウィンドウに属します。

ペインのセッターは `SetMinimumThickness`、`SetMaximumThickness`、`SetPreferredThicknessFraction`、`SetHoldingPriority`、`SetCollapsible`、`SetCanCollapseFromWindowResize`、`SetCollapsed`、`Toggle`、`IsCollapsed`、`OnCollapsedChange` です。`SetAutosaveName` は、起動をまたいで境界線の位置を保存します。

メインペインの `SetContentLayout` では、`MacContentLayoutBelowToolbar` と `MacContentLayoutEdgeToEdge` のいずれかを選びます。`MacContentLayoutAutomatic` は `MacWindow.ContentLayout` を継承し、さらに `TitleBar.FullSizeContent` に従います。端から端まで広げる配置にすると、macOS 26 以降では AppKit がツールバーの下でスクロール端のエフェクトを適用できます。

## サイドバー

`MacSidebar` はネイティブのソースリストです。ルート行、セクション、任意の深さまで入れ子にした行を保持します。後で行を更新できるよう、返されたハンドルを保持してください。

```go
sidebar := application.NewMacSidebar()

// A root row above the sections.
sidebar.AddItem("All Notes").
    SetSymbol("tray.full").
    SetBadge(12).
    OnClick(func(*application.Context) {})

notes := sidebar.AddSection("Notes")
draft := notes.AddItem("Saturday, slowly").
    SetSymbol("doc.text").
    SetAccessorySymbol("pin.fill").
    SetTooltip("Field notes").
    SetEditable(true).
    OnRename(func(_ *application.Context, label string) {
        // the row label is already updated
    })

// Nested rows with a tinted symbol.
tags := sidebar.AddSection("Tags")
tint := application.NewRGB(0, 122, 255)
work := tags.AddItem("Work").SetSymbol("briefcase").SetTintColor(&tint)
work.AddItem("Meetings").SetSymbol("tag")
work.SetExpanded(true)

sidebar.SetSelectedItem(draft)
```

行のセッターは `SetLabel`、`SetSymbol`、`SetTooltip`、`SetEnabled`、`SetHidden`、`SetBadge`、`SetAccessorySymbol`、`SetTintColor`、`SetEditable`、`SetExpanded` です。`OnClick` は AppKit が行を選択したとき、`OnExpandedChange` はユーザーが入れ子の行を開閉したとき、`OnRename` はインラインでの名前変更が確定した後に呼び出されます。

### 選択

デフォルトは単一選択です。`SetSelectedItem` は `OnClick` を発火させずに行を選択します。複数選択が有効な場合も、クリックされた行に対して `OnClick` が発火し、`OnSelectionChange` は選択された行の集合全体を報告します。

```go
sidebar.SetAllowsMultipleSelection(true)
sidebar.OnSelectionChange(func(_ *application.Context, items []*application.MacSidebarItem) {
    for _, item := range items {
        _ = item.Section()
    }
})
```

### コンテキストメニュー

右クリック時には、行自身の `SetContextMenu`、サイドバーの `OnContextMenu` コールバック、サイドバーの代替 `SetContextMenu` の順にメニューを探します。コールバックは AppKit が待機する間にアプリケーションスレッドで実行されるため、短時間で処理してください。

```go
sidebar.OnContextMenu(func(_ *application.Context, item *application.MacSidebarItem) *application.Menu {
    if item == nil {
        return nil // fall back to the menu set with SetContextMenu
    }
    menu := application.NewMenu()
    menu.Add("Delete").OnClick(func(*application.Context) { item.Remove() })
    return menu
})

empty := application.NewMenu()
empty.Add("New Note").OnClick(func(*application.Context) {})
sidebar.SetContextMenu(empty)
```

### ドラッグによる並べ替え

`SetReorderable` を使うと、ユーザーはセクション内、セクション間、ルートとの間で行をドラッグできます。`OnMove` が発火する前に Go のモデルが更新されます。

```go
sidebar.SetReorderable(true)
sidebar.OnMove(func(_ *application.Context, item *application.MacSidebarItem, section *application.MacSidebarSection, index int) {
    // section is nil when the row was dropped at the sidebar root
})
```

### 削除

行を削除すると、入れ子の行も削除されます。その後、ハンドルは機能しなくなります。

```go
notes.Remove(draft)
sidebar.RemoveSection(tags)
```

## コンテンツリスト

`MacContentList` は Finder、Mail、ドキュメントブラウザの中央の列に相当します。列を設定しない場合、タイトル、サブタイトル、先頭側のシンボル、末尾側の詳細、件数バッジを含む情報豊富な行を表示します。

```go
list := application.NewMacContentList().
    SetStyle(application.MacContentListStyleInset).
    SetEmptyText("No notes match")

row := list.AddRow("Saturday, slowly").
    SetSubtitle("A slow day is still a day well spent.").
    SetDetail("Yesterday").
    SetSymbol("doc.text").
    SetBadge(2)

list.OnSelectionChange(func(_ *application.Context, rows []*application.MacContentListRow) {})
list.OnActivate(func(_ *application.Context, row *application.MacContentListRow) {
    // double-click or Return
})
list.SetSelectedRow(row)
```

スタイルは `MacContentListStyleAutomatic`、`MacContentListStyleInset`、`MacContentListStyleSourceList`、`MacContentListStylePlain`、`MacContentListStyleFullWidth` です。表示に関するその他のオプションは `SetRowHeight`、`SetAlternatingRowBackgrounds`、`SetHeaderVisible`、`SetAllowsMultipleSelection` です。行では `SetHidden`、`SetEnabled`、`SetTooltip`、指定位置への `InsertRow`、`Remove`、`RemoveAll` が使えます。

### 列

`SetColumns` はテーブルモードに切り替えます。その後、行は列ごとに 1 つずつ `SetCells` の値を表示します。`Sortable` と指定された列には並べ替えインジケーターが表示され、`SetSortable` はヘッダーのクリックによる並べ替えを有効にします。`OnSort` コールバックがない場合、リストは `SortBy` を使って列のテキストで自動的に並べ替えます。

```go
table := application.NewMacContentList().
    SetColumns(
        application.MacContentListColumn{Title: "Name", Width: 220, Sortable: true},
        application.MacContentListColumn{Title: "Size", Width: 80, Alignment: application.MacContentListAlignTrailing},
        application.MacContentListColumn{Title: "Modified", Sortable: true},
    ).
    SetSortable(true).
    SetAlternatingRowBackgrounds(true)

table.AddRow("notes.txt").SetCells("notes.txt", "4 KB", "Today")
table.AddRow("ideas.txt").SetCells("ideas.txt", "1 KB", "Yesterday")

// Without OnSort the list sorts itself by the column text.
table.OnSort(func(_ *application.Context, column int, ascending bool) {
    table.SortRows(func(a, b *application.MacContentListRow) bool {
        return (a.Cells()[column] < b.Cells()[column]) == ascending
    })
})
```

### コンテキストメニュー

コンテキストメニューはサイドバーと同じ順序で決まります。行の `SetContextMenu`、`OnContextMenu`、リストの代替メニューの順です。

```go
list.OnContextMenu(func(_ *application.Context, row *application.MacContentListRow) *application.Menu {
    if row == nil {
        return nil
    }
    menu := application.NewMenu()
    menu.Add("Delete").OnClick(func(*application.Context) { row.Remove() })
    return menu
})
```

## インスペクタ

`MacInspector` はネイティブコントロールで構成する末尾側のプロパティパネルで、セクションごとにグループ化されます。各 `Add` メソッドは、種類に応じたセッターとコールバックを持つ `MacInspectorControl` ハンドルを返します。

```go
inspector := application.NewMacInspector()

document := inspector.AddSection("Document")
document.AddTextField("Title", "Saturday, slowly").
    OnTextChange(func(_ *application.Context, value string) {})
document.AddPopup("Category", []string{"Personal", "Work"}, 0).
    OnSelectionChange(func(_ *application.Context, index int, value string) {})
document.AddCheckbox("Pinned", false).
    OnToggle(func(_ *application.Context, checked bool) {})

appearance := inspector.AddSection("Appearance").SetCollapsible(true)
priority := appearance.AddSlider("Priority", 0, 5, 0)
priority.OnValueChange(func(_ *application.Context, value float64) {})
appearance.AddStepper("Indent", 0, 8, 1, 2)
appearance.AddSegmented("Align", []string{"Left", "Centre", "Right"}, 0)
appearance.AddColorWell("Tint", application.NewRGB(0, 122, 255)).
    OnColorChange(func(_ *application.Context, colour application.RGBA) {})
appearance.AddDatePicker("Due", time.Time{}).
    OnDateChange(func(_ *application.Context, t time.Time) {})
appearance.AddButton("Reset").OnClick(func(*application.Context) {
    priority.SetFloatValue(0)
})

statistics := inspector.AddSection("Statistics")
words := statistics.AddLabel("Words", "0")

// Programmatic setters never fire the callbacks.
words.SetValue("128")
```

コントロールの種類とセッターは次のとおりです。

| コントロール | セッター | コールバック |
|---------|---------|----------|
| `AddLabel` | `SetValue` | なし |
| `AddTextField` | `SetValue` | `OnTextChange` |
| `AddCheckbox` | `SetChecked` | `OnToggle` |
| `AddPopup`, `AddSegmented` | `SetOptions`, `SetSelectedIndex` | `OnSelectionChange` |
| `AddSlider`, `AddStepper` | `SetFloatValue`, `SetRange`, `SetStep` | `OnValueChange` |
| `AddColorWell` | `SetColor` | `OnColorChange` |
| `AddDatePicker` | `SetDate` | `OnDateChange` |
| `AddButton` | `SetLabel` | `OnClick` |

`SetLabel`、`SetTooltip`、`SetEnabled`、`SetHidden` はすべての種類に適用できます。セクションは折りたたみ可能にでき、セクションとコントロールはどちらもいつでも移動または削除できます。

```go
appearance.SetCollapsed(true)
inspector.MoveSection(statistics, 0)
appearance.Remove(priority)
inspector.RemoveSection(statistics)
```

## アクセサリ

`MacAccessory` はネイティブコントロールの帯です。作成時にレイアウトを選択します。`MacAccessoryLayoutLeading` はウィンドウボタンの横、`MacAccessoryLayoutTrailing` はタイトルバーの末尾側、`MacAccessoryLayoutBottom` はタイトルバーとツールバーの下で全幅に広がり、`MacAccessoryLayoutTop` は分割ビューペインの上部に配置されます。アクセサリをアタッチする前に、すべてのコントロールを追加してください。

```go
// Next to the window buttons.
folders := application.NewMacAccessory(application.MacAccessoryLayoutLeading)
folders.AddSegmented([]string{"Inbox", "Starred"}, 0).
    SetSegmentSymbols("tray", "star").
    OnSelectionChange(func(_ *application.Context, index int, label string) {})

// At the trailing edge of the titlebar.
tools := application.NewMacAccessory(application.MacAccessoryLayoutTrailing)
tools.AddSearch("Search mail").
    SetIncremental(true).
    SetWidth(200).
    OnSearch(func(_ *application.Context, query string) {})
tools.AddSymbolButton("square.and.pencil").
    SetTooltip("Compose").
    OnClick(func(*application.Context) {})

// A full-width strip beneath the titlebar and toolbar.
status := application.NewMacAccessory(application.MacAccessoryLayoutBottom).
    SetHeight(30).
    SetPreferredScrollEdgeEffectStyle(application.MacScrollEdgeEffectStyleSoft)
label := status.AddLabel("Saved").SetSymbol("checkmark.circle")
status.AddFlexibleSpace()
status.AddButton("Mark All Read").OnClick(func(*application.Context) {})

for _, accessory := range []*application.MacAccessory{folders, tools, status} {
    if err := window.AddTitlebarAccessory(accessory); err != nil {
        window.Error("titlebar accessory: %s", err)
    }
}

// Live updates and lifecycle.
label.SetText("Edited").SetSymbol("pencil.circle")
status.SetHidden(true)
tools.Remove()
```

コントロールには `AddSearch`、`AddSegmented`、`AddButton`、`AddSymbolButton`、`AddMenuButton`、`AddLabel`、`AddFlexibleSpace`、ネイティブ統合用の `AddNativeView` があります。`AddTitlebarAccessory` は `WebviewWindow` と `NativeWindow` の両方で使えます。ウィンドウの作成前の呼び出しは保留され、作成時に適用されます。`Remove` はアクセサリを取り外して別の場所へ再アタッチできるようにし、`SetHidden` はその場で折りたたみます。

### ペインアクセサリ

macOS 26 以降では、アクセサリを分割ペインの上部または下部に配置できます。Finder のサイドバーのフィルターフィールドはこの位置にあります。`Top` または `Bottom` レイアウトでアクセサリを作成し、ペインの `AddTopAccessory` または `AddBottomAccessory` でアタッチします。

```go
filter := application.NewMacAccessory(application.MacAccessoryLayoutTop)
filter.AddSearch("Filter").
    SetIncremental(true).
    OnSearch(func(_ *application.Context, query string) { filterNotes(query) })

if err := sidebarPane.AddTopAccessory(filter); err != nil {
    window.Error("sidebar filter: %s", err)
}
```

`SetPreferredScrollEdgeEffectStyle` は、アクセサリの背後をスクロールするコンテンツに適用する AppKit の `Automatic`、`Soft`、`Hard` のスタイルを選びます。明示的なスタイルには macOS 26.1 が必要です。それ以前のリリースでは、要求はウィンドウのエラーハンドラーを通じて報告され、自動スタイルが引き続き適用されます。

## 組み合わせる

このプログラムは、サイドバー、WebView、インスペクタからなる 3 ペインのウィンドウと、両方の境界線に追従するツールバーを構築します。

```go title="main.go"
package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets
var assets embed.FS

func main() {
	app := application.New(application.Options{
		Name:   "Notes",
		Assets: application.AssetOptions{Handler: application.AssetFileServerFS(assets)},
	})

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Notes",
		Width:  1100,
		Height: 700,
		URL:    "/",
		Mac: application.MacWindow{
			ContentLayout: application.MacContentLayoutEdgeToEdge,
			TitleBar: application.MacTitleBar{
				FullSizeContent:      true,
				HideToolbarSeparator: true,
				ToolbarStyle:         application.MacToolbarStyleUnified,
			},
		},
	})

	sidebar := application.NewMacSidebar()
	notes := sidebar.AddSection("Notes")
	notes.AddItem("Saturday, slowly").SetSymbol("doc.text").OnClick(func(*application.Context) {
		app.Event.Emit("note:selected", "saturday")
	})

	inspector := application.NewMacInspector()
	inspector.AddSection("Document").AddTextField("Title", "Saturday, slowly").
		OnTextChange(func(_ *application.Context, value string) {
			app.Event.Emit("note:title", value)
		})

	split := application.NewMacSplitView().SetAutosaveName("notes.main")
	split.AddSidebar(sidebar).SetMinimumThickness(200).SetCollapsible(true)
	split.AddPrimaryContent()
	split.AddInspector(inspector).SetMinimumThickness(240).SetCollapsible(true)
	window.SetSplitView(split)

	toolbar := application.NewMacToolbar().SetDisplayMode(application.MacToolbarDisplayModeIconOnly)
	toolbar.AddSidebarToggle()
	toolbar.AddSidebarTrackingSeparator()
	toolbar.AddButton("New").SetSymbol("square.and.pencil").SetBordered(true).
		OnClick(func(*application.Context) {
			notes.AddItem("Untitled").SetSymbol("doc.text")
		})
	toolbar.AddFlexibleSpace()
	toolbar.AddInspectorTrackingSeparator()
	toolbar.AddInspectorToggle()
	window.SetToolbar(toolbar)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
```

ネイティブのウィンドウ装飾とフロントエンドは、通常の Wails イベントとサービスを通じて通信します。サイドバーは `note:selected`、インスペクタは `note:title` を発行し、ページはランタイムの `Events.On` で受信します。

## ネイティブウィンドウ

@note{type="caution" title="実験的機能"}
`NativeWindow`、`NativeWindowManager`、`MacTextEditor` は v3 では実験的機能です。API は意図的に小さくしてあり、v4 向けに共通のウィンドウ API を再設計する際に変更される可能性があります。
@end

`NativeWindow` には WebView がありません。メインコンテンツは `NSScrollView` 内の `NSTextView` である `MacTextEditor` です。`WebviewWindow` と同じツールバー、分割ビュー、アクセサリの型を受け取ります。`app.NativeWindow.New` または `app.NativeWindow.NewWithOptions` で作成し、後で `Get` または `GetByID` で取得できます。

ネイティブウィンドウは、コンテンツが設定されて初めて作成されます。`AddTextEditor` でメインペインを追加した分割ビューを、`NativeWindowOptions.SplitView` または `SetSplitView` で指定してください。`app.Run` の前はレイアウトが保留されます。実行中のアプリケーションでは、`SetSplitView` は直ちにウィンドウを作成して表示し（`Hidden` が設定されている場合を除く）、作成エラーがあれば返します。レイアウトのないウィンドウは作成が保留されたままで、`Run` は `ErrNativeWindowContentRequired` を返します。テキストエディタのないレイアウトは `ErrNativeWindowEditorRequired` で拒否されます。ネイティブウィンドウが無視する `Mac.Toolbar` と `Mac.SplitView` フィールドではなく、`NativeWindowOptions.Toolbar` と `NativeWindowOptions.SplitView` を使用してください。

```go title="main.go"
package main

import (
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
	app := application.New(application.Options{
		Name:       "Native Notes",
		NativeOnly: true,
	})

	editor := application.NewMacTextEditor()
	editor.OnChange(func(*application.Context) {
		// mark the document dirty; call editor.Text() only when saving
	})

	sidebar := application.NewMacSidebar()
	sidebar.AddSection("Files").AddItem("README.txt").
		SetSymbol("doc.plaintext").
		OnClick(func(*application.Context) {
			editor.SetText("Hello from AppKit")
		})

	split := application.NewMacSplitView().SetAutosaveName("native-notes.main")
	split.AddSidebar(sidebar).SetMinimumThickness(200).SetCollapsible(true)
	split.AddTextEditor(editor).SetMinimumThickness(400)

	window := app.NativeWindow.NewWithOptions(application.NativeWindowOptions{
		Title:  "Native Notes",
		Width:  900,
		Height: 600,
		Mac: application.MacWindow{
			TitleBar: application.MacTitleBar{
				FullSizeContent: true,
				ToolbarStyle:    application.MacToolbarStyleUnified,
			},
		},
	})
	if err := window.SetSplitView(split); err != nil {
		log.Fatal(err)
	}

	toolbar := application.NewMacToolbar()
	toolbar.AddSidebarToggle()
	toolbar.AddSidebarTrackingSeparator()
	toolbar.AddFlexibleSpace()
	toolbar.AddButton("Save").SetSymbol("square.and.arrow.down").SetBordered(true).
		OnClick(func(*application.Context) {
			log.Printf("%d bytes", len(editor.Text()))
		})
	if err := window.SetToolbar(toolbar); err != nil {
		log.Fatal(err)
	}

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
```

実行中のアプリケーションからは、ウィンドウ装飾をオプションとして渡し、同じウィンドウを 1 回の呼び出しで作成できます。

```go
window := app.NativeWindow.NewWithOptions(application.NativeWindowOptions{
    Title:     "Native Notes",
    SplitView: split,
    Toolbar:   toolbar,
})
```

`MacTextEditor` は `SetText`、`Text`、`SetEditable`、`OnChange`、`Focus` を提供します。プログラムからの `SetText` 呼び出しでは `OnChange` が発火しないため、ファイルを読み込んでも未保存の変更として扱われません。`Text` は AppKit から文書全体を読み取るので、変更のたびではなく、内容が必要なときに呼び出してください。

次の 2 つのオプションで、ネイティブ専用アプリケーションを軽量に保てます。

- `application.Options` の `NativeOnly: true` は、実行時にフロントエンドの通信機構とアセットサーバーを省きます。設定した場合、`WebviewWindow` を作成しないでください。
- `wails_native` ビルドタグは、WebView、フロントエンド、アップデーターのコードをバイナリから完全に除外し、`NativeOnly` を自動的に設定します。

```sh
go build -tags wails_native .
```

`wails_native` ビルドでは単一インスタンスのサポートも除外されます。必要な場合は `wails_single_instance` タグを追加してください。

## ウィンドウタブ

macOS ではウィンドウをタブにまとめられます。タブ化モードはウィンドウの作成時に固定されるため、参加させる各ウィンドウで `Mac.TabbingMode` を `MacWindowTabbingModePreferred` または `MacWindowTabbingModeAutomatic` に設定してください。`MacWindowTabbingModeDisallowed`（および未設定のデフォルト）では、ウィンドウはタブグループに入りません。

```go
first := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Document 1",
    URL:   "/",
    Mac: application.MacWindow{
        TabbingMode: application.MacWindowTabbingModePreferred,
    },
})
```

タブ操作には作成済みのネイティブウィンドウが必要です。`app.Run` の開始後に実行されるメニューハンドラー、サービスメソッドなどから呼び出してください。`TabGroup` は、呼び出しのたびにネイティブグループを解決する `MacWindowTabGroup` ハンドルを返します。nil ハンドルも安全に使用でき、そのすべてのメソッドはゼロ値を返します。

```go
second := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Document 2",
    URL:   "/",
    Mac: application.MacWindow{
        TabbingMode: application.MacWindowTabbingModeAutomatic,
    },
})
if err := first.AddTab(second, application.MacTabOrderAbove); err != nil {
    app.Logger.Error("add tab", "error", err)
}
second.SetTabTitle("Draft")

if group := first.TabGroup(); group != nil {
    group.SelectNext()
    group.ToggleTabBar()
    for _, member := range group.Windows() {
        app.Logger.Info("tab", "name", member.Name())
    }
}

first.MoveTabToNewWindow()
first.MergeAllWindows()
```

`MacWindowTabGroup` は `Identifier`、`Count`、`Windows`、`NativeWindows`、`SelectedWindow`、`Select`、`SelectNative`、`SelectNext`、`SelectPrevious`、`IsTabBarVisible`、`ToggleTabBar`、`IsOverviewVisible`、`ToggleTabOverview` を提供します。`app.Window.TabGroups` は `WebviewWindow` を含むすべてのグループを列挙します。`AddNativeTab` は `NativeWindow` を WebView ウィンドウのグループに追加し、`SetTabTooltip` はタブのホバーテキストを設定します。macOS 以外では、タブのメソッドは `ErrMacWindowTabsUnsupported` を返します。

## バージョン要件

このページの機能はすべて macOS 専用です。古いリリースでは以下のとおり機能が段階的に縮退しますが、Go API はどの環境でも同じです。

| 機能 | 最小 macOS バージョン | それ以前のリリースでの動作 |
|---------|---------------|-------------------------------|
| ウィンドウタブ | 10.12 | 利用不可 |
| ツールバーグループ、枠付き項目、`AddMenu` | 10.15 | メニュー項目は省略され、グループは従来の表示方式を使用 |
| SF Symbols（ツールバー、サイドバー、リスト、アクセサリの項目に対する `SetSymbol`） | 11 | 画像は表示されない |
| `NSSearchToolbarItem` としての `AddSearch`、`SetNavigational`、サイドバーのトラッキングセパレーター | 11 | 検索は通常の検索フィールドに切り替わり、セパレーターは省略される |
| インスペクタとコンテンツリストの分割ロール | 11 | 同じペインが通常の分割項目に収容される |
| コンテンツリストのスタイル | 11 | 無視される |
| `SetCenteredItems` | 13 | 無視される |
| インスペクタの切り替えボタンとトラッキングセパレーター | 14 | Wails がネイティブの切り替えボタンを提供し、セパレーターは省略される |
| ツールバー項目の `SetBadgeCount`、`SetProminent`、`SetTintColor` | 26 | 保存され、利用可能になったときに適用される |
| ペインアクセサリ（`AddTopAccessory`、`AddBottomAccessory`） | 26 | `ErrMacSplitItemAccessoryUnavailable` |
| 明示的なスクロール端エフェクトのスタイル | 26.1 | ウィンドウのエラーハンドラーを通じて報告され、自動スタイルが維持される |

タイトルバーアクセサリ、分割ビュー、サイドバー、インスペクタ、テキストエディタには、Wails の最小要件以外のバージョン要件はありません。

## サンプル

各サンプルは、そのまま実行できる完全なアプリケーションです。

- [`v3/examples/mac-toolbar`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-toolbar): ツールバー、サイドバー、コンテンツリスト、インスペクタ、共有プロバイダー、ツールバーのカスタマイズ、ペインアクセサリを備えたメモエディタ。
- [`v3/examples/mac-titlebar-accessory`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-titlebar-accessory): メールボックス表示を操作する、先頭側、末尾側、下部のタイトルバーアクセサリ。
- [`v3/examples/mac-window-tabs`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-window-tabs): タブ化モード、`AddTab`、`AddNativeTab`、メニューから使うタブグループ API。
- [`v3/examples/mac-native-editor`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-native-editor): `wails_native` タグを使った、WebView のない `NativeWindow` テキストエディタ。

## さらに詳しく

このページのウィンドウ装飾は、ネイティブな macOS アプリケーションを構成する一面です。もう一面はアプリの動作で、プロキシアイコンとカスケード配置を備えたドキュメントウィンドウ、シンボルとバッジ付きのメニュー、進捗表示付きの Dock メニュー、ネイティブのアラートとパネル、取り外し可能なステータス項目、触覚フィードバックと音声、豊富な形式に対応するクリップボードと外向きのドラッグ、権限・電源・ロケールに関するシステム情報などがあります。これらの API は [macOS プラットフォーム統合](/guides/macos-platform-integration) ガイドで説明します。
