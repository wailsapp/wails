---
title: "Jendela Tanpa Bingkai"
description: "Buat dekorasi jendela khusus dengan jendela tanpa bingkai"
slug: "features/windows/frameless"
sourcePath: "features/windows/frameless.md"
---

## Jendela Tanpa Bingkai

Wails menyediakan **dukungan jendela tanpa bingkai** dengan area seret berbasis CSS dan perilaku native platform. Hapus bilah judul native platform untuk mendapatkan kendali penuh atas dekorasi jendela, desain khusus, dan pengalaman pengguna yang unik, sekaligus mempertahankan fungsi penting seperti menyeret, mengubah ukuran, dan kontrol sistem.

![Aplikasi pemula TypeScript Wails v3 bawaan yang berjalan sebagai jendela tanpa bingkai dengan sudut native macOS](/assets/screenshots/frameless-v3-native-corners-macos.png)

Contoh di atas adalah aplikasi pemula TypeScript Wails v3 bawaan dengan `Frameless: true` diaktifkan.

## Mulai Cepat

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:     "Frameless App",
    Width:     800,
    Height:    600,
    Frameless: true,
})
```

**CSS untuk bilah judul yang dapat diseret:**

```css
.titlebar {
    --wails-draggable: drag;
    height: 40px;
    background: #333;
}

.titlebar button {
    --wails-draggable: no-drag;
}
```

**HTML:**

```html
<div class="titlebar">
    <span>My Application</span>
    <button onclick="window.close()">×</button>
</div>
```

**Selesai!** Kini Anda memiliki bilah judul khusus.

## Membuat Jendela Tanpa Bingkai

### Radius Sudut (macOS)

Secara bawaan, jendela tanpa bingkai mempertahankan sudut membulat standar macOS dari AppKit. Atur `Mac.CornerRadius` untuk menggunakan radius khusus (dalam poin):

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Mac: application.MacWindow{
        CornerRadius: 16,
    },
})
```

Atur `Mac.CornerType` ke `MacWindowCornerTypeSquare` untuk sudut siku-siku. Pengaturan ini mengabaikan `CornerRadius`:

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Mac: application.MacWindow{
        CornerType: application.MacWindowCornerTypeSquare,
    },
})
```

### Jendela Tanpa Bingkai Dasar

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Width:     800,
    Height:    600,
})
```

**Yang Anda dapatkan:**

- Tanpa bilah judul
- Tanpa batas jendela
- Tanpa tombol sistem
- Latar belakang transparan (opsional)

**Yang perlu Anda implementasikan:**

- Area yang dapat diseret
- Tombol tutup/minimalkan/maksimalkan
- Handel pengubahan ukuran (jika ukurannya dapat diubah)

### Dengan Latar Belakang Transparan

**API privat di macOS:** atur `Mac.Backdrop: application.MacBackdropTransparent` dan lakukan build dengan `-tags private_mac_apis` agar webview transparan. Tanpa tag tersebut, webview native tetap opak meskipun latar belakang HTML/CSS transparan. `Frameless` dan `TitleBar.AppearsTransparent` sendiri menggunakan API publik. Lihat [API macOS privat](/guides/build/private-macos-apis/#webview-transparency-and-background).

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless:      true,
    BackgroundType: application.BackgroundTypeTransparent,
    Mac: application.MacWindow{
        Backdrop: application.MacBackdropTransparent, // Requires -tags private_mac_apis.
    },
})
```

**Kasus penggunaan:**

- Sudut membulat
- Bentuk khusus
- Jendela overlay
- Layar pembuka

## Area Seret

### Penyeretan Berbasis CSS

Gunakan properti CSS `--wails-draggable`:

```css
/* Draggable area */
.titlebar {
    --wails-draggable: drag;
}

/* Non-draggable elements within draggable area */
.titlebar button {
    --wails-draggable: no-drag;
}
```

**Nilai:**

- `drag` - Area dapat diseret
- `no-drag` - Area tidak dapat diseret (meskipun induknya dapat diseret)

### Contoh Bilah Judul Lengkap

```html
<div class="titlebar">
    <div class="title">My Application</div>
    <div class="controls">
        <button class="minimize">−</button>
        <button class="maximize">□</button>
        <button class="close">×</button>
    </div>
</div>
```

```css
.titlebar {
    --wails-draggable: drag;
    display: flex;
    justify-content: space-between;
    align-items: center;
    height: 40px;
    background: #2c2c2c;
    color: white;
    padding: 0 16px;
}

.title {
    font-size: 14px;
    user-select: none;
}

.controls {
    display: flex;
    gap: 8px;
}

.controls button {
    --wails-draggable: no-drag;
    width: 32px;
    height: 32px;
    border: none;
    background: transparent;
    color: white;
    font-size: 16px;
    cursor: pointer;
    border-radius: 4px;
}

.controls button:hover {
    background: rgba(255, 255, 255, 0.1);
}

.controls .close:hover {
    background: #e81123;
}
```

**JavaScript untuk tombol:**

```javascript
import { Window } from '@wailsio/runtime'

document.querySelector('.minimize').addEventListener('click', () => Window.Minimise())
document.querySelector('.maximize').addEventListener('click', () => Window.Maximise())
document.querySelector('.close').addEventListener('click', () => Window.Close())
```

## Area Nonklien Native di Windows

Windows dapat memperlakukan bagian dari bilah judul khusus sebagai area nonklien native. Dengan demikian, Anda dapat menggambar bilah judul dan tombol bilah judul dengan desain HTML/CSS apa pun sekaligus mempertahankan perilaku native Windows: area bilah judul menyeret jendela, tombol maksimalkan dapat menampilkan Snap Assist / Snap Layouts Windows 11, dan tombol minimalkan, maksimalkan, serta tutup menerima hit testing dan status mouse native.

Video di bawah menunjukkan bilah judul HTML/CSS khusus yang menggunakan hit testing native Windows, termasuk Snap Assist / Snap Layouts Windows 11 pada tombol maksimalkan khusus.

<video src="/assets/windows-native-non-client-regions/wails-app-region.mp4" controls muted playsInline></video>

Wails mendukung dua mekanisme khusus Windows:

- `app-region` melalui dukungan area nonklien native WebView2
- `--wails-non-client-region` melalui pelacakan runtime Wails untuk tombol bilah judul khusus

### Memilih Mode

@note{type="caution" title="Eksperimental"}
`WebView2CompositionHosting` mengubah cara jendela meng-host dan berinteraksi dengan WebView2 di balik layar. Alih-alih menggunakan pengontrol WebView2 bawaan yang di-host oleh HWND, Wails menggunakan hosting pengontrol komposisi dan meneruskan input secara eksplisit. Mode ini mungkin mengalami masalah rendering, input, fokus, atau kompatibilitas WebView2 Runtime. Aktifkan hanya jika Anda memerlukan perilaku native untuk tombol bilah judul khusus, lalu uji aplikasi Anda dengan cermat pada versi Windows dan WebView2 Runtime yang Anda dukung.

@end

Pilih berdasarkan kebutuhan Anda dari Windows:

- Gunakan `NonClientRegionSupport` untuk penyeretan aplikasi native sederhana dengan `app-region: drag` dan `app-region: no-drag` dari WebView2.
- Gunakan `WebView2CompositionHosting` jika tombol minimalkan, maksimalkan, dan tutup khusus Anda perlu berperilaku seperti tombol bilah judul native Windows.
- Aktifkan keduanya jika jendela yang sama memerlukan dukungan `app-region` native WebView2 dan area tombol bilah judul khusus yang dikelola Wails.

`NonClientRegionSupport` adalah alternatif native yang ringan untuk pelacakan `--wails-draggable` milik Wails. Anda menandai area yang dapat dan tidak dapat diseret dengan CSS, WebView2 menentukan piksel yang termasuk dalam area bilah judul, lalu Wails meminta area native tersebut dari WebView2 saat melakukan hit testing.

Itulah keseluruhan cakupan mode ini saat ini. Mode ini tidak membuat tombol minimalkan, maksimalkan, atau tutup khusus berperilaku seperti tombol bilah judul native Windows, dan tidak mengaktifkan Snap Assist / Snap Layouts Windows 11 untuk tombol maksimalkan khusus. Gunakan mode ini jika Anda memerlukan penyeretan aplikasi native sederhana tanpa mekanisme tambahan dari `--wails-draggable`.

`WebView2CompositionHosting` digunakan untuk tombol bilah judul kustom dengan perilaku native. Wails melacak persegi panjang DOM yang ditandai dengan `--wails-non-client-region`, memetakannya ke nilai uji hit Windows seperti `HTMINBUTTON`, `HTMAXBUTTON`, dan `HTCLOSE`, lalu meneruskan input mouse kembali ke permukaan WebView2 yang dihosting melalui komposisi. Hal inilah yang memungkinkan tombol maksimalkan kustom berpartisipasi dalam Snap Assist / Snap Layouts Windows 11 sambil tetap mempertahankan desain visual apa pun yang Anda pilih.

Dengan kata lain: `NonClientRegionSupport` adalah dukungan region CSS native WebView2. `WebView2CompositionHosting` berarti Wails bertanggung jawab atas komposisi yang dimiliki host dan uji hit non-klien kustom.

### app-region WebView2

Aktifkan dukungan region non-klien native WebView2 untuk jendela:

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Windows: application.WindowsWindow{
        NonClientRegionSupport: true,
    },
})
```

Kemudian tandai area yang dapat diseret dengan properti CSS `app-region`:

```css
.titlebar {
    app-region: drag;
}

.titlebar button,
.titlebar input,
.titlebar select,
.titlebar textarea {
    app-region: no-drag;
}
```

Gunakan ini jika Anda hanya memerlukan penyeretan bilah judul secara native dan kontrol bilah judul ditangani melalui klik frontend biasa.

Keterbatasan mode ini adalah cakupannya dibatasi oleh dukungan region non-klien WebView2 sendiri. Dalam rilis WebView2 saat ini, ini berarti hanya region seret dan jangan-seret. Mode ini tidak ditujukan untuk memodelkan tombol bilah judul frontend yang sepenuhnya kustom dengan peran native yang berbeda untuk meminimalkan, memaksimalkan, dan menutup.

### Tombol Bilah Judul Kustom dengan Perilaku Native

Untuk tombol minimalkan, maksimalkan, dan tutup kustom yang harus berperilaku seperti tombol bilah judul sistem, aktifkan hosting komposisi:

@note{type="caution" title="Eksperimental"}
`WebView2CompositionHosting` menggunakan hosting pengontrol komposisi WebView2 dengan DirectComposition. Lihat [Memilih Mode](#memilih-mode) sebelum mengaktifkannya.

@end

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Windows: application.WindowsWindow{
        WebView2CompositionHosting: true,
    },
})
```

Kemudian tandai setiap region frontend dengan `--wails-non-client-region`:

```html
<div class="titlebar">
    <div class="title">My Application</div>
    <div class="window-controls">
        <button class="window-button minimize" aria-label="Minimize"></button>
        <button class="window-button maximize" aria-label="Maximize"></button>
        <button class="window-button close" aria-label="Close"></button>
    </div>
</div>
```

```css
.titlebar {
    --wails-non-client-region: caption;
    height: 40px;
}

.window-controls {
    display: flex;
    height: 100%;
}

.window-button {
    width: 46px;
    border: 0;
    background: transparent;
}

.window-button.minimize {
    --wails-non-client-region: minimize;
}

.window-button.maximize {
    --wails-non-client-region: maximize;
}

.window-button.close {
    --wails-non-client-region: close;
}
```

Nilai `--wails-non-client-region` yang didukung:

- `caption` - area bilah judul yang dapat diseret
- `minimize` - target hit tombol minimalkan native
- `maximize` - target hit tombol maksimalkan native, termasuk perilaku saat penunjuk diarahkan untuk Snap Assist / Snap Layouts Windows 11
- `close` - target hit tombol tutup native

Runtime Wails mengamati perubahan DOM, gaya, ukuran, gulir, dan viewport, lalu mengirim snapshot region ke jendela native. Geometri region diukur dalam piksel CSS dan dikonversi menjadi piksel fisik untuk uji hit Windows.

Desain visual tetap sepenuhnya menjadi pilihan Anda. Region hanya memberi tahu Windows arti setiap persegi panjang; bentuk tombol, ikon, warna, jarak, gaya saat penunjuk diarahkan, dan tata letak tetap berasal dari frontend Anda.

### Menggabungkan Keduanya

Anda dapat mengaktifkan kedua opsi jika menginginkan dukungan `app-region` WebView2 dan region tombol bilah judul yang dikelola Wails dalam jendela yang sama:

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Windows: application.WindowsWindow{
        NonClientRegionSupport:    true,
        WebView2CompositionHosting: true,
    },
})
```

## Tombol Sistem

### Mengimplementasikan Tutup/Minimalkan/Maksimalkan

**Sisi Go:**

```go
type WindowControls struct {
    window *application.WebviewWindow
}

func (wc *WindowControls) Minimise() {
    wc.window.Minimise()
}

func (wc *WindowControls) Maximise() {
    if wc.window.IsMaximised() {
        wc.window.UnMaximise()
    } else {
        wc.window.Maximise()
    }
}

func (wc *WindowControls) Close() {
    wc.window.Close()
}
```

**Sisi JavaScript:**

```javascript
import { Minimise, Maximise, Close } from './bindings/WindowControls'

document.querySelector('.minimize').addEventListener('click', Minimise)
document.querySelector('.maximize').addEventListener('click', Maximise)
document.querySelector('.close').addEventListener('click', Close)
```

**Atau gunakan metode runtime:**

```javascript
import { Window } from '@wailsio/runtime'

document.querySelector('.minimize').addEventListener('click', () => Window.Minimise())
document.querySelector('.maximize').addEventListener('click', () => Window.Maximise())
document.querySelector('.close').addEventListener('click', () => Window.Close())
```

### Mengalihkan Status Maksimalkan

Lacak status maksimalkan untuk ikon tombol:

```javascript
import { Window } from '@wailsio/runtime'

async function toggleMaximise() {
    const isMaximised = await Window.IsMaximised()

    if (isMaximised) {
        await Window.Restore()
    } else {
        await Window.Maximise()
    }

    updateMaximiseButton()
}

async function updateMaximiseButton() {
    const isMaximised = await Window.IsMaximised()
    const button = document.querySelector('.maximize')
    button.textContent = isMaximised ? '❐' : '□'
}
```

## Pegangan Ubah Ukuran

### Ubah Ukuran Berbasis CSS

Wails menyediakan pegangan ubah ukuran otomatis untuk jendela tanpa bingkai:

```css
/* Enable resize on all edges */
body {
    --wails-resize: all;
}

/* Or specific edges */
.resize-top {
    --wails-resize: top;
}

.resize-bottom {
    --wails-resize: bottom;
}

.resize-left {
    --wails-resize: left;
}

.resize-right {
    --wails-resize: right;
}

/* Corners */
.resize-top-left {
    --wails-resize: top-left;
}

.resize-top-right {
    --wails-resize: top-right;
}

.resize-bottom-left {
    --wails-resize: bottom-left;
}

.resize-bottom-right {
    --wails-resize: bottom-right;
}
```

**Nilai:**

- `all` - Ubah ukuran dari semua tepi
- `top`, `bottom`, `left`, `right` - Tepi tertentu
- `top-left`, `top-right`, `bottom-left`, `bottom-right` - Sudut
- `none` - Jangan ubah ukuran

### Contoh Pegangan Ubah Ukuran

```html
<div class="window">
    <div class="titlebar">...</div>
    <div class="content">...</div>
    <div class="resize-handle resize-bottom-right"></div>
</div>
```

```css
.resize-handle {
    position: absolute;
    width: 16px;
    height: 16px;
}

.resize-bottom-right {
    --wails-resize: bottom-right;
    bottom: 0;
    right: 0;
    cursor: nwse-resize;
}
```

## Perilaku Khusus Platform

@tabs{sync-key="platform"}
[Windows]
**Jendela tanpa bingkai Windows:**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Windows: application.WindowsWindow{
        DisableFramelessWindowDecorations: false,
    },
})
```

**Fitur:**

- Bayangan jatuh otomatis
- Dukungan Snap Layouts (Windows 11)
- Dukungan Aero Snap
- Penskalaan DPI

**Nonaktifkan dekorasi:**

```go
Windows: application.WindowsWindow{
    DisableFramelessWindowDecorations: true,
},
```

**Snap Assist:**

```go
// Trigger Windows 11 Snap Assist
window.SnapAssist()
```

Ini memicu Snap Layouts melalui jalur tombol pintas Windows. Untuk tombol maksimalkan HTML kustom dengan Snap Layouts native saat penunjuk diarahkan, gunakan [Region Non-Klien Native di Windows](#area-nonklien-native-di-windows) sebagai gantinya.

**Tinggi bilah judul kustom:** Windows secara otomatis mendeteksi region seret dari CSS.

[macOS]
**Jendela tanpa bingkai macOS:**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Mac: application.MacWindow{
        TitleBar: application.MacTitleBar{
            AppearsTransparent: true,
        },
        InvisibleTitleBarHeight: 40,
    },
})
```

**Fitur:**

- Dukungan layar penuh native
- Tombol lampu lalu lintas (opsional)
- Efek vibrancy
- Bilah judul transparan

**Sembunyikan bilah judul sepenuhnya** (gunakan varian preset yang diekspor dari paket `application` — tidak ada bidang `TitleBarStyle` ataupun konstanta `MacTitleBarStyleHidden`):

```go
Mac: application.MacWindow{
    TitleBar: application.MacTitleBarHidden,
},
```

Preset lainnya mencakup `MacTitleBarDefault`, `MacTitleBarHiddenInset`, dan `MacTitleBarHiddenInsetUnified`.

**Bilah judul tak terlihat:** Memungkinkan jendela diseret saat bilah judul disembunyikan. Ini hanya berlaku jika jendela tidak berbingkai atau menggunakan `AppearsTransparent`:

```go
Mac: application.MacWindow{
    InvisibleTitleBarHeight: 40,
},
```

[Linux]
**Jendela tanpa bingkai di Linux:**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
})
```

**Fitur:**

- Dukungan dasar untuk jendela tanpa bingkai
- Area seret CSS
- Bervariasi menurut lingkungan desktop

**Catatan lingkungan desktop:**

- **GNOME:** Dukungan baik
- **KDE Plasma:** Dukungan baik
- **XFCE:** Dukungan dasar
- **WM tiling:** Dukungan terbatas

**Kompositor diperlukan:** Transparansi memerlukan kompositor (sebagian besar lingkungan desktop modern memilikinya).

@end

## Pola Umum

### Pola 1: Bilah Judul Modern

```html
<div class="modern-titlebar">
    <div class="app-icon">
        <img src="/icon.png" alt="App Icon">
    </div>
    <div class="title">My Application</div>
    <div class="controls">
        <button class="minimize">−</button>
        <button class="maximize">□</button>
        <button class="close">×</button>
    </div>
</div>
```

```css
.modern-titlebar {
    --wails-draggable: drag;
    display: flex;
    align-items: center;
    height: 40px;
    background: linear-gradient(to bottom, #3a3a3a, #2c2c2c);
    border-bottom: 1px solid #1a1a1a;
    padding: 0 16px;
}

.app-icon {
    --wails-draggable: no-drag;
    width: 24px;
    height: 24px;
    margin-right: 12px;
}

.title {
    flex: 1;
    font-size: 13px;
    color: #e0e0e0;
    user-select: none;
}

.controls {
    display: flex;
    gap: 1px;
}

.controls button {
    --wails-draggable: no-drag;
    width: 46px;
    height: 32px;
    border: none;
    background: transparent;
    color: #e0e0e0;
    font-size: 14px;
    cursor: pointer;
    transition: background 0.2s;
}

.controls button:hover {
    background: rgba(255, 255, 255, 0.1);
}

.controls .close:hover {
    background: #e81123;
    color: white;
}
```

### Pola 2: Layar Pembuka

```go
splash := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:          "Loading...",
    Width:          400,
    Height:         300,
    Frameless:      true,
    AlwaysOnTop:    true,
    BackgroundType: application.BackgroundTypeTransparent,
    Mac: application.MacWindow{
        Backdrop: application.MacBackdropTransparent, // Requires -tags private_mac_apis.
    },
})
```

```css
body {
    background: transparent;
    display: flex;
    justify-content: center;
    align-items: center;
}

.splash {
    background: white;
    border-radius: 12px;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
    padding: 40px;
    text-align: center;
}
```

### Pola 3: Jendela Bersudut Bulat

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless:      true,
    BackgroundType: application.BackgroundTypeTransparent,
    Mac: application.MacWindow{
        Backdrop: application.MacBackdropTransparent, // Requires -tags private_mac_apis.
    },
})
```

```css
body {
    background: transparent;
    margin: 8px;
}

.window {
    background: white;
    border-radius: 16px;
    box-shadow: 0 4px 24px rgba(0, 0, 0, 0.15);
    overflow: hidden;
    height: calc(100vh - 16px);
}

.titlebar {
    --wails-draggable: drag;
    background: #f5f5f5;
    border-bottom: 1px solid #e0e0e0;
}
```

### Pola 4: Jendela Overlay

```go
overlay := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless:      true,
    AlwaysOnTop:    true,
    BackgroundType: application.BackgroundTypeTransparent,
    Mac: application.MacWindow{
        Backdrop: application.MacBackdropTransparent, // Requires -tags private_mac_apis.
    },
})
```

```css
body {
    background: transparent;
}

.overlay {
    background: rgba(0, 0, 0, 0.8);
    backdrop-filter: blur(10px);
    border-radius: 8px;
    padding: 20px;
}
```

## Contoh Lengkap

Berikut adalah jendela tanpa bingkai yang siap digunakan dalam produksi:

**Go:**

```go
package main

import (
    _ "embed"
    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
    app := application.New(application.Options{
        Name: "Frameless App",
        Assets: application.AssetOptions{
            Handler: application.AssetFileServerFS(assets),
        },
    })

    window := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:     "Frameless Application",
        Width:     1000,
        Height:    700,
        MinWidth:  800,
        MinHeight: 600,
        Frameless: true,

        Mac: application.MacWindow{
            TitleBar: application.MacTitleBar{
                AppearsTransparent: true,
            },
            InvisibleTitleBarHeight: 40,
        },

        Windows: application.WindowsWindow{
            DisableFramelessWindowDecorations: false,
        },
    })

    window.Center()
    window.Show()

    app.Run()
}
```

**HTML:**

```html
<!DOCTYPE html>
<html>
<head>
    <link rel="stylesheet" href="/style.css">
</head>
<body>
    <div class="window">
        <div class="titlebar">
            <div class="title">Frameless Application</div>
            <div class="controls">
                <button class="minimize" title="Minimise">−</button>
                <button class="maximize" title="Maximise">□</button>
                <button class="close" title="Close">×</button>
            </div>
        </div>
        <div class="content">
            <h1>Hello from Frameless Window!</h1>
        </div>
    </div>
    <script src="/main.js" type="module"></script>
</body>
</html>
```

**CSS:**

```css
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

body {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    background: #f5f5f5;
}

.window {
    height: 100vh;
    display: flex;
    flex-direction: column;
}

.titlebar {
    --wails-draggable: drag;
    display: flex;
    justify-content: space-between;
    align-items: center;
    height: 40px;
    background: #ffffff;
    border-bottom: 1px solid #e0e0e0;
    padding: 0 16px;
}

.title {
    font-size: 13px;
    font-weight: 500;
    color: #333;
    user-select: none;
}

.controls {
    display: flex;
    gap: 8px;
}

.controls button {
    --wails-draggable: no-drag;
    width: 32px;
    height: 32px;
    border: none;
    background: transparent;
    color: #666;
    font-size: 16px;
    cursor: pointer;
    border-radius: 4px;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s;
}

.controls button:hover {
    background: #f0f0f0;
    color: #333;
}

.controls .close:hover {
    background: #e81123;
    color: white;
}

.content {
    flex: 1;
    padding: 40px;
    overflow: auto;
}
```

**JavaScript:**

```javascript
import { Window } from '@wailsio/runtime'

// Minimise button
document.querySelector('.minimize').addEventListener('click', () => {
    Window.Minimise()
})

// Maximise/restore button
const maximiseBtn = document.querySelector('.maximize')
maximiseBtn.addEventListener('click', async () => {
    const isMaximised = await Window.IsMaximised()

    if (isMaximised) {
        await Window.Restore()
    } else {
        await Window.Maximise()
    }

    updateMaximiseButton()
})

// Close button
document.querySelector('.close').addEventListener('click', () => {
    Window.Close()
})

// Update maximise button icon
async function updateMaximiseButton() {
    const isMaximised = await Window.IsMaximised()
    maximiseBtn.textContent = isMaximised ? '❐' : '□'
    maximiseBtn.title = isMaximised ? 'Restore' : 'Maximise'
}

// Initial state
updateMaximiseButton()
```

## Praktik Terbaik

### ✅ Lakukan

- **Sediakan area yang dapat diseret** - Pengguna perlu memindahkan jendela
- **Implementasikan tombol sistem** - Tutup, minimalkan, maksimalkan
- **Tetapkan ukuran minimum** - Cegah tata letak yang tidak dapat digunakan
- **Uji di semua platform** - Perilakunya bervariasi
- **Gunakan CSS untuk area seret** - Fleksibel dan mudah dipelihara
- **Berikan umpan balik visual** - Status saat penunjuk berada di atas tombol

### ❌ Jangan Lakukan

- **Jangan lupakan gagang pengubah ukuran** - Jika ukuran jendela dapat diubah
- **Jangan jadikan seluruh jendela dapat diseret** - Ini mencegah interaksi
- **Jangan lupa menonaktifkan seret pada tombol** - Tombol tidak akan berfungsi
- **Jangan gunakan area seret yang terlalu kecil** - Sulit untuk diseret
- **Jangan abaikan perbedaan antarplatform** - Uji secara menyeluruh

## Pemecahan Masalah

### Jendela Tidak Dapat Diseret

**Penyebab:** `--wails-draggable: drag` tidak ada

**Solusi:**

```css
.titlebar {
    --wails-draggable: drag;
}
```

### Tombol Tidak Berfungsi

**Penyebab:** Tombol berada di area yang dapat diseret

**Solusi:**

```css
.titlebar button {
    --wails-draggable: no-drag;
}
```

### Ukuran Jendela Tidak Dapat Diubah

**Penyebab:** Gagang pengubah ukuran tidak ada

**Solusi:**

```css
body {
    --wails-resize: all;
}
```

## Langkah Berikutnya

@cards{cols="2"}
▣ Dasar-Dasar Jendela
Pelajari dasar-dasar pengelolaan jendela.

[Pelajari Selengkapnya →](/features/windows/basics/)

---
⚙ Opsi Jendela
Referensi lengkap untuk opsi jendela.

[Pelajari Selengkapnya →](/features/windows/options/)

---
🚀 Peristiwa Jendela
Tangani peristiwa siklus hidup jendela.

[Pelajari Selengkapnya →](/features/windows/events/)

---
◆ Beberapa Jendela
Pola untuk aplikasi multijendela.

[Pelajari Selengkapnya →](/features/windows/multiple/)

@end

---

**Ada pertanyaan?** Tanyakan di [Discord](https://discord.gg/JDdSxwjhGf) atau lihat [contoh jendela tanpa bingkai](https://github.com/wailsapp/wails/tree/master/v3/examples/frameless).
