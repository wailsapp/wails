---
title: "Dépôt de fichiers"
description: "Acceptez les fichiers glissés depuis le système d’exploitation dans votre application"
slug: "features/drag-and-drop/files"
sourcePath: "features/drag-and-drop/files.md"
---

Wails permet aux utilisateurs de faire glisser des fichiers depuis le système d’exploitation (gestionnaire de fichiers, bureau) vers votre application. Contrairement au glisser-déposer HTML5, qui ne fonctionne qu’au sein du navigateur, cette fonctionnalité vous donne accès aux chemins réels des fichiers sur le disque.

![Exemple de glisser-déposer Wails avec une zone de dépôt de fichiers externes sous macOS](/assets/screenshots/file-drop-macos.png)

La zone de dépôt externe fait partie de la vue web, tandis que Wails fournit les événements natifs du système d’exploitation liés au dépôt de fichiers.

## Activer le dépôt de fichiers

Le dépôt de fichiers est désactivé par défaut. Pour l’activer, définissez `EnableFileDrop: true` dans les options de votre fenêtre :

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:          "My App",
    Width:          800,
    Height:         600,
    EnableFileDrop: true,
})
```

Lorsque `EnableFileDrop` vaut `false` (valeur par défaut), les fichiers glissés depuis le système d’exploitation sont bloqués : ils ne s’ouvrent pas dans la vue web et ne déclenchent aucun événement. Cela empêche toute navigation accidentelle lorsque les utilisateurs font glisser des fichiers au-dessus de votre application.

## Définir les zones de dépôt

Les zones de dépôt indiquent à Wails quels éléments doivent accepter les fichiers. Les fichiers déposés en dehors d’une zone de dépôt sont ignorés.

Ajoutez l’attribut `data-file-drop-target` à n’importe quel élément :

```html
<div id="upload" class="drop-zone" data-file-drop-target>
    Drop files here
</div>
```

Vous pouvez définir plusieurs zones de dépôt. L’attribut `id` et les classes CSS de l’élément sont transmis à votre code Go, ce qui vous permet de traiter les dépôts différemment selon l’endroit où les fichiers sont déposés.

## Styliser le survol lors du glissement

Lorsque des fichiers sont glissés au-dessus d’une zone de dépôt, Wails ajoute la classe `file-drop-target-active`. Vous pouvez ainsi fournir un retour visuel pour indiquer aux utilisateurs où ils peuvent déposer les fichiers :

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

La classe est supprimée automatiquement lorsque les fichiers quittent la zone ou y sont déposés.

## Détecter les fichiers déposés

Lorsque des fichiers sont déposés dans une zone de dépôt valide, Wails déclenche un événement `WindowFilesDropped`. Le contexte de l’événement contient les chemins complets, dans le système de fichiers, de tous les fichiers déposés :

```go
import "github.com/wailsapp/wails/v3/pkg/events"

window.OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {
    files := event.Context().DroppedFiles()
    
    for _, file := range files {
        fmt.Println("Dropped:", file)
    }
})
```

Les chemins sont absolus, comme `/home/user/documents/report.pdf` ou `C:\Users\Name\Documents\report.pdf`.

## Obtenir les informations sur la cible du dépôt

Si vous avez plusieurs zones de dépôt, vous pouvez déterminer laquelle a reçu les fichiers à l’aide de `DropTargetDetails()` :

```go
window.OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {
    files := event.Context().DroppedFiles()
    details := event.Context().DropTargetDetails()
    
    fmt.Printf("Dropped on element: id=%s, classes=%v\n", 
        details.ElementID, details.ClassList)
    fmt.Printf("Position: x=%d, y=%d\n", details.X, details.Y)
})
```

Vous pouvez ainsi acheminer les fichiers vers différents gestionnaires :

```go
switch details.ElementID {
case "images":
    handleImageUpload(files)
case "documents":
    handleDocumentUpload(files)
}
```

## Exemple complet

**Go :**

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

**HTML :**

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

## Dépôt dans toute la fenêtre

Si vous souhaitez permettre le dépôt de fichiers n’importe où dans votre application, ajoutez l’attribut à l’élément body :

```html
<body data-file-drop-target>
    <!-- Your app content -->
</body>
```

Vous pouvez utiliser une superposition CSS pour indiquer que toute la fenêtre est une cible de dépôt :

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

## Combiner avec le glisser-déposer HTML

Vous pouvez utiliser à la fois le dépôt de fichiers externes et le glisser-déposer HTML interne dans la même application. Lorsque `EnableFileDrop` vaut `true`, Wails intercepte les glissements de fichiers externes, mais laisse les glissements HTML5 internes s’effectuer normalement.

Pour les distinguer dans les gestionnaires de vos zones de dépôt HTML, vérifiez si le glissement contient des fichiers :

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

Ainsi, vos gestionnaires de dépôt HTML ne répondent qu’aux glissements internes, par exemple lors du déplacement d’éléments d’une liste, tandis que Wails traite séparément les dépôts de fichiers externes au moyen de l’événement `WindowFilesDropped`.

## Étapes suivantes

- [Glisser-déposer HTML](/features/drag-and-drop/html/) — Faites glisser des éléments au sein de votre application
- [Options de fenêtre](/features/windows/options/) — Toutes les options de configuration des fenêtres
