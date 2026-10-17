---
title: "Mengapa Wails?"
description: "Pahami mengapa Wails merupakan pilihan yang tepat untuk aplikasi desktop Anda"
slug: "quick-start/why-wails"
sourcePath: "quick-start/why-wails.md"
---

Wails memadukan **performa dan kesederhanaan Go** dengan **fleksibilitas UI web modern**, sehingga Anda dapat membuat aplikasi desktop native yang indah dengan alat yang sudah Anda kuasai.

## Performa yang Dapat Dirasakan Pengguna

**Aplikasi Wails:**

- **Biner berukuran ~15 MB** (dibandingkan dengan Electron yang berukuran 150 MB)
- **Penggunaan memori dasar ~10 MB** (dibandingkan dengan Electron yang menggunakan 100 MB+)
- Waktu mulai **&lt;0.5 detik** (dibandingkan dengan 2-3 detik pada Electron)
- **Rendering native** menggunakan WebView yang disediakan OS

Pengguna akan merasakan bahwa aplikasi Anda cepat, ringan, dan profesional.

## Pengalaman Developer

**Tulis Sekali, Jalankan di Mana Saja:**

- Satu basis kode Go untuk Windows, macOS, dan Linux
- Gunakan framework web apa pun (React, Vue, Svelte, JavaScript biasa)
- Hot reload selama pengembangan
- Binding TypeScript dibuat secara otomatis dari kode Go

Rilis lebih cepat dengan lebih sedikit kode yang perlu dipelihara.

## Fitur Siap Produksi

**Semua yang Anda perlukan:**

- Beberapa jendela dengan siklus hidup masing-masing
- Menu native (aplikasi, konteks, baki sistem)
- Dialog file dengan UI native platform
- Integrasi sistem (notifikasi, papan klip, pintasan papan ketik)
- Penandatanganan kode dan pengemasan untuk semua platform

Bangun aplikasi profesional, bukan purwarupa.

## Pengembangan Lebih Cepat

- **Satu basis kode, tiga platform** - Tulis sekali, buat build untuk Windows, macOS, dan Linux
- **Gunakan keahlian yang sudah Anda miliki** - Go untuk backend, HTML/CSS/JS untuk UI
- **Umpan balik seketika** - Hot reload selama pengembangan, dengan waktu kompilasi yang hanya beberapa detik
- **Biner berukuran kecil** - Aplikasi berukuran 15 MB memungkinkan proses build, pengunduhan, dan iterasi yang lebih cepat

## Kapan Sebaiknya Memilih Wails

**Wails Sangat Cocok Untuk:**

- **Aplikasi bisnis** (CRM, inventaris, dasbor, alat administrasi)
- **Alat developer** (klien basis data, penguji API, alat deployment)
- **Aplikasi produktivitas** (pencatatan, pengelola tugas, pelacak waktu)
- **Alat kreatif** (editor gambar, pemroses video, utilitas desain)
- **Alat internal** (aplikasi khusus perusahaan, alat otomatisasi)

## Kisah Sukses di Dunia Nyata

@note{type="tip" title="Aplikasi Produksi"}
Wails mendukung aplikasi nyata yang digunakan oleh ribuan pengguna:

- **Alat pengelolaan basis data** dengan UI yang kompleks
- **Dasbor keuangan** yang memproses data waktu nyata
- **Alat pengeditan video** dengan performa native
- **Utilitas pengembangan** yang digunakan oleh tim engineering

[Lihat galeri →](/community/showcase/)

@end

## Cara Kerja Wails

Tidak seperti Electron yang menyertakan seluruh browser dan runtime Node.js, Wails menggunakan pendekatan yang secara mendasar berbeda: kode Go Anda dikompilasi menjadi biner native, sedangkan UI Anda berjalan di WebView bawaan sistem operasi. Arsitektur ini menghasilkan biner berukuran kecil, waktu mulai yang cepat, dan penggunaan memori yang rendah sehingga aplikasi Wails terasa native.

### Arsitektur

Aplikasi Wails terdiri atas dua bagian utama yang berkomunikasi tanpa hambatan: backend Go yang menangani logika bisnis dan operasi sistem, serta frontend berbasis web untuk antarmuka pengguna Anda. WebView yang disediakan OS merender UI Anda tanpa menyertakan browser, sedangkan lapisan binding menyediakan komunikasi yang aman terhadap tipe antara Go dan JavaScript.

<div style="display: flex; justify-content: center; align-items: center; margin: 2rem 0;">
  <img src="/img/architecture.svg" alt="Arsitektur Wails: backend Go dan UI web Anda dikompilasi menjadi satu biner native, dihubungkan melalui binding yang dihasilkan, dan dirender oleh WebView sistem operasi" style="max-width: 640px; width: 100%;" />
</div>

Arsitektur sederhana ini memungkinkan kode JavaScript memanggil fungsi Go secara langsung (melalui binding yang dibuat secara otomatis), sedangkan Go dapat mengirimkan event dan data kembali ke frontend. Kedua lapisan berkomunikasi melalui bridge dalam memori yang efisien dengan overhead kurang dari satu milidetik.

**Cara Wails mencapai performa tinggi:**

1. **Tidak menyertakan runtime** - Menggunakan biner Go yang telah dikompilasi
2. **WebView native** - Mesin rendering yang disediakan OS
3. **Bridge langsung Go ↔ JS** - Komunikasi dalam memori tanpa overhead jaringan
4. **Biner terkompilasi** - Mulai seketika tanpa kompilasi JIT

## Langkah Berikutnya

Setelah memahami apa yang disediakan Wails, mari siapkan lingkungan Anda:

1. **Instal Wails** - Siapkan lingkungan pengembangan Anda dalam 5 menit [Panduan Instalasi →](/quick-start/installation/)

2. **Bangun Aplikasi Pertama Anda** - Buat aplikasi yang berfungsi dan pahami dasar-dasarnya [Tutorial Aplikasi Pertama →](/quick-start/first-app/)

3. **Jelajahi Fitur** - Temukan kemampuan Wails untuk aplikasi Anda [Ikhtisar Fitur →](/quick-start/next-steps/)

---

**Masih punya pertanyaan?** Bergabunglah dengan [komunitas Discord](https://discord.gg/JDdSxwjhGf) kami dan tanyakan langsung kepada tim.
