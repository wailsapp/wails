---
title: "Manager-API"
description: "Strukturierte API mit spezialisierten Manager-Schnittstellen"
slug: "concepts/manager-api"
sourcePath: "concepts/manager-api.md"
---

Die Manager-API von Wails v3 bietet einen strukturierten und leicht auffindbaren Zugriff auf Anwendungsfunktionen über spezialisierte Manager-Structs, die in öffentlichen Feldern von `*application.App` gruppiert sind. Wails 3 stellt einen klaren Bruch mit v2 dar – es gibt keine Wrapper-Schicht für einzelne Aufrufe, die die Kompatibilität mit der alten API im Stil von `app.NewWebviewWindow(...)` aufrechterhält. Verwenden Sie daher die folgenden Manager, um die Anwendung zu steuern.

## Überblick

Die Manager-API gliedert die Anwendungsfunktionen in zwölf spezialisierte Bereiche (einen Logger und elf Manager):

- **`app.Window`** – Erstellen und Verwalten von Fenstern sowie Callbacks
- **`app.ContextMenu`** – Registrieren und Verwalten von Kontextmenüs\
- **`app.KeyBinding`** – Verwalten globaler Tastenbelegungen
- **`app.Browser`** – Browser-Integration (Öffnen von URLs und Dateien)
- **`app.Env`** – Informationen zur Umgebung und zum Systemzustand
- **`app.Dialog`** – Datei- und Meldungsdialoge
- **`app.Event`** – Verarbeiten benutzerdefinierter Ereignisse und Anwendungsereignisse
- **`app.Menu`** – Verwalten des Anwendungsmenüs
- **`app.Screen`** – Verwalten von Bildschirmen und Transformieren von Koordinaten
- **`app.Clipboard`** – Textoperationen in der Zwischenablage
- **`app.SystemTray`** – Erstellen und Verwalten von Symbolen im Infobereich
- **`app.Autostart`** – Registrieren der Anwendung für den Start bei der Benutzeranmeldung

## Vorteile

- **Bessere Auffindbarkeit** – Die automatische IDE-Vervollständigung zeigt eine strukturierte API-Oberfläche an
- **Verbesserte Codestruktur** – Zusammengehörige Methoden sind gruppiert
- **Bessere Wartbarkeit** – Zuständigkeiten sind auf die Manager verteilt
- **Künftige Erweiterbarkeit** – Neue Funktionen lassen sich einfacher zu bestimmten Bereichen hinzufügen

## Verwendung

Die Manager-API bietet strukturierten Zugriff auf alle Anwendungsfunktionen:

```go
// Events and custom event handling
app.Event.Emit("custom", data)
app.Event.On("custom", func(e *CustomEvent) { ... })

// Window management
window, _ := app.Window.GetByName("main")
app.Window.OnCreate(func(window Window) { ... })

// Browser integration
app.Browser.OpenURL("https://wails.io")

// Menu management
menu := app.Menu.New()
app.Menu.Set(menu)

// System tray
systray := app.SystemTray.New()
```

## Manager-Referenz

### Fenster-Manager

Verwaltet das Erstellen und Abrufen von Fenstern sowie Callbacks für deren Lebenszyklus.

```go
// Create windows
window := app.Window.New()
window := app.Window.NewWithOptions(options)
current := app.Window.Current()

// Find windows
window, exists := app.Window.GetByName("main")
windows := app.Window.GetAll()

// Window callbacks
app.Window.OnCreate(func(window Window) {
    // Handle window creation
})
```

### Ereignis-Manager

Verarbeitet benutzerdefinierte Ereignisse und überwacht Anwendungsereignisse.

```go
// Custom events
app.Event.Emit("userAction", data)
cancelFunc := app.Event.On("userAction", func(e *CustomEvent) {
    // Handle event
})
app.Event.Off("userAction")
app.Event.Reset() // Remove all listeners

// Application events
app.Event.OnApplicationEvent(events.Common.ThemeChanged, func(e *ApplicationEvent) {
    // Handle system theme change
})
```

### Browser-Manager

Stellt die Browser-Integration zum Öffnen von URLs und Dateien bereit.

```go
// Open URLs and files in default browser
err := app.Browser.OpenURL("https://wails.io")
err := app.Browser.OpenFile("/path/to/document.pdf")
```

### Umgebungs-Manager

Ermöglicht den Zugriff auf Informationen zur Systemumgebung.

```go
// Get environment info
env := app.Env.Info()
fmt.Printf("OS: %s, Arch: %s\n", env.OS, env.Arch)

// Check system theme
if app.Env.IsDarkMode() {
    // Dark mode is active
}

// Open file manager
err := app.Env.OpenFileManager("/path/to/folder", false)
```

### Dialog-Manager

Bietet strukturierten Zugriff auf Datei- und Meldungsdialoge.

```go
// File dialogs
result, err := app.Dialog.OpenFile().
    AddFilter("Text Files", "*.txt").
    PromptForSingleSelection()

result, err = app.Dialog.SaveFile().
    SetFilename("document.txt").
    PromptForSingleSelection()

// Message dialogs
app.Dialog.Info().
    SetTitle("Information").
    SetMessage("Operation completed successfully").
    Show()

app.Dialog.Error().
    SetTitle("Error").
    SetMessage("An error occurred").
    Show()
```

### Menü-Manager

Erstellt und verwaltet das Anwendungsmenü.

```go
// Create and set application menu
menu := app.Menu.New()
fileMenu := menu.AddSubmenu("File")
fileMenu.Add("New").OnClick(func(ctx *Context) {
    // Handle menu click
})

app.Menu.Set(menu)

// Show about dialog
app.Menu.ShowAbout()
```

### Tastenbelegungs-Manager

Verwaltet globale Tastenbelegungen dynamisch.

```go
// Add key bindings
app.KeyBinding.Add("ctrl+n", func(window application.Window) {
    // Handle Ctrl+N
})

app.KeyBinding.Add("ctrl+q", func(window application.Window) {
    app.Quit()
})

// Remove key bindings
app.KeyBinding.Remove("ctrl+n")

// Get all bindings
bindings := app.KeyBinding.GetAll()
```

### Kontextmenü-Manager

Erweiterte Verwaltung von Kontextmenüs (für Bibliotheksautoren).

```go
// Create and register context menu
menu := app.ContextMenu.New()
app.ContextMenu.Add("myMenu", menu)

// Retrieve context menu
menu, exists := app.ContextMenu.Get("myMenu")

// Remove context menu
app.ContextMenu.Remove("myMenu")
```

### Bildschirm-Manager

Verwaltet Bildschirme und transformiert Koordinaten für Konfigurationen mit mehreren Monitoren.

```go
// Get screen information
screens := app.Screen.GetAll()
primary := app.Screen.GetPrimary()

// Coordinate transformations
physicalPoint := app.Screen.DipToPhysicalPoint(logicalPoint)
logicalPoint := app.Screen.PhysicalToDipPoint(physicalPoint)

// Screen detection
screen := app.Screen.ScreenNearestDipPoint(point)
screen = app.Screen.ScreenNearestDipRect(rect)
```

### Zwischenablage-Manager

Operationen zum Lesen und Schreiben von Text in der Zwischenablage.

```go
// Set text to clipboard
success := app.Clipboard.SetText("Hello World")
if !success {
    // Handle error
}

// Get text from clipboard
text, ok := app.Clipboard.Text()
if !ok {
    // Handle error
} else {
    // Use the text
}
```

### SystemTray-Manager

Erstellt und verwaltet Symbole im Infobereich.

```go
// Create system tray
systray := app.SystemTray.New()
systray.SetLabel("My App")
systray.SetIcon(iconBytes)

// Add menu to system tray
menu := app.Menu.New()
menu.Add("Open").OnClick(func(ctx *Context) {
    // Handle click
})
systray.SetMenu(menu)

// Destroy system tray when done
systray.Destroy()
```

### Autostart-Manager

Registriert die Anwendung für den Start bei der Benutzeranmeldung. Wählt für jede Plattform den geeigneten nativen Mechanismus aus: SMAppService oder eine LaunchAgent-plist-Datei unter macOS, den Registrierungsschlüssel `HKCU\…\Run` unter Windows und einen XDG-`.desktop`-Eintrag unter Linux.

```go
// Register to launch at login
err := app.Autostart.Enable()

// With extra launch-time arguments and a custom identifier
err = app.Autostart.EnableWithOptions(application.AutostartOptions{
    Identifier: "com.example.myapp",
    Arguments:  []string{"--hidden"},
})

// Check / remove
enabled, err := app.Autostart.IsEnabled()
status, err := app.Autostart.Status()  // includes Path + Strategy
err = app.Autostart.Disable()
```

Informationen zum plattformspezifischen Verhalten, zu den Regeln für Bezeichner und zur Garantie für die Erkennung veralteter Einträge finden Sie auf der Seite [Autostart-Funktion](/features/autostart/basics/).
