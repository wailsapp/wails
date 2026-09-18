---
title: "流式事件"
description: "在 Go 中读取 SSE 并批量更新 Svelte，避免渲染过载"
slug: "features/events/streaming"
sourcePath: "features/events/streaming.md"
---

长时间运行的数据流往往会传来大量小更新。每收到一个片段就修改响应式状态，可能导致渲染过载。本指南在 Go 中读取外部服务的服务器发送事件（SSE）响应，并通过 `requestAnimationFrame` 合并 Svelte 5 的状态更新。

## 选择传输方式

对于新的持续点对点数据流，优先使用[原生 Streams](/guides/streams/)：注册 `app.HandleStream`，从 `c.Context()` 派生上游请求，检查 `c.SendJSON` 的错误，并使用 `@wailsio/runtime` 中的 `JSONStream` 接收对象。关闭连接会取消其上下文。`Stream` 默认传输字节，不会解析 SSE，也不会连接外部 URL。请将上游 HTTP 客户端留在 Go 端。传输差异请参阅[从 WebSocket 迁移到 Streams](/guides/streams-from-websockets/)。

下面的事件示例适合接入已有的事件驱动应用。它假设一个窗口中只有一个负责管理的组件，服务同时只有一个活动请求。应用事件会广播给监听器，并不是每个连接独享的私有数据流。按帧合并可减少界面工作，但不会为传输提供背压。对于持续的大流量或独立消费者，请使用原生 Streams，并保留相同的前端批处理方式。

## 后端服务

该服务从使用 LF 或 CRLF 换行的 UTF-8 SSE 响应中读取多行 `data` 字段。它启动 goroutine，使 `StartStream` 立即返回，并在新请求开始时取消上一个请求。请使用应用配置的可信端点；下面的 URL 只是占位示例。

这是一个有容量限制、只读取数据的读取器，并非完整的 `EventSource` 实现：它忽略 `event`、`id` 和 `retry`，不自动重连，并要求以空行结束后才派发每条记录。按照 [SSE 规范](https://html.spec.whatwg.org/multipage/server-sent-events.html#event-stream-interpretation)，到达 EOF 时未完成的记录会被丢弃。单行和累计记录的大小限制约为 1 MiB；如果提供方需要其他协议功能，请使用完整的 SSE 客户端。提供方特有的 JSON 或完成标记应先解析，再追加文本。

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

创建应用后注册服务：

```go title="main.go"
app := application.New(application.Options{
    Name: "Streaming Example",
})

app.RegisterService(application.NewService(NewStreamService(app)))
```

从 `app.Context()` 派生请求上下文，可确保应用退出时也停止请求。`StopStream` 只取消匹配的请求，因此延迟执行的清理不会停止替代请求。前端为每次请求提供新的流 ID，并忽略已替换或停止的请求中仍在传输的事件。这个服务在整个应用中只管理一个请求；若每个窗口需要自己的数据流，请使用独立的原生流连接。

## Svelte 组件

将收到的片段存入普通数组。最多只安排一个动画帧回调，然后在该帧中合并排队的片段并更新一次 Svelte 状态。组件会串行处理启动和停止调用，并处理绑定调用被拒绝的情况。示例使用 Svelte 5 的符文。

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

将 `changeme` 替换为应用的 Go 模块名，并在注册服务后生成绑定。确切的绑定路径取决于包含 `StreamService` 的包。将示例 URL 替换为你的 SSE 端点，并使用其实际响应格式进行测试。

监听器收到的是 `WailsEvent`，因此 Go 载荷位于 `event.data`。保留 `Events.On` 返回的取消订阅函数，并在组件销毁时调用。清理还会取消动画帧和匹配的后端请求；如果启动调用稍后才完成，会再次触发取消。

## 为什么按帧合并有效

接收数据时，组件在每个动画帧中最多对 `output` 赋值一次；停止时则显式刷新待处理文本。在 60 Hz 下，两帧之间收到的片段可以合并成一次响应式更新。渲染次数取决于数据到达时间、屏幕刷新率和框架，并不是每个响应都固定渲染几次。

隐藏窗口中的动画帧可能暂停。这个小示例没有限制等待数组和累计输出的大小。对于长期运行的数据流，请限制保留的历史和待处理数据，并决定是暂停生产者还是丢弃旧的显示数据。原生 Streams 提供传输背压，但任何传输方式都不会限制前端接收数据之后的应用状态大小。
