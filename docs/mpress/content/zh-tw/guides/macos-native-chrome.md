---
title: "原生 macOS 視窗介面"
description: "在 Wails 視窗周圍建構原生 AppKit 工具列、側邊欄、內容清單、檢閱器、附加控制列與視窗分頁"
slug: "guides/macos-native-chrome"
sourcePath: "guides/macos-native-chrome.md"
---

適用平台：macOS

Wails v3 可以用真正的 AppKit 視窗介面包覆 WebView。工具列、側邊欄、內容清單、檢閱器與標題列控制列都是從 Go 建立的原生控制項，不含任何 HTML，因此能直接使用 AppKit 的材質、鍵盤操作、動畫與持久化功能，讓前端專注於內容。

本頁所有 API 都位於 `application` 套件，並以 `Mac` 為前綴。相同程式碼也能在 Windows 和 Linux 上編譯：建構函式與設定方法在所有平台都可運作，附加操作不執行任何動作或傳回錯誤，而視窗會保有一般的單一 WebView。

## 視窗結構

完整配置的視窗從起始邊緣到結束邊緣包含以下部分：

| 部分 | 型別 | AppKit 類別 |
|------|------|--------------|
| 工具列 | `MacToolbar` | `NSToolbar` |
| 側邊欄 | `MacSidebar` | 位於側邊欄分割項目中的 `NSOutlineView` 來源清單 |
| 內容清單 | `MacContentList` | 位於內容清單分割項目中的 `NSTableView` |
| 主要內容 | 您的 WebView，或 `MacTextEditor` | `WKWebView` 或 `NSTextView` |
| 檢閱器 | `MacInspector` | 位於檢閱器分割項目中的原生屬性控制項 |
| 附加控制列 | `MacAccessory` | `NSTitlebarAccessoryViewController` 或 `NSSplitViewItemAccessoryViewController` |

各窗格由 `MacSplitView` 排列，它是一個 `NSSplitViewController`。先建立各部分，依起始邊緣到結束邊緣的順序加入分割檢視，將分割檢視附加至視窗，然後附加工具列。您可以使用 `SetSplitView` 和 `SetToolbar` 完成，也可以透過 `Mac.SplitView` 和 `Mac.Toolbar` 視窗選項一次完成。

典型的視窗設定會搭配下列視窗選項，讓內容能在整合式工具列下方捲動：

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

## 工具列

`NewMacToolbar` 會建立 `NSToolbar`。使用 `Add` 方法加入項目，在傳回的控制代碼上鏈結設定方法與回呼，並使用 `SetToolbar` 附加工具列。識別碼會自動產生。

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

項目種類如下：

- `AddButton` 加入按鈕。每個按鈕都必須在工具列附加前設定 `OnClick`，否則 `SetToolbar` 會回報錯誤並保留先前的工具列。
- `AddSearch` 加入 `NSSearchToolbarItem`。必須設定 `OnSearch`。使用 `SetSearchPlaceholder`、`SetSearchIncremental`、用於保存最近搜尋選單的 `SetSearchRecentsKey`，以及用於放大鏡後方自訂選單的 `SetSearchMenu`。
- `AddShare` 加入系統分享項目（見下文）。
- `AddGroup` 加入分段式 `NSToolbarItemGroup`。使用群組的 `AddButton` 加入成員，並選擇 `ToolbarGroupSelectOne`、`ToolbarGroupMomentary` 或 `ToolbarGroupSelectAny`。
- `AddMenu` 加入由一般 `Menu` 驅動的下拉式 `NSMenuToolbarItem`。`SetShowsIndicator(false)` 會隱藏箭頭。
- `AddSpace` 和 `AddFlexibleSpace` 加入標準間隔項目。
- `AddSidebarToggle` 和 `AddSidebarTrackingSeparator` 加入 AppKit 的側邊欄項目。分隔線會讓其前方所有項目對齊側邊欄分隔線上方，因此視窗必須有包含側邊欄窗格的分割檢視。
- `AddInspectorToggle` 和 `AddInspectorTrackingSeparator` 對檢閱器窗格執行相同操作。

`SetDisplayMode` 可選擇 `MacToolbarDisplayModeIconAndLabel`（預設值）、`MacToolbarDisplayModeIconOnly`、`MacToolbarDisplayModeLabelOnly` 或 `MacToolbarDisplayModeDefault`。

### 即時更新

每個控制代碼也都可用於即時更新。附加後呼叫的設定方法會在應用程式執行緒上更新原生項目。

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

其他實用的設定方法包括 `SetLabel`、`SetTooltip`、`SetEnabled`、`SetHidden`、`SetTintColor` 和 `SetNavigational`；最後一項會將項目保留在起始邊緣，方式與 Safari 保留上一頁和下一頁按鈕相同。`SetVisibilityPriority` 決定視窗變窄時哪些項目會先移入溢出選單。

### 使用者自訂

`SetCustomizable` 會啟用標準的「自訂工具列...」面板，並以您提供的鍵儲存使用者版面配置。請為每個項目指定穩定的 `SetPersistenceKey`，讓儲存的版面配置在重新啟動後仍有效；對於應只在使用者加入後才顯示的項目，使用 `SetInDefaultSet(false)`。請在附加工具列前呼叫 `SetCustomizable`。

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

### 分享

`AddShare` 會傳回 `MacToolbarShareItem`。在 `MacShareProvider` 宣告至少一種資料表示形式之前，它會保持停用。Wails 只有在分享服務要求資料時才向提供者索取位元組，因此大型匯出內容會延遲產生。`MacShareProviderFunc` 可將兩個函式轉接為提供者；具備狀態的應用程式可以直接實作該介面。

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

常見內容型別有 `MacShareTypePlainText`、`MacShareTypeHTML`、`MacShareTypePDF`、`MacShareTypePNG` 和 `MacShareTypeJPEG`。也接受任何其他 UTI 字串。

### 附加與卸除

工具列一次只能屬於一個視窗。`WebviewWindow.SetToolbar` 和 `NativeWindow.SetToolbar` 都可在視窗建立前或建立後接收工具列；傳入 `nil` 會移除工具列，並釋出供其他地方使用。在 `WebviewWindow` 上呼叫 `SetToolbar` 會透過 `Window.Error` 回報驗證問題，而 `NativeWindow` 版本會傳回錯誤。`Mac.Toolbar` 視窗選項會在建立時附加工具列；它會在 `Mac.SplitView` 之後套用，讓追蹤分隔線能找到要對齊的側邊欄。

## 分割檢視

`MacSplitView` 負責排列窗格。請依起始邊緣到結束邊緣的順序加入。`AddSidebar`、`AddContentList` 和 `AddInspector` 接收其承載的原生模型，並傳回用於控制尺寸與收合的 `MacSplitPane`。`AddPrimaryContent` 會放入視窗現有的 WebView，並傳回額外提供 `SetContentLayout` 的 `MacSplitWebviewPane`。

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

`SetSplitView` 可在原生視窗建立前或建立後使用。若在建立前呼叫，版面配置會排入佇列，待視窗建立時安裝。若在建立後呼叫，例如從執行中應用程式的選單或系統匣回呼呼叫，版面配置會立即安裝：視窗現有的 WebView 會成為主要窗格，目前的工具列會重新附加，使追蹤分隔線對齊，任何已排入佇列的附加控制列也會附加。

在 `app.Run` 之後建立的視窗，可以透過 `Mac.SplitView` 和 `Mac.Toolbar` 選項一次設定：

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

版面配置須遵守以下規則：

- 版面配置至少需要兩個窗格，且恰好有一個主要窗格（`WebviewWindow` 使用 `AddPrimaryContent`，`NativeWindow` 使用 `AddTextEditor`）。
- 最多只能有一個內容清單，放在側邊欄之後、主要窗格之前。
- 分割檢視附加後，窗格結構便固定。窗格設定、收合狀態，以及側邊欄、清單和檢閱器的內容，仍可隨時變更。
- 已安裝的版面配置不能更換。對同一視窗再次呼叫 `SetSplitView` 會回報 `ErrMacSplitViewAlreadyInstalled`（`WebviewWindow` 透過 `Window.Error`，`NativeWindow` 則作為傳回值），且視窗保持不變。在安裝前傳入 `nil` 可清除待安裝的版面配置。
- 側邊欄、清單、檢閱器或分割檢視一次只能屬於一個視窗。

窗格設定方法有 `SetMinimumThickness`、`SetMaximumThickness`、`SetPreferredThicknessFraction`、`SetHoldingPriority`、`SetCollapsible`、`SetCanCollapseFromWindowResize`、`SetCollapsed`、`Toggle`、`IsCollapsed` 和 `OnCollapsedChange`。`SetAutosaveName` 會跨啟動儲存分隔線位置。

主要窗格上的 `SetContentLayout` 可選擇 `MacContentLayoutBelowToolbar` 或 `MacContentLayoutEdgeToEdge`。`MacContentLayoutAutomatic` 會繼承 `MacWindow.ContentLayout`，後者則遵循 `TitleBar.FullSizeContent`。在 macOS 26 及更新版本中，全邊緣版面配置可讓 AppKit 在工具列下方套用捲動邊緣效果。

## 側邊欄

`MacSidebar` 是原生來源清單，可容納根層列、區段，以及任意深度巢狀排列的列。請保留傳回的控制代碼，以便稍後更新列。

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

列的設定方法有 `SetLabel`、`SetSymbol`、`SetTooltip`、`SetEnabled`、`SetHidden`、`SetBadge`、`SetAccessorySymbol`、`SetTintColor`、`SetEditable` 和 `SetExpanded`。AppKit 選取列時會觸發 `OnClick`；使用者展開或收合其巢狀列時會觸發 `OnExpandedChange`；行內重新命名確認後會觸發 `OnRename`。

### 選取

預設為單選。`SetSelectedItem` 會選取一列，但不觸發其 `OnClick`。啟用多選後，點選的列仍會觸發 `OnClick`，而 `OnSelectionChange` 會回報完整的選取集合。

```go
sidebar.SetAllowsMultipleSelection(true)
sidebar.OnSelectionChange(func(_ *application.Context, items []*application.MacSidebarItem) {
    for _, item := range items {
        _ = item.Section()
    }
})
```

### 快顯選單

按右鍵時會依下列順序尋找選單：列本身的 `SetContextMenu`、側邊欄的 `OnContextMenu` 回呼，最後是側邊欄的備用 `SetContextMenu`。AppKit 等待期間，回呼會在應用程式執行緒上執行，因此應盡快完成。

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

### 拖曳重新排序

`SetReorderable` 讓使用者可在區段內、區段之間，以及根層與區段之間拖曳列。Go 模型會在觸發 `OnMove` 前更新。

```go
sidebar.SetReorderable(true)
sidebar.OnMove(func(_ *application.Context, item *application.MacSidebarItem, section *application.MacSidebarSection, index int) {
    // section is nil when the row was dropped at the sidebar root
})
```

### 移除

移除一列也會移除其巢狀列。之後這些控制代碼不再生效。

```go
notes.Remove(draft)
sidebar.RemoveSection(tags)
```

## 內容清單

`MacContentList` 是 Finder、郵件和文件瀏覽器的中間欄。沒有欄位時，它會顯示資訊豐富的列：標題、副標題、起始圖示、結尾詳細資訊和數量徽章。

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

樣式有 `MacContentListStyleAutomatic`、`MacContentListStyleInset`、`MacContentListStyleSourceList`、`MacContentListStylePlain` 和 `MacContentListStyleFullWidth`。`SetRowHeight`、`SetAlternatingRowBackgrounds`、`SetHeaderVisible` 和 `SetAllowsMultipleSelection` 提供其餘顯示選項。列支援 `SetHidden`、`SetEnabled`、`SetTooltip`、在指定位置使用 `InsertRow`，以及 `Remove` 和 `RemoveAll`。

### 欄位

`SetColumns` 會切換至表格模式。此時每列會顯示其 `SetCells` 值，每欄一個。標記為 `Sortable` 的欄位會顯示排序指示器；`SetSortable` 會啟用點選表頭排序。若沒有 `OnSort` 回呼，清單會透過 `SortBy` 依欄位文字自行排序。

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

### 快顯選單

快顯選單的選擇順序與側邊欄相同：列的 `SetContextMenu`、`OnContextMenu`，最後是清單的備用選單。

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

## 檢閱器

`MacInspector` 是由原生控制項組成、依區段分組的結尾側屬性面板。每個 `Add` 方法都會傳回 `MacInspectorControl` 控制代碼，提供該控制項種類專用的設定方法與回呼。

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

控制項種類及其設定方法如下：

| 控制項 | 設定方法 | 回呼 |
|---------|---------|----------|
| `AddLabel` | `SetValue` | 無 |
| `AddTextField` | `SetValue` | `OnTextChange` |
| `AddCheckbox` | `SetChecked` | `OnToggle` |
| `AddPopup`, `AddSegmented` | `SetOptions`, `SetSelectedIndex` | `OnSelectionChange` |
| `AddSlider`, `AddStepper` | `SetFloatValue`, `SetRange`, `SetStep` | `OnValueChange` |
| `AddColorWell` | `SetColor` | `OnColorChange` |
| `AddDatePicker` | `SetDate` | `OnDateChange` |
| `AddButton` | `SetLabel` | `OnClick` |

`SetLabel`、`SetTooltip`、`SetEnabled` 和 `SetHidden` 適用於所有種類。區段可以收合，區段與控制項也都可隨時移動或移除。

```go
appearance.SetCollapsed(true)
inspector.MoveSection(statistics, 0)
appearance.Remove(priority)
inspector.RemoveSection(statistics)
```

## 附加控制列

`MacAccessory` 是一列原生控制項。建立時請選擇版面配置：`MacAccessoryLayoutLeading` 位於視窗按鈕旁，`MacAccessoryLayoutTrailing` 位於標題列的結尾邊緣，`MacAccessoryLayoutBottom` 橫跨標題列與工具列下方的整個寬度，而 `MacAccessoryLayoutTop` 位於分割檢視窗格頂部。請在附加控制列前加入所有控制項。

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

控制項包括 `AddSearch`、`AddSegmented`、`AddButton`、`AddSymbolButton`、`AddMenuButton`、`AddLabel`、`AddFlexibleSpace`，以及供原生整合使用的 `AddNativeView`。`WebviewWindow` 和 `NativeWindow` 都提供 `AddTitlebarAccessory`；在視窗建立前呼叫時，操作會排入佇列，並在建立視窗時套用。`Remove` 會卸除附加控制列，供其他地方再次附加；`SetHidden` 則會就地將其收合。

### 窗格附加控制列

在 macOS 26 及更新版本中，附加控制列可位於分割窗格的頂部或底部，Finder 的側邊欄篩選欄位便位於此處。使用 `Top` 或 `Bottom` 版面配置建立附加控制列，並透過窗格上的 `AddTopAccessory` 或 `AddBottomAccessory` 附加。

```go
filter := application.NewMacAccessory(application.MacAccessoryLayoutTop)
filter.AddSearch("Filter").
    SetIncremental(true).
    OnSearch(func(_ *application.Context, query string) { filterNotes(query) })

if err := sidebarPane.AddTopAccessory(filter); err != nil {
    window.Error("sidebar filter: %s", err)
}
```

`SetPreferredScrollEdgeEffectStyle` 可為附加控制列後方的內容捲動選擇 AppKit 的 `Automatic`、`Soft` 或 `Hard` 效果。明確指定樣式需要 macOS 26.1；在較舊版本中，此要求會透過視窗的錯誤處理常式回報，並繼續使用自動樣式。

## 組合使用

此程式會建立包含側邊欄、WebView 和檢閱器的三窗格視窗，以及追蹤兩條分隔線的工具列。

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

原生視窗介面與前端透過一般的 Wails 事件和服務通訊。側邊欄發出 `note:selected`，檢閱器發出 `note:title`，頁面則透過執行階段的 `Events.On` 監聽。

## 原生視窗

@note{type="caution" title="實驗性功能"}
`NativeWindow`、`NativeWindowManager` 和 `MacTextEditor` 在 v3 中屬於實驗性功能。此 API 刻意保持精簡，並可能在 v4 重新設計共用視窗 API 時變更。
@end

`NativeWindow` 沒有 WebView。其主要內容是 `MacTextEditor`，也就是位於 `NSScrollView` 中的 `NSTextView`；它接受與 `WebviewWindow` 相同的工具列、分割檢視和附加控制列型別。使用 `app.NativeWindow.New` 或 `app.NativeWindow.NewWithOptions` 建立，之後可透過 `Get` 或 `GetByID` 查找。

原生視窗只有在具備內容後才會建立。請提供以 `AddTextEditor` 加入主要窗格的分割檢視，方式可以是透過 `NativeWindowOptions.SplitView`，也可以是呼叫 `SetSplitView`。在 `app.Run` 前，版面配置會排入佇列；在執行中的應用程式內，`SetSplitView` 會立即建立並顯示視窗（除非設定了 `Hidden`），並傳回任何建立錯誤。沒有版面配置的視窗會保持延後建立，而 `Run` 會傳回 `ErrNativeWindowContentRequired`；沒有文字編輯器的版面配置則會以 `ErrNativeWindowEditorRequired` 拒絕。請使用 `NativeWindowOptions.Toolbar` 和 `NativeWindowOptions.SplitView`，不要使用原生視窗會忽略的 `Mac.Toolbar` 和 `Mac.SplitView` 欄位。

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

在執行中的應用程式內，也可以將視窗介面作為選項傳入，一次建立相同的視窗：

```go
window := app.NativeWindow.NewWithOptions(application.NativeWindowOptions{
    Title:     "Native Notes",
    SplitView: split,
    Toolbar:   toolbar,
})
```

`MacTextEditor` 提供 `SetText`、`Text`、`SetEditable`、`OnChange` 和 `Focus`。程式呼叫 `SetText` 不會觸發 `OnChange`，因此載入檔案不會將其標記為已修改。`Text` 會從 AppKit 讀取完整文件，因此請在需要內容時呼叫，而非每次變更時都呼叫。

以下兩個選項可讓純原生應用程式保持精簡：

- 在 `application.Options` 中設定 `NativeOnly: true`，會在執行時略過前端傳輸層和資源伺服器。設定後請勿建立 `WebviewWindow`。
- `wails_native` 建置標籤會將 WebView、前端和更新程式碼完全排除在二進位檔之外，並自動設定 `NativeOnly`：

```sh
go build -tags wails_native .
```

`wails_native` 建置也會排除單一執行個體支援；需要此功能時，請加入 `wails_single_instance` 標籤。

## 視窗分頁

macOS 可以將視窗群組成分頁。分頁模式在建立視窗時便固定，因此每個要參與的視窗都應將 `Mac.TabbingMode` 設為 `MacWindowTabbingModePreferred` 或 `MacWindowTabbingModeAutomatic`。`MacWindowTabbingModeDisallowed`（以及未設定時的預設值）會讓視窗不加入分頁群組。

```go
first := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Document 1",
    URL:   "/",
    Mac: application.MacWindow{
        TabbingMode: application.MacWindowTabbingModePreferred,
    },
})
```

分頁操作需要已建立的原生視窗，因此請從選單處理常式、服務方法，或任何在 `app.Run` 開始後執行的程式碼呼叫。`TabGroup` 傳回 `MacWindowTabGroup` 控制代碼，每次呼叫都會重新取得原生群組；nil 控制代碼可安全使用，其每個方法都會傳回對應的零值。

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

`MacWindowTabGroup` 提供 `Identifier`、`Count`、`Windows`、`NativeWindows`、`SelectedWindow`、`Select`、`SelectNative`、`SelectNext`、`SelectPrevious`、`IsTabBarVisible`、`ToggleTabBar`、`IsOverviewVisible` 和 `ToggleTabOverview`。`app.Window.TabGroups` 會列出包含 `WebviewWindow` 的所有群組。`AddNativeTab` 會將 `NativeWindow` 加入 WebView 視窗的群組，而 `SetTabTooltip` 會設定分頁的懸停文字。在 macOS 以外的平台上，分頁方法會傳回 `ErrMacWindowTabsUnsupported`。

## 版本需求

本頁所有功能都僅適用於 macOS。在較舊版本上，功能會依下表所述適度降級；Go API 在所有平台上都相同。

| 功能 | 最低 macOS 版本 | 較舊版本的行為 |
|---------|---------------|-------------------------------|
| 視窗分頁 | 10.12 | 不可用 |
| 工具列群組、有邊框項目、`AddMenu` | 10.15 | 省略選單項目；群組使用舊版呈現方式 |
| SF Symbols（工具列、側邊欄、清單和附加控制列項目上的 `SetSymbol`） | 11 | 不顯示影像 |
| 作為 `NSSearchToolbarItem` 的 `AddSearch`、`SetNavigational`、側邊欄追蹤分隔線 | 11 | 搜尋功能退回一般搜尋欄位；省略分隔線 |
| 檢閱器和內容清單分割角色 | 11 | 相同窗格由一般分割項目承載 |
| 內容清單樣式 | 11 | 忽略 |
| `SetCenteredItems` | 13 | 忽略 |
| 檢閱器切換按鈕和追蹤分隔線 | 14 | Wails 提供原生切換按鈕；省略分隔線 |
| 工具列項目上的 `SetBadgeCount`、`SetProminent`、`SetTintColor` | 26 | 儲存設定，並在可用時套用 |
| 窗格附加控制列（`AddTopAccessory`、`AddBottomAccessory`） | 26 | `ErrMacSplitItemAccessoryUnavailable` |
| 明確指定的捲動邊緣效果樣式 | 26.1 | 透過視窗錯誤處理常式回報；維持自動樣式 |

標題列附加控制列、分割檢視、側邊欄、檢閱器和文字編輯器，除了 Wails 的最低版本需求外，沒有其他版本需求。

## 範例

每個範例都是完整且可執行的應用程式：

- [`v3/examples/mac-toolbar`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-toolbar)：具有工具列、側邊欄、內容清單、檢閱器、分享提供者、工具列自訂功能和窗格附加控制列的筆記編輯器。
- [`v3/examples/mac-titlebar-accessory`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-titlebar-accessory)：由起始、結尾及底部標題列附加控制列操作的信箱檢視。
- [`v3/examples/mac-window-tabs`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-window-tabs)：透過選單操作分頁模式、`AddTab`、`AddNativeTab` 和分頁群組 API。
- [`v3/examples/mac-native-editor`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-native-editor)：使用 `wails_native` 標籤建置、不含 WebView 的 `NativeWindow` 文字編輯器。

## 進一步探索

本頁的視窗介面只是原生 macOS 應用程式的一部分。另一部分是應用程式的行為：具有代理圖示與階梯式排列的文件視窗、帶有圖示與徽章的選單、顯示進度的 Dock 選單、原生警示和面板、可移除的狀態項目、觸覺回饋與語音、豐富的剪貼簿功能與向外拖曳，以及權限、電源和地區設定等系統資訊。這些 API 都收錄於 [macOS 平台整合](/guides/macos-platform-integration)指南。
