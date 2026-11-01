---
title: "Wake"
description: "Buat proyek Wake pertama Anda, ikuti wizard penyiapan, build dan jalankan proyek tersebut, lalu migrasikan proyek Taskfile yang sudah ada."
slug: "experimental/wake"
sourcePath: "experimental/wake.md"
---

Wake adalah sistem build eksperimental Wails. Satu `wails.yaml` mendeskripsikan aplikasi, frontend, platform target, paket, dan pengaturan penandatanganan Anda. Wake menggunakan konfigurasi tersebut untuk merencanakan dan menjalankan build, menyimpan pekerjaan yang tidak berubah dalam cache, serta melaporkan berkas yang dihasilkan.

Dalam tutorial ini, Anda akan membuat aplikasi desktop kecil, mengonfigurasinya di browser, mem-build-nya untuk komputer Anda sendiri, dan menjalankannya dalam mode pengembangan. Mulailah dengan satu target lokal dan build tanpa tanda tangan. Tambahkan kompilasi silang dan distribusi setelah build pertama tersebut berhasil.

Sudah memiliki aplikasi? Buka [Memigrasikan proyek yang sudah ada](#migrasikan-proyek-yang-sudah-ada). Untuk masing-masing kolom dan alur kerja tingkat lanjut, gunakan [panduan build YAML](/experimental/yaml-builds/) dan [referensi manifes](/reference/wails-yaml-fields/).

## Hal yang dapat diharapkan dari eksperimen ini

@note{type="caution" title="Gunakan proyek pengujian atau cabang terpisah"}
Wake bersifat opsional dan konfigurasi, perintah, serta perilakunya dapat berubah di antara build pratinjau. Wake belum memiliki jaminan kompatibilitas seperti sistem build stabil. Simpan revisi proyek yang diketahui berfungsi dan catat versi CLI serta framework yang Anda uji. Uji paket pada sistem tujuannya sebelum mendistribusikannya.
@end

Wake dan backend **Pure Go** merupakan pilihan terpisah. Mengaktifkan Wake tidak otomatis mengaktifkan Pure Go. Pelajaran ini menggunakan backend standar; karena itu, build GUI macOS dan Linux memerlukan alat pengembangan native yang biasanya digunakan.

Wizard memeriksa persyaratan dan menyimpan konfigurasi. Menyimpan bukanlah proses build dan tidak membuktikan bahwa setiap target yang dipilih sudah siap. Alat atau kredensial yang belum tersedia dapat ditambahkan nanti, misalnya pada host CI. Kombinasi target atau paket yang tidak valid harus diperbaiki sebelum menyimpan.

Migrasi bersifat **satu arah**. Wake tidak menjalankan Taskfile, memperluas templat Task, ataupun menerjemahkan YAML kembali menjadi Taskfile. Proyek yang belum dimigrasikan tetap menggunakan runtime Task lama yang terpisah. `wails.yaml` yang aktif tidak pernah secara diam-diam beralih kembali ke runtime tersebut.

## Sebelum memulai

Anda memerlukan Go, Node.js beserta npm untuk frontend awal, serta [dependensi platform](/getting-started/installation/#platform-specific-dependencies) untuk komputer Anda. Mulailah dari direktori tempat Anda ingin membuat folder proyek baru. Anda juga memerlukan browser untuk menggunakan wizard.

Gunakan revisi CLI dan framework yang menyertakan pratinjau ini. Menginstal CLI `v3.0.0` stabil tidak menambahkan layar penyiapan Wake yang dijelaskan di sini. Perintah berikut menyematkan pratinjau dari [PR #6531](https://github.com/wailsapp/wails/pull/6531):

```sh
go install github.com/wailsapp/wails/v3/cmd/wails3@0243da963
wails3 version
```

Pastikan direktori executable Go tercantum dalam `PATH`; [panduan instalasi](/getting-started/installation/) menjelaskan caranya. Nanti dalam tutorial ini, Anda akan menyematkan dependensi Go aplikasi ke revisi yang sama. Saat menguji pratinjau yang lebih baru, perbarui kedua argumen revisi secara bersamaan.

### Aktifkan Wake di terminal Anda

Di Bash atau Zsh:

```bash
export WAILS_EXP_USE_WAKE=1
```

Di PowerShell:

```powershell
$env:WAILS_EXP_USE_WAKE = '1'
```

Biarkan terminal ini tetap terbuka selama mengikuti tutorial. Tetapkan kembali variabel tersebut di setiap terminal baru atau job CI yang menjalankan perintah Wake. **Keberadaan** variabel tersebut mengaktifkan Wake: menetapkannya ke `false` atau string kosong tidak menonaktifkannya. Nama lama `WAILS_USE_WAKE` diabaikan.

## Buat aplikasi

Jalankan wizard pembuatan proyek:

```sh
wails3 init -ui
```

Wizard proyek akan terbuka di browser. Wizard ini membuat berkas sumber; wizard `wails3 setup` terpisah di bawah mengonfigurasi cara mem-build berkas tersebut.

| Layar | Pilihan untuk tutorial ini | Yang terjadi selanjutnya |
| --- | --- | --- |
| **Proyek — Buat proyek Wails** | Masukkan `wake-demo`, lalu pilih **Lanjutkan**. | Proyek akan dibuat dalam folder `wake-demo` baru di bawah direktori yang ditampilkan. |
| **Framework — Pilih framework** | Pilih **Vanilla**, lalu pilih **Lanjutkan**. | Anda memperoleh frontend kecil tanpa perlu mempelajari langkah khusus framework. |
| **Bahasa — Preferensi bahasa** | Pilih **JavaScript**, lalu pilih **Lanjutkan**. | Layar **Binding** yang hanya tersedia untuk TypeScript akan dilewati. |
| **Detail — Detail proyek** | Gunakan nama produk `Wake Demo`, pengidentifikasi bundle `com.example.wakedemo`, dan versi `0.1.0`. Klik **Buat Proyek**. | Saat Wake diaktifkan, detail tersebut ditulis ke `wails.yaml`. |

Jangan gunakan **Gunakan default global** untuk pelajaran ini karena opsi tersebut melewati pilihan framework dan bahasa. Jika memilih TypeScript dalam proyek lain, layar **Binding** memungkinkan Anda memilih antarmuka atau kelas untuk binding data Go yang dihasilkan.

Kembali ke terminal dan tunggu hingga pembuatan proyek selesai. Kemudian jalankan:

```sh
cd wake-demo
go get github.com/wailsapp/wails/v3@0243da963
```

Sekarang Anda seharusnya memiliki `main.go`, sebuah modul Go, direktori `frontend`, dan `wails.yaml` native. Templat bawaan menggunakan Wake secara langsung; Anda tidak perlu memigrasikan proyek baru ini. Templat komunitas dapat memiliki task lama kustom sehingga mungkin memerlukan peninjauan migrasi yang dijelaskan di bawah.

## Konfigurasikan build, layar demi layar

Dari dalam `wake-demo`, jalankan:

```sh
wails3 setup
```

Browser kini menampilkan **PENYIAPAN WAKE**. Jika yang ditampilkan adalah alur penyiapan umum, pastikan `WAILS_EXP_USE_WAKE` telah ditetapkan di terminal yang menjalankan perintah tersebut dan Anda menggunakan CLI pratinjau.

Tangkapan layar di bawah menampilkan panduan ini di Linux. Target host dan alat yang diperlukan akan berbeda di Windows atau macOS.

### Proyek — hanya jika metadata belum tersedia

Untuk aplikasi yang baru saja Anda buat, layar ini **dilewati** karena `init` telah menulis metadata. Hal ini memang semestinya terjadi.

Saat mengonfigurasi direktori sumber Go yang sudah ada tanpa manifes native atau konfigurasi build lama, **Tentang proyek Anda** meminta nama produk, nama proyek, pengidentifikasi aplikasi, versi, perusahaan, dan deskripsi. Nilai-nilai tersebut tetap berada di wizard hingga Anda memilih **Simpan**. Penyiapan tidak membuat berkas sumber aplikasi. Jika menemukan Taskfile atau metadata lama, lakukan migrasi terlebih dahulu.

### Aplikasi — Apa yang sedang Anda buat?

Pilih **Aplikasi desktop**. Biarkan **Paketkan binary yang sudah ada** tidak dicentang dan pilih **Backend standar**. Klik **Lanjutkan**.

![Layar Aplikasi dengan Aplikasi desktop dan Backend standar dipilih serta pemaketan binary bawaan dinonaktifkan.](/assets/screenshots/wake/application.png)

- **Aplikasi desktop** mem-build GUI dengan frontend web.
- **Aplikasi baris perintah** dan **Server** mendeskripsikan program tanpa jendela desktop. Mengubah pilihan ini tidak menulis ulang kode sumber aplikasi Anda; gunakan proyek sumber yang sesuai dengan jenis yang dipilih.
- **Pure Go (eksperimental)** adalah opsi terpisah yang harus diaktifkan secara eksplisit. Opsi ini menghilangkan persyaratan compiler C untuk kompilasi desktop yang didukung, asalkan dependensi Anda juga dapat di-build tanpa CGo. Opsi ini tidak menghilangkan library runtime target ataupun semua alat pemaketan dan penandatanganan.
- **Kemas biner yang sudah ada** menyiapkan perintah menggunakan `--skip-build`.
  Opsi ini memerlukan biner yang sudah Anda build; biarkan nonaktif untuk proyek baru.

Backend merupakan **pengaturan tingkat proyek**. Perilaku profil lain yang mewarisi pengaturan tersebut dapat berubah saat Anda mengubahnya di sini.

### Target — Di mana aplikasi harus dijalankan?

Pertahankan hanya target untuk komputer yang sedang Anda gunakan, lalu klik **Lanjutkan**.
Pada awalnya, wizard memilih sistem operasi dan arsitektur host.

![Layar Target dengan hanya linux/amd64 yang dipilih untuk host tutorial.](/assets/screenshots/wake/targets.png)

Nama target menggabungkan sistem operasi dan arsitektur CPU: misalnya,
`windows/amd64`, `darwin/arm64` untuk Mac Apple Silicon, atau `linux/amd64`.
`darwin` berarti macOS. Memilih target lain berarti meminta executable untuk mesin tersebut; tindakan ini tidak membuat executable itu dapat dijalankan di komputer Anda saat ini.

Biarkan arsitektur tambahan, build universal macOS, dan target seluler eksperimental untuk proses penyiapan berikutnya.

### Paket — Pilih keluaran rilis

Biarkan kotak format paket tidak dicentang dan nonaktifkan **Tandatangani keluaran rilis**.
Klik **Lanjutkan**.

![Layar Paket tanpa format paket Linux tambahan yang dipilih dan dengan penandatanganan dinonaktifkan.](/assets/screenshots/wake/packages.png)

Jika tidak ada format tambahan yang dipilih, Wake menghasilkan keluaran default target: executable desktop, atau bundel aplikasi macOS untuk GUI macOS. Ini cukup untuk memeriksa aplikasi secara lokal. Memilih format nanti akan menambahkan keluaran distribusi seperti DMG, MSIX, NSIS, DEB, atau AppImage jika didukung.

Penandatanganan bersifat opsional untuk build lokal pertama ini. Build tanpa tanda tangan yang berhasil tidak berarti aplikasi siap didistribusikan kepada publik.

### Alat — Alat untuk build Anda

Baca persyaratan untuk pilihan Anda. Daftar tersebut berasal dari rencana build Wake; daftar itu bukan permintaan untuk memasang setiap alat di setiap platform.

![Layar Alat yang menampilkan persyaratan yang telah ditentukan untuk build dari sumber secara lokal di Linux.](/assets/screenshots/wake/tools.png)

Untuk pelajaran ini, Anda memerlukan Go dan alat frontend, ditambah alat pengembangan native yang diwajibkan oleh backend GUI standar host Anda. Docker bukan persyaratan umum untuk semua build. Cross-build dengan CGo dapat memerlukan compiler dan pustaka target atau toolchain kontainer yang didukung.

Selesaikan baris yang ditandai **Belum siap**, **Perlu dikonfigurasi**, atau **Perlu diperbarui** sebelum mencoba build lokal. Gunakan petunjuk yang ditampilkan jika tersedia, atau [panduan pemasangan platform](/getting-started/installation/), lalu klik **Periksa ulang**. Anda mungkin perlu memulai ulang penyiapan jika pemasangan mengubah `PATH` terminal. Setelah persyaratan host siap, klik **Lanjutkan**.

Wizard dapat menyimpan konfigurasi yang valid meskipun komputer ini tidak memiliki alat yang diperlukan. Hal ini berguna untuk konfigurasi CI, tetapi dalam tutorial ini Anda memerlukan build yang dapat dijalankan di sini.

### Penandatanganan — hanya jika dipilih

Layar ini **dilewati** karena Anda membiarkan penandatanganan nonaktif. Pada proses berikutnya, mengaktifkan **Tandatangani keluaran rilis** akan menampilkan kontrol penandatanganan yang tersedia untuk platform desktop yang Anda pilih. **Siapkan penandatanganan nanti** menonaktifkan penandatanganan untuk profil penyiapan yang disimpan. Penandatanganan seluler dikonfigurasi dalam YAML, bukan melalui kontrol ini.

Pengaturan default penandatanganan berada di pengaturan pengguna Anda. Wake dapat menggunakan file sertifikat terkonfigurasi yang berada di luar proyek tanpa menyalin kunci privat atau kata sandi ke dalam `wails.yaml`. Pengaturan penandatanganan proyek yang eksplisit lebih diutamakan. Ikuti [panduan penandatanganan saat penyiapan](/getting-started/setup/) dan [panduan pengemasan dan penandatanganan](/experimental/yaml-builds/#package-and-sign-for-distribution) sebelum menyiapkan rilis publik.

### Simpan — Simpan konfigurasi build Anda

Periksa ringkasannya: aplikasi desktop, backend standar, build dari sumber, target host Anda, dan penandatanganan dinonaktifkan. Klik **Simpan wails.yaml**.

![Layar Simpan yang merangkum aplikasi GUI, backend standar, build dari sumber, dan target linux/amd64 tanpa tanda tangan.](/assets/screenshots/wake/save.png)

Anda akan melihat **Konfigurasi build disimpan**, jalur ke manifes Anda, dan:

```sh
wails3 build setup
```

Klik **Selesai** dan kembali ke terminal. Penyiapan telah menyimpan pilihan target dan paket di bawah `profiles.setup`. Penyiapan juga memperbarui jenis aplikasi, pengaktifan frontend, dan pengaturan backend. Perintah frontend, program pendamping, hook, dan konfigurasi lain yang sudah ada tetap dipertahankan.

Nama profil penting: `wails3 build setup` menggunakan pilihan yang disimpan tersebut.
`wails3 build` tanpa argumen tambahan tidak otomatis memilih profil `setup`.

## Build dan jalankan aplikasi pertama Anda

Jalankan perintah yang ditampilkan oleh wizard:

```sh
wails3 build setup
```

Build pertama memasang dependensi frontend, menghasilkan binding, dan mengompilasi aplikasi. Build berikutnya dapat menggunakan kembali hasil pekerjaan yang tidak berubah. Tunggu hingga muncul **build berhasil** dan baca jalur keluaran yang dicetak oleh CLI. Direktori keluaran default adalah `bin`; GUI macOS dirakit menjadi bundel `.app`.

Mulai mode pengembangan di komputer Anda sendiri:

```sh
wails3 dev
```

Jendela aplikasi awal akan muncul. Masukkan nama Anda dan coba sapaannya. Ubah placeholder input nama di `frontend/index.html` lalu simpan file; frontend akan dimuat ulang. Edit sapaan di `greetservice.go` untuk mencoba build ulang Go. Tekan **Ctrl+C** di terminal setelah selesai.

Pengembangan merupakan alur kerja di host lokal. Profil rilis dapat memuat beberapa target cross-build tanpa mengubah mesin mana yang dapat menjalankan jendela pengembangan. Untuk mem-build dan menjalankan aplikasi produksi secara lokal, gunakan `wails3 run`; lihat [referensi perintah run](/reference/wails-yaml-run/).

Kini Anda memiliki proyek sumber, profil build yang telah ditinjau, dan siklus pengembangan lokal yang berfungsi. Commit sumber dan `wails.yaml` ke kontrol versi. File kerja `.wails/` yang dihasilkan bukan sumber konfigurasi Anda.

## Yang dapat Anda lakukan selanjutnya

| Tujuan | Dukungan Wake | Yang perlu diperiksa terlebih dahulu |
| --- | --- | --- |
| Build untuk desktop lain | Target Windows, macOS, dan Linux, termasuk kedua arsitektur macOS serta keluaran universal. | Pure Go dan CGo memiliki persyaratan compiler yang berbeda. Uji di sistem operasi dan arsitektur tujuan. |
| Buat installer atau paket | MSIX, NSIS, aplikasi/DMG macOS, DEB, RPM, Arch, Void, dan AppImage. | NSIS memerlukan `makensis`. AppImage memerlukan runtime target; cross-packaging GUI juga memerlukan pohon runtime Linux yang disediakan. |
| Tandatangani rilis | Penandatanganan portabel dengan kunci yang diekspor untuk Windows, macOS, dan DEB/RPM/Arch; penyimpanan identitas native tetap tersedia. | Setiap identitas, format, dan host memiliki batasan. AppImage dan Void tidak memiliki alur penandatanganan bawaan. |
| Notarisasikan aplikasi macOS | Pengiriman dengan kunci API tim dan penyematan tiket, atau profil keychain macOS. | Tanda tangan Developer ID dan stempel waktu yang sesuai tetap diperlukan. Lolosnya tanda tangan pengujian yang dihasilkan Go dalam pemeriksaan native bukanlah bukti bahwa Apple akan menerima rilis Anda. |
| Build beberapa program | Program utama berupa GUI, CLI, atau server, ditambah program pendamping bernama. | Deklarasikan paket dan jenis program pendamping dalam YAML. Wizard tidak menghasilkan kode sumbernya. |
| Kemas file yang di-build di tempat lain | `wails3 build setup --skip-build`. | Sediakan nama executable dan arsitektur yang sesuai di `build.output`. Flag ini merupakan pilihan satu kali untuk perintah tersebut, bukan pengaturan penyiapan permanen. |
| Tambahkan otomatisasi rilis | Profil bernama, hook siklus hidup terbatas, caching, dan pengerjaan paket secara paralel. | Gunakan YAML literal dan skrip milik proyek yang dideklarasikan. Perilaku Taskfile arbitrer bukan mekanisme ekstensi Wake. |

Buka kembali `wails3 setup` untuk mengubah profil `setup`. Profil bernama lainnya dan pengaturan lanjutan diedit di `wails.yaml`. Lihat [panduan build YAML](/experimental/yaml-builds/) untuk contoh lengkap dan persyaratan paket. Aplikasi GUI tetap memerlukan pustaka sistem pada tujuan; Pure Go tidak menjadikannya aplikasi yang sepenuhnya statis.

iOS dan Android masih bersifat eksperimental, memerlukan SDK platform masing-masing dan CGo, serta berada di luar cakupan pelajaran desktop pertama ini. Wizard tidak membuat build seluler menjadi portabel di sembarang host.

## Migrasikan proyek yang sudah ada

Kerjakan di branch terpisah atau salinan dari working copy proyek **Wails v3** yang berfungsi. Ini adalah migrasi konfigurasi build, bukan upgrade dari Wails v2. Pertahankan hasil build dan pengujian yang ada agar Anda dapat membandingkan perilakunya setelah aktivasi.

@note{type="note" title="Mulai dengan migrasi yang telah ditinjau"}
Proyek baru yang dihasilkan dari templat lama bawaan saat ini dapat dimigrasikan sebelum frontend-nya dibangun. Wake mengenali pemilih tugas standar dan siklus hidup pengembangannya tanpa menjalankan Taskfile atau mengevaluasi ekspresi khusus.

Tugas khusus, penggantian nilai variabel, dan kebijakan pengembangan tetap perlu ditinjau. Jika analisis tidak dapat merepresentasikan proyek Anda, tetap gunakan sistem build lama dan laporkan hasil diagnostiknya. Jangan hapus perilaku yang Anda perlukan hanya agar aktivasi berhasil. Tutorial proyek baru di atas tidak memerlukan migrasi.
@end

### 1. Analisis tanpa mengubah build aktif

Aktifkan Wake dan gunakan CLI pratinjau dari awal halaman ini. Dari root proyek yang sudah ada, jalankan:

```sh
wails3 migrate --dry-run
```

Baca pengaturan yang didukung dan penghambat yang dilaporkan. Penganalisis membaca deklarasi lama; alat ini tidak menjalankan perintahnya. Dry run tidak membuat manifes aktif. Untuk diagnostik yang dapat dibaca mesin, gunakan `wails3 migrate --json`; gunakan `--dry-run --json` jika Anda juga ingin menghindari penulisan draf.

### 2. Buat dan tinjau proposal yang tidak aktif

```sh
wails3 migrate
```

Langkah ini membuat **`wails.migrated.yaml`**, bukan `wails.yaml`. Periksa file tersebut bersama Taskfile dan `build/config.yml` yang sudah ada. Periksa metadata proyek, perintah frontend, arsitektur target, isi paket, aset, dan pengaturan penandatanganan. Simpan rahasia dalam variabel lingkungan atau penyimpanan kredensial Anda.

Task bawaan yang dikenali dapat dipetakan ke pengaturan native. Task siklus hidup konservatif tanpa argumen yang memanggil skrip milik proyek dapat dijadikan hook. Perilaku task kustom lain yang dapat dijangkau merupakan penghambat, bukan sesuatu yang dihapus secara diam-diam oleh Wake atau dijalankan melalui Task sebagai fallback.

Jika draf sudah ada, migrasi tidak akan menimpa hasil edit yang telah Anda tinjau. Simpan salinannya sebelum sengaja mengganti atau membuatnya ulang. Jangan ganti nama draf yang terhambat menjadi `wails.yaml` untuk melewati pemeriksaan aktivasi. Kurangi atau ganti perilaku lama yang tidak didukung, lalu jalankan ulang analisis sebelum aktivasi. Jika proyek belum dapat direpresentasikan, pertahankan build lamanya dan laporkan kekurangan tersebut.

### 3. Aktifkan secara eksplisit

Setelah meninjau draf dan mengatasi penghambat, jalankan:

```sh
wails3 migrate --activate --dry-run
wails3 migrate --activate
```

Perintah pertama memeriksa aktivasi tanpa menerapkannya. Perintah kedua memeriksa ulang deklarasi lama dan memvalidasi draf sebelum membuat `wails.yaml` secara atomik; aktivasi yang berhasil menghapus `wails.migrated.yaml` yang tidak aktif. Perintah ini tidak mengubah Taskfile, `build/config.yml`, skrip, maupun aset. Sejak titik ini, manifes aktif mengendalikan build Wake. File lama bukan sumber konfigurasi kedua.

Beberapa proyek lama sudah memiliki `wails.yaml` khusus metadata dengan bagian `info:` pada tingkat teratas. Meskipun bukan manifes Wake native, file tersebut tetap menempati tujuan aktivasi. Pertahankan dan pindahkan file lama itu secara sengaja sebelum aktivasi; migrator tidak akan menimpa atau mengganti namanya untuk Anda.

### 4. Verifikasi aplikasi yang telah dimigrasikan

```sh
wails3 config check
wails3 build
wails3 dev
```

Bandingkan isi aplikasi dan paket dengan build sebelumnya, termasuk perilaku rilis kustom Anda. Jalankan pengujian proyek. Periksa setiap profil yang dimigrasikan secara eksplisit, misalnya `wails3 config check release` dan `wails3 build release` jika draf berisi profil `release`.

Sekarang Anda dapat menjalankan `wails3 setup` untuk menambahkan atau mengganti profil khusus `setup`. Perintah ini tidak otomatis memilih atau mengedit profil `release` yang sudah ada.

### Jika Anda perlu menghentikan eksperimen

Menghapus pengaturan `WAILS_EXP_USE_WAKE` akan menonaktifkan keikutsertaan, tetapi `wails.yaml` yang aktif tetap memblokir perutean otomatis ke Taskfile. YAML yang tidak valid juga menghentikan proses dengan galat. Tidak ada perintah migrasi balik.

Untuk kembali ke alur kerja lama, pulihkan revisi lengkap proyek sebelum migrasi dalam kontrol versi, termasuk kondisi tanpa manifes aktif, lalu gunakan CLI dan dependensi yang kompatibel dengannya. Pertahankan setiap perubahan sumber baru yang Anda perlukan sebelum melakukannya. Menggunakan branch uji coba terpisah sangat memudahkan pemulihan yang disengaja ini.

## Diagnosis masalah dan berikan umpan balik

| Gejala | Pemeriksaan pertama |
| --- | --- |
| Penyiapan menampilkan wizard umum | Tetapkan `WAILS_EXP_USE_WAKE` di terminal yang digunakan untuk menjalankan proses, lalu periksa executable `wails3` mana yang tersedia di `PATH`. |
| Penyiapan meminta migrasi | Anda membuka proyek lama atau templat komunitas dengan konfigurasi lama. Tinjau migrasi sebelum aktivasi. |
| Build yang disimpan melaporkan alat yang tidak tersedia | Penyimpanan mencatat maksud konfigurasi. Kembali ke **Alat**, penuhi persyaratan host, lalu periksa kembali. |
| Executable hasil cross-build tidak dapat dijalankan di sini | Jalankan di sistem operasi dan arsitektur yang dipilih; gunakan target host untuk pengembangan lokal. |
| Mode prebuilt tidak dapat menemukan binary | Periksa [nama file dan path prebuilt](/experimental/yaml-builds/#use-prebuilt-desktop-executables), termasuk subdirektori untuk beberapa target. |
| Perubahan manifes gagal | Jalankan `wails3 config check` untuk profil yang relevan dan ikuti diagnostik field. Wake tidak menggunakan Taskfile sebagai fallback. |

Periksa build sebelum menjalankannya dengan `wails3 build setup --plan --json`.
`wails3 config check` memvalidasi konfigurasi, sedangkan `wails3 eject` menulis snapshot `wails.ejected.yaml` yang tidak aktif tanpa mengubah manifes aktif.
Tinjau output perintah sebelum membagikannya: jalur proyek dan lokasi kredensial dapat mengidentifikasi lingkungan lokal Anda.

Bagikan umpan balik di [Diskusi Eksperimen](https://github.com/wailsapp/wails/discussions/categories/experiments) atau [laporkan bug yang dapat direproduksi](https://github.com/wailsapp/wails/issues). Sertakan revisi CLI dan framework, OS/arsitektur host, target yang dipilih, perintah persisnya, apakah Anda memilih Pure Go, serta diagnostik yang relevan. Sebutkan apakah kegagalan terjadi saat kompilasi, pengemasan, penandatanganan, atau saat menjalankan aplikasi. Sertakan manifes minimal tanpa rahasia dan output `wails3 doctor` jika berguna.
