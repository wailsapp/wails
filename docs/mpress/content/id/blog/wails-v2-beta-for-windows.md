---
title: "Wails v2 Beta untuk Windows"
description: "Catatan rilis dan pengumuman untuk Wails"
authors: ["leaanthony"]
tags: ["wails","v2"]
date: "2021-09-27"
slug: "blog/wails-v2-beta-for-windows"
image: "/assets/blog-images/wails.webp"
sourcePath: "blog/wails-v2-beta-for-windows.md"
---

![tangkapan layar Wails](/assets/blog-images/wails.webp)

Ketika saya pertama kali mengumumkan Wails di Reddit, lebih dari 2 tahun lalu dari sebuah kereta di Sydney, saya tidak menyangka pengumuman itu akan mendapat banyak perhatian. Beberapa hari kemudian, seorang vlogger teknologi yang produktif merilis video tutorial, memberikan ulasan positif, dan sejak saat itu minat terhadap proyek ini meningkat pesat.

Jelas bahwa orang-orang antusias menambahkan frontend web ke proyek Go mereka, dan mereka hampir seketika mendorong proyek ini melampaui tahap bukti konsep yang telah saya buat. Saat itu, Wails menggunakan proyek [webview](https://github.com/webview/webview) untuk menangani frontend, dan satu-satunya opsi untuk Windows adalah renderer IE11. Banyak laporan bug berakar pada keterbatasan ini: dukungan JavaScript/CSS yang buruk dan tidak adanya alat pengembang untuk melakukan debugging. Pengalaman pengembangan ini membuat frustrasi, tetapi tidak banyak yang dapat dilakukan untuk memperbaikinya.

Untuk waktu yang lama, saya sangat yakin bahwa Microsoft pada akhirnya harus membereskan masalah browser mereka. Dunia terus bergerak maju, pengembangan frontend berkembang pesat, dan IE sudah tidak memadai. Ketika Microsoft mengumumkan peralihan ke Chromium sebagai dasar arah baru browser mereka, saya tahu bahwa hanya tinggal menunggu waktu sebelum Wails dapat menggunakannya dan membawa pengalaman pengembang Windows ke tingkat berikutnya.

Hari ini, dengan senang hati saya mengumumkan: **Wails v2 Beta untuk Windows**! Ada sangat banyak hal yang perlu dibahas dalam rilis ini, jadi siapkan minuman, duduklah, dan mari kita mulai...

## Tanpa Dependensi CGO!

Tidak, saya tidak bercanda: *Tanpa* *dependensi* *CGO* 🤯! Masalahnya, tidak seperti MacOS dan Linux, Windows tidak dilengkapi compiler bawaan. Selain itu, CGO memerlukan compiler mingw dan tersedia sangat banyak opsi instalasi yang berbeda. Penghapusan persyaratan CGO telah sangat menyederhanakan penyiapan sekaligus membuat debugging jauh lebih mudah. Walaupun saya telah mencurahkan cukup banyak upaya agar ini berfungsi, sebagian besar penghargaan patut diberikan kepada [John Chadwick](https://github.com/jchv), bukan hanya karena memulai beberapa proyek yang memungkinkan hal ini, tetapi juga karena bersedia membiarkan orang lain mengambil dan mengembangkan proyek-proyek tersebut. Penghargaan juga patut diberikan kepada [Tad Vizbaras](https://github.com/tadvi), yang proyek [winc](https://github.com/tadvi/winc)-nya mengawali perjalanan saya di jalur ini.

### Renderer Chromium WebView2

![tangkapan layar alat pengembang](/assets/blog-images/devtools.png)

Akhirnya, pengembang Windows mendapatkan mesin rendering kelas satu untuk aplikasi mereka! Masa ketika Anda harus memelintir kode frontend agar dapat berfungsi di Windows telah berlalu. Selain itu, Anda juga mendapatkan pengalaman alat pengembang kelas satu!

Namun, komponen WebView2 mengharuskan `WebView2Loader.dll` berada di samping berkas biner. Hal ini membuat distribusi sedikit lebih merepotkan daripada yang biasa kami, para gopher, hadapi. Semua solusi dan pustaka yang menggunakan WebView2 (sejauh yang saya ketahui) memiliki dependensi ini.

Namun, dengan sangat antusias saya mengumumkan bahwa aplikasi Wails *tidak memiliki persyaratan tersebut*! Berkat keahlian luar biasa [John Chadwick](https://github.com/jchv), kami dapat membundel DLL ini ke dalam berkas biner dan membuat Windows memuatnya seolah-olah DLL tersebut tersedia di disk.

Bersukacitalah, para gopher! Impian satu berkas biner tetap hidup!

### Fitur Baru

![tangkapan layar menu Wails](/assets/blog-images/wails-menus.webp)

Ada banyak permintaan untuk dukungan menu native. Wails akhirnya menyediakannya untuk Anda. Menu aplikasi kini tersedia dan mendukung sebagian besar fitur menu native. Dukungan ini mencakup item menu standar, kotak centang, grup tombol radio, submenu, dan pemisah.

Pada v1, ada sangat banyak permintaan untuk kemampuan mengendalikan jendela itu sendiri dengan lebih leluasa. Dengan senang hati saya mengumumkan bahwa kini tersedia API runtime baru yang dibuat khusus untuk keperluan ini. API tersebut kaya fitur dan mendukung konfigurasi multi-monitor. API dialog juga telah ditingkatkan: kini Anda dapat menggunakan dialog native modern dengan konfigurasi lengkap untuk memenuhi seluruh kebutuhan dialog Anda.

Kini tersedia opsi untuk menghasilkan konfigurasi IDE bersama proyek Anda. Artinya, jika Anda membuka proyek di IDE yang didukung, IDE tersebut sudah dikonfigurasi untuk membangun dan melakukan debugging aplikasi Anda. Saat ini VSCode didukung, tetapi kami berharap dapat segera mendukung IDE lain seperti Goland.

![tangkapan layar VSCode](/assets/blog-images/vscode.webp)

### Tidak Perlu Membundel Aset

Salah satu kendala terbesar pada v1 adalah keharusan memadatkan seluruh aplikasi Anda menjadi satu berkas JS dan satu berkas CSS. Dengan senang hati saya mengumumkan bahwa pada v2, aset sama sekali tidak perlu dibundel dengan cara apa pun. Ingin memuat gambar lokal? Gunakan tag `<img>` dengan jalur src lokal. Ingin menggunakan font yang keren? Salin font tersebut dan tambahkan jalurnya di CSS Anda.

> Wah, kedengarannya seperti server web...

Ya, cara kerjanya persis seperti server web, tetapi sebenarnya bukan server web.

> Jadi, bagaimana cara menyertakan aset saya?

Cukup teruskan satu `embed.FS` yang berisi seluruh aset Anda ke konfigurasi aplikasi. Aset tersebut bahkan tidak perlu berada di direktori teratas—Wails akan menanganinya untuk Anda.

### Pengalaman Pengembangan Baru

![tangkapan layar browser](/assets/blog-images/browser.webp)

Kini setelah aset tidak perlu dibundel, terbuka pengalaman pengembangan yang sama sekali baru. Perintah `wails dev` yang baru akan membangun dan menjalankan aplikasi Anda, tetapi alih-alih menggunakan aset dalam `embed.FS`, perintah tersebut memuatnya langsung dari disk.

Perintah ini juga menyediakan fitur tambahan berikut:

- Hot reload—Setiap perubahan pada aset frontend akan memicu pemuatan ulang otomatis pada frontend aplikasi
- Pembangunan ulang otomatis—Setiap perubahan pada kode Go Anda akan membangun ulang dan meluncurkan kembali aplikasi

Selain itu, server web akan dimulai pada port 34115. Server ini akan menyajikan aplikasi Anda ke setiap browser yang terhubung. Semua browser web yang terhubung akan merespons peristiwa sistem seperti hot reload ketika aset berubah.

Dalam Go, kita terbiasa menggunakan struct di dalam aplikasi. Sering kali, mengirim struct ke frontend dan menggunakannya sebagai state dalam aplikasi sangatlah berguna. Pada v1, proses ini harus dilakukan secara manual dan agak membebani pengembang. Dengan senang hati saya mengumumkan bahwa pada v2, setiap aplikasi yang dijalankan dalam mode pengembangan akan secara otomatis menghasilkan model TypeScript untuk semua struct yang menjadi parameter input atau output bagi metode yang diikat. Hal ini memungkinkan pertukaran model data yang mulus antara kedua lingkungan tersebut.

Selain itu, modul JS lain dibuat secara dinamis untuk membungkus semua metode yang Anda ikat. Modul ini menyediakan JSDoc untuk metode Anda sehingga IDE dapat menyediakan pelengkapan kode dan petunjuk. Sangat keren ketika model data diimpor secara otomatis setelah Anda menekan tombol tab dalam modul yang dibuat otomatis untuk membungkus kode Go Anda!

### Templat Jarak Jauh

![tangkapan layar remote](/assets/blog-images/remote.webp)

Menyiapkan dan menjalankan aplikasi dengan cepat selalu menjadi tujuan utama proyek Wails. Saat pertama kali diluncurkan, kami berusaha mendukung banyak framework modern pada waktu itu: React, Vue, dan Angular. Dunia pengembangan frontend penuh dengan perbedaan pendapat, bergerak cepat, dan sulit untuk terus diikuti! Akibatnya, template dasar kami cepat menjadi usang dan menimbulkan beban pemeliharaan. Ini juga berarti kami tidak memiliki template modern yang menarik untuk kumpulan teknologi terbaru dan tercanggih.

Dengan v2, saya ingin memberdayakan komunitas dengan memberi Anda kemampuan untuk membuat dan meng-host template sendiri, tanpa bergantung pada proyek Wails. Jadi, sekarang Anda dapat membuat proyek menggunakan template yang didukung komunitas! Saya berharap ini akan menginspirasi para pengembang untuk menciptakan ekosistem template proyek yang dinamis. Saya sangat antusias melihat apa yang dapat diciptakan oleh komunitas pengembang kita!

### Kesimpulan

Wails v2 menjadi fondasi baru bagi proyek ini. Rilis ini bertujuan untuk memperoleh masukan tentang pendekatan baru tersebut dan memperbaiki semua bug sebelum rilis penuh. Masukan Anda akan sangat kami hargai. Silakan sampaikan masukan di papan diskusi [Beta v2](https://github.com/wailsapp/wails/discussions/828).

Ada banyak lika-liku, perubahan arah, dan langkah mundur untuk mencapai titik ini. Hal ini sebagian disebabkan oleh keputusan teknis awal yang perlu diubah, dan sebagian lagi karena beberapa masalah inti yang sebelumnya kami atasi dengan menghabiskan waktu untuk membuat solusi sementara akhirnya diperbaiki di upstream: fitur embed milik Go adalah contoh yang tepat. Untungnya, semuanya terwujud pada saat yang tepat, dan kini kami memiliki solusi terbaik yang bisa kami capai. Saya yakin penantian ini sepadan—hal ini bahkan belum mungkin diwujudkan 2 bulan yang lalu.

Saya juga harus mengucapkan terima kasih yang sebesar-besarnya :pray: kepada orang-orang berikut karena tanpa mereka, rilis ini tidak akan pernah terwujud:

- [Misite Bao](https://github.com/misitebao) — Pekerja keras yang luar biasa dalam menggarap terjemahan bahasa Mandarin dan penemu bug yang hebat.
- [John Chadwick](https://github.com/jchv) — Karya luar biasanya pada [go-webview2](https://github.com/jchv/go-webview2) dan [go-winloader](https://github.com/jchv/go-winloader) memungkinkan terwujudnya Wails versi Windows yang kita miliki saat ini.
- [Tad Vizbaras](https://github.com/tadvi) — Eksperimen dengan proyek [winc](https://github.com/tadvi/winc) miliknya merupakan langkah pertama menuju Wails yang sepenuhnya menggunakan Go.
- [Mat Ryer](https://github.com/matryer) — Dukungan, dorongan, dan masukannya benar-benar membantu mendorong kemajuan proyek ini.

Terakhir, saya ingin menyampaikan terima kasih khusus kepada semua [sponsor proyek](/credits/#sponsors), termasuk [JetBrains](https://www.jetbrains.com?from=Wails), yang dukungannya mendorong proyek ini dalam berbagai cara di balik layar.

Saya tidak sabar melihat apa yang akan dibuat orang-orang dengan Wails pada fase baru proyek yang menarik ini!

Lea.

PS: Pengguna MacOS dan Linux tidak perlu merasa ditinggalkan—proses porting ke fondasi baru ini sedang aktif dikerjakan dan sebagian besar pekerjaan sulitnya sudah selesai. Bersabarlah sedikit lagi!

PPS: Jika Anda atau perusahaan Anda merasa Wails bermanfaat, mohon pertimbangkan untuk [menjadi sponsor proyek](https://github.com/sponsors/leaanthony). Terima kasih!
