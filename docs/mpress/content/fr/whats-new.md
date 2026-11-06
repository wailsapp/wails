---
title: "Nouveautés de Wails v3"
description: "Découvrez les principales améliorations et nouveautés de Wails v3"
slug: "whats-new"
sourcePath: "whats-new.md"
---

Wails v3 apporte des changements importants par rapport à v2. L’API déclarative à fenêtre unique est remplacée par une approche procédurale plus flexible. Cette nouvelle conception de l’API améliore la lisibilité du code et simplifie le développement, en particulier pour les applications complexes à plusieurs fenêtres.

Wails v3 marque une évolution majeure dans la manière de créer des applications de bureau avec Go et les technologies web.

## Fenêtres multiples

Wails v3 permet désormais de créer et de gérer plusieurs fenêtres au sein d’une même application. Cette fonctionnalité permet aux développeurs de concevoir des interfaces utilisateur plus complexes et polyvalentes, sans les limitations des applications à fenêtre unique.

Chaque fenêtre peut être configurée indépendamment, notamment en matière de taille, de position, de contenu et de comportement. Il est ainsi possible de créer des applications dotées de fenêtres distinctes pour différentes fonctions, telles que l’interface principale, les panneaux de paramètres ou les vues auxiliaires.

Les développeurs peuvent créer, manipuler et gérer ces fenêtres par programmation afin de produire des interfaces utilisateur dynamiques qui s’adaptent aux besoins des utilisateurs et aux états de l’application.

@note{type="tip" title="Fenêtres multiples"}
@details{title="Exemple"}
```go
package main

import (
   "embed"
   "log"
   
   "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets/*
var assets embed.FS

func main() {

   app := application.New(application.Options{
        Name:   "Multi Window Demo",
        Assets: application.AssetOptions{
            Handler: application.AssetFileServerFS(assets),
        },
   })
   
   window1 := app.Window.NewWithOptions(application.WebviewWindowOptions{
       Title:  "Window 1",
   })
   
   window2 := app.Window.NewWithOptions(application.WebviewWindowOptions{
       Title:  "Window 2",
   })
   
   // load the embedded html from the embed.FS
   window1.SetURL("/")
   window1.Center()
   
   // Load an external URL
   window2.SetURL("https://wails.io")
   
   err := app.Run()

   if err != nil {
	   log.Fatal(err.Error())
   }
}
```

@end

@end

## Intégration à la zone de notification

Wails v3 prend désormais pleinement en charge la zone de notification, ce qui permet à votre application de rester accessible en permanence sur le bureau de l’utilisateur. Cette fonctionnalité est particulièrement utile pour les applications qui doivent s’exécuter en arrière-plan ou offrir un accès rapide à des fonctions essentielles.

Les principales fonctionnalités de l’intégration à la zone de notification dans Wails v3 sont les suivantes :

1. Association d’une fenêtre : vous pouvez associer une fenêtre à l’icône de la zone de notification. Lorsqu’elle est activée, cette fenêtre est centrée par rapport à la position de l’icône, ce qui offre un excellent moyen d’accéder rapidement à votre application.

2. Prise en charge complète des menus : créez des menus riches et interactifs auxquels les utilisateurs peuvent accéder directement depuis l’icône de la zone de notification. Ils peuvent ainsi effectuer des actions rapides sans avoir à ouvrir la fenêtre complète de l’application.

3. Affichage adaptatif de l’icône : la prise en charge d’icônes pour les modes clair et sombre garantit que l’icône de votre application dans la zone de notification reste visible et esthétique avec les différents thèmes du système. Les icônes de modèle sont également prises en charge sous macOS.

@note{type="tip" title="Zone de notification"}
@details{title="Exemple"}
```go
package main

import (
    "log"
    "runtime"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/icons"
)

func main() {
    app := application.New(application.Options{
        Name:        "Systray Demo",
        Mac: application.MacOptions{
            ActivationPolicy: application.ActivationPolicyAccessory,
        },
    })

    window := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Width:       500,
        Height:      800,
        Frameless:   true,
        AlwaysOnTop: true,
        Hidden:      true,
        Windows: application.WindowsWindow{
            HiddenOnTaskbar: true,
        },
    })

    systemTray := app.SystemTray.New()

    // Support for template icons on macOS
    if runtime.GOOS == "darwin" {
        systemTray.SetTemplateIcon(icons.SystrayMacTemplate)
    } else {
        // Support for light/dark mode icons
        systemTray.SetDarkModeIcon(icons.SystrayDark)
        systemTray.SetIcon(icons.SystrayLight)
    }

    // Support for menu
    myMenu := app.Menu.New()
    myMenu.Add("Hello World!").OnClick(func(_ *application.Context) {
        println("Hello World!")
    })
    systemTray.SetMenu(myMenu)

    // This will center the window to the systray icon with a 5px offset
    // It will automatically be shown when the systray icon is clicked
    // and hidden when the window loses focus
    systemTray.AttachWindow(window).WindowOffset(5)

    err := app.Run()
    if err != nil {
        log.Fatal(err)
    }
}
```

@end

@end

## Amélioration de la génération des liaisons

Wails v3 améliore considérablement la génération des liaisons de votre projet. Les liaisons assurent la connexion entre votre backend Go et votre frontend, et permettent une communication fluide entre les deux.

La génération des liaisons repose désormais sur un analyseur statique sophistiqué qui améliore radicalement le processus. Celui-ci accélère la génération et préserve la qualité du code en conservant les commentaires et les noms des paramètres.

Le processus de génération des liaisons a été simplifié et ne nécessite plus qu’une seule commande : `wails3 generate bindings`.

@note{type="tip" title="Liaisons"}
@details{title="Exemple"}
```js
// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

// Generated layout (excerpt): frontend/bindings/<full-go-import-path>/greetservice.js
import { Call as $Call, Create as $Create } from "/wails/runtime.js";

/**
 * Greet greets a person
 * @param {string} $0
 * @returns {Promise<string>}
 */
export function Greet($0) {
    return $Call.ByID(1411160069, $0);
}

/**
 * GreetPerson greets a person
 * @param {main.Person} $0
 * @returns {Promise<string>}
 */
export function GreetPerson($0) {
    return $Call.ByID(4021313248, $0);
}
```

@end

@end

## Amélioration du système de compilation

Wails v3 introduit un système de compilation plus flexible et transparent qui remédie aux limites de son prédécesseur. Dans v2, le processus de compilation était en grande partie opaque et difficile à personnaliser, ce qui pouvait frustrer les développeurs souhaitant mieux maîtriser le processus de compilation de leur projet.

Toutes les opérations complexes effectuées par le système de compilation de v2, telles que la génération des icônes et la création du manifeste, ont été ajoutées à la CLI sous forme de commandes d’outils. Nous avons intégré [Taskfile](https://taskfile.dev) à la CLI afin d’orchestrer ces appels et d’offrir la même expérience de développement qu’avec v2. Cette approche offre toutefois un équilibre optimal entre flexibilité et simplicité d’utilisation, puisque vous pouvez désormais adapter le processus de compilation à vos besoins.

Vous pouvez même utiliser make si vous le souhaitez !

@note{type="tip" title="Taskfile.yml"}
@details{title="Exemple"}
```yaml {title="build/Taskfile.darwin.yml"}
darwin:build:
  summary: Builds the application for macOS
  platforms:
    - darwin
  cmds:
    - task: common:go:mod:tidy
    - task: common:build:frontend
    - task: common:generate:icons
    - task: darwin:build:app
  env:
    CGO_CFLAGS: "-mmacosx-version-min=10.15"
    CGO_LDFLAGS: "-mmacosx-version-min=10.15"
    MACOSX_DEPLOYMENT_TARGET: "10.15"
```

@end

@end

## Amélioration des événements

Wails émet désormais des événements pour diverses opérations d’exécution et activités du système. Votre application peut ainsi réagir à ces événements en temps réel. Des événements multiplateformes (communs) sont également disponibles, ce qui vous permet d’écrire des méthodes cohérentes de gestion des événements qui fonctionnent sur différents systèmes d’exploitation.

Vous pouvez enregistrer des gestionnaires d’événements pour traiter certains événements de manière synchrone. Contrairement à la méthode `On`, ces gestionnaires permettent d’annuler l’événement si nécessaire. Un cas d’usage courant consiste à afficher une boîte de dialogue de confirmation avant de fermer une fenêtre. Vous maîtrisez ainsi davantage le déroulement des événements et l’expérience utilisateur.

@note{type="tip" title="Exemple de gestion des événements"}
@details{title="Exemple"}
```go
package main

import (
    "embed"
    "log"
    "time"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed assets
var assets embed.FS

func main() {

    app := application.New(application.Options{
        Name:        "Events Demo",
        Description: "A demo of the Events API",
        Assets: application.AssetOptions{
            Handler: application.AssetFileServerFS(assets),
        },
        Mac: application.MacOptions{
            ApplicationShouldTerminateAfterLastWindowClosed: true,
        },
    })

    // Custom event handling — App.Event.On(name, func(e *CustomEvent))
    app.Event.On("myevent", func(e *application.CustomEvent) {
        log.Printf("[Go] CustomEvent received: %+v\n", e)
    })

    // OS-specific application events — App.Event.OnApplicationEvent(eventType, func(e *ApplicationEvent))
    app.Event.OnApplicationEvent(events.Mac.ApplicationDidFinishLaunching, func(event *application.ApplicationEvent) {
        println("events.Mac.ApplicationDidFinishLaunching fired!")
    })

    // Platform-agnostic events
    app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(event *application.ApplicationEvent) {
        println("events.Common.ApplicationStarted fired!")
    })

    win1 := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title: "Takes 3 attempts to close me!",
    })

    var countdown = 3

    // Register a hook to cancel the window closing
    win1.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
        countdown--
        if countdown == 0 {
            println("Closing!")
            return
        }
        println("Nope! Not closing!")
        e.Cancel()
    })

    win1.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
        println("[Event] Window focus!")
    })

    err := app.Run()

    if err != nil {
        log.Fatal(err.Error())
    }
}
```

@end

@end

## Langage de balisage Wails (wml)

Une fonctionnalité expérimentale permettant d’appeler des méthodes d’exécution à l’aide de HTML standard, à la manière de [htmx](https://htmx.org).

@note{type="tip" title="Exemple de wml"}
@details{title="Exemple"}
```html
<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <title>Wails ML Demo</title>
  </head>
  <body style="margin-top:50px; color: white; background-color: #191919">
    <h2>Wails ML Demo</h2>
    <p>This application contains no Javascript!</p>
    <button wml-event="button-pressed">Press me!</button>
    <button wml-event="delete-things" wml-confirm="Are you sure?">
      Delete all the things!
    </button>
    <button wml-window="Close" wml-confirm="Are you sure?">
      Close the Window?
    </button>
    <button wml-window="Center">Center</button>
    <button wml-window="Minimise">Minimise</button>
    <button wml-window="Maximise">Maximise</button>
    <button wml-window="UnMaximise">UnMaximise</button>
    <button wml-window="Fullscreen">Fullscreen</button>
    <button wml-window="UnFullscreen">UnFullscreen</button>
    <button wml-window="Restore">Restore</button>
    <div
      style="width: 200px; height: 200px; border: 2px solid white;"
      wml-event="hover"
      wml-trigger="mouseover"
    >
      Hover over me
    </div>
  </body>
</html>
```

@end

@end

## Exemples

D’autres exemples sont disponibles dans le répertoire [examples](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples). Consultez-les !
