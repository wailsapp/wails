---
title: "Snippet Expander"
description: "Um aplicativo para desktop criado com Wails"
slug: "community/showcase/snippetexpander"
sourcePath: "community/showcase/snippetexpander.md"
---

![Captura de tela do Snippet Expander](/assets/showcase-images/snippetexpandergui-select-snippet.png)

Captura de tela da janela Select Snippet do Snippet Expander

![Captura de tela do Snippet Expander](/assets/showcase-images/snippetexpandergui-add-snippet.png)

Captura de tela da tela Add Snippet do Snippet Expander

![Captura de tela do Snippet Expander](/assets/showcase-images/snippetexpandergui-search-and-paste.png)

Captura de tela da janela Search & Paste do Snippet Expander

[Snippet Expander](https://snippetexpander.org) é o «seu pequeno assistente de expansão de trechos de texto» para Linux.

O Snippet Expander consiste em um aplicativo com interface gráfica criado com Wails para gerenciar trechos e configurações, além de um modo de janela Search & Paste para selecionar e colar rapidamente um trecho.

A interface gráfica baseada em Wails, a CLI em go-lang e o daemon de expansão automática em vala-lang se comunicam com um daemon em go-lang via D-Bus. O daemon realiza a maior parte do trabalho: gerencia o banco de dados de trechos e as configurações comuns e fornece serviços para expandir e colar trechos, entre outras tarefas.

Consulte o [código-fonte](https://git.sr.ht/~ianmjones/snippetexpander/tree/trunk/item/cmd/snippetexpandergui/app.go#L38) para ver como o aplicativo Wails envia mensagens da interface do usuário ao backend, que depois são encaminhadas ao daemon, e como assina um evento D-Bus para monitorar alterações feitas nos trechos por outra instância do aplicativo ou pela CLI e exibi-las instantaneamente na interface por meio de um evento do Wails.
