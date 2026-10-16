---
title: "Fenster-API"
description: "Vollständige Referenz für die Fenster-API"
slug: "reference/window"
sourcePath: "reference/window.md"
---

## Überblick

Die Fenster-API stellt Methoden bereit, mit denen Sie Darstellung, Verhalten und Lebenszyklus von Fenstern steuern können. Der Zugriff erfolgt über Fensterinstanzen oder den Manager `app.Window`.

`Window` ist eine von `*application.WebviewWindow` implementierte Schnittstelle; die folgenden Methodensignaturen gehören zu `*WebviewWindow`. Viele verändernde Methoden geben für die Verkettung `Window` zurück. Die Rückgabewerte sind bei der jeweiligen Methode dokumentiert.

**Häufige Vorgänge:**

- Fenster erstellen und anzeigen
- Größe, Position und Zustand steuern
- Fensterereignisse verarbeiten
- Fensterinhalt verwalten
- Darstellung und Verhalten konfigurieren

## Sichtbarkeit

### Show()

Zeigt das Fenster an. War das Fenster ausgeblendet, wird es wieder sichtbar. Gibt den Empfänger zur Verkettung zurück.

```go
func (w *WebviewWindow) Show() Window
```

**Beispiel:**

```go
window := app.Window.New()
window.Show()
```

### Hide()

Blendet das Fenster aus, ohne es zu schließen. Das Fenster verbleibt im Speicher und kann erneut angezeigt werden. Gibt den Empfänger zur Verkettung zurück.

```go
func (w *WebviewWindow) Hide() Window
```

**Beispiel:**

```go
// Hide window temporarily
window.Hide()

// Show it again later
window.Show()
```

**Anwendungsfälle:**

- Infobereich-Anwendungen, die Fenster in den Infobereich ausblenden
- Assistentenabläufe, in denen Fenster wiederverwendet werden
- Vorübergehendes Ausblenden während Vorgängen

### Close()

Schließt das Fenster. Dadurch wird das Ereignis `WindowClosing` ausgelöst.

```go
func (w *WebviewWindow) Close()
```

**Beispiel:**

```go
window.Close()
```

**Hinweis:** Wenn ein registrierter Hook `event.Cancel()` aufruft, wird das Schließen verhindert.

## Fenstereigenschaften

### SetTitle()

Legt den Text in der Titelleiste des Fensters fest. Gibt den Empfänger zur Verkettung zurück.

```go
func (w *WebviewWindow) SetTitle(title string) Window
```

**Parameter:**

- `title` – Der neue Fenstertitel

**Beispiel:**

```go
window.SetTitle("My Application - Document.txt")
```

### Name()

Gibt den eindeutigen Namensbezeichner des Fensters zurück.

```go
func (w *WebviewWindow) Name() string
```

**Beispiel:**

```go
name := window.Name()
fmt.Println("Window name:", name)

// Retrieve window by name later
if w, ok := app.Window.GetByName(name); ok {
    w.Focus()
}
```

## Größe und Position

### SetSize()

Legt die Fensterabmessungen in Pixeln fest. Gibt den Empfänger zur Verkettung zurück.

```go
func (w *WebviewWindow) SetSize(width, height int) Window
```

**Parameter:**

- `width` – Fensterbreite in Pixeln
- `height` – Fensterhöhe in Pixeln

**Beispiel:**

```go
window.SetSize(1024, 768)
```

### Size()

Gibt die aktuellen Fensterabmessungen zurück.

```go
func (w *WebviewWindow) Size() (width, height int)
```

**Beispiel:**

```go
width, height := window.Size()
fmt.Printf("Window is %dx%d\n", width, height)
```

### SetMinSize() / SetMaxSize()

Legt die minimalen und maximalen Fensterabmessungen fest. Beide Methoden geben den Empfänger zur Verkettung zurück.

```go
func (w *WebviewWindow) SetMinSize(width, height int) Window
func (w *WebviewWindow) SetMaxSize(width, height int) Window
```

**Beispiel:**

```go
// Prevent window from being too small
window.SetMinSize(800, 600)

// Prevent window from being too large
window.SetMaxSize(1920, 1080)
```

### SetPosition()

Legt die Fensterposition relativ zur linken oberen Ecke des Bildschirms fest.

```go
func (w *WebviewWindow) SetPosition(x, y int)
```

**Parameter:**

- `x` – Horizontale Position in Pixeln
- `y` – Vertikale Position in Pixeln

**Beispiel:**

```go
// Position window at top-left
window.SetPosition(0, 0)

// Position window 100px from top-left
window.SetPosition(100, 100)
```

### Position()

Gibt die aktuelle Fensterposition zurück.

```go
func (w *WebviewWindow) Position() (x, y int)
```

**Beispiel:**

```go
x, y := window.Position()
fmt.Printf("Window is at (%d, %d)\n", x, y)
```

### Center()

Zentriert das Fenster auf dem Bildschirm.

```go
func (w *WebviewWindow) Center()
```

**Beispiel:**

```go
window := app.Window.New()
window.Center()
window.Show()
```

**Hinweis:** Zentriert das Fenster auf dem primären Monitor. Informationen zu Konfigurationen mit mehreren Monitoren finden Sie in den Bildschirm-APIs.

### Focus()

Bringt das Fenster in den Vordergrund und weist ihm den Tastaturfokus zu.

```go
func (w *WebviewWindow) Focus()
```

**Beispiel:**

```go
// Bring window to front
window.Focus()
```

## Fensterzustand

### Minimise() / UnMinimise()

Minimiert das Fenster in die Taskleiste bzw. das Dock oder stellt es wieder her. `Minimise()` gibt den Empfänger zur Verkettung zurück; `UnMinimise()` gibt nichts zurück.

```go
func (w *WebviewWindow) Minimise() Window
func (w *WebviewWindow) UnMinimise()
```

**Beispiel:**

```go
// Minimise window
window.Minimise()

// Restore from minimised state
window.UnMinimise()
```

### Maximise() / UnMaximise()

Maximiert das Fenster auf die volle Bildschirmgröße oder stellt seine vorherige Größe wieder her. `Maximise()` gibt den Empfänger zur Verkettung zurück; `UnMaximise()` gibt nichts zurück.

```go
func (w *WebviewWindow) Maximise() Window
func (w *WebviewWindow) UnMaximise()
```

**Beispiel:**

```go
// Maximise window
window.Maximise()

// Restore to previous size
window.UnMaximise()
```

### Fullscreen() / UnFullscreen() / ToggleFullscreen()

Aktiviert oder beendet den Vollbildmodus. `Fullscreen()` gibt den Empfänger zur Verkettung zurück.

```go
func (w *WebviewWindow) Fullscreen() Window
func (w *WebviewWindow) UnFullscreen()
func (w *WebviewWindow) ToggleFullscreen()
```

**Beispiel:**

```go
// Enter fullscreen
window.Fullscreen()

// Exit fullscreen
window.UnFullscreen()

// Or toggle
window.ToggleFullscreen()
```

Es gibt keine Methode `SetFullscreen(bool)`.

### IsMinimised() / IsMaximised() / IsFullscreen()

Prüft den aktuellen Fensterzustand.

```go
func (w *WebviewWindow) IsMinimised() bool
func (w *WebviewWindow) IsMaximised() bool
func (w *WebviewWindow) IsFullscreen() bool
```

**Beispiel:**

```go
if window.IsMinimised() {
    window.UnMinimise()
}

if window.IsMaximised() {
    fmt.Println("Window is maximised")
}

if window.IsFullscreen() {
    window.UnFullscreen()
}
```

## Fensterinhalt

### SetURL()

Navigiert innerhalb des Fensters zu einer bestimmten URL. Gibt den Empfänger zur Methodenverkettung zurück.

```go
func (w *WebviewWindow) SetURL(url string) Window
```

**Parameter:**

- `url` – URL, zu der navigiert werden soll (für eingebettete Assets kann `http://wails.localhost/` verwendet werden)

**Beispiel:**

```go
// Navigate to embedded page
window.SetURL("http://wails.localhost/settings.html")

// Navigate to external URL (if allowed)
window.SetURL("https://wails.io")
```

### SetHTML()

Legt den Fensterinhalt direkt aus einer HTML-Zeichenfolge fest. Gibt den Empfänger zur Methodenverkettung zurück.

```go
func (w *WebviewWindow) SetHTML(html string) Window
```

**Parameter:**

- `html` – anzuzeigender HTML-Inhalt

**Beispiel:**

```go
html := `
<!DOCTYPE html>
<html>
<head><title>Dynamic Content</title></head>
<body>
    <h1>Hello from Go!</h1>
    <p>This content was generated dynamically.</p>
</body>
</html>
`
window.SetHTML(html)
```

**Anwendungsfälle:**

- Dynamische Inhaltserzeugung
- Einfache Fenster ohne Frontend-Buildprozess
- Fehlerseiten oder Startbildschirme

### Reload()

Lädt den aktuellen Fensterinhalt neu.

```go
func (w *WebviewWindow) Reload()
```

**Beispiel:**

```go
// Reload current page
window.Reload()
```

**Hinweis:** Nützlich während der Entwicklung oder wenn der Inhalt aktualisiert werden muss.

## Fensterereignisse

Wails bietet zwei Methoden zur Behandlung von Fensterereignissen:

- **OnWindowEvent()** – Lauscht auf Fensterereignisse (kann sie nicht verhindern).
- **RegisterHook()** – Bindet sich in Fensterereignisse ein (kann sie durch Aufrufen von `event.Cancel()` verhindern).

### OnWindowEvent()

Registriert einen Callback für Fensterereignisse. Gibt eine Funktion zum Abbestellen zurück.

```go
func (w *WebviewWindow) OnWindowEvent(
    eventType events.WindowEventType,
    callback func(event *WindowEvent),
) func()
```

**Beispiel:**

```go
import "github.com/wailsapp/wails/v3/pkg/events"

// Listen for window focus
window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
    app.Logger.Info("Window gained focus")
})

// Listen for window lost focus
window.OnWindowEvent(events.Common.WindowLostFocus, func(e *application.WindowEvent) {
    app.Logger.Info("Window lost focus")
})

// Listen for window resize
window.OnWindowEvent(events.Common.WindowDidResize, func(e *application.WindowEvent) {
    app.Logger.Info("Window resized")
})
```

**Häufige Fensterereignisse:**

- `events.Common.WindowClosing` – Fenster wird gleich geschlossen
- `events.Common.WindowFocus` – Fenster hat den Fokus erhalten
- `events.Common.WindowLostFocus` – Fenster hat den Fokus verloren
- `events.Common.WindowDidMove` – Fenster wurde verschoben
- `events.Common.WindowDidResize` – Fenstergröße wurde geändert
- `events.Common.WindowMinimise` – Fenster wurde minimiert
- `events.Common.WindowMaximise` – Fenster wurde maximiert
- `events.Common.WindowFullscreen` – Fenster ist in den Vollbildmodus gewechselt
- `events.Common.WindowRuntimeReady` – Laufzeitumgebung im Fenster wurde initialisiert

### RegisterHook()

Registriert einen Hook für Fensterereignisse. Hooks werden vor Listenern ausgeführt und können das Ereignis durch Aufrufen von `event.Cancel()` verhindern. Gibt eine Funktion zum Abbestellen zurück.

```go
func (w *WebviewWindow) RegisterHook(
    eventType events.WindowEventType,
    callback func(event *WindowEvent),
) func()
```

**Beispiel – Schließen des Fensters verhindern:**

```go
import "github.com/wailsapp/wails/v3/pkg/events"

window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    confirm := app.Dialog.Question().
        SetTitle("Confirm Close").
        SetMessage("Are you sure you want to close this window?")

    yes := confirm.AddButton("Yes")
    no := confirm.AddButton("No")
    confirm.SetDefaultButton(yes)
    confirm.SetCancelButton(no)

    no.OnClick(func() {
        e.Cancel() // Prevent window from closing
    })

    confirm.Show()
})
```

**Beispiel – Vor dem Schließen speichern:**

```go
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    if !hasUnsavedChanges {
        return
    }

    dlg := app.Dialog.Question().
        SetTitle("Unsaved Changes").
        SetMessage("Save changes before closing?")

    save := dlg.AddButton("Save")
    discard := dlg.AddButton("Don't Save")
    cancel := dlg.AddButton("Cancel")
    dlg.SetDefaultButton(save)
    dlg.SetCancelButton(cancel)

    save.OnClick(func() { saveData() })
    cancel.OnClick(func() { e.Cancel() })
    _ = discard // allow close

    dlg.Show()
})
```

### EmitEvent()

Sendet ein benutzerdefiniertes Ereignis an das Frontend des Fensters. Gibt `true` zurück, wenn das Senden durch einen Hook abgebrochen wurde.

```go
func (w *WebviewWindow) EmitEvent(name string, data ...any) bool
```

**Parameter:**

- `name` – Ereignisname
- `data` – optionale Daten, die mit dem Ereignis gesendet werden sollen

**Beispiel:**

```go
// Send data to specific window
window.EmitEvent("data-updated", map[string]any{
    "count":  42,
    "status": "success",
})
```

**Frontend (JavaScript):**

```javascript
import { Events } from '@wailsio/runtime'

Events.On('data-updated', (data) => {
    console.log('Count:', data.count)
    console.log('Status:', data.status)
})
```

## Weitere Methoden

### SetEnabled()

Aktiviert oder deaktiviert die Benutzerinteraktion mit dem Fenster.

```go
func (w *WebviewWindow) SetEnabled(enabled bool)
```

**Beispiel:**

```go
// Disable window during long operation
window.SetEnabled(false)

// Perform operation
performLongOperation()

// Re-enable window
window.SetEnabled(true)
```

### SetBackgroundColour()

Legt die Hintergrundfarbe des Fensters fest, die vor dem Laden des Inhalts angezeigt wird. Gibt den Empfänger zur Methodenverkettung zurück.

```go
func (w *WebviewWindow) SetBackgroundColour(colour RGBA) Window
```

`RGBA` ist `application.RGBA{Red, Green, Blue, Alpha uint8}`. Verwenden Sie die Hilfsfunktionen `application.NewRGB(r, g, b)` (Alpha 255) oder `application.NewRGBA(r, g, b, a)`.

**Beispiel:**

```go
// White background
window.SetBackgroundColour(application.NewRGB(255, 255, 255))

// Dark background with full alpha
window.SetBackgroundColour(application.NewRGBA(30, 30, 30, 255))
```

### SetResizable()

Steuert, ob der Benutzer die Fenstergröße ändern kann. Gibt den Empfänger zur Methodenverkettung zurück.

```go
func (w *WebviewWindow) SetResizable(resizable bool) Window
```

**Beispiel:**

```go
// Make window fixed size
window.SetResizable(false)
```

### SetAlwaysOnTop()

Legt fest, ob das Fenster über anderen Fenstern bleibt. Gibt den Empfänger zur Methodenverkettung zurück.

```go
func (w *WebviewWindow) SetAlwaysOnTop(alwaysOnTop bool) Window
```

**Beispiel:**

```go
// Keep window on top
window.SetAlwaysOnTop(true)
```

### Print()

Öffnet den nativen Druckdialog für den Fensterinhalt.

```go
func (w *WebviewWindow) Print() error
```

**Rückgabewert:** Fehler, wenn das Drucken fehlschlägt.

**Beispiel:**

```go
if err := window.Print(); err != nil {
    log.Println("Print failed:", err)
}
```

### AttachModal()

Hängt ein zweites Fenster als modales Dialogfenster in Form eines Sheets an.

```go
func (w *WebviewWindow) AttachModal(modalWindow Window)
```

**Parameter:**

- `modalWindow` – Das als modales Dialogfenster anzuhängende Fenster

**Plattformunterstützung:**

- **macOS**: Vollständig unterstützt (Darstellung als Sheet)
- **Windows**: Nicht unterstützt
- **Linux**: Nicht unterstützt

**Beispiel:**

```go
modalWindow := app.Window.New()
window.AttachModal(modalWindow)
```

## Plattformspezifische Optionen

### Linux

Linux-Fenster unterstützen über `LinuxWindow` die folgenden plattformspezifischen Optionen:

#### MenuStyle

Steuert die Darstellung des Anwendungsmenüs. Diese Option ist im standardmäßigen GTK4-Build verfügbar und wird in älteren `-tags gtk3`-Builds ignoriert.

| Wert | Beschreibung |
| --- | --- |
| `LinuxMenuStyleMenuBar` | Herkömmliche Menüleiste unterhalb der Titelleiste (Standard) |
| `LinuxMenuStylePrimaryMenu` | Primäre Menüschaltfläche in der Kopfleiste (GNOME-Stil) |

**Beispiel:**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "My Application",
    Linux: application.LinuxWindow{
        MenuStyle: application.LinuxMenuStylePrimaryMenu,
    },
})
window.SetMenu(menu)
```

**Hinweis:** Beim primären Menüstil wird gemäß den GNOME Human Interface Guidelines eine Hamburger-Schaltfläche (☰) in der Kopfleiste angezeigt. Dieser Stil wird für moderne GNOME-Anwendungen empfohlen.

## Vollständiges Beispiel

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)

func main() {
    app := application.New(application.Options{
        Name: "Window API Demo",
    })

    // Create window with options
    window := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:            "My Application",
        Width:            1024,
        Height:           768,
        MinWidth:         800,
        MinHeight:        600,
        BackgroundColour: application.NewRGB(255, 255, 255),
        URL:              "http://wails.localhost/",
    })

    // Configure window behaviour
    window.SetResizable(true)
    window.SetMinSize(800, 600)
    window.SetMaxSize(1920, 1080)

    // Confirm-before-close hook
    window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
        dlg := app.Dialog.Question().
            SetTitle("Confirm Close").
            SetMessage("Are you sure you want to close this window?")

        yes := dlg.AddButton("Yes")
        no := dlg.AddButton("No")
        dlg.SetDefaultButton(yes)
        dlg.SetCancelButton(no)
        no.OnClick(func() { e.Cancel() })

        dlg.Show()
    })

    // Listen for window events
    window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
        window.SetTitle("My Application (Active)")
        app.Logger.Info("Window gained focus")
    })

    window.OnWindowEvent(events.Common.WindowLostFocus, func(e *application.WindowEvent) {
        window.SetTitle("My Application")
        app.Logger.Info("Window lost focus")
    })

    // Position and show window
    window.Center()
    window.Show()

    app.Run()
}
```
