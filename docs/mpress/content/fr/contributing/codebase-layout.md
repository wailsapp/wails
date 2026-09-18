---
title: "Organisation du code source"
description: "Organisation du dépôt Wails v3 et articulation de ses différents éléments"
slug: "contributing/codebase-layout"
sourcePath: "contributing/codebase-layout.md"
---

Wails v3 se trouve dans un **monorepo** qui contient l’environnement d’exécution du framework, la CLI, les exemples, la documentation et la chaîne d’outils de compilation. Cette page présente la *structure des répertoires* utile à toute personne souhaitant explorer le fonctionnement interne.

## Vue d’ensemble

```
wails/
├── v3/               # ⬅️ Everything specific to Wails v3 lives here
├── v2/               # Legacy v2 implementation (can be ignored for v3 work)
├── docs/             # M-Press-powered v3 docs site (this page!)
├── website/          # Docusaurus v2 site and marketing pages (main site)
├── scripts/          # Misc helper scripts (e.g. sponsor image generator)
└── *.md              # Project-wide meta files (CHANGELOG, LICENSE, …)
```

Nous allons maintenant examiner plus en détail l’arborescence **`v3/`**.

## Racine `v3/`

```
v3/
├── cmd/          # Compilable commands (currently only the wails3 CLI)
├── internal/     # Framework implementation (not public API)
├── pkg/          # Public Go packages — the API surface
├── tasks/        # Taskfile-based release / generation utilities
├── wep/          # RFC-style proposals (Wails Enhancement Proposals)
├── tests/        # Integration test harness
├── go.mod
└── go.sum
```

> Les modèles de projet sont fournis sous `internal/templates/` (un dossier par pile technologique
>
> de framework, plus `base/`, `_common/` et `ios/`). Il n’existe aucun répertoire `v3/templates/`
>
> au niveau supérieur.

### Modèle mental

1. **`pkg/`** expose *ce que les développeurs d’applications importent*\
2. **`internal/`** contient *la manière dont la magie est mise en œuvre*\
3. **`cmd/wails3`** pilote *le cycle de vie du projet et les compilations*\

Tous les autres éléments soutiennent ces trois piliers.

---

## `cmd/` – Commandes

| Chemin | Remarques |
| --- | --- |
| `v3/cmd/wails3` | Le **point d’entrée de la CLI**. Un minuscule `main.go` délègue toute la logique aux packages de `internal/commands`. |
| `internal/commands/*` | Sous-commandes (init, dev, build, doctor, …). Chacune se trouve dans son propre fichier afin d’être facile à repérer. |
| `internal/commands/task_wrapper.go` | Fait le lien entre les options de la CLI et le pipeline de compilation Taskfile. |

La CLI prend en charge :

- **Génération de la structure initiale du projet** (`init`, génération de modèles)\
- **Orchestration du serveur de développement** (`dev`, rechargement à chaud)\
- **Compilations de production et création des packages** (`build`, `package`, enveloppes propres aux plateformes)\
- **Diagnostics** (`doctor`)\

---

## `internal/` – Le cœur du moteur

```
internal/
├── assetserver/  # Serving & embedding web assets
├── buildinfo/    # Reproducible build metadata
├── commands/     # CLI mechanics (see above)
├── runtime/      # Build-tag glue + embedded JS runtime sources
├── generator/    # Static analysis & binding generator
├── templates/    # Project templates (frontend stacks)
├── packager/     # nfpm wrapper used by `wails3 tool package`
├── capabilities/ # Host OS capability probing
├── dbus/         # Generic D-Bus helper
├── service/      # Service-template scaffolding (`wails3 generate service`)
└── ...           # [other helper sub-packages: flags, hash, term, …]
```

### Principaux sous-packages

| Package | Responsabilité | Point de connexion |
| --- | --- | --- |
| `runtime` | Contient la petite couche de liaison `runtime{,_darwin,_linux,_windows,_android,_dev,_prod}.go` fondée sur des balises de compilation, ainsi que l’environnement d’exécution JS intégré sous `runtime/desktop/`. Le code propre à chaque système d’exploitation pour les fenêtres, le presse-papiers, les boîtes de dialogue et la zone de notification se trouve dans `pkg/application/*_{darwin,linux,windows}.go`. | Importé indirectement via `pkg/application`. |
| `assetserver` | Serveur de fichiers à deux modes :<br />• Développement : sert les fichiers depuis le disque et agit comme proxy vers Vite (`build_dev.go`)<br />• Production : intègre les ressources via `go:embed` (`build_production.go`) | Initialisé par `pkg/application` au démarrage. |
| `generator` | Analyse le code source Go afin de créer les **métadonnées de liaison** qui servent ensuite à produire les fichiers de liaison TypeScript/JS et les constantes d’événements. Points d’entrée : `generator.Generate` / `generator.Generator` au-dessus de `collect/` + `render/`. | Déclenché par `wails3 generate bindings`. |
| `packager` | Enveloppe `nfpm` utilisée pour générer les artefacts Linux `deb`/`rpm`/`archlinux` (pilotée par les configurations nfpm `myapp.DEB`/`.RPM`/`.ARCHLINUX` sous `internal/commands/`). | Invoqué par `wails3 tool package`. Le code de prise en charge des formats DMG pour macOS et MSIX pour Windows se trouve sous `internal/commands/{dmg,msix.go,webview2/}`. |

Les utilitaires complémentaires (par exemple `s/`, `hash/` et `flags/`) assurent le découplage des préoccupations internes.

---

## `pkg/` – API publique

```
pkg/
├── application/  # Core API: App, windows, menus, dialogs, events, managers
├── events/       # Event constants (Common/Mac/Windows/Linux) + generator
├── services/     # Optional built-in services (notifications, kvstore, …)
├── doctor-ng/    # New-style `wails3 doctor-ng` checks
├── errs/         # Shared error types
├── icons/        # Default platform icons
├── mac/          # macOS-only helpers
└── w32/          # Windows Win32 helpers
```

> Il n’existe aucun package `pkg/runtime/`, `pkg/options/` ou `pkg/menu/`. Les options des fenêtres et des menus
>
> se trouvent aux côtés de `pkg/application` (par exemple `WebviewWindowOptions`, `Menu`,
>
> `MenuItem`), tandis que `assetserver/` se trouve sous `internal/`.

`pkg/application` initialise un programme Wails :

```go
func main() {
    app := application.New(application.Options{
        Name: "MyApp",
        Assets: application.AssetOptions{
            Handler: application.AssetFileServerFS(assetsFS),
        },
    })
    window := app.Window.New()
    window.SetTitle("Hello").SetSize(1024, 768)
    _ = app.Run()
}
```

En interne, il :

1. relie la couche de liaison `internal/runtime` fondée sur des balises de compilation au code propre à chaque système d’exploitation dans `pkg/application/`
2. configure une instance de `internal/assetserver`
3. enregistre tous les processeurs de messages pilotés par les liaisons
4. entre dans le thread principal du système d’exploitation

---

## `internal/templates/` – Modèles de structure initiale

`internal/templates/` fournit des **modèles de base** (structure Go sous `base/`, `_common/`, `ios/`) et des **habillages d’interface** (`vanilla[-ts]`, `react[-ts]`, `react-swc[-ts]`, `lit[-ts]`, `preact[-ts]`, `qwik[-ts]`, `solid[-ts]`, `svelte[-ts]`, `sveltekit[-ts]`, `vue[-ts]`).

Lors de l’exécution de `wails3 init -t react`, la CLI :

1. copie les fichiers Go de `_common`
2. fusionne le pack d’interface souhaité
3. Exécute `go mod tidy` (peut être ignoré avec `--skipgomodtidy`)

La modification des modèles **n’affecte pas** les applications existantes, mais uniquement les futurs `init`. Les exemples publics se trouvent sous `v3/examples/` ; ils ne remplacent pas les suites de tests automatisés décrites dans la documentation destinée aux contributeurs.

---

## `tasks/` – Automatisation des versions

Les Taskfiles encapsulent les opérations complexes de compilation croisée, de mise à jour des versions et de génération du journal des modifications. Ils sont utilisés par programmation par `internal/commands/task.go`, afin que la même logique alimente la **CLI** et la **CI**.

---

## Interaction entre les composants

```d2
direction: down
CLI: CLI wails3
Generator: internal/generator
AssetDev: assetserver (développement)
Packager: internal/packager
AppRuntime: {
  label: Environnement d’exécution de l’application
  ApplicationPkg: pkg.application
  InternalRuntime: internal.runtime
  OSAPIs: API du système d’exploitation
}
CLI -> Generator: compiler / générer
CLI -> AssetDev: développement
CLI -> Packager: empaqueter
Generator -> ApplicationPkg: liaisons
ApplicationPkg -> InternalRuntime
InternalRuntime -> OSAPIs
ApplicationPkg -> AssetDev
ApplicationPkg.label: ApplicationPkg
InternalRuntime.label: InternalRuntime
OSAPIs.label: OSAPIs
```

*CLI → générateur → runtime* constitue le chemin principal qui mène du **code source** à l’**application de bureau en cours d’exécution**.

---

## Conseils d’orientation

| Vous devez comprendre… | Consultez… |
| --- | --- |
| Couches d’adaptation aux plateformes | `pkg/application/*_darwin.go`, `*_linux.go`, `*_windows.go` (fenêtre, presse-papiers, boîtes de dialogue, zone de notification, thread principal, events_common). cgo sous Linux : `pkg/application/linux_cgo*.go`. |
| Protocole de passerelle | `pkg/application/messageprocessor*.go` |
| Flux de gestion des ressources | `internal/assetserver/` (`build_dev.go` ou `build_production.go`) |
| Flux d’empaquetage | `internal/commands/{appimage,msix,dot_desktop,dmg/}.go`, `internal/packager/` |
| Moteur de modèles | `internal/templates/` (`templates.Install`, `templates.GetDefaultTemplates`) |
| Analyse statique | `internal/generator/{generate.go,collect/,render/}` |

---

Vous disposez désormais d’une **carte mentale** du dépôt. Utilisez-la avec `ripgrep`, les commandes « Go to File/Symbol » de votre IDE et les exemples d’applications pour explorer plus en profondeur n’importe quelle fonctionnalité. Bon développement !
