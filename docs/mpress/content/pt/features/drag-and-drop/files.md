---
title: "Soltar arquivos"
description: "Aceite arquivos arrastados do sistema operacional para o seu aplicativo"
slug: "features/drag-and-drop/files"
sourcePath: "features/drag-and-drop/files.md"
---

O Wails permite que os usuários arrastem arquivos do sistema operacional (gerenciador de arquivos, área de trabalho) para o seu aplicativo. Ao contrário do recurso de arrastar e soltar do HTML5, que funciona apenas dentro do navegador, isso permite acessar os caminhos reais dos arquivos no disco.

![Um exemplo de arrastar e soltar no Wails com uma zona para soltar arquivos externos no macOS](/assets/screenshots/file-drop-macos.png)

A zona para soltar arquivos externos faz parte da webview, enquanto o Wails fornece os eventos nativos de arquivos soltos do sistema operacional.

## Ativar o recurso de soltar arquivos

O recurso de soltar arquivos fica desativado por padrão. Para ativá-lo, defina `EnableFileDrop: true` nas opções da janela:

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:          "My App",
    Width:          800,
    Height:         600,
    EnableFileDrop: true,
})
```

Quando `EnableFileDrop` é `false` (o padrão), os arquivos arrastados do sistema operacional são bloqueados: eles não são abertos na webview nem disparam eventos. Isso evita a navegação acidental quando os usuários arrastam arquivos sobre o aplicativo.

## Definir zonas para soltar arquivos

As zonas para soltar arquivos indicam ao Wails quais elementos devem aceitar arquivos. Os arquivos soltos fora de uma dessas zonas são ignorados.

Adicione o atributo `data-file-drop-target` a qualquer elemento:

```html
<div id="upload" class="drop-zone" data-file-drop-target>
    Drop files here
</div>
```

Você pode ter várias zonas para soltar arquivos. O `id` e as classes CSS do elemento são passados ao seu código Go, permitindo processar os arquivos de maneira diferente conforme o local em que forem soltos.

## Estilizar ao arrastar sobre a zona

Quando arquivos são arrastados sobre uma zona para soltar arquivos, o Wails adiciona a classe `file-drop-target-active`. Isso permite fornecer um retorno visual para que os usuários saibam onde podem soltá-los:

```css
.drop-zone {
    border: 2px dashed #ccc;
    padding: 40px;
    text-align: center;
    transition: all 0.2s ease;
}

.drop-zone.file-drop-target-active {
    border-color: #007bff;
    background-color: rgba(0, 123, 255, 0.1);
}
```

A classe é removida automaticamente quando os arquivos saem da zona ou são soltos.

## Detectar arquivos soltos

Quando arquivos são soltos em uma zona válida, o Wails dispara um evento `WindowFilesDropped`. O contexto do evento contém os caminhos completos de todos os arquivos soltos no sistema de arquivos:

```go
import "github.com/wailsapp/wails/v3/pkg/events"

window.OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {
    files := event.Context().DroppedFiles()
    
    for _, file := range files {
        fmt.Println("Dropped:", file)
    }
})
```

Os caminhos são absolutos, como `/home/user/documents/report.pdf` ou `C:\Users\Name\Documents\report.pdf`.

## Obter informações sobre o destino

Quando houver várias zonas para soltar arquivos, você poderá descobrir qual delas recebeu os arquivos usando `DropTargetDetails()`:

```go
window.OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {
    files := event.Context().DroppedFiles()
    details := event.Context().DropTargetDetails()
    
    fmt.Printf("Dropped on element: id=%s, classes=%v\n", 
        details.ElementID, details.ClassList)
    fmt.Printf("Position: x=%d, y=%d\n", details.X, details.Y)
})
```

Isso permite encaminhar os arquivos para diferentes manipuladores:

```go
switch details.ElementID {
case "images":
    handleImageUpload(files)
case "documents":
    handleDocumentUpload(files)
}
```

## Exemplo completo

**Go:**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:          "File Uploader",
    EnableFileDrop: true,
})

window.OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {
    files := event.Context().DroppedFiles()
    details := event.Context().DropTargetDetails()
    
    // Send to frontend
    app.Event.Emit("files-dropped", map[string]any{
        "files":   files,
        "target":  details.ElementID,
    })
})
```

**HTML:**

```html
<div id="images" class="drop-zone" data-file-drop-target>
    Drop images here
</div>

<div id="documents" class="drop-zone" data-file-drop-target>
    Drop documents here
</div>

<style>
    .drop-zone {
        border: 2px dashed #ccc;
        border-radius: 8px;
        padding: 40px;
        text-align: center;
        margin: 20px;
        transition: all 0.2s ease;
    }
    
    .drop-zone.file-drop-target-active {
        border-color: #007bff;
        background-color: rgba(0, 123, 255, 0.1);
    }
</style>
```

## Soltar arquivos em toda a janela

Se quiser que seja possível soltar arquivos em qualquer parte do aplicativo, adicione o atributo ao elemento body:

```html
<body data-file-drop-target>
    <!-- Your app content -->
</body>
```

Você pode usar uma sobreposição CSS para indicar que toda a janela é um destino para soltar arquivos:

```css
body.file-drop-target-active::after {
    content: "Drop files anywhere";
    position: fixed;
    inset: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 24px;
    color: #007bff;
    background: rgba(255, 255, 255, 0.9);
    pointer-events: none;
}
```

## Combinar com o recurso de arrastar e soltar do HTML

Você pode usar tanto o recurso de soltar arquivos externos quanto o de arrastar e soltar elementos HTML internos no mesmo aplicativo. Quando `EnableFileDrop` é `true`, o Wails intercepta os arquivos externos arrastados, mas permite que os elementos HTML5 internos arrastados sejam processados normalmente.

Para distingui-los nos manipuladores das zonas para soltar arquivos em HTML, verifique se os dados arrastados contêm arquivos:

```javascript
zone.addEventListener('dragenter', (e) => {
    // Skip external file drags - Wails handles these
    if (e.dataTransfer?.types.includes('Files')) {
        return;
    }
    // Handle internal HTML5 drags
    zone.classList.add('drag-over');
});

zone.addEventListener('drop', (e) => {
    // Skip external file drops - Wails handles these
    if (e.dataTransfer?.types.includes('Files')) {
        return;
    }
    e.preventDefault();
    zone.classList.remove('drag-over');
    // Handle internal drop
});
```

Isso garante que seus manipuladores de elementos soltos em HTML respondam apenas a elementos internos arrastados (como ao mover itens de uma lista), enquanto o Wails processa separadamente os arquivos externos soltos por meio do evento `WindowFilesDropped`.

## Próximas etapas

- [Arrastar e soltar com HTML](/features/drag-and-drop/html/) — Arraste elementos dentro do seu aplicativo
- [Opções da janela](/features/windows/options/) — Todas as opções de configuração da janela
