---
title: "Jendela Notch"
description: "Buat jendela macOS native yang melekat pada casing kamera."
slug: "features/windows/notch-windows"
sourcePath: "features/windows/notch-windows.md"
---

`NewNotchWindow` membuat panel macOS berbentuk yang tidak mengaktifkan aplikasi dan melekat pada casing kamera. Wails mengelola penempatan native, sayap sambungan hitam, kanvas luar transparan, level jendela, perilaku Spaces, serta animasi tampil dan sembunyi opsional. Konten web Anda hanya menempati persegi panjang bagian dalam yang diminta. Memindahkan penunjuk ke dalam jendela akan segera menjadikan webview-nya sebagai jendela utama untuk input tanpa mengaktifkan aplikasi.

@note{type="caution" title="Transparansi webview menggunakan API privat"}
`NewNotchWindow` memerlukan `-tags private_mac_apis` agar konten web transparan dapat menampilkan bentuk notch native. Tanpa tag tersebut, panel, penempatan, dan animasi tetap berfungsi, tetapi webview tetap buram. Lihat [API macOS privat](/guides/build/private-macos-apis/#webview-transparency-and-background).

@end

<figure>
  <img
    src="/images/notch-notification.gif"
    alt="Notifikasi notch Wails yang bergeser turun dari casing kamera MacBook, menampilkan metrik sistem secara langsung, lalu kembali tersembunyi di bawah notch"
    loading="lazy"
    decoding="async"
    style="width: 100%; border-radius: 0.75rem"
  />
  <figcaption>
    Notifikasi notch beranimasi yang menggunakan webview persisten dengan transisi tampil dan sembunyi native.
  </figcaption>
</figure>

```go
alert := app.Window.NewNotchWindow(application.NotchWindowOptions{
    Width:    660,
    Height:   92,
    Animated: true,
    WindowOptions: application.WebviewWindowOptions{
        Name: "alert",
        URL:  "/alert",
    },
})

alert.Show()
alert.Hide()
visible := alert.Visibility()
alert.Close()
```

## Opsi

| Bidang | Tipe | Default | Deskripsi |
| --- | --- | --- | --- |
| `Width` | `int` | `660` | Lebar webview yang dapat digunakan di dalam tepi berbentuk native. |
| `Height` | `int` | `92` | Tinggi webview yang dapat digunakan di dalam tepi berbentuk native. |
| `Animated` | `bool` | `false` | Menggeser jendela ke bawah saat `Show` dan ke atas saat `Hide`. |
| `AnimationSpeed` | `time.Duration` | `420ms` | Durasi tampil. Durasi sembunyi menggunakan dua pertiga dari nilai ini, secara default `280ms`. |
| `Screen` | `*Screen` | layar utama | Menargetkan layar tertentu. Jika tidak ditentukan, Wails menggunakan layar utama. |
| `WindowOptions` | `WebviewWindowOptions` | nilai default | Menyediakan nama, URL atau HTML, CSS, JavaScript, pengikatan tombol, dan perilaku webview lainnya. |

`NewNotchWindow` mengelola ukuran luar, posisi, bingkai, transparansi, kebijakan pengubahan ukuran, kelas panel native, level jendela, dan perilaku koleksi. Nilai untuk bidang-bidang tersebut dalam `WindowOptions` sengaja diganti. Bidang lainnya dipertahankan. Penyeretan latar belakang native dan area seret CSS dinonaktifkan agar jendela tetap melekat pada casing kamera.

`NotchWindow` yang dikembalikan sengaja hanya mengekspos `Show`, `Hide`, `Visibility`, dan `Close`; geometri native tidak dapat diubah melalui handle tingkat tinggi.

@note{type="note"}
Jendela notch memerlukan macOS. Pada Mac tanpa casing kamera, Wails menempatkan jendela di tengah atas, di bawah bilah menu. Pada platform yang tidak didukung, `NewNotchWindow` mengembalikan handle nonaktif yang metode siklus hidupnya aman dan tidak melakukan apa pun, serta `Visibility`-nya selalu bernilai false.

@end

## Siklus Hidup

- `Show` menampilkan jendela native yang sudah ada. Jika animasi diaktifkan, jendela bergeser turun dari atas layar.
- `Hide` mempertahankan jendela agar tetap aktif untuk digunakan kembali. Jika animasi diaktifkan, jendela bergeser kembali ke atas layar sebelum dikeluarkan dari tampilan. Webview, status JavaScript, pengikatan, dan pemroses peristiwa tetap dimuat, tetapi jendela tersembunyi tidak memiliki target hover; aplikasi harus memanggil `Show` untuk menampilkannya kembali.
- `Visibility` melaporkan visibilitas native saat ini.
- `Close` memusnahkan jendela native secara permanen. Buat jendela baru sebelum menampilkan notifikasi tersebut lagi.
- Saat penunjuk memasuki jendela notch, jendela tersebut dibawa ke depan dan webview-nya difokuskan agar dapat langsung menerima input papan ketik, sembari mempertahankan perilaku panel yang tidak mengaktifkan aplikasi.

Setiap pemanggilan `NewNotchWindow` membuat jendela independen dengan konten, visibilitas, dan status animasinya sendiri. Semua jendela menggunakan level native yang sama, sehingga instans yang paling baru ditampilkan atau dimasuki penunjuk akan muncul di depan dan dapat menutupi instans sebelumnya. Jika jendela menargetkan layar yang berbeda, setiap jendela dipusatkan pada casing kamera layar tersebut. macOS tidak mengoordinasikan jendela notch milik aplikasi yang berbeda; jika beberapa aplikasi menampilkan jendela pada lokasi dan level yang sama, jendela yang terakhir ditempatkan dalam urutan tampilan akan muncul di depan.

Untuk beban kerja notifikasi, aplikasi umumnya menggunakan kembali jendela tersembunyi atau mengelola antreannya sendiri dalam satu webview. Wails tidak memberlakukan antrean, penggantian, penutupan otomatis, ataupun kebijakan satu jendela.

@note{type="note"}
Permukaan diciutkan persisten yang terbuka kembali saat di-hover sengaja dipisahkan dari penyembunyian notifikasi dan dilacak di [#6009](https://github.com/wailsapp/wails/issues/6009).

@end

Lihat [contoh notch-notification](https://github.com/wailsapp/wails/tree/master/v3/examples/notch-notification) untuk monitor sistem ringkas yang status JavaScript aktifnya tetap dipertahankan melalui siklus tampil dan sembunyi notifikasi yang berulang.
