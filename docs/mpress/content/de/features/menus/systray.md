---
title: "System-Tray-Menüs"
description: "Integriere den Infobereich der Taskleiste in deine Anwendung"
slug: "features/menus/systray"
sourcePath: "features/menus/systray.md"
---

## System-Tray-Menüs

Wails bietet **einheitliche System-Tray-APIs**, die auf allen Plattformen funktionieren. Erstelle System-Tray-Symbole mit Menüs, verknüpfe Fenster und verarbeite Klicks mit nativem Plattformverhalten für Hintergrundanwendungen, Dienste und Dienstprogramme für den Schnellzugriff.

![Ein Wails-System-Tray-Menü, geöffnet über die macOS-Menüleiste](/assets/screenshots/systray-menu-macos.png)

Unter macOS erscheint ein Wails-System-Tray-Element in der Menüleiste und öffnet ein natives Menü. Das Beispiel enthält deaktivierte Einträge, Kontrollkästchen, Optionsfelder, Untermenüs und Aktionseinträge.

## Schnellstart

```go
package main

import (
    _ "embed"
    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets/icon.png
var icon []byte

func main() {
    app := application.New(application.Options{
        Name: "Tray App",
    })

    // Create system tray
    systray := app.SystemTray.New()
    systray.SetIcon(icon)
    systray.SetLabel("My App")

    // Add menu
    menu := app.NewMenu()
    menu.Add("Show").OnClick(func(ctx *application.Context) {
        // Show main window
    })
    menu.Add("Quit").OnClick(func(ctx *application.Context) {
        app.Quit()
    })
    systray.SetMenu(menu)

    // Create hidden window
    window := app.Window.New()
    window.Hide()

    app.Run()
}
```

**Ergebnis:** System-Tray-Symbol mit Menü auf allen Plattformen.

## System-Tray-Symbol erstellen

### Einfaches System-Tray-Symbol

```go
// Create system tray
systray := app.SystemTray.New()

// Set icon
systray.SetIcon(iconBytes)

// Set label (macOS) / tooltip (Windows)
systray.SetLabel("My Application")
```

### Mit Symbol

Symbole sollten eingebettet werden:

```go
import _ "embed"

//go:embed assets/icon.png
var icon []byte

//go:embed assets/icon-dark.png
var iconDark []byte

func main() {
    app := application.New(application.Options{
        Name: "My App",
    })

    systray := app.SystemTray.New()
    systray.SetIcon(icon)
    systray.SetDarkModeIcon(iconDark)  // Windows and macOS dark mode
    
    app.Run()
}
```

**Anforderungen an Symbole:**

| Plattform | Größe | Format | Hinweise |
| --- | --- | --- | --- |
| **Windows** | 16x16 oder 32x32 | PNG, ICO | Infobereich der Taskleiste |
| **macOS** | 18x18 bis 22x22 | PNG | Menüleiste, Vorlage empfohlen |
| **Linux** | 22x22 bis 48x48 | PNG, SVG | Je nach Desktop-Umgebung |

### Vorlagensymbole (macOS)

Vorlagensymbole passen sich automatisch an den hellen oder dunklen Modus an:

```go
systray.SetTemplateIcon(iconBytes)
```

**Richtlinien für Vorlagensymbole:**

- Verwende ausschließlich Schwarz und transparente Farben
- Schwarz wird im dunklen Modus zu Weiß
- Benenne die Datei mit dem Suffix `Template`: `iconTemplate.png`
- [Gestaltungsrichtlinien](https://bjango.com/articles/designingmenubarextras/)

## Menüs hinzufügen

System-Tray-Menüs funktionieren wie Anwendungsmenüs:

```go
menu := app.NewMenu()

// Add items
menu.Add("Open").OnClick(func(ctx *application.Context) {
    showMainWindow()
})

menu.AddSeparator()

menu.AddCheckbox("Start at Login", false).OnClick(func(ctx *application.Context) {
    enabled := ctx.ClickedMenuItem().Checked()
    setStartAtLogin(enabled)
})

menu.AddSeparator()

menu.Add("Quit").OnClick(func(ctx *application.Context) {
    app.Quit()
})

// Set menu
systray.SetMenu(menu)
```

**Alle Menüeintragstypen** findest du in der [Menüreferenz](/features/menus/reference/).

## Fenster verknüpfen

Verknüpfe ein Fenster mit dem System-Tray-Symbol, damit es automatisch ein- und ausgeblendet wird:

```go
// Create window
window := app.Window.New()

// Attach to tray
systray.AttachWindow(window)

// Configure behaviour — these are setters that return the receiver for chaining.
systray.WindowOffset(10)                          // Pixels from tray icon
systray.WindowDebounce(200 * time.Millisecond)    // Click debounce
```

**Verhalten:**

- Das Fenster ist anfangs ausgeblendet
- **Linksklick auf das System-Tray-Symbol** → Sichtbarkeit des Fensters umschalten
- **Rechtsklick auf das System-Tray-Symbol** → Menü anzeigen (sofern festgelegt)
- Das Fenster wird in der Nähe des System-Tray-Symbols positioniert

**Beispiel: Popup-Fenster**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:           "Quick Access",
    Width:           300,
    Height:          400,
    Frameless:       true, // No title bar
    AlwaysOnTop:     true, // Stay on top
    HideOnFocusLost: true, // Dismiss when another window receives focus
    HideOnEscape:    true, // Dismiss when the user presses Escape
})

systray.AttachWindow(window)
systray.WindowOffset(5)
```

`HideOnFocusLost` ist für System-Tray-Popups unter Windows und macOS sowie auf Linux-Desktops mit Fokus durch Klicken nützlich. Wails deaktiviert dieses Verhalten unter Linux in Umgebungen, in denen der Fokus der Maus folgt (darunter gängige Hyprland-, Sway- und i3- Konfigurationen). Andernfalls könnte das Popup beim Verlassen ausgeblendet werden, bevor es verwendet werden kann. `HideOnEscape` bleibt in diesen Umgebungen verfügbar.

Die oben beschriebenen Verhaltensweisen für Links- und Rechtsklicks sind sinnvolle Standardwerte. Ein expliziter `OnClick`- oder `OnRightClick`-Handler ersetzt den jeweiligen Standardwert. Plattformprüfungen und Sonderfälle findest du in der [manuellen Systray-Testsuite](https://github.com/wailsapp/wails/tree/master/v3/test/manual/systray) und im [Systray-Stresstestbeispiel](https://github.com/wailsapp/wails/tree/master/v3/examples/systray-stress).

## Klick-Handler

Verarbeite Klicks auf das System-Tray-Symbol:

```go
systray := app.SystemTray.New()

// Left click
systray.OnClick(func() {
    fmt.Println("Tray icon clicked")
})

// Right click
systray.OnRightClick(func() {
    fmt.Println("Tray icon right-clicked")
})

// Double click
systray.OnDoubleClick(func() {
    fmt.Println("Tray icon double-clicked")
})

// Mouse enter/leave
systray.OnMouseEnter(func() {
    fmt.Println("Mouse entered tray icon")
})

systray.OnMouseLeave(func() {
    fmt.Println("Mouse left tray icon")
})
```

**Plattformunterstützung:**

| Ereignis | Windows | macOS | Linux |
| --- | --- | --- | --- |
| OnClick | ✅ | ✅ | ✅ |
| OnRightClick | ✅ | ✅ | ✅ |
| OnDoubleClick | ✅ | ✅ | ⚠️ Unterschiedlich |
| OnMouseEnter | ✅ | ✅ | ⚠️ Unterschiedlich |
| OnMouseLeave | ✅ | ✅ | ⚠️ Unterschiedlich |

## Dynamische Aktualisierungen

Aktualisiere System-Tray-Symbol und Menü dynamisch:

### Symbol ändern

```go
var isActive bool

func updateTrayIcon() {
    if isActive {
        systray.SetIcon(activeIcon)
        systray.SetLabel("Active")
    } else {
        systray.SetIcon(inactiveIcon)
        systray.SetLabel("Inactive")
    }
}
```

### Menü aktualisieren

```go
var isPaused bool

pauseMenuItem := menu.Add("Pause")

pauseMenuItem.OnClick(func(ctx *application.Context) {
    isPaused = !isPaused
    
    if isPaused {
        pauseMenuItem.SetLabel("Resume")
    } else {
        pauseMenuItem.SetLabel("Pause")
    }
    
    menu.Update()  // Important!
})
```

@note{type="caution" title="Update() immer aufrufen"}
Rufe nach einer Änderung des Menüzustands **die Methode `menu.Update()`** auf. Weitere Informationen findest du in der [Menüreferenz](/features/menus/reference/#enabled-state).

@end

### Menü neu erstellen

Erstellen Sie bei umfangreichen Änderungen das gesamte Menü neu:

```go
func rebuildTrayMenu(status string) {
    menu := app.NewMenu()
    
    // Status-specific items
    switch status {
    case "syncing":
        menu.Add("Syncing...").SetEnabled(false)
        menu.Add("Pause Sync").OnClick(pauseSync)
    case "synced":
        menu.Add("Up to date ✓").SetEnabled(false)
        menu.Add("Sync Now").OnClick(startSync)
    case "error":
        menu.Add("Sync Error").SetEnabled(false)
        menu.Add("Retry").OnClick(retrySync)
    }
    
    menu.AddSeparator()
    menu.Add("Quit").OnClick(func(ctx *application.Context) {
        app.Quit()
    })
    
    systray.SetMenu(menu)
}
```

## Plattformspezifische Funktionen

@tabs{sync-key="platform"}
[macOS]
**Integration in die Menüleiste:**

```go
// Set label (appears next to icon)
systray.SetLabel("My App")

// Use template icon (adapts to dark mode)
systray.SetTemplateIcon(iconBytes)

// Set icon position — uses AppKit NSImage placement constants.
systray.SetIconPosition(application.NSImageRight)
```

**Symbolpositionen** (entsprechend `NSImagePosition`):

- `application.NSImageLeft` – Symbol links von der Beschriftung.
- `application.NSImageRight` – Symbol rechts von der Beschriftung.
- `application.NSImageOnly` – Nur Symbol, keine Beschriftung.
- `application.NSImageNone` – Nur Beschriftung, kein Symbol.

**Bewährte Vorgehensweisen:**

- Vorlagensymbole verwenden (schwarz und transparent)
- Beschriftungen kurz halten (3-5 Zeichen)
- 18x18 bis 22x22 Pixel für Retina-Displays
- Sowohl im hellen als auch im dunklen Modus testen

[Windows]
**Integration in den Infobereich:**

```go
// Set tooltip (appears on hover)
systray.SetTooltip("My Application")

// Or use SetLabel (same as tooltip on Windows)
systray.SetLabel("My Application")

// Show/Hide functionality (fully functional)
systray.Show()  // Show tray icon
systray.Hide()  // Hide tray icon
```

**Anforderungen an Symbole:**

- 16x16 oder 32x32 Pixel
- PNG- oder ICO-Format
- Transparenter Hintergrund

**Beschränkungen für Tooltips:**

- Maximal 127 UTF-16-Zeichen
- Längere Tooltips werden abgeschnitten
- Für eine optimale Benutzererfahrung kurz halten

**Plattformfunktionen:**

- System-Tray-Symbol bleibt nach Neustarts von Windows Explorer erhalten
- Die Methoden Show() und Hide() sind vollständig funktionsfähig
- Korrekte Lebenszyklusverwaltung

**Bewährte Vorgehensweisen:**

- Für Displays mit hoher Pixeldichte 32x32 verwenden
- Tooltips auf weniger als 127 Zeichen beschränken
- Auf verschiedenen Windows-Versionen testen
- Den Überlauf des Infobereichs berücksichtigen
- Show/Hide verwenden, um das Taskleistensymbol bedingt ein- oder auszublenden

[Linux]
**Integration in den Systembereich:**

Verwendet die StatusNotifierItem-Spezifikation (in den meisten modernen Desktop-Umgebungen).

```go
systray.SetIcon(iconBytes)
systray.SetLabel("My App")
```

**Unterstützung für Desktop-Umgebungen:**

- **GNOME**: Obere Leiste (mit Erweiterung)
- **KDE Plasma**: Systembereich
- **XFCE**: Infobereich
- **Andere**: Unterschiedlich

**Bewährte Vorgehensweisen:**

- 22x22 oder 24x24 Pixel verwenden
- SVG-Symbole lassen sich besser skalieren
- Auf den Ziel-Desktop-Umgebungen testen
- Eine Ausweichlösung für nicht unterstützte Desktop-Umgebungen bereitstellen

@end

## Vollständiges Beispiel

Hier sehen Sie eine produktionsreife Anwendung mit Taskleistensymbol:

```go
package main

import (
    _ "embed"
    "fmt"
    "time"
    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets/icon.png
var icon []byte

//go:embed assets/icon-active.png
var iconActive []byte

type TrayApp struct {
    app     *application.App
    systray *application.SystemTray
    window  *application.WebviewWindow
    menu    *application.Menu
    isActive bool
}

func main() {
    app := application.New(application.Options{
        Name: "Tray Application",
        Mac: application.MacOptions{
            ApplicationShouldTerminateAfterLastWindowClosed: false,
        },
    })

    trayApp := &TrayApp{app: app}
    trayApp.setup()

    app.Run()
}

func (t *TrayApp) setup() {
    // Create system tray
    t.systray = t.app.SystemTray.New()
    t.systray.SetIcon(icon)
    t.systray.SetLabel("Inactive")
    
    // Create menu
    t.createMenu()
    
    // Create window (hidden by default)
    t.window = t.app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:  "Tray Application",
        Width:  400,
        Height: 600,
        Hidden: true,
    })
    
    // Attach window to tray
    t.systray.AttachWindow(t.window)
    t.systray.WindowOffset(10)
    
    // Handle tray clicks
    t.systray.OnRightClick(func() {
        t.systray.OpenMenu()
    })
    
    // Start background task
    go t.backgroundTask()
}

func (t *TrayApp) createMenu() {
    t.menu = t.app.NewMenu()
    
    // Status item (disabled)
    statusItem := t.menu.Add("Status: Inactive")
    statusItem.SetEnabled(false)
    
    t.menu.AddSeparator()
    
    // Toggle active
    t.menu.Add("Start").OnClick(func(ctx *application.Context) {
        t.toggleActive()
    })
    
    // Show window
    t.menu.Add("Show Window").OnClick(func(ctx *application.Context) {
        t.window.Show()
        t.window.Focus()
    })
    
    t.menu.AddSeparator()
    
    // Settings
    t.menu.AddCheckbox("Start at Login", false).OnClick(func(ctx *application.Context) {
        enabled := ctx.ClickedMenuItem().Checked()
        t.setStartAtLogin(enabled)
    })
    
    t.menu.AddSeparator()
    
    // Quit
    t.menu.Add("Quit").OnClick(func(ctx *application.Context) {
        t.app.Quit()
    })
    
    t.systray.SetMenu(t.menu)
}

func (t *TrayApp) toggleActive() {
    t.isActive = !t.isActive
    t.updateTray()
}

func (t *TrayApp) updateTray() {
    if t.isActive {
        t.systray.SetIcon(iconActive)
        t.systray.SetLabel("Active")
    } else {
        t.systray.SetIcon(icon)
        t.systray.SetLabel("Inactive")
    }
    
    // Rebuild menu with new status
    t.createMenu()
}

func (t *TrayApp) backgroundTask() {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        if t.isActive {
            fmt.Println("Background task running...")
            // Do work
        }
    }
}

func (t *TrayApp) setStartAtLogin(enabled bool) {
    // Implementation varies by platform
    fmt.Printf("Start at login: %v\n", enabled)
}
```

## Sichtbarkeit steuern

Blenden Sie das Taskleistensymbol dynamisch ein oder aus:

```go
// Hide tray icon
systray.Hide()

// Show tray icon
systray.Show()
```

Es gibt keinen `IsVisible()`-Getter. Verwalten Sie die Sichtbarkeit bei Bedarf im eigenen Anwendungszustand.

**Plattformunterstützung:**

| Plattform | Hide() | Show() | Hinweise |
| --- | --- | --- | --- |
| **Windows** | ✅ | ✅ | Vollständig funktionsfähig – das Symbol wird im Infobereich ein- oder ausgeblendet |
| **macOS** | ✅ | ✅ | Menüleistenelement wird ein- oder ausgeblendet |
| **Linux** | ✅ | ✅ | Je nach Desktop-Umgebung unterschiedlich |

**Anwendungsfälle:**

- Taskleistensymbol entsprechend den Benutzereinstellungen vorübergehend ausblenden
- Headless-Modus, in dem das Taskleistensymbol nur bei Bedarf erscheint
- Sichtbarkeit abhängig vom Anwendungszustand umschalten

**Beispiel – bedingte Sichtbarkeit des Taskleistensymbols:**

```go
func (t *TrayApp) setTrayVisibility(visible bool) {
    if visible {
        t.systray.Show()
    } else {
        t.systray.Hide()
    }
}

// Show tray only when updates are available
func (t *TrayApp) checkForUpdates() {
    if hasUpdates {
        t.systray.Show()
        t.systray.SetLabel("Update Available")
    } else {
        t.systray.Hide()
    }
}
```

## Bereinigung

Zerstören Sie das Taskleistensymbol, wenn es nicht mehr benötigt wird:

```go
// In OnShutdown
app := application.New(application.Options{
    OnShutdown: func() {
        if systray != nil {
            systray.Destroy()
        }
    },
})
```

**Wichtig:** Zerstören Sie beim Herunterfahren stets das Taskleistensymbol, um Ressourcen freizugeben.

## Bewährte Vorgehensweisen

### ✅ Empfohlen

- **Vorlagensymbole unter macOS verwenden** – Passt sich dem Dunkelmodus an
- **Beschriftungen kurz halten** – maximal 3-5 Zeichen
- **Unter Windows Tooltips bereitstellen** – Hilft Benutzern, Ihre App zu erkennen
- **Auf allen Plattformen testen** – Das Verhalten variiert
- **Klicks angemessen behandeln** – Linksklick für die Hauptaktion, Rechtsklick für das Menü
- **Symbol entsprechend dem Status aktualisieren** – Visuelles Feedback ist wichtig
- **Beim Beenden zerstören** – Gibt Ressourcen frei

### ❌ Nicht empfohlen

- **Keine großen Symbole verwenden** – Plattformrichtlinien beachten
- **Keine langen Beschriftungen verwenden** – Sie werden abgeschnitten
- **Dunkelmodus nicht vergessen** – Im Dunkelmodus von Windows und macOS testen
- **Klick-Handler nicht blockieren** – Schnell ausführen lassen
- **menu.Update() nicht vergessen** – Nach einer Änderung des Menüzustands
- **Tray-Unterstützung nicht voraussetzen** – Einige Linux-Desktop-Umgebungen unterstützen sie nicht

## Fehlerbehebung

### Tray-Symbol wird nicht angezeigt

**Mögliche Ursachen:**

1. Symbolformat wird nicht unterstützt
2. Symbol ist zu groß oder zu klein
3. System-Tray wird nicht unterstützt (Linux)

**Lösung:**

Es gibt keine `SystemTraySupported()`-Hilfsfunktion. Erstellen Sie stattdessen das Tray, prüfen Sie die Plattform und stellen Sie bei fehlender Unterstützung eine sinnvolle Alternative bereit:

```go
// Probe support: on Linux without a notification-area extension, the tray
// will simply not appear. Defensive code can fall back to window-only mode
// based on runtime.GOOS or after a short timeout if no tray events arrive.
systray := app.SystemTray.New()
systray.SetIcon(iconBytes)
```

### Symbol wird unter macOS falsch dargestellt

**Ursache:** Es wird kein Vorlagensymbol verwendet

**Lösung:**

```go
// Use template icon
systray.SetTemplateIcon(iconBytes)

// Or design icon as template (black + transparent)
```

### Menü wird nicht aktualisiert

**Ursache:** Der Aufruf von `menu.Update()` wurde vergessen

**Lösung:**

```go
menuItem.SetLabel("New Label")
menu.Update()  // Add this!
```

## Nächste Schritte

@cards{cols="2"}
📖 Menüreferenz
Vollständige Referenz zu Menüeintragstypen und -eigenschaften.

[Mehr erfahren →](/features/menus/reference/)

---
☰ Anwendungsmenüs
Menüleisten für Anwendungen erstellen.

[Mehr erfahren →](/features/menus/application/)

---
◆ Kontextmenüs
Kontextmenüs für Rechtsklicks erstellen.

[Mehr erfahren →](/features/menus/context/)

---
📖 System-Tray-Beispiel
Eine vollständige System-Tray-Anwendung erkunden.

[Mehr erfahren →](https://github.com/wailsapp/wails/tree/master/v3/examples/systray-basic)

@end

---

**Fragen?** Fragen Sie auf [Discord](https://discord.gg/JDdSxwjhGf) nach oder sehen Sie sich die [System-Tray-Beispiele](https://github.com/wailsapp/wails/tree/master/v3/examples/systray-basic) an.
