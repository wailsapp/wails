---
title: "Testes e integração contínua"
description: "Como o Wails v3 garante a qualidade com testes unitários, suítes de integração, detecção de condições de corrida e CI do GitHub Actions."
slug: "contributing/testing-ci"
sourcePath: "contributing/testing-ci.md"
---

Frameworks robustos para desktop exigem testes extremamente confiáveis. O Wails v3 adota uma **estratégia em camadas**:

| Camada | Objetivo | Ferramentas |
| --- | --- | --- |
| Testes unitários | Feedback rápido para funções isoladas | `go test ./...` |
| Testes do gerador e da CLI | Validar `wails3 generate bindings` e a infraestrutura da CLI | `task test:generator`, `task test:cli` |
| Testes de templates | Garantir que todos os templates distribuídos continuem sendo compilados | `task test:templates` |
| Detecção de condições de corrida | Detectar condições de corrida de dados no runtime e na ponte | `go test -race ./...` |
| Matriz de CI | Confiança entre sistemas operacionais em cada PR | GitHub Actions |

Este documento explica **onde ficam os testes**, **como executá-los** e **o que o Taskfile orquestra**.

`v3/TESTING.md` documenta a compilação de exemplos multiplataforma; não é um catálogo de supressões de condições de corrida.

---

## 1. Convenções de diretórios

```
v3/
├── internal/.../_test.go     # Unit tests for internal packages
├── pkg/.../_test.go          # Public API tests
├── tasks/events/generate.go  # Code generator for event constants (NOT a test harness)
├── tests/                    # Top-level integration test harness
└── TESTING.md                # Cross-platform example-build guidance
```

Diretrizes:

- **Mantenha os testes unitários junto ao código** (`foo.go` ↔ `foo_test.go`).
- Use o **estilo caixa-preta** nos pacotes `pkg/` (`package application_test`) quando isso melhorar a organização da API.
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
- Quando necessário, isole as peculiaridades de plataforma por trás de tags de compilação (`foo_windows_test.go`, `foo_darwin_test.go`, …).

### Expectativa de cobertura

Espera-se que a lógica nova e alterada tenha 100% de cobertura de instruções Go. Meça o pacote alterado em vez de depender de uma porcentagem para todo o repositório:

```bash
go test ./pkg/application -coverprofile=coverage.out
go tool cover -func=coverage.out
```

Alguns caminhos não podem ser testados de forma razoável em um ambiente de testes normal — por exemplo, falhas exclusivas de uma plataforma, comportamentos dependentes de hardware ou um fallback defensivo que não pode ser provocado com segurança. Mantenha essas exceções restritas e explique cada caminho sem cobertura na descrição do PR.

### Execução local

```bash
cd v3
go test ./... -cover
```

Ou pelo Taskfile (targets reais — não há o atalho `task test`):

```
task test:cli            # CLI plumbing tests
task test:generator      # bindings generator round-trip tests
task test:templates      # build every shipped template
task test:infrastructure # supporting helpers
task test:examples       # exercise the example matrix (downloads as needed)
task test:all            # everything above
task sanity              # default GTK4 smoke check (legacy: sanity:gtk3)
task precommit           # what you should run before pushing
```

---

## 3. Testes de integração

`v3/tests/` contém a infraestrutura de integração entre pacotes. Os targets `test:example:*` e `test:examples:*` do Taskfile executam verificações de compilação em darwin / windows / linux (incluindo matrizes GTK3 / GTK4 baseadas em Docker no Linux).

> Os exemplos executáveis ficam em `v3/examples/`. Os targets de teste selecionam e compilam
>
> os exemplos adequados à plataforma do host ou à matriz de CI.

Execute a suíte de smoke tests da plataforma do host com:

```
task test:examples       # host
task test:examples:all   # full matrix (slow)
```

---

## 4. Detecção de condições de corrida

Condições de corrida de dados são fatais em runtimes de GUI.

### Como investigar uma condição de corrida

Inspecione os dois acessos conflitantes e reproduza o problema na plataforma afetada. Não considere um relatório inofensivo nem o suprima sem evidências. `v3/TESTING.md` descreve a compilação de exemplos; o sucesso dessas compilações não comprova a ausência de condições de corrida nem o comportamento do runtime nativo.

### Suíte local de detecção de condições de corrida

```
go test -race ./...
```

> `wails3 dev` não tem a flag `-race` — suas flags de CLI são `--config`, `--port`,
>
> e `-s` (habilita HTTPS). Para exercitar o runtime com o detector de condições de corrida,
>
> compile um aplicativo de teste com `go build -race` e execute-o diretamente.

---

## 5. Workflows do GitHub Actions

Arquivos de workflow reais em `.github/workflows/` (verificados na árvore):

| Arquivo | Finalidade |
| --- | --- |
| `build-and-test-v3.yml` | Matriz do Go 1.26/1.27, verificações do runtime e do gerador, compilação de templates e exemplos e verificações de compilação de aplicativos. Por padrão, o Linux usa GTK4/WebKitGTK 6.0; a variante legada define `BUILD_TAGS=gtk3`. A inicialização e a interação nativas devem ser verificadas separadamente. |
| `cross-compile-test-v3.yml` | Verificações básicas de compilação cruzada |
| `auto-changelog-v3.yml`, `changelog-v3.yml` | Automação do changelog |
| `nightly-release-v3.yml` | Artefatos da versão nightly v3 |
| `bump-webview2-v3.yml`, `release-webview2.yml` | Gerenciamento da dependência e das versões do WebView2 |
| `build-cross-image.yml` | Compila a imagem de contêiner do compilador cruzado |
| `publish-npm.yml` | Publica no npm o runtime JS `@wailsio/runtime` incorporado |
| `pr-master.yml` | Verificações do PR em relação à branch `master` |
| `semgrep.yml` | Análise estática com Semgrep |
| `stale-issues.yml`, `generate-sponsor-image.yml`, `build-and-test.yml`, `weekly-release-v2.yml` | Manutenção do repositório e workflows de v2; os antigos workflows de sincronização e upload do Crowdin foram descontinuados. |

As etapas de CI correspondem aos targets do Taskfile acima (`task test:cli`, `task test:generator`, `task test:templates`, `task test:examples`, …), portanto você pode reproduzir a CI localmente de forma idêntica. A etapa de smoke test `wails3 build` em `build-and-test-v3.yml` é invocada **sem flags adicionais** — não há a flag `-skip-package` em `wails3 build`.

---

## 6. Paridade local com a CI

Não há um único alvo abrangente `task ci`. Reproduza a CI encadeando os alvos reais:

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
| **Condição de corrida em `webview_window_darwin.go`** | Alteração do estado da janela fora da thread principal | Encaminhe por `application.InvokeAsync` / `Invoke` para que a chamada seja executada na thread principal |
| **O teste no Linux trava na CI sem interface gráfica** | O GTK precisa de um servidor de exibição | Execute com `xvfb-run`, por exemplo, `xvfb-run task test:examples:linux` |
| **A compilação do template falha** | O arquivo de lock do frontend está desatualizado | Execute `wails3 init` novamente em um diretório limpo para atualizar o template |
| **Erros de coverpkg** | Teste de integração que importa `main` | Mude para a tag de compilação `//go:build integration` e condicione a importação a ela |

---

## 8. Adição de novos testes

1. **Unitário** — crie `*_test.go` e execute `go test ./...`
2. **Gerador / CLI** — amplie os casos em `internal/generator/testcases/` ou `internal/commands/*_test.go` e execute novamente `task test:generator` / `task test:cli`
3. **Templates / exemplos** — verifique se os templates distribuídos ainda são compilados com `task test:templates`

---

## 9. Mapa dos principais arquivos

| O quê | Caminho |
| --- | --- |
| Teste de ida e volta do gerador | `internal/generator/generate_test.go` |
| Teste de compilação de assets | `internal/commands/build-assets_test.go` |
| Compilações de exemplos multiplataforma | `v3/TESTING.md` |
| Alvos de teste do Taskfile | `v3/Taskfile.yaml` |
| Gerador de constantes de eventos | `v3/tasks/events/generate.go` |
| Fluxo de trabalho de CI | `.github/workflows/build-and-test-v3.yml` (1.26/1.27 do Go via `actions/setup-go@v5`) |
| Análise estática | `.github/workflows/semgrep.yml` |
| Publicação do runtime no npm | `.github/workflows/publish-npm.yml` |

---

A qualidade não é algo secundário no Wails v3. Com testes unitários, suítes do gerador e de templates, detecção de condições de corrida e uma matriz de CI multiplataforma, você pode contribuir com confiança, sabendo que suas alterações passam em todos os sistemas operacionais compatíveis. Bons testes!
