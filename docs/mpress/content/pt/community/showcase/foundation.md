---
title: "Foundation"
description: "Scaffold de app desktop para Wails v3.1 com React 19 e Mantine 9"
slug: "community/showcase/foundation"
sourcePath: "community/showcase/foundation.md"
---

![Configurações de layout do Foundation com nove predefinições alternáveis](/assets/showcase-images/foundation-layouts.webp)

![Página inicial do Foundation no tema escuro](/assets/showcase-images/foundation-home-dark.webp)

![Gráficos X-Pro do Foundation no tema Obsidian](/assets/showcase-images/foundation-charts.webp)

**[Foundation](https://github.com/FxRayHughes/foundation)** é um **scaffold de aplicativo desktop para Wails v3.1**, construído com **Go, React 19, TypeScript e Mantine 9**. Ele oferece um backend Go em camadas, janelas adaptadas a cada plataforma, persistência com SQLite, temas, internacionalização e testes de ponta a ponta como ponto de partida para novos projetos.

## Recursos

- **Janelas por plataforma** — janela sem bordas com barra de título desenhada pelo frontend no Windows e no Linux; no macOS, barra de título imersiva que mantém os botões nativos da janela.
- **Nove predefinições de layout** — barra de ícones, barra lateral, navegação superior, barra de título mesclada, barra à direita, dock inferior, dock flutuante arrastável, cartão recuado e modo foco, sobre uma estrutura CSS Grid; é possível registrar predefinições personalizadas.
- **Temas e internacionalização** — predefinições clara, escura e Obsidian e um editor de paleta personalizada; interface em inglês e chinês; tema e idioma ficam sincronizados entre a janela principal, as janelas filhas e o painel da bandeja.
- **Persistência com SQLite** — GORM AutoMigrate sem arquivos de migração, local do banco de dados escolhido pelo usuário e gráficos de uso por tabela.
- **Integração com o desktop** — painel na bandeja do sistema, janelas filhas, diálogos nativos de arquivos e mensagens, e exportação CSV / XLSX pelo diálogo nativo "Salvar como".
- **Utilitários** — cliente HTTP, subprocessos vinculados a um Job Object do Windows, criptografia AES-GCM, logs com rotação e gravações atômicas.
- **Testes de ponta a ponta** — o Playwright é executado contra um build do Wails em modo servidor (backend Go e frontend reais); a mesma suíte gera as capturas de tela do README.

O frontend chama os serviços Go pelos bindings gerados pelo Wails, e cada plataforma gera um único binário. O frontend segue uma estrutura MVVM, e o README traz um guia passo a passo para renomear o scaffold para o seu próprio projeto.

Funciona no Windows, macOS e Linux.

[GitHub](https://github.com/FxRayHughes/foundation) |
[Versões](https://github.com/FxRayHughes/foundation/releases/latest)
