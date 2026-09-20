---
title: "Fensteroptionen"
description: "Vollständige Referenz für WebviewWindowOptions"
slug: "features/windows/options"
sourcePath: "features/windows/options.md"
---

## Optionen zur Fensterkonfiguration

Wails bietet eine umfassende Fensterkonfiguration mit zahlreichen Optionen für Größe, Position, Darstellung und Verhalten. Diese **vollständige Referenz** für `WebviewWindowOptions` behandelt alle verfügbaren Optionen unter Windows, macOS und Linux. Sie enthält jede Option für jede Plattform einschließlich Beispielen und Einschränkungen.

## Struktur von WebviewWindowOptions

```go
type WebviewWindowOptions struct {
    // Identity
    Name  string
    Title string

    // Size and Position
    Width           int
    Height          int
    X               int
    Y               int
    MinWidth        int
    MinHeight       int
    MaxWidth        int
    MaxHeight       int
    InitialPosition WindowStartPosition // WindowCentered (default) or WindowXY
    Screen          *Screen             // target screen for initial placement

    // Initial State
    Hidden        bool
    Frameless     bool
    DisableResize bool        // inverted vs v2's `Resizable`
    AlwaysOnTop   bool
    StartState    WindowState // WindowStateNormal | Minimised | Maximised | Fullscreen

    // Appearance
    BackgroundColour RGBA
    BackgroundType   BackgroundType
    Zoom             float64
    ZoomControlEnabled bool

    // Content
    URL  string
    HTML string
    JS   string
    CSS  string

    // Behaviour
    EnableFileDrop              bool
    IgnoreMouseEvents           bool
    HideOnFocusLost             bool
    HideOnEscape                bool
    DevToolsEnabled             bool
    DefaultContextMenuDisabled  bool
    ContentProtectionEnabled    bool
    KeyBindings                 map[string]func(window *WebviewWindow)

    // Permissions
    Permissions map[PermissionType]Permission

    // Window-control button states
    MinimiseButtonState ButtonState
    MaximiseButtonState ButtonState
    CloseButtonState    ButtonState

    // Menu
    UseApplicationMenu bool

    // Platform-specific (per-window)
    Mac     MacWindow
    Windows WindowsWindow
    Linux   LinuxWindow
}
```

`WebviewWindowOptions` besitzt **kein** Feld `Parent`. Verwenden Sie für Eltern-/Modalbeziehungen `parentWindow.AttachModal(childWindow)`. Ebenso besitzt es **kein** Feld `Assets`. Die Asset-Konfiguration befindet sich in `application.Options` (`Assets AssetOptions`).

Vollständiger Quellcode: [`v3/pkg/application/webview_window_options.go`](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window_options.go).

## Kernoptionen

### Name

**Typ:** `string` **Standardwert:** Automatisch generierte UUID **Plattform:** Alle

```go
Name: "main-window"
```

**Zweck:** Eindeutige Kennung, über die das Fenster später gefunden werden kann.

**Bewährte Vorgehensweisen:**

- Verwenden Sie aussagekräftige Namen: `"main"`, `"settings"`, `"about"`
- Verwenden Sie Kebab-Case: `"file-browser"`, `"color-picker"`
- Halten Sie den Namen kurz und einprägsam

**Beispiel:**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name: "settings-window",
})

// Later...
if settings, ok := app.Window.GetByName("settings-window"); ok {
    settings.Focus()
}
```

### Title

**Typ:** `string` **Standardwert:** Anwendungsname **Plattform:** Alle

```go
Title: "My Application"
```

**Zweck:** Text, der in der Titelleiste und Taskleiste angezeigt wird.

**Dynamische Aktualisierungen:**

```go
window.SetTitle("My Application - Document.txt")
```

### Width / Height

**Typ:** `int` (Pixel) **Standardwert:** 800 × 600 **Plattform:** Alle **Einschränkungen:** Muss positiv sein

```go
Width:  1200,
Height: 800,
```

**Zweck:** Anfängliche Fenstergröße in logischen Pixeln.

**Hinweise:**

- Wails berücksichtigt die DPI-Skalierung automatisch
- Verwenden Sie logische statt physischer Pixel
- Berücksichtigen Sie die minimale Bildschirmauflösung (1024x768)

**Beispielgrößen:**

| Anwendungsfall | Breite | Höhe |
| --- | --- | --- |
| Kleines Dienstprogramm | 400 | 300 |
| Standardanwendung | 1024 | 768 |
| Große Anwendung | 1440 | 900 |
| Full HD | 1920 | 1080 |

### X / Y

**Typ:** `int` (Pixel) **Standardwert:** Auf dem Bildschirm zentriert **Plattform:** Alle

```go
X: 100,  // 100px from left edge
Y: 100,  // 100px from top edge
```

**Zweck:** Anfängliche Fensterposition.

**Koordinatensystem:**

- (0, 0) liegt oben links auf dem primären Bildschirm
- Positive X-Werte verlaufen nach rechts
- Positive Y-Werte verlaufen nach unten

**Beispiel:**

`X` und `Y` werden nur wirksam, wenn `InitialPosition: application.WindowXY` festgelegt ist. Andernfalls verwendet `InitialPosition` standardmäßig `WindowCentered` und `X`/`Y` werden ignoriert.

```go
settings := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:            "coordinate-window",
    InitialPosition: application.WindowXY, // opt into X/Y coordinates
    X:               100,
    Y:               100,
})
```

**Bewährte Vorgehensweise:** Wenn bestimmte Koordinaten keine Rolle spielen, zentrieren Sie das Fenster nach seiner Erstellung mit `Center()`:

```go
window := app.Window.New()
window.Center()
```

### MinWidth / MinHeight

**Typ:** `int` (Pixel) **Standardwert:** 0 (kein Minimum) **Plattform:** Alle

```go
MinWidth:  400,
MinHeight: 300,
```

**Zweck:** Verhindert, dass das Fenster zu klein wird.

**Anwendungsfälle:**

- Fehlerhafte Layouts verhindern
- Benutzbarkeit gewährleisten
- Seitenverhältnis beibehalten

**Beispiel:**

```go
// Prevent window smaller than 400x300
MinWidth:  400,
MinHeight: 300,
```

### MaxWidth / MaxHeight

**Typ:** `int` (Pixel) **Standardwert:** 0 (kein Maximum) **Plattform:** Alle

```go
MaxWidth:  1920,
MaxHeight: 1080,
```

**Zweck:** Verhindert, dass das Fenster zu groß wird.

**Anwendungsfälle:**

- Anwendungen mit fester Größe
- Übermäßige Ressourcennutzung verhindern
- Designvorgaben einhalten

## Statusoptionen

### Hidden

**Typ:** `bool` **Standardwert:** `false` **Plattform:** Alle

```go
Hidden: true,
```

**Zweck:** Fenster erstellen, ohne es anzuzeigen.

**Anwendungsfälle:**

- Hintergrundfenster
- Bei Bedarf angezeigte Fenster
- Startbildschirme (erstellen, laden und anschließend anzeigen)
- Weißes Aufblitzen beim Laden von Inhalten verhindern

**Plattformspezifische Verbesserungen:**

- **Windows:** Weißes Aufblitzen des Fensters behoben – das Fenster bleibt unsichtbar, bis `Show()` aufgerufen wird
- **macOS:** Vollständig unterstützt
- **Linux:** Vollständig unterstützt

**Empfohlenes Muster für flüssiges Laden:**

```go
// Create hidden window
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:             "main-window",
    Hidden:           true,
    BackgroundColour: application.NewRGB(30, 30, 30), // Match your theme
})

// Load content while hidden
// ... content loads ...

// Show when ready (no flash!)
window.Show()
```

**Beispiel:**

```go
settings := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:   "settings",
    Hidden: true,
})

// Show when needed
settings.Show()
```

### Frameless

**Typ:** `bool` **Standardwert:** `false` **Plattform:** Alle

```go
Frameless: true,
```

**Zweck:** Titelleiste und Fensterrahmen entfernen.

**Anwendungsfälle:**

- Benutzerdefinierte Fensterdekoration
- Startbildschirme
- Kioskanwendungen
- Individuell gestaltete Fenster

**Wichtig:** Folgendes müssen Sie implementieren:

- Verschieben des Fensters
- Schaltflächen zum Schließen, Minimieren und Maximieren
- Griffe zur Größenänderung (wenn die Größe veränderbar ist)

**Weitere Informationen finden Sie unter [Rahmenlose Fenster](/features/windows/frameless/).**

### DisableResize

**Typ:** `bool` **Standardwert:** `false` (die Fenstergröße ist standardmäßig veränderbar) **Plattform:** Alle

```go
DisableResize: true,
```

**Zweck:** Größenänderung des Fensters verhindern. Beachten Sie, dass dieses Feld die **Umkehrung** von `Resizable` aus v2 ist. Legen Sie `DisableResize: true` fest, damit die Größe eines Fensters nicht veränderbar ist.

**Anwendungsfälle:**

- Anwendungen mit fester Fenstergröße
- Startbildschirme
- Dialogfelder

**Hinweis:** Benutzer können das Fenster weiterhin maximieren oder in den Vollbildmodus versetzen, sofern Sie dies nicht zusätzlich über `MaximiseButtonState` bzw. das Verhalten der Collection deaktivieren.

### AlwaysOnTop

**Typ:** `bool` **Standardwert:** `false` **Plattform:** Alle

```go
AlwaysOnTop: true,
```

**Zweck:** Fenster über allen anderen Fenstern halten.

**Anwendungsfälle:**

- Schwebende Symbolleisten
- Benachrichtigungen
- Bild-im-Bild
- Zeitgeber

**Plattformhinweise:**

- **macOS:** Vollständig unterstützt
- **Windows:** Vollständig unterstützt
- **Linux:** Abhängig vom Fenstermanager

### StartState

**Typ:** `WindowState`-Enum **Standardwert:** `WindowStateNormal` **Plattform:** Alle

```go
StartState: application.WindowStateMaximised,
```

**Zweck:** Anfangszustand des Fensters bei der Anzeige.

**Werte:**

- `WindowStateNormal` – Normales Fenster
- `WindowStateMinimised` – Minimiert
- `WindowStateMaximised` – Maximiert
- `WindowStateFullscreen` – Vollbild

Es gibt keine `WindowStateHidden`-Konstante. Verwenden Sie das boolesche Feld `Hidden`, um das Fenster zunächst unsichtbar zu öffnen.

**Vollbildmodus zur Laufzeit umschalten:**

```go
window.Fullscreen()
window.UnFullscreen()
window.ToggleFullscreen() // there is no SetFullscreen(bool)
```

## Darstellungsoptionen

### BackgroundColour

**Typ:** `RGBA`-Struktur **Standard:** Weiß **Plattform:** Alle

```go
BackgroundColour: application.RGBA{Red: 0, Green: 0, Blue: 0, Alpha: 255},
```

Die Felder von `RGBA` heißen `Red, Green, Blue, Alpha` und sind jeweils vom Typ uint8. Verwenden Sie vorzugsweise die Hilfsfunktionen `application.NewRGB(r, g, b)` (Alpha 255) oder `application.NewRGBA(r, g, b, a)`.

**Zweck:** Hintergrundfarbe des Fensters, bevor der Inhalt geladen wird.

**Anwendungsfälle:**

- An das Theme Ihrer Anwendung anpassen
- Weißes Aufblitzen bei dunklen Themes verhindern
- Für einen flüssigen Ladevorgang sorgen

**Beispiel:**

```go
// Dark theme
BackgroundColour: application.NewRGB(30, 30, 30),

// Light theme
BackgroundColour: application.NewRGB(255, 255, 255),
```

**Hilfsmethode:**

```go
window.SetBackgroundColour(application.NewRGB(30, 30, 30))
```

### BackgroundType

**Typ:** `BackgroundType`-Enumeration **Standard:** `BackgroundTypeSolid` **Plattform:** macOS, Windows (teilweise)

```go
BackgroundType: application.BackgroundTypeTranslucent,
```

**Werte:**

- `BackgroundTypeSolid` – Einfarbig
- `BackgroundTypeTransparent` – Vollständig transparent
- `BackgroundTypeTranslucent` – Halbtransparente Unschärfe

**Plattformunterstützung:**

- **macOS:** Konfigurieren Sie `Mac.Backdrop`; für die Transparenz der Webview ist [`private_mac_apis`](/guides/build/private-macos-apis/#webview-transparency-and-background) erforderlich. Ohne diese Einstellung bleibt die Webview undurchsichtig.
- **Windows:** Transparent und halbtransparent (Windows 11+)
- **Linux:** Nur einfarbig

**Beispiel (macOS):**

```go
BackgroundType: application.BackgroundTypeTranslucent,
Mac: application.MacWindow{
    Backdrop: application.MacBackdropTranslucent,
},
```

### OpenInspectorOnStartup und OpenDevTools

**Private API unter macOS:** Für `OpenInspectorOnStartup: true`, Go-`window.OpenDevTools()` und JavaScript-`Window.OpenDevTools()` ist `private_mac_apis` erforderlich, um den Inspektor programmgesteuert zu öffnen. Ohne diese Einstellung bewirken diese Operationen nichts. Produktions-Builds benötigen außerdem `devtools`. Für die öffentliche Safari-Inspektion unter macOS 13.3+ sind keine privaten APIs erforderlich; für die Aktivierung des Inspektors unter älteren macOS-Versionen hingegen schon. Weitere Informationen finden Sie in der [Build-Matrix für den Web Inspector](/guides/build/private-macos-apis/#web-inspector).

## Inhaltsoptionen

### URL

**Typ:** `string` **Standard:** Leer (lädt aus Assets) **Plattform:** Alle

```go
URL: "https://example.com",
```

**Zweck:** Eine externe URL anstelle eingebetteter Assets laden.

**Anwendungsfälle:**

- Entwicklung (vom Entwicklungsserver laden)
- Webbasierte Anwendungen
- Hybride Anwendungen

**Beispiel:**

```go
// Development — point the window at the Vite dev server
URL: "http://localhost:9245",

// Production — embedded assets are configured at the application level
// (Assets is application.Options.Assets, not a WebviewWindowOptions field).
```

Anwendungsweites Codefragment für den Produktionseinsatz:

```go
app := application.New(application.Options{
    Name: "My App",
    Assets: application.AssetOptions{
        Handler: application.AssetFileServerFS(assets),
    },
})
```

### HTML

**Typ:** `string` **Standard:** Leer **Plattform:** Alle

```go
HTML: "<h1>Hello World</h1>",
```

**Zweck:** Eine HTML-Zeichenkette direkt laden.

**Anwendungsfälle:**

- Einfache Fenster
- Generierte Inhalte
- Tests

**Beispiel:**

```go
HTML: `
<!DOCTYPE html>
<html>
<head><title>Simple Window</title></head>
<body><h1>Hello from Wails!</h1></body>
</html>
`,
```

### Assets (nur auf Anwendungsebene)

Die Asset-Konfiguration ist **kein** Feld von `WebviewWindowOptions`. Die Anwendung selbst stellt Frontend-Assets über `application.Options.Assets` (`AssetOptions`) bereit; jedes Fenster übernimmt diesen Asset-Server.

```go
app := application.New(application.Options{
    Assets: application.AssetOptions{
        Handler: application.AssetFileServerFS(assets),
    },
})
```

**Weitere Informationen finden Sie unter [Build-System](/concepts/build-system/).**

### UseApplicationMenu

**Typ:** `bool` **Standard:** `false` **Plattform:** Windows, Linux (keine Auswirkung unter macOS)

```go
UseApplicationMenu: true,
```

**Zweck:** Das über `app.Menu.Set()` festgelegte Anwendungsmenü für dieses Fenster verwenden.

Unter **macOS** hat diese Option keine Auswirkung, da macOS immer ein globales Anwendungsmenü am oberen Bildschirmrand verwendet.

Unter **Windows** und **Linux** zeigen Fenster standardmäßig kein Menü an. Durch das Festlegen von `UseApplicationMenu: true` verwendet das Fenster das anwendungsweite Menü. Dies bietet eine einfache plattformübergreifende Lösung.

**Beispiel:**

```go
// Set the application menu once
menu := app.NewMenu()
menu.AddRole(application.FileMenu)
menu.AddRole(application.EditMenu)
app.Menu.Set(menu)

// All windows with UseApplicationMenu will display this menu
app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:              "Main Window",
    UseApplicationMenu: true,
})

app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:              "Second Window",
    UseApplicationMenu: true,  // Also gets the app menu
})
```

**Hinweise:**

- Wenn sowohl `UseApplicationMenu` als auch ein fensterspezifisches Menü festgelegt sind, hat das fensterspezifische Menü Vorrang
- Dies vereinfacht plattformübergreifenden Code, da keine Betriebssystemprüfung zur Laufzeit erforderlich ist
- Die vollständige Menüdokumentation finden Sie unter [Anwendungsmenüs](/features/menus/application/)

## Eingabeoptionen

### EnableFileDrop

**Typ:** `bool` **Standardwert:** `false` **Plattform:** Alle

```go
EnableFileDrop: true,
```

**Zweck:** Ermöglicht das Ziehen und Ablegen von Dateien aus dem Betriebssystem in das Fenster.

Wenn aktiviert:

- Aus Dateimanagern gezogene Dateien können in Ihrer Anwendung abgelegt werden
- Das Ereignis `WindowFilesDropped` wird mit den Pfaden der abgelegten Dateien ausgelöst
- Elemente mit dem Attribut `data-file-drop-target` stellen detaillierte Informationen zum Ablegevorgang bereit

**Anwendungsfälle:**

- Oberflächen zum Hochladen von Dateien
- Dokumenteditoren
- Tools zum Importieren von Mediendateien
- Alle Anwendungen, die Dateien akzeptieren

**Beispiel:**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:          "File Uploader",
    EnableFileDrop: true,
})

// Handle dropped files
window.OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {
    files := event.Context().DroppedFiles()
    details := event.Context().DropTargetDetails()
    
    for _, file := range files {
        fmt.Println("Dropped:", file)
    }
})
```

**HTML-Ablagebereiche:**

```html
<!-- Mark elements as drop targets -->
<div id="upload" data-file-drop-target>
    Drop files here
</div>
```

**Die vollständige Dokumentation finden Sie unter [Dateiablage](/features/drag-and-drop/files/).**

## Sicherheitsoptionen

### ContentProtectionEnabled

**Typ:** `bool` **Standardwert:** `false` **Plattform:** Windows (10+), macOS

```go
ContentProtectionEnabled: true,
```

**Zweck:** Verhindert Bildschirmaufnahmen des Fensterinhalts.

**Plattformunterstützung:**

- **Windows:** Windows 10 ab Build 19041 (vollständig), ältere Versionen (teilweise)
- **macOS:** Vollständig unterstützt
- **Linux:** Nicht unterstützt

**Anwendungsfälle:**

- Banking-Anwendungen
- Passwortmanager
- Patientenakten
- Vertrauliche Dokumente

**Wichtige Hinweise:**

1. Verhindert nicht das Abfotografieren des Bildschirms
2. Einige Tools können den Schutz möglicherweise umgehen
3. Teil eines umfassenden Sicherheitskonzepts, kein alleiniger Schutz
4. DevTools-Fenster werden nicht automatisch geschützt

**Beispiel:**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Secure Window",
    ContentProtectionEnabled: true,
})

// Toggle at runtime
window.SetContentProtection(true)
```

### Berechtigungen

**Typ:** `map[PermissionType]Permission` **Standardwert:** `nil` (plattformseitige Standardbehandlung) **Plattform:** Linux, Windows (macOS überlässt dies TCC)

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionAllow,
    application.PermissionCamera:     application.PermissionAllow,
},
```

**Zweck:** Steuert deklarativ und ohne plattformspezifischen Code, wie Berechtigungsanfragen des Webinhalts im Fenster für Funktionen wie Kamera, Mikrofon, Standortbestimmung, Benachrichtigungen und das Lesen der Zwischenablage behandelt werden.

**PermissionType-Werte:** `PermissionMicrophone`, `PermissionCamera`, `PermissionGeolocation`, `PermissionNotifications`, `PermissionClipboardRead`

**Berechtigungswerte:**

- `PermissionDefault` (0) — plattformeigene Behandlung: Aufforderung durch das Betriebssystem bzw. WebView2 unter macOS/Windows; unter Linux werden Kamera und Mikrofon zugelassen, alles andere wird verweigert
- `PermissionAllow` (1) — ohne Aufforderung gewähren (Linux: Nur Kamera und Mikrofon sind implementiert; andere Typen bleiben verweigert)
- `PermissionDeny` (2) — ohne Aufforderung verweigern

**Wichtig — Windows:** Bevor diese Option verfügbar war, gewährte Wails stillschweigend alle WebView2-Berechtigungen. Sobald Sie nun einen Eintrag in `Permissions` festlegen, wird diese pauschale Gewährung deaktiviert. Für nicht aufgeführte Berechtigungen zeigt WebView2 seine native Aufforderung an, anstatt sie automatisch zu gewähren. Führen Sie jede von Ihrer Anwendung benötigte Berechtigung ausdrücklich auf.

**Den vollständigen Leitfaden, die Plattformmatrix und Beispiele finden Sie unter [Berechtigungen](/features/windows/permissions/).**

## Fensterlebenszyklus-Ereignisse

Fensterlebenszyklus-Ereignisse werden mit `OnWindowEvent` und `RegisterHook` behandelt. Diese Methoden ermöglichen eine präzise Steuerung des Schließ- und Zerstörungsverhaltens von Fenstern.

### Schließen eines Fensters abbrechen

Um das Schließen eines Fensters zu verhindern, etwa bei nicht gespeicherten Änderungen, verwenden Sie `RegisterHook` mit dem Ereignis `WindowClosing` und rufen `event.Cancel()` auf:

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:  "main-window",
    Title: "My Application",
})

// Register a hook to intercept the closing event
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

**Wichtige Punkte:**

- `RegisterHook` fängt Ereignisse ab, bevor sie eintreten
- Rufen Sie `event.Cancel()` auf, um das Schließen des Fensters zu verhindern
- Nach dem Abbrechen bleibt das Fenster geöffnet

### Schließen des Fensters behandeln

Um beim Schließen eines Fensters Aufräumarbeiten auszuführen, verwenden Sie `OnWindowEvent` mit dem Ereignis `WindowClosing`:

```go
window.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
    // Cleanup code runs here
    fmt.Printf("Window %s is closing\n", window.Name())

    // Close database connection
    if db != nil {
        db.Close()
    }

    // Remove from window list
    removeWindow(window.ID())
})
```

**Wichtige Punkte:**

- `OnWindowEvent` behandelt Ereignisse, die unmittelbar bevorstehen
- Die Aufräumarbeiten werden ausgeführt, bevor das Fenster zerstört wird
- Das Schließen kann hier nicht abgebrochen werden (verwenden Sie dafür `RegisterHook`)

### Bereinigungsmuster für Singleton-Fenster

Verwenden Sie bei Singleton-Fenstern (um nur eine Instanz zuzulassen) `WindowClosing`, um die Referenz zu bereinigen:

```go
var settingsWindow *application.WebviewWindow

func ShowSettings(app *application.App) {
    // Create if doesn't exist
    if settingsWindow == nil {
        settingsWindow = app.Window.NewWithOptions(application.WebviewWindowOptions{
            Name:  "settings",
            Title: "Settings",
            Width: 600,
            Height: 400,
        })

        // Cleanup on close
        settingsWindow.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
            settingsWindow = nil
        })
    }

    // Show and focus
    settingsWindow.Show()
    settingsWindow.Focus()
}
```

## Plattformspezifische Optionen

### macOS-Optionen

```go
Mac: application.MacWindow{
    TitleBar: application.MacTitleBar{
        AppearsTransparent: true,
        Hide:               false,
        HideTitle:          true,
        FullSizeContent:    true,
    },
    Backdrop:                application.MacBackdropTranslucent,
    InvisibleTitleBarHeight: 50,
    WindowClass:             application.MacWindowClassPanel,
    PanelPreferences: application.MacPanelPreferences{
        NonActivating:          true,
        FloatingPanel:          true,
        BecomesKeyOnlyIfNeeded: false,
        UtilityWindow:          false,
    },
    WindowLevel:             application.MacWindowLevelFloating,
    CollectionBehavior:      application.MacWindowCollectionBehaviorDefault,
    TabbingMode:             application.MacWindowTabbingModeDisallowed,
},
```

**TitleBar** (`MacTitleBar`)

- `AppearsTransparent` – Macht die Titelleiste transparent; der Inhalt erstreckt sich in den Bereich der Titelleiste
- `Hide` – Blendet die Titelleiste vollständig aus
- `HideTitle` – Blendet nur den Titeltext aus
- `FullSizeContent` – Erweitert den Inhalt auf die gesamte Fenstergröße

Unter macOS erfordern die Transparenz der Webview und das programmgesteuerte Öffnen des Inspektors das Build-Tag `private_mac_apis`. Ohne dieses Tag bleiben dieselben Optionen gültig, aber Operationen, die private APIs verwenden, haben keine Wirkung. Auch die Liquid-Glass-Gruppierung wird ignoriert, und für Stile werden öffentliche Alternativen verwendet. Build-Befehle und das genaue Verhalten finden Sie unter [Private macOS-APIs](/guides/build/private-macos-apis/).

**Backdrop** (`MacBackdrop`)

- `MacBackdropNormal` – Undurchsichtiger Standardhintergrund
- `MacBackdropTranslucent` – **Für die Transparenz der Webview ist eine private API erforderlich.** Ohne das Tag bleibt die native Unschärfe hinter einer undurchsichtigen Webview.
- `MacBackdropTransparent` – **Für die Transparenz der Webview ist eine private API erforderlich.** Ohne das Tag bleibt die Webview undurchsichtig.
- `MacBackdropLiquidGlass` – **Für die Transparenz der Webview ist eine private API erforderlich.** Ohne das Tag bleibt die Glasschicht hinter einer undurchsichtigen Webview; für Stile werden öffentliche Alternativen verwendet.

**LiquidGlass** (`MacLiquidGlass`)

| Feld oder Wert | Abhängigkeit von privaten APIs unter macOS |
| --- | --- |
| `Style: LiquidGlassStyleAutomatic` | Der native reguläre Stil ist öffentlich; für die Transparenz der Backdrop-Webview ist `private_mac_apis` erforderlich. |
| `Style: LiquidGlassStyleLight` | Das Tag behält die vorhandene native Zuordnung zum transparenten Stil bei; ohne das Tag verwendet Wails reguläres Glas mit Aqua-Darstellung. |
| `Style: LiquidGlassStyleDark` | **Private API:** undokumentierter nativer Stilwert `2`; ohne das Tag reguläres Glas mit Dark-Aqua-Darstellung. |
| `Style: LiquidGlassStyleVibrant` | Die native Zuordnung zum transparenten Stil ist öffentlich; für die Transparenz der Backdrop-Webview ist das Tag erforderlich. |
| `GroupID` | **Private API:** Nicht leere Werte fordern eine Gruppierung an; ohne das Tag werden sie ignoriert. |
| `GroupSpacing` | **Private API:** Positive Werte fordern einen Gruppenabstand an; ohne das Tag werden sie ignoriert. |
| `Material`, `CornerRadius`, `TintColor` | Sie besitzen selbst keine Abhängigkeit von privaten APIs. |

Die nativen Stilwerte und ihre Betriebssystemverfügbarkeit finden Sie unter [Liquid-Glass-Werte](/guides/build/private-macos-apis/#liquid-glass-values).

**InvisibleTitleBarHeight** (`int`)

- Höhe des unsichtbaren Titelleistenbereichs (zum Ziehen)
- Wird nur wirksam, wenn der native Ziehbereich der Titelleiste ausgeblendet ist – also wenn das Fenster rahmenlos ist (`Frameless: true`) oder eine transparente Titelleiste verwendet (`AppearsTransparent: true`)
- Hat keine Wirkung auf Standardfenster mit sichtbarer Titelleiste

**WindowClass** (`MacWindowClass`)

- `MacWindowClassWindow` – Standardverhalten von `NSWindow` (Voreinstellung)
- `MacWindowClassPanel` – Ein zusätzliches `NSPanel`, das niemals zum Hauptfenster der Anwendung wird

`PanelPreferences` gilt nur für `MacWindowClassPanel`:

- `NonActivating` fügt `NSWindowStyleMaskNonactivatingPanel` hinzu. Durch Anzeigen oder Fokussieren des Panels wird die Wails-Anwendung nicht aktiviert. Das Panel kann jedoch weiterhin den Key-Status für Steuerelemente und Texteingaben erhalten.
- `FloatingPanel` aktiviert das Verhalten eines schwebenden Panels von AppKit.
- `BecomesKeyOnlyIfNeeded` erhält den Key-Status nur, wenn die angeklickte Ansicht eine Tastatureingabe anfordert.
- `UtilityWindow` wendet den nativen Stil eines Dienstprogrammfensters an.

Wails-Panels bleiben sichtbar, wenn die Anwendung deaktiviert wird, und werden beim Schließen freigegeben. Dies entspricht dem von `WebviewWindow` erwarteten Lebenszyklus. Dabei handelt es sich um beabsichtigte Überschreibungen der gegenteiligen Standardwerte von `NSPanel`.

Fensterklasse, Fensterebene, Aktivierungsrichtlinie und Collection-Verhalten lösen jeweils unterschiedliche Probleme:

- `WindowClass` wählt `NSWindow` oder `NSPanel` aus und steuert die Semantik von Haupt- und Key-Fenstern.
- `WindowLevel` steuert die Z-Reihenfolge. Verwenden Sie `MacWindowLevelPopUpMenu` für Overlays der Menüleiste.
- `MacOptions.ActivationPolicy` steuert die gesamte Anwendung einschließlich ihrer Darstellung im Dock und in der Menüleiste. Ein nicht aktivierendes Panel erfordert keine Accessory-Aktivierungsrichtlinie. Eine Anwendung kann diese jedoch weiterhin verwenden, um ihr Dock-Symbol auszublenden.
- `CollectionBehavior` steuert die Teilnahme an Spaces und am Vollbildmodus.

```go
// Spotlight/menu-bar panel that leaves the current application active.
Mac: application.MacWindow{
    WindowClass: application.MacWindowClassPanel,
    PanelPreferences: application.MacPanelPreferences{
        NonActivating: true,
    },
    WindowLevel: application.MacWindowLevelPopUpMenu,
    CollectionBehavior: application.MacWindowCollectionBehaviorCanJoinAllSpaces |
        application.MacWindowCollectionBehaviorFullScreenAuxiliary |
        application.MacWindowCollectionBehaviorStationary,
},
```

**WindowLevel** (`MacWindowLevel`)

- `MacWindowLevelNormal` – Standardfensterebene (Standardwert)
- `MacWindowLevelFloating` – Schwebt über normalen Fenstern
- `MacWindowLevelTornOffMenu` – Ebene für abgetrennte Menüs
- `MacWindowLevelModalPanel` – Ebene für modale Panels
- `MacWindowLevelMainMenu` – Ebene des Hauptmenüs
- `MacWindowLevelStatus` – Ebene für Statusfenster
- `MacWindowLevelPopUpMenu` – Ebene für Pop-up-Menüs
- `MacWindowLevelScreenSaver` – Ebene des Bildschirmschoners

Ein explizites `WindowLevel` hat Vorrang vor `AlwaysOnTop` und `PanelPreferences.FloatingPanel`. Ohne explizite Ebene werden `AlwaysOnTop` oder ein schwebendes Panel zu `MacWindowLevelFloating` aufgelöst; andernfalls lautet die Ebene `MacWindowLevelNormal`. Ein späterer Aufruf von `SetAlwaysOnTop` bleibt eine explizite Änderung zur Laufzeit.

**CollectionBehavior** (`MacWindowCollectionBehavior`)

Steuert das Verhalten des Fensters in macOS Spaces und im Vollbildmodus. Diese Werte sind Bitmasken und können mit einem bitweisen ODER (`|`) kombiniert werden.

**Space-Verhalten:**

- `MacWindowCollectionBehaviorDefault` – Verwendet FullScreenPrimary (Standardwert, abwärtskompatibel)
- `MacWindowCollectionBehaviorCanJoinAllSpaces` – Das Fenster erscheint auf allen Spaces
- `MacWindowCollectionBehaviorMoveToActiveSpace` – Wechselt beim Anzeigen zum aktiven Space
- `MacWindowCollectionBehaviorManaged` – Standardverhalten für verwaltete Fenster
- `MacWindowCollectionBehaviorTransient` – Temporäres/transientes Fenster
- `MacWindowCollectionBehaviorStationary` – Bleibt beim Wechsel zwischen Spaces an seiner Position

**Fensterwechsel:**

- `MacWindowCollectionBehaviorParticipatesInCycle` – In den Fensterwechsel mit Cmd+` einbezogen
- `MacWindowCollectionBehaviorIgnoresCycle` – Vom Fensterwechsel mit Cmd+` ausgeschlossen

**Vollbildverhalten:**

- `MacWindowCollectionBehaviorFullScreenPrimary` – Kann in den Vollbildmodus wechseln
- `MacWindowCollectionBehaviorFullScreenAuxiliary` – Kann Anwendungen im Vollbildmodus überlagern
- `MacWindowCollectionBehaviorFullScreenNone` – Deaktiviert die Vollbildfunktion
- `MacWindowCollectionBehaviorFullScreenAllowsTiling` – Ermöglicht die Anordnung nebeneinander (macOS 10.11+)
- `MacWindowCollectionBehaviorFullScreenDisallowsTiling` – Verhindert die Anordnung nebeneinander (macOS 10.11+)

**Beispiel – Spotlight-ähnliches Fenster:**

```go
// Window that appears on all Spaces AND can overlay fullscreen apps
Mac: application.MacWindow{
	WindowClass: application.MacWindowClassPanel,
	PanelPreferences: application.MacPanelPreferences{
		NonActivating: true,
	},
    CollectionBehavior: application.MacWindowCollectionBehaviorCanJoinAllSpaces |
                        application.MacWindowCollectionBehaviorFullScreenAuxiliary,
    WindowLevel:        application.MacWindowLevelFloating,
},
```

**Beispiel – Einzelnes Verhalten:**

```go
// Window that can appear over fullscreen applications
Mac: application.MacWindow{
    CollectionBehavior: application.MacWindowCollectionBehaviorFullScreenAuxiliary,
},
```

**TabbingMode** (`MacWindowTabbingMode`)

Steuert das Verhalten von Fenster-Tabs unter macOS 10.12 und neuer. Mit Fenster-Tabs können mehrere Fenster als Tabs gruppiert werden.

**Optionen:**

- `MacWindowTabbingModeDefault` – Sentinel-Nullwert (nicht explizit festgelegt). Zur Laufzeit sind Fenster-Tabs standardmäßig nicht zulässig
- `MacWindowTabbingModeAutomatic` – Das System bestimmt das Verhalten von Fenster-Tabs
- `MacWindowTabbingModePreferred` – Das Fenster bevorzugt den Tab-Modus
- `MacWindowTabbingModeDisallowed` – Deaktiviert Fenster-Tabs

**Beispiel – Fenster-Tabs deaktivieren:**

```go
Mac: application.MacWindow{
    TabbingMode: application.MacWindowTabbingModeDisallowed,
},
```

**Beispiel – Fenster-Tabs bevorzugen:**

```go
Mac: application.MacWindow{
    TabbingMode: application.MacWindowTabbingModePreferred,
},
```

**WebviewPreferences** (`MacWebviewPreferences`)

Ermöglicht die detaillierte Steuerung der zugrunde liegenden `WKWebView`-Konfiguration. Alle Felder sind optional – nicht festgelegte Felder ändern den WebKit-Standardwert nicht.

```go
Mac: application.MacWindow{
    WebviewPreferences: application.MacWebviewPreferences{
        TabFocusesLinks:                       optional.True,
        TextInteractionEnabled:                optional.True,
        FullscreenEnabled:                     optional.False,
        AllowsBackForwardNavigationGestures:   optional.False,
        AllowsMagnification:                   optional.True,
        AllowsAirPlayForMediaPlayback:         optional.True,
        JavaScriptCanOpenWindowsAutomatically: optional.False,
        MinimumFontSize:                       optional.NewVar(12.0),
        ApplicationNameForUserAgent:           "MyApp",
        EnableAutoplayWithoutUserAction:       optional.True,
    },
},
```

- `TabFocusesLinks` — Wenn `true`, verschiebt die Tabulatortaste den Fokus auf Links und Formularsteuerelemente (Standard: `false`)
- `TextInteractionEnabled` — Wenn `true`, können Benutzer Text in der Webview auswählen und mit ihm interagieren (Standard: `true`)
- `FullscreenEnabled` — Wenn `true`, können Webinhalte über die HTML Fullscreen API in den Vollbildmodus wechseln (Standard: `false`). Erfordert macOS 12.3 oder neuer.
- `AllowsBackForwardNavigationGestures` — Wenn `true`, lösen horizontale Wischgesten die Rückwärts-/Vorwärtsnavigation aus (Standard: `false`)
- `AllowsMagnification` — Wenn `true`, ist Pinch-to-Zoom in der Webview aktiviert (Standard: `false`)
- `AllowsAirPlayForMediaPlayback` — Wenn `true`, können Medien an AirPlay-Geräte gestreamt werden (Standard: `true`)
- `JavaScriptCanOpenWindowsAutomatically` — Wenn `true`, kann JavaScript ohne Benutzergeste neue Fenster öffnen (Standard: `false`)
- `MinimumFontSize` — Mindestschriftgröße in Punkt. Legen Sie sie mit `optional.NewVar(12.0)` fest. Wenn sie nicht festgelegt ist, bleibt der WebKit-Standard erhalten.
- `ApplicationNameForUserAgent` — Überschreibt das Suffix des Anwendungsnamens in der User-Agent-Zeichenfolge von WebKit. Nützlich, wenn Websites die standardmäßige Kennung `"wails.io"` ablehnen (z. B. YouTube-Einbettungen). Lassen Sie den Wert leer, um den Standard beizubehalten.
- `EnableAutoplayWithoutUserAction` — Wenn `true`, können Audio und Video ohne Benutzergeste automatisch wiedergegeben werden. Wird auf `WKWebViewConfiguration.mediaTypesRequiringUserActionForPlayback = WKAudiovisualMediaTypeNone` abgebildet (Standard: `false`)

### Windows-Optionen (pro Fenster)

Die Struktur pro Fenster ist `application.WindowsWindow` — **nicht** `WindowsOptions` (dies ist die Struktur auf <em>Anwendungs</em>ebene).

```go
Windows: application.WindowsWindow{
    DisableIcon:                       false,
    DisableMenu:                       false,
    BackdropType:                      application.Auto,
    CustomTheme:                       application.ThemeSettings{},
    DisableFramelessWindowDecorations: false,
    NonClientRegionSupport:            false,
    WebView2CompositionHosting:        false,
},
```

**DisableIcon** (`bool`)

- Entfernt das Symbol aus der Titelleiste.

**DisableMenu** (`bool`)

- Deaktiviert die Menüleiste für das Fenster. Wenn `true`, zeigt das Fenster auch dann keine Menüleiste an, wenn eine konfiguriert ist.
- Standard: `false`

**BackdropType** (`BackdropType`)

- `application.Auto` – Systemstandard
- `application.None` – Kein Hintergrundeffekt
- `application.Mica` – Mica-Material (Windows 11)
- `application.Acrylic` – Acrylmaterial (Windows 11)
- `application.Tabbed` – Material mit Registerkartenoptik (Windows 11)

Es gibt keine Konstanten im Stil von `WindowsBackdropTypeMica` — verwenden Sie `application.Mica` usw.

**CustomTheme** (`ThemeSettings`)

- Wert (kein Zeiger). Benutzerdefinierte Farben für den Dunkel-/Hellmodus für Fensterrahmen, Text und Hintergrund der Titelleiste sowie die Menüleiste.

**DisableFramelessWindowDecorations** (`bool`)

- Deaktiviert die standardmäßigen Dekorationen rahmenloser Fenster (Aero-Schatten, abgerundete Ecken).

**NonClientRegionSupport** (`bool`)

- Aktiviert die native WebView2-Unterstützung für `app-region: drag` / `app-region: no-drag` bei rahmenlosen benutzerdefinierten Titelleisten.
- Dies dient nur zum einfachen nativen Ziehen der Anwendung. Es bietet weder natives Verhalten für benutzerdefinierte Titelleistenschaltflächen noch Windows 11 Snap Assist / Snap Layouts für benutzerdefinierte Maximierungsschaltflächen.

**WebView2CompositionHosting** (`bool`)

- Aktiviert die von Wails verwaltete Unterstützung für `--wails-non-client-region` bei benutzerdefinierten Titelleistenschaltflächen mit nativem Windows-Verhalten, einschließlich Windows 11 Snap Assist / Snap Layouts bei benutzerdefinierten Maximierungsschaltflächen.
- Experimentell. Dadurch wird WebView2 über `ICoreWebView2CompositionController` und DirectComposition statt über den standardmäßigen HWND-gehosteten Controller gehostet.
- Kann mit `NonClientRegionSupport` kombiniert werden, wenn ein Fenster sowohl native WebView2-Unterstützung für `app-region` als auch von Wails verwaltete benutzerdefinierte Bereiche für Titelleistenschaltflächen benötigt.

**Beispiel:**

```go
Windows: application.WindowsWindow{
    BackdropType: application.Mica,
    DisableIcon:  true,
},
```

**Beispiel – Benutzerdefinierte Bereiche der Windows-Titelleiste:**

```go
Windows: application.WindowsWindow{
    NonClientRegionSupport:    true,
    WebView2CompositionHosting: true,
},
```

Ausführliche Informationen zum Verhalten, zu den Kompromissen und zum passenden CSS finden Sie unter [Rahmenlose Fenster](/features/windows/frameless/#native-non-client-regions-on-windows).

### Linux-Optionen (pro Fenster)

Die Struktur pro Fenster ist `application.LinuxWindow` — **nicht** `LinuxOptions`.

```go
Linux: application.LinuxWindow{
    Icon:                []byte{/* PNG data */},
    WindowIsTranslucent: false,
},
```

**Icon** (`[]byte`)

- Fenstersymbol (PNG-Format).

**WindowIsTranslucent** (`bool`)

- Erfordert Unterstützung durch den Compositor.

**Beispiel:**

```go
//go:embed icon.png
var icon []byte

Linux: application.LinuxWindow{
    Icon: icon,
},
```

## Windows-Optionen auf Anwendungsebene

Einige Windows-spezifische Optionen müssen auf Anwendungsebene statt für einzelne Fenster konfiguriert werden. Der Grund dafür ist, dass WebView2 für jeden Benutzerdatenpfad eine gemeinsame Browserumgebung verwendet.

### Browser-Flags

WebView2-Browser-Flags steuern experimentelle Funktionen und das Verhalten für **alle Fenster** Ihrer Anwendung. Sie müssen in `application.Options.Windows` festgelegt werden:

```go
app := application.New(application.Options{
    Name: "My App",
    Windows: application.WindowsOptions{
        // Enable experimental WebView2 features
        EnabledFeatures: []string{
            "msWebView2EnableDraggableRegions",
        },

        // Disable specific features
        DisabledFeatures: []string{
            "msSmartScreenProtection",  // Always disabled by Wails
        },

        // Additional Chromium command-line arguments
        AdditionalBrowserArgs: []string{
            "--disable-gpu",
            "--remote-debugging-port=9222",
        },
    },
})
```

**EnabledFeatures** (`[]string`)

- Liste der zu aktivierenden WebView2-Feature-Flags
- Verfügbare Flags finden Sie unter [WebView2-Browser-Flags](https://learn.microsoft.com/en-us/microsoft-edge/webview2/concepts/webview-features-flags)
- Beispiel: `"msWebView2EnableDraggableRegions"`

**DisabledFeatures** (`[]string`)

- Liste der zu deaktivierenden WebView2-Feature-Flags
- Wails deaktiviert automatisch `msSmartScreenProtection`
- Beispiel: `"msExperimentalFeature"`

**AdditionalBrowserArgs** (`[]string`)

- Chromium-Befehlszeilenargumente, die an den Browserprozess übergeben werden
- Muss das Präfix `--` enthalten (z. B. `"--remote-debugging-port=9222"`)
- Verfügbare Argumente finden Sie unter [Chromium-Befehlszeilenschalter](https://peter.sh/experiments/chromium-command-line-switches/)

@note{type="caution" title="Wichtig"}
Diese Flags gelten global für ALLE Fenster, da WebView2 für jeden Benutzerdatenpfad eine gemeinsame Browserumgebung verwendet. Für verschiedene Fenster können keine unterschiedlichen Browser-Flags festgelegt werden.

@end

**Vollständiges Beispiel:**

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "My App",
        Windows: application.WindowsOptions{
            // Enable draggable regions feature
            EnabledFeatures: []string{
                "msWebView2EnableDraggableRegions",
            },
            // Enable remote debugging
            AdditionalBrowserArgs: []string{
                "--remote-debugging-port=9222",
            },
        },
    })

    // All windows will use the browser flags configured above
    window := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:  "Main Window",
        Width:  1024,
        Height: 768,
    })

    window.Show()
    app.Run()
}
```

## Vollständiges Beispiel

Hier sehen Sie eine produktionsreife Fensterkonfiguration:

```go
package main

import (
    _ "embed"
    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var assets embed.FS

//go:embed icon.png
var icon []byte

func main() {
    app := application.New(application.Options{
        Name: "My Application",
    })

    window := app.Window.NewWithOptions(application.WebviewWindowOptions{
        // Identity
        Name:  "main-window",
        Title: "My Application",

        // Size and Position
        Width:     1200,
        Height:    800,
        MinWidth:  800,
        MinHeight: 600,

        // Initial State
        StartState: application.WindowStateNormal,

        // Appearance
        BackgroundColour: application.NewRGB(255, 255, 255),

        // Platform-Specific (per-window structs)
        Mac: application.MacWindow{
            TitleBar: application.MacTitleBar{
                AppearsTransparent: true,
            },
            Backdrop: application.MacBackdropTranslucent,
        },

        Windows: application.WindowsWindow{
            BackdropType: application.Mica,
            DisableIcon:  false,
        },

        Linux: application.LinuxWindow{
            Icon: icon,
        },
    })

    window.Center()
    window.Show()

    app.Run()
}
```

Frontend-Assets werden auf Anwendungsebene (`application.Options.Assets`) und nicht für einzelne Fenster bereitgestellt.

## Nächste Schritte

- [Fenstergrundlagen](/features/windows/basics/) – Fenster erstellen und steuern
- [Mehrere Fenster](/features/windows/multiple/) – Muster für Anwendungen mit mehreren Fenstern
- [Rahmenlose Fenster](/features/windows/frameless/) – Benutzerdefinierte Fensterdekoration
- [Fensterereignisse](/features/windows/events/) – Lebenszyklusereignisse

---

**Fragen?** Stellen Sie sie auf [Discord](https://discord.gg/JDdSxwjhGf) oder sehen Sie sich die [Beispiele](https://github.com/wailsapp/wails/tree/master/v3/examples) an.
