---
title: "Kontextmenüs"
description: "Erstellen Sie Kontextmenüs für Ihre Anwendung, die per Rechtsklick geöffnet werden"
slug: "features/menus/context"
sourcePath: "features/menus/context.md"
---

## Das Problem

Benutzer erwarten bei einem Rechtsklick Menüs mit kontextspezifischen Aktionen. Unterschiedliche Elemente benötigen unterschiedliche Menüs:

- **Text**: Ausschneiden, Kopieren, Einfügen
- **Bilder**: Speichern, Kopieren, Öffnen
- **Benutzerdefinierte Elemente**: anwendungsspezifische Aktionen

Wenn Sie Kontextmenüs manuell erstellen, müssen Sie Mausereignisse, Positionierung und Plattformunterschiede selbst behandeln.

## Die Wails-Lösung

Wails stellt mithilfe von CSS-Eigenschaften **deklarative Kontextmenüs** bereit. Ordnen Sie Menüs HTML-Elementen zu, übergeben Sie Daten und verarbeiten Sie Klicks – alles mit nativem Plattformverhalten.

![Ein benutzerdefiniertes Wails-Kontextmenü, das unter macOS über einer Webview angezeigt wird](/assets/screenshots/context-menu-macos.png)

Das Menü ist plattformnativ, während das Element, das es geöffnet hat, Teil Ihrer Webview bleibt. Diese macOS-Aufnahme verwendet die Registrierung des benutzerdefinierten Kontextmenüs aus dem folgenden Beispiel.

## Schnellstart

**Go-Code:**

```go
// Create context menu
contextMenu := app.ContextMenu.New()
contextMenu.Add("Cut").OnClick(handleCut)
contextMenu.Add("Copy").OnClick(handleCopy)
contextMenu.Add("Paste").OnClick(handlePaste)

// Register with ID
app.ContextMenu.Add("editor-menu", contextMenu)
```

**HTML:**

```html
<textarea style="--custom-contextmenu: editor-menu">
    Right-click me!
</textarea>
```

**Das ist alles!** Ein Rechtsklick auf den Textbereich zeigt Ihr benutzerdefiniertes Menü an.

## Kontextmenüs erstellen

### Einfaches Kontextmenü

```go
// Create menu
contextMenu := app.ContextMenu.New()

// Add items
contextMenu.Add("Cut").SetAccelerator("CmdOrCtrl+X").OnClick(func(ctx *application.Context) {
    // Handle cut
})

contextMenu.Add("Copy").SetAccelerator("CmdOrCtrl+C").OnClick(func(ctx *application.Context) {
    // Handle copy
})

contextMenu.Add("Paste").SetAccelerator("CmdOrCtrl+V").OnClick(func(ctx *application.Context) {
    // Handle paste
})

// Register with unique ID
app.ContextMenu.Add("text-menu", contextMenu)
```

**Menü-ID:** Muss eindeutig sein. Sie dient dazu, das Menü HTML-Elementen zuzuordnen.

### Mit Untermenüs

```go
contextMenu := app.ContextMenu.New()

// Add regular items
contextMenu.Add("Open").OnClick(handleOpen)
contextMenu.Add("Delete").OnClick(handleDelete)

contextMenu.AddSeparator()

// Add submenu
exportMenu := contextMenu.AddSubmenu("Export As")
exportMenu.Add("PNG").OnClick(exportPNG)
exportMenu.Add("JPEG").OnClick(exportJPEG)
exportMenu.Add("SVG").OnClick(exportSVG)

app.ContextMenu.Add("image-menu", contextMenu)
```

### Mit Kontrollkästchen und Optionsgruppen

```go
contextMenu := app.ContextMenu.New()

// Checkbox
contextMenu.AddCheckbox("Show Grid", true).OnClick(func(ctx *application.Context) {
    showGrid := ctx.ClickedMenuItem().Checked()
    // Toggle grid
})

contextMenu.AddSeparator()

// Radio group
contextMenu.AddRadio("Small", false).OnClick(handleSize)
contextMenu.AddRadio("Medium", true).OnClick(handleSize)
contextMenu.AddRadio("Large", false).OnClick(handleSize)

app.ContextMenu.Add("view-menu", contextMenu)
```

Informationen zu **allen Menüelementtypen** finden Sie in der [Menüreferenz](/features/menus/reference/).

## HTML-Elementen zuordnen

Verwenden Sie benutzerdefinierte CSS-Eigenschaften, um Kontextmenüs zuzuweisen:

### Einfache Zuordnung

```html
<div style="--custom-contextmenu: menu-id">
    Right-click me!
</div>
```

**CSS-Eigenschaft:** `--custom-contextmenu: <menu-id>`

### Mit Kontextdaten

Übergeben Sie Daten aus HTML an Go:

```html
<div style="--custom-contextmenu: file-menu; --custom-contextmenu-data: file-123">
    Right-click this file
</div>
```

**Go-Handler:**

```go
contextMenu := app.ContextMenu.New()
contextMenu.Add("Open").OnClick(func(ctx *application.Context) {
    fileID := ctx.ContextMenuData()  // "file-123"
    openFile(fileID)
})

app.ContextMenu.Add("file-menu", contextMenu)
```

**CSS-Eigenschaften:**

- `--custom-contextmenu: <menu-id>` – Das anzuzeigende Menü
- `--custom-contextmenu-data: <data>` – An Handler zu übergebende Daten

### Dynamische Daten

Generieren Sie Daten dynamisch in JavaScript:

```html
<div id="file-item" style="--custom-contextmenu: file-menu">
    File.txt
</div>

<script>
// Set data dynamically
const fileItem = document.getElementById('file-item')
fileItem.style.setProperty('--custom-contextmenu-data', 'file-' + fileId)
</script>
```

### Mehrere Elemente, dasselbe Menü

```html
<div class="file-item" style="--custom-contextmenu: file-menu; --custom-contextmenu-data: file-1">
    Document.pdf
</div>

<div class="file-item" style="--custom-contextmenu: file-menu; --custom-contextmenu-data: file-2">
    Image.png
</div>

<div class="file-item" style="--custom-contextmenu: file-menu; --custom-contextmenu-data: file-3">
    Video.mp4
</div>
```

**Ein Menü, unterschiedliche Daten für jedes Element.**

## Kontextdaten

### Auf Kontextdaten zugreifen

```go
contextMenu.Add("Process").OnClick(func(ctx *application.Context) {
    data := ctx.ContextMenuData()  // Get data from HTML
    
    // Use the data
    processItem(data)
})
```

**Datentyp:** Immer `string`. Parsen Sie die Daten nach Bedarf.

### Komplexe Daten übergeben

Verwenden Sie JSON für komplexe Daten:

```html
<div style="--custom-contextmenu: item-menu; --custom-contextmenu-data: {&quot;id&quot;:123,&quot;type&quot;:&quot;image&quot;}">
    Image.png
</div>
```

**Go-Handler:**

```go
import "encoding/json"

type ItemData struct {
    ID   int    `json:"id"`
    Type string `json:"type"`
}

contextMenu.Add("Process").OnClick(func(ctx *application.Context) {
    dataStr := ctx.ContextMenuData()
    
    var data ItemData
    if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
        log.Printf("Invalid data: %v", err)
        return
    }
    
    processItem(data.ID, data.Type)
})
```

@note{type="caution" title="Sicherheit"}
**Validieren Sie Kontextdaten aus dem Frontend immer**. Benutzer können CSS-Eigenschaften manipulieren; behandeln Sie die Daten daher als nicht vertrauenswürdige Eingabe.

@end

### Validierungsbeispiel

```go
contextMenu.Add("Delete").OnClick(func(ctx *application.Context) {
    fileID := ctx.ContextMenuData()
    
    // Validate
    if !isValidFileID(fileID) {
        log.Printf("Invalid file ID: %s", fileID)
        return
    }
    
    // Check permissions
    if !canDeleteFile(fileID) {
        showError("Permission denied")
        return
    }
    
    // Safe to proceed
    deleteFile(fileID)
})
```

## Standardkontextmenü

Die WebView stellt ein integriertes Kontextmenü für Standardoperationen wie Kopieren, Einfügen und Untersuchen bereit. Steuern Sie es mit `--default-contextmenu`:

### Standardmenü ausblenden

```html
<div style="--default-contextmenu: hide">
    No default menu here
</div>
```

**Anwendungsfall:** Benutzerdefinierte UI-Elemente, für die das Standardmenü nicht sinnvoll ist.

### Standardmenü anzeigen

```html
<div style="--default-contextmenu: show">
    Default menu always shown
</div>
```

**Anwendungsfall:** Textbereiche, Eingabefelder und bearbeitbare Inhalte.

### Automatischer (intelligenter) Modus

```html
<div style="--default-contextmenu: auto">
    Smart context menu
</div>
```

**Standardverhalten.** Zeigt das Standardmenü an, wenn:

- Text ausgewählt ist
- In Texteingabefeldern
- In bearbeitbaren Inhalten (`contenteditable`)

Blendet das Standardmenü andernfalls aus.

### Benutzerdefiniertes und Standardmenü kombinieren

```html
<!-- Custom menu + default menu -->
<textarea style="--custom-contextmenu: editor-menu; --default-contextmenu: show">
    Both menus available
</textarea>
```

**Verhalten:**

1. Das benutzerdefinierte Menü wird zuerst angezeigt
2. Wenn das benutzerdefinierte Menü leer ist oder nicht gefunden wird, wird das Standardmenü angezeigt
3. Beide können gleichzeitig vorhanden sein (plattformabhängig)

## Dynamische Kontextmenüs

Aktualisieren Sie Menüs anhand des Anwendungszustands:

### Menüelemente aktivieren/deaktivieren

```go
var cutMenuItem *application.MenuItem
var copyMenuItem *application.MenuItem

func createContextMenu() {
    contextMenu := app.ContextMenu.New()
    
    cutMenuItem = contextMenu.Add("Cut")
    cutMenuItem.SetEnabled(false)  // Initially disabled
    cutMenuItem.OnClick(handleCut)
    
    copyMenuItem = contextMenu.Add("Copy")
    copyMenuItem.SetEnabled(false)
    copyMenuItem.OnClick(handleCopy)
    
    app.ContextMenu.Add("editor-menu", contextMenu)
}

func onSelectionChanged(hasSelection bool) {
    cutMenuItem.SetEnabled(hasSelection)
    copyMenuItem.SetEnabled(hasSelection)
    contextMenu.Update()  // Important!
}
```

@note{type="caution" title="Update() immer aufrufen"}
Rufen Sie nach einer Änderung des Menüzustands **`contextMenu.Update()`** auf. Dies ist unter Windows von entscheidender Bedeutung.

Weitere Informationen finden Sie in der [Menüreferenz](/features/menus/reference/#enabled-state).

@end

### Beschriftungen ändern

```go
playMenuItem := contextMenu.Add("Play")

playMenuItem.OnClick(func(ctx *application.Context) {
    if isPlaying {
        playMenuItem.SetLabel("Pause")
    } else {
        playMenuItem.SetLabel("Play")
    }
    contextMenu.Update()
})
```

### Menüs neu erstellen

Erstellen Sie bei umfangreichen Änderungen das gesamte Menü neu:

```go
func rebuildContextMenu(fileType string) {
    contextMenu := app.ContextMenu.New()
    
    // Common items
    contextMenu.Add("Open").OnClick(handleOpen)
    contextMenu.Add("Delete").OnClick(handleDelete)
    
    contextMenu.AddSeparator()
    
    // Type-specific items
    switch fileType {
    case "image":
        contextMenu.Add("Edit Image").OnClick(editImage)
        contextMenu.Add("Set as Wallpaper").OnClick(setWallpaper)
    case "video":
        contextMenu.Add("Play").OnClick(playVideo)
        contextMenu.Add("Extract Audio").OnClick(extractAudio)
    case "document":
        contextMenu.Add("Print").OnClick(printDocument)
        contextMenu.Add("Export PDF").OnClick(exportPDF)
    }
    
    app.ContextMenu.Add("file-menu", contextMenu)
}
```

## Plattformverhalten

Kontextmenüs sind **plattformnativ**:

@tabs{sync-key="platform"}
[macOS]
**Native macOS-Kontextmenüs:**

- Systemanimationen und Übergänge
- Rechtsklick = Control+Klick (automatisch)
- Passt sich dem Erscheinungsbild des Systems an (hell/dunkel)
- Standardmäßige Textoperationen im Standardmenü
- Natives Scrollen bei langen Menüs

**macOS-Konventionen:**

- Verwenden Sie für Menüeinträge Satzschreibung
- Verwenden Sie für Einträge, die Dialogfelder öffnen, Auslassungspunkte (...)
- Gängige Tastenkürzel: ⌘C (Kopieren), ⌘V (Einfügen)

[Windows]
**Native Windows-Kontextmenüs:**

- Nativer Windows-Stil
- Entspricht dem Windows-Design
- Standardmäßige Windows-Operationen im Standardmenü
- Unterstützung für Touch- und Stifteingaben

**Windows-Konventionen:**

- Schreiben Sie bei Menüeinträgen die wichtigen Wörter groß
- Verwenden Sie für Einträge, die Dialogfelder öffnen, Auslassungspunkte (...)
- Gängige Tastenkürzel: Ctrl+C (Kopieren), Ctrl+V (Einfügen)

[Linux]
**Integration in die Desktop-Umgebung:**

- Passt sich dem Desktop-Design an (GTK, Qt usw.)
- Das Rechtsklickverhalten entspricht den Systemeinstellungen
- Der Inhalt des Standardmenüs variiert je nach Umgebung
- Die Positionierung entspricht den Konventionen der Desktop-Umgebung

**Hinweise zu Linux:**

- Testen Sie auf den Ziel-Desktop-Umgebungen
- GTK und Qt verhalten sich unterschiedlich
- Einige Desktop-Umgebungen passen Kontextmenüs an

@end

## Vollständiges Beispiel

**Go-Code:**

```go
package main

import (
    "encoding/json"
    "log"
    "github.com/wailsapp/wails/v3/pkg/application"
)

type FileData struct {
    ID   string `json:"id"`
    Type string `json:"type"`
    Name string `json:"name"`
}

func main() {
    app := application.New(application.Options{
        Name: "Context Menu Demo",
    })

    // Create file context menu
    fileMenu := createFileMenu(app)
    app.ContextMenu.Add("file-menu", fileMenu)

    // Create image context menu
    imageMenu := createImageMenu(app)
    app.ContextMenu.Add("image-menu", imageMenu)

    // Create text context menu
    textMenu := createTextMenu(app)
    app.ContextMenu.Add("text-menu", textMenu)

    app.Window.New()
    app.Run()
}

func createFileMenu(app *application.App) *application.ContextMenu {
    menu := app.ContextMenu.New()
    
    menu.Add("Open").OnClick(func(ctx *application.Context) {
        data := parseFileData(ctx.ContextMenuData())
        openFile(data.ID)
    })
    
    menu.Add("Rename").OnClick(func(ctx *application.Context) {
        data := parseFileData(ctx.ContextMenuData())
        renameFile(data.ID)
    })
    
    menu.AddSeparator()
    
    menu.Add("Delete").OnClick(func(ctx *application.Context) {
        data := parseFileData(ctx.ContextMenuData())
        deleteFile(data.ID)
    })
    
    return menu
}

func createImageMenu(app *application.App) *application.ContextMenu {
    menu := app.ContextMenu.New()
    
    menu.Add("View").OnClick(func(ctx *application.Context) {
        data := parseFileData(ctx.ContextMenuData())
        viewImage(data.ID)
    })
    
    menu.Add("Edit").OnClick(func(ctx *application.Context) {
        data := parseFileData(ctx.ContextMenuData())
        editImage(data.ID)
    })
    
    menu.AddSeparator()
    
    exportMenu := menu.AddSubmenu("Export As")
    exportMenu.Add("PNG").OnClick(exportPNG)
    exportMenu.Add("JPEG").OnClick(exportJPEG)
    exportMenu.Add("WebP").OnClick(exportWebP)
    
    return menu
}

func createTextMenu(app *application.App) *application.ContextMenu {
    menu := app.ContextMenu.New()
    
    menu.Add("Cut").SetAccelerator("CmdOrCtrl+X").OnClick(handleCut)
    menu.Add("Copy").SetAccelerator("CmdOrCtrl+C").OnClick(handleCopy)
    menu.Add("Paste").SetAccelerator("CmdOrCtrl+V").OnClick(handlePaste)
    
    return menu
}

func parseFileData(dataStr string) FileData {
    var data FileData
    if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
        log.Printf("Invalid file data: %v", err)
    }
    return data
}

// Handler implementations...
func openFile(id string) { /* ... */ }
func renameFile(id string) { /* ... */ }
func deleteFile(id string) { /* ... */ }
func viewImage(id string) { /* ... */ }
func editImage(id string) { /* ... */ }
func exportPNG(ctx *application.Context) { /* ... */ }
func exportJPEG(ctx *application.Context) { /* ... */ }
func exportWebP(ctx *application.Context) { /* ... */ }
func handleCut(ctx *application.Context) { /* ... */ }
func handleCopy(ctx *application.Context) { /* ... */ }
func handlePaste(ctx *application.Context) { /* ... */ }
```

**HTML:**

```html
<!DOCTYPE html>
<html>
<head>
    <style>
        .file-item {
            padding: 10px;
            margin: 5px;
            border: 1px solid #ccc;
            cursor: pointer;
        }
        
        .file-item:hover {
            background: #f0f0f0;
        }
        
        textarea {
            width: 100%;
            height: 200px;
        }
    </style>
</head>
<body>
    <h2>Files</h2>
    
    <!-- Regular file -->
    <div class="file-item" 
         style="--custom-contextmenu: file-menu; 
                --custom-contextmenu-data: {&quot;id&quot;:&quot;file-1&quot;,&quot;type&quot;:&quot;document&quot;,&quot;name&quot;:&quot;Report.pdf&quot;}">
        📄 Report.pdf
    </div>
    
    <!-- Image file -->
    <div class="file-item" 
         style="--custom-contextmenu: image-menu; 
                --custom-contextmenu-data: {&quot;id&quot;:&quot;file-2&quot;,&quot;type&quot;:&quot;image&quot;,&quot;name&quot;:&quot;Photo.jpg&quot;}">
        🖼️ Photo.jpg
    </div>
    
    <h2>Text Editor</h2>
    
    <!-- Text area with custom menu + default menu -->
    <textarea 
        style="--custom-contextmenu: text-menu; --default-contextmenu: show"
        placeholder="Type here, then right-click...">
    </textarea>
    
    <h2>No Context Menu</h2>
    
    <!-- Disable default menu -->
    <div style="--default-contextmenu: hide; padding: 20px; border: 1px solid #ccc;">
        Right-click here - no menu appears
    </div>
</body>
</html>
```

## Bewährte Vorgehensweisen

### ✅ Empfohlen

- **Menüs übersichtlich halten** – Nur für das Element relevante Aktionen anbieten
- **Kontextdaten validieren** – Als nicht vertrauenswürdige Eingabe behandeln
- **Eindeutige Beschriftungen verwenden** – „Datei löschen“ statt „Löschen“
- **menu.Update() aufrufen** – Nachdem Sie den Menüzustand geändert haben
- **Auf allen Plattformen testen** – Das Verhalten variiert
- **Tastenkürzel bereitstellen** – Für häufige Aktionen
- **Zusammengehörige Einträge gruppieren** – Trennlinien verwenden

### ❌ Nicht empfohlen

- **Kontextdaten nicht vertrauen** – Immer validieren
- **Menüs nicht zu lang machen** – Höchstens 7-10 Einträge
- **menu.Update() nicht vergessen** – Andernfalls funktionieren Menüs nicht ordnungsgemäß
- **Nicht zu tief verschachteln** – Höchstens 2 Ebenen
- **Keine Fachsprache verwenden** – Beschriftungen benutzerfreundlich halten
- **Handler nicht blockieren** – Handler schnell halten

## Fehlerbehebung

### Kontextmenü wird nicht angezeigt

**Mögliche Ursachen:**

1. Menü-ID stimmt nicht überein
2. Tippfehler in der CSS-Eigenschaft
3. Runtime nicht initialisiert

**Lösung:**

```go
// Check menu is registered
app.ContextMenu.Add("my-menu", contextMenu)
```

```html
<!-- Check ID matches -->
<div style="--custom-contextmenu: my-menu">
```

### Kontextdaten werden nicht empfangen

**Mögliche Ursachen:**

1. CSS-Eigenschaft ist nicht festgelegt
2. Daten enthalten Sonderzeichen

**Lösung:**

```html
<!-- Escape quotes in JSON -->
<div style="--custom-contextmenu-data: {&quot;id&quot;:123}">
```

Oder verwenden Sie JavaScript:

```javascript
element.style.setProperty('--custom-contextmenu-data', JSON.stringify(data))
```

### Menüeinträge reagieren nicht

**Ursache:** Nach dem Aktivieren wurde der Aufruf von `menu.Update()` vergessen

**Lösung:**

```go
menuItem.SetEnabled(true)
contextMenu.Update()  // Add this!
```

## Nächste Schritte

@cards{cols="2"}
📖 Menüreferenz
Vollständige Referenz zu Menüelementtypen und -eigenschaften.

[Mehr erfahren →](/features/menus/reference/)

---
☰ Anwendungsmenüs
Erstellen Sie Menüleisten für Anwendungen.

[Mehr erfahren →](/features/menus/application/)

---
★ Infobereichsmenüs
Fügen Sie eine Integration in den Infobereich beziehungsweise die Menüleiste hinzu.

[Mehr erfahren →](/features/menus/systray/)

---
📖 Menümuster
Gängige Menümuster und bewährte Vorgehensweisen.

[Mehr erfahren →](/guides/menus/)

@end

---

**Fragen?** Fragen Sie auf [Discord](https://discord.gg/JDdSxwjhGf) nach oder sehen Sie sich das [Kontextmenübeispiel](https://github.com/wailsapp/wails/tree/master/v3/examples/contextmenus) an.
