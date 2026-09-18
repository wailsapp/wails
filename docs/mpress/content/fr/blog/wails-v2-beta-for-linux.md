---
title: "Wails v2 bêta pour Linux"
description: "Notes de version et annonces concernant Wails"
authors: ["leaanthony"]
tags: ["wails","v2"]
date: "2022-02-22"
slug: "blog/wails-v2-beta-for-linux"
image: "/assets/blog-images/wails-linux.webp"
sourcePath: "blog/wails-v2-beta-for-linux.md"
---

![capture d’écran de wails-linux](/assets/blog-images/wails-linux.webp)

J’ai le plaisir de vous annoncer enfin que Wails v2 est désormais disponible en version bêta pour Linux ! Il est assez ironique que les toutes premières expérimentations avec la v2 aient eu lieu sous Linux, alors que cette plateforme est finalement la dernière à bénéficier d’une version. Cela dit, la v2 actuelle est très différente de ces premières expérimentations. Sans plus attendre, passons en revue les nouvelles fonctionnalités :

## Nouvelles fonctionnalités

![capture d’écran de wails-menus-linux](/assets/blog-images/wails-menus-linux.webp)

La prise en charge des menus natifs a fait l’objet de nombreuses demandes. Wails y répond enfin. Les menus d’application sont désormais disponibles et prennent en charge la plupart des fonctionnalités des menus natifs, notamment les éléments de menu standard, les cases à cocher, les groupes de boutons radio, les sous-menus et les séparateurs.

Dans la v1, vous avez été très nombreux à demander davantage de contrôle sur la fenêtre elle-même. J’ai le plaisir de vous annoncer que de nouvelles API d’exécution sont spécialement prévues à cet effet. Riches en fonctionnalités, elles prennent en charge les configurations à plusieurs écrans. L’API des boîtes de dialogue a également été améliorée : vous pouvez désormais utiliser des boîtes de dialogue modernes et natives, assorties de nombreuses options de configuration pour répondre à tous vos besoins.

### Aucune obligation de regrouper les ressources

L’un des principaux points faibles de la v1 était la nécessité de condenser toute votre application dans un seul fichier JS et un seul fichier CSS. J’ai le plaisir de vous annoncer qu’avec la v2, il n’est absolument plus nécessaire de regrouper les ressources. Vous souhaitez charger une image locale ? Utilisez une balise `<../../../assets/blog-images>` avec un chemin src local. Vous souhaitez utiliser une police sympa ? Copiez-la et ajoutez son chemin dans votre CSS.

> Cela ressemble beaucoup à un serveur web…

Oui, cela fonctionne exactement comme un serveur web, sauf que ce n’en est pas un.

> Alors, comment inclure mes ressources ?

Il vous suffit de transmettre à la configuration de votre application un seul `embed.FS` contenant toutes vos ressources. Elles n’ont même pas besoin de se trouver dans le répertoire racine : Wails s’occupe de tout pour vous.

### Nouvelle expérience de développement

Comme il n’est désormais plus nécessaire de regrouper les ressources, une toute nouvelle expérience de développement devient possible. La nouvelle commande `wails dev` compile et exécute votre application, mais au lieu d’utiliser les ressources du `embed.FS`, elle les charge directement depuis le disque.

Elle offre également les fonctionnalités supplémentaires suivantes :

- Rechargement à chaud : toute modification des ressources du frontend déclenche le rechargement automatique du frontend de l’application
- Recompilation automatique : toute modification de votre code Go entraîne la recompilation et le redémarrage de votre application

En outre, un serveur web démarre sur le port 34115. Il sert votre application à tout navigateur qui s’y connecte. Tous les navigateurs web connectés réagissent aux événements système, par exemple en effectuant un rechargement à chaud lorsque les ressources sont modifiées.

En Go, nous avons l’habitude de manipuler des structures dans nos applications. Il est souvent utile de transmettre des structures au frontend et de les utiliser comme état de l’application. Dans la v1, ce processus était très manuel et assez contraignant pour les développeurs. J’ai le plaisir de vous annoncer que, dans la v2, toute application exécutée en mode développement génère automatiquement des modèles TypeScript pour toutes les structures utilisées comme paramètres d’entrée ou de sortie des méthodes liées. Cela permet d’échanger les modèles de données de manière fluide entre ces deux univers.

En outre, un autre module JS est généré dynamiquement afin d’encapsuler toutes vos méthodes liées. Il fournit la documentation JSDoc de vos méthodes, ce qui permet à votre IDE de proposer la saisie semi-automatique et des indications. C’est vraiment très pratique de voir les modèles de données importés automatiquement lorsque vous appuyez sur la touche de tabulation dans un module généré automatiquement qui encapsule votre code Go !

### Modèles distants

![capture d’écran de remote-linux](/assets/blog-images/remote-linux.webp)

Pouvoir mettre rapidement une application en service a toujours été un objectif essentiel du projet Wails. Lors du lancement, nous avons essayé de prendre en charge de nombreux frameworks modernes de l’époque : react, vue et angular. Le monde du développement frontend est très marqué par les préférences de chacun, évolue rapidement et reste difficile à suivre ! Nos modèles de base sont donc devenus assez vite obsolètes, ce qui a compliqué leur maintenance. Cela signifiait également que nous ne disposions pas de modèles modernes et attrayants pour les toutes dernières piles technologiques.

Avec la v2, je souhaitais donner davantage de moyens à la communauté en vous permettant de créer et d’héberger vous-mêmes des modèles, au lieu de dépendre du projet Wails. Vous pouvez donc désormais créer des projets à partir de modèles maintenus par la communauté ! J’espère que cela incitera les développeurs à bâtir un écosystème dynamique de modèles de projet. J’ai vraiment hâte de découvrir ce que notre communauté de développeurs va créer !

### Compilation croisée pour Windows

Comme Wails v2 pour Windows est entièrement écrit en Go, vous pouvez produire des builds ciblant Windows sans docker.

![capture d’écran de build-cross-windows](/assets/blog-images/linux-build-cross-windows.webp)

### Conclusion

Comme je l’ai indiqué dans les notes de version pour Windows, Wails v2 constitue une nouvelle base pour le projet. Cette version vise à recueillir vos retours sur la nouvelle approche et à corriger les éventuels bogues avant la publication de la version finale. Vos avis seront les bienvenus ! Merci de transmettre tous vos commentaires sur le forum de discussion [v2 bêta](https://github.com/wailsapp/wails/discussions/828).

La prise en charge de Linux est **difficile**. Nous nous attendons à ce que la version bêta présente quelques particularités. Aidez-nous à vous aider en envoyant des rapports de bogue détaillés !

Enfin, je tiens à remercier tout particulièrement tous les [sponsors du projet](/credits/#sponsors), dont le soutien fait avancer le projet de bien des façons en coulisses.

J’ai hâte de découvrir ce que chacun créera avec Wails au cours de cette nouvelle phase passionnante du projet !

Lea.

P.-S. : la sortie de la v2 est désormais toute proche !

P.-P.-S. : si vous ou votre entreprise trouvez Wails utile, envisagez de [soutenir financièrement le projet](https://github.com/sponsors/leaanthony). Merci !
