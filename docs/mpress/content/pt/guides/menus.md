---
title: "Menus"
description: "Um guia para criar e personalizar menus no Wails v3"
slug: "guides/menus"
sourcePath: "guides/menus.md"
---

O Wails v3 oferece um poderoso sistema de menus que permite criar tanto menus de aplicativo quanto menus de contexto. Este guia apresenta os diversos recursos e funcionalidades do sistema de menus.

## Como criar um menu

Para criar um menu, use o método `New()` do gerenciador Menus:

```go
menu := app.Menu.New()
```

### Como adicionar itens de menu

O Wails oferece suporte a vários tipos de itens de menu, cada um com uma finalidade específica:

#### Itens de menu comuns

Os itens de menu comuns são os componentes básicos dos menus. Eles exibem texto e podem acionar ações quando clicados:

```go
menuItem := menu.Add("Click Me")
```

#### Caixas de seleção

Os itens de menu com caixa de seleção oferecem um estado alternável, útil para ativar ou desativar recursos ou configurações:

```go
checkbox := menu.AddCheckbox("My checkbox", true)  // true = initially checked
```

#### Grupos de botões de opção

Os grupos de botões de opção permitem selecionar uma alternativa em um conjunto de opções mutuamente exclusivas. Eles são criados automaticamente quando itens desse tipo são posicionados lado a lado:

```go
menu.AddRadio("Option 1", true)   // true = initially selected
menu.AddRadio("Option 2", false)
menu.AddRadio("Option 3", false)
```

#### Separadores

Os separadores são linhas horizontais que ajudam a organizar os itens de menu em grupos lógicos:

```go
menu.AddSeparator()
```

#### Submenus

Submenus são menus aninhados que aparecem ao passar o cursor sobre um item de menu ou clicar nele. Eles são úteis para organizar estruturas de menu complexas:

```go
submenu := menu.AddSubmenu("File")
submenu.Add("Open")
submenu.Add("Save")
```

#### Como combinar menus

Um menu pode ser adicionado a outro menu, no início ou no final.

```go
menu := app.Menu.New()
menu.Add("First Menu")

secondaryMenu := app.Menu.New()
secondaryMenu.Add("Second Menu")

// insert 'secondaryMenu' after 'menu'
menu.Append(secondaryMenu)

// insert 'secondaryMenu' before 'menu'
menu.Prepend(secondaryMenu)

// update the menu
menu.Update()
```

@note{type="info"}
Por padrão, `prepend` e `append` compartilham o estado com o menu original. Para criar um menu com estado próprio, você pode chamar `.Clone()` no menu.

Por exemplo: `menu.Append(secondaryMenu.Clone())`

@end

#### Como limpar um menu

Em alguns casos, será melhor criar um menu totalmente novo se você estiver trabalhando com uma quantidade variável de itens de menu.

Isso remove todos os itens de um menu existente e permite adicionar itens novamente.

```go
menu := app.Menu.New()
menu.Add("Waiting for update...")

// after certain logic, the menu has to be updated
menu.Clear()
menu.Add("Update complete!")
menu.Update()
```

@note{type="info"}
A limpeza de um menu remove apenas os itens de menu do nível superior. Embora os submenus deixem de ficar visíveis, eles ainda ocuparão memória; portanto, gerencie seus menus com cuidado.

@end

#### Como destruir um menu

Para limpar e liberar um menu, use o método `Destroy()`:

```go
menu := app.Menu.New()
menu.Add("Waiting for update...")

// after certain logic, the menu has to be destroyed
menu.Destroy()
```

### Propriedades dos itens de menu

Os itens de menu têm várias propriedades que podem ser configuradas:

| Propriedade | Método | Descrição |
| --- | --- | --- |
| Rótulo | `SetLabel(string)` | Define o texto exibido |
| Ativado | `SetEnabled(bool)` | Ativa ou desativa o item |
| Marcado | `SetChecked(bool)` | Define o estado de seleção (para caixas de seleção ou botões de opção) |
| Dica de ferramenta | `SetTooltip(string)` | Define o texto da dica de ferramenta |
| Oculto | `SetHidden(bool)` | Exibe ou oculta o item |
| Tecla de atalho | `SetAccelerator(string)` | Define o atalho de teclado |

### Estados dos itens de menu

Os itens de menu podem ter diferentes estados que controlam sua visibilidade e interatividade:

#### Visibilidade

Os itens de menu podem ser exibidos ou ocultados dinamicamente usando o método `SetHidden()`:

```go
menuItem := menu.Add("Dynamic Item")

// Hide the menu item
menuItem.SetHidden(true)

// Show the menu item
menuItem.SetHidden(false)

// Check current visibility
isHidden := menuItem.Hidden()
```

Os itens de menu ocultos são completamente removidos do menu até serem exibidos novamente. Isso é útil para itens de menu contextuais que devem aparecer somente em determinados estados do aplicativo.

#### Estado de ativação

Os itens de menu podem ser ativados ou desativados usando o método `SetEnabled()`:

```go
menuItem := menu.Add("Save")

// Disable the menu item
menuItem.SetEnabled(false)  // Item appears grayed out and cannot be clicked

// Enable the menu item
menuItem.SetEnabled(true)   // Item becomes clickable again

// Check current enabled state
isEnabled := menuItem.Enabled()
```

Os itens de menu desativados permanecem visíveis, mas aparecem esmaecidos e não podem ser clicados. Isso geralmente é usado para indicar que uma ação está indisponível no momento, como nos seguintes casos:

- Desativar "Salvar" quando não houver alterações para salvar
- Desativar "Copiar" quando nada estiver selecionado
- Desativar "Desfazer" quando não houver nenhuma ação para desfazer

#### Gerenciamento dinâmico de estados

Você pode combinar esses estados com manipuladores de eventos para criar menus dinâmicos:

```go
saveMenuItem := menu.Add("Save")

// Initially disable the Save menu item
saveMenuItem.SetEnabled(false)

// Enable Save only when there are unsaved changes
documentChanged := func() {
    saveMenuItem.SetEnabled(true)
    menu.Update()  // Remember to update the menu after changing states
}

// Disable Save after saving
documentSaved := func() {
    saveMenuItem.SetEnabled(false)
    menu.Update()
}
```

### Tratamento de eventos

Os itens de menu podem responder a eventos de clique usando o método `OnClick`:

```go
menuItem.OnClick(func(ctx *application.Context) {
    // Handle the click event
    println("Menu item clicked!")
})
```

O contexto fornece informações sobre o item de menu clicado:

```go
menuItem.OnClick(func(ctx *application.Context) {
    // Get the clicked menu item
    clickedItem := ctx.ClickedMenuItem()
    // Get its current state
    isChecked := clickedItem.Checked()
})
```

### Itens de menu baseados em funções

O Wails fornece um conjunto de funções de menu predefinidas que criam automaticamente itens de menu com funcionalidades padrão. Estas são as funções de menu compatíveis:

#### Estruturas de menu completas

Estas funções criam estruturas de menu completas com funcionalidades comuns:

| Função | Descrição | Observações sobre plataformas |
| --- | --- | --- |
| `AppMenu` | Menu do aplicativo com Sobre, Serviços, Ocultar/Mostrar e Sair | Somente macOS |
| `EditMenu` | Menu Editar padrão com Desfazer, Refazer, Recortar, Copiar, Colar etc. | Todas as plataformas |
| `ViewMenu` | Menu Exibir com controles de recarregamento, zoom e tela cheia | Todas as plataformas |
| `WindowMenu` | Controles de janela (Minimizar, Zoom etc.) | Todas as plataformas |
| `HelpMenu` | Menu Ajuda com o link "Saiba mais" para o site do Wails | Todas as plataformas |

#### Itens de menu individuais

Estas funções podem ser usadas para adicionar itens de menu individuais:

| Função | Descrição | Observações sobre plataformas |
| --- | --- | --- |
| `About` | Exibir a caixa de diálogo Sobre do aplicativo | Todas as plataformas |
| `Hide` | Ocultar o aplicativo | Somente macOS |
| `HideOthers` | Ocultar outros aplicativos | Somente macOS |
| `UnHide` | Exibir o aplicativo oculto | Somente macOS |
| `CloseWindow` | Fechar a janela atual | Todas as plataformas |
| `Minimise` | Minimizar a janela | Todas as plataformas |
| `Zoom` | Aplicar zoom à janela | Somente macOS |
| `Front` | Trazer a janela para a frente | Somente macOS |
| `Quit` | Encerrar o aplicativo | Todas as plataformas |
| `Undo` | Desfazer a última ação | Todas as plataformas |
| `Redo` | Refazer a última ação | Todas as plataformas |
| `Cut` | Recortar seleção | Todas as plataformas |
| `Copy` | Copiar seleção | Todas as plataformas |
| `Paste` | Colar da área de transferência | Todas as plataformas |
| `PasteAndMatchStyle` | Colar e aplicar o mesmo estilo | Somente macOS |
| `SelectAll` | Selecionar tudo | Todas as plataformas |
| `Delete` | Excluir seleção | Todas as plataformas |
| `Reload` | Recarregar a página atual | Todas as plataformas |
| `ForceReload` | Forçar o recarregamento da página atual | Todas as plataformas |
| `ToggleFullscreen` | Alternar o modo de tela cheia | Todas as plataformas |
| `ResetZoom` | Redefinir o nível de zoom | Todas as plataformas |
| `ZoomIn` | Aumentar o zoom | Todas as plataformas |
| `ZoomOut` | Diminuir o zoom | Todas as plataformas |

Veja um exemplo de como usar menus completos e funções individuais:

```go
menu := app.Menu.New()

// Add complete menu structures
menu.AddRole(application.AppMenu)    // macOS only
menu.AddRole(application.EditMenu)   // Common edit operations
menu.AddRole(application.ViewMenu)   // View controls
menu.AddRole(application.WindowMenu) // Window controls

// Add individual role-based items to a custom menu
fileMenu := menu.AddSubmenu("File")
fileMenu.AddRole(application.CloseWindow)
fileMenu.AddSeparator()
fileMenu.AddRole(application.Quit)
```

## Menus do aplicativo

Os menus do aplicativo são exibidos na parte superior da janela do aplicativo (Windows/Linux) ou na parte superior da tela (macOS).

### Comportamento do menu do aplicativo

Quando você define um menu do aplicativo usando `app.Menu.Set()`, ele se torna o menu principal no macOS. No Windows/Linux, os menus são definidos individualmente para cada janela.

```go
app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Custom Menu Window",
    Windows: application.WindowsWindow{
        Menu: customMenu,  // Override application menu for this window
    },
})
```

Veja um exemplo completo desses diferentes comportamentos de menu:

```go
func main() {
    app := application.New(application.Options{})

    // Create application menu
    appMenu := app.Menu.New()
    fileMenu := appMenu.AddSubmenu("File")
    fileMenu.Add("New").OnClick(func(ctx *application.Context) {
        // This will be available in all windows unless overridden
        window := app.Window.Current()
        window.SetTitle("New Window")
    })
    
    // Set as application menu - this is for macOS
    app.Menu.Set(appMenu)

    // Window with custom menu on Windows
    customMenu := app.Menu.New()
    customMenu.Add("Custom Action")
    app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title: "Custom Menu",
        Windows: application.WindowsWindow{
            Menu: customMenu,
        },
    })

    app.Run()
}
```

## Menus de contexto

Menus de contexto são menus pop-up exibidos quando você clica com o botão direito do mouse em elementos do aplicativo. Eles fornecem acesso rápido às ações relevantes para o elemento clicado.

### Menu de contexto padrão

O menu de contexto padrão é o menu de contexto integrado da webview, que fornece operações no nível do sistema, como:

- Copiar, Recortar e Colar para manipulação de texto
- Controles de seleção de texto
- Opções de verificação ortográfica

#### Como controlar o menu de contexto padrão

Você pode controlar quando o menu de contexto padrão é exibido usando a propriedade CSS `--default-contextmenu`:

```html
<!-- Always show default context menu -->
<div style="--default-contextmenu: show">
    <input type="text" placeholder="Right-click for text operations"/>
    <textarea>Standard text operations available here</textarea>
</div>

<!-- Hide default context menu -->
<div style="--default-contextmenu: hide">
    <div class="custom-component">Custom context menu only</div>
</div>

<!-- Smart context menu behaviour (default) -->
<div style="--default-contextmenu: auto">
    <!-- Shows default menu when text is selected or in input fields -->
    <p>Select this text to see the default menu</p>
    <input type="text" placeholder="Default menu for input operations"/>
</div>
```

@note{type="info"}
Esse recurso só funcionará conforme esperado depois que o [runtime de frontend estiver pronto](/reference/frontend-runtime/).

@end

#### Comportamento de menus de contexto aninhados

Ao usar a propriedade `--default-contextmenu` em elementos aninhados, aplicam-se as seguintes regras:

1. Os elementos filhos herdam a configuração do menu de contexto do elemento pai, a menos que ela seja explicitamente substituída
2. A configuração mais específica (mais próxima) tem precedência
3. O valor `auto` pode ser usado para restaurar o comportamento padrão

Exemplo do comportamento de menus de contexto aninhados:

```html
<!-- Parent sets hide -->
<div style="--default-contextmenu: hide">
    <!-- This inherits hide -->
    <p>No context menu here</p>
    
    <!-- This overrides to show -->
    <div style="--default-contextmenu: show">
        <p>Context menu shown here</p>
        
        <!-- This inherits show -->
        <span>Also has context menu</span>
        
        <!-- This resets to automatic behaviour -->
        <div style="--default-contextmenu: auto">
            <p>Shows menu only when text is selected</p>
        </div>
    </div>
</div>
```

### Menus de contexto personalizados

Os menus de contexto personalizados permitem disponibilizar ações específicas do aplicativo que sejam relevantes para o elemento clicado. Eles são particularmente úteis para:

- Operações de arquivo em um gerenciador de documentos
- Ferramentas de manipulação de imagens
- Ações personalizadas em uma grade de dados
- Operações específicas do componente

#### Como criar um menu de contexto personalizado

Ao criar um menu de contexto personalizado, forneça um identificador exclusivo (nome) que vincule o menu aos elementos HTML:

```go
// Create a context menu with identifier "imageMenu"
contextMenu := app.ContextMenu.New()
app.ContextMenu.Add("imageMenu", contextMenu)
```

O parâmetro de nome ("imageMenu" neste exemplo) funciona como um identificador exclusivo que será usado para:

1. Vincular elementos HTML a esse menu de contexto específico
2. Identificar qual menu deve ser exibido ao clicar com o botão direito
3. Permitir a atualização e a limpeza do menu

#### Dados de contexto

Ao tratar eventos do menu de contexto, você pode acessar tanto o item de menu clicado quanto os dados de contexto associados a ele:

```go
contextMenu.Add("Process").OnClick(func(ctx *application.Context) {
    // Get the clicked menu item
    menuItem := ctx.ClickedMenuItem()
    
    // Get the context data as a string
    contextData := ctx.ContextMenuData()
    
    // Check if the menu item is checked (for checkbox/radio items)
    isChecked := ctx.IsChecked()
    
    // Use the data
    if contextData != "" {
        processItem(contextData)
    }
})
```

Os dados de contexto são transmitidos pela propriedade `--custom-contextmenu-data` do elemento HTML e ficam disponíveis no manipulador de clique por meio de `ctx.ContextMenuData()`. Isso é particularmente útil para:

- Trabalhar com listas ou grades nas quais cada item precisa de uma identificação exclusiva
- Executar operações em componentes ou elementos específicos
- Transmitir estado ou metadados do frontend para o backend

#### Gerenciamento de menus de contexto

Após fazer alterações em um menu de contexto, chame o método `Update()` para aplicá-las:

```go
contextMenu.Update()
```

Quando não precisar mais de um menu de contexto, você poderá destruí-lo:

```go
contextMenu.Destroy()
```

@note{type="danger" title="Aviso"}
Após chamar `Destroy()`, reutilizar a referência ao menu de contexto causará um panic.

@end

### Exemplo real: galeria de imagens

Veja um exemplo completo de implementação de um menu de contexto personalizado para uma galeria de imagens:

```go
// Backend: Create the context menu
imageMenu := app.ContextMenu.New()
app.ContextMenu.Add("imageMenu", imageMenu)

// Add relevant operations
imageMenu.Add("View Full Size").OnClick(func(ctx *application.Context) {
    // Get the image ID from context data
    if imageID := ctx.ContextMenuData(); imageID != "" {
        openFullSizeImage(imageID)
    }
})

imageMenu.Add("Download").OnClick(func(ctx *application.Context) {
    if imageID := ctx.ContextMenuData(); imageID != "" {
        downloadImage(imageID)
    }
})

imageMenu.Add("Share").OnClick(func(ctx *application.Context) {
    if imageID := ctx.ContextMenuData(); imageID != "" {
        showShareDialog(imageID)
    }
})
```

```html
<!-- Frontend: Image gallery implementation -->
<div class="gallery">
    <!-- Each image container with context menu -->
    <div class="image-container" 
         style="--custom-contextmenu: imageMenu; --custom-contextmenu-data: img_123">
        <img src="/images/img_123.jpg" alt="Gallery Image"/>
        <span class="caption">Nature Photo</span>
    </div>
    
    <div class="image-container"
         style="--custom-contextmenu: imageMenu; --custom-contextmenu-data: img_124">
        <img src="/images/img_124.jpg" alt="Gallery Image"/>
        <span class="caption">City Photo</span>
    </div>
</div>
```

Neste exemplo:

1. O menu de contexto é criado com o identificador "imageMenu"
2. Cada contêiner de imagem é vinculado ao menu por meio de `--custom-contextmenu: imageMenu`
3. Cada contêiner fornece o ID da respectiva imagem como dado de contexto por meio de `--custom-contextmenu-data`
4. O backend recebe o ID da imagem nos manipuladores de clique e pode executar operações específicas
5. O mesmo menu é reutilizado para todas as imagens, mas os dados de contexto indicam em qual imagem a operação deve ser executada

Esse padrão é particularmente eficaz para:

- Grades de dados cujas linhas exigem operações específicas
- Gerenciadores de arquivos nos quais os arquivos exigem ações específicas ao contexto
- Ferramentas de design nas quais diferentes elementos exigem operações diferentes
- Qualquer componente em que as mesmas operações sejam aplicadas a várias instâncias
