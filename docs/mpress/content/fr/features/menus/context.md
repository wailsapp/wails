---
title: "Menus contextuels"
description: "Créez des menus contextuels accessibles par clic droit pour votre application"
slug: "features/menus/context"
sourcePath: "features/menus/context.md"
---

## Le problème

Les utilisateurs s’attendent à disposer de menus accessibles par clic droit proposant des actions adaptées au contexte. Chaque type d’élément nécessite un menu différent :

- **Texte** : Couper, Copier, Coller
- **Images** : Enregistrer, Copier, Ouvrir
- **Éléments personnalisés** : actions propres à l’application

Créer manuellement des menus contextuels implique de gérer les événements de la souris, le positionnement et les différences entre plateformes.

## La solution Wails

Wails fournit des **menus contextuels déclaratifs** au moyen de propriétés CSS. Associez des menus à des éléments HTML, transmettez des données et gérez les clics, tout en bénéficiant du comportement natif de la plateforme.

![Un menu contextuel personnalisé Wails affiché au-dessus d’une WebView sous macOS](/assets/screenshots/context-menu-macos.png)

Le menu est natif à la plateforme, tandis que l’élément qui l’a ouvert reste intégré à votre WebView. Cette capture réalisée sous macOS utilise l’enregistrement du menu contextuel personnalisé de l’exemple ci-dessous.

## Démarrage rapide

**Code Go :**

```go
// Create context menu
contextMenu := app.ContextMenu.New()
contextMenu.Add("Cut").OnClick(handleCut)
contextMenu.Add("Copy").OnClick(handleCopy)
contextMenu.Add("Paste").OnClick(handlePaste)

// Register with ID
app.ContextMenu.Add("editor-menu", contextMenu)
```

**HTML :**

```html
<textarea style="--custom-contextmenu: editor-menu">
    Right-click me!
</textarea>
```

**C’est tout !** Un clic droit dans la zone de texte affiche votre menu personnalisé.

## Création de menus contextuels

### Menu contextuel de base

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

**Identifiant du menu :** il doit être unique. Il sert à associer le menu aux éléments HTML.

### Avec des sous-menus

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

### Avec des cases à cocher et des groupes de boutons radio

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

**Pour tous les types d’éléments de menu**, consultez la [Référence des menus](/features/menus/reference/).

## Association à des éléments HTML

Utilisez des propriétés personnalisées CSS pour associer des menus contextuels :

### Association de base

```html
<div style="--custom-contextmenu: menu-id">
    Right-click me!
</div>
```

**Propriété CSS :** `--custom-contextmenu: <menu-id>`

### Avec des données contextuelles

Transmettez des données du HTML vers Go :

```html
<div style="--custom-contextmenu: file-menu; --custom-contextmenu-data: file-123">
    Right-click this file
</div>
```

**Gestionnaire Go :**

```go
contextMenu := app.ContextMenu.New()
contextMenu.Add("Open").OnClick(func(ctx *application.Context) {
    fileID := ctx.ContextMenuData()  // "file-123"
    openFile(fileID)
})

app.ContextMenu.Add("file-menu", contextMenu)
```

**Propriétés CSS :**

- `--custom-contextmenu: <menu-id>` - Menu à afficher
- `--custom-contextmenu-data: <data>` - Données à transmettre aux gestionnaires

### Données dynamiques

Générez les données dynamiquement en JavaScript :

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

### Plusieurs éléments, un même menu

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

**Un seul menu, avec des données différentes pour chaque élément.**

## Données contextuelles

### Accès aux données contextuelles

```go
contextMenu.Add("Process").OnClick(func(ctx *application.Context) {
    data := ctx.ContextMenuData()  // Get data from HTML
    
    // Use the data
    processItem(data)
})
```

**Type de données :** toujours `string`. Analysez-les selon vos besoins.

### Transmission de données complexes

Utilisez JSON pour les données complexes :

```html
<div style="--custom-contextmenu: item-menu; --custom-contextmenu-data: {&quot;id&quot;:123,&quot;type&quot;:&quot;image&quot;}">
    Image.png
</div>
```

**Gestionnaire Go :**

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

@note{type="caution" title="Sécurité"}
**Validez toujours les données contextuelles** provenant du frontend. Les utilisateurs peuvent modifier les propriétés CSS ; traitez donc ces données comme des entrées non fiables.

@end

### Exemple de validation

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

## Menu contextuel par défaut

La WebView fournit un menu contextuel intégré pour les opérations standard (copier, coller, inspecter). Contrôlez-le avec `--default-contextmenu` :

### Masquer le menu par défaut

```html
<div style="--default-contextmenu: hide">
    No default menu here
</div>
```

**Cas d’utilisation :** éléments d’interface utilisateur personnalisés pour lesquels le menu par défaut n’est pas pertinent.

### Afficher le menu par défaut

```html
<div style="--default-contextmenu: show">
    Default menu always shown
</div>
```

**Cas d’utilisation :** zones de texte, champs de saisie et contenu modifiable.

### Mode automatique (intelligent)

```html
<div style="--default-contextmenu: auto">
    Smart context menu
</div>
```

**Comportement par défaut.** Affiche le menu par défaut dans les cas suivants :

- Du texte est sélectionné
- Dans les champs de saisie de texte
- Dans le contenu modifiable (`contenteditable`)

Dans les autres cas, masque le menu par défaut.

### Combinaison des menus personnalisé et par défaut

```html
<!-- Custom menu + default menu -->
<textarea style="--custom-contextmenu: editor-menu; --default-contextmenu: show">
    Both menus available
</textarea>
```

**Comportement :**

1. Le menu personnalisé s’affiche en premier
2. Si le menu personnalisé est vide ou introuvable, le menu par défaut s’affiche
3. Les deux peuvent coexister (selon la plateforme)

## Menus contextuels dynamiques

Mettez à jour les menus en fonction de l’état de l’application :

### Activation et désactivation des éléments

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

@note{type="caution" title="Toujours appeler Update()"}
Après avoir modifié l’état du menu, **appelez `contextMenu.Update()`**. Cette opération est indispensable sous Windows.

Consultez la [Référence des menus](/features/menus/reference/#enabled-state) pour plus de détails.

@end

### Modifier les libellés

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

### Reconstruire les menus

Pour apporter des modifications importantes, reconstruisez l’intégralité du menu :

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

## Comportement selon la plateforme

Les menus contextuels sont **natifs de chaque plateforme** :

@tabs{sync-key="platform"}
[macOS]
**Menus contextuels natifs de macOS :**

- Animations et transitions du système
- Clic droit = Contrôle+clic (automatique)
- Adaptation à l’apparence du système (claire/sombre)
- Opérations de texte standard dans le menu par défaut
- Défilement natif pour les menus longs

**Conventions de macOS :**

- Utilisez la casse de phrase pour les éléments de menu
- Utilisez des points de suspension (...) pour les éléments qui ouvrent des boîtes de dialogue
- Raccourcis courants : ⌘C (Copier), ⌘V (Coller)

[Windows]
**Menus contextuels natifs de Windows :**

- Style natif de Windows
- Respect du thème Windows
- Opérations Windows standard dans le menu par défaut
- Prise en charge de la saisie tactile et au stylet

**Conventions de Windows :**

- Utilisez une majuscule initiale pour les mots principaux des éléments de menu
- Utilisez des points de suspension (...) pour les éléments qui ouvrent des boîtes de dialogue
- Raccourcis courants : Ctrl+C (Copier), Ctrl+V (Coller)

[Linux]
**Intégration à l’environnement de bureau :**

- Adaptation au thème du bureau (GTK, Qt, etc.)
- Comportement du clic droit conforme aux paramètres système
- Contenu du menu par défaut variable selon l’environnement
- Positionnement conforme aux conventions de l’environnement de bureau

**Points à prendre en compte sous Linux :**

- Effectuez des tests dans les environnements de bureau ciblés
- GTK et Qt ont des comportements différents
- Certains environnements de bureau personnalisent les menus contextuels

@end

## Exemple complet

**Code Go :**

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

**HTML :**

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

## Bonnes pratiques

### ✅ À faire

- **Gardez des menus ciblés** — Proposez uniquement les actions pertinentes pour l’élément
- **Validez les données de contexte** — Traitez-les comme des entrées non fiables
- **Utilisez des libellés explicites** — « Supprimer le fichier » plutôt que « Supprimer »
- **Appelez menu.Update()** — Après avoir modifié l’état du menu
- **Testez sur toutes les plateformes** — Le comportement varie
- **Fournissez des raccourcis clavier** — Pour les actions courantes
- **Regroupez les éléments associés** — Utilisez des séparateurs

### ❌ À ne pas faire

- **Ne faites pas confiance aux données de contexte** — Validez-les systématiquement
- **Ne créez pas de menus trop longs** — 7-10 éléments au maximum
- **N’oubliez pas menu.Update()** — Les menus ne fonctionneront pas correctement
- **N’imbriquez pas les menus trop profondément** — 2 niveaux au maximum
- **N’utilisez pas de jargon** — Gardez des libellés faciles à comprendre
- **Ne bloquez pas les gestionnaires** — Veillez à ce qu’ils s’exécutent rapidement

## Résolution des problèmes

### Le menu contextuel ne s’affiche pas

**Causes possibles :**

1. L’identifiant du menu ne correspond pas
2. Faute de frappe dans la propriété CSS
3. Runtime non initialisé

**Solution :**

```go
// Check menu is registered
app.ContextMenu.Add("my-menu", contextMenu)
```

```html
<!-- Check ID matches -->
<div style="--custom-contextmenu: my-menu">
```

### Les données de contexte ne sont pas reçues

**Causes possibles :**

1. Propriété CSS non définie
2. Les données contiennent des caractères spéciaux

**Solution :**

```html
<!-- Escape quotes in JSON -->
<div style="--custom-contextmenu-data: {&quot;id&quot;:123}">
```

Vous pouvez également utiliser JavaScript :

```javascript
element.style.setProperty('--custom-contextmenu-data', JSON.stringify(data))
```

### Les éléments de menu ne répondent pas

**Cause :** oubli de l’appel à `menu.Update()` après l’activation

**Solution :**

```go
menuItem.SetEnabled(true)
contextMenu.Update()  // Add this!
```

## Étapes suivantes

@cards{cols="2"}
📖 Référence des menus
Référence complète des types et des propriétés des éléments de menu.

[En savoir plus →](/features/menus/reference/)

---
☰ Menus de l’application
Créez des barres de menus pour l’application.

[En savoir plus →](/features/menus/application/)

---
★ Menus de la zone de notification
Ajoutez l’intégration à la zone de notification ou à la barre des menus.

[En savoir plus →](/features/menus/systray/)

---
📖 Modèles de menus
Modèles de menus courants et bonnes pratiques.

[En savoir plus →](/guides/menus/)

@end

---

**Des questions ?** Posez-les sur [Discord](https://discord.gg/JDdSxwjhGf) ou consultez l’[exemple de menu contextuel](https://github.com/wailsapp/wails/tree/master/v3/examples/contextmenus).
