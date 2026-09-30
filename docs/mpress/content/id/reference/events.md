---
title: "API Peristiwa"
description: "Referensi lengkap untuk API Peristiwa"
slug: "reference/events"
sourcePath: "reference/events.md"
---

## Ikhtisar

API Peristiwa menyediakan metode untuk memancarkan dan mendengarkan peristiwa, sehingga memungkinkan komunikasi antara berbagai bagian aplikasi Anda.

**Jenis Peristiwa:**

- **Peristiwa Aplikasi** - Peristiwa siklus hidup aplikasi (dimulai, dihentikan)
- **Peristiwa Jendela** - Perubahan status jendela (mendapatkan fokus, kehilangan fokus, perubahan ukuran)
- **Peristiwa Kustom** - Peristiwa yang ditentukan pengguna untuk komunikasi khusus aplikasi

**Pola Komunikasi:**

- **Go ke Frontend** - Pancarkan peristiwa dari Go, dengarkan di JavaScript
- **Frontend ke Go** - Tidak secara langsung (gunakan binding layanan sebagai gantinya)
- **Frontend ke Frontend** - Melalui Go atau peristiwa runtime lokal
- **Jendela ke Jendela** - Targetkan jendela tertentu atau siarkan ke semua jendela

## Metode Peristiwa (Go)

### app.Event.Emit()

Memancarkan peristiwa kustom ke semua jendela. Mengembalikan `true` jika hook membatalkan pemancaran.

```go
func (em *EventManager) Emit(name string, data ...any) bool
```

**Parameter:**

- `name` - Nama peristiwa
- `data` - Data opsional yang akan dikirim bersama peristiwa

**Contoh:**

```go
// Emit simple event
app.Event.Emit("user-logged-in")

// Emit with data
app.Event.Emit("data-updated", map[string]interface{}{
    "count": 42,
    "status": "success",
})

// Emit multiple values
app.Event.Emit("progress", 75, "Processing files...")
```

### app.Event.On()

Mendengarkan peristiwa kustom di Go.

```go
func (em *EventManager) On(name string, callback func(*CustomEvent)) func()
```

**Parameter:**

- `name` - Nama peristiwa yang akan didengarkan
- `callback` - Fungsi yang dipanggil saat peristiwa dipancarkan

**Nilai kembalian:** Fungsi pembersihan untuk menghapus listener peristiwa

**Contoh:**

```go
// Listen for events
cleanup := app.Event.On("user-action", func(e *application.CustomEvent) {
    data := e.Data.(map[string]interface{})
    action := data["action"].(string)
    app.Logger.Info("User action", "action", action)
})

// Later, remove listener
cleanup()
```

### Peristiwa Khusus Jendela

Pancarkan peristiwa ke jendela tertentu:

```go
// Emit to specific window
window.EmitEvent("notification", "Hello from Go!")

// Emit to all windows
app.Event.Emit("global-update", data)
```

## Metode Peristiwa (Frontend)

### On()

Mendengarkan peristiwa dari Go.

```javascript
import { Events } from '@wailsio/runtime'

Events.On(eventName, callback)
```

**Parameter:**

- `eventName` - Nama peristiwa yang akan didengarkan
- `callback` - Fungsi yang dipanggil saat peristiwa diterima

**Nilai kembalian:** Fungsi pembersihan

**Contoh:**

```javascript
import { Events } from '@wailsio/runtime'

// Listen for events
const cleanup = Events.On('data-updated', (data) => {
    console.log('Count:', data.count)
    console.log('Status:', data.status)
    updateUI(data)
})

// Later, remove listener
cleanup()
```

### Once()

Mendengarkan satu kemunculan peristiwa.

```javascript
import { Events } from '@wailsio/runtime'

Events.Once(eventName, callback)
```

**Contoh:**

```javascript
import { Events } from '@wailsio/runtime'

// Listen for first occurrence only
Events.Once('initialization-complete', (data) => {
    console.log('App initialized!', data)
    // This will only fire once
})
```

### Off()

Menghapus semua listener untuk satu atau beberapa peristiwa. `Off` menerima string nama peristiwa variadik—fungsi ini **tidak** menerima callback. Untuk menghapus satu listener, simpan fungsi berhenti berlangganan yang dikembalikan oleh `Events.On(...)` lalu panggil fungsi tersebut.

```typescript
import { Events } from '@wailsio/runtime'

Events.Off(...eventNames: string[]): void
```

**Contoh:**

```javascript
import { Events } from '@wailsio/runtime'

// Preferred: keep the unsubscribe fn from On()
const unsubscribe = Events.On('my-event', (data) => {
    console.log('Event received:', data)
})

// Later — remove just this listener
unsubscribe()

// Or: remove every listener for one or more events
Events.Off('my-event', 'another-event')
```

### OffAll()

Menghapus **semua** listener peristiwa. Tidak menerima argumen.

```typescript
import { Events } from '@wailsio/runtime'

Events.OffAll(): void
```

**Contoh:**

```javascript
import { Events } from '@wailsio/runtime'

// Remove all listeners — typically used during teardown.
Events.OffAll()
```

### OnMultiple()

Mendengarkan peristiwa hingga `max` kali, lalu berhenti berlangganan secara otomatis.

```typescript
Events.OnMultiple(eventName: string, callback, max: number): () => void
```

**Contoh:**

```javascript
Events.OnMultiple('progress', (data) => {
    console.log('progress', data)
}, 5)
```

## Peristiwa Aplikasi

### app.Event.OnApplicationEvent()

Mendengarkan peristiwa siklus hidup aplikasi.

```go
func (em *EventManager) OnApplicationEvent(
    eventType events.ApplicationEventType,
    callback func(*ApplicationEvent),
) func()
```

Konstanta peristiwa berada dalam paket `events`—`events.Common.*` untuk peristiwa lintas platform, sedangkan `events.Mac.*` / `events.Windows.*` / `events.Linux.*` untuk peristiwa khusus platform.

**Peristiwa aplikasi umum:**

- `events.Common.ApplicationStarted` - Aplikasi telah selesai diluncurkan.
- `events.Common.ThemeChanged` - Tema sistem beralih antara terang dan gelap.
- `events.Common.ApplicationOpenedWithFile` - Diluncurkan melalui asosiasi berkas.
- `events.Common.ApplicationLaunchedWithUrl` - Diluncurkan melalui skema URL.

**Tidak ada** konstanta peristiwa "penghentian aplikasi" generik—daftarkan pembersihan saat aplikasi dihentikan melalui `application.Options.OnShutdown` atau `app.OnShutdown(func())`.

**Contoh:**

```go
import "github.com/wailsapp/wails/v3/pkg/events"

// Handle application startup
app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(e *application.ApplicationEvent) {
    app.Logger.Info("Application started")
})

// React to theme changes
app.Event.OnApplicationEvent(events.Common.ThemeChanged, func(e *application.ApplicationEvent) {
    if e.Context().IsDarkMode() {
        app.Logger.Info("Dark mode enabled")
    }
})

// Shutdown cleanup is configured on the application, not as an event:
app.OnShutdown(func() {
    database.Close()
    saveSettings()
})
```

## Peristiwa Jendela

### OnWindowEvent()

Mendengarkan peristiwa khusus jendela.

```go
func (w *WebviewWindow) OnWindowEvent(
    eventType events.WindowEventType,
    callback func(*WindowEvent),
) func()
```

Konstanta peristiwa jendela berada dalam paket `events`: `events.Common.*` untuk peristiwa lintas platform (dan `events.Mac.*` / `events.Windows.*` / `events.Linux.*` untuk peristiwa khusus platform).

**Peristiwa jendela umum:**

- `events.Common.WindowFocus` - Jendela mendapatkan fokus.
- `events.Common.WindowLostFocus` - Jendela kehilangan fokus.
- `events.Common.WindowClosing` - Jendela akan ditutup (dapat dibatalkan melalui `RegisterHook`).
- `events.Common.WindowDidResize` - Ukuran jendela diubah.
- `events.Common.WindowDidMove` - Jendela dipindahkan.
- `events.Common.WindowMinimise` / `WindowUnMinimise` / `WindowMaximise` / `WindowUnMaximise` / `WindowFullscreen` / `WindowUnFullscreen`.
- `events.Common.WindowRuntimeReady` - Event sudah aman dikirim ke runtime dalam jendela.

**Contoh:**

```go
import "github.com/wailsapp/wails/v3/pkg/events"

// Handle window focus
window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
    app.Logger.Info("Window focused")
})

// Handle window resize
window.OnWindowEvent(events.Common.WindowDidResize, func(e *application.WindowEvent) {
    width, height := window.Size()
    app.Logger.Info("Window resized", "width", width, "height", height)
})
```

## Pola Umum

Pola-pola ini menunjukkan pendekatan yang telah terbukti untuk menggunakan event dalam aplikasi nyata. Setiap pola mengatasi tantangan komunikasi tertentu antara backend Go dan frontend Anda, sehingga membantu Anda membangun aplikasi yang responsif dan terstruktur dengan baik.

### Pola Permintaan/Respons

Gunakan pola ini saat Anda ingin memberi tahu frontend bahwa operasi backend telah selesai, misalnya setelah pengambilan data, pemrosesan berkas, atau tugas latar belakang. Binding layanan mengembalikan data secara langsung, sedangkan event memberikan notifikasi tambahan untuk pembaruan UI, seperti menampilkan pesan toast atau memuat ulang daftar.

**Go:**

```go
// Service method
type DataService struct {
    app *application.App
}

func (s *DataService) FetchData(query string) ([]Item, error) {
    items := fetchFromDatabase(query)

    // Emit event when done
    s.app.Event.Emit("data-fetched", map[string]interface{}{
        "query": query,
        "count": len(items),
    })

    return items, nil
}
```

**JavaScript:**

```javascript
import { FetchData } from './bindings/DataService'
import { Events } from '@wailsio/runtime'

// Listen for completion event
Events.On('data-fetched', (data) => {
    console.log(`Fetched ${data.count} items for query: ${data.query}`)
    showNotification(`Found ${data.count} results`)
})

// Call service method
const items = await FetchData("search term")
displayItems(items)
```

### Pembaruan Progres

Ideal untuk operasi yang berjalan lama, seperti pengunggahan berkas, pemrosesan batch, impor data berukuran besar, atau pengodean video. Kirim event progres selama operasi berlangsung untuk memperbarui bilah progres, teks status, atau indikator langkah pada UI, sehingga pengguna memperoleh umpan balik secara waktu nyata.

**Go:**

```go
func (s *Service) ProcessFiles(files []string) error {
    total := len(files)

    for i, file := range files {
        // Process file
        processFile(file)

        // Emit progress event
        s.app.Event.Emit("progress", map[string]interface{}{
            "current": i + 1,
            "total":   total,
            "percent": float64(i+1) / float64(total) * 100,
            "file":    file,
        })
    }

    s.app.Event.Emit("processing-complete")
    return nil
}
```

**JavaScript:**

```javascript
import { Events } from '@wailsio/runtime'

// Update progress bar
Events.On('progress', (data) => {
    progressBar.style.width = `${data.percent}%`
    statusText.textContent = `Processing ${data.file}... (${data.current}/${data.total})`
})

// Handle completion
Events.Once('processing-complete', () => {
    progressBar.style.width = '100%'
    statusText.textContent = 'Complete!'
    setTimeout(() => hideProgressBar(), 2000)
})
```

### Komunikasi Antarjendela

Sangat cocok untuk aplikasi dengan beberapa jendela, seperti panel pengaturan, dasbor, atau penampil dokumen. Siarkan event untuk menyinkronkan status di seluruh jendela (perubahan tema, preferensi pengguna), atau kirim event tertarget ke jendela tertentu untuk pembaruan khusus jendela tersebut.

**Go:**

```go
// Broadcast to all windows
app.Event.Emit("theme-changed", "dark")

// Send to specific window
preferencesWindow.EmitEvent("settings-updated", settings)

// Per-window listener — receive on the global event bus, but
// gate by the event's source-window name (set automatically when
// a window emits via window.EmitEvent).
app.Event.On("request-data", func(e *application.CustomEvent) {
    if e.Sender != window1.Name() {
        return
    }
    window1.EmitEvent("data-response", data)
})
```

**JavaScript:**

```javascript
import { Events } from '@wailsio/runtime'

// Listen in any window
Events.On('theme-changed', (theme) => {
    document.body.className = theme
})
```

### Sinkronisasi Status

Gunakan pola ini saat Anda perlu menjaga agar status frontend dan backend tetap sinkron, misalnya untuk sesi pengguna, konfigurasi aplikasi, atau fitur kolaboratif. Saat status berubah di backend, kirim event untuk memperbarui semua frontend yang terhubung, sehingga konsistensi di seluruh aplikasi tetap terjaga.

**Go:**

```go
type StateService struct {
    app   *application.App
    state map[string]interface{}
    mu    sync.RWMutex
}

func (s *StateService) UpdateState(key string, value interface{}) {
    s.mu.Lock()
    s.state[key] = value
    s.mu.Unlock()

    // Notify all windows
    s.app.Event.Emit("state-updated", map[string]interface{}{
        "key":   key,
        "value": value,
    })
}

func (s *StateService) GetState(key string) interface{} {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.state[key]
}
```

**JavaScript:**

```javascript
import { Events } from '@wailsio/runtime'
import { GetState } from './bindings/StateService'

// Keep local state in sync
let localState = {}

Events.On('state-updated', async (data) => {
    localState[data.key] = data.value
    updateUI(data.key, data.value)
})

// Initialize state
const initialState = await GetState("all")
localState = initialState
```

### Notifikasi Berbasis Event

Paling sesuai untuk menampilkan umpan balik kepada pengguna, seperti konfirmasi keberhasilan, peringatan kesalahan, atau pesan informasi. Alih-alih memanggil kode UI secara langsung dari layanan, kirim event notifikasi yang ditangani frontend secara konsisten. Dengan demikian, Anda dapat dengan mudah mengubah gaya notifikasi atau menambahkan fitur seperti riwayat notifikasi.

**Go:**

```go
type NotificationService struct {
    app *application.App
}

func (s *NotificationService) Success(message string) {
    s.app.Event.Emit("notification", map[string]interface{}{
        "type":    "success",
        "message": message,
    })
}

func (s *NotificationService) Error(message string) {
    s.app.Event.Emit("notification", map[string]interface{}{
        "type":    "error",
        "message": message,
    })
}

func (s *NotificationService) Info(message string) {
    s.app.Event.Emit("notification", map[string]interface{}{
        "type":    "info",
        "message": message,
    })
}
```

**JavaScript:**

```javascript
import { Events } from '@wailsio/runtime'

// Unified notification handler
Events.On('notification', (data) => {
    const toast = document.createElement('div')
    toast.className = `toast toast-${data.type}`
    toast.textContent = data.message

    document.body.appendChild(toast)

    setTimeout(() => {
        toast.classList.add('fade-out')
        setTimeout(() => toast.remove(), 300)
    }, 3000)
})
```

## Contoh Lengkap

**Go:**

```go
package main

import (
    "sync"
    "time"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)

type EventDemoService struct {
    app *application.App
    mu  sync.Mutex
}

func NewEventDemoService(app *application.App) *EventDemoService {
    service := &EventDemoService{app: app}

    // Listen for custom events
    app.Event.On("user-action", func(e *application.CustomEvent) {
        data := e.Data.(map[string]interface{})
        app.Logger.Info("User action received", "data", data)
    })

    return service
}

func (s *EventDemoService) StartLongTask() {
    go func() {
        s.app.Event.Emit("task-started")

        for i := 1; i <= 10; i++ {
            time.Sleep(500 * time.Millisecond)

            s.app.Event.Emit("task-progress", map[string]interface{}{
                "step":    i,
                "total":   10,
                "percent": i * 10,
            })
        }

        s.app.Event.Emit("task-completed", map[string]interface{}{
            "message": "Task finished successfully!",
        })
    }()
}

func (s *EventDemoService) BroadcastMessage(message string) {
    s.app.Event.Emit("broadcast", message)
}

func main() {
    app := application.New(application.Options{
        Name: "Event Demo",
    })

    // Handle application lifecycle
    app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(e *application.ApplicationEvent) {
        app.Logger.Info("Application started!")
    })

    app.OnShutdown(func() {
        app.Logger.Info("Application shutting down...")
    })

    // Register service (RegisterService returns nothing).
    service := NewEventDemoService(app)
    app.RegisterService(application.NewService(service))

    // Create window
    window := app.Window.New()

    // Handle window events
    window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
        window.EmitEvent("window-state", "focused")
    })

    window.OnWindowEvent(events.Common.WindowLostFocus, func(e *application.WindowEvent) {
        window.EmitEvent("window-state", "blurred")
    })

    window.Show()
    app.Run()
}
```

**JavaScript:**

```javascript
import { Events } from '@wailsio/runtime'
import { StartLongTask, BroadcastMessage } from './bindings/EventDemoService'

// Task events
Events.On('task-started', () => {
    console.log('Task started...')
    document.getElementById('status').textContent = 'Running...'
})

Events.On('task-progress', (data) => {
    const progressBar = document.getElementById('progress')
    progressBar.style.width = `${data.percent}%`
    console.log(`Step ${data.step} of ${data.total}`)
})

Events.Once('task-completed', (data) => {
    console.log('Task completed!', data.message)
    document.getElementById('status').textContent = data.message
})

// Broadcast events
Events.On('broadcast', (message) => {
    console.log('Broadcast:', message)
    alert(message)
})

// Window state events
Events.On('window-state', (state) => {
    console.log('Window is now:', state)
    document.body.dataset.windowState = state
})

// Trigger long task
document.getElementById('startTask').addEventListener('click', async () => {
    await StartLongTask()
})

// Send broadcast
document.getElementById('broadcast').addEventListener('click', async () => {
    const message = document.getElementById('message').value
    await BroadcastMessage(message)
})
```

## Event Bawaan

Wails menyediakan event sistem bawaan untuk siklus hidup aplikasi dan jendela. Event ini dikirim secara otomatis oleh framework.

### Event Umum vs Event Native Platform

Wails menyediakan dua jenis event sistem:

**Event Umum** (`events.Common.*`) adalah abstraksi lintas platform yang bekerja secara konsisten di macOS, Windows, dan Linux. Gunakan event ini dalam aplikasi Anda untuk mendapatkan portabilitas maksimal.

**Event Native Platform** (`events.Mac.*`, `events.Windows.*`, `events.Linux.*`) adalah event khusus sistem operasi yang mendasari pemetaan Event Umum. Event ini menyediakan akses ke perilaku dan kasus khusus platform.

**Cara Kerjanya:**

```go
import "github.com/wailsapp/wails/v3/pkg/events"

// ✅ RECOMMENDED: Use Common Events for cross-platform code
window.OnWindowEvent(events.Common.WindowClosing, func(e *application.WindowEvent) {
    // This works on all platforms
})

// Platform-specific events for advanced use cases
window.OnWindowEvent(events.Mac.WindowWillClose, func(e *application.WindowEvent) {
    // macOS-specific "will close" event (before WindowClosing)
})

window.OnWindowEvent(events.Windows.WindowClosing, func(e *application.WindowEvent) {
    // Windows-specific close event
})
```

**Pemetaan Event:**

Event native platform dipetakan secara otomatis ke Event Umum:

- macOS: `events.Mac.WindowShouldClose` → `events.Common.WindowClosing`
- Windows: `events.Windows.WindowClosing` → `events.Common.WindowClosing`
- Linux: `events.Linux.WindowDeleteEvent` → `events.Common.WindowClosing`

Pemetaan ini berlangsung secara otomatis di latar belakang. Karena itu, saat Anda memantau `events.Common.WindowClosing`, Anda akan menerimanya di platform apa pun.

**Kapan Menggunakan Masing-Masing Jenis:**

- **Gunakan Event Umum** untuk 99% kode aplikasi Anda karena event ini memberikan perilaku yang konsisten di berbagai platform
- **Gunakan Event Native Platform** hanya saat Anda memerlukan fungsionalitas khusus platform yang tidak tersedia dalam Event Umum (misalnya, event siklus hidup jendela khusus macOS atau event manajemen daya Windows)

### Event Aplikasi

| Event | Deskripsi | Waktu Dikirim | Dapat Dibatalkan |
| --- | --- | --- | --- |
| `ApplicationOpenedWithFile` | Aplikasi dibuka dengan sebuah berkas | Saat aplikasi diluncurkan dengan sebuah file (misalnya, melalui asosiasi file) | Tidak |
| `ApplicationStarted` | Aplikasi telah selesai diluncurkan | Setelah inisialisasi aplikasi selesai dan aplikasi siap digunakan | Tidak |
| `ApplicationLaunchedWithUrl` | Aplikasi diluncurkan dengan URL | Saat aplikasi diluncurkan melalui skema URL | Tidak |
| `ThemeChanged` | Tema sistem berubah | Saat tema OS beralih antara mode terang dan gelap | Tidak |
| `SystemWillSleep` | Sistem akan segera ditangguhkan | Tepat sebelum OS ditangguhkan (macOS / Windows / Linux dengan logind) | Tidak |
| `SystemDidWake` | Sistem dilanjutkan setelah ditangguhkan | Segera setelah sistem aktif kembali dari mode tidur (macOS / Windows / Linux dengan logind) | Tidak |

**Penggunaan:**

```go
import "github.com/wailsapp/wails/v3/pkg/events"

app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(e *application.ApplicationEvent) {
    app.Logger.Info("Application ready!")
})

app.Event.OnApplicationEvent(events.Common.ThemeChanged, func(e *application.ApplicationEvent) {
    // Update app theme
})
```

### Peristiwa Jendela

| Peristiwa | Deskripsi | Waktu Dipancarkan | Dapat Dibatalkan |
| --- | --- | --- | --- |
| `WindowClosing` | Jendela akan segera ditutup | Sebelum jendela ditutup (pengguna mengeklik X atau Close() dipanggil) | Ya |
| `WindowDidMove` | Jendela dipindahkan ke posisi baru | Setelah posisi jendela berubah (dengan debounce) | Tidak |
| `WindowDidResize` | Ukuran jendela diubah | Setelah ukuran jendela berubah | Tidak |
| `WindowDPIChanged` | Penskalaan DPI jendela berubah | Saat jendela dipindahkan di antara monitor dengan DPI berbeda (Windows) | Tidak |
| `WindowFilesDropped` | File dijatuhkan melalui seret dan jatuhkan native OS | Setelah file dari OS dijatuhkan ke jendela | Tidak |
| `WindowFocus` | Jendela memperoleh fokus | Saat jendela menjadi aktif | Tidak |
| `WindowFullscreen` | Jendela memasuki mode layar penuh | Setelah Fullscreen() atau pengguna memasuki mode layar penuh | Tidak |
| `WindowHide` | Jendela disembunyikan | Setelah Hide() atau jendela menjadi terhalang | Tidak |
| `WindowLostFocus` | Jendela kehilangan fokus | Saat jendela menjadi tidak aktif | Tidak |
| `WindowMaximise` | Jendela dimaksimalkan | Setelah Maximise() atau pengguna memaksimalkan jendela | Ya (macOS) |
| `WindowMinimise` | Jendela diminimalkan | Setelah Minimise() atau pengguna meminimalkan jendela | Ya (macOS) |
| `WindowRestore` | Jendela dipulihkan dari keadaan diminimalkan/dimaksimalkan | Setelah Restore() (terutama di Windows) | Tidak |
| `WindowRuntimeReady` | Runtime Wails telah dimuat dan siap | Saat inisialisasi runtime JavaScript selesai | Tidak |
| `WindowShow` | Jendela menjadi terlihat | Setelah Show() dipanggil atau jendela menjadi terlihat | Tidak |
| `WindowUnFullscreen` | Jendela keluar dari mode layar penuh | Setelah UnFullscreen() dipanggil atau pengguna keluar dari mode layar penuh | Tidak |
| `WindowUnMaximise` | Jendela keluar dari status dimaksimalkan | Setelah UnMaximise() dipanggil atau pengguna membatalkan pemaksimalan jendela | Ya (macOS) |
| `WindowUnMinimise` | Jendela keluar dari status diminimalkan | Setelah UnMinimise()/Restore() dipanggil atau pengguna memulihkan jendela | Ya (macOS) |
| `WindowZoomIn` | Perbesaran konten jendela meningkat | Setelah ZoomIn() dipanggil (terutama di macOS) | Ya (macOS) |
| `WindowZoomOut` | Perbesaran konten jendela menurun | Setelah ZoomOut() dipanggil (terutama di macOS) | Ya (macOS) |
| `WindowZoomReset` | Perbesaran konten jendela diatur ulang ke 100% | Setelah ZoomReset() dipanggil (terutama di macOS) | Ya (macOS) |
| `WindowDropZoneFilesDropped` | File dijatuhkan ke zona penjatuhan yang ditentukan dalam JS | Saat file dijatuhkan ke elemen yang memiliki zona penjatuhan | Tidak |

**Penggunaan:**

```go
import "github.com/wailsapp/wails/v3/pkg/events"

// Listen for window events
window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
    app.Logger.Info("Window focused")
})

// Cancel window close
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    dlg := app.Dialog.Question().SetMessage("Close window?")
    yes := dlg.AddButton("Yes")
    no := dlg.AddButton("No")
    dlg.SetDefaultButton(yes)
    dlg.SetCancelButton(no)
    no.OnClick(func() { e.Cancel() })

    dlg.Show()
})

// Wait for runtime ready
window.OnWindowEvent(events.Common.WindowRuntimeReady, func(e *application.WindowEvent) {
    app.Logger.Info("Runtime ready, safe to emit events to frontend")
    window.EmitEvent("app-initialized", data)
})
```

**Catatan Penting:**

- **WindowRuntimeReady** sangat penting—tunggu peristiwa ini sebelum memancarkan peristiwa ke frontend
- **WindowDidMove** dan **WindowDidResize** menerapkan debounce (nilai bawaan 50 ms) untuk mencegah banjir peristiwa
- **Peristiwa yang dapat dibatalkan** dapat dicegah dengan memanggil `event.Cancel()` dalam handler `RegisterHook()`
- **WindowFilesDropped** digunakan untuk penjatuhan file native OS; **WindowDropZoneFilesDropped** digunakan untuk zona penjatuhan berbasis web
- Beberapa peristiwa bersifat khusus platform (misalnya, WindowDPIChanged di Windows dan peristiwa zoom terutama di macOS)

## Konvensi Penamaan Peristiwa

```go
// Good - descriptive and specific
app.Event.Emit("user:logged-in", user)
app.Event.Emit("data:fetch:complete", results)
app.Event.Emit("ui:theme:changed", theme)

// Bad - vague and unclear
app.Event.Emit("event1", data)
app.Event.Emit("update", stuff)
app.Event.Emit("e", value)
```

## Pertimbangan Performa

### Debounce untuk Peristiwa Berfrekuensi Tinggi

```go
type Service struct {
    app            *application.App
    lastEmit       time.Time
    debounceWindow time.Duration
}

func (s *Service) EmitWithDebounce(event string, data interface{}) {
    now := time.Now()
    if now.Sub(s.lastEmit) < s.debounceWindow {
        return // Skip this emission
    }

    s.app.Event.Emit(event, data)
    s.lastEmit = now
}
```

### Throttling Peristiwa

```javascript
import { Events } from '@wailsio/runtime'

let lastUpdate = 0
const throttleMs = 100

Events.On('high-frequency-event', (data) => {
    const now = Date.now()
    if (now - lastUpdate < throttleMs) {
        return // Skip this update
    }

    processUpdate(data)
    lastUpdate = now
})
```
