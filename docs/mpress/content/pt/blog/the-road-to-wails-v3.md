---
title: "O caminho até o Wails v3"
description: "Notas de versão e anúncios do Wails"
authors: ["leaanthony"]
tags: ["wails","v3"]
date: "2023-01-17"
slug: "blog/the-road-to-wails-v3"
image: "/assets/blog-images/multiwindow.webp"
sourcePath: "blog/the-road-to-wails-v3.md"
---

![captura de tela com várias janelas](/assets/blog-images/multiwindow.webp)

## Introdução

O Wails é um projeto que simplifica o desenvolvimento de aplicativos desktop multiplataforma com Go. Ele usa componentes webview nativos no frontend (não navegadores incorporados), levando a potência do sistema de interface de usuário mais popular do mundo ao Go sem deixar de ser leve.

A versão 2 foi lançada em 22 de setembro de 2022 e trouxe muitas melhorias, incluindo:

- Desenvolvimento em tempo real com o popular projeto Vite
- Recursos avançados para gerenciar janelas e criar menus
- Componente WebView2 da Microsoft
- Geração de modelos Typescript que espelham suas structs Go
- Criação de instalador NSIS
- Builds ofuscados

Atualmente, o Wails v2 oferece ferramentas avançadas para criar aplicativos desktop multiplataforma sofisticados.

Esta publicação busca analisar o estado atual do projeto e o que podemos melhorar daqui para frente.

## Onde estamos agora?

Tem sido incrível ver a popularidade do Wails crescer desde o lançamento da v2. A criatividade da comunidade e as coisas maravilhosas que estão sendo criadas com o Wails me surpreendem constantemente. Quanto maior a popularidade, mais atenção o projeto recebe. E, com isso, surgem mais solicitações de recursos e relatos de bugs.

Com o tempo, consegui identificar alguns dos problemas mais urgentes enfrentados pelo projeto. Também consegui identificar alguns dos fatores que estão limitando seu avanço.

## Problemas atuais

Identifiquei as seguintes áreas que, na minha opinião, estão limitando o avanço do projeto:

- A API
- Geração de bindings
- O sistema de build

### A API

Atualmente, a API para criar um aplicativo Wails consiste em 2 partes:

- A API de aplicativo
- A API de runtime

Como se sabe, a API de aplicativo tem apenas 1 função: `Run()`, que recebe inúmeras opções que determinam como o aplicativo funcionará. Embora seja muito simples de usar, ela também é bastante limitada. É uma abordagem "declarativa" que oculta grande parte da complexidade subjacente. Por exemplo, não há uma referência à janela principal, portanto você não pode interagir diretamente com ela. Para isso, é necessário usar a API de runtime. Isso se torna um problema quando você começa a querer fazer coisas mais complexas, como criar várias janelas.

A API de runtime oferece muitas funções utilitárias ao desenvolvedor. Isso inclui:

- Gerenciamento de janelas
- Caixas de diálogo
- Menus
- Eventos
- Logs

Há várias coisas na API de runtime que não me agradam. A primeira é que ela exige que um "contexto" seja repassado. Isso é frustrante e confuso para novos desenvolvedores, que fornecem um contexto e depois recebem um erro de runtime.

O maior problema da API de runtime é que ela foi projetada para aplicativos que usam apenas uma janela. Com o tempo, a demanda por várias janelas aumentou, e a API não é adequada para isso.

### Considerações sobre a API da v3

Não seria ótimo se pudéssemos fazer algo assim?

```go
func main() {
    app := wails.NewApplication(options.App{})
    myWindow := app.NewWindow(options.Window{})
    myWindow.SetTitle("My Window")
    myWindow.On(events.Window.Close, func() {
        app.Quit()
    })
    app.Run()
}
```

Essa abordagem programática é muito mais intuitiva e permite que o desenvolvedor interaja diretamente com os elementos do aplicativo. Todos os métodos atuais do runtime relacionados a janelas simplesmente se tornariam métodos do objeto de janela. Quanto aos outros métodos do runtime, poderíamos movê-los para o objeto de aplicativo, desta forma:

```go
app := wails.NewApplication(options.App{})
app.NewInfoDialog(options.InfoDialog{})
app.Log.Info("Hello World")
```

Essa API é muito mais poderosa e permitirá criar aplicativos mais complexos. Ela também permite criar várias janelas, [o recurso mais votado no GitHub](https://github.com/wailsapp/wails/issues/1480):

```go
func main() {
    app := wails.NewApplication(options.App{})
    myWindow := app.NewWindow(options.Window{})
    myWindow.SetTitle("My Window")
    myWindow.On(events.Window.Close, func() {
        app.Quit()
    })
    myWindow2 := app.NewWindow(options.Window{})
    myWindow2.SetTitle("My Window 2")
    myWindow2.On(events.Window.Close, func() {
        app.Quit()
    })
    app.Run()
}
```

### Geração de bindings

Um dos principais recursos do Wails é gerar bindings para seus métodos Go, permitindo que sejam chamados pelo Javascript. O método atual para fazer isso é uma espécie de gambiarra. Ele envolve compilar o aplicativo com uma flag especial e depois executar o binário resultante, que usa reflexão para determinar o que foi vinculado. Isso cria uma espécie de dilema do ovo e da galinha: não é possível compilar o aplicativo sem os bindings, nem gerar os bindings sem compilar o aplicativo. Há muitas maneiras de contornar isso, mas a melhor seria não usar essa abordagem.

Houve várias tentativas de criar um analisador estático para projetos Wails, mas elas não avançaram muito. Mais recentemente, isso se tornou um pouco mais fácil graças à maior disponibilidade de material sobre o assunto.

Em comparação com a reflexão, a abordagem baseada em AST é muito mais rápida, porém significativamente mais complexa. Inicialmente, talvez seja necessário impor certas restrições à forma de especificar os bindings no código. O objetivo é oferecer suporte aos casos de uso mais comuns e ampliá-lo posteriormente.

### O sistema de build

Assim como a abordagem declarativa da API, o sistema de build foi criado para ocultar as complexidades da compilação de um aplicativo desktop. Quando você executa `wails build`, ele realiza muitas tarefas nos bastidores:

- Compila o binário do backend para os bindings e gera os bindings
- Instala as dependências do frontend
- Compila os recursos do frontend
- Verifica se o ícone do aplicativo está presente e, em caso afirmativo, o incorpora
- Compila o binário final
- Se o build for para `darwin/universal`, compila 2 binários, um para `darwin/amd64` e outro para `darwin/arm64`, e então cria um binário universal usando `lipo`
- Se a compactação for necessária, compacta o binário com UPX
- Determina se esse binário deve ser empacotado e, em caso afirmativo:
  - Garante que o ícone e o manifesto do aplicativo sejam compilados no binário (Windows)
  - Monta o pacote do aplicativo, gera o pacote de ícones e copia esse pacote, o binário e o Info.plist para o pacote do aplicativo (Mac)

- Se um instalador NSIS for necessário, ele o cria

Todo esse processo, embora seja muito poderoso, também é muito opaco. É muito difícil personalizá-lo e depurá-lo.

Para resolver isso na v3, gostaria de adotar um sistema de build externo ao Wails. Depois de usar o [Task](https://taskfile.dev/) por algum tempo, tornei-me um grande entusiasta da ferramenta. Ela é excelente para configurar sistemas de build e deve ser razoavelmente familiar para quem já usou Makefiles.

O sistema de build seria configurado por meio de um arquivo `Taskfile.yml`, gerado por padrão com qualquer um dos templates compatíveis. Esse arquivo conteria todas as etapas necessárias para realizar todas as tarefas atuais, como compilar ou empacotar a aplicação, permitindo uma personalização fácil.

Não haverá nenhum requisito externo para essa ferramenta, pois ela fará parte da CLI do Wails. Isso significa que você ainda poderá usar `wails build`, que continuará fazendo tudo o que faz atualmente. No entanto, se quiser personalizar o processo de build, você poderá editar o arquivo `Taskfile.yml`. Isso também permitirá compreender facilmente as etapas do build e usar seu próprio sistema de build, se desejar.

A peça que falta no quebra-cabeça do build são as operações atômicas do processo, como geração de ícones, compactação e empacotamento. Exigir várias ferramentas externas não proporcionaria uma boa experiência ao desenvolvedor. Para resolver isso, a CLI do Wails oferecerá todos esses recursos como parte da própria CLI. Assim, os builds continuarão funcionando como esperado, sem ferramentas externas adicionais, mas você poderá substituir qualquer etapa do build pela ferramenta que preferir.

Esse sistema de build será muito mais transparente, facilitará a personalização e resolverá muitos dos problemas apontados a seu respeito.

## Os benefícios

Essas mudanças positivas trarão enormes benefícios ao projeto:

- A nova API será muito mais intuitiva e permitirá criar aplicações mais complexas.
- O uso de análise estática para gerar bindings será muito mais rápido e reduzirá grande parte da complexidade do processo atual.
- O uso de um sistema de build externo e consolidado tornará o processo de build totalmente transparente, permitindo uma personalização avançada.

Os benefícios para os mantenedores do projeto serão:

- A nova API será muito mais fácil de manter e adaptar a novos recursos e plataformas.
- O novo sistema de build será muito mais fácil de manter e ampliar. Espero que isso dê origem a um novo ecossistema de pipelines de build desenvolvidos pela comunidade.
- Uma melhor separação de responsabilidades dentro do projeto. Isso facilitará a adição de novos recursos e plataformas.

## O plano

Grande parte dos experimentos necessários já foi realizada, e os resultados parecem promissores. Ainda não há um cronograma para esse trabalho, mas espero que, até o final do primeiro trimestre de 2023, haja uma versão alfa para Mac, permitindo que a comunidade a teste, faça experimentos e envie feedback.

## Resumo

- A API da v2 é declarativa, oculta muita coisa do desenvolvedor e não é adequada para recursos como múltiplas janelas. Será criada uma nova API, mais simples, intuitiva e poderosa.
- O sistema de build é opaco e difícil de personalizar; por isso, adotaremos um sistema de build externo que tornará tudo transparente.
- A geração de bindings é lenta e complexa; por isso, adotaremos a análise estática, que eliminará grande parte da complexidade do método atual.

Muito trabalho foi dedicado aos componentes internos da v2, que estão sólidos. Agora é hora de aprimorar a camada construída sobre eles e proporcionar uma experiência muito melhor ao desenvolvedor.

Espero que você esteja tão empolgado com isso quanto eu. Aguardo seus comentários e seu feedback.

Atenciosamente,

&dash; Lea

PS: Se você ou sua empresa consideram o Wails útil, pense na possibilidade de [patrocinar o projeto](https://github.com/sponsors/leaanthony). Obrigado!

PPS: Sim, essa é uma captura de tela autêntica de uma aplicação com múltiplas janelas criada com o Wails. Não é um mockup. É real. É incrível. Em breve estará disponível.
