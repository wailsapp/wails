---
title: "Sortie de Wails v2"
description: "Notes de version et annonces concernant Wails"
authors: ["leaanthony"]
tags: ["wails","v2"]
date: "2022-09-22"
slug: "blog/wails-v2-released"
image: "/assets/blog-images/montage.png"
sourcePath: "blog/wails-v2-released.md"
---

![capture d’écran du montage](/assets/blog-images/montage.png)

## Wails v2 est là !

Aujourd’hui marque la sortie de [Wails](https://wails.io) v2. Environ 18 mois se sont écoulés depuis la première version alpha de la v2, et environ un an depuis la première version bêta. Je suis profondément reconnaissant envers toutes les personnes qui ont participé à l’évolution du projet.

Si cela a pris autant de temps, c’est notamment parce que nous voulions atteindre un certain degré d’achèvement avant de l’appeler officiellement v2. En réalité, il n’existe jamais de moment idéal pour étiqueter une version : il reste toujours des problèmes à résoudre ou « juste une dernière » fonctionnalité à intégrer. Étiqueter une version majeure imparfaite apporte toutefois un peu de stabilité aux utilisateurs du projet, tout en permettant aux développeurs de repartir sur de nouvelles bases.

Cette version dépasse tout ce que j’avais imaginé. J’espère que son utilisation vous procurera autant de plaisir que son développement nous en a donné.

## Qu’*est*-ce que Wails ?

Si vous ne connaissez pas Wails, il s’agit d’un projet qui permet aux programmeurs Go de créer des interfaces riches pour leurs programmes Go à l’aide de technologies web familières. C’est une alternative légère à Electron, écrite en Go. Vous trouverez bien plus d’informations sur le [site officiel](https://wails.io/docs/introduction).

## Quoi de neuf ?

La version v2 constitue une avancée majeure pour le projet et résout de nombreuses difficultés rencontrées avec la v1. Si vous n’avez lu aucun des billets de blog consacrés aux versions bêta pour [macOS](/blog/wails-v2-beta-for-mac/), [Windows](/blog/wails-v2-beta-for-windows/) ou [Linux](/blog/wails-v2-beta-for-linux/), je vous encourage à le faire, car ils présentent plus en détail tous les changements majeurs. En résumé :

- Composant WebView2 pour Windows, compatible avec les normes web modernes et doté de fonctionnalités de débogage.
- [Thèmes sombre et clair](https://wails.io/docs/reference/options#theme), ainsi que [thèmes personnalisés](https://wails.io/docs/reference/options#customtheme) sous Windows.
- Windows ne nécessite désormais plus CGO.
- Prise en charge immédiate des modèles de projet Svelte, Vue, React, Preact, Lit et Vanilla.
- Intégration de [Vite](https://vitejs.dev/), qui fournit un environnement de développement avec rechargement à chaud pour votre application.
- [Menus](https://wails.io/docs/guides/application-development#application-menu) et [boîtes de dialogue](https://wails.io/docs/reference/runtime/dialog) natifs pour les applications.
- Effets natifs de translucidité des fenêtres pour [Windows](https://wails.io/docs/reference/options#windowistranslucent) et [macOS](https://wails.io/docs/reference/options#windowistranslucent-1). Prise en charge des arrière-plans Mica et Acrylic.
- Générez facilement un [programme d’installation NSIS](https://wails.io/docs/guides/windows-installer) pour les déploiements sous Windows.
- Une riche [bibliothèque d’exécution](https://wails.io/docs/reference/runtime/intro) qui fournit des méthodes utilitaires pour manipuler les fenêtres et gérer les événements, les boîtes de dialogue, les menus et la journalisation.
- Prise en charge de l’[obfuscation](https://wails.io/docs/guides/obfuscated) de votre application avec [garble](https://github.com/burrowers/garble).
- Prise en charge de la compression de votre application avec [UPX](https://upx.github.io/).
- Génération automatique de code TypeScript à partir des structures Go. Plus d’informations [ici](https://wails.io/docs/howdoesitwork#calling-bound-go-methods).
- Il n’est pas nécessaire de distribuer de bibliothèque ni de DLL supplémentaire avec votre application, quelle que soit la plateforme.
- Il n’est pas nécessaire d’incorporer les ressources du frontend. Développez simplement votre application comme n’importe quelle autre application web.

## Crédits et remerciements

Parvenir à la v2 a demandé un travail considérable. Entre la première version alpha et la publication d’aujourd’hui, ~2200 commits ont été réalisés par 89 contributeurs, auxquels s’ajoutent de très nombreuses personnes qui ont fourni des traductions, effectué des tests, donné leur avis et apporté leur aide sur les forums de discussion ainsi que dans l’outil de suivi des problèmes. Je suis infiniment reconnaissant envers chacun et chacune d’entre vous. Je tiens également à remercier tout particulièrement tous les sponsors du projet pour leurs orientations, leurs conseils et leurs retours. Tout ce que vous faites est extrêmement apprécié.

Je tiens à mentionner tout particulièrement quelques personnes :

Tout d’abord, un **immense** merci à [@stffabi](https://github.com/stffabi), qui a apporté tant de contributions dont nous bénéficions tous et fourni une aide considérable sur de nombreux problèmes. Il est à l’origine de fonctionnalités essentielles, telles que la prise en charge d’un serveur de développement externe, qui a transformé notre mode de développement en nous permettant d’exploiter les super-pouvoirs de [Vite](https://vitejs.dev/). Il est juste de dire que Wails v2 serait une version bien moins enthousiasmante sans ses [incroyables contributions](https://github.com/wailsapp/wails/commits?author=stffabi&since=2020-01-04). Merci infiniment, @stffabi !

Je tiens également à saluer chaleureusement [@misitebao](https://github.com/misitebao), qui assure inlassablement la maintenance du site web, fournit les traductions chinoises, gère Crowdin et aide les nouveaux traducteurs à prendre leurs marques. C’est une tâche extrêmement importante, et je suis profondément reconnaissant pour tout le temps et tous les efforts qui y ont été consacrés ! Tu assures !

Enfin, un immense merci à Mat Ryer, qui a apporté ses conseils et son soutien pendant le développement de la v2. Développer xBar ensemble à l’aide d’une première version alpha de la v2 a contribué à définir l’orientation de la v2 et m’a permis de comprendre certains défauts de conception des premières versions. Je suis heureux d’annoncer qu’à compter d’aujourd’hui, nous allons commencer à porter xBar vers Wails v2, et que xBar deviendra l’application phare du projet. Merci, Mat !

## Enseignements tirés

Le parcours jusqu’à la v2 nous a apporté plusieurs enseignements qui façonneront la suite du développement.

## Des versions plus petites, plus rapides et mieux ciblées

Au cours du développement de la v2, de nombreuses fonctionnalités et corrections de bogues ont été développées au cas par cas. Cela a allongé les cycles de publication et rendu le débogage plus difficile. À l’avenir, nous publierons plus souvent des versions comprenant un nombre réduit de fonctionnalités. Chaque version s’accompagnera d’une mise à jour de la documentation et de tests approfondis. Nous espérons que ces versions plus petites, plus rapides et mieux ciblées entraîneront moins de régressions et amélioreront la qualité de la documentation.

## Encourager la participation

Lorsque j’ai lancé ce projet, je voulais aider immédiatement chaque personne qui rencontrait un problème. Je prenais les problèmes « personnellement » et voulais les résoudre aussi vite que possible. Cette approche n’est pas viable et nuit en définitive à la pérennité du projet. À l’avenir, je laisserai davantage de place aux personnes qui souhaitent participer en répondant aux questions et en évaluant les problèmes. Il serait utile de disposer d’outils pour nous y aider. Si vous avez des suggestions, participez à la discussion [ici](https://github.com/wailsapp/wails/discussions/1855).

## Apprendre à dire non

Plus un projet open source suscite de participation, plus il reçoit de demandes de fonctionnalités supplémentaires, qui peuvent être utiles ou non à la majorité des utilisateurs. Le développement et le débogage initiaux de ces fonctionnalités demandent du temps, puis leur maintenance engendre des coûts permanents. Je suis moi-même le premier coupable, car je cherche souvent à « faire bouillir l’océan » plutôt qu’à fournir la fonctionnalité minimale viable. À l’avenir, nous devrons dire « non » un peu plus souvent à l’ajout de fonctionnalités au cœur du projet et concentrer nos efforts sur un moyen de permettre aux développeurs de fournir eux-mêmes ces fonctionnalités. Nous étudions sérieusement le recours à des plugins dans ce cas. Chacun pourra ainsi étendre le projet comme il l’entend, tout en disposant d’un moyen simple d’y contribuer.

## Regard vers l’avenir

Nous envisageons déjà d’ajouter de très nombreuses fonctionnalités essentielles à Wails lors du prochain cycle majeur de développement. La [feuille de route](https://github.com/wailsapp/wails/discussions/1484) regorge d’idées intéressantes, et j’ai hâte de commencer à travailler dessus. La prise en charge de plusieurs fenêtres fait partie des demandes les plus importantes. C’est un sujet complexe et, pour bien faire les choses, nous devrons peut-être proposer une autre API, car l’API actuelle n’a pas été conçue dans cette optique. Au vu des premières idées et des retours reçus, je pense que la direction que nous envisageons vous plaira.

La perspective de faire fonctionner des applications Wails sur les appareils mobiles m’enthousiasme tout particulièrement. Nous disposons déjà d’un projet de démonstration qui prouve qu’il est possible d’exécuter une application Wails sur Android. J’ai donc vraiment hâte d’explorer les possibilités qui s’offrent à nous !

Un dernier point que j’aimerais aborder concerne la parité des fonctionnalités. Depuis longtemps, l’un de nos principes fondamentaux consiste à ne rien ajouter au projet sans une prise en charge multiplateforme complète. Bien que cela se soit révélé jusqu’ici réalisable dans l’ensemble, ce principe a considérablement retardé la publication de nouvelles fonctionnalités. À l’avenir, nous adopterons une approche légèrement différente : toute nouvelle fonctionnalité qui ne peut pas être publiée immédiatement sur toutes les plateformes le sera au moyen d’une configuration ou d’une API expérimentale. Les premiers utilisateurs de certaines plateformes pourront ainsi l’essayer et fournir des retours qui contribueront à sa conception définitive. Cela signifie bien sûr que la stabilité de l’API ne sera pas garantie tant que la fonctionnalité ne sera pas entièrement prise en charge par toutes les plateformes sur lesquelles elle peut l’être, mais cette approche permettra au moins de débloquer son développement.

## Le mot de la fin

Je suis vraiment fier de ce que nous avons réussi à accomplir avec la version V2. C’est formidable de voir ce que les utilisateurs ont déjà pu créer jusqu’à présent avec les versions bêta : des applications de qualité comme [Varly](https://varly.app/), [Surge](https://getsurge.io/) et [October](https://october.utf9k.net/). Je vous invite à les découvrir.

Cette version a vu le jour grâce au travail acharné de nombreux contributeurs. Son téléchargement et son utilisation sont gratuits, mais sa création n’a pas été sans coût. Ne vous y trompez pas : ce projet a représenté un coût considérable. Il ne s’agit pas seulement du temps que j’y ai consacré et de celui de chacune et chacun des contributeurs, mais aussi du temps que toutes ces personnes ont passé loin de leurs proches et de leur famille. C’est pourquoi j’éprouve une immense gratitude pour chaque seconde consacrée à la réalisation de ce projet. Plus nous aurons de contributeurs, plus cet effort pourra être partagé et plus nous pourrons accomplir de choses ensemble. Je vous invite toutes et tous à trouver une façon de contribuer, qu’il s’agisse de confirmer le bogue signalé par quelqu’un, de proposer un correctif, de modifier la documentation ou d’aider une personne qui en a besoin. Toutes ces petites actions ont un impact considérable ! Ce serait formidable que vous participiez, vous aussi, à l’aventure qui nous mènera à v3.

Profitez-en bien !

&dash; Lea

P.-S. : si vous ou votre entreprise trouvez Wails utile, pensez à [parrainer le projet](https://github.com/sponsors/leaanthony). Merci !
