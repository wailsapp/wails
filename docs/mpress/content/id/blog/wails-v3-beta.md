---
title: "Wails v3 Beta: fondasi baru untuk aplikasi desktop Go"
description: "Wails v3 Beta memperkenalkan model aplikasi yang lebih langsung, binding yang lebih kaya, dan fondasi yang lebih jelas untuk aplikasi desktop Go."
authors: ["leaanthony"]
tags: ["wails","v3","beta"]
date: "2026-08-02"
slug: "blog/wails-v3-beta"
image: "/assets/screenshots/frameless-v3-native-corners-macos.png"
sourcePath: "blog/wails-v3-beta.md"
---

![Jendela Wails v3 native tanpa bingkai di macOS](/assets/screenshots/frameless-v3-native-corners-macos.png)

Hari ini kami merilis Wails v3 Beta.

Wails memungkinkan pengembang Go membangun aplikasi desktop dengan alat frontend web yang sudah mereka kenal, menggunakan WebView native di setiap platform alih-alih peramban tersemat. v3 merupakan langkah maju yang signifikan: versi ini menyediakan API yang lebih langsung bagi aplikasi, model build yang lebih jelas, dan fondasi yang lebih baik untuk aplikasi desktop yang selama ini diharapkan pengguna dapat didukung oleh Wails.

Ini adalah rilis beta, bukan rilis final 3.0. API desktop sudah stabil dan sejumlah tim telah menggunakan v3 dalam produksi, tetapi Anda harus mengujinya secara menyeluruh sebelum melakukan deployment. Kami menggunakan periode beta untuk menemukan masalah kompatibilitas dan alur kerja terakhir bersama komunitas. Wails v2 tetap menjadi rilis stabil saat ini dan akan terus menerima perbaikan.

## Dokumentasi selama masa beta

Selama masa beta, kami mempertahankan dokumentasi berbahasa Inggris sebagai sumber acuan utama sementara API dan alur kerja menjalani validasi akhir. Pada tahap ini, kami belum menerima PR terjemahan. Penerjemahan akan dilanjutkan sebelum ketersediaan umum, setelah dokumentasi cukup stabil sehingga penerjemah dapat bekerja tanpa harus berulang kali mengerjakan ulang terjemahan.

## Yang tersedia dalam v3

- API aplikasi dan jendela yang eksplisit, termasuk dukungan bawaan untuk banyak jendela
- Layanan Go dengan analisis sumber statis yang menghasilkan binding TypeScript lebih kaya, sekaligus mempertahankan komentar dan nama parameter yang bermakna
- Layanan yang dapat membundel aset dan skrip frontend bersama API backend-nya—fondasi untuk plugin yang dapat diinstal dan lebih kaya
- Sistem build berbasis Taskfile yang transparan serta dapat Anda periksa, perluas, dan debug
- Build server untuk menjalankan aplikasi dan layanan yang sama tanpa jendela desktop native
- Dukungan desktop modern untuk macOS, Windows, dan Linux pada Intel, Apple Silicon, amd64, dan arm64 jika didukung
- Dukungan seluler eksperimental untuk iOS dan Android, tersedia untuk dieksplorasi tetapi tidak tercakup dalam jaminan kompatibilitas beta desktop

## Mengapa v3

Wails v2 memudahkan pembuatan aplikasi Go dengan frontend web modern. Versi tersebut telah melayani proyek ini—dan sangat banyak aplikasi—dengan baik. Namun, runtime berjendela tunggal yang berbasis konteks dan proses build yang dikelola secara ketat membuat beberapa pekerjaan desktop umum lebih sulit daripada semestinya.

v3 berangkat dari model yang berbeda. Aplikasi, jendela, layanan, peristiwa, dan kapabilitas platform merupakan objek eksplisit. Dengan demikian, framework lebih mudah dipahami seiring berkembangnya aplikasi, dan fitur seperti banyak jendela menjadi bagian normal dari model aplikasi, bukan solusi sementara.

## Yang baru

### API aplikasi yang dirancang untuk perangkat lunak desktop sesungguhnya

v3 menggantikan gaya konfigurasi v2 `wails.Run(...)` dengan siklus hidup aplikasi yang eksplisit. Anda membuat aplikasi, mendaftarkan layanan, membuat jendela, dan berinteraksi dengan objek yang memiliki perilaku yang Anda perlukan.

Hal ini meniadakan banyak penerusan konteks implisit. Operasi jendela menjadi milik jendela; operasi yang berlaku di seluruh aplikasi menjadi milik aplikasi. Model ini lebih alami untuk aplikasi berjendela banyak dan lebih sesuai untuk pengujian serta pemeliharaan basis kode yang lebih besar.

Tanggapan dari pengguna v3 selama masa alfa sangat positif. Secara khusus, para pengembang menyambut baik model eksplisit tersebut: kode menjadi lebih mudah diikuti, kepemilikan menjadi lebih jelas, dan aplikasi desktop yang kompleks memiliki ruang untuk berkembang tanpa harus melawan framework.

### Dukungan bawaan untuk banyak jendela

Banyak jendela merupakan kapabilitas inti v3. Setiap jendela memiliki siklus hidupnya sendiri serta dapat dibuat, dikelola, dan ditutup saat runtime. Hasilnya adalah jalur yang lebih jelas untuk membuat perangkat lunak desktop yang memerlukan editor, inspektur, preferensi, jendela alat, atau beberapa bagian UI yang independen.

### Layanan dan binding yang dihasilkan

Layanan Go menggantikan model binding lama. Layanan mempertahankan logika aplikasi sebagai kode Go biasa dan memperjelas batas dengan frontend. Binding dihasilkan dalam struktur yang mencerminkan aplikasi beserta layanannya, sehingga API yang diekspos ke frontend lebih mudah ditemukan dan digunakan.

v3 menghasilkan binding tersebut melalui analisis sumber statis. Artinya, generator dapat mempertahankan informasi yang dimasukkan pengembang ke dalam kode—termasuk komentar dan nama parameter yang bermakna—alih-alih menemukan program yang sudah dibangun melalui refleksi. Hasilnya adalah API frontend yang lebih kaya dan berguna, serta proses pembuatan yang lebih mudah dipahami dan dipelihara.

Layanan juga dapat memasang aset dan skrip frontend bersama kode Go-nya. Dengan demikian, sebuah kapabilitas memiliki satu tempat yang terpadu: API backend-nya, JavaScript atau UI yang diperlukannya, dan titik integrasi untuk aplikasi host. Hal ini membuka jalan bagi plugin Wails yang langsung menyediakan fungsionalitas kaya—instal plugin, pasang layanannya, lalu gunakan fiturnya—alih-alih merangkai sendiri kumpulan binding dan dependensi frontend yang terpisah. Sistem plugin umum bukan bagian dari beta ini, tetapi v3 membuat arah tersebut praktis dengan cara yang tidak dimungkinkan oleh model binding v2.

### Sistem build yang dapat Anda periksa dan sesuaikan

v3 membuat struktur build proyek terlihat jelas. Alih-alih menyembunyikan setiap keputusan build di dalam satu perintah, proyek memiliki tata letak konvensional dan konfigurasi build berbasis Taskfile yang dapat dipahami, diperluas, dan di-debug bersama aplikasi.

### Fondasi desktop lintas platform yang lebih kuat

Beta ini mendukung Windows pada amd64 dan arm64, macOS pada Intel dan Apple Silicon, serta Linux pada amd64 dan arm64. GTK4 dengan WebKitGTK 6.0 merupakan stack Linux default; GTK3 tetap tersedia sebagai opsi lama di seluruh seri v3.0. Dukungan seluler menjanjikan, tetapi masih bersifat eksperimental dan tidak termasuk dalam jaminan kompatibilitas beta desktop.

Rilis ini juga mencakup pekerjaan yang diperlukan agar pengalaman sehari-hari lebih andal: perilaku platform yang ditingkatkan, model jendela yang lebih mumpuni, diagnostik yang lebih jelas, serta artefak rilis yang dilengkapi checksum dan informasi asal-usul.

## Bermigrasi dari v2

v3 adalah versi mayor baru, dan migrasinya merupakan porting nyata, bukan sekadar perubahan nomor versi. Perubahan konseptual utamanya mencakup siklus hidup aplikasi dan jendela, layanan sebagai pengganti binding yang terikat konteks, API aplikasi dan jendela langsung sebagai pengganti paket runtime v2, serta pembuatan ulang binding frontend.

Kami telah menerbitkan panduan migrasi [v2-ke-v3](/migration/v2-to-v3/) yang menguraikan perubahan tersebut serta menyertakan pemetaan fitur dan daftar periksa pengujian. Panduan manual tersebut merupakan jalur migrasi yang didukung untuk beta ini. Jangan berharap setiap proyek v2 dapat dikonversi tanpa peninjauan: uji hasilnya, porting panggilan runtime Anda secara cermat, dan pertahankan v2 hingga aplikasi baru siap.

Kami juga sedang mengevaluasi asisten migrasi eksperimental. Alat tersebut bukan bagian dari rilis beta ini, dan kami hanya akan merekomendasikannya setelah divalidasi terhadap proyek v2 dunia nyata yang representatif.

## Catatan terbuka tentang perjalanan ini

Tag alfa pertama v3 diterbitkan pada 18 Januari 2023. Itu merupakan waktu yang lama untuk berada dalam tahap alfa dan layak mendapatkan pengakuan yang lebih dari sekadar pernyataan samar.

Proyek ini telah berubah sangat besar selama periode tersebut. Ketika beta pertama Wails v2 dirilis untuk Windows pada September 2021, repositorinya memiliki sekitar 4000 bintang. Pada rilis v2 pada September 2022, jumlahnya mencapai sekitar 10300. Kini jumlahnya melebihi 35000. Pertumbuhan tersebut merupakan suatu kehormatan, tetapi juga mengubah bentuk pengelolaan yang baik: semakin banyak pengguna bergantung pada keputusan rilis, semakin banyak kontributor membutuhkan jalur partisipasi yang jelas, dan semakin banyak pekerjaan berfokus pada menjadikan proyek dapat diprediksi, bukan sekadar menambahkan fitur berikutnya.

Saya tidak selalu menyesuaikan proses kami secepat atau sejelas yang dituntut oleh pertumbuhan tersebut. Saya bertanggung jawab atas hal itu. Jawabannya bukanlah membuat janji-janji besar atau mengubah setiap keputusan menjadi seremoni, melainkan memperjelas apa yang didukung, apa yang bersifat eksperimental, bagaimana keputusan dibuat, dan apa yang akan kami lakukan selanjutnya.

Upaya itu dimulai dengan versi beta ini. Kini kami memiliki tonggak beta, kandidat rilis, dan GA yang eksplisit; komitmen kompatibilitas yang lebih jelas; kebijakan keamanan yang diperbarui; serta proses WEP (Wails Enhancement Proposal) untuk perubahan pada perilaku publik dan kemampuan baru. Kami akan terus menyempurnakan peta jalan dan meninjau tata kelola proyek seiring pertumbuhan Wails. Tujuannya adalah proyek yang terasa lebih mudah diandalkan dan lebih mudah dikontribusikan—bukan proyek yang lebih sulit dikembangkan.

Kami juga meluncurkan kembali [subreddit Wails](https://www.reddit.com/r/wails/) dan akan memeliharanya secara aktif sebagai tempat lain untuk diskusi praktis, pertanyaan, dan umpan balik selama v3 bergerak menuju ketersediaan umum.

## Instal dan coba versi beta

Setelah rilis dipublikasikan, instal CLI v3 terbaru dengan:

```sh
go install github.com/wailsapp/wails/v3/cmd/wails3@latest
```

Sebelum membuat proyek, jalankan wizard penyiapan terpandu. Wizard ini memeriksa lingkungan pengembangan lokal Anda dan membantu mengonfigurasi dependensi yang diperlukan Wails:

```sh
wails3 setup
```

![Wizard penyiapan Wails v3 dalam mode gelap](/assets/screenshots/wails3-setup-wizard-dark.png)

### Buat proyek

Setelah semuanya siap, buat proyek dengan:

```sh
wails3 init
```

- Dokumentasi: <https://v3.wails.io/>
- Panduan migrasi: <https://v3.wails.io/migration/v2-to-v3/>
- Reddit: <https://www.reddit.com/r/wails/>
- Catatan rilis: [VERSION PLACEHOLDER](https://github.com/wailsapp/wails/releases)

Jika Anda menemukan bug yang dapat direproduksi, harap laporkan bug tersebut dengan menyertakan keluaran `wails3 doctor` dan, jika memungkinkan, contoh minimal. Jika Anda ingin mengusulkan kemampuan baru atau perubahan pada perilaku publik, mulailah dengan draf PR WEP, bukan issue permintaan fitur. Kedua jalur tersebut membantu kami memberikan tanggapan yang jelas dan menjaga kemajuan versi beta.

## Terima kasih

Wails v3 hadir berkat orang-orang yang menguji build yang belum lengkap, melaporkan bug yang sulit, menerjemahkan dokumentasi, menjawab pertanyaan, menyumbangkan kode, dan terus mendorong proyek ini agar menjadi lebih baik. Terima kasih.

Dan terima kasih yang sangat istimewa dan tulus kepada para sponsor yang telah menopang Wails selama transisi panjang ini. Dukungan Anda bukan sekadar menjaga proyek ini tetap berjalan: dukungan tersebut memberi kami kesempatan untuk mencurahkan waktu secara berkelanjutan pada arsitektur, alat bantu, pengujian, dan dokumentasi yang memungkinkan proyek ini melaju menuju v3. Setiap penguji dan kontributor telah membantu membentuk rilis ini, tetapi para sponsorlah yang memungkinkan pekerjaan tersebut memperoleh perhatian yang semestinya.

Versi beta ini merupakan undangan untuk membantu kami menuntaskan v3 dengan baik. Cobalah, bangun aplikasi dengannya, beri tahu kami bagian yang bermasalah, dan bantu kami menjadikan perjalanan menuju 3.0 singkat sekaligus cermat.
