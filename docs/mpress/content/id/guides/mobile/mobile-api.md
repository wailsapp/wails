---
title: "API Seluler"
description: "Pengelola lintas platform application.Mobile — satu titik masuk dengan build guard untuk kapabilitas seluler native yang digunakan bersama oleh iOS dan Android"
slug: "guides/mobile/mobile-api"
sourcePath: "guides/mobile/mobile-api.md"
---

@note{type="caution" title="Fitur Eksperimental"}
Dukungan seluler bersifat eksperimental dan dapat berubah dalam rilis mendatang.

@end

Kapabilitas seluler native diekspos melalui dua cara:

- **Pengelola per platform** — `application.IOS` (dalam file `//go:build ios`) dan `application.Android` (dalam file `//go:build android`). Gunakan pengelola ini untuk semua hal yang khusus untuk platform tertentu. Lihat referensi [iOS](/guides/mobile/ios/) dan [Android](/guides/mobile/android/) untuk mengetahui seluruh API per platform.
- **`application.Mobile`** — satu pengelola dengan build guard yang mencakup subset kapabilitas yang berperilaku identik pada kedua platform. Gunakan ini jika Anda menginginkan satu jalur kode yang dapat dikompilasi dan dijalankan di semua platform.

## `application.Mobile`

`application.Mobile` meneruskan pemanggilan ke `IOS` di iOS, ke `Android` di Android, dan ke stub tanpa operasi di desktop. Karena tidak memiliki batasan build, Anda dapat memanggilnya dari kode Go biasa yang tidak bergantung pada platform — tanpa perlu membuat file `//go:build` sendiri:

```go
// Works in any file, on any target.
// On desktop this returns "" (no-op); on device it returns the real path.
dbDir := application.Mobile.StoragePath()
if dbDir == "" {
    // Off-device, or the directory could not be created — handle accordingly.
    return
}
db, _ := sql.Open("sqlite", filepath.Join(dbDir, "app.db"))
```

`StoragePath()` mengembalikan path absolut ke direktori file privat aplikasi — `getFilesDir()` di Android dan direktori Application Support di iOS — yang direkomendasikan sebagai lokasi penyimpanan database dan file persisten lainnya. Fungsi ini mengembalikan string kosong di desktop, serta di perangkat jika direktori tidak tersedia (di iOS, jika direktori tidak dapat dibuat), jadi periksa `""` sebelum menggunakannya.

@note{type="note"}
Di luar perangkat (build desktop), setiap metode `Mobile` tidak melakukan operasi dan setiap kueri mengembalikan nilai nolnya (`""` untuk string). Perilaku inilah yang memungkinkan kode lintas platform memanggil `application.Mobile.*` tanpa syarat. Jika Anda juga memerlukan path yang nyata di desktop, buat percabangan berdasarkan platform dan gunakan `os.UserConfigDir()` atau API serupa sebagai fallback.

@end

## Kapabilitas

Pengelola `Mobile` mengekspos kapabilitas yang memiliki signature identik di iOS dan Android:

| Kapabilitas | API | Catatan |
| --- | --- | --- |
| Lembar berbagi | `Mobile.Share(json)` | `{text, url}` |
| Buka URL secara eksternal | `Mobile.OpenURL(url)` | Browser sistem |
| Pertahankan layar tetap aktif | `Mobile.SetKeepAwake(bool)` |  |
| Lampu senter | `Mobile.SetTorch(bool)` | → `common:torch` |
| Inset area aman | `Mobile.SafeAreaJSON()` | `{top,bottom,left,right}` |
| Informasi aplikasi | `Mobile.AppInfoJSON()` | `{name,version,build,bundleId}` |
| Kunci orientasi | `Mobile.SetOrientation(mode)` | `portrait` / `landscape` / `auto` |
| Bilah status | `Mobile.SetStatusBar(json)` | gaya + visibilitas |
| Informasi penyimpanan | `Mobile.StorageJSON()` | `{free,total}` byte |
| Path penyimpanan | `Mobile.StoragePath()` | Direktori file privat aplikasi |
| Daya / baterai | `Mobile.PowerJSON()` | `{level,charging,lowPower}` |
| Status jaringan | `Mobile.NetworkJSON()` | `{connected,type}` |
| Biometrik | `Mobile.BiometricAuthenticate(reason)` | → `common:biometric` |
| Penyimpanan aman | `Mobile.SecureGet(key)` / `Mobile.SecureDelete(key)` | Keychain / `EncryptedSharedPreferences` |
| Geolokasi | `Mobile.GetLocation()` | sekali jalan → `common:location` |
| Haptik | `Mobile.Haptic(type)` | dampak / notifikasi / pilihan |
| Akselerometer | `Mobile.SetMotion(bool)` | → `common:motion` |
| Proksimitas | `Mobile.SetProximity(bool)` | → `common:proximity` |
| Teks ke ucapan | `Mobile.Speak(text)` / `Mobile.StopSpeak()` |  |
| Inset papan ketik | `Mobile.SetKeyboardWatch(bool)` | → `common:keyboard` |
| Penangkapan layar | `Mobile.SetScreenProtect(bool)` | → `common:screenCapture` |
| Kamera | `Mobile.CapturePhoto()` / `Mobile.CaptureVideo()` | → `common:capture` |

Hasil asinkron diterima sebagai peristiwa `common:*`, sama persis seperti pada pengelola khusus tiap platform — lihat [Peristiwa](/guides/mobile/ios/#events) untuk payload-nya.

## Yang tetap khusus untuk platform

Kapabilitas yang bentuknya berbeda antara iOS dan Android **tidak** tersedia di `Mobile`; panggil kapabilitas tersebut melalui `application.IOS` / `application.Android` dari file dengan tag build:

| Aspek | iOS | Android |
| --- | --- | --- |
| Kecerahan (atur) | `IOS.SetBrightness(0.0-1.0)` | `Android.SetBrightness(0-100)` |
| Kecerahan / orientasi (dapatkan) | `IOS.GetBrightness()` / `IOS.GetOrientation()` | `Android.BrightnessJSON()` / `Android.OrientationJSON()` |
| Notifikasi lokal | `IOS.PostNotification(json)` | `Android.Notify(json)` |
| Penyimpanan aman (tulis) | `IOS.SecureSet(key, value)` | `Android.SecureSet(json)` |
| Eksekusi di latar belakang | `IOS.BeginBackgroundTask` / `EndBackgroundTask` | `Android.StartForegroundService` / `StopForegroundService` |

@note{type="tip"}
Antarmuka `MobileManager` adalah kontrak yang mendasari `application.Mobile`. Karena kedua pengelola platform harus mematuhinya, setiap metode yang tercantum di atas dijamin mempertahankan signature yang identik di iOS dan Android — jika keduanya sampai menyimpang, build platform akan gagal dikompilasi.

@end
