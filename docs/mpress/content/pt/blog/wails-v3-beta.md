---
title: "Wails v3 Beta: uma nova base para aplicações desktop em Go"
description: "O Wails v3 Beta apresenta um modelo de aplicação mais direto, bindings mais completos e uma base mais clara para aplicações desktop em Go."
authors: ["leaanthony"]
tags: ["wails","v3","beta"]
date: "2026-08-02"
slug: "blog/wails-v3-beta"
image: "/assets/screenshots/frameless-v3-native-corners-macos.png"
sourcePath: "blog/wails-v3-beta.md"
---

![Uma janela nativa sem moldura do Wails v3 no macOS](/assets/screenshots/frameless-v3-native-corners-macos.png)

Hoje estamos lançando o Wails v3 Beta.

O Wails permite que desenvolvedores Go criem aplicações desktop com as ferramentas de frontend web que já conhecem, usando a WebView nativa de cada plataforma em vez de um navegador incorporado. O v3 representa um avanço substancial: oferece às aplicações uma API mais direta, um modelo de build mais claro e uma base melhor para as aplicações desktop que as pessoas vêm pedindo que o Wails ofereça suporte.

Esta é uma versão beta, não a versão final do 3.0. A API para desktop está estável, e equipes já usam o v3 em produção, mas você deve testar tudo minuciosamente antes da implantação. Estamos usando o período beta para identificar, junto à comunidade, os últimos problemas de compatibilidade e fluxo de trabalho. O Wails v2 continua sendo a versão estável atual e continuará recebendo correções.

## Documentação durante o beta

Durante o beta, manteremos a documentação em inglês como fonte oficial enquanto a API e os fluxos de trabalho passam pela validação final. Não estamos aceitando PRs de tradução nesta etapa. O trabalho de tradução será retomado antes da disponibilidade geral, quando a documentação estiver estável o suficiente para que os tradutores trabalhem sem alterações repetidas.

## O que está incluído no v3

- Uma API explícita para aplicações e janelas, incluindo suporte de primeira classe a múltiplas janelas
- Serviços Go com análise estática do código-fonte que gera bindings TypeScript mais completos, preservando comentários e nomes de parâmetros significativos
- Serviços que podem agrupar recursos e scripts de frontend com sua API de backend — a base para plugins instaláveis e mais completos
- Um sistema de build visível, baseado em Taskfile, que você pode inspecionar, estender e depurar
- Builds de servidor para executar a mesma aplicação e os mesmos serviços sem uma janela desktop nativa
- Suporte a ambientes desktop modernos no macOS, Windows e Linux, em Intel, Apple Silicon, amd64 e arm64, quando compatíveis
- Suporte móvel experimental para iOS e Android, disponível para exploração, mas fora da garantia de compatibilidade do beta para desktop

## Por que o v3

O Wails v2 facilitou a criação de aplicações Go com um frontend web moderno. Ele atendeu bem ao projeto — e a muitas aplicações. No entanto, seu runtime de janela única orientado por contexto e seu processo de build rigidamente gerenciado tornavam algumas tarefas comuns no desenvolvimento desktop mais difíceis do que deveriam ser.

O v3 parte de um modelo diferente. Aplicações, janelas, serviços, eventos e recursos da plataforma são objetos explícitos. Isso facilita compreender o framework à medida que uma aplicação cresce e torna recursos como múltiplas janelas uma parte normal do modelo da aplicação, em vez de uma solução alternativa.

## Novidades

### Uma API de aplicação criada para software desktop de verdade

O v3 substitui o estilo de configuração v2 `wails.Run(...)` por um ciclo de vida explícito da aplicação. Você cria uma aplicação, registra serviços, cria janelas e interage com os objetos responsáveis pelo comportamento de que precisa.

Isso elimina grande parte da passagem implícita de contexto. As operações de janela pertencem às janelas; as operações que abrangem toda a aplicação pertencem à aplicação. É um modelo mais natural para aplicações com múltiplas janelas e mais adequado para testar e manter bases de código maiores.

O retorno das pessoas que usaram o v3 durante a fase alfa foi extremamente positivo. Em especial, os desenvolvedores receberam bem o modelo explícito: ele facilita acompanhar o código, deixa mais claras as responsabilidades e permite que aplicações desktop complexas cresçam sem entrar em conflito com o framework.

### Múltiplas janelas de primeira classe

Múltiplas janelas são um recurso central do v3. As janelas têm seu próprio ciclo de vida e podem ser criadas, gerenciadas e fechadas em tempo de execução. O resultado é um caminho mais claro para o tipo de software desktop que precisa de editores, inspetores, preferências, janelas de ferramentas ou vários elementos independentes de interface.

### Serviços e bindings gerados

Os serviços Go substituem o modelo anterior de bindings. Eles mantêm a lógica da aplicação como código Go convencional e tornam explícita a fronteira com o frontend. Os bindings são gerados em uma estrutura que reflete a aplicação e seus serviços, facilitando localizar e usar a API exposta ao frontend.

O v3 gera esses bindings por meio da análise estática do código-fonte. Isso significa que o gerador pode preservar as informações que os desenvolvedores inserem no código — incluindo comentários e nomes de parâmetros significativos — em vez de descobrir, por reflexão, um programa já compilado. O resultado é uma API de frontend mais completa e útil, além de um processo de geração mais fácil de compreender e manter.

Os serviços também podem montar recursos e scripts de frontend junto ao código Go. Isso reúne um recurso em um único lugar coerente: sua API de backend, o JavaScript ou a interface de que precisa e o ponto de integração com a aplicação host. Isso abre caminho para plugins do Wails que ofereçam funcionalidades completas prontas para uso — instale um plugin, monte seu serviço e use o recurso — em vez de você mesmo reunir um conjunto disperso de bindings e dependências de frontend. Um sistema geral de plugins não faz parte deste beta, mas o v3 torna essa direção viável de uma forma que o modelo de bindings do v2 não permitia.

### Um sistema de build que você pode inspecionar e adaptar

O v3 torna visível a estrutura de build do projeto. Em vez de ocultar todas as decisões de build em um único comando, os projetos têm um layout convencional e uma configuração de build baseada em Taskfile, que podem ser compreendidos, estendidos e depurados junto à aplicação.

### Uma base desktop multiplataforma mais sólida

O beta oferece suporte ao Windows em amd64 e arm64, ao macOS em Intel e Apple Silicon e ao Linux em amd64 e arm64. GTK4 com WebKitGTK 6.0 é a pilha padrão no Linux; GTK3 continuará disponível como opção legada durante toda a série v3.0. O suporte móvel é promissor, mas permanece experimental e não faz parte da garantia de compatibilidade do beta para desktop.

Esta versão também inclui o trabalho necessário para tornar a experiência cotidiana mais confiável: comportamento aprimorado nas plataformas, um modelo de janelas mais completo, diagnósticos mais claros e artefatos de versão com somas de verificação e proveniência.

## Migração do v2

O v3 é uma nova versão principal, e a migração é uma portabilidade de verdade, não apenas uma alteração no número da versão. As principais mudanças conceituais são o ciclo de vida da aplicação e das janelas, serviços em vez de bindings vinculados ao contexto, APIs diretas para aplicação e janelas em vez do pacote de runtime v2 e bindings de frontend gerados novamente.

Publicamos um [guia de migração do v2 para o v3](/migration/v2-to-v3/) que apresenta essas mudanças e inclui um mapeamento de recursos e uma lista de verificação para testes. Esse guia manual é o caminho de migração oficialmente compatível para este beta. Não espere que todos os projetos v2 sejam convertidos sem revisão: teste o resultado, faça a portabilidade das chamadas de runtime de forma deliberada e mantenha o v2 em uso até que a nova aplicação esteja pronta.

Também estamos avaliando um assistente de migração experimental. Ele não faz parte desta versão beta, e só o recomendaremos depois que for validado com projetos v2 reais e representativos.

## Uma observação franca sobre a trajetória

A primeira tag alfa do v3 foi publicada em 18 de janeiro de 2023. É muito tempo em fase alfa, e isso merece mais do que um reconhecimento vago.

O projeto mudou enormemente durante esse período. Quando o primeiro beta do Wails v2 foi lançado para Windows em setembro de 2021, o repositório tinha aproximadamente 4000 estrelas. Na versão v2, lançada em setembro de 2022, tinha aproximadamente 10300. Hoje tem mais de 35000. Esse crescimento é um privilégio, mas também muda o que significa cuidar bem do projeto: mais usuários dependem das decisões de lançamento, mais colaboradores precisam de caminhos claros para participar e uma parcela maior do trabalho consiste em tornar o projeto previsível, em vez de simplesmente adicionar o próximo recurso.

Nem sempre adaptei nossos processos com a rapidez ou a clareza que esse crescimento exigia. Assumo essa responsabilidade. A resposta não é fazer grandes promessas nem transformar cada decisão em uma cerimônia, mas deixar mais claro o que tem suporte, o que é experimental, como as decisões são tomadas e o que faremos a seguir.

Esse trabalho começa com esta versão beta. Agora temos marcos explícitos para as fases beta, candidata a lançamento e de disponibilidade geral (GA); compromissos de compatibilidade mais claros; uma política de segurança atualizada; e um processo de WEP (Wails Enhancement Proposal) para alterações no comportamento público e novos recursos. Continuaremos aprimorando o roteiro e revisando a governança do projeto à medida que o Wails crescer. O objetivo é ter um projeto no qual seja mais fácil confiar e para o qual seja mais fácil contribuir, não um projeto mais difícil de fazer avançar.

Também estamos relançando o [subreddit do Wails](https://www.reddit.com/r/wails/) e vamos mantê-lo ativamente como mais um espaço para discussões práticas, perguntas e feedback à medida que o v3 avança em direção à disponibilidade geral.

## Instale e experimente a versão beta

Depois que a versão for publicada, instale a CLI v3 mais recente com:

```sh
go install github.com/wailsapp/wails/v3/cmd/wails3@latest
```

Antes de criar um projeto, execute o assistente de configuração guiada. Ele verifica seu ambiente de desenvolvimento local e ajuda a configurar as dependências necessárias para o Wails:

```sh
wails3 setup
```

![O assistente de configuração do Wails v3 no modo escuro](/assets/screenshots/wails3-setup-wizard-dark.png)

### Crie um projeto

Quando tudo estiver pronto, crie um projeto com:

```sh
wails3 init
```

- Documentação: <https://v3.wails.io/>
- Guia de migração: <https://v3.wails.io/migration/v2-to-v3/>
- Reddit: <https://www.reddit.com/r/wails/>
- Notas de lançamento: [ESPAÇO RESERVADO PARA A VERSÃO](https://github.com/wailsapp/wails/releases)

Se você encontrar um bug reproduzível, informe-o com a saída de `wails3 doctor` e, quando possível, um exemplo mínimo. Se quiser propor um novo recurso ou uma alteração no comportamento público, abra um PR de rascunho de WEP em vez de uma issue de solicitação de recurso. Os dois caminhos nos ajudam a responder com clareza e a manter o avanço da versão beta.

## Agradecemos

O Wails v3 existe graças às pessoas que testaram builds incompletos, relataram bugs difíceis, traduziram a documentação, responderam a perguntas, contribuíram com código e continuaram cobrando melhorias do projeto. Agradecemos a todos.

E um agradecimento muito especial e sincero aos patrocinadores que sustentaram o Wails durante esta longa transição. O apoio de vocês fez mais do que manter o projeto funcionando: deu-nos a oportunidade de dedicar tempo contínuo à arquitetura, às ferramentas, aos testes e à documentação que permitiram ao projeto acelerar rumo ao v3. Cada pessoa que testou e contribuiu ajudou a moldar esta versão, mas foram os patrocinadores que tornaram possível dar a esse trabalho a atenção que ele merecia.

A versão beta é um convite para nos ajudar a concluir o v3 da maneira certa. Experimente-a, crie projetos com ela, conte-nos onde ela falha e ajude-nos a tornar o caminho até o 3.0 curto e cuidadoso.
