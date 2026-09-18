---
title: "Native macOS-Fensterelemente"
description: "Erstellen Sie native AppKit-Symbolleisten, Seitenleisten, Inhaltslisten, Inspektoren, Zusatzleisten und Fenstertabs für Ihr Wails-Fenster"
slug: "guides/macos-native-chrome"
sourcePath: "guides/macos-native-chrome.md"
---

Relevante Plattformen: macOS

Wails v3 kann Ihre WebView mit echten AppKit-Fensterelementen umgeben. Symbolleiste, Seitenleiste, Inhaltsliste, Inspektor und Titelleisten sind native Steuerelemente, die aus Go heraus erstellt werden. Sie enthalten kein HTML und nutzen dadurch die Materialien, Tastatursteuerung, Animationen und Zustandsspeicherung von AppKit. Ihr Frontend kann sich auf den Inhalt konzentrieren.

Alle APIs auf dieser Seite befinden sich im Paket `application` und tragen das Präfix `Mac`. Derselbe Code lässt sich unter Windows und Linux kompilieren: Konstruktoren und Setter funktionieren überall, Aufrufe zum Anhängen bewirken nichts oder geben einen Fehler zurück, und das Fenster behält seine gewöhnliche einzelne WebView.

## Aufbau des Fensters

Ein vollständig ausgestattetes Fenster besteht vom vorderen bis zum hinteren Rand aus diesen Teilen:

| Teil | Typ | AppKit-Klasse |
|------|------|--------------|
| Symbolleiste | `MacToolbar` | `NSToolbar` |
| Seitenleiste | `MacSidebar` | `NSOutlineView` als Quellliste in einem Seitenleisten-Teilbereich |
| Inhaltsliste | `MacContentList` | `NSTableView` in einem Inhaltslisten-Teilbereich |
| Hauptinhalt | Ihre WebView oder `MacTextEditor` | `WKWebView` oder `NSTextView` |
| Inspektor | `MacInspector` | native Eigenschaftssteuerelemente in einem Inspektor-Teilbereich |
| Zusatzleisten | `MacAccessory` | `NSTitlebarAccessoryViewController` oder `NSSplitViewItemAccessoryViewController` |

Die Teilbereiche werden von einer `MacSplitView` angeordnet, einem `NSSplitViewController`. Erstellen Sie zuerst die Komponenten, fügen Sie sie der Split-View in der Reihenfolge vom vorderen zum hinteren Rand hinzu, hängen Sie die Split-View an das Fenster und anschließend die Symbolleiste. Dafür können Sie `SetSplitView` und `SetToolbar` verwenden oder beides in einem Schritt über die Fensteroptionen `Mac.SplitView` und `Mac.Toolbar` festlegen.

Bei einer typischen Fensterkonfiguration werden die nativen Fensterelemente mit diesen Fensteroptionen kombiniert, damit Inhalt unter einer einheitlichen Symbolleiste hindurchscrollen kann:

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

## Symbolleiste

`NewMacToolbar` erstellt eine `NSToolbar`. Fügen Sie mit den `Add`-Methoden Elemente hinzu, verketten Sie Setter und Callbacks auf den zurückgegebenen Handles und hängen Sie die Symbolleiste mit `SetToolbar` an. Kennungen werden automatisch erzeugt.

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

Es gibt folgende Elementtypen:

- `AddButton` fügt eine Schaltfläche hinzu. Für jede Schaltfläche muss vor dem Anhängen der Symbolleiste ein `OnClick` festgelegt sein. Andernfalls meldet `SetToolbar` einen Fehler und die bisherige Symbolleiste bleibt erhalten.
- `AddSearch` fügt ein `NSSearchToolbarItem` hinzu. `OnSearch` ist erforderlich. Verwenden Sie `SetSearchPlaceholder`, `SetSearchIncremental`, `SetSearchRecentsKey` für ein dauerhaft gespeichertes Menü zuletzt verwendeter Suchanfragen und `SetSearchMenu` für ein eigenes Menü hinter dem Lupensymbol.
- `AddShare` fügt das systemeigene Teilen-Element hinzu (siehe unten).
- `AddGroup` fügt eine segmentierte `NSToolbarItemGroup` hinzu. Fügen Sie mit `AddButton` der Gruppe Mitglieder hinzu und wählen Sie `ToolbarGroupSelectOne`, `ToolbarGroupMomentary` oder `ToolbarGroupSelectAny`.
- `AddMenu` fügt ein Dropdown-`NSMenuToolbarItem` hinzu, das von einem gewöhnlichen `Menu` gesteuert wird. `SetShowsIndicator(false)` blendet den Pfeil aus.
- `AddSpace` und `AddFlexibleSpace` fügen die Standardabstände hinzu.
- `AddSidebarToggle` und `AddSidebarTrackingSeparator` fügen die Seitenleistenelemente von AppKit hinzu. Der Trenner hält alle Elemente davor über der Seitenleisten-Trennlinie ausgerichtet; das Fenster muss daher eine Split-View mit einem Seitenleisten-Teilbereich haben.
- `AddInspectorToggle` und `AddInspectorTrackingSeparator` erfüllen dieselbe Funktion für einen Inspektor-Teilbereich.

`SetDisplayMode` wählt zwischen `MacToolbarDisplayModeIconAndLabel` (Standard), `MacToolbarDisplayModeIconOnly`, `MacToolbarDisplayModeLabelOnly` und `MacToolbarDisplayModeDefault`.

### Änderungen zur Laufzeit

Jedes Handle dient auch zum Ändern zur Laufzeit. Nach dem Anhängen aufgerufene Setter aktualisieren das native Element im Anwendungsthread.

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

Weitere nützliche Setter sind `SetLabel`, `SetTooltip`, `SetEnabled`, `SetHidden`, `SetTintColor` und `SetNavigational`. Letzterer hält ein Element am vorderen Rand, wie Safari seine Zurück- und Vorwärts-Schaltflächen. `SetVisibilityPriority` bestimmt, welche Elemente zuerst ins Überlaufmenü verschoben werden, wenn das Fenster schmaler wird.

### Anpassung durch Benutzer

`SetCustomizable` aktiviert das standardmäßige Dialogblatt „Symbolleiste anpassen …“ und speichert die Anordnung des Benutzers unter dem angegebenen Schlüssel. Geben Sie jedem Element mit `SetPersistenceKey` einen stabilen Schlüssel, damit die gespeicherte Anordnung einen Neustart übersteht. Verwenden Sie `SetInDefaultSet(false)` für Elemente, die erst erscheinen sollen, nachdem der Benutzer sie hinzugefügt hat. Rufen Sie `SetCustomizable` auf, bevor die Symbolleiste angehängt wird.

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

### Teilen

`AddShare` gibt ein `MacToolbarShareItem` zurück. Es bleibt deaktiviert, bis ein `MacShareProvider` mindestens eine Repräsentation anbietet. Wails fordert die Bytes erst dann vom Provider an, wenn ein Teilen-Dienst sie benötigt. Dadurch werden große Exporte erst bei Bedarf erzeugt. `MacShareProviderFunc` passt zwei Funktionen als Provider an; Anwendungen mit eigenem Zustand können die Schnittstelle direkt implementieren.

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

Übliche Inhaltstypen sind `MacShareTypePlainText`, `MacShareTypeHTML`, `MacShareTypePDF`, `MacShareTypePNG` und `MacShareTypeJPEG`. Auch jede andere UTI-Zeichenfolge wird akzeptiert.

### Anhängen und Entfernen

Eine Symbolleiste gehört jeweils nur zu einem Fenster. `WebviewWindow.SetToolbar` und `NativeWindow.SetToolbar` akzeptieren eine Symbolleiste vor oder nach der Erstellung des Fensters. Mit `nil` wird sie entfernt und für die Verwendung an anderer Stelle freigegeben. `SetToolbar` auf einem `WebviewWindow` meldet Validierungsfehler über `Window.Error`, während die Variante für `NativeWindow` den Fehler zurückgibt. Die Fensteroption `Mac.Toolbar` hängt eine Symbolleiste bei der Erstellung an. Sie wird nach `Mac.SplitView` angewendet, damit ein mitlaufender Trenner die Seitenleiste findet, an der er ausgerichtet wird.

## Split-View

`MacSplitView` ordnet die Teilbereiche an. Fügen Sie sie in der Reihenfolge vom vorderen zum hinteren Rand hinzu. `AddSidebar`, `AddContentList` und `AddInspector` übernehmen jeweils das native Modell, das sie enthalten, und geben einen `MacSplitPane` zur Steuerung von Größe und Einklappen zurück. `AddPrimaryContent` platziert die vorhandene WebView des Fensters und gibt einen `MacSplitWebviewPane` zurück, der zusätzlich `SetContentLayout` bietet.

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

`SetSplitView` funktioniert vor und nach der Erstellung des nativen Fensters. Bei einem Aufruf davor wird das Layout vorgemerkt und bei der Fenstererstellung installiert. Bei einem späteren Aufruf, etwa aus einem Menü- oder Tray-Callback in einer laufenden Anwendung, wird es sofort installiert: Die vorhandene WebView des Fensters wird zum Hauptteilbereich, die aktuelle Symbolleiste wird erneut angehängt, damit ihre mitlaufenden Trenner ausgerichtet sind, und alle vorgemerkten Zusatzleisten werden angehängt.

Ein nach `app.Run` erstelltes Fenster kann mit den Optionen `Mac.SplitView` und `Mac.Toolbar` in einem Aufruf konfiguriert werden:

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

Für ein Layout gelten folgende Regeln:

- Ein Layout benötigt mindestens zwei Teilbereiche und genau einen Hauptteilbereich (`AddPrimaryContent` für ein `WebviewWindow`, `AddTextEditor` für ein `NativeWindow`).
- Es darf höchstens eine Inhaltsliste geben, die nach der Seitenleiste und vor dem Hauptteilbereich liegt.
- Die Struktur der Teilbereiche ist festgelegt, sobald die Split-View angehängt wurde. Einstellungen und Einklappzustand der Teilbereiche sowie die Inhalte von Seitenleiste, Liste und Inspektor können weiterhin jederzeit geändert werden.
- Ein installiertes Layout kann nicht ersetzt werden. Ein zweiter Aufruf von `SetSplitView` für dasselbe Fenster meldet `ErrMacSplitViewAlreadyInstalled` (bei einem `WebviewWindow` über `Window.Error`, bei einem `NativeWindow` als Rückgabewert) und lässt das Fenster unverändert. Wird vor der Installation `nil` übergeben, wird ein vorgemerktes Layout gelöscht.
- Eine Seitenleiste, Liste, ein Inspektor oder eine Split-View gehört jeweils nur zu einem Fenster.

Setter für Teilbereiche sind `SetMinimumThickness`, `SetMaximumThickness`, `SetPreferredThicknessFraction`, `SetHoldingPriority`, `SetCollapsible`, `SetCanCollapseFromWindowResize`, `SetCollapsed`, `Toggle`, `IsCollapsed` und `OnCollapsedChange`. `SetAutosaveName` speichert die Positionen der Trennlinien über Anwendungsstarts hinweg.

`SetContentLayout` im Hauptteilbereich wählt zwischen `MacContentLayoutBelowToolbar` und `MacContentLayoutEdgeToEdge`. `MacContentLayoutAutomatic` übernimmt `MacWindow.ContentLayout`, das wiederum `TitleBar.FullSizeContent` folgt. Die randfüllende Anordnung ermöglicht AppKit ab macOS 26, den Effekt am Scrollrand unter der Symbolleiste anzuwenden.

## Seitenleiste

`MacSidebar` ist eine native Quellliste. Sie enthält Zeilen auf der obersten Ebene, Abschnitte und beliebig tief verschachtelte Zeilen. Bewahren Sie die zurückgegebenen Handles auf, um Zeilen später zu aktualisieren.

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

Setter für Zeilen sind `SetLabel`, `SetSymbol`, `SetTooltip`, `SetEnabled`, `SetHidden`, `SetBadge`, `SetAccessorySymbol`, `SetTintColor`, `SetEditable` und `SetExpanded`. `OnClick` wird ausgelöst, wenn AppKit die Zeile auswählt, `OnExpandedChange`, wenn der Benutzer ihre untergeordneten Zeilen öffnet oder schließt, und `OnRename`, nachdem eine direkte Umbenennung übernommen wurde.

### Auswahl

Standardmäßig kann nur ein Element ausgewählt werden. `SetSelectedItem` wählt eine Zeile aus, ohne deren `OnClick` auszulösen. Ist Mehrfachauswahl aktiviert, wird `OnClick` weiterhin für die angeklickte Zeile ausgelöst, und `OnSelectionChange` meldet die gesamte Auswahl.

```go
sidebar.SetAllowsMultipleSelection(true)
sidebar.OnSelectionChange(func(_ *application.Context, items []*application.MacSidebarItem) {
    for _, item := range items {
        _ = item.Section()
    }
})
```

### Kontextmenüs

Bei einem Rechtsklick wird in dieser Reihenfolge nach einem Menü gesucht: zuerst das mit `SetContextMenu` gesetzte Menü der Zeile, dann der `OnContextMenu`-Callback der Seitenleiste und schließlich das mit `SetContextMenu` gesetzte Ersatzmenü der Seitenleiste. Der Callback läuft im Anwendungsthread, während AppKit wartet; halten Sie ihn daher kurz.

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

### Umordnen durch Ziehen

`SetReorderable` ermöglicht es Benutzern, Zeilen innerhalb eines Abschnitts, zwischen Abschnitten sowie zur oder von der obersten Ebene zu ziehen. Das Go-Modell wird aktualisiert, bevor `OnMove` ausgelöst wird.

```go
sidebar.SetReorderable(true)
sidebar.OnMove(func(_ *application.Context, item *application.MacSidebarItem, section *application.MacSidebarSection, index int) {
    // section is nil when the row was dropped at the sidebar root
})
```

### Entfernen

Beim Entfernen einer Zeile werden auch ihre untergeordneten Zeilen entfernt. Die Handles sind danach wirkungslos.

```go
notes.Remove(draft)
sidebar.RemoveSection(tags)
```

## Inhaltsliste

`MacContentList` bildet die mittlere Spalte von Finder, Mail und Dokumentbrowsern. Ohne Spalten zeigt sie umfangreiche Zeilen mit Titel, Untertitel, vorangestelltem Symbol, nachgestelltem Detail und Zähler-Badge.

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

Die Stile sind `MacContentListStyleAutomatic`, `MacContentListStyleInset`, `MacContentListStyleSourceList`, `MacContentListStylePlain` und `MacContentListStyleFullWidth`. `SetRowHeight`, `SetAlternatingRowBackgrounds`, `SetHeaderVisible` und `SetAllowsMultipleSelection` ergänzen die Darstellungsoptionen. Zeilen unterstützen `SetHidden`, `SetEnabled`, `SetTooltip`, `InsertRow` an einer bestimmten Position, `Remove` und `RemoveAll`.

### Spalten

`SetColumns` schaltet in den Tabellenmodus. Zeilen zeigen dann ihre `SetCells`-Werte, einen pro Spalte. Mit `Sortable` markierte Spalten zeigen einen Sortierindikator; `SetSortable` aktiviert das Sortieren durch Klick auf die Spaltenüberschrift. Ohne `OnSort`-Callback sortiert sich die Liste mit `SortBy` selbst nach dem Spaltentext.

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

### Kontextmenüs

Kontextmenüs werden in derselben Reihenfolge wie in der Seitenleiste bestimmt: zuerst `SetContextMenu` der Zeile, dann `OnContextMenu` und schließlich das Ersatzmenü der Liste.

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

## Inspektor

`MacInspector` ist ein am hinteren Rand angeordnetes Eigenschaftenfenster aus nativen Steuerelementen, die in Abschnitte gruppiert sind. Jede `Add`-Methode gibt ein `MacInspectorControl`-Handle mit typspezifischen Settern und Callbacks zurück.

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

Steuerelementtypen und ihre Setter:

| Steuerelement | Setter | Callback |
|---------|---------|----------|
| `AddLabel` | `SetValue` | keiner |
| `AddTextField` | `SetValue` | `OnTextChange` |
| `AddCheckbox` | `SetChecked` | `OnToggle` |
| `AddPopup`, `AddSegmented` | `SetOptions`, `SetSelectedIndex` | `OnSelectionChange` |
| `AddSlider`, `AddStepper` | `SetFloatValue`, `SetRange`, `SetStep` | `OnValueChange` |
| `AddColorWell` | `SetColor` | `OnColorChange` |
| `AddDatePicker` | `SetDate` | `OnDateChange` |
| `AddButton` | `SetLabel` | `OnClick` |

`SetLabel`, `SetTooltip`, `SetEnabled` und `SetHidden` gelten für jeden Typ. Abschnitte können einklappbar sein; sowohl Abschnitte als auch Steuerelemente lassen sich jederzeit verschieben oder entfernen.

```go
appearance.SetCollapsed(true)
inspector.MoveSection(statistics, 0)
appearance.Remove(priority)
inspector.RemoveSection(statistics)
```

## Zusatzleisten

`MacAccessory` ist eine Leiste mit nativen Steuerelementen. Wählen Sie beim Erstellen ein Layout: `MacAccessoryLayoutLeading` liegt neben den Fensterschaltflächen, `MacAccessoryLayoutTrailing` am hinteren Rand der Titelleiste, `MacAccessoryLayoutBottom` erstreckt sich über die gesamte Breite unter Titelleiste und Symbolleiste, und `MacAccessoryLayoutTop` liegt oben in einem Split-View-Teilbereich. Fügen Sie alle Steuerelemente hinzu, bevor Sie die Zusatzleiste anhängen.

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

Verfügbare Steuerelemente sind `AddSearch`, `AddSegmented`, `AddButton`, `AddSymbolButton`, `AddMenuButton`, `AddLabel`, `AddFlexibleSpace` und für native Integrationen `AddNativeView`. `AddTitlebarAccessory` gibt es sowohl für `WebviewWindow` als auch für `NativeWindow`. Aufrufe vor der Erstellung des Fensters werden vorgemerkt und bei der Erstellung angewendet. `Remove` löst eine Zusatzleiste, damit sie anderswo erneut angehängt werden kann; `SetHidden` klappt sie an ihrem Platz ein.

### Zusatzleisten für Teilbereiche

Ab macOS 26 kann eine Zusatzleiste oben oder unten in einem Split-Teilbereich sitzen, etwa dort, wo Finder sein Filterfeld für die Seitenleiste platziert. Erstellen Sie die Zusatzleiste mit dem Layout `Top` oder `Bottom` und hängen Sie sie mit `AddTopAccessory` oder `AddBottomAccessory` an den Teilbereich.

```go
filter := application.NewMacAccessory(application.MacAccessoryLayoutTop)
filter.AddSearch("Filter").
    SetIncremental(true).
    OnSearch(func(_ *application.Context, query string) { filterNotes(query) })

if err := sidebarPane.AddTopAccessory(filter); err != nil {
    window.Error("sidebar filter: %s", err)
}
```

`SetPreferredScrollEdgeEffectStyle` wählt die Darstellung `Automatic`, `Soft` oder `Hard` von AppKit für Inhalte, die hinter der Zusatzleiste scrollen. Explizite Stile erfordern macOS 26.1. Unter älteren Versionen wird die Anforderung über den Fehlerhandler des Fensters gemeldet, und der automatische Stil bleibt aktiv.

## Alles zusammensetzen

Dieses Programm erstellt ein Fenster mit drei Teilbereichen: einer Seitenleiste, der WebView und einem Inspektor. Seine Symbolleiste folgt beiden Trennlinien.

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

Die nativen Fensterelemente und das Frontend kommunizieren über die üblichen Wails-Ereignisse und -Dienste. Die Seitenleiste sendet `note:selected`, der Inspektor sendet `note:title`, und die Seite lauscht mit `Events.On` der Laufzeitumgebung darauf.

## Native Fenster

@note{type="caution" title="Experimentell"}
`NativeWindow`, `NativeWindowManager` und `MacTextEditor` sind in v3 experimentell. Die API ist bewusst klein gehalten und kann sich ändern, wenn die gemeinsame Fenster-API für v4 neu gestaltet wird.
@end

Ein `NativeWindow` hat keine WebView. Sein Hauptinhalt ist ein `MacTextEditor`, also ein `NSTextView` innerhalb eines `NSScrollView`. Es akzeptiert dieselben Typen für Symbolleiste, Split-View und Zusatzleisten wie ein `WebviewWindow`. Erstellen Sie es mit `app.NativeWindow.New` oder `app.NativeWindow.NewWithOptions` und rufen Sie es später mit `Get` oder `GetByID` ab.

Ein natives Fenster wird erst erstellt, wenn es Inhalt hat. Stellen Sie eine Split-View bereit, deren Hauptteilbereich mit `AddTextEditor` hinzugefügt wurde, entweder über `NativeWindowOptions.SplitView` oder durch Aufruf von `SetSplitView`. Vor `app.Run` wird das Layout vorgemerkt. In einer laufenden Anwendung erstellt und zeigt `SetSplitView` das Fenster sofort an (sofern `Hidden` nicht gesetzt ist) und gibt etwaige Erstellungsfehler zurück. Ein Fenster ohne Layout bleibt zurückgestellt, und `Run` gibt `ErrNativeWindowContentRequired` zurück. Ein Layout ohne Texteditor wird mit `ErrNativeWindowEditorRequired` abgelehnt. Verwenden Sie `NativeWindowOptions.Toolbar` und `NativeWindowOptions.SplitView` statt der Felder `Mac.Toolbar` und `Mac.SplitView`, die ein natives Fenster ignoriert.

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

Dasselbe Fenster kann in einer laufenden Anwendung mit einem einzigen Aufruf erstellt werden, indem die nativen Fensterelemente als Optionen übergeben werden:

```go
window := app.NativeWindow.NewWithOptions(application.NativeWindowOptions{
    Title:     "Native Notes",
    SplitView: split,
    Toolbar:   toolbar,
})
```

`MacTextEditor` bietet `SetText`, `Text`, `SetEditable`, `OnChange` und `Focus`. Programmatische Aufrufe von `SetText` lösen `OnChange` nicht aus. Das Laden einer Datei markiert sie daher nicht als geändert. `Text` liest das gesamte Dokument aus AppKit. Rufen Sie es nur auf, wenn Sie den Inhalt benötigen, und nicht bei jeder Änderung.

Zwei Optionen halten eine rein native Anwendung schlank:

- `NativeOnly: true` in `application.Options` überspringt zur Laufzeit den Frontend-Transport und den Asset-Server. Erstellen Sie kein `WebviewWindow`, wenn diese Option gesetzt ist.
- Das Build-Tag `wails_native` entfernt beim Kompilieren WebView-, Frontend- und Updater-Code vollständig aus der Binärdatei und setzt `NativeOnly` automatisch:

```sh
go build -tags wails_native .
```

Die Unterstützung für eine einzelne Anwendungsinstanz ist ebenfalls aus einem `wails_native`-Build ausgeschlossen. Fügen Sie das Tag `wails_single_instance` hinzu, wenn Sie sie benötigen.

## Fenstertabs

macOS kann Fenster in Tabs gruppieren. Der Tab-Modus wird bei der Erstellung eines Fensters festgelegt. Setzen Sie daher `Mac.TabbingMode` für jedes teilnehmende Fenster auf `MacWindowTabbingModePreferred` oder `MacWindowTabbingModeAutomatic`. `MacWindowTabbingModeDisallowed` und der nicht gesetzte Standardwert schließen ein Fenster von Tab-Gruppen aus.

```go
first := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Document 1",
    URL:   "/",
    Mac: application.MacWindow{
        TabbingMode: application.MacWindowTabbingModePreferred,
    },
})
```

Tab-Operationen benötigen aktive native Fenster. Rufen Sie sie daher aus einem Menühandler, einer Dienstmethode oder anderem Code auf, der nach dem Start von `app.Run` ausgeführt wird. `TabGroup` gibt ein `MacWindowTabGroup`-Handle zurück, das die native Gruppe bei jedem Aufruf erneut ermittelt. Ein nil-Handle kann gefahrlos verwendet werden; jede seiner Methoden gibt den jeweiligen Nullwert zurück.

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

`MacWindowTabGroup` bietet `Identifier`, `Count`, `Windows`, `NativeWindows`, `SelectedWindow`, `Select`, `SelectNative`, `SelectNext`, `SelectPrevious`, `IsTabBarVisible`, `ToggleTabBar`, `IsOverviewVisible` und `ToggleTabOverview`. `app.Window.TabGroups` listet alle Gruppen auf, die ein `WebviewWindow` enthalten. `AddNativeTab` fügt ein `NativeWindow` zur Gruppe eines WebView-Fensters hinzu, und `SetTabTooltip` legt den Text fest, der beim Zeigen auf einen Tab erscheint. Außerhalb von macOS geben die Tab-Methoden `ErrMacWindowTabsUnsupported` zurück.

## Versionsanforderungen

Alle Funktionen auf dieser Seite sind nur für macOS verfügbar. Auf älteren Versionen werden Funktionen wie unten beschrieben eingeschränkt; die Go-API ist überall identisch.

| Funktion | Mindestversion von macOS | Verhalten unter älteren Versionen |
|---------|---------------|-------------------------------|
| Fenstertabs | 10.12 | nicht verfügbar |
| Symbolleistengruppen, umrandete Elemente, `AddMenu` | 10.15 | Menüelemente werden weggelassen; Gruppen verwenden die ältere Darstellung |
| SF Symbols (`SetSymbol` für Elemente in Symbolleiste, Seitenleiste, Liste und Zusatzleiste) | 11 | kein Bild wird angezeigt |
| `AddSearch` als `NSSearchToolbarItem`, `SetNavigational`, mitlaufender Seitenleisten-Trenner | 11 | die Suche verwendet ersatzweise ein einfaches Suchfeld; der Trenner wird weggelassen |
| Split-Rollen für Inspektor und Inhaltsliste | 11 | dieselben Teilbereiche werden in gewöhnlichen Split-Elementen angezeigt |
| Stile für Inhaltslisten | 11 | ignoriert |
| `SetCenteredItems` | 13 | ignoriert |
| Inspektor-Umschalter und mitlaufender Trenner | 14 | Wails stellt eine native Umschaltfläche bereit; der Trenner wird weggelassen |
| `SetBadgeCount`, `SetProminent`, `SetTintColor` für Symbolleistenelemente | 26 | gespeichert und angewendet, sobald verfügbar |
| Zusatzleisten für Teilbereiche (`AddTopAccessory`, `AddBottomAccessory`) | 26 | `ErrMacSplitItemAccessoryUnavailable` |
| Explizite Stile für den Scrollrand-Effekt | 26.1 | über den Fenster-Fehlerhandler gemeldet; der automatische Stil bleibt aktiv |

Für Titelleisten-Zusatzleisten, Split-Views, Seitenleisten, Inspektoren und den Texteditor gelten über die Mindestanforderung von Wails hinaus keine Versionsanforderungen.

## Beispiele

Jedes Beispiel ist eine vollständige, ausführbare Anwendung:

- [`v3/examples/mac-toolbar`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-toolbar): ein Notizeditor mit Symbolleiste, Seitenleiste, Inhaltsliste, Inspektor, Share-Provider, anpassbarer Symbolleiste und einer Zusatzleiste für einen Teilbereich.
- [`v3/examples/mac-titlebar-accessory`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-titlebar-accessory): Zusatzleisten am vorderen Rand, hinteren Rand und unterhalb der Titelleiste zur Steuerung einer Postfachansicht.
- [`v3/examples/mac-window-tabs`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-window-tabs): Tab-Modi, `AddTab`, `AddNativeTab` und die Tab-Gruppen-API über ein Menü.
- [`v3/examples/mac-native-editor`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-native-editor): ein `NativeWindow`-Texteditor ohne WebView, erstellt mit dem Tag `wails_native`.

## Weitere Möglichkeiten

Die hier beschriebenen Fensterelemente sind ein Teil einer nativen macOS-Anwendung. Der andere Teil ist das Verhalten der Anwendung: Dokumentfenster mit Proxy-Symbolen und gestaffelter Anordnung, Menüs mit Symbolen und Badges, ein Dock-Menü mit Fortschrittsanzeige, native Warnmeldungen und Dialoge, entfernbare Statusobjekte, haptisches Feedback und Sprachausgabe, die Zwischenablage für vielfältige Inhalte und das Herausziehen von Elementen sowie Systeminformationen zu Berechtigungen, Stromversorgung und Gebietsschema. Diese APIs behandelt der Leitfaden [macOS-Plattformintegration](/guides/macos-platform-integration).
