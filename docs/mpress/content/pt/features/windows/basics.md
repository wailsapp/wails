---
title: "Noções básicas de janelas"
description: "Criação e gerenciamento de janelas de aplicativos no Wails"
slug: "features/windows/basics"
sourcePath: "features/windows/basics.md"
---

## Gerenciamento de janelas

O Wails fornece uma **API unificada de gerenciamento de janelas** que funciona em todas as plataformas. Crie janelas, controle seu comportamento e gerencie várias janelas com controle total sobre criação, aparência, comportamento e ciclo de vida.

## Início rápido

```go
package main

import "github.com/wailsapp/wails/v3/pkg/application"

func main() {
    app := application.New(application.Options{
        Name: "My App",
    })

    // Create a window
    window := app.Window.New()
    
    // Configure it
    window.SetTitle("Hello Wails")
    window.SetSize(800, 600)
    window.Center()
    
    // Show it
    window.Show()

    app.Run()
}
```

**É só isso!** Agora você tem uma janela multiplataforma.

## Criação de janelas

### Janela básica

A maneira mais simples de criar uma janela:

```go
window := app.Window.New()
```

**O que você obtém:**

- Tamanho padrão (800x600)
- Título padrão (nome do aplicativo)
- WebView pronta para o seu frontend
- Aparência nativa da plataforma

### Janela com opções

Crie uma janela com configuração personalizada:

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "My Application",
    Width:  1200,
    Height: 800,
    X:      100,   // Position from left
    Y:      100,   // Position from top
    AlwaysOnTop: false,
    Frameless: false,
    Hidden: false,
    MinWidth: 400,
    MinHeight: 300,
    MaxWidth: 1920,
    MaxHeight: 1080,
})
```

**Opções comuns:**

| Opção | Tipo | Descrição |
| --- | --- | --- |
| `Title` | `string` | Título da janela |
| `Width` | `int` | Largura da janela em pixels |
| `Height` | `int` | Altura da janela em pixels |
| `X` | `int` | Posição X (a partir da esquerda) |
| `Y` | `int` | Posição Y (a partir do topo) |
| `AlwaysOnTop` | `bool` | Manter a janela acima das demais |
| `Frameless` | `bool` | Remover a barra de título e as bordas |
| `Hidden` | `bool` | Iniciar oculta |
| `MinWidth` | `int` | Largura mínima |
| `MinHeight` | `int` | Altura mínima |
| `MaxWidth` | `int` | Largura máxima |
| `MaxHeight` | `int` | Altura máxima |

**Consulte [Opções de janela](/features/windows/options/) para ver a lista completa.**

### Janelas nomeadas

Atribua nomes às janelas para recuperá-las facilmente:

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:  "main-window",
    Title: "Main Application",
})

// Later, find it by name
if mainWindow, ok := app.Window.GetByName("main-window"); ok {
    mainWindow.Show()
}
```

**Casos de uso:**

- Várias janelas (principal, configurações, sobre)
- Localização de janelas em diferentes partes do código
- Comunicação entre janelas

## Controle de janelas

### Exibir e ocultar

```go
// Show window
window.Show()

// Hide window
window.Hide()

// Check if visible
if window.IsVisible() {
    fmt.Println("Window is visible")
}
```

**Casos de uso:**

- Telas de abertura (exibir e depois ocultar)
- Janelas de configurações (ocultar quando não forem necessárias)
- Janelas pop-up (exibir sob demanda)

### Posição e tamanho

```go
// Set size
window.SetSize(1024, 768)

// Set position
window.SetPosition(100, 100)

// Centre on screen
window.Center()

// Get current size
width, height := window.Size()

// Get current position
x, y := window.Position()
```

**Sistema de coordenadas:**

- (0, 0) é o canto superior esquerdo da tela principal
- Valores positivos de X avançam para a direita
- Valores positivos de Y avançam para baixo

### Estado da janela

```go
// Minimise
window.Minimise()

// Maximise
window.Maximise()

// Fullscreen
window.Fullscreen()

// Restore to normal
window.Restore()

// Check state
if window.IsMinimised() {
    fmt.Println("Window is minimised")
}

if window.IsMaximised() {
    fmt.Println("Window is maximised")
}

if window.IsFullscreen() {
    fmt.Println("Window is fullscreen")
}
```

**Transições de estado:**

```
Normal ←→ Minimised
Normal ←→ Maximised
Normal ←→ Fullscreen
```

### Título e aparência

```go
// Set title
window.SetTitle("My Application - Document.txt")

// Set background colour — RGBA value (helper for RGB)
window.SetBackgroundColour(application.NewRGBA(0, 0, 0, 255))

// Set always on top
window.SetAlwaysOnTop(true)

// Set resizable
window.SetResizable(false)
```

### Fechamento de janelas

```go
// Close window — dispatches WindowClosing; a RegisterHook can call e.Cancel().
window.Close()
```

Não há um método `window.Destroy()` em v3 — use `Close()` e escute com `OnWindowEvent` (não é possível cancelar) ou intercepte com `RegisterHook` (é possível chamar `e.Cancel()` para manter a janela aberta).

## Localização de janelas

### Por nome

```go
if window, ok := app.Window.GetByName("settings"); ok {
    window.Show()
}
```

### Por ID

Cada janela tem um ID exclusivo:

```go
id := window.ID()
fmt.Printf("Window ID: %d\n", id)

// Find by ID
if found, ok := app.Window.GetByID(id); ok {
    found.Focus()
}
```

### Janela atual

Obtenha a janela que está com o foco:

```go
current := app.Window.Current()
if current != nil {
    current.SetTitle("Active Window")
}
```

### Todas as janelas

Obtenha todas as janelas:

```go
windows := app.Window.GetAll()
fmt.Printf("Total windows: %d\n", len(windows))

for _, w := range windows {
    fmt.Printf("Window: %s (ID: %d)\n", w.Name(), w.ID())
}
```

## Ciclo de vida da janela

### Criação

```go
app.Window.OnCreate(func(window application.Window) {
    fmt.Printf("Window created: %s\n", window.Name())

    // Configure new windows
    window.SetMinSize(400, 300)
})
```

### Fechamento

Para impedir que uma janela seja fechada, use `RegisterHook` com o evento `WindowClosing`:

```go
window.RegisterHook(events.Common.WindowClosing, func(event *application.WindowEvent) {
    if hasUnsavedChanges() {
        // Ask user for confirmation
        result := showConfirmDialog("Unsaved changes. Close anyway?")
        if result != "yes" {
            // Cancel the close event
            event.Cancel()
        }
    }
})
```

**Importante:** `RegisterHook` intercepta o evento de fechamento antes que ele ocorra. Chame `event.Cancel()` para impedir o fechamento da janela. Isso funciona para fechamentos iniciados pelo usuário (ao clicar no botão X).

### Destruição

Para executar a limpeza quando uma janela for fechada, use `OnWindowEvent` com o evento `WindowClosing`:

```go
window.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
    fmt.Println("Window is closing")
    // Cleanup resources
})
```

## Várias janelas

### Criação de várias janelas

```go
// Main window
mainWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:   "main",
    Title:  "Main Application",
    Width:  1200,
    Height: 800,
})

// Settings window
settingsWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:   "settings",
    Title:  "Settings",
    Width:  600,
    Height: 400,
    Hidden: true,  // Start hidden
})

// Show settings when needed
settingsWindow.Show()
```

### Comunicação entre janelas

As janelas podem se comunicar por meio de eventos:

```go
// In main window
app.Event.Emit("data-updated", map[string]interface{}{
    "value": 42,
})

// In settings window
app.Event.On("data-updated", func(event *application.CustomEvent) {
    data := event.Data.(map[string]interface{})
    value := data["value"].(int)
    fmt.Printf("Received: %d\n", value)
})
```

**Consulte [Eventos](/features/events/system/) para saber mais.**

### Janelas pai e filhas

`WebviewWindowOptions` não possui um campo `Parent`. Crie a janela filha como uma janela normal e anexe-a a uma janela pai como uma folha modal:

```go
// Create child window
childWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Child Window",
})

// Attach to the parent — presents as a sheet on macOS.
mainWindow.AttachModal(childWindow)
```

**Comportamento:**

- A janela filha permanece acima da janela pai.
- A janela filha é modal — bloqueia a interação com a janela pai.

**Compatibilidade com plataformas:**

- **macOS:** compatibilidade total (exibida como uma folha modal).
- **Windows:** não compatível.
- **Linux:** não compatível.

## Recursos específicos de cada plataforma

@tabs{sync-key="platform"}
[Windows]
**Recursos específicos do Windows:**

```go
// Flash taskbar button
window.Flash(true)  // Start flashing
window.Flash(false) // Stop flashing

// Trigger Windows 11 Snap Assist (Win+Z)
window.SnapAssist()
```

Não há `SetIcon` por janela — o ícone do aplicativo é definido no próprio aplicativo por meio de `app.SetIcon([]byte)` (ou, para um ícone de janela específico do Linux, pelo campo `application.LinuxWindow.Icon` durante a criação da janela).

**Snap Assist:** exibe as opções de layout de ajuste 11 do Windows por meio do atalho do sistema. Para usar um botão HTML personalizado de maximização com os Snap Layouts nativos ao passar o ponteiro, use [Regiões não clientes nativas no Windows](/features/windows/frameless/#native-non-client-regions-on-windows).

**Piscar na barra de tarefas:** útil para notificações quando a janela está minimizada.

[macOS]
**Recursos específicos do macOS:**

```go
// Transparent title bar
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Mac: application.MacWindow{
        TitleBar: application.MacTitleBar{
            AppearsTransparent: true,
        },
        Backdrop: application.MacBackdropTranslucent,
    },
})
```

**Tipos de plano de fundo:**

- `MacBackdropNormal` — Janela padrão
- `MacBackdropTranslucent` — Plano de fundo translúcido; uma **API privada é necessária** para a transparência da webview.
- `MacBackdropTransparent` — Totalmente transparente; uma **API privada é necessária** para a transparência da webview.
- `MacBackdropLiquidGlass` — Plano de fundo de vidro; uma **API privada é necessária** para a transparência da webview.

Compile com `-tags private_mac_apis` para que esses efeitos apareçam através da webview. Sem essa opção, a webview permanece opaca. O próprio `TitleBar.AppearsTransparent` usa APIs públicas. Consulte [APIs privadas do macOS](/guides/build/private-macos-apis/).

**Comportamento da coleção:** controle como as janelas se comportam entre os Spaces:

- `MacWindowCollectionBehaviorCanJoinAllSpaces` — Visível em todos os Spaces
- `MacWindowCollectionBehaviorFullScreenAuxiliary` — Pode sobrepor aplicativos em tela cheia

**Tela cheia nativa:** o modo de tela cheia do macOS cria um novo Space (área de trabalho virtual).

[Linux]
**Recursos específicos do Linux:**

```go
// Set window icon (per-window struct is LinuxWindow, not the app-level LinuxOptions)
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Linux: application.LinuxWindow{
        Icon: iconBytes,
    },
})
```

**Observações sobre ambientes de desktop:**

- GNOME: compatibilidade total
- KDE Plasma: compatibilidade total
- XFCE: compatibilidade parcial
- Outros: varia

**Gerenciadores de janelas em mosaico (Hyprland, Sway, i3 etc.):**

- `Minimise()` e `Maximise()` podem não funcionar como esperado — o gerenciador de janelas controla a geometria das janelas
- As solicitações de `SetSize()` e `SetPosition()` são apenas recomendações e podem ser ignoradas
- `Fullscreen()` geralmente funciona como esperado
- Alguns gerenciadores de janelas não oferecem suporte a janelas sempre visíveis

@end

## Padrões comuns

### Tela de abertura

```go
// Create splash screen
splash := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:     "Loading...",
    Width:     400,
    Height:    300,
    Frameless: true,
    AlwaysOnTop: true,
})

// Show splash
splash.Show()

// Initialise application
time.Sleep(2 * time.Second)

// Hide splash, show main window
splash.Close()
mainWindow.Show()
```

### Janela de configurações

```go
var settingsWindow *application.WebviewWindow

func showSettings() {
    if settingsWindow == nil {
        settingsWindow = app.Window.NewWithOptions(application.WebviewWindowOptions{
            Name:   "settings",
            Title:  "Settings",
            Width:  600,
            Height: 400,
        })
    }
    
    settingsWindow.Show()
    settingsWindow.Focus()
}
```

### Confirmar antes de fechar

```go
window.RegisterHook(events.Common.WindowClosing, func(event *application.WindowEvent) {
    if hasUnsavedChanges() {
        // Show dialog
        result := showConfirmDialog("Unsaved changes. Close anyway?")
        if result != "yes" {
            // Cancel the close event
            event.Cancel()
        }
    }
})
```

## Boas práticas

### ✅ Faça

- **Dê nomes às janelas importantes** — facilita encontrá-las depois
- **Defina um tamanho mínimo** — evita layouts inutilizáveis
- **Centralize as janelas** — oferece uma experiência melhor do que uma posição aleatória
- **Trate os eventos de fechamento** — evita a perda de dados
- **Teste em todas as plataformas** — o comportamento varia
- **Use tamanhos adequados** — considere diferentes tamanhos de tela

### ❌ Não faça

- **Não crie janelas demais** — isso confunde os usuários
- **Não se esqueça de fechar as janelas** — isso causa vazamentos de memória
- **Não fixe posições no código** — os tamanhos de tela variam
- **Não ignore as diferenças entre plataformas** — teste minuciosamente
- **Não bloqueie a thread da interface** — use goroutines para operações demoradas

## Solução de problemas

### A janela não aparece

**Possíveis causas:**

1. A janela foi criada como oculta
2. A janela está fora da tela
3. A janela está atrás de outras janelas

**Solução:**

```go
window.Show()
window.Center()
window.Focus()
```

### A janela está com o tamanho incorreto

**Causa:** dimensionamento de DPI no Windows/Linux

**Solução:**

```go
// Wails handles DPI automatically
// Just use logical pixels
window.SetSize(800, 600)
```

### A janela fecha imediatamente

**Causa:** o aplicativo é encerrado quando a última janela é fechada

**Solução:**

```go
app := application.New(application.Options{
    Mac: application.MacOptions{
        ApplicationShouldTerminateAfterLastWindowClosed: false,
    },
})
```

## Próximas etapas

@cards{cols="2"}
⚙ Opções de janela
Referência completa de todas as opções de janela.

[Saiba mais →](/features/windows/options/)

---
▣ Várias janelas
Padrões para aplicativos com várias janelas.

[Saiba mais →](/features/windows/multiple/)

---
★ Janelas sem moldura
Crie uma decoração de janela personalizada.

[Saiba mais →](/features/windows/frameless/)

---
🚀 Eventos de janela
Trate os eventos do ciclo de vida das janelas.

[Saiba mais →](/features/windows/events/)

@end

---

**Tem dúvidas?** Pergunte no [Discord](https://discord.gg/JDdSxwjhGf) ou consulte os [exemplos de janelas](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples).
