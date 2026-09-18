---
title: "Perbaiki dokumentasi"
description: "Kirim PR koreksi untuk dokumentasi Wails v3 menggunakan M-Press."
sourcePath: "contributing/documentation.md"
---

PR koreksi dipersilakan. Perbaiki kesalahan ketik, tautan rusak, contoh yang sudah usang,  
penjelasan yang tidak jelas, atau terjemahan. Anda tidak memerlukan issue atau pengujian  
kode yang gagal untuk koreksi yang hanya menyangkut dokumentasi.

## Pratinjau secara lokal

Fork [wailsapp/wails](https://github.com/wailsapp/wails/fork), kloning fork Anda,  
lalu buat branch dari `master`.

Instal generator dokumentasi dengan versi yang telah ditetapkan:

```sh
go install github.com/leaanthony/mpress/cmd/mpress@v1.0.17
```

Anda juga dapat mengunduh berkas biner terverifikasi dari  
[rilis M-Press v1.0.17](https://github.com/leaanthony/mpress/releases/tag/v1.0.17).

Dari direktori root repositori Wails:

```sh
mpress version
mpress dev
```

Edit file sumber `.md` di `docs/mpress/content/`. Bahasa Inggris adalah bahasa default  
dan berada langsung di direktori tersebut. Terjemahan yang sudah ada berada di folder  
bahasa seperti `fr/` dan `id/`. Pratinjau dibuat ulang saat Anda menyimpan perubahan.

Pertahankan blok metadata di bagian atas setiap halaman dan pasangan komponen `@...` / `@end`.  
Paragraf biasa, judul, daftar, dan kode berpagar dapat diedit sebagai teks.  
Jangan edit file yang dihasilkan di `docs/mpress/site/`.

## Periksa koreksi

```sh
python3 docs/mpress/scripts/check_translations.py
mpress build --strict --no-purge-css
mpress check
# Use `python` on Windows if `python3` is unavailable.
python3 docs/mpress/scripts/check_site.py docs/mpress/site
```

Periksa halaman yang diubah di browser dan jalankan setiap contoh kode yang Anda ubah.  
Untuk koreksi terjemahan, bandingkan seluruh bagian yang diubah dengan versi bahasa Inggris.  
Pertahankan perintah, nama API, tautan, sampel kode, dan koneksi diagram.  
Gunakan bahasa teknis yang alami, pertahankan persyaratan dan catatan penting, serta terjemahkan  
label diagram yang terlihat, label navigasi, dan deskripsi gambar, selain teks prosa.

Setiap bahasa yang dipublikasikan harus memiliki terjemahan lengkap untuk setiap halaman berbahasa Inggris.  
Jangan gunakan placeholder berbahasa Inggris atau halaman fallback. Koreksi pada satu terjemahan  
dapat hanya mengubah bahasa tersebut. Jika Anda mengubah makna teks bahasa Inggris, perbarui halaman yang sesuai  
dalam bahasa lain yang dipublikasikan; halaman yang tidak terkait tidak perlu dibuat ulang.

## Kirim pull request

Buka PR terhadap `master`. Jelaskan apa yang salah, uraikan koreksi Anda,  
dan cantumkan pemeriksaan yang Anda jalankan. Sertakan tangkapan layar untuk perubahan tata letak yang terlihat  
dan detail platform/versi untuk contoh kode yang diubah.

Kredensial Cloudflare dan layanan privat tidak diperlukan. Pemeriksaan PR publik  
membangun dan memvalidasi situs statis tanpa kredensial deployment.

Untuk perubahan kode dan proposal fitur, lihat [Berkontribusi pada Wails](/contributing/).  
Untuk informasi internal, lihat [Gambaran Umum Teknis](/contributing/overview/).
