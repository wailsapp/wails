---
title: "Anwendungsmenüs"
description: "Native Menüleisten für Ihre Desktopanwendung erstellen"
slug: "features/menus/application"
sourcePath: "features/menus/application.md"
---

## Das Problem

Professionelle Desktopanwendungen benötigen Menüleisten – Datei, Bearbeiten, Ansicht, Hilfe. Menüs funktionieren jedoch auf jeder Plattform anders:

- **macOS**: Globale Menüleiste am oberen Bildschirmrand
- **Windows**: Menüleiste in der Titelleiste des Fensters
- **Linux**: Abhängig von der Desktopumgebung

Plattformgerechte Menüs manuell zu erstellen, ist mühsam und fehleranfällig.

## Die Wails-Lösung

Wails stellt eine **einheitliche API** bereit, die automatisch plattformnative Menüs erstellt. Einmal schreiben und auf allen Plattformen natives Verhalten erhalten.

![Ein Wails-Anwendungsmenü unter macOS mit Standard-, Kontrollkästchen-, Options- und Untermenüeinträgen](/assets/screenshots/application-menu-macos.png)

Unter macOS wird das Anwendungsmenü in der globalen Menüleiste angezeigt. Diese Aufnahme zeigt das von der Wails-Menü-API gerenderte native Menü mit deaktivierten Einträgen sowie Kontrollkästchen-, Options- und Untermenüeinträgen.

## Schnellstart

```go
package main

import (
    "runtime"
    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "My App",
    })

    // Create menu
    menu := app.NewMenu()

    // Add standard menus (platform-appropriate)
    if runtime.GOOS == "darwin" {
        menu.AddRole(application.AppMenu)  // macOS only
    }
    menu.AddRole(application.FileMenu)
    menu.AddRole(application.EditMenu)
    menu.AddRole(application.WindowMenu)
    menu.AddRole(application.HelpMenu)

    // Set the application menu
    app.Menu.Set(menu)

    // Create window with UseApplicationMenu to inherit the menu on Windows/Linux
    app.Window.NewWithOptions(application.WebviewWindowOptions{
        UseApplicationMenu: true,
    })

    app.Run()
}
```

**Das ist alles!** Sie verfügen nun über plattformnative Menüs mit Standardeinträgen. Die Option `UseApplicationMenu` sorgt dafür, dass Windows- und Linux-Fenster das Menü ohne zusätzlichen Code anzeigen.

## Menüs erstellen

### Grundlegende Menüerstellung

```go
// Create a new menu
menu := app.NewMenu()

// Add a top-level menu
fileMenu := menu.AddSubmenu("File")

// Add menu items
fileMenu.Add("New").OnClick(func(ctx *application.Context) {
    // Handle New
})

fileMenu.Add("Open").OnClick(func(ctx *application.Context) {
    // Handle Open
})

fileMenu.AddSeparator()

fileMenu.Add("Quit").OnClick(func(ctx *application.Context) {
    app.Quit()
})
```

### Menü festlegen

**Empfohlener Ansatz** – Verwenden Sie `UseApplicationMenu` für plattformübergreifende Konsistenz:

```go
// Set the application menu once
app.Menu.Set(menu)

// Create windows that inherit the menu on Windows/Linux
app.Window.NewWithOptions(application.WebviewWindowOptions{
    UseApplicationMenu: true,  // Window uses the app menu
})
```

Dieser Ansatz bewirkt Folgendes:

- Unter **macOS**: Das Menü erscheint am oberen Bildschirmrand (Standardverhalten)
- Unter **Windows/Linux**: Jedes Fenster mit `UseApplicationMenu: true` zeigt das Anwendungsmenü an

**Plattformspezifische Details:**

@tabs{sync-key="platform"}
[macOS]
**Globale Menüleiste** (eine pro Anwendung):

```go
app.Menu.Set(menu)
```

Das Menü erscheint am oberen Bildschirmrand und bleibt auch dann bestehen, wenn alle Fenster geschlossen sind. Die Option `UseApplicationMenu` hat unter macOS keine Auswirkung, da alle Anwendungen das globale Menü verwenden.

[Windows]
**Fensterspezifische Menüleiste**:

```go
// Option 1: Use application menu (recommended)
app.Menu.Set(menu)
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    UseApplicationMenu: true,
})

// Option 2: Set menu directly on window
window.SetMenu(menu)
```

Jedes Fenster kann ein eigenes Menü haben oder das Anwendungsmenü übernehmen. Das Menü erscheint in der Titelleiste des Fensters.

[Linux]
**Fensterspezifische Menüleiste** (üblicherweise):

```go
// Option 1: Use application menu (recommended)
app.Menu.Set(menu)
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    UseApplicationMenu: true,
})

// Option 2: Set menu directly on window
window.SetMenu(menu)
```

Das Verhalten hängt von der Desktopumgebung ab. Einige Umgebungen wie Unity unterstützen globale Menüs.

@end

@note{type="tip" title="Plattformübergreifende Menüs vereinfachen"}
Durch die Verwendung von `UseApplicationMenu: true` entfällt plattformspezifischer Code wie dieser:

```go
// Old approach - no longer needed
if runtime.GOOS == "darwin" {
    app.Menu.Set(menu)
} else {
    window.SetMenu(menu)
}
```

@end

**Benutzerdefinierte fensterspezifische Menüs:**

Wenn ein Fenster ein anderes Menü als das Anwendungsmenü benötigt, legen Sie es direkt fest:

```go
window.SetMenu(customMenu)  // Overrides UseApplicationMenu
```

**Menüleiste eines Fensters entfernen:**

Übergeben Sie `nil`, um die Menüleiste aus einem Fenster zu entfernen:

```go
window.SetMenu(nil)
```

Unter Windows und Linux zeigt das Fenster anschließend keine Menüleiste mehr an – auch nicht das Anwendungsmenü, das es mit `UseApplicationMenu` übernommen hatte. Außerdem funktionieren die Tastenkürzel des entfernten Menüs in diesem Fenster nicht mehr. Erfolgt der Aufruf, bevor das Fenster gestartet wurde, wird das Fenster ohne Menüleiste erstellt. Ein späterer Aufruf von `SetMenu(menu)` fügt wieder eine Menüleiste hinzu.

Unter macOS gehört die Menüleiste zur Anwendung und nicht zu einem Fenster. Daher hat `SetMenu` dort keine Wirkung. Verwenden Sie `app.Menu.Set`, um das Anwendungsmenü zu ändern.

## Menürollen

Wails stellt **vordefinierte Menürollen** bereit, die automatisch plattformgerechte Menüstrukturen erstellen.

### Verfügbare Rollen

| Rolle | Beschreibung | Plattformhinweise |
| --- | --- | --- |
| `AppMenu` | Anwendungsmenü mit Info, Einstellungen und Beenden | **Nur macOS** |
| `FileMenu` | Dateioperationen (Neu, Öffnen, Speichern usw.) | Alle Plattformen |
| `EditMenu` | Textbearbeitung (Rückgängig, Wiederholen, Ausschneiden, Kopieren, Einfügen) | Alle Plattformen |
| `WindowMenu` | Fensterverwaltung (Minimieren, Zoomen usw.) | Alle Plattformen |
| `HelpMenu` | Hilfe und Informationen | Alle Plattformen |

### Rollen verwenden

```go
menu := app.NewMenu()

// macOS: Add application menu
if runtime.GOOS == "darwin" {
    menu.AddRole(application.AppMenu)
}

// All platforms: Add standard menus
menu.AddRole(application.FileMenu)
menu.AddRole(application.EditMenu)
menu.AddRole(application.WindowMenu)
menu.AddRole(application.HelpMenu)
```

**Das erhalten Sie:**

@tabs{sync-key="platform"}
[macOS]
**AppMenu** (mit Anwendungsname):

- Über [Anwendungsname]
- Einstellungen... (⌘,)
- ---
- Dienste
- ---
- [Anwendungsname] ausblenden (⌘H)
- Andere ausblenden (⌥⌘H)
- Alle einblenden
- ---
- [Anwendungsname] beenden (⌘Q)

**FileMenu**:

- Neu (⌘N)
- Öffnen... (⌘O)
- ---
- Fenster schließen (⌘W)

**EditMenu**:

- Rückgängig (⌘Z)
- Wiederholen (⇧⌘Z)
- ---
- Ausschneiden (⌘X)
- Kopieren (⌘C)
- Einfügen (⌘V)
- Alles auswählen (⌘A)

**WindowMenu**:

- Minimieren (⌘M)
- Zoomen
- ---
- Alle nach vorne bringen

**HelpMenu**:

- Hilfe für [Anwendungsname]

[Windows]
**FileMenu**:

- Neu (Ctrl+N)
- Öffnen... (Ctrl+O)
- ---
- Beenden (Alt+F4)

**EditMenu**:

- Rückgängig (Ctrl+Z)
- Wiederholen (Ctrl+Y)
- ---
- Ausschneiden (Ctrl+X)
- Kopieren (Ctrl+C)
- Einfügen (Ctrl+V)
- Alles auswählen (Ctrl+A)

**WindowMenu**:

- Minimieren
- Maximieren

**HelpMenu**:

- Über [App-Name]

[Linux]
Ähnlich wie unter Windows, die Tastenkürzel können jedoch je nach Desktop-Umgebung variieren.

@end

### Rollenmenüs anpassen

`Menu.AddRole(role)` gibt das **Empfänger**-Menü (das Menü der obersten Ebene) zurück, **nicht** das Untermenü der Rolle. Um dem Untermenü der Rolle Einträge hinzuzufügen, suchen Sie mit `FindByRole` nach dem eingefügten Rolleneintrag und rufen Sie darauf `GetSubmenu()` auf:

```go
menu.AddRole(application.FileMenu)

fileMenu := menu.FindByRole(application.FileMenu).GetSubmenu()
fileMenu.Add("Import...").OnClick(handleImport)
fileMenu.Add("Export...").OnClick(handleExport)
```

## Benutzerdefinierte Menüs

Erstellen Sie eigene Menüs für anwendungsspezifische Funktionen:

```go
// Add a custom top-level menu
toolsMenu := menu.AddSubmenu("Tools")

// Add items
toolsMenu.Add("Settings").OnClick(func(ctx *application.Context) {
    showSettingsWindow()
})

toolsMenu.AddSeparator()

// Add checkbox
toolsMenu.AddCheckbox("Dark Mode", false).OnClick(func(ctx *application.Context) {
    isDark := ctx.ClickedMenuItem().Checked()
    setTheme(isDark)
})

// Add radio group
toolsMenu.AddRadio("Small", true).OnClick(handleFontSize)
toolsMenu.AddRadio("Medium", false).OnClick(handleFontSize)
toolsMenu.AddRadio("Large", false).OnClick(handleFontSize)

// Add submenu
advancedMenu := toolsMenu.AddSubmenu("Advanced")
advancedMenu.Add("Configure...").OnClick(showAdvancedSettings)
```

Weitere Menüeintragstypen finden Sie unter **Weitere Menüeintragstypen** in der [Menüreferenz](/features/menus/reference/).

## Dynamische Menüs

Aktualisieren Sie Menüs abhängig vom Anwendungszustand:

### Einträge aktivieren/deaktivieren

```go
var saveMenuItem *application.MenuItem

func createMenu() {
    menu := app.NewMenu()
    fileMenu := menu.AddSubmenu("File")
    
    saveMenuItem = fileMenu.Add("Save")
    saveMenuItem.SetEnabled(false)  // Initially disabled
    saveMenuItem.OnClick(handleSave)
    
    app.Menu.Set(menu)
}

func onDocumentChanged() {
    saveMenuItem.SetEnabled(hasUnsavedChanges())
    menu.Update()  // Important!
}
```

@note{type="caution" title="Immer menu.Update() aufrufen"}
Rufen Sie nach jeder Änderung des Menüzustands (aktiviert/deaktiviert, Beschriftung, ausgewählt) **immer `menu.Update()`** auf. Dies ist besonders unter Windows entscheidend, da Menüs dort neu aufgebaut werden.

Weitere Informationen finden Sie in der [Menüreferenz](/features/menus/reference/#enabled-state).

@end

### Beschriftungen ändern

```go
updateMenuItem := menu.Add("Check for Updates")

updateMenuItem.OnClick(func(ctx *application.Context) {
    updateMenuItem.SetLabel("Checking...")
    menu.Update()
    
    checkForUpdates()
    
    updateMenuItem.SetLabel("Check for Updates")
    menu.Update()
})
```

### Menüs neu aufbauen

Bauen Sie bei umfangreichen Änderungen das gesamte Menü neu auf:

```go
func rebuildFileMenu() {
    menu := app.NewMenu()
    fileMenu := menu.AddSubmenu("File")
    
    fileMenu.Add("New").OnClick(handleNew)
    fileMenu.Add("Open").OnClick(handleOpen)
    
    // Add recent files dynamically
    if hasRecentFiles() {
        recentMenu := fileMenu.AddSubmenu("Open Recent")
        for _, file := range getRecentFiles() {
            filePath := file  // Capture for closure
            recentMenu.Add(filepath.Base(file)).OnClick(func(ctx *application.Context) {
                openFile(filePath)
            })
        }
        recentMenu.AddSeparator()
        recentMenu.Add("Clear Recent").OnClick(clearRecentFiles)
    }
    
    fileMenu.AddSeparator()
    fileMenu.Add("Quit").OnClick(func(ctx *application.Context) {
        app.Quit()
    })
    
    app.Menu.Set(menu)
}
```

## Fenster über Menüs steuern

Menüeinträge können Fenster steuern:

```go
viewMenu := menu.AddSubmenu("View")

// Toggle fullscreen
viewMenu.Add("Toggle Fullscreen").OnClick(func(ctx *application.Context) {
    if window, ok := app.Window.GetByName("main"); ok {
        window.ToggleFullscreen()
    }
})

// Zoom controls
viewMenu.Add("Zoom In").SetAccelerator("CmdOrCtrl++").OnClick(func(ctx *application.Context) {
    // Increase zoom
})

viewMenu.Add("Zoom Out").SetAccelerator("CmdOrCtrl+-").OnClick(func(ctx *application.Context) {
    // Decrease zoom
})

viewMenu.Add("Reset Zoom").SetAccelerator("CmdOrCtrl+0").OnClick(func(ctx *application.Context) {
    // Reset zoom
})
```

**Aktives Fenster abrufen:**

```go
menuItem.OnClick(func(ctx *application.Context) {
    window := application.Get().Window.Current() // the window the menu was invoked from
    // Use window
})
```

## Plattformspezifische Aspekte

### macOS

**Verhalten der Menüleiste:**

- Wird **am oberen Bildschirmrand** angezeigt (global)
- Bleibt bestehen, wenn alle Fenster geschlossen sind
- Das erste Menü ist **immer das Anwendungsmenü**
- Verwenden Sie `menu.AddRole(application.AppMenu)` für Standardeinträge

**Standardpositionen:**

- **Über**: Anwendungsmenü
- **Einstellungen**: Anwendungsmenü (⌘,)
- **Beenden**: Anwendungsmenü (⌘Q)
- **Hilfe**: Hilfemenü

**Beispiel:**

```go
if runtime.GOOS == "darwin" {
    menu.AddRole(application.AppMenu)  // Adds About, Preferences, Quit
    
    // Don't add Quit to File menu on macOS
    // Don't add About to Help menu on macOS
}
```

### Windows

**Verhalten der Menüleiste:**

- Wird in der **Fenstertitelleiste** angezeigt
- Jedes Fenster hat ein eigenes Menü
- Kein Anwendungsmenü

**Standardpositionen:**

- **Beenden**: Dateimenü (Alt+F4)
- **Einstellungen**: Menü „Extras“ oder „Bearbeiten“
- **Über**: Hilfemenü

**Beispiel:**

```go
if runtime.GOOS == "windows" {
    menu.AddRole(application.FileMenu) // Exit is added automatically
    menu.AddRole(application.HelpMenu) // About is added automatically
}
```

### Linux

**Verhalten der Menüleiste:**

- Normalerweise fensterbezogen (wie unter Windows)
- Einige Desktop-Umgebungen unterstützen globale Menüs (Unity, GNOME mit Erweiterung)
- Das Erscheinungsbild variiert je nach Desktop-Umgebung

**Bewährte Vorgehensweise:** Befolgen Sie die Windows-Konventionen und testen Sie auf den Ziel-Desktop-Umgebungen.

## Vollständiges Beispiel

Hier sehen Sie eine produktionsreife Menüstruktur:

```go
package main

import (
    "runtime"
    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "My Application",
    })

    // Create and set menu
    createMenu(app)

    // Create main window with UseApplicationMenu for cross-platform menu support
    app.Window.NewWithOptions(application.WebviewWindowOptions{
        UseApplicationMenu: true,
    })

    app.Run()
}

func createMenu(app *application.App) {
    menu := app.NewMenu()

    // Platform-specific application menu (macOS only)
    if runtime.GOOS == "darwin" {
        menu.AddRole(application.AppMenu)
    }

    // File menu — AddRole returns the receiver menu, not the role submenu.
    // To add items into the File submenu, look it up via FindByRole + GetSubmenu.
    menu.AddRole(application.FileMenu)
    fileMenu := menu.FindByRole(application.FileMenu).GetSubmenu()
    fileMenu.Add("Import...").SetAccelerator("CmdOrCtrl+I").OnClick(handleImport)
    fileMenu.Add("Export...").SetAccelerator("CmdOrCtrl+E").OnClick(handleExport)

    // Edit menu
    menu.AddRole(application.EditMenu)

    // View menu
    viewMenu := menu.AddSubmenu("View")
    viewMenu.Add("Toggle Fullscreen").SetAccelerator("F11").OnClick(toggleFullscreen)
    viewMenu.AddSeparator()
    viewMenu.AddCheckbox("Show Sidebar", true).OnClick(toggleSidebar)
    viewMenu.AddCheckbox("Show Toolbar", true).OnClick(toggleToolbar)

    // Tools menu
    toolsMenu := menu.AddSubmenu("Tools")
    
    // Settings location varies by platform
    if runtime.GOOS == "darwin" {
        // On macOS, Preferences is in Application menu (added by AppMenu role)
    } else {
        toolsMenu.Add("Settings").SetAccelerator("CmdOrCtrl+,").OnClick(showSettings)
    }
    
    toolsMenu.AddSeparator()
    toolsMenu.AddCheckbox("Dark Mode", false).OnClick(toggleDarkMode)

    // Window menu
    menu.AddRole(application.WindowMenu)

    // Help menu
    helpMenu := menu.AddRole(application.HelpMenu)
    helpMenu.Add("Documentation").OnClick(openDocumentation)
    
    // About location varies by platform
    if runtime.GOOS == "darwin" {
        // On macOS, About is in Application menu (added by AppMenu role)
    } else {
        helpMenu.AddSeparator()
        helpMenu.Add("About").OnClick(showAbout)
    }

    // Set the application menu
    app.Menu.Set(menu)
}

func handleImport(ctx *application.Context) {
    // Implementation
}

func handleExport(ctx *application.Context) {
    // Implementation
}

func toggleFullscreen(ctx *application.Context) {
    window := application.Get().Window.Current()
    window.ToggleFullscreen()
}

func toggleSidebar(ctx *application.Context) {
    // Implementation
}

func toggleToolbar(ctx *application.Context) {
    // Implementation
}

func showSettings(ctx *application.Context) {
    // Implementation
}

func toggleDarkMode(ctx *application.Context) {
    isDark := ctx.ClickedMenuItem().Checked()
    // Apply theme
}

func openDocumentation(ctx *application.Context) {
    // Open browser
}

func showAbout(ctx *application.Context) {
    // Show about dialog
}
```

## Bewährte Vorgehensweisen

### ✅ Empfohlen

- Verwenden Sie **Menürollen** für Standardmenüs (Datei, Bearbeiten usw.)
- Befolgen Sie bei der Menüstruktur die **Plattformkonventionen**
- Fügen Sie häufig verwendeten Aktionen **Tastenkürzel** hinzu
- Rufen Sie nach Änderungen des Menüzustands **menu.Update() auf**
- **Testen Sie auf allen Plattformen** – das Verhalten variiert
- **Halten Sie Menüs flach** – maximal 2-3 Ebenen
- **Verwenden Sie eindeutige Beschriftungen** – „Projekt speichern“ statt „Speichern“

### ❌ Nicht empfohlen

- **Codieren Sie plattformspezifische Tastenkürzel nicht fest** – verwenden Sie `CmdOrCtrl`
- **Platzieren Sie „Beenden“ unter macOS nicht im Dateimenü** – der Eintrag gehört ins Anwendungsmenü
- **Platzieren Sie „Über“ unter macOS nicht im Hilfemenü** – der Eintrag gehört ins Anwendungsmenü
- **Vergessen Sie menu.Update() nicht** – andernfalls funktionieren Menüs nicht ordnungsgemäß
- **Verschachteln Sie Menüs nicht zu tief** – Benutzer verlieren sonst die Orientierung
- **Verwenden Sie keinen Fachjargon** – halten Sie Beschriftungen benutzerfreundlich

## Nächste Schritte

@cards{cols="2"}
📖 Menüreferenz
Vollständige Referenz der Menüeintragstypen und Eigenschaften.

[Mehr erfahren →](/features/menus/reference/)

---
◆ Kontextmenüs
Erstellen Sie Kontextmenüs für Rechtsklicks.

[Mehr erfahren →](/features/menus/context/)

---
★ Infobereichsmenüs
Fügen Sie eine Integration für den Infobereich beziehungsweise die Menüleiste hinzu.

[Mehr erfahren →](/features/menus/systray/)

---
📖 Menümuster
Gängige Menümuster und bewährte Vorgehensweisen.

[Mehr erfahren →](/guides/menus/)

@end

---

**Fragen?** Stellen Sie sie auf [Discord](https://discord.gg/JDdSxwjhGf) oder sehen Sie sich das [Menübeispiel](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/menu) an.
