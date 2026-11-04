---
title: "Langkah Berikutnya"
description: "Langkah selanjutnya setelah membuat aplikasi pertama Anda"
slug: "quick-start/next-steps"
sourcePath: "quick-start/next-steps.md"
---

Anda telah membuat aplikasi Wails pertama Anda. Ikuti jalur pembelajaran di bawah ini, atau langsung buka panduan jika Anda sudah tahu apa yang ingin dibuat.

## Pelajari Lebih Lanjut

Lanjutkan belajar dengan membuat satu aplikasi lengkap pada setiap tahap:

1. **[Buat layanan kode QR](/tutorials/01-creating-a-service/).** Berlatihlah mengekspos logika Go ke frontend Anda melalui layanan.
2. **[Buat daftar TODO](/tutorials/02-todo-vanilla/).** Tambahkan, ubah, dan hapus data aplikasi melalui pemanggilan layanan.
3. **[Buat aplikasi Catatan](/tutorials/03-notes-vanilla/).** Bekerjalah dengan file, dialog native, dan data tersimpan.

Setelah siap mempelajari rilis, ikuti [tutorial aplikasi yang memperbarui dirinya sendiri](/tutorials/04-self-update-a-wails-app/). Lihat [semua tutorial](/tutorials/overview/) untuk mengetahui hasil pembelajaran setiap materi.

## Fitur

**Sudah membuat aplikasi Anda sendiri?** Gunakan panduan praktis untuk tugas yang sedang Anda kerjakan:

| Saya ingin… | Mulai di sini |
| --- | --- |
| Memanggil Go dari frontend saya | [Memanggil Go dari JavaScript](/features/bindings/methods/) |
| Membuka atau mengelola jendela | [Dasar-dasar jendela](/features/windows/basics/) |
| Menambahkan menu native atau ikon baki sistem | [Menu aplikasi](/features/menus/application/) atau [menu baki sistem](/features/menus/systray/) |
| Meminta pengguna membuka atau menyimpan file | [Dialog file](/features/dialogs/file/) |
| Mengirim peristiwa antarkomponen | [Sistem peristiwa](/features/events/system/) |
| Menyalin teks atau menerima file yang dilepas | [Papan klip](/features/clipboard/basics/) atau [seret dan lepas file](/features/drag-and-drop/files/) |
| Menambahkan pintasan papan ketik | [Pintasan papan ketik](/features/keyboard/shortcuts/) |
| Membangun dan mengemas aplikasi saya | [Membangun aplikasi](/guides/build/building/), lalu mengemasnya untuk [Windows](/guides/build/windows/), [macOS](/guides/build/macos/), atau [Linux](/guides/build/linux/) |

Halaman-halaman ini berada di bagian **Panduan Praktis** pada bilah sisi. Anda dapat membuka topik yang relevan tanpa harus membaca seluruh bagian secara berurutan.

## Pahami Cara Kerja Wails

Gunakan **Konsep Inti** saat Anda ingin memahami desain yang mendasari kode:

- [Cara kerja Wails](/concepts/architecture/) menjelaskan hubungan antara backend Go, frontend, dan WebView native.
- [Jembatan Go–frontend](/concepts/bridge/) menjelaskan cara pemanggilan diteruskan antara Go dan JavaScript.
- [Siklus hidup aplikasi](/concepts/lifecycle/) menjelaskan proses memulai dan menghentikan aplikasi.
- [Sistem build](/concepts/build-system/) menjelaskan peran tugas proyek dan alat frontend.

## Cari API

Gunakan [Referensi API](/reference/overview/) saat Anda memerlukan tipe, metode, opsi, atau perintah, bukan materi pembelajaran. Untuk bekerja melalui baris perintah, mulailah dengan [ikhtisar perintah CLI](/reference/cli/); gunakan [flag dan contoh CLI](/guides/cli/) untuk melihat detail setiap perintah.

## Dapatkan Bantuan

Periksa [pertanyaan umum](/faq/), telusuri [contoh yang dapat dijalankan](https://github.com/wailsapp/wails/tree/releases/v3-beta/v3/examples), atau ajukan pertanyaan di [Discord](https://discord.gg/JDdSxwjhGf). Masalah khusus platform dikelompokkan di bagian **Pemecahan Masalah** pada bilah sisi.
