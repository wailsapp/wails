---
title: "Zentro"
description: "IDE SQL de desktop centrada em projetos para PostgreSQL, MySQL, SQL Server e SQLite"
slug: "community/showcase/zentro"
sourcePath: "community/showcase/zentro.md"
---

![Editor SQL do Zentro com preenchimento automático baseado no esquema, o explorador de esquema e um painel de registro](/assets/showcase-images/zentro.png)

**[Zentro](https://zentro-ide.com)** é uma **IDE SQL de desktop centrada em projetos** para PostgreSQL, MySQL, SQL Server e SQLite, criada com **Go, Wails v2 e React**. Cada projeto sabe para qual ambiente aponta, e o Zentro pergunta antes de executar qualquer coisa irreversível, para que a consulta destinada ao ambiente de homologação não chegue à produção.

## Destaques

- **Consciência de ambiente** - os projetos carregam seus ambientes (desenvolvimento, homologação, produção) com uma sobreposição clara de produção, classificação de risco das instruções e uma prévia das linhas afetadas antes de uma gravação destrutiva
- **Editor SQL** - preenchimento automático baseado no esquema, diagnósticos em linha e correções rápidas
- **Explorador de esquema e ERD** - navegue por tabelas, colunas, chaves e índices e visualize os relacionamentos
- **Rastreamento com Git** - scripts salvos e instantâneos do esquema versionados como código, com histórico e diferenças
- **Túneis SSH e terminal integrado** - acesse um banco de dados por um túnel e mantenha um terminal local ou SSH ao lado das consultas
- **Assistente de IA com sua própria chave** - conecta-se ao endpoint do modelo que você configurar, nunca por um servidor operado pelo Zentro, além de um servidor MCP para agentes externos
- **Foco no teclado** - navegação por foco no explorador e na grade de resultados, paleta de comandos e atalhos em sequência
- **Privado por padrão** - as credenciais dos bancos de dados ficam no chaveiro do seu sistema operacional

O Zentro é distribuído como um único binário nativo por plataforma para Windows, macOS e Linux. O núcleo é de código aberto sob a GPL-3.0 e um nível Pro pago opcional acrescenta proteções de produção e mais, com teste gratuito de 30 dias; o código-fonte da Community Edition é público.

[Site](https://zentro-ide.com) |
[GitHub](https://github.com/ZentroIDE/Zentro) |
[Baixar](https://zentro-ide.com/download)
