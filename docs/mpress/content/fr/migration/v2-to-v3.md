---
title: "Migration de la v2 vers la v3"
description: "Guide complet pour migrer votre application Wails v2 vers la v3"
slug: "migration/v2-to-v3"
sourcePath: "migration/v2-to-v3.md"
---

Wails v3 est une **réécriture complète** qui apporte des améliorations majeures en matière d’architecture, de performances et d’expérience de développement. Ce guide vous aide à migrer votre application de la v2 vers la v3.

**Principales modifications :**

- Nouvelle structure d’application
- Système de liaisons amélioré
- Gestion des fenêtres améliorée
- Meilleur système d’événements
- Configuration simplifiée

**Durée de la migration :** 1-4 heures pour une application classique

## Modifications incompatibles

### Initialisation de l’application

Dans la v2, la configuration de l’application, la configuration des fenêtres et l’exécution étaient toutes regroupées dans un seul appel à `wails.Run()`. Cette approche monolithique compliquait la création de plusieurs fenêtres, la gestion des erreurs aux différentes étapes et le test séparé des composants de votre application.

La v3 sépare ces responsabilités en plusieurs phases distinctes : création de l’application, création des fenêtres et exécution. Cette séparation vous donne un contrôle explicite sur chaque étape du cycle de vie de votre application et rend le code plus modulaire et plus facile à tester.

**v2 :**

```go
err := wails.Run(&options.App{
    Title:  "My App",
    Width:  1024,
    Height: 768,
    Bind: []interface{}{
        &GreetService{},
    },
})
```

**v3 :**

```go
app := application.New(application.Options{
    Name: "My App",
    Services: []application.Service{
        application.NewService(&GreetService{}),
    },
})

window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "My App",
    Width:  1024,
    Height: 768,
})

app.Run()
```

**Avantages de cette approche :**

- **Prise en charge de plusieurs fenêtres** : vous pouvez créer des fenêtres dynamiquement à tout moment, et pas seulement au démarrage
- **Meilleure gestion des erreurs** : chaque phase peut être validée séparément avec une gestion appropriée des erreurs
- **Code plus clair** : cette séparation permet de comprendre immédiatement ce qui se passe à chaque étape
- **Testabilité accrue** : vous pouvez tester la configuration de l’application sans exécuter la boucle d’événements
- **Flexibilité accrue** : les fenêtres peuvent être créées, détruites et recréées tout au long du cycle de vie de l’application

### Liaisons

Dans la v2, chaque structure liée devait comporter un champ de contexte et une méthode `startup(ctx)` pour recevoir le contexte d’exécution. Cela créait un couplage étroit entre votre logique métier et l’environnement d’exécution de Wails, ce qui rendait le code plus difficile à tester et à comprendre.

La v3 introduit le modèle de services, dans lequel vos structures sont entièrement autonomes et n’ont pas besoin de stocker le contexte d’exécution. Si un service doit accéder à l’instance de l’application, il la reçoit explicitement par injection de dépendances plutôt que par transmission implicite du contexte.

**v2 :**

```go
type App struct {
    ctx context.Context
}

func (a *App) startup(ctx context.Context) {
    a.ctx = ctx
}

func (a *App) Greet(name string) string {
    return "Hello " + name
}
```

**v3 :**

```go
type GreetService struct{}

func (g *GreetService) Greet(name string) string {
    return "Hello " + name
}

// Register as service
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&GreetService{}),
    },
})
```

**Avantages de cette approche :**

- **Aucune dépendance implicite** : les services sont de simples structures Go sans dépendances cachées à l’environnement d’exécution
- **Tests plus simples** : vous pouvez tester les méthodes des services sans simuler de contexte Wails
- **Code plus clair** : les dépendances sont explicites, car elles sont transmises comme arguments du constructeur, plutôt que cachées dans un champ de contexte
- **Meilleure organisation** : les services peuvent être regroupés par domaine au lieu d’être tous réunis dans une seule structure `App`
- **Initialisation appropriée** : utilisez la méthode `ServiceStartup()` lorsqu’une initialisation est nécessaire, afin de la rendre explicite

### Environnement d’exécution

Dans la v2, toutes les opérations d’exécution nécessitaient de transmettre un contexte à des fonctions globales du package `runtime`. Cela créait un couplage étroit avec l’objet de contexte dans toute la base de code et donnait à l’API un caractère procédural plutôt qu’orienté objet.

La v3 remplace l’environnement d’exécution fondé sur le contexte par des appels directs de méthodes sur les objets d’application et de fenêtre. Les opérations sont appelées directement sur les objets qu’elles affectent, ce qui rend le code plus intuitif et davantage orienté objet.

**v2 :**

```go
import "github.com/wailsapp/wails/v2/pkg/runtime"

runtime.WindowSetTitle(a.ctx, "New Title")
runtime.EventsEmit(a.ctx, "event-name", data)
```

**v3 :**

```go
// Store app reference
type MyService struct {
    app *application.App
}

func (s *MyService) UpdateTitle() {
    window := s.app.Window.Current()
    window.SetTitle("New Title")
}

func (s *MyService) EmitEvent() {
    s.app.Event.Emit("event-name", data)
}
```

**Avantages de cette approche :**

- **Conception orientée objet** : les méthodes sont appelées sur les objets qu’elles affectent (fenêtre, application, menu, etc.)
- **Intention plus claire** : `window.SetTitle()` est plus explicite que `runtime.WindowSetTitle(ctx, ...)`
- **Meilleure prise en charge par les IDE** : la saisie semi-automatique fonctionne correctement lorsque les méthodes appartiennent aux objets
- **Gestion plus claire de plusieurs fenêtres** : lorsque plusieurs fenêtres sont ouvertes, vous choisissez explicitement celle sur laquelle agir
- **Aucune transmission du contexte** : vous n’avez pas besoin de transmettre le contexte à chaque fonction

### Liaisons de l’interface utilisateur

Dans la v2, les liaisons étaient organisées selon le package Go et le nom de la structure, ce qui produisait généralement des chemins tels que `wailsjs/go/main/App`. Cette structure ne reflétait pas le regroupement logique et compliquait la recherche des fonctionnalités connexes.

La v3 organise les liaisons par nom de service et par module d’application, créant ainsi une structure logique plus claire. Les liaisons sont générées dans un répertoire `bindings`, organisé selon le nom de votre application et les noms des services, ce qui permet de comprendre plus facilement les fonctionnalités disponibles.

**v2 :**

```javascript
import { Greet } from '../wailsjs/go/main/App'

const result = await Greet("World")
```

**v3 :**

```javascript
import { Greet } from './bindings/changeme/greetservice'

const result = await Greet("World")
```

**Avantages de cette approche :**

- **Organisation logique** : les liaisons sont regroupées par nom de service plutôt que selon la structure des packages Go
- **Importations plus claires** : le chemin reflète la logique du domaine (greetservice), et non la structure des fichiers (main/App)
- **Repérage facilité** : vous pouvez parcourir les bindings par fonctionnalité plutôt que par structure technique
- **Nommage cohérent** : l’organisation par service correspond à l’architecture de votre backend
- **Chemins simplifiés** : plus de préfixe `../wailsjs/go`, seulement `./bindings`

### Événements

Dans la v2, les événements utilisaient des paramètres variadiques `interface{}` et exigeaient de transmettre le contexte à chaque fonction d’événement. Les gestionnaires d’événements recevaient des données non typées qui nécessitaient des assertions de type manuelles, ce qui rendait le système d’événements sujet aux erreurs et difficile à déboguer.

La v3 introduit des objets d’événement typés et supprime l’obligation de fournir le contexte. Les gestionnaires d’événements reçoivent un véritable objet d’événement contenant des données typées, ce qui rend le système plus fiable et plus facile à utiliser.

**v2 :**

```go
runtime.EventsOn(ctx, "event-name", func(data ...interface{}) {
    // Handle event
})

runtime.EventsEmit(ctx, "event-name", data)
```

**v3 :**

```go
app.Event.On("event-name", func(e *application.CustomEvent) {
    data := e.Data
    // Handle event
})

app.Event.Emit("event-name", data)
```

**Pourquoi cette approche est préférable :**

- **Sécurité des types** : les événements utilisent de véritables objets d’événement au lieu de `...interface{}`
- **Débogage facilité** : les objets d’événement contiennent des métadonnées telles que le nom de l’événement, ce qui facilite le débogage
- **API plus claire** : `app.Event.On()` et `app.Event.Emit()` sont plus intuitifs que les fonctions du runtime
- **Aucun contexte nécessaire** : les événements fonctionnent directement sur l’objet d’application, sans avoir à faire transiter le contexte
- **Gestionnaires simplifiés** : les gestionnaires d’événements ont une signature claire au lieu de paramètres variadiques

### Fenêtres

La v2 ne prenait en charge qu’une seule fenêtre par application. La fenêtre était créée au démarrage et toutes les opérations la concernant étaient effectuées au moyen de fonctions du runtime qui ciblaient implicitement cette unique fenêtre.

La v3 introduit la prise en charge native de plusieurs fenêtres en tant que fonctionnalité fondamentale. Chaque fenêtre est un objet à part entière, doté de ses propres méthodes et de son propre cycle de vie. Vous pouvez créer, gérer et détruire dynamiquement plusieurs fenêtres pendant toute la durée de vie de votre application.

**v2 :**

```go
// Single window only
runtime.WindowSetSize(ctx, 800, 600)
```

**v3 :**

```go
// Multiple windows supported
window1 := app.Window.New()
window1.SetSize(800, 600)

window2 := app.Window.New()
window2.SetSize(1024, 768)
```

**Pourquoi cette approche est préférable :**

- **Applications multifenêtres** : créez des applications comportant plusieurs fenêtres indépendantes (tableaux de bord, préférences, outils, etc.)
- **Références explicites aux fenêtres** : chaque fenêtre est un objet que vous pouvez stocker et manipuler directement
- **Création dynamique de fenêtres** : créez et détruisez des fenêtres à tout moment pendant l’exécution
- **État indépendant de chaque fenêtre** : chaque fenêtre possède ses propres événements, propriétés et cycle de vie
- **Meilleure architecture** : la gestion des fenêtres repose sur des objets plutôt que sur le contexte

## Étapes de migration

### Étape 1 : mettez à jour les dépendances

**go.mod :**

```go
module myapp

go 1.25.0

require (
    github.com/wailsapp/wails/v3 v3.0.0-beta.0
)
```

**Mettez à jour :**

```bash
go get github.com/wailsapp/wails/v3@latest
go mod tidy
```

### Étape 2 : mettez à jour main.go

**v2 :**

```go
package main

import (
    "embed"
    "github.com/wailsapp/wails/v2/pkg/options"
    "github.com/wailsapp/wails/v2/pkg/options/assetserver"
    "github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
    app := NewApp()

    err := wails.Run(&options.App{
        Title:  "My App",
        Width:  1024,
        Height: 768,
        AssetServer: &assetserver.Options{
            Assets: assets,
        },
        Bind: []interface{}{
            app,
        },
        Windows: &windows.Options{
            WebviewIsTransparent: false,
        },
    })

    if err != nil {
        println("Error:", err.Error())
    }
}
```

**v3 :**

```go
package main

import (
    "embed"
    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
    app := application.New(application.Options{
        Name: "My App",
        Services: []application.Service{
            application.NewService(&MyService{}),
        },
        Assets: application.AssetOptions{
            Handler: application.AssetFileServerFS(assets),
        },
    })

    app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:  "My App",
        Width:  1024,
        Height: 768,
    })

    err := app.Run()
    if err != nil {
        panic(err)
    }
}
```

### Étape 3 : convertissez la structure App en service

**v2 :**

```go
type App struct {
    ctx context.Context
}

func NewApp() *App {
    return &App{}
}

func (a *App) startup(ctx context.Context) {
    a.ctx = ctx
    // Initialisation
}

func (a *App) Greet(name string) string {
    return "Hello " + name
}
```

**v3 :**

```go
type MyService struct {
    app *application.App
}

func NewMyService(app *application.App) *MyService {
    return &MyService{app: app}
}

func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    // Initialisation
    return nil
}

func (s *MyService) Greet(name string) string {
    return "Hello " + name
}

// Register after app creation
app := application.New(application.Options{})
app.RegisterService(application.NewService(NewMyService(app)))
```

### Étape 4 : mettez à jour les appels au runtime

**v2 :**

```go
func (a *App) DoSomething() {
    runtime.WindowSetTitle(a.ctx, "New Title")
    runtime.EventsEmit(a.ctx, "update", data)
    runtime.LogInfo(a.ctx, "Message")
}
```

**v3 :**

```go
func (s *MyService) DoSomething() {
    window := s.app.Window.Current()
    window.SetTitle("New Title")
    
    s.app.Event.Emit("update", data)
    
    s.app.Logger.Info("Message")
}
```

### Étape 5 : mettez à jour le frontend

**Générez les nouveaux bindings :**

```bash
wails3 generate bindings
```

**Mettez à jour les imports :**

```javascript
// v2
import { Greet } from '../wailsjs/go/main/App'

// v3
import { Greet } from './bindings/changeme/myservice'
```

**Mettez à jour la gestion des événements :**

```javascript
// v2
import { EventsOn, EventsEmit } from '../wailsjs/runtime/runtime'

EventsOn("update", (data) => {
    console.log(data)
})

EventsEmit("action", data)

// v3
import { Events } from '@wailsio/runtime'

Events.On("update", (data) => {
    console.log(data)
})

Events.Emit("action", data)
```

### Étape 6 : mettez à jour la configuration

**v2 (wails.json) :**

```json
{
  "name": "myapp",
  "outputfilename": "myapp",
  "frontend:install": "npm install",
  "frontend:build": "npm run build",
  "frontend:dev:watcher": "npm run dev",
  "frontend:dev:serverUrl": "auto"
}
```

**v3 (wails.json) :**

```json
{
  "name": "myapp",
  "frontend": {
    "dir": "./frontend",
    "install": "npm install",
    "build": "npm run build",
    "dev": "npm run dev",
    "devServerUrl": "http://localhost:5173"
  }
}
```

## Correspondance des fonctionnalités

### Boîtes de dialogue

**v2 :**

```go
selection, err := runtime.OpenFileDialog(ctx, runtime.OpenDialogOptions{
    Title: "Select File",
})
```

**v3 :**

```go
selection, err := app.Dialog.OpenFileWithOptions(&application.OpenFileDialogOptions{
    Title: "Select File",
}).PromptForSingleSelection()
```

### Menus

**v2 :**

```go
menu := menu.NewMenu()
menu.Append(menu.Text("File", nil, []*menu.MenuItem{
    menu.Text("Quit", nil, func(_ *menu.CallbackData) {
        runtime.Quit(ctx)
    }),
}))
```

**v3 :**

```go
menu := app.NewMenu()
fileMenu := menu.AddSubmenu("File")
fileMenu.Add("Quit").OnClick(func(ctx *application.Context) {
    app.Quit()
})
```

### Zone de notification

**v2 :**

```go
// Not available in v2
```

**v3 :**

```go
systray := app.SystemTray.New()
systray.SetIcon(iconBytes)
systray.SetLabel("My App")

menu := app.NewMenu()
menu.Add("Show").OnClick(showWindow)
menu.Add("Quit").OnClick(app.Quit)
systray.SetMenu(menu)
```

## Problèmes courants

### Problème : bindings introuvables

**Problème :** erreurs d’importation après la migration

**Solution :**

```bash
# Regenerate bindings
wails3 generate bindings

# Check output directory
ls frontend/bindings
```

### Problème : erreurs de contexte

**Problème :** `ctx` n’est pas disponible

**Solution :**

Stockez plutôt une référence à l’application :

```go
type MyService struct {
    app *application.App
}

func NewMyService(app *application.App) *MyService {
    return &MyService{app: app}
}
```

### Problème : les méthodes de fenêtre ne fonctionnent pas

**Problème :** `runtime.WindowSetTitle()` n’existe pas

**Solution :**

Utilisez directement les méthodes de la fenêtre :

```go
window := s.app.Window.Current()
window.SetTitle("New Title")
```

### Problème : les événements ne se déclenchent pas

**Problème :** les événements sont enregistrés, mais ne sont pas reçus

**Solution :**

Vérifiez que les noms des événements correspondent exactement :

```go
// Go
app.Event.Emit("my-event", data)

// JavaScript
OnEvent("my-event", handler)  // Must match exactly
```

## Test de la migration

### Liste de contrôle

- [ ] L’application démarre sans erreur
- [ ] Toutes les liaisons fonctionnent
- [ ] Les événements sont envoyés et reçus
- [ ] Les fenêtres s’ouvrent et se ferment correctement
- [ ] Les menus fonctionnent (le cas échéant)
- [ ] Les boîtes de dialogue fonctionnent (le cas échéant)
- [ ] La zone de notification fonctionne (le cas échéant)
- [ ] Le processus de compilation fonctionne
- [ ] La compilation de production fonctionne

### Commandes de test

```bash
# Development
wails3 dev

# Build
wails3 build

# Generate bindings
wails3 generate bindings
```

## Avantages de la v3

### Performances

- **Démarrage plus rapide** – Initialisation optimisée
- **Consommation de mémoire réduite** – Utilisation efficace des ressources
- **Meilleur pont** – Surcharge de &lt;1 ms par appel

### Fonctionnalités

- **Multifenêtre** – Prise en charge native
- **Zone de notification** – Intégrée
- **Événements améliorés** – API typée et plus simple
- **Services** – Meilleure organisation du code

### Expérience de développement

- **Sûreté du typage** – Prise en charge complète de TypeScript
- **Erreurs améliorées** – Messages d’erreur clairs
- **Rechargement à chaud** – Développement plus rapide
- **Documentation améliorée** – Guides complets

## Obtenir de l’aide

### Ressources

- [Documentation](/quick-start/why-wails/)
- [Communauté Discord](https://discord.gg/JDdSxwjhGf)
- [Tickets GitHub](https://github.com/wailsapp/wails/issues)
- [Exemples](https://github.com/wailsapp/wails/tree/master/v3/examples)

### Questions fréquentes

**Q : Puis-je exécuter la v2 et la v3 côte à côte ?** R : Oui, elles utilisent des chemins d’importation différents.

**Q : La v3 est-elle prête pour la production ?** R : La v3 est un logiciel bêta doté d’une API de bureau stable. Des applications l’utilisent en production, mais testez-la minutieusement avant le déploiement. La v2 reste la version stable actuelle.

**Q : La v2 sera-t-elle maintenue ?** R : Oui, la v2 recevra les mises à jour critiques.

**Q : Combien de temps prend la migration ?** R : 1-4 heures pour les applications classiques.

## Étapes suivantes

@cards{cols="2"}
🚀 Démarrage rapide
Commencez à utiliser Wails v3.

[En savoir plus →](/quick-start/installation/)

---
★ Concepts fondamentaux
Comprenez l’architecture de la v3.

[En savoir plus →](/concepts/architecture/)

---
◆ Liaisons
Découvrez le nouveau système de liaisons.

[En savoir plus →](/features/bindings/methods/)

---
📖 Exemples
Consultez des exemples complets pour la v3.

[Voir les exemples →](https://github.com/wailsapp/wails/tree/master/v3/examples)

@end

---

**Des questions ?** Posez-les sur [Discord](https://discord.gg/JDdSxwjhGf) ou [ouvrez un ticket](https://github.com/wailsapp/wails/issues).
