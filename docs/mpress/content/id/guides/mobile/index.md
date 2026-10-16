---
title: "Ikhtisar Seluler"
description: "Bangun aplikasi iOS dan Android dari basis kode Go yang sama dengan aplikasi desktop Anda"
slug: "guides/mobile"
sourcePath: "guides/mobile/index.md"
---

Wails v3 berjalan di **iOS dan Android** menggunakan `main.go` dan frontend yang sama dengan yang sudah Anda tulis untuk desktop. Tidak ada proyek seluler terpisah, jembatan berbagi kode, atau penulisan ulang: biner Go dikompilasi untuk target seluler dan WebView native merender frontend Anda yang sudah ada.

@cards{cols="2"}
iOS
Host WKWebView + UIKit. Aset disajikan melalui skema `wails://` khusus — tanpa port terbuka. Memerlukan **macOS** dengan Xcode lengkap.

[Panduan iOS →](/guides/mobile/ios/)

---
Android
Android WebView + `WebViewAssetLoader`. Go dikompilasi sebagai `libwails.so` melalui NDK. Berfungsi di macOS, Linux, dan Windows.

[Panduan Android →](/guides/mobile/android/)

@end

## Lihat cara kerjanya: contoh Kitchen Sink

Cara terbaik untuk memahami kemungkinannya adalah dengan melihat **Kitchen Sink** — satu aplikasi Wails yang berjalan secara identik di iOS, Android, dan desktop dari satu basis kode:

@linkcard{title="Mobile Kitchen Sink — GitHub" href="https://github.com/wailsapp/wails/tree/master/v3/examples/mobile" description="Binding · Peristiwa · Dialog · Haptik · Geolokasi · Biometrik · Notifikasi · Penyimpanan aman · dan lainnya — semuanya dari satu main.go"}
Contoh ini mendemonstrasikan setiap antarmuka API seluler utama dalam 7 tab — dan juga berjalan di desktop. Tab **Seluler** dan **Perangkat Keras** disembunyikan di desktop melalui pemeriksaan platform di frontend; sisi Go tidak mendaftarkan handler untuk peristiwa seluler `common:*` saat dibuat untuk desktop. Ini adalah pola yang direkomendasikan untuk mendistribusikan satu basis kode ke semua platform.

| Tab | Platform | Yang ditampilkan |
| --- | --- | --- |
| **Binding** | semua | Panggilan layanan JS → Go yang mengembalikan nilai, struct, dan kesalahan |
| **Peristiwa** | semua | Jam Go → JS, ping/pong JS → Go → JS, peristiwa sistem OS (baterai, jaringan, tema) |
| **Dialog** | semua | Dialog pesan native di setiap platform |
| **Sistem** | semua | Papan klip, metrik layar, informasi perangkat |
| **Seluler** | iOS + Android | Lembar berbagi, menjaga perangkat tetap aktif, senter, kecerahan, biometrik, notifikasi lokal, penyimpanan aman |
| **Perangkat Keras** | iOS + Android | Haptik, geolokasi, akselerometer, sensor kedekatan, text-to-speech |
| **Native** | iOS + Android | iOS: haptik + pengalih WKWebView · Android: getaran + toast |

Untuk menjalankannya sendiri:

```bash
git clone https://github.com/wailsapp/wails.git
cd wails/v3/examples/mobile

wails3 task ios:run        # iOS Simulator (macOS + Xcode required)
wails3 task android:run    # Android Emulator
wails3 task run            # Desktop
```

## Cara kerjanya

Model aplikasi yang sama berlaku di setiap platform:

1. **Backend Go** — layanan, handler peristiwa, dan logika aplikasi Anda dikompilasi tanpa perubahan untuk `GOOS=ios` dan `GOOS=android`.
2. **Frontend** — HTML/JS/CSS yang sama persis. Paket `@wailsio/runtime` berfungsi secara identik; binding layanan, peristiwa, dialog, dan papan klip semuanya dirutekan melalui transport dalam proses yang sama.
3. **Host WebView** — di iOS, `WKWebView` di dalam `UIViewController`; di Android, `WebView` di dalam `Activity`. Wails menyiapkan jembatan pesan secara otomatis.
4. **Penyajian aset dalam proses** — aset disajikan langsung dari memori Go, bukan dari server localhost. Tanpa port terbuka, tanpa loopback, tanpa latensi tambahan.

Perilaku khusus platform ditempatkan dalam file yang dilindungi oleh `//go:build ios` atau `//go:build android`, sehingga kode bersama Anda tetap bersih.

## Ringkasan prasyarat

| Persyaratan | iOS | Android |
| --- | --- | --- |
| Sistem operasi | Hanya macOS | macOS, Linux, Windows |
| Toolchain | Xcode lengkap (bukan hanya alat CLI) | Android SDK + NDK 26.3.x + JDK |
| Go | 1.25+ | 1.25+ |
| npm | ✅ | ✅ |
| Verifikasi dengan | `wails3 doctor` | `wails3 doctor` |

@note{type="tip"}
Jalankan `wails3 doctor` setelah menyiapkan toolchain Anda — perintah ini menunjukkan secara tepat apa yang ditemukan dan apa yang belum tersedia untuk setiap platform.

@end

## Yang didukung

Kedua platform memiliki rangkaian fitur inti yang sama:

| Fitur | iOS | Android |
| --- | --- | --- |
| Binding layanan (JS → Go) | ✅ | ✅ |
| Peristiwa (dua arah) | ✅ | ✅ |
| Dialog pesan | ✅ UIAlertController | ✅ AlertDialog |
| Dialog buka file | ✅ UIDocumentPicker | ✅ Storage Access Framework |
| Dialog simpan file | ❌ sebagai gantinya, tulis ke sandbox | ❌ sebagai gantinya, tulis ke sandbox |
| Papan klip | ✅ UIPasteboard | ✅ ClipboardManager |
| Metrik layar / area aman | ✅ | ✅ |
| Peristiwa siklus hidup | ✅ `events.IOS.*` | ✅ `events.Android.*` |
| Umpan balik haptik | ✅ `IOS.Haptics.*` | ✅ `Android.Haptics.Vibrate` |
| Informasi perangkat | ✅ `IOS.Device.Info()` | ✅ `Android.Device.Info()` |
| Tab native (iOS) | ✅ UITabBar | — |
| Pesan toast (Android) | — | ✅ `Android.Toast.Show` |
| Beberapa jendela | ❌ hanya jendela pertama | ❌ hanya jendela pertama |
| Geometri jendela / menu / baki sistem | operasi kosong yang disengaja | operasi kosong yang disengaja |

## Aturan tag build

Dua aturan penting yang perlu diketahui saat menulis kode kondisional berdasarkan platform:

- **`ios` mengimplikasikan `darwin`** — file yang diberi tag `//go:build darwin` juga akan dikompilasi untuk iOS. Untuk menargetkan macOS saja, gunakan `//go:build darwin && !ios`.
- **`android` mengimplikasikan `linux`** — file yang diberi tag `//go:build linux` juga akan dikompilasi untuk Android. Untuk menargetkan Linux desktop saja, gunakan `//go:build linux && !android`.

Saat runtime, `runtime.GOOS` masing-masing mengembalikan `"ios"` dan `"android"`.

## Deteksi platform saat runtime

Tag build digunakan untuk kode yang hanya dapat *dikompilasi* pada platform tertentu. Untuk percabangan biasa dalam kode bersama, gunakan `application.System` — fitur ini tersedia di setiap build (tanpa memerlukan tag build), sehingga file yang sama dapat digunakan di semua platform:

```go
import "github.com/wailsapp/wails/v3/pkg/application"

if application.System.IsMobile() {
    // iOS or Android
} else if application.System.IsDesktop() {
    // macOS, Windows or Linux
}

// Or test a single target directly:
if application.System.IsPlatform(application.PlatformIOS) {
    // iOS only
}
```

Tersedia: `IsMobile()`, `IsDesktop()`, `IsServer()` (tag build `server`), dan `IsPlatform(application.PlatformMacOS | PlatformWindows | PlatformLinux | PlatformIOS | PlatformAndroid | PlatformServer)`.

Frontend menyediakan fungsi pembantu yang sesuai di `@wailsio/runtime`:

```js
import { System } from "@wailsio/runtime";

if (System.IsMobile()) { /* iOS or Android */ }
if (System.IsIOS()) { /* … */ }      // also IsAndroid, IsMac, IsWindows, IsLinux, IsDesktop
```

## Langkah berikutnya

@cards{cols="2"}
🚀 Aplikasi Seluler Pertama Anda
Jalankan aplikasi Wails desktop di iOS Simulator atau Android Emulator hanya dalam hitungan menit.

[Mulai →](/guides/mobile/first-mobile-app/)

---
Panduan iOS
Penyiapan lengkap toolchain iOS, simulator, build perangkat, penandatanganan, konfigurasi, dan referensi API.

[Panduan iOS →](/guides/mobile/ios/)

---
Panduan Android
Penyiapan lengkap Android SDK/NDK, emulator, penandatanganan APK, pengemasan untuk Play Store, dan referensi API.

[Panduan Android →](/guides/mobile/android/)

@end
