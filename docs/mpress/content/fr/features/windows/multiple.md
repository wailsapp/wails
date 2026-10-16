---
title: "Fenêtres multiples"
description: "Modèles et bonnes pratiques pour les applications multifenêtres"
slug: "features/windows/multiple"
sourcePath: "features/windows/multiple.md"
---

## Applications multifenêtres

Wails v3 offre une **prise en charge native de plusieurs fenêtres** pour créer des fenêtres de paramètres, de documents, de palettes d’outils et d’inspection. Suivez les fenêtres, permettez-leur de communiquer entre elles et gérez leur cycle de vie à l’aide d’API simples et cohérentes.

### Fenêtre principale et fenêtre de paramètres

```go
package main

import "github.com/wailsapp/wails/v3/pkg/application"

type App struct {
    app            *application.App
    mainWindow     *application.WebviewWindow
    settingsWindow *application.WebviewWindow
}

func main() {
    app := &App{}
    
    app.app = application.New(application.Options{
        Name: "Multi-Window App",
    })
    
    // Create main window
    app.mainWindow = app.app.Window.NewWithOptions(application.WebviewWindowOptions{
        Name:   "main",
        Title:  "Main Application",
        Width:  1200,
        Height: 800,
    })
    
    // Create settings window (hidden initially)
    app.settingsWindow = app.app.Window.NewWithOptions(application.WebviewWindowOptions{
        Name:   "settings",
        Title:  "Settings",
        Width:  600,
        Height: 400,
        Hidden: true,
    })
    
    app.app.Run()
}

// Show settings from main window
func (a *App) ShowSettings() {
    if a.settingsWindow != nil {
        a.settingsWindow.Show()
        a.settingsWindow.Focus()
    }
}
```

**Points clés :**

- Fenêtre principale toujours visible
- Fenêtre de paramètres créée, mais masquée
- Affichage des paramètres à la demande
- Réutilisation de la même fenêtre (ne pas en créer plusieurs)

## Suivi des fenêtres

### Obtenir toutes les fenêtres

```go
windows := app.Window.GetAll()
fmt.Printf("Total windows: %d\n", len(windows))

for _, window := range windows {
    fmt.Printf("- %s (ID: %d)\n", window.Name(), window.ID())
}
```

### Rechercher une fenêtre précise

```go
// By name
if settings, ok := app.Window.GetByName("settings"); ok {
    settings.Show()
}

// By ID
if window, ok := app.Window.GetByID(123); ok {
    window.Focus()
}

// Current (focused) window
current := app.Window.Current()
```

### Modèle de registre de fenêtres

Suivez les fenêtres de votre application :

```go
type WindowManager struct {
    windows map[string]*application.WebviewWindow
    mu      sync.RWMutex
}

func (wm *WindowManager) Register(name string, window *application.WebviewWindow) {
    wm.mu.Lock()
    defer wm.mu.Unlock()
    wm.windows[name] = window
}

func (wm *WindowManager) Get(name string) *application.WebviewWindow {
    wm.mu.RLock()
    defer wm.mu.RUnlock()
    return wm.windows[name]
}

func (wm *WindowManager) Remove(name string) {
    wm.mu.Lock()
    defer wm.mu.Unlock()
    delete(wm.windows, name)
}
```

## Communication entre les fenêtres

### Utiliser des événements

Les fenêtres communiquent au moyen du système d’événements :

```go
// In main window - emit event
app.Event.Emit("settings-changed", map[string]interface{}{
    "theme": "dark",
    "fontSize": 14,
})

// In settings window - listen for event
app.Event.On("settings-changed", func(event *application.CustomEvent) {
    data := event.Data.(map[string]interface{})
    theme := data["theme"].(string)
    fontSize := data["fontSize"].(int)
    
    // Update UI
    updateSettings(theme, fontSize)
})
```

### Modèle d’état partagé

Utilisez un gestionnaire d’état partagé :

```go
type AppState struct {
    theme    string
    fontSize int
    mu       sync.RWMutex
}

var state = &AppState{
    theme:    "light",
    fontSize: 12,
}

func (s *AppState) SetTheme(theme string) {
    s.mu.Lock()
    s.theme = theme
    s.mu.Unlock()
    
    // Notify all windows
    app.Event.Emit("theme-changed", theme)
}

func (s *AppState) GetTheme() string {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.theme
}
```

### Messages entre fenêtres

Envoyez des messages entre des fenêtres précises :

```go
// Get target window
if targetWindow, ok := app.Window.GetByName("preview"); ok {
    // Emit event to specific window
    targetWindow.EmitEvent("update-preview", previewData)
}
```

## Modèles courants

### Modèle 1 : fenêtres à instance unique

Veillez à ce qu’il n’existe qu’une seule instance d’une fenêtre :

```go
var settingsWindow *application.WebviewWindow

func ShowSettings(app *application.App) {
    // Create if doesn't exist
    if settingsWindow == nil {
        settingsWindow = app.Window.NewWithOptions(application.WebviewWindowOptions{
            Name:   "settings",
            Title:  "Settings",
            Width:  600,
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

### Modèle 2 : fenêtres de documents

Créez plusieurs instances du même type de fenêtre :

```go
type DocumentWindow struct {
    window   *application.WebviewWindow
    filePath string
    modified bool
}

var documents = make(map[string]*DocumentWindow)

func OpenDocument(app *application.App, filePath string) {
    // Check if already open
    if doc, exists := documents[filePath]; exists {
        doc.window.Show()
        doc.window.Focus()
        return
    }
    
    // Create new document window
    window := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:  filepath.Base(filePath),
        Width:  800,
        Height: 600,
    })
    
    doc := &DocumentWindow{
        window:   window,
        filePath: filePath,
        modified: false,
    }
    
    documents[filePath] = doc
    
    // Cleanup on close
    window.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
        delete(documents, filePath)
    })
    
    // Load document
    loadDocument(window, filePath)
}
```

### Modèle 3 : palettes d’outils

Créez des fenêtres flottantes qui restent au premier plan :

```go
func CreateToolPalette(app *application.App) *application.WebviewWindow {
    palette := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Name:          "tools",
        Title:         "Tools",
        Width:         200,
        Height:        400,
        AlwaysOnTop:   true,
        DisableResize: true,
    })

    return palette
}
```

### Modèle 4 : boîtes de dialogue modales (macOS uniquement)

Créez des fenêtres enfants qui bloquent leur fenêtre parente :

```go
func ShowModaldialog(parent *application.WebviewWindow, title string) {
    dialog := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:         title,
        Width:         400,
        Height:        200,
        AlwaysOnTop:   true,
        DisableResize: true,
    })

    parent.AttachModal(dialog)
}
```

### Modèle 5 : fenêtres d’inspection ou d’aperçu

Créez des fenêtres liées qui se mettent à jour ensemble :

```go
type EditorApp struct {
    editor  *application.WebviewWindow
    preview *application.WebviewWindow
}

func (e *EditorApp) UpdatePreview(content string) {
    if e.preview != nil && e.preview.IsVisible() {
        e.preview.EmitEvent("content-changed", content)
    }
}

func (e *EditorApp) TogglePreview() {
    if e.preview == nil {
        e.preview = app.Window.NewWithOptions(application.WebviewWindowOptions{
            Name:   "preview",
            Title:  "Preview",
            Width:  600,
            Height: 800,
        })
        
        e.preview.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
            e.preview = nil
        })
    }
    
    if e.preview.IsVisible() {
        e.preview.Hide()
    } else {
        e.preview.Show()
    }
}
```

## Relations parent-enfant

### Créer des fenêtres enfants

```go
childWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "Child Window",
})

parentWindow.AttachModal(childWindow)
```

**Comportement :**

- La fenêtre enfant reste au-dessus de la fenêtre parente
- La fenêtre enfant se déplace avec la fenêtre parente
- La fenêtre enfant bloque les interactions avec la fenêtre parente

**Prise en charge des plateformes :**

| macOS | Windows | Linux |
| --- | --- | --- |
| ✅ | ❌ | ❌ |

### Comportement modal

Créez un comportement similaire à celui d’une fenêtre modale :

```go
func ShowModaldialog(parent *application.WebviewWindow, title string) {
    dialog := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:       title,
        Width:       400,
        Height:      200,
    })
    
    parent.AttachModal(dialog)
}
```

## Gestion du cycle de vie des fenêtres

### Fonctions de rappel de création

Recevez une notification lors de la création de fenêtres :

```go
app.Window.OnCreate(func(window application.Window) {
    fmt.Printf("Window created: %s\n", window.Name())

    // Configure all new windows
    window.SetMinSize(400, 300)
})
```

### Fonctions de rappel de destruction

Effectuez le nettoyage lors de la destruction des fenêtres :

```go
window.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
    fmt.Printf("Window %s is closing\n", window.Name())
    
    // Cleanup resources
    cleanup(window.ID())
    
    // Remove from tracking
    removeFromRegistry(window.Name())
})
```

### Comportement à la fermeture de l’application

Contrôlez le moment où l’application se ferme :

```go
app := application.New(application.Options{
    Mac: application.MacOptions{
        // Don't quit when last window closes
        ApplicationShouldTerminateAfterLastWindowClosed: false,
    },
})
```

**Cas d’utilisation :**

- Applications de la zone de notification système
- Services en arrière-plan
- Applications de la barre des menus (macOS)

## Gestion de la mémoire

### Prévenir les fuites

Nettoyez toujours les références aux fenêtres :

```go
var windows = make(map[string]*application.WebviewWindow)

func CreateWindow(name string) {
    window := app.Window.New()
    windows[name] = window
    
    // IMPORTANT: Clean up on close
    window.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
        delete(windows, name)
    })
}
```

### Fermer les fenêtres

```go
// Close — dispatches WindowClosing; a RegisterHook can call e.Cancel().
window.Close()
```

Il n’existe pas de `window.Destroy()` dans la v3. Pour empêcher la fermeture, utilisez `RegisterHook(events.Common.WindowClosing, ...)` et appelez `event.Cancel()`.

### Nettoyage des ressources

```go
type ManagedWindow struct {
    window     *application.WebviewWindow
    resources  []io.Closer
}

func (mw *ManagedWindow) Destroy() {
    // Close all resources
    for _, resource := range mw.resources {
        resource.Close()
    }

    // Close the window — WindowClosing listeners run before the window goes away.
    mw.window.Close()
}
```

## Modèles avancés

### Pool de fenêtres

Réutilisez les fenêtres au lieu d’en créer de nouvelles :

```go
type WindowPool struct {
    available []*application.WebviewWindow
    inUse     map[uint]*application.WebviewWindow
    mu        sync.Mutex
}

func (wp *WindowPool) Acquire() *application.WebviewWindow {
    wp.mu.Lock()
    defer wp.mu.Unlock()
    
    // Reuse available window
    if len(wp.available) > 0 {
        window := wp.available[0]
        wp.available = wp.available[1:]
        wp.inUse[window.ID()] = window
        return window
    }
    
    // Create new window
    window := app.Window.New()
    wp.inUse[window.ID()] = window
    return window
}

func (wp *WindowPool) Release(window *application.WebviewWindow) {
    wp.mu.Lock()
    defer wp.mu.Unlock()
    
    delete(wp.inUse, window.ID())
    window.Hide()
    wp.available = append(wp.available, window)
}
```

### Groupes de fenêtres

Gérez ensemble les fenêtres associées :

```go
type WindowGroup struct {
    name    string
    windows []*application.WebviewWindow
}

func (wg *WindowGroup) Add(window *application.WebviewWindow) {
    wg.windows = append(wg.windows, window)
}

func (wg *WindowGroup) ShowAll() {
    for _, window := range wg.windows {
        window.Show()
    }
}

func (wg *WindowGroup) HideAll() {
    for _, window := range wg.windows {
        window.Hide()
    }
}

func (wg *WindowGroup) CloseAll() {
    for _, window := range wg.windows {
        window.Close()
    }
}
```

### Gestion des espaces de travail

Enregistrez et restaurez la disposition des fenêtres :

```go
type WindowLayout struct {
    Windows []WindowState `json:"windows"`
}

type WindowState struct {
    Name   string `json:"name"`
    X      int    `json:"x"`
    Y      int    `json:"y"`
    Width  int    `json:"width"`
    Height int    `json:"height"`
}

func SaveLayout() *WindowLayout {
    layout := &WindowLayout{}
    
    for _, window := range app.Window.GetAll() {
        x, y := window.Position()
        width, height := window.Size()
        
        layout.Windows = append(layout.Windows, WindowState{
            Name:   window.Name(),
            X:      x,
            Y:      y,
            Width:  width,
            Height: height,
        })
    }
    
    return layout
}

func RestoreLayout(layout *WindowLayout) {
    for _, state := range layout.Windows {
        if window, ok := app.Window.GetByName(state.Name); ok {
            window.SetPosition(state.X, state.Y)
            window.SetSize(state.Width, state.Height)
        }
    }
}
```

## Exemple complet

Voici une application multifenêtre prête pour la production :

```go
package main

import (
    "encoding/json"
    "os"
    "sync"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)

type MultiWindowApp struct {
    app     *application.App
    windows map[string]*application.WebviewWindow
    mu      sync.RWMutex
}

func main() {
    mwa := &MultiWindowApp{
        windows: make(map[string]*application.WebviewWindow),
    }
    
    mwa.app = application.New(application.Options{
        Name: "Multi-Window Application",
        Mac: application.MacOptions{
            ApplicationShouldTerminateAfterLastWindowClosed: false,
        },
    })
    
    // Create main window
    mwa.CreateMainWindow()
    
    // Load saved layout
    mwa.LoadLayout()
    
    mwa.app.Run()
}

func (mwa *MultiWindowApp) CreateMainWindow() {
    window := mwa.app.Window.NewWithOptions(application.WebviewWindowOptions{
        Name:   "main",
        Title:  "Main Application",
        Width:  1200,
        Height: 800,
    })
    
    mwa.RegisterWindow("main", window)
}

func (mwa *MultiWindowApp) ShowSettings() {
    if window := mwa.GetWindow("settings"); window != nil {
        window.Show()
        window.Focus()
        return
    }
    
    window := mwa.app.Window.NewWithOptions(application.WebviewWindowOptions{
        Name:   "settings",
        Title:  "Settings",
        Width:  600,
        Height: 400,
    })
    
    mwa.RegisterWindow("settings", window)
}

func (mwa *MultiWindowApp) OpenDocument(path string) {
    name := "doc-" + path
    
    if window := mwa.GetWindow(name); window != nil {
        window.Show()
        window.Focus()
        return
    }
    
    window := mwa.app.Window.NewWithOptions(application.WebviewWindowOptions{
        Name:  name,
        Title: path,
        Width: 800,
        Height: 600,
    })
    
    mwa.RegisterWindow(name, window)
}

func (mwa *MultiWindowApp) RegisterWindow(name string, window *application.WebviewWindow) {
    mwa.mu.Lock()
    mwa.windows[name] = window
    mwa.mu.Unlock()
    
    window.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
        mwa.UnregisterWindow(name)
    })
}

func (mwa *MultiWindowApp) UnregisterWindow(name string) {
    mwa.mu.Lock()
    delete(mwa.windows, name)
    mwa.mu.Unlock()
}

func (mwa *MultiWindowApp) GetWindow(name string) *application.WebviewWindow {
    mwa.mu.RLock()
    defer mwa.mu.RUnlock()
    return mwa.windows[name]
}

func (mwa *MultiWindowApp) SaveLayout() {
    layout := make(map[string]WindowState)
    
    mwa.mu.RLock()
    for name, window := range mwa.windows {
        x, y := window.Position()
        width, height := window.Size()
        
        layout[name] = WindowState{
            X:      x,
            Y:      y,
            Width:  width,
            Height: height,
        }
    }
    mwa.mu.RUnlock()
    
    data, _ := json.Marshal(layout)
    os.WriteFile("layout.json", data, 0644)
}

func (mwa *MultiWindowApp) LoadLayout() {
    data, err := os.ReadFile("layout.json")
    if err != nil {
        return
    }
    
    var layout map[string]WindowState
    if err := json.Unmarshal(data, &layout); err != nil {
        return
    }
    
    for name, state := range layout {
        if window := mwa.GetWindow(name); window != nil {
            window.SetPosition(state.X, state.Y)
            window.SetSize(state.Width, state.Height)
        }
    }
}

type WindowState struct {
    X      int `json:"x"`
    Y      int `json:"y"`
    Width  int `json:"width"`
    Height int `json:"height"`
}
```

## Bonnes pratiques

### ✅ À faire

- **Suivez les fenêtres** – Conservez des références pour y accéder facilement
- **Nettoyez les ressources lors de la destruction** – Évitez les fuites de mémoire
- **Utilisez des événements pour la communication** – Architecture découplée
- **Réutilisez les fenêtres** – Ne créez pas de doublons
- **Enregistrez et restaurez les dispositions** – Pour une meilleure expérience utilisateur
- **Gérez la fermeture des fenêtres** – Demandez confirmation avant de fermer une fenêtre contenant des données non enregistrées

### ❌ À ne pas faire

- **Ne créez pas un nombre illimité de fenêtres** – Problèmes de mémoire et de performances
- **N’oubliez pas de nettoyer les ressources** – Fuites de mémoire
- **N’utilisez pas les variables globales sans précaution** – Problèmes de sûreté des accès concurrents
- **Ne bloquez pas la création des fenêtres** – Créez-les de manière asynchrone si nécessaire
- **N’ignorez pas les différences entre les plateformes** – Effectuez des tests sur toutes les plateformes

## Étapes suivantes

- [Principes de base des fenêtres](/features/windows/basics/) – Découvrez les fondamentaux de la gestion des fenêtres
- [Événements de fenêtre](/features/windows/events/) – Gérez les événements du cycle de vie des fenêtres
- [Système d’événements](/features/events/system/) – Explorez le système d’événements en profondeur
- [Fenêtres sans cadre](/features/windows/frameless/) – Créez un habillage de fenêtre personnalisé

---

**Des questions ?** Posez-les sur [Discord](https://discord.gg/JDdSxwjhGf) ou consultez l’[exemple multifenêtre](https://github.com/wailsapp/wails/tree/master/v3/examples/multi-window).
