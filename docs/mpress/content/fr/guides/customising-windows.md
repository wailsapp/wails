---
title: "Personnalisation des fenêtres dans Wails"
description: "Personnalisez l’apparence et le comportement des fenêtres de vos applications Wails"
slug: "guides/customising-windows"
sourcePath: "guides/customising-windows.md"
---

Plateformes concernées : <span class="mpress-badge mpress-badge-note">Windows</span> <span class="mpress-badge mpress-badge-success">macOS</span>

<br/>

Wails fournit une API permettant de contrôler l’apparence et le fonctionnement des commandes d’une fenêtre. Cette fonctionnalité est disponible sous Windows et macOS, mais pas sous Linux.

## Définition de l’état des boutons de la fenêtre

L’état des boutons est défini par l’énumération `ButtonState` :

```go
type ButtonState int

const (
    ButtonEnabled   ButtonState = 0
    ButtonDisabled  ButtonState = 1
    ButtonHidden    ButtonState = 2
)
```

- `ButtonEnabled` : le bouton est activé et visible.
- `ButtonDisabled` : le bouton est visible, mais désactivé (grisé).
- `ButtonHidden` : le bouton est masqué dans la barre de titre.

Vous pouvez définir l’état des boutons lors de la création de la fenêtre ou pendant l’exécution.

### Définition de l’état des boutons lors de la création de la fenêtre

Lors de la création d’une fenêtre, vous pouvez définir l’état initial des boutons à l’aide de la structure `WebviewWindowOptions` :

```go {title="main.go"}
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name:        "My Application",
    })

    app.Window.NewWithOptions(application.WebviewWindowOptions{
        MinimiseButtonState:   application.ButtonHidden,
        MaximiseButtonState:   application.ButtonDisabled,
        CloseButtonState:      application.ButtonEnabled,
        FullscreenButtonState: application.ButtonEnabled,
    })

    app.Run()
}
```

Dans l’exemple ci-dessus, le bouton de réduction est masqué, le bouton d’agrandissement est désactivé (grisé) et le bouton de fermeture est activé.

### Définition de l’état des boutons pendant l’exécution

Vous pouvez également modifier l’état des boutons pendant l’exécution à l’aide des méthodes suivantes de l’interface `Window` :

```go
window.SetMinimiseButtonState(wails.ButtonHidden)
window.SetMaximiseButtonState(wails.ButtonEnabled)
window.SetCloseButtonState(wails.ButtonDisabled)
window.SetFullscreenButtonState(wails.ButtonEnabled)
```

### macOS : MaximiseButtonState et FullscreenButtonState partagent le même bouton

Sous macOS, le bouton vert de type « feu de signalisation » (`NSWindowZoomButton`) constitue la même commande physique pour l’agrandissement et le plein écran. Sans mesure particulière, attribuer des valeurs différentes à `MaximiseButtonState` et `FullscreenButtonState` lors de la création de la fenêtre entraînerait un remplacement silencieux selon le principe de la dernière écriture prioritaire.

Pour éviter cela, Wails applique, lors de l’initialisation, l’état **le plus restrictif** des deux, selon l’ordre `ButtonEnabled` < `ButtonDisabled` < `ButtonHidden`.

| `MaximiseButtonState` | `FullscreenButtonState` | État effectif sous macOS |
| --- | --- | --- |
| `ButtonEnabled` | `ButtonEnabled` | `ButtonEnabled` |
| `ButtonDisabled` | `ButtonEnabled` | `ButtonDisabled` |
| `ButtonEnabled` | `ButtonHidden` | `ButtonHidden` |
| `ButtonDisabled` | `ButtonHidden` | `ButtonHidden` |

Pendant l’exécution, `SetMaximiseButtonState` et `SetFullscreenButtonState` ciblent tous deux `NSWindowZoomButton` sous macOS ; le dernier appel l’emporte donc.

### Différences entre les plateformes

La gestion de l’état des boutons se comporte légèrement différemment sous Windows et macOS :

|  | Windows | Mac |
| --- | --- | --- |
| Désactiver Réduire/Agrandir/Fermer | Désactive Réduire/Agrandir/Fermer | Désactive Réduire/Agrandir/Fermer |
| Masquer Réduire | Désactive Réduire | Masque le bouton Réduire |
| Masquer Agrandir | Désactive Agrandir | Masque le bouton Agrandir |
| Masquer Fermer | Masque toutes les commandes | Masque Fermer |
| `FullscreenButtonState` | Sans effet | Cible le bouton de zoom (vert) |

Remarque : sous Windows, il n’est pas possible de masquer séparément les boutons Réduire et Agrandir. Toutefois, leur désactivation simultanée masque ces deux contrôles et n’affiche que le bouton Fermer. La barre de titre standard de Windows ne comporte aucun bouton dédié au plein écran ; `FullscreenButtonState` n’y a donc aucun effet.

### Contrôle du style de la fenêtre (Windows)

Pour contrôler le style de la barre de titre sous Windows, utilisez le champ `ExStyle` de la structure `WebviewWindowOptions` :

Exemple :

```go {title="main.go"}
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/w32"
)

func main() {
    app := application.New(application.Options{
        Name: "My Application",
    })

    app.Window.NewWithOptions(application.WebviewWindowOptions{
        Windows: application.WindowsWindow{
            ExStyle: w32.WS_EX_TOOLWINDOW | w32.WS_EX_NOREDIRECTIONBITMAP | w32.WS_EX_TOPMOST,
        },
    })

    app.Run()
}
```

Ce paramètre remplace les autres options qui affectent le style étendu d’une fenêtre :

- HiddenOnTaskbar
- AlwaysOnTop
- IgnoreMouseEvents
- BackgroundType
