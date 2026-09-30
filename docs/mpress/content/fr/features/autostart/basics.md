---
title: "Démarrage automatique"
description: "Enregistrez votre application pour qu’elle se lance à la connexion de l’utilisateur sous macOS, Windows et Linux"
slug: "features/autostart/basics"
sourcePath: "features/autostart/basics.md"
---

## Démarrage automatique

`app.Autostart` enregistre votre application pour qu’elle se lance automatiquement lorsque l’utilisateur se connecte. Il sélectionne le mécanisme natif adapté à chaque plateforme et résout les chemins d’installation comportant des liens symboliques (Homebrew, Scoop), afin que les enregistrements restent valides lors de la mise à niveau du binaire.

L’enregistrement prend effet à la **prochaine connexion**, et non immédiatement.

## Démarrage rapide

```go
import "github.com/wailsapp/wails/v3/pkg/application"

// Register to launch at login
if err := app.Autostart.Enable(); err != nil {
    app.Logger.Error("autostart enable failed", "error", err)
}

// Stop launching at login
if err := app.Autostart.Disable(); err != nil {
    app.Logger.Error("autostart disable failed", "error", err)
}

// Check status
enabled, err := app.Autostart.IsEnabled()
```

## API

### `Enable`

Enregistre l’application pour qu’elle se lance à la connexion avec les options par défaut.

```go
func (m *AutostartManager) Enable() error
```

Vous pouvez appeler `Enable` plusieurs fois sans risque : l’enregistrement est remplacé à chaque appel. Vous pouvez donc l’appeler à chaque démarrage si vous avez conservé la préférence de l’utilisateur.

### `EnableWithOptions`

Effectue l’enregistrement avec des options personnalisées.

```go
func (m *AutostartManager) EnableWithOptions(opts AutostartOptions) error
```

**`AutostartOptions` :**

| Champ | Type | Description |
| --- | --- | --- |
| `Identifier` | `string` | Remplace l’identifiant d’enregistrement généré automatiquement. Consultez la section « Identifiant » ci-dessous. |
| `Arguments` | `[]string` | Arguments supplémentaires ajoutés au chemin de l’exécutable lors du lancement à la connexion (par exemple, `--hidden`). |

### `Disable`

Supprime l’enregistrement du démarrage automatique. Renvoie `nil` si l’application n’était pas enregistrée : la désactivation est idempotente.

```go
func (m *AutostartManager) Disable() error
```

### `IsEnabled`

Indique s’il existe un enregistrement. Cette opération est rapide, car elle ne valide pas le chemin enregistré.

```go
func (m *AutostartManager) IsEnabled() (bool, error)
```

### `Status`

Renvoie l’état complet de l’enregistrement.

```go
func (m *AutostartManager) Status() (AutostartStatus, error)
```

**`AutostartStatus` :**

| Champ | Type | Description |
| --- | --- | --- |
| `Enabled` | `bool` | Indique s’il existe un enregistrement. |
| `Path` | `string` | Emplacement sur le disque de l’artefact d’enregistrement (chemin du fichier plist, chemin `.desktop` ou sous-clé du Registre). Vide lorsque `Enabled` vaut false. |
| `Strategy` | `AutostartStrategy` | Mécanisme ayant enregistré l’application (consultez la section [Comportement selon la plateforme](#comportement-selon-la-plateforme)). |

## Comportement selon la plateforme

@tabs{sync-key="platform"}
[macOS]
Deux mécanismes sont utilisés selon la manière dont l’application est empaquetée :

- **macOS 13 ou version ultérieure, `.app`** empaquetée : `SMAppService.mainAppService`. Fonctionne avec les applications en bac à sable et les versions destinées au Mac App Store. Aucune invite d’automatisation TCC ne s’affiche (l’ancienne approche fondée sur AppleScript en déclenchait une).
- **Versions de macOS antérieures à 13, ou binaire non empaqueté** : un fichier plist LaunchAgent est écrit dans `~/Library/LaunchAgents/<identifier>.plist` avec `RunAtLoad=true`.

`Status()` renvoie `AutostartStrategySMAppService` ou `AutostartStrategyLaunchAgent` afin que le code appelant puisse déterminer le mécanisme utilisé.

Lorsque l’application passe d’une version non empaquetée à une version empaquetée lors d’une mise à niveau, `Status()` vérifie les deux mécanismes et `Disable()` les nettoie, afin qu’un LaunchAgent orphelin ne continue pas à lancer l’ancienne version.

[Windows]
Une valeur de Registre est ajoutée sous `HKCU\Software\Microsoft\Windows\CurrentVersion\Run`. Son nom correspond à l’`Identifier` de démarrage automatique et ses données contiennent le chemin de l’exécutable entre guillemets, suivi des éventuels `Arguments`.

La mise entre guillemets des arguments respecte les règles de `CommandLineToArgvW` (les barres obliques inverses sont doublées avant les guillemets), afin que les chemins contenant des espaces ou des guillemets soient restitués correctement.

`Status().Strategy` renvoie `AutostartStrategyRegistryRun`.

[Linux]
Une entrée de démarrage automatique XDG est écrite dans `$XDG_CONFIG_HOME/autostart/<identifier>.desktop` (`~/.config/autostart/` par défaut) avec :

```ini
Type=Application
Hidden=false
X-GNOME-Autostart-enabled=true
Exec=<executable> <arguments>
```

Le champ `Exec` est échappé conformément à la [spécification Desktop Entry de freedesktop.org](https://specifications.freedesktop.org/desktop-entry-spec/) : les caractères réservés (`"`, `` ` ``, `$`, `\\`) sont précédés d’une barre oblique inverse, et la valeur est placée entre guillemets doubles lorsqu’elle contient des caractères blancs.

`Status().Strategy` renvoie `AutostartStrategyXDGAutostart`.

[iOS / Android / serveur]
Non pris en charge. Toutes les méthodes renvoient `ErrAutostartNotSupported`. Utilisez `errors.Is(err, application.ErrAutostartNotSupported)` pour détecter proprement ce cas :

```go
if err := app.Autostart.Enable(); err != nil {
    if errors.Is(err, application.ErrAutostartNotSupported) {
        // hide the toggle in the UI
        return
    }
    // real failure — surface it
}
```

@end

## Identifiant

Si `Options.Identifier` est vide, une valeur par défaut est dérivée du nom de votre application :

| Plateforme | Identifiant par défaut |
| --- | --- |
| macOS (empaquetée) | Identifiant de bundle de l’application, par exemple `com.example.MyApp` |
| macOS (non empaquetée) | `wails.autostart.<slug>`, où `<slug>` est dérivé de `application.Options.Name` |
| Windows | Slug de `application.Options.Name` (en minuscules, caractères autres que `A-Za-z0-9._-` supprimés, espaces remplacés par des tirets) |
| Linux | Même slug que sous Windows |

Les identifiants doivent correspondre à `^[A-Za-z0-9._-]+$` et ne pas dépasser 200 caractères. Le format DNS inversé est recommandé sous macOS (il correspond à la convention d’écriture des Labels de launchd).

Lorsque `AutostartOptions.Identifier` est remplacé, le même identifiant est réutilisé comme nom de valeur du Registre sous Windows et comme nom de fichier `.desktop` sous Linux. Une seule chaîne permet donc d’identifier l’enregistrement sur toutes les plateformes.

## Détection des enregistrements obsolètes

`Disable()` et `Status()` localisent l’enregistrement en **comparant le chemin de l’exécutable enregistré à `os.Executable()` (après résolution des éventuels liens symboliques)**, et non en recherchant l’identifiant. Cela signifie que :

- **Vous pouvez modifier l’identifiant entre deux versions sans risque.** L’ancien enregistrement reste détectable par `Status()` et est supprimé par `Disable()`, à condition que le chemin de l’exécutable reste identique.
- **Une seconde copie de l’application située à un autre emplacement n’écrasera pas l’enregistrement de la première.** Chaque emplacement du binaire est suivi indépendamment.
- **Les installations utilisant des liens symboliques (Homebrew, Scoop) restent stables.** `filepath.EvalSymlinks` est appliqué à `os.Executable()` avant la comparaison. Ainsi, une mise à niveau Homebrew qui remplace la cible ne laisse pas l’entrée orpheline.

Ce mécanisme ne couvre *pas* le cas suivant : si l’utilisateur déplace ou renomme le binaire vers un chemin sans rapport avec le précédent, l’ancien enregistrement devient orphelin (il pointe vers le fichier désormais absent). Les applications distribuées depuis un chemin d’installation stable n’ont pas à s’en préoccuper. Celles distribuées sous forme de binaires portables constitués d’un seul fichier devraient appeler `Disable()` avant de se déplacer elles-mêmes, ou toujours être lancées au moyen d’un lien symbolique stable.

## Exemple

```go
package main

import (
    "errors"

    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "My App",
    })

    // Restore the user's preference on startup
    if userPrefersAutostart() {
        if err := app.Autostart.Enable(); err != nil {
            if !errors.Is(err, application.ErrAutostartNotSupported) {
                app.Logger.Error("autostart", "error", err)
            }
        }
    }

    app.Run()
}
```

Un exemple complet et exécutable comportant des boutons d’état, d’activation et de désactivation se trouve dans [`examples/autostart/`](https://github.com/wailsapp/wails/tree/master/v3/examples/autostart).
