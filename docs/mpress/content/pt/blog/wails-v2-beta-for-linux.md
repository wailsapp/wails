---
title: "Wails v2 Beta para Linux"
description: "Notas de versão e anúncios do Wails"
authors: ["leaanthony"]
tags: ["wails","v2"]
date: "2022-02-22"
slug: "blog/wails-v2-beta-for-linux"
image: "/assets/blog-images/wails-linux.webp"
sourcePath: "blog/wails-v2-beta-for-linux.md"
---

![captura de tela do wails-linux](/assets/blog-images/wails-linux.webp)

Tenho o prazer de finalmente anunciar que o Wails v2 agora está em beta para Linux! É um tanto irônico que os primeiros experimentos com a v2 tenham sido realizados no Linux e, ainda assim, ela tenha acabado sendo a última versão lançada. Dito isso, a v2 que temos hoje é muito diferente daqueles primeiros experimentos. Portanto, sem mais delongas, vamos conhecer os novos recursos:

## Novos recursos

![captura de tela do wails-menus-linux](/assets/blog-images/wails-menus-linux.webp)

Recebemos muitos pedidos de suporte a menus nativos. Finalmente, o Wails oferece esse recurso. Agora há menus de aplicativo disponíveis, com suporte à maioria dos recursos de menus nativos. Isso inclui itens de menu padrão, caixas de seleção, grupos de opções, submenus e separadores.

Na v1, recebemos uma enorme quantidade de pedidos para oferecer maior controle sobre a própria janela. Tenho o prazer de anunciar que há novas APIs de runtime específicas para isso. Elas oferecem muitos recursos e são compatíveis com configurações de vários monitores. Também há uma API de caixas de diálogo aprimorada: agora você pode ter caixas de diálogo modernas e nativas, com amplas opções de configuração para atender a todas as suas necessidades.

### Não é necessário empacotar os recursos

Um grande problema da v1 era a necessidade de condensar todo o aplicativo em arquivos únicos de JS e CSS. Tenho o prazer de anunciar que, na v2, não é necessário empacotar os recursos de forma alguma. Quer carregar uma imagem local? Use uma tag `<../../../assets/blog-images>` com um caminho local no atributo src. Quer usar uma fonte bacana? Copie-a para o projeto e adicione o caminho correspondente ao CSS.

> Uau, isso parece um servidor web...

Sim, funciona como um servidor web, mas não é um.

> Então, como incluo meus recursos?

Basta passar um único `embed.FS` contendo todos os seus recursos para a configuração do aplicativo. Eles nem sequer precisam estar no diretório principal — o Wails cuidará disso para você.

### Nova experiência de desenvolvimento

Como os recursos não precisam mais ser empacotados, agora é possível ter uma experiência de desenvolvimento totalmente nova. O novo comando `wails dev` compila e executa o aplicativo, mas, em vez de usar os recursos no `embed.FS`, ele os carrega diretamente do disco.

Ele também oferece estes recursos adicionais:

- Recarregamento a quente — qualquer alteração nos recursos do frontend aciona o recarregamento automático do frontend do aplicativo
- Recompilação automática — qualquer alteração no código Go recompila e reinicia o aplicativo

Além disso, um servidor web será iniciado na porta 34115. Ele disponibilizará o aplicativo para qualquer navegador que se conectar. Todos os navegadores conectados responderão a eventos do sistema, como o recarregamento a quente quando um recurso for alterado.

Em Go, estamos acostumados a trabalhar com structs nos aplicativos. Muitas vezes, é útil enviar structs ao frontend e usá-las como estado no aplicativo. Na v1, esse processo era bastante manual e um tanto trabalhoso para quem desenvolvia. Tenho o prazer de anunciar que, na v2, qualquer aplicativo executado no modo de desenvolvimento gerará automaticamente modelos TypeScript para todas as structs usadas como parâmetros de entrada ou saída de métodos vinculados. Isso permite o intercâmbio transparente de modelos de dados entre os dois ambientes.

Além disso, outro módulo JS é gerado dinamicamente para encapsular todos os métodos vinculados. Ele fornece JSDoc para esses métodos, permitindo o preenchimento automático de código e a exibição de sugestões no IDE. É muito bacana ver os modelos de dados serem importados automaticamente ao pressionar Tab em um módulo gerado automaticamente que encapsula o código Go!

### Modelos remotos

![captura de tela do remote-linux](/assets/blog-images/remote-linux.webp)

Colocar um aplicativo em funcionamento rapidamente sempre foi um objetivo central do projeto Wails. Quando lançamos o projeto, tentamos abranger muitos dos frameworks modernos da época: react, vue e angular. O mundo do desenvolvimento frontend é cheio de opiniões fortes, evolui rapidamente e é difícil acompanhar! Como resultado, nossos modelos básicos ficavam desatualizados com muita rapidez, o que causava problemas de manutenção. Isso também significava que não tínhamos modelos modernos e interessantes para as tecnologias mais recentes e avançadas.

Com a v2, eu queria dar mais autonomia à comunidade, permitindo que vocês mesmos criassem e hospedassem modelos, em vez de depender do projeto Wails. Agora, portanto, você pode criar projetos usando modelos mantidos pela comunidade! Espero que isso inspire os desenvolvedores a criar um ecossistema dinâmico de modelos de projeto. Estou realmente entusiasmado para ver o que nossa comunidade de desenvolvedores será capaz de criar!

### Compilação cruzada para Windows

Como o Wails v2 para Windows é escrito inteiramente em Go, você pode gerar builds para Windows sem usar docker.

![captura de tela do build-cross-windows](/assets/blog-images/linux-build-cross-windows.webp)

### Conclusão

Como mencionei nas notas de versão do Windows, o Wails v2 representa uma nova base para o projeto. O objetivo desta versão é receber feedback sobre a nova abordagem e corrigir eventuais bugs antes do lançamento da versão final. Sua contribuição será muito bem-vinda! Envie seus comentários ao fórum de discussão da [v2 Beta](https://github.com/wailsapp/wails/discussions/828).

É **difícil** oferecer suporte ao Linux. Esperamos que a versão beta apresente algumas peculiaridades. Ajude-nos a ajudar você enviando relatórios de bugs detalhados!

Por fim, gostaria de fazer um agradecimento especial a todos os [patrocinadores do projeto](/credits/#sponsors), cujo apoio impulsiona o projeto de muitas maneiras nos bastidores.

Mal posso esperar para ver o que as pessoas criarão com o Wails nesta nova e empolgante fase do projeto!

Lea.

PS: o lançamento da v2 está próximo!

PPS: se você ou sua empresa considera o Wails útil, considere [patrocinar o projeto](https://github.com/sponsors/leaanthony). Agradeço!
