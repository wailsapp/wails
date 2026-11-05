---
title: "API des fenêtres"
description: "Référence complète de l’API des fenêtres"
slug: "reference/window"
sourcePath: "reference/window.md"
---

## Vue d’ensemble

L’API des fenêtres fournit des méthodes permettant de contrôler l’apparence, le comportement et le cycle de vie des fenêtres. Accédez-y à partir des instances de fenêtre ou du gestionnaire `app.Window`.

`Window` est une interface implémentée par `*application.WebviewWindow`. Les signatures de méthodes ci-dessous appartiennent à `*WebviewWindow`. De nombreuses méthodes de modification renvoient `Window` afin de permettre le chaînage ; les valeurs de retour sont documentées pour chaque méthode.

**Opérations courantes :**

- Créer et afficher des fenêtres
- Contrôler la taille, la position et l’état
- Gérer les événements de fenêtre
- Gérer le contenu des fenêtres
- Configurer l’apparence et le comportement

## Visibilité

### Show()

Affiche la fenêtre. Si elle était masquée, elle devient visible. Renvoie le récepteur pour permettre le chaînage.

```go
func (w *WebviewWindow) Show() Window
```

**Exemple :**

```go
window := app.Window.New()
window.Show()
```

### Hide()

Masque la fenêtre sans la fermer. Elle reste en mémoire et peut être affichée de nouveau. Renvoie le récepteur pour permettre le chaînage.

```go
func (w *WebviewWindow) Hide() Window
```

**Exemple :**

```go
// Hide window temporarily
window.Hide()

// Show it again later
window.Show()
```

**Cas d’utilisation :**

- Applications de la zone de notification qui se masquent dans celle-ci
- Assistants dans lesquels les fenêtres sont réutilisées
- Masquage temporaire pendant des opérations

### Close()

Ferme la fenêtre. Cette opération déclenche l’événement `WindowClosing`.

```go
func (w *WebviewWindow) Close()
```

**Exemple :**

```go
window.Close()
```

**Remarque :** si un hook enregistré appelle `event.Cancel()`, la fermeture est empêchée.

## Propriétés des fenêtres

### SetTitle()

Définit le texte de la barre de titre de la fenêtre. Renvoie le récepteur pour permettre le chaînage.

```go
func (w *WebviewWindow) SetTitle(title string) Window
```

**Paramètres :**

- `title` - Nouveau titre de la fenêtre

**Exemple :**

```go
window.SetTitle("My Application - Document.txt")
```

### Name()

Renvoie l’identifiant de nom unique de la fenêtre.

```go
func (w *WebviewWindow) Name() string
```

**Exemple :**

```go
name := window.Name()
fmt.Println("Window name:", name)

// Retrieve window by name later
if w, ok := app.Window.GetByName(name); ok {
    w.Focus()
}
```

## Taille et position

### SetSize()

Définit les dimensions de la fenêtre en pixels. Renvoie le récepteur pour permettre le chaînage.

```go
func (w *WebviewWindow) SetSize(width, height int) Window
```

**Paramètres :**

- `width` - Largeur de la fenêtre en pixels
- `height` - Hauteur de la fenêtre en pixels

**Exemple :**

```go
window.SetSize(1024, 768)
```

### Size()

Renvoie les dimensions actuelles de la fenêtre.

```go
func (w *WebviewWindow) Size() (width, height int)
```

**Exemple :**

```go
width, height := window.Size()
fmt.Printf("Window is %dx%d\n", width, height)
```

### SetMinSize() / SetMaxSize()

Définit les dimensions minimales et maximales de la fenêtre. Ces deux méthodes renvoient le récepteur pour permettre le chaînage.

```go
func (w *WebviewWindow) SetMinSize(width, height int) Window
func (w *WebviewWindow) SetMaxSize(width, height int) Window
```

**Exemple :**

```go
// Prevent window from being too small
window.SetMinSize(800, 600)

// Prevent window from being too large
window.SetMaxSize(1920, 1080)
```

### SetPosition()

Définit la position de la fenêtre par rapport au coin supérieur gauche de l’écran.

```go
func (w *WebviewWindow) SetPosition(x, y int)
```

**Paramètres :**

- `x` - Position horizontale en pixels
- `y` - Position verticale en pixels

**Exemple :**

```go
// Position window at top-left
window.SetPosition(0, 0)

// Position window 100px from top-left
window.SetPosition(100, 100)
```

### Position()

Renvoie la position actuelle de la fenêtre.

```go
func (w *WebviewWindow) Position() (x, y int)
```

**Exemple :**

```go
x, y := window.Position()
fmt.Printf("Window is at (%d, %d)\n", x, y)
```

### Center()

Centre la fenêtre sur l’écran.

```go
func (w *WebviewWindow) Center()
```

**Exemple :**

```go
window := app.Window.New()
window.Center()
window.Show()
```

**Remarque :** la fenêtre est centrée sur le moniteur principal. Pour les configurations à plusieurs moniteurs, consultez les API d’écran.

### Focus()

Place la fenêtre au premier plan et lui donne le focus clavier.

```go
func (w *WebviewWindow) Focus()
```

**Exemple :**

```go
// Bring window to front
window.Focus()
```

## État de la fenêtre

### Minimise() / UnMinimise()

Réduit la fenêtre dans la barre des tâches ou le Dock, ou la restaure. `Minimise()` renvoie le récepteur pour permettre le chaînage ; `UnMinimise()` ne renvoie rien.

```go
func (w *WebviewWindow) Minimise() Window
func (w *WebviewWindow) UnMinimise()
```

**Exemple :**

```go
// Minimise window
window.Minimise()

// Restore from minimised state
window.UnMinimise()
```

### Maximise() / UnMaximise()

Agrandit la fenêtre pour qu’elle occupe tout l’écran ou lui rend sa taille précédente. `Maximise()` renvoie le récepteur pour permettre le chaînage ; `UnMaximise()` ne renvoie rien.

```go
func (w *WebviewWindow) Maximise() Window
func (w *WebviewWindow) UnMaximise()
```

**Exemple :**

```go
// Maximise window
window.Maximise()

// Restore to previous size
window.UnMaximise()
```

### Fullscreen() / UnFullscreen() / ToggleFullscreen()

Active ou quitte le mode plein écran. `Fullscreen()` renvoie le récepteur pour permettre le chaînage.

```go
func (w *WebviewWindow) Fullscreen() Window
func (w *WebviewWindow) UnFullscreen()
func (w *WebviewWindow) ToggleFullscreen()
```

**Exemple :**

```go
// Enter fullscreen
window.Fullscreen()

// Exit fullscreen
window.UnFullscreen()

// Or toggle
window.ToggleFullscreen()
```

Il n’existe aucune méthode `SetFullscreen(bool)`.

### IsMinimised() / IsMaximised() / IsFullscreen()

Vérifie l’état actuel de la fenêtre.

```go
func (w *WebviewWindow) IsMinimised() bool
func (w *WebviewWindow) IsMaximised() bool
func (w *WebviewWindow) IsFullscreen() bool
```

**Exemple :**

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

## Contenu de la fenêtre

### SetURL()

Accède à une URL précise dans la fenêtre. Renvoie le récepteur pour permettre le chaînage.

```go
func (w *WebviewWindow) SetURL(url string) Window
```

**Paramètres :**

- `url` — URL à laquelle accéder (peut être `http://wails.localhost/` pour les ressources intégrées)

**Exemple :**

```go
// Navigate to embedded page
window.SetURL("http://wails.localhost/settings.html")

// Navigate to external URL (if allowed)
window.SetURL("https://wails.io")
```

### SetHTML()

Définit directement le contenu de la fenêtre à partir d’une chaîne HTML. Renvoie le récepteur pour permettre le chaînage.

```go
func (w *WebviewWindow) SetHTML(html string) Window
```

**Paramètres :**

- `html` — contenu HTML à afficher

**Exemple :**

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

**Cas d’utilisation :**

- Génération de contenu dynamique
- Fenêtres simples sans processus de compilation du frontend
- Pages d’erreur ou écrans de démarrage

### Reload()

Recharge le contenu actuel de la fenêtre.

```go
func (w *WebviewWindow) Reload()
```

**Exemple :**

```go
// Reload current page
window.Reload()
```

**Remarque :** utile pendant le développement ou lorsque le contenu doit être actualisé.

## Événements de fenêtre

Wails fournit deux méthodes pour gérer les événements de fenêtre :

- **OnWindowEvent()** — écoute les événements de fenêtre (sans pouvoir les empêcher).
- **RegisterHook()** — intercepte les événements de fenêtre (et peut les empêcher en appelant `event.Cancel()`).

### OnWindowEvent()

Enregistre une fonction de rappel pour les événements de fenêtre. Renvoie une fonction de désabonnement.

```go
func (w *WebviewWindow) OnWindowEvent(
    eventType events.WindowEventType,
    callback func(event *WindowEvent),
) func()
```

**Exemple :**

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

**Événements de fenêtre courants :**

- `events.Common.WindowClosing` — la fenêtre est sur le point de se fermer
- `events.Common.WindowFocus` — la fenêtre a obtenu le focus
- `events.Common.WindowLostFocus` — la fenêtre a perdu le focus
- `events.Common.WindowDidMove` — la fenêtre a été déplacée
- `events.Common.WindowDidResize` — la fenêtre a été redimensionnée
- `events.Common.WindowMinimise` — la fenêtre a été réduite
- `events.Common.WindowMaximise` — la fenêtre a été agrandie
- `events.Common.WindowFullscreen` — la fenêtre est passée en plein écran
- `events.Common.WindowRuntimeReady` — l’environnement d’exécution intégré à la fenêtre a été initialisé

### RegisterHook()

Enregistre un hook pour les événements de fenêtre. Les hooks s’exécutent avant les écouteurs et peuvent empêcher l’événement en appelant `event.Cancel()`. Renvoie une fonction de désabonnement.

```go
func (w *WebviewWindow) RegisterHook(
    eventType events.WindowEventType,
    callback func(event *WindowEvent),
) func()
```

**Exemple — empêcher la fermeture de la fenêtre :**

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

**Exemple — enregistrer avant la fermeture :**

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

Émet un événement personnalisé vers le frontend de la fenêtre. Renvoie `true` si l’émission a été annulée par un hook.

```go
func (w *WebviewWindow) EmitEvent(name string, data ...any) bool
```

**Paramètres :**

- `name` — nom de l’événement
- `data` — données facultatives à envoyer avec l’événement

**Exemple :**

```go
// Send data to specific window
window.EmitEvent("data-updated", map[string]any{
    "count":  42,
    "status": "success",
})
```

**Frontend (JavaScript) :**

```javascript
import { Events } from '@wailsio/runtime'

Events.On('data-updated', (data) => {
    console.log('Count:', data.count)
    console.log('Status:', data.status)
})
```

## Autres méthodes

### SetEnabled()

Active ou désactive l’interaction de l’utilisateur avec la fenêtre.

```go
func (w *WebviewWindow) SetEnabled(enabled bool)
```

**Exemple :**

```go
// Disable window during long operation
window.SetEnabled(false)

// Perform operation
performLongOperation()

// Re-enable window
window.SetEnabled(true)
```

### SetBackgroundColour()

Définit la couleur d’arrière-plan de la fenêtre (affichée avant le chargement du contenu). Renvoie le récepteur pour permettre le chaînage.

```go
func (w *WebviewWindow) SetBackgroundColour(colour RGBA) Window
```

`RGBA` est `application.RGBA{Red, Green, Blue, Alpha uint8}`. Utilisez les fonctions auxiliaires `application.NewRGB(r, g, b)` (alpha 255) ou `application.NewRGBA(r, g, b, a)`.

**Exemple :**

```go
// White background
window.SetBackgroundColour(application.NewRGB(255, 255, 255))

// Dark background with full alpha
window.SetBackgroundColour(application.NewRGBA(30, 30, 30, 255))
```

### SetResizable()

Détermine si l’utilisateur peut redimensionner la fenêtre. Renvoie le récepteur pour permettre le chaînage.

```go
func (w *WebviewWindow) SetResizable(resizable bool) Window
```

**Exemple :**

```go
// Make window fixed size
window.SetResizable(false)
```

### SetAlwaysOnTop()

Détermine si la fenêtre reste au-dessus des autres fenêtres. Renvoie le récepteur pour permettre le chaînage.

```go
func (w *WebviewWindow) SetAlwaysOnTop(alwaysOnTop bool) Window
```

**Exemple :**

```go
// Keep window on top
window.SetAlwaysOnTop(true)
```

### Print()

Ouvre la boîte de dialogue d’impression native pour le contenu de la fenêtre.

```go
func (w *WebviewWindow) Print() error
```

**Renvoie :** une erreur si l’impression échoue.

**Exemple :**

```go
if err := window.Print(); err != nil {
    log.Println("Print failed:", err)
}
```

### AttachModal()

Attache une seconde fenêtre en tant que fenêtre modale de type feuille.

```go
func (w *WebviewWindow) AttachModal(modalWindow Window)
```

**Paramètres :**

- `modalWindow` - Fenêtre à attacher en tant que fenêtre modale

**Prise en charge selon la plateforme :**

- **macOS** : prise en charge complète (affichage sous forme de feuille)
- **Windows** : non pris en charge
- **Linux** : non pris en charge

**Exemple :**

```go
modalWindow := app.Window.New()
window.AttachModal(modalWindow)
```

## Options propres à chaque plateforme

### Linux

Sous Linux, les fenêtres prennent en charge les options propres à la plateforme suivantes via `LinuxWindow` :

#### MenuStyle

Contrôle le mode d’affichage du menu de l’application. Cette option est disponible dans la version GTK4 par défaut et est ignorée dans les anciennes versions `-tags gtk3`.

| Valeur | Description |
| --- | --- |
| `LinuxMenuStyleMenuBar` | Barre de menus classique sous la barre de titre (par défaut) |
| `LinuxMenuStylePrimaryMenu` | Bouton du menu principal dans la barre d’en-tête (style GNOME) |

**Exemple :**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "My Application",
    Linux: application.LinuxWindow{
        MenuStyle: application.LinuxMenuStylePrimaryMenu,
    },
})
window.SetMenu(menu)
```

**Remarque :** le style de menu principal affiche un bouton de menu hamburger (☰) dans la barre d’en-tête, conformément aux recommandations GNOME relatives aux interfaces humaines. Ce style est recommandé pour les applications GNOME modernes.

## Exemple complet

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
