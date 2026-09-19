---
title: "Menangani Panic"
description: "Cara menangani panic dalam aplikasi Wails Anda"
slug: "guides/panic-handling"
sourcePath: "guides/panic-handling.md"
---

Dalam aplikasi Go, panic dapat terjadi saat runtime ketika sesuatu yang tidak terduga terjadi. Panduan ini menjelaskan cara menangani panic baik dalam kode Go secara umum maupun secara khusus dalam aplikasi Wails Anda.

## Memahami Panic di Go

Sebelum membahas penanganan panic khusus Wails, penting untuk memahami cara kerja panic di Go:

1. Panic digunakan untuk galat yang tidak dapat dipulihkan dan seharusnya tidak terjadi selama operasi normal
2. Ketika panic terjadi dalam sebuah goroutine, hanya goroutine tersebut yang terpengaruh
3. Panic dapat dipulihkan menggunakan `defer` dan `recover()`

Berikut adalah contoh dasar penanganan panic di Go:

```go
func doSomething() {
    // Deferred functions run even when a panic occurs
    defer func() {
        if r := recover(); r != nil {
            fmt.Printf("Recovered from panic: %v\n", r)
        }
    }()
    
    // Your code that might panic
    panic("something went wrong")
}
```

Untuk informasi lebih terperinci tentang panic dan recover di Go, lihat [Blog Go: Defer, Panic, dan Recover](https://go.dev/blog/defer-panic-and-recover).

## Penanganan Panic di Wails

Wails secara otomatis menangani panic yang terjadi dalam metode Service Anda ketika metode tersebut dipanggil dari frontend. Artinya, Anda tidak perlu menambahkan pemulihan panic ke metode ini—Wails akan menangkap panic tersebut dan memprosesnya melalui penangan panic yang telah Anda konfigurasi.

Penangan panic dirancang secara khusus untuk menangkap:

- Panic dalam metode service terikat yang dipanggil dari frontend
- Panic internal dari runtime Wails

Untuk skenario lain, seperti goroutine latar belakang atau kode Go mandiri, Anda sebaiknya menangani panic sendiri menggunakan mekanisme pemulihan panic standar Go.

## Struct PanicDetails

Ketika panic terjadi, Wails menyimpan informasi penting tentang panic tersebut dalam sebuah struct `PanicDetails`:

```go
type PanicDetails struct {
    StackTrace     string    // The stack trace of where the panic occurred. Potentially trimmed to provide more context
    Error          error     // The error that caused the panic
    Time           time.Time // The time when the panic occurred
    FullStackTrace string    // The complete stack trace including runtime frames
}
```

Struktur ini menyediakan informasi lengkap tentang panic tersebut:

- `StackTrace`: String berformat yang menampilkan tumpukan panggilan yang menyebabkan panic
- `Error`: Galat aktual atau pesan panic
- `Time`: Waktu persis saat panic terjadi
- `FullStackTrace`: Pelacakan tumpukan lengkap, termasuk frame runtime

@note{type="info" title="Panic dalam Kode Service"}
Ketika panic ditangkap dalam kode Service Anda setelah dipanggil dari frontend, pelacakan tumpukan dipangkas agar berfokus tepat pada lokasi terjadinya panic dalam kode Anda. Jika ingin melihat pelacakan tumpukan lengkap, Anda dapat menggunakan field `FullStackTrace`.

@end

## Penangan Panic Bawaan

Jika Anda tidak menentukan penangan panic khusus, Wails akan menggunakan penangan bawaannya, yang menampilkan informasi galat dalam pesan log berformat lalu keluar. Contoh:

```
************************ FATAL ******************************
* There has been a catastrophic failure in your application *
********************* Error Details *************************
panic error: oh no! something went wrong deep in my service! :(
main.(*WindowService).call2
	at E:/wails/v3/examples/panic-handling/main.go:23
main.(*WindowService).call1
	at E:/wails/v3/examples/panic-handling/main.go:19
main.(*WindowService).GeneratePanic
	at E:/wails/v3/examples/panic-handling/main.go:15
*************************************************************
```

## Penangan Panic Khusus

Anda dapat mengimplementasikan penangan panic sendiri dengan menetapkan opsi `PanicHandler` saat membuat aplikasi. Berikut contohnya:

```go
app := application.New(application.Options{
    Name: "My App",
    PanicHandler: func(panicDetails *application.PanicDetails) {
        fmt.Printf("*** Custom Panic Handler ***\n")
        fmt.Printf("Time: %s\n", panicDetails.Time)
        fmt.Printf("Error: %s\n", panicDetails.Error)
        fmt.Printf("Stacktrace: %s\n", panicDetails.StackTrace)
        fmt.Printf("Full Stacktrace: %s\n", panicDetails.FullStackTrace)
        
        // You could also:
        // - Log to a file
        // - Send to a crash reporting service
        // - Show a user-friendly error dialog
        // - Attempt to recover or restart the application
    },
})
```

## Menangani Panic di Goroutine Anda Sendiri {#user-goroutines}

Wails hanya dapat memasang pemulihan tertunda pada lokasi yang dikendalikannya — metode layanan terikat yang dipanggil dari frontend, callback runtime internal, dan sebagainya. Jika kode Anda sendiri melakukan ini:

```go
go func() {
    // your work
}()
```

Wails tidak dapat menyisipkan `defer handlePanic()` ke goroutine tersebut. Jika terjadi panic, seluruh proses berhenti sesuai perilaku default Go — `PanicHandler` yang Anda daftarkan **tidak** dipanggil.

Untuk mengarahkan panic goroutine pengguna ke handler yang sama dengan panic yang ditangkap Wails, tambahkan fungsi pembantu kecil yang membuat `PanicDetails` secara manual dan memanggil fungsi yang Anda daftarkan sebagai `PanicHandler`:

```go
import (
    "fmt"
    "runtime/debug"
    "time"

    "github.com/wailsapp/wails/v3/pkg/application"
)

// reportPanic is the function you register as application.Options.PanicHandler.
func reportPanic(pd *application.PanicDetails) {
    // log to file / send to Sentry / show dialog / etc.
}

// recoverAndReport funnels goroutine panics to reportPanic. Defer it as the
// first statement of every goroutine you spawn in user code.
func recoverAndReport() {
    r := recover()
    if r == nil {
        return
    }
    err, ok := r.(error)
    if !ok {
        err = fmt.Errorf("%v", r)
    }
    stack := string(debug.Stack())
    reportPanic(&application.PanicDetails{
        Error:          err,
        Time:           time.Now(),
        StackTrace:     stack,
        FullStackTrace: stack,
    })
}
```

Gunakan di awal setiap goroutine yang Anda kelola:

```go
go func() {
    defer recoverAndReport()
    // your work
}()
```

Jalur panic yang ditangkap Wails maupun jalur goroutine pengguna kini berakhir di `reportPanic`, sehingga pelaporan tetap terpusat.

Wails menyertakan dua contoh yang dapat dijalankan:

- `v3/examples/panic-handling` — minimal: hanya panic pada metode terikat, diarahkan melalui `PanicHandler`.
- `v3/examples/user-panic-handling` — panic pada metode terikat dan goroutine latar belakang mengalir melalui handler yang sama.

@note{type="caution" title="Keutuhan stack trace"}

Wails memangkas `PanicDetails.StackTrace` untuk menyembunyikan frame pembungkusnya sendiri agar kode Anda berada di bagian atas trace. Saat Anda membuat `PanicDetails` sendiri dari `runtime/debug.Stack()`, tidak ada pemangkasan — `StackTrace` dan `FullStackTrace` akan identik dan menyertakan seluruh stack goroutine. Biasanya inilah yang diperlukan untuk goroutine yang Anda buat.

@end

## Mengumpulkan Diagnostik Saat Panic {#panic-diagnostics}

`PanicHandler` menerima panic itu sendiri, tetapi sering kali Anda juga perlu menangkap seluruh keadaan proses — informasi sistem, informasi build, keadaan proses/memori/modul, serta minidump di Windows — agar dapat merekonstruksi kejadian nanti.

Wails menyediakan paket khusus untuk itu: [Men-debug Crash](/guides/debugging-crashes/). Integrasi umumnya:

```go
import "github.com/wailsapp/wails/v3/pkg/debug"

app := application.New(application.Options{
    PanicHandler: func(pd *application.PanicDetails) {
        report, _ := debug.Report(debug.WithDump())
        // pd holds the wails-side panic info; report adds system context
        // and (on Windows) a minidump at report.DumpPath.
        mycrashservice.Upload(pd, report)
    },
})
```

`debug.Report` harus diaktifkan secara eksplisit — Wails tidak akan menyimpan dump atau snapshot sistem secara otomatis karena pelaporan crash sering memerlukan persetujuan pengguna atau penyamaran informasi identitas pribadi. Lihat panduan [Men-debug Crash](/guides/debugging-crashes/) untuk API lengkap.

## Catatan Akhir

Ingat bahwa penangan panic Wails ditujukan secara khusus untuk mengelola panic dalam metode terikat dan galat runtime internal. Untuk bagian lain aplikasi Anda, gunakan pola penanganan galat standar Go dan mekanisme pemulihan panic jika sesuai. Seperti pada semua aplikasi Go, sebaiknya cegah panic melalui penanganan galat yang tepat jika memungkinkan.
