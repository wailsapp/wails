---
title: "Próximos passos"
description: "O que fazer depois de criar seu primeiro aplicativo"
slug: "quick-start/next-steps"
sourcePath: "quick-start/next-steps.md"
---

Você criou seu primeiro aplicativo Wails. Siga o percurso de aprendizagem abaixo ou vá direto para um guia se já souber o que deseja criar.

## Saiba mais

Continue aprendendo com um aplicativo completo de cada vez:

1. **[Crie um serviço de QR Code](/tutorials/01-creating-a-service/).** Pratique como expor a lógica em Go ao frontend por meio de um serviço.
2. **[Crie uma lista de tarefas](/tutorials/02-todo-vanilla/).** Adicione, altere e exclua dados do aplicativo por meio de chamadas de serviço.
3. **[Crie um aplicativo de anotações](/tutorials/03-notes-vanilla/).** Trabalhe com arquivos, caixas de diálogo nativas e dados salvos.

Quando estiver pronto para explorar as versões, siga o [tutorial sobre aplicativos com atualização automática](/tutorials/04-self-update-a-wails-app/). Consulte [todos os tutoriais](/tutorials/overview/) para ver os objetivos de aprendizagem.

## Recursos

**Já está criando seu próprio aplicativo?** Use um guia prático para a tarefa que precisa realizar:

| Quero… | Comece aqui |
| --- | --- |
| Chamar Go pelo meu frontend | [Chamar Go pelo JavaScript](/features/bindings/methods/) |
| Abrir ou gerenciar uma janela | [Noções básicas sobre janelas](/features/windows/basics/) |
| Adicionar menus nativos ou um ícone na bandeja do sistema | [Menus do aplicativo](/features/menus/application/) ou [menus da bandeja do sistema](/features/menus/systray/) |
| Pedir ao usuário que abra ou salve um arquivo | [Caixas de diálogo de arquivos](/features/dialogs/file/) |
| Enviar eventos entre componentes | [Sistema de eventos](/features/events/system/) |
| Copiar texto ou aceitar arquivos arrastados | [Área de transferência](/features/clipboard/basics/) ou [arrastar e soltar arquivos](/features/drag-and-drop/files/) |
| Adicionar atalhos de teclado | [Atalhos de teclado](/features/keyboard/shortcuts/) |
| Compilar e empacotar meu aplicativo | [Compilar aplicativos](/guides/build/building/) e depois empacotá-los para [Windows](/guides/build/windows/), [macOS](/guides/build/macos/) ou [Linux](/guides/build/linux/) |

Essas páginas ficam em **Guias práticos** na barra lateral. Você pode abrir o tópico relevante sem precisar ler toda a seção na ordem.

## Entenda como o Wails funciona

Consulte **Conceitos fundamentais** quando quiser entender o design por trás do código:

- [Como o Wails funciona](/concepts/architecture/) mostra a conexão entre o backend em Go, o frontend e a WebView nativa.
- [Ponte entre Go e o frontend](/concepts/bridge/) explica como as chamadas transitam entre Go e JavaScript.
- [Ciclo de vida do aplicativo](/concepts/lifecycle/) explica a inicialização e o encerramento.
- [Sistema de compilação](/concepts/build-system/) explica a função das tarefas do projeto e das ferramentas de frontend.

## Consulte uma API

Consulte a [Referência da API](/reference/overview/) quando precisar de um tipo, método, opção ou comando, em vez de uma lição. Para trabalhar na linha de comando, comece pela [visão geral dos comandos da CLI](/reference/cli/); consulte [opções e exemplos da CLI](/guides/cli/) para obter detalhes de cada comando.

## Obtenha ajuda

Consulte as [perguntas frequentes](/faq/), explore os [exemplos funcionais](https://github.com/wailsapp/wails/tree/releases/v3-beta/v3/examples) ou faça uma pergunta no [Discord](https://discord.gg/JDdSxwjhGf). Os problemas específicos de cada plataforma estão agrupados em **Solução de problemas** na barra lateral.
