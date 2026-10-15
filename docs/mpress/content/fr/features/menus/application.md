---
title: "Menus de l’application"
description: "Créez des barres de menus natives pour votre application de bureau"
slug: "features/menus/application"
sourcePath: "features/menus/application.md"
---

## Le problème

Les applications de bureau professionnelles ont besoin de barres de menus — Fichier, Édition, Présentation, Aide. Toutefois, les menus fonctionnent différemment sur chaque plateforme :

- **macOS** : barre de menus globale en haut de l’écran
- **Windows** : barre de menus dans la barre de titre de la fenêtre
- **Linux** : varie selon l’environnement de bureau

Créer manuellement des menus adaptés à chaque plateforme est fastidieux et source d’erreurs.

## La solution Wails

Wails fournit une **API unifiée** qui crée automatiquement des menus natifs adaptés à chaque plateforme. Écrivez le code une seule fois et bénéficiez d’un comportement natif sur toutes les plateformes.

![Menu d’une application Wails sous macOS avec des éléments standard, des cases à cocher, des boutons radio et des sous-menus](/assets/screenshots/application-menu-macos.png)

Sous macOS, le menu de l’application se trouve dans la barre de menus globale. Cette capture montre le menu natif rendu par l’API de menus Wails, avec des éléments désactivés, des cases à cocher, des boutons radio et des sous-menus.

## Démarrage rapide

```go
package main

import (
    "runtime"
    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "My App",
    })

    // Create menu
    menu := app.NewMenu()

    // Add standard menus (platform-appropriate)
    if runtime.GOOS == "darwin" {
        menu.AddRole(application.AppMenu)  // macOS only
    }
    menu.AddRole(application.FileMenu)
    menu.AddRole(application.EditMenu)
    menu.AddRole(application.WindowMenu)
    menu.AddRole(application.HelpMenu)

    // Set the application menu
    app.Menu.Set(menu)

    // Create window with UseApplicationMenu to inherit the menu on Windows/Linux
    app.Window.NewWithOptions(application.WebviewWindowOptions{
        UseApplicationMenu: true,
    })

    app.Run()
}
```

**C’est tout !** Vous disposez maintenant de menus natifs adaptés à chaque plateforme, avec des éléments standard. L’option `UseApplicationMenu` garantit que les fenêtres Windows et Linux affichent le menu sans code supplémentaire.

## Création de menus

### Création d’un menu de base

```go
// Create a new menu
menu := app.NewMenu()

// Add a top-level menu
fileMenu := menu.AddSubmenu("File")

// Add menu items
fileMenu.Add("New").OnClick(func(ctx *application.Context) {
    // Handle New
})

fileMenu.Add("Open").OnClick(func(ctx *application.Context) {
    // Handle Open
})

fileMenu.AddSeparator()

fileMenu.Add("Quit").OnClick(func(ctx *application.Context) {
    app.Quit()
})
```

### Définition du menu

**Approche recommandée** — Utilisez `UseApplicationMenu` pour garantir un comportement cohérent sur toutes les plateformes :

```go
// Set the application menu once
app.Menu.Set(menu)

// Create windows that inherit the menu on Windows/Linux
app.Window.NewWithOptions(application.WebviewWindowOptions{
    UseApplicationMenu: true,  // Window uses the app menu
})
```

Cette approche produit le comportement suivant :

- Sous **macOS** : le menu apparaît en haut de l’écran (comportement standard)
- Sous **Windows/Linux** : chaque fenêtre avec `UseApplicationMenu: true` affiche le menu de l’application

**Détails propres à chaque plateforme :**

@tabs{sync-key="platform"}
[macOS]
**Barre de menus globale** (une par application) :

```go
app.Menu.Set(menu)
```

Le menu apparaît en haut de l’écran et reste affiché même lorsque toutes les fenêtres sont fermées. L’option `UseApplicationMenu` n’a aucun effet sous macOS, car toutes les applications utilisent le menu global.

[Windows]
**Barre de menus propre à chaque fenêtre** :

```go
// Option 1: Use application menu (recommended)
app.Menu.Set(menu)
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    UseApplicationMenu: true,
})

// Option 2: Set menu directly on window
window.SetMenu(menu)
```

Chaque fenêtre peut avoir son propre menu ou hériter du menu de l’application. Le menu apparaît dans la barre de titre de la fenêtre.

[Linux]
**Barre de menus propre à chaque fenêtre** (généralement) :

```go
// Option 1: Use application menu (recommended)
app.Menu.Set(menu)
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    UseApplicationMenu: true,
})

// Option 2: Set menu directly on window
window.SetMenu(menu)
```

Le comportement varie selon l’environnement de bureau. Certains, comme Unity, prennent en charge les menus globaux.

@end

@note{type="tip" title="Simplifiez les menus multiplateformes"}
L’utilisation de `UseApplicationMenu: true` évite d’écrire du code propre à chaque plateforme, comme :

```go
// Old approach - no longer needed
if runtime.GOOS == "darwin" {
    app.Menu.Set(menu)
} else {
    window.SetMenu(menu)
}
```

@end

**Menus personnalisés propres à chaque fenêtre :**

Si une fenêtre nécessite un menu différent de celui de l’application, définissez-le directement :

```go
window.SetMenu(customMenu)  // Overrides UseApplicationMenu
```

## Rôles de menu

Wails fournit des **rôles de menu prédéfinis** qui créent automatiquement des structures de menus adaptées à chaque plateforme.

### Rôles disponibles

| Rôle | Description | Remarques sur les plateformes |
| --- | --- | --- |
| `AppMenu` | Menu de l’application avec À propos, Préférences et Quitter | **macOS uniquement** |
| `FileMenu` | Opérations sur les fichiers (Nouveau, Ouvrir, Enregistrer, etc.) | Toutes les plateformes |
| `EditMenu` | Édition de texte (Annuler, Rétablir, Couper, Copier, Coller) | Toutes les plateformes |
| `WindowMenu` | Gestion des fenêtres (Réduire, Zoom, etc.) | Toutes les plateformes |
| `HelpMenu` | Aide et informations | Toutes les plateformes |

### Utilisation des rôles

```go
menu := app.NewMenu()

// macOS: Add application menu
if runtime.GOOS == "darwin" {
    menu.AddRole(application.AppMenu)
}

// All platforms: Add standard menus
menu.AddRole(application.FileMenu)
menu.AddRole(application.EditMenu)
menu.AddRole(application.WindowMenu)
menu.AddRole(application.HelpMenu)
```

**Résultat obtenu :**

@tabs{sync-key="platform"}
[macOS]
**AppMenu** (avec le nom de l’application) :

- À propos de [Nom de l’application]
- Préférences… (⌘,)
- ---
- Services
- ---
- Masquer [Nom de l’application] (⌘H)
- Masquer les autres (⌥⌘H)
- Tout afficher
- ---
- Quitter [Nom de l’application] (⌘Q)

**FileMenu** :

- Nouveau (⌘N)
- Ouvrir… (⌘O)
- ---
- Fermer la fenêtre (⌘W)

**EditMenu** :

- Annuler (⌘Z)
- Rétablir (⇧⌘Z)
- ---
- Couper (⌘X)
- Copier (⌘C)
- Coller (⌘V)
- Tout sélectionner (⌘A)

**WindowMenu** :

- Réduire (⌘M)
- Zoom
- ---
- Tout ramener au premier plan

**HelpMenu** :

- Aide de [Nom de l’application]

[Windows]
**FileMenu** :

- Nouveau (Ctrl+N)
- Ouvrir... (Ctrl+O)
- ---
- Quitter (Alt+F4)

**EditMenu** :

- Annuler (Ctrl+Z)
- Rétablir (Ctrl+Y)
- ---
- Couper (Ctrl+X)
- Copier (Ctrl+C)
- Coller (Ctrl+V)
- Tout sélectionner (Ctrl+A)

**WindowMenu** :

- Réduire
- Agrandir

**HelpMenu** :

- À propos de [Nom de l’application]

[Linux]
Similaire à Windows, mais les raccourcis clavier peuvent varier selon l’environnement de bureau.

@end

### Personnaliser les menus de rôle

`Menu.AddRole(role)` renvoie le menu **récepteur** (le menu de premier niveau), et **non** le sous-menu du rôle. Pour ajouter des éléments au sous-menu du rôle, recherchez l’élément de rôle inséré avec `FindByRole`, puis appelez `GetSubmenu()` sur celui-ci :

```go
menu.AddRole(application.FileMenu)

fileMenu := menu.FindByRole(application.FileMenu).GetSubmenu()
fileMenu.Add("Import...").OnClick(handleImport)
fileMenu.Add("Export...").OnClick(handleExport)
```

## Menus personnalisés

Créez vos propres menus pour les fonctionnalités propres à l’application :

```go
// Add a custom top-level menu
toolsMenu := menu.AddSubmenu("Tools")

// Add items
toolsMenu.Add("Settings").OnClick(func(ctx *application.Context) {
    showSettingsWindow()
})

toolsMenu.AddSeparator()

// Add checkbox
toolsMenu.AddCheckbox("Dark Mode", false).OnClick(func(ctx *application.Context) {
    isDark := ctx.ClickedMenuItem().Checked()
    setTheme(isDark)
})

// Add radio group
toolsMenu.AddRadio("Small", true).OnClick(handleFontSize)
toolsMenu.AddRadio("Medium", false).OnClick(handleFontSize)
toolsMenu.AddRadio("Large", false).OnClick(handleFontSize)

// Add submenu
advancedMenu := toolsMenu.AddSubmenu("Advanced")
advancedMenu.Add("Configure...").OnClick(showAdvancedSettings)
```

**Pour découvrir d’autres types d’éléments de menu**, consultez la [référence des menus](/features/menus/reference/).

## Menus dynamiques

Mettez à jour les menus en fonction de l’état de l’application :

### Activer ou désactiver des éléments

```go
var saveMenuItem *application.MenuItem

func createMenu() {
    menu := app.NewMenu()
    fileMenu := menu.AddSubmenu("File")
    
    saveMenuItem = fileMenu.Add("Save")
    saveMenuItem.SetEnabled(false)  // Initially disabled
    saveMenuItem.OnClick(handleSave)
    
    app.Menu.Set(menu)
}

func onDocumentChanged() {
    saveMenuItem.SetEnabled(hasUnsavedChanges())
    menu.Update()  // Important!
}
```

@note{type="caution" title="Toujours appeler menu.Update()"}
Après avoir modifié l’état du menu (activation/désactivation, libellé ou état coché), **appelez toujours `menu.Update()`**. Cette étape est particulièrement importante sous Windows, où les menus sont reconstruits.

Pour plus de détails, consultez la [référence des menus](/features/menus/reference/#enabled-state).

@end

### Modifier les libellés

```go
updateMenuItem := menu.Add("Check for Updates")

updateMenuItem.OnClick(func(ctx *application.Context) {
    updateMenuItem.SetLabel("Checking...")
    menu.Update()
    
    checkForUpdates()
    
    updateMenuItem.SetLabel("Check for Updates")
    menu.Update()
})
```

### Reconstruire les menus

Pour des modifications importantes, reconstruisez l’intégralité du menu :

```go
func rebuildFileMenu() {
    menu := app.NewMenu()
    fileMenu := menu.AddSubmenu("File")
    
    fileMenu.Add("New").OnClick(handleNew)
    fileMenu.Add("Open").OnClick(handleOpen)
    
    // Add recent files dynamically
    if hasRecentFiles() {
        recentMenu := fileMenu.AddSubmenu("Open Recent")
        for _, file := range getRecentFiles() {
            filePath := file  // Capture for closure
            recentMenu.Add(filepath.Base(file)).OnClick(func(ctx *application.Context) {
                openFile(filePath)
            })
        }
        recentMenu.AddSeparator()
        recentMenu.Add("Clear Recent").OnClick(clearRecentFiles)
    }
    
    fileMenu.AddSeparator()
    fileMenu.Add("Quit").OnClick(func(ctx *application.Context) {
        app.Quit()
    })
    
    app.Menu.Set(menu)
}
```

## Contrôler les fenêtres depuis les menus

Les éléments de menu peuvent contrôler les fenêtres :

```go
viewMenu := menu.AddSubmenu("View")

// Toggle fullscreen
viewMenu.Add("Toggle Fullscreen").OnClick(func(ctx *application.Context) {
    if window, ok := app.Window.GetByName("main"); ok {
        window.ToggleFullscreen()
    }
})

// Zoom controls
viewMenu.Add("Zoom In").SetAccelerator("CmdOrCtrl++").OnClick(func(ctx *application.Context) {
    // Increase zoom
})

viewMenu.Add("Zoom Out").SetAccelerator("CmdOrCtrl+-").OnClick(func(ctx *application.Context) {
    // Decrease zoom
})

viewMenu.Add("Reset Zoom").SetAccelerator("CmdOrCtrl+0").OnClick(func(ctx *application.Context) {
    // Reset zoom
})
```

**Obtenir la fenêtre active :**

```go
menuItem.OnClick(func(ctx *application.Context) {
    window := application.Get().Window.Current() // the window the menu was invoked from
    // Use window
})
```

## Considérations propres à chaque plateforme

### macOS

**Comportement de la barre de menus :**

- S’affiche **en haut de l’écran** (barre globale)
- Reste affichée lorsque toutes les fenêtres sont fermées
- Le premier menu est **toujours le menu de l’application**
- Utilisez `menu.AddRole(application.AppMenu)` pour les éléments standard

**Emplacements standard :**

- **À propos** : menu de l’application
- **Préférences** : menu de l’application (⌘,)
- **Quitter** : menu de l’application (⌘Q)
- **Aide** : menu Aide

**Exemple :**

```go
if runtime.GOOS == "darwin" {
    menu.AddRole(application.AppMenu)  // Adds About, Preferences, Quit
    
    // Don't add Quit to File menu on macOS
    // Don't add About to Help menu on macOS
}
```

### Windows

**Comportement de la barre de menus :**

- S’affiche dans la **barre de titre de la fenêtre**
- Chaque fenêtre possède son propre menu
- Aucun menu d’application

**Emplacements standard :**

- **Quitter** : menu Fichier (Alt+F4)
- **Paramètres** : menu Outils ou Édition
- **À propos** : menu Aide

**Exemple :**

```go
if runtime.GOOS == "windows" {
    menu.AddRole(application.FileMenu) // Exit is added automatically
    menu.AddRole(application.HelpMenu) // About is added automatically
}
```

### Linux

**Comportement de la barre de menus :**

- Généralement propre à chaque fenêtre (comme sous Windows)
- Certains environnements de bureau prennent en charge les menus globaux (Unity, GNOME avec une extension)
- L’apparence varie selon l’environnement de bureau

**Bonne pratique :** suivez les conventions de Windows et testez sur les environnements de bureau cibles.

## Exemple complet

Voici une structure de menus prête pour la production :

```go
package main

import (
    "runtime"
    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "My Application",
    })

    // Create and set menu
    createMenu(app)

    // Create main window with UseApplicationMenu for cross-platform menu support
    app.Window.NewWithOptions(application.WebviewWindowOptions{
        UseApplicationMenu: true,
    })

    app.Run()
}

func createMenu(app *application.App) {
    menu := app.NewMenu()

    // Platform-specific application menu (macOS only)
    if runtime.GOOS == "darwin" {
        menu.AddRole(application.AppMenu)
    }

    // File menu — AddRole returns the receiver menu, not the role submenu.
    // To add items into the File submenu, look it up via FindByRole + GetSubmenu.
    menu.AddRole(application.FileMenu)
    fileMenu := menu.FindByRole(application.FileMenu).GetSubmenu()
    fileMenu.Add("Import...").SetAccelerator("CmdOrCtrl+I").OnClick(handleImport)
    fileMenu.Add("Export...").SetAccelerator("CmdOrCtrl+E").OnClick(handleExport)

    // Edit menu
    menu.AddRole(application.EditMenu)

    // View menu
    viewMenu := menu.AddSubmenu("View")
    viewMenu.Add("Toggle Fullscreen").SetAccelerator("F11").OnClick(toggleFullscreen)
    viewMenu.AddSeparator()
    viewMenu.AddCheckbox("Show Sidebar", true).OnClick(toggleSidebar)
    viewMenu.AddCheckbox("Show Toolbar", true).OnClick(toggleToolbar)

    // Tools menu
    toolsMenu := menu.AddSubmenu("Tools")
    
    // Settings location varies by platform
    if runtime.GOOS == "darwin" {
        // On macOS, Preferences is in Application menu (added by AppMenu role)
    } else {
        toolsMenu.Add("Settings").SetAccelerator("CmdOrCtrl+,").OnClick(showSettings)
    }
    
    toolsMenu.AddSeparator()
    toolsMenu.AddCheckbox("Dark Mode", false).OnClick(toggleDarkMode)

    // Window menu
    menu.AddRole(application.WindowMenu)

    // Help menu
    helpMenu := menu.AddRole(application.HelpMenu)
    helpMenu.Add("Documentation").OnClick(openDocumentation)
    
    // About location varies by platform
    if runtime.GOOS == "darwin" {
        // On macOS, About is in Application menu (added by AppMenu role)
    } else {
        helpMenu.AddSeparator()
        helpMenu.Add("About").OnClick(showAbout)
    }

    // Set the application menu
    app.Menu.Set(menu)
}

func handleImport(ctx *application.Context) {
    // Implementation
}

func handleExport(ctx *application.Context) {
    // Implementation
}

func toggleFullscreen(ctx *application.Context) {
    window := application.Get().Window.Current()
    window.ToggleFullscreen()
}

func toggleSidebar(ctx *application.Context) {
    // Implementation
}

func toggleToolbar(ctx *application.Context) {
    // Implementation
}

func showSettings(ctx *application.Context) {
    // Implementation
}

func toggleDarkMode(ctx *application.Context) {
    isDark := ctx.ClickedMenuItem().Checked()
    // Apply theme
}

func openDocumentation(ctx *application.Context) {
    // Open browser
}

func showAbout(ctx *application.Context) {
    // Show about dialog
}
```

## Bonnes pratiques

### ✅ À faire

- **Utilisez les rôles de menu** pour les menus standard (Fichier, Édition, etc.)
- **Respectez les conventions de la plateforme** pour structurer les menus
- **Ajoutez des raccourcis clavier** aux actions courantes
- **Appelez menu.Update()** après avoir modifié l’état du menu
- **Testez sur toutes les plateformes** : le comportement varie
- **Limitez la profondeur des menus** à 2-3 niveaux maximum
- **Utilisez des libellés explicites** : « Enregistrer le projet » plutôt que « Enregistrer »

### ❌ À éviter

- **Ne codez pas en dur les raccourcis propres à chaque plateforme** : utilisez `CmdOrCtrl`
- **Ne placez pas Quitter dans le menu Fichier sous macOS** : cette commande se trouve dans le menu de l’application
- **Ne placez pas À propos dans le menu Aide sous macOS** : cette commande se trouve dans le menu de l’application
- **N’oubliez pas d’appeler menu.Update()** : les menus ne fonctionneront pas correctement
- **N’imbriquez pas trop profondément les menus** : les utilisateurs s’y perdent
- **N’utilisez pas de jargon** : choisissez des libellés faciles à comprendre

## Étapes suivantes

@cards{cols="2"}
📖 Référence des menus
Référence complète des types et propriétés des éléments de menu.

[En savoir plus →](/features/menus/reference/)

---
◆ Menus contextuels
Créez des menus contextuels accessibles par clic droit.

[En savoir plus →](/features/menus/context/)

---
★ Menus de la zone de notification
Ajoutez une intégration à la zone de notification ou à la barre des menus.

[En savoir plus →](/features/menus/systray/)

---
📖 Modèles de menus
Modèles de menus courants et bonnes pratiques.

[En savoir plus →](/guides/menus/)

@end

---

**Des questions ?** Posez-les sur [Discord](https://discord.gg/JDdSxwjhGf) ou consultez l’[exemple de menu](https://github.com/wailsapp/wails/tree/master/v3/examples/menu).
