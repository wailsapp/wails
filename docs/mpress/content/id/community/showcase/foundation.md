---
title: "Foundation"
description: "Kerangka aplikasi desktop untuk Wails v3.1 dengan React 19 dan Mantine 9"
slug: "community/showcase/foundation"
sourcePath: "community/showcase/foundation.md"
---

![Pengaturan tata letak Foundation dengan sembilan preset yang dapat dialihkan](/assets/showcase-images/foundation-layouts.webp)

![Halaman beranda Foundation dengan tema gelap](/assets/showcase-images/foundation-home-dark.webp)

![Grafik X-Pro Foundation dengan tema Obsidian](/assets/showcase-images/foundation-charts.webp)

**[Foundation](https://github.com/FxRayHughes/foundation)** adalah **kerangka aplikasi desktop untuk Wails v3.1**, dibangun dengan **Go, React 19, TypeScript, dan Mantine 9**. Kerangka ini menyediakan backend Go berlapis, tampilan jendela khusus tiap platform, persistensi SQLite, tema, internasionalisasi, dan pengujian end-to-end sebagai titik awal untuk proyek baru.

## Fitur

- **Jendela khusus tiap platform** — jendela tanpa bingkai dengan bilah judul yang digambar oleh frontend di Windows dan Linux; di macOS, bilah judul imersif yang tetap mempertahankan tombol jendela bawaan.
- **Sembilan preset tata letak** — bilah ikon, bilah samping, navigasi atas, navigasi di bilah judul, bilah kanan, dock bawah, dock mengambang yang dapat diseret, kartu inset, dan mode fokus di atas kerangka CSS Grid; preset kustom dapat didaftarkan.
- **Tema dan internasionalisasi** — preset terang, gelap, dan Obsidian, serta editor palet kustom; antarmuka dalam bahasa Inggris dan Mandarin; tema dan bahasa tetap sinkron di jendela utama, jendela anak, dan panel tray.
- **Persistensi SQLite** — GORM AutoMigrate tanpa file migrasi, lokasi database yang dapat dipilih pengguna, dan grafik penggunaan per tabel.
- **Integrasi desktop** — panel system tray, jendela anak, dialog file dan pesan native, serta ekspor CSV / XLSX melalui dialog native "Simpan Sebagai".
- **Utilitas** — klien HTTP, subproses yang terikat ke Windows Job Object, enkripsi AES-GCM, log dengan rotasi, dan penulisan atomik.
- **Pengujian end-to-end** — Playwright dijalankan terhadap build Wails mode server (backend Go dan frontend sungguhan); rangkaian pengujian yang sama menghasilkan tangkapan layar README.

Frontend memanggil layanan Go melalui binding yang dihasilkan Wails, dan setiap platform dibangun menjadi satu biner. Frontend mengikuti struktur MVVM, dan README menyertakan panduan langkah demi langkah untuk mengganti nama kerangka ini menjadi proyek Anda sendiri.

Berjalan di Windows, macOS, dan Linux.

[GitHub](https://github.com/FxRayHughes/foundation) |
[Rilis](https://github.com/FxRayHughes/foundation/releases/latest)
