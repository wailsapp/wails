---
title: "API da aplicação"
description: "Referência completa da API da aplicação"
slug: "reference/application"
sourcePath: "reference/application.md"
---

## Visão geral

`Application` é o núcleo do seu aplicativo Wails. Ele gerencia janelas, serviços e eventos e fornece acesso a todos os recursos da plataforma.

## Criação de uma aplicação

```go
import "github.com/wailsapp/wails/v3/pkg/application"

app := application.New(application.Options{
    Name:        "My App",
    Description: "My awesome application",
    Services: []application.Service{
        application.NewService(&MyService{}),
    },
})
```

## Métodos principais

### Run()

Inicia o loop de eventos da aplicação.

```go
func (a *App) Run() error
```

**Exemplo:**

```go
err := app.Run()
if err != nil {
    log.Fatal(err)
}
```

**Retorna:** um erro se a inicialização falhar

### Quit()

Encerra a aplicação de forma controlada.

```go
func (a *App) Quit()
```

**Exemplo:**

```go
// In a menu handler
menu.Add("Quit").OnClick(func(ctx *application.Context) {
    app.Quit()
})
```

### Config()

Retorna a configuração da aplicação.

```go
func (a *App) Config() Options
```

**Exemplo:**

```go
config := app.Config()
fmt.Println("App name:", config.Name)
```

## Gerenciamento de janelas

### app.Window.New()

Cria uma nova janela de webview com as opções padrão.

```go
func (wm *WindowManager) New() *WebviewWindow
```

**Exemplo:**

```go
window := app.Window.New()
window.Show()
```

### app.Window.NewWithOptions()

Cria uma nova janela de webview com opções personalizadas.

```go
func (wm *WindowManager) NewWithOptions(options WebviewWindowOptions) *WebviewWindow
```

**Exemplo:**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "My Window",
    Width:  800,
    Height: 600,
    BackgroundColour: application.NewRGB(255, 255, 255),
})
```

### app.Window.GetByName()

Obtém uma janela pelo nome. Retorna a janela e indica se ela foi encontrada.

```go
func (wm *WindowManager) GetByName(name string) (Window, bool)
```

**Exemplo:**

```go
if window, ok := app.Window.GetByName("main"); ok {
    window.Show()
}
```

### app.Window.GetAll()

Retorna todas as janelas da aplicação.

```go
func (wm *WindowManager) GetAll() []Window
```

**Exemplo:**

```go
windows := app.Window.GetAll()
for _, window := range windows {
    fmt.Println("Window:", window.Name())
}
```

## Gerenciadores

A aplicação fornece acesso a vários gerenciadores por meio de propriedades:

```go
app.Window       // Window management
app.Menu         // Menu management
app.Dialog       // Dialog management
app.Event        // Event management
app.Clipboard    // Clipboard operations
app.Screen       // Screen information
app.SystemTray   // System tray
app.Browser      // Browser operations
app.Env          // Environment variables
app.ContextMenu  // Context-menu management
app.KeyBinding   // Global keyboard shortcuts
app.Logger       // *slog.Logger
```

### Exemplo de uso

```go
// Create window
window := app.Window.New()

// Show dialog
app.Dialog.Info().SetMessage("Hello!").Show()

// Copy to clipboard
app.Clipboard.SetText("Copied text")

// Get screens
screens := app.Screen.GetAll()
```

## Gerenciamento de serviços

### RegisterService()

Registra um serviço na aplicação.

```go
func (a *App) RegisterService(service Service)
```

`RegisterService` não retorna nada; os erros de inicialização do serviço são manifestados por falhas de `ServiceStartup` durante `app.Run()`.

**Exemplo:**

```go
type MyService struct {
    app *application.App
}

func NewMyService(app *application.App) *MyService {
    return &MyService{app: app}
}

// Register after app creation
app.RegisterService(application.NewService(NewMyService(app)))
```

## Gerenciamento de eventos

### app.Event.Emit()

Emite um evento personalizado. Retorna `true` se um hook cancelar a emissão.

```go
func (em *EventManager) Emit(name string, data ...any) bool
```

**Exemplo:**

```go
// Emit event with data
app.Event.Emit("user-logged-in", map[string]interface{}{
    "username": "john",
    "timestamp": time.Now(),
})
```

### app.Event.On()

Monitora eventos personalizados. Retorna uma função `func()` para cancelar a inscrição.

```go
func (em *EventManager) On(name string, callback func(*CustomEvent)) func()
```

**Exemplo:**

```go
app.Event.On("user-logged-in", func(e *application.CustomEvent) {
    data := e.Data.(map[string]interface{})
    username := data["username"].(string)
    fmt.Println("User logged in:", username)
})
```

### app.Event.OnApplicationEvent()

Monitora eventos do ciclo de vida da aplicação. O parâmetro `eventType` é do tipo `events.ApplicationEventType` (do pacote `events`).

```go
func (em *EventManager) OnApplicationEvent(
    eventType events.ApplicationEventType,
    callback func(*ApplicationEvent),
) func()
```

**Exemplo:**

```go
import "github.com/wailsapp/wails/v3/pkg/events"

// Listen for app-started
app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(e *application.ApplicationEvent) {
    fmt.Println("Application started")
})

// Application shutdown is NOT an event constant; register cleanup via:
app.OnShutdown(func() {
    fmt.Println("Application shutting down")
})
```

## Métodos de diálogo

Os diálogos são acessados por meio do gerenciador `app.Dialog`. Consulte a [API de diálogos](/reference/dialogs/) para obter a referência completa.

### Diálogos de mensagem

```go
// Information dialog
app.Dialog.Info().
    SetTitle("Success").
    SetMessage("Operation completed!").
    Show()

// Error dialog
app.Dialog.Error().
    SetTitle("Error").
    SetMessage("Something went wrong.").
    Show()

// Warning dialog
app.Dialog.Warning().
    SetTitle("Warning").
    SetMessage("This action cannot be undone.").
    Show()
```

### Diálogos de pergunta

Os diálogos de pergunta usam callbacks de botões para processar as respostas do usuário:

```go
dialog := app.Dialog.Question().
    SetTitle("Confirm").
    SetMessage("Continue?")

yes := dialog.AddButton("Yes")
yes.OnClick(func() {
    // Handle yes
})

no := dialog.AddButton("No")
no.OnClick(func() {
    // Handle no
})

dialog.SetDefaultButton(yes)
dialog.SetCancelButton(no)
dialog.Show()
```

### Diálogos de arquivo

```go
// Open file dialog
path, err := app.Dialog.OpenFile().
    SetTitle("Select File").
    AddFilter("Images", "*.png;*.jpg").
    PromptForSingleSelection()

// Save file dialog
path, err := app.Dialog.SaveFile().
    SetTitle("Save File").
    SetFilename("document.pdf").
    AddFilter("PDF", "*.pdf").
    PromptForSingleSelection()

// Folder selection (use OpenFile with directory options)
path, err := app.Dialog.OpenFile().
    SetTitle("Select Folder").
    CanChooseDirectories(true).
    CanChooseFiles(false).
    PromptForSingleSelection()
```

## Logger

A aplicação fornece um logger estruturado:

```go
app.Logger.Info("Message", "key", "value")
app.Logger.Error("Error occurred", "error", err)
app.Logger.Debug("Debug info")
app.Logger.Warn("Warning message")
```

**Exemplo:**

```go
func (s *MyService) ProcessData(data string) error {
    s.app.Logger.Info("Processing data", "length", len(data))
    
    if err := process(data); err != nil {
        s.app.Logger.Error("Processing failed", "error", err)
        return err
    }
    
    s.app.Logger.Info("Processing complete")
    return nil
}
```

## Processamento de mensagens brutas

Para aplicações que precisam de controle direto e de baixo nível sobre a comunicação entre o frontend e o backend, o Wails fornece a opção `RawMessageHandler`. Essa opção ignora o sistema de bindings padrão.

@note{type="info"}
Mensagens brutas só devem ser usadas como último recurso. O sistema de bindings padrão é altamente otimizado e suficiente para quase todas as aplicações. Use mensagens brutas apenas se você tiver analisado o desempenho da aplicação com um profiler e confirmado que os bindings são um gargalo.

@end

### RawMessageHandler

`RawMessageHandler` é um campo de `application.Options`, não um método. O runtime o invoca para cada mensagem bruta enviada pelo frontend por meio de `System.invoke()`.

```go
type Options struct {
    // ... other fields ...
    RawMessageHandler func(window Window, message string, originInfo *OriginInfo)
}
```

`OriginInfo` contém `Origin`, `TopOrigin` e `IsMainFrame` (as plataformas preenchem subconjuntos diferentes — consulte o [Guia de mensagens brutas](/guides/raw-messages/) para ver a matriz de cada plataforma).

**Exemplo:**

```go
app := application.New(application.Options{
    Name: "My App",
    RawMessageHandler: func(window application.Window, message string, originInfo *application.OriginInfo) {
        // Handle the raw message
        fmt.Printf("Received from %s (%s): %s\n", window.Name(), originInfo.Origin, message)

        // You can respond using events
        window.EmitEvent("response", processMessage(message))
    },
})
```

Para obter mais detalhes, consulte o [Guia de mensagens brutas](/guides/raw-messages/).

## Opções específicas da plataforma

### Opções do Windows

Configure o comportamento específico do Windows no nível da aplicação:

```go
app := application.New(application.Options{
    Name: "My App",
    Windows: application.WindowsOptions{
        // WebView2 browser flags (apply to ALL windows)
        EnabledFeatures:       []string{"msWebView2EnableDraggableRegions"},
        DisabledFeatures:      []string{"msExperimentalFeature"},
        AdditionalBrowserArgs: []string{"--remote-debugging-port=9222"},

        // Other Windows options
        WndClass:                      "MyAppClass",
        WebviewUserDataPath:           "",  // Default: %APPDATA%\[BinaryName.exe]
        WebviewBrowserPath:            "",  // Default: system WebView2
        DisableQuitOnLastWindowClosed: false,
    },
})
```

**Flags do navegador:**

- `EnabledFeatures` - flags de recursos do WebView2 a serem ativadas
- `DisabledFeatures` - flags de recursos do WebView2 a serem desativadas
- `AdditionalBrowserArgs` - argumentos de linha de comando do Chromium

Consulte [Opções de janela — Opções do Windows no nível do aplicativo](/features/windows/options/#application-level-windows-options) para obter a documentação detalhada.

### Opções do Mac

```go
app := application.New(application.Options{
    Name: "My App",
    Mac: application.MacOptions{
        ActivationPolicy: application.ActivationPolicyRegular,
        ApplicationShouldTerminateAfterLastWindowClosed: true,
    },
})
```

### Opções do Linux

```go
app := application.New(application.Options{
    Name: "My App",
    Linux: application.LinuxOptions{
        ProgramName:                   "my-app",
        DisableQuitOnLastWindowClosed: false,
    },
})
```

## Exemplo completo de aplicativo

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name:        "My Application",
        Description: "A demo application",
        Mac: application.MacOptions{
            ApplicationShouldTerminateAfterLastWindowClosed: true,
        },
    })

    // Create main window
    window := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:            "My App",
        Width:            1024,
        Height:           768,
        MinWidth:         800,
        MinHeight:        600,
        BackgroundColour: application.NewRGB(255, 255, 255),
        URL:              "http://wails.localhost/",
    })

    window.Center()
    window.Show()

    app.Run()
}
```
