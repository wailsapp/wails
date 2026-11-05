---
title: "Modal File Manager"
description: "Une application de bureau développée avec Wails"
slug: "community/showcase/modalfilemanager"
sourcePath: "community/showcase/modalfilemanager.md"
---

![Modal File Manager](/assets/showcase-images/modalfilemanager.webp)

[Modal File Manager](https://github.com/raguay/ModalFileManager) est un gestionnaire de fichiers à deux volets utilisant les technologies web. Ma conception d’origine reposait sur NW.js et se trouve [ici](https://github.com/raguay/ModalFileManager-NWjs). Cette version utilise le même code d’interface basé sur Svelte, mais celui-ci a été considérablement modifié depuis l’abandon de NW.js, tandis que le backend est une implémentation de [Wails 2](https://wails.io/). Grâce à cette implémentation, je n’utilise plus les commandes en ligne de commande `rm`, `cp`, etc. Toutefois, Git doit être installé sur le système pour télécharger des thèmes et des extensions. Cette version est entièrement écrite en Go et s’exécute beaucoup plus rapidement que les versions précédentes.

Ce gestionnaire de fichiers repose sur le même principe que Vim : des actions au clavier contrôlées par des états. Le nombre d’états n’est pas fixe, mais se configure très librement. Il est donc possible de créer et d’utiliser un nombre illimité de configurations de clavier. C’est ce qui le distingue principalement des autres gestionnaires de fichiers. Des thèmes et des extensions peuvent être téléchargés depuis GitHub.
