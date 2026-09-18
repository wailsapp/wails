---
title: "Personnalisation de la compilation"
description: "Personnalisez votre processus de compilation avec Task et Taskfile.yml"
slug: "guides/build/customization"
sourcePath: "guides/build/customization.md"
---

## Vue d’ensemble

Le système de compilation de Wails est un outil flexible et puissant conçu pour simplifier le processus de compilation de vos applications Wails. Il s’appuie sur [Task](https://taskfile.dev), un exécuteur de tâches qui vous permet de définir et d’exécuter facilement des tâches. Bien que le système de compilation v3 soit utilisé par défaut, Wails encourage une approche dans laquelle vous pouvez utiliser vos propres outils, ce qui permet aux développeurs de personnaliser leur processus de compilation selon leurs besoins.

Pour en savoir plus sur l’utilisation de Task, consultez la [documentation officielle](https://taskfile.dev/usage/).

## Task : au cœur du système de compilation

[Task](https://taskfile.dev) est une alternative moderne à Make, écrite en Go. Il utilise un fichier YAML pour définir les tâches et leurs dépendances. Dans le système de compilation de Wails, [Task](https://taskfile.dev) joue un rôle central dans l’orchestration du processus de compilation.

Le fichier `Taskfile.yml` principal se trouve à la racine du projet, tandis que les tâches propres à chaque plateforme sont définies dans des fichiers `build/<platform>/Taskfile.yml`. Un fichier commun `Taskfile.yml` dans le répertoire `build` contient les tâches communes partagées entre les plateformes.

@filetree

- Project Root
  - Taskfile.yml
  - build
    - windows/Taskfile.yml
    - darwin/Taskfile.yml
    - linux/Taskfile.yml
    - Taskfile.yml
@end

## Taskfile.yml

Le fichier `Taskfile.yml` situé à la racine du projet constitue le point d’entrée principal du système de compilation. Il définit les tâches et leurs dépendances. Voici le fichier `Taskfile.yml` par défaut :

```yaml
version: '3'

includes:
  common: ./build/Taskfile.yml
  windows: ./build/windows/Taskfile.yml
  darwin: ./build/darwin/Taskfile.yml
  linux: ./build/linux/Taskfile.yml

vars:
  APP_NAME: "myproject"
  BIN_DIR: "bin"
  VITE_PORT: '{{.WAILS_VITE_PORT | default 9245}}'

tasks:
  build:
    summary: Builds the application
    cmds:
      - task: "{{OS}}:build"

  package:
    summary: Packages a production build of the application
    cmds:
      - task: "{{OS}}:package"

  run:
    summary: Runs the application
    cmds:
      - task: "{{OS}}:run"

  dev:
    summary: Runs the application in development mode
    cmds:
      - wails3 dev -config ./build/config.yml -port {{.VITE_PORT}}


```

## Taskfiles propres à chaque plateforme

Chaque plateforme possède son propre Taskfile, situé dans les répertoires de plateforme sous le répertoire `build`. Ces fichiers définissent les tâches principales de cette plateforme. Chaque Taskfile inclut les tâches communes du fichier `build/Taskfile.yml`.

### Windows

Emplacement : `build/windows/Taskfile.yml`

Le Taskfile propre à Windows comprend des tâches pour compiler, empaqueter et exécuter l’application sous Windows. Ses principales fonctionnalités sont les suivantes :

- Compilation avec des options de production facultatives
- Génération du fichier d’icône `.ico`
- Génération du fichier Windows `.syso`
- Création d’un programme d’installation NSIS pour l’empaquetage

### Linux

Emplacement : `build/linux/Taskfile.yml`

Le Taskfile propre à Linux comprend des tâches pour compiler, empaqueter et exécuter l’application sous Linux. Ses principales fonctionnalités sont les suivantes :

- Compilation avec des options de production facultatives
- Création d’un AppImage et de paquets deb, rpm et Arch Linux
- Génération du fichier `.desktop` pour les applications Linux

### macOS

Emplacement : `build/darwin/Taskfile.yml`

Le Taskfile propre à macOS comprend des tâches pour compiler, empaqueter et exécuter l’application sous macOS. Ses principales fonctionnalités sont les suivantes :

- Compilation de fichiers binaires pour les architectures amd64, arm64 et universelle (les deux)
- Génération du fichier d’icône `.icns`
- Création d’un paquet `.app` pour la distribution
- Signature ad hoc des paquets `.app`
- Définition des options de compilation et des variables d’environnement propres à macOS

## Exécution des tâches et alias de commandes

La commande `wails3 task` est une version intégrée de [Taskfile](https://taskfile.dev) qui exécute les tâches définies dans votre fichier `Taskfile.yml`.

Les commandes `wails3 build` et `wails3 package` sont respectivement des alias de `wails3 task build` et `wails3 task package`. Lorsque vous exécutez ces commandes, Wails les traduit en interne en l’exécution de tâche appropriée :

- `wails3 build` → `wails3 task build`
- `wails3 package` → `wails3 task package`

### Transmission de paramètres aux tâches

Vous pouvez transmettre des variables de ligne de commande aux tâches au format `KEY=VALUE`. Ces variables sont transmises par l’intermédiaire des commandes alias :

```bash
# These are equivalent:
wails3 build PLATFORM=linux CONFIG=production
wails3 task build PLATFORM=linux CONFIG=production

# Package with custom version:
wails3 package VERSION=2.0.0 OUTPUT=myapp.pkg
```

Dans votre fichier `Taskfile.yml`, vous pouvez accéder à ces variables à l’aide de la syntaxe des modèles Go :

```yaml
tasks:
  build:
    cmds:
      - echo "Building for {{.PLATFORM | default "darwin"}}"
      - go build -tags {{.CONFIG | default "debug"}} -o myapp
```

### Variables de personnalisation de la compilation Go

Les projets générés proposent des variables additives pour les tags de compilation Go, les options de l’éditeur de liens et CGO. Les tags doivent être des noms séparés par des virgules, sans l’option `-tags` ; Wails les combine avec les tags requis par le mode de compilation et la plateforme choisis.

| Variable | Portée |
| --- | --- |
| `APP_TAGS` | Toutes les compilations pilotées par des Taskfiles |
| `APP_TAGS_LINUX` | Compilations Linux |
| `APP_TAGS_DARWIN` | Compilations macOS |
| `APP_TAGS_WINDOWS` | Compilations Windows |
| `APP_TAGS_ANDROID` | Compilations Android |
| `APP_TAGS_IOS` | Compilations iOS |
| `APP_TAGS_SERVER` | Compilations en mode serveur |
| `APP_LDFLAGS` | Options supplémentaires de l’éditeur de liens pour toutes les compilations pilotées par des Taskfiles |
| `APP_CGO_ENABLED` | Remplace le réglage CGO par `0` ou `1` pour les compilations de bureau et serveur |
| `EXTRA_TAGS` | Tags supplémentaires pour une seule invocation |

Par exemple :

```bash
wails3 build APP_TAGS=sqlite_fts5,netgo APP_TAGS_LINUX=myapp_linux
wails3 build EXTRA_TAGS=diagnostics
APP_LDFLAGS='-X example.com/myapp/internal/version.Value=1.2.3' wails3 build
```

Définissez une valeur `APP_*` dans le fichier racine `Taskfile.yml` pour en faire une valeur par défaut permanente du projet. Un argument `KEY=value` passé avec la commande a la priorité la plus élevée et remplace cette variable pour l’invocation. Une valeur littérale du Taskfile racine prime sur la variable d’environnement du processus portant le même nom. Si vous conservez l’expression autoréférente avec valeur par défaut présentée ci-dessous, l’environnement du processus est utilisé lorsqu’aucune valeur n’est fournie en ligne de commande. Utilisez `EXTRA_TAGS` pour ajouter des tags à une compilation sans remplacer la valeur permanente de `APP_TAGS`.

Lorsque `APP_CGO_ENABLED` est vide, Linux et macOS utilisent `1` par défaut, Windows utilise `0`, les compilations serveur natives conservent la valeur par défaut de Go sur l’hôte et les compilations serveur Docker utilisent `0`. Les compilations Android et iOS utilisent toujours CGO et conservent `CGO_ENABLED=1` ; `APP_CGO_ENABLED` ne remplace pas le réglage de ces chaînes d’outils mobiles. La compilation croisée de bureau et les compilations serveur via Docker reçoivent les mêmes valeurs `APP_*` applicables que leurs équivalents Taskfile natifs.

@note{type="info" title="Projets existants"}
Le fichier racine `Taskfile.yml` appartient au projet : `wails3 update build-assets` ne l’écrase donc pas. Les projets créés avant l’introduction de ces variables doivent ajouter manuellement les entrées suivantes à leur bloc `vars` racine :

```yaml
vars:
  APP_TAGS: '{{.APP_TAGS | default ""}}'
  APP_TAGS_LINUX: '{{.APP_TAGS_LINUX | default ""}}'
  APP_TAGS_DARWIN: '{{.APP_TAGS_DARWIN | default ""}}'
  APP_TAGS_WINDOWS: '{{.APP_TAGS_WINDOWS | default ""}}'
  APP_TAGS_ANDROID: '{{.APP_TAGS_ANDROID | default ""}}'
  APP_TAGS_IOS: '{{.APP_TAGS_IOS | default ""}}'
  APP_TAGS_SERVER: '{{.APP_TAGS_SERVER | default ""}}'
  APP_LDFLAGS: '{{.APP_LDFLAGS | default ""}}'
  APP_CGO_ENABLED: '{{.APP_CGO_ENABLED | default ""}}'
```

Remplacez une expression autoréférente avec valeur par défaut par une valeur littérale pour en faire une valeur par défaut permanente du projet.

@end

Ces variables sont utilisées par les compilations pilotées par des Taskfiles. La phase de compilation du projet Xcode généré appelle Go directement et ne passe pas par cette personnalisation des Taskfiles ; les compilations lancées depuis Xcode n’utilisent donc pas ces variables actuellement.

## Processus de compilation commun

Sur toutes les plateformes, le processus de compilation comprend généralement les étapes suivantes :

1. Nettoyage des modules Go
2. Compilation du frontend
3. Génération des icônes
4. Compilation du code Go avec des options propres à la plateforme
5. Empaquetage de l’application (propre à la plateforme)

## Personnalisation du processus de compilation

Bien que le système de compilation v3 fournisse une configuration par défaut robuste, vous pouvez facilement le personnaliser en fonction des besoins de votre projet. En modifiant le fichier `Taskfile.yml` et les Taskfiles propres à chaque plateforme, vous pouvez :

- Ajouter de nouvelles tâches
- Modifier les tâches existantes
- Modifier l’ordre d’exécution des tâches
- Intégrer d’autres outils et scripts

Cette flexibilité vous permet d’adapter le processus de compilation à vos exigences particulières, tout en bénéficiant de la structure fournie par le système de compilation de Wails.

@note{type="tip" title="Découvrir Taskfile"}
Nous vous recommandons vivement de lire la documentation de [Taskfile](https://taskfile.dev) afin de comprendre comment l’utiliser efficacement. Pour connaître la version de Taskfile intégrée à la CLI Wails, exécutez `wails3 task --version`.

@end

## Mode de développement

Le système de compilation de Wails comprend un puissant mode de développement qui améliore l’expérience des développeurs grâce au rechargement en direct et au remplacement à chaud des modules. Ce mode est activé à l’aide de la commande `wails3 dev`.

### Fonctionnement

Lorsque vous exécutez `wails3 dev`, le processus suivant se déroule :

1. La commande recherche un port disponible et utilise 9245 par défaut si aucun port n’est spécifié.
2. Elle configure les variables d’environnement du serveur de développement frontend (Vite).
3. Elle démarre la surveillance des fichiers à l’aide de la bibliothèque [refresh](https://github.com/atterpac/refresh).

La bibliothèque [refresh](https://github.com/atterpac/refresh) surveille les modifications de fichiers et déclenche les reconstructions. Elle utilise la configuration définie sous la clé `dev_mode` du fichier `./build/config.yml`.  
Vous pouvez la configurer pour ignorer certains répertoires et fichiers, déterminer les fichiers à surveiller et définir les actions à effectuer lorsque des modifications sont détectées.  
La configuration par défaut fonctionne plutôt bien, mais n’hésitez pas à l’adapter à vos besoins.

### Configuration

Voici un exemple de sa structure :

```yaml
dev_mode:
  root_path: .
  log_level: warn
  debounce: 1000
  ignore:
    dir:
      - .git
      - node_modules
      - frontend
      - bin
    file:
      - .DS_Store
      - .gitignore
      - .gitkeep
    watched_extension:
      - "*.go"
    git_ignore: true
  executes:
    - cmd: wails3 task common:install:frontend:deps
      type: once
    - cmd: wails3 task common:dev:frontend
      type: background
    - cmd: go mod tidy
      type: blocking
    - cmd: wails3 task build
      type: blocking
    - cmd: wails3 task run
      type: primary
```

Ce fichier de configuration vous permet d’effectuer les opérations suivantes :

- Définir le chemin racine pour la surveillance des fichiers
- Configurer le niveau de journalisation
- Définir un délai anti-rebond pour les événements de modification de fichiers
- Ignorer des répertoires, des fichiers ou des extensions de fichier spécifiques
- Définir les commandes à exécuter lors de modifications de fichiers

### Personnalisation du mode de développement

Vous pouvez personnaliser le mode de développement en modifiant ces valeurs dans le fichier `config.yml`.

Vous pouvez notamment :

1. Modifier les répertoires ou fichiers surveillés
2. Ajuster le délai anti-rebond pour contrôler la rapidité avec laquelle le système réagit aux modifications
3. Ajouter ou modifier les commandes à exécuter en fonction des besoins de votre projet

### Utilisation d’un navigateur pour le développement

Bien que Wails v2 prenne entièrement en charge l’utilisation d’un navigateur pour le développement, cela entraînait beaucoup de confusion. Une application qui fonctionnait dans le navigateur ne fonctionnait pas nécessairement dans l’application de bureau, car les vues web ne donnent pas accès à toutes les API de navigateur.

Pour les travaux de développement axés sur l’interface utilisateur, vous pouvez toujours utiliser un navigateur dans la v3 en accédant à l’URL de Vite à l’adresse `http://localhost:9245` en mode de développement. Vous disposez ainsi des puissants outils de développement du navigateur lorsque vous travaillez sur les styles et la mise en page. Sachez que les liaisons Go *ne fonctionneront pas* dans ce mode.  
Lorsque vous êtes prêt à tester des fonctionnalités telles que les liaisons et les événements, passez simplement à la vue de bureau afin de vérifier que tout fonctionne parfaitement dans l’environnement de production.
