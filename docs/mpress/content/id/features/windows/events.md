---
title: "Peristiwa Jendela"
description: "Menangani peristiwa siklus hidup dan perubahan status jendela"
slug: "features/windows/events"
sourcePath: "features/windows/events.md"
---

## Peristiwa Jendela

Wails mengirimkan peristiwa siklus hidup dan perubahan status jendela melalui satu API: `OnWindowEvent` untuk listener pasif dan `RegisterHook` untuk hook yang dapat dibatalkan guna mencegah tindakan bawaan.

```go
import "github.com/wailsapp/wails/v3/pkg/events"

// Listener — observes the event, cannot cancel it.
window.OnWindowEvent(events.Common.WindowDidResize, func(e *application.WindowEvent) { /* ... */ })

// Hook — runs before listeners; can call e.Cancel() to suppress the default action.
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    if hasUnsavedChanges() {
        e.Cancel() // prevent the window from closing
    }
})
```

Kedua pemanggilan tersebut mengembalikan `unsubscribe func()` yang dapat Anda panggil untuk menghapus handler.

Peristiwa lintas platform berada di `events.Common.*`. Peristiwa khusus platform berada di `events.Mac.*`, `events.Windows.*`, dan `events.Linux.*`. Daftar lengkapnya dibuat di `v3/pkg/events/events.go`.

## Peristiwa Siklus Hidup

### Pembuatan Jendela

Jalankan callback setiap kali jendela dibuat menggunakan `app.Window.OnCreate`:

```go
app.Window.OnCreate(func(window application.Window) {
    fmt.Printf("Window created: %s (ID: %d)\n", window.Name(), window.ID())

    // Configure all new windows
    window.SetMinSize(400, 300)

    // Register a hook for this window
    window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
        if !confirmClose() {
            e.Cancel()
        }
    })
})
```

Callback menerima `application.Window` (sebuah antarmuka). Callback `OnCreate` dipanggil sekali untuk setiap jendela, setelah runtime jendela diinisialisasi.

### WindowClosing

Dipicu saat pengguna mencoba menutup jendela (mengeklik X, menekan ⌘W, Alt+F4, dan sebagainya).

Gunakan **hook** (`RegisterHook`) — hanya hook yang dapat membatalkan penutupan. Listener mengamati peristiwa tersebut, tetapi tidak dapat mencegahnya.

```go
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    if hasUnsavedChanges() {
        // Cancel the close — the window stays open.
        e.Cancel()
    }
})
```

**Penting:**

- `WindowClosing` dikirim saat pengguna mencoba menutup jendela.
- Callback `RegisterHook` dapat memanggil `e.Cancel()` agar jendela tetap terbuka.
- Callback `OnWindowEvent` untuk `WindowClosing` berfungsi sebagai pengamat — callback tersebut dipicu, tetapi tidak dapat membatalkan.

**Penutupan terprogram:** panggil `window.Close()`. Tidak ada `window.Destroy()`.

### WindowRuntimeReady

Dipicu setelah runtime dalam jendela selesai diinisialisasi — ini adalah saat yang aman untuk memanggil konteks JS jendela:

```go
window.OnWindowEvent(events.Common.WindowRuntimeReady, func(e *application.WindowEvent) {
    window.EmitEvent("app-ready", nil)
})
```

## Peristiwa Fokus

### WindowFocus

Dipanggil saat jendela memperoleh fokus:

```go
window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
    fmt.Println("Window gained focus")
    updateTitleBar(true)
    app.Event.Emit("window-focused", window.ID())
})
```

### WindowLostFocus

Dipanggil saat jendela kehilangan fokus:

```go
window.OnWindowEvent(events.Common.WindowLostFocus, func(e *application.WindowEvent) {
    fmt.Println("Window lost focus")
    updateTitleBar(false)
    saveCurrentState()
})
```

**Contoh: UI yang peka terhadap fokus:**

```go
type FocusAwareWindow struct {
    window  *application.WebviewWindow
    focused bool
}

func (fw *FocusAwareWindow) Setup() {
    fw.window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
        fw.focused = true
        fw.updateAppearance()
    })

    fw.window.OnWindowEvent(events.Common.WindowLostFocus, func(e *application.WindowEvent) {
        fw.focused = false
        fw.updateAppearance()
    })
}

func (fw *FocusAwareWindow) updateAppearance() {
    if fw.focused {
        fw.window.EmitEvent("update-theme", "active")
    } else {
        fw.window.EmitEvent("update-theme", "inactive")
    }
}
```

## Peristiwa Perubahan Status

### WindowMinimise / WindowUnMinimise

```go
window.OnWindowEvent(events.Common.WindowMinimise, func(e *application.WindowEvent) {
    pauseRendering()
    saveWindowState()
})

window.OnWindowEvent(events.Common.WindowUnMinimise, func(e *application.WindowEvent) {
    resumeRendering()
    refreshContent()
})
```

### WindowMaximise / WindowUnMaximise

```go
window.OnWindowEvent(events.Common.WindowMaximise, func(e *application.WindowEvent) {
    window.EmitEvent("layout-mode", "maximised")
})

window.OnWindowEvent(events.Common.WindowUnMaximise, func(e *application.WindowEvent) {
    window.EmitEvent("layout-mode", "normal")
})
```

### WindowFullscreen / WindowUnFullscreen

```go
window.OnWindowEvent(events.Common.WindowFullscreen, func(e *application.WindowEvent) {
    window.EmitEvent("chrome-visibility", false)
    window.EmitEvent("layout-mode", "fullscreen")
})

window.OnWindowEvent(events.Common.WindowUnFullscreen, func(e *application.WindowEvent) {
    window.EmitEvent("chrome-visibility", true)
    window.EmitEvent("layout-mode", "normal")
})
```

Masuk/keluar dari layar penuh secara terprogram dengan `window.Fullscreen()` / `window.UnFullscreen()` / `window.ToggleFullscreen()` — tidak ada `SetFullscreen(bool)`. Periksa status dengan `window.IsFullscreen() bool`.

## Peristiwa Posisi dan Ukuran

### WindowDidMove

```go
window.OnWindowEvent(events.Common.WindowDidMove, func(e *application.WindowEvent) {
    x, y := window.Position()
    fmt.Printf("Window moved to: %d, %d\n", x, y)
    saveWindowPosition(x, y)
})
```

### WindowDidResize

```go
window.OnWindowEvent(events.Common.WindowDidResize, func(e *application.WindowEvent) {
    width, height := window.Size()
    fmt.Printf("Window resized to: %dx%d\n", width, height)
    saveWindowSize(width, height)
    window.EmitEvent("window-size", map[string]int{
        "width":  width,
        "height": height,
    })
})
```

`WindowDidResize` dan `WindowDidMove` tidak membawa koordinat pada peristiwanya sendiri — periksa jendela melalui `window.Size()` / `window.Position()` di dalam callback.

**Contoh: Tata letak responsif:**

```go
window.OnWindowEvent(events.Common.WindowDidResize, func(e *application.WindowEvent) {
    width, _ := window.Size()
    var layout string
    switch {
    case width < 600:
        layout = "compact"
    case width < 1200:
        layout = "normal"
    default:
        layout = "wide"
    }
    window.EmitEvent("layout-changed", layout)
})
```

## Contoh Lengkap

Jendela siap produksi dengan penanganan peristiwa lengkap:

```go
package main

import (
    "encoding/json"
    "fmt"
    "os"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)

type WindowState struct {
    X          int  `json:"x"`
    Y          int  `json:"y"`
    Width      int  `json:"width"`
    Height     int  `json:"height"`
    Maximised  bool `json:"maximised"`
    Fullscreen bool `json:"fullscreen"`
}

type ManagedWindow struct {
    app    *application.App
    window *application.WebviewWindow
    state  WindowState
    dirty  bool
}

func main() {
    app := application.New(application.Options{
        Name: "Event Demo",
    })

    mw := &ManagedWindow{app: app}
    mw.CreateWindow()
    mw.LoadState()
    mw.SetupEventHandlers()

    app.Run()
}

func (mw *ManagedWindow) CreateWindow() {
    mw.window = mw.app.Window.NewWithOptions(application.WebviewWindowOptions{
        Name:   "main",
        Title:  "Event Demo",
        Width:  800,
        Height: 600,
    })
}

func (mw *ManagedWindow) SetupEventHandlers() {
    // Focus events
    mw.window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
        mw.window.EmitEvent("focus-state", true)
    })

    mw.window.OnWindowEvent(events.Common.WindowLostFocus, func(e *application.WindowEvent) {
        mw.window.EmitEvent("focus-state", false)
    })

    // State change events
    mw.window.OnWindowEvent(events.Common.WindowMinimise, func(e *application.WindowEvent) {
        mw.SaveState()
    })

    mw.window.OnWindowEvent(events.Common.WindowMaximise, func(e *application.WindowEvent) {
        mw.state.Maximised = true
        mw.dirty = true
    })

    mw.window.OnWindowEvent(events.Common.WindowUnMaximise, func(e *application.WindowEvent) {
        mw.state.Maximised = false
        mw.dirty = true
    })

    mw.window.OnWindowEvent(events.Common.WindowFullscreen, func(e *application.WindowEvent) {
        mw.state.Fullscreen = true
        mw.dirty = true
    })

    mw.window.OnWindowEvent(events.Common.WindowUnFullscreen, func(e *application.WindowEvent) {
        mw.state.Fullscreen = false
        mw.dirty = true
    })

    // Position and size events
    mw.window.OnWindowEvent(events.Common.WindowDidMove, func(e *application.WindowEvent) {
        mw.state.X, mw.state.Y = mw.window.Position()
        mw.dirty = true
    })

    mw.window.OnWindowEvent(events.Common.WindowDidResize, func(e *application.WindowEvent) {
        mw.state.Width, mw.state.Height = mw.window.Size()
        mw.dirty = true
    })

    // Cancellable close — use a hook, not a listener.
    mw.window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
        if mw.dirty {
            mw.SaveState()
        }
    })
}

func (mw *ManagedWindow) LoadState() {
    data, err := os.ReadFile("window-state.json")
    if err != nil {
        return
    }

    if err := json.Unmarshal(data, &mw.state); err != nil {
        return
    }

    // Restore window state
    mw.window.SetPosition(mw.state.X, mw.state.Y)
    mw.window.SetSize(mw.state.Width, mw.state.Height)

    if mw.state.Maximised {
        mw.window.Maximise()
    }

    if mw.state.Fullscreen {
        mw.window.Fullscreen()
    }
}

func (mw *ManagedWindow) SaveState() {
    data, err := json.Marshal(mw.state)
    if err != nil {
        return
    }

    os.WriteFile("window-state.json", data, 0644)
    mw.dirty = false

    fmt.Println("Window state saved")
}
```

## Koordinasi Peristiwa

### Peristiwa Lintas Jendela

Koordinasikan beberapa jendela menggunakan bus peristiwa aplikasi:

```go
// In main window
mainWindow.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
    app.Event.Emit("main-window-focused", nil)
})

// In other windows
app.Event.On("main-window-focused", func(event *application.CustomEvent) {
    updateRelativeToMain()
})
```

### Rantai Peristiwa

```go
window.OnWindowEvent(events.Common.WindowMaximise, func(e *application.WindowEvent) {
    saveWindowState()
    window.EmitEvent("layout-changed", "maximised")
    app.Event.Emit("window-maximised", window.ID())
})
```

### Peristiwa dengan Debounce

```go
var resizeTimer *time.Timer

window.OnWindowEvent(events.Common.WindowDidResize, func(e *application.WindowEvent) {
    if resizeTimer != nil {
        resizeTimer.Stop()
    }

    resizeTimer = time.AfterFunc(500*time.Millisecond, func() {
        w, h := window.Size()
        saveWindowSize(w, h)
    })
})
```

## Praktik Terbaik

### ✅ Lakukan

- **Gunakan hook untuk pembatalan** — hanya callback `RegisterHook` yang dapat memanggil `e.Cancel()`.
- **Simpan status saat jendela ditutup** — pulihkan posisi/ukuran jendela saat aplikasi dijalankan berikutnya.
- **Terapkan debounce pada peristiwa yang sering terjadi** — `WindowDidResize` dan `WindowDidMove` dipicu dengan cepat.
- **Tangani perubahan fokus** — perbarui UI sebagaimana mestinya.
- **Lakukan koordinasi melalui peristiwa** — gunakan `app.Event.Emit` untuk pengiriman pesan lintas jendela.
- **Berhenti berlangganan** saat handler tidak lagi diperlukan — `OnWindowEvent`/`RegisterHook` sama-sama mengembalikan `func()` untuk berhenti berlangganan.

### ❌ Jangan Lakukan

- **Jangan memblokir handler peristiwa** — pastikan handler tetap berjalan cepat.
- **Jangan mencoba membatalkan dari `OnWindowEvent`** — gunakan `RegisterHook`.
- **Jangan gunakan `window.Destroy()`** — itu tidak ada; gunakan `window.Close()`.
- **Jangan menyimpan pada setiap peristiwa** — terapkan debounce terlebih dahulu.
- **Jangan menganggap koordinat tersedia pada peristiwa** — panggil `window.Position()` / `window.Size()`.

## Pemecahan Masalah

### Hook WindowClosing Tidak Mencegah Penutupan

**Penyebab:** Anda menggunakan `OnWindowEvent` alih-alih `RegisterHook`.

**Solusi:** Hanya callback `RegisterHook` yang dapat memanggil `e.Cancel()`. Callback `OnWindowEvent` hanya mengamati; callback tersebut tidak dapat membatalkan.

```go
// ❌ Cannot cancel — this is a listener.
window.OnWindowEvent(events.Common.WindowClosing, func(e *application.WindowEvent) {
    e.Cancel() // no effect
})

// ✅ Can cancel — this is a hook.
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    e.Cancel()
})
```

### Peristiwa Tidak Dipicu

**Penyebab:** Handler didaftarkan setelah peristiwa terjadi.

**Solusi:** Daftarkan handler segera setelah jendela dibuat (atau dalam callback `app.Window.OnCreate`).

```go
window := app.Window.New()
window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) { /* ... */ })
```

### Kebocoran Memori

**Penyebab:** Handler berumur panjang disimpan oleh state berumur pendek.

**Solusi:** Simpan fungsi berhenti berlangganan `func()` yang dikembalikan oleh `OnWindowEvent`/`RegisterHook`, lalu panggil fungsi tersebut saat pembersihan.

```go
unsub := window.OnWindowEvent(events.Common.WindowDidResize, handler)
// ...later:
unsub()
```

## Langkah Berikutnya

**Dasar-Dasar Jendela** - Pelajari dasar-dasar pengelolaan jendela [Pelajari Selengkapnya →](/features/windows/basics/)

**Beberapa Jendela** - Pola untuk aplikasi dengan beberapa jendela [Pelajari Selengkapnya →](/features/windows/multiple/)

**Sistem Peristiwa** - Pelajari sistem peristiwa secara mendalam [Pelajari Selengkapnya →](/features/events/system/)

**Siklus Hidup Aplikasi** - Pahami siklus hidup aplikasi [Pelajari Selengkapnya →](/concepts/lifecycle/)

---

**Ada pertanyaan?** Tanyakan di [Discord](https://discord.gg/JDdSxwjhGf) atau lihat [contoh](https://github.com/wailsapp/wails/tree/master/v3/examples).
