---
title: "Gambaran Umum Teknis"
description: "Arsitektur tingkat tinggi dan peta jalan untuk basis kode Wails v3"
slug: "contributing/overview"
sourcePath: "contributing/overview.md"
---

## Selamat Datang di Dokumentasi Teknis Wails v3

Bagian ini **bukan** membahas pedoman komunitas atau cara membuka pull request. Sebaliknya, bagian ini mengulas secara mendalam **cara Wails v3 dibangun** agar Anda dapat dengan cepat memahami basis kode dan mulai mengutak-atiknya dengan percaya diri.

Baik Anda berencana memperbaiki runtime, memperluas CLI, membuat templat baru, maupun sekadar memahami mekanisme internalnya, halaman-halaman berikut menyediakan konteks teknis yang Anda perlukan.

---

## Arsitektur Tingkat Tinggi

@cards{cols="2"}
◇ Backend Go
Inti setiap aplikasi Wails adalah kode Go yang dikompilasi menjadi berkas eksekutabel native. Backend ini menangani logika aplikasi, integrasi sistem, dan operasi yang sangat memengaruhi performa.

---
▤ Frontend Web
UI ditulis dengan teknologi web standar (React, Vue, Svelte, Vanilla, …) dan dirender oleh WebView sistem yang ringan (WebKit di Linux/macOS, WebView2 di Windows).

---
◆ Lapisan Penghubung
Penghubung dalam memori tanpa penyalinan memungkinkan pemanggilan **Go⇄JavaScript** dengan konversi tipe otomatis, penerusan peristiwa, dan penerusan galat.

---
▸ CLI & Peralatan
`wails3` mengatur pembuatan proyek, server pengembangan dengan pemuatan ulang langsung, penggabungan aset, kompilasi silang, dan pemaketan (deb, rpm, AppImage, msi, dmg…).

@end

---

## Gambaran Umum Arsitektur

**Wails v3 – Alur Menyeluruh**

**[Placeholder Diagram Alur Menyeluruh]**

Diagram tersebut menunjukkan **alur menyeluruh**:

1. **CLI** mengendalikan pembuatan kode, server pengembangan, kompilasi, dan pemaketan.\
2. **Sistem Binding** menghasilkan kode penghubung yang memungkinkan **Frontend Web** memanggil **Backend Go**.\
3. Selama pengembangan, **Server Aset** meneruskan permintaan melalui proksi ke server pengembangan framework; dalam produksi, server ini menyajikan berkas yang disematkan.\
4. Saat runtime, **Runtime Desktop** mengelola jendela dan API sistem operasi, sedangkan **Penghubung** mengirimkan pesan antara Go dan JavaScript.

---

## Cakupan Dokumentasi Ini

| Topik | Alasan Pentingnya |
| --- | --- |
| **Struktur Basis Kode** | Peta direktori `/v3` dan cara modul berinteraksi. |
| **Mekanisme Internal Runtime** | Pengelolaan jendela, API sistem, pemroses pesan, dan lapisan kompatibilitas platform. |
| **Server Aset & Pengembangan** | Cara aset web disajikan dalam pengembangan dan disematkan dalam produksi. |
| **Alur Build & Pemaketan** | Alur kerja berbasis Taskfile, kompilasi lintas platform, dan pembuatan penginstal. |
| **Sistem Binding** | Alur analisis statis yang menghasilkan binding Go⇄TS yang aman secara tipe. |
| **Sistem Templat** | Arsitektur generator yang mendukung `wails3 init -t <framework>`. |
| **Pengujian & CI** | Kerangka pengujian unit/integrasi, GitHub Actions, dan panduan pendeteksi race condition. |
| **Memperluas Wails** | Menambahkan layanan, templat, atau subperintah CLI. |

Setiap halaman berikutnya mengulas area-area ini secara mendalam dengan contoh kode konkret, diagram, dan referensi ke berkas sumber yang relevan.

---

@note{type="info"}
Prasyarat: Anda sebaiknya sudah terbiasa dengan **Go 1.26+**, TypeScript dasar, dan alat build frontend modern. Jika Anda baru mengenal Go, pertimbangkan untuk membaca sekilas tur resminya terlebih dahulu.

@end

Selamat menjelajah—dan selamat datang di mekanisme internal Wails v3!
