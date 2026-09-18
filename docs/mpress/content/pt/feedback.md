---
title: "Feedback"
description: "Como enviar feedback e relatar problemas do Wails v3"
slug: "feedback"
sourcePath: "feedback.md"
---

Agradecemos (e incentivamos) seu feedback! Pesquise problemas ou discussões existentes antes de criar novos. Estas são as diferentes formas de contribuir:

@tabs
[Bugs]
Se você encontrar um bug, [abra uma issue](https://github.com/wailsapp/wails/issues/new/choose) no GitHub usando o modelo de relatório de bug.

- Descreva claramente o bug com um exemplo simples e reproduzível. Se a documentação não deixar claro o que *deveria* acontecer, inclua essa informação no relatório.
- Inclua a saída de `wails3 doctor` no relatório.
- Se o bug for um comportamento que não está de acordo com a documentação atual, faça também o seguinte:
  - Atualize um exemplo existente no diretório `v3/examples` ou crie um novo exemplo que demonstre claramente o problema.
  - Abra um [PR](https://github.com/wailsapp/wails/pulls) que faça referência à issue.


@note{type="caution"}
*Lembre-se*: um comportamento inesperado não é necessariamente um bug — talvez ele apenas não faça o que você espera. Para esses casos, use `Suggestions`.

@end

Você também pode discutir bugs no canal [#v3](https://discord.gg/bdj28QNHmT) do Discord.

[Correções]
Se você tiver uma correção de bug ou uma melhoria para a documentação:

- Abra um pull request no [repositório do Wails](https://github.com/wailsapp/wails) seguindo as [diretrizes de contribuição](https://github.com/wailsapp/wails/blob/master/CONTRIBUTING.md).
- Faça referência a todas as issues relacionadas na descrição do PR.

[Aprimoramentos]
Novas funcionalidades e alterações no comportamento público devem ser propostas por meio de um pull request de rascunho de uma **WEP (Proposta de Aprimoramento do Wails)**, e não por uma issue de solicitação de recurso.

- Leia o [processo de WEP](https://github.com/wailsapp/wails/blob/master/v3/wep/README.md).
- Copie o modelo e abra um PR de rascunho com o título `[WEP] <title>`, contendo apenas a WEP e o material de apoio.
- Você pode primeiro discutir informalmente uma ideia nas [Discussões do GitHub](https://github.com/wailsapp/wails/discussions) ou no canal [#v3](https://discord.gg/bdj28QNHmT) do Discord, mas um PR de WEP é obrigatório para que um mantenedor tome uma decisão.

[Votos positivos]
- Demonstre apoio a bugs, WEPs e discussões usando a reação :thumbsup: no GitHub.
- *Não* adicione apenas comentários como "+1" ou "eu também".
- Adicione um comentário se tiver algo relevante com que contribuir, como "esse bug também afeta compilações para ARM" ou "Outra abordagem seria...".

@end

Problemas conhecidos e trabalhos em andamento podem ser encontrados [aqui](https://github.com/orgs/wailsapp/projects/6).
