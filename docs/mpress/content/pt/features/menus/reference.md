---
title: "Referência de menus"
description: "Referência completa dos tipos, propriedades e métodos de itens de menu"
slug: "features/menus/reference"
sourcePath: "features/menus/reference.md"
---

## Referência de menus

Referência completa dos tipos, propriedades e comportamentos dinâmicos de itens de menu. Crie menus profissionais e responsivos com caixas de seleção, grupos de opções, separadores e atualizações dinâmicas.

## Tipos de itens de menu

### Itens de menu comuns

O tipo mais comum — exibe texto e aciona uma ação:

```go
menuItem := menu.Add("Click Me")
menuItem.OnClick(func(ctx *application.Context) {
    fmt.Println("Menu item clicked!")
})
```

**Use para:** comandos, ações e abertura de janelas

### Caixas de seleção

Itens de menu alternáveis, com estado marcado ou desmarcado:

```go
checkbox := menu.AddCheckbox("Enable Feature", true)  // true = initially checked
checkbox.OnClick(func(ctx *application.Context) {
    isChecked := ctx.ClickedMenuItem().Checked()
    fmt.Printf("Feature is now: %v\n", isChecked)
})
```

**Use para:** configurações booleanas, ativação e desativação de recursos e opções de visualização

**Importante:** o estado marcado é alternado automaticamente ao clicar.

### Grupos de opções

Opções mutuamente exclusivas — apenas uma pode ser selecionada:

```go
menu.AddRadio("Small", true)   // true = initially selected
menu.AddRadio("Medium", false)
menu.AddRadio("Large", false)
```

**Use para:** escolhas mutuamente exclusivas (tamanho, tema, modo)

**Como funciona o agrupamento:**

- Itens de opção adjacentes formam um grupo automaticamente
- Selecionar um item desmarca os demais no grupo
- Separe os grupos com um separador ou um item comum

**Exemplo com vários grupos:**

```go
// Group 1: Size
menu.AddRadio("Small", true)
menu.AddRadio("Medium", false)
menu.AddRadio("Large", false)

menu.AddSeparator()

// Group 2: Theme
menu.AddRadio("Light", true)
menu.AddRadio("Dark", false)
```

### Submenus

Estruturas de menu aninhadas para organização:

```go
submenu := menu.AddSubmenu("More Options")
submenu.Add("Submenu Item 1").OnClick(func(ctx *application.Context) {
    // Handle click
})
submenu.Add("Submenu Item 2")
```

**Use para:** agrupar itens relacionados e reduzir a poluição visual

**Limite de aninhamento:** a maioria das plataformas oferece suporte a 2-3 níveis. Evite um aninhamento mais profundo.

### Separadores

Divisores visuais entre itens de menu:

```go
menu.Add("Item 1")
menu.AddSeparator()
menu.Add("Item 2")
```

**Use para:** agrupar visualmente itens relacionados

**Boa prática:** não inicie nem termine menus com separadores.

## Propriedades dos itens de menu

### Rótulo

O texto exibido no item de menu:

```go
menuItem := menu.Add("Initial Label")
menuItem.SetLabel("New Label")

// Get current label
label := menuItem.Label()
```

**Rótulos dinâmicos:**

```go
updateMenuItem := menu.Add("Check for Updates")
updateMenuItem.OnClick(func(ctx *application.Context) {
    updateMenuItem.SetLabel("Checking...")
    menu.Update()  // Important on Windows!
    
    // Perform update check
    checkForUpdates()
    
    updateMenuItem.SetLabel("Check for Updates")
    menu.Update()
})
```

### Estado habilitado

Controle se o usuário pode interagir com o item de menu:

```go
menuItem := menu.Add("Save")
menuItem.SetEnabled(false)  // Greyed out, can't click

// Enable it later
menuItem.SetEnabled(true)
menu.Update()  // Important: Call this after changing enabled state!

// Check current state
isEnabled := menuItem.Enabled()
```

@note{type="caution" title="Comportamento dos menus no Windows"}
No Windows, os menus precisam ser reconstruídos quando seu estado muda. **Sempre chame `menu.Update()` após habilitar ou desabilitar itens de menu**, especialmente se o item tiver sido criado enquanto estava desabilitado.

**Motivo:** no Windows, os menus são recriados do zero quando atualizados. Se você não chamar `Update()`, os manipuladores de clique não serão acionados corretamente.

@end

**Exemplo: habilitação e desabilitação dinâmicas**

```go
var hasSelection bool

cutMenuItem := menu.Add("Cut")
cutMenuItem.SetEnabled(false)  // Initially disabled

copyMenuItem := menu.Add("Copy")
copyMenuItem.SetEnabled(false)

// When selection changes
func onSelectionChanged(selected bool) {
    hasSelection = selected
    cutMenuItem.SetEnabled(hasSelection)
    copyMenuItem.SetEnabled(hasSelection)
    menu.Update()  // Critical on Windows!
}
```

**Padrão comum: habilitar sob determinada condição**

```go
saveMenuItem := menu.Add("Save")

func updateSaveMenuItem() {
    canSave := hasUnsavedChanges() && !isSaving()
    saveMenuItem.SetEnabled(canSave)
    menu.Update()
}

// Call whenever state changes
onDocumentChanged(func() {
    updateSaveMenuItem()
})
```

### Estado marcado

Para caixas de seleção e itens de opção, controle ou consulte o estado marcado:

```go
checkbox := menu.AddCheckbox("Feature", false)
checkbox.SetChecked(true)
menu.Update()

// Query state
isChecked := checkbox.Checked()
```

**Alternância automática:** as caixas de seleção alternam automaticamente ao serem clicadas. Não é necessário chamar `SetChecked()` no manipulador de clique.

**Controle manual:**

```go
checkbox := menu.AddCheckbox("Auto-save", false)

// Sync with external state
func syncAutoSave(enabled bool) {
    checkbox.SetChecked(enabled)
    menu.Update()
}
```

### Aceleradores (atalhos de teclado)

Adicione atalhos de teclado aos itens de menu:

```go
saveMenuItem := menu.Add("Save")
saveMenuItem.SetAccelerator("CmdOrCtrl+S")

quitMenuItem := menu.Add("Quit")
quitMenuItem.SetAccelerator("CmdOrCtrl+Q")
```

**Formato dos aceleradores:**

- `CmdOrCtrl` — Cmd no macOS, Ctrl no Windows/Linux
- `Shift`, `Alt`, `Option` — teclas modificadoras
- `A-Z`, `0-9` — teclas de letras/números
- `F1-F12` — teclas de função
- `Enter`, `Space`, `Backspace` etc. — teclas especiais

**Exemplos:**

```go
"CmdOrCtrl+S"           // Save
"CmdOrCtrl+Shift+S"     // Save As
"CmdOrCtrl+W"           // Close Window
"CmdOrCtrl+Q"           // Quit
"F5"                    // Refresh
"CmdOrCtrl+,"           // Preferences (macOS convention)
"Alt+F4"                // Close (Windows convention)
```

**Aceleradores específicos da plataforma:**

```go
if runtime.GOOS == "darwin" {
    prefsMenuItem.SetAccelerator("Cmd+,")
} else {
    prefsMenuItem.SetAccelerator("Ctrl+P")
}
```

### Dica de ferramenta

Adicione aos itens de menu um texto exibido ao passar o cursor (o suporte varia conforme a plataforma):

```go
menuItem := menu.Add("Advanced Options")
menuItem.SetTooltip("Configure advanced settings")
```

**Suporte por plataforma:**

- **Windows:** ✅ compatível
- **macOS:** ❌ incompatível (dicas de ferramenta não são padrão em menus)
- **Linux:** ⚠️ varia conforme o ambiente de área de trabalho

### Estado oculto

Oculte itens de menu sem removê-los:

```go
debugMenuItem := menu.Add("Debug Mode")
debugMenuItem.SetHidden(true)  // Hidden

// Show in debug builds
if isDebugBuild {
    debugMenuItem.SetHidden(false)
    menu.Update()
}
```

**Use para:** opções de depuração, sinalizadores de recursos e recursos condicionais

## Tratamento de eventos

### Manipulador OnClick

Execute código quando o item de menu for clicado:

```go
menuItem := menu.Add("Click Me")
menuItem.OnClick(func(ctx *application.Context) {
    // Handle click
    fmt.Println("Clicked!")
})
```

**O contexto fornece:**

- `ctx.ClickedMenuItem()` — O item de menu que foi clicado
- Contexto da janela (se for proveniente do menu da janela)
- Contexto do aplicativo

**Exemplo: acessar o item de menu no manipulador**

```go
checkbox := menu.AddCheckbox("Feature", false)
checkbox.OnClick(func(ctx *application.Context) {
    item := ctx.ClickedMenuItem()
    isChecked := item.Checked()
    fmt.Printf("Feature is now: %v\n", isChecked)
})
```

### Vários manipuladores

Você pode definir vários manipuladores (o último prevalece):

```go
menuItem := menu.Add("Action")
menuItem.OnClick(func(ctx *application.Context) {
    fmt.Println("First handler")
})

// This replaces the first handler
menuItem.OnClick(func(ctx *application.Context) {
    fmt.Println("Second handler - this one runs")
})
```

**Prática recomendada:** defina o manipulador uma única vez e, se necessário, use lógica condicional dentro dele.

## Menus dinâmicos

### Atualização de itens de menu

**Regra de ouro:** sempre chame `menu.Update()` depois de alterar o estado do menu.

```go
// ✅ Correct
menuItem.SetEnabled(true)
menu.Update()

// ❌ Wrong (especially on Windows)
menuItem.SetEnabled(true)
// Forgot to call Update() - click handlers may not work!
```

**Por que isso é importante:**

- **Windows:** os menus são reconstruídos quando atualizados
- **macOS/Linux:** é menos importante, mas ainda recomendado
- **Manipuladores de clique:** não serão acionados corretamente sem Update()

### Reconstrução de menus

Para alterações significativas, reconstrua o menu inteiro:

```go
func rebuildFileMenu() {
    menu := app.Menu.New()
    
    menu.Add("New").OnClick(handleNew)
    menu.Add("Open").OnClick(handleOpen)
    
    if hasRecentFiles() {
        recentMenu := menu.AddSubmenu("Open Recent")
        for _, file := range getRecentFiles() {
            recentMenu.Add(file).OnClick(func(ctx *application.Context) {
                openFile(file)
            })
        }
    }
    
    menu.AddSeparator()
    menu.Add("Quit").OnClick(handleQuit)
    
    // Set the new menu
    window.SetMenu(menu)
}
```

**Quando reconstruir:**

- A lista de arquivos recentes muda
- Os menus de plugins mudam
- Transições de estado significativas

**Quando atualizar:**

- Habilitar ou desabilitar itens
- Alterar rótulos
- Marcar ou desmarcar caixas de seleção

### Menus sensíveis ao contexto

Ajuste os menus com base no estado do aplicativo:

```go
func updateEditMenu() {
    cutMenuItem.SetEnabled(hasSelection())
    copyMenuItem.SetEnabled(hasSelection())
    pasteMenuItem.SetEnabled(hasClipboardContent())
    undoMenuItem.SetEnabled(canUndo())
    redoMenuItem.SetEnabled(canRedo())
    menu.Update()
}

// Call whenever state changes
onSelectionChanged(updateEditMenu)
onClipboardChanged(updateEditMenu)
onUndoStackChanged(updateEditMenu)
```

## Diferenças entre plataformas

### Localização da barra de menus

| Plataforma | Localização | Observações |
| --- | --- | --- |
| **macOS** | Parte superior da tela | Barra de menus global |
| **Windows** | Parte superior da janela | Menu por janela |
| **Linux** | Parte superior da janela | Por janela (geralmente) |

### Menus padrão

**macOS:**

- Tem o menu "Aplicativo" (com o nome do aplicativo)
- "Preferências" no menu Aplicativo
- "Encerrar" no menu Aplicativo

**Windows/Linux:**

- Não há menu Aplicativo
- "Preferências" no menu Editar ou Ferramentas
- "Sair" no menu Arquivo

**Exemplo: estrutura apropriada para a plataforma**

```go
menu := app.Menu.New()

// macOS gets Application menu
if runtime.GOOS == "darwin" {
    menu.AddRole(application.AppMenu)
}

// File menu
fileMenu := menu.AddSubmenu("File")
fileMenu.Add("New")
fileMenu.Add("Open")

// Preferences location varies
if runtime.GOOS == "darwin" {
    // On macOS, preferences are in Application menu (added by AppMenu role)
} else {
    // On Windows/Linux, add to Edit or Tools menu
    editMenu := menu.AddSubmenu("Edit")
    editMenu.Add("Preferences")
}
```

### Convenções de atalhos de teclado

**macOS:**

- `Cmd+` para a maioria dos atalhos
- `Cmd+,` para Preferências
- `Cmd+Q` para Encerrar

**Windows:**

- `Ctrl+` para a maioria dos atalhos
- `Ctrl+P` ou `Ctrl+,` para Preferências
- `Alt+F4` para Sair (ou `Ctrl+Q`)

**Linux:**

- Geralmente segue as convenções do Windows
- O ambiente de desktop pode substituí-las

## Práticas recomendadas

### ✅ Faça

- **Chame menu.Update()** depois de alterar o estado do menu (especialmente no Windows)
- **Use grupos de botões de opção** para opções mutuamente exclusivas
- **Use caixas de seleção** para recursos que podem ser ativados ou desativados
- **Adicione atalhos de teclado** às ações comuns
- **Agrupe itens relacionados** usando separadores
- **Teste em todas as plataformas** — o comportamento varia

### ❌ Não faça

- **Não se esqueça de menu.Update()** — os manipuladores de clique não funcionarão corretamente
- **Não crie muitos níveis de aninhamento** — no máximo 2-3 níveis
- **Não comece nem termine com separadores** — isso parece pouco profissional
- **Não use dicas de ferramenta no macOS** — esse recurso não é compatível
- **Não fixe atalhos específicos da plataforma no código** — use `CmdOrCtrl`

## Solução de problemas

### Os itens de menu não respondem

**Sintoma:** os manipuladores de clique não são acionados

**Causa:** a chamada a `menu.Update()` não foi feita após habilitar o item

**Solução:**

```go
menuItem.SetEnabled(true)
menu.Update()  // Add this!
```

### Os itens de menu aparecem esmaecidos

**Sintoma:** não é possível clicar nos itens de menu

**Causa:** os itens estão desabilitados

**Solução:**

```go
menuItem.SetEnabled(true)
menu.Update()
```

### Os atalhos de teclado não funcionam

**Sintoma:** os atalhos de teclado não acionam os itens de menu

**Causas:**

1. Formato incorreto do atalho de teclado
2. Conflito com atalhos do sistema
3. A janela não está com o foco

**Solução:**

```go
// Check format
menuItem.SetAccelerator("CmdOrCtrl+S")  // ✅ Correct
menuItem.SetAccelerator("Ctrl+S")       // ❌ Wrong (macOS uses Cmd)

// Avoid conflicts
// ❌ Cmd+H (Hide Window on macOS - system shortcut)
// ✅ Cmd+Shift+H (Custom shortcut)
```

## Próximas etapas

- [Menus do aplicativo](/features/menus/application/) — crie barras de menus para o aplicativo
- [Menus de contexto](/features/menus/context/) — menus de contexto abertos com o botão direito
- [Menus da bandeja do sistema](/features/menus/systray/) — menus da bandeja do sistema ou da barra de menus
- [Padrões de menus](/guides/menus/) — padrões comuns de menus e práticas recomendadas

---

**Dúvidas?** Pergunte no [Discord](https://discord.gg/JDdSxwjhGf) ou consulte os [exemplos de menus](https://github.com/wailsapp/wails/tree/master/v3/examples/menu).
