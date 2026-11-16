---
title: "Personnalisation de la compilation"
description: "Personnalisez votre processus de compilation avec Task et Taskfile.yml"
slug: "guides/build/customization"
sourcePath: "guides/build/customization.md"
---

## Présentation

Le système de compilation de Wails est un outil flexible et puissant conçu pour simplifier la compilation de vos applications Wails. Il s’appuie sur [Task](https://taskfile.dev), un exécuteur de tâches qui permet de définir et d’exécuter facilement des tâches. Bien que le système de compilation v3 soit utilisé par défaut, Wails encourage une approche consistant à « utiliser vos propres outils », afin que les développeurs puissent personnaliser leur processus de compilation selon leurs besoins.

Pour en savoir plus sur l’utilisation de Task, consultez la [documentation officielle](https://taskfile.dev/usage/).

## Task : au cœur du système de compilation

[Task](https://taskfile.dev) est une solution moderne écrite en Go qui remplace Make. Il utilise un fichier YAML pour définir les tâches et leurs dépendances. Dans le système de compilation de Wails, [Task](https://taskfile.dev) joue un rôle central dans l’orchestration du processus de compilation.

Le fichier `Taskfile.yml` principal se trouve à la racine du projet, tandis que les tâches propres à chaque plateforme sont définies dans des fichiers `build/<platform>/Taskfile.yml`. Un fichier `Taskfile.yml` commun, situé dans le répertoire `build`, contient les tâches partagées entre les plateformes.

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

Le fichier `Taskfile.yml` situé à la racine du projet constitue le point d’entrée principal du système de compilation. Il définit les tâches et leurs dépendances. Voici le fichier `Taskfile.yml` par défaut :

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

## Taskfiles propres aux plateformes

Chaque plateforme possède son propre Taskfile, situé dans l’un des répertoires de plateforme sous le répertoire `build`. Ces fichiers définissent les tâches principales de la plateforme correspondante. Chaque Taskfile inclut les tâches communes du fichier `build/Taskfile.yml`.

### Windows

Emplacement : `build/windows/Taskfile.yml`

Le Taskfile propre à Windows comprend des tâches permettant de compiler, de empaqueter et d’exécuter l’application sous Windows. Ses principales fonctionnalités sont les suivantes :

- Compilation avec des options de production facultatives
- Génération du fichier d’icône `.ico`
- Génération du fichier Windows `.syso`
- Création d’un programme d’installation NSIS pour l’empaquetage

### Linux

Emplacement : `build/linux/Taskfile.yml`

Le Taskfile propre à Linux comprend des tâches permettant de compiler, d’empaqueter et d’exécuter l’application sous Linux. Ses principales fonctionnalités sont les suivantes :

- Compilation avec des options de production facultatives
- Création de paquets AppImage, deb, rpm et Arch Linux
- Génération du fichier `.desktop` pour les applications Linux

### macOS

Emplacement : `build/darwin/Taskfile.yml`

Le Taskfile propre à macOS comprend des tâches permettant de compiler, d’empaqueter et d’exécuter l’application sous macOS. Ses principales fonctionnalités sont les suivantes :

- Compilation de fichiers binaires pour les architectures amd64, arm64 et universelle (les deux)
- Génération du fichier d’icône `.icns`
- Création d’un paquet `.app` destiné à la distribution
- Signature ad hoc des paquets `.app`
- Définition d’options de compilation et de variables d’environnement propres à macOS

## Exécution des tâches et alias de commandes

La commande `wails3 task` est une version intégrée de [Taskfile](https://taskfile.dev) qui exécute les tâches définies dans votre `Taskfile.yml`.

Les commandes `wails3 build` et `wails3 package` sont respectivement des alias de `wails3 task build` et `wails3 task package`. Lorsque vous exécutez ces commandes, Wails les traduit en interne en l’exécution de tâche correspondante :

- `wails3 build` → `wails3 task build`
- `wails3 package` → `wails3 task package`

### Transmission de paramètres aux tâches

Vous pouvez transmettre des variables de ligne de commande aux tâches au format `KEY=VALUE`. Ces variables sont transmises par l’intermédiaire des commandes alias :

```bash
# These are equivalent:
wails3 build GOOS=linux
wails3 task build GOOS=linux
```

Les tâches fournies par défaut utilisent `GOOS`, `ARCH` et `DEV`. Les variables telles que `VERSION`, `OUTPUT`, `PLATFORM` ou `CONFIG` n’ont d’effet sur une tâche que si votre propre Taskfile les lit. L’exemple suivant illustre un Taskfile personnalisé ; il ne s’agit pas de la tâche de compilation fournie par défaut avec Wails :

```yaml
tasks:
  build:
    cmds:
      - echo "Building for {{.PLATFORM | default "darwin"}}"
      - go build -tags {{.CONFIG | default "debug"}} -o myapp
```

Si votre tâche de création de paquet personnalisée lit `VERSION` et `OUTPUT`, vous pouvez les transmettre de la même manière :

```bash
wails3 package VERSION=2.0.0 OUTPUT=myapp.pkg
```

## Processus de compilation commun

Sur toutes les plateformes, le processus de compilation comprend généralement les étapes suivantes :

1. Nettoyage des modules Go
2. Compilation du frontend
3. Génération des icônes
4. Compilation du code Go avec des options propres à la plateforme
5. Empaquetage de l’application (propre à la plateforme)

## Personnalisation du processus de compilation

Bien que le système de compilation v3 fournisse une configuration par défaut robuste, vous pouvez facilement le personnaliser selon les besoins de votre projet. En modifiant le fichier `Taskfile.yml` et les Taskfiles propres aux plateformes, vous pouvez :

- Ajouter de nouvelles tâches
- Modifier les tâches existantes
- Changer l’ordre d’exécution des tâches
- Intégrer d’autres outils et scripts

Cette flexibilité vous permet d’adapter le processus de compilation à vos exigences particulières tout en bénéficiant de la structure fournie par le système de compilation de Wails.

@note{type="tip" title="Découvrir Taskfile"}
Nous vous recommandons vivement de lire la documentation de [Taskfile](https://taskfile.dev) afin de comprendre comment l’utiliser efficacement. Pour connaître la version de Taskfile intégrée à l’interface en ligne de commande de Wails, exécutez `wails3 task --version`.

@end

## Mode de développement

Le système de compilation de Wails comprend un puissant mode de développement qui améliore l’expérience des développeurs grâce au rechargement automatique et au remplacement à chaud des modules. Activez ce mode avec la commande `wails3 dev`.

### Fonctionnement

Lorsque vous exécutez `wails3 dev`, le processus suivant se déroule :

1. La commande recherche un port disponible et utilise 9245 par défaut si aucun port n’est indiqué.
2. Elle configure les variables d’environnement du serveur de développement frontend (Vite).
3. Elle démarre la surveillance des fichiers à l’aide de la bibliothèque [refresh](https://github.com/atterpac/refresh).

La bibliothèque [refresh](https://github.com/atterpac/refresh) surveille les modifications de fichiers et déclenche les recompilations. Elle utilise la configuration définie sous la clé `dev_mode` du fichier `./build/config.yml`. Vous pouvez la configurer pour ignorer certains répertoires et fichiers, déterminer les fichiers à surveiller et définir les actions à exécuter lorsque des modifications sont détectées. La configuration par défaut fonctionne très bien, mais n’hésitez pas à l’adapter à vos besoins.

### Configuration

Voici un exemple de sa structure :

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

Ce fichier de configuration vous permet de :

- Définir le chemin racine de la surveillance des fichiers
- Configurer le niveau de journalisation
- Définir un délai anti-rebond pour les événements de modification de fichiers
- Ignorer des répertoires, des fichiers ou des extensions de fichiers précis
- Définir les commandes à exécuter lors de modifications de fichiers

### Personnaliser le mode de développement

Vous pouvez personnaliser le mode de développement en modifiant ces valeurs dans le fichier `config.yml`.

Vous pouvez notamment le personnaliser comme suit :

1. Modifier les répertoires ou fichiers surveillés
2. Ajuster le délai d’anti-rebond pour contrôler la rapidité avec laquelle le système réagit aux modifications
3. Ajouter ou modifier les commandes d’exécution en fonction des besoins de votre projet

### Utiliser un navigateur pour le développement

Bien que Wails v2 ait entièrement pris en charge l’utilisation d’un navigateur pour le développement, cette possibilité a suscité beaucoup de confusion. Une application fonctionnant dans le navigateur ne fonctionnait pas nécessairement dans l’application de bureau, car toutes les API de navigateur ne sont pas disponibles dans les vues web.

Pour les travaux de développement axés sur l’interface utilisateur, vous pouvez toujours utiliser un navigateur dans v3 en accédant, en mode développement, à l’URL Vite disponible à l’adresse `http://localhost:9245`. Vous bénéficiez ainsi des puissants outils de développement du navigateur lorsque vous travaillez sur les styles et la mise en page. Sachez que les liaisons Go *ne fonctionneront pas* dans ce mode. Lorsque vous êtes prêt à tester des fonctionnalités telles que les liaisons et les événements, passez simplement à la vue de bureau pour vérifier que tout fonctionne parfaitement dans l’environnement de production.
