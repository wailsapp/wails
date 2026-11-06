---
title: "Fonctionnement de Wails"
description: "Comprendre l’architecture de Wails et la manière dont elle offre des performances natives"
slug: "concepts/architecture"
sourcePath: "concepts/architecture.md"
---

Wails est un framework permettant de créer des applications de bureau avec **Go pour le backend** et des **technologies web pour le frontend**. Mais contrairement à Electron, Wails n’intègre pas de navigateur : il utilise la **WebView native du système d’exploitation**.

```d2
direction: left

Wails App: {
  shape: sequence_diagram
  label: Application Wails

  frontend: Frontend
  backend: Backend Go
  os: Système d’exploitation

  Initialisation: Initialisation {
    shape: sequence_diagram
    backend."Serves Static Web App": Sert l’application web statique
    backend -> frontend: HTML / JS / CSS
    frontend."Render Site via OS-native WebView": Affiche le site via la WebView native du système
  }
  Regular Communication: Communication normale {
    shape: sequence_diagram
    frontend."Make API-style call": Effectue un appel de type API
    frontend -> backend.a: JSON
    backend.a."Service processes request": Le service traite la requête
    backend.a -> os: Appelle les API système
    backend.a."Generate Response": Génère la réponse
    backend.a -> frontend: JSON
    frontend."Process response": Traite la réponse
  }
  backend.a.label: a
}
```

**Principales différences avec Electron :**

| Aspect | Wails | Electron |
| --- | --- | --- |
| **Navigateur** | WebView fournie par le système d’exploitation | Chromium intégré (~100 Mo) |
| **Backend** | Go (compilé) | Node.js (interprété) |
| **Communication** | Pont en mémoire | IPC (interprocessus) |
| **Taille du paquet** | ~15 Mo | ~150 Mo |
| **Mémoire** | ~10 Mo | ~100 Mo+ |
| **Démarrage** | &lt;0.5 s | 2-3 s |

## Composants principaux

### 1. WebView native

Wails utilise le moteur de rendu web intégré au système d’exploitation :

@tabs{sync-key="platform"}
[Windows]
**WebView2** (Microsoft Edge WebView2)

- Basé sur Chromium (comme le navigateur Edge)
- Préinstallé sur Windows 10/11
- Mises à jour automatiques via Windows Update
- Prise en charge complète des standards web modernes

[macOS]
**WebKit** (moteur de rendu de Safari)

- Intégré à macOS
- Même moteur que le navigateur Safari
- Excellentes performances et autonomie
- Prise en charge complète des standards web modernes

[Linux]
**WebKitGTK** (portage GTK de WebKit)

- Installé via le gestionnaire de paquets
- Même moteur que GNOME Web (Epiphany)
- Bonne prise en charge des standards
- Léger et performant

@end

**Pourquoi est-ce important ?**

- **Aucun navigateur intégré** → Application moins volumineuse
- **Natif au système d’exploitation** → Meilleures intégration et performances
- **Mises à jour automatiques** → Correctifs de sécurité fournis par les mises à jour du système d’exploitation
- **Rendu familier** → Identique à celui du navigateur système

### 2. Le pont Wails

Le pont est au cœur de Wails : il permet une **communication directe** entre Go et JavaScript.

```d2
direction: down

Frontend: Frontend (JavaScript) {
  shape: rectangle
  style.fill: "#8B5CF6"
}

Bridge: Pont Wails {
  Encoder: Encodeur JSON {
    shape: rectangle
  }

  Router: Routeur de méthodes {
    shape: diamond
    style.fill: "#10B981"
  }

  Decoder: Décodeur JSON {
    shape: rectangle
  }
}

Backend: Backend (Go) {
  Services: Services enregistrés {
    shape: rectangle
    style.fill: "#00ADD8"
  }
}

Frontend -> Bridge.Encoder: "1. Appeler la méthode Go\nGreet('Alice')"
Bridge.Encoder -> Bridge.Router: "2. Encoder en JSON\n{method: 'Greet', args: ['Alice']}"
Bridge.Router -> Backend.Services: "3. Acheminer vers le service\nGreetService.Greet('Alice')"
Backend.Services -> Bridge.Decoder: "4. Renvoyer le résultat\n'Hello, Alice!'"
Bridge.Decoder -> Frontend: "5. Décoder en JS\nLa promesse est résolue"
```

**Fonctionnement :**

1. **Le frontend appelle une méthode Go** (via une liaison générée automatiquement)
2. **Le pont encode l’appel** en JSON (nom de la méthode + arguments)
3. **Le routeur trouve la méthode Go** dans les services enregistrés
4. **La méthode Go s’exécute** et renvoie une valeur
5. **Le pont décode le résultat** et le renvoie au frontend
6. **La promesse est résolue** dans JavaScript avec le résultat

**Caractéristiques de performance :**

- **En mémoire** : aucun surcoût réseau, aucun HTTP
- **Sans copie** lorsque cela est possible (pour les données volumineuses)
- **Asynchrone par défaut** : non bloquant des deux côtés
- **À typage sûr** : définitions TypeScript générées automatiquement

### 3. Système de services

Les services constituent la méthode recommandée pour exposer les fonctionnalités Go au frontend.

```go
// Define a service (just a regular Go struct)
type GreetService struct {
    prefix string
}

// Methods with exported names are automatically available
func (g *GreetService) Greet(name string) string {
    return g.prefix + name + "!"
}

func (g *GreetService) GetTime() time.Time {
    return time.Now()
}

// Register the service
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&GreetService{prefix: "Hello, "}),
    },
})
```

**Découverte des services :**

- Wails **analyse votre structure** au démarrage
- Les **méthodes exportées** peuvent être appelées depuis le frontend
- Les **informations de type** sont extraites pour les bindings TypeScript
- La **gestion des erreurs** est automatique (erreurs Go → exceptions JS)

**Binding TypeScript généré :**

```typescript
// Auto-generated in frontend/bindings/GreetService.ts
export function Greet(name: string): Promise<string>
export function GetTime(): Promise<Date>
```

**Pourquoi utiliser des services ?**

- **Typage sûr** : prise en charge complète de TypeScript
- **Découverte automatique** : aucun enregistrement manuel des méthodes
- **Organisation** : regroupez les fonctionnalités associées
- **Testabilité** : les services sont de simples structures Go

[En savoir plus sur les services →](/features/bindings/services/)

### 4. Système d’événements

Les événements permettent une **communication par publication-abonnement** entre les composants.

```d2
direction: left

Wails Event System: Système d’événements Wails {
  shape: sequence_diagram

  window1: Fenêtre 1
  window2: Fenêtre 2
  backend: Backend Go

  Event Driver: Gestionnaire d’événements {
    shape: sequence_diagram
    window1."Subscribe to 'data-updated' events": "S’abonner aux événements 'data-updated'"
    window2."Subscribe to 'data-updated' events": "S’abonner aux événements 'data-updated'"
    backend.a."App Emit('data-updated', data)": "L’application émet Emit('data-updated', data)"
    backend.a -> window1.a: Bus d’événements JSON
    backend.a -> window2: Bus d’événements JSON
    window1.a."Subscriber processes On('data-updated', handler)": "L’abonné traite On('data-updated', handler)"
    window2."Subscriber processes On('data-updated', handler)": "L’abonné traite On('data-updated', handler)"
  }
  backend.a.label: a
  window1.a.label: a
}
```

**Cas d’utilisation :**

- **Communication entre fenêtres** : une fenêtre en informe d’autres
- **Tâches en arrière-plan** : le service Go informe l’interface utilisateur de la progression
- **Synchronisation de l’état** : maintenez plusieurs fenêtres synchronisées
- **Couplage faible** : les composants n’ont pas besoin de références directes

**Exemple :**

```go
// Go: Emit an event
app.Event.Emit("user-logged-in", user)
```

```javascript
// JavaScript: Listen for event
import { Events } from '@wailsio/runtime'

Events.On('user-logged-in', (user) => {
    console.log('User logged in:', user)
})
```

[En savoir plus sur les événements →](/features/events/system/)

## Cycle de vie de l’application

Comprendre le cycle de vie vous aide à déterminer quand initialiser et libérer les ressources.

```d2
direction: down

Start: Démarrage de l’application {
  shape: oval
  style.fill: "#10B981"
}

Init: Initialisation {
  Create: Créer l’application {
    shape: rectangle
  }

  Register: Enregistrer les services {
    shape: rectangle
  }

  Setup: Configurer les fenêtres/menus {
    shape: rectangle
  }
}

Run: Boucle d’événements {
  Events: Traiter les événements {
    shape: rectangle
  }

  Messages: Gérer les messages {
    shape: rectangle
  }

  Render: Mettre à jour l’interface utilisateur {
    shape: rectangle
  }
}

Shutdown: Arrêt {
  Cleanup: Libérer les ressources {
    shape: rectangle
  }

  Save: Enregistrer l’état {
    shape: rectangle
  }
}

End: Fin de l’application {
  shape: oval
  style.fill: "#EF4444"
}

Start -> Init.Create
Init.Create -> Init.Register
Init.Register -> Init.Setup
Init.Setup -> Run.Events
Run.Events -> Run.Messages
Run.Messages -> Run.Render
Run.Render -> Run.Events: Boucle
Run.Events -> Shutdown.Cleanup: Signal de fermeture
Shutdown.Cleanup -> Shutdown.Save
Shutdown.Save -> End
```

**Hooks du cycle de vie :**

```go
app := application.New(application.Options{
    Name: "My App",

    // Cleanly intercept quit requests (e.g. unsaved changes).
    ShouldQuit: func() bool { return true },

    // Called when the app is confirmed to be quitting — save state, close connections, etc.
    OnShutdown: func() {},
})
```

`application.Options` ne comporte aucun champ `OnStartup`. Les opérations de démarrage doivent être placées dans le `ServiceStartup(ctx, options)` d’un service, dans un rappel enregistré via `app.Event.OnApplicationEvent(events.Common.ApplicationStarted, ...)` ou simplement avant `app.Run()`.

[En savoir plus sur le cycle de vie →](/concepts/lifecycle/)

## Processus de build

Voici comment Wails construit votre application :

```d2
direction: down

Source: Code source {
  Go: "Code Go\n(main.go, services)" {
    shape: rectangle
    style.fill: "#00ADD8"
  }

  Frontend: "Code frontend\n(HTML/CSS/JS)" {
    shape: rectangle
    style.fill: "#8B5CF6"
  }
}

Build: Processus de build {
  AnalyseGo: Analyser le code Go {
    shape: rectangle
  }

  GenerateBindings: Générer les bindings {
    shape: rectangle
  }

  BuildFrontend: Construire le frontend {
    shape: rectangle
  }

  CompileGo: Compiler le code Go {
    shape: rectangle
  }

  Embed: Incorporer les ressources {
    shape: rectangle
  }
}

Output: Sortie {
  Binary: "Binaire natif\n(myapp.exe/.app)" {
    shape: rectangle
    style.fill: "#10B981"
  }
}

Source.Go -> Build.AnalyseGo
Build.AnalyseGo -> Build.GenerateBindings: Extraire les types
Build.GenerateBindings -> Source.Frontend: Bindings TypeScript
Source.Frontend -> Build.BuildFrontend: Compiler (Vite/webpack)
Build.BuildFrontend -> Build.Embed: Ressources regroupées
Source.Go -> Build.CompileGo
Build.CompileGo -> Build.Embed
Build.Embed -> Output.Binary
```

**Étapes du build :**

1. **Analyser le code Go**
  - Analyser les services pour trouver les méthodes exportées
  - Extraire les types des paramètres et des valeurs de retour
  - Générer les signatures de méthodes


2. **Générer les bindings TypeScript**
  - Créer les fichiers `.ts` de chaque service
  - Inclure les définitions de types complètes
  - Ajouter des commentaires JSDoc


3. **Construire le frontend**
  - Exécuter votre bundler (Vite, webpack, etc.)
  - Minifier et optimiser
  - Écrire la sortie dans `frontend/dist/`


4. **Compiler le code Go**
  - Compiler avec les optimisations (`-ldflags="-s -w"`)
  - Inclure les métadonnées du build
  - Effectuer la compilation propre à la plateforme


5. **Incorporer les ressources**
  - Incorporer les fichiers du frontend au binaire Go
  - Compresser les ressources
  - Créer un exécutable unique


**Résultat :** un exécutable natif unique contenant tous les éléments incorporés.

[En savoir plus sur le build →](/guides/build/building/)

## Développement et production

Wails se comporte différemment en développement et en production :

@tabs{sync-key="mode"}
[Développement (wails3 dev)]
**Caractéristiques :**

- **Rechargement à chaud** : les modifications du frontend sont rechargées instantanément
- **Mappages de sources** : débogage à partir du code source d’origine
- **DevTools** : les outils de développement du navigateur sont disponibles
- **Journalisation** : la journalisation détaillée est activée
- **Frontend externe** : servi par le serveur de développement (Vite)

**Fonctionnement :**

```d2
direction: right

WailsApp: Application Wails {
  shape: rectangle
  style.fill: "#00ADD8"
}

DevServer: "Serveur de développement Vite\n(localhost:5173)" {
  shape: rectangle
  style.fill: "#8B5CF6"
}

WebView: WebView {
  shape: rectangle
  style.fill: "#6B7280"
}

WailsApp -> DevServer: Transmettre les requêtes via le proxy
DevServer -> WebView: Servir avec HMR
WebView -> WailsApp: Appeler les méthodes Go
```

**Avantages :**

- Retour instantané sur les modifications
- Fonctionnalités de débogage complètes
- Itérations plus rapides

[Production (wails3 build)]
**Caractéristiques :**

- **Ressources intégrées** : le frontend est intégré au binaire
- **Optimisé** : minifié et compressé
- **Aucun DevTools** : désactivé par défaut
- **Journalisation minimale** : erreurs uniquement
- **Fichier unique** : tout est regroupé dans un seul exécutable

**Fonctionnement :**

```d2
direction: right

Binary: "Binaire unique\n(myapp.exe)" {
  GoCode: Code Go compilé {
    shape: rectangle
    style.fill: "#00ADD8"
  }

  Assets: "Ressources intégrées\n(HTML/CSS/JS)" {
    shape: rectangle
    style.fill: "#8B5CF6"
  }
}

WebView: WebView {
  shape: rectangle
  style.fill: "#6B7280"
}

Binary.Assets -> WebView: Servir depuis la mémoire
WebView -> Binary.GoCode: Appeler les méthodes Go
```

**Avantages :**

- Distribution sous la forme d’un fichier unique
- Taille réduite (minification)
- Meilleures performances
- Aucune dépendance externe

@end

## Modèle de mémoire

Comprendre l’utilisation de la mémoire vous aide à créer des applications efficaces.

**Zones de mémoire :**

1. **Tas Go**
  - Vos services et l’état de l’application
  - Géré par le ramasse-miettes de Go
  - Généralement 5-10 Mo pour les applications simples


2. **Mémoire de la WebView**
  - DOM, tas JavaScript et CSS
  - Gérée par le moteur de la WebView
  - Généralement 10-20 Mo pour les applications simples


3. **Mémoire du pont**
  - Tampons de messages pour la communication
  - Surcharge minimale (<1 Mo)
  - Aucune copie pour les grandes quantités de données lorsque cela est possible


**Conseils d’optimisation :**

- **Évitez de transférer de grandes quantités de données** : transmettez des identifiants et récupérez les détails à la demande
- **Utilisez des événements pour les mises à jour** : n’interrogez pas périodiquement le backend depuis le frontend
- **Diffusez les fichiers volumineux en continu** : ne les chargez pas entièrement en mémoire
- **Supprimez les écouteurs devenus inutiles** : retirez les écouteurs d’événements lorsque vous n’en avez plus besoin

[En savoir plus sur les performances →](/guides/performance/)

## Modèle de sécurité

Wails fournit une architecture sécurisée par défaut :

```d2
direction: down

Frontend: Frontend (non fiable) {
  shape: rectangle
  style.fill: "#EF4444"
}

Bridge: Pont Wails (validation) {
  shape: diamond
  style.fill: "#F59E0B"
}

Backend: Backend (fiable) {
  shape: rectangle
  style.fill: "#10B981"
}

Frontend -> Bridge: Appeler la méthode
Bridge -> Bridge: "Valider :\n- La méthode existe-t-elle ?\n- Les types sont-ils corrects ?\n- L’accès est-il autorisé ?"
Bridge -> Backend: Exécuter si valide
Backend -> Bridge: Renvoyer le résultat
Bridge -> Frontend: Envoyer la réponse
```

**Fonctionnalités de sécurité :**

1. **Liste blanche des méthodes**
  - Seules les méthodes exportées peuvent être appelées
  - Les méthodes privées sont inaccessibles
  - L’enregistrement explicite des services est obligatoire


2. **Validation des types**
  - Les arguments sont vérifiés par rapport aux types Go
  - Les types non valides sont rejetés
  - Empêche les attaques par injection


3. **Aucun eval()**
  - Le frontend ne peut pas exécuter de code Go arbitraire
  - Seules les méthodes prédéfinies peuvent être appelées
  - Aucune exécution dynamique de code


4. **Isolation des contextes**
  - Chaque fenêtre possède son propre contexte
  - Les services peuvent vérifier le contexte de l’appelant
  - Des autorisations par fenêtre sont possibles


**Bonnes pratiques :**

- **Validez les saisies utilisateur** dans Go (ne faites pas confiance au frontend)
- **Utilisez le contexte** pour l’authentification et l’autorisation
- **Assainissez les chemins de fichiers** avant toute opération sur les fichiers
- **Limitez la fréquence** des opérations coûteuses

[En savoir plus sur la sécurité →](/guides/security/)

## Étapes suivantes

**Cycle de vie de l’application** - Comprenez le démarrage, l’arrêt et les hooks du cycle de vie [En savoir plus →](/concepts/lifecycle/)

**Pont entre Go et le frontend** - Découvrez en détail le fonctionnement du pont [En savoir plus →](/concepts/bridge/)

**Système de build** - Comprenez comment Wails construit votre application [En savoir plus →](/concepts/build-system/)

**Commencez à développer** – Mettez en pratique ce que vous avez appris dans un tutoriel [Tutoriels →](/tutorials/03-notes-vanilla/)

---

**Des questions sur l’architecture ?** Posez-les sur [Discord](https://discord.gg/JDdSxwjhGf) ou consultez la [référence de l’API](/reference/overview/).
