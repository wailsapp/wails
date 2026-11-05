---
title: "Criação de uma camada de transporte personalizada"
description: "Saiba como criar e personalizar sua própria camada de transporte IPC personalizada do Wails v3"
slug: "guides/custom-transport"
sourcePath: "guides/custom-transport.md"
---

O Wails v3 permite fornecer uma camada de transporte IPC personalizada, preservando todos os bindings gerados e a comunicação de eventos. Assim, você pode substituir o transporte padrão baseado em solicitações HTTP fetch por WebSockets, protocolos personalizados ou qualquer outro mecanismo de transporte.

## Visão geral

Por padrão, o Wails usa solicitações HTTP fetch do frontend para se comunicar com o backend por meio de `/wails/runtime`. A API de transporte personalizado permite:

- Substituir o transporte HTTP por WebSockets, gRPC ou qualquer protocolo personalizado
- Manter compatibilidade total com a geração de código do Wails
- Manter todos os bindings, eventos, diálogos e outros recursos existentes do Wails
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

## Executar o exemplo completo de WebSocket

O [exemplo de transporte WebSocket](https://github.com/wailsapp/wails/tree/master/v3/examples/websocket-transport) inclui as duas partes do protocolo: IDs de solicitação e respostas, gerenciamento do ciclo de vida, cliente JavaScript e entrega de eventos. Comece por esses arquivos, em vez de usar um transporte que implemente apenas `Start` e `Stop`.

Em uma cópia de trabalho do repositório do Wails com as dependências nativas da plataforma instaladas:

```bash
cd v3/examples/websocket-transport
go run .
```

O exemplo abre uma janela nativa e inicia um endpoint WebSocket na porta 9099. Use-o em um ambiente local de desenvolvimento confiável. Antes de implantar uma variante acessível pela rede, restrinja o listener, valide as origens e autentique os clientes. A política de origem do exemplo é permissiva e não constitui uma política de segurança para produção.

### Contrato do backend

Um transporte deve implementar os três métodos obrigatórios. O `WebSocketTransport` do exemplo implementa `Start`, `JSClient` e `Stop`; `JSClient` fornece o adaptador JavaScript que o Wails injeta na página. Inclua em sua implementação uma verificação de interface em tempo de compilação:

```go
var _ application.Transport = (*WebSocketTransport)(nil)
```

Configure a aplicação usando o construtor do exemplo:

```go
Transport: NewWebSocketTransport("127.0.0.1:9099"),
```

Esse campo pertence a `application.Options`, junto com os assets e serviços comuns. Copie os arquivos de implementação, além da chamada ao construtor. `Start` recebe o `MessageProcessor`; `Stop` deve fechar listeners, conexões ativas e trabalhos pendentes.

### Contrato do frontend

O `websocket-transport.js` do exemplo implementa o transporte do runtime, incluindo a correlação entre solicitações e respostas e a rejeição de chamadas que falharam. Depois de disponibilizar esse arquivo por meio de seus assets, inicialize-o antes de invocar os bindings de serviço gerados:

```javascript
import { setTransport } from '/wails/runtime.js'
import { createWebSocketTransport } from '/websocket-transport.js'

setTransport(createWebSocketTransport('ws://127.0.0.1:9099/wails/ws'))
```

Os bindings de serviço gerados permanecem inalterados. Se o seu transporte substituir o cliente, ele deverá oferecer os métodos usados pela aplicação e preservar a semântica de erros deles. A entrega de eventos exige o contrato opcional de transporte de eventos; ele não é fornecido automaticamente por uma implementação arbitrária de `Transport`. O exemplo completo de WebSocket demonstra esse contrato.

Para implantações remotas, use TLS e políticas explícitas de autenticação e de origem. A hospedagem exclusiva em navegador também exige a disponibilização dos assets da aplicação, conforme descrito abaixo.

## Referência da API

### Interface Transport

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

Para implantações baseadas em navegador ou quando quiser disponibilizar tanto os assets quanto o IPC por meio de seu transporte personalizado, implemente a interface `AssetServerTransport`:

```go
type AssetServerTransport interface {
    Transport

    // ServeAssets configures the transport to serve assets alongside IPC.
    // The assetHandler is Wails' internal asset server that handles all assets,
    // runtime.js, capabilities, flags, etc.
    ServeAssets(assetHandler http.Handler) error
}
```

**Quando implementar essa interface:**

- Executar a aplicação em um navegador em vez de uma webview
- Disponibilizar assets por HTTP junto com seu transporte IPC personalizado
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

- Todos os assets estáticos (HTML, CSS, JS, imagens etc.)
- `/wails/runtime.js` — A biblioteca de runtime do Wails

## Veja também

- `transport.go` — Interfaces e tipos principais de transporte
- `messageprocessor.go` — O processador de mensagens subjacente que gerencia todo o IPC do Wails
