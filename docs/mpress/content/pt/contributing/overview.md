---
title: "Visão geral técnica"
description: "Arquitetura de alto nível e roteiro do código-base do Wails v3"
slug: "contributing/overview"
sourcePath: "contributing/overview.md"
---

## Boas-vindas à documentação técnica do Wails v3

Esta seção **não** trata de diretrizes da comunidade nem de como abrir uma pull request. Em vez disso, ela se aprofunda em **como o Wails v3 é desenvolvido**, para que você possa se orientar rapidamente no código-base e começar a modificá-lo com confiança.

Seja para corrigir o runtime, estender a CLI, criar novos templates ou apenas entender os componentes internos, as páginas a seguir fornecem o contexto técnico necessário.

---

## Arquitetura de alto nível

@cards{cols="2"}
◇ Backend em Go
O núcleo de todo aplicativo Wails é um código Go compilado em um executável nativo. Ele é responsável pela lógica do aplicativo, pela integração com o sistema e pelas operações críticas para o desempenho.

---
▤ Frontend web
A interface é escrita com tecnologias web padrão (React, Vue, Svelte, JavaScript puro, …) e renderizada por uma WebView leve do sistema (WebKit no Linux/macOS e WebView2 no Windows).

---
◆ Camada de integração
Uma ponte em memória e sem cópias permite chamadas **Go⇄JavaScript** com conversão automática de tipos, propagação de eventos e encaminhamento de erros.

---
▸ CLI e ferramentas
`wails3` orquestra a criação de projetos, o servidor de desenvolvimento com recarregamento em tempo real, o empacotamento de recursos, a compilação cruzada e a geração de pacotes (deb, rpm, AppImage, msi, dmg…).

@end

---

## Visão geral da arquitetura

**Wails v3 – Fluxo de ponta a ponta**

**[Espaço reservado para o diagrama do fluxo de ponta a ponta]**

O diagrama mostra o **fluxo de ponta a ponta**:

1. A **CLI** controla a geração, o servidor de desenvolvimento, a compilação e a criação de pacotes.\
2. O **sistema de bindings** gera o código de integração que permite ao **frontend web** chamar o **backend em Go**.\
3. Durante o desenvolvimento, o **servidor de recursos** atua como proxy para o servidor de desenvolvimento do framework; em produção, ele disponibiliza os arquivos incorporados.\
4. Em tempo de execução, o **runtime de desktop** gerencia as janelas e as APIs do sistema operacional, enquanto a **ponte** transporta mensagens entre Go e JavaScript.

---

## O que esta documentação aborda

| Tópico | Por que é importante |
| --- | --- |
| **Estrutura do código-base** | Mapa dos diretórios de `/v3` e de como os módulos interagem. |
| **Componentes internos do runtime** | Gerenciamento de janelas, APIs do sistema, processador de mensagens e camadas de compatibilidade entre plataformas. |
| **Servidor de recursos e de desenvolvimento** | Como os recursos web são disponibilizados durante o desenvolvimento e incorporados em produção. |
| **Pipeline de compilação e empacotamento** | Fluxo de trabalho baseado em Taskfile, compilação multiplataforma e geração de instaladores. |
| **Sistema de bindings** | Pipeline de análise estática que gera bindings Go⇄TS com tipagem segura. |
| **Sistema de templates** | Arquitetura do gerador que viabiliza `wails3 init -t <framework>`. |
| **Testes e CI** | Infraestrutura de testes unitários e de integração, GitHub Actions e orientações sobre o detector de condições de corrida. |
| **Como estender o Wails** | Adição de serviços, templates ou subcomandos da CLI. |

Cada página subsequente se aprofunda nessas áreas com exemplos concretos de código, diagramas e referências aos arquivos-fonte relevantes.

---

@note{type="info"}
Pré-requisitos: você deve ter familiaridade com **Go 1.26+**, noções básicas de TypeScript e ferramentas modernas de compilação de frontend. Se você ainda não conhece Go, considere consultar primeiro o tour oficial.

@end

Boa exploração — e boas-vindas aos componentes internos do Wails v3!
