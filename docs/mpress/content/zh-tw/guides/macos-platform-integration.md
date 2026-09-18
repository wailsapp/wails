---
title: "macOS 平台整合"
description: "macOS 上的文件視窗、Dock 與選單擴充功能、原生面板、狀態項目、回饋功能、豐富的剪貼簿功能、向外拖曳、權限、電源與生命週期"
slug: "guides/macos-platform-integration"
sourcePath: "guides/macos-platform-integration.md"
---

適用平台：macOS

Wails v3 讓您的應用程式具備 macOS 使用者期待的原生行為：具有代理圖示與階梯式排列的文件視窗、帶有圖示與徽章的選單、顯示進度的 Dock 選單、原生警示和面板、可移除的狀態項目、觸覺回饋與語音、支援向外拖曳的豐富剪貼簿功能、權限、電源和地區設定等系統資訊，以及與「服務」選單、Handoff、AppleScript 和 Quick Look 的整合。所有功能都透過 `application` 套件從 Go 控制。

相同程式碼也能在 Windows 和 Linux 上編譯。設定方法會儲存其值，查詢會傳回零值，而需要 macOS 的操作會傳回文件所述的錯誤，例如 `ErrMacOnly`、`ErrDialogNotSupported` 或 `ErrClipboardNotSupported`。下方的[平台說明](#platform-notes)列出各功能在 macOS 以外平台上的行為。

關於原生視窗介面（工具列、側邊欄、檢閱器、附加控制列和視窗分頁），請參閱[原生 macOS 視窗介面](/guides/macos-native-chrome)指南。

## 文件視窗

文件視窗會在標題列顯示其代表的檔案、在關閉按鈕上以圓點標示尚未儲存的變更，並以階梯式排列開啟新視窗。這些功能都是 `WebviewWindow` 的方法及 `MacWindow` 的選項。

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

- `SetRepresentedFile` 會在標題列顯示檔案的代理圖示。使用者可將圖示拖曳至其他應用程式，或按住 Command 點選圖示以顯示路徑。傳入 `""` 可移除圖示。`RepresentedFile` 可讀回其值。
- `SetDocumentEdited` 會在關閉按鈕顯示未儲存變更的圓點，並將代理圖示變暗。`IsDocumentEdited` 可讀回其值。
- `SetSubtitle` 會在 macOS 11 及更新版本的標題下方顯示第二行。
- `InitialPosition: application.WindowCascade` 會將視窗放在上一個階梯式排列視窗的右下方，如同開啟新文件一樣。`CascadeFrom(other)` 可對現有視窗執行相同操作，並更新後續視窗的階梯式排列位置。
- `Mac.FrameAutosaveName` 會在視窗首次顯示前還原儲存的位置與大小，並在視窗移動時持續儲存。還原的視窗範圍優先於 `X`、`Y`、`Width`、`Height` 和 `InitialPosition`。`SetFrameAutosaveName` 可變更現有視窗使用的名稱。
- `MacTitleBar.WindowButtonsOffset` 會以指定點數移動關閉、最小化和縮放按鈕。`SetWindowButtonsOffset` 和 `ResetWindowButtonsOffset` 可在執行時變更此設定。

這三個設定方法都可在原生視窗建立前呼叫；其值會在建立視窗時套用。

### 請求使用者注意

應用程式在背景執行時，`RequestAttention` 會讓 Dock 圖示彈跳。資訊性要求會彈跳一次。緊急要求會持續彈跳，直到使用者啟用應用程式或您取消要求。

```go
request := window.RequestAttention(true)

// Once the work that needed attention is done:
request.Cancel()
```

`Flash` 仍是跨平台請求使用者注意一次的方法。

### 列印與匯出

`PrintWithOptions` 會以明確指定的頁面設定列印 WebView。使用零值時，會以共用列印設定顯示列印面板。`Print` 保持其既有行為（橫向、30 點邊界）。

`ExportPDF` 會將頁面轉譯為 PDF 文件，`Snapshot` 則會擷取為 PNG。兩者都會等待 WebKit，因此請從 goroutine 呼叫，絕不可從應用程式執行緒呼叫；在該執行緒呼叫會傳回 `ErrMacExportOnMainThread`。

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

`PrintOptions` 也接受 `PrinterName`、`PaperName`（例如 `"iso-a4"` 這類 PostScript 名稱）及 `Scale`。`PDFExportOptions` 和 `SnapshotOptions` 可接受選用的 `Rect` 來限制擷取範圍，以及預設為 `DefaultMacExportTimeout`（30 秒）的 `Timeout`。

## 附屬面板

附屬面板是附加在父視窗頂部的第二個視窗，類似儲存面板。任何 `WebviewWindow` 都能透過 `PresentSheet` 作為另一個視窗的附屬面板顯示，並透過 `EndSheet` 搭配回應碼結束；回應碼會傳至附屬面板的 `OnSheetEnd` 回呼。建立附屬面板視窗時請設定 `Hidden`，避免附加前在畫面上閃現；結束時 AppKit 會再次將其移出畫面，因此同一個視窗可重複顯示。

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

`PresentCriticalSheet` 會將附屬面板顯示在任何已附加的一般附屬面板前方，而非排隊等在其後。`PresentNativeSheet` 對 `NativeWindow` 執行相同操作。`IsSheet`、`SheetParent`、`AttachedSheet`、`AttachedNativeSheet` 和 `HasAttachedSheet` 可描述目前狀態。若關閉附屬面板視窗而非結束附屬面板，會傳送 `MacSheetResponseStop`。

## 彈出面板

`MacPopover` 是 `NSPopover`：暫時顯示的面板，可錨定至視窗中的矩形、工具列項目或選單列中的狀態項目。其內容是原生 `MacAccessory` 控制列，與[原生 macOS 視窗介面](/guides/macos-native-chrome)指南用於標題列附加控制列的型別相同。請在首次顯示前加入所有控制項。

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

`MacToolbarItem.ShowPopover` 和 `SystemTray.ShowPopover` 可將同一個彈出面板錨定至工具列項目或狀態項目。`MacPopoverBehaviorTransient` 會在面板外任何位置被點選時關閉面板，`Semitransient` 則只會在呈現它的視窗中發生點選時關閉，而預設行為會讓面板保持開啟，直到呼叫 `Close`。`MacRectEdge` 決定彈出面板出現在哪一側。`SetContentSize` 和 `SetBehavior` 可調整已顯示的彈出面板，而 `Destroy` 會釋放原生彈出面板，讓內容控制列可在其他地方使用。

## 狀態還原

macOS 會在應用程式當機、強制結束或重新開機後重新開啟其視窗；若系統設定中的「結束應用程式時關閉視窗」已關閉，正常結束後也會重新開啟。請為視窗指定 `Mac.RestorationID`，使用 `SetRestorationData` 儲存重建視窗所需資料，並註冊 `app.Window.OnRestore`，以便下次啟動時再次建立視窗。

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

只有應用程式終止時可見的視窗會被儲存。`RestorationState.Data` 是字串映射；請只用它儲存識別碼、路徑和位置。`InteractionState` 會將 WebView 的上一頁／下一頁清單及捲動位置傳回為不透明資料區塊（macOS 12+）；`RestoreInteractionState` 可將其套用到重建的視窗，通常會以 base64 編碼後儲存在還原資料中。`SetRestorationID` 和 `RestorationID` 可變更及讀取現有視窗的識別碼。

## 呈現選項

`MacPresentationOptions` 對應 `NSApplication.presentationOptions`：它是位元遮罩，可在應用程式作用期間隱藏 Dock 或選單列，並停用切換處理程序、強制結束、登出或「隱藏」命令。在應用程式選項中設定 `Mac.PresentationOptions` 可於啟動時套用，或使用 `SetPresentationOptions` 在執行時變更。無效的組合會在傳至 AppKit 前遭拒絕，並傳回包裝了 `ErrMacPresentationOptionsInvalid` 的錯誤。

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

隱藏選單列（`HideMenuBar` 或 `AutoHideMenuBar`）需要同時設定一個 Dock 選項，而 `AutoHideToolbar` 同時需要 `FullScreen` 和 `AutoHideMenuBar`。`Validate` 會回報值違反的第一項規則，`Has` 則可檢查個別旗標。

## 選單與 Dock

選單項目可使用 SF Symbols、徽章、區段標題、色彩調色盤、混合勾選狀態、替代項目和縮排。這些都是 `MenuItem` 和 `Menu` 的方法，因此可用於應用程式選單、快顯選單、系統匣選單和 Dock 選單。

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

- `SetSymbol` 會在標題旁顯示 SF Symbol（macOS 11+）。它會取代透過 `SetBitmap` 設定的影像。
- `SetBadge` 會顯示數量，`SetBadgeText` 則會在標題後顯示簡短字串（macOS 14+）。`ClearBadge` 可將其移除；`BadgeCount` 和 `BadgeText` 可讀回其值。
- `AddSectionHeader` 會加入無法互動的標題（macOS 14+）。在較舊版本上，它是具有相同標題的停用項目。
- `AddPalette` 會加入由 `NSMenu` 調色盤選單支援的一列色票（macOS 14+）。請為每個色票傳入一個圖示（每種顏色一個），或傳入空切片以顯示實心圓。沒有標籤時，調色盤會直接顯示在父選單中；`SetLabel` 則會將其呈現為有標題的子選單。`PaletteSelected` 會傳回選取的索引。
- `SetMixed` 會將核取方塊設為以短橫線繪製的混合狀態。點選後會依 AppKit 的行為將其完全勾選。
- `SetAlternate(true)` 會在按住不同修飾鍵時，以該項目取代其上方項目。這兩個項目必須使用相同按鍵，且修飾鍵不同。
- `SetIndentationLevel` 可將標題縮排最多 15 個層級。

### 開啟最近使用的項目

`fileMenu.AddRole(application.OpenRecent)` 會加入標準的「開啟最近使用的項目」子選單。在 macOS 上，`NSDocumentController` 每次開啟該選單時都會填入內容，並包含「清除選單」項目。使用 `app.Menu.AddRecentDocument` 加入檔案、`RecentDocuments` 列出檔案，以及 `ClearRecentDocuments` 清空清單。清單在重新啟動後仍會保留。

選擇最近使用的檔案會產生與從 Finder 開啟檔案相同的事件：

```go
app.Event.OnApplicationEvent(events.Common.ApplicationOpenedWithFile, func(event *application.ApplicationEvent) {
    path := event.Context().Filename()
    log.Println("open", path)
})
```

### Dock 選單

`app.Menu.SetDockMenu` 會安裝在 Dock 圖示上按右鍵時顯示的靜態選單。`OnDockMenu` 則會在每次即將顯示時依需求建立選單，適合項目須反映變動狀態的情況。若有建構函式，它會優先於靜態選單。

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

### Dock 進度

Dock 服務除了現有的徽章支援，還可在 Dock 圖示上繪製進度列。將 `dock.New()` 註冊為服務，並呼叫 `SetProgress`，傳入介於 0 與 1 之間的比例值。

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

`GetProgress` 會傳回目前比例值；未顯示進度列時則傳回 `nil`。

## 對話框

訊息、開啟和儲存對話框都接受 macOS 選項；對話框管理器還新增文字輸入提示，以及系統色彩與字型面板。

### 警示

`SetSuppression` 會加入「不要再顯示此訊息」核取方塊，`SetHelp` 則會顯示說明按鈕。可在按鈕回呼中使用 `Suppressed` 讀取核取方塊狀態，或註冊 `OnSuppression` 先接收其狀態。

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

### 文字輸入提示

`Prompt` 會顯示含文字欄位的警示，並阻塞直到警示關閉，因此請從 goroutine 或繫結的方法呼叫。`Secure` 會將欄位變成密碼欄位，`Window` 則會將警示呈現為附屬面板。

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

### 檔案面板

`AddContentType` 可在開啟及儲存對話框中依統一型別識別碼篩選。它與 `AddFilter` 並用，因此 `"public.image"` 可比對系統已知的所有影像型別，同時仍可透過篩選條件依副檔名比對 PDF。

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

儲存面板新增格式彈出式選單、名稱欄位的自訂標籤，以及 Finder 標籤。使用者變更彈出式選單時，`SetFormats` 會切換允許的型別及名稱欄位的副檔名；`SelectedFormat` 則會回報最終選擇。

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

### 色彩與字型面板

`PickColor` 和 `PickFont` 會開啟共用系統面板，並阻塞直到面板關閉。面板開啟期間，`OnChange` 會傳送每次選取結果，讓頁面即時預覽。每種類型一次只能開啟一個面板；第二次呼叫會傳回 `ErrDialogInProgress`。

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

## 狀態項目與回饋

### 狀態項目

macOS 上的系統匣項目是 `NSStatusItem`。它可使用 SF Symbol 繪製、顯示工具提示，也可像內建項目一樣由使用者移除。

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

- `SetSymbol` 會將圖示轉譯為樣板影像，使其隨選單列外觀變化（macOS 11+）。`SetSymbolConfiguration` 可設定點數大小和字重。
- `SetTooltip` 設定懸停文字。`Tooltip` 可讀回其值。
- `SetRemovable(true, name)` 讓使用者可按住 Command，將項目拖出選單列。請提供穩定的自動儲存名稱，讓 macOS 跨啟動記住移除狀態。`Show` 或 `SetVisible(true)` 可將其帶回。
- `IsVisible` 讀取 `NSStatusItem.visible`，因此使用者移除項目後會傳回 false。`OnVisibilityChange` 會回報每次變更。

### 觸覺回饋

應用程式作用期間，`app.Haptics.Perform` 會在 Force Touch 觸控式軌跡板或 Magic Trackpad 上播放回饋模式。

```go
app.Haptics.Perform(application.HapticAlignment)
```

種類包括 `HapticGeneric`、`HapticAlignment`（項目對齊定位）和 `HapticLevelChange`（段位或點按階段）。`IsSupported` 會回報平台是否支援播放回饋。

### 音效

`app.Sound` 可播放警示音、指定名稱的系統音效或音訊檔案。

```go
app.Sound.Beep()

if err := app.Sound.Play("Glass"); err != nil {
    log.Println(err)
}

for _, name := range app.Sound.SystemSounds() {
    log.Println(name)
}
```

`Play` 接受 `SystemSounds` 中的名稱，或 Core Audio 可解碼的任何檔案之絕對路徑。`PlayData` 會從記憶體播放完整音訊檔案。

### 語音

`app.Speech.Speak` 會將文字排入佇列，交由系統語音朗讀，並傳回 `Utterance`。語音內容會依序播放；`Stop` 會移除其中一段，`StopAll` 則會清空佇列。`Voices` 會列出已安裝的語音及其識別碼和語言。

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

`Recognize` 會使用 `SFSpeechRecognizer` 轉錄預設麥克風的輸入。首次呼叫會要求麥克風與語音辨識權限，並阻塞直到使用者回應，因此請從 goroutine 呼叫。部分轉錄文字會透過 `OnPartial` 傳送；`Stop` 會結束擷取並傳回最終文字。

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

語音辨識需要已封裝的應用程式，且其 `Info.plist` 必須宣告 `NSSpeechRecognitionUsageDescription` 和 `NSMicrophoneUsageDescription`。缺少這些鍵時，macOS 會拒絕存取，而 `Recognize` 會傳回 `ErrSpeechRecognitionUsageDescription`。

## 剪貼簿與拖曳

### 豐富的剪貼簿功能

除了純文字，`app.Clipboard` 也能依任何統一型別識別碼讀寫影像、檔案參照、HTML、RTF 和原始資料。`Types` 會列出剪貼簿上的型別，`OnChange` 則會回報任何應用程式造成的變更。

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

`SetImage` 和 `Image` 使用 PNG 位元組；其他應用程式以 TIFF 複製的影像會自動轉換。系統不會通知剪貼簿變更，因此只要至少有一個監聽器，`OnChange` 就會每 500 ms 輪詢變更計數。

### 向外拖曳

`StartDrag` 會從視窗開始系統拖曳，效果如同使用者從 Finder 拿起項目。它可提供現有檔案、只在目的地接受放置時才產生內容的檔案承諾，或純文字。請在滑鼠手勢期間啟動：繫結 Go 方法，並從頁面上可拖曳元素的 `mousedown` 或 `pointerdown` 處理常式呼叫，同時將 HTML `draggable` 屬性設為 `false`，避免 WebKit 開始自己的拖曳操作。

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

在手勢以外呼叫 `StartDrag` 會傳回 `ErrDragOutNoGesture`。`DragItems.Image` 和 `ImageOffset` 可設定游標下方的圖片。

### 來自其他應用程式的放置

檔案放置仍使用 `WindowFilesDropped` 事件。若要接受從其他應用程式拖入的文字、URL 或影像，請在 `DropTypes` 中列出型別，並註冊 `OnDrop`。這些放置事件會傳送至 Go，而非頁面本身的 HTML5 放置處理常式。

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

## 系統

### 權限

`app.Permissions` 可回報及要求系統隱私權限（相機、麥克風、螢幕錄製、輔助使用、定位、通知、輸入監控和完整磁碟存取）。`Status` 絕不會顯示權限提示。`Request` 會針對狀態尚未確定的權限種類顯示提示，並阻塞直到使用者回應，因此請從 goroutine 呼叫。`OpenSystemSettings` 會開啟對應的「隱私權與安全性」設定面板。

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

完整磁碟存取權無法透過程式要求，會傳回 `ErrPermissionNotRequestable`；請引導使用者前往設定面板。如果要求始終未獲回應，會傳回 `ErrPermissionRequestTimeout`；在 macOS 上，這通常表示 `Info.plist` 缺少該權限種類的用途說明鍵。

macOS 現在會遵循 `Permissions` 視窗選項。它決定如何處理頁面發出的 `getUserMedia` 要求：`PermissionAllow` 會略過 WebView 自身的提示，`PermissionDeny` 會直接拒絕，`PermissionDefault` 則會顯示提示。首次使用相機或麥克風時，系統層級的 TCC 提示仍會出現。

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

`app.Power.PreventSleep` 會讓系統保持喚醒；若使用 `Display`，也會讓螢幕保持開啟，直到呼叫傳回的釋放函式。保持喚醒的要求會累計，因此應用程式的多個部分可同時提出要求。原因會顯示在「活動監視器」中。

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

變更會以 `events.Mac.ApplicationDidChangePowerState`（低耗電模式切換）和 `events.Mac.ApplicationDidChangeThermalState` 事件傳送。

### 生命週期

應用程式若透過 `NSSupportsSuddenTermination` 選擇加入，macOS 可在登出或關機時立即終止閒置的應用程式；若透過 `NSSupportsAutomaticTermination` 選擇加入，macOS 可結束沒有視窗的閒置應用程式。`app.Lifecycle.HoldTermination` 可在儲存檔案等關鍵區段暫停這兩種終止行為。

```go
release := app.Lifecycle.HoldTermination("Saving document")
defer release()
// write the file
```

`SetSuddenTerminationEnabled` 可在執行時切換突然終止功能；`SuddenTerminationEnabled` 會回報目前狀態，其初始值取自 `Info.plist` 鍵。

### 環境

`app.Env` 新增三項使用者設定查詢。

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

- `Accessibility` 反映「減少動態效果」、「減少透明度」、「增加對比」、「不依賴顏色區分」、「反轉顏色」、VoiceOver 和「切換控制」設定。
- `KeyboardLayout` 會傳回作用中的輸入來源，以及其識別碼、在地化名稱和語言。
- `Locale` 會傳回 AppKit 為應用程式選擇的地區設定，以及使用者完整且有順序的 `Preferred` 清單。`Identifier` 只會反映套件在 `CFBundleLocalizations` 中宣告的語言；若要自行選擇語言，請使用 `Preferred`。

### 事件

以下是新增的應用程式事件。每個事件都透過 `app.Event.OnApplicationEvent` 傳送；請查詢對應的管理器以取得最新值。

| 事件 | 觸發時機 | 讀取方式 |
|-------|------------|-----------|
| `events.Mac.ApplicationDidChangePowerState` | 切換低耗電模式 | `app.Power.State()` |
| `events.Mac.ApplicationDidChangeThermalState` | 散熱壓力變更 | `app.Power.State()` |
| `events.Common.AccessibilitySettingsChanged` | 輔助使用顯示設定變更 | `app.Env.Accessibility()` |
| `events.Mac.ApplicationDidChangeAccessibilitySettings` | 相同變更的 macOS 形式 | `app.Env.Accessibility()` |
| `events.Mac.ApplicationDidChangeKeyboardLayout` | 輸入來源變更 | `app.Env.KeyboardLayout()` |
| `events.Mac.ApplicationDidChangeLocale` | 地區設定變更 | `app.Env.Locale()` |

```go
app.Event.OnApplicationEvent(events.Mac.ApplicationDidChangeThermalState, func(*application.ApplicationEvent) {
    app.Event.Emit("system:power", app.Power.State())
})
app.Event.OnApplicationEvent(events.Common.AccessibilitySettingsChanged, func(*application.ApplicationEvent) {
    app.Event.Emit("system:accessibility", app.Env.Accessibility())
})
```

## 整合

### 服務選單

`app.ServicesProvider.Register` 會在每個 macOS 應用程式針對選取文字或檔案顯示的「服務」子選單中加入項目。處理常式會以 `ServiceRequest` 接收剪貼簿，並傳回 `ServiceResponse` 以寫回內容；空回應會讓選取內容保持不變。AppKit 會在主執行緒等待處理常式，因此請盡快完成。

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

`Name` 是 AppKit 傳送的訊息，且必須是一般識別碼。`SendTypes` 和 `ReturnTypes` 是剪貼簿型別；服務至少需要其中一項。僅註冊服務不會讓它顯示：應用程式套件的 `Info.plist` 必須在 `NSServices` 下宣告。`InfoPlistXML` 會傳回可直接貼上的區塊，`InfoPlistEntries` 則以映射形式傳回相同資料，供 plist 序列化工具使用。這些項目中的 `NSPortName` 是應用程式的 `Name`，必須與 `CFBundleName` 相符。

使用 Wails CLI 建置的專案可在 `build/config.yml` 中宣告服務一次，由封裝程序產生 `NSServices` 區塊：

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

每個項目都必須與 Go 中註冊、具有相同 `Name` 的 `ServiceDefinition` 相符。安裝新建置版本後，請執行 `pbs -update`，讓「服務」選單無須登出即可取得變更。

### Handoff 與使用者活動

`app.Activity.Publish` 會將 `NSUserActivity` 設為目前活動，讓使用者可在其他裝置上繼續、透過 Spotlight 找到，或由 Siri 建議。傳回的 `PublishedActivity` 可隨狀態變更更新，並可在文件關閉時使其失效。發佈新活動會取代先前的活動。

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

傳入的活動會透過 `OnContinue` 傳送。通用連結的型別為 `UserActivityTypeBrowsingWeb`，頁面位於 `WebpageURL`；它們也會以 `events.Common.ApplicationLaunchedWithUrl` 傳送，讓應用程式可共用同一套 URL 處理流程。`OnWillContinue`、`OnFailed` 和 `OnUpdated` 涵蓋委派的其餘部分。

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

每種活動型別都必須列於 `Info.plist` 的 `NSUserActivityTypes` 下。通用連結還需要包含 `applinks:example.com` 項目的 `com.apple.developer.associated-domains` 權限，以及該網域上相符的 `apple-app-site-association` 檔案。

### Apple Events

`app.AppleEvents.Handle` 會為事件類別與 ID 註冊處理常式，讓 AppleScript、捷徑和其他應用程式可控制此應用程式。代碼是四個字元的字串。直接參數會解碼為 Go 值（`string`、檔案路徑的 `[]string`、`int64`、`float64`、`bool`、`[]any` 或 `AppleEventRawData`），回覆的 `Result` 接受相同種類。事件暫停期間，處理常式會在各自的 goroutine 上執行。

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

指令碼可使用原始事件語法直接呼叫處理常式：

```applescript
tell application id "com.example.notes" to «event WAILnote» "hello"
```

`ScriptingDefinition` 會產生最小化的 `.sdef`，為每個處理常式提供命令名稱。請將它放入 `Contents/Resources`，並透過 `NSAppleScriptEnabled` 和 `OSAScriptingDefinition` 在 `Info.plist` 中指向它；之後「腳本編輯程式」便會在「檔案 > 開啟辭典」下顯示。Wails 已處理自訂 URL 配置的「取得 URL」事件；針對 `"GURL"`/`"GURL"` 的處理常式會接續其處理流程，而針對 `"aevt"`/`"odoc"` 的處理常式則會取代內建的「開啟文件」事件傳送。`Send` 會依套件識別碼指定執行中的應用程式；在已封裝的應用程式中需要 `NSAppleEventsUsageDescription`，並會阻塞呼叫端 goroutine，直到收到回覆。

### Quick Look

`app.QuickLook.Preview` 可針對一個或多個檔案開啟共用 Quick Look 面板；若有多個路徑，面板會顯示切換檔案的箭頭。`Thumbnail` 會透過系統縮圖提供者轉譯檔案並傳回 PNG，因此文件、影像、PDF 和影片都適用。

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

路徑必須是絕對路徑，且檔案必須存在。`ClosePreview` 和 `IsPreviewOpen` 可管理面板。`ThumbnailOptions.IconMode` 會繪製 Finder 風格的文件邊框，而 `Scale: 2` 會產生 Retina 影像。`Thumbnail` 會阻塞呼叫端 goroutine，因此請從 goroutine 或繫結的方法呼叫。

### 工作區輔助方法

`app.Browser` 新增三個與 `NSWorkspace` 相關的輔助方法。`OpenWith` 會使用以套件識別碼或套件路徑指定的應用程式開啟檔案。`ApplicationsForFile` 會列出可開啟檔案的已安裝應用程式，預設處理程式排在第一位。`ActivateApplication` 會將執行中的應用程式帶到最前方。

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

`app.Spotlight.Index` 會透過 Core Spotlight 將應用程式內容加入系統搜尋索引。每個 `SearchableItem` 都有 `ID`、`Title`，並可選擇性提供用於批次移除的 `Domain`、`Description`、`Keywords`、`ContentType`、PNG 縮圖、深層連結 `URL` 和到期時間。使用者在 Spotlight 中選取其中一個項目，或使用查詢選擇「在 App 中搜尋」時，會呼叫 `OnOpen`。

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

`IsAvailable` 會回報索引是否接受項目。建立索引需要已封裝的應用程式：由未封裝的 `go run` 二進位檔建立索引的項目，不會出現在 Spotlight 中。`DeleteAll` 會移除應用程式建立索引的所有項目。

## 即將推出

涵蓋附屬面板、彈出面板、呈現選項和狀態還原的 `mac-windows-extra` 範例正在新增中，完成後會在此提供連結。

## 版本需求

本頁所有功能都僅適用於 macOS，而 Go API 在所有平台上都相同。Wails 支援 macOS 10.13 及更新版本；需要較新版本的功能會依下表所述降級。

| 功能 | 最低 macOS 版本 | 較舊版本的行為 |
|---------|---------------|-------------------------------|
| 相機和麥克風權限狀態 | 10.14 | 回報為已授權（較舊版本不會限制擷取裝置） |
| `Speech.Speak` 和 `Voices` | 10.14 | `ErrSpeechNotSupported` |
| 螢幕錄製和輸入監控權限 | 10.15 | 回報為已授權 |
| `Speech.Recognize` | 10.15 | `ErrSpeechRecognitionNotSupported` |
| `QuickLook.Thumbnail` | 10.15 | `ErrQuickLookNotSupported` |
| `SetSubtitle` | 11 | 忽略，並記錄偵錯日誌 |
| `ExportPDF` | 11 | `ErrMacExportUnsupported` |
| 選單項目和狀態項目上的 `SetSymbol` | 11 | 不顯示影像 |
| 依 UTI 使用的 `AddContentType`、`SetFormats` | 11 | 相同識別碼會透過舊版允許檔案型別 API 套用 |
| `StartDrag` 中的檔案承諾圖示 | 11 | 通用文件圖示 |
| `PowerState.LowPowerMode` 及其事件 | 12 | 永遠為 false；事件不會觸發 |
| `InteractionState` 和 `RestoreInteractionState` | 12 | `ErrMacInteractionStateUnsupported` |
| `OpenSystemSettings` 中的通知面板 | 13 | 開啟舊版通知偏好設定面板 |
| 選單徽章、區段標題、調色盤 | 14 | 不顯示徽章；標題為停用項目；隱藏調色盤 |
| 沒有自訂檢視的項目上的 `MacToolbarItem.ShowPopover` | 14 | `ErrMacPopoverAnchorUnavailable` |

其他所有功能除了 Wails 的最低版本需求外，沒有其他版本需求。

## Info.plist 鍵

數項功能仰賴應用程式 `Info.plist` 中的鍵。用途說明會在權限提示中顯示給使用者；若缺少對應的鍵，macOS 就不會顯示提示，而要求會逾時。

| 鍵 | 所需功能 |
|-----|-----------|
| `NSCameraUsageDescription` | `Permissions.Request(PermissionKindCamera)`、頁面的相機存取 |
| `NSMicrophoneUsageDescription` | `Permissions.Request(PermissionKindMicrophone)`、頁面的麥克風存取、`Speech.Recognize` |
| `NSSpeechRecognitionUsageDescription` | `Speech.Recognize` |
| `NSLocationUsageDescription` | `Permissions.Request(PermissionKindLocation)` |
| `NSSupportsSuddenTermination` | `Lifecycle.SuddenTerminationEnabled` 的初始狀態；`HoldTermination` 可暫停此功能 |
| `NSSupportsAutomaticTermination` | 允許 macOS 結束閒置應用程式；`HoldTermination` 可暫停此功能 |
| `CFBundleLocalizations` | `Env.Locale().Identifier` 可回報的語言 |
| `NSServices` | 每個透過 `app.ServicesProvider` 註冊的服務各需一個項目；可使用 `InfoPlistXML` 產生區塊 |
| `NSUserActivityTypes` | 透過 `app.Activity` 發佈或繼續的每個 `UserActivity.Type` |
| `NSAppleScriptEnabled` 和 `OSAScriptingDefinition` | 將應用程式標記為可由指令碼控制，並指定由 `app.AppleEvents.ScriptingDefinition` 寫出的 `.sdef` |
| `NSAppleEventsUsageDescription` | 對其他應用程式呼叫 `app.AppleEvents.Send` |

通知權限與語音辨識也要求應用程式以具有套件識別碼的套件形式執行；未封裝的 `go run` 二進位檔會對通知回報 `PermissionStatusUnsupported`。

<a id="platform-notes"></a>

## 平台說明

本頁所有 API 都能在 Windows 和 Linux 上編譯。在 macOS 以外的平台：

- 視窗附加功能：`SetRepresentedFile`、`SetDocumentEdited`、`SetSubtitle`、`CascadeFrom`、`SetFrameAutosaveName` 和 `SetWindowButtonsOffset` 不執行任何動作。`WindowCascade` 的行為與 `WindowCentered` 相同。`RequestAttention` 傳回的控制代碼，其 `Cancel` 不執行任何動作。`PrintWithOptions` 會呼叫 `Print`。`ExportPDF` 和 `Snapshot` 傳回 `ErrMacOnly`。
- 選單：圖示、徽章、混合狀態、替代項目和縮排會儲存，但不會繪製。區段標題是停用項目。調色盤會隱藏。「開啟最近使用的項目」子選單會在建立選單時從 Go 的最近使用清單填入。Dock 選單不會顯示。
- 對話框：`SetSuppression`、`SetHelp`、`SetNameFieldLabel` 和 `SetTags` 會被忽略。若已知對應轉換方式，`AddContentType` 和 `SetFormats` 會轉為副檔名篩選條件。`Prompt`、`PickColor` 和 `PickFont` 傳回 `ErrDialogNotSupported`。
- 狀態項目：`SetSymbol`、`SetRemovable` 和 `OnVisibilityChange` 不生效。`IsVisible` 反映最近一次 `Show` 或 `Hide` 呼叫。
- 回饋：觸覺回饋支援 iOS 和 Android；在 Windows 和 Linux 上不執行任何動作。`Sound.Beep` 和 `Sound.Play` 可在 Windows 上搭配 WAV 檔案與登錄別名運作；在其他平台上，`Play` 傳回 `ErrSoundNotSupported`。語音功能傳回 `ErrSpeechNotSupported` 和 `ErrSpeechRecognitionNotSupported`。
- 剪貼簿：豐富內容方法傳回 `ErrClipboardNotSupported`，`Types` 為空，`ChangeCount` 為 0，且 `OnChange` 永不觸發。
- 拖曳：`StartDrag` 傳回 `ErrDragOutUnsupported`。非檔案放置型別不會傳送；檔案放置仍可透過 `WindowFilesDropped` 運作。
- 系統：`Permissions.Status` 回報 `PermissionStatusUnsupported`，`Request` 傳回 `ErrPermissionsUnsupported`。`PreventSleep` 傳回 `ErrPreventSleepUnsupported`，並提供不執行任何動作的釋放函式。`HoldTermination` 傳回不執行任何動作的釋放函式。`Accessibility` 的所有值都是 false，`KeyboardLayout` 為零值，`Locale` 則由 `LC_ALL`、`LC_MESSAGES` 和 `LANG` 推導。`Permissions` 視窗選項可跨平台使用。
- 整合：`ServicesProvider.Register` 傳回 `ErrServicesUnsupported`，`Activity.Publish` 傳回 `ErrActivityUnsupported`；`InfoPlistXML`、`InfoPlistEntries` 和活動處理常式仍可運作。`AppleEvents` 的 `Handle` 和 `Send` 傳回 `ErrAppleEventsNotSupported`；`ScriptingDefinition` 可在所有平台產生。`QuickLook.Preview` 和 `Thumbnail` 傳回 `ErrQuickLookNotSupported`。`Spotlight` 建立索引的方法傳回 `ErrSpotlightNotSupported`，且 `OnOpen` 永不觸發。`Browser.OpenWith` 會啟動指定的可執行檔，並將路徑作為引數；`ApplicationsForFile` 為空，`ActivateApplication` 傳回 `ErrApplicationNotRunning`。
- 呈現選項：`SetPresentationOptions` 傳回 `ErrMacOnly`，而 `PresentationOptions` 為 `MacPresentationDefault`。
- 附屬面板：`PresentSheet`、`PresentCriticalSheet` 和 `PresentNativeSheet` 傳回 `ErrMacSheetUnsupported`；`EndSheet` 不執行任何動作，查詢方法則回報沒有附屬面板。
- 彈出面板：`NewMacPopover` 可運作，顯示方法傳回 `ErrMacPopoverUnsupported`，而 `IsShown` 為 false。
- 狀態還原：`SetRestorationID` 和 `SetRestorationData` 不執行任何動作，`OnRestore` 永不呼叫，而 `InteractionState` 和 `RestoreInteractionState` 傳回 `ErrMacOnly`。

## 範例

每個範例都是完整且可執行的應用程式：

- [`v3/examples/mac-toolbar`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-toolbar)：`windowextras.go` 將已修改圓點、副標題、PDF 匯出、列印選項和階梯式排列視窗整合進筆記編輯器。
- [`v3/examples/mac-menus-dock`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-menus-dock)：圖示、徽章、區段標題、混合狀態、替代項目、調色盤、「開啟最近使用的項目」、動態 Dock 選單和 Dock 進度。
- [`v3/examples/mac-dialogs`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-dialogs)：不再顯示核取方塊與說明按鈕、文字輸入提示、內容型別、格式彈出式選單、Finder 標籤，以及色彩與字型面板。
- [`v3/examples/mac-feedback`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-feedback)：可移除的 SF Symbol 狀態項目、觸覺回饋、系統音效、文字轉語音和語音辨識。
- [`v3/examples/mac-clipboard-drag`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-clipboard-drag)：具備變更追蹤的豐富剪貼簿功能、使用檔案承諾向外拖曳，以及文字、URL 和影像放置。
- [`v3/examples/mac-system`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-system)：權限、防止睡眠、暫停終止、電源狀態、輔助使用、鍵盤配置和地區設定，並具備即時更新。
- [`v3/examples/mac-integration`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-integration)：服務選單項目、具有接續處理常式的 Handoff 活動，以及 `app.Browser` 上的工作區輔助方法。
- [`v3/examples/mac-search-preview`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-search-preview)：具有 `OnOpen` 的 Spotlight 索引、Quick Look 預覽和縮圖，以及具有指令碼定義的自訂 Apple Event。
