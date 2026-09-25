---
title: "Création de modèles personnalisés"
description: "Comment générer, personnaliser et héberger vos propres modèles de projet Wails v3"
slug: "guides/advanced/custom-templates"
sourcePath: "guides/advanced/custom-templates.md"
---

Wails est fourni avec un ensemble de modèles intégrés, mais vous pouvez créer les vôtres et les partager avec la communauté. Un modèle personnalisé est simplement un dépôt Git : une fois hébergé publiquement, il permet à quiconque de générer la structure initiale d’un projet à l’aide d’une seule commande.

## Générer la structure initiale d’un modèle

La commande `wails3 generate template` génère un répertoire de modèle prêt à être personnalisé :

```bash
wails3 generate template -name MyTemplate
```

Toutes les options :

| Option | Description | Valeur par défaut |
| --- | --- | --- |
| `-name` | Nom du modèle (obligatoire) | — |
| `-author` | Nom de l’auteur | — |
| `-description` | Brève description affichée dans l’interface en ligne de commande | — |
| `-helpurl` | URL de la documentation de ce modèle | — |
| `-version` | Version initiale | `v0.0.1` |
| `-frontend` | Copier un répertoire frontend existant dans le modèle | — |
| `-dir` | Emplacement où créer le répertoire du modèle | Répertoire courant |

Exemple avec toutes les options :

```bash
wails3 generate template \
  -name "My Template" \
  -author "Your Name" \
  -description "React + custom setup" \
  -helpurl "https://github.com/yourname/my-template" \
  -version "v1.0.0" \
  -frontend ./my-existing-frontend
```

Le répertoire généré se présente comme suit :

```
MyTemplate/
├── template.yaml          # Template metadata — edit this
├── NEXTSTEPS.md           # Guidance for you as the template author — delete before publishing
├── README.md              # Shown to users after they create a project
├── main.go.tmpl           # Application entry point
├── greetservice.go        # Example Go service
├── go.mod.tmpl            # Go module file
├── go.sum.tmpl            # Go checksums
├── gitignore.tmpl         # Becomes .gitignore in generated projects
├── Taskfile.tmpl.yml      # Build task definitions
└── frontend/              # Your frontend code
```

@note{type="tip" title="Lire NEXTSTEPS.md"}
Le fichier `NEXTSTEPS.md` généré contient des instructions détaillées sur chaque partie du modèle. Lisez-le avant toute personnalisation. Supprimez-le avant la publication : il ne doit pas figurer dans les projets créés à partir de votre modèle.

@end

## Configurer les métadonnées du modèle

Ouvrez `template.yaml` pour définir les métadonnées de votre modèle :

```yaml
# yaml-language-server: $schema=https://v3.wails.io/schemas/template.v3.json
name: "My Template"
shortname: my-template
author: Your Name
description: A template with my preferred setup
helpurl: https://github.com/yourname/my-template
version: v1.0.0
wailsVersion: 3
```

Le champ `wailsVersion` est **obligatoire** et doit avoir la valeur `3`. Le commentaire `# yaml-language-server` situé en haut active l’autocomplétion et la validation à la volée dans VS Code (avec l’[extension YAML](https://marketplace.visualstudio.com/items?itemName=redhat.vscode-yaml)) et les IDE JetBrains. Vous pouvez le conserver ou le supprimer ; il n’a aucun effet à l’exécution.

## Personnaliser le modèle

### Frontend

Le répertoire `frontend/` est copié à l’identique dans chaque projet créé à partir de votre modèle. Remplacez le contenu temporaire par votre véritable frontend :

@tabs
[Partir de zéro]
```bash
cd MyTemplate/frontend
npm create vite@latest .
```

Suivez les invites, puis installez les dépendances :

```bash
npm install
```

[Utiliser un projet existant]
Lors de la génération du modèle, transmettez `-frontend` pour copier un frontend existant en une seule étape :

```bash
wails3 generate template -name MyTemplate -frontend ./my-app/frontend
```

Vous pouvez aussi le copier manuellement dans le répertoire `frontend/` par la suite.

@end

### Tâches de compilation

`Taskfile.tmpl.yml` définit le processus de compilation. Adaptez les tâches `install:frontend:deps` et `build:frontend` à votre chaîne d’outils frontend :

```yaml
tasks:
  install:frontend:deps:
    dir: frontend
    cmds:
      - npm install       # replace with pnpm install, yarn, etc.

  build:frontend:
    dir: frontend
    deps: [install:frontend:deps, generate:bindings]
    cmds:
      - npm run build     # replace with your build command
```

### Application Go

Le fichier `main.go.tmpl` constitue le point d’entrée de l’application. Le moteur de modèles de Wails le traite lors de la création d’un projet : les variables de modèle telles que `{{.ProductName}}` sont remplacées par les valeurs fournies par l’utilisateur.

Pour le modifier comme un véritable fichier Go, avec la prise en charge de l’IDE, renommez-le temporairement en `main.go`, apportez vos modifications, puis rétablissez son nom `main.go.tmpl` avant d’effectuer le commit.

#### Variables de modèle

Ces variables sont disponibles dans tout fichier `.tmpl` :

| Variable | Description | Exemple |
| --- | --- | --- |
| `{{.ProjectName}}` | Nom du projet fourni par l’utilisateur | `"MyApp"` |
| `{{.BinaryName}}` | Nom du fichier binaire | `"myapp"` |
| `{{.ProductName}}` | Nom d’affichage du produit | `"My Application"` |
| `{{.ProductDescription}}` | Description du produit | `"An awesome application"` |
| `{{.ProductVersion}}` | Version du produit | `"1.0.0"` |
| `{{.ProductCompany}}` | Nom de l’entreprise ou de l’auteur | `"My Company Ltd"` |
| `{{.ProductCopyright}}` | Mention de copyright | `"Copyright 2024 My Company Ltd"` |
| `{{.ProductComments}}` | Commentaires supplémentaires sur le produit | `"Built with Wails"` |
| `{{.ProductIdentifier}}` | Identifiant du produit au format DNS inversé | `"com.mycompany.myapp"` |
| `{{.ModulePath}}` | Chemin du module Go | `"github.com/you/myapp"` |
| `{{.WailsVersion}}` | Version de Wails utilisée pour créer le projet | `"3.0.0"` |
| `{{.Typescript}}` | `true` si le nom du modèle se termine par `-ts` | `true` |
| `{{.Opn}}` | `{{` littéral — à échapper dans les modèles | `{{` |
| `{{.Cls}}` | `}}` littéral — à échapper dans les modèles | `}}` |

@note{type="tip"}
N’importe quel fichier de votre modèle peut être un fichier `.tmpl`, y compris les fichiers HTML, JSON et YAML. Les fichiers sans le suffixe `.tmpl` sont copiés à l’identique.

@end

## Tester votre modèle localement

Avant de le publier, testez le modèle en créant un projet à partir d’un chemin local :

```bash
wails3 init -n testproject -t /path/to/MyTemplate
```

Vérifiez ensuite que le projet fonctionne :

```bash
cd testproject
wails3 dev    # development mode with hot reload
wails3 build  # production binary
```

Vérifiez les points suivants :

- Le rechargement à chaud du frontend fonctionne
- Les modifications du code Go entraînent la recompilation et le redémarrage de l’application
- Le binaire de production situé dans `bin/` s’exécute correctement

## Publier sur GitHub

@steps
### **Créez un dépôt GitHub public** pour votre modèle. La racine du dépôt doit contenir `template.yaml`.
### **Supprimez `NEXTSTEPS.md`** — ce fichier fournit des conseils aux auteurs de modèles et ne doit pas figurer dans les projets que les utilisateurs créent à partir de votre modèle.
### **Validez et poussez** le contenu du répertoire du modèle à la racine du dépôt :
```bash
git init
git add .
git commit -m "Initial template"
git remote add origin https://github.com/yourname/my-template.git
git push -u origin main
```

### **Créez une étiquette de version** en respectant la gestion sémantique des versions :
```bash
git tag v1.0.0
git push origin v1.0.0
```

@end

Les utilisateurs peuvent désormais créer des projets à partir de votre modèle :

```bash
# Latest commit on the default branch
wails3 init -n myapp -t https://github.com/yourname/my-template

# Pinned to a specific release tag
wails3 init -n myapp -t https://github.com/yourname/my-template@v1.0.0
```

@note{type="caution" title="Avertissement concernant les modèles tiers"}
Lorsqu’un utilisateur installe un modèle distant, Wails affiche un avertissement indiquant que le modèle est du code tiers et que le projet Wails décline toute responsabilité quant à son contenu. L’utilisateur doit confirmer explicitement avant la création du projet.

En tant qu’auteur du modèle, vous êtes responsable de la sécurité et du bon fonctionnement de tout le code qu’il contient.

@end

## Bonnes pratiques

- **Rédigez un `README.md`** clair — il est présenté aux utilisateurs après la création d’un projet. Expliquez comment exécuter, compiler et personnaliser le projet.
- **Renseignez `helpurl`** — ajoutez un lien vers votre dépôt ou une documentation dédiée. Les utilisateurs le voient dans la liste des modèles de la CLI Wails.
- **Figez les versions des dépendances du frontend** dans `package.json` afin d’éviter que des mises à jour en amont ne rendent les installations impossibles.
- **Testez avant de créer l’étiquette de version** — créez un nouveau projet à partir de la version étiquetée avant de l’annoncer à la communauté.
- **Conservez `wailsVersion: 3`** — ce champ indique à Wails la version majeure ciblée par le modèle. Ne le modifiez pas.
- **Effectuez des mises à jour régulières** — maintenez les dépendances à jour et effectuez des tests avec les nouvelles versions de Wails.
