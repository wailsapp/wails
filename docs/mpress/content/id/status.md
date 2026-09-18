---
title: "Status proyek"
description: "Kompatibilitas beta Wails v3, dukungan keamanan, dan panduan peningkatan versi"
slug: "status"
sourcePath: "status.md"
---

## Status Saat Ini: Beta

Lihat [Catatan Perubahan](/changelog/) untuk mengetahui status terbaru.

Tujuan kami adalah rilis v3.0 yang stabil. Wails v2 tetap menjadi rilis stabil saat ini dan terus menerima perbaikan. Uji rilis beta dengan aplikasi Anda sebelum penerapan.

## Jaminan Kompatibilitas Beta

Kontrak Beta v3 mencakup aplikasi desktop:

| Platform | Target yang didukung | Persyaratan dan catatan |
| --- | --- | --- |
| Windows | amd64 dan arm64 | Runtime WebView2 |
| macOS | Intel dan Apple Silicon | Versi macOS dan WebKit yang tercantum dalam panduan instalasi |
| Linux | amd64 dan arm64 | GTK4 + WebKitGTK 6.0 secara default; GTK3 + WebKit2GTK 4.1 tetap tersedia sebagai opsi lama `-tags gtk3` hingga v3.0.x dan dihapus pada v3.1 |

Semua target memerlukan Go 1.25 atau yang lebih baru untuk pengembangan. Dukungan Android dan iOS bersifat eksperimental dan tidak menghambat Beta desktop. API Beta ditujukan agar stabil, tetapi cacat prarilis dan perubahan yang diumumkan secara eksplisit masih dapat diperbaiki sebelum v3.0.0.

## Cara Anda Dapat Berkontribusi

- Uji rilis beta terbaru dan laporkan bug yang dapat direproduksi
- Berkontribusi pada dokumentasi dan contoh
- Berpartisipasi dalam diskusi dan berikan umpan balik tentang draf WEP
- Kirim pull request untuk perbaikan bug, dokumentasi, atau WEP yang telah diterima

Kami menyambut kontribusi dari komunitas. Jika Anda ingin membantu mencapai tujuan ini, bergabunglah dalam diskusi komunitas. Proposal untuk fungsionalitas baru harus diajukan melalui PR WEP, bukan issue permintaan fitur.

## Umpan Balik dan Pembaruan

Laporkan masalah yang dapat direproduksi sebagai issue; usulkan fungsionalitas baru melalui PR [WEP (Wails Enhancement Proposal)](https://github.com/wailsapp/wails/blob/master/v3/wep/README.md).

## Menggunakan beta

Tetapkan versi yang pasti untuk CLI, modul Go, dan runtime frontend alih-alih mengikuti `latest`. Untuk proyek alpha yang sudah ada, ikuti [panduan peningkatan dari alpha ke beta](/migration/alpha-to-beta/).

[Kebijakan keamanan](https://github.com/wailsapp/wails/blob/master/SECURITY.md) mencantumkan rilis beta v3 sebagai versi yang didukung dan rilis alpha sebagai versi yang tidak didukung. Laporkan kerentanan melalui [pelaporan kerentanan privat](https://github.com/wailsapp/wails/security/advisories/new), bukan issue publik.

## Pekerjaan yang dilacak

- [Bug terbuka berlabel v3](https://github.com/wailsapp/wails/issues?q=is%3Aissue+is%3Aopen+label%3ABug+label%3Av3)
- [Issue v3 terbuka berlabel P0 atau P1](https://github.com/wailsapp/wails/issues?q=is%3Aissue+is%3Aopen+label%3Av3+label%3AP0%2CP1)
- [Milestone rilis](https://github.com/wailsapp/wails/milestones)

Kueri langsung ini bergantung pada label issue; hasilnya bukan daftar lengkap atau janji tanggal maupun cakupan rilis. Baca issue untuk menilai dampaknya pada proyek Anda.
