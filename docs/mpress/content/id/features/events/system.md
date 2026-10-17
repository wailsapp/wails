---
title: "Sistem Peristiwa"
description: "Berkomunikasi antar komponen dengan sistem peristiwa"
slug: "features/events/system"
sourcePath: "features/events/system.md"
---

## Sistem Peristiwa

Wails menyediakan **sistem peristiwa terpadu** untuk komunikasi pub/sub. Pancarkan peristiwa dari mana saja dan dengarkan dari mana saja—dari Go ke JavaScript, JavaScript ke Go, maupun satu jendela ke jendela lain—sehingga memungkinkan arsitektur yang terpisah dengan peristiwa bertipe dan hook siklus hidup.

## Mulai Cepat

**Go (memancarkan):**

```go
app.Event.Emit("user-logged-in", map[string]interface{}{
    "userId": 123,
    "name": "Alice",
})
```

**JavaScript (mendengarkan):**

```javascript
import { Events } from '@wailsio/runtime'

Events.On("user-logged-in", (event) => {
    console.log(`User ${event.data.name} logged in`)
})
```

**Selesai!** Pub/sub lintas bahasa.

## Jenis Peristiwa

### Peristiwa Kustom

Peristiwa khusus aplikasi Anda:

```go
// Emit from Go
app.Event.Emit("order-created", order)
app.Event.Emit("payment-processed", payment)
app.Event.Emit("notification", message)
```

```javascript
// Listen in JavaScript
Events.On("order-created", handleOrder)
Events.On("payment-processed", handlePayment)
Events.On("notification", showNotification)
```

### Peristiwa Sistem

Peristiwa bawaan OS dan aplikasi:

```go
import "github.com/wailsapp/wails/v3/pkg/events"

// Theme changes
app.Event.OnApplicationEvent(events.Common.ThemeChanged, func(e *application.ApplicationEvent) {
    if e.Context().IsDarkMode() {
        app.Logger.Info("Dark mode enabled")
    }
})

// Application lifecycle
app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(e *application.ApplicationEvent) {
    app.Logger.Info("Application started")
})
```

### Peristiwa Jendela

Peristiwa khusus jendela:

```go
window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
    app.Logger.Info("Window focused")
})

window.OnWindowEvent(events.Common.WindowClosing, func(e *application.WindowEvent) {
    app.Logger.Info("Window closing")
})
```

## Memancarkan Peristiwa

### Dari Go

**Pemancaran dasar:**

```go
app.Event.Emit("event-name", data)
```

**Dengan berbagai jenis data:**

```go
// String
app.Event.Emit("message", "Hello")

// Number
app.Event.Emit("count", 42)

// Struct
app.Event.Emit("user", User{ID: 1, Name: "Alice"})

// Map
app.Event.Emit("config", map[string]interface{}{
    "theme": "dark",
    "fontSize": 14,
})

// Array
app.Event.Emit("items", []string{"a", "b", "c"})
```

**Ke jendela tertentu:**

```go
window.EmitEvent("window-specific-event", data)
```

### Dari JavaScript

```javascript
import { Events } from '@wailsio/runtime'

// Emit to Go
Events.Emit("button-clicked", { buttonId: "submit" })

// Emit to all windows
Events.Emit("broadcast-message", "Hello everyone")
```

## Mendengarkan Peristiwa

### Di Go

**Peristiwa aplikasi:**

```go
app.Event.On("custom-event", func(e *application.CustomEvent) {
    data := e.Data
    // Handle event
})
```

**Dengan penegasan tipe:**

```go
app.Event.On("user-updated", func(e *application.CustomEvent) {
    user := e.Data.(User)
    app.Logger.Info("User updated", "name", user.Name)
})
```

**Beberapa handler:**

```go
// All handlers will be called
app.Event.On("order-created", logOrder)
app.Event.On("order-created", sendEmail)
app.Event.On("order-created", updateInventory)
```

### Di JavaScript

**Listener dasar:**

```javascript
import { Events } from '@wailsio/runtime'

Events.On("event-name", (event) => {
    console.log("Event received:", event.data)
})
```

**Dengan pembersihan:**

```javascript
const unsubscribe = Events.On("event-name", handleEvent)

// Later, stop listening
unsubscribe()
```

**Beberapa handler:**

```javascript
Events.On("data-updated", updateUI)
Events.On("data-updated", saveToCache)
Events.On("data-updated", logChange)
```

**Handler sekali pakai:**

```javascript
Events.Once("data-updated", updateVariable)
```

## Peristiwa Sistem

### Peristiwa Aplikasi

**Peristiwa umum (lintas platform):**

```go
import "github.com/wailsapp/wails/v3/pkg/events"

// Application started
app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(e *application.ApplicationEvent) {
    app.Logger.Info("App started")
})

// Theme changed
app.Event.OnApplicationEvent(events.Common.ThemeChanged, func(e *application.ApplicationEvent) {
    isDark := e.Context().IsDarkMode()
    app.Event.Emit("theme-changed", isDark)
})

// System about to suspend
app.Event.OnApplicationEvent(events.Common.SystemWillSleep, func(e *application.ApplicationEvent) {
    flushPendingWrites()
})

// System resumed from suspend
app.Event.OnApplicationEvent(events.Common.SystemDidWake, func(e *application.ApplicationEvent) {
    reconnectSockets()
    refreshState()
})

// File opened
app.Event.OnApplicationEvent(events.Common.ApplicationOpenedWithFile, func(e *application.ApplicationEvent) {
    // Single-file: e.Context().Filename() returns ""; for multi-file launches
    // (Finder "Open With..." → multiple selection) use OpenedFiles().
    filePath := e.Context().Filename()
    openFile(filePath)
})
```

**Peristiwa khusus platform:**

@tabs{sync-key="platform"}
[macOS]
```go
// Application became active
app.Event.OnApplicationEvent(events.Mac.ApplicationDidBecomeActive, func(e *application.ApplicationEvent) {
    app.Logger.Info("App became active")
})

// Application will terminate
app.Event.OnApplicationEvent(events.Mac.ApplicationWillTerminate, func(e *application.ApplicationEvent) {
    cleanup()
})

// System about to sleep / resumed
app.Event.OnApplicationEvent(events.Mac.ApplicationWillSleep, func(e *application.ApplicationEvent) {
    flushPendingWrites()
})
app.Event.OnApplicationEvent(events.Mac.ApplicationDidWake, func(e *application.ApplicationEvent) {
    refreshState()
})

// Displays sleeping/waking — distinct from system sleep (e.g. lid lowered)
app.Event.OnApplicationEvent(events.Mac.ApplicationScreensDidSleep, func(e *application.ApplicationEvent) {
    pauseRendering()
})
app.Event.OnApplicationEvent(events.Mac.ApplicationScreensDidWake, func(e *application.ApplicationEvent) {
    resumeRendering()
})
```

[Windows]
```go
// Power status changed
app.Event.OnApplicationEvent(events.Windows.APMPowerStatusChange, func(e *application.ApplicationEvent) {
    app.Logger.Info("Power status changed")
})

// System suspending
app.Event.OnApplicationEvent(events.Windows.APMSuspend, func(e *application.ApplicationEvent) {
    saveState()
})

// System resumed. APMResumeAutomatic always fires on resume;
// APMResumeSuspend fires after it when the wake was triggered by
// user input. Prefer Common.SystemDidWake if you don't need to
// distinguish.
app.Event.OnApplicationEvent(events.Windows.APMResumeAutomatic, func(e *application.ApplicationEvent) {
    refreshState()
})
```

[Linux]
```go
// Application startup
app.Event.OnApplicationEvent(events.Linux.ApplicationStartup, func(e *application.ApplicationEvent) {
    app.Logger.Info("App starting")
})

// Theme changed
app.Event.OnApplicationEvent(events.Linux.SystemThemeChanged, func(e *application.ApplicationEvent) {
    updateTheme()
})

// System sleep/wake — sourced from systemd-logind's PrepareForSleep
// signal on the system bus. Requires logind/elogind exposing
// org.freedesktop.login1, so it may be unavailable on distros/setups
// without it (Alpine, Void, some Devuan setups).
app.Event.OnApplicationEvent(events.Linux.SystemWillSleep, func(e *application.ApplicationEvent) {
    flushPendingWrites()
})
app.Event.OnApplicationEvent(events.Linux.SystemDidWake, func(e *application.ApplicationEvent) {
    refreshState()
})
```

@end

### Peristiwa Jendela

**Peristiwa jendela yang umum:**

```go
// Window focus
window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
    app.Logger.Info("Window focused")
})

// Window blur
window.OnWindowEvent(events.Common.WindowLostFocus, func(e *application.WindowEvent) {
    app.Logger.Info("Window blurred")
})

// Window closing
window.OnWindowEvent(events.Common.WindowClosing, func(e *application.WindowEvent) {
    if hasUnsavedChanges() {
        e.Cancel()  // Prevent close
    }
})

// Window closed
window.OnWindowEvent(events.Common.WindowClosing, func(e *application.WindowEvent) {
    cleanup()
})
```

## Hook Peristiwa

Hook dijalankan **sebelum** listener standar dan dapat **membatalkan** peristiwa:

```go
// Hook - runs first, can cancel
window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
    if hasUnsavedChanges() {
        result := showConfirmdialog("Unsaved changes. Close anyway?")
        if result != "yes" {
            e.Cancel()  // Prevent window close
        }
    }
})

// Standard listener - runs after hooks
window.OnWindowEvent(events.Common.WindowClosing, func(e *application.WindowEvent) {
    app.Logger.Info("Window closing")
})
```

**Perbedaan utama:**

| Fitur | Hook | Listener Standar |
| --- | --- | --- |
| Urutan eksekusi | Pertama, sesuai urutan pendaftaran | Setelah hook, tanpa jaminan urutan |
| Pemblokiran | Sinkron, memblokir hook berikutnya | Asinkron, tanpa pemblokiran |
| Dapat membatalkan | Ya | Tidak (sudah disebarkan) |
| Kasus penggunaan | Kontrol alur, validasi | Pencatatan log, efek samping |

## Pola Peristiwa

### Pola Pub/Sub

```go
// Publisher (service)
type OrderService struct {
    app *application.App
}

func (o *OrderService) CreateOrder(items []Item) (*Order, error) {
    order := &Order{Items: items}
    
    if err := o.saveOrder(order); err != nil {
        return nil, err
    }
    
    // Publish event
    o.app.Event.Emit("order-created", order)
    
    return order, nil
}

// Subscribers
app.Event.On("order-created", func(e *application.CustomEvent) {
    order := e.Data.(*Order)
    sendConfirmationEmail(order)
})

app.Event.On("order-created", func(e *application.CustomEvent) {
    order := e.Data.(*Order)
    updateInventory(order)
})

app.Event.On("order-created", func(e *application.CustomEvent) {
    order := e.Data.(*Order)
    logOrder(order)
})
```

### Pola Permintaan/Respons

```go
// Frontend requests data
Emit("get-user-data", { userId: 123 })

// Backend responds
app.Event.On("get-user-data", func(e *application.CustomEvent) {
    data := e.Data.(map[string]interface{})
    userId := int(data["userId"].(float64))
    
    user := getUserFromDB(userId)
    
    // Send response
    app.Event.Emit("user-data-response", user)
})
```

```javascript
// Frontend receives response
Events.On("user-data-response", (event) => {
    const user = event.data
    displayUser(user)
})
```

**Catatan:** Untuk permintaan/respons, **binding lebih sesuai**. Gunakan peristiwa untuk notifikasi.

### Pola Siaran

```go
// Broadcast to all windows
app.Event.Emit("global-notification", "System update available")

// Each window handles it
Events.On("global-notification", (event) => {
    const message = event.data
    showNotification(message)
})
```

### Agregasi Peristiwa

```go
type EventAggregator struct {
    events []Event
    mu     sync.Mutex
}

func (ea *EventAggregator) Add(event Event) {
    ea.mu.Lock()
    defer ea.mu.Unlock()
    
    ea.events = append(ea.events, event)
    
    // Emit batch every 100 events
    if len(ea.events) >= 100 {
        app.Event.Emit("event-batch", ea.events)
        ea.events = nil
    }
}
```

## Contoh Lengkap

**Go:**

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/events"
)

type NotificationService struct {
    app *application.App
}

func (n *NotificationService) Notify(message string) {
    // Emit to all windows
    n.app.Event.Emit("notification", map[string]interface{}{
        "message":   message,
        "timestamp": time.Now(),
    })
}

func main() {
    app := application.New(application.Options{
        Name: "Event Demo",
    })
    
    notifService := &NotificationService{app: app}
    
    // System events
    app.Event.OnApplicationEvent(events.Common.ThemeChanged, func(e *application.ApplicationEvent) {
        isDark := e.Context().IsDarkMode()
        app.Event.Emit("theme-changed", isDark)
    })
    
    // Custom events from frontend
    app.Event.On("user-action", func(e *application.CustomEvent) {
        data := e.Data.(map[string]interface{})
        action := data["action"].(string)
        
        app.Logger.Info("User action", "action", action)
        
        // Respond
        notifService.Notify("Action completed: " + action)
    })
    
    // Window events
    window := app.Window.New()
    
    window.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
        app.Event.Emit("window-focused", window.Name())
    })
    
    window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
        // Confirm before close
        app.Event.Emit("confirm-close", nil)
        e.Cancel()  // Wait for confirmation
    })
    
    app.Run()
}
```

**JavaScript:**

```javascript
import { Events } from '@wailsio/runtime'

// Listen for notifications
Events.On("notification", (event) => {
    showNotification(event.data.message)
})

// Listen for theme changes
Events.On("theme-changed", (event) => {
    const isDark = event.data
    document.body.classList.toggle('dark', isDark)
})

// Listen for window focus
Events.On("window-focused", (event) => {
    const windowName = event.data
    console.log(`Window ${windowName} focused`)
})

// Handle close confirmation
Events.On("confirm-close", (event) => {
    if (confirm("Close window?")) {
        Events.Emit("close-confirmed", true)
    }
})

// Emit user actions
document.getElementById('button').addEventListener('click', () => {
    Events.Emit("user-action", { action: "button-clicked" })
})
```

## Praktik Terbaik

### ✅ Lakukan

- **Gunakan peristiwa untuk notifikasi** - Komunikasi satu arah
- **Gunakan binding untuk permintaan** - Komunikasi dua arah
- **Pertahankan konsistensi nama peristiwa** - Gunakan kebab-case
- **Dokumentasikan data peristiwa** - Kolom apa saja yang disertakan?
- **Hentikan langganan setelah selesai** - Cegah kebocoran memori
- **Gunakan hook untuk validasi** - Kendalikan alur peristiwa

### ❌ Jangan Lakukan

- **Jangan gunakan peristiwa untuk RPC** - Gunakan binding sebagai gantinya
- **Jangan memancarkan event terlalu sering** - Kelompokkan jika diperlukan
- **Jangan memblokir di dalam handler** - Pastikan handler tetap cepat
- **Jangan lupa berhenti berlangganan** - Dapat menyebabkan kebocoran memori
- **Jangan gunakan event untuk data berukuran besar** - Gunakan binding
- **Jangan membuat perulangan event** - A memancarkan B, B memancarkan A

## Langkah Berikutnya

@cards{cols="2"}
★ Panduan Event
Pelajari pola event dan pembuatan event yang aman terhadap tipe.

[Pelajari Selengkapnya →](/guides/events-reference/)

---
◆ API Event
Pelajari API event aplikasi dan jendela secara lengkap.

[Pelajari Selengkapnya →](/reference/events/)

---
🚀 Binding
Gunakan binding untuk permintaan/respons.

[Pelajari Selengkapnya →](/features/bindings/methods/)

---
▣ Event Jendela
Tangani event siklus hidup jendela.

[Pelajari Selengkapnya →](/features/windows/events/)

@end

---

**Ada pertanyaan?** Tanyakan di [Discord](https://discord.gg/JDdSxwjhGf) atau lihat [contoh event](https://github.com/wailsapp/wails/tree/master/v3/examples/events).
