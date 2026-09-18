---
title: Native macOS Chrome
description: Build native AppKit toolbars, sidebars, content lists, inspectors, accessories and window tabs around your Wails window
slug: "guides/macos-native-chrome"
sourcePath: "guides/macos-native-chrome.md"
---

Relevant Platforms: macOS

Wails v3 can wrap your WebView in real AppKit window chrome. The toolbar,
sidebar, content list, inspector and titlebar strips are native controls
created from Go. Nothing in them is HTML, so they get AppKit's materials,
keyboard handling, animation and persistence for free, and your frontend
stays focused on the content.

Every API on this page lives in the `application` package and is prefixed
with `Mac`. The same code compiles on Windows and Linux: constructors and
setters work everywhere, attach calls are no-ops or return an error, and the
window keeps its ordinary single WebView.

## Window anatomy

A fully dressed window has these parts, from leading edge to trailing edge:

| Part | Type | AppKit class |
|------|------|--------------|
| Toolbar | `MacToolbar` | `NSToolbar` |
| Sidebar | `MacSidebar` | `NSOutlineView` source list in a sidebar split item |
| Content list | `MacContentList` | `NSTableView` in a content-list split item |
| Primary content | your WebView, or `MacTextEditor` | `WKWebView` or `NSTextView` |
| Inspector | `MacInspector` | native property controls in an inspector split item |
| Accessories | `MacAccessory` | `NSTitlebarAccessoryViewController` or `NSSplitViewItemAccessoryViewController` |

The panes are arranged by a `MacSplitView`, which is an
`NSSplitViewController`. Build the pieces first, add them to the split view
in leading-to-trailing order, attach the split view to the window, and then
attach the toolbar. You can do that with `SetSplitView` and `SetToolbar`, or
in one step through the `Mac.SplitView` and `Mac.Toolbar` window options.

A typical window setup pairs the chrome with these window options so content
can scroll beneath a unified toolbar:

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

## Toolbar

`NewMacToolbar` creates an `NSToolbar`. Add items with the `Add` methods,
chain setters and callbacks on the returned handles, and attach the toolbar
with `SetToolbar`. Identifiers are generated for you.

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

The item kinds are:

- `AddButton` adds a push button. Every button needs an `OnClick` before the
  toolbar is attached, otherwise `SetToolbar` reports an error and leaves the
  previous toolbar in place.
- `AddSearch` adds an `NSSearchToolbarItem`. `OnSearch` is required. Use
  `SetSearchPlaceholder`, `SetSearchIncremental`, `SetSearchRecentsKey` for a
  persistent recent-searches menu, and `SetSearchMenu` for a custom menu
  behind the magnifier.
- `AddShare` adds the system share item (see below).
- `AddGroup` adds a segmented `NSToolbarItemGroup`. Add members with the
  group's `AddButton` and choose `ToolbarGroupSelectOne`,
  `ToolbarGroupMomentary` or `ToolbarGroupSelectAny`.
- `AddMenu` adds a dropdown `NSMenuToolbarItem` driven by an ordinary
  `Menu`. `SetShowsIndicator(false)` hides the chevron.
- `AddSpace` and `AddFlexibleSpace` add the standard spacers.
- `AddSidebarToggle` and `AddSidebarTrackingSeparator` add AppKit's sidebar
  items. The separator keeps everything before it aligned above the sidebar
  divider, so the window must have a split view with a sidebar pane.
- `AddInspectorToggle` and `AddInspectorTrackingSeparator` do the same for
  an inspector pane.

`SetDisplayMode` chooses between `MacToolbarDisplayModeIconAndLabel` (the
default), `MacToolbarDisplayModeIconOnly`, `MacToolbarDisplayModeLabelOnly`
and `MacToolbarDisplayModeDefault`.

### Live updates

Every handle is also a live update handle. Setters applied after attachment
update the native item on the application thread.

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

Other useful setters are `SetLabel`, `SetTooltip`, `SetEnabled`,
`SetHidden`, `SetTintColor` and `SetNavigational`, which keeps an item at the
leading edge in the way Safari keeps its back and forward buttons.
`SetVisibilityPriority` decides which items move into the overflow menu first
when the window narrows.

### User customisation

`SetCustomizable` enables the standard "Customize Toolbar..." sheet and
saves the user's layout under the key you supply. Give each item a stable
`SetPersistenceKey` so the saved layout survives a relaunch, and use
`SetInDefaultSet(false)` for items that should only appear once the user adds
them. Call `SetCustomizable` before the toolbar is attached.

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

### Sharing

`AddShare` returns a `MacToolbarShareItem`. It stays disabled until a
`MacShareProvider` advertises at least one representation. Wails asks the
provider for bytes only when a sharing service requests them, so large
exports are rendered lazily. `MacShareProviderFunc` adapts two functions into
a provider; stateful applications can implement the interface directly.

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

The common content types are `MacShareTypePlainText`, `MacShareTypeHTML`,
`MacShareTypePDF`, `MacShareTypePNG` and `MacShareTypeJPEG`. Any other UTI
string is accepted as well.

### Attaching and detaching

A toolbar belongs to one window at a time. `WebviewWindow.SetToolbar` and
`NativeWindow.SetToolbar` accept a toolbar before or after the window exists;
passing `nil` removes it and releases it for use elsewhere. `SetToolbar` on a
`WebviewWindow` reports validation problems through `Window.Error`, while the
`NativeWindow` version returns the error. The `Mac.Toolbar` window option
attaches a toolbar at creation time; it is applied after `Mac.SplitView` so a
tracking separator can find the sidebar it aligns with.

## Split view

`MacSplitView` arranges the panes. Add them in leading-to-trailing order.
`AddSidebar`, `AddContentList` and `AddInspector` take the native model they
host and return a `MacSplitPane` for sizing and collapse control.
`AddPrimaryContent` places the window's existing WebView and returns a
`MacSplitWebviewPane`, which adds `SetContentLayout`.

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

`SetSplitView` works before or after the native window exists. Called
before, the layout is queued and installed when the window is created.
Called afterwards, for example from a menu or tray callback in a running
application, the layout is installed immediately: the window's existing
WebView becomes the primary pane, the current toolbar is re-attached so its
tracking separators line up, and any queued accessories are attached.

A window created after `app.Run` can be configured in one call with the
`Mac.SplitView` and `Mac.Toolbar` options:

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

The rules for a layout are:

- A layout needs at least two panes and exactly one primary pane
  (`AddPrimaryContent` for a `WebviewWindow`, `AddTextEditor` for a
  `NativeWindow`).
- At most one content list, placed after the sidebar and before the primary
  pane.
- The pane structure is frozen once the split view is attached. Pane
  settings, collapse state and the contents of the sidebar, list and
  inspector can still change at any time.
- An installed layout cannot be replaced. A second `SetSplitView` on the
  same window reports `ErrMacSplitViewAlreadyInstalled` (through
  `Window.Error` on a `WebviewWindow`, as the return value on a
  `NativeWindow`) and leaves the window unchanged. Passing `nil` before
  installation clears a pending layout.
- A sidebar, list, inspector or split view belongs to one window at a time.

Pane setters are `SetMinimumThickness`, `SetMaximumThickness`,
`SetPreferredThicknessFraction`, `SetHoldingPriority`, `SetCollapsible`,
`SetCanCollapseFromWindowResize`, `SetCollapsed`, `Toggle`, `IsCollapsed` and
`OnCollapsedChange`. `SetAutosaveName` persists divider positions between
launches.

`SetContentLayout` on the primary pane chooses between
`MacContentLayoutBelowToolbar` and `MacContentLayoutEdgeToEdge`.
`MacContentLayoutAutomatic` inherits `MacWindow.ContentLayout`, which in turn
follows `TitleBar.FullSizeContent`. Edge-to-edge is the arrangement that lets
AppKit apply its scroll-edge effect beneath the toolbar on macOS 26 and newer.

## Sidebar

`MacSidebar` is a native source list. It holds root rows, sections and rows
nested to any depth. Keep the returned handles to update rows later.

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

Row setters are `SetLabel`, `SetSymbol`, `SetTooltip`, `SetEnabled`,
`SetHidden`, `SetBadge`, `SetAccessorySymbol`, `SetTintColor`,
`SetEditable` and `SetExpanded`. `OnClick` fires when AppKit selects the row,
`OnExpandedChange` when the user opens or closes its nested rows and
`OnRename` after an inline rename commits.

### Selection

Single selection is the default. `SetSelectedItem` selects a row without
firing its `OnClick`. With multiple selection enabled, `OnClick` still fires
for the clicked row and `OnSelectionChange` reports the whole set.

```go
sidebar.SetAllowsMultipleSelection(true)
sidebar.OnSelectionChange(func(_ *application.Context, items []*application.MacSidebarItem) {
    for _, item := range items {
        _ = item.Section()
    }
})
```

### Context menus

A right-click looks for a menu in this order: the row's own `SetContextMenu`,
then the sidebar's `OnContextMenu` callback, then the sidebar's fallback
`SetContextMenu`. The callback runs on the application thread while AppKit
waits, so keep it quick.

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

### Drag reorder

`SetReorderable` lets the user drag rows within a section, between sections
and to or from the root. The Go model is updated before `OnMove` fires.

```go
sidebar.SetReorderable(true)
sidebar.OnMove(func(_ *application.Context, item *application.MacSidebarItem, section *application.MacSidebarSection, index int) {
    // section is nil when the row was dropped at the sidebar root
})
```

### Removal

Removing a row also removes its nested rows. Handles are inert afterwards.

```go
notes.Remove(draft)
sidebar.RemoveSection(tags)
```

## Content list

`MacContentList` is the middle column of Finder, Mail and document browsers.
Without columns it shows rich rows: a title, a subtitle, a leading symbol, a
trailing detail and a count badge.

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

Styles are `MacContentListStyleAutomatic`, `MacContentListStyleInset`,
`MacContentListStyleSourceList`, `MacContentListStylePlain` and
`MacContentListStyleFullWidth`. `SetRowHeight`,
`SetAlternatingRowBackgrounds`, `SetHeaderVisible` and
`SetAllowsMultipleSelection` complete the presentation options. Rows support
`SetHidden`, `SetEnabled`, `SetTooltip`, `InsertRow` at a position, `Remove`
and `RemoveAll`.

### Columns

`SetColumns` switches to table mode. Rows then show their `SetCells` values,
one per column. Columns marked `Sortable` show a sort indicator; `SetSortable`
turns header clicks on. Without an `OnSort` callback the list sorts itself by
the column text with `SortBy`.

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

### Context menus

Context menus resolve in the same order as the sidebar: the row's
`SetContextMenu`, then `OnContextMenu`, then the list's fallback menu.

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

## Inspector

`MacInspector` is a trailing property panel built from native controls,
grouped into sections. Each `Add` method returns a `MacInspectorControl`
handle with kind-specific setters and callbacks.

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

The control kinds and their setters:

| Control | Setters | Callback |
|---------|---------|----------|
| `AddLabel` | `SetValue` | none |
| `AddTextField` | `SetValue` | `OnTextChange` |
| `AddCheckbox` | `SetChecked` | `OnToggle` |
| `AddPopup`, `AddSegmented` | `SetOptions`, `SetSelectedIndex` | `OnSelectionChange` |
| `AddSlider`, `AddStepper` | `SetFloatValue`, `SetRange`, `SetStep` | `OnValueChange` |
| `AddColorWell` | `SetColor` | `OnColorChange` |
| `AddDatePicker` | `SetDate` | `OnDateChange` |
| `AddButton` | `SetLabel` | `OnClick` |

`SetLabel`, `SetTooltip`, `SetEnabled` and `SetHidden` apply to every kind.
Sections can be collapsible, and both sections and controls can be moved or
removed at any time.

```go
appearance.SetCollapsed(true)
inspector.MoveSection(statistics, 0)
appearance.Remove(priority)
inspector.RemoveSection(statistics)
```

## Accessories

`MacAccessory` is a strip of native controls. Choose a layout when you create
it: `MacAccessoryLayoutLeading` sits next to the window buttons,
`MacAccessoryLayoutTrailing` at the trailing edge of the titlebar,
`MacAccessoryLayoutBottom` spans the full width beneath the titlebar and
toolbar, and `MacAccessoryLayoutTop` is the top of a split-view pane. Add
every control before attaching the accessory.

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

Controls are `AddSearch`, `AddSegmented`, `AddButton`, `AddSymbolButton`,
`AddMenuButton`, `AddLabel`, `AddFlexibleSpace` and, for native
integrations, `AddNativeView`. `AddTitlebarAccessory` exists on both
`WebviewWindow` and `NativeWindow`; calls made before the window exists are
queued and applied when it is created. `Remove` detaches an accessory so it
can be attached again elsewhere, and `SetHidden` collapses it in place.

### Pane accessories

On macOS 26 and newer an accessory can sit at the top or bottom of a split
pane, which is where Finder keeps its sidebar filter field. Create the
accessory with the `Top` or `Bottom` layout and attach it with
`AddTopAccessory` or `AddBottomAccessory` on the pane.

```go
filter := application.NewMacAccessory(application.MacAccessoryLayoutTop)
filter.AddSearch("Filter").
    SetIncremental(true).
    OnSearch(func(_ *application.Context, query string) { filterNotes(query) })

if err := sidebarPane.AddTopAccessory(filter); err != nil {
    window.Error("sidebar filter: %s", err)
}
```

`SetPreferredScrollEdgeEffectStyle` picks AppKit's `Automatic`, `Soft` or
`Hard` treatment for content scrolling behind the accessory. Explicit styles
need macOS 26.1; on earlier releases the request is reported through the
window's error handler and the automatic style remains in effect.

## Putting it together

This program builds a three-pane window with a sidebar, the WebView and an
inspector, and a toolbar that tracks both dividers.

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

The native chrome and the frontend talk through the usual Wails events and
services. The sidebar emits `note:selected`, the inspector emits
`note:title`, and the page listens with the runtime's `Events.On`.

## Native windows

@note{type="caution" title="Experimental"}
`NativeWindow`, `NativeWindowManager` and `MacTextEditor` are experimental in
v3. The API is deliberately small and may change when the common window API
is redesigned for v4.
@end

A `NativeWindow` has no WebView. Its primary content is a `MacTextEditor`,
an `NSTextView` inside an `NSScrollView`, and it accepts the same toolbar,
split view and accessory types as a `WebviewWindow`. Create one with
`app.NativeWindow.New` or `app.NativeWindow.NewWithOptions`, and look it up
later with `Get` or `GetByID`.

A native window is only created once it has content. Supply a split view
whose primary pane was added with `AddTextEditor`, either through
`NativeWindowOptions.SplitView` or by calling `SetSplitView`. Before
`app.Run` the layout is queued; in a running application `SetSplitView`
creates and shows the window immediately (unless `Hidden` is set) and
returns any creation error. A window without a layout stays deferred and
`Run` returns `ErrNativeWindowContentRequired`; a layout without a text
editor is rejected with `ErrNativeWindowEditorRequired`. Use
`NativeWindowOptions.Toolbar` and `NativeWindowOptions.SplitView` rather than
the `Mac.Toolbar` and `Mac.SplitView` fields, which a native window ignores.

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

The same window can be created in one call from a running application by
passing the chrome as options:

```go
window := app.NativeWindow.NewWithOptions(application.NativeWindowOptions{
    Title:     "Native Notes",
    SplitView: split,
    Toolbar:   toolbar,
})
```

`MacTextEditor` offers `SetText`, `Text`, `SetEditable`, `OnChange` and
`Focus`. Programmatic `SetText` calls do not fire `OnChange`, so loading a
file never marks it dirty. `Text` reads the complete document from AppKit, so
call it when you need the contents rather than on every change.

Two options keep a native-only application lean:

- `NativeOnly: true` in `application.Options` skips the frontend transport
  and asset server at runtime. Do not create a `WebviewWindow` when it is set.
- The `wails_native` build tag compiles the WebView, frontend and updater
  code out of the binary entirely and sets `NativeOnly` automatically:

```sh
go build -tags wails_native .
```

Single-instance support is also excluded from a `wails_native` build; add the
`wails_single_instance` tag when you need it.

## Window tabs

macOS can group windows into tabs. Tabbing mode is fixed when a window is
created, so set `Mac.TabbingMode` to `MacWindowTabbingModePreferred` or
`MacWindowTabbingModeAutomatic` on every window that should take part.
`MacWindowTabbingModeDisallowed` (and the unset default) keeps a window out
of tab groups.

```go
first := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Document 1",
    URL:   "/",
    Mac: application.MacWindow{
        TabbingMode: application.MacWindowTabbingModePreferred,
    },
})
```

The tab operations need live native windows, so call them from a menu
handler, a service method or any other code that runs after `app.Run` has
started. `TabGroup` returns a `MacWindowTabGroup` handle that resolves the
native group on every call; a nil handle is safe to use and every method on
it returns its zero value.

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

`MacWindowTabGroup` provides `Identifier`, `Count`, `Windows`,
`NativeWindows`, `SelectedWindow`, `Select`, `SelectNative`, `SelectNext`,
`SelectPrevious`, `IsTabBarVisible`, `ToggleTabBar`, `IsOverviewVisible` and
`ToggleTabOverview`. `app.Window.TabGroups` lists every group that contains a
`WebviewWindow`. `AddNativeTab` adds a `NativeWindow` to a WebView window's
group, and `SetTabTooltip` sets the hover text on a tab. Off macOS the tab
methods return `ErrMacWindowTabsUnsupported`.

## Version requirements

Everything on this page is macOS only. Features degrade gracefully on older
releases as noted below; the Go API is identical everywhere.

| Feature | Minimum macOS | Behaviour on earlier releases |
|---------|---------------|-------------------------------|
| Window tabs | 10.12 | not available |
| Toolbar groups, bordered items, `AddMenu` | 10.15 | menu items are omitted; groups use the legacy presentation |
| SF Symbols (`SetSymbol` on toolbar, sidebar, list and accessory items) | 11 | no image is shown |
| `AddSearch` as `NSSearchToolbarItem`, `SetNavigational`, sidebar tracking separator | 11 | search falls back to a plain search field; the separator is omitted |
| Inspector and content-list split roles | 11 | the same panes are hosted in regular split items |
| Content list styles | 11 | ignored |
| `SetCenteredItems` | 13 | ignored |
| Inspector toggle and tracking separator | 14 | Wails supplies a native toggle button; the separator is omitted |
| `SetBadgeCount`, `SetProminent`, `SetTintColor` on toolbar items | 26 | stored and applied when available |
| Pane accessories (`AddTopAccessory`, `AddBottomAccessory`) | 26 | `ErrMacSplitItemAccessoryUnavailable` |
| Explicit scroll-edge effect styles | 26.1 | reported through the window error handler; automatic style stays |

Titlebar accessories, split views, sidebars, inspectors and the text editor
have no version requirement beyond the Wails minimum.

## Examples

Each example is a complete, runnable application:

- [`v3/examples/mac-toolbar`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-toolbar):
  a notes editor with a toolbar, sidebar, content list, inspector, share
  provider, toolbar customisation and a pane accessory.
- [`v3/examples/mac-titlebar-accessory`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-titlebar-accessory):
  leading, trailing and bottom titlebar accessories driving a mailbox view.
- [`v3/examples/mac-window-tabs`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-window-tabs):
  tabbing modes, `AddTab`, `AddNativeTab` and the tab group API from a menu.
- [`v3/examples/mac-native-editor`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-native-editor):
  a WebView-free `NativeWindow` text editor built with the `wails_native` tag.

## Going further

The chrome on this page is one half of a native macOS application. The
other half is how the app behaves: document windows with proxy icons and
cascading, menus with symbols and badges, a Dock menu with progress, native
alerts and panels, removable status items, haptics and speech, the rich
clipboard and drag out, and system information about permissions, power
and locale. Those APIs are covered in the
[macOS Platform Integration](/guides/macos-platform-integration) guide.
