---
title: "API de menus"
description: "Referência completa da API de menus"
slug: "reference/menu"
sourcePath: "reference/menu.md"
---

## Visão geral

A API de menus fornece métodos para criar e gerenciar menus do aplicativo, menus de contexto e menus da bandeja do sistema.

**Tipos de menu:**

- **Menus do aplicativo** - Barra de menus superior (Arquivo, Editar etc.)
- **Menus de contexto** - Menus abertos com o botão direito do mouse
- **Menus da bandeja do sistema** - Menus na bandeja do sistema/área de notificação

## Criação de menus

### NewMenu()

Cria um novo menu.

```go
func (a *App) NewMenu() *Menu
```

**Exemplo:**

```go
menu := app.NewMenu()
```

## Métodos de menu

### Add()

Adiciona um item ao menu.

```go
func (m *Menu) Add(label string) *MenuItem
```

**Parâmetros:**

- `label` - O texto exibido no item de menu

**Retorna:** O item de menu criado

**Exemplo:**

```go
item := menu.Add("Open File")
item.OnClick(func(ctx *application.Context) {
    // Handle click
})
```

### AddSubmenu()

Adiciona um submenu ao menu.

```go
func (m *Menu) AddSubmenu(label string) *Menu
```

**Parâmetros:**

- `label` - O rótulo do submenu

**Retorna:** O submenu criado

**Exemplo:**

```go
fileMenu := menu.AddSubmenu("File")
fileMenu.Add("New")
fileMenu.Add("Open")
fileMenu.Add("Save")
```

### AddSeparator()

Adiciona uma linha separadora visual entre os itens de menu.

```go
func (m *Menu) AddSeparator()
```

**Exemplo:**

```go
menu.Add("Copy")
menu.Add("Paste")
menu.AddSeparator()
menu.Add("Select All")
```

**Prática recomendada:** Use separadores para agrupar itens de menu relacionados.

### AddCheckbox()

Adiciona um item de menu que pode ser marcado.

```go
func (m *Menu) AddCheckbox(label string, checked bool) *MenuItem
```

**Parâmetros:**

- `label` - O rótulo da caixa de seleção
- `checked` - O estado inicial de seleção

**Exemplo:**

```go
darkMode := menu.AddCheckbox("Dark Mode", false)
darkMode.OnClick(func(ctx *application.Context) {
    isChecked := darkMode.Checked()
    // Toggle dark mode
})
```

### AddRadio()

Adiciona um item de menu do tipo botão de opção (grupo mutuamente exclusivo).

```go
func (m *Menu) AddRadio(label string, checked bool) *MenuItem
```

**Parâmetros:**

- `label` - O rótulo do botão de opção
- `checked` - O estado inicial de seleção

**Exemplo:**

```go
// Create radio group for view modes
viewMenu := menu.AddSubmenu("View")
listView := viewMenu.AddRadio("List View", true)
gridView := viewMenu.AddRadio("Grid View", false)
treeView := viewMenu.AddRadio("Tree View", false)

listView.OnClick(func(ctx *application.Context) {
    setViewMode("list")
})
gridView.OnClick(func(ctx *application.Context) {
    setViewMode("grid")
})
```

### Update()

Atualiza o menu para refletir todas as alterações feitas nos itens de menu.

```go
func (m *Menu) Update()
```

**Exemplo:**

```go
item.SetEnabled(false)
menu.Update()  // Must call to apply changes
```

**Importante:** Sempre chame `Update()` depois de modificar as propriedades de um item de menu.

## Métodos de itens de menu

### OnClick()

Registra um manipulador de clique para o item de menu.

```go
func (mi *MenuItem) OnClick(callback func(ctx *application.Context)) *MenuItem
```

**Parâmetros:**

- `callback` - Função chamada quando o item é clicado

**Retorna:** O item de menu (para encadeamento)

**Exemplo:**

```go
item.OnClick(func(ctx *application.Context) {
    fmt.Println("Menu item clicked")
    app.Logger.Info("User clicked menu item")
})
```

### SetLabel()

Altera o rótulo do item de menu.

```go
func (mi *MenuItem) SetLabel(label string) *MenuItem
```

**Exemplo:**

```go
item.SetLabel("Save As...")
menu.Update()
```

### SetEnabled()

Habilita ou desabilita o item de menu.

```go
func (mi *MenuItem) SetEnabled(enabled bool) *MenuItem
```

**Exemplo:**

```go
// Disable save when no document is open
saveItem.SetEnabled(hasOpenDocument)
menu.Update()
```

**Padrão comum:**

```go
// Update menu state based on application state
func updateMenuState() {
    saveItem.SetEnabled(hasUnsavedChanges)
    undoItem.SetEnabled(canUndo)
    redoItem.SetEnabled(canRedo)
    menu.Update()
}
```

### SetChecked()

Define o estado de seleção de itens de menu do tipo caixa de seleção ou botão de opção.

```go
func (mi *MenuItem) SetChecked(checked bool) *MenuItem
```

**Exemplo:**

```go
darkModeItem.SetChecked(isDarkModeEnabled)
menu.Update()
```

### Checked()

Retorna o estado de seleção atual.

```go
func (mi *MenuItem) Checked() bool
```

**Exemplo:**

```go
if darkModeItem.Checked() {
    // Dark mode is enabled
}
```

### SetAccelerator()

Define um atalho de teclado para o item de menu.

```go
func (mi *MenuItem) SetAccelerator(accelerator string) *MenuItem
```

**Parâmetros:**

- `accelerator` - Atalho de teclado (por exemplo, "Ctrl+S", "Cmd+Q")

**Formato do acelerador:**

- **Modificadores:** `Ctrl`, `Cmd`, `Alt`, `Shift`
- **Teclas:** `A-Z`, `0-9`, `F1-F12`, `Enter`, `Backspace` etc.
- **Plataforma:** Use `Cmd` no macOS e `Ctrl` no Windows/Linux

**Exemplo:**

```go
saveItem.SetAccelerator("Ctrl+S")
quitItem.SetAccelerator("Ctrl+Q")
newItem.SetAccelerator("Ctrl+N")
```

**Exemplo adaptado à plataforma:**

```go
import "runtime"

var quitShortcut string
if runtime.GOOS == "darwin" {
    quitShortcut = "Cmd+Q"
} else {
    quitShortcut = "Ctrl+Q"
}
quitItem.SetAccelerator(quitShortcut)
```

### SetTooltip()

Define uma dica de ferramenta exibida ao passar o cursor sobre o item de menu.

```go
func (mi *MenuItem) SetTooltip(tooltip string) *MenuItem
```

**Exemplo:**

```go
item.SetTooltip("Opens a file from disk")
```

### SetHidden()

Exibe ou oculta o item de menu.

```go
func (mi *MenuItem) SetHidden(hidden bool) *MenuItem
```

**Exemplo:**

```go
// Hide debug menu in production
debugItem.SetHidden(!isDevelopment)
menu.Update()
```

## Menu do aplicativo

### app.Menu.Set()

Define a barra de menus principal do aplicativo.

```go
func (mm *MenuManager) Set(menu *Menu)
```

**Exemplo:**

```go
menu := app.NewMenu()

// File menu
fileMenu := menu.AddSubmenu("File")
fileMenu.Add("New").SetAccelerator("Ctrl+N").OnClick(newFile)
fileMenu.Add("Open").SetAccelerator("Ctrl+O").OnClick(openFile)
fileMenu.Add("Save").SetAccelerator("Ctrl+S").OnClick(saveFile)
fileMenu.AddSeparator()
fileMenu.Add("Exit").SetAccelerator("Ctrl+Q").OnClick(func(ctx *application.Context) {
    app.Quit()
})

// Edit menu
editMenu := menu.AddSubmenu("Edit")
editMenu.Add("Undo").SetAccelerator("Ctrl+Z").OnClick(undo)
editMenu.Add("Redo").SetAccelerator("Ctrl+Y").OnClick(redo)
editMenu.AddSeparator()
editMenu.Add("Cut").SetAccelerator("Ctrl+X").OnClick(cut)
editMenu.Add("Copy").SetAccelerator("Ctrl+C").OnClick(copy)
editMenu.Add("Paste").SetAccelerator("Ctrl+V").OnClick(paste)

app.Menu.Set(menu)
```

**Observações sobre as plataformas:**

- **macOS:** o menu aparece na barra de menus superior
- **Windows/Linux:** o menu aparece na barra de título da janela
- **macOS:** adiciona automaticamente o menu do aplicativo com o nome do aplicativo

## Menus de contexto

### app.ContextMenu.New() / app.ContextMenu.Add()

Crie um `*ContextMenu` por meio do gerenciador e registre-o com um nome. `ContextMenuManager.Add` recebe `*ContextMenu` — **não** `*Menu` — e **não** há nenhum método `app.RegisterContextMenu`.

```go
func (cm *ContextMenuManager) New() *ContextMenu
func (cm *ContextMenuManager) Add(name string, menu *ContextMenu)
func (cm *ContextMenuManager) Get(name string) (*ContextMenu, bool)
func (cm *ContextMenuManager) Remove(name string)
```

Como alternativa, `application.NewContextMenu(name string) *ContextMenu`, no nível do pacote, cria e registra um menu de contexto em uma única etapa.

**Go:**

```go
// Build the context menu via the manager
contextMenu := app.ContextMenu.New()
contextMenu.Add("Cut").OnClick(cut)
contextMenu.Add("Copy").OnClick(copy)
contextMenu.Add("Paste").OnClick(paste)
contextMenu.AddSeparator()
contextMenu.Add("Select All").OnClick(selectAll)

// Register under a name; HTML opts into it via the CSS custom property below.
app.ContextMenu.Add("editor", contextMenu)
```

**HTML/CSS:**

O runtime aciona um menu de contexto registrado quando o elemento clicado com o botão direito (ou qualquer ancestral dele) tem a propriedade personalizada CSS `--custom-contextmenu` definida com esse nome. A propriedade opcional `--custom-contextmenu-data` é passada ao callback do Go por meio de `ctx.ContextMenuData()`. Para suprimir o menu de contexto padrão do navegador, defina `--default-contextmenu: hide` (ou `auto`/`show`).

```html
<!-- Trigger context menu on right-click -->
<div style="--custom-contextmenu: editor; --default-contextmenu: hide">
    Right-click here for context menu
</div>
```

Não há nenhum atributo `data-wails-context-menu="..."` — ele nunca foi integrado ao runtime.

**Menus de contexto dinâmicos:**

```go
// Update context menu based on selection
func updateContextMenu() {
    contextMenu := app.ContextMenu.New()

    if hasSelection {
        contextMenu.Add("Cut").OnClick(cut)
        contextMenu.Add("Copy").OnClick(copy)
    }

    contextMenu.Add("Paste").SetEnabled(hasClipboardContent).OnClick(paste)

    app.ContextMenu.Add("editor", contextMenu)
}
```

## Menu da bandeja do sistema

### app.SystemTray.New()

Cria um novo ícone na bandeja do sistema.

```go
func (sm *SystemTrayManager) New() *SystemTray
```

**Exemplo:**

```go
tray := app.SystemTray.New()
```

### SetIcon()

Define o ícone da bandeja do sistema.

```go
func (st *SystemTray) SetIcon(icon []byte) *SystemTray
```

**Exemplo:**

```go
iconData, _ := os.ReadFile("icon.png")
tray.SetIcon(iconData)
```

### SetMenu()

Define o menu da bandeja do sistema.

```go
func (st *SystemTray) SetMenu(menu *Menu) *SystemTray
```

**Exemplo:**

```go
trayMenu := app.NewMenu()
trayMenu.Add("Show Window").OnClick(func(ctx *application.Context) {
    window.Show()
    window.Focus()
})
trayMenu.Add("Settings").OnClick(openSettings)
trayMenu.AddSeparator()
trayMenu.Add("Quit").OnClick(func(ctx *application.Context) {
    app.Quit()
})

tray.SetMenu(trayMenu)
```

### SetTooltip()

Define a dica de ferramenta exibida ao passar o cursor sobre o ícone da bandeja. Não retorna nenhum valor.

```go
func (st *SystemTray) SetTooltip(tooltip string)
```

**Exemplo:**

```go
tray.SetTooltip("My Application - Running")
```

### OnClick()

Trata o clique com o botão esquerdo no ícone da bandeja.

```go
func (st *SystemTray) OnClick(callback func()) *SystemTray
```

**Exemplo:**

```go
tray.OnClick(func() {
    if window.IsVisible() {
        window.Hide()
    } else {
        window.Show()
        window.Focus()
    }
})
```

## Exemplos completos

### Menu padrão do aplicativo

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
)

func createMenu(app *application.App) *application.Menu {
    menu := app.NewMenu()

    // File menu
    fileMenu := menu.AddSubmenu("File")
    fileMenu.Add("New").
        SetAccelerator("Ctrl+N").
        OnClick(func(ctx *application.Context) {
            // Create new document
        })
    fileMenu.Add("Open").
        SetAccelerator("Ctrl+O").
        OnClick(func(ctx *application.Context) {
            // Open file dialog
        })
    fileMenu.Add("Save").
        SetAccelerator("Ctrl+S").
        OnClick(func(ctx *application.Context) {
            // Save document
        })
    fileMenu.AddSeparator()
    fileMenu.Add("Exit").OnClick(func(ctx *application.Context) {
        app.Quit()
    })

    // Edit menu
    editMenu := menu.AddSubmenu("Edit")
    editMenu.Add("Undo").SetAccelerator("Ctrl+Z")
    editMenu.Add("Redo").SetAccelerator("Ctrl+Y")
    editMenu.AddSeparator()
    editMenu.Add("Cut").SetAccelerator("Ctrl+X")
    editMenu.Add("Copy").SetAccelerator("Ctrl+C")
    editMenu.Add("Paste").SetAccelerator("Ctrl+V")

    // View menu
    viewMenu := menu.AddSubmenu("View")
    darkMode := viewMenu.AddCheckbox("Dark Mode", false)
    darkMode.OnClick(func(ctx *application.Context) {
        // Toggle dark mode
        isChecked := darkMode.Checked()
        app.Logger.Info("Dark mode", "enabled", isChecked)
    })
    viewMenu.AddSeparator()
    viewMenu.AddRadio("List View", true)
    viewMenu.AddRadio("Grid View", false)
    viewMenu.AddRadio("Detail View", false)

    // Help menu
    helpMenu := menu.AddSubmenu("Help")
    helpMenu.Add("Documentation").OnClick(func(ctx *application.Context) {
        // Open docs
    })
    helpMenu.Add("About").OnClick(func(ctx *application.Context) {
        // Show about dialog
    })

    return menu
}

func main() {
    app := application.New(application.Options{
        Name: "Menu Demo",
    })

    menu := createMenu(app)
    app.Menu.Set(menu)

    window := app.Window.New()
    window.Show()

    app.Run()
}
```

### Aplicativo com bandeja do sistema

```go
func setupSystemTray(app *application.App, window application.Window) {
    // Create system tray
    tray := app.SystemTray.New()

    // Set icon
    iconData, _ := os.ReadFile("icon.png")
    tray.SetIcon(iconData)
    tray.SetTooltip("My App - Running")

    // Handle left-click on tray icon
    tray.OnClick(func() {
        if window.IsVisible() {
            window.Hide()
        } else {
            window.Show()
            window.Focus()
        }
    })

    // Create tray menu
    trayMenu := app.NewMenu()

    showItem := trayMenu.Add("Show Window")
    showItem.OnClick(func(ctx *application.Context) {
        window.Show()
        window.Focus()
    })

    trayMenu.AddSeparator()

    trayMenu.Add("Settings").OnClick(func(ctx *application.Context) {
        // Open settings window
    })

    trayMenu.Add("About").OnClick(func(ctx *application.Context) {
        // Show about dialog
    })

    trayMenu.AddSeparator()

    trayMenu.Add("Quit").OnClick(func(ctx *application.Context) {
        app.Quit()
    })

    tray.SetMenu(trayMenu)
}
```

### Atualizações dinâmicas de menus

```go
type Editor struct {
    app         *application.App
    menu        *application.Menu
    undoItem    *application.MenuItem
    redoItem    *application.MenuItem
    saveItem    *application.MenuItem
    undoStack   []string
    redoStack   []string
    hasChanges  bool
}

func (e *Editor) createMenu() {
    e.menu = e.app.NewMenu()

    fileMenu := e.menu.AddSubmenu("File")
    e.saveItem = fileMenu.Add("Save").SetAccelerator("Ctrl+S")
    e.saveItem.OnClick(func(ctx *application.Context) {
        e.save()
    })

    editMenu := e.menu.AddSubmenu("Edit")
    e.undoItem = editMenu.Add("Undo").SetAccelerator("Ctrl+Z")
    e.undoItem.OnClick(func(ctx *application.Context) {
        e.undo()
    })

    e.redoItem = editMenu.Add("Redo").SetAccelerator("Ctrl+Y")
    e.redoItem.OnClick(func(ctx *application.Context) {
        e.redo()
    })

    e.updateMenuState()
    e.app.Menu.Set(e.menu)
}

func (e *Editor) updateMenuState() {
    // Update menu items based on current state
    e.saveItem.SetEnabled(e.hasChanges)
    e.undoItem.SetEnabled(len(e.undoStack) > 0)
    e.redoItem.SetEnabled(len(e.redoStack) > 0)
    e.menu.Update()
}

func (e *Editor) onChange() {
    e.hasChanges = true
    e.updateMenuState()
}

func (e *Editor) save() {
    // Save logic
    e.hasChanges = false
    e.updateMenuState()
}
```

## Boas práticas

### ✅ Faça

- **Use aceleradores padrão** — siga as convenções da plataforma (Ctrl+C para copiar etc.)
- **Chame Update() após as alterações** — caso contrário, o menu não refletirá as alterações
- **Agrupe itens relacionados** — use separadores para organizar os itens de menu
- **Desabilite ações indisponíveis** — não as oculte; desabilite-as com SetEnabled(false)
- **Use rótulos claros** — seja conciso e descritivo
- **Siga as convenções da plataforma** — os padrões de menu diferem entre macOS e Windows/Linux

### ❌ Não faça

- **Não se esqueça de Update()** — esse é o erro mais comum
- **Não crie níveis demais de aninhamento** — mantenha os menus com no máximo 2-3 níveis
- **Não use rótulos ambíguos** — "Processar" versus "Processar documento"
- **Não complique demais** — mantenha os menus simples e focados
- **Não misture metáforas** — mantenha a nomenclatura e a organização consistentes

## Observações específicas de cada plataforma

### macOS

- O menu do aplicativo é adicionado automaticamente com o nome do aplicativo
- Use `Cmd` em vez de `Ctrl` nos aceleradores
- Por padrão, "Sobre", "Preferências" e "Sair" ficam no menu do aplicativo

### Windows/Linux

- Sem menu de aplicativo automático
- Use `Ctrl` para os atalhos de teclado
- “Sair” geralmente fica no menu Arquivo
