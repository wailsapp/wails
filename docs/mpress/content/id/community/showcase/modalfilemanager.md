---
title: "Modal File Manager"
description: "Aplikasi desktop yang dibangun dengan Wails"
slug: "community/showcase/modalfilemanager"
sourcePath: "community/showcase/modalfilemanager.md"
---

![Modal File Manager](/assets/showcase-images/modalfilemanager.webp)

[Modal File Manager](https://github.com/raguay/ModalFileManager) adalah pengelola berkas dua panel yang menggunakan teknologi web. Desain awal saya berbasis NW.js dan dapat ditemukan [di sini](https://github.com/raguay/ModalFileManager-NWjs). Versi ini menggunakan kode frontend berbasis Svelte yang sama (tetapi telah banyak dimodifikasi sejak beralih dari NW.js), sedangkan backend-nya merupakan implementasi [Wails 2](https://wails.io/). Dengan implementasi ini, saya tidak lagi menggunakan perintah baris perintah `rm`, `cp`, dan sebagainya, tetapi git harus terinstal pada sistem untuk mengunduh tema dan ekstensi. Aplikasi ini sepenuhnya dikembangkan menggunakan Go dan berjalan jauh lebih cepat daripada versi-versi sebelumnya.

Pengelola berkas ini dirancang berdasarkan prinsip yang sama dengan Vim: tindakan papan ketik yang dikendalikan oleh status. Jumlah status tidak tetap, tetapi sangat mudah diprogram. Karena itu, konfigurasi papan ketik dalam jumlah tak terbatas dapat dibuat dan digunakan. Inilah perbedaan utamanya dibandingkan pengelola berkas lain. Tema dan ekstensi tersedia untuk diunduh dari GitHub.
