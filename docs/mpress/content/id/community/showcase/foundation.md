---
title: "Foundation"
description: "Kerangka aplikasi desktop Wails v3 yang siap produksi"
slug: "community/showcase/foundation"
sourcePath: "community/showcase/foundation.md"
---

![Pengaturan tata letak Foundation dengan sembilan preset yang dapat dialihkan](/assets/showcase-images/foundation-layouts.webp)

![Halaman beranda Foundation dengan tema gelap](/assets/showcase-images/foundation-home-dark.webp)

![Grafik X-Pro Foundation dengan tema Obsidian](/assets/showcase-images/foundation-charts.webp)

**[Foundation](https://github.com/FxRayHughes/foundation)** adalah **kerangka aplikasi desktop untuk Wails v3** yang siap produksi, dibangun dengan **Go, React 19, TypeScript, dan Mantine 9**. Kerangka ini menyediakan bagian-bagian yang biasanya butuh berminggu-minggu untuk dibangun — backend Go berlapis, tampilan jendela khusus tiap platform, persistensi, tema, internasionalisasi, dan pengujian end-to-end — sehingga Anda bisa langsung fokus pada logika bisnis.

## Fitur utama

- **Jendela native di setiap platform** — bilah judul tanpa bingkai yang digambar oleh frontend di Windows dan Linux; bilah judul imersif di macOS yang tetap mempertahankan tombol jendela bawaan.
- **Sembilan tata letak yang dapat dialihkan** — bilah ikon, bilah samping, navigasi atas, navigasi di bilah judul, bilah kanan, dock bawah, dock mengambang yang dapat diseret, kartu inset, dan mode fokus, semuanya di atas kerangka CSS Grid; preset kustom dapat didaftarkan dengan satu baris.
- **Tema dan internasionalisasi** — preset terang, gelap, dan Obsidian, ditambah editor palet kustom; bahasa Inggris dan Mandarin tersedia langsung; tema dan bahasa tetap sinkron di jendela utama, jendela anak, dan panel tray.
- **Persistensi SQLite** — GORM AutoMigrate tanpa file migrasi, lokasi database yang dapat dipilih pengguna, dan grafik penggunaan per tabel.
- **Integrasi desktop** — panel system tray, jendela anak, dialog file dan pesan native, serta ekspor CSV / XLSX melalui dialog native "Simpan Sebagai".
- **Utilitas yang andal** — klien HTTP, subproses yang terikat ke Windows Job Object, enkripsi AES-GCM, log dengan rotasi, dan penulisan atomik.
- **Teruji end-to-end** — Playwright menjalankan build Wails mode server (backend Go dan frontend sungguhan), dan rangkaian pengujian yang sama menghasilkan tangkapan layar README.

Dengan Wails, semuanya dikemas dalam satu biner untuk tiap platform: layanan Go dipanggil langsung melalui binding yang dihasilkan, sementara frontend mengikuti struktur MVVM yang ketat sehingga tetap mudah dirawat seiring bertambahnya fitur.

Tersedia untuk Windows, macOS, dan Linux.

[GitHub](https://github.com/FxRayHughes/foundation) |
[Rilis](https://github.com/FxRayHughes/foundation/releases/latest)
