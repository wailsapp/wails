---
title: "Android"
description: "Bangun dan jalankan aplikasi Wails di Android — penyiapan toolchain, emulator, penandatanganan APK, pemaketan untuk Play Store, dan referensi API"
slug: "guides/mobile/android"
sourcePath: "guides/mobile/android.md"
---

@note{type="caution" title="Fitur Eksperimental"}
Dukungan Android bersifat eksperimental dan dapat berubah dalam rilis mendatang.

@end

@note{type="tip"}
Baru menggunakan Wails untuk perangkat seluler? Mulailah dengan [Aplikasi Seluler Pertama Anda →](/guides/mobile/first-mobile-app/) untuk panduan langkah demi langkah, lalu kembali ke sini untuk membaca referensi lengkap.

@end

Aplikasi Wails v3 berjalan di Android sebagai aplikasi native: sebuah `WebView` merender frontend, aset disajikan **dalam proses** melalui sebuah `WebViewAssetLoader` yang didukung oleh server aset Go (tanpa server localhost dan tanpa port terbuka), sedangkan `@wailsio/runtime` standar berfungsi tanpa perubahan — binding layanan, peristiwa, dialog, dan papan klip dirutekan melalui pemroses pesan Go.

`main.go` yang sama dapat dibangun untuk desktop dan Android. Kode Go dikompilasi sebagai pustaka bersama C (`libwails.so`, `GOOS=android` + toolchain NDK) dan dimuat oleh host Java kecil. Perilaku khusus Android berada dalam file Go per platform yang dilindungi oleh `//go:build android`.

## Persyaratan

- **Android SDK** dengan platform-tools, platform SDK (API 35), build-tools, dan **NDK** (26.3.x) — `wails3 doctor` menampilkan komponen yang ditemukannya
- **JDK** (misalnya OpenJDK 21) untuk Gradle; atur `JAVA_HOME` jika `java` tidak ada di `PATH` Anda
- Go 1.25+ dan npm
- `ANDROID_HOME` (atau `ANDROID_SDK_ROOT`) yang mengarah ke SDK

Instal komponen SDK dengan alat baris perintah:

```bash
sdkmanager "platform-tools" "platforms;android-35" "build-tools;35.0.0" \
           "ndk;26.3.11579264" "emulator" \
           "system-images;android-35;google_apis;arm64-v8a"
avdmanager create avd --name wails \
           --package "system-images;android-35;google_apis;arm64-v8a" \
           --device pixel_7
```

## Menjalankan di Emulator

Dari direktori proyek Anda:

```bash
wails3 task android:run
```

Perintah ini menjalankan emulator jika belum ada yang berjalan, menghasilkan binding, membangun frontend, mengompilasi kode Go Anda menjadi `libwails.so` untuk ABI emulator, merakit APK debug dengan Gradle, lalu menginstal dan menjalankannya.

Perintah pendamping yang berguna:

```bash
wails3 task android:logs    # stream the app's logcat output
```

Dalam build debug, WebView dapat diperiksa dari Chrome di `chrome://inspect`.

## Pemaketan

```bash
wails3 task android:package             # production release APK
wails3 task android:deploy-emulator     # install + launch it
wails3 task android:bundle              # production release AAB (Android App Bundle)
wails3 task android:bundle:fat          # release AAB containing all ABIs
wails3 task android:run:device          # debug install + launch on a physical device
wails3 task android:deploy-device       # install + launch on a physical device
DEVICE_ID=<serial> wails3 task android:run:device
DEVICE_ID=<serial> wails3 task android:deploy-device
```

Build produksi menggunakan `-tags production,android`, menghapus simbol yang tidak diperlukan, dan tidak menyertakan diagnostik internal framework saat kompilasi. `wails3 task android:package:fat` membangun `arm64-v8a` dan `x86_64` ke dalam satu APK.

Google Play mewajibkan format Android App Bundle (`.aab`) untuk pengajuan aplikasi baru, dan pengajuan baru harus menargetkan Android 15 (API 35) atau yang lebih tinggi; templat proyek menetapkan `compileSdk` dan `targetSdk` ke 35 di `build/android/app/build.gradle`. `wails3 task android:bundle:fat` menghasilkan `bin/<AppName>.aab` yang menyertakan kedua ABI; Google Play menggunakannya untuk menghasilkan APK per perangkat yang dioptimalkan, sehingga bundle gemuk tersebut merupakan artefak yang tepat untuk diunggah ke toko aplikasi. APK tetap menjadi cara tercepat untuk pengujian lokal dan emulator karena `.aab` tidak dapat diinstal secara langsung dengan `adb`.

`android:run` dan `android:deploy-emulator` adalah task yang ditujukan untuk emulator. Untuk perangkat Android fisik, gunakan `android:run:device` untuk APK debug atau `android:deploy-device` untuk APK rilis. Keduanya membangun untuk `arm64`, memilih entri non-emulator pertama yang terhubung dari `adb devices`, menginstalnya, lalu menjalankan `com.wails.app.MainActivity`. Teruskan `DEVICE_ID=<serial>` untuk menargetkan perangkat tertentu.

## Penandatanganan & build rilis

Tanpa keystore, build rilis ditandatangani dengan keystore **debug** Android agar dapat diinstal untuk pengujian. Untuk menandatanganinya dengan keystore Anda sendiri, tetapkan:

```bash
ANDROID_KEYSTORE_FILE=/path/to/release.jks \
ANDROID_KEYSTORE_PASSWORD=... \
ANDROID_KEY_ALIAS=... \
ANDROID_KEY_PASSWORD=... \
  wails3 task android:package
```

Variabel yang sama juga menandatangani App Bundle: jalankan `wails3 task android:bundle:fat` setelah menetapkan variabel tersebut untuk menghasilkan `.aab` yang siap untuk Play. Tanpa variabel itu, bundle ditandatangani dengan keystore debug dan akan ditolak oleh Google Play, sehingga task tersebut mencetak peringatan.

@note{type="tip"}
Dengan [Play App Signing](https://support.google.com/googleplay/android-developer/answer/9842756), keystore yang Anda gunakan untuk menandatangani secara lokal adalah **kunci upload** Anda: Google menggunakannya untuk memverifikasi upload Anda, lalu menandatangani ulang aplikasi dengan kunci penandatanganan aplikasi yang dikelolanya. Perhatikan juga bahwa Google Play mewajibkan `versionCode` yang lebih tinggi untuk setiap upload; tingkatkan nilainya di `build/android/app/build.gradle`.

@end

## Konfigurasi

Frontend mengendalikan fitur Android saat runtime melalui objek runtime `Android`: `Android.Haptics.Vibrate(durationMs)`, `Android.Device.Info()`, `Android.Toast.Show(message)`. Nama paket dikendalikan oleh `APP_ID` dalam task build.

## Yang berfungsi dan yang tidak

| Area | Status |
| --- | --- |
| WebView + aset dalam proses (`WebViewAssetLoader`) | ✅ |
| Binding layanan, peristiwa (dua arah) | ✅ |
| Dialog pesan | ✅ AlertDialog dengan callback tombol |
| Dialog untuk membuka satu/beberapa file | ✅ Storage Access Framework (file diimpor sebagai salinan cache) |
| Dialog untuk membuka direktori/menyimpan file | ❌ Mengembalikan galat — sebagai gantinya, tulis di dalam sandbox aplikasi |
| Papan klip | ✅ ClipboardManager |
| API Screens | ✅ WindowMetrics termasuk area kerja di luar bilah sistem |
| Peristiwa siklus hidup (`events.Android.*`) | ✅ |
| Haptik, informasi perangkat, toast | ✅ API runtime `Android.*` |
| Build untuk emulator + perangkat fisik | ✅ `android:run`, `android:run:device`, `android:deploy-emulator`, `android:deploy-device` |
| Geometri jendela, menu, baki sistem | Sengaja tidak melakukan apa pun |
| Beberapa jendela | Hanya jendela pertama yang ditampilkan |

## Catatan porting

- Kode desktop dikompilasi tanpa perubahan di bawah `GOOS=android`; pemanggilan geometri/menu/baki tidak melakukan apa pun karena aplikasi Android berjalan dalam layar penuh.
- `android` **menyiratkan tag build `linux`** (Android menggunakan kernel Linux): file yang hanya ditujukan untuk Linux desktop memerlukan `//go:build linux && !android`, dan saat runtime `runtime.GOOS` bernilai `"android"`.
- Ganti dialog penyimpanan file dan pemilihan direktori dengan penulisan ke sandbox aplikasi beserta alur berbagi melalui intent. Dialog pembukaan file berfungsi dan mengimpor dokumen yang dipilih sebagai salinan di direktori cache, sehingga Anda memperoleh jalur sistem file yang sebenarnya.
- Aplikasi nyata selalu dibangun dengan `CGO_ENABLED=1` dan NDK; jalur non-cgo hanya tersedia agar alat seperti `wails3 generate bindings` dapat memuat paket tersebut.
- Rancang frontend secara responsif; area kerja `Screens` tidak mencakup bilah status dan navigasi.
