---
title: "Wails v2 lançado"
description: "Notas de versão e anúncios do Wails"
authors: ["leaanthony"]
tags: ["wails","v2"]
date: "2022-09-22"
slug: "blog/wails-v2-released"
image: "/assets/blog-images/montage.png"
sourcePath: "blog/wails-v2-released.md"
---

![captura de tela da montagem](/assets/blog-images/montage.png)

## Chegou!

Hoje marca o lançamento do [Wails](https://wails.io) v2. Já se passaram cerca de 18 meses desde a primeira versão alfa do v2 e aproximadamente um ano desde a primeira versão beta. Sou muito grato a todos que participaram da evolução do projeto.

Parte da razão para ter demorado tanto foi o desejo de alcançar alguma definição de completude antes de chamá-lo oficialmente de v2. A verdade é que nunca existe um momento perfeito para marcar uma versão: sempre há problemas pendentes ou “só mais um” recurso para incluir. No entanto, marcar uma versão principal imperfeita proporciona certa estabilidade aos usuários do projeto, além de representar uma espécie de recomeço para os desenvolvedores.

Esta versão é muito mais do que eu jamais esperava que fosse. Espero que ela proporcione a você tanta satisfação quanto o desenvolvimento dela proporcionou a nós.

## O que *é* o Wails?

Se você ainda não conhece o Wails, ele é um projeto que permite aos programadores Go criar interfaces avançadas para seus programas em Go usando tecnologias web conhecidas. É uma alternativa leve ao Electron, escrita em Go. Você encontra muito mais informações no [site oficial](https://v2.wails.io/docs/introduction).

## O que há de novo?

A versão v2 representa um enorme avanço para o projeto e resolve muitos dos pontos problemáticos da v1. Se você ainda não leu as publicações do blog sobre as versões beta para [macOS](/blog/wails-v2-beta-for-mac/), [Windows](/blog/wails-v2-beta-for-windows/) ou [Linux](/blog/wails-v2-beta-for-linux/), recomendo que as leia, pois elas abordam com mais detalhes todas as principais mudanças. Em resumo:

- Componente Webview2 para Windows, com suporte a padrões modernos da web e recursos de depuração.
- [Tema escuro/claro](https://v2.wails.io/docs/reference/options#theme) e [temas personalizados](https://v2.wails.io/docs/reference/options#customtheme) no Windows.
- O Windows agora não exige CGO.
- Suporte pronto para uso a modelos de projeto Svelte, Vue, React, Preact, Lit e Vanilla.
- Integração com o [Vite](https://vitejs.dev/), que oferece um ambiente de desenvolvimento com recarregamento automático para seu aplicativo.
- [Menus](https://v2.wails.io/docs/guides/application-development#application-menu) e [caixas de diálogo](https://v2.wails.io/docs/reference/runtime/dialog) nativos do aplicativo.
- Efeitos nativos de translucidez de janelas para [Windows](https://v2.wails.io/docs/reference/options#windowistranslucent) e [macOS](https://v2.wails.io/docs/reference/options#windowistranslucent-1). Suporte aos planos de fundo Mica e Acrylic.
- Gere facilmente um [instalador NSIS](https://v2.wails.io/docs/guides/windows-installer) para implantações no Windows.
- Uma completa [biblioteca de runtime](https://v2.wails.io/docs/reference/runtime/intro) que oferece métodos utilitários para manipulação de janelas, eventos, caixas de diálogo, menus e registro de logs.
- Suporte para [ofuscar](https://v2.wails.io/docs/guides/obfuscated) seu aplicativo usando o [garble](https://github.com/burrowers/garble).
- Suporte para compactar seu aplicativo usando o [UPX](https://upx.github.io/).
- Geração automática de TypeScript a partir de structs Go. Mais informações [aqui](https://v2.wails.io/docs/howdoesitwork#calling-bound-go-methods).
- Não é necessário distribuir bibliotecas adicionais nem DLLs com seu aplicativo, em nenhuma plataforma.
- Não é necessário empacotar os recursos do frontend. Basta desenvolver seu aplicativo como qualquer outro aplicativo web.

## Créditos e agradecimentos

Chegar à versão v2 exigiu um esforço enorme. Entre a versão alfa inicial e o lançamento de hoje, houve cerca de 2200 commits de 89 colaboradores, além de muitas outras pessoas que contribuíram com traduções, testes, feedback e ajuda nos fóruns de discussão e no rastreador de problemas. Sou imensamente grato a cada um de vocês. Também gostaria de fazer um agradecimento muito especial a todos os patrocinadores do projeto que ofereceram orientação, aconselhamento e feedback. Tudo o que vocês fazem é profundamente valorizado.

Há algumas pessoas que eu gostaria de mencionar especialmente:

Primeiramente, um agradecimento **enorme** ao [@stffabi](https://github.com/stffabi), que fez tantas contribuições das quais todos nós nos beneficiamos, além de oferecer muito suporte em diversos problemas. Ele implementou alguns recursos essenciais, como o suporte a servidores externos de desenvolvimento, que transformou nosso modo de desenvolvimento ao nos permitir aproveitar os superpoderes do [Vite](https://vitejs.dev/). É justo dizer que o Wails v2 seria uma versão muito menos empolgante sem suas [contribuições incríveis](https://github.com/wailsapp/wails/commits?author=stffabi&since=2020-01-04). Muito obrigado, @stffabi!

Também gostaria de fazer um enorme agradecimento ao [@misitebao](https://github.com/misitebao), que tem trabalhado incansavelmente na manutenção do site, além de fornecer traduções para o chinês, gerenciar o Crowdin e ajudar novos tradutores a se familiarizarem com o projeto. Essa é uma tarefa extremamente importante, e sou imensamente grato por todo o tempo e esforço dedicados a ela! Você é incrível!

Por último, mas não menos importante, um enorme agradecimento a Mat Ryer, que ofereceu aconselhamento e suporte durante o desenvolvimento da versão v2. Escrever o xBar juntos usando uma versão alfa inicial da v2 ajudou a definir a direção da v2 e também me permitiu compreender algumas falhas de design das primeiras versões. Tenho o prazer de anunciar que, a partir de hoje, começaremos a portar o xBar para o Wails v2, e ele se tornará o principal aplicativo de referência do projeto. Obrigado, Mat!

## Lições aprendidas

O caminho até a versão v2 trouxe várias lições que orientarão o desenvolvimento daqui em diante.

## Versões menores, mais rápidas e focadas

Durante o desenvolvimento da versão v2, muitos recursos e correções de bugs foram desenvolvidos de maneira pontual. Isso resultou em ciclos de lançamento mais longos e tornou a depuração mais difícil. Daqui em diante, lançaremos versões com mais frequência e com um número menor de recursos. Cada lançamento incluirá atualizações da documentação e testes completos. Esperamos que essas versões menores, mais rápidas e focadas resultem em menos regressões e em uma documentação de melhor qualidade.

## Incentivar a participação

Quando iniciei este projeto, eu queria ajudar imediatamente todas as pessoas que enfrentassem algum problema. Eu levava os problemas para o lado “pessoal” e queria resolvê-los o mais rápido possível. Isso não é sustentável e, no fim, prejudica a longevidade do projeto. Daqui em diante, darei mais espaço para que outras pessoas participem respondendo a perguntas e fazendo a triagem de problemas. Seria bom contar com algumas ferramentas para ajudar nisso; portanto, se você tiver alguma sugestão, participe da discussão [aqui](https://github.com/wailsapp/wails/discussions/1855).

## Aprender a dizer não

Quanto mais pessoas participam de um projeto de código aberto, mais solicitações surgem por recursos adicionais que podem ou não ser úteis para a maioria. Esses recursos exigem um investimento inicial de tempo para desenvolvimento e depuração e, a partir daí, geram um custo contínuo de manutenção. Eu mesmo sou quem mais comete esse erro, pois frequentemente quero “resolver tudo de uma vez” em vez de oferecer o recurso mínimo viável. Daqui em diante, precisaremos dizer “não” com um pouco mais de frequência à inclusão de recursos no núcleo e concentrar nossos esforços em uma forma de permitir que os próprios desenvolvedores ofereçam essas funcionalidades. Estamos considerando seriamente o uso de plugins nesse cenário. Isso permitirá que qualquer pessoa estenda o projeto como achar adequado, além de proporcionar uma maneira simples de contribuir com ele.

## Olhando para o futuro

Já estamos considerando muitos recursos essenciais para adicionar ao Wails no próximo grande ciclo de desenvolvimento. O [roteiro](https://github.com/wailsapp/wails/discussions/1484) está repleto de ideias interessantes, e estou ansioso para começar a trabalhar nelas. Uma das grandes solicitações tem sido o suporte a várias janelas. É algo complicado de implementar corretamente, e talvez precisemos oferecer uma API alternativa, pois a atual não foi projetada com isso em mente. Com base em algumas ideias preliminares e no feedback recebido, acredito que você gostará da direção que estamos considerando.

Pessoalmente, estou muito empolgado com a perspectiva de executar aplicativos Wails em dispositivos móveis. Já temos um projeto de demonstração que mostra ser possível executar um aplicativo Wails no Android, então estou muito interessado em explorar até onde podemos chegar com isso!

Um último ponto que gostaria de abordar é a paridade de recursos. Há muito tempo, um princípio fundamental é não adicionar nada ao projeto sem que haja suporte multiplataforma completo. Embora isso tenha se mostrado (em grande parte) viável até agora, essa postura realmente atrasou o lançamento de novos recursos no projeto. Daqui em diante, adotaremos uma abordagem um pouco diferente: qualquer recurso novo que não possa ser lançado imediatamente para todas as plataformas será disponibilizado por meio de uma configuração ou API experimental. Isso permitirá que os primeiros usuários em determinadas plataformas testem o recurso e forneçam feedback que contribuirá para o design final. É claro que isso significa que não haverá garantias de estabilidade da API até que o recurso tenha suporte completo em todas as plataformas nas quais for possível oferecer suporte, mas, pelo menos, isso desbloqueará o desenvolvimento.

## Considerações finais

Tenho muito orgulho do que conseguimos alcançar com a versão V2. É incrível ver o que as pessoas já conseguiram desenvolver usando as versões beta até agora. Aplicativos de qualidade como [Varly](https://varly.app/), [Surge](https://getsurge.io/) e [October](https://october.utf9k.net/). Recomendo que você os conheça.

Esta versão foi viabilizada pelo trabalho árduo de muitas pessoas que contribuíram para o projeto. Embora possa ser baixada e usada gratuitamente, ela não foi produzida sem custos. Não se engane: este projeto teve um custo considerável. Não envolveu apenas o meu tempo e o de cada uma dessas pessoas, mas também o tempo que todas elas passaram longe de amigos e familiares. Por isso, sou extremamente grato por cada segundo dedicado a tornar este projeto realidade. Quanto mais pessoas contribuírem, mais poderemos distribuir esse esforço e mais conseguiremos realizar juntos. Gostaria de incentivar cada um de vocês a escolher algo com que possa contribuir, seja confirmando o bug relatado por alguém, sugerindo uma correção, alterando a documentação ou ajudando alguém que precise. Todas essas pequenas ações causam um impacto enorme! Seria fantástico se você também fizesse parte da história da jornada até a versão v3.

Aproveite!

&dash; Lea

PS: se você ou sua empresa consideram o Wails útil, pensem em [patrocinar o projeto](https://github.com/sponsors/leaanthony). Obrigado!
