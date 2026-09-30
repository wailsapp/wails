---
title: "Instans Tunggal"
description: "Membatasi aplikasi Anda agar hanya menjalankan satu instans"
slug: "guides/single-instance"
sourcePath: "guides/single-instance.md"
---

Penguncian instans tunggal adalah mekanisme yang mencegah beberapa instans aplikasi Anda berjalan secara bersamaan. Mekanisme ini berguna untuk aplikasi yang dirancang untuk membuka file dari baris perintah atau penjelajah file OS.

## Penggunaan

Untuk mengaktifkan fungsi instans tunggal di aplikasi Anda, berikan struct `SingleInstanceOptions` saat membuat aplikasi:

```go
app := application.New(application.Options{
    // ... other options ...
    SingleInstance: &application.SingleInstanceOptions{
        UniqueID: "com.myapp.unique-id",
        OnSecondInstanceLaunch: func(data application.SecondInstanceData) {
            log.Printf("Second instance launched with args: %v", data.Args)
            log.Printf("Working directory: %s", data.WorkingDir)
            log.Printf("Additional data: %v", data.AdditionalData)
        },
        // Optional: Pass additional data to second instance
        AdditionalData: map[string]string{
            "launchtime": time.Now().String(),
        },
    },
})
```

Struct `SingleInstanceOptions` memiliki bidang berikut:

- `UniqueID`: Pengidentifikasi unik untuk aplikasi Anda. Nilainya sebaiknya berupa string unik, biasanya dalam notasi domain terbalik (misalnya, "com.company.appname").
- `EncryptionKey`: Array 32 byte opsional untuk mengenkripsi data yang diteruskan antarinstans menggunakan AES-256-GCM. Jika diberikan sebagai array yang tidak seluruh elemennya nol, semua komunikasi antarinstans akan dienkripsi.
- `OnSecondInstanceLaunch`: Fungsi callback yang dipanggil saat instans kedua aplikasi Anda diluncurkan. Callback menerima struct `SecondInstanceData` yang berisi:
  - `Args`: Argumen baris perintah yang diteruskan ke instans kedua
  - `WorkingDir`: Direktori kerja instans kedua
  - `AdditionalData`: Data tambahan apa pun yang diteruskan dari instans kedua (jika diberikan)

- `AdditionalData`: Map opsional berisi pasangan kunci-nilai string yang akan diteruskan ke instans pertama saat instans berikutnya diluncurkan

@note{type="danger" title="Peringatan"}
Fitur Instans Tunggal menerapkan protokol enkripsi opsional menggunakan AES-256-GCM. Jika enkripsi tidak diaktifkan, data yang diteruskan antarinstans tidak aman. Saat menggunakan fitur instans tunggal tanpa enkripsi, aplikasi Anda sebaiknya memperlakukan setiap data yang diteruskan kepadanya melalui callback instans kedua sebagai data yang tidak tepercaya. Anda sebaiknya memverifikasi bahwa argumen yang diterima valid dan tidak mengandung data berbahaya.

@end

### Komunikasi Aman

Untuk mengaktifkan komunikasi aman antarinstans, berikan kunci enkripsi sepanjang 32 byte. Kunci ini harus sama untuk semua instans aplikasi Anda:

```go
// Define your encryption key (must be exactly 32 bytes)
var encryptionKey = [32]byte{
    0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
    0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
    0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
    0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
}

// Use the key in SingleInstanceOptions
SingleInstance: &application.SingleInstanceOptions{
    UniqueID: "com.myapp.unique-id",
    // Enable encryption for instance communication
    EncryptionKey: encryptionKey,
    // ... other options ...
}
```

@note{type="tip" title="Praktik Terbaik Keamanan"}
- Gunakan kunci unik untuk aplikasi Anda
- Simpan kunci dengan aman jika memuatnya dari konfigurasi
- Jangan gunakan kunci contoh yang ditampilkan di atas—buat kunci Anda sendiri!

@end

### Pengelolaan Jendela

Saat menangani peluncuran instans kedua, Anda sering kali ingin membawa jendela aplikasi ke depan. Anda dapat melakukannya menggunakan metode `Focus()` milik jendela. Jika jendela diminimalkan, Anda mungkin perlu memulihkannya terlebih dahulu:

```go

    var mainWindow *application.WebviewWindow

    SingleInstance: &application.SingleInstanceOptions{
        // Other options...
        OnSecondInstanceLaunch: func(data application.SecondInstanceData) {
            // Focus the window if needed
            if mainWindow != nil {
                mainWindow.Restore()
                mainWindow.Focus()
            }
        },
    }
```

## Cara kerjanya

@tabs{sync-key="platform"}
[Mac]
Penguncian instans tunggal menggunakan mutex bernama. Nama mutex dibuat dari ID unik yang Anda berikan. Data diteruskan ke instans pertama melalui [NSDistributedNotificationCenter](https://developer.apple.com/documentation/foundation/nsdistributednotificationcenter)

[Windows]
Penguncian instans tunggal menggunakan mutex bernama. Nama mutex dibuat dari ID unik yang Anda berikan. Data diteruskan ke instans pertama melalui jendela bersama menggunakan [SendMessage](https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-sendmessage)

[Linux]
Penguncian instans tunggal menggunakan [dbus](https://www.freedesktop.org/wiki/Software/dbus/). Nama dbus dibuat dari ID unik yang Anda berikan. Data diteruskan ke instans pertama melalui [dbus](https://www.freedesktop.org/wiki/Software/dbus/)

@end
