---
title: "Menüs"
description: "Eine Anleitung zum Erstellen und Anpassen von Menüs in Wails v3"
slug: "guides/menus"
sourcePath: "guides/menus.md"
---

Wails v3 bietet ein leistungsstarkes Menüsystem, mit dem Sie sowohl Anwendungsmenüs als auch Kontextmenüs erstellen können. Diese Anleitung führt Sie durch die verschiedenen Funktionen und Möglichkeiten des Menüsystems.

## Ein Menü erstellen

Verwenden Sie zum Erstellen eines neuen Menüs die Methode `New()` des Menüs-Managers:

```go
menu := app.Menu.New()
```

### Menüeinträge hinzufügen

Wails unterstützt mehrere Arten von Menüeinträgen, die jeweils einem bestimmten Zweck dienen:

#### Normale Menüeinträge

Normale Menüeinträge sind die grundlegenden Bausteine von Menüs. Sie zeigen Text an und können beim Anklicken Aktionen auslösen:

```go
menuItem := menu.Add("Click Me")
```

#### Kontrollkästchen

Menüeinträge mit Kontrollkästchen bieten einen umschaltbaren Zustand und eignen sich zum Aktivieren oder Deaktivieren von Funktionen oder Einstellungen:

```go
checkbox := menu.AddCheckbox("My checkbox", true)  // true = initially checked
```

#### Optionsgruppen

Mit Optionsgruppen können Benutzer eine Option aus mehreren sich gegenseitig ausschließenden Optionen auswählen. Sie werden automatisch erstellt, wenn Optionsfelder direkt nebeneinander angeordnet werden:

```go
menu.AddRadio("Option 1", true)   // true = initially selected
menu.AddRadio("Option 2", false)
menu.AddRadio("Option 3", false)
```

#### Trennlinien

Trennlinien sind horizontale Linien, mit denen sich Menüeinträge in logische Gruppen gliedern lassen:

```go
menu.AddSeparator()
```

#### Untermenüs

Untermenüs sind verschachtelte Menüs, die erscheinen, wenn der Mauszeiger über einen Menüeintrag bewegt oder dieser angeklickt wird. Sie eignen sich zum Strukturieren komplexer Menüs:

```go
submenu := menu.AddSubmenu("File")
submenu.Add("Open")
submenu.Add("Save")
```

#### Menüs kombinieren

Ein Menü kann an ein anderes Menü angehängt oder diesem vorangestellt werden.

```go
menu := app.Menu.New()
menu.Add("First Menu")

secondaryMenu := app.Menu.New()
secondaryMenu.Add("Second Menu")

// insert 'secondaryMenu' after 'menu'
menu.Append(secondaryMenu)

// insert 'secondaryMenu' before 'menu'
menu.Prepend(secondaryMenu)

// update the menu
menu.Update()
```

@note{type="info"}
Standardmäßig teilen `prepend` und `append` ihren Zustand mit dem ursprünglichen Menü. Wenn Sie ein neues Menü mit eigenem Zustand erstellen möchten, können Sie für das Menü `.Clone()` aufrufen.

Beispiel: `menu.Append(secondaryMenu.Clone())`

@end

#### Ein Menü leeren

Wenn Sie mit einer variablen Anzahl von Menüeinträgen arbeiten, ist es in manchen Fällen besser, ein vollständig neues Menü zu erstellen.

Dadurch werden alle Einträge eines vorhandenen Menüs entfernt, sodass Sie anschließend wieder Einträge hinzufügen können.

```go
menu := app.Menu.New()
menu.Add("Waiting for update...")

// after certain logic, the menu has to be updated
menu.Clear()
menu.Add("Update complete!")
menu.Update()
```

@note{type="info"}
Beim Leeren eines Menüs werden lediglich die Menüeinträge auf der obersten Ebene entfernt. Untermenüs sind dann zwar nicht mehr sichtbar, belegen aber weiterhin Speicher. Verwalten Sie Ihre Menüs daher sorgfältig.

@end

#### Ein Menü zerstören

Wenn Sie ein Menü leeren und freigeben möchten, verwenden Sie die Methode `Destroy()`:

```go
menu := app.Menu.New()
menu.Add("Waiting for update...")

// after certain logic, the menu has to be destroyed
menu.Destroy()
```

### Eigenschaften von Menüeinträgen

Für Menüeinträge lassen sich mehrere Eigenschaften konfigurieren:

| Eigenschaft | Methode | Beschreibung |
| --- | --- | --- |
| Beschriftung | `SetLabel(string)` | Legt den Anzeigetext fest |
| Aktiviert | `SetEnabled(bool)` | Aktiviert oder deaktiviert den Eintrag |
| Ausgewählt | `SetChecked(bool)` | Legt den Auswahlzustand fest (für Kontrollkästchen und Optionsfelder) |
| Tooltip | `SetTooltip(string)` | Legt den Tooltip-Text fest |
| Ausgeblendet | `SetHidden(bool)` | Blendet den Eintrag ein oder aus |
| Tastenkürzel | `SetAccelerator(string)` | Legt das Tastenkürzel fest |

### Zustände von Menüeinträgen

Menüeinträge können verschiedene Zustände aufweisen, die ihre Sichtbarkeit und Interaktivität steuern:

#### Sichtbarkeit

Menüeinträge können mit der Methode `SetHidden()` dynamisch ein- oder ausgeblendet werden:

```go
menuItem := menu.Add("Dynamic Item")

// Hide the menu item
menuItem.SetHidden(true)

// Show the menu item
menuItem.SetHidden(false)

// Check current visibility
isHidden := menuItem.Hidden()
```

Ausgeblendete Menüeinträge werden vollständig aus dem Menü entfernt, bis sie wieder eingeblendet werden. Dies eignet sich für kontextabhängige Menüeinträge, die nur in bestimmten Anwendungszuständen erscheinen sollen.

#### Aktivierungszustand

Menüeinträge können mit der Methode `SetEnabled()` aktiviert oder deaktiviert werden:

```go
menuItem := menu.Add("Save")

// Disable the menu item
menuItem.SetEnabled(false)  // Item appears grayed out and cannot be clicked

// Enable the menu item
menuItem.SetEnabled(true)   // Item becomes clickable again

// Check current enabled state
isEnabled := menuItem.Enabled()
```

Deaktivierte Menüeinträge bleiben sichtbar, werden jedoch ausgegraut dargestellt und können nicht angeklickt werden. Dies wird üblicherweise verwendet, um anzuzeigen, dass eine Aktion derzeit nicht verfügbar ist, zum Beispiel:

- „Speichern“ deaktivieren, wenn keine zu speichernden Änderungen vorhanden sind
- „Kopieren“ deaktivieren, wenn nichts ausgewählt ist
- „Rückgängig“ deaktivieren, wenn keine Aktion rückgängig gemacht werden kann

#### Dynamische Zustandsverwaltung

Sie können diese Zustände mit Ereignishandlern kombinieren, um dynamische Menüs zu erstellen:

```go
saveMenuItem := menu.Add("Save")

// Initially disable the Save menu item
saveMenuItem.SetEnabled(false)

// Enable Save only when there are unsaved changes
documentChanged := func() {
    saveMenuItem.SetEnabled(true)
    menu.Update()  // Remember to update the menu after changing states
}

// Disable Save after saving
documentSaved := func() {
    saveMenuItem.SetEnabled(false)
    menu.Update()
}
```

### Ereignisbehandlung

Menüeinträge können mit der Methode `OnClick` auf Klickereignisse reagieren:

```go
menuItem.OnClick(func(ctx *application.Context) {
    // Handle the click event
    println("Menu item clicked!")
})
```

Der Kontext enthält Informationen über den angeklickten Menüeintrag:

```go
menuItem.OnClick(func(ctx *application.Context) {
    // Get the clicked menu item
    clickedItem := ctx.ClickedMenuItem()
    // Get its current state
    isChecked := clickedItem.Checked()
})
```

### Rollenbasierte Menüeinträge

Wails stellt vordefinierte Menürollen bereit, die automatisch Menüeinträge mit Standardfunktionen erstellen. Folgende Menürollen werden unterstützt:

#### Vollständige Menüstrukturen

Diese Rollen erstellen vollständige Menüstrukturen mit gängigen Funktionen:

| Rolle | Beschreibung | Plattformhinweise |
| --- | --- | --- |
| `AppMenu` | Anwendungsmenü mit „Über“, „Dienste“, „Ausblenden/Einblenden“ und „Beenden“ | Nur macOS |
| `EditMenu` | Standardmenü „Bearbeiten“ mit Rückgängig, Wiederholen, Ausschneiden, Kopieren, Einfügen usw. | Alle Plattformen |
| `ViewMenu` | Menü „Ansicht“ mit Steuerelementen zum Neuladen, Zoomen und Umschalten in den Vollbildmodus | Alle Plattformen |
| `WindowMenu` | Fenstersteuerung (Minimieren, Zoomen usw.) | Alle Plattformen |
| `HelpMenu` | Hilfemenü mit dem Link „Mehr erfahren“ zur Wails-Website | Alle Plattformen |

#### Einzelne Menüeinträge

Mit diesen Rollen können einzelne Menüeinträge hinzugefügt werden:

| Rolle | Beschreibung | Plattformhinweise |
| --- | --- | --- |
| `About` | „Über“-Dialog der Anwendung anzeigen | Alle Plattformen |
| `Hide` | Anwendung ausblenden | Nur macOS |
| `HideOthers` | Andere Anwendungen ausblenden | Nur macOS |
| `UnHide` | Ausgeblendete Anwendung anzeigen | Nur macOS |
| `CloseWindow` | Aktuelles Fenster schließen | Alle Plattformen |
| `Minimise` | Fenster minimieren | Alle Plattformen |
| `Zoom` | Fenster zoomen | Nur macOS |
| `Front` | Fenster in den Vordergrund bringen | Nur macOS |
| `Quit` | Anwendung beenden | Alle Plattformen |
| `Undo` | Letzte Aktion rückgängig machen | Alle Plattformen |
| `Redo` | Letzte Aktion wiederholen | Alle Plattformen |
| `Cut` | Auswahl ausschneiden | Alle Plattformen |
| `Copy` | Auswahl kopieren | Alle Plattformen |
| `Paste` | Aus der Zwischenablage einfügen | Alle Plattformen |
| `PasteAndMatchStyle` | Einfügen und Formatierung anpassen | Nur macOS |
| `SelectAll` | Alles auswählen | Alle Plattformen |
| `Delete` | Auswahl löschen | Alle Plattformen |
| `Reload` | Aktuelle Seite neu laden | Alle Plattformen |
| `ForceReload` | Neuladen der aktuellen Seite erzwingen | Alle Plattformen |
| `ToggleFullscreen` | Vollbildmodus umschalten | Alle Plattformen |
| `ResetZoom` | Zoomstufe zurücksetzen | Alle Plattformen |
| `ZoomIn` | Vergrößern | Alle Plattformen |
| `ZoomOut` | Verkleinern | Alle Plattformen |

Das folgende Beispiel zeigt, wie du sowohl vollständige Menüs als auch einzelne Rollen verwendest:

```go
menu := app.Menu.New()

// Add complete menu structures
menu.AddRole(application.AppMenu)    // macOS only
menu.AddRole(application.EditMenu)   // Common edit operations
menu.AddRole(application.ViewMenu)   // View controls
menu.AddRole(application.WindowMenu) // Window controls

// Add individual role-based items to a custom menu
fileMenu := menu.AddSubmenu("File")
fileMenu.AddRole(application.CloseWindow)
fileMenu.AddSeparator()
fileMenu.AddRole(application.Quit)
```

## Anwendungsmenüs

Anwendungsmenüs werden oben im Anwendungsfenster (Windows/Linux) oder am oberen Bildschirmrand (macOS) angezeigt.

### Verhalten von Anwendungsmenüs

Wenn du mit `app.Menu.Set()` ein Anwendungsmenü festlegst, wird es unter macOS zum Hauptmenü. Unter Windows/Linux werden Menüs für jedes Fenster separat festgelegt.

```go
app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Custom Menu Window",
    Windows: application.WindowsWindow{
        Menu: customMenu,  // Override application menu for this window
    },
})
```

Das folgende vollständige Beispiel zeigt diese unterschiedlichen Menüverhaltensweisen:

```go
func main() {
    app := application.New(application.Options{})

    // Create application menu
    appMenu := app.Menu.New()
    fileMenu := appMenu.AddSubmenu("File")
    fileMenu.Add("New").OnClick(func(ctx *application.Context) {
        // This will be available in all windows unless overridden
        window := app.Window.Current()
        window.SetTitle("New Window")
    })
    
    // Set as application menu - this is for macOS
    app.Menu.Set(appMenu)

    // Window with custom menu on Windows
    customMenu := app.Menu.New()
    customMenu.Add("Custom Action")
    app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title: "Custom Menu",
        Windows: application.WindowsWindow{
            Menu: customMenu,
        },
    })

    app.Run()
}
```

## Kontextmenüs

Kontextmenüs sind Popup-Menüs, die erscheinen, wenn du mit der rechten Maustaste auf Elemente in deiner Anwendung klickst. Sie ermöglichen den schnellen Zugriff auf Aktionen, die für das angeklickte Element relevant sind.

### Standardkontextmenü

Das Standardkontextmenü ist das integrierte Kontextmenü des Webviews und stellt unter anderem folgende Operationen auf Systemebene bereit:

- Kopieren, Ausschneiden und Einfügen zur Textbearbeitung
- Steuerelemente zur Textauswahl
- Optionen für die Rechtschreibprüfung

#### Standardkontextmenü steuern

Mit der CSS-Eigenschaft `--default-contextmenu` kannst du steuern, wann das Standardkontextmenü angezeigt wird:

```html
<!-- Always show default context menu -->
<div style="--default-contextmenu: show">
    <input type="text" placeholder="Right-click for text operations"/>
    <textarea>Standard text operations available here</textarea>
</div>

<!-- Hide default context menu -->
<div style="--default-contextmenu: hide">
    <div class="custom-component">Custom context menu only</div>
</div>

<!-- Smart context menu behaviour (default) -->
<div style="--default-contextmenu: auto">
    <!-- Shows default menu when text is selected or in input fields -->
    <p>Select this text to see the default menu</p>
    <input type="text" placeholder="Default menu for input operations"/>
</div>
```

@note{type="info"}
Diese Funktion arbeitet erst wie erwartet, nachdem die [Frontend-Laufzeit bereit ist](/reference/frontend-runtime/).

@end

#### Verhalten verschachtelter Kontextmenüs

Wenn du die Eigenschaft `--default-contextmenu` für verschachtelte Elemente verwendest, gelten die folgenden Regeln:

1. Untergeordnete Elemente erben die Kontextmenüeinstellung ihres übergeordneten Elements, sofern sie nicht ausdrücklich überschrieben wird
2. Die spezifischste (nächstgelegene) Einstellung hat Vorrang
3. Mit dem Wert `auto` kannst du das Standardverhalten wiederherstellen

Beispiel für das Verhalten verschachtelter Kontextmenüs:

```html
<!-- Parent sets hide -->
<div style="--default-contextmenu: hide">
    <!-- This inherits hide -->
    <p>No context menu here</p>
    
    <!-- This overrides to show -->
    <div style="--default-contextmenu: show">
        <p>Context menu shown here</p>
        
        <!-- This inherits show -->
        <span>Also has context menu</span>
        
        <!-- This resets to automatic behaviour -->
        <div style="--default-contextmenu: auto">
            <p>Shows menu only when text is selected</p>
        </div>
    </div>
</div>
```

### Benutzerdefinierte Kontextmenüs

Mit benutzerdefinierten Kontextmenüs kannst du anwendungsspezifische Aktionen bereitstellen, die für das angeklickte Element relevant sind. Sie eignen sich besonders für:

- Dateioperationen in einer Dokumentverwaltung
- Werkzeuge zur Bildbearbeitung
- Benutzerdefinierte Aktionen in einem Datenraster
- Komponentenspezifische Operationen

#### Benutzerdefiniertes Kontextmenü erstellen

Beim Erstellen eines benutzerdefinierten Kontextmenüs geben Sie eine eindeutige Kennung (einen Namen) an, die das Menü mit HTML-Elementen verknüpft:

```go
// Create a context menu with identifier "imageMenu"
contextMenu := app.ContextMenu.New()
app.ContextMenu.Add("imageMenu", contextMenu)
```

Der Parameter name (in diesem Beispiel „imageMenu“) dient als eindeutige Kennung für folgende Zwecke:

1. HTML-Elemente mit diesem bestimmten Kontextmenü verknüpfen
2. Bestimmen, welches Menü bei einem Rechtsklick angezeigt werden soll
3. Aktualisieren und Bereinigen des Menüs ermöglichen

#### Kontextdaten

Beim Verarbeiten von Kontextmenüereignissen können Sie sowohl auf den angeklickten Menüeintrag als auch auf die zugehörigen Kontextdaten zugreifen:

```go
contextMenu.Add("Process").OnClick(func(ctx *application.Context) {
    // Get the clicked menu item
    menuItem := ctx.ClickedMenuItem()
    
    // Get the context data as a string
    contextData := ctx.ContextMenuData()
    
    // Check if the menu item is checked (for checkbox/radio items)
    isChecked := ctx.IsChecked()
    
    // Use the data
    if contextData != "" {
        processItem(contextData)
    }
})
```

Die Kontextdaten werden über die Eigenschaft `--custom-contextmenu-data` des HTML-Elements übergeben und sind im Click-Handler über `ctx.ContextMenuData()` verfügbar. Dies ist besonders nützlich beim:

- Arbeiten mit Listen oder Rastern, in denen jeder Eintrag eindeutig identifiziert werden muss
- Ausführen von Operationen für bestimmte Komponenten oder Elemente
- Übergeben von Zustand oder Metadaten vom Frontend an das Backend

#### Kontextmenüverwaltung

Rufen Sie nach Änderungen an einem Kontextmenü die Methode `Update()` auf, um die Änderungen anzuwenden:

```go
contextMenu.Update()
```

Wenn Sie ein Kontextmenü nicht mehr benötigen, können Sie es zerstören:

```go
contextMenu.Destroy()
```

@note{type="danger" title="Warnung"}
Wenn Sie nach dem Aufruf von `Destroy()` die bestehende Kontextmenüreferenz erneut verwenden, tritt eine Panic auf.

@end

### Praxisbeispiel: Bildergalerie

Das folgende vollständige Beispiel zeigt die Implementierung eines benutzerdefinierten Kontextmenüs für eine Bildergalerie:

```go
// Backend: Create the context menu
imageMenu := app.ContextMenu.New()
app.ContextMenu.Add("imageMenu", imageMenu)

// Add relevant operations
imageMenu.Add("View Full Size").OnClick(func(ctx *application.Context) {
    // Get the image ID from context data
    if imageID := ctx.ContextMenuData(); imageID != "" {
        openFullSizeImage(imageID)
    }
})

imageMenu.Add("Download").OnClick(func(ctx *application.Context) {
    if imageID := ctx.ContextMenuData(); imageID != "" {
        downloadImage(imageID)
    }
})

imageMenu.Add("Share").OnClick(func(ctx *application.Context) {
    if imageID := ctx.ContextMenuData(); imageID != "" {
        showShareDialog(imageID)
    }
})
```

```html
<!-- Frontend: Image gallery implementation -->
<div class="gallery">
    <!-- Each image container with context menu -->
    <div class="image-container" 
         style="--custom-contextmenu: imageMenu; --custom-contextmenu-data: img_123">
        <img src="/images/img_123.jpg" alt="Gallery Image"/>
        <span class="caption">Nature Photo</span>
    </div>
    
    <div class="image-container"
         style="--custom-contextmenu: imageMenu; --custom-contextmenu-data: img_124">
        <img src="/images/img_124.jpg" alt="Gallery Image"/>
        <span class="caption">City Photo</span>
    </div>
</div>
```

In diesem Beispiel:

1. Das Kontextmenü wird mit der Kennung „imageMenu“ erstellt
2. Jeder Bildcontainer wird über `--custom-contextmenu: imageMenu` mit dem Menü verknüpft
3. Jeder Container stellt seine Bild-ID über `--custom-contextmenu-data` als Kontextdaten bereit
4. Das Backend empfängt die Bild-ID in Click-Handlern und kann bestimmte Operationen ausführen
5. Dasselbe Menü wird für alle Bilder wiederverwendet; die Kontextdaten geben jedoch an, für welches Bild die Operation ausgeführt werden soll

Dieses Muster eignet sich besonders für:

- Datenraster, deren Zeilen jeweils eigene Operationen erfordern
- Dateimanager, in denen Dateien kontextspezifische Aktionen erfordern
- Designwerkzeuge, in denen unterschiedliche Elemente unterschiedliche Operationen erfordern
- Alle Komponenten, bei denen dieselben Operationen auf mehrere Instanzen angewendet werden
