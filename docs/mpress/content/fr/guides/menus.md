---
title: "Menus"
description: "Guide de création et de personnalisation des menus dans Wails v3"
slug: "guides/menus"
sourcePath: "guides/menus.md"
---

Wails v3 fournit un puissant système de menus qui vous permet de créer aussi bien des menus d’application que des menus contextuels. Ce guide vous présente les différentes fonctionnalités et possibilités du système de menus.

## Créer un menu

Pour créer un menu, utilisez la méthode `New()` du gestionnaire Menus :

```go
menu := app.Menu.New()
```

### Ajouter des éléments de menu

Wails prend en charge plusieurs types d’éléments de menu, chacun répondant à un besoin précis :

#### Éléments de menu ordinaires

Les éléments de menu ordinaires constituent les composants de base des menus. Ils affichent du texte et peuvent déclencher des actions lorsque vous cliquez dessus :

```go
menuItem := menu.Add("Click Me")
```

#### Cases à cocher

Les éléments de menu de type case à cocher fournissent un état activable ou désactivable, utile pour activer ou désactiver des fonctionnalités ou des paramètres :

```go
checkbox := menu.AddCheckbox("My checkbox", true)  // true = initially checked
```

#### Groupes de boutons radio

Les groupes de boutons radio permettent aux utilisateurs de sélectionner une option parmi plusieurs choix mutuellement exclusifs. Ils sont créés automatiquement lorsque des éléments de type bouton radio sont placés les uns à côté des autres :

```go
menu.AddRadio("Option 1", true)   // true = initially selected
menu.AddRadio("Option 2", false)
menu.AddRadio("Option 3", false)
```

#### Séparateurs

Les séparateurs sont des lignes horizontales qui permettent d’organiser les éléments de menu en groupes logiques :

```go
menu.AddSeparator()
```

#### Sous-menus

Les sous-menus sont des menus imbriqués qui apparaissent lorsque vous survolez un élément de menu ou cliquez dessus. Ils servent à organiser les structures de menus complexes :

```go
submenu := menu.AddSubmenu("File")
submenu.Add("Open")
submenu.Add("Save")
```

#### Combiner des menus

Vous pouvez ajouter un menu à un autre en l’insérant à la fin ou au début.

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
Par défaut, `prepend` et `append` partagent l’état du menu d’origine. Pour créer un menu doté de son propre état,  vous pouvez appeler `.Clone()` sur le menu.

Ex. : `menu.Append(secondaryMenu.Clone())`

@end

#### Vider un menu

Dans certains cas, si le nombre d’éléments de menu est variable, il peut être préférable de construire un tout nouveau menu.

Cette opération supprime tous les éléments d’un menu existant et vous permet d’en ajouter de nouveau.

```go
menu := app.Menu.New()
menu.Add("Waiting for update...")

// after certain logic, the menu has to be updated
menu.Clear()
menu.Add("Update complete!")
menu.Update()
```

@note{type="info"}
Vider un menu supprime uniquement ses éléments de premier niveau. Les sous-menus ne sont alors plus visibles, mais occupent toujours de la mémoire ;  veillez donc à gérer soigneusement vos menus.

@end

#### Détruire un menu

Pour vider et libérer un menu, utilisez la méthode `Destroy()` :

```go
menu := app.Menu.New()
menu.Add("Waiting for update...")

// after certain logic, the menu has to be destroyed
menu.Destroy()
```

### Propriétés des éléments de menu

Les éléments de menu possèdent plusieurs propriétés configurables :

| Propriété | Méthode | Description |
| --- | --- | --- |
| Libellé | `SetLabel(string)` | Définit le texte affiché |
| Activé | `SetEnabled(bool)` | Active ou désactive l’élément |
| Coché | `SetChecked(bool)` | Définit l’état coché (pour les cases à cocher et les boutons radio) |
| Info-bulle | `SetTooltip(string)` | Définit le texte de l’info-bulle |
| Masqué | `SetHidden(bool)` | Affiche ou masque l’élément |
| Raccourci clavier | `SetAccelerator(string)` | Définit le raccourci clavier |

### États des éléments de menu

Les éléments de menu peuvent avoir différents états qui déterminent leur visibilité et leur interactivité :

#### Visibilité

Vous pouvez afficher ou masquer dynamiquement les éléments de menu à l’aide de la méthode `SetHidden()` :

```go
menuItem := menu.Add("Dynamic Item")

// Hide the menu item
menuItem.SetHidden(true)

// Show the menu item
menuItem.SetHidden(false)

// Check current visibility
isHidden := menuItem.Hidden()
```

Les éléments de menu masqués sont entièrement retirés du menu jusqu’à ce qu’ils soient de nouveau affichés. Cette possibilité est utile pour les éléments contextuels qui ne doivent apparaître que dans certains états de l’application.

#### État d’activation

Vous pouvez activer ou désactiver les éléments de menu à l’aide de la méthode `SetEnabled()` :

```go
menuItem := menu.Add("Save")

// Disable the menu item
menuItem.SetEnabled(false)  // Item appears grayed out and cannot be clicked

// Enable the menu item
menuItem.SetEnabled(true)   // Item becomes clickable again

// Check current enabled state
isEnabled := menuItem.Enabled()
```

Les éléments de menu désactivés restent visibles, mais apparaissent grisés et ne peuvent pas être sélectionnés. Cet état sert généralement à indiquer qu’une action est actuellement indisponible, par exemple :

- Désactiver « Enregistrer » lorsqu’aucune modification n’est à enregistrer
- Désactiver « Copier » lorsqu’aucun élément n’est sélectionné
- Désactiver « Annuler » lorsqu’aucune action ne peut être annulée

#### Gestion dynamique des états

Vous pouvez combiner ces états avec des gestionnaires d’événements pour créer des menus dynamiques :

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

### Gestion des événements

Les éléments de menu peuvent réagir aux événements de clic à l’aide de la méthode `OnClick` :

```go
menuItem.OnClick(func(ctx *application.Context) {
    // Handle the click event
    println("Menu item clicked!")
})
```

Le contexte fournit des informations sur l’élément de menu sélectionné :

```go
menuItem.OnClick(func(ctx *application.Context) {
    // Get the clicked menu item
    clickedItem := ctx.ClickedMenuItem()
    // Get its current state
    isChecked := clickedItem.Checked()
})
```

### Éléments de menu basés sur des rôles

Wails fournit un ensemble de rôles de menu prédéfinis qui créent automatiquement des éléments de menu dotés de fonctionnalités standard. Voici les rôles de menu pris en charge :

#### Structures de menu complètes

Ces rôles créent des structures de menu complètes dotées de fonctionnalités courantes :

| Rôle | Description | Remarques sur les plateformes |
| --- | --- | --- |
| `AppMenu` | Menu de l’application comprenant À propos, Services, Masquer/Afficher et Quitter | macOS uniquement |
| `EditMenu` | Menu Édition standard comprenant Annuler, Rétablir, Couper, Copier, Coller, etc. | Toutes les plateformes |
| `ViewMenu` | Menu Affichage comprenant des commandes pour recharger, zoomer et passer en plein écran | Toutes les plateformes |
| `WindowMenu` | Commandes de la fenêtre (Réduire, Zoomer, etc.) | Toutes les plateformes |
| `HelpMenu` | Menu Aide comprenant un lien « En savoir plus » vers le site web de Wails | Toutes les plateformes |

#### Éléments de menu individuels

Ces rôles permettent d’ajouter des éléments de menu individuels :

| Rôle | Description | Remarques sur les plateformes |
| --- | --- | --- |
| `About` | Afficher la boîte de dialogue À propos de l’application | Toutes les plateformes |
| `Hide` | Masquer l’application | macOS uniquement |
| `HideOthers` | Masquer les autres applications | macOS uniquement |
| `UnHide` | Afficher l’application masquée | macOS uniquement |
| `CloseWindow` | Fermer la fenêtre active | Toutes les plateformes |
| `Minimise` | Réduire la fenêtre | Toutes les plateformes |
| `Zoom` | Zoomer la fenêtre | macOS uniquement |
| `Front` | Placer la fenêtre au premier plan | macOS uniquement |
| `Quit` | Quitter l’application | Toutes les plateformes |
| `Undo` | Annuler la dernière action | Toutes les plateformes |
| `Redo` | Rétablir la dernière action | Toutes les plateformes |
| `Cut` | Couper la sélection | Toutes les plateformes |
| `Copy` | Copier la sélection | Toutes les plateformes |
| `Paste` | Coller depuis le presse-papiers | Toutes les plateformes |
| `PasteAndMatchStyle` | Coller et adapter le style | macOS uniquement |
| `SelectAll` | Tout sélectionner | Toutes les plateformes |
| `Delete` | Supprimer la sélection | Toutes les plateformes |
| `Reload` | Recharger la page actuelle | Toutes les plateformes |
| `ForceReload` | Forcer le rechargement de la page actuelle | Toutes les plateformes |
| `ToggleFullscreen` | Activer ou désactiver le mode plein écran | Toutes les plateformes |
| `ResetZoom` | Réinitialiser le niveau de zoom | Toutes les plateformes |
| `ZoomIn` | Augmenter le zoom | Toutes les plateformes |
| `ZoomOut` | Réduire le zoom | Toutes les plateformes |

Voici un exemple montrant comment utiliser à la fois des menus complets et des rôles individuels :

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

## Menus de l’application

Les menus de l’application apparaissent en haut de la fenêtre de votre application (Windows/Linux) ou en haut de l’écran (macOS).

### Comportement du menu de l’application

Lorsque vous définissez un menu d’application avec `app.Menu.Set()`, il devient le menu principal sous macOS. Sous Windows/Linux, les menus sont définis individuellement pour chaque fenêtre.

```go
app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Custom Menu Window",
    Windows: application.WindowsWindow{
        Menu: customMenu,  // Override application menu for this window
    },
})
```

Voici un exemple complet illustrant ces différents comportements de menu :

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

## Menus contextuels

Les menus contextuels sont des menus surgissants qui apparaissent lorsque vous cliquez avec le bouton droit sur des éléments de votre application. Ils permettent d’accéder rapidement aux actions pertinentes pour l’élément sur lequel vous avez cliqué.

### Menu contextuel par défaut

Le menu contextuel par défaut est le menu contextuel intégré à la vue web. Il propose des opérations au niveau du système, telles que :

- Copier, Couper et Coller pour manipuler du texte
- Commandes de sélection de texte
- Options de vérification orthographique

#### Contrôle du menu contextuel par défaut

Vous pouvez contrôler l’affichage du menu contextuel par défaut à l’aide de la propriété CSS `--default-contextmenu` :

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
Cette fonctionnalité ne fonctionnera comme prévu qu’une fois que l’[environnement d’exécution frontend sera prêt](/reference/frontend-runtime/).

@end

#### Comportement des menus contextuels imbriqués

Lorsque vous utilisez la propriété `--default-contextmenu` sur des éléments imbriqués, les règles suivantes s’appliquent :

1. Les éléments enfants héritent du réglage de menu contextuel de leur parent, sauf s’il est explicitement remplacé
2. Le réglage le plus spécifique, c’est-à-dire le plus proche, est prioritaire
3. La valeur `auto` permet de rétablir le comportement par défaut

Exemple de comportement de menus contextuels imbriqués :

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

### Menus contextuels personnalisés

Les menus contextuels personnalisés vous permettent de proposer des actions propres à l’application et pertinentes pour l’élément sur lequel l’utilisateur clique. Ils sont particulièrement utiles pour :

- Les opérations sur les fichiers dans un gestionnaire de documents
- Les outils de manipulation d’images
- Les actions personnalisées dans une grille de données
- Opérations propres au composant

#### Créer un menu contextuel personnalisé

Lorsque vous créez un menu contextuel personnalisé, fournissez un identifiant unique (nom) qui associe le menu aux éléments HTML :

```go
// Create a context menu with identifier "imageMenu"
contextMenu := app.ContextMenu.New()
app.ContextMenu.Add("imageMenu", contextMenu)
```

Le paramètre de nom (« imageMenu » dans cet exemple) sert d’identifiant unique permettant de :

1. Associer des éléments HTML à ce menu contextuel précis
2. Déterminer le menu à afficher lors d’un clic droit
3. Permettre la mise à jour et la suppression du menu

#### Données de contexte

Lors du traitement des événements d’un menu contextuel, vous pouvez accéder à la fois à l’élément de menu sélectionné et aux données de contexte qui lui sont associées :

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

Les données de contexte proviennent de la propriété `--custom-contextmenu-data` de l’élément HTML et sont accessibles dans le gestionnaire de clic via `ctx.ContextMenuData()`. Cela est particulièrement utile pour :

- Utiliser des listes ou des grilles dans lesquelles chaque élément doit être identifié de manière unique
- Effectuer des opérations sur des composants ou des éléments précis
- Transmettre un état ou des métadonnées du frontend au backend

#### Gestion des menus contextuels

Après avoir modifié un menu contextuel, appelez la méthode `Update()` pour appliquer les modifications :

```go
contextMenu.Update()
```

Lorsque vous n’avez plus besoin d’un menu contextuel, vous pouvez le détruire :

```go
contextMenu.Destroy()
```

@note{type="danger" title="Avertissement"}
Après avoir appelé `Destroy()`, toute nouvelle utilisation de la référence au menu contextuel provoquera une panique.

@end

### Exemple concret : galerie d’images

Voici un exemple complet d’implémentation d’un menu contextuel personnalisé pour une galerie d’images :

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

Dans cet exemple :

1. Le menu contextuel est créé avec l’identifiant « imageMenu »
2. Chaque conteneur d’image est associé au menu à l’aide de `--custom-contextmenu: imageMenu`
3. Chaque conteneur fournit l’identifiant de son image comme données de contexte à l’aide de `--custom-contextmenu-data`
4. Le backend reçoit l’identifiant de l’image dans les gestionnaires de clic et peut effectuer des opérations spécifiques
5. Le même menu est réutilisé pour toutes les images, mais les données de contexte indiquent l’image sur laquelle effectuer l’opération

Ce modèle est particulièrement efficace pour :

- Les grilles de données dont les lignes nécessitent des opérations spécifiques
- Les gestionnaires de fichiers dans lesquels les fichiers nécessitent des actions propres au contexte
- Les outils de conception dans lesquels différents éléments nécessitent différentes opérations
- Tout composant dont plusieurs instances partagent les mêmes opérations
