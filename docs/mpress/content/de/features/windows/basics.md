---
title: "Fenstergrundlagen"
description: "Anwendungsfenster in Wails erstellen und verwalten"
slug: "features/windows/basics"
sourcePath: "features/windows/basics.md"
---

## Fensterverwaltung

Wails bietet eine **einheitliche API zur Fensterverwaltung**, die auf allen Plattformen funktioniert. Erstellen Sie Fenster, steuern Sie deren Verhalten und verwalten Sie mehrere Fenster mit vollständiger Kontrolle über Erstellung, Darstellung, Verhalten und Lebenszyklus.

## Schnellstart

```go
package main

import "github.com/wailsapp/wails/v3/pkg/application"

func main() {
    app := application.New(application.Options{
        Name: "My App",
    })

    // Create a window
    window := app.Window.New()
    
    // Configure it
    window.SetTitle("Hello Wails")
    window.SetSize(800, 600)
    window.Center()
    
    // Show it
    window.Show()

    app.Run()
}
```

**Das ist alles!** Sie haben nun ein plattformübergreifendes Fenster.

## Fenster erstellen

### Einfaches Fenster

So erstellen Sie ein Fenster am einfachsten:

```go
window := app.Window.New()
```

**Sie erhalten:**

- Standardgröße (800x600)
- Standardtitel (Anwendungsname)
- Für Ihr Frontend bereites WebView
- Plattformnative Darstellung

### Fenster mit Optionen

Erstellen Sie ein Fenster mit einer benutzerdefinierten Konfiguration:

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "My Application",
    Width:  1200,
    Height: 800,
    X:      100,   // Position from left
    Y:      100,   // Position from top
    AlwaysOnTop: false,
    Frameless: false,
    Hidden: false,
    MinWidth: 400,
    MinHeight: 300,
    MaxWidth: 1920,
    MaxHeight: 1080,
})
```

**Gängige Optionen:**

| Option | Typ | Beschreibung |
| --- | --- | --- |
| `Title` | `string` | Fenstertitel |
| `Width` | `int` | Fensterbreite in Pixeln |
| `Height` | `int` | Fensterhöhe in Pixeln |
| `X` | `int` | X-Position (von links) |
| `Y` | `int` | Y-Position (von oben) |
| `AlwaysOnTop` | `bool` | Fenster über anderen Fenstern halten |
| `Frameless` | `bool` | Titelleiste und Rahmen entfernen |
| `Hidden` | `bool` | Ausgeblendet starten |
| `MinWidth` | `int` | Mindestbreite |
| `MinHeight` | `int` | Mindesthöhe |
| `MaxWidth` | `int` | Maximalbreite |
| `MaxHeight` | `int` | Maximalhöhe |

**Die vollständige Liste finden Sie unter [Fensteroptionen](/features/windows/options/).**

### Benannte Fenster

Vergeben Sie Namen für Fenster, damit Sie sie leicht abrufen können:

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:  "main-window",
    Title: "Main Application",
})

// Later, find it by name
if mainWindow, ok := app.Window.GetByName("main-window"); ok {
    mainWindow.Show()
}
```

**Anwendungsfälle:**

- Mehrere Fenster (Hauptfenster, Einstellungen, Info)
- Fenster aus verschiedenen Teilen Ihres Codes abrufen
- Kommunikation zwischen Fenstern

## Fenster steuern

### Ein- und Ausblenden

```go
// Show window
window.Show()

// Hide window
window.Hide()

// Check if visible
if window.IsVisible() {
    fmt.Println("Window is visible")
}
```

**Anwendungsfälle:**

- Startbildschirme (einblenden, dann ausblenden)
- Einstellungsfenster (ausblenden, wenn sie nicht benötigt werden)
- Popupfenster (bei Bedarf einblenden)

### Position und Größe

```go
// Set size
window.SetSize(1024, 768)

// Set position
window.SetPosition(100, 100)

// Centre on screen
window.Center()

// Get current size
width, height := window.Size()

// Get current position
x, y := window.Position()
```

**Koordinatensystem:**

- (0, 0) ist die obere linke Ecke des primären Bildschirms
- Positive X-Werte verlaufen nach rechts
- Positive Y-Werte verlaufen nach unten

### Fensterzustand

```go
// Minimise
window.Minimise()

// Maximise
window.Maximise()

// Fullscreen
window.Fullscreen()

// Restore to normal
window.Restore()

// Check state
if window.IsMinimised() {
    fmt.Println("Window is minimised")
}

if window.IsMaximised() {
    fmt.Println("Window is maximised")
}

if window.IsFullscreen() {
    fmt.Println("Window is fullscreen")
}
```

**Zustandsübergänge:**

```
Normal ←→ Minimised
Normal ←→ Maximised
Normal ←→ Fullscreen
```

### Titel und Erscheinungsbild

```go
// Set title
window.SetTitle("My Application - Document.txt")

// Set background colour — RGBA value (helper for RGB)
window.SetBackgroundColour(application.NewRGBA(0, 0, 0, 255))

// Set always on top
window.SetAlwaysOnTop(true)

// Set resizable
window.SetResizable(false)
```

### Fenster schließen

```go
// Close window — dispatches WindowClosing; a RegisterHook can call e.Cancel().
window.Close()
```

In v3 gibt es keine Methode `window.Destroy()`. Verwenden Sie `Close()` und registrieren Sie entweder mit `OnWindowEvent` einen Listener (Schließen kann nicht abgebrochen werden) oder mit `RegisterHook` einen Hook (kann `e.Cancel()` aufrufen, um das Fenster geöffnet zu lassen).

## Fenster suchen

### Nach Name

```go
if window, ok := app.Window.GetByName("settings"); ok {
    window.Show()
}
```

### Nach ID

Jedes Fenster hat eine eindeutige ID:

```go
id := window.ID()
fmt.Printf("Window ID: %d\n", id)

// Find by ID
if found, ok := app.Window.GetByID(id); ok {
    found.Focus()
}
```

### Aktuelles Fenster

Rufen Sie das aktuell fokussierte Fenster ab:

```go
current := app.Window.Current()
if current != nil {
    current.SetTitle("Active Window")
}
```

### Alle Fenster

Rufen Sie alle Fenster ab:

```go
windows := app.Window.GetAll()
fmt.Printf("Total windows: %d\n", len(windows))

for _, w := range windows {
    fmt.Printf("Window: %s (ID: %d)\n", w.Name(), w.ID())
}
```

## Fensterlebenszyklus

### Erstellung

```go
app.Window.OnCreate(func(window application.Window) {
    fmt.Printf("Window created: %s\n", window.Name())

    // Configure new windows
    window.SetMinSize(400, 300)
})
```

### Schließen

Um das Schließen eines Fensters zu verhindern, verwenden Sie `RegisterHook` mit dem Ereignis `WindowClosing`:

```go
window.RegisterHook(events.Common.WindowClosing, func(event *application.WindowEvent) {
    if hasUnsavedChanges() {
        // Ask user for confirmation
        result := showConfirmDialog("Unsaved changes. Close anyway?")
        if result != "yes" {
            // Cancel the close event
            event.Cancel()
        }
    }
})
```

**Wichtig:** `RegisterHook` fängt das Schließereignis ab, bevor es eintritt. Rufen Sie `event.Cancel()` auf, um das Schließen des Fensters zu verhindern. Dies funktioniert beim Schließen durch den Benutzer (Klick auf die X-Schaltfläche).

### Zerstörung

Um beim Schließen eines Fensters Aufräumarbeiten auszuführen, verwenden Sie `OnWindowEvent` mit dem Ereignis `WindowClosing`:

```go
window.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
    fmt.Println("Window is closing")
    // Cleanup resources
})
```

## Mehrere Fenster

### Mehrere Fenster erstellen

```go
// Main window
mainWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:   "main",
    Title:  "Main Application",
    Width:  1200,
    Height: 800,
})

// Settings window
settingsWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:   "settings",
    Title:  "Settings",
    Width:  600,
    Height: 400,
    Hidden: true,  // Start hidden
})

// Show settings when needed
settingsWindow.Show()
```

### Fensterkommunikation

Fenster können über Ereignisse kommunizieren:

```go
// In main window
app.Event.Emit("data-updated", map[string]interface{}{
    "value": 42,
})

// In settings window
app.Event.On("data-updated", func(event *application.CustomEvent) {
    data := event.Data.(map[string]interface{})
    value := data["value"].(int)
    fmt.Printf("Received: %d\n", value)
})
```

**Weitere Informationen finden Sie unter [Ereignisse](/features/events/system/).**

### Über- und untergeordnete Fenster

`WebviewWindowOptions` hat kein Feld `Parent`. Erstellen Sie das untergeordnete Fenster als normales Fenster und hängen Sie es als modales Sheet an ein übergeordnetes Fenster an:

```go
// Create child window
childWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Child Window",
})

// Attach to the parent — presents as a sheet on macOS.
mainWindow.AttachModal(childWindow)
```

**Verhalten:**

- Das untergeordnete Fenster bleibt über dem übergeordneten Fenster.
- Das untergeordnete Fenster ist modal und blockiert die Interaktion mit dem übergeordneten Fenster.

**Plattformunterstützung:**

- **macOS:** Vollständig unterstützt (Darstellung als Sheet).
- **Windows:** Nicht unterstützt.
- **Linux:** Nicht unterstützt.

## Plattformspezifische Funktionen

@tabs{sync-key="platform"}
[Windows]
**Windows-spezifische Funktionen:**

```go
// Flash taskbar button
window.Flash(true)  // Start flashing
window.Flash(false) // Stop flashing

// Trigger Windows 11 Snap Assist (Win+Z)
window.SnapAssist()
```

Es gibt kein fensterspezifisches `SetIcon`. Das Anwendungssymbol wird über `app.SetIcon([]byte)` für die App festgelegt (oder bei der Fenstererstellung über das Feld `application.LinuxWindow.Icon`, wenn ein Linux-spezifisches Fenstersymbol benötigt wird).

**Snap Assist:** Zeigt über den Pfad für Systemtastenkürzel die Optionen des Windows-11-Andocklayouts an. Verwenden Sie stattdessen [Native Nicht-Client-Bereiche unter Windows](/features/windows/frameless/#native-non-client-regions-on-windows), wenn eine benutzerdefinierte HTML-Maximierungsschaltfläche beim Darüberfahren native Andocklayouts anzeigen soll.

**Blinken der Taskleiste:** Nützlich für Benachrichtigungen, wenn das Fenster minimiert ist.

[macOS]
**macOS-spezifische Funktionen:**

```go
// Transparent title bar
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Mac: application.MacWindow{
        TitleBar: application.MacTitleBar{
            AppearsTransparent: true,
        },
        Backdrop: application.MacBackdropTranslucent,
    },
})
```

**Hintergrundtypen:**

- `MacBackdropNormal` – Standardfenster
- `MacBackdropTranslucent` – Durchscheinender Hintergrund; für die Transparenz der Webview ist eine **private API erforderlich**.
- `MacBackdropTransparent` – Vollständig transparent; für die Transparenz der Webview ist eine **private API erforderlich**.
- `MacBackdropLiquidGlass` – Glashintergrund; für die Transparenz der Webview ist eine **private API erforderlich**.

Erstellen Sie den Build mit `-tags private_mac_apis`, damit diese Effekte durch die Webview sichtbar sind. Ohne diese Option bleibt die Webview undurchsichtig. `TitleBar.AppearsTransparent` selbst verwendet öffentliche APIs. Siehe [Private macOS-APIs](/guides/build/private-macos-apis/).

**Sammlungsverhalten:** Steuern Sie das Verhalten von Fenstern über mehrere Spaces hinweg:

- `MacWindowCollectionBehaviorCanJoinAllSpaces` – Auf allen Spaces sichtbar
- `MacWindowCollectionBehaviorFullScreenAuxiliary` – Kann Vollbild-Apps überlagern

**Natives Vollbild:** Der macOS-Vollbildmodus erstellt einen neuen Space (virtuellen Desktop).

[Linux]
**Linux-spezifische Funktionen:**

```go
// Set window icon (per-window struct is LinuxWindow, not the app-level LinuxOptions)
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Linux: application.LinuxWindow{
        Icon: iconBytes,
    },
})
```

**Hinweise zu Desktop-Umgebungen:**

- GNOME: Vollständig unterstützt
- KDE Plasma: Vollständig unterstützt
- XFCE: Teilweise unterstützt
- Andere: Unterschiedlich

**Kachelnde Fenstermanager (Hyprland, Sway, i3 usw.):**

- `Minimise()` und `Maximise()` funktionieren möglicherweise nicht wie erwartet – der Fenstermanager steuert die Fenstergeometrie
- Anforderungen für `SetSize()` und `SetPosition()` sind unverbindlich und können ignoriert werden
- `Fullscreen()` funktioniert in der Regel wie erwartet
- Einige Fenstermanager unterstützen „Immer im Vordergrund“ nicht

@end

## Gängige Muster

### Startbildschirm

```go
// Create splash screen
splash := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:     "Loading...",
    Width:     400,
    Height:    300,
    Frameless: true,
    AlwaysOnTop: true,
})

// Show splash
splash.Show()

// Initialise application
time.Sleep(2 * time.Second)

// Hide splash, show main window
splash.Close()
mainWindow.Show()
```

### Einstellungsfenster

```go
var settingsWindow *application.WebviewWindow

func showSettings() {
    if settingsWindow == nil {
        settingsWindow = app.Window.NewWithOptions(application.WebviewWindowOptions{
            Name:   "settings",
            Title:  "Settings",
            Width:  600,
            Height: 400,
        })
    }
    
    settingsWindow.Show()
    settingsWindow.Focus()
}
```

### Schließen bestätigen

```go
window.RegisterHook(events.Common.WindowClosing, func(event *application.WindowEvent) {
    if hasUnsavedChanges() {
        // Show dialog
        result := showConfirmDialog("Unsaved changes. Close anyway?")
        if result != "yes" {
            // Cancel the close event
            event.Cancel()
        }
    }
})
```

## Bewährte Vorgehensweisen

### ✅ Empfohlen

- **Wichtige Fenster benennen** – So lassen sie sich später leichter finden
- **Mindestgröße festlegen** – Verhindert unbrauchbare Layouts
- **Fenster zentrieren** – Bessere Benutzerfreundlichkeit als eine zufällige Position
- **Schließereignisse behandeln** – Verhindert Datenverlust
- **Auf allen Plattformen testen** – Das Verhalten variiert
- **Geeignete Größen verwenden** – Unterschiedliche Bildschirmgrößen berücksichtigen

### ❌ Nicht empfohlen

- **Nicht zu viele Fenster erstellen** – Das verwirrt Benutzer
- **Nicht vergessen, Fenster zu schließen** – Sonst entstehen Speicherlecks
- **Positionen nicht fest codieren** – Bildschirmgrößen unterscheiden sich
- **Plattformunterschiede nicht ignorieren** – Gründlich testen
- **UI-Thread nicht blockieren** – Für lange Vorgänge Goroutinen verwenden

## Fehlerbehebung

### Fenster wird nicht angezeigt

**Mögliche Ursachen:**

1. Fenster wurde ausgeblendet erstellt
2. Fenster befindet sich außerhalb des sichtbaren Bildschirmbereichs
3. Fenster befindet sich hinter anderen Fenstern

**Lösung:**

```go
window.Show()
window.Center()
window.Focus()
```

### Fenster hat die falsche Größe

**Ursache:** DPI-Skalierung unter Windows/Linux

**Lösung:**

```go
// Wails handles DPI automatically
// Just use logical pixels
window.SetSize(800, 600)
```

### Fenster wird sofort geschlossen

**Ursache:** Die Anwendung wird beendet, wenn das letzte Fenster geschlossen wird

**Lösung:**

```go
app := application.New(application.Options{
    Mac: application.MacOptions{
        ApplicationShouldTerminateAfterLastWindowClosed: false,
    },
})
```

## Nächste Schritte

@cards{cols="2"}
⚙ Fensteroptionen
Vollständige Referenz für alle Fensteroptionen.

[Mehr erfahren →](/features/windows/options/)

---
▣ Mehrere Fenster
Muster für Anwendungen mit mehreren Fenstern.

[Mehr erfahren →](/features/windows/multiple/)

---
★ Rahmenlose Fenster
Benutzerdefinierte Fensterdekorationen erstellen.

[Mehr erfahren →](/features/windows/frameless/)

---
🚀 Fensterereignisse
Ereignisse im Fensterlebenszyklus behandeln.

[Mehr erfahren →](/features/windows/events/)

@end

---

**Fragen?** Stelle sie auf [Discord](https://discord.gg/JDdSxwjhGf) oder sieh dir die [Fensterbeispiele](https://github.com/wailsapp/wails/tree/master/v3/examples) an.
