---
title: "Installation"
description: "Installez Wails et configurez votre environnement de développement"
slug: "getting-started/installation"
sourcePath: "getting-started/installation.md"
---

## Plateformes prises en charge

- Windows AMD64/ARM64
- macOS 10.15+ AMD64 (déploiement possible sur macOS 10.13+)
- macOS 11.0+ ARM64
- Ubuntu 24.04 AMD64/ARM64 (d’autres distributions Linux peuvent également fonctionner !)

## Dépendances

Wails nécessite plusieurs dépendances communes avant son installation.

@note{type="tip"}
Après avoir installé l’interface en ligne de commande de Wails, vous pouvez exécuter `wails3 setup` pour vérifier automatiquement ces dépendances et obtenir de l’aide pour les installer.

@end

@tabs
[Go (1.24 au minimum)]
Téléchargez Go depuis la [page de téléchargement de Go](https://go.dev/dl/).

Veillez à suivre les [instructions officielles d’installation de Go](https://go.dev/doc/install). Vérifiez également que votre variable d’environnement `PATH` contient le chemin vers votre répertoire `~/go/bin`. Redémarrez votre terminal, puis effectuez les vérifications suivantes :

- Vérifiez que Go est correctement installé : `go version`
- Vérifiez que `~/go/bin` figure dans votre variable PATH
  - Mac/Linux : `echo $PATH | grep go/bin`
  - Windows : `$env:PATH -split ';' | Where-Object { $_ -like '*\go\bin' }`


[npm (facultatif)]
Bien que Wails ne nécessite pas l’installation de npm, la plupart des modèles fournis en ont besoin.

Téléchargez la dernière version du programme d’installation de Node depuis la [page de téléchargement de Node](https://nodejs.org/en/download/). Il est préférable d’utiliser la version la plus récente, car c’est généralement celle que nous testons.

Exécutez `npm --version` pour vérifier l’installation.

@note{type="info"}
Si vous préférez un autre gestionnaire de paquets à npm, vous pouvez l’utiliser. Vous devrez adapter les fichiers Taskfile du projet en conséquence.

@end

@end

## Dépendances propres à chaque plateforme

Vous devrez également installer les dépendances propres à votre plateforme :

@tabs{sync-key="platform"}
[Mac]
Wails nécessite l’installation des outils en ligne de commande de Xcode. Pour les installer, exécutez :

```sh
xcode-select --install
```

[Windows]
Wails nécessite l’installation de [WebView2 Runtime](https://developer.microsoft.com/en-us/microsoft-edge/webview2/). Celui-ci est déjà installé sur la quasi-totalité des systèmes Windows. Vous pouvez le vérifier à l’aide de la commande `wails doctor`.

[Linux]
Linux nécessite les outils de compilation `gcc` standard, ainsi que `gtk4` et `webkitgtk-6.0`. Après l’installation, exécutez <code>wails3 doctor</code> pour savoir comment installer les dépendances. L’ancienne pile GTK3/WebKit2GTK 4.1 reste disponible via `-tags gtk3` (voir [Création de paquets Linux – Prise en charge de l’ancien GTK3](/guides/build/linux/#legacy-gtk3-support)) jusqu’à la version v3.1. Si votre distribution ou votre gestionnaire de paquets n’est pas pris en charge, veuillez nous le signaler sur Discord.

@end

## Installation

Pour installer l’interface en ligne de commande de Wails à l’aide des modules Go, exécutez les commandes suivantes :

```shell
go install -v github.com/wailsapp/wails/v3/cmd/wails3@latest
```

Pour installer la dernière version de développement, exécutez les commandes suivantes :

```shell
git clone https://github.com/wailsapp/wails.git
cd wails
cd v3/cmd/wails3
go install
```

Lorsque vous utilisez la version de développement, tous les projets générés utilisent la directive [replace](https://go.dev/ref/mod#go-mod-file-replace) de Go afin de garantir qu’ils utilisent la version de développement de Wails.

## Étapes suivantes

Après avoir installé l’interface en ligne de commande, lancez l’assistant de configuration pour configurer votre environnement de développement :

```shell
wails3 setup
```

@note{type="caution" title="Expérimental"}
L’assistant de configuration est récent et a principalement été testé sous Linux. Si vous rencontrez des problèmes, veuillez [les signaler](https://github.com/wailsapp/wails/issues/4904) et suivre les étapes d’installation manuelle des dépendances ci-dessous.

@end

L’assistant de configuration effectuera les opérations suivantes :

- Vérifier les dépendances de la plateforme et vous aider à les installer
- Configurer les valeurs par défaut du projet (informations sur l’auteur, préfixe de l’identifiant du bundle)
- Configurer éventuellement Docker pour les compilations multiplateformes
- Configurer la signature du code (si nécessaire)

Pour en savoir plus, consultez le [guide de configuration](/getting-started/setup/).

## Installation manuelle des dépendances

Si vous préférez installer les dépendances manuellement, ou si l’assistant de configuration ne fonctionne pas sur votre système, suivez les instructions propres à votre plateforme ci-dessus, puis exécutez :

```shell
wails3 doctor
```

Cette commande vérifiera que les dépendances appropriées sont installées et vous indiquera celles qui manquent.

## La commande `wails3` semble introuvable ?

Si votre système indique que la commande `wails3` est introuvable, vérifiez les points suivants :

- Assurez-vous d’avoir correctement suivi le **guide d’installation de Go** ci-dessus et que le répertoire `go/bin` figure dans la variable d’environnement `PATH`.
- Fermez puis rouvrez les terminaux actuels afin qu’ils prennent en compte la nouvelle variable `PATH`.
