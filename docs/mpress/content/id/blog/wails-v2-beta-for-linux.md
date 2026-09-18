---
title: "Wails v2 Beta untuk Linux"
description: "Catatan rilis dan pengumuman untuk Wails"
authors: ["leaanthony"]
tags: ["wails","v2"]
date: "2022-02-22"
slug: "blog/wails-v2-beta-for-linux"
image: "/assets/blog-images/wails-linux.webp"
sourcePath: "blog/wails-v2-beta-for-linux.md"
---

![tangkapan layar wails-linux](/assets/blog-images/wails-linux.webp)

Dengan senang hati akhirnya saya mengumumkan bahwa Wails v2 kini tersedia dalam versi beta untuk Linux! Agak ironis bahwa eksperimen paling awal dengan v2 dilakukan di Linux, tetapi Linux justru menjadi platform terakhir yang mendapatkan rilis. Meskipun demikian, v2 yang kita miliki saat ini sangat berbeda dari eksperimen awal tersebut. Jadi, tanpa berpanjang lebar lagi, mari kita bahas fitur-fitur barunya:

## Fitur Baru

![tangkapan layar wails-menus-linux](/assets/blog-images/wails-menus-linux.webp)

Ada banyak permintaan untuk dukungan menu native. Wails akhirnya menyediakannya untuk Anda. Menu aplikasi kini tersedia dan mendukung sebagian besar fitur menu native. Dukungan ini mencakup item menu standar, kotak centang, grup tombol radio, submenu, dan pemisah.

Pada v1, ada sangat banyak permintaan agar pengguna dapat memiliki kendali yang lebih besar atas jendela aplikasi. Dengan senang hati saya mengumumkan bahwa kini tersedia API runtime baru yang dirancang khusus untuk itu. API ini kaya fitur dan mendukung konfigurasi multi-monitor. API dialog juga telah disempurnakan: kini Anda dapat menggunakan dialog native modern dengan konfigurasi lengkap untuk memenuhi semua kebutuhan dialog Anda.

### Tidak perlu membundel aset

Salah satu masalah terbesar pada v1 adalah keharusan memadatkan seluruh aplikasi Anda menjadi satu file JS dan satu file CSS. Dengan senang hati saya mengumumkan bahwa pada v2, aset sama sekali tidak perlu dibundel dalam bentuk apa pun. Ingin memuat gambar lokal? Gunakan tag `<../../../assets/blog-images>` dengan jalur src lokal. Ingin menggunakan font yang keren? Salin font tersebut lalu tambahkan jalurnya ke CSS Anda.

> Wah, kedengarannya seperti server web...

Ya, cara kerjanya memang seperti server web, tetapi sebenarnya bukan.

> Jadi, bagaimana cara menyertakan aset saya?

Anda cukup meneruskan satu `embed.FS` yang berisi semua aset Anda ke konfigurasi aplikasi. Aset tersebut bahkan tidak harus berada di direktori teratas—Wails akan menanganinya untuk Anda.

### Pengalaman Pengembangan Baru

Karena aset kini tidak perlu dibundel, tersedia pengalaman pengembangan yang benar-benar baru. Perintah `wails dev` yang baru akan membangun dan menjalankan aplikasi Anda, tetapi alih-alih menggunakan aset di dalam `embed.FS`, perintah tersebut memuatnya langsung dari disk.

Perintah ini juga menyediakan fitur tambahan berikut:

- Hot reload—Setiap perubahan pada aset frontend akan memicu pemuatan ulang frontend aplikasi secara otomatis
- Pembangunan ulang otomatis—Setiap perubahan pada kode Go Anda akan membangun ulang dan menjalankan kembali aplikasi

Selain itu, server web akan dijalankan pada port 34115. Server ini akan menyajikan aplikasi Anda ke setiap browser yang terhubung. Semua browser web yang terhubung akan merespons peristiwa sistem, seperti hot reload saat aset berubah.

Dalam Go, kita terbiasa menggunakan struct di aplikasi. Sering kali, mengirim struct ke frontend dan menggunakannya sebagai state aplikasi sangatlah berguna. Pada v1, proses ini harus dilakukan secara manual dan cukup membebani developer. Dengan senang hati saya mengumumkan bahwa pada v2, setiap aplikasi yang dijalankan dalam mode dev akan otomatis menghasilkan model TypeScript untuk semua struct yang menjadi parameter input atau output metode yang diikat. Hal ini memungkinkan pertukaran model data yang lancar antara kedua lingkungan tersebut.

Selain itu, modul JS lain akan dibuat secara dinamis untuk membungkus semua metode yang diikat. Modul ini menyediakan JSDoc untuk metode Anda sehingga IDE dapat memberikan pelengkapan dan petunjuk kode. Rasanya benar-benar keren ketika model data diimpor secara otomatis saat Anda menekan Tab dalam modul yang dibuat otomatis untuk membungkus kode Go Anda!

### Templat Jarak Jauh

![tangkapan layar remote-linux](/assets/blog-images/remote-linux.webp)

Membuat aplikasi siap dan berjalan dengan cepat selalu menjadi salah satu tujuan utama proyek Wails. Saat pertama kali diluncurkan, kami berupaya mendukung banyak framework modern pada waktu itu: react, vue, dan angular. Dunia pengembangan frontend penuh dengan beragam pandangan, bergerak cepat, dan sulit untuk terus diikuti! Akibatnya, templat dasar kami cukup cepat menjadi usang sehingga menyulitkan pemeliharaan. Hal ini juga berarti kami tidak memiliki templat modern yang keren untuk stack teknologi terbaru dan terbaik.

Dengan v2, saya ingin memberdayakan komunitas dengan memberi Anda kemampuan untuk membuat dan meng-host templat sendiri tanpa harus bergantung pada proyek Wails. Jadi, kini Anda dapat membuat proyek menggunakan templat yang didukung komunitas! Saya berharap hal ini akan menginspirasi developer untuk menciptakan ekosistem templat proyek yang dinamis. Saya sangat antusias menantikan apa yang dapat dibuat oleh komunitas developer kita!

### Kompilasi Silang untuk Windows

Karena Wails v2 untuk Windows dibuat sepenuhnya dengan Go, Anda dapat menargetkan build Windows tanpa Docker.

![tangkapan layar build-cross-windows](/assets/blog-images/linux-build-cross-windows.webp)

### Kesimpulan

Seperti yang saya sampaikan dalam catatan rilis Windows, Wails v2 menjadi fondasi baru bagi proyek ini. Tujuan rilis ini adalah mendapatkan umpan balik mengenai pendekatan baru tersebut dan mengatasi semua bug sebelum rilis penuh. Masukan Anda akan sangat kami hargai! Silakan sampaikan umpan balik di papan diskusi [v2 Beta](https://github.com/wailsapp/wails/discussions/828).

Linux memang **sulit** untuk didukung. Kami memperkirakan versi beta ini akan menunjukkan sejumlah perilaku yang tidak semestinya. Bantu kami membantu Anda dengan mengirimkan laporan bug yang terperinci!

Terakhir, saya ingin menyampaikan terima kasih khusus kepada semua [sponsor proyek](/credits/#sponsors) yang dukungannya mendorong kemajuan proyek ini dalam berbagai cara di balik layar.

Saya tidak sabar melihat apa yang akan dibuat orang-orang dengan Wails pada fase baru proyek yang menarik ini!

Lea.

NB: Rilis v2 kini sudah tidak lama lagi!

NB 2: Jika Anda atau perusahaan Anda merasa Wails bermanfaat, pertimbangkan untuk [mensponsori proyek ini](https://github.com/sponsors/leaanthony). Terima kasih!
