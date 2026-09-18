---
title: "Visão geral técnica"
description: "Arquitetura de alto nível e roteiro da base de código do Wails v3"
slug: "contributing/overview"
sourcePath: "contributing/overview.md"
---

## Boas-vindas à documentação técnica do Wails v3

Esta seção **não** trata das diretrizes da comunidade nem de como abrir um pull request. Em vez disso, ela explica em detalhes **como o Wails v3 é construído** para que você possa se orientar rapidamente na base de código e começar a modificá-la com confiança.

Se você pretende corrigir o runtime, ampliar a CLI, criar novos templates ou simplesmente entender os mecanismos internos, as páginas a seguir fornecem o contexto técnico necessário.

---

## Arquitetura de alto nível

@cards{cols="2"}
◇ Backend em Go
O núcleo de cada aplicativo Wails é o código Go compilado em um executável nativo. Ele é responsável pela lógica do aplicativo, pela integração com o sistema e pelas operações críticas para o desempenho.

---
▤ Frontend web
A interface do usuário é escrita com tecnologias web padrão (React, Vue, Svelte, Vanilla, …) e renderizada por uma WebView leve do sistema (WebKit no Linux/macOS e WebView2 no Windows).

---
◆ Camada de integração
Uma ponte em memória e sem cópias permite chamadas **Go⇄JavaScript** com conversão automática de tipos, propagação de eventos e encaminhamento de erros.

---
▸ CLI e ferramentas
`wails3` orquestra a criação de projetos, o servidor de desenvolvimento com recarregamento ao vivo, o empacotamento de assets, a compilação cruzada e a criação de pacotes (deb, rpm, AppImage, msi, dmg…).

@end

---

## Visão geral da arquitetura

**Wails v3 – Fluxo de ponta a ponta**

**[Espaço reservado para o diagrama do fluxo de ponta a ponta]**

O diagrama mostra o **fluxo de ponta a ponta**:

1. A **CLI** controla a geração, o servidor de desenvolvimento, a compilação e a criação de pacotes.\
2. O **sistema de bindings** produz código de integração que permite ao **frontend web** chamar o **backend em Go**.\
3. Durante o desenvolvimento, o **servidor de assets** atua como proxy para o servidor de desenvolvimento do framework; em produção, ele disponibiliza arquivos incorporados.\
4. Durante a execução, o **runtime de desktop** gerencia as janelas e as APIs do sistema operacional, enquanto a **ponte** transporta mensagens entre Go e JavaScript.

---

## O que esta documentação aborda

| Tópico | Por que é importante |
| --- | --- |
| **Estrutura da base de código** | Mapa dos diretórios de `/v3` e de como os módulos interagem. |
| **Mecanismos internos do runtime** | Gerenciamento de janelas, APIs do sistema, processador de mensagens e camadas de compatibilidade entre plataformas. |
| **Servidor de assets e de desenvolvimento** | Como os assets da web são disponibilizados no ambiente de desenvolvimento e incorporados em produção. |
| **Pipeline de build e criação de pacotes** | Fluxo de trabalho baseado em Taskfile, compilação multiplataforma e geração de instaladores. |
| **Sistema de bindings** | Pipeline de análise estática que gera bindings Go⇄TS com segurança de tipos. |
| **Sistema de templates** | Arquitetura do gerador que sustenta `wails3 init -t <framework>`. |
| **Testes e CI** | Infraestrutura de testes unitários e de integração, GitHub Actions e orientações sobre o detector de condições de corrida. |
| **Extensão do Wails** | Adição de serviços, templates ou subcomandos da CLI. |

Cada página subsequente detalha essas áreas com exemplos concretos de código, diagramas e referências aos arquivos de código-fonte pertinentes.

---

@note{type="info"}
Pré-requisitos: você deve ter familiaridade com **Go 1.25+**, TypeScript básico e ferramentas modernas de build para frontend. Se você não conhece Go, considere consultar primeiro, por alto, o tour oficial.

@end

Boa exploração — e boas-vindas aos mecanismos internos do Wails v3!
