---
title: "API de eventos"
description: "Referência completa da API de eventos"
slug: "reference/events"
sourcePath: "reference/events.md"
---

## Visão geral

A API de eventos fornece métodos para emitir e escutar eventos, permitindo a comunicação entre diferentes partes do aplicativo.

**Tipos de eventos:**

- **Eventos do aplicativo** — Eventos do ciclo de vida do aplicativo (inicialização, encerramento)
- **Eventos de janela** — Alterações no estado da janela (obtenção de foco, perda de foco, redimensionamento)
- **Eventos personalizados** — Eventos definidos pelo usuário para comunicação específica do aplicativo

**Padrões de comunicação:**

- **Do Go para o frontend** — Emita eventos no Go e escute-os no JavaScript
- **Do frontend para o Go** — Não diretamente (use vinculações de serviço)
- **De frontend para frontend** — Por meio do Go ou de eventos locais do runtime
- **De janela para janela** — Direcione a janelas específicas ou transmita para todas

## Métodos de eventos (Go)

### app.Event.Emit()

Emite um evento personalizado para todas as janelas. Retorna `true` se um hook cancelar a emissão.

```go
func (em *EventManager) Emit(name string, data ...any) bool
```

**Parâmetros:**

- `name` — Nome do evento
- `data` — Dados opcionais a serem enviados com o evento

**Exemplo:**

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

Escuta eventos personalizados no Go.

```go
func (em *EventManager) On(name string, callback func(*CustomEvent)) func()
```

**Parâmetros:**

- `name` — Nome do evento a ser escutado
- `callback` — Função chamada quando o evento é emitido

**Retorno:** função de limpeza que remove o listener do evento

**Exemplo:**

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

### Eventos específicos de uma janela

Emita eventos para uma janela específica:

```go
// Emit to specific window
window.EmitEvent("notification", "Hello from Go!")

// Emit to all windows
app.Event.Emit("global-update", data)
```

## Métodos de eventos (frontend)

### On()

Escuta eventos provenientes do Go.

```javascript
import { Events } from '@wailsio/runtime'

Events.On(eventName, callback)
```

**Parâmetros:**

- `eventName` — Nome do evento a ser escutado
- `callback` — Função chamada quando o evento é recebido

**Retorno:** função de limpeza

**Exemplo:**

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

Escuta uma única ocorrência de um evento.

```javascript
import { Events } from '@wailsio/runtime'

Events.Once(eventName, callback)
```

**Exemplo:**

```javascript
import { Events } from '@wailsio/runtime'

// Listen for first occurrence only
Events.Once('initialization-complete', (data) => {
    console.log('App initialized!', data)
    // This will only fire once
})
```

### Off()

Remove todos os listeners de um ou mais eventos. `Off` aceita strings variádicas com nomes de eventos — **não** aceita um callback. Para remover um único listener, armazene a função de cancelamento de inscrição retornada por `Events.On(...)` e chame-a.

```typescript
import { Events } from '@wailsio/runtime'

Events.Off(...eventNames: string[]): void
```

**Exemplo:**

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

Remove **todos** os listeners de eventos. Não aceita argumentos.

```typescript
import { Events } from '@wailsio/runtime'

Events.OffAll(): void
```

**Exemplo:**

```javascript
import { Events } from '@wailsio/runtime'

// Remove all listeners — typically used during teardown.
Events.OffAll()
```

### OnMultiple()

Escuta um evento até `max` vezes e depois cancela automaticamente a inscrição.

```typescript
Events.OnMultiple(eventName: string, callback, max: number): () => void
```

**Exemplo:**

```javascript
Events.OnMultiple('progress', (data) => {
    console.log('progress', data)
}, 5)
```

## Eventos do aplicativo

### app.Event.OnApplicationEvent()

Escuta eventos do ciclo de vida do aplicativo.

```go
func (em *EventManager) OnApplicationEvent(
    eventType events.ApplicationEventType,
    callback func(*ApplicationEvent),
) func()
```

As constantes de eventos ficam no pacote `events` — `events.Common.*` para eventos multiplataforma e `events.Mac.*` / `events.Windows.*` / `events.Linux.*` para eventos específicos de cada plataforma.

**Eventos comuns do aplicativo:**

- `events.Common.ApplicationStarted` — O aplicativo concluiu a inicialização.
- `events.Common.ThemeChanged` — O tema do sistema alternou entre claro e escuro.
- `events.Common.ApplicationOpenedWithFile` — Iniciado por meio de uma associação de arquivo.
- `events.Common.ApplicationLaunchedWithUrl` — Iniciado por meio de um esquema de URL.

**Não** existe uma constante genérica para o evento de "encerramento do aplicativo" — registre a limpeza de encerramento por meio de `application.Options.OnShutdown` ou `app.OnShutdown(func())`.

**Exemplo:**

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

## Eventos de janela

### OnWindowEvent()

Escuta eventos específicos de uma janela.

```go
func (w *WebviewWindow) OnWindowEvent(
    eventType events.WindowEventType,
    callback func(*WindowEvent),
) func()
```

As constantes de eventos de janela ficam no pacote `events`: `events.Common.*` para eventos multiplataforma (e `events.Mac.*` / `events.Windows.*` / `events.Linux.*` para eventos específicos de cada plataforma).

**Eventos comuns de janela:**

- `events.Common.WindowFocus` — A janela obteve o foco.
- `events.Common.WindowLostFocus` - A janela perdeu o foco.
- `events.Common.WindowClosing` - A janela está prestes a fechar (pode ser cancelado por meio de `RegisterHook`).
- `events.Common.WindowDidResize` - A janela foi redimensionada.
- `events.Common.WindowDidMove` - A janela foi movida.
- `events.Common.WindowMinimise` / `WindowUnMinimise` / `WindowMaximise` / `WindowUnMaximise` / `WindowFullscreen` / `WindowUnFullscreen`.
- `events.Common.WindowRuntimeReady` - É seguro emitir eventos para o runtime dentro da janela.

**Exemplo:**

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

## Padrões comuns

Estes padrões demonstram abordagens comprovadas para usar eventos em aplicações reais. Cada padrão resolve um desafio específico de comunicação entre o backend em Go e o frontend, ajudando você a criar aplicações responsivas e bem estruturadas.

### Padrão de solicitação e resposta

Use este padrão quando quiser notificar o frontend sobre a conclusão de operações do backend, como após a busca de dados, o processamento de arquivos ou tarefas em segundo plano. O binding do serviço retorna os dados diretamente, enquanto os eventos fornecem notificações adicionais para atualizar a interface, como exibir mensagens toast ou atualizar listas.

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

### Atualizações de progresso

Ideal para operações demoradas, como uploads de arquivos, processamento em lote, importações de grandes volumes de dados ou codificação de vídeo. Emita eventos de progresso durante a operação para atualizar barras de progresso, textos de status ou indicadores de etapas na interface, fornecendo feedback em tempo real aos usuários.

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

### Comunicação entre várias janelas

Perfeito para aplicações com várias janelas, como painéis de configurações, dashboards ou visualizadores de documentos. Transmita eventos para sincronizar o estado entre todas as janelas (alterações de tema, preferências do usuário) ou envie eventos direcionados a janelas específicas para atualizações próprias de cada janela.

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

### Sincronização de estado

Use quando precisar manter sincronizados os estados do frontend e do backend, como em sessões de usuário, configurações da aplicação ou recursos colaborativos. Quando o estado mudar no backend, emita eventos para atualizar todos os frontends conectados, garantindo a consistência em toda a aplicação.

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

### Notificações orientadas a eventos

Ideal para exibir feedback ao usuário, como confirmações de sucesso, alertas de erro ou mensagens informativas. Em vez de chamar o código da interface diretamente nos serviços, emita eventos de notificação que o frontend processe de maneira consistente, facilitando a alteração dos estilos das notificações ou a adição de recursos como um histórico de notificações.

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

## Exemplo completo

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

## Eventos integrados

O Wails fornece eventos de sistema integrados para o ciclo de vida da aplicação e das janelas. Esses eventos são emitidos automaticamente pelo framework.

### Eventos comuns versus eventos nativos da plataforma

O Wails fornece dois tipos de eventos de sistema:

**Eventos comuns** (`events.Common.*`) são abstrações multiplataforma que funcionam de maneira consistente no macOS, Windows e Linux. Estes são os eventos que você deve usar na aplicação para obter a máxima portabilidade.

**Eventos nativos da plataforma** (`events.Mac.*`, `events.Windows.*`, `events.Linux.*`) são os eventos subjacentes específicos do sistema operacional dos quais os Eventos comuns são mapeados. Eles fornecem acesso a comportamentos e casos extremos específicos da plataforma.

**Como funcionam:**

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

**Mapeamento de eventos:**

Os eventos nativos da plataforma são mapeados automaticamente para Eventos comuns:

- macOS: `events.Mac.WindowShouldClose` → `events.Common.WindowClosing`
- Windows: `events.Windows.WindowClosing` → `events.Common.WindowClosing`
- Linux: `events.Linux.WindowDeleteEvent` → `events.Common.WindowClosing`

Esse mapeamento ocorre automaticamente em segundo plano. Portanto, ao escutar `events.Common.WindowClosing`, você receberá esse evento independentemente da plataforma.

**Quando usar cada tipo:**

- **Use Eventos comuns** em 99% do código da aplicação — eles proporcionam um comportamento consistente entre plataformas
- **Use Eventos nativos da plataforma** somente quando precisar de uma funcionalidade específica da plataforma que não esteja disponível nos Eventos comuns (por exemplo, eventos do ciclo de vida das janelas específicos do macOS ou eventos de gerenciamento de energia do Windows)

### Eventos da aplicação

| Evento | Descrição | Quando é emitido | Cancelável |
| --- | --- | --- | --- |
| `ApplicationOpenedWithFile` | Aplicação aberta com um arquivo | Quando o aplicativo é iniciado com um arquivo (por exemplo, por associação de arquivo) | Não |
| `ApplicationStarted` | O aplicativo concluiu a inicialização | Depois que a inicialização do aplicativo é concluída e ele está pronto | Não |
| `ApplicationLaunchedWithUrl` | Aplicativo iniciado com uma URL | Quando o aplicativo é iniciado por meio de um esquema de URL | Não |
| `ThemeChanged` | O tema do sistema foi alterado | Quando o tema do sistema operacional alterna entre os modos claro e escuro | Não |
| `SystemWillSleep` | O sistema está prestes a ser suspenso | Pouco antes de o sistema operacional ser suspenso (macOS / Windows / Linux com logind) | Não |
| `SystemDidWake` | O sistema retomou após a suspensão | Imediatamente após sair do modo de suspensão (macOS / Windows / Linux com logind) | Não |

**Uso:**

```go
import "github.com/wailsapp/wails/v3/pkg/events"

app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(e *application.ApplicationEvent) {
    app.Logger.Info("Application ready!")
})

app.Event.OnApplicationEvent(events.Common.ThemeChanged, func(e *application.ApplicationEvent) {
    // Update app theme
})
```

### Eventos de janela

| Evento | Descrição | Quando é emitido | Cancelável |
| --- | --- | --- | --- |
| `WindowClosing` | A janela está prestes a ser fechada | Antes de a janela ser fechada (o usuário clicou no X ou Close() foi chamado) | Sim |
| `WindowDidMove` | A janela foi movida para uma nova posição | Depois que a posição da janela é alterada (com debounce) | Não |
| `WindowDidResize` | A janela foi redimensionada | Depois que o tamanho da janela é alterado | Não |
| `WindowDPIChanged` | A escala de DPI da janela foi alterada | Ao mover a janela entre monitores com DPIs diferentes (Windows) | Não |
| `WindowFilesDropped` | Arquivos soltos por meio do recurso nativo de arrastar e soltar do sistema operacional | Depois que arquivos são arrastados do sistema operacional e soltos sobre a janela | Não |
| `WindowFocus` | A janela recebeu foco | Quando a janela se torna ativa | Não |
| `WindowFullscreen` | A janela entrou no modo de tela cheia | Depois que Fullscreen() é chamado ou o usuário entra no modo de tela cheia | Não |
| `WindowHide` | A janela foi ocultada | Depois que Hide() é chamado ou a janela fica obstruída | Não |
| `WindowLostFocus` | A janela perdeu o foco | Quando a janela se torna inativa | Não |
| `WindowMaximise` | A janela foi maximizada | Depois que Maximise() é chamado ou o usuário maximiza a janela | Sim (macOS) |
| `WindowMinimise` | A janela foi minimizada | Depois que Minimise() é chamado ou o usuário minimiza a janela | Sim (macOS) |
| `WindowRestore` | A janela foi restaurada do estado minimizado ou maximizado | Depois que Restore() é chamado (principalmente no Windows) | Não |
| `WindowRuntimeReady` | O runtime do Wails foi carregado e está pronto | Quando a inicialização do runtime JavaScript é concluída | Não |
| `WindowShow` | A janela ficou visível | Após Show() ou quando a janela se torna visível | Não |
| `WindowUnFullscreen` | A janela saiu do modo de tela cheia | Após UnFullscreen() ou quando o usuário sai do modo de tela cheia | Não |
| `WindowUnMaximise` | A janela saiu do estado maximizado | Após UnMaximise() ou quando o usuário desmaximiza a janela | Sim (macOS) |
| `WindowUnMinimise` | A janela saiu do estado minimizado | Após UnMinimise()/Restore() ou quando o usuário restaura a janela | Sim (macOS) |
| `WindowZoomIn` | O zoom do conteúdo da janela aumentou | Após a chamada de ZoomIn() (principalmente no macOS) | Sim (macOS) |
| `WindowZoomOut` | O zoom do conteúdo da janela diminuiu | Após a chamada de ZoomOut() (principalmente no macOS) | Sim (macOS) |
| `WindowZoomReset` | O zoom do conteúdo da janela foi redefinido para 100% | Após a chamada de ZoomReset() (principalmente no macOS) | Sim (macOS) |
| `WindowDropZoneFilesDropped` | Arquivos soltos em uma zona de soltura definida em JS | Quando arquivos são soltos em um elemento com uma zona de soltura | Não |

**Uso:**

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

**Observações importantes:**

- **WindowRuntimeReady** é essencial — aguarde esse evento antes de emitir eventos para o frontend
- **WindowDidMove** e **WindowDidResize** usam debounce (50 ms por padrão) para evitar uma sobrecarga de eventos
- **Eventos canceláveis** podem ser impedidos chamando `event.Cancel()` em um manipulador de `RegisterHook()`
- **WindowFilesDropped** destina-se a solturas de arquivos nativas do sistema operacional; **WindowDropZoneFilesDropped** destina-se a zonas de soltura baseadas na Web
- Alguns eventos são específicos da plataforma (por exemplo, WindowDPIChanged no Windows e eventos de zoom principalmente no macOS)

## Convenções de nomenclatura de eventos

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

## Considerações de desempenho

### Aplicação de debounce a eventos de alta frequência

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

### Limitação da frequência de eventos

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
