---
title: "API de gerenciadores"
description: "Estrutura de API organizada com interfaces de gerenciadores especializadas"
slug: "concepts/manager-api"
sourcePath: "concepts/manager-api.md"
---

A API de gerenciadores do Wails v3 oferece uma maneira organizada e fácil de explorar para acessar as funcionalidades do aplicativo por meio de structs de gerenciadores especializados, agrupados em campos públicos de `*application.App`. O Wails 3 rompe completamente com a v2 — não há uma camada de encapsulamento por chamada para manter a compatibilidade com a antiga API no estilo `app.NewWebviewWindow(...)`; portanto, os gerenciadores abaixo são a forma de controlar o aplicativo.

## Visão geral

A API de gerenciadores organiza as funcionalidades do aplicativo em doze áreas especializadas (um logger e onze gerenciadores):

- **`app.Window`** - Criação e gerenciamento de janelas e callbacks
- **`app.ContextMenu`** - Registro e gerenciamento de menus de contexto\
- **`app.KeyBinding`** - Gerenciamento de atalhos de teclado globais
- **`app.Browser`** - Integração com o navegador (abertura de URLs e arquivos)
- **`app.Env`** - Informações do ambiente e estado do sistema
- **`app.Dialog`** - Operações com caixas de diálogo de arquivos e mensagens
- **`app.Event`** - Tratamento de eventos personalizados e eventos do aplicativo
- **`app.Menu`** - Gerenciamento do menu do aplicativo
- **`app.Screen`** - Gerenciamento de telas e transformações de coordenadas
- **`app.Clipboard`** - Operações de texto na área de transferência
- **`app.SystemTray`** - Criação e gerenciamento do ícone da bandeja do sistema
- **`app.Autostart`** - Registro do aplicativo para iniciar quando o usuário fizer login

## Benefícios

- **Maior facilidade de descoberta** - O preenchimento automático da IDE exibe uma superfície de API organizada
- **Melhor organização do código** - Os métodos relacionados ficam agrupados
- **Maior facilidade de manutenção** - Separação de responsabilidades entre os gerenciadores
- **Extensibilidade futura** - Maior facilidade para adicionar novos recursos a áreas específicas

## Uso

A API de gerenciadores oferece acesso organizado a todas as funcionalidades do aplicativo:

```go
// Events and custom event handling
app.Event.Emit("custom", data)
app.Event.On("custom", func(e *CustomEvent) { ... })

// Window management
window, _ := app.Window.GetByName("main")
app.Window.OnCreate(func(window Window) { ... })

// Browser integration
app.Browser.OpenURL("https://wails.io")

// Menu management
menu := app.Menu.New()
app.Menu.Set(menu)

// System tray
systray := app.SystemTray.New()
```

## Referência dos gerenciadores

### Gerenciador de janelas

Gerencia a criação e a obtenção de janelas, além dos callbacks do ciclo de vida.

```go
// Create windows
window := app.Window.New()
window := app.Window.NewWithOptions(options)
current := app.Window.Current()

// Find windows
window, exists := app.Window.GetByName("main")
windows := app.Window.GetAll()

// Window callbacks
app.Window.OnCreate(func(window Window) {
    // Handle window creation
})
```

### Gerenciador de eventos

Trata eventos personalizados e o monitoramento de eventos do aplicativo.

```go
// Custom events
app.Event.Emit("userAction", data)
cancelFunc := app.Event.On("userAction", func(e *CustomEvent) {
    // Handle event
})
app.Event.Off("userAction")
app.Event.Reset() // Remove all listeners

// Application events
app.Event.OnApplicationEvent(events.Common.ThemeChanged, func(e *ApplicationEvent) {
    // Handle system theme change
})
```

### Gerenciador do navegador

Oferece integração com o navegador para abrir URLs e arquivos.

```go
// Open URLs and files in default browser
err := app.Browser.OpenURL("https://wails.io")
err := app.Browser.OpenFile("/path/to/document.pdf")
```

### Gerenciador de ambiente

Fornece acesso às informações do ambiente do sistema.

```go
// Get environment info
env := app.Env.Info()
fmt.Printf("OS: %s, Arch: %s\n", env.OS, env.Arch)

// Check system theme
if app.Env.IsDarkMode() {
    // Dark mode is active
}

// Open file manager
err := app.Env.OpenFileManager("/path/to/folder", false)
```

### Gerenciador de caixas de diálogo

Oferece acesso organizado a caixas de diálogo de arquivos e mensagens.

```go
// File dialogs
result, err := app.Dialog.OpenFile().
    AddFilter("Text Files", "*.txt").
    PromptForSingleSelection()

result, err = app.Dialog.SaveFile().
    SetFilename("document.txt").
    PromptForSingleSelection()

// Message dialogs
app.Dialog.Info().
    SetTitle("Information").
    SetMessage("Operation completed successfully").
    Show()

app.Dialog.Error().
    SetTitle("Error").
    SetMessage("An error occurred").
    Show()
```

### Gerenciador de menus

Criação e gerenciamento do menu do aplicativo.

```go
// Create and set application menu
menu := app.Menu.New()
fileMenu := menu.AddSubmenu("File")
fileMenu.Add("New").OnClick(func(ctx *Context) {
    // Handle menu click
})

app.Menu.Set(menu)

// Show about dialog
app.Menu.ShowAbout()
```

### Gerenciador de atalhos de teclado

Gerenciamento dinâmico de atalhos de teclado globais.

```go
// Add key bindings
app.KeyBinding.Add("ctrl+n", func(window application.Window) {
    // Handle Ctrl+N
})

app.KeyBinding.Add("ctrl+q", func(window application.Window) {
    app.Quit()
})

// Remove key bindings
app.KeyBinding.Remove("ctrl+n")

// Get all bindings
bindings := app.KeyBinding.GetAll()
```

### Gerenciador de menus de contexto

Gerenciamento avançado de menus de contexto (para autores de bibliotecas).

```go
// Create and register context menu
menu := app.ContextMenu.New()
app.ContextMenu.Add("myMenu", menu)

// Retrieve context menu
menu, exists := app.ContextMenu.Get("myMenu")

// Remove context menu
app.ContextMenu.Remove("myMenu")
```

### Gerenciador de telas

Gerenciamento de telas e transformações de coordenadas para configurações com vários monitores.

```go
// Get screen information
screens := app.Screen.GetAll()
primary := app.Screen.GetPrimary()

// Coordinate transformations
physicalPoint := app.Screen.DipToPhysicalPoint(logicalPoint)
logicalPoint := app.Screen.PhysicalToDipPoint(physicalPoint)

// Screen detection
screen := app.Screen.ScreenNearestDipPoint(point)
screen = app.Screen.ScreenNearestDipRect(rect)
```

### Gerenciador da área de transferência

Operações da área de transferência para leitura e gravação de texto.

```go
// Set text to clipboard
success := app.Clipboard.SetText("Hello World")
if !success {
    // Handle error
}

// Get text from clipboard
text, ok := app.Clipboard.Text()
if !ok {
    // Handle error
} else {
    // Use the text
}
```

### Gerenciador SystemTray

Criação e gerenciamento do ícone da bandeja do sistema.

```go
// Create system tray
systray := app.SystemTray.New()
systray.SetLabel("My App")
systray.SetIcon(iconBytes)

// Add menu to system tray
menu := app.Menu.New()
menu.Add("Open").OnClick(func(ctx *Context) {
    // Handle click
})
systray.SetMenu(menu)

// Destroy system tray when done
systray.Destroy()
```

### Gerenciador de inicialização automática

Registra o aplicativo para iniciar quando o usuário fizer login. Seleciona o mecanismo nativo adequado para cada plataforma: SMAppService ou um plist de LaunchAgent no macOS, a chave de Registro `HKCU\…\Run` no Windows e uma entrada XDG `.desktop` no Linux.

```go
// Register to launch at login
err := app.Autostart.Enable()

// With extra launch-time arguments and a custom identifier
err = app.Autostart.EnableWithOptions(application.AutostartOptions{
    Identifier: "com.example.myapp",
    Arguments:  []string{"--hidden"},
})

// Check / remove
enabled, err := app.Autostart.IsEnabled()
status, err := app.Autostart.Status()  // includes Path + Strategy
err = app.Autostart.Disable()
```

Consulte a [página do recurso de inicialização automática](/features/autostart/basics/) para conhecer o comportamento em cada plataforma, as regras de identificadores e a garantia de detecção de registros obsoletos.
