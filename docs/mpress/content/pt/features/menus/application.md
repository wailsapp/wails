---
title: "Menus do aplicativo"
description: "Crie barras de menus nativas para seu aplicativo desktop"
slug: "features/menus/application"
sourcePath: "features/menus/application.md"
---

## O problema

Aplicativos desktop profissionais precisam de barras de menus — Arquivo, Editar, Visualizar, Ajuda. No entanto, os menus funcionam de forma diferente em cada plataforma:

- **macOS**: barra de menus global na parte superior da tela
- **Windows**: barra de menus na barra de título da janela
- **Linux**: varia conforme o ambiente de desktop

Criar manualmente menus adequados a cada plataforma é trabalhoso e propenso a erros.

## A solução do Wails

O Wails fornece uma **API unificada** que cria automaticamente menus nativos de cada plataforma. Escreva uma vez e obtenha comportamento nativo em todas as plataformas.

![Um menu de aplicativo Wails no macOS com itens padrão, caixas de seleção, botões de opção e submenus](/assets/screenshots/application-menu-macos.png)

No macOS, o menu do aplicativo é colocado na barra de menus global. Esta captura mostra o menu nativo renderizado pela API de menus do Wails, incluindo itens desabilitados, caixas de seleção, botões de opção e submenus.

## Início rápido

```go
package main

import (
    "runtime"
    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "My App",
    })

    // Create menu
    menu := app.NewMenu()

    // Add standard menus (platform-appropriate)
    if runtime.GOOS == "darwin" {
        menu.AddRole(application.AppMenu)  // macOS only
    }
    menu.AddRole(application.FileMenu)
    menu.AddRole(application.EditMenu)
    menu.AddRole(application.WindowMenu)
    menu.AddRole(application.HelpMenu)

    // Set the application menu
    app.Menu.Set(menu)

    // Create window with UseApplicationMenu to inherit the menu on Windows/Linux
    app.Window.NewWithOptions(application.WebviewWindowOptions{
        UseApplicationMenu: true,
    })

    app.Run()
}
```

**É só isso!** Agora você tem menus nativos de cada plataforma com itens padrão. A opção `UseApplicationMenu` garante que as janelas no Windows e no Linux exibam o menu sem código adicional.

## Criação de menus

### Criação básica de menus

```go
// Create a new menu
menu := app.NewMenu()

// Add a top-level menu
fileMenu := menu.AddSubmenu("File")

// Add menu items
fileMenu.Add("New").OnClick(func(ctx *application.Context) {
    // Handle New
})

fileMenu.Add("Open").OnClick(func(ctx *application.Context) {
    // Handle Open
})

fileMenu.AddSeparator()

fileMenu.Add("Quit").OnClick(func(ctx *application.Context) {
    app.Quit()
})
```

### Definição do menu

**Abordagem recomendada** — use `UseApplicationMenu` para manter a consistência entre plataformas:

```go
// Set the application menu once
app.Menu.Set(menu)

// Create windows that inherit the menu on Windows/Linux
app.Window.NewWithOptions(application.WebviewWindowOptions{
    UseApplicationMenu: true,  // Window uses the app menu
})
```

Essa abordagem:

- No **macOS**: o menu aparece na parte superior da tela (comportamento padrão)
- No **Windows/Linux**: cada janela com `UseApplicationMenu: true` exibe o menu do aplicativo

**Detalhes específicos de cada plataforma:**

@tabs{sync-key="platform"}
[macOS]
**Barra de menus global** (uma por aplicativo):

```go
app.Menu.Set(menu)
```

O menu aparece na parte superior da tela e permanece lá mesmo quando todas as janelas estão fechadas. A opção `UseApplicationMenu` não tem efeito no macOS, pois todos os aplicativos usam o menu global.

[Windows]
**Barra de menus por janela**:

```go
// Option 1: Use application menu (recommended)
app.Menu.Set(menu)
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    UseApplicationMenu: true,
})

// Option 2: Set menu directly on window
window.SetMenu(menu)
```

Cada janela pode ter seu próprio menu ou herdar o menu do aplicativo. O menu aparece na barra de título da janela.

[Linux]
**Barra de menus por janela** (geralmente):

```go
// Option 1: Use application menu (recommended)
app.Menu.Set(menu)
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    UseApplicationMenu: true,
})

// Option 2: Set menu directly on window
window.SetMenu(menu)
```

O comportamento varia conforme o ambiente de desktop. Alguns ambientes, como o Unity, oferecem suporte a menus globais.

@end

@note{type="tip" title="Simplifique os menus multiplataforma"}
Usar `UseApplicationMenu: true` elimina a necessidade de código específico para cada plataforma, como:

```go
// Old approach - no longer needed
if runtime.GOOS == "darwin" {
    app.Menu.Set(menu)
} else {
    window.SetMenu(menu)
}
```

@end

**Menus personalizados por janela:**

Se uma janela precisar de um menu diferente do menu do aplicativo, defina-o diretamente:

```go
window.SetMenu(customMenu)  // Overrides UseApplicationMenu
```

## Funções de menu

O Wails fornece **funções de menu predefinidas** que criam automaticamente estruturas de menu adequadas a cada plataforma.

### Funções disponíveis

| Função | Descrição | Observações sobre as plataformas |
| --- | --- | --- |
| `AppMenu` | Menu do aplicativo com Sobre, Preferências e Sair | **Somente macOS** |
| `FileMenu` | Operações de arquivo (Novo, Abrir, Salvar etc.) | Todas as plataformas |
| `EditMenu` | Edição de texto (Desfazer, Refazer, Recortar, Copiar, Colar) | Todas as plataformas |
| `WindowMenu` | Gerenciamento de janelas (Minimizar, Ampliar etc.) | Todas as plataformas |
| `HelpMenu` | Ajuda e informações | Todas as plataformas |

### Uso das funções

```go
menu := app.NewMenu()

// macOS: Add application menu
if runtime.GOOS == "darwin" {
    menu.AddRole(application.AppMenu)
}

// All platforms: Add standard menus
menu.AddRole(application.FileMenu)
menu.AddRole(application.EditMenu)
menu.AddRole(application.WindowMenu)
menu.AddRole(application.HelpMenu)
```

**O que você obtém:**

@tabs{sync-key="platform"}
[macOS]
**AppMenu** (com o nome do aplicativo):

- Sobre [Nome do aplicativo]
- Preferências... (⌘,)
- ---
- Serviços
- ---
- Ocultar [Nome do aplicativo] (⌘H)
- Ocultar Outros (⌥⌘H)
- Mostrar Todos
- ---
- Encerrar [Nome do aplicativo] (⌘Q)

**FileMenu**:

- Novo (⌘N)
- Abrir... (⌘O)
- ---
- Fechar Janela (⌘W)

**EditMenu**:

- Desfazer (⌘Z)
- Refazer (⇧⌘Z)
- ---
- Recortar (⌘X)
- Copiar (⌘C)
- Colar (⌘V)
- Selecionar tudo (⌘A)

**WindowMenu**:

- Minimizar (⌘M)
- Zoom
- ---
- Trazer todas para a frente

**HelpMenu**:

- Ajuda do [Nome do aplicativo]

[Windows]
**FileMenu**:

- Novo (Ctrl+N)
- Abrir... (Ctrl+O)
- ---
- Sair (Alt+F4)

**EditMenu**:

- Desfazer (Ctrl+Z)
- Refazer (Ctrl+Y)
- ---
- Recortar (Ctrl+X)
- Copiar (Ctrl+C)
- Colar (Ctrl+V)
- Selecionar tudo (Ctrl+A)

**WindowMenu**:

- Minimizar
- Maximizar

**HelpMenu**:

- Sobre o [Nome do aplicativo]

[Linux]
Semelhante ao Windows, mas os atalhos de teclado podem variar conforme o ambiente de desktop.

@end

### Personalização dos menus de função

`Menu.AddRole(role)` retorna o menu **receptor** (o menu de nível superior), e **não** o submenu da função. Para adicionar itens ao submenu da função, procure o item de função inserido com `FindByRole` e chame `GetSubmenu()` nele:

```go
menu.AddRole(application.FileMenu)

fileMenu := menu.FindByRole(application.FileMenu).GetSubmenu()
fileMenu.Add("Import...").OnClick(handleImport)
fileMenu.Add("Export...").OnClick(handleExport)
```

## Menus personalizados

Crie seus próprios menus para recursos específicos do aplicativo:

```go
// Add a custom top-level menu
toolsMenu := menu.AddSubmenu("Tools")

// Add items
toolsMenu.Add("Settings").OnClick(func(ctx *application.Context) {
    showSettingsWindow()
})

toolsMenu.AddSeparator()

// Add checkbox
toolsMenu.AddCheckbox("Dark Mode", false).OnClick(func(ctx *application.Context) {
    isDark := ctx.ClickedMenuItem().Checked()
    setTheme(isDark)
})

// Add radio group
toolsMenu.AddRadio("Small", true).OnClick(handleFontSize)
toolsMenu.AddRadio("Medium", false).OnClick(handleFontSize)
toolsMenu.AddRadio("Large", false).OnClick(handleFontSize)

// Add submenu
advancedMenu := toolsMenu.AddSubmenu("Advanced")
advancedMenu.Add("Configure...").OnClick(showAdvancedSettings)
```

**Para ver mais tipos de itens de menu**, consulte a [Referência de menus](/features/menus/reference/).

## Menus dinâmicos

Atualize os menus com base no estado do aplicativo:

### Ativar/desativar itens

```go
var saveMenuItem *application.MenuItem

func createMenu() {
    menu := app.NewMenu()
    fileMenu := menu.AddSubmenu("File")
    
    saveMenuItem = fileMenu.Add("Save")
    saveMenuItem.SetEnabled(false)  // Initially disabled
    saveMenuItem.OnClick(handleSave)
    
    app.Menu.Set(menu)
}

func onDocumentChanged() {
    saveMenuItem.SetEnabled(hasUnsavedChanges())
    menu.Update()  // Important!
}
```

@note{type="caution" title="Sempre chame menu.Update()"}
Depois de alterar o estado do menu (ativado/desativado, rótulo ou marcação), **sempre chame `menu.Update()`**. Isso é especialmente importante no Windows, onde os menus são reconstruídos.

Consulte a [Referência de menus](/features/menus/reference/#enabled-state) para obter detalhes.

@end

### Alterar rótulos

```go
updateMenuItem := menu.Add("Check for Updates")

updateMenuItem.OnClick(func(ctx *application.Context) {
    updateMenuItem.SetLabel("Checking...")
    menu.Update()
    
    checkForUpdates()
    
    updateMenuItem.SetLabel("Check for Updates")
    menu.Update()
})
```

### Reconstruir menus

Para alterações significativas, reconstrua todo o menu:

```go
func rebuildFileMenu() {
    menu := app.NewMenu()
    fileMenu := menu.AddSubmenu("File")
    
    fileMenu.Add("New").OnClick(handleNew)
    fileMenu.Add("Open").OnClick(handleOpen)
    
    // Add recent files dynamically
    if hasRecentFiles() {
        recentMenu := fileMenu.AddSubmenu("Open Recent")
        for _, file := range getRecentFiles() {
            filePath := file  // Capture for closure
            recentMenu.Add(filepath.Base(file)).OnClick(func(ctx *application.Context) {
                openFile(filePath)
            })
        }
        recentMenu.AddSeparator()
        recentMenu.Add("Clear Recent").OnClick(clearRecentFiles)
    }
    
    fileMenu.AddSeparator()
    fileMenu.Add("Quit").OnClick(func(ctx *application.Context) {
        app.Quit()
    })
    
    app.Menu.Set(menu)
}
```

## Controle de janelas por meio de menus

Os itens de menu podem controlar janelas:

```go
viewMenu := menu.AddSubmenu("View")

// Toggle fullscreen
viewMenu.Add("Toggle Fullscreen").OnClick(func(ctx *application.Context) {
    if window, ok := app.Window.GetByName("main"); ok {
        window.ToggleFullscreen()
    }
})

// Zoom controls
viewMenu.Add("Zoom In").SetAccelerator("CmdOrCtrl++").OnClick(func(ctx *application.Context) {
    // Increase zoom
})

viewMenu.Add("Zoom Out").SetAccelerator("CmdOrCtrl+-").OnClick(func(ctx *application.Context) {
    // Decrease zoom
})

viewMenu.Add("Reset Zoom").SetAccelerator("CmdOrCtrl+0").OnClick(func(ctx *application.Context) {
    // Reset zoom
})
```

**Obtenha a janela ativa:**

```go
menuItem.OnClick(func(ctx *application.Context) {
    window := application.Get().Window.Current() // the window the menu was invoked from
    // Use window
})
```

## Considerações específicas de cada plataforma

### macOS

**Comportamento da barra de menus:**

- Aparece na **parte superior da tela** (global)
- Permanece quando todas as janelas estão fechadas
- O primeiro menu é **sempre o menu do aplicativo**
- Use `menu.AddRole(application.AppMenu)` para itens padrão

**Locais padrão:**

- **Sobre**: menu do aplicativo
- **Preferências**: menu do aplicativo (⌘,)
- **Encerrar**: menu do aplicativo (⌘Q)
- **Ajuda**: menu Ajuda

**Exemplo:**

```go
if runtime.GOOS == "darwin" {
    menu.AddRole(application.AppMenu)  // Adds About, Preferences, Quit
    
    // Don't add Quit to File menu on macOS
    // Don't add About to Help menu on macOS
}
```

### Windows

**Comportamento da barra de menus:**

- Aparece na **barra de título da janela**
- Cada janela tem seu próprio menu
- Não há menu do aplicativo

**Locais padrão:**

- **Sair**: menu Arquivo (Alt+F4)
- **Configurações**: menu Ferramentas ou Editar
- **Sobre**: menu Ajuda

**Exemplo:**

```go
if runtime.GOOS == "windows" {
    menu.AddRole(application.FileMenu) // Exit is added automatically
    menu.AddRole(application.HelpMenu) // About is added automatically
}
```

### Linux

**Comportamento da barra de menus:**

- Geralmente é específica de cada janela (como no Windows)
- Alguns ambientes de desktop oferecem suporte a menus globais (Unity, GNOME com extensão)
- A aparência varia conforme o ambiente de desktop

**Prática recomendada:** siga as convenções do Windows e teste nos ambientes de desktop de destino.

## Exemplo completo

Esta é uma estrutura de menu pronta para produção:

```go
package main

import (
    "runtime"
    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "My Application",
    })

    // Create and set menu
    createMenu(app)

    // Create main window with UseApplicationMenu for cross-platform menu support
    app.Window.NewWithOptions(application.WebviewWindowOptions{
        UseApplicationMenu: true,
    })

    app.Run()
}

func createMenu(app *application.App) {
    menu := app.NewMenu()

    // Platform-specific application menu (macOS only)
    if runtime.GOOS == "darwin" {
        menu.AddRole(application.AppMenu)
    }

    // File menu — AddRole returns the receiver menu, not the role submenu.
    // To add items into the File submenu, look it up via FindByRole + GetSubmenu.
    menu.AddRole(application.FileMenu)
    fileMenu := menu.FindByRole(application.FileMenu).GetSubmenu()
    fileMenu.Add("Import...").SetAccelerator("CmdOrCtrl+I").OnClick(handleImport)
    fileMenu.Add("Export...").SetAccelerator("CmdOrCtrl+E").OnClick(handleExport)

    // Edit menu
    menu.AddRole(application.EditMenu)

    // View menu
    viewMenu := menu.AddSubmenu("View")
    viewMenu.Add("Toggle Fullscreen").SetAccelerator("F11").OnClick(toggleFullscreen)
    viewMenu.AddSeparator()
    viewMenu.AddCheckbox("Show Sidebar", true).OnClick(toggleSidebar)
    viewMenu.AddCheckbox("Show Toolbar", true).OnClick(toggleToolbar)

    // Tools menu
    toolsMenu := menu.AddSubmenu("Tools")
    
    // Settings location varies by platform
    if runtime.GOOS == "darwin" {
        // On macOS, Preferences is in Application menu (added by AppMenu role)
    } else {
        toolsMenu.Add("Settings").SetAccelerator("CmdOrCtrl+,").OnClick(showSettings)
    }
    
    toolsMenu.AddSeparator()
    toolsMenu.AddCheckbox("Dark Mode", false).OnClick(toggleDarkMode)

    // Window menu
    menu.AddRole(application.WindowMenu)

    // Help menu
    helpMenu := menu.AddRole(application.HelpMenu)
    helpMenu.Add("Documentation").OnClick(openDocumentation)
    
    // About location varies by platform
    if runtime.GOOS == "darwin" {
        // On macOS, About is in Application menu (added by AppMenu role)
    } else {
        helpMenu.AddSeparator()
        helpMenu.Add("About").OnClick(showAbout)
    }

    // Set the application menu
    app.Menu.Set(menu)
}

func handleImport(ctx *application.Context) {
    // Implementation
}

func handleExport(ctx *application.Context) {
    // Implementation
}

func toggleFullscreen(ctx *application.Context) {
    window := application.Get().Window.Current()
    window.ToggleFullscreen()
}

func toggleSidebar(ctx *application.Context) {
    // Implementation
}

func toggleToolbar(ctx *application.Context) {
    // Implementation
}

func showSettings(ctx *application.Context) {
    // Implementation
}

func toggleDarkMode(ctx *application.Context) {
    isDark := ctx.ClickedMenuItem().Checked()
    // Apply theme
}

func openDocumentation(ctx *application.Context) {
    // Open browser
}

func showAbout(ctx *application.Context) {
    // Show about dialog
}
```

## Práticas recomendadas

### ✅ Faça

- **Use funções de menu** para menus padrão (Arquivo, Editar etc.)
- **Siga as convenções da plataforma** para a estrutura dos menus
- **Adicione atalhos de teclado** às ações comuns
- **Chame menu.Update()** após alterar o estado do menu
- **Teste em todas as plataformas** — o comportamento varia
- **Mantenha os menus pouco profundos** — no máximo 2-3 níveis
- **Use rótulos claros** — "Salvar projeto", não "Salvar"

### ❌ Não faça

- **Não codifique diretamente atalhos específicos da plataforma** — use `CmdOrCtrl`
- **Não coloque Sair no menu Arquivo no macOS** — essa opção fica no menu do aplicativo
- **Não coloque Sobre no menu Ajuda no macOS** — essa opção fica no menu do aplicativo
- **Não se esqueça de chamar menu.Update()** — os menus não funcionarão corretamente
- **Não crie muitos níveis de aninhamento** — os usuários se perdem
- **Não use jargão** — use rótulos fáceis de entender

## Próximas etapas

@cards{cols="2"}
📖 Referência de menus
Referência completa dos tipos e das propriedades dos itens de menu.

[Saiba mais →](/features/menus/reference/)

---
◆ Menus de contexto
Crie menus de contexto acionados com o botão direito do mouse.

[Saiba mais →](/features/menus/context/)

---
★ Menus da área de notificação
Adicione integração com a área de notificação ou a barra de menus.

[Saiba mais →](/features/menus/systray/)

---
📖 Padrões de menus
Padrões comuns de menus e práticas recomendadas.

[Saiba mais →](/guides/menus/)

@end

---

**Tem dúvidas?** Pergunte no [Discord](https://discord.gg/JDdSxwjhGf) ou consulte o [exemplo de menu](https://github.com/wailsapp/wails/tree/master/v3/examples/menu).
