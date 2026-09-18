---
title: "Izin"
description: "Mengontrol permintaan akses kamera, mikrofon, geolokasi, dan kapabilitas lain dari konten web"
slug: "features/windows/permissions"
sourcePath: "features/windows/permissions.md"
---

Konten web yang memanggil `navigator.mediaDevices.getUserMedia()`, API Geolocation, atau API Notifications memerlukan aplikasi host untuk mengizinkan atau menolak permintaan tersebut. Wails menyediakan peta `Permissions` lintas platform pada `WebviewWindowOptions` yang memungkinkan Anda mengaturnya secara deklaratif—tanpa memerlukan kode khusus platform.

## Mulai Cepat

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "My App",
    Permissions: map[application.PermissionType]application.Permission{
        application.PermissionMicrophone: application.PermissionAllow,
        application.PermissionCamera:     application.PermissionAllow,
    },
})
```

Permintaan akses kamera dan mikrofon dari konten web jendela tersebut diizinkan tanpa menampilkan permintaan izin browser.

## Jenis Izin

`PermissionType` (uint8) mengidentifikasi kapabilitas yang dapat diminta oleh konten web.

| Konstanta | Kapabilitas |
| --- | --- |
| `PermissionMicrophone` | `getUserMedia({audio: true})` |
| `PermissionCamera` | `getUserMedia({video: true})` |
| `PermissionGeolocation` | `navigator.geolocation` |
| `PermissionNotifications` | `Notification.requestPermission()` |
| `PermissionClipboardRead` | `navigator.clipboard.readText()` |

## Nilai Izin

`Permission` (uint8) adalah kebijakan yang diterapkan pada jenis tertentu.

| Konstanta | Nilai | Arti |
| --- | --- | --- |
| `PermissionDefault` | 0 | Gunakan penanganan bawaan platform (lihat di bawah) |
| `PermissionAllow` | 1 | Izinkan tanpa meminta konfirmasi |
| `PermissionDeny` | 2 | Tolak tanpa meminta konfirmasi |

`PermissionDefault` adalah nilai nol, sehingga entri yang tidak ditetapkan dalam peta menggunakan perilaku bawaan.

## Perilaku Platform

Setiap platform menangani `PermissionDefault` secara berbeda karena webview yang mendasarinya memiliki perilaku bawaan yang berbeda.

### Linux (WebKitGTK)

WebKitGTK **tidak memiliki permintaan izin bawaan**. Tanpa handler yang terpasang, WebKitGTK secara diam-diam menolak setiap permintaan—karena itulah `getUserMedia` selalu mengembalikan `NotAllowedError` sebelum fitur ini ditambahkan.

Saat ini, Wails menangani permintaan akses **kamera dan mikrofon** di Linux. Geolokasi, notifikasi, dan pembacaan papan klip belum dihubungkan dan tetap ditolak, apa pun kebijakan yang Anda tetapkan.

| Kebijakan | Kamera/Mikrofon | Geolokasi, Notifikasi, Papan Klip |
| --- | --- | --- |
| `PermissionDefault` | **Diizinkan** (memulihkan getUserMedia) | Selalu ditolak |
| `PermissionAllow` | Diizinkan | Selalu ditolak (belum diimplementasikan) |
| `PermissionDeny` | Ditolak | Selalu ditolak |

### Windows (WebView2)

WebView2 memiliki permintaan izin bawaan dan API izin untuk setiap jenis. Kelima jenis kapabilitas didukung sepenuhnya.

| Kebijakan | Perilaku |
| --- | --- |
| `PermissionDefault` | WebView2 menampilkan permintaan izin bawaan OS/browser |
| `PermissionAllow` | Diizinkan tanpa pemberitahuan |
| `PermissionDeny` | Ditolak tanpa pemberitahuan |

**Penting:** Sebelum fitur ini tersedia, Wails memanggil `SetGlobalPermission(Allow)` tanpa syarat—sehingga secara diam-diam mengizinkan semua kapabilitas. Sekarang, jika ada entri apa pun dalam `Permissions`, izin menyeluruh tersebut **tidak** ditetapkan. Kapabilitas yang tidak ditetapkan akan ditangani oleh permintaan izin bawaan WebView2 dan tidak lagi diizinkan secara otomatis.

Artinya, jika Anda mengonfigurasi `Permissions` di Windows, kapabilitas apa pun yang tidak Anda cantumkan secara eksplisit akan menampilkan permintaan izin, bukan diizinkan tanpa pemberitahuan. Tetapkan secara eksplisit kapabilitas yang Anda perlukan.

### macOS (TCC)

macOS mengelola akses kamera, mikrofon, geolokasi, dan notifikasi melalui kerangka kerja privasi sistemnya. Permintaan izin OS muncul secara otomatis saat konten web pertama kali meminta suatu kapabilitas, dan pilihan pengguna disimpan per aplikasi di System Settings → Privacy & Security.

Ini berfungsi dengan benar tanpa konfigurasi `Permissions` apa pun. Peta tersebut **saat ini diabaikan di macOS**—semua permintaan diproses melalui TCC, apa pun yang Anda tetapkan. Keterbatasan praktisnya adalah `PermissionDeny` tidak berpengaruh di macOS: Anda tidak dapat mencegah webview menggunakan kapabilitas yang sudah diizinkan TCC pada tingkat sistem.

Pastikan `Info.plist` Anda menyertakan kunci deskripsi penggunaan yang sesuai:

```xml
<key>NSMicrophoneUsageDescription</key>
<string>Used for voice input</string>
<key>NSCameraUsageDescription</key>
<string>Used for video calls</string>
```

## Pola Umum

### Aplikasi perekaman media

Izinkan akses kamera dan mikrofon di semua platform:

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionAllow,
    application.PermissionCamera:     application.PermissionAllow,
},
```

Di **Linux**, pengaturan ini secara eksplisit mengizinkan kedua perangkat; kapabilitas lainnya tetap ditolak. Di **Windows**, pengaturan ini mengizinkan keduanya; kapabilitas lain yang tidak Anda cantumkan akan menampilkan dialog izin native. Di **macOS**, pengaturan ini tidak berpengaruh; TCC menangani semuanya.

### Menolak pengambilan media di Linux

Secara default, Linux mengizinkan kamera dan mikrofon. Untuk menonaktifkannya:

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionDeny,
    application.PermissionCamera:     application.PermissionDeny,
},
```

### Mengizinkan semua kapabilitas di Windows

Untuk memberikan semua kapabilitas tanpa menampilkan dialog izin:

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone:    application.PermissionAllow,
    application.PermissionCamera:        application.PermissionAllow,
    application.PermissionGeolocation:   application.PermissionAllow,
    application.PermissionNotifications: application.PermissionAllow,
    application.PermissionClipboardRead: application.PermissionAllow,
},
```

### Kebijakan per jendela

Setiap jendela dapat memiliki kebijakan yang berbeda:

```go
// Main app window — full media access
mainWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "App",
    Permissions: map[application.PermissionType]application.Permission{
        application.PermissionMicrophone: application.PermissionAllow,
        application.PermissionCamera:     application.PermissionAllow,
    },
})

// Settings window — no special capabilities
settingsWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Settings",
    // No Permissions entry — uses platform defaults
})

// Embedded content window — deny media capture
embeddedWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Embedded",
    Permissions: map[application.PermissionType]application.Permission{
        application.PermissionMicrophone: application.PermissionDeny,
        application.PermissionCamera:     application.PermissionDeny,
    },
})
```

## Penimpaan Khusus Windows

Kolom `Windows.Permissions` per jendela (`map[CoreWebView2PermissionKind]CoreWebView2PermissionState`) tetap berfungsi dan dapat menimpa masing-masing kapabilitas setelah peta lintas platform diterapkan. Gunakan kolom ini jika Anda memerlukan akses ke jenis izin WebView2 yang tidak memiliki padanan lintas platform (misalnya, `CoreWebView2PermissionKindOtherSensors`).

```go
Windows: application.WindowsWindow{
    Permissions: map[application.CoreWebView2PermissionKind]application.CoreWebView2PermissionState{
        application.CoreWebView2PermissionKindOtherSensors: application.CoreWebView2PermissionStateAllow,
    },
},
```

Urutan evaluasi di Windows adalah:

1. Peta `Permissions` lintas platform (menetapkan status per jenis melalui `SetPermission`)
2. Peta `Windows.Permissions` (menimpa masing-masing jenis)
3. Untuk jenis apa pun yang tidak tercakup oleh kedua peta: dialog izin native WebView2 (jika kebijakan dikonfigurasi) atau izin otomatis (jika tidak ada kebijakan yang dikonfigurasi—perilaku lama)

## Matriks Dukungan Platform

| Kapabilitas | Linux | Windows | macOS |
| --- | --- | --- | --- |
| Mikrofon | ✅ | ✅ | Hanya TCC |
| Kamera | ✅ | ✅ | Hanya TCC |
| Geolokasi | ❌ belum didukung | ✅ | Hanya TCC |
| Notifikasi | ❌ belum didukung | ✅ | Hanya TCC |
| Pembacaan Papan Klip | ❌ belum didukung | ✅ | Hanya TCC |

## Pemecahan Masalah

**`getUserMedia` masih gagal di Linux setelah pemutakhiran**

Pastikan Anda tidak menetapkan `PermissionMicrophone: PermissionDeny` atau `PermissionCamera: PermissionDeny` secara eksplisit. Pengaturan default (tidak ditetapkan) mengizinkan pengambilan media di Linux.

**Windows menampilkan dialog untuk izin yang tidak saya konfigurasi**

Setelah ada entri apa pun di `Permissions`, Wails tidak lagi menetapkan pemberian izin menyeluruh `Allow`. Kapabilitas yang tidak Anda cantumkan akan menampilkan dialog izin native WebView2. Tambahkan entri `PermissionAllow` secara eksplisit untuk setiap kapabilitas yang digunakan aplikasi Anda.

**Izin macOS tidak berfungsi**

Peta `Permissions` tidak berpengaruh di macOS. Pastikan `Info.plist` Anda menyertakan kunci deskripsi penggunaan yang benar (`NSMicrophoneUsageDescription`, `NSCameraUsageDescription`, dan sebagainya) serta pengguna telah memberikan akses di Pengaturan Sistem → Privasi & Keamanan.

**Geolokasi/notifikasi/papan klip tidak berpengaruh di Linux**

Saat ini, hanya kamera dan mikrofon yang ditangani di Linux. Dukungan untuk jenis kapabilitas lainnya belum diimplementasikan—kapabilitas tersebut tetap ditolak, apa pun kebijakan yang Anda tetapkan.
