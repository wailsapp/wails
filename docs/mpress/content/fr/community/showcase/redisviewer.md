---
title: "Redis Viewer"
description: "Une interface graphique de bureau pour Redis conçue avec Wails"
slug: "community/showcase/redisviewer"
sourcePath: "community/showcase/redisviewer.md"
---

![Capture d’écran de RedisViewer](/assets/showcase-images/redisviewer-overview1.webp)

![Capture d’écran de RedisViewer](/assets/showcase-images/redisviewer-overview2.webp)

[RedisViewer](https://redisviewer.com/) est une interface graphique de bureau moderne pour Redis, conçue avec Wails. Utilisez-la pour inspecter des valeurs complexes, exécuter des commandes et analyser les performances de Redis sans sacrifier la qualité des interactions.

Conçue autour de l’architecture WebView + Go de Wails, l’application conserve les grands espaces de clés et les charges utiles volumineuses dans le backend Go au lieu de tout transférer vers le frontend. Elle tire ainsi parti du modèle mémoire de Go tout en évitant la pression sur le tas et les risques de fuite courants dans les clients qui reposent fortement sur JS. L’interface utilisateur est optimisée pour n’afficher que ce qui apparaît à l’écran, afin que la navigation dans d’immenses jeux de données reste réactive et offre une expérience proche de celle d’une application de bureau native.

[Visiter le site web du projet](https://redisviewer.com/)
