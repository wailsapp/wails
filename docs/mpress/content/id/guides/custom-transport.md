---
title: "Membuat Lapisan Transport Kustom"
description: "Pelajari cara membuat dan menyesuaikan lapisan transport IPC kustom Wails v3 Anda sendiri"
slug: "guides/custom-transport"
sourcePath: "guides/custom-transport.md"
---

Wails v3 memungkinkan Anda menyediakan lapisan transport IPC kustom sambil tetap mempertahankan semua binding yang dihasilkan dan komunikasi peristiwa. Dengan demikian, Anda dapat mengganti transport default berbasis HTTP fetch dengan WebSocket, protokol kustom, atau mekanisme transport lainnya.

## Gambaran Umum

Secara default, Wails menggunakan permintaan HTTP fetch dari frontend untuk berkomunikasi dengan backend melalui `/wails/runtime`. API transport kustom memungkinkan Anda untuk:

- Mengganti transport HTTP dengan WebSocket, gRPC, atau protokol kustom apa pun
- Mempertahankan kompatibilitas penuh dengan pembuatan kode Wails
- Mempertahankan semua binding, peristiwa, dialog, dan fitur Wails lainnya yang sudah ada
- Mengimplementasikan pengelolaan koneksi, autentikasi, dan penanganan kesalahan Anda sendiri

## Arsitektur

```text
┌─────────────────────────────────────────────────┐
│  Frontend (TypeScript)                          │
│  - Generated bindings still work                │
│  - Your custom client transport                 │
└──────────────────┬──────────────────────────────┘
                   │
                   │ Your Protocol (WebSocket/etc)
                   │
┌──────────────────▼──────────────────────────────┐
│  Backend (Go)                                   │
│  - Your Transport implementation                │
│  - Wails MessageProcessor                       │
│  - All existing Wails infrastructure            │
└─────────────────────────────────────────────────┘
```

## Menjalankan contoh WebSocket lengkap

[Contoh transport WebSocket](https://github.com/wailsapp/wails/tree/master/v3/examples/websocket-transport) mencakup kedua sisi protokol: ID permintaan dan respons, pengelolaan siklus hidup, klien JavaScript, serta pengiriman peristiwa. Mulailah dari file-file tersebut, bukan dari transport yang hanya mengimplementasikan `Start` dan `Stop`.

Dari working tree repositori Wails yang telah dilengkapi dependensi platform native:

```bash
cd v3/examples/websocket-transport
go run .
```

Contoh tersebut membuka jendela native dan memulai endpoint WebSocket pada port 9099. Gunakan contoh ini di lingkungan pengembangan lokal tepercaya. Sebelum menerapkan varian yang dapat diakses melalui jaringan, batasi listener, validasi origin, dan autentikasi klien. Kebijakan origin pada contoh bersifat permisif dan bukan merupakan kebijakan keamanan untuk produksi.

### Kontrak backend

Transport harus mengimplementasikan ketiga metode wajib. `WebSocketTransport` dalam contoh mengimplementasikan `Start`, `JSClient`, dan `Stop`; `JSClient` menyediakan shim JavaScript yang disuntikkan Wails ke dalam halaman. Sertakan pemeriksaan antarmuka pada waktu kompilasi dalam implementasi Anda:

```go
var _ application.Transport = (*WebSocketTransport)(nil)
```

Konfigurasikan aplikasi menggunakan konstruktor dari contoh:

```go
Transport: NewWebSocketTransport("127.0.0.1:9099"),
```

Field ini berada di `application.Options`, bersama aset dan layanan biasa. Salin file implementasi beserta pemanggilan konstruktornya. `Start` menerima `MessageProcessor`; `Stop` harus menutup listener, koneksi aktif, dan pekerjaan yang tertunda.

### Kontrak frontend

`websocket-transport.js` dalam contoh mengimplementasikan transport runtime, termasuk korelasi permintaan/respons dan penolakan panggilan yang gagal. Setelah menyajikan file tersebut dari aset Anda, inisialisasikan transport sebelum memanggil binding layanan yang dihasilkan:

```javascript
import { setTransport } from '/wails/runtime.js'
import { createWebSocketTransport } from '/websocket-transport.js'

setTransport(createWebSocketTransport('ws://127.0.0.1:9099/wails/ws'))
```

Binding layanan yang dihasilkan tetap tidak berubah. Jika transport Anda menggantikan klien, transport tersebut harus membawa metode yang digunakan aplikasi Anda dan mempertahankan semantik kesalahannya. Pengiriman peristiwa memerlukan kontrak transport peristiwa opsional; kontrak ini tidak disediakan secara otomatis oleh implementasi `Transport` sembarang. Contoh WebSocket lengkap mendemonstrasikannya.

Untuk deployment jarak jauh, gunakan TLS serta kebijakan autentikasi dan origin yang eksplisit. Hosting khusus browser juga mengharuskan aset aplikasi disajikan seperti yang dijelaskan di bawah ini.

## Referensi API

### Antarmuka Transport

```go
type Transport interface {
    Start(ctx context.Context, messageProcessor *application.MessageProcessor) error
    // JSClient returns the JavaScript shim that the runtime injects into the
    // window so that frontend code can call into the transport.
    JSClient() []byte
    Stop() error
}
```

### Antarmuka AssetServerTransport (Opsional)

Untuk deployment berbasis browser, atau jika Anda ingin menyajikan aset dan IPC melalui transport kustom yang sama, implementasikan antarmuka `AssetServerTransport`:

```go
type AssetServerTransport interface {
    Transport

    // ServeAssets configures the transport to serve assets alongside IPC.
    // The assetHandler is Wails' internal asset server that handles all assets,
    // runtime.js, capabilities, flags, etc.
    ServeAssets(assetHandler http.Handler) error
}
```

**Kapan harus mengimplementasikan antarmuka ini:**

- Menjalankan aplikasi di browser, bukan di webview
- Menyajikan aset melalui HTTP bersama transport IPC kustom Anda
- Membangun aplikasi yang dapat diakses melalui jaringan

**Contoh implementasi:**

```go
func (t *MyTransport) ServeAssets(assetHandler http.Handler) error {
    mux := http.NewServeMux()

    // Mount your IPC endpoint
    mux.HandleFunc("/my/ipc/endpoint", t.handleIPC)

    // Mount Wails asset server for everything else
    mux.Handle("/", assetHandler)

    // Start HTTP server
    t.httpServer.Handler = mux
    go t.httpServer.ListenAndServe()

    return nil
}
```

Ketika `ServeAssets()` dipanggil, assetHandler menyediakan:

- Semua aset statis (HTML, CSS, JS, gambar, dan sebagainya)
- `/wails/runtime.js` - Pustaka runtime Wails

## Lihat Juga

- `transport.go` - Antarmuka dan tipe transport inti
- `messageprocessor.go` - Pemroses pesan dasar yang menangani seluruh IPC Wails
