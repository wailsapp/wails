---
title: "Wails v2 bêta pour Windows"
description: "Notes de version et annonces de Wails"
authors: ["leaanthony"]
tags: ["wails","v2"]
date: "2021-09-27"
slug: "blog/wails-v2-beta-for-windows"
image: "/assets/blog-images/wails.webp"
sourcePath: "blog/wails-v2-beta-for-windows.md"
---

![capture d’écran de Wails](/assets/blog-images/wails.webp)

Lorsque j’ai annoncé Wails pour la première fois sur Reddit, il y a un peu plus de 2 ans, depuis un train à Sydney, je ne m’attendais pas à ce que le projet suscite beaucoup d’intérêt. Quelques jours plus tard, un vlogueur technique prolifique a publié un tutoriel vidéo et en a fait une critique positive. À partir de là, l’intérêt pour le projet a explosé.

Il était clair que la possibilité d’ajouter des interfaces web à des projets Go enthousiasmait les utilisateurs, qui ont presque immédiatement poussé le projet au-delà de la preuve de concept que j’avais créée. À l’époque, Wails utilisait le projet [webview](https://github.com/webview/webview) pour gérer l’interface, et le moteur de rendu IE11 était la seule option sous Windows. Cette limitation était à l’origine de nombreux rapports de bogues : mauvaise prise en charge de JavaScript et de CSS, et absence d’outils de développement pour le débogage. L’expérience de développement était frustrante, mais il n’était guère possible d’y remédier.

Pendant longtemps, j’ai été fermement convaincu que Microsoft finirait par devoir régler son problème de navigateur. Le monde avançait, le développement frontend était en plein essor et IE n’était plus à la hauteur. Lorsque Microsoft a annoncé que Chromium servirait de base à sa nouvelle stratégie en matière de navigateur, j’ai su que Wails finirait par pouvoir l’utiliser et ferait ainsi passer l’expérience des développeurs Windows au niveau supérieur.

Aujourd’hui, j’ai le plaisir de vous annoncer : **Wails v2 bêta pour Windows** ! Cette version apporte énormément de nouveautés à découvrir. Prenez donc un verre, installez-vous et commençons…

## Aucune dépendance à CGO !

Non, ce n’est pas une plaisanterie : *aucune* *dépendance* à *CGO* 🤯 ! Contrairement à macOS et Linux, Windows n’est pas fourni avec un compilateur par défaut. En outre, CGO nécessite un compilateur mingw, pour lequel il existe une multitude d’options d’installation. La suppression de cette exigence a considérablement simplifié la configuration et grandement facilité le débogage. Même si j’ai consacré beaucoup d’efforts à faire fonctionner tout cela, l’essentiel du mérite revient à [John Chadwick](https://github.com/jchv), non seulement pour avoir lancé plusieurs projets qui ont rendu cela possible, mais aussi pour avoir accepté que quelqu’un les reprenne et les développe. Merci également à [Tad Vizbaras](https://github.com/tadvi), dont le projet [winc](https://github.com/tadvi/winc) m’a engagé sur cette voie.

### Moteur de rendu Chromium WebView2

![capture d’écran des outils de développement](/assets/blog-images/devtools.png)

Enfin, les développeurs Windows disposent d’un moteur de rendu de premier ordre pour leurs applications ! Fini le temps où il fallait contorsionner le code de votre interface pour qu’il fonctionne sous Windows. De plus, vous bénéficiez d’outils de développement de premier ordre !

Le composant WebView2 exige toutefois que `WebView2Loader.dll` se trouve à côté du binaire. Cela rend la distribution un peu plus pénible que ce à quoi nous autres gophers sommes habitués. Toutes les solutions et bibliothèques utilisant WebView2 que je connais ont cette dépendance.

Cependant, je suis vraiment ravi d’annoncer que les applications Wails *n’ont aucune exigence de ce type* ! Grâce aux prouesses de [John Chadwick](https://github.com/jchv), nous pouvons intégrer cette DLL au binaire et faire en sorte que Windows la charge comme si elle était présente sur le disque.

Réjouissez-vous, gophers ! Le rêve du binaire unique perdure !

### Nouvelles fonctionnalités

![capture d’écran des menus Wails](/assets/blog-images/wails-menus.webp)

La prise en charge des menus natifs a fait l’objet de nombreuses demandes. Wails répond enfin à ce besoin. Les menus d’application sont désormais disponibles et prennent en charge la plupart des fonctionnalités des menus natifs, notamment les éléments de menu standard, les cases à cocher, les groupes de boutons radio, les sous-menus et les séparateurs.

Dans la v1, de très nombreux utilisateurs ont demandé à pouvoir mieux contrôler la fenêtre elle-même. J’ai le plaisir d’annoncer que de nouvelles API d’exécution sont spécialement prévues à cet effet. Elles offrent de nombreuses fonctionnalités et prennent en charge les configurations à plusieurs moniteurs. L’API de boîtes de dialogue a également été améliorée : vous pouvez maintenant utiliser des boîtes de dialogue modernes et natives, assorties de nombreuses options de configuration pour répondre à tous vos besoins.

Vous pouvez désormais générer la configuration de l’IDE avec votre projet. Ainsi, lorsque vous ouvrez celui-ci dans un IDE compatible, il est déjà configuré pour compiler et déboguer votre application. VSCode est actuellement pris en charge, mais nous espérons bientôt en faire autant pour d’autres IDE tels que Goland.

![capture d’écran de VSCode](/assets/blog-images/vscode.webp)

### Aucune obligation de regrouper les ressources

L’un des principaux points faibles de la v1 était la nécessité de condenser toute l’application en un seul fichier JS et un seul fichier CSS. J’ai le plaisir d’annoncer qu’avec la v2, vous n’avez absolument plus besoin de regrouper les ressources. Vous souhaitez charger une image locale ? Utilisez une balise `<img>` avec un chemin src local. Vous souhaitez utiliser une police attrayante ? Copiez-la et ajoutez son chemin dans votre CSS.

> Tiens, on dirait un serveur web…

Oui, cela fonctionne exactement comme un serveur web, sauf que ce n’en est pas un.

> Alors, comment inclure mes ressources ?

Il vous suffit de transmettre à la configuration de votre application un unique `embed.FS` contenant toutes vos ressources. Elles n’ont même pas besoin de se trouver dans le répertoire racine : Wails s’en charge pour vous.

### Nouvelle expérience de développement

![capture d’écran du navigateur](/assets/blog-images/browser.webp)

Comme les ressources n’ont désormais plus besoin d’être regroupées, une toute nouvelle expérience de développement devient possible. La nouvelle commande `wails dev` compile et exécute votre application, mais au lieu d’utiliser les ressources du `embed.FS`, elle les charge directement depuis le disque.

Elle offre également les fonctionnalités supplémentaires suivantes :

- Rechargement à chaud : toute modification des ressources frontend déclenche le rechargement automatique de l’interface de l’application
- Recompilation automatique : toute modification de votre code Go entraîne la recompilation et le redémarrage de votre application

En outre, un serveur web démarre sur le port 34115. Il fournit votre application à tout navigateur qui s’y connecte. Tous les navigateurs web connectés réagissent aux événements système, tels que le rechargement à chaud après la modification d’une ressource.

En Go, nous avons l’habitude de manipuler des structs dans nos applications. Il est souvent utile d’envoyer ces structs à l’interface et de les utiliser comme état de l’application. Dans la v1, ce processus était très manuel et quelque peu contraignant pour le développeur. J’ai le plaisir d’annoncer que, dans la v2, toute application exécutée en mode développement génère automatiquement des modèles TypeScript pour toutes les structs utilisées comme paramètres d’entrée ou de sortie de méthodes liées. Cela permet un échange fluide des modèles de données entre les deux environnements.

En outre, un autre module JS est généré dynamiquement pour encapsuler toutes vos méthodes liées. Il fournit la documentation JSDoc de vos méthodes, ce qui permet à votre IDE de proposer la complétion du code et des indications. C’est vraiment impressionnant de voir les modèles de données être importés automatiquement lorsque vous appuyez sur la touche Tab dans un module généré automatiquement qui encapsule votre code Go !

### Modèles distants

![capture d’écran distante](/assets/blog-images/remote.webp)

Permettre de démarrer rapidement une application a toujours été un objectif essentiel du projet Wails. Lors du lancement, nous avons essayé de prendre en charge de nombreux frameworks modernes de l’époque : react, vue et angular. Le monde du développement frontend est très clivé, évolue rapidement et il est difficile d’en suivre le rythme ! Par conséquent, nos modèles de base devenaient rapidement obsolètes, ce qui compliquait leur maintenance. Cela signifiait également que nous ne disposions pas de modèles modernes et attrayants pour les toutes dernières piles technologiques.

Avec la v2, je souhaitais donner davantage de moyens à la communauté en vous permettant de créer et d’héberger vous-mêmes des modèles, au lieu de dépendre du projet Wails. Vous pouvez donc désormais créer des projets à partir de modèles pris en charge par la communauté ! J’espère que cela incitera les développeurs à créer un écosystème dynamique de modèles de projet. J’ai vraiment hâte de découvrir ce que notre communauté de développeurs saura créer !

### En conclusion

Wails v2 constitue une nouvelle fondation pour le projet. Cette version vise à recueillir des retours sur la nouvelle approche et à corriger les éventuels bogues avant la publication d’une version définitive. Vos avis seront les bienvenus. Veuillez transmettre vos retours sur le forum de discussion [Bêta de la v2](https://github.com/wailsapp/wails/discussions/828).

Le chemin jusqu’ici a connu de nombreux rebondissements, changements de cap et demi-tours. Cela tient en partie à des décisions techniques initiales qu’il a fallu revoir, et en partie au fait que certains problèmes fondamentaux, pour lesquels nous avions consacré du temps à élaborer des solutions de contournement, ont été corrigés en amont : la fonctionnalité embed de Go en est un bon exemple. Heureusement, tout s’est mis en place au bon moment et nous disposons aujourd’hui de la meilleure solution possible. Je pense que l’attente en valait la peine : cela n’aurait même pas été possible il y a 2 mois.

Je tiens également à remercier chaleureusement :pray: les personnes suivantes, car sans elles, cette version n’existerait tout simplement pas :

- [Misite Bao](https://github.com/misitebao) — Une force de travail exceptionnelle pour les traductions chinoises et une personne incroyablement douée pour débusquer les bogues.
- [John Chadwick](https://github.com/jchv) — Son travail remarquable sur [go-webview2](https://github.com/jchv/go-webview2) et [go-winloader](https://github.com/jchv/go-winloader) a rendu possible la version Windows dont nous disposons aujourd’hui.
- [Tad Vizbaras](https://github.com/tadvi) — Les expérimentations menées avec son projet [winc](https://github.com/tadvi/winc) ont constitué la première étape vers une version de Wails entièrement écrite en Go.
- [Mat Ryer](https://github.com/matryer) — Son soutien, ses encouragements et ses retours ont réellement contribué à faire avancer le projet.

Enfin, je tiens à remercier tout particulièrement tous les [sponsors du projet](/credits/#sponsors), notamment [JetBrains](https://www.jetbrains.com?from=Wails), dont le soutien contribue de multiples façons, en coulisses, à faire avancer le projet.

J’ai hâte de découvrir ce que chacun créera avec Wails au cours de cette nouvelle phase passionnante du projet !

Lea.

PS : Les utilisateurs de MacOS et de Linux ne doivent pas se sentir oubliés : le portage vers cette nouvelle fondation est en cours et la majeure partie du travail difficile a déjà été accomplie. Encore un peu de patience !

P.-P.-S. : Si vous ou votre entreprise trouvez Wails utile, pensez à [soutenir financièrement le projet](https://github.com/sponsors/leaanthony). Merci !
