---
title: "Build para servidor"
description: "Execute aplicações Wails como servidores HTTP sem uma janela de interface gráfica nativa"
slug: "guides/server-build"
sourcePath: "guides/server-build.md"
---

O Wails v3 oferece suporte ao modo servidor, permitindo executar sua aplicação como um servidor HTTP puro, sem criar janelas nativas nem exigir dependências de interface gráfica. Isso permite implantar a mesma aplicação Wails em servidores, contêineres e navegadores web.

O modo servidor é útil para:

- **Implantações em Docker/contêineres** — execute sem dependências do X11/Wayland
- **Aplicações do lado do servidor** — implante como um servidor web acessível pelo navegador
- **Acesso somente pela web** — compartilhe a mesma base de código entre desktop e web
- **Testes de CI/CD** — execute testes de integração sem um servidor de exibição
- **Microsserviços** — use bindings do Wails em serviços de backend sem interface gráfica

## Início rápido

O modo servidor é ativado pela tag de build `server`. O código da aplicação permanece o mesmo — basta fazer o build com a tag:

```bash
# Using Taskfile (recommended)
wails3 task build:server
wails3 task run:server

# Or build directly with Go
go build -tags server -o myapp-server .
```

Veja um exemplo mínimo:

```go
package main

import (
    "embed"
    "log"

    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
    app := application.New(application.Options{
        Name: "My App",
        // Server options are used when built with -tags server
        Server: application.ServerOptions{
            Host: "localhost",
            Port: 8080,
        },
        Services: []application.Service{
            application.NewService(&MyService{}),
        },
        Assets: application.AssetOptions{
            Handler: application.AssetFileServerFS(assets),
        },
    })

    log.Println("Starting application...")
    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}
```

O mesmo código pode ser compilado para o modo desktop (sem a tag) ou para o modo servidor (com `-tags server`).

## Configuração

### ServerOptions

Configure o servidor HTTP com `ServerOptions`:

```go
Server: application.ServerOptions{
    // Host to bind to. Default: "localhost"
    // Use "0.0.0.0" to listen on all interfaces
    Host: "localhost",

    // Port to listen on. Default: 8080
    Port: 8080,

    // Request read timeout. Default: 30s
    ReadTimeout: 30 * time.Second,

    // Response write timeout. Default: 30s
    WriteTimeout: 30 * time.Second,

    // Idle connection timeout. Default: 120s
    IdleTimeout: 120 * time.Second,

    // Graceful shutdown timeout. Default: 30s
    ShutdownTimeout: 30 * time.Second,

    // Additional origins allowed to open WebSocket connections.
    // Same-origin connections are always allowed.
    WebSocketOriginPatterns: []string{"app.example.com"},

    // Disable WebSocket origin checks. Unsafe; default: false.
    WebSocketAllowAllOrigins: false,

    // TLS configuration (optional)
    TLS: &application.TLSOptions{
        CertFile: "/path/to/cert.pem",
        KeyFile:  "/path/to/key.pem",
    },
},
```

## Recursos

### Endpoint de verificação de integridade

Um endpoint de verificação de integridade fica disponível automaticamente em `/health`:

```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

Isso é útil para:

- Probes de atividade e prontidão do Kubernetes
- Verificações de integridade do balanceador de carga
- Sistemas de monitoramento

### Bindings de serviços

Todos os bindings de serviços funcionam da mesma forma que no modo desktop:

```go
type GreetService struct{}

func (g *GreetService) Greet(name string) string {
    return "Hello, " + name + "!"
}

// Register in options
Services: []application.Service{
    application.NewService(&GreetService{}),
},
```

O frontend pode chamar esses bindings usando o runtime padrão do Wails:

```javascript
const greeting = await wails.Call.ByName('main.GreetService.Greet', 'World');
```

### Eventos

No modo servidor, os eventos funcionam de forma bidirecional:

- **Do frontend para o backend**: os eventos emitidos pelo navegador são enviados via HTTP e recebidos pelos seus handlers de eventos em Go
- **Do backend para o frontend**: os eventos emitidos pelo Go são transmitidos a todos os navegadores conectados via WebSocket

Cada aba do navegador é representada como uma "janela" com um nome exclusivo (`browser-1`, `browser-2` etc.), acessível por meio de `event.Sender`:

```go
// Listen for events from browsers
app.Event.On("user-action", func(event *application.CustomEvent) {
    log.Printf("Event from %s: %v", event.Sender, event.Data)
    // event.Sender will be "browser-1", "browser-2", etc.
})

// Emit events to all connected browsers
app.Event.Emit("server-update", data)
```

No frontend:

```javascript
// Emit event to server (and all other browsers)
await wails.Events.Emit('user-action', { action: 'click' });

// Listen for events from server
wails.Events.On('server-update', (event) => {
    console.log('Update from server:', event.data);
});
```

### Encerramento gracioso

O servidor trata os sinais `SIGINT` e `SIGTERM` de forma graciosa:

1. Para de aceitar novas conexões
2. Aguarda a conclusão das requisições ativas (por até `ShutdownTimeout`)
3. Executa os hooks `OnShutdown`
4. Encerra os serviços na ordem inversa

## Diferenças em relação ao modo desktop

| Recurso | Modo desktop | Modo servidor |
| --- | --- | --- |
| Janelas nativas | Criadas | Janelas do navegador (`browser-N`) |
| Bandeja do sistema | Disponível | Indisponível |
| Caixas de diálogo nativas | Disponíveis | Indisponíveis |
| Menu da aplicação | Disponível | Indisponível |
| Informações da tela | Disponíveis | Retorna erro |
| Bindings de serviços | Funciona | Funciona |
| Eventos | Funcionam | Funcionam (via WebSocket) |
| Recursos estáticos | Via webview | Via HTTP |
| CGO obrigatório | Sim | Não |

### Comportamento da API de janelas

No modo servidor, as APIs relacionadas a janelas são tratadas com segurança:

- `app.Window.NewWithOptions()` — registra um aviso e retorna nil
- `app.Hide()` / `app.Show()` — não realiza nenhuma operação
- `app.Screen.GetPrimary()` — retorna um erro

Isso permite que o código que referencia janelas seja executado sem falhar, embora as operações de janela não tenham efeito.

## Build para produção

### Usando o Task (recomendado)

Os projetos criados com `wails3 init` incluem uma tarefa `build:server`:

```bash
# Build for server mode
wails3 task build:server

# Build and run
wails3 task run:server
```

### Compilação manual

```bash
# Build with server mode
go build -tags server -o myapp-server .
```

### Docker

Os projetos Wails incluem uma configuração do Docker pronta para uso. Para compilar e executar seu aplicativo em um contêiner:

```bash
# Build the Docker image
wails3 task build:docker

# Run it
wails3 task run:docker
```

Pronto! Seu aplicativo estará disponível em `http://localhost:8080`.

Você pode personalizar a compilação com algumas opções:

```bash
# Use a custom image tag
wails3 task build:docker TAG=myapp:v1.0.0

# Run on a different port
wails3 task run:docker PORT=3000
```

O `Dockerfile.server` gerado cria uma imagem mínima baseada em distroless. Ele gerencia o vínculo de rede automaticamente, portanto seu aplicativo poderá ser acessado de fora do contêiner.

### Docker Compose

Para implantações mais complexas, veja uma configuração do Docker Compose com verificações de integridade:

```yaml
services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - WAILS_SERVER_HOST=0.0.0.0
    healthcheck:
      test: ["CMD", "wget", "-q", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
```

@note{type="info"}
O exemplo de verificação de integridade usa `wget`. Se você estiver usando uma imagem base distroless, precisará incluir um binário de verificação de integridade na imagem ou usar um mecanismo externo de verificação de integridade (por exemplo, a opção `curl` do Docker ou um contêiner sidecar).

@end

### Dockerfile personalizado

Se precisar de mais controle, você poderá criar seu próprio Dockerfile. O principal é definir `WAILS_SERVER_HOST=0.0.0.0` para que o servidor aceite conexões de fora do contêiner:

```dockerfile
# Build stage
FROM golang:alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY . .
RUN go mod tidy
RUN go build -tags server -ldflags="-s -w" -o server .

# Runtime stage
FROM gcr.io/distroless/static-debian12
COPY --from=builder /app/server /server
COPY --from=builder /app/frontend/dist /frontend/dist
EXPOSE 8080
ENV WAILS_SERVER_HOST=0.0.0.0
ENTRYPOINT ["/server"]
```

## Considerações de segurança

Ao implantar aplicativos no modo servidor:

1. **Vincule ao localhost por padrão** — Use `0.0.0.0` somente quando necessário
2. **Use TLS em produção** — Configure `ServerOptions.TLS`
3. **Coloque atrás de um proxy reverso** — Use nginx/traefik para obter segurança adicional
4. **Mantenha os WebSockets na mesma origem** — Adicione somente origens confiáveis com `WebSocketOriginPatterns`; evite `WebSocketAllowAllOrigins`
5. **Valide todas as entradas** — Adote as mesmas práticas de segurança usadas em qualquer aplicativo web

## Exemplo

Um exemplo completo está disponível em `v3/examples/server/`:

```bash
cd v3/examples/server

# Using Taskfile
task dev

# Or run directly
go run -tags server .

# Open http://localhost:8080 in browser
```

## Variáveis de ambiente

Para cenários de implantação nos quais você precisa substituir a configuração do servidor sem alterar o código, o Wails reconhece estas variáveis de ambiente:

| Variável | Descrição | Padrão |
| --- | --- | --- |
| `WAILS_SERVER_HOST` | Interface de rede à qual vincular | `localhost` |
| `WAILS_SERVER_PORT` | Porta na qual escutar | `8080` |

Essas variáveis têm precedência sobre `ServerOptions` no seu código. Por isso, os exemplos do Docker definem `WAILS_SERVER_HOST=0.0.0.0`: isso permite que o contêiner aceite conexões externas sem exigir alterações no aplicativo.

## Veja também

- [Transporte personalizado](/guides/custom-transport/) — Para personalização avançada da IPC
- [Serviços](/features/bindings/services/) — Documentação sobre vinculação de serviços
- [Eventos](/guides/events-reference/) — Documentação do sistema de eventos

### Tamanho das solicitações do runtime

As solicitações para `/wails/runtime` são limitadas a 64 MiB antes do processamento do JSON. Uma solicitação comum acima desse limite recebe HTTP 413, inclusive solicitações sem `Content-Length`. Os uploads fragmentados do runtime mantêm seus limites separados de 1 MiB por fragmento e 64 MiB para a carga útil montada. Use um middleware do aplicativo ou um proxy reverso para impor um limite menor quando apropriado.
