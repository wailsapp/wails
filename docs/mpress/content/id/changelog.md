---
title: "Log perubahan"
description: "Riwayat versi dan catatan rilis untuk Wails v3"
slug: "changelog"
sourcePath: "changelog.md"
---

Keterangan:

-  - macOS
- ⊞ - Windows
- 🐧 - Linux

/_-- Semua perubahan penting pada proyek ini akan didokumentasikan dalam berkas ini.

Format ini didasarkan pada [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), dan proyek ini mengikuti [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

- `Added` untuk fitur baru.
- `Changed` untuk perubahan pada fungsionalitas yang sudah ada.
- `Deprecated` untuk fitur yang akan segera dihapus.
- `Removed` untuk fitur yang kini telah dihapus.
- `Fixed` untuk setiap perbaikan bug.
- `Security` jika terdapat kerentanan.

_/

/_   * MOHON JANGAN PERBARUI BERKAS INI *   Pembaruan sebaiknya ditambahkan ke `v3/UNRELEASED_CHANGELOG.md`   Terima kasih! _/

## [Belum dirilis]

## v3.0.0-beta.21 - 2026-09-13

## Ditambahkan

- Sajikan dokumentasi v3 Wails dengan M-Press dalam [PR](https://github.com/wailsapp/wails/pull/6116) oleh @leaanthony

## Diperbaiki

- Uraikan nilai slug JSON dari frontmatter MPD untuk pembuatan changelog dalam [PR](https://github.com/wailsapp/wails/pull/6118) oleh @leaanthony
- Updater menghapus variabel lingkungan helper dan menjalankan ulang target asli setelah kegagalan pencadangan dalam [PR](https://github.com/wailsapp/wails/pull/6080) oleh @cnmax
- Mulai penangan sinyal bawaan selama App.Run dalam [PR](https://github.com/wailsapp/wails/pull/6098) oleh @leaanthony
- Menu Windows menangani menu nil, membebaskan sumber daya yang diganti, dan menggambar ulang bilah menu dalam [PR](https://github.com/wailsapp/wails/pull/6112) oleh @taliesin-ai
- Pulihkan pemaketan MSIX untuk proyek baru yang menggunakan konfigurasi YAML bersama dalam [PR](https://github.com/wailsapp/wails/pull/6115) oleh @leaanthony
- Perbaiki binding JavaScript dan TypeScript yang dihasilkan yang gagal dimuat ketika pembuat model generik mereferensikan deklarasi helper yang muncul kemudian, serta cegah stack overflow saat membuat model generik yang saling bergantung (#6062)

## v3.0.0-beta.20 - 2026-09-10

## Diubah

- Perbarui tautan showcase Clave agar mengarah ke situs web dan repositori saat ini dalam [PR](https://github.com/wailsapp/wails/pull/6082) oleh @01xR4in

## Diperbaiki

- Batalkan permintaan aset Windows yang dihentikan, termasuk permintaan worker, sambil mempertahankan handler keepalive selama navigasi. Teruskan konteks permintaan native melalui wrapper aplikasi pada platform Apple. (#5963, #5969)
- Pertahankan entri log perubahan saat terjadi push yang bersaing melalui percobaan ulang dalam [PR](https://github.com/wailsapp/wails/pull/6094) oleh @leaanthony
- Perbaiki kegagalan `go mod vendor` dengan `pattern arm64/WebView2Loader.dll: no matching files found` pada setiap platform dengan menghapus embed yang merujuk ke berkas biner yang tidak pernah disertakan dalam modul, sehingga memperbaiki [#5782](https://github.com/wailsapp/wails/issues/5782) dan [#5376](https://github.com/wailsapp/wails/issues/5376), dalam [PR](https://github.com/wailsapp/wails/pull/6031) oleh @Grantmartin2002

## Dihapus

- Hapus dukungan loader WebView2 native yang telah digantikan oleh loader Go murni. Perubahan ini menghapus berkas biner `WebView2Loader.dll` yang disematkan dan dependensi `github.com/jchv/go-winloader`. Build tag `native_webview2loader` masih diterima dan tidak lagi menimbulkan galat, tetapi tidak berpengaruh pada build v3, dalam [PR](https://github.com/wailsapp/wails/pull/6031) oleh @Grantmartin2002
- Hapus build tag yang tidak digunakan dan opsi FPS dari panduan API macOS dalam [PR](https://github.com/wailsapp/wails/pull/6097) oleh @leaanthony

## v3.0.0-beta.19 - 2026-09-09

## Ditambahkan

- Batasi akses ke API privat macOS dengan build tag agar penggunaannya harus diaktifkan secara eksplisit — lihat [dokumentasi](https://v3.wails.io/features/browser/integration) dan [dokumentasi](https://v3.wails.io/features/environment/info) dan [dokumentasi](https://v3.wails.io/features/windows/basics) dan [dokumentasi](https://v3.wails.io/features/windows/frameless) dan [dokumentasi](https://v3.wails.io/features/windows/notch-windows) dan [dokumentasi](https://v3.wails.io/features/windows/options) dan [dokumentasi](https://v3.wails.io/guides/build/macos) dan [dokumentasi](https://v3.wails.io/guides/build/private-macos-apis) dan [dokumentasi](https://v3.wails.io/reference/overview) dalam [PR](https://github.com/wailsapp/wails/pull/6087) oleh @leaanthony

## Diperbaiki

- Tolak permintaan runtime berukuran lebih dari 64 MiB dengan HTTP 413 dalam [PR](https://github.com/wailsapp/wails/pull/6091) oleh @leaanthony

## Keamanan

- Perkuat keamanan origin MCP dan akses jarak jauh dengan autentikasi token dalam [PR](https://github.com/wailsapp/wails/pull/6092) oleh @leaanthony

## v3.0.0-beta.18 - 2026-09-08

## Diperbaiki

- Perbaiki kebocoran memori Calloc di Linux dan Darwin dengan menggunakan pointer receiver dalam [PR](https://github.com/wailsapp/wails/pull/6083) oleh @4RH1T3CT0R7

## v3.0.0-beta.17 - 2026-09-06

## Diperbaiki

- Windows: `GetRequest` yang gagal atau nil dalam handler WebResourceRequested tidak lagi menghentikan proses (`log.Fatal` / panic akibat dereferensi nil) — sebagai gantinya, permintaan dibuang dan dicatat dalam log, dalam [PR](https://github.com/wailsapp/wails/pull/6006) oleh @midagedev

## v3.0.0-beta.16 - 2026-08-29

## Diubah

- Minta kata sandi notarisasi di jendela terminal baru dalam [PR](https://github.com/wailsapp/wails/pull/6029) oleh @leaanthony

## Diperbaiki

- Tangani jenis klik systray dengan benar di macOS dalam [PR](https://github.com/wailsapp/wails/pull/5919) oleh @ChewbaccaCookie
- CI menghapus repositori apt Microsoft yang tidak digunakan sebelum melakukan pembaruan dalam [PR](https://github.com/wailsapp/wails/pull/6041) oleh @Grantmartin2002

## v3.0.0-beta.15 - 2026-08-27

## Diperbaiki

- Tingkatkan batas waktu penyematan WebView2 menjadi 60 detik dalam [PR](https://github.com/wailsapp/wails/pull/6043) oleh @Grantmartin2002

## v3.0.0-beta.14 - 2026-08-26

## Diperbaiki

- Beri nama penekanan tombol Control-huruf dengan benar di macOS dalam [PR](https://github.com/wailsapp/wails/pull/6032) oleh @taliesin-ai
- Perbaiki ikon tray ICO dan ikuti tema taskbar di Windows dalam [PR](https://github.com/wailsapp/wails/pull/6016) oleh @nik9play

## v3.0.0-beta.13 - 2026-08-25

## Diperbaiki

- Terus layani pekerjaan thread utama di macOS saat loop modal berjalan dalam [PR](https://github.com/wailsapp/wails/pull/6026) oleh @leaanthony
- Menjadikan penyimpanan aman seluler dapat melaporkan kegagalan dan tetap memblokir akses saat terjadi kegagalan dalam [PR](https://github.com/wailsapp/wails/pull/5923) oleh @mortenolsrud
- Jalankan hook peristiwa aplikasi meskipun tidak ada listener yang terdaftar dalam [PR](https://github.com/wailsapp/wails/pull/5999) oleh @archy-rock3t-cloud
- Perbaiki kesalahan ketik dalam komentar dan dokumentasi yang dilokalkan dalam [PR](https://github.com/wailsapp/wails/pull/6023) oleh @haoku123
- Hapus berkas biner macOS prakompilasi yang di-commit di bawah `v3/examples` dalam [PR](https://github.com/wailsapp/wails/pull/6025) oleh @4RH1T3CT0R7

## v3.0.0-beta.12 - 2026-08-21

## Ditambahkan

- Tambahkan jendela notifikasi notch macOS beserta contoh siklus hidup dan telemetri — lihat [dokumentasi](https://v3.wails.io/features/windows/notch-windows) dalam [PR](https://github.com/wailsapp/wails/pull/6010) oleh @leaanthony
- Menambahkan dukungan jendela NSPanel macOS dengan opsi baru dan integrasi native — lihat [dokumentasi](https://v3.wails.io/features/windows/options) dalam [PR](https://github.com/wailsapp/wails/pull/6008) oleh @leaanthony

## Diperbaiki

- Mencegah Prepare SQLite macet saat dipanggil secara bersamaan dalam [PR](https://github.com/wailsapp/wails/pull/5998) oleh @archy-rock3t-cloud

## v3.0.0-beta.11 - 2026-08-20

## Dihapus

- Menghapus pelacak implementasi yang sudah tidak digunakan dari dokumentasi dalam [PR](https://github.com/wailsapp/wails/pull/6005) oleh @leaanthony

## v3.0.0-beta.10 - 2026-08-19

## Diperbaiki

- Memperbaiki host GTK4 Linux yang mengabaikan argumen peluncuran protokol khusus dan asosiasi berkas dalam [PR](https://github.com/wailsapp/wails/pull/6000) oleh @midagedev
- Menangani baris yang dihapus dan koreksi dari sumber yang sama dengan benar dalam validasi changelog pada [PR](https://github.com/wailsapp/wails/pull/5993) oleh @taliesin-ai

## v3.0.0-beta.9 - 2026-08-16

## Ditambahkan

- Menambahkan server MCP wails3 yang aman untuk pengelolaan proyek dengan bantuan agen dalam [PR](https://github.com/wailsapp/wails/pull/5896) oleh @leaanthony
- Menambahkan dokumentasi untuk model dalam binding — lihat [dokumentasi](https://v3.wails.io/features/bindings/models) dalam [PR](https://github.com/wailsapp/wails/pull/5988) oleh @taliesin-ai
- Menambahkan dukungan instalasi rpm-ostree untuk sistem Linux atomik dalam [PR](https://github.com/wailsapp/wails/pull/5987) oleh @leaanthony
- Menambahkan pembuatan dan publikasi native mingguan untuk grafik riwayat bintang — lihat [dokumentasi](https://v3.wails.io/credits), [dokumentasi](https://v3.wails.io/de/credits), [dokumentasi](https://v3.wails.io/fr/credits), [dokumentasi](https://v3.wails.io/id/credits), [dokumentasi](https://v3.wails.io/ja/credits), [dokumentasi](https://v3.wails.io/ko/credits), [dokumentasi](https://v3.wails.io/pt/credits), [dokumentasi](https://v3.wails.io/ru/credits), [dokumentasi](https://v3.wails.io/zh-cn/credits), dan [dokumentasi](https://v3.wails.io/zh-tw/credits) dalam [PR](https://github.com/wailsapp/wails/pull/5986) oleh @leaanthony
- Menambahkan paket mac khusus Darwin untuk menemukan sumber daya bundle aplikasi — lihat [dokumentasi](https://v3.wails.io/guides/build/macos) dalam [PR](https://github.com/wailsapp/wails/pull/5965) oleh @leaanthony
- Menambahkan halaman showcase Condui dan entri indeks — lihat [dokumentasi](https://v3.wails.io/community/showcase/condui) dan [dokumentasi](https://v3.wails.io/community/showcase) dalam [PR](https://github.com/wailsapp/wails/pull/5962) oleh @mgueregath
- Menambahkan halaman showcase Redis Viewer beserta tangkapan layar dan tautan proyek — lihat [dokumentasi](https://v3.wails.io/community/showcase) dan [dokumentasi](https://v3.wails.io/community/showcase/redisviewer) dalam [PR](https://github.com/wailsapp/wails/pull/5984) oleh @redisviewer

## Diubah

- Memperbarui flag aplikasi GTK menjadi G<em>APPLICATION</em>NON_UNIQUE untuk Linux dalam [PR](https://github.com/wailsapp/wails/pull/5971) oleh @overlordtm
- Mencatat peristiwa jendela yang tidak ditemukan pada tingkat debug, bukan peringatan, dalam [PR](https://github.com/wailsapp/wails/pull/5914) oleh @julianstorer

## Diperbaiki

- Memungkinkan pintasan papan ketik macOS yang terdaftar untuk diprioritaskan daripada webview dalam [PR](https://github.com/wailsapp/wails/pull/5902) oleh @julianstorer
- Memperbaiki tautan bilah sisi dokumentasi yang rusak dalam [PR](https://github.com/wailsapp/wails/pull/5937) oleh @northes
- Membatalkan konteks permintaan aset macOS dan iOS saat WebKit membatalkan tugas skema khusus yang sesuai (#5963)
- Menangani pesan WindowSetFullscreenButtonEnabled dalam [PR](https://github.com/wailsapp/wails/pull/5976) oleh @archy-rock3t-cloud
- Mengimpor Fragment dalam templat preact-ts untuk mengatasi kegagalan build dalam [PR](https://github.com/wailsapp/wails/pull/5979) oleh @haoku123
- Mencegah aplikasi lama berbasis GTK3 yang hanya menjalankan layanan mengalami crash ketika penemuan layar berjalan sebelum jendela atau tampilan aktif tersedia (#5966)
- Memungkinkan proses rilis dengan versi eksplisit dilanjutkan ketika changelog yang belum dirilis kosong (#5977)

## Keamanan

- Memperbarui lockfile nanoid situs web ke versi 3.3.18 yang telah ditambal untuk mengatasi peringatan keamanan dalam [PR](https://github.com/wailsapp/wails/pull/5985) oleh @taliesin-ai

## v3.0.0-beta.8 - 2026-08-12

## Ditambahkan

- Menambahkan pembuatan URL dokumentasi ke entri changelog otomatis dalam [PR](https://github.com/wailsapp/wails/pull/5957) oleh @taliesin-ai
- Menambahkan Streams: stream byte dua arah antara Go dan JavaScript dengan model pemrograman WebSocket dan tanpa soket pendengar. Deklarasikan stream di Go dengan `app.HandleStream(name, handler)` dan hubungkan dari frontend dengan `Stream(name)`, yang mengembalikan objek berbentuk `WebSocket`. Go→JS dikirim melalui satu polling tertahan per jendela lewat server aset, sedangkan JS→Go melalui POST biasa; tidak ada yang mengikat port TCP dan tidak ada yang melewati `evaluateJavaScript`. Dalam build server (`-tags server`), handler yang sama dilayani melalui WebSocket sungguhan sebagai gantinya, sehingga kode aplikasi identik di semua build. oleh @leaanthony
- Memindahkan entri changelog mailbox ke Belum Dirilis dalam [PR](https://github.com/wailsapp/wails/pull/5935) oleh @leaanthony

## Diubah

- Memperbarui pembuatan otomatis bilah sisi dokumentasi dan derivasi tipe penulis blog dalam [PR](https://github.com/wailsapp/wails/pull/5938) oleh @leaanthony

## Diperbaiki

- Inisialisasi WebView2 menggunakan tenggat waktu dan pompa pesan dalam [PR](https://github.com/wailsapp/wails/pull/5952) oleh @leaanthony
- Pengujian cookie WebView2 dilewati di CI kecuali diaktifkan secara eksplisit, dan eksekusinya dikunci ke thread OS saat ini dalam [PR](https://github.com/wailsapp/wails/pull/5951) oleh @leaanthony
- Builder menu Windows memulihkan ID perintah untuk item induk submenu dalam [PR](https://github.com/wailsapp/wails/pull/5944) oleh @gilad-ch
- Menyelaraskan image kompilasi silang resmi dengan baseline dukungan Linux GTK 4.14+ (#5928)
- Mengonfigurasi proyek Xcode iOS untuk mempertahankan flag linker yang diwarisi dan menambahkan -ObjC dalam [PR](https://github.com/wailsapp/wails/pull/5915) oleh @mortenolsrud
- Memperbaiki pergantian koneksi TCP yang berlebihan pada proksi aset `wails3 dev` untuk frontend besar, yang dapat menghabiskan port sementara host dan menyebabkan proses yang tidak terkait gagal dengan `EADDRNOTAVAIL`
- Mengantrekan JavaScript peristiwa per jendela agar dikirim secara berurutan dan menerapkan backpressure dalam [PR](https://github.com/wailsapp/wails/pull/5934) oleh @leaanthony

## Dihapus

- Menghapus pipeline rilis biner desktop: rilis v3 hanya menggunakan tag dan CLI `wails3` diinstal dengan `go install`. Menghapus `release-v3.yml` dan langkah nightly yang menjalankannya dalam [PR](https://github.com/wailsapp/wails/pull/5946) oleh @leaanthony

## v3.0.0-beta.7 - 2026-08-11

## Ditambahkan

- Tambahkan preferensi putar otomatis macOS untuk menonaktifkan persyaratan tindakan pengguna bagi pemutaran media dalam [PR](https://github.com/wailsapp/wails/pull/5512) oleh @Eyalm321
- Pindahkan entri log perubahan mailbox ke Belum Dirilis dalam [PR](https://github.com/wailsapp/wails/pull/5935) oleh @leaanthony

## Diubah

- Animasi zoom macOS menggunakan CADisplayLink atau NSTimer agar kinerjanya lebih halus dalam [PR](https://github.com/wailsapp/wails/pull/5945) oleh @savely-krasovsky

## Diperbaiki

- Konfigurasikan proyek Xcode iOS agar mempertahankan flag linker yang diwariskan dan menambahkan -ObjC dalam [PR](https://github.com/wailsapp/wails/pull/5915) oleh @mortenolsrud
- Perbaiki pergantian koneksi TCP yang berlebihan pada proksi aset `wails3 dev` untuk frontend besar, yang dapat menghabiskan port ephemeral host dan menyebabkan proses yang tidak terkait gagal dengan `EADDRNOTAVAIL`
- Antrekan JavaScript peristiwa per jendela untuk pengiriman berurutan dan backpressure dalam [PR](https://github.com/wailsapp/wails/pull/5934) oleh @leaanthony

### Ditambahkan

- Implementasikan mailbox FIFO asinkron generik untuk pengiriman peristiwa secara berurutan dalam [PR](https://github.com/wailsapp/wails/pull/5851) oleh @savely-krasovsky dan @DevLumuz

## v3.0.0-beta.6 - 2026-08-09

## Ditambahkan

- Implementasikan penyimpanan terbatas di sisi host untuk peristiwa berukuran terlalu besar dan pengiriman JavaScript secara berurutan dalam [PR](https://github.com/wailsapp/wails/pull/5930) oleh @leaanthony
- Implementasikan pantulan Dock macOS untuk kedipan jendela dalam [PR](https://github.com/wailsapp/wails/pull/5921) oleh @julianstorer

## Diperbaiki

- Server aset mempertahankan galat pendeteksi jenis konten dan prefiks yang belum ditulis selama flushing dalam [PR](https://github.com/wailsapp/wails/pull/5931) oleh @leaanthony
- Cegah aplikasi macOS mengalami crash saat menu aplikasi diganti dari callback Wails
- Perbaiki menu native yang tidak terbaca pada Windows 10 1809 / Windows Server 2019 (build 17763). Ekspor uxtheme mode gelap dibatasi untuk build 18334, sehingga aktivasi mode gelap pada tingkat aplikasi tidak pernah dijalankan pada host tersebut: latar belakang menu digambar gelap, tetapi Windows tetap menggambar teks menu menggunakan tema terang, sehingga teks gelap tampil di atas latar belakang gelap. Ordinal tersebut tersedia mulai 17763, sehingga batasannya kini sesuai.
- Perbaiki `w32.GetStockObject` yang memanggil `GetDeviceCaps` alih-alih `GetStockObject`, yang menyebabkannya mengembalikan 0 untuk setiap objek bawaan.
- Tingkatkan penanganan dan pelaporan galat pengunduhan bootstrapper WebView2 dalam [PR](https://github.com/wailsapp/wails/pull/5924) oleh @jannskiee

## v3.0.0-beta.5 - 2026-08-07

## Diperbaiki

- Aktivasi aplikasi macOS kini mematuhi kebijakan aktivasi hanya untuk aplikasi reguler dalam [PR](https://github.com/wailsapp/wails/pull/5897) oleh @julianstorer
- Lindungi jendela GTK yang belum diinisialisasi dalam build Linux dalam [PR](https://github.com/wailsapp/wails/pull/5898) oleh @julianstorer
- Tetapkan warna latar belakang opak secara eksplisit untuk jendela WebKit Linux sebelum URL dimuat dalam [PR](https://github.com/wailsapp/wails/pull/5899) oleh @julianstorer

## v3.0.0-beta.4 - 2026-08-05

## Diubah

- Tugas build Android menggunakan arm64 secara default, dan deploy-emulator memilih arsitektur host dalam [PR](https://github.com/wailsapp/wails/pull/5890) oleh @mortenolsrud

## Diperbaiki

- Pertahankan status zoom jendela macOS selama penyeretan dan kurangi gerakan dalam [PR](https://github.com/wailsapp/wails/pull/5900) oleh @leaanthony
- Perbaiki build mode server Windows dengan menambahkan `!server` ke batasan build `webview_window_windows_nonclient.go`

## v3.0.0-beta.3 - 2026-08-03

## Ditambahkan

- Dokumentasikan penyelesaian verifikasi beta Fase 10 dalam detail implementasi pada [PR](https://github.com/wailsapp/wails/pull/5881) oleh @leaanthony

## Diperbaiki

- Teruskan handle jendela ke API mode gelap Windows dan validasi argumen dalam [PR](https://github.com/wailsapp/wails/pull/5877) oleh @leaanthony
- Pusatkan resolusi status tombol bilah judul macOS untuk jendela tanpa bingkai dalam [PR](https://github.com/wailsapp/wails/pull/5870) oleh @taliesin-ai
- Cegah teks menu native menjadi tidak terbaca ketika aplikasi Windows meminta mode gelap sementara tema aplikasi Windows menggunakan mode terang. Menu kini menggunakan latar belakang native terang yang sesuai hingga Windows dapat merender teks menu gelap.
- Perbaiki menu native yang tidak terbaca pada Windows 10 1809 / Windows Server 2019 (build 17763). Ekspor uxtheme mode gelap hanya diaktifkan mulai build 18334, sehingga aktivasi mode gelap pada tingkat aplikasi tidak pernah dijalankan pada host tersebut: latar belakang menu digambar gelap, tetapi Windows tetap menggambar teks menu menggunakan tema terang, sehingga teks gelap tampil di atas latar belakang gelap. Ordinal tersebut tersedia mulai 17763, sehingga batasannya kini sesuai.

## v3.0.0-beta.2 - 2026-08-02

## Diubah

- Promosikan v3 dari alfa ke beta
- Dokumentasikan default cerdas baki sistem dan perilaku sembunyi otomatis popup, beserta cakupan regresi untuk pemilihan pengendali klik (#5840).
- Updater GitHub secara default mengecualikan aset penginstal Windows dalam [PR](https://github.com/wailsapp/wails/pull/5861) oleh @leaanthony
- Tambahkan dukungan untuk sudut membulat, persegi, dan radius khusus pada jendela macOS tanpa bingkai dalam [PR](https://github.com/wailsapp/wails/pull/5866) oleh @leaanthony

## Diperbaiki

- Laporkan ukuran jendela GTK4 secara langsung dan pancarkan peristiwa perubahan ukuran serta status maksimalkan, minimalkan, dan layar penuh dari surface yang dikonfigurasi (#5830).
- Perbaiki crash WebKit Linux saat mengirim Blob atau FormData dalam permintaan fetch pada [PR](https://github.com/wailsapp/wails/pull/5854) oleh @taliesin-ai
- Shim fetch meneruskan undefined untuk header Blob/FormData yang tidak ada dalam [PR](https://github.com/wailsapp/wails/pull/5865) oleh @leaanthony

## v3.0.0-alpha2.122 - 2026-08-01

## Ditambahkan

## Diubah

- Tambahkan dukungan untuk sudut membulat, persegi, dan radius khusus pada jendela macOS tanpa bingkai dalam [PR](https://github.com/wailsapp/wails/pull/5866) oleh @leaanthony

## Diperbaiki

- Shim fetch meneruskan undefined untuk header Blob/FormData yang tidak ada dalam [PR](https://github.com/wailsapp/wails/pull/5865) oleh @leaanthony

## v3.0.0-alpha2.121 - 2026-07-31

## Ditambahkan

- Menambahkan dukungan pengemasan DMG macOS beserta opsi dan tugas build baru dalam [PR](https://github.com/wailsapp/wails/pull/5857) oleh @leaanthony

## Diubah

- Updater GitHub secara default mengecualikan aset penginstal Windows dalam [PR](https://github.com/wailsapp/wails/pull/5861) oleh @leaanthony

## Diperbaiki

- Memperbaiki crash WebKit Linux saat mengirim Blob atau FormData dalam permintaan fetch dalam [PR](https://github.com/wailsapp/wails/pull/5854) oleh @taliesin-ai

## v3.0.0-alpha2.120 - 2026-07-31

## Ditambahkan

- Mengimplementasikan tindakan klik ganda pada bilah judul macOS untuk memaksimalkan atau meminimalkan jendela dalam [PR](https://github.com/wailsapp/wails/pull/5853) oleh @taliesin-ai

## Diubah

- Memperbarui tutorial layanan QR agar menggunakan NewServiceWithOptions dan menambahkan spasi dalam [PR](https://github.com/wailsapp/wails/pull/5849) oleh @jeongkyu

## Diperbaiki

- Menjaga WKWebView tetap responsif selama zoom di macOS dalam [PR](https://github.com/wailsapp/wails/pull/5856) oleh @leaanthony
- Memperbaiki kueri ukuran jendela GTK4 dan memancarkan peristiwa perubahan ukuran serta status maksimalkan, minimalkan, dan layar penuh dari `GdkSurface` yang dikonfigurasi.

## v3.0.0-alpha2.119 - 2026-07-27

## Diperbaiki

- Memperbarui dokumentasi dalam berbagai bahasa agar menyertakan diagram arsitektur dalam [PR](https://github.com/wailsapp/wails/pull/5833) oleh @taliesin-ai

## v3.0.0-alpha2.118 - 2026-07-26

## Ditambahkan

- Menyediakan jalur default untuk masukan dan keluaran pembuatan ikon dalam [PR](https://github.com/wailsapp/wails/pull/5825) oleh @taliesin-ai
- Menambahkan modul entri sumber ke sideEffects dalam package.json paket runtime dalam [PR](https://github.com/wailsapp/wails/pull/5797) oleh @savely-krasovsky
- Menambahkan bagian lisensi dan asal-usul ke panduan kontribusi dalam [PR](https://github.com/wailsapp/wails/pull/5816) oleh @taliesin-ai

## Diperbaiki

- Menerapkan CSS tanpa bingkai GTK4 dengan cakupan terbatas untuk menghapus radius batas dalam [PR](https://github.com/wailsapp/wails/pull/5800) oleh @savely-krasovsky
- Menangani kegagalan posisi kursor dengan baik di Windows untuk menu pop-up dan enumerasi layar dalam [PR](https://github.com/wailsapp/wails/pull/5789) oleh @wayneforrest
- Dialog buka berkas macOS memfilter ekstensi dengan benar dan memvalidasi berkas yang diizinkan berdasarkan sufiks dalam [PR](https://github.com/wailsapp/wails/pull/5678) oleh @phergul
- Melindungi inisialisasi mode gelap Windows dari pemanggilan API nil dalam [PR](https://github.com/wailsapp/wails/pull/5793) oleh @roachadam
- Memperbaiki kegagalan build 32-bit pada updater: konstanta `maxArchiveTotalSize` (2 GiB) melampaui kapasitas `int` platform ketika diteruskan ke `fmt.Errorf` pada `GOARCH=386`. Kini konstanta tersebut secara eksplisit diberi tipe `int64`.
- Memperbaiki panic akibat pointer nil saat startup ketika sebuah jendela menggunakan bilah judul Gelap (atau gelap-sistem) pada build Windows yang tidak memuat API uxtheme mode gelap, seperti Windows 10 1809 / Windows Server 2019 (build 17763). Pemanggilan `AllowDarkModeForWindow` dalam penyiapan tema jendela kini dilindungi dari nilai nil, sesuai dengan perlindungan yang sudah digunakan dalam `w32.SetMenuTheme`.

## v3.0.0-alpha2.117 - 2026-07-08

## Ditambahkan

- Mengimplementasikan logika hit-test khusus untuk area nonklien di Windows dalam [PR](https://github.com/wailsapp/wails/pull/5462) oleh @savely-krasovsky

## Diubah

- Mengonfigurasi deteksi skala monitor WebView2 berdasarkan UseVisualHosting dalam [PR](https://github.com/wailsapp/wails/pull/5761) oleh @wayneforrest

## v3.0.0-alpha2.116 - 2026-07-07

## Ditambahkan

- Memperbarui dokumentasi FAQ agar berfokus pada fitur dan panduan Wails v3 dalam [PR](https://github.com/wailsapp/wails/pull/5763) oleh @taliesin-ai

## v3.0.0-alpha2.115 - 2026-07-06

## Diperbaiki

- Memperbaiki `Menu.Update()` yang tidak membangun ulang menu native pada GTK4 Linux (#5659, didiagnosis dan diperbaiki secara independen oleh @puneetdixit200 dalam #5539)
- Memperbaiki crash saat mengenumerasi layar macOS ketika tampilan berubah dengan menyalin string id/nama layar dan mengambil snapshot jumlahnya (#5565, didiagnosis dan diperbaiki secara independen oleh @x-haose dalam #5584)
- Memperbaiki crash di Windows ketika `WM_ERASEBKGND` menggambar latar belakang berwarna solid selama transisi minimalkan/pulihkan saat `GetClientRect` mengembalikan nil (perlindungan dilaporkan oleh @sinspired dalam #5636)
- Memperbaiki kesalahan binding frontend yang selalu diurai sebagai teks oleh @mbaklor dalam #5690
- Memperbaiki kegagalan build Windows saat menggunakan tag build `server`, yang disebabkan oleh berkas GUI Windows yang tidak memiliki batasan build `!server` yang sudah dimiliki oleh berkas ekuivalen macOS dan Linux (#5680)

## v3.0.0-alpha2.114 - 2026-07-05

## Ditambahkan

- Mengimplementasikan protokol Update Manifest dan penyedia endpoint dalam [PR](https://github.com/wailsapp/wails/pull/5720) oleh @taliesin-ai

## Diubah

- Menggabungkan binding `webview2` ke dalam modul v3 sebagai `v3/internal/webview2`, serta menghapus modul mandiri, alur kerja rilis/sinkronisasi nightly, dan mekanisme perubahan versi go.mod miliknya (v3 adalah satu-satunya penggunanya) dalam [PR](https://github.com/wailsapp/wails/pull/5711) oleh @taliesin-ai

## Diperbaiki

- Memindahkan perbaikan deteksi skala monitor WebView2 dan sinkronisasi ulang host saat perubahan DPI ke bagian Belum Dirilis dalam [PR](https://github.com/wailsapp/wails/pull/5750) oleh @taliesin-ai
- Memperbarui marshaling COM WebView2 untuk parameter float64 dan BOOL dalam [PR](https://github.com/wailsapp/wails/pull/5741) oleh @wayneforrest
- Mencegah panic dan dereferensi nil saat memperbarui dan menghancurkan ikon baki sistem Windows dalam [PR](https://github.com/wailsapp/wails/pull/5703) oleh @wayneforrest
- Memperbaiki jendela tersembunyi yang tidak disembunyikan kembali dengan benar di Windows dalam [PR](https://github.com/wailsapp/wails/pull/5743) oleh @wayneforrest
- Menyinkronkan visibilitas pengontrol WebView2 dengan tindakan meminimalkan, memaksimalkan, dan memulihkan jendela dalam [PR](https://github.com/wailsapp/wails/pull/5742) oleh @wayneforrest

### Diperbaiki

- Mengaktifkan kembali deteksi skala monitor WebView2 dan membatasi sinkronisasi ulang host agar hanya dilakukan saat DPI berubah dalam [PR](https://github.com/wailsapp/wails/pull/5734) oleh @taliesin-ai, berdasarkan perbaikan yang divalidasi oleh @randalmurphal, dengan verifikasi akar penyebab oleh @eleclin dan pengujian perangkat keras oleh @qq540491950

## v3.0.0-alpha2.113 - 2026-07-04

## Ditambahkan

- Tambahkan peringatan saat membangun AAB rilis tanpa menetapkan `ANDROID_KEYSTORE_FILE` (Google Play menolak bundle yang ditandatangani dengan kunci debug), serta dokumentasikan pengemasan dan penandatanganan App Bundle dalam [PR](https://github.com/wailsapp/wails/pull/5730) oleh @taliesin-ai
- Tambahkan dokumentasi Why Wails dalam beberapa bahasa melalui [PR](https://github.com/wailsapp/wails/pull/5739) oleh @taliesin-ai
- Tambahkan dukungan untuk memetakan Go time.Time ke JS Date atau string dalam binding melalui [PR](https://github.com/wailsapp/wails/pull/5398) oleh @fbbdev
- Tambahkan tugas pengemasan Android App Bundle (AAB) (`bundle`, `bundle:fat`, `assemble:aab`, `assemble:aab:release`) untuk pengiriman ke Play Store—tugas APK tetap tersedia untuk pengujian lokal/emulator—melalui [PR](https://github.com/wailsapp/wails/pull/5728) oleh @mortenolsrud (memperbaiki [#5726](https://github.com/wailsapp/wails/issues/5726))
- Tambahkan target tugas untuk perangkat fisik Android dan pemulihan izin kamera/lokasi melalui [PR](https://github.com/wailsapp/wails/pull/5735) oleh @taliesin-ai

## Diubah

- Tingkatkan `webview2` ke v1.0.28 ([catatan rilis](https://github.com/wailsapp/wails/releases/tag/webview2%2Fv1.0.28)).
- Tingkatkan `compileSdk`/`targetSdk` pada templat Android dari 34 ke 35, yang diwajibkan Google Play untuk pengiriman aplikasi baru, melalui [PR](https://github.com/wailsapp/wails/pull/5730) oleh @taliesin-ai

## Diperbaiki

- Perbaiki mask avatar yang telah diproses dalam sponsorkit melalui [PR](https://github.com/wailsapp/wails/pull/5745) oleh @leaanthony
- Perbaiki pembuatan otomatis AVD Android yang memilih image sistem atau versi cmdline-tools yang salah akibat pengurutan versi secara leksikografis melalui [PR](https://github.com/wailsapp/wails/pull/5730) oleh @taliesin-ai
- Perbaiki wizard penyiapan yang menyarankan versi Android NDK usang (sekarang 26.3.11579264, sesuai dengan persyaratan yang didokumentasikan) melalui [PR](https://github.com/wailsapp/wails/pull/5730) oleh @taliesin-ai
- Perbarui dokumentasi bahasa Prancis untuk SvelteKit dan opsi melalui [PR](https://github.com/wailsapp/wails/pull/5744) oleh @leaanthony
- Perbaiki SIGSEGV saat enumerasi layar macOS selama perubahan tampilan melalui [PR](https://github.com/wailsapp/wails/pull/5516) oleh @flofreud

## v3.0.0-alpha2.112 - 2026-07-03

## Ditambahkan

- Tambahkan generator SVG kontributor berbasis Go dan perbarui halaman kredit dokumentasi/situs web melalui [PR](https://github.com/wailsapp/wails/pull/5724) oleh @taliesin-ai
- Tambahkan dukungan untuk memetakan Go time.Time ke JS Date atau string dalam binding melalui [PR](https://github.com/wailsapp/wails/pull/5398) oleh @fbbdev

## Diubah

- Ganti pipeline gambar sponsor berbasis Node dengan generator Go melalui [PR](https://github.com/wailsapp/wails/pull/5719) oleh @taliesin-ai

## Diperbaiki

- Perbaiki skrip instalasi dependensi aset build Android melalui [PR](https://github.com/wailsapp/wails/pull/5729) oleh @taliesin-ai
- Tolak karakter kontrol U+0085 (NEXT LINE) dalam `ValidateAndSanitizeURL`, sehingga cakupan karakter spasi pada validator URL menjadi lengkap
- Hitung ulang frame DWM saat DPI berubah untuk jendela tanpa bingkai melalui [PR](https://github.com/wailsapp/wails/pull/4785) oleh @leaanthony
- Perbaiki kegagalan deteksi dropzone DnD pada penskalaan selain 100% di Windows melalui [PR](https://github.com/wailsapp/wails/pull/4632) oleh @yulesxoxo
- Tambahkan pengelolaan memori Objective-C eksplisit untuk objek Cocoa di seluruh dialog, menu, tray, dan notifikasi Darwin melalui [PR](https://github.com/wailsapp/wails/pull/5714) oleh @taliesin-ai
- Perbaiki bug backend CGO Linux dan masalah system tray melalui [PR](https://github.com/wailsapp/wails/pull/5718) oleh @taliesin-ai

## v3.0.0-alpha2.111 - 2026-07-01

## Ditambahkan

- Tambahkan HappyTools ke galeri komunitas melalui [PR](https://github.com/wailsapp/wails/pull/5061) oleh @Aliuyanfeng
- Tambahkan dukungan bahasa Indonesia dan dokumentasi yang komprehensif melalui [PR](https://github.com/wailsapp/wails/pull/5643) oleh @triadmoko
- Tambahkan opsi DisableMenu ke WindowsWindow melalui [PR](https://github.com/wailsapp/wails/pull/4813) oleh @leaanthony

## Diubah

- Perbarui templat Taskfile dan CLI agar meneruskan tugas build/package dengan GOOS dan ARCH melalui [PR](https://github.com/wailsapp/wails/pull/5617) oleh @leaanthony

## Diperbaiki

- Perbaiki masalah pengelompokan tab jendela macOS melalui [PR](https://github.com/wailsapp/wails/pull/5708) oleh @taliesin-ai

## Dihapus

- Hapus file MDX terjemahan bahasa Jerman dari bagian kontribusi, fitur, dan panduan melalui [PR](https://github.com/wailsapp/wails/pull/5702) oleh @taliesin-ai

## v3.0.0-alpha2.110 - 2026-06-30

## Ditambahkan

- Implementasikan muat ulang dan muat ulang paksa WebView di macOS, serta tambahkan pemulihan setelah proses WebContent dihentikan melalui [PR](https://github.com/wailsapp/wails/pull/5129) oleh @wayneforrest
- Tambahkan dokumentasi bahasa Jerman yang komprehensif untuk kontribusi, fitur, dan panduan melalui [PR](https://github.com/wailsapp/wails/pull/5396) oleh @leaanthony
- Tingkatkan notifikasi dengan suara, lampiran, penjadwalan, dan API pembaruan melalui [PR](https://github.com/wailsapp/wails/pull/5333) oleh @popaprozac

## Diperbaiki

- Hitung ulang frame DWM saat DPI berubah untuk jendela tanpa bingkai melalui [PR](https://github.com/wailsapp/wails/pull/4785) oleh @leaanthony
- Perbaiki kegagalan deteksi dropzone DnD pada penskalaan selain 100% di Windows melalui [PR](https://github.com/wailsapp/wails/pull/4632) oleh @yulesxoxo

## v3.0.0-alpha2.109 - 2026-06-29

## Ditambahkan

- Tambahkan contoh kode ke dokumentasi EventsEmit melalui [PR](https://github.com/wailsapp/wails/pull/5026) oleh @iamhabbeboy
- Tambahkan opsi hosting visual WebView2 di Windows melalui [PR](https://github.com/wailsapp/wails/pull/5380) oleh @MerIijn
- Tambahkan Klustr ke dokumentasi galeri komunitas melalui [PR](https://github.com/wailsapp/wails/pull/5536) oleh @SametKUM
- Tambahkan Kira ke galeri komunitas beserta halaman baru dan entri changelog melalui [PR](https://github.com/wailsapp/wails/pull/5685) oleh @thiennguyen93
- Tambahkan bagian umpan balik ke panduan layanan MCP melalui [PR](https://github.com/wailsapp/wails/pull/5694) oleh @taliesin-ai

## Diubah

- Mode server kini memiliki build produksi kelas satu, yang konsisten dengan tugas build desktop (#5693). `task build:server` secara default membuat biner produksi (`-tags server,production`, `-trimpath`, simbol dihapus), serta menerima `DEV=true` (server pengembangan), `OBFUSCATED=true` (garble), dan `EXTRA_TAGS`. `task run:server` menjalankan server pengembangan. `Dockerfile.server` / `task build:docker` terlebih dahulu membuat server produksi (`-tags server,production`) dan frontend produksi; image secara default menggunakan build statis Go murni di distroless/static, dengan `CGO_ENABLED`, `GO_IMAGE`, dan `RUNTIME_IMAGE` yang diekspos sebagai argumen build yang dapat diganti untuk aplikasi CGO.

## Diperbaiki

- Mencegah crash saat menutup jendela yang memiliki panggilan asinkron tertunda dalam [PR](https://github.com/wailsapp/wails/pull/4435) oleh @leaanthony
- Mencegah aktivasi jendela saat membuka aplikasi tersembunyi di Windows dalam [PR](https://github.com/wailsapp/wails/pull/5249) oleh @leaanthony
- Memastikan metadata permintaan WebKit, penyelesaian respons, dan penanganan stream isi dijalankan di thread utama GTK dalam [PR](https://github.com/wailsapp/wails/pull/5668) oleh @taliesin-ai
- Memperbaiki `Menu.Update()` yang tidak membangun ulang menu native di GTK4 Linux (#5659, didiagnosis dan diperbaiki secara terpisah oleh @puneetdixit200 dalam #5539)
- Memperbaiki crash saat mengenumerasi layar macOS ketika tampilan berubah dengan menyalin string ID/nama layar dan mengambil snapshot jumlahnya (#5565, didiagnosis dan diperbaiki secara terpisah oleh @x-haose dalam #5584)
- Memperbaiki konten WebView2 yang mengecil lalu menghilang setelah jendela diseret melintasi monitor dengan DPI berbeda di Windows dengan menetapkan ulang batas pengontrol dalam handler `WM_DPICHANGED`, mengikuti sinkronisasi ulang DPI saat keluar dari keadaan diminimalkan (#5677)

## v3.0.0-alpha2.108 - 2026-06-28

## Ditambahkan

- Menambahkan pintasan papan ketik global (seluruh sistem) melalui `app.GlobalShortcut` (`Register`, `Unregister`, `UnregisterAll`, `IsRegistered`, `GetAll`). Pintasan tetap terpicu meskipun aplikasi tidak memiliki fokus. Diimplementasikan secara native untuk setiap platform tanpa dependensi pihak ketiga: hot key Carbon di macOS, `RegisterHotKey` di Windows, `XGrabKey` di X11, dan antarmuka pintasan global XDG Desktop Portal di Wayland.
- Menambahkan server MCP bawaan: server Model Context Protocol yang dimulai secara otomatis ketika aplikasi dibuat dengan tag `mcp`, sehingga agen LLM dapat menguji dan mengendalikan aplikasi Wails yang sedang berjalan—mengendalikan jendela, memeriksa DOM, mengevaluasi JavaScript, memanggil metode terikat, menangani peristiwa, serta menyimulasikan input tetikus/papan ketik yang ditampilkan dengan kursor animasi di layar. Tidak memerlukan kode pengguna: tag `mcp` ditambahkan secara otomatis oleh `wails3 build`/`wails3 dev` ketika `WAILS_MCP=1` ditetapkan. Seluruh konfigurasi dilakukan melalui variabel lingkungan (`WAILS_MCP_HOST`, `WAILS_MCP_PORT`, `WAILS_MCP_TIMEOUT`, `WAILS_MCP_HIDE_CURSOR`).

## Diperbaiki

- Memperbaiki `Menu.Update()` yang tidak membangun ulang menu native di GTK4 Linux (#5659, didiagnosis dan diperbaiki secara terpisah oleh @puneetdixit200 dalam #5539)
- Memperbaiki crash saat mengenumerasi layar macOS ketika tampilan berubah dengan menyalin string ID/nama layar dan mengambil snapshot jumlahnya (#5565, didiagnosis dan diperbaiki secara terpisah oleh @x-haose dalam #5584)

## v3.0.0-alpha2.107 - 2026-06-27

## Ditambahkan

- Menambahkan dokumentasi eksperimental Wake beserta navigasi bilah sisi dalam [PR](https://github.com/wailsapp/wails/pull/5613) oleh @leaanthony

## v3.0.0-alpha2.106 - 2026-06-24

## Diubah

- Menaikkan versi `webview2` ke v1.0.27.
  - ci(webview2): memperbaiki build rilis (kompilasi silang Windows + melengkapi go.sum) (#5671)\

  **Diff lengkap:** https://github.com/wailsapp/wails/compare/webview2/v1.0.26...webview2/v1.0.27

- Menghapus go vet dari kompilasi silang alur kerja rilis webview2 dalam [PR](https://github.com/wailsapp/wails/pull/5672) oleh @taliesin-ai
- Memperbarui model OpenRouter untuk auto-changelog ke google/gemini-2.5-flash-lite dalam [PR](https://github.com/wailsapp/wails/pull/5670) oleh @taliesin-ai
- Menaikkan versi `webview2` ke v1.0.26.

### Perbaikan

- **Memulihkan diri dari error COM runtime sementara alih-alih keluar** (#5658, #5580). Sebelumnya, `Chromium.errorCallback` memanggil `os.Exit(1)` untuk *setiap* error COM, sehingga gangguan sementara yang dapat dipulihkan setelah startup menghentikan seluruh aplikasi. Jalur runtime (`Resize`/`GetClientRect`, `Navigate`/`NavigateToString`, `Init`, `MessageReceived`, `PutZoomFactor`, `OpenDevToolsWindow`) kini mencatat error dan melakukan pemulihan. Secara khusus, pesan web yang rusak/tidak tepercaya dalam `MessageReceived` kini dibuang alih-alih menghentikan proses. Ini mengatasi kelas crash saat melintasi monitor dengan DPI berbeda (#5544, #5650). Jalur pembuatan lingkungan/pengontrol tetap bersifat fatal.\

**Diff lengkap:** https://github.com/wailsapp/wails/compare/webview2/v1.0.25...webview2/v1.0.26

## Diperbaiki

- Memperbaiki alur kerja release-webview2 agar menangani file go.sum dengan benar dalam [PR](https://github.com/wailsapp/wails/pull/5671) oleh @taliesin-ai
- Memperbaiki pembaruan menu GTK4 Linux dengan mengosongkan dan membangun ulang menu native dalam [PR](https://github.com/wailsapp/wails/pull/5659) oleh @taliesin-ai

## v3.0.0-alpha2.105 - 2026-06-21

## Ditambahkan

- Menambahkan `application.System` untuk mendeteksi platform saat runtime dari kode bersama: `System.IsMobile()` (iOS/Android), `System.IsDesktop()` (macOS/Windows/Linux), `System.IsServer()` (tag build `server`), dan `System.IsPlatform(application.PlatformMacOS|PlatformWindows|PlatformLinux|PlatformIOS|PlatformAndroid|PlatformServer)` untuk menguji satu target secara langsung. Fitur ini dapat dikompilasi pada setiap target sehingga Anda dapat membuat percabangan tanpa tag build. Helper frontend yang sesuai (`System.IsMobile/IsDesktop/IsIOS/IsAndroid/...`) tersedia di `@wailsio/runtime`
- Menambahkan panduan "Menggunakan Framework Frontend Lain" yang menunjukkan cara menempatkan proyek Vite Anda sendiri ke dalam `frontend/` (mencakup Solid, Preact, Lit, SvelteKit, Qwik, Angular, dan lainnya)
- Wizard `wails3 setup` kini memeriksa toolchain seluler (iOS/Android)—Xcode dan runtime iOS Simulator, JDK, Android SDK/NDK, serta emulator—dengan pemasangan sekali klik dan perbaikan konfigurasi shell yang dapat disalin jika tersedia
- Proyek yang dihasilkan menyertakan `frontend/.npmrc` yang menetapkan `minimum-release-age` selama 7 hari untuk mengurangi paparan terhadap paket yang baru diterbitkan (dan mungkin telah disusupi) (dipatuhi oleh pnpm dan bun; diabaikan tanpa dampak oleh npm)

## Diubah

- Mendesain ulang semua templat awal bawaan dengan tampilan hero pegunungan neon baru (web, iOS, dan Android)
- **TypeScript kini menjadi default untuk templat awal dan menggunakan nama templat tanpa akhiran.** `wails3 init` (tanpa `-t`) membuat kerangka proyek TypeScript; `-t vanilla`, `-t react`, `-t vue`, dan `-t svelte` menggunakan TypeScript, dengan varian JavaScript di `-t vanilla-js`, `-t react-js`, `-t vue-js`, dan `-t svelte-js`. Templat bawaan mendeklarasikan bahasanya melalui `typescript:` dalam `template.yaml`; templat komunitas yang menggunakan akhiran `-ts` tetap berfungsi sebagai fallback
- Mendesain ulang wizard `wails3 setup` dengan tema neon "Wails digital" (efek frosted glass yang hidup di atas latar pegunungan)

## Diperbaiki

- Memperbaiki crash di Windows saat memulihkan aplikasi yang diminimalkan cukup lama hingga WebView2 ditangguhkan atau proses render/GPU-nya didaur ulang. Sinkronisasi ulang DPI saat meminimalkan/memulihkan (#5544) kini hanya menyentuh pengontrol WebView2 jika DPI jendela benar-benar berubah, sehingga menghindari pemanggilan COM fatal ke pengontrol yang ditangguhkan pada pemulihan dengan DPI yang sama, yang merupakan kasus umum (#5605)
- Memperbaiki crash native `SIGABRT`/`SIGSEGV` yang berulang (biasanya di dalam `g_object_unref` selama loop utama GTK) pada aplikasi Linux yang berjalan lama dan sering memuat aset/media. Server aset menyelesaikan `WebKitURISchemeRequest` dari goroutine pekerja sehingga fungsi WebKit2GTK yang tidak aman terhadap thread dipanggil di luar thread utama GTK; penyelesaian (`webkit_uri_scheme_request_finish_with_response`/`finish_error`) kini dijalankan di thread utama. Menyempurnakan perbaikan parsial dalam #5566. Memengaruhi build GTK3 maupun GTK4/WebKitGTK 6.0 (#5631, #5557)
- Memperbaiki `fatal error: invalid pointer found on stack` yang terjadi sesekali dalam `setupSignalHandlers` di Linux/GTK3. ID jendela yang diteruskan sebagai `user_data` untuk sinyal disimpan dalam variabel lokal Go `unsafe.Pointer`, sehingga garbage collector membatalkan proses ketika memindai nilai (non-pointer) tersebut selama penyalinan stack. ID tersebut kini dipertahankan sebagai tipe integer (`uintptr_t`) di sisi Go, dengan menerapkan kembali ke jalur GTK3 lama perbaikan yang sama dengan yang diterapkan #4958 pada jalur GTK4 (yang mengubah fungsi sinyal C menjadi `uintptr_t` untuk mengatasi galat `-race`/checkptr) (#5631)

## Dihapus

- Menghapus templat pemula `react-swc`, `preact`, `lit`, `solid`, `qwik`, dan `sveltekit` (beserta varian `-ts`-nya). Kumpulan bawaan yang kini didukung adalah `vanilla`, `react`, `vue`, dan `svelte`—masing-masing menggunakan TypeScript secara default, dengan varian JavaScript `-js`. Framework lain tetap dapat digunakan dengan [membawa frontend Anda sendiri](https://v3.wails.io/guides/dev/frontend-frameworks) atau melalui templat khusus

## v3.0.0-alpha2.104 - 2026-06-18

## Diperbaiki

- Memperbaiki crash iOS (SIGABRT) ketika metode layanan Go yang diikat mengembalikan string kosong. Penulis respons aset iOS memeriksa pointer isi dengan `buf != nil` alih-alih panjangnya, sehingga isi dengan panjang nol menyebabkan `&buf[0]` mengalami panic; kini pemeriksaan didasarkan pada panjang, selaras dengan penulis untuk desktop

## v3.0.0-alpha2.103 - 2026-06-15

## Diubah

- Memindahkan fitur native iOS dan Android ke pengelola platform: panggil fitur tersebut melalui `application.IOS.*` dan `application.Android.*` (misalnya `application.IOS.Haptic("medium")`, `application.Android.Share(payload)`), bukan melalui fungsi bebas `application.IOS*`/`application.Android*` yang lama (#5602)
- Mengganti nama event bridge seluler: event lintas platform kini menggunakan prefiks `common:*` (misalnya `common:haptic`, `common:location`), sedangkan event khusus platform menggunakan `ios:*` / `android:*` (misalnya `ios:backgroundTask`, `android:foregroundService`); prefiks `native:*` tidak lagi digunakan (#5602)

## v3.0.0-alpha.102 - 2026-06-14

## Ditambahkan

- Menambahkan wizard eksperimental `wails3 setup` untuk penyiapan proyek interaktif dan pemeriksaan dependensi
- Menambahkan flag `--json` ke `wails3 doctor` untuk keluaran yang dapat dibaca mesin
- Menambahkan bagian status penandatanganan ke perintah `wails3 doctor`

## Diperbaiki

- Memperbaiki deteksi npm di Linux agar memeriksa PATH selain pengelola paket

## v3.0.0-alpha.101 - 2026-06-13

## Ditambahkan

- iOS: dialog pesan native (UIAlertController) serta dialog untuk membuka file/beberapa file/direktori (UIDocumentPickerViewController); dialog penyimpanan mengembalikan galat eksplisit
- iOS: dukungan papan klip melalui UIPasteboard
- iOS: metrik layar sebenarnya melalui UIScreen (titik, piksel, skala, area kerja safe area)
- iOS: build perangkat (`IOS_PLATFORM=device`), dukungan identitas penandatanganan kode/profil provisi/entitlement, pemaketan `.ipa`, dan `deploy-device` melalui devicectl
- iOS: versi minimum iOS yang dapat dikonfigurasi (`ios.minIOSVersion` dalam build/config.yml)
- iOS: `wails3 doctor` melaporkan ketersediaan Xcode dan iOS SDK di macOS
- iOS: event sistem—baterai, jaringan, tema, penguncian layar, dan memori rendah diekspos sebagai event aplikasi `events.IOS.*` dan `events.Common.*` yang netral terhadap platform
- iOS: bridge fitur seluler native (`application.IOS*` yang diekspor)—lembar berbagi, membuka URL, menjaga perangkat tetap aktif, senter, inset safe area, kecerahan, informasi aplikasi, penguncian orientasi, bilah status, biometrik (Face ID/Touch ID), notifikasi lokal, dan penyimpanan aman Keychain
- iOS: sensor & perangkat keras — haptik, pengambilan geolokasi satu kali per permintaan, akselerometer, proksimitas, teks ke ucapan, informasi penyimpanan, status daya/baterai, status jaringan, inset papan ketik, dan deteksi tangkapan layar
- iOS: dokumentasi (IOS.md dan panduan situs dokumentasi)
- Android: dialog pesan native (AlertDialog) serta dialog untuk membuka file/beberapa file (Storage Access Framework, diimpor sebagai salinan cache); dialog pembukaan direktori dan penyimpanan mengembalikan galat eksplisit
- Android: dukungan papan klip melalui ClipboardManager
- Android: metrik layar sebenarnya melalui WindowMetrics/DisplayMetrics (dp, piksel, skala, area kerja bilah sistem)
- Android: metode runtime untuk haptik (`Android.Haptics.Vibrate`), informasi perangkat (`Android.Device.Info`), dan toast (`Android.Toast.Show`)
- Android: event siklus hidup bertipe (`events.Android.*`, dihasilkan dari events.txt), dengan `ActivityCreated` dipetakan ke `Common.ApplicationStarted`
- Android: pipeline build menghasilkan APK debug dan rilis yang dapat diinstal (`android:run`, `android:package`, `android:package:fat`); penandatanganan rilis secara default menggunakan keystore debug, atau keystore sebenarnya melalui variabel lingkungan `ANDROID_KEYSTORE_*`
- Android: `wails3 doctor` melaporkan Android SDK, NDK, dan JDK
- Android: event sistem—baterai, jaringan, tema, penguncian layar, dan memori rendah diekspos sebagai event aplikasi `events.Android.*` dan `events.Common.*` yang netral terhadap platform
- Android: bridge fitur seluler native (`application.Android*` yang diekspor)—berbagi, membuka URL, menjaga perangkat tetap aktif, senter, inset safe area, kecerahan, informasi aplikasi, penguncian orientasi, bilah status, biometrik (BiometricPrompt), notifikasi lokal, dan penyimpanan aman EncryptedSharedPreferences
- Android: sensor & perangkat keras — haptik, pengambilan geolokasi satu kali per permintaan, akselerometer, proksimitas, teks ke ucapan, informasi penyimpanan, status daya/baterai, status jaringan, inset papan ketik, dan pemblokiran tangkapan layar dengan FLAG_SECURE
- Android: dokumentasi (ANDROID.md dan panduan situs dokumentasi)
- Contoh: kitchen sink `mobile` kini memiliki tab Seluler dan Perangkat Keras yang mendemonstrasikan bridge fitur native di iOS dan Android (tab berbentuk pil berlanjut ke beberapa baris)
- Seluler: baterai—akselerometer, sensor proksimitas, senter, dan jam berkala pada contoh dijeda ketika aplikasi berjalan di latar belakang lalu dipulihkan saat aplikasi kembali aktif (Android mempertahankan proses tetap berjalan di latar belakang, sedangkan senter merupakan status perangkat keras yang tetap dipertahankan di iOS), dan penerima event sistem Android hanya didaftarkan selama aplikasi berada di latar depan
- iOS: pengambilan gambar dengan kamera — `application.IOSCapturePhoto`/`IOSCaptureVideo` (UIImagePickerController → peristiwa `native:capture` dengan gambar mini base64)
- iOS: eksekusi latar belakang — `application.IOSBeginBackgroundTask`/`IOSEndBackgroundTask` (periode waktu untuk tugas latar belakang UIApplication) dan `ios.backgroundModes` yang dapat dikonfigurasi (build/config.yml), yang menyisipkan `UIBackgroundModes` melalui templat ke dalam Info.plist yang dihasilkan
- Android: pengambilan gambar dengan kamera — `application.AndroidCapturePhoto`/`AndroidCaptureVideo` (kamera sistem melalui FileProvider → peristiwa `native:capture`)
- Android: layanan latar depan — `application.AndroidStartForegroundService`/`AndroidStopForegroundService` (`WailsForegroundService` dengan notifikasi berkelanjutan menjaga proses tetap berjalan untuk pekerjaan latar belakang berdurasi panjang)
- Contoh: tab Kamera yang mendemonstrasikan pengambilan foto/video dan eksekusi latar belakang (layanan latar depan di Android, periode waktu untuk tugas latar belakang di iOS)

## Diperbaiki

- Memperbaiki `getUserMedia` yang selalu gagal dengan `NotAllowedError` di Linux: WebKitGTK menolak permintaan izin yang tidak ditangani, dan sinyal `permission-request` belum terhubung. Kamera/mikrofon kini ditangani melalui peta lintas platform `WebviewWindowOptions.Permissions` baru (`map[PermissionType]Permission`), yang diterapkan di Linux (WebKitGTK) maupun Windows (WebView2). Di Linux, yang tidak memiliki dialog native, kamera/mikrofon diizinkan secara default (memulihkan `getUserMedia`) dan dapat dinonaktifkan dengan `PermissionDeny` (#5552)
- iOS: `GOOS=ios` dapat dikompilasi kembali (`events.IOS` yang diekspor, stub nama metode seluler), dan build bertag produksi dapat dikompilasi (perbaikan tag build di pkg/application dan beberapa layanan)
- iOS: peristiwa Go→JS dan ExecJS kini berfungsi — halaman tidak lagi dimuat dua kali saat dimulai dan handshake `wails:runtime:ready` tidak lagi dapat hilang
- iOS: `ApplicationDidFinishLaunching`/`ApplicationStarted` tidak lagi mengalami kondisi balapan dengan proses mulai aplikasi; penundaan awal tetap selama 2 detik telah dihapus
- iOS: memperbaiki kebocoran string C pada setiap eksekusi JavaScript Go→JS
- iOS: `hasListeners` kini mencerminkan pendaftaran listener yang sebenarnya
- iOS: pencatatan log debug framework tidak disertakan saat mengompilasi build produksi
- Android: `GOOS=android` dapat dikompilasi kembali — mendefinisikan `events.Android`, menghapus larik listener `events_android.go` yang melampaui batas, menambahkan stub nama metode seluler, dan mencegah berkas Linux desktop (`linux_cgo.*`, `events_linux.*`, `environment_linux.go`) ikut masuk ke build Android
- Android: binding JS→Go kini berfungsi — WebView tidak dapat mengirimkan isi POST `fetch()` ke `shouldInterceptRequest`, sehingga panggilan runtime dirutekan melalui transport JavascriptInterface (`nativeHandleRuntimeCall`), alih-alih mengalami crash akibat isi permintaan nil
- Android: panggilan runtime `Screens.*` mengembalikan data sebenarnya — ScreenManager kini diisi saat aplikasi dimulai (sebelumnya tidak pernah dihubungkan sehingga `GetAll` mengembalikan nil)
- Android: pencatatan log debug framework tidak disertakan saat mengompilasi build produksi dan, dalam build debug, diarahkan melalui logcat dengan tag `Wails`
- Android: registri `hasListeners` yang sebenarnya, penanganan referensi/pengecualian JNI, dan siklus hidup halaman yang hanya memuat sekali (tanpa navigasi ganda)
- Memperbaiki kegagalan `wails3 generate bindings` dengan pesan "Akses ditolak" di Windows saat server pengembangan Vite berjalan, dengan menyinkronkan berkas yang dihasilkan ke direktori keluaran, alih-alih menimpa direktori tersebut melalui operasi penggantian nama (#5515)
- Memperbaiki crash fatal yang terjadi sesekali di macOS saat membaca informasi layar setelah perubahan tampilan: ID dan nama layar menyimpan pointer ke buffer `UTF8String` yang dirilis otomatis dan dapat dibebaskan sebelum disalin oleh Go (use-after-free). String tersebut kini di-`strdup` dan dibebaskan setelah konversi, sedangkan enumerasi layar dijalankan dalam autorelease pool eksplisit agar tidak lagi mengalami kebocoran saat dipanggil dari goroutine Go (#5556)
- Memperbaiki SIGSEGV yang terjadi sesekali di Linux saat assetserver menutup `WebKitURISchemeRequest`: `g_object_unref` terakhir dijalankan pada goroutine assetserver sehingga menyelesaikan GObject WebKit di luar thread utama GTK. Operasi unref kini diteruskan ke konteks utama GTK melalui `g_main_context_invoke` (#5557)

## v3.0.0-alpha.100 - 2026-06-13

## Ditambahkan

- Memperluas `MacWebviewPreferences` dengan opsi konfigurasi WKWebView tambahan: `EnableAutoplayWithoutUserAction`, `AllowsAirPlayForMediaPlayback`, `AllowsMagnification`, `JavaScriptCanOpenWindowsAutomatically`, `MinimumFontSize`, dan `ApplicationNameForUserAgent` (#5549)

## Diperbaiki

- Memperbaiki kegagalan `wails3 generate bindings` dengan pesan "Akses ditolak" di Windows saat server pengembangan Vite berjalan, dengan menyinkronkan berkas yang dihasilkan ke direktori keluaran, alih-alih menimpa direktori tersebut melalui operasi penggantian nama (#5561)
- Memperbaiki peristiwa perubahan ukuran JS yang tidak dipicu untuk jendela tanpa bingkai di Linux; memperbaiki deteksi tepi bilah gulir untuk jendela tanpa bingkai (#5368)
- Memperbaiki updater di Windows yang gagal dengan pesan "tautan lintas perangkat tidak valid" saat direktori sementara berada di volume yang berbeda dari direktori instalasi (#5560)

## v3.0.0-alpha.99 - 2026-06-10

## Diperbaiki

- Memperbaiki kegagalan `wails3 generate bindings` dengan pesan "Akses ditolak" di Windows saat server pengembangan Vite berjalan, dengan menyinkronkan berkas yang dihasilkan ke direktori keluaran, alih-alih menimpa direktori tersebut melalui operasi penggantian nama (#5515)

## v3.0.0-alpha.98 - 2026-06-03

## Diperbaiki

- Memperbaiki antarmuka WebKit yang membeku di Linux saat tidak aktif (misalnya ketika inspector terbuka) dengan tidak lagi memaksakan `SA_ONSTACK` pada `SIGUSR1`, yang merusak sinkronisasi thread GC JavaScriptCore (#5527)

## v3.0.0-alpha.97 - 2026-05-31

## Ditambahkan

- Menambahkan halaman debugging dan panduan bekerja dengan `runtime/trace`

## Diubah

- Menghapus beberapa impor `_ "embed"` yang tidak diperlukan untuk sedikit merapikan kode

## Diperbaiki

- Memperbaiki batasan lebar/tinggi minimum yang tidak diterapkan setelah jendela keluar dari keadaan dimaksimalkan di Windows (#4593)
- Memperbaiki klik mouse yang menembus jendela dalam mode layar penuh dengan opsi jendela Frameless + Transparent (#4408)

## v3.0.0-alpha.96 - 2026-05-25

## Ditambahkan

- Menambahkan dukungan obfuscation Garble ([#4563](https://github.com/wailsapp/wails/issues/4563)): ID metode binding yang stabil, integrasi build/Taskfile (`build --obfuscated --garbleargs`, `generate bindings -obfuscated`), dan tag struct JSON pada setiap payload yang berhadapan dengan runtime (`EnvironmentInfo`, `OSInfo`, `Screen`, `Rect`, `Point`, `Size`, `Capabilities`) agar format wire tetap berfungsi setelah Garble mengganti nama field yang diekspor.

## v3.0.0-alpha.95 - 2026-05-20

## Ditambahkan

- Menambahkan halaman struktur proyek yang sebelumnya tidak tersedia

## Diubah

- Dokumentasi: mengubah beberapa diagram pada halaman arsitektur agar menggunakan diagram urutan sehingga tampilannya lebih rapi
- Dokumentasi: Tambahkan catatan tentang penginstalan D2 sebagai prasyarat untuk menjalankan

## Diperbaiki

- Perbaiki `wails3 generate appimage` pada konfigurasi default GTK4: bundler kini mendeteksi stack GTK dari biner sebelum mencari berkas runtime, sehingga memilih `libwebkitgtkinjectedbundle.so` (di bawah `webkitgtk-6.0/`) untuk build GTK4 dan `libwebkit2gtkinjectedbundle.so` (di bawah `webkit2gtk-4.1/`) untuk build `-tags gtk3`. Pemeriksaan `.relr.dyn` juga memeriksa `libgtk-4.so.1` sehingga stripping dinonaktifkan dengan benar pada toolchain modern apa pun stack-nya. (#5475)
- Perbaiki kegagalan `wails3 generate appimage` ketika dipanggil dengan `-builddir` relatif: bundler kini terlebih dahulu mengubah `-binary`, `-icon`, `-desktopfile`, `-builddir`, dan `-outputdir` menjadi path absolut agar `s.CD` di tengah alur tidak mengganggu goroutine pengunduhan AppRun atau pemeriksaan `ldd` setelah penyalinan.
- Perbaiki kegagalan `wails3 generate appimage` memindahkan AppImage akhir ke `-outputdir` ketika kolom `Name=` desktop tidak cocok dengan nama dasar biner: bundler kini memaksa plugin appimage milik linuxdeploy (melalui variabel lingkungan `OUTPUT`) menulis AppImage ke `<binary>-<arch>.AppImage`, bukan ke nama yang diturunkan dari berkas desktop.
- Perbaiki `events.Common.ApplicationStarted`, `Common.ThemeChanged`, `Common.SystemWillSleep`, dan `Common.SystemDidWake` yang tidak terpicu di Linux setelah stack GTK4 + WebKitGTK 6.0 dijadikan default pada alpha.93. `application_linux.go` `run()` default yang baru tidak memanggil `setupCommonEvents()` (yang meneruskan peristiwa `Linux.*` ke padanan `Common.*`-nya) ataupun `monitorPowerEvents()`. Pembantu pemantauan daya DBus kini digunakan bersama oleh jalur build GTK3 dan GTK4 melalui `application_linux_dbus.go`. (#5474)

## v3.0.0-alpha.94 - 2026-05-19

## Diperbaiki

- Perbaiki `events.Common.ApplicationStarted`, `Common.ThemeChanged`, `Common.SystemWillSleep`, dan `Common.SystemDidWake` yang tidak terpicu di Linux setelah stack GTK4 + WebKitGTK 6.0 dijadikan default pada alpha.93. `application_linux.go` `run()` default yang baru tidak memanggil `setupCommonEvents()` (yang meneruskan peristiwa `Linux.*` ke padanan `Common.*`-nya) ataupun `monitorPowerEvents()`. Pembantu pemantauan daya DBus kini digunakan bersama oleh jalur build GTK3 dan GTK4 melalui `application_linux_dbus.go`. (#5474)

## v3.0.0-alpha.93 - 2026-05-17

## Ditambahkan

- Tambahkan `XDG_SESSION_TYPE` ke keluaran `wails3 doctor` di Linux oleh @leaanthony

## Diperbaiki

- Perbaiki crash menu jendela di Wayland yang disebabkan oleh appmenu-gtk-module saat mengakses jendela yang belum direalisasikan (#4769) oleh @leaanthony
- Perbaiki crash aplikasi GTK ketika nama aplikasi berisi karakter yang tidak valid (spasi, tanda kurung, dan sebagainya) oleh @leaanthony
- Perbaiki galat "memori tidak cukup" saat menginisialisasi seret dan lepas di Windows (#4701) oleh @overlordtm
- Perbaiki kondisi balapan pada penyimpanan callback mainthread yang menggunakan RLock yang salah untuk menghapus entri peta (Linux, macOS, iOS) (#4424) oleh @leaanthony
- Perbaiki penanganan variabel saat meneruskan argumen baris perintah ke tugas. Variabel CLI yang ditentukan sebagai pasangan KEY=VALUE kini diinisialisasi dan disebarkan dengan benar selama eksekusi tugas.
- Perbaiki konflik NSWindowZoomButton di macOS: `MaximiseButtonState` dan `FullscreenButtonState` kini menerapkan status yang lebih ketat, baik saat awal dijalankan maupun saat runtime; kedua setter tidak lagi dapat secara diam-diam menimpa satu sama lain (#5319)
- Perbaiki sekumpulan bug lama yang sudah ada pada jalur build GTK3 lawas (`-tags gtk3`) dan terungkap oleh CodeRabbit pada #5463: peluncuran melalui asosiasi berkas tidak lagi melewati handler startup; `getTheme` kini aman terhadap batas dan tipe; `appName` tidak lagi membebaskan memori milik GLib; `clipboardGet` tidak lagi membocorkan `gchar*` yang dikembalikan oleh GTK; `Calloc` kini menggunakan receiver pointer (dan `NewCalloc` mengembalikan `*Calloc`) sehingga pool benar-benar melacak alokasi; `zoomOut` menggunakan kebalikan dari `zoomInFactor`, bukan pengali negatif yang dibatasi menjadi 1.0; `execJS` menggunakan kembali nama world kosong yang telah dialokasikan sebelumnya, bukan membocorkan satu `C.CString("")` per pemanggilan; `fmt.Println` untuk pengembangan telah dihapus dari `menuItem.setAccelerator`. Menyelesaikan #5465.
- Perbaiki kebocoran akibat receiver nilai `Calloc` yang sama pada jalur build GTK4 default (`linux_cgo.go`): gunakan receiver pointer + `NewCalloc() *Calloc` agar alokasi `c.String(...)` per jendela benar-benar dilacak dan dibebaskan.

## v3.0.0-alpha.92 - 2026-05-15

## Ditambahkan

- Ubah Taskfile agar pengelola paket frontend yang digunakan dapat dikendalikan melalui opsi `PACKAGE_MANAGER`
- Perkaya data templat dengan `{{.Opn}}` dan `{{.Cls}}` agar penulisan templat Taskfile lebih mudah diprediksi

## Diubah

- Ubah beberapa Taskfile yang ada agar menggunakan `{{.Opn}} and {{.Cls}}`

## Diperbaiki

- Perbaiki galat fatal runtime `concurrent map read and map write` di `linuxSystemTray` ketika menu tray diperbarui saat panel membacanya.
- Gunakan `log` alih-alih `fmt` untuk keluaran galat dan pelacakan tumpukan WebView2 agar pesan tidak hilang ketika aplikasi berjalan tanpa konsol yang terpasang di Windows.

## v3.0.0-alpha.91 - 2026-05-12

## Diubah

- Memperbarui SVG sponsor dalam [PR](https://github.com/wailsapp/wails/pull/5414) oleh `@github-actions[bot]`
- **PERUBAHAN INKOMPATIBEL (macOS):** Normalkan sistem koordinat macOS agar `GetScreens`, `Position`, dan `SetPosition` semuanya menggunakan ruang yang sama—titik logis, sumbu Y mengarah ke bawah, dengan `(0,0)` di sudut kiri atas layar utama. Ini selaras dengan Windows, GTK, serta API publik Electron dan web. Layar yang secara fisik berada di atas layar utama kini melaporkan `Bounds.Y` negatif (sebelumnya positif), dan nilai `Position()`/`SetPosition()` kini dinyatakan dalam titik logis, bukan `points × primaryScale`. Konsistensi konversi pulang-pergi `Position()` → `SetPosition()` tetap dipertahankan; nilai absolut yang dicatat dari build alpha sebelumnya atau solusi sementara yang dihitung secara manual (misalnya mengalikan dengan `primaryScale` atau membalik Y terhadap tinggi layar) perlu diperbarui. Menyelesaikan [#5117](https://github.com/wailsapp/wails/issues/5117).

## Diperbaiki

- Validasi nama sinyal DBus dan panjang body secara defensif untuk mencegah panic dalam [PR](https://github.com/wailsapp/wails/pull/5416) oleh @leaanthony
- Perbaiki masalah keamanan memori dalam penanganan menu GTK di Linux pada [PR](https://github.com/wailsapp/wails/pull/5363) oleh @leaanthony
- Deteksi GPU NVIDIA dan nonaktifkan renderer DMA-BUF di Linux dalam [PR](https://github.com/wailsapp/wails/pull/5295) oleh @leaanthony
- Perbaiki konversi Y lintas layar `SetPosition` di macOS: gunakan tinggi layar utama sebagai referensi global agar jendela ditempatkan pada posisi yang benar di monitor yang bergeser secara vertikal dari layar utama dalam [#5117](https://github.com/wailsapp/wails/issues/5117)
- Perbaiki templat PR git agar mengarah ke URL umpan balik yang benar dalam [PR](https://github.com/wailsapp/wails/pull/5109) oleh @wayneforrest
- Memperbaiki serangkaian crash `SetMenu` pada systray Windows yang disebabkan oleh syscall `DestroyMenu` yang rusak karena meneruskan empat argumen, bukan satu, sehingga setiap panggilan mengembalikan FALSE dan tidak membebaskan apa pun. Juga membebaskan handle HMENU dan HBITMAP (termasuk yang dialokasikan saat runtime melalui `MenuItem.SetBitmap`) ketika menu dibuat ulang, mereset peta kotak centang/tombol radio yang kedaluwarsa di `Win32Menu.Update`, serta menghapus panggilan `Update()` yang redundan di `systemtray.updateMenu` karena panggilan tersebut menyebabkan alokasi berlipat ganda. Aplikasi systray yang berjalan lama tidak lagi membocorkan objek GDI/USER setiap kali menu dibuat ulang.

## v3.0.0-alpha.90 - 2026-05-11

## Ditambahkan

- Menambahkan nama aplikasi yang dapat dikonfigurasi untuk User-Agent WKWebView di macOS dalam [PR](https://github.com/wailsapp/wails/pull/5261) oleh @vinhvoit225
- Menambahkan dependensi tidak langsung github.com/coder/websocket ke contoh gin-service dalam [PR](https://github.com/wailsapp/wails/pull/5400) oleh @taliesin-ai
- Menambahkan dukungan perbandingan kesetaraan mendalam ke pengujian aset build dalam [PR](https://github.com/wailsapp/wails/pull/5402) oleh @leaanthony

## Diubah

- Menggabungkan output build ke dalam direktori assets dalam [PR](https://github.com/wailsapp/wails/pull/5401) oleh @taliesin-ai
- Memperbarui SVG sponsor dalam [PR](https://github.com/wailsapp/wails/pull/5399) oleh `@github-actions[bot]`

## Diperbaiki

- Menggunakan objek notifikasi untuk pesan instans tunggal macOS dalam [PR](https://github.com/wailsapp/wails/pull/5289) oleh @overlordtm
- Memproses callback Windows secara berkelompok untuk mencegah hilangnya promise saat beban tinggi dalam [PR](https://github.com/wailsapp/wails/pull/5383) oleh @taliesin-ai

## v3.0.0-alpha.89 - 2026-05-10

## Ditambahkan

- Menambahkan job go<em>test</em>results untuk mengagregasikan hasil pengujian Go dalam [PR](https://github.com/wailsapp/wails/pull/5316) oleh @leaanthony

## Diubah

- Membagi payload RPC berukuran besar secara kondisional menjadi potongan-potongan yang dikirim melalui permintaan POST dalam [PR](https://github.com/wailsapp/wails/pull/5369) oleh @leaanthony
- Meningkatkan Vite dari 5.x.x ke 8.0.0 di semua templat frontend dalam [PR](https://github.com/wailsapp/wails/pull/5386) oleh @leaanthony
- Memigrasikan konfigurasi port server pengembangan Vite ke variabel lingkungan dalam [PR](https://github.com/wailsapp/wails/pull/5365) oleh @leaanthony
- Mengonfigurasi server pengembangan Vite agar terikat ke 127.0.0.1 di semua templat dalam [PR](https://github.com/wailsapp/wails/pull/5361) oleh @leaanthony
- Memperbarui SVG sponsor dalam [PR](https://github.com/wailsapp/wails/pull/5384) oleh `@github-actions[bot]`

## Diperbaiki

- Membersihkan stub templat Info.plist selama pembaruan build-assets dalam [PR](https://github.com/wailsapp/wails/pull/5312) oleh @leaanthony
- Memperbaiki status kedaluwarsa pada menu macOS dengan menerapkan mutator item menu (`setMenuItemChecked()`, `setMenuItemLabel()`, `setMenuItemDisabled()`, `setMenuItemHidden()`, `setMenuItemTooltip()`) secara sinkron di thread utama, sehingga menghilangkan race `dispatch_async` yang menyebabkan menu merender status sebelumnya ketika dibuka kembali dengan cepat (#5002)
- Mengabaikan file `*_test.go` dalam mode pengembangan untuk mencegah build ulang yang tidak diperlukan dalam [PR](https://github.com/wailsapp/wails/pull/5203) oleh @leaanthony
- Mencegah segfault Menu.Update() ketika aplikasi tidak berjalan dalam [PR](https://github.com/wailsapp/wails/pull/5291) oleh @wucm667
- Menggunakan lastSizeWParam untuk mengatur kapan bilah menu digambar ulang di Windows dalam [PR](https://github.com/wailsapp/wails/pull/5382) oleh @taliesin-ai

## v3.0.0-alpha.88 - 2026-05-09

## Diubah

- Mengubah HiddenOnTaskbar agar menggunakan WS<em>EX</em>TOOLWINDOW dalam [PR](https://github.com/wailsapp/wails/pull/5371) oleh @leaanthony
- Mengurutkan ulang dependensi dan menghapus direktif replace untuk webview2 di go.mod dalam [PR](https://github.com/wailsapp/wails/pull/5370) oleh @atterpac
- Memperbarui SVG sponsor dalam [PR](https://github.com/wailsapp/wails/pull/5358) oleh `@github-actions[bot]`

## Diperbaiki

- Menghapus alias indirection generik dan menggabungkan tipe kunci map dalam [PR](https://github.com/wailsapp/wails/pull/5331) oleh @fbbdev

## Dihapus

- Menghapus alur kerja PR-master, termasuk dokumentasi, pengujian Go, dan fungsi untuk melewati pengujian, dalam [PR](https://github.com/wailsapp/wails/pull/5377) oleh @leaanthony

## v3.0.0-alpha.87 - 2026-05-07

## Ditambahkan

- Menambahkan dokumentasi bahasa Korea untuk Wails v3 dalam [PR](https://github.com/wailsapp/wails/pull/5352) oleh @leaanthony
- Menambahkan dokumentasi bahasa Prancis untuk instalasi dan panduan memulai cepat dalam [PR](https://github.com/wailsapp/wails/pull/5354) oleh @leaanthony
- Menambahkan dokumentasi bahasa Portugis untuk panduan memulai cepat, konsep, dan komunitas dalam [PR](https://github.com/wailsapp/wails/pull/5355) oleh @leaanthony

## v3.0.0-alpha.86 - 2026-05-06

## Ditambahkan

- Menambahkan pelokalan dokumentasi bahasa Prancis dalam [PR](https://github.com/wailsapp/wails/pull/5328) oleh @leaanthony
- Menambahkan locale bahasa Jerman ke situs dokumentasi dalam [PR](https://github.com/wailsapp/wails/pull/5343) oleh @leaanthony

## Diubah

- Mendaftarkan semua 8 locale terjemahan dalam konfigurasi dokumentasi dalam [PR](https://github.com/wailsapp/wails/pull/5347) oleh @leaanthony
- Memperbarui berbagai file terkait Windows untuk WebView2 dalam [PR](https://github.com/wailsapp/wails/pull/5317) oleh @leaanthony

## Diperbaiki

- Memisahkan pengiriman dialog antara GTK3 dan GTK4 untuk Linux dalam [PR](https://github.com/wailsapp/wails/pull/5340) oleh @leaanthony
- Memastikan callback dialog dijalankan di thread GTK, sehingga memperbaiki segfault, dalam [PR](https://github.com/wailsapp/wails/pull/5339) oleh @leaanthony

## v3.0.0-alpha.85 - 2026-05-05

## Ditambahkan

- Menambahkan URL templat PR ke repositori dalam [PR](https://github.com/wailsapp/wails/pull/5179) oleh @leaanthony
- Menambahkan dokumentasi bahasa Jerman untuk Wails v3 dalam [PR](https://github.com/wailsapp/wails/pull/5330) oleh @leaanthony

## v3.0.0-alpha.84 - 2026-05-03

## Ditambahkan

- Menambahkan opsi untuk menonaktifkan keluarnya mode layar penuh melalui tombol Escape di macOS dalam [PR](https://github.com/wailsapp/wails/pull/5307) oleh @leaanthony
- Menambahkan opsi untuk menonaktifkan keluarnya mode layar penuh melalui tombol Escape di macOS dalam [PR](https://github.com/wailsapp/wails/pull/5310) oleh @leaanthony
- Menambahkan dokumentasi etalase komunitas Pausa dalam [PR](https://github.com/wailsapp/wails/pull/5288) oleh @yuseferi

## Diubah

- Memperbarui SVG sponsor dalam [PR](https://github.com/wailsapp/wails/pull/5308) oleh `@github-actions[bot]`
- Memperbarui perintah pembuatan ikon agar dapat menangani platform yang tidak didukung dalam [PR](https://github.com/wailsapp/wails/pull/5309) oleh @leaanthony
- Mengganti API layar penuh boolean dengan ButtonState tiga status dan mengimplementasikan binding platform dalam [PR](https://github.com/wailsapp/wails/pull/5224) oleh @leaanthony

## Diperbaiki

- Melindungi operasi fokus WebView2 dari status pengontrol nil dalam [PR](https://github.com/wailsapp/wails/pull/5315) oleh @leaanthony
- Memperbarui alur kerja GitHub Actions agar merujuk cabang dasar PR dengan benar dalam [PR](https://github.com/wailsapp/wails/pull/5313) oleh @leaanthony
- Mengabaikan file `*_test.go` dalam mode pengembangan untuk mencegah build ulang yang tidak diperlukan dalam [PR](https://github.com/wailsapp/wails/pull/5203) oleh @leaanthony
- Mencegah segfault pada Menu.Update() saat aplikasi tidak berjalan dalam [PR](https://github.com/wailsapp/wails/pull/5291) oleh @wucm667

## v3.0.0-alpha.83 - 2026-05-02

## Ditambahkan

- Menambahkan flag InstallScope dan opsi build untuk penginstalan tingkat mesin/pengguna dalam [PR](https://github.com/wailsapp/wails/pull/5094) oleh @symball
- Menambahkan metode SetScreen tanpa operasi ke BrowserWindow agar memenuhi antarmuka Window dalam [PR](https://github.com/wailsapp/wails/pull/5294) oleh @leaanthony

## Diperbaiki

- Mendeteksi GPU NVIDIA dan menonaktifkan perender DMA-BUF di Linux dalam [PR](https://github.com/wailsapp/wails/pull/5295) oleh @leaanthony
- Memperbaiki templat PR Git agar mengarah ke URL umpan balik yang benar dalam [PR](https://github.com/wailsapp/wails/pull/5109) oleh @wayneforrest
- Memperbaiki serangkaian crash `SetMenu` pada baki sistem Windows yang disebabkan oleh syscall `DestroyMenu` yang rusak dan meneruskan empat argumen alih-alih satu, sehingga setiap panggilan mengembalikan FALSE dan tidak membebaskan apa pun. Juga membebaskan handle HMENU dan HBITMAP (termasuk yang dialokasikan saat runtime melalui `MenuItem.SetBitmap`) ketika menu dibuat ulang, mereset pemetaan kotak centang/tombol radio yang sudah usang dalam `Win32Menu.Update`, serta menghapus panggilan `Update()` redundan dalam `systemtray.updateMenu` yang menggandakan alokasi. Aplikasi baki sistem yang berjalan lama tidak lagi membocorkan objek GDI/USER setiap kali menu dibuat ulang.

## v3.0.0-alpha.82 - 2026-05-01

## Diperbaiki

- Memperbaiki pembuatan file desktop agar menangani nama desktop dengan benar dalam [PR](https://github.com/wailsapp/wails/pull/5232) oleh @leaanthony

## v3.0.0-alpha.81 - 2026-04-30

## Diubah

- Menyesuaikan jadwal rilis nightly menjadi pukul 15:00 UTC dalam [PR](https://github.com/wailsapp/wails/pull/5286) oleh @leaanthony

## Diperbaiki

- Memperbaiki nilai Screen Bounds, WorkArea, dan Size yang menjadi setengah pada Mac Retina -  (#5168)

## v3.0.0-alpha.80 - 2026-04-29

## Diubah

- Memperbarui dependensi dokumentasi dan pemuat koleksi konten dalam [PR](https://github.com/wailsapp/wails/pull/5285) oleh @leaanthony

## v3.0.0-alpha.79 - 2026-04-29

## Ditambahkan

- Memberikan izin actions: write kepada tugas trigger-release dalam [PR](https://github.com/wailsapp/wails/pull/5270) oleh @leaanthony

## Diubah

- Tugas rilis kini secara default menggunakan cabang master dan memperbarui redaksi changelog dalam [PR](https://github.com/wailsapp/wails/pull/5283) oleh @leaanthony
- Memperbarui alur kerja auto-changelog agar menggunakan versi terbaru dalam [PR](https://github.com/wailsapp/wails/pull/5282) oleh @leaanthony
- Meningkatkan efisiensi alur kerja dengan menambahkan filter jalur dan menghapus alur kerja yang tidak lagi digunakan dalam [PR](https://github.com/wailsapp/wails/pull/5280) oleh @leaanthony
- Memperbarui dokumentasi agar tautan contoh merujuk ke cabang master dalam [PR](https://github.com/wailsapp/wails/pull/5274) oleh @leaanthony
- Memperbarui dokumentasi dan contoh untuk v3 dalam [PR](https://github.com/wailsapp/wails/pull/5272) oleh @leaanthony

## Diperbaiki

- Meningkatkan proksi balik dengan logika percobaan ulang dan pemaksaan IPv4 untuk pengembangan dalam [PR](https://github.com/wailsapp/wails/pull/5265) oleh @AkagiYui
- Menulis ulang alur kerja pemicu changelog yang belum dirilis dalam [PR](https://github.com/wailsapp/wails/pull/5281) oleh @leaanthony

## Dihapus

- Menghapus skrip pengujian shell untuk berbagai keperluan pengujian dalam [PR](https://github.com/wailsapp/wails/pull/5267) oleh @leaanthony
- Menghapus alur kerja deployment dokumentasi v3-alpha dan catatan CNAME dalam [PR](https://github.com/wailsapp/wails/pull/5266) oleh @leaanthony

### Ditambahkan

- Menambahkan entri Perutean Frontend ke navigasi bilah samping dalam [PR](https://github.com/wailsapp/wails/pull/5196) oleh @leaanthony
- Menambahkan panduan perutean frontend dengan rekomendasi khusus framework dalam [PR](https://github.com/wailsapp/wails/pull/5185) oleh @leaanthony
- Menambahkan dukungan untuk lembar modal (macOS)
- Menaikkan versi ghw untuk dukungan perangkat Apple yang lebih baik oleh @leaanthony (#4977)
- Menambahkan metode `GetBadge` ke layanan dock
- Menambahkan flag `-tags` ke perintah `wails3 build` untuk meneruskan tag build Go khusus (misalnya, `wails3 build -tags gtk4`) (#4957)
- Menambahkan dokumentasi untuk pembuatan enum otomatis dalam generator binding, termasuk halaman Enums khusus dan navigasi bilah samping (#4972)
- Menambahkan flag `-tags` ke perintah `wails3 build` untuk meneruskan tag build Go khusus (misalnya, `wails3 build -tags gtk4`) (#4957)
- Menambahkan contoh Web API dalam `v3/examples/web-apis/` yang mendemonstrasikan 41 API browser, termasuk Penyimpanan (localStorage, sessionStorage, IndexedDB, Cache API), Jaringan (Fetch, WebSocket, XMLHttpRequest, EventSource, Beacon), Media (Canvas, WebGL, Web Audio, MediaDevices, MediaRecorder, Speech Synthesis), Perangkat (Geolocation, Clipboard, Fullscreen, Device Orientation, Vibration, Gamepad), Performa (Performance API, Mutation Observer, Intersection/Resize Observer), UI (Web Components, Pointer Events, Selection, Dialog, Drag and Drop), dan lainnya
- Menambahkan contoh pemeriksa kompatibilitas WebView API (`v3/examples/webview-api-check/`) yang menguji lebih dari 200 API browser di berbagai platform
- Menambahkan paket `internal/libpath` untuk menemukan jalur pustaka native di Linux dengan pencarian paralel, caching, serta dukungan untuk Flatpak/Snap/Nix
- **WIP:** Menambahkan dukungan eksperimental WebKitGTK 6.0 / GTK4 untuk Linux, tersedia melalui `-tags gtk4` (GTK3/WebKit2GTK 4.1 tetap menjadi pilihan default)
- Catatan: Pada pengelola jendela tiling (misalnya Hyprland dan Sway), operasi Minimize/Maximize mungkin tidak berfungsi sebagaimana mestinya karena geometri jendela dikendalikan oleh pengelola jendela
- Menambahkan cara menggunakan **Penangan Sekali Pakai** dalam dokumentasi **Mendengarkan Peristiwa di JavaScript** oleh @AbdelhadiSeddar
- Menambahkan opsi `UseApplicationMenu` ke `WebviewWindowOptions` agar jendela di Windows/Linux dapat mewarisi menu aplikasi yang ditetapkan melalui `app.Menu.Set()` oleh @leaanthony
- Menambahkan dukungan untuk menggunakan file `.icon` (format Apple Icon Composer) guna menghasilkan ikon Liquid Glass dan katalog aset (macOS) (#4934) oleh @wimaha
- Menambahkan mode server eksperimental untuk deployment headless/web (`-tags server`). Memungkinkan aplikasi Wails dijalankan sebagai server HTTP tanpa dependensi GUI native. Build dengan `wails3 task build:server`. Lihat `examples/server` untuk detailnya.
- Menambahkan paket `internal/libpath` untuk menemukan jalur pustaka native di Linux melalui pencarian paralel, caching, serta dukungan untuk Flatpak/Snap/Nix
- Menambahkan opsi `CollectionBehavior` ke `MacWindow` untuk mengendalikan perilaku jendela di seluruh Spaces dan layar penuh macOS (#4756) oleh @leaanthony
- Menambahkan pengujian unit untuk pkg/application oleh @leaanthony
- Menambahkan dukungan protokol khusus ke pemaketan MSIX oleh @leaanthony
- Menambahkan deteksi lingkungan desktop di Linux [PR #4797](https://github.com/wailsapp/wails/pull/4797)
- Menambahkan metode `Window.Print()` ke runtime JavaScript untuk memicu dialog cetak dari frontend (#4290) oleh @leaanthony
- Menambahkan `XDG_SESSION_TYPE` ke keluaran `wails3 doctor` di Linux oleh @leaanthony
- Menambahkan peristiwa perubahan pemuatan WebKit2 tambahan untuk Linux: `WindowLoadStarted`, `WindowLoadRedirected`, `WindowLoadCommitted`, `WindowLoadFinished` (#3896) oleh @leaanthony
- Menambahkan `XDG_SESSION_TYPE` ke keluaran `wails3 doctor` di Linux oleh @leaanthony
- Menghasilkan file `.desktop` selama build Linux, bukan hanya saat pemaketan (#4575)
- Menambahkan dokumentasi dependensi runtime Linux dengan nama paket khusus distribusi dan contoh pemaketan nfpm (#4339) oleh @leaanthony
- Menambahkan informasi versi driver NVIDIA ke keluaran `wails3 doctor` di Linux oleh @leaanthony
- Menambahkan origin ke handler pesan mentah oleh @APshenkin dalam [PR](https://github.com/wailsapp/wails/pull/4710)
- Menambahkan dukungan universal link untuk macOS oleh @APshenkin dalam [PR](https://github.com/wailsapp/wails/pull/4712)
- Merefaktor lapisan transport binding oleh @APshenkin dalam [PR](https://github.com/wailsapp/wails/pull/4702)
- Menambahkan pengenal aria-label ke templat helloworld agar aplikasi contoh dapat diuji dengan mudah oleh klien pengujian Appium, oleh @chinenual dalam [PR](https://github.com/wailsapp/wails/pull/4760)
- Menambahkan origin ke handler pesan mentah oleh @APshenkin dalam [PR](https://github.com/wailsapp/wails/pull/4710)
- Menambahkan dukungan universal link untuk macOS oleh @APshenkin dalam [PR](https://github.com/wailsapp/wails/pull/4712)
- Merefaktor lapisan transport binding oleh @APshenkin dalam [PR](https://github.com/wailsapp/wails/pull/4702)
- Peristiwa Bertipe oleh @fbbdev dan @ianvs dalam [#4633](https://github.com/wailsapp/wails/pull/4633)
- Menambahkan contoh `systray-clock` yang menampilkan tray headless dengan pembaruan tooltip secara langsung (#4653).
- Menambahkan templat Protokol NSIS untuk Windows oleh @Tolfx dalam #4510
- Menambahkan pengujian untuk build-assets oleh @Tolfx dalam #4510
- macOS: Menampilkan kontrol jendela native pada bilah menu dalam [#4588](https://github.com/wailsapp/wails/pull/4588) oleh @nidib
- Menambahkan layanan Dock macOS untuk menyembunyikan/menampilkan ikon aplikasi di Dock oleh @popaprozac dalam [PR](https://github.com/wailsapp/wails/pull/4451)
- Menambahkan layanan Dock macOS untuk menyembunyikan/menampilkan ikon aplikasi di Dock oleh @popaprozac dalam [PR](https://github.com/wailsapp/wails/pull/4451)
- Menambahkan dukungan efek Liquid Glass native untuk macOS dengan NSGlassEffectView (macOS 15.0+) dan fallback NSVisualEffectView, termasuk opsi penyesuaian material yang komprehensif, oleh @leaanthony dalam [#4534](https://github.com/wailsapp/wails/pull/4534)
- Sanitasi URL Browser oleh @leaanthony dalam [#4500](https://github.dev/wailsapp/wails/pull/4500). Berdasarkan [#4484](https://github.com/wailsapp/wails/pull/4484) oleh @APShenkin.
- Menambahkan Perlindungan Konten di Windows/Mac oleh [@leaanthony](https://github.com/leaanthony) berdasarkan karya asli [@Taiterbase](https://github.com/Taiterbase) dalam [PR](https://github.com/wailsapp/wails/pull/4241) ini
- Menambahkan dukungan untuk meneruskan variabel CLI ke perintah Task melalui alias `wails3 build` dan `wails3 package` (#4422) oleh @leaanthony dalam [PR](https://github.com/wailsapp/wails/pull/4488)
- Dukungan dropzone dengan event sourcing untuk data elemen yang dilepas oleh [@atterpac](https://github.com/atterpac) dalam [#4318](https://github.com/wailsapp/wails/pull/4318)
- Menambahkan `AdditionalLaunchArgs` ke opsi `WindowsWindow` agar argumen baris perintah tambahan dapat diteruskan ke browser WebView2, dalam [PR](https://github.com/wailsapp/wails/pull/4467)
- Menambahkan eksekusi go mod tidy secara otomatis setelah wails init oleh [@triadmoko](https://github.com/triadmoko) dalam [PR](https://github.com/wailsapp/wails/pull/4286)
- Fitur Snapassist Windows oleh @leaanthony dalam [PR](https://github.dev/wailsapp/wails/pull/4463)
- Menambahkan `AdditionalLaunchArgs` ke opsi `WindowsWindow` agar argumen baris perintah tambahan dapat diteruskan ke browser WebView2, dalam [PR](https://github.com/wailsapp/wails/pull/4467)
- Menambahkan eksekusi go mod tidy secara otomatis setelah wails init oleh [@triadmoko](https://github.com/triadmoko) dalam [PR](https://github.com/wailsapp/wails/pull/4286)
- Fitur Snapassist Windows oleh @leaanthony dalam [PR](https://github.dev/wailsapp/wails/pull/4463)
- Menambahkan implementasi `getAccentColor` Windows oleh [@almas-x](https://github.com/almas-x) dalam [PR](https://github.com/wailsapp/wails/pull/4427)
- Menambahkan implementasi `getAccentColor` Windows oleh [@almas-x](https://github.com/almas-x) dalam [PR](https://github.com/wailsapp/wails/pull/4427)
- Menu dan bilah menu bertema gelap di Windows. Oleh @leaanthony dalam [a29b4f0861b1d0a700e9eb213c6f1076ec40efd5](https://github.com/wailsapp/wails/commit/a29b4f0861b1d0a700e9eb213c6f1076ec40efd5)
- Mengganti nama layanan bawaan agar binding JS/TS lebih jelas, oleh @popaprozac dalam [PR](https://github.com/wailsapp/wails/pull/4405)
- `app.Env.GetAccentColor` untuk mendapatkan warna aksen sistem pengguna. Berfungsi di macOS. Oleh [@etesam913](https://github.com/etesam913)
- Menambahkan API `window.ToggleFrameless()` oleh [@atterpac](https://github.com/atterpac) dalam [#4137](https://github.com/wailsapp/wails/pull/4137)
- Tambahkan dependensi build khusus distribusi untuk Linux oleh @leaanthony dalam [PR](https://github.com/wailsapp/wails/pull/4345)
- Panduan binding ditambahkan oleh @atterpac dalam [PR](https://github.com/wailsapp/wails/pull/4404)
- **Infrastruktur Pengujian yang Tertata**: Memindahkan file pengujian Docker ke direktori khusus `test/docker/` dengan image yang dioptimalkan dan keandalan build yang ditingkatkan oleh [@leaanthony](https://github.com/leaanthony) dalam [#4359](https://github.com/wailsapp/wails/pull/4359)
- **Pola Pengelolaan Sumber Daya yang Disempurnakan**: Menambahkan pembersihan event handler yang tepat dan pengelolaan goroutine yang peka konteks dalam contoh oleh [@leaanthony](https://github.com/leaanthony) dalam [#4359](https://github.com/wailsapp/wails/pull/4359)
- Mendukung build AppImage aarch64 oleh [@AkshayKalose](https://github.com/AkshayKalose) dalam [#3981](https://github.com/wailsapp/wails/pull/3981)
- Tambahkan bagian diagnostik ke `wails doctor` oleh [@leaanthony](https://github.com/leaanthony)
- Tambahkan jendela ke konteks saat memanggil metode layanan oleh [@leaanthony](https://github.com/leaanthony)
- Tambahkan contoh `window-call` untuk menunjukkan cara mengetahui jendela mana yang memanggil suatu layanan oleh [@leaanthony](https://github.com/leaanthony)
- Panduan Menu baru oleh [@leaanthony](https://github.com/leaanthony)
- Penanganan panic yang lebih baik oleh [@leaanthony](https://github.com/leaanthony)
- Panduan Menu baru oleh [@leaanthony](https://github.com/leaanthony)
- Tambahkan komentar dokumentasi untuk Service API oleh [@fbbdev](https://github.com/fbbdev) dalam [#4024](https://github.com/wailsapp/wails/pull/4024)
- Tambahkan fungsi `application.NewServiceWithOptions` untuk menginisialisasi layanan dengan konfigurasi tambahan oleh [@leaanthony](https://github.com/leaanthony) dalam [#4024](https://github.com/wailsapp/wails/pull/4024)
- Kontrol menu ditingkatkan oleh [@FalcoG](https://github.com/FalcoG) dan [@leaanthony](https://github.com/leaanthony) dalam [#4031](https://github.com/wailsapp/wails/pull/4031)
- Dokumentasi tambahan oleh [@leaanthony](https://github.com/leaanthony)
- Mendukung pembatalan event dalam listener event standar oleh [@leaanthony](https://github.com/leaanthony)
- Dukungan Systray untuk `Hide`, `Show`, dan `Destroy` oleh [@leaanthony](https://github.com/leaanthony)
- Dukungan Systray untuk `SetTooltip` oleh [@leaanthony](https://github.com/leaanthony). Ide awal oleh [@lujihong](https://github.com/wailsapp/wails/issues/3487#issuecomment-2633242304)
- Laporkan path paket dalam peringatan generator binding tentang tipe yang tidak didukung oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Tambahkan dukungan generator binding untuk alias generik oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Tambahkan dukungan generator binding untuk flag JSON `omitzero` oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Tambahkan direktif `//wails:ignore` untuk mencegah pembuatan binding bagi metode layanan yang dipilih oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Tambahkan direktif `//wails:internal` pada layanan dan model untuk mengizinkan tipe yang diekspor di Go tetapi tidak di JS/TS oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Tambahkan dukungan generator binding untuk konstanta bertipe alias agar enum dengan pengetikan lemah dapat digunakan oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Tambahkan pengujian generator binding untuk fitur Go 1.24 oleh [@fbbdev](https://github.com/fbbdev) dalam [#4068](https://github.com/wailsapp/wails/pull/4068)
- Tambahkan dukungan untuk macOS 15 "Sequoia" ke `OSInfo.Branding` guna meningkatkan deteksi versi OS dalam [#4065](https://github.com/wailsapp/wails/pull/4065)
- Tambahkan hook `PostShutdown` untuk menjalankan kode khusus setelah proses shutdown selesai oleh [@fbbdev](https://github.com/fbbdev) dalam [#4066](https://github.com/wailsapp/wails/pull/4066)
- Tambahkan struct `FatalError` untuk mendukung deteksi kesalahan fatal dalam handler kesalahan khusus oleh [@fbbdev](https://github.com/fbbdev) dalam [#4066](https://github.com/wailsapp/wails/pull/4066)
- Standarkan dan dokumentasikan urutan startup dan shutdown layanan oleh [@fbbdev](https://github.com/fbbdev) dalam [#4066](https://github.com/wailsapp/wails/pull/4066)
- Tambahkan kerangka pengujian untuk urutan startup/shutdown aplikasi serta pengujian startup/shutdown layanan oleh [@fbbdev](https://github.com/fbbdev) dalam [#4066](https://github.com/wailsapp/wails/pull/4066)
- Tambahkan metode `RegisterService` untuk mendaftarkan layanan setelah aplikasi dibuat oleh [@fbbdev](https://github.com/fbbdev) dalam [#4066](https://github.com/wailsapp/wails/pull/4066)
- Tambahkan field `MarshalError` dalam opsi aplikasi dan layanan untuk penanganan kesalahan khusus dalam pemanggilan binding oleh [@fbbdev](https://github.com/fbbdev) dalam [#4066](https://github.com/wailsapp/wails/pull/4066)
- Tambahkan wrapper promise yang dapat dibatalkan dan meneruskan permintaan pembatalan melalui rantai promise oleh [@fbbdev](https://github.com/fbbdev) dalam [#4100](https://github.com/wailsapp/wails/pull/4100)
- Tambahkan kemampuan untuk mengaitkan pembatalan pemanggilan binding dengan `AbortSignal` oleh [@fbbdev](https://github.com/fbbdev) dalam [#4100](https://github.com/wailsapp/wails/pull/4100)
- Dukungan untuk atribut `data-wml-*` pada WML ditambahkan oleh [@leaanthony](https://github.com/leaanthony), di samping atribut `wml-*` yang biasa digunakan
- Tambahkan metode `Configure` pada semua layanan untuk konfigurasi tahap akhir/konfigurasi ulang dinamis oleh [@fbbdev](https://github.com/fbbdev) dalam [#4067](https://github.com/wailsapp/wails/pull/4067)
- Layanan `fileserver` mengirimkan respons 503 Service Unavailable saat belum dikonfigurasi oleh [@fbbdev](https://github.com/fbbdev) dalam [#4067](https://github.com/wailsapp/wails/pull/4067)
- Layanan `kvstore` secara default menyediakan penyimpanan kunci-nilai dalam memori saat belum dikonfigurasi; kontribusi oleh [@fbbdev](https://github.com/fbbdev) dalam [#4067](https://github.com/wailsapp/wails/pull/4067)
- Tambahkan metode `Load` pada layanan `kvstore` untuk memuat ulang data dari file setelah perubahan konfigurasi oleh [@fbbdev](https://github.com/fbbdev) dalam [#4067](https://github.com/wailsapp/wails/pull/4067)
- Tambahkan metode `Clear` pada layanan `kvstore` untuk menghapus semua kunci oleh [@fbbdev](https://github.com/fbbdev) dalam [#4067](https://github.com/wailsapp/wails/pull/4067)
- Tambahkan tipe `Level` dalam layanan `log` untuk menyediakan konstanta level log di sisi JS oleh [@fbbdev](https://github.com/fbbdev) dalam [#4067](https://github.com/wailsapp/wails/pull/4067)
- Tambahkan metode `Log` pada layanan `log` untuk menentukan level log secara dinamis oleh [@fbbdev](https://github.com/fbbdev) dalam [#4067](https://github.com/wailsapp/wails/pull/4067)
- Layanan `sqlite` secara default menyediakan DB dalam memori saat belum dikonfigurasi oleh [@fbbdev](https://github.com/fbbdev) dalam [#4067](https://github.com/wailsapp/wails/pull/4067)
- Tambahkan metode `Close` pada layanan `sqlite` untuk menutup DB secara manual oleh [@fbbdev](https://github.com/fbbdev) dalam [#4067](https://github.com/wailsapp/wails/pull/4067)
- Tambahkan dukungan pembatalan untuk metode kueri pada layanan `sqlite` oleh [@fbbdev](https://github.com/fbbdev) dalam [#4067](https://github.com/wailsapp/wails/pull/4067)
- Tambahkan dukungan prepared statement ke layanan `sqlite` dengan binding JS oleh [@fbbdev](https://github.com/fbbdev) dalam [#4067](https://github.com/wailsapp/wails/pull/4067)
- Dukungan Gin oleh [Lea Anthony](https://github.com/leaanthony) dalam [PR](https://github.com/wailsapp/wails/pull/3537), berdasarkan karya asli [@AnalogJ](https://github.com/AnalogJ) dalam [PR](https://github.com/wailsapp/wails/pull/3537) ini
- Perbaiki penyimpanan otomatis dan penyimpanan otomatis kata sandi yang selalu diaktifkan oleh [@oSethoum](https://github.com/osethoum) dalam [#4134](https://github.com/wailsapp/wails/pull/4134)
- Tambahkan `SetMenu()` pada jendela agar menu dapat ditetapkan pada jendela oleh [@leaanthony](https://github.com/leaanthony)
- Tambahkan dukungan Notifikasi oleh [@popaprozac](https://github.com/popaprozac) dalam [#4098](https://github.com/wailsapp/wails/pull/4098)
-  Tambahkan dukungan Asosiasi File untuk mac oleh [@wimaha](https://github.com/wimaha) dalam [#4177](https://github.com/wailsapp/wails/pull/4177)
- Tambahkan `wails3 tool version` untuk menaikkan versi semantik oleh [@leaanthony](https://github.com/leaanthony)
- Tambahkan dukungan lencana untuk macOS dan Windows oleh [@popaprozac](https://github.com/popaprozac) dalam [#](https://github.com/wailsapp/wails/pull/4234)
- Tambahkan dukungan untuk event terdaftar/bertipe ketat oleh [@fbbdev](https://github.com/fbbdev) dan [@IanVS](https://github.com/IanVS) dalam [#4161](https://github.com/wailsapp/wails/pull/4161)
- Tambahkan kemampuan untuk mendaftarkan hook bagi event khusus oleh [@fbbdev](https://github.com/fbbdev) dan [@IanVS](https://github.com/IanVS) dalam [#4161](https://github.com/wailsapp/wails/pull/4161)
- `app.OpenFileManager(path string, selectFile bool)` untuk membuka pengelola file sistem pada jalur `path` dengan penyorotan opsional melalui `selectFile` oleh [@Krzysztofz01](https://github.com/Krzysztofz01) [@rcalixte](https://github.com/rcalixte)
- Flag `-git` baru untuk perintah `wails3 init` oleh [@leaanthony](https://github.com/leaanthony)
- Perintah `wails3 generate webview2bootstrapper` baru oleh [@leaanthony](https://github.com/leaanthony)
- Tambahkan metode `init()` di runtime untuk memungkinkan inisialisasi runtime secara manual oleh [@leaanthony](https://github.com/leaanthony)
- Tambahkan opsi `WindowDidMoveDebounceMS` ke WindowOptions milik Window oleh [@leaanthony](https://github.com/leaanthony)
- Tambahkan fitur Instans Tunggal oleh [@leaanthony](https://github.com/leaanthony). Berdasarkan [PR v2](https://github.com/wailsapp/wails/pull/2951) oleh @APshenkin.
- Perintah `wails3 generate template` oleh [@leaanthony](https://github.com/leaanthony)
- Perintah `wails3 releasenotes` oleh [@leaanthony](https://github.com/leaanthony)
- Perintah `wails3 update cli` oleh [@leaanthony](https://github.com/leaanthony)
- Opsi `-clean` untuk perintah `wails3 generate bindings` oleh [@leaanthony](https://github.com/leaanthony)
- Izinkan build AppImage Linux untuk aarch64 (arm64) oleh [@AkshayKalose](https://github.com/AkshayKalose) dalam [#3981](https://github.com/wailsapp/wails/pull/3981)
- Tambahkan hyperlink untuk sponsor oleh @ansxuman dalam [#3958](https://github.com/wailsapp/wails/pull/3958)
- Dukungan pemaketan Linux untuk build pemaket deb, rpm, dan Arch Linux oleh
- Tambahkan dukungan untuk build dan paket universal Darwin oleh
- Dokumentasi event pada situs web oleh
- Templat untuk sveltekit dan sveltekit-ts yang dikonfigurasi untuk pengembangan non-SSR
- Perbarui aset build menggunakan perintah `wails3 update build-assets` baru oleh
- Contoh untuk menguji HTML Drag and Drop API oleh
- Dukungan Asosiasi File oleh [leaanthony](https://github.com/leaanthony) dalam
- Perintah `wails3 generate runtime` baru oleh
- Opsi `InitialPosition` baru untuk menentukan apakah jendela harus dipusatkan atau
- Tambahkan metode `Path` & `Paths` ke paket `application` oleh
- Menambahkan opsi Windows `GeneralAutofillEnabled` dan `PasswordAutosaveEnabled`
- Menambahkan kemampuan untuk mengambil jendela yang memanggil metode layanan, oleh
- Menambahkan opsi `EnabledFeatures` dan `DisabledFeatures` untuk WebView2 oleh
- ⊞ Sistem DIP baru untuk Dukungan Monitor DPI Tinggi yang Disempurnakan oleh
- ⊞ Opsi nama kelas jendela oleh [windom](https://github.com/windom/) dalam
- Layanan telah diperluas untuk menyediakan fungsionalitas plugin. Oleh
- 🐧 Event WindowDidMove / WindowDidResize dalam
- ⊞ Event WindowDidResize dalam
-  tambahkan Event ApplicationShouldHandleReopen agar dapat menangani dock
-  tambahkan getPrimaryScreen/getScreens ke implementasi oleh @tmclane dalam
-  tambahkan opsi untuk menampilkan bilah alat dalam mode layar penuh di macOS oleh
- 🐧 tambahkan logika onKeyPress untuk mengonversi penekanan tombol Linux menjadi accelerator
- 🐧 tambahkan tugas `run:linux` oleh
- Ekspor metode `SetIcon` oleh [@almas-x](https://github.com/almas-x) dalam
- Tingkatkan `OnShutdown` oleh [@almas-x](https://github.com/almas-x) dalam
- Pulihkan metode `ToggleMaximise` dalam antarmuka `Window` oleh
- Menambahkan informasi lebih lanjut ke `Environment()`. Oleh @leaanthony dalam
- Ekspos metode `WebviewWindow.IsFocused` pada antarmuka `Window` oleh
- Dukung beberapa event pemicu yang dipisahkan spasi dalam sistem WML oleh
- Tambahkan ekspor ESM dari skrip runtime JS terbundel oleh
- Tambahkan flag generator binding untuk menggunakan skrip runtime JS terbundel sebagai pengganti
- Implementasikan `setIcon` di Linux oleh [@abichinger](https://github.com/abichinger)
- Tambahkan flag `-port` ke perintah dev dan dukung variabel lingkungan
- Tambahkan pengujian untuk pemanggilan metode terikat oleh
- ⊞ tambahkan `SetIgnoreMouseEvents` untuk jendela yang sudah dibuat oleh
-  Tambahkan kemampuan untuk mengatur tingkat penumpukan (urutan) jendela oleh

### Diperbaiki

- Perbaiki `Screen.Bounds`, `WorkArea`, dan `Size` yang nilainya menjadi separuh pada Mac Retina dengan mengonversi nilai titik NSScreen menjadi piksel perangkat dalam bidang `Physical*`, serta isi `Screen.X`/`Y` tingkat teratas agar deteksi monitor yang bersinggungan dan penempatan area kerja berfungsi dengan benar dalam [PR](https://github.com/wailsapp/wails/pull/5168) oleh @wayneforrest
- Perbaiki data race di ScreenManager yang menyebabkan deadlock WebKit DisplayLink saat konfigurasi layar berubah (misalnya, monitor eksternal dihubungkan saat proses tidur/bangun)
- Tetapkan CFBundleIconName secara langsung ke appicon jika Assets.car ada dalam [PR](https://github.com/wailsapp/wails/pull/5154) oleh @symball
- Perbaiki `wails3 doctor` yang melaporkan paket WebKitGTK yang salah di Fedora, openSUSE, Arch, dan NixOS — entri fallback 4.0 telah dihapus karena v3 memerlukan API 4.1 pada waktu kompilasi (#5071)
- Perbaiki nama paket webkit2gtk pada doctor untuk openSUSE (`webkit2gtk4_1-devel` → `webkit2gtk3-devel`, nama paket openSUSE yang benar) (#5071)
- Perbaiki galat `Unexpected token '<'` ketika `/wails/custom.js` tidak ada dalam mode pengembangan desktop. Tambahkan handler 404 eksplisit untuk `/wails/custom.js` dan validasi `Content-Type` tanpa membedakan huruf besar-kecil dalam `loadOptionalScript` guna mencegah fallback SPA HTML disisipkan sebagai JavaScript. ([#5068](https://github.com/wailsapp/wails/issues/5068))
- Perbaiki status sorotan menu baki sistem di macOS — ikon kini menampilkan status terpilih saat menu terbuka (#4910)
- Perbaiki jendela yang terpasang pada baki sistem dan muncul di belakang jendela lain di macOS — kini menggunakan tingkat jendela popup yang tepat (#4910)
- Perbaiki contoh impor `@wailsio/runtime` yang salah di seluruh dokumentasi (#4989)
- Perbaiki jendela tanpa bingkai yang tidak dapat diminimalkan di darwin (#4294)
- Memperbaiki hang selama 20-30 menit saat `wails3 build` dan `wails3 dev` dengan mengecualikan `node_modules/` dari pemeriksaan kemutakhiran go-task. Sebelumnya, glob `sources: "**/*"` menyebabkan go-task mendata setiap berkas dan menghitung checksum-nya dalam `node_modules/` (50000-100000+ berkas dengan dependensi berat seperti MUI), yang sangat lambat khususnya di Windows/NTFS (#4939)
- Perbaiki kegagalan build GTK4 akibat typedef C `Screen` yang bertabrakan dengan X11 Xlib.h (#4957)
- Perbaiki konsistensi metode lencana dock di macOS
- Perbaiki `InvisibleTitleBarHeight` yang diterapkan pada semua jendela macOS, bukan hanya jendela tanpa bingkai atau dengan bilah judul transparan (#4960)
- Perbaiki jendela yang berguncang/bergetar saat ukurannya diubah dari sudut atas dengan `InvisibleTitleBarHeight` diaktifkan, dengan melewati inisiasi seret di dekat tepi jendela (#4960)
- Perbaiki pembuatan tipe terpetakan dengan kunci enum dalam binding JS/TS (#4437) oleh @fbbdev
- Perbaiki seret dan lepas berkas di Windows yang tidak berfungsi pada penskalaan layar selain 100%
- Perbaiki seret dan lepas internal HTML5 yang rusak ketika pelepasan berkas diaktifkan di Windows
- Perbaiki koordinat pelepasan berkas yang berada dalam ruang piksel yang salah di Windows (piksel fisik versus piksel CSS)
- Perbaiki seret dan lepas berkas di Linux yang tidak berfungsi secara andal dengan efek hover
- Perbaiki seret dan lepas internal HTML5 yang rusak ketika pelepasan berkas diaktifkan di Linux
- Perbaiki tindakan menampilkan/menyembunyikan jendela di Linux/GTK4 yang terkadang memulihkan jendela ke status diminimalkan dengan menggunakan `gtk_window_present()` (#4957)
- Perbaiki operasi mendapatkan/mengatur posisi jendela di Linux/GTK4 yang selalu mengembalikan 0,0 dengan menambahkan dukungan bersyarat X11 melalui `XTranslateCoordinates`/`XMoveWindow` (#4957)
- Perbaiki ukuran maksimum jendela yang tidak diberlakukan di Linux/GTK4 dengan menambahkan pembatasan ukuran berbasis sinyal untuk menggantikan `gtk_window_set_geometry_hints` yang telah dihapus (#4957)
- Perbaiki penskalaan DPI di Linux/GTK4 dengan menerapkan penghitungan PhysicalBounds yang tepat dan dukungan penskalaan pecahan melalui `gdk_monitor_get_scale` (GTK 4.14+)
- Perbaiki item menu yang terduplikasi saat membuat jendela baru di Linux/GTK4
- Perbaiki pembuatan tipe terpetakan dengan kunci enum dalam binding JS/TS (#4437) oleh @fbbdev
- Perbaiki seret dan lepas berkas di Windows yang tidak berfungsi pada penskalaan layar selain 100%
- Perbaiki seret dan lepas internal HTML5 yang rusak ketika pelepasan berkas diaktifkan di Windows
- Perbaiki koordinat pelepasan berkas yang berada dalam ruang piksel yang salah di Windows (piksel fisik versus piksel CSS)
- Perbaiki seret dan lepas berkas di Linux yang tidak berfungsi secara andal dengan efek hover
- Perbaiki seret dan lepas internal HTML5 yang rusak ketika pelepasan berkas diaktifkan di Linux
- Perbaiki penskalaan DPI di Linux/GTK4 dengan menerapkan penghitungan PhysicalBounds yang tepat dan dukungan penskalaan pecahan melalui `gdk_monitor_get_scale` (GTK 4.14+)
- Perbaiki item menu yang terduplikasi saat membuat jendela baru di Linux/GTK4
- Perbaiki pembuatan tipe terpetakan dengan kunci enum dalam binding JS/TS (#4437) oleh @fbbdev
- Perbaiki masalah "jendela hantu" di macOS akibat AppKit API tidak diakses dari Main Thread dalam App.Window.Current() (#4947) oleh @wimaha
- Perbaiki `<input type="file">` HTML yang tidak berfungsi di macOS dengan mengimplementasikan WKUIDelegate runOpenPanelWithParameters (#4862)
- Perbaiki seret dan lepas berkas native yang tidak berfungsi saat menggunakan modul npm `@wailsio/runtime` di macOS/Linux (#4953) oleh @leaanthony
- Perbaiki pembuatan binding untuk alias tipe lintas paket (#4578) oleh @fbbdev
- Perbaiki crash OpenFileDialog di Linux akibat pelanggaran keamanan thread GTK (#3683) oleh @ddmoney420
- Perbaiki crash SIGSEGV saat memanggil `Focus()` pada jendela yang disembunyikan atau telah dihancurkan (#4890) oleh @ddmoney420
- Perbaiki potensi panic saat mengatur ikon atau bitmap kosong di Linux (#4923) oleh @ddmoney420
- Perbaiki crash ErrorDialog saat dipanggil dari binding layanan di macOS (#3631) oleh @leaanthony
- Tampilkan menu di sistem operasi Windows dalam `v3\examples\dialogs` oleh @ndianabasi
- Perbaiki race condition yang menyebabkan TypeError saat halaman dimuat ulang (#4872) oleh @ddmoney420
- Perbaiki keluaran yang salah dari pengujian generator binding dengan menghapus status global dalam metode `Collector.IsVoidAlias()` (#4941) oleh @fbbdev
- Perbaiki pemilih berkas `<input type="file">` yang tidak berfungsi di macOS (#4862) oleh @leaanthony
- Perbaiki `Position()` dan `SetPosition()` yang menggunakan sistem koordinat tidak konsisten di macOS, sehingga posisi jendela bergeser saat menyimpan/memulihkan status (#4816) oleh @leaanthony
- Perbaiki galat SetProcessDpiAwarenessContext "Akses ditolak" ketika kesadaran DPI sudah ditetapkan melalui manifes aplikasi (#4803)
- Perbarui halaman dokumentasi pintasan papan ketik dan koreksi tipe parameter callback untuk `KeyBinding.Add` oleh @ndianabasi
- Perbaiki dokumentasi tentang pembuatan binding khusus; harus menggunakan `-d String`, bukan `-o String`
- Perbaiki menu yang tidak menghapus elemen anak saat `menu.Update()`
- Perbaiki referensi Manager API yang kedaluwarsa dalam dokumentasi (31 berkas diperbarui agar menggunakan pola baru seperti `app.Window.New()`, `app.Event.Emit()`, dan seterusnya) oleh @leaanthony
- Perbaiki crash Linux saat terjadi panic dalam metode Go yang terikat ke JS akibat WebKit menimpa penangan sinyal (#3965) oleh @leaanthony
- Perbaiki SaveFileDialog.SetFilename() yang tidak berpengaruh di Linux (#4841) oleh @samstanier
- Perbaiki koordinat pelepasan yang ditampilkan sebagai undefined dalam contoh seret dan lepas
- Perbaiki kegagalan pembuatan bundel aplikasi macOS ketika APP_NAME berisi spasi (masalah ekspansi kurung kurawal)
- Perbaiki panic indeks di luar batas pada Windows saat memanggil metode layanan (batalkan penggunaan goccy/go-json)
- Perbaiki seret dan lepas file pada Windows yang tidak berfungsi pada penskalaan tampilan selain 100%
- Perbaiki seret dan lepas internal HTML5 yang rusak ketika pelepasan file diaktifkan pada Windows
- Perbaiki koordinat pelepasan file yang menggunakan ruang piksel keliru pada Windows (piksel fisik dibandingkan piksel CSS)
- Perbaiki seret dan lepas file pada Linux yang tidak berfungsi secara andal dengan efek hover
- Perbaiki seret dan lepas internal HTML5 yang rusak ketika pelepasan file diaktifkan pada Linux
- Perbarui semua perintah dalam file Taskfile.yml untuk seluruh sistem operasi agar dapat menangani spasi dalam variabel seperti `APP_NAME` oleh @ndianabasi
- Perbaiki galat argumen perintah saat menjalankan tugas 'build:universal:lipo:go' di Linux oleh @wux1an
- Perbaiki galat Docker "undefined symbol: **<em>ubsan</em>handle_xxxxxxx" saat menjalankan 'wails3 build GOOS=darwin GOARCH=arm64' di Linux oleh @wux1an
- Satukan dokumentasi protokol khusus dan tambahkan bagian Universal Links oleh @leaanthony
- Perbaiki crash menu baki sistem Windows saat ikon diklik berulang kali dengan menambahkan pelindung terhadap pemanggilan TrackPopupMenuEx secara bersamaan (#4151) oleh @leaanthony
- Cegah aplikasi mengalami crash saat systray.Run() dipanggil sebelum app.Run() oleh @leaanthony
- Perbaiki crash pada macOS saat visibilitas jendela dialihkan melalui Hide()/Show() dengan ApplicationShouldTerminateAfterLastWindowClosed diaktifkan (#4389) oleh @leaanthony
- Perbaiki kebocoran memori pada menu konteks di macOS dan Windows ketika menu dibuka berulang kali (#4012) oleh @leaanthony
- Perbaiki sumber daya native menu konteks yang tidak digunakan kembali pada macOS, sehingga menu baru dibuat setiap kali ditampilkan (#4012) oleh @leaanthony
- Perbaiki klik ikon dock macOS yang tidak menampilkan jendela tersembunyi ketika aplikasi dimulai dengan `Hidden: true` (#4583) oleh @leaanthony
- Perbaiki dialog cetak macOS yang tidak terbuka akibat tipe pointer jendela yang keliru dalam pemanggilan CGO (#4290) oleh @leaanthony
- Perbaiki crash menu jendela pada Wayland akibat appmenu-gtk-module mengakses jendela yang belum direalisasikan (#4769) oleh @leaanthony
- Perbaiki crash aplikasi GTK ketika nama aplikasi berisi karakter yang tidak valid (spasi, tanda kurung, dan sebagainya) oleh @leaanthony
- Perbaiki galat "memori tidak cukup" saat menginisialisasi seret dan lepas pada Windows (#4701) oleh @overlordtm
- Perbaiki penjelajah file yang membuka direktori keliru pada Linux akibat escaping URI yang tidak tepat (#4397) oleh @leaanthony
- Perbaiki kegagalan build AppImage pada distribusi Linux modern (Arch, Fedora 39+, Ubuntu 24.04+) dengan mendeteksi otomatis bagian ELF `.relr.dyn` dan menonaktifkan stripping (#4642) oleh @leaanthony
- Perbaiki `wails doctor` yang secara keliru melaporkan paket webkit sebagai telah terpasang pada Fedora/sistem berbasis DNF (#4457) oleh @leaanthony
- Perbaiki `config.yml` default yang menjalankan `wails3 dev` dengan build produksi oleh @mbaklor
- Perbaiki stub layanan iOS yang menyebabkan kegagalan build akibat mengimpor paket yang tidak ada oleh @leaanthony
- Perbaiki pencatatan terstruktur dalam metode debug/info yang menyebabkan galat "tidak ada direktif pemformatan" oleh @leaanthony
- Hapus pernyataan cetak debug sementara yang tidak sengaja disertakan dari penggabungan platform seluler oleh @leaanthony
- Perbaiki crash WebKitGTK pada Wayland dengan GPU NVIDIA (Galat 71 Protocol error) dengan menonaktifkan perender DMA-BUF secara otomatis oleh @leaanthony
- Atasi nilai alfa yang diabaikan dalam `application.WebviewWindowOptions.BackgroundColour` pada Linux ([#4722](https://github.com/wailsapp/wails/pull/4722), @BradHacker)
- Perbaiki ikon baki sistem Windows yang tidak menggunakan ikon aplikasi secara default ketika tidak ada ikon khusus yang diberikan (#4704)
- Lacak kepemilikan `HICON` agar hanya handle buatan pengguna yang dimusnahkan, sehingga mencegah crash saat Explorer melakukan daur ulang (#4653).
- Lepaskan listener tema sistem Windows dan ikon baki yang dipertahankan saat pemusnahan untuk menghentikan kebocoran goroutine dan konteks perangkat (#4653).
- Potong tooltip baki pada 127 unit UTF-16 agar pasangan surrogate dan glif multibita tidak rusak (#4653).
- Perbaiki kegagalan tugas paket Windows (#4667)
- Perbaiki variabel appicon AppImage Linux dalam taskfile Linux [PR #4644](https://github.com/wailsapp/wails/pull/4644)
- Perbaiki galat build Windows akibat perubahan signature go-webview2 v1.0.22 (#4513, #4645)
- Perbaiki variabel appicon AppImage Linux dalam taskfile Linux [PR #4644](https://github.com/wailsapp/wails/pull/4644)
- Memperbaiki rentang protokol desktop.tmpl Linux dengan mengubah `<.Info.Protocol>` menjadi `<.Protocol>` oleh @Tolfx dalam #4510
- Perbaiki galat redefinisi untuk demo liquid glass dalam [#4542](https://github.com/wailsapp/wails/pull/4542) oleh @Etesam913
- Perbaiki pembaruan menu baki sistem pada Linux [#4604](https://github.com/wailsapp/wails/issues/4604) oleh [@JackDoan](https://github.com/JackDoan)
- Perbaiki jendela putih yang muncul pada Windows saat membuat jendela tersembunyi oleh @leaanthony dalam [#4612](https://github.com/wailsapp/wails/pull/4612)
- Perbaiki jalur impor paket notifikasi dalam dokumentasi oleh @rxliuli dalam [#4617](https://github.com/wailsapp/wails/pull/4617)
- Memperbaiki fungsi seret dan lepas yang tidak berfungsi saat menggunakan paket npm @wailsio/runtime (#4489), oleh @leaanthony dalam #4616
- Windows: Memperbaiki jendela yang berkedip saat dimulai dan jendela tersembunyi yang ditampilkan secara tidak semestinya dalam [PR](https://github.com/wailsapp/wails/pull/4600) oleh @leaanthony.
- Memperbaiki masalah ukuran jendela saat dimaksimalkan di Wayland (https://github.com/wailsapp/wails/issues/4429), oleh [@samstanier](https://github.com/samstanier)
- Memperbaiki masalah ukuran jendela saat dimaksimalkan di Wayland (https://github.com/wailsapp/wails/issues/4429), oleh [@samstanier](https://github.com/samstanier)
- Memperbaiki kesalahan redefinisi pada demo liquid glass dalam [#4542](https://github.com/wailsapp/wails/pull/4542), oleh @Etesam913
- Memperbaiki masalah yang dapat menyebabkan AssetServer mengalami crash di MacOS dalam [#4576](https://github.com/wailsapp/wails/pull/4576), oleh @jghiloni
- Memperbaiki masalah kompilasi saat membangun dengan NextJs. Diperbaiki dalam [#4585](https://github.com/wailsapp/wails/pull/4585) oleh @rev42
- Memperbaiki pipeline untuk rilis nightly dalam [#4597](https://github.com/wailsapp/wails/pull/4597), oleh @riadafridishibly
- Memperbaiki kesalahan redefinisi pada demo liquid glass dalam [#4542](https://github.com/wailsapp/wails/pull/4542), oleh @Etesam913
- Memperbaiki masalah yang dapat menyebabkan AssetServer mengalami crash di MacOS dalam [#4576](https://github.com/wailsapp/wails/pull/4576), oleh @jghiloni
- Memperbaiki masalah kompilasi saat membangun dengan NextJs. Diperbaiki dalam [#4585](https://github.com/wailsapp/wails/pull/4585) oleh @rev42
- Memperbaiki pipeline untuk rilis nightly dalam [#4597](https://github.com/wailsapp/wails/pull/4597), oleh @riadafridishibly
- Memperbaiki kesalahan redefinisi pada demo liquid glass dalam [#4542](https://github.com/wailsapp/wails/pull/4542), oleh @Etesam913
- Memperbaiki SetBackgroundColour di Windows, oleh @PPTGamer dalam [PR](https://github.com/wailsapp/wails/pull/4492)
- Memperbarui dokumentasi agar mencerminkan perubahan dari pemfaktoran ulang Manager API, oleh @yulesxoxo dalam [PR #4476](https://github.com/wailsapp/wails/pull/4476)
- Memperbaiki variabel appicon pada berkas .desktop Linux di Taskfile Linux [PR #4477](https://github.com/wailsapp/wails/pull/4477)
- Memperbarui dokumentasi agar mencerminkan perubahan dari pemfaktoran ulang Manager API, oleh @yulesxoxo dalam [PR #4476](https://github.com/wailsapp/wails/pull/4476)
- Memperbaiki bug dereferensi pointer nil di Windows yang dilaporkan dalam [#4456](https://github.com/wailsapp/wails/issues/4456), oleh @leaanthony dalam [#4460](https://github.com/wailsapp/wails/pull/4460)
- Menambahkan dukungan untuk `allowsBackForwardNavigationGestures` di WKWebView macOS guna mengaktifkan gestur navigasi usap dua jari (#1857)
- Memperbaiki masalah onClick yang tidak berfungsi untuk item menu yang pada awalnya ditetapkan sebagai dinonaktifkan, oleh @leaanthony dalam [PR #4469](https://github.com/wailsapp/wails/pull/4469). Terima kasih kepada @IanVS atas investigasi awalnya.
- Memperbaiki server Vite yang tidak dibersihkan ketika proses build gagal (#4403)
- Memperbaiki panic saat menutup atau membatalkan `SaveFileDialog` di Windows. Diperbaiki dalam [PR](https://github.com/wailsapp/wails/pull/4284) oleh @hkhere
- Memperbaiki fungsi seret dan lepas pada tingkat HTML di Windows, oleh [@mbaklor](https://github.com/mbaklor) dalam [#4259](https://github.com/wailsapp/wails/pull/4259)
- Menambahkan dukungan untuk `allowsBackForwardNavigationGestures` di WKWebView macOS guna mengaktifkan gestur navigasi usap dua jari (#1857)
- Memperbaiki masalah onClick yang tidak berfungsi untuk item menu yang pada awalnya ditetapkan sebagai dinonaktifkan, oleh @leaanthony dalam [PR #4469](https://github.com/wailsapp/wails/pull/4469). Terima kasih kepada @IanVS atas investigasi awalnya.
- Memperbaiki server Vite yang tidak dibersihkan ketika proses build gagal (#4403)
- Memperbaiki penguraian notifikasi di Windows, oleh @popaprozac dalam [PR](https://github.com/wailsapp/wails/pull/4450)
- Memperbaiki perintah doctor agar memeriksa dependensi Windows SDK, oleh [@kodumulo](https://github.com/kodumulo) dalam [#4390](https://github.com/wailsapp/wails/issues/4390)
- Memperbaiki dereferensi pointer nil dalam processURLRequest untuk Mac, oleh [@etesam913](https://github.com/etesam913) dalam [#4366](https://github.com/wailsapp/wails/pull/4366)
- Memperbaiki bug Linux yang membuat dialog dengan filter tidak dapat digunakan, oleh [@bh90210](https://github.com/bh90210) dalam [#4287](https://github.com/wailsapp/wails/pull/4287)
- Memperbaiki masalah menu Edit di Windows dan Linux, oleh [@leaanthony](https://github.com/leaanthony) dalam [#3f78a3a](https://github.com/wailsapp/wails/commit/3f78a3a8ce7837e8b32242c8edbbed431c68c062)
- Memperbarui versi sistem minimum dalam berkas .plist macOS dari 10.13.0 menjadi 10.15.0, oleh [@AkshayKalose](https://github.com/AkshayKalose) dalam [#3981](https://github.com/wailsapp/wails/pull/3981)
- Memperbaiki masalah ID jendela yang terlewati, oleh [@leaanthony](https://github.com/leaanthony)
- Memperbaiki masalah menu nil saat memanggil RegisterContextMenu, oleh [@leaanthony](https://github.com/leaanthony)
- Memperbaiki siklus dependensi dalam keluaran generator binding, oleh [@fbbdev](https://github.com/fbbdev) dalam [#4001](https://github.com/wailsapp/wails/pull/4001)
- Memperbaiki kesalahan penggunaan sebelum definisi dalam keluaran generator binding, oleh [@fbbdev](https://github.com/fbbdev) dalam [#4001](https://github.com/wailsapp/wails/pull/4001)
- Meneruskan flag build ke generator binding, oleh [@fbbdev](https://github.com/fbbdev) dalam [#4023](https://github.com/wailsapp/wails/pull/4023)
- Mengubah garis miring pada path di Taskfile Windows menjadi garis miring ke depan agar berfungsi pada platform selain Windows, oleh [@leaanthony](https://github.com/leaanthony)
- Memperbaiki event Mac dan event JS Mac, oleh [@leaanthony](https://github.com/leaanthony)
- Memperbaiki deadlock event di macOS, oleh [@leaanthony](https://github.com/leaanthony)
- Memperbaiki kesalahan `Parameter incorrect` saat inisialisasi Window di Windows ketika HTML disediakan tetapi JS tidak tersedia, oleh [@leaanthony](https://github.com/leaanthony)
- Memperbaiki ukuran prefiks respons yang digunakan untuk mendeteksi jenis konten di asset server, oleh [@fbbdev](https://github.com/fbbdev) dalam [#4049](https://github.com/wailsapp/wails/pull/4049)
- Memperbaiki penanganan respons non-404 pada path indeks root di asset server, oleh [@fbbdev](https://github.com/fbbdev) dalam [#4049](https://github.com/wailsapp/wails/pull/4049)
- Memperbaiki perilaku tak terdefinisi dalam generator binding saat menguji properti tipe generik, oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Memperbaiki keluaran generator binding untuk model ketika tipe dasarnya tidak memiliki properti yang sama dengan wrapper bernama, oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Memperbaiki keluaran generator binding untuk tipe kunci map dan prapemrosesan, oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Memperbaiki keluaran generator binding untuk struct yang mengimplementasikan antarmuka marshaler, oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Memperbaiki deteksi siklus tipe yang melibatkan tipe generik dalam generator binding, oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Memperbaiki referensi yang tidak valid ke model yang tidak diekspor dalam keluaran generator binding oleh [@fbbdev](https://github.com/fbbdev) di [#4045](https://github.com/wailsapp/wails/pull/4045)
- Memindahkan kode yang diinjeksi ke akhir file layanan oleh [@fbbdev](https://github.com/fbbdev) di [#4045](https://github.com/wailsapp/wails/pull/4045)
- Memperbaiki penanganan galat dari operasi penutupan file dalam generator binding oleh [@fbbdev](https://github.com/fbbdev) di [#4045](https://github.com/wailsapp/wails/pull/4045)
- Menonaktifkan peringatan untuk layanan yang mendefinisikan metode siklus hidup atau HTTP, tetapi tidak memiliki metode terikat lainnya, oleh [@fbbdev](https://github.com/fbbdev) di [#4045](https://github.com/wailsapp/wails/pull/4045)
- Memperbaiki templat non-React yang gagal menampilkan footer Hello World saat menggunakan skema warna sistem terang oleh [@marcus-crane](https://github.com/marcus-crane) di [#4056](https://github.com/wailsapp/wails/pull/4056)
- Memperbaiki item menu tersembunyi di macOS oleh [@leaanthony](https://github.com/leaanthony)
- Memperbaiki penanganan dan pemformatan galat dalam pemroses pesan oleh [@fbbdev](https://github.com/fbbdev) di [#4066](https://github.com/wailsapp/wails/pull/4066)
-  Memperbaiki penghentian layanan yang terlewat saat keluar dari aplikasi oleh [@fbbdev](https://github.com/fbbdev) di [#4066](https://github.com/wailsapp/wails/pull/4066)
-  Memastikan pembaruan menu dilakukan di thread utama oleh [@leaanthony](https://github.com/leaanthony)
- Mekanisme penyeretan dan pengubahan ukuran kini lebih tangguh dan lebih sesuai dengan perilaku platform yang diharapkan, oleh [@fbbdev](https://github.com/fbbdev) di [#4100](https://github.com/wailsapp/wails/pull/4100)
- Memperbaiki [#4097](https://github.com/wailsapp/wails/issues/4097): Webpack/Angular membuang kode inisialisasi runtime, oleh [@fbbdev](https://github.com/fbbdev) di [#4100](https://github.com/wailsapp/wails/pull/4100)
- Memperbaiki masalah pada item menu yang awalnya tersembunyi; perbaikan oleh [@IanVS](https://github.com/IanVS) dalam [#4116](https://github.com/wailsapp/wails/pull/4116)
- Memperbaiki assetFileServer yang tidak menyajikan file `.html` untuk permintaan tanpa ekstensi ketika `[request]` tidak ada, tetapi `[request].html` ada
- Memperbaiki jalur pembuatan ikon oleh [@robin-samuel](https://github.com/robin-samuel) di [#4125](https://github.com/wailsapp/wails/pull/4125)
- Memperbaiki masalah event `fullscreen`, `unfullscreen`, `unminimise`, dan `unmaximise` yang tidak dipancarkan; perbaikan oleh [@oSethoum](https://github.com/osethoum) dalam [#4130](https://github.com/wailsapp/wails/pull/4130)
- Memperbaiki Galat NSIS akibat prefiks yang salah pada versi bawaan dalam konfigurasi oleh [@robin-samuel](https://github.com/robin-samuel) di [#4126](https://github.com/wailsapp/wails/pull/4126)
- Memperbaiki fungsi runtime Dialogs yang mengembalikan jalur dengan karakter escape di Windows oleh [TheGB0077](https://github.com/TheGB0077) di [#4188](https://github.com/wailsapp/wails/pull/4188)
- Memperbaiki jalur deteksi Webview2 di HKCU oleh [@leaanthony](https://github.com/leaanthony).
- Memperbaiki masalah input di macOS oleh [@leaanthony](https://github.com/leaanthony).
- Memperbaiki nama file task pembuatan ikon Windows oleh [@yulesxoxo](https://github.com/yulesxoxo) di [#4219](https://github.com/wailsapp/wails/pull/4219).
- Memperbaiki masalah transparansi pada jendela tanpa bingkai oleh [@leaanthony](https://github.com/leaanthony) berdasarkan pekerjaan @kron.
- Memperbaiki pemanggilan fokus saat jendela dinonaktifkan atau diminimalkan; perbaikan oleh [@leaanthony](https://github.com/leaanthony) berdasarkan pekerjaan @kron.
- Memperbaiki ikon baki sistem yang tidak muncul setelah taskbar dimulai ulang; perbaikan oleh [@leaanthony](https://github.com/leaanthony) berdasarkan pekerjaan @kron.
- Memperbaiki fallbackResponseWriter yang tidak mengimplementasikan Flush() di [#4245](https://github.com/wailsapp/wails/pull/4245)
- Memperbaiki fallbackResponseWriter yang tidak mengimplementasikan Flush() oleh [@superDingda] di [#4236](https://github.com/wailsapp/wails/issues/4236)
- Memperbaiki crash saat menutup jendela macOS ketika masih ada pemanggilan fungsi asinkron yang terikat ke Go dan tertunda, oleh [@joshhardy](https://github.com/joshhardy) di [#4354](https://github.com/wailsapp/wails/pull/4354)
- Memperbaiki kondisi balapan saat memulai mode Efisiensi Windows oleh [@leaanthony](https://github.com/leaanthony)
- Memperbaiki pembersihan handle ikon Windows oleh [@leaanthony](https://github.com/leaanthony).
- Memperbaiki `OpenFileManager` di Windows oleh [@PPTGamer](https://github.com/PPTGamer) di [#4375](https://github.com/wailsapp/wails/pull/4375).
- Memperbaiki opsi lebar minimum/maksimum untuk Linux oleh @atterpac di [#3979](https://github.com/wailsapp/wails/pull/3979)
- Memperbaiki definisi tipe pada templat TypeScript melalui peningkatan versi npm oleh @atterpac di [#3966](https://github.com/wailsapp/wails/pull/3966)
- Memperbaiki referensi CSS templat SvelteKit oleh @atterpac di [#3945](https://github.com/wailsapp/wails/pull/3945)
- Memastikan callback utama dalam run() jendela dipanggil di thread utama oleh [@leaanthony](https://github.com/leaanthony)
- Memperbaiki contoh pemilih direktori dialog oleh [@leaanthony](https://github.com/leaanthony)
- Membuat halaman galat baru dalam bahasa Tionghoa ketika index.html tidak ada, oleh [@leaanthony](https://github.com/leaanthony)
-  Memastikan callback `windowDidBecomeKey` berjalan di thread utama oleh [@leaanthony](https://github.com/leaanthony)
-  Menambahkan dukungan layar penuh untuk jendela tanpa bingkai oleh [@leaanthony](https://github.com/leaanthony)
-  Meningkatkan logika penghancuran jendela oleh [@leaanthony](https://github.com/leaanthony)
-  Memperbaiki logika posisi jendela saat ditautkan ke baki sistem oleh [@leaanthony](https://github.com/leaanthony)
-  Menambahkan dukungan layar penuh untuk jendela tanpa bingkai oleh [@leaanthony](https://github.com/leaanthony)
- Memperbaiki penanganan event oleh [@leaanthony](https://github.com/leaanthony)
- Memperbaiki logika penghentian jendela oleh [@leaanthony](https://github.com/leaanthony)
- Taskfile umum kini secara bawaan menghasilkan binding TypeScript untuk templat TypeScript, oleh [@leaanthony](https://github.com/leaanthony)
- Memperbaiki penutupan aplikasi saat menerima pesan WM_CLOSE ketika tidak ada jendela yang terbuka/aplikasi hanya menggunakan baki sistem, oleh [@mmalcek](https://github.com/mmalcek) di [#3990](https://github.com/wailsapp/wails/pull/3990)
- Memperbaiki build garble oleh @5aaee9 di [#3192](https://github.com/wailsapp/wails/pull/3192)
- Memperbaiki build NSIS Windows oleh [@leaanthony](https://github.com/leaanthony)
- Memperbaiki deadlock dalam dialog Linux untuk beberapa pilihan yang disebabkan oleh sesuatu yang tidak ditutup
- Memperbaiki pembersihan lintas platform untuk file .syso selama build Windows oleh
- Memperbaiki kompilasi AppImage amd64 oleh @atterpac di
- Memperbaiki pembaruan aset build oleh @ansxuman di
- Memperbaiki implementasi `OnClick` dan `OnRightClick` pada systray Linux oleh @atterpac
- Memperbaiki `AlwaysOnTop` yang tidak berfungsi di Mac oleh
-  Memperbaiki `application.NewEditMenu` yang menyertakan duplikat
- 🐧 Memperbaiki kompilasi aarch64
- ⊞ Memperbaiki item menu grup radio oleh
- Memperbaiki kesalahan saat membangun .app yang dapat dijalankan di MacOS ketika 'name' dan 'outputfilename'
- Memperbaiki bug penggunaan customEventProcessor dalam contoh drag-n-drop oleh
- 🐧 Memperbaiki kesalahan kompilasi Linux yang muncul akibat penambahan IgnoreMouseEvents oleh
- ⊞ Memperbaiki bug pembuatan berkas ikon syso oleh
- 🐧 Memperbaiki agar dapat berjalan secara native di Wayland, diintegrasikan dari
- Jangan ikat metode layanan internal di
- ⊞ Memperbaiki panic saat systray dimulai di
- Jangan ikat metode layanan internal di
- ⊞ Memperbaiki panic saat systray dimulai di
- Refaktor besar pada item menu dan penanganan peristiwa. Untuk saat ini, terutama meningkatkan macOS. Oleh
- Memperbaiki pengujian setelah refaktor plugin dan peristiwa di
- ⊞ Memperbaiki peringatan `Failed to unregister class Chrome_WidgetWin_0`. Oleh
- Masalah modul
- Memperbaiki pengiriman pesan peristiwa pengubahan ukuran oleh [atterpac](https://github.com/atterpac) di
- 🐧 Memperbaiki kesalahan penanganan tema di NixOS oleh
- Memperbaiki instalasi proyek lintas volume untuk Windows oleh
- Memperbaiki CSS templat React agar footer ditampilkan oleh
- Memperbaiki proses zombi saat bekerja dalam devmode dengan memperbarui refresh ke versi terbaru
- Memperbaiki pengambilan sumber berkas WebKit AppImage oleh [Atterpac](https://github.com/atterpac)
- Memperbaiki verifikasi paket apt oleh Doctor, oleh [Atterpac](https://github.com/Atterpac) di
- Memperbaiki aplikasi yang membeku saat ditutup (Darwin) oleh @5aaee9 di
- Memperbaiki warna latar belakang contoh di Windows oleh
- Memperbaiki menu konteks bawaan oleh [mmghv](https://github.com/mmghv) di
- Memperbaiki nilai heksadesimal untuk tombol panah di Darwin oleh
- Membuat drag-n-drop berfungsi di Windows. Ditambahkan oleh
- Memperbaiki bug Linux di Doctor ketika pengguna tidak memiliki driver yang sesuai
- Memperbaiki penskalaan DPI saat dimulai (Windows). Diubah oleh [@almas-x](https://github.com/almas-x) di
- Memperbaiki baris penggantian di `go.mod` agar menggunakan jalur relatif. Memperbaiki jalur Windows dengan
- Memperbaiki penanganan klik systray MacOS ketika tidak ada jendela yang terpasang oleh
- Memperbaiki kegagalan build Windows akibat opsi yang tidak dikenal oleh
- Memperbaiki crash di Windows saat mengeklik kiri ikon systray ketika tidak memiliki
- Memperbaiki baseURL yang salah ketika jendela dibuka dua kali oleh @5aaee9 dalam PR
- Memperbaiki urutan cabang if dalam metode `WebviewWindow.Restore` oleh
- Menghitung `startURL` dengan benar pada beberapa pemanggilan `GetStartURL` ketika
- Memperbaiki tipe JS struct `Screen` agar sesuai dengan padanannya di Go oleh
- Memperbaiki metode `WML.Reload` untuk memastikan pembersihan yang tepat atas peristiwa yang terdaftar
- Memperbaiki menu konteks khusus yang langsung tertutup di Linux oleh
- Memperbaiki jalur keluaran dan ekstensi berkas model yang dihasilkan oleh binding
- Memperbaiki jalur impor berkas model dalam kode JS yang dihasilkan oleh binding
- Memperbaiki drag-n-drop pada beberapa distro Linux oleh
- Memperbaiki task yang hilang untuk macOS ketika menggunakan `wails3 task dev` oleh
- Memperbaiki pendaftaran peristiwa yang menyebabkan penetapan pada map nil oleh
- Memperbaiki unmarshaling parameter metode terikat oleh
- Memperbaiki penanganan beberapa nilai kembalian dari metode terikat oleh
- Memperbaiki deteksi npm oleh Doctor ketika npm tidak diinstal dengan pengelola paket sistem
- Memperbaiki MicrosoftEdgeWebview2Setup.exe yang hilang. Terima kasih kepada
- Memperbaiki crash acak di Linux akibat penanganan ID jendela oleh @leaanthony. Berdasarkan
- Memperbaiki systemTray.setIcon yang menyebabkan crash di Linux oleh
- Memastikan bingkai jendela diterapkan pada pemanggilan pertama dalam fungsi `setFrameless` di

### Diubah

- **PERUBAHAN TIDAK KOMPATIBEL**: Kunci map dalam binding JS/TS yang dihasilkan kini ditandai sebagai opsional agar mencerminkan semantik map Go secara akurat. Akses nilai map di Typescript kini mengembalikan `T | undefined`, bukan `T`, sehingga memerlukan pemeriksaan null atau assertion (#4943) oleh `@fbbdev`
- Mengubah penggunaan `Event` menjadi `Events` sesuai dengan perubahan pada `@wailsio/runtime`, serta menyesuaikan pemanggilan fungsi dalam dokumentasi di `Features/Events/Event System` oleh @AbdelhadiSeddar
- Memindahkan `EnabledFeatures`, `DisabledFeatures`, dan `AdditionalBrowserArgs` dari opsi per jendela ke `Options.Windows` tingkat aplikasi (#4559) oleh @leaanthony
- Memperbarui README untuk contoh `Drag N Drop` dan menegaskan bahwa `Internal Drag and Drop` didemonstrasikan melalui contoh tersebut oleh @ndianabasi
- Mengubah berbagai log debug dari Info menjadi Debug (oleh @mbaklor)
- **PERUBAHAN TIDAK KOMPATIBEL:** Mengganti nama `EnableDragAndDrop` menjadi `EnableFileDrop` dalam opsi jendela
- **PERUBAHAN TIDAK KOMPATIBEL:** Mengganti nama `DropZoneDetails` menjadi `DropTargetDetails` dalam konteks peristiwa
- **PERUBAHAN TIDAK KOMPATIBEL:** Mengganti nama metode `DropZoneDetails()` menjadi `DropTargetDetails()` pada `WindowEventContext`
- **PERUBAHAN TIDAK KOMPATIBEL:** Menghapus peristiwa `WindowDropZoneFilesDropped`; gunakan `WindowFilesDropped` sebagai gantinya
- **PERUBAHAN YANG MEMUTUS KOMPATIBILITAS:** Ubah atribut HTML dari `data-wails-dropzone` menjadi `data-file-drop-target`
- **PERUBAHAN YANG MEMUTUS KOMPATIBILITAS:** Ubah kelas hover CSS dari `wails-dropzone-hover` menjadi `file-drop-target-active`
- **PERUBAHAN YANG MEMUTUS KOMPATIBILITAS:** Hapus opsi `DragEffect`, `OnEnterEffect`, dan `OnOverEffect` dari Windows (sebelumnya merupakan bagian dari IDropTarget yang telah dihapus)
- Beralih ke goccy/go-json untuk seluruh pemrosesan JSON saat runtime (binding metode, peristiwa, permintaan webview, notifikasi, kvstore), sehingga meningkatkan performa sebesar 21-63% dan mengurangi alokasi memori sebesar 40-60%
- Optimalkan tata letak struct BoundMethod dan cache flag isVariadic untuk mengurangi overhead setiap panggilan
- Gunakan buffer argumen yang dialokasikan di stack untuk metode dengan `<=8` argumen guna menghindari alokasi heap
- Optimalkan pengumpulan hasil dalam panggilan metode untuk menghindari alokasi slice bagi nilai kembalian tunggal
- Gunakan sync.Map untuk cache tipe MIME guna meningkatkan performa konkurensi
- Gunakan pool buffer untuk membaca isi permintaan transport HTTP
- Alokasikan channel CloseNotify secara lazy dalam pendeteksi tipe konten untuk mengurangi alokasi per permintaan
- Hapus pencatatan log debug CSS dari server aset
- Perluas peta ekstensi tipe MIME agar mencakup 50+ format web umum (font, audio, video, dan sebagainya)
- Perbarui dokumentasi untuk opsi Window `X/Y` oleh @ruhuang2001
- Perbarui dokumentasi `Frontend Runtime` dengan menambahkan lebih banyak opsi untuk menghasilkan binding frontend oleh @ndianabasi
- Perbarui halaman dokumentasi Asset Server Wails v3 oleh @ndianabasi
- **PERUBAHAN YANG MEMUTUS KOMPATIBILITAS**: Hapus fungsi dialog tingkat paket (`application.InfoDialog()`, `application.QuestionDialog()`, dan sebagainya). Sebagai gantinya, gunakan pengelola `app.Dialog`: `app.Dialog.Info()`, `app.Dialog.Question()`, `app.Dialog.Warning()`, `app.Dialog.Error()`, `app.Dialog.OpenFile()`, `app.Dialog.SaveFile()`
- Perbarui dokumentasi dialog agar sesuai dengan API yang sebenarnya: gunakan `app.Dialog.*`, `AddButton()` dengan callback (bukan `SetButtons()`), `SetDefaultButton(*Button)` (bukan string), `AddFilter()` (bukan `SetFilters()`), `SetFilename()` (bukan `SetDefaultFilename()`), serta `app.Dialog.OpenFile().CanChooseDirectories(true)` untuk memilih folder
- **PERUBAHAN YANG MEMUTUS KOMPATIBILITAS**: Build produksi kini menjadi default. Untuk membuat build pengembangan, tetapkan `DEV=true` dalam Taskfile Anda. Buat proyek baru untuk melihat contohnya. Oleh @leaanthony
- Saat memancarkan peristiwa khusus dengan nol atau satu argumen data, nilai data akan langsung ditetapkan ke field Data tanpa membungkusnya dalam slice, oleh [@fbbdev](https://github.com/fbbdev) dalam [#4633](https://github.com/wailsapp/wails/pull/4633)
- Tray Windows kini mematuhi `SystemTray.Show()`/`Hide()` dengan mengalihkan `NIS_HIDDEN`, sehingga aplikasi benar-benar dapat menghilang dan muncul kembali (#4653).
- Pendaftaran tray menggunakan kembali ikon yang telah di-resolve, menetapkan `NOTIFYICON_VERSION_4` satu kali, dan mengaktifkan `NIF_SHOWTIP` agar tooltip pulih setelah Explorer dimulai ulang (#4653).
- macOS: Gunakan `visibleFrame` sebagai pengganti `frame` untuk memusatkan jendela dengan mengecualikan area bilah menu dan Dock
- macOS: Gunakan `visibleFrame` sebagai pengganti `frame` untuk memusatkan jendela dengan mengecualikan area bilah menu dan Dock
- Saat menjalankan `wails3 update build-assets` dengan parameter `-config`, nilai yang ditetapkan melalui parameter `-product*` akan
- `window.NativeWindowHandle()` -> `window.NativeWindow()` oleh @leaanthony dalam [#4471](https://github.com/wailsapp/wails/pull/4471)
- Refaktor penanganan jendela internal oleh @leaanthony dalam [#4471](https://github.com/wailsapp/wails/pull/4471)
- Hapus `application.WindowIDKey` dan `application.WindowNameKey` (digantikan oleh `application.WindowKey`), oleh [@leaanthony](https://github.com/leaanthony)
- ContextMenuData kini mengembalikan string, bukan any, oleh [@leaanthony](https://github.com/leaanthony)
- Dalam binding JS/TS, field kelas bertipe array dengan panjang tetap kini diinisialisasi dengan panjang yang diharapkan dan tidak lagi dibiarkan kosong, oleh [@fbbdev](https://github.com/fbbdev) dalam [#4001](https://github.com/wailsapp/wails/pull/4001)
- ContextMenuData kini mengembalikan string, bukan any, oleh [@leaanthony](https://github.com/leaanthony)
- `application.NewService` tidak lagi menerima options sebagai parameter opsional (gunakan `application.NewServiceWithOptions` sebagai gantinya), oleh [@leaanthony](https://github.com/leaanthony) dalam [#4024](https://github.com/wailsapp/wails/pull/4024)
- Hapus dependensi `nanoid` oleh [@leaanthony](https://github.com/leaanthony)
- Perbarui contoh Window untuk gaya jendela mica/acrylic/tabbed oleh [@leaanthony](https://github.com/leaanthony)
- Dalam binding JS/TS, file model `internal.js/ts` telah dihapus; semua model kini dapat ditemukan di `models.js/ts`, oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Dalam binding JS/TS, tipe bernama tidak pernah dirender sebagai alias bagi tipe bernama lainnya; perilaku lama kini dibatasi hanya untuk alias, oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Dalam binding JS/TS pada mode kelas, field struct yang tipenya merupakan parameter tipe ditandai sebagai opsional dan tidak pernah diinisialisasi secara otomatis, oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Hapus ESLint dari templat oleh [@IanVS](https://github.com/IanVS) dalam [#4059](https://github.com/wailsapp/wails/pull/4059)
- Perbarui tahun hak cipta menjadi 2025 oleh [@IanVS](https://github.com/IanVS) dalam [#4037](https://github.com/wailsapp/wails/pull/4037)
- Tambahkan dokumentasi untuk event.Sender oleh [@IanVS](https://github.com/IanVS) dalam [#4075](https://github.com/wailsapp/wails/pull/4075)
- Dukungan Go 1.24 oleh [@leaanthony](https://github.com/leaanthony)
- Hook `ServiceStartup` kini dipanggil saat `App.Run` dipanggil, bukan dalam `application.New`, oleh [@fbbdev](https://github.com/fbbdev) dalam [#4066](https://github.com/wailsapp/wails/pull/4066)
- Error `ServiceStartup` kini dikembalikan dari `App.Run` alih-alih menghentikan proses, oleh [@fbbdev](https://github.com/fbbdev) dalam [#4066](https://github.com/wailsapp/wails/pull/4066)
- Panggilan binding dan dialog dari JS kini ditolak dengan objek error, bukan string, oleh [@fbbdev](https://github.com/fbbdev) dalam [#4066](https://github.com/wailsapp/wails/pull/4066)
- Tingkatkan penempatan menu systray pada Windows oleh [@leaanthony](https://github.com/leaanthony)
- Runtime JS telah dipindahkan ke TypeScript oleh [@fbbdev](https://github.com/fbbdev) dalam [#4100](https://github.com/wailsapp/wails/pull/4100)
- Runtime diinisialisasi segera setelah diimpor; tidak perlu menunggu jendela dimuat oleh [@fbbdev](https://github.com/fbbdev) di [#4100](https://github.com/wailsapp/wails/pull/4100)
- Runtime tidak lagi mengekspor metode init. Impor khusus efek samping dapat digunakan untuk menginisialisasinya oleh [@fbbdev](https://github.com/fbbdev) di [#4100](https://github.com/wailsapp/wails/pull/4100)
- Metode terikat kini mengembalikan `CancellablePromise` yang ditolak dengan `CancelError` jika dibatalkan. Hasil aktual panggilan tersebut dibuang oleh [@fbbdev](https://github.com/fbbdev) di [#4100](https://github.com/wailsapp/wails/pull/4100)
- Tipe layanan bawaan kini secara konsisten disebut `Service` oleh [@fbbdev](https://github.com/fbbdev) di [#4067](https://github.com/wailsapp/wails/pull/4067)
- Fungsi pembuatan layanan bawaan dengan opsi kini secara konsisten disebut `NewWithConfig` oleh [@fbbdev](https://github.com/fbbdev) di [#4067](https://github.com/wailsapp/wails/pull/4067)
- Metode `Select` pada layanan `sqlite` kini dinamai `Query` agar konsisten dengan API Go oleh [@fbbdev](https://github.com/fbbdev) di [#4067](https://github.com/wailsapp/wails/pull/4067)
- Templat: memindahkan runtime ke "dependencies" dan menata file package.json oleh [@IanVS](https://github.com/IanVS) di [#4133](https://github.com/wailsapp/wails/pull/4133)
- Membuat dan menandatangani ad hoc bundel aplikasi dalam mode pengembangan untuk mengaktifkan API macOS tertentu oleh [@popaprozac](https://github.com/popaprozac) di [#4171](https://github.com/wailsapp/wails/pull/4171)
- Memindahkan aset build ke direktori khusus platform oleh [@leaanthony](https://github.com/leaanthony)
- Memindahkan dan mengganti nama Taskfile ke direktori khusus platform oleh [@leaanthony](https://github.com/leaanthony)
- Meningkatkan pengalaman secara signifikan ketika `index.html` tidak tersedia oleh [@leaanthony](https://github.com/leaanthony)
- [Windows] Meningkatkan performa saat meminimalkan dan memulihkan jendela oleh [@leaanthony](https://github.com/leaanthony). Berdasarkan [PR](https://github.com/wailsapp/wails/pull/3955) asli oleh [562589540](https://github.com/562589540)
- Menghapus opsi `ShouldClose` (sebagai gantinya, daftarkan hook untuk events.Common.WindowClosing) oleh [@leaanthony](https://github.com/leaanthony)
- [Windows] Mengurangi kedipan saat membuka jendela oleh [@leaanthony](https://github.com/leaanthony)
- Menghapus `Window.Destroy` karena fungsi ini dimaksudkan sebagai fungsi internal oleh [@leaanthony](https://github.com/leaanthony)
- Mengganti nama event `WindowClose` menjadi `WindowClosing` oleh [@leaanthony](https://github.com/leaanthony)
- Build frontend kini menggunakan lingkungan vite "development" atau "production", bergantung pada tipe build, oleh [@leaanthony](https://github.com/leaanthony)
- Memperbarui ke go-webview2 v1.19 oleh [@leaanthony](https://github.com/leaanthony)
- Memastikan fork taskfile digunakan oleh @leaanthony
- Memperbarui fork Taskfile untuk memperbaiki masalah versi saat menginstal menggunakan
- Menggunakan fork Taskfile untuk memperbaiki masalah versi saat menginstal menggunakan
- `service.OnStartup` kini menghentikan aplikasi saat terjadi kesalahan dan menjalankan
- Memfaktorkan ulang pengiriman pesan klik systray agar lebih selaras dengan interaksi pengguna oleh
- Penyematan aset kini menyertakan `all:frontend/dist` untuk mendukung framework yang menghasilkan
- Pemfaktoran ulang Taskfile oleh [leaanthony](https://github.com/leaanthony) di
- Memutakhirkan ke `go-webview2` v1.0.16 oleh
- Memperbaiki tipe `Screen` agar menyertakan `ID`, bukan `Id`, oleh
- Memperbarui versi wails `go.mod.tmpl` untuk mendukung `application.ServiceOptions` oleh
- Memperbaiki penentuan nama layanan oleh [windom](https://github.com/windom/) di
- mkdocs serve kini menggunakan docker oleh [leaanthony](https://github.com/leaanthony)
- Menggabungkan konfigurasi pengembangan ke dalam `config.yml` oleh
- Dialog systray kini secara default menggunakan ikon aplikasi jika tersedia (Windows) oleh
- Meningkatkan pelaporan GPU + Memori untuk macOS oleh
- Menghapus `WebviewGpuIsDisabled` dan `EnableFraudulentWebsiteWarnings`
- Perubahan API Events: `On`/`Emit` -> event pengguna, `OnApplicationEvent` ->
- Perbaikan API Events di Linux oleh [TheGB0077](https://github.com/TheGB0077) di
- [CI] penyempurnaan action dan pengaktifan action agar juga dapat dijalankan di fork dan
- Mengganti nama `AbsolutePosition()` menjadi `Position()` oleh
- Memperbarui dependensi webkit Linux ke webkit2gtk-4.1 dari webkitgtk2-4.0 untuk
- Skrip runtime JS yang dibundel kini merupakan modul ESM: tag script yang mengimpornya
- Paket `@wailsio/runtime` tidak memublikasikan API-nya pada `window.wails`
- Modul Window API `@wailsio/runtime/src/window` kini mengekspos objek yang memuatnya
- API jendela JS telah diperbarui agar sesuai dengan `WebviewWindow` Go saat ini
- Generator binding kini menggunakan panggilan berdasarkan ID secara default. Opsi CLI `-id`
- Tata letak baru kode binding: sebelumnya file keluaran ditata dalam folder
- Nama field struct `application.Options.Bind` telah diubah menjadi
- Sintaks baru untuk mengikat layanan: instance layanan kini harus dibungkus dalam
- Menonaktifkan spinner di lingkungan nonterminal atau CI oleh

### Dihapus

- **BREAKING**: Menghapus `EnabledFeatures`, `DisabledFeatures`, dan `AdditionalLaunchArgs` dari opsi `WindowsWindow` per jendela. Sebagai gantinya, gunakan `Options.Windows.EnabledFeatures`, `Options.Windows.DisabledFeatures`, dan `Options.Windows.AdditionalBrowserArgs` tingkat aplikasi. Flag ini berlaku secara global pada lingkungan WebView2 bersama (#4559) oleh @leaanthony
- Menghapus implementasi native `IDropTarget` di Windows dan menggantinya dengan pendekatan berbasis JavaScript (sesuai dengan perilaku v2)
- Menghapus dependensi github.com/wailsapp/mimetype dan menggantinya dengan peta ekstensi yang diperluas + http.DetectContentType dari pustaka standar, sehingga ukuran biner berkurang sekitar 1.2MB
- Menghapus dependensi gopkg.in/ini.v1 dengan mengimplementasikan parser minimal untuk file .desktop bagi penjelajah file Linux, sehingga menghemat sekitar 45KB
- Menghapus samber/lo dari kode runtime dengan menggunakan paket slices pustaka standar Go 1.21+ dan helper internal minimal, sehingga menghemat ~310 KB
- Menghapus pernyataan printf debug dari handler skema URL Darwin (#4834)
- **BREAKING**: Menghapus event `linux:WindowLoadChanged`—gunakan `linux:WindowLoadFinished` sebagai gantinya untuk mendeteksi saat WebView selesai dimuat (#3896), oleh @leaanthony

### Perubahan yang Merusak Kompatibilitas

- **Pemfaktoran Ulang API Manager**: Menata ulang API aplikasi dari struktur datar menjadi sejumlah manager yang terorganisasi agar pengorganisasian kode dan kemudahan penemuannya lebih baik, oleh [@leaanthony](https://github.com/leaanthony) dalam [#4359](https://github.com/wailsapp/wails/pull/4359)
- `app.NewWebviewWindow()` → `app.Window.New()`
- `app.CurrentWindow()` → `app.Window.Current()`
- `app.GetAllWindows()` → `app.Window.GetAll()`
- `app.WindowByName()` → `app.Window.GetByName()`
- `app.EmitEvent()` → `app.Event.Emit()`
- `app.OnApplicationEvent()` → `app.Event.OnApplicationEvent()`
- `app.OnWindowEvent()` → `app.Event.OnWindowEvent()`
- `app.SetApplicationMenu()` → `app.Menu.SetApplicationMenu()`
- `app.OpenFileDialog()` → `app.Dialog.OpenFile()`
- `app.SaveFileDialog()` → `app.Dialog.SaveFile()`
- `app.MessageDialog()` → `app.Dialog.Message()`
- `app.InfoDialog()` → `app.Dialog.Info()`
- `app.WarningDialog()` → `app.Dialog.Warning()`
- `app.ErrorDialog()` → `app.Dialog.Error()`
- `app.QuestionDialog()` → `app.Dialog.Question()`
- `app.NewSystemTray()` → `app.SystemTray.New()`
- `app.GetSystemTray()` → `app.SystemTray.Get()`
- `app.ShowContextMenu()` → `app.ContextMenu.Show()`
- `app.RegisterKeybinding()` → `app.KeyBinding.Register()`
- `app.UnregisterKeybinding()` → `app.KeyBinding.Unregister()`
- `app.GetPrimaryScreen()` → `app.Screen.GetPrimary()`
- `app.GetAllScreens()` → `app.Screen.GetAll()`
- `app.BrowserOpenURL()` → `app.Browser.OpenURL()`
- `app.Environment()` → `app.Env.GetAll()`
- `app.ClipboardGetText()` → `app.Clipboard.Text()`
- `app.ClipboardSetText()` → `app.Clipboard.SetText()`
- Mengganti nama metode Service: `Name` -> `ServiceName`, `OnStartup` -> `ServiceStartup`, `OnShutdown` -> `ServiceShutdown`, oleh [@leaanthony](https://github.com/leaanthony)
- Memindahkan metode `Path` dan `Paths` ke paket `application`, oleh [@leaanthony](https://github.com/leaanthony)
- Menu aplikasi kini hanya tersedia di macOS, oleh [@leaanthony](https://github.com/leaanthony)

## v3.0.0-alpha.78 - 2026-04-21

## Ditambahkan

## Diperbaiki

## v3.0.0-alpha.77 - 2026-04-18

## Diperbaiki

## v3.0.0-alpha.76 - 2026-04-17

## Diperbaiki

## v3.0.0-alpha.75 - 2026-04-16

## Diperbaiki

## v3.0.0-alpha.74 - 2026-03-01

## Ditambahkan

## Diperbaiki

## v3.0.0-alpha.73 - 2026-02-27

## Diperbaiki

## v3.0.0-alpha.72 - 2026-02-16

## Diperbaiki

## v3.0.0-alpha.71 - 2026-02-10

## Ditambahkan

## Diperbaiki

## v3.0.0-alpha.70 - 2026-02-09

## Ditambahkan

## Diperbaiki

## v3.0.0-alpha.69 - 2026-02-08

## Ditambahkan

## Diperbaiki

## v3.0.0-alpha.68 - 2026-02-07

## Ditambahkan

## Diubah

## Diperbaiki

## v3.0.0-alpha.67 - 2026-02-04

## Ditambahkan

## Diubah

## Diperbaiki

## v3.0.0-alpha.66 - 2026-02-03

## Ditambahkan

## Diubah

## Diperbaiki

## Dihapus

## v3.0.0-alpha.65 - 2026-02-01

## Ditambahkan

## v3.0.0-alpha.64 - 2026-01-26

## Ditambahkan

## v3.0.0-alpha.63 - 2026-01-25

## Diperbaiki

## v3.0.0-alpha.62 - 2026-01-22

## Diperbaiki

## v3.0.0-alpha.61 - 2026-01-20

## Diperbaiki

## v3.0.0-alpha.60 - 2026-01-14

## Diperbaiki

## v3.0.0-alpha.59 - 2026-01-11

## Diubah

## v3.0.0-alpha.58 - 2026-01-09

## Diperbaiki

## v3.0.0-alpha.57 - 2026-01-05

## Diubah

## Diperbaiki

## v3.0.0-alpha.56 - 2026-01-04

## Ditambahkan

## Diubah

## Diperbaiki

## Dihapus

## v3.0.0-alpha.55 - 2026-01-02

## Diubah

## Diperbaiki

## Dihapus

## v3.0.0-alpha.54 - 2025-12-29

## Ditambahkan

## Diperbaiki

## Dihapus

## v3.0.0-alpha.53 - 2025-12-27

## Ditambahkan

## Diperbaiki

## v3.0.0-alpha.52 - 2025-12-26

## Diperbaiki

## v3.0.0-alpha.51 - 2025-12-23

## Diperbaiki

## v3.0.0-alpha.50 - 2025-12-21

## Diubah

## v3.0.0-alpha.49 - 2025-12-18

## Diubah

## v3.0.0-alpha.48 - 2025-12-16

## Ditambahkan

## Diubah

## Diperbaiki

## v3.0.0-alpha.47 - 2025-12-15

## Ditambahkan

## Diperbaiki

## v3.0.0-alpha.46 - 2025-12-14

## Ditambahkan

## Dihapus

## v3.0.0-alpha.45 - 2025-12-13

## Ditambahkan

## Diperbaiki

## v3.0.0-alpha.44 - 2025-12-12

## Ditambahkan

## Diubah

## Diperbaiki

## v3.0.0-alpha.43 - 2025-12-11

## Ditambahkan

## v3.0.0-alpha.42 - 2025-12-10

## Ditambahkan

## v3.0.0-alpha.41 - 2025-11-23

## Diperbaiki

## v3.0.0-alpha.40 - 2025-11-13

## Diperbaiki

## v3.0.0-alpha.39 - 2025-11-12

## Ditambahkan

## Diubah

## v3.0.0-alpha.38 - 2025-11-04

## Ditambahkan

## Diubah

## Diperbaiki

## v3.0.0-alpha.37 - 2025-11-02

## Diperbaiki

## v3.0.0-alpha.36 - 2025-10-15

## Diperbaiki

## v3.0.0-alpha.35 - 2025-10-14

## Diperbaiki

## v3.0.0-alpha.34 - 2025-10-06

## Ditambahkan

## Diperbaiki

## v3.0.0-alpha.33 - 2025-10-04

## Diperbaiki

## v3.0.0-alpha.32 - 2025-10-02

## Diperbaiki

## v3.0.0-alpha.31 - 2025-09-27

## Diperbaiki

## v3.0.0-alpha.30 - 2025-09-26

## Diperbaiki

## v3.0.0-alpha.29 - 2025-09-25

## Ditambahkan

## Diubah

## Diperbaiki

## v3.0.0-alpha.29 - 2025-09-25

## Ditambahkan

## Diubah

## Diperbaiki

## v3.0.0-alpha.27 - 2025-09-07

## Diperbaiki

## v3.0.0-alpha.26 - 2025-08-24

## Ditambahkan

## v3.0.0-alpha.25 - 2025-08-16

## Diubah

tidak lagi diabaikan dan menimpa nilai konfigurasi.

## v3.0.0-alpha.24 - 2025-08-13

## Ditambahkan

## v3.0.0-alpha.23 - 2025-08-11

## Diperbaiki

## v3.0.0-alpha.22 - 2025-08-10

## Ditambahkan

## Diubah

+ Perbaiki dependensi paket Linux yang terlalu luas dan dependensi RPM yang sudah usang.

## v3.0.0-alpha.21 - 2025-08-07

## Diperbaiki

## v3.0.0-alpha.20 - 2025-08-06

## Diperbaiki

## v3.0.0-alpha.19 - 2025-08-05

## Ditambahkan

## Diperbaiki

## v3.0.0-alpha.18 - 2025-08-03

## Ditambahkan

## Diperbaiki

## v3.0.0-alpha.17 - 2025-07-31

## Diperbaiki

## v3.0.0-alpha.16 - 2025-07-25

## Ditambahkan

## v3.0.0-alpha.15 - 2025-07-25

## Ditambahkan

## v3.0.0-alpha.14 - 2025-07-25

## Ditambahkan

## v3.0.0-alpha.12 - 2025-07-15

### Ditambahkan

### Diperbaiki

## v3.0.0-alpha.11 - 2025-07-12

## Ditambahkan

## v3.0.0-alpha.10 - 2025-07-06

### Perubahan yang Tidak Kompatibel

### Ditambahkan

### Diperbaiki

### Diubah

## v3.0.0-alpha.9 - 2025-01-13

### Ditambahkan

### Diperbaiki

### Diubah

## v3.0.0-alpha.8.3 - 2024-12-07

### Diubah

## v3.0.0-alpha.8.2 - 2024-12-07

### Diubah

`go install` oleh @leaanthony

## v3.0.0-alpha.8.1 - 2024-12-07

### Diubah

`go install` oleh @leaanthony

## v3.0.0-alpha.8 - 2024-12-06

### Ditambahkan

@atterpac dalam [#3909](https://github.com/wailsapp/wails/3909)   [ansxuman](https://github.com/ansxuman) dalam   [#3902](https://github.com/wailsapp/wails/pull/3902)   [atterpac](https://github.com/atterpac) dalam   [#3867](https://github.com/wailsapp/wails/pull/3867)   oleh [atterpac](https://github.com/atterpac) dalam   [#3829](https://github.com/wailsapp/wails/pull/3829)   [leaanthony](https://github.com/leaanthony)   [FerroO2000](https://github.com/FerroO2000) dalam   [#3856](https://github.com/wailsapp/wails/pull/3856)   [#3873](https://github.com/wailsapp/wails/pull/3873)   [leaanthony](https://github.com/leaanthony)   ditempatkan pada lokasi X/Y yang ditentukan — diimplementasikan oleh   [leaanthony](https://github.com/leaanthony) dalam   [#3885](https://github.com/wailsapp/wails/pull/3885)   [ansxuman](https://github.com/ansxuman) dan   [leaanthony](https://github.com/leaanthony) dalam   [#3823](https://github.com/wailsapp/wails/pull/3823)   oleh [leaanthony](https://github.com/leaanthony) dalam   [#3766](https://github.com/wailsapp/wails/pull/3766)   [leaanthony](https://github.com/leaanthony) dalam   [#3888](https://github.com/wailsapp/wails/pull/3888)   [leaanthony](https://github.com/leaanthony). -

### Diubah

`service.OnShutdown`untuk setiap layanan yang telah dimulai sebelumnya — perubahan oleh @atterpac dalam   [#3920](https://github.com/wailsapp/wails/pull/3920)   @atterpac dalam [#3907](https://github.com/wailsapp/wails/pull/3907)   subfolder oleh @atterpac dalam   [#3887](https://github.com/wailsapp/wails/pull/3887)   [#3748](https://github.com/wailsapp/wails/pull/3748)   [leaanthony](https://github.com/leaanthony)   [etesam913](https://github.com/etesam913) dalam   [#3778](https://github.com/wailsapp/wails/pull/3778)   [northes](https://github.com/northes) dalam   [#3836](https://github.com/wailsapp/wails/pull/3836)   [#3827](https://github.com/wailsapp/wails/pull/3827)   [leaanthony](https://github.com/leaanthony)   [@leaanthony](https://github.com/leaanthony)   [@leaanthony](https://github.com/leaanthony)   (digantikan oleh opsi `EnabledFeatures` dan `DisabledFeatures`) oleh   [leaanthony](https://github.com/leaanthony)

### Diperbaiki

variabel channel oleh @michael-freling dalam   [#3925](https://github.com/wailsapp/wails/pull/3925)   [ansxuman](https://github.com/ansxuman) dalam   [#3924](https://github.com/wailsapp/wails/pull/3924)   [#3898](https://github.com/wailsapp/wails/pull/3898)   [#3901](https://github.com/wailsapp/wails/pull/3901)   dalam [#3886](https://github.com/wailsapp/wails/pull/3886)   [leaanthony](https://github.com/leaanthony) dalam   [#3841](https://github.com/wailsapp/wails/pull/3841)   peran `PasteAndMatchStyle` dalam menu edit di Darwin oleh   [johnmccabe](https://github.com/johnmccabe) dalam   [#3839](https://github.com/wailsapp/wails/pull/3839)   [#3840](https://github.com/wailsapp/wails/issues/3840) dalam   [#3854](https://github.com/wailsapp/wails/pull/3854) oleh   [kodflow](https://github.com/kodflow)   [@leaanthony](https://github.com/leaanthony)   berbeda. oleh @nickisworking dalam   [#3789](https://github.com/wailsapp/wails/pull/3789)

## v3.0.0-alpha.7 - 2024-09-18

### Ditambahkan

[mmghv](https://github.com/mmghv) dalam   [#3665](https://github.com/wailsapp/wails/pull/3665)   [#3682](https://github.com/wailsapp/wails/pull/3682)   [atterpac](https://github.com/atterpac) dan   [leaanthony](https://github.com/leaanthony) dalam   [#3570](https://github.com/wailsapp/wails/pull/3570)

### Diubah

Event Aplikasi `OnWindowEvent` -> Event Jendela, oleh   [leaanthony](https://github.com/leaanthony)   [#3734](https://github.com/wailsapp/wails/pull/3734)   branch dengan prefiks `v3/` atau `v3-` oleh   [stendler](https://github.com/stendler) dalam   [#3747](https://github.com/wailsapp/wails/pull/3747)

### Diperbaiki

[etesam913](https://github.com/etesam913) dalam   [#3742](https://github.com/wailsapp/wails/pull/3742)   [atterpac](https://github.com/atterpac) dalam   [#3721](https://github.com/wailsapp/wails/pull/3721)   [atterpac](https://github.com/atterpac) dalam   [#3675](https://github.com/wailsapp/wails/pull/3675)   [#1811](https://github.com/wailsapp/wails/pull/1811) dalam   [#3614](https://github.com/wailsapp/wails/pull/3614) oleh   [@stendler](https://github.com/stendler)   [#3720](https://github.com/wailsapp/wails/pull/3720) oleh   [leaanthony](https://github.com/leaanthony)   [#3693](https://github.com/wailsapp/wails/issues/3693) oleh   [@DeltaLaboratory](https://github.com/DeltaLaboratory)   [#3720](https://github.com/wailsapp/wails/pull/3720) oleh   [leaanthony](https://github.com/leaanthony)   [#3693](https://github.com/wailsapp/wails/issues/3693) oleh   [@DeltaLaboratory](https://github.com/DeltaLaboratory)   [leaanthony](https://github.com/leaanthony)   [#3746](https://github.com/wailsapp/wails/pull/3746) oleh   [@stendler](https://github.com/stendler)   [leaanthony](https://github.com/leaanthony)

## v3.0.0-alpha.6 - 2024-07-30

### Diperbaiki

## v3.0.0-alpha.5 - 2024-07-30

### Ditambahkan

[#3580](https://github.com/wailsapp/wails/pull/3580)   [#3580](https://github.com/wailsapp/wails/pull/3580)   klik ikon oleh @5aaee9 dalam [#2991](https://github.com/wailsapp/wails/pull/2991)   [#2618](https://github.com/wailsapp/wails/pull/2618)   [@fbbdev](https://github.com/fbbdev) dalam   [#3282](https://github.com/wailsapp/wails/pull/3282)   @[Atterpac](https://github.com/Atterpac)   dalam [#3022](https://github.com/wailsapp/wails/pull/3022])   [@marcus-crane](https://github.com/marcus-crane) dalam   [#3146](https://github.com/wailsapp/wails/pull/3146)   [PR](https://github.com/wailsapp/wails/pull/3147)   [PR](https://github.com/wailsapp/wails/pull/3189)   [@fbbdev](https://github.com/fbbdev) dalam   [#3281](https://github.com/wailsapp/wails/pull/3281)   [aba82cc](https://github.com/wailsapp/wails/commit/aba82cc52787c97fb99afa58b8b63a0004b7ff6c)   berdasarkan [PR](https://github.com/wailsapp/wails/pull/2044) oleh @Mai-Lapyst   [@fbbdev](https://github.com/fbbdev) dalam   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@fbbdev](https://github.com/fbbdev) dalam   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@fbbdev](https://github.com/fbbdev) dalam   [#3295](https://github.com/wailsapp/wails/pull/3295)   paket npm oleh [@fbbdev](https://github.com/fbbdev) dalam   [#3334](https://github.com/wailsapp/wails/pull/3334)   dalam [#3354](https://github.com/wailsapp/wails/pull/3354)   `WAILS_VITE_PORT` oleh [@abichinger](https://github.com/abichinger) dalam   [#3429](https://github.com/wailsapp/wails/pull/3429)   [@abichinger](https://github.com/abichinger) dalam   [#3431](https://github.com/wailsapp/wails/pull/3431)   [@bruxaodev](https://github.com/bruxaodev) dalam   [#3667](https://github.com/wailsapp/wails/pull/3667)   [@OlegGulevskyy](https://github.com/OlegGulevskyy) dalam   [#3674](https://github.com/wailsapp/wails/pull/3674)

### Diperbaiki

[#3606](https://github.com/wailsapp/wails/pull/3606)   [tmclane](https://github.com/tmclane) dalam   [#3515](https://github.com/wailsapp/wails/pull/3515)   [atterpac](https://github.com/atterac) dalam   [#3512](https://github.com/wailsapp/wails/pull/3512)   [atterpac](https://github.com/atterpac) dalam   [#3477](https://github.com/wailsapp/wails/pull/3477)   oleh [Atterpac](https://github.com/atterpac) dalam   [#3320](https://github.com/wailsapp/wails/pull/3320).   dalam [#3306](https://github.com/wailsapp/wails/pull/3306).   [#2972](https://github.com/wailsapp/wails/pull/2972).   [#2982](https://github.com/wailsapp/wails/pull/2982)   [mmghv](https://github.com/mmghv) dalam   [#2750](https://github.com/wailsapp/wails/pull/2750).   [#2753](https://github.com/wailsapp/wails/pull/2753).   [jaybeecave](https://github.com/jaybeecave) dalam   [#3052](https://github.com/wailsapp/wails/pull/3052).   [@pylotlight](https://github.com/pylotlight) dalam   [PR](https://github.com/wailsapp/wails/pull/3039)   terinstal. Ditambahkan oleh [@pylotlight](https://github.com/pylotlight) dalam   [PR](https://github.com/wailsapp/wails/pull/3032)   [PR](https://github.com/wailsapp/wails/pull/3145)   spasi - @leaanthony.   [thomas-senechal](https://github.com/thomas-senechal) dalam PR   [#3207](https://github.com/wailsapp/wails/pull/3207)   [thomas-senechal](https://github.com/thomas-senechal) dalam PR   [#3208](https://github.com/wailsapp/wails/pull/3208)   jendela yang terhubung [tw1nk](https://github.com/tw1nk) dalam PR   [#3271](https://github.com/wailsapp/wails/pull/3271)   [#3273](https://github.com/wailsapp/wails/pull/3273)   [@fbbdev](https://github.com/fbbdev) dalam   [#3279](https://github.com/wailsapp/wails/pull/3279)   `FRONTEND_DEVSERVER_URL` tersedia.   [#3299](https://github.com/wailsapp/wails/pull/3299)   [@fbbdev](https://github.com/fbbdev) dalam   [#3295](https://github.com/wailsapp/wails/pull/3295)   listener oleh [@fbbdev](https://github.com/fbbdev) dalam   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@abichinger](https://github.com/abichinger) dalam   [#3330](https://github.com/wailsapp/wails/pull/3330)   generator oleh [@fbbdev](https://github.com/fbbdev) dalam   [#3334](https://github.com/wailsapp/wails/pull/3334)   generator oleh [@fbbdev](https://github.com/fbbdev) dalam   [#3334](https://github.com/wailsapp/wails/pull/3334)   [@abichinger](https://github.com/abichinger) dalam   [#3346](https://github.com/wailsapp/wails/pull/3346)   [@hfoxy](https://github.com/hfoxy) dalam   [#3417](https://github.com/wailsapp/wails/pull/3417)   [@hfoxy](https://github.com/hfoxy) dalam   [#3426](https://github.com/wailsapp/wails/pull/3426)   [@fbbdev](https://github.com/fbbdev) dalam   [#3431](https://github.com/wailsapp/wails/pull/3431)   [@fbbdev](https://github.com/fbbdev) dalam   [#3431](https://github.com/wailsapp/wails/pull/3431)   oleh [@pekim](https://github.com/pekim) dalam   [#3458](https://github.com/wailsapp/wails/pull/3458)   [@robin-samuel](https://github.com/robin-samuel).   PR [#3466](https://github.com/wailsapp/wails/pull/3622) oleh   [@5aaee9](https://github.com/5aaee9).   [@windom](https://github.com/windom/) dalam   [#3636](https://github.com/wailsapp/wails/pull/3636).   Windows oleh [@bruxaodev](https://github.com/bruxaodev/) dalam   [#3691](https://github.com/wailsapp/wails/pull/3691).

### Diubah

[mmghv](https://github.com/mmghv) dalam   [#3611](https://github.com/wailsapp/wails/pull/3611)   mendukung Ubuntu 24.04 LTS oleh [atterpac](https://github.com/atterpac) dalam   [#3461](https://github.com/wailsapp/wails/pull/3461)   harus memiliki atribut `type="module"`. Oleh   [@fbbdev](https://github.com/fbbdev) dalam   [#3295](https://github.com/wailsapp/wails/pull/3295)   objek, dan tidak memulai sistem WML. Perubahan ini dilakukan untuk meningkatkan   enkapsulasi. Jika diinginkan, sistem WML dapat dimulai secara manual dengan memanggil   metode baru `WML.Enable`. Skrip runtime JS yang dibundel tetap menjalankan kedua   operasi tersebut secara otomatis. Oleh [@fbbdev](https://github.com/fbbdev) dalam   [#3295](https://github.com/wailsapp/wails/pull/3295)   objek window sebagai ekspor default. Metode individual tidak dapat lagi diimpor   melalui sintaks impor bernama atau namespace ESM.   API. Beberapa metode telah berganti nama atau prototipe, khususnya: `Screen`   menjadi `GetScreen`; `GetZoomLevel`/`SetZoomLevel` menjadi `GetZoom`/`SetZoom`;   `GetZoom`, `Width`, dan `Height` kini mengembalikan nilai secara langsung, bukan membungkusnya   di dalam objek. Oleh [@fbbdev](https://github.com/fbbdev) dalam   [#3295](https://github.com/wailsapp/wails/pull/3295)   telah dihapus. Gunakan opsi CLI `-names` untuk kembali menggunakan pemanggilan berdasarkan nama.   Oleh [@fbbdev](https://github.com/fbbdev) dalam   [#3468](https://github.com/wailsapp/wails/pull/3468)   diberi nama berdasarkan paket yang memuatnya; kini path impor Go lengkap digunakan,   termasuk path modul. Oleh [@fbbdev](https://github.com/fbbdev) dalam   [#3468](https://github.com/wailsapp/wails/pull/3468)   `application.Options.Services`. Oleh [@fbbdev](https://github.com/fbbdev) dalam   [#3468](https://github.com/wailsapp/wails/pull/3468)   pemanggilan ke `application.NewService`. Oleh [@fbbdev](https://github.com/fbbdev) dalam   [#3468](https://github.com/wailsapp/wails/pull/3468)   [@DeltaLaboratory](https://github.com/DeltaLaboratory) dalam   [#3574](https://github.com/wailsapp/wails/pull/3574)
