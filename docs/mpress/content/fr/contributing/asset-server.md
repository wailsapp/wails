---
title: "Serveur de ressources"
description: "Comment Wails v3 sert et intègre vos ressources web en développement et en production"
slug: "contributing/asset-server"
sourcePath: "contributing/asset-server.md"
---

## Vue d’ensemble

Chaque application Wails est distribuée sous la forme d’un **unique exécutable natif** qui réunit :

1. Votre backend *Go*
2. Un frontend *web* (HTML + JS + CSS)

Le **serveur de ressources** est le lien qui rend cela possible. Il possède **deux modes de fonctionnement**, sélectionnés lors de la compilation au moyen de tags de build Go :

| Mode | Tag | Objectif |
| --- | --- | --- |
| **Développement** | `//go:build !production` | Itérations rapides avec rechargement à chaud |
| **Production** | `//go:build production` | Ressources intégrées, sans dépendance |

L’implémentation se trouve dans `v3/internal/assetserver/`, avec une répartition claire entre les fichiers :

```
build_dev.go              # ⬅️ dev-only entrypoint (!production build tag)
build_production.go       # ⬅️ production-only entrypoint (production build tag)
assetserver.go            # Shared core
assetserver_dev.go        # Dev proxy/disk handler
assetserver_webview.go    # WebView-side adapter
assetserver_darwin.go     # OS-specific helpers (also linux/windows variants)
asset_fileserver.go       # Shared static file logic
content_type_sniffer.go   # MIME type detection
mimecache.go              # Cached MIME lookups
ringqueue.go              # Tiny in-memory LRU
options.go                # Configuration struct
middleware.go             # http.Handler middleware type
bundled_assetserver.go    # Hand-written wrapper around embedded bundles
bundledassets/            # Embedded runtime JS assets
```

---

## Mode développement

### Cycle de vie

1. `wails3 dev` démarre et **lance votre serveur de développement frontend** (Vite, SvelteKit, React-SWC…) en exécutant la tâche définie dans `build/Taskfile.yml` (généralement `npm run dev`).
2. La CLI définit `WAILS_VITE_PORT` sur le port de développement de Wails et `FRONTEND_DEVSERVER_URL` sur l’URL **complète** (`http://host:port` / `https://host:port`) qui pointe vers le serveur de développement actif du framework. Consultez `internal/commands/dev.go`.
3. Le serveur de ressources de développement (inclus dans la compilation au moyen de `//go:build !production` dans `build_dev.go`) lit `FRONTEND_DEVSERVER_URL` par l’intermédiaire de `GetDevServerURL()` et lui transmet par proxy inverse le trafic qui ne concerne pas le runtime.
4. Les fichiers statiques (`/assets/logo.svg`) peuvent être **servis directement depuis le disque** par l’intermédiaire de `asset_fileserver.go` (pour plus de rapidité), tandis que tout élément inconnu est **transmis par proxy** au serveur de développement du framework, ce qui permet le remplacement de modules à chaud *instantané*.

```
┌─────────┐  /wails/runtime.js     ┌─────────────┐
│ Browser │ ── embedded runtime ──▶│   Runtime   │
├─────────┤                        └─────────────┘
│   JS    │  / (index.html)        proxy / -> Vite via FRONTEND_DEVSERVER_URL
└─────────┘ ◀─────────────┐
              AssetServer │
                          ▼
                   ┌────────────┐
                   │  Vite Dev  │
                   │   Server   │
                   └────────────┘
```

### Fonctionnalités

- **Rechargement à chaud** — Vite, SvelteKit, etc. injectent le HMR par WebSocket ; le serveur de ressources de développement le transmet de manière transparente.
- **Prise en charge des source maps** — comme les ressources ne sont pas regroupées, les outils de développement de votre navigateur associent les erreurs au code source d’origine.
- **Aucune recompilation de Go** — seul le frontend est recompilé ; le code Go continue de s’exécuter jusqu’à ce que vous modifiiez des fichiers `.go`.

### Changer de framework

Le proxy de développement est **indépendant du framework**. La CLI Wails publie deux variables d’environnement lorsqu’elle lance votre tâche de développement :

| Variable d’environnement | Source | Signification |
| --- | --- | --- |
| `WAILS_VITE_PORT` | `internal/commands/dev.go` (constante `wailsVitePort`) | Port de développement par défaut (9245, sauf si `--port` est fourni) — votre configuration Vite devrait le respecter |
| `FRONTEND_DEVSERVER_URL` | `internal/commands/dev.go` | URL complète vers laquelle Wails transmettra les requêtes ; lue en Go par l’intermédiaire de `assetserver.GetDevServerURL()` (`build_dev.go`) |

Il n’existe aucune variable d’environnement `VITE_PORT`, `FRONTEND_DEV_PORT` ou `WAILSDEV_VERBOSE` dans l’arborescence v3.

Ajoutez un nouveau modèle → définissez sa tâche de développement → le serveur de ressources fonctionne immédiatement.

---

## Mode production

Lorsque vous exécutez `wails3 build`, le pipeline :

1. Exécute le **build de production** du frontend (`npm run build`), qui génère `frontend/dist/**`.
2. **Intègre** ce répertoire à l’application au moyen de `go:embed` dans le propre package de votre application (généralement `//go:embed all:frontend/dist` à côté de `main.go`).
3. Compile le binaire Go avec `-tags production` (transmis par l’intermédiaire de `EXTRA_TAGS` par le wrapper Taskfile).

`internal/assetserver/build_production.go` est le stub associé au tag de build qui active le chemin de code de production. `internal/assetserver/bundled_assetserver.go` est **écrit à la main** : il encapsule le JS du runtime dans `bundledassets/` et n’est pas un fichier généré.

### Traitement des requêtes

Le gestionnaire proprement dit est `internal/assetserver/assetserver.go` / `asset_fileserver.go`. Conceptuellement :

1. Recherche les ressources statiques intégrées au chemin demandé.
2. Se rabat sur `index.html` pour le routage de la SPA.
3. Détecte le type de contenu si l’extension est inconnue (`content_type_sniffer.go`).
4. Définit des en-têtes de cache appropriés.

- **Détection du type MIME** — pour les fichiers sans extension, le type de contenu est détecté à partir des ~512 premiers octets (`content_type_sniffer.go`), puis le résultat est mis en cache dans `mimecache.go` / `ringqueue.go`.
- **En-têtes de sécurité** — interdit la navigation `file://` et définit `nosniff`.

Comme tout est intégré, le binaire distribué ne possède **aucune dépendance externe** (même sous Windows).

---

## Faire le lien entre le développement et la production

Du point de vue de `pkg/application`, les deux modes exposent la **même interface publique** : une structure `AssetOptions` dotée d’un `Handler http.Handler`, ainsi que des middlewares et le raccordement au cycle de vie dans `internal/assetserver/`. Le passage du mode développement au mode production s’effectue entièrement au moyen des tags de build Go ; le code de l’application est donc identique dans les deux modes.

```go
app := application.New(application.Options{
    Assets: application.AssetOptions{
        Handler: application.AssetFileServerFS(assetsFS),
    },
})
```

---

## Intégration des frameworks frontend

### Modèles

Chaque modèle fourni (React, Vue, Svelte, Solid, Vanilla…) contient :

- `build/Taskfile.yml`
- `frontend/vite.config.ts` (ou équivalent)

Leur configuration Vite (ou équivalente) lit `WAILS_VITE_PORT` et lie le serveur de développement à ce port. La CLI publie ensuite la valeur réelle de `FRONTEND_DEVSERVER_URL` afin que le proxy intégré à l’application puisse l’utiliser.

Les frameworks restent totalement découplés de Go :

- Il n’est pas nécessaire d’importer un SDK JavaScript Wails lors de la compilation : `/wails/runtime.js` est servi par le serveur de ressources lors de l’exécution.
- Tout framework doté d’un serveur de développement HTTP peut s’intégrer.

---

## Extension et personnalisation

Vous avez besoin d’en-têtes personnalisés, d’une authentification ou de gzip ?

1. Définissez un `middleware.Middleware` (alias de `func(http.Handler) http.Handler`, déclaré dans `internal/assetserver/middleware.go`).
2. Intégrez-le à votre `application.AssetOptions` au moyen de la configuration exposée par `internal/assetserver/options.go`.
3. Le comportement est identique en développement et en production : il n’existe pas de liste de middlewares propre à chaque mode.

---

## Principaux fichiers sources

| Fichier | Rôle |
| --- | --- |
| `build_dev.go` / `build_production.go` | Enveloppes avec contraintes de compilation sélectionnant le mode développement ou production |
| `assetserver.go` / `asset_fileserver.go` | Gestionnaire HTTP principal |
| `assetserver_dev.go` | Proxy inverse vers `FRONTEND_DEVSERVER_URL` |
| `bundled_assetserver.go` | Enveloppe écrite manuellement autour de `bundledassets/` |
| `options.go` | Configuration destinée à `application.AssetOptions` |
| `mimecache.go` / `ringqueue.go` | Cache MIME et petit cache LRU |

---

## Pièges et débogage

- **Écran blanc en production** — généralement dû au routage de l’application monopage (SPA) : vérifiez que votre serveur de développement sert `index.html` pour les chemins inconnus et que le mécanisme de repli du gestionnaire de production intégré est bien atteint.
- **404 en développement** — votre configuration Vite ne se lie pas à `WAILS_VITE_PORT`, ou la CLI n’a pas pu joindre le serveur de développement pour renseigner `FRONTEND_DEVSERVER_URL`.
- **Ressources volumineuses** — leur intégration augmente la taille du binaire. Servez les médias volumineux depuis une origine distincte ou diffusez-les via un `http.Handler` personnalisé.

---

Vous savez maintenant comment le **serveur de ressources** de Wails fournit votre code web à la fenêtre native, aussi bien en **développement** qu’en **production**. Maîtrisez cette couche pour diagnostiquer les problèmes de chargement, ajouter des middlewares ou même adopter en toute confiance une chaîne d’outils frontend entièrement différente.
