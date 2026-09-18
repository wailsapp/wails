---
title: "Streaming-Ereignisse"
description: "SSE in Go lesen und Svelte-Aktualisierungen bündeln, ohne das Rendering zu überlasten"
slug: "features/events/streaming"
sourcePath: "features/events/streaming.md"
---

Lang laufende Datenströme liefern oft viele kleine Aktualisierungen. Wird der reaktive Zustand für jedes Fragment geändert, kann das Rendering überlastet werden. Dieses Rezept liest eine externe Server-Sent-Events-Antwort (SSE) in Go und bündelt Zustandsänderungen in Svelte 5 mit `requestAnimationFrame`.

## Transport auswählen

Für einen neuen kontinuierlichen Punkt-zu-Punkt-Datenstrom empfehlen sich [native Streams](/guides/streams/): Registrieren Sie `app.HandleStream`, leiten Sie die vorgelagerte Anfrage von `c.Context()` ab, prüfen Sie Fehler von `c.SendJSON` und empfangen Sie Objekte mit `JSONStream` aus `@wailsio/runtime`. Das Schließen der Verbindung bricht ihren Kontext ab. `Stream` überträgt standardmäßig Bytes; es parst weder SSE noch verbindet es sich mit einer externen URL. Belassen Sie den vorgelagerten HTTP-Client in Go. Die Transportunterschiede erklärt [Migration eines WebSockets zu Streams](/guides/streams-from-websockets/).

Das folgende ereignisbasierte Beispiel eignet sich für die Integration in eine bestehende ereignisgesteuerte Anwendung. Es setzt eine zuständige Komponente in einem Fenster und eine aktive Anfrage des Dienstes voraus. Anwendungsereignisse werden an Listener verteilt; sie sind keine privaten Datenströme pro Verbindung. Die Bündelung pro Bild reduziert die Oberflächenarbeit, bietet aber keine Flusskontrolle im Transport. Für dauerhaft hohe Datenmengen oder unabhängige Empfänger verwenden Sie native Streams und behalten die Bündelung im Frontend bei.

## Backend-Dienst

Der Dienst liest mehrzeilige `data`-Felder einer UTF-8-SSE-Antwort mit LF- oder CRLF-Zeilenenden. Er startet eine Goroutine, damit `StartStream` sofort zurückkehrt, und bricht beim Start einer neuen Anfrage die vorherige ab. Verwenden Sie einen vertrauenswürdigen, von Ihrer Anwendung konfigurierten Endpunkt; die folgende URL ist ein Platzhalter.

Dies ist ein begrenzter Leser für Datenfelder, keine vollständige `EventSource`-Implementierung: Er ignoriert `event`, `id` und `retry`, verbindet sich nicht erneut und benötigt eine Leerzeile zur Ausgabe jedes Datensatzes. Ein unvollständiger Datensatz am Dateiende wird gemäß der [SSE-Spezifikation](https://html.spec.whatwg.org/multipage/server-sent-events.html#event-stream-interpretation) verworfen. Einzelne Zeilen und angesammelte Datensätze sind auf etwa 1 MiB begrenzt; benötigt Ihr Anbieter weitere Protokollfunktionen, verwenden Sie einen vollständigen SSE-Client. Anbieterspezifisches JSON oder Abschlussmarker müssen vor dem Anhängen von Text interpretiert werden.

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

Registrieren Sie den Dienst nach dem Erstellen der Anwendung:

```go title="main.go"
app := application.New(application.Options{
    Name: "Streaming Example",
})

app.RegisterService(application.NewService(NewStreamService(app)))
```

Durch Ableiten des Anfragekontexts von `app.Context()` beendet auch das Herunterfahren der Anwendung die Anfrage. `StopStream` bricht nur die passende Anfrage ab, sodass verspätetes Aufräumen keinen Nachfolger stoppt. Das Frontend liefert eine neue Stream-ID und ignoriert noch unterwegs befindliche Ereignisse ersetzter oder gestoppter Anfragen. Dieser Dienst verwaltet eine Anfrage für die gesamte Anwendung; benötigt jedes Fenster einen eigenen Datenstrom, verwenden Sie getrennte native Stream-Verbindungen.

## Svelte-Komponente

Speichern Sie ankommende Fragmente in einem gewöhnlichen Array. Planen Sie höchstens einen Animationsframe, verbinden Sie dann die wartenden Fragmente und aktualisieren Sie den Svelte-Zustand einmal pro Frame. Die Komponente serialisiert Start- und Stoppaufrufe und behandelt abgelehnte Binding-Aufrufe. Das Beispiel verwendet Svelte-5-Runes.

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

Ersetzen Sie `changeme` durch den Go-Modulnamen Ihrer Anwendung und generieren Sie die Bindings nach der Registrierung des Dienstes. Der genaue Pfad hängt vom Paket mit `StreamService` ab. Ersetzen Sie die Beispiel-URL durch Ihren SSE-Endpunkt und testen Sie dessen tatsächliches Antwortformat.

Der Listener erhält ein `WailsEvent`; die Go-Nutzdaten stehen daher in `event.data`. Bewahren Sie die von `Events.On` zurückgegebene Abmeldefunktion auf und rufen Sie sie beim Zerstören der Komponente auf. Das Aufräumen bricht auch den Animationsframe und die passende Backend-Anfrage ab; ein verspätet abgeschlossener Start löst die Abbruchanforderung erneut aus.

## Warum die Bündelung pro Frame hilft

Während des Empfangs führt die Komponente höchstens eine Zuweisung an `output` pro Animationsframe aus; beim Stoppen wird wartender Text ausdrücklich ausgegeben. Bei 60 Hz können zwischen zwei Frames empfangene Fragmente zu einer reaktiven Aktualisierung werden. Die Anzahl der Renderdurchläufe hängt von Ankunftszeit, Bildwiederholrate und Framework ab; sie ist nicht pro Antwort festgelegt.

Animationsframes können in einem verborgenen Fenster pausieren. Das Warte-Array und die angesammelte Ausgabe dieses kleinen Beispiels sind nicht begrenzt. Begrenzen Sie bei lang laufenden Datenströmen Verlauf und wartende Daten und entscheiden Sie, ob der Erzeuger pausieren oder alte Anzeigedaten verworfen werden sollen. Native Streams bieten Flusskontrolle im Transport; kein Transport begrenzt jedoch den Anwendungszustand nach dem Empfang im Frontend.
