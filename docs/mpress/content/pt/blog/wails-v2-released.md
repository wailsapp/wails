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

Hoje é o lançamento do [Wails](https://wails.io) v2. Já se passaram cerca de 18 meses desde a primeira versão alfa da v2 e aproximadamente um ano desde o lançamento da primeira versão beta. Sou verdadeiramente grato a todos que participaram da evolução do projeto.

Parte do motivo de ter levado tanto tempo foi o desejo de atingir algum nível de completude antes de chamá-la oficialmente de v2. A verdade é que nunca há um momento perfeito para marcar uma versão: sempre existem problemas pendentes ou “só mais um” recurso para incluir. No entanto, marcar uma versão principal imperfeita proporciona alguma estabilidade aos usuários do projeto e também uma espécie de recomeço para os desenvolvedores.

Esta versão superou todas as minhas expectativas. Espero que ela proporcione a você tanto prazer quanto seu desenvolvimento proporcionou a nós.

## O que *é* o Wails?

Se você ainda não conhece o Wails, ele é um projeto que permite aos programadores Go criar interfaces avançadas para seus programas em Go usando tecnologias web conhecidas. É uma alternativa leve ao Electron, escrita em Go. Há muito mais informações no [site oficial](https://wails.io/docs/introduction).

## Quais são as novidades?

A versão v2 representa um enorme avanço para o projeto e resolve muitos dos pontos problemáticos da v1. Se você ainda não leu as publicações do blog sobre as versões beta para [macOS](/blog/wails-v2-beta-for-mac/), [Windows](/blog/wails-v2-beta-for-windows/) ou [Linux](/blog/wails-v2-beta-for-linux/), recomendo que as leia, pois elas abordam todas as principais mudanças em mais detalhes. Em resumo:

- Componente WebView2 para Windows, compatível com padrões web modernos e recursos de depuração.
- [Tema escuro/claro](https://wails.io/docs/reference/options#theme) + [temas personalizados](https://wails.io/docs/reference/options#customtheme) no Windows.
- O Windows agora não exige CGO.
- Compatibilidade pronta para uso com modelos de projeto Svelte, Vue, React, Preact, Lit e Vanilla.
- Integração com o [Vite](https://vitejs.dev/), que oferece um ambiente de desenvolvimento com recarregamento automático para seu aplicativo.
- [Menus](https://wails.io/docs/guides/application-development#application-menu) e [caixas de diálogo](https://wails.io/docs/reference/runtime/dialog) nativos do aplicativo.
- Efeitos nativos de translucidez de janelas para [Windows](https://wails.io/docs/reference/options#windowistranslucent) e [macOS](https://wails.io/docs/reference/options#windowistranslucent-1). Compatibilidade com planos de fundo Mica e Acrylic.
- Gere facilmente um [instalador NSIS](https://wails.io/docs/guides/windows-installer) para implantações no Windows.
- Uma ampla [biblioteca de runtime](https://wails.io/docs/reference/runtime/intro) que fornece métodos utilitários para manipulação de janelas, eventos, caixas de diálogo, menus e registro de logs.
- Compatibilidade com a [ofuscação](https://wails.io/docs/guides/obfuscated) do seu aplicativo usando o [garble](https://github.com/burrowers/garble).
- Compatibilidade com a compactação do seu aplicativo usando o [UPX](https://upx.github.io/).
- Geração automática de TypeScript a partir de structs Go. Mais informações [aqui](https://wails.io/docs/howdoesitwork#calling-bound-go-methods).
- Não é necessário distribuir bibliotecas adicionais nem DLLs com seu aplicativo. Em nenhuma plataforma.
- Não é necessário empacotar os ativos do frontend. Basta desenvolver seu aplicativo como qualquer outro aplicativo web.

## Créditos e agradecimentos

Chegar à v2 exigiu um esforço enorme. Houve cerca de 2200 commits feitos por 89 colaboradores entre a versão alfa inicial e o lançamento de hoje, além de muitas outras pessoas que forneceram traduções, testes, feedback e ajuda tanto nos fóruns de discussão quanto no rastreador de problemas. Sou imensamente grato a cada um de vocês. Também gostaria de agradecer de maneira especialmente calorosa a todos os patrocinadores do projeto que ofereceram orientação, aconselhamento e feedback. Tudo o que vocês fazem é muito valorizado.

Há algumas pessoas que eu gostaria de mencionar em especial:

Primeiramente, um agradecimento **enorme** ao [@stffabi](https://github.com/stffabi), que fez tantas contribuições das quais todos nós nos beneficiamos, além de oferecer muito suporte em diversos problemas. Ele implementou alguns recursos essenciais, como a compatibilidade com servidores de desenvolvimento externos, que transformou nosso modo de desenvolvimento ao permitir a integração com os superpoderes do [Vite](https://vitejs.dev/). É justo dizer que o Wails v2 seria uma versão muito menos empolgante sem suas [contribuições incríveis](https://github.com/wailsapp/wails/commits?author=stffabi&since=2020-01-04). Muito obrigado, @stffabi!

Também gostaria de fazer um enorme agradecimento ao [@misitebao](https://github.com/misitebao), que vem trabalhando incansavelmente na manutenção do site, além de fornecer traduções para o chinês, gerenciar o Crowdin e ajudar novos tradutores a se familiarizarem com o projeto. Essa é uma tarefa extremamente importante, e sou imensamente grato por todo o tempo e esforço dedicados a ela! Você é demais!

Por último, mas não menos importante, um enorme agradecimento a Mat Ryer, que ofereceu orientação e suporte durante o desenvolvimento da v2. Escrever o xBar em conjunto usando uma versão alfa inicial da v2 ajudou a definir os rumos da v2 e também me permitiu compreender algumas falhas de design das primeiras versões. Tenho a satisfação de anunciar que, a partir de hoje, começaremos a migrar o xBar para o Wails v2, e ele se tornará o principal aplicativo de referência do projeto. Obrigado, Mat!

## Lições aprendidas

O caminho até a v2 trouxe várias lições que orientarão o desenvolvimento daqui em diante.

## Versões menores, mais rápidas e focadas

Durante o desenvolvimento da v2, muitos recursos e correções de bugs foram desenvolvidos de maneira pontual. Isso resultou em ciclos de lançamento mais longos e mais difíceis de depurar. Daqui em diante, criaremos versões com mais frequência e com um número menor de recursos. Cada lançamento incluirá atualizações da documentação e testes rigorosos. Esperamos que essas versões menores, mais rápidas e focadas resultem em menos regressões e em uma documentação de melhor qualidade.

## Incentivar a participação

Quando iniciei este projeto, eu queria ajudar imediatamente todas as pessoas que enfrentassem algum problema. Os problemas eram “pessoais”, e eu queria resolvê-los o mais rápido possível. Isso não é sustentável e, em última análise, prejudica a longevidade do projeto. Daqui em diante, darei mais espaço para que outras pessoas participem respondendo a perguntas e fazendo a triagem de problemas. Seria bom contar com algumas ferramentas para ajudar nisso; portanto, se você tiver alguma sugestão, participe da discussão [aqui](https://github.com/wailsapp/wails/discussions/1855).

## Aprender a dizer não

Quanto mais pessoas participam de um projeto de código aberto, mais solicitações surgem por recursos adicionais que podem ou não ser úteis para a maioria. Esses recursos exigem inicialmente tempo para desenvolvimento e depuração e, a partir daí, geram um custo contínuo de manutenção. Eu mesmo sou quem mais comete esse erro, pois muitas vezes quero “resolver tudo de uma vez” em vez de oferecer o recurso mínimo viável. Daqui em diante, precisaremos dizer “não” com um pouco mais de frequência à inclusão de recursos no núcleo do projeto e concentrar nossos esforços em uma maneira de capacitar os desenvolvedores a fornecerem essa funcionalidade por conta própria. Estamos considerando seriamente o uso de plugins nesse cenário. Isso permitirá que qualquer pessoa estenda o projeto como considerar adequado, além de oferecer uma maneira fácil de contribuir com o projeto.

## Olhando para o futuro

Já há muitos recursos essenciais que estamos considerando adicionar ao Wails no próximo ciclo principal de desenvolvimento. O  
[roadmap](https://github.com/wailsapp/wails/discussions/1484) está repleto de  
ideias interessantes, e estou ansioso para começar a trabalhar nelas. Uma das grandes solicitações tem sido o suporte a várias janelas. É algo complicado de implementar corretamente, e talvez precisemos considerar uma API alternativa, pois a atual não foi projetada com isso em mente. Com base em algumas ideias preliminares e no feedback recebido, acho que você vai gostar do rumo que estamos considerando para esse recurso.

Pessoalmente, estou muito empolgado com a possibilidade de executar aplicativos Wails em dispositivos móveis. Já temos um projeto de demonstração que mostra que é possível executar um aplicativo Wails no Android, então estou muito ansioso para explorar até onde podemos chegar com isso!

Um último ponto que gostaria de abordar é a paridade de recursos. Há muito tempo, um dos nossos princípios fundamentais é não adicionar nada ao projeto sem que haja suporte multiplataforma completo. Embora isso tenha se mostrado, em grande parte, viável até agora, essa exigência realmente atrasou o lançamento de novos recursos no projeto. Daqui em diante, adotaremos uma abordagem um pouco diferente: qualquer recurso novo que não possa ser lançado imediatamente para todas as plataformas será disponibilizado por meio de uma configuração ou API experimental. Isso permitirá que os primeiros usuários em determinadas plataformas testem o recurso e forneçam feedback que contribuirá para seu design final. Isso significa, é claro, que não haverá garantias de estabilidade da API até que o recurso seja totalmente compatível com todas as plataformas nas quais possa ser implementado, mas pelo menos isso destravará o desenvolvimento.

## Palavras finais

Tenho muito orgulho do que conseguimos alcançar com a versão V2. É incrível ver o que as pessoas já conseguiram criar usando as versões beta até agora. Aplicativos de qualidade como [Varly](https://varly.app/),  
[Surge](https://getsurge.io/) e [October](https://october.utf9k.net/). Recomendo que você os conheça.

Esta versão foi possível graças ao trabalho árduo de muitos colaboradores. Embora possa ser baixada e usada gratuitamente, ela não surgiu sem custos. Não se engane: este projeto teve um custo considerável. Esse custo não inclui apenas o meu tempo e o de cada colaborador, mas também o tempo que todas essas pessoas deixaram de passar com amigos e familiares. É por isso que sou extremamente grato por cada segundo dedicado a tornar este projeto uma realidade. Quanto mais colaboradores tivermos, mais poderemos distribuir esse esforço e mais poderemos alcançar juntos. Gostaria de incentivar todos vocês a escolher algo com que possam contribuir, seja confirmando o bug relatado por alguém, sugerindo uma correção, alterando a documentação ou ajudando alguém que precise. Todas essas pequenas coisas causam um impacto enorme! Seria incrível se você também fizesse parte da história da nossa jornada até a v3.

Aproveite!

&dash; Lea

PS: Se você ou sua empresa consideram o Wails útil, pensem na possibilidade de  
[patrocinar o projeto](https://github.com/sponsors/leaanthony). Obrigado!
