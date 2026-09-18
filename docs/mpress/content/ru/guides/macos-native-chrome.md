---
title: "Нативное оформление окна macOS"
description: "Создавайте нативные панели инструментов, боковые панели, списки содержимого, инспекторы, дополнительные элементы и вкладки окна AppKit вокруг окна Wails"
slug: "guides/macos-native-chrome"
sourcePath: "guides/macos-native-chrome.md"
---

Поддерживаемые платформы: macOS

Wails v3 позволяет окружить WebView настоящими нативными элементами окна AppKit. Панель инструментов, боковая панель, список содержимого, инспектор и элементы строки заголовка — это нативные элементы управления, создаваемые из Go. В них нет HTML, поэтому они используют материалы AppKit, обработку клавиатуры, анимацию и сохранение состояния без дополнительных усилий, а фронтенд может сосредоточиться на содержимом.

Все API на этой странице находятся в пакете `application` и имеют префикс `Mac`. Тот же код компилируется в Windows и Linux: конструкторы и сеттеры работают везде, вызовы подключения ничего не делают либо возвращают ошибку, а окно сохраняет свой обычный единственный WebView.

## Устройство окна

Полностью оформленное окно содержит следующие части, от переднего края к заднему:

| Часть | Тип | Класс AppKit |
|------|------|--------------|
| Панель инструментов | `MacToolbar` | `NSToolbar` |
| Боковая панель | `MacSidebar` | список источников `NSOutlineView` в элементе разделённого представления боковой панели |
| Список содержимого | `MacContentList` | `NSTableView` в элементе разделённого представления списка содержимого |
| Основное содержимое | ваш WebView или `MacTextEditor` | `WKWebView` или `NSTextView` |
| Инспектор | `MacInspector` | нативные элементы управления свойствами в элементе разделённого представления инспектора |
| Дополнительные элементы | `MacAccessory` | `NSTitlebarAccessoryViewController` или `NSSplitViewItemAccessoryViewController` |

Области располагаются с помощью `MacSplitView`, представляющего собой `NSSplitViewController`. Сначала создайте части, добавьте их в разделённое представление в порядке от переднего края к заднему, подключите разделённое представление к окну, а затем подключите панель инструментов. Для этого можно использовать `SetSplitView` и `SetToolbar` либо выполнить всё за один шаг через параметры окна `Mac.SplitView` и `Mac.Toolbar`.

Обычно вместе с нативным оформлением задают следующие параметры окна, чтобы содержимое могло прокручиваться под объединённой панелью инструментов:

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

## Панель инструментов

`NewMacToolbar` создаёт `NSToolbar`. Добавляйте элементы методами `Add`, вызывайте сеттеры и назначайте обратные вызовы через возвращённые дескрипторы, а затем подключайте панель инструментов с помощью `SetToolbar`. Идентификаторы создаются автоматически.

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

Доступны следующие виды элементов:

- `AddButton` добавляет кнопку. Для каждой кнопки до подключения панели инструментов необходимо задать `OnClick`, иначе `SetToolbar` сообщит об ошибке и оставит прежнюю панель на месте.
- `AddSearch` добавляет `NSSearchToolbarItem`. Для него необходим `OnSearch`. Используйте `SetSearchPlaceholder`, `SetSearchIncremental`, `SetSearchRecentsKey` для сохраняемого меню недавних поисковых запросов и `SetSearchMenu` для собственного меню за значком лупы.
- `AddShare` добавляет системный элемент для отправки данных (см. ниже).
- `AddGroup` добавляет сегментированную группу `NSToolbarItemGroup`. Добавляйте элементы методом группы `AddButton` и выбирайте `ToolbarGroupSelectOne`, `ToolbarGroupMomentary` или `ToolbarGroupSelectAny`.
- `AddMenu` добавляет раскрывающееся меню `NSMenuToolbarItem`, использующее обычный `Menu`. `SetShowsIndicator(false)` скрывает стрелку.
- `AddSpace` и `AddFlexibleSpace` добавляют стандартные разделители.
- `AddSidebarToggle` и `AddSidebarTrackingSeparator` добавляют элементы боковой панели AppKit. Разделитель удерживает все предшествующие ему элементы над границей боковой панели, поэтому в окне должно быть разделённое представление с областью боковой панели.
- `AddInspectorToggle` и `AddInspectorTrackingSeparator` выполняют ту же роль для области инспектора.

`SetDisplayMode` позволяет выбрать `MacToolbarDisplayModeIconAndLabel` (по умолчанию), `MacToolbarDisplayModeIconOnly`, `MacToolbarDisplayModeLabelOnly` или `MacToolbarDisplayModeDefault`.

### Обновление во время работы

Каждый дескриптор также позволяет обновлять элемент во время работы. Сеттеры, вызванные после подключения, обновляют нативный элемент в потоке приложения.

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

Другие полезные сеттеры: `SetLabel`, `SetTooltip`, `SetEnabled`, `SetHidden`, `SetTintColor` и `SetNavigational`. Последний удерживает элемент у переднего края, как Safari удерживает кнопки «Назад» и «Вперёд». `SetVisibilityPriority` определяет, какие элементы первыми перемещаются в меню переполнения при сужении окна.

### Настройка пользователем

`SetCustomizable` включает стандартную панель «Настроить панель инструментов...» и сохраняет заданное пользователем расположение под указанным вами ключом. Назначьте каждому элементу постоянный `SetPersistenceKey`, чтобы сохранённое расположение восстанавливалось после перезапуска, и используйте `SetInDefaultSet(false)` для элементов, которые должны появляться только после добавления пользователем. Вызывайте `SetCustomizable` до подключения панели инструментов.

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

### Отправка данных

`AddShare` возвращает `MacToolbarShareItem`. Элемент остаётся отключённым, пока `MacShareProvider` не объявит хотя бы одно представление данных. Wails запрашивает байты у провайдера только тогда, когда они нужны службе отправки, поэтому большие экспортируемые данные формируются по требованию. `MacShareProviderFunc` превращает две функции в провайдер; приложения с состоянием могут реализовать интерфейс напрямую.

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

Распространённые типы содержимого: `MacShareTypePlainText`, `MacShareTypeHTML`, `MacShareTypePDF`, `MacShareTypePNG` и `MacShareTypeJPEG`. Также принимается любая другая строка UTI.

### Подключение и отключение

Панель инструментов может принадлежать только одному окну одновременно. `WebviewWindow.SetToolbar` и `NativeWindow.SetToolbar` принимают панель инструментов как до создания окна, так и после; передача `nil` удаляет её и освобождает для использования в другом месте. `SetToolbar` у `WebviewWindow` сообщает о проблемах проверки через `Window.Error`, а версия для `NativeWindow` возвращает ошибку. Параметр окна `Mac.Toolbar` подключает панель инструментов при создании; он применяется после `Mac.SplitView`, чтобы отслеживающий разделитель мог найти боковую панель, по которой выравнивается.

## Разделённое представление

`MacSplitView` размещает области. Добавляйте их в порядке от переднего края к заднему. `AddSidebar`, `AddContentList` и `AddInspector` принимают размещаемую нативную модель и возвращают `MacSplitPane` для управления размером и сворачиванием. `AddPrimaryContent` размещает существующий WebView окна и возвращает `MacSplitWebviewPane`, добавляющий `SetContentLayout`.

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

`SetSplitView` работает как до создания нативного окна, так и после. Если вызвать его до создания окна, конфигурация сохраняется в очереди и устанавливается при создании. Если вызвать его после, например из обратного вызова меню или значка в области уведомлений работающего приложения, конфигурация устанавливается сразу: существующий WebView окна становится основной областью, текущая панель инструментов подключается повторно для выравнивания отслеживающих разделителей, а ожидающие дополнительные элементы подключаются.

Окно, создаваемое после `app.Run`, можно настроить одним вызовом с параметрами `Mac.SplitView` и `Mac.Toolbar`:

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

Правила компоновки:

- Компоновка должна содержать не менее двух областей и ровно одну основную область (`AddPrimaryContent` для `WebviewWindow`, `AddTextEditor` для `NativeWindow`).
- Допускается не более одного списка содержимого; он располагается после боковой панели и перед основной областью.
- После подключения разделённого представления структура областей фиксируется. Настройки областей, состояние сворачивания и содержимое боковой панели, списка и инспектора по-прежнему можно менять в любое время.
- Установленную компоновку нельзя заменить. Повторный вызов `SetSplitView` для того же окна сообщает `ErrMacSplitViewAlreadyInstalled` (через `Window.Error` у `WebviewWindow` или возвращаемым значением у `NativeWindow`) и не изменяет окно. Передача `nil` до установки удаляет ожидающую компоновку.
- Боковая панель, список, инспектор и разделённое представление могут принадлежать только одному окну одновременно.

Сеттеры областей: `SetMinimumThickness`, `SetMaximumThickness`, `SetPreferredThicknessFraction`, `SetHoldingPriority`, `SetCollapsible`, `SetCanCollapseFromWindowResize`, `SetCollapsed`, `Toggle`, `IsCollapsed` и `OnCollapsedChange`. `SetAutosaveName` сохраняет положения разделителей между запусками.

`SetContentLayout` основной области выбирает между `MacContentLayoutBelowToolbar` и `MacContentLayoutEdgeToEdge`. `MacContentLayoutAutomatic` наследует `MacWindow.ContentLayout`, который, в свою очередь, следует `TitleBar.FullSizeContent`. Компоновка от края до края позволяет AppKit применять эффект края прокрутки под панелью инструментов в macOS 26 и новее.

## Боковая панель

`MacSidebar` — нативный список источников. Он содержит корневые строки, разделы и строки с любой глубиной вложенности. Сохраните возвращённые дескрипторы, чтобы позже обновлять строки.

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

Сеттеры строк: `SetLabel`, `SetSymbol`, `SetTooltip`, `SetEnabled`, `SetHidden`, `SetBadge`, `SetAccessorySymbol`, `SetTintColor`, `SetEditable` и `SetExpanded`. `OnClick` вызывается, когда AppKit выбирает строку, `OnExpandedChange` — когда пользователь раскрывает или закрывает вложенные строки, а `OnRename` — после подтверждения переименования непосредственно в строке.

### Выбор элементов

По умолчанию можно выбрать один элемент. `SetSelectedItem` выбирает строку без вызова её `OnClick`. При включённом множественном выборе `OnClick` по-прежнему вызывается для строки, на которую нажали, а `OnSelectionChange` сообщает обо всём наборе выбранных строк.

```go
sidebar.SetAllowsMultipleSelection(true)
sidebar.OnSelectionChange(func(_ *application.Context, items []*application.MacSidebarItem) {
    for _, item := range items {
        _ = item.Section()
    }
})
```

### Контекстные меню

При щелчке правой кнопкой меню ищется в следующем порядке: собственный `SetContextMenu` строки, затем обратный вызов боковой панели `OnContextMenu`, затем резервный `SetContextMenu` боковой панели. Обратный вызов выполняется в потоке приложения, пока AppKit ожидает результата, поэтому он должен работать быстро.

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

### Изменение порядка перетаскиванием

`SetReorderable` позволяет пользователю перетаскивать строки внутри раздела, между разделами, а также в корень и из корня. Модель Go обновляется до вызова `OnMove`.

```go
sidebar.SetReorderable(true)
sidebar.OnMove(func(_ *application.Context, item *application.MacSidebarItem, section *application.MacSidebarSection, index int) {
    // section is nil when the row was dropped at the sidebar root
})
```

### Удаление

При удалении строки также удаляются вложенные в неё строки. После этого дескрипторы перестают действовать.

```go
notes.Remove(draft)
sidebar.RemoveSection(tags)
```

## Список содержимого

`MacContentList` — средняя колонка в Finder, Mail и обозревателях документов. Без столбцов он показывает строки с расширенным содержимым: заголовком, подзаголовком, значком у переднего края, дополнительной информацией у заднего края и значком с количеством.

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

Доступные стили: `MacContentListStyleAutomatic`, `MacContentListStyleInset`, `MacContentListStyleSourceList`, `MacContentListStylePlain` и `MacContentListStyleFullWidth`. Параметры отображения дополняют `SetRowHeight`, `SetAlternatingRowBackgrounds`, `SetHeaderVisible` и `SetAllowsMultipleSelection`. Для строк поддерживаются `SetHidden`, `SetEnabled`, `SetTooltip`, `InsertRow` для вставки в заданную позицию, `Remove` и `RemoveAll`.

### Столбцы

`SetColumns` переключает список в режим таблицы. После этого строки показывают значения `SetCells` — по одному на столбец. Столбцы с пометкой `Sortable` показывают индикатор сортировки; `SetSortable` включает сортировку по нажатию на заголовок. Если обратный вызов `OnSort` не задан, список сам сортирует строки по тексту столбца с помощью `SortBy`.

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

### Контекстные меню

Контекстные меню выбираются в том же порядке, что и на боковой панели: `SetContextMenu` строки, затем `OnContextMenu`, затем резервное меню списка.

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

## Инспектор

`MacInspector` — панель свойств у заднего края, состоящая из нативных элементов управления, сгруппированных по разделам. Каждый метод `Add` возвращает дескриптор `MacInspectorControl` с сеттерами и обратными вызовами, соответствующими виду элемента.

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

Виды элементов управления и их сеттеры:

| Элемент управления | Сеттеры | Обратный вызов |
|---------|---------|----------|
| `AddLabel` | `SetValue` | нет |
| `AddTextField` | `SetValue` | `OnTextChange` |
| `AddCheckbox` | `SetChecked` | `OnToggle` |
| `AddPopup`, `AddSegmented` | `SetOptions`, `SetSelectedIndex` | `OnSelectionChange` |
| `AddSlider`, `AddStepper` | `SetFloatValue`, `SetRange`, `SetStep` | `OnValueChange` |
| `AddColorWell` | `SetColor` | `OnColorChange` |
| `AddDatePicker` | `SetDate` | `OnDateChange` |
| `AddButton` | `SetLabel` | `OnClick` |

`SetLabel`, `SetTooltip`, `SetEnabled` и `SetHidden` применимы ко всем видам. Разделы могут сворачиваться; и разделы, и элементы управления можно перемещать или удалять в любое время.

```go
appearance.SetCollapsed(true)
inspector.MoveSection(statistics, 0)
appearance.Remove(priority)
inspector.RemoveSection(statistics)
```

## Дополнительные элементы

`MacAccessory` — полоса нативных элементов управления. Выберите расположение при создании: `MacAccessoryLayoutLeading` размещается рядом с кнопками окна, `MacAccessoryLayoutTrailing` — у заднего края строки заголовка, `MacAccessoryLayoutBottom` занимает всю ширину под строкой заголовка и панелью инструментов, а `MacAccessoryLayoutTop` размещается в верхней части области разделённого представления. Добавьте все элементы управления до подключения полосы.

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

Доступные элементы управления: `AddSearch`, `AddSegmented`, `AddButton`, `AddSymbolButton`, `AddMenuButton`, `AddLabel`, `AddFlexibleSpace` и, для нативных интеграций, `AddNativeView`. `AddTitlebarAccessory` есть и у `WebviewWindow`, и у `NativeWindow`; вызовы до создания окна ставятся в очередь и применяются при его создании. `Remove` отключает дополнительный элемент, чтобы его можно было подключить в другом месте, а `SetHidden` сворачивает его на месте.

### Дополнительные элементы областей

В macOS 26 и новее дополнительный элемент можно разместить сверху или снизу области разделённого представления — именно там Finder размещает поле фильтра боковой панели. Создайте дополнительный элемент с расположением `Top` или `Bottom` и подключите его к области с помощью `AddTopAccessory` или `AddBottomAccessory`.

```go
filter := application.NewMacAccessory(application.MacAccessoryLayoutTop)
filter.AddSearch("Filter").
    SetIncremental(true).
    OnSearch(func(_ *application.Context, query string) { filterNotes(query) })

if err := sidebarPane.AddTopAccessory(filter); err != nil {
    window.Error("sidebar filter: %s", err)
}
```

`SetPreferredScrollEdgeEffectStyle` выбирает оформление AppKit `Automatic`, `Soft` или `Hard` для содержимого, прокручиваемого за дополнительным элементом. Для явно заданных стилей требуется macOS 26.1; в более ранних версиях запрос передаётся обработчику ошибок окна, а автоматический стиль продолжает действовать.

## Собираем всё вместе

Эта программа создаёт окно из трёх областей: боковой панели, WebView и инспектора, а также панель инструментов, отслеживающую оба разделителя.

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

Нативное оформление и фронтенд взаимодействуют через обычные события и службы Wails. Боковая панель отправляет `note:selected`, инспектор — `note:title`, а страница подписывается на них с помощью `Events.On` среды выполнения.

## Нативные окна

@note{type="caution" title="Экспериментальная возможность"}
`NativeWindow`, `NativeWindowManager` и `MacTextEditor` являются экспериментальными в v3. API намеренно невелик и может измениться при переработке общего API окон для v4.
@end

У `NativeWindow` нет WebView. Его основное содержимое — `MacTextEditor`, то есть `NSTextView` внутри `NSScrollView`. Оно поддерживает те же типы панели инструментов, разделённого представления и дополнительных элементов, что и `WebviewWindow`. Создайте окно через `app.NativeWindow.New` или `app.NativeWindow.NewWithOptions`, а позже найдите его через `Get` или `GetByID`.

Нативное окно создаётся только после получения содержимого. Передайте разделённое представление, основная область которого добавлена с помощью `AddTextEditor`, либо через `NativeWindowOptions.SplitView`, либо вызовом `SetSplitView`. До `app.Run` компоновка ставится в очередь; в работающем приложении `SetSplitView` сразу создаёт и показывает окно (если не задан `Hidden`) и возвращает возможную ошибку создания. Окно без компоновки остаётся отложенным, а `Run` возвращает `ErrNativeWindowContentRequired`; компоновка без текстового редактора отклоняется с `ErrNativeWindowEditorRequired`. Используйте `NativeWindowOptions.Toolbar` и `NativeWindowOptions.SplitView` вместо полей `Mac.Toolbar` и `Mac.SplitView`, которые нативное окно игнорирует.

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

В работающем приложении то же окно можно создать одним вызовом, передав нативное оформление в параметрах:

```go
window := app.NativeWindow.NewWithOptions(application.NativeWindowOptions{
    Title:     "Native Notes",
    SplitView: split,
    Toolbar:   toolbar,
})
```

`MacTextEditor` предоставляет `SetText`, `Text`, `SetEditable`, `OnChange` и `Focus`. Программные вызовы `SetText` не вызывают `OnChange`, поэтому загрузка файла не помечает его как изменённый. `Text` считывает весь документ из AppKit, поэтому вызывайте его, когда нужно содержимое, а не при каждом изменении.

Два параметра помогают уменьшить объём приложения, использующего только нативный интерфейс:

- `NativeOnly: true` в `application.Options` отключает транспорт фронтенда и сервер ресурсов во время выполнения. При этом значении не создавайте `WebviewWindow`.
- Тег сборки `wails_native` полностью исключает из исполняемого файла код WebView, фронтенда и обновления и автоматически задаёт `NativeOnly`:

```sh
go build -tags wails_native .
```

Поддержка единственного экземпляра также исключена из сборки `wails_native`; при необходимости добавьте тег `wails_single_instance`.

## Вкладки окон

macOS может объединять окна во вкладки. Режим вкладок фиксируется при создании окна, поэтому задайте `Mac.TabbingMode` равным `MacWindowTabbingModePreferred` или `MacWindowTabbingModeAutomatic` для каждого окна, которое должно участвовать. `MacWindowTabbingModeDisallowed` (как и значение по умолчанию, когда параметр не задан) исключает окно из групп вкладок.

```go
first := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Document 1",
    URL:   "/",
    Mac: application.MacWindow{
        TabbingMode: application.MacWindowTabbingModePreferred,
    },
})
```

Для операций со вкладками нужны существующие нативные окна, поэтому вызывайте их из обработчика меню, метода службы или другого кода, выполняющегося после запуска `app.Run`. `TabGroup` возвращает дескриптор `MacWindowTabGroup`, который при каждом вызове находит нативную группу; дескриптор со значением nil безопасен в использовании, и каждый его метод возвращает нулевое значение.

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

`MacWindowTabGroup` предоставляет `Identifier`, `Count`, `Windows`, `NativeWindows`, `SelectedWindow`, `Select`, `SelectNative`, `SelectNext`, `SelectPrevious`, `IsTabBarVisible`, `ToggleTabBar`, `IsOverviewVisible` и `ToggleTabOverview`. `app.Window.TabGroups` перечисляет все группы, содержащие `WebviewWindow`. `AddNativeTab` добавляет `NativeWindow` в группу окна WebView, а `SetTabTooltip` задаёт текст подсказки при наведении на вкладку. Вне macOS методы вкладок возвращают `ErrMacWindowTabsUnsupported`.

## Требования к версии

Всё на этой странице относится только к macOS. В более ранних версиях возможности работают с указанными ниже ограничениями; API Go везде одинаков.

| Возможность | Минимальная версия macOS | Поведение в более ранних версиях |
|---------|---------------|-------------------------------|
| Вкладки окон | 10.12 | недоступны |
| Группы панели инструментов, элементы с рамкой, `AddMenu` | 10.15 | элементы меню отсутствуют; группы используют прежний способ отображения |
| SF Symbols (`SetSymbol` для элементов панели инструментов, боковой панели, списка и дополнительных элементов) | 11 | изображение не показывается |
| `AddSearch` в виде `NSSearchToolbarItem`, `SetNavigational`, отслеживающий разделитель боковой панели | 11 | поиск использует обычное поле поиска; разделитель отсутствует |
| Роли инспектора и списка содержимого в разделённом представлении | 11 | те же области размещаются в обычных элементах разделённого представления |
| Стили списка содержимого | 11 | игнорируются |
| `SetCenteredItems` | 13 | игнорируется |
| Переключатель и отслеживающий разделитель инспектора | 14 | Wails предоставляет нативную кнопку переключения; разделитель отсутствует |
| `SetBadgeCount`, `SetProminent`, `SetTintColor` для элементов панели инструментов | 26 | значения сохраняются и применяются, когда это становится возможно |
| Дополнительные элементы областей (`AddTopAccessory`, `AddBottomAccessory`) | 26 | `ErrMacSplitItemAccessoryUnavailable` |
| Явно заданные стили эффекта края прокрутки | 26.1 | передаётся обработчику ошибок окна; сохраняется автоматический стиль |

Для дополнительных элементов строки заголовка, разделённых представлений, боковых панелей, инспекторов и текстового редактора нет требований к версии сверх минимальной версии Wails.

## Примеры

Каждый пример — полноценное приложение, которое можно запустить:

- [`v3/examples/mac-toolbar`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-toolbar): редактор заметок с панелью инструментов, боковой панелью, списком содержимого, инспектором, провайдером отправки данных, настройкой панели инструментов и дополнительным элементом области.
- [`v3/examples/mac-titlebar-accessory`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-titlebar-accessory): дополнительные элементы строки заголовка у переднего и заднего краёв и снизу, управляющие представлением почтового ящика.
- [`v3/examples/mac-window-tabs`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-window-tabs): режимы вкладок, `AddTab`, `AddNativeTab` и API групп вкладок из меню.
- [`v3/examples/mac-native-editor`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-native-editor): текстовый редактор `NativeWindow` без WebView, собранный с тегом `wails_native`.

## Дальнейшие возможности

Нативное оформление, описанное на этой странице, — одна из сторон нативного приложения macOS. Другая сторона — его поведение: окна документов со значками файлов и каскадным размещением, меню со значками и индикаторами, меню Dock с индикатором хода выполнения, нативные предупреждения и панели, удаляемые элементы строки меню, тактильная обратная связь и речь, расширенный буфер обмена и перетаскивание из приложения, а также системные сведения о разрешениях, питании и локали. Эти API описаны в руководстве [Интеграция с платформой macOS](/guides/macos-platform-integration).
