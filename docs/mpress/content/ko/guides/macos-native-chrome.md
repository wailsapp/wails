---
title: "네이티브 macOS 창 구성 요소"
description: "Wails 창에 네이티브 AppKit 도구 막대, 사이드바, 콘텐츠 목록, 속성 패널, 보조 컨트롤 및 창 탭 구성하기"
slug: "guides/macos-native-chrome"
sourcePath: "guides/macos-native-chrome.md"
---

관련 플랫폼: macOS

Wails v3에서는 WebView를 실제 AppKit 창 구성 요소로 감쌀 수 있습니다. 도구 막대, 사이드바, 콘텐츠 목록, 속성 패널 및 제목 표시줄의 컨트롤은 Go에서 생성하는 네이티브 컨트롤입니다. HTML로 구성된 부분이 없으므로 AppKit의 재질, 키보드 처리, 애니메이션 및 상태 유지 기능을 그대로 활용할 수 있고, 프런트엔드는 콘텐츠에 집중할 수 있습니다.

이 페이지의 모든 API는 `application` 패키지에 있으며 이름이 `Mac`으로 시작합니다. 동일한 코드는 Windows와 Linux에서도 컴파일됩니다. 생성자와 설정 메서드는 모든 플랫폼에서 작동하고, 연결 호출은 아무 작업도 하지 않거나 오류를 반환하며, 창에는 일반적인 단일 WebView가 유지됩니다.

## 창 구성

모든 요소를 갖춘 창은 시작 가장자리부터 끝 가장자리까지 다음 부분으로 구성됩니다:

| 부분 | 유형 | AppKit 클래스 |
|------|------|--------------|
| 도구 막대 | `MacToolbar` | `NSToolbar` |
| 사이드바 | `MacSidebar` | 사이드바 분할 항목의 소스 목록 `NSOutlineView` |
| 콘텐츠 목록 | `MacContentList` | 콘텐츠 목록 분할 항목의 `NSTableView` |
| 기본 콘텐츠 | WebView 또는 `MacTextEditor` | `WKWebView` 또는 `NSTextView` |
| 속성 패널 | `MacInspector` | 속성 패널 분할 항목의 네이티브 속성 컨트롤 |
| 보조 컨트롤 | `MacAccessory` | `NSTitlebarAccessoryViewController` 또는 `NSSplitViewItemAccessoryViewController` |

패널은 `NSSplitViewController`인 `MacSplitView`로 배치합니다. 먼저 각 구성 요소를 만들고 시작 가장자리에서 끝 가장자리 순서로 분할 뷰에 추가한 다음, 분할 뷰를 창에 연결하고 도구 막대를 연결하세요. `SetSplitView`와 `SetToolbar`를 사용하거나 `Mac.SplitView`와 `Mac.Toolbar` 창 옵션으로 한 번에 설정할 수 있습니다.

일반적인 창 설정에서는 콘텐츠가 통합 도구 막대 아래로 스크롤될 수 있도록 다음 창 옵션을 함께 사용합니다:

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

## 도구 막대

`NewMacToolbar`는 `NSToolbar`를 생성합니다. `Add` 메서드로 항목을 추가하고 반환된 핸들에 설정 메서드와 콜백을 연결한 뒤 `SetToolbar`로 도구 막대를 연결하세요. 식별자는 자동으로 생성됩니다.

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

항목 종류는 다음과 같습니다:

- `AddButton`은 누름 버튼을 추가합니다. 도구 막대를 연결하기 전에 모든 버튼에 `OnClick`을 설정해야 합니다. 그렇지 않으면 `SetToolbar`가 오류를 보고하고 이전 도구 막대를 그대로 유지합니다.
- `AddSearch`는 `NSSearchToolbarItem`을 추가합니다. `OnSearch`가 필요합니다. `SetSearchPlaceholder`, `SetSearchIncremental`, 최근 검색어 메뉴를 유지하는 `SetSearchRecentsKey`, 돋보기 뒤에 사용자 지정 메뉴를 두는 `SetSearchMenu`를 사용하세요.
- `AddShare`는 시스템 공유 항목을 추가합니다(아래 참조).
- `AddGroup`은 분할된 `NSToolbarItemGroup`을 추가합니다. 그룹의 `AddButton`으로 구성원을 추가하고 `ToolbarGroupSelectOne`, `ToolbarGroupMomentary` 또는 `ToolbarGroupSelectAny`를 선택하세요.
- `AddMenu`는 일반 `Menu`로 구동되는 드롭다운 `NSMenuToolbarItem`을 추가합니다. `SetShowsIndicator(false)`는 펼침 표시를 숨깁니다.
- `AddSpace`와 `AddFlexibleSpace`는 표준 간격 항목을 추가합니다.
- `AddSidebarToggle`과 `AddSidebarTrackingSeparator`는 AppKit의 사이드바 항목을 추가합니다. 구분선은 그 앞의 모든 항목을 사이드바 경계선 위에 맞춰 유지하므로, 창에 사이드바 패널이 있는 분할 뷰가 있어야 합니다.
- `AddInspectorToggle`과 `AddInspectorTrackingSeparator`는 속성 패널에 대해 같은 역할을 합니다.

`SetDisplayMode`는 `MacToolbarDisplayModeIconAndLabel`(기본값), `MacToolbarDisplayModeIconOnly`, `MacToolbarDisplayModeLabelOnly`, `MacToolbarDisplayModeDefault` 중 하나를 선택합니다.

### 실시간 업데이트

모든 핸들은 실시간 업데이트 핸들이기도 합니다. 연결 후 설정 메서드를 적용하면 애플리케이션 스레드에서 네이티브 항목이 업데이트됩니다.

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

그 밖에 유용한 설정 메서드로는 `SetLabel`, `SetTooltip`, `SetEnabled`, `SetHidden`, `SetTintColor`, `SetNavigational`이 있습니다. `SetNavigational`은 Safari가 뒤로 및 앞으로 버튼을 배치하는 방식처럼 항목을 시작 가장자리에 유지합니다. `SetVisibilityPriority`는 창이 좁아질 때 어떤 항목부터 더보기 메뉴로 이동할지 결정합니다.

### 사용자 지정

`SetCustomizable`은 표준 "도구 막대 사용자화..." 시트를 활성화하고 사용자가 제공한 키로 사용자의 배치를 저장합니다. 각 항목에 안정적인 `SetPersistenceKey`를 지정하면 앱을 다시 실행해도 저장된 배치가 유지됩니다. 사용자가 추가한 후에만 나타나야 하는 항목에는 `SetInDefaultSet(false)`를 사용하세요. 도구 막대를 연결하기 전에 `SetCustomizable`을 호출하세요.

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

### 공유

`AddShare`는 `MacToolbarShareItem`을 반환합니다. `MacShareProvider`가 하나 이상의 표현 형식을 제공한다고 알릴 때까지 이 항목은 비활성화됩니다. Wails는 공유 서비스가 요청할 때만 제공자에게 바이트를 요청하므로 큰 내보내기 결과를 지연 생성할 수 있습니다. `MacShareProviderFunc`는 두 함수를 제공자로 변환하며, 상태를 관리하는 애플리케이션은 인터페이스를 직접 구현할 수 있습니다.

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

일반적인 콘텐츠 유형은 `MacShareTypePlainText`, `MacShareTypeHTML`, `MacShareTypePDF`, `MacShareTypePNG`, `MacShareTypeJPEG`입니다. 다른 UTI 문자열도 사용할 수 있습니다.

### 연결 및 분리

도구 막대는 한 번에 하나의 창에만 속할 수 있습니다. `WebviewWindow.SetToolbar`와 `NativeWindow.SetToolbar`는 창이 생성되기 전이나 후에 도구 막대를 받을 수 있습니다. `nil`을 전달하면 도구 막대를 제거하고 다른 곳에서 사용할 수 있도록 해제합니다. `WebviewWindow`의 `SetToolbar`는 검증 문제를 `Window.Error`를 통해 보고하고, `NativeWindow` 버전은 오류를 반환합니다. `Mac.Toolbar` 창 옵션은 창을 생성할 때 도구 막대를 연결합니다. 추적 구분선이 정렬할 사이드바를 찾을 수 있도록 이 옵션은 `Mac.SplitView` 다음에 적용됩니다.

## 분할 뷰

`MacSplitView`는 패널을 배치합니다. 시작 가장자리에서 끝 가장자리 순서로 추가하세요. `AddSidebar`, `AddContentList`, `AddInspector`는 각각 표시할 네이티브 모델을 받고 크기 및 접기 제어에 사용하는 `MacSplitPane`을 반환합니다. `AddPrimaryContent`는 창의 기존 WebView를 배치하고 `SetContentLayout`을 추가로 제공하는 `MacSplitWebviewPane`을 반환합니다.

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

`SetSplitView`는 네이티브 창이 생성되기 전이나 후에 작동합니다. 생성 전에 호출하면 레이아웃을 대기열에 넣고 창이 생성될 때 설치합니다. 생성 후, 예를 들어 실행 중인 애플리케이션의 메뉴나 트레이 콜백에서 호출하면 레이아웃을 즉시 설치합니다. 이때 창의 기존 WebView가 기본 패널이 되고, 추적 구분선이 맞춰지도록 현재 도구 막대를 다시 연결하며, 대기 중인 보조 컨트롤도 연결합니다.

`app.Run` 이후 생성하는 창은 `Mac.SplitView`와 `Mac.Toolbar` 옵션을 사용해 한 번의 호출로 구성할 수 있습니다:

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

레이아웃에는 다음 규칙이 적용됩니다:

- 레이아웃에는 패널이 최소 두 개 필요하고 기본 패널은 정확히 하나여야 합니다(`WebviewWindow`에는 `AddPrimaryContent`, `NativeWindow`에는 `AddTextEditor`).
- 콘텐츠 목록은 최대 하나이며 사이드바 뒤, 기본 패널 앞에 배치해야 합니다.
- 분할 뷰를 연결하면 패널 구조가 고정됩니다. 패널 설정, 접힘 상태, 사이드바·목록·속성 패널의 내용은 이후에도 언제든 변경할 수 있습니다.
- 설치된 레이아웃은 교체할 수 없습니다. 같은 창에서 `SetSplitView`를 두 번째로 호출하면 `ErrMacSplitViewAlreadyInstalled`가 보고되고(`WebviewWindow`에서는 `Window.Error`를 통해, `NativeWindow`에서는 반환값으로), 창은 변경되지 않습니다. 설치 전에 `nil`을 전달하면 대기 중인 레이아웃이 지워집니다.
- 사이드바, 목록, 속성 패널 또는 분할 뷰는 한 번에 하나의 창에만 속할 수 있습니다.

패널 설정 메서드는 `SetMinimumThickness`, `SetMaximumThickness`, `SetPreferredThicknessFraction`, `SetHoldingPriority`, `SetCollapsible`, `SetCanCollapseFromWindowResize`, `SetCollapsed`, `Toggle`, `IsCollapsed`, `OnCollapsedChange`입니다. `SetAutosaveName`은 실행 간 구분선 위치를 유지합니다.

기본 패널의 `SetContentLayout`은 `MacContentLayoutBelowToolbar`와 `MacContentLayoutEdgeToEdge` 중 하나를 선택합니다. `MacContentLayoutAutomatic`은 `MacWindow.ContentLayout`을 따르고, 이는 다시 `TitleBar.FullSizeContent`를 따릅니다. 가장자리까지 확장하는 배치를 사용하면 macOS 26 이상에서 AppKit이 도구 막대 아래에 스크롤 가장자리 효과를 적용할 수 있습니다.

## 사이드바

`MacSidebar`는 네이티브 소스 목록입니다. 최상위 행, 섹션 및 깊이 제한 없이 중첩된 행을 담습니다. 나중에 행을 업데이트할 수 있도록 반환된 핸들을 보관하세요.

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

행 설정 메서드는 `SetLabel`, `SetSymbol`, `SetTooltip`, `SetEnabled`, `SetHidden`, `SetBadge`, `SetAccessorySymbol`, `SetTintColor`, `SetEditable`, `SetExpanded`입니다. AppKit이 행을 선택하면 `OnClick`, 사용자가 중첩된 행을 펼치거나 접으면 `OnExpandedChange`, 인라인 이름 변경이 확정되면 `OnRename`이 실행됩니다.

### 선택

기본값은 단일 선택입니다. `SetSelectedItem`은 `OnClick`을 실행하지 않고 행을 선택합니다. 다중 선택이 활성화되어 있어도 클릭한 행에 대해 `OnClick`이 실행되며 `OnSelectionChange`는 전체 선택 집합을 보고합니다.

```go
sidebar.SetAllowsMultipleSelection(true)
sidebar.OnSelectionChange(func(_ *application.Context, items []*application.MacSidebarItem) {
    for _, item := range items {
        _ = item.Section()
    }
})
```

### 컨텍스트 메뉴

마우스 오른쪽 버튼을 클릭하면 다음 순서로 메뉴를 찾습니다. 행의 `SetContextMenu`, 사이드바의 `OnContextMenu` 콜백, 사이드바의 대체 `SetContextMenu` 순입니다. 콜백은 AppKit이 기다리는 동안 애플리케이션 스레드에서 실행되므로 빠르게 완료되도록 하세요.

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

### 드래그로 순서 변경

`SetReorderable`을 사용하면 사용자가 섹션 내부와 섹션 간에 행을 드래그하거나 최상위 수준으로 또는 최상위 수준에서 이동할 수 있습니다. `OnMove`가 실행되기 전에 Go 모델이 업데이트됩니다.

```go
sidebar.SetReorderable(true)
sidebar.OnMove(func(_ *application.Context, item *application.MacSidebarItem, section *application.MacSidebarSection, index int) {
    // section is nil when the row was dropped at the sidebar root
})
```

### 제거

행을 제거하면 중첩된 행도 제거됩니다. 이후 핸들은 아무 동작도 하지 않습니다.

```go
notes.Remove(draft)
sidebar.RemoveSection(tags)
```

## 콘텐츠 목록

`MacContentList`는 Finder, Mail, 문서 브라우저의 가운데 열에 해당합니다. 열이 없으면 제목, 부제목, 앞쪽 기호, 뒤쪽 세부 정보 및 개수 배지가 포함된 풍부한 행을 표시합니다.

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

스타일은 `MacContentListStyleAutomatic`, `MacContentListStyleInset`, `MacContentListStyleSourceList`, `MacContentListStylePlain`, `MacContentListStyleFullWidth`입니다. `SetRowHeight`, `SetAlternatingRowBackgrounds`, `SetHeaderVisible`, `SetAllowsMultipleSelection`으로 표시 방식을 추가로 설정할 수 있습니다. 행에는 `SetHidden`, `SetEnabled`, `SetTooltip`, 특정 위치에 삽입하는 `InsertRow`, `Remove`, `RemoveAll`을 사용할 수 있습니다.

### 열

`SetColumns`는 표 모드로 전환합니다. 그러면 행은 열마다 하나씩 `SetCells` 값을 표시합니다. `Sortable`로 표시된 열에는 정렬 표시가 나타나며, `SetSortable`은 머리글 클릭을 활성화합니다. `OnSort` 콜백이 없으면 목록이 `SortBy`를 사용해 열의 텍스트를 기준으로 자체 정렬합니다.

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

### 컨텍스트 메뉴

컨텍스트 메뉴는 사이드바와 같은 순서로 결정됩니다. 행의 `SetContextMenu`, `OnContextMenu`, 목록의 대체 메뉴 순입니다.

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

## 속성 패널

`MacInspector`는 네이티브 컨트롤을 섹션으로 묶어 만든 끝쪽 속성 패널입니다. 각 `Add` 메서드는 종류별 설정 메서드와 콜백을 제공하는 `MacInspectorControl` 핸들을 반환합니다.

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

컨트롤 종류와 설정 메서드는 다음과 같습니다:

| 컨트롤 | 설정 메서드 | 콜백 |
|---------|---------|----------|
| `AddLabel` | `SetValue` | 없음 |
| `AddTextField` | `SetValue` | `OnTextChange` |
| `AddCheckbox` | `SetChecked` | `OnToggle` |
| `AddPopup`, `AddSegmented` | `SetOptions`, `SetSelectedIndex` | `OnSelectionChange` |
| `AddSlider`, `AddStepper` | `SetFloatValue`, `SetRange`, `SetStep` | `OnValueChange` |
| `AddColorWell` | `SetColor` | `OnColorChange` |
| `AddDatePicker` | `SetDate` | `OnDateChange` |
| `AddButton` | `SetLabel` | `OnClick` |

`SetLabel`, `SetTooltip`, `SetEnabled`, `SetHidden`은 모든 종류에 적용됩니다. 섹션은 접을 수 있도록 설정할 수 있으며 섹션과 컨트롤은 언제든 이동하거나 제거할 수 있습니다.

```go
appearance.SetCollapsed(true)
inspector.MoveSection(statistics, 0)
appearance.Remove(priority)
inspector.RemoveSection(statistics)
```

## 보조 컨트롤

`MacAccessory`는 네이티브 컨트롤의 띠입니다. 생성할 때 레이아웃을 선택하세요. `MacAccessoryLayoutLeading`은 창 버튼 옆, `MacAccessoryLayoutTrailing`은 제목 표시줄의 끝 가장자리, `MacAccessoryLayoutBottom`은 제목 표시줄과 도구 막대 아래의 전체 너비에 배치됩니다. `MacAccessoryLayoutTop`은 분할 뷰 패널의 맨 위에 배치됩니다. 보조 컨트롤을 연결하기 전에 모든 컨트롤을 추가하세요.

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

사용할 수 있는 컨트롤은 `AddSearch`, `AddSegmented`, `AddButton`, `AddSymbolButton`, `AddMenuButton`, `AddLabel`, `AddFlexibleSpace`와 네이티브 통합을 위한 `AddNativeView`입니다. `AddTitlebarAccessory`는 `WebviewWindow`와 `NativeWindow` 모두에 있습니다. 창이 생성되기 전의 호출은 대기열에 넣었다가 생성 시 적용합니다. `Remove`는 보조 컨트롤을 분리하여 다른 곳에 다시 연결할 수 있게 하고, `SetHidden`은 제자리에서 접습니다.

### 패널 보조 컨트롤

macOS 26 이상에서는 보조 컨트롤을 분할 패널의 위나 아래에 배치할 수 있습니다. Finder의 사이드바 필터 필드가 이런 위치에 있습니다. `Top` 또는 `Bottom` 레이아웃으로 보조 컨트롤을 생성하고 패널의 `AddTopAccessory` 또는 `AddBottomAccessory`로 연결하세요.

```go
filter := application.NewMacAccessory(application.MacAccessoryLayoutTop)
filter.AddSearch("Filter").
    SetIncremental(true).
    OnSearch(func(_ *application.Context, query string) { filterNotes(query) })

if err := sidebarPane.AddTopAccessory(filter); err != nil {
    window.Error("sidebar filter: %s", err)
}
```

`SetPreferredScrollEdgeEffectStyle`은 보조 컨트롤 뒤로 스크롤되는 콘텐츠에 적용할 AppKit의 `Automatic`, `Soft`, `Hard` 처리 방식을 선택합니다. 명시적 스타일에는 macOS 26.1이 필요합니다. 이전 릴리스에서는 요청이 창의 오류 처리기를 통해 보고되고 자동 스타일이 계속 적용됩니다.

## 모두 조합하기

이 프로그램은 사이드바, WebView, 속성 패널로 구성된 3개 패널 창과 두 경계선을 모두 따라가는 도구 막대를 만듭니다.

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

네이티브 창 구성 요소와 프런트엔드는 일반적인 Wails 이벤트 및 서비스를 통해 통신합니다. 사이드바는 `note:selected`, 속성 패널은 `note:title`을 발생시키며 페이지는 런타임의 `Events.On`으로 이를 수신합니다.

## 네이티브 창

@note{type="caution" title="실험적 기능"}
`NativeWindow`, `NativeWindowManager`, `MacTextEditor`는 v3에서 실험적 기능입니다. API는 의도적으로 작게 설계되었으며 v4에서 공통 창 API를 다시 설계할 때 변경될 수 있습니다.
@end

`NativeWindow`에는 WebView가 없습니다. 기본 콘텐츠는 `NSScrollView` 안의 `NSTextView`인 `MacTextEditor`이며, `WebviewWindow`와 동일한 도구 막대, 분할 뷰 및 보조 컨트롤 유형을 받습니다. `app.NativeWindow.New` 또는 `app.NativeWindow.NewWithOptions`로 생성하고, 나중에 `Get` 또는 `GetByID`로 찾으세요.

네이티브 창은 콘텐츠가 있어야 생성됩니다. `NativeWindowOptions.SplitView` 또는 `SetSplitView` 호출을 통해 `AddTextEditor`로 기본 패널을 추가한 분할 뷰를 제공하세요. `app.Run` 전에는 레이아웃이 대기열에 들어갑니다. 실행 중인 애플리케이션에서는 `SetSplitView`가 즉시 창을 생성하고 표시하며(`Hidden`이 설정된 경우 제외), 생성 오류가 있으면 반환합니다. 레이아웃이 없는 창은 생성이 보류된 상태로 남고 `Run`은 `ErrNativeWindowContentRequired`를 반환합니다. 텍스트 편집기가 없는 레이아웃은 `ErrNativeWindowEditorRequired`로 거부됩니다. 네이티브 창에서는 무시되는 `Mac.Toolbar` 및 `Mac.SplitView` 필드 대신 `NativeWindowOptions.Toolbar`와 `NativeWindowOptions.SplitView`를 사용하세요.

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

실행 중인 애플리케이션에서는 창 구성 요소를 옵션으로 전달하여 같은 창을 한 번의 호출로 생성할 수 있습니다:

```go
window := app.NativeWindow.NewWithOptions(application.NativeWindowOptions{
    Title:     "Native Notes",
    SplitView: split,
    Toolbar:   toolbar,
})
```

`MacTextEditor`는 `SetText`, `Text`, `SetEditable`, `OnChange`, `Focus`를 제공합니다. 코드에서 호출한 `SetText`는 `OnChange`를 실행하지 않으므로 파일을 불러와도 수정된 것으로 표시되지 않습니다. `Text`는 AppKit에서 문서 전체를 읽으므로 변경 때마다 호출하지 말고 내용이 필요할 때 호출하세요.

다음 두 옵션으로 네이티브 전용 애플리케이션을 간결하게 유지할 수 있습니다:

- `application.Options`의 `NativeOnly: true`는 런타임에 프런트엔드 전송 계층과 에셋 서버를 건너뜁니다. 이 옵션이 설정되어 있으면 `WebviewWindow`를 생성하지 마세요.
- `wails_native` 빌드 태그는 WebView, 프런트엔드 및 업데이터 코드를 바이너리에서 완전히 제외하고 `NativeOnly`를 자동으로 설정합니다:

```sh
go build -tags wails_native .
```

`wails_native` 빌드에서는 단일 인스턴스 지원도 제외됩니다. 필요하다면 `wails_single_instance` 태그를 추가하세요.

## 창 탭

macOS에서는 창을 탭으로 묶을 수 있습니다. 탭 모드는 창을 생성할 때 고정되므로 참여할 모든 창에서 `Mac.TabbingMode`를 `MacWindowTabbingModePreferred` 또는 `MacWindowTabbingModeAutomatic`으로 설정하세요. `MacWindowTabbingModeDisallowed`와 설정되지 않은 기본값은 창을 탭 그룹에서 제외합니다.

```go
first := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Document 1",
    URL:   "/",
    Mac: application.MacWindow{
        TabbingMode: application.MacWindowTabbingModePreferred,
    },
})
```

탭 작업에는 생성된 네이티브 창이 필요합니다. 따라서 `app.Run`이 시작된 후 실행되는 메뉴 처리기, 서비스 메서드 또는 다른 코드에서 호출하세요. `TabGroup`은 호출할 때마다 네이티브 그룹을 확인하는 `MacWindowTabGroup` 핸들을 반환합니다. nil 핸들도 안전하게 사용할 수 있으며, 모든 메서드는 해당하는 0 값을 반환합니다.

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

`MacWindowTabGroup`은 `Identifier`, `Count`, `Windows`, `NativeWindows`, `SelectedWindow`, `Select`, `SelectNative`, `SelectNext`, `SelectPrevious`, `IsTabBarVisible`, `ToggleTabBar`, `IsOverviewVisible`, `ToggleTabOverview`를 제공합니다. `app.Window.TabGroups`는 `WebviewWindow`가 포함된 모든 그룹을 나열합니다. `AddNativeTab`은 `NativeWindow`를 WebView 창의 그룹에 추가하고, `SetTabTooltip`은 탭에 마우스를 올렸을 때 표시할 텍스트를 설정합니다. macOS 이외의 플랫폼에서는 탭 메서드가 `ErrMacWindowTabsUnsupported`를 반환합니다.

## 버전 요구 사항

이 페이지의 모든 기능은 macOS 전용입니다. 이전 릴리스에서는 아래와 같이 기능이 적절히 제한되며 Go API는 모든 플랫폼에서 동일합니다.

| 기능 | 최소 macOS 버전 | 이전 릴리스에서의 동작 |
|---------|---------------|-------------------------------|
| 창 탭 | 10.12 | 사용할 수 없음 |
| 도구 막대 그룹, 테두리가 있는 항목, `AddMenu` | 10.15 | 메뉴 항목은 생략되고 그룹은 기존 표시 방식을 사용함 |
| SF Symbols(도구 막대, 사이드바, 목록 및 보조 컨트롤 항목의 `SetSymbol`) | 11 | 이미지가 표시되지 않음 |
| `NSSearchToolbarItem`으로서의 `AddSearch`, `SetNavigational`, 사이드바 추적 구분선 | 11 | 검색은 일반 검색 필드로 대체되고 구분선은 생략됨 |
| 속성 패널 및 콘텐츠 목록의 분할 역할 | 11 | 동일한 패널이 일반 분할 항목에 표시됨 |
| 콘텐츠 목록 스타일 | 11 | 무시됨 |
| `SetCenteredItems` | 13 | 무시됨 |
| 속성 패널 토글 및 추적 구분선 | 14 | Wails가 네이티브 토글 버튼을 제공하고 구분선은 생략됨 |
| 도구 막대 항목의 `SetBadgeCount`, `SetProminent`, `SetTintColor` | 26 | 저장되었다가 사용할 수 있을 때 적용됨 |
| 패널 보조 컨트롤(`AddTopAccessory`, `AddBottomAccessory`) | 26 | `ErrMacSplitItemAccessoryUnavailable` |
| 명시적 스크롤 가장자리 효과 스타일 | 26.1 | 창 오류 처리기를 통해 보고되고 자동 스타일이 유지됨 |

제목 표시줄 보조 컨트롤, 분할 뷰, 사이드바, 속성 패널 및 텍스트 편집기에는 Wails 최소 버전 외의 추가 버전 요구 사항이 없습니다.

## 예제

각 예제는 완전하게 실행할 수 있는 애플리케이션입니다:

- [`v3/examples/mac-toolbar`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-toolbar): 도구 막대, 사이드바, 콘텐츠 목록, 속성 패널, 공유 제공자, 도구 막대 사용자화 및 패널 보조 컨트롤이 있는 메모 편집기입니다.
- [`v3/examples/mac-titlebar-accessory`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-titlebar-accessory): 편지함 뷰를 제어하는 제목 표시줄의 시작·끝·아래쪽 보조 컨트롤입니다.
- [`v3/examples/mac-window-tabs`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-window-tabs): 메뉴에서 사용하는 탭 모드, `AddTab`, `AddNativeTab` 및 탭 그룹 API입니다.
- [`v3/examples/mac-native-editor`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-native-editor): `wails_native` 태그로 빌드하는 WebView 없는 `NativeWindow` 텍스트 편집기입니다.

## 더 알아보기

이 페이지의 창 구성 요소는 네이티브 macOS 애플리케이션의 한 부분입니다. 다른 부분은 애플리케이션의 동작입니다. 여기에는 프록시 아이콘과 계단식 배치를 갖춘 문서 창, 기호와 배지가 있는 메뉴, 진행률을 표시하는 Dock 메뉴, 네이티브 경고와 패널, 제거 가능한 상태 항목, 햅틱과 음성, 다양한 형식을 지원하는 클립보드와 외부로 드래그하기, 권한·전원·로캘에 관한 시스템 정보가 포함됩니다. 이러한 API는 [macOS 플랫폼 통합](/guides/macos-platform-integration) 가이드에서 다룹니다.
