---
title: "Vue d’ensemble des commandes de la CLI"
description: "Trouvez les commandes de la CLI Wails selon leur fonction : projets, builds, signature, outils et services"
slug: "reference/cli"
sourcePath: "reference/cli.md"
---

Commencez ici pour trouver une commande selon sa fonction. Pour connaître les options propres à chaque commande et consulter des exemples d’utilisation, voir [Options de la CLI et exemples](/guides/cli/).

## Vue d’ensemble

La CLI Wails (`wails3`) constitue le point d’entrée en ligne de commande pour créer, développer, compiler, signer, empaqueter et inspecter les applications Wails 3. La majeure partie de l’orchestration des builds est déléguée aux Taskfiles propres à chaque projet (situés sous `build/` dans votre projet) : de nombreuses commandes `wails3` sont de simples enveloppes qui invoquent une tâche précise.

Pour obtenir l’aide la plus récente sur une commande, exécutez :

```bash
wails3 --help
wails3 <command> --help
```

## Cycle de vie du projet

| Commande | Description |
| --- | --- |
| `wails3 init` | Créez un projet à partir d’un modèle. Options : `-n` (nom du projet), `-t` (modèle, valeur par défaut : `vanilla`), `-p` (nom du paquet Go, valeur par défaut : `main`), `-d` (répertoire du projet, valeur par défaut : `.`), `-q` (mode silencieux), `-l` (afficher les modèles), `-mod` (chemin du module Go), `--git` (URL du dépôt Git), `--skipgomodtidy`, `-s` (ignorer l’avertissement relatif aux modèles distants), `--productname`/`--productdescription`/`--productversion`/`--productcompany`/`--productcopyright`/`--productcomments`/`--productidentifier`. |
| `wails3 dev` | Exécutez l’application en mode développement avec rechargement à chaud du frontend. Options : `--config` (valeur par défaut : `./build/config.yml`), `--port` (port du serveur de développement Vite), `-s` (activer HTTPS). |
| `wails3 build` | Compilez le projet. Il s’agit d’une simple enveloppe autour de la tâche `build` du Taskfile. Options : `--tags` (transmise sous la forme `EXTRA_TAGS=`), `--obfuscated` (compiler avec Garble ; voir [Builds obfusqués](/guides/build/obfuscation/)), `--garbleargs` (options supplémentaires transmises à `garble` avant la sous-commande `build`). |
| `wails3 package` | Exécutez la tâche `package` du Taskfile propre à la plateforme. |
| `wails3 task [name]` | Exécutez n’importe quelle tâche du Taskfile ; sans nom, `--list` affiche toutes les tâches enregistrées. |
| `wails3 mcp` | Démarrez le serveur MCP du projet. Il utilise automatiquement l’entrée-sortie standard pour les processus lancés par un agent, ou HTTP Streamable sur l’interface de bouclage lors d’une utilisation interactive dans un terminal. |
| `wails3 doctor` | Affichez un rapport de diagnostic de votre environnement. |
| `wails3 doctor-ng` | Variante TUI plus récente de `doctor`. |
| `wails3 version` | Affichez la version de la CLI. |
| `wails3 releasenotes` | Affichez les notes de version récentes. |
| `wails3 docs` | Ouvrez le site de documentation dans votre navigateur. |
| `wails3 sponsor` | Ouvrez la page de parrainage. |

## Génération

`wails3 generate <subcommand>`:

| Sous-commande | Description |
| --- | --- |
| `generate bindings` | Générez les liaisons de Go vers le frontend. Options : `-d` (répertoire de sortie), `-models`, `-index`, `-ts`, `-i` (interfaces), `-b` (bundle), `-names` (émettre `Call.ByName`), `-noevents`, `-noindex`, `-dry`, `-silent`, `-v`, `-clean` (valeur par défaut : `true`), `-f`, `-obfuscated` (générer `wails_obfuscated.gen.go` avec des identifiants de liaison stables pour les builds Garble ; voir [Builds obfusqués](/guides/build/obfuscation/)), `-obfuscated-output` (répertoire du fichier généré ; utilise par défaut le répertoire du paquet principal). Accepte les motifs de paquets (par exemple `./...`) ; utilise le répertoire courant si aucun n’est fourni. |
| `generate icons` | Convertissez une image PNG source aux formats d’icône des différentes plateformes. Options : `-input`, `-windowsfilename`, `-macfilename`, `-iconcomposerinput`, `-macassetdir`. |
| `generate build-assets` | Générez le contenu du répertoire `build/` (extraits de Taskfile, fichiers NSIS, `Info.plist`, modèle `.desktop`, etc.) à partir de `build/config.yml`. |
| `generate runtime` | Régénérez le fichier `/wails/runtime.js` précompilé fourni à la webview. |
| `generate syso` | Générez le fichier de ressources Windows `.syso` (icône, manifeste et informations de version). |
| `generate webview2bootstrapper` | Générez un programme d’installation d’amorçage WebView2 pour Windows. |
| `generate constants` | Générez les constantes JavaScript des noms d’événements à partir des types d’événements Go. |
| `generate template` | Créez la structure d’un nouveau modèle de projet. |
| `generate .desktop` | Générez un fichier Linux `.desktop` (utilisé par AppImage, DEB et RPM). |
| `generate appimage` | Générez le répertoire de build AppImage. |

## Mise à jour

`wails3 update <subcommand>`:

| Sous-commande | Description |
| --- | --- |
| `update build-assets` | Actualisez le répertoire `build/` à partir de `build/config.yml` (les modifications de l’utilisateur sont conservées dans la mesure du possible). |
| `update cli` | Mettez à jour le binaire `wails3` lui-même. |

## Signature du code et empaquetage

| Commande | Description |
| --- | --- |
| `wails3 signing` | Assistant interactif qui configure la signature pour les plateformes qu’il détecte dans `build/`. Options : `--platform` (répétable ; par défaut, détection automatique à partir du répertoire de build). |
| `wails3 entitlements` | Assistant interactif pour les autorisations macOS. Option : `--output` (chemin ; valeur par défaut : `build/darwin/entitlements.plist`). |
| `wails3 sign [GOOS=…]` | Enveloppe qui exécute la tâche `*:sign` du Taskfile propre à la plateforme pour le système d’exploitation actuel (ou celui indiqué via `GOOS`). |
| `wails3 tool sign` | Point d’entrée de bas niveau pour la signature directe. Options : `--input`, `--output`, `--verbose`, `--certificate`, `--password`, `--thumbprint`, `--timestamp`, `--identity`, `--entitlements`, `--hardened-runtime`, `--notarize`, `--keychain-profile`, `--pgp-key`, `--pgp-password`, `--role`. |

Dans cette version, l’analyseur de la CLI expose `signing` et `entitlements` au niveau supérieur. Consultez `wails3 signing --help` et `wails3 entitlements --help` ; `wails3 setup signing` ne sélectionne pas l’assistant de signature. Pour gérer directement les identifiants, utilisez `xcrun notarytool store-credentials` ou `gpg`.

## Outils

`wails3 tool <subcommand>`:

| Sous-commande | Description |
| --- | --- |
| `tool checkport` | Vérifie si un port TCP est ouvert (utile pour attendre Vite). |
| `tool watcher` | Exécute une commande chaque fois que les fichiers surveillés changent. |
| `tool cp` | Copie des fichiers de manière multiplateforme. |
| `tool buildinfo` | Affiche les informations de compilation Go intégrées à un fichier binaire. |
| `tool package` | Crée un paquet Linux (`deb`, `rpm`, `archlinux`) à partir de `build/linux/nfpm`. |
| `tool version` | Incrémente la version sémantique d’un projet. |
| `tool lipo` | Combine des fichiers binaires pour plusieurs architectures macOS en un binaire universel. |
| `tool capabilities` | Détecte la disponibilité de GTK3/GTK4 et de WebKit sur le système. |
| `tool sign` | (Voir [Signature du code et création de paquets](#signature-du-code-et-empaquetage).) |

## Services

`wails3 service <subcommand>`:

| Sous-commande | Description |
| --- | --- |
| `service init` | Génère la structure initiale d’un nouveau paquet de service. |

## iOS

`wails3 ios <subcommand>`:

| Sous-commande | Description |
| --- | --- |
| `ios overlay:gen` | Génère la superposition Go pour la couche d’adaptation du pont iOS. |
| `ios xcode:gen` | Génère un projet Xcode dans le répertoire de sortie. |

## Chemins des sorties de compilation

- Les fichiers binaires natifs sont placés dans `bin/<APP_NAME>` (ou dans `bin/<APP_NAME>.exe` sous Windows). Il n’existe aucun `build/bin/`.
- Les sorties empaquetées (`.app`, `.dmg`, programme d’installation NSIS, MSIX, DEB/RPM/AppImage) sont également placées dans `bin/` (ou dans les sous-répertoires propres à chaque plateforme créés par la tâche Taskfile correspondante).

## Options communes

| Option | S’applique à | Description |
| --- | --- | --- |
| `-nocolour` | Commandes qui l’exposent dans `--help` (par exemple `init` et `update cli`) | Désactive les couleurs ANSI dans la sortie de la CLI. Cette option n’est pas acceptée par toutes les commandes, notamment `version`. |

---

**Des questions ?** Posez-les sur [Discord](https://discord.gg/JDdSxwjhGf) ou consultez les [exemples](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples).
