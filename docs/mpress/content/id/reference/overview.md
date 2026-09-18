---
title: "Referensi API"
description: "Dokumentasi API lengkap untuk Wails v3"
slug: "reference/overview"
sourcePath: "reference/overview.md"
---

## Tentang Referensi Ini

Ini adalah referensi API lengkap untuk Wails v3. Referensi ini mendokumentasikan setiap tipe, metode, dan opsi publik yang tersedia dalam kerangka kerja ini.

**Organisasi:**

- [Aplikasi](/reference/application/) - API aplikasi inti
- [Jendela](/reference/window/) - Pembuatan dan pengelolaan jendela
- [Menu](/reference/menu/) - Menu aplikasi, konteks, dan baki sistem
- [Peristiwa](/reference/events/) - Sistem peristiwa dan peristiwa bawaan
- [Dialog](/reference/dialogs/) - Dialog file dan pesan
- [Runtime Frontend](/reference/frontend-runtime/) - API runtime frontend
- [CLI](/reference/cli/) - Antarmuka baris perintah

## Konvensi API

@details{title="Konvensi API Go - Untuk pengembang yang baru mengenal Go"}
### Penamaan

- <strong></strong>Tipe<strong></strong>: PascalCase (misalnya, `WebviewWindow`)
- <strong></strong>Metode<strong></strong>: PascalCase (misalnya, `SetTitle()`)
- <strong></strong>Opsi<strong></strong>: struct PascalCase (misalnya, `WindowOptions`)
- <strong></strong>Konstanta<strong></strong>: PascalCase (misalnya, `WindowStartStateMaximised`)

#### Penanganan Kesalahan

Sebagian besar metode yang dapat gagal mengembalikan `error` sebagai nilai kembalian terakhir. `app.Run()` memblokir hingga aplikasi ditutup dan mengembalikan galat apa pun yang terjadi saat aplikasi dimulai:

```go
if err := app.Run(); err != nil {
    log.Fatal(err)
}
```

Pembuatan jendela tidak mengembalikan galat — `app.Window.New()` mengembalikan `*WebviewWindow` secara langsung.

#### Konteks

Metode siklus hidup layanan menerima sebuah `context.Context`:

```go
func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    // ctx is cancelled when the application is shutting down.
    return nil
}
```

Konteks masa aktif aplikasi tersedia melalui `app.Context()`. Tidak ada `RunWithContext` — panggil `app.Run()`.

#### Pola Opsi

Konfigurasi menggunakan struct opsi:

```go
app := application.New(application.Options{
    Name: "My App",
    Description: "A demo application",
    Services: []application.Service{
        application.NewService(&MyService{}),
    },
})
```

@end

### Konvensi API JavaScript

#### Penamaan

- **Fungsi**: camelCase (misalnya, `setTitle()`)
- **Konstanta**: SCREAMING<em>SNAKE</em>CASE (misalnya, `WINDOW_EVENT_FOCUS`)

#### Asinkron secara Bawaan

Semua pemanggilan metode Go mengembalikan Promise:

```javascript
// Async/await (recommended)
const result = await MyService.DoSomething()

// Promise chain
MyService.DoSomething()
    .then(result => console.log(result))
    .catch(error => console.error(error))
```

#### Penanganan Galat

Kesalahan Go menjadi pengecualian JavaScript:

```javascript
try {
    await MyService.MightFail()
} catch (error) {
    console.error('Go error:', error)
}
```

#### Keamanan Tipe

Definisi TypeScript dibuat secara otomatis:

```typescript
// Fully typed
import { Greet } from './bindings/GreetService'

const message: string = await Greet("World")
```

## Struktur Paket

```
github.com/wailsapp/wails/v3/pkg/
├── application/          # Core application package
│   ├── application.go    # App type
│   ├── webview_window.go # Window management
│   ├── menu.go           # Menu types
│   ├── event_manager.go  # Event system
│   └── dialogs.go        # Dialog APIs
├── events/               # Event constants
└── services/             # Built-in services
    ├── dock/             # macOS dock (includes badge support)
    ├── fileserver/       # File-server service
    ├── kvstore/          # Key/value store
    ├── log/              # Structured logging service
    ├── notifications/    # Notifications service
    └── sqlite/           # SQLite service
```

## Jalur Impor

### Go

```go
import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)
```

### JavaScript

```javascript
// Auto-generated bindings
import { MyMethod } from './bindings/MyService'

// Runtime APIs
import { Events, Window } from '@wailsio/runtime'
```

## Referensi Tipe

### Tipe Umum

@tabs{sync-key="lang"}
[Go]
```go
// Application
type App struct { /* ... */ }
type Options struct { /* ... */ }

// Window
type WebviewWindow struct { /* ... */ } // implements the Window interface
type WebviewWindowOptions struct { /* ... */ }

// Menu
type Menu struct { /* ... */ }
type MenuItem struct { /* ... */ }

// Events — there is no generic Event type; events are typed by source.
type ApplicationEvent struct { /* ... */ }
type WindowEvent struct { /* ... */ }
type CustomEvent struct { /* ... */ }
type EventListener struct { /* ... */ }

// Dialogs
type OpenFileDialogOptions struct { /* ... */ }
type SaveFileDialogOptions struct { /* ... */ }
```

[TypeScript]
```typescript
// Window runtime
interface WindowOptions {
    title?: string
    width?: number
    height?: number
    // ...
}

// Events
type EventCallback = (data: any) => void

// Bindings (auto-generated)
export function MyMethod(arg: string): Promise<string>
```

@end

## Perbedaan Antarplatform

Beberapa API berperilaku berbeda pada platform yang berbeda:

| Fitur | Windows | macOS | Linux |
| --- | --- | --- | --- |
| **Menu Aplikasi** | Bilah menu jendela | Bilah menu global | Bilah menu jendela |
| **Baki Sistem** | Area notifikasi | Bilah menu | Baki sistem |
| **Dock** | Tidak berlaku | ✅ Tersedia | Tidak berlaku |
| **Dialog file** | Native | Native | Bawaan (GTK) |
| **Transparansi** | ✅ Penuh | Memerlukan [`private_mac_apis`](/guides/build/private-macos-apis/#webview-transparency-and-background) | ⚠️ Terbatas |

Perilaku khusus platform didokumentasikan di setiap bagian API.

## Pembuatan Versi

Wails v3 mengikuti pembuatan versi semantik:

- **Mayor** (v3.x.x): Perubahan yang merusak kompatibilitas
- **Minor** (v3.x.x): Fitur baru, kompatibel dengan versi sebelumnya
- **Patch** (v3.x.x): Perbaikan bug, kompatibel dengan versi sebelumnya

**Status saat ini:** Beta (API stabil, penyempurnaan masih berlangsung)

## Kebijakan Deprekasi

Ketika API dideprekasi:

1. **Ditandai dalam dokumentasi** dengan pemberitahuan penghentian penggunaan
2. **Alternatif disediakan** beserta panduan migrasi
3. **Tetap dipelihara selama 1 versi mayor** sebelum dihapus
4. **Peringatan compiler** (jika memungkinkan)

## Stabilitas API

### API Stabil ✅

API berikut stabil dan aman digunakan dalam produksi:

- API inti aplikasi
- Pengelolaan jendela
- Sistem menu
- Sistem peristiwa
- Dialog berkas
- Binding layanan

### API Tidak Stabil ⚠️

API berikut dapat berubah sebelum rilis final:

- Beberapa opsi jendela tingkat lanjut
- Fitur khusus platform
- Fitur eksperimental

API yang tidak stabil ditandai dalam dokumentasi.

## Mendapatkan Bantuan

### Pertanyaan tentang API

1. **Periksa referensi ini** - Dokumentasi API lengkap
2. **Periksa contoh** - [Contoh di GitHub](https://github.com/wailsapp/wails/tree/master/v3/examples)
3. **Cari di Discord** - [Server Discord](https://discord.gg/JDdSxwjhGf)
4. **Tanyakan kepada komunitas** - Kanal #help di Discord

### Melaporkan Masalah API

Menemukan bug atau inkonsistensi?

1. **Periksa masalah yang sudah ada** - [Masalah di GitHub](https://github.com/wailsapp/wails/issues)
2. **Buat laporan terperinci** - Sertakan kode, pesan kesalahan, dan platform
3. **Berikan langkah reproduksi** - Contoh minimal yang menunjukkan masalah

## Dokumentasi Terkait

- [Tutorial](/tutorials/overview/) - Pelajari dengan membangun aplikasi nyata
- [Panduan](/guides/architecture/) - Panduan berorientasi tugas untuk skenario umum
- [Fitur](/features/windows/basics/) - Dokumentasi untuk setiap fitur
- [Contoh](https://github.com/wailsapp/wails/tree/master/v3/examples) - Contoh kode yang berfungsi dan tersedia di GitHub

---

**Telusuri API:** Gunakan navigasi di sebelah kiri untuk menjelajahi API tertentu.
