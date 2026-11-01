---
title: "Wake"
description: "Buat proyek Wake pertama Anda, ikuti wizard pembuatan proyek, build dan jalankan proyek tersebut, serta migrasikan proyek Taskfile yang sudah ada."
slug: "experimental/wake"
sourcePath: "experimental/wake.md"
---

Wake adalah sistem build eksperimental Wails. Satu `wails.yaml` mendeskripsikan aplikasi, frontend, platform target, paket, dan pengaturan penandatanganan Anda. Wake menggunakan konfigurasi tersebut untuk merencanakan dan menjalankan build, menyimpan pekerjaan yang tidak berubah dalam cache, serta melaporkan file yang dihasilkan.

Dalam tutorial ini, Anda akan membuat aplikasi desktop kecil, mengonfigurasinya di browser, mem-build-nya untuk komputer Anda sendiri, dan menjalankannya dalam mode pengembangan. Mulailah dengan satu target lokal dan build tanpa tanda tangan. Tambahkan kompilasi silang dan distribusi setelah build pertama tersebut berhasil.

Sudah memiliki aplikasi? Buka [Migrasikan proyek yang sudah ada](#migrasikan-proyek-yang-sudah-ada). Untuk masing-masing bidang dan alur kerja tingkat lanjut, gunakan [panduan build YAML](/experimental/yaml-builds/) dan [referensi manifes](/reference/wails-yaml-fields/).

## Hal yang dapat diharapkan dari eksperimen ini

@note{type="caution" title="Gunakan proyek pengujian atau branch terpisah"}
Wake bersifat opsional, dan konfigurasi, perintah, serta perilakunya dapat berubah di antara build pratinjau. Wake belum memiliki jaminan kompatibilitas seperti sistem build yang stabil. Simpan revisi proyek yang diketahui berfungsi dan catat versi CLI serta framework yang Anda uji. Uji paket pada sistem tujuannya sebelum mendistribusikannya.
@end

Wake dan backend **Pure Go** merupakan pilihan yang terpisah. Mengaktifkan Wake tidak otomatis mengaktifkan Pure Go. Pelajaran ini menggunakan backend standar; karena itu, build GUI macOS dan Linux memerlukan alat pengembangan native yang biasanya digunakan.

Wizard memeriksa persyaratan dan menyimpan konfigurasi. Menyimpan bukanlah proses build dan tidak membuktikan bahwa setiap target yang dipilih sudah siap. Alat atau kredensial yang belum tersedia dapat ditambahkan nanti, misalnya pada host CI. Kombinasi target atau paket yang tidak valid harus diperbaiki sebelum penyimpanan.

Migrasi bersifat **satu arah**. Wake tidak mengeksekusi Taskfile, memperluas templat Task, atau menerjemahkan YAML kembali menjadi Taskfile. Proyek yang belum dimigrasikan tetap menggunakan runtime Task lama yang terpisah. `wails.yaml` yang aktif tidak pernah diam-diam beralih kembali ke runtime tersebut.

## Sebelum memulai

Anda memerlukan Go, Node.js beserta npm untuk frontend awal, dan [dependensi platform](/getting-started/installation/#platform-specific-dependencies) untuk komputer Anda. Mulailah dari direktori tempat Anda ingin folder proyek baru dibuat. Anda juga memerlukan browser untuk menggunakan wizard.

Gunakan revisi CLI dan framework yang menyertakan pratinjau ini. Menginstal CLI `v3.0.0` stabil tidak menambahkan layar pembuatan proyek Wake yang dijelaskan di sini. Perintah berikut menetapkan pratinjau dari [PR #6531](https://github.com/wailsapp/wails/pull/6531):

```sh
go install github.com/wailsapp/wails/v3/cmd/wails3@b1daf95f3
wails3 version
```

Pastikan direktori executable Go berada di `PATH`; [panduan instalasi](/getting-started/installation/) menjelaskan hal ini. Nanti dalam tutorial, Anda akan menetapkan dependensi Go aplikasi ke revisi yang sama. Saat menguji pratinjau yang lebih baru, perbarui kedua argumen revisi secara bersamaan.

### Aktifkan Wake di terminal Anda

Di Bash atau Zsh:

```bash
export WAILS_EXP_USE_WAKE=1
```

Di PowerShell:

```powershell
$env:WAILS_EXP_USE_WAKE = '1'
```

Biarkan terminal ini tetap terbuka selama mengikuti tutorial. Tetapkan kembali variabel tersebut di setiap terminal atau job CI baru yang menjalankan perintah Wake. **Keberadaan** variabel tersebut mengaktifkan Wake: menetapkannya ke `false` atau string kosong tidak menonaktifkannya. Nama lama `WAILS_USE_WAKE` diabaikan.

## Buat aplikasi

Jalankan wizard pembuatan proyek:

```sh
wails3 init -ui
```

Browser akan membuka wizard proyek. Saat Wake diaktifkan, wizard mengumpulkan pilihan sumber dan build dalam satu alur. Tidak ada yang dibuat hingga tombol terakhir **Buat Proyek** diklik. `wails3 setup` adalah wizard terpisah untuk preferensi global, alat, dan kredensial penandatanganan; wizard tersebut tidak mengonfigurasi profil Wake proyek ini.

| Layar | Pilihan untuk tutorial ini | Yang terjadi berikutnya |
| --- | --- | --- |
| **Proyek — Buat proyek Wails** | Masukkan `wake-demo`, lalu pilih **Lanjutkan**. | Proyek akan dibuat dalam folder `wake-demo` baru di bawah direktori yang ditampilkan. |
| **Framework — Pilih framework** | Pilih **Vanilla**, lalu **Lanjutkan**. | Anda mendapatkan frontend kecil tanpa perlu mempelajari langkah khusus framework. |
| **Bahasa — Preferensi Bahasa** | Pilih **JavaScript**, lalu **Lanjutkan**. | Layar **Binding** yang hanya tersedia untuk TypeScript akan dilewati. |
| **Detail — Detail proyek** | Gunakan nama produk `Wake Demo`, pengenal bundle `com.example.wakedemo`, dan versi `0.1.0`. Klik **Lanjutkan**. | Metadata Anda tetap tersimpan dalam wizard saat Anda memilih pengaturan build. |

Jangan gunakan **Gunakan default global** untuk pelajaran ini karena opsi tersebut melewati pilihan framework dan bahasa. Jika Anda memilih TypeScript dalam proyek lain, layar **Binding** memungkinkan Anda memilih interface atau class untuk binding data Go yang dihasilkan.

## Pilih pengaturan build sebelum membuat proyek

Lanjutkan di jendela browser yang sama. Tangkapan layar di bawah menunjukkan panduan ini di Linux. Target host dan alat yang diperlukan akan berbeda di Windows atau macOS.

### Aplikasi — Apa yang sedang Anda buat?

Pertahankan pilihan **Aplikasi desktop** dan pilih **Backend standar**. Klik **Lanjutkan**.

![Layar Aplikasi dalam wizard pembuatan proyek dengan Aplikasi desktop dan Backend standar dipilih.](/assets/screenshots/wake/application.png)

Pelajaran ini membuat GUI dengan frontend web. **Pure Go (eksperimental)** merupakan opsi terpisah yang harus diaktifkan secara eksplisit. Opsi ini menghapus persyaratan compiler C untuk kompilasi desktop yang didukung, asalkan dependensi Anda juga dapat di-build tanpa CGo. Opsi ini tidak menghapus kebutuhan akan pustaka runtime target atau semua alat pengemasan dan penandatanganan.

Backend merupakan **pengaturan untuk seluruh proyek**. Profil lain yang mewarisi pengaturan tersebut juga menggunakan backend yang dipilih.

### Target — Di mana aplikasi harus dijalankan?

Pertahankan hanya target untuk komputer yang sedang Anda gunakan, lalu klik **Lanjutkan**. Pada awalnya, wizard memilih sistem operasi dan arsitektur host.

![Layar Target dengan hanya linux/amd64 dipilih untuk host tutorial.](/assets/screenshots/wake/targets.png)

Nama target menggabungkan sistem operasi dan arsitektur CPU: misalnya, `windows/amd64`, `darwin/arm64` untuk Mac Apple Silicon, atau `linux/amd64`. `darwin` berarti macOS. Memilih target lain meminta executable untuk mesin tersebut; tindakan ini tidak membuat executable itu dapat dijalankan di komputer Anda saat ini.

Tunda arsitektur tambahan, build universal macOS, dan target seluler eksperimental untuk nanti. Anda dapat menambahkannya ke `wails.yaml` setelah build lokal pertama.

### Paket — Pilih output rilis

Biarkan kotak format paket tidak dicentang dan nonaktifkan **Tandatangani output rilis**. Klik **Lanjutkan**.

![Layar Paket tanpa format paket Linux tambahan yang dipilih dan dengan penandatanganan dinonaktifkan.](/assets/screenshots/wake/packages.png)

Tanpa memilih format tambahan, Wake menghasilkan keluaran default target: berkas executable desktop, atau bundel aplikasi macOS untuk GUI macOS. Ini cukup untuk memeriksa aplikasi secara lokal. Pemilihan format akan menambahkan keluaran distribusi seperti DMG, MSIX, NSIS, DEB, atau AppImage jika didukung.

Penandatanganan bersifat opsional untuk build lokal pertama ini. Build tanpa tanda tangan yang berhasil tidak berarti aplikasi siap didistribusikan kepada publik. Untuk proyek bertanda tangan, konfigurasikan kredensial yang dapat digunakan kembali dengan `wails3 setup` sebelum pembuatan, atau konfigurasikan penandatanganan di `wails.yaml` setelahnya. Wizard pembuatan menggunakan kredensial yang telah Anda konfigurasi; wizard tersebut tidak mengubah preferensi penandatanganan global Anda.

Wake dapat menggunakan berkas sertifikat yang telah dikonfigurasi di luar proyek tanpa menyalin kunci privat atau kata sandi ke `wails.yaml`. Pengaturan penandatanganan proyek yang eksplisit diprioritaskan. Ikuti [panduan penyiapan penandatanganan](/getting-started/setup/) dan [panduan pemaketan dan penandatanganan](/experimental/yaml-builds/#package-and-sign-for-distribution) sebelum menyiapkan rilis publik. Penandatanganan seluler dikonfigurasi dalam YAML.

### Alat — Alat untuk build Anda

Baca persyaratan untuk pilihan Anda. Daftar tersebut berasal dari rencana build Wake; daftar itu bukan permintaan untuk menginstal semua alat di setiap platform.

![Layar Alat yang menampilkan persyaratan yang telah ditentukan untuk build sumber Linux lokal.](/assets/screenshots/wake/tools.png)

Untuk pelajaran ini, Anda memerlukan Go dan alat frontend, serta alat pengembangan native yang diperlukan oleh backend GUI standar host Anda. Docker bukan persyaratan umum. Build silang dengan CGo mungkin memerlukan compiler dan pustaka target atau toolchain kontainer yang didukung.

Selesaikan baris yang ditandai **Belum siap**, **Perlu dikonfigurasi**, atau **Perlu diperbarui** sebelum mencoba build lokal. Gunakan petunjuk yang ditampilkan jika tersedia, atau [panduan instalasi platform](/getting-started/installation/), lalu klik **Periksa ulang persyaratan**. Jika instalasi mengubah `PATH` terminal, mulai ulang `wails3 init -ui` dari terminal yang telah diperbarui dan ulangi pilihan Anda.

Wizard dapat membuat proyek yang valid meskipun komputer ini tidak memiliki alat atau kredensial penandatanganan. Hal ini berguna untuk konfigurasi CI, tetapi untuk tutorial ini Anda memerlukan build yang dapat berjalan di komputer ini. Konfigurasi yang tidak valid harus diperbaiki sebelum pembuatan.

### Buat proyek

Periksa pilihan Anda: aplikasi desktop, backend standar, target host Anda, dan penandatanganan dinonaktifkan. Di **Alat**, klik **Buat Proyek** lalu tunggu hingga muncul **Pilihan proyek diterima**. Kembali ke terminal dan tunggu hingga CLI selesai membuat proyek; konfirmasi di browser hanya menerima pilihan Anda.

![Layar Pilihan proyek diterima yang menampilkan jalur aplikasi yang direncanakan dan perintah yang digunakan setelah pembuatan berhasil.](/assets/screenshots/wake/save.png)

Wizard membuat sumber dan menyimpan pilihan target, paket, serta penandatanganan di bawah `profiles.setup` dalam `wails.yaml`. Wizard juga menetapkan jenis aplikasi, pengaktifan frontend, dan backend untuk seluruh proyek. Anda dapat menutup browser dan kembali ke terminal. Tunggu hingga pembuatan proyek selesai, lalu jalankan:

```sh
cd wake-demo
go get github.com/wailsapp/wails/v3@b1daf95f3
```

Sekarang Anda seharusnya memiliki `main.go`, sebuah modul Go, direktori `frontend`, dan `wails.yaml` native. Templat bawaan menggunakan Wake secara langsung; Anda tidak perlu memigrasikan proyek baru ini. Templat komunitas dapat memiliki tugas lama kustom dan mungkin memerlukan peninjauan migrasi yang dijelaskan di bawah.

Layar penyelesaian menampilkan `wails3 build setup` dan `wails3 dev`. Nama profil penting: `wails3 build setup` menggunakan pilihan yang baru saja Anda simpan. `wails3 build` biasa tidak secara otomatis memilih profil `setup`. Nama profil tersebut tidak berarti Anda harus menjalankan `wails3 setup`.

## Build dan jalankan aplikasi pertama Anda

Jalankan perintah yang ditampilkan oleh wizard:

```sh
wails3 build setup
```

Build pertama menginstal dependensi frontend, menghasilkan binding, dan mengompilasi aplikasi. Build berikutnya dapat menggunakan kembali pekerjaan yang tidak berubah. Tunggu hingga **build berhasil** dan baca jalur keluaran yang dicetak oleh CLI. Direktori keluaran default adalah `bin`; GUI macOS dirakit menjadi bundel `.app`.

Mulai mode pengembangan di komputer Anda sendiri:

```sh
wails3 dev
```

Anda seharusnya melihat jendela aplikasi awal. Masukkan nama Anda dan coba salamnya. Ubah placeholder input nama di `frontend/index.html` dan simpan berkas tersebut; frontend seharusnya dimuat ulang. Edit salam di `greetservice.go` untuk mencoba build ulang Go. Tekan **Ctrl+C** di terminal setelah selesai.

Pengembangan adalah alur kerja di host lokal. Profil rilis dapat memuat beberapa target build silang tanpa mengubah mesin mana yang dapat menjalankan jendela pengembangan. Untuk mem-build dan meluncurkan aplikasi produksi secara lokal, gunakan `wails3 run`; lihat [referensi menjalankan aplikasi](/reference/wails-yaml-run/).

Sekarang Anda memiliki proyek sumber, profil build yang telah ditinjau, dan siklus pengembangan lokal yang berfungsi. Commit sumber dan `wails.yaml` ke kontrol versi. Berkas kerja `.wails/` yang dihasilkan bukan sumber konfigurasi Anda.

## Yang dapat Anda lakukan selanjutnya

| Tujuan | Dukungan Wake | Yang perlu diperiksa terlebih dahulu |
| --- | --- | --- |
| Build untuk desktop lain | Target Windows, macOS, dan Linux, termasuk kedua arsitektur macOS dan keluaran universal. | Pure Go dan CGo memiliki persyaratan compiler yang berbeda. Uji pada OS dan arsitektur tujuan. |
| Buat penginstal atau paket | MSIX, NSIS, aplikasi/DMG macOS, DEB, RPM, Arch, Void, dan AppImage. | NSIS memerlukan `makensis`. AppImage memerlukan runtime target; pemaketan silang GUI juga memerlukan pohon runtime Linux yang disediakan. |
| Tandatangani rilis | Penandatanganan portabel dengan kunci yang diekspor untuk Windows, macOS, dan DEB/RPM/Arch; penyimpanan identitas native tetap tersedia. | Setiap identitas, format, dan host memiliki batasan. AppImage dan Void tidak memiliki alur penandatanganan bawaan. |
| Notarisasikan aplikasi macOS | Pengiriman menggunakan kunci API tim dan penempelan tiket, atau profil keychain macOS. | Tanda tangan Developer ID dan stempel waktu yang sesuai tetap diperlukan. Tanda tangan pengujian yang dihasilkan Go dan lolos pemeriksaan native bukan bukti bahwa Apple akan menerima rilis Anda. |
| Build beberapa program | Program utama GUI, CLI, atau server beserta program pendamping bernama. | Deklarasikan paket dan jenis program pendamping dalam YAML. Wizard tidak menghasilkan kode sumbernya. |
| Paketkan berkas yang dibuat di tempat lain | `wails3 build setup --skip-build`. | Sediakan nama executable dan arsitektur yang cocok di `build.output`. Flag ini merupakan pilihan perintah satu kali, bukan pengaturan profil permanen. |
| Tambahkan otomatisasi rilis | Profil bernama, hook siklus hidup terbatas, caching, dan pekerjaan pemaketan paralel. | Gunakan YAML literal dan skrip milik proyek yang dideklarasikan. Perilaku Taskfile arbitrer bukan mekanisme ekstensi Wake. |

Edit `profiles.setup`, profil bernama lainnya, dan pengaturan lanjutan di `wails.yaml`. Gunakan `wails3 setup` untuk mengubah preferensi dan kredensial bersama. Gunakan [panduan build YAML](/experimental/yaml-builds/) untuk contoh lengkap dan persyaratan paket. Aplikasi GUI tetap memerlukan pustaka sistem di tujuan; Pure Go tidak menjadikannya aplikasi yang sepenuhnya statis.

iOS dan Android masih bersifat eksperimental, memerlukan SDK platform masing-masing dan CGo, serta berada di luar cakupan pelajaran desktop pertama ini. Wizard tidak membuat build seluler menjadi portabel di sembarang host.

## Migrasikan proyek yang sudah ada

Kerjakan pada cabang terpisah atau salinan dari proyek **Wails v3** yang masih berfungsi. Ini adalah migrasi konfigurasi build, bukan peningkatan dari Wails v2. Pertahankan hasil build dan pengujian yang ada agar Anda dapat membandingkan perilaku setelah aktivasi.

@note{type="note" title="Mulai dengan migrasi yang telah ditinjau"}
Proyek baru yang dihasilkan dari templat lama bawaan saat ini dapat dimigrasikan sebelum frontend-nya dibangun. Wake mengenali pemilih task standar dan siklus hidup pengembangannya tanpa mengeksekusi Taskfile atau mengevaluasi ekspresi khusus.

Task khusus, penimpaan variabel, dan kebijakan pengembangan masih perlu ditinjau. Jika analisis tidak dapat merepresentasikan proyek Anda, pertahankan build lama dan laporkan diagnostiknya. Jangan hapus perilaku yang Anda perlukan hanya agar aktivasi berhasil. Tutorial proyek baru di atas tidak memerlukan migrasi.
@end

### 1. Analisis tanpa mengubah build aktif

Aktifkan Wake dan gunakan CLI pratinjau seperti yang dijelaskan di awal halaman ini. Dari direktori akar proyek yang ada, jalankan:

```sh
wails3 migrate --dry-run
```

Baca pengaturan yang dilaporkan sebagai didukung serta faktor penghambatnya. Penganalisis membaca deklarasi lama; penganalisis tidak mengeksekusi perintah di dalamnya. Dry run tidak membuat manifes aktif. Untuk diagnostik yang dapat dibaca mesin, gunakan `wails3 migrate --json`; gunakan `--dry-run --json` jika Anda juga ingin menghindari penulisan draf.

### 2. Buat dan tinjau proposal yang tidak aktif

```sh
wails3 migrate
```

Tindakan ini membuat **`wails.migrated.yaml`**, bukan `wails.yaml`. Periksa file tersebut bersama Taskfile yang ada dan `build/config.yml`. Periksa metadata proyek, perintah frontend, arsitektur target, isi paket, aset, dan pengaturan penandatanganan. Simpan rahasia dalam variabel lingkungan atau penyimpanan kredensial Anda.

Task bawaan yang dikenal dapat dipetakan ke pengaturan native. Task siklus hidup konservatif tanpa argumen yang memanggil skrip milik proyek dapat dijadikan hook. Perilaku task khusus lain yang dapat dijangkau merupakan penghambat, bukan sesuatu yang diam-diam diabaikan oleh Wake atau dijalankan melalui Task sebagai fallback.

Jika draf sudah ada, migrasi tidak akan menimpa hasil edit yang telah Anda tinjau. Simpan salinannya sebelum sengaja mengganti atau membuat ulang draf tersebut. Jangan ganti nama draf yang terhambat menjadi `wails.yaml` untuk melewati pemeriksaan aktivasi. Kurangi atau ganti perilaku lama yang tidak didukung, lalu jalankan kembali analisis sebelum aktivasi. Jika proyek belum dapat direpresentasikan, pertahankan build lamanya dan laporkan kesenjangan tersebut.

### 3. Aktifkan secara eksplisit

Setelah meninjau draf dan mengatasi faktor penghambat, jalankan:

```sh
wails3 migrate --activate --dry-run
wails3 migrate --activate
```

Perintah pertama memeriksa aktivasi tanpa menerapkannya. Perintah kedua memeriksa ulang deklarasi lama dan memvalidasi draf sebelum membuat `wails.yaml` secara atomik; aktivasi yang berhasil menghapus `wails.migrated.yaml` yang tidak aktif. Taskfile, `build/config.yml`, skrip, dan aset tidak diubah. Sejak saat ini, manifes aktif mengendalikan build Wake. File lama bukan sumber konfigurasi kedua.

Beberapa proyek lama sudah memiliki `wails.yaml` khusus metadata dengan bagian `info:` pada tingkat teratas. Meskipun bukan manifes native Wake, file tersebut tetap menempati tujuan aktivasi. Pertahankan lalu pindahkan file lama tersebut secara sengaja sebelum aktivasi; migrator tidak akan menimpa atau mengganti namanya untuk Anda.

### 4. Verifikasi aplikasi yang telah dimigrasikan

```sh
wails3 config check
wails3 build
wails3 dev
```

Bandingkan aplikasi dan isi paket dengan build sebelumnya, termasuk perilaku rilis khusus Anda. Jalankan pengujian proyek Anda. Periksa setiap profil yang dimigrasikan secara eksplisit, misalnya `wails3 config check release` dan `wails3 build release` jika draf berisi profil `release`.

Edit profil yang dimigrasikan secara langsung di `wails.yaml`. `wails3 init -ui` membuat proyek baru; `wails3 setup` mengelola preferensi global. Keduanya bukan editor untuk profil rilis proyek yang sudah ada.

### Jika Anda perlu menghentikan eksperimen

Menghapus pengaturan `WAILS_EXP_USE_WAKE` akan menonaktifkan keikutsertaan, tetapi `wails.yaml` yang aktif tetap memblokir perutean otomatis ke Taskfile. YAML yang tidak valid juga menyebabkan proses berhenti dengan kesalahan. Tidak ada perintah migrasi balik.

Untuk kembali ke alur kerja lama, pulihkan revisi lengkap proyek sebelum migrasi di sistem kontrol versi, termasuk kondisi tanpa manifes aktif, lalu gunakan CLI dan dependensi yang kompatibel dengannya. Sebelum melakukannya, pertahankan setiap perubahan sumber baru yang Anda perlukan. Menggunakan cabang uji coba terpisah sangat memudahkan pemulihan yang disengaja ini.

## Diagnosis masalah dan berikan umpan balik

| Gejala | Pemeriksaan pertama |
| --- | --- |
| Pembuatan proyek melewati Application, Targets, Packages, dan Tools | Tetapkan `WAILS_EXP_USE_WAKE` di terminal yang digunakan untuk menjalankan proses, lalu periksa executable `wails3` mana yang terdapat di `PATH`. `wails3 setup` global selalu mempertahankan alur preferensinya. |
| Templat komunitas perlu dimigrasikan | Tinjau konfigurasi lamanya dengan perintah migrasi sebelum aktivasi. |
| Build yang disimpan melaporkan alat yang tidak tersedia | Pembuatan proyek mencatat konfigurasi yang diinginkan. Instal persyaratan host yang dilaporkan sebelum membangun; pembuatan proyek tidak menginstal setiap alat yang dipilih. |
| Executable hasil cross-build tidak dapat dijalankan di sini | Jalankan executable tersebut pada sistem operasi dan arsitektur yang dipilih; gunakan target host untuk pengembangan lokal. |
| Mode prebuilt tidak dapat menemukan biner | Periksa [nama file dan jalur prebuilt](/experimental/yaml-builds/#use-prebuilt-desktop-executables), termasuk subdirektori untuk beberapa target. |
| Perubahan manifes gagal | Jalankan `wails3 config check` untuk profil yang relevan dan ikuti diagnostik bidangnya. Wake tidak beralih kembali ke Taskfile sebagai fallback. |

Periksa build sebelum menjalankannya dengan `wails3 build setup --plan --json`. `wails3 config check` memvalidasi konfigurasi, dan `wails3 eject` menulis snapshot `wails.ejected.yaml` yang tidak aktif tanpa mengubah manifes aktif. Tinjau keluaran perintah sebelum membagikannya: jalur proyek dan lokasi kredensial dapat mengidentifikasi lingkungan lokal Anda.

Bagikan umpan balik di [Diskusi Eksperimen](https://github.com/wailsapp/wails/discussions/categories/experiments) atau [laporkan bug yang dapat direproduksi](https://github.com/wailsapp/wails/issues). Sertakan revisi CLI dan framework Anda, OS/arsitektur host, target yang dipilih, perintah persisnya, apakah Anda memilih Pure Go, serta diagnostik yang relevan. Nyatakan apakah kegagalan terjadi saat kompilasi, pemaketan, penandatanganan, atau menjalankan aplikasi. Jika berguna, sertakan manifes minimal yang rahasianya telah dihapus serta keluaran `wails3 doctor`.
