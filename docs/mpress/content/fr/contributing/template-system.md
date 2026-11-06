---
title: "Système de modèles"
description: "Comment Wails v3 génère la structure de nouveaux projets, comment les modèles sont organisés et comment créer les vôtres."
slug: "contributing/template-system"
sourcePath: "contributing/template-system.md"
---

Wails fournit un **système de modèles** qui permet à `wails3 init` de produire un projet prêt à l’emploi. Un ensemble volontairement restreint de frameworks dispose de modèles intégrés (Vanilla, React, Vue, Svelte) ; vous pouvez utiliser tout autre framework en [apportant votre propre frontend](/guides/dev/frontend-frameworks/) ou en publiant un [modèle personnalisé](/guides/advanced/custom-templates/).

Cette page aborde les sujets suivants :

1. Structure du répertoire des modèles
2. Comment la CLI sélectionne et génère les modèles
3. Création pas à pas d’un nouveau modèle
4. Mise à jour ou remplacement des modèles existants
5. Dépannage et bonnes pratiques

---

## 1. Emplacement des modèles

```
v3/internal/templates/
├── _common/        # Files copied into EVERY project (Taskfile.yml, build/, etc.)
├── base/           # Minimal frontend skeleton (frontend/ + NEXTSTEPS.md)
├── ios/            # iOS bootstrapper
├── vanilla/        vanilla-js/   # TypeScript (default) + JavaScript variant
├── react/          react-js/     # TypeScript (default) + JavaScript variant
├── vue/                          # TypeScript only
├── svelte/                       # TypeScript only
└── templates.go    # Registry + Install/Get APIs (no auto-registration via embed)
```

- **`_common/`** — structure de base universelle (Taskfile, répertoire `build/` et infrastructure partagée) fusionnée dans chaque projet.
- **`base/`** — squelette de frontend minimal destiné aux modèles personnalisés, et non couche backend automatique. Lors de la création normale d’un projet, le modèle du framework sélectionné est superposé à `_common/`.
- **Dossiers des frameworks** — contiennent le frontend (`frontend/`), la configuration du framework et un fichier `template.yaml` qui décrit les métadonnées du modèle.
- Les noms des dossiers correspondent à l’**identifiant du modèle** transmis à la CLI (`wails3 init -t react`).
- **Convention de nommage par langage :** TypeScript est utilisé par défaut et reçoit le nom sans suffixe (`react`) ; lorsqu’une variante JavaScript existe, son nom porte le suffixe `-js` (`react-js`). Les modèles intégrés déclarent explicitement leur langage avec `typescript: true|false` dans `template.yaml`. Les modèles de la communauté peuvent encore employer l’ancien suffixe `-ts`, qui reste pris en charge comme solution de repli.

> L’intégralité du répertoire `internal/templates/` est compilée dans le binaire de la CLI
>
> via `//go:embed *`, afin que les utilisateurs puissent générer la structure de projets hors ligne.

---

Le chargeur accepte également les anciennes métadonnées `template.json`. Les nouveaux modèles intégrés utilisent YAML, notamment pour définir explicitement le paramètre `typescript`.

## 2. Comment `wails3 init` utilise les modèles

Chaîne d’appels (sans `cmd/wails3/init.go` — la CLI est directement configurée dans `cmd/wails3/main.go`) :

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

Il n’existe aucune API `Template.Load()` / `Template.CopyTo()` / `Template.Validate()` : la création de la copie de travail est gérée par `gosod` (`github.com/leaanthony/gosod`) à partir de `fs.FS` intégré.

### Options de `wails3 init`

Définies dans `internal/flags/init.go` :

| Option | Rôle | Valeur par défaut |
| --- | --- | --- |
| `-p` | Nom du paquet | `main` |
| `-t` | Nom d’un modèle intégré, chemin local ou URL | `vanilla` |
| `-n` | Nom du projet | (vide) |
| `-d` | Répertoire du projet | `.` |
| `-q` | Masquer la sortie de la console | false |
| `-l` | Répertorier les modèles | false |
| `-skipgomodtidy` | Ne pas exécuter `go mod tidy` après la création de la copie de travail | false |
| `-git` | URL du dépôt Git à initialiser | (vide) |
| `-mod` | Chemin du module Go (dérivé de `-git` s’il n’est pas défini) | (vide) |
| `-s` | Ignorer l’avertissement lors de l’utilisation de modèles distants | false |
| `-productname` / `-productdescription` / `-productversion` / `-productcompany` / `-productcopyright` / `-productcomments` / `-productidentifier` | Métadonnées intégrées aux ressources générées lors de la compilation | valeurs par défaut raisonnables |

Il n’existe ni alias long `-list` (uniquement `-l`), ni option `--help` propre à chaque modèle.

### Substitutions

Les espaces réservés sont des directives de modèle Go standard — le caractère `.` initial fait partie de l’accès au champ :

| Espace réservé | Exemple | Source |
| --- | --- | --- |
| `{{.ProjectName}}` | `myapp` | Option `-n` / nom du répertoire |
| `{{.ModulePath}}` | `github.com/me/myapp` | Option `-mod` ou valeur dérivée de `-git` |
| `{{.WailsVersion}}` | `v3.0.0` | Constante intégrée à la compilation depuis `internal/version` |
| `{{.ProductName}}`, `{{.ProductDescription}}`, `{{.ProductVersion}}`, `{{.ProductCompany}}`, `{{.ProductCopyright}}`, `{{.ProductComments}}`, `{{.ProductIdentifier}}` | Métadonnées définies lors de la compilation | options `-product*` correspondantes |

Si vous avez besoin d’un nouvel espace réservé, ajoutez un champ aux données du modèle dans `internal/templates/templates.go`, ainsi qu’un champ ou une option correspondant dans `internal/flags/init.go` (ou définissez sa valeur depuis `internal/commands/init.go`).

### Hook après copie

Une fois que `gosod` a terminé de créer la copie de travail à partir du modèle, la CLI exécute :

```
go mod tidy
```

sauf si vous transmettez `-skipgomodtidy`. Il n’existe aucune étape `task deps`.

---

## 3. Création d’un nouveau modèle

> Exemple : ajouter un modèle **Solid**

### 3.1 Dossier et identifiant

```
internal/templates/solid/
```

Le nom du dossier correspond à l’identifiant du modèle. Utilisez le format **kebab-case**.

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

Commencez par copier `react`, puis supprimez les fichiers inutiles. N’oubliez pas de créer un `template.yaml` — définissez `typescript: true` pour un modèle TypeScript — `base/` est le seul dossier qui n’en possède pas.

### 3.3 Mettre à jour les espaces réservés

Recherchez et remplacez les valeurs d’exemple littérales par des directives de modèle Go, par exemple :

- `myapp` → `{{.ProjectName}}`
- `github.com/you/myapp` → `{{.ModulePath}}`

### 3.4 Connecter le modèle

Comme `templates.go` parcourt le système de fichiers intégré lors de l’initialisation, il suffit généralement d’ajouter un nouveau dossier sous `internal/templates/<id>/` : aucun appel d’enregistrement manuel n’est nécessaire. Si vous avez besoin d’une logique supplémentaire (validation personnalisée, étapes après la copie), ajoutez-la à `templates.Install` dans `internal/templates/templates.go`.

### 3.5 Tester

```bash
wails3 init -n demo -t solid
cd demo
wails3 dev
```

Vérifiez les points suivants :

- Le serveur de développement démarre sur le port indiqué dans `WAILS_VITE_PORT`
- Les liaisons générées apparaissent sous `frontend/bindings/...`
- Le rechargement à chaud fonctionne

---

## 4. Modifier des modèles existants

1. Modifiez les fichiers sous `internal/templates/<id>/`.
2. Recompilez la CLI (`cd v3 && go build -o ../wails3 ./cmd/wails3`) ; la directive `//go:embed *` prend en compte le nouveau contenu.
3. Mettez à niveau toutes les **versions des dépendances** dans `frontend/package.json` et `Taskfile.yml`.
4. Si le comportement change, mettez à jour la description `template.yaml` du modèle.

### Modifications courantes

| Tâche | Emplacement |
| --- | --- |
| Changer le port du serveur de développement | `frontend/vite.config.ts` — lire `WAILS_VITE_PORT` |
| Ajouter des variables d’environnement | `build/Taskfile.yml` ou `frontend/.env` |
| Remplacer le gestionnaire de paquets JS | Dans `build/Taskfile.yml`, remplacer `npm` par `pnpm`/`bun` |

---

## 5. Conseils pour créer des modèles

- **Conservez un frontend générique** — évitez de faire référence à des variables globales propres à Wails ; `/wails/runtime.js` est servi par le serveur de ressources à l’exécution.
- **N’incluez aucun artefact compilé** — excluez `node_modules`, `dist` et `.DS_Store` du répertoire intégré (ou ajoutez-les à `.gitignore` pour qu’ils ne soient jamais validés dans le dépôt).
- **Documentez les prérequis** — version de Node, outils CLI supplémentaires, etc., dans `template.yaml` ou `NEXTSTEPS.md`.
- **Évitez les changements incompatibles** — si la refonte est importante, créez un nouvel ID de modèle au lieu de modifier un modèle existant.

---

## 6. Résolution des problèmes

| Symptôme | Cause | Solution |
| --- | --- | --- |
| `unknown template name` | Faute de frappe dans `-t` ou modèle non intégré | Exécutez `wails3 init -l` pour afficher la liste des modèles disponibles |
| Les espaces réservés ne sont pas remplacés | `{{ProjectName}}` a été utilisé au lieu de `{{.ProjectName}}` | Ajoutez le `.` initial (accès à un champ de modèle Go) |
| Le serveur de développement ouvre une page vide | La configuration de Vite ne lit pas `WAILS_VITE_PORT` | Vérifiez votre `vite.config.ts` |
| La compilation du frontend échoue en production | Chemin Vite `base` oublié | Définissez `base: "./"` dans `vite.config.ts` |

---

## 7. Carte des principaux fichiers sources

| Fichier | Responsabilité |
| --- | --- |
| `internal/templates/templates.go` | Intègre le système de fichiers des modèles et expose `Install(options *flags.Init) error`, `GetDefaultTemplates()` et `ValidTemplateName(name)` |
| `internal/templates/<id>/**` | Contenu réel du modèle |
| `internal/commands/init.go` | Code de liaison de la CLI : sélectionne le modèle, renseigne les métadonnées et appelle `templates.Install` |
| `internal/commands/generate_template.go` | `wails3 generate template` — utilitaire permettant d’*exporter* un projet actif vers un modèle (pratique pour les mises à jour) |
| `internal/flags/init.go` | Définitions des options pour `wails3 init` |

---

## 8. Récapitulatif

- Les modèles se trouvent dans **`internal/templates/`** et sont intégrés à la CLI via `//go:embed *`.
- `wails3 init -t <id>` extrait le modèle via `gosod` et exécute `go mod tidy`, sauf si cette étape est ignorée avec `-skipgomodtidy`.
- Pour créer un modèle, il suffit d’**ajouter un dossier**, d’y placer les fichiers ainsi qu’un `template.yaml`, puis d’utiliser des espaces réservés de type `{{.ProjectName}}`.
- Le système est **extensible** et **autonome** — il est idéal pour partager des piles personnalisées avec votre équipe ou la communauté.

Bonne création de modèles !
