---
title: "Menus de contexto"
description: "Crie menus de contexto de clique com o botão direito para seu aplicativo"
slug: "features/menus/context"
sourcePath: "features/menus/context.md"
---

## O problema

Os usuários esperam menus de clique com o botão direito que ofereçam ações específicas ao contexto. Elementos diferentes precisam de menus diferentes:

- **Texto**: Recortar, Copiar, Colar
- **Imagens**: Salvar, Copiar, Abrir
- **Elementos personalizados**: Ações específicas do aplicativo

Criar menus de contexto manualmente exige lidar com eventos do mouse, posicionamento e diferenças entre plataformas.

## A solução do Wails

O Wails fornece **menus de contexto declarativos** por meio de propriedades CSS. Associe menus a elementos HTML, passe dados e trate cliques — tudo com o comportamento nativo da plataforma.

![Um menu de contexto personalizado do Wails exibido sobre uma webview no macOS](/assets/screenshots/context-menu-macos.png)

O menu é nativo da plataforma, enquanto o elemento que o abriu continua fazendo parte da sua webview. Esta captura do macOS usa o registro de menu de contexto personalizado do exemplo abaixo.

## Início rápido

**Código Go:**

```go
// Create context menu
contextMenu := app.ContextMenu.New()
contextMenu.Add("Cut").OnClick(handleCut)
contextMenu.Add("Copy").OnClick(handleCopy)
contextMenu.Add("Paste").OnClick(handlePaste)

// Register with ID
app.ContextMenu.Add("editor-menu", contextMenu)
```

**HTML:**

```html
<textarea style="--custom-contextmenu: editor-menu">
    Right-click me!
</textarea>
```

**É só isso!** Clicar com o botão direito na área de texto exibe seu menu personalizado.

## Criação de menus de contexto

### Menu de contexto básico

```go
// Create menu
contextMenu := app.ContextMenu.New()

// Add items
contextMenu.Add("Cut").SetAccelerator("CmdOrCtrl+X").OnClick(func(ctx *application.Context) {
    // Handle cut
})

contextMenu.Add("Copy").SetAccelerator("CmdOrCtrl+C").OnClick(func(ctx *application.Context) {
    // Handle copy
})

contextMenu.Add("Paste").SetAccelerator("CmdOrCtrl+V").OnClick(func(ctx *application.Context) {
    // Handle paste
})

// Register with unique ID
app.ContextMenu.Add("text-menu", contextMenu)
```

**ID do menu:** deve ser exclusivo. É usado para associar o menu a elementos HTML.

### Com submenus

```go
contextMenu := app.ContextMenu.New()

// Add regular items
contextMenu.Add("Open").OnClick(handleOpen)
contextMenu.Add("Delete").OnClick(handleDelete)

contextMenu.AddSeparator()

// Add submenu
exportMenu := contextMenu.AddSubmenu("Export As")
exportMenu.Add("PNG").OnClick(exportPNG)
exportMenu.Add("JPEG").OnClick(exportJPEG)
exportMenu.Add("SVG").OnClick(exportSVG)

app.ContextMenu.Add("image-menu", contextMenu)
```

### Com caixas de seleção e grupos de opções

```go
contextMenu := app.ContextMenu.New()

// Checkbox
contextMenu.AddCheckbox("Show Grid", true).OnClick(func(ctx *application.Context) {
    showGrid := ctx.ClickedMenuItem().Checked()
    // Toggle grid
})

contextMenu.AddSeparator()

// Radio group
contextMenu.AddRadio("Small", false).OnClick(handleSize)
contextMenu.AddRadio("Medium", true).OnClick(handleSize)
contextMenu.AddRadio("Large", false).OnClick(handleSize)

app.ContextMenu.Add("view-menu", contextMenu)
```

Consulte a [Referência de menus](/features/menus/reference/) **para conhecer todos os tipos de item de menu**.

## Associação a elementos HTML

Use propriedades personalizadas de CSS para anexar menus de contexto:

### Associação básica

```html
<div style="--custom-contextmenu: menu-id">
    Right-click me!
</div>
```

**Propriedade CSS:** `--custom-contextmenu: <menu-id>`

### Com dados de contexto

Passe dados do HTML para o Go:

```html
<div style="--custom-contextmenu: file-menu; --custom-contextmenu-data: file-123">
    Right-click this file
</div>
```

**Manipulador Go:**

```go
contextMenu := app.ContextMenu.New()
contextMenu.Add("Open").OnClick(func(ctx *application.Context) {
    fileID := ctx.ContextMenuData()  // "file-123"
    openFile(fileID)
})

app.ContextMenu.Add("file-menu", contextMenu)
```

**Propriedades CSS:**

- `--custom-contextmenu: <menu-id>` — Qual menu exibir
- `--custom-contextmenu-data: <data>` — Dados a serem passados aos manipuladores

### Dados dinâmicos

Gere dados dinamicamente em JavaScript:

```html
<div id="file-item" style="--custom-contextmenu: file-menu">
    File.txt
</div>

<script>
// Set data dynamically
const fileItem = document.getElementById('file-item')
fileItem.style.setProperty('--custom-contextmenu-data', 'file-' + fileId)
</script>
```

### Vários elementos, o mesmo menu

```html
<div class="file-item" style="--custom-contextmenu: file-menu; --custom-contextmenu-data: file-1">
    Document.pdf
</div>

<div class="file-item" style="--custom-contextmenu: file-menu; --custom-contextmenu-data: file-2">
    Image.png
</div>

<div class="file-item" style="--custom-contextmenu: file-menu; --custom-contextmenu-data: file-3">
    Video.mp4
</div>
```

**Um menu, com dados diferentes para cada elemento.**

## Dados de contexto

### Acesso aos dados de contexto

```go
contextMenu.Add("Process").OnClick(func(ctx *application.Context) {
    data := ctx.ContextMenuData()  // Get data from HTML
    
    // Use the data
    processItem(data)
})
```

**Tipo de dado:** sempre `string`. Faça a análise conforme necessário.

### Passagem de dados complexos

Use JSON para dados complexos:

```html
<div style="--custom-contextmenu: item-menu; --custom-contextmenu-data: {&quot;id&quot;:123,&quot;type&quot;:&quot;image&quot;}">
    Image.png
</div>
```

**Manipulador Go:**

```go
import "encoding/json"

type ItemData struct {
    ID   int    `json:"id"`
    Type string `json:"type"`
}

contextMenu.Add("Process").OnClick(func(ctx *application.Context) {
    dataStr := ctx.ContextMenuData()
    
    var data ItemData
    if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
        log.Printf("Invalid data: %v", err)
        return
    }
    
    processItem(data.ID, data.Type)
})
```

@note{type="caution" title="Segurança"}
**Sempre valide os dados de contexto** recebidos do frontend. Os usuários podem manipular propriedades CSS; portanto, trate esses dados como entrada não confiável.

@end

### Exemplo de validação

```go
contextMenu.Add("Delete").OnClick(func(ctx *application.Context) {
    fileID := ctx.ContextMenuData()
    
    // Validate
    if !isValidFileID(fileID) {
        log.Printf("Invalid file ID: %s", fileID)
        return
    }
    
    // Check permissions
    if !canDeleteFile(fileID) {
        showError("Permission denied")
        return
    }
    
    // Safe to proceed
    deleteFile(fileID)
})
```

## Menu de contexto padrão

A WebView fornece um menu de contexto integrado para operações padrão (copiar, colar, inspecionar). Controle-o com `--default-contextmenu`:

### Ocultar o menu padrão

```html
<div style="--default-contextmenu: hide">
    No default menu here
</div>
```

**Caso de uso:** elementos personalizados da interface nos quais o menu padrão não faz sentido.

### Exibir o menu padrão

```html
<div style="--default-contextmenu: show">
    Default menu always shown
</div>
```

**Caso de uso:** áreas de texto, campos de entrada e conteúdo editável.

### Modo automático (inteligente)

```html
<div style="--default-contextmenu: auto">
    Smart context menu
</div>
```

**Comportamento padrão.** Exibe o menu padrão quando:

- Há texto selecionado
- O usuário está em campos de entrada de texto
- O usuário está em conteúdo editável (`contenteditable`)

Caso contrário, oculta o menu padrão.

### Combinação dos menus personalizado e padrão

```html
<!-- Custom menu + default menu -->
<textarea style="--custom-contextmenu: editor-menu; --default-contextmenu: show">
    Both menus available
</textarea>
```

**Comportamento:**

1. O menu personalizado é exibido primeiro
2. Se o menu personalizado estiver vazio ou não for encontrado, o menu padrão será exibido
3. Ambos podem coexistir (dependendo da plataforma)

## Menus de contexto dinâmicos

Atualize os menus de acordo com o estado do aplicativo:

### Ativar ou desativar itens

```go
var cutMenuItem *application.MenuItem
var copyMenuItem *application.MenuItem

func createContextMenu() {
    contextMenu := app.ContextMenu.New()
    
    cutMenuItem = contextMenu.Add("Cut")
    cutMenuItem.SetEnabled(false)  // Initially disabled
    cutMenuItem.OnClick(handleCut)
    
    copyMenuItem = contextMenu.Add("Copy")
    copyMenuItem.SetEnabled(false)
    copyMenuItem.OnClick(handleCopy)
    
    app.ContextMenu.Add("editor-menu", contextMenu)
}

func onSelectionChanged(hasSelection bool) {
    cutMenuItem.SetEnabled(hasSelection)
    copyMenuItem.SetEnabled(hasSelection)
    contextMenu.Update()  // Important!
}
```

@note{type="caution" title="Sempre chame Update()"}
Depois de alterar o estado do menu, **chame `contextMenu.Update()`**. Isso é essencial no Windows.

Consulte a [Referência de menus](/features/menus/reference/#enabled-state) para obter detalhes.

@end

### Alterar rótulos

```go
playMenuItem := contextMenu.Add("Play")

playMenuItem.OnClick(func(ctx *application.Context) {
    if isPlaying {
        playMenuItem.SetLabel("Pause")
    } else {
        playMenuItem.SetLabel("Play")
    }
    contextMenu.Update()
})
```

### Recriar menus

Para alterações significativas, recrie o menu inteiro:

```go
func rebuildContextMenu(fileType string) {
    contextMenu := app.ContextMenu.New()
    
    // Common items
    contextMenu.Add("Open").OnClick(handleOpen)
    contextMenu.Add("Delete").OnClick(handleDelete)
    
    contextMenu.AddSeparator()
    
    // Type-specific items
    switch fileType {
    case "image":
        contextMenu.Add("Edit Image").OnClick(editImage)
        contextMenu.Add("Set as Wallpaper").OnClick(setWallpaper)
    case "video":
        contextMenu.Add("Play").OnClick(playVideo)
        contextMenu.Add("Extract Audio").OnClick(extractAudio)
    case "document":
        contextMenu.Add("Print").OnClick(printDocument)
        contextMenu.Add("Export PDF").OnClick(exportPDF)
    }
    
    app.ContextMenu.Add("file-menu", contextMenu)
}
```

## Comportamento nas plataformas

Os menus de contexto são **nativos da plataforma**:

@tabs{sync-key="platform"}
[macOS]
**Menus de contexto nativos do macOS:**

- Animações e transições do sistema
- Clique com o botão direito = Control+Clique (automático)
- Adaptação à aparência do sistema (clara/escura)
- Operações de texto padrão no menu predefinido
- Rolagem nativa para menus longos

**Convenções do macOS:**

- Use apenas a primeira palavra com inicial maiúscula nos itens de menu
- Use reticências (...) nos itens que abrem caixas de diálogo
- Atalhos comuns: ⌘C (Copiar), ⌘V (Colar)

[Windows]
**Menus de contexto nativos do Windows:**

- Estilo nativo do Windows
- Acompanha o tema do Windows
- Operações padrão do Windows no menu predefinido
- Compatibilidade com entrada por toque e caneta

**Convenções do Windows:**

- Use iniciais maiúsculas nas palavras principais de cada item de menu
- Use reticências (...) nos itens que abrem caixas de diálogo
- Atalhos comuns: Ctrl+C (Copiar), Ctrl+V (Colar)

[Linux]
**Integração com o ambiente de desktop:**

- Adaptação ao tema do desktop (GTK, Qt etc.)
- O comportamento do clique com o botão direito segue as configurações do sistema
- O conteúdo do menu predefinido varia conforme o ambiente
- O posicionamento segue as convenções do ambiente de desktop

**Considerações sobre o Linux:**

- Teste nos ambientes de desktop de destino
- GTK e Qt têm comportamentos diferentes
- Alguns ambientes de desktop personalizam os menus de contexto

@end

## Exemplo completo

**Código Go:**

```go
package main

import (
    "encoding/json"
    "log"
    "github.com/wailsapp/wails/v3/pkg/application"
)

type FileData struct {
    ID   string `json:"id"`
    Type string `json:"type"`
    Name string `json:"name"`
}

func main() {
    app := application.New(application.Options{
        Name: "Context Menu Demo",
    })

    // Create file context menu
    fileMenu := createFileMenu(app)
    app.ContextMenu.Add("file-menu", fileMenu)

    // Create image context menu
    imageMenu := createImageMenu(app)
    app.ContextMenu.Add("image-menu", imageMenu)

    // Create text context menu
    textMenu := createTextMenu(app)
    app.ContextMenu.Add("text-menu", textMenu)

    app.Window.New()
    app.Run()
}

func createFileMenu(app *application.App) *application.ContextMenu {
    menu := app.ContextMenu.New()
    
    menu.Add("Open").OnClick(func(ctx *application.Context) {
        data := parseFileData(ctx.ContextMenuData())
        openFile(data.ID)
    })
    
    menu.Add("Rename").OnClick(func(ctx *application.Context) {
        data := parseFileData(ctx.ContextMenuData())
        renameFile(data.ID)
    })
    
    menu.AddSeparator()
    
    menu.Add("Delete").OnClick(func(ctx *application.Context) {
        data := parseFileData(ctx.ContextMenuData())
        deleteFile(data.ID)
    })
    
    return menu
}

func createImageMenu(app *application.App) *application.ContextMenu {
    menu := app.ContextMenu.New()
    
    menu.Add("View").OnClick(func(ctx *application.Context) {
        data := parseFileData(ctx.ContextMenuData())
        viewImage(data.ID)
    })
    
    menu.Add("Edit").OnClick(func(ctx *application.Context) {
        data := parseFileData(ctx.ContextMenuData())
        editImage(data.ID)
    })
    
    menu.AddSeparator()
    
    exportMenu := menu.AddSubmenu("Export As")
    exportMenu.Add("PNG").OnClick(exportPNG)
    exportMenu.Add("JPEG").OnClick(exportJPEG)
    exportMenu.Add("WebP").OnClick(exportWebP)
    
    return menu
}

func createTextMenu(app *application.App) *application.ContextMenu {
    menu := app.ContextMenu.New()
    
    menu.Add("Cut").SetAccelerator("CmdOrCtrl+X").OnClick(handleCut)
    menu.Add("Copy").SetAccelerator("CmdOrCtrl+C").OnClick(handleCopy)
    menu.Add("Paste").SetAccelerator("CmdOrCtrl+V").OnClick(handlePaste)
    
    return menu
}

func parseFileData(dataStr string) FileData {
    var data FileData
    if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
        log.Printf("Invalid file data: %v", err)
    }
    return data
}

// Handler implementations...
func openFile(id string) { /* ... */ }
func renameFile(id string) { /* ... */ }
func deleteFile(id string) { /* ... */ }
func viewImage(id string) { /* ... */ }
func editImage(id string) { /* ... */ }
func exportPNG(ctx *application.Context) { /* ... */ }
func exportJPEG(ctx *application.Context) { /* ... */ }
func exportWebP(ctx *application.Context) { /* ... */ }
func handleCut(ctx *application.Context) { /* ... */ }
func handleCopy(ctx *application.Context) { /* ... */ }
func handlePaste(ctx *application.Context) { /* ... */ }
```

**HTML:**

```html
<!DOCTYPE html>
<html>
<head>
    <style>
        .file-item {
            padding: 10px;
            margin: 5px;
            border: 1px solid #ccc;
            cursor: pointer;
        }
        
        .file-item:hover {
            background: #f0f0f0;
        }
        
        textarea {
            width: 100%;
            height: 200px;
        }
    </style>
</head>
<body>
    <h2>Files</h2>
    
    <!-- Regular file -->
    <div class="file-item" 
         style="--custom-contextmenu: file-menu; 
                --custom-contextmenu-data: {&quot;id&quot;:&quot;file-1&quot;,&quot;type&quot;:&quot;document&quot;,&quot;name&quot;:&quot;Report.pdf&quot;}">
        📄 Report.pdf
    </div>
    
    <!-- Image file -->
    <div class="file-item" 
         style="--custom-contextmenu: image-menu; 
                --custom-contextmenu-data: {&quot;id&quot;:&quot;file-2&quot;,&quot;type&quot;:&quot;image&quot;,&quot;name&quot;:&quot;Photo.jpg&quot;}">
        🖼️ Photo.jpg
    </div>
    
    <h2>Text Editor</h2>
    
    <!-- Text area with custom menu + default menu -->
    <textarea 
        style="--custom-contextmenu: text-menu; --default-contextmenu: show"
        placeholder="Type here, then right-click...">
    </textarea>
    
    <h2>No Context Menu</h2>
    
    <!-- Disable default menu -->
    <div style="--default-contextmenu: hide; padding: 20px; border: 1px solid #ccc;">
        Right-click here - no menu appears
    </div>
</body>
</html>
```

## Práticas recomendadas

### ✅ Faça

- **Mantenha os menus focados** — Inclua apenas ações relevantes para o elemento
- **Valide os dados de contexto** — Trate-os como entrada não confiável
- **Use rótulos claros** — Use "Excluir arquivo", não "Excluir"
- **Chame menu.Update()** — Após alterar o estado do menu
- **Teste em todas as plataformas** — O comportamento varia
- **Forneça atalhos de teclado** — Para ações comuns
- **Agrupe itens relacionados** — Use separadores

### ❌ Não faça

- **Não confie nos dados de contexto** — Sempre os valide
- **Não crie menus longos demais** — No máximo 7-10 itens
- **Não se esqueça de menu.Update()** — Os menus não funcionarão corretamente
- **Não crie aninhamentos muito profundos** — No máximo 2 níveis
- **Não use jargão** — Mantenha os rótulos fáceis de entender
- **Não bloqueie os manipuladores** — Mantenha-os rápidos

## Solução de problemas

### O menu de contexto não aparece

**Possíveis causas:**

1. Incompatibilidade no ID do menu
2. Erro de digitação na propriedade CSS
3. Runtime não inicializado

**Solução:**

```go
// Check menu is registered
app.ContextMenu.Add("my-menu", contextMenu)
```

```html
<!-- Check ID matches -->
<div style="--custom-contextmenu: my-menu">
```

### Dados de contexto não recebidos

**Possíveis causas:**

1. Propriedade CSS não definida
2. Os dados contêm caracteres especiais

**Solução:**

```html
<!-- Escape quotes in JSON -->
<div style="--custom-contextmenu-data: {&quot;id&quot;:123}">
```

Ou use JavaScript:

```javascript
element.style.setProperty('--custom-contextmenu-data', JSON.stringify(data))
```

### Os itens de menu não respondem

**Causa:** não chamar `menu.Update()` após habilitar

**Solução:**

```go
menuItem.SetEnabled(true)
contextMenu.Update()  // Add this!
```

## Próximas etapas

@cards{cols="2"}
📖 Referência de menus
Referência completa dos tipos e das propriedades dos itens de menu.

[Saiba mais →](/features/menus/reference/)

---
☰ Menus do aplicativo
Crie barras de menu para o aplicativo.

[Saiba mais →](/features/menus/application/)

---
★ Menus da bandeja do sistema
Adicione integração com a bandeja do sistema ou a barra de menus.

[Saiba mais →](/features/menus/systray/)

---
📖 Padrões de menus
Padrões comuns de menus e práticas recomendadas.

[Saiba mais →](/guides/menus/)

@end

---

**Tem alguma dúvida?** Pergunte no [Discord](https://discord.gg/JDdSxwjhGf) ou consulte o [exemplo de menu de contexto](https://github.com/wailsapp/wails/tree/master/v3/examples/contextmenus).
