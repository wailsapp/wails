---
title: "Pesan Mentah"
description: "Implementasikan komunikasi khusus dari frontend ke backend untuk aplikasi yang sangat mementingkan performa"
slug: "guides/raw-messages"
sourcePath: "guides/raw-messages.md"
---

Pesan mentah menyediakan kanal komunikasi tingkat rendah antara frontend dan backend Anda dengan melewati sistem binding standar. Pendekatan ini mengorbankan kemudahan demi kecepatan.

## Kapan Menggunakan Pesan Mentah

Pesan mentah paling sesuai untuk kasus ekstrem yang sangat jarang terjadi:

- **Pembaruan berfrekuensi sangat tinggi** - Ribuan pesan per detik ketika setiap mikrodetik sangat berarti
- **Protokol pesan khusus** - Ketika Anda memerlukan kendali penuh atas format data yang ditransmisikan

@note{type="tip"}
Untuk hampir semua kasus penggunaan, [binding layanan](/features/bindings/services/) standar direkomendasikan karena menyediakan keamanan tipe, serialisasi otomatis, dan pengalaman pengembangan yang lebih baik dengan overhead yang dapat diabaikan.

@end

## Penyiapan Backend

Konfigurasikan `RawMessageHandler` dalam opsi aplikasi Anda:

```go
package main

import (
    "encoding/json"
    "fmt"

    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "Raw Message Demo",
        Assets: application.AssetOptions{
            Handler: application.BundledAssetFileServer(assets),
        },
        RawMessageHandler: func(window application.Window, message string, originInfo *application.OriginInfo) {
            fmt.Printf("Raw message from window '%s': %s (origin: %+v)\n", window.Name(), message, originInfo.Origin)

            // Process the message and respond via events
            response := processMessage(message)
            window.EmitEvent("raw-response", response)
        },
    })

    app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title: "My App",
        Name:  "main",
    })

    app.Run()
}

func processMessage(message string) map[string]any {
    // Your custom message processing logic
    return map[string]any{
        "received": message,
        "status":   "processed",
    }
}
```

### Signature Handler

```go
RawMessageHandler func(window Window, message string, originInfo *application.OriginInfo)
```

| Parameter | Tipe | Deskripsi |
| --- | --- | --- |
| `window` | `Window` | Jendela yang mengirim pesan |
| `message` | `string` | Konten pesan mentah |
| `originInfo` | `*application.OriginInfo` | Informasi asal tentang sumber pesan |

#### Struktur OriginInfo

```go
type OriginInfo struct {
	Origin      string
	TopOrigin   string
	IsMainFrame bool
}
```

| Field | Tipe | Deskripsi |
| --- | --- | --- |
| `Origin` | `string` | URL asal dokumen yang mengirim pesan |
| `TopOrigin` | `string` | URL asal tingkat teratas (mungkin berbeda dari Origin dalam iframe) |
| `IsMainFrame` | `bool` | Apakah pesan berasal dari frame utama |

#### Ketersediaan Khusus Platform

- **macOS**: `Origin` dan `IsMainFrame` disediakan
- **Windows**: `Origin` dan `TopOrigin` disediakan
- **Linux**: Hanya `Origin` yang disediakan

### Validasi Asal

@note{type="caution"}
Jangan pernah menganggap pesan aman hanya karena pesan tersebut masuk ke handler Anda. Informasi asal harus divalidasi sebelum memproses operasi sensitif atau operasi yang mengubah status.

@end

**Selalu verifikasi asal pesan masuk sebelum memprosesnya.** Parameter `originInfo` menyediakan informasi keamanan penting yang harus divalidasi untuk mencegah akses tanpa izin. Konten berbahaya, konten yang telah disusupi, atau skrip yang tidak dimaksudkan dapat mengirim pesan mentah. Tanpa validasi asal, Anda mungkin memproses perintah dari sumber yang tidak tepercaya. Gunakan `originInfo` untuk memastikan pesan berasal dari sumber yang diharapkan.

### Poin Validasi Utama

- **Selalu periksa `Origin`** - Pastikan asal cocok dengan sumber tepercaya yang Anda harapkan (biasanya `wails://wails` atau `http://wails.localhost` untuk aset lokal atau asal khusus aplikasi Anda)
- **Validasi `IsMainFrame`** (macOS) - Perhatikan jika pesan berasal dari iframe karena hal ini dapat menunjukkan konten tersemat dengan konteks keamanan yang berbeda
- **Gunakan `TopOrigin`** (Windows) - Verifikasi asal tingkat teratas saat menangani konten dalam frame
- **Tolak asal yang tidak diharapkan** - Terapkan kegagalan yang aman dengan menolak pesan dari asal yang tidak Anda izinkan secara eksplisit

@note{type="info"}
Pesan yang diawali dengan `wails:` dicadangkan untuk komunikasi internal Wails dan tidak akan diteruskan ke handler Anda.

@end

## Penyiapan Frontend

Kirim pesan mentah menggunakan `System.invoke()`:

```html
<!DOCTYPE html>
<html>
<head>
    <script type="module">
        import { System, Events } from '@wailsio/runtime'

        // Send raw message
        document.getElementById('send').addEventListener('click', () => {
            const message = document.getElementById('input').value
            System.invoke(message)
        })

        // Listen for response
        Events.On('raw-response', (event) => {
            console.log('Response:', event.data)
        })
    </script>
</head>
<body>
    <input type="text" id="input" placeholder="Enter message" />
    <button id="send">Send</button>
</body>
</html>
```

### Menggunakan Bundle Siap Pakai

Jika Anda tidak menggunakan npm, akses `invoke` melalui objek global `wails`:

```html
<script type="module" src="/wails/runtime.js"></script>
<script>
    window.onload = function() {
        document.getElementById('send').onclick = function() {
            wails.System.invoke('my-message')
        }
    }
</script>
```

## Pesan Terstruktur

Untuk data kompleks, lakukan serialisasi ke JSON:

### Frontend

```javascript
import { System } from '@wailsio/runtime'

const command = {
    action: 'update',
    payload: {
        id: 123,
        value: 'new value'
    }
}

System.invoke(JSON.stringify(command))
```

### Backend

```go
RawMessageHandler: func(window application.Window, message string, originInfo *application.OriginInfo) {
    var cmd struct {
        Action  string `json:"action"`
        Payload struct {
            ID    int    `json:"id"`
            Value string `json:"value"`
        } `json:"payload"`
    }

    if err := json.Unmarshal([]byte(message), &cmd); err != nil {
        window.EmitEvent("error", err.Error())
        return
    }

    switch cmd.Action {
    case "update":
        // Handle update
        result := handleUpdate(cmd.Payload.ID, cmd.Payload.Value)
        window.EmitEvent("update-complete", result)
    default:
        window.EmitEvent("error", "unknown action")
    }
}
```

## Perbandingan Performa

| Pendekatan | Overhead | Keamanan Tipe | Kasus Penggunaan |
| --- | --- | --- | --- |
| Binding Layanan | Lebih Tinggi | Penuh | Tujuan umum |
| Pesan Mentah | Minimal | Manual | Frekuensi tinggi, kritis terhadap performa |

### Contoh Tolok Ukur

Untuk payload sederhana, pesan mentah dapat memproses jauh lebih banyak pesan per detik dibandingkan binding layanan:

```go
// Raw message handler - minimal overhead
RawMessageHandler: func(window application.Window, message string, originInfo *application.OriginInfo) {
    // Direct string processing, no reflection or marshaling
    counter++
}
```

## Contoh Lengkap

Berikut contoh lengkap yang mengimplementasikan protokol perintah sederhana:

### main.go

```go
package main

import (
    "embed"
    "encoding/json"
    "fmt"
    "time"

    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets
var assets embed.FS

type Command struct {
    Type string          `json:"type"`
    Data json.RawMessage `json:"data"`
}

func main() {
    app := application.New(application.Options{
        Name: "Raw Message Demo",
        Assets: application.AssetOptions{
            Handler: application.BundledAssetFileServer(assets),
        },
        Mac: application.MacOptions{
            ApplicationShouldTerminateAfterLastWindowClosed: true,
        },
        RawMessageHandler: func(window application.Window, message string, originInfo *application.OriginInfo) {
            var cmd Command
            if err := json.Unmarshal([]byte(message), &cmd); err != nil {
                window.EmitEvent("error", map[string]string{"error": err.Error()})
                return
            }

            switch cmd.Type {
            case "ping":
                window.EmitEvent("pong", map[string]any{
                    "time":   time.Now().UnixMilli(),
                    "window": window.Name(),
                })
            case "echo":
                var text string
                json.Unmarshal(cmd.Data, &text)
                window.EmitEvent("echo", text)
            default:
                window.EmitEvent("error", map[string]string{
                    "error": fmt.Sprintf("unknown command: %s", cmd.Type),
                })
            }
        },
    })

    app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title: "Raw Message Demo",
        Name:  "main",
        Width: 400,
        Height: 300,
    })

    app.Run()
}
```

### assets/index.html

```html
<!DOCTYPE html>
<html>
<head>
    <title>Raw Message Demo</title>
    <style>
        body { font-family: sans-serif; padding: 20px; }
        button { margin: 5px; padding: 10px 20px; }
        #output { margin-top: 20px; padding: 10px; background: #f0f0f0; }
    </style>
</head>
<body>
    <h1>Raw Message Demo</h1>

    <button id="ping">Ping</button>
    <button id="echo">Echo "Hello"</button>

    <div id="output">Waiting for response...</div>

    <script type="module">
        import { System, Events } from '@wailsio/runtime'

        const output = document.getElementById('output')

        function send(type, data) {
            System.invoke(JSON.stringify({ type, data }))
        }

        document.getElementById('ping').onclick = () => send('ping')
        document.getElementById('echo').onclick = () => send('echo', 'Hello')

        Events.On('pong', (e) => {
            output.textContent = `Pong from ${e.data.window} at ${e.data.time}`
        })

        Events.On('echo', (e) => {
            output.textContent = `Echo: ${e.data}`
        })

        Events.On('error', (e) => {
            output.textContent = `Error: ${e.data.error}`
        })
    </script>
</body>
</html>
```

## Praktik Terbaik

### Lakukan

- Gunakan pesan mentah untuk jalur yang benar-benar kritis terhadap performa
- Terapkan penanganan kesalahan yang tepat di handler Anda
- Gunakan event untuk mengirim respons kembali ke frontend
- Pertimbangkan JSON untuk data terstruktur
- Pastikan pemrosesan pesan tetap cepat agar tidak memblokir

### Jangan Lakukan

- Jangan gunakan pesan mentah jika binding layanan sudah memadai
- Jangan lupa memvalidasi pesan yang masuk
- Jangan memblokir handler dengan operasi yang berjalan lama (gunakan goroutine)
- Jangan abaikan parameter window ketika respons perlu ditujukan ke window tertentu

## Pertimbangan Multi-Window

Parameter `window` mengidentifikasi window yang mengirim pesan, sehingga Anda dapat:

- Mengirim respons ke window yang tepat
- Menerapkan perilaku khusus untuk setiap window
- Melacak sumber pesan untuk debugging

```go
RawMessageHandler: func(window application.Window, message string, originInfo *application.OriginInfo) {
    // Respond only to the sending window
    window.EmitEvent("response", result)

    // Or broadcast to all windows
    app.Event.Emit("broadcast", result)
}
```

## Langkah Berikutnya

- [Binding Layanan](/features/bindings/services/) - Pendekatan standar untuk sebagian besar aplikasi
- [Event](/guides/events-reference/) - Sistem event untuk komunikasi dari backend ke frontend
- [Performa](/guides/performance/) - Pengoptimalan performa secara umum
