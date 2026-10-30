---
title: "Wails v3.0.0 Dirilis"
description: "Wails v3.0.0 telah hadir: dukungan untuk banyak jendela, baki sistem, menu native, layanan dan binding yang dihasilkan, pemutakhir, dukungan perangkat seluler, serta situs dokumentasi baru."
authors: ["leaanthony"]
tags: ["wails","v3","release"]
date: "2026-10-26"
slug: "blog/wails-v3-released"
image: "/assets/videos/wails-v3-release-clean-poster.webp"
sourcePath: "blog/wails-v3-released.md"
---

@video{src="/assets/videos/wails-v3-release-clean-av1.mp4" poster="/assets/videos/wails-v3-release-clean-poster.webp" title="Film peluncuran Wails v3.0.0 — paling baik ditonton dalam layar penuh" width="1920" height="1080"}

## Akhirnya hadir!

Hari ini saya merilis Wails v3.0.0.

Terus terang, perjalanannya panjang. Tag alfa pertama v3 dirilis pada Januari 2023, dan sejak saat itu begitu banyak orang telah membuat, menguji, merusak, dan memperbaiki berbagai hal di sepanjang perjalanan. Saya sungguh berterima kasih kepada Anda semua. Jika Anda ingin membaca kisah apa adanya tentang perjalanan tersebut, saya menuliskannya dalam [pengumuman beta v3](/blog/wails-v3-beta/).

## Apa *itu* Wails?

Jika Anda baru mengenal Wails: Wails memungkinkan pemrogram Go membuat aplikasi desktop dengan alat frontend web yang sudah mereka kenal, menggunakan WebView native di setiap platform alih-alih browser tersemat. Anda menulis backend dalam Go dan UI dalam HTML, CSS, serta JavaScript (atau framework apa pun yang Anda sukai), lalu Wails mengemas semuanya menjadi satu aplikasi native.

## Mengapa v3?

Wails v2 memudahkan pembuatan aplikasi Go dengan frontend web modern dan telah digunakan dengan baik oleh banyak aplikasi. Namun, runtime-nya yang berbasis konteks dan hanya mendukung satu jendela membuat beberapa pekerjaan desktop sehari-hari lebih sulit daripada seharusnya.

v3 dimulai dengan model yang berbeda. Aplikasi, jendela, layanan, peristiwa, dan kemampuan platform merupakan objek eksplisit. Anda membuat aplikasi, mendaftarkan layanan, membuat jendela, lalu berinteraksi dengan objek yang memiliki perilaku yang Anda perlukan. Model ini lebih mudah diikuti seiring berkembangnya aplikasi, dan hal-hal seperti banyak jendela menjadi bagian normal dari model, bukan solusi sementara.

## Yang tersedia dalam v3.0.0

- **Banyak jendela.** Buat, kelola, dan tutup jendela saat runtime, masing-masing dengan konfigurasi dan siklus hidupnya sendiri.
- **Baki sistem.** Pertahankan keberadaan aplikasi Anda dengan menu, ikon terang dan gelap, ikon templat di macOS, serta opsi untuk menautkan jendela ke ikon baki.
- **Menu dan dialog native.** Menu aplikasi, menu konteks, dan dialog native di setiap platform.
- **Layanan dan binding yang dihasilkan.** Backend Anda tetap berupa kode Go biasa. Binding dihasilkan melalui analisis statis kode sumber, sehingga komentar dan nama parameter yang bermakna diteruskan ke frontend hanya dengan satu perintah `wails3 generate bindings`.
- **Sistem build yang transparan.** Proyek menggunakan build berbasis [Taskfile](https://taskfile.dev) yang dapat Anda periksa, perluas, dan debug, alih-alih satu perintah yang tidak transparan.
- **Lintas platform.** Windows (amd64 dan arm64), macOS (Intel dan Apple Silicon), serta Linux (amd64 dan arm64), dengan GTK4 dan WebKitGTK 6.0 sebagai pilihan default di Linux. Stack GTK3 lama tetap tersedia dengan `-tags gtk3` hingga v3.1.
- **iOS dan Android.** `main.go` dan frontend yang sama dapat berjalan di perangkat seluler. Dukungan perangkat seluler masih ditandai eksperimental dalam dokumentasi, jadi perlakukan dukungan ini sebagai sesuatu untuk dieksplorasi dan diberi masukan, bukan sebagai jaminan stabil.
- **Dibuat dengan Go 1.26.** Modul v3 dibuat dengan Go 1.26.
- **Pemutakhir dalam aplikasi.** Penyedia yang dapat diganti, verifikasi kriptografis, pertukaran atomik, serta UI default yang dapat Anda sesuaikan temanya atau ganti. Lihat [panduan pemutakhir](/guides/updater/).
- **Pengemasan dan penandatanganan.** Panduan untuk [penandatanganan kode](/guides/build/signing/) di setiap platform, [penginstal](/guides/installers/), dan [pengemasan MSIX](/guides/build/msix/).
- **Build server.** Jalankan aplikasi dan layanan yang sama sebagai server HTTP tanpa jendela native.

## Sorotan sejak versi beta

Tulisan tentang versi beta diterbitkan pada 2 Agustus, dan sejak saat itu changelog terus dipenuhi pembaruan. Berikut beberapa hal yang ingin saya soroti:

- **Situs dokumentasi baru.** Dokumentasi kini berjalan di [M-Press](https://m-press.me), alat dokumentasi kami sendiri, dan terjemahannya telah diselesaikan serta dipindahkan dari komunitas. Dokumentasi berbahasa Inggris tetap menjadi sumber acuan utama.
- **Pengalaman yang lebih tangguh di Windows.** Wails kini dapat pulih dari kegagalan proses WebView2 alih-alih meninggalkan jendela kosong, dan sejumlah masalah siklus hidup COM WebView2 juga telah diperbaiki.
- **Pintasan global yang mengikuti tata letak papan ketik Anda.** Di macOS, akselerator huruf kini mengikuti tombol berlabel pada tata letak seperti AZERTY, QWERTZ, dan Dvorak.
- **Penguatan pemutakhir.** Antara lain, perbaikan untuk pertukaran berkas biner antarsistem berkas di Linux dan untuk pembaruan AppImage.
- **Penyiapan dan alat bantu.** Wizard `wails3 setup` telah diselesaikan, dan `wails3 dev` kini menghentikan proses latar belakangnya saat aplikasi Anda ditutup.
- **Liquid Glass native di macOS.** Dukungan untuk efek Liquid Glass di macOS 26 dan versi yang lebih baru, dengan fallback untuk versi sebelumnya.

## Beralih dari v2

v3 adalah versi mayor baru, sehingga migrasinya merupakan porting yang sesungguhnya, bukan sekadar menaikkan versi. Perubahan utamanya mencakup siklus hidup aplikasi dan jendela, layanan yang menggantikan binding berbasis konteks, API aplikasi dan jendela langsung yang menggantikan paket runtime v2, serta binding frontend yang dihasilkan ulang.

[Panduan migrasi dari v2 ke v3](/migration/v2-to-v3/) menjelaskan semuanya dengan contoh sebelum dan sesudah. Lakukan porting dengan cermat, uji selagi Anda mengerjakannya, dan pertahankan v2 hingga aplikasi v3 Anda siap.

**Wails v2 tidak akan dihentikan.** Versi ini akan terus menerima pembaruan kritis. Untuk proyek baru, gunakan v3.0.0. Kedua CLI merupakan berkas biner terpisah (`wails` dan `wails3`) dengan jalur impor modul yang berbeda, sehingga keduanya dapat terpasang di satu mesin.

## Instalasi

Instal CLI v3.0.0 dengan:

```sh
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0
```

Kemudian jalankan wizard penyiapan terpandu, yang akan memeriksa lingkungan Anda dan membantu menyiapkan dependensi yang diperlukan Wails:

```sh
wails3 setup
```

![Wizard penyiapan Wails v3 dalam mode gelap](/assets/screenshots/wails3-setup-wizard-dark.png)

Lalu buat proyek pertama Anda:

```sh
wails3 init
```

- Dokumentasi: <https://v3.wails.io/>
- Panduan migrasi: <https://v3.wails.io/migration/v2-to-v3/>
- Catatan rilis: [Wails v3.0.0 di GitHub](https://github.com/wailsapp/wails/releases/tag/v3.0.0)

## Dukungan dan rilis

Wails telah berkembang jauh melampaui perkiraan saya, dan sulit untuk menangani banyaknya issue serta pull request. Rilis terkadang dilakukan secara ad hoc, dan perbaikan serta kontribusi harus menunggu terlalu lama sebelum sampai kepada pengguna. Kami telah menyiapkan sistem untuk lebih banyak mengotomatiskan pekerjaan dasar dalam triase dan persiapan rilis, sehingga kami dapat merespons masalah dengan lebih cepat dan meluangkan lebih banyak waktu untuk merilis penyempurnaan serta fitur baru.

Setelah v3.0.0, kami akan beralih ke **prarilis nightly** dengan tag bertanggal seperti `v3.1.0-nightly.20261026`, disertai **rilis minor mingguan** seperti `v3.1.0` dan `v3.2.0`. Rilis nightly menyediakan perbaikan terbaru lebih cepat, sedangkan rilis mingguan memberi Anda jadwal pembaruan teruji yang teratur.

Jika menemukan bug yang dapat direproduksi, laporkan bug tersebut dengan menyertakan output `wails3 doctor` dan contoh minimal. Jika ingin mengusulkan kemampuan baru atau perubahan pada perilaku publik, buat PR draf WEP (Wails Enhancement Proposal), bukan issue permintaan fitur.

## Terjemahan

Kepada semua orang yang telah membantu menerjemahkan dokumentasi Wails v2: terima kasih yang sebesar-besarnya dari lubuk hati kami. Anda telah meluangkan waktu agar para pengembang dapat membaca dokumentasi dalam bahasa mereka sendiri, dan hal itu sangat berarti bagi banyak orang. Kami sungguh berterima kasih.

Kami merasa Crowdin sangat sulit dikelola, sehingga kami tidak lagi menggunakannya. Mulai v3, dokumentasi diterjemahkan secara otomatis dan hasil terjemahannya disimpan di repositori bersama halaman berbahasa Inggris.

Jika Anda menginginkan terjemahan dalam bahasa yang belum tersedia, atau menemukan kesalahan dalam terjemahan yang sudah ada, silakan buka issue dan beri tahu kami. Kami sangat senang menerima masukan dari Anda dan berterima kasih atas setiap laporan.

## Penghargaan & Terima Kasih

Wails v3 terwujud berkat orang-orang yang menguji build yang belum selesai, melaporkan bug yang sulit, menerjemahkan dokumentasi, menjawab pertanyaan, dan menyumbangkan kode. Terima kasih. Daftar lengkap kontributor tersedia di [halaman penghargaan](/credits/), dan saya mendorong Anda untuk melihatnya.

Pertama-tama, terima kasih yang sebesar-besarnya kepada [Atterpac](https://github.com/atterpac), maintainer Wails, atas alat bantu pengembangan, integrasi desktop Linux, dan dukungan komunitas yang telah membantu membawa v3 hingga tahap ini.

Terima kasih khusus kepada [Fabio Massaioli (fbbdev)](https://github.com/fbbdev) karena telah membuat penganalisis kode sumber statis dan generator binding yang luar biasa. Keduanya berperan besar dalam membuat v3 berfungsi: layanan Go biasa menjadi API JavaScript dan TypeScript, lengkap dengan tipe dan dokumentasi yang diteruskan ke frontend Anda.

Terima kasih juga kepada [Travis McLane](https://github.com/tmclane) atas kontribusi awalnya yang luas pada platform; [Stffabi](https://github.com/stffabi) atas server aset dan fondasi platform native; [Zach Botterman (popaprozac)](https://github.com/popaprozac) atas notifikasi, integrasi Dock, dan fitur desktop lainnya; [Ian VanSchooten (IanVS)](https://github.com/IanVS) atas event bertipe dan peningkatan runtime; serta [Lyimmi](https://github.com/Lyimmi) atas kontribusinya untuk Linux. Kontribusi ini tetap menjadi bagian dari v3, bersama karya semua orang yang telah membantu membangun dan menyempurnakannya.

Terima kasih yang sungguh tulus kepada para sponsor yang telah mendukung Wails sepanjang transisi panjang ini. Dukungan Anda memberi kami waktu untuk mengerjakan arsitektur, alat bantu, pengujian, dan dokumentasi yang membawa kami hingga tahap ini. Terima kasih khusus juga kepada JetBrains dan ZSA, serta kepada [CodeRabbit](https://www.coderabbit.ai/oss) yang mendukung proyek ini melalui peninjauan kode.

Terakhir, terima kasih khusus kepada Valerio Tesei atas segala dukungannya terhadap proyek ini selama bertahun-tahun.

## Penutup

Banyak pekerjaan sudah menanti untuk kami kerjakan. Ada sejumlah perkembangan menarik di dunia teknologi, dan kami siap mengeksplorasi paradigma baru serta menghadirkan berbagai fitur yang sepenuhnya baru ke Wails. Namun, terlebih dahulu kami perlu menjadikan v3 sebagai fondasi stabil yang dapat Anda andalkan. Memastikan hal itu tercapai lebih penting daripada meluncurkan apa yang berikutnya.

Saya sangat bangga dengan v3. Saya harap Anda merasa sama senangnya membangun aplikasi dengannya seperti kami saat membuatnya. Cobalah, buat sesuatu, lalu sampaikan pendapat Anda kepada kami. Mari menyapa di [Discord](https://discord.gg/JDdSxwjhGf), [GitHub Discussions](https://github.com/wailsapp/wails/discussions), atau [subreddit](https://www.reddit.com/r/wails/).

Selamat menikmati!

&dash;&nbsp;Lea

PS: Jika Anda atau perusahaan Anda merasa Wails bermanfaat, harap pertimbangkan untuk [mensponsori proyek ini](https://github.com/sponsors/leaanthony). Terima kasih!
