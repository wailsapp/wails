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

### macOS (WKWebView + TCC)

Dua lapisan harus sepakat di macOS. WKWebView menanyakan aplikasi sebelum memulai sesi capture, dan permintaan itulah yang dijawab oleh map `Permissions`. Di bawahnya, framework privasi sistem (TCC) menjaga perangkatnya sendiri: prompt OS muncul saat aplikasi benar-benar mengakses kamera atau mikrofon untuk pertama kali, dan pilihan pengguna diingat per-aplikasi di System Settings → Privacy & Security.

Wails menangani permintaan **kamera dan mikrofon** di macOS 12 dan yang lebih baru. Geolokasi, notifikasi, dan clipboard read tidak memiliki padanan `WKUIDelegate`, sehingga belum dihubungkan dan diserahkan sepenuhnya ke TCC — seperti di Linux, kebijakan yang Anda setel untuk ketiganya tidak berpengaruh.

| Kebijakan | Kamera / Mikrofon | Geolokasi, Notifikasi, Clipboard |
| --- | --- | --- |
| `PermissionDefault` | WebKit menampilkan prompt permission-nya sendiri | TCC saja |
| `PermissionAllow` | Prompt WebKit dilewati — **TCC tetap berlaku** | TCC saja |
| `PermissionDeny` | Ditolak sebelum perangkat disentuh | TCC saja |

`PermissionAllow` mengizinkan permintaan webview, bukan perangkatnya. Capture pertama tetap memunculkan prompt TCC, dan aplikasi yang ditolak pengguna di System Settings tetap ditolak — tidak ada aplikasi yang dapat memberikan akses perangkat kepada dirinya sendiri. Yang dihilangkan `PermissionAllow` adalah prompt WebKit di depannya.

Di bawah macOS 12 metode delegate tersebut tidak ada, sehingga map diabaikan di sana dan setiap permintaan kembali ke prompt WebKit.

Pastikan `Info.plist` Anda menyertakan kunci deskripsi penggunaan yang sesuai:

```xml
<key>NSMicrophoneUsageDescription</key>
<string>Used for voice input</string>
<key>NSCameraUsageDescription</key>
<string>Used for video calls</string>
```

@note{type="caution" title="Tolak apa yang tidak dideklarasikan Info.plist Anda"}

Kemampuan yang mencapai AVFoundation tanpa kunci deskripsi penggunaannya tidak gagal — macOS menghentikan aplikasi.

`PermissionDefault` adalah nilai nol, jadi `{PermissionMicrophone: PermissionAllow}` saja membiarkan kamera pada prompt WebKit. Jika pengguna menerima prompt itu dan aplikasi hanya mendeklarasikan `NSMicrophoneUsageDescription`, aplikasi akan mati. Setel `PermissionDeny` secara eksplisit untuk setiap kemampuan yang tidak Anda sertakan deskripsi penggunaannya:

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionAllow,
    application.PermissionCamera:     application.PermissionDeny,
},
```

@end

## Pola Umum

### Aplikasi perekaman media

Izinkan akses kamera dan mikrofon di semua platform:

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionAllow,
    application.PermissionCamera:     application.PermissionAllow,
},
```

Di **Linux**, pengaturan ini secara eksplisit mengizinkan kedua perangkat; kapabilitas lainnya tetap ditolak. Di **Windows**, pengaturan ini mengizinkan keduanya; kapabilitas lain yang tidak Anda cantumkan akan menampilkan dialog izin native. Di **macOS** ini mengizinkan keduanya di lapisan WebKit, sehingga tidak ada prompt browser yang muncul; TCC tetap menanyakan perangkatnya sendiri saat pertama kali digunakan, dan kedua kunci deskripsi penggunaan harus ada di `Info.plist`.

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
| Mikrofon | ✅ | ✅ | ✅ (macOS 12+) |
| Kamera | ✅ | ✅ | ✅ (macOS 12+) |
| Geolokasi | ❌ belum didukung | ✅ | ❌ belum didukung |
| Notifikasi | ❌ belum didukung | ✅ | ❌ belum didukung |
| Pembacaan Papan Klip | ❌ belum didukung | ✅ | ❌ belum didukung |

Di tempat macOS bertanda ✅, kebijakan menjawab permintaan WebKit; TCC tetap menjaga perangkat di atasnya. Di tempat bertanda ❌, kemampuan diserahkan sepenuhnya ke TCC.

## Pemecahan Masalah

**`getUserMedia` masih gagal di Linux setelah pemutakhiran**

Pastikan Anda tidak menetapkan `PermissionMicrophone: PermissionDeny` atau `PermissionCamera: PermissionDeny` secara eksplisit. Pengaturan default (tidak ditetapkan) mengizinkan pengambilan media di Linux.

**Windows menampilkan dialog untuk izin yang tidak saya konfigurasi**

Setelah ada entri apa pun di `Permissions`, Wails tidak lagi menetapkan pemberian izin menyeluruh `Allow`. Kapabilitas yang tidak Anda cantumkan akan menampilkan dialog izin native WebView2. Tambahkan entri `PermissionAllow` secara eksplisit untuk setiap kapabilitas yang digunakan aplikasi Anda.

**Izin macOS tidak berfungsi**

`Permissions` mencakup kamera dan mikrofon di macOS 12 dan yang lebih baru; geolokasi, notifikasi, dan clipboard read belum dihubungkan dan mengabaikan kebijakan. Di tempat kebijakan dihormati, TCC tetap menjaga perangkat di atasnya: `PermissionAllow` menghilangkan prompt WebKit, bukan prompt sistem. Pastikan `Info.plist` Anda menyertakan kunci deskripsi penggunaan yang benar (`NSMicrophoneUsageDescription`, `NSCameraUsageDescription`) dan pengguna telah memberikan akses di System Settings → Privacy & Security.

**Aplikasi macOS saya keluar saat konten web meminta kamera atau mikrofon**

macOS menghentikan aplikasi yang mencapai AVFoundation tanpa kunci deskripsi penggunaan yang sesuai. Tambahkan kuncinya, atau setel `PermissionDeny` untuk kemampuan tersebut agar permintaannya tidak pernah sampai sejauh itu — `PermissionDefault` membiarkannya pada prompt WebKit, yang dapat diterima pengguna.

**Geolokasi/notifikasi/papan klip tidak berpengaruh di Linux**

Saat ini, hanya kamera dan mikrofon yang ditangani di Linux. Dukungan untuk jenis kapabilitas lainnya belum diimplementasikan—kapabilitas tersebut tetap ditolak, apa pun kebijakan yang Anda tetapkan.
