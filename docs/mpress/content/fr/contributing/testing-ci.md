---
title: "Tests et intégration continue"
description: "Comment Wails v3 garantit la qualité grâce aux tests unitaires, aux suites d’intégration, à la détection des accès concurrents non synchronisés et à l’intégration continue avec GitHub Actions."
slug: "contributing/testing-ci"
sourcePath: "contributing/testing-ci.md"
---

Les frameworks de bureau robustes exigent des tests à toute épreuve. Wails v3 emploie une **stratégie à plusieurs niveaux** :

| Niveau | Objectif | Outils |
| --- | --- | --- |
| Tests unitaires | Retour rapide sur les fonctions isolées | `go test ./...` |
| Tests du générateur et de la CLI | Valider `wails3 generate bindings` et les mécanismes internes de la CLI | `task test:generator`, `task test:cli` |
| Tests des modèles | Vérifier que chaque modèle distribué peut toujours être compilé | `task test:templates` |
| Détection des accès concurrents non synchronisés | Détecter les accès concurrents non synchronisés aux données dans le runtime et le pont | `go test -race ./...` |
| Matrice d’intégration continue | Confiance multiplateforme pour chaque PR | GitHub Actions |

Ce document explique **où se trouvent les tests**, **comment les exécuter** et **ce que le Taskfile orchestre**.

> Les recommandations sur les accès concurrents non synchronisés auxquelles les anciennes versions faisaient référence sous le nom `pkg/application/RACE.md`
>
> se trouvent aujourd’hui dans `v3/TESTING.md`.

---

## 1. Conventions relatives aux répertoires

```
v3/
├── internal/.../_test.go     # Unit tests for internal packages
├── pkg/.../_test.go          # Public API tests
├── tasks/events/generate.go  # Code generator for event constants (NOT a test harness)
├── tests/                    # Top-level integration test harness
└── TESTING.md                # Race / Cgo testing guidance
```

Consignes :

- **Conservez les tests unitaires à côté du code** (`foo.go` ↔ `foo_test.go`).
- Utilisez le style **boîte noire** pour les packages `pkg/` (`package application_test`) lorsqu’il améliore la propreté de l’API.
- Placez les fixtures partagées là où elles sont utilisées (cet arbre ne contient aucun package `internal/testutil/` central : raccordez plutôt les fonctions auxiliaires dans chaque package).

---

## 2. Tests unitaires

### Écriture des tests

```go
func TestEventConstants(t *testing.T) {
    assert.NotEmpty(t, events.Common.WindowFocus)
}
```

Recommandations :

- Utilisez [`stretchr/testify`](https://github.com/stretchr/testify), déjà présent dans `go.mod`.
- Privilégiez les tests **pilotés par table** lorsque plusieurs entrées, cas limites ou résultats attendus mettent en œuvre le même comportement. Donnez un nom descriptif à chaque cas.
- Si nécessaire, remplacez les comportements propres aux plateformes par des implémentations factices sélectionnées à l’aide de balises de compilation (`foo_windows_test.go`, `foo_darwin_test.go`, …).

### Couverture attendue

La logique nouvelle ou modifiée doit atteindre une couverture de 100 % des instructions Go. Mesurez le package que vous avez modifié au lieu de vous fier à un pourcentage portant sur l’ensemble du dépôt :

```bash
go test ./pkg/application -coverprofile=coverage.out
go tool cover -func=coverage.out
```

Certains chemins ne peuvent raisonnablement pas être testés dans un environnement de test normal, par exemple les échecs propres à une plateforme, les comportements dépendant du matériel ou un mécanisme de repli défensif qu’il est impossible de déclencher sans risque. Limitez strictement ces exceptions et expliquez chaque chemin non couvert dans la description de la PR.

### Exécution locale

```bash
cd v3
go test ./... -cover
```

Ou via le Taskfile (avec les véritables cibles : il n’existe aucun raccourci `task test`) :

```
task test:cli            # CLI plumbing tests
task test:generator      # bindings generator round-trip tests
task test:templates      # build every shipped template
task test:infrastructure # supporting helpers
task test:examples       # exercise the example matrix (downloads as needed)
task test:all            # everything above
task sanity              # quick smoke check (also: sanity:gtk4)
task precommit           # what you should run before pushing
```

---

## 3. Tests d’intégration

`v3/tests/` héberge l’infrastructure de tests d’intégration entre packages. Les cibles `test:example:*` et `test:examples:*` du Taskfile pilotent les vérifications de compilation et de lancement sous darwin / windows / linux (y compris les matrices GTK3 / GTK4 basées sur Docker sous Linux).

> Les exemples exécutables se trouvent dans `v3/examples/`. Les cibles de test sélectionnent et compilent
>
> les exemples adaptés à la plateforme hôte ou à la matrice d’intégration continue.

Exécutez la suite de tests de bon fonctionnement de la plateforme hôte avec :

```
task test:examples       # host
task test:examples:all   # full matrix (slow)
```

---

## 4. Détection des accès concurrents non synchronisés

Les accès concurrents non synchronisés aux données sont fatals dans les runtimes d’interface graphique.

### Guide sur les accès concurrents non synchronisés

Consultez `v3/TESTING.md` pour connaître :

- les accès concurrents non synchronisés connus comme bénins et les raisons de leur suppression
- la manière d’interpréter les traces de pile qui traversent les frontières Cgo (GTK sous Linux + WebKit2GTK)

### Suite locale de détection des accès concurrents non synchronisés

```
go test -race ./...
```

> `wails3 dev` ne possède aucun indicateur `-race` : ses indicateurs de CLI sont `--config`, `--port`,
>
> et `-s` (pour activer HTTPS). Pour tester le runtime avec le détecteur d’accès concurrents non synchronisés,
>
> compilez une application de test avec `go build -race`, puis exécutez-la directement.

---

## 5. Workflows GitHub Actions

Fichiers de workflow réels sous `.github/workflows/` (vérifiés dans l’arbre) :

| Fichier | Objectif |
| --- | --- |
| `build-and-test-v3.yml` | Matrice principale de compilation et de test de la v3. Utilise `actions/setup-go@v5` avec `go-version: 1.25`. Exécute `task runtime:check`, `task runtime:test`, `task runtime:build`, `task test:examples` (ainsi que `BUILD_TAGS=gtk4 task test:examples` pour le chemin GTK4), `task generator:test:check`, `task install`, puis `wails3 build` comme test de bon fonctionnement. Les tâches Linux installent `libgtk-3-dev libwebkit2gtk-4.1-dev libwayland-dev build-essential pkg-config xvfb x11-xserver-utils at-spi2-core xdg-desktop-portal-gtk` et exécutent la suite de tests sous `dbus-run-session -- xvfb-run`. |
| `cross-compile-test-v3.yml` | Vérifications de cohérence de la compilation croisée |
| `auto-changelog-v3.yml`, `changelog-v3.yml` | Automatisation du journal des modifications |
| `nightly-release-v3.yml` | Artefacts des versions nocturnes de la v3 |
| `bump-webview2-v3.yml`, `release-webview2.yml` | Gestion des dépendances et des versions de WebView2 |
| `build-cross-image.yml` | Construit l’image de conteneur du compilateur croisé |
| `publish-npm.yml` | Publie sur npm le runtime JS intégré `@wailsio/runtime` |
| `pr-master.yml` | Vérifications des PR par rapport à la branche `master` |
| `semgrep.yml` | Analyse statique avec Semgrep |
| `stale-issues.yml`, `issue-labeler.yml`, `file-labeler.yml`, `claude.yml`, `generate-sponsor-image.yml`, `sync-translated-documents.yml`, `upload-source-documents.yml`, `build-and-test.yml`, `weekly-release-v2.yml` | Maintenance du dépôt et workflows propres à la v2 |

Cet arbre ne contient **aucun** `qodana.yaml` ni **aucun** `runtime.yml` : d’anciennes versions préliminaires de cette page mentionnaient les deux, mais seul `semgrep.yml` couvre l’analyse statique, et le paquet JS du runtime est publié par l’intermédiaire de `publish-npm.yml`.

Les étapes de CI correspondent aux cibles du Taskfile ci-dessus (`task test:cli`, `task test:generator`, `task test:templates`, `task test:examples`, …), ce qui vous permet de reproduire localement la CI à l’identique. L’étape de vérification rapide `wails3 build` dans `build-and-test-v3.yml` est appelée **sans option supplémentaire** : `wails3 build` ne possède aucune option `-skip-package`.

---

## 6. Reproduction locale de la CI

Il n’existe aucune cible globale `task ci` unique. Reproduisez la CI en enchaînant les véritables cibles :

```
task precommit
task test:cli
task test:generator
task test:templates
task test:examples
```

---

## 7. Résolution des échecs de tests

| Symptôme | Cause probable | Solution |
| --- | --- | --- |
| **Condition de concurrence dans `webview_window_darwin.go`** | Modification de l’état de la fenêtre hors du thread principal | Transmettez l’appel via `application.InvokeAsync` / `Invoke` afin qu’il s’exécute sur le thread principal |
| **Le test Linux se bloque dans une CI sans affichage** | GTK nécessite un affichage | Exécutez sous `xvfb-run`, par exemple `xvfb-run task test:examples:linux` |
| **La compilation du modèle échoue** | Le fichier de verrouillage du frontend n’est plus à jour | Réexécutez `wails3 init` sur un répertoire propre afin d’actualiser le modèle |
| **Erreurs de coverpkg** | Le test d’intégration importe `main` | Passez au tag de compilation `//go:build integration` et conditionnez l’importation |

---

## 8. Ajout de nouveaux tests

1. **Unitaire** — créez `*_test.go`, puis exécutez `go test ./...`
2. **Générateur / CLI** — complétez les cas sous `internal/generator/testcases/` ou `internal/commands/*_test.go`, puis réexécutez `task test:generator` / `task test:cli`
3. **Modèles / exemples** — vérifiez que les modèles distribués se compilent toujours avec `task test:templates`

---

## 9. Carte des fichiers principaux

| Élément | Chemin |
| --- | --- |
| Test aller-retour du générateur | `internal/generator/generate_test.go` |
| Test de génération des ressources | `internal/commands/build-assets_test.go` |
| Guide sur les conditions de concurrence et Cgo | `v3/TESTING.md` |
| Cibles de test du Taskfile | `v3/Taskfile.yaml` |
| Générateur de constantes d’événements | `v3/tasks/events/generate.go` |
| Workflow de CI | `.github/workflows/build-and-test-v3.yml` (Go 1.25 via `actions/setup-go@v5`) |
| Analyse statique | `.github/workflows/semgrep.yml` |
| Publication du runtime sur npm | `.github/workflows/publish-npm.yml` |

---

Dans Wails v3, la qualité n’est pas une considération secondaire. Grâce aux tests unitaires, aux suites de tests des générateurs et des modèles, à la détection des conditions de concurrence et à une matrice de CI multiplateforme, vous pouvez contribuer en toute confiance, sachant que vos modifications passent avec succès sur chaque système d’exploitation pris en charge. Bons tests !
