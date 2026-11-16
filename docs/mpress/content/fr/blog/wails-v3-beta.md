---
title: "Wails v3 bêta : une nouvelle base pour les applications de bureau en Go"
description: "Wails v3 bêta introduit un modèle d’application plus direct, des liaisons plus riches et une base plus claire pour les applications de bureau en Go."
authors: ["leaanthony"]
tags: ["wails","v3","beta"]
date: "2026-08-02"
slug: "blog/wails-v3-beta"
image: "/assets/screenshots/frameless-v3-native-corners-macos.png"
sourcePath: "blog/wails-v3-beta.md"
---

![Une fenêtre Wails v3 native sans cadre sous macOS](/assets/screenshots/frameless-v3-native-corners-macos.png)

Nous publions aujourd’hui Wails v3 bêta.

Wails permet aux développeurs Go de créer des applications de bureau avec les outils frontend web qu’ils connaissent déjà, en utilisant la WebView native de chaque plateforme plutôt qu’un navigateur intégré. v3 constitue une avancée majeure : cette version fournit aux applications une API plus directe, un modèle de compilation plus clair et une meilleure base pour les applications de bureau dont les utilisateurs demandent la prise en charge par Wails.

Il s’agit d’une version bêta, et non de la version 3.0 finale. L’API de bureau est stable et des équipes utilisent déjà v3 en production, mais vous devez effectuer des tests approfondis avant tout déploiement. Nous profitons de la période bêta pour résoudre avec la communauté les derniers problèmes de compatibilité et de flux de travail. Wails v2 reste la version stable actuelle et continuera à recevoir des correctifs.

## Documentation pendant la phase bêta

Pendant la phase bêta, nous maintenons la documentation anglaise comme source de référence tandis que l’API et les flux de travail font l’objet de leur validation finale. Nous n’acceptons pas encore les PR de traduction. Le travail de traduction reprendra avant la disponibilité générale, lorsque la documentation sera suffisamment stable pour que les traducteurs puissent travailler sans subir de remaniements répétés.

## Contenu de v3

- Une API explicite pour les applications et les fenêtres, avec notamment une prise en charge native de plusieurs fenêtres
- Des services Go dont l’analyse statique du code source génère des liaisons TypeScript plus riches, tout en préservant les commentaires et les noms de paramètres explicites
- Des services capables de regrouper des ressources et des scripts frontend avec leur API backend, constituant ainsi la base de plugins installables plus riches
- Un système de compilation visible, fondé sur Taskfile, que vous pouvez inspecter, étendre et déboguer
- Des compilations serveur permettant d’exécuter la même application et les mêmes services sans fenêtre de bureau native
- Une prise en charge moderne des environnements de bureau macOS, Windows et Linux sur Intel, Apple Silicon, amd64 et arm64, lorsque ces architectures sont prises en charge
- Une prise en charge mobile expérimentale pour iOS et Android, disponible à des fins d’exploration, mais exclue de l’engagement de compatibilité de la version bêta pour ordinateurs de bureau

## Pourquoi v3

Wails v2 a simplifié la création d’applications Go dotées d’un frontend web moderne. Cette version a rendu de grands services au projet, ainsi qu’à un très grand nombre d’applications. Toutefois, son environnement d’exécution à fenêtre unique piloté par le contexte et son processus de compilation étroitement encadré compliquaient certaines tâches courantes de développement d’applications de bureau plus qu’ils ne l’auraient dû.

v3 repose sur un modèle différent. Les applications, fenêtres, services, événements et fonctionnalités propres aux plateformes sont des objets explicites. Il devient ainsi plus facile de comprendre le framework à mesure qu’une application se développe, et des fonctionnalités telles que les fenêtres multiples font naturellement partie du modèle d’application au lieu de nécessiter une solution de contournement.

## Nouveautés

### Une API d’application conçue pour de véritables logiciels de bureau

v3 remplace le style de configuration v2 `wails.Run(...)` par un cycle de vie explicite de l’application. Vous créez une application, enregistrez des services, créez des fenêtres et interagissez avec les objets qui détiennent le comportement dont vous avez besoin.

Cela élimine une grande partie de la transmission implicite du contexte. Les opérations sur les fenêtres appartiennent aux fenêtres ; celles qui concernent l’ensemble de l’application appartiennent à l’application. Ce modèle est plus naturel pour les applications à plusieurs fenêtres et mieux adapté aux tests ainsi qu’à la maintenance des bases de code plus importantes.

Les retours des personnes ayant utilisé v3 pendant la phase alpha ont été extrêmement positifs. Les développeurs ont particulièrement apprécié le modèle explicite : il facilite la compréhension du code, clarifie la responsabilité des objets et permet aux applications de bureau complexes d’évoluer sans avoir à lutter contre le framework.

### Prise en charge native de plusieurs fenêtres

La gestion de plusieurs fenêtres est une fonctionnalité essentielle de v3. Chaque fenêtre possède son propre cycle de vie et peut être créée, gérée et fermée pendant l’exécution. Il en résulte une voie plus claire pour créer le type de logiciel de bureau qui nécessite des éditeurs, des inspecteurs, des préférences, des fenêtres d’outils ou plusieurs éléments d’interface utilisateur indépendants.

### Services et liaisons générées

Les services Go remplacent l’ancien modèle de liaison. Ils conservent la logique applicative sous forme de code Go ordinaire et rendent explicite la frontière avec le frontend. Les liaisons sont générées selon une structure qui reflète l’application et ses services, ce qui facilite la recherche et l’utilisation de l’API exposée au frontend.

v3 génère ces liaisons par analyse statique du code source. Le générateur peut donc conserver les informations que les développeurs ont ajoutées à leur code, notamment les commentaires et les noms de paramètres explicites, au lieu d’examiner par réflexion un programme déjà compilé. Il en résulte une API frontend plus riche et plus utile, ainsi qu’un processus de génération plus facile à comprendre et à maintenir.

Les services peuvent également monter des ressources et des scripts frontend aux côtés de leur code Go. Une fonctionnalité dispose ainsi d’un emplacement unique et cohérent réunissant son API backend, le JavaScript ou l’interface utilisateur dont elle a besoin et le point d’intégration destiné à l’application hôte. Cela ouvre la voie à des plugins Wails qui fournissent immédiatement des fonctionnalités riches : installez un plugin, montez son service et utilisez la fonctionnalité, au lieu d’assembler vous-même un ensemble disparate de liaisons et de dépendances frontend. Cette version bêta ne comprend pas de système général de plugins, mais v3 rend cette orientation réalisable, contrairement au modèle de liaison de v2.

### Un système de compilation que vous pouvez inspecter et adapter

v3 rend visible la structure de compilation du projet. Au lieu de masquer chaque décision de compilation derrière une seule commande, les projets adoptent une organisation conventionnelle et une configuration de compilation fondée sur Taskfile, que vous pouvez comprendre, étendre et déboguer avec l’application.

### Une base multiplateforme plus robuste pour les applications de bureau

La version bêta prend en charge Windows sur amd64 et arm64, macOS sur Intel et Apple Silicon, ainsi que Linux sur amd64 et arm64. GTK4 avec WebKitGTK 6.0 constitue la pile Linux par défaut ; GTK3 reste disponible comme option héritée pendant toute la série v3.0. La prise en charge des plateformes mobiles est prometteuse, mais reste expérimentale et ne fait pas partie de l’engagement de compatibilité de la version bêta pour ordinateurs de bureau.

Cette version comprend également les améliorations nécessaires pour rendre l’utilisation quotidienne plus fiable : un meilleur comportement sur chaque plateforme, un modèle de fenêtre plus performant, des diagnostics plus clairs et des artefacts de publication accompagnés de sommes de contrôle et d’informations de provenance.

## Migration depuis v2

v3 est une nouvelle version majeure : sa migration constitue un véritable portage et non un simple changement de numéro de version. Les principaux changements conceptuels concernent le cycle de vie des applications et des fenêtres, l’emploi de services à la place de liaisons associées au contexte, l’utilisation directe des API d’application et de fenêtre à la place du paquet d’exécution v2, ainsi que la régénération des liaisons frontend.

Nous avons publié un [guide de migration de v2 vers v3](/migration/v2-to-v3/) qui détaille ces changements et comprend une table de correspondance des fonctionnalités ainsi qu’une liste de contrôle pour les tests. Pour cette version bêta, ce guide manuel constitue la méthode de migration prise en charge. Ne vous attendez pas à pouvoir convertir tous les projets v2 sans vérification : testez le résultat, portez délibérément vos appels à l’environnement d’exécution et conservez v2 jusqu’à ce que la nouvelle application soit prête.

Nous évaluons également un assistant de migration expérimental. Il ne fait pas partie de cette version bêta et nous ne le recommanderons qu’après l’avoir validé sur des projets v2 réels et représentatifs.

## Un regard franc sur le chemin parcouru

La première balise alpha de v3 a été publiée le 18 janvier 2023. Cette phase alpha a duré longtemps et mérite plus qu’une vague reconnaissance.

Le projet a énormément changé pendant cette période. Lorsque la première version bêta de Wails v2 pour Windows est sortie en septembre 2021, le dépôt comptait environ 4000 étoiles. Lors de la publication de v2 en septembre 2022, il en comptait environ 10300. Il en compte aujourd’hui plus de 35000. Cette croissance est un privilège, mais elle transforme aussi ce que signifie une bonne gestion du projet : davantage d’utilisateurs dépendent des décisions de publication, davantage de contributeurs ont besoin de voies de participation claires et une plus grande part du travail consiste à rendre le projet prévisible plutôt qu’à simplement ajouter la fonctionnalité suivante.

Je n’ai pas toujours adapté nos processus aussi rapidement ni aussi clairement que cette croissance l’exigeait. J’en assume la responsabilité. La réponse n’est pas de faire de grandes promesses ni de transformer chaque décision en cérémonie, mais d’indiquer plus clairement ce qui est pris en charge, ce qui est expérimental, comment les décisions sont prises et ce que nous allons faire ensuite.

Ce travail commence avec cette version bêta. Nous disposons désormais de jalons explicites pour la bêta, les versions candidates et la disponibilité générale, d’engagements plus clairs en matière de compatibilité, d’une politique de sécurité mise à jour et d’un processus WEP (Wails Enhancement Proposal) pour les modifications du comportement public et les nouvelles fonctionnalités. Nous continuerons à améliorer la feuille de route et à revoir la gouvernance du projet à mesure que Wails se développe. L’objectif est de rendre le projet plus fiable et plus accessible aux contributions, sans entraver son évolution.

Nous relançons également le [subreddit Wails](https://www.reddit.com/r/wails/) et le maintiendrons activement afin d’offrir un espace supplémentaire pour les discussions pratiques, les questions et les retours à mesure que v3 se rapproche de la disponibilité générale.

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
- Notes de version : [VERSION PLACEHOLDER](https://github.com/wailsapp/wails/releases)

Si vous trouvez un bug reproductible, veuillez le signaler en joignant la sortie de `wails3 doctor` et, si possible, un exemple minimal. Si vous souhaitez proposer une nouvelle fonctionnalité ou une modification du comportement public, commencez plutôt par une PR WEP à l’état de brouillon que par un ticket de demande de fonctionnalité. Ces deux voies nous aident à répondre clairement et à faire avancer la bêta.

## Merci

Wails v3 existe grâce aux personnes qui ont testé des versions incomplètes, signalé des bugs difficiles, traduit la documentation, répondu aux questions, contribué au code et continué à pousser le projet à s’améliorer. Merci.

Et un merci tout particulier et sincère aux sponsors qui ont soutenu Wails tout au long de cette longue transition. Votre soutien a fait bien plus que nous permettre de poursuivre le projet : il nous a donné la possibilité de consacrer durablement du temps à l’architecture, aux outils, aux tests et à la documentation qui ont permis au projet d’accélérer vers v3. Chaque personne ayant testé la version ou contribué au projet a participé à façonner cette version, mais ce sont les sponsors qui nous ont permis d’accorder à ce travail l’attention qu’il méritait.

La bêta est une invitation à nous aider à finaliser correctement v3. Essayez-la, développez avec elle, dites-nous où elle échoue et aidez-nous à parcourir le chemin vers 3.0 rapidement et avec soin.
