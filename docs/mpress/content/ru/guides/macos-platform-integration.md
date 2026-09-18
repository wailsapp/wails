---
title: "Интеграция с платформой macOS"
description: "Окна документов, дополнительные возможности Dock и меню, нативные панели, элементы строки меню, обратная связь, расширенный буфер обмена, перетаскивание из приложения, разрешения, питание и жизненный цикл в macOS"
slug: "guides/macos-platform-integration"
sourcePath: "guides/macos-platform-integration.md"
---

Поддерживаемые платформы: macOS

Wails v3 обеспечивает поведение, которого пользователи macOS ожидают от нативного приложения: окна документов со значками файлов и каскадным размещением, меню со значками и индикаторами, меню Dock с индикатором хода выполнения, нативные предупреждения и панели, удаляемые элементы строки меню, тактильную обратную связь и речь, расширенный буфер обмена и перетаскивание из приложения, системные сведения о разрешениях, питании и локали, а также интеграцию с меню «Службы», Handoff, AppleScript и Quick Look. Всё управляется из Go через пакет `application`.

Тот же код компилируется в Windows и Linux. Сеттеры сохраняют значения, запросы возвращают нулевые значения, а операции, которым необходима macOS, возвращают документированную ошибку, например `ErrMacOnly`, `ErrDialogNotSupported` или `ErrClipboardNotSupported`. В разделе [Примечания о платформах](#platform-notes) ниже описано поведение каждой области вне macOS.

О нативном оформлении окна (панелях инструментов, боковых панелях, инспекторах, дополнительных элементах и вкладках окон) читайте в руководстве [Нативное оформление окна macOS](/guides/macos-native-chrome).

## Окна документов

Окно документа показывает связанный с ним файл в строке заголовка, отмечает несохранённые изменения точкой на кнопке закрытия и открывает новые окна каскадом. Всё это доступно через методы `WebviewWindow` и параметры `MacWindow`.

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

- `SetRepresentedFile` показывает значок связанного файла в строке заголовка. Пользователь может перетащить значок в другое приложение или щёлкнуть по нему с нажатой клавишей Command, чтобы увидеть путь. Передайте `""`, чтобы удалить значок. `RepresentedFile` возвращает заданный файл.
- `SetDocumentEdited` показывает точку несохранённых изменений на кнопке закрытия и затемняет значок файла. `IsDocumentEdited` возвращает это состояние.
- `SetSubtitle` показывает вторую строку под заголовком в macOS 11 и новее.
- `InitialPosition: application.WindowCascade` размещает окно ниже и правее последнего окна в каскаде, как при открытии новых документов. `CascadeFrom(other)` делает то же для существующего окна и обновляет точку каскада для последующих окон.
- `Mac.FrameAutosaveName` восстанавливает сохранённое положение и размер до первого показа окна и продолжает сохранять их при перемещении окна. Восстановленные параметры имеют приоритет над `X`, `Y`, `Width`, `Height` и `InitialPosition`. `SetFrameAutosaveName` меняет имя у работающего окна.
- `MacTitleBar.WindowButtonsOffset` смещает кнопки закрытия, сворачивания и изменения размера на заданное число пунктов. `SetWindowButtonsOffset` и `ResetWindowButtonsOffset` меняют смещение во время работы.

Все три сеттера можно вызывать до создания нативного окна; значения применяются при его создании.

### Запросы внимания

`RequestAttention` заставляет значок в Dock подпрыгивать, пока приложение находится в фоновом режиме. Информационный запрос вызывает один скачок. Критический запрос продолжает анимацию, пока пользователь не активирует приложение или вы не отмените запрос.

```go
request := window.RequestAttention(true)

// Once the work that needed attention is done:
request.Cancel()
```

`Flash` остаётся кроссплатформенным способом однократно привлечь внимание.

### Печать и экспорт

`PrintWithOptions` печатает содержимое WebView с явно заданными параметрами страницы. Нулевое значение показывает панель печати с общими настройками печати. `Print` сохраняет прежнее поведение (альбомная ориентация, поля по 30 пунктов).

`ExportPDF` преобразует страницу в PDF-документ, а `Snapshot` создаёт снимок в формате PNG. Оба метода ожидают WebKit, поэтому вызывайте их из горутины и никогда из потока приложения; вызов из этого потока возвращает `ErrMacExportOnMainThread`.

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

`PrintOptions` также принимает `PrinterName`, `PaperName` (имя PostScript, например `"iso-a4"`) и `Scale`. `PDFExportOptions` и `SnapshotOptions` принимают необязательный `Rect` для ограничения захватываемой области и `Timeout`, значение которого по умолчанию равно `DefaultMacExportTimeout` (30 секунд).

## Модальные панели

Модальная панель — второе окно, прикреплённое к верхней части родительского, как панель сохранения. Любой `WebviewWindow` можно показать как модальную панель другого окна с помощью `PresentSheet` и закрыть с помощью `EndSheet`, передав код ответа в обратные вызовы панели `OnSheetEnd`. Создавайте окно панели с установленным `Hidden`, чтобы оно не мелькало на экране до прикрепления; после закрытия AppKit снова убирает его с экрана, поэтому одно окно можно показывать многократно.

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

`PresentCriticalSheet` показывает панель перед любой уже прикреплённой обычной панелью вместо постановки в очередь за ней. `PresentNativeSheet` делает то же для `NativeWindow`. `IsSheet`, `SheetParent`, `AttachedSheet`, `AttachedNativeSheet` и `HasAttachedSheet` описывают текущее состояние. Если закрыть окно панели вместо завершения её работы, передаётся `MacSheetResponseStop`.

## Всплывающие панели

`MacPopover` — это `NSPopover`: временная панель, привязанная к прямоугольной области окна, элементу панели инструментов или элементу строки меню. Её содержимое — нативная полоса элементов управления `MacAccessory`, тот же тип, который используется для дополнительных элементов строки заголовка в руководстве [Нативное оформление окна macOS](/guides/macos-native-chrome). Добавьте все элементы управления до первого показа.

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

`MacToolbarItem.ShowPopover` и `SystemTray.ShowPopover` привязывают одну и ту же всплывающую панель к элементу панели инструментов или элементу строки меню. `MacPopoverBehaviorTransient` закрывает панель при любом щелчке вне её, `Semitransient` — только при щелчке в окне, из которого она открыта, а режим по умолчанию оставляет её открытой до вызова `Close`. `MacRectEdge` выбирает сторону появления панели. `SetContentSize` и `SetBehavior` настраивают уже открытую панель, а `Destroy` освобождает нативную панель и её полосу содержимого для использования в другом месте.

## Восстановление состояния

macOS повторно открывает окна приложения после сбоя, принудительного завершения или перезагрузки, а также после обычного завершения, если в Системных настройках выключен параметр «Закрывать окна при завершении приложения». Задайте окну `Mac.RestorationID`, сохраните данные, необходимые для его воссоздания, с помощью `SetRestorationData` и зарегистрируйте `app.Window.OnRestore`, чтобы при следующем запуске создать окно заново.

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

Сохраняются только окна, видимые в момент завершения приложения. `RestorationState.Data` — строковая карта; храните в ней идентификаторы, пути и позиции. `InteractionState` возвращает список переходов WebView назад и вперёд и позиции прокрутки в виде непрозрачного блока данных (macOS 12+), который `RestoreInteractionState` применяет к воссозданному окну; обычно его сохраняют в данных восстановления в кодировке base64. `SetRestorationID` и `RestorationID` изменяют и считывают идентификатор работающего окна.

## Параметры отображения

`MacPresentationOptions` соответствует `NSApplication.presentationOptions`: битовой маске, которая скрывает Dock или строку меню и отключает переключение процессов, принудительное завершение, выход из учётной записи или команду скрытия, пока приложение активно. Задайте `Mac.PresentationOptions` в параметрах приложения для применения при запуске или меняйте значение во время работы с помощью `SetPresentationOptions`. Недопустимые сочетания отклоняются до передачи в AppKit с ошибкой, оборачивающей `ErrMacPresentationOptionsInvalid`.

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

Для скрытия строки меню (`HideMenuBar` или `AutoHideMenuBar`) требуется один из параметров Dock, а для `AutoHideToolbar` — одновременно `FullScreen` и `AutoHideMenuBar`. `Validate` сообщает о первом нарушенном правиле, а `Has` проверяет отдельные флаги.

## Меню и Dock

Пункты меню могут содержать SF Symbols, индикаторы, заголовки разделов, цветовые палитры, смешанное состояние флажка, альтернативные варианты и отступы. Всё это методы `MenuItem` и `Menu`, поэтому они работают в меню приложения, контекстных меню, меню значка в области уведомлений и меню Dock.

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

- `SetSymbol` показывает SF Symbol рядом с названием (macOS 11+). Он заменяет изображение, заданное через `SetBitmap`.
- `SetBadge` показывает число, а `SetBadgeText` — короткую строку после названия (macOS 14+). `ClearBadge` удаляет индикатор; `BadgeCount` и `BadgeText` считывают его значение.
- `AddSectionHeader` добавляет неинтерактивный заголовок (macOS 14+). В более ранних версиях это отключённый пункт с тем же названием.
- `AddPalette` добавляет строку образцов цвета на основе меню палитры `NSMenu` (macOS 14+). Передайте по одному символу для каждого образца, по одному для каждого цвета, либо пустой срез для заполненных кружков. Без метки палитра встроена в родительское меню; `SetLabel` показывает её как подменю с заголовком. `PaletteSelected` возвращает индекс выбранного образца.
- `SetMixed` переводит флажок в смешанное состояние, обозначенное чертой. При нажатии он полностью включается, как в AppKit.
- `SetAlternate(true)` показывает пункт вместо расположенного над ним, пока нажата отличающаяся клавиша-модификатор. Оба пункта должны иметь одну клавишу и различаться модификаторами.
- `SetIndentationLevel` задаёт отступ названия до 15 уровней.

### Недавние файлы

`fileMenu.AddRole(application.OpenRecent)` добавляет стандартное подменю «Недавние файлы». В macOS `NSDocumentController` заполняет его при каждом открытии; оно содержит пункт очистки меню. Добавляйте файлы с помощью `app.Menu.AddRecentDocument`, получайте список через `RecentDocuments` и очищайте его через `ClearRecentDocuments`. Список сохраняется между запусками.

Выбор недавнего файла вызывает то же событие, что и открытие файла из Finder:

```go
app.Event.OnApplicationEvent(events.Common.ApplicationOpenedWithFile, func(event *application.ApplicationEvent) {
    path := event.Context().Filename()
    log.Println("open", path)
})
```

### Меню Dock

`app.Menu.SetDockMenu` устанавливает статическое меню, открываемое щелчком правой кнопкой по значку Dock. `OnDockMenu` создаёт меню по требованию перед каждым показом; используйте его, когда пункты должны отражать меняющееся состояние. Генератор имеет приоритет над статическим меню.

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

### Индикатор выполнения в Dock

Служба Dock рисует индикатор выполнения поверх значка Dock в дополнение к существующей поддержке значков с числом. Зарегистрируйте `dock.New()` как службу и вызовите `SetProgress` с долей от 0 до 1.

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

`GetProgress` возвращает текущую долю или `nil`, если индикатор не показан.

## Диалоговые окна

Диалоговые окна сообщений, открытия и сохранения принимают параметры macOS, а диспетчер диалоговых окон также предоставляет запрос текста и системные панели цвета и шрифта.

### Предупреждения

`SetSuppression` добавляет флажок «Больше не показывать это сообщение», а `SetHelp` показывает кнопку справки. Считывайте состояние флажка через `Suppressed` в обратном вызове кнопки либо зарегистрируйте `OnSuppression`, чтобы получить его раньше.

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

### Запрос текста

`Prompt` показывает предупреждение с текстовым полем и блокирует выполнение до его закрытия, поэтому вызывайте его из горутины или привязанного метода. `Secure` превращает поле в поле пароля, а `Window` показывает предупреждение как модальную панель.

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

### Панели выбора файлов

`AddContentType` фильтрует файлы по универсальному идентификатору типа как в диалоговом окне открытия, так и в окне сохранения. Он используется вместе с `AddFilter`: например, `"public.image"` соответствует всем известным системе типам изображений, а фильтр по расширению по-прежнему позволяет выбрать PDF.

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

Панели сохранения получают раскрывающийся список «Формат», собственную метку поля имени и теги Finder. Когда пользователь меняет значение в списке, `SetFormats` меняет разрешённый тип и расширение в поле имени; `SelectedFormat` сообщает окончательный выбор.

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

### Панели цвета и шрифта

`PickColor` и `PickFont` открывают общие системные панели и блокируют выполнение до их закрытия. `OnChange` передаёт каждое выбранное значение, пока панель открыта, позволяя странице показывать результат сразу. Одновременно может быть открыта только одна панель каждого вида; повторный вызов возвращает `ErrDialogInProgress`.

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

## Элементы строки меню и обратная связь

### Элементы строки меню

В macOS элемент области уведомлений представлен `NSStatusItem`. Его можно нарисовать с помощью SF Symbol, снабдить всплывающей подсказкой и позволить пользователю удалять так же, как встроенные элементы.

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

- `SetSymbol` отображает символ как шаблонное изображение, чтобы он соответствовал оформлению строки меню (macOS 11+). `SetSymbolConfiguration` задаёт размер в пунктах и насыщенность.
- `SetTooltip` задаёт текст при наведении. `Tooltip` считывает его.
- `SetRemovable(true, name)` позволяет пользователю перетащить элемент из строки меню с нажатой клавишей Command. Назначьте постоянное имя для автосохранения, чтобы macOS помнила удаление между запусками. `Show` или `SetVisible(true)` возвращает элемент.
- `IsVisible` считывает `NSStatusItem.visible`, поэтому после удаления элемента пользователем возвращает false. `OnVisibilityChange` сообщает о каждом изменении.

### Тактильная обратная связь

`app.Haptics.Perform` воспроизводит шаблон тактильного отклика на трекпаде Force Touch или Magic Trackpad, пока приложение активно.

```go
app.Haptics.Perform(application.HapticAlignment)
```

Виды отклика: `HapticGeneric`, `HapticAlignment` (элемент встал на место) и `HapticLevelChange` (переход через фиксированное положение или ступень щелчка). `IsSupported` сообщает, может ли платформа вообще воспроизводить тактильный отклик.

### Звуки

`app.Sound` воспроизводит звук предупреждения, системный звук по имени или аудиофайл.

```go
app.Sound.Beep()

if err := app.Sound.Play("Glass"); err != nil {
    log.Println(err)
}

for _, name := range app.Sound.SystemSounds() {
    log.Println(name)
}
```

`Play` принимает имя из `SystemSounds` или абсолютный путь к любому файлу, который умеет декодировать Core Audio. `PlayData` воспроизводит из памяти целый аудиофайл.

### Речь

`app.Speech.Speak` ставит текст в очередь для системного голоса и возвращает `Utterance`. Высказывания воспроизводятся одно за другим; `Stop` удаляет одно, а `StopAll` очищает очередь. `Voices` перечисляет установленные голоса с их идентификаторами и языками.

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

`Recognize` распознаёт речь с микрофона по умолчанию через `SFSpeechRecognizer`. При первом вызове запрашиваются разрешения на доступ к микрофону и распознавание речи; вызов блокируется до ответа пользователя, поэтому выполняйте его из горутины. Промежуточные результаты поступают через `OnPartial`; `Stop` завершает запись и возвращает окончательный текст.

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

Для распознавания нужно приложение в составе пакета, где в `Info.plist` объявлены `NSSpeechRecognitionUsageDescription` и `NSMicrophoneUsageDescription`. Без них macOS отказывает в доступе, а `Recognize` возвращает `ErrSpeechRecognitionUsageDescription`.

## Буфер обмена и перетаскивание

### Расширенный буфер обмена

`app.Clipboard` читает и записывает изображения, ссылки на файлы, HTML, RTF и необработанные данные под любым универсальным идентификатором типа в дополнение к обычному тексту. `Types` перечисляет содержимое буфера обмена, а `OnChange` сообщает об изменениях, сделанных любым приложением.

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

`SetImage` и `Image` работают с байтами PNG; изображения, скопированные другими приложениями в TIFF, преобразуются автоматически. Системного уведомления об изменениях буфера обмена нет, поэтому `OnChange` проверяет счётчик изменений каждые 500 мс, пока существует хотя бы один слушатель.

### Перетаскивание из приложения

`StartDrag` начинает системное перетаскивание из окна, словно пользователь взял элементы в Finder. Можно перетаскивать существующие файлы, обещанные файлы, содержимое которых создаётся только после принятия их целевым приложением, или обычный текст. Запускайте его во время жеста мышью: привяжите метод Go и вызовите его из обработчика `mousedown` или `pointerdown` перетаскиваемого элемента на странице, установив HTML-атрибут `draggable` равным `false`, чтобы WebKit не начал собственное перетаскивание.

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

Вызов `StartDrag` вне жеста возвращает `ErrDragOutNoGesture`. `DragItems.Image` и `ImageOffset` задают изображение под указателем.

### Данные, перетаскиваемые из других приложений

Для перетаскиваемых файлов по-прежнему используется событие `WindowFilesDropped`. Чтобы принимать текст, URL или изображения из других приложений, перечислите типы в `DropTypes` и зарегистрируйте `OnDrop`. Эти данные передаются Go, а не собственным обработчикам HTML5-события drop на странице.

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

## Система

### Разрешения

`app.Permissions` сообщает о системных разрешениях на доступ к личным данным и запрашивает их: камера, микрофон, запись экрана, универсальный доступ, местоположение, уведомления, мониторинг ввода и полный доступ к диску. `Status` никогда не показывает запрос. `Request` запрашивает разрешения с ещё не определённым состоянием и блокируется до ответа пользователя, поэтому вызывайте его из горутины. `OpenSystemSettings` открывает соответствующий раздел «Конфиденциальность и безопасность».

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

Полный доступ к диску нельзя запросить программно: возвращается `ErrPermissionNotRequestable`; направьте пользователя в раздел настроек. Если ответ на запрос так и не поступает, возвращается `ErrPermissionRequestTimeout`. В macOS это обычно означает, что в `Info.plist` отсутствует ключ с описанием использования для данного вида разрешения.

Параметр окна `Permissions` теперь учитывается в macOS. Он определяет обработку запросов `getUserMedia` со страницы: `PermissionAllow` пропускает собственный запрос WebView, `PermissionDeny` отказывает без запроса, а `PermissionDefault` показывает запрос. Системный запрос TCC всё равно появляется при первом использовании камеры или микрофона.

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

### Питание

`app.Power.PreventSleep` не даёт системе переходить в режим сна, а с `Display` удерживает включённым и экран до вызова возвращённой функции освобождения. Удержания подсчитываются, поэтому несколько частей приложения могут использовать их одновременно. Причина отображается в Мониторинге системы.

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

Изменения поступают как `events.Mac.ApplicationDidChangePowerState` (переключение режима энергосбережения) и `events.Mac.ApplicationDidChangeThermalState`.

### Жизненный цикл

macOS может мгновенно завершить неактивное приложение при выходе из учётной записи или выключении, если оно разрешило это через `NSSupportsSuddenTermination`, и может закрыть неактивное приложение без окон, если оно разрешило это через `NSSupportsAutomaticTermination`. `app.Lifecycle.HoldTermination` приостанавливает оба механизма на время критического участка, например сохранения файла.

```go
release := app.Lifecycle.HoldTermination("Saving document")
defer release()
// write the file
```

`SetSuddenTerminationEnabled` переключает внезапное завершение во время работы; `SuddenTerminationEnabled` сообщает текущее состояние, начальное значение которого берётся из ключа `Info.plist`.

### Окружение

`app.Env` предоставляет три запроса о настройках пользователя.

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

- `Accessibility` отражает настройки уменьшения движения, уменьшения прозрачности, увеличения контрастности, различения без цвета, инверсии цветов, VoiceOver и управления переключателями.
- `KeyboardLayout` возвращает активный источник ввода с его идентификатором, локализованным названием и языками.
- `Locale` возвращает локаль, выбранную AppKit для приложения, и полный упорядоченный список пользовательских предпочтений `Preferred`. `Identifier` отражает только язык, объявленный пакетом в `CFBundleLocalizations`; для самостоятельного выбора языка используйте `Preferred`.

### События

Эти события приложения добавлены недавно. Каждое передаётся через `app.Event.OnApplicationEvent`; чтобы получить актуальное значение, обратитесь к соответствующему менеджеру.

| Событие | Когда возникает | Как получить значение |
|-------|------------|-----------|
| `events.Mac.ApplicationDidChangePowerState` | переключается режим энергосбережения | `app.Power.State()` |
| `events.Mac.ApplicationDidChangeThermalState` | меняется тепловая нагрузка | `app.Power.State()` |
| `events.Common.AccessibilitySettingsChanged` | меняется настройка отображения для универсального доступа | `app.Env.Accessibility()` |
| `events.Mac.ApplicationDidChangeAccessibilitySettings` | вариант того же изменения для macOS | `app.Env.Accessibility()` |
| `events.Mac.ApplicationDidChangeKeyboardLayout` | меняется источник ввода | `app.Env.KeyboardLayout()` |
| `events.Mac.ApplicationDidChangeLocale` | меняется локаль | `app.Env.Locale()` |

```go
app.Event.OnApplicationEvent(events.Mac.ApplicationDidChangeThermalState, func(*application.ApplicationEvent) {
    app.Event.Emit("system:power", app.Power.State())
})
app.Event.OnApplicationEvent(events.Common.AccessibilitySettingsChanged, func(*application.ApplicationEvent) {
    app.Event.Emit("system:accessibility", app.Env.Accessibility())
})
```

## Интеграция

### Меню «Службы»

`app.ServicesProvider.Register` добавляет пункт в подменю «Службы», которое каждое приложение macOS показывает для выделенного текста или файлов. Обработчик получает буфер обмена как `ServiceRequest` и возвращает `ServiceResponse` для записи результата; пустой ответ оставляет выделение без изменений. Обработчик должен работать быстро, поскольку AppKit ожидает его в главном потоке.

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

`Name` — сообщение, которое отправляет AppKit; это должен быть простой идентификатор. `SendTypes` и `ReturnTypes` — типы буфера обмена; службе нужен хотя бы один из них. Одной регистрации недостаточно для появления службы: `Info.plist` пакета должен объявлять её в `NSServices`. `InfoPlistXML` возвращает готовый для вставки блок, а `InfoPlistEntries` возвращает те же данные в виде карт для сериализатора plist. `NSPortName` в этих записях — это `Name` приложения, которое должно совпадать с `CFBundleName`.

Проекты, собранные с помощью Wails CLI, могут один раз объявить те же службы в `build/config.yml`, чтобы при упаковке создавался блок `NSServices`:

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

Каждая запись должна соответствовать `ServiceDefinition`, зарегистрированному в Go с тем же `Name`. После установки новой сборки выполните `pbs -update`, чтобы меню «Службы» подхватило изменение без выхода из учётной записи.

### Handoff и действия пользователя

`app.Activity.Publish` делает `NSUserActivity` текущим действием, чтобы пользователь мог продолжить его на другом устройстве, найти в Spotlight или получить предложение от Siri. Возвращённый `PublishedActivity` можно обновлять при изменении состояния и аннулировать при закрытии документа. Публикация нового действия заменяет предыдущее.

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

Входящие действия поступают через `OnContinue`. Универсальные ссылки имеют тип `UserActivityTypeBrowsingWeb`, а адрес страницы находится в `WebpageURL`; они также передаются как `events.Common.ApplicationLaunchedWithUrl`, чтобы приложение могло обрабатывать URL одним способом. `OnWillContinue`, `OnFailed` и `OnUpdated` охватывают остальные методы делегата.

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

Каждый тип действия должен быть указан в `NSUserActivityTypes` файла `Info.plist`. Для универсальных ссылок также требуется право `com.apple.developer.associated-domains` с записью `applinks:example.com` и соответствующий файл `apple-app-site-association` в этом домене.

### Apple Events

`app.AppleEvents.Handle` регистрирует обработчик класса и идентификатора события, чтобы AppleScript, Shortcuts и другие приложения могли управлять приложением. Коды представляют собой строки из четырёх символов. Прямой параметр декодируется в значение Go (`string`, `[]string` с путями к файлам, `int64`, `float64`, `bool`, `[]any` или `AppleEventRawData`), а `Result` ответа принимает те же виды значений. Обработчики работают в отдельных горутинах, пока событие приостановлено.

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

Скрипт может сразу вызвать обработчик, используя необработанный синтаксис события:

```applescript
tell application id "com.example.notes" to «event WAILnote» "hello"
```

`ScriptingDefinition` создаёт минимальный `.sdef`, назначающий каждому обработчику имя команды. Поместите его в `Contents/Resources` и укажите в `Info.plist` с помощью `NSAppleScriptEnabled` и `OSAScriptingDefinition`; после этого Редактор скриптов покажет его в меню «Файл > Открыть словарь». Wails уже обрабатывает событие Get URL для собственных схем URL; обработчик для `"GURL"`/`"GURL"` вызывается в цепочке с ним, а обработчик для `"aevt"`/`"odoc"` заменяет встроенную доставку событий открытия документов. `Send` обращается к работающему приложению по идентификатору пакета, требует `NSAppleEventsUsageDescription` в приложении, собранном как пакет, и блокирует вызывающую горутину до получения ответа.

### Quick Look

`app.QuickLook.Preview` открывает общую панель Quick Look для одного или нескольких файлов; если путей несколько, на панели появляются стрелки для перехода между ними. `Thumbnail` создаёт миниатюру файла через системные провайдеры и возвращает PNG, поэтому поддерживаются документы, изображения, PDF и видео.

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

Пути должны быть абсолютными и указывать на существующие файлы. `ClosePreview` и `IsPreviewOpen` управляют панелью. `ThumbnailOptions.IconMode` рисует рамку документа в стиле Finder, а `Scale: 2` создаёт изображение Retina. `Thumbnail` блокирует вызывающую горутину, поэтому вызывайте его из горутины или привязанного метода.

### Вспомогательные средства для работы с приложениями

`app.Browser` предоставляет три вспомогательных метода на основе `NSWorkspace`. `OpenWith` открывает файл в конкретном приложении, указанном идентификатором пакета или путём к пакету. `ApplicationsForFile` перечисляет установленные приложения, способные открыть файл, начиная с обработчика по умолчанию. `ActivateApplication` выводит работающее приложение на передний план.

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

`app.Spotlight.Index` добавляет содержимое приложения в системный поисковый индекс через Core Spotlight. У каждого `SearchableItem` есть `ID`, `Title` и, необязательно, `Domain` для массового удаления, `Description`, `Keywords`, `ContentType`, миниатюра PNG, `URL` глубокой ссылки и срок действия. `OnOpen` вызывается, когда пользователь выбирает элемент в Spotlight или выбирает «Искать в приложении» с запросом.

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

`IsAvailable` сообщает, принимает ли индекс элементы. Для индексирования приложение должно быть собрано как пакет: элементы, добавленные исполняемым файлом без пакета через `go run`, никогда не появляются в Spotlight. `DeleteAll` удаляет всё, что приложение добавило в индекс.

## Следующие дополнения

Пример `mac-windows-extra`, посвящённый модальным и всплывающим панелям, параметрам отображения и восстановлению состояния, готовится к добавлению; ссылка появится здесь после его появления.

## Требования к версии

Всё на этой странице относится только к macOS, а API Go везде одинаков. Wails поддерживает macOS 10.13 и новее; возможности, требующие более новой версии, работают с указанными ограничениями.

| Возможность | Минимальная версия macOS | Поведение в более ранних версиях |
|---------|---------------|-------------------------------|
| Состояние разрешений для камеры и микрофона | 10.14 | сообщается, что доступ разрешён (в более ранних версиях доступ к устройствам захвата не ограничивается) |
| `Speech.Speak` и `Voices` | 10.14 | `ErrSpeechNotSupported` |
| Разрешения на запись экрана и мониторинг ввода | 10.15 | сообщается, что доступ разрешён |
| `Speech.Recognize` | 10.15 | `ErrSpeechRecognitionNotSupported` |
| `QuickLook.Thumbnail` | 10.15 | `ErrQuickLookNotSupported` |
| `SetSubtitle` | 11 | игнорируется с записью в журнал отладки |
| `ExportPDF` | 11 | `ErrMacExportUnsupported` |
| `SetSymbol` для пунктов меню и элементов строки меню | 11 | изображение не показывается |
| `AddContentType`, `SetFormats` по UTI | 11 | те же идентификаторы применяются через прежний API разрешённых типов файлов |
| Значки обещанных файлов в `StartDrag` | 11 | общий значок документа |
| `PowerState.LowPowerMode` и его событие | 12 | всегда false; событие никогда не возникает |
| `InteractionState` и `RestoreInteractionState` | 12 | `ErrMacInteractionStateUnsupported` |
| Раздел уведомлений в `OpenSystemSettings` | 13 | открывается прежняя панель настроек уведомлений |
| Индикаторы меню, заголовки разделов, палитры | 14 | индикаторы не показываются; заголовки становятся отключёнными пунктами; палитры скрыты |
| `MacToolbarItem.ShowPopover` для элементов без собственного представления | 14 | `ErrMacPopoverAnchorUnavailable` |

Для всего остального нет требований сверх минимальной версии Wails.

## Ключи Info.plist

Некоторые возможности зависят от ключей в `Info.plist` приложения. Описания использования показываются пользователю в запросе разрешения; без соответствующего ключа macOS не показывает запрос, и ожидание ответа истекает.

| Ключ | Для чего нужен |
|-----|-----------|
| `NSCameraUsageDescription` | `Permissions.Request(PermissionKindCamera)`, доступ к камере со страницы |
| `NSMicrophoneUsageDescription` | `Permissions.Request(PermissionKindMicrophone)`, доступ к микрофону со страницы, `Speech.Recognize` |
| `NSSpeechRecognitionUsageDescription` | `Speech.Recognize` |
| `NSLocationUsageDescription` | `Permissions.Request(PermissionKindLocation)` |
| `NSSupportsSuddenTermination` | начальное состояние `Lifecycle.SuddenTerminationEnabled`; `HoldTermination` приостанавливает его |
| `NSSupportsAutomaticTermination` | позволяет macOS завершать неактивное приложение; `HoldTermination` приостанавливает его |
| `CFBundleLocalizations` | языки, о которых может сообщать `Env.Locale().Identifier` |
| `NSServices` | по одной записи на каждую службу, зарегистрированную через `app.ServicesProvider`; создайте блок с помощью `InfoPlistXML` |
| `NSUserActivityTypes` | каждый `UserActivity.Type`, опубликованный или продолженный через `app.Activity` |
| `NSAppleScriptEnabled` и `OSAScriptingDefinition` | отмечают приложение как поддерживающее скрипты и указывают `.sdef`, созданный через `app.AppleEvents.ScriptingDefinition` |
| `NSAppleEventsUsageDescription` | `app.AppleEvents.Send` для обращения к другим приложениям |

Для разрешения на уведомления и распознавания речи приложение также должно запускаться как пакет с идентификатором пакета; исполняемый файл без пакета, запущенный через `go run`, сообщает `PermissionStatusUnsupported` для уведомлений.

<a id="platform-notes"></a>

## Примечания о платформах

Все API на этой странице компилируются в Windows и Linux. Вне macOS:

- Дополнительные возможности окон: `SetRepresentedFile`, `SetDocumentEdited`, `SetSubtitle`, `CascadeFrom`, `SetFrameAutosaveName` и `SetWindowButtonsOffset` ничего не делают. `WindowCascade` ведёт себя как `WindowCentered`. `RequestAttention` возвращает дескриптор, чей `Cancel` ничего не делает. `PrintWithOptions` вызывает `Print`. `ExportPDF` и `Snapshot` возвращают `ErrMacOnly`.
- Меню: символы, индикаторы, смешанное состояние, альтернативные пункты и отступы сохраняются, но не отображаются. Заголовки разделов становятся отключёнными пунктами. Палитры скрыты. Подменю недавних файлов заполняется из списка Go при создании меню. Меню Dock никогда не показывается.
- Диалоговые окна: `SetSuppression`, `SetHelp`, `SetNameFieldLabel` и `SetTags` игнорируются. `AddContentType` и `SetFormats` превращаются в фильтры по расширениям, если известен способ преобразования. `Prompt`, `PickColor` и `PickFont` возвращают `ErrDialogNotSupported`.
- Элементы строки меню: `SetSymbol`, `SetRemovable` и `OnVisibilityChange` не действуют. `IsVisible` отражает последний вызов `Show` или `Hide`.
- Обратная связь: тактильная обратная связь работает в iOS и Android; в Windows и Linux вызовы ничего не делают. `Sound.Beep` и `Sound.Play` работают в Windows с файлами WAV и псевдонимами реестра; на других платформах `Play` возвращает `ErrSoundNotSupported`. Речевые функции возвращают `ErrSpeechNotSupported` и `ErrSpeechRecognitionNotSupported`.
- Буфер обмена: расширенные методы возвращают `ErrClipboardNotSupported`, `Types` пуст, `ChangeCount` равен 0, а `OnChange` никогда не вызывается.
- Перетаскивание: `StartDrag` возвращает `ErrDragOutUnsupported`. Данные, отличные от файлов, не передаются; перетаскивание файлов продолжает работать через `WindowFilesDropped`.
- Система: `Permissions.Status` сообщает `PermissionStatusUnsupported`, а `Request` возвращает `ErrPermissionsUnsupported`. `PreventSleep` возвращает `ErrPreventSleepUnsupported` и функцию освобождения, которая ничего не делает. `HoldTermination` возвращает функцию освобождения, которая ничего не делает. Все значения `Accessibility` равны false, `KeyboardLayout` имеет нулевое значение, а `Locale` определяется по `LC_ALL`, `LC_MESSAGES` и `LANG`. Параметр окна `Permissions` работает на разных платформах.
- Интеграция: `ServicesProvider.Register` возвращает `ErrServicesUnsupported`, а `Activity.Publish` — `ErrActivityUnsupported`; `InfoPlistXML`, `InfoPlistEntries` и обработчики действий продолжают работать. `Handle` и `Send` у `AppleEvents` возвращают `ErrAppleEventsNotSupported`; `ScriptingDefinition` создаётся везде. `QuickLook.Preview` и `Thumbnail` возвращают `ErrQuickLookNotSupported`. Методы индексирования `Spotlight` возвращают `ErrSpotlightNotSupported`, а `OnOpen` никогда не вызывается. `Browser.OpenWith` запускает указанный исполняемый файл, передавая путь аргументом; `ApplicationsForFile` пуст, а `ActivateApplication` возвращает `ErrApplicationNotRunning`.
- Параметры отображения: `SetPresentationOptions` возвращает `ErrMacOnly`, а `PresentationOptions` имеет значение `MacPresentationDefault`.
- Модальные панели: `PresentSheet`, `PresentCriticalSheet` и `PresentNativeSheet` возвращают `ErrMacSheetUnsupported`; `EndSheet` ничего не делает, а методы запросов сообщают об отсутствии панели.
- Всплывающие панели: `NewMacPopover` работает, методы показа возвращают `ErrMacPopoverUnsupported`, а `IsShown` равен false.
- Восстановление состояния: `SetRestorationID` и `SetRestorationData` ничего не делают, `OnRestore` никогда не вызывается, а `InteractionState` и `RestoreInteractionState` возвращают `ErrMacOnly`.

## Примеры

Каждый пример — полноценное приложение, которое можно запустить:

- [`v3/examples/mac-toolbar`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-toolbar): `windowextras.go` добавляет в редактор заметок точку изменённого документа, подзаголовок, экспорт в PDF, параметры печати и каскадное размещение окон.
- [`v3/examples/mac-menus-dock`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-menus-dock): символы, индикаторы, заголовки разделов, смешанное состояние, альтернативные пункты, палитра, недавние файлы, динамическое меню Dock и индикатор выполнения в Dock.
- [`v3/examples/mac-dialogs`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-dialogs): флажок «Больше не показывать» и кнопка справки, запросы текста, типы содержимого, список «Формат», теги Finder и панели цвета и шрифта.
- [`v3/examples/mac-feedback`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-feedback): удаляемый элемент строки меню с SF Symbol, тактильная обратная связь, системные звуки, синтез и распознавание речи.
- [`v3/examples/mac-clipboard-drag`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-clipboard-drag): расширенный буфер обмена с отслеживанием изменений, перетаскивание из приложения с обещанными файлами и приём текста, URL и изображений.
- [`v3/examples/mac-system`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-system): разрешения, предотвращение сна, удержание от завершения, состояние питания, настройки универсального доступа, раскладка клавиатуры и локаль с обновлениями во время работы.
- [`v3/examples/mac-integration`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-integration): пункт меню «Службы», действие Handoff с обработчиками продолжения и вспомогательные методы для работы с приложениями в `app.Browser`.
- [`v3/examples/mac-search-preview`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-search-preview): индексирование Spotlight с `OnOpen`, предварительный просмотр и миниатюры Quick Look, а также собственное событие Apple Event с его определением для скриптов.
