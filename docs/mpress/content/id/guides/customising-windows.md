---
title: "Menyesuaikan Jendela di Wails"
description: "Sesuaikan tampilan dan perilaku jendela dalam aplikasi Wails Anda"
slug: "guides/customising-windows"
sourcePath: "guides/customising-windows.md"
---

Platform yang Relevan: <span class="mpress-badge mpress-badge-note">Windows</span> <span class="mpress-badge mpress-badge-success">macOS</span>

<br/>

Wails menyediakan API untuk mengontrol tampilan dan fungsi kontrol jendela. Fungsionalitas ini tersedia di Windows dan macOS, tetapi tidak di Linux.

## Mengatur Status Tombol Jendela

Status tombol ditentukan oleh enum `ButtonState`:

```go
type ButtonState int

const (
    ButtonEnabled   ButtonState = 0
    ButtonDisabled  ButtonState = 1
    ButtonHidden    ButtonState = 2
)
```

- `ButtonEnabled`: Tombol diaktifkan dan terlihat.
- `ButtonDisabled`: Tombol terlihat, tetapi dinonaktifkan (berwarna abu-abu).
- `ButtonHidden`: Tombol disembunyikan dari bilah judul.

Status tombol dapat diatur saat jendela dibuat atau pada saat runtime.

### Mengatur Status Tombol Saat Membuat Jendela

Saat membuat jendela baru, Anda dapat mengatur status awal tombol menggunakan struct `WebviewWindowOptions`:

```go {title="main.go"}
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name:        "My Application",
    })

    app.Window.NewWithOptions(application.WebviewWindowOptions{
        MinimiseButtonState:   application.ButtonHidden,
        MaximiseButtonState:   application.ButtonDisabled,
        CloseButtonState:      application.ButtonEnabled,
        FullscreenButtonState: application.ButtonEnabled,
    })

    app.Run()
}
```

Dalam contoh di atas, tombol minimalkan disembunyikan, tombol maksimalkan tidak aktif (berwarna abu-abu), dan tombol tutup aktif.

### Mengatur Status Tombol pada Saat Runtime

Anda juga dapat mengubah status tombol pada saat runtime menggunakan metode berikut pada antarmuka `Window`:

```go
window.SetMinimiseButtonState(wails.ButtonHidden)
window.SetMaximiseButtonState(wails.ButtonEnabled)
window.SetCloseButtonState(wails.ButtonDisabled)
window.SetFullscreenButtonState(wails.ButtonEnabled)
```

### macOS: MaximiseButtonState dan FullscreenButtonState Berbagi Tombol yang Sama

Di macOS, tombol lampu lalu lintas hijau (`NSWindowZoomButton`) merupakan kontrol fisik yang sama untuk fungsi maksimalkan dan layar penuh. Tanpa mekanisme pencegahan ini, jika `MaximiseButtonState` dan `FullscreenButtonState` diatur ke nilai yang berbeda saat membuat jendela, pengaturan terakhir akan menggantikan pengaturan sebelumnya tanpa pemberitahuan.

Untuk menghindari hal ini, saat inisialisasi Wails menerapkan status yang **lebih ketat** dari kedua status tersebut, dengan urutan `ButtonEnabled` < `ButtonDisabled` < `ButtonHidden`.

| `MaximiseButtonState` | `FullscreenButtonState` | Status efektif di macOS |
| --- | --- | --- |
| `ButtonEnabled` | `ButtonEnabled` | `ButtonEnabled` |
| `ButtonDisabled` | `ButtonEnabled` | `ButtonDisabled` |
| `ButtonEnabled` | `ButtonHidden` | `ButtonHidden` |
| `ButtonDisabled` | `ButtonHidden` | `ButtonHidden` |

Pada saat runtime, `SetMaximiseButtonState` dan `SetFullscreenButtonState` sama-sama menargetkan `NSWindowZoomButton` di macOS, sehingga pemanggilan terakhir akan berlaku.

### Perbedaan Antarplatform

Fungsionalitas status tombol berperilaku sedikit berbeda di Windows dan macOS:

|  | Windows | Mac |
| --- | --- | --- |
| Nonaktifkan Minimalkan/Maksimalkan/Tutup | Menonaktifkan Minimalkan/Maksimalkan/Tutup | Menonaktifkan Minimalkan/Maksimalkan/Tutup |
| Sembunyikan Minimalkan | Menonaktifkan Minimalkan | Menyembunyikan tombol Minimalkan |
| Sembunyikan Maksimalkan | Menonaktifkan Maksimalkan | Menyembunyikan tombol Maksimalkan |
| Sembunyikan Tutup | Menyembunyikan semua kontrol | Menyembunyikan Tutup |
| `FullscreenButtonState` | Tidak melakukan apa pun | Menargetkan tombol zoom (hijau) |

Catatan: Di Windows, tombol Min/Maks tidak dapat disembunyikan satu per satu. Namun, menonaktifkan keduanya akan menyembunyikan kedua kontrol tersebut dan hanya menampilkan tombol tutup. Windows tidak memiliki tombol layar penuh khusus pada bilah judul standar, sehingga `FullscreenButtonState` tidak melakukan apa pun di sana.

### Mengontrol Gaya Jendela (Windows)

Untuk mengontrol gaya bilah judul di Windows, Anda dapat menggunakan bidang `ExStyle` dalam struct `WebviewWindowOptions`:

Contoh:

```go {title="main.go"}
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/w32"
)

func main() {
    app := application.New(application.Options{
        Name: "My Application",
    })

    app.Window.NewWithOptions(application.WebviewWindowOptions{
        Windows: application.WindowsWindow{
            ExStyle: w32.WS_EX_TOOLWINDOW | w32.WS_EX_NOREDIRECTIONBITMAP | w32.WS_EX_TOPMOST,
        },
    })

    app.Run()
}
```

Opsi lain yang memengaruhi Gaya Diperluas jendela akan ditimpa oleh pengaturan ini:

- HiddenOnTaskbar
- AlwaysOnTop
- IgnoreMouseEvents
- BackgroundType
