---
title: "Fensterereignisse"
description: "Lebenszyklus- und Zustandsänderungsereignisse von Fenstern verarbeiten"
slug: "features/windows/events"
sourcePath: "features/windows/events.md"
---

## Fensterereignisse

Wails löst Fensterlebenszyklus- und Zustandsänderungsereignisse über eine einheitliche API aus: `OnWindowEvent` für passive Listener und `RegisterHook` für abbrechbare Hooks, die die Standardaktion verhindern können.

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

Beide Aufrufe geben eine `unsubscribe func()` zurück, die Sie aufrufen können, um den Handler zu entfernen.

Plattformübergreifende Ereignisse befinden sich in `events.Common.*`. Plattformspezifische Ereignisse befinden sich in `events.Mac.*`, `events.Windows.*` und `events.Linux.*`. Die vollständige Liste wird in `v3/pkg/events/events.go` generiert.

## Lebenszyklusereignisse

### Fenstererstellung

Führen Sie mit `app.Window.OnCreate` bei jeder Fenstererstellung einen Callback aus:

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

Der Callback erhält `application.Window` (ein Interface). Der `OnCreate`-Callback wird einmal pro Fenster aufgerufen, nachdem dessen Runtime initialisiert wurde.

### WindowClosing

Wird ausgelöst, wenn der Benutzer versucht, das Fenster zu schließen (durch Klicken auf X, ⌘W, Alt+F4 usw.).

Verwenden Sie einen **Hook** (`RegisterHook`) – nur Hooks können das Schließen abbrechen. Listener beobachten das Ereignis, können es aber nicht verhindern.

```go
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    if hasUnsavedChanges() {
        // Cancel the close — the window stays open.
        e.Cancel()
    }
})
```

**Wichtig:**

- `WindowClosing` wird bei vom Benutzer initiierten Schließversuchen ausgelöst.
- Ein `RegisterHook`-Callback kann `e.Cancel()` aufrufen, um das Fenster geöffnet zu halten.
- `OnWindowEvent`-Callbacks für `WindowClosing` sind Beobachter: Sie werden ausgelöst, können den Vorgang aber nicht abbrechen.

**Programmatisches Schließen:** Rufen Sie `window.Close()` auf. `window.Destroy()` existiert nicht.

### WindowRuntimeReady

Wird ausgelöst, sobald die Runtime im Fenster vollständig initialisiert ist – ab diesem Zeitpunkt können Sie sicher den JS-Kontext des Fensters aufrufen:

```go
window.OnWindowEvent(events.Common.WindowRuntimeReady, func(e *application.WindowEvent) {
    window.EmitEvent("app-ready", nil)
})
```

## Fokusereignisse

### WindowFocus

Wird aufgerufen, wenn das Fenster den Fokus erhält:

```go
window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
    fmt.Println("Window gained focus")
    updateTitleBar(true)
    app.Event.Emit("window-focused", window.ID())
})
```

### WindowLostFocus

Wird aufgerufen, wenn das Fenster den Fokus verliert:

```go
window.OnWindowEvent(events.Common.WindowLostFocus, func(e *application.WindowEvent) {
    fmt.Println("Window lost focus")
    updateTitleBar(false)
    saveCurrentState()
})
```

**Beispiel: fokusabhängige Benutzeroberfläche:**

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

## Zustandsänderungsereignisse

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

Aktivieren oder beenden Sie den Vollbildmodus programmgesteuert mit `window.Fullscreen()` / `window.UnFullscreen()` / `window.ToggleFullscreen()` – `SetFullscreen(bool)` existiert nicht. Fragen Sie den Zustand mit `window.IsFullscreen() bool` ab.

## Positions- und Größenereignisse

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

`WindowDidResize` und `WindowDidMove` enthalten im Ereignis selbst keine Koordinaten. Fragen Sie das Fenster innerhalb des Callbacks über `window.Size()` / `window.Position()` ab.

**Beispiel: responsives Layout:**

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

## Vollständiges Beispiel

Ein produktionsreifes Fenster mit vollständiger Ereignisverarbeitung:

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

## Ereigniskoordination

### Fensterübergreifende Ereignisse

Koordinieren Sie mehrere Fenster über den Ereignisbus der Anwendung:

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

### Ereignisketten

```go
window.OnWindowEvent(events.Common.WindowMaximise, func(e *application.WindowEvent) {
    saveWindowState()
    window.EmitEvent("layout-changed", "maximised")
    app.Event.Emit("window-maximised", window.ID())
})
```

### Entprellte Ereignisse

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

## Bewährte Methoden

### ✅ Empfohlen

- **Verwenden Sie Hooks zum Abbrechen** – nur `RegisterHook`-Callbacks können `e.Cancel()` aufrufen.
- **Speichern Sie beim Schließen den Zustand** – stellen Sie beim nächsten Start Position und Größe des Fensters wieder her.
- **Entprellen Sie häufige Ereignisse** – `WindowDidResize` und `WindowDidMove` werden in schneller Folge ausgelöst.
- **Verarbeiten Sie Fokusänderungen** – aktualisieren Sie die Benutzeroberfläche entsprechend.
- **Koordinieren Sie über Ereignisse** – verwenden Sie `app.Event.Emit` für die fensterübergreifende Kommunikation.
- **Melden Sie Handler ab**, wenn sie nicht mehr benötigt werden – `OnWindowEvent` und `RegisterHook` geben jeweils eine `func()` zum Abmelden zurück.

### ❌ Nicht tun

- **Blockieren Sie Ereignishandler nicht** – halten Sie deren Ausführung kurz.
- **Versuchen Sie nicht, den Vorgang über `OnWindowEvent`** abzubrechen – verwenden Sie `RegisterHook`.
- **Verwenden Sie `window.Destroy()`** nicht – es existiert nicht; verwenden Sie `window.Close()`.
- **Speichern Sie nicht bei jedem Ereignis** – entprellen Sie es zuerst.
- **Gehen Sie nicht davon aus, dass das Ereignis Koordinaten enthält** – rufen Sie `window.Position()` / `window.Size()` auf.

## Fehlerbehebung

### WindowClosing-Hook verhindert das Schließen nicht

**Ursache:** Sie haben `OnWindowEvent` anstelle von `RegisterHook` verwendet.

**Lösung:** Nur `RegisterHook`-Callbacks können `e.Cancel()` aufrufen. `OnWindowEvent`-Callbacks beobachten das Ereignis; sie können es nicht abbrechen.

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

### Ereignisse werden nicht ausgelöst

**Ursache:** Der Handler wurde erst registriert, nachdem das Ereignis eingetreten war.

**Lösung:** Registrieren Sie Handler unmittelbar nach der Fenstererstellung (oder im `app.Window.OnCreate`-Callback).

```go
window := app.Window.New()
window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) { /* ... */ })
```

### Speicherlecks

**Ursache:** Langlebige Handler werden von kurzlebigem Zustand referenziert.

**Lösung:** Bewahren Sie die von `OnWindowEvent`/`RegisterHook` zurückgegebene Funktion `func()` zum Abmelden auf und rufen Sie sie bei der Bereinigung auf.

```go
unsub := window.OnWindowEvent(events.Common.WindowDidResize, handler)
// ...later:
unsub()
```

## Nächste Schritte

**Fenstergrundlagen** – Lernen Sie die Grundlagen der Fensterverwaltung kennen [Mehr erfahren →](/features/windows/basics/)

**Mehrere Fenster** – Muster für Anwendungen mit mehreren Fenstern [Mehr erfahren →](/features/windows/multiple/)

**Ereignissystem** – Vertiefender Einblick in das Ereignissystem [Mehr erfahren →](/features/events/system/)

**Anwendungslebenszyklus** – Machen Sie sich mit dem Anwendungslebenszyklus vertraut [Mehr erfahren →](/concepts/lifecycle/)

---

**Fragen?** Fragen Sie auf [Discord](https://discord.gg/JDdSxwjhGf) nach oder sehen Sie sich die [Beispiele](https://github.com/wailsapp/wails/tree/master/v3/examples) an.
