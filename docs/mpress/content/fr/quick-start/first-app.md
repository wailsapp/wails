---
title: "Votre première application"
description: "Créez une application Wails fonctionnelle en 10 minutes"
slug: "quick-start/first-app"
sourcePath: "quick-start/first-app.md"
---

Nous allons créer une application de salutation simple qui illustre les concepts fondamentaux de Wails :

- Backend Go gérant la logique
- Frontend appelant des fonctions Go
- Liaisons avec typage sûr
- Rechargement à chaud pendant le développement

**Temps nécessaire :** 10 minutes

## Créez votre projet

@steps
### Générez le projet
```bash
wails3 init -n myapp -t vanilla-js
cd myapp
```

Cette leçon sélectionne explicitement Vanilla + Vite en JavaScript afin que les fichiers sources ci-dessous correspondent à votre projet généré. Le modèle `vanilla` sans option supplémentaire utilise TypeScript.

@note{type="tip" title="Autres modèles"}
Essayez `-t react`, `-t vue` ou `-t svelte` selon le framework que vous préférez. Ces modèles utilisent TypeScript par défaut ; pour employer du JavaScript simple, utilisez `-t vanilla-js` ou `-t react-js`. Exécutez `wails3 init -l` pour afficher tous les modèles disponibles, ou [utilisez votre propre framework frontend](/guides/dev/frontend-frameworks/).

@end

### Comprenez la structure du projet
```
myapp/
├── main.go              # Application entry point
├── greetservice.go      # Greet service
├── frontend/            # Your UI code
│   ├── index.html       # HTML entry point
│   ├── src/
│   │   └── main.js      # Frontend JavaScript
│   ├── public/
│   │   └── style.css    # Styles
│   ├── package.json     # Frontend dependencies
│   ├── tsconfig.json    # Editor and JavaScript tooling configuration
│   └── vite.config.js   # Vite bundler config
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

L’application s’ouvre et affiche une interface de salutation. Saisissez votre nom et cliquez sur « Greet » : le backend Go traite votre saisie et renvoie une salutation.

@end

@note{type="tip" title="Conseil de performances pour les utilisateurs de Windows 11"}
Envisagez d’utiliser [Dev Drive](https://learn.microsoft.com/en-us/windows/dev-drive/) pour stocker vos projets. Les Dev Drives sont optimisés pour les charges de travail de développement et peuvent améliorer considérablement les temps de compilation et les vitesses d’accès au disque, jusqu’à 30 % par rapport aux lecteurs NTFS classiques.

@end

## Fonctionnement

Examinons le code qui assure ce fonctionnement.

### Le backend Go

Ouvrez `greetservice.go` :

```go {title="greetservice.go"}
package main

type GreetService struct{}

func (g *GreetService) Greet(name string) string {
	return "Hello " + name + "!"
}
```

**Concepts clés :**

1. **Service** — Une structure Go dotée de méthodes exportées
2. **Méthode exportée** — `Greet` commence par une majuscule, ce qui la rend accessible au frontend
3. **Logique simple** — Reçoit un nom et renvoie une salutation
4. **Sûreté du typage** — Les types d’entrée et de sortie sont définis

@note{type="tip" title="Comprendre les services et les liaisons"}
Les **services** sont des modules Go autonomes qui exposent des fonctionnalités à votre frontend. Il s’agit simplement de structures Go ordinaires, dotées de méthodes exportées, que vous enregistrez dans le champ `Services` de la configuration de votre application.

Les **liaisons** constituent le SDK TypeScript/JavaScript généré automatiquement qui permet à votre frontend d’appeler ces services. Lorsque vous exécutez `wails3 dev` ou `wails3 build`, Wails analyse les services enregistrés et génère des liaisons avec typage sûr dans `frontend/bindings/`.

Considérez les services comme l’API de votre backend, et les liaisons comme la bibliothèque cliente qui communique avec elle.

@end

### Enregistrez le service

Ouvrez `main.go` et recherchez l’enregistrement du service :

```go {title="main.go" highlight="4-6"}
app := application.New(application.Options{
    Name: "myapp",
    Services: []application.Service{
        application.NewService(&GreetService{}),
    },
    // ... other options
})
```

Cette opération enregistre votre `GreetService` auprès de Wails et rend toutes ses méthodes exportées accessibles au frontend.

### Le frontend

Ouvrez `frontend/src/main.js` :

```javascript {title="frontend/src/main.js"}
import {GreetService} from "../bindings/changeme";

window.doGreet = async () => {
    const nameElement = document.getElementById('name');

    const name = nameElement.value;
    if (!name) {
        return;
    }

    try {
        const result = await GreetService.Greet(name);
        showToast(result);
    } catch (err) {
        console.error(err);
    }
};
```

Ce gestionnaire simplifié remplace le `window.doGreet` du projet initial ; conservez l’utilitaire `showToast` existant ainsi que le reste du code de démarrage. Le bouton HTML généré appelle `doGreet`, et `showToast` rend son résultat visible. Le chemin d’importation reprend le nom du module Go (`changeme` sauf s’il a été remplacé).

**Concepts clés :**

1. **Liaisons générées automatiquement** — `GreetService` est importé depuis le code généré
2. **Appels avec typage sûr** — Les noms et les signatures des méthodes correspondent à votre code Go
3. **Asynchrones par défaut** — Tous les appels Go renvoient des promesses
4. **Gestion des erreurs** — Les erreurs provenant de Go sont interceptées dans un bloc try/catch

@note{type="info" title="Où se trouvent les liaisons ?"}
Les liaisons générées se trouvent dans `frontend/bindings/`. Elles sont créées automatiquement lorsque vous exécutez `wails3 dev` ou `wails3 build`.

**Ne modifiez jamais ces fichiers manuellement** : ils sont régénérés à chaque compilation.

@end

## Personnalisez votre application

Ajoutons une fonctionnalité pour comprendre le flux de travail.

### Ajoutez une fonctionnalité « Greet Many »

@steps
### Ajoutez la méthode à GreetService
Ajoutez ceci à `greetservice.go` :

```go {title="greetservice.go"}
func (g *GreetService) GreetMany(names []string) []string {
    greetings := make([]string, len(names))
    for i, name := range names {
        greetings[i] = "Hello " + name + "!"
    }
    return greetings
}
```

### L’application sera recompilée automatiquement
Enregistrez le fichier : `wails3 dev` recompilera automatiquement votre code Go et redémarrera l’application.

@note{type="info" title="Recompilation automatique"}
Les modifications du code Go déclenchent automatiquement une recompilation et un redémarrage. Les modifications du frontend sont rechargées à chaud sans redémarrage.

@end

### Utilisez-la dans le frontend
Ajoutez ceci à `frontend/src/main.js` :

```javascript {title="frontend/src/main.js"}
window.greetMany = async () => {
    const names = ['Alice', 'Bob', 'Charlie'];
    const greetings = await GreetService.GreetMany(names);
    console.log(greetings);
};
```

Ouvrez la console du navigateur et appelez `greetMany()` : le tableau de salutations s’affiche.

@end

## Compilez pour la production

Lorsque vous êtes prêt à distribuer votre application :

```bash
wails3 build
```

**Résultat de cette opération :**

- Compile le code Go avec des optimisations
- Compile le frontend pour la production en le minifiant
- Crée un exécutable natif dans `bin/`

@tabs{sync-key="os"}
[Windows]
**Sortie :** `bin/myapp.exe`

Double-cliquez dessus pour l’exécuter. Aucune dépendance n’est nécessaire, car WebView2 fait partie de Windows.

[macOS]
**Sortie :** `bin/myapp.app`

Faites-le glisser dans le dossier Applications ou double-cliquez dessus pour l’exécuter.

[Linux]
**Sortie :** `bin/myapp`

Exécutez-le avec `./bin/myapp` ou créez un fichier `.desktop` pour votre lanceur.

@end

@note{type="tip" title="Compilations multiplateformes"}
Vous souhaitez compiler pour d’autres plateformes ? Consultez [Compilations multiplateformes →](/guides/build/cross-platform/)

@end

## Ce que nous avons appris

**Structure du projet**

- `main.go` pour le backend Go
- `frontend/` pour le code de l’interface utilisateur
- `Taskfile.yml` pour les tâches de build

**Services**

- Créez des structures Go avec des méthodes exportées
- Enregistrez-les avec `application.NewService()`
- Les méthodes sont automatiquement accessibles dans le frontend

**Bindings**

- Définitions TypeScript générées automatiquement
- Appels de fonctions avec typage sûr
- Asynchrones par défaut (promesses)

**Flux de développement**

- `wails3 dev` pour le rechargement à chaud
- Les modifications du code Go déclenchent automatiquement une nouvelle compilation et un redémarrage
- Les modifications du frontend sont immédiatement rechargées à chaud

---

**Des questions ?** Rejoignez [Discord](https://discord.gg/JDdSxwjhGf) et posez-les à la communauté.
