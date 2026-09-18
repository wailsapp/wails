---
title: "API Events"
description: "Référence complète de l’API Events"
slug: "reference/events"
sourcePath: "reference/events.md"
---

## Présentation

L’API Events fournit des méthodes permettant d’émettre et d’écouter des événements afin d’assurer la communication entre les différentes parties de votre application.

**Types d’événements :**

- **Événements d’application** — Événements du cycle de vie de l’application (démarrage, arrêt)
- **Événements de fenêtre** — Changements d’état de la fenêtre (obtention ou perte du focus, redimensionnement)
- **Événements personnalisés** — Événements définis par l’utilisateur pour les communications propres à l’application

**Modes de communication :**

- **De Go vers le frontend** — Émettez des événements depuis Go et écoutez-les en JavaScript
- **Du frontend vers Go** — Pas directement (utilisez plutôt les liaisons de services)
- **D’un frontend à un autre** — Via Go ou les événements locaux du runtime
- **D’une fenêtre à une autre** — Ciblez des fenêtres précises ou diffusez l’événement à toutes les fenêtres

## Méthodes relatives aux événements (Go)

### app.Event.Emit()

Émet un événement personnalisé vers toutes les fenêtres. Renvoie `true` si un hook a annulé l’émission.

```go
func (em *EventManager) Emit(name string, data ...any) bool
```

**Paramètres :**

- `name` — Nom de l’événement
- `data` — Données facultatives à envoyer avec l’événement

**Exemple :**

```go
// Emit simple event
app.Event.Emit("user-logged-in")

// Emit with data
app.Event.Emit("data-updated", map[string]interface{}{
    "count": 42,
    "status": "success",
})

// Emit multiple values
app.Event.Emit("progress", 75, "Processing files...")
```

### app.Event.On()

Écoute les événements personnalisés dans Go.

```go
func (em *EventManager) On(name string, callback func(*CustomEvent)) func()
```

**Paramètres :**

- `name` — Nom de l’événement à écouter
- `callback` — Fonction appelée lors de l’émission de l’événement

**Valeur renvoyée :** fonction de nettoyage permettant de supprimer l’écouteur d’événement

**Exemple :**

```go
// Listen for events
cleanup := app.Event.On("user-action", func(e *application.CustomEvent) {
    data := e.Data.(map[string]interface{})
    action := data["action"].(string)
    app.Logger.Info("User action", "action", action)
})

// Later, remove listener
cleanup()
```

### Événements propres à une fenêtre

Émettez des événements vers une fenêtre précise :

```go
// Emit to specific window
window.EmitEvent("notification", "Hello from Go!")

// Emit to all windows
app.Event.Emit("global-update", data)
```

## Méthodes relatives aux événements (frontend)

### On()

Écoute les événements provenant de Go.

```javascript
import { Events } from '@wailsio/runtime'

Events.On(eventName, callback)
```

**Paramètres :**

- `eventName` — Nom de l’événement à écouter
- `callback` — Fonction appelée à la réception de l’événement

**Valeur renvoyée :** fonction de nettoyage

**Exemple :**

```javascript
import { Events } from '@wailsio/runtime'

// Listen for events
const cleanup = Events.On('data-updated', (data) => {
    console.log('Count:', data.count)
    console.log('Status:', data.status)
    updateUI(data)
})

// Later, remove listener
cleanup()
```

### Once()

Attend une seule occurrence d’un événement.

```javascript
import { Events } from '@wailsio/runtime'

Events.Once(eventName, callback)
```

**Exemple :**

```javascript
import { Events } from '@wailsio/runtime'

// Listen for first occurrence only
Events.Once('initialization-complete', (data) => {
    console.log('App initialized!', data)
    // This will only fire once
})
```

### Off()

Supprime tous les écouteurs d’un ou de plusieurs événements. `Off` accepte un nombre variable de chaînes correspondant aux noms d’événements ; cette fonction n’accepte **pas** de fonction de rappel. Pour supprimer un seul écouteur, conservez la fonction de désabonnement renvoyée par `Events.On(...)` et appelez-la.

```typescript
import { Events } from '@wailsio/runtime'

Events.Off(...eventNames: string[]): void
```

**Exemple :**

```javascript
import { Events } from '@wailsio/runtime'

// Preferred: keep the unsubscribe fn from On()
const unsubscribe = Events.On('my-event', (data) => {
    console.log('Event received:', data)
})

// Later — remove just this listener
unsubscribe()

// Or: remove every listener for one or more events
Events.Off('my-event', 'another-event')
```

### OffAll()

Supprime **tous** les écouteurs d’événements. N’accepte aucun argument.

```typescript
import { Events } from '@wailsio/runtime'

Events.OffAll(): void
```

**Exemple :**

```javascript
import { Events } from '@wailsio/runtime'

// Remove all listeners — typically used during teardown.
Events.OffAll()
```

### OnMultiple()

Écoute un événement jusqu’à `max` occurrences, puis se désabonne automatiquement.

```typescript
Events.OnMultiple(eventName: string, callback, max: number): () => void
```

**Exemple :**

```javascript
Events.OnMultiple('progress', (data) => {
    console.log('progress', data)
}, 5)
```

## Événements d’application

### app.Event.OnApplicationEvent()

Écoute les événements du cycle de vie de l’application.

```go
func (em *EventManager) OnApplicationEvent(
    eventType events.ApplicationEventType,
    callback func(*ApplicationEvent),
) func()
```

Les constantes d’événements se trouvent dans le package `events` : `events.Common.*` pour les événements multiplateformes, et `events.Mac.*`, `events.Windows.*` et `events.Linux.*` pour ceux propres à chaque plateforme.

**Événements d’application courants :**

- `events.Common.ApplicationStarted` — Le lancement de l’application est terminé.
- `events.Common.ThemeChanged` — Le thème du système est passé du mode clair au mode sombre, ou inversement.
- `events.Common.ApplicationOpenedWithFile` — L’application a été lancée par l’intermédiaire d’une association de fichiers.
- `events.Common.ApplicationLaunchedWithUrl` — L’application a été lancée par l’intermédiaire d’un schéma d’URL.

Il n’existe **aucune** constante d’événement générique « arrêt de l’application » : enregistrez les opérations de nettoyage à effectuer lors de l’arrêt au moyen de `application.Options.OnShutdown` ou de `app.OnShutdown(func())`.

**Exemple :**

```go
import "github.com/wailsapp/wails/v3/pkg/events"

// Handle application startup
app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(e *application.ApplicationEvent) {
    app.Logger.Info("Application started")
})

// React to theme changes
app.Event.OnApplicationEvent(events.Common.ThemeChanged, func(e *application.ApplicationEvent) {
    if e.Context().IsDarkMode() {
        app.Logger.Info("Dark mode enabled")
    }
})

// Shutdown cleanup is configured on the application, not as an event:
app.OnShutdown(func() {
    database.Close()
    saveSettings()
})
```

## Événements de fenêtre

### OnWindowEvent()

Écoute les événements propres à une fenêtre.

```go
func (w *WebviewWindow) OnWindowEvent(
    eventType events.WindowEventType,
    callback func(*WindowEvent),
) func()
```

Les constantes d’événements de fenêtre se trouvent dans le package `events` : `events.Common.*` pour les événements multiplateformes (et `events.Mac.*`, `events.Windows.*` et `events.Linux.*` pour ceux propres à chaque plateforme).

**Événements de fenêtre courants :**

- `events.Common.WindowFocus` — La fenêtre a obtenu le focus.
- `events.Common.WindowLostFocus` - La fenêtre a perdu le focus.
- `events.Common.WindowClosing` - La fenêtre est sur le point de se fermer (fermeture annulable via `RegisterHook`).
- `events.Common.WindowDidResize` - La fenêtre a été redimensionnée.
- `events.Common.WindowDidMove` - La fenêtre a été déplacée.
- `events.Common.WindowMinimise` / `WindowUnMinimise` / `WindowMaximise` / `WindowUnMaximise` / `WindowFullscreen` / `WindowUnFullscreen`.
- `events.Common.WindowRuntimeReady` - Les événements peuvent être émis en toute sécurité vers le runtime de la fenêtre.

**Exemple :**

```go
import "github.com/wailsapp/wails/v3/pkg/events"

// Handle window focus
window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
    app.Logger.Info("Window focused")
})

// Handle window resize
window.OnWindowEvent(events.Common.WindowDidResize, func(e *application.WindowEvent) {
    width, height := window.Size()
    app.Logger.Info("Window resized", "width", width, "height", height)
})
```

## Modèles courants

Ces modèles présentent des approches éprouvées pour utiliser les événements dans des applications réelles. Chacun répond à un besoin de communication précis entre votre backend Go et votre frontend, afin de vous aider à créer des applications réactives et bien structurées.

### Modèle requête-réponse

Utilisez ce modèle pour signaler au frontend la fin d’opérations du backend, par exemple après la récupération de données, le traitement de fichiers ou l’exécution de tâches en arrière-plan. La liaison de service renvoie directement les données, tandis que les événements fournissent des notifications supplémentaires pour mettre à jour l’interface utilisateur, par exemple en affichant des notifications toast ou en actualisant des listes.

**Go :**

```go
// Service method
type DataService struct {
    app *application.App
}

func (s *DataService) FetchData(query string) ([]Item, error) {
    items := fetchFromDatabase(query)

    // Emit event when done
    s.app.Event.Emit("data-fetched", map[string]interface{}{
        "query": query,
        "count": len(items),
    })

    return items, nil
}
```

**JavaScript :**

```javascript
import { FetchData } from './bindings/DataService'
import { Events } from '@wailsio/runtime'

// Listen for completion event
Events.On('data-fetched', (data) => {
    console.log(`Fetched ${data.count} items for query: ${data.query}`)
    showNotification(`Found ${data.count} results`)
})

// Call service method
const items = await FetchData("search term")
displayItems(items)
```

### Mises à jour de la progression

Ce modèle convient parfaitement aux opérations longues, telles que le téléversement de fichiers, le traitement par lots, l’importation de grands volumes de données ou l’encodage vidéo. Émettez des événements de progression pendant l’opération afin de mettre à jour les barres de progression, le texte d’état ou les indicateurs d’étape de l’interface utilisateur, et ainsi fournir aux utilisateurs un retour en temps réel.

**Go :**

```go
func (s *Service) ProcessFiles(files []string) error {
    total := len(files)

    for i, file := range files {
        // Process file
        processFile(file)

        // Emit progress event
        s.app.Event.Emit("progress", map[string]interface{}{
            "current": i + 1,
            "total":   total,
            "percent": float64(i+1) / float64(total) * 100,
            "file":    file,
        })
    }

    s.app.Event.Emit("processing-complete")
    return nil
}
```

**JavaScript :**

```javascript
import { Events } from '@wailsio/runtime'

// Update progress bar
Events.On('progress', (data) => {
    progressBar.style.width = `${data.percent}%`
    statusText.textContent = `Processing ${data.file}... (${data.current}/${data.total})`
})

// Handle completion
Events.Once('processing-complete', () => {
    progressBar.style.width = '100%'
    statusText.textContent = 'Complete!'
    setTimeout(() => hideProgressBar(), 2000)
})
```

### Communication entre plusieurs fenêtres

Ce modèle est idéal pour les applications comportant plusieurs fenêtres, comme des panneaux de paramètres, des tableaux de bord ou des visionneuses de documents. Diffusez des événements pour synchroniser l’état entre toutes les fenêtres (changements de thème, préférences utilisateur), ou envoyez des événements ciblés à des fenêtres précises pour effectuer des mises à jour qui leur sont propres.

**Go :**

```go
// Broadcast to all windows
app.Event.Emit("theme-changed", "dark")

// Send to specific window
preferencesWindow.EmitEvent("settings-updated", settings)

// Per-window listener — receive on the global event bus, but
// gate by the event's source-window name (set automatically when
// a window emits via window.EmitEvent).
app.Event.On("request-data", func(e *application.CustomEvent) {
    if e.Sender != window1.Name() {
        return
    }
    window1.EmitEvent("data-response", data)
})
```

**JavaScript :**

```javascript
import { Events } from '@wailsio/runtime'

// Listen in any window
Events.On('theme-changed', (theme) => {
    document.body.className = theme
})
```

### Synchronisation de l’état

Utilisez ce modèle lorsque vous devez synchroniser l’état du frontend et celui du backend, par exemple pour les sessions utilisateur, la configuration de l’application ou les fonctionnalités collaboratives. Lorsque l’état change dans le backend, émettez des événements pour mettre à jour tous les frontends connectés et garantir ainsi la cohérence dans toute l’application.

**Go :**

```go
type StateService struct {
    app   *application.App
    state map[string]interface{}
    mu    sync.RWMutex
}

func (s *StateService) UpdateState(key string, value interface{}) {
    s.mu.Lock()
    s.state[key] = value
    s.mu.Unlock()

    // Notify all windows
    s.app.Event.Emit("state-updated", map[string]interface{}{
        "key":   key,
        "value": value,
    })
}

func (s *StateService) GetState(key string) interface{} {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.state[key]
}
```

**JavaScript :**

```javascript
import { Events } from '@wailsio/runtime'
import { GetState } from './bindings/StateService'

// Keep local state in sync
let localState = {}

Events.On('state-updated', async (data) => {
    localState[data.key] = data.value
    updateUI(data.key, data.value)
})

// Initialize state
const initialState = await GetState("all")
localState = initialState
```

### Notifications pilotées par les événements

Ce modèle est particulièrement adapté à l’affichage de retours destinés à l’utilisateur, comme des confirmations de réussite, des alertes d’erreur ou des messages d’information. Au lieu d’appeler directement le code de l’interface utilisateur depuis les services, émettez des événements de notification que le frontend traite de manière cohérente. Vous pourrez ainsi modifier facilement le style des notifications ou ajouter des fonctionnalités telles qu’un historique des notifications.

**Go :**

```go
type NotificationService struct {
    app *application.App
}

func (s *NotificationService) Success(message string) {
    s.app.Event.Emit("notification", map[string]interface{}{
        "type":    "success",
        "message": message,
    })
}

func (s *NotificationService) Error(message string) {
    s.app.Event.Emit("notification", map[string]interface{}{
        "type":    "error",
        "message": message,
    })
}

func (s *NotificationService) Info(message string) {
    s.app.Event.Emit("notification", map[string]interface{}{
        "type":    "info",
        "message": message,
    })
}
```

**JavaScript :**

```javascript
import { Events } from '@wailsio/runtime'

// Unified notification handler
Events.On('notification', (data) => {
    const toast = document.createElement('div')
    toast.className = `toast toast-${data.type}`
    toast.textContent = data.message

    document.body.appendChild(toast)

    setTimeout(() => {
        toast.classList.add('fade-out')
        setTimeout(() => toast.remove(), 300)
    }, 3000)
})
```

## Exemple complet

**Go :**

```go
package main

import (
    "sync"
    "time"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)

type EventDemoService struct {
    app *application.App
    mu  sync.Mutex
}

func NewEventDemoService(app *application.App) *EventDemoService {
    service := &EventDemoService{app: app}

    // Listen for custom events
    app.Event.On("user-action", func(e *application.CustomEvent) {
        data := e.Data.(map[string]interface{})
        app.Logger.Info("User action received", "data", data)
    })

    return service
}

func (s *EventDemoService) StartLongTask() {
    go func() {
        s.app.Event.Emit("task-started")

        for i := 1; i <= 10; i++ {
            time.Sleep(500 * time.Millisecond)

            s.app.Event.Emit("task-progress", map[string]interface{}{
                "step":    i,
                "total":   10,
                "percent": i * 10,
            })
        }

        s.app.Event.Emit("task-completed", map[string]interface{}{
            "message": "Task finished successfully!",
        })
    }()
}

func (s *EventDemoService) BroadcastMessage(message string) {
    s.app.Event.Emit("broadcast", message)
}

func main() {
    app := application.New(application.Options{
        Name: "Event Demo",
    })

    // Handle application lifecycle
    app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(e *application.ApplicationEvent) {
        app.Logger.Info("Application started!")
    })

    app.OnShutdown(func() {
        app.Logger.Info("Application shutting down...")
    })

    // Register service (RegisterService returns nothing).
    service := NewEventDemoService(app)
    app.RegisterService(application.NewService(service))

    // Create window
    window := app.Window.New()

    // Handle window events
    window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
        window.EmitEvent("window-state", "focused")
    })

    window.OnWindowEvent(events.Common.WindowLostFocus, func(e *application.WindowEvent) {
        window.EmitEvent("window-state", "blurred")
    })

    window.Show()
    app.Run()
}
```

**JavaScript :**

```javascript
import { Events } from '@wailsio/runtime'
import { StartLongTask, BroadcastMessage } from './bindings/EventDemoService'

// Task events
Events.On('task-started', () => {
    console.log('Task started...')
    document.getElementById('status').textContent = 'Running...'
})

Events.On('task-progress', (data) => {
    const progressBar = document.getElementById('progress')
    progressBar.style.width = `${data.percent}%`
    console.log(`Step ${data.step} of ${data.total}`)
})

Events.Once('task-completed', (data) => {
    console.log('Task completed!', data.message)
    document.getElementById('status').textContent = data.message
})

// Broadcast events
Events.On('broadcast', (message) => {
    console.log('Broadcast:', message)
    alert(message)
})

// Window state events
Events.On('window-state', (state) => {
    console.log('Window is now:', state)
    document.body.dataset.windowState = state
})

// Trigger long task
document.getElementById('startTask').addEventListener('click', async () => {
    await StartLongTask()
})

// Send broadcast
document.getElementById('broadcast').addEventListener('click', async () => {
    const message = document.getElementById('message').value
    await BroadcastMessage(message)
})
```

## Événements intégrés

Wails fournit des événements système intégrés pour le cycle de vie de l’application et des fenêtres. Ces événements sont émis automatiquement par le framework.

### Événements communs et événements natifs de la plateforme

Wails fournit deux types d’événements système :

Les **événements communs** (`events.Common.*`) sont des abstractions multiplateformes qui fonctionnent de manière cohérente sous macOS, Windows et Linux. Pour assurer une portabilité maximale, utilisez ces événements dans votre application.

Les **événements natifs de la plateforme** (`events.Mac.*`, `events.Windows.*`, `events.Linux.*`) sont les événements sous-jacents propres au système d’exploitation à partir desquels les événements communs sont mis en correspondance. Ils donnent accès aux comportements et aux cas limites propres à chaque plateforme.

**Fonctionnement :**

```go
import "github.com/wailsapp/wails/v3/pkg/events"

// ✅ RECOMMENDED: Use Common Events for cross-platform code
window.OnWindowEvent(events.Common.WindowClosing, func(e *application.WindowEvent) {
    // This works on all platforms
})

// Platform-specific events for advanced use cases
window.OnWindowEvent(events.Mac.WindowWillClose, func(e *application.WindowEvent) {
    // macOS-specific "will close" event (before WindowClosing)
})

window.OnWindowEvent(events.Windows.WindowClosing, func(e *application.WindowEvent) {
    // Windows-specific close event
})
```

**Mise en correspondance des événements :**

Les événements natifs de la plateforme sont automatiquement mis en correspondance avec les événements communs :

- macOS : `events.Mac.WindowShouldClose` → `events.Common.WindowClosing`
- Windows : `events.Windows.WindowClosing` → `events.Common.WindowClosing`
- Linux : `events.Linux.WindowDeleteEvent` → `events.Common.WindowClosing`

Cette mise en correspondance s’effectue automatiquement en arrière-plan. Vous recevrez donc `events.Common.WindowClosing` quelle que soit la plateforme lorsque vous l’écoutez.

**Quand utiliser chaque type :**

- **Utilisez les événements communs** pour 99 % du code de votre application : ils offrent un comportement cohérent sur toutes les plateformes
- **Utilisez les événements natifs de la plateforme** uniquement lorsqu’une fonctionnalité propre à une plateforme n’est pas disponible dans les événements communs (par exemple, les événements du cycle de vie des fenêtres propres à macOS ou les événements de gestion de l’alimentation de Windows)

### Événements de l’application

| Événement | Description | Moment de l’émission | Annulable |
| --- | --- | --- | --- |
| `ApplicationOpenedWithFile` | Application ouverte avec un fichier | Lorsque l’application est lancée avec un fichier (par exemple, via une association de fichiers) | Non |
| `ApplicationStarted` | L’application a terminé son lancement | Une fois l’initialisation de l’application terminée et l’application prête | Non |
| `ApplicationLaunchedWithUrl` | Application lancée avec une URL | Lorsque l’application est lancée via un schéma d’URL | Non |
| `ThemeChanged` | Le thème système a changé | Lorsque le thème du système d’exploitation bascule entre les modes clair et sombre | Non |
| `SystemWillSleep` | Le système est sur le point d’être mis en veille | Juste avant la mise en veille du système d’exploitation (macOS / Windows / Linux avec logind) | Non |
| `SystemDidWake` | Le système est sorti de veille | Immédiatement à la sortie de veille (macOS / Windows / Linux avec logind) | Non |

**Utilisation :**

```go
import "github.com/wailsapp/wails/v3/pkg/events"

app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(e *application.ApplicationEvent) {
    app.Logger.Info("Application ready!")
})

app.Event.OnApplicationEvent(events.Common.ThemeChanged, func(e *application.ApplicationEvent) {
    // Update app theme
})
```

### Événements de fenêtre

| Événement | Description | Moment de l’émission | Annulable |
| --- | --- | --- | --- |
| `WindowClosing` | La fenêtre est sur le point de se fermer | Avant la fermeture de la fenêtre (clic de l’utilisateur sur X ou appel de Close()) | Oui |
| `WindowDidMove` | La fenêtre a été déplacée vers une nouvelle position | Après le changement de position de la fenêtre (avec temporisation anti-rebond) | Non |
| `WindowDidResize` | La fenêtre a été redimensionnée | Après le changement de taille de la fenêtre | Non |
| `WindowDPIChanged` | La mise à l’échelle liée au DPI de la fenêtre a changé | Lors du déplacement entre des moniteurs ayant des résolutions en PPP différentes (Windows) | Non |
| `WindowFilesDropped` | Fichiers déposés via la fonction native de glisser-déposer du système d’exploitation | Après le dépôt de fichiers depuis le système d’exploitation sur la fenêtre | Non |
| `WindowFocus` | La fenêtre a obtenu le focus | Lorsque la fenêtre devient active | Non |
| `WindowFullscreen` | La fenêtre est passée en plein écran | Après un appel de Fullscreen() ou le passage en plein écran par l’utilisateur | Non |
| `WindowHide` | La fenêtre a été masquée | Après un appel de Hide() ou lorsque la fenêtre devient occultée | Non |
| `WindowLostFocus` | La fenêtre a perdu le focus | Lorsque la fenêtre devient inactive | Non |
| `WindowMaximise` | La fenêtre a été agrandie | Après un appel de Maximise() ou l’agrandissement par l’utilisateur | Oui (macOS) |
| `WindowMinimise` | La fenêtre a été réduite | Après un appel de Minimise() ou la réduction par l’utilisateur | Oui (macOS) |
| `WindowRestore` | La fenêtre a été restaurée depuis son état réduit ou agrandi | Après un appel de Restore() (principalement sous Windows) | Non |
| `WindowRuntimeReady` | L’environnement d’exécution Wails est chargé et prêt | Lorsque l’initialisation de l’environnement d’exécution JavaScript est terminée | Non |
| `WindowShow` | La fenêtre est devenue visible | Après Show() ou lorsque la fenêtre devient visible | Non |
| `WindowUnFullscreen` | La fenêtre a quitté le mode plein écran | Après UnFullscreen() ou lorsque l’utilisateur quitte le mode plein écran | Non |
| `WindowUnMaximise` | La fenêtre a quitté l’état agrandi | Après UnMaximise() ou lorsque l’utilisateur restaure la fenêtre | Oui (macOS) |
| `WindowUnMinimise` | La fenêtre a quitté l’état réduit | Après UnMinimise()/Restore() ou lorsque l’utilisateur restaure la fenêtre | Oui (macOS) |
| `WindowZoomIn` | Le niveau de zoom du contenu de la fenêtre a augmenté | Après l’appel de ZoomIn() (principalement sur macOS) | Oui (macOS) |
| `WindowZoomOut` | Le niveau de zoom du contenu de la fenêtre a diminué | Après l’appel de ZoomOut() (principalement sur macOS) | Oui (macOS) |
| `WindowZoomReset` | Le zoom du contenu de la fenêtre a été réinitialisé à 100 % | Après l’appel de ZoomReset() (principalement sur macOS) | Oui (macOS) |
| `WindowDropZoneFilesDropped` | Fichiers déposés dans une zone de dépôt définie en JavaScript | Lorsque des fichiers sont déposés sur un élément doté d’une zone de dépôt | Non |

**Utilisation :**

```go
import "github.com/wailsapp/wails/v3/pkg/events"

// Listen for window events
window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
    app.Logger.Info("Window focused")
})

// Cancel window close
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    dlg := app.Dialog.Question().SetMessage("Close window?")
    yes := dlg.AddButton("Yes")
    no := dlg.AddButton("No")
    dlg.SetDefaultButton(yes)
    dlg.SetCancelButton(no)
    no.OnClick(func() { e.Cancel() })

    dlg.Show()
})

// Wait for runtime ready
window.OnWindowEvent(events.Common.WindowRuntimeReady, func(e *application.WindowEvent) {
    app.Logger.Info("Runtime ready, safe to emit events to frontend")
    window.EmitEvent("app-initialized", data)
})
```

**Remarques importantes :**

- **WindowRuntimeReady** est essentiel : attendez cet événement avant d’émettre des événements vers le frontend
- **WindowDidMove** et **WindowDidResize** font l’objet d’un anti-rebond (50 ms par défaut) afin d’éviter une saturation par les événements
- Il est possible d’empêcher les **événements annulables** en appelant `event.Cancel()` dans un gestionnaire `RegisterHook()`
- **WindowFilesDropped** concerne le dépôt de fichiers natif du système d’exploitation ; **WindowDropZoneFilesDropped** concerne les zones de dépôt web
- Certains événements sont propres à une plateforme (par exemple, WindowDPIChanged sous Windows et les événements de zoom principalement sous macOS)

## Conventions de nommage des événements

```go
// Good - descriptive and specific
app.Event.Emit("user:logged-in", user)
app.Event.Emit("data:fetch:complete", results)
app.Event.Emit("ui:theme:changed", theme)

// Bad - vague and unclear
app.Event.Emit("event1", data)
app.Event.Emit("update", stuff)
app.Event.Emit("e", value)
```

## Considérations relatives aux performances

### Anti-rebond des événements à haute fréquence

```go
type Service struct {
    app            *application.App
    lastEmit       time.Time
    debounceWindow time.Duration
}

func (s *Service) EmitWithDebounce(event string, data interface{}) {
    now := time.Now()
    if now.Sub(s.lastEmit) < s.debounceWindow {
        return // Skip this emission
    }

    s.app.Event.Emit(event, data)
    s.lastEmit = now
}
```

### Limitation de la fréquence des événements

```javascript
import { Events } from '@wailsio/runtime'

let lastUpdate = 0
const throttleMs = 100

Events.On('high-frequency-event', (data) => {
    const now = Date.now()
    if (now - lastUpdate < throttleMs) {
        return // Skip this update
    }

    processUpdate(data)
    lastUpdate = now
})
```
