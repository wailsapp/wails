---
title: "Votre première application"
description: "Créez une application Wails fonctionnelle en 10 minutes"
slug: "quick-start/first-app"
sourcePath: "quick-start/first-app.md"
---

Nous allons créer une application simple de salutations qui illustre les concepts fondamentaux de Wails :

- Un backend Go qui gère la logique
- Un frontend qui appelle des fonctions Go
- Des liaisons avec sûreté des types
- Le rechargement à chaud pendant le développement

**Durée nécessaire :** 10 minutes

@note{type="tip" title="Conseil de performances pour les utilisateurs de Windows 11"}
Envisagez d’utiliser un [Dev Drive](https://learn.microsoft.com/en-us/windows/dev-drive/) pour stocker vos projets. Les Dev Drives sont optimisés pour les charges de travail de développement et peuvent améliorer considérablement les temps de compilation et les vitesses d’accès au disque, jusqu’à 30 % par rapport aux lecteurs NTFS classiques.

@end

## Créez votre projet

@steps
### Générez le projet
```bash
wails3 init -n myapp
cd myapp
```

Cette commande crée un projet avec le modèle Vanilla + Vite par défaut (HTML/CSS/TypeScript avec l’outil de regroupement Vite).

@note{type="tip" title="Autres modèles"}
Essayez `-t react`, `-t vue` ou `-t svelte` pour utiliser le framework de votre choix. Par défaut, ces modèles utilisent TypeScript ; pour utiliser JavaScript sans TypeScript, choisissez `-t vanilla-js` ou `-t react-js`. Exécutez `wails3 init -l` pour afficher tous les modèles disponibles, ou [utilisez votre propre framework frontend](/guides/dev/frontend-frameworks/).

@end

### Comprenez la structure du projet
```
myapp/
├── main.go              # Application entry point
├── greetservice.go      # Greet service
├── frontend/            # Your UI code
│   ├── index.html       # HTML entry point
│   ├── src/
│   │   └── main.ts      # Frontend TypeScript
│   ├── public/
│   │   └── style.css    # Styles
│   ├── package.json     # Frontend dependencies
│   ├── tsconfig.json    # TypeScript configuration
│   └── vite.config.ts   # Vite bundler config
├── build/               # Build configuration
└── Taskfile.yml         # Build tasks
```

### Exécutez l’application
```bash
wails3 dev
```

@note{type="info" title="Première exécution"}
La première exécution peut prendre plus de temps que prévu, car elle installe les dépendances du frontend, génère les liaisons, etc. Les exécutions suivantes sont beaucoup plus rapides.

@end

L’application s’ouvre et affiche une interface de salutations. Saisissez votre nom et cliquez sur « Greet » : le backend Go traite votre saisie et renvoie une salutation.

@end

## Fonctionnement

Examinons le code qui assure ce fonctionnement.

### Le backend Go

Ouvrez `greetservice.go` :

```go {title="greetservice.go"}
package main

import (
	"fmt"
)

type GreetService struct{}

func (g *GreetService) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
```

**Concepts clés :**

1. **Service** – Une structure Go avec des méthodes exportées
2. **Méthode exportée** – `Greet` commence par une majuscule, ce qui la rend accessible au frontend
3. **Logique simple** – Reçoit un nom et renvoie une salutation
4. **Sûreté des types** – Les types d’entrée et de sortie sont définis

@note{type="tip" title="Comprendre les services et les liaisons"}
Les **services** sont des modules Go autonomes qui exposent des fonctionnalités à votre frontend. Il s’agit simplement de structures Go classiques dotées de méthodes exportées, que vous enregistrez dans le champ `Services` de la configuration de votre application.

Les **liaisons** constituent le SDK TypeScript/JavaScript généré automatiquement qui permet à votre frontend d’appeler ces services. Lorsque vous exécutez `wails3 dev` ou `wails3 build`, Wails analyse les services enregistrés et génère des liaisons avec sûreté des types dans `frontend/bindings/`.

Considérez les services comme l’API de votre backend et les liaisons comme la bibliothèque cliente qui communique avec elle.

@end

### Enregistrement du service

Ouvrez `main.go` et recherchez l’enregistrement du service :

```go {title="main.go" highlight="4-6"}
err := application.New(application.Options{
    Name: "myapp",
    Services: []application.Service{
        application.NewService(&GreetService{}),
    },
    // ... other options
})
```

Cela enregistre votre `GreetService` auprès de Wails et rend toutes ses méthodes exportées accessibles au frontend.

### Le frontend

Ouvrez `frontend/src/main.js` :

```javascript {title="frontend/src/main.js"}
import {GreetService} from "../bindings/changeme";

window.greet = async () => {
    const nameElement = document.getElementById('name');
    const resultElement = document.getElementById('result');

    const name = nameElement.value;
    if (!name) {
        return;
    }

    try {
        const result = await GreetService.Greet(name);
        resultElement.innerText = result;
    } catch (err) {
        console.error(err);
    }
};
```

**Concepts clés :**

1. **Liaisons générées automatiquement** – `GreetService` est importé depuis le code généré
2. **Appels avec sûreté des types** – Les noms et signatures des méthodes correspondent à votre code Go
3. **Asynchrones par défaut** – Tous les appels Go renvoient des promesses
4. **Gestion des erreurs** – Les erreurs provenant de Go sont interceptées dans un bloc try/catch

@note{type="info" title="Où se trouvent les liaisons ?"}
Les liaisons générées se trouvent dans `frontend/bindings/`. Elles sont créées automatiquement lorsque vous exécutez `wails3 dev` ou `wails3 build`.

**Ne modifiez jamais ces fichiers manuellement** : ils sont générés à nouveau à chaque compilation.

@end

## Personnalisez votre application

Ajoutons une fonctionnalité pour comprendre le flux de travail.

### Ajoutez une fonctionnalité « Saluer plusieurs personnes »

@steps
### Ajoutez la méthode à GreetService
Ajoutez ceci à `greetservice.go` :

```go {title="greetservice.go"}
func (g *GreetService) GreetMany(names []string) []string {
    greetings := make([]string, len(names))
    for i, name := range names {
        greetings[i] = fmt.Sprintf("Hello %s!", name)
    }
    return greetings
}
```

### L’application sera recompilée automatiquement
Enregistrez le fichier : `wails3 dev` recompilera automatiquement votre code Go et redémarrera l’application.

@note{type="info" title="Recompilation automatique"}
Les modifications du code Go déclenchent automatiquement une recompilation et un redémarrage. Les modifications du frontend sont rechargées à chaud sans redémarrage.

@end

### Utilisez-la dans le frontend
Ajoutez ceci à `frontend/src/main.js` :

```javascript {title="frontend/src/main.js"}
window.greetMany = async () => {
    const names = ['Alice', 'Bob', 'Charlie'];
    const greetings = await GreetService.GreetMany(names);
    console.log(greetings);
};
```

Ouvrez la console du navigateur et appelez `greetMany()` : le tableau de salutations s’affichera.

@end

## Compilez pour la production

Lorsque vous êtes prêt à distribuer votre application :

```bash
wails3 build
```

**Résultat de cette opération :**

- Compile le code Go avec des optimisations
- Compile le frontend pour la production (version minifiée)
- Crée un exécutable natif dans `bin/`

@tabs{sync-key="os"}
[Windows]
**Sortie :** `bin/myapp.exe`

Double-cliquez pour lancer l’application. Aucune dépendance n’est nécessaire (WebView2 fait partie de Windows).

[macOS]
**Sortie :** `bin/myapp.app`

Faites glisser l’application dans le dossier Applications ou double-cliquez dessus pour la lancer.

[Linux]
**Sortie :** `bin/myapp`

Lancez l’application avec `./bin/myapp` ou créez un fichier `.desktop` pour votre lanceur.

@end

@note{type="tip" title="Compilations multiplateformes"}
Vous souhaitez compiler pour d’autres plateformes ? Consultez [Compilations multiplateformes →](/guides/build/cross-platform/)

@end

## Ce que nous avons appris

**Structure du projet**

- `main.go` pour le backend Go
- `frontend/` pour le code de l’interface utilisateur
- `Taskfile.yml` pour les tâches de compilation

**Services**

- Créez des structures Go avec des méthodes exportées
- Enregistrez-les avec `application.NewService()`
- Les méthodes sont automatiquement disponibles dans le frontend

**Liaisons**

- Définitions TypeScript générées automatiquement
- Appels de fonctions avec vérification des types
- Asynchrones par défaut (promesses)

**Flux de développement**

- `wails3 dev` pour le rechargement à chaud
- Les modifications du code Go déclenchent automatiquement une nouvelle compilation et un redémarrage
- Les modifications du frontend sont rechargées à chaud instantanément

---

**Des questions ?** Rejoignez [Discord](https://discord.gg/JDdSxwjhGf) et posez-les à la communauté.
