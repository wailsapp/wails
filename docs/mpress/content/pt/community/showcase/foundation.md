---
title: "Foundation"
description: "Um scaffold de app desktop Wails v3 pronto para produção"
slug: "community/showcase/foundation"
sourcePath: "community/showcase/foundation.md"
---

![Configurações de layout do Foundation com nove predefinições alternáveis](/assets/showcase-images/foundation-layouts.webp)

![Página inicial do Foundation no tema escuro](/assets/showcase-images/foundation-home-dark.webp)

![Gráficos X-Pro do Foundation no tema Obsidian](/assets/showcase-images/foundation-charts.webp)

**[Foundation](https://github.com/FxRayHughes/foundation)** é um **scaffold de aplicativo desktop para Wails v3** pronto para produção, construído com **Go, React 19, TypeScript e Mantine 9**. Ele entrega a um novo projeto as partes que normalmente levam semanas para construir — um backend Go em camadas, janelas adaptadas a cada plataforma, persistência, temas, internacionalização e testes de ponta a ponta — para que você vá direto à lógica de negócio.

## Principais recursos

- **Janela nativa em todas as plataformas** — barra de título sem bordas desenhada pelo frontend no Windows e no Linux; barra de título imersiva no macOS que mantém os botões nativos da janela.
- **Nove layouts alternáveis** — barra de ícones, barra lateral, navegação superior, barra de título mesclada, barra à direita, dock inferior, dock flutuante arrastável, cartão recuado e modo foco, tudo sobre uma estrutura CSS Grid; registre sua própria predefinição em uma linha.
- **Temas e internacionalização** — predefinições clara, escura e Obsidian, além de um editor de paleta personalizada; inglês e chinês prontos para uso; tema e idioma ficam sincronizados entre a janela principal, as janelas filhas e o painel da bandeja.
- **Persistência com SQLite** — GORM AutoMigrate sem arquivos de migração, local do banco de dados escolhido pelo usuário e gráficos de uso por tabela.
- **Integração com o desktop** — painel na bandeja do sistema, janelas filhas, diálogos nativos de arquivos e mensagens, e exportação CSV / XLSX pelo diálogo nativo "Salvar como".
- **Utilitários robustos** — cliente HTTP, subprocessos vinculados a um Job Object do Windows, criptografia AES-GCM, logs com rotação e gravações atômicas.
- **Testado de ponta a ponta** — o Playwright controla um build do Wails em modo servidor (backend Go real e frontend real), e a mesma suíte gera as capturas de tela do README.

Com o Wails, tudo fica em um único binário por plataforma: os serviços Go são chamados diretamente pelos bindings gerados, enquanto o frontend segue uma estrutura MVVM rigorosa que continua fácil de manter à medida que novos recursos são adicionados.

Disponível para Windows, macOS e Linux.

[GitHub](https://github.com/FxRayHughes/foundation) |
[Versões](https://github.com/FxRayHughes/foundation/releases/latest)
