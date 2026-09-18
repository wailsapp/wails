---
title: "Peristiwa streaming"
description: "Baca SSE di Go dan kelompokkan pembaruan Svelte tanpa membebani perenderan"
slug: "features/events/streaming"
sourcePath: "features/events/streaming.md"
---

Aliran yang berjalan lama sering mengirim banyak pembaruan kecil. Memperbarui keadaan reaktif untuk setiap potongan dapat membebani perenderan. Resep ini membaca respons server-sent events (SSE) eksternal di Go dan mengelompokkan pembaruan keadaan Svelte 5 dengan `requestAnimationFrame`.

## Pilih transport

Untuk aliran berkelanjutan antardua pihak yang baru, utamakan [Streams bawaan](/guides/streams/): daftarkan `app.HandleStream`, turunkan konteks permintaan ke layanan sumber dari `c.Context()`, periksa galat dari `c.SendJSON`, dan terima objek dengan `JSONStream` dari `@wailsio/runtime`. Menutup koneksi membatalkan konteksnya. Secara bawaan, `Stream` mengirim byte; ia tidak mengurai SSE atau terhubung ke URL eksternal. Tempatkan klien HTTP layanan sumber di Go. Lihat [Migrasi WebSocket ke Streams](/guides/streams-from-websockets/) untuk perbedaan transport.

Contoh berbasis peristiwa berikut berguna untuk integrasi dengan aplikasi berbasis peristiwa yang sudah ada. Contoh ini mengasumsikan satu komponen pemilik dalam satu jendela dan satu permintaan aktif pada layanan. Peristiwa aplikasi disiarkan ke pendengar; bukan aliran privat per koneksi. Pengelompokan per bingkai mengurangi pekerjaan antarmuka, tetapi tidak menambahkan tekanan balik pada transport. Untuk volume tinggi yang berkelanjutan atau konsumen independen, gunakan Streams bawaan sambil mempertahankan teknik pengelompokan frontend yang sama.

## Layanan backend

Layanan membaca kolom `data` multibaris dari respons SSE UTF-8 dengan akhir baris LF atau CRLF. Layanan memulai goroutine agar `StartStream` segera kembali dan membatalkan permintaan sebelumnya saat permintaan baru dimulai. Gunakan endpoint tepercaya yang dikonfigurasi aplikasi; URL di bawah hanyalah contoh pengganti.

Ini adalah pembaca khusus data dengan batas ukuran, bukan implementasi `EventSource` lengkap: ia mengabaikan `event`, `id`, dan `retry`, tidak menyambung ulang, serta membutuhkan baris kosong untuk mengirim setiap rekaman. Rekaman yang belum selesai saat EOF dibuang sesuai [spesifikasi SSE](https://html.spec.whatwg.org/multipage/server-sent-events.html#event-stream-interpretation). Baris individual dan rekaman yang terkumpul dibatasi sekitar 1 MiB; gunakan klien SSE lengkap jika penyedia membutuhkan fitur protokol lain. JSON atau penanda selesai khusus penyedia harus ditafsirkan sebelum teks ditambahkan.

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

Daftarkan layanan setelah membuat aplikasi:

```go title="main.go"
app := application.New(application.Options{
    Name: "Streaming Example",
})

app.RegisterService(application.NewService(NewStreamService(app)))
```

Menurunkan konteks permintaan dari `app.Context()` memastikan penutupan aplikasi juga menghentikan permintaan. `StopStream` hanya membatalkan permintaan yang cocok, sehingga pembersihan terlambat tidak menghentikan penggantinya. Frontend memberikan ID aliran baru dan mengabaikan peristiwa yang masih dalam perjalanan dari permintaan yang diganti atau dihentikan. Layanan ini memiliki satu permintaan untuk seluruh aplikasi; gunakan koneksi aliran bawaan terpisah jika setiap jendela membutuhkan alirannya sendiri.

## Komponen Svelte

Simpan potongan masuk dalam array biasa. Jadwalkan paling banyak satu bingkai animasi, lalu gabungkan potongan yang menunggu dan perbarui keadaan Svelte sekali pada bingkai itu. Komponen mengurutkan pemanggilan mulai dan berhenti serta menangani penolakan binding. Contoh menggunakan rune Svelte 5.

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

Ganti `changeme` dengan nama modul Go aplikasi dan buat binding setelah mendaftarkan layanan. Jalur binding yang tepat bergantung pada paket yang berisi `StreamService`. Ganti URL contoh dengan endpoint SSE Anda dan uji menggunakan format respons sebenarnya.

Pendengar menerima `WailsEvent`, sehingga muatan Go tersedia pada `event.data`. Simpan fungsi berhenti berlangganan yang dikembalikan `Events.On` dan panggil saat komponen dihancurkan. Pembersihan juga membatalkan bingkai animasi serta permintaan backend yang cocok; penyelesaian mulai yang terlambat memicu pembatalan lagi.

## Mengapa pengelompokan per bingkai membantu

Selama menerima aliran, komponen melakukan paling banyak satu penetapan `output` per bingkai animasi; penghentian secara eksplisit menampilkan teks tertunda. Pada 60 Hz, potongan yang diterima di antara bingkai dapat menjadi satu pembaruan reaktif. Jumlah perenderan bergantung pada waktu kedatangan, laju penyegaran layar, dan framework; bukan jumlah tetap per respons.

Bingkai animasi dapat berhenti sementara pada jendela tersembunyi. Array antrean dan keluaran yang terkumpul dalam contoh kecil ini tidak dibatasi. Untuk aliran berumur panjang, batasi riwayat yang disimpan dan data tertunda, lalu tentukan apakah produsen perlu dijeda atau data tampilan lama dibuang. Streams bawaan menyediakan tekanan balik transport, tetapi transport mana pun tidak membatasi keadaan aplikasi setelah frontend menerima data.
