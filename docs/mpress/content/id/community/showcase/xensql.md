---
title: "XenSQL"
description: "Workbench SQL yang mengutamakan penyimpanan lokal, dengan dukungan berbagai basis data dan alat kueri canggih"
slug: "community/showcase/xensql"
sourcePath: "community/showcase/xensql.md"
---

![Editor tabel XenSQL dengan pengeditan langsung](/assets/showcase-images/xensql-1.png) ![Editor kueri XenSQL dengan saran tabel dan kolom](/assets/showcase-images/xensql-2.png) ![Beberapa kueri transaksi di XenSQL](/assets/showcase-images/xensql-3.png)

**[XenSQL](https://github.com/Bare7a/XenSQL)** adalah **workbench SQL desktop yang cepat dan mengutamakan penyimpanan lokal**, yang dibuat dengan **Go, Wails, dan React**. Aplikasi ini menyatukan PostgreSQL, MySQL/MariaDB, dan SQLite dalam satu antarmuka yang rapi dan terasa natif—tanpa cloud, telemetri, ataupun akun.

## Keunggulan Utama

- **Editor SQL yang Canggih**—berbasis Monaco dengan pelengkapan otomatis cerdas yang memahami skema, eksekusi beberapa pernyataan, hasil streaming, dan tab hasil untuk setiap pernyataan
- **Penampil Data Tingkat Lanjut**—pemeriksa JSON interaktif, editor sel yang memahami sintaks (JSON, XML, HTML, teks), pengeditan langsung, dan pemeriksaan rekaman secara menyeluruh
- **Pengeditan Data yang Mulus**—telusuri tabel, siapkan perubahan secara langsung, lakukan operasi massal, serta gunakan `INSERT`/`UPDATE`/`DELETE` dengan aman dan didukung oleh `RETURNING`
- **Fitur Produktivitas**—penjelajah skema, kueri tersimpan, riwayat kueri, pencarian cepat (`Ctrl+P`), dan alur kerja yang mengutamakan keyboard
- **Opsi Ekspor**—CSV, JSON, Markdown, SQL INSERT

Sepenuhnya luring dan portabel. Semuanya disimpan secara lokal dalam satu folder `XenSQL-data/` yang dapat dibawa bersama aplikasi.

**Basis Data yang Didukung**: PostgreSQL, MySQL, MariaDB, dan SQLite (dengan mode hanya-baca dan opsi koneksi aman).

Dirancang bagi pengembang yang menginginkan kecepatan, kejelasan, dan kendali tanpa beban berlebihan dari alat SQL tradisional.
