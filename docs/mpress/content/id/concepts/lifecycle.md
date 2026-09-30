---
title: "Siklus Hidup Aplikasi"
description: "Memahami siklus hidup aplikasi Wails dari mulai dijalankan hingga dihentikan"
slug: "concepts/lifecycle"
sourcePath: "concepts/lifecycle.md"
---

## Memahami Siklus Hidup Aplikasi

Aplikasi desktop memiliki siklus hidup dari mulai dijalankan hingga dihentikan. Wails v3 menyediakan **layanan**, **peristiwa**, dan **hook** untuk mengelola siklus hidup ini secara efektif.

## Tahapan Siklus Hidup

```d2
direction: down

Start: Aplikasi Dimulai {
  shape: oval
  style.fill: "#10B981"
}

Init: Inisialisasi {
  Parse: Uraikan Opsi {
    shape: rectangle
  }
  Register: Daftarkan Layanan {
    shape: rectangle
  }
  Setup: Siapkan Runtime {
    shape: rectangle
  }
}

AppRun: app.Run() {
  shape: rectangle
  style.fill: "#3B82F6"
}

ServiceStartup: Layanan Dimulai {
  shape: rectangle
  style.fill: "#8B5CF6"
}

EventLoop: Perulangan Peristiwa {
  Process: Proses Peristiwa {
    shape: rectangle
  }
  Handle: Tangani Pesan {
    shape: rectangle
  }
  Update: Perbarui UI {
    shape: rectangle
  }
}

QuitSignal: Sinyal Keluar {
  shape: diamond
  style.fill: "#F59E0B"
}

ShouldQuit: Pemeriksaan ShouldQuit {
  shape: rectangle
  style.fill: "#3B82F6"
}

OnShutdown: Callback OnShutdown {
  shape: rectangle
  style.fill: "#3B82F6"
}

ServiceShutdown: Penghentian Layanan {
  shape: rectangle
  style.fill: "#8B5CF6"
}

Cleanup: Pembersihan {
  Close: Tutup Jendela {
    shape: rectangle
  }
  Release: Lepaskan Sumber Daya {
    shape: rectangle
  }
}

End: Aplikasi Berakhir {
  shape: oval
  style.fill: "#EF4444"
}

Start -> Init.Parse
Init.Parse -> Init.Register
Init.Register -> Init.Setup
Init.Setup -> AppRun
AppRun -> ServiceStartup
ServiceStartup -> EventLoop.Process
EventLoop.Process -> EventLoop.Handle
EventLoop.Handle -> EventLoop.Update
EventLoop.Update -> EventLoop.Process: Ulangi
EventLoop.Process -> QuitSignal: Pengguna keluar
QuitSignal -> ShouldQuit: Periksa apakah diizinkan
ShouldQuit -> EventLoop.Process: Ditolak
ShouldQuit -> OnShutdown: Diizinkan
OnShutdown -> ServiceShutdown
ServiceShutdown -> Cleanup.Close
Cleanup.Close -> Cleanup.Release
Cleanup.Release -> End
```

### 1. Pembuatan Aplikasi

Buat aplikasi Anda dengan `application.New()`:

```go
app := application.New(application.Options{
    Name:        "My App",
    Description: "An application built with Wails",
    Services: []application.Service{
        application.NewService(&MyService{}),
    },
    Assets: application.AssetOptions{
        Handler: application.BundledAssetFileServer(assets),
    },
})
```

**Yang terjadi:**

1. Opsi diuraikan dan divalidasi
2. Layanan didaftarkan (tetapi belum dimulai)
3. Server aset dikonfigurasi
4. Runtime disiapkan

### 2. Menjalankan Aplikasi

Panggil `app.Run()` untuk memulai aplikasi:

```go
err := app.Run()  // Blocks until quit
if err != nil {
    log.Fatal(err)
}
```

**Yang terjadi:**

1. Layanan dimulai sesuai urutan pendaftaran
2. Listener peristiwa diaktifkan
3. Jendela dapat dibuat
4. Perulangan peristiwa dimulai

### 3. Perulangan Peristiwa

Aplikasi memasuki perulangan peristiwa, tempat aplikasi menghabiskan sebagian besar waktunya:

- Peristiwa OS diproses (peristiwa mouse, papan ketik, dan jendela)
- Pesan dari Go ke JS ditangani
- Panggilan dari JS ke Go dijalankan
- Pembaruan UI dirender

### 4. Penghentian

Saat aplikasi dihentikan:

1. Callback `ShouldQuit` diperiksa (jika ditetapkan)
2. Callback `OnShutdown` dijalankan
3. Layanan dihentikan dalam urutan terbalik
4. Jendela ditutup
5. Sumber daya dilepaskan

## Siklus Hidup Layanan

Layanan merupakan cara utama untuk mengelola siklus hidup di Wails v3. Layanan menyediakan hook yang dipanggil saat proses mulai dan penghentian melalui antarmuka. Untuk dokumentasi lengkap tentang layanan, lihat [panduan Layanan](/features/bindings/services/).

### Membuat Layanan

```go
type MyService struct {
    db *sql.DB
}

// ServiceStartup is called when the application starts
func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    var err error
    s.db, err = sql.Open("sqlite3", "app.db")
    if err != nil {
        return err  // Startup aborts if error returned
    }

    // Run migrations
    if err := s.runMigrations(); err != nil {
        return err
    }

    return nil
}

// ServiceShutdown is called when the application shuts down
func (s *MyService) ServiceShutdown() error {
    if s.db != nil {
        return s.db.Close()
    }
    return nil
}
```

### Mendaftarkan Layanan

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&MyService{}),
        application.NewService(&AnotherService{}),
    },
})
```

**Poin utama:**

- Layanan dimulai sesuai urutan pendaftaran
- Layanan dihentikan dalam urutan pendaftaran **terbalik**
- Jika `ServiceStartup` milik suatu layanan mengembalikan kesalahan, aplikasi dibatalkan
- `ctx` yang diteruskan ke `ServiceStartup` dibatalkan saat proses penghentian dimulai

### Menggunakan Konteks Aplikasi

Konteks yang diteruskan ke `ServiceStartup` berlaku selama masa aktif aplikasi:

```go
func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    // Start a background task that respects shutdown
    go func() {
        ticker := time.NewTicker(5 * time.Minute)
        defer ticker.Stop()

        for {
            select {
            case <-ticker.C:
                s.performBackgroundSync()
            case <-ctx.Done():
                // Application is shutting down
                return
            }
        }
    }()

    return nil
}
```

Anda juga dapat mengakses konteks dari instans aplikasi:

```go
app := application.Get()
ctx := app.Context()
```

## Hook Tingkat Aplikasi

Callback praktis dalam `application.Options` ini memungkinkan Anda mengaitkan kode ke siklus hidup aplikasi tanpa membuat layanan lengkap. Callback ini berguna untuk tugas pembersihan sederhana, konfirmasi keluar, atau saat Anda perlu menjalankan kode pada titik tertentu dalam urutan penonaktifan.

Untuk pengelolaan siklus hidup yang lebih kompleks dengan logika pengaktifan, injeksi dependensi, atau sumber daya berstatus, gunakan [Layanan](#siklus-hidup-layanan) sebagai gantinya.

### ShouldQuit

Callback `ShouldQuit` dipanggil setiap kali ada permintaan keluar—baik karena pengguna menutup jendela terakhir, menekan Cmd+Q (macOS) / Alt+F4 (Windows), maupun memanggil `app.Quit()` secara terprogram.

**Nilai kembalian:**

- Kembalikan `true` untuk mengizinkan proses keluar dilanjutkan (aplikasi akan dinonaktifkan)
- Kembalikan `false` untuk membatalkan proses keluar (aplikasi tetap berjalan)

Pada tahap ini, Anda dapat mencegat dan, bila diperlukan, membatalkan permintaan keluar, misalnya untuk memberi tahu pengguna tentang perubahan yang belum disimpan:

```go
app := application.New(application.Options{
    ShouldQuit: func() bool {
        if !hasUnsavedChanges() {
            return true // No unsaved changes, allow quit
        }

        // Prompt the user — MessageDialog.Show() blocks and returns nothing.
        // The button's OnClick callback fires for whichever button the user picks.
        shouldQuit := false

        dlg := application.Get().Dialog.Question().
            SetTitle("Unsaved Changes").
            SetMessage("You have unsaved changes. Quit anyway?")

        quit := dlg.AddButton("Quit")
        cancel := dlg.AddButton("Cancel")
        dlg.SetDefaultButton(cancel)
        dlg.SetCancelButton(cancel)

        quit.OnClick(func() { shouldQuit = true })

        dlg.Show()
        return shouldQuit
    },
})
```

Jika `ShouldQuit` tidak ditetapkan, aplikasi akan langsung keluar saat diminta.

**Kapan ShouldQuit dipanggil:**

- Pengguna menutup jendela terakhir (kecuali jika `DisableQuitOnLastWindowClosed` ditetapkan)
- Pengguna menekan Cmd+Q di macOS
- Pengguna menekan Alt+F4 di Windows (saat fokus berada pada jendela terakhir)
- Kode memanggil `app.Quit()`

**Kapan ShouldQuit TIDAK dipanggil:**

- Proses dihentikan secara paksa (SIGKILL, penghentian paksa melalui Task Manager)
- `os.Exit()` dipanggil secara langsung

### OnShutdown

Callback `OnShutdown` dipanggil setelah dipastikan bahwa aplikasi akan keluar (setelah `ShouldQuit` mengembalikan `true`, jika ditetapkan). Gunakan callback ini untuk tugas pembersihan seperti menyimpan status, menutup koneksi basis data, atau melepaskan sumber daya.

```go
app := application.New(application.Options{
    OnShutdown: func() {
        // Save application state
        saveState()

        // Close connections
        cleanup()
    },
})
```

Anda juga dapat mendaftarkan callback penonaktifan tambahan secara terprogram kapan saja selama masa aktif aplikasi:

```go
app.OnShutdown(func() {
    log.Println("Application shutting down...")
})
```

Beberapa callback dijalankan sesuai urutan pendaftarannya. Proses penonaktifan akan terblokir hingga semua callback selesai.

**Penting:** Pastikan callback penonaktifan berjalan cepat (kurang dari 1 detik). Sistem operasi dapat menghentikan paksa aplikasi yang memerlukan waktu terlalu lama untuk keluar, sehingga dapat menginterupsi proses pembersihan dan menyebabkan kehilangan data.

### PostShutdown

Callback `PostShutdown` dipanggil setelah semua tugas penonaktifan selesai, tepat sebelum proses berakhir. Pada tahap ini, instans aplikasi tidak dapat digunakan lagi—semua jendela telah ditutup, layanan telah dinonaktifkan, dan sumber daya telah dilepaskan.

Callback ini terutama berguna untuk:

- Pencatatan log akhir yang harus dilakukan setelah semua pembersihan lainnya
- Pengujian dan penelusuran kesalahan perilaku penonaktifan
- Platform tempat `app.Run()` tidak mengembalikan kontrol (callback memastikan kode Anda dijalankan)

```go
app := application.New(application.Options{
    PostShutdown: func() {
        // Final logging
        log.Println("Application terminated cleanly")

        // Flush any buffered logs
        logger.Sync()
    },
})
```

**Catatan:** Jangan mencoba menggunakan fitur aplikasi (jendela, dialog, dan sebagainya) dalam `PostShutdown`—fitur tersebut sudah tidak tersedia.

## Siklus Hidup Berbasis Peristiwa

Wails menyediakan sistem peristiwa yang memberi tahu Anda saat sesuatu terjadi dalam aplikasi—jendela dibuka, aplikasi dimulai, tema berubah, dan lainnya. Anda dapat memantau peristiwa ini untuk merespons perubahan siklus hidup tanpa memblokir atau mencegatnya.

Untuk peristiwa jendela, Anda juga dapat menggunakan `RegisterHook` sebagai pengganti `OnWindowEvent` guna mencegat dan membatalkan tindakan—misalnya, mencegah jendela ditutup. Lihat [Hook Jendela](#hook-jendela-peristiwa-yang-dapat-dibatalkan) di bawah ini.

Untuk dokumentasi lengkap tentang sistem peristiwa, lihat [panduan Peristiwa](/features/events/system/).

### Peristiwa Aplikasi

Pantau peristiwa siklus hidup aplikasi:

```go
app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(event *application.ApplicationEvent) {
    app.Logger.Info("Application has started!")
})
```

Peristiwa khusus platform juga tersedia:

```go
// macOS
app.Event.OnApplicationEvent(events.Mac.ApplicationDidFinishLaunching, func(event *application.ApplicationEvent) {
    // Handle macOS launch
})

app.Event.OnApplicationEvent(events.Mac.ApplicationWillTerminate, func(event *application.ApplicationEvent) {
    // Handle macOS termination
})

// Windows
app.Event.OnApplicationEvent(events.Windows.ApplicationStarted, func(event *application.ApplicationEvent) {
    // Handle Windows start
})
```

### Peristiwa Jendela

Pantau peristiwa siklus hidup jendela:

```go
window := app.Window.New()

window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
    app.Logger.Info("Window gained focus")
})

window.OnWindowEvent(events.Common.WindowClosing, func(e *application.WindowEvent) {
    app.Logger.Info("Window is closing")
})
```

### Hook Jendela (Peristiwa yang Dapat Dibatalkan)

Gunakan `RegisterHook` sebagai pengganti `OnWindowEvent` saat Anda perlu **membatalkan** suatu peristiwa:

```go
window := app.Window.New()

var countdown = 3

window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    countdown--
    if countdown > 0 {
        app.Logger.Info("Not closing yet!", "remaining", countdown)
        e.Cancel()  // Prevent the window from closing
        return
    }
    app.Logger.Info("Window closing now")
})
```

**Perbedaan antara OnWindowEvent dan RegisterHook:**

- `OnWindowEvent`: Memberi tahu Anda saat suatu peristiwa terjadi (tidak dapat dibatalkan)
- `RegisterHook`: Memungkinkan Anda mencegat dan, bila diperlukan, membatalkan peristiwa

## Siklus Hidup Jendela

Jendela memiliki siklus hidupnya sendiri, mulai dari pembuatan hingga penghancuran. Setiap jendela memuat konten frontend-nya secara independen dan dapat ditampilkan, disembunyikan, atau ditutup kapan saja. Saat pengguna mencoba menutup jendela, Anda dapat mencegat tindakan ini dengan `RegisterHook` untuk meminta konfirmasi atau menyembunyikan jendela alih-alih menghancurkannya.

Untuk dokumentasi lengkap tentang jendela, lihat [panduan Jendela](/features/windows/basics/).

```d2
direction: down

Create: Buat Jendela {
  shape: oval
  style.fill: "#10B981"
}

Load: Muat Frontend {
  shape: rectangle
}

Show: Tampilkan Jendela {
  shape: rectangle
}

Active: Jendela Aktif {
  Events: Tangani Peristiwa {
    shape: rectangle
  }
}

CloseRequest: Permintaan Penutupan {
  shape: diamond
  style.fill: "#F59E0B"
}

Hook: Hook WindowClosing {
  shape: rectangle
  style.fill: "#3B82F6"
}

Destroy: Hancurkan Jendela {
  shape: rectangle
}

End: Jendela Ditutup {
  shape: oval
  style.fill: "#EF4444"
}

Create -> Load
Load -> Show
Show -> Active.Events
Active.Events -> Active.Events: Ulangi
Active.Events -> CloseRequest: Pengguna menutup
CloseRequest -> Hook
Hook -> Active.Events: Dibatalkan
Hook -> Destroy: Diizinkan
Destroy -> End
```

### Membuat Jendela

```go
window := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "My Window",
    Width:  800,
    Height: 600,
})
```

### Mencegah Jendela Ditutup

```go
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    if !hasUnsavedChanges() {
        return
    }

    // MessageDialog.Show() returns nothing; per-button OnClick handlers fire.
    dlg := application.Get().Dialog.Question().
        SetTitle("Unsaved Changes").
        SetMessage("Save before closing?")

    save := dlg.AddButton("Save")
    discard := dlg.AddButton("Discard")
    cancel := dlg.AddButton("Cancel")
    dlg.SetDefaultButton(save)
    dlg.SetCancelButton(cancel)

    save.OnClick(func() { saveChanges() })
    cancel.OnClick(func() { e.Cancel() }) // Prevent close
    _ = discard                            // "Discard" falls through and allows close

    dlg.Show()
})
```

### Menyembunyikan Alih-alih Menutup

Pola yang umum untuk aplikasi baki sistem:

```go
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    window.Hide()  // Hide instead of destroy
    e.Cancel()     // Prevent actual close
})
```

## Siklus Hidup Multi-Jendela

Dengan beberapa jendela:

```go
mainWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title: "Main Window",
})

settingsWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "Settings",
    Width:  400,
    Height: 600,
    Hidden: true,  // Start hidden
})
```

**Perilaku default berbeda-beda menurut platform:**

| Platform | Perilaku default saat jendela terakhir ditutup |
| --- | --- |
| macOS | Aplikasi tetap berjalan (bilah menu tetap ada) |
| Windows | Aplikasi berhenti |
| Linux | Aplikasi berhenti |

macOS mengikuti konvensi bawaan platform, yaitu aplikasi biasanya tetap aktif di bilah menu meskipun tidak ada jendela. Secara default, aplikasi berhenti di Windows dan Linux.

**Buat aplikasi berhenti di semua platform saat jendela terakhir ditutup:**

```go
app := application.New(application.Options{
    Mac: application.MacOptions{
        ApplicationShouldTerminateAfterLastWindowClosed: true,
    },
})
```

**Buat aplikasi tetap berjalan di semua platform saat jendela terakhir ditutup:**

Ini berguna untuk aplikasi baki sistem atau aplikasi yang sebaiknya tetap berjalan di latar belakang.

```go
app := application.New(application.Options{
    Windows: application.WindowsOptions{
        DisableQuitOnLastWindowClosed: true,
    },
    Linux: application.LinuxOptions{
        DisableQuitOnLastWindowClosed: true,
    },
})
```

## Pola Umum

### Pola 1: Layanan Basis Data

```go
type DatabaseService struct {
    db *sql.DB
}

func (s *DatabaseService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    var err error
    s.db, err = sql.Open("sqlite3", "app.db")
    if err != nil {
        return fmt.Errorf("failed to open database: %w", err)
    }

    if err := s.db.PingContext(ctx); err != nil {
        return fmt.Errorf("failed to connect to database: %w", err)
    }

    return nil
}

func (s *DatabaseService) ServiceShutdown() error {
    if s.db != nil {
        return s.db.Close()
    }
    return nil
}

// Exported methods are available to the frontend
func (s *DatabaseService) GetUsers() ([]User, error) {
    // Query implementation
}
```

### Pola 2: Layanan Konfigurasi

```go
type ConfigService struct {
    config *Config
    path   string
}

func (s *ConfigService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    s.path = "config.json"

    data, err := os.ReadFile(s.path)
    if err != nil {
        if os.IsNotExist(err) {
            s.config = &Config{} // Default config
            return nil
        }
        return err
    }

    return json.Unmarshal(data, &s.config)
}

func (s *ConfigService) ServiceShutdown() error {
    data, err := json.MarshalIndent(s.config, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(s.path, data, 0644)
}
```

### Pola 3: Pekerja Latar Belakang

```go
type WorkerService struct {
    cancel context.CancelFunc
}

func (s *WorkerService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    workerCtx, cancel := context.WithCancel(ctx)
    s.cancel = cancel

    go s.runWorker(workerCtx)

    return nil
}

func (s *WorkerService) ServiceShutdown() error {
    if s.cancel != nil {
        s.cancel()
    }
    return nil
}

func (s *WorkerService) runWorker(ctx context.Context) {
    ticker := time.NewTicker(time.Minute)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            s.doWork()
        case <-ctx.Done():
            return
        }
    }
}
```

## Referensi Siklus Hidup

| Hook/Antarmuka | Waktu Dipanggil | Dapat Membatalkan? | Kegunaan |
| --- | --- | --- | --- |
| `ServiceStartup` | Selama `app.Run()`, sebelum perulangan peristiwa | Tidak (kembalikan galat untuk membatalkan) | Inisialisasi |
| `ServiceShutdown` | Selama proses penghentian, setelah `OnShutdown` | Tidak | Pembersihan |
| `OnShutdown` | Saat keluar dikonfirmasi | Tidak | Pembersihan aplikasi |
| `ShouldQuit` | Saat keluar diminta | Ya (kembalikan false) | Konfirmasi keluar |
| `RegisterHook(WindowClosing)` | Saat penutupan jendela diminta | Ya (`e.Cancel()`) | Mencegah jendela ditutup |
| `OnWindowEvent` | Saat peristiwa terjadi | Tidak | Menanggapi peristiwa |
| `OnApplicationEvent` | Saat peristiwa terjadi | Tidak | Menanggapi peristiwa |

## Perbedaan Platform

### macOS

- **Menu aplikasi** tetap tersedia meskipun tidak ada jendela
- **Cmd+Q** memicu aplikasi untuk berhenti (melalui `ShouldQuit`)
- **Ikon Dock** tetap ditampilkan kecuali disembunyikan
- Gunakan `ApplicationShouldTerminateAfterLastWindowClosed` untuk mengendalikan perilaku saat aplikasi berhenti

### Windows

- **Tidak ada menu aplikasi** tanpa jendela
- **Alt+F4** menutup jendela (dapat dicegah dengan `RegisterHook`)
- **Baki sistem** dapat membuat aplikasi tetap berjalan

### Linux

- **Perilaku bervariasi** menurut lingkungan desktop
- **Secara umum mirip dengan Windows**

## Men-debug Masalah Siklus Hidup

### Masalah: Aplikasi Tidak Dapat Berhenti

**Penyebab:**

1. `ShouldQuit` mengembalikan `false`
2. `OnShutdown` memerlukan waktu terlalu lama
3. Goroutine latar belakang tidak berhenti

**Solusi:**

```go
// 1. Check ShouldQuit logic
ShouldQuit: func() bool {
    log.Println("ShouldQuit called")
    return true
}

// 2. Keep OnShutdown fast
OnShutdown: func() {
    log.Println("OnShutdown started")
    // Fast cleanup only
    log.Println("OnShutdown finished")
}

// 3. Use context for background tasks
func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    go func() {
        <-ctx.Done()
        log.Println("Context cancelled, stopping background work")
    }()
    return nil
}
```

### Masalah: Proses Memulai Layanan Gagal

**Solusi:** Kembalikan pesan kesalahan yang deskriptif:

```go
func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    if err := s.init(); err != nil {
        return fmt.Errorf("failed to initialise: %w", err)
    }
    return nil
}
```

Kesalahan akan dicatat dan aplikasi tidak akan dimulai.

## Praktik Terbaik

### Lakukan

- **Gunakan layanan untuk mengelola siklus hidup** — Layanan menyediakan hook yang tepat untuk memulai dan menghentikan aplikasi
- **Pastikan proses penghentian berlangsung cepat** — Targetkan kurang dari 1 detik untuk seluruh pembersihan
- **Gunakan context untuk pembatalan** — Hentikan tugas latar belakang dengan benar
- **Tangani kesalahan saat memulai aplikasi** — Kembalikan kesalahan agar proses dapat dibatalkan dengan aman
- **Catat peristiwa siklus hidup** — Hal ini membantu proses debugging

### Jangan Lakukan

- **Jangan lakukan operasi yang memblokir saat layanan dimulai** — Pastikan inisialisasi berlangsung cepat (kurang dari 2 detik)
- **Jangan tampilkan dialog saat aplikasi dihentikan** — Aplikasi sedang berhenti sehingga UI mungkin tidak berfungsi
- **Jangan abaikan context** — Selalu periksa `ctx.Done()` dalam goroutine
- **Jangan biarkan sumber daya bocor** — Selalu implementasikan `ServiceShutdown`

## Langkah Berikutnya

**Layanan** — Pelajari lebih lanjut tentang sistem layanan [Pelajari Selengkapnya →](/features/bindings/services/)

**Sistem Peristiwa** — Gunakan peristiwa untuk berkomunikasi [Pelajari Selengkapnya →](/features/events/system/)

**Pengelolaan Jendela** — Buat dan kelola jendela [Pelajari Selengkapnya →](/features/windows/basics/)

---

**Ada pertanyaan tentang siklus hidup?** Tanyakan di [Discord](https://discord.gg/JDdSxwjhGf) atau lihat [contoh](https://github.com/wailsapp/wails/tree/master/v3/examples).
