---
title: "Mensagens brutas"
description: "Implemente uma comunicação personalizada do frontend para o backend em aplicações com requisitos críticos de desempenho"
slug: "guides/raw-messages"
sourcePath: "guides/raw-messages.md"
---

As mensagens brutas fornecem um canal de comunicação de baixo nível entre o frontend e o backend, contornando o sistema padrão de bindings. Isso oferece mais velocidade em troca de conveniência.

## Quando usar mensagens brutas

As mensagens brutas são mais adequadas para casos extremos:

- **Atualizações em frequência extremamente alta** — milhares de mensagens por segundo, quando cada microssegundo importa
- **Protocolos de mensagens personalizados** — quando você precisa de controle total sobre o formato de transmissão

@note{type="tip"}
Para quase todos os casos de uso, os [bindings de serviço](/features/bindings/services/) padrão são recomendados, pois oferecem segurança de tipos, serialização automática e uma experiência de desenvolvimento melhor, com sobrecarga desprezível.

@end

## Configuração do backend

Configure `RawMessageHandler` nas opções da aplicação:

```go
package main

import (
    "encoding/json"
    "fmt"

    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "Raw Message Demo",
        Assets: application.AssetOptions{
            Handler: application.BundledAssetFileServer(assets),
        },
        RawMessageHandler: func(window application.Window, message string, originInfo *application.OriginInfo) {
            fmt.Printf("Raw message from window '%s': %s (origin: %+v)\n", window.Name(), message, originInfo.Origin)

            // Process the message and respond via events
            response := processMessage(message)
            window.EmitEvent("raw-response", response)
        },
    })

    app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title: "My App",
        Name:  "main",
    })

    app.Run()
}

func processMessage(message string) map[string]any {
    // Your custom message processing logic
    return map[string]any{
        "received": message,
        "status":   "processed",
    }
}
```

### Assinatura do manipulador

```go
RawMessageHandler func(window Window, message string, originInfo *application.OriginInfo)
```

| Parâmetro | Tipo | Descrição |
| --- | --- | --- |
| `window` | `Window` | A janela que enviou a mensagem |
| `message` | `string` | O conteúdo bruto da mensagem |
| `originInfo` | `*application.OriginInfo` | Informações de origem sobre a fonte da mensagem |

#### Estrutura OriginInfo

```go
type OriginInfo struct {
	Origin      string
	TopOrigin   string
	IsMainFrame bool
}
```

| Campo | Tipo | Descrição |
| --- | --- | --- |
| `Origin` | `string` | A URL de origem do documento que enviou a mensagem |
| `TopOrigin` | `string` | A URL de origem de nível superior (pode ser diferente de Origin em iframes) |
| `IsMainFrame` | `bool` | Indica se a mensagem se originou no frame principal |

#### Disponibilidade específica por plataforma

- **macOS**: `Origin` e `IsMainFrame` são fornecidos
- **Windows**: `Origin` e `TopOrigin` são fornecidos
- **Linux**: somente `Origin` é fornecido

### Validação da origem

@note{type="caution"}
Nunca presuma que uma mensagem é segura apenas porque chegou ao manipulador. As informações de origem devem ser validadas antes do processamento de operações sensíveis ou que modifiquem o estado.

@end

**Sempre verifique a origem das mensagens recebidas antes de processá-las.** O parâmetro `originInfo` fornece informações de segurança essenciais que devem ser validadas para impedir o acesso não autorizado. Conteúdo malicioso, conteúdo comprometido ou scripts não previstos podem enviar mensagens brutas. Sem a validação da origem, você pode processar comandos provenientes de fontes não confiáveis. Use `originInfo` para garantir que as mensagens venham das fontes esperadas.

### Principais pontos de validação

- **Sempre verifique `Origin`** — confirme se a origem corresponde às fontes confiáveis esperadas (normalmente `wails://wails` ou `http://wails.localhost` para recursos locais ou à origem específica da aplicação)
- **Valide `IsMainFrame`** (macOS) — verifique se a mensagem vem de um iframe, pois isso pode indicar conteúdo incorporado com contextos de segurança diferentes
- **Use `TopOrigin`** (Windows) — verifique a origem de nível superior ao lidar com conteúdo em frames
- **Rejeite origens inesperadas** — mantenha a segurança rejeitando mensagens de origens que você não permite explicitamente

@note{type="info"}
Mensagens prefixadas com `wails:` são reservadas para a comunicação interna do Wails e não serão encaminhadas ao seu manipulador.

@end

## Configuração do frontend

Envie mensagens brutas usando `System.invoke()`:

```html
<!DOCTYPE html>
<html>
<head>
    <script type="module">
        import { System, Events } from '@wailsio/runtime'

        // Send raw message
        document.getElementById('send').addEventListener('click', () => {
            const message = document.getElementById('input').value
            System.invoke(message)
        })

        // Listen for response
        Events.On('raw-response', (event) => {
            console.log('Response:', event.data)
        })
    </script>
</head>
<body>
    <input type="text" id="input" placeholder="Enter message" />
    <button id="send">Send</button>
</body>
</html>
```

### Como usar o bundle pré-compilado

Se você não usa npm, acesse `invoke` por meio do objeto global `wails`:

```html
<script type="module" src="/wails/runtime.js"></script>
<script>
    window.onload = function() {
        document.getElementById('send').onclick = function() {
            wails.System.invoke('my-message')
        }
    }
</script>
```

## Mensagens estruturadas

Para dados complexos, serialize-os como JSON:

### Frontend

```javascript
import { System } from '@wailsio/runtime'

const command = {
    action: 'update',
    payload: {
        id: 123,
        value: 'new value'
    }
}

System.invoke(JSON.stringify(command))
```

### Backend

```go
RawMessageHandler: func(window application.Window, message string, originInfo *application.OriginInfo) {
    var cmd struct {
        Action  string `json:"action"`
        Payload struct {
            ID    int    `json:"id"`
            Value string `json:"value"`
        } `json:"payload"`
    }

    if err := json.Unmarshal([]byte(message), &cmd); err != nil {
        window.EmitEvent("error", err.Error())
        return
    }

    switch cmd.Action {
    case "update":
        // Handle update
        result := handleUpdate(cmd.Payload.ID, cmd.Payload.Value)
        window.EmitEvent("update-complete", result)
    default:
        window.EmitEvent("error", "unknown action")
    }
}
```

## Comparação de desempenho

| Abordagem | Sobrecarga | Segurança de tipos | Caso de uso |
| --- | --- | --- | --- |
| Associações de serviços | Maior | Completa | Uso geral |
| Mensagens brutas | Mínima | Manual | Alta frequência, desempenho crítico |

### Exemplo de benchmark

Para cargas simples, as mensagens brutas podem processar significativamente mais mensagens por segundo do que as associações de serviços:

```go
// Raw message handler - minimal overhead
RawMessageHandler: func(window application.Window, message string, originInfo *application.OriginInfo) {
    // Direct string processing, no reflection or marshaling
    counter++
}
```

## Exemplo completo

Veja um exemplo completo que implementa um protocolo de comandos simples:

### main.go

```go
package main

import (
    "embed"
    "encoding/json"
    "fmt"
    "time"

    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets
var assets embed.FS

type Command struct {
    Type string          `json:"type"`
    Data json.RawMessage `json:"data"`
}

func main() {
    app := application.New(application.Options{
        Name: "Raw Message Demo",
        Assets: application.AssetOptions{
            Handler: application.BundledAssetFileServer(assets),
        },
        Mac: application.MacOptions{
            ApplicationShouldTerminateAfterLastWindowClosed: true,
        },
        RawMessageHandler: func(window application.Window, message string, originInfo *application.OriginInfo) {
            var cmd Command
            if err := json.Unmarshal([]byte(message), &cmd); err != nil {
                window.EmitEvent("error", map[string]string{"error": err.Error()})
                return
            }

            switch cmd.Type {
            case "ping":
                window.EmitEvent("pong", map[string]any{
                    "time":   time.Now().UnixMilli(),
                    "window": window.Name(),
                })
            case "echo":
                var text string
                json.Unmarshal(cmd.Data, &text)
                window.EmitEvent("echo", text)
            default:
                window.EmitEvent("error", map[string]string{
                    "error": fmt.Sprintf("unknown command: %s", cmd.Type),
                })
            }
        },
    })

    app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title: "Raw Message Demo",
        Name:  "main",
        Width: 400,
        Height: 300,
    })

    app.Run()
}
```

### assets/index.html

```html
<!DOCTYPE html>
<html>
<head>
    <title>Raw Message Demo</title>
    <style>
        body { font-family: sans-serif; padding: 20px; }
        button { margin: 5px; padding: 10px 20px; }
        #output { margin-top: 20px; padding: 10px; background: #f0f0f0; }
    </style>
</head>
<body>
    <h1>Raw Message Demo</h1>

    <button id="ping">Ping</button>
    <button id="echo">Echo "Hello"</button>

    <div id="output">Waiting for response...</div>

    <script type="module">
        import { System, Events } from '@wailsio/runtime'

        const output = document.getElementById('output')

        function send(type, data) {
            System.invoke(JSON.stringify({ type, data }))
        }

        document.getElementById('ping').onclick = () => send('ping')
        document.getElementById('echo').onclick = () => send('echo', 'Hello')

        Events.On('pong', (e) => {
            output.textContent = `Pong from ${e.data.window} at ${e.data.time}`
        })

        Events.On('echo', (e) => {
            output.textContent = `Echo: ${e.data}`
        })

        Events.On('error', (e) => {
            output.textContent = `Error: ${e.data.error}`
        })
    </script>
</body>
</html>
```

## Práticas recomendadas

### O que fazer

- Use mensagens brutas em fluxos que realmente tenham desempenho crítico
- Implemente o tratamento adequado de erros no manipulador
- Use eventos para enviar respostas de volta ao frontend
- Considere usar JSON para dados estruturados
- Mantenha o processamento de mensagens rápido para evitar bloqueios

### O que não fazer

- Não use mensagens brutas quando as associações de serviços forem suficientes
- Não se esqueça de validar as mensagens recebidas
- Não bloqueie o manipulador com operações demoradas (use goroutines)
- Não ignore o parâmetro de janela quando as respostas precisarem ser direcionadas a janelas específicas

## Considerações para várias janelas

O parâmetro `window` identifica qual janela enviou a mensagem, permitindo que você:

- Envie respostas à janela correta
- Implemente comportamentos específicos para cada janela
- Rastreie as origens das mensagens para depuração

```go
RawMessageHandler: func(window application.Window, message string, originInfo *application.OriginInfo) {
    // Respond only to the sending window
    window.EmitEvent("response", result)

    // Or broadcast to all windows
    app.Event.Emit("broadcast", result)
}
```

## Próximas etapas

- [Associações de serviços](/features/bindings/services/) - Abordagem padrão para a maioria das aplicações
- [Eventos](/guides/events-reference/) - Sistema de eventos para comunicação do backend com o frontend
- [Desempenho](/guides/performance/) - Otimização geral de desempenho
