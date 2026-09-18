---
title: "Opções de janela"
description: "Referência completa de WebviewWindowOptions"
slug: "features/windows/options"
sourcePath: "features/windows/options.md"
---

## Opções de configuração da janela

O Wails oferece uma configuração abrangente de janelas, com dezenas de opções de tamanho, posição, aparência e comportamento. Esta é a **referência completa** de todas as opções disponíveis para Windows, macOS e Linux referentes a `WebviewWindowOptions`. Todas as opções, todas as plataformas, com exemplos e restrições.

## Estrutura de WebviewWindowOptions

```go
type WebviewWindowOptions struct {
    // Identity
    Name  string
    Title string

    // Size and Position
    Width           int
    Height          int
    X               int
    Y               int
    MinWidth        int
    MinHeight       int
    MaxWidth        int
    MaxHeight       int
    InitialPosition WindowStartPosition // WindowCentered (default) or WindowXY
    Screen          *Screen             // target screen for initial placement

    // Initial State
    Hidden        bool
    Frameless     bool
    DisableResize bool        // inverted vs v2's `Resizable`
    AlwaysOnTop   bool
    StartState    WindowState // WindowStateNormal | Minimised | Maximised | Fullscreen

    // Appearance
    BackgroundColour RGBA
    BackgroundType   BackgroundType
    Zoom             float64
    ZoomControlEnabled bool

    // Content
    URL  string
    HTML string
    JS   string
    CSS  string

    // Behaviour
    EnableFileDrop              bool
    IgnoreMouseEvents           bool
    HideOnFocusLost             bool
    HideOnEscape                bool
    DevToolsEnabled             bool
    DefaultContextMenuDisabled  bool
    ContentProtectionEnabled    bool
    KeyBindings                 map[string]func(window *WebviewWindow)

    // Permissions
    Permissions map[PermissionType]Permission

    // Window-control button states
    MinimiseButtonState ButtonState
    MaximiseButtonState ButtonState
    CloseButtonState    ButtonState

    // Menu
    UseApplicationMenu bool

    // Platform-specific (per-window)
    Mac     MacWindow
    Windows WindowsWindow
    Linux   LinuxWindow
}
```

`WebviewWindowOptions` **não** tem um campo `Parent` — para relações de janela pai/modal, use `parentWindow.AttachModal(childWindow)`. Também **não** tem um campo `Assets` — a configuração de ativos fica em `application.Options` (`Assets AssetOptions`).

Código-fonte completo: [`v3/pkg/application/webview_window_options.go`](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window_options.go).

## Opções principais

### Nome

**Tipo:** `string` **Padrão:** UUID gerado automaticamente **Plataforma:** Todas

```go
Name: "main-window"
```

**Finalidade:** Identificador exclusivo para localizar janelas posteriormente.

**Práticas recomendadas:**

- Use nomes descritivos: `"main"`, `"settings"`, `"about"`
- Use kebab-case: `"file-browser"`, `"color-picker"`
- Use um nome curto e fácil de lembrar

**Exemplo:**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name: "settings-window",
})

// Later...
if settings, ok := app.Window.GetByName("settings-window"); ok {
    settings.Focus()
}
```

### Título

**Tipo:** `string` **Padrão:** Nome do aplicativo **Plataforma:** Todas

```go
Title: "My Application"
```

**Finalidade:** Texto exibido na barra de título e na barra de tarefas.

**Atualizações dinâmicas:**

```go
window.SetTitle("My Application - Document.txt")
```

### Largura / Altura

**Tipo:** `int` (pixels) **Padrão:** 800 x 600 **Plataforma:** Todas **Restrições:** Deve ser positivo

```go
Width:  1200,
Height: 800,
```

**Finalidade:** Tamanho inicial da janela em pixels lógicos.

**Observações:**

- O Wails processa automaticamente o dimensionamento de DPI
- Use pixels lógicos, não pixels físicos
- Considere a resolução mínima da tela (1024x768)

**Exemplos de tamanho:**

| Caso de uso | Largura | Altura |
| --- | --- | --- |
| Utilitário pequeno | 400 | 300 |
| Aplicativo padrão | 1024 | 768 |
| Aplicativo grande | 1440 | 900 |
| Full HD | 1920 | 1080 |

### X / Y

**Tipo:** `int` (pixels) **Padrão:** Centralizada na tela **Plataforma:** Todas

```go
X: 100,  // 100px from left edge
Y: 100,  // 100px from top edge
```

**Finalidade:** Posição inicial da janela.

**Sistema de coordenadas:**

- (0, 0) corresponde ao canto superior esquerdo da tela principal
- Os valores positivos de X avançam para a direita
- Os valores positivos de Y avançam para baixo

**Exemplo:**

`X` e `Y` só entram em vigor quando `InitialPosition: application.WindowXY` está definido. Caso contrário, `InitialPosition` usa `WindowCentered` como padrão, e `X`/`Y` são ignorados.

```go
settings := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:            "coordinate-window",
    InitialPosition: application.WindowXY, // opt into X/Y coordinates
    X:               100,
    Y:               100,
})
```

**Prática recomendada:** Se você não precisar de coordenadas específicas, use `Center()` para centralizar uma janela após criá-la:

```go
window := app.Window.New()
window.Center()
```

### MinWidth / MinHeight

**Tipo:** `int` (pixels) **Padrão:** 0 (sem mínimo) **Plataforma:** Todas

```go
MinWidth:  400,
MinHeight: 300,
```

**Finalidade:** Impedir que a janela fique pequena demais.

**Casos de uso:**

- Evitar layouts quebrados
- Garantir a usabilidade
- Manter a proporção

**Exemplo:**

```go
// Prevent window smaller than 400x300
MinWidth:  400,
MinHeight: 300,
```

### MaxWidth / MaxHeight

**Tipo:** `int` (pixels) **Padrão:** 0 (sem máximo) **Plataforma:** Todas

```go
MaxWidth:  1920,
MaxHeight: 1080,
```

**Finalidade:** Impedir que a janela fique grande demais.

**Casos de uso:**

- Aplicativos de tamanho fixo
- Evitar o uso excessivo de recursos
- Manter as restrições de design

## Opções de estado

### Hidden

**Tipo:** `bool` **Padrão:** `false` **Plataforma:** Todas

```go
Hidden: true,
```

**Finalidade:** Criar a janela sem exibi-la.

**Casos de uso:**

- Janelas em segundo plano
- Janelas exibidas sob demanda
- Telas de abertura (criar, carregar e depois exibir)
- Evitar um clarão branco durante o carregamento do conteúdo

**Melhorias específicas de plataforma:**

- **Windows:** O clarão branco da janela foi corrigido — a janela permanece invisível até que `Show()` seja chamado
- **macOS:** Compatibilidade total
- **Linux:** Compatibilidade total

**Padrão recomendado para um carregamento suave:**

```go
// Create hidden window
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:             "main-window",
    Hidden:           true,
    BackgroundColour: application.NewRGB(30, 30, 30), // Match your theme
})

// Load content while hidden
// ... content loads ...

// Show when ready (no flash!)
window.Show()
```

**Exemplo:**

```go
settings := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:   "settings",
    Hidden: true,
})

// Show when needed
settings.Show()
```

### Frameless

**Tipo:** `bool` **Padrão:** `false` **Plataforma:** Todas

```go
Frameless: true,
```

**Finalidade:** Remover a barra de título e as bordas da janela.

**Casos de uso:**

- Decoração de janela personalizada
- Telas de abertura
- Aplicativos de quiosque
- Janelas com design personalizado

**Importante:** Você precisará implementar:

- Arrastar a janela
- Botões de fechar, minimizar e maximizar
- Alças de redimensionamento (se a janela for redimensionável)

**Consulte [Janelas sem moldura](/features/windows/frameless/) para obter detalhes.**

### DisableResize

**Tipo:** `bool` **Padrão:** `false` (por padrão, a janela é redimensionável) **Plataforma:** Todas

```go
DisableResize: true,
```

**Finalidade:** Impedir o redimensionamento da janela. Observe que o campo é o **inverso** de `Resizable` na v2 — defina `DisableResize: true` para impedir que uma janela seja redimensionada.

**Casos de uso:**

- Aplicativos de tamanho fixo
- Telas de abertura
- Caixas de diálogo

**Observação:** Os usuários ainda podem maximizar a janela ou colocá-la em tela cheia, a menos que você também desative essas opções por meio de `MaximiseButtonState` ou do comportamento da coleção.

### AlwaysOnTop

**Tipo:** `bool` **Padrão:** `false` **Plataforma:** Todas

```go
AlwaysOnTop: true,
```

**Finalidade:** Manter a janela acima de todas as outras.

**Casos de uso:**

- Barras de ferramentas flutuantes
- Notificações
- Picture-in-picture
- Temporizadores

**Observações específicas de plataforma:**

- **macOS:** Compatibilidade total
- **Windows:** Compatibilidade total
- **Linux:** Depende do gerenciador de janelas

### StartState

**Tipo:** enum `WindowState` **Padrão:** `WindowStateNormal` **Plataforma:** Todas

```go
StartState: application.WindowStateMaximised,
```

**Finalidade:** Definir o estado inicial da janela quando ela for exibida.

**Valores:**

- `WindowStateNormal` — Janela normal
- `WindowStateMinimised` — Minimizada
- `WindowStateMaximised` — Maximizada
- `WindowStateFullscreen` — Tela cheia

Não há uma constante `WindowStateHidden` — use o campo booleano `Hidden` para iniciar a janela invisível.

**Alternar o modo de tela cheia em tempo de execução:**

```go
window.Fullscreen()
window.UnFullscreen()
window.ToggleFullscreen() // there is no SetFullscreen(bool)
```

## Opções de aparência

### BackgroundColour

**Tipo:** struct `RGBA` **Padrão:** Branco **Plataforma:** Todas

```go
BackgroundColour: application.RGBA{Red: 0, Green: 0, Blue: 0, Alpha: 255},
```

Os campos de `RGBA` são `Red, Green, Blue, Alpha` (uint8). Prefira as funções auxiliares `application.NewRGB(r, g, b)` (alfa 255) ou `application.NewRGBA(r, g, b, a)`.

**Finalidade:** Cor de fundo da janela antes do carregamento do conteúdo.

**Casos de uso:**

- Combinar com o tema do aplicativo
- Evitar um clarão branco em temas escuros
- Proporcionar uma experiência de carregamento suave

**Exemplo:**

```go
// Dark theme
BackgroundColour: application.NewRGB(30, 30, 30),

// Light theme
BackgroundColour: application.NewRGB(255, 255, 255),
```

**Método auxiliar:**

```go
window.SetBackgroundColour(application.NewRGB(30, 30, 30))
```

### BackgroundType

**Tipo:** enum `BackgroundType` **Padrão:** `BackgroundTypeSolid` **Plataforma:** macOS, Windows (parcial)

```go
BackgroundType: application.BackgroundTypeTranslucent,
```

**Valores:**

- `BackgroundTypeSolid` — Cor sólida
- `BackgroundTypeTransparent` — Totalmente transparente
- `BackgroundTypeTranslucent` — Desfoque semitransparente

**Compatibilidade com plataformas:**

- **macOS:** Configure `Mac.Backdrop`; a transparência da webview requer [`private_mac_apis`](/guides/build/private-macos-apis/#webview-transparency-and-background). Sem isso, a webview permanece opaca.
- **Windows:** Transparente e translúcido (Windows 11+)
- **Linux:** Somente sólido

**Exemplo (macOS):**

```go
BackgroundType: application.BackgroundTypeTranslucent,
Mac: application.MacWindow{
    Backdrop: application.MacBackdropTranslucent,
},
```

### OpenInspectorOnStartup e OpenDevTools

**API privada no macOS:** `OpenInspectorOnStartup: true`, `window.OpenDevTools()` do Go e `Window.OpenDevTools()` do JavaScript exigem `private_mac_apis` para abrir o inspetor por meio de programação. Sem isso, essas operações não têm efeito. As compilações de produção também precisam de `devtools`. A inspeção pública do Safari no macOS 13.3+ não precisa de APIs privadas; a ativação do inspetor em versões anteriores do macOS precisa. Consulte a [matriz de compilação do Inspetor Web](/guides/build/private-macos-apis/#web-inspector).

## Opções de conteúdo

### URL

**Tipo:** `string` **Padrão:** Vazio (carrega a partir de Assets) **Plataforma:** Todas

```go
URL: "https://example.com",
```

**Finalidade:** Carregar uma URL externa em vez dos recursos incorporados.

**Casos de uso:**

- Desenvolvimento (carregar a partir do servidor de desenvolvimento)
- Aplicativos baseados na Web
- Aplicativos híbridos

**Exemplo:**

```go
// Development — point the window at the Vite dev server
URL: "http://localhost:9245",

// Production — embedded assets are configured at the application level
// (Assets is application.Options.Assets, not a WebviewWindowOptions field).
```

Trecho no nível do aplicativo para o caso de produção:

```go
app := application.New(application.Options{
    Name: "My App",
    Assets: application.AssetOptions{
        Handler: application.AssetFileServerFS(assets),
    },
})
```

### HTML

**Tipo:** `string` **Padrão:** Vazio **Plataforma:** Todas

```go
HTML: "<h1>Hello World</h1>",
```

**Finalidade:** Carregar diretamente uma string HTML.

**Casos de uso:**

- Janelas simples
- Conteúdo gerado
- Testes

**Exemplo:**

```go
HTML: `
<!DOCTYPE html>
<html>
<head><title>Simple Window</title></head>
<body><h1>Hello from Wails!</h1></body>
</html>
`,
```

### Assets (somente no nível do aplicativo)

A configuração de recursos **não** é um campo de `WebviewWindowOptions`. Os recursos do frontend são servidos pelo próprio aplicativo por meio de `application.Options.Assets` (`AssetOptions`); todas as janelas herdam esse servidor de recursos.

```go
app := application.New(application.Options{
    Assets: application.AssetOptions{
        Handler: application.AssetFileServerFS(assets),
    },
})
```

**Consulte [Sistema de compilação](/concepts/build-system/) para obter detalhes.**

### UseApplicationMenu

**Tipo:** `bool` **Padrão:** `false` **Plataforma:** Windows, Linux (sem efeito no macOS)

```go
UseApplicationMenu: true,
```

**Finalidade:** Usar o menu do aplicativo (definido por meio de `app.Menu.Set()`) nesta janela.

No **macOS**, esta opção não tem efeito porque o macOS sempre usa um menu global do aplicativo na parte superior da tela.

No **Windows** e no **Linux**, as janelas não exibem um menu por padrão. Definir `UseApplicationMenu: true` instrui a janela a usar o menu no nível do aplicativo, oferecendo uma solução multiplataforma simples.

**Exemplo:**

```go
// Set the application menu once
menu := app.NewMenu()
menu.AddRole(application.FileMenu)
menu.AddRole(application.EditMenu)
app.Menu.Set(menu)

// All windows with UseApplicationMenu will display this menu
app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:              "Main Window",
    UseApplicationMenu: true,
})

app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:              "Second Window",
    UseApplicationMenu: true,  // Also gets the app menu
})
```

**Observações:**

- Se `UseApplicationMenu` e um menu específico da janela estiverem definidos, o menu específico da janela terá prioridade
- Isso simplifica o código multiplataforma, eliminando a necessidade de verificar o sistema operacional em tempo de execução
- Consulte [Menus do aplicativo](/features/menus/application/) para obter a documentação completa sobre menus

## Opções de entrada

### EnableFileDrop

**Tipo:** `bool` **Padrão:** `false` **Plataforma:** Todas

```go
EnableFileDrop: true,
```

**Finalidade:** Permitir arrastar e soltar arquivos do sistema operacional na janela.

Quando ativado:

- Os arquivos arrastados de gerenciadores de arquivos podem ser soltos no aplicativo
- O evento `WindowFilesDropped` é disparado com os caminhos dos arquivos soltos
- Os elementos com o atributo `data-file-drop-target` fornecem informações detalhadas sobre a operação de soltar arquivos

**Casos de uso:**

- Interfaces de upload de arquivos
- Editores de documentos
- Importadores de mídia
- Qualquer aplicativo que aceite arquivos

**Exemplo:**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:          "File Uploader",
    EnableFileDrop: true,
})

// Handle dropped files
window.OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {
    files := event.Context().DroppedFiles()
    details := event.Context().DropTargetDetails()
    
    for _, file := range files {
        fmt.Println("Dropped:", file)
    }
})
```

**Zonas de soltura em HTML:**

```html
<!-- Mark elements as drop targets -->
<div id="upload" data-file-drop-target>
    Drop files here
</div>
```

**Consulte [Soltura de arquivos](/features/drag-and-drop/files/) para obter a documentação completa.**

## Opções de segurança

### ContentProtectionEnabled

**Tipo:** `bool` **Padrão:** `false` **Plataforma:** Windows (10+), macOS

```go
ContentProtectionEnabled: true,
```

**Finalidade:** Impedir a captura do conteúdo da janela.

**Compatibilidade com plataformas:**

- **Windows:** Windows 10, compilação 19041+ (compatibilidade total); versões anteriores (compatibilidade parcial)
- **macOS:** Compatibilidade total
- **Linux:** Sem compatibilidade

**Casos de uso:**

- Aplicativos bancários
- Gerenciadores de senhas
- Prontuários médicos
- Documentos confidenciais

**Observações importantes:**

1. Não impede fotografias físicas da tela
2. Algumas ferramentas podem contornar a proteção
3. Faz parte de uma estratégia de segurança abrangente; não deve ser a única proteção
4. As janelas do DevTools não são protegidas automaticamente

**Exemplo:**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Secure Window",
    ContentProtectionEnabled: true,
})

// Toggle at runtime
window.SetContentProtection(true)
```

### Permissões

**Tipo:** `map[PermissionType]Permission` **Padrão:** `nil` (tratamento padrão da plataforma) **Plataforma:** Linux, Windows (o macOS delega ao TCC)

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionAllow,
    application.PermissionCamera:     application.PermissionAllow,
},
```

**Finalidade:** Controlar declarativamente, sem código específico da plataforma, como são tratadas as solicitações de recursos (câmera, microfone, geolocalização, notificações e leitura da área de transferência) provenientes do conteúdo web da janela.

**Valores de PermissionType:** `PermissionMicrophone`, `PermissionCamera`, `PermissionGeolocation`, `PermissionNotifications`, `PermissionClipboardRead`

**Valores de Permission:**

- `PermissionDefault` (0) — tratamento nativo da plataforma: solicitação do sistema operacional/WebView2 no macOS/Windows; no Linux, o acesso à câmera e ao microfone é permitido, e todos os demais são negados
- `PermissionAllow` (1) — conceder sem solicitar confirmação (Linux: apenas câmera e microfone estão implementados; os demais tipos continuam negados)
- `PermissionDeny` (2) — negar sem solicitar confirmação

**Importante — Windows:** Antes da existência desta opção, o Wails concedia silenciosamente todos os recursos do WebView2. Agora, definir qualquer entrada em `Permissions` desativa essa concessão geral. Os recursos não listados exibirão a solicitação nativa do WebView2, em vez de serem permitidos automaticamente. Liste explicitamente todos os recursos necessários ao aplicativo.

**Consulte [Permissões](/features/windows/permissions/) para acessar o guia completo, a matriz de plataformas e os exemplos.**

## Eventos do ciclo de vida da janela

Os eventos do ciclo de vida da janela são tratados usando `OnWindowEvent` e `RegisterHook`. Esses métodos permitem controlar detalhadamente o comportamento de fechamento e destruição da janela.

### Cancelamento do fechamento da janela

Para impedir que uma janela seja fechada (por exemplo, quando houver alterações não salvas), use `RegisterHook` com o evento `WindowClosing` e chame `event.Cancel()`:

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:  "main-window",
    Title: "My Application",
})

// Register a hook to intercept the closing event
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

**Pontos principais:**

- `RegisterHook` intercepta eventos antes que ocorram
- Chame `event.Cancel()` para impedir que a janela seja fechada
- A janela permanecerá aberta após o cancelamento

### Tratamento do fechamento da janela

Para realizar a limpeza quando uma janela for fechada, use `OnWindowEvent` com o evento `WindowClosing`:

```go
window.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
    // Cleanup code runs here
    fmt.Printf("Window %s is closing\n", window.Name())

    // Close database connection
    if db != nil {
        db.Close()
    }

    // Remove from window list
    removeWindow(window.ID())
})
```

**Pontos principais:**

- `OnWindowEvent` trata eventos que estão prestes a ocorrer
- A limpeza é executada antes que a janela seja destruída
- Não é possível cancelar o fechamento aqui (para isso, use `RegisterHook`)

### Padrão de limpeza de janela singleton

Para janelas singleton (para garantir apenas uma instância), use `WindowClosing` para limpar a referência:

```go
var settingsWindow *application.WebviewWindow

func ShowSettings(app *application.App) {
    // Create if doesn't exist
    if settingsWindow == nil {
        settingsWindow = app.Window.NewWithOptions(application.WebviewWindowOptions{
            Name:  "settings",
            Title: "Settings",
            Width: 600,
            Height: 400,
        })

        // Cleanup on close
        settingsWindow.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
            settingsWindow = nil
        })
    }

    // Show and focus
    settingsWindow.Show()
    settingsWindow.Focus()
}
```

## Opções específicas da plataforma

### Opções do Mac

```go
Mac: application.MacWindow{
    TitleBar: application.MacTitleBar{
        AppearsTransparent: true,
        Hide:               false,
        HideTitle:          true,
        FullSizeContent:    true,
    },
    Backdrop:                application.MacBackdropTranslucent,
    InvisibleTitleBarHeight: 50,
    WindowClass:             application.MacWindowClassPanel,
    PanelPreferences: application.MacPanelPreferences{
        NonActivating:          true,
        FloatingPanel:          true,
        BecomesKeyOnlyIfNeeded: false,
        UtilityWindow:          false,
    },
    WindowLevel:             application.MacWindowLevelFloating,
    CollectionBehavior:      application.MacWindowCollectionBehaviorDefault,
    TabbingMode:             application.MacWindowTabbingModeDisallowed,
},
```

**TitleBar** (`MacTitleBar`)

- `AppearsTransparent` — Torna a barra de título transparente e estende o conteúdo até a área da barra de título
- `Hide` — Oculta completamente a barra de título
- `HideTitle` — Oculta apenas o texto do título
- `FullSizeContent` — Estende o conteúdo por toda a janela

No macOS, a transparência da webview e a abertura programática do inspetor exigem a tag de compilação `private_mac_apis`. Sem ela, as mesmas opções continuam válidas, mas as operações que dependem de APIs privadas não fazem nada. O agrupamento Liquid Glass também é ignorado, e os estilos usam alternativas públicas. Consulte [APIs privadas do macOS](/guides/build/private-macos-apis/) para ver os comandos de compilação e o comportamento exato.

**Backdrop** (`MacBackdrop`)

- `MacBackdropNormal` — Fundo opaco padrão
- `MacBackdropTranslucent` — **É necessária uma API privada para a transparência da webview.** Sem a tag, o desfoque nativo permanece atrás de uma webview opaca.
- `MacBackdropTransparent` — **É necessária uma API privada para a transparência da webview.** Sem a tag, a webview permanece opaca.
- `MacBackdropLiquidGlass` — **É necessária uma API privada para a transparência da webview.** Sem a tag, a camada de vidro permanece atrás de uma webview opaca; os estilos usam alternativas públicas.

**LiquidGlass** (`MacLiquidGlass`)

| Campo ou valor | Dependência de API privada no macOS |
| --- | --- |
| `Style: LiquidGlassStyleAutomatic` | O estilo regular nativo é público; a transparência da webview sobre o plano de fundo exige `private_mac_apis`. |
| `Style: LiquidGlassStyleLight` | A tag preserva o mapeamento existente para o estilo nativo transparente; sem ela, o Wails usa vidro regular com a aparência Aqua. |
| `Style: LiquidGlassStyleDark` | **API privada:** valor de estilo nativo não documentado `2`; sem a tag, usa vidro regular com a aparência Dark Aqua. |
| `Style: LiquidGlassStyleVibrant` | O mapeamento para o estilo nativo transparente é público; a transparência da webview sobre o plano de fundo exige a tag. |
| `GroupID` | **API privada:** valores não vazios solicitam agrupamento; sem a tag, são ignorados. |
| `GroupSpacing` | **API privada:** valores positivos solicitam espaçamento entre grupos; sem a tag, são ignorados. |
| `Material`, `CornerRadius`, `TintColor` | Não possuem dependência própria de API privada. |

Consulte [Valores de Liquid Glass](/guides/build/private-macos-apis/#liquid-glass-values) para ver os valores de estilo nativo e a disponibilidade nos sistemas operacionais.

**InvisibleTitleBarHeight** (`int`)

- Altura da área invisível da barra de título (para arrastar)
- Só tem efeito quando a área nativa de arrastar da barra de título está oculta, ou seja, quando a janela não tem moldura (`Frameless: true`) ou usa uma barra de título transparente (`AppearsTransparent: true`)
- Não tem efeito em janelas padrão com barra de título visível

**WindowClass** (`MacWindowClass`)

- `MacWindowClassWindow` — Comportamento padrão de `NSWindow` (padrão)
- `MacWindowClassPanel` — Uma `NSPanel` auxiliar que nunca se torna a janela principal do aplicativo

`PanelPreferences` aplica-se apenas a `MacWindowClassPanel`:

- `NonActivating` adiciona `NSWindowStyleMaskNonactivatingPanel`. Exibir ou colocar o painel em foco não ativa o aplicativo Wails, mas o painel ainda pode se tornar a janela que recebe a entrada do teclado para controles e entrada de texto.
- `FloatingPanel` habilita o comportamento de painel flutuante do AppKit.
- `BecomesKeyOnlyIfNeeded` assume o foco de entrada do teclado somente quando a visualização clicada solicita entrada pelo teclado.
- `UtilityWindow` aplica o estilo nativo de janela de utilitário.

Os painéis do Wails permanecem visíveis quando o aplicativo é desativado e são liberados ao serem fechados, de acordo com o ciclo de vida esperado por `WebviewWindow`. Essas são sobrescritas intencionais dos padrões opostos de `NSPanel`.

A classe, o nível, a política de ativação e o comportamento de coleção da janela resolvem problemas distintos:

- `WindowClass` seleciona `NSWindow` ou `NSPanel` e controla a semântica de janela principal e janela-chave.
- `WindowLevel` controla a ordem Z. Use `MacWindowLevelPopUpMenu` para sobreposições na barra de menus.
- `MacOptions.ActivationPolicy` controla o aplicativo como um todo, incluindo sua apresentação no Dock e na barra de menus. Um painel que não ativa o aplicativo não exige uma política de ativação de acessório, embora um aplicativo ainda possa usá-la para ocultar seu ícone no Dock.
- `CollectionBehavior` controla a participação nos Spaces e no modo de tela cheia.

```go
// Spotlight/menu-bar panel that leaves the current application active.
Mac: application.MacWindow{
    WindowClass: application.MacWindowClassPanel,
    PanelPreferences: application.MacPanelPreferences{
        NonActivating: true,
    },
    WindowLevel: application.MacWindowLevelPopUpMenu,
    CollectionBehavior: application.MacWindowCollectionBehaviorCanJoinAllSpaces |
        application.MacWindowCollectionBehaviorFullScreenAuxiliary |
        application.MacWindowCollectionBehaviorStationary,
},
```

**WindowLevel** (`MacWindowLevel`)

- `MacWindowLevelNormal` - Nível de janela padrão (padrão)
- `MacWindowLevelFloating` - Flutua acima das janelas normais
- `MacWindowLevelTornOffMenu` - Nível de menu destacado
- `MacWindowLevelModalPanel` - Nível de painel modal
- `MacWindowLevelMainMenu` - Nível do menu principal
- `MacWindowLevelStatus` - Nível de janela de status
- `MacWindowLevelPopUpMenu` - Nível de menu pop-up
- `MacWindowLevelScreenSaver` - Nível do protetor de tela

Um `WindowLevel` explícito tem precedência sobre `AlwaysOnTop` e `PanelPreferences.FloatingPanel`. Sem um nível explícito, `AlwaysOnTop` ou um painel flutuante resulta em `MacWindowLevelFloating`; caso contrário, o nível é `MacWindowLevelNormal`. Uma chamada posterior a `SetAlwaysOnTop` continua sendo uma alteração explícita em tempo de execução.

**CollectionBehavior** (`MacWindowCollectionBehavior`)

Controla como a janela se comporta nos Spaces do macOS e no modo de tela cheia. Esses são valores de máscara de bits que podem ser combinados usando a operação OR bit a bit (`|`).

**Comportamento nos Spaces:**

- `MacWindowCollectionBehaviorDefault` - Usa FullScreenPrimary (padrão, compatível com versões anteriores)
- `MacWindowCollectionBehaviorCanJoinAllSpaces` - A janela aparece em todos os Spaces
- `MacWindowCollectionBehaviorMoveToActiveSpace` - Move-se para o Space ativo quando exibida
- `MacWindowCollectionBehaviorManaged` - Comportamento padrão de janela gerenciada
- `MacWindowCollectionBehaviorTransient` - Janela temporária/transitória
- `MacWindowCollectionBehaviorStationary` - Permanece imóvel durante a alternância entre Spaces

**Alternância cíclica entre janelas:**

- `MacWindowCollectionBehaviorParticipatesInCycle` - Incluída na alternância cíclica com Cmd+`
- `MacWindowCollectionBehaviorIgnoresCycle` - Excluída da alternância cíclica com Cmd+`

**Comportamento em tela cheia:**

- `MacWindowCollectionBehaviorFullScreenPrimary` - Pode entrar no modo de tela cheia
- `MacWindowCollectionBehaviorFullScreenAuxiliary` - Pode ser sobreposta a aplicativos em tela cheia
- `MacWindowCollectionBehaviorFullScreenNone` - Desabilita o recurso de tela cheia
- `MacWindowCollectionBehaviorFullScreenAllowsTiling` - Permite o posicionamento lado a lado (macOS 10.11+)
- `MacWindowCollectionBehaviorFullScreenDisallowsTiling` - Impede o posicionamento lado a lado (macOS 10.11+)

**Exemplo — janela semelhante ao Spotlight:**

```go
// Window that appears on all Spaces AND can overlay fullscreen apps
Mac: application.MacWindow{
	WindowClass: application.MacWindowClassPanel,
	PanelPreferences: application.MacPanelPreferences{
		NonActivating: true,
	},
    CollectionBehavior: application.MacWindowCollectionBehaviorCanJoinAllSpaces |
                        application.MacWindowCollectionBehaviorFullScreenAuxiliary,
    WindowLevel:        application.MacWindowLevelFloating,
},
```

**Exemplo — comportamento único:**

```go
// Window that can appear over fullscreen applications
Mac: application.MacWindow{
    CollectionBehavior: application.MacWindowCollectionBehaviorFullScreenAuxiliary,
},
```

**TabbingMode** (`MacWindowTabbingMode`)

Controla o comportamento das abas de janela no macOS 10.12 e posteriores. As abas permitem agrupar várias janelas como abas.

**Opções:**

- `MacWindowTabbingModeDefault` - Valor sentinela zero (não definido explicitamente). Em tempo de execução, o padrão é não permitir abas
- `MacWindowTabbingModeAutomatic` - O sistema determina o comportamento das abas
- `MacWindowTabbingModePreferred` - A janela prefere permanecer no modo de abas
- `MacWindowTabbingModeDisallowed` - Desabilita as abas de janela

**Exemplo — desabilitar abas de janela:**

```go
Mac: application.MacWindow{
    TabbingMode: application.MacWindowTabbingModeDisallowed,
},
```

**Exemplo — dar preferência a abas de janela:**

```go
Mac: application.MacWindow{
    TabbingMode: application.MacWindowTabbingModePreferred,
},
```

**WebviewPreferences** (`MacWebviewPreferences`)

Controle detalhado sobre a configuração subjacente de `WKWebView`. Todos os campos são opcionais — os campos não definidos mantêm inalterado o padrão do WebKit.

```go
Mac: application.MacWindow{
    WebviewPreferences: application.MacWebviewPreferences{
        TabFocusesLinks:                       optional.True,
        TextInteractionEnabled:                optional.True,
        FullscreenEnabled:                     optional.False,
        AllowsBackForwardNavigationGestures:   optional.False,
        AllowsMagnification:                   optional.True,
        AllowsAirPlayForMediaPlayback:         optional.True,
        JavaScriptCanOpenWindowsAutomatically: optional.False,
        MinimumFontSize:                       optional.NewVar(12.0),
        ApplicationNameForUserAgent:           "MyApp",
        EnableAutoplayWithoutUserAction:       optional.True,
    },
},
```

- `TabFocusesLinks` — Quando `true`, pressionar Tab move o foco para links e controles de formulário (padrão: `false`)
- `TextInteractionEnabled` — Quando `true`, os usuários podem selecionar e interagir com o texto na webview (padrão: `true`)
- `FullscreenEnabled` — Quando `true`, o conteúdo da Web pode entrar em tela cheia por meio da API Fullscreen do HTML (padrão: `false`). Requer macOS 12.3 ou posterior.
- `AllowsBackForwardNavigationGestures` — Quando `true`, gestos de deslizamento horizontal acionam a navegação para trás/para a frente (padrão: `false`)
- `AllowsMagnification` — Quando `true`, o zoom por gesto de pinça é habilitado na webview (padrão: `false`)
- `AllowsAirPlayForMediaPlayback` — Quando `true`, é possível transmitir mídia para dispositivos AirPlay (padrão: `true`)
- `JavaScriptCanOpenWindowsAutomatically` — Quando `true`, o JavaScript pode abrir novas janelas sem um gesto do usuário (padrão: `false`)
- `MinimumFontSize` — Tamanho mínimo da fonte em pontos. Use `optional.NewVar(12.0)` para defini-lo. Se não for definido, o padrão do WebKit será mantido.
- `ApplicationNameForUserAgent` — Substitui o sufixo do nome do aplicativo na string de agente do usuário do WebKit. Útil quando sites rejeitam o identificador `"wails.io"` padrão (por exemplo, incorporações do YouTube). Deixe vazio para manter o padrão.
- `EnableAutoplayWithoutUserAction` — Quando `true`, áudio e vídeo podem ser reproduzidos automaticamente sem um gesto do usuário. Corresponde a `WKWebViewConfiguration.mediaTypesRequiringUserActionForPlayback = WKAudiovisualMediaTypeNone` (padrão: `false`)

### Opções do Windows (por janela)

A struct por janela é `application.WindowsWindow` — **não** `WindowsOptions` (essa é a struct no nível do *aplicativo*).

```go
Windows: application.WindowsWindow{
    DisableIcon:                       false,
    DisableMenu:                       false,
    BackdropType:                      application.Auto,
    CustomTheme:                       application.ThemeSettings{},
    DisableFramelessWindowDecorations: false,
    NonClientRegionSupport:            false,
    WebView2CompositionHosting:        false,
},
```

**DisableIcon** (`bool`)

- Remove o ícone da barra de título.

**DisableMenu** (`bool`)

- Desabilita a barra de menus da janela. Quando `true`, a janela não exibe uma barra de menus, mesmo que uma esteja configurada.
- Padrão: `false`

**BackdropType** (`BackdropType`)

- `application.Auto` - Padrão do sistema
- `application.None` - Sem plano de fundo
- `application.Mica` - Material Mica (Windows 11)
- `application.Acrylic` - Material Acrylic (Windows 11)
- `application.Tabbed` - Material Tabbed (Windows 11)

Não há constantes no estilo `WindowsBackdropTypeMica` — use `application.Mica` etc.

**CustomTheme** (`ThemeSettings`)

- Valor (não um ponteiro). Cores personalizadas para os modos escuro/claro da borda da janela, do texto/plano de fundo da barra de título e da barra de menus.

**DisableFramelessWindowDecorations** (`bool`)

- Desabilita as decorações padrão de janelas sem moldura (sombra Aero e cantos arredondados).

**NonClientRegionSupport** (`bool`)

- Habilita o suporte nativo do WebView2 a `app-region: drag` / `app-region: no-drag` para barras de título personalizadas sem moldura.
- Destina-se somente ao arraste nativo simples do aplicativo. Não oferece comportamento nativo para botões de legenda personalizados nem o Snap Assist / Snap Layouts do Windows 11 para botões de maximização personalizados.

**WebView2CompositionHosting** (`bool`)

- Habilita o suporte gerenciado pelo Wails a `--wails-non-client-region` para botões de legenda personalizados com comportamento nativo do Windows, incluindo o Snap Assist / Snap Layouts do Windows 11 em botões de maximização personalizados.
- Experimental. Hospeda o WebView2 por meio de `ICoreWebView2CompositionController` e DirectComposition, em vez do controlador padrão hospedado em HWND.
- Pode ser combinado com `NonClientRegionSupport` quando uma janela precisar tanto do suporte nativo do WebView2 a `app-region` quanto de regiões de botões de legenda personalizados gerenciadas pelo Wails.

**Exemplo:**

```go
Windows: application.WindowsWindow{
    BackdropType: application.Mica,
    DisableIcon:  true,
},
```

**Exemplo — Regiões personalizadas da barra de título do Windows:**

```go
Windows: application.WindowsWindow{
    NonClientRegionSupport:    true,
    WebView2CompositionHosting: true,
},
```

Consulte [Janelas sem moldura](/features/windows/frameless/#native-non-client-regions-on-windows) para ver o comportamento detalhado, as vantagens e desvantagens e o CSS correspondente.

### Opções do Linux (por janela)

A struct por janela é `application.LinuxWindow` — **não** `LinuxOptions`.

```go
Linux: application.LinuxWindow{
    Icon:                []byte{/* PNG data */},
    WindowIsTranslucent: false,
},
```

**Icon** (`[]byte`)

- Ícone da janela (formato PNG).

**WindowIsTranslucent** (`bool`)

- Requer suporte do compositor.

**Exemplo:**

```go
//go:embed icon.png
var icon []byte

Linux: application.LinuxWindow{
    Icon: icon,
},
```

## Opções do Windows no nível do aplicativo

Algumas opções específicas do Windows devem ser configuradas no nível do aplicativo, e não por janela. Isso ocorre porque o WebView2 compartilha um único ambiente de navegador por caminho de dados do usuário.

### Flags do navegador

As flags do navegador WebView2 controlam recursos experimentais e o comportamento em **todas as janelas** do aplicativo. Elas devem ser definidas em `application.Options.Windows`:

```go
app := application.New(application.Options{
    Name: "My App",
    Windows: application.WindowsOptions{
        // Enable experimental WebView2 features
        EnabledFeatures: []string{
            "msWebView2EnableDraggableRegions",
        },

        // Disable specific features
        DisabledFeatures: []string{
            "msSmartScreenProtection",  // Always disabled by Wails
        },

        // Additional Chromium command-line arguments
        AdditionalBrowserArgs: []string{
            "--disable-gpu",
            "--remote-debugging-port=9222",
        },
    },
})
```

**EnabledFeatures** (`[]string`)

- Lista de flags de recursos do WebView2 a serem habilitadas
- Consulte [flags do navegador WebView2](https://learn.microsoft.com/en-us/microsoft-edge/webview2/concepts/webview-features-flags) para ver as flags disponíveis
- Exemplo: `"msWebView2EnableDraggableRegions"`

**DisabledFeatures** (`[]string`)

- Lista de flags de recursos do WebView2 a serem desabilitadas
- O Wails desabilita automaticamente `msSmartScreenProtection`
- Exemplo: `"msExperimentalFeature"`

**AdditionalBrowserArgs** (`[]string`)

- Argumentos de linha de comando do Chromium passados ao processo do navegador
- Devem incluir o prefixo `--` (por exemplo, `"--remote-debugging-port=9222"`)
- Consulte [opções de linha de comando do Chromium](https://peter.sh/experiments/chromium-command-line-switches/) para ver os argumentos disponíveis

@note{type="caution" title="Importante"}
Essas flags se aplicam globalmente a TODAS as janelas, pois o WebView2 compartilha um único ambiente de navegador por caminho de dados do usuário. Não é possível usar flags de navegador diferentes para janelas diferentes.

@end

**Exemplo completo:**

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "My App",
        Windows: application.WindowsOptions{
            // Enable draggable regions feature
            EnabledFeatures: []string{
                "msWebView2EnableDraggableRegions",
            },
            // Enable remote debugging
            AdditionalBrowserArgs: []string{
                "--remote-debugging-port=9222",
            },
        },
    })

    // All windows will use the browser flags configured above
    window := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:  "Main Window",
        Width:  1024,
        Height: 768,
    })

    window.Show()
    app.Run()
}
```

## Exemplo completo

Veja uma configuração de janela pronta para produção:

```go
package main

import (
    _ "embed"
    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var assets embed.FS

//go:embed icon.png
var icon []byte

func main() {
    app := application.New(application.Options{
        Name: "My Application",
    })

    window := app.Window.NewWithOptions(application.WebviewWindowOptions{
        // Identity
        Name:  "main-window",
        Title: "My Application",

        // Size and Position
        Width:     1200,
        Height:    800,
        MinWidth:  800,
        MinHeight: 600,

        // Initial State
        StartState: application.WindowStateNormal,

        // Appearance
        BackgroundColour: application.NewRGB(255, 255, 255),

        // Platform-Specific (per-window structs)
        Mac: application.MacWindow{
            TitleBar: application.MacTitleBar{
                AppearsTransparent: true,
            },
            Backdrop: application.MacBackdropTranslucent,
        },

        Windows: application.WindowsWindow{
            BackdropType: application.Mica,
            DisableIcon:  false,
        },

        Linux: application.LinuxWindow{
            Icon: icon,
        },
    })

    window.Center()
    window.Show()

    app.Run()
}
```

Os recursos do frontend são servidos no nível do aplicativo (`application.Options.Assets`), e não por janela.

## Próximos passos

- [Conceitos básicos de janelas](/features/windows/basics/) - Como criar e controlar janelas
- [Várias janelas](/features/windows/multiple/) - Padrões para várias janelas
- [Janelas sem moldura](/features/windows/frameless/) - Decoração de janela personalizada
- [Eventos de janela](/features/windows/events/) - Eventos do ciclo de vida

---

**Dúvidas?** Pergunte no [Discord](https://discord.gg/JDdSxwjhGf) ou consulte os [exemplos](https://github.com/wailsapp/wails/tree/master/v3/examples).
