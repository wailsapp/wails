---
title: "Sistema de bindings"
description: "Como o Wails v3 permite que Go e JavaScript chamem um ao outro sem nenhum código repetitivo"
slug: "contributing/binding-system"
sourcePath: "contributing/binding-system.md"
---

> “Bindings” são o **contrato com segurança de tipos** que permite escrever:

```go
msg, err := chatService.Send("Hello")
```

em Go *e*

```ts
import { Send } from "../bindings/github.com/you/yourapp/services/chatservice";

const msg = await Send("Hello");
```

em TypeScript **sem escrever manualmente nenhum código de integração de IPC**. Este documento detalha *como* isso acontece, desde a **análise estática** durante a compilação, passando pela **geração de código**, até a **ponte de runtime** que transfere bytes pelo WebView.

> Consulte [`contributing/architecture/bindings`](/contributing/architecture/bindings/) para ver
>
> a análise detalhada e oficial do pipeline do gerador — esta página é uma
>
> visão geral voltada a colaboradores.

---

## 1. Visão geral em 30 segundos

| Etapa | Componente | Saída |
| --- | --- | --- |
| **Coleta/análise** | `internal/generator/collect/`, `internal/generator/analyse.go` | Modelo em memória dos serviços Go exportados, métodos, parâmetros, tipos de retorno e modelos |
| **Geração** | `internal/generator/render/templates/*.tmpl` (`service.{js,ts}.tmpl`, `models.{js,ts}.tmpl`, `index.tmpl`, `eventcreate.js.tmpl`, `eventdata.d.ts.tmpl`, `newline.tmpl`) | Módulos ES por serviço em `frontend/bindings/<full Go import path>/...` |
| **Runtime** | `pkg/application/messageprocessor*.go` + o runtime JS incorporado em `internal/runtime/desktop/@wailsio/runtime/src/` (`calls.ts`, `events.ts`, …) | Mensagens de chamada/evento pela ponte nativa do WebView |

O fluxo é orquestrado pelo comando `wails3 generate bindings`, que executa `generator.Generate` (definido em `internal/generator/generate.go`) em um conjunto de pacotes Go.

```
wails3 generate bindings
        │
        ▼
internal/generator/generate.go: Generator.Generate(patterns...)
        │
        ├── internal/generator/collect/   // load.go, collector.go, service.go, model.go, …
        ├── internal/generator/analyse.go // semantic checks
        └── internal/generator/render/    // template execution → frontend/bindings/**
```

---

## 2. Análise estática

### Ponto de entrada

```
internal/generator/generate.go        // Generator + Generate(patterns…)
internal/generator/analyse.go         // semantic validation
internal/generator/collect/load.go    // go/packages loader
internal/generator/collect/collector.go
```

A etapa de coleta percorre cada pacote carregado e registra:

- `collect.ServiceInfo` — um para cada struct Go exportado e vinculado.
- `collect.ServiceMethodInfo` / `collect.MethodInfo` — informações da assinatura de cada método (nome, parâmetros, resultados, posição do erro, receptor e documentação).
- `collect.ModelInfo` / `collect.StructInfo` — emitidos como modelos TS/JS.
- Comentários de diretiva, como `//wails:inject`, `//wails:include`, `//wails:internal`, `//wails:ignore` e `//wails:id <hex>` (consulte `internal/generator/collect/directive.go`).

Tipos não compatíveis geram um erro no gerador para que os enganos apareçam durante a compilação, e não em runtime.

### Identificadores de modelo

O envelope de chamada do runtime identifica um método por um **hash FNV-1a determinístico** de seu nome totalmente qualificado (`pkg.Struct.Method`). Ele aparece como `$Call.ByID(<numeric-id>, …)` nos bindings gerados ou como `$Call.ByName("pkg.Struct.Method", …)` quando a geração é executada com `-names`.

---

## 3. Geração de código

### Templates

`internal/generator/render/templates/`:

| Template | Finalidade |
| --- | --- |
| `service.js.tmpl` | Um módulo JS por serviço vinculado |
| `service.ts.tmpl` | Arquivo complementar TypeScript (com `-ts`) |
| `models.js.tmpl` | Saída de classes de modelo (por pacote) |
| `models.ts.tmpl` | Saída de `.d.ts` do modelo (por pacote) |
| `index.tmpl` | Reexportações agregadas de `index.{js,ts}` por pacote |
| `eventcreate.js.tmpl` / `eventdata.d.ts.tmpl` | Construtor de eventos / tipagens de payload |
| `newline.tmpl` | Normalizador de quebra de linha final |

A saída é gravada em `frontend/bindings/<full Go import path>/...` — por exemplo, um serviço definido em `github.com/you/yourapp/services/chat` é gravado em `frontend/bindings/github.com/you/yourapp/services/chat/`. Não há um diretório `frontend/src/wailsjs/` na v3.

### Saída JavaScript

Os bindings gerados são módulos ES que importam os auxiliares de runtime de `/wails/runtime.js`:

```js
import { Call as $Call, Create as $Create } from "/wails/runtime.js";

/**
 * @param {string} msg
 * @returns {Promise<string> & { cancel(): void }}
 */
export function Send(msg) {
    return $Call.ByID(2042131923, msg);
}
```

Quando a geração é executada com `-names`, `$Call.ByName("pkg.Struct.Method", ...)` é emitido em seu lugar — sempre **totalmente qualificado**, nunca apenas `"Method"`.

As classes de modelo geradas usam um padrão de construtor `$$source`, com valores padrão `if (!("X" in $$source))` para cada campo, nomes de campos entre aspas e um `static createFrom(...)` que executa `JSON.parse` em entradas de string.

### Principais mapeamentos de tipos

Verificado em relação a `internal/generator/render/`:

| Go | TypeScript |
| --- | --- |
| `map[string]V` | `{ [_: string]: V }` |
| `map[K]V` (`K` que não seja string) | `{ [_ in K]?: V }` (não `Map<K, V>`, nem `Record<K, V>`) |
| `[]byte` | `Uint8Array` |
| `time.Time` | `string` (JSON no formato ISO 8601) |
| `error` (na posição de retorno) | promessa rejeitada |

### Observação sobre reflexão

`pkg/application/bindings.go` é **escrito manualmente** e usa `reflect` para controlar o despacho de métodos a partir de um registro `BoundMethod`. Não interprete literalmente as afirmações antigas de "zero reflexão em tempo de execução": o gerador evita reflexão, mas o despachante em tempo de execução a utiliza.

---

## 4. Protocolo de invocação em tempo de execução

### Lado do JavaScript

```ts
import { Call } from "/wails/runtime.js";

await Call.ByID(0x7a1201d3 /* ChatService.Send */, "Hello");
// or, with -names:
await Call.ByName("chatservice.ChatService.Send", "Hello");
```

Os auxiliares de tempo de execução ficam em `internal/runtime/desktop/@wailsio/runtime/src/calls.ts` (despacho de chamadas), `events.ts` (eventos) e arquivos relacionados — não há nenhum `invoke.ts` nem `errors.ts` nesta árvore. O envelope exato do protocolo é codificado por `calls.ts` no lado do JS e decodificado por `pkg/application/messageprocessor_call.go` no lado do Go; consulte esses dois arquivos em conjunto ao depurar a ponte.

### Lado do Go

1. `pkg/application/messageprocessor_call.go` recebe a mensagem de chamada.
2. Procura o método vinculado por ID ou nome em `pkg/application/bindings.go` (orientado por `reflect`).
3. Invoca o método vinculado e serializa `{result, error}` de volta para o JS.

### Mapeamento de erros

| Go | JavaScript |
| --- | --- |
| `error == nil` | `Promise` é resolvida com o resultado |
| `error != nil` | `Promise` é rejeitada com um `Error` cujo `message` contém a string do erro do Go |

---

## 5. Como chamar JavaScript a partir do Go

O gerador de bindings funciona em uma única direção (métodos Go expostos ao JS). Para a comunicação Go → JS, use o barramento de eventos ou execute JS em uma janela:

```go
app.Event.Emit("chat:new-message", msg)
window.ExecJS(`window.dispatchEvent(new CustomEvent("ping"))`)
```

No lado do JS, faça a assinatura com `Events.On(name, cb)` de `/wails/runtime.js`.

---

## 6. Extensão e solução de problemas

### Erro de tipo não compatível

```
error: field "Client" uses unsupported type: chan struct{}
```

→ encapsule o canal por trás de uma API de método ou marque o campo com `//wails:internal` para que o gerador o ignore.

### Bindings desatualizados

A saída gerada é sobrescrita a cada `wails3 generate bindings` / `wails3 dev` / `wails3 build`. Se o IntelliSense da IDE exibir stubs desatualizados, exclua `frontend/bindings/` e execute novamente o gerador. A opção `-clean` (cujo padrão é `true` nas compilações atuais) limpa o diretório de bindings antes de cada execução.

### Dicas de desempenho

- Evite transmitir grandes slices de bytes pela ponte — em vez disso, disponibilize-os pelo servidor de assets.
- Quando a latência for importante, agrupe várias chamadas rápidas em um único método.
- Prefira receptores por valor para structs de parâmetros pequenas a fim de reduzir as alocações.

---

## 7. Mapa dos principais arquivos

| Área | Arquivo |
| --- | --- |
| Orquestração do gerador | `internal/generator/generate.go` |
| Verificações semânticas | `internal/generator/analyse.go` |
| Coleta (serviços, métodos e modelos) | `internal/generator/collect/{service,method,model,struct,package}.go` |
| Templates de renderização | `internal/generator/render/templates/*.tmpl` |
| Local dos bindings gerados | `frontend/bindings/<full Go import path>/...` |
| Despachante do lado do Go | `pkg/application/bindings.go`, `messageprocessor_call.go` |
| Runtime do JS | `internal/runtime/desktop/@wailsio/runtime/src/{calls,events,index}.ts` |

Mantenha esta referência rápida à mão ao rastrear um bug na ponte.

---

## 8. Recapitulação

1. O **coletor** examina seu código Go → modelo semântico em memória.
2. Os **templates** geram módulos ES por serviço e arquivos de modelo/índice por pacote.
3. O **processador de mensagens** despacha chamadas no lado do Go por meio do registro de bindings.
4. O **runtime do JS** encapsula tudo em promises idiomáticas com cancelamento.

Tudo isso sem que você escreva uma única linha de código repetitivo de IPC. Esse é o sistema de bindings do Wails v3. Agora, mãos à obra com os bindings!
