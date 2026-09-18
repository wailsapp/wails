---
title: "串流事件"
description: "在 Go 中讀取 SSE 並批次更新 Svelte，避免算繪過載"
slug: "features/events/streaming"
sourcePath: "features/events/streaming.md"
---

長時間執行的串流往往會傳來大量小幅更新。每收到一個片段就修改反應式狀態，可能導致算繪過載。本指南在 Go 中讀取外部服務的伺服器傳送事件（SSE）回應，並透過 `requestAnimationFrame` 合併 Svelte 5 的狀態更新。

## 選擇傳輸方式

對於新的持續點對點資料流，優先使用[原生 Streams](/guides/streams/)：註冊 `app.HandleStream`，從 `c.Context()` 衍生上游請求，檢查 `c.SendJSON` 的錯誤，並使用 `@wailsio/runtime` 中的 `JSONStream` 接收物件。關閉連線會取消其上下文。`Stream` 預設傳輸位元組，不會解析 SSE，也不會連接外部 URL。請將上游 HTTP 用戶端留在 Go 端。傳輸差異請參閱[從 WebSocket 遷移至 Streams](/guides/streams-from-websockets/)。

以下事件範例適合整合到既有的事件驅動應用程式。它假設一個視窗中只有一個負責管理的元件，服務同時只有一個有效請求。應用程式事件會廣播給監聽器，並非每個連線專用的私有串流。按畫面影格合併可減少介面工作，但不會為傳輸提供背壓。對於持續的大流量或獨立消費者，請使用原生 Streams，並保留相同的前端批次處理方式。

## 後端服務

服務從使用 LF 或 CRLF 換行的 UTF-8 SSE 回應中讀取多行 `data` 欄位。它啟動 goroutine，讓 `StartStream` 立即返回，並在新請求開始時取消上一個請求。請使用應用程式設定的可信任端點；以下 URL 只是需要替換的範例。

這是一個有容量限制、僅讀取資料的讀取器，並非完整的 `EventSource` 實作：它忽略 `event`、`id` 與 `retry`，不自動重新連線，且必須以空行結束才能派送每筆記錄。依照 [SSE 規範](https://html.spec.whatwg.org/multipage/server-sent-events.html#event-stream-interpretation)，到達 EOF 時未完成的記錄會被捨棄。單行及累積記錄的大小限制約為 1 MiB；若提供者需要其他協定功能，請使用完整的 SSE 用戶端。提供者特有的 JSON 或完成標記應先解析，再附加文字。

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

建立應用程式後註冊服務：

```go title="main.go"
app := application.New(application.Options{
    Name: "Streaming Example",
})

app.RegisterService(application.NewService(NewStreamService(app)))
```

從 `app.Context()` 衍生請求上下文，可確保應用程式結束時也停止請求。`StopStream` 僅取消相符的請求，因此延遲執行的清理不會停止替代請求。前端為每次請求提供新的串流 ID，並忽略已替換或停止的請求中仍在傳輸的事件。這個服務在整個應用程式中只管理一個請求；若每個視窗需要自己的資料流，請使用獨立的原生串流連線。

## Svelte 元件

將收到的片段保存在一般陣列中。最多只安排一個動畫影格回呼，然後在該影格中合併排隊的片段並更新一次 Svelte 狀態。元件會依序處理啟動與停止呼叫，並處理繫結呼叫遭拒的情況。範例使用 Svelte 5 的符文。

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

將 `changeme` 替換為應用程式的 Go 模組名稱，並在註冊服務後產生繫結。確切路徑取決於包含 `StreamService` 的套件。將範例 URL 替換為你的 SSE 端點，並使用其實際回應格式測試。

監聽器收到的是 `WailsEvent`，因此 Go 酬載位於 `event.data`。保留 `Events.On` 傳回的取消訂閱函式，並在元件銷毀時呼叫。清理也會取消動畫影格與相符的後端請求；若啟動呼叫稍後才完成，會再次觸發取消。

## 為何按影格合併有效

接收資料時，元件在每個動畫影格中最多對 `output` 賦值一次；停止時則明確刷新待處理文字。在 60 Hz 下，兩個影格之間收到的片段可以合併成一次反應式更新。算繪次數取決於資料抵達時間、螢幕更新率與框架，並非每個回應都有固定次數。

隱藏視窗中的動畫影格可能暫停。這個小範例並未限制等待陣列與累積輸出的大小。對於長期執行的資料流，請限制保留的歷史與待處理資料，並決定要暫停產生者或捨棄舊的顯示資料。原生 Streams 提供傳輸背壓，但任何傳輸方式都不會限制前端收到資料之後的應用程式狀態大小。
