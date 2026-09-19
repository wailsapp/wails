---
title: "Membuat Lapisan Transport Kustom"
description: "Pelajari cara membuat dan menyesuaikan lapisan transport IPC kustom Anda sendiri untuk Wails v3"
slug: "guides/custom-transport"
sourcePath: "guides/custom-transport.md"
---

Wails v3 memungkinkan Anda menyediakan lapisan transport IPC kustom sambil tetap mempertahankan semua binding yang dihasilkan dan komunikasi peristiwa. Dengan demikian, Anda dapat mengganti transport default berbasis HTTP fetch dengan WebSocket, protokol kustom, atau mekanisme transport lainnya.

## Ringkasan

Secara default, Wails menggunakan permintaan HTTP fetch dari frontend untuk berkomunikasi dengan backend melalui `/wails/runtime`. API transport kustom memungkinkan Anda untuk:

- Mengganti transport HTTP dengan WebSocket, gRPC, atau protokol kustom apa pun
- Mempertahankan kompatibilitas penuh dengan pembuatan kode Wails
- Mempertahankan semua binding, peristiwa, dialog, dan fitur Wails lainnya yang sudah ada
- Mengimplementasikan sendiri pengelolaan koneksi, autentikasi, dan penanganan kesalahan

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

## Penggunaan

### 1. Implementasikan Antarmuka Transport

Buat transport kustom dengan mengimplementasikan antarmuka `Transport`:

```go
package main

import (
    "context"
    "github.com/wailsapp/wails/v3/pkg/application"
)

type MyCustomTransport struct {
    // Your fields
}

func (t *MyCustomTransport) Start(ctx context.Context, processor *application.MessageProcessor) error {
    // Initialize your transport (WebSocket server, gRPC server, etc.)
    // When you receive requests, call processor.HandleRuntimeCallWithIDs()
    return nil
}

func (t *MyCustomTransport) Stop() error {
    // Clean up your transport
    return nil
}
```

### 2. Konfigurasikan Aplikasi Anda

Teruskan transport kustom Anda ke opsi aplikasi:

```go
func main() {
    app := application.New(application.Options{
        Name: "My App",
        Transport: &MyCustomTransport{},
        // ... other options
    })

    err := app.Run()
    if err != nil {
        log.Fatal(err)
    }
}
```

### 3. Ubah Runtime Frontend

Jika menggunakan transport kustom, Anda perlu mengubah runtime frontend agar menggunakan transport Anda sebagai pengganti HTTP fetch. Implementasikan antarmuka `RuntimeTransport` yang akan digunakan untuk menangani permintaan:

```typescript
const { setTransport } = await import('/wails/runtime.js');

class MyRuntimeTransport {
  call(objectID: number, method: number, windowName: string, args: any): Promise<any> {
    // TODO: implement IPC call with your transport protocol

    return resp;
  }
}

const myTransport = new MyRuntimeTransport();
setTransport(myTransport);
```

## Catatan

- Transport HTTP default tetap berfungsi jika tidak ada transport kustom yang ditentukan
- Binding yang dihasilkan tetap tidak berubah—hanya lapisan transport yang berubah
- Peristiwa, dialog, papan klip, dan semua fitur Wails lainnya berfungsi secara transparan
- Anda bertanggung jawab atas penanganan kesalahan, logika penyambungan ulang, dan keamanan dalam transport kustom Anda
- Contoh WebSocket yang disediakan ditujukan untuk demonstrasi dan mungkin perlu diperkuat agar layak digunakan dalam produksi

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

Untuk deployment berbasis browser atau jika Anda ingin menyajikan aset dan IPC melalui transport kustom Anda, implementasikan antarmuka `AssetServerTransport`:

```go
type AssetServerTransport interface {
    Transport

    // ServeAssets configures the transport to serve assets alongside IPC.
    // The assetHandler is Wails' internal asset server that handles all assets,
    // runtime.js, capabilities, flags, etc.
    ServeAssets(assetHandler http.Handler) error
}
```

**Kapan antarmuka ini perlu diimplementasikan:**

- Menjalankan aplikasi di browser sebagai pengganti webview
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

Saat `ServeAssets()` dipanggil, assetHandler menyediakan:

- Semua aset statis (HTML, CSS, JS, gambar, dan sebagainya)
- `/wails/runtime.js` - Pustaka runtime Wails

## Lihat Juga

- `transport.go` - Antarmuka dan tipe transport inti
- `messageprocessor.go` - Pemroses pesan dasar yang menangani seluruh IPC Wails
