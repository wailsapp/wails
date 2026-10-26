---
title: "Sortie de Wails v3.0.0"
description: "Wails v3.0.0 est disponible avec la gestion de plusieurs fenêtres, une icône dans la zone de notification, des menus natifs, des services et des liaisons générées, un outil de mise à jour, la prise en charge des appareils mobiles et un nouveau site de documentation."
authors: ["leaanthony"]
tags: ["wails","v3","release"]
date: "2026-10-26"
slug: "blog/wails-v3-released"
image: "/assets/videos/wails-v3-release-clean-poster.webp"
sourcePath: "blog/wails-v3-released.md"
---

@video{src="/assets/videos/wails-v3-release-clean-av1.mp4" poster="/assets/videos/wails-v3-release-clean-poster.webp" title="Vidéo de lancement de Wails v3.0.0 — à regarder de préférence en plein écran" width="1920" height="1080"}

## Le voici !

Aujourd’hui, je publie Wails v3.0.0.

Je vais être franc : le chemin a été long. La première version alpha de v3 a été publiée en janvier 2023 et, depuis, un très grand nombre de personnes ont développé, testé, cassé et corrigé des choses tout au long du parcours. Je vous en suis sincèrement reconnaissant à toutes et à tous. Si vous souhaitez connaître la version sans détour de cette aventure, je l’ai racontée dans l’[v3annonce de la version bêta](/blog/wails-v3-beta/).

## *Qu’est-ce que* Wails ?

Si vous découvrez Wails : il permet aux développeurs Go de créer des applications de bureau avec les outils frontend web qu’ils connaissent déjà, en utilisant la WebView native de chaque plateforme plutôt qu’un navigateur intégré. Vous écrivez le backend en Go et l’interface utilisateur en HTML, CSS et JavaScript (ou avec le framework de votre choix), puis Wails regroupe le tout dans une application native unique.

## Pourquoi v3 ?

Wails v2 permettait de créer facilement une application Go dotée d’un frontend web moderne et a rendu de fiers services à de nombreuses applications. Cependant, son environnement d’exécution à fenêtre unique, fondé sur le contexte, compliquait inutilement certaines tâches courantes des applications de bureau.

v3 repose sur un modèle différent. Les applications, fenêtres, services, événements et fonctionnalités propres aux plateformes sont des objets explicites. Vous créez une application, enregistrez des services, créez des fenêtres et communiquez avec l’objet responsable du comportement dont vous avez besoin. Ce modèle reste plus facile à suivre à mesure que l’application se développe, et les fonctionnalités telles que la gestion de plusieurs fenêtres en font naturellement partie au lieu de nécessiter une solution de contournement.

## Contenu de v3.0.0

- **Plusieurs fenêtres.** Créez, gérez et fermez des fenêtres à l’exécution, chacune disposant de sa propre configuration et de son propre cycle de vie.
- **Zone de notification.** Assurez une présence permanente à votre application grâce à des menus, des icônes pour les thèmes clair et sombre, des icônes modèles sous macOS et la possibilité d’associer une fenêtre à l’icône de la zone de notification.
- **Menus et boîtes de dialogue natifs.** Des menus d’application, des menus contextuels et des boîtes de dialogue natives sur chaque plateforme.
- **Services et liaisons générées.** Votre backend reste du code Go ordinaire. Les liaisons sont générées par analyse statique du code source : les commentaires et les noms de paramètres explicites sont ainsi transmis au frontend à l’aide d’une seule commande `wails3 generate bindings`.
- **Un système de compilation transparent.** Au lieu d’une commande opaque, les projets utilisent un système de compilation fondé sur [Taskfile](https://taskfile.dev), que vous pouvez examiner, étendre et déboguer.
- **Multiplateforme.** Windows (amd64 et arm64), macOS (Intel et Apple Silicon) et Linux (amd64 et arm64), avec GTK4 et WebKitGTK 6.0 par défaut sous Linux. L’ancienne pile GTK3 reste disponible avec `-tags gtk3` jusqu’à v3.1.
- **iOS et Android.** Le même `main.go` et le même frontend peuvent s’exécuter sur les appareils mobiles. La prise en charge des appareils mobiles est encore signalée comme expérimentale dans la documentation : considérez-la donc comme une fonctionnalité à explorer et sur laquelle donner votre avis, et non comme une garantie de stabilité.
- **Compilé avec Go 1.26.** Le module v3 est compilé avec Go 1.26.
- **Un outil de mise à jour intégré à l’application.** Des fournisseurs interchangeables, une vérification cryptographique, un remplacement atomique et une interface utilisateur par défaut que vous pouvez personnaliser ou remplacer. Consultez le [guide de l’outil de mise à jour](/guides/updater/).
- **Création de paquets et signature.** Des guides sur la [signature du code](/guides/build/signing/) pour chaque plateforme, les [programmes d’installation](/guides/installers/) et la [création de paquets MSIX](/guides/build/msix/).
- **Compilations pour serveur.** Exécutez la même application et les mêmes services sous forme de serveur HTTP, sans fenêtre native.

## Principales nouveautés depuis la version bêta

Le billet consacré à la version bêta a été publié le 2 août et le journal des modifications s’est depuis bien étoffé. Voici quelques éléments que je tiens à souligner :

- **Un nouveau site de documentation.** La documentation repose désormais sur [M-Press](https://m-press.me), notre propre outil de documentation, et les traductions ont été achevées et reprises de celles de la communauté. La documentation anglaise reste la référence.
- **Une expérience plus robuste sous Windows.** Wails récupère désormais après une défaillance d’un processus WebView2 au lieu de vous laisser face à une fenêtre vide. Plusieurs problèmes liés au cycle de vie des objets COM de WebView2 ont également été corrigés.
- **Des raccourcis globaux adaptés à votre clavier.** Sous macOS, les raccourcis utilisant des lettres suivent désormais la touche étiquetée sur les dispositions telles qu’AZERTY, QWERTZ et Dvorak.
- **Renforcement de l’outil de mise à jour.** Entre autres améliorations, des correctifs concernent le remplacement de fichiers binaires entre différents systèmes de fichiers sous Linux et les mises à jour AppImage.
- **Configuration et outils.** L’assistant `wails3 setup` a été finalisé et `wails3 dev` arrête désormais ses processus en arrière-plan lorsque votre application se ferme.
- **Liquid Glass natif sous macOS.** Prise en charge de l’effet Liquid Glass sous macOS 26 et versions ultérieures, avec une solution de repli pour les versions antérieures.

## Migration depuis v2

v3 est une nouvelle version majeure : sa migration constitue un véritable portage, et non une simple mise à jour de version. Les principaux changements concernent le cycle de vie de l’application et des fenêtres, le remplacement des liaisons associées au contexte par des services, le remplacement du paquet d’exécution v2 par des API directes pour les applications et les fenêtres, ainsi que la régénération des liaisons frontend.

Le [v2guide de migration de v3](/migration/v2-to-v3/) présente l’ensemble du processus à l’aide d’exemples avant et après. Effectuez le portage avec méthode, testez au fur et à mesure et conservez v2 jusqu’à ce que votre application v3 soit prête.

**Wails v2 ne va pas disparaître.** Il continuera de recevoir les mises à jour critiques. Pour les nouveaux projets, utilisez v3.0.0. Les deux interfaces en ligne de commande sont des fichiers binaires distincts (`wails` et `wails3`) et leurs chemins d’importation de modules diffèrent : vous pouvez donc installer les deux sur une même machine.

## Installation

Installez l’interface en ligne de commande v3.0.0 avec :

```sh
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0
```

Exécutez ensuite l’assistant de configuration guidée, qui vérifie votre environnement et vous aide à installer les dépendances nécessaires à Wails :

```sh
wails3 setup
```

![L’assistant de configuration de Wails v3 en mode sombre](/assets/screenshots/wails3-setup-wizard-dark.png)

Créez ensuite votre premier projet :

```sh
wails3 init
```

- Documentation : <https://v3.wails.io/>
- Guide de migration : <https://v3.wails.io/migration/v2-to-v3/>
- Notes de version : [Wails v3.0.0 sur GitHub](https://github.com/wailsapp/wails/releases/tag/v3.0.0)

## Assistance et versions

Wails a pris une ampleur qui dépasse largement mes attentes et il est devenu difficile de suivre le volume de tickets et de demandes d’intégration. Les publications ont parfois été effectuées au cas par cas, et les correctifs comme les contributions ont mis trop de temps à parvenir aux utilisateurs. Nous avons mis en place des systèmes afin d’automatiser une plus grande partie du travail préparatoire au triage et à la préparation des versions. Nous pouvons ainsi réagir plus rapidement aux problèmes et consacrer davantage de temps à livrer des améliorations et de nouvelles fonctionnalités.

Après v3.0.0, nous passerons à des **préversions nocturnes** portant des étiquettes datées telles que `v3.1.0-nightly.20261026`, parallèlement à des **versions mineures hebdomadaires** telles que `v3.1.0` et `v3.2.0`. Les versions nocturnes rendent les correctifs récents disponibles plus rapidement, tandis que les versions hebdomadaires vous offrent un rythme régulier de mises à jour testées.

Si vous trouvez un bug reproductible, veuillez le signaler en joignant la sortie de `wails3 doctor` et un exemple minimal. Si vous souhaitez proposer une nouvelle fonctionnalité ou une modification du comportement public, ouvrez une demande d’intégration à l’état de brouillon pour une WEP (Wails Enhancement Proposal), plutôt qu’un ticket de demande de fonctionnalité.

## Traductions

À toutes les personnes qui ont contribué à traduire la documentation de Wails v2 : merci du fond du cœur. Vous avez donné de votre temps pour que les développeurs puissent lire la documentation dans leur propre langue, et cela a beaucoup compté pour de nombreuses personnes. Nous vous en sommes profondément reconnaissants.

Nous avons trouvé Crowdin particulièrement difficile à gérer et avons cessé de l’utiliser. À partir de v3, la documentation est traduite automatiquement et les traductions se trouvent dans le dépôt, aux côtés des pages en anglais.

Si vous souhaitez une traduction dans une langue que nous ne proposons pas encore, ou si vous avez repéré une erreur dans une traduction existante, veuillez ouvrir un ticket pour nous en informer. Nous serions ravis de vous lire et vous remercions pour chaque signalement.

## Crédits et remerciements

Wails v3 existe grâce aux personnes qui ont testé des versions inachevées, signalé des bogues difficiles à résoudre, traduit la documentation, répondu aux questions et contribué au code. Merci. La liste complète des contributeurs figure sur la [page des crédits](/credits/), et je vous encourage à la consulter.

Tout d’abord, un immense merci à [Atterpac](https://github.com/atterpac), mainteneur de Wails, pour les outils de développement, l’intégration aux environnements de bureau Linux et le soutien à la communauté qui ont contribué à amener v3 jusqu’ici.

Un merci tout particulier à [Fabio Massaioli (fbbdev)](https://github.com/fbbdev) pour avoir créé un analyseur statique du code source et un générateur de bindings aussi remarquables. Ils jouent un rôle essentiel dans le fonctionnement de v3 : des services Go ordinaires deviennent des API JavaScript et TypeScript, et leurs types ainsi que leur documentation sont transmis à votre frontend.

Merci également à [Travis McLane](https://github.com/tmclane) pour son important travail initial sur les plateformes ; à [Stffabi](https://github.com/stffabi) pour le serveur de ressources et les fondations des plateformes natives ; à [Zach Botterman (popaprozac)](https://github.com/popaprozac) pour les notifications, l’intégration au Dock et d’autres fonctionnalités de bureau ; à [Ian VanSchooten (IanVS)](https://github.com/IanVS) pour les événements typés et les améliorations du runtime ; et à [Lyimmi](https://github.com/Lyimmi) pour son travail sur Linux. Ces contributions font toujours partie de v3, aux côtés du travail de toutes les personnes qui ont contribué à le développer et à le perfectionner.

Un immense merci, du fond du cœur, aux sponsors qui ont soutenu Wails tout au long de cette longue transition. Votre soutien nous a donné le temps de travailler sur l’architecture, les outils, les tests et la documentation qui nous ont permis d’arriver jusqu’ici. Merci tout particulièrement à JetBrains et ZSA, ainsi qu’à [CodeRabbit](https://www.coderabbit.ai/oss) pour son soutien au projet au moyen de revues de code.

Enfin, un grand merci à Valerio Tesei pour tout le soutien qu’il a apporté au projet au fil des années.

## Le mot de la fin

Beaucoup de travail nous attend déjà. Des avancées passionnantes ont lieu dans l’ensemble du monde technologique, et nous sommes prêts à explorer de nouveaux paradigmes et à doter Wails de fonctionnalités entièrement nouvelles. Mais nous devons d’abord faire de v3 une base stable sur laquelle vous pouvez compter. Il faut y parvenir avant de déployer la suite.

Je suis vraiment fière de v3. J’espère que vous prendrez autant de plaisir à créer avec que nous en avons pris à le développer. Essayez-le, créez quelque chose et dites-nous ce que vous en pensez. Venez nous saluer sur [Discord](https://discord.gg/JDdSxwjhGf), dans les [GitHub Discussions](https://github.com/wailsapp/wails/discussions) ou sur le [subreddit](https://www.reddit.com/r/wails/).

Amusez-vous bien !

&dash;&nbsp;Lea

P.-S. : si vous ou votre entreprise trouvez Wails utile, pensez à [sponsoriser le projet](https://github.com/sponsors/leaanthony). Merci !
