---
title: "Wails v3 Beta : une nouvelle base pour les applications de bureau en Go"
description: "Wails v3 Beta introduit un modèle d’application plus direct, des bindings plus riches et une base plus claire pour les applications de bureau en Go."
authors: ["leaanthony"]
tags: ["wails","v3","beta"]
date: "2026-08-02"
slug: "blog/wails-v3-beta"
image: "/assets/screenshots/frameless-v3-native-corners-macos.png"
sourcePath: "blog/wails-v3-beta.md"
---

![Une fenêtre Wails v3 native sans cadre sous macOS](/assets/screenshots/frameless-v3-native-corners-macos.png)

Nous publions aujourd’hui Wails v3 Beta.

Wails permet aux développeurs Go de créer des applications de bureau avec les outils de développement web frontend qu’ils connaissent déjà, en utilisant la WebView native de chaque plateforme plutôt qu’un navigateur embarqué. v3 constitue une avancée majeure : cette version fournit aux applications une API plus directe, un modèle de build plus clair et une meilleure base pour les applications de bureau dont les utilisateurs demandent la prise en charge par Wails.

Il s’agit d’une version bêta, et non de la version 3.0 finale. L’API de bureau est stable et des équipes utilisent déjà v3 en production, mais vous devriez effectuer des tests approfondis avant tout déploiement. Avec la communauté, nous consacrons cette période bêta à la détection des derniers problèmes de compatibilité et de workflow. Wails v2 reste la version stable actuelle et continuera de recevoir des correctifs.

## Documentation pendant la phase bêta

Pendant la phase bêta, nous faisons de la documentation anglaise la source de référence tandis que l’API et les workflows font l’objet de leur validation finale. À ce stade, nous n’acceptons pas les pull requests de traduction. Le travail de traduction reprendra avant la disponibilité générale, une fois la documentation suffisamment stable pour que les traducteurs puissent travailler sans avoir à reprendre constamment leurs traductions.

## Contenu de v3

- Une API explicite pour les applications et les fenêtres, avec notamment une prise en charge native de plusieurs fenêtres
- Des services Go analysés statiquement au niveau du code source afin de générer des bindings TypeScript plus riches, tout en préservant les commentaires et les noms de paramètres explicites
- Des services capables de regrouper des ressources frontend et des scripts avec leur API backend, formant ainsi la base de plugins installables plus riches
- Un système de build visible, fondé sur Taskfile, que vous pouvez examiner, étendre et déboguer
- Des builds serveur permettant d’exécuter la même application et les mêmes services sans fenêtre de bureau native
- Une prise en charge moderne des environnements de bureau macOS, Windows et Linux sur Intel, Apple Silicon, amd64 et arm64, selon les plateformes compatibles
- Une prise en charge expérimentale des appareils mobiles sous iOS et Android, disponible à des fins d’exploration, mais exclue de la garantie de compatibilité de la version bêta pour ordinateurs de bureau

## Pourquoi v3 ?

Wails v2 a simplifié la création d’applications Go dotées d’un frontend web moderne. Cette version a rendu de grands services au projet ainsi qu’à un très grand nombre d’applications. Toutefois, son environnement d’exécution à fenêtre unique piloté par le contexte et son processus de build étroitement contrôlé ont rendu certaines tâches courantes de développement d’applications de bureau plus difficiles qu’elles n’auraient dû l’être.

v3 repose sur un modèle différent. Les applications, fenêtres, services, événements et fonctionnalités propres aux plateformes sont des objets explicites. Il devient ainsi plus facile de comprendre le framework à mesure qu’une application se développe, et des fonctionnalités telles que la gestion de plusieurs fenêtres font désormais partie intégrante du modèle d’application au lieu de constituer une solution de contournement.

## Nouveautés

### Une API d’application conçue pour de véritables logiciels de bureau

v3 remplace le style de configuration `wails.Run(...)` de v2 par un cycle de vie explicite de l’application. Vous créez une application, enregistrez des services, créez des fenêtres et interagissez avec les objets responsables du comportement dont vous avez besoin.

Ce modèle élimine une grande partie de la transmission implicite du contexte. Les opérations sur les fenêtres relèvent des fenêtres ; les opérations qui concernent toute l’application relèvent de l’application. Ce modèle est plus naturel pour les applications à plusieurs fenêtres et mieux adapté aux tests et à la maintenance des bases de code volumineuses.

Les retours des personnes ayant utilisé v3 pendant la phase alpha ont été extrêmement positifs. Les développeurs ont particulièrement apprécié le modèle explicite : il facilite la compréhension du code, clarifie les responsabilités et permet aux applications de bureau complexes de se développer sans avoir à lutter contre le framework.

### Prise en charge native de plusieurs fenêtres

La gestion de plusieurs fenêtres est une fonctionnalité fondamentale de v3. Les fenêtres ont leur propre cycle de vie et peuvent être créées, gérées et fermées pendant l’exécution. Il est ainsi plus facile de créer les types de logiciels de bureau qui nécessitent des éditeurs, des inspecteurs, des préférences, des fenêtres d’outils ou plusieurs parties indépendantes de l’interface utilisateur.

### Services et bindings générés

Les services Go remplacent l’ancien modèle de bindings. Ils conservent la logique applicative sous forme de code Go ordinaire et rendent explicite la frontière avec le frontend. Les bindings sont générés selon une structure qui reflète l’application et ses services, ce qui facilite la recherche et l’utilisation de l’API exposée au frontend.

v3 génère ces bindings par analyse statique du code source. Le générateur peut donc conserver les informations que les développeurs ont intégrées à leur code, notamment les commentaires et les noms de paramètres explicites, au lieu d’examiner par réflexion un programme déjà compilé. Il en résulte une API frontend plus riche et plus utile, ainsi qu’un processus de génération plus facile à comprendre et à maintenir.

Les services peuvent également monter des ressources frontend et des scripts aux côtés de leur code Go. Une fonctionnalité dispose ainsi d’un emplacement unique et cohérent regroupant son API backend, le code JavaScript ou l’interface utilisateur dont elle a besoin, ainsi que son point d’intégration à l’application hôte. Cela ouvre la voie à des plugins Wails offrant d’emblée des fonctionnalités riches : installez un plugin, montez son service et utilisez la fonctionnalité, au lieu d’assembler vous-même un ensemble disparate de bindings et de dépendances frontend. Cette version bêta ne comprend pas de système général de plugins, mais v3 rend cette orientation réalisable d’une manière que le modèle de bindings de v2 ne permettait pas.

### Un système de build que vous pouvez examiner et adapter

v3 rend visible la structure de build du projet. Au lieu de masquer chaque décision de build derrière une commande unique, les projets disposent d’une organisation conventionnelle et d’une configuration de build fondée sur Taskfile, que vous pouvez comprendre, étendre et déboguer en même temps que l’application.

### Une base multiplateforme plus robuste pour les applications de bureau

La version bêta prend en charge Windows sur amd64 et arm64, macOS sur Intel et Apple Silicon, ainsi que Linux sur amd64 et arm64. GTK4 avec WebKitGTK 6.0 constitue la pile Linux par défaut ; GTK3 reste disponible comme option héritée pendant toute la série v3.0. La prise en charge des appareils mobiles est prometteuse, mais reste expérimentale et ne fait pas partie de la garantie de compatibilité de la version bêta pour ordinateurs de bureau.

Cette version comprend également les améliorations nécessaires pour rendre l’expérience quotidienne plus fiable : un meilleur comportement sur chaque plateforme, un modèle de fenêtre plus performant, des diagnostics plus clairs, ainsi que des artefacts de version accompagnés de sommes de contrôle et d’informations de provenance.

## Migration depuis v2

v3 est une nouvelle version majeure : la migration nécessite un véritable portage, et non un simple changement de numéro de version. Les principaux changements conceptuels concernent le cycle de vie des applications et des fenêtres, les services qui remplacent les bindings liés au contexte, les API directes des applications et des fenêtres qui remplacent le package d’exécution de v2, ainsi que la régénération des bindings frontend.

Nous avons publié un [guide de migration de v2 vers v3](/migration/v2-to-v3/) qui détaille ces changements et comprend une correspondance des fonctionnalités ainsi qu’une liste de contrôle pour les tests. Pour cette version bêta, ce guide manuel constitue la procédure de migration prise en charge. Ne vous attendez pas à pouvoir convertir chaque projet v2 sans vérification : testez le résultat, portez méthodiquement vos appels à l’environnement d’exécution et conservez v2 jusqu’à ce que la nouvelle application soit prête.

Nous évaluons également un assistant de migration expérimental. Il ne fait pas partie de cette version bêta et nous ne le recommanderons qu’après l’avoir validé sur un ensemble représentatif de projets v2 réels.

## Un point franc sur le chemin parcouru

La première version alpha de v3 a été publiée le 18 janvier 2023. La phase alpha a donc duré longtemps, ce qui mérite mieux qu’une vague reconnaissance.

Le projet a énormément évolué pendant cette période. Lorsque la première version bêta de Wails v2 pour Windows est sortie en septembre 2021, le dépôt comptait environ 4000 étoiles. Lors de la sortie de la v2 en septembre 2022, il en comptait environ 10300. Aujourd’hui, il en compte plus de 35000. Cette croissance est un privilège, mais elle redéfinit aussi ce qu’implique une bonne gestion du projet : davantage d’utilisateurs dépendent des décisions relatives aux versions, davantage de contributeurs ont besoin de voies clairement définies pour participer, et une part croissante du travail consiste à rendre le projet prévisible plutôt qu’à simplement ajouter la fonctionnalité suivante.

Je n’ai pas toujours adapté nos processus avec la rapidité ou la clarté qu’exigeait cette croissance. J’en assume la responsabilité. La solution n’est pas de faire de grandes promesses ni de transformer chaque décision en cérémonie, mais d’expliquer plus clairement ce qui est pris en charge, ce qui est expérimental, comment les décisions sont prises et ce que nous allons faire ensuite.

Ce travail commence avec cette version bêta. Nous avons désormais défini des jalons explicites pour la bêta, les versions candidates et la disponibilité générale ; clarifié nos engagements en matière de compatibilité ; mis à jour notre politique de sécurité ; et instauré un processus WEP (Wails Enhancement Proposal) pour les modifications du comportement public et les nouvelles fonctionnalités. À mesure que Wails se développera, nous continuerons à améliorer la feuille de route et à réexaminer la gouvernance du projet. Notre objectif est de rendre le projet plus fiable et plus facile à enrichir, sans entraver sa progression.

Nous relançons également le [subreddit Wails](https://www.reddit.com/r/wails/) et l’administrerons activement afin d’offrir un autre espace de discussions pratiques, de questions et de retours à mesure que v3 se rapproche de la disponibilité générale.

## Installer et essayer la version bêta

Une fois la version publiée, installez la dernière CLI v3 avec :

```sh
go install github.com/wailsapp/wails/v3/cmd/wails3@latest
```

Avant de créer un projet, lancez l’assistant de configuration guidée. Il vérifie votre environnement de développement local et vous aide à configurer les dépendances nécessaires à Wails :

```sh
wails3 setup
```

![L’assistant de configuration de Wails v3 en mode sombre](/assets/screenshots/wails3-setup-wizard-dark.png)

### Créer un projet

Une fois la configuration terminée, créez un projet avec :

```sh
wails3 init
```

- Documentation : <https://v3.wails.io/>
- Guide de migration : <https://v3.wails.io/migration/v2-to-v3/>
- Reddit : <https://www.reddit.com/r/wails/>
- Notes de version : [ESPACE RÉSERVÉ À LA VERSION](https://github.com/wailsapp/wails/releases)

Si vous trouvez un bogue reproductible, signalez-le en joignant la sortie de `wails3 doctor` et, si possible, un exemple minimal. Si vous souhaitez proposer une nouvelle fonctionnalité ou une modification du comportement public, ouvrez une PR de WEP en mode brouillon plutôt qu’un ticket de demande de fonctionnalité. Ces deux voies nous aident à répondre clairement et à faire progresser la bêta.

## Merci

Wails v3 existe grâce aux personnes qui ont testé des versions incomplètes, signalé des bogues difficiles, traduit la documentation, répondu aux questions, contribué au code et continué à pousser le projet à s’améliorer. Merci.

Nous adressons aussi des remerciements tout particuliers et sincères aux sponsors qui ont soutenu Wails tout au long de cette longue transition. Votre soutien ne s’est pas contenté d’assurer la continuité du projet : il nous a permis de consacrer durablement du temps à l’architecture, aux outils, aux tests et à la documentation qui ont accéléré sa progression vers v3. Chaque testeur et chaque contributeur a contribué à façonner cette version, mais les sponsors nous ont permis d’accorder à ce travail l’attention qu’il méritait.

Cette version bêta est une invitation à nous aider à finaliser correctement v3. Essayez-la, développez avec elle, indiquez-nous où elle échoue et aidez-nous à parcourir rapidement et prudemment le chemin jusqu’à 3.0.
