---
title: "Jendela Tanpa Bingkai"
description: "Buat dekorasi jendela kustom dengan jendela tanpa bingkai"
slug: "features/windows/frameless"
sourcePath: "features/windows/frameless.md"
---

## Jendela Tanpa Bingkai

Wails menyediakan **dukungan jendela tanpa bingkai** dengan area seret berbasis CSS dan perilaku native platform. Hapus bilah judul native platform untuk mendapatkan kendali penuh atas dekorasi jendela, desain kustom, dan pengalaman pengguna yang unik, sembari mempertahankan fungsi penting seperti menyeret, mengubah ukuran, dan kontrol sistem.

![Aplikasi awal TypeScript v3 bawaan Wails yang berjalan sebagai jendela tanpa bingkai dengan sudut native macOS](/assets/screenshots/frameless-v3-native-corners-macos.png)

Contoh di atas adalah aplikasi awal TypeScript v3 bawaan Wails dengan `Frameless: true` yang diaktifkan.

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

**Selesai!** Kini Anda memiliki bilah judul kustom.

## Membuat Jendela Tanpa Bingkai

### Radius Sudut (macOS)

Secara default, jendela tanpa bingkai mempertahankan sudut membulat standar macOS dari AppKit. Atur `Mac.CornerRadius` untuk menggunakan radius kustom (dalam poin):

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Mac: application.MacWindow{
        CornerRadius: 16,
    },
})
```

Atur `Mac.CornerType` menjadi `MacWindowCornerTypeSquare` untuk sudut siku-siku. Pengaturan ini mengabaikan `CornerRadius`:

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Frameless: true,
    Mac: application.MacWindow{
        CornerType: application.MacWindowCornerTypeSquare,
    },
})
```

### Jendela Dasar Tanpa Bingkai

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
- Gagang pengubah ukuran (jika ukurannya dapat diubah)

### Dengan Latar Belakang Transparan

**API privat di macOS:** atur `Mac.Backdrop: application.MacBackdropTransparent` dan lakukan build dengan `-tags private_mac_apis` untuk transparansi webview. Tanpa tag tersebut, webview native tetap buram meskipun latar belakang HTML/CSS transparan. `Frameless` dan `TitleBar.AppearsTransparent` sendiri menggunakan API publik. Lihat [API macOS privat](/guides/build/private-macos-apis/#webview-transparency-and-background).

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
- Bentuk kustom
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
- `no-drag` - Area tidak dapat diseret (meskipun elemen induknya dapat diseret)

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

## Area Non-Klien Native di Windows

Windows dapat memperlakukan bagian dari bilah judul kustom sebagai area non-klien native. Dengan demikian, Anda dapat menggambar bilah judul dan tombol keterangan menggunakan desain HTML/CSS apa pun sembari mempertahankan perilaku native Windows: area keterangan menyeret jendela, tombol maksimalkan dapat menampilkan 11 Snap Assist / Snap Layouts Windows, dan tombol minimalkan, maksimalkan, serta tutup menerima pengujian hit dan status mouse native.

Video di bawah ini menampilkan bilah judul HTML/CSS kustom yang menggunakan pengujian hit native Windows, termasuk 11 Snap Assist / Snap Layouts Windows pada tombol maksimalkan kustom.

<video src="/assets/windows-native-non-client-regions/wails-app-region.mp4" controls muted playsInline></video>

Wails mendukung dua mekanisme khusus Windows:

- `app-region` melalui dukungan area non-klien native WebView2
- `--wails-non-client-region` melalui pelacakan runtime Wails untuk tombol keterangan kustom

### Memilih Mode

@note{type="caution" title="Eksperimental"}
`WebView2CompositionHosting` mengubah cara jendela meng-host dan berinteraksi dengan WebView2 di balik layar. Alih-alih menggunakan pengontrol WebView2 bawaan yang di-host HWND, Wails menggunakan hosting pengontrol komposisi dan meneruskan input secara eksplisit. Mode ini mungkin mengalami masalah rendering, input, fokus, atau kompatibilitas WebView2 Runtime. Aktifkan hanya jika Anda memerlukan perilaku native pada tombol keterangan kustom, lalu uji aplikasi secara cermat pada versi Windows dan WebView2 Runtime yang Anda dukung.

@end

Pilih berdasarkan kebutuhan Anda dari Windows:

- Gunakan `NonClientRegionSupport` untuk penyeretan aplikasi native sederhana dengan `app-region: drag` dan `app-region: no-drag` milik WebView2.
- Gunakan `WebView2CompositionHosting` jika tombol minimalkan, maksimalkan, dan tutup kustom Anda harus berperilaku seperti tombol keterangan native Windows.
- Aktifkan keduanya jika jendela yang sama memerlukan dukungan `app-region` native WebView2 dan area tombol keterangan kustom yang dikelola Wails.

`NonClientRegionSupport` adalah alternatif native yang ringan untuk pelacakan `--wails-draggable` Wails. Anda menandai area yang dapat dan tidak dapat diseret dengan CSS, WebView2 menentukan piksel mana yang termasuk dalam keterangan, lalu Wails meminta area native kepada WebView2 saat melakukan pengujian hit.

Itulah seluruh cakupan mode ini saat ini. Mode ini tidak membuat tombol minimalkan, maksimalkan, atau tutup kustom berperilaku seperti tombol keterangan native Windows, dan tidak mengaktifkan 11 Snap Assist / Snap Layouts Windows untuk tombol maksimalkan kustom. Gunakan mode ini jika Anda memerlukan penyeretan aplikasi native sederhana tanpa mekanisme tambahan dari `--wails-draggable`.

`WebView2CompositionHosting` ditujukan untuk tombol keterangan kustom dengan perilaku native. Wails melacak persegi panjang DOM yang ditandai dengan `--wails-non-client-region`, memetakannya ke nilai pengujian hit Windows seperti `HTMINBUTTON`, `HTMAXBUTTON`, dan `HTCLOSE`, lalu meneruskan input mouse kembali ke permukaan WebView2 yang di-host melalui komposisi. Mekanisme inilah yang memungkinkan tombol maksimalkan kustom berpartisipasi dalam 11 Snap Assist / Snap Layouts Windows sembari mempertahankan desain visual apa pun yang Anda pilih.

Dengan kata lain: `NonClientRegionSupport` adalah dukungan area CSS native WebView2. Pada `WebView2CompositionHosting`, Wails bertanggung jawab atas komposisi yang dimiliki host dan pengujian hit non-klien kustom.

### app-region WebView2

Aktifkan dukungan area non-klien native WebView2 untuk jendela:

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

Gunakan mode ini jika Anda hanya memerlukan penyeretan keterangan native dan kontrol bilah judul Anda ditangani melalui klik frontend biasa.

Keterbatasan mode ini adalah cakupannya dibatasi oleh dukungan area non-klien WebView2 itu sendiri. Pada rilis WebView2 saat ini, dukungan tersebut hanya mencakup area seret dan tanpa seret. Mode ini tidak ditujukan untuk memodelkan tombol keterangan frontend yang sepenuhnya kustom dengan peran native minimalkan, maksimalkan, dan tutup yang berbeda.

### Tombol Keterangan Kustom dengan Perilaku Native

Untuk tombol minimalkan, maksimalkan, dan tutup kustom yang harus berperilaku seperti tombol keterangan sistem, aktifkan hosting komposisi:

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

Kemudian tandai setiap area frontend dengan `--wails-non-client-region`:

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
- `minimize` - target klik tombol perkecil native
- `maximize` - target klik tombol maksimalkan native, termasuk perilaku saat mengarahkan penunjuk pada Snap Assist / Snap Layouts Windows 11
- `close` - target klik tombol tutup native

Runtime Wails mengamati perubahan DOM, gaya, ukuran, gulir, dan viewport, lalu mengirim snapshot region ke jendela native. Geometri region diukur dalam piksel CSS dan dikonversi menjadi piksel fisik untuk pengujian klik Windows.

Desain visual tetap sepenuhnya Anda tentukan. Region hanya memberi tahu Windows arti setiap persegi panjang; bentuk tombol, ikon, warna, jarak, gaya saat penunjuk diarahkan, dan tata letak tetap berasal dari frontend Anda.

### Menggabungkan Keduanya

Anda dapat mengaktifkan kedua opsi saat menginginkan dukungan `app-region` WebView2 dan region tombol bilah judul yang dikelola Wails dalam jendela yang sama:

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

### Mengimplementasikan Tutup/Perkecil/Maksimalkan

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

### Mengalihkan Status Maksimal

Lacak status maksimal untuk ikon tombol:

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

## Gagang Pengubah Ukuran

### Pengubahan Ukuran Berbasis CSS

Wails menyediakan gagang pengubah ukuran otomatis untuk jendela tanpa bingkai:

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
- `none` - Tidak dapat diubah ukurannya

### Contoh Gagang Pengubah Ukuran

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
**Jendela tanpa bingkai di Windows:**

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

Ini memicu Snap Layouts melalui jalur tombol pintas Windows. Untuk tombol maksimalkan HTML khusus dengan Snap Layouts native saat penunjuk diarahkan, gunakan [Region Non-Klien Native di Windows](#area-non-klien-native-di-windows) sebagai gantinya.

**Tinggi bilah judul khusus:** Windows secara otomatis mendeteksi region seret dari CSS.

[macOS]
**Jendela tanpa bingkai di macOS:**

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

**Sembunyikan bilah judul sepenuhnya** (gunakan varian preset yang diekspor dari paket `application` — tidak ada field `TitleBarStyle` atau konstanta `MacTitleBarStyleHidden`):

```go
Mac: application.MacWindow{
    TitleBar: application.MacTitleBarHidden,
},
```

Preset lainnya mencakup `MacTitleBarDefault`, `MacTitleBarHiddenInset`, dan `MacTitleBarHiddenInsetUnified`.

**Bilah judul tak terlihat:** Memungkinkan jendela diseret sambil menyembunyikan bilah judul. Ini hanya berlaku saat jendela tidak memiliki bingkai atau menggunakan `AppearsTransparent`:

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
- Region seret CSS
- Bervariasi menurut lingkungan desktop

**Catatan lingkungan desktop:**

- **GNOME:** Dukungan baik
- **KDE Plasma:** Dukungan baik
- **XFCE:** Dukungan dasar
- **WM tiling:** Dukungan terbatas

**Memerlukan compositor:** Transparansi memerlukan compositor (sebagian besar DE modern memilikinya).

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
- **Implementasikan tombol sistem** - Tutup, perkecil, maksimalkan
- **Tetapkan ukuran minimum** - Cegah tata letak yang tidak dapat digunakan
- **Uji pada semua platform** - Perilakunya berbeda-beda
- **Gunakan CSS untuk region seret** - Fleksibel dan mudah dipelihara
- **Berikan umpan balik visual** - Status saat penunjuk diarahkan pada tombol

### ❌ Jangan Lakukan

- **Jangan lupakan gagang pengubah ukuran** - Jika ukuran jendela dapat diubah
- **Jangan jadikan seluruh jendela dapat diseret** - Menghalangi interaksi
- **Jangan lupakan no-drag pada tombol** - Tombol tidak akan berfungsi
- **Jangan gunakan area seret yang terlalu kecil** - Sulit diraih
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

**Penyebab:** Handel pengubah ukuran tidak ada

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
Pola untuk aplikasi dengan beberapa jendela.

[Pelajari Selengkapnya →](/features/windows/multiple/)

@end

---

**Ada pertanyaan?** Tanyakan di [Discord](https://discord.gg/JDdSxwjhGf) atau lihat [contoh jendela tanpa bingkai](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/frameless).
