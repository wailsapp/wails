---
title: "Mulai Otomatis"
description: "Daftarkan aplikasi Anda agar diluncurkan saat pengguna masuk di macOS, Windows, dan Linux"
slug: "features/autostart/basics"
sourcePath: "features/autostart/basics.md"
---

## Mulai Otomatis

`app.Autostart` mendaftarkan aplikasi Anda agar diluncurkan secara otomatis saat pengguna masuk. API ini memilih mekanisme native yang tepat untuk setiap platform dan menyelesaikan jalur instalasi yang menggunakan symlink (Homebrew, Scoop), sehingga pendaftaran tidak rusak ketika biner diperbarui.

Pendaftaran mulai berlaku saat **login berikutnya**, bukan seketika.

## Mulai Cepat

```go
import "github.com/wailsapp/wails/v3/pkg/application"

// Register to launch at login
if err := app.Autostart.Enable(); err != nil {
    app.Logger.Error("autostart enable failed", "error", err)
}

// Stop launching at login
if err := app.Autostart.Disable(); err != nil {
    app.Logger.Error("autostart disable failed", "error", err)
}

// Check status
enabled, err := app.Autostart.IsEnabled()
```

## API

### `Enable`

Mendaftarkan aplikasi agar diluncurkan saat login dengan opsi default.

```go
func (m *AutostartManager) Enable() error
```

Memanggil `Enable` berulang kali aman—pendaftaran ditimpa setiap kali, sehingga Anda dapat memanggilnya pada setiap startup jika preferensi pengguna telah disimpan secara persisten.

### `EnableWithOptions`

Mendaftarkan aplikasi dengan opsi khusus.

```go
func (m *AutostartManager) EnableWithOptions(opts AutostartOptions) error
```

**`AutostartOptions`:**

| Kolom | Tipe | Deskripsi |
| --- | --- | --- |
| `Identifier` | `string` | Mengganti ID pendaftaran yang diturunkan secara otomatis. Lihat "Pengidentifikasi" di bawah. |
| `Arguments` | `[]string` | Argumen tambahan yang ditambahkan ke jalur executable ketika diluncurkan saat login (misalnya, `--hidden`). |

### `Disable`

Menghapus pendaftaran mulai otomatis. Mengembalikan `nil` jika aplikasi tidak terdaftar—operasi penonaktifan bersifat idempoten.

```go
func (m *AutostartManager) Disable() error
```

### `IsEnabled`

Melaporkan apakah pendaftaran tersedia. Operasi ini cepat karena tidak memvalidasi jalur yang terdaftar.

```go
func (m *AutostartManager) IsEnabled() (bool, error)
```

### `Status`

Mengembalikan status pendaftaran lengkap.

```go
func (m *AutostartManager) Status() (AutostartStatus, error)
```

**`AutostartStatus`:**

| Kolom | Tipe | Deskripsi |
| --- | --- | --- |
| `Enabled` | `bool` | Apakah pendaftaran tersedia. |
| `Path` | `string` | Lokasi artefak pendaftaran pada disk (jalur plist, jalur `.desktop`, atau subkunci registri). Kosong ketika `Enabled` bernilai false. |
| `Strategy` | `AutostartStrategy` | Mekanisme yang mendaftarkan aplikasi (lihat [Perilaku Platform](#perilaku-platform)). |

## Perilaku Platform

@tabs{sync-key="platform"}
[macOS]
Salah satu dari dua mekanisme digunakan, tergantung pada cara aplikasi dikemas:

- **macOS 13+, `.app` yang dibundel**: `SMAppService.mainAppService`. Berfungsi untuk aplikasi dalam sandbox dan build Mac App Store. Tidak ada prompt otomatisasi TCC (pendekatan AppleScript lama memicunya).
- **macOS sebelum 13 atau biner yang tidak dibundel**: plist LaunchAgent ditulis ke `~/Library/LaunchAgents/<identifier>.plist` dengan `RunAtLoad=true`.

`Status()` mengembalikan `AutostartStrategySMAppService` atau `AutostartStrategyLaunchAgent` agar pemanggil dapat mengetahui jalur yang digunakan.

Ketika aplikasi ditingkatkan dari tidak dibundel menjadi dibundel, kedua jalur diperiksa oleh `Status()` dan dibersihkan oleh `Disable()` agar LaunchAgent yatim tidak terus meluncurkan build lama.

[Windows]
Sebuah nilai registri ditambahkan di bawah `HKCU\Software\Microsoft\Windows\CurrentVersion\Run`, dengan nama nilai ditetapkan ke `Identifier` mulai otomatis dan datanya ditetapkan ke jalur executable yang diapit tanda kutip beserta setiap `Arguments`.

Pengutipan argumen mengikuti aturan `CommandLineToArgvW` (garis miring terbalik digandakan sebelum tanda kutip), sehingga jalur yang mengandung spasi atau tanda kutip dapat diproses bolak-balik dengan benar.

`Status().Strategy` mengembalikan `AutostartStrategyRegistryRun`.

[Linux]
Entri mulai otomatis XDG ditulis ke `$XDG_CONFIG_HOME/autostart/<identifier>.desktop` (secara default ke `~/.config/autostart/`) dengan:

```ini
Type=Application
Hidden=false
X-GNOME-Autostart-enabled=true
Exec=<executable> <arguments>
```

Kolom `Exec` di-escape sesuai [spesifikasi Desktop Entry freedesktop.org](https://specifications.freedesktop.org/desktop-entry-spec/)—karakter khusus (`"`, `` ` ``, `$`, `\\`) di-escape dengan garis miring terbalik dan nilainya diapit tanda kutip ganda jika mengandung spasi.

`Status().Strategy` mengembalikan `AutostartStrategyXDGAutostart`.

[iOS / Android / server]
Tidak didukung. Semua metode mengembalikan `ErrAutostartNotSupported`. Gunakan `errors.Is(err, application.ErrAutostartNotSupported)` untuk mendeteksinya secara tepat:

```go
if err := app.Autostart.Enable(); err != nil {
    if errors.Is(err, application.ErrAutostartNotSupported) {
        // hide the toggle in the UI
        return
    }
    // real failure — surface it
}
```

@end

## Pengidentifikasi

Jika `Options.Identifier` kosong, nilai default diturunkan dari nama aplikasi Anda:

| Platform | Pengidentifikasi default |
| --- | --- |
| macOS (dibundel) | Pengidentifikasi bundle aplikasi, misalnya `com.example.MyApp` |
| macOS (tidak dibundel) | `wails.autostart.<slug>`, dengan `<slug>` diturunkan dari `application.Options.Name` |
| Windows | Slug dari `application.Options.Name` (huruf kecil, karakter selain `A-Za-z0-9._-` dihapus, dan spasi diubah menjadi tanda hubung) |
| Linux | Slug yang sama seperti di Windows |

Pengidentifikasi harus cocok dengan `^[A-Za-z0-9._-]+$` dan panjangnya tidak boleh melebihi 200 karakter. Format DNS terbalik direkomendasikan untuk macOS (sesuai dengan cara Label launchd biasanya ditulis).

Ketika `AutostartOptions.Identifier` diganti, pengidentifikasi yang sama digunakan kembali sebagai nama nilai registri di Windows dan nama file `.desktop` di Linux, sehingga satu string dapat mengidentifikasi pendaftaran lintas platform.

## Deteksi Entri Usang

`Disable()` dan `Status()` menemukan pendaftaran dengan **mencocokkan jalur executable yang terdaftar dengan `os.Executable()` (setelah menyelesaikan setiap symlink)**, bukan dengan mencari pengidentifikasi. Artinya:

- **Mengubah pengidentifikasi antar-rilis aman.** Pendaftaran lama masih dapat ditemukan oleh `Status()` dan dibersihkan oleh `Disable()`—selama jalur executable tetap sama.
- **Salinan kedua aplikasi di jalur yang berbeda tidak akan menimpa pendaftaran salinan pertama.** Setiap lokasi biner dilacak secara independen.
- **Instalasi yang menggunakan symlink (Homebrew, Scoop) tetap stabil.** `filepath.EvalSymlinks` diterapkan ke `os.Executable()` sebelum pencocokan, sehingga peningkatan Homebrew yang mengganti target tidak membuat entri tersebut telantar.

Hal yang *tidak* ditangani oleh mekanisme ini: jika pengguna memindahkan atau mengganti nama berkas biner ke jalur lain yang tidak terkait, pendaftaran lama menjadi yatim (pendaftaran tersebut menunjuk ke berkas yang kini tidak ada). Aplikasi yang didistribusikan dari jalur instalasi tetap tidak perlu mengkhawatirkan hal ini; aplikasi yang didistribusikan sebagai berkas biner portabel tunggal harus memanggil `Disable()` sebelum memindahkan dirinya sendiri, atau selalu dijalankan melalui symlink tetap.

## Contoh

```go
package main

import (
    "errors"

    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "My App",
    })

    // Restore the user's preference on startup
    if userPrefersAutostart() {
        if err := app.Autostart.Enable(); err != nil {
            if !errors.Is(err, application.ErrAutostartNotSupported) {
                app.Logger.Error("autostart", "error", err)
            }
        }
    }

    app.Run()
}
```

Contoh lengkap yang dapat dijalankan, dengan tombol status / aktifkan / nonaktifkan, tersedia di [`examples/autostart/`](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/autostart).
