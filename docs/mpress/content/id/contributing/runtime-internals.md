---
title: "Internal Runtime"
description: "Pembahasan mendalam tentang cara Wails v3 melakukan boot, berjalan, dan berkomunikasi dengan OS"
slug: "contributing/runtime-internals"
sourcePath: "contributing/runtime-internals.md"
---

**runtime** adalah lapisan yang mengubah fungsi Go biasa menjadi aplikasi desktop lintas platform. Dokumen ini menjelaskan komponen-komponen yang akan Anda temui saat menelusuri kode sumber.

---

## 1. Siklus Hidup Aplikasi

| Fase | Jalur Kode | Yang Terjadi |
| --- | --- | --- |
| **Bootstrap** | `pkg/application/application.go:init()` | Mendaftarkan data waktu build dan membuat singleton `application` global. |
| **New()** | `application.New(...)` | Memvalidasi `Options`, menjalankan **AssetServer**, dan menginisialisasi pencatatan log. |
| **Run()** | `application.(*App).Run()` | 1. Memanggil `mainthread.X()` platform untuk memasuki thread UI OS.<br />2. Melakukan boot pada **runtime** (`internal/runtime`).<br />3. Memblokir hingga jendela terakhir ditutup atau `Quit()` dipanggil. |
| **Shutdown** | `application.(*App).Quit()` | Menyiarkan peristiwa `application:shutdown`, menulis entri log yang masih berada di buffer ke tujuan log, lalu menghentikan jendela dan layanan. |

Siklus hidup ini secara ketat bersifat **sekali masuk**: Anda dapat membuat banyak jendela, tetapi objek aplikasi itu sendiri hanya diinisialisasi sekali.

---

## 2. Pengelolaan Jendela

### API Publik

```go
win := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "Dashboard",
    Width:  1280,
    Height: 720,
})
win.Show()
```

> `app.Window.New()` **tidak menerima argumen**; gunakan `NewWithOptions(...)` jika Anda
>
> perlu meneruskan struct `application.WebviewWindowOptions` (berdasarkan nilai).

`app.Window.New[WithOptions]()` mendelegasikan ke `pkg/application/webview_window_*.go`, tempat implementasi khusus platform berada:

```
pkg/application/
├── webview_window_darwin.go    // WKWebView
├── webview_window_linux.go     // GTK + WebKitGTK (plus linux_cgo*.go)
└── webview_window_windows.go   // WebView2
```

Setiap file:

1. Membuat webview native (WKWebView, WebKitGTK, WebView2).
2. Mendaftarkan callback **Pemroses Pesan** (`pkg/application/messageprocessor*.go`).
3. Memetakan peristiwa Wails (`WindowDidResize`, `WindowFocus`, `WindowFilesDropped`, …) ke konstanta dalam `pkg/events`.

`internal/runtime/` dicadangkan untuk kode penghubung kecil dengan build tag (`runtime{,_darwin,_linux,_windows,_android,_dev,_prod}.go`) dan runtime JS tertanam di bawah `internal/runtime/desktop/`.

Jendela aktif dilacak oleh `pkg/application/window_manager.go` / `webview_window.go`. `pkg/application/screenmanager.go` digunakan untuk metadata **layar** (resolusi, skala, area kerja)—bukan untuk mengelola jendela.

---

## 3. Alur Pemrosesan Pesan

Jembatan antara JavaScript dan Go diimplementasikan oleh keluarga **Pemroses Pesan** dalam `pkg/application/messageprocessor_*.go`.

Alur:

1. **JavaScript** memanggil `Call.ByID(<fnv-id>, ...args)` dari `/wails/runtime.js` (diimplementasikan dalam `internal/runtime/desktop/@wailsio/runtime/src/calls.ts`)—atau `Call.ByName("pkg.Struct.Method", ...args)` untuk build mode nama.
2. Helper runtime mengemas panggilan tersebut dan mengirimkannya melalui jembatan native khusus platform ke Go.
3. **Go** menerima pesan tersebut dalam `pkg/application/messageprocessor_call.go`.
4. Pemroses mencari metode terikat dalam `pkg/application/bindings.go` (ditulis manual dan berbasis `reflect`), lalu memanggilnya.
5. Hasil atau galat dimarshalkan kembali ke JS, tempat `Promise` diselesaikan atau ditolak.

> Envelope JSON yang tepat ditentukan oleh helper runtime di sisi JS dan
>
> oleh `messageprocessor_call.go` di sisi Go—draf lama halaman ini
>
> mencantumkan bentuk `{"t":"c","id":"123","m":"Greet","p":[…]}`, tetapi bentuk tersebut tidak
>
> sesuai dengan implementasi saat ini. Baca kedua file secara bersamaan saat menelusuri
>
> bug format data dalam komunikasi.

Pemroses khusus:

| File | Tujuan |
| --- | --- |
| `messageprocessor_window.go` | Tindakan jendela (sembunyikan, maksimalkan, …) |
| `messageprocessor_dialog.go` | Dialog native (`OpenFile`, `MessageBox`, …) |
| `messageprocessor_clipboard.go` | Membaca/menulis papan klip |
| `messageprocessor_events.go` | Berlangganan/memancarkan peristiwa |
| `messageprocessor_browser.go` | Navigasi browser, alat pengembang |

Pemroses bersifat **tanpa status**—semua yang dibutuhkan diambil dari `ApplicationContext` yang diteruskan bersama setiap pesan.

---

## 4. Sistem Peristiwa

Peristiwa berupa string dengan namespace yang dikirimkan melalui tiga lapisan:

1. **Peristiwa aplikasi**: siklus hidup global (`application:ready`, `application:shutdown`).
2. **Peristiwa jendela**: per jendela (`window:focus`, `window:resize`).
3. **Peristiwa kustom**: ditentukan pengguna (`chat:new-message`).

Detail implementasi:

- Konstanta peristiwa berada di `pkg/events/` (`defaults.go`, `known_events.go`, `events.txt`). Konstanta tersebut dihasilkan oleh `v3/tasks/events/generate.go` dan diekspos sebagai `events.Common.*`, `events.Mac.*`, `events.Windows.*`, `events.Linux.*`. Anda dapat membuatnya ulang dengan `wails3 generate constants`.
- Sisi Go (peristiwa aplikasi):
  ```go
  app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(e *application.ApplicationEvent) {})
  ```

- Sisi Go (peristiwa jendela):
  ```go
  window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {})
  ```

- Sisi Go (peristiwa kustom):
  ```go
  app.Event.On("chat:new-message", func(e *application.CustomEvent) {})
  ```

- Sisi JS:
  ```js
  import { Events } from "/wails/runtime.js";
  Events.On("chat:new-message", (e) => { /* … */ });
  ```


Semua peristiwa aplikasi, jendela, dan kustom mengalir melalui `pkg/application/event_manager.go`. Langganan peristiwa jendela dibatasi pada jendelanya, sehingga menutup jendela akan membatalkan pendaftaran handler-nya secara otomatis.

---

## 5. Implementasi Khusus Platform

Kompilasi bersyarat menjaga API publik tetap identik sekaligus menyembunyikan perbedaan OS.

| Aspek | Darwin | Linux | Windows |
| --- | --- | --- | --- |
| Thread Utama | `mainthread_darwin.go` (Cgo ke Foundation) | `mainthread_linux.go` (GTK) | `mainthread_windows.go` (Win32 `AttachThreadInput`) |
| Dialog | `dialogs_darwin.*` (NSAlert) | `dialogs_linux.go` (GtkFileChooser) | `dialogs_windows.go` (IFileOpenDialog) |
| Papan klip | `clipboard_darwin.go` | `clipboard_linux.go` | `clipboard_windows.go` |
| Ikon Baki Sistem | `systemtray_darwin.*` | `systemtray_linux.go` (DBus) | `systemtray_windows.go` (Shell_NotifyIcon) |

Prinsip utama:

- **macOS dan Windows** menggunakan Cgo secara terbatas (terutama melalui `pkg/mac/` dan pembungkus Win32 `w32` di `pkg/w32`).
- **Linux banyak menggunakan Cgo** karena kebutuhan—`pkg/application/linux_cgo.go` (~69 KB) ditambah `linux_cgo_gtk4.{c,go,h}` (~50 KB+) mengendalikan GTK/WebKitGTK secara langsung.
- Gunakan **tag build** (`//go:build darwin`, `//go:build linux`, …) agar berkas untuk setiap OS tetap mudah dibaca.
- `internal/capabilities/` tersedia untuk flag kapabilitas per platform, tetapi framework **tidak** mengekspor sentinel `ErrCapability`—pembatasan fitur dilakukan melalui nilai kembalian stub khusus platform.

---

## 6. Panduan Berkas

| Berkas | Alasan Anda Perlu Mengubahnya |
| --- | --- |
| `internal/runtime/runtime_*.go` | Ubah lapisan stub kecil bertag build (pengembangan vs produksi, kode penghubung khusus OS). |
| `pkg/application/webview_window_*.go` | Implementasikan petunjuk atau perilaku jendela baru. |
| `pkg/application/messageprocessor*.go` | Tambahkan perintah bridge baru yang dapat dipanggil dari JS. |
| `pkg/events/*.go` | Perluas definisi peristiwa bawaan (lalu jalankan kembali `wails3 generate constants`). |
| `internal/assetserver/*` | Sesuaikan penanganan aset untuk pengembangan/produksi. |
| `internal/runtime/desktop/@wailsio/runtime/src/*` | Edit runtime JS tersemat (pengiriman panggilan/peristiwa, dialog, seret, …). |

---

## 7. Kiat Debugging

- Konfigurasikan `Options.LogLevel` (misalnya `slog.LevelDebug`) dan periksa output `Options.Logger`—tidak ada variabel lingkungan `WAILS_LOG_LEVEL`.
- Flag `wails3 dev` adalah `--config`, `--port`, `-s` (mengaktifkan HTTPS), dan `--no-colour` global. Tidak ada flag `-verbose`.
- Di macOS, jalankan di bawah `lldb --` untuk mendeteksi pengecualian Objective-C sejak dini.
- Untuk masalah Chromium di Windows, aktifkan log debug WebView2: `set WEBVIEW2_ADDITIONAL_BROWSER_ARGUMENTS=--remote-debugging-port=9222`

---

## 8. Memperluas Runtime

1. Opsional: deklarasikan setiap flag kapabilitas baru di `internal/capabilities/`.
2. Implementasikan fitur di setiap varian `pkg/application/*_{darwin,linux,windows}.go` menggunakan tag build. Sediakan stub pada platform yang tidak mendukungnya.
3. Tambahkan API publik di `pkg/application` (antarmuka + metode konkret `WebviewWindow`, struct opsi, dan sebagainya).
4. Daftarkan metode pemroses pesan baru (`pkg/application/messageprocessor*.go`) dan helper yang sesuai di runtime JS jika JS perlu memanggilnya.
5. Jika Anda menambahkan peristiwa, deklarasikan konstantanya di `pkg/events/` dan jalankan `wails3 generate constants` untuk memperbarui file yang dihasilkan.

Ikuti daftar periksa ini agar kontrak lintas platform tetap utuh.

---

## 9. Seret dan Lepas

Fitur seret dan lepas file menggunakan pendekatan **yang mengutamakan JavaScript** di semua platform. Lapisan native mencegat peristiwa seret dari OS, tetapi penanganan pelepasan dan interaksi DOM yang sebenarnya berlangsung di JavaScript.

### Alur

1. Pengguna menyeret file dari OS ke atas jendela Wails
2. Lapisan native mendeteksi tindakan seret dan memberi tahu JavaScript untuk menampilkan efek saat penunjuk berada di atas target
3. Pengguna melepaskan file
4. Lapisan native mengirimkan jalur file + koordinat ke JavaScript
5. JavaScript menemukan elemen target pelepasan (`data-file-drop-target`)
6. JavaScript mengirimkan jalur file + detail elemen ke backend Go
7. Go memancarkan peristiwa `WindowFilesDropped` beserta konteks lengkap

### Implementasi Platform

| Platform | Lapisan Native | Tantangan Utama |
| --- | --- | --- |
| **Windows** | Dukungan seret bawaan WebView2 | Koordinat dalam piksel CSS, tanpa perlu konversi |
| **macOS** | Delegat seret NSWindow | Konversikan koordinat relatif terhadap jendela menjadi relatif terhadap webview |
| **Linux** | Sinyal seret GTK3 | Harus membedakan penyeretan file dari penyeretan HTML5 internal |

### Linux: Membedakan Jenis Penyeretan

GTK dan WebKit sama-sama ingin menangani peristiwa seret. Kuncinya adalah memeriksa jenis target seret:

```c
static gboolean is_file_drag(GdkDragContext *context) {
    GList *targets = gdk_drag_context_list_targets(context);
    for (GList *l = targets; l != NULL; l = l->next) {
        GdkAtom atom = GDK_POINTER_TO_ATOM(l->data);
        gchar *name = gdk_atom_name(atom);
        if (name && g_strcmp0(name, "text/uri-list") == 0) {
            g_free(name);
            return TRUE;  // External file drag
        }
        g_free(name);
    }
    return FALSE;  // Internal HTML5 drag
}
```

Penangan sinyal mengembalikan `FALSE` untuk penyeretan internal (agar WebKit menanganinya) dan `TRUE` untuk penyeretan file (agar kita menanganinya sendiri).

### Memblokir Pelepasan File

Ketika `EnableFileDrop` bernilai `false`, kita tetap harus mencegah browser membuka file yang dilepaskan. Setiap platform menanganinya secara berbeda:

- **Windows**: JavaScript memanggil `preventDefault()` pada peristiwa seret
- **macOS**: JavaScript memanggil `preventDefault()` pada peristiwa seret\
- **Linux**: Penangan sinyal GTK mencegat dan menolak penyeretan file pada lapisan native

### File Utama

| File | Tujuan |
| --- | --- |
| `pkg/application/linux_cgo.go` | Penangan sinyal seret GTK (kode C dalam preambule cgo) |
| `pkg/application/webview_window_darwin.go` | Delegat seret macOS |
| `pkg/application/webview_window_windows.go` | Penanganan pesan WebView2 |
| `internal/runtime/desktop/@wailsio/runtime/src/window.ts` | Penanganan pelepasan di JavaScript |

### Debugging

- **Linux**: Tambahkan `printf` dalam kode C (ingat `fflush(stdout)`)
- **Windows**: Gunakan `globalApplication.debug()`
- **JavaScript**: Periksa konsol browser, aktifkan mode debug

Masalah umum:

1. **Penyeretan HTML5 internal tidak berfungsi**: Penangan native mencegatnya (kembalikan `FALSE` untuk penyeretan nonfile)
2. **Efek saat penunjuk berada di atas target tidak muncul**: Penangan JavaScript tidak dipanggil
3. **Koordinat salah**: Periksa konversi ruang koordinat

---

Kini Anda telah mengikuti tur terpandu mengenai cara kerja internal runtime. Gabungkan pengetahuan ini dengan peta **Tata Letak Basis Kode** dan dokumentasi **Server Aset** agar Anda dapat menavigasi dengan percaya diri dan memberikan kontribusi yang berdampak. Selamat membuat kode!
