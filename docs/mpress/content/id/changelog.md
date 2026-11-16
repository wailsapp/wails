---
title: "Log perubahan"
description: "Riwayat versi dan catatan rilis untuk Wails v3"
slug: "changelog"
sourcePath: "changelog.md"
---

Legenda:

-  - macOS
- ⊞ - Windows
- 🐧 - Linux

/_-- Semua perubahan penting pada proyek ini akan didokumentasikan dalam berkas ini.

Format ini didasarkan pada [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), dan proyek ini mengikuti [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

- `Added` untuk fitur baru.
- `Changed` untuk perubahan pada fungsionalitas yang sudah ada.
- `Deprecated` untuk fitur yang akan segera dihapus.
- `Removed` untuk fitur yang kini telah dihapus.
- `Fixed` untuk semua perbaikan bug.
- `Security` jika terdapat kerentanan.

_/

/_   * DIPERBARUI SECARA OTOMATIS *   Catatan rilis dibuat dari PR yang telah digabungkan. Jangan edit berkas ini dalam PR kontribusi biasa atau membuat ulang `v3/UNRELEASED_CHANGELOG.md` yang telah dihapus. Lihat `AGENTS.md` untuk pemeliharaan log perubahan. _/

## [Belum dirilis]

## v3.2.0-nightly.20261103 - 2026-11-03

## Diubah
- [v3] perbarui menu baki sistem saat Menu.Update dalam [PR](https://github.com/wailsapp/wails/pull/6314) oleh @johnmaguire
- [v3/Windows] Implementasikan SetIcon untuk jendela dalam [PR](https://github.com/wailsapp/wails/pull/6502) oleh @hypherionmc

## v3.2.0-nightly.20261102 - 2026-11-02

## Diubah
- Gunakan DisableOverscroll sebagai pengganti DisableBounce dalam [PR](https://github.com/wailsapp/wails/pull/6077) oleh @lethusang

## v3.2.0-nightly.20261101 - 2026-11-01

## Ditambahkan
- tambahkan backend macOS dan Linux eksperimental tanpa CGO dalam [PR](https://github.com/wailsapp/wails/pull/6521) oleh @taliesin-ai
- build YAML native, pemaketan portabel, dan penandatanganan dalam [PR](https://github.com/wailsapp/wails/pull/6531) oleh @taliesin-ai

## Diperbaiki
- Perbaiki berbagai masalah build nightly v3 dan rapikan dependensi dalam [PR](https://github.com/wailsapp/wails/pull/6519) oleh @taliesin-ai
- perbaiki validasi rilis dan klasifikasi fixture v3 dalam [PR](https://github.com/wailsapp/wails/pull/6533) oleh @taliesin-ai
- atasi kegagalan validasi rilis lintas platform dalam [PR](https://github.com/wailsapp/wails/pull/6538) oleh @taliesin-ai

## v3.1.0 - 2026-11-01

## Ditambahkan
- tambahkan backend macOS dan Linux eksperimental tanpa CGO dalam [PR](https://github.com/wailsapp/wails/pull/6521) oleh @taliesin-ai
- build YAML native, pemaketan portabel, dan penandatanganan dalam [PR](https://github.com/wailsapp/wails/pull/6531) oleh @taliesin-ai

## Diubah
- hindari menyalin metadata desktop yang dibuang dalam [PR](https://github.com/wailsapp/wails/pull/6391) oleh @taliesin-ai
- format panjang konten tanpa fmt dalam [PR](https://github.com/wailsapp/wails/pull/6393) oleh @taliesin-ai
- pilih root indeks tanpa mengumpulkan kandidat dalam [PR](https://github.com/wailsapp/wails/pull/6404) oleh @taliesin-ai
- gunakan kembali buffer penyalinan berkas yang tidak dapat di-seek dalam [PR](https://github.com/wailsapp/wails/pull/6408) oleh @taliesin-ai
- hindari mengalokasikan reader untuk mendeteksi prefiks dalam [PR](https://github.com/wailsapp/wails/pull/6410) oleh @taliesin-ai
- lewati pembacaan deteksi MIME untuk ekstensi yang telah dikenal dalam [PR](https://github.com/wailsapp/wails/pull/6414) oleh @taliesin-ai
- cegah pertumbuhan antrean dan osilasi penyusutan dalam [PR](https://github.com/wailsapp/wails/pull/6415) oleh @taliesin-ai
- hindari mengalokasikan konteks permintaan yang dibuang dalam [PR](https://github.com/wailsapp/wails/pull/6419) oleh @taliesin-ai
- lewati penyiapan log debug yang dinonaktifkan dalam [PR](https://github.com/wailsapp/wails/pull/6420) oleh @taliesin-ai
- alokasikan buffer deteksi MIME secara malas dalam [PR](https://github.com/wailsapp/wails/pull/6428) oleh @taliesin-ai
- kurangi alokasi pemformatan hyperlink dalam [PR](https://github.com/wailsapp/wails/pull/6443) oleh @taliesin-ai
- hapus alokasi sementara untuk pemformatan butir dalam [PR](https://github.com/wailsapp/wails/pull/6449) oleh @taliesin-ai
- kurangi konversi string UUID menjadi satu alokasi dalam [PR](https://github.com/wailsapp/wails/pull/6453) oleh @taliesin-ai
- kurangi alokasi pemformatan handle pernyataan dalam [PR](https://github.com/wailsapp/wails/pull/6454) oleh @taliesin-ai
- hilangkan alokasi digest UUID berbasis nama sementara dalam [PR](https://github.com/wailsapp/wails/pull/6455) oleh @taliesin-ai
- hindari menyalin baris keluaran perintah yang belum selesai dalam [PR](https://github.com/wailsapp/wails/pull/6460) oleh @taliesin-ai
- alirkan penyalinan executable selama penyiapan paket dalam [PR](https://github.com/wailsapp/wails/pull/6472) oleh @taliesin-ai
- hindari memindai ulang baris keluaran yang belum selesai dalam [PR](https://github.com/wailsapp/wails/pull/6473) oleh @taliesin-ai
- gunakan kembali pola perutean perintah yang telah dikompilasi dalam [PR](https://github.com/wailsapp/wails/pull/6475) oleh @taliesin-ai
- hindari pencarian metadata selama penemuan runtime dalam [PR](https://github.com/wailsapp/wails/pull/6482) oleh @taliesin-ai
- kurangi alokasi halaman fallback yang dilokalkan dalam [PR](https://github.com/wailsapp/wails/pull/6485) oleh @taliesin-ai
- hindari mengompilasi ulang pola untuk setiap pemroses dalam [PR](https://github.com/wailsapp/wails/pull/6490) oleh @taliesin-ai

## Diperbaiki
- implementasikan SetTooltip untuk baki sistem dalam [PR](https://github.com/wailsapp/wails/pull/6017) oleh @fan711
- baca mode gelap dari portal pengaturan dalam [PR](https://github.com/wailsapp/wails/pull/6072) oleh @TechHutTV
- izinkan Window.SetZoom turun di bawah 1.0 dalam [PR](https://github.com/wailsapp/wails/pull/6166) oleh @APshenkin
- kirim ke setiap listener Go sesuai urutan emit dalam [PR](https://github.com/wailsapp/wails/pull/6275) oleh @DevLumuz
- baca hook peristiwa di bawah hookLock dalam Emit dalam [PR](https://github.com/wailsapp/wails/pull/6276) oleh @DevLumuz
- selalu publikasikan properti Menu baki agar baki yang dibuat setelah Run tidak dibuang dalam [PR](https://github.com/wailsapp/wails/pull/6284) oleh @AlbinoGeek
- tahan pekerjaan thread utama yang dikirim sebelum Run hingga Run dimulai dalam [PR](https://github.com/wailsapp/wails/pull/6287) oleh @AlbinoGeek
- gunakan URL aset yang benar dalam penyedia pembaruan GitHub dalam [PR](https://github.com/wailsapp/wails/pull/6304) oleh @gavidroselj
- pertahankan isi permintaan saat WebView2 tidak menyertakan Content-Length dalam [PR](https://github.com/wailsapp/wails/pull/6326) oleh @strraaa
- laporkan layar macOS dan GTK dalam ruang DIP (batas DPI campuran) dalam [PR](https://github.com/wailsapp/wails/pull/6346) oleh @julianstorer
- jangan laporkan notifikasi yang ditutup sebagai klik dalam [PR](https://github.com/wailsapp/wails/pull/6367) oleh @bercedev
- hentikan proses pengembangan saat terminal ditutup dalam [PR](https://github.com/wailsapp/wails/pull/6371) oleh @xmasdev
- petakan proses log perubahan ke commit squash-merge dalam pemeriksaan bukti pada [PR](https://github.com/wailsapp/wails/pull/6373) oleh @taliesin-ai
- tahan debounce pemulaian ulang bilah tugas selama satu detik dalam [PR](https://github.com/wailsapp/wails/pull/6386) oleh @taliesin-ai
- pertahankan batasan ukuran saat beralih ke layar penuh dari kondisi tersebut dalam [PR](https://github.com/wailsapp/wails/pull/6387) oleh @taliesin-ai
- jangan ubah ukuran jendela layar penuh atau yang dimaksimalkan agar sesuai dengan batasan baru dalam [PR](https://github.com/wailsapp/wails/pull/6388) oleh @taliesin-ai
- jangan tampilkan dialog sebelum Run dalam contoh liquid-glass dalam [PR](https://github.com/wailsapp/wails/pull/6389) oleh @taliesin-ai
- dialog yang diminta sebelum Run tidak melakukan apa pun alih-alih menunggu selamanya dalam [PR](https://github.com/wailsapp/wails/pull/6390) oleh @taliesin-ai
- gunakan fallback saat entri desktop tidak memiliki perintah dalam [PR](https://github.com/wailsapp/wails/pull/6392) oleh @taliesin-ai
- deteksi konten saat ini tanpa menyimpan nama file dalam [PR](https://github.com/wailsapp/wails/pull/6394) oleh @taliesin-ai
- kenali ekstensi aset dengan campuran huruf besar dan kecil dalam [PR](https://github.com/wailsapp/wails/pull/6395) oleh @taliesin-ai
- lepaskan buffer pendeteksi konten yang telah digunakan dalam [PR](https://github.com/wailsapp/wails/pull/6396) oleh @taliesin-ai
- lepaskan referensi ringqueue yang telah dihapus dalam [PR](https://github.com/wailsapp/wails/pull/6397) oleh @taliesin-ai
- patuhi penghapusan header sebelum memilih fallback dalam [PR](https://github.com/wailsapp/wails/pull/6398) oleh @taliesin-ai
- commit header respons fallback sebelum flush dalam [PR](https://github.com/wailsapp/wails/pull/6399) oleh @taliesin-ai
- temukan hanya file indeks yang namanya cocok persis dalam [PR](https://github.com/wailsapp/wails/pull/6400) oleh @taliesin-ai
- hindari body dan pembacaan penuh untuk HEAD yang tidak mendukung seek dalam [PR](https://github.com/wailsapp/wails/pull/6401) oleh @taliesin-ai
- isi prefiks deteksi MIME hingga penuh meski pembacaan berlangsung singkat dalam [PR](https://github.com/wailsapp/wails/pull/6402) oleh @taliesin-ai
- commit respons kosong yang berhasil dan tidak mendukung seek dalam [PR](https://github.com/wailsapp/wails/pull/6403) oleh @taliesin-ai
- hindari duplikasi target URI yang telah diperluas dalam [PR](https://github.com/wailsapp/wails/pull/6405) oleh @taliesin-ai
- pertahankan path selama perluasan field dalam [PR](https://github.com/wailsapp/wails/pull/6406) oleh @taliesin-ai
- patuhi direktori pencarian file desktop XDG dalam [PR](https://github.com/wailsapp/wails/pull/6407) oleh @taliesin-ai
- commit header pendeteksi saat flush kosong dalam [PR](https://github.com/wailsapp/wails/pull/6409) oleh @taliesin-ai
- lewati direktori yang membayangi entri desktop dalam [PR](https://github.com/wailsapp/wails/pull/6411) oleh @taliesin-ai
- pertahankan target yang ditemukan di dalam path executable dalam [PR](https://github.com/wailsapp/wails/pull/6412) oleh @taliesin-ai
- lepaskan variabel yang ditangkap callback setelah timer terpicu dalam [PR](https://github.com/wailsapp/wails/pull/6413) oleh @taliesin-ai
- hindari kecocokan pustaka ldconfig yang keliru dalam [PR](https://github.com/wailsapp/wails/pull/6416) oleh @taliesin-ai
- batasi kecocokan fallback pustaka sistem file dalam [PR](https://github.com/wailsapp/wails/pull/6417) oleh @taliesin-ai
- hilangkan data race pada invalidasi cache dalam [PR](https://github.com/wailsapp/wails/pull/6418) oleh @taliesin-ai
- patuhi konfigurasi logger opsional dalam [PR](https://github.com/wailsapp/wails/pull/6421) oleh @taliesin-ai
- hindari mutasi URL awal bersama dan port kedaluwarsa dalam [PR](https://github.com/wailsapp/wails/pull/6422) oleh @taliesin-ai
- hentikan percobaan ulang proksi pengembangan saat permintaan dibatalkan dalam [PR](https://github.com/wailsapp/wails/pull/6423) oleh @taliesin-ai
- kembalikan 404 untuk aset terbundel yang tidak ditemukan dalam [PR](https://github.com/wailsapp/wails/pull/6424) oleh @taliesin-ai
- tangani permintaan HEAD untuk runtime terbundel dalam [PR](https://github.com/wailsapp/wails/pull/6425) oleh @taliesin-ai
- pertahankan path yang di-escape saat merutekan layanan dalam [PR](https://github.com/wailsapp/wails/pull/6426) oleh @taliesin-ai
- urutkan pengidentifikasi prerelease numerik berukuran sangat besar dengan benar dalam [PR](https://github.com/wailsapp/wails/pull/6427) oleh @taliesin-ai
- tolak galat HTTP sebelum menyimpan unduhan alat dalam [PR](https://github.com/wailsapp/wails/pull/6429) oleh @taliesin-ai
- patuhi tenggat peluncuran selama penemuan desktop dalam [PR](https://github.com/wailsapp/wails/pull/6430) oleh @taliesin-ai
- pertahankan body permintaan selama percobaan ulang proksi pengembangan dalam [PR](https://github.com/wailsapp/wails/pull/6432) oleh @taliesin-ai
- tutup body respons HTTP pada setiap jalur pengembalian dalam [PR](https://github.com/wailsapp/wails/pull/6434) oleh @taliesin-ai
- pertahankan galat pembatalan dalam pencarian pustaka paralel dalam [PR](https://github.com/wailsapp/wails/pull/6435) oleh @taliesin-ai
- tolak respons tag GitHub yang gagal sebelum decoding dalam [PR](https://github.com/wailsapp/wails/pull/6436) oleh @taliesin-ai
- pertahankan path templat literal yang diapit tanda kurung siku dalam [PR](https://github.com/wailsapp/wails/pull/6437) oleh @taliesin-ai
- instal wails3 satu kali dan pertahankan galat penginstal dalam [PR](https://github.com/wailsapp/wails/pull/6438) oleh @taliesin-ai
- abaikan entri tersembunyi dan header dalam blok komentar dalam [PR](https://github.com/wailsapp/wails/pull/6439) oleh @taliesin-ai
- tampilkan galat pencarian rilis stabil dalam [PR](https://github.com/wailsapp/wails/pull/6440) oleh @taliesin-ai
- cetak pesan polos tanpa interpretasi format dalam [PR](https://github.com/wailsapp/wails/pull/6441) oleh @taliesin-ai
- hindari panic saat menulis setelah memuat null dari JSON dalam [PR](https://github.com/wailsapp/wails/pull/6442) oleh @taliesin-ai
- kembalikan snapshot untuk pembacaan seluruh penyimpanan dalam [PR](https://github.com/wailsapp/wails/pull/6444) oleh @taliesin-ai
- propagasikan galat setelah iterasi sebagian baris dalam [PR](https://github.com/wailsapp/wails/pull/6446) oleh @taliesin-ai
- pertahankan path tujuan ekstraksi secara literal dalam [PR](https://github.com/wailsapp/wails/pull/6447) oleh @taliesin-ai
- lepaskan handle prepared statement saat ditutup dan dibuka kembali dalam [PR](https://github.com/wailsapp/wails/pull/6448) oleh @taliesin-ai
- tangani race antara penyiapan dan penutupan atau pembukaan kembali dalam [PR](https://github.com/wailsapp/wails/pull/6450) oleh @taliesin-ai
- patuhi permintaan HEAD untuk halaman indeks fallback dalam [PR](https://github.com/wailsapp/wails/pull/6451) oleh @taliesin-ai
- wajibkan handle statement heksadesimal yang lengkap dalam [PR](https://github.com/wailsapp/wails/pull/6452) oleh @taliesin-ai
- pertahankan basis data dalam memori default saat terjadi konkurensi dalam [PR](https://github.com/wailsapp/wails/pull/6456) oleh @taliesin-ai
- pisahkan stderr dari output HEAD yang di-parse dalam [PR](https://github.com/wailsapp/wails/pull/6457) oleh @taliesin-ai
- hasilkan aset runtime dari bundel tersemat dalam [PR](https://github.com/wailsapp/wails/pull/6458) oleh @taliesin-ai
- deteksi placeholder TODO, FIXME, dan TBD dalam [PR](https://github.com/wailsapp/wails/pull/6459) oleh @taliesin-ai
- validasi respons API catatan rilis dalam [PR](https://github.com/wailsapp/wails/pull/6461) oleh @taliesin-ai
- patuhi direktori tugas dalam pemeriksaan status dalam [PR](https://github.com/wailsapp/wails/pull/6462) oleh @taliesin-ai
- cegah kehilangan data saat menyalin file ke file yang sama dalam [PR](https://github.com/wailsapp/wails/pull/6463) oleh @taliesin-ai
- sertakan direktori tugas dalam identitas cache dalam [PR](https://github.com/wailsapp/wails/pull/6464) oleh @taliesin-ai
- probe port satu kali per pemeriksaan dalam [PR](https://github.com/wailsapp/wails/pull/6465) oleh @taliesin-ai
- hindari panic saat entri cache didekode menjadi null dalam [PR](https://github.com/wailsapp/wails/pull/6466) oleh @taliesin-ai
- naikkan prerelease numerik tanpa batas integer dalam [PR](https://github.com/wailsapp/wails/pull/6467) oleh @taliesin-ai
- mencegah pemeriksaan lama menyembunyikan cache miss dalam [PR](https://github.com/wailsapp/wails/pull/6468) oleh @taliesin-ai
- menjalankan tugas default di direktori dasar executor dalam [PR](https://github.com/wailsapp/wails/pull/6469) oleh @taliesin-ai
- mempertahankan tanda persen literal dalam perintah explorer dalam [PR](https://github.com/wailsapp/wails/pull/6470) oleh @taliesin-ai
- menyelaraskan direktori perintah prasyarat dan tugas dalam [PR](https://github.com/wailsapp/wails/pull/6471) oleh @taliesin-ai
- mengganti symlink dengan file reguler yang dihasilkan dalam [PR](https://github.com/wailsapp/wails/pull/6474) oleh @taliesin-ai
- mempertahankan DMG yang ada saat build gagal dalam [PR](https://github.com/wailsapp/wails/pull/6476) oleh @taliesin-ai
- mempertahankan regex dalam tanda kutip dan garis miring terbalik pada path dalam [PR](https://github.com/wailsapp/wails/pull/6477) oleh @taliesin-ai
- memvalidasi ukuran ikon yang dihasilkan sebelum alokasi gambar dalam [PR](https://github.com/wailsapp/wails/pull/6478) oleh @taliesin-ai
- meneruskan rangkaian perintah dan pengalihan melalui shell dalam [PR](https://github.com/wailsapp/wails/pull/6479) oleh @taliesin-ai
- mempertahankan paket Linux sebelumnya saat build gagal dalam [PR](https://github.com/wailsapp/wails/pull/6480) oleh @taliesin-ai
- menggunakan eksekusi shell fallback khusus platform dalam [PR](https://github.com/wailsapp/wails/pull/6481) oleh @taliesin-ai
- menyajikan aset dengan garis miring awal path yang berulang dalam [PR](https://github.com/wailsapp/wails/pull/6483) oleh @taliesin-ai
- mengenali versi stabil yang memiliki metadata build dalam [PR](https://github.com/wailsapp/wails/pull/6484) oleh @taliesin-ai
- melaporkan kegagalan pembuatan sumber daya secara akurat dalam [PR](https://github.com/wailsapp/wails/pull/6486) oleh @taliesin-ai
- mencegah race condition saat mengganti reporter aktif dalam [PR](https://github.com/wailsapp/wails/pull/6487) oleh @taliesin-ai
- menyerialkan penyimpanan file cache secara bersamaan dalam [PR](https://github.com/wailsapp/wails/pull/6488) oleh @taliesin-ai
- mendukung beberapa kandidat kesetaraan templat dalam [PR](https://github.com/wailsapp/wails/pull/6489) oleh @taliesin-ai
- menormalisasi peta kategori tersimpan yang bernilai null dalam [PR](https://github.com/wailsapp/wails/pull/6491) oleh @taliesin-ai
- menampilkan diagnostik dari perintah shell variabel yang gagal dalam [PR](https://github.com/wailsapp/wails/pull/6492) oleh @taliesin-ai
- mengecualikan rangkaian perintah shell dari cache Go implisit dalam [PR](https://github.com/wailsapp/wails/pull/6493) oleh @taliesin-ai
- menyerialkan persistensi kategori Linux dalam [PR](https://github.com/wailsapp/wails/pull/6494) oleh @taliesin-ai
- menyimpan output Go dengan spasi dalam tanda kutip atau yang di-escape ke cache dalam [PR](https://github.com/wailsapp/wails/pull/6495) oleh @taliesin-ai
- mempertahankan ikon Mac saat pengodean gagal dalam [PR](https://github.com/wailsapp/wails/pull/6496) oleh @taliesin-ai
- berhenti mengabaikan flag jumlah pengujian Go yang tidak valid dalam [PR](https://github.com/wailsapp/wails/pull/6497) oleh @taliesin-ai
- mempertahankan izin sumber untuk salinan file baru dalam [PR](https://github.com/wailsapp/wails/pull/6498) oleh @taliesin-ai
- Memperbaiki berbagai masalah build nightly v3 dan merapikan dependensi dalam [PR](https://github.com/wailsapp/wails/pull/6519) oleh @taliesin-ai
- memperbaiki validasi rilis dan klasifikasi fixture v3 dalam [PR](https://github.com/wailsapp/wails/pull/6533) oleh @taliesin-ai
- mengatasi kegagalan validasi rilis lintas platform dalam [PR](https://github.com/wailsapp/wails/pull/6538) oleh @taliesin-ai

## v3.1.0-nightly.20261028 - 2026-10-28

## Diubah
- Meningkatkan kecepatan pembuatan changelog dengan menggunakan kembali pola regex yang telah dikompilasi dalam [PR](https://github.com/wailsapp/wails/pull/6490) oleh @taliesin-ai
- Mengurangi alokasi saat menyajikan halaman fallback yang dilokalkan dalam [PR](https://github.com/wailsapp/wails/pull/6485) oleh @taliesin-ai
- Mempercepat penemuan runtime AppImage dengan menghindari pencarian metadata yang tidak diperlukan dalam [PR](https://github.com/wailsapp/wails/pull/6482) oleh @taliesin-ai
- Menggunakan kembali pola perutean perintah yang telah dikompilasi untuk meningkatkan performa wake dalam [PR](https://github.com/wailsapp/wails/pull/6475) oleh @taliesin-ai
- Meningkatkan performa dengan menghindari pemindaian ulang baris output yang belum selesai dalam [PR](https://github.com/wailsapp/wails/pull/6473) oleh @taliesin-ai
- Mengalirkan penyalinan executable selama penyiapan paket MSIX untuk mengurangi penggunaan memori dan mempercepat build dalam [PR](https://github.com/wailsapp/wails/pull/6472) oleh @taliesin-ai
- Menghindari penyalinan baris output perintah yang belum selesai untuk meningkatkan performa dalam [PR](https://github.com/wailsapp/wails/pull/6460) oleh @taliesin-ai
- Mengurangi alokasi saat menghasilkan UUID berbasis nama dalam [PR](https://github.com/wailsapp/wails/pull/6455) oleh @taliesin-ai
- Mengurangi alokasi saat memformat handle statement SQLite dalam [PR](https://github.com/wailsapp/wails/pull/6454) oleh @taliesin-ai
- Mengurangi konversi UUID menjadi string hingga hanya memerlukan satu alokasi dalam [PR](https://github.com/wailsapp/wails/pull/6453) oleh @taliesin-ai
- Mengurangi alokasi saat memformat entri poin changelog dalam [PR](https://github.com/wailsapp/wails/pull/6449) oleh @taliesin-ai
- Mengurangi alokasi saat memformat hyperlink terminal dalam [PR](https://github.com/wailsapp/wails/pull/6443) oleh @taliesin-ai
- Mengurangi alokasi saat menyajikan aset dengan mendeteksi buffer MIME hanya ketika diperlukan dalam [PR](https://github.com/wailsapp/wails/pull/6428) oleh @taliesin-ai
- Melewati persiapan log debug di server aset saat pencatatan debug dinonaktifkan dalam [PR](https://github.com/wailsapp/wails/pull/6420) oleh @taliesin-ai
- Mengurangi alokasi dalam penanganan permintaan webview dengan menghindari konteks permintaan yang dibuang dalam [PR](https://github.com/wailsapp/wails/pull/6419) oleh @taliesin-ai
- Mempercepat penyajian aset dengan melewati deteksi MIME untuk ekstensi file yang dikenal dalam [PR](https://github.com/wailsapp/wails/pull/6414) oleh @taliesin-ai
- Mengurangi alokasi saat menyajikan aset dengan menghindari penggunaan reader untuk mendeteksi konten dalam [PR](https://github.com/wailsapp/wails/pull/6410) oleh @taliesin-ai
- Mengurangi alokasi saat menyajikan file aset yang tidak mendukung seek dalam [PR](https://github.com/wailsapp/wails/pull/6408) oleh @taliesin-ai
- Mempercepat pemilihan root indeks server aset dengan menghindari pengumpulan kandidat dalam [PR](https://github.com/wailsapp/wails/pull/6404) oleh @taliesin-ai
- Meningkatkan performa server aset dengan memformat panjang konten tanpa fmt dalam [PR](https://github.com/wailsapp/wails/pull/6393) oleh @taliesin-ai
- Mengurangi penyalinan memori saat membaca metadata file explorer desktop dalam [PR](https://github.com/wailsapp/wails/pull/6391) oleh @taliesin-ai

## Diperbaiki
- Mempertahankan izin file sumber saat menyalin ke file tujuan baru dalam [PR](https://github.com/wailsapp/wails/pull/6498) oleh @taliesin-ai
- Mencegah jendela layar penuh dan termaksimalkan berubah ukuran saat batasan ukuran minimum atau maksimum berubah dalam [PR](https://github.com/wailsapp/wails/pull/6388) oleh @taliesin-ai
- Menerbitkan pembaruan terjadwal riwayat bintang untuk repositori wailsapp/wails dalam [PR](https://github.com/wailsapp/wails/pull/6500) oleh @taliesin-ai
- Berhenti mengabaikan flag jumlah pengujian Go yang tidak valid di wake dalam [PR](https://github.com/wailsapp/wails/pull/6497) oleh @taliesin-ai
- Mempertahankan ikon Mac saat pengodean gagal dalam [PR](https://github.com/wailsapp/wails/pull/6496) oleh @taliesin-ai
- Memperbaiki cache miss build untuk output Go dengan spasi dalam tanda kutip atau yang di-escape dalam [PR](https://github.com/wailsapp/wails/pull/6495) oleh @taliesin-ai
- Menyerialkan persistensi kategori notifikasi di Linux untuk mencegah race condition penulisan serentak dalam [PR](https://github.com/wailsapp/wails/pull/6494) oleh @taliesin-ai
- Mengecualikan rangkaian perintah shell dari cache build Go implisit dalam [PR](https://github.com/wailsapp/wails/pull/6493) oleh @taliesin-ai
- Perbaiki kegagalan pemuatan peta kategori notifikasi ketika data yang disimpan bernilai null dalam [PR](https://github.com/wailsapp/wails/pull/6491) oleh @taliesin-ai
- Tampilkan diagnostik dari perintah shell variabel yang gagal dalam [PR](https://github.com/wailsapp/wails/pull/6492) oleh @taliesin-ai
- Perbaiki pencocokan templat wake ketika terdapat beberapa kandidat kesetaraan dalam [PR](https://github.com/wailsapp/wails/pull/6489) oleh @taliesin-ai
- Serialkan penyimpanan file cache serentak untuk mencegah penulisan cache yang rusak dalam [PR](https://github.com/wailsapp/wails/pull/6488) oleh @taliesin-ai
- Cegah race condition saat mengganti pelapor aktif dalam [PR](https://github.com/wailsapp/wails/pull/6487) oleh @taliesin-ai
- Laporkan kegagalan pembuatan sumber daya secara akurat di CLI dalam [PR](https://github.com/wailsapp/wails/pull/6486) oleh @taliesin-ai
- Kenali versi stabil yang menyertakan metadata build di CLI dalam [PR](https://github.com/wailsapp/wails/pull/6484) oleh @taliesin-ai
- Sajikan aset yang jalur permintaannya diawali garis miring berulang dalam [PR](https://github.com/wailsapp/wails/pull/6483) oleh @taliesin-ai
- Perbaiki eksekusi fallback shell untuk wake pada setiap platform dalam [PR](https://github.com/wailsapp/wails/pull/6481) oleh @taliesin-ai
- Pertahankan paket Linux sebelumnya ketika build gagal dalam [PR](https://github.com/wailsapp/wails/pull/6480) oleh @taliesin-ai
- Arahkan rangkaian perintah dan pengalihan melalui shell dalam [PR](https://github.com/wailsapp/wails/pull/6479) oleh @taliesin-ai
- Pertahankan regex dalam tanda kutip dan garis miring terbalik pada jalur saat mengurai perintah wake dalam [PR](https://github.com/wailsapp/wails/pull/6477) oleh @taliesin-ai
- Pertahankan file DMG yang ada ketika build pemaketan gagal dalam [PR](https://github.com/wailsapp/wails/pull/6476) oleh @taliesin-ai
- Ganti symlink dengan file biasa yang dihasilkan dalam keluaran binding dalam [PR](https://github.com/wailsapp/wails/pull/6474) oleh @taliesin-ai
- Pertahankan tanda persen literal dalam perintah penjelajah file Linux dalam [PR](https://github.com/wailsapp/wails/pull/6470) oleh @taliesin-ai
- Selaraskan direktori prasyarat dan perintah tugas agar keduanya diresolusikan secara konsisten dalam [PR](https://github.com/wailsapp/wails/pull/6471) oleh @taliesin-ai
- Jalankan tugas bawaan di direktori dasar eksekutor dalam [PR](https://github.com/wailsapp/wails/pull/6469) oleh @taliesin-ai
- Cegah pemeriksaan lama menyembunyikan cache miss dalam [PR](https://github.com/wailsapp/wails/pull/6468) oleh @taliesin-ai
- Perbaiki penomoran prarilis CLI agar dapat bertambah tanpa batas overflow bilangan bulat dalam [PR](https://github.com/wailsapp/wails/pull/6467) oleh @taliesin-ai
- Hindari panic ketika entri cache didekode sebagai null dalam [PR](https://github.com/wailsapp/wails/pull/6466) oleh @taliesin-ai
- Periksa port satu kali untuk setiap pemeriksaan di CLI dalam [PR](https://github.com/wailsapp/wails/pull/6465) oleh @taliesin-ai
- Sertakan direktori tugas dalam identitas cache untuk menghindari hasil wake yang kedaluwarsa dalam [PR](https://github.com/wailsapp/wails/pull/6464) oleh @taliesin-ai
- Cegah kehilangan data ketika proses pemaketan menyalin file ke dirinya sendiri dalam [PR](https://github.com/wailsapp/wails/pull/6463) oleh @taliesin-ai
- Pastikan pemeriksaan status mematuhi direktori tugas dalam [PR](https://github.com/wailsapp/wails/pull/6462) oleh @taliesin-ai
- Deteksi placeholder TODO, FIXME, dan TBD dalam entri changelog dalam [PR](https://github.com/wailsapp/wails/pull/6459) oleh @taliesin-ai
- Hasilkan aset runtime dari bundel tersemat di CLI dalam [PR](https://github.com/wailsapp/wails/pull/6458) oleh @taliesin-ai
- Pisahkan stderr dari keluaran HEAD yang diurai dalam perintah git dalam [PR](https://github.com/wailsapp/wails/pull/6457) oleh @taliesin-ai
- Pertahankan basis data SQLite dalam memori bawaan saat diakses secara serentak dalam [PR](https://github.com/wailsapp/wails/pull/6456) oleh @taliesin-ai
- Kembalikan respons yang benar untuk permintaan HEAD pada halaman indeks fallback dalam [PR](https://github.com/wailsapp/wails/pull/6451) oleh @taliesin-ai
- Tangani race condition antara penyiapan statement SQLite dan penutupan atau pembukaan kembali basis data dalam [PR](https://github.com/wailsapp/wails/pull/6450) oleh @taliesin-ai
- Lepaskan handle prepared statement ketika basis data SQLite ditutup atau dibuka kembali dalam [PR](https://github.com/wailsapp/wails/pull/6448) oleh @taliesin-ai
- Pertahankan jalur tujuan ekstraksi secara literal ketika menggunakan gosod dalam [PR](https://github.com/wailsapp/wails/pull/6447) oleh @taliesin-ai
- Propagasikan galat dari kueri SQLite setelah iterasi baris parsial dalam [PR](https://github.com/wailsapp/wails/pull/6446) oleh @taliesin-ai
- Kembalikan snapshot untuk pembacaan seluruh penyimpanan di kvstore dalam [PR](https://github.com/wailsapp/wails/pull/6444) oleh @taliesin-ai
- Cegah panic saat menulis di kvstore setelah memuat nilai JSON null dalam [PR](https://github.com/wailsapp/wails/pull/6442) oleh @taliesin-ai
- Cetak pesan terminal teks biasa tanpa menafsirkan verba format dalam [PR](https://github.com/wailsapp/wails/pull/6441) oleh @taliesin-ai
- Tampilkan galat saat mencari rilis stabil terbaru di updater dalam [PR](https://github.com/wailsapp/wails/pull/6440) oleh @taliesin-ai
- Abaikan entri tersembunyi dan header dalam blok komentar changelog dalam [PR](https://github.com/wailsapp/wails/pull/6439) oleh @taliesin-ai
- Instal wails3 satu kali dan tampilkan galat penginstal di updater dalam [PR](https://github.com/wailsapp/wails/pull/6438) oleh @taliesin-ai
- Pertahankan jalur templat literal yang diapit tanda kurung siku di gosod dalam [PR](https://github.com/wailsapp/wails/pull/6437) oleh @taliesin-ai
- Tolak respons tag GitHub yang gagal sebelum mendekodenya di updater dalam [PR](https://github.com/wailsapp/wails/pull/6436) oleh @taliesin-ai
- Pertahankan galat pembatalan selama pencarian pustaka paralel dalam [PR](https://github.com/wailsapp/wails/pull/6435) oleh @taliesin-ai
- Tutup body respons HTTP pada setiap jalur pengembalian untuk mencegah kebocoran sumber daya dalam [PR](https://github.com/wailsapp/wails/pull/6434) oleh @taliesin-ai
- Pertahankan body permintaan ketika proksi pengembangan mencoba ulang permintaan dalam [PR](https://github.com/wailsapp/wails/pull/6432) oleh @taliesin-ai
- Patuhi tenggat peluncuran selama penemuan penjelajah file desktop dalam [PR](https://github.com/wailsapp/wails/pull/6430) oleh @taliesin-ai
- Tolak galat HTTP sebelum menyimpan unduhan alat AppImage dalam [PR](https://github.com/wailsapp/wails/pull/6429) oleh @taliesin-ai
- Urutkan pengidentifikasi prarilis numerik berukuran sangat besar dengan benar dalam perbandingan semver dalam [PR](https://github.com/wailsapp/wails/pull/6427) oleh @taliesin-ai
- Pertahankan jalur yang di-escape saat mengarahkan layanan di server aset dalam [PR](https://github.com/wailsapp/wails/pull/6426) oleh @taliesin-ai
- Tangani permintaan HEAD untuk runtime yang dibundel di server aset dalam [PR](https://github.com/wailsapp/wails/pull/6425) oleh @taliesin-ai
- Kembalikan 404 untuk aset bundel yang tidak ditemukan, bukan respons yang keliru, dalam [PR](https://github.com/wailsapp/wails/pull/6424) oleh @taliesin-ai
- Hentikan percobaan ulang proksi pengembangan ketika permintaan dibatalkan dalam [PR](https://github.com/wailsapp/wails/pull/6423) oleh @taliesin-ai
- Hindari mutasi URL awal bersama dan port kedaluwarsa di server aset dalam [PR](https://github.com/wailsapp/wails/pull/6422) oleh @taliesin-ai
- Patuhi konfigurasi logger opsional di server aset dalam [PR](https://github.com/wailsapp/wails/pull/6421) oleh @taliesin-ai
- Hilangkan data race dalam invalidasi cache jalur pustaka dalam [PR](https://github.com/wailsapp/wails/pull/6418) oleh @taliesin-ai
- Hindari kecocokan pustaka palsu saat mencari keluaran ldconfig dalam [PR](https://github.com/wailsapp/wails/pull/6416) oleh @taliesin-ai
- Cegah pertumbuhan dan osilasi penyusutan antrean server aset saat dibebani dalam [PR](https://github.com/wailsapp/wails/pull/6415) oleh @taliesin-ai
- Lepaskan capture callback setelah timer debounce terpicu agar referensi kedaluwarsa tidak dipertahankan dalam [PR](https://github.com/wailsapp/wails/pull/6413) oleh @taliesin-ai
- Pertahankan target penjelajah file yang ditemukan di dalam jalur executable dalam [PR](https://github.com/wailsapp/wails/pull/6412) oleh @taliesin-ai
- Lewati direktori yang menutupi entri desktop di penjelajah berkas dalam [PR](https://github.com/wailsapp/wails/pull/6411) oleh @taliesin-ai
- Perbaiki header server aset yang tidak dikirim saat isi respons kosong dalam [PR](https://github.com/wailsapp/wails/pull/6409) oleh @taliesin-ai
- Patuhi direktori pencarian berkas desktop XDG saat mencari ikon aplikasi dalam [PR](https://github.com/wailsapp/wails/pull/6407) oleh @taliesin-ai
- Pertahankan jalur berkas selama ekspansi bidang di penjelajah berkas dalam [PR](https://github.com/wailsapp/wails/pull/6406) oleh @taliesin-ai
- Hindari duplikasi target URI yang telah diekspansi di penjelajah berkas dalam [PR](https://github.com/wailsapp/wails/pull/6405) oleh @taliesin-ai
- Komit respons aset non-seekable kosong yang berhasil alih-alih membiarkannya tertunda dalam [PR](https://github.com/wailsapp/wails/pull/6403) oleh @taliesin-ai
- Isi prefiks pendeteksian MIME dengan benar pada pembacaan singkat berturut-turut di server aset dalam [PR](https://github.com/wailsapp/wails/pull/6402) oleh @taliesin-ai
- Hindari membaca seluruh isi respons untuk permintaan HEAD non-seekable di server aset dalam [PR](https://github.com/wailsapp/wails/pull/6401) oleh @taliesin-ai
- Sajikan hanya berkas indeks yang sama persis saat menyelesaikan permintaan direktori di server aset dalam [PR](https://github.com/wailsapp/wails/pull/6400) oleh @taliesin-ai
- Perbaiki notifikasi Linux yang ditutup tetapi dilaporkan sebagai klik pada server tanpa dukungan aksi dalam [PR](https://github.com/wailsapp/wails/pull/6367) oleh @bercedev
- Hentikan proses pengembangan saat terminal ditutup dalam [PR](https://github.com/wailsapp/wails/pull/6371) oleh @xmasdev
- Pertahankan debounce mulai ulang taskbar Windows selama satu detik sebelum menghapus flag dalam [PR](https://github.com/wailsapp/wails/pull/6386) oleh @taliesin-ai
- Perbaiki header respons fallback yang hilang di server aset dalam [PR](https://github.com/wailsapp/wails/pull/6399) oleh @taliesin-ai
- Terapkan penghapusan header sebelum memilih fallback di server aset dalam [PR](https://github.com/wailsapp/wails/pull/6398) oleh @taliesin-ai
- Lepaskan referensi antrean cincin yang telah dihapus di server aset untuk mencegah kebocoran memori dalam [PR](https://github.com/wailsapp/wails/pull/6397) oleh @taliesin-ai
- Lepaskan buffer pendeteksi konten yang telah digunakan di server aset untuk mengurangi penggunaan memori dalam [PR](https://github.com/wailsapp/wails/pull/6396) oleh @taliesin-ai
- Kenali ekstensi berkas aset dengan kombinasi huruf besar-kecil saat menyajikan aset dalam [PR](https://github.com/wailsapp/wails/pull/6395) oleh @taliesin-ai
- Deteksi jenis konten dengan benar di server aset tanpa menyimpan nama berkas dalam [PR](https://github.com/wailsapp/wails/pull/6394) oleh @taliesin-ai
- Gunakan nilai default sebagai fallback saat entri desktop tidak memiliki perintah dalam [PR](https://github.com/wailsapp/wails/pull/6392) oleh @taliesin-ai
- Cegah contoh liquid-glass menampilkan dialog sebelum aplikasi dimulai pada platform selain macOS dalam [PR](https://github.com/wailsapp/wails/pull/6389) oleh @taliesin-ai
- Pertahankan batasan ukuran jendela saat mode layar penuh dimulai dari ukuran terbatas di Windows dalam [PR](https://github.com/wailsapp/wails/pull/6387) oleh @taliesin-ai
- Perbaiki dialog yang dipanggil sebelum Run agar kembali alih-alih terblokir selamanya dalam [PR](https://github.com/wailsapp/wails/pull/6390) oleh @taliesin-ai

## Keamanan
- Validasi ukuran ikon yang dihasilkan sebelum alokasi gambar untuk mencegah penggunaan memori berlebihan dalam [PR](https://github.com/wailsapp/wails/pull/6478) oleh @taliesin-ai
- Validasi respons API catatan rilis sebelum digunakan dalam [PR](https://github.com/wailsapp/wails/pull/6461) oleh @taliesin-ai
- Wajibkan handle pernyataan heksadesimal yang lengkap di plugin SQLite dalam [PR](https://github.com/wailsapp/wails/pull/6452) oleh @taliesin-ai
- Batasi kecocokan fallback pustaka sistem berkas untuk mencegah pemuatan pustaka yang tidak dimaksudkan dalam [PR](https://github.com/wailsapp/wails/pull/6417) oleh @taliesin-ai

## v3.1.0-nightly.20261027 - 2026-10-27

## Diubah
- Tautkan halaman 404 ke halaman dokumentasi Wails v2 yang sesuai di v2.wails.io dalam [PR](https://github.com/wailsapp/wails/pull/6366) oleh @taliesin-ai
- Pindahkan situs dokumentasi v2 ke v2.wails.io dalam [PR](https://github.com/wailsapp/wails/pull/6365) oleh @taliesin-ai
- Sajikan dokumentasi v3 dari wails.io dan dokumentasi v2 dari v2.wails.io dalam [PR](https://github.com/wailsapp/wails/pull/6364) oleh @taliesin-ai
- Perbarui catatan rilis v3.0.0 dan koreksi klaim advis keamanan dalam [PR](https://github.com/wailsapp/wails/pull/6363) oleh @taliesin-ai

## Diperbaiki
- Kirim callback listener peristiwa Go secara berurutan sesuai urutan emisi dalam [PR](https://github.com/wailsapp/wails/pull/6275) oleh @DevLumuz
- Selalu publikasikan properti Menu pada tray agar tray yang dibuat setelah Run tidak diabaikan di Linux dalam [PR](https://github.com/wailsapp/wails/pull/6284) oleh @AlbinoGeek
- Tahan pekerjaan thread utama yang dikirim sebelum Run hingga aplikasi dimulai dalam [PR](https://github.com/wailsapp/wails/pull/6287) oleh @AlbinoGeek
- Gunakan URL API aset yang benar saat mengunduh aset rilis pemutakhiran GitHub dalam [PR](https://github.com/wailsapp/wails/pull/6304) oleh @gavidroselj
- Perbaiki batas dan posisi layar pada konfigurasi macOS dan Linux dengan DPI campuran dalam [PR](https://github.com/wailsapp/wails/pull/6346) oleh @julianstorer
- Pertahankan isi permintaan di Windows saat WebView2 tidak menyertakan header Content-Length dalam [PR](https://github.com/wailsapp/wails/pull/6326) oleh @strraaa
- Perbaiki data race dalam emisi hook peristiwa dengan mengambil snapshot hook saat terkunci dalam [PR](https://github.com/wailsapp/wails/pull/6276) oleh @DevLumuz
- Perbaiki `SystemTray.SetTooltip()` yang tidak melakukan apa pun di Linux: deskripsi tooltip StatusNotifierItem kini ditetapkan (#6017)
- Jalankan pemeriksaan bukti changelog dalam alur kerja rilis mingguan dalam [PR](https://github.com/wailsapp/wails/pull/6374) oleh @taliesin-ai
- Perbaiki pemeriksaan bukti changelog rilis agar memetakan proses ke commit squash-merge dalam [PR](https://github.com/wailsapp/wails/pull/6373) oleh @taliesin-ai
- Perbaiki deteksi mode gelap Linux yang menampilkan mode terang pada desktop gelap hingga skema warna berubah, serta `linux:SystemThemeChanged` yang tidak pernah dipicu pada backend GTK4. Kedua backend kini menggunakan satu implementasi bersama yang didukung oleh namespace portal `org.freedesktop.appearance` (#6072)
- `Window.SetZoom` tidak lagi menaikkan faktor di bawah 1.0 menjadi 1.0 pada jendela berbingkai (Windows, macOS, Linux): WebView2 `ZoomFactor`, WKWebView `magnification`, dan tingkat zoom WebKitGTK menerima nilai positif apa pun, sedangkan opsi jendela `Zoom` sudah meneruskan nilai tersebut. Jendela tanpa bingkai tetap mempertahankan batas bawah 1.0 karena wilayah seret dan ubah ukurannya tidak disesuaikan dengan zoom ([#4590](https://github.com/wailsapp/wails/issues/4590)); faktor nonpositif, NaN, dan Inf kini diabaikan alih-alih dinaikkan menjadi 1.0 — oleh @APshenkin

## v3.0.0 - 2026-10-26

## Diubah
- Pra-rilis nightly diterbitkan sebagai `vX.Y.0-nightly.YYYYMMDD` untuk menampilkan pratinjau rilis mingguan berikutnya, dan rilis mingguan menaikkan versi minor (`vX.Y.0`), sehingga `@latest` milik Go selalu mengarah ke rilis mingguan atau hotfix; versi paket npm `@wailsio/runtime` selalu sama dengan versi rilis Wails, dengan nightly menggunakan dist-tag `nightly`
- Templat proyek menyematkan `@wailsio/runtime` ke versi CLI yang menghasilkannya
- Wails v3 memerlukan Go 1.26 atau yang lebih baru (Go 1.26.10 atau 1.27.3 dan versi setelahnya menyertakan perbaikan keamanan `crypto/x509`)

## Diperbaiki
- Perbaiki `wails3 update cli` yang tidak menemukan rilis v3: fitur tersebut hanya mempertimbangkan tag v2 dan menginstal dari jalur yang salah
- Memperbaiki aplikasi Windows yang mengalami crash dengan `panic: SetTextColor failed` saat menggambar bilah menu aplikasi pada beberapa build Windows: bilah menu kini digambar pada konteks perangkat yang valid atau menggunakan penggambaran bawaan sebagai fallback, `UAHMENUITEM.Position` kini sesuai dengan `int` C, dan penggambaran menu tidak pernah menyebabkan panic (#6311)
- Menjaga keamanan pembaruan yang telah disiapkan selama proses mulai ulang tertunda: pengunduhan atau proses mulai ulang kedua yang dimulai saat aplikasi dimatikan tidak dapat menghapus atau memperebutkan payload yang sedang ditunggu oleh pembantu penukaran untuk diinstal, oleh @julianstorer
- Memperbaiki aplikasi Windows yang mengalami crash di `applyDarkMode` ketika inisialisasi `pkg/w32` membongkar `uxtheme.dll` selagi pointer fungsi ordinalnya masih digunakan: modul tersebut kini tetap dimuat selama masa hidup proses (#6356)
- Memvalidasi payload yang telah disiapkan sebelum penukaran dan aplikasi yang telah diinstal sebelum cadangan dibuang, sehingga pembaruan yang rusak akan dikembalikan alih-alih menggantikan aplikasi yang berfungsi, oleh @julianstorer

## Dihapus
- Menghapus integrasi Crowdin setelah sinkronisasi terjemahan terakhir dalam [PR](https://github.com/wailsapp/wails/pull/5736) oleh @github-actions[bot]

## Keamanan
- Memblokir pesan `wails:` yang dicadangkan dari frame anak lintas origin di Linux (GTK3 dan GTK4) dengan menyuntikkan skrip penjaga ke semua frame, sehingga memitigasi GHSA-2p5f-2gx8-mh8m
- Memperbarui `github.com/klauspost/compress` ke 1.18.7 (GHSA-259r-337f-4rfw) dan `golang.org/x/image` ke 0.45.0, serta melakukan build dengan Go 1.26.10 untuk menyertakan perbaikan `crypto/x509` (GO-2026-6809)
- Memperbarui `vite` ke 8.2.4 atau versi yang lebih baru dalam paket runtime, wizard penyiapan, contoh, dan templat proyek (GHSA-9jrq-w75r-8gcw, GHSA-rq7h-c2jc-7f22, GHSA-vfpm-58rq-9qcg, GHSA-r825-w8mm-5pxw, GHSA-587h-7wmj-h5pm, GHSA-698v-4r3x-3vcv)

### Keterbatasan yang diketahui
- **Linux:** penjaga frame anak tidak dapat mencegah frame anak dengan *origin yang sama* (misalnya iframe `srcdoc` atau `about:blank`) mengakses pengendali pesan native melalui induknya karena WebKitGTK tidak memberikan informasi frame kepada pengendali. Frame lintas origin dan frame yang di-sandbox tanpa `allow-same-origin` akan diblokir. Jangan sematkan konten tidak tepercaya dengan origin yang sama di jendela Wails.
- **Android (eksperimental):** bridge JavaScript `window.wails` diekspos ke setiap frame dalam WebView tanpa pemeriksaan frame utama atau origin. Konten dalam iframe tersemat, termasuk konten lintas origin, dapat memanggil metode terikat dan mengirim pesan `wails:` yang dicadangkan. Jangan sematkan konten pihak ketiga yang tidak tepercaya di WebView aplikasi. Perilaku ini berbeda dari macOS dan iOS, yang mengabaikan pesan `wails:` dari frame anak. Perbaikan direncanakan untuk 3.0.1.
- **Semua platform:** panggilan metode terikat dari frame anak dikirim melalui HTTP ke `/wails/runtime`, bukan melalui pengendali pesan native, sehingga tidak dibatasi oleh penjaga. Perlakukan setiap iframe yang Anda sematkan seolah-olah memiliki akses ke metode Go terikat Anda.

## v3.0.0-beta.34 - 2026-10-20

## Diperbaiki
- Memperbaiki webview Windows yang tidak memperoleh kembali fokus saat jendela diaktifkan, serta tombol Tab yang terperangkap di dalam webview alih-alih memindahkan fokus ke jendela native lalu kembali lagi (#3783)
- Memperbaiki struct vtable COM `internal/webview2/pkg/webview2` (`ICoreWebView2_2`..`_27` dan antarmuka turunan lainnya) yang tidak menyertakan rantai metode warisan, sehingga panggilan masuk ke slot vtable yang salah (#6129)
- Memperbaiki dialog pertanyaan Windows yang mengabaikan label tombol khusus dan tidak pernah menjalankan callback tombol; tombol khusus kini ditampilkan dengan TaskDialog (menggunakan MessageBox sebagai fallback dengan pemetaan callback berdasarkan posisi ketika comctl32 v6 tidak tersedia) (#5794)
- Memperbaiki `wails3 dev` yang macet selama 60 detik saat digunakan dengan Vite: gerbang kesiapan frontend otomatis kini menggunakan probe TCP, bukan HTTP GET yang tidak mengirim header `Accept` dan menerima 404 ([#6209](https://github.com/wailsapp/wails/issues/6209))
- Menjaga menu systray tetap berfungsi ketika diperbarui selagi terbuka dalam [PR](https://github.com/wailsapp/wails/pull/6301) oleh @taliesin-ai
- Menangani error COM WebView2 sementara secara baik dalam operasi runtime pada [PR](https://github.com/wailsapp/wails/pull/6302) oleh @leaanthony

## v3.0.0-beta.33 - 2026-10-19

## Ditambahkan
- Generator binding kini menghasilkan modul `bindings/events` yang mengekspor konstanta nama bertipe untuk setiap event yang dapat ditemukan dan didaftarkan dengan `application.RegisterEvent` menggunakan nama konstanta waktu kompilasi (misalnya `UserLogin` untuk `user:login`), yang dapat digunakan dengan `Events.On`/`Emit` tanpa plugin Vite (#5645)

## Diperbaiki
- Memperbaiki dialog pertanyaan Windows yang mengabaikan label tombol khusus dan tidak pernah menjalankan callback tombol; tombol khusus kini ditampilkan dengan TaskDialog (menggunakan MessageBox sebagai fallback dengan pemetaan callback berdasarkan posisi ketika comctl32 v6 tidak tersedia) (#5794)
- Memperbaiki `wails3 dev` yang macet selama 60 detik saat digunakan dengan Vite: gerbang kesiapan frontend otomatis kini menggunakan probe TCP, bukan HTTP GET yang tidak mengirim header `Accept` dan menerima 404 ([#6209](https://github.com/wailsapp/wails/issues/6209))
- Menjaga menu systray tetap berfungsi ketika diperbarui selagi terbuka dalam [PR](https://github.com/wailsapp/wails/pull/6301) oleh @taliesin-ai
- Menangani error COM WebView2 sementara secara baik dalam operasi runtime pada [PR](https://github.com/wailsapp/wails/pull/6302) oleh @leaanthony

## v3.0.0-beta.32 - 2026-10-18

## Ditambahkan
- Menambahkan header bagian menu native macOS untuk macOS 14+ — lihat [situs dokumentasi](https://v3.wails.io) (20 halaman diperbarui) dalam [PR](https://github.com/wailsapp/wails/pull/6297) oleh @leaanthony

## Diperbaiki
- Membuat SetRelativePosition relatif terhadap area kerja layar di macOS dalam [PR](https://github.com/wailsapp/wails/pull/6299) oleh @leaanthony
- Pembantu pembaru meluncurkan ulang target yang tidak berubah setelah waktu tunggu induk habis atau penggantian gagal dalam [PR](https://github.com/wailsapp/wails/pull/6298) oleh @taliesin-ai
- Mengunci pnpm ke versi 10 untuk dokumentasi berversi dalam [PR](https://github.com/wailsapp/wails/pull/6293) oleh @taliesin-ai

## v3.0.0-beta.31 - 2026-10-18

## Ditambahkan
- Mengekspos ekspor subpath `@wailsio/runtime/cancellable` agar `CancellablePromise` dapat diimpor tanpa mengevaluasi titik masuk runtime root (#6153)

## Diubah
- Mewajibkan Go 1.26 atau versi yang lebih baru: modul v3 kini mendeklarasikan `go 1.26.0`, dan proyek yang dihasilkan oleh `wails3 init` serta `wails3 service init` mendeklarasikan versi minimum yang sama

## Diperbaiki
- Memperbaiki keluaran generator binding yang berubah saat melakukan build dengan Go 1.27, ketika tipe `encoding/json` seperti `Marshaler`, `Unmarshaler`, dan `RawMessage` menjadi alias tipe `encoding/json/v2` dan `encoding/json/jsontext`: tipe tersebut kembali dirender sebagai `any` dan tidak ada model yang dihasilkan untuk `encoding/json/v2` (#6292)
- Memvalidasi URL templat jarak jauh yang diteruskan ke `wails3 init -t` sebelum melakukan clone: hanya URL `https://`, `http://`, `ssh://`, `git://`, dan `user@host:path` yang diterima, nilai diteruskan ke git setelah `--`, dan git dibatasi hanya pada transport tersebut selama proses clone
- Menolak nilai header `x-wails-window-id` di luar rentang dalam transport runtime HTTP, alih-alih diam-diam memotongnya menjadi ID milik jendela lain
- Memperkuat ekstraksi arsip pembaru terhadap zip-slip: path entri kini diperiksa dengan `filepath.IsLocal` beserta pemeriksaan bahwa path hasil penggabungan tetap berada di dalam lokasi tujuan, dengan pengujian untuk entri `..` bertingkat, garis miring terbalik, path absolut, dan penulisan melalui symlink dalam arsip zip dan tar.gz
- Memperbaiki AppImage dari `wails3 generate appimage` yang mengalami crash saat dimulai di luar Debian dan Ubuntu dengan pesan "Unable to spawn a new child process: … WebKitNetworkProcess". WebKit yang dibundel kini menjalankan proses pembantu yang disertakan dalam AppImage. Sandbox bubblewrap WebKit tetap aktif
- Memperbaiki build `-tags server` yang gagal dikompilasi untuk Android dan iOS, memperbaiki [#6271](https://github.com/wailsapp/wails/issues/6271)
- Memperbaiki updater yang gagal melakukan pertukaran di Linux ketika direktori sementara berada di sistem berkas lain (seperti tmpfs `/tmp`), dengan menyiapkan berkas di samping biner yang terpasang dan menyalinnya lintas sistem berkas
- Memperbaiki updater yang mencoba mengganti mount hanya-baca milik AppImage, bukan berkas `.AppImage`
- Memperbaiki pesan kontrol jendela `wails:` yang dicadangkan agar tidak diterima dari iframe di macOS dan iOS. Pesan tersebut kini hanya diterima dari frame utama jendela; pesan lain dari iframe tetap mencapai `RawMessageHandler` beserta informasi asalnya
- Menghentikan pengiriman `WindowDidUpdate` dan `WindowWillUpdate` di macOS: `NSWindowDidUpdateNotification` dipicu pada setiap lintasan run loop, sehingga halaman menerima sekitar 120 panggilan `evaluateJavaScript:` kosong per detik, dan di macOS 27.0 setiap panggilan tersebut menghabiskan aktivasi sementara dari klik pengguna—`navigator.clipboard.writeText()` dan API lain yang memerlukan aktivasi menolak klik nyata dengan `NotAllowedError` ([#6170](https://github.com/wailsapp/wails/issues/6170)) oleh @APshenkin
- Memperbaiki bahasa default info.json Windows menjadi en-US (0409) agar Explorer menampilkan detail versi (#6210)
- Memperbaiki perender DMA-BUF Linux yang dinonaktifkan pada laptop dengan dGPU NVIDIA di samping iGPU AMD atau Intel: di Wayland, perender kini hanya dinonaktifkan ketika WebKit merender menggunakan NVIDIA (aplikasi dialihkan ke dGPU atau compositor berjalan di sana), bukan setiap kali modul nvidia dimuat, yang membuat perenderan pada iGPU menjadi sangat lambat; sesi X11 mempertahankan perilaku lama karena GPU perender Xorg ditetapkan melalui konfigurasinya
- Menangani kegagalan penggantian nama lintas perangkat dalam updater Linux dengan fallback penyalinan bertahap dalam [PR](https://github.com/wailsapp/wails/pull/6135) oleh @john-okeefe
- Memperbaiki `WebviewWindow.SetMenu` di Windows yang tidak menggambar ulang bilah menu kecuali tema bilah menu telah dikonfigurasi, memperbaiki [#6103](https://github.com/wailsapp/wails/issues/6103)
- Memperbaiki `WebviewWindow.SetMenu(nil)` yang mengalami panic di Windows alih-alih tidak melakukan apa pun seperti di macOS, memperbaiki [#6104](https://github.com/wailsapp/wails/issues/6104)
- Mempertahankan entri changelog saat terjadi push yang bersaing melalui percobaan ulang dalam [PR](https://github.com/wailsapp/wails/pull/6094) oleh @leaanthony
- Memperbaiki `go mod vendor` yang gagal dengan `pattern arm64/WebView2Loader.dll: no matching files found` di setiap platform dengan menghapus embed yang merujuk pada biner yang tidak pernah disertakan dalam modul, memperbaiki [#5782](https://github.com/wailsapp/wails/issues/5782) dan [#5376](https://github.com/wailsapp/wails/issues/5376), dalam [PR](https://github.com/wailsapp/wails/pull/6031) oleh @Grantmartin2002
- Mencegah pembongkaran jendela ganda di macOS dalam [PR](https://github.com/wailsapp/wails/pull/6159) oleh @charlesbalderson
- Menangani WM_CONTEXTMENU agar menu baki sistem terbuka saat diklik kanan dalam [PR](https://github.com/wailsapp/wails/pull/6162) oleh @wzk2025
- Memperbaiki kelima varian `InvokeSync*` yang memblokir pemanggilnya selamanya ketika callback yang didispatch mengalami panic, memperbaiki [#6107](https://github.com/wailsapp/wails/issues/6107), dalam [PR](https://github.com/wailsapp/wails/pull/6108) oleh @RALIST
- Membangun ulang submenu GTK3 yang telah dihancurkan secara aman dalam [PR](https://github.com/wailsapp/wails/pull/6236) oleh @leaanthony
- Mengimplementasikan shim penerusan worker Windows untuk worker bersarang dalam [PR](https://github.com/wailsapp/wails/pull/6261) oleh @taliesin-ai
- Memperbaiki `OpenFile()` di Windows yang mengabaikan `PromptForMultipleSelection()` ketika `CanChooseDirectories(true)`: pengguna kini dapat memilih beberapa folder sekaligus, seperti di macOS dan Linux (#6263)
- Memperbaiki argumen `double` berdasarkan nilai untuk metode COM WebView2 pada windows/arm64—`ZoomFactor`, `RasterizationScale`, `Expires` cookie, serta skala, ukuran halaman, dan margin pengaturan cetak terdampak di sana (setter pengaturan cetak sebelumnya mengembalikan `ErrDoubleArgUnsupported`): panggilan yang terdampak kini melewati trampoline `runtime.cgocall` yang memuat `d0`, yang tidak dapat dilakukan oleh jalur syscall Go (golang.org/issue/62583) oleh @APshenkin

## Keamanan
- Mode server (`-tags server`) kini hanya menerima panggilan `/wails/runtime` sebagai permintaan `POST` dari origin server sendiri atau origin yang tercantum dalam `ServerOptions.WebSocketOriginPatterns`; metode lain menerima HTTP 405 dan origin browser lain menerima HTTP 403. Deployment di belakang reverse proxy yang menulis ulang header `Host` harus menambahkan origin publiknya ke `WebSocketOriginPatterns`
- Mode server (`-tags server`) yang terikat ke alamat loopback (default `localhost`, `127.0.0.0/8`, atau `::1`) kini menolak permintaan dengan HTTP 403 kecuali header `Host` adalah `localhost`, literal IP, host bind, atau host yang cocok dengan `ServerOptions.WebSocketOriginPatterns`, sehingga memblokir serangan DNS rebinding terhadap `/wails/runtime` dan endpoint WebSocket
- Notifikasi Windows tidak lagi beralih ke PowerShell sebagai fallback ketika API toast WinRT gagal; error kini dikembalikan

## v3.0.0-beta.30 - 2026-10-15

## Diperbaiki
- Memperbaiki updater yang mencoba mengganti mount hanya-baca milik AppImage, bukan berkas `.AppImage`
- Menghentikan pengiriman `WindowDidUpdate` dan `WindowWillUpdate` di macOS: `NSWindowDidUpdateNotification` dipicu pada setiap lintasan run loop, sehingga halaman menerima sekitar 120 panggilan `evaluateJavaScript:` kosong per detik, dan di macOS 27.0 setiap panggilan tersebut menghabiskan aktivasi sementara dari klik pengguna—`navigator.clipboard.writeText()` dan API lain yang memerlukan aktivasi menolak klik nyata dengan `NotAllowedError` ([#6170](https://github.com/wailsapp/wails/issues/6170)) oleh @APshenkin
- Memperbaiki bahasa default info.json Windows menjadi en-US (0409) agar Explorer menampilkan detail versi (#6210)
- Menangani kegagalan penggantian nama lintas perangkat dalam updater Linux dengan fallback penyalinan bertahap dalam [PR](https://github.com/wailsapp/wails/pull/6135) oleh @john-okeefe
- Memperbaiki `WebviewWindow.SetMenu` di Windows yang tidak menggambar ulang bilah menu kecuali tema bilah menu telah dikonfigurasi, memperbaiki [#6103](https://github.com/wailsapp/wails/issues/6103)
- Memperbaiki `WebviewWindow.SetMenu(nil)` yang mengalami panic di Windows alih-alih tidak melakukan apa pun seperti di macOS, memperbaiki [#6104](https://github.com/wailsapp/wails/issues/6104)
- Mempertahankan entri changelog saat terjadi push yang bersaing melalui percobaan ulang dalam [PR](https://github.com/wailsapp/wails/pull/6094) oleh @leaanthony
- Memperbaiki `go mod vendor` yang gagal dengan `pattern arm64/WebView2Loader.dll: no matching files found` di setiap platform dengan menghapus embed yang merujuk pada biner yang tidak pernah disertakan dalam modul, memperbaiki [#5782](https://github.com/wailsapp/wails/issues/5782) dan [#5376](https://github.com/wailsapp/wails/issues/5376), dalam [PR](https://github.com/wailsapp/wails/pull/6031) oleh @Grantmartin2002
- Mencegah pembongkaran jendela ganda di macOS dalam [PR](https://github.com/wailsapp/wails/pull/6159) oleh @charlesbalderson
- Menangani WM_CONTEXTMENU agar menu baki sistem terbuka saat diklik kanan dalam [PR](https://github.com/wailsapp/wails/pull/6162) oleh @wzk2025
- Memperbaiki kelima varian `InvokeSync*` yang memblokir pemanggilnya selamanya ketika callback yang didispatch mengalami panic, memperbaiki [#6107](https://github.com/wailsapp/wails/issues/6107), dalam [PR](https://github.com/wailsapp/wails/pull/6108) oleh @RALIST

## v3.0.0-beta.29 - 2026-10-13

## Ditambahkan
- Pintasan global di macOS kini menyesuaikan dengan tata letak papan ketik: akselerator huruf terikat ke tombol berlabel pada tata letak non-QWERTY (AZERTY, QWERTZ, Dvorak, ...) melalui `UCKeyTranslate`, bukan ke posisi fisik QWERTY yang tetap.

## Diperbaiki
- Menginisialisasi menu native sebelum menetapkannya di Linux dalam [PR](https://github.com/wailsapp/wails/pull/6233) oleh @abhicodes-007
- Mengatasi regresi kesiapan GA macOS dalam penanganan output dan wizard penyiapan dalam [PR](https://github.com/wailsapp/wails/pull/6237) oleh @taliesin-ai
- Memperbaiki wrapper `AddRef`/`Release` COM WebView2 di Windows yang melaporkan last-error thread yang sudah kedaluwarsa sebagai kegagalan. Wrapper tersebut kini mengembalikan jumlah referensi, sehingga pelepasan tidak lagi dapat menghentikan aplikasi atau menggagalkan permintaan aset (#3931)
- Perbaiki `ICoreWebView2Deferral.Release` yang memanggil slot `AddRef`, yang membocorkan penundaan setiap permintaan aset tertunda di Windows (#3931)

## v3.0.0-beta.28 - 2026-10-05

## Diperbaiki
- Simpan output aplikasi macOS di luar penyedia file dalam [PR](https://github.com/wailsapp/wails/pull/6225) oleh @taliesin-ai
- Selesaikan wizard penyiapan dan tampilkan detail sistem dalam [PR](https://github.com/wailsapp/wails/pull/6224) oleh @taliesin-ai
- Tingkatkan penanganan galat pemulihan dan pembersihan sumber daya WebView2 di Windows dalam [PR](https://github.com/wailsapp/wails/pull/6222) oleh @leaanthony
- Tingkatkan pemulihan WebView2 di Windows untuk mencegah aplikasi mengalami crash dan pengontrol terlantar dalam [PR](https://github.com/wailsapp/wails/pull/6220) oleh @taliesin-ai
- Koreksi panduan otorisasi media macOS dalam [PR](https://github.com/wailsapp/wails/pull/6221) oleh @taliesin-ai
- Perbaiki klik kanan pada ikon baki sistem Linux yang menjalankan handler klik sekaligus membuka menu, sehingga jendela yang terpasang beralih status pada setiap klik kanan (#6018)
- Pertahankan handler callback WebView2 tetap aktif hingga pelepasan native di Windows dalam [PR](https://github.com/wailsapp/wails/pull/6184) oleh @taliesin-ai

## v3.0.0-beta.27 - 2026-10-01

## Diperbaiki
- Perbaiki opsi `Permissions` yang diabaikan pada macOS 12 dan versi lebih baru: metode pengambilan media `WKUIDelegate` kini telah diimplementasikan, sehingga `PermissionAllow` dan `PermissionDeny` berlaku untuk permintaan kamera dan mikrofon sebagaimana di Linux dan Windows. `NSCameraUsageDescription` / `NSMicrophoneUsageDescription` serta, jika menggunakan sandbox, entitlement perangkat yang sesuai tetap diperlukan

## v3.0.0-beta.26 - 2026-09-25

## Ditambahkan
- Watcher menerapkan pemeriksaan kesiapan HTTP pada server pengembangan frontend dalam [PR](https://github.com/wailsapp/wails/pull/6150) oleh @atterpac

## Diperbaiki
- Pulihkan aplikasi dari kegagalan proses WebView2 alih-alih membiarkan jendela kosong dalam [PR](https://github.com/wailsapp/wails/pull/6002) oleh @taliesin-ai

## v3.0.0-beta.25 - 2026-09-22

## Diperbaiki
- Hentikan `wails3 dev` dan proses latar belakangnya ketika aplikasi utama keluar (#6048)

## v3.0.0-beta.24 - 2026-09-20

## Keamanan
- Batasi izin GITHUB_TOKEN dalam alur kerja build dan publikasi agar mengikuti prinsip hak akses minimum
- Perbaiki traversal jalur dan keluarnya symlink dari batas yang ditetapkan dalam middleware aset contoh layar
- Pertahankan fungsi instalasi dependensi wizard penyiapan sambil membatasinya pada perintah pengelola paket yang disetujui

## v3.0.0-beta.23 - 2026-09-16

## Ditambahkan

- Izinkan penggantian ID aplikasi GTK di Linux dalam [PR](https://github.com/wailsapp/wails/pull/5972) oleh @overlordtm

## Diperbaiki

- Mekanisme instans tunggal Linux mengklaim nama bus dengan ID aplikasi itu sendiri dalam [PR](https://github.com/wailsapp/wails/pull/5973) oleh @overlordtm
- Batalkan permintaan skema khusus yang dihentikan di Linux dalam [PR](https://github.com/wailsapp/wails/pull/6126) oleh @leaanthony

## v3.0.0-beta.22 - 2026-09-14

## Ditambahkan

- Lengkapi terjemahan dan porting media komunitas — lihat [situs dokumentasi](https://v3.wails.io) (1812 halaman diperbarui) dalam [PR](https://github.com/wailsapp/wails/pull/6122) oleh @leaanthony
- Dokumentasikan tombol Wails yang dapat disematkan dengan berbagai gaya dan ukuran dalam [PR](https://github.com/wailsapp/wails/pull/6120) oleh @leaanthony

## v3.0.0-beta.21 - 2026-09-13

## Ditambahkan

- Sajikan dokumentasi Wails v3 dengan M-Press dalam [PR](https://github.com/wailsapp/wails/pull/6116) oleh @leaanthony

## Diperbaiki

- Urai nilai slug JSON dari frontmatter MPD untuk pembuatan changelog dalam [PR](https://github.com/wailsapp/wails/pull/6118) oleh @leaanthony
- Updater menghapus variabel lingkungan helper dan meluncurkan kembali target asli setelah kegagalan pencadangan dalam [PR](https://github.com/wailsapp/wails/pull/6080) oleh @cnmax
- Mulai handler sinyal default selama App.Run dalam [PR](https://github.com/wailsapp/wails/pull/6098) oleh @leaanthony
- Menu Windows menangani menu nil, membebaskan sumber daya yang diganti, dan menggambar ulang bilah menu dalam [PR](https://github.com/wailsapp/wails/pull/6112) oleh @taliesin-ai
- Pulihkan pemaketan MSIX untuk proyek baru yang menggunakan konfigurasi YAML bersama dalam [PR](https://github.com/wailsapp/wails/pull/6115) oleh @leaanthony
- Perbaiki binding JavaScript dan TypeScript hasil pembuatan yang gagal dimuat ketika pembuat model generik mereferensikan deklarasi helper setelahnya, serta cegah stack overflow saat membuat model generik yang saling bergantung (#6062)

## v3.0.0-beta.20 - 2026-09-10

## Diubah

- Perbarui tautan showcase Clave ke situs web dan repositori saat ini dalam [PR](https://github.com/wailsapp/wails/pull/6082) oleh @01xR4in

## Diperbaiki

- Batalkan permintaan aset Windows yang dihentikan, termasuk permintaan worker, sambil mempertahankan handler keepalive selama navigasi. Teruskan konteks permintaan native melalui wrapper aplikasi pada platform Apple. (#5963, #5969)
- Pertahankan entri changelog ketika terjadi push yang bersaing dengan mekanisme percobaan ulang dalam [PR](https://github.com/wailsapp/wails/pull/6094) oleh @leaanthony
- Perbaiki `go mod vendor` yang selalu gagal dengan `pattern arm64/WebView2Loader.dll: no matching files found` pada setiap platform dengan menghapus sematan yang mereferensikan biner yang tidak pernah disertakan dalam modul, sehingga memperbaiki [#5782](https://github.com/wailsapp/wails/issues/5782) dan [#5376](https://github.com/wailsapp/wails/issues/5376), dalam [PR](https://github.com/wailsapp/wails/pull/6031) oleh @Grantmartin2002

## Dihapus

- Hapus dukungan loader WebView2 native yang telah digantikan oleh loader Go murni. Perubahan ini menghapus biner `WebView2Loader.dll` yang disematkan dan dependensi `github.com/jchv/go-winloader`. Build tag `native_webview2loader` masih diterima dan tidak lagi menyebabkan galat, tetapi tidak berpengaruh pada build v3, dalam [PR](https://github.com/wailsapp/wails/pull/6031) oleh @Grantmartin2002
- Hapus build tag yang tidak digunakan dan opsi FPS dari panduan API macOS dalam [PR](https://github.com/wailsapp/wails/pull/6097) oleh @leaanthony

## v3.0.0-beta.19 - 2026-09-09

## Ditambahkan

- Lindungi API macOS privat dengan build tag agar penggunaannya harus diaktifkan secara eksplisit — lihat [dokumentasi](https://v3.wails.io/features/browser/integration) dan [dokumentasi](https://v3.wails.io/features/environment/info) dan [dokumentasi](https://v3.wails.io/features/windows/basics) dan [dokumentasi](https://v3.wails.io/features/windows/frameless) dan [dokumentasi](https://v3.wails.io/features/windows/notch-windows) dan [dokumentasi](https://v3.wails.io/features/windows/options) dan [dokumentasi](https://v3.wails.io/guides/build/macos) dan [dokumentasi](https://v3.wails.io/guides/build/private-macos-apis) dan [dokumentasi](https://v3.wails.io/reference/overview) dalam [PR](https://github.com/wailsapp/wails/pull/6087) oleh @leaanthony

## Diperbaiki

- Tolak permintaan runtime yang melebihi 64 MiB dengan HTTP 413 dalam [PR](https://github.com/wailsapp/wails/pull/6091) oleh @leaanthony

## Keamanan

- Perkuat pembatasan origin MCP dan akses jarak jauh dengan autentikasi token dalam [PR](https://github.com/wailsapp/wails/pull/6092) oleh @leaanthony

## v3.0.0-beta.18 - 2026-09-08

## Diperbaiki

- Perbaiki kebocoran memori Calloc di Linux dan Darwin dengan menggunakan receiver pointer dalam [PR](https://github.com/wailsapp/wails/pull/6083) oleh @4RH1T3CT0R7

## v3.0.0-beta.17 - 2026-09-06

## Diperbaiki

- Windows: `GetRequest` yang gagal atau nil dalam handler WebResourceRequested tidak lagi menghentikan proses (`log.Fatal` / panic akibat dereferensi nil) — sebagai gantinya, permintaan dibatalkan dan dicatat dalam log melalui [PR](https://github.com/wailsapp/wails/pull/6006) oleh @midagedev

## v3.0.0-beta.16 - 2026-08-29

## Diubah

- Minta kata sandi notarisasi di jendela terminal baru dalam [PR](https://github.com/wailsapp/wails/pull/6029) oleh @leaanthony

## Diperbaiki

- Tangani jenis klik system tray dengan benar di macOS dalam [PR](https://github.com/wailsapp/wails/pull/5919) oleh @ChewbaccaCookie
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

- Tetap layani pekerjaan thread utama di macOS saat loop modal berjalan dalam [PR](https://github.com/wailsapp/wails/pull/6026) oleh @leaanthony
- Jadikan penyimpanan aman seluler dapat mengalami kegagalan dan gagal secara tertutup dalam [PR](https://github.com/wailsapp/wails/pull/5923) oleh @mortenolsrud
- Jalankan hook peristiwa aplikasi meskipun tidak ada listener yang terdaftar dalam [PR](https://github.com/wailsapp/wails/pull/5999) oleh @archy-rock3t-cloud
- Perbaiki kesalahan ketik dalam komentar dan dokumentasi yang dilokalkan dalam [PR](https://github.com/wailsapp/wails/pull/6023) oleh @haoku123
- Hapus biner macOS prakompilasi yang di-commit di bawah `v3/examples` dalam [PR](https://github.com/wailsapp/wails/pull/6025) oleh @4RH1T3CT0R7

## v3.0.0-beta.12 - 2026-08-21

## Ditambahkan

- Tambahkan jendela notifikasi notch macOS beserta contoh siklus hidup dan telemetri — lihat [dokumentasi](https://v3.wails.io/features/windows/notch-windows) dalam [PR](https://github.com/wailsapp/wails/pull/6010) oleh @leaanthony
- Tambahkan dukungan jendela NSPanel macOS dengan opsi baru dan integrasi native — lihat [dokumentasi](https://v3.wails.io/features/windows/options) dalam [PR](https://github.com/wailsapp/wails/pull/6008) oleh @leaanthony

## Diperbaiki

- Cegah SQLite Prepare mengalami hang saat dipanggil secara bersamaan dalam [PR](https://github.com/wailsapp/wails/pull/5998) oleh @archy-rock3t-cloud

## v3.0.0-beta.11 - 2026-08-20

## Dihapus

- Hapus pelacak implementasi yang sudah usang dari dokumentasi dalam [PR](https://github.com/wailsapp/wails/pull/6005) oleh @leaanthony

## v3.0.0-beta.10 - 2026-08-19

## Diperbaiki

- Perbaiki host Linux GTK4 yang menghilangkan argumen peluncuran protokol khusus dan asosiasi file dalam [PR](https://github.com/wailsapp/wails/pull/6000) oleh @midagedev
- Tangani baris yang dihapus dan koreksi dari sumber yang sama dengan benar dalam validasi changelog melalui [PR](https://github.com/wailsapp/wails/pull/5993) oleh @taliesin-ai

## v3.0.0-beta.9 - 2026-08-16

## Ditambahkan

- Tambahkan server wails3 mcp yang aman untuk pengelolaan proyek dengan bantuan agen dalam [PR](https://github.com/wailsapp/wails/pull/5896) oleh @leaanthony
- Tambahkan dokumentasi untuk model dalam binding — lihat [dokumentasi](https://v3.wails.io/features/bindings/models) dalam [PR](https://github.com/wailsapp/wails/pull/5988) oleh @taliesin-ai
- Dukung instalasi rpm-ostree untuk sistem Linux atomik dalam [PR](https://github.com/wailsapp/wails/pull/5987) oleh @leaanthony
- Tambahkan pembuatan dan publikasi native mingguan untuk grafik riwayat bintang — lihat [dokumentasi](https://v3.wails.io/credits) dan [dokumentasi](https://v3.wails.io/de/credits) dan [dokumentasi](https://v3.wails.io/fr/credits) dan [dokumentasi](https://v3.wails.io/id/credits) dan [dokumentasi](https://v3.wails.io/ja/credits) dan [dokumentasi](https://v3.wails.io/ko/credits) dan [dokumentasi](https://v3.wails.io/pt/credits) dan [dokumentasi](https://v3.wails.io/ru/credits) dan [dokumentasi](https://v3.wails.io/zh-cn/credits) dan [dokumentasi](https://v3.wails.io/zh-tw/credits) dalam [PR](https://github.com/wailsapp/wails/pull/5986) oleh @leaanthony
- Tambahkan paket mac khusus Darwin untuk menemukan resource bundle aplikasi — lihat [dokumentasi](https://v3.wails.io/guides/build/macos) dalam [PR](https://github.com/wailsapp/wails/pull/5965) oleh @leaanthony
- Tambahkan halaman showcase Condui dan entri indeksnya — lihat [dokumentasi](https://v3.wails.io/community/showcase/condui) dan [dokumentasi](https://v3.wails.io/community/showcase) dalam [PR](https://github.com/wailsapp/wails/pull/5962) oleh @mgueregath
- Tambahkan halaman showcase Redis Viewer dengan tangkapan layar dan tautan proyek — lihat [dokumentasi](https://v3.wails.io/community/showcase) dan [dokumentasi](https://v3.wails.io/community/showcase/redisviewer) dalam [PR](https://github.com/wailsapp/wails/pull/5984) oleh @redisviewer

## Diubah

- Perbarui flag aplikasi GTK menjadi G<em>APPLICATION</em>NON_UNIQUE untuk Linux dalam [PR](https://github.com/wailsapp/wails/pull/5971) oleh @overlordtm
- Catat peristiwa jendela yang tidak ditemukan pada level debug, bukan warning, dalam [PR](https://github.com/wailsapp/wails/pull/5914) oleh @julianstorer

## Diperbaiki

- Izinkan pintasan papan ketik macOS yang terdaftar diprioritaskan di atas webview dalam [PR](https://github.com/wailsapp/wails/pull/5902) oleh @julianstorer
- Perbaiki tautan sidebar dokumentasi yang rusak dalam [PR](https://github.com/wailsapp/wails/pull/5937) oleh @northes
- Batalkan konteks permintaan aset macOS dan iOS saat WebKit membatalkan tugas skema khusus yang sesuai (#5963)
- Tangani pesan WindowSetFullscreenButtonEnabled dalam [PR](https://github.com/wailsapp/wails/pull/5976) oleh @archy-rock3t-cloud
- Impor Fragment dalam templat preact-ts untuk mengatasi kegagalan build dalam [PR](https://github.com/wailsapp/wails/pull/5979) oleh @haoku123
- Cegah aplikasi lama khusus layanan GTK3 mengalami crash ketika penemuan layar berjalan sebelum jendela atau tampilan aktif tersedia (#5966)
- Izinkan proses rilis dengan versi eksplisit dilanjutkan ketika changelog yang belum dirilis kosong (#5977)

## Keamanan

- Perbarui lockfile nanoid situs web ke 3.3.18 yang telah ditambal untuk mengatasi advis keamanan dalam [PR](https://github.com/wailsapp/wails/pull/5985) oleh @taliesin-ai

## v3.0.0-beta.8 - 2026-08-12

## Ditambahkan

- Tambahkan pembuatan URL dokumentasi ke entri changelog otomatis dalam [PR](https://github.com/wailsapp/wails/pull/5957) oleh @taliesin-ai
- Tambahkan Streams: aliran byte dua arah antara Go dan JavaScript dengan model pemrograman WebSocket tanpa soket yang mendengarkan. Deklarasikan stream di Go dengan `app.HandleStream(name, handler)` dan hubungkan dari frontend dengan `Stream(name)`, yang mengembalikan objek berbentuk `WebSocket`. Go→JS dibawa oleh satu polling tertahan per jendela melalui server aset, sedangkan JS→Go melalui POST biasa; tidak ada yang mengikat port TCP dan tidak ada yang melewati `evaluateJavaScript`. Dalam build server (`-tags server`), handler yang sama disajikan melalui WebSocket nyata, sehingga kode aplikasi identik di seluruh build. oleh @leaanthony
- Pindahkan entri log perubahan mailbox ke Belum Dirilis dalam [PR](https://github.com/wailsapp/wails/pull/5935) oleh @leaanthony

## Diubah

- Perbarui pembuatan otomatis bilah sisi dokumentasi dan derivasi tipe penulis blog dalam [PR](https://github.com/wailsapp/wails/pull/5938) oleh @leaanthony

## Diperbaiki

- Inisialisasi WebView2 menggunakan tenggat waktu dan pompa pesan dalam [PR](https://github.com/wailsapp/wails/pull/5952) oleh @leaanthony
- Pengujian cookie WebView2 dilewati di CI kecuali diaktifkan secara eksplisit dan eksekusinya dikunci ke thread OS saat ini dalam [PR](https://github.com/wailsapp/wails/pull/5951) oleh @leaanthony
- Builder menu Windows memulihkan ID perintah untuk item induk submenu dalam [PR](https://github.com/wailsapp/wails/pull/5944) oleh @gilad-ch
- Selaraskan image kompilasi silang resmi dengan baseline dukungan GTK 4.14+ Linux (#5928)
- Konfigurasikan proyek Xcode iOS agar mempertahankan flag linker yang diwariskan dan tambahkan -ObjC dalam [PR](https://github.com/wailsapp/wails/pull/5915) oleh @mortenolsrud
- Perbaiki pergantian koneksi TCP yang berlebihan pada proksi aset `wails3 dev` untuk frontend besar, yang dapat menghabiskan port sementara host dan menyebabkan proses yang tidak terkait gagal dengan `EADDRNOTAVAIL`
- Antrekan JavaScript peristiwa per jendela untuk pengiriman terurut dan backpressure dalam [PR](https://github.com/wailsapp/wails/pull/5934) oleh @leaanthony

## Dihapus

- Hapus pipeline rilis biner desktop: rilis v3 hanya berupa tag dan CLI `wails3` diinstal dengan `go install`. Menghapus `release-v3.yml` dan langkah nightly yang memicunya dalam [PR](https://github.com/wailsapp/wails/pull/5946) oleh @leaanthony

## v3.0.0-beta.7 - 2026-08-11

## Ditambahkan

- Tambahkan preferensi putar otomatis macOS untuk menonaktifkan persyaratan tindakan pengguna bagi pemutaran media dalam [PR](https://github.com/wailsapp/wails/pull/5512) oleh @Eyalm321
- Pindahkan entri log perubahan mailbox ke Belum Dirilis dalam [PR](https://github.com/wailsapp/wails/pull/5935) oleh @leaanthony

## Diubah

- Animasi zoom macOS menggunakan CADisplayLink atau NSTimer untuk kinerja yang lebih mulus dalam [PR](https://github.com/wailsapp/wails/pull/5945) oleh @savely-krasovsky

## Diperbaiki

- Konfigurasikan proyek Xcode iOS agar mempertahankan flag linker yang diwariskan dan tambahkan -ObjC dalam [PR](https://github.com/wailsapp/wails/pull/5915) oleh @mortenolsrud
- Perbaiki pergantian koneksi TCP yang berlebihan pada proksi aset `wails3 dev` untuk frontend besar, yang dapat menghabiskan port sementara host dan menyebabkan proses yang tidak terkait gagal dengan `EADDRNOTAVAIL`
- Antrekan JavaScript peristiwa per jendela untuk pengiriman terurut dan backpressure dalam [PR](https://github.com/wailsapp/wails/pull/5934) oleh @leaanthony

### Ditambahkan

- Implementasikan mailbox FIFO asinkron generik untuk pengiriman peristiwa secara terurut dalam [PR](https://github.com/wailsapp/wails/pull/5851) oleh @savely-krasovsky dan @DevLumuz

## v3.0.0-beta.6 - 2026-08-09

## Ditambahkan

- Implementasikan penyimpanan terbatas di sisi host untuk peristiwa berukuran terlalu besar dan pengiriman JavaScript secara terurut dalam [PR](https://github.com/wailsapp/wails/pull/5930) oleh @leaanthony
- Implementasikan efek pantulan Dock macOS untuk kedipan jendela dalam [PR](https://github.com/wailsapp/wails/pull/5921) oleh @julianstorer

## Diperbaiki

- Server aset mempertahankan galat pendeteksi jenis konten dan prefiks yang belum ditulis selama flushing dalam [PR](https://github.com/wailsapp/wails/pull/5931) oleh @leaanthony
- Cegah aplikasi macOS mengalami crash saat mengganti menu aplikasi dari callback Wails
- Perbaiki menu native yang tidak terbaca pada Windows 10 1809 / Windows Server 2019 (build 17763). Ekspor uxtheme mode gelap dibatasi untuk build 18334, sehingga pengaktifan mode gelap pada tingkat aplikasi tidak pernah dijalankan di host tersebut: latar belakang menu digambar gelap, tetapi Windows tetap menggambar teks menu menggunakan tema terang sehingga menghasilkan teks gelap di atas latar belakang gelap. Nomor ordinal tersebut tersedia sejak 17763, sehingga batasannya kini sesuai.
- Perbaiki `w32.GetStockObject` yang memanggil `GetDeviceCaps` alih-alih `GetStockObject`, yang menyebabkannya mengembalikan 0 untuk setiap objek bawaan.
- Tingkatkan penanganan dan pelaporan galat pengunduhan bootstrapper WebView2 dalam [PR](https://github.com/wailsapp/wails/pull/5924) oleh @jannskiee

## v3.0.0-beta.5 - 2026-08-07

## Diperbaiki

- Aktivasi aplikasi macOS kini mematuhi kebijakan aktivasi hanya untuk aplikasi reguler dalam [PR](https://github.com/wailsapp/wails/pull/5897) oleh @julianstorer
- Lindungi jendela GTK yang belum diinisialisasi dalam build Linux pada [PR](https://github.com/wailsapp/wails/pull/5898) oleh @julianstorer
- Tetapkan warna latar belakang opak secara eksplisit untuk jendela WebKit Linux sebelum URL dimuat dalam [PR](https://github.com/wailsapp/wails/pull/5899) oleh @julianstorer

## v3.0.0-beta.4 - 2026-08-05

## Diubah

- Tugas build Android menggunakan arm64 secara default dan deploy-emulator memilih arsitektur host dalam [PR](https://github.com/wailsapp/wails/pull/5890) oleh @mortenolsrud

## Diperbaiki

- Pertahankan status zoom jendela macOS selama penyeretan dan kurangi gerakan dalam [PR](https://github.com/wailsapp/wails/pull/5900) oleh @leaanthony
- Perbaiki build mode server Windows dengan menambahkan `!server` ke batasan build `webview_window_windows_nonclient.go`

## v3.0.0-beta.3 - 2026-08-03

## Ditambahkan

- Dokumentasikan penyelesaian verifikasi beta Fase 10 dalam detail implementasi pada [PR](https://github.com/wailsapp/wails/pull/5881) oleh @leaanthony

## Diperbaiki

- Teruskan handle jendela ke API mode gelap Windows dan validasi argumen dalam [PR](https://github.com/wailsapp/wails/pull/5877) oleh @leaanthony
- Pusatkan resolusi status tombol bilah judul macOS untuk jendela tanpa bingkai dalam [PR](https://github.com/wailsapp/wails/pull/5870) oleh @taliesin-ai
- Cegah teks menu native tidak terbaca ketika aplikasi Windows meminta mode gelap sementara tema aplikasi Windows terang. Menu kini menggunakan latar belakang native terang yang sesuai hingga Windows dapat merender teks menu gelap.
- Perbaiki menu native yang tidak terbaca pada Windows 10 1809 / Windows Server 2019 (build 17763). Ekspor uxtheme mode gelap dibatasi untuk build 18334, sehingga pengaktifan mode gelap pada tingkat aplikasi tidak pernah dijalankan di host tersebut: latar belakang menu digambar gelap, tetapi Windows tetap menggambar teks menu menggunakan tema terang sehingga menghasilkan teks gelap di atas latar belakang gelap. Nomor ordinal tersebut tersedia sejak 17763, sehingga batasannya kini sesuai.

## v3.0.0-beta.2 - 2026-08-02

## Diubah

- Promosikan v3 dari alfa ke beta
- Dokumentasikan default cerdas baki sistem dan perilaku sembunyi otomatis popup, beserta cakupan regresi untuk pemilihan handler klik (#5840).
- Updater GitHub mengecualikan aset penginstal Windows secara default dalam [PR](https://github.com/wailsapp/wails/pull/5861) oleh @leaanthony
- Tambahkan dukungan untuk sudut membulat, siku, dan berjari-jari khusus pada jendela macOS tanpa bingkai dalam [PR](https://github.com/wailsapp/wails/pull/5866) oleh @leaanthony

## Diperbaiki

- Laporkan ukuran jendela GTK4 secara langsung dan pancarkan peristiwa perubahan ukuran, pemaksimalan, peminimalan, serta status layar penuh dari permukaan yang dikonfigurasi (#5830).
- Perbaiki crash WebKit di Linux saat mengirim Blob atau FormData dalam permintaan fetch pada [PR](https://github.com/wailsapp/wails/pull/5854) oleh @taliesin-ai
- Shim fetch meneruskan undefined untuk header Blob/FormData yang tidak tersedia dalam [PR](https://github.com/wailsapp/wails/pull/5865) oleh @leaanthony

## v3.0.0-alpha2.122 - 2026-08-01

## Ditambahkan

## Diubah

- Tambahkan dukungan untuk sudut membulat, siku, dan berjari-jari khusus pada jendela macOS tanpa bingkai dalam [PR](https://github.com/wailsapp/wails/pull/5866) oleh @leaanthony

## Diperbaiki

- Shim fetch meneruskan undefined untuk header Blob/FormData yang tidak tersedia dalam [PR](https://github.com/wailsapp/wails/pull/5865) oleh @leaanthony

## v3.0.0-alpha2.121 - 2026-07-31

## Ditambahkan

- Tambahkan dukungan pemaketan DMG macOS dengan opsi dan tugas build baru dalam [PR](https://github.com/wailsapp/wails/pull/5857) oleh @leaanthony

## Diubah

- Pembaruan GitHub secara default mengecualikan aset penginstal Windows dalam [PR](https://github.com/wailsapp/wails/pull/5861) oleh @leaanthony

## Diperbaiki

- Perbaiki crash WebKit di Linux saat mengirim Blob atau FormData dalam permintaan fetch pada [PR](https://github.com/wailsapp/wails/pull/5854) oleh @taliesin-ai

## v3.0.0-alpha2.120 - 2026-07-31

## Ditambahkan

- Implementasikan tindakan klik ganda pada bilah judul macOS untuk memaksimalkan atau meminimalkan jendela dalam [PR](https://github.com/wailsapp/wails/pull/5853) oleh @taliesin-ai

## Diubah

- Perbarui tutorial layanan QR agar menggunakan NewServiceWithOptions dan tambahkan spasi dalam [PR](https://github.com/wailsapp/wails/pull/5849) oleh @jeongkyu

## Diperbaiki

- Pertahankan responsivitas WKWebView selama zoom di macOS dalam [PR](https://github.com/wailsapp/wails/pull/5856) oleh @leaanthony
- Perbaiki kueri ukuran jendela GTK4 dan pancarkan peristiwa perubahan ukuran, pemaksimalan, peminimalan, serta status layar penuh dari `GdkSurface` yang dikonfigurasi.

## v3.0.0-alpha2.119 - 2026-07-27

## Diperbaiki

- Perbarui dokumentasi dalam beberapa bahasa untuk menyertakan diagram arsitektur dalam [PR](https://github.com/wailsapp/wails/pull/5833) oleh @taliesin-ai

## v3.0.0-alpha2.118 - 2026-07-26

## Ditambahkan

- Sediakan jalur default untuk masukan dan keluaran pembuatan ikon dalam [PR](https://github.com/wailsapp/wails/pull/5825) oleh @taliesin-ai
- Tambahkan modul entri sumber ke sideEffects dalam package.json runtime pada [PR](https://github.com/wailsapp/wails/pull/5797) oleh @savely-krasovsky
- Tambahkan bagian lisensi dan asal-usul ke panduan kontribusi dalam [PR](https://github.com/wailsapp/wails/pull/5816) oleh @taliesin-ai

## Diperbaiki

- Terapkan CSS tanpa bingkai GTK4 dengan cakupan terbatas untuk menghapus radius batas dalam [PR](https://github.com/wailsapp/wails/pull/5800) oleh @savely-krasovsky
- Tangani kegagalan memperoleh posisi kursor dengan baik di Windows untuk menu pop-up dan enumerasi layar dalam [PR](https://github.com/wailsapp/wails/pull/5789) oleh @wayneforrest
- Dialog buka berkas macOS memfilter ekstensi dengan benar dan memvalidasi berkas yang diizinkan berdasarkan sufiks dalam [PR](https://github.com/wailsapp/wails/pull/5678) oleh @phergul
- Lindungi inisialisasi mode gelap Windows dari pemanggilan API nil dalam [PR](https://github.com/wailsapp/wails/pull/5793) oleh @roachadam
- Perbaiki kegagalan build 32-bit pada pembaru: konstanta `maxArchiveTotalSize` (2 GiB) melampaui kapasitas `int` platform saat diteruskan ke `fmt.Errorf` pada `GOARCH=386`. Kini konstanta tersebut secara eksplisit bertipe `int64`.
- Perbaiki panic akibat pointer nil saat startup ketika jendela menggunakan bilah judul Gelap (atau gelap sistem) pada build Windows yang tidak memuat API uxtheme mode gelap, seperti Windows 10 1809 / Windows Server 2019 (build 17763). Pemanggilan `AllowDarkModeForWindow` dalam penyiapan tema jendela kini dilindungi dari nil, sesuai dengan perlindungan yang sudah digunakan dalam `w32.SetMenuTheme`.

## v3.0.0-alpha2.117 - 2026-07-08

## Ditambahkan

- Implementasikan logika hit-test khusus untuk wilayah nonklien di Windows dalam [PR](https://github.com/wailsapp/wails/pull/5462) oleh @savely-krasovsky

## Diubah

- Konfigurasikan deteksi skala monitor WebView2 berdasarkan UseVisualHosting dalam [PR](https://github.com/wailsapp/wails/pull/5761) oleh @wayneforrest

## v3.0.0-alpha2.116 - 2026-07-07

## Ditambahkan

- Perbarui dokumentasi FAQ agar berfokus pada fitur dan panduan Wails v3 dalam [PR](https://github.com/wailsapp/wails/pull/5763) oleh @taliesin-ai

## v3.0.0-alpha2.115 - 2026-07-06

## Diperbaiki

- Perbaiki `Menu.Update()` yang tidak membangun ulang menu native di GTK4 Linux (#5659, didiagnosis dan diperbaiki secara independen oleh @puneetdixit200 dalam #5539)
- Perbaiki crash saat mengenumerasi layar macOS ketika tampilan berubah dengan menyalin string id/nama layar dan mengambil snapshot jumlah layar (#5565, didiagnosis dan diperbaiki secara independen oleh @x-haose dalam #5584)
- Perbaiki crash di Windows ketika `WM_ERASEBKGND` menggambar latar belakang warna solid selama transisi minimalkan/pulihkan saat `GetClientRect` mengembalikan nil (perlindungan dilaporkan oleh @sinspired dalam #5636)
- Perbaiki kesalahan binding frontend yang selalu diurai sebagai teks oleh @mbaklor dalam #5690
- Perbaiki kegagalan build Windows saat menggunakan tag build `server`, yang disebabkan oleh tidak adanya batasan build `!server` pada berkas GUI Windows meskipun batasan tersebut sudah ada pada berkas ekuivalen macOS dan Linux (#5680)

## v3.0.0-alpha2.114 - 2026-07-05

## Ditambahkan

- Implementasikan protokol Update Manifest dan penyedia endpoint dalam [PR](https://github.com/wailsapp/wails/pull/5720) oleh @taliesin-ai

## Diubah

- Gabungkan binding `webview2` ke dalam modul v3 sebagai `v3/internal/webview2`, serta hapus modul mandiri, alur kerja rilis/sinkronisasi nightly, dan mekanisme perubahan versi go.mod-nya (v3 adalah satu-satunya konsumennya) dalam [PR](https://github.com/wailsapp/wails/pull/5711) oleh @taliesin-ai

## Diperbaiki

- Pindahkan deteksi skala monitor WebView2 dan perbaikan sinkronisasi ulang host saat DPI berubah ke bagian Belum Dirilis dalam [PR](https://github.com/wailsapp/wails/pull/5750) oleh @taliesin-ai
- Perbarui marshaling COM WebView2 untuk parameter float64 dan BOOL dalam [PR](https://github.com/wailsapp/wails/pull/5741) oleh @wayneforrest
- Cegah panic dan dereferensi nil saat memperbarui serta menghancurkan ikon baki sistem Windows dalam [PR](https://github.com/wailsapp/wails/pull/5703) oleh @wayneforrest
- Perbaiki jendela tersembunyi yang tidak disembunyikan kembali dengan benar di Windows dalam [PR](https://github.com/wailsapp/wails/pull/5743) oleh @wayneforrest
- Sinkronkan visibilitas pengontrol WebView2 dengan tindakan meminimalkan, memaksimalkan, dan memulihkan jendela dalam [PR](https://github.com/wailsapp/wails/pull/5742) oleh @wayneforrest

### Diperbaiki

- Aktifkan kembali deteksi skala monitor WebView2 dan jalankan sinkronisasi ulang host hanya saat DPI berubah dalam [PR](https://github.com/wailsapp/wails/pull/5734) oleh @taliesin-ai, berdasarkan perbaikan yang divalidasi oleh @randalmurphal, dengan verifikasi akar penyebab oleh @eleclin dan pengujian perangkat keras oleh @qq540491950

## v3.0.0-alpha2.113 - 2026-07-04

## Ditambahkan

- Tambahkan peringatan saat membangun AAB rilis tanpa menetapkan `ANDROID_KEYSTORE_FILE` (Google Play menolak bundle yang ditandatangani dengan kunci debug) serta dokumentasikan pengemasan dan penandatanganan App Bundle dalam [PR](https://github.com/wailsapp/wails/pull/5730) oleh @taliesin-ai
- Tambahkan dokumentasi Why Wails dalam beberapa bahasa dalam [PR](https://github.com/wailsapp/wails/pull/5739) oleh @taliesin-ai
- Dukung pemetaan Go time.Time ke JS Date atau string dalam binding pada [PR](https://github.com/wailsapp/wails/pull/5398) oleh @fbbdev
- Tambahkan tugas pengemasan Android App Bundle (AAB) (`bundle`, `bundle:fat`, `assemble:aab`, `assemble:aab:release`) untuk pengiriman ke Play Store—tugas APK tetap tersedia untuk pengujian lokal/emulator—dalam [PR](https://github.com/wailsapp/wails/pull/5728) oleh @mortenolsrud (memperbaiki [#5726](https://github.com/wailsapp/wails/issues/5726))
- Tambahkan target tugas untuk perangkat Android fisik dan lanjutkan kembali izin kamera/lokasi dalam [PR](https://github.com/wailsapp/wails/pull/5735) oleh @taliesin-ai

## Diubah

- Tingkatkan `webview2` ke v1.0.28 ([catatan rilis](https://github.com/wailsapp/wails/releases/tag/webview2%2Fv1.0.28)).
- Tingkatkan `compileSdk`/`targetSdk` pada templat Android dari 34 ke 35, yang diwajibkan Google Play untuk pengiriman aplikasi baru, dalam [PR](https://github.com/wailsapp/wails/pull/5730) oleh @taliesin-ai

## Diperbaiki

- Perbaiki mask avatar bawaan dalam sponsorkit pada [PR](https://github.com/wailsapp/wails/pull/5745) oleh @leaanthony
- Perbaiki pembuatan otomatis AVD Android yang memilih system image atau versi cmdline-tools yang salah akibat pengurutan versi secara leksikografis dalam [PR](https://github.com/wailsapp/wails/pull/5730) oleh @taliesin-ai
- Perbaiki wizard penyiapan yang menyarankan versi Android NDK usang (sekarang 26.3.11579264, sesuai dengan persyaratan yang didokumentasikan) dalam [PR](https://github.com/wailsapp/wails/pull/5730) oleh @taliesin-ai
- Perbarui dokumentasi bahasa Prancis untuk SvelteKit dan opsi dalam [PR](https://github.com/wailsapp/wails/pull/5744) oleh @leaanthony
- Perbaiki SIGSEGV saat mengenumerasi layar macOS ketika tampilan berubah dalam [PR](https://github.com/wailsapp/wails/pull/5516) oleh @flofreud

## v3.0.0-alpha2.112 - 2026-07-03

## Ditambahkan

- Tambahkan generator SVG kontributor berbasis Go dan perbarui halaman kredit dokumentasi/situs web dalam [PR](https://github.com/wailsapp/wails/pull/5724) oleh @taliesin-ai
- Dukung pemetaan Go time.Time ke JS Date atau string dalam binding pada [PR](https://github.com/wailsapp/wails/pull/5398) oleh @fbbdev

## Diubah

- Ganti pipeline gambar sponsor berbasis Node dengan generator Go dalam [PR](https://github.com/wailsapp/wails/pull/5719) oleh @taliesin-ai

## Diperbaiki

- Perbaiki skrip instalasi dependensi aset build Android dalam [PR](https://github.com/wailsapp/wails/pull/5729) oleh @taliesin-ai
- Tolak karakter kontrol U+0085 (NEXT LINE) dalam `ValidateAndSanitizeURL`, sehingga cakupan karakter spasi pada validator URL menjadi lengkap
- Hitung ulang frame DWM saat DPI berubah untuk jendela tanpa frame dalam [PR](https://github.com/wailsapp/wails/pull/4785) oleh @leaanthony
- Perbaiki kegagalan deteksi zona pelepasan DnD pada penskalaan selain 100% di Windows dalam [PR](https://github.com/wailsapp/wails/pull/4632) oleh @yulesxoxo
- Tambahkan pengelolaan memori Objective-C secara eksplisit untuk objek Cocoa pada dialog, menu, baki sistem, dan notifikasi Darwin dalam [PR](https://github.com/wailsapp/wails/pull/5714) oleh @taliesin-ai
- Perbaiki bug backend CGO Linux dan masalah baki sistem dalam [PR](https://github.com/wailsapp/wails/pull/5718) oleh @taliesin-ai

## v3.0.0-alpha2.111 - 2026-07-01

## Ditambahkan

- Tambahkan HappyTools ke etalase komunitas dalam [PR](https://github.com/wailsapp/wails/pull/5061) oleh @Aliuyanfeng
- Tambahkan dukungan locale bahasa Indonesia dan dokumentasi lengkap dalam [PR](https://github.com/wailsapp/wails/pull/5643) oleh @triadmoko
- Tambahkan opsi DisableMenu ke WindowsWindow dalam [PR](https://github.com/wailsapp/wails/pull/4813) oleh @leaanthony

## Diubah

- Perbarui templat Taskfile dan CLI agar mendistribusikan tugas build/package dengan GOOS dan ARCH dalam [PR](https://github.com/wailsapp/wails/pull/5617) oleh @leaanthony

## Diperbaiki

- Perbaiki masalah pengelompokan jendela ke dalam tab di macOS dalam [PR](https://github.com/wailsapp/wails/pull/5708) oleh @taliesin-ai

## Dihapus

- Hapus file MDX terjemahan bahasa Jerman dari bagian kontribusi, fitur, dan panduan dalam [PR](https://github.com/wailsapp/wails/pull/5702) oleh @taliesin-ai

## v3.0.0-alpha2.110 - 2026-06-30

## Ditambahkan

- Implementasikan pemuatan ulang dan pemuatan ulang paksa WebView di macOS, serta tambahkan pemulihan saat proses WebContent dihentikan dalam [PR](https://github.com/wailsapp/wails/pull/5129) oleh @wayneforrest
- Tambahkan dokumentasi lengkap berbahasa Jerman untuk bagian kontribusi, fitur, dan panduan dalam [PR](https://github.com/wailsapp/wails/pull/5396) oleh @leaanthony
- Tingkatkan notifikasi dengan suara, lampiran, penjadwalan, dan API pembaruan dalam [PR](https://github.com/wailsapp/wails/pull/5333) oleh @popaprozac

## Diperbaiki

- Hitung ulang frame DWM saat DPI berubah untuk jendela tanpa frame dalam [PR](https://github.com/wailsapp/wails/pull/4785) oleh @leaanthony
- Perbaiki kegagalan deteksi zona pelepasan DnD pada penskalaan selain 100% di Windows dalam [PR](https://github.com/wailsapp/wails/pull/4632) oleh @yulesxoxo

## v3.0.0-alpha2.109 - 2026-06-29

## Ditambahkan

- Tambahkan contoh kode ke dokumentasi EventsEmit dalam [PR](https://github.com/wailsapp/wails/pull/5026) oleh @iamhabbeboy
- Tambahkan opsi visual hosting WebView2 untuk Windows dalam [PR](https://github.com/wailsapp/wails/pull/5380) oleh @MerIijn
- Tambahkan Klustr ke dokumentasi etalase komunitas dalam [PR](https://github.com/wailsapp/wails/pull/5536) oleh @SametKUM
- Tambahkan Kira ke etalase komunitas beserta halaman baru dan entri changelog dalam [PR](https://github.com/wailsapp/wails/pull/5685) oleh @thiennguyen93
- Tambahkan bagian umpan balik ke panduan layanan MCP dalam [PR](https://github.com/wailsapp/wails/pull/5694) oleh @taliesin-ai

## Diubah

- Mode server kini memiliki build produksi sebagai fitur utama, konsisten dengan tugas build desktop (#5693). Secara default, `task build:server` membangun binary produksi (`-tags server,production`, `-trimpath`, dipangkas), serta menerima `DEV=true` (server pengembangan), `OBFUSCATED=true` (garble), dan `EXTRA_TAGS`. `task run:server` menjalankan server pengembangan. `Dockerfile.server` / `task build:docker` membangun frontend produksi terlebih dahulu, lalu server produksi (`-tags server,production`); secara default, image menggunakan build statis Go murni pada distroless/static, dengan `CGO_ENABLED`, `GO_IMAGE`, dan `RUNTIME_IMAGE` yang diekspos sebagai argumen build yang dapat ditimpa untuk aplikasi CGO.

## Diperbaiki

- Cegah crash saat menutup jendela yang masih memiliki panggilan asinkron tertunda dalam [PR](https://github.com/wailsapp/wails/pull/4435) oleh @leaanthony
- Cegah aktivasi jendela saat membuka aplikasi tersembunyi di Windows dalam [PR](https://github.com/wailsapp/wails/pull/5249) oleh @leaanthony
- Pastikan penanganan metadata permintaan WebKit, penyelesaian respons, dan stream body berjalan di thread utama GTK dalam [PR](https://github.com/wailsapp/wails/pull/5668) oleh @taliesin-ai
- Perbaiki `Menu.Update()` yang tidak membangun ulang menu native di GTK4 Linux (#5659, didiagnosis dan diperbaiki secara independen oleh @puneetdixit200 di #5539)
- Perbaiki crash saat mengiterasi layar macOS ketika tampilan berubah dengan menyalin string ID/nama layar dan mengambil snapshot jumlahnya (#5565, didiagnosis dan diperbaiki secara independen oleh @x-haose di #5584)
- Perbaiki konten WebView2 yang menyusut lalu menghilang setelah jendela diseret melintasi monitor dengan DPI berbeda di Windows dengan menetapkan ulang batas pengontrol dalam handler `WM_DPICHANGED`, mengikuti sinkronisasi ulang DPI saat keluar dari keadaan diminimalkan (#5677)

## v3.0.0-alpha2.108 - 2026-06-28

## Ditambahkan

- Tambahkan pintasan papan ketik global (seluruh sistem) melalui `app.GlobalShortcut` (`Register`, `Unregister`, `UnregisterAll`, `IsRegistered`, `GetAll`). Pintasan tetap dipicu meskipun aplikasi tidak memiliki fokus. Diimplementasikan secara native untuk setiap platform tanpa dependensi pihak ketiga: hot key Carbon di macOS, `RegisterHotKey` di Windows, `XGrabKey` di X11, dan antarmuka pintasan global XDG Desktop Portal di Wayland.
- Tambahkan server MCP bawaan: server Model Context Protocol yang dimulai secara otomatis ketika aplikasi dibangun dengan tag `mcp`, sehingga agen LLM dapat menguji dan mengendalikan aplikasi Wails yang sedang berjalan—termasuk mengendalikan jendela, memeriksa DOM, mengevaluasi JavaScript, memanggil metode terikat, menangani peristiwa, serta menyimulasikan masukan tetikus/papan ketik yang ditampilkan dengan kursor animasi di layar. Tidak memerlukan kode pengguna: tag `mcp` ditambahkan secara otomatis oleh `wails3 build`/`wails3 dev` ketika `WAILS_MCP=1` ditetapkan. Seluruh konfigurasi dilakukan melalui variabel lingkungan (`WAILS_MCP_HOST`, `WAILS_MCP_PORT`, `WAILS_MCP_TIMEOUT`, `WAILS_MCP_HIDE_CURSOR`).

## Diperbaiki

- Perbaiki `Menu.Update()` yang tidak membangun ulang menu native di GTK4 Linux (#5659, didiagnosis dan diperbaiki secara independen oleh @puneetdixit200 di #5539)
- Perbaiki crash saat mengiterasi layar macOS ketika tampilan berubah dengan menyalin string ID/nama layar dan mengambil snapshot jumlahnya (#5565, didiagnosis dan diperbaiki secara independen oleh @x-haose di #5584)

## v3.0.0-alpha2.107 - 2026-06-27

## Ditambahkan

- Tambahkan dokumentasi eksperimental Wake dengan navigasi bilah sisi dalam [PR](https://github.com/wailsapp/wails/pull/5613) oleh @leaanthony

## v3.0.0-alpha2.106 - 2026-06-24

## Diubah

- Tingkatkan `webview2` ke v1.0.27.
  - ci(webview2): perbaiki build rilis (kompilasi silang Windows + go.sum lengkap) (#5671)\

  **Diff lengkap:** https://github.com/wailsapp/wails/compare/webview2/v1.0.26...webview2/v1.0.27

- Hapus go vet dari kompilasi silang alur kerja rilis webview2 dalam [PR](https://github.com/wailsapp/wails/pull/5672) oleh @taliesin-ai
- Perbarui model OpenRouter untuk changelog otomatis ke google/gemini-2.5-flash-lite dalam [PR](https://github.com/wailsapp/wails/pull/5670) oleh @taliesin-ai
- Tingkatkan `webview2` ke v1.0.26.

### Perbaikan

- **Pulihkan aplikasi dari error COM runtime sementara alih-alih keluar** (#5658, #5580). Sebelumnya, `Chromium.errorCallback` memanggil `os.Exit(1)` untuk *setiap* error COM, sehingga gangguan sementara yang dapat dipulihkan setelah aplikasi dimulai akan menghentikan seluruh aplikasi. Jalur runtime (`Resize`/`GetClientRect`, `Navigate`/`NavigateToString`, `Init`, `MessageReceived`, `PutZoomFactor`, `OpenDevToolsWindow`) kini mencatat error dan melakukan pemulihan. Secara khusus, pesan web yang rusak/tidak tepercaya dalam `MessageReceived` kini dibuang dan tidak lagi menghentikan proses. Ini menangani kategori crash saat melintasi monitor dengan DPI berbeda (#5544, #5650). Jalur pembuatan lingkungan/pengontrol tetap bersifat fatal.\

**Diff lengkap:** https://github.com/wailsapp/wails/compare/webview2/v1.0.25...webview2/v1.0.26

## Diperbaiki

- Perbaiki alur kerja release-webview2 agar menangani file go.sum dengan benar dalam [PR](https://github.com/wailsapp/wails/pull/5671) oleh @taliesin-ai
- Perbaiki pembaruan menu GTK4 Linux dengan mengosongkan dan membangun ulang menu native dalam [PR](https://github.com/wailsapp/wails/pull/5659) oleh @taliesin-ai

## v3.0.0-alpha2.105 - 2026-06-21

## Ditambahkan

- Tambahkan `application.System` untuk mendeteksi platform saat runtime dari kode bersama: `System.IsMobile()` (iOS/Android), `System.IsDesktop()` (macOS/Windows/Linux), `System.IsServer()` (tag build `server`), dan `System.IsPlatform(application.PlatformMacOS|PlatformWindows|PlatformLinux|PlatformIOS|PlatformAndroid|PlatformServer)` untuk menguji satu target secara langsung. Kode ini dapat dikompilasi di setiap target, sehingga Anda dapat membuat percabangan tanpa tag build. Fungsi pembantu frontend yang sesuai (`System.IsMobile/IsDesktop/IsIOS/IsAndroid/...`) tersedia di `@wailsio/runtime`
- Tambahkan panduan "Menggunakan Framework Frontend Lain" yang menunjukkan cara memasukkan proyek Vite Anda sendiri ke `frontend/` (mencakup Solid, Preact, Lit, SvelteKit, Qwik, Angular, dan lainnya)
- Wizard `wails3 setup` kini memeriksa toolchain seluler (iOS/Android)—Xcode dan runtime iOS Simulator, JDK, Android SDK/NDK, serta emulator—dengan instalasi sekali klik dan perbaikan konfigurasi shell yang dapat disalin jika tersedia
- Proyek yang dihasilkan menyertakan pengaturan `frontend/.npmrc` dengan `minimum-release-age` selama 7 hari untuk mengurangi paparan terhadap paket yang baru dipublikasikan (dan mungkin telah disusupi) (dipatuhi oleh pnpm dan bun; diabaikan tanpa dampak oleh npm)

## Diubah

- Desain ulang semua templat awal bawaan dengan tampilan hero pegunungan neon baru (web, iOS, dan Android)
- **TypeScript kini menjadi bahasa default untuk templat awal dan menggunakan nama templat dasar.** `wails3 init` (tanpa `-t`) membuat struktur awal proyek TypeScript; `-t vanilla`, `-t react`, `-t vue`, dan `-t svelte` menggunakan TypeScript, sedangkan varian JavaScript tersedia di `-t vanilla-js`, `-t react-js`, `-t vue-js`, dan `-t svelte-js`. Templat bawaan mendeklarasikan bahasanya dengan `typescript:` dalam `template.yaml`; templat komunitas yang menggunakan sufiks `-ts` tetap berfungsi sebagai fallback
- Desain ulang wizard `wails3 setup` dengan tema neon "digital Wails" (efek kaca buram nan semarak di atas latar pegunungan)

## Diperbaiki

- Perbaiki crash di Windows saat memulihkan aplikasi yang diminimalkan cukup lama hingga WebView2 ditangguhkan atau proses rendering/GPU-nya didaur ulang. Sinkronisasi ulang DPI saat meminimalkan/memulihkan (#5544) kini hanya mengakses pengontrol WebView2 jika DPI jendela benar-benar berubah, sehingga menghindari pemanggilan COM fatal pada pengontrol yang ditangguhkan dalam pemulihan umum dengan DPI yang sama (#5605)
- Perbaiki crash native `SIGABRT`/`SIGSEGV` yang berulang (biasanya di dalam `g_object_unref` selama perulangan utama GTK) pada aplikasi Linux yang berjalan lama dan sering memuat aset/media. Server aset menyelesaikan `WebKitURISchemeRequest` dari goroutine pekerja, sehingga fungsi WebKit2GTK yang tidak aman terhadap thread dipanggil di luar thread utama GTK; penyelesaian (`webkit_uri_scheme_request_finish_with_response`/`finish_error`) kini dijalankan di thread utama. Menyempurnakan perbaikan parsial dalam #5566. Memengaruhi build 6.0 GTK3 maupun GTK4/WebKitGTK (#5631, #5557)
- Perbaiki `fatal error: invalid pointer found on stack` yang terjadi sesekali dalam `setupSignalHandlers` di Linux/GTK3. ID jendela yang diteruskan sebagai `user_data` sinyal disimpan dalam variabel lokal Go bertipe `unsafe.Pointer`, sehingga garbage collector membatalkan proses ketika memindai nilai (non-pointer) tersebut selama penyalinan stack. ID tersebut kini dipertahankan sebagai tipe integer (`uintptr_t`) di sisi Go, dengan melakukan backport ke jalur GTK3 lama atas perbaikan yang sama yang diterapkan #4958 pada jalur GTK4 (yang mengubah fungsi sinyal C menjadi `uintptr_t` untuk mengatasi error `-race`/checkptr) (#5631)

## Dihapus

- Hapus templat awal `react-swc`, `preact`, `lit`, `solid`, `qwik`, dan `sveltekit` (beserta varian `-ts`-nya). Kumpulan bawaan yang didukung kini terdiri dari `vanilla`, `react`, `vue`, dan `svelte`—masing-masing menggunakan TypeScript secara default, dengan varian JavaScript `-js`. Framework lain apa pun tetap dapat digunakan dengan [membawa frontend Anda sendiri](https://v3.wails.io/guides/dev/frontend-frameworks) atau melalui templat khusus

## v3.0.0-alpha2.104 - 2026-06-18

## Diperbaiki

- Perbaiki crash iOS (SIGABRT) saat metode layanan Go yang terikat mengembalikan string kosong. Penulis respons aset iOS memeriksa pointer isi dengan `buf != nil` alih-alih panjangnya, sehingga isi dengan panjang nol membuat `&buf[0]` mengalami panic; sekarang pemeriksaan dilakukan berdasarkan panjang, selaras dengan penulis respons desktop

## v3.0.0-alpha2.103 - 2026-06-15

## Diubah

- Pindahkan fitur native iOS dan Android ke pengelola platform: panggil melalui `application.IOS.*` dan `application.Android.*` (misalnya `application.IOS.Haptic("medium")`, `application.Android.Share(payload)`), bukan melalui fungsi bebas lama `application.IOS*`/`application.Android*` (#5602)
- Ganti nama event bridge seluler: event lintas platform kini menggunakan prefiks `common:*` (misalnya `common:haptic`, `common:location`), sedangkan event khusus platform menggunakan `ios:*` / `android:*` (misalnya `ios:backgroundTask`, `android:foregroundService`); prefiks `native:*` tidak lagi digunakan (#5602)

## v3.0.0-alpha.102 - 2026-06-14

## Ditambahkan

- Tambahkan wizard eksperimental `wails3 setup` untuk penyiapan proyek interaktif dan pemeriksaan dependensi
- Tambahkan flag `--json` ke `wails3 doctor` untuk keluaran yang dapat dibaca mesin
- Tambahkan bagian status penandatanganan ke perintah `wails3 doctor`

## Diperbaiki

- Perbaiki deteksi npm di Linux agar memeriksa PATH selain pengelola paket

## v3.0.0-alpha.101 - 2026-06-13

## Ditambahkan

- iOS: dialog pesan native (UIAlertController) serta dialog untuk membuka file, beberapa file, dan direktori (UIDocumentPickerViewController); dialog penyimpanan mengembalikan error secara eksplisit
- iOS: dukungan papan klip melalui UIPasteboard
- iOS: metrik layar aktual melalui UIScreen (point, piksel, skala, area kerja dalam safe area)
- iOS: build untuk perangkat (`IOS_PLATFORM=device`), dukungan identitas penandatanganan kode / profil provisioning / entitlement, pemaketan `.ipa`, dan `deploy-device` melalui devicectl
- iOS: versi minimum iOS yang dapat dikonfigurasi (`ios.minIOSVersion` dalam build/config.yml)
- iOS: `wails3 doctor` melaporkan ketersediaan Xcode dan iOS SDK di macOS
- iOS: event sistem — baterai, jaringan, tema, penguncian layar, dan memori rendah diekspos sebagai event aplikasi `events.IOS.*` dan `events.Common.*` yang netral terhadap platform
- iOS: bridge fitur seluler native (`application.IOS*` yang diekspor) — lembar berbagi, membuka URL, menjaga perangkat tetap aktif, senter, inset safe area, kecerahan, info aplikasi, penguncian orientasi, bilah status, biometrik (Face ID/Touch ID), notifikasi lokal, dan penyimpanan aman Keychain
- iOS: sensor & perangkat keras — haptik, geolokasi sekali ambil, akselerometer, proksimitas, text-to-speech, info penyimpanan, status daya/baterai, status jaringan, inset papan ketik, dan deteksi tangkapan layar
- iOS: dokumentasi (IOS.md dan panduan di situs dokumentasi)
- Android: dialog pesan native (AlertDialog) serta dialog untuk membuka file dan beberapa file (Storage Access Framework, diimpor sebagai salinan cache); dialog untuk membuka direktori dan menyimpan mengembalikan error secara eksplisit
- Android: dukungan papan klip melalui ClipboardManager
- Android: metrik layar aktual melalui WindowMetrics/DisplayMetrics (dp, piksel, skala, area kerja bilah sistem)
- Android: metode runtime untuk haptik (`Android.Haptics.Vibrate`), info perangkat (`Android.Device.Info`), dan toast (`Android.Toast.Show`)
- Android: event siklus hidup bertipe (`events.Android.*`, dihasilkan dari events.txt), dengan `ActivityCreated` dipetakan ke `Common.ApplicationStarted`
- Android: pipeline build menghasilkan APK debug dan rilis yang dapat diinstal (`android:run`, `android:package`, `android:package:fat`); penandatanganan rilis secara default menggunakan debug keystore atau keystore sebenarnya melalui variabel lingkungan `ANDROID_KEYSTORE_*`
- Android: `wails3 doctor` melaporkan Android SDK, NDK, dan JDK
- Android: event sistem — baterai, jaringan, tema, penguncian layar, dan memori rendah diekspos sebagai event aplikasi `events.Android.*` dan `events.Common.*` yang netral terhadap platform
- Android: bridge fitur seluler native (`application.Android*` yang diekspor) — berbagi, membuka URL, menjaga perangkat tetap aktif, senter, inset safe area, kecerahan, info aplikasi, penguncian orientasi, bilah status, biometrik (BiometricPrompt), notifikasi lokal, dan penyimpanan aman EncryptedSharedPreferences
- Android: sensor & perangkat keras — haptik, geolokasi sekali ambil, akselerometer, proksimitas, text-to-speech, info penyimpanan, status daya/baterai, status jaringan, inset papan ketik, dan pemblokiran tangkapan layar dengan FLAG_SECURE
- Android: dokumentasi (ANDROID.md dan panduan di situs dokumentasi)
- Contoh: kitchen sink `mobile` memperoleh tab Mobile dan Hardware yang mendemonstrasikan bridge fitur native di iOS dan Android (tab berbentuk pil membungkus menjadi beberapa baris)
- Seluler: baterai — akselerometer, sensor proksimitas, senter, dan jam berkala dalam contoh dijeda saat aplikasi berada di latar belakang dan dipulihkan ketika aplikasi kembali aktif (Android mempertahankan proses tetap berjalan di latar belakang, sedangkan senter merupakan status perangkat keras yang tetap bertahan di iOS), dan penerima event sistem Android hanya didaftarkan saat aplikasi berada di latar depan
- iOS: pengambilan gambar kamera — `application.IOSCapturePhoto`/`IOSCaptureVideo` (UIImagePickerController → event `native:capture` dengan thumbnail base64)
- iOS: eksekusi di latar belakang — `application.IOSBeginBackgroundTask`/`IOSEndBackgroundTask` (jendela tugas latar belakang UIApplication) dan `ios.backgroundModes` yang dapat dikonfigurasi (build/config.yml), yang menerapkan `UIBackgroundModes` ke Info.plist yang dihasilkan melalui templat
- Android: pengambilan gambar kamera — `application.AndroidCapturePhoto`/`AndroidCaptureVideo` (kamera sistem melalui FileProvider → event `native:capture`)
- Android: layanan latar depan — `application.AndroidStartForegroundService`/`AndroidStopForegroundService` (`WailsForegroundService` dengan notifikasi berkelanjutan menjaga proses tetap aktif untuk pekerjaan latar belakang yang berjalan lama)
- Contoh: tab Camera yang mendemonstrasikan pengambilan foto/video dan eksekusi di latar belakang (layanan latar depan di Android, jendela tugas latar belakang di iOS)

## Diperbaiki

- Perbaiki `getUserMedia` yang selalu gagal dengan `NotAllowedError` di Linux: WebKitGTK menolak permintaan izin yang tidak ditangani siapa pun, dan sinyal `permission-request` belum terhubung. Kamera/mikrofon kini ditangani berdasarkan peta lintas platform `WebviewWindowOptions.Permissions` yang baru (`map[PermissionType]Permission`), yang dipatuhi di Linux (WebKitGTK) maupun Windows (WebView2). Di Linux, yang tidak memiliki prompt native, kamera/mikrofon diizinkan secara default (memulihkan `getUserMedia`) dan dapat dinonaktifkan dengan `PermissionDeny` (#5552)
- iOS: `GOOS=ios` dapat dikompilasi kembali (`events.IOS` yang diekspor, stub nama metode seluler), dan build dengan tag produksi dapat dikompilasi (perbaikan tag build di pkg/application dan beberapa layanan)
- iOS: event Go→JS dan ExecJS kini berfungsi — halaman tidak lagi dimuat dua kali saat aplikasi dimulai dan handshake `wails:runtime:ready` tidak lagi dapat terlewat
- iOS: `ApplicationDidFinishLaunching`/`ApplicationStarted` tidak lagi mengalami race dengan proses dimulainya aplikasi; jeda tetap selama 2 detik saat aplikasi dimulai telah dihapus
- iOS: perbaiki kebocoran C-string pada setiap eksekusi JavaScript Go→JS
- iOS: `hasListeners` kini mencerminkan pendaftaran listener yang sebenarnya
- iOS: logging debug framework tidak disertakan saat mengompilasi build produksi
- Android: `GOOS=android` dapat dikompilasi kembali — `events.Android` telah didefinisikan, array listener `events_android.go` yang mengakses di luar batas telah dihapus, stub nama metode seluler telah ditambahkan, dan file Linux desktop (`linux_cgo.*`, `events_linux.*`, `environment_linux.go`) tidak lagi ikut masuk ke build Android
- Android: binding JS→Go kini berfungsi — WebView tidak dapat mengirimkan isi POST `fetch()` ke `shouldInterceptRequest`, sehingga panggilan runtime dirutekan melalui transport JavascriptInterface (`nativeHandleRuntimeCall`), bukan mengalami crash akibat isi permintaan nil
- Android: panggilan runtime `Screens.*` mengembalikan data sebenarnya — ScreenManager kini diisi saat aplikasi dimulai (sebelumnya tidak pernah dihubungkan, sehingga `GetAll` mengembalikan nil)
- Android: logging debug framework tidak disertakan dalam build produksi dan diarahkan melalui logcat dengan tag `Wails` dalam build debug
- Android: registri `hasListeners` yang sebenarnya, penanganan referensi/pengecualian JNI, dan siklus hidup halaman dengan pemuatan tunggal (tanpa navigasi ganda)
- Memperbaiki kegagalan `wails3 generate bindings` dengan pesan "Access is denied" pada Windows ketika server pengembangan Vite berjalan, dengan menyinkronkan berkas yang dihasilkan ke direktori keluaran alih-alih mengganti namanya ke atas direktori tersebut (#5515)
- Memperbaiki crash fatal yang terjadi sesekali pada macOS saat membaca informasi layar setelah perubahan tampilan: ID dan nama layar menyimpan pointer ke buffer `UTF8String` yang dirilis otomatis dan dapat dibebaskan sebelum disalin oleh Go (use-after-free). String tersebut kini di-`strdup` dan dibebaskan setelah konversi, sementara enumerasi layar dijalankan dalam autorelease pool eksplisit sehingga tidak lagi membocorkan memori ketika dipanggil dari goroutine Go (#5556)
- Memperbaiki SIGSEGV yang terjadi sesekali pada Linux ketika assetserver menutup `WebKitURISchemeRequest`: `g_object_unref` terakhir berjalan pada goroutine assetserver sehingga menyelesaikan finalisasi GObject WebKit di luar thread utama GTK. Operasi unref kini diteruskan ke konteks utama GTK melalui `g_main_context_invoke` (#5557)

## v3.0.0-alpha.100 - 2026-06-13

## Ditambahkan

- Memperluas `MacWebviewPreferences` dengan opsi konfigurasi WKWebView tambahan: `EnableAutoplayWithoutUserAction`, `AllowsAirPlayForMediaPlayback`, `AllowsMagnification`, `JavaScriptCanOpenWindowsAutomatically`, `MinimumFontSize`, dan `ApplicationNameForUserAgent` (#5549)

## Diperbaiki

- Memperbaiki kegagalan `wails3 generate bindings` dengan pesan "Access is denied" pada Windows ketika server pengembangan Vite berjalan, dengan menyinkronkan berkas yang dihasilkan ke direktori keluaran alih-alih mengganti namanya ke atas direktori tersebut (#5561)
- Memperbaiki peristiwa perubahan ukuran JS yang tidak dipicu untuk jendela tanpa bingkai pada Linux; memperbaiki deteksi tepi bilah gulir untuk jendela tanpa bingkai (#5368)
- Memperbaiki updater pada Windows yang gagal dengan pesan "invalid cross-device link" ketika direktori sementara berada di volume yang berbeda dari direktori instalasi (#5560)

## v3.0.0-alpha.99 - 2026-06-10

## Diperbaiki

- Memperbaiki kegagalan `wails3 generate bindings` dengan pesan "Access is denied" pada Windows ketika server pengembangan Vite berjalan, dengan menyinkronkan berkas yang dihasilkan ke direktori keluaran alih-alih mengganti namanya ke atas direktori tersebut (#5515)

## v3.0.0-alpha.98 - 2026-06-03

## Diperbaiki

- Memperbaiki antarmuka WebKit yang membeku pada Linux saat menganggur (misalnya ketika inspector terbuka) dengan tidak lagi memaksakan `SA_ONSTACK` pada `SIGUSR1`, yang merusak sinkronisasi thread GC JavaScriptCore (#5527)

## v3.0.0-alpha.97 - 2026-05-31

## Ditambahkan

- Menambahkan halaman debugging dan panduan bekerja dengan `runtime/trace`

## Diubah

- Menghapus beberapa impor `_ "embed"` yang tidak diperlukan sehingga kode menjadi sedikit lebih rapi

## Diperbaiki

- Memperbaiki batasan lebar/tinggi minimum yang tidak diterapkan setelah jendela keluar dari keadaan dimaksimalkan pada Windows (#4593)
- Memperbaiki klik mouse yang menembus jendela dalam mode layar penuh saat menggunakan opsi jendela Frameless + Transparent (#4408)

## v3.0.0-alpha.96 - 2026-05-25

## Ditambahkan

- Menambahkan dukungan obfuscation Garble ([#4563](https://github.com/wailsapp/wails/issues/4563)): ID metode binding yang stabil, integrasi build/Taskfile (`build --obfuscated --garbleargs`, `generate bindings -obfuscated`), dan tag struct JSON pada setiap payload yang berhadapan dengan runtime (`EnvironmentInfo`, `OSInfo`, `Screen`, `Rect`, `Point`, `Size`, `Capabilities`) agar format wire tetap berfungsi setelah Garble mengganti nama field yang diekspor.

## v3.0.0-alpha.95 - 2026-05-20

## Ditambahkan

- Menambahkan halaman struktur proyek yang sebelumnya belum tersedia

## Diubah

- Dokumentasi: Mengubah beberapa diagram pada halaman arsitektur agar menggunakan diagram urutan untuk tampilan yang lebih rapi
- Dokumentasi: Menyertakan catatan tentang penginstalan D2 sebagai prasyarat untuk menjalankannya

## Diperbaiki

- Memperbaiki `wails3 generate appimage` pada konfigurasi bawaan GTK4: bundler kini mendeteksi stack GTK dari biner sebelum mencari berkas runtime sehingga memilih `libwebkitgtkinjectedbundle.so` (di bawah `webkitgtk-6.0/`) untuk build GTK4 dan `libwebkit2gtkinjectedbundle.so` (di bawah `webkit2gtk-4.1/`) untuk build `-tags gtk3`. Pemeriksaan `.relr.dyn` juga memeriksa `libgtk-4.so.1` sehingga stripping dinonaktifkan dengan benar pada toolchain modern apa pun stack-nya. (#5475)
- Memperbaiki kegagalan `wails3 generate appimage` ketika dipanggil dengan `-builddir` relatif: bundler kini terlebih dahulu mengubah `-binary`, `-icon`, `-desktopfile`, `-builddir`, dan `-outputdir` menjadi path absolut agar `s.CD` di tengah alur tidak mengganggu goroutine pengunduhan AppRun maupun pemeriksaan `ldd` setelah penyalinan.
- Memperbaiki kegagalan `wails3 generate appimage` memindahkan AppImage akhir ke `-outputdir` ketika field `Name=` desktop tidak cocok dengan nama dasar biner: bundler kini memaksa plugin appimage milik linuxdeploy (melalui variabel lingkungan `OUTPUT`) untuk menulis AppImage ke `<binary>-<arch>.AppImage` alih-alih menggunakan nama yang diturunkan dari berkas desktop.
- Memperbaiki `events.Common.ApplicationStarted`, `Common.ThemeChanged`, `Common.SystemWillSleep`, dan `Common.SystemDidWake` yang tidak dipicu pada Linux setelah stack GTK4 + WebKitGTK 6.0 dijadikan bawaan pada alpha.93. `run()` `application_linux.go` bawaan yang baru tidak memanggil `setupCommonEvents()` (yang meneruskan peristiwa `Linux.*` ke padanan `Common.*`-nya) maupun `monitorPowerEvents()`. Helper pemantau daya DBus kini digunakan bersama oleh jalur build GTK3 dan GTK4 melalui `application_linux_dbus.go`. (#5474)

## v3.0.0-alpha.94 - 2026-05-19

## Diperbaiki

- Memperbaiki `events.Common.ApplicationStarted`, `Common.ThemeChanged`, `Common.SystemWillSleep`, dan `Common.SystemDidWake` yang tidak dipicu pada Linux setelah stack GTK4 + WebKitGTK 6.0 dijadikan bawaan pada alpha.93. `run()` `application_linux.go` bawaan yang baru tidak memanggil `setupCommonEvents()` (yang meneruskan peristiwa `Linux.*` ke padanan `Common.*`-nya) maupun `monitorPowerEvents()`. Helper pemantau daya DBus kini digunakan bersama oleh jalur build GTK3 dan GTK4 melalui `application_linux_dbus.go`. (#5474)

## v3.0.0-alpha.93 - 2026-05-17

## Ditambahkan

- Menambahkan `XDG_SESSION_TYPE` ke keluaran `wails3 doctor` pada Linux oleh @leaanthony

## Diperbaiki

- Memperbaiki crash menu jendela pada Wayland yang disebabkan oleh appmenu-gtk-module saat mengakses jendela yang belum direalisasikan (#4769) oleh @leaanthony
- Memperbaiki crash aplikasi GTK ketika nama aplikasi mengandung karakter yang tidak valid (spasi, tanda kurung, dan sebagainya) oleh @leaanthony
- Memperbaiki galat "not enough memory" saat menginisialisasi seret dan lepas pada Windows (#4701) oleh @overlordtm
- Memperbaiki kondisi balapan pada penyimpanan callback mainthread akibat penggunaan RLock yang keliru untuk menghapus entri map (Linux, macOS, iOS) (#4424) oleh @leaanthony
- Memperbaiki penanganan variabel ketika argumen baris perintah diteruskan ke task. Variabel CLI yang ditentukan sebagai pasangan KEY=VALUE kini diinisialisasi dengan benar dan disebarkan ke seluruh eksekusi task.
- Memperbaiki konflik NSWindowZoomButton pada macOS: `MaximiseButtonState` dan `FullscreenButtonState` kini menerapkan keadaan yang lebih ketat, baik saat startup maupun runtime; tidak satu pun setter dapat menggantikan yang lain secara diam-diam (#5319)
- Memperbaiki sekumpulan bug lama yang sudah ada pada jalur build GTK3 lawas (`-tags gtk3`) yang ditemukan oleh CodeRabbit pada #5463: peluncuran melalui asosiasi file tidak lagi melewati handler startup; `getTheme` kini aman terhadap batas dan tipe; `appName` tidak lagi membebaskan memori milik GLib; `clipboardGet` tidak lagi membocorkan `gchar*` yang dikembalikan oleh GTK; `Calloc` kini menggunakan receiver pointer (dan `NewCalloc` mengembalikan `*Calloc`) sehingga pool benar-benar melacak alokasi; `zoomOut` menggunakan kebalikan dari `zoomInFactor` alih-alih pengali negatif yang dibatasi menjadi 1.0; `execJS` menggunakan kembali nama world kosong yang telah dialokasikan sebelumnya, alih-alih membocorkan satu `C.CString("")` per pemanggilan; satu `fmt.Println` pengembangan telah dihapus dari `menuItem.setAccelerator`. Menyelesaikan #5465.
- Memperbaiki kebocoran receiver nilai `Calloc` yang sama pada jalur build GTK4 default (`linux_cgo.go`): receiver pointer + `NewCalloc() *Calloc` agar alokasi `c.String(...)` per jendela benar-benar dilacak dan dibebaskan.

## v3.0.0-alpha.92 - 2026-05-15

## Ditambahkan

- Mengubah Taskfile agar pengelola paket frontend yang digunakan dapat dikontrol melalui opsi `PACKAGE_MANAGER`
- Memperkaya data templat dengan `{{.Opn}}` dan `{{.Cls}}` agar penulisan templat Taskfile lebih mudah diprediksi

## Diubah

- Mengubah beberapa Taskfile yang ada agar menggunakan `{{.Opn}} and {{.Cls}}`

## Diperbaiki

- Memperbaiki kegagalan fatal runtime `concurrent map read and map write` di `linuxSystemTray` saat menu tray diperbarui ketika panel sedang membacanya.
- Menggunakan `log` alih-alih `fmt` untuk keluaran galat dan pelacakan tumpukan WebView2 agar pesan tidak hilang ketika aplikasi berjalan tanpa konsol yang terhubung di Windows.

## v3.0.0-alpha.91 - 2026-05-12

## Diubah

- Memperbarui SVG sponsor dalam [PR](https://github.com/wailsapp/wails/pull/5414) oleh `@github-actions[bot]`
- **PERUBAHAN INKOMPATIBEL (macOS):** Menyeragamkan sistem koordinat macOS agar `GetScreens`, `Position`, dan `SetPosition` semuanya menggunakan ruang yang sama—titik logis, Y ke bawah, dengan `(0,0)` di kiri atas layar utama. Ini sesuai dengan Windows, GTK, serta API publik Electron dan web. Layar yang secara fisik berada di atas layar utama kini melaporkan `Bounds.Y` negatif (sebelumnya positif), dan nilai `Position()`/`SetPosition()` kini dinyatakan dalam titik logis, bukan `points × primaryScale`. Konversi bolak-balik `Position()` → `SetPosition()` tetap dipertahankan; nilai absolut yang dicatat dari build alfa sebelumnya atau solusi sementara yang dihitung manual (misalnya mengalikan dengan `primaryScale` atau membalik Y terhadap tinggi layar) perlu diperbarui. Menyelesaikan [#5117](https://github.com/wailsapp/wails/issues/5117).

## Diperbaiki

- Memvalidasi nama sinyal DBus dan panjang body secara defensif untuk mencegah panic dalam [PR](https://github.com/wailsapp/wails/pull/5416) oleh @leaanthony
- Memperbaiki masalah keamanan memori pada penanganan menu GTK di Linux dalam [PR](https://github.com/wailsapp/wails/pull/5363) oleh @leaanthony
- Mendeteksi GPU NVIDIA dan menonaktifkan renderer DMA-BUF di Linux dalam [PR](https://github.com/wailsapp/wails/pull/5295) oleh @leaanthony
- Memperbaiki konversi Y lintas layar `SetPosition` di macOS: menggunakan tinggi layar utama sebagai acuan global agar jendela ditempatkan pada posisi yang tepat di monitor yang bergeser secara vertikal dari layar utama dalam [#5117](https://github.com/wailsapp/wails/issues/5117)
- Memperbaiki templat PR Git agar mengarah ke URL umpan balik yang benar dalam [PR](https://github.com/wailsapp/wails/pull/5109) oleh @wayneforrest
- Memperbaiki serangkaian crash `SetMenu` pada systray Windows yang disebabkan oleh syscall `DestroyMenu` yang rusak karena meneruskan empat argumen alih-alih satu, sehingga setiap pemanggilan mengembalikan FALSE dan tidak membebaskan apa pun. Juga membebaskan handle HMENU dan HBITMAP (termasuk yang dialokasikan saat runtime melalui `MenuItem.SetBitmap`) ketika menu dibangun ulang, mereset map kotak centang/radio yang usang di `Win32Menu.Update`, serta menghapus pemanggilan `Update()` berlebih di `systemtray.updateMenu` yang menggandakan alokasi. Aplikasi systray yang berjalan lama tidak lagi membocorkan objek GDI/USER setiap kali menu dibangun ulang.

## v3.0.0-alpha.90 - 2026-05-11

## Ditambahkan

- Menambahkan nama aplikasi yang dapat dikonfigurasi untuk User-Agent WKWebView di macOS dalam [PR](https://github.com/wailsapp/wails/pull/5261) oleh @vinhvoit225
- Menambahkan dependensi tidak langsung github.com/coder/websocket ke contoh gin-service dalam [PR](https://github.com/wailsapp/wails/pull/5400) oleh @taliesin-ai
- Menambahkan dukungan perbandingan kesetaraan mendalam ke pengujian aset build dalam [PR](https://github.com/wailsapp/wails/pull/5402) oleh @leaanthony

## Diubah

- Menggabungkan keluaran build ke dalam direktori assets dalam [PR](https://github.com/wailsapp/wails/pull/5401) oleh @taliesin-ai
- Memperbarui SVG sponsor dalam [PR](https://github.com/wailsapp/wails/pull/5399) oleh `@github-actions[bot]`

## Diperbaiki

- Menggunakan objek notifikasi untuk pesan instans tunggal macOS dalam [PR](https://github.com/wailsapp/wails/pull/5289) oleh @overlordtm
- Memproses callback Windows secara batch untuk mencegah hilangnya promise saat beban tinggi dalam [PR](https://github.com/wailsapp/wails/pull/5383) oleh @taliesin-ai

## v3.0.0-alpha.89 - 2026-05-10

## Ditambahkan

- Menambahkan job go<em>test</em>results untuk mengagregasi hasil pengujian Go dalam [PR](https://github.com/wailsapp/wails/pull/5316) oleh @leaanthony

## Diubah

- Memecah payload RPC besar secara kondisional menjadi permintaan POST berbasis chunk dalam [PR](https://github.com/wailsapp/wails/pull/5369) oleh @leaanthony
- Meningkatkan Vite dari 5.x.x ke 8.0.0 di semua templat frontend dalam [PR](https://github.com/wailsapp/wails/pull/5386) oleh @leaanthony
- Memigrasikan konfigurasi port server pengembangan Vite ke variabel lingkungan dalam [PR](https://github.com/wailsapp/wails/pull/5365) oleh @leaanthony
- Mengonfigurasi server pengembangan Vite agar mengikat ke 127.0.0.1 di semua templat dalam [PR](https://github.com/wailsapp/wails/pull/5361) oleh @leaanthony
- Memperbarui SVG sponsor dalam [PR](https://github.com/wailsapp/wails/pull/5384) oleh `@github-actions[bot]`

## Diperbaiki

- Membersihkan stub templat Info.plist selama pembaruan build-assets dalam [PR](https://github.com/wailsapp/wails/pull/5312) oleh @leaanthony
- Memperbaiki status usang pada menu macOS dengan menerapkan mutator item menu (`setMenuItemChecked()`, `setMenuItemLabel()`, `setMenuItemDisabled()`, `setMenuItemHidden()`, `setMenuItemTooltip()`) secara sinkron di thread utama, sehingga menghilangkan race `dispatch_async` yang menyebabkan menu merender status sebelumnya ketika dibuka kembali dengan cepat (#5002)
- Mengabaikan file `*_test.go` dalam mode pengembangan untuk mencegah build ulang yang tidak diperlukan dalam [PR](https://github.com/wailsapp/wails/pull/5203) oleh @leaanthony
- Mencegah segfault Menu.Update() saat aplikasi tidak berjalan dalam [PR](https://github.com/wailsapp/wails/pull/5291) oleh @wucm667
- Menggunakan lastSizeWParam untuk membatasi penggambaran ulang bilah menu di Windows dalam [PR](https://github.com/wailsapp/wails/pull/5382) oleh @taliesin-ai

## v3.0.0-alpha.88 - 2026-05-09

## Diubah

- Mengalihkan HiddenOnTaskbar agar menggunakan WS<em>EX</em>TOOLWINDOW dalam [PR](https://github.com/wailsapp/wails/pull/5371) oleh @leaanthony
- Mengurutkan ulang dependensi dan menghapus direktif replace webview2 di go.mod dalam [PR](https://github.com/wailsapp/wails/pull/5370) oleh @atterpac
- Memperbarui SVG sponsor dalam [PR](https://github.com/wailsapp/wails/pull/5358) oleh `@github-actions[bot]`

## Diperbaiki

- Menghapus alias indirection generik dan menggabungkan tipe kunci map dalam [PR](https://github.com/wailsapp/wails/pull/5331) oleh @fbbdev

## Dihapus

- Hapus alur kerja PR-master, termasuk dokumentasi, pengujian Go, dan pengujian yang dilewati, dalam [PR](https://github.com/wailsapp/wails/pull/5377) oleh @leaanthony

## v3.0.0-alpha.87 - 2026-05-07

## Ditambahkan

- Tambahkan dokumentasi Wails v3 dalam bahasa Korea melalui [PR](https://github.com/wailsapp/wails/pull/5352) oleh @leaanthony
- Tambahkan dokumentasi instalasi dan panduan memulai cepat dalam bahasa Prancis melalui [PR](https://github.com/wailsapp/wails/pull/5354) oleh @leaanthony
- Tambahkan dokumentasi panduan memulai cepat, konsep, dan komunitas dalam bahasa Portugis melalui [PR](https://github.com/wailsapp/wails/pull/5355) oleh @leaanthony

## v3.0.0-alpha.86 - 2026-05-06

## Ditambahkan

- Tambahkan pelokalan dokumentasi dalam bahasa Prancis melalui [PR](https://github.com/wailsapp/wails/pull/5328) oleh @leaanthony
- Tambahkan lokal bahasa Jerman ke situs dokumentasi melalui [PR](https://github.com/wailsapp/wails/pull/5343) oleh @leaanthony

## Diubah

- Daftarkan seluruh 8 lokal terjemahan dalam konfigurasi dokumentasi melalui [PR](https://github.com/wailsapp/wails/pull/5347) oleh @leaanthony
- Perbarui berbagai berkas terkait Windows untuk WebView2 melalui [PR](https://github.com/wailsapp/wails/pull/5317) oleh @leaanthony

## Diperbaiki

- Pisahkan pengiriman dialog antara GTK3 dan GTK4 untuk Linux melalui [PR](https://github.com/wailsapp/wails/pull/5340) oleh @leaanthony
- Pastikan callback dialog dijalankan pada thread GTK untuk memperbaiki segfault melalui [PR](https://github.com/wailsapp/wails/pull/5339) oleh @leaanthony

## v3.0.0-alpha.85 - 2026-05-05

## Ditambahkan

- Tambahkan URL templat PR ke repositori melalui [PR](https://github.com/wailsapp/wails/pull/5179) oleh @leaanthony
- Tambahkan dokumentasi Wails v3 dalam bahasa Jerman melalui [PR](https://github.com/wailsapp/wails/pull/5330) oleh @leaanthony

## v3.0.0-alpha.84 - 2026-05-03

## Ditambahkan

- Tambahkan opsi untuk mencegah tombol Escape keluar dari mode layar penuh di macOS melalui [PR](https://github.com/wailsapp/wails/pull/5307) oleh @leaanthony
- Tambahkan opsi untuk mencegah tombol Escape keluar dari mode layar penuh di macOS melalui [PR](https://github.com/wailsapp/wails/pull/5310) oleh @leaanthony
- Tambahkan dokumentasi etalase komunitas Pausa melalui [PR](https://github.com/wailsapp/wails/pull/5288) oleh @yuseferi

## Diubah

- Perbarui SVG sponsor melalui [PR](https://github.com/wailsapp/wails/pull/5308) oleh `@github-actions[bot]`
- Perbarui perintah pembuatan ikon agar menangani platform yang tidak didukung melalui [PR](https://github.com/wailsapp/wails/pull/5309) oleh @leaanthony
- Ganti API layar penuh boolean dengan ButtonState tiga status dan implementasikan binding platform melalui [PR](https://github.com/wailsapp/wails/pull/5224) oleh @leaanthony

## Diperbaiki

- Lindungi operasi fokus WebView2 dari status controller nil melalui [PR](https://github.com/wailsapp/wails/pull/5315) oleh @leaanthony
- Perbarui alur kerja GitHub Actions agar merujuk cabang dasar PR dengan benar melalui [PR](https://github.com/wailsapp/wails/pull/5313) oleh @leaanthony
- Abaikan berkas `*_test.go` dalam mode pengembangan untuk mencegah pembangunan ulang yang tidak diperlukan melalui [PR](https://github.com/wailsapp/wails/pull/5203) oleh @leaanthony
- Cegah segfault Menu.Update() saat aplikasi tidak berjalan melalui [PR](https://github.com/wailsapp/wails/pull/5291) oleh @wucm667

## v3.0.0-alpha.83 - 2026-05-02

## Ditambahkan

- Tambahkan flag InstallScope dan opsi build untuk instalasi tingkat mesin/pengguna melalui [PR](https://github.com/wailsapp/wails/pull/5094) oleh @symball
- Tambahkan metode SetScreen tanpa operasi ke BrowserWindow agar memenuhi antarmuka Window melalui [PR](https://github.com/wailsapp/wails/pull/5294) oleh @leaanthony

## Diperbaiki

- Deteksi GPU NVIDIA dan nonaktifkan perender DMA-BUF di Linux melalui [PR](https://github.com/wailsapp/wails/pull/5295) oleh @leaanthony
- Perbaiki templat PR Git agar mengarah ke URL umpan balik yang benar melalui [PR](https://github.com/wailsapp/wails/pull/5109) oleh @wayneforrest
- Perbaiki serangkaian crash systray Windows `SetMenu` yang disebabkan oleh syscall `DestroyMenu` yang rusak karena meneruskan empat argumen alih-alih satu, sehingga setiap panggilan mengembalikan FALSE dan tidak membebaskan apa pun. Lepaskan juga handle HMENU dan HBITMAP (termasuk yang dialokasikan saat runtime melalui `MenuItem.SetBitmap`) ketika menu dibangun ulang, atur ulang map kotak centang/tombol radio yang kedaluwarsa dalam `Win32Menu.Update`, dan hapus panggilan `Update()` yang redundan dalam `systemtray.updateMenu` karena menggandakan alokasi. Aplikasi systray yang berjalan lama kini tidak lagi membocorkan objek GDI/USER setiap kali menu dibangun ulang.

## v3.0.0-alpha.82 - 2026-05-01

## Diperbaiki

- Perbaiki pembuatan berkas desktop agar menangani nama desktop dengan benar melalui [PR](https://github.com/wailsapp/wails/pull/5232) oleh @leaanthony

## v3.0.0-alpha.81 - 2026-04-30

## Diubah

- Sesuaikan jadwal rilis nightly menjadi pukul 15:00 UTC melalui [PR](https://github.com/wailsapp/wails/pull/5286) oleh @leaanthony

## Diperbaiki

- Perbaiki nilai Screen Bounds, WorkArea, dan Size yang berkurang separuh pada Mac Retina -  (#5168)

## v3.0.0-alpha.80 - 2026-04-29

## Diubah

- Perbarui dependensi dokumentasi dan pemuat koleksi konten melalui [PR](https://github.com/wailsapp/wails/pull/5285) oleh @leaanthony

## v3.0.0-alpha.79 - 2026-04-29

## Ditambahkan

- Berikan izin actions: write kepada job trigger-release melalui [PR](https://github.com/wailsapp/wails/pull/5270) oleh @leaanthony

## Diubah

- Tugas rilis kini menggunakan cabang master secara default dan memperbarui redaksi changelog melalui [PR](https://github.com/wailsapp/wails/pull/5283) oleh @leaanthony
- Perbarui alur kerja auto-changelog agar menggunakan versi terbaru melalui [PR](https://github.com/wailsapp/wails/pull/5282) oleh @leaanthony
- Tingkatkan efisiensi alur kerja dengan menambahkan filter jalur dan menghapus alur kerja yang tidak digunakan melalui [PR](https://github.com/wailsapp/wails/pull/5280) oleh @leaanthony
- Perbarui dokumentasi agar tautan contoh merujuk ke cabang master melalui [PR](https://github.com/wailsapp/wails/pull/5274) oleh @leaanthony
- Perbarui dokumentasi dan contoh untuk v3 melalui [PR](https://github.com/wailsapp/wails/pull/5272) oleh @leaanthony

## Diperbaiki

- Tingkatkan proksi balik dengan logika percobaan ulang dan pemaksaan IPv4 untuk pengembangan melalui [PR](https://github.com/wailsapp/wails/pull/5265) oleh @AkagiYui
- Tulis ulang alur kerja pemicu changelog yang belum dirilis melalui [PR](https://github.com/wailsapp/wails/pull/5281) oleh @leaanthony

## Dihapus

- Hapus skrip pengujian shell untuk berbagai keperluan pengujian melalui [PR](https://github.com/wailsapp/wails/pull/5267) oleh @leaanthony
- Hapus alur kerja penerapan dokumentasi v3-alpha dan record CNAME melalui [PR](https://github.com/wailsapp/wails/pull/5266) oleh @leaanthony

### Ditambahkan

- Tambahkan entri Perutean Frontend ke navigasi bilah sisi melalui [PR](https://github.com/wailsapp/wails/pull/5196) oleh @leaanthony
- Tambahkan panduan perutean frontend dengan rekomendasi khusus framework melalui [PR](https://github.com/wailsapp/wails/pull/5185) oleh @leaanthony
- Menambahkan dukungan untuk sheet modal (macOS)
- Meningkatkan versi ghw untuk dukungan perangkat Apple yang lebih baik oleh @leaanthony (#4977)
- Menambahkan metode `GetBadge` ke layanan dock
- Menambahkan flag `-tags` ke perintah `wails3 build` untuk meneruskan tag build Go khusus (misalnya, `wails3 build -tags gtk4`) (#4957)
- Menambahkan dokumentasi untuk pembuatan enum otomatis dalam generator binding, termasuk halaman khusus Enum dan navigasi bilah sisi (#4972)
- Menambahkan flag `-tags` ke perintah `wails3 build` untuk meneruskan tag build Go khusus (misalnya, `wails3 build -tags gtk4`) (#4957)
- Menambahkan contoh Web API dalam `v3/examples/web-apis/` yang mendemonstrasikan API peramban 41, termasuk Penyimpanan (localStorage, sessionStorage, IndexedDB, Cache API), Jaringan (Fetch, WebSocket, XMLHttpRequest, EventSource, Beacon), Media (Canvas, WebGL, Web Audio, MediaDevices, MediaRecorder, Speech Synthesis), Perangkat (Geolocation, Clipboard, Fullscreen, Device Orientation, Vibration, Gamepad), Performa (Performance API, Mutation Observer, Intersection/Resize Observer), UI (Web Components, Pointer Events, Selection, Dialog, Drag and Drop), dan lainnya
- Menambahkan contoh pemeriksa kompatibilitas WebView API (`v3/examples/webview-api-check/`) yang menguji lebih dari 200 API peramban di berbagai platform
- Menambahkan paket `internal/libpath` untuk menemukan jalur pustaka native di Linux dengan pencarian paralel, caching, dan dukungan untuk Flatpak/Snap/Nix
- **WIP:** Menambahkan dukungan eksperimental WebKitGTK 6.0 / GTK4 untuk Linux, yang tersedia melalui `-tags gtk4` (GTK3/WebKit2GTK 4.1 tetap menjadi default)
- Catatan: Pada window manager tiling (misalnya, Hyprland, Sway), operasi Minimize/Maximize mungkin tidak berfungsi seperti yang diharapkan karena WM mengontrol geometri jendela
- Menambahkan cara menggunakan **Handler Satu Kali** dalam dokumentasi **Memantau Event di JavaScript** oleh @AbdelhadiSeddar
- Menambahkan opsi `UseApplicationMenu` ke `WebviewWindowOptions` agar jendela di Windows/Linux dapat mewarisi menu aplikasi yang ditetapkan melalui `app.Menu.Set()` oleh @leaanthony
- Menambahkan dukungan untuk menggunakan file `.icon` (format Apple Icon Composer) guna menghasilkan ikon Liquid Glass dan katalog aset (macOS) (#4934) oleh @wimaha
- Menambahkan mode server eksperimental untuk deployment headless/web (`-tags server`). Memungkinkan aplikasi Wails berjalan sebagai server HTTP tanpa dependensi GUI native. Build dengan `wails3 task build:server`. Lihat `examples/server` untuk detailnya.
- Menambahkan paket `internal/libpath` untuk menemukan jalur pustaka native di Linux dengan pencarian paralel, caching, dan dukungan untuk Flatpak/Snap/Nix
- Menambahkan opsi `CollectionBehavior` ke `MacWindow` untuk mengontrol perilaku jendela di macOS Spaces dan layar penuh (#4756) oleh @leaanthony
- Menambahkan pengujian unit untuk pkg/application oleh @leaanthony
- Menambahkan dukungan protokol khusus ke pemaketan MSIX oleh @leaanthony
- Menambahkan deteksi lingkungan desktop di Linux [PR #4797](https://github.com/wailsapp/wails/pull/4797)
- Menambahkan metode `Window.Print()` ke runtime JavaScript untuk memicu dialog cetak dari frontend (#4290) oleh @leaanthony
- Menambahkan `XDG_SESSION_TYPE` ke output `wails3 doctor` di Linux oleh @leaanthony
- Menambahkan event perubahan pemuatan WebKit2 tambahan untuk Linux: `WindowLoadStarted`, `WindowLoadRedirected`, `WindowLoadCommitted`, `WindowLoadFinished` (#3896) oleh @leaanthony
- Menambahkan `XDG_SESSION_TYPE` ke output `wails3 doctor` di Linux oleh @leaanthony
- Menghasilkan file `.desktop` saat build Linux, bukan hanya saat pemaketan (#4575)
- Menambahkan dokumentasi dependensi runtime Linux dengan nama paket khusus distro dan contoh pemaketan nfpm (#4339) oleh @leaanthony
- Menambahkan informasi versi driver NVIDIA ke output `wails3 doctor` di Linux oleh @leaanthony
- Menambahkan origin ke handler pesan mentah oleh @APshenkin dalam [PR](https://github.com/wailsapp/wails/pull/4710)
- Menambahkan dukungan universal link untuk macOS oleh @APshenkin dalam [PR](https://github.com/wailsapp/wails/pull/4712)
- Memfaktorkan ulang lapisan transport binding oleh @APshenkin dalam [PR](https://github.com/wailsapp/wails/pull/4702)
- Menambahkan pengidentifikasi aria-label ke templat helloworld agar aplikasi contoh dapat diuji dengan mudah oleh klien pengujian Appium, oleh @chinenual dalam [PR](https://github.com/wailsapp/wails/pull/4760)
- Menambahkan origin ke handler pesan mentah oleh @APshenkin dalam [PR](https://github.com/wailsapp/wails/pull/4710)
- Menambahkan dukungan universal link untuk macOS oleh @APshenkin dalam [PR](https://github.com/wailsapp/wails/pull/4712)
- Memfaktorkan ulang lapisan transport binding oleh @APshenkin dalam [PR](https://github.com/wailsapp/wails/pull/4702)
- Event Bertipe oleh @fbbdev dan @ianvs dalam [#4633](https://github.com/wailsapp/wails/pull/4633)
- Menambahkan contoh `systray-clock` yang menampilkan tray headless dengan pembaruan tooltip secara langsung (#4653).
- Menambahkan templat Protokol NSIS untuk Windows oleh @Tolfx dalam #4510
- Menambahkan pengujian untuk build-assets oleh @Tolfx dalam #4510
- macOS: Menampilkan kontrol jendela native di bilah menu dalam [#4588](https://github.com/wailsapp/wails/pull/4588) oleh @nidib
- Menambahkan layanan Dock macOS untuk menyembunyikan/menampilkan ikon aplikasi di dock oleh @popaprozac dalam [PR](https://github.com/wailsapp/wails/pull/4451)
- Menambahkan layanan Dock macOS untuk menyembunyikan/menampilkan ikon aplikasi di dock oleh @popaprozac dalam [PR](https://github.com/wailsapp/wails/pull/4451)
- Menambahkan dukungan efek Liquid Glass native untuk macOS dengan NSGlassEffectView (macOS 15.0+) serta fallback NSVisualEffectView, termasuk opsi kustomisasi material yang komprehensif, oleh @leaanthony dalam [#4534](https://github.com/wailsapp/wails/pull/4534)
- Sanitasi URL Peramban oleh @leaanthony dalam [#4500](https://github.dev/wailsapp/wails/pull/4500). Berdasarkan [#4484](https://github.com/wailsapp/wails/pull/4484) oleh @APShenkin.
- Menambahkan Perlindungan Konten di Windows/Mac oleh [@leaanthony](https://github.com/leaanthony) berdasarkan karya awal [@Taiterbase](https://github.com/Taiterbase) dalam [PR](https://github.com/wailsapp/wails/pull/4241) ini
- Menambahkan dukungan untuk meneruskan variabel CLI ke perintah Task melalui alias `wails3 build` dan `wails3 package` (#4422) oleh @leaanthony dalam [PR](https://github.com/wailsapp/wails/pull/4488)
- Dukungan untuk dropzone dengan event yang menyediakan data elemen yang dijatuhkan, oleh [@atterpac](https://github.com/atterpac) dalam [#4318](https://github.com/wailsapp/wails/pull/4318)
- Menambahkan `AdditionalLaunchArgs` ke opsi `WindowsWindow` agar argumen baris perintah tambahan dapat diteruskan ke peramban WebView2, dalam [PR](https://github.com/wailsapp/wails/pull/4467)
- Menambahkan eksekusi go mod tidy secara otomatis setelah wails init oleh [@triadmoko](https://github.com/triadmoko) dalam [PR](https://github.com/wailsapp/wails/pull/4286)
- Fitur Windows Snapassist oleh @leaanthony dalam [PR](https://github.dev/wailsapp/wails/pull/4463)
- Menambahkan `AdditionalLaunchArgs` ke opsi `WindowsWindow` agar argumen baris perintah tambahan dapat diteruskan ke peramban WebView2, dalam [PR](https://github.com/wailsapp/wails/pull/4467)
- Menambahkan eksekusi go mod tidy secara otomatis setelah wails init oleh [@triadmoko](https://github.com/triadmoko) dalam [PR](https://github.com/wailsapp/wails/pull/4286)
- Fitur Windows Snapassist oleh @leaanthony dalam [PR](https://github.dev/wailsapp/wails/pull/4463)
- Menambahkan implementasi `getAccentColor` Windows oleh [@almas-x](https://github.com/almas-x) dalam [PR](https://github.com/wailsapp/wails/pull/4427)
- Menambahkan implementasi `getAccentColor` Windows oleh [@almas-x](https://github.com/almas-x) dalam [PR](https://github.com/wailsapp/wails/pull/4427)
- Menu dan bilah menu bertema gelap di Windows. Oleh @leaanthony dalam [a29b4f0861b1d0a700e9eb213c6f1076ec40efd5](https://github.com/wailsapp/wails/commit/a29b4f0861b1d0a700e9eb213c6f1076ec40efd5)
- Mengganti nama layanan bawaan agar binding JS/TS lebih jelas oleh @popaprozac dalam [PR](https://github.com/wailsapp/wails/pull/4405)
- `app.Env.GetAccentColor` untuk mendapatkan warna aksen sistem pengguna. Berfungsi di macOS. Oleh [@etesam913](https://github.com/etesam913)
- Menambahkan API `window.ToggleFrameless()` oleh [@atterpac](https://github.com/atterpac) dalam [#4137](https://github.com/wailsapp/wails/pull/4137)
- Menambahkan dependensi build khusus distribusi untuk Linux oleh @leaanthony dalam [PR](https://github.com/wailsapp/wails/pull/4345)
- Menambahkan panduan binding oleh @atterpac dalam [PR](https://github.com/wailsapp/wails/pull/4404)
- **Infrastruktur Pengujian yang Tertata**: Memindahkan berkas pengujian Docker ke direktori khusus `test/docker/` dengan image yang dioptimalkan serta keandalan build yang ditingkatkan oleh [@leaanthony](https://github.com/leaanthony) dalam [#4359](https://github.com/wailsapp/wails/pull/4359)
- **Pola Pengelolaan Sumber Daya yang Disempurnakan**: Menambahkan pembersihan event handler yang tepat dan pengelolaan goroutine yang memperhitungkan konteks dalam contoh oleh [@leaanthony](https://github.com/leaanthony) dalam [#4359](https://github.com/wailsapp/wails/pull/4359)
- Mendukung build AppImage aarch64 oleh [@AkshayKalose](https://github.com/AkshayKalose) dalam [#3981](https://github.com/wailsapp/wails/pull/3981)
- Menambahkan bagian diagnostik ke `wails doctor` oleh [@leaanthony](https://github.com/leaanthony)
- Menambahkan jendela ke konteks saat memanggil metode layanan oleh [@leaanthony](https://github.com/leaanthony)
- Menambahkan contoh `window-call` untuk menunjukkan cara mengetahui jendela mana yang memanggil layanan oleh [@leaanthony](https://github.com/leaanthony)
- Panduan Menu baru oleh [@leaanthony](https://github.com/leaanthony)
- Penanganan panic yang lebih baik oleh [@leaanthony](https://github.com/leaanthony)
- Panduan Menu baru oleh [@leaanthony](https://github.com/leaanthony)
- Menambahkan komentar dokumentasi untuk API Service oleh [@fbbdev](https://github.com/fbbdev) dalam [#4024](https://github.com/wailsapp/wails/pull/4024)
- Menambahkan fungsi `application.NewServiceWithOptions` untuk menginisialisasi layanan dengan konfigurasi tambahan oleh [@leaanthony](https://github.com/leaanthony) dalam [#4024](https://github.com/wailsapp/wails/pull/4024)
- Meningkatkan kontrol menu oleh [@FalcoG](https://github.com/FalcoG) dan [@leaanthony](https://github.com/leaanthony) dalam [#4031](https://github.com/wailsapp/wails/pull/4031)
- Dokumentasi tambahan oleh [@leaanthony](https://github.com/leaanthony)
- Mendukung pembatalan event dalam event listener standar oleh [@leaanthony](https://github.com/leaanthony)
- Dukungan `Hide`, `Show`, dan `Destroy` untuk systray oleh [@leaanthony](https://github.com/leaanthony)
- Dukungan `SetTooltip` untuk systray oleh [@leaanthony](https://github.com/leaanthony). Ide awal oleh [@lujihong](https://github.com/wailsapp/wails/issues/3487#issuecomment-2633242304)
- Melaporkan path paket dalam peringatan generator binding tentang tipe yang tidak didukung oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Menambahkan dukungan generator binding untuk alias generik oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Menambahkan dukungan generator binding untuk flag JSON `omitzero` oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Menambahkan direktif `//wails:ignore` untuk mencegah pembuatan binding bagi metode layanan yang dipilih oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Menambahkan direktif `//wails:internal` pada layanan dan model untuk mengizinkan tipe yang diekspor di Go tetapi tidak di JS/TS oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Menambahkan dukungan generator binding untuk konstanta bertipe alias guna memungkinkan enum dengan pengetikan longgar oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Menambahkan pengujian generator binding untuk fitur Go 1.24 oleh [@fbbdev](https://github.com/fbbdev) dalam [#4068](https://github.com/wailsapp/wails/pull/4068)
- Menambahkan dukungan untuk macOS 15 "Sequoia" ke `OSInfo.Branding` guna meningkatkan deteksi versi OS dalam [#4065](https://github.com/wailsapp/wails/pull/4065)
- Menambahkan hook `PostShutdown` untuk menjalankan kode khusus setelah proses shutdown selesai oleh [@fbbdev](https://github.com/fbbdev) dalam [#4066](https://github.com/wailsapp/wails/pull/4066)
- Menambahkan struct `FatalError` untuk mendukung deteksi kesalahan fatal dalam error handler khusus oleh [@fbbdev](https://github.com/fbbdev) dalam [#4066](https://github.com/wailsapp/wails/pull/4066)
- Menstandarkan dan mendokumentasikan urutan startup dan shutdown layanan oleh [@fbbdev](https://github.com/fbbdev) dalam [#4066](https://github.com/wailsapp/wails/pull/4066)
- Menambahkan test harness untuk urutan startup/shutdown aplikasi dan pengujian startup/shutdown layanan oleh [@fbbdev](https://github.com/fbbdev) dalam [#4066](https://github.com/wailsapp/wails/pull/4066)
- Menambahkan metode `RegisterService` untuk mendaftarkan layanan setelah aplikasi dibuat oleh [@fbbdev](https://github.com/fbbdev) dalam [#4066](https://github.com/wailsapp/wails/pull/4066)
- Menambahkan field `MarshalError` dalam opsi aplikasi dan layanan untuk penanganan kesalahan khusus pada pemanggilan binding oleh [@fbbdev](https://github.com/fbbdev) dalam [#4066](https://github.com/wailsapp/wails/pull/4066)
- Menambahkan wrapper promise yang dapat dibatalkan dan meneruskan permintaan pembatalan melalui rantai promise oleh [@fbbdev](https://github.com/fbbdev) dalam [#4100](https://github.com/wailsapp/wails/pull/4100)
- Menambahkan kemampuan untuk mengaitkan pembatalan pemanggilan binding dengan `AbortSignal` oleh [@fbbdev](https://github.com/fbbdev) dalam [#4100](https://github.com/wailsapp/wails/pull/4100)
- Mendukung atribut `data-wml-*` untuk WML selain atribut `wml-*` yang biasa digunakan oleh [@leaanthony](https://github.com/leaanthony)
- Menambahkan metode `Configure` pada semua layanan untuk konfigurasi tahap lanjut/konfigurasi ulang dinamis oleh [@fbbdev](https://github.com/fbbdev) dalam [#4067](https://github.com/wailsapp/wails/pull/4067)
- Layanan `fileserver` mengirim respons 503 Service Unavailable ketika belum dikonfigurasi oleh [@fbbdev](https://github.com/fbbdev) dalam [#4067](https://github.com/wailsapp/wails/pull/4067)
- Layanan `kvstore` secara default menyediakan penyimpanan nilai-kunci dalam memori ketika belum dikonfigurasi oleh [@fbbdev](https://github.com/fbbdev) dalam [#4067](https://github.com/wailsapp/wails/pull/4067)
- Menambahkan metode `Load` pada layanan `kvstore` untuk memuat ulang data dari berkas setelah konfigurasi berubah oleh [@fbbdev](https://github.com/fbbdev) dalam [#4067](https://github.com/wailsapp/wails/pull/4067)
- Menambahkan metode `Clear` pada layanan `kvstore` untuk menghapus semua kunci oleh [@fbbdev](https://github.com/fbbdev) dalam [#4067](https://github.com/wailsapp/wails/pull/4067)
- Menambahkan tipe `Level` dalam layanan `log` untuk menyediakan konstanta level log di sisi JS oleh [@fbbdev](https://github.com/fbbdev) dalam [#4067](https://github.com/wailsapp/wails/pull/4067)
- Menambahkan metode `Log` pada layanan `log` untuk menentukan level log secara dinamis oleh [@fbbdev](https://github.com/fbbdev) dalam [#4067](https://github.com/wailsapp/wails/pull/4067)
- Layanan `sqlite` secara default menyediakan DB dalam memori jika tidak dikonfigurasi, oleh [@fbbdev](https://github.com/fbbdev) dalam [#4067](https://github.com/wailsapp/wails/pull/4067)
- Tambahkan metode `Close` pada layanan `sqlite` untuk menutup DB secara manual, oleh [@fbbdev](https://github.com/fbbdev) dalam [#4067](https://github.com/wailsapp/wails/pull/4067)
- Tambahkan dukungan pembatalan untuk metode kueri pada layanan `sqlite`, oleh [@fbbdev](https://github.com/fbbdev) dalam [#4067](https://github.com/wailsapp/wails/pull/4067)
- Tambahkan dukungan prepared statement ke layanan `sqlite` dengan binding JS, oleh [@fbbdev](https://github.com/fbbdev) dalam [#4067](https://github.com/wailsapp/wails/pull/4067)
- Dukungan Gin oleh [Lea Anthony](https://github.com/leaanthony) dalam [PR](https://github.com/wailsapp/wails/pull/3537), berdasarkan karya awal [@AnalogJ](https://github.com/AnalogJ) dalam [PR](https://github.com/wailsapp/wails/pull/3537) ini
- Perbaiki penyimpanan otomatis dan penyimpanan otomatis kata sandi yang selalu diaktifkan, oleh [@oSethoum](https://github.com/osethoum) dalam [#4134](https://github.com/wailsapp/wails/pull/4134)
- Tambahkan `SetMenu()` pada jendela agar menu dapat ditetapkan pada jendela, oleh [@leaanthony](https://github.com/leaanthony)
- Tambahkan dukungan Notifikasi oleh [@popaprozac](https://github.com/popaprozac) dalam [#4098](https://github.com/wailsapp/wails/pull/4098)
-  Tambahkan dukungan Asosiasi File untuk Mac oleh [@wimaha](https://github.com/wimaha) dalam [#4177](https://github.com/wailsapp/wails/pull/4177)
- Tambahkan `wails3 tool version` untuk menaikkan versi semantik, oleh [@leaanthony](https://github.com/leaanthony)
- Tambahkan dukungan badge untuk macOS dan Windows, oleh [@popaprozac](https://github.com/popaprozac) dalam [#](https://github.com/wailsapp/wails/pull/4234)
- Tambahkan dukungan untuk event terdaftar/bertipe ketat, oleh [@fbbdev](https://github.com/fbbdev) dan [@IanVS](https://github.com/IanVS) dalam [#4161](https://github.com/wailsapp/wails/pull/4161)
- Tambahkan kemampuan untuk mendaftarkan hook bagi event khusus, oleh [@fbbdev](https://github.com/fbbdev) dan [@IanVS](https://github.com/IanVS) dalam [#4161](https://github.com/wailsapp/wails/pull/4161)
- `app.OpenFileManager(path string, selectFile bool)` untuk membuka pengelola file sistem pada jalur `path` dengan penyorotan opsional melalui `selectFile`, oleh [@Krzysztofz01](https://github.com/Krzysztofz01) [@rcalixte](https://github.com/rcalixte)
- Flag `-git` baru untuk perintah `wails3 init`, oleh [@leaanthony](https://github.com/leaanthony)
- Perintah `wails3 generate webview2bootstrapper` baru oleh [@leaanthony](https://github.com/leaanthony)
- Tambahkan metode `init()` dalam runtime agar runtime dapat diinisialisasi secara manual, oleh [@leaanthony](https://github.com/leaanthony)
- Tambahkan opsi `WindowDidMoveDebounceMS` ke WindowOptions milik Window, oleh [@leaanthony](https://github.com/leaanthony)
- Tambahkan fitur Instans Tunggal oleh [@leaanthony](https://github.com/leaanthony). Berdasarkan [PR v2](https://github.com/wailsapp/wails/pull/2951) oleh @APshenkin.
- Perintah `wails3 generate template` oleh [@leaanthony](https://github.com/leaanthony)
- Perintah `wails3 releasenotes` oleh [@leaanthony](https://github.com/leaanthony)
- Perintah `wails3 update cli` oleh [@leaanthony](https://github.com/leaanthony)
- Opsi `-clean` untuk perintah `wails3 generate bindings`, oleh [@leaanthony](https://github.com/leaanthony)
- Izinkan build AppImage Linux aarch64 (arm64), oleh [@AkshayKalose](https://github.com/AkshayKalose) dalam [#3981](https://github.com/wailsapp/wails/pull/3981)
- Tambahkan pranala untuk sponsor oleh @ansxuman dalam [#3958](https://github.com/wailsapp/wails/pull/3958)
- Dukungan pemaketan Linux untuk build paket deb, rpm, dan Arch Linux oleh
- Tambahkan dukungan untuk build dan paket universal Darwin oleh
- Dokumentasi event pada situs web oleh
- Templat untuk sveltekit dan sveltekit-ts yang dikonfigurasi untuk pengembangan non-SSR
- Perbarui aset build menggunakan perintah `wails3 update build-assets` baru oleh
- Contoh untuk menguji HTML Drag and Drop API oleh
- Dukungan Asosiasi File oleh [leaanthony](https://github.com/leaanthony) dalam
- Perintah `wails3 generate runtime` baru oleh
- Opsi `InitialPosition` baru untuk menentukan apakah jendela harus dipusatkan atau
- Tambahkan metode `Path` & `Paths` ke paket `application` oleh
- Tambahkan opsi Windows `GeneralAutofillEnabled` dan `PasswordAutosaveEnabled`
- Tambahkan kemampuan untuk mengambil jendela yang memanggil metode layanan, oleh
- Tambahkan opsi `EnabledFeatures` dan `DisabledFeatures` untuk WebView2, oleh
- ⊞ Sistem DIP baru untuk Dukungan Monitor DPI Tinggi yang Ditingkatkan, oleh
- ⊞ Opsi nama kelas jendela oleh [windom](https://github.com/windom/) dalam
- Layanan telah diperluas untuk menyediakan fungsionalitas plugin. Oleh
- 🐧 Event WindowDidMove / WindowDidResize dalam
- ⊞ Event WindowDidResize dalam
-  tambahkan event ApplicationShouldHandleReopen agar dock dapat ditangani
-  tambahkan getPrimaryScreen/getScreens ke implementasi oleh @tmclane dalam
-  tambahkan opsi untuk menampilkan bilah alat dalam mode layar penuh di macOS oleh
- 🐧 tambahkan logika onKeyPress untuk mengonversi penekanan tombol Linux menjadi akselerator
- 🐧 tambahkan task `run:linux` oleh
- Ekspor metode `SetIcon` oleh [@almas-x](https://github.com/almas-x) dalam
- Tingkatkan `OnShutdown` oleh [@almas-x](https://github.com/almas-x) dalam
- Pulihkan metode `ToggleMaximise` dalam antarmuka `Window` oleh
- Tambahkan informasi lebih lanjut ke `Environment()`. Oleh @leaanthony dalam
- Ekspos metode `WebviewWindow.IsFocused` pada antarmuka `Window` oleh
- Dukung beberapa event pemicu yang dipisahkan spasi dalam sistem WML oleh
- Tambahkan ekspor ESM dari skrip runtime JS yang dibundel oleh
- Tambahkan flag generator binding untuk menggunakan skrip runtime JS yang dibundel sebagai pengganti
- Implementasikan `setIcon` di Linux oleh [@abichinger](https://github.com/abichinger)
- Tambahkan flag `-port` ke perintah dev dan dukung variabel lingkungan
- Tambahkan pengujian untuk pemanggilan metode terikat oleh
- ⊞ tambahkan `SetIgnoreMouseEvents` untuk jendela yang telah dibuat oleh
-  Tambahkan kemampuan untuk mengatur tingkat tumpukan (urutan) jendela oleh

### Diperbaiki

- Perbaiki `Screen.Bounds`, `WorkArea`, dan `Size` yang menjadi setengahnya pada Mac Retina dengan mengonversi nilai titik NSScreen menjadi piksel perangkat dalam bidang `Physical*`, serta isi `Screen.X`/`Y` tingkat teratas agar deteksi monitor yang saling bersentuhan dan penempatan area kerja sudah benar, dalam [PR](https://github.com/wailsapp/wails/pull/5168) oleh @wayneforrest
- Perbaiki data race di ScreenManager yang menyebabkan deadlock WebKit DisplayLink saat konfigurasi tampilan berubah (misalnya monitor eksternal dipasang saat proses tidur/bangun)
- Menetapkan CFBundleIconName secara langsung ke appicon ketika Assets.car ada dalam [PR](https://github.com/wailsapp/wails/pull/5154) oleh @symball
- Memperbaiki `wails3 doctor` yang melaporkan paket WebKitGTK yang salah pada Fedora, openSUSE, Arch, dan NixOS — entri fallback 4.0 telah dihapus karena v3 memerlukan API 4.1 saat kompilasi (#5071)
- Memperbaiki nama paket webkit2gtk untuk pemeriksaan doctor di openSUSE (`webkit2gtk4_1-devel` → `webkit2gtk3-devel`, nama paket openSUSE yang benar) (#5071)
- Memperbaiki galat `Unexpected token '<'` ketika `/wails/custom.js` tidak ada dalam mode pengembangan desktop. Menambahkan penangan 404 eksplisit untuk `/wails/custom.js` dan validasi `Content-Type` tanpa membedakan huruf besar-kecil dalam `loadOptionalScript` untuk mencegah fallback SPA HTML disisipkan sebagai JavaScript. ([#5068](https://github.com/wailsapp/wails/issues/5068))
- Memperbaiki status sorotan menu baki sistem pada macOS—ikon kini menampilkan status terpilih ketika menu terbuka (#4910)
- Memperbaiki jendela yang terpasang pada baki sistem dan muncul di belakang jendela lain pada macOS—kini menggunakan tingkat jendela popup yang tepat (#4910)
- Memperbaiki contoh impor `@wailsio/runtime` yang salah di seluruh dokumentasi (#4989)
- Memperbaiki jendela tanpa bingkai yang tidak dapat diminimalkan pada darwin (#4294)
- Memperbaiki macet selama 20-30 menit saat `wails3 build` dan `wails3 dev` dengan mengecualikan `node_modules/` dari pemeriksaan kemutakhiran go-task. Sebelumnya, glob `sources: "**/*"` menyebabkan go-task menghitung dan membuat checksum untuk setiap berkas dalam `node_modules/` (lebih dari 50000-100000 berkas dengan dependensi berat seperti MUI), yang sangat lambat terutama pada Windows/NTFS (#4939)
- Memperbaiki kegagalan build GTK4 akibat typedef `Screen` C yang bertabrakan dengan X11 Xlib.h (#4957)
- Memperbaiki konsistensi metode lencana dock pada macOS
- Memperbaiki `InvisibleTitleBarHeight` yang diterapkan pada semua jendela macOS, bukan hanya jendela tanpa bingkai atau dengan bilah judul transparan (#4960)
- Memperbaiki jendela yang berguncang/bergetar saat ukurannya diubah dari sudut atas ketika `InvisibleTitleBarHeight` diaktifkan, dengan melewati inisiasi seret di dekat tepi jendela (#4960)
- Memperbaiki pembuatan tipe terpetakan dengan kunci enum dalam binding JS/TS (#4437) oleh @fbbdev
- Memperbaiki seret dan lepas berkas pada Windows yang tidak berfungsi pada skala tampilan selain 100%
- Memperbaiki seret dan lepas internal HTML5 yang rusak ketika pelepasan berkas diaktifkan pada Windows
- Memperbaiki koordinat pelepasan berkas yang menggunakan ruang piksel yang salah pada Windows (piksel fisik dibandingkan piksel CSS)
- Memperbaiki seret dan lepas berkas pada Linux yang tidak berfungsi secara andal dengan efek hover
- Memperbaiki seret dan lepas internal HTML5 yang rusak ketika pelepasan berkas diaktifkan pada Linux
- Memperbaiki tindakan menampilkan/menyembunyikan jendela pada Linux/GTK4 yang terkadang memulihkan jendela ke status diminimalkan dengan menggunakan `gtk_window_present()` (#4957)
- Memperbaiki pengambilan/penetapan posisi jendela pada Linux/GTK4 yang selalu mengembalikan 0,0 dengan menambahkan dukungan bersyarat X11 melalui `XTranslateCoordinates`/`XMoveWindow` (#4957)
- Memperbaiki ukuran maksimum jendela yang tidak diberlakukan pada Linux/GTK4 dengan menambahkan pembatasan ukuran berbasis sinyal untuk menggantikan `gtk_window_set_geometry_hints` yang telah dihapus (#4957)
- Memperbaiki penskalaan DPI pada Linux/GTK4 dengan menerapkan penghitungan PhysicalBounds yang tepat dan dukungan penskalaan pecahan melalui `gdk_monitor_get_scale` (GTK 4.14+)
- Memperbaiki item menu yang terduplikasi saat membuat jendela baru pada Linux/GTK4
- Memperbaiki pembuatan tipe terpetakan dengan kunci enum dalam binding JS/TS (#4437) oleh @fbbdev
- Memperbaiki seret dan lepas berkas pada Windows yang tidak berfungsi pada skala tampilan selain 100%
- Memperbaiki seret dan lepas internal HTML5 yang rusak ketika pelepasan berkas diaktifkan pada Windows
- Memperbaiki koordinat pelepasan berkas yang menggunakan ruang piksel yang salah pada Windows (piksel fisik dibandingkan piksel CSS)
- Memperbaiki seret dan lepas berkas pada Linux yang tidak berfungsi secara andal dengan efek hover
- Memperbaiki seret dan lepas internal HTML5 yang rusak ketika pelepasan berkas diaktifkan pada Linux
- Memperbaiki penskalaan DPI pada Linux/GTK4 dengan menerapkan penghitungan PhysicalBounds yang tepat dan dukungan penskalaan pecahan melalui `gdk_monitor_get_scale` (GTK 4.14+)
- Memperbaiki item menu yang terduplikasi saat membuat jendela baru pada Linux/GTK4
- Memperbaiki pembuatan tipe terpetakan dengan kunci enum dalam binding JS/TS (#4437) oleh @fbbdev
- Memperbaiki masalah "jendela hantu" pada macOS akibat API AppKit tidak diakses dari Main Thread dalam App.Window.Current() (#4947) oleh @wimaha
- Memperbaiki `<input type="file">` HTML yang tidak berfungsi pada macOS dengan menerapkan WKUIDelegate runOpenPanelWithParameters (#4862)
- Memperbaiki seret dan lepas berkas native yang tidak berfungsi saat menggunakan modul npm `@wailsio/runtime` pada macOS/Linux (#4953) oleh @leaanthony
- Memperbaiki pembuatan binding untuk alias tipe lintas paket (#4578) oleh @fbbdev
- Memperbaiki crash OpenFileDialog pada Linux akibat pelanggaran keamanan thread GTK (#3683) oleh @ddmoney420
- Memperbaiki crash SIGSEGV saat memanggil `Focus()` pada jendela yang disembunyikan atau telah dihancurkan (#4890) oleh @ddmoney420
- Memperbaiki potensi panic saat menetapkan ikon atau bitmap kosong pada Linux (#4923) oleh @ddmoney420
- Memperbaiki crash ErrorDialog ketika dipanggil dari binding layanan pada macOS (#3631) oleh @leaanthony
- Membuat menu ditampilkan pada sistem operasi Windows dalam `v3\examples\dialogs` oleh @ndianabasi
- Memperbaiki kondisi balapan yang menyebabkan TypeError saat halaman dimuat ulang (#4872) oleh @ddmoney420
- Memperbaiki keluaran yang salah dari pengujian generator binding dengan menghapus status global dalam metode `Collector.IsVoidAlias()` (#4941) oleh @fbbdev
- Memperbaiki pemilih berkas `<input type="file">` yang tidak berfungsi pada macOS (#4862) oleh @leaanthony
- Memperbaiki `Position()` dan `SetPosition()` yang menggunakan sistem koordinat tidak konsisten pada macOS sehingga posisi jendela bergeser saat status disimpan/dipulihkan (#4816) oleh @leaanthony
- Memperbaiki galat SetProcessDpiAwarenessContext "Access is denied" ketika kesadaran DPI telah ditetapkan melalui manifes aplikasi (#4803)
- Memperbarui halaman dokumentasi pintasan papan ketik dan memperbaiki tipe parameter callback untuk `KeyBinding.Add` oleh @ndianabasi
- Memperbaiki dokumentasi tentang pembuatan binding khusus yang harus menggunakan `-d String` alih-alih `-o String`
- Memperbaiki menu yang tidak menghapus elemen turunannya pada `menu.Update()`
- Memperbaiki referensi Manager API yang usang dalam dokumentasi (berkas 31 diperbarui agar menggunakan pola baru seperti `app.Window.New()`, `app.Event.Emit()`, dan sebagainya) oleh @leaanthony
- Memperbaiki crash pada Linux akibat panic dalam metode Go yang terikat ke JS karena WebKit menimpa penangan sinyal (#3965) oleh @leaanthony
- Memperbaiki SaveFileDialog.SetFilename() yang tidak berpengaruh pada Linux (#4841) oleh @samstanier
- Memperbaiki koordinat pelepasan yang ditampilkan sebagai undefined dalam contoh seret dan lepas
- Memperbaiki kegagalan pembuatan bundel aplikasi macOS ketika APP_NAME mengandung spasi (masalah brace expansion)
- Memperbaiki panic indeks di luar batas pada Windows saat memanggil metode layanan (mengembalikan goccy/go-json)
- Memperbaiki seret dan lepas berkas pada Windows yang tidak berfungsi pada skala tampilan selain 100%
- Memperbaiki seret dan lepas internal HTML5 yang rusak ketika pelepasan berkas diaktifkan pada Windows
- Memperbaiki koordinat pelepasan berkas yang menggunakan ruang piksel yang salah pada Windows (piksel fisik dibandingkan piksel CSS)
- Memperbaiki seret-dan-lepas berkas di Linux yang tidak berfungsi dengan andal saat menggunakan efek hover
- Memperbaiki seret-dan-lepas internal HTML5 yang rusak ketika fitur pelepasan berkas diaktifkan di Linux
- Memperbarui semua perintah dalam berkas Taskfile.yml untuk semua sistem operasi agar mendukung spasi dalam variabel seperti `APP_NAME` oleh @ndianabasi
- Memperbaiki kesalahan argumen perintah saat menjalankan tugas 'build:universal:lipo:go' di Linux oleh @wux1an
- Memperbaiki kesalahan Docker "undefined symbol: **<em>ubsan</em>handle_xxxxxxx" saat menjalankan 'wails3 build GOOS=darwin GOARCH=arm64' di Linux oleh @wux1an
- Menyatukan dokumentasi protokol khusus dan menambahkan bagian Universal Links oleh @leaanthony
- Memperbaiki crash menu baki sistem Windows saat ikon diklik berulang kali dengan menambahkan pengaman terhadap pemanggilan TrackPopupMenuEx secara bersamaan (#4151) oleh @leaanthony
- Mencegah aplikasi mengalami crash ketika systray.Run() dipanggil sebelum app.Run() oleh @leaanthony
- Memperbaiki crash di macOS saat visibilitas jendela dialihkan melalui Hide()/Show() ketika ApplicationShouldTerminateAfterLastWindowClosed diaktifkan (#4389) oleh @leaanthony
- Memperbaiki kebocoran memori pada menu konteks di macOS dan Windows ketika menu dibuka berulang kali (#4012) oleh @leaanthony
- Memperbaiki sumber daya native menu konteks yang tidak digunakan kembali di macOS sehingga menu baru dibuat setiap kali ditampilkan (#4012) oleh @leaanthony
- Memperbaiki klik ikon Dock macOS yang tidak menampilkan jendela tersembunyi ketika aplikasi dimulai dengan `Hidden: true` (#4583) oleh @leaanthony
- Memperbaiki dialog cetak macOS yang tidak terbuka akibat tipe penunjuk jendela yang salah dalam pemanggilan CGO (#4290) oleh @leaanthony
- Memperbaiki crash menu jendela di Wayland yang disebabkan oleh appmenu-gtk-module saat mengakses jendela yang belum direalisasikan (#4769) oleh @leaanthony
- Memperbaiki crash aplikasi GTK ketika nama aplikasi berisi karakter yang tidak valid (spasi, tanda kurung, dan sebagainya) oleh @leaanthony
- Memperbaiki kesalahan "not enough memory" saat menginisialisasi seret-dan-lepas di Windows (#4701) oleh @overlordtm
- Memperbaiki penjelajah berkas yang membuka direktori yang salah di Linux akibat pengodean escape URI yang keliru (#4397) oleh @leaanthony
- Memperbaiki kegagalan build AppImage pada distribusi Linux modern (Arch, Fedora 39+, Ubuntu 24.04+) dengan mendeteksi bagian ELF `.relr.dyn` secara otomatis dan menonaktifkan stripping (#4642) oleh @leaanthony
- Memperbaiki `wails doctor` yang keliru melaporkan paket webkit sebagai telah terinstal pada Fedora/sistem berbasis DNF (#4457) oleh @leaanthony
- Memperbaiki `config.yml` bawaan yang menjalankan `wails3 dev` dengan build produksi oleh @mbaklor
- Memperbaiki stub layanan iOS yang menyebabkan kegagalan build akibat impor paket yang tidak ada oleh @leaanthony
- Memperbaiki logging terstruktur dalam metode debug/info yang menyebabkan kesalahan "no formatting directives" oleh @leaanthony
- Menghapus pernyataan cetak debug sementara yang tidak sengaja disertakan dari penggabungan platform seluler oleh @leaanthony
- Memperbaiki crash WebKitGTK di Wayland dengan GPU NVIDIA (Kesalahan protokol 71) dengan menonaktifkan perender DMA-BUF secara otomatis oleh @leaanthony
- Memperbaiki nilai alfa yang diabaikan dalam `application.WebviewWindowOptions.BackgroundColour` di Linux ([#4722](https://github.com/wailsapp/wails/pull/4722), @BradHacker)
- Memperbaiki ikon baki sistem Windows yang tidak menggunakan ikon aplikasi secara bawaan ketika ikon khusus tidak disediakan (#4704)
- Melacak kepemilikan `HICON` agar hanya handle buatan pengguna yang dihancurkan, sehingga mencegah crash saat Explorer didaur ulang (#4653).
- Melepas listener tema sistem Windows dan ikon baki yang dipertahankan saat objek dihancurkan untuk menghentikan kebocoran goroutine dan konteks perangkat (#4653).
- Memotong tooltip baki pada 127 unit UTF-16 untuk menghindari kerusakan pasangan surrogate dan glif multibita (#4653).
- Memperbaiki kegagalan tugas paket Windows (#4667)
- Memperbaiki variabel appicon AppImage Linux dalam taskfile Linux [PR #4644](https://github.com/wailsapp/wails/pull/4644)
- Memperbaiki kesalahan build Windows akibat perubahan signature go-webview2 v1.0.22 (#4513, #4645)
- Memperbaiki variabel appicon AppImage Linux dalam taskfile Linux [PR #4644](https://github.com/wailsapp/wails/pull/4644)
- Memperbaiki rentang protokol desktop.tmpl Linux dengan menghapus `<.Info.Protocol>` hingga `<.Protocol>` oleh @Tolfx dalam #4510
- Memperbaiki kesalahan redefinisi untuk demo liquid glass dalam [#4542](https://github.com/wailsapp/wails/pull/4542) oleh @Etesam913
- Memperbaiki pembaruan menu baki sistem di Linux [#4604](https://github.com/wailsapp/wails/issues/4604) oleh [@JackDoan](https://github.com/JackDoan)
- Memperbaiki jendela putih yang muncul di Windows saat membuat jendela tersembunyi oleh @leaanthony dalam [#4612](https://github.com/wailsapp/wails/pull/4612)
- Memperbaiki jalur impor paket notifikasi dalam dokumentasi oleh @rxliuli dalam [#4617](https://github.com/wailsapp/wails/pull/4617)
- Memperbaiki seret-dan-lepas yang tidak berfungsi saat menggunakan paket npm @wailsio/runtime (#4489) oleh @leaanthony dalam #4616
- Windows: Memperbaiki jendela yang berkedip saat aplikasi dimulai dan jendela tersembunyi yang ditampilkan secara keliru dalam [PR](https://github.com/wailsapp/wails/pull/4600) oleh @leaanthony.
- Memperbaiki masalah ukuran jendela saat dimaksimalkan di Wayland (https://github.com/wailsapp/wails/issues/4429) oleh [@samstanier](https://github.com/samstanier)
- Memperbaiki masalah ukuran jendela saat dimaksimalkan di Wayland (https://github.com/wailsapp/wails/issues/4429) oleh [@samstanier](https://github.com/samstanier)
- Memperbaiki kesalahan redefinisi untuk demo liquid glass dalam [#4542](https://github.com/wailsapp/wails/pull/4542) oleh @Etesam913
- Memperbaiki masalah yang dapat menyebabkan AssetServer mengalami crash di macOS dalam [#4576](https://github.com/wailsapp/wails/pull/4576) oleh @jghiloni
- Memperbaiki masalah kompilasi saat membangun dengan NextJs. Diperbaiki dalam [#4585](https://github.com/wailsapp/wails/pull/4585) oleh @rev42
- Memperbaiki pipeline untuk rilis nightly dalam [#4597](https://github.com/wailsapp/wails/pull/4597) oleh @riadafridishibly
- Memperbaiki kesalahan redefinisi untuk demo liquid glass dalam [#4542](https://github.com/wailsapp/wails/pull/4542) oleh @Etesam913
- Memperbaiki masalah yang dapat menyebabkan AssetServer mengalami crash di macOS dalam [#4576](https://github.com/wailsapp/wails/pull/4576) oleh @jghiloni
- Memperbaiki masalah kompilasi saat membangun dengan NextJs. Diperbaiki dalam [#4585](https://github.com/wailsapp/wails/pull/4585) oleh @rev42
- Memperbaiki pipeline untuk rilis nightly dalam [#4597](https://github.com/wailsapp/wails/pull/4597) oleh @riadafridishibly
- Memperbaiki kesalahan redefinisi untuk demo liquid glass dalam [#4542](https://github.com/wailsapp/wails/pull/4542) oleh @Etesam913
- Memperbaiki SetBackgroundColour di Windows oleh @PPTGamer dalam [PR](https://github.com/wailsapp/wails/pull/4492)
- Memperbarui dokumentasi agar mencerminkan perubahan dari refaktorisasi Manager API oleh @yulesxoxo dalam [PR #4476](https://github.com/wailsapp/wails/pull/4476)
- Memperbaiki variabel appicon berkas .desktop Linux dalam taskfile Linux [PR #4477](https://github.com/wailsapp/wails/pull/4477)
- Memperbarui dokumentasi agar mencerminkan perubahan dari refaktorisasi Manager API oleh @yulesxoxo dalam [PR #4476](https://github.com/wailsapp/wails/pull/4476)
- Memperbaiki bug dereferensi penunjuk nil di Windows yang dilaporkan dalam [#4456](https://github.com/wailsapp/wails/issues/4456) oleh @leaanthony dalam [#4460](https://github.com/wailsapp/wails/pull/4460)
- Menambahkan dukungan untuk `allowsBackForwardNavigationGestures` di WKWebView macOS guna mengaktifkan gestur navigasi geser dua jari (#1857)
- Memperbaiki masalah onClick yang tidak berfungsi untuk item menu yang awalnya ditetapkan sebagai dinonaktifkan oleh @leaanthony dalam [PR #4469](https://github.com/wailsapp/wails/pull/4469). Terima kasih kepada @IanVS atas investigasi awalnya.
- Memperbaiki server Vite yang tidak dibersihkan saat build gagal (#4403)
- Memperbaiki panic saat menutup atau membatalkan `SaveFileDialog` di Windows. Diperbaiki dalam [PR](https://github.com/wailsapp/wails/pull/4284) oleh @hkhere
- Memperbaiki seret dan lepas pada tingkat HTML di Windows oleh [@mbaklor](https://github.com/mbaklor) dalam [#4259](https://github.com/wailsapp/wails/pull/4259)
- Menambahkan dukungan untuk `allowsBackForwardNavigationGestures` di WKWebView macOS guna mengaktifkan gestur navigasi usap dua jari (#1857)
- Memperbaiki masalah onClick yang tidak berfungsi untuk item menu yang awalnya ditetapkan sebagai dinonaktifkan oleh @leaanthony dalam [PR #4469](https://github.com/wailsapp/wails/pull/4469). Terima kasih kepada @IanVS atas investigasi awalnya.
- Memperbaiki server Vite yang tidak dibersihkan saat build gagal (#4403)
- Memperbaiki penguraian notifikasi di Windows oleh @popaprozac dalam [PR](https://github.com/wailsapp/wails/pull/4450)
- Memperbaiki perintah doctor agar memeriksa dependensi Windows SDK oleh [@kodumulo](https://github.com/kodumulo) dalam [#4390](https://github.com/wailsapp/wails/issues/4390)
- Memperbaiki dereferensi pointer nil di processURLRequest untuk Mac oleh [@etesam913](https://github.com/etesam913) dalam [#4366](https://github.com/wailsapp/wails/pull/4366)
- Memperbaiki bug Linux yang menghalangi dialog berfilter oleh [@bh90210](https://github.com/bh90210) dalam [#4287](https://github.com/wailsapp/wails/pull/4287)
- Memperbaiki masalah menu Edit di Windows dan Linux oleh [@leaanthony](https://github.com/leaanthony) dalam [#3f78a3a](https://github.com/wailsapp/wails/commit/3f78a3a8ce7837e8b32242c8edbbed431c68c062)
- Memperbarui versi sistem minimum dalam file .plist macOS dari 10.13.0 menjadi 10.15.0 oleh [@AkshayKalose](https://github.com/AkshayKalose) dalam [#3981](https://github.com/wailsapp/wails/pull/3981)
- Masalah ID jendela yang terlewati oleh [@leaanthony](https://github.com/leaanthony)
- Memperbaiki masalah menu nil saat memanggil RegisterContextMenu oleh [@leaanthony](https://github.com/leaanthony)
- Memperbaiki siklus dependensi dalam keluaran generator binding oleh [@fbbdev](https://github.com/fbbdev) dalam [#4001](https://github.com/wailsapp/wails/pull/4001)
- Memperbaiki galat penggunaan sebelum definisi dalam keluaran generator binding oleh [@fbbdev](https://github.com/fbbdev) dalam [#4001](https://github.com/wailsapp/wails/pull/4001)
- Meneruskan flag build ke generator binding oleh [@fbbdev](https://github.com/fbbdev) dalam [#4023](https://github.com/wailsapp/wails/pull/4023)
- Mengubah jalur dalam Taskfile Windows agar menggunakan garis miring ke depan untuk memastikan Taskfile berfungsi pada platform selain Windows oleh [@leaanthony](https://github.com/leaanthony)
- Event Mac dan JS Mac kini telah diperbaiki oleh [@leaanthony](https://github.com/leaanthony)
- Memperbaiki deadlock event untuk macOS oleh [@leaanthony](https://github.com/leaanthony)
- Memperbaiki galat `Parameter incorrect` dalam inisialisasi Window di Windows ketika HTML disediakan tanpa JS oleh [@leaanthony](https://github.com/leaanthony)
- Memperbaiki ukuran prefiks respons yang digunakan untuk mendeteksi jenis konten di server aset oleh [@fbbdev](https://github.com/fbbdev) dalam [#4049](https://github.com/wailsapp/wails/pull/4049)
- Memperbaiki penanganan respons non-404 pada jalur indeks root di server aset oleh [@fbbdev](https://github.com/fbbdev) dalam [#4049](https://github.com/wailsapp/wails/pull/4049)
- Memperbaiki perilaku tak terdefinisi dalam generator binding saat menguji properti tipe generik oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Memperbaiki keluaran generator binding untuk model ketika tipe dasarnya tidak memiliki properti yang sama dengan pembungkus bernama oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Memperbaiki keluaran generator binding untuk tipe kunci map dan prapemrosesan oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Memperbaiki keluaran generator binding untuk struct yang mengimplementasikan antarmuka marshaler oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Memperbaiki deteksi siklus tipe yang melibatkan tipe generik dalam generator binding oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Memperbaiki referensi yang tidak valid ke model yang tidak diekspor dalam keluaran generator binding oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Memindahkan kode yang diinjeksi ke bagian akhir file layanan oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Memperbaiki penanganan galat dari operasi penutupan file dalam generator binding oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Menonaktifkan peringatan untuk layanan yang mendefinisikan metode siklus hidup atau HTTP, tetapi tidak memiliki metode terikat lainnya, oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Memperbaiki templat non-React yang gagal menampilkan footer Hello World saat menggunakan skema warna sistem terang oleh [@marcus-crane](https://github.com/marcus-crane) dalam [#4056](https://github.com/wailsapp/wails/pull/4056)
- Memperbaiki item menu tersembunyi di macOS oleh [@leaanthony](https://github.com/leaanthony)
- Memperbaiki penanganan dan pemformatan galat dalam pemroses pesan oleh [@fbbdev](https://github.com/fbbdev) dalam [#4066](https://github.com/wailsapp/wails/pull/4066)
-  Memperbaiki penghentian layanan yang terlewati saat keluar dari aplikasi oleh [@fbbdev](https://github.com/fbbdev) dalam [#4066](https://github.com/wailsapp/wails/pull/4066)
-  Memastikan pembaruan menu terjadi di thread utama oleh [@leaanthony](https://github.com/leaanthony)
- Mekanisme penyeretan dan pengubahan ukuran kini lebih andal serta lebih sesuai dengan perilaku platform yang diharapkan oleh [@fbbdev](https://github.com/fbbdev) dalam [#4100](https://github.com/wailsapp/wails/pull/4100)
- Memperbaiki [#4097](https://github.com/wailsapp/wails/issues/4097): Webpack/Angular membuang kode inisialisasi runtime, oleh [@fbbdev](https://github.com/fbbdev) dalam [#4100](https://github.com/wailsapp/wails/pull/4100)
- Memperbaiki item menu yang awalnya disembunyikan oleh [@IanVS](https://github.com/IanVS) dalam [#4116](https://github.com/wailsapp/wails/pull/4116)
- Memperbaiki assetFileServer yang tidak menyajikan file `.html` untuk permintaan tanpa ekstensi ketika `[request]` tidak ada tetapi `[request].html` ada
- Memperbaiki jalur pembuatan ikon oleh [@robin-samuel](https://github.com/robin-samuel) dalam [#4125](https://github.com/wailsapp/wails/pull/4125)
- Memperbaiki event `fullscreen`, `unfullscreen`, `unminimise`, dan `unmaximise` yang tidak dipancarkan oleh [@oSethoum](https://github.com/osethoum) dalam [#4130](https://github.com/wailsapp/wails/pull/4130)
- Memperbaiki galat NSIS akibat prefiks yang salah pada versi default dalam konfigurasi oleh [@robin-samuel](https://github.com/robin-samuel) dalam [#4126](https://github.com/wailsapp/wails/pull/4126)
- Memperbaiki fungsi runtime Dialogs yang mengembalikan jalur yang di-escape di Windows oleh [TheGB0077](https://github.com/TheGB0077) dalam [#4188](https://github.com/wailsapp/wails/pull/4188)
- Memperbaiki jalur deteksi Webview2 di HKCU oleh [@leaanthony](https://github.com/leaanthony).
- Memperbaiki masalah input di macOS oleh [@leaanthony](https://github.com/leaanthony).
- Memperbaiki nama file tugas pembuatan ikon Windows oleh [@yulesxoxo](https://github.com/yulesxoxo) di [#4219](https://github.com/wailsapp/wails/pull/4219).
- Memperbaiki masalah transparansi pada jendela tanpa bingkai oleh [@leaanthony](https://github.com/leaanthony) berdasarkan karya @kron.
- Memperbaiki pemanggilan fokus saat jendela dinonaktifkan atau diminimalkan oleh [@leaanthony](https://github.com/leaanthony) berdasarkan karya @kron.
- Memperbaiki baki sistem yang tidak muncul setelah bilah tugas dimulai ulang oleh [@leaanthony](https://github.com/leaanthony) berdasarkan karya @kron.
- Memperbaiki fallbackResponseWriter yang tidak mengimplementasikan Flush() di [#4245](https://github.com/wailsapp/wails/pull/4245)
- Memperbaiki fallbackResponseWriter yang tidak mengimplementasikan Flush() oleh [@superDingda] di [#4236](https://github.com/wailsapp/wails/issues/4236)
- Memperbaiki crash saat jendela macOS ditutup ketika pemanggilan asinkron fungsi yang terikat ke Go masih tertunda oleh [@joshhardy](https://github.com/joshhardy) di [#4354](https://github.com/wailsapp/wails/pull/4354)
- Memperbaiki kondisi berpacu saat memulai mode Efisiensi Windows oleh [@leaanthony](https://github.com/leaanthony)
- Memperbaiki pembersihan handle ikon Windows oleh [@leaanthony](https://github.com/leaanthony).
- Memperbaiki `OpenFileManager` pada Windows oleh [@PPTGamer](https://github.com/PPTGamer) di [#4375](https://github.com/wailsapp/wails/pull/4375).
- Memperbaiki opsi lebar minimum/maksimum untuk Linux oleh @atterpac di [#3979](https://github.com/wailsapp/wails/pull/3979)
- Definisi tipe templat TypeScript melalui peningkatan versi npm oleh @atterpac di [#3966](https://github.com/wailsapp/wails/pull/3966)
- Memperbaiki referensi CSS templat SvelteKit oleh @atterpac di [#3945](https://github.com/wailsapp/wails/pull/3945)
- Memastikan callback utama dalam run() jendela dipanggil di thread utama oleh [@leaanthony](https://github.com/leaanthony)
- Memperbaiki contoh pemilih direktori dialog oleh [@leaanthony](https://github.com/leaanthony)
- Membuat halaman galat baru berbahasa Tionghoa saat index.html tidak ada oleh [@leaanthony](https://github.com/leaanthony)
-  Memastikan callback `windowDidBecomeKey` berjalan di thread utama oleh [@leaanthony](https://github.com/leaanthony)
-  Mendukung layar penuh untuk jendela tanpa bingkai oleh [@leaanthony](https://github.com/leaanthony)
-  Meningkatkan logika penghancuran jendela oleh [@leaanthony](https://github.com/leaanthony)
-  Memperbaiki logika posisi jendela saat ditautkan ke baki sistem oleh [@leaanthony](https://github.com/leaanthony)
-  Mendukung layar penuh untuk jendela tanpa bingkai oleh [@leaanthony](https://github.com/leaanthony)
- Memperbaiki penanganan peristiwa oleh [@leaanthony](https://github.com/leaanthony)
- Memperbaiki logika penghentian jendela oleh [@leaanthony](https://github.com/leaanthony)
- Taskfile umum kini secara default menghasilkan binding TypeScript untuk templat TypeScript oleh [@leaanthony](https://github.com/leaanthony)
- Memperbaiki penutupan aplikasi saat menerima pesan WM_CLOSE ketika tidak ada jendela yang terbuka/hanya ada baki sistem oleh [@mmalcek](https://github.com/mmalcek) di [#3990](https://github.com/wailsapp/wails/pull/3990)
- Memperbaiki build garble oleh @5aaee9 di [#3192](https://github.com/wailsapp/wails/pull/3192)
- Memperbaiki build NSIS Windows oleh [@leaanthony](https://github.com/leaanthony)
- Memperbaiki deadlock pada dialog Linux untuk beberapa pilihan yang disebabkan oleh bagian yang tidak ditutup
- Memperbaiki pembersihan lintas platform untuk file .syso selama build Windows oleh
- Memperbaiki kompilasi AppImage amd64 oleh @atterpac di
- Memperbaiki pembaruan aset build oleh @ansxuman di
- Memperbaiki implementasi `OnClick` dan `OnRightClick` pada baki sistem Linux oleh @atterpac
- Memperbaiki `AlwaysOnTop` yang tidak berfungsi di Mac oleh
-  Memperbaiki `application.NewEditMenu` yang menyertakan duplikat
- 🐧 Memperbaiki kompilasi aarch64
- ⊞ Memperbaiki item menu grup tombol radio oleh
- Memperbaiki galat saat membangun .app yang dapat dijalankan di macOS ketika 'name' dan 'outputfilename'
- Memperbaiki bug penggunaan customEventProcessor dalam contoh seret dan lepas oleh
- 🐧 Memperbaiki galat kompilasi Linux yang muncul akibat penambahan IgnoreMouseEvents oleh
- ⊞ Memperbaiki bug pembuatan file ikon syso oleh
- 🐧 Perbaikan agar berjalan secara native di Wayland yang diambil dari
- Jangan ikat metode layanan internal di
- ⊞ Memperbaiki panic saat baki sistem dimulai di
- Jangan ikat metode layanan internal di
- ⊞ Memperbaiki panic saat baki sistem dimulai di
- Refaktor besar pada item menu dan penanganan peristiwa. Untuk saat ini, terutama meningkatkan macOS. Oleh
- Memperbaiki pengujian setelah refaktor plugin dan peristiwa di
- ⊞ Memperbaiki peringatan `Failed to unregister class Chrome_WidgetWin_0`. Oleh
- Masalah modul
- Memperbaiki pengiriman pesan peristiwa perubahan ukuran oleh [atterpac](https://github.com/atterpac) di
- 🐧 Memperbaiki galat penanganan tema di NixOS oleh
- Memperbaiki instalasi proyek lintas volume untuk Windows oleh
- Memperbaiki CSS templat React agar footer ditampilkan oleh
- Memperbaiki proses zombie saat bekerja dalam mode pengembangan dengan memperbarui refresh ke versi terbaru
- Memperbaiki pengambilan file WebKit AppImage oleh [Atterpac](https://github.com/atterpac)
- Memperbaiki verifikasi paket apt oleh Doctor oleh [Atterpac](https://github.com/Atterpac) di
- Memperbaiki aplikasi yang macet saat ditutup (Darwin) oleh @5aaee9 di
- Memperbaiki warna latar belakang contoh di Windows oleh
- Memperbaiki menu konteks default oleh [mmghv](https://github.com/mmghv) di
- Memperbaiki nilai heksadesimal untuk tombol panah di Darwin oleh
- Membuat seret dan lepas berfungsi di Windows. Ditambahkan oleh
- Memperbaiki bug Linux pada Doctor ketika pengguna tidak memiliki driver yang semestinya
- Memperbaiki penskalaan DPI saat memulai (Windows). Diubah oleh [@almas-x](https://github.com/almas-x) di
- Memperbaiki baris penggantian di `go.mod` agar menggunakan jalur relatif. Memperbaiki jalur Windows dengan
- Memperbaiki penanganan klik baki sistem macOS ketika tidak ada jendela yang ditautkan oleh
- Memperbaiki kegagalan build Windows akibat opsi yang tidak dikenal oleh
- Memperbaiki crash di Windows saat mengeklik kiri ikon baki sistem ketika tidak memiliki
- Memperbaiki baseURL yang salah saat membuka jendela dua kali oleh @5aaee9 dalam PR
- Memperbaiki urutan cabang if dalam metode `WebviewWindow.Restore` oleh
- Menghitung `startURL` dengan benar pada beberapa pemanggilan `GetStartURL` ketika
- Memperbaiki tipe JS struct `Screen` agar sesuai dengan padanannya di Go oleh
- Memperbaiki metode `WML.Reload` untuk memastikan pembersihan yang tepat atas peristiwa yang terdaftar
- Memperbaiki menu konteks kustom yang langsung tertutup di Linux oleh
- Perbaiki jalur keluaran dan ekstensi file model yang dihasilkan oleh binding
- Perbaiki jalur impor file model dalam kode JS yang dihasilkan oleh binding
- Perbaiki seret dan lepas pada beberapa distro Linux oleh
- Perbaiki task yang hilang untuk macOS saat menggunakan `wails3 task dev` oleh
- Perbaiki pendaftaran event yang menyebabkan penetapan pada map nil oleh
- Perbaiki unmarshaling parameter metode terikat oleh
- Perbaiki penanganan beberapa nilai kembalian dari metode terikat oleh
- Perbaiki deteksi npm oleh doctor ketika npm tidak diinstal menggunakan pengelola paket sistem
- Perbaiki MicrosoftEdgeWebview2Setup.exe yang hilang. Terima kasih kepada
- Perbaiki crash acak di Linux akibat penanganan ID jendela oleh @leaanthony. Berdasarkan
- Perbaiki systemTray.setIcon yang menyebabkan crash di Linux oleh
- Pastikan bingkai jendela diterapkan pada pemanggilan pertama fungsi `setFrameless` di

### Diubah

- **PERUBAHAN BESAR**: Kunci map dalam binding JS/TS yang dihasilkan kini ditandai sebagai opsional agar mencerminkan semantik map Go secara akurat. Akses nilai map di TypeScript kini mengembalikan `T | undefined`, bukan `T`, sehingga memerlukan pemeriksaan null atau assertion (#4943) oleh `@fbbdev`
- Ubah penggunaan `Event` menjadi `Events` sesuai dengan perubahan pada `@wailsio/runtime`, serta sesuaikan pemanggilan fungsi dalam dokumentasi di `Features/Events/Event System` oleh @AbdelhadiSeddar
- Pindahkan `EnabledFeatures`, `DisabledFeatures`, dan `AdditionalBrowserArgs` dari opsi per jendela ke `Options.Windows` tingkat aplikasi (#4559) oleh @leaanthony
- Perbarui README untuk contoh `Drag N Drop` dan tegaskan bahwa `Internal Drag and Drop` didemonstrasikan dalam contoh tersebut oleh @ndianabasi
- Ubah berbagai log debug dari Info menjadi Debug (oleh @mbaklor)
- **PERUBAHAN BESAR:** Ganti nama `EnableDragAndDrop` menjadi `EnableFileDrop` dalam opsi jendela
- **PERUBAHAN BESAR:** Ganti nama `DropZoneDetails` menjadi `DropTargetDetails` dalam konteks event
- **PERUBAHAN BESAR:** Ganti nama metode `DropZoneDetails()` menjadi `DropTargetDetails()` pada `WindowEventContext`
- **PERUBAHAN BESAR:** Hapus event `WindowDropZoneFilesDropped`; gunakan `WindowFilesDropped` sebagai gantinya
- **PERUBAHAN BESAR:** Ubah atribut HTML dari `data-wails-dropzone` menjadi `data-file-drop-target`
- **PERUBAHAN BESAR:** Ubah kelas hover CSS dari `wails-dropzone-hover` menjadi `file-drop-target-active`
- **PERUBAHAN BESAR:** Hapus opsi `DragEffect`, `OnEnterEffect`, dan `OnOverEffect` dari Windows (sebelumnya merupakan bagian dari IDropTarget yang telah dihapus)
- Beralih ke goccy/go-json untuk seluruh pemrosesan JSON runtime (binding metode, event, permintaan webview, notifikasi, kvstore), sehingga meningkatkan performa sebesar 21-63% dan mengurangi alokasi memori sebesar 40-60%
- Optimalkan tata letak struct BoundMethod dan cache flag isVariadic untuk mengurangi overhead setiap pemanggilan
- Gunakan buffer argumen yang dialokasikan pada stack untuk metode dengan `<=8` argumen guna menghindari alokasi heap
- Optimalkan pengumpulan hasil dalam pemanggilan metode untuk menghindari alokasi slice bagi satu nilai kembalian
- Gunakan sync.Map untuk cache jenis MIME guna meningkatkan performa konkuren
- Gunakan pool buffer untuk membaca body permintaan transport HTTP
- Alokasikan channel CloseNotify secara lazy dalam pendeteksi jenis konten untuk mengurangi alokasi per permintaan
- Hapus logging CSS debug dari server aset
- Perluas map ekstensi jenis MIME agar mencakup 50+ format web umum (font, audio, video, dan sebagainya)
- Perbarui dokumentasi untuk opsi Window `X/Y` oleh @ruhuang2001
- Perbarui dokumentasi `Frontend Runtime` dengan menambahkan lebih banyak opsi untuk menghasilkan binding frontend oleh @ndianabasi
- Perbarui halaman dokumentasi Server Aset v3 Wails oleh @ndianabasi
- **PERUBAHAN BESAR**: Hapus fungsi dialog tingkat paket (`application.InfoDialog()`, `application.QuestionDialog()`, dan sebagainya). Gunakan pengelola `app.Dialog` sebagai gantinya: `app.Dialog.Info()`, `app.Dialog.Question()`, `app.Dialog.Warning()`, `app.Dialog.Error()`, `app.Dialog.OpenFile()`, `app.Dialog.SaveFile()`
- Perbarui dokumentasi dialog agar sesuai dengan API sebenarnya: gunakan `app.Dialog.*`, `AddButton()` dengan callback (bukan `SetButtons()`), `SetDefaultButton(*Button)` (bukan string), `AddFilter()` (bukan `SetFilters()`), `SetFilename()` (bukan `SetDefaultFilename()`), dan `app.Dialog.OpenFile().CanChooseDirectories(true)` untuk pemilihan folder
- **PERUBAHAN BESAR**: Build produksi kini menjadi default. Untuk membuat build pengembangan, tetapkan `DEV=true` dalam Taskfile Anda. Hasilkan proyek baru untuk contoh oleh @leaanthony
- Saat memancarkan event khusus dengan nol atau satu argumen data, nilai data akan ditetapkan langsung ke field Data tanpa membungkusnya dalam slice oleh [@fbbdev](https://github.com/fbbdev) di [#4633](https://github.com/wailsapp/wails/pull/4633)
- Tray Windows kini mematuhi `SystemTray.Show()`/`Hide()` dengan mengalihkan `NIS_HIDDEN`, sehingga aplikasi dapat benar-benar menghilang dan muncul kembali (#4653).
- Pendaftaran tray menggunakan kembali ikon yang telah di-resolve, menetapkan `NOTIFYICON_VERSION_4` satu kali, dan mengaktifkan `NIF_SHOWTIP` agar tooltip pulih setelah Explorer dimulai ulang (#4653).
- macOS: Gunakan `visibleFrame`, bukan `frame`, untuk memusatkan jendela dengan mengecualikan area bilah menu dan dock
- macOS: Gunakan `visibleFrame`, bukan `frame`, untuk memusatkan jendela dengan mengecualikan area bilah menu dan dock
- Saat menjalankan `wails3 update build-assets` dengan parameter `-config`, nilai yang ditetapkan melalui parameter `-product*` adalah
- `window.NativeWindowHandle()` -> `window.NativeWindow()` oleh @leaanthony di [#4471](https://github.com/wailsapp/wails/pull/4471)
- Refaktor penanganan jendela internal oleh @leaanthony di [#4471](https://github.com/wailsapp/wails/pull/4471)
- Hapus `application.WindowIDKey` dan `application.WindowNameKey` (digantikan oleh `application.WindowKey`) oleh [@leaanthony](https://github.com/leaanthony)
- ContextMenuData kini mengembalikan string, bukan any, oleh [@leaanthony](https://github.com/leaanthony)
- Dalam binding JS/TS, field kelas bertipe array dengan panjang tetap kini diinisialisasi dengan panjang yang diharapkan, bukan dalam keadaan kosong, oleh [@fbbdev](https://github.com/fbbdev) di [#4001](https://github.com/wailsapp/wails/pull/4001)
- ContextMenuData kini mengembalikan string, bukan any, oleh [@leaanthony](https://github.com/leaanthony)
- `application.NewService` tidak lagi menerima options sebagai parameter opsional (gunakan `application.NewServiceWithOptions` sebagai gantinya) oleh [@leaanthony](https://github.com/leaanthony) di [#4024](https://github.com/wailsapp/wails/pull/4024)
- Hapus dependensi `nanoid` oleh [@leaanthony](https://github.com/leaanthony)
- Perbarui contoh Window untuk gaya jendela mica/acrylic/tabbed oleh [@leaanthony](https://github.com/leaanthony)
- Dalam binding JS/TS, file model `internal.js/ts` telah dihapus; semua model kini dapat ditemukan di `models.js/ts` oleh [@fbbdev](https://github.com/fbbdev) di [#4045](https://github.com/wailsapp/wails/pull/4045)
- Dalam binding JS/TS, tipe bernama tidak pernah dirender sebagai alias untuk tipe bernama lain; perilaku lama kini dibatasi hanya untuk alias oleh [@fbbdev](https://github.com/fbbdev) di [#4045](https://github.com/wailsapp/wails/pull/4045)
- Dalam binding JS/TS, pada mode kelas, field struct yang tipenya merupakan parameter tipe ditandai sebagai opsional dan tidak pernah diinisialisasi secara otomatis oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Menghapus ESLint dari templat oleh [@IanVS](https://github.com/IanVS) dalam [#4059](https://github.com/wailsapp/wails/pull/4059)
- Memperbarui tahun hak cipta menjadi 2025 oleh [@IanVS](https://github.com/IanVS) dalam [#4037](https://github.com/wailsapp/wails/pull/4037)
- Menambahkan dokumentasi untuk event.Sender oleh [@IanVS](https://github.com/IanVS) dalam [#4075](https://github.com/wailsapp/wails/pull/4075)
- Dukungan Go 1.24 oleh [@leaanthony](https://github.com/leaanthony)
- Hook `ServiceStartup` kini dipanggil ketika `App.Run` dipanggil, bukan dalam `application.New`, oleh [@fbbdev](https://github.com/fbbdev) dalam [#4066](https://github.com/wailsapp/wails/pull/4066)
- Error `ServiceStartup` kini dikembalikan dari `App.Run` alih-alih menghentikan proses, oleh [@fbbdev](https://github.com/fbbdev) dalam [#4066](https://github.com/wailsapp/wails/pull/4066)
- Pemanggilan binding dan dialog dari JS kini ditolak dengan objek error, bukan string, oleh [@fbbdev](https://github.com/fbbdev) dalam [#4066](https://github.com/wailsapp/wails/pull/4066)
- Meningkatkan pemosisian menu baki sistem pada Windows oleh [@leaanthony](https://github.com/leaanthony)
- Runtime JS telah diporting ke TypeScript oleh [@fbbdev](https://github.com/fbbdev) dalam [#4100](https://github.com/wailsapp/wails/pull/4100)
- Runtime diinisialisasi segera setelah diimpor, tanpa perlu menunggu jendela dimuat, oleh [@fbbdev](https://github.com/fbbdev) dalam [#4100](https://github.com/wailsapp/wails/pull/4100)
- Runtime tidak lagi mengekspor metode init. Impor efek samping dapat digunakan untuk menginisialisasinya, oleh [@fbbdev](https://github.com/fbbdev) dalam [#4100](https://github.com/wailsapp/wails/pull/4100)
- Metode terikat kini mengembalikan `CancellablePromise` yang ditolak dengan `CancelError` jika dibatalkan. Hasil aktual pemanggilan tersebut dibuang, oleh [@fbbdev](https://github.com/fbbdev) dalam [#4100](https://github.com/wailsapp/wails/pull/4100)
- Tipe layanan bawaan kini secara konsisten disebut `Service` oleh [@fbbdev](https://github.com/fbbdev) dalam [#4067](https://github.com/wailsapp/wails/pull/4067)
- Fungsi pembuatan layanan bawaan dengan opsi kini secara konsisten disebut `NewWithConfig` oleh [@fbbdev](https://github.com/fbbdev) dalam [#4067](https://github.com/wailsapp/wails/pull/4067)
- Metode `Select` pada layanan `sqlite` kini dinamai `Query` agar konsisten dengan API Go, oleh [@fbbdev](https://github.com/fbbdev) dalam [#4067](https://github.com/wailsapp/wails/pull/4067)
- Templat: memindahkan runtime ke "dependencies" dan menata file package.json, oleh [@IanVS](https://github.com/IanVS) dalam [#4133](https://github.com/wailsapp/wails/pull/4133)
- Membuat dan menandatangani secara ad hoc bundel aplikasi dalam mode dev untuk mengaktifkan API macOS tertentu, oleh [@popaprozac](https://github.com/popaprozac) dalam [#4171](https://github.com/wailsapp/wails/pull/4171)
- Memindahkan aset build ke direktori khusus platform oleh [@leaanthony](https://github.com/leaanthony)
- Memindahkan dan mengganti nama Taskfile ke direktori khusus platform oleh [@leaanthony](https://github.com/leaanthony)
- Meningkatkan pengalaman secara signifikan ketika `index.html` tidak tersedia, oleh [@leaanthony](https://github.com/leaanthony)
- [Windows] Meningkatkan performa saat meminimalkan dan memulihkan jendela oleh [@leaanthony](https://github.com/leaanthony). Berdasarkan [PR](https://github.com/wailsapp/wails/pull/3955) asli oleh [562589540](https://github.com/562589540)
- Menghapus opsi `ShouldClose` (sebagai gantinya, daftarkan hook untuk events.Common.WindowClosing) oleh [@leaanthony](https://github.com/leaanthony)
- [Windows] Mengurangi kedipan saat membuka jendela oleh [@leaanthony](https://github.com/leaanthony)
- Menghapus `Window.Destroy` karena fungsi ini dimaksudkan sebagai fungsi internal, oleh [@leaanthony](https://github.com/leaanthony)
- Mengganti nama event `WindowClose` menjadi `WindowClosing` oleh [@leaanthony](https://github.com/leaanthony)
- Build frontend kini menggunakan lingkungan vite "development" atau "production" sesuai jenis build, oleh [@leaanthony](https://github.com/leaanthony)
- Memperbarui ke go-webview2 v1.19 oleh [@leaanthony](https://github.com/leaanthony)
- Memastikan fork taskfile digunakan oleh @leaanthony
- Memperbarui fork Taskfile untuk memperbaiki masalah versi saat menginstal menggunakan
- Menggunakan fork Taskfile untuk memperbaiki masalah versi saat menginstal menggunakan
- `service.OnStartup` kini mematikan aplikasi ketika terjadi error dan menjalankan
- Merefactor pengiriman pesan klik baki sistem agar lebih selaras dengan interaksi pengguna oleh
- Penyematan aset menyertakan `all:frontend/dist` untuk mendukung framework yang menghasilkan
- Refactor Taskfile oleh [leaanthony](https://github.com/leaanthony) dalam
- Meningkatkan ke `go-webview2` v1.0.16 oleh
- Memperbaiki tipe `Screen` agar menyertakan `ID`, bukan `Id`, oleh
- Memperbarui versi Wails `go.mod.tmpl` untuk mendukung `application.ServiceOptions` oleh
- Memperbaiki penentuan nama layanan oleh [windom](https://github.com/windom/) dalam
- mkdocs serve kini menggunakan docker oleh [leaanthony](https://github.com/leaanthony)
- Menggabungkan konfigurasi dev ke dalam `config.yml` oleh
- Dialog baki sistem kini secara default menggunakan ikon aplikasi jika tersedia (Windows), oleh
- Meningkatkan pelaporan GPU + Memori untuk macOS oleh
- Menghapus `WebviewGpuIsDisabled` dan `EnableFraudulentWebsiteWarnings`
- Perubahan API Event: `On`/`Emit` -> event pengguna, `OnApplicationEvent` ->
- Perbaikan untuk API Event pada Linux oleh [TheGB0077](https://github.com/TheGB0077) dalam
- [CI] meningkatkan action & memungkinkan action juga dijalankan dalam fork dan
- Mengganti nama `AbsolutePosition()` menjadi `Position()` oleh
- Memperbarui dependensi webkit Linux ke webkit2gtk-4.1 dari webkitgtk2-4.0 untuk
- Skrip runtime JS yang dibundel kini merupakan modul ESM: tag script yang mengimpornya
- Paket `@wailsio/runtime` tidak memublikasikan API-nya pada `window.wails`
- Modul API Window `@wailsio/runtime/src/window` kini mengekspos objek yang memuatnya
- API jendela JS telah diperbarui agar sesuai dengan `WebviewWindow` Go saat ini
- Generator binding kini secara default menggunakan pemanggilan berdasarkan ID. Opsi CLI `-id`
- Tata letak kode binding baru: sebelumnya, file output ditata dalam folder
- Field struct `application.Options.Bind` telah diganti namanya menjadi
- Sintaks baru untuk layanan binding: instans layanan kini harus dibungkus dalam
- Menonaktifkan spinner pada lingkungan nonterminal atau CI oleh

### Dihapus

- **PERUBAHAN YANG TIDAK KOMPATIBEL**: Menghapus `EnabledFeatures`, `DisabledFeatures`, dan `AdditionalLaunchArgs` dari opsi `WindowsWindow` per jendela. Sebagai gantinya, gunakan `Options.Windows.EnabledFeatures`, `Options.Windows.DisabledFeatures`, dan `Options.Windows.AdditionalBrowserArgs` pada tingkat aplikasi. Flag ini berlaku secara global pada lingkungan WebView2 bersama (#4559) oleh @leaanthony
- Menghapus implementasi native `IDropTarget` di Windows dan menggantinya dengan pendekatan berbasis JavaScript (sesuai dengan perilaku v2)
- Menghapus dependensi github.com/wailsapp/mimetype dan menggantinya dengan peta ekstensi yang diperluas serta http.DetectContentType dari pustaka standar, sehingga ukuran biner berkurang sekitar 1.2 MB
- Menghapus dependensi gopkg.in/ini.v1 dengan mengimplementasikan parser minimal untuk file .desktop bagi penjelajah berkas Linux, sehingga menghemat sekitar 45 KB
- Menghapus samber/lo dari kode runtime dengan menggunakan paket slices dari pustaka standar Go 1.21+ dan helper internal minimal, sehingga menghemat sekitar 310 KB
- Menghapus pernyataan printf debug dari handler skema URL Darwin (#4834)
- **PERUBAHAN YANG TIDAK KOMPATIBEL**: Menghapus event `linux:WindowLoadChanged`—gunakan `linux:WindowLoadFinished` sebagai gantinya untuk mendeteksi saat WebView selesai dimuat (#3896) oleh @leaanthony

### Perubahan yang Tidak Kompatibel

- **Refaktorisasi API Manager**: Menata ulang API aplikasi dari struktur datar menjadi manager yang terorganisasi agar kode lebih tertata dan mudah ditemukan, oleh [@leaanthony](https://github.com/leaanthony) dalam [#4359](https://github.com/wailsapp/wails/pull/4359)
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
- Mengganti nama metode Service: `Name` -> `ServiceName`, `OnStartup` -> `ServiceStartup`, `OnShutdown` -> `ServiceShutdown` oleh [@leaanthony](https://github.com/leaanthony)
- Memindahkan metode `Path` dan `Paths` ke paket `application` oleh [@leaanthony](https://github.com/leaanthony)
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

tidak lagi diabaikan dan menggantikan nilai konfigurasi.

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

@atterpac in [#3909](https://github.com/wailsapp/wails/3909)   [ansxuman](https://github.com/ansxuman) in   [#3902](https://github.com/wailsapp/wails/pull/3902)   [atterpac](https://github.com/atterpac) in   [#3867](https://github.com/wailsapp/wails/pull/3867)   by [atterpac](https://github.com/atterpac) in   [#3829](https://github.com/wailsapp/wails/pull/3829)   [leaanthony](https://github.com/leaanthony)   [FerroO2000](https://github.com/FerroO2000) in   [#3856](https://github.com/wailsapp/wails/pull/3856)   [#3873](https://github.com/wailsapp/wails/pull/3873)   [leaanthony](https://github.com/leaanthony)   positioned at the given X/Y location by   [leaanthony](https://github.com/leaanthony) in   [#3885](https://github.com/wailsapp/wails/pull/3885)   [ansxuman](https://github.com/ansxuman) and   [leaanthony](https://github.com/leaanthony) in   [#3823](https://github.com/wailsapp/wails/pull/3823)   by [leaanthony](https://github.com/leaanthony) in   [#3766](https://github.com/wailsapp/wails/pull/3766)   [leaanthony](https://github.com/leaanthony) in   [#3888](https://github.com/wailsapp/wails/pull/3888)   [leaanthony](https://github.com/leaanthony). -

### Diubah

`service.OnShutdown`untuk semua layanan sebelumnya yang dimulai oleh @atterpac di   [#3920](https://github.com/wailsapp/wails/pull/3920)   @atterpac di [#3907](https://github.com/wailsapp/wails/pull/3907)   subfolder oleh @atterpac di   [#3887](https://github.com/wailsapp/wails/pull/3887)   [#3748](https://github.com/wailsapp/wails/pull/3748)   [leaanthony](https://github.com/leaanthony)   [etesam913](https://github.com/etesam913) di   [#3778](https://github.com/wailsapp/wails/pull/3778)   [northes](https://github.com/northes) di   [#3836](https://github.com/wailsapp/wails/pull/3836)   [#3827](https://github.com/wailsapp/wails/pull/3827)   [leaanthony](https://github.com/leaanthony)   [@leaanthony](https://github.com/leaanthony)   [@leaanthony](https://github.com/leaanthony)   (digantikan oleh opsi `EnabledFeatures` dan `DisabledFeatures`) oleh   [leaanthony](https://github.com/leaanthony)

### Diperbaiki

variabel channel oleh @michael-freling di   [#3925](https://github.com/wailsapp/wails/pull/3925)   [ansxuman](https://github.com/ansxuman) di   [#3924](https://github.com/wailsapp/wails/pull/3924)   [#3898](https://github.com/wailsapp/wails/pull/3898)   [#3901](https://github.com/wailsapp/wails/pull/3901)   di [#3886](https://github.com/wailsapp/wails/pull/3886)   [leaanthony](https://github.com/leaanthony) di   [#3841](https://github.com/wailsapp/wails/pull/3841)   `PasteAndMatchStyle` peran dalam menu edit pada Darwin oleh   [johnmccabe](https://github.com/johnmccabe) di   [#3839](https://github.com/wailsapp/wails/pull/3839)   [#3840](https://github.com/wailsapp/wails/issues/3840) di   [#3854](https://github.com/wailsapp/wails/pull/3854) oleh   [kodflow](https://github.com/kodflow)   [@leaanthony](https://github.com/leaanthony)   berbeda. oleh @nickisworking di   [#3789](https://github.com/wailsapp/wails/pull/3789)

## v3.0.0-alpha.7 - 2024-09-18

### Ditambahkan

[mmghv](https://github.com/mmghv) di   [#3665](https://github.com/wailsapp/wails/pull/3665)   [#3682](https://github.com/wailsapp/wails/pull/3682)   [atterpac](https://github.com/atterpac) dan   [leaanthony](https://github.com/leaanthony) di   [#3570](https://github.com/wailsapp/wails/pull/3570)

### Diubah

Peristiwa Aplikasi `OnWindowEvent` -> Peristiwa Jendela, oleh   [leaanthony](https://github.com/leaanthony)   [#3734](https://github.com/wailsapp/wails/pull/3734)   cabang yang diawali dengan `v3/` atau `v3-` oleh   [stendler](https://github.com/stendler) di   [#3747](https://github.com/wailsapp/wails/pull/3747)

### Diperbaiki

[etesam913](https://github.com/etesam913) di   [#3742](https://github.com/wailsapp/wails/pull/3742)   [atterpac](https://github.com/atterpac) di   [#3721](https://github.com/wailsapp/wails/pull/3721)   [atterpac](https://github.com/atterpac) di   [#3675](https://github.com/wailsapp/wails/pull/3675)   [#1811](https://github.com/wailsapp/wails/pull/1811) di   [#3614](https://github.com/wailsapp/wails/pull/3614) oleh   [@stendler](https://github.com/stendler)   [#3720](https://github.com/wailsapp/wails/pull/3720) oleh   [leaanthony](https://github.com/leaanthony)   [#3693](https://github.com/wailsapp/wails/issues/3693) oleh   [@DeltaLaboratory](https://github.com/DeltaLaboratory)   [#3720](https://github.com/wailsapp/wails/pull/3720) oleh   [leaanthony](https://github.com/leaanthony)   [#3693](https://github.com/wailsapp/wails/issues/3693) oleh   [@DeltaLaboratory](https://github.com/DeltaLaboratory)   [leaanthony](https://github.com/leaanthony)   [#3746](https://github.com/wailsapp/wails/pull/3746) oleh   [@stendler](https://github.com/stendler)   [leaanthony](https://github.com/leaanthony)

## v3.0.0-alpha.6 - 2024-07-30

### Diperbaiki

## v3.0.0-alpha.5 - 2024-07-30

### Ditambahkan

[#3580](https://github.com/wailsapp/wails/pull/3580)   [#3580](https://github.com/wailsapp/wails/pull/3580)   klik ikon oleh @5aaee9 di [#2991](https://github.com/wailsapp/wails/pull/2991)   [#2618](https://github.com/wailsapp/wails/pull/2618)   [@fbbdev](https://github.com/fbbdev) di   [#3282](https://github.com/wailsapp/wails/pull/3282)   @[Atterpac](https://github.com/Atterpac)   di [#3022](https://github.com/wailsapp/wails/pull/3022])   [@marcus-crane](https://github.com/marcus-crane) di   [#3146](https://github.com/wailsapp/wails/pull/3146)   [PR](https://github.com/wailsapp/wails/pull/3147)   [PR](https://github.com/wailsapp/wails/pull/3189)   [@fbbdev](https://github.com/fbbdev) di   [#3281](https://github.com/wailsapp/wails/pull/3281)   [aba82cc](https://github.com/wailsapp/wails/commit/aba82cc52787c97fb99afa58b8b63a0004b7ff6c)   berdasarkan [PR](https://github.com/wailsapp/wails/pull/2044) oleh @Mai-Lapyst   [@fbbdev](https://github.com/fbbdev) di   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@fbbdev](https://github.com/fbbdev) di   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@fbbdev](https://github.com/fbbdev) di   [#3295](https://github.com/wailsapp/wails/pull/3295)   paket npm oleh [@fbbdev](https://github.com/fbbdev) di   [#3334](https://github.com/wailsapp/wails/pull/3334)   di [#3354](https://github.com/wailsapp/wails/pull/3354)   `WAILS_VITE_PORT` oleh [@abichinger](https://github.com/abichinger) di   [#3429](https://github.com/wailsapp/wails/pull/3429)   [@abichinger](https://github.com/abichinger) di   [#3431](https://github.com/wailsapp/wails/pull/3431)   [@bruxaodev](https://github.com/bruxaodev) di   [#3667](https://github.com/wailsapp/wails/pull/3667)   [@OlegGulevskyy](https://github.com/OlegGulevskyy) di   [#3674](https://github.com/wailsapp/wails/pull/3674)

### Diperbaiki

[#3606](https://github.com/wailsapp/wails/pull/3606)   [tmclane](https://github.com/tmclane) dalam   [#3515](https://github.com/wailsapp/wails/pull/3515)   [atterpac](https://github.com/atterac) dalam   [#3512](https://github.com/wailsapp/wails/pull/3512)   [atterpac](https://github.com/atterpac) dalam   [#3477](https://github.com/wailsapp/wails/pull/3477)   oleh [Atterpac](https://github.com/atterpac) dalam   [#3320](https://github.com/wailsapp/wails/pull/3320).   dalam [#3306](https://github.com/wailsapp/wails/pull/3306).   [#2972](https://github.com/wailsapp/wails/pull/2972).   [#2982](https://github.com/wailsapp/wails/pull/2982)   [mmghv](https://github.com/mmghv) dalam   [#2750](https://github.com/wailsapp/wails/pull/2750).   [#2753](https://github.com/wailsapp/wails/pull/2753).   [jaybeecave](https://github.com/jaybeecave) dalam   [#3052](https://github.com/wailsapp/wails/pull/3052).   [@pylotlight](https://github.com/pylotlight) dalam   [PR](https://github.com/wailsapp/wails/pull/3039)   terpasang. Ditambahkan oleh [@pylotlight](https://github.com/pylotlight) dalam   [PR](https://github.com/wailsapp/wails/pull/3032)   [PR](https://github.com/wailsapp/wails/pull/3145)   spasi - @leaanthony.   [thomas-senechal](https://github.com/thomas-senechal) dalam PR   [#3207](https://github.com/wailsapp/wails/pull/3207)   [thomas-senechal](https://github.com/thomas-senechal) dalam PR   [#3208](https://github.com/wailsapp/wails/pull/3208)   jendela terlampir [tw1nk](https://github.com/tw1nk) dalam PR   [#3271](https://github.com/wailsapp/wails/pull/3271)   [#3273](https://github.com/wailsapp/wails/pull/3273)   [@fbbdev](https://github.com/fbbdev) dalam   [#3279](https://github.com/wailsapp/wails/pull/3279)   `FRONTEND_DEVSERVER_URL` tersedia.   [#3299](https://github.com/wailsapp/wails/pull/3299)   [@fbbdev](https://github.com/fbbdev) dalam   [#3295](https://github.com/wailsapp/wails/pull/3295)   listener oleh [@fbbdev](https://github.com/fbbdev) dalam   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@abichinger](https://github.com/abichinger) dalam   [#3330](https://github.com/wailsapp/wails/pull/3330)   generator oleh [@fbbdev](https://github.com/fbbdev) dalam   [#3334](https://github.com/wailsapp/wails/pull/3334)   generator oleh [@fbbdev](https://github.com/fbbdev) dalam   [#3334](https://github.com/wailsapp/wails/pull/3334)   [@abichinger](https://github.com/abichinger) dalam   [#3346](https://github.com/wailsapp/wails/pull/3346)   [@hfoxy](https://github.com/hfoxy) dalam   [#3417](https://github.com/wailsapp/wails/pull/3417)   [@hfoxy](https://github.com/hfoxy) dalam   [#3426](https://github.com/wailsapp/wails/pull/3426)   [@fbbdev](https://github.com/fbbdev) dalam   [#3431](https://github.com/wailsapp/wails/pull/3431)   [@fbbdev](https://github.com/fbbdev) dalam   [#3431](https://github.com/wailsapp/wails/pull/3431)   oleh [@pekim](https://github.com/pekim) dalam   [#3458](https://github.com/wailsapp/wails/pull/3458)   [@robin-samuel](https://github.com/robin-samuel).   PR [#3466](https://github.com/wailsapp/wails/pull/3622) oleh   [@5aaee9](https://github.com/5aaee9).   [@windom](https://github.com/windom/) dalam   [#3636](https://github.com/wailsapp/wails/pull/3636).   Windows oleh [@bruxaodev](https://github.com/bruxaodev/) dalam   [#3691](https://github.com/wailsapp/wails/pull/3691).

### Diubah

[mmghv](https://github.com/mmghv) dalam   [#3611](https://github.com/wailsapp/wails/pull/3611)   mendukung Ubuntu 24.04 LTS oleh [atterpac](https://github.com/atterpac) dalam   [#3461](https://github.com/wailsapp/wails/pull/3461)   harus memiliki atribut `type="module"`. Oleh   [@fbbdev](https://github.com/fbbdev) dalam   [#3295](https://github.com/wailsapp/wails/pull/3295)   objek, dan tidak memulai sistem WML. Perubahan ini dilakukan untuk meningkatkan   enkapsulasi. Jika diinginkan, sistem WML dapat dimulai secara manual dengan memanggil   metode `WML.Enable` yang baru. Skrip runtime JS yang disertakan tetap menjalankan kedua   operasi secara otomatis. Oleh [@fbbdev](https://github.com/fbbdev) dalam   [#3295](https://github.com/wailsapp/wails/pull/3295)   objek window sebagai ekspor default. Metode individual tidak dapat lagi diimpor   melalui sintaks impor bernama atau namespace ESM.   API. Beberapa metode telah berubah nama atau prototipe, khususnya: `Screen`   menjadi `GetScreen`; `GetZoomLevel`/`SetZoomLevel` menjadi `GetZoom`/`SetZoom`;   `GetZoom`, `Width`, dan `Height` kini mengembalikan nilai secara langsung, bukan membungkusnya   di dalam objek. Oleh [@fbbdev](https://github.com/fbbdev) dalam   [#3295](https://github.com/wailsapp/wails/pull/3295)   telah dihapus. Gunakan opsi CLI `-names` untuk kembali menggunakan pemanggilan berdasarkan nama.   Oleh [@fbbdev](https://github.com/fbbdev) dalam   [#3468](https://github.com/wailsapp/wails/pull/3468)   dinamai berdasarkan paket yang memuatnya; kini path impor Go lengkap digunakan,   termasuk path modul. Oleh [@fbbdev](https://github.com/fbbdev) dalam   [#3468](https://github.com/wailsapp/wails/pull/3468)   `application.Options.Services`. Oleh [@fbbdev](https://github.com/fbbdev) dalam   [#3468](https://github.com/wailsapp/wails/pull/3468)   pemanggilan ke `application.NewService`. Oleh [@fbbdev](https://github.com/fbbdev) dalam   [#3468](https://github.com/wailsapp/wails/pull/3468)   [@DeltaLaboratory](https://github.com/DeltaLaboratory) dalam   [#3574](https://github.com/wailsapp/wails/pull/3574)
