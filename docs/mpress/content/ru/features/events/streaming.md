---
title: "Потоковые события"
description: "Чтение SSE в Go и группировка обновлений Svelte без перегрузки отрисовки"
slug: "features/events/streaming"
sourcePath: "features/events/streaming.md"
---

Длительные потоки часто доставляют множество небольших обновлений. Изменение реактивного состояния для каждого фрагмента может перегрузить отрисовку. В этом рецепте Go читает ответ внешнего сервиса с серверными событиями (SSE), а обновления состояния Svelte 5 группируются с помощью `requestAnimationFrame`.

## Выбор транспорта

Для нового непрерывного потока между двумя сторонами предпочтительны [встроенные Streams](/guides/streams/): зарегистрируйте `app.HandleStream`, создавайте запрос к внешнему сервису на основе `c.Context()`, проверяйте ошибки `c.SendJSON` и принимайте объекты через `JSONStream` из `@wailsio/runtime`. Закрытие соединения отменяет его контекст. По умолчанию `Stream` передаёт байты; он не разбирает SSE и не подключается к внешнему URL. Оставьте HTTP-клиент внешнего сервиса в Go. Различия транспорта описаны в руководстве [Переход с WebSocket на Streams](/guides/streams-from-websockets/).

Пример на событиях ниже полезен при интеграции с существующим событийным приложением. Он предполагает один управляющий компонент в одном окне и один активный запрос сервиса. События приложения рассылаются слушателям; это не закрытые потоки отдельных соединений. Группировка по кадрам снижает нагрузку на интерфейс, но не обеспечивает обратное давление транспорта. Для постоянно большого объёма данных или независимых потребителей используйте встроенные Streams, сохранив группировку на стороне интерфейса.

## Сервис бэкенда

Сервис читает многострочные поля `data` из ответа SSE в UTF-8 с окончаниями строк LF или CRLF. Он запускает горутину, чтобы `StartStream` возвращался сразу, и отменяет предыдущий запрос при запуске нового. Используйте доверенную конечную точку из настроек приложения; URL ниже служит заполнителем.

Это ограниченный читатель полей данных, а не полная реализация `EventSource`: он игнорирует `event`, `id` и `retry`, не переподключается и требует пустую строку для передачи каждой записи. Незавершённая запись в конце файла отбрасывается согласно [спецификации SSE](https://html.spec.whatwg.org/multipage/server-sent-events.html#event-stream-interpretation). Отдельные строки и накопленные записи ограничены примерно 1 MiB; если поставщику нужны другие возможности протокола, используйте полноценный SSE-клиент. Специфичные для поставщика JSON или маркеры завершения нужно обработать до добавления текста.

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

Зарегистрируйте сервис после создания приложения:

```go title="main.go"
app := application.New(application.Options{
    Name: "Streaming Example",
})

app.RegisterService(application.NewService(NewStreamService(app)))
```

Контекст запроса, полученный из `app.Context()`, обеспечивает остановку запроса при завершении приложения. `StopStream` отменяет только соответствующий запрос, поэтому запоздалая очистка не остановит заменивший его запрос. Интерфейс передаёт новый идентификатор потока и игнорирует ещё находящиеся в пути события заменённых или остановленных запросов. Сервис управляет одним запросом на всё приложение; если каждому окну нужен собственный поток, используйте отдельные встроенные соединения.

## Компонент Svelte

Храните входящие фрагменты в обычном массиве. Планируйте не более одного кадра анимации, затем объединяйте ожидающие фрагменты и обновляйте состояние Svelte один раз в этом кадре. Компонент выполняет вызовы запуска и остановки последовательно и обрабатывает отклонённые вызовы привязок. В примере используются руны Svelte 5.

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

Замените `changeme` именем модуля Go вашего приложения и сгенерируйте привязки после регистрации сервиса. Точный путь зависит от пакета, содержащего `StreamService`. Замените пример URL своей конечной точкой SSE и проверьте фактический формат её ответа.

Слушатель получает `WailsEvent`, поэтому полезная нагрузка Go доступна в `event.data`. Сохраните функцию отписки, возвращаемую `Events.On`, и вызывайте её при уничтожении компонента. Очистка также отменяет кадр анимации и соответствующий запрос бэкенда; позднее завершение запуска повторно инициирует отмену.

## Почему группировка по кадрам помогает

Во время приёма компонент присваивает `output` значение не чаще одного раза за кадр анимации; остановка явно выводит ожидающий текст. При 60 Гц фрагменты, полученные между кадрами, могут объединиться в одно реактивное обновление. Количество отрисовок зависит от времени поступления, частоты обновления экрана и фреймворка; фиксированного числа на ответ нет.

В скрытом окне кадры анимации могут приостанавливаться. Массив ожидания и накопленный вывод в этом небольшом примере не ограничены. Для длительных потоков ограничьте историю и ожидающие данные и решите, приостанавливать ли источник или отбрасывать старые данные отображения. Встроенные Streams обеспечивают обратное давление транспорта, но ни один транспорт не ограничивает состояние приложения после получения данных интерфейсом.
