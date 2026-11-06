---
title: "Feedback"
description: "Como enviar feedback e relatar problemas do Wails v3"
slug: "feedback"
sourcePath: "feedback.md"
---

Agradecemos (e incentivamos) seu feedback! Antes de criar um novo issue ou iniciar uma nova discussão, pesquise os já existentes. Estas são as diferentes formas de contribuir:

@tabs
[Bugs]
Se encontrar um bug, [abra um issue](https://github.com/wailsapp/wails/issues/new/choose) no GitHub usando o modelo de relatório de bug.

- Descreva claramente o bug com um exemplo simples e reproduzível. Se a documentação não deixar claro o que *deveria* acontecer, mencione isso no relatório.
- Inclua a saída de `wails3 doctor` no relatório.
- Se o bug for um comportamento que não esteja de acordo com a documentação atual, faça também o seguinte:
  - Atualize um exemplo existente no diretório `v3/examples` ou crie um novo exemplo que demonstre claramente o problema.
  - Abra um [PR](https://github.com/wailsapp/wails/pulls) que faça referência ao issue.


@note{type="caution"}
*Lembre-se*: um comportamento inesperado não é necessariamente um bug — ele pode apenas não corresponder ao que você esperava. Nesse caso, use `Suggestions`.

@end

Você também pode discutir bugs no canal [#v3](https://discord.gg/bdj28QNHmT) do Discord.

[Correções]
Se você tiver uma correção para um bug ou uma melhoria na documentação:

- Abra um pull request no [repositório do Wails](https://github.com/wailsapp/wails) seguindo as [diretrizes de contribuição](https://github.com/wailsapp/wails/blob/master/CONTRIBUTING.md).
- Faça referência a todos os issues relacionados na descrição do PR.

[Melhorias]
Novas funcionalidades e alterações no comportamento público devem ser propostas por meio de um pull request de rascunho de uma **WEP (Wails Enhancement Proposal)**, e não por um issue de solicitação de funcionalidade.

- Leia sobre o [processo de WEP](https://github.com/wailsapp/wails/blob/v3.0.0/v3/wep/README.md).
- Copie o modelo e abra um PR de rascunho intitulado `[WEP] <title>` que contenha apenas a WEP e os materiais de apoio.
- Você pode primeiro discutir uma ideia informalmente nas [Discussões do GitHub](https://github.com/wailsapp/wails/discussions) ou no canal [#v3](https://discord.gg/bdj28QNHmT) do Discord, mas um PR de WEP é obrigatório para que um mantenedor tome uma decisão.

[Votação positiva]
- Demonstre apoio a bugs, WEPs e discussões usando a reação :thumbsup: no GitHub.
- *Não* adicione apenas comentários como "+1" ou "eu também".
- Adicione um comentário se tiver algo relevante a contribuir, como "este bug também afeta compilações para ARM" ou "Outra abordagem seria...".

@end

Os problemas conhecidos e os trabalhos em andamento estão disponíveis [aqui](https://github.com/orgs/wailsapp/projects/6).
