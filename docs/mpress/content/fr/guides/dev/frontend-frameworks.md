---
title: "Utiliser d’autres frameworks frontend"
description: "Comment utiliser un framework sans modèle intégré en plaçant votre propre projet Vite dans le répertoire frontend"
slug: "guides/dev/frontend-frameworks"
sourcePath: "guides/dev/frontend-frameworks.md"
---

Wails fournit des modèles de démarrage intégrés pour un ensemble volontairement restreint de frameworks :

| Modèle | Langage |
| --- | --- |
| `vanilla` | TypeScript (par défaut) |
| `vanilla-js` | JavaScript |
| `react` | TypeScript |
| `react-js` | JavaScript |
| `vue` | TypeScript |
| `svelte` | TypeScript |

Cela ne signifie pas que ce sont vos seules options. Le frontend d’une application Wails est **simplement un projet web** : tout ce qui produit du HTML, du CSS et du JS statiques fonctionnera. Si le framework de votre choix (Solid, Preact, Lit, Qwik, SvelteKit, Angular, etc.) ne dispose pas de modèle, vous pouvez générer vous-même sa structure en quelques minutes.

## Utilisation du répertoire `frontend/`

Wails ne se préoccupe pas du framework qui se trouve dans `frontend/`. Il repose uniquement sur un contrat minimal et indépendant du framework :

- **`frontend/dist/` correspond à ce qui est distribué.** `main.go` incorpore le frontend compilé avec `//go:embed all:frontend/dist` et le fournit depuis le serveur de ressources. Votre processus de compilation doit produire un bundle statique dans `frontend/dist/`, le répertoire de sortie par défaut de Vite.
- **La compilation est pilotée par `frontend/package.json`.** Pendant `wails3 build`, Wails exécute le script `build` du frontend ; pendant `wails3 dev`, il exécute `dev` et utilise un proxy vers le serveur de développement Vite pour permettre le rechargement à chaud.
- **Les liaisons sont générées dans `frontend/bindings/`.** Wails examine vos services Go enregistrés et y écrit un SDK avec typage sûr. Importez-le comme n’importe quel autre module :
  ```js
  import { GreetService } from "./bindings/changeme";
  ```


- **Le serveur de développement s’exécute sur un port fixe.** `wails3 dev` utilise un proxy vers Vite sur le port défini dans `WAILS_VITE_PORT` (`9245` par défaut) ; réglez donc `server.port` sur ce port et activez `strictPort: true`. Il s’agit d’une configuration Vite ordinaire : aucun plugin Wails n’intervient.
- **Facultatif : événements personnalisés typés.** Les modèles intégrés enregistrent également le plugin `@wailsio/runtime/plugins/vite`. Il n’est nécessaire que si vous utilisez des événements personnalisés *typés* : il injecte dans le runtime les définitions de types d’événements générées et fait échouer la compilation tant que les liaisons n’ont pas été générées. Si vous utilisez uniquement l’API `Events.On("time", …)` fondée sur des chaînes de caractères, vous pouvez l’omettre.

Tout le reste — composants, routage, état et styles — relève entièrement de votre framework.

## Générer la structure de n’importe quel framework avec Vite

Le moyen le plus rapide consiste à partir d’un modèle intégré afin de disposer de `main.go`, de `Taskfile`, des ressources de compilation et d’un service Go fonctionnel, puis à remplacer `frontend/` par un nouveau projet Vite pour votre framework.

@steps
### Créer un projet à partir du modèle par défaut
```bash
wails3 init -n myapp
cd myapp
```

### Remplacer `frontend/` par une application Vite pour votre framework
Vite peut générer la structure de la plupart des frameworks avec une seule commande. Choisissez un modèle :

```bash
# From the project root — e.g. Solid, Preact, Lit, Svelte, Vue, React, Vanilla
rm -rf frontend
npm create vite@latest frontend -- --template solid
```

Remplacez `solid` par n’importe quel modèle Vite : `preact`, `lit`, `svelte`, `vue`, `react`, `vanilla` ou leurs variantes `-ts` (`solid-ts`, `preact-ts`, etc.).

### Installer le runtime et configurer Vite pour le serveur de développement Wails
```bash
cd frontend
npm install @wailsio/runtime
```

`@wailsio/runtime` fournit les API JS (`Events`, `Browser`, boîtes de dialogue, etc.). La seule modification *obligatoire* dans `vite.config` concerne le port du serveur de développement, afin que `wails3 dev` puisse le trouver :

```ts {title="frontend/vite.config.ts"}
import { defineConfig } from "vite";

export default defineConfig({
  server: {
    host: "127.0.0.1",
    port: Number(process.env.WAILS_VITE_PORT) || 9245,
    strictPort: true,
  },
});
```

Uniquement si vous prévoyez d’utiliser des **événements personnalisés typés**, ajoutez également le plugin : il injecte les types d’événements générés et exige que les liaisons existent avant la compilation :

```ts {title="frontend/vite.config.ts" highlight="2,5"}
import { defineConfig } from "vite";
import wails from "@wailsio/runtime/plugins/vite";

export default defineConfig({
  plugins: [wails("./bindings")],
  server: { host: "127.0.0.1", port: Number(process.env.WAILS_VITE_PORT) || 9245, strictPort: true },
});
```

### Appeler vos services Go
Générez une première fois les liaisons, puis importez-les où vous le souhaitez dans vos composants :

```bash
wails3 generate bindings
```

```js
import { GreetService } from "./bindings/changeme";

const greeting = await GreetService.Greet("World");
```

### Exécuter le projet
```bash
wails3 dev
```

@end

@note{type="tip" title="Générer un projet JavaScript sans TypeScript"}
La même commande génère la structure d’un projet sans TypeScript : utilisez simplement un modèle Vite dont le nom ne porte pas le suffixe `-ts` :

```bash
npm create vite@latest frontend -- --template solid
```

@end

## Frameworks dotés de leur propre outil de génération

Quelques frameworks ne sont pas créés au moyen des modèles `create` de Vite et disposent de leurs propres outils. Ils restent compatibles : générez simplement leur structure avec leur commande native, puis ajoutez le plugin Wails :

- **SvelteKit :** `npx sv create frontend`. Utilisez l’adaptateur statique (`@sveltejs/adapter-static`) afin de produire un bundle statique et désactivez le rendu côté serveur (SSR).
- **Qwik :** `npm create qwik@latest`. Utilisez l’adaptateur statique (SSG).
- **Angular :** générez la structure avec `ng new`, définissez `outputPath` sur `dist` et configurez le script de compilation pour qu’il cible `ng build`.

La règle reste toujours la même : produisez une compilation statique dans `frontend/dist/`, conservez le plugin Vite `@wailsio/runtime` (ou importez directement le runtime) et importez vos liaisons Go depuis `frontend/bindings/`.

@note{type="info"}
Si vous créez une configuration aboutie pour un framework, envisagez de la publier sous forme de [modèle personnalisé](/guides/advanced/custom-templates/) afin que d’autres puissent l’utiliser directement avec `wails3 init -t`.

@end
