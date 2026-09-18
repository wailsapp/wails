---
title: "macOS 平台集成"
description: "macOS 上的文档窗口、Dock 和菜单扩展、原生面板、状态项、反馈、支持丰富内容的剪贴板、向外拖动、权限、电源和生命周期"
slug: "guides/macos-platform-integration"
sourcePath: "guides/macos-platform-integration.md"
---

适用平台：macOS

Wails v3 让应用程序具备 macOS 用户所期望的原生行为：带有代理图标和级联排列的文档窗口、带有符号和徽章的菜单、显示进度的 Dock 菜单、原生警告和面板、可移除的状态项、触觉反馈和语音、支持向外拖动的丰富内容剪贴板、有关权限、电源和区域设置的系统信息，以及与服务菜单、Handoff、AppleScript 和 Quick Look 的集成。所有功能都通过 `application` 包由 Go 驱动。

同一份代码也能在 Windows 和 Linux 上编译。设置方法会保存其值，查询返回零值，需要 macOS 的操作会返回已记录的错误，例如 `ErrMacOnly`、`ErrDialogNotSupported` 或 `ErrClipboardNotSupported`。下方的[平台说明](#platform-notes)列出了各功能在 macOS 以外平台上的行为。

有关原生窗口界面（工具栏、侧边栏、检查器、附加控件和窗口标签页），请参阅[原生 macOS 窗口界面](/guides/macos-native-chrome)指南。

## 文档窗口

文档窗口会在标题栏中显示它所代表的文件，在关闭按钮上用圆点标记未保存的更改，并以级联方式打开新窗口。这些功能都由 `WebviewWindow` 上的方法和 `MacWindow` 上的选项提供。

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

- `SetRepresentedFile` 在标题栏中显示文件的代理图标。用户可以将图标拖到其他应用程序，或按住 Command 键点击图标查看路径。传入 `""` 可将其移除。`RepresentedFile` 可读取该值。
- `SetDocumentEdited` 在关闭按钮上显示表示未保存更改的圆点，并使代理图标变暗。`IsDocumentEdited` 可读取该状态。
- `SetSubtitle` 在 macOS 11 及更高版本上于标题下方显示第二行。
- `InitialPosition: application.WindowCascade` 将窗口放在上一个级联窗口的右下方，与打开新文档时相同。`CascadeFrom(other)` 对现有窗口执行同样的操作，并更新后续窗口的级联位置。
- `Mac.FrameAutosaveName` 在窗口首次显示前恢复保存的位置和尺寸，并在窗口移动时持续保存。恢复的窗口边框设置优先于 `X`、`Y`、`Width`、`Height` 和 `InitialPosition`。`SetFrameAutosaveName` 可切换已创建窗口的保存名称。
- `MacTitleBar.WindowButtonsOffset` 按指定的点数移动关闭、最小化和缩放按钮。`SetWindowButtonsOffset` 和 `ResetWindowButtonsOffset` 可在运行时更改它。

这三个设置方法都可以在原生窗口创建前调用；创建窗口时会应用其值。

### 请求用户关注

应用程序在后台运行时，`RequestAttention` 会使 Dock 图标弹跳。信息类请求弹跳一次。紧急请求会持续弹跳，直到用户激活应用程序或你取消请求。

```go
request := window.RequestAttention(true)

// Once the work that needed attention is done:
request.Cancel()
```

`Flash` 仍是跨平台请求用户关注一次的方法。

### 打印和导出

`PrintWithOptions` 使用明确指定的页面设置打印 WebView。传入零值会显示使用共享打印设置的打印面板。`Print` 保留原有行为（横向、30 点页边距）。

`ExportPDF` 将页面渲染为 PDF 文档，`Snapshot` 则将其捕获为 PNG。两者都会等待 WebKit，因此应从 goroutine 中调用，绝不要在应用程序线程上调用；在该线程上调用会返回 `ErrMacExportOnMainThread`。

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

`PrintOptions` 还接受 `PrinterName`、`PaperName`（例如 `"iso-a4"` 这样的 PostScript 名称）和 `Scale`。`PDFExportOptions` 和 `SnapshotOptions` 接受可选的 `Rect` 来限定捕获范围，以及默认值为 `DefaultMacExportTimeout`（30 秒）的 `Timeout`。

## 工作表

工作表是附加在父窗口顶部的第二个窗口，类似“保存”面板。任何 `WebviewWindow` 都可以通过 `PresentSheet` 作为另一个窗口的工作表显示，并通过 `EndSheet` 结束；传入的响应代码会送达工作表的 `OnSheetEnd` 回调。创建工作表窗口时应设置 `Hidden`，避免它在附加前短暂显示；结束时 AppKit 会再次将其从屏幕上移走，因此同一个窗口可以重复显示为工作表。

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

`PresentCriticalSheet` 会将工作表显示在已附加的普通工作表之前，而不是排队等待。`PresentNativeSheet` 对 `NativeWindow` 执行同样的操作。`IsSheet`、`SheetParent`、`AttachedSheet`、`AttachedNativeSheet` 和 `HasAttachedSheet` 描述当前状态。关闭工作表窗口而不是结束它，会产生 `MacSheetResponseStop`。

## 弹出面板

`MacPopover` 是一个 `NSPopover`：锚定在窗口中的矩形、工具栏项目或菜单栏中的状态项上的临时面板。其内容是原生 `MacAccessory` 控件区域，与[原生 macOS 窗口界面](/guides/macos-native-chrome)指南用于标题栏附加控件的类型相同。请在首次显示前添加所有控件。

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

`MacToolbarItem.ShowPopover` 和 `SystemTray.ShowPopover` 可将同一个弹出面板锚定到工具栏项目或状态项。`MacPopoverBehaviorTransient` 会在点击面板外任意位置时关闭面板，`Semitransient` 只会在点击呈现它的窗口时关闭；默认行为则会保持打开，直到调用 `Close`。`MacRectEdge` 选择弹出面板出现的一侧。`SetContentSize` 和 `SetBehavior` 可调整正在显示的弹出面板，`Destroy` 会释放原生弹出面板，并使内容区域可在别处使用。

## 状态恢复

macOS 会在崩溃、强制退出或重启后重新打开应用程序窗口；如果系统设置中的“退出应用程序时关闭窗口”处于关闭状态，正常退出后也会重新打开。为窗口指定 `Mac.RestorationID`，使用 `SetRestorationData` 保存重建窗口所需的信息，并注册 `app.Window.OnRestore`，以便下次启动时重新构建窗口。

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

只有应用程序终止时可见的窗口才会被保存。`RestorationState.Data` 是字符串映射；其中应只保存标识符、路径和位置。`InteractionState` 将 WebView 的前进后退列表和滚动位置作为不透明数据块返回（macOS 12+），`RestoreInteractionState` 可将其应用到重建的窗口；通常会将其编码为 base64 并存入恢复数据。`SetRestorationID` 和 `RestorationID` 可更改和读取已创建窗口的标识符。

## 呈现选项

`MacPresentationOptions` 对应 `NSApplication.presentationOptions`：一个位掩码，可在应用程序处于活动状态时隐藏 Dock 或菜单栏，并禁止切换进程、强制退出、注销或使用“隐藏”命令。在应用程序选项中设置 `Mac.PresentationOptions` 可在启动时应用，也可以使用 `SetPresentationOptions` 在运行时更改。无效的组合会在到达 AppKit 前被拒绝，并返回封装了 `ErrMacPresentationOptionsInvalid` 的错误。

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

隐藏菜单栏（`HideMenuBar` 或 `AutoHideMenuBar`）需要设置一个 Dock 选项，而 `AutoHideToolbar` 同时需要 `FullScreen` 和 `AutoHideMenuBar`。`Validate` 报告值违反的第一条规则，`Has` 检查单个标志。

## 菜单和 Dock

菜单项支持 SF Symbols、徽章、分区标题、调色板、混合勾选状态、替代项和缩进。这些都是 `MenuItem` 和 `Menu` 上的方法，因此适用于应用程序菜单、上下文菜单、托盘菜单和 Dock 菜单。

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

- `SetSymbol` 在标题旁显示 SF Symbol（macOS 11+）。它会替换通过 `SetBitmap` 设置的图像。
- `SetBadge` 显示计数，`SetBadgeText` 在标题后显示短字符串（macOS 14+）。`ClearBadge` 将其移除；`BadgeCount` 和 `BadgeText` 可读取其值。
- `AddSectionHeader` 添加不可交互的标题（macOS 14+）。在更早版本上，它是具有相同标题的禁用菜单项。
- `AddPalette` 添加一行由 `NSMenu` 的调色板菜单支持的色块（macOS 14+）。可为每个色块传入一个符号，即每种颜色一个，或传入空切片以显示填充圆形。没有标签时，调色板内嵌在父菜单中；`SetLabel` 会将其显示为带标题的子菜单。`PaletteSelected` 返回选中的索引。
- `SetMixed` 将复选框设为以短横线绘制的混合状态。点击后会像 AppKit 一样将其完全选中。
- 按住不同的修饰键时，`SetAlternate(true)` 会使该项取代其上方的项目显示。这两个项目必须使用相同的按键，并具有不同的修饰键。
- `SetIndentationLevel` 可使标题最多缩进 15 级。

### 最近打开的项目

`fileMenu.AddRole(application.OpenRecent)` 添加标准的“最近打开的项目”子菜单。在 macOS 上，`NSDocumentController` 每次打开该菜单时都会填充内容，其中包括“清除菜单”项目。使用 `app.Menu.AddRecentDocument` 添加文件，使用 `RecentDocuments` 列出文件，使用 `ClearRecentDocuments` 清空列表。该列表在重新启动后仍会保留。

选择最近使用的文件时，会产生与从 Finder 打开文件相同的事件：

```go
app.Event.OnApplicationEvent(events.Common.ApplicationOpenedWithFile, func(event *application.ApplicationEvent) {
    path := event.Context().Filename()
    log.Println("open", path)
})
```

### Dock 菜单

`app.Menu.SetDockMenu` 为右键点击 Dock 图标安装静态菜单。`OnDockMenu` 会在每次即将显示时按需构建菜单，适合菜单项需要反映变化状态的情况。构建器的优先级高于静态菜单。

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

### Dock 进度

Dock 服务除了现有的徽章支持外，还可以在 Dock 图标上绘制进度条。将 `dock.New()` 注册为服务，并使用介于 0 和 1 之间的比例值调用 `SetProgress`。

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

`GetProgress` 返回当前比例值；未显示进度条时返回 `nil`。

## 对话框

消息、打开和保存对话框接受 macOS 专用选项，对话框管理器还增加了文本输入提示以及系统颜色和字体面板。

### 警告

`SetSuppression` 添加“不要再次显示此消息”复选框，`SetHelp` 显示帮助按钮。可以在按钮回调中通过 `Suppressed` 读取复选框状态，也可以注册 `OnSuppression` 以先接收该状态。

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

### 文本输入提示

`Prompt` 显示带文本字段的警告，并阻塞至其关闭，因此应从 goroutine 或绑定方法中调用。`Secure` 将该字段变为密码字段，`Window` 将警告作为工作表显示。

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

### 文件面板

`AddContentType` 在打开和保存对话框中按统一类型标识符筛选。它可以与 `AddFilter` 并用，因此 `"public.image"` 会匹配系统已知的所有图像类型，而筛选器仍可通过扩展名匹配 PDF。

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

保存面板增加了“格式”弹出菜单、名称字段的自定义标签和 Finder 标签。用户更改弹出菜单选项时，`SetFormats` 会切换允许的类型和名称字段中的扩展名；`SelectedFormat` 报告最终选择。

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

### 颜色和字体面板

`PickColor` 和 `PickFont` 打开共享的系统面板，并阻塞至面板关闭。面板打开期间，`OnChange` 会传递每次选择，使页面可以实时预览。每种面板同一时间只能打开一个；第二次调用会返回 `ErrDialogInProgress`。

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

## 状态项和反馈

### 状态项

macOS 上的系统托盘项目是一个 `NSStatusItem`。它可以使用 SF Symbol 绘制，带有工具提示，并允许用户像移除系统内置项目一样将其移除。

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

- `SetSymbol` 将符号渲染为模板图像，使其随菜单栏外观变化（macOS 11+）。`SetSymbolConfiguration` 设置点大小和字重。
- `SetTooltip` 设置悬停文字。`Tooltip` 可读取该值。
- `SetRemovable(true, name)` 允许用户按住 Command 键将项目拖出菜单栏。请为它提供稳定的自动保存名称，以便 macOS 在应用重新启动后仍记住移除状态。`Show` 或 `SetVisible(true)` 可将其恢复。
- `IsVisible` 读取 `NSStatusItem.visible`，因此用户移除项目后会返回 false。`OnVisibilityChange` 报告每次变化。

### 触觉反馈

应用程序处于活动状态时，`app.Haptics.Perform` 会在 Force Touch 触控板或 Magic Trackpad 上播放触觉反馈模式。

```go
app.Haptics.Perform(application.HapticAlignment)
```

类型包括 `HapticGeneric`、`HapticAlignment`（项目吸附到位）和 `HapticLevelChange`（档位或点击阶段）。`IsSupported` 报告平台是否能够播放触觉反馈。

### 声音

`app.Sound` 可播放警告声音、指定名称的系统声音或音频文件。

```go
app.Sound.Beep()

if err := app.Sound.Play("Glass"); err != nil {
    log.Println(err)
}

for _, name := range app.Sound.SystemSounds() {
    log.Println(name)
}
```

`Play` 接受 `SystemSounds` 中的名称，或 Core Audio 可以解码的任意文件的绝对路径。`PlayData` 从内存中播放完整的音频文件。

### 语音

`app.Speech.Speak` 将文本排入系统语音队列，并返回一个 `Utterance`。各段语音依次播放；`Stop` 丢弃一段，`StopAll` 清空队列。`Voices` 列出已安装语音及其标识符和语言。

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

`Recognize` 使用 `SFSpeechRecognizer` 转录默认麦克风的输入。首次调用会请求麦克风和语音识别权限，并阻塞至用户作出回应，因此应从 goroutine 中调用。部分转录结果通过 `OnPartial` 传递；`Stop` 结束采集并返回最终文本。

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

语音识别需要打包后的应用程序，并且其 `Info.plist` 必须声明 `NSSpeechRecognitionUsageDescription` 和 `NSMicrophoneUsageDescription`。缺少这些声明时，macOS 会拒绝访问，`Recognize` 返回 `ErrSpeechRecognitionUsageDescription`。

## 剪贴板和拖动

### 丰富内容剪贴板

除纯文本外，`app.Clipboard` 还可读写图像、文件引用、HTML、RTF，以及任意统一类型标识符下的原始数据。`Types` 列出粘贴板上的内容类型，`OnChange` 报告任意应用程序所做的更改。

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

`SetImage` 和 `Image` 使用 PNG 字节；其他应用程序以 TIFF 格式复制的图像会自动转换。系统不会通知剪贴板更改，因此只要存在至少一个监听器，`OnChange` 就会每隔 500 ms 轮询更改计数。

### 向外拖动

`StartDrag` 从窗口发起系统拖动，效果如同用户在 Finder 中拿起这些项目。它可提供现有文件、仅在目标位置接受放置后才生成内容的文件承诺，或纯文本。应在鼠标手势期间启动：绑定一个 Go 方法，并在可拖动元素的 `mousedown` 或 `pointerdown` 处理器中从页面调用它，同时将 HTML `draggable` 属性设为 `false`，防止 WebKit 自行开始拖动。

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

在手势之外调用 `StartDrag` 会返回 `ErrDragOutNoGesture`。`DragItems.Image` 和 `ImageOffset` 用于设置光标下的图像。

### 接收其他应用程序的拖放

文件拖放继续使用 `WindowFilesDropped` 事件。若要接受从其他应用程序拖来的文本、URL 或图像，请在 `DropTypes` 中列出类型并注册 `OnDrop`。这些拖放会交给 Go，而不会交给页面自身的 HTML5 放置处理器。

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

## 系统

### 权限

`app.Permissions` 可查询和请求系统隐私权限（摄像头、麦克风、屏幕录制、辅助功能、位置、通知、输入监控和完全磁盘访问权限）。`Status` 从不弹出权限提示。`Request` 会对尚未确定的权限类型弹出提示，并阻塞至用户作出回应，因此应从 goroutine 中调用。`OpenSystemSettings` 打开对应的“隐私与安全性”设置面板。

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

完全磁盘访问权限不能通过请求获取，调用会返回 `ErrPermissionNotRequestable`；请引导用户前往设置面板。始终没有得到回应的请求会返回 `ErrPermissionRequestTimeout`，在 macOS 上这通常意味着 `Info.plist` 缺少该类型的用途说明键。

macOS 现在会遵循 `Permissions` 窗口选项。它决定如何处理页面发出的 `getUserMedia` 请求：`PermissionAllow` 跳过 WebView 自身的提示，`PermissionDeny` 不询问就拒绝请求，`PermissionDefault` 显示提示。首次使用摄像头或麦克风时，系统级 TCC 提示仍会出现。

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

### 电源

`app.Power.PreventSleep` 会使系统保持唤醒；若指定 `Display`，屏幕也会保持唤醒，直到调用返回的释放函数。保持唤醒的请求会计数，因此应用程序的多个部分可以同时持有。原因会显示在活动监视器中。

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

变化会以 `events.Mac.ApplicationDidChangePowerState`（低电量模式切换）和 `events.Mac.ApplicationDidChangeThermalState` 事件传递。

### 生命周期

如果应用程序通过 `NSSupportsSuddenTermination` 选择启用突然终止，macOS 可以在注销或关机时立即终止空闲应用程序；如果通过 `NSSupportsAutomaticTermination` 选择启用自动终止，macOS 可以退出没有窗口的空闲应用程序。`app.Lifecycle.HoldTermination` 可在保存文件等关键代码段期间暂停这两种终止行为。

```go
release := app.Lifecycle.HoldTermination("Saving document")
defer release()
// write the file
```

`SetSuddenTerminationEnabled` 可在运行时切换突然终止；`SuddenTerminationEnabled` 报告当前状态，初始状态由 `Info.plist` 中的键决定。

### 环境

`app.Env` 增加了三个用于查询用户设置的方法。

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

- `Accessibility` 反映“减少动态效果”“降低透明度”“提高对比度”“不使用颜色进行区分”“反转颜色”、VoiceOver 和切换控制的设置。
- `KeyboardLayout` 返回当前输入源及其标识符、本地化名称和语言。
- `Locale` 返回 AppKit 为应用程序选择的区域设置，以及用户完整且有序的 `Preferred` 列表。`Identifier` 只反映应用包在 `CFBundleLocalizations` 中声明的语言；如需自行选择语言，请使用 `Preferred`。

### 事件

以下是新增的应用程序事件。每个事件都通过 `app.Event.OnApplicationEvent` 传递；请查询相应的管理器以获取最新值。

| 事件 | 触发时机 | 读取方式 |
|-------|------------|-----------|
| `events.Mac.ApplicationDidChangePowerState` | 低电量模式切换 | `app.Power.State()` |
| `events.Mac.ApplicationDidChangeThermalState` | 热压力变化 | `app.Power.State()` |
| `events.Common.AccessibilitySettingsChanged` | 辅助功能显示设置变化 | `app.Env.Accessibility()` |
| `events.Mac.ApplicationDidChangeAccessibilitySettings` | 同一变化的 macOS 专用事件 | `app.Env.Accessibility()` |
| `events.Mac.ApplicationDidChangeKeyboardLayout` | 输入源变化 | `app.Env.KeyboardLayout()` |
| `events.Mac.ApplicationDidChangeLocale` | 区域设置变化 | `app.Env.Locale()` |

```go
app.Event.OnApplicationEvent(events.Mac.ApplicationDidChangeThermalState, func(*application.ApplicationEvent) {
    app.Event.Emit("system:power", app.Power.State())
})
app.Event.OnApplicationEvent(events.Common.AccessibilitySettingsChanged, func(*application.ApplicationEvent) {
    app.Event.Emit("system:accessibility", app.Env.Accessibility())
})
```

## 集成

### 服务菜单

`app.ServicesProvider.Register` 在每个 macOS 应用程序针对选中文本或文件显示的“服务”子菜单中添加项目。处理器以 `ServiceRequest` 接收粘贴板，并返回要写回的 `ServiceResponse`；空响应会保持所选内容不变。请让处理器快速完成，因为 AppKit 会在主线程上等待它。

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

`Name` 是 AppKit 发送的消息，必须是普通标识符。`SendTypes` 和 `ReturnTypes` 是粘贴板类型；服务至少需要其中一个。仅注册并不会让服务可见：应用包的 `Info.plist` 必须在 `NSServices` 下声明服务。`InfoPlistXML` 返回可直接粘贴的配置块，`InfoPlistEntries` 则以映射形式返回相同数据，供 plist 序列化器使用。这些条目中的 `NSPortName` 是应用程序的 `Name`，必须与 `CFBundleName` 一致。

使用 Wails CLI 构建的项目可以在 `build/config.yml` 中声明一次相同的服务，由打包过程生成 `NSServices` 配置块：

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

每个条目都必须与 Go 中注册的 `ServiceDefinition` 具有相同的 `Name`。安装新构建后运行 `pbs -update`，使“服务”菜单无需注销即可识别更改。

### Handoff 和用户活动

`app.Activity.Publish` 将一个 `NSUserActivity` 设为当前活动，使用户可以在另一台设备上继续操作、在 Spotlight 中找到它，或收到 Siri 的建议。返回的 `PublishedActivity` 可以随状态变化而更新，并在文档关闭时失效。发布新活动会取代前一个活动。

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

传入的活动通过 `OnContinue` 送达。通用链接的类型为 `UserActivityTypeBrowsingWeb`，页面地址位于 `WebpageURL`；它们也会作为 `events.Common.ApplicationLaunchedWithUrl` 传递，使应用程序可以共用一条 URL 处理路径。`OnWillContinue`、`OnFailed` 和 `OnUpdated` 涵盖其余委托回调。

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

每种活动类型都必须列在 `Info.plist` 的 `NSUserActivityTypes` 下。通用链接还需要 `com.apple.developer.associated-domains` 授权，其中包含 `applinks:example.com` 条目，并且该域名上必须放置匹配的 `apple-app-site-association` 文件。

### Apple Events

`app.AppleEvents.Handle` 为事件类和 ID 注册处理器，使 AppleScript、Shortcuts 和其他应用程序能够驱动应用程序。这些代码是四字符字符串。直接参数会解码为 Go 值（`string`、由文件路径组成的 `[]string`、`int64`、`float64`、`bool`、`[]any` 或 `AppleEventRawData`），回复中的 `Result` 接受相同类型。事件暂停期间，处理器在自己的 goroutine 上运行。

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

脚本可以使用原始事件语法立即调用处理器：

```applescript
tell application id "com.example.notes" to «event WAILnote» "hello"
```

`ScriptingDefinition` 生成一个简化的 `.sdef`，为每个处理器提供命令名称。将其放入 `Contents/Resources`，并通过 `NSAppleScriptEnabled` 和 `OSAScriptingDefinition` 在 `Info.plist` 中指向它；之后脚本编辑器会在“文件 > 打开词典”中显示它。Wails 已处理自定义 URL 方案的“获取 URL”事件；为 `"GURL"`/`"GURL"` 注册的处理器会接续内置处理，而为 `"aevt"`/`"odoc"` 注册的处理器会替代内置的“打开文档”事件传递。`Send` 通过应用包标识符定位运行中的应用程序；在打包后的应用程序中，它需要 `NSAppleEventsUsageDescription`，并会阻塞调用方的 goroutine，直到收到回复。

### Quick Look

`app.QuickLook.Preview` 为一个或多个文件打开共享的 Quick Look 面板；传入多个路径时，面板会显示用于切换文件的箭头。`Thumbnail` 通过系统的缩略图提供者渲染文件并返回 PNG，因此适用于文档、图像、PDF 和影片。

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

路径必须是绝对路径，且文件必须存在。`ClosePreview` 和 `IsPreviewOpen` 用于管理面板。`ThumbnailOptions.IconMode` 会绘制 Finder 风格的文档边框，`Scale: 2` 则生成 Retina 图像。`Thumbnail` 会阻塞调用方的 goroutine，因此应从 goroutine 或绑定方法中调用。

### 工作区辅助方法

`app.Browser` 增加了三个基于 `NSWorkspace` 的辅助方法。`OpenWith` 使用由应用包标识符或应用包路径指定的应用程序打开文件。`ApplicationsForFile` 列出能够打开文件的已安装应用程序，默认处理程序排在最前面。`ActivateApplication` 将运行中的应用程序切换到前台。

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

`app.Spotlight.Index` 通过 Core Spotlight 将应用程序内容添加到系统搜索索引。每个 `SearchableItem` 都有 `ID` 和 `Title`，还可选填用于批量删除的 `Domain`、`Description`、`Keywords`、`ContentType`、PNG 缩略图、深层链接 `URL` 和到期时间。用户在 Spotlight 中选择其中一个项目，或带着查询选择“在 App 中搜索”时，会调用 `OnOpen`。

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

`IsAvailable` 报告索引是否接受项目。索引需要打包后的应用程序：未打包的 `go run` 二进制文件所索引的项目绝不会出现在 Spotlight 中。`DeleteAll` 会移除该应用程序索引的所有内容。

## 即将推出

正在添加一个涵盖工作表、弹出面板、呈现选项和状态恢复的 `mac-windows-extra` 示例；完成后会在此处添加链接。

## 版本要求

本页所有功能仅适用于 macOS，Go API 在所有平台上相同。Wails 支持 macOS 10.13 及更高版本；需要更新系统版本的功能会按下表所述降级。

| 功能 | 最低 macOS 版本 | 在更早版本上的行为 |
|---------|---------------|-------------------------------|
| 摄像头和麦克风权限状态 | 10.14 | 报告为已授权（更早版本不限制采集设备） |
| `Speech.Speak` 和 `Voices` | 10.14 | `ErrSpeechNotSupported` |
| 屏幕录制和输入监控权限 | 10.15 | 报告为已授权 |
| `Speech.Recognize` | 10.15 | `ErrSpeechRecognitionNotSupported` |
| `QuickLook.Thumbnail` | 10.15 | `ErrQuickLookNotSupported` |
| `SetSubtitle` | 11 | 忽略，并记录调试日志 |
| `ExportPDF` | 11 | `ErrMacExportUnsupported` |
| 对菜单项和状态项调用 `SetSymbol` | 11 | 不显示图像 |
| 按 UTI 使用 `AddContentType`、`SetFormats` | 11 | 通过旧版允许文件类型 API 应用相同的标识符 |
| `StartDrag` 中的文件承诺图标 | 11 | 通用文档图标 |
| `PowerState.LowPowerMode` 及其事件 | 12 | 始终为 false；事件从不触发 |
| `InteractionState` 和 `RestoreInteractionState` | 12 | `ErrMacInteractionStateUnsupported` |
| `OpenSystemSettings` 中的通知面板 | 13 | 打开旧版“通知”偏好设置面板 |
| 菜单徽章、分区标题、调色板 | 14 | 不显示徽章；标题是禁用项目；隐藏调色板 |
| 对没有自定义视图的项目调用 `MacToolbarItem.ShowPopover` | 14 | `ErrMacPopoverAnchorUnavailable` |

其他所有功能除 Wails 的最低版本要求外，没有额外要求。

## Info.plist 键

多项功能依赖应用程序 `Info.plist` 中的键。权限提示会向用户显示用途说明；如果缺少相应的键，macOS 就不会显示提示，请求最终会超时。

| 键 | 所需功能 |
|-----|-----------|
| `NSCameraUsageDescription` | `Permissions.Request(PermissionKindCamera)`、页面的摄像头访问 |
| `NSMicrophoneUsageDescription` | `Permissions.Request(PermissionKindMicrophone)`、页面的麦克风访问、`Speech.Recognize` |
| `NSSpeechRecognitionUsageDescription` | `Speech.Recognize` |
| `NSLocationUsageDescription` | `Permissions.Request(PermissionKindLocation)` |
| `NSSupportsSuddenTermination` | `Lifecycle.SuddenTerminationEnabled` 的初始状态；`HoldTermination` 可暂停该功能 |
| `NSSupportsAutomaticTermination` | 允许 macOS 退出空闲应用程序；`HoldTermination` 可暂停该功能 |
| `CFBundleLocalizations` | `Env.Locale().Identifier` 可以报告的语言 |
| `NSServices` | 为通过 `app.ServicesProvider` 注册的每项服务添加一个条目；可使用 `InfoPlistXML` 生成该配置块 |
| `NSUserActivityTypes` | 通过 `app.Activity` 发布或继续的每种 `UserActivity.Type` |
| `NSAppleScriptEnabled` 和 `OSAScriptingDefinition` | 将应用程序标记为可通过脚本操作，并指定由 `app.AppleEvents.ScriptingDefinition` 写出的 `.sdef` |
| `NSAppleEventsUsageDescription` | 向其他应用程序调用 `app.AppleEvents.Send` |

通知权限和语音识别还要求应用程序以带有应用包标识符的应用包形式运行；未打包的 `go run` 二进制文件会对通知报告 `PermissionStatusUnsupported`。

<a id="platform-notes"></a>

## 平台说明

本页所有 API 都可以在 Windows 和 Linux 上编译。在 macOS 以外的平台上：

- 窗口附加功能：`SetRepresentedFile`、`SetDocumentEdited`、`SetSubtitle`、`CascadeFrom`、`SetFrameAutosaveName` 和 `SetWindowButtonsOffset` 不执行任何操作。`WindowCascade` 的行为与 `WindowCentered` 相同。`RequestAttention` 返回的句柄中，`Cancel` 不执行任何操作。`PrintWithOptions` 调用 `Print`。`ExportPDF` 和 `Snapshot` 返回 `ErrMacOnly`。
- 菜单：符号、徽章、混合状态、替代项和缩进会保存，但不会绘制。分区标题是禁用项目。调色板会隐藏。创建菜单时，“最近打开的项目”子菜单会根据 Go 中的最近文件列表填充。Dock 菜单始终不会显示。
- 对话框：`SetSuppression`、`SetHelp`、`SetNameFieldLabel` 和 `SetTags` 会被忽略。对于已知的类型映射，`AddContentType` 和 `SetFormats` 会转为扩展名筛选器。`Prompt`、`PickColor` 和 `PickFont` 返回 `ErrDialogNotSupported`。
- 状态项：`SetSymbol`、`SetRemovable` 和 `OnVisibilityChange` 不起作用。`IsVisible` 反映最近一次 `Show` 或 `Hide` 调用。
- 反馈：触觉反馈可用于 iOS 和 Android；在 Windows 和 Linux 上不执行任何操作。`Sound.Beep` 和 `Sound.Play` 在 Windows 上支持 WAV 文件和注册表别名；在其他平台上，`Play` 返回 `ErrSoundNotSupported`。语音功能返回 `ErrSpeechNotSupported` 和 `ErrSpeechRecognitionNotSupported`。
- 剪贴板：丰富内容相关方法返回 `ErrClipboardNotSupported`，`Types` 为空，`ChangeCount` 为 0，`OnChange` 从不触发。
- 拖动：`StartDrag` 返回 `ErrDragOutUnsupported`。非文件拖放类型不会传递；文件拖放继续通过 `WindowFilesDropped` 工作。
- 系统：`Permissions.Status` 报告 `PermissionStatusUnsupported`，`Request` 返回 `ErrPermissionsUnsupported`。`PreventSleep` 返回 `ErrPreventSleepUnsupported` 和不执行任何操作的释放函数。`HoldTermination` 返回不执行任何操作的释放函数。`Accessibility` 的所有值均为 false，`KeyboardLayout` 是零值，`Locale` 根据 `LC_ALL`、`LC_MESSAGES` 和 `LANG` 推导。`Permissions` 窗口选项是跨平台的。
- 集成：`ServicesProvider.Register` 返回 `ErrServicesUnsupported`，`Activity.Publish` 返回 `ErrActivityUnsupported`；`InfoPlistXML`、`InfoPlistEntries` 和活动处理器仍可工作。`AppleEvents` 的 `Handle` 和 `Send` 返回 `ErrAppleEventsNotSupported`；`ScriptingDefinition` 在所有平台上都能生成。`QuickLook.Preview` 和 `Thumbnail` 返回 `ErrQuickLookNotSupported`。`Spotlight` 索引方法返回 `ErrSpotlightNotSupported`，`OnOpen` 从不触发。`Browser.OpenWith` 会启动指定的可执行文件，并将路径作为参数传入；`ApplicationsForFile` 为空，`ActivateApplication` 返回 `ErrApplicationNotRunning`。
- 呈现选项：`SetPresentationOptions` 返回 `ErrMacOnly`，`PresentationOptions` 为 `MacPresentationDefault`。
- 工作表：`PresentSheet`、`PresentCriticalSheet` 和 `PresentNativeSheet` 返回 `ErrMacSheetUnsupported`；`EndSheet` 不执行任何操作，查询方法报告没有工作表。
- 弹出面板：`NewMacPopover` 可调用，显示方法返回 `ErrMacPopoverUnsupported`，`IsShown` 为 false。
- 状态恢复：`SetRestorationID` 和 `SetRestorationData` 不执行任何操作，`OnRestore` 从不调用，`InteractionState` 和 `RestoreInteractionState` 返回 `ErrMacOnly`。

## 示例

每个示例都是完整且可运行的应用程序：

- [`v3/examples/mac-toolbar`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-toolbar)：`windowextras.go` 将已修改标记、副标题、PDF 导出、打印选项和级联窗口集成到笔记编辑器中。
- [`v3/examples/mac-menus-dock`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-menus-dock)：符号、徽章、分区标题、混合状态、替代项、调色板、最近打开的项目、动态 Dock 菜单和 Dock 进度。
- [`v3/examples/mac-dialogs`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-dialogs)：不再显示提示和帮助按钮、文本输入提示、内容类型、“格式”弹出菜单、Finder 标签，以及颜色和字体面板。
- [`v3/examples/mac-feedback`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-feedback)：可移除的 SF Symbol 状态项、触觉反馈、系统声音、文本转语音和语音识别。
- [`v3/examples/mac-clipboard-drag`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-clipboard-drag)：带更改跟踪的丰富内容剪贴板、使用文件承诺的向外拖动，以及文本、URL 和图像拖放。
- [`v3/examples/mac-system`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-system)：权限、阻止休眠、暂停终止、电源状态、辅助功能、键盘布局和区域设置及其实时更新。
- [`v3/examples/mac-integration`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-integration)：服务菜单项、带继续处理器的 Handoff 活动，以及 `app.Browser` 上的工作区辅助方法。
- [`v3/examples/mac-search-preview`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-search-preview)：带 `OnOpen` 的 Spotlight 索引、Quick Look 预览和缩略图，以及带脚本定义的自定义 Apple Event。
