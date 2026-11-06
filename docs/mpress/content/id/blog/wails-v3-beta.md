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

Wails memungkinkan pengembang Go membangun aplikasi desktop dengan alat frontend web yang sudah mereka kenal, menggunakan WebView native di setiap platform alih-alih browser tersemat. v3 merupakan kemajuan besar: versi ini menyediakan API yang lebih langsung bagi aplikasi, model build yang lebih jelas, dan fondasi yang lebih baik untuk aplikasi desktop yang selama ini diharapkan dapat didukung oleh Wails.

Ini adalah rilis beta, bukan rilis 3.0 final. API desktop sudah stabil dan sejumlah tim telah menggunakan v3 dalam produksi, tetapi Anda sebaiknya melakukan pengujian menyeluruh sebelum menerapkannya. Selama periode beta, kami bekerja bersama komunitas untuk menemukan masalah kompatibilitas dan alur kerja yang tersisa. Wails v2 tetap menjadi rilis stabil saat ini dan akan terus menerima perbaikan.

## Dokumentasi selama periode beta

Selama periode beta, kami mempertahankan dokumentasi berbahasa Inggris sebagai sumber acuan utama sementara API dan alur kerja menjalani validasi akhir. Pada tahap ini, kami tidak menerima PR terjemahan. Pekerjaan penerjemahan akan dilanjutkan sebelum ketersediaan umum, setelah dokumentasi cukup stabil sehingga penerjemah dapat bekerja tanpa harus berulang kali menyesuaikan perubahan.

## Yang tersedia di v3

- API aplikasi dan jendela yang eksplisit, termasuk dukungan kelas utama untuk banyak jendela
- Layanan Go dengan analisis sumber statis yang menghasilkan binding TypeScript yang lebih kaya, dengan tetap mempertahankan komentar dan nama parameter yang bermakna
- Layanan yang dapat membundel aset dan skrip frontend bersama API backend-nya—fondasi bagi plugin yang dapat diinstal dan lebih kaya
- Sistem build berbasis Taskfile yang transparan serta dapat Anda periksa, perluas, dan debug
- Build server untuk menjalankan aplikasi dan layanan yang sama tanpa jendela desktop native
- Dukungan desktop modern untuk macOS, Windows, dan Linux pada Intel, Apple Silicon, amd64, dan arm64 jika didukung
- Dukungan seluler eksperimental untuk iOS dan Android, tersedia untuk dieksplorasi tetapi tidak tercakup dalam jaminan kompatibilitas beta desktop

## Mengapa v3

Wails v2 memudahkan pembuatan aplikasi Go dengan frontend web modern. Versi ini telah mendukung proyek Wails—dan sangat banyak aplikasi—dengan baik. Namun, runtime satu jendelanya yang digerakkan oleh konteks dan proses build yang dikelola secara ketat membuat beberapa pekerjaan desktop umum menjadi lebih sulit daripada semestinya.

v3 dimulai dengan model yang berbeda. Aplikasi, jendela, layanan, peristiwa, dan kapabilitas platform merupakan objek yang eksplisit. Dengan demikian, framework lebih mudah dipahami seiring berkembangnya aplikasi, dan fitur seperti banyak jendela menjadi bagian normal dari model aplikasi, bukan solusi sementara.

## Yang baru

### API aplikasi yang dirancang untuk perangkat lunak desktop sesungguhnya

v3 menggantikan gaya konfigurasi `wails.Run(...)` pada v2 dengan siklus hidup aplikasi yang eksplisit. Anda membuat aplikasi, mendaftarkan layanan, membuat jendela, dan berinteraksi dengan objek yang memiliki perilaku yang Anda perlukan.

Hal ini menghilangkan banyak penerusan konteks implisit. Operasi jendela menjadi milik jendela; operasi yang mencakup seluruh aplikasi menjadi milik aplikasi. Model ini lebih alami untuk aplikasi berjendela banyak dan lebih sesuai untuk pengujian serta pemeliharaan basis kode yang lebih besar.

Tanggapan dari orang-orang yang menggunakan v3 selama fase alfa sangat positif. Secara khusus, para pengembang menyambut baik model eksplisit tersebut: model ini membuat kode lebih mudah diikuti, memperjelas kepemilikan, dan memberi ruang bagi aplikasi desktop yang kompleks untuk berkembang tanpa harus berhadapan dengan keterbatasan framework.

### Dukungan kelas utama untuk banyak jendela

Banyak jendela merupakan kapabilitas inti v3. Setiap jendela memiliki siklus hidupnya sendiri dan dapat dibuat, dikelola, serta ditutup saat runtime. Hasilnya adalah jalur yang lebih jelas untuk membangun perangkat lunak desktop yang memerlukan editor, inspektor, preferensi, jendela alat, atau beberapa bagian antarmuka pengguna yang independen.

### Layanan dan binding yang dihasilkan

Layanan Go menggantikan model binding lama. Layanan tersebut mempertahankan logika aplikasi sebagai kode Go biasa dan membuat batas dengan frontend menjadi eksplisit. Binding dihasilkan dalam struktur yang mencerminkan aplikasi beserta layanannya, sehingga API yang diekspos ke frontend lebih mudah ditemukan dan digunakan.

v3 menghasilkan binding tersebut dengan analisis sumber statis. Artinya, generator dapat mempertahankan informasi yang dicantumkan pengembang dalam kode mereka—termasuk komentar dan nama parameter yang bermakna—alih-alih menemukan program yang sudah selesai di-build melalui refleksi. Hasilnya adalah API frontend yang lebih kaya dan lebih berguna, serta proses pembuatan binding yang lebih mudah dipahami dan dipelihara.

Layanan juga dapat memasang aset dan skrip frontend berdampingan dengan kode Go-nya. Dengan demikian, satu kapabilitas memiliki satu tempat yang terpadu: API backend-nya, JavaScript atau antarmuka pengguna yang diperlukan, dan titik integrasi untuk aplikasi host. Hal ini membuka jalan bagi plugin Wails yang menyediakan fungsionalitas kaya dan siap digunakan—instal plugin, pasang layanannya, lalu gunakan fiturnya—alih-alih merangkai sendiri sekumpulan binding dan dependensi frontend yang terpisah. Sistem plugin umum bukan bagian dari versi beta ini, tetapi v3 menjadikan arah tersebut praktis dengan cara yang tidak dapat diwujudkan oleh model binding v2.

### Sistem build yang dapat Anda periksa dan sesuaikan

v3 membuat struktur build proyek terlihat. Alih-alih menyembunyikan setiap keputusan build di dalam satu perintah, proyek memiliki tata letak konvensional dan konfigurasi build berbasis Taskfile yang dapat dipahami, diperluas, serta di-debug bersama aplikasi.

### Landasan desktop lintas platform yang lebih kuat

Versi beta ini mendukung Windows pada amd64 dan arm64, macOS pada Intel dan Apple Silicon, serta Linux pada amd64 dan arm64. GTK4 dengan WebKitGTK 6.0 merupakan stack Linux default; GTK3 tetap tersedia sebagai opsi lama di seluruh seri v3.0. Dukungan seluler menjanjikan, tetapi masih bersifat eksperimental dan tidak termasuk dalam jaminan kompatibilitas beta desktop.

Rilis ini juga mencakup pekerjaan yang diperlukan agar pengalaman sehari-hari lebih andal: perilaku platform yang ditingkatkan, model jendela yang lebih mumpuni, diagnostik yang lebih jelas, serta artefak rilis yang dilengkapi checksum dan informasi asal-usul.

## Beralih dari v2

v3 adalah versi mayor baru, dan migrasinya merupakan porting yang sesungguhnya, bukan sekadar perubahan nomor versi. Perubahan konseptual utamanya meliputi siklus hidup aplikasi dan jendela, layanan sebagai pengganti binding yang terikat konteks, API aplikasi dan jendela langsung sebagai pengganti paket runtime v2, serta binding frontend yang dibuat ulang.

Kami telah menerbitkan [panduan migrasi v2 ke v3](/migration/v2-to-v3/) yang menguraikan perubahan tersebut serta menyertakan pemetaan fitur dan daftar periksa pengujian. Panduan manual tersebut merupakan jalur migrasi yang didukung untuk versi beta ini. Jangan berasumsi bahwa setiap proyek v2 dapat dikonversi tanpa peninjauan: uji hasilnya, porting panggilan runtime Anda secara terencana, dan pertahankan v2 hingga aplikasi baru siap.

Kami juga sedang mengevaluasi asisten migrasi eksperimental. Asisten ini bukan bagian dari rilis beta ini, dan kami hanya akan merekomendasikannya setelah divalidasi terhadap proyek v2 dunia nyata yang representatif.

## Catatan terbuka tentang perjalanan ini

Tag alfa v3 pertama diterbitkan pada 18 Januari 2023. Itu merupakan waktu yang lama untuk berada dalam tahap alfa dan patut mendapatkan pengakuan yang lebih jelas daripada sekadar pernyataan samar.

Proyek ini telah mengalami perubahan yang sangat besar selama periode tersebut. Ketika beta Wails v2 pertama untuk Windows dirilis pada September 2021, repositori tersebut memiliki sekitar 4000 bintang. Saat v2 dirilis pada September 2022, jumlahnya sekitar 10300. Kini jumlahnya telah melampaui 35000. Pertumbuhan itu merupakan suatu kehormatan, tetapi juga mengubah seperti apa pengelolaan yang baik: makin banyak pengguna bergantung pada keputusan rilis, makin banyak kontributor membutuhkan jalur yang jelas untuk berpartisipasi, dan makin banyak pekerjaan yang berfokus pada upaya membuat proyek dapat diprediksi, bukan sekadar menambahkan fitur berikutnya.

Saya tidak selalu menyesuaikan proses kami secepat atau sejelas yang dituntut oleh pertumbuhan tersebut. Saya bertanggung jawab atas hal itu. Jawabannya bukanlah membuat janji besar atau mengubah setiap keputusan menjadi serangkaian formalitas, melainkan memperjelas apa yang didukung, apa yang masih eksperimental, bagaimana keputusan dibuat, dan apa yang akan kami lakukan selanjutnya.

Pekerjaan itu dimulai dengan versi beta ini. Kini kami memiliki tonggak pencapaian yang eksplisit untuk versi beta, kandidat rilis, dan GA; komitmen kompatibilitas yang lebih jelas; kebijakan keamanan yang diperbarui; serta proses WEP (Wails Enhancement Proposal) untuk perubahan perilaku publik dan kemampuan baru. Kami akan terus menyempurnakan peta jalan dan meninjau tata kelola proyek seiring pertumbuhan Wails. Tujuannya adalah agar proyek terasa lebih mudah diandalkan dan orang merasa lebih mudah berkontribusi pada proyek ini—bukan agar proyek semakin sulit dikembangkan.

Kami juga meluncurkan kembali [subreddit Wails](https://www.reddit.com/r/wails/) dan akan mengelolanya secara aktif sebagai tempat lain untuk diskusi praktis, pertanyaan, dan umpan balik seiring v3 menuju ketersediaan umum.

## Instal dan coba versi beta

Setelah rilis diterbitkan, instal CLI v3 terbaru dengan:

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
- Catatan rilis: [PLACEHOLDER VERSI](https://github.com/wailsapp/wails/releases)

Jika menemukan bug yang dapat direproduksi, laporkan bug tersebut beserta keluaran `wails3 doctor` dan, jika memungkinkan, contoh minimal. Jika ingin mengusulkan kemampuan baru atau perubahan perilaku publik, buat draf PR WEP, bukan issue permintaan fitur. Kedua jalur tersebut membantu kami memberikan tanggapan yang jelas dan menjaga kemajuan versi beta.

## Terima kasih

Wails v3 terwujud berkat orang-orang yang menguji build yang belum selesai, melaporkan bug yang sulit, menerjemahkan dokumentasi, menjawab pertanyaan, menyumbangkan kode, dan terus mendorong proyek ini agar menjadi lebih baik. Terima kasih.

Dan terima kasih yang sangat istimewa dan tulus kepada para sponsor yang telah menopang Wails selama masa transisi yang panjang ini. Dukungan Anda tidak sekadar menjaga proyek ini tetap berjalan: dukungan tersebut memberi kami kesempatan untuk mencurahkan waktu secara berkelanjutan pada arsitektur, alat bantu, pengujian, dan dokumentasi yang memungkinkan proyek ini bergerak lebih cepat menuju v3. Setiap penguji dan kontributor telah membantu membentuk rilis ini, tetapi para sponsor memungkinkan pekerjaan tersebut mendapatkan perhatian yang semestinya.

Versi beta ini merupakan undangan untuk membantu kami menuntaskan v3 dengan baik. Cobalah, bangun aplikasi dengannya, beri tahu kami di mana masalah terjadi, dan bantu kami menjadikan perjalanan menuju 3.0 singkat dan tetap dilakukan dengan cermat.
