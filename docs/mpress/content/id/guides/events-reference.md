---
title: "Panduan Event"
description: "Panduan praktis menggunakan event di Wails v3 untuk komunikasi aplikasi dan pengelolaan siklus hidup"
slug: "guides/events-reference"
sourcePath: "guides/events-reference.md"
---

**CATATAN: Panduan ini masih dalam proses penyusunan**

## Panduan Event

Event merupakan inti komunikasi dalam aplikasi Wails. Event memungkinkan berbagai bagian aplikasi Anda berkomunikasi tanpa saling terikat erat. Panduan ini akan menjelaskan semua hal yang perlu Anda ketahui untuk menggunakan event secara efektif dalam aplikasi Wails.

## Memahami Event Wails

Anggap event sebagai pesan yang disiarkan ke seluruh aplikasi Anda. Setiap bagian aplikasi dapat mendengarkan pesan ini dan meresponsnya sesuai kebutuhan. Hal ini sangat berguna untuk:

- **Merespons perubahan jendela**: Mengetahui saat jendela diminimalkan, dimaksimalkan, atau dipindahkan
- **Menangani event sistem**: Merespons perubahan tema atau event daya
- **Logika aplikasi khusus**: Membuat event sendiri untuk fitur seperti pembaruan data atau tindakan pengguna
- **Komunikasi lintas komponen**: Memungkinkan berbagai bagian aplikasi berkomunikasi tanpa dependensi langsung

## Konvensi Penamaan Event

Semua event Wails mengikuti pola namespace agar asalnya ditunjukkan dengan jelas:

- `common:` - Event lintas platform yang berfungsi di Windows, macOS, dan Linux
- `windows:` - Event khusus Windows
- `mac:` - Event khusus macOS\
- `linux:` - Event khusus Linux

Contoh:

- `common:WindowFocus` - Jendela mendapatkan fokus (berfungsi di semua platform)
- `windows:APMSuspend` - Sistem akan ditangguhkan (khusus Windows)
- `mac:ApplicationDidBecomeActive` - Aplikasi menjadi aktif (khusus macOS)

## Memulai Penggunaan Event

### Mendengarkan Event (Frontend)

Kasus penggunaan yang paling umum adalah mendengarkan event dalam kode frontend Anda:

```javascript
import { Events } from '@wailsio/runtime';

// Listen for when the window gains focus
Events.On('common:WindowFocus', () => {
    console.log('Window is now focused!');
    // Maybe refresh some data or resume animations
});

// Listen for theme changes
Events.On('common:ThemeChanged', (event) => {
    console.log('Theme changed:', event.data);
    // Update your app's theme accordingly
});

// Listen for custom events from your Go backend
Events.On('my-app:data-updated', (event) => {
    console.log('Data updated:', event.data);
    // Update your UI with the new data
});
```

### Memancarkan Event (Backend)

Dari kode Go, Anda dapat memancarkan event yang dapat didengarkan oleh frontend:

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "time"
)

func (s *Service) UpdateData() {
    // Do some data processing...

    // Notify the frontend
    app := application.Get()
    app.Event.Emit("my-app:data-updated",
        map[string]interface{}{
            "timestamp": time.Now(),
            "count": 42,
        },
	)
}
```

### Memancarkan Event (Frontend)

Meskipun tidak terlalu umum digunakan, Anda juga dapat memancarkan event dari frontend yang dapat didengarkan oleh kode Go:

```javascript
import { Events } from '@wailsio/runtime';

// Event without data
Events.Emit('myapp:close-window')

// Event with data
Events.Emit('myapp:disconnect-requested', 'id-123')
```

Jika Anda menggunakan TypeScript di frontend dan [mendaftarkan event bertipe](#event-bertipe-dengan-keamanan-tipe) dalam kode Go, Anda akan mendapatkan pelengkapan otomatis/pemeriksaan nama event serta pemeriksaan tipe data.

### Menghapus Listener Event

Selalu bersihkan listener event saat tidak lagi diperlukan:

```javascript
import { Events } from '@wailsio/runtime';

// Store the handler reference
const focusHandler = () => {
    console.log('Window focused');
};

// Add the listener — capture the unsubscribe function it returns
const unsubscribe = Events.On('common:WindowFocus', focusHandler);

// Later, remove this specific listener via the returned unsubscribe
unsubscribe();

// Or remove ALL listeners for one (or more) event names — Events.Off takes only event-name strings
Events.Off('common:WindowFocus');
// Events.Off('common:WindowFocus', 'common:WindowLostFocus'); // variadic
// Events.OffAll(); // remove every listener for every event (no args)
```

## Kasus Penggunaan Umum

### 1. Menjeda/Melanjutkan Berdasarkan Fokus Jendela

Banyak aplikasi perlu menjeda aktivitas tertentu saat jendela kehilangan fokus:

```javascript
import { Events } from '@wailsio/runtime';

let animationRunning = true;

Events.On('common:WindowLostFocus', () => {
    animationRunning = false;
    pauseBackgroundTasks();
});

Events.On('common:WindowFocus', () => {
    animationRunning = true;
    resumeBackgroundTasks();
});
```

### 2. Merespons Perubahan Tema

Selaraskan aplikasi Anda dengan tema sistem:

```javascript
import { Events } from '@wailsio/runtime';

Events.On('common:ThemeChanged', (event) => {
    const isDarkMode = event.data.isDarkMode;

    if (isDarkMode) {
        document.body.classList.add('dark-theme');
        document.body.classList.remove('light-theme');
    } else {
        document.body.classList.add('light-theme');
        document.body.classList.remove('dark-theme');
    }
});
```

### 3. Menangani File yang Dilepas

Buat aplikasi Anda menerima file yang diseret:

```javascript
import { Events } from '@wailsio/runtime';

Events.On('common:WindowFilesDropped', (event) => {
    const files = event.data.files;

    files.forEach(file => {
        console.log('File dropped:', file);
        // Process the dropped files
        handleFileUpload(file);
    });
});
```

### 4. Mengelola Siklus Hidup Jendela

Respons perubahan status jendela:

```javascript
import { Events } from '@wailsio/runtime';

Events.On('common:WindowClosing', () => {
    // Save user data before closing
    saveApplicationState();

    // You could also prevent closing by returning false
    // from a registered window close handler
});

Events.On('common:WindowMaximise', () => {
    // Adjust UI for maximized view
    adjustLayoutForMaximized();
});

Events.On('common:WindowRestore', () => {
    // Return UI to normal state
    adjustLayoutForNormal();
});
```

### 5. Fitur Khusus Platform

Tangani event khusus platform saat diperlukan:

```javascript
import { Events } from '@wailsio/runtime';

// Windows-specific power management
Events.On('windows:APMSuspend', () => {
    console.log('System is going to sleep');
    saveState();
});

Events.On('windows:APMResumeSuspend', () => {
    console.log('System woke up');
    refreshData();
});

// macOS-specific app lifecycle
Events.On('mac:ApplicationWillTerminate', () => {
    console.log('App is about to quit');
    performCleanup();
});
```

## Membuat Event Khusus

Anda dapat membuat event sendiri untuk kebutuhan khusus aplikasi.

### Backend (Go)

```go
// Emit a custom event when data changes

func (s *Service) ProcessUserData(userData UserData) error {
    // Process the data...

    app := application.Get()
    // Notify all listeners
    app.Event.Emit("user:data-processed",
        map[string]interface{}{
            "userId": userData.ID,
            "status": "completed",
            "timestamp": time.Now(),
        },
    )
    return nil
}

// Emit periodic updates
func (s *Service) StartMonitoring() {
    app := application.Get()
    ticker := time.NewTicker(5 * time.Second)
    go func() {
        for range ticker.C {
            stats := s.collectStats()
            app.Event.Emit("monitor:stats-updated", stats)
        }
    }()
}
```

### Frontend (JavaScript)

```javascript
import { Events } from '@wailsio/runtime';

// Listen for your custom events
Events.On('user:data-processed', (event) => {
    const { userId, status, timestamp } = event.data;

    showNotification(`User ${userId} processing ${status}`);
    updateUIWithNewData();
});

Events.On('monitor:stats-updated', (event) => {
    updateDashboard(event.data);
});
```

## Event Bertipe dengan Keamanan Tipe

Wails v3 mendukung event bertipe dengan keamanan tipe TypeScript penuh melalui pendaftaran event dan pembuatan binding otomatis.

### Mendaftarkan Event Khusus

Panggil `application.RegisterEvent` saat inisialisasi untuk mendaftarkan nama event khusus beserta tipe datanya:

```go
package main

import "github.com/wailsapp/wails/v3/pkg/application"

type UserLoginData struct {
    UserID   string
    Username string
    LoginTime string
}

type MonitorStats struct {
    CPUUsage    float64
    MemoryUsage float64
}

func init() {
    // Register events with their data types
    application.RegisterEvent[UserLoginData]("user:login")
    application.RegisterEvent[MonitorStats]("monitor:stats")

    // Register events without data (void events)
    application.RegisterEvent[application.Void]("app:ready")
}
```

@note{type="caution"}
`RegisterEvent` dimaksudkan untuk dipanggil saat inisialisasi dan akan mengalami panic jika:

- Argumen tidak valid
- Nama event yang sama didaftarkan dua kali dengan tipe data berbeda

@end

@note{type="info"}
Event yang sama aman untuk didaftarkan beberapa kali selama tipe datanya selalu sama. Hal ini dapat berguna untuk memastikan event telah terdaftar saat salah satu dari beberapa package dimuat.

@end

### Manfaat Pendaftaran Event

Setelah didaftarkan, tipe argumen data yang diteruskan ke `Event.Emit` akan diperiksa terhadap tipe yang ditentukan. Jika tidak cocok:

- Error dipancarkan dan dicatat dalam log (atau diteruskan ke handler error yang terdaftar)
- Event yang bermasalah tidak akan diteruskan
- Hal ini memastikan field data pada event terdaftar selalu dapat ditetapkan ke tipe yang dideklarasikan

### Mode Ketat

Gunakan build tag `strictevents` untuk mengaktifkan peringatan bagi event yang tidak terdaftar selama pengembangan:

```bash
go build -tags strictevents
```

Saat mode ketat diaktifkan, runtime memancarkan maksimal satu peringatan untuk setiap nama event yang tidak terdaftar agar log tidak dibanjiri.

### Pembuatan Binding TypeScript

Generator binding menghasilkan definisi TypeScript dan kode penghubung untuk menyediakan dukungan event bertipe secara transparan di frontend.

#### 1. Menyiapkan Plugin Vite

Dalam `vite.config.ts` Anda:

```typescript
import { defineConfig } from 'vite'
import wails from '@wailsio/runtime/plugins/vite'

export default defineConfig({
  plugins: [wails()],
})
```

#### 2. Membuat Binding

Jalankan generator binding:

```bash
wails3 generate bindings
```

Tindakan ini membuat file TypeScript dalam direktori frontend Anda yang berisi definisi data event bertipe, serta modul `bindings/events` yang mengekspor satu konstanta untuk setiap nama event terdaftar. Konstanta tersebut merupakan konstanta waktu kompilasi yang dapat ditemukan oleh generator (misalnya, `user:login` menjadi `UserLogin`).

#### 3. Menggunakan Event Bertipe di Frontend

```typescript
import { Events } from '@wailsio/runtime'
import { UserLogin, MonitorStats } from './bindings/events'

// Type-safe event emission with autocomplete
Events.Emit(UserLogin, {
    UserID: "123",
    Username: "john_doe",
    LoginTime: new Date().toISOString()
})

// Type-safe event listening
Events.On(UserLogin, (event) => {
    // event.data is typed as UserLoginData
    console.log(`User ${event.data.Username} logged in`)
})

Events.On(MonitorStats, (event) => {
    // event.data is typed as MonitorStats
    updateDashboard({
        cpu: event.data.CPUUsage,
        memory: event.data.MemoryUsage
    })
})
```

Event bertipe menyediakan:

- **Pelengkapan otomatis** untuk nama event
- **Pemeriksaan tipe** untuk data event
- **Error waktu kompilasi** untuk tipe data yang tidak cocok
- Dokumentasi **IntelliSense**

## Referensi Event

### Event Umum (Lintas Platform)

Event berikut berfungsi di semua platform:

| Event | Deskripsi | Waktu Penggunaan |
| --- | --- | --- |
| `common:ApplicationStarted` | Aplikasi telah sepenuhnya dimulai | Inisialisasi aplikasi, muat status tersimpan |
| `common:WindowRuntimeReady` | Runtime Wails siap | Mulai lakukan panggilan API Wails |
| `common:ThemeChanged` | Tema sistem berubah | Perbarui tampilan aplikasi |
| `common:SystemWillSleep` | Sistem akan segera ditangguhkan | Simpan status, tutup soket |
| `common:SystemDidWake` | Sistem dilanjutkan setelah penangguhan | Hubungkan kembali, segarkan data usang |
| `common:WindowFocus` | Jendela mendapatkan fokus | Lanjutkan aktivitas, segarkan data |
| `common:WindowLostFocus` | Jendela kehilangan fokus | Jeda aktivitas, simpan status |
| `common:WindowMinimise` | Jendela diminimalkan | Jeda rendering, kurangi penggunaan sumber daya |
| `common:WindowMaximise` | Jendela dimaksimalkan | Sesuaikan tata letak untuk layar penuh |
| `common:WindowRestore` | Jendela dipulihkan dari keadaan diminimalkan/dimaksimalkan | Kembalikan ke tata letak normal |
| `common:WindowClosing` | Jendela akan segera ditutup | Simpan data, bersihkan sumber daya |
| `common:WindowFilesDropped` | File dijatuhkan ke jendela | Tangani impor file |
| `common:WindowDidResize` | Ukuran jendela diubah | Sesuaikan tata letak, render ulang diagram |
| `common:WindowDidMove` | Jendela dipindahkan | Perbarui fitur yang bergantung pada posisi |

### Peristiwa Khusus Platform

#### Peristiwa Windows

Peristiwa utama untuk aplikasi Windows:

| Peristiwa | Deskripsi | Kasus Penggunaan |
| --- | --- | --- |
| `windows:SystemThemeChanged` | Tema Windows berubah | Perbarui warna aplikasi |
| `windows:APMSuspend` | Sistem sedang ditangguhkan | Simpan status, jeda operasi |
| `windows:APMResumeAutomatic` | Sistem dilanjutkan (selalu dipicu saat sistem dilanjutkan) | Pulihkan status, segarkan data |
| `windows:APMResumeSuspend` | Sistem dilanjutkan melalui input pengguna (setelah `APMResumeAutomatic`) | Bedakan pengaktifan yang dipicu pengguna |
| `windows:APMPowerStatusChange` | Status daya berubah | Sesuaikan pengaturan performa |

#### Peristiwa macOS

Peristiwa penting untuk aplikasi macOS:

| Peristiwa | Deskripsi | Kasus Penggunaan |
| --- | --- | --- |
| `mac:ApplicationDidBecomeActive` | Aplikasi menjadi aktif | Lanjutkan operasi |
| `mac:ApplicationDidResignActive` | Aplikasi menjadi tidak aktif | Jeda operasi |
| `mac:ApplicationWillTerminate` | Aplikasi akan ditutup | Pembersihan akhir |
| `mac:ApplicationWillSleep` | Sistem akan segera ditangguhkan | Simpan status, tutup soket |
| `mac:ApplicationDidWake` | Sistem dilanjutkan | Hubungkan kembali, segarkan |
| `mac:ApplicationScreensDidSleep` | Layar memasuki mode tidur | Jeda rendering (berbeda dari mode tidur sistem) |
| `mac:ApplicationScreensDidWake` | Layar aktif kembali | Lanjutkan rendering |
| `mac:WindowDidEnterFullScreen` | Memasuki mode layar penuh | Sesuaikan UI untuk layar penuh |
| `mac:WindowDidExitFullScreen` | Keluar dari mode layar penuh | Pulihkan UI normal |

#### Peristiwa Linux

Peristiwa utama jendela Linux:

| Peristiwa | Deskripsi | Kasus Penggunaan |
| --- | --- | --- |
| `linux:SystemThemeChanged` | Tema desktop berubah | Perbarui tema aplikasi |
| `linux:SystemWillSleep` | Sistem akan segera ditangguhkan (logind) | Simpan status |
| `linux:SystemDidWake` | Sistem dilanjutkan (logind) | Hubungkan kembali, segarkan |
| `linux:WindowFocusIn` | Jendela mendapatkan fokus | Lanjutkan aktivitas |
| `linux:WindowFocusOut` | Jendela kehilangan fokus | Jeda aktivitas |
| `linux:WindowLoadStarted` | WebView mulai memuat | Tampilkan indikator pemuatan |
| `linux:WindowLoadRedirected` | WebView dialihkan | Lacak pengalihan navigasi |
| `linux:WindowLoadCommitted` | WebView mengonfirmasi pemuatan | Konten sedang diterima |
| `linux:WindowLoadFinished` | WebView selesai dimuat | Sembunyikan indikator pemuatan, injeksikan JS/CSS |

## Praktik Terbaik

### 1. Gunakan Namespace Peristiwa

Saat membuat peristiwa khusus, gunakan namespace untuk menghindari konflik:

```javascript
import { Events } from '@wailsio/runtime';

// Good - namespaced events
Events.Emit('myapp:user:login');
Events.Emit('myapp:data:updated');
Events.Emit('myapp:network:connected');

// Avoid - generic names that might conflict
Events.Emit('login');
Events.Emit('update');
```

### 2. Bersihkan Listener

Selalu hapus listener peristiwa saat komponen dilepas:

```javascript
import { Events } from '@wailsio/runtime';

// React example
useEffect(() => {
    const handler = (event) => {
        // Handle event
    };

    const off = Events.On('common:WindowDidResize', handler);

    // Cleanup — call the unsubscribe returned by Events.On
    return () => {
        off();
    };
}, []);
```

### 3. Tangani Perbedaan Platform

Periksa ketersediaan pada platform saat menggunakan peristiwa khusus platform:

```javascript
import { Events } from '@wailsio/runtime';

// Platform-specific events can be registered unconditionally;
// they will simply never fire on unsupported platforms.
Events.On('windows:APMSuspend', handleSuspend);
Events.On('mac:ApplicationWillTerminate', handleTerminate);
```

### 4. Jangan Gunakan Peristiwa Secara Berlebihan

Meskipun peristiwa sangat berguna, jangan gunakan untuk segala hal:

- ✅ Gunakan peristiwa untuk: Notifikasi sistem, perubahan siklus hidup, pembaruan siaran
- ❌ Hindari peristiwa untuk: Nilai kembalian fungsi langsung, pembaruan satu komponen, operasi sinkron

## Men-debug Peristiwa

Untuk men-debug masalah peristiwa:

```javascript
import { Events } from '@wailsio/runtime';

// Log all events (development only)
if (isDevelopment) {
    const originalOn = Events.On;
    Events.On = function(eventName, handler) {
        console.log(`[Event Registered] ${eventName}`);
        return originalOn.call(this, eventName, function(event) {
            console.log(`[Event Fired] ${eventName}`, event);
            return handler(event);
        });
    };
}
```

## Sumber Acuan

Daftar lengkap peristiwa yang tersedia dapat ditemukan dalam kode sumber Wails:

- Peristiwa frontend: [`v3/internal/runtime/desktop/@wailsio/runtime/src/event_types.ts`](https://github.com/wailsapp/wails/blob/main/v3/internal/runtime/desktop/@wailsio/runtime/src/event_types.ts)
- Peristiwa backend: [`v3/pkg/events/events.go`](https://github.com/wailsapp/wails/blob/main/v3/pkg/events/events.go)

Selalu rujuk file-file ini untuk mengetahui nama dan ketersediaan peristiwa terbaru.

## Ringkasan

Peristiwa di Wails menyediakan cara yang efektif dan tidak saling bergantung untuk menangani komunikasi dalam aplikasi Anda. Dengan mengikuti pola dan praktik dalam panduan ini, Anda dapat membangun aplikasi responsif dan sadar platform yang bereaksi dengan lancar terhadap perubahan sistem serta interaksi pengguna.

Ingat: mulailah dengan peristiwa umum untuk kompatibilitas lintas platform, tambahkan peristiwa khusus platform bila diperlukan, dan selalu bersihkan listener peristiwa Anda untuk mencegah kebocoran memori.
