---
title: "Wails v2 bêta pour macOS"
description: "Notes de version et annonces concernant Wails"
authors: ["leaanthony"]
tags: ["wails","v2"]
date: "2021-11-08"
slug: "blog/wails-v2-beta-for-mac"
image: "/assets/blog-images/wails-mac.webp"
sourcePath: "blog/wails-v2-beta-for-mac.md"
---

![capture d’écran de wails-mac](/assets/blog-images/wails-mac.webp)

Aujourd’hui marque la sortie de la première version bêta de Wails v2 pour Mac ! Il nous a fallu un certain temps pour en arriver là, et j’espère que la version publiée aujourd’hui vous sera raisonnablement utile. Le chemin a été riche en rebondissements et j’espère, avec votre aide, résoudre les derniers problèmes et peaufiner le portage Mac en vue de la version finale de Wails v2.

Cela signifie que cette version n’est pas prête pour la production ? Elle l’est peut-être tout à fait pour votre cas d’usage, mais plusieurs problèmes connus subsistent. Gardez donc un œil sur [ce tableau de projet](https://github.com/wailsapp/wails/projects/7). Si vous souhaitez contribuer, votre aide sera la bienvenue !

Quelles sont donc les nouveautés de Wails v2 pour Mac par rapport à la v1 ? Indice : elles ressemblent beaucoup à celles de la version bêta pour Windows :wink:

## Nouvelles fonctionnalités

![capture d’écran des menus Mac de Wails](/assets/blog-images/wails-menus-mac.webp)

La prise en charge des menus natifs a fait l’objet de nombreuses demandes. Wails répond enfin à ce besoin. Les menus d’application sont désormais disponibles et prennent en charge la plupart des fonctionnalités des menus natifs, notamment les éléments de menu standard, les cases à cocher, les groupes de boutons radio, les sous-menus et les séparateurs.

Avec la v1, vous avez été extrêmement nombreux à demander davantage de contrôle sur la fenêtre elle-même. J’ai le plaisir d’annoncer que de nouvelles API d’exécution sont spécialement prévues à cet effet. Elles offrent de nombreuses fonctionnalités et prennent en charge les configurations à plusieurs moniteurs. L’API de boîtes de dialogue a également été améliorée : vous pouvez désormais utiliser des boîtes de dialogue modernes et natives, avec de nombreuses options de configuration pour répondre à tous vos besoins.

### Options propres au Mac

Outre les options d’application habituelles, Wails v2 pour Mac apporte également quelques fonctionnalités supplémentaires propres au Mac :

- Donnez à votre fenêtre un aspect original et translucide, comme toutes ces jolies applications Swift !
- Barre de titre hautement personnalisable
- Prise en charge des options NSAppearance pour l’application
- Configuration simple permettant de créer automatiquement un menu « À propos »

### Aucune obligation de regrouper les ressources

L’une des principales difficultés de la v1 était l’obligation de condenser toute votre application en un seul fichier JS et un seul fichier CSS. J’ai le plaisir d’annoncer qu’avec la v2, il n’est absolument plus nécessaire de regrouper les ressources. Vous souhaitez charger une image locale ? Utilisez une balise `<img>` avec un chemin src local. Vous souhaitez utiliser une police élégante ? Copiez-la et ajoutez son chemin dans votre CSS.

> Ouah, on dirait un serveur web…

Oui, cela fonctionne exactement comme un serveur web, sauf que ce n’en est pas un.

> Alors, comment inclure mes ressources ?

Il vous suffit de transmettre à la configuration de votre application un unique `embed.FS` contenant toutes vos ressources. Elles n’ont même pas besoin de se trouver dans le répertoire racine : Wails s’en charge pour vous.

### Nouvelle expérience de développement

Comme il n’est désormais plus nécessaire de regrouper les ressources, une toute nouvelle expérience de développement devient possible. La nouvelle commande `wails dev` compile et exécute votre application, mais au lieu d’utiliser les ressources du `embed.FS`, elle les charge directement depuis le disque.

Elle offre également les fonctionnalités supplémentaires suivantes :

- Rechargement à chaud — Toute modification des ressources de l’interface déclenche le rechargement automatique de l’interface de l’application
- Recompilation automatique — Toute modification de votre code Go entraîne la recompilation et le redémarrage de votre application

En outre, un serveur web démarre sur le port 34115. Il sert votre application à tous les navigateurs qui s’y connectent. Tous les navigateurs web connectés réagissent aux événements système, comme le rechargement à chaud lors de la modification d’une ressource.

En Go, nous avons l’habitude de manipuler des structures dans nos applications. Il est souvent utile d’envoyer des structures à l’interface et de les utiliser comme état de l’application. Dans la v1, cette opération était très manuelle et assez contraignante pour les développeurs. J’ai le plaisir d’annoncer qu’avec la v2, toute application exécutée en mode développement génère automatiquement des modèles TypeScript pour toutes les structures utilisées comme paramètres d’entrée ou de sortie des méthodes liées. Cela permet un échange fluide des modèles de données entre les deux environnements.

En outre, un autre module JS est généré dynamiquement pour encapsuler toutes vos méthodes liées. Il fournit la documentation JSDoc de vos méthodes, ce qui permet à votre IDE de proposer la saisie semi-automatique et des indications. C’est vraiment génial de voir les modèles de données importés automatiquement lorsque vous appuyez sur la touche de tabulation dans un module généré automatiquement qui encapsule votre code Go !

### Modèles distants

![capture d’écran de remote-mac](/assets/blog-images/remote-mac.webp)

Permettre de mettre rapidement une application en service a toujours été un objectif essentiel du projet Wails. Lors du lancement, nous avons essayé de prendre en charge de nombreux frameworks modernes de l’époque : react, vue et angular. Le développement d’interfaces est un domaine où les avis sont très tranchés, qui évolue rapidement et dont il est difficile de suivre le rythme ! Nos modèles de base sont donc devenus assez vite obsolètes, ce qui a compliqué leur maintenance. Cela signifiait également que nous ne disposions pas de modèles modernes et attrayants pour les piles technologiques les plus récentes et performantes.

Avec la v2, je souhaitais donner plus d’autonomie à la communauté en vous permettant de créer et d’héberger vous-mêmes des modèles, sans dépendre du projet Wails. Vous pouvez donc désormais créer des projets à l’aide de modèles maintenus par la communauté ! J’espère que cela incitera les développeurs à créer un écosystème dynamique de modèles de projets. J’ai vraiment hâte de découvrir ce que notre communauté de développeurs va créer !

### Prise en charge native des M1

Grâce au soutien exceptionnel de [Mat Ryer](https://github.com/matryer/), le projet Wails prend désormais en charge les builds natifs pour M1 :

![capture d’écran de build-darwin-arm](/assets/blog-images/build-darwin-arm.webp)

Vous pouvez également spécifier `darwin/amd64` comme cible :

![capture d’écran de build-darwin-amd](/assets/blog-images/build-darwin-amd.webp)

Ah, j’allais oublier… Vous pouvez aussi utiliser `darwin/universal`… :wink:

![capture d’écran de build-darwin-universal](/assets/blog-images/build-darwin-universal.webp)

### Compilation croisée pour Windows

Comme Wails v2 pour Windows est entièrement écrit en Go, vous pouvez cibler des builds Windows sans docker.

![capture d’écran de build-cross-windows](/assets/blog-images/build-cross-windows.webp) bu

### Moteur de rendu WKWebView

V1 reposait sur un composant WebView, désormais obsolète. V2 utilise le composant WKWebKit le plus récent : attendez-vous donc à bénéficier des toutes dernières avancées d’Apple.

### Conclusion

Comme je l’ai indiqué dans les notes de version pour Windows, Wails v2 pose de nouvelles bases pour le projet. Cette version vise à recueillir des retours sur cette nouvelle approche et à corriger les éventuels bogues avant la publication d’une version finale. Vos contributions seront les bienvenues ! Veuillez transmettre vos commentaires sur le forum de discussion [v2 Beta](https://github.com/wailsapp/wails/discussions/828).

Enfin, je tiens à remercier tout particulièrement tous les [sponsors du projet](/credits/#sponsors), notamment [JetBrains](https://www.jetbrains.com?from=Wails), dont le soutien contribue de nombreuses façons au projet en coulisses.

J’ai hâte de découvrir ce que chacun créera avec Wails au cours de cette nouvelle phase passionnante du projet !

Lea.

P.-S. : utilisateurs de Linux, vous êtes les prochains !

P.-P.-S. : si vous ou votre entreprise trouvez Wails utile, envisagez de [soutenir financièrement le projet](https://github.com/sponsors/leaanthony). Merci !
