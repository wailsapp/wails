---
title: "API des menus"
description: "Référence complète de l’API des menus"
slug: "reference/menu"
sourcePath: "reference/menu.md"
---

## Présentation

L’API des menus fournit des méthodes permettant de créer et de gérer les menus d’application, les menus contextuels et les menus de la zone de notification.

**Types de menus :**

- **Menus d’application** — Barre de menus supérieure (Fichier, Édition, etc.)
- **Menus contextuels** — Menus accessibles par clic droit
- **Menus de la zone de notification** — Menus situés dans la zone de notification système

## Création de menus

### NewMenu()

Crée un menu.

```go
func (a *App) NewMenu() *Menu
```

**Exemple :**

```go
menu := app.NewMenu()
```

## Méthodes des menus

### Add()

Ajoute un élément au menu.

```go
func (m *Menu) Add(label string) *MenuItem
```

**Paramètres :**

- `label` — Texte affiché pour l’élément de menu

**Valeur renvoyée :** l’élément de menu créé

**Exemple :**

```go
item := menu.Add("Open File")
item.OnClick(func(ctx *application.Context) {
    // Handle click
})
```

### AddSubmenu()

Ajoute un sous-menu au menu.

```go
func (m *Menu) AddSubmenu(label string) *Menu
```

**Paramètres :**

- `label` — Libellé du sous-menu

**Valeur renvoyée :** le sous-menu créé

**Exemple :**

```go
fileMenu := menu.AddSubmenu("File")
fileMenu.Add("New")
fileMenu.Add("Open")
fileMenu.Add("Save")
```

### AddSeparator()

Ajoute une ligne de séparation visuelle entre les éléments de menu.

```go
func (m *Menu) AddSeparator()
```

**Exemple :**

```go
menu.Add("Copy")
menu.Add("Paste")
menu.AddSeparator()
menu.Add("Select All")
```

**Bonne pratique :** utilisez des séparateurs pour regrouper les éléments de menu associés.

### AddCheckbox()

Ajoute un élément de menu à cocher.

```go
func (m *Menu) AddCheckbox(label string, checked bool) *MenuItem
```

**Paramètres :**

- `label` — Libellé de la case à cocher
- `checked` — État coché initial

**Exemple :**

```go
darkMode := menu.AddCheckbox("Dark Mode", false)
darkMode.OnClick(func(ctx *application.Context) {
    isChecked := darkMode.Checked()
    // Toggle dark mode
})
```

### AddRadio()

Ajoute un élément de menu de type bouton radio (dans un groupe à choix mutuellement exclusifs).

```go
func (m *Menu) AddRadio(label string, checked bool) *MenuItem
```

**Paramètres :**

- `label` — Libellé du bouton radio
- `checked` — État coché initial

**Exemple :**

```go
// Create radio group for view modes
viewMenu := menu.AddSubmenu("View")
listView := viewMenu.AddRadio("List View", true)
gridView := viewMenu.AddRadio("Grid View", false)
treeView := viewMenu.AddRadio("Tree View", false)

listView.OnClick(func(ctx *application.Context) {
    setViewMode("list")
})
gridView.OnClick(func(ctx *application.Context) {
    setViewMode("grid")
})
```

### Update()

Met à jour le menu afin de refléter les modifications apportées aux éléments de menu.

```go
func (m *Menu) Update()
```

**Exemple :**

```go
item.SetEnabled(false)
menu.Update()  // Must call to apply changes
```

**Important :** appelez toujours `Update()` après avoir modifié les propriétés d’un élément de menu.

## Méthodes des éléments de menu

### OnClick()

Enregistre un gestionnaire de clic pour l’élément de menu.

```go
func (mi *MenuItem) OnClick(callback func(ctx *application.Context)) *MenuItem
```

**Paramètres :**

- `callback` — Fonction appelée lorsque l’utilisateur clique sur l’élément

**Valeur renvoyée :** l’élément de menu (pour le chaînage)

**Exemple :**

```go
item.OnClick(func(ctx *application.Context) {
    fmt.Println("Menu item clicked")
    app.Logger.Info("User clicked menu item")
})
```

### SetLabel()

Modifie le libellé de l’élément de menu.

```go
func (mi *MenuItem) SetLabel(label string) *MenuItem
```

**Exemple :**

```go
item.SetLabel("Save As...")
menu.Update()
```

### SetEnabled()

Active ou désactive l’élément de menu.

```go
func (mi *MenuItem) SetEnabled(enabled bool) *MenuItem
```

**Exemple :**

```go
// Disable save when no document is open
saveItem.SetEnabled(hasOpenDocument)
menu.Update()
```

**Modèle courant :**

```go
// Update menu state based on application state
func updateMenuState() {
    saveItem.SetEnabled(hasUnsavedChanges)
    undoItem.SetEnabled(canUndo)
    redoItem.SetEnabled(canRedo)
    menu.Update()
}
```

### SetChecked()

Définit l’état coché des éléments de menu de type case à cocher ou bouton radio.

```go
func (mi *MenuItem) SetChecked(checked bool) *MenuItem
```

**Exemple :**

```go
darkModeItem.SetChecked(isDarkModeEnabled)
menu.Update()
```

### Checked()

Renvoie l’état coché actuel.

```go
func (mi *MenuItem) Checked() bool
```

**Exemple :**

```go
if darkModeItem.Checked() {
    // Dark mode is enabled
}
```

### SetAccelerator()

Définit un raccourci clavier pour l’élément de menu.

```go
func (mi *MenuItem) SetAccelerator(accelerator string) *MenuItem
```

**Paramètres :**

- `accelerator` — Raccourci clavier (par exemple, "Ctrl+S", "Cmd+Q")

**Format du raccourci clavier :**

- **Touches de modification :** `Ctrl`, `Cmd`, `Alt`, `Shift`
- **Touches :** `A-Z`, `0-9`, `F1-F12`, `Enter`, `Backspace`, etc.
- **Plateforme :** utilisez `Cmd` sous macOS et `Ctrl` sous Windows/Linux.

**Exemple :**

```go
saveItem.SetAccelerator("Ctrl+S")
quitItem.SetAccelerator("Ctrl+Q")
newItem.SetAccelerator("Ctrl+N")
```

**Exemple adapté à la plateforme :**

```go
import "runtime"

var quitShortcut string
if runtime.GOOS == "darwin" {
    quitShortcut = "Cmd+Q"
} else {
    quitShortcut = "Ctrl+Q"
}
quitItem.SetAccelerator(quitShortcut)
```

### SetTooltip()

Définit une infobulle qui s’affiche au survol de l’élément de menu.

```go
func (mi *MenuItem) SetTooltip(tooltip string) *MenuItem
```

**Exemple :**

```go
item.SetTooltip("Opens a file from disk")
```

### SetHidden()

Affiche ou masque l’élément de menu.

```go
func (mi *MenuItem) SetHidden(hidden bool) *MenuItem
```

**Exemple :**

```go
// Hide debug menu in production
debugItem.SetHidden(!isDevelopment)
menu.Update()
```

## Menu de l’application

### app.Menu.Set()

Définit la barre de menus principale de l’application.

```go
func (mm *MenuManager) Set(menu *Menu)
```

**Exemple :**

```go
menu := app.NewMenu()

// File menu
fileMenu := menu.AddSubmenu("File")
fileMenu.Add("New").SetAccelerator("Ctrl+N").OnClick(newFile)
fileMenu.Add("Open").SetAccelerator("Ctrl+O").OnClick(openFile)
fileMenu.Add("Save").SetAccelerator("Ctrl+S").OnClick(saveFile)
fileMenu.AddSeparator()
fileMenu.Add("Exit").SetAccelerator("Ctrl+Q").OnClick(func(ctx *application.Context) {
    app.Quit()
})

// Edit menu
editMenu := menu.AddSubmenu("Edit")
editMenu.Add("Undo").SetAccelerator("Ctrl+Z").OnClick(undo)
editMenu.Add("Redo").SetAccelerator("Ctrl+Y").OnClick(redo)
editMenu.AddSeparator()
editMenu.Add("Cut").SetAccelerator("Ctrl+X").OnClick(cut)
editMenu.Add("Copy").SetAccelerator("Ctrl+C").OnClick(copy)
editMenu.Add("Paste").SetAccelerator("Ctrl+V").OnClick(paste)

app.Menu.Set(menu)
```

**Remarques propres aux plateformes :**

- **macOS :** le menu apparaît dans la barre de menus en haut de l’écran
- **Windows/Linux :** le menu apparaît dans la barre de titre de la fenêtre
- **macOS :** ajoute automatiquement un menu d’application portant le nom de l’application

## Menus contextuels

### app.ContextMenu.New() / app.ContextMenu.Add()

Créez un `*ContextMenu` via le gestionnaire, puis enregistrez-le sous un nom. `ContextMenuManager.Add` accepte `*ContextMenu` — **pas** `*Menu` — et il n’existe **aucune** méthode `app.RegisterContextMenu`.

```go
func (cm *ContextMenuManager) New() *ContextMenu
func (cm *ContextMenuManager) Add(name string, menu *ContextMenu)
func (cm *ContextMenuManager) Get(name string) (*ContextMenu, bool)
func (cm *ContextMenuManager) Remove(name string)
```

Vous pouvez également utiliser la fonction `application.NewContextMenu(name string) *ContextMenu` au niveau du paquet pour créer et enregistrer un menu contextuel en une seule étape.

**Go :**

```go
// Build the context menu via the manager
contextMenu := app.ContextMenu.New()
contextMenu.Add("Cut").OnClick(cut)
contextMenu.Add("Copy").OnClick(copy)
contextMenu.Add("Paste").OnClick(paste)
contextMenu.AddSeparator()
contextMenu.Add("Select All").OnClick(selectAll)

// Register under a name; HTML opts into it via the CSS custom property below.
app.ContextMenu.Add("editor", contextMenu)
```

**HTML/CSS :**

Le runtime déclenche un menu contextuel enregistré lorsque la cible du clic droit (ou l’un de ses ancêtres) possède la propriété personnalisée CSS `--custom-contextmenu` définie sur ce nom. La propriété facultative `--custom-contextmenu-data` est transmise à la fonction de rappel Go via `ctx.ContextMenuData()`. Pour supprimer le menu contextuel par défaut du navigateur, définissez `--default-contextmenu: hide` (ou `auto`/`show`).

```html
<!-- Trigger context menu on right-click -->
<div style="--custom-contextmenu: editor; --default-contextmenu: hide">
    Right-click here for context menu
</div>
```

Il n’existe aucun attribut `data-wails-context-menu="..."` : il n’a jamais été raccordé au runtime.

**Menus contextuels dynamiques :**

```go
// Update context menu based on selection
func updateContextMenu() {
    contextMenu := app.ContextMenu.New()

    if hasSelection {
        contextMenu.Add("Cut").OnClick(cut)
        contextMenu.Add("Copy").OnClick(copy)
    }

    contextMenu.Add("Paste").SetEnabled(hasClipboardContent).OnClick(paste)

    app.ContextMenu.Add("editor", contextMenu)
}
```

## Menu de la zone de notification

### app.SystemTray.New()

Crée une nouvelle icône dans la zone de notification.

```go
func (sm *SystemTrayManager) New() *SystemTray
```

**Exemple :**

```go
tray := app.SystemTray.New()
```

### SetIcon()

Définit l’icône de la zone de notification.

```go
func (st *SystemTray) SetIcon(icon []byte) *SystemTray
```

**Exemple :**

```go
iconData, _ := os.ReadFile("icon.png")
tray.SetIcon(iconData)
```

### SetMenu()

Définit le menu de la zone de notification.

```go
func (st *SystemTray) SetMenu(menu *Menu) *SystemTray
```

**Exemple :**

```go
trayMenu := app.NewMenu()
trayMenu.Add("Show Window").OnClick(func(ctx *application.Context) {
    window.Show()
    window.Focus()
})
trayMenu.Add("Settings").OnClick(openSettings)
trayMenu.AddSeparator()
trayMenu.Add("Quit").OnClick(func(ctx *application.Context) {
    app.Quit()
})

tray.SetMenu(trayMenu)
```

### SetTooltip()

Définit l’infobulle affichée au survol de l’icône de la zone de notification. Ne renvoie rien.

```go
func (st *SystemTray) SetTooltip(tooltip string)
```

**Exemple :**

```go
tray.SetTooltip("My Application - Running")
```

### OnClick()

Gère le clic gauche sur l’icône de la zone de notification.

```go
func (st *SystemTray) OnClick(callback func()) *SystemTray
```

**Exemple :**

```go
tray.OnClick(func() {
    if window.IsVisible() {
        window.Hide()
    } else {
        window.Show()
        window.Focus()
    }
})
```

## Exemples complets

### Menu d’application standard

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
)

func createMenu(app *application.App) *application.Menu {
    menu := app.NewMenu()

    // File menu
    fileMenu := menu.AddSubmenu("File")
    fileMenu.Add("New").
        SetAccelerator("Ctrl+N").
        OnClick(func(ctx *application.Context) {
            // Create new document
        })
    fileMenu.Add("Open").
        SetAccelerator("Ctrl+O").
        OnClick(func(ctx *application.Context) {
            // Open file dialog
        })
    fileMenu.Add("Save").
        SetAccelerator("Ctrl+S").
        OnClick(func(ctx *application.Context) {
            // Save document
        })
    fileMenu.AddSeparator()
    fileMenu.Add("Exit").OnClick(func(ctx *application.Context) {
        app.Quit()
    })

    // Edit menu
    editMenu := menu.AddSubmenu("Edit")
    editMenu.Add("Undo").SetAccelerator("Ctrl+Z")
    editMenu.Add("Redo").SetAccelerator("Ctrl+Y")
    editMenu.AddSeparator()
    editMenu.Add("Cut").SetAccelerator("Ctrl+X")
    editMenu.Add("Copy").SetAccelerator("Ctrl+C")
    editMenu.Add("Paste").SetAccelerator("Ctrl+V")

    // View menu
    viewMenu := menu.AddSubmenu("View")
    darkMode := viewMenu.AddCheckbox("Dark Mode", false)
    darkMode.OnClick(func(ctx *application.Context) {
        // Toggle dark mode
        isChecked := darkMode.Checked()
        app.Logger.Info("Dark mode", "enabled", isChecked)
    })
    viewMenu.AddSeparator()
    viewMenu.AddRadio("List View", true)
    viewMenu.AddRadio("Grid View", false)
    viewMenu.AddRadio("Detail View", false)

    // Help menu
    helpMenu := menu.AddSubmenu("Help")
    helpMenu.Add("Documentation").OnClick(func(ctx *application.Context) {
        // Open docs
    })
    helpMenu.Add("About").OnClick(func(ctx *application.Context) {
        // Show about dialog
    })

    return menu
}

func main() {
    app := application.New(application.Options{
        Name: "Menu Demo",
    })

    menu := createMenu(app)
    app.Menu.Set(menu)

    window := app.Window.New()
    window.Show()

    app.Run()
}
```

### Application avec icône dans la zone de notification

```go
func setupSystemTray(app *application.App, window application.Window) {
    // Create system tray
    tray := app.SystemTray.New()

    // Set icon
    iconData, _ := os.ReadFile("icon.png")
    tray.SetIcon(iconData)
    tray.SetTooltip("My App - Running")

    // Handle left-click on tray icon
    tray.OnClick(func() {
        if window.IsVisible() {
            window.Hide()
        } else {
            window.Show()
            window.Focus()
        }
    })

    // Create tray menu
    trayMenu := app.NewMenu()

    showItem := trayMenu.Add("Show Window")
    showItem.OnClick(func(ctx *application.Context) {
        window.Show()
        window.Focus()
    })

    trayMenu.AddSeparator()

    trayMenu.Add("Settings").OnClick(func(ctx *application.Context) {
        // Open settings window
    })

    trayMenu.Add("About").OnClick(func(ctx *application.Context) {
        // Show about dialog
    })

    trayMenu.AddSeparator()

    trayMenu.Add("Quit").OnClick(func(ctx *application.Context) {
        app.Quit()
    })

    tray.SetMenu(trayMenu)
}
```

### Mises à jour dynamiques des menus

```go
type Editor struct {
    app         *application.App
    menu        *application.Menu
    undoItem    *application.MenuItem
    redoItem    *application.MenuItem
    saveItem    *application.MenuItem
    undoStack   []string
    redoStack   []string
    hasChanges  bool
}

func (e *Editor) createMenu() {
    e.menu = e.app.NewMenu()

    fileMenu := e.menu.AddSubmenu("File")
    e.saveItem = fileMenu.Add("Save").SetAccelerator("Ctrl+S")
    e.saveItem.OnClick(func(ctx *application.Context) {
        e.save()
    })

    editMenu := e.menu.AddSubmenu("Edit")
    e.undoItem = editMenu.Add("Undo").SetAccelerator("Ctrl+Z")
    e.undoItem.OnClick(func(ctx *application.Context) {
        e.undo()
    })

    e.redoItem = editMenu.Add("Redo").SetAccelerator("Ctrl+Y")
    e.redoItem.OnClick(func(ctx *application.Context) {
        e.redo()
    })

    e.updateMenuState()
    e.app.Menu.Set(e.menu)
}

func (e *Editor) updateMenuState() {
    // Update menu items based on current state
    e.saveItem.SetEnabled(e.hasChanges)
    e.undoItem.SetEnabled(len(e.undoStack) > 0)
    e.redoItem.SetEnabled(len(e.redoStack) > 0)
    e.menu.Update()
}

func (e *Editor) onChange() {
    e.hasChanges = true
    e.updateMenuState()
}

func (e *Editor) save() {
    // Save logic
    e.hasChanges = false
    e.updateMenuState()
}
```

## Bonnes pratiques

### ✅ À faire

- **Utilisez les raccourcis clavier standard** — respectez les conventions de la plateforme (Ctrl+C pour copier, etc.).
- **Appelez Update() après les modifications** — sinon, le menu ne les reflétera pas.
- **Regroupez les éléments associés** — utilisez des séparateurs pour organiser les éléments de menu.
- **Désactivez les actions indisponibles** — ne les masquez pas ; désactivez-les avec SetEnabled(false).
- **Utilisez des libellés clairs** — soyez concis et explicite.
- **Respectez les conventions de la plateforme** — les structures de menus diffèrent entre macOS et Windows/Linux.

### ❌ À éviter

- **N’oubliez pas Update()** — c’est l’erreur la plus courante.
- **N’imbriquez pas trop profondément les menus** — limitez-les à 2-3 niveaux.
- **N’utilisez pas de libellés ambigus** — « Traiter » contre « Traiter le document ».
- **Ne compliquez pas inutilement les menus** — gardez-les simples et ciblés.
- **Ne mélangez pas les métaphores** — adoptez une dénomination et une organisation cohérentes.

## Remarques propres aux plateformes

### macOS

- Un menu d’application portant le nom de l’application est ajouté automatiquement.
- Utilisez `Cmd` plutôt que `Ctrl` pour les raccourcis clavier.
- Par défaut, « À propos », « Préférences » et « Quitter » figurent dans le menu de l’application.

### Windows/Linux

- Aucun menu d’application automatique
- Utilisez `Ctrl` pour les raccourcis clavier
- « Quitter » se trouve généralement dans le menu Fichier
