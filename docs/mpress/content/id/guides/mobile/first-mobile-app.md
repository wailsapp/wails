---
title: "Aplikasi Seluler Pertama Anda"
description: "Jalankan aplikasi Wails Anda di iOS Simulator atau Android Emulator dalam hitungan menit"
slug: "guides/mobile/first-mobile-app"
sourcePath: "guides/mobile/first-mobile-app.md"
---

Panduan ini menjalankan aplikasi desktop Wails standar di iOS Simulator atau Android Emulator. **Anda tidak perlu mengubah kode Go.** `main.go` yang sama dapat dibangun untuk semua target.

**Waktu penyelesaian:** 15–30 menit (sebagian besar waktu digunakan untuk menginstal toolchain saat pertama kali dijalankan)

## Mulai dari proyek desktop

Jika belum memilikinya, buat proyek baru:

```bash
wails3 init -n mymobileapp
cd mymobileapp
```

Pastikan aplikasi desktop berfungsi terlebih dahulu:

```bash
wails3 dev
```

Setelah terbuka, tutup aplikasi dan lanjutkan. Semua yang berjalan di desktop juga berjalan di perangkat seluler—untuk panduan ini, Anda tidak perlu mengubah `main.go` ataupun kode Go apa pun.

---

## Pilih platform Anda

@tabs{sync-key="mobile-platform"}
[iOS Simulator]
### Persyaratan

- **macOS** (build iOS hanya dapat dilakukan di macOS)
- **Xcode lengkap**—bukan hanya alat baris perintah. Instal dari App Store, lalu:
  ```bash
  sudo xcode-select -s /Applications/Xcode.app/Contents/Developer
  sudo xcodebuild -license accept
  ```


- **Go 1.26+** dan **npm** (sudah terinstal jika Anda telah menjalankan `wails3 init`)

Jalankan `wails3 doctor` untuk memverifikasi—perintah ini mencantumkan SDK iOS yang dapat ditemukannya.

### Jalankan di Simulator

@steps
### Luncurkan aplikasi
```bash
wails3 task ios:run
```

Selesai—perintah ini membangun aplikasi Anda, menjalankan simulator jika belum ada yang berjalan, lalu meluncurkan aplikasi.

@note{type="tip"}
Proses pertama memerlukan waktu beberapa menit (Wails sedang mengompilasi dan menyimpan framework Wails untuk iOS ke cache). Proses berikutnya jauh lebih cepat.

@end

Setelah diluncurkan, aplikasi desktop Anda yang belum dimodifikasi berjalan di iOS Simulator—`main.go` yang sama, frontend yang sama:

![Aplikasi Wails bawaan yang berjalan di iOS Simulator](/assets/ios-simulator-first-app.png)

### Tampilkan log secara langsung
Di terminal terpisah:

```bash
wails3 task ios:logs:dev
```

Perintah ini mengikuti log simulator secara langsung dan memfilternya untuk aplikasi Anda. Output `fmt.Println` dan `log.Println` muncul di sini.

### Periksa WebView
Di Safari: **Develop → Simulator → aplikasi Anda**. Seluruh fitur Web Inspector dapat digunakan—konsol, debugger, panel jaringan, dan semuanya.

### Buat perubahan
Edit file frontend apa pun (`frontend/src/main.js`, `index.html`, dan sebagainya), lalu jalankan kembali `wails3 task ios:run`. Wails membangun ulang frontend dan meluncurkan ulang aplikasi.

Untuk perubahan Go, jalankan kembali `wails3 task ios:run` juga. Kompilasi ulang Go bersifat inkremental, sehingga hanya paket yang berubah yang dibangun ulang.

@end

### Buka di Xcode (opsional)

```bash
wails3 task ios:xcode
```

Perintah ini membuka `build/ios/` di Xcode. Anda dapat menggunakan Xcode untuk deployment ke perangkat, profiling tingkat lanjut, atau mengelola profil provisioning. Wails membuat ulang proyek Xcode pada setiap build, jadi jangan ubah file yang dihasilkan secara langsung.

[Android Emulator]
### Persyaratan

Anda memerlukan **Android SDK**, **NDK**, dan **JDK**. Cara termudah adalah menggunakan Android Studio atau alat baris perintah:

@steps
### Instal alat baris perintah Android
Unduh dari [developer.android.com/studio#command-line-tools-only](https://developer.android.com/studio#command-line-tools-only), lalu ekstrak ke `~/android-sdk/cmdline-tools/latest/`.

### Instal komponen SDK
```bash
sdkmanager "platform-tools" \
           "platforms;android-35" \
           "build-tools;35.0.0" \
           "ndk;26.3.11579264" \
           "emulator" \
           "system-images;android-35;google_apis;arm64-v8a"
```

### Buat emulator
```bash
avdmanager create avd \
  --name wails \
  --package "system-images;android-35;google_apis;arm64-v8a" \
  --device pixel_7
```

### Atur variabel lingkungan
Tambahkan ke `~/.zshrc` atau `~/.bashrc`:

```bash
export ANDROID_HOME=~/android-sdk
export ANDROID_SDK_ROOT=~/android-sdk
export PATH=$PATH:$ANDROID_HOME/platform-tools:$ANDROID_HOME/cmdline-tools/latest/bin
```

Muat ulang: `source ~/.zshrc`

### Instal JDK
```bash
# macOS
brew install openjdk@21
export JAVA_HOME=$(brew --prefix openjdk@21)

# Ubuntu/Debian
sudo apt install openjdk-21-jdk
export JAVA_HOME=/usr/lib/jvm/java-21-openjdk-amd64

# Windows (scoop)
scoop install openjdk21
```

@end

Jalankan `wails3 doctor` untuk memastikan semuanya ditemukan.

### Jalankan di Emulator

@steps
### Luncurkan aplikasi
```bash
wails3 task android:run
```

Saat pertama kali dijalankan, perintah ini:

- Menjalankan emulator jika belum ada yang berjalan
- Menghasilkan binding dan membangun frontend
- Mengompilasi kode Go Anda menjadi `libwails.so` melalui cross-compiler NDK
- Merakit APK debug dengan Gradle
- Menginstal dan meluncurkannya di emulator

@note{type="tip"}
Build pertama mengunduh Gradle dan mengompilasi toolchain NDK—proses ini memerlukan sekitar 5–10 menit. Build berikutnya bersifat inkremental dan selesai dalam waktu kurang dari satu menit.

@end

### Tampilkan log secara langsung
Di terminal terpisah:

```bash
wails3 task android:logs
```

Perintah ini menjalankan `adb logcat` yang difilter untuk aplikasi Anda. Output `fmt.Println` muncul di sini.

### Periksa WebView
Buka Chrome dan akses `chrome://inspect`. WebView aplikasi Anda muncul di bagian **Remote Target**—klik **inspect** untuk membuka DevTools.

### Buat perubahan
Edit file apa pun, lalu jalankan kembali `wails3 task android:run`. Karena build Gradle bersifat inkremental, hanya kode yang berubah yang dikompilasi ulang.

@end

@end

---

## Memahami proses yang terjadi

`main.go` Anda sama sekali tidak berubah. Wails menangani semuanya:

- **Sistem build**—`Taskfile.yml` dalam proyek Anda berisi task `ios:*` dan `android:*` yang menjalankan toolchain khusus platform.
- **Kompilasi silang Go**—`GOOS=ios` atau `GOOS=android` dengan `GOARCH` dan sysroot yang sesuai.
- **Host native**—proyek Xcode (iOS) atau proyek Gradle (Android) yang dihasilkan untuk menyematkan kode Go hasil kompilasi Anda dan meng-host WebView.
- **Penyajian aset**—`frontend/dist/` Anda disematkan ke dalam biner Go dan disajikan di dalam proses yang sama. Server localhost tidak diperlukan.

---

## Sesuaikan aplikasi Anda untuk perangkat seluler

Aplikasi Anda sudah berfungsi, tetapi di layar ponsel tampilannya masih seperti aplikasi desktop. Beberapa perubahan kecil akan memberikan perbedaan besar.

### CSS responsif

Layar perangkat seluler lebih sempit dan menggunakan pola input yang berbeda. Di `frontend/public/style.css` (atau file yang setara):

```css
/* Prevent horizontal scrolling */
body {
  overflow-x: hidden;
}

/* Touch-friendly tap targets */
button {
  min-height: 44px;
  min-width: 44px;
}

/* Respect the iOS safe area (notch, home indicator) */
body {
  padding-top: env(safe-area-inset-top);
  padding-bottom: env(safe-area-inset-bottom);
  padding-left: env(safe-area-inset-left);
  padding-right: env(safe-area-inset-right);
}
```

### Deteksi platform di Go

Gunakan tag build untuk menambahkan perilaku khusus platform tanpa membuat kode bersama menjadi berantakan.

Buat `mobile_ios.go` untuk kode khusus iOS:

```go {title="mobile_ios.go"}
//go:build ios

package main

import "github.com/wailsapp/wails/v3/pkg/application"

func platformOptions() application.IOSOptions {
    return application.IOSOptions{
        DisableBounce: true,
    }
}
```

Buat `mobile_android.go` untuk kode khusus Android:

```go {title="mobile_android.go"}
//go:build android

package main

import "github.com/wailsapp/wails/v3/pkg/application"

func platformOptions() application.AndroidOptions {
    return application.AndroidOptions{}
}
```

Buat `mobile_desktop.go` sebagai stub agar kode bersama juga dapat dikompilasi di desktop:

```go {title="mobile_desktop.go"}
//go:build !ios && !android

package main

type mobileOptions struct{}

func platformOptions() mobileOptions { return mobileOptions{} }
```

### Deteksi platform di JavaScript dan batasi UI khusus perangkat seluler

Objek runtime `IOS.*` dan `Android.*` hanya tersedia di platform masing-masing. Memanggilnya di desktop akan menyebabkan error. Pola yang tepat—seperti yang digunakan Kitchen Sink—adalah mendeteksi platform satu kali dan menyembunyikan seluruh kontrol khusus perangkat seluler:

```javascript
// Detect platform from the bridge the host injects into the WebView
const platform = (() => {
  if (typeof window.wails?.platform === 'function') return window.wails.platform(); // Android
  if (window.webkit?.messageHandlers?.external) return 'ios';
  return 'desktop';
})();

const isIOS     = platform === 'ios';
const isAndroid = platform === 'android';
const isMobile  = isIOS || isAndroid;

// Hide any element marked as mobile-only
document.querySelectorAll('.mobile-only').forEach(el => {
  el.style.display = isMobile ? '' : 'none';
});
```

Kemudian, dalam HTML Anda:

```html
<section class="mobile-only">
  <button id="btnHaptic">Haptic feedback</button>
</section>
```

Dengan cara ini, tombol khusus perangkat seluler tidak pernah dirender di desktop, dan Anda tidak perlu melindungi setiap pemanggilan satu per satu dengan pemeriksaan `if (isMobile)`.

Di sisi Go, padukan ini dengan stub bertag build agar handler peristiwa hanya didaftarkan pada platform yang memerlukannya:

```go {title="native_desktop.go"}
//go:build !ios && !android

package main

import "github.com/wailsapp/wails/v3/pkg/application"

// No-op on desktop — mobile tabs are hidden in the frontend so these
// events are never emitted.
func registerNativeFeatures(app *application.App) {}
```

```go {title="native_ios.go"}
//go:build ios

package main

import "github.com/wailsapp/wails/v3/pkg/application"

func registerNativeFeatures(app *application.App) {
    app.Event.On("common:haptic", func(e *application.CustomEvent) {
        // only compiled and called on iOS
        application.IOS.Haptic("medium")
    })
    // ... other handlers
}
```

Inilah pola persis yang digunakan oleh [Kitchen Sink](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples/mobile) — lihat `native_features_stub.go`, `native_features_ios.go`, dan `native_features_android.go`.

@note{type="note" title="Penamaan API fitur native \u0026 peristiwa"}
Ada dua konvensi yang perlu diketahui:

- **Fitur native di sisi Go menggunakan pengelola platform.** Panggil fitur tersebut melalui singleton `application.IOS.*` dan `application.Android.*` — misalnya `application.IOS.Haptic("medium")` atau `application.Android.Share(payload)`. Setiap pengelola hanya tersedia pada platformnya sendiri, sehingga pemanggilannya ditempatkan dalam file `//go:build ios` / `//go:build android`.
- **Peristiwa diberi namespace berdasarkan cakupannya.** Semua peristiwa yang dipahami oleh kedua platform menggunakan prefiks `common:*` (`common:haptic`, `common:location`, …); peristiwa yang hanya dapat dihasilkan atau ditangani oleh satu platform menggunakan `ios:*` atau `android:*` (misalnya `ios:backgroundTask`, `android:foregroundService`). Karena hampir semua fitur seluler digunakan bersama, frontend Anda cukup mempertahankan satu listener untuk setiap peristiwa di bawah `common:*`.

@end

### Tambahkan umpan balik haptik (iOS)

```javascript
import { IOS } from '@wailsio/runtime';

async function onButtonTap() {
  if (isIOS) {
    await IOS.Haptics.Impact({ style: 'medium' });
  }
  // ... rest of your handler
}
```

### Tambahkan getaran (Android)

```javascript
import { Android } from '@wailsio/runtime';

async function onButtonTap() {
  if (isAndroid) {
    await Android.Haptics.Vibrate(50); // 50ms
  }
}
```

---

## Build untuk produksi

@tabs{sync-key="mobile-platform"}
[iOS]
**Build simulator** (untuk pengujian pada simulator, tidak memerlukan penandatanganan):

```bash
wails3 task ios:package
wails3 task ios:deploy-simulator
```

**Build perangkat** (memerlukan identitas penandatanganan dan profil penyediaan):

```bash
wails3 task ios:package \
  IOS_PLATFORM=device \
  CODESIGN_IDENTITY="Apple Development: You (TEAMID)" \
  PROVISIONING_PROFILE=path/to/profile.mobileprovision

wails3 task ios:deploy-device   # installs via xcrun devicectl
```

**IPA distribusi** (untuk App Store atau TestFlight):

```bash
wails3 task ios:package:ipa IOS_PLATFORM=device \
  CODESIGN_IDENTITY="..." \
  PROVISIONING_PROFILE=path/to/distribution.mobileprovision
```

@note{type="tip"}
Untuk mengunggah ke App Store Connect, gunakan `wails3 task ios:xcode` dan biarkan Xcode mengelola penandatanganan serta pengarsipan — Xcode secara otomatis menangani kerumitan sertifikat, profil, dan notarisasi.

@end

[Android]
**APK debug** (ditandatangani dengan keystore debug Android dan dapat langsung diinstal):

```bash
wails3 task android:package
wails3 task android:deploy-emulator
```

**APK rilis** (ditandatangani dengan keystore Anda sendiri):

```bash
ANDROID_KEYSTORE_FILE=/path/to/release.jks \
ANDROID_KEYSTORE_PASSWORD=yourpassword \
ANDROID_KEY_ALIAS=youralias \
ANDROID_KEY_PASSWORD=yourkeypassword \
  wails3 task android:package
```

**APK universal** (arm64 + x86_64 dalam satu file):

```bash
wails3 task android:package:fat
```

@note{type="tip"}
Untuk mengunggah ke Play Store, buat `.aab` (Android App Bundle), bukan APK — buka `build/android/` di Android Studio, lalu gunakan **Build → Generate Signed Bundle / APK**.

@end

@end

---

## Pemecahan masalah

### `wails3 task ios:run` gagal dengan pesan "no iOS SDKs found"

Xcode versi lengkap harus diinstal dan dipilih:

```bash
sudo xcode-select -s /Applications/Xcode.app/Contents/Developer
xcode-select -p  # should print the Xcode path
```

#### `wails3 task android:run` gagal dengan pesan "SDK not found"

Pastikan `ANDROID_HOME` telah ditetapkan dan diekspor. Verifikasi dengan:

```bash
echo $ANDROID_HOME
ls $ANDROID_HOME/platform-tools/adb
```

#### Simulator tidak melakukan boot

Tampilkan daftar simulator yang tersedia, lalu lakukan boot pada salah satunya secara manual:

```bash
xcrun simctl list devices available
xcrun simctl boot "iPhone 16"
```

#### `chrome://inspect` tidak menampilkan target

WebView harus berada dalam mode debug (default untuk `android:run`). Pastikan Anda menjalankan build debug, bukan build produksi. Pastikan juga `adb devices` menampilkan emulator sebagai perangkat yang terhubung.

#### Inset area aman tidak diterapkan

Pastikan HTML Anda menyertakan tag meta viewport:

```html
<meta name="viewport" content="width=device-width, initial-scale=1, viewport-fit=cover">
```

---

## Jelajahi Kitchen Sink

Setelah aplikasi pertama Anda berjalan, contoh **Kitchen Sink** adalah cara tercepat untuk mempelajari kemungkinan lainnya. Ini adalah aplikasi Wails lengkap yang berjalan di iOS, Android, dan desktop dari satu basis kode, serta mencakup fitur haptik, geolokasi, biometrik, notifikasi lokal, penyimpanan aman, dan lainnya:

```bash
git clone https://github.com/wailsapp/wails.git
cd wails/v3/examples/mobile

wails3 task ios:run        # iOS Simulator
wails3 task android:run    # Android Emulator
wails3 task run            # Desktop
```

Telusuri kode sumbernya di [`v3/examples/mobile`](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples/mobile) — file `native_features_ios.go` dan `native_features_android.go` sangat berguna sebagai titik awal yang dapat disalin dan ditempel untuk fitur khusus platform.

## Langkah berikutnya

@cards{cols="2"}
Panduan iOS
Referensi lengkap: opsi konfigurasi, tab native, pengaturan WKWebView, build perangkat, dan penandatanganan.

[Panduan iOS →](/guides/mobile/ios/)

---
Panduan Android
Referensi lengkap: konfigurasi, toast, pemaketan untuk Play Store, dan detail NDK.

[Panduan Android →](/guides/mobile/android/)

---
📖 Kode sumber Kitchen Sink
Fitur haptik, geolokasi, biometrik, notifikasi, dan penyimpanan aman — semuanya dalam satu aplikasi yang dapat dijalankan.

[Lihat di GitHub →](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples/mobile)

@end
