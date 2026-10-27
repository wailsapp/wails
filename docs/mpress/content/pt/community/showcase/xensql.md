---
title: "XenSQL"
description: "Workbench SQL com abordagem local-first, suporte a vários bancos de dados e ferramentas avançadas para consultas"
slug: "community/showcase/xensql"
sourcePath: "community/showcase/xensql.md"
---

![Editor de tabelas do XenSQL com edição em linha](/assets/showcase-images/xensql-1.png) ![Editor de consultas do XenSQL com sugestões de tabelas e colunas](/assets/showcase-images/xensql-2.png) ![Várias consultas em uma transação no XenSQL](/assets/showcase-images/xensql-3.png)

O **[XenSQL](https://github.com/Bare7a/XenSQL)** é um **workbench SQL rápido para desktop, com abordagem local-first**, desenvolvido com **Go, Wails e React**. Ele reúne PostgreSQL, MySQL/MariaDB e SQLite em uma única interface limpa, com aparência nativa — sem nuvem, sem telemetria e sem contas.

## Principais destaques

- **Editor SQL avançado** — Baseado no Monaco, com preenchimento automático inteligente e ciente do esquema, execução de várias instruções, transmissão contínua de resultados e guias de resultados para cada instrução
- **Visualizador de dados avançado** — Inspetor JSON interativo, editor de células ciente da sintaxe (JSON, XML, HTML e texto), edição em linha e inspeção completa de registros
- **Edição de dados integrada** — Navegação por tabelas, preparação de alterações em linha, operações em massa e `INSERT`/`UPDATE`/`DELETE` seguros, com suporte a `RETURNING`
- **Recursos de produtividade** — Explorador de esquemas, consultas salvas, histórico de consultas, pesquisa rápida (`Ctrl+P`) e fluxo de trabalho orientado pelo teclado
- **Opções de exportação** — CSV, JSON, Markdown e instruções SQL INSERT

Totalmente offline e portátil. Tudo é armazenado localmente em uma única pasta `XenSQL-data/`, que acompanha o aplicativo.

**Bancos de dados compatíveis**: PostgreSQL, MySQL, MariaDB e SQLite (com modo somente leitura e opções de transporte seguro).

Projetado para desenvolvedores que buscam velocidade, clareza e controle, sem o excesso de recursos das ferramentas SQL tradicionais.
