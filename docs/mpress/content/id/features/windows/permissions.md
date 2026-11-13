---
title: "Izin"
description: "Mengontrol permintaan akses kamera, mikrofon, geolokasi, dan kapabilitas lain dari konten web"
slug: "features/windows/permissions"
sourcePath: "features/windows/permissions.md"
---

Konten web yang memanggil `navigator.mediaDevices.getUserMedia()`, Geolocation API, atau Notifications API memerlukan aplikasi host untuk mengizinkan atau menolak permintaan tersebut. Wails menyediakan peta `Permissions` lintas platform pada `WebviewWindowOptions` yang memungkinkan Anda mengaturnya secara deklaratif — tanpa memerlukan kode khusus platform.

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

Permintaan kamera dan mikrofon dari konten web jendela tersebut diizinkan tanpa prompt browser.

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
| `PermissionAllow` | 1 | Izinkan tanpa prompt |
| `PermissionDeny` | 2 | Tolak tanpa prompt |

`PermissionDefault` adalah nilai nol, sehingga entri yang tidak ditetapkan dalam peta berperilaku sesuai nilai default.

## Perilaku Platform

Setiap platform menangani `PermissionDefault` secara berbeda karena webview yang mendasarinya memiliki perilaku bawaan yang berbeda.

### Linux (WebKitGTK)

WebKitGTK **tidak memiliki prompt izin bawaan**. Tanpa handler yang terpasang, WebKitGTK menolak setiap permintaan secara diam-diam — karena itulah `getUserMedia` selalu mengembalikan `NotAllowedError` sebelum fitur ini ditambahkan.

Saat ini Wails menangani permintaan **kamera dan mikrofon** di Linux. Geolokasi, notifikasi, dan pembacaan papan klip belum terhubung dan tetap ditolak, apa pun kebijakan yang Anda tetapkan.

| Kebijakan | Kamera/Mikrofon | Geolokasi, Notifikasi, Papan Klip |
| --- | --- | --- |
| `PermissionDefault` | **Diizinkan** (memulihkan getUserMedia) | Selalu ditolak |
| `PermissionAllow` | Diizinkan | Selalu ditolak (belum diimplementasikan) |
| `PermissionDeny` | Ditolak | Selalu ditolak |

### Windows (WebView2)

WebView2 memiliki prompt izin bawaan dan API izin untuk setiap jenis. Kelima jenis kapabilitas didukung sepenuhnya.

| Kebijakan | Perilaku |
| --- | --- |
| `PermissionDefault` | WebView2 menampilkan prompt izin bawaan sistem operasi/browser |
| `PermissionAllow` | Diizinkan secara diam-diam |
| `PermissionDeny` | Ditolak secara diam-diam |

**Penting:** Sebelum fitur ini tersedia, Wails memanggil `SetGlobalPermission(Allow)` tanpa syarat — sehingga semua kapabilitas diizinkan secara diam-diam. Sekarang, jika ada entri apa pun dalam `Permissions`, izin menyeluruh tersebut **tidak** ditetapkan. Kapabilitas yang tidak ditetapkan akan menggunakan prompt bawaan WebView2, bukan diizinkan secara otomatis.

Artinya, jika Anda mengonfigurasi `Permissions` sama sekali di Windows, setiap kapabilitas yang tidak Anda cantumkan secara eksplisit akan menampilkan prompt, bukan diizinkan secara diam-diam. Tetapkan secara eksplisit kapabilitas yang Anda perlukan.

### macOS (WKWebView + TCC)

Di macOS, dua lapisan harus sama-sama menyetujui. WKWebView meminta izin kepada aplikasi sebelum memulai sesi pengambilan, dan permintaan itulah yang dijawab oleh peta `Permissions`. Di bawahnya, kerangka kerja privasi sistem (TCC) membatasi akses ke perangkat itu sendiri: prompt sistem muncul saat pertama kali aplikasi benar-benar mencoba mengakses kamera atau mikrofon, dan pilihan pengguna diingat untuk setiap aplikasi di Pengaturan Sistem → Privasi & Keamanan.

Wails menangani permintaan **kamera dan mikrofon** di macOS 12 dan versi yang lebih baru. Geolokasi, notifikasi, dan pembacaan papan klip tidak memiliki padanan `WKUIDelegate`, sehingga tidak terhubung dan hanya ditangani oleh TCC — seperti di Linux, kebijakan yang Anda tetapkan untuk kapabilitas tersebut tidak berpengaruh.

| Kebijakan | Kamera/Mikrofon | Geolokasi, Notifikasi, Papan Klip |
| --- | --- | --- |
| `PermissionDefault` | WebKit menampilkan prompt izinnya sendiri | Hanya TCC |
| `PermissionAllow` | Prompt WebKit dilewati — **TCC tetap berlaku** | Hanya TCC |
| `PermissionDeny` | Ditolak oleh WebKit; otorisasi sistem ditangani secara terpisah | Hanya TCC |

`PermissionAllow` mengizinkan permintaan webview, bukan akses perangkat. Otorisasi sistem tetap berlaku secara independen, dan aplikasi yang ditolak di Pengaturan Sistem akan tetap ditolak. Hanya prompt WebKit sendiri yang dilewati.

Pada versi sebelum macOS 12, metode delegate tersebut tidak tersedia, sehingga peta diabaikan dan setiap permintaan kembali menggunakan prompt WebKit.

Pastikan `Info.plist` Anda menyertakan kunci deskripsi penggunaan yang sesuai:

```xml
<key>NSMicrophoneUsageDescription</key>
<string>Used for voice input</string>
<key>NSCameraUsageDescription</key>
<string>Used for video calls</string>
```

@note{type="caution" title="Deskripsi penggunaan wajib disertakan"}

WebKit mungkin meminta otorisasi sistem sebelum memanggil delegate izin Wails. `PermissionDeny` memblokir permintaan WebKit, tetapi tidak menjamin bahwa prompt TCC atau permintaan otorisasi AVFoundation tidak akan muncul.

Sertakan deskripsi penggunaan yang sesuai untuk perangkat yang diminta aplikasi Anda, serta entitlement perangkat yang cocok jika aplikasi berjalan dalam sandbox. Jangan gunakan `PermissionDeny` sebagai pengganti deklarasi tersebut.

@end

## Pola Umum

### Aplikasi pengambilan media

Izinkan kamera dan mikrofon di semua platform:

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionAllow,
    application.PermissionCamera:     application.PermissionAllow,
},
```

Di **Linux**, pengaturan ini secara eksplisit mengizinkan kedua perangkat tersebut; kapabilitas lainnya tetap ditolak. Di **Windows**, pengaturan ini mengizinkan keduanya; kapabilitas lain yang tidak Anda cantumkan akan menampilkan dialog izin native. Di **macOS**, pengaturan ini mengizinkan keduanya pada lapisan WebKit sehingga dialog browser tidak muncul; TCC tetap meminta izin untuk perangkat itu sendiri saat pertama kali digunakan, dan kedua kunci deskripsi penggunaan harus tersedia di `Info.plist`.

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

## Penggantian khusus Windows

Kolom `Windows.Permissions` per jendela (`map[CoreWebView2PermissionKind]CoreWebView2PermissionState`) tetap berfungsi dan dapat mengganti pengaturan masing-masing kapabilitas setelah peta lintas platform diterapkan. Gunakan kolom ini jika Anda memerlukan akses ke jenis izin WebView2 yang tidak memiliki padanan lintas platform (misalnya, `CoreWebView2PermissionKindOtherSensors`).

```go
Windows: application.WindowsWindow{
    Permissions: map[application.CoreWebView2PermissionKind]application.CoreWebView2PermissionState{
        application.CoreWebView2PermissionKindOtherSensors: application.CoreWebView2PermissionStateAllow,
    },
},
```

Urutan evaluasi di Windows adalah:

1. Peta lintas platform `Permissions` (menetapkan status per jenis melalui `SetPermission`)
2. Peta `Windows.Permissions` (mengganti pengaturan masing-masing jenis)
3. Untuk jenis apa pun yang tidak tercakup oleh kedua peta: dialog izin native WebView2 (ketika kebijakan dikonfigurasi) atau izin otomatis (ketika tidak ada kebijakan yang dikonfigurasi—perilaku lama)

## Matriks dukungan platform

| Kapabilitas | Linux | Windows | macOS |
| --- | --- | --- | --- |
| Mikrofon | ✅ | ✅ | ✅ (macOS 12+) |
| Kamera | ✅ | ✅ | ✅ (macOS 12+) |
| Geolokasi | ❌ belum didukung | ✅ | ❌ belum didukung |
| Notifikasi | ❌ belum didukung | ✅ | ❌ belum didukung |
| Pembacaan papan klip | ❌ belum didukung | ✅ | ❌ belum didukung |

Jika macOS ditandai ✅, kebijakan tersebut menanggapi permintaan WebKit; selain itu, TCC tetap membatasi akses ke perangkat. Jika ditandai ❌, kapabilitas tersebut hanya dikendalikan oleh TCC.

## Pemecahan masalah

**`getUserMedia` masih gagal di Linux setelah pemutakhiran**

Pastikan Anda tidak menetapkan `PermissionMicrophone: PermissionDeny` atau `PermissionCamera: PermissionDeny` secara eksplisit. Pengaturan default (tidak ditetapkan) mengizinkan pengambilan media di Linux.

**Windows menampilkan dialog untuk izin yang tidak saya konfigurasi**

Setelah `Permissions` memiliki entri apa pun, Wails tidak lagi menetapkan izin menyeluruh `Allow`. Kapabilitas yang tidak Anda cantumkan akan menampilkan dialog izin native WebView2. Tambahkan entri `PermissionAllow` secara eksplisit untuk setiap kapabilitas yang digunakan aplikasi Anda.

**Izin macOS tidak berfungsi**

`Permissions` mencakup kamera dan mikrofon pada macOS 12 dan versi yang lebih baru; geolokasi, notifikasi, dan pembacaan papan klip belum terhubung serta mengabaikan kebijakan tersebut. Jika kebijakan diterapkan, TCC tetap membatasi akses ke perangkat: `PermissionAllow` menghilangkan dialog WebKit, tetapi tidak menghilangkan dialog sistem. Pastikan `Info.plist` Anda menyertakan kunci deskripsi penggunaan yang benar (`NSMicrophoneUsageDescription`, `NSCameraUsageDescription`) dan pengguna telah memberikan akses di Pengaturan Sistem → Privasi & Keamanan.

**Aplikasi macOS saya berhenti ketika konten web meminta akses ke kamera atau mikrofon**

Pastikan `Info.plist` menyertakan `NSCameraUsageDescription` atau `NSMicrophoneUsageDescription` untuk setiap perangkat yang diminta aplikasi Anda. Deskripsi yang tidak tersedia dapat menyebabkan permintaan ditolak; mengakses AVFoundation tanpa deskripsi yang diwajibkan dapat menghentikan aplikasi. `PermissionDeny` mengendalikan keputusan WebKit dan tidak menggantikan kunci-kunci tersebut.

**Geolokasi/notifikasi/papan klip tidak berpengaruh di Linux**

Saat ini, hanya kamera dan mikrofon yang ditangani di Linux. Dukungan untuk jenis kapabilitas lain belum diterapkan—kapabilitas tersebut tetap ditolak, apa pun kebijakan yang Anda tetapkan.
