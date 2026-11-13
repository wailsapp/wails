---
title: "Publication de Wails v2"
description: "Notes de version et annonces concernant Wails"
authors: ["leaanthony"]
tags: ["wails","v2"]
date: "2022-09-22"
slug: "blog/wails-v2-released"
image: "/assets/blog-images/montage.png"
sourcePath: "blog/wails-v2-released.md"
---

![capture d’écran du montage](/assets/blog-images/montage.png)

## Ça y est !

Aujourd’hui marque la sortie de [Wails](https://wails.io) v2. Environ 18 mois se sont écoulés depuis la première version alpha de v2, et environ un an depuis la première version bêta. Je suis profondément reconnaissant envers toutes les personnes qui ont contribué à l’évolution du projet.

Si cela a pris autant de temps, c’est notamment parce que nous souhaitions atteindre un certain niveau d’achèvement avant de l’appeler officiellement v2. En réalité, il n’existe jamais de moment idéal pour baliser une version : il reste toujours des problèmes en suspens ou « juste une dernière » fonctionnalité à y intégrer. Cependant, le balisage d’une version majeure imparfaite apporte une certaine stabilité aux utilisateurs du projet, tout en permettant aux développeurs de repartir sur de nouvelles bases.

Cette version dépasse tout ce que j’avais imaginé. J’espère que vous prendrez autant de plaisir à l’utiliser que nous en avons pris à la développer.

## Qu’*est*-ce que Wails ?

Si vous ne connaissez pas Wails, il s’agit d’un projet qui permet aux développeurs Go de créer des interfaces riches pour leurs programmes Go à l’aide de technologies web familières. C’est une alternative légère à Electron, écrite en Go. Vous trouverez bien plus d’informations sur le [site officiel](https://v2.wails.io/docs/introduction).

## Quoi de neuf ?

La version v2 représente une avancée considérable pour le projet et résout de nombreuses difficultés rencontrées avec v1. Si vous n’avez lu aucun des billets de blog consacrés aux versions bêta pour [macOS](/blog/wails-v2-beta-for-mac/), [Windows](/blog/wails-v2-beta-for-windows/) ou [Linux](/blog/wails-v2-beta-for-linux/), je vous encourage à le faire, car ils présentent plus en détail tous les changements majeurs. En résumé :

- Composant Webview2 pour Windows, compatible avec les normes web modernes et doté de fonctionnalités de débogage.
- [Thèmes sombre et clair](https://v2.wails.io/docs/reference/options#theme), ainsi que [thèmes personnalisés](https://v2.wails.io/docs/reference/options#customtheme) sous Windows.
- Windows ne requiert désormais plus CGO.
- Prise en charge immédiate des modèles de projet Svelte, Vue, React, Preact, Lit et Vanilla.
- Intégration de [Vite](https://vitejs.dev/), qui fournit un environnement de développement avec rechargement à chaud pour votre application.
- [Menus](https://v2.wails.io/docs/guides/application-development#application-menu) et [boîtes de dialogue](https://v2.wails.io/docs/reference/runtime/dialog) d’application natifs.
- Effets natifs de transparence des fenêtres pour [Windows](https://v2.wails.io/docs/reference/options#windowistranslucent) et [macOS](https://v2.wails.io/docs/reference/options#windowistranslucent-1). Prise en charge des arrière-plans Mica et Acrylic.
- Générez facilement un [programme d’installation NSIS](https://v2.wails.io/docs/guides/windows-installer) pour les déploiements sous Windows.
- Une riche [bibliothèque d’exécution](https://v2.wails.io/docs/reference/runtime/intro) fournissant des méthodes utilitaires pour manipuler les fenêtres et gérer les événements, les boîtes de dialogue, les menus et la journalisation.
- Prise en charge de l’[obscurcissement](https://v2.wails.io/docs/guides/obfuscated) de votre application à l’aide de [garble](https://github.com/burrowers/garble).
- Prise en charge de la compression de votre application à l’aide de [UPX](https://upx.github.io/).
- Génération automatique de TypeScript à partir des structures Go. Plus d’informations [ici](https://v2.wails.io/docs/howdoesitwork#calling-bound-go-methods).
- Aucune bibliothèque ni DLL supplémentaire ne doit être distribuée avec votre application, quelle que soit la plateforme.
- Il n’est pas nécessaire de regrouper les ressources de l’interface. Développez simplement votre application comme n’importe quelle autre application web.

## Remerciements

Parvenir à v2 a demandé un effort considérable. Entre la première version alpha et la version publiée aujourd’hui, 89 contributeurs ont effectué environ 2200 commits. Beaucoup, beaucoup d’autres personnes ont fourni des traductions, réalisé des tests, formulé des retours et apporté leur aide sur les forums de discussion ainsi que dans le système de suivi des problèmes. Je suis infiniment reconnaissant envers chacune et chacun d’entre vous. Je tiens également à remercier tout particulièrement les sponsors du projet, qui nous ont apporté conseils, orientations et retours. Tout ce que vous faites est extrêmement apprécié.

Je souhaite mentionner tout particulièrement quelques personnes :

Tout d’abord, un **immense** merci à [@stffabi](https://github.com/stffabi), dont les très nombreuses contributions profitent à tout le monde et qui a également apporté une aide considérable sur de nombreux problèmes. Il a fourni plusieurs fonctionnalités essentielles, notamment la prise en charge d’un serveur de développement externe, qui a transformé notre mode de développement en nous permettant de tirer parti des super-pouvoirs de [Vite](https://vitejs.dev/). On peut affirmer sans exagérer que la version Wails v2 serait bien moins enthousiasmante sans ses [incroyables contributions](https://github.com/wailsapp/wails/commits?author=stffabi&since=2020-01-04). Merci infiniment, @stffabi !

Je tiens également à remercier chaleureusement [@misitebao](https://github.com/misitebao), qui assure sans relâche la maintenance du site web, fournit les traductions chinoises, gère Crowdin et aide les nouveaux traducteurs à prendre leurs marques. Cette tâche est extrêmement importante et je suis profondément reconnaissant pour tout le temps et tous les efforts qui lui ont été consacrés ! Tu assures !

Enfin, et surtout, un immense merci à Mat Ryer, qui a prodigué conseils et soutien pendant le développement de v2. Le développement conjoint de xBar à l’aide d’une première version alpha de v2 a contribué à orienter v2 et m’a permis de mieux comprendre certaines failles de conception des premières versions. Je suis heureux d’annoncer qu’à compter d’aujourd’hui, nous allons commencer à porter xBar vers Wails v2, qui deviendra l’application phare du projet. Merci, Mat !

## Leçons tirées

Le chemin parcouru pour parvenir à v2 nous a enseigné plusieurs leçons qui orienteront la suite du développement.

## Des versions plus petites, plus rapides et mieux ciblées

Au cours du développement de v2, de nombreuses fonctionnalités et corrections de bogues ont été développées de manière ponctuelle. Cela a allongé les cycles de publication et compliqué le débogage. À l’avenir, nous publierons plus souvent des versions comportant moins de fonctionnalités. Chaque publication comprendra la mise à jour de la documentation ainsi que des tests approfondis. Nous espérons que ces versions plus petites, plus rapides et mieux ciblées entraîneront moins de régressions et amélioreront la qualité de la documentation.

## Encourager la participation

Lorsque j’ai lancé ce projet, je voulais aider immédiatement chaque personne rencontrant un problème. Je prenais les problèmes à cœur et souhaitais les résoudre aussi vite que possible. Cette approche n’est pas viable et finit par nuire à la pérennité du projet. À l’avenir, je laisserai davantage de place aux personnes souhaitant participer aux réponses aux questions et au triage des problèmes. Il serait utile de disposer d’outils pour nous y aider. Si vous avez des suggestions, rejoignez la discussion [ici](https://github.com/wailsapp/wails/discussions/1855).

## Apprendre à dire non

Plus un projet open source suscite de participation, plus il reçoit de demandes de fonctionnalités supplémentaires, qui peuvent ou non être utiles à la majorité. Le développement et le débogage initiaux de ces fonctionnalités prennent du temps, puis elles engendrent des coûts de maintenance permanents. Je suis moi-même le premier coupable : je cherche souvent à « faire bouillir l’océan » au lieu de proposer la fonctionnalité minimale viable. À l’avenir, nous devrons dire un peu plus souvent « non » à l’ajout de fonctionnalités au cœur du projet et consacrer nos efforts à donner aux développeurs les moyens de fournir eux-mêmes ces fonctionnalités. Nous étudions sérieusement le recours aux plugins dans ce cas. Chacun pourra ainsi étendre le projet selon ses besoins et y contribuer facilement.

## Regard vers l’avenir

Nous envisageons déjà d’ajouter de très nombreuses fonctionnalités essentielles à Wails au cours du prochain cycle majeur de développement. La [feuille de route](https://github.com/wailsapp/wails/discussions/1484) regorge d’idées intéressantes et j’ai hâte de commencer à travailler dessus. La prise en charge de plusieurs fenêtres fait partie des demandes les plus importantes. C’est un sujet délicat et, pour bien faire les choses, nous devrons peut-être proposer une autre API, car l’API actuelle n’a pas été conçue dans cette optique. Au vu des premières idées et des premiers retours, je pense que l’orientation que nous envisageons vous plaira.

À titre personnel, la perspective de faire fonctionner des applications Wails sur les appareils mobiles m’enthousiasme beaucoup. Nous disposons déjà d’un projet de démonstration prouvant qu’il est possible d’exécuter une application Wails sous Android ; j’ai donc vraiment hâte d’explorer les possibilités qui s’offrent à nous !

Un dernier point que j’aimerais aborder est celui de la parité fonctionnelle. Depuis longtemps, l’un des principes fondamentaux du projet est de ne rien y ajouter sans une prise en charge multiplateforme complète. Bien que cela se soit révélé jusqu’ici réalisable dans l’ensemble, ce principe a fortement freiné la publication de nouvelles fonctionnalités. À l’avenir, nous adopterons une approche légèrement différente : toute nouvelle fonctionnalité qui ne pourra pas être publiée immédiatement sur toutes les plateformes le sera au moyen d’une configuration ou d’une API expérimentale. Les utilisateurs précoces de certaines plateformes pourront ainsi l’essayer et fournir des retours qui contribueront à sa conception définitive. Bien entendu, cela signifie que la stabilité de l’API ne sera pas garantie tant que la fonctionnalité ne sera pas entièrement prise en charge sur toutes les plateformes compatibles, mais cela permettra au moins de débloquer son développement.

## Le mot de la fin

Je suis vraiment fier de ce que nous avons accompli avec la version V2. C’est formidable de voir ce que certains ont déjà réussi à créer jusqu’ici avec les versions bêta : des applications de qualité comme [Varly](https://varly.app/), [Surge](https://getsurge.io/) et [October](https://october.utf9k.net/). Je vous invite à les découvrir.

Cette version a vu le jour grâce au travail acharné de nombreux contributeurs. Bien qu’elle soit téléchargeable et utilisable gratuitement, sa réalisation n’a pas été sans coût. Ne vous y trompez pas : ce projet a représenté un coût considérable. Il ne s’agit pas seulement du temps que j’y ai consacré et de celui de chacun des contributeurs, mais aussi du temps que toutes ces personnes n’ont pas pu passer avec leurs proches. C’est pourquoi je suis extrêmement reconnaissant pour chaque seconde consacrée à la réalisation de ce projet. Plus nous aurons de contributeurs, plus cet effort pourra être réparti et plus nous pourrons accomplir de choses ensemble. Je vous invite tous à choisir une façon de contribuer, qu’il s’agisse de confirmer le bogue signalé par quelqu’un, de proposer un correctif, de modifier la documentation ou d’aider une personne qui en a besoin. Toutes ces petites contributions ont un impact considérable ! Ce serait vraiment formidable que vous participiez vous aussi à l’aventure qui nous mènera à v3.

Profitez-en bien !

&dash; Lea

P.-S. : si vous ou votre entreprise trouvez Wails utile, envisagez de [parrainer le projet](https://github.com/sponsors/leaanthony). Merci !
