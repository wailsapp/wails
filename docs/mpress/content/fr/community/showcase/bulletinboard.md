---
title: "BulletinBoard"
description: "Une application de bureau créée avec Wails"
slug: "community/showcase/bulletinboard"
sourcePath: "community/showcase/bulletinboard.md"
---

![BulletinBoard](/assets/showcase-images/bboard.webp)

L’application [BulletinBoard](https://github.com/raguay/BulletinBoard) est un tableau d’affichage polyvalent qui permet de présenter des messages statiques ou des boîtes de dialogue afin de recueillir auprès de l’utilisateur des informations destinées à un script. Elle dispose d’une interface utilisateur en mode texte (TUI) pour créer de nouvelles boîtes de dialogue qui pourront ensuite servir à recueillir des informations auprès de l’utilisateur. Elle est conçue pour rester active sur votre système, afficher les informations lorsque nécessaire, puis se masquer. Sur mon système, j’utilise un processus qui surveille un fichier et envoie son contenu à BulletinBoard lorsqu’il est modifié. Il s’intègre parfaitement à mes flux de travail. Un [workflow Alfred](https://github.com/raguay/MyAlfred/blob/master/Alfred%205/EmailIt.alfredworkflow) permet également d’envoyer des informations au programme. Ce workflow fonctionne aussi avec [EmailIt](https://github.com/raguay/EmailIt).
