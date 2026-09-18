---
title: "Référence des menus"
description: "Référence complète des types, propriétés et méthodes des éléments de menu"
slug: "features/menus/reference"
sourcePath: "features/menus/reference.md"
---

## Référence des menus

Référence complète des types, propriétés et comportements dynamiques des éléments de menu. Créez des menus professionnels et réactifs avec des cases à cocher, des groupes de boutons radio, des séparateurs et des mises à jour dynamiques.

## Types d’éléments de menu

### Éléments de menu standard

Le type le plus courant : il affiche du texte et déclenche une action :

```go
menuItem := menu.Add("Click Me")
menuItem.OnClick(func(ctx *application.Context) {
    fmt.Println("Menu item clicked!")
})
```

**Utilisation :** commandes, actions, ouverture de fenêtres

### Cases à cocher

Éléments de menu basculables entre les états coché et décoché :

```go
checkbox := menu.AddCheckbox("Enable Feature", true)  // true = initially checked
checkbox.OnClick(func(ctx *application.Context) {
    isChecked := ctx.ClickedMenuItem().Checked()
    fmt.Printf("Feature is now: %v\n", isChecked)
})
```

**Utilisation :** paramètres booléens, activation ou désactivation de fonctionnalités, options d’affichage

**Important :** l’état coché bascule automatiquement lors d’un clic.

### Groupes de boutons radio

Options mutuellement exclusives : une seule peut être sélectionnée :

```go
menu.AddRadio("Small", true)   // true = initially selected
menu.AddRadio("Medium", false)
menu.AddRadio("Large", false)
```

**Utilisation :** choix mutuellement exclusifs (taille, thème, mode)

**Fonctionnement du regroupement :**

- Les éléments radio adjacents forment automatiquement un groupe
- La sélection d’un élément désélectionne les autres éléments du groupe
- Séparez les groupes avec un séparateur ou un élément standard

**Exemple avec plusieurs groupes :**

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

### Sous-menus

Structures de menus imbriquées pour faciliter l’organisation :

```go
submenu := menu.AddSubmenu("More Options")
submenu.Add("Submenu Item 1").OnClick(func(ctx *application.Context) {
    // Handle click
})
submenu.Add("Submenu Item 2")
```

**Utilisation :** regroupement d’éléments associés, réduction de l’encombrement

**Limite d’imbrication :** la plupart des plateformes prennent en charge 2-3 niveaux. Évitez une imbrication plus profonde.

### Séparateurs

Séparations visuelles entre les éléments de menu :

```go
menu.Add("Item 1")
menu.AddSeparator()
menu.Add("Item 2")
```

**Utilisation :** regroupement visuel d’éléments associés

**Bonne pratique :** ne placez pas de séparateur au début ou à la fin d’un menu.

## Propriétés des éléments de menu

### Libellé

Texte affiché pour l’élément de menu :

```go
menuItem := menu.Add("Initial Label")
menuItem.SetLabel("New Label")

// Get current label
label := menuItem.Label()
```

**Libellés dynamiques :**

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

### État activé

Déterminez si l’utilisateur peut interagir avec l’élément de menu :

```go
menuItem := menu.Add("Save")
menuItem.SetEnabled(false)  // Greyed out, can't click

// Enable it later
menuItem.SetEnabled(true)
menu.Update()  // Important: Call this after changing enabled state!

// Check current state
isEnabled := menuItem.Enabled()
```

@note{type="caution" title="Comportement des menus sous Windows"}
Sous Windows, les menus doivent être reconstruits lorsque leur état change. **Appelez toujours `menu.Update()` après avoir activé ou désactivé des éléments de menu**, en particulier si l’élément a été créé alors qu’il était désactivé.

**Pourquoi :** sous Windows, les menus sont entièrement reconstruits lors de leur mise à jour. Si vous n’appelez pas `Update()`, les gestionnaires de clic ne se déclencheront pas correctement.

@end

**Exemple : activation et désactivation dynamiques**

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

**Schéma courant : activation sous condition**

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

### État coché

Pour les cases à cocher et les éléments radio, contrôlez ou consultez leur état coché :

```go
checkbox := menu.AddCheckbox("Feature", false)
checkbox.SetChecked(true)
menu.Update()

// Query state
isChecked := checkbox.Checked()
```

**Basculement automatique :** les cases à cocher basculent automatiquement lors d’un clic. Il n’est pas nécessaire d’appeler `SetChecked()` dans le gestionnaire de clic.

**Contrôle manuel :**

```go
checkbox := menu.AddCheckbox("Auto-save", false)

// Sync with external state
func syncAutoSave(enabled bool) {
    checkbox.SetChecked(enabled)
    menu.Update()
}
```

### Accélérateurs (raccourcis clavier)

Ajoutez des raccourcis clavier aux éléments de menu :

```go
saveMenuItem := menu.Add("Save")
saveMenuItem.SetAccelerator("CmdOrCtrl+S")

quitMenuItem := menu.Add("Quit")
quitMenuItem.SetAccelerator("CmdOrCtrl+Q")
```

**Format des accélérateurs :**

- `CmdOrCtrl` — Cmd sous macOS, Ctrl sous Windows/Linux
- `Shift`, `Alt`, `Option` — touches de modification
- `A-Z`, `0-9` — touches de lettres ou de chiffres
- `F1-F12` — touches de fonction
- `Enter`, `Space`, `Backspace`, etc. — touches spéciales

**Exemples :**

```go
"CmdOrCtrl+S"           // Save
"CmdOrCtrl+Shift+S"     // Save As
"CmdOrCtrl+W"           // Close Window
"CmdOrCtrl+Q"           // Quit
"F5"                    // Refresh
"CmdOrCtrl+,"           // Preferences (macOS convention)
"Alt+F4"                // Close (Windows convention)
```

**Accélérateurs propres à chaque plateforme :**

```go
if runtime.GOOS == "darwin" {
    prefsMenuItem.SetAccelerator("Cmd+,")
} else {
    prefsMenuItem.SetAccelerator("Ctrl+P")
}
```

### Infobulle

Ajoutez du texte qui s’affiche au survol des éléments de menu (la prise en charge varie selon la plateforme) :

```go
menuItem := menu.Add("Advanced Options")
menuItem.SetTooltip("Configure advanced settings")
```

**Prise en charge par plateforme :**

- **Windows :** ✅ pris en charge
- **macOS :** ❌ non pris en charge (les infobulles ne sont pas standard dans les menus)
- **Linux :** ⚠️ varie selon l’environnement de bureau

### État masqué

Masquez des éléments de menu sans les supprimer :

```go
debugMenuItem := menu.Add("Debug Mode")
debugMenuItem.SetHidden(true)  // Hidden

// Show in debug builds
if isDebugBuild {
    debugMenuItem.SetHidden(false)
    menu.Update()
}
```

**Utilisation :** options de débogage, indicateurs de fonctionnalités, fonctionnalités conditionnelles

## Gestion des événements

### Gestionnaire OnClick

Exécutez du code lorsque l’utilisateur clique sur un élément de menu :

```go
menuItem := menu.Add("Click Me")
menuItem.OnClick(func(ctx *application.Context) {
    // Handle click
    fmt.Println("Clicked!")
})
```

**Le contexte fournit :**

- `ctx.ClickedMenuItem()` — L’élément de menu sur lequel l’utilisateur a cliqué
- Contexte de la fenêtre (si l’action provient du menu de la fenêtre)
- Contexte de l’application

**Exemple : accéder à l’élément de menu dans le gestionnaire**

```go
checkbox := menu.AddCheckbox("Feature", false)
checkbox.OnClick(func(ctx *application.Context) {
    item := ctx.ClickedMenuItem()
    isChecked := item.Checked()
    fmt.Printf("Feature is now: %v\n", isChecked)
})
```

### Plusieurs gestionnaires

Vous pouvez définir plusieurs gestionnaires (le dernier prévaut) :

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

**Bonne pratique :** définissez le gestionnaire une seule fois et, si nécessaire, utilisez une logique conditionnelle à l’intérieur.

## Menus dynamiques

### Mise à jour des éléments de menu

**Règle d’or :** appelez toujours `menu.Update()` après avoir modifié l’état du menu.

```go
// ✅ Correct
menuItem.SetEnabled(true)
menu.Update()

// ❌ Wrong (especially on Windows)
menuItem.SetEnabled(true)
// Forgot to call Update() - click handlers may not work!
```

**Pourquoi est-ce important ?**

- **Windows :** les menus sont reconstruits lors de leur mise à jour
- **macOS/Linux :** moins indispensable, mais toujours recommandé
- **Gestionnaires de clic :** ils ne se déclenchent pas correctement sans Update()

### Reconstruction des menus

Pour les modifications importantes, reconstruisez l’intégralité du menu :

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

**Quand reconstruire le menu ?**

- La liste des fichiers récents change
- Les menus des plugins changent
- Transitions d’état majeures

**Quand mettre le menu à jour ?**

- Activer ou désactiver des éléments
- Modifier les libellés
- Cocher ou décocher des cases

### Menus contextuels

Adaptez les menus à l’état de l’application :

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

## Différences entre les plateformes

### Emplacement de la barre de menus

| Plateforme | Emplacement | Remarques |
| --- | --- | --- |
| **macOS** | En haut de l’écran | Barre de menus globale |
| **Windows** | En haut de la fenêtre | Menu propre à chaque fenêtre |
| **Linux** | En haut de la fenêtre | Propre à chaque fenêtre (généralement) |

### Menus standard

**macOS :**

- Comporte un menu « Application » (portant le nom de l’application)
- « Préférences » se trouve dans le menu Application
- « Quitter » se trouve dans le menu Application

**Windows/Linux :**

- Aucun menu Application
- « Préférences » se trouve dans le menu Édition ou Outils
- « Quitter » se trouve dans le menu Fichier

**Exemple : structure adaptée à la plateforme**

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

### Conventions relatives aux raccourcis clavier

**macOS :**

- `Cmd+` pour la plupart des raccourcis
- `Cmd+,` pour Préférences
- `Cmd+Q` pour Quitter

**Windows :**

- `Ctrl+` pour la plupart des raccourcis
- `Ctrl+P` ou `Ctrl+,` pour Préférences
- `Alt+F4` pour Quitter (ou `Ctrl+Q`)

**Linux :**

- Suit généralement les conventions de Windows
- L’environnement de bureau peut les remplacer

## Bonnes pratiques

### ✅ À faire

- **Appelez menu.Update()** après avoir modifié l’état du menu (en particulier sous Windows)
- **Utilisez des groupes de boutons radio** pour les options mutuellement exclusives
- **Utilisez des cases à cocher** pour les fonctionnalités activables et désactivables
- **Ajoutez des raccourcis clavier** aux actions courantes
- **Regroupez les éléments associés** à l’aide de séparateurs
- **Testez sur toutes les plateformes** — le comportement varie

### ❌ À éviter

- **N’oubliez pas menu.Update()** — les gestionnaires de clic ne fonctionneront pas correctement
- **N’imbriquez pas trop profondément** — 2-3 niveaux au maximum
- **Ne commencez et ne terminez pas par des séparateurs** — cela manque de professionnalisme
- **N’utilisez pas d’infobulles sous macOS** — elles ne sont pas prises en charge
- **Ne codez pas en dur les raccourcis propres à chaque plateforme** — utilisez `CmdOrCtrl`

## Dépannage

### Les éléments de menu ne répondent pas

**Symptôme :** les gestionnaires de clic ne se déclenchent pas

**Cause :** oubli de l’appel à `menu.Update()` après l’activation de l’élément

**Solution :**

```go
menuItem.SetEnabled(true)
menu.Update()  // Add this!
```

### Les éléments de menu sont grisés

**Symptôme :** impossible de cliquer sur les éléments de menu

**Cause :** les éléments sont désactivés

**Solution :**

```go
menuItem.SetEnabled(true)
menu.Update()
```

### Les raccourcis clavier ne fonctionnent pas

**Symptôme :** les raccourcis clavier ne déclenchent pas les éléments de menu

**Causes :**

1. Format du raccourci clavier incorrect
2. Conflit avec les raccourcis système
3. La fenêtre n’a pas le focus

**Solution :**

```go
// Check format
menuItem.SetAccelerator("CmdOrCtrl+S")  // ✅ Correct
menuItem.SetAccelerator("Ctrl+S")       // ❌ Wrong (macOS uses Cmd)

// Avoid conflicts
// ❌ Cmd+H (Hide Window on macOS - system shortcut)
// ✅ Cmd+Shift+H (Custom shortcut)
```

## Étapes suivantes

- [Menus d’application](/features/menus/application/) — créez des barres de menus d’application
- [Menus contextuels](/features/menus/context/) — créez des menus contextuels accessibles par clic droit
- [Menus de la zone de notification](/features/menus/systray/) — créez des menus pour la zone de notification ou la barre des menus
- [Modèles de menus](/guides/menus/) — découvrez les modèles de menus courants et les bonnes pratiques

---

**Des questions ?** Posez-les sur [Discord](https://discord.gg/JDdSxwjhGf) ou consultez les [exemples de menus](https://github.com/wailsapp/wails/tree/master/v3/examples/menu).
