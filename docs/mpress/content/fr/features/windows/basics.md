---
title: "Principes de base des fenêtres"
description: "Création et gestion des fenêtres d’application dans Wails"
slug: "features/windows/basics"
sourcePath: "features/windows/basics.md"
---

## Gestion des fenêtres

Wails fournit une **API unifiée de gestion des fenêtres** qui fonctionne sur toutes les plateformes. Créez des fenêtres, contrôlez leur comportement et gérez plusieurs fenêtres en maîtrisant entièrement leur création, leur apparence, leur comportement et leur cycle de vie.

## Démarrage rapide

```go
package main

import "github.com/wailsapp/wails/v3/pkg/application"

func main() {
    app := application.New(application.Options{
        Name: "My App",
    })

    // Create a window
    window := app.Window.New()
    
    // Configure it
    window.SetTitle("Hello Wails")
    window.SetSize(800, 600)
    window.Center()
    
    // Show it
    window.Show()

    app.Run()
}
```

**Voilà, c’est tout !** Vous disposez maintenant d’une fenêtre multiplateforme.

## Création de fenêtres

### Fenêtre de base

La méthode la plus simple pour créer une fenêtre :

```go
window := app.Window.New()
```

**Ce que vous obtenez :**

- Taille par défaut (800x600)
- Titre par défaut (nom de l’application)
- WebView prête pour votre interface frontend
- Apparence native de la plateforme

### Fenêtre avec options

Créez une fenêtre avec une configuration personnalisée :

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "My Application",
    Width:  1200,
    Height: 800,
    X:      100,   // Position from left
    Y:      100,   // Position from top
    AlwaysOnTop: false,
    Frameless: false,
    Hidden: false,
    MinWidth: 400,
    MinHeight: 300,
    MaxWidth: 1920,
    MaxHeight: 1080,
})
```

**Options courantes :**

| Option | Type | Description |
| --- | --- | --- |
| `Title` | `string` | Titre de la fenêtre |
| `Width` | `int` | Largeur de la fenêtre en pixels |
| `Height` | `int` | Hauteur de la fenêtre en pixels |
| `X` | `int` | Position X (depuis la gauche) |
| `Y` | `int` | Position Y (depuis le haut) |
| `AlwaysOnTop` | `bool` | Maintenir la fenêtre au-dessus des autres |
| `Frameless` | `bool` | Supprimer la barre de titre et les bordures |
| `Hidden` | `bool` | Démarrer avec la fenêtre masquée |
| `MinWidth` | `int` | Largeur minimale |
| `MinHeight` | `int` | Hauteur minimale |
| `MaxWidth` | `int` | Largeur maximale |
| `MaxHeight` | `int` | Hauteur maximale |

**Consultez la page [Options des fenêtres](/features/windows/options/) pour obtenir la liste complète.**

### Fenêtres nommées

Attribuez un nom aux fenêtres pour les retrouver facilement :

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:  "main-window",
    Title: "Main Application",
})

// Later, find it by name
if mainWindow, ok := app.Window.GetByName("main-window"); ok {
    mainWindow.Show()
}
```

**Cas d’utilisation :**

- Plusieurs fenêtres (principale, paramètres, à propos)
- Recherche de fenêtres depuis différentes parties de votre code
- Communication entre les fenêtres

## Contrôle des fenêtres

### Afficher et masquer

```go
// Show window
window.Show()

// Hide window
window.Hide()

// Check if visible
if window.IsVisible() {
    fmt.Println("Window is visible")
}
```

**Cas d’utilisation :**

- Écrans de démarrage (afficher, puis masquer)
- Fenêtres de paramètres (masquer lorsqu’elles ne sont pas nécessaires)
- Fenêtres contextuelles (afficher à la demande)

### Position et taille

```go
// Set size
window.SetSize(1024, 768)

// Set position
window.SetPosition(100, 100)

// Centre on screen
window.Center()

// Get current size
width, height := window.Size()

// Get current position
x, y := window.Position()
```

**Système de coordonnées :**

- (0, 0) correspond au coin supérieur gauche de l’écran principal
- Les valeurs X positives vont vers la droite
- Les valeurs Y positives vont vers le bas

### État de la fenêtre

```go
// Minimise
window.Minimise()

// Maximise
window.Maximise()

// Fullscreen
window.Fullscreen()

// Restore to normal
window.Restore()

// Check state
if window.IsMinimised() {
    fmt.Println("Window is minimised")
}

if window.IsMaximised() {
    fmt.Println("Window is maximised")
}

if window.IsFullscreen() {
    fmt.Println("Window is fullscreen")
}
```

**Transitions d’état :**

```
Normal ←→ Minimised
Normal ←→ Maximised
Normal ←→ Fullscreen
```

### Titre et apparence

```go
// Set title
window.SetTitle("My Application - Document.txt")

// Set background colour — RGBA value (helper for RGB)
window.SetBackgroundColour(application.NewRGBA(0, 0, 0, 255))

// Set always on top
window.SetAlwaysOnTop(true)

// Set resizable
window.SetResizable(false)
```

### Fermeture des fenêtres

```go
// Close window — dispatches WindowClosing; a RegisterHook can call e.Cancel().
window.Close()
```

Il n’existe aucune méthode `window.Destroy()` dans la v3 : utilisez `Close()` et écoutez l’événement avec `OnWindowEvent` (sans possibilité d’annulation) ou interceptez-le avec `RegisterHook` (qui permet d’appeler `e.Cancel()` pour maintenir la fenêtre ouverte).

## Recherche de fenêtres

### Par nom

```go
if window, ok := app.Window.GetByName("settings"); ok {
    window.Show()
}
```

### Par ID

Chaque fenêtre possède un ID unique :

```go
id := window.ID()
fmt.Printf("Window ID: %d\n", id)

// Find by ID
if found, ok := app.Window.GetByID(id); ok {
    found.Focus()
}
```

### Fenêtre actuelle

Obtenez la fenêtre qui a actuellement le focus :

```go
current := app.Window.Current()
if current != nil {
    current.SetTitle("Active Window")
}
```

### Toutes les fenêtres

Obtenez toutes les fenêtres :

```go
windows := app.Window.GetAll()
fmt.Printf("Total windows: %d\n", len(windows))

for _, w := range windows {
    fmt.Printf("Window: %s (ID: %d)\n", w.Name(), w.ID())
}
```

## Cycle de vie des fenêtres

### Création

```go
app.Window.OnCreate(func(window application.Window) {
    fmt.Printf("Window created: %s\n", window.Name())

    // Configure new windows
    window.SetMinSize(400, 300)
})
```

### Fermeture

Pour empêcher la fermeture d’une fenêtre, utilisez `RegisterHook` avec l’événement `WindowClosing` :

```go
window.RegisterHook(events.Common.WindowClosing, func(event *application.WindowEvent) {
    if hasUnsavedChanges() {
        // Ask user for confirmation
        result := showConfirmDialog("Unsaved changes. Close anyway?")
        if result != "yes" {
            // Cancel the close event
            event.Cancel()
        }
    }
})
```

**Important :** `RegisterHook` intercepte l’événement de fermeture avant qu’il ne se produise. Appelez `event.Cancel()` pour empêcher la fermeture de la fenêtre. Cela fonctionne pour les fermetures déclenchées par l’utilisateur (clic sur le bouton X).

### Destruction

Pour effectuer le nettoyage lorsqu’une fenêtre se ferme, utilisez `OnWindowEvent` avec l’événement `WindowClosing` :

```go
window.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
    fmt.Println("Window is closing")
    // Cleanup resources
})
```

## Fenêtres multiples

### Création de plusieurs fenêtres

```go
// Main window
mainWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:   "main",
    Title:  "Main Application",
    Width:  1200,
    Height: 800,
})

// Settings window
settingsWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:   "settings",
    Title:  "Settings",
    Width:  600,
    Height: 400,
    Hidden: true,  // Start hidden
})

// Show settings when needed
settingsWindow.Show()
```

### Communication entre les fenêtres

Les fenêtres peuvent communiquer au moyen d’événements :

```go
// In main window
app.Event.Emit("data-updated", map[string]interface{}{
    "value": 42,
})

// In settings window
app.Event.On("data-updated", func(event *application.CustomEvent) {
    data := event.Data.(map[string]interface{})
    value := data["value"].(int)
    fmt.Printf("Received: %d\n", value)
})
```

**Pour en savoir plus, consultez [Événements](/features/events/system/).**

### Fenêtres parentes et enfants

`WebviewWindowOptions` ne comporte aucun champ `Parent`. Créez la fenêtre enfant comme une fenêtre normale, puis rattachez-la à une fenêtre parente sous la forme d’une feuille modale :

```go
// Create child window
childWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Child Window",
})

// Attach to the parent — presents as a sheet on macOS.
mainWindow.AttachModal(childWindow)
```

**Comportement :**

- La fenêtre enfant reste au-dessus de la fenêtre parente.
- La fenêtre enfant est modale : elle bloque les interactions avec la fenêtre parente.

**Prise en charge par plateforme :**

- **macOS :** prise en charge complète (affichage sous forme de feuille).
- **Windows :** non pris en charge.
- **Linux :** non pris en charge.

## Fonctionnalités propres aux plateformes

@tabs{sync-key="platform"}
[Windows]
**Fonctionnalités propres à Windows :**

```go
// Flash taskbar button
window.Flash(true)  // Start flashing
window.Flash(false) // Stop flashing

// Trigger Windows 11 Snap Assist (Win+Z)
window.SnapAssist()
```

Il n’existe aucun `SetIcon` propre à chaque fenêtre : l’icône de l’application est définie au niveau de l’application avec `app.SetIcon([]byte)` (ou, pour une icône de fenêtre propre à Linux, avec le champ `application.LinuxWindow.Icon` lors de la création de la fenêtre).

**Snap Assist :** Affiche les options de disposition d’ancrage de Windows 11 au moyen du raccourci système. Pour créer un bouton HTML personnalisé d’agrandissement qui affiche les dispositions d’ancrage natives au survol, utilisez plutôt [Régions non clientes natives sous Windows](/features/windows/frameless/#native-non-client-regions-on-windows).

**Clignotement dans la barre des tâches :** Utile pour signaler des notifications lorsque la fenêtre est réduite.

[macOS]
**Fonctionnalités propres à macOS :**

```go
// Transparent title bar
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Mac: application.MacWindow{
        TitleBar: application.MacTitleBar{
            AppearsTransparent: true,
        },
        Backdrop: application.MacBackdropTranslucent,
    },
})
```

**Types d’arrière-plan :**

- `MacBackdropNormal` - Fenêtre standard
- `MacBackdropTranslucent` - Arrière-plan translucide ; **une API privée est requise** pour rendre la WebView transparente.
- `MacBackdropTransparent` - Entièrement transparent ; **une API privée est requise** pour rendre la WebView transparente.
- `MacBackdropLiquidGlass` - Arrière-plan en verre ; **une API privée est requise** pour rendre la WebView transparente.

Effectuez la compilation avec `-tags private_mac_apis` pour que ces effets soient visibles à travers la WebView. Sans cette option, la WebView reste opaque. `TitleBar.AppearsTransparent` utilise elle-même des API publiques. Consultez [API privées de macOS](/guides/build/private-macos-apis/).

**Comportement de regroupement :** Contrôlez le comportement des fenêtres dans les différents Spaces :

- `MacWindowCollectionBehaviorCanJoinAllSpaces` - Visible dans tous les Spaces
- `MacWindowCollectionBehaviorFullScreenAuxiliary` - Peut se superposer aux applications en plein écran

**Plein écran natif :** Le mode plein écran de macOS crée un nouveau Space (bureau virtuel).

[Linux]
**Fonctionnalités propres à Linux :**

```go
// Set window icon (per-window struct is LinuxWindow, not the app-level LinuxOptions)
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Linux: application.LinuxWindow{
        Icon: iconBytes,
    },
})
```

**Remarques sur les environnements de bureau :**

- GNOME : prise en charge complète
- KDE Plasma : prise en charge complète
- XFCE : prise en charge partielle
- Autres : prise en charge variable

**Gestionnaires de fenêtres en mosaïque (Hyprland, Sway, i3, etc.) :**

- `Minimise()` et `Maximise()` peuvent ne pas fonctionner comme prévu, car le gestionnaire de fenêtres contrôle la géométrie des fenêtres
- Les requêtes `SetSize()` et `SetPosition()` sont données à titre indicatif et peuvent être ignorées
- `Fullscreen()` fonctionne généralement comme prévu
- Certains gestionnaires de fenêtres ne prennent pas en charge le mode toujours au premier plan

@end

## Modèles courants

### Écran de démarrage

```go
// Create splash screen
splash := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:     "Loading...",
    Width:     400,
    Height:    300,
    Frameless: true,
    AlwaysOnTop: true,
})

// Show splash
splash.Show()

// Initialise application
time.Sleep(2 * time.Second)

// Hide splash, show main window
splash.Close()
mainWindow.Show()
```

### Fenêtre des paramètres

```go
var settingsWindow *application.WebviewWindow

func showSettings() {
    if settingsWindow == nil {
        settingsWindow = app.Window.NewWithOptions(application.WebviewWindowOptions{
            Name:   "settings",
            Title:  "Settings",
            Width:  600,
            Height: 400,
        })
    }
    
    settingsWindow.Show()
    settingsWindow.Focus()
}
```

### Demander confirmation avant la fermeture

```go
window.RegisterHook(events.Common.WindowClosing, func(event *application.WindowEvent) {
    if hasUnsavedChanges() {
        // Show dialog
        result := showConfirmDialog("Unsaved changes. Close anyway?")
        if result != "yes" {
            // Cancel the close event
            event.Cancel()
        }
    }
})
```

## Bonnes pratiques

### ✅ À faire

- **Nommez les fenêtres importantes** – Elles seront plus faciles à retrouver par la suite
- **Définissez une taille minimale** – Vous éviterez les mises en page inutilisables
- **Centrez les fenêtres** – L’expérience utilisateur est meilleure qu’avec une position aléatoire
- **Gérez les événements de fermeture** – Vous éviterez les pertes de données
- **Testez sur toutes les plateformes** – Le comportement varie selon la plateforme
- **Utilisez des dimensions adaptées** – Tenez compte des différentes tailles d’écran

### ❌ À éviter

- **Ne créez pas trop de fenêtres** – Cela déroute les utilisateurs
- **N’oubliez pas de fermer les fenêtres** – Vous éviterez les fuites de mémoire
- **Ne codez pas les positions en dur** – Les tailles d’écran diffèrent
- **N’ignorez pas les différences entre les plateformes** – Effectuez des tests approfondis
- **Ne bloquez pas le thread de l’interface utilisateur** – Utilisez des goroutines pour les opérations longues

## Dépannage

### La fenêtre ne s’affiche pas

**Causes possibles :**

1. La fenêtre a été créée en mode masqué
2. La fenêtre se trouve hors de l’écran
3. La fenêtre se trouve derrière d’autres fenêtres

**Solution :**

```go
window.Show()
window.Center()
window.Focus()
```

### La fenêtre n’a pas la bonne taille

**Cause :** mise à l’échelle selon la résolution (DPI) sous Windows/Linux

**Solution :**

```go
// Wails handles DPI automatically
// Just use logical pixels
window.SetSize(800, 600)
```

### La fenêtre se ferme immédiatement

**Cause :** l’application se ferme lorsque la dernière fenêtre est fermée

**Solution :**

```go
app := application.New(application.Options{
    Mac: application.MacOptions{
        ApplicationShouldTerminateAfterLastWindowClosed: false,
    },
})
```

## Étapes suivantes

@cards{cols="2"}
⚙ Options de fenêtre
Référence complète de toutes les options de fenêtre.

[En savoir plus →](/features/windows/options/)

---
▣ Fenêtres multiples
Modèles pour les applications à plusieurs fenêtres.

[En savoir plus →](/features/windows/multiple/)

---
★ Fenêtres sans cadre
Créez une décoration de fenêtre personnalisée.

[En savoir plus →](/features/windows/frameless/)

---
🚀 Événements de fenêtre
Gérez les événements du cycle de vie des fenêtres.

[En savoir plus →](/features/windows/events/)

@end

---

**Des questions ?** Posez-les sur [Discord](https://discord.gg/JDdSxwjhGf) ou consultez les [exemples de fenêtres](https://github.com/wailsapp/wails/tree/master/v3/examples).
