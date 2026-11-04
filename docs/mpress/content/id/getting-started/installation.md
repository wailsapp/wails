---
title: "Instalasi"
description: "Instal Wails dan siapkan lingkungan pengembangan Anda"
slug: "getting-started/installation"
sourcePath: "getting-started/installation.md"
---

## Platform yang Didukung

- Windows AMD64/ARM64
- macOS AMD64/ARM64. Tugas build yang dihasilkan menargetkan macOS 12.0 atau yang lebih baru; ini adalah target deployment untuk build standar, bukan pernyataan bahwa setiap konfigurasi build khusus mendukung sistem yang lebih lama.
- Ubuntu 24.04 AMD64/ARM64 (distribusi Linux lain mungkin juga dapat digunakan!)

## Dependensi

Wails memiliki sejumlah dependensi umum yang harus tersedia sebelum instalasi.

@note{type="tip"}
Setelah menginstal Wails CLI, Anda dapat menjalankan `wails3 setup` untuk memeriksa dependensi ini secara otomatis dan membantu menginstalnya.

@end

@tabs
[Go (Minimal 1.26)]
Wails v3 memilih toolchain Go 1.26 atau yang lebih baru. Peluncur Go versi lama mungkin mengunduhnya secara otomatis dengan `GOTOOLCHAIN=auto`; jika tidak, instal secara eksplisit. Deklarasi versi Go yang lebih rendah dalam modul awal tidak menggantikan dependensi Wails.

Unduh Go dari [Halaman Unduhan Go](https://go.dev/dl/).

Pastikan Anda mengikuti [petunjuk instalasi Go](https://go.dev/doc/install) resmi. Pastikan juga variabel lingkungan `PATH` Anda menyertakan jalur ke direktori `~/go/bin`. Mulai ulang terminal Anda, lalu lakukan pemeriksaan berikut:

- Pastikan Go telah terinstal dengan benar: `go version`
- Pastikan `~/go/bin` tercantum dalam variabel PATH Anda
  - Mac / Linux: `echo $PATH | grep go/bin`
  - Windows: `$env:PATH -split ';' | Where-Object { $_ -like '*\go\bin' }`


[npm (Opsional)]
Meskipun Wails tidak mengharuskan npm terinstal, sebagian besar templat bawaan memerlukannya.

Unduh penginstal Node terbaru dari [Halaman Unduhan Node](https://nodejs.org/en/download/). Sebaiknya gunakan rilis terbaru karena versi itulah yang umumnya kami gunakan dalam pengujian.

Jalankan `npm --version` untuk memverifikasinya.

@note{type="info"}
Jika Anda lebih memilih pengelola paket selain npm, silakan gunakan pengelola tersebut. Anda harus memperbarui Taskfile proyek agar menggunakannya.

@end

@end

## Dependensi Khusus Platform

Anda juga harus menginstal dependensi khusus platform:

@tabs{sync-key="platform"}
[Mac]
Wails mengharuskan alat baris perintah Xcode terinstal. Anda dapat menginstalnya dengan menjalankan:

```sh
xcode-select --install
```

[Windows]
Wails mengharuskan [WebView2 Runtime](https://developer.microsoft.com/en-us/microsoft-edge/webview2/) terinstal. Komponen ini sudah terinstal pada hampir semua instalasi Windows. Anda dapat memeriksanya menggunakan perintah `wails doctor`.

[Linux]
Linux memerlukan alat build standar `gcc` serta `gtk4` dan `webkitgtk-6.0`. Setelah instalasi, jalankan <code>wails3 doctor</code> untuk melihat cara menginstal dependensi tersebut. Stack GTK3 / WebKit2GTK 4.1 versi lama tetap tersedia melalui `-tags gtk3` (lihat [Pemaketan Linux - Dukungan GTK3 Versi Lama](/guides/build/linux/#legacy-gtk3-support)) hingga v3.1. Jika distro/pengelola paket Anda tidak didukung, beri tahu kami di Discord.

@end

## Instalasi

Untuk menginstal Wails CLI menggunakan Go Modules, jalankan perintah berikut:

```shell
go install -v github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0
```

Jika Anda ingin menginstal versi pengembangan terbaru, jalankan perintah berikut:

```shell
git clone https://github.com/wailsapp/wails.git
cd wails
cd v3/cmd/wails3
go install
```

Saat menggunakan versi pengembangan, semua proyek yang dihasilkan akan menggunakan direktif [replace](https://go.dev/ref/mod#go-mod-file-replace) milik Go untuk memastikan proyek menggunakan Wails versi pengembangan.

## Langkah Berikutnya

Setelah menginstal CLI, jalankan wizard penyiapan untuk mengonfigurasi lingkungan pengembangan Anda:

```shell
wails3 setup
```

@note{type="caution" title="Eksperimental"}
Wizard penyiapan ini masih baru dan terutama telah diuji di Linux. Jika mengalami masalah, silakan [laporkan masalah tersebut](https://github.com/wailsapp/wails/issues/4904) dan gunakan langkah-langkah instalasi dependensi manual di bawah ini.

@end

Wizard penyiapan akan:

- Memeriksa dependensi platform dan membantu menginstalnya
- Mengonfigurasi pengaturan default proyek (informasi pembuat, prefiks ID bundle)
- Secara opsional menyiapkan Docker untuk build lintas platform
- Mengonfigurasi penandatanganan kode (jika diperlukan)

Lihat [Panduan Penyiapan](/getting-started/setup/) untuk detail selengkapnya.

## Instalasi Dependensi Manual

Jika Anda lebih memilih menginstal dependensi secara manual, atau jika wizard penyiapan tidak berfungsi di sistem Anda, ikuti petunjuk khusus platform di atas, lalu jalankan:

```shell
wails3 doctor
```

Perintah ini akan memeriksa apakah dependensi yang benar telah terinstal dan memberi tahu komponen yang masih belum tersedia.

## Perintah `wails3` tampaknya tidak tersedia?

Jika sistem Anda melaporkan bahwa perintah `wails3` tidak tersedia, periksa hal-hal berikut:

- Pastikan Anda telah mengikuti **panduan instalasi Go** di atas dengan benar dan direktori `go/bin` tercantum dalam variabel lingkungan `PATH`.
- Tutup dan buka kembali terminal yang sedang digunakan agar variabel `PATH` yang baru diterapkan.
