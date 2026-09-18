---
title: "Perguntas frequentes"
description: "Respostas às perguntas mais comuns sobre o desenvolvimento de aplicativos com o Wails v3"
slug: "faq"
sourcePath: "faq.md"
---

## Geral

### O que é o Wails?

O Wails é um framework para desenvolver aplicativos desktop com Go e tecnologias web. A lógica do aplicativo é escrita em Go, a interface é criada com HTML, CSS e JavaScript (ou qualquer framework de frontend), e o Wails a renderiza na webview nativa do sistema operacional. O resultado é um aplicativo pequeno, rápido e com aparência nativa: sem navegador incluído, com baixo uso de memória e um único arquivo binário que normalmente ocupa cerca de 10 MB.

### Quais plataformas são compatíveis com o Wails?

| Plataforma | Requisitos |
| --- | --- |
| Windows | AMD64 e ARM64. Usa o [runtime do WebView2](https://developer.microsoft.com/microsoft-edge/webview2/). |
| macOS | 10.15 ou posterior em processadores Intel (os aplicativos podem ter como destino a versão 10.13 ou posterior) e 11.0 ou posterior no Apple Silicon. Há suporte a binários universais. |
| Linux | AMD64 e ARM64. A pilha padrão é o GTK4 com WebKitGTK 6.0 (Ubuntu 24.04 ou posterior, Debian 13 ou posterior, Fedora 40 ou posterior e similares). Distribuições que oferecem apenas o WebKit2GTK 4.1, como Ubuntu 22.04, Debian 12 e RHEL 9, são compatíveis por meio da compilação legada `-tags gtk3` (disponível até a v3.1). Distribuições que oferecem apenas o WebKit2GTK 4.0 não são compatíveis. Consulte o [guia de compilação para Linux](/guides/build/linux/). |
| iOS e Android | Experimental. Consulte os [guias para dispositivos móveis](/guides/mobile/). |

Você também pode disponibilizar seu aplicativo como um aplicativo web comum usando a [compilação para servidor](/guides/server-build/).

Execute `wails3 doctor` a qualquer momento para verificar seu sistema e obter instruções de instalação específicas para a plataforma.

### Do que preciso para começar?

- Go 1.25 ou posterior
- Node.js e npm (para a compilação do frontend)
- Cadeia de ferramentas da plataforma: WebView2 no Windows (pré-instalado no 10/11), Xcode Command Line Tools no macOS e `gcc`, além dos pacotes de desenvolvimento do GTK/WebKit, no Linux

`wails3 doctor` verifica tudo isso para você e informa exatamente o que está faltando. Consulte [Instalação](/quick-start/installation/) para ver o passo a passo completo.

### O Wails v3 está pronto para produção?

O Wails v3 é um software beta com uma API estável para desktop. Já há aplicativos executando-o em produção, mas você deve realizar testes rigorosos antes da implantação enquanto concluímos os ajustes finais para 3.0. Consulte a [página de status do projeto](/status/) para conhecer a situação atual. O Wails v2 é a versão estável atual e continua recebendo correções.

## Desenvolvimento

### Preciso saber Go?

Conhecimentos básicos de Go ajudam, mas você não precisa ser especialista. A lógica do aplicativo fica em métodos Go comuns, e os [tutoriais](/tutorials/overview/) orientam você em todo o restante. Muitos desenvolvedores aprendem Go enquanto criam seu primeiro aplicativo com o Wails.

### Posso usar meu framework de frontend preferido?

Sim. Se ele gerar HTML, CSS e JavaScript, funcionará com o Wails. Há modelos para React, Vue, Svelte e JavaScript puro (cada um com variantes em TypeScript), e qualquer outra opção pode ser integrada em poucos minutos. Consulte [Frameworks de frontend](/guides/dev/frontend-frameworks/).

### Como chamo funções Go a partir do JavaScript?

Registre um serviço, e o Wails gerará bindings tipados para ele:

```go
// Go
type GreetService struct{}

func (g *GreetService) Greet(name string) string {
    return "Hello " + name
}
```

```javascript
// JavaScript
import { GreetService } from "./bindings/changeme";

const message = await GreetService.Greet("World");
```

Os bindings são gerados novamente de forma automática durante `wails3 dev` ou, sob demanda, com `wails3 generate bindings`. Consulte [Serviços](/features/bindings/services/).

### Posso usar TypeScript?

Sim. O gerador de bindings produz definições TypeScript para seus serviços e os respectivos tipos, portanto as chamadas ao Go são totalmente tipadas.

### Como envio eventos entre Go e JavaScript?

```go
// Go
app.Event.Emit("time", time.Now().Format(time.RFC1123))
```

```javascript
// JavaScript
import { Events } from "@wailsio/runtime";

Events.On("time", (event) => {
    console.log(event.data);
});
```

Os nomes dos eventos devem ser exatamente iguais. Consulte a [Referência de eventos](/guides/events-reference/).

### Como depuro meu aplicativo?

Execute `wails3 dev` e clique com o botão direito na janela para abrir as ferramentas de desenvolvimento do navegador, exatamente como faria na web. O servidor de desenvolvimento também oferece recarregamento automático do frontend. Consulte [Depuração](/guides/dev/debugging/).

## Compilação e distribuição

### Como faço uma compilação para produção?

```bash
wails3 build
```

O arquivo binário é gerado em `bin/`. As compilações de produção já aplicam configurações padrão adequadas (tags de compilação, `-trimpath` e remoção de símbolos), portanto não são necessárias flags adicionais para obter um binário enxuto.

### Posso fazer compilação cruzada?

Com algumas limitações. A compilação cruzada de Go puro não se aplica porque cada plataforma usa bibliotecas nativas de webview, mas há bom suporte para casos comuns:

```bash
# Different architecture, same OS
wails3 build GOOS=windows GOARCH=arm64

# macOS universal binary
wails3 task darwin:build:universal
```

A compilação para Linux a partir de outro sistema operacional usa uma cadeia de ferramentas baseada em Docker. Consulte [Compilações multiplataforma](/guides/build/cross-platform/) para ver a matriz completa.

### Como crio um instalador ou pacote?

```bash
wails3 package
```

Isso produz o formato nativo da plataforma, e o [guia de instaladores](/guides/installers/) aborda o NSIS no Windows, os pacotes `.app` e DMGs no macOS e os pacotes do Linux.

### Como assino digitalmente meu aplicativo?

A assinatura no Windows e no macOS, incluindo a notarização, é abordada passo a passo no [guia de assinatura](/guides/build/signing/).

## Recursos

### Posso criar várias janelas?

Sim, o suporte a várias janelas é nativo na v3:

```go
window1 := app.Window.New()
window2 := app.Window.New()
```

Consulte [Várias janelas](/features/windows/multiple/).

### O Wails oferece suporte à bandeja do sistema?

Sim, incluindo menus e manipuladores de cliques:

```go
systemTray := app.SystemTray.New()
systemTray.SetIcon(iconBytes)
systemTray.SetMenu(myMenu)
```

Consulte [Bandeja do sistema](/features/menus/systray/).

### Posso usar caixas de diálogo nativas?

Sim. As caixas de diálogo de arquivo, de mensagem e de pergunta usam implementações nativas:

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select File").
    PromptForSingleSelection()
```

Consulte [Caixas de diálogo](/features/dialogs/overview/).

### O Wails oferece suporte a atualizações automáticas?

Sim. O Wails v3 inclui um atualizador automático integrado (`app.Updater`), com provedores conectáveis para GitHub Releases, keygen.sh e Sparkle AppCast, verificação de assinaturas criptográficas e uma interface padrão que você pode personalizar com temas ou substituir. Consulte o guia [Atualizador no aplicativo](/guides/updater/) e o tutorial [Aplicativo Wails com atualização automática](/tutorials/04-self-update-a-wails-app/).

## Solução de problemas

### Algo não está funcionando. Por onde começo?

```bash
wails3 doctor
```

Ele verifica sua cadeia de ferramentas, lista as dependências ausentes com os comandos de instalação e exibe as informações de versão que você deve incluir em qualquer relatório de bug.

### Minha compilação falha

Estas são as soluções mais comuns, na ordem:

1. `go mod tidy`
2. `cd frontend && npm install` (a ausência de `node_modules` é a causa mais comum)
3. Atualize a CLI: `go install github.com/wailsapp/wails/v3/cmd/wails3@latest`
4. No Linux, verifique em `wails3 doctor` se estão faltando pacotes do GTK/WebKit

### Minhas vinculações estão ausentes ou desatualizadas

```bash
wails3 generate bindings
```

As vinculações são regeneradas automaticamente no modo de desenvolvimento. Se você adicionou um novo serviço ou alterou assinaturas de métodos fora de `wails3 dev`, regenere-as manualmente.

### Os eventos não estão sendo disparados

Os nomes dos eventos devem ser exatamente iguais entre `app.Event.Emit("name", ...)` no Go e `Events.On("name", ...)` no JavaScript. Primeiro, verifique se há erros de digitação e diferenças entre maiúsculas e minúsculas.

### Encontrei um bug

[Abra uma issue](https://github.com/wailsapp/wails/issues) e inclua a saída de `wails3 doctor`. O [guia de feedback](/feedback/) explica o que torna um relatório fácil de tratar.

## Migração da v2

### Devo migrar da v2 para a v3?

A v3 oferece suporte a várias janelas, uma API mais organizada baseada em serviços, um atualizador integrado, um sistema de compilação muito mais flexível e melhor desempenho. Novos projetos devem começar na v3. Para projetos existentes, o [Guia de migração](/migration/v2-to-v3/) apresenta as diferenças passo a passo.

### A v2 continuará recebendo manutenção?

Sim. A v2 continuará recebendo correções enquanto a v3 avança em direção à sua versão estável.

### Posso executar a v2 e a v3 lado a lado?

Sim. As CLIs são binários separados (`wails` e `wails3`), e os módulos têm caminhos de importação diferentes. Portanto, projetos em versões principais diferentes podem coexistir sem problemas na mesma máquina.

## Comunidade

### Como obtenho ajuda?

- [Discord](https://discord.gg/JDdSxwjhGf) para perguntas rápidas e discussões
- [GitHub Discussions](https://github.com/wailsapp/wails/discussions) para perguntas mais detalhadas
- [GitHub Issues](https://github.com/wailsapp/wails/issues) para bugs

### Como posso contribuir?

Consulte o [Guia de contribuição](/contributing/). Correções de bugs são bem-vindas a qualquer momento. Novas funcionalidades e alterações no comportamento público usam um PR de rascunho de [WEP (Proposta de Aprimoramento do Wails)](https://github.com/wailsapp/wails/blob/master/v3/wep/README.md). A discussão informal no Discord ou no GitHub Discussions é opcional.

### Onde posso encontrar exemplos?

O repositório inclui mais de 60 exemplos executáveis que abrangem janelas, caixas de diálogo, eventos, bandeja do sistema, serviços e muito mais: [v3/examples](https://github.com/wailsapp/wails/tree/master/v3/examples).

## Ainda tem dúvidas?

Pergunte no [Discord](https://discord.gg/JDdSxwjhGf) ou [abra uma discussão](https://github.com/wailsapp/wails/discussions).
