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

/_   * DIPERBARUI SECARA OTOMATIS *   Catatan rilis dihasilkan dari PR yang telah digabungkan. Jangan edit berkas ini dalam PR kontribusi biasa atau membuat ulang `v3/UNRELEASED_CHANGELOG.md` yang telah dihapus. Lihat `AGENTS.md` untuk pemeliharaan log perubahan. _/

## [Belum dirilis]

## v3.2.0-nightly.20261102 - 2026-11-02

## Diubah
- Gunakan DisableOverscroll sebagai pengganti DisableBounce dalam [PR](https://github.com/wailsapp/wails/pull/6077) oleh @lethusang

## v3.2.0-nightly.20261101 - 2026-11-01

## Ditambahkan
- Tambahkan backend eksperimental macOS dan Linux tanpa CGO dalam [PR](https://github.com/wailsapp/wails/pull/6521) oleh @taliesin-ai
- Build YAML native, pemaketan portabel, dan penandatanganan dalam [PR](https://github.com/wailsapp/wails/pull/6531) oleh @taliesin-ai

## Diperbaiki
- Perbaiki berbagai masalah build nightly v3 dan rapikan dependensi dalam [PR](https://github.com/wailsapp/wails/pull/6519) oleh @taliesin-ai
- Perbaiki validasi rilis v3 dan klasifikasi fixture dalam [PR](https://github.com/wailsapp/wails/pull/6533) oleh @taliesin-ai
- Atasi kegagalan validasi rilis lintas platform dalam [PR](https://github.com/wailsapp/wails/pull/6538) oleh @taliesin-ai

## v3.1.0 - 2026-11-01

## Ditambahkan
- Tambahkan backend eksperimental macOS dan Linux tanpa CGO dalam [PR](https://github.com/wailsapp/wails/pull/6521) oleh @taliesin-ai
- Build YAML native, pemaketan portabel, dan penandatanganan dalam [PR](https://github.com/wailsapp/wails/pull/6531) oleh @taliesin-ai

## Diubah
- Hindari menyalin metadata desktop yang dibuang dalam [PR](https://github.com/wailsapp/wails/pull/6391) oleh @taliesin-ai
- Format panjang konten tanpa fmt dalam [PR](https://github.com/wailsapp/wails/pull/6393) oleh @taliesin-ai
- Pilih root indeks tanpa mengumpulkan kandidat dalam [PR](https://github.com/wailsapp/wails/pull/6404) oleh @taliesin-ai
- Gunakan kembali buffer penyalinan berkas yang tidak mendukung seek dalam [PR](https://github.com/wailsapp/wails/pull/6408) oleh @taliesin-ai
- Hindari mengalokasikan reader untuk mendeteksi prefiks dalam [PR](https://github.com/wailsapp/wails/pull/6410) oleh @taliesin-ai
- Lewati pembacaan deteksi MIME untuk ekstensi yang telah diketahui dalam [PR](https://github.com/wailsapp/wails/pull/6414) oleh @taliesin-ai
- Cegah pertumbuhan antrean dan osilasi penyusutan dalam [PR](https://github.com/wailsapp/wails/pull/6415) oleh @taliesin-ai
- Hindari mengalokasikan konteks permintaan yang akan dibuang dalam [PR](https://github.com/wailsapp/wails/pull/6419) oleh @taliesin-ai
- Lewati penyiapan log debug yang dinonaktifkan dalam [PR](https://github.com/wailsapp/wails/pull/6420) oleh @taliesin-ai
- Alokasikan buffer deteksi MIME secara lazy dalam [PR](https://github.com/wailsapp/wails/pull/6428) oleh @taliesin-ai
- Kurangi alokasi pemformatan hyperlink dalam [PR](https://github.com/wailsapp/wails/pull/6443) oleh @taliesin-ai
- Hilangkan alokasi sementara untuk pemformatan butir dalam [PR](https://github.com/wailsapp/wails/pull/6449) oleh @taliesin-ai
- Kurangi konversi string UUID menjadi satu alokasi dalam [PR](https://github.com/wailsapp/wails/pull/6453) oleh @taliesin-ai
- Kurangi alokasi pemformatan handle statement dalam [PR](https://github.com/wailsapp/wails/pull/6454) oleh @taliesin-ai
- Hilangkan alokasi digest UUID berbasis nama yang bersifat sementara dalam [PR](https://github.com/wailsapp/wails/pull/6455) oleh @taliesin-ai
- Hindari menyalin baris keluaran perintah yang belum selesai dalam [PR](https://github.com/wailsapp/wails/pull/6460) oleh @taliesin-ai
- Salin executable secara streaming selama staging paket dalam [PR](https://github.com/wailsapp/wails/pull/6472) oleh @taliesin-ai
- Hindari memindai ulang baris keluaran yang belum selesai dalam [PR](https://github.com/wailsapp/wails/pull/6473) oleh @taliesin-ai
- Gunakan kembali pola perutean perintah yang telah dikompilasi dalam [PR](https://github.com/wailsapp/wails/pull/6475) oleh @taliesin-ai
- Hindari pencarian metadata selama penemuan runtime dalam [PR](https://github.com/wailsapp/wails/pull/6482) oleh @taliesin-ai
- Kurangi alokasi halaman fallback yang dilokalkan dalam [PR](https://github.com/wailsapp/wails/pull/6485) oleh @taliesin-ai
- Hindari mengompilasi ulang pola untuk setiap pemroses dalam [PR](https://github.com/wailsapp/wails/pull/6490) oleh @taliesin-ai

## Diperbaiki
- Implementasikan SetTooltip untuk baki sistem dalam [PR](https://github.com/wailsapp/wails/pull/6017) oleh @fan711
- Baca mode gelap dari portal pengaturan dalam [PR](https://github.com/wailsapp/wails/pull/6072) oleh @TechHutTV
- Izinkan Window.SetZoom turun di bawah 1.0 dalam [PR](https://github.com/wailsapp/wails/pull/6166) oleh @APshenkin
- Kirim ke setiap listener Go sesuai urutan emit dalam [PR](https://github.com/wailsapp/wails/pull/6275) oleh @DevLumuz
- Baca hook peristiwa dengan hookLock dalam Emit pada [PR](https://github.com/wailsapp/wails/pull/6276) oleh @DevLumuz
- Selalu iklankan properti Menu pada baki agar baki yang dibuat setelah Run tidak diabaikan dalam [PR](https://github.com/wailsapp/wails/pull/6284) oleh @AlbinoGeek
- Tahan pekerjaan thread utama yang dikirim sebelum Run hingga Run dimulai dalam [PR](https://github.com/wailsapp/wails/pull/6287) oleh @AlbinoGeek
- Gunakan URL aset yang benar pada penyedia pemutakhiran GitHub dalam [PR](https://github.com/wailsapp/wails/pull/6304) oleh @gavidroselj
- Pertahankan body permintaan ketika WebView2 tidak menyertakan Content-Length dalam [PR](https://github.com/wailsapp/wails/pull/6326) oleh @strraaa
- Laporkan layar macOS dan GTK dalam ruang DIP (batas DPI campuran) dalam [PR](https://github.com/wailsapp/wails/pull/6346) oleh @julianstorer
- Jangan laporkan notifikasi yang ditutup sebagai klik dalam [PR](https://github.com/wailsapp/wails/pull/6367) oleh @bercedev
- Hentikan proses dev ketika terminal ditutup dalam [PR](https://github.com/wailsapp/wails/pull/6371) oleh @xmasdev
- Petakan proses pembuatan log perubahan ke commit squash-merge dalam pemeriksaan bukti pada [PR](https://github.com/wailsapp/wails/pull/6373) oleh @taliesin-ai
- Pertahankan debounce mulai ulang taskbar selama satu detik dalam [PR](https://github.com/wailsapp/wails/pull/6386) oleh @taliesin-ai
- Pertahankan batasan ukuran saat memasuki layar penuh yang dimulai dalam keadaan tersebut dalam [PR](https://github.com/wailsapp/wails/pull/6387) oleh @taliesin-ai
- Jangan ubah ukuran jendela layar penuh atau yang dimaksimalkan agar sesuai dengan batasan baru dalam [PR](https://github.com/wailsapp/wails/pull/6388) oleh @taliesin-ai
- jangan tampilkan dialog sebelum Run dalam contoh liquid-glass di [PR](https://github.com/wailsapp/wails/pull/6389) oleh @taliesin-ai
- dialog yang diminta sebelum Run tidak melakukan apa pun alih-alih menunggu selamanya di [PR](https://github.com/wailsapp/wails/pull/6390) oleh @taliesin-ai
- gunakan fallback saat entri desktop tidak memiliki perintah di [PR](https://github.com/wailsapp/wails/pull/6392) oleh @taliesin-ai
- deteksi konten saat ini tanpa menyimpan nama file di [PR](https://github.com/wailsapp/wails/pull/6394) oleh @taliesin-ai
- kenali ekstensi aset dengan kombinasi huruf besar-kecil di [PR](https://github.com/wailsapp/wails/pull/6395) oleh @taliesin-ai
- lepaskan buffer pendeteksi konten yang telah digunakan di [PR](https://github.com/wailsapp/wails/pull/6396) oleh @taliesin-ai
- lepaskan referensi ringqueue yang telah dihapus di [PR](https://github.com/wailsapp/wails/pull/6397) oleh @taliesin-ai
- patuhi penghapusan header sebelum pemilihan fallback di [PR](https://github.com/wailsapp/wails/pull/6398) oleh @taliesin-ai
- commit header respons fallback sebelum flush di [PR](https://github.com/wailsapp/wails/pull/6399) oleh @taliesin-ai
- temukan hanya file indeks yang namanya cocok persis di [PR](https://github.com/wailsapp/wails/pull/6400) oleh @taliesin-ai
- hindari body dan pembacaan penuh untuk HEAD yang tidak mendukung seek di [PR](https://github.com/wailsapp/wails/pull/6401) oleh @taliesin-ai
- isi prefiks deteksi MIME meskipun data diperoleh melalui pembacaan singkat di [PR](https://github.com/wailsapp/wails/pull/6402) oleh @taliesin-ai
- commit respons kosong yang berhasil dan tidak mendukung seek di [PR](https://github.com/wailsapp/wails/pull/6403) oleh @taliesin-ai
- hindari duplikasi target URI yang telah diperluas di [PR](https://github.com/wailsapp/wails/pull/6405) oleh @taliesin-ai
- pertahankan path selama perluasan field di [PR](https://github.com/wailsapp/wails/pull/6406) oleh @taliesin-ai
- patuhi direktori pencarian file desktop XDG di [PR](https://github.com/wailsapp/wails/pull/6407) oleh @taliesin-ai
- commit header pendeteksi saat melakukan flush kosong di [PR](https://github.com/wailsapp/wails/pull/6409) oleh @taliesin-ai
- lewati direktori yang membayangi entri desktop di [PR](https://github.com/wailsapp/wails/pull/6411) oleh @taliesin-ai
- pertahankan target yang ditemukan di dalam path executable di [PR](https://github.com/wailsapp/wails/pull/6412) oleh @taliesin-ai
- lepaskan capture callback setelah timer terpicu di [PR](https://github.com/wailsapp/wails/pull/6413) oleh @taliesin-ai
- hindari kecocokan library ldconfig yang keliru di [PR](https://github.com/wailsapp/wails/pull/6416) oleh @taliesin-ai
- batasi kecocokan fallback library sistem file di [PR](https://github.com/wailsapp/wails/pull/6417) oleh @taliesin-ai
- hilangkan data race pada invalidasi cache di [PR](https://github.com/wailsapp/wails/pull/6418) oleh @taliesin-ai
- patuhi konfigurasi logger opsional di [PR](https://github.com/wailsapp/wails/pull/6421) oleh @taliesin-ai
- hindari mutasi URL awal yang digunakan bersama dan port yang kedaluwarsa di [PR](https://github.com/wailsapp/wails/pull/6422) oleh @taliesin-ai
- hentikan percobaan ulang proxy pengembangan saat permintaan dibatalkan di [PR](https://github.com/wailsapp/wails/pull/6423) oleh @taliesin-ai
- kembalikan 404 untuk aset bundel yang tidak ditemukan di [PR](https://github.com/wailsapp/wails/pull/6424) oleh @taliesin-ai
- tangani permintaan HEAD untuk runtime bundel di [PR](https://github.com/wailsapp/wails/pull/6425) oleh @taliesin-ai
- pertahankan path yang di-escape saat merutekan layanan di [PR](https://github.com/wailsapp/wails/pull/6426) oleh @taliesin-ai
- urutkan identifier prarilis numerik yang terlalu besar dengan benar di [PR](https://github.com/wailsapp/wails/pull/6427) oleh @taliesin-ai
- tolak error HTTP sebelum menyimpan unduhan alat di [PR](https://github.com/wailsapp/wails/pull/6429) oleh @taliesin-ai
- patuhi tenggat peluncuran selama penemuan desktop di [PR](https://github.com/wailsapp/wails/pull/6430) oleh @taliesin-ai
- pertahankan body permintaan selama percobaan ulang proxy pengembangan di [PR](https://github.com/wailsapp/wails/pull/6432) oleh @taliesin-ai
- tutup body respons HTTP pada setiap jalur return di [PR](https://github.com/wailsapp/wails/pull/6434) oleh @taliesin-ai
- pertahankan error pembatalan dalam pencarian library paralel di [PR](https://github.com/wailsapp/wails/pull/6435) oleh @taliesin-ai
- tolak respons tag GitHub yang gagal sebelum decoding di [PR](https://github.com/wailsapp/wails/pull/6436) oleh @taliesin-ai
- pertahankan path templat literal yang diapit tanda kurung siku di [PR](https://github.com/wailsapp/wails/pull/6437) oleh @taliesin-ai
- instal wails3 satu kali dan pertahankan error penginstal di [PR](https://github.com/wailsapp/wails/pull/6438) oleh @taliesin-ai
- abaikan entri dan header tersembunyi dalam blok komentar di [PR](https://github.com/wailsapp/wails/pull/6439) oleh @taliesin-ai
- tampilkan error pencarian rilis stabil di [PR](https://github.com/wailsapp/wails/pull/6440) oleh @taliesin-ai
- cetak pesan polos tanpa interpretasi format di [PR](https://github.com/wailsapp/wails/pull/6441) oleh @taliesin-ai
- hindari panic saat menulis setelah memuat null JSON di [PR](https://github.com/wailsapp/wails/pull/6442) oleh @taliesin-ai
- kembalikan snapshot untuk pembacaan seluruh penyimpanan di [PR](https://github.com/wailsapp/wails/pull/6444) oleh @taliesin-ai
- propagasikan error setelah iterasi sebagian baris di [PR](https://github.com/wailsapp/wails/pull/6446) oleh @taliesin-ai
- pertahankan path tujuan ekstraksi secara literal di [PR](https://github.com/wailsapp/wails/pull/6447) oleh @taliesin-ai
- lepaskan handle prepared statement saat menutup dan membuka kembali di [PR](https://github.com/wailsapp/wails/pull/6448) oleh @taliesin-ai
- tangani race antara persiapan dan penutupan atau pembukaan kembali di [PR](https://github.com/wailsapp/wails/pull/6450) oleh @taliesin-ai
- patuhi permintaan HEAD untuk halaman indeks fallback di [PR](https://github.com/wailsapp/wails/pull/6451) oleh @taliesin-ai
- wajibkan handle statement heksadesimal yang lengkap di [PR](https://github.com/wailsapp/wails/pull/6452) oleh @taliesin-ai
- pertahankan basis data dalam memori default saat terjadi konkurensi di [PR](https://github.com/wailsapp/wails/pull/6456) oleh @taliesin-ai
- pisahkan stderr dari output HEAD yang diurai di [PR](https://github.com/wailsapp/wails/pull/6457) oleh @taliesin-ai
- hasilkan aset runtime dari bundel tersemat di [PR](https://github.com/wailsapp/wails/pull/6458) oleh @taliesin-ai
- deteksi placeholder TODO, FIXME, dan TBD di [PR](https://github.com/wailsapp/wails/pull/6459) oleh @taliesin-ai
- validasi respons API catatan rilis di [PR](https://github.com/wailsapp/wails/pull/6461) oleh @taliesin-ai
- patuhi direktori tugas dalam pemeriksaan status di [PR](https://github.com/wailsapp/wails/pull/6462) oleh @taliesin-ai
- cegah kehilangan data saat menyalin file ke file yang sama di [PR](https://github.com/wailsapp/wails/pull/6463) oleh @taliesin-ai
- sertakan direktori tugas dalam identitas cache di [PR](https://github.com/wailsapp/wails/pull/6464) oleh @taliesin-ai
- periksa port satu kali per pemeriksaan di [PR](https://github.com/wailsapp/wails/pull/6465) oleh @taliesin-ai
- hindari panic saat entri cache didekode menjadi null di [PR](https://github.com/wailsapp/wails/pull/6466) oleh @taliesin-ai
- naikkan nilai prarilis numerik tanpa batasan integer di [PR](https://github.com/wailsapp/wails/pull/6467) oleh @taliesin-ai
- cegah pemeriksaan lama menyembunyikan cache miss di [PR](https://github.com/wailsapp/wails/pull/6468) oleh @taliesin-ai
- jalankan tugas default di direktori dasar eksekutor di [PR](https://github.com/wailsapp/wails/pull/6469) oleh @taliesin-ai
- pertahankan tanda persen literal dalam perintah explorer di [PR](https://github.com/wailsapp/wails/pull/6470) oleh @taliesin-ai
- selaraskan direktori perintah prasyarat dan tugas dalam [PR](https://github.com/wailsapp/wails/pull/6471) oleh @taliesin-ai
- ganti symlink dengan berkas biasa yang dihasilkan dalam [PR](https://github.com/wailsapp/wails/pull/6474) oleh @taliesin-ai
- pertahankan DMG yang ada saat build gagal dalam [PR](https://github.com/wailsapp/wails/pull/6476) oleh @taliesin-ai
- pertahankan regex yang dikutip dan garis miring terbalik pada jalur dalam [PR](https://github.com/wailsapp/wails/pull/6477) oleh @taliesin-ai
- validasi ukuran ikon yang dihasilkan sebelum alokasi gambar dalam [PR](https://github.com/wailsapp/wails/pull/6478) oleh @taliesin-ai
- jalankan rantai perintah dan pengalihan melalui shell dalam [PR](https://github.com/wailsapp/wails/pull/6479) oleh @taliesin-ai
- pertahankan paket Linux sebelumnya saat build gagal dalam [PR](https://github.com/wailsapp/wails/pull/6480) oleh @taliesin-ai
- gunakan eksekusi fallback shell khusus platform dalam [PR](https://github.com/wailsapp/wails/pull/6481) oleh @taliesin-ai
- sajikan aset dengan garis miring awal jalur yang berulang dalam [PR](https://github.com/wailsapp/wails/pull/6483) oleh @taliesin-ai
- kenali versi stabil yang memiliki metadata build dalam [PR](https://github.com/wailsapp/wails/pull/6484) oleh @taliesin-ai
- laporkan kegagalan pembuatan sumber daya secara akurat dalam [PR](https://github.com/wailsapp/wails/pull/6486) oleh @taliesin-ai
- cegah race condition saat mengganti pelapor aktif dalam [PR](https://github.com/wailsapp/wails/pull/6487) oleh @taliesin-ai
- serialisasikan penyimpanan berkas cache secara konkuren dalam [PR](https://github.com/wailsapp/wails/pull/6488) oleh @taliesin-ai
- dukung beberapa kandidat kesetaraan templat dalam [PR](https://github.com/wailsapp/wails/pull/6489) oleh @taliesin-ai
- normalisasikan peta kategori tersimpan yang bernilai null dalam [PR](https://github.com/wailsapp/wails/pull/6491) oleh @taliesin-ai
- tampilkan diagnostik dari perintah shell variabel yang gagal dalam [PR](https://github.com/wailsapp/wails/pull/6492) oleh @taliesin-ai
- jangan masukkan rantai shell ke cache Go implisit dalam [PR](https://github.com/wailsapp/wails/pull/6493) oleh @taliesin-ai
- serialisasikan persistensi kategori Linux dalam [PR](https://github.com/wailsapp/wails/pull/6494) oleh @taliesin-ai
- simpan output Go dengan spasi yang dikutip atau di-escape ke cache dalam [PR](https://github.com/wailsapp/wails/pull/6495) oleh @taliesin-ai
- pertahankan ikon Mac saat pengodean gagal dalam [PR](https://github.com/wailsapp/wails/pull/6496) oleh @taliesin-ai
- jangan lagi mengabaikan flag jumlah pengujian Go yang tidak valid dalam [PR](https://github.com/wailsapp/wails/pull/6497) oleh @taliesin-ai
- pertahankan izin sumber untuk salinan berkas baru dalam [PR](https://github.com/wailsapp/wails/pull/6498) oleh @taliesin-ai
- Perbaiki berbagai masalah build nightly v3 dan rapikan dependensi dalam [PR](https://github.com/wailsapp/wails/pull/6519) oleh @taliesin-ai
- perbaiki validasi rilis v3 dan klasifikasi fixture dalam [PR](https://github.com/wailsapp/wails/pull/6533) oleh @taliesin-ai
- atasi kegagalan validasi rilis lintas platform dalam [PR](https://github.com/wailsapp/wails/pull/6538) oleh @taliesin-ai

## v3.1.0-nightly.20261028 - 2026-10-28

## Diubah
- Percepat pembuatan changelog dengan menggunakan kembali pola regex yang telah dikompilasi dalam [PR](https://github.com/wailsapp/wails/pull/6490) oleh @taliesin-ai
- Kurangi alokasi saat menyajikan halaman fallback yang dilokalkan dalam [PR](https://github.com/wailsapp/wails/pull/6485) oleh @taliesin-ai
- Percepat penemuan runtime AppImage dengan menghindari pencarian metadata yang tidak diperlukan dalam [PR](https://github.com/wailsapp/wails/pull/6482) oleh @taliesin-ai
- Gunakan kembali pola perutean perintah yang telah dikompilasi untuk meningkatkan kinerja wake dalam [PR](https://github.com/wailsapp/wails/pull/6475) oleh @taliesin-ai
- Tingkatkan kinerja dengan menghindari pemindaian ulang baris output yang belum selesai dalam [PR](https://github.com/wailsapp/wails/pull/6473) oleh @taliesin-ai
- Alirkan penyalinan berkas executable selama penyiapan paket MSIX untuk mengurangi penggunaan memori dan mempercepat build dalam [PR](https://github.com/wailsapp/wails/pull/6472) oleh @taliesin-ai
- Hindari penyalinan baris output perintah yang belum selesai untuk meningkatkan kinerja dalam [PR](https://github.com/wailsapp/wails/pull/6460) oleh @taliesin-ai
- Kurangi alokasi saat membuat UUID berbasis nama dalam [PR](https://github.com/wailsapp/wails/pull/6455) oleh @taliesin-ai
- Kurangi alokasi saat memformat handle pernyataan SQLite dalam [PR](https://github.com/wailsapp/wails/pull/6454) oleh @taliesin-ai
- Kurangi konversi string UUID menjadi satu alokasi dalam [PR](https://github.com/wailsapp/wails/pull/6453) oleh @taliesin-ai
- Kurangi alokasi saat memformat entri berpoin changelog dalam [PR](https://github.com/wailsapp/wails/pull/6449) oleh @taliesin-ai
- Kurangi alokasi saat memformat hyperlink terminal dalam [PR](https://github.com/wailsapp/wails/pull/6443) oleh @taliesin-ai
- Kurangi alokasi saat menyajikan aset dengan mendeteksi buffer MIME secara lazy dalam [PR](https://github.com/wailsapp/wails/pull/6428) oleh @taliesin-ai
- Lewati penyiapan log debug di server aset saat pencatatan debug dinonaktifkan dalam [PR](https://github.com/wailsapp/wails/pull/6420) oleh @taliesin-ai
- Kurangi alokasi dalam penanganan permintaan webview dengan menghindari konteks permintaan yang dibuang dalam [PR](https://github.com/wailsapp/wails/pull/6419) oleh @taliesin-ai
- Percepat penyajian aset dengan melewati deteksi MIME untuk ekstensi berkas yang dikenal dalam [PR](https://github.com/wailsapp/wails/pull/6414) oleh @taliesin-ai
- Kurangi alokasi saat menyajikan aset dengan tidak menggunakan reader untuk mendeteksi konten dalam [PR](https://github.com/wailsapp/wails/pull/6410) oleh @taliesin-ai
- Kurangi alokasi saat menyajikan berkas aset yang tidak mendukung seek dalam [PR](https://github.com/wailsapp/wails/pull/6408) oleh @taliesin-ai
- Percepat pemilihan root indeks server aset dengan menghindari pengumpulan kandidat dalam [PR](https://github.com/wailsapp/wails/pull/6404) oleh @taliesin-ai
- Tingkatkan kinerja server aset dengan memformat panjang konten tanpa fmt dalam [PR](https://github.com/wailsapp/wails/pull/6393) oleh @taliesin-ai
- Kurangi penyalinan memori saat membaca metadata penjelajah berkas desktop dalam [PR](https://github.com/wailsapp/wails/pull/6391) oleh @taliesin-ai

## Diperbaiki
- Pertahankan izin berkas sumber saat menyalin ke berkas tujuan baru dalam [PR](https://github.com/wailsapp/wails/pull/6498) oleh @taliesin-ai
- Cegah perubahan ukuran jendela layar penuh dan yang dimaksimalkan saat batasan ukuran minimum atau maksimum berubah dalam [PR](https://github.com/wailsapp/wails/pull/6388) oleh @taliesin-ai
- Publikasikan pembaruan riwayat bintang terjadwal untuk repositori wailsapp/wails dalam [PR](https://github.com/wailsapp/wails/pull/6500) oleh @taliesin-ai
- Jangan lagi mengabaikan flag jumlah pengujian Go yang tidak valid di wake dalam [PR](https://github.com/wailsapp/wails/pull/6497) oleh @taliesin-ai
- Pertahankan ikon Mac saat pengodean gagal dalam [PR](https://github.com/wailsapp/wails/pull/6496) oleh @taliesin-ai
- Perbaiki cache miss build untuk output Go dengan spasi yang dikutip atau di-escape dalam [PR](https://github.com/wailsapp/wails/pull/6495) oleh @taliesin-ai
- Serialisasikan persistensi kategori notifikasi di Linux untuk mencegah race condition pada penulisan konkuren dalam [PR](https://github.com/wailsapp/wails/pull/6494) oleh @taliesin-ai
- Jangan masukkan rantai perintah shell ke cache build Go implisit dalam [PR](https://github.com/wailsapp/wails/pull/6493) oleh @taliesin-ai
- Perbaiki kegagalan pemuatan peta kategori notifikasi saat data tersimpan bernilai null dalam [PR](https://github.com/wailsapp/wails/pull/6491) oleh @taliesin-ai
- Tampilkan diagnostik dari perintah shell variabel yang gagal dalam [PR](https://github.com/wailsapp/wails/pull/6492) oleh @taliesin-ai
- Perbaiki pencocokan templat wake saat terdapat beberapa kandidat kesetaraan dalam [PR](https://github.com/wailsapp/wails/pull/6489) oleh @taliesin-ai
- Serialkan penyimpanan file cache secara bersamaan untuk mencegah penulisan cache yang rusak dalam [PR](https://github.com/wailsapp/wails/pull/6488) oleh @taliesin-ai
- Cegah kondisi race saat mengganti pelapor aktif dalam [PR](https://github.com/wailsapp/wails/pull/6487) oleh @taliesin-ai
- Laporkan kegagalan pembuatan sumber daya secara akurat di CLI dalam [PR](https://github.com/wailsapp/wails/pull/6486) oleh @taliesin-ai
- Kenali versi stabil yang menyertakan metadata build di CLI dalam [PR](https://github.com/wailsapp/wails/pull/6484) oleh @taliesin-ai
- Sajikan aset yang jalur permintaannya memuat garis miring awal berulang dalam [PR](https://github.com/wailsapp/wails/pull/6483) oleh @taliesin-ai
- Perbaiki eksekusi shell fallback untuk wake di setiap platform dalam [PR](https://github.com/wailsapp/wails/pull/6481) oleh @taliesin-ai
- Pertahankan paket Linux sebelumnya ketika build gagal dalam [PR](https://github.com/wailsapp/wails/pull/6480) oleh @taliesin-ai
- Salurkan rangkaian perintah dan pengalihan melalui shell dalam [PR](https://github.com/wailsapp/wails/pull/6479) oleh @taliesin-ai
- Pertahankan regex bertanda kutip dan garis miring terbalik pada jalur saat mengurai perintah wake dalam [PR](https://github.com/wailsapp/wails/pull/6477) oleh @taliesin-ai
- Pertahankan file DMG yang sudah ada ketika pemaketan build gagal dalam [PR](https://github.com/wailsapp/wails/pull/6476) oleh @taliesin-ai
- Ganti symlink dengan file reguler yang dihasilkan dalam keluaran binding pada [PR](https://github.com/wailsapp/wails/pull/6474) oleh @taliesin-ai
- Pertahankan tanda persen literal dalam perintah penjelajah file Linux pada [PR](https://github.com/wailsapp/wails/pull/6470) oleh @taliesin-ai
- Selaraskan direktori prasyarat dan perintah tugas agar di-resolve secara konsisten dalam [PR](https://github.com/wailsapp/wails/pull/6471) oleh @taliesin-ai
- Jalankan tugas default di direktori dasar eksekutor dalam [PR](https://github.com/wailsapp/wails/pull/6469) oleh @taliesin-ai
- Cegah pemeriksaan lama menyembunyikan cache miss dalam [PR](https://github.com/wailsapp/wails/pull/6468) oleh @taliesin-ai
- Perbaiki penomoran prarilis CLI agar dapat bertambah tanpa batas luapan bilangan bulat dalam [PR](https://github.com/wailsapp/wails/pull/6467) oleh @taliesin-ai
- Hindari panic ketika entri cache didekode sebagai null dalam [PR](https://github.com/wailsapp/wails/pull/6466) oleh @taliesin-ai
- Periksa port satu kali per pemeriksaan di CLI dalam [PR](https://github.com/wailsapp/wails/pull/6465) oleh @taliesin-ai
- Sertakan direktori tugas dalam identitas cache untuk menghindari hasil wake yang usang dalam [PR](https://github.com/wailsapp/wails/pull/6464) oleh @taliesin-ai
- Cegah kehilangan data saat pemaketan menyalin file ke file itu sendiri dalam [PR](https://github.com/wailsapp/wails/pull/6463) oleh @taliesin-ai
- Pastikan pemeriksaan status menghormati direktori tugas dalam [PR](https://github.com/wailsapp/wails/pull/6462) oleh @taliesin-ai
- Deteksi placeholder TODO, FIXME, dan TBD dalam entri changelog pada [PR](https://github.com/wailsapp/wails/pull/6459) oleh @taliesin-ai
- Hasilkan aset runtime dari bundel tersemat di CLI dalam [PR](https://github.com/wailsapp/wails/pull/6458) oleh @taliesin-ai
- Pisahkan stderr dari keluaran HEAD yang diurai dalam perintah git pada [PR](https://github.com/wailsapp/wails/pull/6457) oleh @taliesin-ai
- Pertahankan basis data SQLite dalam memori default selama akses bersamaan dalam [PR](https://github.com/wailsapp/wails/pull/6456) oleh @taliesin-ai
- Kembalikan respons yang benar terhadap permintaan HEAD untuk halaman indeks fallback dalam [PR](https://github.com/wailsapp/wails/pull/6451) oleh @taliesin-ai
- Tangani kondisi race antara penyiapan statement SQLite dan penutupan atau pembukaan ulang basis data dalam [PR](https://github.com/wailsapp/wails/pull/6450) oleh @taliesin-ai
- Lepaskan handle prepared statement saat basis data SQLite ditutup atau dibuka ulang dalam [PR](https://github.com/wailsapp/wails/pull/6448) oleh @taliesin-ai
- Pertahankan jalur tujuan ekstraksi secara literal saat menggunakan gosod dalam [PR](https://github.com/wailsapp/wails/pull/6447) oleh @taliesin-ai
- Propagasikan kesalahan dari kueri SQLite setelah iterasi baris parsial dalam [PR](https://github.com/wailsapp/wails/pull/6446) oleh @taliesin-ai
- Kembalikan snapshot untuk pembacaan seluruh penyimpanan di kvstore dalam [PR](https://github.com/wailsapp/wails/pull/6444) oleh @taliesin-ai
- Cegah panic saat menulis ke kvstore setelah memuat nilai JSON null dalam [PR](https://github.com/wailsapp/wails/pull/6442) oleh @taliesin-ai
- Cetak pesan terminal teks biasa tanpa menafsirkan verba format dalam [PR](https://github.com/wailsapp/wails/pull/6441) oleh @taliesin-ai
- Tampilkan kesalahan saat mencari rilis stabil terbaru di updater dalam [PR](https://github.com/wailsapp/wails/pull/6440) oleh @taliesin-ai
- Abaikan entri tersembunyi dan header dalam blok komentar changelog pada [PR](https://github.com/wailsapp/wails/pull/6439) oleh @taliesin-ai
- Instal wails3 satu kali dan tampilkan kesalahan penginstal di updater dalam [PR](https://github.com/wailsapp/wails/pull/6438) oleh @taliesin-ai
- Pertahankan jalur templat literal yang diapit tanda kurung siku di gosod dalam [PR](https://github.com/wailsapp/wails/pull/6437) oleh @taliesin-ai
- Tolak respons tag GitHub yang gagal sebelum mendekodenya di updater dalam [PR](https://github.com/wailsapp/wails/pull/6436) oleh @taliesin-ai
- Pertahankan kesalahan pembatalan selama pencarian pustaka secara paralel dalam [PR](https://github.com/wailsapp/wails/pull/6435) oleh @taliesin-ai
- Tutup body respons HTTP pada setiap jalur pengembalian untuk mencegah kebocoran sumber daya dalam [PR](https://github.com/wailsapp/wails/pull/6434) oleh @taliesin-ai
- Pertahankan body permintaan saat proksi pengembangan mencoba ulang permintaan dalam [PR](https://github.com/wailsapp/wails/pull/6432) oleh @taliesin-ai
- Patuhi tenggat peluncuran selama pencarian penjelajah file desktop dalam [PR](https://github.com/wailsapp/wails/pull/6430) oleh @taliesin-ai
- Tolak kesalahan HTTP sebelum menyimpan unduhan alat AppImage dalam [PR](https://github.com/wailsapp/wails/pull/6429) oleh @taliesin-ai
- Urutkan pengidentifikasi prarilis numerik berukuran sangat besar dengan benar dalam perbandingan semver pada [PR](https://github.com/wailsapp/wails/pull/6427) oleh @taliesin-ai
- Pertahankan jalur yang di-escape saat merutekan layanan di server aset dalam [PR](https://github.com/wailsapp/wails/pull/6426) oleh @taliesin-ai
- Tangani permintaan HEAD untuk runtime yang dibundel di server aset dalam [PR](https://github.com/wailsapp/wails/pull/6425) oleh @taliesin-ai
- Kembalikan 404 untuk aset bundel yang tidak ditemukan, bukan respons yang keliru, dalam [PR](https://github.com/wailsapp/wails/pull/6424) oleh @taliesin-ai
- Hentikan percobaan ulang proksi pengembangan saat permintaan dibatalkan dalam [PR](https://github.com/wailsapp/wails/pull/6423) oleh @taliesin-ai
- Hindari mutasi URL awal bersama dan port usang di server aset dalam [PR](https://github.com/wailsapp/wails/pull/6422) oleh @taliesin-ai
- Patuhi konfigurasi logger opsional di server aset dalam [PR](https://github.com/wailsapp/wails/pull/6421) oleh @taliesin-ai
- Hilangkan data race dalam invalidasi cache jalur pustaka pada [PR](https://github.com/wailsapp/wails/pull/6418) oleh @taliesin-ai
- Hindari kecocokan pustaka palsu saat mencari keluaran ldconfig dalam [PR](https://github.com/wailsapp/wails/pull/6416) oleh @taliesin-ai
- Cegah osilasi pertumbuhan dan penyusutan antrean server aset saat dibebani dalam [PR](https://github.com/wailsapp/wails/pull/6415) oleh @taliesin-ai
- Lepaskan tangkapan callback setelah timer debounce dipicu agar referensi usang tidak dipertahankan dalam [PR](https://github.com/wailsapp/wails/pull/6413) oleh @taliesin-ai
- Pertahankan target penjelajah file yang ditemukan di dalam jalur file executable dalam [PR](https://github.com/wailsapp/wails/pull/6412) oleh @taliesin-ai
- Lewati direktori yang membayangi entri desktop di penjelajah file dalam [PR](https://github.com/wailsapp/wails/pull/6411) oleh @taliesin-ai
- Perbaiki header server aset yang tidak terkirim ketika body respons kosong dalam [PR](https://github.com/wailsapp/wails/pull/6409) oleh @taliesin-ai
- Patuhi direktori pencarian file desktop XDG saat menemukan ikon aplikasi dalam [PR](https://github.com/wailsapp/wails/pull/6407) oleh @taliesin-ai
- Pertahankan jalur file selama ekspansi kolom di penjelajah file dalam [PR](https://github.com/wailsapp/wails/pull/6406) oleh @taliesin-ai
- Hindari duplikasi target URI yang diperluas di penjelajah file dalam [PR](https://github.com/wailsapp/wails/pull/6405) oleh @taliesin-ai
- Selesaikan respons aset kosong non-seekable yang berhasil, alih-alih membiarkannya tertunda, dalam [PR](https://github.com/wailsapp/wails/pull/6403) oleh @taliesin-ai
- Isi prefiks pendeteksian MIME dengan benar pada beberapa pembacaan singkat di server aset dalam [PR](https://github.com/wailsapp/wails/pull/6402) oleh @taliesin-ai
- Hindari membaca seluruh isi respons untuk permintaan HEAD non-seekable di server aset dalam [PR](https://github.com/wailsapp/wails/pull/6401) oleh @taliesin-ai
- Sajikan hanya file indeks yang benar-benar cocok saat menyelesaikan permintaan direktori di server aset dalam [PR](https://github.com/wailsapp/wails/pull/6400) oleh @taliesin-ai
- Perbaiki notifikasi Linux yang ditutup agar tidak dilaporkan sebagai klik pada server tanpa dukungan tindakan dalam [PR](https://github.com/wailsapp/wails/pull/6367) oleh @bercedev
- Hentikan proses pengembangan saat terminal ditutup dalam [PR](https://github.com/wailsapp/wails/pull/6371) oleh @xmasdev
- Pertahankan debounce mulai ulang bilah tugas Windows selama satu detik sebelum menghapus flag dalam [PR](https://github.com/wailsapp/wails/pull/6386) oleh @taliesin-ai
- Perbaiki hilangnya header respons fallback di server aset dalam [PR](https://github.com/wailsapp/wails/pull/6399) oleh @taliesin-ai
- Patuhi penghapusan header sebelum pemilihan fallback di server aset dalam [PR](https://github.com/wailsapp/wails/pull/6398) oleh @taliesin-ai
- Lepaskan referensi antrean cincin yang telah dihapus di server aset untuk mencegah kebocoran memori dalam [PR](https://github.com/wailsapp/wails/pull/6397) oleh @taliesin-ai
- Lepaskan buffer pendeteksi konten yang telah digunakan di server aset untuk mengurangi penggunaan memori dalam [PR](https://github.com/wailsapp/wails/pull/6396) oleh @taliesin-ai
- Kenali ekstensi file aset dengan kombinasi huruf besar-kecil saat menyajikan aset dalam [PR](https://github.com/wailsapp/wails/pull/6395) oleh @taliesin-ai
- Deteksi jenis konten dengan benar di server aset tanpa menyimpan nama file dalam [PR](https://github.com/wailsapp/wails/pull/6394) oleh @taliesin-ai
- Gunakan nilai default sebagai fallback saat entri desktop tidak memiliki perintah dalam [PR](https://github.com/wailsapp/wails/pull/6392) oleh @taliesin-ai
- Cegah contoh liquid-glass menampilkan dialog sebelum aplikasi dimulai pada platform selain macOS dalam [PR](https://github.com/wailsapp/wails/pull/6389) oleh @taliesin-ai
- Pertahankan batasan ukuran jendela saat mode layar penuh dimulai dari ukuran yang dibatasi pada Windows dalam [PR](https://github.com/wailsapp/wails/pull/6387) oleh @taliesin-ai
- Perbaiki dialog yang dipanggil sebelum Run agar kembali, bukan terblokir selamanya, dalam [PR](https://github.com/wailsapp/wails/pull/6390) oleh @taliesin-ai

## Keamanan
- Validasi ukuran ikon yang dihasilkan sebelum alokasi gambar untuk mencegah penggunaan memori berlebihan dalam [PR](https://github.com/wailsapp/wails/pull/6478) oleh @taliesin-ai
- Validasi respons API catatan rilis sebelum digunakan dalam [PR](https://github.com/wailsapp/wails/pull/6461) oleh @taliesin-ai
- Wajibkan handle pernyataan heksadesimal yang lengkap di plugin SQLite dalam [PR](https://github.com/wailsapp/wails/pull/6452) oleh @taliesin-ai
- Batasi kecocokan fallback pustaka sistem file untuk mencegah pemuatan pustaka yang tidak dimaksudkan dalam [PR](https://github.com/wailsapp/wails/pull/6417) oleh @taliesin-ai

## v3.1.0-nightly.20261027 - 2026-10-27

## Diubah
- Tautkan halaman 404 ke halaman dokumentasi Wails v2 yang sesuai di v2.wails.io dalam [PR](https://github.com/wailsapp/wails/pull/6366) oleh @taliesin-ai
- Pindahkan situs dokumentasi v2 ke v2.wails.io dalam [PR](https://github.com/wailsapp/wails/pull/6365) oleh @taliesin-ai
- Sajikan dokumentasi v3 dari wails.io dan dokumentasi v2 dari v2.wails.io dalam [PR](https://github.com/wailsapp/wails/pull/6364) oleh @taliesin-ai
- Perbarui catatan rilis v3.0.0 dan koreksi klaim advis keamanan dalam [PR](https://github.com/wailsapp/wails/pull/6363) oleh @taliesin-ai

## Diperbaiki
- Kirim callback listener peristiwa Go secara serial sesuai urutan emisi dalam [PR](https://github.com/wailsapp/wails/pull/6275) oleh @DevLumuz
- Selalu publikasikan properti Menu tray agar tray yang dibuat setelah Run tidak diabaikan di Linux dalam [PR](https://github.com/wailsapp/wails/pull/6284) oleh @AlbinoGeek
- Tahan pekerjaan thread utama yang dikirim sebelum Run hingga aplikasi dimulai dalam [PR](https://github.com/wailsapp/wails/pull/6287) oleh @AlbinoGeek
- Gunakan URL API aset yang benar saat mengunduh aset rilis pemutakhir GitHub dalam [PR](https://github.com/wailsapp/wails/pull/6304) oleh @gavidroselj
- Perbaiki batas dan posisi layar pada konfigurasi macOS dan Linux dengan DPI campuran dalam [PR](https://github.com/wailsapp/wails/pull/6346) oleh @julianstorer
- Pertahankan isi permintaan pada Windows saat WebView2 menghilangkan header Content-Length dalam [PR](https://github.com/wailsapp/wails/pull/6326) oleh @strraaa
- Perbaiki data race dalam emisi hook peristiwa dengan mengambil snapshot hook saat terkunci dalam [PR](https://github.com/wailsapp/wails/pull/6276) oleh @DevLumuz
- Perbaiki `SystemTray.SetTooltip()` yang tidak melakukan apa pun di Linux: deskripsi tooltip StatusNotifierItem kini ditetapkan (#6017)
- Jalankan pemeriksaan bukti changelog dalam alur kerja rilis mingguan dalam [PR](https://github.com/wailsapp/wails/pull/6374) oleh @taliesin-ai
- Perbaiki pemeriksaan bukti changelog rilis agar memetakan proses ke commit squash-merge dalam [PR](https://github.com/wailsapp/wails/pull/6373) oleh @taliesin-ai
- Perbaiki deteksi mode gelap Linux yang mengembalikan mode terang pada desktop gelap hingga skema warna berubah, serta `linux:SystemThemeChanged` yang tidak pernah dipicu pada backend GTK4. Kedua backend kini berbagi satu implementasi yang didukung oleh namespace portal `org.freedesktop.appearance` (#6072)
- `Window.SetZoom` tidak lagi menaikkan faktor di bawah 1.0 menjadi 1.0 pada jendela berbingkai (Windows, macOS, Linux): WebView2 `ZoomFactor`, WKWebView `magnification`, dan tingkat zoom WebKitGTK menerima nilai positif apa pun, sedangkan opsi jendela `Zoom` sudah meneruskan nilai tersebut. Jendela tanpa bingkai tetap mempertahankan batas bawah 1.0 karena wilayah seret dan ubah ukurannya tidak disesuaikan dengan zoom ([#4590](https://github.com/wailsapp/wails/issues/4590)); faktor nonpositif, NaN, dan Inf kini diabaikan, bukan dinaikkan menjadi 1.0 — oleh @APshenkin

## v3.0.0 - 2026-10-26

## Diubah
- Prarilis nightly diterbitkan sebagai `vX.Y.0-nightly.YYYYMMDD` untuk memberikan pratinjau rilis mingguan berikutnya, dan rilis mingguan menaikkan versi minor (`vX.Y.0`), sehingga `@latest` milik Go selalu mengarah ke rilis mingguan atau hotfix; versi paket npm `@wailsio/runtime` selalu sama dengan versi rilis Wails, dengan nightly menggunakan dist-tag `nightly`
- Templat proyek menetapkan `@wailsio/runtime` ke versi CLI yang menghasilkannya
- Wails v3 memerlukan Go 1.26 atau yang lebih baru (Go 1.26.10 atau 1.27.3 dan yang lebih baru menyertakan perbaikan keamanan `crypto/x509`)

## Diperbaiki
- Perbaiki `wails3 update cli` yang tidak menemukan rilis v3: fitur tersebut hanya mempertimbangkan tag v2 dan menginstal dari jalur yang salah
- Perbaiki aplikasi Windows yang mengalami crash dengan `panic: SetTextColor failed` saat menggambar bilah menu aplikasi pada beberapa build Windows: bilah menu kini digambar pada device context yang valid atau menggunakan gambar default sebagai fallback, `UAHMENUITEM.Position` kini cocok dengan `int` C, dan proses menggambar menu tidak pernah panic (#6311)
- Jaga keamanan pembaruan yang telah disiapkan selama mulai ulang masih tertunda: unduhan atau mulai ulang kedua yang dimulai selama proses penutupan tidak dapat menghapus atau memperebutkan payload yang sedang ditunggu pembantu pertukaran untuk diinstal oleh @julianstorer
- Perbaiki aplikasi Windows yang mengalami crash di `applyDarkMode` ketika init milik `pkg/w32` membongkar `uxtheme.dll` sementara pointer fungsi ordinalnya masih digunakan: modul tersebut kini tetap dimuat selama masa hidup proses (#6356)
- Validasi payload yang telah disiapkan sebelum pertukaran dan aplikasi yang telah terpasang sebelum cadangan dibuang, sehingga pembaruan yang rusak dibatalkan dan dipulihkan alih-alih menggantikan aplikasi yang berfungsi, oleh @julianstorer

## Dihapus
- Hapus integrasi Crowdin setelah sinkronisasi terjemahan terakhir dalam [PR](https://github.com/wailsapp/wails/pull/5736) oleh @github-actions[bot]

## Keamanan
- Blokir pesan `wails:` yang dicadangkan dari frame anak lintas asal di Linux (GTK3 dan GTK4) dengan menyuntikkan skrip pelindung ke semua frame, untuk memitigasi GHSA-2p5f-2gx8-mh8m
- Perbarui `github.com/klauspost/compress` ke 1.18.7 (GHSA-259r-337f-4rfw) dan `golang.org/x/image` ke 0.45.0, lalu bangun dengan Go 1.26.10 untuk menyertakan perbaikan `crypto/x509` (GO-2026-6809)
- Perbarui `vite` ke 8.2.4 atau yang lebih baru dalam paket runtime, wisaya penyiapan, contoh, dan templat proyek (GHSA-9jrq-w75r-8gcw, GHSA-rq7h-c2jc-7f22, GHSA-vfpm-58rq-9qcg, GHSA-r825-w8mm-5pxw, GHSA-587h-7wmj-h5pm, GHSA-698v-4r3x-3vcv)

### Keterbatasan yang diketahui
- **Linux:** pelindung frame anak tidak dapat mencegah frame anak *dengan asal yang sama* (misalnya iframe `srcdoc` atau `about:blank`) mencapai penangan pesan native melalui induknya karena WebKitGTK tidak memberikan informasi frame kepada penangan. Frame lintas asal dan frame dalam sandbox tanpa `allow-same-origin` diblokir. Jangan sematkan konten tidak tepercaya dengan asal yang sama di jendela Wails.
- **Android (eksperimental):** bridge JavaScript `window.wails` diekspos ke setiap frame di WebView tanpa pemeriksaan frame utama atau asal. Konten dalam iframe tersemat, termasuk konten lintas asal, dapat memanggil metode terikat dan mengirim pesan `wails:` yang dicadangkan. Jangan sematkan konten pihak ketiga yang tidak tepercaya di WebView aplikasi. Perilaku ini berbeda dari macOS dan iOS, yang mengabaikan pesan `wails:` dari frame anak. Perbaikan direncanakan untuk 3.0.1.
- **Semua platform:** panggilan metode terikat dari frame anak dikirim melalui HTTP ke `/wails/runtime`, bukan melalui penangan pesan native, sehingga pelindung tidak membatasinya. Perlakukan setiap iframe yang Anda sematkan seolah-olah memiliki akses ke metode Go terikat Anda.

## v3.0.0-beta.34 - 2026-10-20

## Diperbaiki
- Perbaiki webview Windows yang tidak kembali mendapatkan fokus saat jendela diaktifkan, serta tombol Tab yang terperangkap di dalam webview alih-alih memindahkan fokus ke jendela native dan kembali lagi (#3783)
- Perbaiki struct vtable COM `internal/webview2/pkg/webview2` (`ICoreWebView2_2`..`_27` dan antarmuka turunan lainnya) yang tidak menyertakan rantai metode yang diwarisi, sehingga panggilan masuk ke slot vtable yang salah (#6129)
- Perbaiki dialog pertanyaan Windows yang mengabaikan label tombol khusus dan tidak pernah menjalankan callback tombol; tombol khusus kini ditampilkan dengan TaskDialog (dengan fallback ke MessageBox dan pemetaan callback berdasarkan posisi ketika comctl32 v6 tidak tersedia) (#5794)
- Perbaiki `wails3 dev` yang macet selama 60 detik saat digunakan dengan Vite: gerbang kesiapan frontend otomatis kini menggunakan probe TCP, bukan HTTP GET yang tidak mengirim header `Accept` dan menerima 404 ([#6209](https://github.com/wailsapp/wails/issues/6209))
- Pertahankan agar menu baki sistem tetap berfungsi ketika diperbarui saat sedang terbuka dalam [PR](https://github.com/wailsapp/wails/pull/6301) oleh @taliesin-ai
- Tangani error COM WebView2 sementara dengan baik dalam operasi runtime di [PR](https://github.com/wailsapp/wails/pull/6302) oleh @leaanthony

## v3.0.0-beta.33 - 2026-10-19

## Ditambahkan
- Generator binding kini menghasilkan modul `bindings/events` yang mengekspor konstanta nama bertipe untuk setiap event yang dapat ditemukan dan didaftarkan dengan `application.RegisterEvent` menggunakan nama konstanta waktu kompilasi (misalnya `UserLogin` untuk `user:login`), yang dapat digunakan dengan `Events.On`/`Emit` tanpa plugin Vite (#5645)

## Diperbaiki
- Perbaiki dialog pertanyaan Windows yang mengabaikan label tombol khusus dan tidak pernah menjalankan callback tombol; tombol khusus kini ditampilkan dengan TaskDialog (dengan fallback ke MessageBox dan pemetaan callback berdasarkan posisi ketika comctl32 v6 tidak tersedia) (#5794)
- Perbaiki `wails3 dev` yang macet selama 60 detik saat digunakan dengan Vite: gerbang kesiapan frontend otomatis kini menggunakan probe TCP, bukan HTTP GET yang tidak mengirim header `Accept` dan menerima 404 ([#6209](https://github.com/wailsapp/wails/issues/6209))
- Pertahankan agar menu baki sistem tetap berfungsi ketika diperbarui saat sedang terbuka dalam [PR](https://github.com/wailsapp/wails/pull/6301) oleh @taliesin-ai
- Tangani error COM WebView2 sementara dengan baik dalam operasi runtime di [PR](https://github.com/wailsapp/wails/pull/6302) oleh @leaanthony

## v3.0.0-beta.32 - 2026-10-18

## Ditambahkan
- Tambahkan header bagian menu native macOS untuk macOS 14+ — lihat [situs dokumentasi](https://v3.wails.io) (20 halaman diperbarui) dalam [PR](https://github.com/wailsapp/wails/pull/6297) oleh @leaanthony

## Diperbaiki
- Buat SetRelativePosition menjadi relatif terhadap area kerja layar di macOS dalam [PR](https://github.com/wailsapp/wails/pull/6299) oleh @leaanthony
- Helper pembaru meluncurkan kembali target yang tidak berubah setelah waktu tunggu induk habis atau penggantian gagal dalam [PR](https://github.com/wailsapp/wails/pull/6298) oleh @taliesin-ai
- Kunci pnpm ke versi 10 untuk dokumentasi berversi dalam [PR](https://github.com/wailsapp/wails/pull/6293) oleh @taliesin-ai

## v3.0.0-beta.31 - 2026-10-18

## Ditambahkan
- Ekspos ekspor subpath `@wailsio/runtime/cancellable` agar `CancellablePromise` dapat diimpor tanpa mengevaluasi entry runtime root (#6153)

## Diubah
- Wajibkan Go 1.26 atau yang lebih baru: modul v3 kini mendeklarasikan `go 1.26.0`, dan proyek yang dihasilkan oleh `wails3 init` serta `wails3 service init` mendeklarasikan versi minimum yang sama

## Diperbaiki
- Perbaiki keluaran generator binding yang berubah ketika dibangun dengan Go 1.27, saat tipe `encoding/json` seperti `Marshaler`, `Unmarshaler`, dan `RawMessage` menjadi alias bagi tipe `encoding/json/v2` dan `encoding/json/jsontext`: tipe tersebut kembali dirender sebagai `any` dan tidak ada model yang dihasilkan untuk `encoding/json/v2` (#6292)
- Validasi URL templat jarak jauh yang diteruskan ke `wails3 init -t` sebelum melakukan kloning: hanya URL `https://`, `http://`, `ssh://`, `git://`, dan `user@host:path` yang diterima; nilai diteruskan ke git setelah `--`; dan git dibatasi hanya pada transport tersebut selama kloning
- Tolak nilai header `x-wails-window-id` yang berada di luar rentang dalam transport runtime HTTP, alih-alih secara diam-diam memotongnya menjadi ID milik jendela lain
- Perkuat ekstraksi arsip pembaru terhadap zip-slip: jalur entri kini diperiksa dengan `filepath.IsLocal` serta pemeriksaan bahwa jalur gabungan tetap berada di dalam direktori tujuan, disertai pengujian untuk entri `..` bertingkat, garis miring terbalik, jalur absolut, dan penulisan melalui symlink dalam arsip zip dan tar.gz
- Perbaiki AppImage dari `wails3 generate appimage` yang crash saat dimulai di luar Debian dan Ubuntu dengan pesan "Tidak dapat membuat proses anak baru: … WebKitNetworkProcess". WebKit yang dibundel kini menjalankan proses helper yang disertakan dalam AppImage. Sandbox bubblewrap WebKit tetap aktif
- Perbaiki build `-tags server` yang gagal dikompilasi untuk Android dan iOS, sekaligus memperbaiki [#6271](https://github.com/wailsapp/wails/issues/6271)
- Perbaiki pembaru yang gagal melakukan pertukaran di Linux ketika direktori sementara berada di sistem berkas lain (seperti tmpfs `/tmp`) dengan menyiapkan file di sebelah biner yang terpasang dan menyalinnya lintas sistem berkas
- Perbaiki pembaru yang mencoba mengganti mount hanya-baca milik AppImage, bukan file `.AppImage`
- Memperbaiki pesan kontrol jendela khusus `wails:` yang diterima dari iframe di macOS dan iOS. Pesan tersebut kini hanya diterima dari frame utama jendela; pesan lain dari iframe tetap diteruskan ke `RawMessageHandler` beserta informasi asalnya
- Menghentikan pengiriman `WindowDidUpdate` dan `WindowWillUpdate` di macOS: `NSWindowDidUpdateNotification` dipicu pada setiap lintasan run loop sehingga halaman menerima sekitar 120 panggilan kosong `evaluateJavaScript:` per detik, dan di macOS, 27.0 pada setiap panggilan tersebut menghabiskan aktivasi sementara dari klik pengguna—`navigator.clipboard.writeText()` dan API lain yang memerlukan aktivasi menolak klik nyata dengan `NotAllowedError` ([#6170](https://github.com/wailsapp/wails/issues/6170)) oleh @APshenkin
- Memperbaiki bahasa default info.json Windows menjadi en-US (0409) agar Explorer menampilkan detail versi (#6210)
- Memperbaiki perender DMA-BUF Linux yang dinonaktifkan pada laptop dengan dGPU NVIDIA di samping iGPU AMD atau Intel: di Wayland, perender tersebut kini hanya dinonaktifkan ketika WebKit merender pada NVIDIA (aplikasi dialihkan ke dGPU atau compositor berjalan di sana), bukan setiap kali modul nvidia dimuat, yang membuat perenderan pada iGPU menjadi sangat lambat; sesi X11 tetap menggunakan perilaku lama karena GPU render Xorg ditentukan oleh konfigurasinya
- Menangani kegagalan penggantian nama lintas perangkat dalam pemutakhir Linux dengan fallback penyalinan bertahap dalam [PR](https://github.com/wailsapp/wails/pull/6135) oleh @john-okeefe
- Memperbaiki `WebviewWindow.SetMenu` di Windows yang tidak menggambar ulang bilah menu kecuali tema bilah menu telah dikonfigurasi, sehingga memperbaiki [#6103](https://github.com/wailsapp/wails/issues/6103)
- Memperbaiki `WebviewWindow.SetMenu(nil)` di Windows yang mengalami panic alih-alih tidak melakukan apa pun seperti di macOS, sehingga memperbaiki [#6104](https://github.com/wailsapp/wails/issues/6104)
- Mempertahankan entri log perubahan saat terjadi push yang bersaing dengan mencoba ulang dalam [PR](https://github.com/wailsapp/wails/pull/6094) oleh @leaanthony
- Memperbaiki kegagalan `go mod vendor` dengan `pattern arm64/WebView2Loader.dll: no matching files found` pada setiap platform dengan menghapus penyematan yang merujuk pada biner yang tidak pernah disertakan dalam modul, sehingga memperbaiki [#5782](https://github.com/wailsapp/wails/issues/5782) dan [#5376](https://github.com/wailsapp/wails/issues/5376), dalam [PR](https://github.com/wailsapp/wails/pull/6031) oleh @Grantmartin2002
- Mencegah pembongkaran jendela duplikat di macOS dalam [PR](https://github.com/wailsapp/wails/pull/6159) oleh @charlesbalderson
- Menangani WM_CONTEXTMENU agar menu baki sistem terbuka saat diklik kanan dalam [PR](https://github.com/wailsapp/wails/pull/6162) oleh @wzk2025
- Memperbaiki kelima varian `InvokeSync*` yang memblokir pemanggilnya selamanya ketika callback yang dikirim mengalami panic, sehingga memperbaiki [#6107](https://github.com/wailsapp/wails/issues/6107), dalam [PR](https://github.com/wailsapp/wails/pull/6108) oleh @RALIST
- Membangun ulang submenu GTK3 yang telah dihancurkan dengan aman dalam [PR](https://github.com/wailsapp/wails/pull/6236) oleh @leaanthony
- Mengimplementasikan shim penerusan worker Windows untuk worker bersarang dalam [PR](https://github.com/wailsapp/wails/pull/6261) oleh @taliesin-ai
- Memperbaiki `OpenFile()` di Windows yang mengabaikan `PromptForMultipleSelection()` ketika `CanChooseDirectories(true)`: pengguna kini dapat memilih beberapa folder sekaligus, seperti di macOS dan Linux (#6263)
- Memperbaiki argumen `double` yang diteruskan berdasarkan nilai ke metode COM WebView2 pada windows/arm64—`ZoomFactor`, `RasterizationScale`, `Expires` cookie, serta skala, ukuran halaman, dan margin pengaturan cetak terdampak di sana (setter pengaturan cetak sebelumnya mengembalikan `ErrDoubleArgUnsupported`): panggilan yang terdampak kini melewati trampoline `runtime.cgocall` yang memuat `d0`, yang tidak dapat dilakukan oleh jalur syscall Go (golang.org/issue/62583) oleh @APshenkin

## Keamanan
- Mode server (`-tags server`) kini hanya menerima panggilan `/wails/runtime` sebagai permintaan `POST` dari origin server sendiri atau origin yang tercantum dalam `ServerOptions.WebSocketOriginPatterns`; metode lain menerima HTTP 405 dan origin browser lain menerima HTTP 403. Deployment di belakang reverse proxy yang menulis ulang header `Host` harus menambahkan origin publiknya ke `WebSocketOriginPatterns`
- Mode server (`-tags server`) yang terikat ke alamat loopback (default `localhost`, `127.0.0.0/8`, atau `::1`) kini menolak permintaan dengan HTTP 403 kecuali header `Host` berupa `localhost`, literal IP, host pengikatan, atau host yang cocok dengan `ServerOptions.WebSocketOriginPatterns`, sehingga memblokir serangan DNS rebinding terhadap `/wails/runtime` dan endpoint WebSocket
- Notifikasi Windows tidak lagi menggunakan PowerShell sebagai fallback ketika API toast WinRT gagal; galat kini dikembalikan

## v3.0.0-beta.30 - 2026-10-15

## Diperbaiki
- Memperbaiki pemutakhir yang mencoba mengganti mount hanya-baca milik AppImage alih-alih berkas `.AppImage`
- Menghentikan pengiriman `WindowDidUpdate` dan `WindowWillUpdate` di macOS: `NSWindowDidUpdateNotification` dipicu pada setiap lintasan run loop sehingga halaman menerima sekitar 120 panggilan kosong `evaluateJavaScript:` per detik, dan di macOS, 27.0 pada setiap panggilan tersebut menghabiskan aktivasi sementara dari klik pengguna—`navigator.clipboard.writeText()` dan API lain yang memerlukan aktivasi menolak klik nyata dengan `NotAllowedError` ([#6170](https://github.com/wailsapp/wails/issues/6170)) oleh @APshenkin
- Memperbaiki bahasa default info.json Windows menjadi en-US (0409) agar Explorer menampilkan detail versi (#6210)
- Menangani kegagalan penggantian nama lintas perangkat dalam pemutakhir Linux dengan fallback penyalinan bertahap dalam [PR](https://github.com/wailsapp/wails/pull/6135) oleh @john-okeefe
- Memperbaiki `WebviewWindow.SetMenu` di Windows yang tidak menggambar ulang bilah menu kecuali tema bilah menu telah dikonfigurasi, sehingga memperbaiki [#6103](https://github.com/wailsapp/wails/issues/6103)
- Memperbaiki `WebviewWindow.SetMenu(nil)` di Windows yang mengalami panic alih-alih tidak melakukan apa pun seperti di macOS, sehingga memperbaiki [#6104](https://github.com/wailsapp/wails/issues/6104)
- Mempertahankan entri log perubahan saat terjadi push yang bersaing dengan mencoba ulang dalam [PR](https://github.com/wailsapp/wails/pull/6094) oleh @leaanthony
- Memperbaiki kegagalan `go mod vendor` dengan `pattern arm64/WebView2Loader.dll: no matching files found` pada setiap platform dengan menghapus penyematan yang merujuk pada biner yang tidak pernah disertakan dalam modul, sehingga memperbaiki [#5782](https://github.com/wailsapp/wails/issues/5782) dan [#5376](https://github.com/wailsapp/wails/issues/5376), dalam [PR](https://github.com/wailsapp/wails/pull/6031) oleh @Grantmartin2002
- Mencegah pembongkaran jendela duplikat di macOS dalam [PR](https://github.com/wailsapp/wails/pull/6159) oleh @charlesbalderson
- Menangani WM_CONTEXTMENU agar menu baki sistem terbuka saat diklik kanan dalam [PR](https://github.com/wailsapp/wails/pull/6162) oleh @wzk2025
- Memperbaiki kelima varian `InvokeSync*` yang memblokir pemanggilnya selamanya ketika callback yang dikirim mengalami panic, sehingga memperbaiki [#6107](https://github.com/wailsapp/wails/issues/6107), dalam [PR](https://github.com/wailsapp/wails/pull/6108) oleh @RALIST

## v3.0.0-beta.29 - 2026-10-13

## Ditambahkan
- Pintasan global di macOS kini memperhitungkan tata letak papan ketik: akselerator huruf diikat ke tombol berlabel pada tata letak non-QWERTY (AZERTY, QWERTZ, Dvorak, ...) melalui `UCKeyTranslate`, bukan ke posisi fisik QWERTY yang tetap.

## Diperbaiki
- Menginisialisasi menu native sebelum menetapkannya di Linux dalam [PR](https://github.com/wailsapp/wails/pull/6233) oleh @abhicodes-007
- Mengatasi regresi kesiapan GA macOS dalam penanganan output dan wizard penyiapan dalam [PR](https://github.com/wailsapp/wails/pull/6237) oleh @taliesin-ai
- Memperbaiki wrapper `AddRef`/`Release` COM WebView2 di Windows yang melaporkan galat terakhir thread yang sudah kedaluwarsa sebagai kegagalan. Wrapper tersebut kini mengembalikan jumlah referensi sehingga operasi release tidak lagi dapat menghentikan aplikasi atau menggagalkan permintaan aset (#3931)
- Memperbaiki `ICoreWebView2Deferral.Release` yang memanggil slot `AddRef`, yang membocorkan deferral dari setiap permintaan aset tertunda di Windows (#3931)

## v3.0.0-beta.28 - 2026-10-05

## Diperbaiki
- Menjaga output aplikasi macOS tetap berada di luar penyedia berkas dalam [PR](https://github.com/wailsapp/wails/pull/6225) oleh @taliesin-ai
- Menyelesaikan wizard penyiapan dan menampilkan detail sistem dalam [PR](https://github.com/wailsapp/wails/pull/6224) oleh @taliesin-ai
- Meningkatkan penanganan galat pemulihan WebView2 dan pembersihan sumber daya di Windows dalam [PR](https://github.com/wailsapp/wails/pull/6222) oleh @leaanthony
- Tingkatkan pemulihan WebView2 di Windows untuk mencegah crash aplikasi dan controller yang terlantar dalam [PR](https://github.com/wailsapp/wails/pull/6220) oleh @taliesin-ai
- Perbaiki panduan otorisasi media macOS dalam [PR](https://github.com/wailsapp/wails/pull/6221) oleh @taliesin-ai
- Perbaiki klik kanan pada ikon baki sistem Linux yang menjalankan penangan klik sekaligus membuka menu, sehingga jendela yang terpasang beralih status pada setiap klik kanan (#6018)
- Pertahankan penangan callback WebView2 hingga dirilis secara native di Windows dalam [PR](https://github.com/wailsapp/wails/pull/6184) oleh @taliesin-ai

## v3.0.0-beta.27 - 2026-10-01

## Diperbaiki
- Perbaiki opsi `Permissions` yang diabaikan pada macOS 12 dan versi yang lebih baru: metode pengambilan media `WKUIDelegate` kini telah diimplementasikan, sehingga `PermissionAllow` dan `PermissionDeny` berlaku untuk permintaan kamera dan mikrofon seperti di Linux dan Windows. `NSCameraUsageDescription` / `NSMicrophoneUsageDescription` serta, jika menggunakan sandbox, entitlement perangkat yang sesuai tetap diperlukan

## v3.0.0-beta.26 - 2026-09-25

## Ditambahkan
- Watcher menerapkan pemeriksaan kesiapan HTTP pada server pengembangan frontend dalam [PR](https://github.com/wailsapp/wails/pull/6150) oleh @atterpac

## Diperbaiki
- Pulihkan kegagalan proses WebView2 agar tidak meninggalkan jendela kosong dalam [PR](https://github.com/wailsapp/wails/pull/6002) oleh @taliesin-ai

## v3.0.0-beta.25 - 2026-09-22

## Diperbaiki
- Hentikan `wails3 dev` dan proses latar belakangnya saat aplikasi utama keluar (#6048)

## v3.0.0-beta.24 - 2026-09-20

## Keamanan
- Batasi izin GITHUB_TOKEN dalam alur kerja build dan publikasi agar mengikuti prinsip hak akses minimum
- Perbaiki traversal jalur dan pelolosan symlink pada middleware aset contoh layar
- Pertahankan fungsi penginstalan dependensi wizard penyiapan sembari membatasinya pada perintah pengelola paket yang disetujui

## v3.0.0-beta.23 - 2026-09-16

## Ditambahkan

- Izinkan penggantian ID aplikasi GTK di Linux dalam [PR](https://github.com/wailsapp/wails/pull/5972) oleh @overlordtm

## Diperbaiki

- Instans tunggal Linux mengklaim nama bus menggunakan ID aplikasi itu sendiri dalam [PR](https://github.com/wailsapp/wails/pull/5973) oleh @overlordtm
- Batalkan permintaan skema khusus yang dihentikan di Linux dalam [PR](https://github.com/wailsapp/wails/pull/6126) oleh @leaanthony

## v3.0.0-beta.22 - 2026-09-14

## Ditambahkan

- Lengkapi terjemahan dan port media komunitas — lihat [situs dokumentasi](https://v3.wails.io) (1812 halaman diperbarui) dalam [PR](https://github.com/wailsapp/wails/pull/6122) oleh @leaanthony
- Dokumentasikan tombol Wails yang dapat disematkan dengan beragam gaya dan ukuran dalam [PR](https://github.com/wailsapp/wails/pull/6120) oleh @leaanthony

## v3.0.0-beta.21 - 2026-09-13

## Ditambahkan

- Sajikan dokumentasi Wails v3 dengan M-Press dalam [PR](https://github.com/wailsapp/wails/pull/6116) oleh @leaanthony

## Diperbaiki

- Urai nilai slug JSON dari frontmatter MPD untuk pembuatan log perubahan dalam [PR](https://github.com/wailsapp/wails/pull/6118) oleh @leaanthony
- Updater menghapus variabel lingkungan helper dan menjalankan kembali target asli setelah kegagalan pencadangan dalam [PR](https://github.com/wailsapp/wails/pull/6080) oleh @cnmax
- Mulai penangan sinyal default selama App.Run dalam [PR](https://github.com/wailsapp/wails/pull/6098) oleh @leaanthony
- Menu Windows menangani menu nil, membebaskan sumber daya yang diganti, dan menggambar ulang bilah menu dalam [PR](https://github.com/wailsapp/wails/pull/6112) oleh @taliesin-ai
- Pulihkan pemaketan MSIX untuk proyek baru yang menggunakan konfigurasi YAML bersama dalam [PR](https://github.com/wailsapp/wails/pull/6115) oleh @leaanthony
- Perbaiki binding JavaScript dan TypeScript hasil pembuatan yang gagal dimuat ketika pembuat model generik merujuk deklarasi helper yang muncul kemudian, serta cegah stack overflow saat membuat model generik yang saling bergantung (#6062)

## v3.0.0-beta.20 - 2026-09-10

## Diubah

- Perbarui tautan showcase Clave ke situs web dan repositori saat ini dalam [PR](https://github.com/wailsapp/wails/pull/6082) oleh @01xR4in

## Diperbaiki

- Batalkan permintaan aset Windows yang dihentikan, termasuk permintaan worker, sembari mempertahankan penangan keepalive selama navigasi. Teruskan konteks permintaan native melalui pembungkus aplikasi pada platform Apple. (#5963, #5969)
- Pertahankan entri log perubahan dari push yang bersaing dengan percobaan ulang dalam [PR](https://github.com/wailsapp/wails/pull/6094) oleh @leaanthony
- Perbaiki kegagalan `go mod vendor` dengan `pattern arm64/WebView2Loader.dll: no matching files found` pada semua platform dengan menghapus penyematan yang merujuk pada biner yang tidak pernah disertakan dalam modul, sehingga memperbaiki [#5782](https://github.com/wailsapp/wails/issues/5782) dan [#5376](https://github.com/wailsapp/wails/issues/5376), dalam [PR](https://github.com/wailsapp/wails/pull/6031) oleh @Grantmartin2002

## Dihapus

- Hapus dukungan loader WebView2 native yang telah digantikan oleh loader Go murni. Perubahan ini menghapus biner `WebView2Loader.dll` yang disematkan dan dependensi `github.com/jchv/go-winloader`. Tag build `native_webview2loader` masih diterima dan tidak lagi menghasilkan error, tetapi tidak berpengaruh pada build v3, dalam [PR](https://github.com/wailsapp/wails/pull/6031) oleh @Grantmartin2002
- Hapus tag build yang tidak digunakan dan opsi FPS dari panduan API macOS dalam [PR](https://github.com/wailsapp/wails/pull/6097) oleh @leaanthony

## v3.0.0-beta.19 - 2026-09-09

## Ditambahkan

- Batasi API privat macOS di balik tag build agar penggunaannya bersifat opt-in — lihat [dokumentasi](https://v3.wails.io/features/browser/integration), [dokumentasi](https://v3.wails.io/features/environment/info), [dokumentasi](https://v3.wails.io/features/windows/basics), [dokumentasi](https://v3.wails.io/features/windows/frameless), [dokumentasi](https://v3.wails.io/features/windows/notch-windows), [dokumentasi](https://v3.wails.io/features/windows/options), [dokumentasi](https://v3.wails.io/guides/build/macos), [dokumentasi](https://v3.wails.io/guides/build/private-macos-apis), dan [dokumentasi](https://v3.wails.io/reference/overview) dalam [PR](https://github.com/wailsapp/wails/pull/6087) oleh @leaanthony

## Diperbaiki

- Tolak permintaan runtime yang melebihi 64 MiB dengan HTTP 413 dalam [PR](https://github.com/wailsapp/wails/pull/6091) oleh @leaanthony

## Keamanan

- Perkuat origin MCP dan akses jarak jauh dengan autentikasi token dalam [PR](https://github.com/wailsapp/wails/pull/6092) oleh @leaanthony

## v3.0.0-beta.18 - 2026-09-08

## Diperbaiki

- Perbaiki kebocoran memori Calloc di Linux dan Darwin dengan menggunakan pointer receiver dalam [PR](https://github.com/wailsapp/wails/pull/6083) oleh @4RH1T3CT0R7

## v3.0.0-beta.17 - 2026-09-06

## Diperbaiki

- Windows: `GetRequest` yang gagal atau nil dalam penangan WebResourceRequested tidak lagi menghentikan proses (`log.Fatal` / panic akibat dereferensi nil) — sebagai gantinya, permintaan dibuang dan dicatat dalam log, dalam [PR](https://github.com/wailsapp/wails/pull/6006) oleh @midagedev

## v3.0.0-beta.16 - 2026-08-29

## Diubah

- Meminta kata sandi notarisasi di jendela terminal baru dalam [PR](https://github.com/wailsapp/wails/pull/6029) oleh @leaanthony

## Diperbaiki

- Menangani jenis klik systray dengan benar di macOS dalam [PR](https://github.com/wailsapp/wails/pull/5919) oleh @ChewbaccaCookie
- CI menghapus repositori apt Microsoft yang tidak digunakan sebelum melakukan pembaruan dalam [PR](https://github.com/wailsapp/wails/pull/6041) oleh @Grantmartin2002

## v3.0.0-beta.15 - 2026-08-27

## Diperbaiki

- Meningkatkan batas waktu penyematan WebView2 menjadi 60 detik dalam [PR](https://github.com/wailsapp/wails/pull/6043) oleh @Grantmartin2002

## v3.0.0-beta.14 - 2026-08-26

## Diperbaiki

- Memberi nama penekanan tombol Control-huruf dengan benar di macOS dalam [PR](https://github.com/wailsapp/wails/pull/6032) oleh @taliesin-ai
- Memperbaiki ikon baki ICO dan mengikuti tema bilah tugas di Windows dalam [PR](https://github.com/wailsapp/wails/pull/6016) oleh @nik9play

## v3.0.0-beta.13 - 2026-08-25

## Diperbaiki

- Tetap melayani pekerjaan thread utama di macOS saat loop modal berjalan dalam [PR](https://github.com/wailsapp/wails/pull/6026) oleh @leaanthony
- Membuat penyimpanan aman seluler dapat mengalami kegagalan dan gagal secara tertutup dalam [PR](https://github.com/wailsapp/wails/pull/5923) oleh @mortenolsrud
- Menjalankan hook peristiwa aplikasi meskipun tidak ada listener yang terdaftar dalam [PR](https://github.com/wailsapp/wails/pull/5999) oleh @archy-rock3t-cloud
- Memperbaiki kesalahan ketik dalam komentar dan dokumentasi yang dilokalkan dalam [PR](https://github.com/wailsapp/wails/pull/6023) oleh @haoku123
- Menghapus biner macOS prakompilasi yang di-commit di bawah `v3/examples` dalam [PR](https://github.com/wailsapp/wails/pull/6025) oleh @4RH1T3CT0R7

## v3.0.0-beta.12 - 2026-08-21

## Ditambahkan

- Menambahkan jendela notifikasi notch macOS beserta contoh siklus hidup dan telemetri — lihat [dokumentasi](https://v3.wails.io/features/windows/notch-windows) dalam [PR](https://github.com/wailsapp/wails/pull/6010) oleh @leaanthony
- Menambahkan dukungan jendela NSPanel macOS dengan opsi baru dan integrasi native — lihat [dokumentasi](https://v3.wails.io/features/windows/options) dalam [PR](https://github.com/wailsapp/wails/pull/6008) oleh @leaanthony

## Diperbaiki

- Mencegah SQLite Prepare mengalami hang saat dipanggil secara bersamaan dalam [PR](https://github.com/wailsapp/wails/pull/5998) oleh @archy-rock3t-cloud

## v3.0.0-beta.11 - 2026-08-20

## Dihapus

- Menghapus pelacak implementasi yang sudah usang dari dokumentasi dalam [PR](https://github.com/wailsapp/wails/pull/6005) oleh @leaanthony

## v3.0.0-beta.10 - 2026-08-19

## Diperbaiki

- Memperbaiki host Linux GTK4 yang menghilangkan argumen peluncuran protokol khusus dan asosiasi file dalam [PR](https://github.com/wailsapp/wails/pull/6000) oleh @midagedev
- Menangani baris yang dihapus dan koreksi dari sumber yang sama dengan benar dalam validasi changelog dalam [PR](https://github.com/wailsapp/wails/pull/5993) oleh @taliesin-ai

## v3.0.0-beta.9 - 2026-08-16

## Ditambahkan

- Menambahkan server wails3 mcp yang aman untuk manajemen proyek dengan bantuan agen dalam [PR](https://github.com/wailsapp/wails/pull/5896) oleh @leaanthony
- Menambahkan dokumentasi untuk model dalam binding — lihat [dokumentasi](https://v3.wails.io/features/bindings/models) dalam [PR](https://github.com/wailsapp/wails/pull/5988) oleh @taliesin-ai
- Mendukung instalasi rpm-ostree untuk sistem Linux atomik dalam [PR](https://github.com/wailsapp/wails/pull/5987) oleh @leaanthony
- Menambahkan pembuatan dan publikasi native mingguan untuk grafik riwayat bintang — lihat [dokumentasi](https://v3.wails.io/credits) dan [dokumentasi](https://v3.wails.io/de/credits) dan [dokumentasi](https://v3.wails.io/fr/credits) dan [dokumentasi](https://v3.wails.io/id/credits) dan [dokumentasi](https://v3.wails.io/ja/credits) dan [dokumentasi](https://v3.wails.io/ko/credits) dan [dokumentasi](https://v3.wails.io/pt/credits) dan [dokumentasi](https://v3.wails.io/ru/credits) dan [dokumentasi](https://v3.wails.io/zh-cn/credits) dan [dokumentasi](https://v3.wails.io/zh-tw/credits) dalam [PR](https://github.com/wailsapp/wails/pull/5986) oleh @leaanthony
- Menambahkan paket mac khusus Darwin untuk menemukan resource bundel aplikasi — lihat [dokumentasi](https://v3.wails.io/guides/build/macos) dalam [PR](https://github.com/wailsapp/wails/pull/5965) oleh @leaanthony
- Menambahkan halaman showcase Condui dan entri indeks — lihat [dokumentasi](https://v3.wails.io/community/showcase/condui) dan [dokumentasi](https://v3.wails.io/community/showcase) dalam [PR](https://github.com/wailsapp/wails/pull/5962) oleh @mgueregath
- Menambahkan halaman showcase Redis Viewer beserta tangkapan layar dan tautan proyek — lihat [dokumentasi](https://v3.wails.io/community/showcase) dan [dokumentasi](https://v3.wails.io/community/showcase/redisviewer) dalam [PR](https://github.com/wailsapp/wails/pull/5984) oleh @redisviewer

## Diubah

- Memperbarui flag aplikasi GTK menjadi G<em>APPLICATION</em>NON_UNIQUE untuk Linux dalam [PR](https://github.com/wailsapp/wails/pull/5971) oleh @overlordtm
- Mencatat peristiwa jendela yang tidak ditemukan pada level debug, bukan peringatan, dalam [PR](https://github.com/wailsapp/wails/pull/5914) oleh @julianstorer

## Diperbaiki

- Memungkinkan pintasan papan ketik macOS yang terdaftar untuk diprioritaskan daripada webview dalam [PR](https://github.com/wailsapp/wails/pull/5902) oleh @julianstorer
- Memperbaiki tautan bilah sisi dokumentasi yang rusak dalam [PR](https://github.com/wailsapp/wails/pull/5937) oleh @northes
- Membatalkan konteks permintaan aset macOS dan iOS ketika WebKit membatalkan tugas skema khusus yang sesuai (#5963)
- Menangani pesan WindowSetFullscreenButtonEnabled dalam [PR](https://github.com/wailsapp/wails/pull/5976) oleh @archy-rock3t-cloud
- Mengimpor Fragment dalam templat preact-ts untuk mengatasi kegagalan build dalam [PR](https://github.com/wailsapp/wails/pull/5979) oleh @haoku123
- Mencegah aplikasi khusus layanan GTK3 lama mengalami crash ketika pencarian layar berjalan sebelum jendela atau tampilan aktif tersedia (#5966)
- Memungkinkan proses rilis dengan versi eksplisit untuk dilanjutkan ketika changelog yang belum dirilis kosong (#5977)

## Keamanan

- Memperbarui lockfile nanoid situs web ke 3.3.18 yang telah ditambal untuk mengatasi advisori keamanan dalam [PR](https://github.com/wailsapp/wails/pull/5985) oleh @taliesin-ai

## v3.0.0-beta.8 - 2026-08-12

## Ditambahkan

- Menambahkan pembuatan URL dokumentasi ke entri changelog otomatis dalam [PR](https://github.com/wailsapp/wails/pull/5957) oleh @taliesin-ai
- Menambahkan Streams: stream byte dua arah antara Go dan JavaScript dengan model pemrograman WebSocket dan tanpa soket yang mendengarkan. Deklarasikan stream di Go dengan `app.HandleStream(name, handler)` dan hubungkan dari frontend dengan `Stream(name)`, yang mengembalikan objek berbentuk `WebSocket`. Go→JS dibawa oleh satu polling tertahan per jendela melalui server aset, sedangkan JS→Go melalui POST biasa; tidak ada yang mengikat port TCP dan tidak ada yang melewati `evaluateJavaScript`. Dalam build server (`-tags server`), handler yang sama disajikan melalui WebSocket sungguhan sebagai gantinya, sehingga kode aplikasi identik di semua build. oleh @leaanthony
- Memindahkan entri changelog kotak surat ke Belum Dirilis dalam [PR](https://github.com/wailsapp/wails/pull/5935) oleh @leaanthony

## Diubah

- Memperbarui pembuatan otomatis bilah sisi dokumentasi dan penurunan tipe penulis blog dalam [PR](https://github.com/wailsapp/wails/pull/5938) oleh @leaanthony

## Diperbaiki

- Inisialisasi WebView2 menggunakan tenggat waktu dan pompa pesan dalam [PR](https://github.com/wailsapp/wails/pull/5952) oleh @leaanthony
- Pengujian cookie WebView2 dilewati di CI kecuali diaktifkan secara eksplisit dan eksekusinya dikunci ke thread OS saat ini dalam [PR](https://github.com/wailsapp/wails/pull/5951) oleh @leaanthony
- Pembuat menu Windows memulihkan ID perintah untuk item induk submenu dalam [PR](https://github.com/wailsapp/wails/pull/5944) oleh @gilad-ch
- Selaraskan citra kompilasi silang resmi dengan baseline dukungan GTK 4.14+ Linux (#5928)
- Konfigurasikan proyek Xcode iOS agar mempertahankan flag linker yang diwarisi dan tambahkan -ObjC dalam [PR](https://github.com/wailsapp/wails/pull/5915) oleh @mortenolsrud
- Perbaiki pergantian koneksi TCP yang berlebihan pada proksi aset `wails3 dev` untuk frontend berukuran besar, yang dapat menghabiskan port efemeral host dan menyebabkan proses yang tidak terkait gagal dengan `EADDRNOTAVAIL`
- Antrekan JavaScript peristiwa per jendela untuk pengiriman berurutan dan backpressure dalam [PR](https://github.com/wailsapp/wails/pull/5934) oleh @leaanthony

## Dihapus

- Hapus pipeline rilis biner desktop: rilis v3 hanya berupa tag dan CLI `wails3` diinstal dengan `go install`. Menghapus `release-v3.yml` dan langkah nightly yang menjalankannya dalam [PR](https://github.com/wailsapp/wails/pull/5946) oleh @leaanthony

## v3.0.0-beta.7 - 2026-08-11

## Ditambahkan

- Tambahkan preferensi pemutaran otomatis macOS untuk menonaktifkan persyaratan tindakan pengguna bagi pemutaran media dalam [PR](https://github.com/wailsapp/wails/pull/5512) oleh @Eyalm321
- Pindahkan entri changelog mailbox ke Belum Dirilis dalam [PR](https://github.com/wailsapp/wails/pull/5935) oleh @leaanthony

## Diubah

- Animasi zoom macOS menggunakan CADisplayLink atau NSTimer agar performanya lebih mulus dalam [PR](https://github.com/wailsapp/wails/pull/5945) oleh @savely-krasovsky

## Diperbaiki

- Konfigurasikan proyek Xcode iOS agar mempertahankan flag linker yang diwarisi dan tambahkan -ObjC dalam [PR](https://github.com/wailsapp/wails/pull/5915) oleh @mortenolsrud
- Perbaiki pergantian koneksi TCP yang berlebihan pada proksi aset `wails3 dev` untuk frontend berukuran besar, yang dapat menghabiskan port efemeral host dan menyebabkan proses yang tidak terkait gagal dengan `EADDRNOTAVAIL`
- Antrekan JavaScript peristiwa per jendela untuk pengiriman berurutan dan backpressure dalam [PR](https://github.com/wailsapp/wails/pull/5934) oleh @leaanthony

### Ditambahkan

- Implementasikan mailbox FIFO asinkron generik untuk pengiriman peristiwa secara berurutan dalam [PR](https://github.com/wailsapp/wails/pull/5851) oleh @savely-krasovsky dan @DevLumuz

## v3.0.0-beta.6 - 2026-08-09

## Ditambahkan

- Implementasikan penyimpanan terbatas di sisi host untuk peristiwa berukuran terlalu besar dan pengiriman JavaScript secara berurutan dalam [PR](https://github.com/wailsapp/wails/pull/5930) oleh @leaanthony
- Implementasikan pantulan Dock macOS untuk kedipan jendela dalam [PR](https://github.com/wailsapp/wails/pull/5921) oleh @julianstorer

## Diperbaiki

- Server aset mempertahankan galat pendeteksi jenis konten dan prefiks yang belum ditulis selama flushing dalam [PR](https://github.com/wailsapp/wails/pull/5931) oleh @leaanthony
- Cegah aplikasi macOS mengalami crash saat mengganti menu aplikasi dari callback Wails
- Perbaiki menu native yang tidak terbaca pada Windows 10 1809 / Windows Server 2019 (build 17763). Ekspor uxtheme mode gelap dibatasi untuk build 18334, sehingga pengaktifan mode gelap pada tingkat aplikasi tidak pernah dijalankan pada host tersebut: latar belakang menu digambar gelap, tetapi Windows tetap menggambar teks menu dengan tema terang sehingga menghasilkan teks gelap pada latar belakang gelap. Ordinal tersebut tersedia sejak 17763, sehingga batasannya kini sesuai.
- Perbaiki `w32.GetStockObject` yang memanggil `GetDeviceCaps` alih-alih `GetStockObject`, yang menyebabkannya mengembalikan 0 untuk setiap objek stok.
- Tingkatkan penanganan dan pelaporan galat pengunduhan bootstrapper WebView2 dalam [PR](https://github.com/wailsapp/wails/pull/5924) oleh @jannskiee

## v3.0.0-beta.5 - 2026-08-07

## Diperbaiki

- Aktivasi aplikasi macOS kini mematuhi kebijakan aktivasi hanya untuk aplikasi reguler dalam [PR](https://github.com/wailsapp/wails/pull/5897) oleh @julianstorer
- Lindungi jendela GTK yang belum diinisialisasi dalam build Linux dalam [PR](https://github.com/wailsapp/wails/pull/5898) oleh @julianstorer
- Tetapkan warna latar belakang opak secara eksplisit untuk jendela WebKit Linux sebelum URL dimuat dalam [PR](https://github.com/wailsapp/wails/pull/5899) oleh @julianstorer

## v3.0.0-beta.4 - 2026-08-05

## Diubah

- Tugas build Android menggunakan arm64 secara default dan deploy-emulator memilih arsitektur host dalam [PR](https://github.com/wailsapp/wails/pull/5890) oleh @mortenolsrud

## Diperbaiki

- Pertahankan status zoom jendela macOS selama penyeretan dan kurangi gerakan dalam [PR](https://github.com/wailsapp/wails/pull/5900) oleh @leaanthony
- Perbaiki build mode server Windows dengan menambahkan `!server` ke batasan build `webview_window_windows_nonclient.go`

## v3.0.0-beta.3 - 2026-08-03

## Ditambahkan

- Dokumentasikan penyelesaian verifikasi beta Fase 10 dalam detail implementasi di [PR](https://github.com/wailsapp/wails/pull/5881) oleh @leaanthony

## Diperbaiki

- Teruskan handle jendela ke API mode gelap Windows dan validasi argumen dalam [PR](https://github.com/wailsapp/wails/pull/5877) oleh @leaanthony
- Pusatkan resolusi status tombol bilah judul macOS untuk jendela tanpa bingkai dalam [PR](https://github.com/wailsapp/wails/pull/5870) oleh @taliesin-ai
- Cegah teks menu native yang tidak terbaca saat aplikasi Windows meminta mode gelap sementara tema aplikasi Windows menggunakan tema terang. Menu kini menggunakan latar belakang native terang yang sesuai hingga Windows dapat merender teks menu gelap.
- Perbaiki menu native yang tidak terbaca pada Windows 10 1809 / Windows Server 2019 (build 17763). Ekspor uxtheme mode gelap dibatasi untuk build 18334, sehingga pengaktifan mode gelap pada tingkat aplikasi tidak pernah dijalankan pada host tersebut: latar belakang menu digambar gelap, tetapi Windows tetap menggambar teks menu dengan tema terang sehingga menghasilkan teks gelap pada latar belakang gelap. Ordinal tersebut tersedia sejak 17763, sehingga batasannya kini sesuai.

## v3.0.0-beta.2 - 2026-08-02

## Diubah

- Promosikan v3 dari alfa menjadi beta
- Dokumentasikan nilai default cerdas system tray dan perilaku sembunyikan otomatis popup, beserta cakupan regresi untuk pemilihan click handler (#5840).
- Updater GitHub mengecualikan aset penginstal Windows secara default dalam [PR](https://github.com/wailsapp/wails/pull/5861) oleh @leaanthony
- Tambahkan dukungan untuk sudut membulat, persegi, dan radius khusus pada jendela macOS tanpa bingkai dalam [PR](https://github.com/wailsapp/wails/pull/5866) oleh @leaanthony

## Diperbaiki

- Laporkan ukuran jendela GTK4 secara langsung dan emit peristiwa perubahan ukuran, maksimalkan, minimalkan, serta status layar penuh dari surface yang dikonfigurasi (#5830).
- Perbaiki crash WebKit Linux saat mengirim Blob atau FormData dalam permintaan fetch di [PR](https://github.com/wailsapp/wails/pull/5854) oleh @taliesin-ai
- Shim fetch meneruskan undefined untuk header Blob/FormData yang tidak ada dalam [PR](https://github.com/wailsapp/wails/pull/5865) oleh @leaanthony

## v3.0.0-alpha2.122 - 2026-08-01

## Ditambahkan

## Diubah

- Tambahkan dukungan untuk sudut membulat, persegi, dan radius khusus pada jendela macOS tanpa bingkai dalam [PR](https://github.com/wailsapp/wails/pull/5866) oleh @leaanthony

## Diperbaiki

- Shim fetch meneruskan undefined untuk header Blob/FormData yang tidak ada dalam [PR](https://github.com/wailsapp/wails/pull/5865) oleh @leaanthony

## v3.0.0-alpha2.121 - 2026-07-31

## Ditambahkan

- Menambahkan dukungan pemaketan DMG macOS dengan opsi dan tugas build baru dalam [PR](https://github.com/wailsapp/wails/pull/5857) oleh @leaanthony

## Diubah

- Updater GitHub secara default mengecualikan aset penginstal Windows dalam [PR](https://github.com/wailsapp/wails/pull/5861) oleh @leaanthony

## Diperbaiki

- Memperbaiki crash WebKit Linux saat mengirim Blob atau FormData dalam permintaan fetch dalam [PR](https://github.com/wailsapp/wails/pull/5854) oleh @taliesin-ai

## v3.0.0-alpha2.120 - 2026-07-31

## Ditambahkan

- Mengimplementasikan tindakan klik dua kali pada bilah judul macOS untuk memaksimalkan atau meminimalkan dalam [PR](https://github.com/wailsapp/wails/pull/5853) oleh @taliesin-ai

## Diubah

- Memperbarui tutorial layanan QR agar menggunakan NewServiceWithOptions dan menambahkan spasi dalam [PR](https://github.com/wailsapp/wails/pull/5849) oleh @jeongkyu

## Diperbaiki

- Menjaga WKWebView tetap responsif selama zoom di macOS dalam [PR](https://github.com/wailsapp/wails/pull/5856) oleh @leaanthony
- Memperbaiki kueri ukuran jendela GTK4 serta memancarkan peristiwa perubahan ukuran, pemaksimalan, peminimalan, dan status layar penuh dari `GdkSurface` yang dikonfigurasi.

## v3.0.0-alpha2.119 - 2026-07-27

## Diperbaiki

- Memperbarui dokumentasi dalam beberapa bahasa agar menyertakan diagram arsitektur dalam [PR](https://github.com/wailsapp/wails/pull/5833) oleh @taliesin-ai

## v3.0.0-alpha2.118 - 2026-07-26

## Ditambahkan

- Menyediakan jalur default untuk masukan dan keluaran pembuatan ikon dalam [PR](https://github.com/wailsapp/wails/pull/5825) oleh @taliesin-ai
- Menambahkan modul entri sumber ke sideEffects package.json runtime dalam [PR](https://github.com/wailsapp/wails/pull/5797) oleh @savely-krasovsky
- Menambahkan bagian lisensi dan asal-usul ke panduan kontribusi dalam [PR](https://github.com/wailsapp/wails/pull/5816) oleh @taliesin-ai

## Diperbaiki

- Menerapkan CSS tanpa bingkai GTK4 dengan cakupan terbatas untuk menghapus radius batas dalam [PR](https://github.com/wailsapp/wails/pull/5800) oleh @savely-krasovsky
- Menangani kegagalan mendapatkan posisi kursor secara baik di Windows untuk menu pop-up dan enumerasi layar dalam [PR](https://github.com/wailsapp/wails/pull/5789) oleh @wayneforrest
- Dialog buka berkas macOS kini memfilter ekstensi dengan benar dan memvalidasi berkas yang diizinkan berdasarkan sufiks dalam [PR](https://github.com/wailsapp/wails/pull/5678) oleh @phergul
- Melindungi inisialisasi mode gelap Windows dari pemanggilan API nil dalam [PR](https://github.com/wailsapp/wails/pull/5793) oleh @roachadam
- Memperbaiki kegagalan build 32-bit pada updater: konstanta `maxArchiveTotalSize` (2 GiB) melampaui kapasitas `int` platform saat diteruskan ke `fmt.Errorf` pada `GOARCH=386`. Sekarang konstanta tersebut secara eksplisit bertipe `int64`.
- Memperbaiki panic akibat penunjuk nil saat startup ketika sebuah jendela menggunakan bilah judul Gelap (atau gelap-sistem) pada build Windows yang tidak memuat API uxtheme mode gelap, seperti Windows 10 1809 / Windows Server 2019 (build 17763). Pemanggilan `AllowDarkModeForWindow` dalam penyiapan tema jendela kini dilindungi dari nil, sesuai dengan perlindungan yang sudah digunakan dalam `w32.SetMenuTheme`.

## v3.0.0-alpha2.117 - 2026-07-08

## Ditambahkan

- Mengimplementasikan logika hit-test khusus untuk wilayah non-klien di Windows dalam [PR](https://github.com/wailsapp/wails/pull/5462) oleh @savely-krasovsky

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
- Memperbaiki kegagalan build Windows saat menggunakan tag build `server`, yang disebabkan oleh berkas GUI Windows yang tidak memiliki batasan build `!server` yang sudah dimiliki padanannya di macOS dan Linux (#5680)

## v3.0.0-alpha2.114 - 2026-07-05

## Ditambahkan

- Mengimplementasikan protokol Update Manifest dan penyedia endpoint dalam [PR](https://github.com/wailsapp/wails/pull/5720) oleh @taliesin-ai

## Diubah

- Menggabungkan binding `webview2` ke dalam modul v3 sebagai `v3/internal/webview2`, serta menghapus modul mandiri, alur kerja rilis/sinkronisasi nightly, dan pergantian versi go.mod-nya (v3 adalah satu-satunya konsumennya) dalam [PR](https://github.com/wailsapp/wails/pull/5711) oleh @taliesin-ai

## Diperbaiki

- Memindahkan perbaikan deteksi skala monitor WebView2 dan sinkronisasi ulang host saat perubahan DPI ke bagian Belum Dirilis dalam [PR](https://github.com/wailsapp/wails/pull/5750) oleh @taliesin-ai
- Memperbarui marshaling COM WebView2 untuk parameter float64 dan BOOL dalam [PR](https://github.com/wailsapp/wails/pull/5741) oleh @wayneforrest
- Mencegah panic dan dereferensi nil saat memperbarui dan memusnahkan ikon baki sistem Windows dalam [PR](https://github.com/wailsapp/wails/pull/5703) oleh @wayneforrest
- Memperbaiki jendela tersembunyi yang tidak kembali disembunyikan dengan benar di Windows dalam [PR](https://github.com/wailsapp/wails/pull/5743) oleh @wayneforrest
- Menyinkronkan visibilitas pengontrol WebView2 dengan tindakan meminimalkan, memaksimalkan, dan memulihkan jendela dalam [PR](https://github.com/wailsapp/wails/pull/5742) oleh @wayneforrest

### Diperbaiki

- Mengaktifkan kembali deteksi skala monitor WebView2 dan membatasi sinkronisasi ulang host agar hanya dilakukan saat DPI berubah dalam [PR](https://github.com/wailsapp/wails/pull/5734) oleh @taliesin-ai, berdasarkan perbaikan yang divalidasi oleh @randalmurphal, dengan verifikasi akar penyebab oleh @eleclin dan pengujian perangkat keras oleh @qq540491950

## v3.0.0-alpha2.113 - 2026-07-04

## Ditambahkan

- Menambahkan peringatan saat membangun AAB rilis tanpa menetapkan `ANDROID_KEYSTORE_FILE` (Google Play menolak bundle yang ditandatangani dengan kunci debug) serta mendokumentasikan pemaketan dan penandatanganan App Bundle dalam [PR](https://github.com/wailsapp/wails/pull/5730) oleh @taliesin-ai
- Menambahkan dokumentasi Mengapa Wails dalam beberapa bahasa dalam [PR](https://github.com/wailsapp/wails/pull/5739) oleh @taliesin-ai
- Mendukung pemetaan Go time.Time ke JS Date atau string dalam binding dalam [PR](https://github.com/wailsapp/wails/pull/5398) oleh @fbbdev
- Tambahkan tugas pemaketan Android App Bundle (AAB) (`bundle`, `bundle:fat`, `assemble:aab`, `assemble:aab:release`) untuk pengiriman ke Play Store — tugas APK tetap tersedia untuk pengujian lokal/emulator dalam [PR](https://github.com/wailsapp/wails/pull/5728) oleh @mortenolsrud (memperbaiki [#5726](https://github.com/wailsapp/wails/issues/5726))
- Tambahkan target tugas untuk perangkat Android fisik serta lanjutkan izin kamera/lokasi dalam [PR](https://github.com/wailsapp/wails/pull/5735) oleh @taliesin-ai

## Diubah

- Tingkatkan `webview2` ke v1.0.28 ([catatan rilis](https://github.com/wailsapp/wails/releases/tag/webview2%2Fv1.0.28)).
- Tingkatkan `compileSdk`/`targetSdk` pada templat Android dari 34 ke 35, yang diwajibkan Google Play untuk pengiriman aplikasi baru, dalam [PR](https://github.com/wailsapp/wails/pull/5730) oleh @taliesin-ai

## Diperbaiki

- Perbaiki mask avatar yang telah dipanggang di sponsorkit dalam [PR](https://github.com/wailsapp/wails/pull/5745) oleh @leaanthony
- Perbaiki pembuatan otomatis AVD Android yang memilih citra sistem atau versi cmdline-tools yang salah akibat pengurutan versi secara leksikografis dalam [PR](https://github.com/wailsapp/wails/pull/5730) oleh @taliesin-ai
- Perbaiki wizard penyiapan yang menyarankan versi Android NDK lama (sekarang 26.3.11579264, sesuai dengan persyaratan yang didokumentasikan) dalam [PR](https://github.com/wailsapp/wails/pull/5730) oleh @taliesin-ai
- Perbarui dokumentasi bahasa Prancis untuk SvelteKit dan opsi dalam [PR](https://github.com/wailsapp/wails/pull/5744) oleh @leaanthony
- Perbaiki SIGSEGV saat menginventarisasi layar macOS ketika tampilan berubah dalam [PR](https://github.com/wailsapp/wails/pull/5516) oleh @flofreud

## v3.0.0-alpha2.112 - 2026-07-03

## Ditambahkan

- Tambahkan generator SVG kontributor berbasis Go dan perbarui halaman kredit dokumentasi/situs web dalam [PR](https://github.com/wailsapp/wails/pull/5724) oleh @taliesin-ai
- Dukung pemetaan Go time.Time ke JS Date atau string dalam binding melalui [PR](https://github.com/wailsapp/wails/pull/5398) oleh @fbbdev

## Diubah

- Ganti pipeline gambar sponsor berbasis Node dengan generator Go dalam [PR](https://github.com/wailsapp/wails/pull/5719) oleh @taliesin-ai

## Diperbaiki

- Perbaiki skrip instalasi dependensi aset build Android dalam [PR](https://github.com/wailsapp/wails/pull/5729) oleh @taliesin-ai
- Tolak karakter kontrol U+0085 (NEXT LINE) dalam `ValidateAndSanitizeURL`, sehingga cakupan karakter spasi pada validator URL menjadi lengkap
- Hitung ulang bingkai DWM saat DPI berubah untuk jendela tanpa bingkai dalam [PR](https://github.com/wailsapp/wails/pull/4785) oleh @leaanthony
- Perbaiki kegagalan deteksi zona pelepasan DnD pada penskalaan selain 100% di Windows dalam [PR](https://github.com/wailsapp/wails/pull/4632) oleh @yulesxoxo
- Tambahkan pengelolaan memori Objective-C secara eksplisit untuk objek Cocoa pada dialog, menu, baki sistem, dan notifikasi Darwin dalam [PR](https://github.com/wailsapp/wails/pull/5714) oleh @taliesin-ai
- Perbaiki bug backend CGO Linux dan masalah baki sistem dalam [PR](https://github.com/wailsapp/wails/pull/5718) oleh @taliesin-ai

## v3.0.0-alpha2.111 - 2026-07-01

## Ditambahkan

- Tambahkan HappyTools ke etalase komunitas dalam [PR](https://github.com/wailsapp/wails/pull/5061) oleh @Aliuyanfeng
- Tambahkan dukungan lokal Indonesia dan dokumentasi lengkap dalam [PR](https://github.com/wailsapp/wails/pull/5643) oleh @triadmoko
- Tambahkan opsi DisableMenu ke WindowsWindow dalam [PR](https://github.com/wailsapp/wails/pull/4813) oleh @leaanthony

## Diubah

- Perbarui templat Taskfile dan CLI untuk meneruskan tugas build/paket dengan GOOS dan ARCH dalam [PR](https://github.com/wailsapp/wails/pull/5617) oleh @leaanthony

## Diperbaiki

- Perbaiki masalah tab jendela Mac dalam [PR](https://github.com/wailsapp/wails/pull/5708) oleh @taliesin-ai

## Dihapus

- Hapus file MDX terjemahan bahasa Jerman dari bagian kontribusi, fitur, dan panduan dalam [PR](https://github.com/wailsapp/wails/pull/5702) oleh @taliesin-ai

## v3.0.0-alpha2.110 - 2026-06-30

## Ditambahkan

- Implementasikan pemuatan ulang dan pemuatan ulang paksa WebView di macOS, serta tambahkan pemulihan saat proses WebContent dihentikan, dalam [PR](https://github.com/wailsapp/wails/pull/5129) oleh @wayneforrest
- Tambahkan dokumentasi lengkap berbahasa Jerman untuk kontribusi, fitur, dan panduan dalam [PR](https://github.com/wailsapp/wails/pull/5396) oleh @leaanthony
- Tingkatkan notifikasi dengan suara, lampiran, penjadwalan, dan API pembaruan dalam [PR](https://github.com/wailsapp/wails/pull/5333) oleh @popaprozac

## Diperbaiki

- Hitung ulang bingkai DWM saat DPI berubah untuk jendela tanpa bingkai dalam [PR](https://github.com/wailsapp/wails/pull/4785) oleh @leaanthony
- Perbaiki kegagalan deteksi zona pelepasan DnD pada penskalaan selain 100% di Windows dalam [PR](https://github.com/wailsapp/wails/pull/4632) oleh @yulesxoxo

## v3.0.0-alpha2.109 - 2026-06-29

## Ditambahkan

- Tambahkan contoh kode ke dokumentasi EventsEmit dalam [PR](https://github.com/wailsapp/wails/pull/5026) oleh @iamhabbeboy
- Tambahkan opsi hosting visual WebView2 di Windows dalam [PR](https://github.com/wailsapp/wails/pull/5380) oleh @MerIijn
- Tambahkan Klustr ke dokumentasi etalase komunitas dalam [PR](https://github.com/wailsapp/wails/pull/5536) oleh @SametKUM
- Tambahkan Kira ke etalase komunitas beserta halaman baru dan entri log perubahan dalam [PR](https://github.com/wailsapp/wails/pull/5685) oleh @thiennguyen93
- Tambahkan bagian umpan balik ke panduan layanan MCP dalam [PR](https://github.com/wailsapp/wails/pull/5694) oleh @taliesin-ai

## Diubah

- Mode server kini memiliki build produksi kelas utama yang konsisten dengan tugas build desktop (#5693). Secara default, `task build:server` membuat biner produksi (`-tags server,production`, `-trimpath`, simbol dihapus), serta menerima `DEV=true` (server pengembangan), `OBFUSCATED=true` (garble), dan `EXTRA_TAGS`. `task run:server` menjalankan server pengembangan. `Dockerfile.server` / `task build:docker` terlebih dahulu membuat server produksi (`-tags server,production`) dan frontend produksi; secara default, citranya menggunakan build statis Go murni pada distroless/static, dengan `CGO_ENABLED`, `GO_IMAGE`, dan `RUNTIME_IMAGE` yang diekspos sebagai argumen build yang dapat ditimpa untuk aplikasi CGO.

## Diperbaiki

- Cegah crash saat menutup jendela yang masih memiliki panggilan asinkron tertunda dalam [PR](https://github.com/wailsapp/wails/pull/4435) oleh @leaanthony
- Cegah aktivasi jendela saat membuka aplikasi tersembunyi di Windows dalam [PR](https://github.com/wailsapp/wails/pull/5249) oleh @leaanthony
- Pastikan penanganan metadata permintaan WebKit, penyelesaian respons, dan aliran isi berjalan di thread utama GTK dalam [PR](https://github.com/wailsapp/wails/pull/5668) oleh @taliesin-ai
- Perbaiki `Menu.Update()` yang tidak membangun ulang menu native di GTK4 Linux (#5659, didiagnosis dan diperbaiki secara independen oleh @puneetdixit200 dalam #5539)
- Perbaiki crash saat menginventarisasi layar macOS ketika tampilan berubah dengan menyalin string id/nama layar dan membuat snapshot jumlahnya (#5565, didiagnosis dan diperbaiki secara independen oleh @x-haose dalam #5584)
- Perbaiki konten WebView2 yang menyusut lalu menghilang setelah jendela diseret melintasi monitor dengan DPI berbeda di Windows dengan menetapkan ulang batas pengontrol dalam handler `WM_DPICHANGED`, mengikuti sinkronisasi ulang DPI saat keluar dari keadaan diminimalkan (#5677)

## v3.0.0-alpha2.108 - 2026-06-28

## Ditambahkan

- Tambahkan pintasan papan ketik global (berlaku di seluruh sistem) melalui `app.GlobalShortcut` (`Register`, `Unregister`, `UnregisterAll`, `IsRegistered`, `GetAll`). Pintasan tetap dipicu meskipun aplikasi tidak sedang memiliki fokus. Diimplementasikan secara native untuk setiap platform tanpa dependensi pihak ketiga: hot key Carbon di macOS, `RegisterHotKey` di Windows, `XGrabKey` di X11, dan antarmuka pintasan global XDG Desktop Portal di Wayland.
- Tambahkan server MCP bawaan: server Model Context Protocol yang dimulai secara otomatis saat aplikasi dibangun dengan tag `mcp`, sehingga agen LLM dapat menguji dan mengendalikan aplikasi Wails yang sedang berjalan—mengendalikan jendela, memeriksa DOM, mengevaluasi JavaScript, memanggil metode terikat, menangani peristiwa, serta menyimulasikan masukan tetikus/papan ketik yang dirender dengan kursor animasi di layar. Tidak memerlukan kode pengguna: tag `mcp` ditambahkan secara otomatis oleh `wails3 build`/`wails3 dev` ketika `WAILS_MCP=1` ditetapkan. Seluruh konfigurasi dilakukan melalui variabel lingkungan (`WAILS_MCP_HOST`, `WAILS_MCP_PORT`, `WAILS_MCP_TIMEOUT`, `WAILS_MCP_HIDE_CURSOR`).

## Diperbaiki

- Perbaiki `Menu.Update()` yang tidak membangun ulang menu native di GTK4 Linux (#5659, didiagnosis dan diperbaiki secara independen oleh @puneetdixit200 dalam #5539)
- Perbaiki crash saat menginventarisasi layar macOS ketika tampilan berubah dengan menyalin string ID/nama layar dan membuat snapshot jumlah layar (#5565, didiagnosis dan diperbaiki secara independen oleh @x-haose dalam #5584)

## v3.0.0-alpha2.107 - 2026-06-27

## Ditambahkan

- Tambahkan dokumentasi Wake eksperimental dengan navigasi bilah sisi dalam [PR](https://github.com/wailsapp/wails/pull/5613) oleh @leaanthony

## v3.0.0-alpha2.106 - 2026-06-24

## Diubah

- Naikkan `webview2` ke v1.0.27.
  - ci(webview2): perbaiki build rilis (kompilasi silang Windows + go.sum lengkap) (#5671)\

  **Diff lengkap:** https://github.com/wailsapp/wails/compare/webview2/v1.0.26...webview2/v1.0.27

- Hapus go vet dari kompilasi silang alur kerja rilis webview2 dalam [PR](https://github.com/wailsapp/wails/pull/5672) oleh @taliesin-ai
- Perbarui model OpenRouter untuk changelog otomatis ke google/gemini-2.5-flash-lite dalam [PR](https://github.com/wailsapp/wails/pull/5670) oleh @taliesin-ai
- Naikkan `webview2` ke v1.0.26.

### Perbaikan

- **Pulihkan dari galat COM runtime sementara alih-alih keluar** (#5658, #5580). Sebelumnya, `Chromium.errorCallback` memanggil `os.Exit(1)` untuk *setiap* galat COM, sehingga gangguan yang dapat dipulihkan setelah startup menghentikan seluruh aplikasi. Jalur runtime (`Resize`/`GetClientRect`, `Navigate`/`NavigateToString`, `Init`, `MessageReceived`, `PutZoomFactor`, `OpenDevToolsWindow`) kini mencatat galat dan melakukan pemulihan. Secara khusus, pesan web yang rusak/tidak tepercaya dalam `MessageReceived` kini dibuang alih-alih menghentikan proses. Ini menangani kelas crash saat berpindah antar-monitor dengan DPI campuran (#5544, #5650). Jalur pembuatan lingkungan/pengontrol tetap bersifat fatal.\

**Diff lengkap:** https://github.com/wailsapp/wails/compare/webview2/v1.0.25...webview2/v1.0.26

## Diperbaiki

- Perbaiki alur kerja release-webview2 agar menangani file go.sum dengan benar dalam [PR](https://github.com/wailsapp/wails/pull/5671) oleh @taliesin-ai
- Perbaiki pembaruan menu GTK4 Linux dengan mengosongkan dan membangun ulang menu native dalam [PR](https://github.com/wailsapp/wails/pull/5659) oleh @taliesin-ai

## v3.0.0-alpha2.105 - 2026-06-21

## Ditambahkan

- Tambahkan `application.System` untuk mendeteksi platform saat runtime dari kode bersama: `System.IsMobile()` (iOS/Android), `System.IsDesktop()` (macOS/Windows/Linux), `System.IsServer()` (tag build `server`), dan `System.IsPlatform(application.PlatformMacOS|PlatformWindows|PlatformLinux|PlatformIOS|PlatformAndroid|PlatformServer)` untuk menguji satu target secara langsung. Kode ini dapat dikompilasi pada setiap target sehingga Anda dapat membuat percabangan tanpa tag build. Helper frontend yang sesuai (`System.IsMobile/IsDesktop/IsIOS/IsAndroid/...`) tersedia dalam `@wailsio/runtime`
- Tambahkan panduan "Menggunakan Framework Frontend Lain" yang menunjukkan cara menempatkan proyek Vite Anda sendiri ke dalam `frontend/` (mencakup Solid, Preact, Lit, SvelteKit, Qwik, Angular, dan lain-lain)
- Wizard `wails3 setup` kini memeriksa toolchain seluler (iOS/Android)—Xcode dan runtime iOS Simulator, JDK, Android SDK/NDK, serta emulator—dengan instalasi sekali klik dan perbaikan konfigurasi shell yang dapat disalin jika berlaku
- Proyek yang dihasilkan menyertakan `frontend/.npmrc` yang menetapkan `minimum-release-age` selama 7 hari untuk mengurangi paparan terhadap paket yang baru diterbitkan (dan berpotensi disusupi) (dipatuhi oleh pnpm dan bun; diabaikan tanpa dampak oleh npm)

## Diubah

- Desain ulang semua templat awal bawaan dengan tampilan hero pegunungan neon baru (web, iOS, dan Android)
- **TypeScript kini menjadi pilihan default untuk templat awal dan menggunakan nama templat tanpa akhiran.** `wails3 init` (tanpa `-t`) membuat kerangka proyek TypeScript; `-t vanilla`, `-t react`, `-t vue`, dan `-t svelte` menggunakan TypeScript, dengan varian JavaScript pada `-t vanilla-js`, `-t react-js`, `-t vue-js`, dan `-t svelte-js`. Templat bawaan mendeklarasikan bahasanya dengan `typescript:` dalam `template.yaml`; templat komunitas yang menggunakan akhiran `-ts` tetap berfungsi sebagai fallback
- Desain ulang wizard `wails3 setup` dengan tema neon "digital Wails" (efek semarak kaca buram di atas latar pegunungan)

## Diperbaiki

- Perbaiki crash di Windows saat memulihkan aplikasi yang diminimalkan cukup lama hingga WebView2 ditangguhkan atau proses render/GPU-nya didaur ulang. Sinkronisasi ulang DPI saat meminimalkan/memulihkan (#5544) kini hanya menyentuh pengontrol WebView2 ketika DPI jendela benar-benar berubah, sehingga menghindari panggilan COM fatal ke pengontrol yang ditangguhkan pada pemulihan dengan DPI yang sama, yang merupakan kasus umum (#5605)
- Perbaiki crash native `SIGABRT`/`SIGSEGV` yang berulang (biasanya di dalam `g_object_unref` selama loop utama GTK) pada aplikasi Linux yang berjalan lama dengan pemuatan aset/media yang sering. Server aset menyelesaikan `WebKitURISchemeRequest` dari goroutine pekerja, sehingga memanggil fungsi WebKit2GTK yang tidak aman terhadap thread di luar thread utama GTK; penyelesaian (`webkit_uri_scheme_request_finish_with_response`/`finish_error`) kini berjalan di thread utama. Melengkapi perbaikan parsial dalam #5566. Memengaruhi build GTK3 maupun GTK4/WebKitGTK 6.0 (#5631, #5557)
- Perbaiki `fatal error: invalid pointer found on stack` yang terjadi sesekali dalam `setupSignalHandlers` pada Linux/GTK3. ID jendela yang diteruskan sebagai `user_data` sinyal disimpan dalam variabel lokal Go `unsafe.Pointer`, sehingga garbage collector membatalkan proses ketika memindai nilai (non-pointer) tersebut selama penyalinan stack. ID kini dipertahankan sebagai tipe integer (`uintptr_t`) di sisi Go, dengan mem-backport perbaikan yang sama dari #4958 pada jalur GTK4 (yang mengubah fungsi sinyal C menjadi `uintptr_t` untuk menghilangkan galat `-race`/checkptr) ke jalur GTK3 lama (#5631)

## Dihapus

- Hapus templat awal `react-swc`, `preact`, `lit`, `solid`, `qwik`, dan `sveltekit` (beserta varian `-ts`-nya). Kumpulan bawaan yang didukung kini terdiri atas `vanilla`, `react`, `vue`, dan `svelte`—masing-masing menggunakan TypeScript secara default, dengan varian JavaScript `-js`. Framework lain tetap dapat digunakan dengan [membawa frontend Anda sendiri](https://v3.wails.io/guides/dev/frontend-frameworks) atau melalui templat khusus

## v3.0.0-alpha2.104 - 2026-06-18

## Diperbaiki

- Perbaiki crash iOS (SIGABRT) ketika metode layanan Go yang terikat mengembalikan string kosong. Penulis respons aset iOS memeriksa pointer body dengan `buf != nil` alih-alih panjangnya, sehingga body dengan panjang nol membuat `&buf[0]` mengalami panic; kini pemeriksaan dilakukan berdasarkan panjang, sesuai dengan penulis desktop

## v3.0.0-alpha2.103 - 2026-06-15

## Diubah

- Pindahkan fitur native iOS dan Android ke pengelola platform: panggil fitur tersebut melalui `application.IOS.*` dan `application.Android.*` (misalnya `application.IOS.Haptic("medium")`, `application.Android.Share(payload)`), bukan melalui fungsi bebas `application.IOS*`/`application.Android*` yang lama (#5602)
- Ganti nama peristiwa bridge seluler: peristiwa lintas platform kini menggunakan prefiks `common:*` (misalnya, `common:haptic`, `common:location`), sedangkan peristiwa khusus platform menggunakan `ios:*` / `android:*` (misalnya, `ios:backgroundTask`, `android:foregroundService`); prefiks `native:*` tidak lagi digunakan (#5602)

## v3.0.0-alpha.102 - 2026-06-14

## Ditambahkan

- Tambahkan wizard eksperimental `wails3 setup` untuk penyiapan proyek interaktif dan pemeriksaan dependensi
- Tambahkan flag `--json` ke `wails3 doctor` untuk keluaran yang dapat dibaca mesin
- Tambahkan bagian status penandatanganan ke perintah `wails3 doctor`

## Diperbaiki

- Perbaiki deteksi npm di Linux agar memeriksa PATH selain pengelola paket

## v3.0.0-alpha.101 - 2026-06-13

## Ditambahkan

- iOS: dialog pesan native (UIAlertController) serta dialog untuk membuka file, beberapa file, dan direktori (UIDocumentPickerViewController); dialog penyimpanan mengembalikan galat eksplisit
- iOS: dukungan papan klip melalui UIPasteboard
- iOS: metrik layar nyata melalui UIScreen (titik, piksel, skala, area kerja safe area)
- iOS: build perangkat (`IOS_PLATFORM=device`), dukungan identitas penandatanganan kode / profil penyediaan / entitlement, pemaketan `.ipa`, dan `deploy-device` melalui devicectl
- iOS: versi minimum iOS yang dapat dikonfigurasi (`ios.minIOSVersion` dalam build/config.yml)
- iOS: `wails3 doctor` melaporkan ketersediaan Xcode dan iOS SDK di macOS
- iOS: peristiwa sistem — baterai, jaringan, tema, penguncian layar, dan memori rendah diekspos sebagai peristiwa aplikasi `events.IOS.*` dan `events.Common.*` yang netral terhadap platform
- iOS: bridge fitur seluler native (`application.IOS*` yang diekspor) — lembar berbagi, membuka URL, menjaga perangkat tetap aktif, senter, inset safe area, kecerahan, info aplikasi, penguncian orientasi, bilah status, biometrik (Face ID/Touch ID), notifikasi lokal, dan penyimpanan aman Keychain
- iOS: sensor & perangkat keras — haptik, geolokasi sekali jalan, akselerometer, proksimitas, text-to-speech, info penyimpanan, status daya/baterai, status jaringan, inset papan ketik, dan deteksi tangkapan layar
- iOS: dokumentasi (IOS.md dan panduan situs dokumentasi)
- Android: dialog pesan native (AlertDialog) serta dialog untuk membuka file dan beberapa file (Storage Access Framework, diimpor sebagai salinan cache); dialog untuk membuka direktori dan menyimpan mengembalikan galat eksplisit
- Android: dukungan papan klip melalui ClipboardManager
- Android: metrik layar nyata melalui WindowMetrics/DisplayMetrics (dp, piksel, skala, area kerja bilah sistem)
- Android: metode runtime untuk haptik (`Android.Haptics.Vibrate`), info perangkat (`Android.Device.Info`), dan toast (`Android.Toast.Show`)
- Android: peristiwa siklus hidup bertipe (`events.Android.*`, dihasilkan dari events.txt), dengan `ActivityCreated` dipetakan ke `Common.ApplicationStarted`
- Android: pipeline build menghasilkan APK debug dan rilis yang dapat diinstal (`android:run`, `android:package`, `android:package:fat`); penandatanganan rilis secara default menggunakan debug keystore atau keystore sebenarnya melalui variabel lingkungan `ANDROID_KEYSTORE_*`
- Android: `wails3 doctor` melaporkan Android SDK, NDK, dan JDK
- Android: peristiwa sistem — baterai, jaringan, tema, penguncian layar, dan memori rendah diekspos sebagai peristiwa aplikasi `events.Android.*` dan `events.Common.*` yang netral terhadap platform
- Android: bridge fitur seluler native (`application.Android*` yang diekspor) — berbagi, membuka URL, menjaga perangkat tetap aktif, senter, inset safe area, kecerahan, info aplikasi, penguncian orientasi, bilah status, biometrik (BiometricPrompt), notifikasi lokal, dan penyimpanan aman EncryptedSharedPreferences
- Android: sensor & perangkat keras — haptik, geolokasi sekali jalan, akselerometer, proksimitas, text-to-speech, info penyimpanan, status daya/baterai, status jaringan, inset papan ketik, dan pemblokiran tangkapan layar dengan FLAG_SECURE
- Android: dokumentasi (ANDROID.md dan panduan situs dokumentasi)
- Contoh: kitchen sink `mobile` kini memiliki tab Seluler dan Perangkat Keras yang mendemonstrasikan bridge fitur native di iOS dan Android (tab berbentuk pil dibungkus menjadi beberapa baris)
- Seluler: baterai — akselerometer, sensor proksimitas, senter, dan jam periodik dalam contoh dijeda saat aplikasi berada di latar belakang lalu dipulihkan ketika aplikasi kembali (Android mempertahankan proses tetap berjalan di latar belakang, sedangkan senter merupakan status perangkat keras yang tetap bertahan di iOS), dan penerima peristiwa sistem Android hanya didaftarkan saat aplikasi berada di latar depan
- iOS: pengambilan gambar kamera — `application.IOSCapturePhoto`/`IOSCaptureVideo` (UIImagePickerController → peristiwa `native:capture` dengan thumbnail base64)
- iOS: eksekusi latar belakang — `application.IOSBeginBackgroundTask`/`IOSEndBackgroundTask` (jendela tugas latar belakang UIApplication) dan `ios.backgroundModes` yang dapat dikonfigurasi (build/config.yml), yang memasukkan `UIBackgroundModes` sebagai templat ke dalam Info.plist yang dihasilkan
- Android: pengambilan gambar kamera — `application.AndroidCapturePhoto`/`AndroidCaptureVideo` (kamera sistem melalui FileProvider → peristiwa `native:capture`)
- Android: layanan latar depan — `application.AndroidStartForegroundService`/`AndroidStopForegroundService` (`WailsForegroundService` dengan notifikasi berkelanjutan mempertahankan proses tetap hidup untuk pekerjaan latar belakang yang berjalan lama)
- Contoh: tab Kamera yang mendemonstrasikan pengambilan foto/video dan eksekusi latar belakang (layanan latar depan di Android, jendela tugas latar belakang di iOS)

## Diperbaiki

- Perbaiki `getUserMedia` yang selalu gagal dengan `NotAllowedError` di Linux: WebKitGTK menolak permintaan izin yang tidak ditangani oleh siapa pun, dan sinyal `permission-request` tidak terhubung. Kamera/mikrofon kini ditangani berdasarkan peta lintas platform `WebviewWindowOptions.Permissions` baru (`map[PermissionType]Permission`), yang dipatuhi di Linux (WebKitGTK) maupun Windows (WebView2). Di Linux, yang tidak memiliki prompt native, kamera/mikrofon secara default diizinkan (memulihkan `getUserMedia`) dan dapat dinonaktifkan dengan `PermissionDeny` (#5552)
- iOS: `GOOS=ios` dapat dikompilasi kembali (`events.IOS` yang diekspor, stub nama metode seluler), dan build dengan tag production dapat dikompilasi (perbaikan tag build di pkg/application dan beberapa layanan)
- iOS: peristiwa Go→JS dan ExecJS kini berfungsi — halaman tidak lagi dimuat dua kali saat startup dan handshake `wails:runtime:ready` tidak dapat lagi hilang
- iOS: `ApplicationDidFinishLaunching`/`ApplicationStarted` tidak lagi mengalami kondisi balapan dengan startup aplikasi; jeda startup tetap selama 2 detik telah dihapus
- iOS: perbaiki kebocoran string C pada setiap eksekusi JavaScript Go→JS
- iOS: `hasListeners` kini mencerminkan pendaftaran listener yang sebenarnya
- iOS: logging debug framework tidak disertakan saat build production dikompilasi
- Android: `GOOS=android` dapat dikompilasi kembali — mendefinisikan `events.Android`, menghapus array listener `events_android.go` yang berada di luar batas, menambahkan stub nama metode seluler, dan mencegah file Linux desktop (`linux_cgo.*`, `events_linux.*`, `environment_linux.go`) ikut masuk ke build Android
- Android: binding JS→Go kini berfungsi — WebView tidak dapat mengirimkan body POST `fetch()` ke `shouldInterceptRequest`, sehingga panggilan runtime dirutekan melalui transport JavascriptInterface (`nativeHandleRuntimeCall`), bukan mengalami crash akibat body permintaan nil
- Android: panggilan runtime `Screens.*` mengembalikan data sebenarnya — ScreenManager kini diisi saat startup (sebelumnya tidak pernah dihubungkan, sehingga `GetAll` mengembalikan nil)
- Android: logging debug framework tidak disertakan saat build production dikompilasi dan, dalam build debug, diarahkan melalui logcat dengan tag `Wails`
- Android: registri `hasListeners` yang sebenarnya, penanganan referensi/pengecualian JNI, dan siklus hidup halaman yang hanya memuat sekali (tanpa navigasi ganda)
- Perbaiki kegagalan `wails3 generate bindings` dengan pesan "Akses ditolak" di Windows saat server pengembangan Vite berjalan, dengan menyinkronkan file yang dihasilkan ke direktori keluaran alih-alih mengganti namanya di atas direktori tersebut (#5515)
- Memperbaiki crash fatal yang terjadi secara sporadis di macOS saat membaca informasi layar setelah perubahan tampilan: ID dan nama layar menyimpan pointer ke buffer `UTF8String` yang dirilis otomatis dan dapat dibebaskan sebelum disalin oleh Go (use-after-free). String tersebut kini di-`strdup` dan dibebaskan setelah konversi, sedangkan enumerasi layar dijalankan dalam autorelease pool eksplisit agar tidak lagi menyebabkan kebocoran memori saat dipanggil dari goroutine Go (#5556)
- Memperbaiki SIGSEGV yang terjadi secara sporadis di Linux saat assetserver menutup `WebKitURISchemeRequest`: `g_object_unref` terakhir berjalan pada goroutine assetserver sehingga menyelesaikan finalisasi GObject WebKit di luar thread utama GTK. Operasi unref kini diteruskan ke konteks utama GTK melalui `g_main_context_invoke` (#5557)

## v3.0.0-alpha.100 - 2026-06-13

## Ditambahkan

- Memperluas `MacWebviewPreferences` dengan opsi konfigurasi WKWebView tambahan: `EnableAutoplayWithoutUserAction`, `AllowsAirPlayForMediaPlayback`, `AllowsMagnification`, `JavaScriptCanOpenWindowsAutomatically`, `MinimumFontSize`, dan `ApplicationNameForUserAgent` (#5549)

## Diperbaiki

- Memperbaiki kegagalan `wails3 generate bindings` dengan pesan "Akses ditolak" di Windows saat server pengembangan Vite berjalan, dengan menyinkronkan berkas yang dihasilkan ke direktori keluaran alih-alih mengganti nama berkas lama dengan berkas tersebut (#5561)
- Memperbaiki peristiwa pengubahan ukuran JS yang tidak dipicu untuk jendela tanpa bingkai di Linux; memperbaiki deteksi tepi bilah gulir untuk jendela tanpa bingkai (#5368)
- Memperbaiki kegagalan updater di Windows dengan pesan "invalid cross-device link" ketika direktori sementara berada di volume yang berbeda dari direktori instalasi (#5560)

## v3.0.0-alpha.99 - 2026-06-10

## Diperbaiki

- Memperbaiki kegagalan `wails3 generate bindings` dengan pesan "Akses ditolak" di Windows saat server pengembangan Vite berjalan, dengan menyinkronkan berkas yang dihasilkan ke direktori keluaran alih-alih mengganti nama berkas lama dengan berkas tersebut (#5515)

## v3.0.0-alpha.98 - 2026-06-03

## Diperbaiki

- Memperbaiki UI WebKit yang macet di Linux saat tidak aktif (misalnya ketika inspector terbuka) dengan tidak lagi memaksakan `SA_ONSTACK` pada `SIGUSR1`, yang merusak sinkronisasi thread GC JavaScriptCore (#5527)

## v3.0.0-alpha.97 - 2026-05-31

## Ditambahkan

- Menambahkan halaman debugging dan panduan bekerja dengan `runtime/trace`

## Diubah

- Menghapus beberapa impor `_ "embed"` yang tidak diperlukan sehingga kode menjadi sedikit lebih rapi

## Diperbaiki

- Memperbaiki batasan lebar/tinggi minimum yang tidak diterapkan setelah jendela keluar dari keadaan dimaksimalkan di Windows (#4593)
- Memperbaiki klik tetikus yang menembus jendela dalam mode layar penuh saat menggunakan opsi jendela Frameless + Transparent (#4408)

## v3.0.0-alpha.96 - 2026-05-25

## Ditambahkan

- Menambahkan dukungan obfuscation Garble ([#4563](https://github.com/wailsapp/wails/issues/4563)): ID metode binding yang stabil, integrasi build/Taskfile (`build --obfuscated --garbleargs`, `generate bindings -obfuscated`), dan tag struct JSON pada setiap payload yang digunakan runtime (`EnvironmentInfo`, `OSInfo`, `Screen`, `Rect`, `Point`, `Size`, `Capabilities`) agar format wire tetap utuh ketika Garble mengganti nama field yang diekspor.

## v3.0.0-alpha.95 - 2026-05-20

## Ditambahkan

- Menambahkan halaman struktur proyek yang sebelumnya belum tersedia

## Diubah

- Dokumentasi: Mengubah beberapa diagram pada halaman arsitektur agar menggunakan diagram urutan sehingga tampilannya lebih rapi
- Dokumentasi: Menyertakan catatan tentang penginstalan D2 sebagai prasyarat untuk menjalankan

## Diperbaiki

- Memperbaiki `wails3 generate appimage` pada GTK4 default: bundler kini mendeteksi stack GTK dari biner sebelum mencari berkas runtime sehingga memilih `libwebkitgtkinjectedbundle.so` (di bawah `webkitgtk-6.0/`) untuk build GTK4 dan `libwebkit2gtkinjectedbundle.so` (di bawah `webkit2gtk-4.1/`) untuk build `-tags gtk3`. Pemeriksaan `.relr.dyn` juga memeriksa `libgtk-4.so.1` sehingga stripping dinonaktifkan dengan benar pada toolchain modern apa pun stack-nya. (#5475)
- Memperbaiki kegagalan `wails3 generate appimage` ketika dipanggil dengan `-builddir` relatif: bundler kini terlebih dahulu mengubah `-binary`, `-icon`, `-desktopfile`, `-builddir`, dan `-outputdir` menjadi path absolut agar `s.CD` di tengah alur tidak mengganggu goroutine pengunduhan AppRun maupun pemeriksaan `ldd` setelah penyalinan.
- Memperbaiki kegagalan `wails3 generate appimage` memindahkan AppImage akhir ke `-outputdir` ketika field `Name=` desktop tidak cocok dengan nama dasar biner: bundler kini memaksa plugin appimage milik linuxdeploy (melalui variabel lingkungan `OUTPUT`) menulis AppImage ke `<binary>-<arch>.AppImage`, bukan ke nama yang diturunkan dari berkas desktop.
- Memperbaiki `events.Common.ApplicationStarted`, `Common.ThemeChanged`, `Common.SystemWillSleep`, dan `Common.SystemDidWake` yang tidak dipicu di Linux setelah stack GTK4 + WebKitGTK 6.0 dijadikan default pada alpha.93. `run()` `application_linux.go` default yang baru tidak memanggil `setupCommonEvents()` (yang meneruskan peristiwa `Linux.*` ke padanan `Common.*`-nya) ataupun `monitorPowerEvents()`. Helper pemantau daya DBus kini digunakan bersama oleh jalur build GTK3 dan GTK4 melalui `application_linux_dbus.go`. (#5474)

## v3.0.0-alpha.94 - 2026-05-19

## Diperbaiki

- Memperbaiki `events.Common.ApplicationStarted`, `Common.ThemeChanged`, `Common.SystemWillSleep`, dan `Common.SystemDidWake` yang tidak dipicu di Linux setelah stack GTK4 + WebKitGTK 6.0 dijadikan default pada alpha.93. `run()` `application_linux.go` default yang baru tidak memanggil `setupCommonEvents()` (yang meneruskan peristiwa `Linux.*` ke padanan `Common.*`-nya) ataupun `monitorPowerEvents()`. Helper pemantau daya DBus kini digunakan bersama oleh jalur build GTK3 dan GTK4 melalui `application_linux_dbus.go`. (#5474)

## v3.0.0-alpha.93 - 2026-05-17

## Ditambahkan

- Menambahkan `XDG_SESSION_TYPE` ke keluaran `wails3 doctor` di Linux oleh @leaanthony

## Diperbaiki

- Memperbaiki crash menu jendela di Wayland akibat appmenu-gtk-module mengakses jendela yang belum direalisasikan (#4769) oleh @leaanthony
- Memperbaiki crash aplikasi GTK ketika nama aplikasi mengandung karakter yang tidak valid (spasi, tanda kurung, dan sebagainya) oleh @leaanthony
- Memperbaiki galat "memori tidak cukup" saat menginisialisasi seret dan lepas di Windows (#4701) oleh @overlordtm
- Memperbaiki race condition pada penyimpanan callback mainthread yang menggunakan RLock yang keliru untuk menghapus entri map (Linux, macOS, iOS) (#4424) oleh @leaanthony
- Memperbaiki penanganan variabel saat meneruskan argumen baris perintah ke task. Variabel CLI yang ditentukan sebagai pasangan KEY=VALUE kini diinisialisasi dan diteruskan dengan benar di sepanjang eksekusi task.
- Memperbaiki konflik NSWindowZoomButton di macOS: `MaximiseButtonState` dan `FullscreenButtonState` kini menerapkan keadaan yang lebih ketat, baik saat startup maupun runtime; kedua setter tidak dapat lagi diam-diam menimpa satu sama lain (#5319)
- Memperbaiki sekumpulan bug lama yang telah ada pada jalur build GTK3 lawas (`-tags gtk3`) dan terungkap oleh CodeRabbit pada #5463: peluncuran melalui asosiasi berkas tidak lagi melewati handler startup; `getTheme` kini aman dari segi batas dan tipe; `appName` tidak lagi membebaskan memori milik GLib; `clipboardGet` tidak lagi membocorkan `gchar*` yang dikembalikan GTK; `Calloc` kini menggunakan pointer receiver (dan `NewCalloc` mengembalikan `*Calloc`) sehingga pool benar-benar melacak alokasi; `zoomOut` menggunakan nilai kebalikan dari `zoomInFactor`, bukan pengali negatif yang dibatasi menjadi 1.0; `execJS` menggunakan kembali world-name kosong yang telah dialokasikan sebelumnya, alih-alih membocorkan satu `C.CString("")` pada setiap pemanggilan; sebuah `fmt.Println` pengembangan telah dihapus dari `menuItem.setAccelerator`. Menyelesaikan #5465.
- Perbaiki kebocoran penerima nilai `Calloc` yang sama pada jalur build GTK4 default (`linux_cgo.go`): gunakan penerima pointer + `NewCalloc() *Calloc` agar alokasi `c.String(...)` per jendela benar-benar dilacak dan dibebaskan.

## v3.0.0-alpha.92 - 2026-05-15

## Ditambahkan

- Ubah Taskfile agar pengelola paket frontend yang digunakan dapat dikontrol melalui opsi `PACKAGE_MANAGER`
- Perkaya data templat dengan `{{.Opn}}` dan `{{.Cls}}` agar penulisan templat Taskfile lebih mudah diprediksi

## Diubah

- Ubah beberapa Taskfile yang ada agar menggunakan `{{.Opn}} and {{.Cls}}`

## Diperbaiki

- Perbaiki galat fatal runtime `concurrent map read and map write` di `linuxSystemTray` saat menu baki diperbarui ketika panel sedang membacanya.
- Gunakan `log` sebagai pengganti `fmt` untuk keluaran galat dan jejak tumpukan WebView2 agar pesan tidak hilang ketika aplikasi berjalan tanpa konsol yang terpasang di Windows.

## v3.0.0-alpha.91 - 2026-05-12

## Diubah

- Perbarui SVG sponsor dalam [PR](https://github.com/wailsapp/wails/pull/5414) oleh `@github-actions[bot]`
- **PERUBAHAN INKOMPATIBEL (macOS):** Normalkan sistem koordinat macOS agar `GetScreens`, `Position`, dan `SetPosition` menggunakan ruang yang sama—titik logis, Y ke bawah, dengan `(0,0)` di kiri atas layar utama. Ini sesuai dengan Windows, GTK, serta API publik Electron dan web. Layar yang secara fisik berada di atas layar utama kini melaporkan `Bounds.Y` negatif (sebelumnya positif), dan nilai `Position()`/`SetPosition()` kini menggunakan titik logis, bukan `points × primaryScale`. Konversi pulang-pergi `Position()` → `SetPosition()` tetap dipertahankan; nilai absolut yang dicatat oleh build alfa sebelumnya atau solusi sementara yang dihitung secara manual (misalnya mengalikan dengan `primaryScale` atau membalik Y terhadap tinggi layar) perlu diperbarui. Menyelesaikan [#5117](https://github.com/wailsapp/wails/issues/5117).

## Diperbaiki

- Validasi nama sinyal DBus dan panjang isi secara defensif untuk mencegah panic dalam [PR](https://github.com/wailsapp/wails/pull/5416) oleh @leaanthony
- Perbaiki masalah keamanan memori dalam penanganan menu GTK di Linux melalui [PR](https://github.com/wailsapp/wails/pull/5363) oleh @leaanthony
- Deteksi GPU NVIDIA dan nonaktifkan perender DMA-BUF di Linux melalui [PR](https://github.com/wailsapp/wails/pull/5295) oleh @leaanthony
- Perbaiki konversi Y lintas layar `SetPosition` di macOS: gunakan tinggi layar utama sebagai acuan global agar jendela ditempatkan pada posisi yang tepat di monitor yang bergeser secara vertikal dari layar utama dalam [#5117](https://github.com/wailsapp/wails/issues/5117)
- Perbaiki templat PR git agar mengarah ke URL umpan balik yang benar dalam [PR](https://github.com/wailsapp/wails/pull/5109) oleh @wayneforrest
- Perbaiki serangkaian crash `SetMenu` pada baki sistem Windows yang disebabkan oleh syscall `DestroyMenu` yang rusak dan meneruskan empat argumen alih-alih satu, sehingga setiap panggilan mengembalikan FALSE dan tidak membebaskan apa pun. Bebaskan juga handle HMENU dan HBITMAP (termasuk yang dialokasikan saat runtime melalui `MenuItem.SetBitmap`) ketika menu dibangun ulang, atur ulang peta kotak centang/radio yang kedaluwarsa di `Win32Menu.Update`, dan hapus panggilan `Update()` yang berlebihan di `systemtray.updateMenu` yang menggandakan alokasi. Aplikasi baki sistem yang berjalan lama tidak lagi membocorkan objek GDI/USER setiap kali menu dibangun ulang.

## v3.0.0-alpha.90 - 2026-05-11

## Ditambahkan

- Tambahkan nama aplikasi yang dapat dikonfigurasi untuk User-Agent WKWebView di macOS melalui [PR](https://github.com/wailsapp/wails/pull/5261) oleh @vinhvoit225
- Tambahkan dependensi tidak langsung github.com/coder/websocket ke contoh gin-service melalui [PR](https://github.com/wailsapp/wails/pull/5400) oleh @taliesin-ai
- Tambahkan dukungan perbandingan kesetaraan mendalam ke pengujian aset build melalui [PR](https://github.com/wailsapp/wails/pull/5402) oleh @leaanthony

## Diubah

- Satukan keluaran build ke dalam direktori aset melalui [PR](https://github.com/wailsapp/wails/pull/5401) oleh @taliesin-ai
- Perbarui SVG sponsor dalam [PR](https://github.com/wailsapp/wails/pull/5399) oleh `@github-actions[bot]`

## Diperbaiki

- Gunakan objek notifikasi untuk pesan instans tunggal macOS melalui [PR](https://github.com/wailsapp/wails/pull/5289) oleh @overlordtm
- Proses callback Windows secara berkelompok untuk mencegah hilangnya promise saat beban berat melalui [PR](https://github.com/wailsapp/wails/pull/5383) oleh @taliesin-ai

## v3.0.0-alpha.89 - 2026-05-10

## Ditambahkan

- Tambahkan job go<em>test</em>results untuk mengagregasi hasil pengujian Go melalui [PR](https://github.com/wailsapp/wails/pull/5316) oleh @leaanthony

## Diubah

- Bagi payload RPC berukuran besar secara kondisional menjadi permintaan POST berpotongan melalui [PR](https://github.com/wailsapp/wails/pull/5369) oleh @leaanthony
- Tingkatkan Vite dari 5.x.x ke 8.0.0 di semua templat frontend melalui [PR](https://github.com/wailsapp/wails/pull/5386) oleh @leaanthony
- Migrasikan konfigurasi port server pengembangan Vite ke variabel lingkungan melalui [PR](https://github.com/wailsapp/wails/pull/5365) oleh @leaanthony
- Konfigurasikan server pengembangan Vite agar melakukan binding ke 127.0.0.1 di semua templat melalui [PR](https://github.com/wailsapp/wails/pull/5361) oleh @leaanthony
- Perbarui SVG sponsor dalam [PR](https://github.com/wailsapp/wails/pull/5384) oleh `@github-actions[bot]`

## Diperbaiki

- Sanitasi stub templat Info.plist saat pembaruan build-assets melalui [PR](https://github.com/wailsapp/wails/pull/5312) oleh @leaanthony
- Perbaiki status kedaluwarsa pada menu macOS dengan menerapkan mutator item menu (`setMenuItemChecked()`, `setMenuItemLabel()`, `setMenuItemDisabled()`, `setMenuItemHidden()`, `setMenuItemTooltip()`) secara sinkron di thread utama, sehingga menghilangkan kondisi balapan `dispatch_async` yang menyebabkan menu merender status sebelumnya ketika dibuka kembali dengan cepat (#5002)
- Abaikan file `*_test.go` dalam mode pengembangan untuk mencegah build ulang yang tidak perlu melalui [PR](https://github.com/wailsapp/wails/pull/5203) oleh @leaanthony
- Cegah segfault Menu.Update() saat aplikasi tidak berjalan melalui [PR](https://github.com/wailsapp/wails/pull/5291) oleh @wucm667
- Gunakan lastSizeWParam untuk membatasi penggambaran ulang bilah menu di Windows melalui [PR](https://github.com/wailsapp/wails/pull/5382) oleh @taliesin-ai

## v3.0.0-alpha.88 - 2026-05-09

## Diubah

- Alihkan HiddenOnTaskbar agar menggunakan WS<em>EX</em>TOOLWINDOW melalui [PR](https://github.com/wailsapp/wails/pull/5371) oleh @leaanthony
- Susun ulang dependensi dan hapus direktif replace webview2 dalam go.mod melalui [PR](https://github.com/wailsapp/wails/pull/5370) oleh @atterpac
- Perbarui SVG sponsor dalam [PR](https://github.com/wailsapp/wails/pull/5358) oleh `@github-actions[bot]`

## Diperbaiki

- Hapus alias pengalihan generik dan satukan tipe kunci peta melalui [PR](https://github.com/wailsapp/wails/pull/5331) oleh @fbbdev

## Dihapus

- Hapus alur kerja PR-master, termasuk dokumentasi, pengujian Go, dan pengujian yang dilewati, melalui [PR](https://github.com/wailsapp/wails/pull/5377) oleh @leaanthony

## v3.0.0-alpha.87 - 2026-05-07

## Ditambahkan

- Tambahkan dokumentasi bahasa Korea untuk Wails v3 melalui [PR](https://github.com/wailsapp/wails/pull/5352) oleh @leaanthony
- Tambahkan dokumentasi bahasa Prancis untuk instalasi dan panduan mulai cepat melalui [PR](https://github.com/wailsapp/wails/pull/5354) oleh @leaanthony
- Tambahkan dokumentasi bahasa Portugis untuk panduan mulai cepat, konsep, dan komunitas melalui [PR](https://github.com/wailsapp/wails/pull/5355) oleh @leaanthony

## v3.0.0-alpha.86 - 2026-05-06

## Ditambahkan

- Menambahkan pelokalan dokumentasi bahasa Prancis dalam [PR](https://github.com/wailsapp/wails/pull/5328) oleh @leaanthony
- Menambahkan locale bahasa Jerman ke situs dokumentasi dalam [PR](https://github.com/wailsapp/wails/pull/5343) oleh @leaanthony

## Diubah

- Mendaftarkan semua locale terjemahan 8 dalam konfigurasi dokumentasi melalui [PR](https://github.com/wailsapp/wails/pull/5347) oleh @leaanthony
- Memperbarui berbagai berkas terkait Windows untuk WebView2 dalam [PR](https://github.com/wailsapp/wails/pull/5317) oleh @leaanthony

## Diperbaiki

- Memisahkan pengiriman dialog antara GTK3 dan GTK4 untuk Linux dalam [PR](https://github.com/wailsapp/wails/pull/5340) oleh @leaanthony
- Memastikan callback dialog dijalankan pada thread GTK sehingga memperbaiki segmentation fault, dalam [PR](https://github.com/wailsapp/wails/pull/5339) oleh @leaanthony

## v3.0.0-alpha.85 - 2026-05-05

## Ditambahkan

- Menambahkan URL templat PR ke repositori dalam [PR](https://github.com/wailsapp/wails/pull/5179) oleh @leaanthony
- Menambahkan dokumentasi bahasa Jerman untuk Wails v3 dalam [PR](https://github.com/wailsapp/wails/pull/5330) oleh @leaanthony

## v3.0.0-alpha.84 - 2026-05-03

## Ditambahkan

- Menambahkan opsi untuk mencegah tombol Escape keluar dari mode layar penuh di macOS dalam [PR](https://github.com/wailsapp/wails/pull/5307) oleh @leaanthony
- Menambahkan opsi untuk mencegah tombol Escape keluar dari mode layar penuh di macOS dalam [PR](https://github.com/wailsapp/wails/pull/5310) oleh @leaanthony
- Menambahkan dokumentasi etalase komunitas Pausa dalam [PR](https://github.com/wailsapp/wails/pull/5288) oleh @yuseferi

## Diubah

- Memperbarui SVG sponsor dalam [PR](https://github.com/wailsapp/wails/pull/5308) oleh `@github-actions[bot]`
- Memperbarui perintah pembuatan ikon agar menangani platform yang tidak didukung dalam [PR](https://github.com/wailsapp/wails/pull/5309) oleh @leaanthony
- Mengganti API layar penuh boolean dengan ButtonState tiga status dan mengimplementasikan binding platform dalam [PR](https://github.com/wailsapp/wails/pull/5224) oleh @leaanthony

## Diperbaiki

- Melindungi operasi fokus WebView2 dari status controller nil dalam [PR](https://github.com/wailsapp/wails/pull/5315) oleh @leaanthony
- Memperbarui alur kerja GitHub Actions agar mereferensikan cabang dasar PR dengan benar dalam [PR](https://github.com/wailsapp/wails/pull/5313) oleh @leaanthony
- Mengabaikan berkas `*_test.go` dalam mode pengembangan untuk mencegah build ulang yang tidak diperlukan, dalam [PR](https://github.com/wailsapp/wails/pull/5203) oleh @leaanthony
- Mencegah segmentation fault pada Menu.Update() ketika aplikasi tidak berjalan, dalam [PR](https://github.com/wailsapp/wails/pull/5291) oleh @wucm667

## v3.0.0-alpha.83 - 2026-05-02

## Ditambahkan

- Menambahkan flag InstallScope dan opsi build untuk instalasi per mesin/pengguna dalam [PR](https://github.com/wailsapp/wails/pull/5094) oleh @symball
- Menambahkan metode SetScreen tanpa operasi ke BrowserWindow untuk memenuhi interface Window, dalam [PR](https://github.com/wailsapp/wails/pull/5294) oleh @leaanthony

## Diperbaiki

- Mendeteksi GPU NVIDIA dan menonaktifkan renderer DMA-BUF di Linux dalam [PR](https://github.com/wailsapp/wails/pull/5295) oleh @leaanthony
- Memperbaiki templat PR Git agar mengarah ke URL umpan balik yang benar dalam [PR](https://github.com/wailsapp/wails/pull/5109) oleh @wayneforrest
- Memperbaiki serangkaian crash systray Windows `SetMenu` yang disebabkan oleh syscall `DestroyMenu` yang rusak karena meneruskan empat argumen alih-alih satu, sehingga setiap panggilan mengembalikan FALSE dan tidak membebaskan apa pun. Juga melepaskan handle HMENU dan HBITMAP (termasuk yang dialokasikan saat runtime melalui `MenuItem.SetBitmap`) saat menu dibuat ulang, mereset peta checkbox/radio usang dalam `Win32Menu.Update`, dan menghapus panggilan `Update()` redundan dalam `systemtray.updateMenu` yang menggandakan alokasi. Aplikasi systray yang berjalan lama tidak lagi membocorkan objek GDI/USER setiap kali menu dibuat ulang.

## v3.0.0-alpha.82 - 2026-05-01

## Diperbaiki

- Memperbaiki pembuatan berkas desktop agar menangani nama desktop dengan benar dalam [PR](https://github.com/wailsapp/wails/pull/5232) oleh @leaanthony

## v3.0.0-alpha.81 - 2026-04-30

## Diubah

- Menyesuaikan jadwal rilis nightly menjadi pukul 15:00 UTC dalam [PR](https://github.com/wailsapp/wails/pull/5286) oleh @leaanthony

## Diperbaiki

- Memperbaiki nilai Screen Bounds, WorkArea, dan Size yang menjadi separuh pada Mac Retina -  (#5168)

## v3.0.0-alpha.80 - 2026-04-29

## Diubah

- Memperbarui dependensi dokumentasi dan loader koleksi konten dalam [PR](https://github.com/wailsapp/wails/pull/5285) oleh @leaanthony

## v3.0.0-alpha.79 - 2026-04-29

## Ditambahkan

- Memberikan izin actions: write kepada job trigger-release dalam [PR](https://github.com/wailsapp/wails/pull/5270) oleh @leaanthony

## Diubah

- Menjadikan cabang master sebagai default tugas rilis dan memperbarui redaksi changelog dalam [PR](https://github.com/wailsapp/wails/pull/5283) oleh @leaanthony
- Memperbarui alur kerja auto-changelog agar menggunakan versi terbaru dalam [PR](https://github.com/wailsapp/wails/pull/5282) oleh @leaanthony
- Meningkatkan efisiensi alur kerja dengan menambahkan filter jalur dan menghapus alur kerja yang tidak lagi digunakan, dalam [PR](https://github.com/wailsapp/wails/pull/5280) oleh @leaanthony
- Memperbarui dokumentasi agar mereferensikan cabang master untuk tautan contoh dalam [PR](https://github.com/wailsapp/wails/pull/5274) oleh @leaanthony
- Memperbarui dokumentasi dan contoh untuk v3 dalam [PR](https://github.com/wailsapp/wails/pull/5272) oleh @leaanthony

## Diperbaiki

- Meningkatkan reverse proxy dengan logika percobaan ulang dan pemaksaan IPv4 untuk pengembangan dalam [PR](https://github.com/wailsapp/wails/pull/5265) oleh @AkagiYui
- Menulis ulang alur kerja pemicu changelog yang belum dirilis dalam [PR](https://github.com/wailsapp/wails/pull/5281) oleh @leaanthony

## Dihapus

- Menghapus skrip pengujian shell untuk berbagai keperluan pengujian dalam [PR](https://github.com/wailsapp/wails/pull/5267) oleh @leaanthony
- Menghapus alur kerja deployment dokumentasi v3-alpha dan record CNAME dalam [PR](https://github.com/wailsapp/wails/pull/5266) oleh @leaanthony

### Ditambahkan

- Menambahkan entri Perutean Frontend ke navigasi bilah samping dalam [PR](https://github.com/wailsapp/wails/pull/5196) oleh @leaanthony
- Menambahkan panduan perutean frontend beserta rekomendasi khusus framework dalam [PR](https://github.com/wailsapp/wails/pull/5185) oleh @leaanthony
- Menambahkan dukungan untuk sheet modal (macOS)
- Menaikkan versi ghw untuk dukungan perangkat Apple yang lebih baik oleh @leaanthony (#4977)
- Menambahkan metode `GetBadge` ke layanan dock
- Menambahkan flag `-tags` ke perintah `wails3 build` untuk meneruskan tag build Go khusus (misalnya, `wails3 build -tags gtk4`) (#4957)
- Menambahkan dokumentasi untuk pembuatan enum otomatis dalam generator binding, termasuk halaman khusus Enum dan navigasi bilah samping (#4972)
- Menambahkan flag `-tags` ke perintah `wails3 build` untuk meneruskan tag build Go khusus (misalnya, `wails3 build -tags gtk4`) (#4957)
- Tambahkan contoh Web API di `v3/examples/web-apis/` yang mendemonstrasikan API browser 41, termasuk Penyimpanan (localStorage, sessionStorage, IndexedDB, Cache API), Jaringan (Fetch, WebSocket, XMLHttpRequest, EventSource, Beacon), Media (Canvas, WebGL, Web Audio, MediaDevices, MediaRecorder, Speech Synthesis), Perangkat (Geolocation, Clipboard, Fullscreen, Device Orientation, Vibration, Gamepad), Performa (Performance API, Mutation Observer, Intersection/Resize Observer), UI (Web Components, Pointer Events, Selection, Dialog, Drag and Drop), dan lainnya
- Tambahkan contoh pemeriksa kompatibilitas WebView API (`v3/examples/webview-api-check/`) yang menguji lebih dari 200 API browser di berbagai platform
- Tambahkan paket `internal/libpath` untuk menemukan jalur pustaka native di Linux dengan pencarian paralel, caching, dan dukungan untuk Flatpak/Snap/Nix
- **WIP:** Tambahkan dukungan eksperimental WebKitGTK 6.0 / GTK4 untuk Linux, yang tersedia melalui `-tags gtk4` (GTK3/WebKit2GTK 4.1 tetap menjadi default)
- Catatan: Pada pengelola jendela tiling (misalnya Hyprland dan Sway), operasi Perkecil/Perbesar mungkin tidak bekerja seperti yang diharapkan karena pengelola jendela mengendalikan geometri jendela
- Tambahkan cara membuat **Handler Satu Kali** dalam dokumentasi **Mendengarkan Peristiwa di JavaScript** oleh @AbdelhadiSeddar
- Tambahkan opsi `UseApplicationMenu` ke `WebviewWindowOptions` agar jendela di Windows/Linux dapat mewarisi menu aplikasi yang ditetapkan melalui `app.Menu.Set()` oleh @leaanthony
- Tambahkan dukungan penggunaan file `.icon` (format Apple Icon Composer) untuk menghasilkan ikon Liquid Glass dan katalog aset (macOS) (#4934) oleh @wimaha
- Tambahkan mode server eksperimental untuk deployment headless/web (`-tags server`). Memungkinkan aplikasi Wails dijalankan sebagai server HTTP tanpa dependensi GUI native. Build dengan `wails3 task build:server`. Lihat `examples/server` untuk detailnya.
- Tambahkan paket `internal/libpath` untuk menemukan jalur pustaka native di Linux dengan pencarian paralel, caching, dan dukungan untuk Flatpak/Snap/Nix
- Tambahkan opsi `CollectionBehavior` ke `MacWindow` untuk mengendalikan perilaku jendela di berbagai macOS Spaces dan layar penuh (#4756) oleh @leaanthony
- Tambahkan pengujian unit untuk pkg/application oleh @leaanthony
- Tambahkan dukungan protokol khusus ke pemaketan MSIX oleh @leaanthony
- Tambahkan deteksi lingkungan desktop di Linux [PR #4797](https://github.com/wailsapp/wails/pull/4797)
- Tambahkan metode `Window.Print()` ke runtime JavaScript untuk memicu dialog cetak dari frontend (#4290) oleh @leaanthony
- Tambahkan `XDG_SESSION_TYPE` ke output `wails3 doctor` di Linux oleh @leaanthony
- Tambahkan peristiwa perubahan pemuatan WebKit2 tambahan untuk Linux: `WindowLoadStarted`, `WindowLoadRedirected`, `WindowLoadCommitted`, `WindowLoadFinished` (#3896) oleh @leaanthony
- Tambahkan `XDG_SESSION_TYPE` ke output `wails3 doctor` di Linux oleh @leaanthony
- Hasilkan file `.desktop` selama build Linux, bukan hanya saat pemaketan (#4575)
- Tambahkan dokumentasi dependensi runtime Linux dengan nama paket khusus distribusi dan contoh pemaketan nfpm (#4339) oleh @leaanthony
- Tambahkan informasi versi driver NVIDIA ke output `wails3 doctor` di Linux oleh @leaanthony
- Tambahkan origin ke handler pesan mentah oleh @APshenkin dalam [PR](https://github.com/wailsapp/wails/pull/4710)
- Tambahkan dukungan universal link untuk macOS oleh @APshenkin dalam [PR](https://github.com/wailsapp/wails/pull/4712)
- Refaktor lapisan transport binding oleh @APshenkin dalam [PR](https://github.com/wailsapp/wails/pull/4702)
- Tambahkan pengenal aria-label ke template helloworld agar aplikasi contoh dapat diuji dengan mudah oleh klien pengujian Appium, oleh @chinenual dalam [PR](https://github.com/wailsapp/wails/pull/4760)
- Tambahkan origin ke handler pesan mentah oleh @APshenkin dalam [PR](https://github.com/wailsapp/wails/pull/4710)
- Tambahkan dukungan universal link untuk macOS oleh @APshenkin dalam [PR](https://github.com/wailsapp/wails/pull/4712)
- Refaktor lapisan transport binding oleh @APshenkin dalam [PR](https://github.com/wailsapp/wails/pull/4702)
- Peristiwa Bertipe oleh @fbbdev dan @ianvs dalam [#4633](https://github.com/wailsapp/wails/pull/4633)
- Tambahkan contoh `systray-clock` yang menampilkan tray headless dengan pembaruan tooltip secara langsung (#4653).
- Tambahkan template Protokol NSIS untuk Windows oleh @Tolfx dalam #4510
- Tambahkan pengujian untuk build-assets oleh @Tolfx dalam #4510
- macOS: Tampilkan kontrol jendela native di bilah menu dalam [#4588](https://github.com/wailsapp/wails/pull/4588) oleh @nidib
- Tambahkan layanan Dock macOS untuk menyembunyikan/menampilkan ikon aplikasi di Dock oleh @popaprozac dalam [PR](https://github.com/wailsapp/wails/pull/4451)
- Tambahkan layanan Dock macOS untuk menyembunyikan/menampilkan ikon aplikasi di Dock oleh @popaprozac dalam [PR](https://github.com/wailsapp/wails/pull/4451)
- Tambahkan dukungan efek Liquid Glass native untuk macOS dengan NSGlassEffectView (macOS 15.0+) dan fallback NSVisualEffectView, termasuk opsi penyesuaian material yang lengkap, oleh @leaanthony dalam [#4534](https://github.com/wailsapp/wails/pull/4534)
- Sanitasi URL Browser oleh @leaanthony dalam [#4500](https://github.dev/wailsapp/wails/pull/4500). Berdasarkan [#4484](https://github.com/wailsapp/wails/pull/4484) oleh @APShenkin.
- Tambahkan Perlindungan Konten di Windows/Mac oleh [@leaanthony](https://github.com/leaanthony) berdasarkan karya asli [@Taiterbase](https://github.com/Taiterbase) dalam [PR](https://github.com/wailsapp/wails/pull/4241) ini
- Tambahkan dukungan untuk meneruskan variabel CLI ke perintah Task melalui alias `wails3 build` dan `wails3 package` (#4422) oleh @leaanthony dalam [PR](https://github.com/wailsapp/wails/pull/4488)
- Dukungan untuk dropzone dengan data elemen yang dijatuhkan sebagai sumber peristiwa oleh [@atterpac](https://github.com/atterpac) dalam [#4318](https://github.com/wailsapp/wails/pull/4318)
- Tambahkan `AdditionalLaunchArgs` ke opsi `WindowsWindow` agar argumen baris perintah tambahan dapat diteruskan ke browser WebView2 dalam [PR](https://github.com/wailsapp/wails/pull/4467)
- Tambahkan eksekusi go mod tidy secara otomatis setelah wails init oleh [@triadmoko](https://github.com/triadmoko) dalam [PR](https://github.com/wailsapp/wails/pull/4286)
- Fitur Snapassist Windows oleh @leaanthony dalam [PR](https://github.dev/wailsapp/wails/pull/4463)
- Tambahkan `AdditionalLaunchArgs` ke opsi `WindowsWindow` agar argumen baris perintah tambahan dapat diteruskan ke browser WebView2 dalam [PR](https://github.com/wailsapp/wails/pull/4467)
- Tambahkan eksekusi go mod tidy secara otomatis setelah wails init oleh [@triadmoko](https://github.com/triadmoko) dalam [PR](https://github.com/wailsapp/wails/pull/4286)
- Fitur Snapassist Windows oleh @leaanthony dalam [PR](https://github.dev/wailsapp/wails/pull/4463)
- Tambahkan implementasi `getAccentColor` Windows oleh [@almas-x](https://github.com/almas-x) dalam [PR](https://github.com/wailsapp/wails/pull/4427)
- Tambahkan implementasi `getAccentColor` Windows oleh [@almas-x](https://github.com/almas-x) dalam [PR](https://github.com/wailsapp/wails/pull/4427)
- Menu dan bilah menu bertema gelap di Windows. Oleh @leaanthony dalam [a29b4f0861b1d0a700e9eb213c6f1076ec40efd5](https://github.com/wailsapp/wails/commit/a29b4f0861b1d0a700e9eb213c6f1076ec40efd5)
- Ganti nama layanan bawaan agar binding JS/TS lebih jelas, oleh @popaprozac dalam [PR](https://github.com/wailsapp/wails/pull/4405)
- `app.Env.GetAccentColor` untuk mendapatkan warna aksen sistem pengguna. Berfungsi di MacOS. Oleh [@etesam913](https://github.com/etesam913)
- Tambahkan API `window.ToggleFrameless()` oleh [@atterpac](https://github.com/atterpac) dalam [#4137](https://github.com/wailsapp/wails/pull/4137)
- Tambahkan dependensi build khusus distribusi untuk Linux oleh @leaanthony dalam [PR](https://github.com/wailsapp/wails/pull/4345)
- Menambahkan panduan binding oleh @atterpac dalam [PR](https://github.com/wailsapp/wails/pull/4404)
- **Infrastruktur Pengujian yang Terorganisasi**: Memindahkan file pengujian Docker ke direktori khusus `test/docker/` dengan image yang dioptimalkan dan keandalan build yang ditingkatkan oleh [@leaanthony](https://github.com/leaanthony) dalam [#4359](https://github.com/wailsapp/wails/pull/4359)
- **Pola Pengelolaan Sumber Daya yang Ditingkatkan**: Menambahkan pembersihan event handler yang tepat dan pengelolaan goroutine yang peka terhadap konteks dalam contoh oleh [@leaanthony](https://github.com/leaanthony) dalam [#4359](https://github.com/wailsapp/wails/pull/4359)
- Mendukung build AppImage aarch64 oleh [@AkshayKalose](https://github.com/AkshayKalose) dalam [#3981](https://github.com/wailsapp/wails/pull/3981)
- Menambahkan bagian diagnostik ke `wails doctor` oleh [@leaanthony](https://github.com/leaanthony)
- Menambahkan jendela ke konteks saat memanggil metode layanan oleh [@leaanthony](https://github.com/leaanthony)
- Menambahkan contoh `window-call` untuk menunjukkan cara mengetahui jendela yang memanggil layanan oleh [@leaanthony](https://github.com/leaanthony)
- Panduan Menu baru oleh [@leaanthony](https://github.com/leaanthony)
- Penanganan panic yang lebih baik oleh [@leaanthony](https://github.com/leaanthony)
- Panduan Menu baru oleh [@leaanthony](https://github.com/leaanthony)
- Menambahkan komentar dokumentasi untuk Service API oleh [@fbbdev](https://github.com/fbbdev) dalam [#4024](https://github.com/wailsapp/wails/pull/4024)
- Menambahkan fungsi `application.NewServiceWithOptions` untuk menginisialisasi layanan dengan konfigurasi tambahan oleh [@leaanthony](https://github.com/leaanthony) dalam [#4024](https://github.com/wailsapp/wails/pull/4024)
- Meningkatkan kontrol menu oleh [@FalcoG](https://github.com/FalcoG) dan [@leaanthony](https://github.com/leaanthony) dalam [#4031](https://github.com/wailsapp/wails/pull/4031)
- Menambahkan dokumentasi oleh [@leaanthony](https://github.com/leaanthony)
- Mendukung pembatalan event dalam event listener standar oleh [@leaanthony](https://github.com/leaanthony)
- Dukungan `Hide`, `Show`, dan `Destroy` untuk systray oleh [@leaanthony](https://github.com/leaanthony)
- Dukungan `SetTooltip` untuk systray oleh [@leaanthony](https://github.com/leaanthony). Ide awal oleh [@lujihong](https://github.com/wailsapp/wails/issues/3487#issuecomment-2633242304)
- Melaporkan path paket dalam peringatan generator binding tentang tipe yang tidak didukung oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Menambahkan dukungan generator binding untuk alias generik oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Menambahkan dukungan generator binding untuk flag JSON `omitzero` oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Menambahkan direktif `//wails:ignore` untuk mencegah pembuatan binding bagi metode layanan yang dipilih oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Menambahkan direktif `//wails:internal` pada layanan dan model untuk mengizinkan tipe yang diekspor di Go tetapi tidak di JS/TS oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Menambahkan dukungan generator binding untuk konstanta bertipe alias agar dapat membuat enum dengan pengetikan lemah oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Menambahkan pengujian generator binding untuk fitur Go 1.24 oleh [@fbbdev](https://github.com/fbbdev) dalam [#4068](https://github.com/wailsapp/wails/pull/4068)
- Menambahkan dukungan macOS 15 "Sequoia" ke `OSInfo.Branding` untuk meningkatkan deteksi versi OS dalam [#4065](https://github.com/wailsapp/wails/pull/4065)
- Menambahkan hook `PostShutdown` untuk menjalankan kode kustom setelah proses shutdown selesai oleh [@fbbdev](https://github.com/fbbdev) dalam [#4066](https://github.com/wailsapp/wails/pull/4066)
- Menambahkan struct `FatalError` untuk mendukung deteksi error fatal dalam error handler kustom oleh [@fbbdev](https://github.com/fbbdev) dalam [#4066](https://github.com/wailsapp/wails/pull/4066)
- Menstandarkan dan mendokumentasikan urutan startup dan shutdown layanan oleh [@fbbdev](https://github.com/fbbdev) dalam [#4066](https://github.com/wailsapp/wails/pull/4066)
- Menambahkan test harness untuk urutan startup/shutdown aplikasi dan pengujian startup/shutdown layanan oleh [@fbbdev](https://github.com/fbbdev) dalam [#4066](https://github.com/wailsapp/wails/pull/4066)
- Menambahkan metode `RegisterService` untuk mendaftarkan layanan setelah aplikasi dibuat oleh [@fbbdev](https://github.com/fbbdev) dalam [#4066](https://github.com/wailsapp/wails/pull/4066)
- Menambahkan field `MarshalError` dalam opsi aplikasi dan layanan untuk penanganan error kustom dalam panggilan binding oleh [@fbbdev](https://github.com/fbbdev) dalam [#4066](https://github.com/wailsapp/wails/pull/4066)
- Menambahkan wrapper promise yang dapat dibatalkan dan meneruskan permintaan pembatalan melalui rantai promise oleh [@fbbdev](https://github.com/fbbdev) dalam [#4100](https://github.com/wailsapp/wails/pull/4100)
- Menambahkan kemampuan untuk mengaitkan pembatalan panggilan binding dengan `AbortSignal` oleh [@fbbdev](https://github.com/fbbdev) dalam [#4100](https://github.com/wailsapp/wails/pull/4100)
- Mendukung atribut `data-wml-*` untuk WML selain atribut `wml-*` yang biasa digunakan oleh [@leaanthony](https://github.com/leaanthony)
- Menambahkan metode `Configure` pada semua layanan untuk konfigurasi lanjutan/konfigurasi ulang dinamis oleh [@fbbdev](https://github.com/fbbdev) dalam [#4067](https://github.com/wailsapp/wails/pull/4067)
- Layanan `fileserver` mengirimkan respons 503 Service Unavailable ketika belum dikonfigurasi oleh [@fbbdev](https://github.com/fbbdev) dalam [#4067](https://github.com/wailsapp/wails/pull/4067)
- Layanan `kvstore` secara default menyediakan penyimpanan nilai-kunci dalam memori ketika belum dikonfigurasi oleh [@fbbdev](https://github.com/fbbdev) dalam [#4067](https://github.com/wailsapp/wails/pull/4067)
- Menambahkan metode `Load` pada layanan `kvstore` untuk memuat ulang data dari file setelah perubahan konfigurasi oleh [@fbbdev](https://github.com/fbbdev) dalam [#4067](https://github.com/wailsapp/wails/pull/4067)
- Menambahkan metode `Clear` pada layanan `kvstore` untuk menghapus semua kunci oleh [@fbbdev](https://github.com/fbbdev) dalam [#4067](https://github.com/wailsapp/wails/pull/4067)
- Menambahkan tipe `Level` dalam layanan `log` untuk menyediakan konstanta level log di sisi JS oleh [@fbbdev](https://github.com/fbbdev) dalam [#4067](https://github.com/wailsapp/wails/pull/4067)
- Menambahkan metode `Log` pada layanan `log` untuk menentukan level log secara dinamis oleh [@fbbdev](https://github.com/fbbdev) dalam [#4067](https://github.com/wailsapp/wails/pull/4067)
- Layanan `sqlite` secara default menyediakan DB dalam memori ketika belum dikonfigurasi oleh [@fbbdev](https://github.com/fbbdev) dalam [#4067](https://github.com/wailsapp/wails/pull/4067)
- Menambahkan metode `Close` pada layanan `sqlite` untuk menutup DB secara manual oleh [@fbbdev](https://github.com/fbbdev) dalam [#4067](https://github.com/wailsapp/wails/pull/4067)
- Menambahkan dukungan pembatalan untuk metode kueri pada layanan `sqlite` oleh [@fbbdev](https://github.com/fbbdev) dalam [#4067](https://github.com/wailsapp/wails/pull/4067)
- Menambahkan dukungan prepared statement ke layanan `sqlite` dengan binding JS oleh [@fbbdev](https://github.com/fbbdev) dalam [#4067](https://github.com/wailsapp/wails/pull/4067)
- Dukungan Gin oleh [Lea Anthony](https://github.com/leaanthony) dalam [PR](https://github.com/wailsapp/wails/pull/3537), berdasarkan karya awal [@AnalogJ](https://github.com/AnalogJ) dalam [PR](https://github.com/wailsapp/wails/pull/3537) ini
- Memperbaiki penyimpanan otomatis dan penyimpanan otomatis kata sandi yang selalu aktif oleh [@oSethoum](https://github.com/osethoum) dalam [#4134](https://github.com/wailsapp/wails/pull/4134)
- Menambahkan `SetMenu()` pada jendela agar menu dapat ditetapkan pada jendela oleh [@leaanthony](https://github.com/leaanthony)
- Menambahkan dukungan Notifikasi oleh [@popaprozac](https://github.com/popaprozac) dalam [#4098](https://github.com/wailsapp/wails/pull/4098)
-  Menambahkan dukungan Asosiasi File untuk Mac oleh [@wimaha](https://github.com/wimaha) dalam [#4177](https://github.com/wailsapp/wails/pull/4177)
- Menambahkan `wails3 tool version` untuk menaikkan versi semantik oleh [@leaanthony](https://github.com/leaanthony)
- Menambahkan dukungan lencana untuk macOS dan Windows oleh [@popaprozac](https://github.com/popaprozac) dalam [#](https://github.com/wailsapp/wails/pull/4234)
- Menambahkan dukungan untuk event terdaftar/bertipe ketat oleh [@fbbdev](https://github.com/fbbdev) dan [@IanVS](https://github.com/IanVS) dalam [#4161](https://github.com/wailsapp/wails/pull/4161)
- Menambahkan kemampuan untuk mendaftarkan hook bagi event khusus oleh [@fbbdev](https://github.com/fbbdev) dan [@IanVS](https://github.com/IanVS) dalam [#4161](https://github.com/wailsapp/wails/pull/4161)
- `app.OpenFileManager(path string, selectFile bool)` untuk membuka pengelola file sistem pada jalur `path` dengan penyorotan opsional melalui `selectFile` oleh [@Krzysztofz01](https://github.com/Krzysztofz01) [@rcalixte](https://github.com/rcalixte)
- Flag `-git` baru untuk perintah `wails3 init` oleh [@leaanthony](https://github.com/leaanthony)
- Perintah `wails3 generate webview2bootstrapper` baru oleh [@leaanthony](https://github.com/leaanthony)
- Menambahkan metode `init()` dalam runtime untuk memungkinkan inisialisasi runtime secara manual oleh [@leaanthony](https://github.com/leaanthony)
- Menambahkan opsi `WindowDidMoveDebounceMS` ke WindowOptions milik Window oleh [@leaanthony](https://github.com/leaanthony)
- Menambahkan fitur Instans Tunggal oleh [@leaanthony](https://github.com/leaanthony). Berdasarkan PR [v2](https://github.com/wailsapp/wails/pull/2951) oleh @APshenkin.
- Perintah `wails3 generate template` oleh [@leaanthony](https://github.com/leaanthony)
- Perintah `wails3 releasenotes` oleh [@leaanthony](https://github.com/leaanthony)
- Perintah `wails3 update cli` oleh [@leaanthony](https://github.com/leaanthony)
- Opsi `-clean` untuk perintah `wails3 generate bindings` oleh [@leaanthony](https://github.com/leaanthony)
- Memungkinkan build AppImage Linux aarch64 (arm64) oleh [@AkshayKalose](https://github.com/AkshayKalose) dalam [#3981](https://github.com/wailsapp/wails/pull/3981)
- Menambahkan hyperlink untuk sponsor oleh @ansxuman dalam [#3958](https://github.com/wailsapp/wails/pull/3958)
- Dukungan pemaketan Linux untuk build pemaket deb, rpm, dan Arch Linux oleh
- Menambahkan dukungan untuk build dan paket universal Darwin oleh
- Dokumentasi event ke situs web oleh
- Templat untuk sveltekit dan sveltekit-ts yang dikonfigurasi untuk pengembangan non-SSR
- Memperbarui aset build menggunakan perintah `wails3 update build-assets` baru oleh
- Contoh untuk menguji HTML Drag and Drop API oleh
- Dukungan Asosiasi File oleh [leaanthony](https://github.com/leaanthony) dalam
- Perintah `wails3 generate runtime` baru oleh
- Opsi `InitialPosition` baru untuk menentukan apakah jendela harus dipusatkan atau
- Menambahkan metode `Path` & `Paths` ke paket `application` oleh
- Menambahkan opsi Windows `GeneralAutofillEnabled` dan `PasswordAutosaveEnabled`
- Menambahkan kemampuan untuk mengambil jendela yang memanggil metode layanan oleh
- Menambahkan opsi `EnabledFeatures` dan `DisabledFeatures` untuk WebView2 oleh
- ⊞ Sistem DIP baru untuk Dukungan Monitor DPI Tinggi yang Disempurnakan oleh
- ⊞ Opsi nama kelas jendela oleh [windom](https://github.com/windom/) dalam
- Layanan telah diperluas untuk menyediakan fungsionalitas plugin. Oleh
- 🐧 Event WindowDidMove / WindowDidResize dalam
- ⊞ Event WindowDidResize dalam
-  menambahkan event ApplicationShouldHandleReopen agar dock dapat ditangani
-  menambahkan getPrimaryScreen/getScreens ke implementasi oleh @tmclane dalam
-  menambahkan opsi untuk menampilkan bilah alat dalam mode layar penuh di macOS oleh
- 🐧 menambahkan logika onKeyPress untuk mengonversi penekanan tombol Linux menjadi akselerator
- 🐧 menambahkan tugas `run:linux` oleh
- Mengekspor metode `SetIcon` oleh [@almas-x](https://github.com/almas-x) dalam
- Meningkatkan `OnShutdown` oleh [@almas-x](https://github.com/almas-x) dalam
- Memulihkan metode `ToggleMaximise` dalam antarmuka `Window` oleh
- Menambahkan lebih banyak informasi ke `Environment()`. Oleh @leaanthony dalam
- Mengekspos metode `WebviewWindow.IsFocused` pada antarmuka `Window` oleh
- Mendukung beberapa event pemicu yang dipisahkan spasi dalam sistem WML oleh
- Menambahkan ekspor ESM dari skrip runtime JS yang dibundel oleh
- Menambahkan flag generator binding untuk menggunakan skrip runtime JS yang dibundel sebagai pengganti
- Mengimplementasikan `setIcon` di Linux oleh [@abichinger](https://github.com/abichinger)
- Menambahkan flag `-port` ke perintah dev dan mendukung variabel lingkungan
- Menambahkan pengujian untuk pemanggilan metode terikat oleh
- ⊞ menambahkan `SetIgnoreMouseEvents` untuk jendela yang sudah dibuat oleh
-  Menambahkan kemampuan untuk menetapkan tingkat penumpukan (urutan) jendela oleh

### Diperbaiki

- Memperbaiki `Screen.Bounds`, `WorkArea`, dan `Size` yang menjadi setengah pada Mac Retina dengan mengonversi nilai titik NSScreen menjadi piksel perangkat dalam kolom `Physical*`, serta mengisi `Screen.X`/`Y` tingkat atas agar deteksi persinggungan antarmonitor dan penempatan area kerja menjadi benar dalam [PR](https://github.com/wailsapp/wails/pull/5168) oleh @wayneforrest
- Memperbaiki data race di ScreenManager yang menyebabkan deadlock WebKit DisplayLink saat konfigurasi layar berubah (misalnya monitor eksternal disambungkan saat siklus tidur/bangun)
- Menetapkan CFBundleIconName secara langsung ke appicon ketika Assets.car ada dalam [PR](https://github.com/wailsapp/wails/pull/5154) oleh @symball
- Memperbaiki `wails3 doctor` yang melaporkan paket WebKitGTK yang salah di Fedora, openSUSE, Arch, dan NixOS — entri fallback 4.0 telah dihapus karena v3 memerlukan API 4.1 pada waktu kompilasi (#5071)
- Memperbaiki nama paket doctor webkit2gtk openSUSE (`webkit2gtk4_1-devel` → `webkit2gtk3-devel`, nama paket openSUSE yang benar) (#5071)
- Perbaiki kesalahan `Unexpected token '<'` saat `/wails/custom.js` tidak ada dalam mode pengembangan desktop. Menambahkan penangan eksplisit 404 untuk `/wails/custom.js` dan validasi `Content-Type` yang tidak peka huruf besar-kecil di `loadOptionalScript` guna mencegah fallback SPA HTML disisipkan sebagai JavaScript. ([#5068](https://github.com/wailsapp/wails/issues/5068))
- Perbaiki status sorotan menu baki sistem di macOS—ikon kini menampilkan status terpilih saat menu terbuka (#4910)
- Perbaiki jendela yang terpasang pada baki sistem dan muncul di belakang jendela lain di macOS—kini menggunakan level jendela popup yang sesuai (#4910)
- Perbaiki contoh impor `@wailsio/runtime` yang keliru di seluruh dokumentasi (#4989)
- Perbaiki jendela tanpa bingkai yang tidak dapat diminimalkan di darwin (#4294)
- Perbaiki macet selama 20-30 menit saat `wails3 build` dan `wails3 dev` dengan mengecualikan `node_modules/` dari pemeriksaan kemutakhiran go-task. Sebelumnya, glob `sources: "**/*"` menyebabkan go-task menghitung dan memeriksa checksum setiap berkas di `node_modules/` (lebih dari 50000-100000 berkas dengan dependensi berat seperti MUI), yang sangat lambat khususnya di Windows/NTFS (#4939)
- Perbaiki kegagalan build GTK4 akibat typedef C `Screen` yang bertabrakan dengan X11 Xlib.h (#4957)
- Perbaiki konsistensi metode lencana dock di macOS
- Perbaiki `InvisibleTitleBarHeight` yang diterapkan ke semua jendela macOS, bukan hanya jendela tanpa bingkai atau dengan bilah judul transparan (#4960)
- Perbaiki jendela yang berguncang/bergetar saat ukurannya diubah dari sudut atas dengan `InvisibleTitleBarHeight` diaktifkan, dengan melewati inisiasi penyeretan di dekat tepi jendela (#4960)
- Perbaiki pembuatan tipe terpetakan dengan kunci enum dalam binding JS/TS (#4437) oleh @fbbdev
- Perbaiki seret dan lepas berkas di Windows yang tidak berfungsi pada skala tampilan selain 100%
- Perbaiki seret dan lepas internal HTML5 yang rusak saat fitur pelepasan berkas diaktifkan di Windows
- Perbaiki koordinat pelepasan berkas yang menggunakan ruang piksel yang salah di Windows (piksel fisik dibandingkan piksel CSS)
- Perbaiki seret dan lepas berkas di Linux yang tidak berfungsi secara andal dengan efek hover
- Perbaiki seret dan lepas internal HTML5 yang rusak saat fitur pelepasan berkas diaktifkan di Linux
- Perbaiki tindakan menampilkan/menyembunyikan jendela di Linux/GTK4 yang terkadang memulihkannya ke status diminimalkan dengan menggunakan `gtk_window_present()` (#4957)
- Perbaiki operasi mengambil/mengatur posisi jendela di Linux/GTK4 yang selalu mengembalikan 0,0 dengan menambahkan dukungan bersyarat X11 melalui `XTranslateCoordinates`/`XMoveWindow` (#4957)
- Perbaiki batas ukuran maksimum jendela yang tidak diterapkan di Linux/GTK4 dengan menambahkan pembatasan ukuran berbasis sinyal untuk menggantikan `gtk_window_set_geometry_hints` yang telah dihapus (#4957)
- Perbaiki penskalaan DPI di Linux/GTK4 dengan menerapkan penghitungan PhysicalBounds yang tepat dan dukungan penskalaan pecahan melalui `gdk_monitor_get_scale` (GTK 4.14+)
- Perbaiki item menu yang terduplikasi saat membuat jendela baru di Linux/GTK4
- Perbaiki pembuatan tipe terpetakan dengan kunci enum dalam binding JS/TS (#4437) oleh @fbbdev
- Perbaiki seret dan lepas berkas di Windows yang tidak berfungsi pada skala tampilan selain 100%
- Perbaiki seret dan lepas internal HTML5 yang rusak saat fitur pelepasan berkas diaktifkan di Windows
- Perbaiki koordinat pelepasan berkas yang menggunakan ruang piksel yang salah di Windows (piksel fisik dibandingkan piksel CSS)
- Perbaiki seret dan lepas berkas di Linux yang tidak berfungsi secara andal dengan efek hover
- Perbaiki seret dan lepas internal HTML5 yang rusak saat fitur pelepasan berkas diaktifkan di Linux
- Perbaiki penskalaan DPI di Linux/GTK4 dengan menerapkan penghitungan PhysicalBounds yang tepat dan dukungan penskalaan pecahan melalui `gdk_monitor_get_scale` (GTK 4.14+)
- Perbaiki item menu yang terduplikasi saat membuat jendela baru di Linux/GTK4
- Perbaiki pembuatan tipe terpetakan dengan kunci enum dalam binding JS/TS (#4437) oleh @fbbdev
- Perbaiki masalah "jendela hantu" di macOS akibat API AppKit tidak diakses dari Main Thread dalam App.Window.Current() (#4947) oleh @wimaha
- Perbaiki `<input type="file">` HTML yang tidak berfungsi di macOS dengan menerapkan WKUIDelegate runOpenPanelWithParameters (#4862)
- Perbaiki seret dan lepas berkas native yang tidak berfungsi saat menggunakan modul npm `@wailsio/runtime` di macOS/Linux (#4953) oleh @leaanthony
- Perbaiki pembuatan binding untuk alias tipe lintas paket (#4578) oleh @fbbdev
- Perbaiki crash OpenFileDialog di Linux akibat pelanggaran keamanan thread GTK (#3683) oleh @ddmoney420
- Perbaiki crash SIGSEGV saat memanggil `Focus()` pada jendela yang disembunyikan atau telah dihancurkan (#4890) oleh @ddmoney420
- Perbaiki potensi panic saat mengatur ikon atau bitmap kosong di Linux (#4923) oleh @ddmoney420
- Perbaiki crash ErrorDialog saat dipanggil dari binding layanan di macOS (#3631) oleh @leaanthony
- Tampilkan menu di sistem operasi Windows dalam `v3\examples\dialogs` oleh @ndianabasi
- Perbaiki kondisi race yang menyebabkan TypeError saat halaman dimuat ulang (#4872) oleh @ddmoney420
- Perbaiki keluaran yang keliru dari pengujian generator binding dengan menghapus status global dalam metode `Collector.IsVoidAlias()` (#4941) oleh @fbbdev
- Perbaiki pemilih berkas `<input type="file">` yang tidak berfungsi di macOS (#4862) oleh @leaanthony
- Perbaiki `Position()` dan `SetPosition()` yang menggunakan sistem koordinat tidak konsisten di macOS sehingga posisi jendela bergeser saat status disimpan/dipulihkan (#4816) oleh @leaanthony
- Perbaiki kesalahan SetProcessDpiAwarenessContext "Access is denied" ketika kesadaran DPI sudah diatur melalui manifes aplikasi (#4803)
- Perbarui halaman dokumentasi pintasan papan ketik dan koreksi tipe parameter callback untuk `KeyBinding.Add` oleh @ndianabasi
- Perbaiki dokumentasi mengenai pembuatan binding kustom: harus menggunakan `-d String` alih-alih `-o String`
- Perbaiki menu yang tidak menghapus elemen anak pada `menu.Update()`
- Perbaiki referensi Manager API yang usang dalam dokumentasi (berkas 31 diperbarui agar menggunakan pola baru seperti `app.Window.New()`, `app.Event.Emit()`, dan sebagainya) oleh @leaanthony
- Perbaiki crash Linux saat terjadi panic dalam metode Go yang terikat ke JS akibat WebKit menimpa penangan sinyal (#3965) oleh @leaanthony
- Perbaiki SaveFileDialog.SetFilename() yang tidak berpengaruh di Linux (#4841) oleh @samstanier
- Perbaiki koordinat pelepasan yang ditampilkan sebagai undefined dalam contoh seret dan lepas
- Perbaiki kegagalan pembuatan bundel aplikasi macOS saat APP_NAME mengandung spasi (masalah ekspansi kurung kurawal)
- Perbaiki panic indeks di luar batas di Windows saat memanggil metode layanan (kembalikan goccy/go-json)
- Perbaiki seret dan lepas berkas di Windows yang tidak berfungsi pada skala tampilan selain 100%
- Perbaiki seret dan lepas internal HTML5 yang rusak saat fitur pelepasan berkas diaktifkan di Windows
- Perbaiki koordinat pelepasan berkas yang menggunakan ruang piksel yang salah di Windows (piksel fisik dibandingkan piksel CSS)
- Perbaiki seret dan lepas berkas di Linux yang tidak berfungsi secara andal dengan efek hover
- Perbaiki seret dan lepas internal HTML5 yang rusak saat fitur pelepasan berkas diaktifkan di Linux
- Perbarui semua perintah dalam berkas Taskfile.yml untuk seluruh sistem operasi agar dapat menangani spasi dalam variabel seperti `APP_NAME` oleh @ndianabasi
- Perbaiki kesalahan argumen perintah saat menjalankan tugas 'build:universal:lipo:go' di Linux oleh @wux1an
- Memperbaiki galat Docker "undefined symbol: **<em>ubsan</em>handle_xxxxxxx" saat menjalankan 'wails3 build GOOS=darwin GOARCH=arm64' di Linux oleh @wux1an
- Menggabungkan dokumentasi protokol khusus dan menambahkan bagian Universal Links oleh @leaanthony
- Memperbaiki crash menu baki sistem Windows saat ikon diklik berulang kali dengan menambahkan pengaman terhadap pemanggilan TrackPopupMenuEx secara bersamaan (#4151) oleh @leaanthony
- Mencegah aplikasi mengalami crash saat systray.Run() dipanggil sebelum app.Run() oleh @leaanthony
- Memperbaiki crash di macOS saat visibilitas jendela dialihkan melalui Hide()/Show() dengan ApplicationShouldTerminateAfterLastWindowClosed diaktifkan (#4389) oleh @leaanthony
- Memperbaiki kebocoran memori pada menu konteks di macOS dan Windows ketika menu dibuka berulang kali (#4012) oleh @leaanthony
- Memperbaiki sumber daya native menu konteks yang tidak digunakan kembali di macOS sehingga menu baru dibuat setiap kali ditampilkan (#4012) oleh @leaanthony
- Memperbaiki klik ikon Dock macOS yang tidak menampilkan jendela tersembunyi ketika aplikasi dimulai dengan `Hidden: true` (#4583) oleh @leaanthony
- Memperbaiki dialog cetak macOS yang tidak terbuka karena tipe pointer jendela yang salah dalam pemanggilan CGO (#4290) oleh @leaanthony
- Memperbaiki crash menu jendela di Wayland akibat appmenu-gtk-module mengakses jendela yang belum direalisasikan (#4769) oleh @leaanthony
- Memperbaiki crash aplikasi GTK ketika nama aplikasi berisi karakter yang tidak valid (spasi, tanda kurung, dan sebagainya) oleh @leaanthony
- Memperbaiki galat "not enough memory" saat menginisialisasi seret dan lepas di Windows (#4701) oleh @overlordtm
- Memperbaiki penjelajah berkas yang membuka direktori yang salah di Linux akibat pengodean escape URI yang keliru (#4397) oleh @leaanthony
- Memperbaiki kegagalan build AppImage pada distribusi Linux modern (Arch, Fedora 39+, Ubuntu 24.04+) dengan mendeteksi bagian ELF `.relr.dyn` secara otomatis dan menonaktifkan stripping (#4642) oleh @leaanthony
- Memperbaiki `wails doctor` yang keliru melaporkan paket webkit sebagai telah terinstal pada Fedora/sistem berbasis DNF (#4457) oleh @leaanthony
- Memperbaiki perilaku bawaan `config.yml` yang menjalankan `wails3 dev` dengan build produksi oleh @mbaklor
- Memperbaiki stub layanan iOS yang menyebabkan kegagalan build akibat mengimpor paket yang tidak ada oleh @leaanthony
- Memperbaiki logging terstruktur dalam metode debug/info yang menyebabkan galat "no formatting directives" oleh @leaanthony
- Menghapus pernyataan cetak debug sementara yang secara tidak sengaja disertakan dari penggabungan platform seluler oleh @leaanthony
- Memperbaiki crash WebKitGTK di Wayland dengan GPU NVIDIA (Galat 71 Protocol error) dengan menonaktifkan perender DMA-BUF secara otomatis oleh @leaanthony
- Mengatasi nilai alfa yang diabaikan dalam `application.WebviewWindowOptions.BackgroundColour` di Linux ([#4722](https://github.com/wailsapp/wails/pull/4722), @BradHacker)
- Memperbaiki ikon baki sistem Windows yang tidak menggunakan ikon aplikasi secara bawaan ketika ikon khusus tidak disediakan (#4704)
- Melacak kepemilikan `HICON` agar hanya handle buatan pengguna yang dimusnahkan, sehingga mencegah crash saat Explorer mendaur ulang sumber daya (#4653).
- Melepas listener tema sistem Windows dan ikon baki yang dipertahankan saat pemusnahan untuk menghentikan kebocoran goroutine dan konteks perangkat (#4653).
- Memangkas tooltip baki pada 127 unit UTF-16 untuk menghindari kerusakan pasangan surrogate dan glif multibita (#4653).
- Memperbaiki kegagalan task paket Windows (#4667)
- Memperbaiki variabel appicon AppImage Linux dalam taskfile Linux [PR #4644](https://github.com/wailsapp/wails/pull/4644)
- Memperbaiki galat build Windows akibat perubahan signature go-webview2 v1.0.22 (#4513, #4645)
- Memperbaiki variabel appicon AppImage Linux dalam taskfile Linux [PR #4644](https://github.com/wailsapp/wails/pull/4644)
- Memperbaiki rentang protokol desktop.tmpl Linux dengan menghapus `<.Info.Protocol>` hingga `<.Protocol>` oleh @Tolfx dalam #4510
- Memperbaiki galat redefinisi untuk demo liquid glass dalam [#4542](https://github.com/wailsapp/wails/pull/4542) oleh @Etesam913
- Memperbaiki pembaruan menu baki sistem di Linux [#4604](https://github.com/wailsapp/wails/issues/4604) oleh [@JackDoan](https://github.com/JackDoan)
- Memperbaiki jendela putih yang muncul di Windows saat membuat jendela tersembunyi oleh @leaanthony dalam [#4612](https://github.com/wailsapp/wails/pull/4612)
- Memperbaiki jalur impor paket notifications dalam dokumentasi oleh @rxliuli dalam [#4617](https://github.com/wailsapp/wails/pull/4617)
- Memperbaiki seret dan lepas yang tidak berfungsi saat menggunakan paket npm @wailsio/runtime (#4489) oleh @leaanthony dalam #4616
- Windows: Memperbaiki jendela yang berkedip saat mulai dan jendela tersembunyi yang ditampilkan secara keliru dalam [PR](https://github.com/wailsapp/wails/pull/4600) oleh @leaanthony.
- Memperbaiki masalah ukuran jendela saat dimaksimalkan di Wayland (https://github.com/wailsapp/wails/issues/4429) oleh [@samstanier](https://github.com/samstanier)
- Memperbaiki masalah ukuran jendela saat dimaksimalkan di Wayland (https://github.com/wailsapp/wails/issues/4429) oleh [@samstanier](https://github.com/samstanier)
- Memperbaiki galat redefinisi untuk demo liquid glass dalam [#4542](https://github.com/wailsapp/wails/pull/4542) oleh @Etesam913
- Memperbaiki masalah yang dapat menyebabkan AssetServer mengalami crash di macOS dalam [#4576](https://github.com/wailsapp/wails/pull/4576) oleh @jghiloni
- Memperbaiki masalah kompilasi saat melakukan build dengan NextJs. Diperbaiki dalam [#4585](https://github.com/wailsapp/wails/pull/4585) oleh @rev42
- Memperbaiki pipeline untuk rilis nightly dalam [#4597](https://github.com/wailsapp/wails/pull/4597) oleh @riadafridishibly
- Memperbaiki galat redefinisi untuk demo liquid glass dalam [#4542](https://github.com/wailsapp/wails/pull/4542) oleh @Etesam913
- Memperbaiki masalah yang dapat menyebabkan AssetServer mengalami crash di macOS dalam [#4576](https://github.com/wailsapp/wails/pull/4576) oleh @jghiloni
- Memperbaiki masalah kompilasi saat melakukan build dengan NextJs. Diperbaiki dalam [#4585](https://github.com/wailsapp/wails/pull/4585) oleh @rev42
- Memperbaiki pipeline untuk rilis nightly dalam [#4597](https://github.com/wailsapp/wails/pull/4597) oleh @riadafridishibly
- Memperbaiki galat redefinisi untuk demo liquid glass dalam [#4542](https://github.com/wailsapp/wails/pull/4542) oleh @Etesam913
- Memperbaiki SetBackgroundColour di Windows oleh @PPTGamer dalam [PR](https://github.com/wailsapp/wails/pull/4492)
- Memperbarui dokumentasi agar mencerminkan perubahan dari refaktorisasi Manager API oleh @yulesxoxo dalam [PR #4476](https://github.com/wailsapp/wails/pull/4476)
- Memperbaiki variabel appicon berkas .desktop Linux dalam taskfile Linux [PR #4477](https://github.com/wailsapp/wails/pull/4477)
- Memperbarui dokumentasi agar mencerminkan perubahan dari refaktorisasi Manager API oleh @yulesxoxo dalam [PR #4476](https://github.com/wailsapp/wails/pull/4476)
- Memperbaiki bug dereferensi pointer nil di Windows yang dilaporkan dalam [#4456](https://github.com/wailsapp/wails/issues/4456) oleh @leaanthony dalam [#4460](https://github.com/wailsapp/wails/pull/4460)
- Menambahkan dukungan untuk `allowsBackForwardNavigationGestures` di WKWebView macOS guna mengaktifkan gestur navigasi usap dua jari (#1857)
- Memperbaiki masalah yang menyebabkan onClick tidak berfungsi untuk item menu yang awalnya ditetapkan sebagai nonaktif oleh @leaanthony dalam [PR #4469](https://github.com/wailsapp/wails/pull/4469). Terima kasih kepada @IanVS atas investigasi awalnya.
- Memperbaiki server Vite yang tidak dibersihkan ketika build gagal (#4403)
- Memperbaiki panic saat menutup atau membatalkan `SaveFileDialog` di Windows. Diperbaiki dalam [PR](https://github.com/wailsapp/wails/pull/4284) oleh @hkhere
- Memperbaiki seret dan lepas tingkat HTML di Windows oleh [@mbaklor](https://github.com/mbaklor) dalam [#4259](https://github.com/wailsapp/wails/pull/4259)
- Menambahkan dukungan untuk `allowsBackForwardNavigationGestures` di WKWebView macOS guna mengaktifkan gestur navigasi usap dua jari (#1857)
- Memperbaiki masalah onClick yang tidak berfungsi untuk item menu yang awalnya ditetapkan sebagai nonaktif oleh @leaanthony dalam [PR #4469](https://github.com/wailsapp/wails/pull/4469). Terima kasih kepada @IanVS atas investigasi awalnya.
- Memperbaiki server Vite yang tidak dibersihkan ketika build gagal (#4403)
- Memperbaiki penguraian notifikasi di Windows oleh @popaprozac dalam [PR](https://github.com/wailsapp/wails/pull/4450)
- Memperbaiki perintah doctor agar memeriksa dependensi Windows SDK oleh [@kodumulo](https://github.com/kodumulo) dalam [#4390](https://github.com/wailsapp/wails/issues/4390)
- Memperbaiki dereferensi pointer nil di processURLRequest untuk Mac oleh [@etesam913](https://github.com/etesam913) dalam [#4366](https://github.com/wailsapp/wails/pull/4366)
- Memperbaiki bug Linux yang mencegah dialog berfilter oleh [@bh90210](https://github.com/bh90210) dalam [#4287](https://github.com/wailsapp/wails/pull/4287)
- Memperbaiki masalah menu Edit di Windows dan Linux oleh [@leaanthony](https://github.com/leaanthony) dalam [#3f78a3a](https://github.com/wailsapp/wails/commit/3f78a3a8ce7837e8b32242c8edbbed431c68c062)
- Memperbarui versi sistem minimum dalam file .plist macOS dari 10.13.0 menjadi 10.15.0 oleh [@AkshayKalose](https://github.com/AkshayKalose) dalam [#3981](https://github.com/wailsapp/wails/pull/3981)
- Masalah ID jendela yang terlewati oleh [@leaanthony](https://github.com/leaanthony)
- Memperbaiki masalah menu nil saat memanggil RegisterContextMenu oleh [@leaanthony](https://github.com/leaanthony)
- Memperbaiki siklus dependensi dalam keluaran generator binding oleh [@fbbdev](https://github.com/fbbdev) dalam [#4001](https://github.com/wailsapp/wails/pull/4001)
- Memperbaiki kesalahan penggunaan sebelum definisi dalam keluaran generator binding oleh [@fbbdev](https://github.com/fbbdev) dalam [#4001](https://github.com/wailsapp/wails/pull/4001)
- Meneruskan flag build ke generator binding oleh [@fbbdev](https://github.com/fbbdev) dalam [#4023](https://github.com/wailsapp/wails/pull/4023)
- Mengubah path dalam Taskfile Windows agar menggunakan garis miring untuk memastikan file tersebut berfungsi pada platform non-Windows oleh [@leaanthony](https://github.com/leaanthony)
- Event JS Mac + Mac kini telah diperbaiki oleh [@leaanthony](https://github.com/leaanthony)
- Memperbaiki deadlock event untuk macOS oleh [@leaanthony](https://github.com/leaanthony)
- Memperbaiki kesalahan `Parameter incorrect` saat inisialisasi Window di Windows ketika HTML disediakan tanpa JS oleh [@leaanthony](https://github.com/leaanthony)
- Memperbaiki ukuran prefiks respons yang digunakan untuk mendeteksi tipe konten di server aset oleh [@fbbdev](https://github.com/fbbdev) dalam [#4049](https://github.com/wailsapp/wails/pull/4049)
- Memperbaiki penanganan respons non-404 pada path indeks root di server aset oleh [@fbbdev](https://github.com/fbbdev) dalam [#4049](https://github.com/wailsapp/wails/pull/4049)
- Memperbaiki perilaku tak terdefinisi dalam generator binding saat menguji properti tipe generik oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Memperbaiki keluaran generator binding untuk model ketika tipe dasarnya tidak memiliki properti yang sama dengan wrapper bernama oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Memperbaiki keluaran generator binding untuk tipe kunci map dan prapemrosesan oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Memperbaiki keluaran generator binding untuk struct yang mengimplementasikan antarmuka marshaler oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Memperbaiki deteksi siklus tipe yang melibatkan tipe generik dalam generator binding oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Memperbaiki referensi tidak valid ke model yang tidak diekspor dalam keluaran generator binding oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Memindahkan kode yang diinjeksi ke akhir file layanan oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Memperbaiki penanganan kesalahan dari operasi penutupan file dalam generator binding oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Menonaktifkan peringatan untuk layanan yang mendefinisikan metode siklus hidup atau http, tetapi tidak memiliki metode terikat lainnya oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Memperbaiki templat non-React yang gagal menampilkan footer Hello World saat menggunakan skema warna sistem terang oleh [@marcus-crane](https://github.com/marcus-crane) dalam [#4056](https://github.com/wailsapp/wails/pull/4056)
- Memperbaiki item menu tersembunyi di macOS oleh [@leaanthony](https://github.com/leaanthony)
- Memperbaiki penanganan dan pemformatan kesalahan dalam pemroses pesan oleh [@fbbdev](https://github.com/fbbdev) dalam [#4066](https://github.com/wailsapp/wails/pull/4066)
-  Memperbaiki penghentian layanan yang terlewati saat keluar dari aplikasi oleh [@fbbdev](https://github.com/fbbdev) dalam [#4066](https://github.com/wailsapp/wails/pull/4066)
-  Memastikan pembaruan menu berlangsung di thread utama oleh [@leaanthony](https://github.com/leaanthony)
- Mekanisme penyeretan dan pengubahan ukuran kini lebih andal dan lebih sesuai dengan perilaku platform yang diharapkan oleh [@fbbdev](https://github.com/fbbdev) dalam [#4100](https://github.com/wailsapp/wails/pull/4100)
- Memperbaiki [#4097](https://github.com/wailsapp/wails/issues/4097): Webpack/angular membuang kode inisialisasi runtime oleh [@fbbdev](https://github.com/fbbdev) dalam [#4100](https://github.com/wailsapp/wails/pull/4100)
- Memperbaiki item menu yang awalnya tersembunyi oleh [@IanVS](https://github.com/IanVS) dalam [#4116](https://github.com/wailsapp/wails/pull/4116)
- Memperbaiki assetFileServer yang tidak menyajikan file `.html` untuk permintaan tanpa ekstensi ketika `[request]` tidak ada, tetapi `[request].html` ada
- Memperbaiki path pembuatan ikon oleh [@robin-samuel](https://github.com/robin-samuel) dalam [#4125](https://github.com/wailsapp/wails/pull/4125)
- Memperbaiki event `fullscreen`, `unfullscreen`, `unminimise`, dan `unmaximise` yang tidak dipancarkan oleh [@oSethoum](https://github.com/osethoum) dalam [#4130](https://github.com/wailsapp/wails/pull/4130)
- Memperbaiki Kesalahan NSIS akibat prefiks yang salah pada versi default dalam konfigurasi oleh [@robin-samuel](https://github.com/robin-samuel) dalam [#4126](https://github.com/wailsapp/wails/pull/4126)
- Memperbaiki fungsi runtime Dialogs yang mengembalikan path dengan karakter escape di Windows oleh [TheGB0077](https://github.com/TheGB0077) dalam [#4188](https://github.com/wailsapp/wails/pull/4188)
- Memperbaiki path deteksi Webview2 di HKCU oleh [@leaanthony](https://github.com/leaanthony).
- Memperbaiki masalah input di macOS oleh [@leaanthony](https://github.com/leaanthony).
- Memperbaiki nama file tugas pembuatan ikon Windows oleh [@yulesxoxo](https://github.com/yulesxoxo) dalam [#4219](https://github.com/wailsapp/wails/pull/4219).
- Memperbaiki masalah transparansi pada jendela tanpa bingkai oleh [@leaanthony](https://github.com/leaanthony) berdasarkan karya @kron.
- Memperbaiki pemanggilan fokus saat jendela dinonaktifkan atau diminimalkan oleh [@leaanthony](https://github.com/leaanthony) berdasarkan karya @kron.
- Memperbaiki baki sistem yang tidak muncul setelah bilah tugas dimulai ulang oleh [@leaanthony](https://github.com/leaanthony) berdasarkan karya @kron.
- Memperbaiki fallbackResponseWriter yang tidak mengimplementasikan Flush() di [#4245](https://github.com/wailsapp/wails/pull/4245)
- Memperbaiki fallbackResponseWriter yang tidak mengimplementasikan Flush() oleh [@superDingda] di [#4236](https://github.com/wailsapp/wails/issues/4236)
- Memperbaiki crash saat menutup jendela macOS ketika masih ada pemanggilan asinkron fungsi terikat Go yang tertunda oleh [@joshhardy](https://github.com/joshhardy) di [#4354](https://github.com/wailsapp/wails/pull/4354)
- Memperbaiki kondisi balapan saat memulai mode Efisiensi Windows oleh [@leaanthony](https://github.com/leaanthony)
- Memperbaiki pembersihan handle ikon Windows oleh [@leaanthony](https://github.com/leaanthony).
- Memperbaiki `OpenFileManager` di Windows oleh [@PPTGamer](https://github.com/PPTGamer) dalam [#4375](https://github.com/wailsapp/wails/pull/4375).
- Memperbaiki opsi lebar minimum/maksimum untuk Linux oleh @atterpac dalam [#3979](https://github.com/wailsapp/wails/pull/3979)
- Definisi tipe templat TypeScript melalui peningkatan versi npm oleh @atterpac dalam [#3966](https://github.com/wailsapp/wails/pull/3966)
- Memperbaiki referensi CSS templat SvelteKit oleh @atterpac dalam [#3945](https://github.com/wailsapp/wails/pull/3945)
- Memastikan callback utama di run() jendela dipanggil pada thread utama oleh [@leaanthony](https://github.com/leaanthony)
- Memperbaiki contoh pemilih direktori dialog oleh [@leaanthony](https://github.com/leaanthony)
- Membuat halaman galat baru berbahasa Tionghoa ketika index.html tidak ada oleh [@leaanthony](https://github.com/leaanthony)
-  Memastikan callback `windowDidBecomeKey` berjalan pada thread utama oleh [@leaanthony](https://github.com/leaanthony)
-  Mendukung layar penuh untuk jendela tanpa bingkai oleh [@leaanthony](https://github.com/leaanthony)
-  Meningkatkan logika penghancuran jendela oleh [@leaanthony](https://github.com/leaanthony)
-  Memperbaiki logika posisi jendela ketika dikaitkan dengan baki sistem oleh [@leaanthony](https://github.com/leaanthony)
-  Mendukung layar penuh untuk jendela tanpa bingkai oleh [@leaanthony](https://github.com/leaanthony)
- Memperbaiki penanganan peristiwa oleh [@leaanthony](https://github.com/leaanthony)
- Memperbaiki logika penghentian jendela oleh [@leaanthony](https://github.com/leaanthony)
- Taskfile umum kini secara default menghasilkan binding TypeScript untuk templat TypeScript oleh [@leaanthony](https://github.com/leaanthony)
- Memperbaiki penutupan aplikasi saat menerima pesan WM_CLOSE ketika tidak ada jendela yang terbuka/hanya ada baki sistem oleh [@mmalcek](https://github.com/mmalcek) dalam [#3990](https://github.com/wailsapp/wails/pull/3990)
- Memperbaiki build garble oleh @5aaee9 dalam [#3192](https://github.com/wailsapp/wails/pull/3192)
- Memperbaiki build NSIS Windows oleh [@leaanthony](https://github.com/leaanthony)
- Memperbaiki deadlock pada dialog Linux untuk beberapa pilihan yang disebabkan oleh tidak ditutupnya
- Memperbaiki pembersihan lintas platform untuk berkas .syso selama build Windows oleh
- Memperbaiki kompilasi AppImage amd64 oleh @atterpac dalam
- Memperbaiki pembaruan aset build oleh @ansxuman dalam
- Memperbaiki implementasi `OnClick` dan `OnRightClick` pada baki sistem Linux oleh @atterpac
- Memperbaiki `AlwaysOnTop` yang tidak berfungsi di Mac oleh
-  Memperbaiki `application.NewEditMenu` yang menyertakan duplikat
- 🐧 Memperbaiki kompilasi aarch64
- ⊞ Memperbaiki item menu grup radio oleh
- Memperbaiki galat saat membangun .app yang dapat dijalankan di macOS ketika 'name' dan 'outputfilename'
- Memperbaiki bug pada penggunaan customEventProcessor dalam contoh seret dan lepas oleh
- 🐧 Memperbaiki galat kompilasi Linux yang disebabkan oleh penambahan IgnoreMouseEvents oleh
- ⊞ Memperbaiki bug pembuatan berkas ikon syso oleh
- 🐧 Mengintegrasikan perbaikan agar berjalan secara native di Wayland dari
- Tidak mengikat metode layanan internal dalam
- ⊞ Memperbaiki panic saat memulai baki sistem dalam
- Tidak mengikat metode layanan internal dalam
- ⊞ Memperbaiki panic saat memulai baki sistem dalam
- Refaktor besar pada item menu dan penanganan peristiwa. Untuk saat ini, terutama meningkatkan macOS. Oleh
- Memperbaiki pengujian setelah refaktor plugin dan peristiwa dalam
- ⊞ Memperbaiki peringatan `Failed to unregister class Chrome_WidgetWin_0`. Oleh
- Masalah modul
- Memperbaiki pengiriman pesan peristiwa pengubahan ukuran oleh [atterpac](https://github.com/atterpac) dalam
- 🐧 Memperbaiki galat penanganan tema di NixOS oleh
- Memperbaiki instalasi proyek lintas volume untuk Windows oleh
- Memperbaiki CSS templat React agar footer ditampilkan oleh
- Memperbaiki proses zombi saat bekerja dalam mode pengembangan dengan memperbarui refresh ke versi terbaru
- Memperbaiki pengambilan berkas WebKit AppImage oleh [Atterpac](https://github.com/atterpac)
- Memperbaiki verifikasi paket apt oleh Doctor oleh [Atterpac](https://github.com/Atterpac) dalam
- Memperbaiki aplikasi yang membeku saat ditutup (Darwin) oleh @5aaee9 dalam
- Memperbaiki warna latar belakang contoh di Windows oleh
- Memperbaiki menu konteks default oleh [mmghv](https://github.com/mmghv) dalam
- Memperbaiki nilai heksadesimal untuk tombol panah di Darwin oleh
- Membuat seret dan lepas berfungsi di Windows. Ditambahkan oleh
- Memperbaiki bug Linux pada doctor ketika pengguna tidak memiliki driver yang sesuai
- Memperbaiki penskalaan DPI saat memulai (Windows). Diubah oleh [@almas-x](https://github.com/almas-x) dalam
- Memperbaiki baris penggantian dalam `go.mod` agar menggunakan path relatif. Memperbaiki path Windows dengan
- Memperbaiki penanganan klik baki sistem macOS ketika tidak ada jendela yang dikaitkan oleh
- Memperbaiki kegagalan build Windows akibat opsi yang tidak dikenal oleh
- Memperbaiki crash di Windows saat mengklik kiri ikon baki sistem ketika tidak memiliki
- Memperbaiki baseURL yang salah ketika membuka jendela dua kali oleh @5aaee9 dalam PR
- Memperbaiki urutan cabang if dalam metode `WebviewWindow.Restore` oleh
- Menghitung `startURL` dengan benar di beberapa pemanggilan `GetStartURL` ketika
- Memperbaiki tipe JS dari struct `Screen` agar sesuai dengan padanannya di Go oleh
- Memperbaiki metode `WML.Reload` untuk memastikan pembersihan yang tepat atas peristiwa yang terdaftar
- Memperbaiki menu konteks khusus yang langsung tertutup di Linux oleh
- Memperbaiki path keluaran dan ekstensi berkas model yang dihasilkan oleh binding
- Memperbaiki path impor berkas model dalam kode JS yang dihasilkan oleh binding
- Memperbaiki seret dan lepas pada beberapa distro Linux oleh
- Memperbaiki task macOS yang hilang saat menggunakan `wails3 task dev` oleh
- Memperbaiki pendaftaran peristiwa yang menyebabkan penetapan pada map nil oleh
- Memperbaiki unmarshaling parameter metode terikat oleh
- Memperbaiki penanganan beberapa nilai kembalian dari metode terikat oleh
- Memperbaiki deteksi npm oleh doctor jika npm tidak diinstal dengan pengelola paket sistem
- Memperbaiki MicrosoftEdgeWebview2Setup.exe yang hilang. Terima kasih kepada
- Memperbaiki crash acak di Linux akibat penanganan ID jendela oleh @leaanthony. Berdasarkan
- Memperbaiki systemTray.setIcon yang menyebabkan crash di Linux oleh
- Memastikan bingkai jendela diterapkan pada pemanggilan pertama dalam fungsi `setFrameless` pada

### Diubah

- **PERUBAHAN INKOMPATIBEL**: Kunci map dalam binding JS/TS yang dihasilkan kini ditandai sebagai opsional agar secara akurat mencerminkan semantik map Go. Akses nilai map dalam Typescript kini mengembalikan `T | undefined`, bukan `T`, sehingga memerlukan pemeriksaan null atau assertion (#4943) oleh `@fbbdev`
- Mengubah penggunaan `Event` menjadi `Events` sesuai dengan perubahan pada `@wailsio/runtime` serta pemanggilan fungsi yang sesuai dalam dokumentasi di `Features/Events/Event System` oleh @AbdelhadiSeddar
- Memindahkan `EnabledFeatures`, `DisabledFeatures`, dan `AdditionalBrowserArgs` dari opsi per jendela ke `Options.Windows` tingkat aplikasi (#4559) oleh @leaanthony
- Memperbarui README untuk contoh `Drag N Drop` dan menegaskan bahwa `Internal Drag and Drop` didemonstrasikan dalam contoh tersebut oleh @ndianabasi
- Mengubah berbagai log debug dari Info menjadi Debug (oleh @mbaklor)
- **PERUBAHAN INKOMPATIBEL:** Mengganti nama `EnableDragAndDrop` menjadi `EnableFileDrop` dalam opsi jendela
- **PERUBAHAN INKOMPATIBEL:** Mengganti nama `DropZoneDetails` menjadi `DropTargetDetails` dalam konteks peristiwa
- **PERUBAHAN INKOMPATIBEL:** Mengganti nama metode `DropZoneDetails()` menjadi `DropTargetDetails()` pada `WindowEventContext`
- **PERUBAHAN INKOMPATIBEL:** Menghapus peristiwa `WindowDropZoneFilesDropped`; gunakan `WindowFilesDropped` sebagai gantinya
- **PERUBAHAN INKOMPATIBEL:** Mengubah atribut HTML dari `data-wails-dropzone` menjadi `data-file-drop-target`
- **PERUBAHAN INKOMPATIBEL:** Mengubah kelas hover CSS dari `wails-dropzone-hover` menjadi `file-drop-target-active`
- **PERUBAHAN INKOMPATIBEL:** Menghapus opsi `DragEffect`, `OnEnterEffect`, dan `OnOverEffect` dari Windows (sebelumnya merupakan bagian dari IDropTarget yang telah dihapus)
- Beralih ke goccy/go-json untuk seluruh pemrosesan JSON runtime (binding metode, peristiwa, permintaan webview, notifikasi, kvstore), sehingga meningkatkan performa sebesar 21-63% dan mengurangi alokasi memori sebesar 40-60%
- Mengoptimalkan tata letak struct BoundMethod dan menyimpan flag isVariadic dalam cache untuk mengurangi overhead per pemanggilan
- Menggunakan buffer argumen yang dialokasikan di stack untuk metode dengan `<=8` argumen guna menghindari alokasi heap
- Mengoptimalkan pengumpulan hasil dalam pemanggilan metode untuk menghindari alokasi slice bagi satu nilai kembalian
- Menggunakan sync.Map untuk cache jenis MIME guna meningkatkan performa konkurensi
- Menggunakan kumpulan buffer untuk membaca body permintaan transport HTTP
- Mengalokasikan channel CloseNotify secara lazy dalam pendeteksi jenis konten untuk mengurangi alokasi per permintaan
- Menghapus logging CSS debug dari server aset
- Memperluas map ekstensi jenis MIME agar mencakup 50+ format web umum (font, audio, video, dan sebagainya)
- Memperbarui dokumentasi untuk opsi Window `X/Y` oleh @ruhuang2001
- Memperbarui dokumentasi `Frontend Runtime` dengan menambahkan lebih banyak opsi untuk menghasilkan binding frontend oleh @ndianabasi
- Memperbarui halaman dokumentasi Asset Server v3 Wails oleh @ndianabasi
- **PERUBAHAN INKOMPATIBEL**: Menghapus fungsi dialog tingkat paket (`application.InfoDialog()`, `application.QuestionDialog()`, dan sebagainya). Gunakan pengelola `app.Dialog` sebagai gantinya: `app.Dialog.Info()`, `app.Dialog.Question()`, `app.Dialog.Warning()`, `app.Dialog.Error()`, `app.Dialog.OpenFile()`, `app.Dialog.SaveFile()`
- Memperbarui dokumentasi dialog agar sesuai dengan API sebenarnya: gunakan `app.Dialog.*`, `AddButton()` dengan callback (bukan `SetButtons()`), `SetDefaultButton(*Button)` (bukan string), `AddFilter()` (bukan `SetFilters()`), `SetFilename()` (bukan `SetDefaultFilename()`), serta `app.Dialog.OpenFile().CanChooseDirectories(true)` untuk memilih folder
- **PERUBAHAN INKOMPATIBEL**: Build produksi kini menjadi default. Untuk membuat build pengembangan, tetapkan `DEV=true` dalam Taskfile Anda. Menghasilkan proyek baru untuk contoh oleh @leaanthony
- Saat memancarkan peristiwa khusus dengan nol atau satu argumen data, nilai data akan ditetapkan langsung ke field Data tanpa membungkusnya dalam slice oleh [@fbbdev](https://github.com/fbbdev) dalam [#4633](https://github.com/wailsapp/wails/pull/4633)
- Tray Windows kini mematuhi `SystemTray.Show()`/`Hide()` dengan mengalihkan `NIS_HIDDEN`, sehingga aplikasi benar-benar dapat menghilang dan muncul kembali (#4653).
- Pendaftaran tray menggunakan kembali ikon yang telah di-resolve, menetapkan `NOTIFYICON_VERSION_4` satu kali, dan mengaktifkan `NIF_SHOWTIP` agar tooltip pulih setelah Explorer dimulai ulang (#4653).
- macOS: Menggunakan `visibleFrame` sebagai pengganti `frame` untuk memusatkan jendela dengan mengecualikan area bilah menu dan dock
- macOS: Menggunakan `visibleFrame` sebagai pengganti `frame` untuk memusatkan jendela dengan mengecualikan area bilah menu dan dock
- Saat menjalankan `wails3 update build-assets` dengan parameter `-config`, nilai yang ditetapkan melalui parameter `-product*` adalah
- `window.NativeWindowHandle()` -> `window.NativeWindow()` oleh @leaanthony dalam [#4471](https://github.com/wailsapp/wails/pull/4471)
- Memfaktorkan ulang penanganan jendela internal oleh @leaanthony dalam [#4471](https://github.com/wailsapp/wails/pull/4471)
- Menghapus `application.WindowIDKey` dan `application.WindowNameKey` (digantikan oleh `application.WindowKey`) oleh [@leaanthony](https://github.com/leaanthony)
- ContextMenuData kini mengembalikan string, bukan any, oleh [@leaanthony](https://github.com/leaanthony)
- Dalam binding JS/TS, field kelas dengan jenis array berukuran tetap kini diinisialisasi dengan panjang yang diharapkan, bukan dibiarkan kosong, oleh [@fbbdev](https://github.com/fbbdev) dalam [#4001](https://github.com/wailsapp/wails/pull/4001)
- ContextMenuData kini mengembalikan string, bukan any, oleh [@leaanthony](https://github.com/leaanthony)
- `application.NewService` tidak lagi menerima options sebagai parameter opsional (gunakan `application.NewServiceWithOptions` sebagai gantinya) oleh [@leaanthony](https://github.com/leaanthony) dalam [#4024](https://github.com/wailsapp/wails/pull/4024)
- Menghapus dependensi `nanoid` oleh [@leaanthony](https://github.com/leaanthony)
- Memperbarui contoh Window untuk gaya jendela mica/acrylic/tabbed oleh [@leaanthony](https://github.com/leaanthony)
- Dalam binding JS/TS, file model `internal.js/ts` telah dihapus; semua model kini dapat ditemukan di `models.js/ts` oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Dalam binding JS/TS, named type tidak pernah dirender sebagai alias untuk named type lain; perilaku lama kini dibatasi hanya untuk alias oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Dalam binding JS/TS pada mode kelas, field struct yang jenisnya merupakan parameter jenis ditandai sebagai opsional dan tidak pernah diinisialisasi secara otomatis oleh [@fbbdev](https://github.com/fbbdev) dalam [#4045](https://github.com/wailsapp/wails/pull/4045)
- Menghapus ESLint dari templat oleh [@IanVS](https://github.com/IanVS) dalam [#4059](https://github.com/wailsapp/wails/pull/4059)
- Memperbarui tahun hak cipta menjadi 2025 oleh [@IanVS](https://github.com/IanVS) dalam [#4037](https://github.com/wailsapp/wails/pull/4037)
- Tambahkan dokumentasi untuk event.Sender oleh [@IanVS](https://github.com/IanVS) di [#4075](https://github.com/wailsapp/wails/pull/4075)
- Dukungan Go 1.24 oleh [@leaanthony](https://github.com/leaanthony)
- Hook `ServiceStartup` kini dipanggil saat `App.Run` dipanggil, bukan di `application.New`, oleh [@fbbdev](https://github.com/fbbdev) di [#4066](https://github.com/wailsapp/wails/pull/4066)
- Kesalahan `ServiceStartup` kini dikembalikan dari `App.Run` alih-alih menghentikan proses, oleh [@fbbdev](https://github.com/fbbdev) di [#4066](https://github.com/wailsapp/wails/pull/4066)
- Panggilan binding dan dialog dari JS kini ditolak dengan objek kesalahan alih-alih string, oleh [@fbbdev](https://github.com/fbbdev) di [#4066](https://github.com/wailsapp/wails/pull/4066)
- Pemosisian menu systray di Windows ditingkatkan oleh [@leaanthony](https://github.com/leaanthony)
- Runtime JS telah diporting ke TypeScript oleh [@fbbdev](https://github.com/fbbdev) di [#4100](https://github.com/wailsapp/wails/pull/4100)
- Runtime diinisialisasi segera setelah diimpor sehingga tidak perlu menunggu jendela dimuat, oleh [@fbbdev](https://github.com/fbbdev) di [#4100](https://github.com/wailsapp/wails/pull/4100)
- Runtime tidak lagi mengekspor metode init. Impor efek samping dapat digunakan untuk menginisialisasinya, oleh [@fbbdev](https://github.com/fbbdev) di [#4100](https://github.com/wailsapp/wails/pull/4100)
- Metode terikat kini mengembalikan `CancellablePromise` yang ditolak dengan `CancelError` jika dibatalkan. Hasil panggilan yang sebenarnya dibuang, oleh [@fbbdev](https://github.com/fbbdev) di [#4100](https://github.com/wailsapp/wails/pull/4100)
- Tipe layanan bawaan kini secara konsisten disebut `Service`, oleh [@fbbdev](https://github.com/fbbdev) di [#4067](https://github.com/wailsapp/wails/pull/4067)
- Fungsi pembuatan layanan bawaan dengan opsi kini secara konsisten disebut `NewWithConfig`, oleh [@fbbdev](https://github.com/fbbdev) di [#4067](https://github.com/wailsapp/wails/pull/4067)
- Metode `Select` pada layanan `sqlite` kini bernama `Query` agar konsisten dengan API Go, oleh [@fbbdev](https://github.com/fbbdev) di [#4067](https://github.com/wailsapp/wails/pull/4067)
- Templat: pindahkan runtime ke "dependencies" dan tata file package.json, oleh [@IanVS](https://github.com/IanVS) di [#4133](https://github.com/wailsapp/wails/pull/4133)
- Buat dan tanda tangani secara ad hoc bundel aplikasi dalam mode pengembangan untuk mengaktifkan API macOS tertentu, oleh [@popaprozac](https://github.com/popaprozac) di [#4171](https://github.com/wailsapp/wails/pull/4171)
- Aset build dipindahkan ke direktori khusus platform oleh [@leaanthony](https://github.com/leaanthony)
- Taskfile dipindahkan dan diganti namanya ke direktori khusus platform oleh [@leaanthony](https://github.com/leaanthony)
- Pengalaman saat `index.html` tidak tersedia ditingkatkan secara signifikan oleh [@leaanthony](https://github.com/leaanthony)
- [Windows] Performa saat meminimalkan dan memulihkan ditingkatkan oleh [@leaanthony](https://github.com/leaanthony). Berdasarkan [PR](https://github.com/wailsapp/wails/pull/3955) asli oleh [562589540](https://github.com/562589540)
- Opsi `ShouldClose` dihapus (sebagai gantinya, daftarkan hook untuk events.Common.WindowClosing) oleh [@leaanthony](https://github.com/leaanthony)
- [Windows] Kedipan saat membuka jendela dikurangi oleh [@leaanthony](https://github.com/leaanthony)
- `Window.Destroy` dihapus karena dimaksudkan sebagai fungsi internal, oleh [@leaanthony](https://github.com/leaanthony)
- Nama event `WindowClose` diubah menjadi `WindowClosing` oleh [@leaanthony](https://github.com/leaanthony)
- Build frontend kini menggunakan environment vite "development" atau "production" sesuai tipe build, oleh [@leaanthony](https://github.com/leaanthony)
- Perbarui ke go-webview2 v1.19 oleh [@leaanthony](https://github.com/leaanthony)
- Pastikan fork taskfile digunakan oleh @leaanthony
- Perbarui fork Taskfile untuk memperbaiki masalah versi saat menginstal menggunakan
- Gunakan fork Taskfile untuk memperbaiki masalah versi saat menginstal menggunakan
- `service.OnStartup` kini mematikan aplikasi ketika terjadi kesalahan dan menjalankan
- Pemrosesan pesan klik systray direfaktor agar lebih selaras dengan interaksi pengguna oleh
- Penyematan aset kini menyertakan `all:frontend/dist` untuk mendukung framework yang menghasilkan
- Refaktor Taskfile oleh [leaanthony](https://github.com/leaanthony) di
- Tingkatkan ke `go-webview2` v1.0.16 oleh
- Perbaiki tipe `Screen` agar menyertakan `ID`, bukan `Id`, oleh
- Perbarui versi Wails `go.mod.tmpl` untuk mendukung `application.ServiceOptions` oleh
- Penentuan nama layanan diperbaiki oleh [windom](https://github.com/windom/) di
- mkdocs serve kini menggunakan Docker, oleh [leaanthony](https://github.com/leaanthony)
- Konfigurasi pengembangan digabungkan ke dalam `config.yml` oleh
- Dialog systray kini secara default menggunakan ikon aplikasi jika tersedia (Windows), oleh
- Pelaporan GPU + Memori untuk macOS ditingkatkan oleh
- Hapus `WebviewGpuIsDisabled` dan `EnableFraudulentWebsiteWarnings`
- Perubahan API Events: `On`/`Emit` -> event pengguna, `OnApplicationEvent` ->
- Perbaikan API Events di Linux oleh [TheGB0077](https://github.com/TheGB0077) di
- [CI] peningkatan pada action serta memungkinkan action dijalankan juga di fork dan
- Nama `AbsolutePosition()` diubah menjadi `Position()` oleh
- Perbarui dependensi WebKit Linux ke webkit2gtk-4.1 dari webkitgtk2-4.0 untuk
- Skrip runtime JS yang dibundel kini merupakan modul ESM: tag script yang mengimpornya
- Paket `@wailsio/runtime` tidak memublikasikan API-nya pada `window.wails`
- Modul API Window `@wailsio/runtime/src/window` kini mengekspos objek yang memuatnya
- API jendela JS telah diperbarui agar sesuai dengan `WebviewWindow` Go saat ini
- Generator binding kini secara default menggunakan panggilan berdasarkan ID. Opsi CLI `-id`
- Tata letak kode binding baru: sebelumnya file output ditata dalam folder
- Nama field struct `application.Options.Bind` telah diubah menjadi
- Sintaks baru untuk layanan binding: instance layanan kini harus dibungkus dalam
- Nonaktifkan spinner di lingkungan nonterminal atau CI oleh

### Dihapus

- **BREAKING**: Hapus `EnabledFeatures`, `DisabledFeatures`, dan `AdditionalLaunchArgs` dari opsi `WindowsWindow` per jendela. Gunakan `Options.Windows.EnabledFeatures`, `Options.Windows.DisabledFeatures`, dan `Options.Windows.AdditionalBrowserArgs` tingkat aplikasi sebagai gantinya. Flag ini berlaku secara global pada environment WebView2 bersama (#4559) oleh @leaanthony
- Hapus implementasi native `IDropTarget` di Windows dan gunakan pendekatan berbasis JavaScript sebagai gantinya (sesuai dengan perilaku v2)
- Hapus dependensi github.com/wailsapp/mimetype dan gunakan peta ekstensi yang diperluas + http.DetectContentType dari pustaka standar sebagai gantinya, sehingga ukuran biner berkurang sekitar 1.2 MB
- Hapus dependensi gopkg.in/ini.v1 dengan mengimplementasikan parser minimal untuk berkas .desktop bagi penjelajah berkas Linux, sehingga menghemat ~45KB
- Hapus samber/lo dari kode runtime dengan menggunakan paket slices dari pustaka standar Go 1.21+ dan helper internal minimal, sehingga menghemat ~310KB
- Hapus pernyataan printf debug dari handler skema URL Darwin (#4834)
- **PERUBAHAN YANG MERUSAK KOMPATIBILITAS**: Hapus event `linux:WindowLoadChanged`—gunakan `linux:WindowLoadFinished` sebagai gantinya untuk mendeteksi saat WebView selesai dimuat (#3896) oleh @leaanthony

### Perubahan yang Merusak Kompatibilitas

- **Refaktorisasi API Manager**: Menata ulang API aplikasi dari struktur datar menjadi manager yang terorganisasi agar pengorganisasian kode dan kemudahan penemuannya menjadi lebih baik, oleh [@leaanthony](https://github.com/leaanthony) dalam [#4359](https://github.com/wailsapp/wails/pull/4359)
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

tidak lagi diabaikan dan menimpa nilai konfigurasi.

## v3.0.0-alpha.24 - 2025-08-13

## Ditambahkan

## v3.0.0-alpha.23 - 2025-08-11

## Diperbaiki

## v3.0.0-alpha.22 - 2025-08-10

## Ditambahkan

## Diubah

+ Memperbaiki dependensi paket Linux yang terlalu luas serta dependensi RPM yang sudah usang.

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

### Perubahan yang Merusak Kompatibilitas

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

`service.OnShutdown`untuk semua layanan sebelumnya yang dimulai oleh @atterpac dalam   [#3920](https://github.com/wailsapp/wails/pull/3920)   @atterpac dalam [#3907](https://github.com/wailsapp/wails/pull/3907)   subfolder oleh @atterpac dalam   [#3887](https://github.com/wailsapp/wails/pull/3887)   [#3748](https://github.com/wailsapp/wails/pull/3748)   [leaanthony](https://github.com/leaanthony)   [etesam913](https://github.com/etesam913) dalam   [#3778](https://github.com/wailsapp/wails/pull/3778)   [northes](https://github.com/northes) dalam   [#3836](https://github.com/wailsapp/wails/pull/3836)   [#3827](https://github.com/wailsapp/wails/pull/3827)   [leaanthony](https://github.com/leaanthony)   [@leaanthony](https://github.com/leaanthony)   [@leaanthony](https://github.com/leaanthony)   (digantikan oleh opsi `EnabledFeatures` dan `DisabledFeatures`) oleh   [leaanthony](https://github.com/leaanthony)

### Diperbaiki

variabel channel oleh @michael-freling dalam   [#3925](https://github.com/wailsapp/wails/pull/3925)   [ansxuman](https://github.com/ansxuman) dalam   [#3924](https://github.com/wailsapp/wails/pull/3924)   [#3898](https://github.com/wailsapp/wails/pull/3898)   [#3901](https://github.com/wailsapp/wails/pull/3901)   dalam [#3886](https://github.com/wailsapp/wails/pull/3886)   [leaanthony](https://github.com/leaanthony) dalam   [#3841](https://github.com/wailsapp/wails/pull/3841)   peran `PasteAndMatchStyle` dalam menu edit di Darwin oleh   [johnmccabe](https://github.com/johnmccabe) dalam   [#3839](https://github.com/wailsapp/wails/pull/3839)   [#3840](https://github.com/wailsapp/wails/issues/3840) dalam   [#3854](https://github.com/wailsapp/wails/pull/3854) oleh   [kodflow](https://github.com/kodflow)   [@leaanthony](https://github.com/leaanthony)   berbeda. oleh @nickisworking dalam   [#3789](https://github.com/wailsapp/wails/pull/3789)

## v3.0.0-alpha.7 - 2024-09-18

### Ditambahkan

[mmghv](https://github.com/mmghv) dalam   [#3665](https://github.com/wailsapp/wails/pull/3665)   [#3682](https://github.com/wailsapp/wails/pull/3682)   [atterpac](https://github.com/atterpac) dan   [leaanthony](https://github.com/leaanthony) dalam   [#3570](https://github.com/wailsapp/wails/pull/3570)

### Diubah

Application Events `OnWindowEvent` -> Window Events, oleh   [leaanthony](https://github.com/leaanthony)   [#3734](https://github.com/wailsapp/wails/pull/3734)   branch dengan prefiks `v3/` atau `v3-` oleh   [stendler](https://github.com/stendler) dalam   [#3747](https://github.com/wailsapp/wails/pull/3747)

### Diperbaiki

[etesam913](https://github.com/etesam913) dalam   [#3742](https://github.com/wailsapp/wails/pull/3742)   [atterpac](https://github.com/atterpac) dalam   [#3721](https://github.com/wailsapp/wails/pull/3721)   [atterpac](https://github.com/atterpac) dalam   [#3675](https://github.com/wailsapp/wails/pull/3675)   [#1811](https://github.com/wailsapp/wails/pull/1811) dalam   [#3614](https://github.com/wailsapp/wails/pull/3614) oleh   [@stendler](https://github.com/stendler)   [#3720](https://github.com/wailsapp/wails/pull/3720) oleh   [leaanthony](https://github.com/leaanthony)   [#3693](https://github.com/wailsapp/wails/issues/3693) oleh   [@DeltaLaboratory](https://github.com/DeltaLaboratory)   [#3720](https://github.com/wailsapp/wails/pull/3720) oleh   [leaanthony](https://github.com/leaanthony)   [#3693](https://github.com/wailsapp/wails/issues/3693) oleh   [@DeltaLaboratory](https://github.com/DeltaLaboratory)   [leaanthony](https://github.com/leaanthony)   [#3746](https://github.com/wailsapp/wails/pull/3746) oleh   [@stendler](https://github.com/stendler)   [leaanthony](https://github.com/leaanthony)

## v3.0.0-alpha.6 - 2024-07-30

### Diperbaiki

## v3.0.0-alpha.5 - 2024-07-30

### Ditambahkan

[#3580](https://github.com/wailsapp/wails/pull/3580)   [#3580](https://github.com/wailsapp/wails/pull/3580)   klik ikon oleh @5aaee9 dalam [#2991](https://github.com/wailsapp/wails/pull/2991)   [#2618](https://github.com/wailsapp/wails/pull/2618)   [@fbbdev](https://github.com/fbbdev) dalam   [#3282](https://github.com/wailsapp/wails/pull/3282)   @[Atterpac](https://github.com/Atterpac)   dalam [#3022](https://github.com/wailsapp/wails/pull/3022])   [@marcus-crane](https://github.com/marcus-crane) dalam   [#3146](https://github.com/wailsapp/wails/pull/3146)   [PR](https://github.com/wailsapp/wails/pull/3147)   [PR](https://github.com/wailsapp/wails/pull/3189)   [@fbbdev](https://github.com/fbbdev) dalam   [#3281](https://github.com/wailsapp/wails/pull/3281)   [aba82cc](https://github.com/wailsapp/wails/commit/aba82cc52787c97fb99afa58b8b63a0004b7ff6c)   berdasarkan [PR](https://github.com/wailsapp/wails/pull/2044) oleh @Mai-Lapyst   [@fbbdev](https://github.com/fbbdev) dalam   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@fbbdev](https://github.com/fbbdev) dalam   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@fbbdev](https://github.com/fbbdev) dalam   [#3295](https://github.com/wailsapp/wails/pull/3295)   paket npm oleh [@fbbdev](https://github.com/fbbdev) dalam   [#3334](https://github.com/wailsapp/wails/pull/3334)   dalam [#3354](https://github.com/wailsapp/wails/pull/3354)   `WAILS_VITE_PORT` oleh [@abichinger](https://github.com/abichinger) dalam   [#3429](https://github.com/wailsapp/wails/pull/3429)   [@abichinger](https://github.com/abichinger) dalam   [#3431](https://github.com/wailsapp/wails/pull/3431)   [@bruxaodev](https://github.com/bruxaodev) dalam   [#3667](https://github.com/wailsapp/wails/pull/3667)   [@OlegGulevskyy](https://github.com/OlegGulevskyy) dalam   [#3674](https://github.com/wailsapp/wails/pull/3674)

### Diperbaiki

[#3606](https://github.com/wailsapp/wails/pull/3606)   [tmclane](https://github.com/tmclane) dalam   [#3515](https://github.com/wailsapp/wails/pull/3515)   [atterpac](https://github.com/atterac) dalam   [#3512](https://github.com/wailsapp/wails/pull/3512)   [atterpac](https://github.com/atterpac) dalam   [#3477](https://github.com/wailsapp/wails/pull/3477)   oleh [Atterpac](https://github.com/atterpac) dalam   [#3320](https://github.com/wailsapp/wails/pull/3320).   dalam [#3306](https://github.com/wailsapp/wails/pull/3306).   [#2972](https://github.com/wailsapp/wails/pull/2972).   [#2982](https://github.com/wailsapp/wails/pull/2982)   [mmghv](https://github.com/mmghv) dalam   [#2750](https://github.com/wailsapp/wails/pull/2750).   [#2753](https://github.com/wailsapp/wails/pull/2753).   [jaybeecave](https://github.com/jaybeecave) dalam   [#3052](https://github.com/wailsapp/wails/pull/3052).   [@pylotlight](https://github.com/pylotlight) dalam   [PR](https://github.com/wailsapp/wails/pull/3039)   telah diinstal. Ditambahkan oleh [@pylotlight](https://github.com/pylotlight) dalam   [PR](https://github.com/wailsapp/wails/pull/3032)   [PR](https://github.com/wailsapp/wails/pull/3145)   spasi - @leaanthony.   [thomas-senechal](https://github.com/thomas-senechal) dalam PR   [#3207](https://github.com/wailsapp/wails/pull/3207)   [thomas-senechal](https://github.com/thomas-senechal) dalam PR   [#3208](https://github.com/wailsapp/wails/pull/3208)   jendela terlampir [tw1nk](https://github.com/tw1nk) dalam PR   [#3271](https://github.com/wailsapp/wails/pull/3271)   [#3273](https://github.com/wailsapp/wails/pull/3273)   [@fbbdev](https://github.com/fbbdev) dalam   [#3279](https://github.com/wailsapp/wails/pull/3279)   `FRONTEND_DEVSERVER_URL` tersedia.   [#3299](https://github.com/wailsapp/wails/pull/3299)   [@fbbdev](https://github.com/fbbdev) dalam   [#3295](https://github.com/wailsapp/wails/pull/3295)   listener oleh [@fbbdev](https://github.com/fbbdev) dalam   [#3295](https://github.com/wailsapp/wails/pull/3295)   [@abichinger](https://github.com/abichinger) dalam   [#3330](https://github.com/wailsapp/wails/pull/3330)   generator oleh [@fbbdev](https://github.com/fbbdev) dalam   [#3334](https://github.com/wailsapp/wails/pull/3334)   generator oleh [@fbbdev](https://github.com/fbbdev) dalam   [#3334](https://github.com/wailsapp/wails/pull/3334)   [@abichinger](https://github.com/abichinger) dalam   [#3346](https://github.com/wailsapp/wails/pull/3346)   [@hfoxy](https://github.com/hfoxy) dalam   [#3417](https://github.com/wailsapp/wails/pull/3417)   [@hfoxy](https://github.com/hfoxy) dalam   [#3426](https://github.com/wailsapp/wails/pull/3426)   [@fbbdev](https://github.com/fbbdev) dalam   [#3431](https://github.com/wailsapp/wails/pull/3431)   [@fbbdev](https://github.com/fbbdev) dalam   [#3431](https://github.com/wailsapp/wails/pull/3431)   oleh [@pekim](https://github.com/pekim) dalam   [#3458](https://github.com/wailsapp/wails/pull/3458)   [@robin-samuel](https://github.com/robin-samuel).   PR [#3466](https://github.com/wailsapp/wails/pull/3622) oleh   [@5aaee9](https://github.com/5aaee9).   [@windom](https://github.com/windom/) dalam   [#3636](https://github.com/wailsapp/wails/pull/3636).   Windows oleh [@bruxaodev](https://github.com/bruxaodev/) dalam   [#3691](https://github.com/wailsapp/wails/pull/3691).

### Diubah

[mmghv](https://github.com/mmghv) dalam   [#3611](https://github.com/wailsapp/wails/pull/3611)   mendukung Ubuntu 24.04 LTS oleh [atterpac](https://github.com/atterpac) dalam   [#3461](https://github.com/wailsapp/wails/pull/3461)   harus memiliki atribut `type="module"`. Oleh   [@fbbdev](https://github.com/fbbdev) dalam   [#3295](https://github.com/wailsapp/wails/pull/3295)   objek, dan tidak memulai sistem WML. Perubahan ini dilakukan untuk meningkatkan   enkapsulasi. Jika diinginkan, sistem WML dapat dimulai secara manual dengan memanggil   metode baru `WML.Enable`. Skrip runtime JS yang disertakan masih menjalankan kedua   operasi secara otomatis. Oleh [@fbbdev](https://github.com/fbbdev) dalam   [#3295](https://github.com/wailsapp/wails/pull/3295)   objek window sebagai ekspor default. Metode individual tidak lagi dapat diimpor   melalui sintaks impor bernama atau namespace ESM.   API. Beberapa metode telah berganti nama atau prototipe, khususnya: `Screen`   menjadi `GetScreen`; `GetZoomLevel`/`SetZoomLevel` menjadi `GetZoom`/`SetZoom`;   `GetZoom`, `Width`, dan `Height` kini mengembalikan nilai secara langsung, bukan membungkusnya   di dalam objek. Oleh [@fbbdev](https://github.com/fbbdev) dalam   [#3295](https://github.com/wailsapp/wails/pull/3295)   telah dihapus. Gunakan opsi CLI `-names` untuk kembali ke pemanggilan berdasarkan nama.   Oleh [@fbbdev](https://github.com/fbbdev) dalam   [#3468](https://github.com/wailsapp/wails/pull/3468)   sebelumnya dinamai berdasarkan paket yang memuatnya; kini jalur impor Go lengkap digunakan,   termasuk jalur modul. Oleh [@fbbdev](https://github.com/fbbdev) dalam   [#3468](https://github.com/wailsapp/wails/pull/3468)   `application.Options.Services`. Oleh [@fbbdev](https://github.com/fbbdev) dalam   [#3468](https://github.com/wailsapp/wails/pull/3468)   pemanggilan ke `application.NewService`. Oleh [@fbbdev](https://github.com/fbbdev) dalam   [#3468](https://github.com/wailsapp/wails/pull/3468)   [@DeltaLaboratory](https://github.com/DeltaLaboratory) dalam   [#3574](https://github.com/wailsapp/wails/pull/3574)
