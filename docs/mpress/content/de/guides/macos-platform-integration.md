---
title: "macOS-Plattformintegration"
description: "Dokumentfenster, Dock- und Menüerweiterungen, native Dialoge, Statusobjekte, Feedback, die Zwischenablage für vielfältige Inhalte, Herausziehen von Elementen, Berechtigungen, Stromversorgung und Lebenszyklus unter macOS"
slug: "guides/macos-platform-integration"
sourcePath: "guides/macos-platform-integration.md"
---

Relevante Plattformen: macOS

Wails v3 gibt Ihrer Anwendung das Verhalten, das macOS-Benutzer von einer nativen Anwendung erwarten: Dokumentfenster mit Proxy-Symbolen und gestaffelter Anordnung, Menüs mit Symbolen und Badges, ein Dock-Menü mit Fortschrittsanzeige, native Warnmeldungen und Dialoge, entfernbare Statusobjekte, haptisches Feedback und Sprachausgabe, eine Zwischenablage für vielfältige Inhalte samt Herausziehen von Elementen, Systeminformationen zu Berechtigungen, Stromversorgung und Gebietsschema sowie Integration in das Dienste-Menü, Handoff, AppleScript und Quick Look. Alles wird aus Go über das Paket `application` gesteuert.

Derselbe Code lässt sich unter Windows und Linux kompilieren. Setter speichern ihre Werte, Abfragen geben Nullwerte zurück, und Operationen, die macOS benötigen, geben einen dokumentierten Fehler wie `ErrMacOnly`, `ErrDialogNotSupported` oder `ErrClipboardNotSupported` zurück. Die [Plattformhinweise](#platform-notes) unten beschreiben das Verhalten der einzelnen Bereiche außerhalb von macOS.

Informationen zu nativen Fensterelementen wie Symbolleisten, Seitenleisten, Inspektoren, Zusatzleisten und Fenstertabs finden Sie im Leitfaden [Native macOS-Fensterelemente](/guides/macos-native-chrome).

## Dokumentfenster

Ein Dokumentfenster zeigt die zugehörige Datei in der Titelleiste an, kennzeichnet ungespeicherte Änderungen mit einem Punkt in der Schließen-Schaltfläche und öffnet neue Fenster gestaffelt. Dies wird über Methoden von `WebviewWindow` und Optionen von `MacWindow` gesteuert.

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

- `SetRepresentedFile` zeigt das Proxy-Symbol der Datei in der Titelleiste. Benutzer können das Symbol in eine andere Anwendung ziehen oder bei gedrückter Befehlstaste darauf klicken, um den Pfad anzuzeigen. Übergeben Sie `""`, um es zu entfernen. `RepresentedFile` liest den Wert aus.
- `SetDocumentEdited` zeigt den Punkt für ungespeicherte Änderungen in der Schließen-Schaltfläche und blendet das Proxy-Symbol ab. `IsDocumentEdited` liest den Wert aus.
- `SetSubtitle` zeigt ab macOS 11 eine zweite Zeile unter dem Titel.
- `InitialPosition: application.WindowCascade` platziert das Fenster unterhalb und rechts vom zuletzt gestaffelten Fenster, wie beim Öffnen neuer Dokumente. `CascadeFrom(other)` tut dies für ein bestehendes Fenster und aktualisiert den Staffelungspunkt für spätere Fenster.
- `Mac.FrameAutosaveName` stellt die gespeicherte Position und Größe wieder her, bevor das Fenster zum ersten Mal angezeigt wird, und speichert sie bei jeder Bewegung des Fensters erneut. Ein wiederhergestellter Fensterrahmen hat Vorrang vor `X`, `Y`, `Width`, `Height` und `InitialPosition`. `SetFrameAutosaveName` ändert den Namen für ein aktives Fenster.
- `MacTitleBar.WindowButtonsOffset` verschiebt die Schaltflächen zum Schließen, Minimieren und Vergrößern um eine Anzahl von Punkten. `SetWindowButtonsOffset` und `ResetWindowButtonsOffset` ändern den Versatz zur Laufzeit.

Alle drei Setter können aufgerufen werden, bevor das native Fenster existiert. Die Werte werden bei seiner Erstellung angewendet.

### Aufmerksamkeit anfordern

`RequestAttention` lässt das Dock-Symbol springen, während die Anwendung im Hintergrund läuft. Eine informative Anforderung lässt es einmal springen. Bei einer dringenden Anforderung springt es weiter, bis der Benutzer die Anwendung aktiviert oder Sie die Anforderung abbrechen.

```go
request := window.RequestAttention(true)

// Once the work that needed attention is done:
request.Cancel()
```

`Flash` bleibt die plattformübergreifende Möglichkeit, einmalig Aufmerksamkeit anzufordern.

### Drucken und Exportieren

`PrintWithOptions` druckt die WebView mit ausdrücklich festgelegten Seiteneinstellungen. Der Nullwert zeigt das Druckfenster mit den gemeinsamen Druckeinstellungen. `Print` behält sein bisheriges Verhalten bei (Querformat, Ränder von 30 Punkten).

`ExportPDF` rendert die Seite als PDF-Dokument, und `Snapshot` erfasst sie als PNG. Beide warten auf WebKit. Rufen Sie sie daher aus einer Goroutine und niemals aus dem Anwendungsthread auf. Ein Aufruf dort gibt `ErrMacExportOnMainThread` zurück.

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

`PrintOptions` akzeptiert außerdem `PrinterName`, `PaperName` (einen PostScript-Namen wie `"iso-a4"`) und `Scale`. `PDFExportOptions` und `SnapshotOptions` akzeptieren ein optionales `Rect`, um die Erfassung zu begrenzen, und ein `Timeout`, dessen Standardwert `DefaultMacExportTimeout` (30 Sekunden) ist.

## Dialogblätter

Ein Dialogblatt ist ein zweites Fenster, das am oberen Rand seines übergeordneten Fensters hängt, wie ein Sichern-Dialog. Jedes `WebviewWindow` kann mit `PresentSheet` als Dialogblatt eines anderen Fensters angezeigt und mit `EndSheet` und einem Antwortcode beendet werden. Der Code erreicht die `OnSheetEnd`-Callbacks des Dialogblatts. Erstellen Sie das Dialogblattfenster mit gesetztem `Hidden`, damit es vor dem Anhängen nicht kurz auf dem Bildschirm erscheint. AppKit blendet es beim Beenden wieder aus, sodass dasselbe Fenster wiederholt angezeigt werden kann.

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

`PresentCriticalSheet` zeigt das Dialogblatt vor einem bereits angehängten gewöhnlichen Dialogblatt, statt es dahinter einzureihen. `PresentNativeSheet` tut dasselbe für ein `NativeWindow`. `IsSheet`, `SheetParent`, `AttachedSheet`, `AttachedNativeSheet` und `HasAttachedSheet` beschreiben den aktuellen Zustand. Wird ein Dialogblattfenster geschlossen statt beendet, wird `MacSheetResponseStop` übermittelt.

## Popovers

Ein `MacPopover` ist ein `NSPopover`: ein vorübergehendes Fenster, das an einem Rechteck in einem Fenster, an einem Symbolleistenelement oder am Statusobjekt in der Menüleiste verankert ist. Sein Inhalt ist eine native `MacAccessory`-Leiste mit Steuerelementen, derselbe Typ, den der Leitfaden [Native macOS-Fensterelemente](/guides/macos-native-chrome) für Titelleisten-Zusatzleisten verwendet. Fügen Sie alle Steuerelemente vor der ersten Anzeige hinzu.

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

`MacToolbarItem.ShowPopover` und `SystemTray.ShowPopover` verankern dasselbe Popover an einem Symbolleistenelement oder Statusobjekt. `MacPopoverBehaviorTransient` schließt es bei jedem Klick außerhalb, `Semitransient` nur bei Klicks im Fenster, von dem aus es angezeigt wird. Beim Standardverhalten bleibt es bis zu `Close` geöffnet. `MacRectEdge` bestimmt, auf welcher Seite das Popover erscheint. `SetContentSize` und `SetBehavior` passen ein aktives Popover an, und `Destroy` gibt das native Popover sowie die Inhaltsleiste zur Verwendung an anderer Stelle frei.

## Zustandswiederherstellung

macOS öffnet die Fenster einer Anwendung nach einem Absturz, einem erzwungenen Beenden oder einem Neustart erneut. Nach einem normalen Beenden geschieht dies, wenn „Fenster beim Beenden einer App schließen“ in den Systemeinstellungen deaktiviert ist. Geben Sie einem Fenster eine `Mac.RestorationID`, speichern Sie die für die Neuerstellung nötigen Daten mit `SetRestorationData` und registrieren Sie `app.Window.OnRestore`, um das Fenster beim nächsten Start neu zu erstellen.

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

Nur Fenster, die beim Beenden der Anwendung sichtbar sind, werden gespeichert. `RestorationState.Data` ist eine String-Map; beschränken Sie ihren Inhalt auf Kennungen, Pfade und Positionen. `InteractionState` gibt die Vorwärts- und Zurück-Liste der WebView sowie Scrollpositionen als undurchsichtigen Datenblock zurück (ab macOS 12), den `RestoreInteractionState` auf ein neu erstelltes Fenster anwendet. Üblicherweise wird er Base64-kodiert in den Wiederherstellungsdaten gespeichert. `SetRestorationID` und `RestorationID` ändern beziehungsweise lesen die Kennung eines aktiven Fensters.

## Darstellungsoptionen

`MacPresentationOptions` entspricht `NSApplication.presentationOptions`: einer Bitmaske, die das Dock oder die Menüleiste ausblendet und während der aktiven Anwendung den Prozesswechsel, das erzwungene Beenden, das Abmelden oder den Befehl zum Ausblenden deaktiviert. Setzen Sie `Mac.PresentationOptions` in den Anwendungsoptionen, um dies beim Start anzuwenden, oder ändern Sie es zur Laufzeit mit `SetPresentationOptions`. Ungültige Kombinationen werden vor der Übergabe an AppKit mit einem Fehler abgelehnt, der `ErrMacPresentationOptionsInvalid` umschließt.

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

Das Ausblenden der Menüleiste (`HideMenuBar` oder `AutoHideMenuBar`) erfordert eine der Dock-Optionen. `AutoHideToolbar` erfordert sowohl `FullScreen` als auch `AutoHideMenuBar`. `Validate` meldet die erste Regel, gegen die ein Wert verstößt, und `Has` prüft einzelne Flags.

## Menüs und Dock

Menüelemente erhalten SF Symbols, Badges, Abschnittsüberschriften, Farbpaletten, den gemischten Auswahlzustand, alternative Elemente und Einrückungen. Dies sind alles Methoden von `MenuItem` und `Menu`. Sie funktionieren daher im Anwendungsmenü, in Kontextmenüs, Tray-Menüs und im Dock-Menü.

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

- `SetSymbol` zeigt ab macOS 11 ein SF Symbol neben dem Titel. Es ersetzt ein mit `SetBitmap` gesetztes Bild.
- `SetBadge` zeigt eine Zahl und `SetBadgeText` eine kurze Zeichenfolge nach dem Titel (ab macOS 14). `ClearBadge` entfernt das Badge; `BadgeCount` und `BadgeText` lesen es aus.
- `AddSectionHeader` fügt eine nicht interaktive Überschrift hinzu (ab macOS 14). Unter älteren Versionen ist sie ein deaktiviertes Element mit demselben Titel.
- `AddPalette` fügt eine Reihe von Farbfeldern hinzu, die vom Palettenmenü von `NSMenu` bereitgestellt wird (ab macOS 14). Übergeben Sie ein Symbol für jedes Farbfeld, eines pro Farbe, oder einen leeren Slice für ausgefüllte Kreise. Ohne Beschriftung erscheint die Palette direkt im übergeordneten Menü; `SetLabel` zeigt sie als Untermenü mit Titel. `PaletteSelected` gibt den ausgewählten Index zurück.
- `SetMixed` versetzt ein Kontrollkästchen in den gemischten Zustand, der als Strich dargestellt wird. Ein Klick aktiviert es vollständig, wie bei AppKit.
- `SetAlternate(true)` zeigt das Element anstelle des darüberliegenden an, solange die abweichende Sondertaste gedrückt wird. Beide Elemente müssen dieselbe Taste verwenden und sich in den Sondertasten unterscheiden.
- `SetIndentationLevel` rückt den Titel um bis zu 15 Ebenen ein.

### Zuletzt geöffnet

`fileMenu.AddRole(application.OpenRecent)` fügt das standardmäßige Untermenü „Zuletzt geöffnet“ hinzu. Unter macOS wird es bei jedem Öffnen von `NSDocumentController` gefüllt und enthält einen Eintrag zum Leeren des Menüs. Fügen Sie Dateien mit `app.Menu.AddRecentDocument` hinzu, listen Sie sie mit `RecentDocuments` auf und leeren Sie die Liste mit `ClearRecentDocuments`. Die Liste bleibt über Neustarts hinweg erhalten.

Die Auswahl einer zuletzt verwendeten Datei löst dasselbe Ereignis aus wie das Öffnen einer Datei aus Finder:

```go
app.Event.OnApplicationEvent(events.Common.ApplicationOpenedWithFile, func(event *application.ApplicationEvent) {
    path := event.Context().Filename()
    log.Println("open", path)
})
```

### Dock-Menü

`app.Menu.SetDockMenu` installiert ein statisches Menü für den Rechtsklick auf das Dock-Symbol. `OnDockMenu` erstellt bei Bedarf jedes Mal vor der Anzeige ein Menü. Das ist sinnvoll, wenn die Elemente einen veränderlichen Zustand widerspiegeln. Ein Builder hat Vorrang vor dem statischen Menü.

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

### Dock-Fortschritt

Der Dock-Dienst zeichnet zusätzlich zur bestehenden Badge-Unterstützung einen Fortschrittsbalken über das Dock-Symbol. Registrieren Sie `dock.New()` als Dienst und rufen Sie `SetProgress` mit einem Anteil zwischen 0 und 1 auf.

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

`GetProgress` gibt den aktuellen Anteil zurück oder `nil`, wenn kein Balken angezeigt wird.

## Dialoge

Meldungs-, Öffnen- und Sichern-Dialoge akzeptieren macOS-Optionen. Der Dialog-Manager erhält außerdem eine Texteingabe sowie die systemeigenen Farb- und Schriftfenster.

### Warnmeldungen

`SetSuppression` fügt ein Kontrollkästchen „Diese Meldung nicht mehr anzeigen“ hinzu, und `SetHelp` zeigt die Hilfe-Schaltfläche. Lesen Sie das Kontrollkästchen mit `Suppressed` in einem Schaltflächen-Callback aus oder registrieren Sie `OnSuppression`, um den Wert zuerst zu erhalten.

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

### Texteingabe

`Prompt` zeigt eine Warnmeldung mit einem Textfeld und blockiert, bis sie geschlossen wird. Rufen Sie es daher aus einer Goroutine oder einer gebundenen Methode auf. `Secure` macht das Feld zu einem Passwortfeld, und `Window` zeigt die Warnmeldung als Dialogblatt.

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

### Dateidialoge

`AddContentType` filtert sowohl Öffnen- als auch Sichern-Dialoge nach Uniform Type Identifiern. Es ergänzt `AddFilter`: `"public.image"` erfasst damit alle dem System bekannten Bildtypen, während ein Filter weiterhin PDFs anhand der Dateiendung erfasst.

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

Sichern-Dialoge erhalten ein Format-Auswahlmenü, eine eigene Beschriftung für das Namensfeld und Finder-Tags. `SetFormats` ändert den zulässigen Typ und die Dateiendung im Namensfeld, wenn der Benutzer die Auswahl im Menü ändert. `SelectedFormat` meldet die endgültige Auswahl.

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

### Farb- und Schriftfenster

`PickColor` und `PickFont` öffnen die gemeinsam genutzten Systemfenster und blockieren, bis das jeweilige Fenster geschlossen wird. `OnChange` übermittelt jede Auswahl bei geöffnetem Fenster, sodass die Seite eine Vorschau in Echtzeit anzeigen kann. Von jedem Fenstertyp kann nur eines gleichzeitig geöffnet sein; ein zweiter Aufruf gibt `ErrDialogInProgress` zurück.

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

## Statusobjekte und Feedback

### Statusobjekte

Ein System-Tray-Element ist unter macOS ein `NSStatusItem`. Es kann mit einem SF Symbol dargestellt werden, einen Tooltip haben und vom Benutzer wie die integrierten Elemente entfernt werden.

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

- `SetSymbol` rendert das Symbol als Vorlagenbild, damit es sich dem Erscheinungsbild der Menüleiste anpasst (ab macOS 11). `SetSymbolConfiguration` legt Punktgröße und Strichstärke fest.
- `SetTooltip` legt den Text fest, der beim Zeigen auf das Element erscheint. `Tooltip` liest ihn aus.
- `SetRemovable(true, name)` ermöglicht es Benutzern, das Element bei gedrückter Befehlstaste aus der Menüleiste zu ziehen. Geben Sie ihm einen stabilen Namen für die automatische Speicherung, damit macOS sich das Entfernen über Anwendungsstarts hinweg merkt. `Show` oder `SetVisible(true)` bringt es zurück.
- `IsVisible` liest `NSStatusItem.visible` aus und ist daher false, nachdem der Benutzer das Element entfernt hat. `OnVisibilityChange` meldet jede Änderung.

### Haptisches Feedback

`app.Haptics.Perform` gibt ein Muster auf einem Force Touch-Trackpad oder Magic Trackpad wieder, während die Anwendung aktiv ist.

```go
app.Haptics.Perform(application.HapticAlignment)
```

Die Typen sind `HapticGeneric`, `HapticAlignment` (ein Element rastet an seiner Position ein) und `HapticLevelChange` (eine Raststufe oder Klickstufe). `IsSupported` meldet, ob die Plattform überhaupt haptisches Feedback ausgeben kann.

### Töne

`app.Sound` spielt den Warnton, einen benannten Systemton oder eine Audiodatei ab.

```go
app.Sound.Beep()

if err := app.Sound.Play("Glass"); err != nil {
    log.Println(err)
}

for _, name := range app.Sound.SystemSounds() {
    log.Println(name)
}
```

`Play` akzeptiert einen Namen aus `SystemSounds` oder einen absoluten Pfad zu einer Datei, die Core Audio dekodieren kann. `PlayData` spielt eine vollständige Audiodatei aus dem Speicher ab.

### Sprache

`app.Speech.Speak` reiht Text für die Systemstimme ein und gibt eine `Utterance` zurück. Äußerungen werden nacheinander abgespielt; `Stop` entfernt eine und `StopAll` leert die Warteschlange. `Voices` listet die installierten Stimmen mit ihren Kennungen und Sprachen auf.

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

`Recognize` transkribiert Audio vom Standardmikrofon mit `SFSpeechRecognizer`. Der erste Aufruf fragt nach Berechtigungen für Mikrofon und Spracherkennung und blockiert, bis der Benutzer antwortet. Rufen Sie ihn daher aus einer Goroutine auf. Teiltranskripte kommen über `OnPartial`; `Stop` beendet die Aufnahme und gibt den endgültigen Text zurück.

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

Für die Spracherkennung ist eine als Bundle gepackte Anwendung erforderlich, deren `Info.plist` `NSSpeechRecognitionUsageDescription` und `NSMicrophoneUsageDescription` deklariert. Ohne diese Einträge verweigert macOS den Zugriff, und `Recognize` gibt `ErrSpeechRecognitionUsageDescription` zurück.

## Zwischenablage und Ziehen

### Zwischenablage für vielfältige Inhalte

`app.Clipboard` liest und schreibt neben einfachem Text auch Bilder, Dateiverweise, HTML, RTF und Rohdaten unter beliebigen Uniform Type Identifiern. `Types` listet die Typen auf der Zwischenablage auf, und `OnChange` meldet Änderungen durch jede Anwendung.

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

`SetImage` und `Image` arbeiten mit PNG-Bytes. Bilder, die andere Anwendungen als TIFF kopiert haben, werden automatisch konvertiert. Es gibt keine Systembenachrichtigung über Änderungen der Zwischenablage. Deshalb fragt `OnChange` den Änderungszähler alle 500 ms ab, solange mindestens ein Listener vorhanden ist.

### Elemente herausziehen

`StartDrag` beginnt einen systemweiten Ziehvorgang aus dem Fenster, als hätte der Benutzer die Elemente in Finder aufgenommen. Es bietet vorhandene Dateien, Dateiversprechen, deren Inhalt erst erzeugt wird, wenn ein Ziel das Ablegen annimmt, oder einfachen Text an. Starten Sie den Vorgang während einer Mausgeste: Binden Sie eine Go-Methode und rufen Sie sie im `mousedown`- oder `pointerdown`-Handler der Seite für das ziehbare Element auf. Setzen Sie dabei das HTML-Attribut `draggable` auf `false`, damit WebKit keinen eigenen Ziehvorgang startet.

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

Wird `StartDrag` außerhalb einer Geste aufgerufen, gibt es `ErrDragOutNoGesture` zurück. `DragItems.Image` und `ImageOffset` legen das Bild unter dem Mauszeiger fest.

### Ablegen aus anderen Anwendungen

Für abgelegte Dateien wird weiterhin das Ereignis `WindowFilesDropped` verwendet. Um Text, URLs oder Bilder anzunehmen, die aus anderen Anwendungen gezogen werden, listen Sie die Typen in `DropTypes` auf und registrieren Sie `OnDrop`. Diese Ablagevorgänge werden an Go statt an die eigenen HTML5-Drop-Handler der Seite übermittelt.

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

## System

### Berechtigungen

`app.Permissions` fragt die Datenschutzberechtigungen des Systems ab und fordert sie an: Kamera, Mikrofon, Bildschirmaufnahme, Bedienungshilfen, Standort, Mitteilungen, Eingabeüberwachung und vollständiger Festplattenzugriff. `Status` zeigt niemals eine Berechtigungsabfrage. `Request` fragt bei noch nicht entschiedenen Berechtigungstypen nach und blockiert, bis der Benutzer antwortet. Rufen Sie es daher aus einer Goroutine auf. `OpenSystemSettings` öffnet den passenden Bereich unter „Datenschutz & Sicherheit“.

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

Vollständiger Festplattenzugriff kann nicht angefordert werden; der Versuch gibt `ErrPermissionNotRequestable` zurück. Verweisen Sie den Benutzer auf den entsprechenden Bereich der Systemeinstellungen. Eine unbeantwortete Anforderung gibt `ErrPermissionRequestTimeout` zurück. Unter macOS bedeutet dies meist, dass der Schlüssel mit der Nutzungsbeschreibung für diesen Berechtigungstyp in `Info.plist` fehlt.

Die Fensteroption `Permissions` wird jetzt unter macOS berücksichtigt. Sie bestimmt, wie `getUserMedia`-Anforderungen von der Seite behandelt werden: `PermissionAllow` überspringt die eigene Abfrage der WebView, `PermissionDeny` lehnt ohne Nachfrage ab, und `PermissionDefault` zeigt die Abfrage an. Die systemweite TCC-Abfrage erscheint weiterhin beim ersten Zugriff auf Kamera oder Mikrofon.

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

### Stromversorgung

`app.Power.PreventSleep` hält das System wach und mit `Display` auch den Bildschirm, bis die zurückgegebene Freigabefunktion aufgerufen wird. Aktive Sperren werden gezählt, sodass mehrere Teile der Anwendung gleichzeitig eine halten können. Der Grund wird in der Aktivitätsanzeige angezeigt.

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

Änderungen werden als `events.Mac.ApplicationDidChangePowerState` (Umschalten des Stromsparmodus) und `events.Mac.ApplicationDidChangeThermalState` gemeldet.

### Lebenszyklus

macOS kann eine inaktive Anwendung beim Abmelden oder Herunterfahren sofort beenden, wenn sie dies mit `NSSupportsSuddenTermination` erlaubt. Mit `NSSupportsAutomaticTermination` kann es eine inaktive Anwendung ohne Fenster beenden. `app.Lifecycle.HoldTermination` setzt beides für einen kritischen Abschnitt wie das Speichern einer Datei aus.

```go
release := app.Lifecycle.HoldTermination("Saving document")
defer release()
// write the file
```

`SetSuddenTerminationEnabled` schaltet das sofortige Beenden zur Laufzeit um. `SuddenTerminationEnabled` meldet den aktuellen Zustand, dessen Anfangswert aus dem Schlüssel in `Info.plist` stammt.

### Umgebung

`app.Env` erhält drei Abfragen zu den Einstellungen des Benutzers.

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

- `Accessibility` bildet „Bewegung reduzieren“, „Transparenz reduzieren“, „Kontrast erhöhen“, „Ohne Farbe differenzieren“, „Farben umkehren“, VoiceOver und Schaltersteuerung ab.
- `KeyboardLayout` gibt die aktive Eingabequelle mit Kennung, lokalisiertem Namen und Sprachen zurück.
- `Locale` gibt das von AppKit für die Anwendung gewählte Gebietsschema sowie die vollständige geordnete Liste `Preferred` des Benutzers zurück. `Identifier` berücksichtigt nur eine Sprache, die das Bundle in `CFBundleLocalizations` deklariert. Verwenden Sie `Preferred`, um selbst eine Sprache auszuwählen.

### Ereignisse

Diese Anwendungsereignisse sind neu. Jedes wird über `app.Event.OnApplicationEvent` übermittelt. Fragen Sie den zugehörigen Manager nach dem aktuellen Wert ab.

| Ereignis | Auslöser | Auslesen mit |
|-------|------------|-----------|
| `events.Mac.ApplicationDidChangePowerState` | der Stromsparmodus wird umgeschaltet | `app.Power.State()` |
| `events.Mac.ApplicationDidChangeThermalState` | die thermische Belastung ändert sich | `app.Power.State()` |
| `events.Common.AccessibilitySettingsChanged` | eine Anzeigeeinstellung der Bedienungshilfen ändert sich | `app.Env.Accessibility()` |
| `events.Mac.ApplicationDidChangeAccessibilitySettings` | die macOS-Form derselben Änderung | `app.Env.Accessibility()` |
| `events.Mac.ApplicationDidChangeKeyboardLayout` | die Eingabequelle ändert sich | `app.Env.KeyboardLayout()` |
| `events.Mac.ApplicationDidChangeLocale` | das Gebietsschema ändert sich | `app.Env.Locale()` |

```go
app.Event.OnApplicationEvent(events.Mac.ApplicationDidChangeThermalState, func(*application.ApplicationEvent) {
    app.Event.Emit("system:power", app.Power.State())
})
app.Event.OnApplicationEvent(events.Common.AccessibilitySettingsChanged, func(*application.ApplicationEvent) {
    app.Event.Emit("system:accessibility", app.Env.Accessibility())
})
```

## Integration

### Dienste-Menü

`app.ServicesProvider.Register` fügt dem Untermenü „Dienste“ einen Eintrag hinzu, den jede macOS-Anwendung für ausgewählten Text oder Dateien anzeigt. Der Handler erhält die Zwischenablage als `ServiceRequest` und gibt eine `ServiceResponse` zum Zurückschreiben zurück. Eine leere Antwort lässt die Auswahl unverändert. Halten Sie den Handler kurz, da AppKit im Hauptthread auf ihn wartet.

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

`Name` ist die von AppKit gesendete Nachricht und muss eine einfache Kennung sein. `SendTypes` und `ReturnTypes` sind Zwischenablagetypen; ein Dienst benötigt mindestens einen davon. Die Registrierung allein macht den Dienst nicht sichtbar: Die `Info.plist` des Bundles muss ihn unter `NSServices` deklarieren. `InfoPlistXML` gibt diesen Block zum direkten Einfügen zurück, und `InfoPlistEntries` liefert dieselben Daten als Maps für einen Plist-Serializer. `NSPortName` in diesen Einträgen ist der `Name` der Anwendung und muss mit `CFBundleName` übereinstimmen.

Mit der Wails-CLI erstellte Projekte können dieselben Dienste einmal in `build/config.yml` deklarieren und den `NSServices`-Block beim Packen erzeugen lassen:

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

Jeder Eintrag muss mit einer in Go registrierten `ServiceDefinition` mit demselben `Name` übereinstimmen. Führen Sie nach der Installation eines neuen Builds `pbs -update` aus, damit das Dienste-Menü die Änderung ohne Abmelden übernimmt.

### Handoff und Benutzeraktivitäten

`app.Activity.Publish` macht eine `NSUserActivity` aktuell, damit Benutzer sie auf einem anderen Gerät fortsetzen, in Spotlight finden oder von Siri vorgeschlagen bekommen können. Die zurückgegebene `PublishedActivity` kann bei Zustandsänderungen aktualisiert und beim Schließen des Dokuments ungültig gemacht werden. Das Veröffentlichen einer neuen Aktivität ersetzt die vorherige.

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

Eingehende Aktivitäten werden über `OnContinue` übermittelt. Universal Links haben den Typ `UserActivityTypeBrowsingWeb`, wobei die Seite in `WebpageURL` steht. Sie werden außerdem als `events.Common.ApplicationLaunchedWithUrl` übermittelt, sodass eine Anwendung URLs über denselben Codepfad behandeln kann. `OnWillContinue`, `OnFailed` und `OnUpdated` decken die übrigen Delegate-Ereignisse ab.

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

Jeder Aktivitätstyp muss unter `NSUserActivityTypes` in `Info.plist` aufgeführt sein. Universal Links benötigen außerdem das Entitlement `com.apple.developer.associated-domains` mit einem Eintrag `applinks:example.com` und die passende Datei `apple-app-site-association` auf dieser Domain.

### Apple Events

`app.AppleEvents.Handle` registriert einen Handler für eine Ereignisklasse und -ID, damit AppleScript, Shortcuts und andere Anwendungen die Anwendung steuern können. Die Codes sind Zeichenfolgen mit vier Zeichen. Der direkte Parameter wird in einen Go-Wert dekodiert (`string`, `[]string` mit Dateipfaden, `int64`, `float64`, `bool`, `[]any` oder `AppleEventRawData`); `Result` der Antwort akzeptiert dieselben Typen. Handler laufen in einer eigenen Goroutine, während das Ereignis ausgesetzt ist.

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

Ein Skript kann den Handler sofort mit der Syntax für rohe Ereignisse aufrufen:

```applescript
tell application id "com.example.notes" to «event WAILnote» "hello"
```

`ScriptingDefinition` erzeugt eine minimale `.sdef`, die jedem Handler einen Befehlsnamen gibt. Legen Sie sie in `Contents/Resources` ab und verweisen Sie in `Info.plist` mit `NSAppleScriptEnabled` und `OSAScriptingDefinition` darauf. Der Skripteditor zeigt sie dann unter „Ablage > Wörterbuch öffnen“ an. Wails behandelt das Ereignis „URL öffnen“ für eigene URL-Schemata bereits selbst. Ein Handler für `"GURL"`/`"GURL"` wird daran angehängt, während ein Handler für `"aevt"`/`"odoc"` die integrierte Übermittlung von „Dokumente öffnen“ ersetzt. `Send` adressiert eine laufende Anwendung über ihre Bundle-Kennung, benötigt in einer als Bundle gepackten Anwendung `NSAppleEventsUsageDescription` und blockiert die aufrufende Goroutine, bis die Antwort eintrifft.

### Quick Look

`app.QuickLook.Preview` öffnet das gemeinsam genutzte Quick Look-Fenster für eine oder mehrere Dateien. Bei mehreren Pfaden zeigt das Fenster Pfeile zum Wechseln zwischen ihnen. `Thumbnail` rendert eine Datei über die Vorschaubild-Provider des Systems und gibt ein PNG zurück; dies funktioniert für Dokumente, Bilder, PDFs und Filme.

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

Pfade müssen absolut sein und existieren. `ClosePreview` und `IsPreviewOpen` verwalten das Fenster. `ThumbnailOptions.IconMode` zeichnet den Dokumentrahmen im Finder-Stil, und `Scale: 2` erzeugt ein Retina-Bild. `Thumbnail` blockiert die aufrufende Goroutine. Rufen Sie es daher aus einer Goroutine oder einer gebundenen Methode auf.

### Workspace-Hilfsfunktionen

`app.Browser` erhält drei Hilfsfunktionen für `NSWorkspace`. `OpenWith` öffnet eine Datei mit einer bestimmten Anwendung, die durch ihre Bundle-Kennung oder ihren Bundle-Pfad angegeben wird. `ApplicationsForFile` listet installierte Anwendungen auf, die eine Datei öffnen können, beginnend mit der Standardanwendung. `ActivateApplication` bringt eine laufende Anwendung in den Vordergrund.

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

`app.Spotlight.Index` fügt Anwendungsinhalte über Core Spotlight zum Systemsuchindex hinzu. Jedes `SearchableItem` hat eine `ID`, einen `Title` und optional eine `Domain` zum gebündelten Entfernen, eine `Description`, `Keywords`, einen `ContentType`, ein PNG-Vorschaubild, eine Deep-Link-`URL` und ein Ablaufdatum. `OnOpen` wird aufgerufen, wenn der Benutzer in Spotlight eines der Elemente auswählt oder mit einer Suchanfrage „In App suchen“ wählt.

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

`IsAvailable` meldet, ob der Index Elemente annimmt. Für die Indizierung ist eine als Bundle gepackte Anwendung erforderlich: Elemente, die eine nicht gebündelte `go run`-Binärdatei indiziert, erscheinen nie in Spotlight. `DeleteAll` entfernt alles, was die Anwendung indiziert hat.

## In Vorbereitung

Ein Beispiel `mac-windows-extra` zu Dialogblättern, Popovers, Darstellungsoptionen und Zustandswiederherstellung wird ergänzt und hier verlinkt, sobald es verfügbar ist.

## Versionsanforderungen

Alle Funktionen auf dieser Seite sind nur für macOS verfügbar; die Go-API ist überall identisch. Wails unterstützt macOS 10.13 und neuer. Funktionen, die eine neuere Version benötigen, verhalten sich wie angegeben eingeschränkt.

| Funktion | Mindestversion von macOS | Verhalten unter älteren Versionen |
|---------|---------------|-------------------------------|
| Berechtigungsstatus für Kamera und Mikrofon | 10.14 | als autorisiert gemeldet (ältere Versionen beschränken den Zugriff auf Aufnahmegeräte nicht) |
| `Speech.Speak` und `Voices` | 10.14 | `ErrSpeechNotSupported` |
| Berechtigungen für Bildschirmaufnahme und Eingabeüberwachung | 10.15 | als autorisiert gemeldet |
| `Speech.Recognize` | 10.15 | `ErrSpeechRecognitionNotSupported` |
| `QuickLook.Thumbnail` | 10.15 | `ErrQuickLookNotSupported` |
| `SetSubtitle` | 11 | ignoriert, mit Eintrag im Debug-Log |
| `ExportPDF` | 11 | `ErrMacExportUnsupported` |
| `SetSymbol` für Menü- und Statusobjekte | 11 | kein Bild wird angezeigt |
| `AddContentType`, `SetFormats` nach UTI | 11 | dieselben Kennungen werden über die ältere API für zulässige Dateitypen angewendet |
| Symbole für Dateiversprechen in `StartDrag` | 11 | ein generisches Dokumentsymbol |
| `PowerState.LowPowerMode` und zugehöriges Ereignis | 12 | immer false; das Ereignis wird nie ausgelöst |
| `InteractionState` und `RestoreInteractionState` | 12 | `ErrMacInteractionStateUnsupported` |
| Mitteilungsbereich in `OpenSystemSettings` | 13 | der ältere Einstellungsbereich für Mitteilungen wird geöffnet |
| Menü-Badges, Abschnittsüberschriften, Paletten | 14 | Badges werden nicht angezeigt; Überschriften sind deaktivierte Elemente; Paletten sind ausgeblendet |
| `MacToolbarItem.ShowPopover` für Elemente ohne eigene Ansicht | 14 | `ErrMacPopoverAnchorUnavailable` |

Für alles Übrige gelten keine Anforderungen über die Mindestversion von Wails hinaus.

## Schlüssel in Info.plist

Mehrere Funktionen hängen von Schlüsseln in der `Info.plist` der Anwendung ab. Die Nutzungsbeschreibungen werden Benutzern in der Berechtigungsabfrage angezeigt. Ohne den jeweiligen Schlüssel zeigt macOS die Abfrage nicht an, und die Anforderung läuft nach einer Zeitüberschreitung ab.

| Schlüssel | Benötigt für |
|-----|-----------|
| `NSCameraUsageDescription` | `Permissions.Request(PermissionKindCamera)`, Kamerazugriff von der Seite |
| `NSMicrophoneUsageDescription` | `Permissions.Request(PermissionKindMicrophone)`, Mikrofonzugriff von der Seite, `Speech.Recognize` |
| `NSSpeechRecognitionUsageDescription` | `Speech.Recognize` |
| `NSLocationUsageDescription` | `Permissions.Request(PermissionKindLocation)` |
| `NSSupportsSuddenTermination` | den Anfangszustand von `Lifecycle.SuddenTerminationEnabled`; `HoldTermination` setzt es aus |
| `NSSupportsAutomaticTermination` | ermöglicht macOS, die inaktive Anwendung zu beenden; `HoldTermination` setzt es aus |
| `CFBundleLocalizations` | die Sprachen, die `Env.Locale().Identifier` melden kann |
| `NSServices` | einen Eintrag pro mit `app.ServicesProvider` registriertem Dienst; den Block mit `InfoPlistXML` erzeugen |
| `NSUserActivityTypes` | jeden `UserActivity.Type`, der über `app.Activity` veröffentlicht oder fortgesetzt wird |
| `NSAppleScriptEnabled` und `OSAScriptingDefinition` | kennzeichnen die Anwendung als skriptfähig und benennen die mit `app.AppleEvents.ScriptingDefinition` erzeugte `.sdef` |
| `NSAppleEventsUsageDescription` | `app.AppleEvents.Send` an andere Anwendungen |

Für die Berechtigung zu Mitteilungen und für die Spracherkennung muss die Anwendung außerdem als Bundle mit einer Bundle-Kennung ausgeführt werden. Eine nicht gebündelte `go run`-Binärdatei meldet für Mitteilungen `PermissionStatusUnsupported`.

<a id="platform-notes"></a>

## Plattformhinweise

Jede API auf dieser Seite lässt sich unter Windows und Linux kompilieren. Außerhalb von macOS gilt:

- Fensterfunktionen: `SetRepresentedFile`, `SetDocumentEdited`, `SetSubtitle`, `CascadeFrom`, `SetFrameAutosaveName` und `SetWindowButtonsOffset` bewirken nichts. `WindowCascade` verhält sich wie `WindowCentered`. `RequestAttention` gibt ein Handle zurück, dessen `Cancel` nichts bewirkt. `PrintWithOptions` ruft `Print` auf. `ExportPDF` und `Snapshot` geben `ErrMacOnly` zurück.
- Menüs: Symbole, Badges, gemischter Zustand, alternative Elemente und Einrückungen werden gespeichert, aber nicht dargestellt. Abschnittsüberschriften sind deaktivierte Elemente. Paletten sind ausgeblendet. Das Untermenü „Zuletzt geöffnet“ wird bei der Erstellung des Menüs aus der Go-Liste zuletzt verwendeter Dateien gefüllt. Dock-Menüs werden nie angezeigt.
- Dialoge: `SetSuppression`, `SetHelp`, `SetNameFieldLabel` und `SetTags` werden ignoriert. `AddContentType` und `SetFormats` werden zu Dateiendungsfiltern, sofern eine Zuordnung bekannt ist. `Prompt`, `PickColor` und `PickFont` geben `ErrDialogNotSupported` zurück.
- Statusobjekte: `SetSymbol`, `SetRemovable` und `OnVisibilityChange` haben keine Wirkung. `IsVisible` spiegelt den letzten Aufruf von `Show` oder `Hide` wider.
- Feedback: Haptisches Feedback funktioniert unter iOS und Android; unter Windows und Linux bewirken die Aufrufe nichts. `Sound.Beep` und `Sound.Play` funktionieren unter Windows mit WAV-Dateien und Registry-Aliassen. Andernorts gibt `Play` `ErrSoundNotSupported` zurück. Sprachfunktionen geben `ErrSpeechNotSupported` und `ErrSpeechRecognitionNotSupported` zurück.
- Zwischenablage: Die Methoden für vielfältige Inhalte geben `ErrClipboardNotSupported` zurück, `Types` ist leer, `ChangeCount` ist 0 und `OnChange` wird nie ausgelöst.
- Ziehen: `StartDrag` gibt `ErrDragOutUnsupported` zurück. Andere abgelegte Typen als Dateien werden nicht übermittelt; das Ablegen von Dateien funktioniert weiterhin über `WindowFilesDropped`.
- System: `Permissions.Status` meldet `PermissionStatusUnsupported`, und `Request` gibt `ErrPermissionsUnsupported` zurück. `PreventSleep` gibt `ErrPreventSleepUnsupported` zusammen mit einer wirkungslosen Freigabefunktion zurück. `HoldTermination` gibt eine wirkungslose Freigabefunktion zurück. Alle Werte von `Accessibility` sind false, `KeyboardLayout` hat den Nullwert, und `Locale` wird aus `LC_ALL`, `LC_MESSAGES` und `LANG` abgeleitet. Die Fensteroption `Permissions` ist plattformübergreifend.
- Integration: `ServicesProvider.Register` gibt `ErrServicesUnsupported` und `Activity.Publish` gibt `ErrActivityUnsupported` zurück. `InfoPlistXML`, `InfoPlistEntries` und die Aktivitäts-Handler funktionieren weiterhin. `Handle` und `Send` von `AppleEvents` geben `ErrAppleEventsNotSupported` zurück; `ScriptingDefinition` wird überall erzeugt. `QuickLook.Preview` und `Thumbnail` geben `ErrQuickLookNotSupported` zurück. Die Indizierungsmethoden von `Spotlight` geben `ErrSpotlightNotSupported` zurück, und `OnOpen` wird nie ausgelöst. `Browser.OpenWith` startet die angegebene ausführbare Datei mit dem Pfad als Argument, `ApplicationsForFile` ist leer und `ActivateApplication` gibt `ErrApplicationNotRunning` zurück.
- Darstellungsoptionen: `SetPresentationOptions` gibt `ErrMacOnly` zurück, und `PresentationOptions` ist `MacPresentationDefault`.
- Dialogblätter: `PresentSheet`, `PresentCriticalSheet` und `PresentNativeSheet` geben `ErrMacSheetUnsupported` zurück. `EndSheet` bewirkt nichts, und die Abfragemethoden melden kein Dialogblatt.
- Popovers: `NewMacPopover` funktioniert, die Anzeigemethoden geben `ErrMacPopoverUnsupported` zurück, und `IsShown` ist false.
- Zustandswiederherstellung: `SetRestorationID` und `SetRestorationData` bewirken nichts, `OnRestore` wird nie aufgerufen, und `InteractionState` und `RestoreInteractionState` geben `ErrMacOnly` zurück.

## Beispiele

Jedes Beispiel ist eine vollständige, ausführbare Anwendung:

- [`v3/examples/mac-toolbar`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-toolbar): `windowextras.go` integriert den Änderungspunkt, den Untertitel, den PDF-Export, Druckoptionen und gestaffelte Fenster in den Notizeditor.
- [`v3/examples/mac-menus-dock`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-menus-dock): Symbole, Badges, Abschnittsüberschriften, gemischter Zustand, alternative Elemente, eine Palette, „Zuletzt geöffnet“, ein dynamisches Dock-Menü und Dock-Fortschritt.
- [`v3/examples/mac-dialogs`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-dialogs): Unterdrückungs- und Hilfe-Schaltflächen, Texteingaben, Inhaltstypen, das Format-Auswahlmenü, Finder-Tags sowie die Farb- und Schriftfenster.
- [`v3/examples/mac-feedback`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-feedback): ein entfernbares Statusobjekt mit SF Symbol, haptisches Feedback, Systemtöne, Text-to-Speech und Spracherkennung.
- [`v3/examples/mac-clipboard-drag`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-clipboard-drag): die Zwischenablage für vielfältige Inhalte mit Änderungsverfolgung, Herausziehen mit Dateiversprechen sowie das Ablegen von Text, URLs und Bildern.
- [`v3/examples/mac-system`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-system): Berechtigungen, Verhindern des Ruhezustands, Aussetzen des Beendens, Stromversorgungszustand, Bedienungshilfen, Tastaturlayout und Gebietsschema mit laufenden Aktualisierungen.
- [`v3/examples/mac-integration`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-integration): ein Eintrag im Dienste-Menü, eine Handoff-Aktivität mit Fortsetzungs-Handlern und die Workspace-Hilfsfunktionen von `app.Browser`.
- [`v3/examples/mac-search-preview`](https://github.com/wailsapp/wails/tree/master/v3/examples/mac-search-preview): Spotlight-Indizierung mit `OnOpen`, Quick Look-Vorschauen und Vorschaubilder sowie ein eigenes Apple Event mit seiner Skriptdefinition.
