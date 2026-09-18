---
title: "Streaming Events"
description: "Read SSE in Go and batch Svelte updates without overwhelming rendering"
slug: "features/events/streaming"
sourcePath: "features/events/streaming.md"
---

Long-running streams often deliver many small updates. Updating reactive state for every chunk can overwhelm rendering. This recipe reads an external server-sent events (SSE) response in Go and batches Svelte 5 state updates with `requestAnimationFrame`.

## Choose the transport

For a new continuous point-to-point feed, prefer [native Streams](/guides/streams/): register `app.HandleStream`, derive the upstream request from `c.Context()`, check errors from `c.SendJSON`, and receive objects with `JSONStream` from `@wailsio/runtime`. Closing that connection cancels its context. `Stream` delivers bytes by default; it does not parse SSE or connect to an external URL. Keep the upstream HTTP client in Go. See [Migrating a WebSocket to Streams](/guides/streams-from-websockets/) for the transport differences.

The events-based example below is useful when integrating an existing event-driven application. It assumes one owning component in one window and one active request for the service. Application events broadcast to listeners; they are not private per-connection streams. Frame batching reduces UI work but adds no transport backpressure. For sustained high-volume data or independent consumers, use native Streams and retain the same frontend batching technique.

## Backend service

The service reads multiline `data` fields from a UTF-8 SSE response using LF or CRLF line endings. It starts a goroutine so `StartStream` returns immediately and cancels the previous request when a new one starts. Use a trusted endpoint configured by your application; the URL below is a placeholder.

This is a bounded data-only reader, not a complete `EventSource` implementation: it ignores `event`, `id`, and `retry`, does not reconnect, and requires a blank line to dispatch each record. An unfinished record at EOF is discarded, as required by the [SSE specification](https://html.spec.whatwg.org/multipage/server-sent-events.html#event-stream-interpretation). Individual lines and accumulated records are limited to about 1 MiB; use a full SSE client if your provider needs other protocol features. Provider-specific JSON or completion markers must be interpreted before appending text.

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

Register the service after creating the application:

```go title="main.go"
app := application.New(application.Options{
    Name: "Streaming Example",
})

app.RegisterService(application.NewService(NewStreamService(app)))
```

Deriving the request context from `app.Context()` ensures that application shutdown also stops the request. `StopStream` cancels only the matching request, so delayed cleanup cannot stop a replacement. The frontend supplies a fresh stream ID and ignores events already in transit from replaced or stopped requests. This service owns one request across the application; use separate native stream connections when each window needs its own feed.

## Svelte component

Keep incoming chunks in a plain array. Schedule at most one animation frame, then join the queued chunks and update Svelte state once in that frame. The component serialises its start and stop calls and handles rejected bindings. The example uses Svelte 5 runes.

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

Replace `changeme` with your application's Go module name and generate the bindings after registering the service. The exact binding path depends on the package containing `StreamService`. Replace the example URL with your SSE endpoint and test against its actual response format.

The event listener receives a `WailsEvent`, so the Go payload is available on `event.data`. Keep the unsubscribe function returned by `Events.On` and call it when the component is destroyed. Cleanup also cancels the animation frame and the matching backend request; a late start completion triggers cancellation again.

## Why frame batching helps

While streaming, the component performs at most one `output` assignment per animation frame; stopping explicitly flushes pending text. At 60 Hz, chunks received between frames can become one reactive update. The number of renders depends on arrival timing, display refresh rate, and the framework; it is not a fixed number per response.

Animation frames may pause in a hidden window. The pending array and accumulated output in this small example are not bounded. For long-lived feeds, limit retained history and pending data, and decide whether to pause the producer or discard old display data. Native Streams provide transport backpressure, but neither transport bounds the application state after the frontend receives it.
