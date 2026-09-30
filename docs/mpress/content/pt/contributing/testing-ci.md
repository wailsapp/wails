---
title: "Testes e integração contínua"
description: "Como o Wails v3 garante a qualidade com testes unitários, suítes de integração, detecção de condições de corrida e CI do GitHub Actions."
slug: "contributing/testing-ci"
sourcePath: "contributing/testing-ci.md"
---

Frameworks robustos para desktop exigem testes extremamente confiáveis. O Wails v3 emprega uma **estratégia em camadas**:

| Camada | Objetivo | Ferramentas |
| --- | --- | --- |
| Testes unitários | Feedback rápido para funções isoladas | `go test ./...` |
| Testes do gerador e da CLI | Validar `wails3 generate bindings` e a infraestrutura da CLI | `task test:generator`, `task test:cli` |
| Testes de templates | Garantir que todos os templates distribuídos continuem sendo compilados | `task test:templates` |
| Detecção de condições de corrida | Detectar condições de corrida de dados no runtime e na ponte | `go test -race ./...` |
| Matriz de CI | Confiança multiplataforma em cada PR | GitHub Actions |

Este documento explica **onde ficam os testes**, **como executá-los** e **o que o Taskfile orquestra**.

> As orientações sobre condições de corrida mencionadas em versões anteriores como `pkg/application/RACE.md`
>
> estão atualmente em `v3/TESTING.md`.

---

## 1. Convenções de diretórios

```
v3/
├── internal/.../_test.go     # Unit tests for internal packages
├── pkg/.../_test.go          # Public API tests
├── tasks/events/generate.go  # Code generator for event constants (NOT a test harness)
├── tests/                    # Top-level integration test harness
└── TESTING.md                # Race / Cgo testing guidance
```

Diretrizes:

- **Mantenha os testes unitários junto ao código** (`foo.go` ↔ `foo_test.go`).
- Use o **estilo de caixa-preta** para pacotes `pkg/` (`package application_test`) quando isso melhorar a organização da API.
- Os fixtures compartilhados ficam onde são usados (não há um pacote `internal/testutil/` central nesta árvore — em vez disso, conecte os helpers em cada pacote).

---

## 2. Testes unitários

### Como escrever testes

```go
func TestEventConstants(t *testing.T) {
    assert.NotEmpty(t, events.Common.WindowFocus)
}
```

Recomendações:

- Use [`stretchr/testify`](https://github.com/stretchr/testify) — já incluído em `go.mod`.
- Prefira testes **orientados por tabela** sempre que várias entradas, casos extremos ou resultados esperados exercitarem o mesmo comportamento. Dê um nome descritivo a cada caso.
- Quando necessário, use stubs para substituir comportamentos específicos de plataforma, condicionando-os a tags de compilação (`foo_windows_test.go`, `foo_darwin_test.go`, …).

### Cobertura esperada

Espera-se que a lógica nova e alterada tenha 100% de cobertura de instruções Go. Meça o pacote que você alterou em vez de depender de uma porcentagem para todo o repositório:

```bash
go test ./pkg/application -coverprofile=coverage.out
go tool cover -func=coverage.out
```

Alguns caminhos não podem ser testados de maneira razoável em um ambiente de testes normal — por exemplo, falhas exclusivas de determinada plataforma, comportamentos dependentes de hardware ou um fallback defensivo que não possa ser induzido com segurança. Mantenha essas exceções restritas e explique cada caminho não coberto na descrição do PR.

### Execução local

```bash
cd v3
go test ./... -cover
```

Ou por meio do Taskfile (alvos reais — não há um atalho `task test`):

```
task test:cli            # CLI plumbing tests
task test:generator      # bindings generator round-trip tests
task test:templates      # build every shipped template
task test:infrastructure # supporting helpers
task test:examples       # exercise the example matrix (downloads as needed)
task test:all            # everything above
task sanity              # quick smoke check (also: sanity:gtk4)
task precommit           # what you should run before pushing
```

---

## 3. Testes de integração

`v3/tests/` hospeda a infraestrutura de integração entre pacotes. Os alvos `test:example:*` e `test:examples:*` do Taskfile executam verificações de compilação e inicialização em darwin / windows / linux (incluindo matrizes GTK3 / GTK4 baseadas em Docker no Linux).

> Os exemplos executáveis ficam em `v3/examples/`. Os alvos de teste selecionam e compilam
>
> os exemplos adequados à plataforma host ou à matriz de CI.

Execute a suíte de testes de fumaça da plataforma host com:

```
task test:examples       # host
task test:examples:all   # full matrix (slow)
```

---

## 4. Detecção de condições de corrida

Condições de corrida de dados são fatais em runtimes de GUI.

### Guia sobre condições de corrida

Consulte `v3/TESTING.md` para saber sobre:

- Condições de corrida benignas conhecidas e justificativa para suprimi-las
- Como interpretar rastreamentos de pilha que atravessam limites do Cgo (GTK no Linux + WebKit2GTK)

### Suíte local de detecção de condições de corrida

```
go test -race ./...
```

> `wails3 dev` não tem a opção `-race` — suas opções de CLI são `--config`, `--port`,
>
> e `-s` (habilita HTTPS). Para exercitar o runtime com o detector de condições de corrida,
>
> compile um aplicativo de teste com `go build -race` e execute-o diretamente.

---

## 5. Fluxos de trabalho do GitHub Actions

Arquivos reais de fluxo de trabalho em `.github/workflows/` (verificados em relação à árvore):

| Arquivo | Finalidade |
| --- | --- |
| `build-and-test-v3.yml` | Matriz principal de compilação e testes do v3. Usa `actions/setup-go@v5` com `go-version: 1.25`. Executa `task runtime:check`, `task runtime:test`, `task runtime:build`, `task test:examples` (e `BUILD_TAGS=gtk4 task test:examples` para o caminho GTK4), `task generator:test:check`, `task install` e, em seguida, `wails3 build` como verificação de fumaça. Os jobs do Linux instalam `libgtk-3-dev libwebkit2gtk-4.1-dev libwayland-dev build-essential pkg-config xvfb x11-xserver-utils at-spi2-core xdg-desktop-portal-gtk` e executam a suíte de testes em `dbus-run-session -- xvfb-run`. |
| `cross-compile-test-v3.yml` | Verificações básicas de compilação cruzada |
| `auto-changelog-v3.yml`, `changelog-v3.yml` | Automação do changelog |
| `nightly-release-v3.yml` | Artefatos da versão nightly do v3 |
| `bump-webview2-v3.yml`, `release-webview2.yml` | Gerenciamento de dependências e versões do WebView2 |
| `build-cross-image.yml` | Compila a imagem de contêiner do compilador cruzado |
| `publish-npm.yml` | Publica no npm o runtime JS `@wailsio/runtime` incorporado |
| `pr-master.yml` | Verificações de PR em relação à branch `master` |
| `semgrep.yml` | Análise estática com Semgrep |
| `stale-issues.yml`, `issue-labeler.yml`, `file-labeler.yml`, `claude.yml`, `generate-sponsor-image.yml`, `sync-translated-documents.yml`, `upload-source-documents.yml`, `build-and-test.yml`, `weekly-release-v2.yml` | Fluxos de manutenção do repositório e relacionados à v2 |

Não há **nenhum** `qodana.yaml` nem **nenhum** `runtime.yml` nesta árvore — versões anteriores desta página mencionavam ambos, mas somente `semgrep.yml` abrange a análise estática, e o pacote JS do runtime é publicado por meio de `publish-npm.yml`.

As etapas de CI correspondem aos alvos do Taskfile acima (`task test:cli`, `task test:generator`, `task test:templates`, `task test:examples`, …), portanto você pode reproduzir localmente cada etapa da CI. A etapa de smoke test `wails3 build` em `build-and-test-v3.yml` é invocada **sem flags adicionais** — não existe a flag `-skip-package` em `wails3 build`.

---

## 6. Paridade entre verificações locais e CI

Não há um único alvo geral `task ci`. Reproduza a CI encadeando os alvos reais:

```
task precommit
task test:cli
task test:generator
task test:templates
task test:examples
```

---

## 7. Solução de problemas em testes com falha

| Sintoma | Causa provável | Correção |
| --- | --- | --- |
| **Condição de corrida em `webview_window_darwin.go`** | Alteração do estado da janela fora da thread principal | Encaminhe por meio de `application.InvokeAsync` / `Invoke` para que a chamada seja executada na thread principal |
| **O teste no Linux trava em uma CI sem interface gráfica** | O GTK precisa de um servidor gráfico | Execute sob `xvfb-run`, por exemplo, `xvfb-run task test:examples:linux` |
| **A compilação do template falha** | O lockfile do frontend está desatualizado | Execute `wails3 init` novamente em um diretório limpo para atualizar o template |
| **Erros de Coverpkg** | Teste de integração importando `main` | Mude para a tag de compilação `//go:build integration` e condicione a importação a ela |

---

## 8. Adição de novos testes

1. **Unitário** — crie `*_test.go` e execute `go test ./...`
2. **Gerador / CLI** — amplie os casos em `internal/generator/testcases/` ou `internal/commands/*_test.go` e execute novamente `task test:generator` / `task test:cli`
3. **Templates / exemplos** — certifique-se de que os templates distribuídos ainda sejam compilados com `task test:templates`

---

## 9. Mapa dos principais arquivos

| O quê | Caminho |
| --- | --- |
| Teste de ida e volta do gerador | `internal/generator/generate_test.go` |
| Teste de compilação de ativos | `internal/commands/build-assets_test.go` |
| Guia de condições de corrida / Cgo | `v3/TESTING.md` |
| Alvos de teste do Taskfile | `v3/Taskfile.yaml` |
| Gerador de constantes de eventos | `v3/tasks/events/generate.go` |
| Fluxo de trabalho de CI | `.github/workflows/build-and-test-v3.yml` (Go 1.25 por meio de `actions/setup-go@v5`) |
| Análise estática | `.github/workflows/semgrep.yml` |
| Publicação do runtime no npm | `.github/workflows/publish-npm.yml` |

---

A qualidade não é deixada para depois no Wails v3. Com testes unitários, suítes de geradores e templates, detecção de condições de corrida e uma matriz de CI multiplataforma, você pode contribuir com confiança, sabendo que suas alterações passam em todos os sistemas operacionais compatíveis. Bons testes!
