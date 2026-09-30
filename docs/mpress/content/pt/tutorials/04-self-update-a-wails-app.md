---
title: "Aplicativo Wails com atualização automática"
description: "Crie um aplicativo Wails v3 que se atualiza pelo GitHub Releases — desde `wails3 init` até a verificação da versão assinada e a substituição no modo auxiliar."
slug: "tutorials/04-self-update-a-wails-app"
sourcePath: "tutorials/04-self-update-a-wails-app.md"
---

Neste tutorial, você adicionará um atualizador integrado a um novo aplicativo Wails v3. Ao final, o aplicativo poderá:

- Verificar o GitHub Releases sob demanda (e, opcionalmente, por um temporizador).
- Baixar o artefato correto para o sistema operacional e a arquitetura em execução.
- Verificar um resumo SHA-256 (e, opcionalmente, uma assinatura Ed25519) em relação aos bytes baixados.
- Exibir as notas da versão na janela de atualização padrão do framework.
- Substituir o binário em execução e reiniciar o aplicativo — tudo isso sem distribuir um executável auxiliar separado.

Usaremos o **GitHub Releases** como fonte de atualizações porque ele é gratuito e não exige infraestrutura. Os mesmos padrões funcionam com o [keygen.sh](/guides/updater/#keygensh--updaterproviderskeygen) e o [Sparkle AppCast](/guides/updater/#sparkle-appcast--updaterprovidersappcast) — quando terminar, consulte o [guia do atualizador](/guides/updater/).

@note{type="tip" title="Pré-requisitos"}
- Go 1.25 ou mais recente
- CLI `wails3` instalada (`go install github.com/wailsapp/wails/v3/cmd/wails3@latest`)
- Um repositório do GitHub no qual você possa publicar versões
- Conhecer o [tutorial do serviço de QR Code](/tutorials/01-creating-a-service/) é útil, mas não obrigatório

@end

<br/>

@steps
### Comece com um novo aplicativo Wails
Gere a estrutura de um novo projeto com o modelo vanilla:

```bash
wails3 init -n updater-tutorial -t vanilla
cd updater-tutorial
```

Agora você deve ter um diretório com `main.go`, `frontend/` e um `Taskfile.yml`. Confirme que ele é compilado e iniciado:

```bash
wails3 task dev
```

Uma janela vazia do Wails deve ser aberta. Feche-a e continue.

### Adicione a importação do atualizador
Abra `main.go` e adicione os dois pacotes do atualizador às importações:

```go {title="main.go" ins="6-7"}
package main

import (
    _ "embed"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/updater"
    "github.com/wailsapp/wails/v3/pkg/updater/providers/github"
)
```

Esses pacotes importam o próprio Updater e o provedor do GitHub Releases.

### Configure o Updater
`app.Updater` já está integrado a cada `*application.App` — basta chamar `Init`:

```go {title="main.go"}
const currentVersion = "1.0.0"

gh, err := github.New(github.Config{
    Repository:    "yourorg/your-repo",   // ← change this
    ChecksumAsset: "SHA256SUMS",          // sibling file with sha256 digests
})
if err != nil {
    log.Fatalf("github.New: %v", err)
}

if err := app.Updater.Init(updater.Config{
    CurrentVersion: currentVersion,
    Providers:      []updater.Provider{gh},
}); err != nil {
    log.Fatalf("Updater.Init: %v", err)
}
```

Coloque isto depois de `application.New` e antes de `app.Run()`.

@note{type="note" title="Formato da string de versão"}
Passe a mesma versão usada para marcar as versões, **sem** o `v` inicial. O provedor remove `v` dos nomes das tags no lado dele. `1.0.0` aqui ↔ `v1.0.0` no GitHub.

@end

### Adicione um item de menu que acione a atualização
No mesmo `main.go`, adicione uma entrada de menu "Verificar atualizações…":

```go {title="main.go"}
menu := app.Menu.New()
app.Menu.SetApplicationMenu(menu)
appMenu := menu.AddSubmenu("App")
appMenu.Add("Check for Updates…").OnClick(func(*application.Context) {
    go func() {
        if err := app.Updater.CheckAndInstall(context.Background()); err != nil {
            app.Logger.Error("update", "error", err)
        }
    }()
})
```

`CheckAndInstall` abre a janela de atualização do framework, executa `Check` e, se encontrar uma versão, executa `DownloadAndInstall` automaticamente. Quando não há novidades, a janela permanece aberta no estado "Atualizado" — o usuário a fecha com o botão **Fechar**.

@note{type="caution" title="Execute em uma goroutine"}
`CheckAndInstall` fica bloqueado até a conclusão da verificação e da instalação. Chamá-lo diretamente no clique do menu bloquearia a thread da interface. Encapsule a chamada em `go func()`.

@end

### Execute uma vez sem versões publicadas
```bash
wails3 task dev
```

Clique em **Aplicativo → Verificar atualizações…**. A janela de atualização deve abrir brevemente, consultar a API do GitHub, não encontrar versões mais recentes que `1.0.0` e permanecer no estado **Atualizado**, com um ✓ verde.

Se ocorrer um erro aqui, geralmente será um destes:

| Sintoma | Correção |
| --- | --- |
| `404 Not Found` | O campo `Repository` está incorreto — deve ser `owner/repo` |
| `403 rate-limited` | Adicione `Token: "ghp_…"` a github.Config (use um PAT com o escopo `public_repo`) |
| Erros de rede | Confirme que o aplicativo em execução consegue acessar `api.github.com` |

### Publique uma versão de teste
Atualize `currentVersion` em `main.go` para `1.0.0` (ou deixe como está). Compile para uma plataforma a fim de obter um binário que você possa anexar a uma versão:

@tabs
[macOS]
```bash
wails3 task build:darwin
# produces bin/updater-tutorial.app
# zip it for the release asset:
cd bin && zip -r updater-tutorial-darwin-arm64.zip updater-tutorial.app && cd ..
```

[Linux]
```bash
wails3 task build:linux
# produces bin/updater-tutorial
mv bin/updater-tutorial bin/updater-tutorial-linux-amd64
```

[Windows]
```bash
wails3 task build:windows
# produces bin/updater-tutorial.exe
mv bin/updater-tutorial.exe bin/updater-tutorial-windows-amd64.exe
```

@end

Gere um arquivo `SHA256SUMS` ao lado do binário:

```bash
cd bin
shasum -a 256 updater-tutorial-* > SHA256SUMS
cat SHA256SUMS
```

Você deve ver uma ou mais linhas como esta:

```
abc123…  updater-tutorial-darwin-arm64.zip
```

Agora publique isso como **v2.0.0** no seu repositório do GitHub:

```bash
gh release create v2.0.0 \
    --title "v2.0.0" \
    --notes "First update for the self-update tutorial.

- **Bold** Markdown renders in the update window
- \`Code spans\` too
- Lists work
- GFM tables work" \
    bin/SHA256SUMS bin/updater-tutorial-*
```

@note{type="note" title="Nomenclatura dos artefatos"}
O seletor de artefatos padrão escolhe pelo trecho `GOOS` + `GOARCH` no nome do arquivo. Desde que o nome do artefato inclua `darwin` (ou `linux` / `windows`) e `arm64` (ou `amd64` / `386`), o seletor o encontrará. Consulte o [guia do atualizador](/guides/updater/#github-releases--updaterprovidersgithub) para saber como usar seletores personalizados.

@end

### Execute o aplicativo e verifique a atualização
Com `currentVersion` ainda definido como `1.0.0`, execute o aplicativo novamente:

```bash
wails3 task dev
```

Clique em **Aplicativo → Verificar atualizações…**. Desta vez, você deve ver algo assim:

![Janela padrão do atualizador no estado Atualização pronta, mostrando o indicador da versão, as notas da versão renderizadas em Markdown e o botão principal Reiniciar e aplicar.](/assets/updater/default-window-ready.png)

- O ícone principal muda de ↓ azul ("Atualização disponível") para ✓ verde ("Atualização pronta").
- O subtítulo mostra `v1.0.0 → v2.0.0 · <size>`.
- O painel de notas da versão renderiza o Markdown com negrito, trechos de código e a tabela.
- A barra de progresso é preenchida durante o download (será rápido — o binário é pequeno).

O Updater prepara o novo binário em um diretório temporário. Para concluir a atualização:

- Clique em **Reiniciar e aplicar**.
- O aplicativo é encerrado, o auxiliar substitui o binário e o novo binário reinicia.
- O aplicativo reiniciado informa `currentVersion = "1.0.0"` (porque definimos esse valor diretamente no código), mas os bytes no disco correspondem à compilação v2.0.0.

Em um aplicativo real, `currentVersion` seria definido no momento da compilação por meio de `-ldflags`, para que o novo binário saiba que agora é v2.0.0 e uma verificação posterior não encontre nenhuma atualização.

### Vincule `currentVersion` à compilação
Substitua a constante por uma variável definida no momento da compilação:

```go {title="main.go" ins="2,4"}
var (
    currentVersion = "dev" // overridden by -ldflags at release time
)
```

Depois, no comando de compilação:

```bash
wails3 task build:darwin -- -ldflags "-X main.currentVersion=2.0.0"
```

Ou adicione `-ldflags` ao seu `Taskfile.yml` para que o valor seja obtido de `git describe --tags`.

### Adicione uma assinatura criptográfica (recomendado para produção)
O caminho SHA256SUMS verifica a *integridade* (os bytes correspondem ao que o GitHub armazenou), mas não a *autenticidade* (se esses bytes foram produzidos pelo seu pipeline de lançamento, e não por meio de uma conta de mantenedor comprometida). Para oferecer resistência à adulteração, assine cada versão com uma chave Ed25519:

```bash
# One-time: generate the keypair
ssh-keygen -t ed25519 -f updater-key -N "" -C "wails-updater"
#   updater-key      — keep secret (build server, HSM, password manager)
#   updater-key.pub  — bundle in your app
```

A cada versão, assine com sua chave privada o resumo SHA-256 de cada artefato. Um pequeno utilitário em Go:

```go {title="cmd/sign-release/main.go"}
package main

import (
    "crypto/ed25519"
    "crypto/sha256"
    "encoding/base64"
    "fmt"
    "io"
    "os"
)

func main() {
    priv, _ := os.ReadFile("updater-key")
    key := ed25519.PrivateKey(priv) // raw 64-byte private key

    f, _ := os.Open(os.Args[1])
    defer f.Close()
    h := sha256.New()
    _, _ = io.Copy(h, f)
    sig := ed25519.Sign(key, h.Sum(nil))
    fmt.Println(base64.StdEncoding.EncodeToString(sig))
}
```

No momento, o provedor padrão do GitHub não busca um arquivo de assinatura separado — você pode [criar um provedor personalizado](/guides/updater/#writing-your-own-provider) que faça isso ou mudar para o **keygen.sh**, que assina cada artefato no servidor e disponibiliza tanto o resumo quanto a assinatura por meio de sua API.

Incorpore a chave pública ao seu aplicativo:

```go {title="main.go" ins="1,6"}
//go:embed updater-key.pub
var updaterPublicKey []byte

app.Updater.Init(updater.Config{
    CurrentVersion: currentVersion,
    PublicKey:      updaterPublicKey,
    Providers:      []updater.Provider{gh},
})
```

Com `PublicKey` definido, qualquer versão que forneça um `Signature` deve ser validada com essa chave. A origem da versão não tem como substituir a chave por outra — esse é justamente o objetivo de fixá-la por um canal independente durante a compilação.

### Personalize a janela
A janela padrão atende ao caso mais comum. Se precisar de mais controle, há três alternativas — escolha uma conforme o grau de personalização desejado:

@tabs
[Somente CSS]
```go
app.Updater.Init(updater.Config{
    // …
    Window: &updater.BuiltinWindow{
        CSS: `:root { --accent: #ff6f00; --radius: 16px; }`,
    },
})
```

Consulte a seção [Tema por meio de variáveis CSS](/guides/updater/#theme-via-css-variables) para ver a lista completa de variáveis.

[HTML personalizado]
```go
//go:embed updater-window.html
var updaterHTML string

app.Updater.Init(updater.Config{
    // …
    Window: &updater.BuiltinWindow{HTML: updaterHTML},
})
```

Seu HTML deve assinar os eventos `updater:*` e emitir as ações `updater:user:*` pelo canal de eventos do Wails. Consulte [Substituir o modelo](/guides/updater/#replace-the-template) para ver o shim de JS.

[Use sua própria janela]
```go
myWin := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:                "My App Updater",
    Width:                520, Height: 460,
    HTML:                 updaterHTML,
    AllowSimpleEventEmit: true,  // required — see security note
})
app.Updater.Init(updater.Config{
    // …
    Window: updater.BYOWindow(myWin.AsUpdaterWindow()),
})
```

Essa opção é útil quando você já tem sua própria infraestrutura de janelas e quer que o atualizador a controle, em vez de abrir outra janela. Um modelo HTML totalmente personalizado — controlado pelos mesmos eventos do atualizador usados pelo modelo padrão — tem esta aparência:

![Uma janela própria para o atualizador, com fundo em gradiente rosa e laranja e um layout personalizado de cartão com cantos arredondados, demonstrando que a interface padrão pode ser totalmente substituída.](/assets/updater/byo-custom-window.png)

@note{type="caution" title="`AllowSimpleEventEmit` é obrigatório"}
O shim de HTML personalizado do atualizador aciona Instalar / Ignorar / Lembrar / Reiniciar pelo atalho `wails:event:emit:` postMessage, que, por segurança, só fica disponível quando esse campo está habilitado. Se você se esquecer dele, os botões não farão nada, sem exibir qualquer aviso. Não o habilite em janelas que carregam HTML que você não controla completamente — consulte a seção [Use sua própria janela](/guides/updater/#bring-your-own-window) do guia para entender o modelo de ameaças.

@end

@end

### Execute verificações automáticas em segundo plano
Para executar a verificação periodicamente, em vez de usar o clique no menu (ou além dele):

```go {ins="5"}
app.Updater.Init(updater.Config{
    CurrentVersion: currentVersion,
    Providers:      []updater.Provider{gh},
    PublicKey:      updaterPublicKey,
    CheckInterval:  6 * time.Hour,
})
```

Cada intervalo executa o mesmo fluxo `CheckAndInstall` de um clique manual. Defina `Window: updater.WindowNone` se quiser que a verificação periódica permaneça silenciosa até que algo seja realmente encontrado — depois, assine `EventUpdateAvailable` por conta própria para decidir qual experiência apresentar ao usuário.

@end

## Tudo pronto

Agora você tem um aplicativo Wails que:

- Verifica se há atualizações nas versões do GitHub sob demanda e periodicamente.
- Renderiza as notas da versão como Markdown em uma janela padrão bem-acabada.
- Valida os downloads usando um resumo SHA-256 publicado por você.
- Opcionalmente, valida uma assinatura Ed25519 usando uma chave pública incorporada durante a compilação.
- Substitui o binário em execução no próprio local e reinicia o aplicativo automaticamente.

## Próximas etapas

- O [guia do atualizador](/guides/updater/) contém a referência completa da API, todos os eventos, todas as opções de configuração e o funcionamento da substituição no modo auxiliar.
- Consulte [`v3/examples/updater`](https://github.com/wailsapp/wails/tree/master/v3/examples/updater) para ver um exemplo funcional completo que você pode clonar.
- O repositório de destino para testes [`wailsapp/updater-demo`](https://github.com/wailsapp/updater-demo) mostra a estrutura recomendada para os artefatos da versão.

## Cuidados no ambiente de produção

- **Assinatura de código no macOS** — o Gatekeeper exige que o binário substituído seja assinado e autenticado por notarização. Assine seu pacote `.app` *antes* de compactá-lo em ZIP para a versão. O atualizador preserva os bytes exatamente como estão; ele não assina nada novamente.
- **Antivírus no Windows** — arquivos `.exe` não assinados baixados da internet podem acionar avisos do SmartScreen. Assine seu binário com um certificado Authenticode ou considere que usuários de máquinas com restrições rigorosas talvez precisem adicionar seu aplicativo à lista de permissões.
- **Versões atômicas** — publique `SHA256SUMS` (e seus binários) juntos, não em commits separados. O atualizador baixa o arquivo auxiliar separadamente do binário; se eles ficarem dessincronizados, a verificação do resumo falhará de modo seguro.
- **Versões ignoradas** — o botão "Ignorar esta versão" da janela padrão registra localmente que a versão deve ser ignorada. Se você lançar uma atualização de segurança crítica, atribua a ela um novo número de versão para que não seja ignorada automaticamente por usuários que dispensaram uma versão anterior.
