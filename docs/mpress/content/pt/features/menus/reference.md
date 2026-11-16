---
title: "Referência de menus"
description: "Referência completa de tipos, propriedades e métodos de itens de menu"
slug: "features/menus/reference"
sourcePath: "features/menus/reference.md"
---

## Referência de menus

Referência completa de tipos, propriedades e comportamento dinâmico dos itens de menu. Crie menus profissionais e responsivos com caixas de seleção, grupos de opções, separadores e atualizações dinâmicas.

## Tipos de itens de menu

### Itens de menu comuns

O tipo mais comum: exibe texto e aciona uma ação:

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

**Importante:** o estado de seleção é alternado automaticamente ao clicar.

### Grupos de opções

Opções mutuamente exclusivas: somente uma pode ser selecionada:

```go
menu.AddRadio("Small", true)   // true = initially selected
menu.AddRadio("Medium", false)
menu.AddRadio("Large", false)
```

**Use para:** escolhas mutuamente exclusivas (tamanho, tema, modo)

**Como funciona o agrupamento:**

- Itens de opção adjacentes formam um grupo automaticamente
- Selecionar um item desmarca os demais itens do grupo
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

**Limite de aninhamento:** a maioria das plataformas aceita 2-3 níveis. Evite um aninhamento mais profundo.

### Separadores

Divisores visuais entre itens de menu:

```go
menu.Add("Item 1")
menu.AddSeparator()
menu.Add("Item 2")
```

**Use para:** agrupar visualmente itens relacionados

**Boa prática:** não comece nem termine menus com separadores.

### Cabeçalhos de seção

Rótulos nativos e não interativos para um grupo de itens de menu no **macOS 14+**:

```go
menu.AddSectionHeader("Build server")
menu.Add("Start build").OnClick(startBuild)
menu.Add("View logs").OnClick(viewLogs)
```

`AddSectionHeader` usa a aparência padrão de cabeçalho de seção do AppKit. Não retorna nenhum item de menu nem associa uma ação. Em versões anteriores do macOS, no Windows, no Linux e em compilações para dispositivos móveis e servidores, a chamada não tem efeito e não insere nenhuma linha.

Use-o em menus da bandeja, de contexto ou do aplicativo para rotular comandos relacionados. Use `AddSeparator()` para separar grupos de opções em todas as plataformas, pois cabeçalhos omitidos não separam grupos de opções. Consulte a [API de menus](/reference/menu#addsectionheader) para ver a assinatura do método e o exemplo de `systray-menu` para ver um cabeçalho com nome de host.

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

**Motivo:** os menus do Windows são reconstruídos do zero quando atualizados. Se você não chamar `Update()`, os manipuladores de clique não serão acionados corretamente.

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

### Estado de seleção

Para caixas de seleção e itens de opção, controle ou consulte o estado de seleção:

```go
checkbox := menu.AddCheckbox("Feature", false)
checkbox.SetChecked(true)
menu.Update()

// Query state
isChecked := checkbox.Checked()
```

**Alternância automática:** as caixas de seleção alternam automaticamente ao serem clicadas. Você não precisa chamar `SetChecked()` no manipulador de clique.

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

Adicione aos itens de menu um texto exibido ao passar o cursor (a compatibilidade varia conforme a plataforma):

```go
menuItem := menu.Add("Advanced Options")
menuItem.SetTooltip("Configure advanced settings")
```

**Compatibilidade com plataformas:**

- **Windows:** ✅ Compatível
- **macOS:** ❌ Não compatível (dicas de ferramenta não são padrão em menus)
- **Linux:** ⚠️ Varia conforme o ambiente de desktop

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

- `ctx.ClickedMenuItem()` — o item de menu que foi clicado
- Contexto da janela (se vier de um menu de janela)
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

**Boa prática:** defina o manipulador uma única vez e, se necessário, use lógica condicional dentro dele.

## Menus dinâmicos

### Atualização de itens de menu

**A regra de ouro:** sempre chame `menu.Update()` após alterar o estado do menu.

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
- **macOS/Linux:** menos crítico, mas ainda recomendado
- **Manipuladores de clique:** não serão acionados corretamente sem Update()

### Reconstrução de menus

Para alterações significativas, reconstrua todo o menu:

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

- Alterações na lista de arquivos recentes
- Alterações nos menus de plugins
- Transições importantes de estado

**Quando atualizar:**

- Habilitar ou desabilitar itens
- Alterar rótulos
- Alternar caixas de seleção

### Menus sensíveis ao contexto

Ajuste os menus com base no estado da aplicação:

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
| **Windows** | Parte superior da janela | Menu específico de cada janela |
| **Linux** | Parte superior da janela | Específico de cada janela (geralmente) |

### Menus padrão

**macOS:**

- Tem o menu "Aplicação" (com o nome da aplicação)
- "Preferências" no menu Aplicação
- "Encerrar" no menu Aplicação

**Windows/Linux:**

- Não têm menu Aplicação
- "Preferências" no menu Editar ou Ferramentas
- "Sair" no menu Arquivo

**Exemplo: estrutura adequada à plataforma**

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

### Convenções de teclas de atalho

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

## Boas práticas

### ✅ Faça

- **Chame menu.Update()** depois de alterar o estado do menu (especialmente no Windows)
- **Use grupos de opções** para opções mutuamente exclusivas
- **Use caixas de seleção** para funcionalidades que podem ser ativadas ou desativadas
- **Adicione teclas de atalho** às ações comuns
- **Agrupe itens relacionados** com separadores
- **Teste em todas as plataformas** — o comportamento varia

### ❌ Não faça

- **Não se esqueça de chamar menu.Update()** — os manipuladores de clique não funcionarão corretamente
- **Não crie níveis demais de aninhamento** — no máximo 2-3 níveis
- **Não comece nem termine com separadores** — isso parece pouco profissional
- **Não use dicas de ferramenta no macOS** — elas não são compatíveis
- **Não fixe atalhos específicos da plataforma no código** — use `CmdOrCtrl`

## Solução de problemas

### Itens de menu não respondem

**Sintoma:** os manipuladores de clique não são acionados

**Causa:** não foi chamado `menu.Update()` depois de habilitar o item

**Solução:**

```go
menuItem.SetEnabled(true)
menu.Update()  // Add this!
```

### Itens de menu aparecem esmaecidos

**Sintoma:** não é possível clicar nos itens de menu

**Causa:** os itens estão desabilitados

**Solução:**

```go
menuItem.SetEnabled(true)
menu.Update()
```

### Teclas de atalho não funcionam

**Sintoma:** os atalhos de teclado não acionam os itens de menu

**Causas:**

1. Formato incorreto da tecla de atalho
2. Conflito com atalhos do sistema
3. A janela não está em foco

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

- [Menus da aplicação](/features/menus/application/) — crie barras de menus para a aplicação
- [Menus de contexto](/features/menus/context/) — menus de contexto acionados com o botão direito
- [Menus da bandeja do sistema](/features/menus/systray/) — menus da bandeja do sistema/barra de menus
- [Padrões de menus](/guides/menus/) — padrões comuns de menus e boas práticas

---

**Dúvidas?** Pergunte no [Discord](https://discord.gg/JDdSxwjhGf) ou consulte os [exemplos de menus](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/menu).
