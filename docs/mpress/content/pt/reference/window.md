---
title: "API de janelas"
description: "Referência completa da API de janelas"
slug: "reference/window"
sourcePath: "reference/window.md"
---

## Visão geral

A API de janelas fornece métodos para controlar a aparência, o comportamento e o ciclo de vida das janelas. Acesse-a por meio das instâncias de janela ou do gerenciador `app.Window`.

`Window` é uma interface implementada por `*application.WebviewWindow`; as assinaturas dos métodos abaixo pertencem a `*WebviewWindow`. Muitos métodos modificadores retornam `Window` para permitir o encadeamento — os valores de retorno estão documentados em cada método.

**Operações comuns:**

- Criar e exibir janelas
- Controlar tamanho, posição e estado
- Tratar eventos de janela
- Gerenciar o conteúdo da janela
- Configurar aparência e comportamento

## Visibilidade

### Show()

Exibe a janela. Se ela estava oculta, torna-se visível. Retorna o receptor para permitir o encadeamento.

```go
func (w *WebviewWindow) Show() Window
```

**Exemplo:**

```go
window := app.Window.New()
window.Show()
```

### Hide()

Oculta a janela sem fechá-la. A janela permanece na memória e pode ser exibida novamente. Retorna o receptor para permitir o encadeamento.

```go
func (w *WebviewWindow) Hide() Window
```

**Exemplo:**

```go
// Hide window temporarily
window.Hide()

// Show it again later
window.Show()
```

**Casos de uso:**

- Aplicativos da bandeja do sistema que são ocultados na bandeja
- Fluxos de assistente nos quais as janelas são reutilizadas
- Ocultação temporária durante operações

### Close()

Fecha a janela. Isso dispara o evento `WindowClosing`.

```go
func (w *WebviewWindow) Close()
```

**Exemplo:**

```go
window.Close()
```

**Observação:** se um hook registrado chamar `event.Cancel()`, o fechamento será impedido.

## Propriedades da janela

### SetTitle()

Define o texto da barra de título da janela. Retorna o receptor para permitir o encadeamento.

```go
func (w *WebviewWindow) SetTitle(title string) Window
```

**Parâmetros:**

- `title` - O novo título da janela

**Exemplo:**

```go
window.SetTitle("My Application - Document.txt")
```

### Name()

Retorna o identificador de nome exclusivo da janela.

```go
func (w *WebviewWindow) Name() string
```

**Exemplo:**

```go
name := window.Name()
fmt.Println("Window name:", name)

// Retrieve window by name later
if w, ok := app.Window.GetByName(name); ok {
    w.Focus()
}
```

## Tamanho e posição

### SetSize()

Define as dimensões da janela em pixels. Retorna o receptor para permitir o encadeamento.

```go
func (w *WebviewWindow) SetSize(width, height int) Window
```

**Parâmetros:**

- `width` - Largura da janela em pixels
- `height` - Altura da janela em pixels

**Exemplo:**

```go
window.SetSize(1024, 768)
```

### Size()

Retorna as dimensões atuais da janela.

```go
func (w *WebviewWindow) Size() (width, height int)
```

**Exemplo:**

```go
width, height := window.Size()
fmt.Printf("Window is %dx%d\n", width, height)
```

### SetMinSize() / SetMaxSize()

Define as dimensões mínima e máxima da janela. Ambos retornam o receptor para permitir o encadeamento.

```go
func (w *WebviewWindow) SetMinSize(width, height int) Window
func (w *WebviewWindow) SetMaxSize(width, height int) Window
```

**Exemplo:**

```go
// Prevent window from being too small
window.SetMinSize(800, 600)

// Prevent window from being too large
window.SetMaxSize(1920, 1080)
```

### SetPosition()

Define a posição da janela em relação ao canto superior esquerdo da tela.

```go
func (w *WebviewWindow) SetPosition(x, y int)
```

**Parâmetros:**

- `x` - Posição horizontal em pixels
- `y` - Posição vertical em pixels

**Exemplo:**

```go
// Position window at top-left
window.SetPosition(0, 0)

// Position window 100px from top-left
window.SetPosition(100, 100)
```

### Position()

Retorna a posição atual da janela.

```go
func (w *WebviewWindow) Position() (x, y int)
```

**Exemplo:**

```go
x, y := window.Position()
fmt.Printf("Window is at (%d, %d)\n", x, y)
```

### Center()

Centraliza a janela na tela.

```go
func (w *WebviewWindow) Center()
```

**Exemplo:**

```go
window := app.Window.New()
window.Center()
window.Show()
```

**Observação:** centraliza no monitor principal. Para configurações com vários monitores, consulte as APIs de tela.

### Focus()

Traz a janela para a frente e atribui a ela o foco do teclado.

```go
func (w *WebviewWindow) Focus()
```

**Exemplo:**

```go
// Bring window to front
window.Focus()
```

## Estado da janela

### Minimise() / UnMinimise()

Minimiza a janela para a barra de tarefas ou o Dock, ou a restaura. `Minimise()` retorna o receptor para permitir o encadeamento; `UnMinimise()` não retorna nada.

```go
func (w *WebviewWindow) Minimise() Window
func (w *WebviewWindow) UnMinimise()
```

**Exemplo:**

```go
// Minimise window
window.Minimise()

// Restore from minimised state
window.UnMinimise()
```

### Maximise() / UnMaximise()

Maximiza a janela para preencher a tela ou a restaura ao tamanho anterior. `Maximise()` retorna o receptor para permitir o encadeamento; `UnMaximise()` não retorna nada.

```go
func (w *WebviewWindow) Maximise() Window
func (w *WebviewWindow) UnMaximise()
```

**Exemplo:**

```go
// Maximise window
window.Maximise()

// Restore to previous size
window.UnMaximise()
```

### Fullscreen() / UnFullscreen() / ToggleFullscreen()

Entra no modo de tela cheia ou sai dele. `Fullscreen()` retorna o receptor para permitir o encadeamento.

```go
func (w *WebviewWindow) Fullscreen() Window
func (w *WebviewWindow) UnFullscreen()
func (w *WebviewWindow) ToggleFullscreen()
```

**Exemplo:**

```go
// Enter fullscreen
window.Fullscreen()

// Exit fullscreen
window.UnFullscreen()

// Or toggle
window.ToggleFullscreen()
```

Não há um método `SetFullscreen(bool)`.

### IsMinimised() / IsMaximised() / IsFullscreen()

Verifica o estado atual da janela.

```go
func (w *WebviewWindow) IsMinimised() bool
func (w *WebviewWindow) IsMaximised() bool
func (w *WebviewWindow) IsFullscreen() bool
```

**Exemplo:**

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

## Conteúdo da janela

### SetURL()

Navega até uma URL específica na janela. Retorna o receptor para permitir o encadeamento de chamadas.

```go
func (w *WebviewWindow) SetURL(url string) Window
```

**Parâmetros:**

- `url` — URL de destino (pode ser `http://wails.localhost/` para recursos incorporados)

**Exemplo:**

```go
// Navigate to embedded page
window.SetURL("http://wails.localhost/settings.html")

// Navigate to external URL (if allowed)
window.SetURL("https://wails.io")
```

### SetHTML()

Define o conteúdo da janela diretamente a partir de uma string HTML. Retorna o receptor para permitir o encadeamento de chamadas.

```go
func (w *WebviewWindow) SetHTML(html string) Window
```

**Parâmetros:**

- `html` — conteúdo HTML a ser exibido

**Exemplo:**

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

**Casos de uso:**

- Geração dinâmica de conteúdo
- Janelas simples sem um processo de build do frontend
- Páginas de erro ou telas de abertura

### Reload()

Recarrega o conteúdo atual da janela.

```go
func (w *WebviewWindow) Reload()
```

**Exemplo:**

```go
// Reload current page
window.Reload()
```

**Observação:** Útil durante o desenvolvimento ou quando o conteúdo precisa ser atualizado.

## Eventos da janela

O Wails fornece dois métodos para tratar eventos da janela:

- **OnWindowEvent()** — monitora eventos da janela (não pode impedi-los).
- **RegisterHook()** — intercepta eventos da janela (pode impedi-los chamando `event.Cancel()`).

### OnWindowEvent()

Registra uma função de callback para eventos da janela. Retorna uma função para cancelar a inscrição.

```go
func (w *WebviewWindow) OnWindowEvent(
    eventType events.WindowEventType,
    callback func(event *WindowEvent),
) func()
```

**Exemplo:**

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

**Eventos comuns da janela:**

- `events.Common.WindowClosing` — a janela está prestes a fechar
- `events.Common.WindowFocus` — a janela ganhou foco
- `events.Common.WindowLostFocus` — a janela perdeu o foco
- `events.Common.WindowDidMove` — a janela foi movida
- `events.Common.WindowDidResize` — a janela foi redimensionada
- `events.Common.WindowMinimise` — a janela foi minimizada
- `events.Common.WindowMaximise` — a janela foi maximizada
- `events.Common.WindowFullscreen` — a janela entrou no modo de tela cheia
- `events.Common.WindowRuntimeReady` — o runtime da janela foi inicializado

### RegisterHook()

Registra um hook para eventos da janela. Os hooks são executados antes dos listeners e podem impedir o evento chamando `event.Cancel()`. Retorna uma função para cancelar a inscrição.

```go
func (w *WebviewWindow) RegisterHook(
    eventType events.WindowEventType,
    callback func(event *WindowEvent),
) func()
```

**Exemplo — impedir o fechamento da janela:**

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

**Exemplo — salvar antes de fechar:**

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

Emite um evento personalizado para o frontend da janela. Retorna `true` se a emissão tiver sido cancelada por um hook.

```go
func (w *WebviewWindow) EmitEvent(name string, data ...any) bool
```

**Parâmetros:**

- `name` — nome do evento
- `data` — dados opcionais a serem enviados com o evento

**Exemplo:**

```go
// Send data to specific window
window.EmitEvent("data-updated", map[string]any{
    "count":  42,
    "status": "success",
})
```

**Frontend (JavaScript):**

```javascript
import { Events } from '@wailsio/runtime'

Events.On('data-updated', (data) => {
    console.log('Count:', data.count)
    console.log('Status:', data.status)
})
```

## Outros métodos

### SetEnabled()

Habilita ou desabilita a interação do usuário com a janela.

```go
func (w *WebviewWindow) SetEnabled(enabled bool)
```

**Exemplo:**

```go
// Disable window during long operation
window.SetEnabled(false)

// Perform operation
performLongOperation()

// Re-enable window
window.SetEnabled(true)
```

### SetBackgroundColour()

Define a cor de fundo da janela (exibida antes do carregamento do conteúdo). Retorna o receptor para permitir o encadeamento de chamadas.

```go
func (w *WebviewWindow) SetBackgroundColour(colour RGBA) Window
```

`RGBA` é `application.RGBA{Red, Green, Blue, Alpha uint8}`. Use as funções auxiliares `application.NewRGB(r, g, b)` (alfa 255) ou `application.NewRGBA(r, g, b, a)`.

**Exemplo:**

```go
// White background
window.SetBackgroundColour(application.NewRGB(255, 255, 255))

// Dark background with full alpha
window.SetBackgroundColour(application.NewRGBA(30, 30, 30, 255))
```

### SetResizable()

Controla se a janela pode ser redimensionada pelo usuário. Retorna o receptor para permitir o encadeamento de chamadas.

```go
func (w *WebviewWindow) SetResizable(resizable bool) Window
```

**Exemplo:**

```go
// Make window fixed size
window.SetResizable(false)
```

### SetAlwaysOnTop()

Define se a janela permanece acima das demais janelas. Retorna o receptor para permitir o encadeamento de chamadas.

```go
func (w *WebviewWindow) SetAlwaysOnTop(alwaysOnTop bool) Window
```

**Exemplo:**

```go
// Keep window on top
window.SetAlwaysOnTop(true)
```

### Print()

Abre a caixa de diálogo nativa de impressão para o conteúdo da janela.

```go
func (w *WebviewWindow) Print() error
```

**Retorno:** um erro se a impressão falhar.

**Exemplo:**

```go
if err := window.Print(); err != nil {
    log.Println("Print failed:", err)
}
```

### AttachModal()

Anexa uma segunda janela como uma folha modal.

```go
func (w *WebviewWindow) AttachModal(modalWindow Window)
```

**Parâmetros:**

- `modalWindow` - A janela a ser anexada como modal

**Compatibilidade com plataformas:**

- **macOS**: Compatibilidade total (exibida como uma folha)
- **Windows**: Sem compatibilidade
- **Linux**: Sem compatibilidade

**Exemplo:**

```go
modalWindow := app.Window.New()
window.AttachModal(modalWindow)
```

## Opções específicas da plataforma

### Linux

As janelas no Linux são compatíveis com as seguintes opções específicas da plataforma por meio de `LinuxWindow`:

#### MenuStyle

Controla como o menu do aplicativo é exibido. Essa opção está disponível na compilação GTK4 padrão e é ignorada nas compilações legadas de `-tags gtk3`.

| Valor | Descrição |
| --- | --- |
| `LinuxMenuStyleMenuBar` | Barra de menus tradicional abaixo da barra de título (padrão) |
| `LinuxMenuStylePrimaryMenu` | Botão do menu principal na barra de cabeçalho (estilo GNOME) |

**Exemplo:**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "My Application",
    Linux: application.LinuxWindow{
        MenuStyle: application.LinuxMenuStylePrimaryMenu,
    },
})
window.SetMenu(menu)
```

**Observação:** O estilo de menu principal exibe um botão de menu hambúrguer (☰) na barra de cabeçalho, de acordo com as Diretrizes de Interface Humana do GNOME. Esse é o estilo recomendado para aplicativos GNOME modernos.

## Exemplo completo

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
