---
title: "macOS プラットフォーム統合"
description: "macOS のドキュメントウィンドウ、Dock とメニューの拡張機能、ネイティブパネル、ステータス項目、フィードバック、豊富な形式に対応するクリップボード、外向きのドラッグ、権限、電源、ライフサイクル"
slug: "guides/macos-platform-integration"
sourcePath: "guides/macos-platform-integration.md"
---

対象プラットフォーム: macOS

Wails v3 では、macOS ユーザーがネイティブアプリに期待する動作をアプリケーションに提供できます。プロキシアイコンとカスケード配置を備えたドキュメントウィンドウ、シンボルとバッジ付きのメニュー、進捗表示付きの Dock メニュー、ネイティブのアラートとパネル、取り外し可能なステータス項目、触覚フィードバックと音声、外向きのドラッグに対応した豊富な形式のクリップボード、権限・電源・ロケールに関するシステム情報、さらに Services メニュー、Handoff、AppleScript、Quick Look との連携です。すべて `application` パッケージを通じて Go から操作できます。

同じコードは Windows と Linux でもコンパイルできます。セッターは値を保存し、問い合わせはゼロ値を返し、macOS を必要とする操作は `ErrMacOnly`、`ErrDialogNotSupported`、`ErrClipboardNotSupported` などの文書化されたエラーを返します。以下の[プラットフォームに関する注記](#platform-notes)には、macOS 以外での各領域の動作を示します。

ネイティブのウィンドウ装飾（ツールバー、サイドバー、インスペクタ、アクセサリ、ウィンドウタブ）については、[macOS のネイティブウィンドウ装飾](/guides/macos-native-chrome) ガイドを参照してください。

## ドキュメントウィンドウ

ドキュメントウィンドウは、表すファイルをタイトルバーに表示し、未保存の変更を閉じるボタンの点で示し、新しいウィンドウをカスケード配置で開きます。これらはすべて `WebviewWindow` のメソッドと `MacWindow` のオプションです。

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:           "Report.md",
    URL:             "/editor",
    InitialPosition: application.WindowCascade,
    Mac: application.MacWindow{
        FrameAutosaveName: "editor.main",
        TitleBar: application.MacTitleBar{
            WindowButtonsOffset: &application.Point{X: 12, Y: 8},
        },
    },
})

window.SetRepresentedFile("/Users/me/Documents/Report.md")
window.SetSubtitle("Documents")
window.SetDocumentEdited(true)
```

- `SetRepresentedFile` はファイルのプロキシアイコンをタイトルバーに表示します。ユーザーはアイコンを別のアプリへドラッグしたり、command キーを押しながらクリックしてパスを表示したりできます。`""` を渡すと削除します。`RepresentedFile` で現在の値を読み取れます。
- `SetDocumentEdited` は閉じるボタンに未保存の変更を示す点を表示し、プロキシアイコンを薄くします。`IsDocumentEdited` で状態を読み取れます。
- `SetSubtitle` は macOS 11 以降でタイトルの下に 2 行目を表示します。
- `InitialPosition: application.WindowCascade` は、新しいドキュメントを開くときと同様に、最後にカスケード配置したウィンドウの右下にウィンドウを配置します。`CascadeFrom(other)` は既存のウィンドウに対して同じことを行い、後続のウィンドウ用にカスケード位置を更新します。
- `Mac.FrameAutosaveName` はウィンドウを初めて表示する前に保存済みの位置とサイズを復元し、移動に応じて保存し続けます。復元されたフレームは `X`、`Y`、`Width`、`Height`、`InitialPosition` より優先されます。`SetFrameAutosaveName` は作成済みウィンドウの名前を切り替えます。
- `MacTitleBar.WindowButtonsOffset` は閉じる、最小化、拡大の各ボタンを指定したポイント数だけ移動します。`SetWindowButtonsOffset` と `ResetWindowButtonsOffset` で実行時に変更できます。

これら 3 つのセッターは、ネイティブウィンドウの作成前にも呼び出せます。値は作成時に適用されます。

### 注意を引く要求

`RequestAttention` はアプリケーションがバックグラウンドにある間、Dock アイコンを跳ねさせます。情報通知の要求では 1 回だけ跳ねます。重大な要求では、ユーザーがアプリケーションをアクティブにするか、要求をキャンセルするまで跳ね続けます。

```go
request := window.RequestAttention(true)

// Once the work that needed attention is done:
request.Cancel()
```

`Flash` は、1 回だけ注意を引くためのクロスプラットフォームの方法です。

### 印刷とエクスポート

`PrintWithOptions` は明示的なページ設定で WebView を印刷します。ゼロ値では、共有の印刷設定を使った印刷パネルが表示されます。`Print` は従来の動作（横向き、余白 30 ポイント）を維持します。

`ExportPDF` はページを PDF 文書として生成し、`Snapshot` は PNG としてキャプチャします。どちらも WebKit の処理を待つため、goroutine から呼び出し、アプリケーションスレッドからは呼び出さないでください。そこで呼び出すと `ErrMacExportOnMainThread` が返されます。

```go
err := window.PrintWithOptions(application.PrintOptions{
    Orientation: application.PrintOrientationPortrait,
    Margins:     application.PrintMargins{Top: 36, Left: 36, Bottom: 36, Right: 36},
    Silent:      false,
})
if err != nil {
    log.Println("print:", err)
}

go func() {
    pdf, err := window.ExportPDF(application.PDFExportOptions{})
    if err != nil {
        log.Println("export:", err)
        return
    }
    if err := os.WriteFile("report.pdf", pdf, 0o644); err != nil {
        log.Println("write:", err)
    }

    png, err := window.Snapshot(application.SnapshotOptions{Width: 800})
    if err != nil {
        log.Println("snapshot:", err)
        return
    }
    if err := os.WriteFile("preview.png", png, 0o644); err != nil {
        log.Println("write:", err)
    }
}()
```

`PrintOptions` は `PrinterName`、`PaperName`（`"iso-a4"` のような PostScript 名）、`Scale` も受け取ります。`PDFExportOptions` と `SnapshotOptions` には、キャプチャ範囲を制限する省略可能な `Rect` と、デフォルトが `DefaultMacExportTimeout`（30 秒）の `Timeout` を指定できます。

## シート

シートは、保存パネルのように親ウィンドウの上部に取り付けられる別のウィンドウです。どの `WebviewWindow` も `PresentSheet` で別のウィンドウのシートとして表示でき、`EndSheet` と応答コードで終了できます。応答コードはシートの `OnSheetEnd` コールバックに渡されます。アタッチ前に画面へ一瞬表示されないよう、シートのウィンドウは `Hidden` を設定して作成してください。終了時には AppKit が再び画面から取り除くため、同じウィンドウを繰り返し表示できます。

```go
sheet := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "Rename",
    URL:    "/rename",
    Width:  420,
    Height: 180,
    Hidden: true,
})
sheet.OnSheetEnd(func(code int) {
    if code == application.MacSheetResponseOK {
        log.Println("renamed")
    }
})

// From a menu item or a bound method:
if err := window.PresentSheet(sheet); err != nil {
    log.Println(err)
}

// From the sheet's own page, through a bound method:
sheet.EndSheet(application.MacSheetResponseOK)
```

`PresentCriticalSheet` は、既にアタッチされている通常のシートの後ろで待機させず、その前面にシートを表示します。`PresentNativeSheet` は `NativeWindow` に対して同じことを行います。`IsSheet`、`SheetParent`、`AttachedSheet`、`AttachedNativeSheet`、`HasAttachedSheet` は現在の状態を示します。シートを終了する代わりにシートウィンドウを閉じると、`MacSheetResponseStop` が通知されます。

## ポップオーバー

`MacPopover` は `NSPopover` です。ウィンドウ内の矩形、ツールバー項目、またはメニューバーのステータス項目を基点とする一時的なパネルです。内容はネイティブの `MacAccessory` コントロールの帯で、[macOS のネイティブウィンドウ装飾](/guides/macos-native-chrome) ガイドでタイトルバーアクセサリに使う型と同じです。初めて表示する前にすべてのコントロールを追加してください。

```go
strip := application.NewMacAccessory(application.MacAccessoryLayoutBottom)
strip.AddLabel("Sort by")
strip.AddSegmented([]string{"Date", "Title"}, 0)
strip.AddFlexibleSpace()
strip.AddButton("Apply").OnClick(func(*application.Context) {
    // apply the sort
})

popover := application.NewMacPopover(application.MacPopoverOptions{
    Width:    320,
    Behavior: application.MacPopoverBehaviorTransient,
    Content:  strip,
})
popover.OnClose(func() {
    log.Println("popover closed")
})

// Anchored to a rectangle in the page, in window content coordinates:
err := popover.ShowRelativeTo(application.Rect{X: 20, Y: 60, Width: 120, Height: 28}, window, application.MacRectEdgeMaxY)
if err != nil {
    log.Println(err)
}
```

`MacToolbarItem.ShowPopover` と `SystemTray.ShowPopover` は、同じポップオーバーをツールバー項目またはステータス項目に固定します。`MacPopoverBehaviorTransient` はポップオーバーの外側をクリックすると閉じ、`Semitransient` は表示元のウィンドウ内をクリックした場合にのみ閉じます。デフォルトでは `Close` まで開いたままです。`MacRectEdge` はポップオーバーを表示する側を選びます。`SetContentSize` と `SetBehavior` は表示中のポップオーバーを調整し、`Destroy` はネイティブのポップオーバーを解放して、コンテンツの帯を別の場所で使えるようにします。

## 状態の復元

macOS は、クラッシュ、強制終了、再起動の後にアプリケーションのウィンドウを再び開きます。また、システム設定の「アプリケーションを終了するときにウィンドウを閉じる」がオフの場合、通常の終了後にも再び開きます。ウィンドウに `Mac.RestorationID` を設定し、`SetRestorationData` で再作成に必要な情報を保存し、次回起動時にウィンドウを再構築する `app.Window.OnRestore` を登録してください。

```go
app.Window.OnRestore(func(id string, state application.RestorationState) application.Window {
    if id != "editor" {
        return nil
    }
    path := state.Get("path")
    restored := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title: filepath.Base(path),
        URL:   "/editor?path=" + path,
        Mac:   application.MacWindow{RestorationID: id},
    })
    restored.SetRestorationData(state.Data)
    return restored
})

window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Report.md",
    URL:   "/editor?path=/Users/me/Documents/Report.md",
    Mac:   application.MacWindow{RestorationID: "editor"},
})
window.SetRestorationData(map[string]string{"path": "/Users/me/Documents/Report.md"})
```

アプリケーション終了時に表示されているウィンドウだけが保存されます。`RestorationState.Data` は文字列のマップです。識別子、パス、位置などに絞ってください。`InteractionState` は WebView の戻る・進む履歴とスクロール位置を不透明なデータとして返します（macOS 12 以降）。`RestoreInteractionState` はこれを再作成したウィンドウに適用します。通常は base64 でエンコードして復元データに保存します。`SetRestorationID` と `RestorationID` は、作成済みウィンドウの識別子を変更、取得します。

## 表示オプション

`MacPresentationOptions` は `NSApplication.presentationOptions` に対応するビットマスクです。アプリケーションがアクティブな間、Dock やメニューバーを隠したり、プロセス切り替え、強制終了、ログアウト、「隠す」コマンドを無効にしたりします。起動時に適用するにはアプリケーションオプションの `Mac.PresentationOptions` を設定し、実行時に変更するには `SetPresentationOptions` を使います。無効な組み合わせは AppKit に渡される前に拒否され、`ErrMacPresentationOptionsInvalid` をラップしたエラーが返されます。

```go
kiosk := application.MacPresentationHideDock |
    application.MacPresentationHideMenuBar |
    application.MacPresentationDisableProcessSwitching |
    application.MacPresentationDisableForceQuit

if err := app.SetPresentationOptions(kiosk); err != nil {
    log.Println(err)
}
log.Println("presentation:", app.PresentationOptions())

// Restore the standard Dock and menu bar:
_ = app.SetPresentationOptions(application.MacPresentationDefault)
```

メニューバーを隠す（`HideMenuBar` または `AutoHideMenuBar`）には、Dock オプションのいずれかが必要です。`AutoHideToolbar` には `FullScreen` と `AutoHideMenuBar` の両方が必要です。`Validate` は値が最初に違反する規則を報告し、`Has` は個々のフラグを調べます。

## メニューと Dock

メニュー項目では、SF Symbols、バッジ、セクションヘッダー、カラーパレット、混合チェック状態、代替項目、インデントを使えます。これらはすべて `MenuItem` と `Menu` のメソッドなので、アプリケーションメニュー、コンテキストメニュー、トレイメニュー、Dock メニューで機能します。

```go
menu := app.NewMenu()
menu.AddRole(application.AppMenu)

fileMenu := menu.AddSubmenu("File")
fileMenu.Add("Open...").SetSymbol("folder").SetAccelerator("CmdOrCtrl+o").
    OnClick(func(*application.Context) {
        // open the file, then note it in the recent list:
        app.Menu.AddRecentDocument("/Users/me/Documents/Report.md")
    })
fileMenu.AddRole(application.OpenRecent)

view := menu.AddSubmenu("View")
view.AddSectionHeader("Mailboxes")
inbox := view.Add("Inbox").SetSymbol("tray.full").SetBadge(3)
inbox.OnClick(func(*application.Context) {
    inbox.ClearBadge()
})
view.Add("Updates").SetBadgeText("New")

view.AddSeparator()
wrap := view.AddCheckbox("Wrap lines", false).SetMixed()
view.Add("Reset").SetIndentationLevel(1).OnClick(func(*application.Context) {
    wrap.SetMixed()
})

view.AddSeparator()
view.Add("Close Tab").SetAccelerator("CmdOrCtrl+w").OnClick(func(*application.Context) {})
view.Add("Close All Tabs").SetAccelerator("CmdOrCtrl+OptionOrAlt+w").SetAlternate(true).
    OnClick(func(*application.Context) {})

colours := []application.RGBA{
    application.NewRGB(255, 59, 48),
    application.NewRGB(255, 149, 0),
    application.NewRGB(52, 199, 89),
}
view.AddPalette([]string{"tag.fill"}, colours, 0, func(_ *application.Context, index int) {
    log.Println("tag colour", index)
}).SetLabel("Tag colour")

app.Menu.Set(menu)
```

- `SetSymbol` はタイトルの横に SF Symbol を表示します（macOS 11 以降）。`SetBitmap` で設定した画像は置き換えられます。
- `SetBadge` は件数を、`SetBadgeText` は短い文字列をタイトルの後に表示します（macOS 14 以降）。`ClearBadge` で削除し、`BadgeCount` と `BadgeText` で値を読み取れます。
- `AddSectionHeader` は操作できないヘッダーを追加します（macOS 14 以降）。それ以前のリリースでは、同じタイトルを持つ無効な項目になります。
- `AddPalette` は `NSMenu` のパレットメニューを使った色見本の行を追加します（macOS 14 以降）。色見本ごとに 1 つのシンボルを渡すか、塗りつぶされた円を使う場合は空のスライスを渡します。ラベルがなければパレットは親メニュー内に表示され、`SetLabel` を使うとタイトル付きのサブメニューになります。`PaletteSelected` は選択されたインデックスを返します。
- `SetMixed` はチェックボックスを、横線で描かれる混合状態にします。クリックすると AppKit と同様に完全なオンになります。
- `SetAlternate(true)` は、異なる修飾キーが押されている間、直前の項目の代わりにその項目を表示します。2 つの項目は同じキーを共有し、修飾キーが異なる必要があります。
- `SetIndentationLevel` はタイトルを最大 15 段階までインデントします。

### 最近使った項目を開く

`fileMenu.AddRole(application.OpenRecent)` は標準の「最近使った項目を開く」サブメニューを追加します。macOS では開くたびに `NSDocumentController` が内容を設定し、「メニューを消去」項目も含まれます。`app.Menu.AddRecentDocument` でファイルを追加し、`RecentDocuments` で一覧を取得し、`ClearRecentDocuments` で一覧を空にします。一覧は再起動後も保持されます。

最近使ったファイルを選ぶと、Finder からファイルを開いた場合と同じイベントが届きます。

```go
app.Event.OnApplicationEvent(events.Common.ApplicationOpenedWithFile, func(event *application.ApplicationEvent) {
    path := event.Context().Filename()
    log.Println("open", path)
})
```

### Dock メニュー

`app.Menu.SetDockMenu` は、Dock アイコンを右クリックしたときに表示する静的メニューを設定します。`OnDockMenu` は表示直前に毎回メニューを構築します。項目が変化する状態を反映する場合に適しています。ビルダーは静的メニューより優先されます。

```go
dockMenu := application.NewMenu()
dockMenu.Add("New Document").SetSymbol("doc.badge.plus").OnClick(func(*application.Context) {
    // create a document
})
app.Menu.SetDockMenu(dockMenu)

app.Menu.OnDockMenu(func() *application.Menu {
    dynamic := application.NewMenu()
    dynamic.AddSectionHeader("Recent")
    for _, path := range app.Menu.RecentDocuments() {
        path := path
        dynamic.Add(filepath.Base(path)).OnClick(func(*application.Context) {
            app.Menu.OpenRecentDocument(path)
        })
    }
    return dynamic
})
```

### Dock の進捗表示

Dock サービスは、既存のバッジ機能に加えて Dock アイコン上に進捗バーを描画します。`dock.New()` をサービスとして登録し、0 から 1 までの割合を指定して `SetProgress` を呼び出します。

```go
import "github.com/wailsapp/wails/v3/pkg/services/dock"

dockService := dock.New()

app := application.New(application.Options{
    Name: "Exporter",
    Services: []application.Service{
        application.NewService(dockService),
    },
})

app.Event.On("export:progress", func(event *application.CustomEvent) {
    if fraction, ok := event.Data.(float64); ok {
        _ = dockService.SetProgress(fraction)
    }
})
app.Event.On("export:done", func(*application.CustomEvent) {
    _ = dockService.ClearProgress()
})
```

`GetProgress` は現在の割合を返します。バーが表示されていない場合は `nil` を返します。

## ダイアログ

メッセージ、開く、保存の各ダイアログは macOS オプションに対応し、ダイアログマネージャーにはテキスト入力プロンプトとシステムのカラー・フォントパネルが追加されます。

### アラート

`SetSuppression` は「このメッセージを再度表示しない」チェックボックスを追加し、`SetHelp` はヘルプボタンを表示します。ボタンのコールバックから `Suppressed` でチェック状態を読み取るか、先に受け取るために `OnSuppression` を登録します。

```go
dialog := app.Dialog.Warning().
    SetTitle("Delete 3 items?").
    SetMessage("The items will be moved to the Bin.").
    SetSuppression("Do not warn me again").
    SetHelp(func() {
        log.Println("help requested")
    }).
    AttachToWindow(window)

dialog.AddButton("Delete").SetAsDefault().OnClick(func() {
    if dialog.Suppressed() {
        // remember not to ask again
    }
})
dialog.AddButton("Cancel").SetAsCancel()
dialog.Show()
```

### テキスト入力プロンプト

`Prompt` はテキストフィールド付きのアラートを表示し、閉じられるまで処理をブロックします。そのため、goroutine またはバインドされたメソッドから呼び出してください。`Secure` はフィールドをパスワード入力欄にし、`Window` はアラートをシートとして表示します。

```go
go func() {
    value, ok, err := app.Dialog.Prompt(application.PromptOptions{
        Title:        "Name this document",
        Message:      "The name is used for the exported file.",
        Placeholder:  "Untitled",
        DefaultValue: "Quarterly report",
        OKLabel:      "Rename",
        Window:       window,
    })
    if err != nil || !ok {
        return
    }
    log.Println("renamed to", value)
}()
```

### ファイルパネル

`AddContentType` は、開くダイアログと保存ダイアログの両方で統一型識別子によるフィルターを設定します。`AddFilter` と併用できるため、`"public.image"` はシステムが認識するすべての画像形式に一致し、別のフィルターで拡張子による PDF の検出もできます。

```go
paths, err := app.Dialog.OpenFile().
    SetTitle("Choose images").
    AddFilter("PDF", "*.pdf").
    AddContentType("public.image").
    AttachToWindow(window).
    PromptForMultipleSelection()
if err == nil {
    log.Println(paths)
}
```

保存パネルでは、形式のポップアップ、名前フィールドのカスタムラベル、Finder タグが使えます。`SetFormats` はユーザーがポップアップの選択を変えるたびに、許可する型と名前フィールドの拡張子を切り替えます。`SelectedFormat` は最終的な選択を報告します。

```go
formats := []application.DialogFormat{
    {Label: "PNG image", Extension: "png", UTI: "public.png"},
    {Label: "PDF document", Extension: "pdf", UTI: "com.adobe.pdf"},
}
save := app.Dialog.SaveFile().
    SetFilename("Quarterly report.png").
    SetNameFieldLabel("Export As:").
    SetTags([]string{"Reports", "Draft"}).
    AttachToWindow(window)
save.SetFormats(formats, 0, nil)

path, err := save.PromptForSingleSelection()
if err == nil && path != "" {
    log.Println("export", path, "as", formats[save.SelectedFormat()].Label)
}
```

### カラー・フォントパネル

`PickColor` と `PickFont` は共有のシステムパネルを開き、パネルが閉じるまで処理をブロックします。`OnChange` はパネルを開いている間の選択を毎回通知するため、ページで選択結果をリアルタイムにプレビューできます。各種類のパネルは同時に 1 つしか開けません。2 回目の呼び出しは `ErrDialogInProgress` を返します。

```go
go func() {
    colour, changed, err := app.Dialog.PickColor(application.ColorPickerOptions{
        Initial:    application.NewRGB(52, 120, 246),
        ShowsAlpha: true,
        Title:      "Accent colour",
        OnChange: func(colour application.RGBA) {
            app.Event.Emit("accent:preview", colour)
        },
    })
    if err == nil && changed {
        log.Println("accent", colour)
    }

    font, changed, err := app.Dialog.PickFont(application.FontPickerOptions{
        Family: "Helvetica Neue",
        Size:   18,
    })
    if err == nil && changed {
        log.Println(font.Family, font.Face, font.PostScriptName, font.Size)
    }
}()
```

## ステータス項目とフィードバック

### ステータス項目

macOS のシステムトレイ項目は `NSStatusItem` です。SF Symbol で描画でき、ツールチップを付けられ、組み込み項目と同じようにユーザーが取り外せます。

```go
tray := app.SystemTray.New()
tray.SetSymbol("waveform.circle").SetSymbolConfiguration(0, application.MacSymbolWeightMedium)
tray.SetTooltip("Recorder. Command-drag to remove.")
tray.SetRemovable(true, "com.example.recorder.tray")
tray.OnVisibilityChange(func(visible bool) {
    log.Println("status item visible:", visible)
})

trayMenu := app.Menu.New()
trayMenu.Add("Show").OnClick(func(*application.Context) {
    window.Show().Focus()
})
tray.SetMenu(trayMenu)
tray.Run()
```

- `SetSymbol` はシンボルをテンプレート画像として描画し、メニューバーの外観に追従させます（macOS 11 以降）。`SetSymbolConfiguration` はポイントサイズと太さを設定します。
- `SetTooltip` はホバーテキストを設定します。`Tooltip` で読み取れます。
- `SetRemovable(true, name)` を使うと、ユーザーは command キーを押しながら項目をメニューバーの外へドラッグできます。macOS が起動をまたいで取り外しを記憶できるよう、安定した自動保存名を指定してください。`Show` または `SetVisible(true)` で再表示します。
- `IsVisible` は `NSStatusItem.visible` を読み取るため、ユーザーが項目を取り外した後は false になります。`OnVisibilityChange` はすべての変化を報告します。

### 触覚フィードバック

`app.Haptics.Perform` は、アプリケーションがアクティブな間、Force Touch トラックパッドまたは Magic Trackpad でパターンを再生します。

```go
app.Haptics.Perform(application.HapticAlignment)
```

種類は `HapticGeneric`、`HapticAlignment`（項目が所定の位置に収まるとき）、`HapticLevelChange`（段階やクリック位置が変わるとき）です。`IsSupported` は、そのプラットフォームでフィードバックを再生できるかどうかを報告します。

### サウンド

`app.Sound` は警告音、名前を指定したシステムサウンド、または音声ファイルを再生します。

```go
app.Sound.Beep()

if err := app.Sound.Play("Glass"); err != nil {
    log.Println(err)
}

for _, name := range app.Sound.SystemSounds() {
    log.Println(name)
}
```

`Play` は `SystemSounds` に含まれる名前、または Core Audio がデコードできるファイルへの絶対パスを受け取ります。`PlayData` はメモリ上の完全な音声ファイルを再生します。

### 音声

`app.Speech.Speak` はシステムの音声で読み上げるテキストをキューに追加し、`Utterance` を返します。発話は順番に再生されます。`Stop` は 1 件を取り除き、`StopAll` はキューを空にします。`Voices` はインストール済みの音声を識別子と言語とともに列挙します。

```go
utterance, err := app.Speech.Speak("Export finished", application.SpeechOptions{
    Voice: "com.apple.voice.compact.en-GB.Daniel",
    Rate:  0.5,
})
if err != nil {
    log.Println(err)
    return
}
utterance.OnFinished(func() {
    log.Println("done, stopped:", utterance.WasStopped())
})
```

`Recognize` は `SFSpeechRecognizer` を使ってデフォルトのマイクから音声を文字起こしします。最初の呼び出しではマイクと音声認識の権限を求め、ユーザーが応答するまで処理をブロックするため、goroutine から呼び出してください。途中の文字起こし結果は `OnPartial` に届きます。`Stop` は収録を終了し、最終的なテキストを返します。

```go
go func() {
    session, err := app.Speech.Recognize(application.RecognitionOptions{
        Locale: "en-US",
        OnPartial: func(text string) {
            app.Event.Emit("dictation:partial", text)
        },
    })
    if err != nil {
        if errors.Is(err, application.ErrSpeechRecognitionDenied) {
            _ = app.Permissions.OpenSystemSettings(application.PermissionKindMicrophone)
        }
        log.Println(err)
        return
    }
    time.Sleep(5 * time.Second)
    text, err := session.Stop()
    if err != nil {
        log.Println(err)
        return
    }
    app.Event.Emit("dictation:final", text)
}()
```

音声認識には、`Info.plist` に `NSSpeechRecognitionUsageDescription` と `NSMicrophoneUsageDescription` を宣言したバンドル形式のアプリケーションが必要です。これらがない場合、macOS はアクセスを拒否し、`Recognize` は `ErrSpeechRecognitionUsageDescription` を返します。

## クリップボードとドラッグ

### 豊富な形式に対応するクリップボード

`app.Clipboard` はプレーンテキストに加え、画像、ファイル参照、HTML、RTF、任意の統一型識別子に対応する生データを読み書きします。`Types` はペーストボード上の型を列挙し、`OnChange` はどのアプリケーションによる変更も報告します。

```go
if err := app.Clipboard.SetHTML("<p>Rich <b>HTML</b></p>", "Rich HTML"); err != nil {
    log.Println(err)
}
_ = app.Clipboard.SetFiles([]string{"/Users/me/Documents/Report.md"})
_ = app.Clipboard.SetData("com.example.record", []byte(`{"id":42}`))

stop := app.Clipboard.OnChange(func() {
    log.Println("clipboard changed", app.Clipboard.ChangeCount(), app.Clipboard.Types())
    if files, err := app.Clipboard.Files(); err == nil && len(files) > 0 {
        log.Println("files:", files)
    }
})
defer stop()
```

`SetImage` と `Image` は PNG バイト列を扱います。他のアプリが TIFF としてコピーした画像は自動的に変換されます。クリップボードの変更にはシステム通知がないため、`OnChange` はリスナーが少なくとも 1 つある間、500 ms ごとに変更回数をポーリングします。

### 外向きのドラッグ

`StartDrag` は、ユーザーが Finder で項目を持ち上げたかのように、ウィンドウからシステムのドラッグを開始します。既存のファイル、ドロップ先が受け入れたときだけ内容が生成されるファイルプロミス、またはプレーンテキストを提供できます。マウス操作中に開始してください。Go メソッドをバインドし、ドラッグ対象要素のページ側 `mousedown` または `pointerdown` ハンドラーから呼び出します。WebKit が独自のドラッグを始めないよう、HTML の `draggable` 属性は `false` に設定します。

```go
// DragService is bound to the page and called from a mousedown handler.
type DragService struct {
    window *application.WebviewWindow
}

func (s *DragService) DragExport() error {
    return s.window.StartDrag(application.DragItems{
        Promises: []application.DragPromise{{
            Filename: "export.csv",
            Data: func() ([]byte, error) {
                return []byte("id,name\n1,Wails\n"), nil
            },
        }},
        Operations: application.DragOperationCopy,
    })
}
```

```go
window.OnDragEnd(func(operation application.DragOperation) {
    log.Println("drag ended:", operation)
})
```

操作中以外に `StartDrag` を呼び出すと `ErrDragOutNoGesture` が返されます。`DragItems.Image` と `ImageOffset` はカーソルの下に表示する画像を設定します。

### 他のアプリケーションからのドロップ

ファイルのドロップには引き続き `WindowFilesDropped` イベントを使います。他のアプリからドラッグされたテキスト、URL、画像を受け取るには、`DropTypes` に型を列挙し、`OnDrop` を登録します。これらのドロップはページ自身の HTML5 ドロップハンドラーではなく、Go に届けられます。

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Inbox",
    URL:   "/",
    DropTypes: []application.DropType{
        application.DropFiles,
        application.DropText,
        application.DropURLs,
        application.DropImages,
    },
})

window.OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {
    log.Println("files:", event.Context().DroppedFiles())
})

window.OnDrop(func(_ *application.Context, data application.DropData) {
    log.Println("text:", data.Text, "urls:", data.URLs, "images:", len(data.Images), "at", data.X, data.Y)
})
```

## システム

### 権限

`app.Permissions` は、システムのプライバシー権限（カメラ、マイク、画面収録、アクセシビリティ、位置情報、通知、入力監視、フルディスクアクセス）の状態を報告し、権限を要求します。`Status` は確認を求めません。`Request` は未決定の種類について確認を求め、ユーザーの応答まで処理をブロックするため、goroutine から呼び出してください。`OpenSystemSettings` は対応する「プライバシーとセキュリティ」の設定画面を開きます。

```go
go func() {
    status := app.Permissions.Status(application.PermissionKindCamera)
    if status == application.PermissionStatusNotDetermined {
        status, _ = app.Permissions.Request(application.PermissionKindCamera)
    }
    if status == application.PermissionStatusDenied {
        _ = app.Permissions.OpenSystemSettings(application.PermissionKindCamera)
    }
}()
```

フルディスクアクセスは要求できず、`ErrPermissionNotRequestable` を返します。ユーザーを設定画面に案内してください。応答が得られない要求は `ErrPermissionRequestTimeout` を返します。macOS では通常、その種類に対応する使用目的の説明キーが `Info.plist` にないことを意味します。

`Permissions` ウィンドウオプションは、macOS でも適用されるようになりました。ページからの `getUserMedia` 要求の処理方法を決めます。`PermissionAllow` は WebView 独自の確認を省略し、`PermissionDeny` は確認せず拒否し、`PermissionDefault` は確認を表示します。カメラまたはマイクを初めて使用するときには、システムレベルの TCC 確認が引き続き表示されます。

```go
app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Meeting",
    URL:   "/",
    Permissions: map[application.PermissionType]application.Permission{
        application.PermissionMicrophone: application.PermissionAllow,
        application.PermissionCamera:     application.PermissionDefault,
    },
})
```

### 電源

`app.Power.PreventSleep` は、返された解放関数が呼び出されるまでシステムをスリープさせず、`Display` を指定すると画面もスリープさせません。保持は個別に数えられるため、アプリの複数の部分から同時に保持できます。理由はアクティビティモニタに表示されます。

```go
release, err := app.Power.PreventSleep("Exporting video", application.PreventSleepOptions{Display: true})
if err != nil {
    log.Println(err)
}
defer release()

state := app.Power.State()
if state.LowPowerMode || state.ThermalState >= application.ThermalStateSerious {
    // trim background work
}
log.Println("battery", state.BatteryLevel, "charging", state.Charging, "on battery", state.OnBattery)
```

変化は `events.Mac.ApplicationDidChangePowerState`（低電力モードの切り替え）と `events.Mac.ApplicationDidChangeThermalState` として届きます。

### ライフサイクル

アプリが `NSSupportsSuddenTermination` で許可すると、macOS はログアウトまたはシャットダウン時にアイドル状態のアプリを直ちに終了できます。また、`NSSupportsAutomaticTermination` で許可すると、ウィンドウのないアイドル状態のアプリを終了できます。`app.Lifecycle.HoldTermination` はファイルの保存などの重要な処理中、両方を一時停止します。

```go
release := app.Lifecycle.HoldTermination("Saving document")
defer release()
// write the file
```

`SetSuddenTerminationEnabled` は実行時に突然の終了を切り替えます。`SuddenTerminationEnabled` は現在の状態を報告し、初期状態は `Info.plist` のキーで決まります。

### 環境

`app.Env` に、ユーザー設定についての 3 つの問い合わせが追加されます。

```go
a11y := app.Env.Accessibility()
if a11y.ReduceMotion || a11y.ReduceTransparency {
    app.Event.Emit("theme:calm", true)
}

layout := app.Env.KeyboardLayout()
log.Println(layout.ID, layout.Name, layout.Languages)

locale := app.Env.Locale()
log.Println(locale.Identifier, locale.Language, locale.Region, locale.Preferred)
```

- `Accessibility` は「視差効果を減らす」「透明度を下げる」「コントラストを上げる」「カラー以外で区別」「色を反転」、VoiceOver、スイッチコントロールの設定を反映します。
- `KeyboardLayout` は、アクティブな入力ソースの識別子、ローカライズされた名前、言語を返します。
- `Locale` は、AppKit がアプリ向けに選択したロケールと、ユーザーの優先順に並んだ `Preferred` の完全な一覧を返します。`Identifier` が反映するのはバンドルが `CFBundleLocalizations` で宣言した言語だけです。言語を自分で選ぶには `Preferred` を使ってください。

### イベント

次のアプリケーションイベントが新たに追加されました。いずれも `app.Event.OnApplicationEvent` を通じて通知されます。最新の値は対応するマネージャーに問い合わせてください。

| イベント | 発生するタイミング | 値の取得元 |
|-------|------------|-----------|
| `events.Mac.ApplicationDidChangePowerState` | 低電力モードが切り替わったとき | `app.Power.State()` |
| `events.Mac.ApplicationDidChangeThermalState` | 熱負荷が変化したとき | `app.Power.State()` |
| `events.Common.AccessibilitySettingsChanged` | アクセシビリティの表示設定が変化したとき | `app.Env.Accessibility()` |
| `events.Mac.ApplicationDidChangeAccessibilitySettings` | 同じ変化の macOS 形式 | `app.Env.Accessibility()` |
| `events.Mac.ApplicationDidChangeKeyboardLayout` | 入力ソースが変化したとき | `app.Env.KeyboardLayout()` |
| `events.Mac.ApplicationDidChangeLocale` | ロケールが変化したとき | `app.Env.Locale()` |

```go
app.Event.OnApplicationEvent(events.Mac.ApplicationDidChangeThermalState, func(*application.ApplicationEvent) {
    app.Event.Emit("system:power", app.Power.State())
})
app.Event.OnApplicationEvent(events.Common.AccessibilitySettingsChanged, func(*application.ApplicationEvent) {
    app.Event.Emit("system:accessibility", app.Env.Accessibility())
})
```

## 統合

### Services メニュー

`app.ServicesProvider.Register` は、選択されたテキストやファイルに対してすべての macOS アプリケーションが表示する Services サブメニューに項目を追加します。ハンドラーはペーストボードを `ServiceRequest` として受け取り、書き戻すための `ServiceResponse` を返します。空の応答では選択内容は変更されません。AppKit はメインスレッドでハンドラーの完了を待つため、処理を短時間で終えてください。

```go
err := app.ServicesProvider.Register(application.ServiceDefinition{
    Name:          "summarise",
    MenuTitle:     "Summarise with Notes",
    SendTypes:     []string{"public.utf8-plain-text"},
    ReturnTypes:   []string{"public.utf8-plain-text"},
    KeyEquivalent: "S",
    Handler: func(_ *application.Context, request application.ServiceRequest) (application.ServiceResponse, error) {
        return application.ServiceResponse{Text: "Summary: " + request.Text}, nil
    },
})
if err != nil {
    log.Println(err)
}

// Paste this inside the top-level <dict> of Info.plist.
log.Println(app.ServicesProvider.InfoPlistXML())
```

`Name` は AppKit が送信するメッセージで、単純な識別子でなければなりません。`SendTypes` と `ReturnTypes` はペーストボードの型で、サービスには少なくとも一方が必要です。登録するだけではサービスは表示されません。バンドルの `Info.plist` で `NSServices` の下に宣言する必要があります。`InfoPlistXML` は貼り付け可能なそのブロックを返し、`InfoPlistEntries` は plist シリアライザー向けに同じデータをマップとして返します。これらの項目の `NSPortName` はアプリケーションの `Name` で、`CFBundleName` と一致する必要があります。

Wails CLI でビルドするプロジェクトでは、同じサービスを `build/config.yml` に一度だけ宣言し、パッケージング時に `NSServices` ブロックを生成できます。

```yaml
services:
  - name: SummariseText
    menuTitle: Summarise with My Product
    sendTypes:
      - public.utf8-plain-text
    returnTypes:
      - public.utf8-plain-text
    keyEquivalent: S
```

各項目は、Go で登録した同じ `Name` の `ServiceDefinition` と一致する必要があります。新しいビルドをインストールした後に `pbs -update` を実行すると、ログアウトせずに Services メニューへ変更が反映されます。

### Handoff とユーザーアクティビティ

`app.Activity.Publish` は `NSUserActivity` を現在のアクティビティとして設定し、ユーザーが別のデバイスで続きを行ったり、Spotlight で見つけたり、Siri から提案を受けたりできるようにします。返された `PublishedActivity` は状態の変化に応じて更新でき、ドキュメントを閉じたときに無効化できます。新しいアクティビティを公開すると、以前のものは置き換えられます。

```go
activity, err := app.Activity.Publish(application.UserActivity{
    Type:               "com.example.notes.editing",
    Title:              "Editing Quarterly report",
    UserInfo:           map[string]any{"note": "quarterly-report"},
    WebpageURL:         "https://example.com/notes/quarterly-report",
    EligibleForHandoff: true,
    EligibleForSearch:  true,
    Keywords:           []string{"report", "quarterly"},
})
if err != nil {
    log.Println(err)
    return
}
if err := activity.Update(map[string]any{"note": "quarterly-report", "cursor": 120}); err != nil {
    log.Println(err)
}
// When the document closes:
activity.Invalidate()
```

受信したアクティビティは `OnContinue` を通じて届きます。ユニバーサルリンクの型は `UserActivityTypeBrowsingWeb` で、ページは `WebpageURL` に入ります。`events.Common.ApplicationLaunchedWithUrl` としても通知されるため、アプリは URL の処理経路を 1 つにまとめられます。`OnWillContinue`、`OnFailed`、`OnUpdated` は残りのデリゲート処理に対応します。

```go
app.Activity.OnContinue(func(_ *application.Context, incoming application.UserActivity) bool {
    if incoming.Type == application.UserActivityTypeBrowsingWeb {
        log.Println("universal link", incoming.WebpageURL)
        return true
    }
    note, _ := incoming.UserInfo["note"].(string)
    log.Println("continue editing", note)
    return true
})
```

すべてのアクティビティ型を `Info.plist` の `NSUserActivityTypes` に列挙する必要があります。ユニバーサルリンクには、`applinks:example.com` 項目を含む `com.apple.developer.associated-domains` エンタイトルメントと、そのドメイン上の対応する `apple-app-site-association` ファイルも必要です。

### Apple Events

`app.AppleEvents.Handle` はイベントのクラスと ID に対するハンドラーを登録し、AppleScript、Shortcuts、他のアプリケーションからアプリを操作できるようにします。コードは 4 文字の文字列です。直接パラメーターは Go の値（`string`、ファイルパスの `[]string`、`int64`、`float64`、`bool`、`[]any`、`AppleEventRawData`）にデコードされ、応答の `Result` も同じ種類を受け取ります。イベントが保留されている間、ハンドラーは専用の goroutine で実行されます。

```go
err := app.AppleEvents.Handle("WAIL", "note", func(_ *application.Context, event application.AppleEvent) (application.AppleEventReply, error) {
    text, _ := event.DirectObject.(string)
    return application.AppleEventReply{Result: "noted: " + text}, nil
})
if err != nil {
    log.Println(err)
}

go func() {
    result, err := app.AppleEvents.Send("com.apple.finder", "misc", "actv", nil)
    log.Println(result, err)
}()

sdef := app.AppleEvents.ScriptingDefinition()
if err := os.WriteFile("Notes.sdef", []byte(sdef), 0o644); err != nil {
    log.Println(err)
}
```

スクリプトは生のイベント構文を使って、直ちにハンドラーを呼び出せます。

```applescript
tell application id "com.example.notes" to «event WAILnote» "hello"
```

`ScriptingDefinition` は、各ハンドラーにコマンド名を与える最小限の `.sdef` を生成します。これを `Contents/Resources` に含め、`Info.plist` の `NSAppleScriptEnabled` と `OSAScriptingDefinition` で参照してください。するとスクリプトエディタの「ファイル」>「用語説明を開く」に表示されます。Wails はカスタム URL スキームの Get URL イベントを既に処理しています。`"GURL"`/`"GURL"` のハンドラーはその処理に連結されますが、`"aevt"`/`"odoc"` のハンドラーは組み込みの Open Documents 通知を置き換えます。`Send` はバンドル識別子で実行中のアプリケーションを指定します。バンドル形式のアプリでは `NSAppleEventsUsageDescription` が必要で、応答が届くまで呼び出し元の goroutine をブロックします。

### Quick Look

`app.QuickLook.Preview` は、1 つ以上のファイルを共有の Quick Look パネルで開きます。複数のパスを指定すると、パネルにファイル間を移動する矢印が表示されます。`Thumbnail` はシステムのサムネイルプロバイダーを使ってファイルを描画し、PNG を返します。ドキュメント、画像、PDF、動画に対応します。

```go
if err := app.QuickLook.Preview([]string{"/Users/me/Documents/Report.pdf"}); err != nil {
    log.Println(err)
}

go func() {
    png, err := app.QuickLook.Thumbnail("/Users/me/Documents/Report.pdf", application.ThumbnailOptions{
        Width: 256,
        Scale: 2,
    })
    if err != nil {
        log.Println(err)
        return
    }
    if err := os.WriteFile("thumbnail.png", png, 0o644); err != nil {
        log.Println(err)
    }
}()
```

パスは絶対パスで、対象が存在する必要があります。`ClosePreview` と `IsPreviewOpen` でパネルを管理します。`ThumbnailOptions.IconMode` は Finder 風の文書の枠を描画し、`Scale: 2` は Retina 画像を生成します。`Thumbnail` は呼び出し元の goroutine をブロックするため、goroutine またはバインドされたメソッドから呼び出してください。

### Workspace ヘルパー

`app.Browser` に、`NSWorkspace` を利用する 3 つのヘルパーが追加されます。`OpenWith` は、バンドル識別子またはバンドルパスで指定したアプリケーションでファイルを開きます。`ApplicationsForFile` は、ファイルを開けるインストール済みアプリケーションを、デフォルトのハンドラーを先頭にして列挙します。`ActivateApplication` は実行中のアプリケーションを前面に出します。

```go
apps := app.Browser.ApplicationsForFile("/Users/me/Documents/Report.md")
for _, info := range apps {
    log.Println(info.Name, info.BundleID, info.Path)
}

if err := app.Browser.OpenWith("/Users/me/Documents/Report.md", "com.apple.TextEdit"); err != nil {
    log.Println(err)
}

if err := app.Browser.ActivateApplication("com.apple.TextEdit"); err != nil {
    log.Println(err)
}
```

### Spotlight

`app.Spotlight.Index` は Core Spotlight を通じてアプリケーションのコンテンツをシステムの検索インデックスに追加します。各 `SearchableItem` には `ID` と `Title` があり、一括削除用の `Domain`、`Description`、`Keywords`、`ContentType`、PNG サムネイル、ディープリンクの `URL`、有効期限を省略可能で指定できます。ユーザーが Spotlight で項目を選ぶか、クエリを指定して「アプリで検索」を選ぶと、`OnOpen` が呼び出されます。

```go
err := app.Spotlight.Index([]application.SearchableItem{{
    ID:          "note:quarterly-report",
    Domain:      "notes",
    Title:       "Quarterly report",
    Description: "Draft for the board meeting",
    Keywords:    []string{"finance", "q3"},
    ContentType: "public.plain-text",
    URL:         "notes://open/quarterly-report",
}})
if err != nil {
    log.Println(err)
}

stop := app.Spotlight.OnOpen(func(_ *application.Context, id string, query string) {
    if query != "" {
        log.Println("search in app:", query)
        return
    }
    log.Println("open item", id)
})
defer stop()

// Later:
_ = app.Spotlight.Delete([]string{"note:quarterly-report"})
_ = app.Spotlight.DeleteDomain("notes")
```

`IsAvailable` はインデックスが項目を受け付けるかどうかを報告します。インデックス作成にはバンドル形式のアプリケーションが必要です。バンドル化されていない `go run` バイナリでインデックスに追加した項目は Spotlight に表示されません。`DeleteAll` はアプリがインデックスに追加したすべての項目を削除します。

## 今後の予定

シート、ポップオーバー、表示オプション、状態の復元を扱う `mac-windows-extra` サンプルを追加中です。公開されたらここにリンクを掲載します。

## バージョン要件

このページの機能はすべて macOS 専用で、Go API はどの環境でも同じです。Wails は macOS 10.13 以降を対象とします。新しいリリースを必要とする機能は、以下のとおり段階的に縮退します。

| 機能 | 最小 macOS バージョン | それ以前のリリースでの動作 |
|---------|---------------|-------------------------------|
| カメラとマイクの権限状態 | 10.14 | 許可済みとして報告される（それ以前のリリースではキャプチャデバイスへのアクセスを制限しない） |
| `Speech.Speak` と `Voices` | 10.14 | `ErrSpeechNotSupported` |
| 画面収録と入力監視の権限 | 10.15 | 許可済みとして報告される |
| `Speech.Recognize` | 10.15 | `ErrSpeechRecognitionNotSupported` |
| `QuickLook.Thumbnail` | 10.15 | `ErrQuickLookNotSupported` |
| `SetSubtitle` | 11 | デバッグログに記録され、無視される |
| `ExportPDF` | 11 | `ErrMacExportUnsupported` |
| メニュー項目とステータス項目の `SetSymbol` | 11 | 画像は表示されない |
| UTI による `AddContentType`、`SetFormats` | 11 | 同じ識別子が従来の許可ファイル型 API を通じて適用される |
| `StartDrag` のファイルプロミスのアイコン | 11 | 汎用の文書アイコン |
| `PowerState.LowPowerMode` とそのイベント | 12 | 常に false で、イベントは発火しない |
| `InteractionState` と `RestoreInteractionState` | 12 | `ErrMacInteractionStateUnsupported` |
| `OpenSystemSettings` の通知設定画面 | 13 | 従来の通知環境設定画面が開く |
| メニューのバッジ、セクションヘッダー、パレット | 14 | バッジは表示されず、ヘッダーは無効な項目になり、パレットは非表示になる |
| カスタムビューのない項目での `MacToolbarItem.ShowPopover` | 14 | `ErrMacPopoverAnchorUnavailable` |

その他の機能には、Wails の最小要件以外の要件はありません。

## Info.plist のキー

いくつかの機能は、アプリケーションの `Info.plist` にあるキーに依存します。使用目的の説明は権限の確認画面でユーザーに表示されます。キーがない場合、macOS は確認画面を表示せず、要求はタイムアウトします。

| キー | 必要とする機能 |
|-----|-----------|
| `NSCameraUsageDescription` | `Permissions.Request(PermissionKindCamera)`、ページからのカメラアクセス |
| `NSMicrophoneUsageDescription` | `Permissions.Request(PermissionKindMicrophone)`、ページからのマイクアクセス、`Speech.Recognize` |
| `NSSpeechRecognitionUsageDescription` | `Speech.Recognize` |
| `NSLocationUsageDescription` | `Permissions.Request(PermissionKindLocation)` |
| `NSSupportsSuddenTermination` | `Lifecycle.SuddenTerminationEnabled` の初期状態。`HoldTermination` はこれを一時停止する |
| `NSSupportsAutomaticTermination` | macOS がアイドル状態のアプリを終了できるようにする。`HoldTermination` はこれを一時停止する |
| `CFBundleLocalizations` | `Env.Locale().Identifier` が報告できる言語 |
| `NSServices` | `app.ServicesProvider` に登録したサービスごとに 1 項目。`InfoPlistXML` でブロックを生成する |
| `NSUserActivityTypes` | `app.Activity` を通じて公開または継続するすべての `UserActivity.Type` |
| `NSAppleScriptEnabled` と `OSAScriptingDefinition` | アプリがスクリプト操作可能であることを示し、`app.AppleEvents.ScriptingDefinition` から書き出した `.sdef` を指定する |
| `NSAppleEventsUsageDescription` | 他のアプリケーションへの `app.AppleEvents.Send` |

通知の権限と音声認識では、アプリをバンドル識別子を持つバンドルとして実行する必要もあります。バンドル化されていない `go run` バイナリでは、通知に対して `PermissionStatusUnsupported` が報告されます。

<a id="platform-notes"></a>

## プラットフォームに関する注記

このページのすべての API は Windows と Linux でもコンパイルできます。macOS 以外では、次のように動作します。

- ウィンドウの追加機能: `SetRepresentedFile`、`SetDocumentEdited`、`SetSubtitle`、`CascadeFrom`、`SetFrameAutosaveName`、`SetWindowButtonsOffset` は何もしません。`WindowCascade` は `WindowCentered` と同じ動作になります。`RequestAttention` は、`Cancel` が何もしないハンドルを返します。`PrintWithOptions` は `Print` を呼び出します。`ExportPDF` と `Snapshot` は `ErrMacOnly` を返します。
- メニュー: シンボル、バッジ、混合状態、代替項目、インデントは保存されますが描画されません。セクションヘッダーは無効な項目です。パレットは非表示です。「最近使った項目を開く」サブメニューは、メニューの作成時に Go の最近使った項目の一覧から設定されます。Dock メニューは表示されません。
- ダイアログ: `SetSuppression`、`SetHelp`、`SetNameFieldLabel`、`SetTags` は無視されます。`AddContentType` と `SetFormats` は、変換が分かっている場合、拡張子フィルターになります。`Prompt`、`PickColor`、`PickFont` は `ErrDialogNotSupported` を返します。
- ステータス項目: `SetSymbol`、`SetRemovable`、`OnVisibilityChange` は効果がありません。`IsVisible` は最後の `Show` または `Hide` の呼び出しを反映します。
- フィードバック: 触覚フィードバックは iOS と Android では動作し、Windows と Linux では何もしません。`Sound.Beep` と `Sound.Play` は Windows で WAV ファイルとレジストリのエイリアスに対応します。その他の環境では `Play` は `ErrSoundNotSupported` を返します。音声機能は `ErrSpeechNotSupported` と `ErrSpeechRecognitionNotSupported` を返します。
- クリップボード: 豊富な形式を扱うメソッドは `ErrClipboardNotSupported` を返し、`Types` は空、`ChangeCount` は 0 で、`OnChange` は発火しません。
- ドラッグ: `StartDrag` は `ErrDragOutUnsupported` を返します。ファイル以外のドロップ型は通知されません。ファイルのドロップは引き続き `WindowFilesDropped` を通じて動作します。
- システム: `Permissions.Status` は `PermissionStatusUnsupported` を報告し、`Request` は `ErrPermissionsUnsupported` を返します。`PreventSleep` は `ErrPreventSleepUnsupported` と、何もしない解放関数を返します。`HoldTermination` は何もしない解放関数を返します。`Accessibility` はすべて false、`KeyboardLayout` はゼロ値で、`Locale` は `LC_ALL`、`LC_MESSAGES`、`LANG` から導出されます。`Permissions` ウィンドウオプションはクロスプラットフォームです。
- 統合: `ServicesProvider.Register` は `ErrServicesUnsupported`、`Activity.Publish` は `ErrActivityUnsupported` を返します。`InfoPlistXML`、`InfoPlistEntries`、アクティビティハンドラーは引き続き動作します。`AppleEvents` の `Handle` と `Send` は `ErrAppleEventsNotSupported` を返します。`ScriptingDefinition` はどの環境でも生成されます。`QuickLook.Preview` と `Thumbnail` は `ErrQuickLookNotSupported` を返します。`Spotlight` のインデックス作成メソッドは `ErrSpotlightNotSupported` を返し、`OnOpen` は発火しません。`Browser.OpenWith` は指定された実行ファイルを、パスを引数として起動します。`ApplicationsForFile` は空で、`ActivateApplication` は `ErrApplicationNotRunning` を返します。
- 表示オプション: `SetPresentationOptions` は `ErrMacOnly` を返し、`PresentationOptions` は `MacPresentationDefault` です。
- シート: `PresentSheet`、`PresentCriticalSheet`、`PresentNativeSheet` は `ErrMacSheetUnsupported` を返します。`EndSheet` は何もせず、問い合わせメソッドはシートがないことを報告します。
- ポップオーバー: `NewMacPopover` は動作しますが、表示メソッドは `ErrMacPopoverUnsupported` を返し、`IsShown` は false です。
- 状態の復元: `SetRestorationID` と `SetRestorationData` は何もせず、`OnRestore` は呼び出されません。`InteractionState` と `RestoreInteractionState` は `ErrMacOnly` を返します。

## サンプル

各サンプルは、そのまま実行できる完全なアプリケーションです。

- [`v3/examples/mac-toolbar`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-toolbar): `windowextras.go` は、変更済みを示す点、サブタイトル、PDF エクスポート、印刷オプション、カスケード配置されたウィンドウをメモエディタに組み込みます。
- [`v3/examples/mac-menus-dock`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-menus-dock): シンボル、バッジ、セクションヘッダー、混合状態、代替項目、パレット、「最近使った項目を開く」、動的な Dock メニュー、Dock の進捗表示。
- [`v3/examples/mac-dialogs`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-dialogs): 再表示を抑制するチェックボックスとヘルプボタン、テキスト入力プロンプト、コンテンツ型、形式のポップアップ、Finder タグ、カラー・フォントパネル。
- [`v3/examples/mac-feedback`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-feedback): 取り外し可能な SF Symbol のステータス項目、触覚フィードバック、システムサウンド、テキスト読み上げ、音声認識。
- [`v3/examples/mac-clipboard-drag`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-clipboard-drag): 変更追跡付きの豊富な形式のクリップボード、ファイルプロミスを使う外向きのドラッグ、テキスト、URL、画像のドロップ。
- [`v3/examples/mac-system`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-system): 権限、スリープ防止、終了の保留、電源状態、アクセシビリティ、キーボードレイアウト、ロケールと、それらのリアルタイム更新。
- [`v3/examples/mac-integration`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-integration): Services メニュー項目、継続処理ハンドラーを備えた Handoff アクティビティ、`app.Browser` の Workspace ヘルパー。
- [`v3/examples/mac-search-preview`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-search-preview): `OnOpen` を使った Spotlight のインデックス作成、Quick Look のプレビューとサムネイル、スクリプト定義を備えたカスタム Apple Event。
