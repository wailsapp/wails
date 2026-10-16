---
title: "Wails v2 Beta untuk MacOS"
description: "Catatan rilis dan pengumuman untuk Wails"
authors: ["leaanthony"]
tags: ["wails","v2"]
date: "2021-11-08"
slug: "blog/wails-v2-beta-for-mac"
image: "/assets/blog-images/wails-mac.webp"
sourcePath: "blog/wails-v2-beta-for-mac.md"
---

![tangkapan layar wails-mac](/assets/blog-images/wails-mac.webp)

Hari ini menandai rilis beta pertama Wails v2 untuk Mac! Perlu waktu cukup lama untuk mencapai tahap ini, dan saya berharap rilis hari ini dapat memberi Anda sesuatu yang cukup bermanfaat. Ada cukup banyak liku-liku untuk sampai ke tahap ini dan saya berharap, dengan bantuan Anda, dapat membereskan berbagai masalah kecil serta menyempurnakan port Mac untuk rilis final v2.

Maksudnya, ini belum siap digunakan dalam produksi? Untuk kasus penggunaan Anda, mungkin saja sudah siap, tetapi masih ada sejumlah masalah yang diketahui. Jadi, terus pantau [papan proyek ini](https://github.com/wailsapp/wails/projects/7). Jika Anda ingin berkontribusi, kami akan dengan senang hati menyambut Anda!

Jadi, apa yang baru di Wails v2 untuk Mac dibandingkan v1? Petunjuk: Cukup mirip dengan versi Beta Windows :wink:

## Fitur Baru

![tangkapan layar wails-menus-mac](/assets/blog-images/wails-menus-mac.webp)

Ada banyak permintaan untuk dukungan menu native. Wails akhirnya menyediakannya untuk Anda. Menu aplikasi kini tersedia dan mendukung sebagian besar fitur menu native. Ini mencakup item menu standar, kotak centang, grup tombol radio, submenu, dan pemisah.

Di v1, ada sangat banyak permintaan agar pengguna dapat memiliki kontrol yang lebih besar atas jendela itu sendiri. Dengan senang hati saya mengumumkan bahwa kini tersedia API runtime baru khusus untuk keperluan tersebut. API ini kaya fitur dan mendukung konfigurasi multi-monitor. Tersedia pula API dialog yang telah ditingkatkan: Kini Anda dapat menggunakan dialog native modern dengan konfigurasi lengkap untuk memenuhi semua kebutuhan dialog Anda.

### Opsi Khusus Mac

Selain opsi aplikasi biasa, Wails v2 untuk Mac juga menghadirkan beberapa fitur tambahan khusus Mac:

- Buat jendela Anda tampil keren dan semitransparan, seperti semua aplikasi Swift yang cantik!
- Bilah judul dengan banyak opsi penyesuaian
- Kami mendukung opsi NSAppearance untuk aplikasi
- Konfigurasi sederhana untuk membuat menu "Tentang" secara otomatis

### Tidak Perlu Membundel Aset

Salah satu kendala besar di v1 adalah keharusan memadatkan seluruh aplikasi Anda menjadi satu file JS dan satu file CSS. Dengan senang hati saya mengumumkan bahwa di v2, aset sama sekali tidak perlu dibundel dalam bentuk apa pun. Ingin memuat gambar lokal? Gunakan tag `<img>` dengan jalur src lokal. Ingin menggunakan font yang keren? Salin font tersebut dan tambahkan jalurnya ke CSS Anda.

> Wah, kedengarannya seperti server web...

Ya, cara kerjanya persis seperti server web, tetapi sebenarnya bukan server web.

> Jadi, bagaimana cara menyertakan aset saya?

Anda cukup meneruskan satu `embed.FS` yang berisi semua aset Anda ke konfigurasi aplikasi. Aset tersebut bahkan tidak harus berada di direktori teratas—Wails akan menanganinya untuk Anda.

### Pengalaman Pengembangan Baru

Karena aset kini tidak perlu dibundel, hal ini memungkinkan pengalaman pengembangan yang benar-benar baru. Perintah `wails dev` yang baru akan membangun dan menjalankan aplikasi Anda, tetapi alih-alih menggunakan aset di dalam `embed.FS`, perintah tersebut memuatnya langsung dari disk.

Perintah tersebut juga menyediakan fitur tambahan berikut:

- Hot reload—Setiap perubahan pada aset frontend akan memicu pemuatan ulang otomatis pada frontend aplikasi
- Pembangunan ulang otomatis—Setiap perubahan pada kode Go Anda akan membangun ulang dan meluncurkan kembali aplikasi

Selain itu, server web akan dimulai pada port 34115. Server ini akan menyajikan aplikasi Anda ke setiap browser yang terhubung. Semua browser web yang terhubung akan merespons peristiwa sistem, seperti hot reload ketika aset berubah.

Dalam Go, kita terbiasa menggunakan struct di aplikasi. Sering kali, mengirim struct ke frontend dan menggunakannya sebagai state dalam aplikasi sangatlah berguna. Di v1, proses ini harus dilakukan secara manual dan cukup membebani pengembang. Dengan senang hati saya mengumumkan bahwa di v2, setiap aplikasi yang dijalankan dalam mode pengembangan akan secara otomatis menghasilkan model TypeScript untuk semua struct yang menjadi parameter input atau output metode yang diikat. Hal ini memungkinkan pertukaran model data yang lancar antara kedua lingkungan tersebut.

Selain itu, modul JS lain dibuat secara dinamis untuk membungkus semua metode yang Anda ikat. Modul ini menyediakan JSDoc untuk metode Anda sehingga IDE dapat menyediakan pelengkapan kode dan petunjuk. Rasanya sangat keren ketika model data diimpor secara otomatis setelah Anda menekan tombol tab di modul yang dibuat otomatis untuk membungkus kode Go Anda!

### Templat Jarak Jauh

![tangkapan layar remote-mac](/assets/blog-images/remote-mac.webp)

Membuat aplikasi dapat segera berjalan selalu menjadi salah satu tujuan utama proyek Wails. Saat pertama kali diluncurkan, kami mencoba mendukung banyak framework modern pada masa itu: react, vue, dan angular. Dunia pengembangan frontend memiliki banyak preferensi yang kuat, bergerak cepat, dan sulit untuk terus diikuti! Akibatnya, templat dasar kami cepat usang dan menimbulkan masalah pemeliharaan. Ini juga berarti kami tidak memiliki templat modern yang keren untuk tumpukan teknologi terbaru dan terbaik.

Dengan v2, saya ingin memberdayakan komunitas dengan memberi Anda kemampuan untuk membuat dan menghosting templat sendiri, alih-alih bergantung pada proyek Wails. Jadi, sekarang Anda dapat membuat proyek menggunakan templat yang didukung komunitas! Saya berharap hal ini akan menginspirasi para pengembang untuk menciptakan ekosistem templat proyek yang dinamis. Saya benar-benar antusias melihat apa yang dapat diciptakan komunitas pengembang kita!

### Dukungan Native M1

Berkat dukungan luar biasa dari [Mat Ryer](https://github.com/matryer/), proyek Wails kini mendukung build native M1:

![tangkapan layar build-darwin-arm](/assets/blog-images/build-darwin-arm.webp)

Anda juga dapat menetapkan `darwin/amd64` sebagai target:

![tangkapan layar build-darwin-amd](/assets/blog-images/build-darwin-amd.webp)

Oh, saya hampir lupa.... Anda juga dapat membuat build dengan target `darwin/universal`.... :wink:

![tangkapan layar build-darwin-universal](/assets/blog-images/build-darwin-universal.webp)

### Kompilasi Silang ke Windows

Karena Wails v2 untuk Windows sepenuhnya dibuat dengan Go, Anda dapat menargetkan build Windows tanpa docker.

![tangkapan layar build-cross-windows](/assets/blog-images/build-cross-windows.webp)  
bu

### Renderer WKWebView

V1 mengandalkan komponen WebView (yang kini tidak lagi dianjurkan penggunaannya). V2 menggunakan komponen WKWebKit terbaru, sehingga Anda dapat mengharapkan teknologi terkini dan terbaik dari Apple.

### Kesimpulan

Seperti yang telah saya sampaikan dalam catatan rilis Windows, Wails v2 menjadi fondasi baru bagi proyek ini. Tujuan rilis ini adalah memperoleh masukan tentang pendekatan baru tersebut dan memperbaiki semua bug sebelum rilis penuh. Masukan Anda akan sangat kami hargai! Silakan sampaikan masukan apa pun melalui papan diskusi [v2 Beta](https://github.com/wailsapp/wails/discussions/828).

Terakhir, saya ingin menyampaikan terima kasih secara khusus kepada semua [sponsor proyek](/credits/#sponsors), termasuk [JetBrains](https://www.jetbrains.com?from=Wails), yang dukungannya turut menggerakkan proyek ini dalam berbagai hal di balik layar.

Saya menantikan karya yang akan dibuat orang-orang dengan Wails dalam fase baru proyek yang menarik ini!

Lea.

PS: Para pengguna Linux, giliran Anda berikutnya!

PPS: Jika Anda atau perusahaan Anda merasa Wails bermanfaat, harap pertimbangkan untuk [mensponsori proyek ini](https://github.com/sponsors/leaanthony). Terima kasih!
