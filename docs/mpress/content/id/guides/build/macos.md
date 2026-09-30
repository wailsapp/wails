---
title: "Pengemasan macOS"
description: "Kemas aplikasi Wails Anda untuk didistribusikan di macOS"
slug: "guides/build/macos"
sourcePath: "guides/build/macos.md"
---

## API Privat macOS

Wails v3 menggunakan API publik macOS secara default. Untuk mengaktifkan fitur yang memerlukan API Apple yang tidak terdokumentasi, build aplikasi Anda dengan satu tag build Go `private_mac_apis`:

```bash
wails3 build -tags private_mac_apis
EXTRA_TAGS=private_mac_apis wails3 dev
wails3 package GOOS=darwin EXTRA_TAGS=private_mac_apis
```

Untuk build Go secara langsung, gunakan `go build -tags private_mac_apis .` (atau `-tags production,private_mac_apis` untuk produksi). Aplikasi yang sudah ada dan bergantung pada perilaku privat harus menambahkan tag ini untuk mempertahankan perilaku tersebut. Tag ini hanya berlaku untuk build desktop macOS.

Lihat [API Privat macOS](/guides/build/private-macos-apis/) untuk mengetahui daftar lengkap fitur dan nilai opsi yang terdampak, fallback build publik yang tepat, pemetaan gaya Liquid Glass, serta kombinasi build inspector. Operasi khusus privat tidak melakukan apa pun tanpa tag tersebut; API Go publik tidak berubah.

## Bundel Aplikasi

Kemas aplikasi Anda sebagai bundel `.app` macOS standar:

```bash
wails3 package GOOS=darwin
```

Tindakan ini membuat `bin/<AppName>.app` yang berisi:

- Biner hasil kompilasi di `Contents/MacOS/`
- Ikon aplikasi di `Contents/Resources/` (dari `icons.icns` atau, jika tersedia, dari katalog aset `Assets.car`)
- `Info.plist` dengan metadata aplikasi

## Sumber Daya Bundel

`Contents/Resources/` adalah lokasi standar untuk file hanya-baca yang disertakan bersama aplikasi macOS. Gunakan lokasi ini untuk templat berukuran besar, data awal, media, paket bahasa, atau muatan lain yang sebaiknya dibuka sesuai kebutuhan, bukan dikompilasi ke dalam executable Go dengan `embed`.

Wails sudah menempatkan ikon aplikasi di direktori ini. Untuk menambahkan file Anda sendiri, tempatkan file tersebut dalam direktori sumber seperti `build/resources/`, lalu tambahkan langkah penyalinan ke task `create:app:bundle` di `build/darwin/Taskfile.yml`:

```yaml
tasks:
  create:app:bundle:
    cmds:
      # Existing bundle creation commands...
      - |
          if [ -d build/resources ]; then
            cp -R build/resources/. "{{.BIN_DIR}}/{{.APP_NAME}}.app/Contents/Resources/"
          fi
```

Jika Anda menggunakan task `darwin:run` dari Taskfile, tambahkan perintah yang setara ke task `run` miliknya dengan target `{{.BIN_DIR}}/{{.APP_NAME}}.dev.app/Contents/Resources/`.

### Membaca Sumber Daya dari Go

Impor paket platform macOS:

```go
import (
	"io/fs"

	"github.com/wailsapp/wails/v3/pkg/mac"
)
```

Untuk file kecil, gunakan `LoadResource`:

```go
func loadSplash() ([]byte, error) {
	return mac.LoadResource("images/splash.png")
}
```

Untuk file yang lebih besar, gunakan `ResourceFS`. Fungsi ini mengembalikan `io/fs.FS` dengan root di `Contents/Resources` sehingga pemanggil dapat membuka dan melakukan streaming sumber daya tanpa terlebih dahulu memuat seluruhnya ke dalam slice byte Go:

```go
func openCatalogue() (fs.File, error) {
	resources, err := mac.ResourceFS()
	if err != nil {
		return nil, err
	}

	return resources.Open("catalogue/defaults.json")
}
```

Nama sumber daya adalah path yang dipisahkan garis miring dan relatif terhadap `Contents/Resources`. `ResourceFS` dan `LoadResource` mengembalikan `mac.ErrNotInAppBundle` kecuali executable dijalankan dari `.app/Contents/MacOS`.

Perlakukan sumber daya bundel sebagai data yang tidak dapat diubah. Mengubah file di dalam aplikasi yang telah ditandatangani akan membatalkan tanda tangan kodenya; sebagai gantinya, simpan data yang diunduh, dihasilkan, atau dapat diedit pengguna di direktori Application Support milik pengguna.

### Biner Universal

Build untuk Mac berbasis Apple Silicon dan Intel:

```bash
wails3 task darwin:package:universal
```

Tindakan ini membuat satu `.app` yang berjalan secara native pada kedua arsitektur. Biner universal dapat dibuat di platform apa pun—di Linux dan Windows, `wails3 tool lipo` digunakan secara otomatis.

## Menyesuaikan Bundel

Edit `build/darwin/Info.plist` untuk menyesuaikan:

- Pengidentifikasi bundel (`CFBundleIdentifier`)
- Nama dan versi aplikasi
- Versi minimum macOS
- Asosiasi file
- Skema URL

Ikon aplikasi dihasilkan dari aset dalam direktori `build/`. Gunakan task `generate:icons`:

```bash
wails3 task common:generate:icons
```

Task ini menggunakan `build/appicon.png` untuk menghasilkan `darwin/icons.icns` dan `windows/icon.ico`. Di macOS, Anda juga dapat menyediakan `build/appicon.icon` (format Icon Composer): task tersebut meneruskan `-iconcomposerinput appicon.icon -macassetdir darwin`, yang menghasilkan `Assets.car` dan `darwin/icons.icns` dari file `.icon` (dilewati pada platform selain macOS). Jika `Assets.car` tersedia, jalankan task `update:build-assets` agar `Info.plist` dan `CFBundleIconName` diperbarui sebagaimana mestinya:

```bash
wails3 task common:update:build-assets
```

Untuk menjalankan perintah ikon secara manual dari direktori `build/`:

```bash
cd build
wails3 generate icons -input appicon.png -macfilename darwin/icons.icns -windowsfilename windows/icon.ico -iconcomposerinput appicon.icon -macassetdir darwin
```

## Penandatanganan Kode

Tandatangani aplikasi Anda untuk didistribusikan:

```bash
# Using the wrapper (auto-detects platform)
wails3 sign GOOS=darwin

# Or using the task directly
wails3 task darwin:sign
```

Konfigurasikan penandatanganan di `build/darwin/Taskfile.yml`:

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  ENTITLEMENTS: "build/darwin/entitlements.plist"
```

### Notarisasi

Untuk aplikasi yang didistribusikan di luar Mac App Store, Apple mewajibkan notarisasi:

```bash
wails3 task darwin:sign:notarize
```

Pertama, simpan kredensial Anda. Jalankan wizard interaktif (`wails3 setup signing`) atau panggil `notarytool` secara langsung:

```bash
xcrun notarytool store-credentials "my-notarize-profile" \
  --apple-id "you@email.com" \
  --team-id "TEAMID" \
  --password "app-specific-password"
```

Konfigurasikan di `build/darwin/Taskfile.yml`:

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  KEYCHAIN_PROFILE: "my-notarize-profile"
```

Lihat [Menandatangani Aplikasi](/guides/build/signing/) untuk detail selengkapnya.

## Penginstal DMG

Templat yang disertakan dengan Wails 3 menyediakan `wails3 task darwin:package:dmg`. Templat ini terlebih dahulu membuat `.app`, lalu membuat DMG bergaya dengan pustaka DMG. Secara default, DMG menggunakan latar belakang gradien bermerek Wails dengan simbol naga merah dan logotipe WAILS.

```bash
wails3 task darwin:package:dmg
```

Task tingkat rendah `darwin:create:dmg` membuat DMG dari bundel `.app` yang sudah ada dan dapat dikonfigurasi langsung dari Taskfile:

```yaml
vars:
  # These are the template defaults; override them when needed.
  DMG_BACKGROUND: build/darwin/dmg-background.png
  DMG_VOLUME_ICON: build/darwin/icons.icns
  DMG_FILE_ICON: build/darwin/dmg-file-icon.icns
  DMG_WINDOW_WIDTH: 540
  DMG_WINDOW_HEIGHT: 380
  DMG_FILES: "Install.command=build/darwin/Install.command,README.txt=README.md"
```

### Tata Letak Default

DMG yang dihasilkan berisi:

- Bundel aplikasi di sebelah kiri
- Tautan `Applications` di sebelah kanan
- Jendela Finder berukuran 540×380 piksel
- Ukuran ikon 96 poin dengan label di bawah setiap ikon
- Latar belakang bermerek Wails dari `build/darwin/dmg-background.png`

Ikon aplikasi dan `Applications` diposisikan relatif terhadap dimensi jendela yang dikonfigurasi. Dengan demikian, mengubah `DMG_WINDOW_WIDTH` atau `DMG_WINDOW_HEIGHT` akan mempertahankan jarak proporsional pada tata letak dua ikon default. Untuk hasil terbaik, gunakan gambar latar belakang dengan dimensi piksel yang sama seperti jendela Finder.

### Mengganti Aset DMG

File yang dihasilkan di bawah `build/darwin/` merupakan aset proyek biasa dan dapat diganti:

- `DMG_BACKGROUND` mengatur gambar yang ditampilkan di belakang isi jendela Finder.
- `DMG_VOLUME_ICON` mengatur ikon yang ditampilkan untuk volume yang dipasang.
- `DMG_FILE_ICON` mengatur ikon yang ditampilkan di Finder untuk file `.dmg` yang dihasilkan.

Ikon volume dan ikon file DMG merupakan sumber daya yang terpisah. Mengganti ikon aplikasi tidak secara otomatis mengganti salah satunya.

### Menambahkan File Tambahan

Gunakan `DMG_FILES` untuk menyertakan skrip penginstal, catatan rilis, lisensi, atau sumber daya lain bersama aplikasi. Nilainya berupa daftar pasangan `name=path` yang dipisahkan dengan koma:

```yaml
vars:
  DMG_FILES: "Install.command=build/darwin/Install.command,README.txt=README.md"
```

Nama sebelum `=` adalah nama file yang ditampilkan di dalam DMG. Jalur setelah `=` adalah file sumber dalam proyek. Spasi kosong di awal dan akhir akan diabaikan.

Setiap nama yang ditampilkan harus unik. File tambahan tidak dapat menggantikan entri yang telah dibuat oleh pemaket, termasuk bundel aplikasi atau entri `Applications`. Nama yang berkonflik menyebabkan proses pemaketan gagal dengan pesan kesalahan, alih-alih menghasilkan DMG yang rusak.

@note{type="note"}
Pembuatan DMG hanya didukung di macOS karena menggunakan alat citra disk dan Finder dari macOS. Bundel `.app` hasil kompilasi silang dapat dibuat di platform lain, tetapi DMG akhir harus dibuat di Mac.

@end

## Pemecahan Masalah

### "Aplikasi rusak dan tidak dapat dibuka"

Aplikasi belum ditandatangani. Tandatangani aplikasi dengan sertifikat Developer ID, atau pengguna dapat melewati Gatekeeper:

```bash
xattr -cr /path/to/YourApp.app
```

### Notarisasi gagal

Masalah umum:

- **Kredensial tidak valid**: Jalankan kembali `xcrun notarytool store-credentials` (atau `wails3 setup signing`)
- **Runtime yang diperkeras diperlukan**: Pastikan entitlement menyertakan `com.apple.security.cs.allow-unsigned-executable-memory` jika diperlukan
- **Stempel waktu tidak ada**: Proses penandatanganan seharusnya menyertakan stempel waktu secara otomatis

### Aplikasi hasil kompilasi silang tidak dapat dijalankan

Biner macOS hasil kompilasi silang tidak ditandatangani. Transfer ke Mac dan tandatangani sebelum pengujian:

```bash
codesign --force --deep --sign - YourApp.app
```
