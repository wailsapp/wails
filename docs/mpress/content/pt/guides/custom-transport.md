---
title: "Criação de uma camada de transporte personalizada"
description: "Saiba como criar e personalizar sua própria camada de transporte IPC no Wails v3"
slug: "guides/custom-transport"
sourcePath: "guides/custom-transport.md"
---

O Wails v3 permite fornecer uma camada de transporte IPC personalizada, preservando todos os bindings gerados e a comunicação de eventos. Isso permite substituir o transporte padrão baseado em requisições HTTP fetch por WebSockets, protocolos personalizados ou qualquer outro mecanismo de transporte.

## Visão geral

Por padrão, o Wails usa requisições HTTP fetch do frontend para se comunicar com o backend por meio de `/wails/runtime`. A API de transporte personalizado permite:

- Substituir o transporte HTTP por WebSockets, gRPC ou qualquer protocolo personalizado
- Manter compatibilidade total com a geração de código do Wails
- Preservar todos os bindings, eventos, caixas de diálogo e outros recursos existentes do Wails
- Implementar seu próprio gerenciamento de conexões, autenticação e tratamento de erros

## Arquitetura

```text
┌─────────────────────────────────────────────────┐
│  Frontend (TypeScript)                          │
│  - Generated bindings still work                │
│  - Your custom client transport                 │
└──────────────────┬──────────────────────────────┘
                   │
                   │ Your Protocol (WebSocket/etc)
                   │
┌──────────────────▼──────────────────────────────┐
│  Backend (Go)                                   │
│  - Your Transport implementation                │
│  - Wails MessageProcessor                       │
│  - All existing Wails infrastructure            │
└─────────────────────────────────────────────────┘
```

## Uso

### 1. Implemente a interface de transporte

Crie um transporte personalizado implementando a interface `Transport`:

```go
package main

import (
    "context"
    "github.com/wailsapp/wails/v3/pkg/application"
)

type MyCustomTransport struct {
    // Your fields
}

func (t *MyCustomTransport) Start(ctx context.Context, processor *application.MessageProcessor) error {
    // Initialize your transport (WebSocket server, gRPC server, etc.)
    // When you receive requests, call processor.HandleRuntimeCallWithIDs()
    return nil
}

func (t *MyCustomTransport) Stop() error {
    // Clean up your transport
    return nil
}
```

### 2. Configure sua aplicação

Passe seu transporte personalizado para as opções da aplicação:

```go
func main() {
    app := application.New(application.Options{
        Name: "My App",
        Transport: &MyCustomTransport{},
        // ... other options
    })

    err := app.Run()
    if err != nil {
        log.Fatal(err)
    }
}
```

### 3. Modifique o runtime do frontend

Se você usar um transporte personalizado, precisará modificar o runtime do frontend para usar esse transporte em vez de HTTP fetch. Implemente a interface `RuntimeTransport`, que será usada para processar as requisições:

```typescript
const { setTransport } = await import('/wails/runtime.js');

class MyRuntimeTransport {
  call(objectID: number, method: number, windowName: string, args: any): Promise<any> {
    // TODO: implement IPC call with your transport protocol

    return resp;
  }
}

const myTransport = new MyRuntimeTransport();
setTransport(myTransport);
```

## Observações

- O transporte HTTP padrão continuará funcionando se nenhum transporte personalizado for especificado
- Os bindings gerados permanecem inalterados — somente a camada de transporte muda
- Eventos, caixas de diálogo, área de transferência e todos os outros recursos do Wails funcionam de forma transparente
- Você é responsável pelo tratamento de erros, pela lógica de reconexão e pela segurança do seu transporte personalizado
- O exemplo de WebSocket fornecido é apenas para demonstração e pode precisar de reforços de segurança e robustez para uso em produção

## Referência da API

### Interface de transporte

```go
type Transport interface {
    Start(ctx context.Context, messageProcessor *application.MessageProcessor) error
    // JSClient returns the JavaScript shim that the runtime injects into the
    // window so that frontend code can call into the transport.
    JSClient() []byte
    Stop() error
}
```

### Interface AssetServerTransport (opcional)

Para implantações baseadas em navegador ou quando você quiser fornecer tanto os ativos quanto o IPC por meio do seu transporte personalizado, implemente a interface `AssetServerTransport`:

```go
type AssetServerTransport interface {
    Transport

    // ServeAssets configures the transport to serve assets alongside IPC.
    // The assetHandler is Wails' internal asset server that handles all assets,
    // runtime.js, capabilities, flags, etc.
    ServeAssets(assetHandler http.Handler) error
}
```

**Quando implementar esta interface:**

- Executar a aplicação em um navegador em vez de uma webview
- Fornecer ativos por HTTP juntamente com seu transporte IPC personalizado
- Criar aplicações acessíveis pela rede

**Exemplo de implementação:**

```go
func (t *MyTransport) ServeAssets(assetHandler http.Handler) error {
    mux := http.NewServeMux()

    // Mount your IPC endpoint
    mux.HandleFunc("/my/ipc/endpoint", t.handleIPC)

    // Mount Wails asset server for everything else
    mux.Handle("/", assetHandler)

    // Start HTTP server
    t.httpServer.Handler = mux
    go t.httpServer.ListenAndServe()

    return nil
}
```

Quando `ServeAssets()` é chamado, o assetHandler fornece:

- Todos os ativos estáticos (HTML, CSS, JS, imagens etc.)
- `/wails/runtime.js` — A biblioteca de runtime do Wails

## Veja também

- `transport.go` — Interfaces e tipos principais de transporte
- `messageprocessor.go` — O processador de mensagens subjacente que processa todo o IPC do Wails
