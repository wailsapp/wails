---
title: "Servidor de ativos"
description: "Como o Wails v3 disponibiliza e incorpora seus ativos web nos ambientes de desenvolvimento e produção"
slug: "contributing/asset-server"
sourcePath: "contributing/asset-server.md"
---

## Visão geral

Cada aplicativo Wails é distribuído como um **único executável nativo** que combina:

1. Seu backend em *Go*
2. Um frontend *web* (HTML + JS + CSS)

O **servidor de ativos** é o elemento que torna isso possível. Ele tem **dois modos de operação**, selecionados em tempo de compilação por meio de tags de build do Go:

| Modo | Tag | Finalidade |
| --- | --- | --- |
| **Desenvolvimento** | `//go:build !production` | Iteração rápida com recarregamento automático |
| **Produção** | `//go:build production` | Ativos incorporados, sem dependências |

A implementação fica em `v3/internal/assetserver/`, com uma separação clara entre os arquivos:

```
build_dev.go              # ⬅️ dev-only entrypoint (!production build tag)
build_production.go       # ⬅️ production-only entrypoint (production build tag)
assetserver.go            # Shared core
assetserver_dev.go        # Dev proxy/disk handler
assetserver_webview.go    # WebView-side adapter
assetserver_darwin.go     # OS-specific helpers (also linux/windows variants)
asset_fileserver.go       # Shared static file logic
content_type_sniffer.go   # MIME type detection
mimecache.go              # Cached MIME lookups
ringqueue.go              # Tiny in-memory LRU
options.go                # Configuration struct
middleware.go             # http.Handler middleware type
bundled_assetserver.go    # Hand-written wrapper around embedded bundles
bundledassets/            # Embedded runtime JS assets
```

---

## Modo de desenvolvimento

### Ciclo de vida

1. `wails3 dev` inicializa e **inicia o servidor de desenvolvimento do seu frontend** (Vite, SvelteKit, React-SWC etc.) executando a tarefa definida em `build/Taskfile.yml` (normalmente `npm run dev`).
2. A CLI define `WAILS_VITE_PORT` como a porta de desenvolvimento do Wails e `FRONTEND_DEVSERVER_URL` como a URL **completa** (`http://host:port` / `https://host:port`) que aponta para o servidor de desenvolvimento em execução do framework. Consulte `internal/commands/dev.go`.
3. O servidor de ativos de desenvolvimento (incluído na compilação por meio de `//go:build !production` em `build_dev.go`) lê `FRONTEND_DEVSERVER_URL` por meio de `GetDevServerURL()` e encaminha para ele, via proxy reverso, o tráfego que não pertence ao runtime.
4. Os arquivos estáticos (`/assets/logo.svg`) podem ser **servidos diretamente do disco** por meio de `asset_fileserver.go` (para maior velocidade), enquanto tudo o que não for reconhecido é **encaminhado por proxy** ao servidor de desenvolvimento do framework, proporcionando substituição *instantânea* de módulos a quente.

```
┌─────────┐  /wails/runtime.js     ┌─────────────┐
│ Browser │ ── embedded runtime ──▶│   Runtime   │
├─────────┤                        └─────────────┘
│   JS    │  / (index.html)        proxy / -> Vite via FRONTEND_DEVSERVER_URL
└─────────┘ ◀─────────────┐
              AssetServer │
                          ▼
                   ┌────────────┐
                   │  Vite Dev  │
                   │   Server   │
                   └────────────┘
```

### Recursos

- **Recarregamento automático** — Vite, SvelteKit etc. injetam HMR por WebSocket; o servidor de ativos de desenvolvimento o encaminha de forma transparente.
- **Suporte a mapas de código-fonte** — como os ativos não são empacotados, as ferramentas de desenvolvimento do navegador associam os erros ao código-fonte original.
- **Sem recompilação do Go** — apenas o frontend é recompilado; o código Go continua em execução até que você altere arquivos `.go`.

### Troca de frameworks

O proxy de desenvolvimento é **independente de framework**. A CLI do Wails disponibiliza duas variáveis de ambiente ao iniciar sua tarefa de desenvolvimento:

| Variável de ambiente | Origem | Significado |
| --- | --- | --- |
| `WAILS_VITE_PORT` | `internal/commands/dev.go` (constante `wailsVitePort`) | Porta de desenvolvimento padrão (9245, a menos que `--port` seja passado) — sua configuração do Vite deve respeitar esse valor |
| `FRONTEND_DEVSERVER_URL` | `internal/commands/dev.go` | URL completa para a qual o Wails encaminhará as solicitações por proxy; lida no Go por meio de `assetserver.GetDevServerURL()` (`build_dev.go`) |

Não há nenhuma variável de ambiente `VITE_PORT`, `FRONTEND_DEV_PORT` nem `WAILSDEV_VERBOSE` na árvore da v3.

Adicione um novo modelo → defina sua tarefa de desenvolvimento → o servidor de ativos simplesmente funciona.

---

## Modo de produção

Quando você executa `wails3 build`, o pipeline:

1. Executa o **build de produção** do frontend (`npm run build`), gerando `frontend/dist/**`.
2. **Incorpora** esse diretório ao aplicativo por meio de `go:embed` no próprio pacote do aplicativo (normalmente `//go:embed all:frontend/dist` ao lado de `main.go`).
3. Compila o binário Go com `-tags production` (repassado por meio de `EXTRA_TAGS` pelo wrapper do Taskfile).

`internal/assetserver/build_production.go` é o stub de tag de build que ativa o caminho de código de produção. `internal/assetserver/bundled_assetserver.go` é **escrito manualmente** — ele encapsula o JavaScript do runtime em `bundledassets/` e não é um arquivo gerado.

### Tratamento de solicitações

O handler efetivo é `internal/assetserver/assetserver.go` / `asset_fileserver.go`. Conceitualmente:

1. Tenta servir os ativos estáticos incorporados no caminho solicitado.
2. Se isso falhar, usa `index.html` para o roteamento da SPA.
3. Detecta o tipo de conteúdo se a extensão for desconhecida (`content_type_sniffer.go`).
4. Define cabeçalhos de cache adequados.

- **Detecção de MIME** — para arquivos sem extensão, o tipo de conteúdo é detectado a partir dos primeiros ~512 bytes (`content_type_sniffer.go`), e o resultado é armazenado em cache em `mimecache.go` / `ringqueue.go`.
- **Cabeçalhos de segurança** — impede a navegação `file://` e define `nosniff`.

Como tudo é incorporado, o binário distribuído **não tem dependências externas** (nem mesmo no Windows).

---

## Integração entre desenvolvimento ↔ produção

Do ponto de vista de `pkg/application`, ambos os modos expõem a **mesma interface pública**: uma struct `AssetOptions` com um `Handler http.Handler`, além de middlewares e da integração do ciclo de vida em `internal/assetserver/`. A alternância entre desenvolvimento e produção ocorre inteiramente por meio de tags de build do Go, portanto o código do aplicativo é idêntico nos dois modos.

```go
app := application.New(application.Options{
    Assets: application.AssetOptions{
        Handler: application.AssetFileServerFS(assetsFS),
    },
})
```

---

## Como os frameworks de frontend se integram

### Modelos

Cada modelo distribuído (React, Vue, Svelte, Solid, Vanilla etc.) contém:

- `build/Taskfile.yml`
- `frontend/vite.config.ts` (ou equivalente)

A configuração do Vite ou equivalente lê `WAILS_VITE_PORT` e vincula o servidor de desenvolvimento a essa porta. Em seguida, a CLI publica o `FRONTEND_DEVSERVER_URL` real para ser usado pelo proxy no aplicativo.

Os frameworks permanecem totalmente desacoplados do Go:

- Não é necessário importar nenhum SDK JavaScript do Wails durante a compilação — `/wails/runtime.js` é servido pelo servidor de ativos em tempo de execução.
- Qualquer framework com um servidor HTTP de desenvolvimento pode ser integrado.

---

## Extensão e personalização

Precisa de cabeçalhos personalizados, autenticação ou gzip?

1. Defina um `middleware.Middleware` (alias de `func(http.Handler) http.Handler`, declarado em `internal/assetserver/middleware.go`).
2. Conecte-o ao seu `application.AssetOptions` por meio da configuração exposta por `internal/assetserver/options.go`.
3. O comportamento é idêntico nos modos de desenvolvimento e produção — não há uma lista de middlewares específica para cada modo.

---

## Principais arquivos de código-fonte

| Arquivo | Função |
| --- | --- |
| `build_dev.go` / `build_production.go` | Wrappers com tags de compilação que selecionam desenvolvimento ou produção |
| `assetserver.go` / `asset_fileserver.go` | Handler HTTP principal |
| `assetserver_dev.go` | Proxy reverso para `FRONTEND_DEVSERVER_URL` |
| `bundled_assetserver.go` | Wrapper escrito manualmente em torno de `bundledassets/` |
| `options.go` | Configuração voltada para `application.AssetOptions` |
| `mimecache.go` / `ringqueue.go` | Cache de MIME e pequeno cache LRU |

---

## Armadilhas e depuração

- **Tela branca em produção** — geralmente é causada pelo roteamento da SPA: verifique se o servidor de desenvolvimento serve `index.html` para caminhos desconhecidos e se o fallback do handler de produção incorporado é alcançado.
- **404 no desenvolvimento** — sua configuração do Vite não está vinculando a `WAILS_VITE_PORT`, ou a CLI não conseguiu acessar o servidor de desenvolvimento para preencher `FRONTEND_DEVSERVER_URL`.
- **Ativos grandes** — incorporá-los aumenta o tamanho do binário. Sirva mídias grandes a partir de uma origem separada ou transmita-as por streaming usando um `http.Handler` personalizado.

---

Agora você sabe como o **servidor de ativos** do Wails fornece seu código web à janela nativa tanto em **desenvolvimento** quanto em **produção**. Domine essa camada para depurar problemas de carregamento, adicionar middlewares ou até mesmo substituir toda a cadeia de ferramentas de frontend por outra com confiança.
