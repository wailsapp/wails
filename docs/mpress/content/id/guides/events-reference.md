---
title: "Panduan Peristiwa"
description: "Panduan praktis menggunakan peristiwa di Wails v3 untuk komunikasi aplikasi dan pengelolaan siklus hidup"
slug: "guides/events-reference"
sourcePath: "guides/events-reference.md"
---

**CATATAN: Panduan ini masih dalam proses penyusunan**

## Panduan Peristiwa

Peristiwa merupakan inti komunikasi dalam aplikasi Wails. Peristiwa memungkinkan berbagai bagian aplikasi berkomunikasi satu sama lain tanpa saling terikat erat. Panduan ini akan menjelaskan semua yang perlu Anda ketahui agar dapat menggunakan peristiwa secara efektif dalam aplikasi Wails.

## Memahami Peristiwa Wails

Anggap peristiwa sebagai pesan yang disiarkan ke seluruh aplikasi. Setiap bagian aplikasi dapat mendengarkan pesan ini dan meresponsnya. Hal ini sangat berguna untuk:

- **Merespons perubahan jendela**: Mengetahui saat jendela diminimalkan, dimaksimalkan, atau dipindahkan
- **Menangani peristiwa sistem**: Merespons perubahan tema atau peristiwa daya
- **Logika aplikasi khusus**: Membuat peristiwa sendiri untuk fitur seperti pembaruan data atau tindakan pengguna
- **Komunikasi antar-komponen**: Memungkinkan berbagai bagian aplikasi berkomunikasi tanpa dependensi langsung

## Konvensi Penamaan Peristiwa

Semua peristiwa Wails mengikuti pola namespace untuk menunjukkan asalnya dengan jelas:

- `common:` - Peristiwa lintas platform yang berfungsi di Windows, macOS, dan Linux
- `windows:` - Peristiwa khusus Windows
- `mac:` - Peristiwa khusus macOS\
- `linux:` - Peristiwa khusus Linux

Contoh:

- `common:WindowFocus` - Jendela mendapatkan fokus (berfungsi di semua platform)
- `windows:APMSuspend` - Sistem akan ditangguhkan (khusus Windows)
- `mac:ApplicationDidBecomeActive` - Aplikasi menjadi aktif (khusus macOS)

## Memulai dengan Peristiwa

### Mendengarkan Peristiwa (Frontend)

Kasus penggunaan yang paling umum adalah mendengarkan peristiwa dalam kode frontend:

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

### Memancarkan Peristiwa (Backend)

Dari kode Go, Anda dapat memancarkan peristiwa yang dapat didengarkan oleh frontend:

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

### Memancarkan Peristiwa (Frontend)

Meskipun tidak terlalu umum digunakan, Anda juga dapat memancarkan peristiwa dari frontend yang dapat didengarkan oleh kode Go:

```javascript
import { Events } from '@wailsio/runtime';

// Event without data
Events.Emit('myapp:close-window')

// Event with data
Events.Emit('myapp:disconnect-requested', 'id-123')
```

Jika menggunakan TypeScript di frontend dan [mendaftarkan peristiwa bertipe](#peristiwa-bertipe-dengan-keamanan-tipe) dalam kode Go, Anda akan memperoleh pelengkapan otomatis/pemeriksaan nama peristiwa serta pemeriksaan tipe data.

### Menghapus Listener Peristiwa

Selalu bersihkan listener peristiwa saat tidak lagi diperlukan:

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

### 1. Menjeda/Melanjutkan saat Fokus Jendela Berubah

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

Pastikan aplikasi selalu selaras dengan tema sistem:

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

Buat aplikasi agar menerima file yang diseret dan dilepas:

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

Tangani peristiwa khusus platform saat diperlukan:

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

## Membuat Peristiwa Khusus

Anda dapat membuat peristiwa sendiri untuk kebutuhan khusus aplikasi.

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

## Peristiwa Bertipe dengan Keamanan Tipe

Wails v3 mendukung peristiwa bertipe dengan keamanan tipe TypeScript penuh melalui pendaftaran peristiwa dan pembuatan binding otomatis.

### Mendaftarkan Peristiwa Khusus

Panggil `application.RegisterEvent` pada saat inisialisasi untuk mendaftarkan nama peristiwa khusus beserta tipe datanya:

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
`RegisterEvent` dimaksudkan untuk dipanggil pada saat inisialisasi dan akan mengalami panic jika:

- Argumen tidak valid
- Nama peristiwa yang sama didaftarkan dua kali dengan tipe data yang berbeda

@end

@note{type="info"}
Peristiwa yang sama aman untuk didaftarkan beberapa kali selama tipe datanya selalu sama. Hal ini dapat berguna untuk memastikan suatu peristiwa terdaftar ketika salah satu dari beberapa paket dimuat.

@end

### Manfaat Pendaftaran Peristiwa

Setelah didaftarkan, tipe argumen data yang diteruskan ke `Event.Emit` akan diperiksa terhadap tipe yang ditentukan. Jika tidak cocok:

- Kesalahan dipancarkan dan dicatat dalam log (atau diteruskan ke handler kesalahan yang terdaftar)
- Peristiwa yang bermasalah tidak akan diteruskan
- Hal ini memastikan bahwa bidang data pada peristiwa terdaftar selalu dapat ditetapkan ke tipe yang dideklarasikan

### Mode Ketat

Gunakan tag build `strictevents` untuk mengaktifkan peringatan bagi event yang tidak terdaftar selama pengembangan:

```bash
go build -tags strictevents
```

Saat mode ketat diaktifkan, runtime mengeluarkan paling banyak satu peringatan untuk setiap nama event yang tidak terdaftar agar log tidak dipenuhi pesan berulang.

### Pembuatan Binding TypeScript

Generator binding menghasilkan definisi TypeScript dan kode penghubung untuk mendukung event bertipe secara transparan di frontend.

#### 1. Siapkan Plugin Vite

Di `vite.config.ts` Anda:

```typescript
import { defineConfig } from 'vite'
import wails from '@wailsio/runtime/plugins/vite'

export default defineConfig({
  plugins: [wails()],
})
```

#### 2. Buat Binding

Jalankan generator binding:

```bash
wails3 generate bindings
```

Tindakan ini membuat file TypeScript di direktori frontend Anda yang berisi pembuat event bertipe dan antarmuka data.

#### 3. Gunakan Event Bertipe di Frontend

```typescript
import { Events } from '@wailsio/runtime'
import { UserLogin, MonitorStats } from './bindings/events'

// Type-safe event emission with autocomplete
Events.Emit(UserLogin({
    UserID: "123",
    Username: "john_doe",
    LoginTime: new Date().toISOString()
}))

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
- **Error saat kompilasi** untuk tipe data yang tidak cocok
- Dokumentasi **IntelliSense**

## Referensi Event

### Event Umum (Lintas Platform)

Event berikut berfungsi di semua platform:

| Event | Deskripsi | Waktu Penggunaan |
| --- | --- | --- |
| `common:ApplicationStarted` | Aplikasi telah dimulai sepenuhnya | Inisialisasi aplikasi Anda dan muat status yang tersimpan |
| `common:WindowRuntimeReady` | Runtime Wails siap | Mulai lakukan panggilan API Wails |
| `common:ThemeChanged` | Tema sistem berubah | Perbarui tampilan aplikasi |
| `common:SystemWillSleep` | Sistem akan segera ditangguhkan | Tulis status tertunda dan tutup soket |
| `common:SystemDidWake` | Sistem dilanjutkan setelah ditangguhkan | Hubungkan kembali dan segarkan data usang |
| `common:WindowFocus` | Jendela mendapat fokus | Lanjutkan aktivitas dan segarkan data |
| `common:WindowLostFocus` | Jendela kehilangan fokus | Jeda aktivitas dan simpan status |
| `common:WindowMinimise` | Jendela diminimalkan | Jeda rendering dan kurangi penggunaan sumber daya |
| `common:WindowMaximise` | Jendela dimaksimalkan | Sesuaikan tata letak untuk layar penuh |
| `common:WindowRestore` | Jendela dipulihkan dari keadaan diminimalkan atau dimaksimalkan | Kembalikan ke tata letak normal |
| `common:WindowClosing` | Jendela akan segera ditutup | Simpan data dan bersihkan sumber daya |
| `common:WindowFilesDropped` | File dijatuhkan ke jendela | Tangani impor file |
| `common:WindowDidResize` | Ukuran jendela diubah | Sesuaikan tata letak dan render ulang bagan |
| `common:WindowDidMove` | Jendela dipindahkan | Perbarui fitur yang bergantung pada posisi |

### Event Khusus Platform

#### Event Windows

Peristiwa utama untuk aplikasi Windows:

| Peristiwa | Deskripsi | Kasus Penggunaan |
| --- | --- | --- |
| `windows:SystemThemeChanged` | Tema Windows berubah | Perbarui warna aplikasi |
| `windows:APMSuspend` | Sistem memasuki mode tidur | Simpan status, jeda operasi |
| `windows:APMResumeAutomatic` | Sistem kembali aktif (selalu dipicu saat kembali aktif) | Pulihkan status, segarkan data |
| `windows:APMResumeSuspend` | Sistem kembali aktif melalui input pengguna (setelah `APMResumeAutomatic`) | Bedakan pengaktifan oleh pengguna |
| `windows:APMPowerStatusChange` | Status daya berubah | Sesuaikan pengaturan performa |

#### Peristiwa macOS

Peristiwa penting aplikasi macOS:

| Peristiwa | Deskripsi | Kasus Penggunaan |
| --- | --- | --- |
| `mac:ApplicationDidBecomeActive` | Aplikasi menjadi aktif | Lanjutkan operasi |
| `mac:ApplicationDidResignActive` | Aplikasi menjadi tidak aktif | Jeda operasi |
| `mac:ApplicationWillTerminate` | Aplikasi akan ditutup | Lakukan pembersihan akhir |
| `mac:ApplicationWillSleep` | Sistem akan segera memasuki mode tidur | Simpan status, tutup soket |
| `mac:ApplicationDidWake` | Sistem kembali aktif | Hubungkan kembali, segarkan |
| `mac:ApplicationScreensDidSleep` | Layar memasuki mode tidur | Jeda rendering (berbeda dari mode tidur sistem) |
| `mac:ApplicationScreensDidWake` | Layar kembali aktif | Lanjutkan rendering |
| `mac:WindowDidEnterFullScreen` | Memasuki layar penuh | Sesuaikan UI untuk layar penuh |
| `mac:WindowDidExitFullScreen` | Keluar dari layar penuh | Pulihkan UI normal |

#### Peristiwa Linux

Peristiwa inti jendela Linux:

| Peristiwa | Deskripsi | Kasus Penggunaan |
| --- | --- | --- |
| `linux:SystemThemeChanged` | Tema desktop berubah | Perbarui tema aplikasi |
| `linux:SystemWillSleep` | Sistem akan segera memasuki mode tidur (logind) | Simpan status |
| `linux:SystemDidWake` | Sistem kembali aktif (logind) | Hubungkan kembali, segarkan |
| `linux:WindowFocusIn` | Jendela memperoleh fokus | Lanjutkan aktivitas |
| `linux:WindowFocusOut` | Jendela kehilangan fokus | Jeda aktivitas |
| `linux:WindowLoadStarted` | WebView mulai memuat | Tampilkan indikator pemuatan |
| `linux:WindowLoadRedirected` | WebView dialihkan | Lacak pengalihan navigasi |
| `linux:WindowLoadCommitted` | WebView telah melakukan commit pemuatan | Konten sedang diterima |
| `linux:WindowLoadFinished` | WebView selesai memuat | Sembunyikan indikator pemuatan, injeksikan JS/CSS |

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

## Sumber Kebenaran

Daftar lengkap peristiwa yang tersedia dapat ditemukan dalam kode sumber Wails:

- Peristiwa frontend: [`v3/internal/runtime/desktop/@wailsio/runtime/src/event_types.ts`](https://github.com/wailsapp/wails/blob/main/v3/internal/runtime/desktop/@wailsio/runtime/src/event_types.ts)
- Peristiwa backend: [`v3/pkg/events/events.go`](https://github.com/wailsapp/wails/blob/main/v3/pkg/events/events.go)

Selalu rujuk file-file ini untuk mengetahui nama dan ketersediaan peristiwa terbaru.

## Ringkasan

Peristiwa di Wails menyediakan cara yang andal dan tidak saling bergantung untuk menangani komunikasi dalam aplikasi Anda. Dengan mengikuti pola dan praktik dalam panduan ini, Anda dapat membuat aplikasi responsif yang memperhitungkan platform serta bereaksi dengan lancar terhadap perubahan sistem dan interaksi pengguna.

Ingat: mulailah dengan peristiwa umum untuk kompatibilitas lintas platform, tambahkan peristiwa khusus platform bila diperlukan, dan selalu bersihkan listener peristiwa Anda untuk mencegah kebocoran memori.
