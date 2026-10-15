---
title: "Modal File Manager"
description: "Um aplicativo para desktop desenvolvido com Wails"
slug: "community/showcase/modalfilemanager"
sourcePath: "community/showcase/modalfilemanager.md"
---

![Modal File Manager](/assets/showcase-images/modalfilemanager.webp)

[Modal File Manager](https://github.com/raguay/ModalFileManager) é um gerenciador de arquivos de dois painéis que usa tecnologias web. Meu projeto original era baseado em NW.js e pode ser encontrado [aqui](https://github.com/raguay/ModalFileManager-NWjs). Esta versão usa o mesmo código de frontend baseado em Svelte (mas ele foi amplamente modificado desde a migração do NW.js), porém o backend é uma implementação em [Wails 2](https://wails.io/). Com esta implementação, não uso mais comandos de linha de comando como `rm`, `cp` etc., mas o Git precisa estar instalado no sistema para baixar temas e extensões. O aplicativo é inteiramente escrito em Go e funciona muito mais rápido do que as versões anteriores.

Este gerenciador de arquivos foi projetado com base no mesmo princípio do Vim: ações de teclado controladas por estado. O número de estados não é fixo, mas é altamente configurável. Portanto, é possível criar e usar uma quantidade ilimitada de configurações de teclado. Essa é a principal diferença em relação a outros gerenciadores de arquivos. Há temas e extensões disponíveis para download no GitHub.
