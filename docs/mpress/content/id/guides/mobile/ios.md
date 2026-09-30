---
title: "iOS"
description: "Bangun dan jalankan aplikasi Wails di iOS — penyiapan, simulator, build perangkat, konfigurasi, dan fitur native"
slug: "guides/mobile/ios"
sourcePath: "guides/mobile/ios.md"
---

@note{type="caution" title="Fitur Eksperimental"}
Dukungan iOS bersifat eksperimental dan dapat berubah dalam rilis mendatang.

@end

@note{type="tip"}
Baru menggunakan Wails untuk perangkat seluler? Mulailah dengan [Aplikasi Seluler Pertama Anda →](/guides/mobile/first-mobile-app/) untuk mengikuti panduan langkah demi langkah, lalu kembali ke sini untuk membaca referensi lengkap.

@end

Aplikasi Wails v3 berjalan di iOS sebagai aplikasi yang sepenuhnya native — dan bagian terbaiknya, aplikasi tersebut bekerja *persis* seperti versi desktop. Backend Go yang sama, frontend yang sama, dan `@wailsio/runtime` yang sama: binding layanan, event, dialog, dan papan klip semuanya berperilaku identik, dengan **tanpa** pengaturan ulang khusus perangkat seluler. Tidak ada basis kode seluler terpisah, lapisan porting, atau API khusus yang harus dipelajari — aplikasi Wails Anda yang sudah ada langsung berjalan di iOS. Proses porting benar-benar mulus: bawa aplikasi Anda apa adanya lalu rilis.

`main.go` yang sama dibuat untuk desktop dan iOS; penyesuaian khusus iOS dikonfigurasi melalui `application.Options.IOS`.

## Persyaratan

- macOS dengan **Xcode lengkap** yang telah terinstal (alat baris perintah saja tidak cukup) — `wails3 doctor` menampilkan SDK iOS yang dapat ditemukannya
- Go 1.25+ dan npm

## Simulator

Dari direktori proyek Anda:

```bash
wails3 task ios:run
```

Perintah ini mem-build aplikasi Anda, menjalankan simulator jika belum ada yang berjalan, lalu meluncurkan aplikasi.

Perintah pendamping yang berguna:

```bash
wails3 task ios:logs:dev    # stream the app's logs from the simulator
wails3 task ios:xcode       # open the generated Xcode project
```

Dalam build debug, webview dapat diperiksa dari menu Develop di Safari.

## Pemaketan

```bash
wails3 task ios:package             # production .app for the simulator
wails3 task ios:deploy-simulator    # install + launch it
```

Ini adalah build produksi yang dioptimalkan dan dipangkas.

## Build perangkat

```bash
wails3 task ios:package IOS_PLATFORM=device \
    CODESIGN_IDENTITY="Apple Development: You (TEAMID)" \
    PROVISIONING_PROFILE=path/to/profile.mobileprovision

wails3 task ios:deploy-device [DEVICE_ID=<udid>]     # install + launch on a device
wails3 task ios:package:ipa IOS_PLATFORM=device ...  # distribution .ipa
```

`IOS_PLATFORM=device` membuat build untuk perangkat fisik. Entitlement berasal dari `build/ios/entitlements.plist` dan hanya berlaku untuk build perangkat — tambahkan kunci kapabilitas yang dibutuhkan aplikasi Anda.

@note{type="tip"}
Untuk penandatanganan, provisioning, dan arsip App Store yang dikelola secara otomatis, buka proyek Xcode yang dihasilkan dengan `wails3 task ios:xcode`, lalu lakukan build dari Xcode.

@end

## Konfigurasi

`build/config.yml`:

```yaml
ios:
  bundleID: com.example.myapp
  displayName: My App
  version: 1.0.0
  minIOSVersion: "15.0"
```

Opsi startup (`application.Options.IOS`) mencakup `DisableScroll`, `DisableBounce`, `DisableScrollIndicators`, `DisableInputAccessoryView`, `EnableBackForwardNavigationGestures`, `DisableLinkPreview`, `EnableInlineMediaPlayback`, `EnableAutoplayWithoutUserAction`, `DisableInspectable`, `UserAgent`, `ApplicationNameForUserAgent`, `BackgroundColour`, serta tab bawah native melalui `EnableNativeTabs` + `NativeTabsItems`.

## Fitur native

Kapabilitas khusus iOS tersedia melalui `application.IOS`, yang dipanggil dari Go di dalam file `//go:build ios` agar kode bersama Anda tetap tidak bergantung pada platform. Android menyediakan kumpulan kapabilitas yang sama melalui `application.Android`.

Tindakan sekali jalan segera mengembalikan kontrol:

```go
//go:build ios

application.IOS.Haptic("impact-medium") // impact-light|impact-medium|impact-heavy|success|warning|error|selection
application.IOS.Share(`{"text":"Hi","url":"https://wails.io"}`)
application.IOS.SetKeepAwake(true)
application.IOS.PostNotification(`{"title":"Done","body":"Build finished","delay":2}`)
application.IOS.SecureSet("token", "abc") // stored securely
```

Helper kueri mengembalikan hasilnya sebagai JSON — `SafeAreaJSON()`, `AppInfoJSON()`, `PowerJSON()`, `NetworkJSON()`, `StorageJSON()`, `GetOrientation()`, `GetBrightness()`. `StoragePath()` mengembalikan path absolut ke direktori Application Support aplikasi — tempat yang cocok untuk basis data dan file persisten lainnya (padanan `getFilesDir()` milik Android di iOS). Direktori tersebut dibuat saat pertama kali diakses; `StoragePath()` mengembalikan string kosong jika direktori tidak dapat dibuat, jadi periksa `""` sebelum menggunakannya.

### Event

Segala sesuatu yang selesai kemudian — prompt izin, stream sensor, pengambilan gambar kamera — mengirimkan hasilnya sebagai **event**, bukan nilai kembalian, dan Anda dapat memantaunya di Go atau frontend. Nama diawali dengan `common:` untuk kapabilitas yang juga tersedia di Android dan `ios:` untuk kapabilitas khusus iOS.

```go
// Go
app.Event.On("common:location", func(e *application.CustomEvent) {
    // e.Data -> {"lat":..,"lng":..,"accuracy":..} or {"error":..}
})
```

```js
// frontend
import { Events } from "@wailsio/runtime";
Events.On("common:notification", (e) => { /* {ok, scheduled, presented, tapped, error} */ });
```

| Event | Dipicu oleh | Payload |
| --- | --- | --- |
| `common:biometric` | `BiometricAuthenticate(reason)` | `{ok, error}` |
| `common:location` | `GetLocation()` | `{lat, lng, accuracy}` / `{error}` |
| `common:motion` | `SetMotion(true)` | `{x, y, z}` |
| `common:proximity` | `SetProximity(true)` | `{near}` |
| `common:keyboard` | `SetKeyboardWatch(true)` | `{visible, height}` |
| `common:torch` | `SetTorch(bool)` | `{on, available}` |
| `common:notification` | `PostNotification(json)` | `{ok, scheduled, presented, tapped, error}` |
| `common:capture` | `CapturePhoto()` / `CaptureVideo()` | `{type, path, size, thumb}` |
| `common:screenCapture` | `SetScreenProtect(true)` | `{screenshot, recording}` |
| `ios:backgroundTask` | `BeginBackgroundTask(seconds)` | `{message, granted}` |

Contoh lengkap di bawah `v3/examples/mobile` menghubungkan setiap fitur di atas dari awal hingga akhir.

## Kontrol WebView

Beberapa perilaku WebView juga dapat diubah saat runtime dari Go:

```go
application.IOS.SetScrollEnabled(false)
application.IOS.SetBounceEnabled(false)
application.IOS.SetScrollIndicatorsEnabled(false)
application.IOS.SetBackForwardGesturesEnabled(true)
application.IOS.SetLinkPreviewEnabled(false)
application.IOS.SetInspectableEnabled(true)
application.IOS.SetCustomUserAgent("MyApp/1.0")
```

`@wailsio/runtime` yang disertakan juga menyediakan namespace iOS **frontend** berukuran kecil:

```js
import { IOS } from "@wailsio/runtime";
await IOS.Haptics.Impact("medium"); // light|medium|heavy|soft|rigid
const info = await IOS.Device.Info();
```

Pilihan tab bawah native diterima sebagai peristiwa `nativeTabSelected` pada `window`.

## Status dukungan

| Area | Status |
| --- | --- |
| Rendering frontend & aset | ✅ |
| Binding layanan, peristiwa (dua arah) | ✅ |
| Dialog pesan | ✅ |
| Dialog untuk membuka file / beberapa file / direktori | ✅ Diimpor sebagai salinan dalam sandbox |
| Dialog penyimpanan file | ❌ Sebagai gantinya, tulis di dalam sandbox aplikasi |
| Papan klip | ✅ |
| API layar | ✅ Mencakup area kerja dalam area aman |
| Peristiwa siklus hidup | ✅ |
| Geometri jendela, menu, baki sistem | Tidak melakukan apa pun di iOS |
| Beberapa jendela | Hanya jendela pertama yang ditampilkan |

## Catatan porting

- Kode desktop dapat di-build untuk iOS tanpa perubahan — pemanggilan untuk jendela, menu, dan baki sistem tidak melakukan apa pun.
- Ganti dialog penyimpanan file dengan penulisan ke dalam sandbox aplikasi, lalu bagikan file tersebut.
- Rancang frontend agar responsif; area aman ditangani secara otomatis.
