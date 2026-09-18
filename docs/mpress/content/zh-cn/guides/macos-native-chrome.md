---
title: "原生 macOS 窗口界面"
description: "围绕 Wails 窗口构建原生 AppKit 工具栏、侧边栏、内容列表、检查器、附加控件和窗口标签页"
slug: "guides/macos-native-chrome"
sourcePath: "guides/macos-native-chrome.md"
---

适用平台：macOS

Wails v3 可以用真正的 AppKit 窗口界面包裹 WebView。工具栏、侧边栏、内容列表、检查器和标题栏区域都是通过 Go 创建的原生控件。其中没有 HTML，因此可以直接获得 AppKit 的材质效果、键盘处理、动画和状态持久化支持，而前端可以专注于内容。

本页的所有 API 都位于 `application` 包中，并以 `Mac` 为前缀。同一份代码也能在 Windows 和 Linux 上编译：构造函数和设置方法在所有平台上都可调用，附加操作要么不执行任何操作，要么返回错误，窗口则保留普通的单个 WebView。

## 窗口结构

一个完整配置的窗口从起始边缘到末端边缘包含以下部分：

| 部分 | 类型 | AppKit 类 |
|------|------|--------------|
| 工具栏 | `MacToolbar` | `NSToolbar` |
| 侧边栏 | `MacSidebar` | 侧边栏分割项中的 `NSOutlineView` 来源列表 |
| 内容列表 | `MacContentList` | 内容列表分割项中的 `NSTableView` |
| 主要内容 | 你的 WebView 或 `MacTextEditor` | `WKWebView` 或 `NSTextView` |
| 检查器 | `MacInspector` | 检查器分割项中的原生属性控件 |
| 附加控件 | `MacAccessory` | `NSTitlebarAccessoryViewController` 或 `NSSplitViewItemAccessoryViewController` |

这些窗格由 `MacSplitView` 排列，它是一个 `NSSplitViewController`。先构建各个部分，再按从起始边缘到末端边缘的顺序将其添加到分割视图，将分割视图附加到窗口，最后附加工具栏。可以使用 `SetSplitView` 和 `SetToolbar` 完成这些操作，也可以通过 `Mac.SplitView` 和 `Mac.Toolbar` 窗口选项一步完成。

典型的窗口配置会将这些界面元素与以下窗口选项搭配使用，使内容可以在统一的工具栏下方滚动：

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

## 工具栏

`NewMacToolbar` 创建 `NSToolbar`。使用 `Add` 方法添加项目，在返回的句柄上链式调用设置方法和回调，然后使用 `SetToolbar` 附加工具栏。标识符会自动生成。

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

项目类型包括：

- `AddButton` 添加按钮。每个按钮都必须在工具栏附加前设置 `OnClick`，否则 `SetToolbar` 会报告错误，并保留原有工具栏。
- `AddSearch` 添加 `NSSearchToolbarItem`。必须设置 `OnSearch`。可以使用 `SetSearchPlaceholder`、`SetSearchIncremental`、用于持久化最近搜索菜单的 `SetSearchRecentsKey`，以及用于设置放大镜后面的自定义菜单的 `SetSearchMenu`。
- `AddShare` 添加系统共享项目（见下文）。
- `AddGroup` 添加分段式 `NSToolbarItemGroup`。使用该组的 `AddButton` 添加成员，并选择 `ToolbarGroupSelectOne`、`ToolbarGroupMomentary` 或 `ToolbarGroupSelectAny`。
- `AddMenu` 添加由普通 `Menu` 驱动的下拉式 `NSMenuToolbarItem`。`SetShowsIndicator(false)` 可隐藏下拉箭头。
- `AddSpace` 和 `AddFlexibleSpace` 添加标准间隔项。
- `AddSidebarToggle` 和 `AddSidebarTrackingSeparator` 添加 AppKit 的侧边栏项目。分隔项会使其之前的所有项目与侧边栏分隔线对齐，因此窗口必须具有包含侧边栏窗格的分割视图。
- `AddInspectorToggle` 和 `AddInspectorTrackingSeparator` 对检查器窗格执行同样的操作。

`SetDisplayMode` 可在 `MacToolbarDisplayModeIconAndLabel`（默认值）、`MacToolbarDisplayModeIconOnly`、`MacToolbarDisplayModeLabelOnly` 和 `MacToolbarDisplayModeDefault` 之间选择。

### 实时更新

每个句柄也都可用于实时更新。附加后调用的设置方法会在应用程序线程上更新原生项目。

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

其他实用的设置方法包括 `SetLabel`、`SetTooltip`、`SetEnabled`、`SetHidden`、`SetTintColor` 和 `SetNavigational`；后者会将项目保持在起始边缘，类似 Safari 的后退和前进按钮。`SetVisibilityPriority` 决定窗口变窄时哪些项目先移入溢出菜单。

### 用户自定义

`SetCustomizable` 启用标准的“自定义工具栏...”面板，并使用你提供的键保存用户的布局。为每个项目设置稳定的 `SetPersistenceKey`，使保存的布局在应用重新启动后仍然有效；对于只应在用户添加后才出现的项目，使用 `SetInDefaultSet(false)`。请在附加工具栏前调用 `SetCustomizable`。

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

### 共享

`AddShare` 返回 `MacToolbarShareItem`。在 `MacShareProvider` 提供至少一种数据表示之前，它始终处于禁用状态。只有共享服务请求数据时，Wails 才会向提供者索取字节，因此大型导出内容会延迟生成。`MacShareProviderFunc` 可将两个函数适配为提供者；有状态的应用程序也可以直接实现该接口。

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

常见的内容类型包括 `MacShareTypePlainText`、`MacShareTypeHTML`、`MacShareTypePDF`、`MacShareTypePNG` 和 `MacShareTypeJPEG`。也接受其他任何 UTI 字符串。

### 附加和移除

一个工具栏同一时间只能属于一个窗口。`WebviewWindow.SetToolbar` 和 `NativeWindow.SetToolbar` 都可以在窗口创建前或创建后接收工具栏；传入 `nil` 会移除工具栏，并使其可供其他窗口使用。在 `WebviewWindow` 上调用 `SetToolbar` 时，验证问题通过 `Window.Error` 报告；`NativeWindow` 版本则返回错误。`Mac.Toolbar` 窗口选项会在创建窗口时附加工具栏；它在 `Mac.SplitView` 之后应用，使跟踪分隔项能找到要对齐的侧边栏。

## 分割视图

`MacSplitView` 用于排列窗格。请按从起始边缘到末端边缘的顺序添加窗格。`AddSidebar`、`AddContentList` 和 `AddInspector` 接收它们承载的原生模型，并返回用于控制尺寸和折叠状态的 `MacSplitPane`。`AddPrimaryContent` 放置窗口现有的 WebView，并返回增加了 `SetContentLayout` 方法的 `MacSplitWebviewPane`。

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

`SetSplitView` 可以在原生窗口创建前或创建后调用。提前调用时，布局会排队等待，并在窗口创建时安装。之后调用时，例如在运行中的应用程序的菜单或托盘回调中，布局会立即安装：窗口现有的 WebView 成为主窗格，当前工具栏会重新附加，使其跟踪分隔项对齐，所有排队等待的附加控件也会附加。

在 `app.Run` 之后创建的窗口可以通过 `Mac.SplitView` 和 `Mac.Toolbar` 选项一次完成配置：

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

布局规则如下：

- 布局至少需要两个窗格，且必须恰好有一个主窗格（`WebviewWindow` 使用 `AddPrimaryContent`，`NativeWindow` 使用 `AddTextEditor`）。
- 内容列表最多只能有一个，放在侧边栏之后、主窗格之前。
- 分割视图附加后，窗格结构即被冻结。窗格设置、折叠状态以及侧边栏、列表和检查器的内容仍可随时更改。
- 已安装的布局不能替换。在同一个窗口上第二次调用 `SetSplitView` 会报告 `ErrMacSplitViewAlreadyInstalled`（`WebviewWindow` 通过 `Window.Error` 报告，`NativeWindow` 则将其作为返回值），窗口保持不变。在安装前传入 `nil` 可清除待安装的布局。
- 侧边栏、列表、检查器或分割视图同一时间只能属于一个窗口。

窗格设置方法包括 `SetMinimumThickness`、`SetMaximumThickness`、`SetPreferredThicknessFraction`、`SetHoldingPriority`、`SetCollapsible`、`SetCanCollapseFromWindowResize`、`SetCollapsed`、`Toggle`、`IsCollapsed` 和 `OnCollapsedChange`。`SetAutosaveName` 可在应用启动之间持久保存分隔线位置。

主窗格上的 `SetContentLayout` 可选择 `MacContentLayoutBelowToolbar` 或 `MacContentLayoutEdgeToEdge`。`MacContentLayoutAutomatic` 继承 `MacWindow.ContentLayout`，后者又遵循 `TitleBar.FullSizeContent`。边到边布局使 AppKit 能在 macOS 26 及更高版本上，对工具栏下方的内容应用滚动边缘效果。

## 侧边栏

`MacSidebar` 是原生来源列表。它包含根行、分区和可任意深度嵌套的行。请保留返回的句柄，以便日后更新行。

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

行设置方法包括 `SetLabel`、`SetSymbol`、`SetTooltip`、`SetEnabled`、`SetHidden`、`SetBadge`、`SetAccessorySymbol`、`SetTintColor`、`SetEditable` 和 `SetExpanded`。AppKit 选中行时触发 `OnClick`；用户展开或收起其嵌套行时触发 `OnExpandedChange`；行内重命名提交后触发 `OnRename`。

### 选择

默认采用单选。`SetSelectedItem` 选中一行，但不触发其 `OnClick`。启用多选后，点击行仍会触发该行的 `OnClick`，而 `OnSelectionChange` 会报告整个选中集合。

```go
sidebar.SetAllowsMultipleSelection(true)
sidebar.OnSelectionChange(func(_ *application.Context, items []*application.MacSidebarItem) {
    for _, item := range items {
        _ = item.Section()
    }
})
```

### 上下文菜单

右键点击时，会按以下顺序查找菜单：行自身的 `SetContextMenu`，然后是侧边栏的 `OnContextMenu` 回调，最后是侧边栏的备用 `SetContextMenu`。AppKit 等待期间，回调会在应用程序线程上运行，因此应尽快完成。

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

### 拖动重排

`SetReorderable` 允许用户在分区内、分区之间以及根层级与分区之间拖动行。在触发 `OnMove` 之前，Go 模型会先更新。

```go
sidebar.SetReorderable(true)
sidebar.OnMove(func(_ *application.Context, item *application.MacSidebarItem, section *application.MacSidebarSection, index int) {
    // section is nil when the row was dropped at the sidebar root
})
```

### 移除

移除一行时，其嵌套行也会被移除。之后相应句柄不再起作用。

```go
notes.Remove(draft)
sidebar.RemoveSection(tags)
```

## 内容列表

`MacContentList` 对应 Finder、Mail 和文档浏览器中的中间列。不设置列时，它会显示信息丰富的行：标题、副标题、前置符号、后置详情和计数徽章。

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

样式包括 `MacContentListStyleAutomatic`、`MacContentListStyleInset`、`MacContentListStyleSourceList`、`MacContentListStylePlain` 和 `MacContentListStyleFullWidth`。`SetRowHeight`、`SetAlternatingRowBackgrounds`、`SetHeaderVisible` 和 `SetAllowsMultipleSelection` 提供其余显示选项。行支持 `SetHidden`、`SetEnabled`、`SetTooltip`、在指定位置调用 `InsertRow`、`Remove` 和 `RemoveAll`。

### 列

`SetColumns` 会切换到表格模式。此时每行显示其 `SetCells` 值，每列一个。标记为 `Sortable` 的列会显示排序指示符；`SetSortable` 启用点击表头排序。未设置 `OnSort` 回调时，列表会通过 `SortBy` 按列文本自行排序。

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

### 上下文菜单

上下文菜单的查找顺序与侧边栏相同：行的 `SetContextMenu`，然后是 `OnContextMenu`，最后是列表的备用菜单。

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

## 检查器

`MacInspector` 是位于末端边缘的属性面板，由按分区组织的原生控件构成。每个 `Add` 方法都会返回 `MacInspectorControl` 句柄，提供对应控件类型的设置方法和回调。

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

控件类型及其设置方法如下：

| 控件 | 设置方法 | 回调 |
|---------|---------|----------|
| `AddLabel` | `SetValue` | 无 |
| `AddTextField` | `SetValue` | `OnTextChange` |
| `AddCheckbox` | `SetChecked` | `OnToggle` |
| `AddPopup`, `AddSegmented` | `SetOptions`, `SetSelectedIndex` | `OnSelectionChange` |
| `AddSlider`, `AddStepper` | `SetFloatValue`, `SetRange`, `SetStep` | `OnValueChange` |
| `AddColorWell` | `SetColor` | `OnColorChange` |
| `AddDatePicker` | `SetDate` | `OnDateChange` |
| `AddButton` | `SetLabel` | `OnClick` |

`SetLabel`、`SetTooltip`、`SetEnabled` 和 `SetHidden` 适用于所有控件类型。分区可以折叠，分区和控件都可以随时移动或移除。

```go
appearance.SetCollapsed(true)
inspector.MoveSection(statistics, 0)
appearance.Remove(priority)
inspector.RemoveSection(statistics)
```

## 附加控件

`MacAccessory` 是一条原生控件区域。创建时需选择布局：`MacAccessoryLayoutLeading` 位于窗口按钮旁，`MacAccessoryLayoutTrailing` 位于标题栏的末端边缘，`MacAccessoryLayoutBottom` 横跨标题栏和工具栏下方的整个宽度，`MacAccessoryLayoutTop` 位于分割视图窗格的顶部。请在附加前添加所有控件。

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

控件包括 `AddSearch`、`AddSegmented`、`AddButton`、`AddSymbolButton`、`AddMenuButton`、`AddLabel`、`AddFlexibleSpace`，以及用于原生集成的 `AddNativeView`。`WebviewWindow` 和 `NativeWindow` 都提供 `AddTitlebarAccessory`；窗口创建前的调用会排队等待，并在创建时应用。`Remove` 会移除附加控件，使其可以在别处重新附加；`SetHidden` 则使其在原位折叠。

### 窗格附加控件

在 macOS 26 及更高版本上，附加控件可以位于分割窗格的顶部或底部，Finder 的侧边栏筛选字段就放在这里。使用 `Top` 或 `Bottom` 布局创建附加控件，然后在窗格上调用 `AddTopAccessory` 或 `AddBottomAccessory` 附加。

```go
filter := application.NewMacAccessory(application.MacAccessoryLayoutTop)
filter.AddSearch("Filter").
    SetIncremental(true).
    OnSearch(func(_ *application.Context, query string) { filterNotes(query) })

if err := sidebarPane.AddTopAccessory(filter); err != nil {
    window.Error("sidebar filter: %s", err)
}
```

`SetPreferredScrollEdgeEffectStyle` 为附加控件后方滚动的内容选择 AppKit 的 `Automatic`、`Soft` 或 `Hard` 效果。显式指定样式需要 macOS 26.1；在更早的版本上，请求会通过窗口的错误处理器报告，并继续使用自动样式。

## 组合使用

这个程序构建一个三窗格窗口，包含侧边栏、WebView 和检查器，以及跟踪两条分隔线的工具栏。

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

原生窗口界面与前端通过常规的 Wails 事件和服务通信。侧边栏发出 `note:selected`，检查器发出 `note:title`，页面通过运行时的 `Events.On` 监听。

## 原生窗口

@note{type="caution" title="实验性功能"}
`NativeWindow`、`NativeWindowManager` 和 `MacTextEditor` 在 v3 中是实验性功能。API 特意保持精简，并可能在 v4 重新设计通用窗口 API 时发生变化。
@end

`NativeWindow` 没有 WebView。它的主要内容是 `MacTextEditor`，即位于 `NSScrollView` 内的 `NSTextView`；它接受与 `WebviewWindow` 相同的工具栏、分割视图和附加控件类型。使用 `app.NativeWindow.New` 或 `app.NativeWindow.NewWithOptions` 创建窗口，之后可使用 `Get` 或 `GetByID` 查找。

原生窗口只有在具备内容后才会创建。请通过 `NativeWindowOptions.SplitView` 或调用 `SetSplitView`，提供一个使用 `AddTextEditor` 添加了主窗格的分割视图。在 `app.Run` 之前，布局会排队等待；在运行中的应用程序内，`SetSplitView` 会立即创建并显示窗口（除非设置了 `Hidden`），并返回任何创建错误。没有布局的窗口会继续延迟创建，`Run` 返回 `ErrNativeWindowContentRequired`；没有文本编辑器的布局会被拒绝，并返回 `ErrNativeWindowEditorRequired`。请使用 `NativeWindowOptions.Toolbar` 和 `NativeWindowOptions.SplitView`，不要使用原生窗口会忽略的 `Mac.Toolbar` 和 `Mac.SplitView` 字段。

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

在运行中的应用程序内，也可以将窗口界面作为选项传入，一次调用创建相同的窗口：

```go
window := app.NativeWindow.NewWithOptions(application.NativeWindowOptions{
    Title:     "Native Notes",
    SplitView: split,
    Toolbar:   toolbar,
})
```

`MacTextEditor` 提供 `SetText`、`Text`、`SetEditable`、`OnChange` 和 `Focus`。通过程序调用 `SetText` 不会触发 `OnChange`，因此加载文件不会将其标记为已修改。`Text` 会从 AppKit 读取完整文档，因此应在需要内容时调用，而不是每次发生更改时都调用。

以下两个选项可使纯原生应用程序保持精简：

- `application.Options` 中的 `NativeOnly: true` 会在运行时跳过前端传输和资源服务器。设置后不要创建 `WebviewWindow`。
- `wails_native` 构建标签会将 WebView、前端和更新程序代码从二进制文件中完全排除，并自动设置 `NativeOnly`：

```sh
go build -tags wails_native .
```

`wails_native` 构建也会排除单实例支持；如需该功能，请添加 `wails_single_instance` 标签。

## 窗口标签页

macOS 可以将多个窗口组合为标签页。标签页模式在窗口创建时确定，因此应为每个需要参与的窗口将 `Mac.TabbingMode` 设为 `MacWindowTabbingModePreferred` 或 `MacWindowTabbingModeAutomatic`。`MacWindowTabbingModeDisallowed`（以及未设置时的默认值）会使窗口不加入标签页组。

```go
first := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Document 1",
    URL:   "/",
    Mac: application.MacWindow{
        TabbingMode: application.MacWindowTabbingModePreferred,
    },
})
```

标签页操作需要已创建的原生窗口，因此请在 `app.Run` 启动后运行的菜单处理器、服务方法或其他代码中调用。`TabGroup` 返回 `MacWindowTabGroup` 句柄，每次调用都会重新解析原生组；使用 nil 句柄是安全的，其所有方法都会返回对应的零值。

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

`MacWindowTabGroup` 提供 `Identifier`、`Count`、`Windows`、`NativeWindows`、`SelectedWindow`、`Select`、`SelectNative`、`SelectNext`、`SelectPrevious`、`IsTabBarVisible`、`ToggleTabBar`、`IsOverviewVisible` 和 `ToggleTabOverview`。`app.Window.TabGroups` 列出包含 `WebviewWindow` 的所有组。`AddNativeTab` 将 `NativeWindow` 添加到 WebView 窗口的组，`SetTabTooltip` 设置标签页的悬停文字。在 macOS 以外的平台上，标签页方法返回 `ErrMacWindowTabsUnsupported`。

## 版本要求

本页所有功能仅适用于 macOS。在较早的版本上，功能会按下表所述平稳降级；Go API 在所有平台上相同。

| 功能 | 最低 macOS 版本 | 在更早版本上的行为 |
|---------|---------------|-------------------------------|
| 窗口标签页 | 10.12 | 不可用 |
| 工具栏组、带边框的项目、`AddMenu` | 10.15 | 菜单项目会省略；组使用旧版显示方式 |
| SF Symbols（对工具栏、侧边栏、列表和附加控件项目调用 `SetSymbol`） | 11 | 不显示图像 |
| 作为 `NSSearchToolbarItem` 的 `AddSearch`、`SetNavigational`、侧边栏跟踪分隔项 | 11 | 搜索退回为普通搜索字段；分隔项会省略 |
| 检查器和内容列表分割角色 | 11 | 相同的窗格由普通分割项承载 |
| 内容列表样式 | 11 | 忽略 |
| `SetCenteredItems` | 13 | 忽略 |
| 检查器切换按钮和跟踪分隔项 | 14 | Wails 提供原生切换按钮；分隔项会省略 |
| 对工具栏项目调用 `SetBadgeCount`、`SetProminent`、`SetTintColor` | 26 | 保存设置，并在可用时应用 |
| 窗格附加控件（`AddTopAccessory`、`AddBottomAccessory`） | 26 | `ErrMacSplitItemAccessoryUnavailable` |
| 显式指定滚动边缘效果样式 | 26.1 | 通过窗口错误处理器报告；继续使用自动样式 |

标题栏附加控件、分割视图、侧边栏、检查器和文本编辑器除 Wails 的最低版本要求外，没有其他版本要求。

## 示例

每个示例都是完整且可运行的应用程序：

- [`v3/examples/mac-toolbar`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-toolbar)：带有工具栏、侧边栏、内容列表、检查器、共享提供者、工具栏自定义功能和窗格附加控件的笔记编辑器。
- [`v3/examples/mac-titlebar-accessory`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-titlebar-accessory)：通过位于起始边缘、末端边缘和底部的标题栏附加控件操作邮箱视图。
- [`v3/examples/mac-window-tabs`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-window-tabs)：从菜单操作标签页模式、`AddTab`、`AddNativeTab` 和标签页组 API。
- [`v3/examples/mac-native-editor`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-native-editor)：使用 `wails_native` 标签构建的不含 WebView 的 `NativeWindow` 文本编辑器。

## 深入了解

本页介绍的窗口界面只是原生 macOS 应用程序的一部分。另一部分是应用程序的行为：带有代理图标和级联排列的文档窗口、带有符号和徽章的菜单、显示进度的 Dock 菜单、原生警告和面板、可移除的状态项、触觉反馈和语音、支持丰富内容的剪贴板和向外拖动，以及有关权限、电源和区域设置的系统信息。这些 API 在 [macOS 平台集成](/guides/macos-platform-integration)指南中介绍。
