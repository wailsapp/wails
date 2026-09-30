---
title: "Por que usar o Wails?"
description: "Entenda por que o Wails é a escolha certa para seu aplicativo de desktop"
slug: "quick-start/why-wails"
sourcePath: "quick-start/why-wails.md"
---

O Wails combina **o desempenho e a simplicidade do Go** com **a flexibilidade das interfaces web modernas**, permitindo criar aplicativos de desktop nativos e atraentes com as ferramentas que você já conhece.

## Desempenho que os usuários percebem

**Aplicativos Wails:**

- **Binários de aproximadamente 15 MB** (em comparação com 150 MB no Electron)
- **Consumo de memória inicial de aproximadamente 10 MB** (em comparação com mais de 100 MB no Electron)
- **&lt;0.5 s para iniciar** (em comparação com 2-3 s no Electron)
- **Renderização nativa** com a WebView fornecida pelo sistema operacional

Os usuários percebem seu aplicativo como rápido, leve e profissional.

## Experiência de desenvolvimento

**Escreva uma vez, execute em qualquer lugar:**

- Uma única base de código Go para Windows, macOS e Linux
- Use qualquer framework web (React, Vue, Svelte ou JavaScript puro)
- Recarga automática durante o desenvolvimento
- Bindings TypeScript gerados automaticamente a partir do código Go

Entregue mais rápido com menos código para manter.

## Recursos prontos para produção

**Tudo de que você precisa:**

- Várias janelas com ciclos de vida independentes
- Menus nativos (do aplicativo, de contexto e da bandeja do sistema)
- Caixas de diálogo de arquivos com interface nativa da plataforma
- Integração com o sistema (notificações, área de transferência e atalhos de teclado)
- Assinatura de código e empacotamento para todas as plataformas

Crie aplicativos profissionais, não protótipos.

## Desenvolvimento mais rápido

- **Uma base de código, três plataformas** — Escreva uma vez e compile para Windows, macOS e Linux
- **Use os conhecimentos que você já tem** — Go no backend e HTML/CSS/JS na interface
- **Feedback instantâneo** — Recarga automática durante o desenvolvimento e tempos de compilação medidos em segundos
- **Binários pequenos** — Aplicativos de 15 MB significam compilações, downloads e iterações mais rápidos

## Quando escolher o Wails

**O Wails é perfeito para:**

- **Aplicativos empresariais** (CRM, estoque, painéis e ferramentas administrativas)
- **Ferramentas para desenvolvedores** (clientes de banco de dados, ferramentas de teste de APIs e de implantação)
- **Aplicativos de produtividade** (anotações, gerenciadores de tarefas e controle de tempo)
- **Ferramentas criativas** (editores de imagens, processadores de vídeo e utilitários de design)
- **Ferramentas internas** (aplicativos específicos da empresa e ferramentas de automação)

## Casos reais de sucesso

@note{type="tip" title="Aplicativos em produção"}
O Wails viabiliza aplicativos reais usados por milhares de usuários:

- **Ferramentas de gerenciamento de bancos de dados** com interfaces complexas
- **Painéis financeiros** que processam dados em tempo real
- **Ferramentas de edição de vídeo** com desempenho nativo
- **Utilitários de desenvolvimento** usados por equipes de engenharia

[Veja os destaques →](/community/showcase/)

@end

## Como o Wails funciona

Ao contrário do Electron, que inclui um navegador completo e o runtime do Node.js, o Wails adota uma abordagem fundamentalmente diferente: seu código Go é compilado em um binário nativo, e sua interface é executada na WebView integrada ao sistema operacional. Essa arquitetura proporciona binários pequenos, inicialização rápida e baixo consumo de memória, fazendo com que os aplicativos Wails tenham uma experiência nativa.

### Arquitetura

Os aplicativos Wails são compostos por duas partes principais que se comunicam de forma transparente: um backend Go, responsável pela lógica de negócios e pelas operações do sistema, e um frontend baseado em tecnologias web para a interface do usuário. A WebView fornecida pelo sistema operacional renderiza a interface sem incluir um navegador no pacote, enquanto a camada de bindings oferece comunicação com segurança de tipos entre Go e JavaScript.

<div style="display: flex; justify-content: center; align-items: center; margin: 2rem 0;">
  <img src="/img/architecture.svg" alt="Arquitetura do Wails: um backend em Go e sua interface web compilados em um único binário nativo, conectados por bindings gerados e renderizados pelo WebView do sistema operacional" style="max-width: 640px; width: 100%;" />
</div>

Essa arquitetura simples permite que o código JavaScript chame funções Go diretamente (por meio de bindings gerados automaticamente), enquanto o Go pode enviar eventos e dados de volta ao frontend. As duas camadas se comunicam por uma ponte eficiente em memória, com sobrecarga inferior a um milissegundo.

**Como o Wails alcança esse desempenho:**

1. **Nenhum runtime incluído** — Usa o binário compilado do Go
2. **WebView nativa** — Mecanismo de renderização fornecido pelo sistema operacional
3. **Ponte direta entre Go ↔ JS** — Comunicação em memória, sem sobrecarga de rede
4. **Binário compilado** — Inicialização instantânea, sem compilação JIT

## Próximos passos

Agora que você entende o que o Wails oferece, vamos configurar seu ambiente:

1. **Instale o Wails** — Configure seu ambiente de desenvolvimento em 5 minutos [Guia de instalação →](/quick-start/installation/)

2. **Crie seu primeiro aplicativo** — Crie um aplicativo funcional e entenda os conceitos básicos [Tutorial do primeiro aplicativo →](/quick-start/first-app/)

3. **Explore os recursos** — Descubra o que o Wails pode fazer pelo seu aplicativo [Visão geral dos recursos →](/quick-start/next-steps/)

---

**Ainda tem dúvidas?** Participe da nossa [comunidade no Discord](https://discord.gg/JDdSxwjhGf) e fale diretamente com a equipe.
