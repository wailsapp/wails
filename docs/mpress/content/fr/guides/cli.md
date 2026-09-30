---
title: "Référence de la CLI"
description: "Référence complète des commandes de la CLI Wails"
slug: "guides/cli"
sourcePath: "guides/cli.md"
---

La CLI Wails fournit un ensemble complet de commandes pour vous aider à développer, compiler et maintenir vos applications Wails.

## Commandes principales

Les commandes principales servent à créer, développer et compiler des projets.

Toutes les commandes de la CLI respectent le format suivant : `wails3 <command>`.

### `init`

Initialise un nouveau projet Wails. Pendant cette initialisation, la commande `go mod tidy` est exécutée pour mettre à jour les paquets du projet. Vous pouvez ignorer cette étape en utilisant l’option `-skipgomodtidy` avec la commande `init`.

```bash
wails3 init [flags]
```

#### Options

| Option | Description | Valeur par défaut |
| --- | --- | --- |
| `-p` | Nom du paquet Go | `main` |
| `-t` | Nom ou URL du modèle | `vanilla` |
| `-n` | Nom du projet |  |
| `-d` | Répertoire du projet | `.` |
| `-q` | Masquer la sortie | `false` |
| `-l` | Afficher la liste des modèles | `false` |
| `-mod` | Chemin du module Go (calculé à partir de `-git` s’il est omis) |  |
| `-git` | URL du dépôt Git |  |
| `-s` | Ignorer l’avertissement lors de l’utilisation d’un modèle distant | `false` |
| `-productname` | Nom du produit | `My Product` |
| `-productdescription` | Description du produit | `My Product Description` |
| `-productversion` | Version du produit | `0.1.0` |
| `-productcompany` | Nom de l’entreprise | `My Company` |
| `-productcopyright` | Mention de droits d’auteur | `© now, My Company` |
| `-productcomments` | Commentaires de fichier | `This is a comment` |
| `-productidentifier` | Identifiant du produit |  |
| `-skipgomodtidy` | Ne pas exécuter go mod tidy | `false` |

L’option `-git` accepte différents formats d’URL Git :

- HTTPS : `https://github.com/username/project`
- SSH : `git@github.com:username/project` ou `ssh://git@github.com/username/project`
- Protocole Git : `git://github.com/username/project`
- Système de fichiers : `file:///path/to/project.git`

Lorsque cette option est fournie, elle effectue les opérations suivantes :

1. Initialiser un dépôt Git dans le répertoire du projet
2. Définir l’URL spécifiée comme dépôt distant origin
3. Mettre à jour le nom du module dans `go.mod` pour qu’il corresponde à l’URL du dépôt
4. Ajouter tous les fichiers

### `dev`

Exécute l’application en mode développement. Vous disposez ainsi d’un aperçu en direct de votre code frontend : vous pouvez le modifier et voir les changements apparaître dans l’application en cours d’exécution sans avoir à recompiler l’ensemble de l’application. Les modifications apportées à votre code Go sont également détectées, puis l’application est automatiquement recompilée et relancée.

```bash
wails3 dev [flags]
```

#### Options

| Option | Description | Valeur par défaut |
| --- | --- | --- |
| `-config` | Chemin du fichier de configuration | `./build/config.yml` |
| `-port` | Port du serveur de développement Vite | `9245` |
| `-s` | Activer HTTPS | `false` |

@note{type="info"}
Cela revient à exécuter `wails3 task dev` et lance la tâche `dev` définie dans le Taskfile principal du projet. Vous pouvez personnaliser ce comportement en modifiant le fichier `Taskfile.yml`.

@end

### `build`

Compile une version de débogage de votre application. Par défaut, la compilation cible la plateforme et l’architecture actuelles.

```bash
wails3 build [flags] [CLI variables...]
```

#### Options

| Option | Description | Valeur par défaut |
| --- | --- | --- |
| `-tags` | Balises de compilation Go supplémentaires (séparées par des virgules) |  |

Vous pouvez transmettre des variables CLI pour personnaliser la compilation :

```bash
wails3 build PLATFORM=linux CONFIG=production
```

Utilisez l’option `-tags` pour transmettre des balises de compilation Go personnalisées :

```bash
# Build with legacy GTK3 + WebKit2GTK 4.1 on Linux (default is GTK4 + WebKitGTK 6.0)
wails3 build -tags gtk3

# Build in server mode (no GUI)
wails3 build -tags server

# Multiple tags
wails3 build -tags gtk3,customtag
```

Les balises sont transmises au Taskfile sous-jacent sous la forme `EXTRA_TAGS`.

@note{type="info"}
Cela équivaut à exécuter `wails3 task build`, qui lance la tâche `build` dans le Taskfile principal du projet. Toutes les variables CLI transmises à `build` sont transférées à la tâche sous-jacente. Vous pouvez personnaliser le processus de compilation en modifiant le fichier `Taskfile.yml`.

@end

### `package`

Crée des paquets propres à chaque plateforme en vue de leur distribution.

```bash
wails3 package [CLI variables...]
```

Vous pouvez transmettre des variables CLI pour personnaliser la création des paquets :

```bash
wails3 package VERSION=2.0.0 OUTPUT=myapp.pkg
```

#### Types de paquets

Les types de paquets suivants sont disponibles pour chaque plateforme :

| Plateforme | Type de paquet |
| --- | --- |
| Windows | `.exe` |
| macOS | `.app`, |
| Linux | `.AppImage`, `.deb`, `.rpm`, `.archlinux` |

@note{type="info"}
Cela équivaut à `wails3 task package`, qui lance la tâche `package` dans le Taskfile principal du projet. Toutes les variables CLI transmises à `package` sont transférées à la tâche sous-jacente. Vous pouvez personnaliser le processus de création des paquets en modifiant le fichier `Taskfile.yml`.

@end

### `task`

Exécute les tâches définies dans le fichier Taskfile.yml de votre projet. Il s’agit d’une version intégrée de [Taskfile](https://taskfile.dev) qui vous permet de définir et d’exécuter des tâches personnalisées de compilation, de test et de déploiement.

```bash
wails3 task [taskname] [CLI variables...] [flags]
```

#### Variables CLI

Vous pouvez transmettre des variables aux tâches au format `KEY=VALUE` :

```bash
wails3 task build PLATFORM=linux CONFIG=production
wails3 task deploy ENV=staging VERSION=1.2.3
```

Ces variables sont accessibles dans votre fichier Taskfile.yml à l’aide de la syntaxe des modèles Go :

```yaml
tasks:
  build:
    cmds:
      - echo "Building for {{.PLATFORM | default "darwin"}}"
      - echo "Config: {{.CONFIG | default "debug"}}"
```

#### Options

| Option | Description | Valeur par défaut |
| --- | --- | --- |
| `-h` | Affiche les informations d’utilisation de Task | `false` |
| `-i` | Crée un nouveau fichier Taskfile.yml | `false` |
| `-list` | Répertorie les tâches accompagnées de leur description | `false` |
| `-list-all` | Répertorie toutes les tâches (avec ou sans description) | `false` |
| `-json` | Formate la liste des tâches en JSON | `false` |
| `-status` | Renvoie un code de sortie non nul si la tâche n’est pas à jour | `false` |
| `-f` | Force l’exécution même lorsque la tâche est à jour | `false` |
| `-w` | Active le mode de surveillance pour la tâche indiquée | `false` |
| `-v` | Active le mode détaillé | `false` |
| `-version` | Affiche la version de Task | `false` |
| `-s` | Désactive l’affichage des commandes exécutées | `false` |
| `-p` | Exécute les tâches en parallèle | `false` |
| `-dry` | Compile et affiche les tâches sans les exécuter | `false` |
| `-summary` | Affiche un résumé d’une tâche | `false` |
| `-x` | Propage le code de sortie de la tâche | `false` |
| `-dir` | Définit le répertoire d’exécution |  |
| `-taskfile` | Choisit le Taskfile à exécuter |  |
| `-output` | Définit le style de sortie : [interleaved|group|prefixed] |  |
| `-c` | Sortie en couleur (activée par défaut) | `true` |
| `-C` | Limite le nombre de tâches exécutées simultanément |  |
| `-interval` | Intervalle de surveillance des modifications (en secondes) |  |

#### Exemples

```bash
# Run the default task
wails3 task

# Run a specific task
wails3 task test

# Run a task with variables
wails3 task build PLATFORM=windows ARCH=amd64

# List all available tasks
wails3 task --list

# Run multiple tasks in parallel
wails3 task -p task1 task2 task3

# Watch for changes and re-run task
wails3 task -w dev
```

### `mcp`

Démarre le serveur MCP du projet Wails pour la gestion du projet assistée par un agent. Celui-ci est distinct du serveur MCP compilé dans une application en cours d’exécution : `wails3 mcp` gère les fichiers du projet et les commandes de son cycle de vie, tandis que le serveur MCP de l’application contrôle la WebView en cours d’exécution.

```bash
wails3 mcp [flags]
```

Le transport est sélectionné automatiquement :

- Lorsqu’un hôte MCP lance Wails avec l’entrée et la sortie standard redirigées par des tubes, le serveur utilise **stdio**.
- Lorsqu’il est exécuté de manière interactive dans un terminal, le serveur utilise **Streamable HTTP** sur `127.0.0.1` et demande un port libre au système d’exploitation.

Utilisez `--stdio` ou `--http` pour sélectionner explicitement un transport. En mode HTTP, utilisez `--port 0` pour choisir un port de bouclage libre.

#### Options MCP

| Option | Description | Valeur par défaut |
| --- | --- | --- |
| `--root` | Racine de projet autorisée. Les chemins et liens symboliques situés en dehors de celle-ci sont rejetés. | Répertoire actuel |
| `--token` | Jeton de session ou jeton Bearer pour les outils de modification et de contrôle des processus. Utilise `WAILS_MCP_TOKEN` comme solution de repli. | Généré de manière sécurisée |
| `--stdio` | Force le transport stdio. | Automatique |
| `--http` | Force le transport Streamable HTTP. | Automatique |
| `--port` | Port HTTP ; `0` sélectionne un port de bouclage libre. | `0` |

En mode HTTP, Wails affiche le point de terminaison et le jeton Bearer sur la sortie d’erreur standard. En mode stdio, le jeton est inclus dans les instructions d’initialisation de MCP. Le serveur ne permet pas l’exécution de commandes shell arbitraires. Les modèles distants et les dépôts Git distants nécessitent une approbation explicite au moyen de l’entrée `allowExternal` de l’outil.

### `doctor`

Effectue une vérification du système et affiche un rapport d’état.

```bash
wails3 doctor
```

## Commandes de génération

Les commandes de génération permettent de créer différents éléments du projet, tels que les liaisons, les icônes et les fichiers de construction. Toutes les commandes de génération utilisent la commande de base : `wails3 generate <command>`.

### `generate bindings`

Génère les liaisons et les modèles correspondant à votre code Go.

```bash
wails3 generate bindings [flags] [patterns...]
```

#### Options

| Option | Description | Valeur par défaut |
| --- | --- | --- |
| `-f` | Options de construction Go supplémentaires |  |
| `-d` | Répertoire de sortie | `frontend/bindings` |
| `-models` | Nom du fichier des modèles | `models` |
| `-index` | Nom du fichier d’index | `index` |
| `-ts` | Génère du TypeScript | `false` |
| `-i` | Utilise des interfaces TypeScript | `false` |
| `-b` | Utilise le runtime intégré | `false` |
| `-names` | Utilise les noms à la place des identifiants | `false` |
| `-noindex` | Ignore les fichiers d’index | `false` |
| `-noevents` | Ignorer la génération des liaisons liées aux événements | `false` |
| `-dry` | Simulation | `false` |
| `-silent` | Mode silencieux | `false` |
| `-v` | Sortie de débogage | `false` |
| `-clean` | Nettoyer le répertoire de sortie avant la génération | `true` |

### `generate build-assets`

Génère les ressources de compilation de votre application.

```bash
wails3 generate build-assets [flags]
```

#### Options

| Option | Description | Valeur par défaut |
| --- | --- | --- |
| `-name` | Nom du projet |  |
| `-dir` | Répertoire de sortie | `build` |
| `-silent` | Supprimer la sortie | `false` |
| `-company` | Nom de l’entreprise |  |
| `-productname` | Nom du produit |  |
| `-description` | Description du produit |  |
| `-version` | Version du produit |  |
| `-identifier` | Identifiant du produit | `com.wails.[name]` |
| `-copyright` | Mention de copyright |  |
| `-comments` | Commentaires du fichier |  |

### `generate icons`

Génère les icônes de l’application.

```bash
wails3 generate icons [flags]
```

#### Options

| Option | Description | Valeur par défaut |
| --- | --- | --- |
| `-input` | Fichier PNG d’entrée | Obligatoire |
| `-windowsfilename` | Nom du fichier de sortie Windows |  |
| `-macfilename` | Nom du fichier de sortie macOS |  |
| `-sizes` | Tailles des icônes (séparées par des virgules) | `256,128,64,48,32,16` |
| `-example` | Générer un exemple d’icône | `false` |
| `-iconcomposerinput` | Fichier Icon Composer d’entrée (`.icon`) |  |
| `-macassetdir` | Répertoire de sortie des ressources Mac (Assets.car + icns) |  |

#### Icon Composer (macOS)

Sous macOS 26 ou version ultérieure, vous pouvez utiliser les fichiers `.icon` d’Icon Composer pour générer `Assets.car` et `icons.icns` :

```bash
wails3 generate icons -iconcomposerinput build/appicon.icon -macassetdir build
```

Cette opération compile le fichier `.icon` à l’aide de la commande `actool` d’Apple. Elle nécessite Xcode avec `actool` version 26 ou ultérieure.

Lorsque vous utilisez Icon Composer, définissez `cfBundleIconName` dans votre `build/config.yml` afin qu’il corresponde au nom du fichier `.icon` (sans l’extension) :

```yaml
info:
  cfBundleIconName: "appicon"
```

Si cette valeur n’est pas définie et que `Assets.car` existe, sa valeur par défaut est `"appicon"`.

### `generate syso`

Génère un fichier Windows .syso.

```bash
wails3 generate syso [flags]
```

#### Options

| Option | Description | Valeur par défaut |
| --- | --- | --- |
| `-manifest` | Chemin du fichier manifeste | Obligatoire |
| `-icon` | Chemin du fichier d’icône | Obligatoire |
| `-info` | Chemin du fichier d’informations de version |  |
| `-arch` | Architecture cible | GOARCH actuel |
| `-out` | Nom du fichier de sortie | `rsrc_windows_[arch].syso` |

### `generate .desktop`

Génère un fichier Linux .desktop.

```bash
wails3 generate .desktop [flags]
```

#### Options

| Option | Description | Valeur par défaut |
| --- | --- | --- |
| `-name` | Nom de l’application | Obligatoire |
| `-exec` | Chemin de l’exécutable | Obligatoire |
| `-icon` | Chemin de l’icône |  |
| `-categories` | Catégories de l’application | `Utility` |
| `-comment` | Commentaire de l’application |  |
| `-terminal` | Exécuter dans un terminal | `false` |
| `-keywords` | Mots-clés de recherche |  |
| `-version` | Version de l’application |  |
| `-genericname` | Nom générique |  |
| `-startupnotify` | Afficher une notification au démarrage | `false` |
| `-mimetype` | Types MIME pris en charge |  |
| `-output` | Nom du fichier de sortie | `[name].desktop` |

### `generate runtime`

Génère la version précompilée du runtime.

```bash
wails3 generate runtime
```

### `generate constants`

Génère des constantes JavaScript à partir du code Go.

```bash
wails3 generate constants
```

### `generate webview2bootstrapper`

Génère un programme d’installation initiale de WebView2 pour Windows destiné à la distribution.

```bash
wails3 generate webview2bootstrapper [flags]
```

### `generate template`

Crée la structure d’un nouveau répertoire de modèle de projet.

```bash
wails3 generate template [flags]
```

### `generate appimage`

Génère une AppImage Linux.

```bash
wails3 generate appimage [flags]
```

#### Options

| Option | Description | Valeur par défaut |
| --- | --- | --- |
| `-binary` | Chemin du fichier binaire | Obligatoire |
| `-icon` | Chemin du fichier d’icône | Obligatoire |
| `-desktop` | Chemin du fichier .desktop | Obligatoire |
| `-builddir` | Répertoire de compilation | Répertoire temporaire |
| `-output` | Répertoire de sortie | `.` |

## Commandes de service

Les commandes de service permettent de gérer les services Wails. Elles utilisent toutes la commande de base : `wails3 service <command>`.

### `service init`

Initialise un nouveau service.

```bash
wails3 service init [flags]
```

#### Options

| Option | Description | Valeur par défaut |
| --- | --- | --- |
| `-n` | Nom du service | `example_service` |
| `-d` | Description du service | `Example service` |
| `-p` | Nom du paquet |  |
| `-o` | Répertoire de sortie | `.` |
| `-q` | Supprimer la sortie | `false` |
| `-a` | Nom de l’auteur |  |
| `-v` | Version |  |
| `-w` | URL du site web |  |
| `-r` | URL du dépôt |  |
| `-l` | Licence |  |

## Commandes d’outils

Les commandes d’outils fournissent des utilitaires de développement et de débogage. Elles utilisent toutes la commande de base : `wails3 tool <command>`.

### `tool checkport`

Vérifie si un port est ouvert. Utile pour vérifier si vite est en cours d’exécution.

```bash
wails3 tool checkport [flags]
```

#### Options

| Option | Description | Valeur par défaut |
| --- | --- | --- |
| `-port` | Port à vérifier | `9245` |
| `-host` | Hôte à vérifier | `localhost` |

### `tool watcher`

Surveille les fichiers et exécute une commande lorsqu’ils sont modifiés.

```bash
wails3 tool watcher [flags]
```

#### Options

| Option | Description | Valeur par défaut |
| --- | --- | --- |
| `-config` | Chemin du fichier de configuration | `./build/config.yml` |
| `-ignore` | Motifs à ignorer |  |
| `-include` | Motifs à inclure |  |

### `tool cp`

Copie des fichiers.

```bash
wails3 tool cp
```

### `tool buildinfo`

Affiche les informations de compilation de l’application.

```bash
wails3 tool buildinfo
```

### `tool version`

Incrémente une version sémantique en fonction des options fournies.

```bash
wails3 tool version [flags]
```

#### Options

| Option | Description | Valeur par défaut |
| --- | --- | --- |
| `-v` | Version actuelle à incrémenter |  |
| `-major` | Incrémenter la version majeure | `false` |
| `-minor` | Incrémenter la version mineure | `false` |
| `-patch` | Incrémenter la version corrective | `false` |
| `-prerelease` | Incrémenter la version de préversion (par exemple, de alpha.5 à alpha.6) | `false` |

La commande respecte l’ordre de priorité suivant : majeure > mineure > corrective > préversion. Elle conserve le préfixe « v » s’il figure dans la version d’entrée, ainsi que les éventuels composants de préversion et de métadonnées.

Exemple d’utilisation :

```bash
wails3 tool version -v 1.2.3 -major      # Output: 2.0.0
wails3 tool version -v v1.2.3 -minor     # Output: v1.3.0
wails3 tool version -v 1.2.3-alpha -patch # Output: 1.2.4-alpha
wails3 tool version -v v3.0.0-alpha.5 -prerelease # Output: v3.0.0-alpha.6
```

### `tool package`

Génère des paquets Linux (deb, rpm, archlinux).

```bash
wails3 tool package [flags]
```

#### Options

| Option | Description | Valeur par défaut |
| --- | --- | --- |
| `-format` | Format du paquet (deb, rpm, archlinux) | `deb` |
| `-name` | Nom de l’exécutable | `myapp` |
| `-config` | Chemin du fichier de configuration |  |
| `-out` | Répertoire de sortie | `.` |

### `tool lipo`

Crée un binaire universel macOS en combinant des binaires propres à chaque architecture.

```bash
wails3 tool lipo [flags]
```

#### Options

| Option | Description | Valeur par défaut |
| --- | --- | --- |
| `-output` | Chemin du binaire de sortie |  |

### `tool capabilities`

Vérifie les capacités de compilation du système (disponibilité de GTK4/GTK3 sous Linux).

```bash
wails3 tool capabilities
```

### `tool docker-mounts`

Génère les options de montage de volumes Docker nécessaires à la compilation croisée. Produit des options `-v` pour le cache des modules Go et pour toutes les directives `replace` locales dans `go.mod`, afin de les utiliser dans les commandes `docker run` du Taskfile.

```bash
wails3 tool docker-mounts
```

### `tool has`

Vérifie si un outil ou une capacité est disponible et affiche `true` ou `false` sur la sortie standard. Cette commande est conçue pour être utilisée dans les variables `sh:` d’un Taskfile, comme solution multiplateforme remplaçant `command -v`.

Utilisez `|` pour vérifier si au moins l’une de plusieurs alternatives est disponible.

```bash
wails3 tool has <tool>
```

#### Exemples

```bash
# Check for a C compiler (gcc or clang)
wails3 tool has gcc|clang

# Check for a specific tool
wails3 tool has git
wails3 tool has node
```

#### Utilisation dans un Taskfile

```yaml
vars:
  HAS_CC:
    sh: 'wails3 tool has gcc|clang'
```

### `tool has-cc`

@note{type="caution" title="Obsolète"}
`wails3 tool has-cc` est obsolète. Mettez à jour votre Taskfile pour utiliser `wails3 tool has gcc|clang` à la place.

@end

Alias rétrocompatible de `wails3 tool has gcc|clang`. Vérifie si `gcc` ou `clang` est disponible dans PATH et affiche `true` ou `false`.

```bash
wails3 tool has-cc
```

## Commandes de mise à jour

Les commandes de mise à jour permettent de gérer et de mettre à jour les ressources du projet. Toutes utilisent la commande de base suivante : `wails3 update <command>`.

### `update cli`

Met à jour la CLI Wails vers une nouvelle version.

```bash
wails3 update cli [flags]
```

#### Options

| Option | Description | Valeur par défaut |
| --- | --- | --- |
| `-pre` | Mettre à jour vers la dernière préversion | `false` |
| `-version` | Mettre à jour vers une version précise |  |
| `-nocolour` | Désactiver la sortie en couleur | `false` |

La commande update cli vous permet de mettre à jour votre installation de la CLI Wails. Par défaut, elle installe la dernière version stable. Utilisez l’option `-pre` pour installer la dernière préversion, ou l’option `-version` pour indiquer une version précise.

Après la mise à jour, pensez à mettre à jour le fichier go.mod de votre projet afin d’utiliser la même version :

```bash
require github.com/wailsapp/wails/v3 v3.x.x
```

### `update build-assets`

Met à jour les ressources de compilation à l’aide du fichier de configuration indiqué.

```bash
wails3 update build-assets [flags]
```

#### Options

| Option | Description | Valeur par défaut |
| --- | --- | --- |
| `-config` | Chemin du fichier de configuration |  |
| `-dir` | Répertoire de sortie | `build` |
| `-silent` | Masquer la sortie | `false` |
| `-company` | Nom de l’entreprise |  |
| `-productname` | Nom du produit |  |
| `-description` | Description du produit |  |
| `-version` | Version du produit |  |
| `-identifier` | Identifiant du produit |  |
| `-copyright` | Mention de droits d’auteur |  |
| `-comments` | Commentaires du fichier |  |

## Commandes utilitaires

Les commandes utilitaires fournissent des raccourcis pratiques pour les tâches courantes. Utilisez-les directement avec la commande de base : `wails3 <command>`.

### `docs`

Ouvre la documentation de Wails dans votre navigateur par défaut.

```bash
wails3 docs
```

### `releasenotes`

Affiche les notes de publication de la version actuelle ou de la version spécifiée.

```bash
wails3 releasenotes [flags]
```

#### Options

| Option | Description | Valeur par défaut |
| --- | --- | --- |
| `-v` | Version dont les notes de publication doivent être affichées |  |
| `-n` | Désactiver la sortie en couleur | `false` |

### `version`

Affiche la version actuelle de Wails.

```bash
wails3 version
```

### `sponsor`

Ouvre la page de parrainage de Wails dans votre navigateur par défaut.

```bash
wails3 sponsor

```
