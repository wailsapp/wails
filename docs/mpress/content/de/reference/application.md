---
title: "Application-API"
description: "Vollständige Referenz der Application-API"
slug: "reference/application"
sourcePath: "reference/application.md"
---

## Überblick

Die `Application` ist das Herzstück Ihrer Wails-Anwendung. Sie verwaltet Fenster, Dienste und Ereignisse und bietet Zugriff auf alle Plattformfunktionen.

## Anwendung erstellen

```go
import "github.com/wailsapp/wails/v3/pkg/application"

app := application.New(application.Options{
    Name:        "My App",
    Description: "My awesome application",
    Services: []application.Service{
        application.NewService(&MyService{}),
    },
})
```

## Kernmethoden

### Run()

Startet die Ereignisschleife der Anwendung.

```go
func (a *App) Run() error
```

**Beispiel:**

```go
err := app.Run()
if err != nil {
    log.Fatal(err)
}
```

**Rückgabewert:** Fehler, wenn der Start fehlschlägt

### Quit()

Beendet die Anwendung ordnungsgemäß.

```go
func (a *App) Quit()
```

**Beispiel:**

```go
// In a menu handler
menu.Add("Quit").OnClick(func(ctx *application.Context) {
    app.Quit()
})
```

### Config()

Gibt die Anwendungskonfiguration zurück.

```go
func (a *App) Config() Options
```

**Beispiel:**

```go
config := app.Config()
fmt.Println("App name:", config.Name)
```

## Fensterverwaltung

### app.Window.New()

Erstellt ein neues Webview-Fenster mit Standardoptionen.

```go
func (wm *WindowManager) New() *WebviewWindow
```

**Beispiel:**

```go
window := app.Window.New()
window.Show()
```

### app.Window.NewWithOptions()

Erstellt ein neues Webview-Fenster mit benutzerdefinierten Optionen.

```go
func (wm *WindowManager) NewWithOptions(options WebviewWindowOptions) *WebviewWindow
```

**Beispiel:**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "My Window",
    Width:  800,
    Height: 600,
    BackgroundColour: application.NewRGB(255, 255, 255),
})
```

### app.Window.GetByName()

Ruft ein Fenster anhand seines Namens ab. Gibt das Fenster sowie die Angabe zurück, ob es gefunden wurde.

```go
func (wm *WindowManager) GetByName(name string) (Window, bool)
```

**Beispiel:**

```go
if window, ok := app.Window.GetByName("main"); ok {
    window.Show()
}
```

### app.Window.GetAll()

Gibt alle Anwendungsfenster zurück.

```go
func (wm *WindowManager) GetAll() []Window
```

**Beispiel:**

```go
windows := app.Window.GetAll()
for _, window := range windows {
    fmt.Println("Window:", window.Name())
}
```

## Manager

Die Application bietet über Eigenschaften Zugriff auf verschiedene Manager:

```go
app.Window       // Window management
app.Menu         // Menu management
app.Dialog       // Dialog management
app.Event        // Event management
app.Clipboard    // Clipboard operations
app.Screen       // Screen information
app.SystemTray   // System tray
app.Browser      // Browser operations
app.Env          // Environment variables
app.ContextMenu  // Context-menu management
app.KeyBinding   // Global keyboard shortcuts
app.Logger       // *slog.Logger
```

### Anwendungsbeispiel

```go
// Create window
window := app.Window.New()

// Show dialog
app.Dialog.Info().SetMessage("Hello!").Show()

// Copy to clipboard
app.Clipboard.SetText("Copied text")

// Get screens
screens := app.Screen.GetAll()
```

## Dienstverwaltung

### RegisterService()

Registriert einen Dienst bei der Anwendung.

```go
func (a *App) RegisterService(service Service)
```

`RegisterService` gibt nichts zurück; Fehler bei der Dienstinitialisierung treten während `app.Run()` durch Fehler von `ServiceStartup` zutage.

**Beispiel:**

```go
type MyService struct {
    app *application.App
}

func NewMyService(app *application.App) *MyService {
    return &MyService{app: app}
}

// Register after app creation
app.RegisterService(application.NewService(NewMyService(app)))
```

## Ereignisverwaltung

### app.Event.Emit()

Löst ein benutzerdefiniertes Ereignis aus. Gibt `true` zurück, wenn ein Hook das Auslösen abgebrochen hat.

```go
func (em *EventManager) Emit(name string, data ...any) bool
```

**Beispiel:**

```go
// Emit event with data
app.Event.Emit("user-logged-in", map[string]interface{}{
    "username": "john",
    "timestamp": time.Now(),
})
```

### app.Event.On()

Lauscht auf benutzerdefinierte Ereignisse. Gibt eine `func()` zum Abbestellen zurück.

```go
func (em *EventManager) On(name string, callback func(*CustomEvent)) func()
```

**Beispiel:**

```go
app.Event.On("user-logged-in", func(e *application.CustomEvent) {
    data := e.Data.(map[string]interface{})
    username := data["username"].(string)
    fmt.Println("User logged in:", username)
})
```

### app.Event.OnApplicationEvent()

Lauscht auf Lebenszyklusereignisse der Anwendung. Der Parameter `eventType` ist vom Typ `events.ApplicationEventType` (aus dem Paket `events`).

```go
func (em *EventManager) OnApplicationEvent(
    eventType events.ApplicationEventType,
    callback func(*ApplicationEvent),
) func()
```

**Beispiel:**

```go
import "github.com/wailsapp/wails/v3/pkg/events"

// Listen for app-started
app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(e *application.ApplicationEvent) {
    fmt.Println("Application started")
})

// Application shutdown is NOT an event constant; register cleanup via:
app.OnShutdown(func() {
    fmt.Println("Application shutting down")
})
```

## Dialogmethoden

Der Zugriff auf Dialoge erfolgt über den Manager `app.Dialog`. Die vollständige Referenz finden Sie in der [Dialogs-API](/reference/dialogs/).

### Meldungsdialoge

```go
// Information dialog
app.Dialog.Info().
    SetTitle("Success").
    SetMessage("Operation completed!").
    Show()

// Error dialog
app.Dialog.Error().
    SetTitle("Error").
    SetMessage("Something went wrong.").
    Show()

// Warning dialog
app.Dialog.Warning().
    SetTitle("Warning").
    SetMessage("This action cannot be undone.").
    Show()
```

### Fragedialoge

Fragedialoge verarbeiten Benutzerantworten mithilfe von Button-Callbacks:

```go
dialog := app.Dialog.Question().
    SetTitle("Confirm").
    SetMessage("Continue?")

yes := dialog.AddButton("Yes")
yes.OnClick(func() {
    // Handle yes
})

no := dialog.AddButton("No")
no.OnClick(func() {
    // Handle no
})

dialog.SetDefaultButton(yes)
dialog.SetCancelButton(no)
dialog.Show()
```

### Dateidialoge

```go
// Open file dialog
path, err := app.Dialog.OpenFile().
    SetTitle("Select File").
    AddFilter("Images", "*.png;*.jpg").
    PromptForSingleSelection()

// Save file dialog
path, err := app.Dialog.SaveFile().
    SetTitle("Save File").
    SetFilename("document.pdf").
    AddFilter("PDF", "*.pdf").
    PromptForSingleSelection()

// Folder selection (use OpenFile with directory options)
path, err := app.Dialog.OpenFile().
    SetTitle("Select Folder").
    CanChooseDirectories(true).
    CanChooseFiles(false).
    PromptForSingleSelection()
```

## Logger

Die Anwendung stellt einen strukturierten Logger bereit:

```go
app.Logger.Info("Message", "key", "value")
app.Logger.Error("Error occurred", "error", err)
app.Logger.Debug("Debug info")
app.Logger.Warn("Warning message")
```

**Beispiel:**

```go
func (s *MyService) ProcessData(data string) error {
    s.app.Logger.Info("Processing data", "length", len(data))
    
    if err := process(data); err != nil {
        s.app.Logger.Error("Processing failed", "error", err)
        return err
    }
    
    s.app.Logger.Info("Processing complete")
    return nil
}
```

## Verarbeitung von Rohnachrichten

Für Anwendungen, die eine direkte Low-Level-Kontrolle über die Kommunikation vom Frontend zum Backend benötigen, stellt Wails die Option `RawMessageHandler` bereit. Sie umgeht das standardmäßige Binding-System.

@note{type="info"}
Raw Messages sollten nur als letztes Mittel verwendet werden. Das standardmäßige Binding-System ist hochgradig optimiert und für fast alle Anwendungen ausreichend. Verwenden Sie Raw Messages nur, wenn Sie Ihre Anwendung profiliert und bestätigt haben, dass Bindings einen Engpass darstellen.

@end

### RawMessageHandler

`RawMessageHandler` ist ein Feld von `application.Options` und keine Methode. Die Runtime ruft es für jede Raw Message auf, die das Frontend über `System.invoke()` sendet.

```go
type Options struct {
    // ... other fields ...
    RawMessageHandler func(window Window, message string, originInfo *OriginInfo)
}
```

`OriginInfo` enthält `Origin`, `TopOrigin` und `IsMainFrame`. Die Plattformen befüllen unterschiedliche Teilmengen – die plattformspezifische Matrix finden Sie im [Leitfaden zu Raw Messages](/guides/raw-messages/).

**Beispiel:**

```go
app := application.New(application.Options{
    Name: "My App",
    RawMessageHandler: func(window application.Window, message string, originInfo *application.OriginInfo) {
        // Handle the raw message
        fmt.Printf("Received from %s (%s): %s\n", window.Name(), originInfo.Origin, message)

        // You can respond using events
        window.EmitEvent("response", processMessage(message))
    },
})
```

Weitere Einzelheiten finden Sie im [Leitfaden zu Raw Messages](/guides/raw-messages/).

## Plattformspezifische Optionen

### Windows-Optionen

Konfigurieren Sie das Windows-spezifische Verhalten auf Anwendungsebene:

```go
app := application.New(application.Options{
    Name: "My App",
    Windows: application.WindowsOptions{
        // WebView2 browser flags (apply to ALL windows)
        EnabledFeatures:       []string{"msWebView2EnableDraggableRegions"},
        DisabledFeatures:      []string{"msExperimentalFeature"},
        AdditionalBrowserArgs: []string{"--remote-debugging-port=9222"},

        // Other Windows options
        WndClass:                      "MyAppClass",
        WebviewUserDataPath:           "",  // Default: %APPDATA%\[BinaryName.exe]
        WebviewBrowserPath:            "",  // Default: system WebView2
        DisableQuitOnLastWindowClosed: false,
    },
})
```

**Browser-Flags:**

- `EnabledFeatures` – zu aktivierende WebView2-Feature-Flags
- `DisabledFeatures` – zu deaktivierende WebView2-Feature-Flags
- `AdditionalBrowserArgs` – Chromium-Befehlszeilenargumente

Eine ausführliche Dokumentation finden Sie unter [Fensteroptionen – Windows-Optionen auf Anwendungsebene](/features/windows/options/#application-level-windows-options).

### Mac-Optionen

```go
app := application.New(application.Options{
    Name: "My App",
    Mac: application.MacOptions{
        ActivationPolicy: application.ActivationPolicyRegular,
        ApplicationShouldTerminateAfterLastWindowClosed: true,
    },
})
```

### Linux-Optionen

```go
app := application.New(application.Options{
    Name: "My App",
    Linux: application.LinuxOptions{
        ProgramName:                   "my-app",
        DisableQuitOnLastWindowClosed: false,
    },
})
```

## Vollständiges Anwendungsbeispiel

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name:        "My Application",
        Description: "A demo application",
        Mac: application.MacOptions{
            ApplicationShouldTerminateAfterLastWindowClosed: true,
        },
    })

    // Create main window
    window := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:            "My App",
        Width:            1024,
        Height:           768,
        MinWidth:         800,
        MinHeight:        600,
        BackgroundColour: application.NewRGB(255, 255, 255),
        URL:              "http://wails.localhost/",
    })

    window.Center()
    window.Show()

    app.Run()
}
```
