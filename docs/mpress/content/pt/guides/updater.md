---
title: "Atualizador"
description: "Atualizações automáticas dentro do aplicativo para Wails v3 — provedores conectáveis, verificação criptográfica, substituição atômica e uma interface padrão que você pode personalizar com temas ou substituir."
slug: "guides/updater"
sourcePath: "guides/updater.md"
---

O atualizador distribui atualizações de software dentro do aplicativo sem obrigar você a criar seu próprio pipeline de download, verificação e substituição. Ele funciona sobre `app.Updater`, aceita um ou mais `Provider`s conectáveis (GitHub Releases, keygen.sh, Sparkle AppCast, o protocolo aberto Wails Update Manifest ou um provedor próprio), autentica os downloads usando uma chave pública configurada, substitui com segurança o binário em execução e expõe todas as transições por meio do barramento de eventos padrão do Wails.

![A janela padrão do atualizador no estado Atualização pronta — ícone de acordo com o estado, indicador de versão (v1.0.0 → v2.0.1 · 8.8 MB), notas da versão renderizadas em Markdown, incluindo uma tabela GFM, e uma única ação principal.](/assets/updater/default-window-ready.png)

## Início rápido

```go {title="main.go"}
package main

import (
    "context"
    "log"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/updater"
    "github.com/wailsapp/wails/v3/pkg/updater/providers/github"
)

func main() {
    app := application.New(application.Options{Name: "Demo"})

    gh, _ := github.New(github.Config{Repository: "myorg/myapp"})
    if err := app.Updater.Init(updater.Config{
        CurrentVersion: "1.0.0",
        Providers:      []updater.Provider{gh},
    }); err != nil {
        log.Fatal(err)
    }

    if err := app.Updater.CheckAndInstall(context.Background()); err != nil {
        log.Printf("update: %v", err)
    }

    _ = app.Run()
}
```

Isso abre a janela de atualização do framework, verifica o GitHub, baixa o artefato da plataforma, verifica-o, substitui o binário e aguarda o usuário reiniciar o aplicativo.

## O ciclo de vida

`app.Updater` é uma máquina de estados com os seguintes estados (`updater.State`):

| Estado | Quando |
| --- | --- |
| `unconfigured` | Antes de `Init` ser chamado |
| `idle` | Depois de `Init`, antes de qualquer verificação |
| `checking` | Um `Check` está em andamento |
| `up-to-date` | A resposta mais recente do provedor informou que o chamador está atualizado |
| `available` | Uma nova versão foi encontrada; o download ainda não começou |
| `downloading` | Os bytes estão sendo transmitidos pelo provedor |
| `verifying` | Download concluído; assinatura/resumo sendo verificado |
| `installing` | Os bytes verificados estão sendo descompactados e renomeados no diretório de preparação |
| `ready` | Atualização preparada; chame `Restart` para aplicá-la |
| `error` | Uma etapa anterior falhou |

Você pode consultar o estado atual a qualquer momento com `app.Updater.State()`. Cada transição também emite um evento do Wails (consulte [Eventos](#eventos)).

`Restart` aguarda o processo auxiliar alcançar `application.New` antes de solicitar o encerramento do aplicativo em execução. O tempo limite padrão de inicialização é de 30 segundos. Se o aplicativo executar uma inicialização demorada antes de `application.New`, defina `Config.HelperReadyTimeout` como uma duração maior, por exemplo, `time.Minute`. Zero seleciona o valor padrão; durações negativas são rejeitadas. Se a inicialização exceder o tempo limite, `Restart` retornará `updater.ErrHelperNotReady` e manterá o aplicativo em execução aberto.

A atualização é preparada em um diretório `wails-update-*` ao lado do binário instalado (no macOS, ao lado do pacote `.app`), portanto esse diretório precisa permitir gravação.

Dentro de um AppImage, o atualizador substitui e reinicia o arquivo `.AppImage` indicado por `$APPIMAGE`, em vez do ponto de montagem somente leitura.

A janela padrão reflete automaticamente o estado atual — por exemplo, quando `Check` não retorna nenhuma atualização, o usuário vê esta mensagem e fecha a janela com **Fechar**:

![A janela padrão do atualizador no estado Atualizado — marca de seleção verde, título “Você está atualizado” e um único botão Fechar.](/assets/updater/default-window-up-to-date.png)

## Provedores

Um `Provider` pode ser qualquer implementação que satisfaça esta interface:

```go
type Provider interface {
    Name() string
    Check(ctx context.Context, req CheckRequest) (*Release, error)
    Download(ctx context.Context, r *Release, dst io.Writer, onProgress func(written, total int64)) error
}
```

Quatro implementações estão incluídas na árvore do repositório.

### GitHub Releases — `updater/providers/github`

```go {title="github provider"}
gh, err := github.New(github.Config{
    Repository:    "myorg/myapp",     // your owner/repo (required)
    Token:         "ghp_…",           // optional; raises rate limit + private repos
    Prerelease:    false,             // include pre-releases in latest lookup
    ChecksumAsset: "SHA256SUMS",      // optional sibling asset for digest verification
    BaseURL:       "",                // optional override (e.g. GitHub Enterprise)
    AssetMatcher:  nil,               // optional custom asset-picker; nil uses DefaultAssetMatcher
    HTTPClient:    nil,               // optional client override
})
```

O seletor de artefatos padrão procura as substrings `GOOS` + `GOARCH` no nome do arquivo e reconhece aliases comuns (`amd64` / `x86_64` / `x64`, `arm64` / `aarch64`, `386` / `i386` / `x86` / `ia32`). Para esquemas de nomenclatura personalizados:

```go
gh, _ := github.New(github.Config{
    Repository: "myorg/myapp",
    AssetMatcher: func(req updater.CheckRequest, assets []github.ReleaseAsset) int {
        for i, a := range assets {
            if strings.Contains(a.Name, "my-naming-convention") &&
               strings.Contains(a.Name, req.Platform) {
                return i
            }
        }
        return -1 // no match
    },
})
```

`ChecksumAsset` é o nome de um artefato irmão da versão cujo conteúdo consiste em linhas `<sha256>  <filename>` (o formato gerado por `sha256sum` e `shasum -a 256`). O provedor o obtém durante `Check`, encontra a linha correspondente ao artefato selecionado e preenche `Release.Verification.Digest` para que o framework verifique o download.

### keygen.sh — `updater/providers/keygen`

```go {title="keygen provider"}
kg, err := keygen.New(keygen.Config{
    Account:    "your-account-slug", // required
    Product:    "product-uuid",      // optional but recommended when account has multiple products
    Package:    "",                  // optional further narrowing
    Channel:    "stable",            // "stable" / "rc" / "beta" / "alpha" / "dev"
    Filetype:   "",                  // optional artifact filetype filter ("dmg", "exe", …)
    Token:      "prod-…",            // product / environment / user / admin token; wins over LicenseKey
    LicenseKey: "",                  // license key auth (used only when Token is empty)
    BaseURL:    "",                  // optional API base override
    HTTPClient: nil,                 // optional client override; redirect-strip wrapper still applied
})
```

O provedor mapeia automaticamente a soma de verificação SHA-512 e a assinatura Ed25519ph de cada artefato do keygen.sh para o bloco `Release.Verification` do framework — nenhuma configuração adicional é necessária.

**Formatos de token:** os tokens do keygen.sh têm um prefixo de função (`admi-` / `prod-` / `envi-` / `user-`). O UUID puro exibido no painel é o *identificador* do token, não seu valor secreto — o segredo só fica visível no momento da criação do token. Consulte a [documentação de autenticação](https://keygen.sh/docs/api/authentication/) do keygen.sh para obter detalhes.

### Sparkle AppCast — `updater/providers/appcast`

```go {title="appcast provider"}
ac, err := appcast.New(appcast.Config{
    URL:        "https://your.app/appcast.xml", // required
    Channel:    "stable",                       // optional sparkle:channel filter
    HTTPClient: nil,                            // optional client override
})
```

Integra-se sem alterações à infraestrutura existente do Sparkle / WinSparkle. Lê `sparkle:shortVersionString`, `<enclosure url type length sparkle:os sparkle:edSignature>` e `sparkle:channel` do feed.

As assinaturas DSA de 1 do Sparkle (`sparkle:dsaSignature`) não são compatíveis — projetos que usam esse esquema de assinatura devem migrar para EdDSA (2 do Sparkle).

### Wails Update Manifest — `updater/providers/endpoint`

```go {title="endpoint provider"}
ep, err := endpoint.New(endpoint.Config{
    URL:        "https://updates.example.com/check", // required; supports {{platform}} / {{arch}} / {{version}} / {{channel}} placeholders
    Channel:    "stable",                            // optional channel filter
    Headers:    nil,                                 // optional headers, e.g. {"Authorization": "License <key>"}
    HTTPClient: nil,                                 // optional client override; redirect-strip wrapper still applied
})
```

Usa o protocolo aberto [Wails Update Manifest](/reference/update-manifest/): um documento JSON que descreve a versão mais recente e seus artefatos por plataforma, com somas de verificação e assinaturas embutidas. O mesmo documento funciona em uma hospedagem de arquivos estáticos (S3, GitHub Pages ou qualquer CDN — publique um manifesto por canal contendo todas as plataformas) ou em um servidor de atualização dinâmico (o provedor envia `platform`, `arch`, `version` e `channel` em cada verificação, permitindo que o servidor retorne exatamente um artefato ou condicione o acesso a uma licença).

Os placeholders de URL permitem configurar layouts estáticos em uma única linha:

```go
ep, _ := endpoint.New(endpoint.Config{
    URL: "https://cdn.example.com/updates/{{platform}}/{{arch}}/stable.json",
})
```

Os cabeçalhos configurados são enviados em todas as solicitações do manifesto; os downloads de artefatos só os reutilizam no próprio host do manifesto, sem rebaixamento de `https` para `http`, e o cabeçalho `Authorization` é removido em qualquer redirecionamento entre origens ou que faça esse rebaixamento.

A CLI cuida da publicação: `wails3 updater manifest` calcula os resumos, assina e descreve os arquivos da versão em um único comando, enquanto `wails3 updater verify` verifica novamente o resultado antes do envio. Consulte [Publicação com a CLI wails3](/reference/update-manifest/#publishing-with-the-wails3-cli).

### Cadeia de fallback

`Config.Providers` é ordenado. O atualizador percorre a lista sequencialmente: vence o primeiro provedor que retorna uma versão, e o primeiro que informa “atualizado” interrompe a cadeia (o fallback serve para o caso de o “provedor principal estar inacessível”, não para quando “os provedores discordam”). Um erro avança para o próximo provedor.

```go
app.Updater.Init(updater.Config{
    CurrentVersion: "1.0.0",
    Providers: []updater.Provider{
        kg, // primary: licensed customers
        gh, // fallback: public mirror
    },
})
```

### Como criar seu próprio provedor

São três métodos e cerca de 150 linhas em uma implementação típica. O Updater é responsável pela verificação, preparação atômica, substituição e janela — o código do provedor determina a próxima versão e transmite os bytes:

```go
type CustomProvider struct { /* config */ }

func (p *CustomProvider) Name() string { return "custom" }

func (p *CustomProvider) Check(ctx context.Context, req updater.CheckRequest) (*updater.Release, error) {
    // Hit your update endpoint, decide whether req.CurrentVersion is current,
    // and return either nil (no upgrade) or a *Release with Artifact + optional
    // Verification populated.
    // Errors here drop through to the next provider in Config.Providers.
}

func (p *CustomProvider) Download(ctx context.Context, r *updater.Release, dst io.Writer, onProgress func(int64, int64)) error {
    // Stream the artifact's bytes to dst. Call onProgress(written, total) as
    // bytes flow past; the Updater debounces emits to ~10 Hz on the event bus.
}
```

Use como referência os provedores incluídos na árvore do repositório — cada um ocupa um único arquivo Go.

## Verificação criptográfica

As versões são autenticadas pelo verificador do framework, usando `Config.PublicKey` como raiz de confiança:

```go
//go:embed publickey.pem
var publicKey []byte

app.Updater.Init(updater.Config{
    CurrentVersion: "1.0.0",
    Providers:      []updater.Provider{...},
    PublicKey:      publicKey,
})
```

Algoritmos compatíveis (`Release.Verification.SignatureAlgo`):

| Algoritmo | O que é assinado | Observações |
| --- | --- | --- |
| `ed25519` | O resumo SHA-256 do artefato | Usado pelo EdDSA do Sparkle |
| `ed25519ph` | O artefato completo por meio do pré-hash do Ed25519ph (SHA-512 internamente) | Usado pelo keygen.sh |
| `ecdsa-p256` | O resumo SHA-256 do artefato | São aceitas assinaturas `r∥s` brutas e DER |

Também há suporte somente a resumo (`DigestAlgo`: `sha256` / `sha512`) quando uma versão inclui um hash, mas não uma assinatura.

`Config.PublicKey` é a ÚNICA âncora de confiança para a verificação de assinaturas — a origem da versão não tem como substituí-la por uma chave própria. Versões que incluem uma `Signature` sem uma `Config.PublicKey` configurada falham de forma segura. O verificador calcula o resumo em uma passagem contínua durante o download; portanto, mesmo em atualizações com vários GB, a verificação não acrescenta outra passagem pelo disco.

@note{type="caution" title="Somente resumo ≠ verificação criptográfica"}
Uma versão que contém apenas `Digest` é autenticada com base no TLS do registro e em qualquer garantia de integridade fornecida pelo próprio registro — não com base em uma raiz criptográfica sob seu controle. Use somente resumo para detectar corrupção de dados; use assinaturas para resistir a adulterações quando o pipeline de lançamento estiver comprometido.

@end

### Como gerar uma chave de assinatura

```bash
wails3 updater genkey
# updater.key       — keep secret, use to sign releases
# updater.key.pub   — bundle in your app via go:embed
```

A chave privada usa PEM PKCS#8, e a chave pública usa PEM PKIX; `Config.PublicKey` aceita o arquivo `.pub` como está (também aceita a chave bruta de 32 bytes ou sua representação em base64, que `genkey` exibe para inclusão direta). Assine as versões com `wails3 updater manifest -key updater.key ...` ou `wails3 updater sign`; consulte [Publicação com a CLI wails3](/reference/update-manifest/#publishing-with-the-wails3-cli).

Ou em Go:

```go
import "crypto/ed25519"
import "crypto/rand"

pub, priv, _ := ed25519.GenerateKey(rand.Reader)
// Persist `priv` securely (HSM, signing CI, etc.); embed `pub` in your binary.
```

## Formatos de artefato

Os provedores transmitem os bytes do arquivo que você publicar; em seguida, o framework o descompacta antes da substituição:

- **Binário único** (por exemplo, `myapp-linux-amd64`) — usado como está. Comum no Linux.
- **`.zip`** — extraído no local de destino. O arquivo compactado deve conter exatamente uma entrada no nível superior (normalmente um pacote `.app` do macOS ou um único binário). É o empacotamento recomendado para macOS.
- **`.tar.gz`** / **`.tgz`** — extraído no local de destino seguindo a mesma regra de uma única entrada no nível superior. Útil para distribuições Linux que fornecem uma árvore de runtime junto ao binário.

Arquivos compactados com mais de uma entrada no nível superior são rejeitados: o framework substitui um único destino no disco, portanto, “colocar este arquivo compactado no lugar” é ambíguo quando ele contém vários itens. `.dmg` e `.pkg` (macOS), bem como `.msi` (Windows), não são compatíveis com v1 — em vez disso, distribua um `.zip` do pacote. A extração aplica proteção contra zip slip, rejeita links simbólicos que escapem da raiz do arquivo compactado e limita o tamanho total descompactado (2 GiB) e o número de entradas (50 000).

## A janela padrão

`app.Updater.CheckAndInstall(ctx)` abre uma janela de 520×540 gerenciada pelo framework com:

- Um ícone de destaque controlado pelo estado (↓ azul para disponível/baixando, ✓ verde para pronta/atualizada e ! vermelho para erro)
- Um indicador de versão: `v1.0.0 → v2.0.1 · 8.8 MB`
- Um painel rolável de notas da versão com Markdown **renderizado** (parágrafos, negrito/itálico, listas, tabelas GFM, código embutido, blocos de código delimitados, h1–h3 e links)
- Uma única ação principal por estado (Instalar / Reiniciar e aplicar / Tentar novamente)
- Ações secundárias em estilo discreto (Ignorar esta versão / Lembrar mais tarde)
- Modo escuro/claro por meio de `prefers-color-scheme`
- Efeito de brilho de progresso indeterminado quando o tamanho total é desconhecido

Ela escuta eventos `updater:*` no barramento de eventos do Wails e emite ações `updater:user:*` de volta para o Go.

### Tema por meio de variáveis CSS

```go
app.Updater.Init(updater.Config{
    // …
    Window: &updater.BuiltinWindow{
        CSS: `:root { --accent: #ff6f00; --bg: #1a1a1a; --fg: #fafafa; }`,
    },
})
```

A folha de estilos padrão expõe estas variáveis — substitua qualquer uma delas:

| Variável | Padrão (claro) | Padrão (escuro) |
| --- | --- | --- |
| `--bg` | `#f8f8fa` | `#1a1a1c` |
| `--surface` | `#ffffff` | `#232326` |
| `--surface-2` | `#f0f0f3` | `#2c2c30` |
| `--fg` | `#1d1d1f` | `#f5f5f7` |
| `--fg-dim` | `#6b6b73` | `#b0b0b8` |
| `--fg-faint` | `#99999f` | `#7a7a82` |
| `--border` | `#d6d6dc` | `#3a3a3e` |
| `--accent` | `#0a84ff` | `#0a84ff` |
| `--accent-fg` | `#ffffff` | — |
| `--success` | `#34c759` | — |
| `--error` | `#ff3b30` | — |
| `--radius` | `10px` | — |
| `--font` | pilha do sistema | — |

### Substitua o modelo

Forneça seu próprio HTML; ele só precisa escutar eventos `wails:updater:*` e emitir ações `wails:updater:user:*`:

```go
app.Updater.Init(updater.Config{
    // …
    Window: &updater.BuiltinWindow{
        HTML: myCustomTemplate,
    },
})
```

As janelas InitialHTML são carregadas sem uma origem do servidor de ativos, portanto não podem buscar `/wails/runtime.js` dinamicamente. Há duas maneiras de se comunicar com o host a partir de uma delas:

1. **Basta escrever HTML.** O framework injeta automaticamente um shim `window.wails.Events` mínimo em qualquer janela aberta com `WebviewWindowOptions.AllowSimpleEventEmit = true` e `HTML` definidos — exatamente como fazem os fluxos integrado e BYO do atualizador. Nenhuma etapa de build é necessária. Este é o fluxo usado no exemplo abaixo.
2. **Empacote `@wailsio/runtime` com o empacotador de sua preferência** (Vite, esbuild, Rollup) e importe-o para seu HTML personalizado durante a compilação. `Events.On` funciona sem configuração adicional porque opera exclusivamente no cliente; `Events.Emit` usa o transporte fetch do runtime, que não funciona com a origem nula — portanto, instale um pequeno transporte postMessage por meio do gancho [`setTransport`](/guides/custom-transport/) do runtime, que encaminha as mensagens por `window._wails.invoke("wails:event:emit:<name>")`. A injeção do framework não faz nada se `window.wails.Events` já estiver no escopo, evitando conflitos entre as duas abordagens.

De qualquer forma, o JavaScript escrito no HTML personalizado chama a mesma API `Events.On` / `Events.Emit`:

```html
<script>
const { On, Emit } = window.wails.Events;

On("wails:updater:update-available", (e) => {
    const rel = e.data ?? e;
    document.getElementById("ver").textContent = rel.version;
});

document.getElementById("install").addEventListener("click",
    () => Emit("wails:updater:user:install"));

// Ask the host to replay the current state so we paint correctly on (re)open.
Emit("wails:updater:window:ready");
</script>
```

O shim expõe o subconjunto do runtime moderno necessário para eventos com nomes simples: `Events.On(name, cb)` retorna uma função de cancelamento da inscrição, e `Events.Emit(nameOrEventObject)` encaminha ao host pelo caminho postMessage protegido por `wails:event:emit:`. Ele é instalado uma única vez durante o carregamento da página, antes da execução de qualquer um dos seus scripts inline.

Se você *quiser* substituir o shim (ou estiver carregando o runtime completo de outra maneira), defina `window.wails.Events` antes da execução da primeira tag `<script>` da página, e a injeção será ignorada.

### Moldura da janela

Substitua as opções da janela (tamanho, sem moldura, sempre visível) sem alterar o HTML:

```go
Window: &updater.BuiltinWindow{
    Options: updater.WindowOptions{
        Title:         "My App Updater",
        Width:         640,
        Height:        480,
        Frameless:     true,
        AlwaysOnTop:   true,
        DisableResize: false,
    },
},
```

### Use sua própria janela

Conduza o fluxo de atualização usando um `*application.WebviewWindow` criado por você. O atualizador chama `Show()` / `Close()` / `EmitEvent()` na sua janela — seu HTML decide o que renderizar:

![Uma janela de atualização própria, com fundo em gradiente rosa e laranja, um único cartão branco arredondado, tipografia personalizada e os mesmos eventos do atualizador controlando o estado visível. Demonstra como a interface padrão pode ser completamente substituída.](/assets/updater/byo-custom-window.png)

```go
myWin := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:                "My Updater",
    Width:                520, Height: 460,
    HTML:                 myCustomHTML,
    AllowSimpleEventEmit: true,           // see security note below
})
app.Updater.Init(updater.Config{
    // …
    Window: updater.BYOWindow(myWin.AsUpdaterWindow()),
})
```

Seu HTML usa `window.wails.Events.On` / `Events.Emit` exatamente como o modelo integrado — o shim injetado automaticamente pelo framework é incluído em qualquer janela com `AllowSimpleEventEmit: true`, independentemente de a janela pertencer ao framework ou a você. Consulte [Substituir o modelo](#substitua-o-modelo) para conhecer a API.

@note{type="caution" title="`AllowSimpleEventEmit` é obrigatório para janelas de atualização próprias"}
Por segurança, o framework condiciona o atalho postMessage `wails:event:emit:` a esse campo: uma janela sem o campo definido não pode gerar eventos personalizados no host. O shim de HTML personalizado do atualizador emite eventos `updater:user:*` por esse atalho; portanto, se o campo for esquecido em uma janela própria, todos os cliques nos botões serão silenciosamente ignorados — o usuário clica em Instalar e nada acontece.

Mantenha `AllowSimpleEventEmit` **desativado** em qualquer janela que carregue HTML que você não controle totalmente (URLs remotas ou conteúdo fornecido pelo usuário). Quando ativado, qualquer JavaScript na página — inclusive pontos vulneráveis a XSS — pode acionar qualquer manipulador `app.Event.On(name, …)`. O atalho transmite apenas nomes simples (sem payload) e não consegue acessar o caminho de binding/Call, mas ainda pode acionar manipuladores privilegiados de eventos personalizados no seu código Go caso eles atuem apenas com base no nome do evento.

A janela *integrada* do atualizador do framework define isso internamente — somente quem usa uma janela própria precisa se lembrar desse campo.

@end

### Sem interface gráfica

```go
app.Updater.Init(updater.Config{
    // …
    Window: updater.WindowNone,
})
```

Nenhuma janela é aberta. Inscreva sua própria interface (ou a janela principal existente) nos eventos `updater:*` e chame `app.Updater.CheckAndInstall(ctx)` em um manipulador de botão. Isso é útil para verificações periódicas em segundo plano que só devem aparecer quando algo for encontrado ou para aplicativos que integram o fluxo de atualização a um painel de configurações personalizado.

## Eventos

Tanto Go quanto JavaScript fazem inscrições pelo barramento de eventos padrão do Wails. **Não digite manualmente as strings usadas na comunicação** — use as constantes exportadas pelo pacote do atualizador (Go) ou pelo pacote do runtime (JS). As duas camadas compartilham o mesmo conjunto de nomes, mantido em sincronia por um teste de regressão.

### No Go

As constantes ficam em `github.com/wailsapp/wails/v3/pkg/updater`. Faça a inscrição por meio de `app.Event.On(name, fn)`; o callback recebe um `*application.CustomEvent` cujo campo `Data` contém o payload tipado listado na [referência de eventos](#referncia-de-eventos) — use uma asserção de tipo, não a decodificação de JSON:

```go {title="main.go"}
import (
    "log"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/updater"
)

// …

app.Event.On(updater.EventUpdateAvailable, func(e *application.CustomEvent) {
    rel, ok := e.Data.(*updater.Release)
    if !ok { return }
    log.Printf("update found: %s", rel.Version)
})

app.Event.On(updater.EventDownloadProgress, func(e *application.CustomEvent) {
    p, ok := e.Data.(updater.Progress)
    if !ok { return }
    log.Printf("%d / %d bytes (%.0f KB/s)", p.Written, p.Total, p.Rate/1024)
})

app.Event.On(updater.EventError, func(e *application.CustomEvent) {
    info, ok := e.Data.(updater.ErrorInfo)
    if !ok { return }
    log.Printf("update failed during %s: %s", info.Stage, info.Message)
})
```

Todas as constantes Go disponíveis:

| Constante | String de comunicação |
| --- | --- |
| `updater.EventCheckStarted` | `wails:updater:check-started` |
| `updater.EventUpdateAvailable` | `wails:updater:update-available` |
| `updater.EventNoUpdate` | `wails:updater:no-update` |
| `updater.EventDownloadStarted` | `wails:updater:download-started` |
| `updater.EventDownloadProgress` | `wails:updater:download-progress` |
| `updater.EventDownloadComplete` | `wails:updater:download-complete` |
| `updater.EventVerifying` | `wails:updater:verifying` |
| `updater.EventInstalling` | `wails:updater:installing` |
| `updater.EventUpdateReady` | `wails:updater:update-ready` |
| `updater.EventError` | `wails:updater:error` |
| `updater.EventMeta` | `wails:updater:meta` |
| `updater.EventWindowReady` | `wails:updater:window:ready` |
| `updater.EventUserInstall` | `wails:updater:user:install` |
| `updater.EventUserSkip` | `wails:updater:user:skip` |
| `updater.EventUserRemind` | `wails:updater:user:remind` |
| `updater.EventUserCancel` | `wails:updater:user:cancel` |
| `updater.EventUserRestart` | `wails:updater:user:restart` |

### No JavaScript

As constantes ficam em `Updater.Events`, dentro de `@wailsio/runtime`. Elas têm os mesmos nomes usados no Go e são organizadas por subnamespace (`User.*`, `Window.*`) para facilitar sua descoberta pelo preenchimento automático:

```js
import { Events, Updater } from "@wailsio/runtime";

Events.On(Updater.Events.UpdateAvailable, (e) => {
    console.log("update found:", e.data.version);
});

Events.On(Updater.Events.DownloadProgress, (e) => {
    const p = e.data;
    console.log(`${p.written} / ${p.total} bytes (${(p.rate/1024).toFixed(1)} KB/s)`);
});

Events.On(Updater.Events.Error, (e) => {
    const info = e.data;
    console.error(`update failed during ${info.stage}: ${info.message}`);
});
```

Os eventos de ações do usuário que seu HTML personalizado emite *de volta* para o host ficam em `Updater.Events.User`:

```js
import { Updater } from "@wailsio/runtime";

document.getElementById("install-btn").addEventListener("click", () => {
    // The framework window does this internally via the postMessage shim;
    // shown here for BYO templates that need to drive the flow themselves.
    window._wails.invoke("wails:event:emit:" + Updater.Events.User.Install);
});
```

### Referência de eventos

Lado da inscrição (host → página):

| Constante (Go) | Constante (JS) | Payload | Quando |
| --- | --- | --- | --- |
| `updater.EventCheckStarted` | `Updater.Events.CheckStarted` | nenhum | Antes de cada operação de ida e volta de `Check` |
| `updater.EventUpdateAvailable` | `Updater.Events.UpdateAvailable` | `*Release` | `Check` encontrou uma versão mais recente |
| `updater.EventNoUpdate` | `Updater.Events.NoUpdate` | nenhum | `Check` confirmou que o aplicativo está atualizado |
| `updater.EventDownloadStarted` | `Updater.Events.DownloadStarted` | `*Release` | Os bytes começam a ser transmitidos |
| `updater.EventDownloadProgress` | `Updater.Events.DownloadProgress` | `Progress` | ~10 Hz durante o download |
| `updater.EventDownloadComplete` | `Updater.Events.DownloadComplete` | `*Release` | Todos os bytes foram gravados, antes da verificação |
| `updater.EventVerifying` | `Updater.Events.Verifying` | `*Release` | A verificação da assinatura/resumo é iniciada |
| `updater.EventInstalling` | `Updater.Events.Installing` | `*Release` | O desempacotamento e a preparação são iniciados |
| `updater.EventUpdateReady` | `Updater.Events.UpdateReady` | `*Release` | Reinicialização pendente |
| `updater.EventError` | `Updater.Events.Error` | `ErrorInfo` | Falha em qualquer etapa |
| `updater.EventMeta` | `Updater.Events.Meta` | `Meta` | Uma vez por sessão, antes da reprodução do snapshot |

No lado da página (página → host) — seu código se inscreve caso você crie um modelo personalizado:

| Constante (Go) | Constante (JS) | Quando |
| --- | --- | --- |
| `updater.EventWindowReady` | `Updater.Events.Window.Ready` | A janela terminou de carregar; o host reproduz o estado atual |
| `updater.EventUserInstall` | `Updater.Events.User.Install` | Ação principal no estado `available` |
| `updater.EventUserRestart` | `Updater.Events.User.Restart` | Ação principal no estado `ready` |
| `updater.EventUserSkip` | `Updater.Events.User.Skip` | "Ignorar esta versão" |
| `updater.EventUserRemind` | `Updater.Events.User.Remind` | "Lembrar mais tarde" |
| `updater.EventUserCancel` | `Updater.Events.User.Cancel` | Botão Fechar |

## Referência da API

### `updater.Config`

| Campo | Tipo | Observações |
| --- | --- | --- |
| `CurrentVersion` | `string` | **Obrigatório.** A mesma string usada para marcar as versões (sem o prefixo `v`) |
| `Providers` | `[]updater.Provider` | **Obrigatório.** Cadeia de fallback ordenada |
| `PublicKey` | `[]byte` | PEM ou bytes brutos. Opcional, mas versões assinadas falham de forma segura sem esse valor |
| `CheckInterval` | `time.Duration` | Um valor diferente de zero inicia um loop de consulta em segundo plano que chama `CheckAndInstall` |
| `Platform` | `string` | Substitui `runtime.GOOS` na seleção de artefatos |
| `Arch` | `string` | Substitui `runtime.GOARCH` na seleção de artefatos |
| `Channel` | `string` | Atualmente apenas informativo; filtragem por canal específica do provedor |
| `Window` | `updater.WindowOption` | `nil` (padrões integrados), `&BuiltinWindow{…}`, `BYOWindow(handle)` ou `WindowNone` |

### Métodos de `*updater.Updater`

| Assinatura | Finalidade |
| --- | --- |
| `Init(cfg Config) error` | Configura. Retorna `ErrAlreadyConfigured` na segunda chamada |
| `State() State` | Fase atual do ciclo de vida |
| `CurrentVersion() string` | A versão passada para `Init` |
| `Check(ctx) (*Release, error)` | Percorre a cadeia de provedores. `(rel, nil)` = encontrada, `(nil, nil)` = atualizada, `(nil, err)` = todos falharam |
| `DownloadAndInstall(ctx) error` | Transmite, verifica, extrai (se for um arquivo compactado) e prepara. Requer uma chamada anterior a `Check` |
| `CheckAndInstall(ctx) error` | Atalho: abre a janela, executa `Check` e, se encontrar, executa `DownloadAndInstall` |
| `Restart(ctx) error` | Inicia o auxiliar, chama `Host.Quit` e encerra; o auxiliar faz a substituição e reinicia o aplicativo |
| `DownloadedPath() string` | Local no disco em que a atualização preparada está armazenada, ou `""` se não houver nenhuma |
| `SkipVersion(v string)` | Registra `v` como ignorada; chamadas subsequentes a `Check` a consideram atualizada |
| `SkippedVersion() string` | Lê a versão atualmente ignorada |
| `StopPeriodicCheck()` | Cancela o temporizador iniciado por `Config.CheckInterval` e aguarda o retorno do loop |

### Erros

| Sentinela | Retornado por |
| --- | --- |
| `ErrAlreadyConfigured` | `Init` após o primeiro sucesso |
| `ErrNotConfigured` | Qualquer operação antes de `Init` |
| `ErrNoPendingRelease` | `DownloadAndInstall` sem um `Check` anterior |
| `ErrDownloadInProgress` | `DownloadAndInstall` chamado enquanto outro está em execução |
| `ErrNotReady` | `Restart` sem nenhuma atualização preparada |

## Como funciona a substituição

`Restart` executa novamente o binário atual com variáveis de ambiente sentinela definidas. `application.New` as detecta durante a inicialização e desvia a execução para o modo auxiliar:

1. O auxiliar aguarda até 30 s para que o PID do processo pai seja encerrado (`platformIsAlive` consulta periodicamente por meio de `syscall.OpenProcess` + `GetExitCodeProcess` no Windows e `os.FindProcess` + `proc.Signal(syscall.Signal(0))` no Unix).
2. O auxiliar cria um backup do destino (cópia para arquivos e cópia recursiva para diretórios de pacotes `.app` do macOS).
3. O auxiliar substitui o destino pelo artefato preparado, tentando novamente até 20 vezes, com uma espera de 500 ms entre as tentativas:
  - **Unix** — `os.RemoveAll(target)` + `os.Rename(newPath, target)`. Descritores de arquivo abertos para o inode antigo continuam válidos.
  - **Windows** — `os.Rename(target, target.old.<nanos>)` + `os.Rename(newPath, target)`. O Windows permite renomear arquivos cuja imagem ainda esteja mapeada, mas não excluí-los; na próxima atualização, o auxiliar remove quaisquer arquivos irmãos `.old.*` remanescentes cujo mapeamento de kernel proprietário já tenha sido liberado.

4. O auxiliar restaura no novo binário o modo de execução do arquivo original (o arquivo baixado foi criado com a umask padrão, que remove `+x` no Unix; no Windows, essa operação não tem efeito).
5. O auxiliar remove as variáveis de ambiente do modo auxiliar e reinicia o binário, agora substituído.
6. O auxiliar é encerrado.

Se a inicialização falhar, o auxiliar restaura o backup. Se o processo pai não for encerrado em até 30 s, o auxiliar aborta antes de modificar o destino (assim, o usuário mantém um aplicativo funcional mesmo se uma caixa de diálogo de encerramento bloquear `Quit`).

Para pacotes `.app` do macOS distribuídos como `.zip` (o empacotamento recomendado), o arquivo é descompactado entre as etapas de verificação e prontidão, para que o auxiliar tenha um diretório real para colocar no lugar do atual.

## Verificação periódica

```go
app.Updater.Init(updater.Config{
    // …
    CheckInterval: 6 * time.Hour,
})
```

Quando `CheckInterval > 0`, uma goroutine em segundo plano chama `CheckAndInstall` no intervalo configurado. Os disparos que ocorrem enquanto outro fluxo já está em andamento (verificação, download, validação ou instalação) são descartados — máquinas de estado simultâneas não são compatíveis.

Para realizar consultas silenciosas em segundo plano que só sejam exibidas quando algo for encontrado, defina `Window: updater.WindowNone` e responda a `EventUpdateAvailable` em sua própria interface.

## Ignorar e lembrar depois

O botão "Ignorar esta versão" da janela padrão registra a versão disponível por meio de `SkipVersion(rel.Version)`. Chamadas posteriores a `Check` encontram a mesma versão e a consideram atualizada (até que o usuário atualize `CurrentVersion`, o que ocorre automaticamente após um `Restart` bem-sucedido). "Lembrar mais tarde" apenas fecha a janela sem registrar nada.

```go
// Reading what the user skipped (e.g. to surface in app settings)
if v := app.Updater.SkippedVersion(); v != "" {
    log.Printf("user skipped %s", v)
}

// Programmatically clearing the skip:
app.Updater.SkipVersion("")
```

## Lista de verificação para distribuição

Antes de publicar uma versão que será instalada pelo atualizador:

1. **Escolha o formato de arquivo correto.** macOS: `.zip` do pacote `.app`. Linux: binário único ou `.tar.gz`. Windows: um único `.exe` ou `.zip`. `.dmg`, `.msi` e `.pkg` não são compatíveis.
2. **Assine o artefato** com a chave privada correspondente a `Config.PublicKey`. Para feeds publicados por provedores (keygen.sh, AppCast), siga o fluxo de assinatura de cada provedor. Para GitHub Releases com `ChecksumAsset`, gere um arquivo `SHA256SUMS` com `sha256sum` / `shasum -a 256`.
3. **Faça a string de versão corresponder.** `Config.CurrentVersion` e a tag de versão da release devem ser exatamente iguais (por exemplo, `1.0.0` ↔ tag `v1.0.0`; o `v` inicial é removido no lado do provedor).
4. **Teste a substituição na plataforma de destino** pelo menos uma vez antes de distribuir — o tratamento de assinatura de código, notarização e Gatekeeper é específico da plataforma e não é realizado pelo próprio atualizador.

## Solução de problemas

**"a assinatura exige uma chave pública, mas nenhuma foi configurada"** — a release tem um campo `Signature`, mas `Config.PublicKey` está vazio. Defina a chave pública ou altere seu pipeline de release para não incluir uma assinatura.

**"incompatibilidade de resumo"** — os bytes baixados não correspondem ao que o provedor informou. Geralmente, isso indica um download parcial (falha momentânea de rede) ou um artefato corrompido. Executar novamente costuma resolver o problema.

**A janela abre, mas desaparece imediatamente, sem Markdown nem progresso** — seu HTML personalizado não chamou `wails:runtime:ready`. Consulte o shim [Substituir o modelo](#substitua-o-modelo).

**A atualização no Windows nunca termina; o log do auxiliar informa "remove old (attempt N): Access is denied"** — isso só acontece em versões desta PR anteriores a `de764fb`; a implementação atual usa renomeação para um nome alternativo, que não apresenta esse problema. Atualize.

**O Gatekeeper do macOS bloqueia o binário substituído** — a assinatura de código deve ser preservada de ponta a ponta. Assine o `.app` original *e* assine novamente o binário reiniciado caso seu pipeline de compilação modifique os direitos durante a atualização.

## Veja também

- Exemplo executável: [`v3/examples/updater`](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/updater)
- Repositório de demonstração para testes: [`wailsapp/updater-demo`](https://github.com/wailsapp/updater-demo)
- Tutorial: [Como adicionar atualizações automáticas a um aplicativo Wails](/tutorials/04-self-update-a-wails-app/)
