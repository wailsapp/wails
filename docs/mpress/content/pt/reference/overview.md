---
title: "Referência da API"
description: "Consulte tipos, métodos, opções e comandos da CLI do Wails"
slug: "reference/overview"
sourcePath: "reference/overview.md"
---

## Sobre esta referência

Use esta seção para consultar tipos, métodos, opções e comandos do Wails enquanto trabalha em um aplicativo. Se precisar de uma lição guiada, comece pelos [Tutoriais](/tutorials/overview/). Para ver as etapas de execução de uma tarefa, use os [Guias práticos](/quick-start/next-steps/#features). Para conhecer os fundamentos de design, consulte os [Conceitos fundamentais](/concepts/architecture/).

**Organização:**

- [Aplicativo](/reference/application/) — APIs principais do aplicativo
- [Janela](/reference/window/) — Criação e gerenciamento de janelas; [opções de janela](/features/windows/options/) lista os campos de configuração
- [Menu](/reference/menu/) — APIs de aplicativo, menu de contexto e bandeja do sistema; a [referência de itens de menu](/features/menus/reference/) aborda os tipos e as propriedades dos itens
- [Eventos](/reference/events/) — Sistema de eventos e eventos integrados
- [Caixas de diálogo](/reference/dialogs/) — Caixas de diálogo de arquivos e mensagens
- [Runtime do frontend](/reference/frontend-runtime/) — APIs do runtime do frontend
- [Visão geral dos comandos da CLI](/reference/cli/) — Comandos agrupados por finalidade; [opções e exemplos](/guides/cli/) apresenta detalhes de cada comando
- [Protocolo de manifesto de atualização](/reference/update-manifest/) — Metadados de atualização e campos de verificação
- [API para dispositivos móveis](/guides/mobile/mobile-api/) — APIs experimentais para dispositivos móveis

## Convenções da API

@details{title="Convenções da API Go — Para desenvolvedores iniciantes em Go"}
### Nomenclatura

- <strong>Tipos</strong>: PascalCase (por exemplo, `WebviewWindow`)
- <strong>Métodos</strong>: PascalCase (por exemplo, `SetTitle()`)
- <strong>Opções</strong>: structs em PascalCase (por exemplo, `WindowOptions`)
- <strong>Constantes</strong>: PascalCase (por exemplo, `WindowStartStateMaximised`)

#### Tratamento de erros

A maioria dos métodos que podem falhar retorna `error` como último valor de retorno. `app.Run()` bloqueia até o encerramento do aplicativo e retorna qualquer erro de inicialização:

```go
if err := app.Run(); err != nil {
    log.Fatal(err)
}
```

A criação de janelas não retorna um erro — `app.Window.New()` retorna `*WebviewWindow` diretamente.

#### Contexto

Os métodos do ciclo de vida do serviço recebem um `context.Context`:

```go
func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    // ctx is cancelled when the application is shutting down.
    return nil
}
```

O contexto do ciclo de vida do aplicativo está disponível por meio de `app.Context()`. Não há `RunWithContext` — chame `app.Run()`.

#### Padrão de opções

A configuração usa structs de opções:

```go
app := application.New(application.Options{
    Name: "My App",
    Description: "A demo application",
    Services: []application.Service{
        application.NewService(&MyService{}),
    },
})
```

@end

### Convenções da API JavaScript

#### Nomenclatura

- **Funções**: camelCase (por exemplo, `setTitle()`)
- **Constantes**: SCREAMING<em>SNAKE</em>CASE (por exemplo, `WINDOW_EVENT_FOCUS`)

#### Assíncrona por padrão

Todas as chamadas de métodos Go retornam Promises:

```javascript
// Async/await (recommended)
const result = await MyService.DoSomething()

// Promise chain
MyService.DoSomething()
    .then(result => console.log(result))
    .catch(error => console.error(error))
```

#### Tratamento de erros

Erros do Go tornam-se exceções do JavaScript:

```javascript
try {
    await MyService.MightFail()
} catch (error) {
    console.error('Go error:', error)
}
```

#### Segurança de tipos

As definições do TypeScript são geradas automaticamente:

```typescript
// Fully typed
import { Greet } from './bindings/GreetService'

const message: string = await Greet("World")
```

## Estrutura de pacotes

```
github.com/wailsapp/wails/v3/pkg/
├── application/          # Core application package
│   ├── application.go    # App type
│   ├── webview_window.go # Window management
│   ├── menu.go           # Menu types
│   ├── event_manager.go  # Event system
│   └── dialogs.go        # Dialog APIs
├── events/               # Event constants
└── services/             # Built-in services
    ├── dock/             # macOS dock (includes badge support)
    ├── fileserver/       # File-server service
    ├── kvstore/          # Key/value store
    ├── log/              # Structured logging service
    ├── notifications/    # Notifications service
    └── sqlite/           # SQLite service
```

## Caminhos de importação

### Go

```go
import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)
```

### JavaScript

```javascript
// Auto-generated bindings
import { MyMethod } from './bindings/MyService'

// Runtime APIs
import { Events, Window } from '@wailsio/runtime'
```

## Referência de tipos

### Tipos comuns

@tabs{sync-key="lang"}
[Go]
```go
// Application
type App struct { /* ... */ }
type Options struct { /* ... */ }

// Window
type WebviewWindow struct { /* ... */ } // implements the Window interface
type WebviewWindowOptions struct { /* ... */ }

// Menu
type Menu struct { /* ... */ }
type MenuItem struct { /* ... */ }

// Events — there is no generic Event type; events are typed by source.
type ApplicationEvent struct { /* ... */ }
type WindowEvent struct { /* ... */ }
type CustomEvent struct { /* ... */ }
type EventListener struct { /* ... */ }

// Dialogs
type OpenFileDialogOptions struct { /* ... */ }
type SaveFileDialogOptions struct { /* ... */ }
```

[TypeScript]
```typescript
// Window runtime
interface WindowOptions {
    title?: string
    width?: number
    height?: number
    // ...
}

// Events
type EventCallback = (data: any) => void

// Bindings (auto-generated)
export function MyMethod(arg: string): Promise<string>
```

@end

## Diferenças entre plataformas

Algumas APIs se comportam de maneira diferente dependendo da plataforma:

| Recurso | Windows | macOS | Linux |
| --- | --- | --- | --- |
| **Menu do aplicativo** | Barra de menus da janela | Barra de menus global | Barra de menus da janela |
| **Bandeja do sistema** | Área de notificação | Barra de menus | Bandeja do sistema |
| **Dock** | N/A | ✅ Disponível | N/A |
| **Caixas de diálogo de arquivos** | Nativa | Nativa | Nativa (GTK) |
| **Transparência** | ✅ Total | Requer [`private_mac_apis`](/guides/build/private-macos-apis/#webview-transparency-and-background) | ⚠️ Limitada |

O comportamento específico de cada plataforma está documentado em cada seção da API.

## Versionamento

O Wails v3 segue o versionamento semântico:

- **Principal** (v4.0.0): alterações incompatíveis
- **Secundária** (v3.1.0): novos recursos, com compatibilidade retroativa
- **Correção** (v3.0.1): correções de bugs, com compatibilidade retroativa

**Wails v3.0.0:** estável

## Política de descontinuação

Quando as APIs são descontinuadas:

1. **Marcadas na documentação** com um aviso de descontinuação
2. **Alternativa fornecida** com um guia de migração
3. **Mantidas por 1 versão principal** antes da remoção
4. **Avisos do compilador** (quando possível)

## Estabilidade da API

### APIs estáveis ✅

Estas APIs são estáveis e seguras para uso em produção:

- APIs principais do aplicativo
- Gerenciamento de janelas
- Sistema de menus
- Sistema de eventos
- Caixas de diálogo de arquivos
- Bindings de serviços

### APIs experimentais ⚠️

Somente as APIs explicitamente documentadas como experimentais estão fora do contrato de compatibilidade estável para desktop:

- Suporte ao Android
- Suporte ao iOS
- Outras APIs explicitamente marcadas como experimentais

As APIs experimentais são identificadas na documentação.

## Como obter ajuda

### Dúvidas sobre a API

1. **Consulte esta referência** — Documentação completa da API
2. **Consulte os exemplos** — [Exemplos no GitHub](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples)
3. **Pesquise no Discord** — [Servidor do Discord](https://discord.gg/JDdSxwjhGf)
4. **Pergunte à comunidade** — Canal #help no Discord

### Como relatar problemas na API

Encontrou um bug ou uma inconsistência?

1. **Consulte os problemas existentes** — [Problemas no GitHub](https://github.com/wailsapp/wails/issues)
2. **Crie um relatório detalhado** — Inclua o código, o erro e a plataforma
3. **Forneça uma reprodução** — Um exemplo mínimo que demonstre o problema

## Documentação relacionada

- [Tutoriais](/tutorials/overview/) — Aprenda criando aplicações reais
- [Guias](/guides/architecture/) — Guias orientados a tarefas para cenários comuns
- [Recursos](/features/windows/basics/) — Documentação de cada recurso
- [Exemplos](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples) — Exemplos de código funcionais no GitHub

---

**Explore a API:** use a navegação à esquerda para conhecer APIs específicas.
