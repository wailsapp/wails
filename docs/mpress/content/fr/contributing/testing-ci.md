---
title: "Tests et intégration continue"
description: "Comment Wails v3 garantit la qualité grâce aux tests unitaires, aux suites d’intégration, à la détection des accès concurrents et à l’intégration continue avec GitHub Actions."
slug: "contributing/testing-ci"
sourcePath: "contributing/testing-ci.md"
---

Les frameworks de bureau robustes exigent des tests à toute épreuve. Wails v3 adopte une **stratégie à plusieurs niveaux** :

| Niveau | Objectif | Outils |
| --- | --- | --- |
| Tests unitaires | Retour rapide sur les fonctions isolées | `go test ./...` |
| Tests du générateur et de la CLI | Valider `wails3 generate bindings` et les mécanismes internes de la CLI | `task test:generator`, `task test:cli` |
| Tests des modèles | Vérifier que chaque modèle distribué peut toujours être compilé | `task test:templates` |
| Détection des accès concurrents | Détecter les accès concurrents aux données dans le runtime et le pont | `go test -race ./...` |
| Matrice d’intégration continue | Garantir la compatibilité entre systèmes d’exploitation pour chaque PR | GitHub Actions |

Ce document explique **où se trouvent les tests**, **comment les exécuter** et **ce que le Taskfile orchestre**.

`v3/TESTING.md` documente la compilation multiplateforme des exemples ; il ne s’agit pas d’un catalogue de suppressions des rapports d’accès concurrents.

---

## 1. Conventions relatives aux répertoires

```
v3/
├── internal/.../_test.go     # Unit tests for internal packages
├── pkg/.../_test.go          # Public API tests
├── tasks/events/generate.go  # Code generator for event constants (NOT a test harness)
├── tests/                    # Top-level integration test harness
└── TESTING.md                # Cross-platform example-build guidance
```

Consignes :

- **Placez les tests unitaires à côté du code** (`foo.go` ↔ `foo_test.go`).
- Pour les packages `pkg/`, utilisez le style **boîte noire** (`package application_test`) lorsqu’il améliore la qualité de l’API.
- Placez les fixtures partagées là où elles sont utilisées (cette arborescence ne comporte aucun package `internal/testutil/` central ; raccordez plutôt les fonctions auxiliaires package par package).

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
- Privilégiez les tests **pilotés par table** lorsque plusieurs entrées, cas limites ou résultats attendus exercent le même comportement. Donnez un nom descriptif à chaque cas.
- Si nécessaire, isolez les particularités propres aux plateformes derrière des balises de compilation (`foo_windows_test.go`, `foo_darwin_test.go`, …).

### Couverture attendue

La logique nouvelle ou modifiée doit atteindre une couverture de 100 % des instructions Go. Mesurez la couverture du package que vous avez modifié plutôt que de vous fier à un pourcentage couvrant l’ensemble du dépôt :

```bash
go test ./pkg/application -coverprofile=coverage.out
go tool cover -func=coverage.out
```

Certains chemins ne peuvent raisonnablement pas être testés dans un environnement de test normal, par exemple les défaillances propres à une plateforme, les comportements dépendant du matériel ou un mécanisme de repli défensif impossible à déclencher sans risque. Limitez strictement ces exceptions et expliquez chaque chemin non couvert dans la description de la PR.

### Exécution locale

```bash
cd v3
go test ./... -cover
```

Ou via le Taskfile (cibles réelles : il n’existe aucun raccourci `task test`) :

```
task test:cli            # CLI plumbing tests
task test:generator      # bindings generator round-trip tests
task test:templates      # build every shipped template
task test:infrastructure # supporting helpers
task test:examples       # exercise the example matrix (downloads as needed)
task test:all            # everything above
task sanity              # default GTK4 smoke check (legacy: sanity:gtk3)
task precommit           # what you should run before pushing
```

---

## 3. Tests d’intégration

`v3/tests/` héberge l’infrastructure de tests d’intégration couvrant plusieurs packages. Les cibles `test:example:*` et `test:examples:*` du Taskfile exécutent les vérifications de compilation sur darwin / windows / linux, notamment les matrices GTK3 / GTK4 basées sur Docker sous Linux.

> Les exemples exécutables se trouvent dans `v3/examples/`. Les cibles de test sélectionnent et compilent
>
> les exemples adaptés à la plateforme hôte ou à la matrice d’intégration continue.

Exécutez la suite de tests rapides de la plateforme hôte avec :

```
task test:examples       # host
task test:examples:all   # full matrix (slow)
```

---

## 4. Détection des accès concurrents

Les accès concurrents aux données sont rédhibitoires dans les runtimes d’interface graphique.

### Analyse d’un accès concurrent

Examinez les deux accès en conflit et reproduisez le problème sur la plateforme concernée. Ne considérez pas un rapport comme bénin et ne le supprimez pas sans preuve. `v3/TESTING.md` décrit la compilation des exemples ; la réussite de ces compilations ne prouve ni l’absence d’accès concurrents ni le comportement du runtime natif.

### Suite locale de détection des accès concurrents

```
go test -race ./...
```

> `wails3 dev` ne comporte aucune option `-race` ; ses options de CLI sont `--config`, `--port`,
>
> et `-s` (activation de HTTPS). Pour exécuter le runtime avec le détecteur d’accès concurrents,
>
> compilez une application de test avec `go build -race`, puis exécutez-la directement.

---

## 5. Workflows GitHub Actions

Fichiers de workflow réels dans `.github/workflows/` (vérifiés dans l’arborescence) :

| Fichier | Objectif |
| --- | --- |
| `build-and-test-v3.yml` | Matrice Go 1.26/1.27, vérifications du runtime et du générateur, compilation des modèles et des exemples, ainsi que vérifications de la compilation des applications. Linux utilise GTK4/WebKitGTK 6.0 par défaut ; la variante héritée définit `BUILD_TAGS=gtk3`. Le lancement natif et les interactions doivent être vérifiés séparément. |
| `cross-compile-test-v3.yml` | Vérifications élémentaires de la compilation croisée |
| `auto-changelog-v3.yml`, `changelog-v3.yml` | Automatisation du journal des modifications |
| `nightly-release-v3.yml` | Artefacts de version nocturne v3 |
| `bump-webview2-v3.yml`, `release-webview2.yml` | Gestion de la dépendance WebView2 et de ses versions |
| `build-cross-image.yml` | Compile l’image de conteneur du compilateur croisé |
| `publish-npm.yml` | Publie sur npm le runtime JS `@wailsio/runtime` intégré |
| `pr-master.yml` | Vérifications côté PR par rapport à la branche `master` |
| `semgrep.yml` | Analyse statique Semgrep |
| `stale-issues.yml`, `generate-sponsor-image.yml`, `build-and-test.yml`, `weekly-release-v2.yml` | Maintenance du dépôt et workflows v2 ; les anciens workflows de synchronisation et de téléversement Crowdin ont été abandonnés. |

Les étapes d’intégration continue reproduisent les cibles du Taskfile ci-dessus (`task test:cli`, `task test:generator`, `task test:templates`, `task test:examples`, …), ce qui vous permet de reproduire localement l’intégration continue à l’identique. Dans `build-and-test-v3.yml`, l’étape de test rapide `wails3 build` est appelée **sans option supplémentaire** : `wails3 build` ne comporte aucune option `-skip-package`.

---

## 6. Reproduction locale de l’intégration continue

Il n’existe pas de cible globale `task ci` unique. Reproduisez la CI en enchaînant les véritables cibles :

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
| **Condition de concurrence dans `webview_window_darwin.go`** | Modification de l’état de la fenêtre en dehors du thread principal | Faites transiter l’appel par `application.InvokeAsync` / `Invoke` afin qu’il s’exécute sur le thread principal |
| **Le test Linux se bloque dans la CI sans affichage** | GTK nécessite un serveur d’affichage | Exécutez le test sous `xvfb-run`, par exemple `xvfb-run task test:examples:linux` |
| **Échec de la compilation du modèle** | Le fichier de verrouillage du frontend n’est plus à jour | Réexécutez `wails3 init` dans un répertoire propre pour actualiser le modèle |
| **Erreurs de Coverpkg** | Le test d’intégration importe `main` | Passez au tag de compilation `//go:build integration` et conditionnez l’importation |

---

## 8. Ajout de nouveaux tests

1. **Test unitaire** — créez `*_test.go`, puis exécutez `go test ./...`
2. **Générateur / CLI** — étendez les cas sous `internal/generator/testcases/` ou `internal/commands/*_test.go`, puis réexécutez `task test:generator` / `task test:cli`
3. **Modèles / exemples** — vérifiez que les modèles distribués se compilent toujours avec `task test:templates`

---

## 9. Carte des fichiers clés

| Élément | Chemin |
| --- | --- |
| Test aller-retour du générateur | `internal/generator/generate_test.go` |
| Test de compilation des ressources | `internal/commands/build-assets_test.go` |
| Compilations d’exemples multiplateformes | `v3/TESTING.md` |
| Cibles de test du Taskfile | `v3/Taskfile.yaml` |
| Générateur de constantes d’événements | `v3/tasks/events/generate.go` |
| Workflow de CI | `.github/workflows/build-and-test-v3.yml` (1.26/1.27 Go via `actions/setup-go@v5`) |
| Analyse statique | `.github/workflows/semgrep.yml` |
| Publication du paquet npm du runtime | `.github/workflows/publish-npm.yml` |

---

Dans Wails v3, la qualité n’est pas une préoccupation secondaire. Grâce aux tests unitaires, aux suites de tests du générateur et des modèles, à la détection des conditions de concurrence et à une matrice de CI multiplateforme, vous pouvez contribuer en toute confiance, sachant que vos modifications passent tous les tests sur chaque système d’exploitation pris en charge. Bons tests !
