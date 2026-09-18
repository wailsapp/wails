---
title: "Wails v2 Beta para Windows"
description: "Notas de versão e anúncios do Wails"
authors: ["leaanthony"]
tags: ["wails","v2"]
date: "2021-09-27"
slug: "blog/wails-v2-beta-for-windows"
image: "/assets/blog-images/wails.webp"
sourcePath: "blog/wails-v2-beta-for-windows.md"
---

![captura de tela do Wails](/assets/blog-images/wails.webp)

Quando anunciei o Wails pela primeira vez no Reddit, há pouco mais de 2 anos, de dentro de um trem em Sydney, não esperava que ele recebesse muita atenção. Alguns dias depois, um prolífico vlogger de tecnologia publicou um vídeo tutorial, fez uma avaliação positiva e, desde então, o interesse pelo projeto disparou.

Ficou claro que as pessoas estavam entusiasmadas com a possibilidade de adicionar frontends web aos seus projetos Go e, quase imediatamente, levaram o projeto além da prova de conceito que eu havia criado. Na época, o Wails usava o projeto [webview](https://github.com/webview/webview) para gerenciar o frontend, e a única opção para Windows era o mecanismo de renderização do IE11. Muitos relatos de bugs tinham origem nessa limitação: suporte precário a JavaScript/CSS e ausência de ferramentas de desenvolvimento para depuração. Isso tornava a experiência de desenvolvimento frustrante, mas não havia muito que pudesse ser feito para corrigir o problema.

Durante muito tempo, acreditei firmemente que a Microsoft acabaria tendo que resolver a situação de seus navegadores. O mundo avançava, o desenvolvimento frontend estava em plena expansão e o IE já não dava conta. Quando a Microsoft anunciou que passaria a usar o Chromium como base de sua nova estratégia para navegadores, soube que seria apenas uma questão de tempo até que o Wails pudesse usá-lo e elevasse a experiência dos desenvolvedores no Windows a outro nível.

Hoje, tenho o prazer de anunciar: **Wails v2 Beta para Windows**! Há muita coisa para explorar nesta versão, então pegue uma bebida, sente-se e vamos começar...

## Sem dependência do CGO!

Não, não estou brincando: *nenhuma* *dependência* do *CGO* 🤯! A questão do Windows é que, ao contrário do MacOS e do Linux, ele não inclui um compilador padrão. Além disso, o CGO exige um compilador mingw, e há inúmeras opções de instalação diferentes. A remoção do requisito do CGO simplificou enormemente a configuração e também facilitou muito a depuração. Embora eu tenha me empenhado bastante para fazer isso funcionar, a maior parte do crédito deve ir para [John Chadwick](https://github.com/jchv), não apenas por iniciar alguns projetos que tornaram isso possível, mas também por permitir que outra pessoa assumisse esses projetos e os desenvolvesse ainda mais. O crédito também vai para [Tad Vizbaras](https://github.com/tadvi), cujo projeto [winc](https://github.com/tadvi/winc) me colocou nesse caminho.

### Mecanismo de renderização Chromium do WebView2

![captura de tela das ferramentas de desenvolvimento](/assets/blog-images/devtools.png)

Finalmente, os desenvolvedores para Windows têm um mecanismo de renderização de primeira classe para seus aplicativos! Acabaram-se os dias de contorcer o código do frontend para que funcionasse no Windows. Além disso, você conta com ferramentas de desenvolvimento de primeira classe!

No entanto, o componente WebView2 exige que o `WebView2Loader.dll` fique ao lado do binário. Isso torna a distribuição um pouco mais trabalhosa do que nós, gophers, estamos acostumados. Todas as soluções e bibliotecas que usam o WebView2 (até onde sei) têm essa dependência.

No entanto, estou muito entusiasmado em anunciar que os aplicativos Wails *não têm esse requisito*! Graças à magia de [John Chadwick](https://github.com/jchv), podemos incorporar essa dll ao binário e fazer o Windows carregá-la como se estivesse presente no disco.

Alegrem-se, gophers! O sonho de um único binário continua vivo!

### Novos recursos

![captura de tela dos menus do Wails](/assets/blog-images/wails-menus.webp)

Houve muitos pedidos de suporte a menus nativos. Finalmente, o Wails oferece esse recurso. Os menus de aplicativo agora estão disponíveis e incluem suporte à maioria dos recursos de menus nativos. Isso inclui itens de menu padrão, caixas de seleção, grupos de botões de opção, submenus e separadores.

Na v1, houve uma enorme quantidade de pedidos para ter maior controle sobre a própria janela. Tenho o prazer de anunciar que há novas APIs de runtime específicas para isso. Elas oferecem muitos recursos e são compatíveis com configurações de vários monitores. Há também uma API de caixas de diálogo aprimorada: agora você pode ter caixas de diálogo modernas e nativas, com ampla configuração para atender a todas as suas necessidades.

Agora existe a opção de gerar a configuração do IDE junto com o projeto. Isso significa que, se você abrir o projeto em um IDE compatível, ele já estará configurado para compilar e depurar o aplicativo. No momento, há suporte ao VSCode, mas esperamos oferecer suporte a outros IDEs, como o Goland, em breve.

![captura de tela do VSCode](/assets/blog-images/vscode.webp)

### Não é necessário empacotar os recursos

Um grande ponto problemático da v1 era a necessidade de condensar todo o aplicativo em arquivos JS e CSS únicos. Tenho o prazer de anunciar que, na v2, não é necessário empacotar os recursos de forma alguma. Quer carregar uma imagem local? Use uma tag `<img>` com um caminho local no atributo src. Quer usar uma fonte interessante? Copie-a para o projeto e adicione o caminho correspondente ao CSS.

> Uau, isso parece um servidor web...

Sim, funciona exatamente como um servidor web, mas não é um.

> Então, como incluo meus recursos?

Basta passar um único `embed.FS` que contenha todos os seus recursos para a configuração do aplicativo. Eles nem precisam estar no diretório principal — o Wails resolve isso para você.

### Nova experiência de desenvolvimento

![captura de tela do navegador](/assets/blog-images/browser.webp)

Agora que os recursos não precisam ser empacotados, tornou-se possível oferecer uma experiência de desenvolvimento totalmente nova. O novo comando `wails dev` compila e executa o aplicativo, mas, em vez de usar os recursos contidos no `embed.FS`, ele os carrega diretamente do disco.

Ele também oferece os seguintes recursos adicionais:

- Recarga automática — qualquer alteração nos recursos do frontend aciona uma recarga automática do frontend do aplicativo
- Recompilação automática — qualquer alteração no código Go recompila e reinicia o aplicativo

Além disso, um servidor web será iniciado na porta 34115. Ele disponibilizará o aplicativo para qualquer navegador que se conectar. Todos os navegadores web conectados responderão a eventos do sistema, como a recarga automática após uma alteração nos recursos.

Em Go, estamos acostumados a trabalhar com structs em nossos aplicativos. Muitas vezes, é útil enviar structs ao frontend e usá-las como estado no aplicativo. Na v1, esse era um processo bastante manual e um tanto trabalhoso para o desenvolvedor. Tenho o prazer de anunciar que, na v2, qualquer aplicativo executado no modo de desenvolvimento gerará automaticamente modelos TypeScript para todas as structs usadas como parâmetros de entrada ou saída de métodos vinculados. Isso permite a troca transparente de modelos de dados entre os dois ambientes.

Além disso, outro módulo JS é gerado dinamicamente para encapsular todos os seus métodos vinculados. Ele fornece JSDoc para os métodos, oferecendo preenchimento de código e dicas no IDE. É muito legal ver os modelos de dados serem importados automaticamente quando você pressiona Tab em um módulo gerado automaticamente que encapsula seu código Go!

### Modelos remotos

![captura de tela remota](/assets/blog-images/remote.webp)

Colocar uma aplicação em funcionamento rapidamente sempre foi um dos principais objetivos do projeto Wails. Quando o lançamos, tentamos abranger muitos dos frameworks modernos da época: React, Vue e Angular. O mundo do desenvolvimento frontend é repleto de opiniões divergentes, evolui rapidamente e é difícil de acompanhar! Por isso, nossos templates básicos ficavam desatualizados muito depressa, o que gerava uma enorme dor de cabeça de manutenção. Isso também significava que não tínhamos templates modernos e interessantes para as melhores e mais recentes pilhas de tecnologia.

Com a v2, eu queria dar mais autonomia à comunidade, permitindo que vocês mesmos criassem e hospedassem templates, em vez de dependerem do projeto Wails. Agora, portanto, vocês podem criar projetos usando templates mantidos pela comunidade! Espero que isso inspire os desenvolvedores a criar um ecossistema dinâmico de templates de projeto. Estou realmente muito empolgada para ver o que nossa comunidade de desenvolvedores pode criar!

### Conclusão

O Wails v2 representa uma nova base para o projeto. O objetivo desta versão é obter feedback sobre a nova abordagem e corrigir quaisquer bugs antes do lançamento da versão final. Sua contribuição será muito bem-vinda. Envie seus comentários para o fórum de discussão [Beta da v2](https://github.com/wailsapp/wails/discussions/828).

Houve muitas reviravoltas, mudanças de rumo e retrocessos até chegarmos a este ponto. Isso ocorreu em parte devido a decisões técnicas iniciais que precisaram ser alteradas e, em parte, porque alguns problemas fundamentais, para os quais havíamos dedicado tempo à criação de soluções alternativas, foram corrigidos nos projetos de origem: o recurso embed do Go é um bom exemplo. Felizmente, tudo se encaixou no momento certo e hoje temos a melhor solução possível. Acredito que a espera valeu a pena — isso não teria sido possível nem mesmo 2 meses atrás.

Também preciso agradecer imensamente :pray: às pessoas a seguir, pois, sem elas, esta versão simplesmente não existiria:

- [Misite Bao](https://github.com/misitebao) — Trabalhou incansavelmente nas traduções para o chinês e é uma pessoa incrível na descoberta de bugs.
- [John Chadwick](https://github.com/jchv) — Seu trabalho excepcional no [go-webview2](https://github.com/jchv/go-webview2) e no [go-winloader](https://github.com/jchv/go-winloader) tornou possível a versão para Windows que temos hoje.
- [Tad Vizbaras](https://github.com/tadvi) — Os experimentos com seu projeto [winc](https://github.com/tadvi/winc) foram o primeiro passo no caminho para um Wails escrito inteiramente em Go.
- [Mat Ryer](https://github.com/matryer) — Seu apoio, incentivo e feedback realmente ajudaram a impulsionar o projeto.

Por fim, gostaria de fazer um agradecimento especial a todos os [patrocinadores do projeto](/credits/#sponsors), incluindo a [JetBrains](https://www.jetbrains.com?from=Wails), cujo apoio impulsiona o projeto de muitas maneiras nos bastidores.

Mal posso esperar para ver o que as pessoas criarão com o Wails nesta nova e empolgante fase do projeto!

Lea.

PS: Usuários do macOS e do Linux não precisam se sentir excluídos — a migração para essa nova base está em andamento, e a maior parte do trabalho difícil já foi concluída. Aguardem só mais um pouco!

PPS: Se você ou sua empresa considera o Wails útil, pense na possibilidade de [patrocinar o projeto](https://github.com/sponsors/leaanthony). Agradeço!
