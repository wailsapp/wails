---
title: "Pintasan Global"
description: "Daftarkan pintasan papan ketik tingkat sistem yang tetap dipicu meskipun aplikasi Anda tidak sedang memiliki fokus"
slug: "features/keyboard/global-shortcuts"
sourcePath: "features/keyboard/global-shortcuts.md"
---

Pintasan global adalah pintasan papan ketik tingkat sistem yang dipicu tanpa bergantung pada aplikasi mana yang sedang memiliki fokus, selama aplikasi Wails Anda berjalan. Pintasan ini ideal untuk tombol pintas tampilkan/sembunyikan, alat perekaman cepat, kontrol media, dan fitur lain yang diharapkan pengguna dapat diakses dari mana saja.

@note{type="info" title="Pintasan global dan pengikatan tombol"}
[Pengikatan tombol](/features/keyboard/shortcuts/) (`app.KeyBinding`) hanya dipicu saat salah satu jendela aplikasi Anda memiliki fokus. Pintasan global (`app.GlobalShortcut`) dipicu di seluruh sistem, bahkan saat aplikasi Anda berjalan di latar belakang. Gunakan yang sesuai dengan kebutuhan Anda.

@end

Pintasan global dibangun langsung di atas fasilitas native setiap platform dan tidak menambahkan dependensi pihak ketiga.

## Mengakses Pengelola Pintasan Global

Pengelola tersedia melalui properti `GlobalShortcut` pada instans aplikasi Anda:

```go
app := application.New(application.Options{
    Name: "Global Shortcuts Demo",
})

globalShortcuts := app.GlobalShortcut
```

## Mendaftarkan Pintasan

`Register` menerima akselerator dan callback. Callback berjalan dalam goroutine tersendiri setiap kali pintasan ditekan.

```go
err := app.GlobalShortcut.Register("CmdOrCtrl+Shift+G", func() {
    // Runs even when another application is focused.
    window.Show()
    window.Focus()
})
if err != nil {
    app.Logger.Error("could not register shortcut", "error", err)
}
```

Anda dapat mendaftarkan pintasan sebelum memanggil `app.Run`. Pengikatan dengan sistem operasi kemudian dilakukan secara otomatis saat aplikasi dimulai.

@note{type="tip" title="Mengakses UI dari callback"}
Callback berjalan di luar thread utama. Jika callback perlu berinteraksi dengan jendela atau UI lain, metode jendela akan menanganinya untuk Anda. Namun, untuk pekerjaan khusus pada thread utama, bungkus pekerjaan tersebut dengan `application.InvokeSync`.

@end

### Format Akselerator

Pintasan global menggunakan format akselerator yang sama dengan akselerator menu dan pengikatan tombol:

```go
"CmdOrCtrl+Shift+G"  // Command on macOS, Control elsewhere
"Ctrl+Alt+K"         // Control + Alt + K
"Cmd+Option+Space"   // Command + Option + Space (macOS)
"Super+D"            // Super / Windows / Logo key + D
"Ctrl+Shift+F5"      // Function keys are supported
```

`CmdOrCtrl` diartikan sebagai Command di macOS serta Control di Windows dan Linux, sehingga praktis digunakan untuk pintasan lintas platform.

## Mengelola Pintasan

```go
// Check whether a shortcut is registered (modifier order does not matter).
registered := app.GlobalShortcut.IsRegistered("Ctrl+Shift+G")

// List every shortcut this application has registered.
for _, accelerator := range app.GlobalShortcut.GetAll() {
    app.Logger.Info("global shortcut", "accelerator", accelerator)
}

// Release a single shortcut.
app.GlobalShortcut.Unregister("Ctrl+Shift+G")

// Release everything (also done automatically on shutdown).
app.GlobalShortcut.UnregisterAll()
```

Semua pintasan yang terdaftar dilepas secara otomatis saat aplikasi keluar, sehingga Anda tidak perlu membersihkannya secara manual.

## Yang Terjadi Saat Pintasan yang Sama Didaftarkan Dua Kali

Ada dua kasus yang berbeda, dan Wails menanganinya dengan cara yang berbeda.

### Aplikasi yang sama mendaftarkan pintasan dua kali

Kasus ini ditangani oleh Wails sendiri dan perilakunya identik di setiap platform. Pemanggilan `Register` yang kedua mengembalikan kesalahan, sedangkan pengikatan awal tetap dipertahankan ("kesalahan dan pertahankan"). Hal ini membuat perilaku tetap dapat diprediksi dan menampilkan kesalahan tersebut alih-alih mengganti pintasan yang berfungsi secara diam-diam.

```go
app.GlobalShortcut.Register("Ctrl+Shift+G", showWindow)        // ok
err := app.GlobalShortcut.Register("Shift+Ctrl+G", doSomething) // err: already registered
// showWindow is still the active callback for this shortcut.
```

Jika ingin mengubah callback untuk suatu pintasan, lakukan `Unregister` terlebih dahulu, lalu lakukan `Register` lagi.

### Aplikasi lain sudah memiliki pintasan tersebut

Hal ini ditentukan oleh sistem operasi, sehingga hasilnya bergantung pada platform:

| Platform | Perilaku saat aplikasi lain memiliki pintasan tersebut |
| --- | --- |
| **macOS** | Pendaftaran berhasil. macOS mengizinkan beberapa aplikasi mendaftarkan tombol pintas yang sama, sehingga callback Anda ditambahkan bersama pemilik yang sudah ada dan tidak ditolak. |
| **Windows** | Pendaftaran gagal dan `Register` mengembalikan kesalahan. Pintasan tetap dimiliki oleh aplikasi yang pertama kali mendaftarkannya. |
| **Linux (X11)** | Pendaftaran gagal dan `Register` mengembalikan kesalahan karena server X menolak pengambilan kedua atas kombinasi yang sama. |
| **Linux (Wayland)** | Compositor bertindak sebagai penentu. Pengguna biasanya diminta menyetujui atau memilih pengikatan melalui dialog pintasan global di lingkungan desktop. |

Karena perbedaan ini, selalu periksa kesalahan yang dikembalikan oleh `Register` dan sediakan pintasan alternatif atau umpan balik kepada pengguna saat suatu pintasan tidak dapat diambil.

## Pertimbangan Platform

@tabs
[macOS]
Pintasan global menggunakan API tombol pintas Carbon Event Manager. Ini adalah mekanisme standar untuk tombol pintas tingkat sistem di macOS dan tidak memerlukan izin Aksesibilitas.

Pintasan huruf menyesuaikan tata letak. API Carbon mengikat ke posisi tombol fisik, sehingga Wails menerjemahkan sebuah huruf ke tombol yang menghasilkan huruf tersebut pada tata letak papan ketik aktif Anda (melalui `UCKeyTranslate`). Karena itu, `Cmd+A` menggunakan tombol berlabel **A** pada tata letak AZERTY, QWERTZ, Dvorak, dan tata letak serupa, bukan posisi QWERTY. Digit, tanda baca, dan tombol bernama (tombol panah, tombol fungsi, dan sebagainya) menggunakan posisi fisik tetapnya. Pengikatan ditentukan saat pintasan didaftarkan. Jika Anda mengganti tata letak papan ketik ketika aplikasi sedang berjalan, pintasan huruf yang sudah ada tetap menggunakan tombol fisik yang awalnya diikat.

@note{type="caution" title="Pintasan sembunyikan dan `ApplicationShouldTerminateAfterLastWindowClosed`"}
`window.Hide()` menggunakan `orderOut:` di macOS, yang membuat jendela tidak terlihat. AppKit menganggap jendela terakhir yang tidak terlihat sebagai jendela yang sudah ditutup. Jadi, jika Anda menetapkan `Mac.ApplicationShouldTerminateAfterLastWindowClosed: true` dan menggunakan pintasan global untuk menyembunyikan satu-satunya jendela, aplikasi akan berhenti alih-alih tetap berjalan di latar belakang. Biarkan opsi tersebut tidak ditetapkan (keadaan bawaan) jika Anda mengandalkan tombol pintas tampilkan/sembunyikan agar jendela dapat disembunyikan dan ditampilkan kembali nanti.

@end

[Windows]
Pintasan global menggunakan API Win32 `RegisterHotKey`. Pengulangan otomatis dinonaktifkan, sehingga menahan tombol hanya memicu callback satu kali, bukan berulang kali.

Pendaftaran gagal jika aplikasi lain sudah memiliki kombinasi tersebut. Karena itu, pilih kombinasi yang lebih jarang digunakan sebagai nilai bawaan.

[Linux]
Dalam sesi **X11**, Wails mengambil pintasan langsung dari server X, sehingga akselerator yang diminta diikat persis seperti yang ditentukan.

Dalam sesi **Wayland**, aplikasi memang dirancang agar tidak dapat mengambil tombol secara langsung. Sebagai gantinya, Wails menggunakan antarmuka XDG Desktop Portal `org.freedesktop.portal.GlobalShortcuts`. Dengan portal tersebut, akselerator yang Anda berikan merupakan pemicu *yang diutamakan*, sedangkan compositor (dan pada akhirnya pengguna) menentukan kombinasi tombol akhir. Callback Anda tetap dipicu ketika pintasan diaktifkan, tetapi tombol yang sebenarnya tidak dijamin sesuai dengan permintaan Anda, dan `IsRegistered`/`GetAll` melaporkan kombinasi yang Anda minta, bukan kombinasi yang diikat oleh compositor.

Portal tersebut memerlukan lingkungan desktop yang mengimplementasikan portal pintasan global (misalnya GNOME atau KDE Plasma versi terbaru).

@end

## Contoh Lengkap

```go
package main

import (
    "log"

    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "Global Shortcuts Demo",
    })

    window := app.Window.New()

    // Bring the window to the front from anywhere.
    if err := app.GlobalShortcut.Register("CmdOrCtrl+Shift+G", func() {
        window.Show()
        window.Focus()
    }); err != nil {
        log.Printf("could not register show shortcut: %v", err)
    }

    // Hide the window from anywhere.
    if err := app.GlobalShortcut.Register("CmdOrCtrl+Shift+H", func() {
        window.Hide()
    }); err != nil {
        log.Printf("could not register hide shortcut: %v", err)
    }

    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}
```

@note{type="danger" title="Hindari pintasan sistem yang penting"}
Beberapa kombinasi dicadangkan oleh sistem operasi atau lingkungan desktop dan tidak dapat diambil oleh aplikasi. Pilih nilai bawaan yang kemungkinan kecil berbenturan, dan selalu tangani kesalahan dari `Register`.

@end
