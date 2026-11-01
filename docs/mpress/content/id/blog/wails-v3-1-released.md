---
title: "Wails v3.1.0: Lebih banyak Go, lebih sedikit hambatan build"
description: "Build desktop eksperimental tanpa CGO dan Wake: cara yang lebih sederhana untuk membangun, mengemas, dan menandatangani aplikasi Wails."
authors: ["leaanthony"]
tags: ["wails", "v3", "release"]
date: "2026-11-01"
slug: "blog/wails-v3-1-released"
sourcePath: "blog/wails-v3-1-released.md"
---

Minggu ini sangat sibuk di Wails: build tanpa CGO, penyiapan proyek yang lebih mudah, serta proses pengemasan dan penandatanganan yang lebih lancar.

Wails v3.1.0 memulai rilis fitur mingguan kami agar Anda dapat memperoleh kemampuan baru secara rutin. Perbaikan bug akan terus dirilis melalui pembaruan patch nightly segera setelah siap.

## Perbaikan dan pengoptimalan minggu ini

Minggu ini menghadirkan **100 perbaikan, 21 peningkatan performa, dan empat perubahan penguatan keamanan** di seluruh framework, CLI, alat build, dan CI. Sorotannya mencakup penempatan jendela yang lebih andal pada konfigurasi DPI campuran, perbaikan baki sistem dan mode gelap di Linux, penyajian aset yang lebih efisien, serta penggunaan memori yang lebih rendah selama pengemasan MSIX. [Lihat beberapa perubahannya di bawah ini](#beberapa-sorotan-minggu-ini).

## Build eksperimental tanpa CGO

Purego memungkinkan Go memanggil pustaka sistem native tanpa CGo, jembatan Go yang lazim digunakan untuk berinteraksi dengan kode C. Kami mulai mengeksplorasinya saat pertama kali hadir, tetapi memelihara dua backend merupakan tantangan. Seiring makin matangnya purego dan AI generatif membuat pemeliharaan tersebut lebih mudah dikelola, kami siap membagikan backend eksperimental untuk macOS dan Linux.

Kini Anda dapat membangun dan melakukan kompilasi silang aplikasi desktop tanpa compiler C:

```sh
WAILS_PUREGO=1 wails3 build
```

Setelah membangun aset frontend, Anda juga dapat membangun secara langsung dengan Go:

```sh
CGO_ENABLED=0 go build -tags purego,production .
```

Backend macOS menggunakan AppKit dan WebKit; Linux mendukung GTK4/WebKitGTK 6.0 dan GTK3. Pengguna tetap memerlukan pustaka runtime native, sedangkan pengemasan dan penandatanganan memiliki persyaratan tersendiri. CGo tetap menjadi pilihan default selama kami mengumpulkan masukan.

[Detail implementasi dan platform](https://github.com/wailsapp/wails/pull/6521)

## Mengapa Wake?

Wails v2 terlalu banyak menyembunyikan proses build. Dalam v3, Taskfile membuka seluruh proses tersebut, tetapi kebanyakan orang hanya menginginkan beberapa penyesuaian tanpa harus memelihara sistem build.

Wake menempatkan pilihan umum dalam `wails.yaml`: aplikasi, alat CLI pendamping, server, target, dan pengemasan. Task tetap menjadi sistem build inti kami dan tersedia saat Anda membutuhkan fleksibilitasnya.

Wake kini mendukung pemuatan ulang langsung, build inkremental, output yang di-cache, dan pengemasan paralel. Sudah memiliki berkas biner? `--skip-build` dapat mengemasnya tanpa memerlukan instalasi Go atau alat frontend.

Aktifkan dengan `WAILS_EXP_USE_WAKE=1`.

### Memulai dengan panduan

Saat Wake diaktifkan, `wails3 init -ui` memandu Anda memilih backend, target, dan paket, memeriksa alat yang tersedia, lalu menyimpan konfigurasi proyek.

@image{src="/assets/screenshots/wake/application.png" alt="Wizard proyek Wake yang menampilkan pilihan aplikasi desktop dan backend, termasuk Pure Go eksperimental." expand}

Pilih format paket dan opsi penandatanganan untuk setiap target:

@image{src="/assets/screenshots/wake/packages.png" alt="Pemilihan paket Wake untuk Linux yang menampilkan format AppImage, DEB, RPM, Void, dan Arch Linux serta opsi untuk menandatangani output rilis." expand}

Pemeriksaan alat di bawah ini menampilkan build Linux standar:

@image{src="/assets/screenshots/wake/tools.png" alt="Wake memeriksa compiler C, Go, npm, dan pkg-config untuk build Linux standar." expand}

### Apa yang dapat saya bangun dari mesin saya?

Kompilasi, pengemasan, dan penandatanganan memiliki persyaratan yang berbeda. Berikut yang dapat Anda lakukan dari setiap host:

@table{header search sort column-separators align=["left", "centre", "centre", "centre", "left"]}
| Operasi | Host Windows | Host macOS | Host Linux | Yang Anda perlukan |
| Kompilasi aplikasi desktop dengan purego | ✓ | ✓ | ✓ | Go dan backend purego eksperimental untuk target macOS/Linux |
| Membuat paket MSIX | ✓ | ✓ | ✓ | Aplikasi yang telah dibangun; tidak memerlukan Windows SDK |
| Membuat installer NSIS | ✓ | ✓ | ✓ | Aplikasi yang telah dibangun dan `makensis` |
| Membuat aplikasi macOS dan DMG | ✓ | ✓ | ✓ | Aplikasi yang telah dibangun; Go menangani pembuatan paket |
| Membuat paket DEB, RPM, Arch, dan Void | ✓ | ✓ | ✓ | Aplikasi yang telah dibangun |
| Membuat AppImage | ✓ | ✓ | ✓ | Runtime AppImage target; pohon runtime GUI Linux untuk pengemasan silang |
| Menandatangani paket Windows dengan identitas yang diekspor | ✓ | ✓ | ✓ | Sertifikat RSA PFX dan kunci privat yang diekspor |
| Menandatangani aplikasi Apple dan DMG dengan identitas yang diekspor | ✓ | ✓ | ✓ | Identitas RSA PKCS#12 yang diekspor; lihat catatan validasi di bawah |
| Menandatangani paket DEB, RPM, dan Arch | ✓ | ✓ | ✓ | Kunci OpenPGP yang diekspor |
| Mengirimkan ke notarisasi Apple dan melampirkan tiket | ✓ | ✓ | ✓ | Kunci API tim Apple dan aplikasi yang ditandatangani dengan semestinya |
| Menggunakan `codesign` Apple atau profil keychain | – | ✓ | – | macOS dan identitas penandatanganan yang sesuai |
@end

Build CGo masih memerlukan toolchain C. Hook khusus yang memodifikasi DMG memerlukan macOS dan `hdiutil`.

### Penandatanganan dan notarisasi

Wake mendukung penandatanganan paket Windows, Apple, dan Linux, termasuk tanda tangan ad hoc Apple. Penandatanganan lintas platform menggunakan identitas yang diekspor; identitas dari penyimpanan sertifikat dan keychain tetap memerlukan alat native. Kredensial tingkat pengguna menjaga agar kunci privat tidak masuk ke konfigurasi proyek. Penandatanganan AppImage dan Void belum disertakan.

Notarisasi Apple dilakukan melalui Go dengan kunci API tim. Kami telah menguji pengiriman dan pelampiran tiket dari Linux, dan Gatekeeper menerima hasilnya. Aplikasi tersebut ditandatangani menggunakan [`codesign`](https://developer.apple.com/library/archive/documentation/Security/Conceptual/CodeSigningGuide/Procedures/Procedures.html) Apple; penerimaan Apple terhadap tanda tangan Developer ID yang dihasilkan Go oleh Wake masih belum diverifikasi.

### Beralih ke Wake

Migrasi mendukung kesembilan templat bawaan. Proses ini satu arah: mengaktifkan `wails.yaml` akan menggantikan perutean Taskfile. Perintah khusus dan override yang tidak didukung akan ditolak untuk ditinjau.

[Fitur, persyaratan, dan detail migrasi Wake](https://github.com/wailsapp/wails/pull/6531)

## Beberapa sorotan minggu ini

- **Jendela ditempatkan sebagaimana mestinya pada layar dengan DPI campuran.** Batas dan posisi layar kini ditangani dengan benar pada konfigurasi macOS dan Linux yang menggunakan skala tampilan berbeda. [#6346](https://github.com/wailsapp/wails/pull/6346)
- **Baki sistem dan deteksi tema di Linux bekerja lebih andal.** Baki sistem yang dibuat setelah aplikasi dimulai tidak lagi diabaikan, tooltip kini terisi, dan deteksi mode gelap membaca pengaturan tampilan desktop sejak aplikasi dimulai.
- **Callback event Go mempertahankan urutannya untuk setiap listener.** Setiap listener berjalan secara serial sesuai urutan emit; listener yang berbeda dapat berjalan secara bersamaan, sedangkan emit yang dilakukan secara bersamaan tidak memiliki urutan yang ditentukan. Hook event disalin sebagai snapshot saat berada di bawah lock untuk menghindari data race. [#6275](https://github.com/wailsapp/wails/pull/6275), [#6276](https://github.com/wailsapp/wails/pull/6276)
- **Cache build menangani lebih banyak proyek dunia nyata.** Spasi yang diapit tanda kutip atau di-escape tidak lagi menyebabkan cache miss pada output Go, dan penyimpanan cache secara bersamaan kini diserialkan untuk mencegah kerusakan berkas. [#6495](https://github.com/wailsapp/wails/pull/6495), [#6488](https://github.com/wailsapp/wails/pull/6488)
- **Penyajian aset mengurangi pekerjaan yang tidak perlu.** Ekstensi berkas yang dikenali melewati deteksi MIME, dan pencatatan debug yang dinonaktifkan tidak lagi menyiapkan pesan yang nantinya dibuang. Kedua perubahan ini mengurangi pekerjaan pada jalur permintaan. [#6414](https://github.com/wailsapp/wails/pull/6414), [#6420](https://github.com/wailsapp/wails/pull/6420)
- **Pengemasan MSIX menghindari buffer sebesar berkas biner.** Berkas executable dialirkan ke dalam paket, bukan dibaca seluruhnya ke memori. Sebelumnya, penyalinan representatif sebesar 16 MiB mengalokasikan sekitar 16.8 MB per operasi. [#6472](https://github.com/wailsapp/wails/pull/6472)

## Silakan coba

Kedua fitur ini masih eksperimental. Cobalah build lintas platform, kemas aplikasi Anda, dan beri tahu kami bagian mana yang masih menyulitkan. Terima kasih kepada semua orang yang telah menguji, melaporkan masalah, dan berkontribusi.

Penghargaan khusus untuk beesophie dari komunitas Discord kami, yang membangun [Cook](https://github.com/SimTXNG/cook), sebuah sistem build luar biasa yang menunjukkan beberapa kekurangan dalam sistem kami sendiri. Melihat berbagai hal yang dimungkinkan oleh Cook membantu kami menyadari bagian mana yang dapat kami permudah dalam proses build dengan Wails dan benar-benar menginspirasi pekerjaan ini.

Terima kasih, Beesophie!
