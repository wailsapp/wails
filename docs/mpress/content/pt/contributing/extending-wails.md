---
title: "Como estender o Wails"
description: "Guia prático para adicionar novos recursos e plataformas ao Wails v3"
slug: "contributing/extending-wails"
sourcePath: "contributing/extending-wails.md"
---

> O Wails foi projetado para ser **modificável**.
>
> Todos os principais subsistemas estão em código Go que você pode ler, modificar e distribuir.
>
> Esta página mostra *por onde* começar e *como* manter a compatibilidade multiplataforma ao:

- Adicionar um **serviço** (notificações, armazenamento de chave-valor, IPC personalizado, …)
- Criar um **novo comando da CLI** (`wails3 <foo>`)
- Estender o **runtime** (API de janelas, caixas de diálogo, eventos)
- Introduzir um **recurso de plataforma** (Wayland, …)
- Manter a **compatibilidade multiplataforma** sem se perder em tags `//go:build`

---

## 1. Como adicionar um serviço

Um "serviço" na v3 é um tipo Go fornecido pelo usuário, registrado por meio de `application.Options.Services` e exposto ao JS por bindings gerados. A base de código da v3 inclui:

- `internal/service/` — estrutura básica para `wails3 generate service`:
  ```
  internal/service/
  ├── service.go              # Install(options *flags.ServiceInit)
  └── template/
      ├── README.tmpl.md
      ├── go.mod.tmpl
      ├── service.go.tmpl
      └── service.tmpl.yml
  ```

- `pkg/services/` — serviços prontos que você pode registrar hoje (notificações, kvstore, sqlite, log, fileserver, dock, …).

Os arquivos do gerador e da CLI mencionados em versões anteriores como `internal/service/template/template.go` e `internal/generator/collect/services.go` não existem — a ferramenta de geração da estrutura básica é `internal/service/service.go` (ponto de entrada `service.Install`), e os metadados de binding dos serviços são coletados em `internal/generator/collect/service.go`.

### 1.1 Defina o serviço

```go
package chat

type Service struct {
    messages []string
}

func New() *Service { return &Service{} }

func (s *Service) Send(msg string) string {
    s.messages = append(s.messages, msg)
    return "ok"
}
```

### 1.2 Implemente as interfaces de ciclo de vida (opcional)

Opcionalmente, um serviço pode implementar estas interfaces (de `pkg/application`):

```go
func (s *Service) ServiceStartup(ctx context.Context, options application.ServiceOptions) error { return nil }
func (s *Service) ServiceShutdown() error                                                       { return nil }
```

> **Importante:** `ServiceShutdown` **não recebe argumentos**. Um método com a
>
> assinatura `ServiceShutdown(ctx context.Context) error` **não** implementa
>
> a interface e, sem qualquer aviso, nunca será chamado.

### 1.3 Registre o serviço no aplicativo

Não há uma chamada global `services.Register(...)`. Os serviços são registrados em tempo de execução por meio de `application.Options.Services`:

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(chat.New()),
    },
})
```

Após o registro, `wails3 generate bindings` gera módulos ES em `frontend/bindings/<your import path>/...` que encapsulam os métodos exportados.

### 1.4 Chame a partir do JS

```js
import { Send } from "../bindings/github.com/you/yourapp/chat";

await Send("hi");
```

Não há um `window.backend.*` global na v3 — as chamadas passam pelos módulos ES gerados, que, por sua vez, chamam `Call.ByID(...)` de `/wails/runtime.js`.

---

## 2. Como escrever um novo comando da CLI

A CLI da v3 usa **`github.com/leaanthony/clir`** (não cobra). A integração fica em `v3/cmd/wails3/main.go`:

```go
import "github.com/leaanthony/clir"

func main() {
    app := clir.NewCli("wails", "The Wails3 CLI", "v3")
    app.NewSubCommand("hello", "Prints Hello Wails").Action(func() error {
        fmt.Println("Hello Wails!")
        return nil
    })
    // ... other subcommands explicitly wired here
    _ = app.Run()
}
```

Não há registro automático baseado em `init()`. Adicione o novo subcomando a `cmd/wails3/main.go`, além de uma função de implementação em `internal/commands/` (e uma struct de flags em `internal/flags/`, caso aceite opções). Recompile a CLI:

```
cd v3
go install ./cmd/wails3
wails3 hello
```

Se o comando precisar da integração com o Taskfile, reutilize os auxiliares em `internal/commands/task_wrapper.go` (`wrapTask("yourtask", args)`).

---

## 3. Como modificar o runtime

Motivos comuns:

- Novo recurso de janela (`SetOpacity`, `Shake`, …)
- Caixa de diálogo adicional (`ColorPicker`)
- API de nível de sistema (brilho da tela)

### 3.1 API pública

Adicione o método a `pkg/application/webview_window.go` (a interface fica em `window.go`):

```go
func (w *WebviewWindow) SetOpacity(o float32) Window {
    InvokeSync(func() { w.impl.setOpacity(o) })
    return w
}
```

Use os auxiliares `InvokeSync`/`InvokeAsync` existentes para garantir que a chamada seja executada na thread principal.

### 3.2 Processador de mensagens

Se o JS precisar invocar o novo método, estenda o arquivo `pkg/application/messageprocessor_*.go` apropriado. O processador de mensagens usa métodos baseados em switch em `MessageProcessor`, e não uma chamada global `register(...)`:

```go
// inside messageprocessor_window.go
case "setOpacity":
    var args struct {
        WindowID uint    `json:"windowID"`
        Opacity  float32 `json:"opacity"`
    }
    if err := json.Unmarshal(payload, &args); err != nil { ... }
    window, _ := m.app.Window.GetByID(args.WindowID)
    window.SetOpacity(args.Opacity)
```

**Não** há nenhum arquivo `messageprocessor_window_opacity.go` nem padrão `register(MsgSetOpacity, ...)` baseado em `init()` na base de código.

### 3.3 Implementação da plataforma

Adicione a implementação a cada arquivo específico de sistema operacional em `pkg/application/`:

```
pkg/application/
├── webview_window_darwin.go   //go:build darwin
├── webview_window_linux.go    //go:build linux
└── webview_window_windows.go  //go:build windows
```

Se uma plataforma não puder oferecer suporte ao recurso, escreva um stub que não execute nenhuma operação. O framework não tem um sentinela `ErrCapability` — indique o suporte na documentação e, se necessário, no campo booleano pertinente de `Options` ou na struct de opções específica da plataforma.

### 3.4 Flag de recurso (opcional)

O pacote `internal/capabilities/` existe para declarar conjuntos de recursos específicos de cada plataforma. Não há uma API pública `application.HasCapability` / `application.CapOpacity`. Se quiser um recurso que possa ser verificado em tempo de execução, adicione-o em `internal/capabilities/` e exponha um getter tipado de `pkg/application`.

---

## 4. Como adicionar novos recursos de plataforma

Exemplo: suporte opcional ao Wayland no Linux.

1. Divida o arquivo `pkg/application/*_linux.go` pertinente em `*_linux_x11.go` (`//go:build linux && !wayland`) e `*_linux_wayland.go` (`//go:build linux && wayland`).
2. Permita que os usuários ativem o recurso com `wails3 build --tags wayland`. Encaminhe as tags adicionais pela infraestrutura `EXTRA_TAGS` existente em `internal/commands/task_wrapper.go`. Não há uma flag `--tags wayland` no nível de `dev` — `wails3 dev` aceita apenas `--config`, `--port` e `-s`.
3. Atualize a documentação e qualquer README específico da plataforma em `pkg/application/`.

> Mantenha as tags de compilação padrão no mínimo necessário; reserve as tags de ativação opcional para recursos de nicho.

---

## 5. Lista de verificação de compatibilidade multiplataforma

| ✅ Etapa | Motivo |
| --- | --- |
| Disponibilize **todos** os métodos públicos em todos os arquivos de plataforma (mesmo que sejam stubs) | Mantém a compilação funcionando em todos os sistemas operacionais |
| Documente a degradação gradual em cada sistema operacional | Os aplicativos podem decidir o fluxo com base em `runtime.GOOS` sem erros ocultos |
| Use primeiro **Go puro**; use Cgo somente quando necessário | Simplifica a compilação cruzada (o Linux já arca com o custo do Cgo) |
| Execute `task test:cli`, `task test:generator` e `task test:templates` | Reproduz a CI localmente |
| Documente as novas tags de compilação na documentação para colaboradores e no README do modelo | Os usuários precisam conhecer as opções que exigem ativação explícita |

---

## 6. Compilações de depuração e velocidade de iteração

- Use `Options.LogLevel = slog.LevelDebug` (`Options.Logger = slog.Default()`) para exibir informações detalhadas sobre a atividade do runtime. Não existe uma variável de ambiente `WAILS_LOG_LEVEL`.
- As opções de `wails3 dev` são `--config`, `--port` e `-s`. Não há uma opção `-race` nem `-verbose` — execute o detector de condições de corrida com `go test -race ./...` ou compilando seu aplicativo com `go build -race` e executando-o diretamente.
- O guia de testes de condições de corrida e Cgo está em `v3/TESTING.md` (rascunhos mais antigos apontavam para `pkg/application/RACE.md`, que não existe).

---

## 7. Contribuições ao projeto upstream

1. Para uma nova funcionalidade ou uma alteração de comportamento público, abra um PR de rascunho com uma **WEP (Wails Enhancement Proposal)** para discutir a ideia e o design. Use uma issue somente para um bug reproduzível ou um problema na documentação.
2. Siga as técnicas acima para fazer a implementação.
3. Adicione:
  - Testes unitários (`*_test.go`)
  - Documentação (este arquivo ou a página relevante de `docs/...`)
  - Um teste de regressão em `internal/generator/testcases/`, caso tenha alterado o gerador de bindings

4. Antes de enviar as alterações, execute localmente `task precommit` e os alvos relevantes de `task test:*`.

---

### Links rápidos

| Área | Localização |
| --- | --- |
| Serviços integrados | `pkg/services/` |
| Gerador de estrutura de serviços | `internal/service/` |
| Integração da CLI | `v3/cmd/wails3/main.go` |
| Implementações dos comandos da CLI | `internal/commands/` |
| Runtime específico de cada sistema operacional | `pkg/application/*_{darwin,linux,windows}.go` |
| Declarações de recursos | `internal/capabilities/` |
| DSL do Taskfile | `v3/Taskfile.yaml` |
| Gerador de constantes de eventos | `v3/tasks/events/generate.go` |

---

Agora você tem um **roteiro** para adaptar o Wails às suas necessidades — adicione serviços, dê um toque de magia à CLI, modifique o runtime ou implemente recursos totalmente novos para sistemas operacionais. Boas extensões!
