---
title: "API privées de macOS"
description: "Toutes les fonctionnalités et options de Wails qui dépendent d’API privées de macOS, avec les commandes d’activation explicite et les solutions de repli pour les builds publics."
slug: "guides/build/private-macos-apis"
sourcePath: "guides/build/private-macos-apis.md"
---

Wails utilise par défaut les API publiques de macOS. L’unique balise de build Go `private_mac_apis` active les appels privés à WebKit et AppKit répertoriés sur cette page. Toutes les options et méthodes Go publiques restent disponibles dans les deux types de builds. Sans cette balise, les opérations exclusivement privées sont sans effet ; les fonctionnalités disposant de solutions publiques utilisent celles-ci.

@note{type="caution" title="Activer explicitement les comportements privés de macOS"}
La définition d’une option de fenêtre n’active pas les API privées. Ajoutez `private_mac_apis` à la commande de build pour les activer. Cette balise s’applique uniquement aux builds de bureau macOS, et non aux builds iOS, Android, Windows, Linux ou serveur.

@end

## Activer les API privées

```bash
# Build an application
wails3 build -tags private_mac_apis

# Run with live reload
EXTRA_TAGS=private_mac_apis wails3 dev

# Package an application
wails3 package GOOS=darwin EXTRA_TAGS=private_mac_apis

# Run a Go-only example from its directory
go run -tags private_mac_apis .
```

Pour un build de production direct, utilisez `go build -tags production,private_mac_apis .`. Pour les exemples de frontend, suivez leur README afin de générer les bindings et les ressources avant l’exécution. Les Taskfiles personnalisés ou anciens doivent transmettre `EXTRA_TAGS` au compilateur Go.

## Inventaire des fonctionnalités

| Fonctionnalité ou valeur | Ce qu’active `private_mac_apis` | Sans la balise |
| --- | --- | --- |
| `Mac.Backdrop: MacBackdropTransparent` | WKWebView transparent au-dessus de la fenêtre native | La fenêtre native est configurée, mais la vue web reste opaque |
| `Mac.Backdrop: MacBackdropTranslucent` | WKWebView transparent laissant apparaître le flou natif | Le flou est configuré derrière une vue web opaque |
| `Mac.Backdrop: MacBackdropLiquidGlass` | WKWebView transparent au-dessus de la couche de verre, avec contrôle privé de l’arrière-plan de la vue web | La couche de verre est configurée derrière une vue web opaque ; le style utilise des solutions publiques |
| Effacement de l’arrière-plan de la vue web pendant la configuration de Liquid Glass | Contrôle WebKit privé `backgroundColor` | `underPageBackgroundColor` public sous macOS 12 ou version ultérieure, ou couleur de couche sous les versions antérieures de macOS ; ne rend pas la vue web transparente |
| `app.Window.NewNotchWindow(...)` | Vue web transparente dans le panneau en forme d’encoche | Le panneau continue de fonctionner, y compris son positionnement et ses animations, mais sa vue web reste opaque |
| `Mac.LiquidGlass.Style` | Mappage existant de Wails vers les styles natifs, y compris une valeur de style sombre non documentée | Utilise les styles publics normal/transparent et les apparences claire/sombre ; consultez le tableau des valeurs ci-dessous |
| `Mac.LiquidGlass.GroupID` | Demande un regroupement privé des effets de verre pour un identifiant non vide | Ignoré ; aucun regroupement n’est demandé |
| `Mac.LiquidGlass.GroupSpacing` | Demande un espacement privé du groupe pour une valeur supérieure à zéro | Ignoré |
| `window.OpenDevTools()` et `Window.OpenDevTools()` JavaScript | Ouvre l’inspecteur WebKit par programmation sous macOS 12 ou version ultérieure | Sans effet |
| `WebviewWindowOptions.OpenInspectorOnStartup: true` | Demande l’ouverture de l’inspecteur par programmation lors du premier affichage de la fenêtre, sous macOS 12 ou version ultérieure | Sans effet |
| Activation historique de l’inspecteur, avant macOS 13.3 | Active les fonctionnalités supplémentaires pour développeurs de WebKit lorsque la prise en charge de l’inspecteur est incluse dans le build | Sans effet ; l’inspection publique avec Safari nécessite macOS 13.3 ou version ultérieure |

## Transparence et arrière-plan de la vue web

**Nécessite des API privées :** la transparence de la vue web utilisée par `MacBackdropTransparent`, `MacBackdropTranslucent`, `MacBackdropLiquidGlass` et les fenêtres à encoche. En interne, Wails définit la clé privée `drawsBackground` de WebKit. Un arrière-plan HTML ou CSS transparent ne peut pas, à lui seul, rendre transparente une WKWebView native opaque.

```go
Mac: application.MacWindow{
    // Requires -tags private_mac_apis for the blur to show through the webview.
    Backdrop: application.MacBackdropTranslucent,
},
```

L’opération privée de définition de la couleur d’arrière-plan de la vue web utilise la clé `backgroundColor` de WebKit. La configuration de Liquid Glass l’utilise pour effacer l’arrière-plan de la vue web. Sans la balise, cette opération interne utilise l’API publique `underPageBackgroundColor` sous macOS 12 ou version ultérieure, ou la couche de la vue sous les versions antérieures de macOS. Ces solutions ne rendent pas la vue web transparente.

`WebviewWindowOptions.BackgroundColour` et `window.SetBackgroundColour()` définissent sous macOS la couleur de la **fenêtre native** et ne nécessitent pas eux-mêmes d’API privées. De même, `Frameless` et `Mac.TitleBar.AppearsTransparent` utilisent des API AppKit publiques ; la dépendance privée concerne la transparence de la vue web, et non celle de la barre de titre. Pour les effets de fond sous macOS, configurez `Mac.Backdrop` au lieu de vous fier uniquement à `BackgroundType`.

Consultez les pages consacrées aux [options de fenêtre](/features/windows/options/#mac-options), aux [fenêtres sans cadre](/features/windows/frameless/#with-transparent-background) et aux [fenêtres à encoche](/features/windows/notch-windows/).

## Valeurs de Liquid Glass

Lorsqu’un `NSGlassEffectView` natif est disponible (macOS 26 ou version ultérieure), les mappages suivants s’appliquent. Seules les valeurs de style natif `0` (normal) et `1` (transparent) sont documentées. Les constantes Go conservent leurs valeurs existantes dans les deux types de builds.

| Valeur de `MacLiquidGlassStyle` | Avec `private_mac_apis` | Sans la balise |
| --- | --- | --- |
| `LiquidGlassStyleAutomatic` (`0`) | Style natif normal (`0`) | Style natif normal (`0`) |
| `LiquidGlassStyleLight` (`1`) | Mappage existant vers le style natif (`1`, transparent) | Style natif normal (`0`) avec l’apparence Aqua |
| `LiquidGlassStyleDark` (`2`) | **Valeur de style native non documentée `2`** | Style standard natif (`0`) avec l’apparence Aqua sombre |
| `LiquidGlassStyleVibrant` (`3`) | Correspond au style clair/translucide natif existant (`1`) | Style translucide natif (`1`) |

Automatic et Vibrant utilisent des valeurs de style natives documentées, mais un arrière-plan Liquid Glass couvrant toute la fenêtre nécessite toujours le tag pour la **transparence de la webview**. L’apparence de Light diffère entre les deux builds. Le mappage de style privé ne garantit pas qu’une future version de macOS produira le même effet.

**Toujours privées :** `GroupID` et `GroupSpacing`. Wails vérifie la présence des sélecteurs privés `setGroupIdentifier:`, `setGroupName:` et `setGroupSpacing:` avant de demander le regroupement. Sans le tag, ces opérations sont sans effet. L’activation du tag ne garantit pas que la version de macOS en cours d’exécution prend en charge ces sélecteurs.

`MacLiquidGlass.Material`, `CornerRadius` et `TintColor` ne nécessitent pas eux-mêmes d’API privées. Sur les versions de macOS dépourvues de prise en charge native de Liquid Glass, Wails configure une solution de repli translucide ; pour qu’elle soit visible à travers la webview, le tag reste nécessaire.

## Inspecteur web

**Nécessite des API privées :** l’appel de `OpenDevTools()` ou la définition de `OpenInspectorOnStartup: true` pour ouvrir l’inspecteur de WebKit depuis l’application. Wails utilise le sélecteur privé `_inspector`. Sans `private_mac_apis`, ces opérations sont silencieusement sans effet.

La prise en charge de l’inspecteur doit également être intégrée lors de la compilation. Les tags `production` et `devtools` existants conservent leur signification :

| Tags de build | Ouverture programmatique de l’inspecteur | Inspection publique avec Safari sous macOS 13.3+ |
| --- | --- | --- |
| Aucun | Sans effet | Activée |
| `private_mac_apis` | Activée sous macOS 12+ | Activée |
| `production` | Sans effet | Désactivée |
| `production,private_mac_apis` | Sans effet | Désactivée |
| `production,devtools` | Sans effet | Activée |
| `production,devtools,private_mac_apis` | Activée sous macOS 12+ | Activée |

Sous macOS 13.3+, Wails active l’inspection avec Safari à l’aide de l’API publique `WKWebView.inspectable` ; cela ne nécessite aucune API privée. Avant macOS 13.3, la solution de repli utilise la préférence privée `developerExtrasEnabled` et nécessite donc `private_mac_apis` ainsi que la prise en charge de l’inspecteur intégrée lors de la compilation.

```bash
# Production build with programmatic inspector support
go build -tags production,devtools,private_mac_apis .
```

## Tenir cette liste à jour

Tous les appels natifs privés sont isolés dans `v3/pkg/application/mac_private_api_darwin.go` ; les builds par défaut sélectionnent `mac_public_api_darwin.go`. Cet inventaire couvre la transparence, la couleur d’arrière-plan de la webview, les styles de verre, le regroupement des effets de verre, l’ouverture de l’inspecteur et l’activation de l’ancien inspecteur. Toute modification de ces implémentations devrait s’accompagner d’une mise à jour conjointe de cette page et de la documentation des options concernées.
