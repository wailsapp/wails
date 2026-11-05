---
title: "API Application"
description: "Référence complète de l’API Application"
slug: "reference/application"
sourcePath: "reference/application.md"
---

## Vue d’ensemble

L’`Application` est au cœur de votre application Wails. Elle gère les fenêtres, les services et les événements, et donne accès à toutes les fonctionnalités de la plateforme.

## Création d’une application

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

## Méthodes principales

### Run()

Démarre la boucle d’événements de l’application.

```go
func (a *App) Run() error
```

**Exemple :**

```go
err := app.Run()
if err != nil {
    log.Fatal(err)
}
```

**Valeur renvoyée :** une erreur si le démarrage échoue

### Quit()

Arrête proprement l’application.

```go
func (a *App) Quit()
```

**Exemple :**

```go
// In a menu handler
menu.Add("Quit").OnClick(func(ctx *application.Context) {
    app.Quit()
})
```

### Config()

Renvoie la configuration de l’application.

```go
func (a *App) Config() Options
```

**Exemple :**

```go
config := app.Config()
fmt.Println("App name:", config.Name)
```

## Gestion des fenêtres

### app.Window.New()

Crée une fenêtre webview avec les options par défaut.

```go
func (wm *WindowManager) New() *WebviewWindow
```

**Exemple :**

```go
window := app.Window.New()
window.Show()
```

### app.Window.NewWithOptions()

Crée une fenêtre webview avec des options personnalisées.

```go
func (wm *WindowManager) NewWithOptions(options WebviewWindowOptions) *WebviewWindow
```

**Exemple :**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "My Window",
    Width:  800,
    Height: 600,
    BackgroundColour: application.NewRGB(255, 255, 255),
})
```

### app.Window.GetByName()

Obtient une fenêtre à partir de son nom. Renvoie la fenêtre et indique si elle a été trouvée.

```go
func (wm *WindowManager) GetByName(name string) (Window, bool)
```

**Exemple :**

```go
if window, ok := app.Window.GetByName("main"); ok {
    window.Show()
}
```

### app.Window.GetAll()

Renvoie toutes les fenêtres de l’application.

```go
func (wm *WindowManager) GetAll() []Window
```

**Exemple :**

```go
windows := app.Window.GetAll()
for _, window := range windows {
    fmt.Println("Window:", window.Name())
}
```

## Gestionnaires

L’Application donne accès à différents gestionnaires au moyen de propriétés :

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

### Exemple d’utilisation

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

## Gestion des services

### RegisterService()

Enregistre un service auprès de l’application.

```go
func (a *App) RegisterService(service Service)
```

`RegisterService` ne renvoie rien ; les erreurs d’initialisation du service se manifestent par des échecs de `ServiceStartup` pendant `app.Run()`.

**Exemple :**

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

## Gestion des événements

### app.Event.Emit()

Émet un événement personnalisé. Renvoie `true` si un hook a annulé l’émission.

```go
func (em *EventManager) Emit(name string, data ...any) bool
```

**Exemple :**

```go
// Emit event with data
app.Event.Emit("user-logged-in", map[string]interface{}{
    "username": "john",
    "timestamp": time.Now(),
})
```

### app.Event.On()

Écoute les événements personnalisés. Renvoie une fonction `func()` de désabonnement.

```go
func (em *EventManager) On(name string, callback func(*CustomEvent)) func()
```

**Exemple :**

```go
app.Event.On("user-logged-in", func(e *application.CustomEvent) {
    data := e.Data.(map[string]interface{})
    username := data["username"].(string)
    fmt.Println("User logged in:", username)
})
```

### app.Event.OnApplicationEvent()

Écoute les événements du cycle de vie de l’application. Le paramètre `eventType` est de type `events.ApplicationEventType` (issu du package `events`).

```go
func (em *EventManager) OnApplicationEvent(
    eventType events.ApplicationEventType,
    callback func(*ApplicationEvent),
) func()
```

**Exemple :**

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

## Méthodes de boîte de dialogue

Les boîtes de dialogue sont accessibles au moyen du gestionnaire `app.Dialog`. Consultez l’[API des boîtes de dialogue](/reference/dialogs/) pour obtenir la référence complète.

### Boîtes de dialogue de message

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

### Boîtes de dialogue de question

Les boîtes de dialogue de question utilisent des fonctions de rappel associées aux boutons pour traiter les réponses de l’utilisateur :

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

### Boîtes de dialogue de fichiers

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

## Journalisation

L’application fournit un journaliseur structuré :

```go
app.Logger.Info("Message", "key", "value")
app.Logger.Error("Error occurred", "error", err)
app.Logger.Debug("Debug info")
app.Logger.Warn("Warning message")
```

**Exemple :**

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

## Traitement des messages bruts

Pour les applications qui nécessitent un contrôle direct et de bas niveau sur la communication entre le frontend et le backend, Wails fournit l’option `RawMessageHandler`. Celle-ci contourne le système de liaison standard.

@note{type="info"}
N’utilisez les messages bruts qu’en dernier recours. Le système de liaison standard est fortement optimisé et suffit à presque toutes les applications. N’utilisez les messages bruts que si vous avez profilé votre application et confirmé que les liaisons constituent un goulot d’étranglement.

@end

### RawMessageHandler

`RawMessageHandler` est un champ de `application.Options`, et non une méthode. Le runtime l’appelle pour chaque message brut envoyé depuis le frontend au moyen de `System.invoke()`.

```go
type Options struct {
    // ... other fields ...
    RawMessageHandler func(window Window, message string, originInfo *OriginInfo)
}
```

`OriginInfo` transporte `Origin`, `TopOrigin` et `IsMainFrame` (les plateformes renseignent des sous-ensembles différents — consultez le [guide des messages bruts](/guides/raw-messages/) pour connaître la matrice propre à chaque plateforme).

**Exemple :**

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

Pour plus de détails, consultez le [guide des messages bruts](/guides/raw-messages/).

## Options propres à chaque plateforme

### Options Windows

Configurez le comportement propre à Windows au niveau de l’application :

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

**Indicateurs du navigateur :**

- `EnabledFeatures` — indicateurs de fonctionnalités WebView2 à activer
- `DisabledFeatures` — indicateurs de fonctionnalités WebView2 à désactiver
- `AdditionalBrowserArgs` - Arguments de ligne de commande Chromium

Consultez [Options de fenêtre - Options Windows au niveau de l’application](/features/windows/options/#application-level-windows-options) pour obtenir une documentation détaillée.

### Options macOS

```go
app := application.New(application.Options{
    Name: "My App",
    Mac: application.MacOptions{
        ActivationPolicy: application.ActivationPolicyRegular,
        ApplicationShouldTerminateAfterLastWindowClosed: true,
    },
})
```

### Options Linux

```go
app := application.New(application.Options{
    Name: "My App",
    Linux: application.LinuxOptions{
        ProgramName:                   "my-app",
        DisableQuitOnLastWindowClosed: false,
    },
})
```

## Exemple d’application complet

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
