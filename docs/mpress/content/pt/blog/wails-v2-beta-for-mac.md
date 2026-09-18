---
title: "Wails v2 Beta para MacOS"
description: "Notas de versão e anúncios do Wails"
authors: ["leaanthony"]
tags: ["wails","v2"]
date: "2021-11-08"
slug: "blog/wails-v2-beta-for-mac"
image: "/assets/blog-images/wails-mac.webp"
sourcePath: "blog/wails-v2-beta-for-mac.md"
---

![captura de tela do wails-mac](/assets/blog-images/wails-mac.webp)

Hoje marca o lançamento da primeira versão beta do Wails v2 para Mac! Levou bastante tempo para chegarmos até aqui, e espero que a versão de hoje ofereça algo razoavelmente útil. Houve várias reviravoltas ao longo do caminho, e espero, com a sua ajuda, corrigir os problemas restantes e deixar a versão para Mac bem-acabada para o lançamento final do v2.

Quer dizer que esta versão não está pronta para produção? Para o seu caso de uso, é bem possível que esteja, mas ainda há vários problemas conhecidos. Portanto, acompanhe [este quadro do projeto](https://github.com/wailsapp/wails/projects/7) e, se quiser contribuir, sua ajuda será muito bem-vinda!

Então, o que há de novo no Wails v2 para Mac em comparação com o v1? Dica: é bem parecido com a versão Beta para Windows :wink:

## Novos recursos

![captura de tela dos menus do Wails no Mac](/assets/blog-images/wails-menus-mac.webp)

Recebemos muitos pedidos de suporte a menus nativos. Finalmente, o Wails oferece esse recurso. Agora, menus de aplicativos estão disponíveis e incluem suporte à maioria dos recursos de menus nativos, como itens de menu padrão, caixas de seleção, grupos de botões de opção, submenus e separadores.

No v1, recebemos uma enorme quantidade de pedidos para oferecer maior controle sobre a própria janela. Tenho o prazer de anunciar que agora há novas APIs de runtime específicas para isso. Elas oferecem muitos recursos e são compatíveis com configurações de vários monitores. Também há uma API de caixas de diálogo aprimorada: agora você pode usar caixas de diálogo modernas e nativas, com diversas opções de configuração para atender a todas as suas necessidades.

### Opções específicas para Mac

Além das opções normais do aplicativo, o Wails v2 para Mac também oferece alguns extras para Mac:

- Deixe sua janela toda estilosa e translúcida, como os belos aplicativos em Swift!
- Barra de título altamente personalizável
- Oferecemos suporte às opções NSAppearance para o aplicativo
- Configuração simples para criar automaticamente um menu "Sobre"

### Não é necessário empacotar os ativos

Um enorme problema do v1 era a necessidade de condensar todo o aplicativo em arquivos JS e CSS únicos. Tenho o prazer de anunciar que, no v2, não é necessário empacotar os ativos de nenhuma maneira. Quer carregar uma imagem local? Use uma tag `<img>` com um caminho src local. Quer usar uma fonte interessante? Copie-a para o projeto e adicione o caminho correspondente ao CSS.

> Uau, isso parece um servidor web...

Sim, funciona como um servidor web, mas não é um.

> Então, como incluo meus ativos?

Basta passar um único `embed.FS` contendo todos os seus ativos para a configuração do aplicativo. Eles nem precisam estar no diretório de nível superior — o Wails cuidará disso para você.

### Nova experiência de desenvolvimento

Agora que os ativos não precisam ser empacotados, tornou-se possível oferecer uma experiência de desenvolvimento totalmente nova. O novo comando `wails dev` compila e executa o aplicativo, mas, em vez de usar os ativos no `embed.FS`, carrega-os diretamente do disco.

Ele também oferece estes recursos adicionais:

- Recarregamento a quente — qualquer alteração nos ativos do frontend aciona o recarregamento automático do frontend do aplicativo
- Recompilação automática — qualquer alteração no código Go recompila e reinicia o aplicativo

Além disso, um servidor web será iniciado na porta 34115. Ele disponibilizará seu aplicativo a qualquer navegador que se conectar. Todos os navegadores conectados responderão a eventos do sistema, como o recarregamento a quente após uma alteração nos ativos.

Em Go, estamos acostumados a trabalhar com structs em nossos aplicativos. Muitas vezes, é útil enviar structs ao frontend e usá-las como estado do aplicativo. No v1, esse processo era muito manual e um tanto trabalhoso para quem desenvolvia. Tenho o prazer de anunciar que, no v2, qualquer aplicativo executado no modo de desenvolvimento gerará automaticamente modelos TypeScript para todas as structs usadas como parâmetros de entrada ou saída de métodos vinculados. Isso permite o intercâmbio transparente de modelos de dados entre os dois ambientes.

Além disso, outro módulo JS é gerado dinamicamente para encapsular todos os seus métodos vinculados. Ele fornece JSDoc para os métodos, permitindo o preenchimento automático de código e a exibição de dicas no IDE. É muito interessante ver os modelos de dados serem importados automaticamente ao pressionar Tab em um módulo gerado automaticamente que encapsula seu código Go!

### Modelos remotos

![captura de tela do remote-mac](/assets/blog-images/remote-mac.webp)

Possibilitar que um aplicativo fosse colocado em funcionamento rapidamente sempre foi um objetivo central do projeto Wails. Quando fizemos o lançamento, tentamos abranger muitos dos frameworks modernos da época: react, vue e angular. O mundo do desenvolvimento frontend é cheio de opiniões fortes, muda rapidamente e é difícil de acompanhar! Como resultado, nossos modelos básicos ficavam desatualizados muito depressa, o que gerava uma grande carga de manutenção. Isso também significava que não tínhamos modelos modernos e interessantes para as melhores e mais recentes stacks de tecnologia.

No v2, eu queria dar mais autonomia à comunidade, permitindo que vocês mesmos criassem e hospedassem modelos, em vez de dependerem do projeto Wails. Agora, portanto, você pode criar projetos usando modelos mantidos pela comunidade! Espero que isso inspire os desenvolvedores a criar um ecossistema dinâmico de modelos de projeto. Estou realmente entusiasmado para ver o que nossa comunidade de desenvolvedores pode criar!

### Suporte nativo ao M1

Graças ao incrível apoio de [Mat Ryer](https://github.com/matryer/), o projeto Wails agora oferece suporte a builds nativos para M1:

![captura de tela do build-darwin-arm](/assets/blog-images/build-darwin-arm.webp)

Você também pode especificar `darwin/amd64` como destino:

![captura de tela do build-darwin-amd](/assets/blog-images/build-darwin-amd.webp)

Ah, quase me esqueci... você também pode usar `darwin/universal`... :wink:

![captura de tela do build-darwin-universal](/assets/blog-images/build-darwin-universal.webp)

### Compilação cruzada para Windows

Como o Wails v2 para Windows é escrito inteiramente em Go, você pode gerar builds para Windows sem usar docker.

![captura de tela do build-cross-windows](/assets/blog-images/build-cross-windows.webp)  
bu

### Renderizador WKWebView

V1 dependia de um componente WebView (agora obsoleto). V2 usa o componente WKWebKit mais recente, portanto você pode esperar os melhores e mais novos recursos da Apple.

### Conclusão

Como eu disse nas notas de versão do Windows, o Wails v2 representa uma nova base para o projeto. O objetivo desta versão é receber feedback sobre a nova abordagem e corrigir eventuais bugs antes do lançamento da versão final. Sua contribuição será muito bem-vinda! Envie seus comentários ao fórum de discussão da [versão beta da v2](https://github.com/wailsapp/wails/discussions/828).

Por fim, gostaria de agradecer especialmente a todos os [patrocinadores do projeto](/credits/#sponsors), incluindo a [JetBrains](https://www.jetbrains.com?from=Wails), cujo apoio impulsiona o projeto de várias maneiras nos bastidores.

Estou ansiosa para ver o que as pessoas criarão com o Wails nesta nova e empolgante fase do projeto!

Lea.

PS: Usuários do Linux, vocês são os próximos!

PPS: Se você ou sua empresa consideram o Wails útil, considere [patrocinar o projeto](https://github.com/sponsors/leaanthony). Obrigada!
