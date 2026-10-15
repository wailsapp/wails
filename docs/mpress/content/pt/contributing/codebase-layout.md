---
title: "Estrutura do código-fonte"
description: "Como o repositório do Wails v3 é organizado e como suas partes se integram"
slug: "contributing/codebase-layout"
sourcePath: "contributing/codebase-layout.md"
---

O Wails v3 fica em um **monorepositório** que contém o runtime do framework, a CLI, exemplos, documentação e a cadeia de ferramentas de build. Esta página apresenta a *estrutura de diretórios* relevante para quem deseja explorar os componentes internos.

## Visão geral

```
wails/
├── v3/               # ⬅️ Everything specific to Wails v3 lives here
├── v2/               # Legacy v2 implementation (can be ignored for v3 work)
├── docs/             # M-Press-powered v3 docs site (this page!)
├── website/          # Docusaurus v2 site and marketing pages (main site)
├── scripts/          # Misc helper scripts (e.g. sponsor image generator)
└── *.md              # Project-wide meta files (CHANGELOG, LICENSE, …)
```

Daqui em diante, examinaremos mais de perto a árvore **`v3/`**.

## Raiz de `v3/`

```
v3/
├── cmd/          # Compilable commands (currently only the wails3 CLI)
├── internal/     # Framework implementation (not public API)
├── pkg/          # Public Go packages — the API surface
├── tasks/        # Taskfile-based release / generation utilities
├── wep/          # RFC-style proposals (Wails Enhancement Proposals)
├── tests/        # Integration test harness
├── go.mod
└── go.sum
```

> Os templates de projeto são fornecidos em `internal/templates/` (uma pasta para cada stack de framework
>
> mais `base/`, `_common/` e `ios/`). Não há nenhum diretório `v3/templates/`
>
> no nível superior.

### Modelo mental

1. **`pkg/`** expõe *o que os desenvolvedores de aplicações importam*\
2. **`internal/`** contém *como a mágica é implementada*\
3. **`cmd/wails3`** controla *o ciclo de vida e os builds do projeto*\

Todo o restante dá suporte a esses três pilares.

---

## `cmd/` – Comandos

| Caminho | Observações |
| --- | --- |
| `v3/cmd/wails3` | O **ponto de entrada da CLI**. Um pequeno `main.go` delega toda a lógica aos pacotes em `internal/commands`. |
| `internal/commands/*` | Subcomandos (init, dev, build, doctor, …). Cada um fica em seu próprio arquivo para facilitar sua localização. |
| `internal/commands/task_wrapper.go` | Faz a integração entre as flags da CLI e o pipeline de build do Taskfile. |

A CLI é responsável por:

- **Criação da estrutura do projeto** (`init`, geração de templates)\
- **Orquestração do servidor de desenvolvimento** (`dev`, recarregamento em tempo real)\
- **Builds de produção e empacotamento** (`build`, `package`, wrappers específicos de plataforma)\
- **Diagnósticos** (`doctor`)\

---

## `internal/` – A sala de máquinas

```
internal/
├── assetserver/  # Serving & embedding web assets
├── buildinfo/    # Reproducible build metadata
├── commands/     # CLI mechanics (see above)
├── runtime/      # Build-tag glue + embedded JS runtime sources
├── generator/    # Static analysis & binding generator
├── templates/    # Project templates (frontend stacks)
├── packager/     # nfpm wrapper used by `wails3 tool package`
├── capabilities/ # Host OS capability probing
├── dbus/         # Generic D-Bus helper
├── service/      # Service-template scaffolding (`wails3 generate service`)
└── ...           # [other helper sub-packages: flags, hash, term, …]
```

### Principais subpacotes

| Pacote | Responsabilidade | Onde se conecta |
| --- | --- | --- |
| `runtime` | Abriga o pequeno código de integração com build tags em `runtime{,_darwin,_linux,_windows,_android,_dev,_prod}.go`, além do runtime JS incorporado em `runtime/desktop/`. O código propriamente dito de janela, área de transferência, caixa de diálogo e bandeja para cada sistema operacional fica em `pkg/application/*_{darwin,linux,windows}.go`. | Importado indiretamente por meio de `pkg/application`. |
| `assetserver` | Servidor de arquivos com dois modos:<br />• Desenvolvimento: serve a partir do disco e atua como proxy do Vite (`build_dev.go`)<br />• Produção: incorpora recursos por meio de `go:embed` (`build_production.go`) | Inicializado por `pkg/application` durante a inicialização. |
| `generator` | Analisa o código-fonte Go para criar os **metadados de bindings**, que posteriormente geram arquivos stub TypeScript/JS e constantes de eventos. Pontos de entrada: `generator.Generate` / `generator.Generator` sobre `collect/` + `render/`. | Acionado por `wails3 generate bindings`. |
| `packager` | Wrapper de `nfpm` usado para gerar artefatos Linux `deb`/`rpm`/`archlinux` (controlado pelas configurações nfpm `myapp.DEB`/`.RPM`/`.ARCHLINUX` em `internal/commands/`). | Invocado por `wails3 tool package`. O DMG do macOS e o MSIX do Windows ficam em `internal/commands/{dmg,msix.go,webview2/}`. |

Os utilitários de suporte (por exemplo, `s/`, `hash/` e `flags/`) mantêm as responsabilidades internas desacopladas.

---

## `pkg/` – API pública

```
pkg/
├── application/  # Core API: App, windows, menus, dialogs, events, managers
├── events/       # Event constants (Common/Mac/Windows/Linux) + generator
├── services/     # Optional built-in services (notifications, kvstore, …)
├── doctor-ng/    # New-style `wails3 doctor-ng` checks
├── errs/         # Shared error types
├── icons/        # Default platform icons
├── mac/          # macOS-only helpers
└── w32/          # Windows Win32 helpers
```

> Não há nenhum pacote `pkg/runtime/`, `pkg/options/` ou `pkg/menu/`. As opções de janela/menu
>
> ficam junto de `pkg/application` (por exemplo, `WebviewWindowOptions`, `Menu`,
>
> `MenuItem`), e `assetserver/` fica em `internal/`.

`pkg/application` inicializa um programa Wails:

```go
func main() {
    app := application.New(application.Options{
        Name: "MyApp",
        Assets: application.AssetOptions{
            Handler: application.AssetFileServerFS(assetsFS),
        },
    })
    window := app.Window.New()
    window.SetTitle("Hello").SetSize(1024, 768)
    _ = app.Run()
}
```

Nos bastidores, ele:

1. Conecta o código de integração com build tags em `internal/runtime` ao código específico de cada sistema operacional em `pkg/application/`
2. Configura uma instância de `internal/assetserver`
3. Registra todos os processadores de mensagens orientados por bindings
4. Entra na thread principal do sistema operacional

---

## `internal/templates/` – Modelos para criação da estrutura

`internal/templates/` fornece **templates básicos** (estrutura Go em `base/`, `_common/`, `ios/`) e **temas de frontend** (`vanilla[-ts]`, `react[-ts]`, `react-swc[-ts]`, `lit[-ts]`, `preact[-ts]`, `qwik[-ts]`, `solid[-ts]`, `svelte[-ts]`, `sveltekit[-ts]`, `vue[-ts]`).

Em `wails3 init -t react`, a CLI:

1. Copia os arquivos Go de `_common`
2. Mescla o pacote de frontend desejado
3. Executa `go mod tidy` (pode ser ignorado com `--skipgomodtidy`)

Editar os modelos **não** afeta os aplicativos existentes, apenas os futuros `init`s. Os exemplos públicos ficam em `v3/examples/`; eles não substituem as suítes de testes automatizados descritas na documentação para colaboradores.

---

## `tasks/` – Automação de versões

Os Taskfiles encapsulam processos complexos de compilação cruzada, atualização de versão e geração do changelog. Eles são consumidos programaticamente por `internal/commands/task.go`, de modo que a mesma lógica seja usada pela **CLI** e pela **CI**.

---

## Como as partes interagem

```d2
direction: down
CLI: CLI wails3
Generator: internal/generator
AssetDev: assetserver (desenvolvimento)
Packager: internal/packager
AppRuntime: {
  label: Runtime do aplicativo
  ApplicationPkg: pkg.application
  InternalRuntime: internal.runtime
  OSAPIs: APIs do SO
}
CLI -> Generator: compilar / gerar
CLI -> AssetDev: desenvolvimento
CLI -> Packager: empacotar
Generator -> ApplicationPkg: vinculações
ApplicationPkg -> InternalRuntime
InternalRuntime -> OSAPIs
ApplicationPkg -> AssetDev
ApplicationPkg.label: ApplicationPkg
InternalRuntime.label: InternalRuntime
OSAPIs.label: OSAPIs
```

*CLI → gerador → runtime* forma o caminho principal do **código-fonte** até o **aplicativo desktop em execução**.

---

## Dicas de orientação

| Precisa entender… | Consulte… |
| --- | --- |
| Camadas de compatibilidade de plataforma | `pkg/application/*_darwin.go`, `*_linux.go`, `*_windows.go` (janela, área de transferência, caixas de diálogo, bandeja do sistema, thread principal, events_common). cgo no Linux: `pkg/application/linux_cgo*.go`. |
| Protocolo da ponte | `pkg/application/messageprocessor*.go` |
| Fluxo de trabalho de recursos | `internal/assetserver/` (`build_dev.go` versus `build_production.go`) |
| Fluxo de empacotamento | `internal/commands/{appimage,msix,dot_desktop,dmg/}.go`, `internal/packager/` |
| Mecanismo de modelos | `internal/templates/` (`templates.Install`, `templates.GetDefaultTemplates`) |
| Análise estática | `internal/generator/{generate.go,collect/,render/}` |

---

Agora você tem um **mapa mental** do repositório. Use-o com `ripgrep`, os comandos “Go to File/Symbol” do seu IDE e os aplicativos de exemplo para explorar mais a fundo qualquer recurso. Bom desenvolvimento!
