---
title: "Menü-API"
description: "Vollständige Referenz zur Menü-API"
slug: "reference/menu"
sourcePath: "reference/menu.md"
---

## Übersicht

Die Menü-API bietet Methoden zum Erstellen und Verwalten von Anwendungsmenüs, Kontextmenüs und Menüs im Infobereich.

**Menütypen:**

- **Anwendungsmenüs** – Obere Menüleiste (Datei, Bearbeiten usw.)
- **Kontextmenüs** – Menüs für Rechtsklicks
- **Infobereichsmenüs** – Menüs im Infobereich der Taskleiste

## Menüs erstellen

### NewMenu()

Erstellt ein neues Menü.

```go
func (a *App) NewMenu() *Menu
```

**Beispiel:**

```go
menu := app.NewMenu()
```

## Menümethoden

### Add()

Fügt dem Menü einen Menüeintrag hinzu.

```go
func (m *Menu) Add(label string) *MenuItem
```

**Parameter:**

- `label` – Der für den Menüeintrag angezeigte Text

**Rückgabewert:** Der erstellte Menüeintrag

**Beispiel:**

```go
item := menu.Add("Open File")
item.OnClick(func(ctx *application.Context) {
    // Handle click
})
```

### AddSubmenu()

Fügt dem Menü ein Untermenü hinzu.

```go
func (m *Menu) AddSubmenu(label string) *Menu
```

**Parameter:**

- `label` – Die Beschriftung des Untermenüs

**Rückgabewert:** Das erstellte Untermenü

**Beispiel:**

```go
fileMenu := menu.AddSubmenu("File")
fileMenu.Add("New")
fileMenu.Add("Open")
fileMenu.Add("Save")
```

### AddSeparator()

Fügt zwischen Menüeinträgen eine sichtbare Trennlinie hinzu.

```go
func (m *Menu) AddSeparator()
```

**Beispiel:**

```go
menu.Add("Copy")
menu.Add("Paste")
menu.AddSeparator()
menu.Add("Select All")
```

**Bewährte Vorgehensweise:** Gruppieren Sie zusammengehörige Menüeinträge mit Trennlinien.

### AddCheckbox()

Fügt einen Menüeintrag mit umschaltbarem Häkchen hinzu.

```go
func (m *Menu) AddCheckbox(label string, checked bool) *MenuItem
```

**Parameter:**

- `label` – Die Beschriftung des Kontrollkästchens
- `checked` – Der anfängliche Auswahlzustand

**Beispiel:**

```go
darkMode := menu.AddCheckbox("Dark Mode", false)
darkMode.OnClick(func(ctx *application.Context) {
    isChecked := darkMode.Checked()
    // Toggle dark mode
})
```

### AddRadio()

Fügt einen Optionsfeld-Menüeintrag hinzu (Gruppe mit gegenseitigem Ausschluss).

```go
func (m *Menu) AddRadio(label string, checked bool) *MenuItem
```

**Parameter:**

- `label` – Die Beschriftung des Optionsfelds
- `checked` – Der anfängliche Auswahlzustand

**Beispiel:**

```go
// Create radio group for view modes
viewMenu := menu.AddSubmenu("View")
listView := viewMenu.AddRadio("List View", true)
gridView := viewMenu.AddRadio("Grid View", false)
treeView := viewMenu.AddRadio("Tree View", false)

listView.OnClick(func(ctx *application.Context) {
    setViewMode("list")
})
gridView.OnClick(func(ctx *application.Context) {
    setViewMode("grid")
})
```

### Update()

Aktualisiert das Menü, sodass Änderungen an Menüeinträgen übernommen werden.

```go
func (m *Menu) Update()
```

**Beispiel:**

```go
item.SetEnabled(false)
menu.Update()  // Must call to apply changes
```

**Wichtig:** Rufen Sie nach dem Ändern von Eigenschaften eines Menüeintrags immer `Update()` auf.

## Methoden für Menüeinträge

### OnClick()

Registriert einen Klick-Handler für den Menüeintrag.

```go
func (mi *MenuItem) OnClick(callback func(ctx *application.Context)) *MenuItem
```

**Parameter:**

- `callback` – Funktion, die beim Klicken auf den Eintrag aufgerufen wird

**Rückgabewert:** Der Menüeintrag (zur Verkettung von Methodenaufrufen)

**Beispiel:**

```go
item.OnClick(func(ctx *application.Context) {
    fmt.Println("Menu item clicked")
    app.Logger.Info("User clicked menu item")
})
```

### SetLabel()

Ändert die Beschriftung des Menüeintrags.

```go
func (mi *MenuItem) SetLabel(label string) *MenuItem
```

**Beispiel:**

```go
item.SetLabel("Save As...")
menu.Update()
```

### SetEnabled()

Aktiviert oder deaktiviert den Menüeintrag.

```go
func (mi *MenuItem) SetEnabled(enabled bool) *MenuItem
```

**Beispiel:**

```go
// Disable save when no document is open
saveItem.SetEnabled(hasOpenDocument)
menu.Update()
```

**Gängiges Muster:**

```go
// Update menu state based on application state
func updateMenuState() {
    saveItem.SetEnabled(hasUnsavedChanges)
    undoItem.SetEnabled(canUndo)
    redoItem.SetEnabled(canRedo)
    menu.Update()
}
```

### SetChecked()

Legt den Auswahlzustand von Kontrollkästchen- und Optionsfeld-Menüeinträgen fest.

```go
func (mi *MenuItem) SetChecked(checked bool) *MenuItem
```

**Beispiel:**

```go
darkModeItem.SetChecked(isDarkModeEnabled)
menu.Update()
```

### Checked()

Gibt den aktuellen Auswahlzustand zurück.

```go
func (mi *MenuItem) Checked() bool
```

**Beispiel:**

```go
if darkModeItem.Checked() {
    // Dark mode is enabled
}
```

### SetAccelerator()

Legt ein Tastenkürzel für den Menüeintrag fest.

```go
func (mi *MenuItem) SetAccelerator(accelerator string) *MenuItem
```

**Parameter:**

- `accelerator` – Tastenkürzel (z. B. „Ctrl+S“, „Cmd+Q“)

**Format des Tastenkürzels:**

- **Modifikatortasten:** `Ctrl`, `Cmd`, `Alt`, `Shift`
- **Tasten:** `A-Z`, `0-9`, `F1-F12`, `Enter`, `Backspace` usw.
- **Plattform:** Verwenden Sie unter macOS `Cmd` und unter Windows/Linux `Ctrl`.

**Beispiel:**

```go
saveItem.SetAccelerator("Ctrl+S")
quitItem.SetAccelerator("Ctrl+Q")
newItem.SetAccelerator("Ctrl+N")
```

**Plattformspezifisches Beispiel:**

```go
import "runtime"

var quitShortcut string
if runtime.GOOS == "darwin" {
    quitShortcut = "Cmd+Q"
} else {
    quitShortcut = "Ctrl+Q"
}
quitItem.SetAccelerator(quitShortcut)
```

### SetTooltip()

Legt einen Tooltip fest, der beim Bewegen des Mauszeigers über den Menüeintrag angezeigt wird.

```go
func (mi *MenuItem) SetTooltip(tooltip string) *MenuItem
```

**Beispiel:**

```go
item.SetTooltip("Opens a file from disk")
```

### SetHidden()

Blendet den Menüeintrag ein oder aus.

```go
func (mi *MenuItem) SetHidden(hidden bool) *MenuItem
```

**Beispiel:**

```go
// Hide debug menu in production
debugItem.SetHidden(!isDevelopment)
menu.Update()
```

## Anwendungsmenü

### app.Menu.Set()

Legt die Hauptmenüleiste der Anwendung fest.

```go
func (mm *MenuManager) Set(menu *Menu)
```

**Beispiel:**

```go
menu := app.NewMenu()

// File menu
fileMenu := menu.AddSubmenu("File")
fileMenu.Add("New").SetAccelerator("Ctrl+N").OnClick(newFile)
fileMenu.Add("Open").SetAccelerator("Ctrl+O").OnClick(openFile)
fileMenu.Add("Save").SetAccelerator("Ctrl+S").OnClick(saveFile)
fileMenu.AddSeparator()
fileMenu.Add("Exit").SetAccelerator("Ctrl+Q").OnClick(func(ctx *application.Context) {
    app.Quit()
})

// Edit menu
editMenu := menu.AddSubmenu("Edit")
editMenu.Add("Undo").SetAccelerator("Ctrl+Z").OnClick(undo)
editMenu.Add("Redo").SetAccelerator("Ctrl+Y").OnClick(redo)
editMenu.AddSeparator()
editMenu.Add("Cut").SetAccelerator("Ctrl+X").OnClick(cut)
editMenu.Add("Copy").SetAccelerator("Ctrl+C").OnClick(copy)
editMenu.Add("Paste").SetAccelerator("Ctrl+V").OnClick(paste)

app.Menu.Set(menu)
```

**Plattformhinweise:**

- **macOS:** Das Menü wird in der oberen Menüleiste angezeigt
- **Windows/Linux:** Das Menü wird in der Titelleiste des Fensters angezeigt
- **macOS:** Fügt automatisch ein Anwendungsmenü mit dem Namen der Anwendung hinzu

## Kontextmenüs

### app.ContextMenu.New() / app.ContextMenu.Add()

Erstellen Sie über den Manager ein `*ContextMenu` und registrieren Sie es unter einem Namen. `ContextMenuManager.Add` akzeptiert `*ContextMenu` – **nicht** `*Menu` – und es gibt **keine** Methode `app.RegisterContextMenu`.

```go
func (cm *ContextMenuManager) New() *ContextMenu
func (cm *ContextMenuManager) Add(name string, menu *ContextMenu)
func (cm *ContextMenuManager) Get(name string) (*ContextMenu, bool)
func (cm *ContextMenuManager) Remove(name string)
```

Alternativ erstellt und registriert das paketweite `application.NewContextMenu(name string) *ContextMenu` ein Kontextmenü in einem Schritt.

**Go:**

```go
// Build the context menu via the manager
contextMenu := app.ContextMenu.New()
contextMenu.Add("Cut").OnClick(cut)
contextMenu.Add("Copy").OnClick(copy)
contextMenu.Add("Paste").OnClick(paste)
contextMenu.AddSeparator()
contextMenu.Add("Select All").OnClick(selectAll)

// Register under a name; HTML opts into it via the CSS custom property below.
app.ContextMenu.Add("editor", contextMenu)
```

**HTML/CSS:**

Die Runtime löst ein registriertes Kontextmenü aus, wenn für das Ziel des Rechtsklicks (oder ein übergeordnetes Element) die benutzerdefinierte CSS-Eigenschaft `--custom-contextmenu` auf diesen Namen gesetzt ist. Die optionale Eigenschaft `--custom-contextmenu-data` wird über `ctx.ContextMenuData()` an den Go-Callback übergeben. Um das standardmäßige Browser-Kontextmenü zu unterdrücken, setzen Sie `--default-contextmenu: hide` (oder `auto`/`show`).

```html
<!-- Trigger context menu on right-click -->
<div style="--custom-contextmenu: editor; --default-contextmenu: hide">
    Right-click here for context menu
</div>
```

Es gibt kein Attribut `data-wails-context-menu="..."` – es war nie mit der Runtime verbunden.

**Dynamische Kontextmenüs:**

```go
// Update context menu based on selection
func updateContextMenu() {
    contextMenu := app.ContextMenu.New()

    if hasSelection {
        contextMenu.Add("Cut").OnClick(cut)
        contextMenu.Add("Copy").OnClick(copy)
    }

    contextMenu.Add("Paste").SetEnabled(hasClipboardContent).OnClick(paste)

    app.ContextMenu.Add("editor", contextMenu)
}
```

## Menü im Infobereich

### app.SystemTray.New()

Erstellt ein neues Symbol im Infobereich.

```go
func (sm *SystemTrayManager) New() *SystemTray
```

**Beispiel:**

```go
tray := app.SystemTray.New()
```

### SetIcon()

Legt das Symbol im Infobereich fest.

```go
func (st *SystemTray) SetIcon(icon []byte) *SystemTray
```

**Beispiel:**

```go
iconData, _ := os.ReadFile("icon.png")
tray.SetIcon(iconData)
```

### SetMenu()

Legt das Menü für den Infobereich fest.

```go
func (st *SystemTray) SetMenu(menu *Menu) *SystemTray
```

**Beispiel:**

```go
trayMenu := app.NewMenu()
trayMenu.Add("Show Window").OnClick(func(ctx *application.Context) {
    window.Show()
    window.Focus()
})
trayMenu.Add("Settings").OnClick(openSettings)
trayMenu.AddSeparator()
trayMenu.Add("Quit").OnClick(func(ctx *application.Context) {
    app.Quit()
})

tray.SetMenu(trayMenu)
```

### SetTooltip()

Legt den Tooltip fest, der beim Bewegen des Mauszeigers über das Symbol im Infobereich angezeigt wird. Gibt nichts zurück.

```go
func (st *SystemTray) SetTooltip(tooltip string)
```

**Beispiel:**

```go
tray.SetTooltip("My Application - Running")
```

### OnClick()

Verarbeitet Linksklicks auf das Symbol im Infobereich.

```go
func (st *SystemTray) OnClick(callback func()) *SystemTray
```

**Beispiel:**

```go
tray.OnClick(func() {
    if window.IsVisible() {
        window.Hide()
    } else {
        window.Show()
        window.Focus()
    }
})
```

## Vollständige Beispiele

### Standard-Anwendungsmenü

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
)

func createMenu(app *application.App) *application.Menu {
    menu := app.NewMenu()

    // File menu
    fileMenu := menu.AddSubmenu("File")
    fileMenu.Add("New").
        SetAccelerator("Ctrl+N").
        OnClick(func(ctx *application.Context) {
            // Create new document
        })
    fileMenu.Add("Open").
        SetAccelerator("Ctrl+O").
        OnClick(func(ctx *application.Context) {
            // Open file dialog
        })
    fileMenu.Add("Save").
        SetAccelerator("Ctrl+S").
        OnClick(func(ctx *application.Context) {
            // Save document
        })
    fileMenu.AddSeparator()
    fileMenu.Add("Exit").OnClick(func(ctx *application.Context) {
        app.Quit()
    })

    // Edit menu
    editMenu := menu.AddSubmenu("Edit")
    editMenu.Add("Undo").SetAccelerator("Ctrl+Z")
    editMenu.Add("Redo").SetAccelerator("Ctrl+Y")
    editMenu.AddSeparator()
    editMenu.Add("Cut").SetAccelerator("Ctrl+X")
    editMenu.Add("Copy").SetAccelerator("Ctrl+C")
    editMenu.Add("Paste").SetAccelerator("Ctrl+V")

    // View menu
    viewMenu := menu.AddSubmenu("View")
    darkMode := viewMenu.AddCheckbox("Dark Mode", false)
    darkMode.OnClick(func(ctx *application.Context) {
        // Toggle dark mode
        isChecked := darkMode.Checked()
        app.Logger.Info("Dark mode", "enabled", isChecked)
    })
    viewMenu.AddSeparator()
    viewMenu.AddRadio("List View", true)
    viewMenu.AddRadio("Grid View", false)
    viewMenu.AddRadio("Detail View", false)

    // Help menu
    helpMenu := menu.AddSubmenu("Help")
    helpMenu.Add("Documentation").OnClick(func(ctx *application.Context) {
        // Open docs
    })
    helpMenu.Add("About").OnClick(func(ctx *application.Context) {
        // Show about dialog
    })

    return menu
}

func main() {
    app := application.New(application.Options{
        Name: "Menu Demo",
    })

    menu := createMenu(app)
    app.Menu.Set(menu)

    window := app.Window.New()
    window.Show()

    app.Run()
}
```

### Anwendung mit Symbol im Infobereich

```go
func setupSystemTray(app *application.App, window application.Window) {
    // Create system tray
    tray := app.SystemTray.New()

    // Set icon
    iconData, _ := os.ReadFile("icon.png")
    tray.SetIcon(iconData)
    tray.SetTooltip("My App - Running")

    // Handle left-click on tray icon
    tray.OnClick(func() {
        if window.IsVisible() {
            window.Hide()
        } else {
            window.Show()
            window.Focus()
        }
    })

    // Create tray menu
    trayMenu := app.NewMenu()

    showItem := trayMenu.Add("Show Window")
    showItem.OnClick(func(ctx *application.Context) {
        window.Show()
        window.Focus()
    })

    trayMenu.AddSeparator()

    trayMenu.Add("Settings").OnClick(func(ctx *application.Context) {
        // Open settings window
    })

    trayMenu.Add("About").OnClick(func(ctx *application.Context) {
        // Show about dialog
    })

    trayMenu.AddSeparator()

    trayMenu.Add("Quit").OnClick(func(ctx *application.Context) {
        app.Quit()
    })

    tray.SetMenu(trayMenu)
}
```

### Dynamische Menüaktualisierungen

```go
type Editor struct {
    app         *application.App
    menu        *application.Menu
    undoItem    *application.MenuItem
    redoItem    *application.MenuItem
    saveItem    *application.MenuItem
    undoStack   []string
    redoStack   []string
    hasChanges  bool
}

func (e *Editor) createMenu() {
    e.menu = e.app.NewMenu()

    fileMenu := e.menu.AddSubmenu("File")
    e.saveItem = fileMenu.Add("Save").SetAccelerator("Ctrl+S")
    e.saveItem.OnClick(func(ctx *application.Context) {
        e.save()
    })

    editMenu := e.menu.AddSubmenu("Edit")
    e.undoItem = editMenu.Add("Undo").SetAccelerator("Ctrl+Z")
    e.undoItem.OnClick(func(ctx *application.Context) {
        e.undo()
    })

    e.redoItem = editMenu.Add("Redo").SetAccelerator("Ctrl+Y")
    e.redoItem.OnClick(func(ctx *application.Context) {
        e.redo()
    })

    e.updateMenuState()
    e.app.Menu.Set(e.menu)
}

func (e *Editor) updateMenuState() {
    // Update menu items based on current state
    e.saveItem.SetEnabled(e.hasChanges)
    e.undoItem.SetEnabled(len(e.undoStack) > 0)
    e.redoItem.SetEnabled(len(e.redoStack) > 0)
    e.menu.Update()
}

func (e *Editor) onChange() {
    e.hasChanges = true
    e.updateMenuState()
}

func (e *Editor) save() {
    // Save logic
    e.hasChanges = false
    e.updateMenuState()
}
```

## Bewährte Vorgehensweisen

### ✅ Empfohlen

- **Standard-Tastenkombinationen verwenden** – Halten Sie sich an die Plattformkonventionen (Ctrl+C zum Kopieren usw.)
- **Nach Änderungen Update() aufrufen** – Andernfalls werden die Änderungen nicht im Menü angezeigt
- **Zusammengehörige Einträge gruppieren** – Verwenden Sie Trennlinien, um Menüeinträge zu organisieren
- **Nicht verfügbare Aktionen deaktivieren** – Blenden Sie sie nicht aus, sondern deaktivieren Sie sie mit SetEnabled(false)
- **Eindeutige Beschriftungen verwenden** – Formulieren Sie kurz und aussagekräftig
- **Plattformkonventionen einhalten** – Beachten Sie die unterschiedlichen Menümuster von macOS und Windows/Linux

### ❌ Nicht empfohlen

- **Update() nicht vergessen** – Dies ist der häufigste Fehler
- **Nicht zu tief verschachteln** – Beschränken Sie Menüs auf höchstens 2-3 Ebenen
- **Keine mehrdeutigen Beschriftungen verwenden** – „Verarbeiten“ gegenüber „Dokument verarbeiten“
- **Nicht unnötig verkomplizieren** – Halten Sie Menüs einfach und fokussiert
- **Keine uneinheitlichen Metaphern verwenden** – Verwenden Sie eine konsistente Benennung und Struktur

## Plattformspezifische Hinweise

### macOS

- Ein Anwendungsmenü mit dem Namen der Anwendung wird automatisch hinzugefügt
- Verwenden Sie für Tastenkombinationen `Cmd` anstelle von `Ctrl`
- „Über“, „Einstellungen“ und „Beenden“ sind standardmäßig im Anwendungsmenü enthalten

### Windows/Linux

- Kein automatisches Anwendungsmenü
- `Ctrl` für Tastenkürzel verwenden
- „Beenden“ befindet sich üblicherweise im Menü „Datei“
