---
title: "Votre première application"
description: "Créez pas à pas votre première application de bureau Wails"
slug: "getting-started/your-first-app"
sourcePath: "getting-started/your-first-app.md"
---

Ce guide vous montre comment créer votre première application Wails v3, de la configuration du projet à sa compilation, en passant par le flux de travail de développement.

<br/>

<br/>

@steps
### Création d’un projet
Ouvrez votre terminal et exécutez la commande suivante pour créer un projet Wails :

```bash
wails3 init -n myfirstapp
```

Cette commande crée un répertoire nommé `myfirstapp` contenant tous les fichiers nécessaires.

   <video src="/assets/wails_init.mp4" controls></video>

### Exploration de la structure du projet
Accédez au répertoire `myfirstapp`. Vous y trouverez plusieurs fichiers et dossiers :

@filetree
- build/           Contient les fichiers utilisés par le processus de compilation
  - appicon.png  Icône de l’application
  - config.yml   Configuration de la compilation
  - Taskfile.yml Build tasks
  - darwin/      Fichiers de compilation propres à macOS
    - Info.dev.plist Development configuration
    - Info.plist    Configuration de production
    - Taskfile.yml  Tâches de compilation pour macOS
    - icons.icns    Icône de l’application macOS
  - linux/       Fichiers de compilation propres à Linux
    - Taskfile.yml  Tâches de compilation pour Linux
    - appimage/     Création du paquet AppImage
      - build.sh  Script de compilation AppImage
    - nfpm/        Création de paquets NFPM
      - nfpm.yaml Package configuration
      - scripts/  Scripts de compilation
  - windows/     Fichiers de compilation propres à Windows
    - Taskfile.yml        Tâches de compilation pour Windows
    - icon.ico           Icône de l’application Windows
    - info.json          Métadonnées de l’application
    - wails.exe.manifest Windows manifest file
    - nsis/              Fichiers du programme d’installation NSIS
      - project.nsi                    Fichier de projet NSIS
      - wails_tools.nsh               Scripts utilitaires NSIS
- frontend/        Fichiers de l’application frontend
  - index.html   Fichier HTML principal
  - main.js      Fichier JavaScript principal
  - package.json NPM package configuration
  - public/      Ressources statiques
  - Inter Font License.txt Font license
- .gitignore      Fichier d’exclusion Git
- README.md       Documentation du projet
- Taskfile.yml    Tâches du projet
- go.mod          Fichier du module Go
- go.sum          Sommes de contrôle du module Go
- greetservice.go Greeting service
- main.go         Code principal de l’application
@end

Prenez le temps d’explorer ces fichiers et de vous familiariser avec la structure.

@note{type="info"}
Bien que Wails v3 utilise [Task](https://taskfile.dev/) comme système de compilation par défaut, rien ne vous empêche d’utiliser `make` ou tout autre système de compilation.

@end

### Compilation de votre application
Pour compiler votre application, exécutez :

```bash
wails3 build
```

Cette commande compile une version de débogage de votre application et l’enregistre dans un nouveau répertoire `bin`.

@note{type="info"}
`wails3 build` est la forme abrégée de `wails3 task build` et exécute la tâche `build` dans `Taskfile.yml`.

@end

     <video src="/assets/wails_build.mp4" controls></video>

Une fois l’application compilée, vous pouvez l’exécuter comme n’importe quelle application classique :

@tabs{sync-key="platform"}
[Mac]
```sh
./bin/myfirstapp
```

[Windows]
```sh
bin\myfirstapp.exe
```

[Linux]
```sh
./bin/myfirstapp
```

@end

Une interface utilisateur simple s’affiche : elle constitue le point de départ de votre application. Comme il s’agit de la version de débogage, des journaux s’affichent également dans la fenêtre de la console. Ils sont utiles pour le débogage.

### Mode développement
Vous pouvez également exécuter l’application en mode développement. Ce mode vous permet de modifier le code de votre frontend et de voir immédiatement les changements dans l’application en cours d’exécution, sans avoir à recompiler toute l’application.

1. Ouvrez une nouvelle fenêtre de terminal.
2. Exécutez `wails3 dev`. L’application sera compilée et exécutée en mode débogage.
3. Ouvrez `frontend/index.html` dans l’éditeur de votre choix.
4. Modifiez le code en remplaçant `Please enter your name below` par `Please enter your name below!!!`.
5. Enregistrez le fichier.

Cette modification apparaîtra immédiatement dans votre application.

Toute modification du code backend déclenchera une recompilation :

1. Ouvrez `greetservice.go`.
2. Dans la ligne contenant `return "Hello " + name + "!"`, remplacez cette valeur par `return "Hello there " + name + "!"`.
3. Enregistrez le fichier.

L’application sera mise à jour en quelques secondes.

     <video src="/assets/wails_dev.mp4" controls></video>

### Création des paquets de votre application
Lorsque votre application est prête à être distribuée, vous pouvez créer des paquets propres à chaque plateforme :

@tabs{sync-key="platform"}
[Mac]
Pour créer un paquet `.app` :

```bash
wails3 package
```

Cette commande crée une version de production et la conditionne dans un paquet `.app` placé dans le répertoire `bin`.

[Windows]
Pour créer un programme d’installation NSIS :

```bash
wails3 package
```

Cette commande crée une version de production et la conditionne dans un programme d’installation NSIS placé dans le répertoire `bin`.

[Linux]
Wails prend en charge plusieurs formats de paquets pour la distribution sous Linux :

```bash
# Create all package types (AppImage, deb, rpm, and Arch Linux)
wails3 package

# Or create specific package types
wails3 task linux:create:appimage  # AppImage format
wails3 task linux:create:deb       # Debian package
wails3 task linux:create:rpm       # Red Hat package
wails3 task linux:create:aur       # Arch Linux package
```

@end

Pour obtenir des informations plus détaillées sur les options et la configuration de la création de paquets, consultez notre [guide de compilation et de création de paquets](/guides/build/building/).

### Configuration du contrôle de version et du nom du module
Votre projet est créé avec le nom de module temporaire `changeme`. Il est recommandé de le modifier pour qu’il corresponde à l’URL de votre dépôt :

1. Créez un dépôt sur GitHub ou sur l’hébergeur Git de votre choix
2. Initialisez Git dans le répertoire de votre projet :
  ```bash
  git init
  git add .
  git commit -m "Initial commit"
  ```

3. Définissez votre dépôt distant en remplaçant l’URL par celle de votre dépôt :
  ```bash
  git remote add origin https://github.com/username/myfirstapp.git
  ```

4. Dans `go.mod`, modifiez le nom de votre module pour qu’il corresponde à l’URL de votre dépôt :
  ```bash
  go mod edit -module github.com/username/myfirstapp
  ```

5. Envoyez votre code vers le dépôt distant :
  ```bash
  git push -u origin main
  ```


Cela garantit que le nom de votre module Go respecte les conventions de nommage des modules Go et facilite le partage de votre code.

@note{type="tip" title="Conseil de pro"}
Vous pouvez automatiser toutes les étapes d’initialisation à l’aide de l’option `-git` lors de la création de votre projet :

```bash
wails3 init -n myfirstapp -git github.com/username/myfirstapp
```

Cette option prend en charge différents formats d’URL Git :

- HTTPS : `https://github.com/username/project`
- SSH : `git@github.com:username/project` ou `ssh://git@github.com/username/project`
- Protocole Git : `git://github.com/username/project`
- Système de fichiers : `file:///path/to/project.git`

@end

@end

## Félicitations !

Vous venez de créer, de développer et de conditionner votre première application Wails. Ce n’est que le début de tout ce que vous pouvez accomplir avec Wails v3.

## Étapes suivantes

Si vous découvrez Wails, nous vous recommandons de consulter ensuite nos tutoriels, qui vous guideront de manière pratique à travers les différentes fonctionnalités de Wails. Le premier tutoriel est [Créer un service](/tutorials/01-creating-a-service/).

Si vous êtes un utilisateur plus expérimenté, consultez le [Guide de compilation et de création de paquets](/guides/build/building/) pour obtenir des informations plus détaillées sur l’utilisation de Wails.
