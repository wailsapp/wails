---
title: "Notifikasi"
description: "Tampilkan notifikasi sistem native dengan tombol tindakan dan input teks"
slug: "features/notifications/overview"
sourcePath: "features/notifications/overview.md"
---

## Pendahuluan

Wails menyediakan sistem notifikasi lintas platform yang komprehensif untuk aplikasi desktop. Layanan ini memungkinkan Anda menampilkan notifikasi sistem native, dengan dukungan untuk:

- Notifikasi dasar dengan judul, subjudul, dan isi
- Notifikasi interaktif dengan tombol tindakan dan balasan teks
- [Kategori notifikasi](#notifikasi-interaktif) yang dapat digunakan kembali untuk tindakan
- [Suara](#suara-khusus) khusus (bawaan, senyap, atau bernama)
- [Lampiran](#lampiran) (gambar di setiap platform; audio/video di macOS)
- [Pengelompokan notifikasi terkait](#utas-dan-pengelompokan) berdasarkan `ThreadID`
- [Prioritas](#tingkat-interupsi) melalui `InterruptionLevel` (`passive` / `active` / `timeSensitive` / `critical`)
- [Pengiriman terjadwal](#pengiriman-terjadwal) (native di macOS; pewaktu dalam proses di Windows + Linux)
- [Pembaruan notifikasi yang sedang aktif](#memperbarui-notifikasi) berdasarkan ID

Setiap kolom opsional baru akan mengalami penurunan fungsi secara wajar ketika suatu platform tidak dapat mendukungnya; lihat [Pertimbangan Platform](#pertimbangan-platform) untuk matriks dukungan setiap fitur.

## Penggunaan Dasar

### Membuat Layanan

Pertama, inisialisasikan layanan notifikasi:

```go
import "github.com/wailsapp/wails/v3/pkg/application"
import "github.com/wailsapp/wails/v3/pkg/services/notifications"

// Create a new notification service
notifier := notifications.New()

//Register the service with the application
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(notifier),
    },
})
```

## Otorisasi Notifikasi

Notifikasi di macOS memerlukan otorisasi pengguna. Minta dan periksa otorisasi:

```go
authorized, err := notifier.CheckNotificationAuthorization()
if err != nil {
    // Handle authorization error
}
if authorized {
    // Send notifications
} else {
    // Request authorization
    authorized, err = notifier.RequestNotificationAuthorization()
}
```

Di Windows dan Linux, ini selalu mengembalikan `true`.

## Jenis Notifikasi

### Notifikasi Dasar

Kirim notifikasi dasar kepada pengguna dengan ID unik, judul, subjudul opsional (macOS dan Linux), serta teks isi:

```go
notifier.SendNotification(notifications.NotificationOptions{
    ID: "unique-id",
    Title: "New Calendar Invite",
    Subtitle: "From: Jane Doe", // Optional
    Body: "Tap to view the event",
})

```

### Notifikasi Interaktif

Kirim notifikasi dengan tombol tindakan dan input teks. Sebelum mengirim notifikasi ini, Anda harus mendaftarkan kategori notifikasi:

```go
// Define a unique category id
categoryID := "unique-category-id"

// Define a category with actions
category := notifications.NotificationCategory{
    ID: categoryID,
    Actions: []notifications.NotificationAction{
        {
            ID:    "OPEN", 
            Title: "Open",
        },
        {
            ID:          "ARCHIVE", 
            Title:       "Archive", 
            Destructive: true,  /* macOS-specific */
        },
    },
    HasReplyField:    true,
    ReplyPlaceholder: "message...",
    ReplyButtonTitle: "Reply",
}

// Register the category
notifier.RegisterNotificationCategory(category)

// Send an interactive notification with the actions registered in the provided category
notifier.SendNotificationWithActions(notifications.NotificationOptions{
    ID:         "unique-id",
    Title:      "New Message",
    Subtitle:   "From: Jane Doe",
    Body:       "Are you able to make it?",
    CategoryID: categoryID,
})
```

## Respons Notifikasi

Proses interaksi pengguna dengan notifikasi:

```go
notifier.OnNotificationResponse(func(result notifications.NotificationResult) {
    response := result.Response
    fmt.Printf("Notification %s was actioned with: %s\n", response.ID, response.ActionIdentifier)

    if response.ActionIdentifier == "TEXT_REPLY" {
        fmt.Printf("User replied: %s\n", response.UserText)
    }

    if data, ok := response.UserInfo["sender"].(string); ok {
        fmt.Printf("Original sender: %s\n", data)
    }

    // Emit an event to the frontend
    app.Event.Emit("notification", result.Response)
})
```

## Penyesuaian Notifikasi

### Metadata Khusus

Notifikasi dasar dan interaktif dapat menyertakan data khusus:

```go
notifier.SendNotification(notifications.NotificationOptions{
    ID: "unique-id",
    Title: "New Calendar Invite",
    Subtitle: "From: Jane Doe", // Optional
    Body: "Tap to view the event",
    Data: map[string]interface{}{
        "sender": "jane.doe@example.com",
        "timestamp": "2025-03-10T15:30:00Z",
    }
})
```

### Suara Khusus

Gunakan `Sound` untuk mengontrol audio yang diputar saat notifikasi dikirimkan. Jika dibiarkan `nil`, suara bawaan platform akan diputar; tetapkan `Silent: true` untuk menonaktifkan suara; tetapkan `Name` untuk memutar suara bernama/yang disertakan dalam bundel.

```go
// Silent
notifier.SendNotification(notifications.NotificationOptions{
    ID:    "silent-id",
    Title: "Background sync complete",
    Sound: &notifications.NotificationSound{Silent: true},
})

// Named sound
notifier.SendNotification(notifications.NotificationOptions{
    ID:    "named-id",
    Title: "New message",
    Body:  "Tap to read",
    Sound: &notifications.NotificationSound{Name: "Ping"},
})
```

Cara `Name` diresolusi pada setiap platform:

- **macOS** — `Name` diteruskan ke `[UNNotificationSound soundNamed:]`; berkas audio harus berada di `Library/Sounds` dalam bundel aplikasi Anda.
- **Windows** — jika `Name` sudah diawali dengan `ms-winsoundevent:` atau `ms-appx:`, nilainya digunakan apa adanya; jika tidak, nilainya dibungkus dalam `ms-winsoundevent:` agar dapat digunakan sebagai nama peristiwa toast bawaan (lihat dokumentasi skema `<audio>` toast Microsoft).
- **Linux** — diteruskan sebagai petunjuk `sound-name` freedesktop; pemutaran suara bergantung pada daemon notifikasi dan tema suara yang aktif.

### Lampiran

`Attachments` menambahkan berkas media ke notifikasi. macOS mendukung beberapa lampiran dengan jenis media apa pun; Windows dan Linux menggunakan lampiran pertama yang berjenis gambar.

```go
notifier.SendNotification(notifications.NotificationOptions{
    ID:    "image-id",
    Title: "Photo uploaded",
    Body:  "Tap to view",
    Attachments: []notifications.NotificationAttachment{
        {
            ID:   "preview",
            Path: "/absolute/path/to/image.png",
            // Type hint:
            //   macOS:   UTI like "public.png" / "public.audio" (often inferred)
            //   Windows: "hero" | "appLogoOverride" | "inline" (default "inline")
            //   Linux:   ignored
            Type: "inline",
        },
    },
})
```

`Path` harus berupa jalur absolut pada sistem berkas. macOS juga menerima URL `file://`.

#### Melampirkan berkas yang disertakan bersama aplikasi Anda

OS membaca lampiran dari disk saat notifikasi dikirimkan, sehingga `Path` harus merujuk ke berkas nyata pada mesin pengguna akhir. Untuk aset yang Anda bundel bersama aplikasi (ikon atau gambar yang disematkan dengan `go:embed`), tidak ada jalur absolut tetap yang dapat ditulis langsung dalam kode karena berkas tersebut berada di dalam biner, bukan di lokasi disk yang diketahui. Tulis berkas tersebut sekali ke direktori yang dapat ditulisi saat aplikasi dimulai, lalu teruskan jalurnya:

```go
import (
    _ "embed"
    "os"
    "path/filepath"
)

//go:embed assets/preview.png
var previewPNG []byte

// Materialise the embedded asset to a stable path the OS can read.
previewPath := filepath.Join(os.TempDir(), "myapp-preview.png")
if err := os.WriteFile(previewPath, previewPNG, 0o644); err != nil {
    // handle error
}

notifier.SendNotification(notifications.NotificationOptions{
    ID:    "image-id",
    Title: "Photo uploaded",
    Body:  "Tap to view",
    Attachments: []notifications.NotificationAttachment{
        {ID: "preview", Path: previewPath, Type: "inline"},
    },
})
```

Berkas yang diberikan pengguna atau diunduh sudah memiliki jalur nyata pada disk, sehingga Anda dapat meneruskannya langsung ke `Path` tanpa langkah ini. Dukungan untuk meneruskan byte lampiran dalam memori mungkin ditambahkan pada rilis mendatang.

### Utas dan Pengelompokan

`ThreadID` mengelompokkan notifikasi terkait agar OS dapat menciutkannya di Notification Center / Action Center.

```go
notifier.SendNotification(notifications.NotificationOptions{
    ID:       "msg-42",
    Title:    "Jane Doe",
    Body:     "Are you free for lunch?",
    ThreadID: "chat:jane.doe",
})
```

### Tingkat Interupsi

`InterruptionLevel` mengontrol prioritas notifikasi. Gunakan salah satu konstanta yang diekspor:

```go
notifier.SendNotification(notifications.NotificationOptions{
    ID:                "alert-id",
    Title:             "Server down",
    Body:              "Investigate immediately",
    InterruptionLevel: notifications.InterruptionLevelTimeSensitive,
})
```

| Konstanta | Nilai | Arti |
| --- | --- | --- |
| `InterruptionLevelPassive` | `"passive"` | Pengiriman senyap; tidak akan menyalakan layar atau memutar suara bawaan |
| `InterruptionLevelActive` | `"active"` | Tingkat bawaan |
| `InterruptionLevelTimeSensitive` | `"timeSensitive"` | Menerobos Focus / Do Not Disturb jika diizinkan |
| `InterruptionLevelCritical` | `"critical"` | Melewati Focus dan dering; macOS memerlukan hak Critical Alert (mengalami penurunan fungsi secara diam-diam tanpa hak tersebut) |

Pemetaan per platform:

- **macOS** — menetapkan `UNNotificationContent.interruptionLevel`. `critical` memerlukan macOS 12+ dan hak Critical Alert.
- **Windows** — dipetakan ke atribut toast `<toast scenario="...">`.
- **Linux** — dipetakan ke petunjuk `urgency` freedesktop.

### Pengiriman Terjadwal

`Schedule` menunda pengiriman. Tetapkan tepat salah satu dari `DelaySeconds` (detik dari sekarang) atau `At` (detik Unix, UTC).

```go
// Deliver in 60 seconds
notifier.SendNotification(notifications.NotificationOptions{
    ID:       "reminder-id",
    Title:    "Stand up and stretch",
    Schedule: &notifications.NotificationSchedule{DelaySeconds: 60},
})

// Deliver at an absolute time
notifier.SendNotification(notifications.NotificationOptions{
    ID:       "alarm-id",
    Title:    "Meeting in 5 minutes",
    Schedule: &notifications.NotificationSchedule{At: time.Now().Add(time.Hour).Unix()},
})
```

@note{type="caution" title="Persistensi"}
Di **macOS**, notifikasi terjadwal menggunakan pemicu native dan tetap tersedia setelah aplikasi dimulai ulang. Di **Windows** dan **Linux**, penjadwalan beralih ke pewaktu `time.AfterFunc` dalam proses dan **akan hilang jika aplikasi ditutup sebelum pengiriman** — baik `wintoast` maupun spesifikasi freedesktop tidak menyediakan primitif pengiriman tertunda.

@end

### Memperbarui Notifikasi

`UpdateNotification` mengganti notifikasi yang sedang aktif dengan `ID` yang sama:

```go
notifier.UpdateNotification(notifications.NotificationOptions{
    ID:    "download-id",
    Title: "Download complete",
    Body:  "report.pdf is ready",
})
```

Perilaku per platform:

- **macOS** — `UNUserNotificationCenter` secara otomatis menghapus duplikasi berdasarkan pengenal, sehingga notifikasi yang ada diperbarui di tempat.
- **Linux** — menggunakan parameter D-Bus `replaces_id` untuk menggantikan notifikasi sebelumnya.
- **Windows** — saat ini mengirimkan ulang sebagai notifikasi baru. Penggantian langsung yang sebenarnya memerlukan dukungan upstream `wintoast` untuk `tag` / `group`.

## Pertimbangan Platform

@tabs
[macOS]
Di macOS, notifikasi:

- Memerlukan otorisasi pengguna
- Mengharuskan aplikasi dikemas dan ditandatangani (dinotariskan untuk distribusi)
- Menggunakan tampilan notifikasi standar sistem
- Mendukung `Subtitle`
- Mendukung input teks pengguna (balasan)
- Mendukung opsi tindakan `Destructive`
- Mendukung beberapa `Attachments` dari semua jenis media (gambar, audio, video)
- Mendukung `ThreadID` untuk pengelompokan di Notification Center
- Mendukung semua nilai `InterruptionLevel` (`critical` memerlukan hak Critical Alert)
- Mendukung pengiriman terjadwal native yang tetap berlaku setelah aplikasi dimulai ulang
- Secara otomatis mendeduplikasi panggilan `UpdateNotification` berdasarkan `ID`
- Menangani mode gelap/terang secara otomatis

[Windows]
Di Windows, notifikasi:

- Menggunakan gaya toast sistem Windows melalui backend `wintoast`
- Menyesuaikan dengan pengaturan tema Windows
- Mendukung input teks pengguna (balasan)
- Mendukung layar DPI tinggi
- Tidak mendukung `Subtitle`
- Mendukung satu `Attachment` gambar dengan petunjuk penempatan `hero`, `appLogoOverride`, atau `inline` (default `inline`)
- Mendukung `ThreadID` untuk pengelompokan di Action Center
- Mendukung `InterruptionLevel` melalui atribut toast `scenario`
- Mendukung pengiriman terjadwal melalui timer dalam proses — **notifikasi terjadwal akan hilang jika aplikasi berhenti sebelum dikirim**
- `UpdateNotification` saat ini mengirimkan ulang sebagai notifikasi baru (penggantian langsung yang sebenarnya masih menunggu dukungan upstream `wintoast` untuk `tag`/`group`)

[Linux]
Di Linux, notifikasi menggunakan antarmuka D-Bus `org.freedesktop.Notifications`. Daemon notifikasi yang kompatibel **harus berjalan** agar notifikasi dapat berfungsi.

@note{type="caution" title="Persyaratan sistem: daemon notifikasi"}
Daemon notifikasi yang kompatibel dengan freedesktop harus terinstal dan berjalan. Pilihan yang umum:

- **dunst** — ringan dan sangat mudah dikonfigurasi (`apt install dunst` / `dnf install dunst`)
- **mako** — native untuk Wayland (`apt install mako-notifier`)
- **GNOME Shell** — mendaftarkan antarmuka secara otomatis di GNOME 43+. Di Ubuntu 24.04 (GNOME Shell 46), antarmuka mungkin tidak mendaftarkan dirinya sendiri saat sesi dimulai; instal `dunst` sebagai solusi cadangan jika notifikasi tidak muncul.
- **xfce4-notifyd** — disertakan dengan desktop XFCE

Jika tidak ada daemon yang berjalan, `SendNotification` akan mengembalikan kesalahan D-Bus: `The name org.freedesktop.Notifications was not provided by any .service files`. Tangani kesalahan ini di aplikasi Anda dan beri tahu pengguna agar menginstal daemon notifikasi.

@end

Di Linux, notifikasi:

- Mengikuti tema lingkungan desktop
- Diposisikan sesuai aturan lingkungan desktop
- Mendukung `Subtitle` (digabungkan ke isi untuk daemon yang tidak merendernya secara terpisah)
- Tidak mendukung input teks pengguna (bukan bagian dari spesifikasi freedesktop)
- Mendukung satu `Attachment` gambar melalui petunjuk `image-path`
- Mendukung `ThreadID` (ditangani oleh daemon jika didukung)
- `Sound.Name` diteruskan sebagai petunjuk `sound-name`; apakah suara diputar bergantung pada daemon yang aktif dan tema suara
- Memetakan `InterruptionLevel` ke petunjuk freedesktop `urgency`
- Mendukung pengiriman terjadwal melalui timer dalam proses — **notifikasi terjadwal akan hilang jika aplikasi berhenti sebelum dikirim**
- `UpdateNotification` menggunakan parameter D-Bus `replaces_id` untuk menggantikan notifikasi sebelumnya secara langsung

@end

## Praktik Terbaik

1. Periksa dan minta otorisasi:
  - macOS memerlukan otorisasi pengguna


2. Berikan notifikasi yang jelas dan ringkas:
  - Gunakan judul, subjudul, teks, dan judul tindakan yang deskriptif


3. Tangani respons notifikasi dengan tepat:
  - Periksa kesalahan dalam respons notifikasi
  - Berikan umpan balik untuk tindakan pengguna


4. Pertimbangkan konvensi platform:
  - Ikuti pola notifikasi khusus platform
  - Patuhi pengaturan sistem


5. Di Linux, tangani dependensi daemon:
  - Periksa kesalahan yang dikembalikan oleh `SendNotification` — daemon yang tidak tersedia menghasilkan kesalahan D-Bus
  - Dokumentasi paket atau README aplikasi harus menyatakan bahwa daemon notifikasi freedesktop diperlukan


## Contoh

Pelajari contoh ini:

- [Notifikasi](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/notifications)

## Referensi API

### Pengelolaan Layanan

| Metode | Deskripsi |
| --- | --- |
| `New()` | Membuat layanan notifikasi baru |

### Otorisasi Notifikasi

| Metode | Deskripsi |
| --- | --- |
| `RequestNotificationAuthorization()` | Meminta izin untuk menampilkan notifikasi (macOS) |
| `CheckNotificationAuthorization()` | Memeriksa status otorisasi notifikasi saat ini (macOS) |

### Mengirim Notifikasi

| Metode | Deskripsi |
| --- | --- |
| `SendNotification(options NotificationOptions)` | Mengirim notifikasi dasar |
| `SendNotificationWithActions(options NotificationOptions)` | Mengirim notifikasi interaktif dengan tindakan |
| `UpdateNotification(options NotificationOptions)` | Memperbarui notifikasi yang sedang diproses berdasarkan `ID` (lihat [Memperbarui Notifikasi](#memperbarui-notifikasi)) |

### Kategori Notifikasi

| Metode | Deskripsi |
| --- | --- |
| `RegisterNotificationCategory(category NotificationCategory)` | Mendaftarkan kategori notifikasi yang dapat digunakan kembali |
| `RemoveNotificationCategory(categoryID string)` | Menghapus kategori yang telah didaftarkan sebelumnya |

### Mengelola Notifikasi

| Metode | Deskripsi |
| --- | --- |
| `RemoveAllPendingNotifications()` | Menghapus semua notifikasi tertunda (khusus macOS dan Linux) |
| `RemovePendingNotification(identifier string)` | Menghapus notifikasi tertunda tertentu (khusus macOS dan Linux) |
| `RemoveAllDeliveredNotifications()` | Menghapus semua notifikasi yang telah ditampilkan (khusus macOS dan Linux) |
| `RemoveDeliveredNotification(identifier string)` | Menghapus notifikasi tertentu yang telah ditampilkan (khusus macOS dan Linux) |
| `RemoveNotification(identifier string)` | Menghapus notifikasi (khusus Linux) |

### Penanganan Peristiwa

| Metode | Deskripsi |
| --- | --- |
| `OnNotificationResponse(callback func(result NotificationResult))` | Mendaftarkan fungsi panggil balik untuk respons notifikasi |

### Struct dan Tipe

#### NotificationOptions

```go
type NotificationOptions struct {
    ID         string                 // Unique identifier for the notification (required)
    Title      string                 // Main notification title (required)
    Subtitle   string                 // Subtitle text (macOS and Linux only)
    Body       string                 // Main notification content
    CategoryID string                 // Category identifier for interactive notifications
    Data       map[string]interface{} // Custom data to associate with the notification

    // Sound controls the sound played on delivery. nil = platform default.
    Sound *NotificationSound

    // Attachments are media files shown alongside the notification.
    // macOS supports multiple attachments of any media type;
    // Windows and Linux honour the first image-typed attachment.
    Attachments []NotificationAttachment

    // ThreadID groups related notifications together in
    // Notification Center / Action Center / the Linux notification daemon.
    ThreadID string

    // InterruptionLevel controls priority. One of "passive",
    // "active" (default), "timeSensitive", "critical".
    InterruptionLevel string

    // Schedule defers delivery. macOS uses a native trigger that survives
    // restarts; Windows and Linux use an in-process timer that does NOT.
    Schedule *NotificationSchedule
}
```

#### NotificationSound

```go
type NotificationSound struct {
    Silent bool   // If true, no sound is played
    Name   string // Named/bundled sound (see "Custom Sound" above)
}
```

#### NotificationAttachment

```go
type NotificationAttachment struct {
    ID   string // Optional identifier
    Path string // Absolute filesystem path (macOS also accepts file:// URLs)
    // Type is an optional placement/UTI hint:
    //   macOS:   UTI like "public.png" / "public.audio" (often inferred)
    //   Windows: "hero" | "appLogoOverride" | "inline" (default "inline")
    //   Linux:   ignored (always image-path hint)
    Type string
}
```

#### NotificationSchedule

```go
// Exactly one of DelaySeconds or At must be set. At is Unix seconds (UTC).
type NotificationSchedule struct {
    DelaySeconds int   // Seconds from now until delivery
    At           int64 // Absolute Unix timestamp (seconds, UTC)
}
```

#### Konstanta InterruptionLevel

```go
const (
    InterruptionLevelPassive       = "passive"
    InterruptionLevelActive        = "active" // default
    InterruptionLevelTimeSensitive = "timeSensitive"
    InterruptionLevelCritical      = "critical"
)
```

#### NotificationCategory

```go
type NotificationCategory struct {
    ID               string                // Unique identifier for the category
    Actions          []NotificationAction  // Button actions for the notification
    HasReplyField    bool                  // Whether to include a text input field
    ReplyPlaceholder string                // Placeholder text for the input field
    ReplyButtonTitle string                // Text for the reply button
}
```

#### NotificationAction

```go
type NotificationAction struct {
    ID          string  // Unique identifier for the action
    Title       string  // Button text
    Destructive bool    // Whether the action is destructive (macOS-specific)
}
```

#### NotificationResponse

```go
type NotificationResponse struct {
    ID               string                  // Notification identifier
    ActionIdentifier string                  // Action that was triggered
    CategoryID       string                  // Category of the notification
    Title            string                  // Title of the notification
    Subtitle         string                  // Subtitle of the notification
    Body             string                  // Body text of the notification
    UserText         string                  // Text entered by the user
    UserInfo         map[string]interface{}  // Custom data from the notification
}
```

#### NotificationResult

```go
type NotificationResult struct {
    Response NotificationResponse  // Response data
    Error    error                 // Any error that occurred
}
```
