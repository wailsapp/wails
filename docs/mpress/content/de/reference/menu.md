---
title: "Menü-API"
description: "Vollständige Referenz der Menü-API"
slug: "reference/menu"
sourcePath: "reference/menu.md"
---

## Überblick

Die Menü-API stellt Methoden zum Erstellen und Verwalten von Anwendungsmenüs, Kontextmenüs und Taskleistenmenüs bereit.

**Menütypen:**

- **Anwendungsmenüs** – obere Menüleiste (Datei, Bearbeiten usw.)
- **Kontextmenüs** – Rechtsklickmenüs
- **Taskleistenmenüs** – Menüs im Infobereich der Taskleiste

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

- `label` – der für den Menüeintrag angezeigte Text

**Rückgabewert:** der erstellte Menüeintrag

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

- `label` – die Beschriftung des Untermenüs

**Rückgabewert:** das erstellte Untermenü

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

**Empfehlung:** Gruppieren Sie zusammengehörige Menüeinträge mit Trennlinien.

### AddSectionHeader()

Fügt unter **macOS 14 und neuer** eine native, nicht interaktive Abschnittsüberschrift hinzu.
Unter älteren macOS-Versionen und auf allen anderen Plattformen (einschließlich Server-Builds) hat dies keine Wirkung: Es wird kein Platzhaltereintrag eingefügt.

```go
func (m *Menu) AddSectionHeader(label string)
```

Verwenden Sie Überschriften, um Gruppen in Taskleisten-, Kontext- oder Anwendungsmenüs zu beschriften:

```go
menu.AddSectionHeader("mac-node1")
menu.Add("Windows VM").OnClick(openWindowsVM)
menu.Add("Linux VM").OnClick(openLinuxVM)
menu.AddSectionHeader("ai-master")
menu.Add("Build VM").OnClick(openBuildVM)
```

Überschriften verwenden das standardmäßige Erscheinungsbild von AppKit für Abschnittsüberschriften. Rufen Sie nach dem erneuten Aufbau eines bereits angezeigten Menüs wie bei anderen strukturellen Änderungen `menu.Update()` auf.
Ein Beispiel für eine über `os.Hostname()` abgerufene Überschrift finden Sie im Beispiel `systray-menu`.

Ausgelassene Überschriften trennen keine Optionsgruppen. Verwenden Sie `AddSeparator()`, wenn Optionsgruppen auf jeder Plattform getrennt bleiben müssen.

### AddCheckbox()

Fügt einen auswählbaren Menüeintrag hinzu.

```go
func (m *Menu) AddCheckbox(label string, checked bool) *MenuItem
```

**Parameter:**

- `label` – die Beschriftung des Kontrollkästchens
- `checked` – anfänglicher Auswahlstatus

**Beispiel:**

```go
darkMode := menu.AddCheckbox("Dark Mode", false)
darkMode.OnClick(func(ctx *application.Context) {
    isChecked := darkMode.Checked()
    // Toggle dark mode
})
```

### AddRadio()

Fügt einen Optionsmenüeintrag hinzu (gegenseitig ausschließende Gruppe).

```go
func (m *Menu) AddRadio(label string, checked bool) *MenuItem
```

**Parameter:**

- `label` – die Beschriftung des Optionsfelds
- `checked` – anfänglicher Auswahlstatus

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

Aktualisiert das Menü, damit alle an Menüeinträgen vorgenommenen Änderungen angezeigt werden.

```go
func (m *Menu) Update()
```

**Beispiel:**

```go
item.SetEnabled(false)
menu.Update()  // Must call to apply changes
```

**Wichtig:** Rufen Sie nach dem Ändern der Eigenschaften von Menüeinträgen immer `Update()` auf.

## Methoden für Menüeinträge

### OnClick()

Registriert einen Klick-Handler für den Menüeintrag.

```go
func (mi *MenuItem) OnClick(callback func(ctx *application.Context)) *MenuItem
```

**Parameter:**

- `callback` – Funktion, die beim Anklicken des Eintrags aufgerufen wird

**Rückgabewert:** der Menüeintrag (zur Methodenverkettung)

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

Legt den Auswahlstatus von Kontrollkästchen- und Optionsmenüeinträgen fest.

```go
func (mi *MenuItem) SetChecked(checked bool) *MenuItem
```

**Beispiel:**

```go
darkModeItem.SetChecked(isDarkModeEnabled)
menu.Update()
```

### Checked()

Gibt den aktuellen Auswahlstatus zurück.

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

Legt einen Tooltip fest, der angezeigt wird, wenn der Mauszeiger über dem Menüeintrag verweilt.

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
- **macOS:** Fügt automatisch ein Anwendungsmenü mit dem Namen der App hinzu

## Kontextmenüs

### app.ContextMenu.New() / app.ContextMenu.Add()

Erstelle über den Manager ein `*ContextMenu` und registriere es unter einem Namen. `ContextMenuManager.Add` akzeptiert `*ContextMenu` – **nicht** `*Menu` – und es gibt **keine** Methode `app.RegisterContextMenu`.

```go
func (cm *ContextMenuManager) New() *ContextMenu
func (cm *ContextMenuManager) Add(name string, menu *ContextMenu)
func (cm *ContextMenuManager) Get(name string) (*ContextMenu, bool)
func (cm *ContextMenuManager) Remove(name string)
```

Alternativ erstellt und registriert die paketweite Funktion `application.NewContextMenu(name string) *ContextMenu` ein Kontextmenü in einem Schritt.

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

Die Runtime löst ein registriertes Kontextmenü aus, wenn für das Ziel des Rechtsklicks oder eines seiner übergeordneten Elemente die benutzerdefinierte CSS-Eigenschaft `--custom-contextmenu` auf diesen Namen gesetzt ist. Die optionale Eigenschaft `--custom-contextmenu-data` wird über `ctx.ContextMenuData()` an den Go-Callback übergeben. Um das standardmäßige Browser-Kontextmenü zu unterdrücken, setze `--default-contextmenu: hide` (oder `auto`/`show`).

```html
<!-- Trigger context menu on right-click -->
<div style="--custom-contextmenu: editor; --default-contextmenu: hide">
    Right-click here for context menu
</div>
```

Es gibt kein Attribut `data-wails-context-menu="..."` – es war nie an die Runtime angebunden.

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

## Infobereichsmenü

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

Legt den Tooltip fest, der beim Bewegen des Mauszeigers über das Symbol im Infobereich angezeigt wird. Gibt keinen Wert zurück.

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

### Anwendung mit Infobereichssymbol

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

- **Standard-Tastenkürzel verwenden** – Halte dich an die Konventionen der Plattform (Ctrl+C zum Kopieren usw.)
- **Nach Änderungen Update() aufrufen** – Andernfalls werden die Änderungen nicht im Menü angezeigt
- **Zusammengehörige Einträge gruppieren** – Verwende Trennlinien, um Menüeinträge zu gliedern
- **Nicht verfügbare Aktionen deaktivieren** – Blende sie nicht aus, sondern deaktiviere sie mit SetEnabled(false)
- **Eindeutige Beschriftungen verwenden** – Formuliere knapp und aussagekräftig
- **Plattformkonventionen befolgen** – Beachte die unterschiedlichen Menümuster von macOS und Windows/Linux

### ❌ Nicht empfohlen

- **Update() nicht vergessen** – Dies ist der häufigste Fehler
- **Menüs nicht zu tief verschachteln** – Beschränke Menüs auf maximal 2-3 Ebenen
- **Keine mehrdeutigen Beschriftungen verwenden** – „Verarbeiten“ statt „Dokument verarbeiten“
- **Nicht unnötig verkomplizieren** – Halte Menüs einfach und zweckgerichtet
- **Keine uneinheitlichen Metaphern verwenden** – Verwende eine einheitliche Benennung und Gliederung

## Plattformspezifische Hinweise

### macOS

- Ein Anwendungsmenü mit dem Namen der App wird automatisch hinzugefügt
- Verwende für Tastenkürzel `Cmd` anstelle von `Ctrl`
- „Über“, „Einstellungen“ und „Beenden“ befinden sich standardmäßig im Anwendungsmenü

### Windows/Linux

- Kein automatisches Anwendungsmenü
- Verwende `Ctrl` für Tastenkürzel
- „Beenden“ befindet sich üblicherweise im Menü „Datei“
