---
title: "Événements en continu"
description: "Lire SSE en Go et regrouper les mises à jour Svelte sans surcharger le rendu"
slug: "features/events/streaming"
sourcePath: "features/events/streaming.md"
---

Les flux de longue durée produisent souvent de nombreuses petites mises à jour. Modifier l’état réactif à chaque fragment peut surcharger le rendu. Cette recette lit en Go une réponse externe d’événements envoyés par le serveur (SSE) et regroupe les mises à jour d’état de Svelte 5 avec `requestAnimationFrame`.

## Choisir le transport

Pour un nouveau flux continu point à point, privilégiez les [Streams natifs](/guides/streams/) : enregistrez `app.HandleStream`, dérivez la requête amont de `c.Context()`, vérifiez les erreurs de `c.SendJSON` et recevez les objets avec `JSONStream` de `@wailsio/runtime`. Fermer cette connexion annule son contexte. `Stream` transmet des octets par défaut ; il n’analyse pas SSE et ne se connecte pas à une URL externe. Gardez le client HTTP amont en Go. Consultez [Migrer un WebSocket vers Streams](/guides/streams-from-websockets/) pour les différences de transport.

L’exemple ci-dessous convient à l’intégration dans une application existante fondée sur les événements. Il suppose un seul composant propriétaire dans une seule fenêtre et une seule requête active pour le service. Les événements de l’application sont diffusés aux écouteurs ; ce ne sont pas des flux privés par connexion. Le regroupement par image réduit le travail de l’interface, sans réguler le débit du transport. Pour un volume élevé et soutenu ou des consommateurs indépendants, utilisez les Streams natifs en conservant ce regroupement côté interface.

## Service backend

Le service lit les champs `data` multilignes d’une réponse SSE UTF-8 avec des fins de ligne LF ou CRLF. Il lance une goroutine pour que `StartStream` revienne immédiatement et annule la requête précédente au démarrage d’une nouvelle. Utilisez un point de terminaison de confiance configuré par votre application ; l’URL ci-dessous est un exemple à remplacer.

Ce lecteur limité aux données et à mémoire bornée n’est pas une implémentation complète d’`EventSource` : il ignore `event`, `id` et `retry`, ne se reconnecte pas et exige une ligne vide pour transmettre chaque enregistrement. Un enregistrement incomplet en fin de fichier est abandonné, conformément à la [spécification SSE](https://html.spec.whatwg.org/multipage/server-sent-events.html#event-stream-interpretation). Les lignes et les enregistrements accumulés sont limités à environ 1 MiB ; utilisez un client SSE complet si votre fournisseur exige d’autres fonctions du protocole. Interprétez le JSON ou les marqueurs de fin propres au fournisseur avant d’ajouter du texte.

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

Enregistrez le service après avoir créé l’application :

```go title="main.go"
app := application.New(application.Options{
    Name: "Streaming Example",
})

app.RegisterService(application.NewService(NewStreamService(app)))
```

Dériver le contexte de la requête d’`app.Context()` garantit que l’arrêt de l’application interrompt aussi la requête. `StopStream` n’annule que la requête correspondante ; un nettoyage tardif ne peut donc pas arrêter sa remplaçante. L’interface fournit un nouvel identifiant de flux et ignore les événements encore en transit des requêtes remplacées ou arrêtées. Ce service possède une seule requête pour toute l’application ; utilisez des connexions Streams natives distinctes si chaque fenêtre nécessite son propre flux.

## Composant Svelte

Conservez les fragments reçus dans un tableau ordinaire. Planifiez au maximum une image d’animation, puis joignez les fragments en attente et mettez à jour l’état Svelte une seule fois dans cette image. Le composant sérialise ses appels de démarrage et d’arrêt et gère les rejets des bindings. L’exemple utilise les runes de Svelte 5.

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

Remplacez `changeme` par le nom du module Go de votre application et générez les bindings après l’enregistrement du service. Leur chemin exact dépend du paquet contenant `StreamService`. Remplacez l’URL d’exemple par votre point de terminaison SSE et testez son format de réponse réel.

L’écouteur reçoit un `WailsEvent` ; la charge utile Go est donc disponible dans `event.data`. Conservez la fonction de désabonnement renvoyée par `Events.On` et appelez-la à la destruction du composant. Le nettoyage annule aussi l’image d’animation et la requête backend correspondante ; si le démarrage se termine tardivement, l’annulation est répétée.

## Pourquoi regrouper par image

Pendant la réception, le composant effectue au maximum une affectation à `output` par image d’animation ; l’arrêt vide explicitement le texte en attente. À 60 Hz, les fragments reçus entre deux images peuvent former une seule mise à jour réactive. Le nombre de rendus dépend de l’arrivée des données, du taux de rafraîchissement et du framework ; il n’est pas fixe par réponse.

Les images d’animation peuvent être suspendues dans une fenêtre masquée. Le tableau d’attente et la sortie accumulée de ce petit exemple ne sont pas bornés. Pour un flux durable, limitez l’historique et les données en attente, puis choisissez de suspendre le producteur ou d’abandonner les anciennes données affichées. Les Streams natifs régulent le débit du transport, mais aucun transport ne borne l’état applicatif après réception par l’interface.
