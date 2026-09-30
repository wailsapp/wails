---
title: "Pipeline de build e empacotamento"
description: "O que acontece nos bastidores quando você executa `wails3 build`, como são produzidos binários multiplataforma e como são gerados instaladores para cada sistema operacional."
slug: "contributing/build-packaging"
sourcePath: "contributing/build-packaging.md"
---

`wails3 build` é intencionalmente **simples**: trata-se de um wrapper de Taskfile que encaminha tags de build adicionais para a tarefa `build` do projeto host. O trabalho pesado fica no próprio `build/Taskfile.yml` do projeto (gerado por `wails3 init`), em `internal/commands/build-assets.go` (que gerencia os ativos incorporados durante o build) e em `internal/packager` (empacotamento nfpm para Linux), além de `internal/commands/appimage.go`, `internal/commands/msix.go`, `internal/commands/dmg/dmg.go`, e `internal/commands/dot_desktop.go` (instaladores específicos de cada plataforma).

Esta página aborda:

1. Os pontos de entrada reais da CLI
2. O fluxo de build orientado pelo Taskfile
3. A incorporação de ativos e a injeção de informações de build
4. Os back-ends de empacotamento de cada plataforma
5. A personalização do pipeline
6. Solução de problemas

---

## 1. Pontos de entrada reais da CLI

```
wails3 build       → internal/commands.Build       (in task_wrapper.go)
wails3 package     → internal/commands.Package     (in task_wrapper.go)
wails3 generate build-assets → GenerateBuildAssets (in build-assets.go)
wails3 update build-assets   → UpdateBuildAssets   (in build-assets.go)
wails3 tool buildinfo        → BuildInfoOptions    (in tool_buildinfo.go)
wails3 tool package          → internal/packager   (nfpm wrapper)
wails3 generate .desktop     → in dot_desktop.go
```

`internal/commands/task_wrapper.go`:

```go
func Build(buildFlags *flags.Build, otherArgs []string) error {
    // forwards --tags / EXTRA_TAGS, then defers to a Taskfile target
    return wrapTask("build", otherArgs)
}
```

`flags.Build` expõe **apenas** uma opção — `--tags` (encaminhada como `EXTRA_TAGS=`). `wails3 build` **não** tem as opções `-platform`, `-o`, `-skipbindings`, `-skip-package`, `-package`, `-ldflags`, `-verbose`, `-debug`, `-devbuild`, `-icon` ou `-clean`. A compilação cruzada, os caminhos de saída, os ícones etc. são configurados em **`Taskfile.yml`**, **`build/config.yml`** e nos comandos auxiliares `wails3 generate icons` / `wails3 generate build-assets`.

`build/build.json` **não** faz parte da v3 — a configuração usa `Taskfile.yml` mais `build/config.yml`.

---

## 2. Fluxo de build orientado pelo Taskfile

Um projeto recém-inicializado inclui um `build/Taskfile.yml` com aproximadamente os seguintes namespaces:

| Namespace | Tarefas (seleção) |
| --- | --- |
| `darwin:` | `build`, `build:universal`, `package`, `run`, `dev` |
| `windows:` | `build`, `package`, `run`, `dev` |
| `linux:` | `build`, `package`, `run`, `dev` |
| `common:` | `update:build-assets`, `generate:icons`, `generate:syso` |

Por padrão, `wails3 build` invoca o namespace `build` do sistema operacional host; por sua vez, o Taskfile do projeto executa `go build` no shell com as opções específicas do host. Para fazer o build para outro sistema operacional, execute diretamente a tarefa correspondente (por exemplo, `wails3 task darwin:build:universal`), em vez de passar uma opção para `wails3 build`.

O diretório de saída padrão é **`bin/<APP_NAME>`** (sem o prefixo `build/bin/`).

---

## 3. Ativos incorporados durante o build e informações de build

| Aspecto | Arquivo |
| --- | --- |
| Geração/atualização de ativos de build | `internal/commands/build-assets.go` |
| Exibição de informações de build (CLI: `wails3 tool buildinfo`) | `internal/commands/tool_buildinfo.go` — exibe informações; **não** é um injetor de `ldflags` |
| Stub de produção | `internal/assetserver/build_production.go` — `//go:build production` |
| Bundles do frontend | incorporados por meio de `//go:embed` no próprio pacote da aplicação (por exemplo, ao lado de `main.go`) |
| Recursos do Windows (`.syso`) | `internal/commands/syso.go` — gera `rsrc_windows_<arch>.syso` |
| MSIX do Windows | `internal/commands/msix.go` + `internal/commands/webview2/` |
| Entradas de DMG do macOS | `internal/commands/dmg/` |
| `.desktop` do Linux | `internal/commands/dot_desktop.go` |

A CLI não gera automaticamente `bundled_assetserver.go` para sua aplicação — `internal/assetserver/bundled_assetserver.go` é **escrito manualmente** e encapsula o runtime JS incorporado em `bundledassets/`.

---

## 4. Back-ends de empacotamento

### Linux

O empacotamento para Linux é gerenciado pelo **nfpm** (e não por `fpm`):

- `internal/packager/packager.go` encapsula `github.com/goreleaser/nfpm/v2` e expõe `CreatePackageFromConfig(pkgType, configPath, output)` / `CreatePackageFromConfigWriter(...)`.
- Os projetos gerados incluem as configurações `myapp.DEB`, `myapp.RPM` e `myapp.ARCHLINUX` no estilo do nfpm em `internal/commands/` (usadas por `wails3 tool package`).
- A geração de AppImage fica em `internal/commands/appimage.go`, que chama `linuxdeploy` + `linuxdeploy-plugin-gtk` (o plugin é fornecido em `internal/commands/linuxdeploy-plugin-gtk.sh`).

`wails3 build` **não** tem uma opção `-package deb`/`rpm`. Use `wails3 tool package` ou o destino do Taskfile específico da plataforma.

### macOS

- `darwin:package` do Taskfile do projeto produz o bundle `.app`.
- Os recursos do DMG ficam em `internal/commands/dmg/`; um projeto pode empacotar o bundle em um DMG com `hdiutil` depois que `darwin:package` terminar (o Taskfile dos modelos mais recentes inclui um auxiliar `dmg`).
- Os identificadores do CFBundle, a versão e os direitos autorais vêm das opções de `-product*` durante `wails3 init` e de `build/config.yml`.

### Windows

- O empacotamento para Windows usa **MSIX** (não WiX/MSI). Consulte `internal/commands/msix.go` e `internal/commands/webview2/` para ver o fluxo de trabalho completo.
- **Não** há `internal/commands/packager.go` e **não** há diretório `internal/commands/windows_resources/`.
- A assinatura de código opcional do executável é realizada por meio de `wails3 tool sign` (Authenticode) — consulte `internal/commands/sign.go`.

---

## 5. Personalização do pipeline

| Necessidade | Abordagem |
| --- | --- |
| Tags de compilação adicionais | `wails3 build --tags myFeature,otherTag` |
| Linter/etapa de pré-compilação | Adicione uma tarefa a `build/Taskfile.yml` e faça a tarefa `build` específica do sistema operacional depender dela |
| Compilação cruzada | Execute a tarefa pertinente ao sistema operacional (por exemplo, `wails3 task linux:build`) — não há uma opção `-platform` |
| Ignorar o empacotamento | Execute apenas a tarefa `build`; `package` é separado |
| Empacotador personalizado | Coloque uma configuração em `internal/commands/myapp.*` e chame `wails3 tool package` com `-config <file>` |
| Remover símbolos | Edite a tarefa `darwin:/windows:/linux:` `build` para passar `-ldflags "-s -w"` diretamente a `go build` — o próprio `wails3 build` não tem uma opção `-ldflags` |

Todos os alvos do Taskfile respeitam as variáveis de ambiente publicadas pelo Wails (`APP_NAME`, `WAILS_VITE_PORT`, `FRONTEND_DEVSERVER_URL`, …), portanto as tarefas personalizadas podem utilizá-las.

---

## 6. Solução de problemas

| Sintoma | Causa provável | Correção |
| --- | --- | --- |
| **`ld: framework not found WebKit` (mac)** | As ferramentas de linha de comando do Xcode estão ausentes | `xcode-select --install` |
| **Janela em branco na compilação de produção** | Falha na compilação do frontend ou no roteamento da SPA | Verifique se `frontend/dist/index.html` existe e se o manipulador de recursos recorre a ele como alternativa |
| **As ferramentas de empacotamento MSIX estão ausentes** | SDK do `WebView2` / ferramentas MSIX não instalados | Execute `wails3 task install:msix:tools` |
| **`linuxdeploy` não encontrado** | O plugin não está no PATH | Instale `linuxdeploy` e execute `internal/commands/linuxdeploy-plugin-gtk.sh` por meio da etapa de instalação automática da CLI |

`wails3 build` não tem uma opção `-verbose`. Defina `TASK_X_VERBOSE=1` (Taskfile) ou inspecione diretamente o alvo da tarefa para ver os comandos executados.

---

## 7. Mapa dos principais arquivos-fonte

| Área | Arquivo |
| --- | --- |
| Wrapper de compilação | `internal/commands/task_wrapper.go` (`Build`, `Package`, `SignWrapper`, `wrapTask`) |
| Geração de recursos de compilação | `internal/commands/build-assets.go` (`GenerateBuildAssets`, `UpdateBuildAssets`) |
| Exibidor de informações de compilação | `internal/commands/tool_buildinfo.go` |
| Gerador de AppImage | `internal/commands/appimage.go` |
| Empacotamento para Linux (nfpm) | `internal/packager/packager.go`, `internal/commands/myapp.{DEB,RPM,ARCHLINUX}` |
| MSIX para Windows | `internal/commands/msix.go`, `internal/commands/webview2/` |
| Gerador de recursos do Windows | `internal/commands/syso.go` |
| Recursos de DMG do macOS | `internal/commands/dmg/` |
| Gerador de `.desktop` | `internal/commands/dot_desktop.go` |
| Constantes de versão | `internal/version/version.go` |

Mantenha esta tabela à mão ao rastrear uma falha de compilação.

---

Agora você tem a visão completa, do **código-fonte** ao **instalador**. Em resumo, o próprio `wails3 build` é apenas um wrapper leve — quase toda personalização ocorre no `Taskfile.yml`/`build/config.yml` do projeto ou por meio dos subcomandos explícitos `wails3 generate …`/`wails3 tool …`. Boa distribuição!
