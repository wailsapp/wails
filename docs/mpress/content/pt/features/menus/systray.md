---
title: "Menus da bandeja do sistema"
description: "Adicione integração com a bandeja do sistema (área de notificação) ao seu aplicativo"
slug: "features/menus/systray"
sourcePath: "features/menus/systray.md"
---

## Menus da bandeja do sistema

O Wails fornece **APIs unificadas para a bandeja do sistema** que funcionam em todas as plataformas. Crie ícones de bandeja com menus, associe janelas e processe cliques com o comportamento nativo de cada plataforma para aplicativos em segundo plano, serviços e utilitários de acesso rápido.

![Um menu da bandeja do sistema do Wails aberto pela barra de menus do macOS](/assets/screenshots/systray-menu-macos.png)

No macOS, um item da bandeja do sistema do Wails aparece na barra de menus e abre um menu nativo. O exemplo inclui itens desabilitados, caixas de seleção, botões de opção, submenus e ações.

## Início rápido

```go
package main

import (
    _ "embed"
    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets/icon.png
var icon []byte

func main() {
    app := application.New(application.Options{
        Name: "Tray App",
    })

    // Create system tray
    systray := app.SystemTray.New()
    systray.SetIcon(icon)
    systray.SetLabel("My App")

    // Add menu
    menu := app.NewMenu()
    menu.Add("Show").OnClick(func(ctx *application.Context) {
        // Show main window
    })
    menu.Add("Quit").OnClick(func(ctx *application.Context) {
        app.Quit()
    })
    systray.SetMenu(menu)

    // Create hidden window
    window := app.Window.New()
    window.Hide()

    app.Run()
}
```

**Resultado:** ícone da bandeja do sistema com menu em todas as plataformas.

## Como criar uma bandeja do sistema

### Bandeja do sistema básica

```go
// Create system tray
systray := app.SystemTray.New()

// Set icon
systray.SetIcon(iconBytes)

// Set label (macOS) / tooltip (Windows)
systray.SetLabel("My Application")
```

### Com ícone

Os ícones devem ser incorporados:

```go
import _ "embed"

//go:embed assets/icon.png
var icon []byte

//go:embed assets/icon-dark.png
var iconDark []byte

func main() {
    app := application.New(application.Options{
        Name: "My App",
    })

    systray := app.SystemTray.New()
    systray.SetIcon(icon)
    systray.SetDarkModeIcon(iconDark)  // Windows and macOS dark mode
    
    app.Run()
}
```

**Requisitos dos ícones:**

| Plataforma | Tamanho | Formato | Observações |
| --- | --- | --- | --- |
| **Windows** | 16x16 ou 32x32 | PNG, ICO | Área de notificação |
| **macOS** | 18x18 a 22x22 | PNG | Barra de menus; recomenda-se usar um ícone de modelo |
| **Linux** | 22x22 a 48x48 | PNG, SVG | Varia conforme o ambiente de desktop |

### Ícones de modelo (macOS)

Os ícones de modelo se adaptam automaticamente aos modos claro e escuro:

```go
systray.SetTemplateIcon(iconBytes)
```

**Diretrizes para ícones de modelo:**

- Use apenas as cores preta e transparente
- O preto se torna branco no modo escuro
- Adicione o sufixo `Template` ao nome do arquivo: `iconTemplate.png`
- [Guia de design](https://bjango.com/articles/designingmenubarextras/)

## Como adicionar menus

Os menus da bandeja do sistema funcionam como os menus do aplicativo:

```go
menu := app.NewMenu()

// Add items
menu.Add("Open").OnClick(func(ctx *application.Context) {
    showMainWindow()
})

menu.AddSeparator()

menu.AddCheckbox("Start at Login", false).OnClick(func(ctx *application.Context) {
    enabled := ctx.ClickedMenuItem().Checked()
    setStartAtLogin(enabled)
})

menu.AddSeparator()

menu.Add("Quit").OnClick(func(ctx *application.Context) {
    app.Quit()
})

// Set menu
systray.SetMenu(menu)
```

Consulte **todos os tipos de item de menu** na [Referência de menus](/features/menus/reference/).

## Como associar janelas

Associe uma janela ao ícone da bandeja para exibi-la e ocultá-la automaticamente:

```go
// Create window
window := app.Window.New()

// Attach to tray
systray.AttachWindow(window)

// Configure behaviour — these are setters that return the receiver for chaining.
systray.WindowOffset(10)                          // Pixels from tray icon
systray.WindowDebounce(200 * time.Millisecond)    // Click debounce
```

**Comportamento:**

- A janela inicia oculta
- **Clique com o botão esquerdo no ícone da bandeja** → Alterna a visibilidade da janela
- **Clique com o botão direito no ícone da bandeja** → Exibe o menu (se configurado)
- A janela é posicionada perto do ícone da bandeja

**Exemplo: janela pop-up**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:           "Quick Access",
    Width:           300,
    Height:          400,
    Frameless:       true, // No title bar
    AlwaysOnTop:     true, // Stay on top
    HideOnFocusLost: true, // Dismiss when another window receives focus
    HideOnEscape:    true, // Dismiss when the user presses Escape
})

systray.AttachWindow(window)
systray.WindowOffset(5)
```

`HideOnFocusLost` é útil para pop-ups da bandeja no Windows, no macOS e em desktops Linux com foco ao clicar. O Wails desabilita esse comportamento em ambientes Linux que usam foco seguindo o mouse (incluindo configurações comuns do Hyprland, Sway e i3), nos quais afastar o ponteiro do pop-up poderia ocultá-lo antes que ele pudesse ser usado. `HideOnEscape` continua disponível nesses ambientes.

Os comportamentos de clique com os botões esquerdo e direito descritos acima são padrões inteligentes. Um manipulador `OnClick` ou `OnRightClick` explícito substitui o padrão correspondente. Para verificações de plataforma e casos extremos, consulte a [suíte manual da bandeja do sistema](https://github.com/wailsapp/wails/tree/master/v3/test/manual/systray) e o [exemplo de teste de estresse da bandeja do sistema](https://github.com/wailsapp/wails/tree/master/v3/examples/systray-stress).

## Manipuladores de clique

Processe os cliques no ícone da bandeja:

```go
systray := app.SystemTray.New()

// Left click
systray.OnClick(func() {
    fmt.Println("Tray icon clicked")
})

// Right click
systray.OnRightClick(func() {
    fmt.Println("Tray icon right-clicked")
})

// Double click
systray.OnDoubleClick(func() {
    fmt.Println("Tray icon double-clicked")
})

// Mouse enter/leave
systray.OnMouseEnter(func() {
    fmt.Println("Mouse entered tray icon")
})

systray.OnMouseLeave(func() {
    fmt.Println("Mouse left tray icon")
})
```

**Compatibilidade com plataformas:**

| Evento | Windows | macOS | Linux |
| --- | --- | --- | --- |
| OnClick | ✅ | ✅ | ✅ |
| OnRightClick | ✅ | ✅ | ✅ |
| OnDoubleClick | ✅ | ✅ | ⚠️ Varia |
| OnMouseEnter | ✅ | ✅ | ⚠️ Varia |
| OnMouseLeave | ✅ | ✅ | ⚠️ Varia |

## Atualizações dinâmicas

Atualize dinamicamente o ícone e o menu da bandeja:

### Alterar o ícone

```go
var isActive bool

func updateTrayIcon() {
    if isActive {
        systray.SetIcon(activeIcon)
        systray.SetLabel("Active")
    } else {
        systray.SetIcon(inactiveIcon)
        systray.SetLabel("Inactive")
    }
}
```

### Atualizar o menu

```go
var isPaused bool

pauseMenuItem := menu.Add("Pause")

pauseMenuItem.OnClick(func(ctx *application.Context) {
    isPaused = !isPaused
    
    if isPaused {
        pauseMenuItem.SetLabel("Resume")
    } else {
        pauseMenuItem.SetLabel("Pause")
    }
    
    menu.Update()  // Important!
})
```

@note{type="caution" title="Sempre chame Update()"}
Após alterar o estado do menu, **chame `menu.Update()`**. Consulte a [Referência de menus](/features/menus/reference/#enabled-state).

@end

### Reconstruir o menu

Para alterações significativas, recrie o menu inteiro:

```go
func rebuildTrayMenu(status string) {
    menu := app.NewMenu()
    
    // Status-specific items
    switch status {
    case "syncing":
        menu.Add("Syncing...").SetEnabled(false)
        menu.Add("Pause Sync").OnClick(pauseSync)
    case "synced":
        menu.Add("Up to date ✓").SetEnabled(false)
        menu.Add("Sync Now").OnClick(startSync)
    case "error":
        menu.Add("Sync Error").SetEnabled(false)
        menu.Add("Retry").OnClick(retrySync)
    }
    
    menu.AddSeparator()
    menu.Add("Quit").OnClick(func(ctx *application.Context) {
        app.Quit()
    })
    
    systray.SetMenu(menu)
}
```

## Recursos específicos de cada plataforma

@tabs{sync-key="platform"}
[macOS]
**Integração com a barra de menus:**

```go
// Set label (appears next to icon)
systray.SetLabel("My App")

// Use template icon (adapts to dark mode)
systray.SetTemplateIcon(iconBytes)

// Set icon position — uses AppKit NSImage placement constants.
systray.SetIconPosition(application.NSImageRight)
```

**Posições do ícone** (correspondentes a `NSImagePosition`):

- `application.NSImageLeft` — Ícone à esquerda do rótulo.
- `application.NSImageRight` — Ícone à direita do rótulo.
- `application.NSImageOnly` — Somente o ícone, sem rótulo.
- `application.NSImageNone` — Somente o rótulo, sem ícone.

**Práticas recomendadas:**

- Use ícones de modelo (preto + transparente)
- Mantenha os rótulos curtos (3-5 caracteres)
- De 18x18 a 22x22 pixels para telas Retina
- Teste nos modos claro e escuro

[Windows]
**Integração com a área de notificação:**

```go
// Set tooltip (appears on hover)
systray.SetTooltip("My Application")

// Or use SetLabel (same as tooltip on Windows)
systray.SetLabel("My Application")

// Show/Hide functionality (fully functional)
systray.Show()  // Show tray icon
systray.Hide()  // Hide tray icon
```

**Requisitos do ícone:**

- 16x16 ou 32x32 pixels
- Formato PNG ou ICO
- Fundo transparente

**Limites da dica de ferramenta:**

- No máximo 127 caracteres UTF-16
- Dicas de ferramenta mais longas serão truncadas
- Seja conciso para oferecer a melhor experiência

**Recursos da plataforma:**

- O ícone da bandeja permanece após reinicializações do Windows Explorer
- Métodos Show() e Hide() totalmente funcionais
- Gerenciamento adequado do ciclo de vida

**Práticas recomendadas:**

- Use 32x32 em telas com DPI alto
- Mantenha as dicas de ferramenta com menos de 127 caracteres
- Teste em diferentes versões do Windows
- Considere a área de estouro da área de notificação
- Use Show/Hide para controlar condicionalmente a visibilidade do ícone da bandeja

[Linux]
**Integração com a bandeja do sistema:**

Usa a especificação StatusNotifierItem (na maioria dos ambientes de desktop modernos).

```go
systray.SetIcon(iconBytes)
systray.SetLabel("My App")
```

**Compatibilidade com ambientes de desktop:**

- **GNOME**: barra superior (com extensão)
- **KDE Plasma**: bandeja do sistema
- **XFCE**: área de notificação
- **Outros**: varia

**Práticas recomendadas:**

- Use 22x22 ou 24x24 pixels
- Ícones SVG oferecem melhor redimensionamento
- Teste nos ambientes de desktop de destino
- Forneça uma alternativa para ambientes de desktop não compatíveis

@end

## Exemplo completo

Veja um aplicativo de bandeja do sistema pronto para produção:

```go
package main

import (
    _ "embed"
    "fmt"
    "time"
    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets/icon.png
var icon []byte

//go:embed assets/icon-active.png
var iconActive []byte

type TrayApp struct {
    app     *application.App
    systray *application.SystemTray
    window  *application.WebviewWindow
    menu    *application.Menu
    isActive bool
}

func main() {
    app := application.New(application.Options{
        Name: "Tray Application",
        Mac: application.MacOptions{
            ApplicationShouldTerminateAfterLastWindowClosed: false,
        },
    })

    trayApp := &TrayApp{app: app}
    trayApp.setup()

    app.Run()
}

func (t *TrayApp) setup() {
    // Create system tray
    t.systray = t.app.SystemTray.New()
    t.systray.SetIcon(icon)
    t.systray.SetLabel("Inactive")
    
    // Create menu
    t.createMenu()
    
    // Create window (hidden by default)
    t.window = t.app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:  "Tray Application",
        Width:  400,
        Height: 600,
        Hidden: true,
    })
    
    // Attach window to tray
    t.systray.AttachWindow(t.window)
    t.systray.WindowOffset(10)
    
    // Handle tray clicks
    t.systray.OnRightClick(func() {
        t.systray.OpenMenu()
    })
    
    // Start background task
    go t.backgroundTask()
}

func (t *TrayApp) createMenu() {
    t.menu = t.app.NewMenu()
    
    // Status item (disabled)
    statusItem := t.menu.Add("Status: Inactive")
    statusItem.SetEnabled(false)
    
    t.menu.AddSeparator()
    
    // Toggle active
    t.menu.Add("Start").OnClick(func(ctx *application.Context) {
        t.toggleActive()
    })
    
    // Show window
    t.menu.Add("Show Window").OnClick(func(ctx *application.Context) {
        t.window.Show()
        t.window.Focus()
    })
    
    t.menu.AddSeparator()
    
    // Settings
    t.menu.AddCheckbox("Start at Login", false).OnClick(func(ctx *application.Context) {
        enabled := ctx.ClickedMenuItem().Checked()
        t.setStartAtLogin(enabled)
    })
    
    t.menu.AddSeparator()
    
    // Quit
    t.menu.Add("Quit").OnClick(func(ctx *application.Context) {
        t.app.Quit()
    })
    
    t.systray.SetMenu(t.menu)
}

func (t *TrayApp) toggleActive() {
    t.isActive = !t.isActive
    t.updateTray()
}

func (t *TrayApp) updateTray() {
    if t.isActive {
        t.systray.SetIcon(iconActive)
        t.systray.SetLabel("Active")
    } else {
        t.systray.SetIcon(icon)
        t.systray.SetLabel("Inactive")
    }
    
    // Rebuild menu with new status
    t.createMenu()
}

func (t *TrayApp) backgroundTask() {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        if t.isActive {
            fmt.Println("Background task running...")
            // Do work
        }
    }
}

func (t *TrayApp) setStartAtLogin(enabled bool) {
    // Implementation varies by platform
    fmt.Printf("Start at login: %v\n", enabled)
}
```

## Controle de visibilidade

Mostre ou oculte dinamicamente o ícone da bandeja:

```go
// Hide tray icon
systray.Hide()

// Show tray icon
systray.Show()
```

Não há um método getter `IsVisible()` — se precisar dessa informação, controle a visibilidade no próprio estado do aplicativo.

**Compatibilidade com plataformas:**

| Plataforma | Hide() | Show() | Observações |
| --- | --- | --- | --- |
| **Windows** | ✅ | ✅ | Totalmente funcional — o ícone aparece/desaparece da área de notificação |
| **macOS** | ✅ | ✅ | O item da barra de menus é exibido/ocultado |
| **Linux** | ✅ | ✅ | Varia conforme o ambiente de desktop |

**Casos de uso:**

- Ocultar temporariamente o ícone da bandeja conforme a preferência do usuário
- Modo sem interface gráfica, com o ícone da bandeja aparecendo somente quando necessário
- Alternar a visibilidade conforme o estado do aplicativo

**Exemplo — Visibilidade condicional do ícone da bandeja:**

```go
func (t *TrayApp) setTrayVisibility(visible bool) {
    if visible {
        t.systray.Show()
    } else {
        t.systray.Hide()
    }
}

// Show tray only when updates are available
func (t *TrayApp) checkForUpdates() {
    if hasUpdates {
        t.systray.Show()
        t.systray.SetLabel("Update Available")
    } else {
        t.systray.Hide()
    }
}
```

## Limpeza

Destrua o ícone da bandeja quando terminar:

```go
// In OnShutdown
app := application.New(application.Options{
    OnShutdown: func() {
        if systray != nil {
            systray.Destroy()
        }
    },
})
```

**Importante:** Sempre destrua a bandeja do sistema ao encerrar o aplicativo para liberar recursos.

## Práticas recomendadas

### ✅ Faça

- **Use ícones de modelo no macOS** — Eles se adaptam ao modo escuro
- **Mantenha os rótulos curtos** — No máximo 3-5 caracteres
- **Forneça dicas de ferramenta no Windows** — Isso ajuda os usuários a identificar seu aplicativo
- **Teste em todas as plataformas** — O comportamento varia
- **Trate os cliques adequadamente** — Use o clique com o botão esquerdo para a ação principal e o clique com o botão direito para o menu
- **Atualize o ícone para indicar o status** — O feedback visual é importante
- **Destrua ao encerrar** — Libere os recursos

### ❌ Não faça

- **Não use ícones grandes** — Siga as diretrizes da plataforma
- **Não use rótulos longos** — Eles são truncados
- **Não se esqueça do modo escuro** — Teste no modo escuro do Windows e do macOS
- **Não bloqueie os manipuladores de clique** — Mantenha-os rápidos
- **Não se esqueça de menu.Update()** — Após alterar o estado do menu
- **Não presuma que haja suporte à bandeja do sistema** — Alguns ambientes de desktop Linux não oferecem suporte

## Solução de problemas

### O ícone da bandeja não aparece

**Possíveis causas:**

1. Formato de ícone não compatível
2. Ícone grande ou pequeno demais
3. Bandeja do sistema não compatível (Linux)

**Solução:**

Não há um auxiliar `SystemTraySupported()`; em vez disso, crie a bandeja, verifique a plataforma e forneça uma alternativa adequada:

```go
// Probe support: on Linux without a notification-area extension, the tray
// will simply not appear. Defensive code can fall back to window-only mode
// based on runtime.GOOS or after a short timeout if no tray events arrive.
systray := app.SystemTray.New()
systray.SetIcon(iconBytes)
```

### O ícone tem uma aparência incorreta no macOS

**Causa:** O ícone de modelo não está sendo usado

**Solução:**

```go
// Use template icon
systray.SetTemplateIcon(iconBytes)

// Or design icon as template (black + transparent)
```

### O menu não é atualizado

**Causa:** A chamada a `menu.Update()` foi esquecida

**Solução:**

```go
menuItem.SetLabel("New Label")
menu.Update()  // Add this!
```

## Próximas etapas

@cards{cols="2"}
📖 Referência de menus
Referência completa dos tipos e das propriedades dos itens de menu.

[Saiba mais →](/features/menus/reference/)

---
☰ Menus do aplicativo
Crie barras de menus para o aplicativo.

[Saiba mais →](/features/menus/application/)

---
◆ Menus de contexto
Crie menus de contexto acionados com o botão direito.

[Saiba mais →](/features/menus/context/)

---
📖 Exemplo de bandeja do sistema
Explore um aplicativo completo com bandeja do sistema.

[Saiba mais →](https://github.com/wailsapp/wails/tree/master/v3/examples/systray-basic)

@end

---

**Dúvidas?** Pergunte no [Discord](https://discord.gg/JDdSxwjhGf) ou consulte os [exemplos de bandeja do sistema](https://github.com/wailsapp/wails/tree/master/v3/examples/systray-basic).
