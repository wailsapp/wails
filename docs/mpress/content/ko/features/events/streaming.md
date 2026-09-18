---
title: "스트리밍 이벤트"
description: "Go에서 SSE를 읽고 렌더링에 과부하를 주지 않도록 Svelte 업데이트 묶기"
slug: "features/events/streaming"
sourcePath: "features/events/streaming.md"
---

오래 실행되는 스트림은 작은 업데이트를 많이 전달하는 경우가 많습니다. 조각마다 반응형 상태를 바꾸면 렌더링에 과부하가 걸릴 수 있습니다. 이 예제는 Go에서 외부 서버 전송 이벤트(SSE) 응답을 읽고 `requestAnimationFrame`으로 Svelte 5 상태 업데이트를 묶습니다.

## 전송 방식 선택

새로운 연속적인 일대일 데이터 전송에는 [네이티브 Streams](/guides/streams/)를 권장합니다. `app.HandleStream`을 등록하고, `c.Context()`에서 외부 요청의 컨텍스트를 만들고, `c.SendJSON`의 오류를 확인하며, `@wailsio/runtime`의 `JSONStream`으로 객체를 받으세요. 연결을 닫으면 해당 컨텍스트가 취소됩니다. `Stream`은 기본적으로 바이트를 전달하며 SSE를 파싱하거나 외부 URL에 연결하지 않습니다. 외부 서비스를 위한 HTTP 클라이언트는 Go에 두세요. 전송 방식의 차이는 [WebSocket에서 Streams로 이전하기](/guides/streams-from-websockets/)를 참조하세요.

아래 이벤트 기반 예제는 기존 이벤트 중심 애플리케이션에 통합할 때 유용합니다. 한 창에 담당 컴포넌트가 하나 있고 서비스에 활성 요청이 하나 있다고 가정합니다. 애플리케이션 이벤트는 리스너들에게 브로드캐스트되며 연결별 비공개 스트림이 아닙니다. 프레임별 묶음 처리는 UI 작업을 줄이지만 전송 계층에 역압을 추가하지는 않습니다. 지속적인 대용량 데이터나 독립적인 소비자가 필요하다면 네이티브 Streams를 사용하면서 동일한 프런트엔드 묶음 처리 기법을 유지하세요.

## 백엔드 서비스

서비스는 LF 또는 CRLF 줄바꿈을 사용하는 UTF-8 SSE 응답에서 여러 줄의 `data` 필드를 읽습니다. goroutine을 시작하므로 `StartStream`이 즉시 반환되고, 새 요청이 시작되면 이전 요청을 취소합니다. 애플리케이션에서 설정한 신뢰할 수 있는 엔드포인트를 사용하세요. 아래 URL은 바꿔 넣을 예시입니다.

이 코드는 크기를 제한한 데이터 전용 리더이며 완전한 `EventSource` 구현이 아닙니다. `event`, `id`, `retry`를 무시하고 재연결하지 않으며 각 레코드를 전달하려면 빈 줄이 필요합니다. [SSE 명세](https://html.spec.whatwg.org/multipage/server-sent-events.html#event-stream-interpretation)에 따라 EOF에서 미완성 레코드는 버립니다. 개별 줄과 누적 레코드는 약 1 MiB로 제한됩니다. 제공자가 다른 프로토콜 기능을 요구하면 완전한 SSE 클라이언트를 사용하세요. 제공자별 JSON이나 완료 표시는 텍스트에 추가하기 전에 해석해야 합니다.

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

애플리케이션을 생성한 뒤 서비스를 등록하세요.

```go title="main.go"
app := application.New(application.Options{
    Name: "Streaming Example",
})

app.RegisterService(application.NewService(NewStreamService(app)))
```

`app.Context()`에서 요청 컨텍스트를 만들면 애플리케이션 종료 시 요청도 중단됩니다. `StopStream`은 일치하는 요청만 취소하므로 늦은 정리가 대체 요청을 중단하지 않습니다. 프런트엔드는 매번 새 스트림 ID를 전달하고 교체되거나 중단된 요청에서 아직 전달 중인 이벤트는 무시합니다. 이 서비스는 애플리케이션 전체에서 요청 하나를 관리합니다. 창마다 자체 데이터 전송이 필요하면 별도의 네이티브 스트림 연결을 사용하세요.

## Svelte 컴포넌트

수신한 조각을 일반 배열에 보관하세요. 한 번에 최대 하나의 애니메이션 프레임을 예약하고, 해당 프레임에서 대기 중인 조각을 합쳐 Svelte 상태를 한 번 업데이트합니다. 컴포넌트는 시작 및 중지 호출을 순서대로 실행하고 거부된 바인딩 호출을 처리합니다. 예제는 Svelte 5 룬을 사용합니다.

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

`changeme`를 애플리케이션의 Go 모듈 이름으로 바꾸고 서비스를 등록한 후 바인딩을 생성하세요. 정확한 경로는 `StreamService`가 포함된 패키지에 따라 달라집니다. 예제 URL을 실제 SSE 엔드포인트로 바꾸고 실제 응답 형식으로 테스트하세요.

리스너는 `WailsEvent`를 받으므로 Go 페이로드는 `event.data`에 있습니다. `Events.On`이 반환하는 구독 해제 함수를 보관하고 컴포넌트가 파괴될 때 호출하세요. 정리 과정은 애니메이션 프레임과 일치하는 백엔드 요청도 취소합니다. 시작 작업이 늦게 완료되면 취소를 다시 실행합니다.

## 프레임별 묶음 처리가 도움이 되는 이유

수신 중에는 애니메이션 프레임마다 `output`에 최대 한 번만 대입하고, 중지할 때는 대기 중인 텍스트를 명시적으로 반영합니다. 60 Hz에서는 프레임 사이에 받은 조각들을 하나의 반응형 업데이트로 묶을 수 있습니다. 렌더링 횟수는 도착 시점, 화면 새로고침 빈도, 프레임워크에 따라 달라지며 응답당 고정된 수가 아닙니다.

숨겨진 창에서는 애니메이션 프레임이 일시 중지될 수 있습니다. 이 작은 예제의 대기 배열과 누적 출력에는 제한이 없습니다. 오래 지속되는 데이터 전송에서는 보관 기록과 대기 데이터의 크기를 제한하고, 생성자를 멈출지 오래된 표시 데이터를 버릴지 결정하세요. 네이티브 Streams는 전송 역압을 제공하지만 어떤 전송 방식도 프런트엔드가 데이터를 받은 후의 애플리케이션 상태 크기를 제한하지는 않습니다.
