---
title: "Triase Pemelihara"
description: "Melakukan triase bug, laporan dokumentasi, dan WEP secara konsisten"
slug: "contributing/maintainer-triage"
sourcePath: "contributing/maintainer-triage.md"
---

## Tujuan

Issue mencatat bug yang dapat direproduksi dan masalah dokumentasi. Pull request WEP (Wails Enhancement Proposal) mencatat usulan kemampuan atau perubahan pada perilaku publik.

## Triase issue

- Pastikan laporan bug menyertakan versi rilis, platform, langkah reproduksi, perilaku yang diharapkan, perilaku aktual, dan output `wails3 doctor`.
- Beri label `Bug` pada laporan yang telah dikonfirmasi dan terapkan label versi serta platform yang relevan. Mintalah reproduksi minimal jika diperlukan.
- Biarkan laporan dokumentasi tetap terbuka jika laporan tersebut mengidentifikasi cacat dokumentasi yang konkret; anjurkan pelapor untuk membuat PR jika mampu melakukan perubahan tersebut.
- Arahkan permintaan fitur ke panduan WEP, lalu tutup permintaan tersebut. Alur kerja pengalihan otomatis menangani issue peningkatan yang baru diberi label; gunakan redaksi yang sama untuk issue lama.
- Pindahkan pertanyaan dan permintaan dukungan ke GitHub Discussions atau Discord.

## Triase WEP

1. Pastikan PR tersebut berupa draf berjudul `[WEP] <title>` dan hanya berisi WEP serta materi pendukung.
2. Pastikan PR tersebut menggunakan templat WEP, mengidentifikasi pelaksana, serta mencakup kompatibilitas, platform, pengujian, pemeliharaan, dan keamanan/privasi.
3. Pertahankan diskusi teknis di PR WEP. Discussions berguna sebagai konteks, tetapi bukan catatan keputusan.
4. Catat keputusan pemelihara dalam komentar PR: diterima, ditolak, atau ditarik, beserta alasan singkat.
5. Untuk WEP yang diterima, tetapkan nomornya, perbarui indeks WEP, gabungkan PR WEP, dan wajibkan PR implementasi untuk menautkan kembali ke WEP tersebut.

## Issue peningkatan yang sudah ada

Jangan menghapus issue peningkatan historis secara diam-diam. Untuk setiap permintaan yang masih relevan, tinggalkan komentar pengalihan dan tutup permintaan tersebut; kontributor yang berminat dapat membuka WEP. Tutup duplikat dengan tautan ke WEP atau keputusan yang sudah ada.
