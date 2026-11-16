---
title: "Sistema de eventos"
description: "Comunique-se entre componentes com o sistema de eventos"
slug: "features/events/system"
sourcePath: "features/events/system.md"
---

## Sistema de eventos

O Wails oferece um **sistema de eventos unificado** para comunicação pub/sub. Emita eventos de qualquer lugar e escute-os em qualquer lugar — do Go para o JavaScript, do JavaScript para o Go e de uma janela para outra —, permitindo uma arquitetura desacoplada com eventos tipados e hooks de ciclo de vida.

## Início rápido

**Go (emitir):**

```go
app.Event.Emit("user-logged-in", map[string]interface{}{
    "userId": 123,
    "name": "Alice",
})
```

**JavaScript (escutar):**

```javascript
import { Events } from '@wailsio/runtime'

Events.On("user-logged-in", (event) => {
    console.log(`User ${event.data.name} logged in`)
})
```

**É só isso!** Pub/sub entre linguagens.

## Tipos de eventos

### Eventos personalizados

Eventos específicos da sua aplicação:

```go
// Emit from Go
app.Event.Emit("order-created", order)
app.Event.Emit("payment-processed", payment)
app.Event.Emit("notification", message)
```

```javascript
// Listen in JavaScript
Events.On("order-created", handleOrder)
Events.On("payment-processed", handlePayment)
Events.On("notification", showNotification)
```

### Eventos do sistema

Eventos integrados do sistema operacional e da aplicação:

```go
import "github.com/wailsapp/wails/v3/pkg/events"

// Theme changes
app.Event.OnApplicationEvent(events.Common.ThemeChanged, func(e *application.ApplicationEvent) {
    if e.Context().IsDarkMode() {
        app.Logger.Info("Dark mode enabled")
    }
})

// Application lifecycle
app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(e *application.ApplicationEvent) {
    app.Logger.Info("Application started")
})
```

### Eventos de janela

Eventos específicos de janela:

```go
window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
    app.Logger.Info("Window focused")
})

window.OnWindowEvent(events.Common.WindowClosing, func(e *application.WindowEvent) {
    app.Logger.Info("Window closing")
})
```

## Emissão de eventos

### No Go

**Emissão básica:**

```go
app.Event.Emit("event-name", data)
```

**Com diferentes tipos de dados:**

```go
// String
app.Event.Emit("message", "Hello")

// Number
app.Event.Emit("count", 42)

// Struct
app.Event.Emit("user", User{ID: 1, Name: "Alice"})

// Map
app.Event.Emit("config", map[string]interface{}{
    "theme": "dark",
    "fontSize": 14,
})

// Array
app.Event.Emit("items", []string{"a", "b", "c"})
```

**Para uma janela específica:**

```go
window.EmitEvent("window-specific-event", data)
```

### No JavaScript

```javascript
import { Events } from '@wailsio/runtime'

// Emit to Go
Events.Emit("button-clicked", { buttonId: "submit" })

// Emit to all windows
Events.Emit("broadcast-message", "Hello everyone")
```

## Escuta de eventos

### No Go

**Eventos da aplicação:**

```go
app.Event.On("custom-event", func(e *application.CustomEvent) {
    data := e.Data
    // Handle event
})
```

**Com asserção de tipo:**

```go
app.Event.On("user-updated", func(e *application.CustomEvent) {
    user := e.Data.(User)
    app.Logger.Info("User updated", "name", user.Name)
})
```

**Vários manipuladores:**

```go
// All handlers will be called
app.Event.On("order-created", logOrder)
app.Event.On("order-created", sendEmail)
app.Event.On("order-created", updateInventory)
```

**Ordem de entrega:** cada manipulador registrado com `On`, `Once` ou `OnMultiple` processa um evento por vez, portanto nunca é executado simultaneamente consigo mesmo e recebe chamadas sequenciais de `Emit` na ordem: se uma chamada de `Emit` retornar antes de a próxima ser feita, todos os manipuladores as receberão nessa ordem. Chamadas simultâneas de `Emit` provenientes de goroutines diferentes não têm uma ordem definida entre si. Manipuladores diferentes são executados de forma independente: um manipulador lento não atrasa os demais, e não há ordenação entre manipuladores. `Emit` nunca aguarda os manipuladores. Depois que um manipulador é removido, os eventos que ainda estiverem enfileirados para ele serão descartados.

### No JavaScript

**Listener básico:**

```javascript
import { Events } from '@wailsio/runtime'

Events.On("event-name", (event) => {
    console.log("Event received:", event.data)
})
```

**Com limpeza:**

```javascript
const unsubscribe = Events.On("event-name", handleEvent)

// Later, stop listening
unsubscribe()
```

**Vários manipuladores:**

```javascript
Events.On("data-updated", updateUI)
Events.On("data-updated", saveToCache)
Events.On("data-updated", logChange)
```

**Manipuladores de execução única:**

```javascript
Events.Once("data-updated", updateVariable)
```

## Eventos do sistema

### Eventos da aplicação

**Eventos comuns (multiplataforma):**

```go
import "github.com/wailsapp/wails/v3/pkg/events"

// Application started
app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(e *application.ApplicationEvent) {
    app.Logger.Info("App started")
})

// Theme changed
app.Event.OnApplicationEvent(events.Common.ThemeChanged, func(e *application.ApplicationEvent) {
    isDark := e.Context().IsDarkMode()
    app.Event.Emit("theme-changed", isDark)
})

// System about to suspend
app.Event.OnApplicationEvent(events.Common.SystemWillSleep, func(e *application.ApplicationEvent) {
    flushPendingWrites()
})

// System resumed from suspend
app.Event.OnApplicationEvent(events.Common.SystemDidWake, func(e *application.ApplicationEvent) {
    reconnectSockets()
    refreshState()
})

// File opened
app.Event.OnApplicationEvent(events.Common.ApplicationOpenedWithFile, func(e *application.ApplicationEvent) {
    // Single-file: e.Context().Filename() returns ""; for multi-file launches
    // (Finder "Open With..." → multiple selection) use OpenedFiles().
    filePath := e.Context().Filename()
    openFile(filePath)
})
```

**Eventos específicos da plataforma:**

@tabs{sync-key="platform"}
[macOS]
```go
// Application became active
app.Event.OnApplicationEvent(events.Mac.ApplicationDidBecomeActive, func(e *application.ApplicationEvent) {
    app.Logger.Info("App became active")
})

// Application will terminate
app.Event.OnApplicationEvent(events.Mac.ApplicationWillTerminate, func(e *application.ApplicationEvent) {
    cleanup()
})

// System about to sleep / resumed
app.Event.OnApplicationEvent(events.Mac.ApplicationWillSleep, func(e *application.ApplicationEvent) {
    flushPendingWrites()
})
app.Event.OnApplicationEvent(events.Mac.ApplicationDidWake, func(e *application.ApplicationEvent) {
    refreshState()
})

// Displays sleeping/waking — distinct from system sleep (e.g. lid lowered)
app.Event.OnApplicationEvent(events.Mac.ApplicationScreensDidSleep, func(e *application.ApplicationEvent) {
    pauseRendering()
})
app.Event.OnApplicationEvent(events.Mac.ApplicationScreensDidWake, func(e *application.ApplicationEvent) {
    resumeRendering()
})
```

[Windows]
```go
// Power status changed
app.Event.OnApplicationEvent(events.Windows.APMPowerStatusChange, func(e *application.ApplicationEvent) {
    app.Logger.Info("Power status changed")
})

// System suspending
app.Event.OnApplicationEvent(events.Windows.APMSuspend, func(e *application.ApplicationEvent) {
    saveState()
})

// System resumed. APMResumeAutomatic always fires on resume;
// APMResumeSuspend fires after it when the wake was triggered by
// user input. Prefer Common.SystemDidWake if you don't need to
// distinguish.
app.Event.OnApplicationEvent(events.Windows.APMResumeAutomatic, func(e *application.ApplicationEvent) {
    refreshState()
})
```

[Linux]
```go
// Application startup
app.Event.OnApplicationEvent(events.Linux.ApplicationStartup, func(e *application.ApplicationEvent) {
    app.Logger.Info("App starting")
})

// Theme changed
app.Event.OnApplicationEvent(events.Linux.SystemThemeChanged, func(e *application.ApplicationEvent) {
    updateTheme()
})

// System sleep/wake — sourced from systemd-logind's PrepareForSleep
// signal on the system bus. Requires logind/elogind exposing
// org.freedesktop.login1, so it may be unavailable on distros/setups
// without it (Alpine, Void, some Devuan setups).
app.Event.OnApplicationEvent(events.Linux.SystemWillSleep, func(e *application.ApplicationEvent) {
    flushPendingWrites()
})
app.Event.OnApplicationEvent(events.Linux.SystemDidWake, func(e *application.ApplicationEvent) {
    refreshState()
})
```

@end

### Eventos de janela

**Eventos comuns de janela:**

```go
// Window focus
window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
    app.Logger.Info("Window focused")
})

// Window blur
window.OnWindowEvent(events.Common.WindowLostFocus, func(e *application.WindowEvent) {
    app.Logger.Info("Window blurred")
})

// Window closing
window.OnWindowEvent(events.Common.WindowClosing, func(e *application.WindowEvent) {
    if hasUnsavedChanges() {
        e.Cancel()  // Prevent close
    }
})

// Window closed
window.OnWindowEvent(events.Common.WindowClosing, func(e *application.WindowEvent) {
    cleanup()
})
```

## Hooks de eventos

Os hooks são executados **antes** dos listeners padrão e podem **cancelar** eventos:

```go
// Hook - runs first, can cancel
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    if hasUnsavedChanges() {
        result := showConfirmdialog("Unsaved changes. Close anyway?")
        if result != "yes" {
            e.Cancel()  // Prevent window close
        }
    }
})

// Standard listener - runs after hooks
window.OnWindowEvent(events.Common.WindowClosing, func(e *application.WindowEvent) {
    app.Logger.Info("Window closing")
})
```

**Principais diferenças:**

| Recurso | Hooks | Listeners padrão |
| --- | --- | --- |
| Ordem de execução | Primeiro, na ordem de registro | Depois dos hooks, sem ordem garantida |
| Bloqueio | Síncrono, bloqueia o próximo hook | Assíncrono, não bloqueante |
| Pode cancelar | Sim | Não (já foi propagado) |
| Caso de uso | Controle de fluxo, validação | Registro de logs, efeitos colaterais |

## Padrões de eventos

### Padrão pub/sub

```go
// Publisher (service)
type OrderService struct {
    app *application.App
}

func (o *OrderService) CreateOrder(items []Item) (*Order, error) {
    order := &Order{Items: items}
    
    if err := o.saveOrder(order); err != nil {
        return nil, err
    }
    
    // Publish event
    o.app.Event.Emit("order-created", order)
    
    return order, nil
}

// Subscribers
app.Event.On("order-created", func(e *application.CustomEvent) {
    order := e.Data.(*Order)
    sendConfirmationEmail(order)
})

app.Event.On("order-created", func(e *application.CustomEvent) {
    order := e.Data.(*Order)
    updateInventory(order)
})

app.Event.On("order-created", func(e *application.CustomEvent) {
    order := e.Data.(*Order)
    logOrder(order)
})
```

### Padrão de solicitação/resposta

```go
// Frontend requests data
Emit("get-user-data", { userId: 123 })

// Backend responds
app.Event.On("get-user-data", func(e *application.CustomEvent) {
    data := e.Data.(map[string]interface{})
    userId := int(data["userId"].(float64))
    
    user := getUserFromDB(userId)
    
    // Send response
    app.Event.Emit("user-data-response", user)
})
```

```javascript
// Frontend receives response
Events.On("user-data-response", (event) => {
    const user = event.data
    displayUser(user)
})
```

**Observação:** para solicitação/resposta, **bindings são mais adequados**. Use eventos para notificações.

### Padrão de transmissão

```go
// Broadcast to all windows
app.Event.Emit("global-notification", "System update available")

// Each window handles it
Events.On("global-notification", (event) => {
    const message = event.data
    showNotification(message)
})
```

### Agregação de eventos

```go
type EventAggregator struct {
    events []Event
    mu     sync.Mutex
}

func (ea *EventAggregator) Add(event Event) {
    ea.mu.Lock()
    defer ea.mu.Unlock()
    
    ea.events = append(ea.events, event)
    
    // Emit batch every 100 events
    if len(ea.events) >= 100 {
        app.Event.Emit("event-batch", ea.events)
        ea.events = nil
    }
}
```

## Exemplo completo

**Go:**

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)

type NotificationService struct {
    app *application.App
}

func (n *NotificationService) Notify(message string) {
    // Emit to all windows
    n.app.Event.Emit("notification", map[string]interface{}{
        "message":   message,
        "timestamp": time.Now(),
    })
}

func main() {
    app := application.New(application.Options{
        Name: "Event Demo",
    })
    
    notifService := &NotificationService{app: app}
    
    // System events
    app.Event.OnApplicationEvent(events.Common.ThemeChanged, func(e *application.ApplicationEvent) {
        isDark := e.Context().IsDarkMode()
        app.Event.Emit("theme-changed", isDark)
    })
    
    // Custom events from frontend
    app.Event.On("user-action", func(e *application.CustomEvent) {
        data := e.Data.(map[string]interface{})
        action := data["action"].(string)
        
        app.Logger.Info("User action", "action", action)
        
        // Respond
        notifService.Notify("Action completed: " + action)
    })
    
    // Window events
    window := app.Window.New()
    
    window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
        app.Event.Emit("window-focused", window.Name())
    })
    
    window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
        // Confirm before close
        app.Event.Emit("confirm-close", nil)
        e.Cancel()  // Wait for confirmation
    })
    
    app.Run()
}
```

**JavaScript:**

```javascript
import { Events } from '@wailsio/runtime'

// Listen for notifications
Events.On("notification", (event) => {
    showNotification(event.data.message)
})

// Listen for theme changes
Events.On("theme-changed", (event) => {
    const isDark = event.data
    document.body.classList.toggle('dark', isDark)
})

// Listen for window focus
Events.On("window-focused", (event) => {
    const windowName = event.data
    console.log(`Window ${windowName} focused`)
})

// Handle close confirmation
Events.On("confirm-close", (event) => {
    if (confirm("Close window?")) {
        Events.Emit("close-confirmed", true)
    }
})

// Emit user actions
document.getElementById('button').addEventListener('click', () => {
    Events.Emit("user-action", { action: "button-clicked" })
})
```

## Práticas recomendadas

### ✅ Faça

- **Use eventos para notificações** — Comunicação unidirecional
- **Use bindings para solicitações** — Comunicação bidirecional
- **Mantenha os nomes de eventos consistentes** — Use kebab-case
- **Documente os dados dos eventos** — Quais campos estão incluídos?
- **Cancele a inscrição quando terminar** — Evite vazamentos de memória
- **Use hooks para validação** — Controle o fluxo de eventos

### ❌ Não faça

- **Não use eventos para RPC** — Use bindings
- **Não emita eventos com frequência excessiva** — Agrupe-os em lotes, se necessário
- **Não bloqueie a execução nos manipuladores** — Mantenha-os rápidos
- **Não se esqueça de cancelar a inscrição** — Isso causa vazamentos de memória
- **Não use eventos para grandes volumes de dados** — Use bindings
- **Não crie loops de eventos** — A emite B, B emite A

## Próximas etapas

@cards{cols="2"}
★ Guia de eventos
Saiba mais sobre padrões de eventos e geração de eventos com segurança de tipos.

[Saiba mais →](/guides/events-reference/)

---
◆ API de eventos
Explore a API completa de eventos da aplicação e das janelas.

[Saiba mais →](/reference/events/)

---
🚀 Bindings
Use bindings para solicitação/resposta.

[Saiba mais →](/features/bindings/methods/)

---
▣ Eventos de janela
Trate eventos do ciclo de vida das janelas.

[Saiba mais →](/features/windows/events/)

@end

---

**Dúvidas?** Pergunte no [Discord](https://discord.gg/JDdSxwjhGf) ou consulte os [exemplos de eventos](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/events).
