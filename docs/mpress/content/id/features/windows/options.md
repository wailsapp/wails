---
title: "Opsi Jendela"
description: "Opsi jendela, nilai default, dan batasan platform"
slug: "features/windows/options"
sourcePath: "features/windows/options.md"
---

## Opsi Konfigurasi Jendela

Wails menyediakan konfigurasi jendela yang lengkap dengan puluhan opsi untuk ukuran, posisi, tampilan, dan perilaku. Referensi ini menjelaskan bidang-bidang umum `WebviewWindowOptions` di Windows, macOS, dan Linux. Struktur di bawah ini merupakan ringkasan bidang yang dipilih; lihat kode sumber Go yang ditautkan untuk mengetahui tipe lengkapnya.

## Struktur WebviewWindowOptions

```go
type WebviewWindowOptions struct {
    // Identity
    Name  string
    Title string

    // Size and Position
    Width           int
    Height          int
    X               int
    Y               int
    MinWidth        int
    MinHeight       int
    MaxWidth        int
    MaxHeight       int
    InitialPosition WindowStartPosition // WindowCentered (default) or WindowXY
    Screen          *Screen             // target screen for initial placement

    // Initial State
    Hidden        bool
    Frameless     bool
    DisableResize bool        // inverted vs v2's `Resizable`
    AlwaysOnTop   bool
    StartState    WindowState // WindowStateNormal | Minimised | Maximised | Fullscreen

    // Appearance
    BackgroundColour RGBA
    BackgroundType   BackgroundType
    Zoom             float64
    ZoomControlEnabled bool

    // Content
    URL  string
    HTML string
    JS   string
    CSS  string

    // Behaviour
    EnableFileDrop              bool
    IgnoreMouseEvents           bool
    HideOnFocusLost             bool
    HideOnEscape                bool
    DevToolsEnabled             bool
    DefaultContextMenuDisabled  bool
    ContentProtectionEnabled    bool
    KeyBindings                 map[string]func(window Window)

    // Permissions
    Permissions map[PermissionType]Permission

    // Window-control button states
    MinimiseButtonState ButtonState
    MaximiseButtonState ButtonState
    CloseButtonState    ButtonState

    // Menu
    UseApplicationMenu bool

    // Platform-specific (per-window)
    Mac     MacWindow
    Windows WindowsWindow
    Linux   LinuxWindow
}
```

`WebviewWindowOptions` tidak memiliki bidang `Parent` — untuk hubungan induk/modal, gunakan `parentWindow.AttachModal(childWindow)`. Struktur ini juga tidak memiliki bidang `Assets` — konfigurasi aset berada di `application.Options` (`Assets AssetOptions`).

Kode sumber lengkap: [`v3/pkg/application/webview_window_options.go`](https://github.com/wailsapp/wails/blob/releases/v3.0.0/v3/pkg/application/webview_window_options.go).

## Opsi Inti

### Nama

**Tipe:** `string` **Default:** `window-<numeric ID>` **Platform:** Semua

```go
Name: "main-window"
```

**Tujuan:** Pengidentifikasi unik untuk menemukan jendela di kemudian hari.

**Praktik terbaik:**

- Gunakan nama deskriptif: `"main"`, `"settings"`, `"about"`
- Gunakan kebab-case: `"file-browser"`, `"color-picker"`
- Buat tetap singkat dan mudah diingat

**Contoh:**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name: "settings-window",
})

// Later...
if settings, ok := app.Window.GetByName("settings-window"); ok {
    settings.Focus()
}
```

### Judul

**Tipe:** `string` **Default:** Nama aplikasi **Platform:** Semua

```go
Title: "My Application"
```

**Tujuan:** Teks yang ditampilkan pada bilah judul dan bilah tugas.

**Pembaruan dinamis:**

```go
window.SetTitle("My Application - Document.txt")
```

### Lebar / Tinggi

**Tipe:** `int` (piksel) **Default:** 800 x 600 **Platform:** Semua **Batasan:** Harus positif

```go
Width:  1200,
Height: 800,
```

**Tujuan:** Ukuran awal jendela dalam piksel logis.

**Catatan:**

- Wails menangani penskalaan DPI secara otomatis
- Gunakan piksel logis, bukan piksel fisik
- Pertimbangkan resolusi layar minimum (1024x768)

**Contoh ukuran:**

| Kasus Penggunaan | Lebar | Tinggi |
| --- | --- | --- |
| Utilitas kecil | 400 | 300 |
| Aplikasi standar | 1024 | 768 |
| Aplikasi besar | 1440 | 900 |
| Full HD | 1920 | 1080 |

### X / Y

**Tipe:** `int` (piksel) **Default:** Berada di tengah layar **Platform:** Semua

```go
X: 100,  // 100px from left edge
Y: 100,  // 100px from top edge
```

**Tujuan:** Posisi awal jendela.

**Sistem koordinat:**

- (0, 0) adalah sudut kiri atas layar utama
- Nilai X positif mengarah ke kanan
- Nilai Y positif mengarah ke bawah

**Contoh:**

`X` dan `Y` hanya berlaku jika `InitialPosition: application.WindowXY` ditetapkan. Jika tidak, nilai default `InitialPosition` adalah `WindowCentered` dan `X`/`Y` diabaikan.

```go
settings := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:            "coordinate-window",
    InitialPosition: application.WindowXY, // opt into X/Y coordinates
    X:               100,
    Y:               100,
})
```

**Praktik terbaik:** Gunakan `Center()` untuk menempatkan jendela di tengah setelah dibuat jika Anda tidak memerlukan koordinat tertentu:

```go
window := app.Window.New()
window.Center()
```

### Lebar Minimum / Tinggi Minimum

**Tipe:** `int` (piksel) **Default:** 0 (tanpa nilai minimum) **Platform:** Semua

```go
MinWidth:  400,
MinHeight: 300,
```

**Tujuan:** Mencegah jendela menjadi terlalu kecil.

**Kasus penggunaan:**

- Mencegah tata letak rusak
- Memastikan kemudahan penggunaan
- Mempertahankan rasio aspek

**Contoh:**

```go
// Prevent window smaller than 400x300
MinWidth:  400,
MinHeight: 300,
```

### Lebar Maksimum / Tinggi Maksimum

**Tipe:** `int` (piksel) **Default:** 0 (tanpa nilai maksimum) **Platform:** Semua

```go
MaxWidth:  1920,
MaxHeight: 1080,
```

**Tujuan:** Mencegah jendela menjadi terlalu besar.

**Kasus penggunaan:**

- Aplikasi berukuran tetap
- Mencegah penggunaan sumber daya yang berlebihan
- Mempertahankan batasan desain

## Opsi Status

### Tersembunyi

**Tipe:** `bool` **Default:** `false` **Platform:** Semua

```go
Hidden: true,
```

**Tujuan:** Membuat jendela tanpa menampilkannya.

**Kasus penggunaan:**

- Jendela latar belakang
- Jendela yang ditampilkan sesuai kebutuhan
- Layar pembuka (buat, muat, lalu tampilkan)
- Mencegah kilatan putih saat memuat konten

**Peningkatan pada platform:**

- **Windows:** Kilatan putih pada jendela telah diperbaiki — jendela tetap tidak terlihat hingga `Show()` dipanggil
- **macOS:** Dukungan penuh
- **Linux:** Dukungan penuh

**Pola yang disarankan untuk pemuatan yang mulus:**

```go
// Create hidden window
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:             "main-window",
    Hidden:           true,
    BackgroundColour: application.NewRGB(30, 30, 30), // Match your theme
})

// Load content while hidden
// ... content loads ...

// Show when ready (no flash!)
window.Show()
```

**Contoh:**

```go
settings := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:   "settings",
    Hidden: true,
})

// Show when needed
settings.Show()
```

### Tanpa Bingkai

**Tipe:** `bool` **Default:** `false` **Platform:** Semua

```go
Frameless: true,
```

**Tujuan:** Menghapus bilah judul dan batas jendela.

**Kasus penggunaan:**

- Chrome jendela khusus
- Layar pembuka
- Aplikasi kios
- Jendela dengan desain khusus

**Penting:** Anda perlu mengimplementasikan:

- Penyeretan jendela
- Tombol tutup/minimalkan/maksimalkan
- Handel pengubahan ukuran (jika ukuran dapat diubah)

**Lihat [Jendela Tanpa Bingkai](/features/windows/frameless/) untuk detailnya.**

### DisableResize

**Tipe:** `bool` **Default:** `false` (ukuran jendela dapat diubah secara default) **Platform:** Semua

```go
DisableResize: true,
```

**Tujuan:** Mencegah perubahan ukuran jendela. Perhatikan bahwa kolom ini merupakan **kebalikan** dari `Resizable` milik v2 — tetapkan `DisableResize: true` agar ukuran jendela tidak dapat diubah.

**Kasus penggunaan:**

- Aplikasi berukuran tetap
- Layar pembuka
- Dialog

**Catatan:** Pengguna tetap dapat memaksimalkan jendela atau mengaktifkan layar penuh, kecuali Anda juga menonaktifkannya melalui `MaximiseButtonState` / perilaku koleksi.

### AlwaysOnTop

**Tipe:** `bool` **Default:** `false` **Platform:** Semua

```go
AlwaysOnTop: true,
```

**Tujuan:** Menjaga jendela tetap berada di atas semua jendela lainnya.

**Kasus penggunaan:**

- Bilah alat mengambang
- Notifikasi
- Gambar-dalam-gambar
- Pengatur waktu

**Catatan platform:**

- **macOS:** Dukungan penuh
- **Windows:** Dukungan penuh
- **Linux:** Bergantung pada pengelola jendela

### StartState

**Tipe:** enum `WindowState` **Default:** `WindowStateNormal` **Platform:** Semua

```go
StartState: application.WindowStateMaximised,
```

**Tujuan:** Status awal jendela saat ditampilkan.

**Nilai:**

- `WindowStateNormal` - Jendela normal
- `WindowStateMinimised` - Diminimalkan
- `WindowStateMaximised` - Dimaksimalkan
- `WindowStateFullscreen` - Layar penuh

Tidak ada konstanta `WindowStateHidden` — gunakan kolom boolean `Hidden` agar jendela tidak terlihat saat dimulai.

**Alihkan mode layar penuh saat runtime:**

```go
window.Fullscreen()
window.UnFullscreen()
window.ToggleFullscreen() // there is no SetFullscreen(bool)
```

## Opsi Tampilan

### BackgroundColour

**Tipe:** struct `RGBA` **Default:** Putih **Platform:** Semua

```go
BackgroundColour: application.RGBA{Red: 0, Green: 0, Blue: 0, Alpha: 255},
```

Kolom `RGBA` bertipe `Red, Green, Blue, Alpha` (uint8). Sebaiknya gunakan fungsi pembantu `application.NewRGB(r, g, b)` (alfa 255) atau `application.NewRGBA(r, g, b, a)`.

**Tujuan:** Warna latar belakang jendela sebelum konten dimuat.

**Kasus penggunaan:**

- Menyesuaikan dengan tema aplikasi Anda
- Mencegah kilatan putih pada tema gelap
- Pengalaman pemuatan yang mulus

**Contoh:**

```go
// Dark theme
BackgroundColour: application.NewRGB(30, 30, 30),

// Light theme
BackgroundColour: application.NewRGB(255, 255, 255),
```

**Metode pembantu:**

```go
window.SetBackgroundColour(application.NewRGB(30, 30, 30))
```

### BackgroundType

**Tipe:** enum `BackgroundType` **Default:** `BackgroundTypeSolid` **Platform:** macOS, Windows (sebagian)

```go
BackgroundType: application.BackgroundTypeTranslucent,
```

**Nilai:**

- `BackgroundTypeSolid` - Warna solid
- `BackgroundTypeTransparent` - Sepenuhnya transparan
- `BackgroundTypeTranslucent` - Efek buram semi-transparan

**Dukungan platform:**

- **macOS:** Konfigurasikan `Mac.Backdrop`; transparansi webview memerlukan [`private_mac_apis`](/guides/build/private-macos-apis/#webview-transparency-and-background). Tanpanya, webview tetap buram.
- **Windows:** Transparan dan translusen (Windows 11+)
- **Linux:** Hanya solid

**Contoh (macOS):**

```go
BackgroundType: application.BackgroundTypeTranslucent,
Mac: application.MacWindow{
    Backdrop: application.MacBackdropTranslucent,
},
```

### OpenInspectorOnStartup dan OpenDevTools

**API privat di macOS:** `OpenInspectorOnStartup: true`, Go `window.OpenDevTools()`, dan JavaScript `Window.OpenDevTools()` memerlukan `private_mac_apis` untuk membuka inspector secara terprogram. Tanpanya, operasi ini tidak melakukan apa pun. Build produksi juga memerlukan `devtools`. Inspeksi Safari publik di macOS 13.3+ tidak memerlukan API privat; pengaktifan inspector pada macOS versi lama memerlukannya. Lihat [matriks build Web Inspector](/guides/build/private-macos-apis/#web-inspector).

## Opsi Konten

### URL

**Tipe:** `string` **Default:** Kosong (dimuat dari Assets) **Platform:** Semua

```go
URL: "https://example.com",
```

**Tujuan:** Memuat URL eksternal sebagai pengganti aset tersemat.

**Kasus penggunaan:**

- Pengembangan (memuat dari server pengembangan)
- Aplikasi berbasis web
- Aplikasi hibrida

**Contoh:**

```go
// Development — point the window at the Vite dev server
URL: "http://localhost:9245",

// Production — embedded assets are configured at the application level
// (Assets is application.Options.Assets, not a WebviewWindowOptions field).
```

Cuplikan tingkat aplikasi untuk kasus produksi:

```go
app := application.New(application.Options{
    Name: "My App",
    Assets: application.AssetOptions{
        Handler: application.AssetFileServerFS(assets),
    },
})
```

### HTML

**Tipe:** `string` **Default:** Kosong **Platform:** Semua

```go
HTML: "<h1>Hello World</h1>",
```

**Tujuan:** Memuat string HTML secara langsung.

**Kasus penggunaan:**

- Jendela sederhana
- Konten yang dihasilkan
- Pengujian

**Contoh:**

```go
HTML: `
<!DOCTYPE html>
<html>
<head><title>Simple Window</title></head>
<body><h1>Hello from Wails!</h1></body>
</html>
`,
```

### Aset (khusus tingkat aplikasi)

Konfigurasi aset **bukan** merupakan field `WebviewWindowOptions`. Aset frontend disajikan oleh aplikasi itu sendiri melalui `application.Options.Assets` (`AssetOptions`); setiap jendela mewarisi server aset tersebut.

```go
app := application.New(application.Options{
    Assets: application.AssetOptions{
        Handler: application.AssetFileServerFS(assets),
    },
})
```

**Lihat [Sistem Build](/concepts/build-system/) untuk detailnya.**

### UseApplicationMenu

**Tipe:** `bool` **Nilai default:** `false` **Platform:** Windows, Linux (tidak berpengaruh pada macOS)

```go
UseApplicationMenu: true,
```

**Tujuan:** Menggunakan menu aplikasi (ditetapkan melalui `app.Menu.Set()`) untuk jendela ini.

Pada **macOS**, opsi ini tidak berpengaruh karena macOS selalu menggunakan menu aplikasi global di bagian atas layar.

Pada **Windows** dan **Linux**, jendela secara default tidak menampilkan menu. Menetapkan `UseApplicationMenu: true` akan menginstruksikan jendela untuk menggunakan menu tingkat aplikasi, sehingga menyediakan solusi lintas platform yang sederhana.

**Contoh:**

```go
// Set the application menu once
menu := app.NewMenu()
menu.AddRole(application.FileMenu)
menu.AddRole(application.EditMenu)
app.Menu.Set(menu)

// All windows with UseApplicationMenu will display this menu
app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:              "Main Window",
    UseApplicationMenu: true,
})

app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:              "Second Window",
    UseApplicationMenu: true,  // Also gets the app menu
})
```

**Catatan:**

- Jika `UseApplicationMenu` dan menu khusus jendela sama-sama ditetapkan, menu khusus jendela akan diprioritaskan
- Ini menyederhanakan kode lintas platform karena tidak memerlukan pemeriksaan OS saat runtime
- Lihat [Menu Aplikasi](/features/menus/application/) untuk dokumentasi menu yang lengkap

## Opsi Input

### EnableFileDrop

**Tipe:** `bool` **Nilai default:** `false` **Platform:** Semua

```go
EnableFileDrop: true,
```

**Tujuan:** Mengaktifkan fitur seret dan lepas untuk file dari sistem operasi ke dalam jendela.

Jika diaktifkan:

- File yang diseret dari pengelola file dapat dilepaskan ke dalam aplikasi Anda
- Peristiwa `WindowFilesDropped` dipicu dengan path file yang dilepaskan
- Elemen dengan atribut `data-file-drop-target` menyediakan informasi mendetail tentang operasi pelepasan

**Kasus penggunaan:**

- Antarmuka unggah file
- Editor dokumen
- Pengimpor media
- Aplikasi apa pun yang menerima file

**Contoh:**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:          "File Uploader",
    EnableFileDrop: true,
})

// Handle dropped files
window.OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {
    files := event.Context().DroppedFiles()
    details := event.Context().DropTargetDetails()
    
    for _, file := range files {
        fmt.Println("Dropped:", file)
    }
})
```

**Zona pelepasan HTML:**

```html
<!-- Mark elements as drop targets -->
<div id="upload" data-file-drop-target>
    Drop files here
</div>
```

**Lihat [Pelepasan File](/features/drag-and-drop/files/) untuk dokumentasi lengkap.**

## Opsi Keamanan

### ContentProtectionEnabled

**Tipe:** `bool` **Nilai default:** `false` **Platform:** Windows (10+), macOS

```go
ContentProtectionEnabled: true,
```

**Tujuan:** Mencegah perekaman konten jendela dari layar.

**Dukungan platform:**

- **Windows:** Windows 10 build 19041+ (penuh), versi yang lebih lama (sebagian)
- **macOS:** Dukungan penuh
- **Linux:** Tidak didukung

**Kasus penggunaan:**

- Aplikasi perbankan
- Pengelola kata sandi
- Rekam medis
- Dokumen rahasia

**Catatan penting:**

1. Tidak mencegah pemotretan fisik
2. Beberapa alat mungkin dapat melewati perlindungan
3. Merupakan bagian dari keamanan menyeluruh, bukan satu-satunya perlindungan
4. Jendela DevTools tidak dilindungi secara otomatis

**Contoh:**

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Secure Window",
    ContentProtectionEnabled: true,
})

// Toggle at runtime
window.SetContentProtection(true)
```

### Izin

**Tipe:** `map[PermissionType]Permission` **Nilai default:** `nil` (penanganan default platform) **Platform:** Linux, Windows (macOS menyerahkannya kepada TCC)

```go
Permissions: map[application.PermissionType]application.Permission{
    application.PermissionMicrophone: application.PermissionAllow,
    application.PermissionCamera:     application.PermissionAllow,
},
```

**Tujuan:** Mengontrol secara deklaratif cara penanganan permintaan kapabilitas (kamera, mikrofon, geolokasi, notifikasi, pembacaan clipboard) dari konten web jendela, tanpa kode khusus platform.

**Nilai PermissionType:** `PermissionMicrophone`, `PermissionCamera`, `PermissionGeolocation`, `PermissionNotifications`, `PermissionClipboardRead`

**Nilai izin:**

- `PermissionDefault` (0) — penanganan native platform: permintaan konfirmasi OS/WebView2 pada macOS/Windows; di Linux, kamera/mikrofon diizinkan, sedangkan semua yang lain ditolak
- `PermissionAllow` (1) — memberikan izin tanpa meminta konfirmasi (Linux: hanya kamera/mikrofon yang telah diimplementasikan; tipe lainnya tetap ditolak)
- `PermissionDeny` (2) — menolak tanpa meminta konfirmasi

**Penting — Windows:** Sebelum opsi ini tersedia, Wails secara diam-diam memberikan semua kapabilitas WebView2. Sekarang, menetapkan entri apa pun dalam `Permissions` akan menonaktifkan pemberian izin menyeluruh tersebut. Kapabilitas yang tidak tercantum akan menampilkan permintaan konfirmasi native WebView2 dan tidak lagi diizinkan secara otomatis. Cantumkan secara eksplisit setiap kapabilitas yang diperlukan aplikasi Anda.

**Lihat [Izin](/features/windows/permissions/) untuk panduan lengkap, matriks platform, dan contoh.**

## Peristiwa Siklus Hidup Jendela

Peristiwa siklus hidup jendela ditangani menggunakan `OnWindowEvent` dan `RegisterHook`. Metode ini memberikan kontrol terperinci atas perilaku penutupan dan pemusnahan jendela.

### Membatalkan Penutupan Jendela

Untuk mencegah jendela ditutup (misalnya, karena ada perubahan yang belum disimpan), gunakan `RegisterHook` dengan peristiwa `WindowClosing`, lalu panggil `event.Cancel()`:

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Name:  "main-window",
    Title: "My Application",
})

// Register a hook to intercept the closing event
window.RegisterHook(events.Common.WindowClosing, func(event *application.WindowEvent) {
    if hasUnsavedChanges() {
        // Ask user for confirmation
        result := showConfirmDialog("Unsaved changes. Close anyway?")
        if result != "yes" {
            // Cancel the close event
            event.Cancel()
        }
    }
})
```

**Poin penting:**

- `RegisterHook` mencegat peristiwa sebelum terjadi
- Panggil `event.Cancel()` untuk mencegah jendela ditutup
- Jendela akan tetap terbuka setelah penutupan dibatalkan

### Menangani Penutupan Jendela

Untuk melakukan pembersihan saat jendela ditutup, gunakan `OnWindowEvent` dengan peristiwa `WindowClosing`:

```go
window.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
    // Cleanup code runs here
    fmt.Printf("Window %s is closing\n", window.Name())

    // Close database connection
    if db != nil {
        db.Close()
    }

    // Remove from window list
    removeWindow(window.ID())
})
```

**Poin penting:**

- `OnWindowEvent` menangani peristiwa yang akan segera terjadi
- Pembersihan dijalankan sebelum jendela dimusnahkan
- Penutupan tidak dapat dibatalkan dari sini (gunakan `RegisterHook` untuk melakukannya)

### Pola Pembersihan Jendela Singleton

Untuk jendela singleton (memastikan hanya ada satu instans), gunakan `WindowClosing` untuk membersihkan referensinya:

```go
var settingsWindow *application.WebviewWindow

func ShowSettings(app *application.App) {
    // Create if doesn't exist
    if settingsWindow == nil {
        settingsWindow = app.Window.NewWithOptions(application.WebviewWindowOptions{
            Name:  "settings",
            Title: "Settings",
            Width: 600,
            Height: 400,
        })

        // Cleanup on close
        settingsWindow.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
            settingsWindow = nil
        })
    }

    // Show and focus
    settingsWindow.Show()
    settingsWindow.Focus()
}
```

## Opsi Khusus Platform

### Opsi Mac

```go
Mac: application.MacWindow{
    TitleBar: application.MacTitleBar{
        AppearsTransparent: true,
        Hide:               false,
        HideTitle:          true,
        FullSizeContent:    true,
    },
    Backdrop:                application.MacBackdropTranslucent,
    InvisibleTitleBarHeight: 50,
    WindowClass:             application.MacWindowClassPanel,
    PanelPreferences: application.MacPanelPreferences{
        NonActivating:          true,
        FloatingPanel:          true,
        BecomesKeyOnlyIfNeeded: false,
        UtilityWindow:          false,
    },
    WindowLevel:             application.MacWindowLevelFloating,
    CollectionBehavior:      application.MacWindowCollectionBehaviorDefault,
    TabbingMode:             application.MacWindowTabbingModeDisallowed,
},
```

**TitleBar** (`MacTitleBar`)

- `AppearsTransparent` - Membuat bilah judul transparan sehingga konten meluas ke area bilah judul
- `Hide` - Menyembunyikan bilah judul sepenuhnya
- `HideTitle` - Hanya menyembunyikan teks judul
- `FullSizeContent` - Memperluas konten hingga memenuhi seluruh ukuran jendela

Di macOS, transparansi webview dan pembukaan inspector secara terprogram memerlukan tag build `private_mac_apis`. Tanpa tag tersebut, opsi yang sama tetap valid, tetapi operasi yang hanya menggunakan API privat tidak melakukan apa pun. Pengelompokan Liquid Glass juga diabaikan dan gaya menggunakan alternatif API publik. Lihat [API macOS privat](/guides/build/private-macos-apis/) untuk mengetahui perintah build dan perilaku yang tepat.

**Backdrop** (`MacBackdrop`)

- `MacBackdropNormal` - Latar belakang opak standar
- `MacBackdropTranslucent` - **API privat diperlukan untuk transparansi webview.** Tanpa tag tersebut, efek buram native tetap berada di belakang webview yang opak.
- `MacBackdropTransparent` - **API privat diperlukan untuk transparansi webview.** Tanpa tag tersebut, webview tetap opak.
- `MacBackdropLiquidGlass` - **API privat diperlukan untuk transparansi webview.** Tanpa tag tersebut, lapisan kaca tetap berada di belakang webview yang opak; gaya menggunakan alternatif API publik.

**LiquidGlass** (`MacLiquidGlass`)

| Bidang atau nilai | Ketergantungan pada API privat di macOS |
| --- | --- |
| `Style: LiquidGlassStyleAutomatic` | Gaya reguler native bersifat publik; transparansi webview pada latar belakang memerlukan `private_mac_apis`. |
| `Style: LiquidGlassStyleLight` | Tag tersebut mempertahankan pemetaan gaya transparan native yang sudah ada; tanpa tag itu, Wails menggunakan kaca reguler dengan tampilan Aqua. |
| `Style: LiquidGlassStyleDark` | **API privat:** nilai gaya native yang tidak terdokumentasi `2`; tanpa tag tersebut, kaca reguler dengan tampilan Dark Aqua. |
| `Style: LiquidGlassStyleVibrant` | Pemetaan gaya transparan native bersifat publik; transparansi webview pada latar belakang memerlukan tag tersebut. |
| `GroupID` | **API privat:** nilai yang tidak kosong meminta pengelompokan; diabaikan tanpa tag tersebut. |
| `GroupSpacing` | **API privat:** nilai positif meminta jarak antarkelompok; diabaikan tanpa tag tersebut. |
| `Material`, `CornerRadius`, `TintColor` | Tidak memiliki ketergantungan sendiri pada API privat. |

Lihat [Nilai Liquid Glass](/guides/build/private-macos-apis/#liquid-glass-values) untuk mengetahui nilai gaya native dan ketersediaannya pada sistem operasi.

**InvisibleTitleBarHeight** (`int`)

- Tinggi area bilah judul tak terlihat (untuk menyeret jendela)
- Hanya berlaku ketika area seret bilah judul native disembunyikan—yaitu ketika jendela tidak memiliki bingkai (`Frameless: true`) atau menggunakan bilah judul transparan (`AppearsTransparent: true`)
- Tidak berpengaruh pada jendela standar dengan bilah judul yang terlihat

**WindowClass** (`MacWindowClass`)

- `MacWindowClassWindow` - Perilaku `NSWindow` standar (default)
- `MacWindowClassPanel` - `NSPanel` tambahan yang tidak pernah menjadi jendela utama aplikasi

`PanelPreferences` hanya berlaku untuk `MacWindowClassPanel`:

- `NonActivating` menambahkan `NSWindowStyleMaskNonactivatingPanel`. Menampilkan atau memfokuskan panel tidak mengaktifkan aplikasi Wails, tetapi panel tetap dapat menjadi jendela key untuk kontrol dan input teks.
- `FloatingPanel` mengaktifkan perilaku panel mengambang AppKit.
- `BecomesKeyOnlyIfNeeded` hanya memperoleh status key ketika tampilan yang diklik meminta input papan ketik.
- `UtilityWindow` menerapkan gaya jendela utilitas native.

Panel Wails tetap terlihat ketika aplikasi dinonaktifkan dan dilepas saat ditutup, sesuai dengan siklus hidup yang diharapkan oleh `WebviewWindow`. Perilaku ini sengaja menggantikan default `NSPanel` yang berlawanan.

Kelas jendela, level, kebijakan aktivasi, dan perilaku koleksi menyelesaikan masalah yang berbeda:

- `WindowClass` memilih `NSWindow` atau `NSPanel` serta mengontrol semantik jendela utama/key.
- `WindowLevel` mengontrol urutan z. Gunakan `MacWindowLevelPopUpMenu` untuk overlay bilah menu.
- `MacOptions.ActivationPolicy` mengontrol aplikasi secara keseluruhan, termasuk penyajian Dock dan bilah menu. Panel yang tidak mengaktifkan aplikasi tidak memerlukan kebijakan aktivasi aksesori, meskipun aplikasi tetap dapat menggunakannya untuk menyembunyikan ikon Dock.
- `CollectionBehavior` mengontrol partisipasi dalam Spaces dan mode layar penuh.

```go
// Spotlight/menu-bar panel that leaves the current application active.
Mac: application.MacWindow{
    WindowClass: application.MacWindowClassPanel,
    PanelPreferences: application.MacPanelPreferences{
        NonActivating: true,
    },
    WindowLevel: application.MacWindowLevelPopUpMenu,
    CollectionBehavior: application.MacWindowCollectionBehaviorCanJoinAllSpaces |
        application.MacWindowCollectionBehaviorFullScreenAuxiliary |
        application.MacWindowCollectionBehaviorStationary,
},
```

**WindowLevel** (`MacWindowLevel`)

- `MacWindowLevelNormal` - Level jendela standar (default)
- `MacWindowLevelFloating` - Mengambang di atas jendela normal
- `MacWindowLevelTornOffMenu` - Level menu yang dilepas
- `MacWindowLevelModalPanel` - Level panel modal
- `MacWindowLevelMainMenu` - Level menu utama
- `MacWindowLevelStatus` - Level jendela status
- `MacWindowLevelPopUpMenu` - Level menu pop-up
- `MacWindowLevelScreenSaver` - Level penghemat layar

`WindowLevel` yang ditetapkan secara eksplisit lebih diutamakan daripada `AlwaysOnTop` dan `PanelPreferences.FloatingPanel`. Tanpa level eksplisit, `AlwaysOnTop` atau panel mengambang akan menghasilkan `MacWindowLevelFloating`; jika tidak, levelnya adalah `MacWindowLevelNormal`. Pemanggilan `SetAlwaysOnTop` setelahnya tetap merupakan perubahan runtime eksplisit.

**CollectionBehavior** (`MacWindowCollectionBehavior`)

Mengontrol perilaku jendela di berbagai macOS Spaces dan mode layar penuh. Nilai-nilai ini berupa bitmask yang dapat digabungkan menggunakan OR bitwise (`|`).

**Perilaku Space:**

- `MacWindowCollectionBehaviorDefault` - Menggunakan FullScreenPrimary (default, kompatibel dengan versi sebelumnya)
- `MacWindowCollectionBehaviorCanJoinAllSpaces` - Jendela muncul di semua Spaces
- `MacWindowCollectionBehaviorMoveToActiveSpace` - Berpindah ke Space aktif ketika ditampilkan
- `MacWindowCollectionBehaviorManaged` - Perilaku default untuk jendela terkelola
- `MacWindowCollectionBehaviorTransient` - Jendela sementara/transien
- `MacWindowCollectionBehaviorStationary` - Tetap diam saat beralih Space

**Peralihan siklus jendela:**

- `MacWindowCollectionBehaviorParticipatesInCycle` - Disertakan dalam pergantian dengan Cmd+`
- `MacWindowCollectionBehaviorIgnoresCycle` - Tidak disertakan dalam pergantian dengan Cmd+`

**Perilaku layar penuh:**

- `MacWindowCollectionBehaviorFullScreenPrimary` - Dapat memasuki mode layar penuh
- `MacWindowCollectionBehaviorFullScreenAuxiliary` - Dapat ditampilkan di atas aplikasi layar penuh
- `MacWindowCollectionBehaviorFullScreenNone` - Menonaktifkan kemampuan layar penuh
- `MacWindowCollectionBehaviorFullScreenAllowsTiling` - Mengizinkan penyusunan berdampingan (macOS 10.11+)
- `MacWindowCollectionBehaviorFullScreenDisallowsTiling` - Mencegah penyusunan (macOS 10.11+)

**Contoh - Jendela seperti Spotlight:**

```go
// Window that appears on all Spaces AND can overlay fullscreen apps
Mac: application.MacWindow{
	WindowClass: application.MacWindowClassPanel,
	PanelPreferences: application.MacPanelPreferences{
		NonActivating: true,
	},
    CollectionBehavior: application.MacWindowCollectionBehaviorCanJoinAllSpaces |
                        application.MacWindowCollectionBehaviorFullScreenAuxiliary,
    WindowLevel:        application.MacWindowLevelFloating,
},
```

**Contoh - Perilaku tunggal:**

```go
// Window that can appear over fullscreen applications
Mac: application.MacWindow{
    CollectionBehavior: application.MacWindowCollectionBehaviorFullScreenAuxiliary,
},
```

**TabbingMode** (`MacWindowTabbingMode`)

Mengontrol perilaku tab jendela pada macOS 10.12 dan versi yang lebih baru. Tab jendela memungkinkan beberapa jendela dikelompokkan sebagai tab.

**Opsi:**

- `MacWindowTabbingModeDefault` - Nilai sentinel nol (tidak ditetapkan secara eksplisit). Saat runtime, nilai bakunya adalah tidak mengizinkan tab
- `MacWindowTabbingModeAutomatic` - Sistem menentukan perilaku tab
- `MacWindowTabbingModePreferred` - Jendela lebih memilih berada dalam mode tab
- `MacWindowTabbingModeDisallowed` - Menonaktifkan tab jendela

**Contoh - Menonaktifkan tab jendela:**

```go
Mac: application.MacWindow{
    TabbingMode: application.MacWindowTabbingModeDisallowed,
},
```

**Contoh - Mengutamakan tab jendela:**

```go
Mac: application.MacWindow{
    TabbingMode: application.MacWindowTabbingModePreferred,
},
```

**WebviewPreferences** (`MacWebviewPreferences`)

Kontrol terperinci atas konfigurasi `WKWebView` yang mendasarinya. Semua bidang bersifat opsional—bidang yang tidak ditetapkan tidak mengubah nilai baku WebKit.

```go
Mac: application.MacWindow{
    WebviewPreferences: application.MacWebviewPreferences{
        TabFocusesLinks:                       optional.True,
        TextInteractionEnabled:                optional.True,
        FullscreenEnabled:                     optional.False,
        AllowsBackForwardNavigationGestures:   optional.False,
        AllowsMagnification:                   optional.True,
        AllowsAirPlayForMediaPlayback:         optional.True,
        JavaScriptCanOpenWindowsAutomatically: optional.False,
        MinimumFontSize:                       optional.NewVar(12.0),
        ApplicationNameForUserAgent:           "MyApp",
        EnableAutoplayWithoutUserAction:       optional.True,
    },
},
```

- `TabFocusesLinks` — Jika `true`, menekan Tab memindahkan fokus ke tautan dan kontrol formulir (nilai baku: `false`)
- `TextInteractionEnabled` — Jika `true`, pengguna dapat memilih dan berinteraksi dengan teks dalam webview (nilai baku: `true`)
- `FullscreenEnabled` — Jika `true`, konten web dapat memasuki layar penuh melalui HTML Fullscreen API (nilai baku: `false`). Memerlukan macOS 12.3+.
- `AllowsBackForwardNavigationGestures` — Jika `true`, gestur usap horizontal memicu navigasi mundur/maju (nilai baku: `false`)
- `AllowsMagnification` — Jika `true`, pembesaran dengan gestur cubit diaktifkan pada webview (nilai baku: `false`)
- `AllowsAirPlayForMediaPlayback` — Jika `true`, media dapat dialirkan ke perangkat AirPlay (nilai baku: `true`)
- `JavaScriptCanOpenWindowsAutomatically` — Jika `true`, JavaScript dapat membuka jendela baru tanpa gestur pengguna (nilai baku: `false`)
- `MinimumFontSize` — Ukuran fon minimum dalam poin. Gunakan `optional.NewVar(12.0)` untuk menetapkannya. Jika tidak ditetapkan, nilai baku WebKit tetap digunakan.
- `ApplicationNameForUserAgent` — Mengganti sufiks nama aplikasi dalam string user-agent WebKit. Berguna ketika situs menolak pengenal bawaan `"wails.io"` (misalnya sematan YouTube). Biarkan kosong untuk mempertahankan nilai baku.
- `EnableAutoplayWithoutUserAction` — Jika `true`, audio dan video dapat diputar otomatis tanpa gestur pengguna. Dipetakan ke `WKWebViewConfiguration.mediaTypesRequiringUserActionForPlayback = WKAudiovisualMediaTypeNone` (nilai baku: `false`)

### Opsi Windows (per jendela)

Struct per jendela adalah `application.WindowsWindow`—**bukan** `WindowsOptions` (yang terakhir adalah struct tingkat *aplikasi*).

```go
Windows: application.WindowsWindow{
    DisableIcon:                       false,
    DisableMenu:                       false,
    BackdropType:                      application.Auto,
    CustomTheme:                       application.ThemeSettings{},
    DisableFramelessWindowDecorations: false,
    NonClientRegionSupport:            false,
    WebView2CompositionHosting:        false,
},
```

**DisableIcon** (`bool`)

- Menghapus ikon dari bilah judul.

**DisableMenu** (`bool`)

- Menonaktifkan bilah menu untuk jendela. Jika `true`, jendela tidak menampilkan bilah menu meskipun bilah tersebut telah dikonfigurasi.
- Nilai baku: `false`

**BackdropType** (`BackdropType`)

- `application.Auto` - Nilai baku sistem
- `application.None` - Tanpa latar belakang
- `application.Mica` - Material Mica (Windows 11)
- `application.Acrylic` - Material Acrylic (Windows 11)
- `application.Tabbed` - Material Tabbed (Windows 11)

Tidak ada konstanta bergaya `WindowsBackdropTypeMica`—gunakan `application.Mica` dan seterusnya.

**CustomTheme** (`ThemeSettings`)

- Nilai (bukan pointer). Warna mode gelap/terang khusus untuk bingkai jendela, teks/latar belakang bilah judul, dan bilah menu.

**DisableFramelessWindowDecorations** (`bool`)

- Menonaktifkan dekorasi bawaan untuk jendela tanpa bingkai (bayangan Aero, sudut membulat).

**NonClientRegionSupport** (`bool`)

- Mengaktifkan dukungan `app-region: drag` / `app-region: no-drag` bawaan WebView2 untuk bilah judul khusus tanpa bingkai.
- Opsi ini hanya untuk penyeretan aplikasi native sederhana. Opsi ini tidak menyediakan perilaku native untuk tombol caption khusus maupun Windows 11 Snap Assist / Snap Layouts pada tombol maksimalkan khusus.

**WebView2CompositionHosting** (`bool`)

- Mengaktifkan dukungan `--wails-non-client-region` yang dikelola Wails untuk tombol caption khusus dengan perilaku native Windows, termasuk Windows 11 Snap Assist / Snap Layouts pada tombol maksimalkan khusus.
- Eksperimental. Opsi ini meng-host WebView2 melalui `ICoreWebView2CompositionController` dan DirectComposition sebagai pengganti pengontrol bawaan yang di-host HWND.
- Dapat digabungkan dengan `NonClientRegionSupport` jika sebuah jendela memerlukan dukungan `app-region` bawaan WebView2 sekaligus area tombol caption khusus yang dikelola Wails.

**Contoh:**

```go
Windows: application.WindowsWindow{
    BackdropType: application.Mica,
    DisableIcon:  true,
},
```

**Contoh - Area bilah judul Windows khusus:**

```go
Windows: application.WindowsWindow{
    NonClientRegionSupport:    true,
    WebView2CompositionHosting: true,
},
```

Lihat [Windows Tanpa Bingkai](/features/windows/frameless/#native-non-client-regions-on-windows) untuk mengetahui detail perilaku, kompromi, dan CSS yang sesuai.

### Opsi Linux (per jendela)

Struct per jendela adalah `application.LinuxWindow`—**bukan** `LinuxOptions`.

```go
Linux: application.LinuxWindow{
    Icon:                []byte{/* PNG data */},
    WindowIsTranslucent: false,
},
```

**Icon** (`[]byte`)

- Ikon jendela (format PNG).

**WindowIsTranslucent** (`bool`)

- Memerlukan dukungan compositor.

**Contoh:**

```go
//go:embed icon.png
var icon []byte

Linux: application.LinuxWindow{
    Icon: icon,
},
```

## Opsi Windows Tingkat Aplikasi

Beberapa opsi khusus Windows harus dikonfigurasi pada tingkat aplikasi, bukan per jendela. Hal ini karena WebView2 menggunakan satu lingkungan browser bersama untuk setiap jalur data pengguna.

### Flag Browser

Flag browser WebView2 mengontrol fitur eksperimental dan perilaku di **semua jendela** dalam aplikasi Anda. Flag ini harus ditetapkan di `application.Options.Windows`:

```go
app := application.New(application.Options{
    Name: "My App",
    Windows: application.WindowsOptions{
        // Enable experimental WebView2 features
        EnabledFeatures: []string{
            "msWebView2EnableDraggableRegions",
        },

        // Disable specific features
        DisabledFeatures: []string{
            "msSmartScreenProtection",  // Always disabled by Wails
        },

        // Additional Chromium command-line arguments
        AdditionalBrowserArgs: []string{
            "--disable-gpu",
            "--remote-debugging-port=9222",
        },
    },
})
```

**EnabledFeatures** (`[]string`)

- Daftar flag fitur WebView2 yang akan diaktifkan
- Lihat [flag browser WebView2](https://learn.microsoft.com/en-us/microsoft-edge/webview2/concepts/webview-features-flags) untuk mengetahui flag yang tersedia
- Contoh: `"msWebView2EnableDraggableRegions"`

**DisabledFeatures** (`[]string`)

- Daftar flag fitur WebView2 yang akan dinonaktifkan
- Wails secara otomatis menonaktifkan `msSmartScreenProtection`
- Contoh: `"msExperimentalFeature"`

**AdditionalBrowserArgs** (`[]string`)

- Argumen baris perintah Chromium yang diteruskan ke proses browser
- Harus menyertakan prefiks `--` (misalnya, `"--remote-debugging-port=9222"`)
- Lihat [sakelar baris perintah Chromium](https://peter.sh/experiments/chromium-command-line-switches/) untuk mengetahui argumen yang tersedia

@note{type="caution" title="Penting"}
Flag ini berlaku secara global untuk SEMUA jendela karena WebView2 menggunakan satu lingkungan browser bersama untuk setiap jalur data pengguna. Anda tidak dapat menggunakan flag browser yang berbeda untuk jendela yang berbeda.

@end

**Contoh Lengkap:**

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "My App",
        Windows: application.WindowsOptions{
            // Enable draggable regions feature
            EnabledFeatures: []string{
                "msWebView2EnableDraggableRegions",
            },
            // Enable remote debugging
            AdditionalBrowserArgs: []string{
                "--remote-debugging-port=9222",
            },
        },
    })

    // All windows will use the browser flags configured above
    window := app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title:  "Main Window",
        Width:  1024,
        Height: 768,
    })

    window.Show()
    app.Run()
}
```

## Contoh Lengkap

Berikut adalah konfigurasi jendela yang siap digunakan dalam produksi:

```go
package main

import (
    _ "embed"
    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var assets embed.FS

//go:embed icon.png
var icon []byte

func main() {
    app := application.New(application.Options{
        Name: "My Application",
    })

    window := app.Window.NewWithOptions(application.WebviewWindowOptions{
        // Identity
        Name:  "main-window",
        Title: "My Application",

        // Size and Position
        Width:     1200,
        Height:    800,
        MinWidth:  800,
        MinHeight: 600,

        // Initial State
        StartState: application.WindowStateNormal,

        // Appearance
        BackgroundColour: application.NewRGB(255, 255, 255),

        // Platform-Specific (per-window structs)
        Mac: application.MacWindow{
            TitleBar: application.MacTitleBar{
                AppearsTransparent: true,
            },
            Backdrop: application.MacBackdropTranslucent,
        },

        Windows: application.WindowsWindow{
            BackdropType: application.Mica,
            DisableIcon:  false,
        },

        Linux: application.LinuxWindow{
            Icon: icon,
        },
    })

    window.Center()
    window.Show()

    app.Run()
}
```

Aset frontend disajikan di tingkat aplikasi (`application.Options.Assets`), bukan per jendela.

## Langkah Berikutnya

- [Dasar-Dasar Jendela](/features/windows/basics/) - Membuat dan mengontrol jendela
- [Beberapa Jendela](/features/windows/multiple/) - Pola multijendela
- [Jendela Tanpa Bingkai](/features/windows/frameless/) - Bingkai jendela kustom
- [Peristiwa Jendela](/features/windows/events/) - Peristiwa siklus hidup

---

**Ada pertanyaan?** Tanyakan di [Discord](https://discord.gg/JDdSxwjhGf) atau lihat [contoh](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples).
