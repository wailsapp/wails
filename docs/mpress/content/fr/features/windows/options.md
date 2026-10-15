---
title: "Options de fenêtre"
description: "Référence complète de WebviewWindowOptions"
slug: "features/windows/options"
sourcePath: "features/windows/options.md"
---

## Options de configuration des fenêtres

Wails fournit une configuration complète des fenêtres, avec des dizaines d’options relatives à leur taille, leur position, leur apparence et leur comportement. Cette page constitue la **référence complète** de `WebviewWindowOptions` et couvre toutes les options disponibles sous Windows, macOS et Linux. Vous y trouverez chaque option pour chaque plateforme, accompagnée d’exemples et de contraintes.

## Structure WebviewWindowOptions

```go
type WebviewWindowOptions struct {
    // Identity
    Name  string
    Title string

    // Size and Position
    Width           int
    Height          int
    X               int
    Y               int
    MinWidth        int
    MinHeight       int
    MaxWidth        int
    MaxHeight       int
    InitialPosition WindowStartPosition // WindowCentered (default) or WindowXY
    Screen          *Screen             // target screen for initial placement

    // Initial State
    Hidden        bool
    Frameless     bool
    DisableResize bool        // inverted vs v2's `Resizable`
    AlwaysOnTop   bool
    StartState    WindowState // WindowStateNormal | Minimised | Maximised | Fullscreen

    // Appearance
    BackgroundColour RGBA
    BackgroundType   BackgroundType
    Zoom             float64
    ZoomControlEnabled bool

    // Content
    URL  string
    HTML string
    JS   string
    CSS  string

    // Behaviour
    EnableFileDrop              bool
    IgnoreMouseEvents           bool
    HideOnFocusLost             bool
    HideOnEscape                bool
    DevToolsEnabled             bool
    DefaultContextMenuDisabled  bool
    ContentProtectionEnabled    bool
    KeyBindings                 map[string]func(window *WebviewWindow)

    // Permissions
    Permissions map[PermissionType]Permission

    // Window-control button states
    MinimiseButtonState ButtonState
    MaximiseButtonState ButtonState
    CloseButtonState    ButtonState

    // Menu
    UseApplicationMenu bool

    // Platform-specific (per-window)
    Mac     MacWindow
    Windows WindowsWindow
    Linux   LinuxWindow
}
```

`WebviewWindowOptions` ne comporte **aucun** champ `Parent` : pour les relations parent/modale, utilisez `parentWindow.AttachModal(childWindow)`. Il ne comporte **pas non plus** de champ `Assets` : la configuration des ressources se trouve dans `application.Options` (`Assets AssetOptions`).

Code source complet : [`v3/pkg/application/webview_window_options.go`](https://github.com/wailsapp/wails/blob/master/v3/pkg/application/webview_window_options.go).

## Options principales

### Nom

**Type :** `string` **Valeur par défaut :** UUID généré automatiquement **Plateforme :** toutes

```go
Name: "main-window"
```

**Objectif :** identifiant unique permettant de retrouver les fenêtres ultérieurement.

**Bonnes pratiques :**

- Utilisez des noms descriptifs : `"main"`, `"settings"`, `"about"`
- Utilisez le format kebab-case : `"file-browser"`, `"color-picker"`
- Choisissez un nom court et facile à retenir

**Exemple :**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name: "settings-window",
})

// Later...
if settings, ok := app.Window.GetByName("settings-window"); ok {
    settings.Focus()
}
```

### Titre

**Type :** `string` **Valeur par défaut :** nom de l’application **Plateforme :** toutes

```go
Title: "My Application"
```

**Objectif :** texte affiché dans la barre de titre et la barre des tâches.

**Mises à jour dynamiques :**

```go
window.SetTitle("My Application - Document.txt")
```

### Largeur / Hauteur

**Type :** `int` (pixels) **Valeur par défaut :** 800 × 600 **Plateforme :** toutes **Contraintes :** la valeur doit être positive

```go
Width:  1200,
Height: 800,
```

**Objectif :** taille initiale de la fenêtre en pixels logiques.

**Remarques :**

- Wails gère automatiquement la mise à l’échelle selon la résolution (DPI)
- Utilisez des pixels logiques, et non des pixels physiques
- Tenez compte de la résolution d’écran minimale (1024x768)

**Exemples de dimensions :**

| Cas d’utilisation | Largeur | Hauteur |
| --- | --- | --- |
| Petit utilitaire | 400 | 300 |
| Application standard | 1024 | 768 |
| Grande application | 1440 | 900 |
| Full HD | 1920 | 1080 |

### X / Y

**Type :** `int` (pixels) **Valeur par défaut :** centrée à l’écran **Plateforme :** toutes

```go
X: 100,  // 100px from left edge
Y: 100,  // 100px from top edge
```

**Objectif :** position initiale de la fenêtre.

**Système de coordonnées :**

- Le point (0, 0) correspond au coin supérieur gauche de l’écran principal
- Les valeurs positives de X vont vers la droite
- Les valeurs positives de Y vont vers le bas

**Exemple :**

`X` et `Y` ne prennent effet que si `InitialPosition: application.WindowXY` est défini. Dans le cas contraire, `InitialPosition` prend par défaut la valeur `WindowCentered`, et `X`/`Y` sont ignorés.

```go
settings := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:            "coordinate-window",
    InitialPosition: application.WindowXY, // opt into X/Y coordinates
    X:               100,
    Y:               100,
})
```

**Bonne pratique :** si des coordonnées précises ne sont pas nécessaires, utilisez `Center()` pour centrer une fenêtre après sa création :

```go
window := app.Window.New()
window.Center()
```

### MinWidth / MinHeight

**Type :** `int` (pixels) **Valeur par défaut :** 0 (aucune valeur minimale) **Plateforme :** toutes

```go
MinWidth:  400,
MinHeight: 300,
```

**Objectif :** empêcher la fenêtre de devenir trop petite.

**Cas d’utilisation :**

- Éviter que la mise en page ne soit détériorée
- Garantir la facilité d’utilisation
- Conserver les proportions

**Exemple :**

```go
// Prevent window smaller than 400x300
MinWidth:  400,
MinHeight: 300,
```

### MaxWidth / MaxHeight

**Type :** `int` (pixels) **Valeur par défaut :** 0 (aucune valeur maximale) **Plateforme :** toutes

```go
MaxWidth:  1920,
MaxHeight: 1080,
```

**Objectif :** empêcher la fenêtre de devenir trop grande.

**Cas d’utilisation :**

- Applications de taille fixe
- Éviter une utilisation excessive des ressources
- Respecter les contraintes de conception

## Options d’état

### Hidden

**Type :** `bool` **Valeur par défaut :** `false` **Plateformes :** toutes

```go
Hidden: true,
```

**Objectif :** créer la fenêtre sans l’afficher.

**Cas d’utilisation :**

- Fenêtres en arrière-plan
- Fenêtres affichées à la demande
- Écrans de démarrage (créer, charger, puis afficher)
- Éviter un flash blanc pendant le chargement du contenu

**Améliorations selon la plateforme :**

- **Windows :** le flash blanc de la fenêtre a été corrigé ; la fenêtre reste invisible jusqu’à l’appel de `Show()`
- **macOS :** prise en charge complète
- **Linux :** prise en charge complète

**Approche recommandée pour un chargement fluide :**

```go
// Create hidden window
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:             "main-window",
    Hidden:           true,
    BackgroundColour: application.NewRGB(30, 30, 30), // Match your theme
})

// Load content while hidden
// ... content loads ...

// Show when ready (no flash!)
window.Show()
```

**Exemple :**

```go
settings := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:   "settings",
    Hidden: true,
})

// Show when needed
settings.Show()
```

### Frameless

**Type :** `bool` **Valeur par défaut :** `false` **Plateformes :** toutes

```go
Frameless: true,
```

**Objectif :** supprimer la barre de titre et les bordures de la fenêtre.

**Cas d’utilisation :**

- Décoration de fenêtre personnalisée
- Écrans de démarrage
- Applications en mode kiosque
- Fenêtres au design personnalisé

**Important :** vous devrez implémenter :

- Le déplacement de la fenêtre
- Les boutons de fermeture, de réduction et d’agrandissement
- Les poignées de redimensionnement (si la fenêtre est redimensionnable)

**Pour plus de détails, consultez la page [Fenêtres sans cadre](/features/windows/frameless/).**

### DisableResize

**Type :** `bool` **Valeur par défaut :** `false` (la fenêtre est redimensionnable par défaut) **Plateformes :** toutes

```go
DisableResize: true,
```

**Objectif :** empêcher le redimensionnement de la fenêtre. Notez que ce champ est l’**inverse** de `Resizable` dans la v2 : définissez `DisableResize: true` pour rendre une fenêtre non redimensionnable.

**Cas d’utilisation :**

- Applications de taille fixe
- Écrans de démarrage
- Boîtes de dialogue

**Remarque :** les utilisateurs peuvent toujours agrandir la fenêtre ou passer en plein écran, sauf si vous désactivez également ces possibilités via `MaximiseButtonState` ou le comportement de collection.

### AlwaysOnTop

**Type :** `bool` **Valeur par défaut :** `false` **Plateformes :** toutes

```go
AlwaysOnTop: true,
```

**Objectif :** maintenir la fenêtre au-dessus de toutes les autres fenêtres.

**Cas d’utilisation :**

- Barres d’outils flottantes
- Notifications
- Mode image dans l’image
- Minuteurs

**Remarques selon la plateforme :**

- **macOS :** prise en charge complète
- **Windows :** prise en charge complète
- **Linux :** dépend du gestionnaire de fenêtres

### StartState

**Type :** énumération `WindowState` **Valeur par défaut :** `WindowStateNormal` **Plateformes :** toutes

```go
StartState: application.WindowStateMaximised,
```

**Objectif :** définir l’état initial de la fenêtre lors de son affichage.

**Valeurs :**

- `WindowStateNormal` - Fenêtre normale
- `WindowStateMinimised` - Réduite
- `WindowStateMaximised` - Agrandie
- `WindowStateFullscreen` - Plein écran

Il n’existe aucune constante `WindowStateHidden` : utilisez le champ booléen `Hidden` pour que la fenêtre soit invisible au démarrage.

**Basculer en mode plein écran lors de l’exécution :**

```go
window.Fullscreen()
window.UnFullscreen()
window.ToggleFullscreen() // there is no SetFullscreen(bool)
```

## Options d’apparence

### BackgroundColour

**Type :** structure `RGBA` **Valeur par défaut :** blanc **Plateforme :** toutes

```go
BackgroundColour: application.RGBA{Red: 0, Green: 0, Blue: 0, Alpha: 255},
```

Les champs de `RGBA` sont `Red, Green, Blue, Alpha`, chacun de type uint8. Préférez les fonctions d’assistance `application.NewRGB(r, g, b)` (alpha 255) ou `application.NewRGBA(r, g, b, a)`.

**Objectif :** définir la couleur d’arrière-plan de la fenêtre avant le chargement du contenu.

**Cas d’utilisation :**

- Reprendre le thème de votre application
- Éviter un flash blanc avec les thèmes sombres
- Assurer un chargement fluide

**Exemple :**

```go
// Dark theme
BackgroundColour: application.NewRGB(30, 30, 30),

// Light theme
BackgroundColour: application.NewRGB(255, 255, 255),
```

**Méthode d’assistance :**

```go
window.SetBackgroundColour(application.NewRGB(30, 30, 30))
```

### BackgroundType

**Type :** énumération `BackgroundType` **Valeur par défaut :** `BackgroundTypeSolid` **Plateforme :** macOS, Windows (prise en charge partielle)

```go
BackgroundType: application.BackgroundTypeTranslucent,
```

**Valeurs :**

- `BackgroundTypeSolid` — couleur unie
- `BackgroundTypeTransparent` — entièrement transparent
- `BackgroundTypeTranslucent` — flou semi-transparent

**Prise en charge par plateforme :**

- **macOS :** configurez `Mac.Backdrop` ; la transparence de la vue web nécessite [`private_mac_apis`](/guides/build/private-macos-apis/#webview-transparency-and-background). Sans cette option, la vue web reste opaque.
- **Windows :** transparent et translucide (Windows 11 ou version ultérieure)
- **Linux :** couleur unie uniquement

**Exemple (macOS) :**

```go
BackgroundType: application.BackgroundTypeTranslucent,
Mac: application.MacWindow{
    Backdrop: application.MacBackdropTranslucent,
},
```

### OpenInspectorOnStartup et OpenDevTools

**API privée sous macOS :** `OpenInspectorOnStartup: true`, la méthode Go `window.OpenDevTools()` et la méthode JavaScript `Window.OpenDevTools()` nécessitent `private_mac_apis` pour ouvrir l’inspecteur par programmation. Sans cette option, ces opérations sont sans effet. Les versions de production nécessitent également `devtools`. L’inspection publique de Safari sous macOS 13.3 ou version ultérieure ne nécessite pas d’API privée, contrairement à l’activation de l’inspecteur sous les versions antérieures de macOS. Consultez la [matrice de compilation de l’inspecteur web](/guides/build/private-macos-apis/#web-inspector).

## Options de contenu

### URL

**Type :** `string` **Valeur par défaut :** vide (chargement depuis Assets) **Plateforme :** toutes

```go
URL: "https://example.com",
```

**Objectif :** charger une URL externe à la place des ressources intégrées.

**Cas d’utilisation :**

- Développement (chargement depuis le serveur de développement)
- Applications web
- Applications hybrides

**Exemple :**

```go
// Development — point the window at the Vite dev server
URL: "http://localhost:9245",

// Production — embedded assets are configured at the application level
// (Assets is application.Options.Assets, not a WebviewWindowOptions field).
```

Extrait de code au niveau de l’application pour le cas de la production :

```go
app := application.New(application.Options{
    Name: "My App",
    Assets: application.AssetOptions{
        Handler: application.AssetFileServerFS(assets),
    },
})
```

### HTML

**Type :** `string` **Valeur par défaut :** vide **Plateforme :** toutes

```go
HTML: "<h1>Hello World</h1>",
```

**Objectif :** charger directement une chaîne HTML.

**Cas d’utilisation :**

- Fenêtres simples
- Contenu généré
- Tests

**Exemple :**

```go
HTML: `
<!DOCTYPE html>
<html>
<head><title>Simple Window</title></head>
<body><h1>Hello from Wails!</h1></body>
</html>
`,
```

### Ressources (au niveau de l’application uniquement)

La configuration des ressources n’est **pas** un champ de `WebviewWindowOptions`. L’application elle-même sert les ressources de l’interface via `application.Options.Assets` (`AssetOptions`) ; chaque fenêtre hérite de ce serveur de ressources.

```go
app := application.New(application.Options{
    Assets: application.AssetOptions{
        Handler: application.AssetFileServerFS(assets),
    },
})
```

**Pour plus de détails, consultez la section [Système de compilation](/concepts/build-system/).**

### UseApplicationMenu

**Type :** `bool` **Valeur par défaut :** `false` **Plateforme :** Windows, Linux (sans effet sous macOS)

```go
UseApplicationMenu: true,
```

**Objectif :** utiliser le menu de l’application (défini via `app.Menu.Set()`) pour cette fenêtre.

Sous **macOS**, cette option est sans effet, car macOS utilise toujours un menu d’application global en haut de l’écran.

Sous **Windows** et **Linux**, les fenêtres n’affichent aucun menu par défaut. Définir `UseApplicationMenu: true` indique à la fenêtre d’utiliser le menu au niveau de l’application, ce qui fournit une solution multiplateforme simple.

**Exemple :**

```go
// Set the application menu once
menu := app.NewMenu()
menu.AddRole(application.FileMenu)
menu.AddRole(application.EditMenu)
app.Menu.Set(menu)

// All windows with UseApplicationMenu will display this menu
app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:              "Main Window",
    UseApplicationMenu: true,
})

app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:              "Second Window",
    UseApplicationMenu: true,  // Also gets the app menu
})
```

**Remarques :**

- Si `UseApplicationMenu` et un menu propre à la fenêtre sont tous deux définis, le menu propre à la fenêtre est prioritaire
- Cela simplifie le code multiplateforme en supprimant la nécessité de vérifier le système d’exploitation à l’exécution
- Consultez [Menus d’application](/features/menus/application/) pour obtenir la documentation complète sur les menus

## Options d’entrée

### EnableFileDrop

**Type :** `bool` **Valeur par défaut :** `false` **Plateforme :** Toutes

```go
EnableFileDrop: true,
```

**Objectif :** autoriser le glisser-déposer de fichiers depuis le système d’exploitation vers la fenêtre.

Lorsque cette option est activée :

- Les fichiers glissés depuis un gestionnaire de fichiers peuvent être déposés dans votre application
- L’événement `WindowFilesDropped` est déclenché avec les chemins des fichiers déposés
- Les éléments dotés de l’attribut `data-file-drop-target` fournissent des informations détaillées sur le dépôt

**Cas d’utilisation :**

- Interfaces de téléversement de fichiers
- Éditeurs de documents
- Outils d’importation de médias
- Toute application acceptant des fichiers

**Exemple :**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:          "File Uploader",
    EnableFileDrop: true,
})

// Handle dropped files
window.OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {
    files := event.Context().DroppedFiles()
    details := event.Context().DropTargetDetails()
    
    for _, file := range files {
        fmt.Println("Dropped:", file)
    }
})
```

**Zones de dépôt HTML :**

```html
<!-- Mark elements as drop targets -->
<div id="upload" data-file-drop-target>
    Drop files here
</div>
```

**Consultez [Dépôt de fichiers](/features/drag-and-drop/files/) pour obtenir la documentation complète.**

## Options de sécurité

### ContentProtectionEnabled

**Type :** `bool` **Valeur par défaut :** `false` **Plateforme :** Windows (10+), macOS

```go
ContentProtectionEnabled: true,
```

**Objectif :** empêcher la capture du contenu de la fenêtre.

**Prise en charge des plateformes :**

- **Windows :** build 19041+ de Windows 10 (prise en charge complète), versions antérieures (prise en charge partielle)
- **macOS :** prise en charge complète
- **Linux :** non pris en charge

**Cas d’utilisation :**

- Applications bancaires
- Gestionnaires de mots de passe
- Dossiers médicaux
- Documents confidentiels

**Remarques importantes :**

1. N’empêche pas de photographier physiquement l’écran
2. Certains outils peuvent contourner la protection
3. Fait partie d’une stratégie de sécurité globale et ne constitue pas une protection suffisante à elle seule
4. Les fenêtres des outils de développement ne sont pas protégées automatiquement

**Exemple :**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Secure Window",
    ContentProtectionEnabled: true,
})

// Toggle at runtime
window.SetContentProtection(true)
```

### Permissions

**Type :** `map[PermissionType]Permission` **Valeur par défaut :** `nil` (gestion par défaut de la plateforme) **Plateforme :** Linux, Windows (macOS délègue la gestion à TCC)

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionAllow,
    application.PermissionCamera:     application.PermissionAllow,
},
```

**Objectif :** contrôler de manière déclarative, sans code propre à chaque plateforme, la gestion des demandes d’accès aux fonctionnalités (caméra, microphone, géolocalisation, notifications et lecture du presse-papiers) émises par le contenu web de la fenêtre.

**Valeurs de PermissionType :** `PermissionMicrophone`, `PermissionCamera`, `PermissionGeolocation`, `PermissionNotifications`, `PermissionClipboardRead`

**Valeurs du type Permission :**

- `PermissionDefault` (0) — gestion native de la plateforme : invite du système d’exploitation ou de WebView2 sous macOS et Windows ; sous Linux, l’accès à la caméra et au microphone est autorisé, tandis que tout le reste est refusé
- `PermissionAllow` (1) — accorder l’autorisation sans invite (sous Linux, seuls la caméra et le microphone sont implémentés ; les autres types restent refusés)
- `PermissionDeny` (2) — refuser sans invite

**Important — Windows :** avant l’existence de cette option, Wails accordait silencieusement l’accès à toutes les fonctionnalités de WebView2. Désormais, la définition de toute entrée dans `Permissions` désactive cette autorisation globale. Pour les fonctionnalités qui ne figurent pas dans la liste, l’invite native de WebView2 s’affiche au lieu d’accorder automatiquement l’autorisation. Répertoriez explicitement chaque fonctionnalité requise par votre application.

**Consultez [Autorisations](/features/windows/permissions/) pour obtenir le guide complet, la matrice des plateformes et des exemples.**

## Événements du cycle de vie des fenêtres

Les événements du cycle de vie des fenêtres sont gérés à l’aide de `OnWindowEvent` et de `RegisterHook`. Ces méthodes permettent de contrôler précisément le comportement de fermeture et de destruction des fenêtres.

### Annulation de la fermeture d’une fenêtre

Pour empêcher la fermeture d’une fenêtre, par exemple en présence de modifications non enregistrées, utilisez `RegisterHook` avec l’événement `WindowClosing` et appelez `event.Cancel()` :

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:  "main-window",
    Title: "My Application",
})

// Register a hook to intercept the closing event
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

**Points clés :**

- `RegisterHook` intercepte les événements avant qu’ils ne se produisent
- Appelez `event.Cancel()` pour empêcher la fermeture de la fenêtre
- La fenêtre reste ouverte après l’annulation

### Gestion de la fermeture d’une fenêtre

Pour effectuer le nettoyage à la fermeture d’une fenêtre, utilisez `OnWindowEvent` avec l’événement `WindowClosing` :

```go
window.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
    // Cleanup code runs here
    fmt.Printf("Window %s is closing\n", window.Name())

    // Close database connection
    if db != nil {
        db.Close()
    }

    // Remove from window list
    removeWindow(window.ID())
})
```

**Points clés :**

- `OnWindowEvent` gère les événements sur le point de se produire
- Le nettoyage s’exécute avant la destruction de la fenêtre
- La fermeture ne peut pas être annulée ici (utilisez `RegisterHook` à cette fin)

### Modèle de nettoyage d’une fenêtre singleton

Pour les fenêtres singleton (afin de garantir une seule instance), utilisez `WindowClosing` pour supprimer la référence :

```go
var settingsWindow *application.WebviewWindow

func ShowSettings(app *application.App) {
    // Create if doesn't exist
    if settingsWindow == nil {
        settingsWindow = app.Window.NewWithOptions(application.WebviewWindowOptions{
            Name:  "settings",
            Title: "Settings",
            Width: 600,
            Height: 400,
        })

        // Cleanup on close
        settingsWindow.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
            settingsWindow = nil
        })
    }

    // Show and focus
    settingsWindow.Show()
    settingsWindow.Focus()
}
```

## Options propres à chaque plateforme

### Options pour Mac

```go
Mac: application.MacWindow{
    TitleBar: application.MacTitleBar{
        AppearsTransparent: true,
        Hide:               false,
        HideTitle:          true,
        FullSizeContent:    true,
    },
    Backdrop:                application.MacBackdropTranslucent,
    InvisibleTitleBarHeight: 50,
    WindowClass:             application.MacWindowClassPanel,
    PanelPreferences: application.MacPanelPreferences{
        NonActivating:          true,
        FloatingPanel:          true,
        BecomesKeyOnlyIfNeeded: false,
        UtilityWindow:          false,
    },
    WindowLevel:             application.MacWindowLevelFloating,
    CollectionBehavior:      application.MacWindowCollectionBehaviorDefault,
    TabbingMode:             application.MacWindowTabbingModeDisallowed,
},
```

**TitleBar** (`MacTitleBar`)

- `AppearsTransparent` — Rend la barre de titre transparente et étend le contenu dans sa zone
- `Hide` — Masque entièrement la barre de titre
- `HideTitle` — Masque uniquement le texte du titre
- `FullSizeContent` — Étend le contenu à toute la fenêtre

Sous macOS, la transparence de la vue web et l’ouverture de l’inspecteur par programmation nécessitent le tag de compilation `private_mac_apis`. Sans ce tag, les mêmes options restent valides, mais les opérations réservées aux API privées sont sans effet. Le regroupement Liquid Glass est également ignoré et les styles utilisent des alternatives publiques. Consultez [API privées de macOS](/guides/build/private-macos-apis/) pour connaître les commandes de compilation et le comportement exact.

**Backdrop** (`MacBackdrop`)

- `MacBackdropNormal` — Arrière-plan opaque standard
- `MacBackdropTranslucent` — **Une API privée est requise pour rendre la vue web transparente.** Sans le tag, le flou natif reste derrière une vue web opaque.
- `MacBackdropTransparent` — **Une API privée est requise pour rendre la vue web transparente.** Sans le tag, la vue web reste opaque.
- `MacBackdropLiquidGlass` — **Une API privée est requise pour rendre la vue web transparente.** Sans le tag, la couche de verre reste derrière une vue web opaque ; les styles utilisent des alternatives publiques.

**LiquidGlass** (`MacLiquidGlass`)

| Champ ou valeur | Dépendance à une API privée sous macOS |
| --- | --- |
| `Style: LiquidGlassStyleAutomatic` | Le style natif standard est public ; la transparence de la vue web avec l’arrière-plan nécessite `private_mac_apis`. |
| `Style: LiquidGlassStyleLight` | Le tag conserve la correspondance existante avec le style natif transparent ; sans celui-ci, Wails utilise le verre standard avec l’apparence Aqua. |
| `Style: LiquidGlassStyleDark` | **API privée :** valeur de style natif non documentée `2` ; sans le tag, verre standard avec l’apparence Dark Aqua. |
| `Style: LiquidGlassStyleVibrant` | La correspondance avec le style natif transparent est publique ; la transparence de la vue web avec l’arrière-plan nécessite le tag. |
| `GroupID` | **API privée :** les valeurs non vides demandent un regroupement ; elles sont ignorées sans le tag. |
| `GroupSpacing` | **API privée :** les valeurs positives demandent un espacement de groupe ; elles sont ignorées sans le tag. |
| `Material`, `CornerRadius`, `TintColor` | Aucune dépendance propre à une API privée. |

Consultez [Valeurs de Liquid Glass](/guides/build/private-macos-apis/#liquid-glass-values) pour connaître les valeurs des styles natifs et leur disponibilité selon le système d’exploitation.

**InvisibleTitleBarHeight** (`int`)

- Hauteur de la zone invisible de la barre de titre (pour déplacer la fenêtre)
- Ne prend effet que lorsque la zone de déplacement de la barre de titre native est masquée, c’est-à-dire lorsque la fenêtre est sans cadre (`Frameless: true`) ou utilise une barre de titre transparente (`AppearsTransparent: true`)
- N’a aucun effet sur les fenêtres standard dont la barre de titre est visible

**WindowClass** (`MacWindowClass`)

- `MacWindowClassWindow` — Comportement standard de `NSWindow` (par défaut)
- `MacWindowClassPanel` — Un `NSPanel` auxiliaire qui ne devient jamais la fenêtre principale de l’application

`PanelPreferences` s’applique uniquement à `MacWindowClassPanel` :

- `NonActivating` ajoute `NSWindowStyleMaskNonactivatingPanel`. L’affichage du panneau ou le fait de lui donner le focus n’active pas l’application Wails, mais le panneau peut tout de même devenir la fenêtre clé pour les contrôles et la saisie de texte.
- `FloatingPanel` active le comportement de panneau flottant d’AppKit.
- `BecomesKeyOnlyIfNeeded` ne donne le statut de fenêtre clé que si la vue sur laquelle vous cliquez demande une saisie au clavier.
- `UtilityWindow` applique le style natif des fenêtres utilitaires.

Les panneaux Wails restent visibles lorsque l’application est désactivée et sont libérés à leur fermeture, conformément au cycle de vie attendu par `WebviewWindow`. Il s’agit de substitutions intentionnelles aux valeurs par défaut opposées de `NSPanel`.

La classe de fenêtre, le niveau, la politique d’activation et le comportement de collection répondent à des besoins distincts :

- `WindowClass` sélectionne `NSWindow` ou `NSPanel` et contrôle la sémantique des fenêtres principales et clés.
- `WindowLevel` contrôle l’ordre d’empilement. Utilisez `MacWindowLevelPopUpMenu` pour les superpositions de la barre des menus.
- `MacOptions.ActivationPolicy` contrôle l’application dans son ensemble, notamment sa présentation dans le Dock et la barre des menus. Un panneau qui ne déclenche pas l’activation ne nécessite pas de politique d’activation en tant qu’accessoire, même si une application peut tout de même l’utiliser pour masquer son icône dans le Dock.
- `CollectionBehavior` contrôle la participation aux Spaces et au mode plein écran.

```go
// Spotlight/menu-bar panel that leaves the current application active.
Mac: application.MacWindow{
    WindowClass: application.MacWindowClassPanel,
    PanelPreferences: application.MacPanelPreferences{
        NonActivating: true,
    },
    WindowLevel: application.MacWindowLevelPopUpMenu,
    CollectionBehavior: application.MacWindowCollectionBehaviorCanJoinAllSpaces |
        application.MacWindowCollectionBehaviorFullScreenAuxiliary |
        application.MacWindowCollectionBehaviorStationary,
},
```

**WindowLevel** (`MacWindowLevel`)

- `MacWindowLevelNormal` — Niveau de fenêtre standard (par défaut)
- `MacWindowLevelFloating` — Flotte au-dessus des fenêtres normales
- `MacWindowLevelTornOffMenu` — Niveau des menus détachés
- `MacWindowLevelModalPanel` — Niveau des panneaux modaux
- `MacWindowLevelMainMenu` — Niveau du menu principal
- `MacWindowLevelStatus` — Niveau des fenêtres d’état
- `MacWindowLevelPopUpMenu` — Niveau des menus locaux
- `MacWindowLevelScreenSaver` — Niveau de l’économiseur d’écran

Une valeur `WindowLevel` explicite prévaut sur `AlwaysOnTop` et `PanelPreferences.FloatingPanel`. En l’absence de niveau explicite, `AlwaysOnTop` ou un panneau flottant produit `MacWindowLevelFloating` ; sinon, le niveau est `MacWindowLevelNormal`. Un appel ultérieur à `SetAlwaysOnTop` reste une modification explicite à l’exécution.

**CollectionBehavior** (`MacWindowCollectionBehavior`)

Contrôle le comportement de la fenêtre dans les Spaces macOS et en plein écran. Ces valeurs de masque de bits peuvent être combinées avec un OU bit à bit (`|`).

**Comportement dans les Spaces :**

- `MacWindowCollectionBehaviorDefault` — Utilise FullScreenPrimary (par défaut, rétrocompatible)
- `MacWindowCollectionBehaviorCanJoinAllSpaces` — La fenêtre apparaît dans tous les Spaces
- `MacWindowCollectionBehaviorMoveToActiveSpace` — Se déplace vers le Space actif lors de son affichage
- `MacWindowCollectionBehaviorManaged` — Comportement par défaut d’une fenêtre gérée
- `MacWindowCollectionBehaviorTransient` — Fenêtre temporaire ou transitoire
- `MacWindowCollectionBehaviorStationary` — Reste immobile lors des changements de Space

**Parcours des fenêtres :**

- `MacWindowCollectionBehaviorParticipatesInCycle` — Incluse dans le parcours avec Cmd+`
- `MacWindowCollectionBehaviorIgnoresCycle` — Exclue du parcours avec Cmd+`

**Comportement en plein écran :**

- `MacWindowCollectionBehaviorFullScreenPrimary` — Peut passer en mode plein écran
- `MacWindowCollectionBehaviorFullScreenAuxiliary` — Peut se superposer aux applications en plein écran
- `MacWindowCollectionBehaviorFullScreenNone` — Désactive la possibilité de passer en plein écran
- `MacWindowCollectionBehaviorFullScreenAllowsTiling` — Autorise la disposition côte à côte (macOS 10.11 ou version ultérieure)
- `MacWindowCollectionBehaviorFullScreenDisallowsTiling` — Empêche la disposition en mosaïque (macOS 10.11 ou version ultérieure)

**Exemple — Fenêtre de type Spotlight :**

```go
// Window that appears on all Spaces AND can overlay fullscreen apps
Mac: application.MacWindow{
	WindowClass: application.MacWindowClassPanel,
	PanelPreferences: application.MacPanelPreferences{
		NonActivating: true,
	},
    CollectionBehavior: application.MacWindowCollectionBehaviorCanJoinAllSpaces |
                        application.MacWindowCollectionBehaviorFullScreenAuxiliary,
    WindowLevel:        application.MacWindowLevelFloating,
},
```

**Exemple — Comportement unique :**

```go
// Window that can appear over fullscreen applications
Mac: application.MacWindow{
    CollectionBehavior: application.MacWindowCollectionBehaviorFullScreenAuxiliary,
},
```

**TabbingMode** (`MacWindowTabbingMode`)

Contrôle le comportement de regroupement des fenêtres en onglets sous macOS 10.12 et les versions ultérieures. Ce regroupement permet de réunir plusieurs fenêtres sous forme d’onglets.

**Options :**

- `MacWindowTabbingModeDefault` — Valeur sentinelle nulle (non définie explicitement). À l’exécution, le regroupement en onglets est interdit par défaut
- `MacWindowTabbingModeAutomatic` — Le système détermine le comportement de regroupement en onglets
- `MacWindowTabbingModePreferred` — La fenêtre privilégie le mode avec onglets
- `MacWindowTabbingModeDisallowed` — Désactive le regroupement des fenêtres en onglets

**Exemple — Désactiver le regroupement des fenêtres en onglets :**

```go
Mac: application.MacWindow{
    TabbingMode: application.MacWindowTabbingModeDisallowed,
},
```

**Exemple — Privilégier le regroupement des fenêtres en onglets :**

```go
Mac: application.MacWindow{
    TabbingMode: application.MacWindowTabbingModePreferred,
},
```

**WebviewPreferences** (`MacWebviewPreferences`)

Permet de contrôler précisément la configuration `WKWebView` sous-jacente. Tous les champs sont facultatifs : les champs non définis ne modifient pas la valeur par défaut de WebKit.

```go
Mac: application.MacWindow{
    WebviewPreferences: application.MacWebviewPreferences{
        TabFocusesLinks:                       optional.True,
        TextInteractionEnabled:                optional.True,
        FullscreenEnabled:                     optional.False,
        AllowsBackForwardNavigationGestures:   optional.False,
        AllowsMagnification:                   optional.True,
        AllowsAirPlayForMediaPlayback:         optional.True,
        JavaScriptCanOpenWindowsAutomatically: optional.False,
        MinimumFontSize:                       optional.NewVar(12.0),
        ApplicationNameForUserAgent:           "MyApp",
        EnableAutoplayWithoutUserAction:       optional.True,
    },
},
```

- `TabFocusesLinks` — Lorsque `true`, la touche Tab déplace le focus vers les liens et les contrôles de formulaire (valeur par défaut : `false`)
- `TextInteractionEnabled` — Lorsque `true`, les utilisateurs peuvent sélectionner le texte dans la vue web et interagir avec lui (valeur par défaut : `true`)
- `FullscreenEnabled` — Lorsque `true`, le contenu web peut passer en plein écran via l’API Fullscreen HTML (valeur par défaut : `false`). Nécessite macOS 12.3 ou version ultérieure.
- `AllowsBackForwardNavigationGestures` — Lorsque `true`, les gestes de balayage horizontal déclenchent la navigation vers la page précédente ou suivante (valeur par défaut : `false`)
- `AllowsMagnification` — Lorsque `true`, le zoom par pincement est activé dans la vue web (valeur par défaut : `false`)
- `AllowsAirPlayForMediaPlayback` — Lorsque `true`, les contenus multimédias peuvent être diffusés vers des appareils AirPlay (valeur par défaut : `true`)
- `JavaScriptCanOpenWindowsAutomatically` — Lorsque `true`, JavaScript peut ouvrir de nouvelles fenêtres sans action de l’utilisateur (valeur par défaut : `false`)
- `MinimumFontSize` — Taille minimale de la police, en points. Utilisez `optional.NewVar(12.0)` pour la définir. Si elle n’est pas définie, la valeur par défaut de WebKit est conservée.
- `ApplicationNameForUserAgent` — Remplace le suffixe du nom de l’application dans la chaîne d’agent utilisateur de WebKit. Utile lorsque des sites rejettent l’identifiant `"wails.io"` par défaut (par exemple, les contenus YouTube intégrés). Laissez ce champ vide pour conserver la valeur par défaut.
- `EnableAutoplayWithoutUserAction` — Lorsque `true`, la lecture automatique des contenus audio et vidéo est possible sans action de l’utilisateur. Correspond à `WKWebViewConfiguration.mediaTypesRequiringUserActionForPlayback = WKAudiovisualMediaTypeNone` (valeur par défaut : `false`)

### Options Windows (par fenêtre)

La structure propre à chaque fenêtre est `application.WindowsWindow` — et **non** `WindowsOptions` (qui est la structure au niveau de l’*application*).

```go
Windows: application.WindowsWindow{
    DisableIcon:                       false,
    DisableMenu:                       false,
    BackdropType:                      application.Auto,
    CustomTheme:                       application.ThemeSettings{},
    DisableFramelessWindowDecorations: false,
    NonClientRegionSupport:            false,
    WebView2CompositionHosting:        false,
},
```

**DisableIcon** (`bool`)

- Supprime l’icône de la barre de titre.

**DisableMenu** (`bool`)

- Désactive la barre de menus de la fenêtre. Lorsque `true`, la fenêtre n’affiche aucune barre de menus, même si une barre a été configurée.
- Valeur par défaut : `false`

**BackdropType** (`BackdropType`)

- `application.Auto` — Valeur par défaut du système
- `application.None` — Aucun arrière-plan
- `application.Mica` — Matériau Mica (Windows 11)
- `application.Acrylic` — Matériau Acrylic (Windows 11)
- `application.Tabbed` — Matériau à onglets (Windows 11)

Il n’existe aucune constante de type `WindowsBackdropTypeMica` : utilisez `application.Mica`, etc.

**CustomTheme** (`ThemeSettings`)

- Valeur (et non pointeur). Couleurs personnalisées des modes sombre et clair pour la bordure de la fenêtre, le texte et l’arrière-plan de la barre de titre, ainsi que la barre de menus.

**DisableFramelessWindowDecorations** (`bool`)

- Désactive les décorations par défaut des fenêtres sans cadre (ombre Aero et angles arrondis).

**NonClientRegionSupport** (`bool`)

- Active la prise en charge native par WebView2 de `app-region: drag` / `app-region: no-drag` pour les barres de titre personnalisées sans cadre.
- Cette option permet uniquement le déplacement natif simple de l’application. Elle ne fournit ni le comportement natif des boutons personnalisés de la barre de titre, ni les fonctions Windows 11 Snap Assist / Snap Layouts pour les boutons d’agrandissement personnalisés.

**WebView2CompositionHosting** (`bool`)

- Active la prise en charge, gérée par Wails, de `--wails-non-client-region` pour les boutons personnalisés de la barre de titre avec le comportement Windows natif, notamment les fonctions Windows 11 Snap Assist / Snap Layouts sur les boutons d’agrandissement personnalisés.
- Expérimental. Cette option héberge WebView2 au moyen de `ICoreWebView2CompositionController` et de DirectComposition, au lieu du contrôleur par défaut hébergé dans un HWND.
- Peut être combiné avec `NonClientRegionSupport` lorsqu’une fenêtre nécessite à la fois la prise en charge native par WebView2 de `app-region` et des régions de boutons personnalisés de la barre de titre gérées par Wails.

**Exemple :**

```go
Windows: application.WindowsWindow{
    BackdropType: application.Mica,
    DisableIcon:  true,
},
```

**Exemple — Régions de barre de titre Windows personnalisée :**

```go
Windows: application.WindowsWindow{
    NonClientRegionSupport:    true,
    WebView2CompositionHosting: true,
},
```

Consultez [Fenêtres sans cadre](/features/windows/frameless/#native-non-client-regions-on-windows) pour connaître en détail le comportement, les compromis et le CSS correspondant.

### Options Linux (par fenêtre)

La structure propre à chaque fenêtre est `application.LinuxWindow` — et **non** `LinuxOptions`.

```go
Linux: application.LinuxWindow{
    Icon:                []byte{/* PNG data */},
    WindowIsTranslucent: false,
},
```

**Icon** (`[]byte`)

- Icône de la fenêtre (format PNG).

**WindowIsTranslucent** (`bool`)

- Nécessite la prise en charge du compositeur.

**Exemple :**

```go
//go:embed icon.png
var icon []byte

Linux: application.LinuxWindow{
    Icon: icon,
},
```

## Options Windows au niveau de l’application

Certaines options propres à Windows doivent être configurées au niveau de l’application plutôt que pour chaque fenêtre. En effet, WebView2 partage un environnement de navigateur unique pour chaque chemin de données utilisateur.

### Indicateurs du navigateur

Les indicateurs du navigateur WebView2 contrôlent les fonctionnalités expérimentales et le comportement de **toutes les fenêtres** de votre application. Ils doivent être définis dans `application.Options.Windows` :

```go
app := application.New(application.Options{
    Name: "My App",
    Windows: application.WindowsOptions{
        // Enable experimental WebView2 features
        EnabledFeatures: []string{
            "msWebView2EnableDraggableRegions",
        },

        // Disable specific features
        DisabledFeatures: []string{
            "msSmartScreenProtection",  // Always disabled by Wails
        },

        // Additional Chromium command-line arguments
        AdditionalBrowserArgs: []string{
            "--disable-gpu",
            "--remote-debugging-port=9222",
        },
    },
})
```

**EnabledFeatures** (`[]string`)

- Liste des indicateurs de fonctionnalité WebView2 à activer
- Consultez les [indicateurs du navigateur WebView2](https://learn.microsoft.com/en-us/microsoft-edge/webview2/concepts/webview-features-flags) pour connaître les indicateurs disponibles
- Exemple : `"msWebView2EnableDraggableRegions"`

**DisabledFeatures** (`[]string`)

- Liste des indicateurs de fonctionnalité WebView2 à désactiver
- Wails désactive automatiquement `msSmartScreenProtection`
- Exemple : `"msExperimentalFeature"`

**AdditionalBrowserArgs** (`[]string`)

- Arguments de ligne de commande Chromium transmis au processus du navigateur
- Doit inclure le préfixe `--` (par exemple, `"--remote-debugging-port=9222"`)
- Consultez les [options de ligne de commande Chromium](https://peter.sh/experiments/chromium-command-line-switches/) pour connaître les arguments disponibles

@note{type="caution" title="Important"}
Ces indicateurs s’appliquent globalement à TOUTES les fenêtres, car WebView2 partage un environnement de navigateur unique pour chaque chemin de données utilisateur. Vous ne pouvez pas utiliser des indicateurs de navigateur différents selon les fenêtres.

@end

**Exemple complet :**

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "My App",
        Windows: application.WindowsOptions{
            // Enable draggable regions feature
            EnabledFeatures: []string{
                "msWebView2EnableDraggableRegions",
            },
            // Enable remote debugging
            AdditionalBrowserArgs: []string{
                "--remote-debugging-port=9222",
            },
        },
    })

    // All windows will use the browser flags configured above
    window := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:  "Main Window",
        Width:  1024,
        Height: 768,
    })

    window.Show()
    app.Run()
}
```

## Exemple complet

Voici une configuration de fenêtre prête pour la production :

```go
package main

import (
    _ "embed"
    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var assets embed.FS

//go:embed icon.png
var icon []byte

func main() {
    app := application.New(application.Options{
        Name: "My Application",
    })

    window := app.Window.NewWithOptions(application.WebviewWindowOptions{
        // Identity
        Name:  "main-window",
        Title: "My Application",

        // Size and Position
        Width:     1200,
        Height:    800,
        MinWidth:  800,
        MinHeight: 600,

        // Initial State
        StartState: application.WindowStateNormal,

        // Appearance
        BackgroundColour: application.NewRGB(255, 255, 255),

        // Platform-Specific (per-window structs)
        Mac: application.MacWindow{
            TitleBar: application.MacTitleBar{
                AppearsTransparent: true,
            },
            Backdrop: application.MacBackdropTranslucent,
        },

        Windows: application.WindowsWindow{
            BackdropType: application.Mica,
            DisableIcon:  false,
        },

        Linux: application.LinuxWindow{
            Icon: icon,
        },
    })

    window.Center()
    window.Show()

    app.Run()
}
```

Les ressources du frontend sont servies au niveau de l’application (`application.Options.Assets`), et non pour chaque fenêtre.

## Étapes suivantes

- [Principes de base des fenêtres](/features/windows/basics/) – Créer et contrôler des fenêtres
- [Fenêtres multiples](/features/windows/multiple/) – Modèles pour applications multifenêtres
- [Fenêtres sans cadre](/features/windows/frameless/) – Habillage de fenêtre personnalisé
- [Événements de fenêtre](/features/windows/events/) – Événements du cycle de vie

---

**Des questions ?** Posez-les sur [Discord](https://discord.gg/JDdSxwjhGf) ou consultez les [exemples](https://github.com/wailsapp/wails/tree/master/v3/examples).
