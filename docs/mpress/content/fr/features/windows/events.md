---
title: "Événements de fenêtre"
description: "Gérer les événements du cycle de vie et de changement d’état des fenêtres"
slug: "features/windows/events"
sourcePath: "features/windows/events.md"
---

## Événements de fenêtre

Wails émet les événements du cycle de vie et de changement d’état des fenêtres par l’intermédiaire d’une API unique : `OnWindowEvent` pour les écouteurs passifs et `RegisterHook` pour les hooks annulables capables d’empêcher l’action par défaut.

```go
import "github.com/wailsapp/wails/v3/pkg/events"

// Listener — observes the event, cannot cancel it.
window.OnWindowEvent(events.Common.WindowDidResize, func(e *application.WindowEvent) { /* ... */ })

// Hook — runs before listeners; can call e.Cancel() to suppress the default action.
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    if hasUnsavedChanges() {
        e.Cancel() // prevent the window from closing
    }
})
```

Ces deux appels renvoient une fonction `unsubscribe func()` que vous pouvez appeler pour supprimer le gestionnaire.

Les événements multiplateformes se trouvent dans `events.Common.*`. Les événements propres à chaque plateforme se trouvent dans `events.Mac.*`, `events.Windows.*` et `events.Linux.*`. La liste complète est générée dans `v3/pkg/events/events.go`.

## Événements du cycle de vie

### Création d’une fenêtre

Exécutez une fonction de rappel à chaque création de fenêtre avec `app.Window.OnCreate` :

```go
app.Window.OnCreate(func(window application.Window) {
    fmt.Printf("Window created: %s (ID: %d)\n", window.Name(), window.ID())

    // Configure all new windows
    window.SetMinSize(400, 300)

    // Register a hook for this window
    window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
        if !confirmClose() {
            e.Cancel()
        }
    })
})
```

La fonction de rappel reçoit `application.Window` (une interface). La fonction de rappel `OnCreate` est appelée une fois par fenêtre, après l’initialisation de l’environnement d’exécution de celle-ci.

### WindowClosing

Se déclenche lorsque l’utilisateur tente de fermer la fenêtre (en cliquant sur X, avec ⌘W, Alt+F4, etc.).

Utilisez un **hook** (`RegisterHook`) : seuls les hooks peuvent annuler la fermeture. Les écouteurs observent l’événement, mais ne peuvent pas l’empêcher.

```go
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    if hasUnsavedChanges() {
        // Cancel the close — the window stays open.
        e.Cancel()
    }
})
```

**Important :**

- `WindowClosing` est émis lors des tentatives de fermeture déclenchées par l’utilisateur.
- Une fonction de rappel `RegisterHook` peut appeler `e.Cancel()` pour maintenir la fenêtre ouverte.
- Les fonctions de rappel `OnWindowEvent` associées à `WindowClosing` sont des observateurs : elles sont exécutées, mais ne peuvent pas annuler l’action.

**Fermeture par programmation :** appelez `window.Close()`. Il n’existe pas de `window.Destroy()`.

### WindowRuntimeReady

Se déclenche une fois l’initialisation de l’environnement d’exécution intégré à la fenêtre terminée ; vous pouvez alors appeler sans risque le contexte JS de la fenêtre :

```go
window.OnWindowEvent(events.Common.WindowRuntimeReady, func(e *application.WindowEvent) {
    window.EmitEvent("app-ready", nil)
})
```

## Événements de focus

### WindowFocus

Appelé lorsque la fenêtre reçoit le focus :

```go
window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
    fmt.Println("Window gained focus")
    updateTitleBar(true)
    app.Event.Emit("window-focused", window.ID())
})
```

### WindowLostFocus

Appelé lorsque la fenêtre perd le focus :

```go
window.OnWindowEvent(events.Common.WindowLostFocus, func(e *application.WindowEvent) {
    fmt.Println("Window lost focus")
    updateTitleBar(false)
    saveCurrentState()
})
```

**Exemple : interface utilisateur sensible au focus :**

```go
type FocusAwareWindow struct {
    window  *application.WebviewWindow
    focused bool
}

func (fw *FocusAwareWindow) Setup() {
    fw.window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
        fw.focused = true
        fw.updateAppearance()
    })

    fw.window.OnWindowEvent(events.Common.WindowLostFocus, func(e *application.WindowEvent) {
        fw.focused = false
        fw.updateAppearance()
    })
}

func (fw *FocusAwareWindow) updateAppearance() {
    if fw.focused {
        fw.window.EmitEvent("update-theme", "active")
    } else {
        fw.window.EmitEvent("update-theme", "inactive")
    }
}
```

## Événements de changement d’état

### WindowMinimise / WindowUnMinimise

```go
window.OnWindowEvent(events.Common.WindowMinimise, func(e *application.WindowEvent) {
    pauseRendering()
    saveWindowState()
})

window.OnWindowEvent(events.Common.WindowUnMinimise, func(e *application.WindowEvent) {
    resumeRendering()
    refreshContent()
})
```

### WindowMaximise / WindowUnMaximise

```go
window.OnWindowEvent(events.Common.WindowMaximise, func(e *application.WindowEvent) {
    window.EmitEvent("layout-mode", "maximised")
})

window.OnWindowEvent(events.Common.WindowUnMaximise, func(e *application.WindowEvent) {
    window.EmitEvent("layout-mode", "normal")
})
```

### WindowFullscreen / WindowUnFullscreen

```go
window.OnWindowEvent(events.Common.WindowFullscreen, func(e *application.WindowEvent) {
    window.EmitEvent("chrome-visibility", false)
    window.EmitEvent("layout-mode", "fullscreen")
})

window.OnWindowEvent(events.Common.WindowUnFullscreen, func(e *application.WindowEvent) {
    window.EmitEvent("chrome-visibility", true)
    window.EmitEvent("layout-mode", "normal")
})
```

Activez ou quittez le mode plein écran par programmation avec `window.Fullscreen()` / `window.UnFullscreen()` / `window.ToggleFullscreen()` ; il n’existe pas de `SetFullscreen(bool)`. Consultez l’état avec `window.IsFullscreen() bool`.

## Événements de position et de taille

### WindowDidMove

```go
window.OnWindowEvent(events.Common.WindowDidMove, func(e *application.WindowEvent) {
    x, y := window.Position()
    fmt.Printf("Window moved to: %d, %d\n", x, y)
    saveWindowPosition(x, y)
})
```

### WindowDidResize

```go
window.OnWindowEvent(events.Common.WindowDidResize, func(e *application.WindowEvent) {
    width, height := window.Size()
    fmt.Printf("Window resized to: %dx%d\n", width, height)
    saveWindowSize(width, height)
    window.EmitEvent("window-size", map[string]int{
        "width":  width,
        "height": height,
    })
})
```

`WindowDidResize` et `WindowDidMove` ne fournissent pas de coordonnées dans l’événement lui-même : interrogez la fenêtre avec `window.Size()` / `window.Position()` dans la fonction de rappel.

**Exemple : mise en page adaptative :**

```go
window.OnWindowEvent(events.Common.WindowDidResize, func(e *application.WindowEvent) {
    width, _ := window.Size()
    var layout string
    switch {
    case width < 600:
        layout = "compact"
    case width < 1200:
        layout = "normal"
    default:
        layout = "wide"
    }
    window.EmitEvent("layout-changed", layout)
})
```

## Exemple complet

Une fenêtre prête pour la production avec une gestion complète des événements :

```go
package main

import (
    "encoding/json"
    "fmt"
    "os"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)

type WindowState struct {
    X          int  `json:"x"`
    Y          int  `json:"y"`
    Width      int  `json:"width"`
    Height     int  `json:"height"`
    Maximised  bool `json:"maximised"`
    Fullscreen bool `json:"fullscreen"`
}

type ManagedWindow struct {
    app    *application.App
    window *application.WebviewWindow
    state  WindowState
    dirty  bool
}

func main() {
    app := application.New(application.Options{
        Name: "Event Demo",
    })

    mw := &ManagedWindow{app: app}
    mw.CreateWindow()
    mw.LoadState()
    mw.SetupEventHandlers()

    app.Run()
}

func (mw *ManagedWindow) CreateWindow() {
    mw.window = mw.app.Window.NewWithOptions(application.WebviewWindowOptions{
        Name:   "main",
        Title:  "Event Demo",
        Width:  800,
        Height: 600,
    })
}

func (mw *ManagedWindow) SetupEventHandlers() {
    // Focus events
    mw.window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
        mw.window.EmitEvent("focus-state", true)
    })

    mw.window.OnWindowEvent(events.Common.WindowLostFocus, func(e *application.WindowEvent) {
        mw.window.EmitEvent("focus-state", false)
    })

    // State change events
    mw.window.OnWindowEvent(events.Common.WindowMinimise, func(e *application.WindowEvent) {
        mw.SaveState()
    })

    mw.window.OnWindowEvent(events.Common.WindowMaximise, func(e *application.WindowEvent) {
        mw.state.Maximised = true
        mw.dirty = true
    })

    mw.window.OnWindowEvent(events.Common.WindowUnMaximise, func(e *application.WindowEvent) {
        mw.state.Maximised = false
        mw.dirty = true
    })

    mw.window.OnWindowEvent(events.Common.WindowFullscreen, func(e *application.WindowEvent) {
        mw.state.Fullscreen = true
        mw.dirty = true
    })

    mw.window.OnWindowEvent(events.Common.WindowUnFullscreen, func(e *application.WindowEvent) {
        mw.state.Fullscreen = false
        mw.dirty = true
    })

    // Position and size events
    mw.window.OnWindowEvent(events.Common.WindowDidMove, func(e *application.WindowEvent) {
        mw.state.X, mw.state.Y = mw.window.Position()
        mw.dirty = true
    })

    mw.window.OnWindowEvent(events.Common.WindowDidResize, func(e *application.WindowEvent) {
        mw.state.Width, mw.state.Height = mw.window.Size()
        mw.dirty = true
    })

    // Cancellable close — use a hook, not a listener.
    mw.window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
        if mw.dirty {
            mw.SaveState()
        }
    })
}

func (mw *ManagedWindow) LoadState() {
    data, err := os.ReadFile("window-state.json")
    if err != nil {
        return
    }

    if err := json.Unmarshal(data, &mw.state); err != nil {
        return
    }

    // Restore window state
    mw.window.SetPosition(mw.state.X, mw.state.Y)
    mw.window.SetSize(mw.state.Width, mw.state.Height)

    if mw.state.Maximised {
        mw.window.Maximise()
    }

    if mw.state.Fullscreen {
        mw.window.Fullscreen()
    }
}

func (mw *ManagedWindow) SaveState() {
    data, err := json.Marshal(mw.state)
    if err != nil {
        return
    }

    os.WriteFile("window-state.json", data, 0644)
    mw.dirty = false

    fmt.Println("Window state saved")
}
```

## Coordination des événements

### Événements entre fenêtres

Coordonnez plusieurs fenêtres au moyen du bus d’événements de l’application :

```go
// In main window
mainWindow.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
    app.Event.Emit("main-window-focused", nil)
})

// In other windows
app.Event.On("main-window-focused", func(event *application.CustomEvent) {
    updateRelativeToMain()
})
```

### Chaînes d’événements

```go
window.OnWindowEvent(events.Common.WindowMaximise, func(e *application.WindowEvent) {
    saveWindowState()
    window.EmitEvent("layout-changed", "maximised")
    app.Event.Emit("window-maximised", window.ID())
})
```

### Événements temporisés

```go
var resizeTimer *time.Timer

window.OnWindowEvent(events.Common.WindowDidResize, func(e *application.WindowEvent) {
    if resizeTimer != nil {
        resizeTimer.Stop()
    }

    resizeTimer = time.AfterFunc(500*time.Millisecond, func() {
        w, h := window.Size()
        saveWindowSize(w, h)
    })
})
```

## Bonnes pratiques

### ✅ À faire

- **Utilisez des hooks pour l’annulation** : seules les fonctions de rappel `RegisterHook` peuvent appeler `e.Cancel()`.
- **Enregistrez l’état lors de la fermeture** : restaurez la position et la taille de la fenêtre au prochain lancement.
- **Temporisez les événements fréquents** : `WindowDidResize` et `WindowDidMove` se déclenchent rapidement.
- **Gérez les changements de focus** : mettez à jour l’interface utilisateur de façon appropriée.
- **Coordonnez les fenêtres au moyen d’événements** : utilisez `app.Event.Emit` pour les messages entre fenêtres.
- **Désabonnez-vous** lorsque les gestionnaires ne sont plus nécessaires : `OnWindowEvent` et `RegisterHook` renvoient tous deux une fonction `func()` de désabonnement.

### ❌ À ne pas faire

- **Ne bloquez pas les gestionnaires d’événements** : veillez à ce qu’ils s’exécutent rapidement.
- **N’essayez pas d’annuler l’action depuis `OnWindowEvent`** : utilisez `RegisterHook`.
- **N’utilisez pas `window.Destroy()`** : cette fonction n’existe pas ; utilisez `window.Close()`.
- **N’enregistrez pas l’état à chaque événement** : appliquez d’abord une temporisation.
- **Ne supposez pas que l’événement contient des coordonnées** : appelez `window.Position()` / `window.Size()`.

## Dépannage

### Le hook WindowClosing n’empêche pas la fermeture

**Cause :** vous avez utilisé `OnWindowEvent` au lieu de `RegisterHook`.

**Solution :** seules les fonctions de rappel `RegisterHook` peuvent appeler `e.Cancel()`. Les fonctions de rappel `OnWindowEvent` observent l’événement ; elles ne peuvent pas l’annuler.

```go
// ❌ Cannot cancel — this is a listener.
window.OnWindowEvent(events.Common.WindowClosing, func(e *application.WindowEvent) {
    e.Cancel() // no effect
})

// ✅ Can cancel — this is a hook.
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    e.Cancel()
})
```

### Les événements ne se déclenchent pas

**Cause :** le gestionnaire a été enregistré après le déclenchement de l’événement.

**Solution :** enregistrez les gestionnaires immédiatement après la création de la fenêtre (ou dans le rappel `app.Window.OnCreate`).

```go
window := app.Window.New()
window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) { /* ... */ })
```

### Fuites de mémoire

**Cause :** des gestionnaires à longue durée de vie sont conservés par un état à courte durée de vie.

**Solution :** conservez la fonction de désabonnement `func()` renvoyée par `OnWindowEvent`/`RegisterHook` et appelez-la lors du nettoyage.

```go
unsub := window.OnWindowEvent(events.Common.WindowDidResize, handler)
// ...later:
unsub()
```

## Étapes suivantes

**Principes de base des fenêtres** – Découvrez les fondamentaux de la gestion des fenêtres [En savoir plus →](/features/windows/basics/)

**Fenêtres multiples** – Modèles pour les applications multifenêtres [En savoir plus →](/features/windows/multiple/)

**Système d’événements** – Explorez en détail le système d’événements [En savoir plus →](/features/events/system/)

**Cycle de vie de l’application** – Comprenez le cycle de vie de l’application [En savoir plus →](/concepts/lifecycle/)

---

**Des questions ?** Posez-les sur [Discord](https://discord.gg/JDdSxwjhGf) ou consultez les [exemples](https://github.com/wailsapp/wails/tree/master/v3/examples).
