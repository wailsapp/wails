---
title: "Référence de la CLI"
description: "Référence complète des commandes de la CLI Wails"
slug: "reference/cli"
sourcePath: "reference/cli.md"
---

## Vue d’ensemble

La CLI Wails (`wails3`) est le point d’entrée en ligne de commande permettant de créer, développer, compiler, signer, empaqueter et inspecter des applications Wails 3. L’essentiel de l’orchestration de la compilation est délégué aux fichiers Taskfile propres à chaque projet (dans le répertoire `build/` de votre projet) : de nombreuses commandes `wails3` sont de simples enveloppes qui invoquent une tâche précise.

Pour obtenir l’aide la plus récente sur une commande, exécutez :

```bash
wails3 --help
wails3 <command> --help
```

## Cycle de vie du projet

| Commande | Description |
| --- | --- |
| `wails3 init` | Créez un projet à partir d’un modèle. Options : `-n` (nom du projet), `-t` (modèle, `vanilla` par défaut), `-p` (nom du paquet Go, `main` par défaut), `-d` (répertoire du projet, `.` par défaut), `-q` (mode silencieux), `-l` (afficher la liste des modèles), `-mod` (chemin du module Go), `--git` (URL du dépôt Git), `--skipgomodtidy`, `-s` (ignorer l’avertissement relatif aux modèles distants), `--productname`/`--productdescription`/`--productversion`/`--productcompany`/`--productcopyright`/`--productcomments`/`--productidentifier`. |
| `wails3 dev` | Exécutez l’application en mode développement avec rechargement à chaud du frontend. Options : `--config` (`./build/config.yml` par défaut), `--port` (port de développement de Vite), `-s` (activer HTTPS). |
| `wails3 build` | Compilez le projet. Cette commande est une simple enveloppe autour de la tâche `build` du Taskfile. Options : `--tags` (transmise sous la forme `EXTRA_TAGS=`), `--obfuscated` (compiler avec Garble ; voir [Compilations obscurcies](/guides/build/obfuscation/)), `--garbleargs` (options supplémentaires transmises à `garble` avant la sous-commande `build`). |
| `wails3 package` | Exécutez la tâche `package` du Taskfile propre à la plateforme. |
| `wails3 task [name]` | Exécutez n’importe quelle tâche du Taskfile ; si aucun nom n’est fourni, `--list` affiche toutes les tâches enregistrées. |
| `wails3 mcp` | Démarrez le serveur MCP du projet. Il utilise automatiquement l’entrée-sortie standard pour les processus lancés par un agent, ou le protocole HTTP Streamable sur l’interface de bouclage lors d’une utilisation interactive dans un terminal. |
| `wails3 doctor` | Affichez un rapport de diagnostic de votre environnement. |
| `wails3 doctor-ng` | Variante TUI plus récente de `doctor`. |
| `wails3 version` | Affichez la version de la CLI. |
| `wails3 releasenotes` | Affichez les notes de version récentes. |
| `wails3 docs` | Ouvrez le site de documentation dans votre navigateur. |
| `wails3 sponsor` | Ouvrez la page de parrainage. |

## Génération

`wails3 generate <subcommand>` :

| Sous-commande | Description |
| --- | --- |
| `generate bindings` | Générez les liaisons entre Go et le frontend. Options : `-d` (répertoire de sortie), `-models`, `-index`, `-ts`, `-i` (interfaces), `-b` (bundle), `-names` (émettre `Call.ByName`), `-noevents`, `-noindex`, `-dry`, `-silent`, `-v`, `-clean` (`true` par défaut), `-f`, `-obfuscated` (générer `wails_obfuscated.gen.go` avec des identifiants de liaison stables pour les compilations avec Garble ; voir [Compilations obscurcies](/guides/build/obfuscation/)), `-obfuscated-output` (répertoire du fichier généré ; le répertoire du paquet principal est utilisé par défaut). Accepte les motifs de paquets (par exemple `./...`) ; si aucun n’est fourni, utilise le répertoire courant. |
| `generate icons` | Convertissez un fichier PNG source aux formats d’icône propres aux différentes plateformes. Options : `-input`, `-windowsfilename`, `-macfilename`, `-iconcomposerinput`, `-macassetdir`. |
| `generate build-assets` | Générez le contenu du répertoire `build/` (extraits de Taskfile, fichiers NSIS, `Info.plist`, modèle `.desktop`, etc.) à partir de `build/config.yml`. |
| `generate runtime` | Régénérez le fichier `/wails/runtime.js` précompilé fourni à la webview. |
| `generate syso` | Générez le fichier de ressources Windows `.syso` (icône, manifeste et informations de version). |
| `generate webview2bootstrapper` | Générez un programme d’installation d’amorçage de WebView2 pour Windows. |
| `generate constants` | Générez les constantes JavaScript des noms d’événements à partir des types d’événements Go. |
| `generate template` | Créez la structure d’un nouveau modèle de projet. |
| `generate .desktop` | Générez un fichier Linux `.desktop` (utilisé par AppImage/DEB/RPM). |
| `generate appimage` | Générez le répertoire de compilation AppImage. |

## Mise à jour

`wails3 update <subcommand>` :

| Sous-commande | Description |
| --- | --- |
| `update build-assets` | Actualisez le répertoire `build/` à partir de `build/config.yml` (en conservant les modifications de l’utilisateur dans la mesure du possible). |
| `update cli` | Mettez automatiquement à jour le binaire `wails3`. |

## Signature du code et empaquetage

| Commande | Description |
| --- | --- |
| `wails3 setup signing` | Assistant interactif qui configure la signature pour les plateformes qu’il détecte dans `build/`. Option : `--platform` (répétable ; par défaut, les plateformes sont détectées automatiquement à partir du répertoire de compilation). |
| `wails3 setup entitlements` | Assistant interactif pour les autorisations macOS. Option : `--output` (chemin ; valeur par défaut : `build/darwin/entitlements.plist`). |
| `wails3 sign [GOOS=…]` | Commande enveloppe qui exécute la tâche Taskfile `*:sign` propre à la plateforme pour le système d’exploitation actuel (ou celui indiqué avec `GOOS`). |
| `wails3 tool sign` | Point d’entrée de bas niveau pour la signature directe. Options : `--input`, `--output`, `--verbose`, `--certificate`, `--password`, `--thumbprint`, `--timestamp`, `--identity`, `--entitlements`, `--hardened-runtime`, `--notarize`, `--keychain-profile`, `--pgp-key`, `--pgp-password`, `--role`. |

Il n’existe **aucune** sous-commande `wails3 signing` : pour les identifiants du trousseau, utilisez directement `xcrun notarytool store-credentials`, et pour les clés PGP, utilisez directement `gpg` (l’assistant `wails3 setup signing` automatise les deux opérations).

## Outils

`wails3 tool <subcommand>` :

| Sous-commande | Description |
| --- | --- |
| `tool checkport` | Vérifier si un port TCP est ouvert (utile pour attendre Vite). |
| `tool watcher` | Exécuter une commande chaque fois que les fichiers surveillés changent. |
| `tool cp` | Copier des fichiers de manière multiplateforme. |
| `tool buildinfo` | Afficher les informations de compilation Go intégrées à un fichier binaire. |
| `tool package` | Créer un paquet Linux (`deb`, `rpm`, `archlinux`) à partir de `build/linux/nfpm`. |
| `tool version` | Incrémenter la version sémantique d’un projet. |
| `tool lipo` | Combiner des fichiers binaires de plusieurs architectures macOS en un binaire universel. |
| `tool capabilities` | Analyser le système pour déterminer la disponibilité de GTK3/GTK4 et de WebKit. |
| `tool sign` | (Voir [Signature du code et création de paquets](#signature-du-code-et-empaquetage).) |

## Services

`wails3 service <subcommand>` :

| Sous-commande | Description |
| --- | --- |
| `service init` | Générer la structure d’un nouveau paquet de service. |

## iOS

`wails3 ios <subcommand>` :

| Sous-commande | Description |
| --- | --- |
| `ios overlay:gen` | Générer l’overlay Go pour la couche d’adaptation du pont iOS. |
| `ios xcode:gen` | Générer un projet Xcode dans le répertoire de sortie. |

## Chemins des sorties de compilation

- Les fichiers binaires natifs sont placés dans `bin/<APP_NAME>` (ou `bin/<APP_NAME>.exe` sous Windows). Il n’existe aucun `build/bin/`.
- Les sorties empaquetées (`.app`, `.dmg`, programme d’installation NSIS, MSIX, DEB/RPM/AppImage) sont également placées dans `bin/` (ou dans les sous-répertoires propres à la plateforme créés par la tâche Taskfile correspondante).

## Options globales

| Option | S’applique à | Description |
| --- | --- | --- |
| `--no-colour` | Toutes les commandes | Désactiver les couleurs ANSI dans la sortie de l’interface en ligne de commande. |

---

**Des questions ?** Posez-les sur [Discord](https://discord.gg/JDdSxwjhGf) ou consultez les [exemples](https://github.com/wailsapp/wails/tree/master/v3/examples).
