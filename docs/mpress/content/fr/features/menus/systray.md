---
title: "Menus de la zone de notification"
description: "Ajoutez l’intégration à la zone de notification de votre système à votre application"
slug: "features/menus/systray"
sourcePath: "features/menus/systray.md"
---

## Menus de la zone de notification

Wails fournit des **API unifiées pour la zone de notification** qui fonctionnent sur toutes les plateformes. Créez des icônes avec des menus dans la zone de notification, associez-y des fenêtres et gérez les clics avec le comportement natif de chaque plateforme pour les applications en arrière-plan, les services et les utilitaires à accès rapide.

![Menu de la zone de notification Wails ouvert depuis la barre des menus de macOS](/assets/screenshots/systray-menu-macos.png)

Sous macOS, un élément Wails de la zone de notification apparaît dans la barre des menus et ouvre un menu natif. L’exemple comprend des éléments désactivés, des cases à cocher, des boutons radio, un sous-menu et des actions.

## Démarrage rapide

```go
package main

import (
    _ "embed"
    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets/icon.png
var icon []byte

func main() {
    app := application.New(application.Options{
        Name: "Tray App",
    })

    // Create system tray
    systray := app.SystemTray.New()
    systray.SetIcon(icon)
    systray.SetLabel("My App")

    // Add menu
    menu := app.NewMenu()
    menu.Add("Show").OnClick(func(ctx *application.Context) {
        // Show main window
    })
    menu.Add("Quit").OnClick(func(ctx *application.Context) {
        app.Quit()
    })
    systray.SetMenu(menu)

    // Create hidden window
    window := app.Window.New()
    window.Hide()

    app.Run()
}
```

**Résultat :** une icône avec un menu dans la zone de notification sur toutes les plateformes.

## Création d’un élément de la zone de notification

### Élément de base de la zone de notification

```go
// Create system tray
systray := app.SystemTray.New()

// Set icon
systray.SetIcon(iconBytes)

// Set label (macOS) / tooltip (Windows)
systray.SetLabel("My Application")
```

### Avec une icône

Les icônes devraient être incorporées :

```go
import _ "embed"

//go:embed assets/icon.png
var icon []byte

//go:embed assets/icon-dark.png
var iconDark []byte

func main() {
    app := application.New(application.Options{
        Name: "My App",
    })

    systray := app.SystemTray.New()
    systray.SetIcon(icon)
    systray.SetDarkModeIcon(iconDark)  // Windows and macOS dark mode
    
    app.Run()
}
```

**Exigences relatives aux icônes :**

| Plateforme | Taille | Format | Remarques |
| --- | --- | --- | --- |
| **Windows** | 16x16 ou 32x32 | PNG, ICO | Zone de notification |
| **macOS** | 18x18 à 22x22 | PNG | Barre des menus, modèle recommandé |
| **Linux** | 22x22 à 48x48 | PNG, SVG | Variable selon l’environnement de bureau |

### Icônes modèles (macOS)

Les icônes modèles s’adaptent automatiquement au mode clair ou sombre :

```go
systray.SetTemplateIcon(iconBytes)
```

**Recommandations pour les icônes modèles :**

- Utilisez uniquement du noir et des couleurs transparentes
- Le noir devient blanc en mode sombre
- Ajoutez le suffixe `Template` au nom du fichier : `iconTemplate.png`
- [Guide de conception](https://bjango.com/articles/designingmenubarextras/)

## Ajout de menus

Les menus de la zone de notification fonctionnent comme les menus de l’application :

```go
menu := app.NewMenu()

// Add items
menu.Add("Open").OnClick(func(ctx *application.Context) {
    showMainWindow()
})

menu.AddSeparator()

menu.AddCheckbox("Start at Login", false).OnClick(func(ctx *application.Context) {
    enabled := ctx.ClickedMenuItem().Checked()
    setStartAtLogin(enabled)
})

menu.AddSeparator()

menu.Add("Quit").OnClick(func(ctx *application.Context) {
    app.Quit()
})

// Set menu
systray.SetMenu(menu)
```

Pour connaître **tous les types d’éléments de menu**, consultez la [référence des menus](/features/menus/reference/).

## Association de fenêtres

Associez une fenêtre à l’icône de la zone de notification pour l’afficher et la masquer automatiquement :

```go
// Create window
window := app.Window.New()

// Attach to tray
systray.AttachWindow(window)

// Configure behaviour — these are setters that return the receiver for chaining.
systray.WindowOffset(10)                          // Pixels from tray icon
systray.WindowDebounce(200 * time.Millisecond)    // Click debounce
```

**Comportement :**

- La fenêtre est initialement masquée
- **Clic gauche sur l’icône de la zone de notification** → Afficher ou masquer la fenêtre
- **Clic droit sur l’icône de la zone de notification** → Afficher le menu, s’il est défini
- La fenêtre est placée près de l’icône de la zone de notification

**Exemple : fenêtre contextuelle**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:           "Quick Access",
    Width:           300,
    Height:          400,
    Frameless:       true, // No title bar
    AlwaysOnTop:     true, // Stay on top
    HideOnFocusLost: true, // Dismiss when another window receives focus
    HideOnEscape:    true, // Dismiss when the user presses Escape
})

systray.AttachWindow(window)
systray.WindowOffset(5)
```

`HideOnFocusLost` est utile pour les fenêtres contextuelles de la zone de notification sous Windows, macOS et les environnements de bureau Linux où le focus est attribué par clic. Wails désactive ce comportement sous les environnements Linux où le focus suit la souris, notamment dans les configurations courantes de Hyprland, Sway et i3, car le simple fait de quitter la fenêtre contextuelle pourrait sinon la masquer avant de pouvoir l’utiliser. `HideOnEscape` reste disponible dans ces environnements.

Les comportements de clic gauche et de clic droit décrits ci-dessus sont des valeurs par défaut intelligentes. Un gestionnaire `OnClick` ou `OnRightClick` explicite remplace le comportement par défaut correspondant. Pour les vérifications propres aux plateformes et les cas limites, consultez la [suite de tests manuels de la zone de notification](https://github.com/wailsapp/wails/tree/master/v3/test/manual/systray) et l’[exemple de test de charge de la zone de notification](https://github.com/wailsapp/wails/tree/master/v3/examples/systray-stress).

## Gestionnaires de clics

Gérez les clics sur l’icône de la zone de notification :

```go
systray := app.SystemTray.New()

// Left click
systray.OnClick(func() {
    fmt.Println("Tray icon clicked")
})

// Right click
systray.OnRightClick(func() {
    fmt.Println("Tray icon right-clicked")
})

// Double click
systray.OnDoubleClick(func() {
    fmt.Println("Tray icon double-clicked")
})

// Mouse enter/leave
systray.OnMouseEnter(func() {
    fmt.Println("Mouse entered tray icon")
})

systray.OnMouseLeave(func() {
    fmt.Println("Mouse left tray icon")
})
```

**Prise en charge selon la plateforme :**

| Événement | Windows | macOS | Linux |
| --- | --- | --- | --- |
| OnClick | ✅ | ✅ | ✅ |
| OnRightClick | ✅ | ✅ | ✅ |
| OnDoubleClick | ✅ | ✅ | ⚠️ Variable |
| OnMouseEnter | ✅ | ✅ | ⚠️ Variable |
| OnMouseLeave | ✅ | ✅ | ⚠️ Variable |

## Mises à jour dynamiques

Mettez à jour dynamiquement l’icône et le menu de la zone de notification :

### Modification de l’icône

```go
var isActive bool

func updateTrayIcon() {
    if isActive {
        systray.SetIcon(activeIcon)
        systray.SetLabel("Active")
    } else {
        systray.SetIcon(inactiveIcon)
        systray.SetLabel("Inactive")
    }
}
```

### Mise à jour du menu

```go
var isPaused bool

pauseMenuItem := menu.Add("Pause")

pauseMenuItem.OnClick(func(ctx *application.Context) {
    isPaused = !isPaused
    
    if isPaused {
        pauseMenuItem.SetLabel("Resume")
    } else {
        pauseMenuItem.SetLabel("Pause")
    }
    
    menu.Update()  // Important!
})
```

@note{type="caution" title="Toujours appeler Update()"}
Après avoir modifié l’état du menu, **appelez `menu.Update()`**. Consultez la [référence des menus](/features/menus/reference/#enabled-state).

@end

### Reconstruction du menu

Pour les modifications importantes, reconstruisez l’intégralité du menu :

```go
func rebuildTrayMenu(status string) {
    menu := app.NewMenu()
    
    // Status-specific items
    switch status {
    case "syncing":
        menu.Add("Syncing...").SetEnabled(false)
        menu.Add("Pause Sync").OnClick(pauseSync)
    case "synced":
        menu.Add("Up to date ✓").SetEnabled(false)
        menu.Add("Sync Now").OnClick(startSync)
    case "error":
        menu.Add("Sync Error").SetEnabled(false)
        menu.Add("Retry").OnClick(retrySync)
    }
    
    menu.AddSeparator()
    menu.Add("Quit").OnClick(func(ctx *application.Context) {
        app.Quit()
    })
    
    systray.SetMenu(menu)
}
```

## Fonctionnalités propres à chaque plateforme

@tabs{sync-key="platform"}
[macOS]
**Intégration à la barre des menus :**

```go
// Set label (appears next to icon)
systray.SetLabel("My App")

// Use template icon (adapts to dark mode)
systray.SetTemplateIcon(iconBytes)

// Set icon position — uses AppKit NSImage placement constants.
systray.SetIconPosition(application.NSImageRight)
```

**Positions de l’icône** (correspondant à `NSImagePosition`) :

- `application.NSImageLeft` — Icône à gauche du libellé.
- `application.NSImageRight` — Icône à droite du libellé.
- `application.NSImageOnly` — Icône uniquement, sans libellé.
- `application.NSImageNone` — Libellé uniquement, sans icône.

**Bonnes pratiques :**

- Utilisez des icônes de modèle (noires avec transparence)
- Utilisez des libellés courts (3-5 caractères)
- 18x18 à 22x22 pixels pour les écrans Retina
- Testez dans les modes clair et sombre

[Windows]
**Intégration à la zone de notification :**

```go
// Set tooltip (appears on hover)
systray.SetTooltip("My Application")

// Or use SetLabel (same as tooltip on Windows)
systray.SetLabel("My Application")

// Show/Hide functionality (fully functional)
systray.Show()  // Show tray icon
systray.Hide()  // Hide tray icon
```

**Exigences relatives à l’icône :**

- 16x16 ou 32x32 pixels
- Format PNG ou ICO
- Arrière-plan transparent

**Limites des infobulles :**

- 127 caractères UTF-16 au maximum
- Les infobulles plus longues seront tronquées
- Restez concis pour une expérience optimale

**Fonctionnalités de la plateforme :**

- L’icône de la zone de notification est conservée après le redémarrage de l’Explorateur Windows
- Méthodes Show() et Hide() entièrement fonctionnelles
- Gestion correcte du cycle de vie

**Bonnes pratiques :**

- Utilisez 32x32 pour les écrans à haute densité de pixels
- Limitez les infobulles à moins de 127 caractères
- Testez sur différentes versions de Windows
- Tenez compte du débordement de la zone de notification
- Utilisez Show/Hide pour afficher ou masquer conditionnellement l’icône de la zone de notification

[Linux]
**Intégration à la zone de notification :**

Utilise la spécification StatusNotifierItem (dans la plupart des environnements de bureau modernes).

```go
systray.SetIcon(iconBytes)
systray.SetLabel("My App")
```

**Prise en charge des environnements de bureau :**

- **GNOME** : barre supérieure (avec une extension)
- **KDE Plasma** : zone de notification
- **XFCE** : zone de notification
- **Autres** : variable

**Bonnes pratiques :**

- Utilisez des dimensions de 22x22 ou 24x24 pixels
- Les icônes SVG s’adaptent mieux aux différentes tailles
- Testez sur les environnements de bureau ciblés
- Prévoyez une solution de secours pour les environnements de bureau non pris en charge

@end

## Exemple complet

Voici une application avec zone de notification prête pour la production :

```go
package main

import (
    _ "embed"
    "fmt"
    "time"
    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets/icon.png
var icon []byte

//go:embed assets/icon-active.png
var iconActive []byte

type TrayApp struct {
    app     *application.App
    systray *application.SystemTray
    window  *application.WebviewWindow
    menu    *application.Menu
    isActive bool
}

func main() {
    app := application.New(application.Options{
        Name: "Tray Application",
        Mac: application.MacOptions{
            ApplicationShouldTerminateAfterLastWindowClosed: false,
        },
    })

    trayApp := &TrayApp{app: app}
    trayApp.setup()

    app.Run()
}

func (t *TrayApp) setup() {
    // Create system tray
    t.systray = t.app.SystemTray.New()
    t.systray.SetIcon(icon)
    t.systray.SetLabel("Inactive")
    
    // Create menu
    t.createMenu()
    
    // Create window (hidden by default)
    t.window = t.app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:  "Tray Application",
        Width:  400,
        Height: 600,
        Hidden: true,
    })
    
    // Attach window to tray
    t.systray.AttachWindow(t.window)
    t.systray.WindowOffset(10)
    
    // Handle tray clicks
    t.systray.OnRightClick(func() {
        t.systray.OpenMenu()
    })
    
    // Start background task
    go t.backgroundTask()
}

func (t *TrayApp) createMenu() {
    t.menu = t.app.NewMenu()
    
    // Status item (disabled)
    statusItem := t.menu.Add("Status: Inactive")
    statusItem.SetEnabled(false)
    
    t.menu.AddSeparator()
    
    // Toggle active
    t.menu.Add("Start").OnClick(func(ctx *application.Context) {
        t.toggleActive()
    })
    
    // Show window
    t.menu.Add("Show Window").OnClick(func(ctx *application.Context) {
        t.window.Show()
        t.window.Focus()
    })
    
    t.menu.AddSeparator()
    
    // Settings
    t.menu.AddCheckbox("Start at Login", false).OnClick(func(ctx *application.Context) {
        enabled := ctx.ClickedMenuItem().Checked()
        t.setStartAtLogin(enabled)
    })
    
    t.menu.AddSeparator()
    
    // Quit
    t.menu.Add("Quit").OnClick(func(ctx *application.Context) {
        t.app.Quit()
    })
    
    t.systray.SetMenu(t.menu)
}

func (t *TrayApp) toggleActive() {
    t.isActive = !t.isActive
    t.updateTray()
}

func (t *TrayApp) updateTray() {
    if t.isActive {
        t.systray.SetIcon(iconActive)
        t.systray.SetLabel("Active")
    } else {
        t.systray.SetIcon(icon)
        t.systray.SetLabel("Inactive")
    }
    
    // Rebuild menu with new status
    t.createMenu()
}

func (t *TrayApp) backgroundTask() {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        if t.isActive {
            fmt.Println("Background task running...")
            // Do work
        }
    }
}

func (t *TrayApp) setStartAtLogin(enabled bool) {
    // Implementation varies by platform
    fmt.Printf("Start at login: %v\n", enabled)
}
```

## Contrôle de la visibilité

Affichez ou masquez dynamiquement l’icône de la zone de notification :

```go
// Hide tray icon
systray.Hide()

// Show tray icon
systray.Show()
```

Il n’existe aucun accesseur `IsVisible()` ; si nécessaire, suivez la visibilité dans l’état de votre propre application.

**Prise en charge des plateformes :**

| Plateforme | Hide() | Show() | Remarques |
| --- | --- | --- | --- |
| **Windows** | ✅ | ✅ | Entièrement fonctionnel — l’icône apparaît dans la zone de notification ou en disparaît |
| **macOS** | ✅ | ✅ | L’élément de la barre des menus s’affiche ou se masque |
| **Linux** | ✅ | ✅ | Varie selon l’environnement de bureau |

**Cas d’utilisation :**

- Masquez temporairement l’icône de la zone de notification selon les préférences de l’utilisateur
- Utilisez un mode sans interface graphique dans lequel l’icône de la zone de notification n’apparaît qu’en cas de besoin
- Basculez la visibilité selon l’état de l’application

**Exemple — Visibilité conditionnelle de l’icône de la zone de notification :**

```go
func (t *TrayApp) setTrayVisibility(visible bool) {
    if visible {
        t.systray.Show()
    } else {
        t.systray.Hide()
    }
}

// Show tray only when updates are available
func (t *TrayApp) checkForUpdates() {
    if hasUpdates {
        t.systray.Show()
        t.systray.SetLabel("Update Available")
    } else {
        t.systray.Hide()
    }
}
```

## Nettoyage

Détruisez l’icône de la zone de notification lorsque vous n’en avez plus besoin :

```go
// In OnShutdown
app := application.New(application.Options{
    OnShutdown: func() {
        if systray != nil {
            systray.Destroy()
        }
    },
})
```

**Important :** détruisez toujours l’icône de la zone de notification lors de l’arrêt afin de libérer les ressources.

## Bonnes pratiques

### ✅ À faire

- **Utilisez des icônes de modèle sur macOS** — Elles s’adaptent au mode sombre
- **Utilisez des libellés courts** — 3-5 caractères maximum
- **Ajoutez des info-bulles sous Windows** — Elles aident les utilisateurs à identifier votre application
- **Testez sur toutes les plateformes** — Le comportement varie
- **Gérez les clics de manière appropriée** — Clic gauche pour l’action principale, clic droit pour le menu
- **Actualisez l’icône en fonction de l’état** — Le retour visuel est important
- **Détruisez l’icône de notification à l’arrêt** — Libérez les ressources

### ❌ À ne pas faire

- **N’utilisez pas de grandes icônes** — Respectez les recommandations de la plateforme
- **N’utilisez pas de libellés longs** — Ils sont tronqués
- **N’oubliez pas le mode sombre** — Testez sous Windows et macOS en mode sombre
- **Ne bloquez pas les gestionnaires de clics** — Veillez à ce qu’ils s’exécutent rapidement
- **N’oubliez pas menu.Update()** — Après avoir modifié l’état du menu
- **Ne supposez pas que la zone de notification est prise en charge** — Certains environnements de bureau Linux ne la prennent pas en charge

## Dépannage

### L’icône n’apparaît pas dans la zone de notification

**Causes possibles :**

1. Format d’icône non pris en charge
2. Icône trop grande ou trop petite
3. Zone de notification non prise en charge (Linux)

**Solution :**

Il n’existe aucun utilitaire `SystemTraySupported()` ; créez plutôt la zone de notification, vérifiez la plateforme et prévoyez un fonctionnement dégradé :

```go
// Probe support: on Linux without a notification-area extension, the tray
// will simply not appear. Defensive code can fall back to window-only mode
// based on runtime.GOOS or after a short timeout if no tray events arrive.
systray := app.SystemTray.New()
systray.SetIcon(iconBytes)
```

### L’icône s’affiche incorrectement sous macOS

**Cause :** aucune icône de modèle n’est utilisée

**Solution :**

```go
// Use template icon
systray.SetTemplateIcon(iconBytes)

// Or design icon as template (black + transparent)
```

### Le menu ne s’actualise pas

**Cause :** oubli de l’appel à `menu.Update()`

**Solution :**

```go
menuItem.SetLabel("New Label")
menu.Update()  // Add this!
```

## Étapes suivantes

@cards{cols="2"}
📖 Référence des menus
Référence complète des types et propriétés des éléments de menu.

[En savoir plus →](/features/menus/reference/)

---
☰ Menus de l’application
Créez des barres de menus pour l’application.

[En savoir plus →](/features/menus/application/)

---
◆ Menus contextuels
Créez des menus contextuels accessibles par clic droit.

[En savoir plus →](/features/menus/context/)

---
📖 Exemple de zone de notification
Explorez une application complète utilisant la zone de notification.

[En savoir plus →](https://github.com/wailsapp/wails/tree/master/v3/examples/systray-basic)

@end

---

**Des questions ?** Posez-les sur [Discord](https://discord.gg/JDdSxwjhGf) ou consultez les [exemples de zones de notification](https://github.com/wailsapp/wails/tree/master/v3/examples/systray-basic).
