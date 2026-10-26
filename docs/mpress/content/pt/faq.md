---
title: "Perguntas frequentes"
description: "Respostas às perguntas mais comuns sobre o desenvolvimento de aplicativos com Wails v3"
slug: "faq"
sourcePath: "faq.md"
---

## Geral

### O que é o Wails?

O Wails é um framework para desenvolver aplicativos para desktop com Go e tecnologias web. A lógica do aplicativo é escrita em Go, a interface é criada com HTML, CSS e JavaScript (ou qualquer framework de frontend), e o Wails a renderiza na webview nativa do sistema operacional. O resultado é um aplicativo pequeno, rápido e com aparência nativa: sem navegador incorporado, com baixo uso de memória e um único binário que normalmente ocupa cerca de 10 MB.

### Quais plataformas são compatíveis com o Wails?

| Plataforma | Requisitos |
| --- | --- |
| Windows | AMD64 e ARM64. Usa o runtime do [WebView2](https://developer.microsoft.com/microsoft-edge/webview2/). |
| macOS | 10.15 ou posterior em processadores Intel (os aplicativos podem ter como destino a versão 10.13 ou posterior) e 11.0 ou posterior em Apple Silicon. Há suporte a binários universais. |
| Linux | AMD64 e ARM64. A pilha padrão é GTK4 com WebKitGTK 6.0 (Ubuntu 24.04 ou posterior, Debian 13 ou posterior, Fedora 40 ou posterior e semelhantes). Distribuições que fornecem apenas WebKit2GTK 4.1, como Ubuntu 22.04, Debian 12 e RHEL 9, são compatíveis por meio da compilação legada `-tags gtk3` (disponível até v3.1). Distribuições que fornecem apenas WebKit2GTK 4.0 não são compatíveis. Consulte o [guia de compilação para Linux](/guides/build/linux/). |
| iOS e Android | Experimental. Consulte os [guias para dispositivos móveis](/guides/mobile/). |

Você também pode disponibilizar o aplicativo como um aplicativo web convencional usando a [compilação para servidor](/guides/server-build/).

Execute `wails3 doctor` a qualquer momento para verificar o sistema e obter instruções de instalação específicas para a plataforma.

### Do que preciso para começar?

- Go 1.26 ou posterior
- Node.js e npm (para a compilação do frontend)
- Cadeia de ferramentas da plataforma: WebView2 no Windows (pré-instalado no 10/11), Xcode Command Line Tools no macOS e `gcc`, além dos pacotes de desenvolvimento do GTK/WebKit, no Linux

`wails3 doctor` verifica tudo isso e informa exatamente o que está faltando. Consulte [Instalação](/quick-start/installation/) para ver o passo a passo completo.

### O Wails v3 está pronto para produção?

O Wails v3.0.0 é a versão estável para aplicativos de desktop no Windows, macOS e Linux. O suporte a Android e iOS continua experimental e não faz parte do contrato de compatibilidade estável para desktop. Consulte a [página de status do projeto](/status/).

## Desenvolvimento

### Preciso saber Go?

Ter conhecimentos básicos de Go ajuda, mas você não precisa ser especialista. A lógica do aplicativo fica em métodos Go comuns, e os [tutoriais](/tutorials/overview/) orientam você em todo o restante. Muitos desenvolvedores aprendem Go enquanto criam seu primeiro aplicativo com Wails.

### Posso usar meu framework de frontend favorito?

Sim. Se ele gera HTML, CSS e JavaScript, funciona com o Wails. Há modelos para React, Vue, Svelte e JavaScript puro (cada um com variantes em TypeScript), e qualquer outra opção pode ser integrada em poucos minutos. Consulte [Frameworks de frontend](/guides/dev/frontend-frameworks/).

### Como chamo funções Go pelo JavaScript?

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

Os nomes dos eventos devem corresponder exatamente. Consulte a [Referência de eventos](/guides/events-reference/).

### Como depuro meu aplicativo?

Execute `wails3 dev` e clique com o botão direito na janela para abrir as ferramentas de desenvolvimento do navegador, exatamente como faria na web. O servidor de desenvolvimento também permite o recarregamento automático do frontend. Consulte [Depuração](/guides/dev/debugging/).

## Compilação e distribuição

### Como faço uma compilação para produção?

```bash
wails3 build
```

O binário é gerado em `bin/`. As compilações de produção já aplicam configurações padrão adequadas (tags de compilação, `-trimpath` e remoção de símbolos), portanto não são necessários parâmetros adicionais para obter um binário enxuto.

### Posso fazer compilação cruzada?

Com algumas limitações. A compilação cruzada de Go puro não se aplica, pois cada plataforma usa bibliotecas nativas de webview, mas há um bom suporte para os casos mais comuns:

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

Isso gera o formato nativo da plataforma, e o [guia de instaladores](/guides/installers/) aborda o NSIS no Windows, pacotes `.app` e DMGs no macOS, além de pacotes para Linux.

### Como assino digitalmente o código do meu aplicativo?

A assinatura no Windows e no macOS (inclusive a notarização) é explicada passo a passo no [guia de assinatura](/guides/build/signing/).

## Recursos

### Posso criar várias janelas?

Sim, v3 oferece suporte nativo a várias janelas:

```go
window1 := app.Window.New()
window2 := app.Window.New()
```

Consulte [Várias janelas](/features/windows/multiple/).

### O Wails oferece suporte à bandeja do sistema?

Sim, inclusive a menus e manipuladores de clique:

```go
systemTray := app.SystemTray.New()
systemTray.SetIcon(iconBytes)
systemTray.SetMenu(myMenu)
```

Consulte [Bandeja do sistema](/features/menus/systray/).

### Posso usar caixas de diálogo nativas?

Sim. As caixas de diálogo de arquivos, mensagens e perguntas usam implementações nativas:

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select File").
    PromptForSingleSelection()
```

Consulte [Caixas de diálogo](/features/dialogs/overview/).

### O Wails oferece atualizações automáticas?

Sim. O Wails v3 inclui um atualizador automático integrado (`app.Updater`), com provedores conectáveis para GitHub Releases, keygen.sh e Sparkle AppCast, verificação de assinatura criptográfica e uma interface de usuário padrão que você pode personalizar com um tema ou substituir. Consulte o guia [Atualizador integrado ao aplicativo](/guides/updater/) e o tutorial [Aplicativo Wails com atualização automática](/tutorials/04-self-update-a-wails-app/).

## Solução de problemas

### Algo não está funcionando. Por onde começo?

```bash
wails3 doctor
```

Ele verifica sua cadeia de ferramentas, lista as dependências ausentes com os comandos de instalação e exibe as informações de versão que você deve incluir em qualquer relatório de bug.

### Minha compilação falha

Estas são as correções habituais, na ordem:

1. `go mod tidy`
2. `cd frontend && npm install` (a ausência de `node_modules` é a causa mais comum)
3. Atualize a CLI: `go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0`
4. No Linux, verifique em `wails3 doctor` se faltam pacotes do GTK/WebKit

### Meus bindings estão ausentes ou desatualizados

```bash
wails3 generate bindings
```

Os bindings são gerados novamente de forma automática no modo de desenvolvimento; se você adicionou um novo serviço ou alterou assinaturas de métodos fora de `wails3 dev`, gere-os novamente de forma manual.

### Os eventos não estão sendo disparados

Os nomes dos eventos devem ser exatamente iguais entre `app.Event.Emit("name", ...)` no Go e `Events.On("name", ...)` no JavaScript. Verifique primeiro se há erros de digitação e diferenças entre maiúsculas e minúsculas.

### Encontrei um bug

[Abra uma issue](https://github.com/wailsapp/wails/issues) e inclua a saída de `wails3 doctor`. O [guia de feedback](/feedback/) explica o que torna um relato fácil de analisar e resolver.

## Migração do v2

### Devo migrar do v2 para o v3?

O v3 oferece suporte a várias janelas, uma API mais simples baseada em serviços, um atualizador integrado, um sistema de build muito mais flexível e melhor desempenho. Novos projetos devem começar com o v3. Para projetos existentes, o [Guia de Migração](/migration/v2-to-v3/) apresenta as diferenças.

### O v2 continuará recebendo manutenção?

Sim. O Wails v2 continua recebendo atualizações críticas. Novos projetos devem usar o Wails v3.0.0.

### Posso executar o v2 e o v3 lado a lado?

Sim. As CLIs são binários separados (`wails` e `wails3`), e os módulos têm caminhos de importação diferentes. Portanto, projetos em versões principais diferentes coexistem sem problemas na mesma máquina.

## Comunidade

### Como obtenho ajuda?

- [Discord](https://discord.gg/JDdSxwjhGf) para perguntas rápidas e discussões
- [Discussões do GitHub](https://github.com/wailsapp/wails/discussions) para perguntas mais detalhadas
- [Issues do GitHub](https://github.com/wailsapp/wails/issues) para bugs

### Como posso contribuir?

Consulte o [Guia de Contribuição](/contributing/). Correções de bugs são bem-vindas a qualquer momento; novas funcionalidades e alterações no comportamento público usam um PR de rascunho de [WEP (Wails Enhancement Proposal)](https://github.com/wailsapp/wails/blob/releases/v3.0.0/v3/wep/README.md). A discussão informal no Discord ou nas Discussões do GitHub é opcional.

### Onde encontro exemplos?

O repositório inclui mais de 60 exemplos executáveis que abrangem janelas, caixas de diálogo, eventos, a bandeja do sistema, serviços e muito mais: [v3/examples](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples).

## Ainda tem dúvidas?

Pergunte no [Discord](https://discord.gg/JDdSxwjhGf) ou [abra uma discussão](https://github.com/wailsapp/wails/discussions).
