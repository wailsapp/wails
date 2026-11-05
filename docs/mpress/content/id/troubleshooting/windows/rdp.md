---
title: "WebView2 macet melalui Remote Desktop (RDP)"
description: "Atasi macetnya UI WebView2 selama beberapa detik saat aplikasi Wails berjalan melalui sesi RDP yang mengubah DPI monitor di tengah sesi."
slug: "troubleshooting/windows/rdp"
sourcePath: "troubleshooting/windows/rdp.md"
---

## Masalah

Saat aplikasi Wails digunakan melalui sesi Remote Desktop (RDP), UI dapat macet selama beberapa detik ketika terjadi interaksi umum berikut:

- Membuka jendela pop-up memerlukan waktu sekitar 4 hingga 8 detik sejak diklik sampai kontennya terlihat.
- Menutup jendela memblokir jendela induk selama sekitar 2 detik.
- Kondisi lambat ini tetap berlanjut setelah koneksi tersambung kembali dan baru hilang setelah mesin host dimulai ulang.

Masalah ini paling sering terjadi pada klien Microsoft Remote Desktop di iOS, yang menyediakan monitor virtual yang dioptimalkan untuk Retina di tengah sesi. Klien RDP apa pun yang menghadirkan monitor dengan konteks DPI berbeda di tengah sesi dapat memicu perilaku yang sama.

## Penyebabnya

Secara default, WebView2 menggunakan hosting berjendela, yaitu permukaan compositornya berada di jendela anak. Saat klien RDP menghadirkan monitor dengan konteks DPI yang berbeda dari konteks sesi, setiap panggilan pengontrol WebView2 (`PutIsVisible`, `MoveFocus`, penggambaran pertama, dan pelepasan permukaan) memaksa proses re-marshalling DirectComposition secara sinkron. Setiap proses re-marshalling memblokir thread UI selama sekitar 2 detik. Karena itu, waktu macet terakumulasi pada aplikasi yang banyak menggunakan pop-up.

Aplikasi Win32 native yang menggunakan WebView2 pada mesin yang sama tidak terpengaruh karena mengandalkan hosting visual. Hal ini menunjukkan bahwa penyebabnya adalah mode hosting, bukan masalah umum pada WebView2 atau compositor Windows.

## Solusi

Aktifkan hosting visual dengan menetapkan `UseVisualHosting` dalam opsi Windows Anda. Dengan hosting visual, permukaan compositor WebView2 dimiliki melalui visual DirectComposition yang dimiliki host, sehingga perubahan konteks DPI tidak lagi memicu proses re-marshalling sinkron.

```go
package main

import "github.com/wailsapp/wails/v3/pkg/application"

func main() {
    app := application.New(application.Options{
        Windows: application.WindowsOptions{
            UseVisualHosting: true,
        },
    })

    // ... create your windows, then:
    app.Run()
}
```

Setelah opsi ini diaktifkan, pop-up terbuka dalam waktu navigasi normal (sekitar 150 hingga 500 ms), dan menutup jendela tidak lagi memblokir jendela induk.

@note{type="caution"}
`UseVisualHosting` harus ditetapkan sebelum `app.Run()`. Wails membacanya saat aplikasi dimulai dan menetapkan variabel lingkungan `COREWEBVIEW2_FORCED_HOSTING_MODE` ke `COREWEBVIEW2_HOSTING_MODE_WINDOW_TO_VISUAL` sebelum lingkungan WebView2 diinisialisasi. Menetapkannya setelah itu tidak akan berpengaruh.

@end

Nilai default opsi ini adalah `false`, sehingga hosting berjendela tetap menjadi mode default. Perilaku aplikasi yang sudah ada tidak berubah kecuali opsi ini diaktifkan secara eksplisit.

## Kapan perlu mengaktifkannya

Tetapkan `UseVisualHosting: true` jika aplikasi Anda sering digunakan melalui RDP, terutama melalui klien Microsoft Remote Desktop di iOS, dan mengalami macet selama beberapa detik saat membuka atau menutup jendela. Jika aplikasi Anda tidak berjalan melalui RDP, opsi ini tidak diperlukan dan Anda dapat mempertahankan nilai default.

## Referensi

- [WebView2: hosting berjendela vs. hosting visual](https://learn.microsoft.com/en-us/microsoft-edge/webview2/concepts/windowed-vs-visual-hosting)
- [Masalah WebView2Feedback #5248](https://github.com/MicrosoftEdge/WebView2Feedback/issues/5248)
- [Masalah WebView2Feedback #4485](https://github.com/MicrosoftEdge/WebView2Feedback/issues/4485)
