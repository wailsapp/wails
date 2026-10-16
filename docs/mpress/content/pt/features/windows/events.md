---
title: "Eventos de janela"
description: "Trate eventos do ciclo de vida e de mudança de estado da janela"
slug: "features/windows/events"
sourcePath: "features/windows/events.md"
---

## Eventos de janela

O Wails despacha eventos do ciclo de vida e de mudança de estado da janela por meio de uma única API: `OnWindowEvent` para listeners passivos e `RegisterHook` para hooks canceláveis que podem impedir a ação padrão.

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

Ambas as chamadas retornam uma função `unsubscribe func()` que você pode chamar para remover o manipulador.

Os eventos multiplataforma ficam em `events.Common.*`. Os eventos específicos de cada plataforma ficam em `events.Mac.*`, `events.Windows.*` e `events.Linux.*`. A lista completa é gerada em `v3/pkg/events/events.go`.

## Eventos do ciclo de vida

### Criação da janela

Execute um callback sempre que uma janela for criada usando `app.Window.OnCreate`:

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

O callback recebe `application.Window` (uma interface). O callback `OnCreate` é invocado uma vez por janela, após a inicialização do runtime da janela.

### WindowClosing

Disparado quando o usuário tenta fechar a janela (clicando no X, pressionando ⌘W, Alt+F4 etc.).

Use um **hook** (`RegisterHook`) — somente hooks podem cancelar o fechamento. Listeners observam o evento, mas não podem impedi-lo.

```go
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    if hasUnsavedChanges() {
        // Cancel the close — the window stays open.
        e.Cancel()
    }
})
```

**Importante:**

- `WindowClosing` é disparado nas tentativas de fechamento iniciadas pelo usuário.
- Um callback de `RegisterHook` pode chamar `e.Cancel()` para manter a janela aberta.
- Os callbacks de `OnWindowEvent` para `WindowClosing` são observadores — eles são executados, mas não podem cancelar.

**Fechamento programático:** chame `window.Close()`. Não existe `window.Destroy()`.

### WindowRuntimeReady

Disparado quando o runtime da janela termina de ser inicializado — este é um ponto seguro para fazer chamadas ao contexto JS da janela:

```go
window.OnWindowEvent(events.Common.WindowRuntimeReady, func(e *application.WindowEvent) {
    window.EmitEvent("app-ready", nil)
})
```

## Eventos de foco

### WindowFocus

Chamado quando a janela recebe foco:

```go
window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
    fmt.Println("Window gained focus")
    updateTitleBar(true)
    app.Event.Emit("window-focused", window.ID())
})
```

### WindowLostFocus

Chamado quando a janela perde o foco:

```go
window.OnWindowEvent(events.Common.WindowLostFocus, func(e *application.WindowEvent) {
    fmt.Println("Window lost focus")
    updateTitleBar(false)
    saveCurrentState()
})
```

**Exemplo: interface sensível ao foco:**

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

## Eventos de mudança de estado

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

Entre e saia do modo de tela cheia programaticamente com `window.Fullscreen()` / `window.UnFullscreen()` / `window.ToggleFullscreen()` — não existe `SetFullscreen(bool)`. Consulte o estado com `window.IsFullscreen() bool`.

## Eventos de posição e tamanho

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

`WindowDidResize` e `WindowDidMove` não incluem coordenadas no próprio evento — consulte a janela por meio de `window.Size()` / `window.Position()` dentro do callback.

**Exemplo: layout responsivo:**

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

## Exemplo completo

Uma janela pronta para produção, com tratamento completo de eventos:

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

## Coordenação de eventos

### Eventos entre janelas

Coordene várias janelas usando o barramento de eventos da aplicação:

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

### Cadeias de eventos

```go
window.OnWindowEvent(events.Common.WindowMaximise, func(e *application.WindowEvent) {
    saveWindowState()
    window.EmitEvent("layout-changed", "maximised")
    app.Event.Emit("window-maximised", window.ID())
})
```

### Eventos com debounce

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

## Práticas recomendadas

### ✅ Faça

- **Use hooks para cancelamento** — somente callbacks de `RegisterHook` podem chamar `e.Cancel()`.
- **Salve o estado ao fechar** — restaure a posição e o tamanho da janela na próxima inicialização.
- **Aplique debounce a eventos frequentes** — `WindowDidResize` e `WindowDidMove` são disparados rapidamente.
- **Trate as mudanças de foco** — atualize a interface adequadamente.
- **Coordene por meio de eventos** — use `app.Event.Emit` para trocar mensagens entre janelas.
- **Cancele a inscrição** quando os manipuladores não forem mais necessários — tanto `OnWindowEvent` quanto `RegisterHook` retornam uma função `func()` para cancelar a inscrição.

### ❌ Não faça

- **Não bloqueie os manipuladores de eventos** — mantenha-os rápidos.
- **Não tente cancelar por meio de `OnWindowEvent`** — use `RegisterHook`.
- **Não use `window.Destroy()`** — ele não existe; use `window.Close()`.
- **Não salve a cada evento** — aplique debounce primeiro.
- **Não presuma que o evento contém coordenadas** — chame `window.Position()` / `window.Size()`.

## Solução de problemas

### O hook WindowClosing não impede o fechamento

**Causa:** você usou `OnWindowEvent` em vez de `RegisterHook`.

**Solução:** somente callbacks de `RegisterHook` podem chamar `e.Cancel()`. Os callbacks de `OnWindowEvent` apenas observam; eles não podem cancelar.

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

### Eventos não disparados

**Causa:** o manipulador foi registrado depois que o evento ocorreu.

**Solução:** registre os manipuladores imediatamente após a criação da janela (ou no callback `app.Window.OnCreate`).

```go
window := app.Window.New()
window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) { /* ... */ })
```

### Vazamentos de memória

**Causa:** manipuladores de longa duração mantidos por um estado de curta duração.

**Solução:** mantenha a função de cancelamento de inscrição `func()` retornada por `OnWindowEvent`/`RegisterHook` e chame-a durante a limpeza.

```go
unsub := window.OnWindowEvent(events.Common.WindowDidResize, handler)
// ...later:
unsub()
```

## Próximos passos

**Noções básicas sobre janelas** - Aprenda os fundamentos do gerenciamento de janelas [Saiba mais →](/features/windows/basics/)

**Várias janelas** - Padrões para aplicações com várias janelas [Saiba mais →](/features/windows/multiple/)

**Sistema de eventos** - Conheça o sistema de eventos em detalhes [Saiba mais →](/features/events/system/)

**Ciclo de vida da aplicação** - Entenda o ciclo de vida da aplicação [Saiba mais →](/concepts/lifecycle/)

---

**Dúvidas?** Pergunte no [Discord](https://discord.gg/JDdSxwjhGf) ou consulte os [exemplos](https://github.com/wailsapp/wails/tree/master/v3/examples).
