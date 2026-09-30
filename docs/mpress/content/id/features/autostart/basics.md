---
title: "Mulai Otomatis"
description: "Daftarkan aplikasi Anda agar diluncurkan saat pengguna masuk di macOS, Windows, dan Linux"
slug: "features/autostart/basics"
sourcePath: "features/autostart/basics.md"
---

## Mulai Otomatis

`app.Autostart` mendaftarkan aplikasi Anda agar diluncurkan secara otomatis saat pengguna masuk. API ini memilih mekanisme native yang tepat untuk setiap platform dan menyelesaikan jalur instalasi yang menggunakan symlink (Homebrew, Scoop), sehingga pendaftaran tidak rusak ketika biner ditingkatkan versinya.

Pendaftaran mulai berlaku pada **proses masuk berikutnya**, bukan saat itu juga.

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

Mendaftarkan aplikasi agar diluncurkan saat pengguna masuk dengan opsi default.

```go
func (m *AutostartManager) Enable() error
```

Memanggil `Enable` berulang kali aman — pendaftaran ditimpa setiap kali, sehingga Anda dapat memanggilnya pada setiap proses awal jika preferensi pengguna telah disimpan secara persisten.

### `EnableWithOptions`

Mendaftarkan aplikasi dengan opsi khusus.

```go
func (m *AutostartManager) EnableWithOptions(opts AutostartOptions) error
```

**`AutostartOptions`:**

| Kolom | Tipe | Deskripsi |
| --- | --- | --- |
| `Identifier` | `string` | Mengganti ID pendaftaran yang dihasilkan secara otomatis. Lihat "Pengidentifikasi" di bawah. |
| `Arguments` | `[]string` | Argumen tambahan yang ditambahkan ke jalur berkas yang dapat dieksekusi saat aplikasi diluncurkan ketika pengguna masuk (misalnya, `--hidden`). |

### `Disable`

Menghapus pendaftaran mulai otomatis. Mengembalikan `nil` jika aplikasi belum didaftarkan — operasi penonaktifan bersifat idempoten.

```go
func (m *AutostartManager) Disable() error
```

### `IsEnabled`

Melaporkan apakah pendaftaran tersedia. Cepat — tidak memvalidasi jalur yang terdaftar.

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
| `Path` | `string` | Lokasi artefak pendaftaran pada disk (jalur plist, jalur `.desktop`, atau subkunci registri). Kosong jika `Enabled` bernilai false. |
| `Strategy` | `AutostartStrategy` | Mekanisme yang digunakan untuk mendaftarkan aplikasi (lihat [Perilaku Platform](#perilaku-platform)). |

## Perilaku Platform

@tabs{sync-key="platform"}
[macOS]
Salah satu dari dua mekanisme digunakan, tergantung pada cara aplikasi dikemas:

- **macOS 13+, `.app`** yang dibundel: `SMAppService.mainAppService`. Berfungsi untuk aplikasi dalam sandbox dan build Mac App Store. Tidak ada dialog permintaan izin otomatisasi TCC (pendekatan AppleScript yang lama memicu dialog tersebut).
- **macOS sebelum 13, atau biner yang tidak dibundel**: plist LaunchAgent yang ditulis ke `~/Library/LaunchAgents/<identifier>.plist` dengan `RunAtLoad=true`.

`Status()` mengembalikan `AutostartStrategySMAppService` atau `AutostartStrategyLaunchAgent` agar pemanggil dapat mengetahui jalur yang digunakan.

Saat aplikasi ditingkatkan dari tidak dibundel menjadi dibundel, kedua jalur diperiksa oleh `Status()` dan dibersihkan oleh `Disable()` agar LaunchAgent yatim tidak terus meluncurkan build lama.

[Windows]
Sebuah nilai registri ditambahkan di bawah `HKCU\Software\Microsoft\Windows\CurrentVersion\Run`, dengan nama nilai ditetapkan ke `Identifier` mulai otomatis dan datanya ditetapkan ke jalur berkas yang dapat dieksekusi dalam tanda kutip beserta setiap `Arguments`.

Pengutipan argumen mengikuti aturan `CommandLineToArgvW` (garis miring terbalik digandakan sebelum tanda kutip), sehingga jalur yang berisi spasi atau tanda kutip dapat diproses bolak-balik dengan benar.

`Status().Strategy` mengembalikan `AutostartStrategyRegistryRun`.

[Linux]
Entri mulai otomatis XDG ditulis ke `$XDG_CONFIG_HOME/autostart/<identifier>.desktop` (secara default ke `~/.config/autostart/`) dengan:

```ini
Type=Application
Hidden=false
X-GNOME-Autostart-enabled=true
Exec=<executable> <arguments>
```

Kolom `Exec` di-escape sesuai [spesifikasi Desktop Entry freedesktop.org](https://specifications.freedesktop.org/desktop-entry-spec/) — karakter khusus (`"`, `` ` ``, `$`, `\\`) di-escape dengan garis miring terbalik dan nilainya diapit tanda kutip ganda jika mengandung spasi kosong.

`Status().Strategy` mengembalikan `AutostartStrategyXDGAutostart`.

[iOS / Android / server]
Tidak didukung. Semua metode mengembalikan `ErrAutostartNotSupported`. Gunakan `errors.Is(err, application.ErrAutostartNotSupported)` untuk mendeteksinya dengan baik:

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
| macOS (dibundel) | Pengidentifikasi bundel aplikasi, misalnya `com.example.MyApp` |
| macOS (tidak dibundel) | `wails.autostart.<slug>`, dengan `<slug>` diturunkan dari `application.Options.Name` |
| Windows | Slug dari `application.Options.Name` (huruf kecil, karakter selain `A-Za-z0-9._-` dihapus, spasi menjadi tanda hubung) |
| Linux | Slug yang sama seperti Windows |

Pengidentifikasi harus cocok dengan `^[A-Za-z0-9._-]+$` dan panjangnya tidak boleh melebihi 200 karakter. Format DNS terbalik direkomendasikan untuk macOS (format ini sesuai dengan cara Label launchd ditulis secara konvensional).

Ketika `AutostartOptions.Identifier` ditimpa, pengenal yang sama digunakan kembali sebagai nama nilai registri di Windows dan nama file `.desktop` di Linux, sehingga satu string dapat mengidentifikasi pendaftaran di berbagai platform.

## Deteksi Entri Usang

`Disable()` dan `Status()` menemukan pendaftaran dengan **mencocokkan jalur berkas yang dapat dieksekusi yang terdaftar dengan `os.Executable()` (setelah semua symlink diurai)**, bukan dengan mencari pengenal. Artinya:

- **Mengubah pengenal di antara rilis tidak menimbulkan masalah.** Pendaftaran lama masih dapat ditemukan oleh `Status()` dan dibersihkan oleh `Disable()` — selama jalur berkas yang dapat dieksekusi tetap sama.
- **Salinan kedua aplikasi di jalur yang berbeda tidak akan menimpa pendaftaran salinan pertama.** Setiap lokasi biner dilacak secara independen.
- **Instalasi melalui symlink (Homebrew, Scoop) tetap stabil.** `filepath.EvalSymlinks` diterapkan pada `os.Executable()` sebelum pencocokan, sehingga pemutakhiran Homebrew yang mengganti target tidak meninggalkan entri tanpa rujukan.

Hal yang *tidak* ditangani oleh mekanisme ini: jika pengguna memindahkan atau mengganti nama biner ke jalur lain yang tidak terkait, pendaftaran lama menjadi yatim (pendaftaran tersebut menunjuk ke berkas yang kini tidak ada). Aplikasi yang didistribusikan dari jalur instalasi tetap tidak perlu mengkhawatirkan hal ini; aplikasi yang didistribusikan sebagai biner portabel satu berkas sebaiknya memanggil `Disable()` sebelum memindahkan dirinya sendiri, atau selalu diluncurkan melalui symlink yang stabil.

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

Contoh lengkap yang dapat dijalankan dengan tombol status / aktifkan / nonaktifkan tersedia di [`examples/autostart/`](https://github.com/wailsapp/wails/tree/master/v3/examples/autostart).
