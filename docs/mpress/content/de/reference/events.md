---
title: "Events-API"
description: "Vollständige Referenz für die Events-API"
slug: "reference/events"
sourcePath: "reference/events.md"
---

## Übersicht

Die Events-API stellt Methoden zum Auslösen und Abonnieren von Ereignissen bereit und ermöglicht so die Kommunikation zwischen verschiedenen Teilen Ihrer Anwendung.

**Ereignistypen:**

- **Anwendungsereignisse** – Ereignisse im Anwendungslebenszyklus (Start, Beenden)
- **Fensterereignisse** – Änderungen des Fensterzustands (Fokus erhalten, Fokus verloren, Größenänderung)
- **Benutzerdefinierte Ereignisse** – benutzerdefinierte Ereignisse für die anwendungsspezifische Kommunikation

**Kommunikationsmuster:**

- **Go zum Frontend** – Ereignisse in Go auslösen und in JavaScript abonnieren
- **Frontend zu Go** – nicht direkt (verwenden Sie stattdessen Service-Bindings)
- **Frontend zu Frontend** – über Go oder lokale Runtime-Ereignisse
- **Fenster zu Fenster** – bestimmte Fenster gezielt ansprechen oder an alle senden

## Ereignismethoden (Go)

### app.Event.Emit()

Löst ein benutzerdefiniertes Ereignis für alle Fenster aus. Gibt `true` zurück, wenn ein Hook das Auslösen abgebrochen hat.

```go
func (em *EventManager) Emit(name string, data ...any) bool
```

**Parameter:**

- `name` – Ereignisname
- `data` – optionale Daten, die mit dem Ereignis gesendet werden

**Beispiel:**

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

Abonniert benutzerdefinierte Ereignisse in Go.

```go
func (em *EventManager) On(name string, callback func(*CustomEvent)) func()
```

**Parameter:**

- `name` – Name des zu abonnierenden Ereignisses
- `callback` – Funktion, die beim Auslösen des Ereignisses aufgerufen wird

**Rückgabewert:** Bereinigungsfunktion zum Entfernen des Ereignis-Listeners

**Beispiel:**

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

### Fensterspezifische Ereignisse

So lösen Sie Ereignisse für ein bestimmtes Fenster aus:

```go
// Emit to specific window
window.EmitEvent("notification", "Hello from Go!")

// Emit to all windows
app.Event.Emit("global-update", data)
```

## Ereignismethoden (Frontend)

### On()

Abonniert Ereignisse aus Go.

```javascript
import { Events } from '@wailsio/runtime'

Events.On(eventName, callback)
```

**Parameter:**

- `eventName` – Name des zu abonnierenden Ereignisses
- `callback` – Funktion, die beim Empfang des Ereignisses aufgerufen wird

**Rückgabewert:** Bereinigungsfunktion

**Beispiel:**

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

Abonniert das einmalige Auftreten eines Ereignisses.

```javascript
import { Events } from '@wailsio/runtime'

Events.Once(eventName, callback)
```

**Beispiel:**

```javascript
import { Events } from '@wailsio/runtime'

// Listen for first occurrence only
Events.Once('initialization-complete', (data) => {
    console.log('App initialized!', data)
    // This will only fire once
})
```

### Off()

Entfernt alle Listener für mindestens ein Ereignis. `Off` akzeptiert eine variable Anzahl von Ereignisnamen als Zeichenfolgen und akzeptiert **keinen** Callback. Um einen einzelnen Listener zu entfernen, speichern Sie die von `Events.On(...)` zurückgegebene Funktion zum Abbestellen und rufen Sie sie auf.

```typescript
import { Events } from '@wailsio/runtime'

Events.Off(...eventNames: string[]): void
```

**Beispiel:**

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

Entfernt **alle** Ereignis-Listener. Akzeptiert keine Argumente.

```typescript
import { Events } from '@wailsio/runtime'

Events.OffAll(): void
```

**Beispiel:**

```javascript
import { Events } from '@wailsio/runtime'

// Remove all listeners — typically used during teardown.
Events.OffAll()
```

### OnMultiple()

Abonniert ein Ereignis für bis zu `max` Auslösungen und meldet es danach automatisch ab.

```typescript
Events.OnMultiple(eventName: string, callback, max: number): () => void
```

**Beispiel:**

```javascript
Events.OnMultiple('progress', (data) => {
    console.log('progress', data)
}, 5)
```

## Anwendungsereignisse

### app.Event.OnApplicationEvent()

Abonniert Ereignisse im Anwendungslebenszyklus.

```go
func (em *EventManager) OnApplicationEvent(
    eventType events.ApplicationEventType,
    callback func(*ApplicationEvent),
) func()
```

Die Ereigniskonstanten befinden sich im Paket `events`: `events.Common.*` für plattformübergreifende Ereignisse und `events.Mac.*` / `events.Windows.*` / `events.Linux.*` für plattformspezifische Ereignisse.

**Häufig verwendete Anwendungsereignisse:**

- `events.Common.ApplicationStarted` – Anwendung wurde vollständig gestartet.
- `events.Common.ThemeChanged` – Das Systemdesign wechselte zwischen hell und dunkel.
- `events.Common.ApplicationOpenedWithFile` – Über eine Dateizuordnung gestartet.
- `events.Common.ApplicationLaunchedWithUrl` – Über ein URL-Schema gestartet.

Es gibt **keine** generische Ereigniskonstante für das Beenden der Anwendung. Registrieren Sie die Bereinigung beim Beenden über `application.Options.OnShutdown` oder `app.OnShutdown(func())`.

**Beispiel:**

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

## Fensterereignisse

### OnWindowEvent()

Abonniert fensterspezifische Ereignisse.

```go
func (w *WebviewWindow) OnWindowEvent(
    eventType events.WindowEventType,
    callback func(*WindowEvent),
) func()
```

Die Konstanten für Fensterereignisse befinden sich im Paket `events`: `events.Common.*` für plattformübergreifende Ereignisse (und `events.Mac.*` / `events.Windows.*` / `events.Linux.*` für plattformspezifische Ereignisse).

**Häufig verwendete Fensterereignisse:**

- `events.Common.WindowFocus` – Fenster hat den Fokus erhalten.
- `events.Common.WindowLostFocus` – Fenster hat den Fokus verloren.
- `events.Common.WindowClosing` – Fenster wird gleich geschlossen (über `RegisterHook` abbrechbar).
- `events.Common.WindowDidResize` – Fenstergröße wurde geändert.
- `events.Common.WindowDidMove` – Fenster wurde verschoben.
- `events.Common.WindowMinimise` / `WindowUnMinimise` / `WindowMaximise` / `WindowUnMaximise` / `WindowFullscreen` / `WindowUnFullscreen`.
- `events.Common.WindowRuntimeReady` – Ereignisse können sicher an die Runtime im Fenster gesendet werden.

**Beispiel:**

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

## Gängige Muster

Diese Muster zeigen bewährte Ansätze für den Einsatz von Ereignissen in realen Anwendungen. Jedes Muster löst eine bestimmte Kommunikationsaufgabe zwischen Go-Backend und Frontend und unterstützt Sie beim Erstellen reaktionsschneller, gut strukturierter Anwendungen.

### Anfrage-Antwort-Muster

Verwenden Sie dieses Muster, um das Frontend über den Abschluss von Backend-Vorgängen zu informieren, etwa nach dem Abrufen von Daten, dem Verarbeiten von Dateien oder dem Ausführen von Hintergrundaufgaben. Die Service-Bindung gibt Daten direkt zurück, während Ereignisse zusätzliche Benachrichtigungen für UI-Aktualisierungen bereitstellen, beispielsweise zum Anzeigen von Toast-Meldungen oder Aktualisieren von Listen.

**Go:**

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

**JavaScript:**

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

### Fortschrittsaktualisierungen

Ideal für lang laufende Vorgänge wie Datei-Uploads, Stapelverarbeitung, den Import großer Datenmengen oder Videokodierung. Senden Sie während des Vorgangs Fortschrittsereignisse, um Fortschrittsbalken, Statustexte oder Schrittanzeigen in der UI zu aktualisieren und Benutzern Echtzeitfeedback zu geben.

**Go:**

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

**JavaScript:**

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

### Kommunikation zwischen mehreren Fenstern

Ideal für Anwendungen mit mehreren Fenstern, etwa Einstellungsdialogen, Dashboards oder Dokumentbetrachtern. Senden Sie Ereignisse an alle Fenster, um deren Zustand zu synchronisieren, beispielsweise bei Änderungen des Designs oder der Benutzereinstellungen. Für fensterspezifische Aktualisierungen können Sie Ereignisse gezielt an einzelne Fenster senden.

**Go:**

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

**JavaScript:**

```javascript
import { Events } from '@wailsio/runtime'

// Listen in any window
Events.On('theme-changed', (theme) => {
    document.body.className = theme
})
```

### Zustandssynchronisierung

Verwenden Sie dieses Muster, wenn Frontend- und Backend-Zustand synchron bleiben müssen, etwa bei Benutzersitzungen, der Anwendungskonfiguration oder Funktionen für die Zusammenarbeit. Wenn sich der Zustand im Backend ändert, senden Sie Ereignisse zur Aktualisierung aller verbundenen Frontends, um die Konsistenz in der gesamten Anwendung sicherzustellen.

**Go:**

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

**JavaScript:**

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

### Ereignisgesteuerte Benachrichtigungen

Ideal zum Anzeigen von Benutzerfeedback wie Erfolgsbestätigungen, Fehlermeldungen oder Hinweisen. Statt UI-Code direkt aus Services aufzurufen, senden Sie Benachrichtigungsereignisse, die das Frontend einheitlich verarbeitet. Dadurch lassen sich Benachrichtigungsstile leicht ändern oder Funktionen wie ein Benachrichtigungsverlauf hinzufügen.

**Go:**

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

**JavaScript:**

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

## Vollständiges Beispiel

**Go:**

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

**JavaScript:**

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

## Integrierte Ereignisse

Wails stellt integrierte Systemereignisse für den Lebenszyklus von Anwendung und Fenstern bereit. Das Framework sendet diese Ereignisse automatisch.

### Gemeinsame und plattformspezifische Ereignisse

Wails stellt zwei Arten von Systemereignissen bereit:

**Gemeinsame Ereignisse** (`events.Common.*`) sind plattformübergreifende Abstraktionen, die unter macOS, Windows und Linux einheitlich funktionieren. Verwenden Sie diese Ereignisse in Ihrer Anwendung, um maximale Portabilität zu erreichen.

**Plattformspezifische Ereignisse** (`events.Mac.*`, `events.Windows.*`, `events.Linux.*`) sind die zugrunde liegenden betriebssystemspezifischen Ereignisse, aus denen gemeinsame Ereignisse abgeleitet werden. Sie ermöglichen den Zugriff auf plattformspezifisches Verhalten und Sonderfälle.

**Funktionsweise:**

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

**Ereigniszuordnung:**

Plattformspezifische Ereignisse werden automatisch gemeinsamen Ereignissen zugeordnet:

- macOS: `events.Mac.WindowShouldClose` → `events.Common.WindowClosing`
- Windows: `events.Windows.WindowClosing` → `events.Common.WindowClosing`
- Linux: `events.Linux.WindowDeleteEvent` → `events.Common.WindowClosing`

Diese Zuordnung erfolgt automatisch im Hintergrund. Wenn Sie auf `events.Common.WindowClosing` lauschen, empfangen Sie das Ereignis daher unabhängig von der Plattform.

**Wann welcher Typ verwendet werden sollte:**

- **Verwenden Sie gemeinsame Ereignisse** für 99 % Ihres Anwendungscodes – sie bieten auf allen Plattformen einheitliches Verhalten.
- **Verwenden Sie plattformspezifische Ereignisse** nur, wenn Sie plattformspezifische Funktionen benötigen, die über gemeinsame Ereignisse nicht verfügbar sind, beispielsweise macOS-spezifische Ereignisse des Fensterlebenszyklus oder Windows-Ereignisse zur Energieverwaltung.

### Anwendungsereignisse

| Ereignis | Beschreibung | Auslösezeitpunkt | Abbrechbar |
| --- | --- | --- | --- |
| `ApplicationOpenedWithFile` | Anwendung mit einer Datei geöffnet | Wenn die App mit einer Datei gestartet wird (z. B. über eine Dateizuordnung) | Nein |
| `ApplicationStarted` | Anwendung wurde vollständig gestartet | Nachdem die Initialisierung der App abgeschlossen und die App bereit ist | Nein |
| `ApplicationLaunchedWithUrl` | Anwendung wurde mit einer URL gestartet | Wenn die App über ein URL-Schema gestartet wird | Nein |
| `ThemeChanged` | Systemdesign wurde geändert | Wenn das Betriebssystem zwischen hellem und dunklem Modus wechselt | Nein |
| `SystemWillSleep` | System wird gleich in den Energiesparmodus versetzt | Unmittelbar bevor das Betriebssystem in den Energiesparmodus versetzt wird (macOS / Windows / Linux mit logind) | Nein |
| `SystemDidWake` | System wurde aus dem Energiesparmodus reaktiviert | Unmittelbar beim Reaktivieren aus dem Energiesparmodus (macOS / Windows / Linux mit logind) | Nein |

**Verwendung:**

```go
import "github.com/wailsapp/wails/v3/pkg/events"

app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(e *application.ApplicationEvent) {
    app.Logger.Info("Application ready!")
})

app.Event.OnApplicationEvent(events.Common.ThemeChanged, func(e *application.ApplicationEvent) {
    // Update app theme
})
```

### Fensterereignisse

| Ereignis | Beschreibung | Auslösezeitpunkt | Abbrechbar |
| --- | --- | --- | --- |
| `WindowClosing` | Fenster wird gleich geschlossen | Bevor das Fenster geschlossen wird (Benutzer hat auf X geklickt oder Close() wurde aufgerufen) | Ja |
| `WindowDidMove` | Fenster wurde an eine neue Position verschoben | Nachdem sich die Fensterposition geändert hat (entprellt) | Nein |
| `WindowDidResize` | Fenstergröße wurde geändert | Nachdem sich die Fenstergröße geändert hat | Nein |
| `WindowDPIChanged` | DPI-Skalierung des Fensters wurde geändert | Beim Verschieben zwischen Monitoren mit unterschiedlichen DPI-Werten (Windows) | Nein |
| `WindowFilesDropped` | Dateien wurden per nativem Drag-and-drop des Betriebssystems abgelegt | Nachdem Dateien aus dem Betriebssystem auf dem Fenster abgelegt wurden | Nein |
| `WindowFocus` | Fenster hat den Fokus erhalten | Wenn das Fenster aktiv wird | Nein |
| `WindowFullscreen` | Fenster ist in den Vollbildmodus gewechselt | Nach Fullscreen() oder wenn der Benutzer in den Vollbildmodus wechselt | Nein |
| `WindowHide` | Fenster wurde ausgeblendet | Nach Hide() oder wenn das Fenster verdeckt wird | Nein |
| `WindowLostFocus` | Fenster hat den Fokus verloren | Wenn das Fenster inaktiv wird | Nein |
| `WindowMaximise` | Fenster wurde maximiert | Nach Maximise() oder wenn der Benutzer das Fenster maximiert | Ja (macOS) |
| `WindowMinimise` | Fenster wurde minimiert | Nach Minimise() oder wenn der Benutzer das Fenster minimiert | Ja (macOS) |
| `WindowRestore` | Fenster wurde aus dem minimierten oder maximierten Zustand wiederhergestellt | Nach Restore() (hauptsächlich unter Windows) | Nein |
| `WindowRuntimeReady` | Wails-Laufzeit wurde geladen und ist bereit | Wenn die Initialisierung der JavaScript-Laufzeit abgeschlossen ist | Nein |
| `WindowShow` | Fenster wurde sichtbar | Nach Show() oder wenn das Fenster sichtbar wird | Nein |
| `WindowUnFullscreen` | Fenster hat den Vollbildmodus verlassen | Nach UnFullscreen() oder wenn der Benutzer den Vollbildmodus verlässt | Nein |
| `WindowUnMaximise` | Fenster hat den maximierten Zustand verlassen | Nach UnMaximise() oder wenn der Benutzer die Maximierung aufhebt | Ja (macOS) |
| `WindowUnMinimise` | Fenster hat den minimierten Zustand verlassen | Nach UnMinimise()/Restore() oder wenn der Benutzer das Fenster wiederherstellt | Ja (macOS) |
| `WindowZoomIn` | Zoom der Fensterinhalte wurde vergrößert | Nach dem Aufruf von ZoomIn() (hauptsächlich unter macOS) | Ja (macOS) |
| `WindowZoomOut` | Zoom der Fensterinhalte wurde verkleinert | Nach dem Aufruf von ZoomOut() (hauptsächlich unter macOS) | Ja (macOS) |
| `WindowZoomReset` | Zoom der Fensterinhalte wurde auf 100 % zurückgesetzt | Nach dem Aufruf von ZoomReset() (hauptsächlich unter macOS) | Ja (macOS) |
| `WindowDropZoneFilesDropped` | Dateien wurden in einer in JS definierten Ablagezone abgelegt | Wenn Dateien auf einem Element mit Ablagezone abgelegt werden | Nein |

**Verwendung:**

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

**Wichtige Hinweise:**

- **WindowRuntimeReady** ist unverzichtbar – warten Sie auf dieses Ereignis, bevor Sie Ereignisse an das Frontend senden
- **WindowDidMove** und **WindowDidResize** werden entprellt (standardmäßig 50 ms), um eine Ereignisflut zu verhindern
- **Abbrechbare Ereignisse** können durch den Aufruf von `event.Cancel()` in einem `RegisterHook()`-Handler verhindert werden
- **WindowFilesDropped** ist für native Dateiablagen des Betriebssystems vorgesehen; **WindowDropZoneFilesDropped** ist für webbasierte Ablagezonen vorgesehen
- Einige Ereignisse sind plattformspezifisch (z. B. WindowDPIChanged unter Windows, Zoomereignisse hauptsächlich unter macOS)

## Konventionen für Ereignisnamen

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

## Leistungsaspekte

### Entprellen hochfrequenter Ereignisse

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

### Drosseln von Ereignissen

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
