---
title: "ストリーミングイベント"
description: "Go で SSE を読み取り、描画を圧迫せずに Svelte の更新をまとめる"
slug: "features/events/streaming"
sourcePath: "features/events/streaming.md"
---

長時間続くストリームでは、小さな更新が大量に届くことがあります。断片ごとにリアクティブな状態を更新すると、描画の負荷が高くなります。このレシピでは、外部サービスのサーバー送信イベント（SSE）応答を Go で読み取り、`requestAnimationFrame` を使って Svelte 5 の状態更新をまとめます。

## トランスポートを選ぶ

新しく継続的な一対一のデータ配信を作る場合は、[ネイティブ Streams](/guides/streams/) を推奨します。`app.HandleStream` を登録し、`c.Context()` から上流へのリクエストを作り、`c.SendJSON` のエラーを確認し、`@wailsio/runtime` の `JSONStream` でオブジェクトを受信します。接続を閉じると、そのコンテキストがキャンセルされます。`Stream` は標準ではバイト列を渡すもので、SSE の解析や外部 URL への接続は行いません。上流サービス用の HTTP クライアントは Go 側に置いてください。トランスポートの違いは [WebSocket から Streams への移行](/guides/streams-from-websockets/) を参照してください。

以下のイベント方式の例は、既存のイベント駆動アプリケーションへの組み込みに役立ちます。一つのウィンドウに管理するコンポーネントが一つあり、サービスが同時に扱うリクエストも一つであることを前提とします。アプリケーションイベントはリスナーに配信され、接続ごとの専用ストリームではありません。フレーム単位の集約は UI の負荷を減らしますが、トランスポートにバックプレッシャーを追加しません。継続的に大量のデータを流す場合や受信側を独立させる場合は、ネイティブ Streams と同じフロントエンドの集約手法を組み合わせてください。

## バックエンドサービス

サービスは、改行が LF または CRLF の UTF-8 SSE 応答から複数行の `data` フィールドを読み取ります。goroutine を起動するため `StartStream` はすぐに戻り、新しいリクエストの開始時には前のリクエストをキャンセルします。アプリケーションで設定した信頼できるエンドポイントを使ってください。以下の URL は置き換え用の例です。

これは容量を制限したデータ専用のリーダーであり、完全な `EventSource` 実装ではありません。`event`、`id`、`retry` を無視し、再接続せず、各レコードの配信には空行を必要とします。[SSE 仕様](https://html.spec.whatwg.org/multipage/server-sent-events.html#event-stream-interpretation) に従い、EOF で未完了のレコードは破棄します。個々の行と蓄積したレコードは約 1 MiB に制限しています。サービスが他のプロトコル機能を必要とする場合は、完全な SSE クライアントを使ってください。サービス固有の JSON や終了マーカーは、テキストに追加する前に解釈する必要があります。

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

アプリケーションの作成後にサービスを登録します。

```go title="main.go"
app := application.New(application.Options{
    Name: "Streaming Example",
})

app.RegisterService(application.NewService(NewStreamService(app)))
```

`app.Context()` からリクエストのコンテキストを作ることで、アプリケーション終了時にもリクエストが停止します。`StopStream` は対応するリクエストだけをキャンセルするため、遅れて実行された後始末が新しいリクエストを止めることはありません。フロントエンドは毎回新しいストリーム ID を渡し、置き換え済みまたは停止済みのリクエストから配送中のイベントを無視します。このサービスが管理するリクエストはアプリケーション全体で一つです。ウィンドウごとに独立した配信が必要なら、別々のネイティブストリーム接続を使ってください。

## Svelte コンポーネント

受信した断片を通常の配列に保存します。同時に予約するアニメーションフレームを一つまでにし、そのフレームで待機中の断片を結合して Svelte の状態を一度更新します。コンポーネントは開始と停止の呼び出しを直列化し、バインディングの失敗を処理します。この例では Svelte 5 のルーンを使います。

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

`changeme` をアプリケーションの Go モジュール名に置き換え、サービス登録後にバインディングを生成してください。正確なパスは `StreamService` を含むパッケージによって異なります。例の URL を SSE エンドポイントに置き換え、実際の応答形式でテストしてください。

リスナーが受け取るのは `WailsEvent` なので、Go のペイロードは `event.data` にあります。`Events.On` が返す購読解除関数を保持し、コンポーネント破棄時に呼び出してください。後始末ではアニメーションフレームと対応するバックエンドのリクエストもキャンセルします。開始処理が遅れて完了した場合は、もう一度キャンセルします。

## フレーム単位でまとめる効果

受信中、コンポーネントが `output` に代入する回数はアニメーションフレームごとに最大一回です。停止時には待機中のテキストを明示的に反映します。60 Hz では、フレームの間に届いた断片を一つのリアクティブ更新にまとめられます。描画回数は到着のタイミング、画面のリフレッシュレート、フレームワークに依存し、応答ごとの固定回数にはなりません。

非表示のウィンドウではアニメーションフレームが停止する場合があります。この小さな例では、待機配列と蓄積された出力に上限がありません。長時間続く配信では、保持する履歴と待機データに制限を設け、生成側を一時停止するか古い表示データを捨てるかを決めてください。ネイティブ Streams はトランスポートのバックプレッシャーを提供しますが、どのトランスポートもフロントエンド受信後のアプリケーション状態の容量は制限しません。
