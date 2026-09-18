---
title: "Meningkatkan versi dari v3 alpha"
description: "Memindahkan proyek Wails v3 alpha yang sudah ada ke versi beta yang ditetapkan"
slug: "migration/alpha-to-beta"
sourcePath: "migration/alpha-to-beta.md"
---

Panduan ini untuk proyek v3 alpha yang sudah ada. Untuk Wails v2, gunakan [panduan v2 ke v3](/migration/v2-to-v3/).

## Sebelum meningkatkan versi

Buat commit atau cadangkan proyek Anda. Baca [catatan perubahan](/changelog/) antara versi alpha Anda dan beta yang dipilih: kode sumber, API, atau konfigurasi build mungkin perlu diubah. Periksa [kebijakan kompatibilitas desktop](/status/) dan persyaratan platform Anda.

Perintah di bawah menggunakan rilis `v3.0.0-beta.23` yang telah diterbitkan sebagai contoh versi pasti, bukan anjuran untuk terus mengikuti rilis terbaru. Jika memilih rilis lain, periksa versi CLI, modul Go, dan runtime npm-nya, lalu perbarui perintah secara bersamaan. Dalam contoh ini, versi npm sama dengan versi Go tanpa awalan `v`.

## 1. Perbarui CLI

```sh
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0-beta.23
wails3 version
```

Pastikan `wails3 version` menampilkan versi yang Anda instal. Berkas biner lama yang muncul lebih awal di `PATH` dapat menutupi CLI baru.

## 2. Perbarui modul Go

Jalankan dari direktori akar proyek. Tinjau perubahan dependensi; jangan meningkatkan semua modul yang tidak terkait secara menyeluruh.

```sh
go get github.com/wailsapp/wails/v3@v3.0.0-beta.23
go mod tidy
```

## 3. Perbarui runtime frontend

Untuk proyek yang menggunakan npm dan direktori `frontend`:

```sh
cd frontend
npm install --save-exact @wailsio/runtime@3.0.0-beta.23
cd ..
```

Pertahankan berkas penguncian dan tinjau perubahannya. Jika frontend menggunakan pengelola paket atau direktori lain, sesuaikan langkah ini sambil tetap menetapkan versi runtime yang pasti.

## 4. Buat ulang, build, dan uji

Dari direktori akar proyek, buat ulang binding dari layanan Go Anda dan jalankan build:

```sh
wails3 generate bindings
wails3 build
```

Jalankan aplikasi hasil build dan uji alur kerja pada setiap platform yang didukung dan menjadi target distribusi Anda. Tinjau dan commit bersama perubahan kode sumber, binding yang dihasilkan, berkas modul, dan berkas penguncian frontend.

## Jika peningkatan versi gagal

Periksa CLI pada `PATH`, versi modul dengan `go list -m github.com/wailsapp/wails/v3`, dan runtime terinstal dengan `npm --prefix frontend ls @wailsio/runtime`. Buat ulang binding setelah menyelesaikan ketidakcocokan versi. Jangan menganggap setiap alpha dapat ditingkatkan tanpa perubahan kode.

Jika masalah tetap terjadi, [laporkan issue yang dapat direproduksi](https://github.com/wailsapp/wails/issues/new/choose) beserta versi lama dan baru, pesan kesalahan yang tepat, serta keluaran `wails3 doctor`. Ikuti [kebijakan keamanan](https://github.com/wailsapp/wails/blob/master/SECURITY.md) untuk kerentanan.
