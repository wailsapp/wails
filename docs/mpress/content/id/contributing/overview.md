---
title: "Ikhtisar Teknis"
description: "Arsitektur tingkat tinggi dan peta jalan untuk memahami basis kode Wails v3"
slug: "contributing/overview"
sourcePath: "contributing/overview.md"
---

## Selamat Datang di Dokumentasi Teknis Wails v3

Bagian ini **bukan** tentang pedoman komunitas atau cara membuka pull request. Sebaliknya, bagian ini membahas secara mendalam **cara Wails v3 dibangun** agar Anda dapat dengan cepat memahami basis kode dan mulai mengutak-atiknya dengan percaya diri.

Baik Anda berencana memperbaiki runtime, memperluas CLI, membuat templat baru, maupun sekadar memahami cara kerja internalnya, halaman-halaman berikut menyediakan konteks teknis yang Anda perlukan.

---

## Arsitektur Tingkat Tinggi

@cards{cols="2"}
◇ Backend Go
Inti setiap aplikasi Wails adalah kode Go yang dikompilasi menjadi berkas executable native. Kode ini menangani logika aplikasi, integrasi sistem, dan operasi yang kritis terhadap performa.

---
▤ Frontend Web
UI ditulis dengan teknologi web standar (React, Vue, Svelte, Vanilla, …) dan dirender oleh WebView sistem yang ringan (WebKit di Linux/macOS, WebView2 di Windows).

---
◆ Lapisan Penghubung
Bridge dalam memori tanpa penyalinan memungkinkan pemanggilan **Go⇄JavaScript** dengan konversi tipe otomatis, propagasi event, dan penerusan error.

---
▸ CLI & Peralatan
`wails3` mengatur pembuatan proyek, server pengembangan dengan pemuatan ulang otomatis, bundling aset, kompilasi silang, dan pengemasan (deb, rpm, AppImage, msi, dmg…).

@end

---

## Ikhtisar Arsitektur

**Wails v3 – Alur Menyeluruh**

**[Placeholder Diagram Alur Menyeluruh]**

Diagram tersebut menunjukkan **alur menyeluruh**:

1. **CLI** mengendalikan pembuatan kode, server pengembangan, kompilasi, dan pengemasan.\
2. **Sistem Binding** menghasilkan kode penghubung yang memungkinkan **Frontend Web** memanggil **Backend Go**.\
3. Selama pengembangan, **Server Aset** meneruskan permintaan ke server pengembangan framework; dalam produksi, server ini menyajikan berkas yang disematkan.\
4. Saat runtime, **Runtime Desktop** mengelola jendela dan API sistem operasi, sedangkan **Bridge** mengirimkan pesan antara Go dan JavaScript.

---

## Cakupan Dokumentasi Ini

| Topik | Mengapa Ini Penting |
| --- | --- |
| **Tata Letak Basis Kode** | Peta direktori `/v3` dan cara modul berinteraksi. |
| **Cara Kerja Internal Runtime** | Pengelolaan jendela, API sistem, pemroses pesan, dan shim platform. |
| **Aset & Server Pengembangan** | Cara aset web disajikan saat pengembangan dan disematkan dalam produksi. |
| **Pipeline Build & Pengemasan** | Alur kerja berbasis Taskfile, kompilasi lintas platform, dan pembuatan installer. |
| **Sistem Binding** | Pipeline analisis statis yang menghasilkan binding Go⇄TS yang aman terhadap tipe. |
| **Sistem Templat** | Arsitektur generator yang mendukung `wails3 init -t <framework>`. |
| **Pengujian & CI** | Kerangka pengujian unit/integrasi, GitHub Actions, dan panduan race detector. |
| **Memperluas Wails** | Menambahkan layanan, templat, atau subperintah CLI. |

Setiap halaman berikutnya membahas area-area ini secara mendalam dengan contoh kode konkret, diagram, dan referensi ke berkas sumber yang relevan.

---

@note{type="info"}
Prasyarat: Anda sebaiknya memahami **Go 1.25+**, TypeScript dasar, dan alat build frontend modern. Jika Anda baru mengenal Go, pertimbangkan untuk membaca sekilas tur resminya terlebih dahulu.

@end

Selamat menjelajah—dan selamat datang di cara kerja internal Wails v3!
