---
title: "API privat macOS"
description: "Semua fitur dan opsi Wails yang bergantung pada API privat macOS, beserta perintah untuk mengaktifkannya dan alternatif untuk build publik."
slug: "guides/build/private-macos-apis"
sourcePath: "guides/build/private-macos-apis.md"
---

Secara default, Wails menggunakan API publik macOS. Satu tag build Go `private_mac_apis` mengaktifkan pemanggilan WebKit dan AppKit privat yang tercantum di halaman ini. Semua opsi dan metode Go publik tetap tersedia pada kedua build. Tanpa tag tersebut, operasi yang hanya tersedia melalui API privat tidak melakukan apa pun; fitur yang memiliki alternatif publik akan menggunakan alternatif tersebut.

@note{type="caution" title="Aktifkan perilaku privat macOS"}
Menetapkan opsi jendela tidak mengaktifkan API privat. Tambahkan `private_mac_apis` ke perintah build untuk mengaktifkannya. Tag ini hanya berlaku untuk build desktop macOS, bukan build iOS, Android, Windows, Linux, atau server.

@end

## Aktifkan API privat

```bash
# Build an application
wails3 build -tags private_mac_apis

# Run with live reload
EXTRA_TAGS=private_mac_apis wails3 dev

# Package an application
wails3 package GOOS=darwin EXTRA_TAGS=private_mac_apis

# Run a Go-only example from its directory
go run -tags private_mac_apis .
```

Untuk build produksi langsung, gunakan `go build -tags production,private_mac_apis .`. Untuk contoh frontend, ikuti README masing-masing guna membangun binding dan aset sebelum menjalankannya. Taskfile khusus atau versi lama harus meneruskan `EXTRA_TAGS` ke compiler Go.

## Inventaris fitur

| Fitur atau nilai | Yang diaktifkan oleh `private_mac_apis` | Tanpa tag tersebut |
| --- | --- | --- |
| `Mac.Backdrop: MacBackdropTransparent` | WKWebView transparan di atas jendela native | Jendela native dikonfigurasi, tetapi webview tetap opak |
| `Mac.Backdrop: MacBackdropTranslucent` | WKWebView transparan agar efek blur native terlihat melaluinya | Efek blur dikonfigurasi di belakang webview yang opak |
| `Mac.Backdrop: MacBackdropLiquidGlass` | WKWebView transparan di atas lapisan kaca, dengan kontrol latar belakang webview privat | Lapisan kaca dikonfigurasi di belakang webview yang opak; penataan menggunakan alternatif publik |
| Pengosongan latar belakang webview selama penyiapan Liquid Glass | Kontrol WebKit privat `backgroundColor` | `underPageBackgroundColor` publik pada macOS 12+, atau warna lapisan pada macOS versi lama; tidak membuat webview menjadi transparan |
| `app.Window.NewNotchWindow(...)` | Webview transparan di dalam panel notch berbentuk khusus | Panel tetap berfungsi, termasuk penempatan dan animasinya, tetapi webview-nya tetap opak |
| `Mac.LiquidGlass.Style` | Pemetaan gaya native Wails yang sudah ada, termasuk nilai gaya gelap yang tidak terdokumentasi | Menggunakan gaya regular/clear publik dan tampilan terang/gelap; lihat tabel nilai di bawah |
| `Mac.LiquidGlass.GroupID` | Meminta pengelompokan kaca privat untuk pengidentifikasi yang tidak kosong | Diabaikan; tidak ada pengelompokan yang diminta |
| `Mac.LiquidGlass.GroupSpacing` | Meminta jarak antarkelompok privat untuk nilai yang lebih besar dari nol | Diabaikan |
| `window.OpenDevTools()` dan JavaScript `Window.OpenDevTools()` | Membuka inspector WebKit secara terprogram pada macOS 12+ | Tidak melakukan apa pun |
| `WebviewWindowOptions.OpenInspectorOnStartup: true` | Meminta pembukaan inspector secara terprogram saat jendela pertama kali ditampilkan, pada macOS 12+ | Tidak melakukan apa pun |
| Pengaktifan inspector lama, sebelum macOS 13.3 | Mengaktifkan fitur tambahan developer WebKit jika dukungan inspector disertakan saat kompilasi | Tidak melakukan apa pun; inspeksi publik melalui Safari memerlukan macOS 13.3+ |

## Transparansi dan latar belakang webview

**Memerlukan API privat:** transparansi webview yang digunakan oleh `MacBackdropTransparent`, `MacBackdropTranslucent`, `MacBackdropLiquidGlass`, dan jendela notch. Secara internal, Wails menetapkan kunci privat WebKit `drawsBackground`. Latar belakang HTML atau CSS yang transparan saja tidak dapat membuat WKWebView native yang opak menjadi transparan.

```go
Mac: application.MacWindow{
    // Requires -tags private_mac_apis for the blur to show through the webview.
    Backdrop: application.MacBackdropTranslucent,
},
```

Operasi warna latar belakang webview privat menggunakan kunci WebKit `backgroundColor`. Penyiapan Liquid Glass menggunakannya untuk mengosongkan latar belakang webview. Tanpa tag tersebut, operasi internal ini menggunakan `underPageBackgroundColor` publik pada macOS 12+, atau lapisan view pada macOS versi lama. Alternatif tersebut tidak membuat webview menjadi transparan.

`WebviewWindowOptions.BackgroundColour` dan `window.SetBackgroundColour()` menetapkan warna **jendela native** pada macOS dan tidak memerlukan API privat. Demikian pula, `Frameless` dan `Mac.TitleBar.AppearsTransparent` menggunakan API publik AppKit; dependensi privatnya adalah transparansi webview, bukan transparansi bilah judul. Untuk efek latar belakang macOS, konfigurasikan `Mac.Backdrop` alih-alih hanya mengandalkan `BackgroundType`.

Lihat [opsi jendela](/features/windows/options/#mac-options), [jendela tanpa bingkai](/features/windows/frameless/#with-transparent-background), dan [jendela notch](/features/windows/notch-windows/).

## Nilai Liquid Glass

Jika `NSGlassEffectView` native tersedia (macOS 26+), pemetaan berikut berlaku. Hanya nilai gaya native `0` (regular) dan `1` (clear) yang terdokumentasi. Konstanta Go mempertahankan nilai yang sudah ada pada kedua build.

| Nilai `MacLiquidGlassStyle` | Dengan `private_mac_apis` | Tanpa tag tersebut |
| --- | --- | --- |
| `LiquidGlassStyleAutomatic` (`0`) | Gaya regular native (`0`) | Gaya regular native (`0`) |
| `LiquidGlassStyleLight` (`1`) | Pemetaan gaya native yang sudah ada (`1`, clear) | Gaya regular native (`0`) dengan tampilan Aqua |
| `LiquidGlassStyleDark` (`2`) | **Nilai gaya native yang tidak terdokumentasi `2`** | Gaya reguler native (`0`) dengan tampilan Dark Aqua |
| `LiquidGlassStyleVibrant` (`3`) | Dipetakan ke gaya terang/transparan native yang sudah ada (`1`) | Gaya transparan native (`1`) |

Automatic dan Vibrant menggunakan nilai gaya native yang terdokumentasi, tetapi latar belakang Liquid Glass untuk seluruh jendela tetap memerlukan tag agar **webview transparan**. Tampilan Light berbeda di antara kedua build. Pemetaan gaya privat tidak menjamin bahwa versi macOS mendatang akan merender efek yang sama.

**Selalu privat:** `GroupID` dan `GroupSpacing`. Wails memeriksa keberadaan selektor privat `setGroupIdentifier:`, `setGroupName:`, dan `setGroupSpacing:` sebelum meminta pengelompokan. Tanpa tag tersebut, operasi ini tidak melakukan apa pun. Mengaktifkan tag tidak menjamin bahwa versi macOS yang sedang berjalan mendukung selektor ini.

`MacLiquidGlass.Material`, `CornerRadius`, dan `TintColor` tidak dengan sendirinya memerlukan API privat. Pada versi macOS tanpa dukungan native untuk Liquid Glass, Wails mengonfigurasi fallback translusen; agar fallback tersebut terlihat melalui webview, tag tetap diperlukan.

## Web Inspector

**Memerlukan API privat:** memanggil `OpenDevTools()` atau menetapkan `OpenInspectorOnStartup: true` untuk membuka inspector WebKit dari aplikasi. Wails menggunakan selektor privat `_inspector`. Tanpa `private_mac_apis`, operasi ini diam-diam tidak melakukan apa pun.

Dukungan inspector juga harus disertakan saat kompilasi. Tag `production` dan `devtools` yang sudah ada tetap mempertahankan artinya:

| Tag build | Pembukaan inspector secara terprogram | Inspeksi publik Safari di macOS 13.3+ |
| --- | --- | --- |
| Tidak ada | Tidak melakukan apa pun | Diaktifkan |
| `private_mac_apis` | Diaktifkan pada macOS 12+ | Diaktifkan |
| `production` | Tidak melakukan apa pun | Dinonaktifkan |
| `production,private_mac_apis` | Tidak melakukan apa pun | Dinonaktifkan |
| `production,devtools` | Tidak melakukan apa pun | Diaktifkan |
| `production,devtools,private_mac_apis` | Diaktifkan pada macOS 12+ | Diaktifkan |

Pada macOS 13.3+, Wails mengaktifkan inspeksi Safari menggunakan `WKWebView.inspectable` publik; hal tersebut tidak memerlukan API privat. Sebelum macOS 13.3, fallback menggunakan preferensi privat `developerExtrasEnabled` sehingga memerlukan `private_mac_apis` serta dukungan inspector yang disertakan saat kompilasi.

```bash
# Production build with programmatic inspector support
go build -tags production,devtools,private_mac_apis .
```

## Memelihara daftar ini

Semua pemanggilan native privat diisolasi dalam `v3/pkg/application/mac_private_api_darwin.go`; build default memilih `mac_public_api_darwin.go`. Inventaris ini mencakup transparansi, warna latar belakang webview, gaya glass, pengelompokan glass, pembukaan inspector, dan pengaktifan inspector lama. Perubahan pada implementasi tersebut sebaiknya disertai pembaruan halaman ini dan dokumentasi opsi yang terdampak secara bersamaan.
