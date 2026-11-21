---
title: "Fenêtres sans cadre"
description: "Créez un habillage de fenêtre personnalisé avec des fenêtres sans cadre"
slug: "features/windows/frameless"
sourcePath: "features/windows/frameless.md"
---

## Fenêtres sans cadre

Wails prend en charge les **fenêtres sans cadre** avec des zones de déplacement définies en CSS et un comportement natif propre à chaque plateforme. Supprimez la barre de titre native de la plateforme pour maîtriser entièrement l’habillage de la fenêtre, créer des designs personnalisés et proposer des expériences utilisateur uniques, tout en conservant les fonctionnalités essentielles telles que le déplacement, le redimensionnement et les contrôles système.

![L’application de démarrage TypeScript v3 par défaut de Wails, exécutée dans une fenêtre sans cadre avec les angles natifs de macOS](/assets/screenshots/frameless-v3-native-corners-macos.png)

L’exemple ci-dessus est l’application de démarrage TypeScript v3 par défaut de Wails, avec `Frameless: true` activé.

## Démarrage rapide

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:     "Frameless App",
    Width:     800,
    Height:    600,
    Frameless: true,
})
```

**CSS pour une barre de titre déplaçable :**

```css
.titlebar {
    --wails-draggable: drag;
    height: 40px;
    background: #333;
}

.titlebar button {
    --wails-draggable: no-drag;
}
```

**HTML :**

```html
<div class="titlebar">
    <span>My Application</span>
    <button onclick="window.close()">×</button>
</div>
```

**C’est tout !** Vous disposez maintenant d’une barre de titre personnalisée.

## Création de fenêtres sans cadre

### Rayon des angles (macOS)

Par défaut, les fenêtres sans cadre conservent les angles arrondis standard d’AppKit sous macOS. Définissez `Mac.CornerRadius` pour utiliser un rayon personnalisé (en points) :

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Mac: application.MacWindow{
        CornerRadius: 16,
    },
})
```

Définissez `Mac.CornerType` sur `MacWindowCornerTypeSquare` pour obtenir des angles droits. Cela ignore `CornerRadius` :

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Mac: application.MacWindow{
        CornerType: application.MacWindowCornerTypeSquare,
    },
})
```

### Fenêtre sans cadre de base

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Width:     800,
    Height:    600,
})
```

**Ce que vous obtenez :**

- Aucune barre de titre
- Aucune bordure de fenêtre
- Aucun bouton système
- Arrière-plan transparent (facultatif)

**Ce que vous devez implémenter :**

- Zone déplaçable
- Boutons de fermeture, de réduction et d’agrandissement
- Poignées de redimensionnement (si la fenêtre est redimensionnable)

### Avec un arrière-plan transparent

**API privée sous macOS :** définissez `Mac.Backdrop: application.MacBackdropTransparent` et effectuez la compilation avec `-tags private_mac_apis` pour rendre la webview transparente. Sans cette balise, la webview native reste opaque, même si l’arrière-plan HTML/CSS est transparent. `Frameless` et `TitleBar.AppearsTransparent` utilisent eux-mêmes des API publiques. Consultez [API privées de macOS](/guides/build/private-macos-apis/#webview-transparency-and-background).

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless:      true,
    BackgroundType: application.BackgroundTypeTransparent,
    Mac: application.MacWindow{
        Backdrop: application.MacBackdropTransparent, // Requires -tags private_mac_apis.
    },
})
```

**Cas d’utilisation :**

- Angles arrondis
- Formes personnalisées
- Fenêtres en superposition
- Écrans de démarrage

## Zones de déplacement

### Déplacement basé sur CSS

Utilisez la propriété CSS `--wails-draggable` :

```css
/* Draggable area */
.titlebar {
    --wails-draggable: drag;
}

/* Non-draggable elements within draggable area */
.titlebar button {
    --wails-draggable: no-drag;
}
```

**Valeurs :**

- `drag` — La zone est déplaçable
- `no-drag` — La zone n’est pas déplaçable (même si son parent l’est)

### Exemple complet de barre de titre

```html
<div class="titlebar">
    <div class="title">My Application</div>
    <div class="controls">
        <button class="minimize">−</button>
        <button class="maximize">□</button>
        <button class="close">×</button>
    </div>
</div>
```

```css
.titlebar {
    --wails-draggable: drag;
    display: flex;
    justify-content: space-between;
    align-items: center;
    height: 40px;
    background: #2c2c2c;
    color: white;
    padding: 0 16px;
}

.title {
    font-size: 14px;
    user-select: none;
}

.controls {
    display: flex;
    gap: 8px;
}

.controls button {
    --wails-draggable: no-drag;
    width: 32px;
    height: 32px;
    border: none;
    background: transparent;
    color: white;
    font-size: 16px;
    cursor: pointer;
    border-radius: 4px;
}

.controls button:hover {
    background: rgba(255, 255, 255, 0.1);
}

.controls .close:hover {
    background: #e81123;
}
```

**JavaScript pour les boutons :**

```javascript
import { Window } from '@wailsio/runtime'

document.querySelector('.minimize').addEventListener('click', () => Window.Minimise())
document.querySelector('.maximize').addEventListener('click', () => Window.Maximise())
document.querySelector('.close').addEventListener('click', () => Window.Close())
```

## Zones non clientes natives sous Windows

Windows peut traiter certaines parties d’une barre de titre personnalisée comme des zones non clientes natives. Vous pouvez ainsi dessiner la barre de titre et les boutons de légende avec n’importe quelle conception HTML/CSS, tout en conservant le comportement natif de Windows : la zone de légende permet de déplacer la fenêtre, le bouton d’agrandissement peut afficher les fonctionnalités 11 Snap Assist / Snap Layouts de Windows, et les boutons de réduction, d’agrandissement et de fermeture bénéficient de la détection de zone et de l’état de souris natifs.

La vidéo ci-dessous présente une barre de titre HTML/CSS personnalisée qui utilise la détection de zone native de Windows, notamment les fonctionnalités 11 Snap Assist / Snap Layouts de Windows sur un bouton d’agrandissement personnalisé.

<video src="/assets/windows-native-non-client-regions/wails-app-region.mp4" controls muted playsInline></video>

Wails prend en charge deux mécanismes propres à Windows :

- `app-region` via la prise en charge native des zones non clientes de WebView2
- `--wails-non-client-region` via le suivi assuré par le runtime Wails pour les boutons de légende personnalisés

### Choix d’un mode

@note{type="caution" title="Expérimental"}
`WebView2CompositionHosting` modifie en interne la manière dont la fenêtre héberge WebView2 et interagit avec celui-ci. Au lieu du contrôleur WebView2 par défaut hébergé par un HWND, Wails utilise l’hébergement par contrôleur de composition et transmet explicitement les entrées. Ce mode peut présenter des problèmes de rendu, d’entrée, de focus ou de compatibilité avec le runtime WebView2. Activez-le uniquement si vous avez besoin du comportement natif de boutons de légende personnalisés, et testez soigneusement votre application avec les versions de Windows et du runtime WebView2 que vous prenez en charge.

@end

Choisissez en fonction du comportement Windows dont vous avez besoin :

- Utilisez `NonClientRegionSupport` pour permettre simplement le déplacement natif de l’application avec les zones `app-region: drag` et `app-region: no-drag` de WebView2.
- Utilisez `WebView2CompositionHosting` lorsque vos boutons personnalisés de réduction, d’agrandissement et de fermeture doivent se comporter comme les boutons de légende natifs de Windows.
- Activez les deux lorsque la même fenêtre doit prendre en charge les zones `app-region` natives de WebView2 et les zones de boutons de légende personnalisés gérées par Wails.

`NonClientRegionSupport` est l’alternative native et légère au suivi `--wails-draggable` de Wails. Vous marquez les zones déplaçables et non déplaçables avec CSS, WebView2 détermine quels pixels appartiennent à la légende, puis Wails demande à WebView2 la zone native lors de la détection de zone.

À ce jour, la portée de ce mode se limite à cela. Il ne permet pas aux boutons personnalisés de réduction, d’agrandissement ou de fermeture de se comporter comme les boutons de légende natifs de Windows, et n’active pas les fonctionnalités 11 Snap Assist / Snap Layouts de Windows pour un bouton d’agrandissement personnalisé. Utilisez-le lorsque vous avez uniquement besoin du déplacement natif de l’application, sans le mécanisme supplémentaire de `--wails-draggable`.

`WebView2CompositionHosting` est destiné aux boutons de légende personnalisés dotés d’un comportement natif. Wails suit les rectangles du DOM marqués avec `--wails-non-client-region`, les associe à des valeurs de détection de zone Windows telles que `HTMINBUTTON`, `HTMAXBUTTON` et `HTCLOSE`, puis retransmet les entrées de la souris à la surface WebView2 hébergée par composition. C’est ce qui permet à un bouton d’agrandissement personnalisé de prendre part aux fonctionnalités 11 Snap Assist / Snap Layouts de Windows, tout en conservant le design visuel de votre choix.

Autrement dit, `NonClientRegionSupport` correspond à la prise en charge des zones CSS natives de WebView2. Avec `WebView2CompositionHosting`, Wails prend en charge la composition détenue par l’hôte et la détection personnalisée des zones non clientes.

### app-region de WebView2

Activez la prise en charge native des zones non clientes de WebView2 pour la fenêtre :

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Windows: application.WindowsWindow{
        NonClientRegionSupport: true,
    },
})
```

Marquez ensuite les zones déplaçables avec la propriété CSS `app-region` :

```css
.titlebar {
    app-region: drag;
}

.titlebar button,
.titlebar input,
.titlebar select,
.titlebar textarea {
    app-region: no-drag;
}
```

Utilisez ce mode si vous avez uniquement besoin du déplacement natif par la légende et si les contrôles de votre barre de titre sont gérés par des clics frontend ordinaires.

Ce mode est limité par la prise en charge des zones non clientes propre à WebView2. Dans les versions actuelles de WebView2, seules les zones déplaçables et non déplaçables sont donc disponibles. Ce mode n’est pas conçu pour représenter des boutons de légende frontend entièrement personnalisés, avec des rôles natifs distincts de réduction, d’agrandissement et de fermeture.

### Boutons de légende personnalisés avec comportement natif

Pour que les boutons personnalisés de réduction, d’agrandissement et de fermeture se comportent comme les boutons de légende système, activez l’hébergement par composition :

@note{type="caution" title="Expérimental"}
`WebView2CompositionHosting` utilise l’hébergement du contrôleur de composition WebView2 avec DirectComposition. Consultez [Choix d’un mode](#choix-dun-mode) avant de l’activer.

@end

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Windows: application.WindowsWindow{
        WebView2CompositionHosting: true,
    },
})
```

Marquez ensuite chaque zone frontend avec `--wails-non-client-region` :

```html
<div class="titlebar">
    <div class="title">My Application</div>
    <div class="window-controls">
        <button class="window-button minimize" aria-label="Minimize"></button>
        <button class="window-button maximize" aria-label="Maximize"></button>
        <button class="window-button close" aria-label="Close"></button>
    </div>
</div>
```

```css
.titlebar {
    --wails-non-client-region: caption;
    height: 40px;
}

.window-controls {
    display: flex;
    height: 100%;
}

.window-button {
    width: 46px;
    border: 0;
    background: transparent;
}

.window-button.minimize {
    --wails-non-client-region: minimize;
}

.window-button.maximize {
    --wails-non-client-region: maximize;
}

.window-button.close {
    --wails-non-client-region: close;
}
```

Valeurs `--wails-non-client-region` prises en charge :

- `caption` — zone de barre de titre permettant le déplacement
- `minimize` — cible de clic du bouton natif de réduction
- `maximize` — cible de clic du bouton natif d’agrandissement, y compris le comportement au survol de Snap Assist / Snap Layouts de Windows 11
- `close` — cible de clic du bouton natif de fermeture

Le runtime Wails surveille les modifications du DOM, des styles, de la taille, du défilement et de la zone d’affichage, puis envoie des instantanés des zones à la fenêtre native. La géométrie des zones est mesurée en pixels CSS, puis convertie en pixels physiques pour les tests de positionnement de Windows.

La conception visuelle reste entièrement sous votre contrôle. Les zones indiquent uniquement à Windows la signification de chaque rectangle ; la forme, l’icône, la couleur, l’espacement, le style au survol et la disposition des boutons proviennent toujours de votre frontend.

### Combiner les deux

Vous pouvez activer les deux options si vous souhaitez utiliser dans la même fenêtre la prise en charge de `app-region` par WebView2 et les zones de boutons de barre de titre gérées par Wails :

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Windows: application.WindowsWindow{
        NonClientRegionSupport:    true,
        WebView2CompositionHosting: true,
    },
})
```

## Boutons système

### Implémenter la fermeture, la réduction et l’agrandissement

**Côté Go :**

```go
type WindowControls struct {
    window *application.WebviewWindow
}

func (wc *WindowControls) Minimise() {
    wc.window.Minimise()
}

func (wc *WindowControls) Maximise() {
    if wc.window.IsMaximised() {
        wc.window.UnMaximise()
    } else {
        wc.window.Maximise()
    }
}

func (wc *WindowControls) Close() {
    wc.window.Close()
}
```

**Côté JavaScript :**

```javascript
import { Minimise, Maximise, Close } from './bindings/WindowControls'

document.querySelector('.minimize').addEventListener('click', Minimise)
document.querySelector('.maximize').addEventListener('click', Maximise)
document.querySelector('.close').addEventListener('click', Close)
```

**Vous pouvez aussi utiliser les méthodes du runtime :**

```javascript
import { Window } from '@wailsio/runtime'

document.querySelector('.minimize').addEventListener('click', () => Window.Minimise())
document.querySelector('.maximize').addEventListener('click', () => Window.Maximise())
document.querySelector('.close').addEventListener('click', () => Window.Close())
```

### Basculer l’état d’agrandissement

Suivez l’état d’agrandissement pour adapter l’icône du bouton :

```javascript
import { Window } from '@wailsio/runtime'

async function toggleMaximise() {
    const isMaximised = await Window.IsMaximised()

    if (isMaximised) {
        await Window.Restore()
    } else {
        await Window.Maximise()
    }

    updateMaximiseButton()
}

async function updateMaximiseButton() {
    const isMaximised = await Window.IsMaximised()
    const button = document.querySelector('.maximize')
    button.textContent = isMaximised ? '❐' : '□'
}
```

## Poignées de redimensionnement

### Redimensionnement avec CSS

Wails fournit des poignées de redimensionnement automatiques pour les fenêtres sans cadre :

```css
/* Enable resize on all edges */
body {
    --wails-resize: all;
}

/* Or specific edges */
.resize-top {
    --wails-resize: top;
}

.resize-bottom {
    --wails-resize: bottom;
}

.resize-left {
    --wails-resize: left;
}

.resize-right {
    --wails-resize: right;
}

/* Corners */
.resize-top-left {
    --wails-resize: top-left;
}

.resize-top-right {
    --wails-resize: top-right;
}

.resize-bottom-left {
    --wails-resize: bottom-left;
}

.resize-bottom-right {
    --wails-resize: bottom-right;
}
```

**Valeurs :**

- `all` — redimensionnement depuis tous les bords
- `top`, `bottom`, `left`, `right` — bords spécifiques
- `top-left`, `top-right`, `bottom-left`, `bottom-right` — coins
- `none` — aucun redimensionnement

### Exemple de poignées de redimensionnement

```html
<div class="window">
    <div class="titlebar">...</div>
    <div class="content">...</div>
    <div class="resize-handle resize-bottom-right"></div>
</div>
```

```css
.resize-handle {
    position: absolute;
    width: 16px;
    height: 16px;
}

.resize-bottom-right {
    --wails-resize: bottom-right;
    bottom: 0;
    right: 0;
    cursor: nwse-resize;
}
```

## Comportement propre à chaque plateforme

@tabs{sync-key="platform"}
[Windows]
**Fenêtres sans cadre sous Windows :**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Windows: application.WindowsWindow{
        DisableFramelessWindowDecorations: false,
    },
})
```

**Fonctionnalités :**

- Ombre portée automatique
- Prise en charge de Snap Layouts (Windows 11)
- Prise en charge d’Aero Snap
- Mise à l’échelle selon le DPI

**Désactiver les décorations :**

```go
Windows: application.WindowsWindow{
    DisableFramelessWindowDecorations: true,
},
```

**Snap Assist :**

```go
// Trigger Windows 11 Snap Assist
window.SnapAssist()
```

Cela déclenche Snap Layouts par l’intermédiaire du raccourci clavier Windows. Pour créer un bouton HTML personnalisé d’agrandissement qui affiche Snap Layouts au survol comme un bouton natif, utilisez plutôt [Zones non clientes natives sous Windows](#zones-non-clientes-natives-sous-windows).

**Hauteur personnalisée de la barre de titre :** Windows détecte automatiquement les zones de déplacement définies avec CSS.

[macOS]
**Fenêtres sans cadre sous macOS :**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Mac: application.MacWindow{
        TitleBar: application.MacTitleBar{
            AppearsTransparent: true,
        },
        InvisibleTitleBarHeight: 40,
    },
})
```

**Fonctionnalités :**

- Prise en charge native du plein écran
- Boutons « feux de circulation » (facultatifs)
- Effets de vibrance
- Barre de titre transparente

**Masquer entièrement la barre de titre** (utilisez les variantes prédéfinies exportées par le paquet `application` — il n’existe ni champ `TitleBarStyle` ni constante `MacTitleBarStyleHidden`) :

```go
Mac: application.MacWindow{
    TitleBar: application.MacTitleBarHidden,
},
```

Les autres préréglages comprennent `MacTitleBarDefault`, `MacTitleBarHiddenInset` et `MacTitleBarHiddenInsetUnified`.

**Barre de titre invisible :** permet de déplacer la fenêtre tout en masquant sa barre de titre. Cette option ne prend effet que si la fenêtre est sans cadre ou utilise `AppearsTransparent` :

```go
Mac: application.MacWindow{
    InvisibleTitleBarHeight: 40,
},
```

[Linux]
**Fenêtres sans cadre sous Linux :**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
})
```

**Fonctionnalités :**

- Prise en charge de base des fenêtres sans cadre
- Zones de déplacement définies avec CSS
- Comportement variable selon l’environnement de bureau

**Remarques sur les environnements de bureau :**

- **GNOME :** bonne prise en charge
- **KDE Plasma :** bonne prise en charge
- **XFCE :** prise en charge de base
- **Gestionnaires de fenêtres en mosaïque :** prise en charge limitée

**Compositeur requis :** la transparence nécessite un compositeur (la plupart des environnements de bureau modernes en possèdent un).

@end

## Modèles courants

### Modèle 1 : barre de titre moderne

```html
<div class="modern-titlebar">
    <div class="app-icon">
        <img src="/icon.png" alt="App Icon">
    </div>
    <div class="title">My Application</div>
    <div class="controls">
        <button class="minimize">−</button>
        <button class="maximize">□</button>
        <button class="close">×</button>
    </div>
</div>
```

```css
.modern-titlebar {
    --wails-draggable: drag;
    display: flex;
    align-items: center;
    height: 40px;
    background: linear-gradient(to bottom, #3a3a3a, #2c2c2c);
    border-bottom: 1px solid #1a1a1a;
    padding: 0 16px;
}

.app-icon {
    --wails-draggable: no-drag;
    width: 24px;
    height: 24px;
    margin-right: 12px;
}

.title {
    flex: 1;
    font-size: 13px;
    color: #e0e0e0;
    user-select: none;
}

.controls {
    display: flex;
    gap: 1px;
}

.controls button {
    --wails-draggable: no-drag;
    width: 46px;
    height: 32px;
    border: none;
    background: transparent;
    color: #e0e0e0;
    font-size: 14px;
    cursor: pointer;
    transition: background 0.2s;
}

.controls button:hover {
    background: rgba(255, 255, 255, 0.1);
}

.controls .close:hover {
    background: #e81123;
    color: white;
}
```

### Modèle 2 : écran de démarrage

```go
splash := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:          "Loading...",
    Width:          400,
    Height:         300,
    Frameless:      true,
    AlwaysOnTop:    true,
    BackgroundType: application.BackgroundTypeTransparent,
    Mac: application.MacWindow{
        Backdrop: application.MacBackdropTransparent, // Requires -tags private_mac_apis.
    },
})
```

```css
body {
    background: transparent;
    display: flex;
    justify-content: center;
    align-items: center;
}

.splash {
    background: white;
    border-radius: 12px;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
    padding: 40px;
    text-align: center;
}
```

### Modèle 3 : fenêtre aux coins arrondis

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless:      true,
    BackgroundType: application.BackgroundTypeTransparent,
    Mac: application.MacWindow{
        Backdrop: application.MacBackdropTransparent, // Requires -tags private_mac_apis.
    },
})
```

```css
body {
    background: transparent;
    margin: 8px;
}

.window {
    background: white;
    border-radius: 16px;
    box-shadow: 0 4px 24px rgba(0, 0, 0, 0.15);
    overflow: hidden;
    height: calc(100vh - 16px);
}

.titlebar {
    --wails-draggable: drag;
    background: #f5f5f5;
    border-bottom: 1px solid #e0e0e0;
}
```

### Modèle 4 : fenêtre superposée

```go
overlay := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless:      true,
    AlwaysOnTop:    true,
    BackgroundType: application.BackgroundTypeTransparent,
    Mac: application.MacWindow{
        Backdrop: application.MacBackdropTransparent, // Requires -tags private_mac_apis.
    },
})
```

```css
body {
    background: transparent;
}

.overlay {
    background: rgba(0, 0, 0, 0.8);
    backdrop-filter: blur(10px);
    border-radius: 8px;
    padding: 20px;
}
```

## Exemple complet

Voici une fenêtre sans cadre prête pour la production :

**Go :**

```go
package main

import (
    _ "embed"
    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
    app := application.New(application.Options{
        Name: "Frameless App",
        Assets: application.AssetOptions{
            Handler: application.AssetFileServerFS(assets),
        },
    })

    window := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:     "Frameless Application",
        Width:     1000,
        Height:    700,
        MinWidth:  800,
        MinHeight: 600,
        Frameless: true,

        Mac: application.MacWindow{
            TitleBar: application.MacTitleBar{
                AppearsTransparent: true,
            },
            InvisibleTitleBarHeight: 40,
        },

        Windows: application.WindowsWindow{
            DisableFramelessWindowDecorations: false,
        },
    })

    window.Center()
    window.Show()

    app.Run()
}
```

**HTML :**

```html
<!DOCTYPE html>
<html>
<head>
    <link rel="stylesheet" href="/style.css">
</head>
<body>
    <div class="window">
        <div class="titlebar">
            <div class="title">Frameless Application</div>
            <div class="controls">
                <button class="minimize" title="Minimise">−</button>
                <button class="maximize" title="Maximise">□</button>
                <button class="close" title="Close">×</button>
            </div>
        </div>
        <div class="content">
            <h1>Hello from Frameless Window!</h1>
        </div>
    </div>
    <script src="/main.js" type="module"></script>
</body>
</html>
```

**CSS :**

```css
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

body {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    background: #f5f5f5;
}

.window {
    height: 100vh;
    display: flex;
    flex-direction: column;
}

.titlebar {
    --wails-draggable: drag;
    display: flex;
    justify-content: space-between;
    align-items: center;
    height: 40px;
    background: #ffffff;
    border-bottom: 1px solid #e0e0e0;
    padding: 0 16px;
}

.title {
    font-size: 13px;
    font-weight: 500;
    color: #333;
    user-select: none;
}

.controls {
    display: flex;
    gap: 8px;
}

.controls button {
    --wails-draggable: no-drag;
    width: 32px;
    height: 32px;
    border: none;
    background: transparent;
    color: #666;
    font-size: 16px;
    cursor: pointer;
    border-radius: 4px;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s;
}

.controls button:hover {
    background: #f0f0f0;
    color: #333;
}

.controls .close:hover {
    background: #e81123;
    color: white;
}

.content {
    flex: 1;
    padding: 40px;
    overflow: auto;
}
```

**JavaScript :**

```javascript
import { Window } from '@wailsio/runtime'

// Minimise button
document.querySelector('.minimize').addEventListener('click', () => {
    Window.Minimise()
})

// Maximise/restore button
const maximiseBtn = document.querySelector('.maximize')
maximiseBtn.addEventListener('click', async () => {
    const isMaximised = await Window.IsMaximised()

    if (isMaximised) {
        await Window.Restore()
    } else {
        await Window.Maximise()
    }

    updateMaximiseButton()
})

// Close button
document.querySelector('.close').addEventListener('click', () => {
    Window.Close()
})

// Update maximise button icon
async function updateMaximiseButton() {
    const isMaximised = await Window.IsMaximised()
    maximiseBtn.textContent = isMaximised ? '❐' : '□'
    maximiseBtn.title = isMaximised ? 'Restore' : 'Maximise'
}

// Initial state
updateMaximiseButton()
```

## Bonnes pratiques

### ✅ À faire

- **Prévoyez une zone de déplacement** — les utilisateurs doivent pouvoir déplacer la fenêtre
- **Implémentez les boutons système** — fermeture, réduction et agrandissement
- **Définissez une taille minimale** — évitez les dispositions inutilisables
- **Testez sur toutes les plateformes** — le comportement varie
- **Utilisez CSS pour les zones de déplacement** — cette solution est souple et facile à maintenir
- **Fournissez un retour visuel** — ajoutez des états au survol aux boutons

### ❌ À éviter

- **N’oubliez pas les poignées de redimensionnement** — si la fenêtre est redimensionnable
- **Ne rendez pas toute la fenêtre déplaçable** — cela empêche les interactions
- **N’oubliez pas d’exclure les boutons des zones de déplacement** — sinon, ils ne fonctionneront pas
- **N’utilisez pas de zones de déplacement minuscules** — elles sont difficiles à saisir
- **N’oubliez pas les différences entre les plateformes** — effectuez des tests approfondis

## Dépannage

### Impossible de déplacer la fenêtre

**Cause :** `--wails-draggable: drag` manquant

**Solution :**

```css
.titlebar {
    --wails-draggable: drag;
}
```

### Les boutons ne fonctionnent pas

**Cause :** les boutons se trouvent dans la zone de déplacement

**Solution :**

```css
.titlebar button {
    --wails-draggable: no-drag;
}
```

### Impossible de redimensionner la fenêtre

**Cause :** poignées de redimensionnement manquantes

**Solution :**

```css
body {
    --wails-resize: all;
}
```

## Étapes suivantes

@cards{cols="2"}
▣ Principes de base des fenêtres
Découvrez les principes fondamentaux de la gestion des fenêtres.

[En savoir plus →](/features/windows/basics/)

---
⚙ Options de fenêtre
Référence complète des options de fenêtre.

[En savoir plus →](/features/windows/options/)

---
🚀 Événements de fenêtre
Gérez les événements du cycle de vie des fenêtres.

[En savoir plus →](/features/windows/events/)

---
◆ Fenêtres multiples
Modèles pour les applications à plusieurs fenêtres.

[En savoir plus →](/features/windows/multiple/)

@end

---

**Des questions ?** Posez-les sur [Discord](https://discord.gg/JDdSxwjhGf) ou consultez l’[exemple de fenêtre sans cadre](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/frameless).
