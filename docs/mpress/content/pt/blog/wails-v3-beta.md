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

O Wails permite que desenvolvedores Go criem aplicações desktop com as ferramentas de frontend web que já conhecem, usando a WebView nativa de cada plataforma em vez de um navegador incorporado. A v3 representa um avanço significativo: oferece às aplicações uma API mais direta, um modelo de compilação mais claro e uma base melhor para as aplicações desktop que as pessoas vêm pedindo que o Wails ofereça suporte.

Esta é uma versão beta, não a versão 3.0 final. A API para desktop está estável, e algumas equipes já usam a v3 em produção, mas você deve realizar testes minuciosos antes da implantação. Estamos usando o período beta para identificar, junto à comunidade, os últimos problemas de compatibilidade e fluxo de trabalho. O Wails v2 continua sendo a versão estável atual e continuará recebendo correções.

## Documentação durante a fase beta

Durante a fase beta, manteremos a documentação em inglês como fonte oficial enquanto a API e os fluxos de trabalho passam pela validação final. Não estamos aceitando PRs de tradução nesta etapa. O trabalho de tradução será retomado antes da disponibilidade geral, quando a documentação estiver estável o suficiente para que os tradutores trabalhem sem alterações recorrentes.

## O que está incluído na v3

- Uma API explícita para aplicações e janelas, incluindo suporte de primeira classe a múltiplas janelas
- Serviços Go com análise estática do código-fonte que gera bindings TypeScript mais completos, preservando comentários e nomes de parâmetros significativos
- Serviços capazes de agrupar recursos e scripts de frontend com sua API de backend — a base para plugins instaláveis e mais completos
- Um sistema de compilação visível, baseado em Taskfile, que você pode inspecionar, estender e depurar
- Compilações para servidor que permitem executar a mesma aplicação e os mesmos serviços sem uma janela desktop nativa
- Suporte moderno a desktop para macOS, Windows e Linux em Intel, Apple Silicon, amd64 e arm64, quando compatíveis
- Suporte móvel experimental para iOS e Android, disponível para exploração, mas fora da garantia de compatibilidade da versão beta para desktop

## Por que a v3

O Wails v2 facilitou a criação de aplicações Go com um frontend web moderno. Ele atendeu bem ao projeto — e a muitas aplicações. No entanto, seu runtime orientado por contexto e limitado a uma única janela, além do processo de compilação rigidamente gerenciado, tornava algumas tarefas comuns em aplicações desktop mais difíceis do que deveriam ser.

A v3 parte de um modelo diferente. Aplicações, janelas, serviços, eventos e recursos das plataformas são objetos explícitos. Isso torna o framework mais fácil de compreender à medida que uma aplicação cresce e faz com que recursos como múltiplas janelas sejam uma parte normal do modelo da aplicação, e não uma solução alternativa.

## Novidades

### Uma API de aplicação criada para software desktop de verdade

A v3 substitui o estilo de configuração `wails.Run(...)` da v2 por um ciclo de vida explícito da aplicação. Você cria uma aplicação, registra serviços, cria janelas e interage com os objetos que controlam o comportamento necessário.

Isso elimina grande parte da passagem implícita de contexto. As operações de janela pertencem às janelas; as operações que abrangem toda a aplicação pertencem à aplicação. É um modelo mais natural para aplicações com múltiplas janelas e mais adequado para testar e manter bases de código maiores.

O retorno das pessoas que usaram a v3 durante a fase alfa foi extremamente positivo. Em especial, os desenvolvedores receberam bem o modelo explícito: ele facilita o acompanhamento do código, deixa mais claro quem é responsável por cada elemento e permite que aplicações desktop complexas cresçam sem entrar em conflito com o framework.

### Múltiplas janelas como recurso de primeira classe

Múltiplas janelas são um recurso central da v3. As janelas têm ciclo de vida próprio e podem ser criadas, gerenciadas e fechadas em tempo de execução. O resultado é um caminho mais claro para o tipo de software desktop que precisa de editores, inspetores, preferências, janelas de ferramentas ou vários elementos independentes de interface.

### Serviços e bindings gerados

Os serviços Go substituem o antigo modelo de bindings. Eles mantêm a lógica da aplicação como código Go convencional e tornam explícita a fronteira com o frontend. Os bindings são gerados em uma estrutura que reflete a aplicação e seus serviços, facilitando a localização e o uso da API exposta ao frontend.

A v3 gera esses bindings por meio da análise estática do código-fonte. Isso significa que o gerador pode preservar as informações que os desenvolvedores inserem no código — incluindo comentários e nomes de parâmetros significativos — em vez de usar reflexão para descobrir um programa já compilado. O resultado é uma API de frontend mais completa e útil, além de um processo de geração mais fácil de compreender e manter.

Os serviços também podem disponibilizar recursos e scripts de frontend junto ao código Go. Isso reúne de forma coerente tudo o que compõe um recurso: sua API de backend, o JavaScript ou a interface de que ele precisa e o ponto de integração com a aplicação hospedeira. Essa estrutura abre caminho para plugins do Wails que ofereçam recursos avançados prontos para uso — instale um plugin, disponibilize seu serviço e use o recurso — em vez de montar por conta própria uma coleção desconexa de bindings e dependências de frontend. Um sistema geral de plugins não faz parte desta versão beta, mas a v3 torna esse caminho viável de uma forma que o modelo de bindings da v2 não permitia.

### Um sistema de compilação que você pode inspecionar e adaptar

A v3 torna visível a estrutura de compilação do projeto. Em vez de ocultar todas as decisões de compilação em um único comando, os projetos têm um layout convencional e uma configuração de compilação baseada em Taskfile, que podem ser compreendidos, estendidos e depurados junto com a aplicação.

### Uma base desktop multiplataforma mais sólida

A versão beta oferece suporte ao Windows em amd64 e arm64, ao macOS em Intel e Apple Silicon e ao Linux em amd64 e arm64. GTK4 com WebKitGTK 6.0 é a pilha padrão no Linux; GTK3 continuará disponível como opção legada durante toda a série v3.0. O suporte móvel é promissor, mas continua experimental e não faz parte da garantia de compatibilidade da versão beta para desktop.

Esta versão também inclui o trabalho necessário para tornar a experiência cotidiana mais confiável: comportamento aprimorado nas plataformas, um modelo de janelas mais completo, diagnósticos mais claros e artefatos de versão com somas de verificação e informações de proveniência.

## Migração da v2

A v3 é uma nova versão principal, e a migração exige uma portabilidade real, não apenas a alteração do número da versão. As principais mudanças conceituais são os ciclos de vida da aplicação e das janelas, serviços em vez de bindings vinculados ao contexto, APIs diretas para aplicações e janelas em vez do pacote de runtime da v2 e bindings de frontend gerados novamente.

Publicamos um [guia de migração da v2 para a v3](/migration/v2-to-v3/) que explica essas mudanças e inclui um mapeamento de recursos e uma lista de verificação para testes. Esse guia manual é o caminho de migração com suporte para esta versão beta. Não espere que todos os projetos v2 sejam convertidos sem revisão: teste o resultado, faça a portabilidade das chamadas ao runtime de forma deliberada e mantenha a v2 em uso até que a nova aplicação esteja pronta.

Também estamos avaliando um assistente experimental de migração. Ele não faz parte desta versão beta, e só o recomendaremos depois que tiver sido validado com projetos v2 reais e representativos.

## Uma observação franca sobre a jornada

A primeira tag alfa da v3 foi publicada em 18 de janeiro de 2023. É muito tempo em fase alfa, e isso merece mais do que um reconhecimento vago.

O projeto mudou enormemente durante esse período. Quando a primeira versão beta do Wails v2 para Windows foi lançada em setembro de 2021, o repositório tinha aproximadamente 4000 estrelas. Na época do lançamento da v2, em setembro de 2022, esse número era de aproximadamente 10300. Hoje, ele ultrapassa 35000. Esse crescimento é um privilégio, mas também muda o que significa administrar bem o projeto: mais usuários dependem das decisões sobre lançamentos, mais colaboradores precisam de caminhos claros para participar, e uma parcela maior do trabalho consiste em tornar o projeto previsível, em vez de simplesmente adicionar o próximo recurso.

Nem sempre adaptei nossos processos com a rapidez ou a clareza que esse crescimento exigia. Assumo essa responsabilidade. A resposta não é fazer grandes promessas nem transformar cada decisão em uma cerimônia; é deixar mais claro o que tem suporte, o que é experimental, como as decisões são tomadas e o que faremos a seguir.

Esse trabalho começa com esta versão beta. Agora temos marcos explícitos para as fases beta, candidata a lançamento e GA; compromissos de compatibilidade mais claros; uma política de segurança atualizada; e um processo de WEP (Wails Enhancement Proposal) para alterações no comportamento público e novos recursos. Continuaremos aprimorando o roteiro e analisando a governança do projeto à medida que o Wails crescer. O objetivo é ter um projeto no qual seja mais fácil confiar e para o qual seja mais fácil contribuir — não um projeto mais difícil de fazer avançar.

Também estamos relançando o [subreddit do Wails](https://www.reddit.com/r/wails/) e vamos mantê-lo ativamente como mais um espaço para discussões práticas, perguntas e feedback à medida que a v3 avançar rumo à disponibilidade geral.

## Instale e experimente a versão beta

Depois que a versão for publicada, instale a CLI mais recente da v3 com:

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

Se você encontrar um bug reproduzível, informe-o incluindo a saída de `wails3 doctor` e, sempre que possível, um exemplo mínimo. Se quiser propor um novo recurso ou uma alteração no comportamento público, abra um PR de rascunho de WEP em vez de uma issue de solicitação de recurso. Os dois caminhos nos ajudam a responder com clareza e a manter o avanço da versão beta.

## Agradecimentos

O Wails v3 existe graças às pessoas que testaram builds incompletos, relataram bugs difíceis, traduziram a documentação, responderam a perguntas, contribuíram com código e continuaram incentivando o projeto a melhorar. Agradecemos a todos.

E um agradecimento muito especial e sincero aos patrocinadores que sustentaram o Wails durante essa longa transição. O apoio de vocês fez mais do que manter o projeto funcionando: deu-nos a oportunidade de dedicar tempo contínuo à arquitetura, às ferramentas, aos testes e à documentação que permitiram acelerar o projeto rumo à v3. Cada pessoa que testou e contribuiu ajudou a moldar esta versão, mas foram os patrocinadores que tornaram possível dedicar a esse trabalho a atenção que ele merecia.

A versão beta é um convite para nos ajudar a concluir a v3 da maneira certa. Experimente-a, desenvolva com ela, conte-nos onde ela falha e ajude-nos a tornar curto e cuidadoso o caminho até 3.0.
