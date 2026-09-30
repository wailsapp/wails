---
title: "Système de modèles"
description: "Découvrez comment Wails v3 génère la structure de nouveaux projets, comment les modèles sont organisés et comment créer les vôtres."
slug: "contributing/template-system"
sourcePath: "contributing/template-system.md"
---

Wails fournit un **système de modèles** qui permet à `wails3 init` de produire un projet prêt à l’emploi. Un nombre volontairement restreint de frameworks disposent de modèles intégrés (Vanilla, React, Vue, Svelte) ; vous pouvez utiliser tout autre framework en [apportant votre propre frontend](/guides/dev/frontend-frameworks/) ou en publiant un [modèle personnalisé](/guides/advanced/custom-templates/).

Cette page aborde les sujets suivants :

1. Structure du répertoire des modèles
2. Sélection et rendu des modèles par la CLI
3. Création pas à pas d’un nouveau modèle
4. Mise à jour ou remplacement des modèles existants
5. Résolution des problèmes et bonnes pratiques

---

## 1. Emplacement des modèles

```
v3/internal/templates/
├── _common/        # Files copied into EVERY project (Taskfile.yml, build/, etc.)
├── base/           # Backend-only "plain Go" base layer (frontend/ + NEXTSTEPS.md)
├── ios/            # iOS bootstrapper
├── vanilla/        vanilla-js/   # TypeScript (default) + JavaScript variant
├── react/          react-js/     # TypeScript (default) + JavaScript variant
├── vue/                          # TypeScript only
├── svelte/                       # TypeScript only
└── templates.go    # Registry + Install/Get APIs (no auto-registration via embed)
```

- **`_common/`** — structure de base universelle (Taskfile, répertoire `build/` et infrastructure partagée), fusionnée dans chaque projet.
- **`base/`** — partie Go qui sert de point de départ à chaque modèle. Remarque : `base/` ne contient **pas** lui-même de `template.json` ; ce fichier se trouve dans chaque modèle propre à un framework.
- **Dossiers des frameworks** — contiennent le frontend (`frontend/`), la configuration du framework et un fichier `template.json` qui décrit les métadonnées du modèle.
- Les noms des dossiers correspondent à l’**identifiant du modèle** transmis à la CLI (`wails3 init -t react`).
- **Convention de nommage des langages :** TypeScript est utilisé par défaut et reprend le nom sans suffixe (`react`) ; lorsqu’une variante JavaScript existe, son nom reçoit le suffixe `-js` (`react-js`). Dans `template.yaml`, les modèles intégrés déclarent explicitement leur langage avec `typescript: true|false`. Les modèles de la communauté peuvent encore utiliser l’ancien suffixe `-ts`, qui reste pris en charge comme solution de repli.

> L’intégralité du répertoire `internal/templates/` est compilée dans le binaire de la CLI
>
> via `//go:embed *`, ce qui permet aux utilisateurs de générer la structure de projets hors ligne.

---

## 2. Utilisation des modèles par `wails3 init`

Chaîne d’appels (sans `cmd/wails3/init.go` — la CLI est directement raccordée dans `cmd/wails3/main.go`) :

```
cmd/wails3/main.go             (clir wiring)
       │
       ▼
internal/commands/init.go      Init(options *flags.Init) error
       │
       ▼
internal/templates/templates.go
       │   templates.Install(options)
       │   templates.GetDefaultTemplates()
       ▼
gosod.New(template.FS).Extract(options.ProjectDir, data)   // file extraction
       │
       ▼
go mod tidy (unless --skipgomodtidy / -skipgomodtidy)
```

Il n’existe aucune API `Template.Load()` / `Template.CopyTo()` / `Template.Validate()` : l’extraction est effectuée par `gosod` (`github.com/leaanthony/gosod`) à partir du `fs.FS` intégré.

### Options de `wails3 init`

Définies dans `internal/flags/init.go` :

| Option | Fonction | Valeur par défaut |
| --- | --- | --- |
| `-p` | Nom du paquet | `main` |
| `-t` | Nom d’un modèle intégré, chemin local ou URL | `vanilla` |
| `-n` | Nom du projet | (vide) |
| `-d` | Répertoire du projet | `.` |
| `-q` | Masquer la sortie de la console | false |
| `-l` | Répertorier les modèles | false |
| `-skipgomodtidy` | Ne pas exécuter `go mod tidy` après l’extraction | false |
| `-git` | URL du dépôt Git à initialiser | (vide) |
| `-mod` | Chemin du module Go (déduit de `-git` s’il n’est pas défini) | (vide) |
| `-s` | Ignorer l’avertissement lors de l’utilisation de modèles distants | false |
| `-productname` / `-productdescription` / `-productversion` / `-productcompany` / `-productcopyright` / `-productcomments` / `-productidentifier` | Métadonnées intégrées aux ressources de build générées | valeurs par défaut adaptées |

Il n’existe **aucun** alias long `-list` (uniquement `-l`), ni aucun `--help` propre à chaque modèle.

### Substitutions

Les espaces réservés sont des directives de modèle Go standard — le `.` initial fait partie de l’accesseur de champ :

| Espace réservé | Exemple | Source |
| --- | --- | --- |
| `{{.ProjectName}}` | `myapp` | Option `-n` / nom du répertoire |
| `{{.ModulePath}}` | `github.com/me/myapp` | Option `-mod` ou valeur dérivée de `-git` |
| `{{.WailsVersion}}` | `v3.0.0-…` | Constante intégrée à la compilation provenant de `internal/version` |
| `{{.ProductName}}`, `{{.ProductDescription}}`, `{{.ProductVersion}}`, `{{.ProductCompany}}`, `{{.ProductCopyright}}`, `{{.ProductComments}}`, `{{.ProductIdentifier}}` | Métadonnées définies lors de l’intégration | Options `-product*` correspondantes |

Si vous avez besoin d’un nouvel espace réservé, ajoutez un champ aux données du modèle dans  
`internal/templates/templates.go` et un champ ou une option correspondant dans  
`internal/flags/init.go` (ou définissez sa valeur à partir de `internal/commands/init.go`).

### Traitement après copie

Une fois que `gosod` a terminé d’extraire le modèle, la CLI exécute :

```
go mod tidy
```

sauf si vous transmettez `-skipgomodtidy`. Il n’existe aucune étape `task deps`.

---

## 3. Création d’un modèle

> Exemple : ajoutez un modèle **Solid**

### 3.1 Dossier et identifiant

```
internal/templates/solid/
```

Le nom du dossier est l’identifiant du modèle. Conservez-le au format **kebab-case**.

### 3.2 Ensemble minimal de fichiers

```
solid/
├── template.yaml    # name, description, wailsVersion, typescript (required)
├── frontend/        # Your web project (no node_modules/dist)
│   ├── src/
│   ├── package.json
│   └── vite.config.ts
└── ...              # Any extra Go files the template wants to inject
```

Commencez par copier `react`, puis supprimez les fichiers inutiles. N’oubliez pas de créer un  
`template.yaml` — définissez `typescript: true` pour un modèle TypeScript — `base/` est  
le seul dossier qui n’en contient pas.

### 3.3 Mise à jour des espaces réservés

Recherchez et remplacez les valeurs d’exemple littérales par des directives de modèle Go, par exemple :

- `myapp` → `{{.ProjectName}}`
- `github.com/you/myapp` → `{{.ModulePath}}`

### 3.4 Intégration

Comme `templates.go` parcourt le système de fichiers intégré lors de l’initialisation, il suffit généralement d’ajouter un nouveau  
dossier sous `internal/templates/<id>/` : aucun appel d’enregistrement  
manuel n’est nécessaire. Si vous avez besoin d’une logique supplémentaire (validation personnalisée,  
étapes après copie), ajoutez-la à `templates.Install` dans  
`internal/templates/templates.go`.

### 3.5 Test

```bash
wails3 init -n demo -t solid
cd demo
wails3 dev
```

Vérifiez les points suivants :

- Le serveur de développement démarre sur le port publié dans `WAILS_VITE_PORT`
- Les liaisons générées apparaissent sous `frontend/bindings/...`
- Le rechargement à chaud fonctionne

---

## 4. Modification des modèles existants

1. Modifiez les fichiers sous `internal/templates/<id>/`.
2. Recompilez la CLI (`cd v3 && go build -o ../wails3 ./cmd/wails3`) ; la directive `//go:embed *` prendra en compte le nouveau contenu.
3. Mettez à jour les **versions de dépendances** concernées dans `frontend/package.json` et `Taskfile.yml`.
4. Mettez à jour la description `template.json` du modèle si son comportement change.

### Modifications courantes

| Tâche | Emplacement |
| --- | --- |
| Changer le port du serveur de développement | `frontend/vite.config.ts` — lire `WAILS_VITE_PORT` |
| Ajouter des variables d’environnement | `build/Taskfile.yml` ou `frontend/.env` |
| Remplacer le gestionnaire de paquets JS | Dans `build/Taskfile.yml`, remplacer `npm` par `pnpm`/`bun` |

---

## 5. Conseils de création de modèles

- **Gardez le frontend générique** — évitez de faire référence à des variables globales propres à Wails ; `/wails/runtime.js` est servi par le serveur de ressources lors de l’exécution.
- **Aucun artefact compilé** — excluez `node_modules`, `dist` et `.DS_Store` du répertoire intégré (ou ajoutez-les à `.gitignore` afin qu’ils ne soient jamais validés dans le dépôt).
- **Documentez les prérequis** — version de Node, outils CLI supplémentaires, etc., dans `template.json` ou `NEXTSTEPS.md`.
- **Évitez les changements incompatibles** — si la refonte est importante, créez un nouvel identifiant de modèle au lieu de modifier un identifiant existant.

---

## 6. Résolution des problèmes

| Symptôme | Cause | Solution |
| --- | --- | --- |
| `unknown template name` | Faute de frappe dans `-t` ou modèle non intégré | Exécutez `wails3 init -l` pour afficher les modèles disponibles |
| Les espaces réservés ne sont pas remplacés | `{{ProjectName}}` a été utilisé à la place de `{{.ProjectName}}` | Ajoutez le caractère `.` initial (accès à un champ de modèle Go) |
| Le serveur de développement ouvre une page vide | La configuration de Vite ne lit pas `WAILS_VITE_PORT` | Vérifiez votre fichier `vite.config.ts` |
| La compilation du frontend échoue en production | Chemin Vite `base` oublié | Définissez `base: "./"` dans `vite.config.ts` |

---

## 7. Carte des principaux fichiers sources

| Fichier | Rôle |
| --- | --- |
| `internal/templates/templates.go` | Intègre le système de fichiers des modèles et expose `Install(options *flags.Init) error`, `GetDefaultTemplates()` et `ValidTemplateName(name)` |
| `internal/templates/<id>/**` | Contenu réel du modèle |
| `internal/commands/init.go` | Code d'intégration de la CLI : sélectionne le modèle, renseigne les métadonnées et appelle `templates.Install` |
| `internal/commands/generate_template.go` | `wails3 generate template` — utilitaire permettant d’*exporter* un projet actif vers un modèle (pratique pour les mises à jour) |
| `internal/flags/init.go` | Définitions des options de `wails3 init` |

---

## 8. Récapitulatif

- Les modèles se trouvent dans **`internal/templates/`** et sont intégrés à la CLI via `//go:embed *`.
- `wails3 init -t <id>` extrait le modèle via `gosod` et exécute `go mod tidy` (cette étape peut être ignorée avec `-skipgomodtidy`).
- Pour créer un modèle, il suffit de **créer un dossier**, d’y ajouter des fichiers ainsi qu’un fichier `template.json`, puis d’utiliser des espaces réservés de type `{{.ProjectName}}`.
- Le système est **extensible** et **autonome** — idéal pour partager des piles technologiques personnalisées avec votre équipe ou la communauté.

Bonne création de modèles !
