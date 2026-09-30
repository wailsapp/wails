---
title: "Kira"
description: "Klien desktop native macOS yang menyatukan AWS — ECS, RDS, S3, DynamoDB, dan lainnya — dalam satu antarmuka berbasis papan ketik"
slug: "community/showcase/kira"
sourcePath: "community/showcase/kira.md"
---

**[Kira](https://kira.thiennguyen.dev)** adalah **klien desktop native macOS untuk AWS** — *"AWS, tanpa hambatan"*. Dibangun dengan **Go, Wails, dan React + TypeScript**, aplikasi ini menyatukan operasi AWS dalam satu aplikasi berbasis papan ketik. Lakukan autentikasi satu kali melalui AWS SSO, lalu kelola infrastruktur di banyak akun tanpa perlu beralih antara konsol, CLI, dan berbagai alat basis data.

Semuanya berawal dari kebutuhan pribadi: berpindah-pindah antara tab peramban, perintah `aws`, dan klien SQL terpisah hanya untuk merilis satu perubahan terasa lambat. Kira menyatukan semua alur kerja tersebut dalam satu jendela native yang cepat — masuk, pilih akun, dan semua yang Anda perlukan dapat diakses hanya dengan satu pintasan papan ketik.

![Ringkasan multiakun Kira — akun produksi, staging, dan pengembangan dalam satu tempat](/assets/showcase-images/kira_screenshot_1.png)

## Fitur Utama

- **AWS SSO multiakun** - Masuk satu kali, lalu beralih antara akun dan region dari satu tempat
- **ECS** - Jelajahi klaster, layanan, dan tugas; deploy ulang, skalakan, dan rollback layanan; lihat definisi tugas; pantau metrik layanan; serta buka shell interaktif melalui ECS Exec
- **Basis data** - Jalankan SQL pada RDS, lakukan kueri dan pemindaian pada DynamoDB, serta hubungkan ke PostgreSQL, MySQL, dan Redshift — dengan tunneling SSH yang aman dan kredensial yang disimpan di macOS Keychain
- **S3** - Telusuri bucket dan prefiks; pratinjau, unggah, unduh, salin, ganti nama, dan hapus objek; serta buat folder
- **Secrets Manager** - Tampilkan daftar rahasia dan ambil nilainya untuk akun aktif
- **CloudWatch Logs** - Pantau secara langsung dan cari aliran log
- **Smart Query** - Pembuatan SQL opsional dengan bantuan AI yang didukung oleh CLI `claude`
- **Ekstensi** - Instal bundel `.kext` khusus yang menambahkan tombol tindakan dengan dukungan skrip Go berukuran kecil
- **Navigasi cepat** - Tombol pintas global untuk memanggil aplikasi, palet perintah `Cmd+K`, dan deep linking `kira://`

## Melihat Lebih Dekat

Pantau layanan ECS Anda secara langsung — kondisi tugas, CPU dan memori, serta status deployment — lalu lakukan deployment ulang, penskalaan, atau rollback tanpa meninggalkan daftar.

![Kira menelusuri layanan ECS dengan metrik langsung](/assets/showcase-images/kira_screenshot_17.png)

Jelajahi S3 seperti menggunakan pengelola file. Pratinjau objek, periksa metadata dan versi, serta unggah, unduh, ganti nama, atau hapus langsung di tempat.

![Peramban objek S3 Kira dengan pratinjau objek dan metadata](/assets/showcase-images/kira_screenshot_5.png)

Didistribusikan sebagai `.dmg` untuk macOS yang ditandatangani dan dinotarisasi.

[Kunjungi Kira](https://kira.thiennguyen.dev) | [Baca dokumentasi](https://docs.kira.thiennguyen.dev)
