---
title: "Peta Jalan"
description: "Status proyek Wails v3, fitur yang direncanakan, dan cara berkontribusi"
slug: "status"
sourcePath: "status.md"
---

## Wails v3.0.0 — Stabil

Lihat [Catatan Perubahan](/changelog/) untuk mengetahui status terbaru.

Dokumentasi ini mencakup Wails v3.0.0. Perbaikan bug diprioritaskan berdasarkan tingkat keparahan dan dikirim melalui pembaruan patch setiap malam. Fitur baru ditujukan untuk v3.1.0; perbaikan bug tetap terpisah dari rilis fitur tersebut.

## Kompatibilitas desktop

Jaminan kompatibilitas v3.0.0 mencakup aplikasi desktop:

| Platform | Target yang didukung | Persyaratan dan catatan |
| --- | --- | --- |
| Windows | amd64 dan arm64 | Runtime WebView2 |
| macOS | Intel dan Apple Silicon | Versi macOS dan WebKit yang tercantum dalam panduan instalasi |
| Linux | amd64 dan arm64 | GTK4 + WebKitGTK 6.0 secara default; GTK3 + WebKit2GTK 4.1 tetap tersedia sebagai opsi lama `-tags gtk3` hingga v3.0.x dan dihapus pada v3.1 |

Semua target memerlukan Go 1.25 atau yang lebih baru untuk pengembangan. Dukungan Android dan iOS tetap eksperimental dan berada di luar jaminan kompatibilitas desktop yang stabil.

## Cara Anda Dapat Berkontribusi

- Uji Wails v3.0.0 dan laporkan bug yang dapat direproduksi
- Berkontribusi pada dokumentasi dan contoh
- Berpartisipasi dalam diskusi dan berikan umpan balik tentang draf WEP
- Kirim pull request untuk perbaikan bug, dokumentasi, atau WEP yang telah diterima

Kami menyambut kontribusi dari komunitas. Jika Anda ingin membantu mencapai tujuan ini, bergabunglah dalam diskusi komunitas. Proposal untuk fungsionalitas baru harus diajukan melalui PR WEP, bukan issue permintaan fitur.

## Umpan Balik dan Pembaruan

Peta jalan ini dapat berubah berdasarkan umpan balik komunitas dan prioritas proyek. Kami akan memperbaruinya secara rutin untuk mencerminkan kemajuan dan perubahan arah. Laporkan masalah yang dapat direproduksi sebagai issue; usulkan fungsionalitas baru melalui PR [WEP (Wails Enhancement Proposal)](https://github.com/wailsapp/wails/blob/releases/v3.0.0/v3/wep/README.md).
