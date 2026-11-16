---
title: "Integrasi Basis Data"
description: "Mengintegrasikan SQLite, PostgreSQL, dan basis data lainnya dengan Wails v3"
slug: "guides/patterns/database"
sourcePath: "guides/patterns/database.md"
---

@note{type="info"}
Halaman ini adalah halaman sementara. Konten lengkap tentang integrasi basis data akan segera tersedia.

@end

Aplikasi Wails v3 adalah aplikasi Go, sehingga Anda dapat menggunakan driver basis data Go standar apa pun — termasuk [SQLite](https://github.com/mattn/go-sqlite3), [PostgreSQL](https://github.com/jackc/pgx), dan [MySQL](https://github.com/go-sql-driver/mysql) — di dalam layanan yang Anda ikat.

[Tutorial TODO](/tutorials/02-todo-vanilla/) mendemonstrasikan layanan terikat dengan penyimpanan dalam memori; tutorial tersebut tidak mengimplementasikan persistensi SQLite. Untuk persistensi berbasis berkas, lihat [tutorial Notes](/tutorials/03-notes-vanilla/). Pengelolaan koneksi, migrasi, dan transaksi khusus basis data tidak dibahas dalam halaman sementara ini.
