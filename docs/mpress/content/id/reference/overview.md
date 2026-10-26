---
title: "Referensi API"
description: "Cari tipe, metode, opsi, dan perintah CLI Wails"
slug: "reference/overview"
sourcePath: "reference/overview.md"
---

## Tentang Referensi Ini

Gunakan bagian ini untuk mencari tipe, metode, opsi, dan perintah Wails saat mengembangkan aplikasi. Jika Anda memerlukan pembelajaran terpandu, mulailah dengan [Tutorial](/tutorials/overview/). Untuk langkah-langkah menyelesaikan tugas, gunakan [Panduan Praktis](/quick-start/next-steps/#features). Untuk memahami latar belakang desain, gunakan [Konsep Inti](/concepts/architecture/).

**Susunan:**

- [Aplikasi](/reference/application/) - API inti aplikasi
- [Jendela](/reference/window/) - Pembuatan dan pengelolaan jendela; [opsi jendela](/features/windows/options/) mencantumkan kolom konfigurasi
- [Menu](/reference/menu/) - API aplikasi, menu konteks, dan baki sistem; [referensi item menu](/features/menus/reference/) membahas tipe dan properti item
- [Peristiwa](/reference/events/) - Sistem peristiwa dan peristiwa bawaan
- [Dialog](/reference/dialogs/) - Dialog file dan pesan
- [Runtime Frontend](/reference/frontend-runtime/) - API runtime frontend
- [Ikhtisar perintah CLI](/reference/cli/) - Perintah yang dikelompokkan berdasarkan tujuan; [flag dan contoh](/guides/cli/) memberikan detail setiap perintah
- [Protokol manifes pembaruan](/reference/update-manifest/) - Metadata pembaruan dan kolom verifikasi
- [API Seluler](/guides/mobile/mobile-api/) - API seluler eksperimental

## Konvensi API

@details{title="Konvensi API Go - Untuk pengembang yang baru mengenal Go"}
### Penamaan

- <strong>Tipe</strong>: PascalCase (misalnya, `WebviewWindow`)
- <strong>Metode</strong>: PascalCase (misalnya, `SetTitle()`)
- <strong>Opsi</strong>: struct PascalCase (misalnya, `WindowOptions`)
- <strong>Konstanta</strong>: PascalCase (misalnya, `WindowStartStateMaximised`)

#### Penanganan Kesalahan

Sebagian besar metode yang dapat gagal mengembalikan `error` sebagai nilai kembalian terakhir. `app.Run()` memblokir hingga aplikasi ditutup dan mengembalikan kesalahan saat memulai, jika ada:

```go
if err := app.Run(); err != nil {
    log.Fatal(err)
}
```

Pembuatan jendela tidak mengembalikan kesalahan — `app.Window.New()` langsung mengembalikan `*WebviewWindow`.

#### Konteks

Metode siklus hidup layanan menerima `context.Context`:

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

#### Asinkron secara Default

Semua pemanggilan metode Go mengembalikan Promise:

```javascript
// Async/await (recommended)
const result = await MyService.DoSomething()

// Promise chain
MyService.DoSomething()
    .then(result => console.log(result))
    .catch(error => console.error(error))
```

#### Penanganan Kesalahan

Kesalahan Go menjadi exception JavaScript:

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

## Path Impor

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
| **Dock** | Tidak tersedia | ✅ Tersedia | Tidak tersedia |
| **Dialog file** | Native | Native | Native (GTK) |
| **Transparansi** | ✅ Penuh | Memerlukan [`private_mac_apis`](/guides/build/private-macos-apis/#webview-transparency-and-background) | ⚠️ Terbatas |

Perilaku khusus platform didokumentasikan di setiap bagian API.

## Pembuatan Versi

Wails v3 mengikuti pembuatan versi semantik:

- **Mayor** (v4.0.0): Perubahan yang merusak kompatibilitas
- **Minor** (v3.1.0): Fitur baru yang kompatibel dengan versi sebelumnya
- **Patch** (v3.0.1): Perbaikan bug yang kompatibel dengan versi sebelumnya

**Wails v3.0.0:** Stabil

## Kebijakan Penghentian Dukungan

Ketika API dihentikan dukungannya:

1. **Ditandai dalam dokumentasi** dengan pemberitahuan penghentian dukungan
2. **Alternatif disediakan** beserta panduan migrasi
3. **Tetap dipelihara selama 1 versi mayor** sebelum dihapus
4. **Peringatan kompilator** (jika memungkinkan)

## Stabilitas API

### API Stabil ✅

API berikut stabil dan aman digunakan dalam produksi:

- API inti aplikasi
- Manajemen jendela
- Sistem menu
- Sistem peristiwa
- Dialog file
- Binding layanan

### API Eksperimental ⚠️

Hanya API yang secara eksplisit didokumentasikan sebagai eksperimental yang berada di luar kontrak kompatibilitas desktop stabil:

- Dukungan Android
- Dukungan iOS
- API lain yang secara eksplisit ditandai sebagai eksperimental

API eksperimental ditandai dalam dokumentasi.

## Mendapatkan Bantuan

### Pertanyaan tentang API

1. **Lihat referensi ini** - Dokumentasi API lengkap
2. **Lihat contoh** - [Contoh di GitHub](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples)
3. **Cari di Discord** - [Server Discord](https://discord.gg/JDdSxwjhGf)
4. **Tanyakan kepada komunitas** - Kanal #help di Discord

### Melaporkan Masalah API

Menemukan bug atau inkonsistensi?

1. **Periksa masalah yang sudah ada** - [Masalah di GitHub](https://github.com/wailsapp/wails/issues)
2. **Buat laporan terperinci** - Sertakan kode, pesan kesalahan, dan platform
3. **Sediakan langkah reproduksi** - Contoh minimal yang menunjukkan masalah

## Dokumentasi Terkait

- [Tutorial](/tutorials/overview/) - Belajar dengan membangun aplikasi nyata
- [Panduan](/guides/architecture/) - Panduan berorientasi tugas untuk skenario umum
- [Fitur](/features/windows/basics/) - Dokumentasi untuk setiap fitur
- [Contoh](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples) - Contoh kode siap pakai di GitHub

---

**Jelajahi API:** Gunakan navigasi di sebelah kiri untuk menjelajahi API tertentu.
