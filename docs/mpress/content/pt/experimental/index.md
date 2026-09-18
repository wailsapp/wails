---
title: "Recursos experimentais"
description: "Um espaço para experimentos em andamento no Wails v3 — o que são, por que existem e como enviar feedback."
slug: "experimental"
sourcePath: "experimental/index.md"
---

@note{type="caution" title="Aqui há experimentos"}
Tudo nesta seção é, por definição, um experimento. Os recursos apresentados aqui são opcionais, vêm desativados por padrão e podem mudar de formato, ser renomeados ou removidos por completo entre versões. Não baseie nada essencial neles sem estar preparado para mudanças frequentes.

@end

## O que significa "Experimental"

O Wails disponibiliza experimentos abertamente. Um experimento é uma ideia que consideramos promissora o bastante para colocar em suas mãos desde cedo — mas com a qual ainda não assumimos um compromisso definitivo. Estamos publicando-a *porque* queremos aprender com o uso real antes de decidir se ela se tornará uma parte permanente e com suporte do Wails.

Isso significa que alguns pontos se aplicam a tudo nesta seção:

- **É opcional.** Os experimentos nunca alteram o comportamento padrão de `wails3`. Você os ativa deliberadamente (geralmente com uma variável de ambiente ou uma flag de compilação) e, enquanto estiverem desativados, nada muda no seu fluxo de trabalho atual.
- **Pode não sobreviver.** Alguns experimentos evoluem para recursos estáveis. Outros são reformulados a ponto de ficarem irreconhecíveis ou são abandonados. Preferimos testar ideias em público e aprender rapidamente a lançar apenas aquilo sobre o qual já temos certeza.
- **A API não está congelada.** Nomes, flags, valores padrão e comportamentos podem mudar entre versões enquanto um experimento toma forma. As notas de versão destacarão as mudanças, mas não espere as garantias de estabilidade oferecidas pelos recursos estáveis.

## Queremos seu feedback

Esta é a parte importante. A continuidade dos experimentos depende do retorno que recebemos de quem realmente os utiliza. Se você testar um deles, queremos mesmo saber:

- Funcionou no seu projeto? Em que pontos falhou?
- Foi mais rápido, mais claro ou mais agradável — ou não valeu a pena fazer a mudança?
- O que precisaria acontecer para você usá-lo por padrão?

O feedback mais útil é concreto: o que você executou, o que esperava e o que realmente aconteceu. Cada experimento tem seu próprio tópico na categoria **Experimentos** do GitHub Discussions:

@container{display="grid" columns="2" gap="1rem"}
@linkcard{title="Discussões sobre experimentos" href="https://github.com/wailsapp/wails/discussions/categories/experiments" description="Encontre o tópico do experimento que você está usando e conte-nos como foi, o que não funcionou ou o que está faltando."}
@end

## Experimentos atuais

@container{display="grid" columns="2" gap="1rem"}
@linkcard{title="Wake" href="/experimental/wake/" description="Um executor de builds alternativo e integrado ao Wails para seus Taskfiles atuais. Builds incrementais mais rápidos, saída estruturada e execução paralela por padrão."}
@linkcard{title="Controle por LLM (MCP)" href="/guides/mcp-service/" description="Um servidor Model Context Protocol integrado que permite aos agentes de LLM inspecionar, testar e controlar um aplicativo Wails em execução — sem exigir código do usuário e ativado por uma tag de compilação."}
@end
