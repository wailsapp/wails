---
title: "Redis Viewer"
description: "Uma interface gráfica de desktop para Redis criada com Wails"
slug: "community/showcase/redisviewer"
sourcePath: "community/showcase/redisviewer.md"
---

![Captura de tela do RedisViewer](/assets/showcase-images/redisviewer-overview1.webp)

![Captura de tela do RedisViewer](/assets/showcase-images/redisviewer-overview2.webp)

[RedisViewer](https://redisviewer.com/) é uma moderna interface gráfica de desktop para Redis criada com Wails. Use-a para inspecionar valores complexos, executar comandos e analisar o desempenho do Redis sem comprometer a qualidade da interação.

Projetado com base na arquitetura WebView + Go do Wails, ele mantém grandes espaços de chaves e cargas úteis pesadas no backend em Go, em vez de enviar tudo para o frontend, aproveitando o modelo de memória do Go e evitando a pressão sobre o heap e os riscos de vazamento comuns em clientes que dependem muito de JS. A interface é otimizada para renderizar apenas o que está na tela; assim, a navegação por conjuntos de dados enormes continua responsiva e proporciona uma experiência próxima à de um aplicativo de desktop nativo.

[Visitar o site do projeto](https://redisviewer.com/)
