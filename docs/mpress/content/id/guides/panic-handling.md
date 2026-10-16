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

Untuk contoh lengkap yang dapat dijalankan mengenai penanganan panic dalam aplikasi Wails, lihat contoh panic-handling di `v3/examples/panic-handling`.

## Catatan Akhir

Ingat bahwa penangan panic Wails ditujukan secara khusus untuk mengelola panic dalam metode terikat dan galat runtime internal. Untuk bagian lain aplikasi Anda, gunakan pola penanganan galat standar Go dan mekanisme pemulihan panic jika sesuai. Seperti pada semua aplikasi Go, sebaiknya cegah panic melalui penanganan galat yang tepat jika memungkinkan.
