---
title: "Eventos em fluxo"
description: "Leia SSE em Go e agrupe atualizações do Svelte sem sobrecarregar a renderização"
slug: "features/events/streaming"
sourcePath: "features/events/streaming.md"
---

Fluxos de longa duração costumam entregar muitas atualizações pequenas. Alterar o estado reativo a cada fragmento pode sobrecarregar a renderização. Esta receita lê em Go uma resposta externa de eventos enviados pelo servidor (SSE) e agrupa as atualizações de estado do Svelte 5 com `requestAnimationFrame`.

## Escolha o transporte

Para um novo fluxo contínuo ponto a ponto, prefira os [Streams nativos](/guides/streams/): registre `app.HandleStream`, derive a requisição ao serviço externo de `c.Context()`, verifique os erros de `c.SendJSON` e receba objetos com `JSONStream` de `@wailsio/runtime`. Fechar essa conexão cancela seu contexto. `Stream` entrega bytes por padrão; não interpreta SSE nem se conecta a uma URL externa. Mantenha o cliente HTTP externo em Go. Consulte [Migrando um WebSocket para Streams](/guides/streams-from-websockets/) para conhecer as diferenças de transporte.

O exemplo baseado em eventos abaixo é útil ao integrar um aplicativo existente orientado a eventos. Ele pressupõe um único componente responsável em uma janela e uma requisição ativa no serviço. Eventos do aplicativo são transmitidos aos ouvintes; não são fluxos privados por conexão. Agrupar por quadro reduz o trabalho da interface, mas não controla a pressão no transporte. Para volumes altos e contínuos ou consumidores independentes, use Streams nativos e mantenha a mesma técnica de agrupamento no frontend.

## Serviço backend

O serviço lê campos `data` com várias linhas de uma resposta SSE UTF-8 com terminações LF ou CRLF. Ele inicia uma goroutine para que `StartStream` retorne imediatamente e cancela a requisição anterior quando uma nova começa. Use um endpoint confiável configurado pelo aplicativo; a URL abaixo é ilustrativa.

Este leitor limitado apenas a dados não é uma implementação completa de `EventSource`: ignora `event`, `id` e `retry`, não reconecta e exige uma linha em branco para despachar cada registro. Um registro incompleto no fim do arquivo é descartado, conforme a [especificação SSE](https://html.spec.whatwg.org/multipage/server-sent-events.html#event-stream-interpretation). Linhas individuais e registros acumulados são limitados a cerca de 1 MiB; use um cliente SSE completo se o provedor exigir outros recursos do protocolo. JSON ou marcadores de conclusão específicos do provedor devem ser interpretados antes de acrescentar texto.

```go title="stream_service.go"
package main

import (
    "bufio"
    "context"
    "fmt"
    "mime"
    "net/http"
    "strings"
    "sync"

    "github.com/wailsapp/wails/v3/pkg/application"
)

type StreamEvent struct {
    StreamID string `json:"streamId"`
    Content  string `json:"content,omitempty"`
    Done     bool   `json:"done,omitempty"`
    Error    string `json:"error,omitempty"`
}

type StreamService struct {
    app *application.App

    mu         sync.Mutex
    cancel     context.CancelFunc
    generation uint64
    streamID   string
}

func NewStreamService(app *application.App) *StreamService {
    return &StreamService{app: app}
}

func (s *StreamService) StartStream(url string, streamID string) {
    s.mu.Lock()
    if s.cancel != nil {
        s.cancel()
    }

    ctx, cancel := context.WithCancel(s.app.Context())
    s.cancel = cancel
    s.streamID = streamID
    s.generation++
    generation := s.generation
    s.mu.Unlock()

    go s.readStream(ctx, cancel, generation, streamID, url)
}

func (s *StreamService) StopStream(streamID string) {
    s.mu.Lock()
    defer s.mu.Unlock()

    if s.cancel != nil && s.streamID == streamID {
        s.cancel()
        s.cancel = nil
    }
}

func (s *StreamService) readStream(
    ctx context.Context,
    cancel context.CancelFunc,
    generation uint64,
    streamID string,
    url string,
) {
    defer s.clearStream(cancel, generation)

    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    if err != nil {
        s.emitError(ctx, streamID, err)
        return
    }
    req.Header.Set("Accept", "text/event-stream")

    response, err := http.DefaultClient.Do(req)
    if err != nil {
        s.emitError(ctx, streamID, err)
        return
    }
    defer response.Body.Close()

    if response.StatusCode != http.StatusOK {
        s.emitError(ctx, streamID, fmt.Errorf("stream returned %s", response.Status))
        return
    }

    mediaType, _, err := mime.ParseMediaType(response.Header.Get("Content-Type"))
    if err != nil || mediaType != "text/event-stream" {
        s.emitError(ctx, streamID, fmt.Errorf("expected text/event-stream"))
        return
    }

    scanner := bufio.NewScanner(response.Body)
    scanner.Buffer(make([]byte, 64*1024), 1024*1024)

    var dataLines []string
    dataBytes := 0
    emitData := func() {
        if len(dataLines) == 0 || ctx.Err() != nil {
            dataLines = dataLines[:0]
            return
        }

        s.app.Event.Emit("stream:chunk", StreamEvent{
            StreamID: streamID,
            Content:  strings.Join(dataLines, "\n"),
        })
        dataLines = dataLines[:0]
        dataBytes = 0
    }

    firstLine := true
    for scanner.Scan() {
        line := scanner.Text()
        if firstLine {
            line = strings.TrimPrefix(line, "\uFEFF")
            firstLine = false
        }
        if line == "" {
            emitData()
            continue
        }

        field, value, found := strings.Cut(line, ":")
        if !found {
            field = line
            value = ""
        }
        if field == "data" {
            value = strings.TrimPrefix(value, " ")
            dataBytes += len(value) + 1
            if dataBytes > 1024*1024 {
                s.emitError(ctx, streamID, fmt.Errorf("SSE record exceeds 1 MiB"))
                return
            }
            dataLines = append(dataLines, value)
        }
    }

    if err := scanner.Err(); err != nil {
        s.emitError(ctx, streamID, err)
        return
    }
    if ctx.Err() == nil {
        s.app.Event.Emit("stream:chunk", StreamEvent{
            StreamID: streamID,
            Done:     true,
        })
    }
}

func (s *StreamService) emitError(
    ctx context.Context,
    streamID string,
    err error,
) {
    // Cancellation is expected when a stream is stopped or replaced.
    if ctx.Err() == nil {
        s.app.Event.Emit("stream:chunk", StreamEvent{
            StreamID: streamID,
            Error:    err.Error(),
        })
    }
}

func (s *StreamService) clearStream(
    cancel context.CancelFunc,
    generation uint64,
) {
    cancel()

    s.mu.Lock()
    defer s.mu.Unlock()

    // Do not clear the cancel function belonging to a newer stream.
    if s.generation == generation {
        s.cancel = nil
    }
}
```

Registre o serviço depois de criar o aplicativo:

```go title="main.go"
app := application.New(application.Options{
    Name: "Streaming Example",
})

app.RegisterService(application.NewService(NewStreamService(app)))
```

Derivar o contexto da requisição de `app.Context()` garante que encerrar o aplicativo também interrompa a requisição. `StopStream` cancela apenas a requisição correspondente, portanto uma limpeza atrasada não interrompe sua substituta. O frontend fornece um novo ID de fluxo e ignora eventos ainda em trânsito de requisições substituídas ou interrompidas. Este serviço controla uma requisição para todo o aplicativo; use conexões nativas separadas quando cada janela precisar do próprio fluxo.

## Componente Svelte

Mantenha os fragmentos recebidos em um array comum. Agende no máximo um quadro de animação; depois una os fragmentos pendentes e atualize o estado do Svelte uma vez nesse quadro. O componente serializa as chamadas de início e parada e trata rejeições dos bindings. O exemplo usa runas do Svelte 5.

```svelte title="StreamOutput.svelte"
<script lang="ts">
    import { onMount } from 'svelte'
    import { Events } from '@wailsio/runtime'
    import {
        StartStream,
        StopStream,
    } from '../bindings/changeme/streamservice'

    type StreamEvent = {
        streamId: string
        content?: string
        done?: boolean
        error?: string
    }

    let output = $state('')
    let status = $state('idle')
    let pendingChunks: string[] = []
    let frame: number | null = null
    let streamId = ''
    let commandPending = $state(false)
    let disposed = false

    function flush() {
        frame = null
        if (pendingChunks.length === 0) {
            return
        }

        output += pendingChunks.join('')
        pendingChunks = []
    }

    function scheduleFlush() {
        if (frame === null) {
            frame = requestAnimationFrame(flush)
        }
    }

    async function start() {
        if (commandPending || disposed) return
        commandPending = true
        streamId = crypto.getRandomValues(new Uint32Array(4)).join('-')
        const currentStreamId = streamId

        if (frame !== null) {
            cancelAnimationFrame(frame)
            frame = null
        }
        pendingChunks = []
        output = ''
        status = 'streaming'
        try {
            await StartStream('https://example.com/events', currentStreamId)
        } catch (error) {
            if (!disposed && streamId === currentStreamId) {
                streamId = ''
                status = `error: ${String(error)}`
            }
        } finally {
            commandPending = false
            // A start binding can finish after the component has unmounted.
            if (disposed) {
                void StopStream(currentStreamId).catch(console.error)
            }
        }
    }

    async function stop() {
        if (commandPending || disposed) return
        commandPending = true
        const stoppedStreamId = streamId
        streamId = ''

        if (frame !== null) {
            cancelAnimationFrame(frame)
            flush()
        }
        status = 'stopping'
        try {
            await StopStream(stoppedStreamId)
            if (!disposed) status = 'stopped'
        } catch (error) {
            if (!disposed) status = `error: ${String(error)}`
        } finally {
            commandPending = false
        }
    }

    onMount(() => {
        const unsubscribe = Events.On('stream:chunk', (event) => {
            const update = event.data as StreamEvent

            if (update.streamId !== streamId) {
                return
            }
            if (update.content) {
                pendingChunks.push(update.content)
                scheduleFlush()
            }
            if (update.error) {
                status = `error: ${update.error}`
            } else if (update.done) {
                status = 'complete'
            }
        })

        return () => {
            disposed = true
            const stoppedStreamId = streamId
            streamId = ''
            unsubscribe()
            void StopStream(stoppedStreamId).catch(console.error)

            if (frame !== null) {
                cancelAnimationFrame(frame)
            }
            pendingChunks = []
        }
    })
</script>

<button onclick={start} disabled={commandPending || status === 'streaming'}>
    Start stream
</button>
<button onclick={stop} disabled={commandPending || status !== 'streaming'}>
    Stop stream
</button>

<p>{status}</p>
<pre>{output}</pre>
```

Substitua `changeme` pelo nome do módulo Go do aplicativo e gere os bindings após registrar o serviço. O caminho exato depende do pacote que contém `StreamService`. Substitua a URL de exemplo pelo seu endpoint SSE e teste o formato real da resposta.

O ouvinte recebe um `WailsEvent`, portanto os dados enviados pelo Go estão em `event.data`. Guarde a função de cancelamento da inscrição retornada por `Events.On` e chame-a quando o componente for destruído. A limpeza também cancela o quadro de animação e a requisição backend correspondente; uma conclusão tardia do início dispara o cancelamento novamente.

## Por que agrupar por quadro ajuda

Durante a recepção, o componente faz no máximo uma atribuição a `output` por quadro de animação; ao parar, esvazia explicitamente o texto pendente. A 60 Hz, fragmentos recebidos entre quadros podem se tornar uma única atualização reativa. A quantidade de renderizações depende do momento de chegada, da taxa de atualização da tela e do framework; não é fixa por resposta.

Quadros de animação podem pausar em uma janela oculta. O array de espera e a saída acumulada deste pequeno exemplo não são limitados. Para fluxos duradouros, limite o histórico e os dados pendentes e decida se deve pausar o produtor ou descartar dados antigos da exibição. Streams nativos oferecem controle de pressão no transporte, mas nenhum transporte limita o estado do aplicativo depois que o frontend recebe os dados.
