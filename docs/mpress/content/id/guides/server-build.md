---
title: "Build Server"
description: "Jalankan aplikasi Wails sebagai server HTTP tanpa jendela GUI native"
slug: "guides/server-build"
sourcePath: "guides/server-build.md"
---

Wails v3 mendukung mode server, yang memungkinkan Anda menjalankan aplikasi sebagai server HTTP murni tanpa membuat jendela native atau memerlukan dependensi GUI. Dengan demikian, aplikasi Wails yang sama dapat diterapkan ke server dan kontainer, serta diakses melalui browser web.

Mode server berguna untuk:

- **Penerapan Docker/kontainer** - Berjalan tanpa dependensi X11/Wayland
- **Aplikasi sisi server** - Terapkan sebagai server web yang dapat diakses melalui browser
- **Akses khusus web** - Gunakan basis kode yang sama untuk desktop dan web
- **Pengujian CI/CD** - Jalankan pengujian integrasi tanpa server tampilan
- **Layanan mikro** - Gunakan binding Wails dalam layanan backend tanpa antarmuka grafis

## Mulai Cepat

Mode server diaktifkan melalui tag build `server`. Kode aplikasi Anda tetap sama—Anda hanya perlu melakukan build dengan tag tersebut:

```bash
# Using Taskfile (recommended)
wails3 task build:server
wails3 task run:server

# Or build directly with Go
go build -tags server -o myapp-server .
```

Berikut contoh minimalnya:

```go
package main

import (
    "embed"
    "log"

    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
    app := application.New(application.Options{
        Name: "My App",
        // Server options are used when built with -tags server
        Server: application.ServerOptions{
            Host: "localhost",
            Port: 8080,
        },
        Services: []application.Service{
            application.NewService(&MyService{}),
        },
        Assets: application.AssetOptions{
            Handler: application.AssetFileServerFS(assets),
        },
    })

    log.Println("Starting application...")
    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}
```

Kode yang sama dapat di-build untuk mode desktop (tanpa tag) atau mode server (dengan `-tags server`).

## Konfigurasi

### ServerOptions

Konfigurasikan server HTTP dengan `ServerOptions`:

```go
Server: application.ServerOptions{
    // Host to bind to. Default: "localhost"
    // Use "0.0.0.0" to listen on all interfaces
    Host: "localhost",

    // Port to listen on. Default: 8080
    Port: 8080,

    // Request read timeout. Default: 30s
    ReadTimeout: 30 * time.Second,

    // Response write timeout. Default: 30s
    WriteTimeout: 30 * time.Second,

    // Idle connection timeout. Default: 120s
    IdleTimeout: 120 * time.Second,

    // Graceful shutdown timeout. Default: 30s
    ShutdownTimeout: 30 * time.Second,

    // Additional origins allowed to open WebSocket connections.
    // Same-origin connections are always allowed.
    WebSocketOriginPatterns: []string{"app.example.com"},

    // Disable WebSocket origin checks. Unsafe; default: false.
    WebSocketAllowAllOrigins: false,

    // TLS configuration (optional)
    TLS: &application.TLSOptions{
        CertFile: "/path/to/cert.pem",
        KeyFile:  "/path/to/key.pem",
    },
},
```

## Fitur

### Endpoint Pemeriksaan Kesehatan

Endpoint pemeriksaan kesehatan tersedia secara otomatis di `/health`:

```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

Endpoint ini berguna untuk:

- Probe liveness/readiness Kubernetes
- Pemeriksaan kesehatan load balancer
- Sistem pemantauan

### Binding Layanan

Semua binding layanan berfungsi sama seperti dalam mode desktop:

```go
type GreetService struct{}

func (g *GreetService) Greet(name string) string {
    return "Hello, " + name + "!"
}

// Register in options
Services: []application.Service{
    application.NewService(&GreetService{}),
},
```

Frontend dapat memanggil binding ini menggunakan runtime Wails standar:

```javascript
const greeting = await wails.Call.ByName('main.GreetService.Greet', 'World');
```

### Peristiwa

Peristiwa berfungsi dua arah dalam mode server:

- **Frontend ke Backend**: Peristiwa yang dipancarkan dari browser dikirim melalui HTTP dan diterima oleh handler peristiwa Go Anda
- **Backend ke Frontend**: Peristiwa yang dipancarkan dari Go disiarkan ke semua browser yang terhubung melalui WebSocket

Setiap tab browser direpresentasikan sebagai "jendela" dengan nama unik (`browser-1`, `browser-2`, dan seterusnya), yang dapat diakses melalui `event.Sender`:

```go
// Listen for events from browsers
app.Event.On("user-action", func(event *application.CustomEvent) {
    log.Printf("Event from %s: %v", event.Sender, event.Data)
    // event.Sender will be "browser-1", "browser-2", etc.
})

// Emit events to all connected browsers
app.Event.Emit("server-update", data)
```

Dari frontend:

```javascript
// Emit event to server (and all other browsers)
await wails.Events.Emit('user-action', { action: 'click' });

// Listen for events from server
wails.Events.On('server-update', (event) => {
    console.log('Update from server:', event.data);
});
```

### Penghentian yang Mulus

Server menangani sinyal `SIGINT` dan `SIGTERM` dengan baik:

1. Berhenti menerima koneksi baru
2. Menunggu permintaan aktif selesai (hingga `ShutdownTimeout`)
3. Menjalankan hook `OnShutdown`
4. Menghentikan layanan dalam urutan terbalik

## Perbedaan dari Mode Desktop

| Fitur | Mode Desktop | Mode Server |
| --- | --- | --- |
| Jendela native | Dibuat | Jendela browser (`browser-N`) |
| Baki sistem | Tersedia | Tidak tersedia |
| Dialog native | Tersedia | Tidak tersedia |
| Menu aplikasi | Tersedia | Tidak tersedia |
| Informasi layar | Tersedia | Mengembalikan galat |
| Binding layanan | Berfungsi | Berfungsi |
| Peristiwa | Berfungsi | Berfungsi (melalui WebSocket) |
| Aset | Melalui webview | Melalui HTTP |
| Memerlukan CGO | Ya | Tidak |

### Perilaku API Jendela

Dalam mode server, API terkait jendela ditangani dengan aman:

- `app.Window.NewWithOptions()` - Mencatat peringatan dan mengembalikan nil
- `app.Hide()` / `app.Show()` - Tidak melakukan apa pun
- `app.Screen.GetPrimary()` - Mengembalikan galat

Hal ini memungkinkan kode yang merujuk ke jendela berjalan tanpa mengalami crash, meskipun operasi jendela tidak memiliki efek apa pun.

## Build untuk Produksi

### Menggunakan Task (Direkomendasikan)

Proyek yang dibuat dengan `wails3 init` menyertakan task `build:server`:

```bash
# Build for server mode
wails3 task build:server

# Build and run
wails3 task run:server
```

### Build Manual

```bash
# Build with server mode
go build -tags server -o myapp-server .
```

### Docker

Proyek Wails menyertakan penyiapan Docker yang siap digunakan. Untuk mem-build dan menjalankan aplikasi Anda di dalam container:

```bash
# Build the Docker image
wails3 task build:docker

# Run it
wails3 task run:docker
```

Selesai! Aplikasi Anda akan tersedia di `http://localhost:8080`.

Anda dapat menyesuaikan proses build dengan beberapa opsi:

```bash
# Use a custom image tag
wails3 task build:docker TAG=myapp:v1.0.0

# Run on a different port
wails3 task run:docker PORT=3000
```

`Dockerfile.server` yang dihasilkan membuat image minimal berbasis distroless. Pengikatan jaringan ditangani secara otomatis sehingga aplikasi Anda dapat diakses dari luar container.

### Docker Compose

Untuk deployment yang lebih kompleks, berikut konfigurasi Docker Compose dengan pemeriksaan kesehatan:

```yaml
services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - WAILS_SERVER_HOST=0.0.0.0
    healthcheck:
      test: ["CMD", "wget", "-q", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
```

@note{type="info"}
Contoh pemeriksaan kesehatan menggunakan `wget`. Jika menggunakan base image distroless, Anda harus menyertakan binary pemeriksaan kesehatan dalam image atau menggunakan mekanisme pemeriksaan kesehatan eksternal (misalnya, opsi `curl` Docker atau container sidecar).

@end

### Dockerfile Kustom

Jika memerlukan kontrol lebih besar, Anda dapat membuat Dockerfile sendiri. Hal utama yang perlu diingat adalah mengatur `WAILS_SERVER_HOST=0.0.0.0` agar server menerima koneksi dari luar container:

```dockerfile
# Build stage
FROM golang:alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY . .
RUN go mod tidy
RUN go build -tags server -ldflags="-s -w" -o server .

# Runtime stage
FROM gcr.io/distroless/static-debian12
COPY --from=builder /app/server /server
COPY --from=builder /app/frontend/dist /frontend/dist
EXPOSE 8080
ENV WAILS_SERVER_HOST=0.0.0.0
ENTRYPOINT ["/server"]
```

## Pertimbangan Keamanan

Saat men-deploy aplikasi dalam mode server:

1. **Ikat ke localhost secara default** - Gunakan `0.0.0.0` hanya jika diperlukan
2. **Gunakan TLS dalam produksi** - Konfigurasikan `ServerOptions.TLS`
3. **Tempatkan di belakang reverse proxy** - Gunakan nginx/traefik untuk keamanan tambahan
4. **Pertahankan WebSocket pada origin yang sama** - Tambahkan hanya origin tepercaya dengan `WebSocketOriginPatterns`; hindari `WebSocketAllowAllOrigins`
5. **Validasi semua input** - Terapkan praktik keamanan yang sama seperti pada aplikasi web lainnya

## Contoh

Contoh lengkap tersedia di `v3/examples/server/`:

```bash
cd v3/examples/server

# Using Taskfile
task dev

# Or run directly
go run -tags server .

# Open http://localhost:8080 in browser
```

## Variabel Lingkungan

Untuk skenario deployment yang mengharuskan Anda mengganti konfigurasi server tanpa mengubah kode, Wails mengenali variabel lingkungan berikut:

| Variabel | Deskripsi | Default |
| --- | --- | --- |
| `WAILS_SERVER_HOST` | Antarmuka jaringan yang akan diikat | `localhost` |
| `WAILS_SERVER_PORT` | Port yang akan digunakan untuk menerima koneksi | `8080` |

Variabel ini diprioritaskan daripada `ServerOptions` dalam kode Anda. Karena itu, contoh Docker mengatur `WAILS_SERVER_HOST=0.0.0.0` agar container dapat menerima koneksi eksternal tanpa memerlukan perubahan apa pun pada aplikasi Anda.

## Lihat Juga

- [Transport Kustom](/guides/custom-transport/) - Untuk penyesuaian IPC tingkat lanjut
- [Layanan](/features/bindings/services/) - Dokumentasi pengikatan layanan
- [Peristiwa](/guides/events-reference/) - Dokumentasi sistem peristiwa

### Ukuran permintaan runtime

Permintaan ke `/wails/runtime` dibatasi hingga 64 MiB sebelum pemrosesan JSON. Permintaan biasa yang melebihi batas ini menerima HTTP 413, termasuk permintaan tanpa `Content-Length`. Unggahan runtime bertahap tetap memiliki batas terpisah sebesar 1 MiB per potongan dan 64 MiB untuk payload yang telah dirangkai. Gunakan middleware aplikasi atau reverse proxy untuk memberlakukan batas yang lebih rendah jika diperlukan.
