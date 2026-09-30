---
title: "Anwendungsmenüs"
description: "Native Menüleisten für Ihre Desktop-Anwendung erstellen"
slug: "features/menus/application"
sourcePath: "features/menus/application.md"
---

## Das Problem

Professionelle Desktop-Anwendungen benötigen Menüleisten – Datei, Bearbeiten, Ansicht, Hilfe. Menüs funktionieren jedoch auf jeder Plattform anders:

- **macOS**: Globale Menüleiste am oberen Bildschirmrand
- **Windows**: Menüleiste in der Titelleiste des Fensters
- **Linux**: Abhängig von der Desktop-Umgebung

Plattformgerechte Menüs manuell zu erstellen, ist mühsam und fehleranfällig.

## Die Wails-Lösung

Wails bietet eine **einheitliche API**, die automatisch plattformspezifische native Menüs erstellt. Einmal schreiben und auf allen Plattformen natives Verhalten erhalten.

![Ein Wails-Anwendungsmenü unter macOS mit Standard-, Kontrollkästchen-, Options- und Untermenüeinträgen](/assets/screenshots/application-menu-macos.png)

Unter macOS befindet sich das Anwendungsmenü in der globalen Menüleiste. Diese Aufnahme zeigt das von der Wails-Menü-API gerenderte native Menü einschließlich deaktivierter Einträge sowie Kontrollkästchen-, Options- und Untermenüeinträgen.

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

**Das ist alles!** Sie verfügen nun über plattformspezifische native Menüs mit Standardeinträgen. Die Option `UseApplicationMenu` stellt sicher, dass Fenster unter Windows und Linux das Menü ohne zusätzlichen Code anzeigen.

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

**Empfohlene Vorgehensweise** – Verwenden Sie für plattformübergreifende Konsistenz `UseApplicationMenu`:

```go
// Set the application menu once
app.Menu.Set(menu)

// Create windows that inherit the menu on Windows/Linux
app.Window.NewWithOptions(application.WebviewWindowOptions{
    UseApplicationMenu: true,  // Window uses the app menu
})
```

Diese Vorgehensweise bewirkt Folgendes:

- Unter **macOS**: Das Menü erscheint am oberen Bildschirmrand (Standardverhalten)
- Unter **Windows/Linux**: Jedes Fenster mit `UseApplicationMenu: true` zeigt das Anwendungsmenü an

**Plattformspezifische Details:**

@tabs{sync-key="platform"}
[macOS]
**Globale Menüleiste** (eine pro Anwendung):

```go
app.Menu.Set(menu)
```

Das Menü erscheint am oberen Bildschirmrand und bleibt auch dann bestehen, wenn alle Fenster geschlossen sind. Die Option `UseApplicationMenu` hat unter macOS keine Wirkung, da alle Anwendungen das globale Menü verwenden.

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

Jedes Fenster kann ein eigenes Menü besitzen oder das Anwendungsmenü übernehmen. Das Menü erscheint in der Titelleiste des Fensters.

[Linux]
**Fensterspezifische Menüleiste** (normalerweise):

```go
// Option 1: Use application menu (recommended)
app.Menu.Set(menu)
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    UseApplicationMenu: true,
})

// Option 2: Set menu directly on window
window.SetMenu(menu)
```

Das Verhalten hängt von der Desktop-Umgebung ab. Einige Desktop-Umgebungen, etwa Unity, unterstützen globale Menüs.

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

## Menürollen

Wails bietet **vordefinierte Menürollen**, die automatisch plattformgerechte Menüstrukturen erstellen.

### Verfügbare Rollen

| Rolle | Beschreibung | Plattformhinweise |
| --- | --- | --- |
| `AppMenu` | Anwendungsmenü mit „Über“, Einstellungen und Beenden | **Nur macOS** |
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

- Über [App Name]
- Einstellungen ... (⌘,)
- ---
- Dienste
- ---
- [App Name] ausblenden (⌘H)
- Andere ausblenden (⌥⌘H)
- Alle einblenden
- ---
- [App Name] beenden (⌘Q)

**FileMenu**:

- Neu (⌘N)
- Öffnen ... (⌘O)
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
- Alle in den Vordergrund

**HelpMenu**:

- Hilfe zu [App Name]

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

- Info zu [App Name]

[Linux]
Ähnlich wie unter Windows, die Tastenkombinationen können jedoch je nach Desktop-Umgebung variieren.

@end

### Rollenmenüs anpassen

`Menu.AddRole(role)` gibt das als **Empfänger** verwendete Menü (das Menü der obersten Ebene) zurück, **nicht** das Untermenü der Rolle. Um dem Untermenü der Rolle Einträge hinzuzufügen, suche den eingefügten Rolleneintrag mit `FindByRole` und rufe darauf `GetSubmenu()` auf:

```go
menu.AddRole(application.FileMenu)

fileMenu := menu.FindByRole(application.FileMenu).GetSubmenu()
fileMenu.Add("Import...").OnClick(handleImport)
fileMenu.Add("Export...").OnClick(handleExport)
```

## Benutzerdefinierte Menüs

Erstelle eigene Menüs für anwendungsspezifische Funktionen:

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

**Weitere Menüeintragstypen** findest du in der [Menüreferenz](/features/menus/reference/).

## Dynamische Menüs

Aktualisiere Menüs entsprechend dem Anwendungszustand:

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
Rufe nach einer Änderung des Menüzustands (aktiviert/deaktiviert, Beschriftung, ausgewählt) **immer `menu.Update()`** auf. Dies ist besonders unter Windows entscheidend, da Menüs dort neu erstellt werden.

Weitere Informationen findest du in der [Menüreferenz](/features/menus/reference/#enabled-state).

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

### Menüs neu erstellen

Erstelle bei umfangreichen Änderungen das gesamte Menü neu:

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
- Verwende `menu.AddRole(application.AppMenu)` für Standardeinträge

**Standardpositionen:**

- **Info**: Anwendungsmenü
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
- **Info**: Hilfemenü

**Beispiel:**

```go
if runtime.GOOS == "windows" {
    menu.AddRole(application.FileMenu) // Exit is added automatically
    menu.AddRole(application.HelpMenu) // About is added automatically
}
```

### Linux

**Verhalten der Menüleiste:**

- Normalerweise fensterspezifisch (wie unter Windows)
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

- **Verwenden Sie Menürollen** für Standardmenüs (Datei, Bearbeiten usw.)
- **Befolgen Sie die plattformspezifischen Konventionen** für die Menüstruktur
- **Fügen Sie häufig verwendeten Aktionen Tastenkombinationen hinzu**
- **Rufen Sie menu.Update() auf**, nachdem Sie den Menüzustand geändert haben
- **Testen Sie auf allen Plattformen** – das Verhalten variiert
- **Halten Sie Menüs flach** – maximal 2-3 Ebenen
- **Verwenden Sie eindeutige Bezeichnungen** – „Projekt speichern“ statt „Speichern“

### ❌ Nicht empfohlen

- **Codieren Sie plattformspezifische Tastenkombinationen nicht fest** – verwenden Sie `CmdOrCtrl`
- **Platzieren Sie „Beenden“ unter macOS nicht im Menü „Datei“** – der Eintrag gehört in das Anwendungsmenü
- **Platzieren Sie „Über“ unter macOS nicht im Menü „Hilfe“** – der Eintrag gehört in das Anwendungsmenü
- **Vergessen Sie menu.Update() nicht** – sonst funktionieren Menüs nicht ordnungsgemäß
- **Verschachteln Sie Menüs nicht zu tief** – sonst verlieren Benutzer den Überblick
- **Verwenden Sie keinen Fachjargon** – halten Sie die Bezeichnungen benutzerfreundlich

## Nächste Schritte

@cards{cols="2"}
📖 Menüreferenz
Vollständige Referenz zu Menüeintragstypen und -eigenschaften.

[Mehr erfahren →](/features/menus/reference/)

---
◆ Kontextmenüs
Erstellen Sie Kontextmenüs, die per Rechtsklick geöffnet werden.

[Mehr erfahren →](/features/menus/context/)

---
★ System-Tray-Menüs
Fügen Sie eine Integration in den Infobereich beziehungsweise die Menüleiste hinzu.

[Mehr erfahren →](/features/menus/systray/)

---
📖 Menümuster
Gängige Menümuster und bewährte Vorgehensweisen.

[Mehr erfahren →](/guides/menus/)

@end

---

**Fragen?** Fragen Sie auf [Discord](https://discord.gg/JDdSxwjhGf) nach oder sehen Sie sich das [Menübeispiel](https://github.com/wailsapp/wails/tree/master/v3/examples/menu) an.
