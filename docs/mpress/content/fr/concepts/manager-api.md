---
title: "API des gestionnaires"
description: "Structure d’API organisée autour d’interfaces de gestion spécialisées"
slug: "concepts/manager-api"
sourcePath: "concepts/manager-api.md"
---

L’API des gestionnaires de Wails v3 permet d’accéder aux fonctionnalités de l’application de façon structurée et intuitive, au moyen de structures de gestion spécialisées regroupées dans les champs publics de `*application.App`. Wails 3 rompt nettement avec la v2 : aucune couche d’adaptation par appel ne maintient la compatibilité avec l’ancienne API de type `app.NewWebviewWindow(...)`. Les gestionnaires ci-dessous constituent donc le moyen de piloter l’application.

## Vue d’ensemble

L’API des gestionnaires organise les fonctionnalités de l’application en douze domaines spécialisés (un journaliseur et onze gestionnaires) :

- **`app.Window`** - Création et gestion des fenêtres, et fonctions de rappel
- **`app.ContextMenu`** - Enregistrement et gestion des menus contextuels\
- **`app.KeyBinding`** - Gestion des raccourcis clavier globaux
- **`app.Browser`** - Intégration au navigateur (ouverture d’URL et de fichiers)
- **`app.Env`** - Informations sur l’environnement et état du système
- **`app.Dialog`** - Opérations relatives aux boîtes de dialogue de fichiers et de messages
- **`app.Event`** - Gestion des événements personnalisés et des événements de l’application
- **`app.Menu`** - Gestion du menu de l’application
- **`app.Screen`** - Gestion des écrans et transformations de coordonnées
- **`app.Clipboard`** - Opérations de texte dans le presse-papiers
- **`app.SystemTray`** - Création et gestion de l’icône de la zone de notification
- **`app.Autostart`** - Enregistrement de l’application pour son lancement à la connexion de l’utilisateur

## Avantages

- **Meilleure facilité de découverte** - La saisie semi-automatique de l’IDE présente une surface d’API structurée
- **Meilleure organisation du code** - Les méthodes connexes sont regroupées
- **Maintenabilité accrue** - Les responsabilités sont séparées entre les gestionnaires
- **Extensibilité future** - L’ajout de nouvelles fonctionnalités dans des domaines précis est facilité

## Utilisation

L’API des gestionnaires offre un accès structuré à toutes les fonctionnalités de l’application :

```go
// Events and custom event handling
app.Event.Emit("custom", data)
app.Event.On("custom", func(e *CustomEvent) { ... })

// Window management
window, _ := app.Window.GetByName("main")
app.Window.OnCreate(func(window Window) { ... })

// Browser integration
app.Browser.OpenURL("https://wails.io")

// Menu management
menu := app.Menu.New()
app.Menu.Set(menu)

// System tray
systray := app.SystemTray.New()
```

## Référence des gestionnaires

### Gestionnaire de fenêtres

Gère la création et la récupération des fenêtres, ainsi que les fonctions de rappel de leur cycle de vie.

```go
// Create windows
window := app.Window.New()
window := app.Window.NewWithOptions(options)
current := app.Window.Current()

// Find windows
window, exists := app.Window.GetByName("main")
windows := app.Window.GetAll()

// Window callbacks
app.Window.OnCreate(func(window Window) {
    // Handle window creation
})
```

### Gestionnaire d’événements

Gère les événements personnalisés et l’écoute des événements de l’application.

```go
// Custom events
app.Event.Emit("userAction", data)
cancelFunc := app.Event.On("userAction", func(e *CustomEvent) {
    // Handle event
})
app.Event.Off("userAction")
app.Event.Reset() // Remove all listeners

// Application events
app.Event.OnApplicationEvent(events.Common.ThemeChanged, func(e *ApplicationEvent) {
    // Handle system theme change
})
```

### Gestionnaire de navigateur

Fournit l’intégration au navigateur nécessaire pour ouvrir des URL et des fichiers.

```go
// Open URLs and files in default browser
err := app.Browser.OpenURL("https://wails.io")
err := app.Browser.OpenFile("/path/to/document.pdf")
```

### Gestionnaire d’environnement

Donne accès aux informations sur l’environnement système.

```go
// Get environment info
env := app.Env.Info()
fmt.Printf("OS: %s, Arch: %s\n", env.OS, env.Arch)

// Check system theme
if app.Env.IsDarkMode() {
    // Dark mode is active
}

// Open file manager
err := app.Env.OpenFileManager("/path/to/folder", false)
```

### Gestionnaire de boîtes de dialogue

Fournit un accès structuré aux boîtes de dialogue de fichiers et de messages.

```go
// File dialogs
result, err := app.Dialog.OpenFile().
    AddFilter("Text Files", "*.txt").
    PromptForSingleSelection()

result, err = app.Dialog.SaveFile().
    SetFilename("document.txt").
    PromptForSingleSelection()

// Message dialogs
app.Dialog.Info().
    SetTitle("Information").
    SetMessage("Operation completed successfully").
    Show()

app.Dialog.Error().
    SetTitle("Error").
    SetMessage("An error occurred").
    Show()
```

### Gestionnaire de menus

Création et gestion du menu de l’application.

```go
// Create and set application menu
menu := app.Menu.New()
fileMenu := menu.AddSubmenu("File")
fileMenu.Add("New").OnClick(func(ctx *Context) {
    // Handle menu click
})

app.Menu.Set(menu)

// Show about dialog
app.Menu.ShowAbout()
```

### Gestionnaire de raccourcis clavier

Gestion dynamique des raccourcis clavier globaux.

```go
// Add key bindings
app.KeyBinding.Add("ctrl+n", func(window application.Window) {
    // Handle Ctrl+N
})

app.KeyBinding.Add("ctrl+q", func(window application.Window) {
    app.Quit()
})

// Remove key bindings
app.KeyBinding.Remove("ctrl+n")

// Get all bindings
bindings := app.KeyBinding.GetAll()
```

### Gestionnaire de menus contextuels

Gestion avancée des menus contextuels (destinée aux auteurs de bibliothèques).

```go
// Create and register context menu
menu := app.ContextMenu.New()
app.ContextMenu.Add("myMenu", menu)

// Retrieve context menu
menu, exists := app.ContextMenu.Get("myMenu")

// Remove context menu
app.ContextMenu.Remove("myMenu")
```

### Gestionnaire d’écrans

Gestion des écrans et transformations de coordonnées pour les configurations à plusieurs moniteurs.

```go
// Get screen information
screens := app.Screen.GetAll()
primary := app.Screen.GetPrimary()

// Coordinate transformations
physicalPoint := app.Screen.DipToPhysicalPoint(logicalPoint)
logicalPoint := app.Screen.PhysicalToDipPoint(physicalPoint)

// Screen detection
screen := app.Screen.ScreenNearestDipPoint(point)
screen = app.Screen.ScreenNearestDipRect(rect)
```

### Gestionnaire du presse-papiers

Opérations de lecture et d’écriture de texte dans le presse-papiers.

```go
// Set text to clipboard
success := app.Clipboard.SetText("Hello World")
if !success {
    // Handle error
}

// Get text from clipboard
text, ok := app.Clipboard.Text()
if !ok {
    // Handle error
} else {
    // Use the text
}
```

### Gestionnaire SystemTray

Création et gestion de l’icône de la zone de notification.

```go
// Create system tray
systray := app.SystemTray.New()
systray.SetLabel("My App")
systray.SetIcon(iconBytes)

// Add menu to system tray
menu := app.Menu.New()
menu.Add("Open").OnClick(func(ctx *Context) {
    // Handle click
})
systray.SetMenu(menu)

// Destroy system tray when done
systray.Destroy()
```

### Gestionnaire de démarrage automatique

Enregistre l’application afin qu’elle se lance à la connexion de l’utilisateur. Sélectionne le mécanisme natif approprié à chaque plateforme : SMAppService ou un fichier plist LaunchAgent sous macOS, la clé de Registre `HKCU\…\Run` sous Windows et une entrée XDG `.desktop` sous Linux.

```go
// Register to launch at login
err := app.Autostart.Enable()

// With extra launch-time arguments and a custom identifier
err = app.Autostart.EnableWithOptions(application.AutostartOptions{
    Identifier: "com.example.myapp",
    Arguments:  []string{"--hidden"},
})

// Check / remove
enabled, err := app.Autostart.IsEnabled()
status, err := app.Autostart.Status()  // includes Path + Strategy
err = app.Autostart.Disable()
```

Consultez la [page consacrée à la fonctionnalité de démarrage automatique](/features/autostart/basics/) pour connaître le comportement propre à chaque plateforme, les règles relatives aux identifiants et la garantie de détection des entrées obsolètes.
