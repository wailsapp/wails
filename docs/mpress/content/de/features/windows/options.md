---
title: "Fensteroptionen"
description: "Fensteroptionen, Standardwerte und Plattformbeschränkungen"
slug: "features/windows/options"
sourcePath: "features/windows/options.md"
---

## Optionen zur Fensterkonfiguration

Wails bietet eine umfassende Fensterkonfiguration mit Dutzenden Optionen für Größe, Position, Darstellung und Verhalten. Diese Referenz erläutert gängige Felder von `WebviewWindowOptions` unter Windows, macOS und Linux. Die folgende Struktur fasst ausgewählte Felder zusammen; den vollständigen Typ finden Sie im verlinkten Go-Quellcode.

## Struktur WebviewWindowOptions

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
    KeyBindings                 map[string]func(window Window)

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

`WebviewWindowOptions` verfügt über kein Feld `Parent` – verwenden Sie für übergeordnete/modale Beziehungen `parentWindow.AttachModal(childWindow)`. Ebenso gibt es kein Feld `Assets` – die Asset-Konfiguration erfolgt über `application.Options` (`Assets AssetOptions`).

Vollständiger Quellcode: [`v3/pkg/application/webview_window_options.go`](https://github.com/wailsapp/wails/blob/releases/v3.0.0/v3/pkg/application/webview_window_options.go).

## Grundlegende Optionen

### Name

**Typ:** `string` **Standardwert:** `window-<numeric ID>` **Plattform:** Alle

```go
Name: "main-window"
```

**Zweck:** Eindeutige Kennung, um Fenster später zu finden.

**Bewährte Vorgehensweisen:**

- Verwenden Sie aussagekräftige Namen: `"main"`, `"settings"`, `"about"`
- Verwenden Sie Kebab-Case: `"file-browser"`, `"color-picker"`
- Wählen Sie einen kurzen und einprägsamen Namen

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

### Titel

**Typ:** `string` **Standardwert:** Anwendungsname **Plattform:** Alle

```go
Title: "My Application"
```

**Zweck:** Text, der in der Titelleiste und Taskleiste angezeigt wird.

**Dynamische Aktualisierungen:**

```go
window.SetTitle("My Application - Document.txt")
```

### Breite/Höhe

**Typ:** `int` (Pixel) **Standardwert:** 800 × 600 **Plattform:** Alle **Beschränkungen:** Muss positiv sein

```go
Width:  1200,
Height: 800,
```

**Zweck:** Anfangsgröße des Fensters in logischen Pixeln.

**Hinweise:**

- Wails berücksichtigt die DPI-Skalierung automatisch
- Verwenden Sie logische statt physischer Pixel
- Berücksichtigen Sie die minimale Bildschirmauflösung (1024 × 768)

**Beispielgrößen:**

| Anwendungsfall | Breite | Höhe |
| --- | --- | --- |
| Kleines Dienstprogramm | 400 | 300 |
| Standardanwendung | 1024 | 768 |
| Große Anwendung | 1440 | 900 |
| Full HD | 1920 | 1080 |

### X/Y

**Typ:** `int` (Pixel) **Standardwert:** Auf dem Bildschirm zentriert **Plattform:** Alle

```go
X: 100,  // 100px from left edge
Y: 100,  // 100px from top edge
```

**Zweck:** Anfangsposition des Fensters.

**Koordinatensystem:**

- (0, 0) ist die linke obere Ecke des primären Bildschirms
- Positive X-Werte verlaufen nach rechts
- Positive Y-Werte verlaufen nach unten

**Beispiel:**

`X` und `Y` werden nur wirksam, wenn `InitialPosition: application.WindowXY` festgelegt ist. Andernfalls verwendet `InitialPosition` standardmäßig `WindowCentered`, und `X`/`Y` werden ignoriert.

```go
settings := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:            "coordinate-window",
    InitialPosition: application.WindowXY, // opt into X/Y coordinates
    X:               100,
    Y:               100,
})
```

**Bewährte Vorgehensweise:** Wenn bestimmte Koordinaten unwichtig sind, zentrieren Sie das Fenster nach seiner Erstellung mit `Center()`:

```go
window := app.Window.New()
window.Center()
```

### MinWidth/MinHeight

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

### MaxWidth/MaxHeight

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

## Zustandsoptionen

### Ausgeblendet

**Typ:** `bool` **Standardwert:** `false` **Plattform:** Alle

```go
Hidden: true,
```

**Zweck:** Erstellt das Fenster, ohne es anzuzeigen.

**Anwendungsfälle:**

- Hintergrundfenster
- Bei Bedarf angezeigte Fenster
- Startbildschirme (erstellen, laden und anschließend anzeigen)
- Weißes Aufblitzen beim Laden von Inhalten verhindern

**Plattformspezifische Verbesserungen:**

- **Windows:** Weißes Aufblitzen des Fensters behoben – das Fenster bleibt unsichtbar, bis `Show()` aufgerufen wird
- **macOS:** Vollständig unterstützt
- **Linux:** Vollständig unterstützt

**Empfohlenes Muster für ein reibungsloses Laden:**

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

### Rahmenlos

**Typ:** `bool` **Standardwert:** `false` **Plattform:** Alle

```go
Frameless: true,
```

**Zweck:** Entfernt die Titelleiste und die Fensterrahmen.

**Anwendungsfälle:**

- Benutzerdefinierte Fensterdekoration
- Startbildschirme
- Kioskanwendungen
- Individuell gestaltete Fenster

**Wichtig:** Sie müssen Folgendes implementieren:

- Verschieben des Fensters
- Schaltflächen zum Schließen, Minimieren und Maximieren
- Griffe zur Größenänderung (falls die Größe veränderbar ist)

**Weitere Informationen finden Sie unter [Rahmenlose Fenster](/features/windows/frameless/).**

### DisableResize

**Typ:** `bool` **Standardwert:** `false` (die Fenstergröße ist standardmäßig veränderbar) **Plattform:** Alle

```go
DisableResize: true,
```

**Zweck:** Verhindert die Größenänderung des Fensters. Beachten Sie, dass dieses Feld die **Umkehrung** von v2s `Resizable` ist. Setzen Sie `DisableResize: true`, damit die Fenstergröße nicht verändert werden kann.

**Anwendungsfälle:**

- Anwendungen mit fester Größe
- Startbildschirme
- Dialogfenster

**Hinweis:** Benutzer können das Fenster weiterhin maximieren oder in den Vollbildmodus versetzen, sofern Sie diese Funktionen nicht ebenfalls über `MaximiseButtonState` beziehungsweise das Collection-Verhalten deaktivieren.

### AlwaysOnTop

**Typ:** `bool` **Standardwert:** `false` **Plattform:** Alle

```go
AlwaysOnTop: true,
```

**Zweck:** Hält das Fenster über allen anderen Fenstern.

**Anwendungsfälle:**

- Schwebende Symbolleisten
- Benachrichtigungen
- Bild-im-Bild
- Timer

**Plattformhinweise:**

- **macOS:** Vollständig unterstützt
- **Windows:** Vollständig unterstützt
- **Linux:** Abhängig vom Fenstermanager

### StartState

**Typ:** `WindowState`-Enum **Standardwert:** `WindowStateNormal` **Plattform:** Alle

```go
StartState: application.WindowStateMaximised,
```

**Zweck:** Anfangszustand des Fensters bei seiner Anzeige.

**Werte:**

- `WindowStateNormal` – Normales Fenster
- `WindowStateMinimised` – Minimiert
- `WindowStateMaximised` – Maximiert
- `WindowStateFullscreen` – Vollbild

Es gibt keine Konstante `WindowStateHidden`. Verwenden Sie das boolesche Feld `Hidden`, um das Fenster anfangs unsichtbar zu machen.

**Vollbildmodus zur Laufzeit umschalten:**

```go
window.Fullscreen()
window.UnFullscreen()
window.ToggleFullscreen() // there is no SetFullscreen(bool)
```

## Darstellungsoptionen

### BackgroundColour

**Typ:** `RGBA`-Struct **Standardwert:** Weiß **Plattform:** Alle

```go
BackgroundColour: application.RGBA{Red: 0, Green: 0, Blue: 0, Alpha: 255},
```

Die Felder von `RGBA` sind vom Typ `Red, Green, Blue, Alpha` (uint8). Verwenden Sie vorzugsweise die Hilfsfunktionen `application.NewRGB(r, g, b)` (Alpha 255) oder `application.NewRGBA(r, g, b, a)`.

**Zweck:** Hintergrundfarbe des Fensters, bevor der Inhalt geladen wird.

**Anwendungsfälle:**

- An das Farbschema Ihrer Anwendung anpassen
- Weißes Aufblitzen bei dunklen Farbschemata verhindern
- Flüssiger Ladevorgang

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

**Typ:** `BackgroundType`-Enum **Standardwert:** `BackgroundTypeSolid` **Plattform:** macOS, Windows (teilweise)

```go
BackgroundType: application.BackgroundTypeTranslucent,
```

**Werte:**

- `BackgroundTypeSolid` – Volltonfarbe
- `BackgroundTypeTransparent` – Vollständig transparent
- `BackgroundTypeTranslucent` – Halbtransparente Unschärfe

**Plattformunterstützung:**

- **macOS:** Konfigurieren Sie `Mac.Backdrop`. Für Webview-Transparenz ist [`private_mac_apis`](/guides/build/private-macos-apis/#webview-transparency-and-background) erforderlich. Andernfalls bleibt die Webview undurchsichtig.
- **Windows:** Transparent und halbtransparent (Windows 11+)
- **Linux:** Nur Volltonfarbe

**Beispiel (macOS):**

```go
BackgroundType: application.BackgroundTypeTranslucent,
Mac: application.MacWindow{
    Backdrop: application.MacBackdropTranslucent,
},
```

### OpenInspectorOnStartup und OpenDevTools

**Private API unter macOS:** `OpenInspectorOnStartup: true`, Go `window.OpenDevTools()` und JavaScript `Window.OpenDevTools()` benötigen `private_mac_apis`, um den Inspektor programmgesteuert zu öffnen. Ohne diese Option haben diese Operationen keine Wirkung. Produktions-Builds benötigen außerdem `devtools`. Die öffentliche Safari-Inspektion unter macOS 13.3+ benötigt keine privaten APIs; die Aktivierung des Inspektors unter älteren macOS-Versionen hingegen schon. Weitere Informationen finden Sie in der [Build-Matrix für den Web Inspector](/guides/build/private-macos-apis/#web-inspector).

## Inhaltsoptionen

### URL

**Typ:** `string` **Standardwert:** Leer (lädt aus Assets) **Plattform:** Alle

```go
URL: "https://example.com",
```

**Zweck:** Lädt anstelle eingebetteter Assets eine externe URL.

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

Codeausschnitt auf Anwendungsebene für den Produktionseinsatz:

```go
app := application.New(application.Options{
    Name: "My App",
    Assets: application.AssetOptions{
        Handler: application.AssetFileServerFS(assets),
    },
})
```

### HTML

**Typ:** `string` **Standardwert:** Leer **Plattform:** Alle

```go
HTML: "<h1>Hello World</h1>",
```

**Zweck:** Lädt direkt eine HTML-Zeichenfolge.

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

Die Asset-Konfiguration ist **kein** `WebviewWindowOptions`-Feld. Frontend-Assets werden von der Anwendung selbst über `application.Options.Assets` (`AssetOptions`) bereitgestellt; jedes Fenster übernimmt diesen Asset-Server.

```go
app := application.New(application.Options{
    Assets: application.AssetOptions{
        Handler: application.AssetFileServerFS(assets),
    },
})
```

**Weitere Informationen finden Sie unter [Build-System](/concepts/build-system/).**

### UseApplicationMenu

**Typ:** `bool` **Standardwert:** `false` **Plattform:** Windows, Linux (ohne Auswirkung unter macOS)

```go
UseApplicationMenu: true,
```

**Zweck:** Für dieses Fenster das Anwendungsmenü verwenden, das über `app.Menu.Set()` festgelegt wurde.

Unter **macOS** hat diese Option keine Auswirkung, da macOS stets ein globales Anwendungsmenü am oberen Bildschirmrand verwendet.

Unter **Windows** und **Linux** zeigen Fenster standardmäßig kein Menü an. Durch Festlegen von `UseApplicationMenu: true` verwendet das Fenster das Menü auf Anwendungsebene. Dies bietet eine einfache plattformübergreifende Lösung.

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
- Dies vereinfacht plattformübergreifenden Code, da keine Betriebssystemprüfungen zur Laufzeit erforderlich sind
- Die vollständige Dokumentation zu Menüs finden Sie unter [Anwendungsmenüs](/features/menus/application/)

## Eingabeoptionen

### EnableFileDrop

**Typ:** `bool` **Standardwert:** `false` **Plattform:** Alle

```go
EnableFileDrop: true,
```

**Zweck:** Das Ziehen und Ablegen von Dateien aus dem Betriebssystem in das Fenster aktivieren.

Wenn diese Option aktiviert ist:

- Aus Dateimanagern gezogene Dateien können in Ihrer Anwendung abgelegt werden
- Das Ereignis `WindowFilesDropped` wird mit den Pfaden der abgelegten Dateien ausgelöst
- Elemente mit dem Attribut `data-file-drop-target` stellen detaillierte Informationen zum Ablegevorgang bereit

**Anwendungsfälle:**

- Oberflächen zum Hochladen von Dateien
- Dokumenteditoren
- Medienimporte
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

**Zweck:** Bildschirmaufnahmen des Fensterinhalts verhindern.

**Plattformunterstützung:**

- **Windows:** Windows 10 Build 19041+ (vollständig), ältere Versionen (teilweise)
- **macOS:** Vollständig unterstützt
- **Linux:** Nicht unterstützt

**Anwendungsfälle:**

- Banking-Anwendungen
- Passwortmanager
- Patientenakten
- Vertrauliche Dokumente

**Wichtige Hinweise:**

1. Verhindert nicht das Abfotografieren des Bildschirms
2. Einige Werkzeuge können den Schutz möglicherweise umgehen
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

**Typ:** `map[PermissionType]Permission` **Standardwert:** `nil` (plattformseitige Standardbehandlung) **Plattform:** Linux, Windows (macOS überlässt die Behandlung TCC)

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionAllow,
    application.PermissionCamera:     application.PermissionAllow,
},
```

**Zweck:** Deklarativ und ohne plattformspezifischen Code steuern, wie Berechtigungsanfragen aus den Webinhalten des Fensters behandelt werden, etwa für Kamera, Mikrofon, Geolokalisierung, Benachrichtigungen und das Lesen der Zwischenablage.

**PermissionType-Werte:** `PermissionMicrophone`, `PermissionCamera`, `PermissionGeolocation`, `PermissionNotifications`, `PermissionClipboardRead`

**Berechtigungswerte:**

- `PermissionDefault` (0) — native Behandlung der Plattform: Aufforderung des Betriebssystems bzw. von WebView2 unter macOS/Windows; unter Linux werden Kamera und Mikrofon zugelassen und alle anderen Berechtigungen verweigert
- `PermissionAllow` (1) — ohne Nachfrage gewähren (Linux: Nur Kamera und Mikrofon sind implementiert; andere Typen bleiben verweigert)
- `PermissionDeny` (2) — ohne Nachfrage verweigern

**Wichtig — Windows:** Bevor diese Option verfügbar war, gewährte Wails stillschweigend alle WebView2-Berechtigungen. Sobald ein beliebiger Eintrag in `Permissions` festgelegt wird, ist diese pauschale Gewährung deaktiviert. Für nicht aufgeführte Berechtigungen zeigt WebView2 seine native Aufforderung an, statt sie automatisch zu gewähren. Führen Sie jede von Ihrer Anwendung benötigte Berechtigung ausdrücklich auf.

**Den vollständigen Leitfaden, die Plattformmatrix und Beispiele finden Sie unter [Berechtigungen](/features/windows/permissions/).**

## Ereignisse im Fensterlebenszyklus

Ereignisse im Fensterlebenszyklus werden mit `OnWindowEvent` und `RegisterHook` behandelt. Diese Methoden ermöglichen eine präzise Steuerung des Verhaltens beim Schließen und Zerstören eines Fensters.

### Schließen eines Fensters abbrechen

Um das Schließen eines Fensters zu verhindern, beispielsweise bei nicht gespeicherten Änderungen, verwenden Sie `RegisterHook` mit dem Ereignis `WindowClosing` und rufen `event.Cancel()` auf:

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

### Schließen eines Fensters behandeln

Um beim Schließen eines Fensters Bereinigungsarbeiten auszuführen, verwenden Sie `OnWindowEvent` mit dem Ereignis `WindowClosing`:

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
- Die Bereinigung wird ausgeführt, bevor das Fenster zerstört wird
- Das Schließen kann hier nicht abgebrochen werden (verwenden Sie dafür `RegisterHook`).

### Bereinigungsmuster für Singleton-Fenster

Verwenden Sie für Singleton-Fenster (um nur eine Instanz zuzulassen) `WindowClosing`, um die Referenz zu bereinigen:

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

**Titelleiste** (`MacTitleBar`)

- `AppearsTransparent` – Macht die Titelleiste transparent; der Inhalt erstreckt sich in den Bereich der Titelleiste
- `Hide` – Blendet die Titelleiste vollständig aus
- `HideTitle` – Blendet nur den Titeltext aus
- `FullSizeContent` – Erweitert den Inhalt auf die gesamte Fenstergröße

Unter macOS erfordern die Webview-Transparenz und das programmgesteuerte Öffnen des Inspectors das Build-Tag `private_mac_apis`. Ohne dieses Tag bleiben dieselben Optionen gültig, aber Operationen, die private APIs voraussetzen, bleiben wirkungslos. Auch die Liquid-Glass-Gruppierung wird ignoriert, und die Stile verwenden öffentliche Alternativen. Build-Befehle und das genaue Verhalten finden Sie unter [Private macOS-APIs](/guides/build/private-macos-apis/).

**Hintergrundeffekt** (`MacBackdrop`)

- `MacBackdropNormal` – Undurchsichtiger Standardhintergrund
- `MacBackdropTranslucent` – **Für Webview-Transparenz ist eine private API erforderlich.** Ohne das Tag bleibt der native Weichzeichnungseffekt hinter einer undurchsichtigen Webview.
- `MacBackdropTransparent` – **Für Webview-Transparenz ist eine private API erforderlich.** Ohne das Tag bleibt die Webview undurchsichtig.
- `MacBackdropLiquidGlass` – **Für Webview-Transparenz ist eine private API erforderlich.** Ohne das Tag bleibt die Glasschicht hinter einer undurchsichtigen Webview; die Stile verwenden öffentliche Alternativen.

**Liquid Glass** (`MacLiquidGlass`)

| Feld oder Wert | Abhängigkeit von privaten APIs unter macOS |
| --- | --- |
| `Style: LiquidGlassStyleAutomatic` | Der native reguläre Stil ist öffentlich; für die Webview-Transparenz des Hintergrundeffekts ist `private_mac_apis` erforderlich. |
| `Style: LiquidGlassStyleLight` | Das Tag erhält die bestehende native Zuordnung des transparenten Stils; ohne das Tag verwendet Wails reguläres Glas mit Aqua-Erscheinungsbild. |
| `Style: LiquidGlassStyleDark` | **Private API:** undokumentierter nativer Stilwert `2`; ohne das Tag wird reguläres Glas mit Dark-Aqua-Erscheinungsbild verwendet. |
| `Style: LiquidGlassStyleVibrant` | Die native Zuordnung des transparenten Stils ist öffentlich; für die Webview-Transparenz des Hintergrundeffekts ist das Tag erforderlich. |
| `GroupID` | **Private API:** Nicht leere Werte fordern eine Gruppierung an; ohne das Tag werden sie ignoriert. |
| `GroupSpacing` | **Private API:** Positive Werte fordern einen Gruppenabstand an; ohne das Tag werden sie ignoriert. |
| `Material`, `CornerRadius`, `TintColor` | Keine eigene Abhängigkeit von privaten APIs. |

Native Stilwerte und die Betriebssystemverfügbarkeit finden Sie unter [Liquid-Glass-Werte](/guides/build/private-macos-apis/#liquid-glass-values).

**Höhe der unsichtbaren Titelleiste** (`int`)

- Höhe des unsichtbaren Titelleistenbereichs (zum Ziehen)
- Wirkt sich nur aus, wenn der native Ziehbereich der Titelleiste ausgeblendet ist, also wenn das Fenster rahmenlos ist (`Frameless: true`) oder eine transparente Titelleiste verwendet (`AppearsTransparent: true`)
- Hat keine Auswirkung auf Standardfenster mit sichtbarer Titelleiste

**Fensterklasse** (`MacWindowClass`)

- `MacWindowClassWindow` – Standardverhalten von `NSWindow` (Voreinstellung)
- `MacWindowClassPanel` – Ein zusätzliches `NSPanel`, das nie zum Hauptfenster der Anwendung wird

`PanelPreferences` gilt nur für `MacWindowClassPanel`:

- `NonActivating` fügt `NSWindowStyleMaskNonactivatingPanel` hinzu. Das Anzeigen oder Fokussieren des Panels aktiviert die Wails-Anwendung nicht; das Panel kann jedoch weiterhin den Key-Status für Steuerelemente und Texteingaben erhalten.
- `FloatingPanel` aktiviert das Verhalten schwebender Panels von AppKit.
- `BecomesKeyOnlyIfNeeded` erhält den Key-Status nur, wenn die angeklickte Ansicht eine Tastatureingabe anfordert.
- `UtilityWindow` wendet den nativen Stil für Dienstprogrammfenster an.

Wails-Panels bleiben sichtbar, wenn die Anwendung deaktiviert wird, und werden beim Schließen freigegeben. Dies entspricht dem von `WebviewWindow` erwarteten Lebenszyklus. Diese Einstellungen überschreiben bewusst die gegenteiligen Voreinstellungen von `NSPanel`.

Fensterklasse, Fensterebene, Aktivierungsrichtlinie und Sammlungsverhalten lösen unterschiedliche Probleme:

- `WindowClass` wählt `NSWindow` oder `NSPanel` aus und steuert die Semantik von Haupt- und Key-Fenstern.
- `WindowLevel` steuert die Z-Reihenfolge. Verwenden Sie `MacWindowLevelPopUpMenu` für Overlays in der Menüleiste.
- `MacOptions.ActivationPolicy` steuert die Anwendung als Ganzes, einschließlich der Darstellung im Dock und in der Menüleiste. Ein nicht aktivierendes Panel erfordert keine Accessory-Aktivierungsrichtlinie; eine Anwendung kann sie jedoch weiterhin verwenden, um ihr Dock-Symbol auszublenden.
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

**Fensterebene** (`MacWindowLevel`)

- `MacWindowLevelNormal` – Standardfensterebene (Voreinstellung)
- `MacWindowLevelFloating` – Schwebt über normalen Fenstern
- `MacWindowLevelTornOffMenu` – Ebene für abgetrennte Menüs
- `MacWindowLevelModalPanel` – Ebene für modale Panels
- `MacWindowLevelMainMenu` – Hauptmenüebene
- `MacWindowLevelStatus` – Ebene für Statusfenster
- `MacWindowLevelPopUpMenu` – Ebene für Einblendmenüs
- `MacWindowLevelScreenSaver` – Bildschirmschonerebene

Ein explizites `WindowLevel` hat Vorrang vor `AlwaysOnTop` und `PanelPreferences.FloatingPanel`. Ohne explizite Ebene wird `AlwaysOnTop` oder ein schwebendes Panel in `MacWindowLevelFloating` aufgelöst; andernfalls ist die Ebene `MacWindowLevelNormal`. Ein späterer Aufruf von `SetAlwaysOnTop` bleibt eine explizite Laufzeitänderung.

**Sammlungsverhalten** (`MacWindowCollectionBehavior`)

Steuert das Verhalten des Fensters über macOS Spaces und den Vollbildmodus hinweg. Diese Bitmaskenwerte lassen sich mit bitweisem ODER (`|`) kombinieren.

**Space-Verhalten:**

- `MacWindowCollectionBehaviorDefault` – Verwendet FullScreenPrimary (Voreinstellung, abwärtskompatibel)
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
- `MacWindowCollectionBehaviorFullScreenAuxiliary` – Kann Vollbild-Apps überlagern
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

Steuert das Verhalten von Fenster-Tabs unter macOS 10.12 und neuer. Mit Fenster-Tabs lassen sich mehrere Fenster als Tabs gruppieren.

**Optionen:**

- `MacWindowTabbingModeDefault` – Nullwert als Sentinel (nicht explizit festgelegt). Zur Laufzeit sind Fenster-Tabs standardmäßig nicht zulässig
- `MacWindowTabbingModeAutomatic` – Das System bestimmt das Verhalten der Fenster-Tabs
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

Ermöglicht die detaillierte Steuerung der zugrunde liegenden `WKWebView`-Konfiguration. Alle Felder sind optional – nicht festgelegte Felder lassen den WebKit-Standard unverändert.

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

- `TabFocusesLinks` – Bei `true` verschiebt die Tabulatortaste den Fokus auf Links und Formularsteuerelemente (Standard: `false`)
- `TextInteractionEnabled` – Bei `true` können Benutzer Text in der Webview auswählen und mit ihm interagieren (Standard: `true`)
- `FullscreenEnabled` – Bei `true` können Webinhalte über die HTML Fullscreen API in den Vollbildmodus wechseln (Standard: `false`). Erfordert macOS 12.3+.
- `AllowsBackForwardNavigationGestures` – Bei `true` lösen horizontale Wischgesten die Rückwärts-/Vorwärtsnavigation aus (Standard: `false`)
- `AllowsMagnification` – Bei `true` ist das Zoomen per Auf- und Zuziehen in der Webview aktiviert (Standard: `false`)
- `AllowsAirPlayForMediaPlayback` – Bei `true` können Medien an AirPlay-Geräte gestreamt werden (Standard: `true`)
- `JavaScriptCanOpenWindowsAutomatically` – Bei `true` kann JavaScript ohne Benutzeraktion neue Fenster öffnen (Standard: `false`)
- `MinimumFontSize` – Mindestschriftgröße in Punkt. Mit `optional.NewVar(12.0)` festlegen. Ist das Feld nicht festgelegt, gilt der WebKit-Standard.
- `ApplicationNameForUserAgent` – Überschreibt das Suffix des Anwendungsnamens in der WebKit-User-Agent-Zeichenfolge. Nützlich, wenn Websites die standardmäßige Kennung `"wails.io"` ablehnen (z. B. YouTube-Einbettungen). Leer lassen, um den Standard beizubehalten.
- `EnableAutoplayWithoutUserAction` – Bei `true` können Audio und Video ohne Benutzeraktion automatisch wiedergegeben werden. Entspricht `WKWebViewConfiguration.mediaTypesRequiringUserActionForPlayback = WKAudiovisualMediaTypeNone` (Standard: `false`)

### Windows-Optionen (pro Fenster)

Die Struktur pro Fenster ist `application.WindowsWindow` – **nicht** `WindowsOptions` (das ist die Struktur auf *Anwendungs*ebene).

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

- Deaktiviert die Menüleiste für das Fenster. Bei `true` zeigt das Fenster selbst dann keine Menüleiste an, wenn eine konfiguriert ist.
- Standard: `false`

**BackdropType** (`BackdropType`)

- `application.Auto` – Systemstandard
- `application.None` – Kein Hintergrundeffekt
- `application.Mica` – Mica-Material (Windows 11)
- `application.Acrylic` – Acrylic-Material (Windows 11)
- `application.Tabbed` – Tabbed-Material (Windows 11)

Es gibt keine Konstanten im Stil von `WindowsBackdropTypeMica` – verwenden Sie `application.Mica` usw.

**CustomTheme** (`ThemeSettings`)

- Wert (kein Zeiger). Benutzerdefinierte Farben für den Dunkel-/Hellmodus des Fensterrahmens, des Texts und Hintergrunds der Titelleiste sowie der Menüleiste.

**DisableFramelessWindowDecorations** (`bool`)

- Deaktiviert die standardmäßigen Dekorationen rahmenloser Fenster (Aero-Schatten, abgerundete Ecken).

**NonClientRegionSupport** (`bool`)

- Aktiviert die native WebView2-Unterstützung für `app-region: drag` / `app-region: no-drag` bei rahmenlosen benutzerdefinierten Titelleisten.
- Dies dient ausschließlich dem einfachen nativen Ziehen der App. Es bietet weder natives Verhalten für benutzerdefinierte Titelleistenschaltflächen noch Windows 11 Snap Assist / Snap Layouts für benutzerdefinierte Maximieren-Schaltflächen.

**WebView2CompositionHosting** (`bool`)

- Aktiviert die von Wails verwaltete `--wails-non-client-region`-Unterstützung für benutzerdefinierte Titelleistenschaltflächen mit nativem Windows-Verhalten, einschließlich Windows 11 Snap Assist / Snap Layouts bei benutzerdefinierten Maximieren-Schaltflächen.
- Experimentell. Hierbei wird WebView2 über `ICoreWebView2CompositionController` und DirectComposition statt über den standardmäßigen HWND-gehosteten Controller bereitgestellt.
- Kann mit `NonClientRegionSupport` kombiniert werden, wenn ein Fenster sowohl native WebView2-Unterstützung für `app-region` als auch von Wails verwaltete Bereiche für benutzerdefinierte Titelleistenschaltflächen benötigt.

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

Ausführliche Informationen zu Verhalten, Abwägungen und dem zugehörigen CSS finden Sie unter [Rahmenlose Fenster unter Windows](/features/windows/frameless/#native-non-client-regions-on-windows).

### Linux-Optionen (pro Fenster)

Die Struktur pro Fenster ist `application.LinuxWindow` – **nicht** `LinuxOptions`.

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

Einige Windows-spezifische Optionen müssen auf Anwendungsebene statt pro Fenster konfiguriert werden. Der Grund ist, dass WebView2 pro Benutzerdatenpfad eine einzige Browserumgebung für alle Fenster verwendet.

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
- Die verfügbaren Flags finden Sie unter [WebView2-Browser-Flags](https://learn.microsoft.com/en-us/microsoft-edge/webview2/concepts/webview-features-flags)
- Beispiel: `"msWebView2EnableDraggableRegions"`

**DisabledFeatures** (`[]string`)

- Liste der zu deaktivierenden WebView2-Feature-Flags
- Wails deaktiviert automatisch `msSmartScreenProtection`
- Beispiel: `"msExperimentalFeature"`

**AdditionalBrowserArgs** (`[]string`)

- Chromium-Befehlszeilenargumente, die an den Browserprozess übergeben werden
- Müssen das Präfix `--` enthalten (z. B. `"--remote-debugging-port=9222"`)
- Die verfügbaren Argumente finden Sie unter [Chromium-Befehlszeilenschalter](https://peter.sh/experiments/chromium-command-line-switches/)

@note{type="caution" title="Wichtig"}
Diese Flags gelten global für ALLE Fenster, da WebView2 pro Benutzerdatenpfad eine einzige Browserumgebung verwendet. Sie können nicht für verschiedene Fenster unterschiedliche Browser-Flags verwenden.

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

Hier ist eine produktionsreife Fensterkonfiguration:

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

Frontend-Assets werden auf Anwendungsebene (`application.Options.Assets`) und nicht für jedes Fenster separat bereitgestellt.

## Nächste Schritte

- [Fenstergrundlagen](/features/windows/basics/) – Fenster erstellen und steuern
- [Mehrere Fenster](/features/windows/multiple/) – Muster für Anwendungen mit mehreren Fenstern
- [Rahmenlose Fenster](/features/windows/frameless/) – Benutzerdefinierte Fensterdekoration
- [Fensterereignisse](/features/windows/events/) – Lebenszyklusereignisse

---

**Fragen?** Stellen Sie sie auf [Discord](https://discord.gg/JDdSxwjhGf) oder sehen Sie sich die [Beispiele](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples) an.
