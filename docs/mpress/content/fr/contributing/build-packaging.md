---
title: "Pipeline de compilation et de packaging"
description: "Ce qui se passe en interne lorsque vous exécutez `wails3 build`, comment les binaires multiplateformes sont produits et comment les programmes d’installation sont générés pour chaque système d’exploitation."
slug: "contributing/build-packaging"
sourcePath: "contributing/build-packaging.md"
---

`wails3 build` est volontairement **minimal** : il s’agit d’une enveloppe Taskfile qui transmet les balises de compilation supplémentaires à la tâche `build` du projet hôte. L’essentiel du travail est effectué dans le propre fichier `build/Taskfile.yml` du projet (généré par `wails3 init`), dans `internal/commands/build-assets.go` (qui gère les ressources intégrées lors de la compilation), ainsi que dans `internal/packager` (packaging Linux avec nfpm) et `internal/commands/appimage.go`, `internal/commands/msix.go`, `internal/commands/dmg/dmg.go`, et `internal/commands/dot_desktop.go` (programmes d’installation propres à chaque plateforme).

Cette page aborde les sujets suivants :

1. Les véritables points d’entrée de la CLI
2. Le flux de compilation piloté par Taskfile
3. L’intégration des ressources et l’injection des informations de compilation
4. Les moteurs de packaging propres à chaque plateforme
5. La personnalisation du pipeline
6. Le dépannage

---

## 1. Véritables points d’entrée de la CLI

```
wails3 build       → internal/commands.Build       (in task_wrapper.go)
wails3 package     → internal/commands.Package     (in task_wrapper.go)
wails3 generate build-assets → GenerateBuildAssets (in build-assets.go)
wails3 update build-assets   → UpdateBuildAssets   (in build-assets.go)
wails3 tool buildinfo        → BuildInfoOptions    (in tool_buildinfo.go)
wails3 tool package          → internal/packager   (nfpm wrapper)
wails3 generate .desktop     → in dot_desktop.go
```

`internal/commands/task_wrapper.go` :

```go
func Build(buildFlags *flags.Build, otherArgs []string) error {
    // forwards --tags / EXTRA_TAGS, then defers to a Taskfile target
    return wrapTask("build", otherArgs)
}
```

`flags.Build` n’expose qu’une **seule** option : `--tags` (transmise sous la forme `EXTRA_TAGS=`). `wails3 build` ne propose **aucune** option `-platform`, `-o`, `-skipbindings`, `-skip-package`, `-package`, `-ldflags`, `-verbose`, `-debug`, `-devbuild`, `-icon` ou `-clean`. La compilation croisée, les chemins de sortie, les icônes, etc. se configurent dans **`Taskfile.yml`**, **`build/config.yml`** et les commandes auxiliaires `wails3 generate icons` / `wails3 generate build-assets`.

`build/build.json` ne fait **pas** partie de la v3 : la configuration repose sur `Taskfile.yml` ainsi que sur `build/config.yml`.

---

## 2. Flux de compilation piloté par Taskfile

Un projet fraîchement initialisé comprend un fichier `build/Taskfile.yml` dont les espaces de noms sont approximativement les suivants :

| Espace de noms | Tâches (sélection) |
| --- | --- |
| `darwin:` | `build`, `build:universal`, `package`, `run`, `dev` |
| `windows:` | `build`, `package`, `run`, `dev` |
| `linux:` | `build`, `package`, `run`, `dev` |
| `common:` | `update:build-assets`, `generate:icons`, `generate:syso` |

Par défaut, `wails3 build` appelle l’espace de noms `build` du système d’exploitation hôte ; le Taskfile du projet exécute ensuite `go build` avec les options propres à l’hôte. Pour compiler pour un autre système d’exploitation, exécutez directement sa tâche (par exemple `wails3 task darwin:build:universal`), au lieu de transmettre une option à `wails3 build`.

Le répertoire de sortie par défaut est **`bin/<APP_NAME>`** (sans préfixe `build/bin/`).

---

## 3. Ressources intégrées lors de la compilation et informations de compilation

| Objet | Fichier |
| --- | --- |
| Génération/mise à jour des ressources de compilation | `internal/commands/build-assets.go` |
| Affichage des informations de compilation (CLI : `wails3 tool buildinfo`) | `internal/commands/tool_buildinfo.go` — affiche des informations ; ce n’est **pas** un injecteur `ldflags` |
| Stub de production | `internal/assetserver/build_production.go` — `//go:build production` |
| Bundles frontend | intégrés via `//go:embed` dans le propre package de l’application (par exemple à côté de `main.go`) |
| Ressources Windows (`.syso`) | `internal/commands/syso.go` — génère `rsrc_windows_<arch>.syso` |
| MSIX Windows | `internal/commands/msix.go` + `internal/commands/webview2/` |
| Entrées DMG macOS | `internal/commands/dmg/` |
| `.desktop` Linux | `internal/commands/dot_desktop.go` |

La CLI ne génère pas automatiquement `bundled_assetserver.go` pour votre application : `internal/assetserver/bundled_assetserver.go` est **écrit à la main** et encapsule l’environnement d’exécution JavaScript intégré sous `bundledassets/`.

---

## 4. Moteurs de packaging

### Linux

Sous Linux, le packaging est piloté par **nfpm** (et non par `fpm`) :

- `internal/packager/packager.go` encapsule `github.com/goreleaser/nfpm/v2` et expose `CreatePackageFromConfig(pkgType, configPath, output)` / `CreatePackageFromConfigWriter(...)`.
- Les projets générés comprennent les configurations de style nfpm `myapp.DEB`, `myapp.RPM` et `myapp.ARCHLINUX` sous `internal/commands/` (utilisées par `wails3 tool package`).
- La génération d’AppImage se trouve dans `internal/commands/appimage.go`, qui appelle `linuxdeploy` + `linuxdeploy-plugin-gtk` (le plugin est fourni à l’emplacement `internal/commands/linuxdeploy-plugin-gtk.sh`).

`wails3 build` ne propose **aucune** option `-package deb`/`rpm`. Utilisez `wails3 tool package` ou la cible Taskfile propre à la plateforme.

### macOS

- `darwin:package`, issu du Taskfile du projet, produit le bundle `.app`.
- Les ressources DMG se trouvent sous `internal/commands/dmg/` ; un projet peut encapsuler le bundle dans un DMG avec `hdiutil` une fois `darwin:package` terminé (le Taskfile des modèles récents comprend un utilitaire `dmg`).
- Les identifiants CFBundle, la version et les mentions de copyright proviennent des options `-product*` lors de `wails3 init`, ainsi que de `build/config.yml`.

### Windows

- Le packaging Windows cible **MSIX** (et non WiX/MSI). Consultez `internal/commands/msix.go` et `internal/commands/webview2/` pour connaître le workflow complet.
- Il n’existe **aucun** `internal/commands/packager.go` ni **aucun** répertoire `internal/commands/windows_resources/`.
- La signature de code facultative de l’exécutable s’effectue via `wails3 tool sign` (Authenticode) — consultez `internal/commands/sign.go`.

---

## 5. Personnaliser le pipeline

| Besoin | Approche |
| --- | --- |
| Tags de build supplémentaires | `wails3 build --tags myFeature,otherTag` |
| Linter / étape préalable au build | Ajoutez une tâche à `build/Taskfile.yml` et faites-en une dépendance de la tâche `build` propre au système d’exploitation |
| Compilation croisée | Exécutez la tâche correspondant au système d’exploitation (par exemple `wails3 task linux:build`) — il n’existe aucune option `-platform` |
| Ignorer le packaging | Exécutez simplement la tâche `build` ; `package` est distincte |
| Outil de packaging personnalisé | Placez une configuration sous `internal/commands/myapp.*` et appelez `wails3 tool package` avec `-config <file>` |
| Supprimer les symboles | Modifiez la tâche `darwin:/windows:/linux:` `build` afin de transmettre `-ldflags "-s -w"` directement à `go build` — `wails3 build` ne possède lui-même aucune option `-ldflags` |

Toutes les cibles du Taskfile respectent les variables d’environnement publiées par Wails (`APP_NAME`, `WAILS_VITE_PORT`, `FRONTEND_DEVSERVER_URL`, …), de sorte que les tâches personnalisées peuvent s’appuyer sur elles.

---

## 6. Résolution des problèmes

| Symptôme | Cause probable | Solution |
| --- | --- | --- |
| **`ld: framework not found WebKit` (mac)** | Outils en ligne de commande de Xcode manquants | `xcode-select --install` |
| **Fenêtre vide dans le build de production** | Échec du build du frontend ou routage de la SPA | Vérifiez que `frontend/dist/index.html` existe et que votre gestionnaire de ressources l’utilise comme solution de repli |
| **Outils de packaging MSIX manquants** | SDK `WebView2` / outils MSIX non installés | Exécutez `wails3 task install:msix:tools` |
| **`linuxdeploy` introuvable** | Plugin absent du PATH | Installez `linuxdeploy`, puis exécutez `internal/commands/linuxdeploy-plugin-gtk.sh` via l’étape d’installation automatique de la CLI |

`wails3 build` ne possède aucune option `-verbose`. Définissez `TASK_X_VERBOSE=1` (Taskfile) ou examinez directement la cible de la tâche pour voir les commandes exécutées.

---

## 7. Carte des principaux fichiers sources

| Fonction | Fichier |
| --- | --- |
| Wrapper de build | `internal/commands/task_wrapper.go` (`Build`, `Package`, `SignWrapper`, `wrapTask`) |
| Génération des ressources de build | `internal/commands/build-assets.go` (`GenerateBuildAssets`, `UpdateBuildAssets`) |
| Affichage des informations de build | `internal/commands/tool_buildinfo.go` |
| Générateur d’AppImage | `internal/commands/appimage.go` |
| Packaging Linux (nfpm) | `internal/packager/packager.go`, `internal/commands/myapp.{DEB,RPM,ARCHLINUX}` |
| MSIX Windows | `internal/commands/msix.go`, `internal/commands/webview2/` |
| Générateur de ressources Windows | `internal/commands/syso.go` |
| Ressources DMG macOS | `internal/commands/dmg/` |
| Générateur `.desktop` | `internal/commands/dot_desktop.go` |
| Constantes de version | `internal/version/version.go` |

Gardez ce tableau à portée de main lorsque vous recherchez la cause d’un échec de build.

---

Vous disposez maintenant d’une vue d’ensemble complète, du **code source** au **programme d’installation**. En résumé, `wails3 build` n’est lui-même qu’un wrapper léger : presque toutes les personnalisations s’effectuent dans les fichiers `Taskfile.yml` / `build/config.yml` du projet, ou au moyen des sous-commandes explicites `wails3 generate …` / `wails3 tool …`. Bonne livraison !
