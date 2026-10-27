---
title: "Référence de l’API"
description: "Consultez les types, méthodes, options et commandes CLI de Wails"
slug: "reference/overview"
sourcePath: "reference/overview.md"
---

## À propos de cette référence

Utilisez cette section pour consulter les types, méthodes, options et commandes de Wails pendant le développement d’une application. Si vous recherchez un apprentissage guidé, commencez par les [tutoriels](/tutorials/overview/). Pour suivre les étapes permettant d’accomplir une tâche, consultez les [guides pratiques](/quick-start/next-steps/#features). Pour découvrir les principes de conception, consultez les [concepts fondamentaux](/concepts/architecture/).

**Organisation :**

- [Application](/reference/application/) — API principales de l’application
- [Fenêtre](/reference/window/) — Création et gestion des fenêtres ; les [options de fenêtre](/features/windows/options/) répertorient les champs de configuration
- [Menu](/reference/menu/) — API des menus d’application et contextuels, ainsi que de la zone de notification ; la [référence des éléments de menu](/features/menus/reference/) décrit les types et propriétés des éléments
- [Événements](/reference/events/) — Système d’événements et événements intégrés
- [Boîtes de dialogue](/reference/dialogs/) — Boîtes de dialogue de fichiers et de messages
- [Runtime frontend](/reference/frontend-runtime/) — API du runtime frontend
- [Vue d’ensemble des commandes CLI](/reference/cli/) — Commandes regroupées par objectif ; les [options et exemples](/guides/cli/) détaillent chaque commande
- [Protocole du manifeste de mise à jour](/reference/update-manifest/) — Métadonnées de mise à jour et champs de vérification
- [API mobile](/guides/mobile/mobile-api/) — API mobiles expérimentales

## Conventions de l’API

@details{title="Conventions de l’API Go — Pour les développeurs qui découvrent Go"}
### Nommage

- <strong>Types</strong> : PascalCase (par exemple, `WebviewWindow`)
- <strong>Méthodes</strong> : PascalCase (par exemple, `SetTitle()`)
- <strong>Options</strong> : structures en PascalCase (par exemple, `WindowOptions`)
- <strong>Constantes</strong> : PascalCase (par exemple, `WindowStartStateMaximised`)

#### Gestion des erreurs

La plupart des méthodes susceptibles d’échouer renvoient `error` comme dernière valeur de retour. `app.Run()` bloque l’exécution jusqu’à la fermeture de l’application et renvoie toute erreur survenue au démarrage :

```go
if err := app.Run(); err != nil {
    log.Fatal(err)
}
```

La construction d’une fenêtre ne renvoie pas d’erreur : `app.Window.New()` renvoie directement `*WebviewWindow`.

#### Contexte

Les méthodes du cycle de vie des services reçoivent un `context.Context` :

```go
func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    // ctx is cancelled when the application is shutting down.
    return nil
}
```

Le contexte correspondant à la durée de vie de l’application est accessible via `app.Context()`. Il n’existe pas de `RunWithContext` : appelez `app.Run()`.

#### Modèle des options

La configuration utilise des structures d’options :

```go
app := application.New(application.Options{
    Name: "My App",
    Description: "A demo application",
    Services: []application.Service{
        application.NewService(&MyService{}),
    },
})
```

@end

### Conventions de l’API JavaScript

#### Nommage

- **Fonctions** : camelCase (par exemple, `setTitle()`)
- **Constantes** : SCREAMING<em>SNAKE</em>CASE (par exemple, `WINDOW_EVENT_FOCUS`)

#### Asynchrone par défaut

Tous les appels de méthodes Go renvoient des promesses :

```javascript
// Async/await (recommended)
const result = await MyService.DoSomething()

// Promise chain
MyService.DoSomething()
    .then(result => console.log(result))
    .catch(error => console.error(error))
```

#### Gestion des erreurs

Les erreurs Go deviennent des exceptions JavaScript :

```javascript
try {
    await MyService.MightFail()
} catch (error) {
    console.error('Go error:', error)
}
```

#### Sûreté des types

Les définitions TypeScript sont générées automatiquement :

```typescript
// Fully typed
import { Greet } from './bindings/GreetService'

const message: string = await Greet("World")
```

## Structure des paquets

```
github.com/wailsapp/wails/v3/pkg/
├── application/          # Core application package
│   ├── application.go    # App type
│   ├── webview_window.go # Window management
│   ├── menu.go           # Menu types
│   ├── event_manager.go  # Event system
│   └── dialogs.go        # Dialog APIs
├── events/               # Event constants
└── services/             # Built-in services
    ├── dock/             # macOS dock (includes badge support)
    ├── fileserver/       # File-server service
    ├── kvstore/          # Key/value store
    ├── log/              # Structured logging service
    ├── notifications/    # Notifications service
    └── sqlite/           # SQLite service
```

## Chemins d’importation

### Go

```go
import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)
```

### JavaScript

```javascript
// Auto-generated bindings
import { MyMethod } from './bindings/MyService'

// Runtime APIs
import { Events, Window } from '@wailsio/runtime'
```

## Référence des types

### Types courants

@tabs{sync-key="lang"}
[Go]
```go
// Application
type App struct { /* ... */ }
type Options struct { /* ... */ }

// Window
type WebviewWindow struct { /* ... */ } // implements the Window interface
type WebviewWindowOptions struct { /* ... */ }

// Menu
type Menu struct { /* ... */ }
type MenuItem struct { /* ... */ }

// Events — there is no generic Event type; events are typed by source.
type ApplicationEvent struct { /* ... */ }
type WindowEvent struct { /* ... */ }
type CustomEvent struct { /* ... */ }
type EventListener struct { /* ... */ }

// Dialogs
type OpenFileDialogOptions struct { /* ... */ }
type SaveFileDialogOptions struct { /* ... */ }
```

[TypeScript]
```typescript
// Window runtime
interface WindowOptions {
    title?: string
    width?: number
    height?: number
    // ...
}

// Events
type EventCallback = (data: any) => void

// Bindings (auto-generated)
export function MyMethod(arg: string): Promise<string>
```

@end

## Différences entre les plateformes

Le comportement de certaines API varie selon la plateforme :

| Fonctionnalité | Windows | macOS | Linux |
| --- | --- | --- | --- |
| **Menu de l’application** | Barre de menus de la fenêtre | Barre de menus globale | Barre de menus de la fenêtre |
| **Zone de notification** | Zone de notification | Barre de menus | Zone de notification |
| **Dock** | Non applicable | ✅ Disponible | Non applicable |
| **Boîtes de dialogue de fichiers** | Natives | Natives | Natives (GTK) |
| **Transparence** | ✅ Complète | Nécessite [`private_mac_apis`](/guides/build/private-macos-apis/#webview-transparency-and-background) | ⚠️ Limitée |

Le comportement propre à chaque plateforme est documenté dans chaque section de l’API.

## Gestion des versions

Wails v3 suit la gestion sémantique des versions :

- **Majeure** (v4.0.0) : changements incompatibles
- **Mineure** (v3.1.0) : nouvelles fonctionnalités rétrocompatibles
- **Corrective** (v3.0.1) : corrections de bogues rétrocompatibles

**Wails v3.0.0 :** stable

## Politique d’obsolescence

Lorsqu’une API devient obsolète :

1. **Signalement dans la documentation** au moyen d’un avis d’obsolescence
2. **Proposition d’une solution de remplacement** accompagnée d’un guide de migration
3. **Maintien pendant 1 version majeure** avant sa suppression
4. **Avertissements du compilateur** (lorsque cela est possible)

## Stabilité de l’API

### API stables ✅

Ces API sont stables et adaptées à une utilisation en production :

- API principales de l’application
- Gestion des fenêtres
- Système de menus
- Système d’événements
- Boîtes de dialogue de fichiers
- Liaisons de services

### API expérimentales ⚠️

Seules les API explicitement documentées comme expérimentales ne sont pas couvertes par le contrat de compatibilité stable pour les applications de bureau :

- Prise en charge d’Android
- Prise en charge d’iOS
- Autres API explicitement signalées comme expérimentales

Les API expérimentales sont signalées dans la documentation.

## Obtenir de l’aide

### Questions sur les API

1. **Consultez cette référence** — Documentation complète des API
2. **Consultez les exemples** — [Exemples sur GitHub](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples)
3. **Effectuez une recherche sur Discord** — [Serveur Discord](https://discord.gg/JDdSxwjhGf)
4. **Interrogez la communauté** — Canal Discord #help

### Signaler des problèmes d’API

Vous avez trouvé un bogue ou une incohérence ?

1. **Consultez les problèmes existants** — [Problèmes sur GitHub](https://github.com/wailsapp/wails/issues)
2. **Créez un rapport détaillé** — Incluez le code, l’erreur et la plateforme
3. **Fournissez une procédure de reproduction** — Un exemple minimal qui illustre le problème

## Documentation connexe

- [Tutoriels](/tutorials/overview/) — Apprenez en créant de véritables applications
- [Guides](/guides/architecture/) — Guides axés sur les tâches pour les scénarios courants
- [Fonctionnalités](/features/windows/basics/) — Documentation détaillée par fonctionnalité
- [Exemples](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples) — Exemples de code fonctionnels sur GitHub

---

**Parcourir les API :** utilisez la navigation à gauche pour explorer les différentes API.
