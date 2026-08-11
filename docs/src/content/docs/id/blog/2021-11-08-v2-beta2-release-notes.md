---
slug: blog/wails-v2-beta-for-mac
title: Wails v2 Beta untuk MacOS
description: "Catatan rilis dan pengumuman untuk Wails"
authors: [leaanthony]
tags: [wails, v2]
date: 2021-11-08
---

![wails-mac screenshot](../../../../assets/blog-images/wails-mac.webp)

Hari ini menandai rilis beta pertama Wails v2 untuk Mac! Butuh cukup lama
untuk sampai ke titik ini dan saya berharap rilis hari ini akan memberi Anda sesuatu
yang cukup berguna. Ada sejumlah liku-liku untuk sampai ke
titik ini dan saya berharap, dengan bantuan Anda, untuk menghilangkan kerutan dan memoles
port Mac untuk rilis v2 final.

Maksud Anda ini belum siap untuk production? Untuk use case Anda, mungkin sudah
siap, tetapi masih ada sejumlah masalah yang diketahui jadi pantau terus
[project board ini](https://github.com/wailsapp/wails/projects/7) dan jika Anda
ingin berkontribusi, Anda sangat welcome!

Jadi apa yang baru untuk Wails v2 untuk Mac vs v1? Hint: Cukup mirip dengan
Windows Beta :wink:

### Fitur Baru

![wails-menus-mac screenshot](../../../../assets/blog-images/wails-menus-mac.webp)

Ada banyak permintaan untuk dukungan menu native. Wails akhirnya
mengcover kebutuhan Anda. Menu aplikasi sekarang tersedia dan mencakup dukungan untuk sebagian besar fitur menu
native. Ini termasuk item menu standar, checkbox, radio group,
submenu, dan separator.

Ada sejumlah besar permintaan di v1 untuk kemampuan memiliki kontrol
lebih besar terhadap window itu sendiri. Saya dengan senang hati mengumumkan bahwa ada API runtime
baru khusus untuk ini. Fitur-rich dan mendukung konfigurasi multi-monitor.
Ada juga API dialog yang ditingkatkan: Sekarang, Anda dapat memiliki dialog
native modern dengan konfigurasi kaya untuk memenuhi semua kebutuhan dialog Anda.

### Opsi Khusus Mac

Selain opsi aplikasi normal, Wails v2 untuk Mac juga membawa
ekstra Mac:

- Buat window Anda funky dan translucent, seperti aplikasi swift yang cantik!
- Titlebar yang sangat dapat dikustomisasi
- Kami mendukung opsi NSAppearance untuk aplikasi
- Konfigurasi sederhana untuk auto-create menu "About"

### Tidak perlu membundel aset

Pain point besar v1 adalah kebutuhan untuk meringkas seluruh aplikasi menjadi
file JS & CSS tunggal. Saya dengan senang hati mengumumkan bahwa untuk v2, tidak ada
persyaratan untuk membundel aset, dalam bentuk apapun. Ingin memuat gambar
lokal? Gunakan tag `<img>` dengan path src lokal. Ingin menggunakan font keren? Salin
dan tambahkan path-nya di CSS Anda.

> Wow, itu terdengar seperti webserver...

Ya, berfungsi seperti webserver, kecuali memang bukan webserver.

> Jadi bagaimana cara memasukkan aset saya?

Anda cukup meneruskan satu `embed.FS` yang berisi semua aset ke
konfigurasi aplikasi Anda. Mereka bahkan tidak perlu berada di direktori teratas -
Wails akan menanganinya untuk Anda.

### Pengalaman Pengembangan Baru

Sekarang aset tidak perlu dibundel, ini memungkinkan pengalaman pengembangan
baru sepenuhnya. Perintah `wails dev` yang baru akan build dan menjalankan aplikasi Anda, tetapi
alih-alih menggunakan aset di `embed.FS`, aset dimuat langsung dari disk.

Perintah ini juga menyediakan fitur tambahan:

- Hot reload - Perubahan apapun pada aset frontend akan memicu dan auto reload
  frontend aplikasi
- Auto rebuild - Perubahan apapun pada kode Go Anda akan rebuild dan meluncurkan ulang
  aplikasi Anda

Selain itu, webserver akan dimulai di port 34115. Ini akan melayani
aplikasi Anda ke browser apapun yang terhubung. Semua web browser yang terhubung akan
merespons event sistem seperti hot reload saat aset berubah.

Di Go, kita terbiasa menangani struct dalam aplikasi kita. Seringkali
berguna mengirim struct ke frontend kita dan menggunakannya sebagai state dalam aplikasi.
Di v1, ini adalah proses yang sangat manual dan sedikit membebani developer.
Saya dengan senang hati mengumumkan bahwa di v2, aplikasi apapun yang dijalankan dalam mode dev akan
secara otomatis menghasilkan model TypeScript untuk semua struct yang menjadi parameter input atau
output dari bound method. Ini memungkinkan pertukaran model data yang mulus
antara kedua dunia.

Selain itu, modul JS lain dihasilkan secara dinamis yang membungkus semua
bound method Anda. Ini menyediakan JSDoc untuk method Anda, memberikan code
completion dan hinting di IDE Anda. Sangat keren ketika Anda mendapatkan model data
auto-imported saat menekan tab di modul auto-generated yang membungkus kode Go
Anda!

### Remote Templates

![remote-mac screenshot](../../../../assets/blog-images/remote-mac.webp)

Membuat aplikasi up and running dengan cepat selalu menjadi tujuan utama proyek
Wails. Ketika kami meluncurkan, kami mencoba mencakup banyak framework modern
saat itu: react, vue dan angular. Dunia pengembangan frontend
sangat opinionated, bergerak cepat dan sulit diikuti! Akibatnya,
kami menemukan template dasar kami menjadi outdated cukup cepat dan ini
menyebabkan headache maintenance. Ini juga berarti kami tidak memiliki template modern keren
untuk tech stack terbaru dan terhebat.

Dengan v2, saya ingin memberdayakan komunitas dengan memberi Anda kemampuan untuk membuat
dan hosting template sendiri, alih-alih bergantung pada proyek Wails. Jadi sekarang Anda
dapat membuat proyek menggunakan template yang didukung komunitas! Saya harap ini akan
menginspirasi developer untuk membuat ekosistem template proyek yang vibrant. Saya
benar-benar antusias dengan apa yang dapat diciptakan komunitas developer kami!

### Dukungan M1 Native

Berkat dukungan luar biasa dari [Mat Ryer](https://github.com/matryer/), proyek
Wails sekarang mendukung build M1 native:

![build-darwin-arm screenshot](../../../../assets/blog-images/build-darwin-arm.webp)

Anda juga dapat menentukan `darwin/amd64` sebagai target:

![build-darwin-amd screenshot](../../../../assets/blog-images/build-darwin-amd.webp)

Oh, hampir lupa.... Anda juga bisa melakukan `darwin/universal`.... :wink:

![build-darwin-universal screenshot](../../../../assets/blog-images/build-darwin-universal.webp)

### Cross Compilation ke Windows

Karena Wails v2 untuk Windows murni Go, Anda dapat menargetkan build Windows tanpa
docker.

![build-cross-windows screenshot](../../../../assets/blog-images/build-cross-windows.webp)

### Renderer WKWebView

V1 bergantung pada komponen WebView (kini deprecated). V2 menggunakan komponen
WKWebKit terbaru jadi harapkan yang terbaik dan terbaru dari Apple.

### Kesimpulan

Seperti yang saya katakan di catatan rilis Windows, Wails v2 mewakili fondasi baru
untuk proyek ini. Tujuan rilis ini adalah mendapatkan feedback tentang pendekatan baru,
dan menghilangkan bug apapun sebelum rilis penuh. Input Anda sangat
welcome! Silakan arahkan feedback apapun ke
discussion board [v2 Beta](https://github.com/wailsapp/wails/discussions/828).

Dan akhirnya, saya ingin memberikan terima kasih khusus kepada semua
[sponsor proyek](/id/credits#sponsors), termasuk
[JetBrains](https://www.jetbrains.com?from=Wails), yang dukungannya mendorong
proyek dalam banyak hal di balik layar.

Saya menantikan melihat apa yang dibangun orang dengan Wails dalam fase
proyek yang menarik ini!

Lea.

PS: Pengguna Linux, giliran Anda berikutnya!

PPS: Jika Anda atau perusahaan Anda merasa Wails berguna, pertimbangkan
[menyponsori proyek ini](https://github.com/sponsors/leaanthony). Terima kasih!
