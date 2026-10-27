---
title: "Wails v3.0.0 lançado"
description: "O Wails v3.0.0 chegou: várias janelas, bandeja do sistema, menus nativos, serviços e bindings gerados, um atualizador, suporte a dispositivos móveis e um novo site de documentação."
authors: ["leaanthony"]
tags: ["wails","v3","release"]
date: "2026-10-26"
slug: "blog/wails-v3-released"
image: "/assets/videos/wails-v3-release-clean-poster.webp"
sourcePath: "blog/wails-v3-released.md"
---

@video{src="/assets/videos/wails-v3-release-clean-av1.mp4" poster="/assets/videos/wails-v3-release-clean-poster.webp" title="Vídeo de lançamento do Wails v3.0.0 — melhor visualizado em tela cheia" width="1920" height="1080"}

## Chegou!

Hoje estou lançando o Wails v3.0.0.

Vou ser sincero: foi uma longa jornada. A primeira tag alfa do v3 foi lançada em janeiro de 2023 e, desde então, muitas pessoas criaram, testaram, quebraram e corrigiram coisas ao longo do caminho. Sou verdadeiramente grato a cada um de vocês. Se quiser conhecer a versão franca dessa jornada, escrevi sobre ela no [v3 anúncio da versão beta](/blog/wails-v3-beta/).

## O que *é* o Wails?

Se você chegou agora: o Wails permite que programadores Go criem aplicações para desktop com as ferramentas de frontend web que já conhecem, usando a WebView nativa de cada plataforma em vez de um navegador incorporado. Você escreve o backend em Go e a interface em HTML, CSS e JavaScript (ou no framework de sua preferência), e o Wails reúne tudo em uma única aplicação nativa.

## Por que o v3?

O Wails v2 facilitou a criação de aplicações Go com um frontend web moderno e atendeu muito bem a muitas aplicações. No entanto, seu runtime de janela única, orientado por contexto, tornou algumas tarefas rotineiras de desktop mais difíceis do que deveriam ser.

O v3 parte de um modelo diferente. Aplicações, janelas, serviços, eventos e recursos da plataforma são objetos explícitos. Você cria uma aplicação, registra serviços, cria janelas e interage com o objeto responsável pelo comportamento de que precisa. Isso fica mais fácil de acompanhar à medida que a aplicação cresce, e recursos como várias janelas são uma parte normal do modelo, não uma solução improvisada.

## O que há no v3.0.0

- **Várias janelas.** Crie, gerencie e feche janelas em tempo de execução, cada uma com sua própria configuração e seu próprio ciclo de vida.
- **Bandeja do sistema.** Dê à sua aplicação uma presença persistente, com menus, ícones para temas claro e escuro, ícones de modelo no macOS e a opção de anexar uma janela ao ícone da bandeja.
- **Menus e caixas de diálogo nativos.** Menus da aplicação, menus de contexto e caixas de diálogo nativas em cada plataforma.
- **Serviços e bindings gerados.** Seu backend continua sendo código Go comum. Os bindings são gerados por meio de análise estática do código-fonte, de modo que comentários e nomes de parâmetros significativos chegam ao frontend com um único comando `wails3 generate bindings`.
- **Um sistema de build que você pode inspecionar.** Os projetos usam um processo de build baseado em [Taskfile](https://taskfile.dev) que você pode inspecionar, estender e depurar, em vez de um único comando opaco.
- **Multiplataforma.** Windows (amd64 e arm64), macOS (Intel e Apple Silicon) e Linux (amd64 e arm64), com GTK4 e WebKitGTK 6.0 como padrão no Linux. A pilha GTK3 legada continua disponível com `-tags gtk3` até v3.1.
- **iOS e Android.** O mesmo `main.go` e o mesmo frontend podem ser executados em dispositivos móveis. O suporte a dispositivos móveis ainda está marcado como experimental na documentação; portanto, considere-o algo a explorar e sobre o qual enviar feedback, não uma promessa de estabilidade.
- **Criado com Go 1.26.** O módulo v3 é compilado com Go 1.26.
- **Um atualizador integrado à aplicação.** Provedores conectáveis, verificação criptográfica, substituição atômica e uma interface padrão que você pode personalizar com um tema ou substituir. Consulte o [guia do atualizador](/guides/updater/).
- **Empacotamento e assinatura.** Guias sobre [assinatura de código](/guides/build/signing/) em todas as plataformas, [instaladores](/guides/installers/) e [empacotamento MSIX](/guides/build/msix/).
- **Builds para servidor.** Execute a mesma aplicação e os mesmos serviços como um servidor HTTP, sem uma janela nativa.

## Destaques desde a versão beta

A publicação sobre a versão beta saiu em 2 de agosto, e o changelog esteve movimentado desde então. Quero destacar alguns pontos:

- **Um novo site de documentação.** Agora a documentação é executada no [M-Press](https://m-press.me), nossa própria ferramenta de documentação, e as traduções foram concluídas e migradas da comunidade. A documentação em inglês continua sendo a fonte oficial.
- **Uma experiência mais robusta no Windows.** Agora o Wails se recupera de falhas de processos do WebView2 em vez de deixar uma janela em branco, e vários problemas de ciclo de vida do COM no WebView2 também foram corrigidos.
- **Atalhos globais que respeitam seu teclado.** No macOS, os aceleradores de letras agora seguem a tecla identificada em layouts como AZERTY, QWERTZ e Dvorak.
- **Maior robustez do atualizador.** Entre outras melhorias, foram feitas correções para a substituição de binários entre sistemas de arquivos no Linux e para atualizações de AppImage.
- **Configuração e ferramentas.** O assistente `wails3 setup` foi concluído, e agora `wails3 dev` encerra seus processos em segundo plano quando a aplicação é fechada.
- **Liquid Glass nativo no macOS.** Suporte ao efeito Liquid Glass no macOS 26 e posteriores, com uma alternativa para versões anteriores.

## Migração do v2

O v3 é uma nova versão principal, e a migração é uma portabilidade de fato, não uma simples atualização de versão. As grandes mudanças são o ciclo de vida da aplicação e das janelas, serviços no lugar de bindings vinculados ao contexto, APIs diretas de aplicação e janela no lugar do pacote de runtime v2 e bindings de frontend gerados novamente.

O [v2guia de migração do v3](/migration/v2-to-v3/) explica tudo isso passo a passo, com exemplos de antes e depois. Faça a portabilidade de forma deliberada, teste durante o processo e mantenha o v2 em uso até que sua aplicação v3 esteja pronta.

**O Wails v2 não vai desaparecer.** Ele continuará recebendo atualizações críticas. Para novos projetos, use o v3.0.0. As duas CLIs são binários separados (`wails` e `wails3`), com caminhos de importação de módulos diferentes; portanto, você pode manter ambas na mesma máquina.

## Instalação

Instale a CLI do v3.0.0 com:

```sh
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0
```

Em seguida, execute o assistente de configuração guiada, que verifica seu ambiente e ajuda com as dependências necessárias ao Wails:

```sh
wails3 setup
```

![O assistente de configuração do Wails v3 no modo escuro](/assets/screenshots/wails3-setup-wizard-dark.png)

E crie seu primeiro projeto:

```sh
wails3 init
```

- Documentação: <https://v3.wails.io/>
- Guia de migração: <https://v3.wails.io/migration/v2-to-v3/>
- Notas de lançamento: [Wails v3.0.0 no GitHub](https://github.com/wailsapp/wails/releases/tag/v3.0.0)

## Suporte e lançamentos

O Wails cresceu muito além do que eu esperava, e tem sido difícil acompanhar o volume de issues e pull requests. Às vezes, os lançamentos ocorreram de forma improvisada, e correções e contribuições demoraram demais para chegar aos usuários. Implementamos sistemas para automatizar uma parte maior do trabalho preparatório de triagem e preparação de lançamentos, para podermos responder mais rapidamente aos problemas e dedicar mais tempo à entrega de melhorias e novos recursos.

Depois do v3.0.0, passaremos a publicar **pré-lançamentos noturnos** com tags datadas, como `v3.1.0-nightly.20261026`, juntamente com **lançamentos secundários semanais**, como `v3.1.0` e `v3.2.0`. Os lançamentos noturnos disponibilizam correções recentes mais cedo, enquanto os lançamentos semanais oferecem uma cadência regular de atualizações testadas.

Se encontrar um bug reproduzível, relate-o com a saída de `wails3 doctor` e um exemplo mínimo. Se quiser propor um novo recurso ou uma alteração no comportamento público, abra um PR de rascunho de WEP (Wails Enhancement Proposal), em vez de uma issue de solicitação de recurso.

## Traduções

A todos que ajudaram a traduzir a documentação do Wails v2: nosso mais sincero agradecimento. Vocês dedicaram seu tempo para que desenvolvedores pudessem ler a documentação no próprio idioma, e isso fez a diferença para muitas pessoas. Somos profundamente gratos.

Achamos o Crowdin particularmente difícil de gerenciar e deixamos de usá-lo. A partir da versão v3, a documentação é traduzida automaticamente, e as traduções ficam no repositório junto às páginas em inglês.

Se você quiser uma tradução para um idioma que ainda não oferecemos ou encontrar um erro em uma das traduções disponíveis, abra uma issue e nos avise. Teremos muito prazer em ouvir você e agradecemos por cada relato.

## Créditos e agradecimentos

O Wails v3 existe graças às pessoas que testaram builds incompletos, relataram bugs difíceis, traduziram a documentação, responderam a perguntas e contribuíram com código. Agradecemos a todos. A lista completa de colaboradores está na [página de créditos](/credits/), e recomendo que você a consulte.

Em primeiro lugar, nosso imenso agradecimento a [Atterpac](https://github.com/atterpac), responsável pela manutenção do Wails, pelas ferramentas de desenvolvimento, pela integração com desktops Linux e pelo suporte à comunidade que ajudaram a levar a versão v3 até este ponto.

Agradecimento especial a [Fabio Massaioli (fbbdev)](https://github.com/fbbdev) por criar um analisador estático de código-fonte e um gerador de bindings tão incríveis. Eles são fundamentais para o funcionamento da versão v3: serviços Go comuns se tornam APIs JavaScript e TypeScript, com os tipos e a documentação disponibilizados no frontend.

Agradecemos também a [Travis McLane](https://github.com/tmclane) pelo extenso trabalho inicial nas plataformas; a [Stffabi](https://github.com/stffabi) pelo servidor de assets e pelas bases das plataformas nativas; a [Zach Botterman (popaprozac)](https://github.com/popaprozac) pelas notificações, pela integração com o Dock e por outros recursos de desktop; a [Ian VanSchooten (IanVS)](https://github.com/IanVS) pelos eventos tipados e pelas melhorias no runtime; e a [Lyimmi](https://github.com/Lyimmi) pelo trabalho no Linux. Essas contribuições continuam fazendo parte da versão v3, junto ao trabalho de todos que ajudaram a desenvolvê-la e aprimorá-la.

Nosso mais sincero agradecimento aos patrocinadores que apoiaram o Wails durante esta longa transição. O apoio de vocês nos deu tempo para trabalhar na arquitetura, nas ferramentas, nos testes e na documentação que nos trouxeram até aqui. Agradecemos especialmente também à JetBrains e à ZSA, bem como ao [CodeRabbit](https://www.coderabbit.ai/oss), por apoiar o projeto com revisões de código.

Por fim, um agradecimento especial a Valerio Tesei por todo o apoio que deu ao projeto ao longo dos anos.

## Considerações finais

Já temos muito trabalho nos esperando. Há novidades empolgantes surgindo no mundo da tecnologia, e estamos prontos para explorar novos paradigmas e trazer recursos completamente novos ao Wails. Antes disso, precisamos que a versão v3 seja uma base estável na qual você possa confiar. Fazer isso direito vem antes de lançar o que está por vir.

Tenho muito orgulho da versão v3. Espero que desenvolvê-la proporcione a você tanto prazer quanto tivemos ao criá-la. Experimente, crie algo e conte-nos o que achou. Venha nos dar um alô no [Discord](https://discord.gg/JDdSxwjhGf), no [GitHub Discussions](https://github.com/wailsapp/wails/discussions) ou no [subreddit](https://www.reddit.com/r/wails/).

Aproveite!

&dash;&nbsp;Lea

PS: Se você ou sua empresa consideram o Wails útil, pensem em [patrocinar o projeto](https://github.com/sponsors/leaanthony). Agradecemos!
