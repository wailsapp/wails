---
title: "Peta Jalan"
description: "Status proyek Wails v3, fitur yang direncanakan, dan cara berkontribusi"
slug: "status"
sourcePath: "status.md"
---

## Wails v3.0.0 — Stabil

Lihat [Catatan Perubahan](/changelog/) untuk mengetahui status terbaru.

Dokumentasi ini membahas Wails v3.0.0. Setelah v3.0.0, perbaikan dan penyempurnaan dirilis setiap minggu, diselingi pratinjau nightly bertanggal serta hotfix untuk masalah mendesak. Lihat [kebijakan rilis](/contributing/release-policy/).

## Kompatibilitas Desktop

Kontrak kompatibilitas v3.0.0 mencakup aplikasi desktop:

| Platform | Target yang didukung | Persyaratan dan catatan |
| --- | --- | --- |
| Windows | amd64 dan arm64 | Runtime WebView2 |
| macOS | Intel dan Apple Silicon | Versi macOS dan WebKit yang tercantum dalam panduan instalasi |
| Linux | amd64 dan arm64 | Secara default, GTK4 + WebKitGTK 6.0; GTK3 + WebKit2GTK 4.1 tetap menjadi opsi lama `-tags gtk3` hingga v3.0.x dan dihapus pada v3.1 |

Untuk pengembangan, semua target memerlukan Go 1.26 atau versi yang lebih baru. Dukungan Android dan iOS masih bersifat eksperimental dan berada di luar kontrak kompatibilitas desktop stabil.

## Cara Anda Dapat Berkontribusi

- Uji Wails v3.0.0 dan laporkan bug yang dapat direproduksi
- Berkontribusi pada dokumentasi dan contoh
- Berpartisipasi dalam diskusi dan memberikan masukan tentang draf WEP
- Kirim pull request untuk perbaikan bug, dokumentasi, atau WEP yang telah diterima

Kami menyambut kontribusi dari komunitas. Jika Anda ingin membantu mencapai tujuan ini, bergabunglah dalam diskusi komunitas. Proposal untuk fungsionalitas baru harus diajukan sebagai PR WEP, bukan sebagai issue permintaan fitur.

## Masukan dan Pembaruan

Peta jalan ini dapat berubah berdasarkan masukan komunitas dan prioritas proyek. Kami akan memperbaruinya secara berkala untuk mencerminkan kemajuan dan perubahan arah. Laporkan masalah yang dapat direproduksi sebagai issue; usulkan fungsionalitas baru melalui PR [WEP (Wails Enhancement Proposal)](https://github.com/wailsapp/wails/blob/releases/v3.0.0/v3/wep/README.md).
